package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type HealthGateStatus string

const (
	StatusMainHealthy    HealthGateStatus = "MAIN_HEALTHY"
	StatusMainUnhealthy  HealthGateStatus = "MAIN_UNHEALTHY"
	StatusMainUnverified HealthGateStatus = "MAIN_UNVERIFIED"

	ExitCodeMainHealthy    = 0
	ExitCodeMainUnhealthy  = 1
	ExitCodeMainUnverified = 2
)

type HealthDimensionResult struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Status  string   `json:"status"` // "PASS", "FAIL", "UNVERIFIED"
	Details []string `json:"details,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type MainlineHealthGateReport struct {
	Status     HealthGateStatus                 `json:"status"`
	ExitCode   int                              `json:"exit_code"`
	Summary    string                           `json:"summary"`
	Timestamp  string                           `json:"timestamp"`
	Root       string                           `json:"root"`
	HeadCommit string                           `json:"head_commit"`
	Dimensions map[string]HealthDimensionResult `json:"dimensions"`
}

func validateCodeAudit(doc DoctorOutput, cat CheckCatalogSummary) []string {
	var errs []string
	if !doc.Healthy {
		errs = append(errs, doc.Warnings...)
		return errs
	}
	if len(doc.Checks) == 0 {
		errs = append(errs, "doctor verified 0 checks")
		return errs
	}
	if cat.Checks == 0 || cat.Mandatory == 0 {
		errs = append(errs, "check catalog has 0 registered or mandatory checks")
		return errs
	}
	return nil
}

func evaluateCodeDimension(root string, state CurrentState, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "code", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	doc, err := doctor(root, state, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("doctor validation failed: %v", err))
		return res
	}

	cat, err := runCheckCatalogAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("check catalog audit failed: %v", err))
		return res
	}

	if errs := validateCodeAudit(doc, cat); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	testedModules, unverifiedErr, err := runRepositoryGoUnitTests(root)
	if unverifiedErr != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, unverifiedErr.Error())
		return res
	}
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res
	}

	testedNodeSuites, nodeUnverifiedErr, err := runRepositoryNodeUnitTests(root)
	if nodeUnverifiedErr != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, nodeUnverifiedErr.Error())
		return res
	}
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("doctor checks passed (%d checks verified)", len(doc.Checks)),
		fmt.Sprintf("check catalog verified (%d checks registered, %d mandatory)", cat.Checks, cat.Mandatory),
		fmt.Sprintf("Go unit tests passed across %d modules (product-runtime + tools/cmdr-dev)", len(testedModules)),
	)
	if len(testedNodeSuites) > 0 {
		res.Details = append(res.Details, fmt.Sprintf("Node unit tests passed across %d suite(s) under product-runtime", len(testedNodeSuites)))
	}
	return res
}

func discoverCodeTestModules(root string) ([]string, error) {
	prInfo, err := statRepoPath(root, "product-runtime")
	if err != nil {
		return nil, fmt.Errorf("inspect product-runtime directory: %w", err)
	}
	if !prInfo.IsDir() {
		return nil, fmt.Errorf("product-runtime is not a directory")
	}
	prDir, err := resolveRepoPath(root, "product-runtime", false)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(prDir)
	if err != nil {
		return nil, fmt.Errorf("read product-runtime directory: %w", err)
	}
	var runtimeModules []string
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("product-runtime entry %s must not be a symlink", entry.Name())
		}
		if !entry.IsDir() {
			continue
		}
		modRel := filepath.ToSlash(filepath.Join("product-runtime", entry.Name(), "go.mod"))
		info, err := lstatRepoEntry(root, modRel)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("inspect module manifest %s: %w", modRel, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("module manifest %s must not be a symlink", modRel)
		}
		if !info.IsDir() {
			runtimeModules = append(runtimeModules, filepath.ToSlash(filepath.Join("product-runtime", entry.Name())))
		}
	}
	if len(runtimeModules) == 0 {
		return nil, fmt.Errorf("no Go modules discovered under product-runtime")
	}
	sort.Strings(runtimeModules)

	cmdrDevModInfo, err := lstatRepoEntry(root, "tools/cmdr-dev/go.mod")
	if err != nil {
		return nil, fmt.Errorf("inspect tools/cmdr-dev/go.mod: %w", err)
	}
	if cmdrDevModInfo.Mode()&os.ModeSymlink != 0 || cmdrDevModInfo.IsDir() {
		return nil, fmt.Errorf("tools/cmdr-dev/go.mod must be a regular file")
	}

	modules := make([]string, 0, 1+len(runtimeModules))
	modules = append(modules, "tools/cmdr-dev")
	modules = append(modules, runtimeModules...)
	return modules, nil
}

func runRepositoryGoUnitTests(root string) ([]string, error, error) {
	modules, err := discoverCodeTestModules(root)
	if err != nil {
		return nil, nil, err
	}
	if _, lookErr := exec.LookPath("go"); lookErr != nil {
		if errors.Is(lookErr, exec.ErrNotFound) {
			return nil, fmt.Errorf("go toolchain unavailable for unit test verification: %w", lookErr), nil
		}
		return nil, fmt.Errorf("resolve go toolchain: %w", lookErr), nil
	}

	workers := 4
	if len(modules) < workers {
		workers = len(modules)
	}
	jobs := make(chan string, len(modules))
	for _, m := range modules {
		jobs <- m
	}
	close(jobs)

	var mu sync.Mutex
	var failures []string
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for modRel := range jobs {
				modDir, pathErr := resolveRepoPath(root, modRel, false)
				if pathErr != nil {
					mu.Lock()
					failures = append(failures, fmt.Sprintf("%s: %v", modRel, pathErr))
					mu.Unlock()
					continue
				}
				cmd := exec.Command("go", "test", "-count=1", "./...") // #nosec G204 -- fixed Go executable and literal flags; module directory is repository-confined.
				cmd.Dir = modDir
				out, runErr := cmd.CombinedOutput()
				if runErr != nil {
					snippet := strings.Join(strings.Fields(string(out)), " ")
					if len(snippet) > 240 {
						snippet = snippet[len(snippet)-240:]
					}
					mu.Lock()
					failures = append(failures, fmt.Sprintf("%s (%v): %s", modRel, runErr, snippet))
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	if len(failures) > 0 {
		sort.Strings(failures)
		return modules, nil, fmt.Errorf("Go unit tests failed in %d module(s): %s", len(failures), strings.Join(failures, "; "))
	}
	return modules, nil, nil
}

type nodeTestSuite struct {
	Directory string
	Files     []string
}

func isNodeTestFile(name string) bool {
	return strings.HasSuffix(name, ".test.mjs") || strings.HasSuffix(name, ".test.cjs") || strings.HasSuffix(name, ".test.js")
}

func discoverCodeNodeTestSuites(root string) ([]nodeTestSuite, error) {
	prInfo, err := statRepoPath(root, "product-runtime")
	if err != nil {
		return nil, fmt.Errorf("inspect product-runtime directory: %w", err)
	}
	if !prInfo.IsDir() {
		return nil, fmt.Errorf("product-runtime is not a directory")
	}
	prDir, err := resolveRepoPath(root, "product-runtime", false)
	if err != nil {
		return nil, err
	}
	byDir := map[string][]string{}
	// #nosec G703 -- prDir is repository-confined by resolveRepoPath; symlink entries are rejected below.
	err = filepath.WalkDir(prDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("product-runtime Node path must not be a symlink: %s", path)
		}
		if entry.IsDir() {
			if path != prDir && (ignoredDependencyDir(entry.Name()) || entry.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !isNodeTestFile(entry.Name()) {
			return nil
		}
		if strings.HasPrefix(entry.Name(), "-") {
			return fmt.Errorf("product-runtime Node test filename must not start with '-': %s", entry.Name())
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		dir := filepath.ToSlash(filepath.Dir(rel))
		byDir[dir] = append(byDir[dir], entry.Name())
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(byDir) == 0 {
		return nil, nil
	}
	dirs := make([]string, 0, len(byDir))
	for dir := range byDir {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	suites := make([]nodeTestSuite, 0, len(dirs))
	for _, dir := range dirs {
		files := byDir[dir]
		sort.Strings(files)
		suites = append(suites, nodeTestSuite{Directory: dir, Files: files})
	}
	return suites, nil
}

func runRepositoryNodeUnitTests(root string) ([]nodeTestSuite, error, error) {
	suites, err := discoverCodeNodeTestSuites(root)
	if err != nil {
		return nil, nil, err
	}
	if len(suites) == 0 {
		return nil, nil, nil
	}
	if _, lookErr := exec.LookPath("node"); lookErr != nil {
		if errors.Is(lookErr, exec.ErrNotFound) {
			return nil, fmt.Errorf("node runtime unavailable for unit test verification: %w", lookErr), nil
		}
		return nil, fmt.Errorf("resolve node runtime: %w", lookErr), nil
	}

	var mu sync.Mutex
	var failures []string
	var wg sync.WaitGroup
	for _, suite := range suites {
		suite := suite
		wg.Add(1)
		go func() {
			defer wg.Done()
			suiteDir, pathErr := resolveRepoPath(root, suite.Directory, false)
			if pathErr != nil {
				mu.Lock()
				failures = append(failures, fmt.Sprintf("%s: %v", suite.Directory, pathErr))
				mu.Unlock()
				return
			}
			args := make([]string, 0, 1+len(suite.Files))
			args = append(args, "--test")
			for _, file := range suite.Files {
				if strings.HasPrefix(file, "-") {
					mu.Lock()
					failures = append(failures, fmt.Sprintf("%s: Node test filename must not start with '-': %s", suite.Directory, file))
					mu.Unlock()
					return
				}
				args = append(args, "./"+file)
			}
			cmd := exec.Command("node", args...) // #nosec G204,G702 -- fixed node executable and --test flag; suite directory and discovered test files are repository-confined.
			cmd.Dir = suiteDir
			out, runErr := cmd.CombinedOutput()
			if runErr != nil {
				snippet := strings.Join(strings.Fields(string(out)), " ")
				if len(snippet) > 240 {
					snippet = snippet[len(snippet)-240:]
				}
				mu.Lock()
				failures = append(failures, fmt.Sprintf("%s (%v): %s", suite.Directory, runErr, snippet))
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(failures) > 0 {
		sort.Strings(failures)
		return suites, nil, fmt.Errorf("Node unit tests failed in %d suite(s): %s", len(failures), strings.Join(failures, "; "))
	}
	return suites, nil, nil
}

func validateSecurityAudit(secScan SecretScanSummary, gates SecurityGateAuditSummary, goVersion string) []string {
	var errs []string
	if secScan.FindingCount > 0 {
		errs = append(errs, fmt.Sprintf("%d secret findings detected across %d paths", secScan.FindingCount, secScan.CandidatePaths))
		return errs
	}
	if secScan.CandidatePaths == 0 {
		errs = append(errs, "secret scan inspected 0 candidate paths")
		return errs
	}
	if gates.Gates == 0 || gates.Blocking == 0 {
		errs = append(errs, "security gate registry has 0 registered or blocking gates")
		return errs
	}
	if goVersion != "1.26.11" {
		errs = append(errs, fmt.Sprintf("engineering Go toolchain version mismatch: expected 1.26.11, got %s", goVersion))
		return errs
	}
	return nil
}

func validateDependencyAudit(deps DependencyAuditSummary) []string {
	var errs []string
	if deps.UnsupportedManifests > 0 {
		errs = append(errs, fmt.Sprintf("%d unsupported product-runtime dependency manifest(s) detected", deps.UnsupportedManifests))
	}
	if deps.Unapproved > 0 {
		errs = append(errs, fmt.Sprintf("%d unapproved third-party product-runtime dependencies detected", deps.Unapproved))
	}
	if deps.RuntimeBoundaries == 0 {
		errs = append(errs, "dependency audit verified 0 runtime boundaries")
	}
	return errs
}

func validateBoundaryEdgeAudit(edges BoundaryEdgeAuditSummary) []string {
	var errs []string
	if edges.RuntimeBoundaries == 0 {
		errs = append(errs, "boundary edge audit verified 0 runtime boundaries")
	}
	if edges.LocalPackages == 0 {
		errs = append(errs, "boundary edge audit discovered 0 local packages")
	}
	return errs
}

func evaluateSecurityDimension(root string) HealthDimensionResult {
	res := HealthDimensionResult{Name: "security", Passed: false, Status: "FAIL"}

	secScan, err := runSecretScan(root, "", true)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("secret scan failed: %v", err))
		return res
	}

	gates, err := runSecurityGateAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("security gates audit failed: %v", err))
		return res
	}

	if errs := validateSecurityAudit(secScan, gates, engineeringGoVersion); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	deps, err := runDependencyAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("dependency audit failed: %v", err))
		return res
	}
	if errs := validateDependencyAudit(deps); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	edges, err := runBoundaryEdgeAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("boundary edge audit failed: %v", err))
		return res
	}
	if errs := validateBoundaryEdgeAudit(edges); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	secEv, err := runCommittedRuntimeSecurityEvidenceAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("runtime security evidence audit failed: %v", err))
		return res
	}

	evidenceDetail := fmt.Sprintf("runtime security evidence valid (status=%s, scopes=%d)", secEv.CoverageStatus, secEv.RegisteredScopes)
	if secEv.GlobalCoveragePercent != nil {
		evidenceDetail = fmt.Sprintf("runtime security evidence valid (status=%s, global_coverage=%.2f%%, scopes=%d)", secEv.CoverageStatus, *secEv.GlobalCoveragePercent, secEv.RegisteredScopes)
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("0 secrets detected across %d candidate paths", secScan.CandidatePaths),
		fmt.Sprintf("security gate registry valid (%d registered gates, %d blocking)", gates.Gates, gates.Blocking),
		fmt.Sprintf("Go toolchain pinned to %s", engineeringGoVersion),
		fmt.Sprintf("runtime dependencies verified (0 unapproved across %d boundaries)", deps.RuntimeBoundaries),
		fmt.Sprintf("architecture boundary edges verified (%d observed, %d cross-boundary across %d boundaries)", edges.ObservedEdges, edges.CrossBoundaryEdges, edges.RuntimeBoundaries),
		evidenceDetail,
	)
	return res
}

func evaluateGitDimension(root string) (HealthDimensionResult, string) {
	res := HealthDimensionResult{Name: "git", Passed: false, Status: "FAIL"}

	cmdHead := exec.Command("git", "-C", root, "rev-parse", "HEAD") // #nosec G204,G702 -- git rev-parse HEAD within confined root
	headBytes, err := cmdHead.Output()
	if err != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, fmt.Sprintf("unable to resolve HEAD commit: %v", err))
		return res, ""
	}
	headCommit := strings.TrimSpace(string(headBytes))
	if !fullCommitIDPattern.MatchString(headCommit) {
		res.Errors = append(res.Errors, fmt.Sprintf("invalid HEAD commit SHA: %q", headCommit))
		return res, headCommit
	}

	cmdStatus := exec.Command("git", "-C", root, "status", "--porcelain") // #nosec G204,G702 -- git status within confined root
	statusBytes, err := cmdStatus.Output()
	if err != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, fmt.Sprintf("git status check failed: %v", err))
		return res, headCommit
	}
	statusStr := strings.TrimSpace(string(statusBytes))
	if len(statusStr) > 0 {
		res.Errors = append(res.Errors, fmt.Sprintf("working tree has uncommitted modifications: %s", statusStr))
		return res, headCommit
	}

	const baselineCommit = "82e630fce3a7f6236e03709116eca6dd0a1acd40"
	if gitCommitExists(root, baselineCommit) || ensureGitCommit(root, baselineCommit) == nil {
		cmdDiff := exec.Command("git", "-C", root, "diff", baselineCommit, "--", "cmdr-product-spec") // #nosec G204,G702 -- invariant diff check
		diffBytes, err := cmdDiff.Output()
		if err != nil {
			res.Status = "UNVERIFIED"
			res.Errors = append(res.Errors, fmt.Sprintf("git diff product spec baseline failed: %v", err))
			return res, headCommit
		}
		if len(bytes.TrimSpace(diffBytes)) > 0 {
			res.Errors = append(res.Errors, "cmdr-product-spec has non-zero diff vs baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40")
			return res, headCommit
		}
	} else {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, "product spec baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40 not available in git history")
		return res, headCommit
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("HEAD commit resolved: %s", headCommit),
		"working tree clean (0 uncommitted changes)",
		"cmdr-product-spec 0-byte invariant verified against baseline 82e630fce3a7f6236e03709116eca6dd0a1acd40",
	)
	return res, headCommit
}

func validateTraceabilityAudit(ledger ImplementationLedgerAuditSummary, dec DecisionRegistryAuditSummary) []string {
	var errs []string
	if ledger.Status != "valid" && ledger.Status != "PASS" {
		errs = append(errs, fmt.Sprintf("implementation ledger status is %s", ledger.Status))
		return errs
	}
	if ledger.Entries == 0 || ledger.VerifiedClaims == 0 {
		errs = append(errs, "implementation ledger has 0 entries or verified claims")
		return errs
	}
	if ledger.VerifiedClaims != ledger.Entries {
		errs = append(errs, fmt.Sprintf("implementation ledger accounting mismatch: %d claims vs %d entries", ledger.VerifiedClaims, ledger.Entries))
		return errs
	}
	if dec.Decisions == 0 {
		errs = append(errs, "decision registry has 0 decisions registered")
		return errs
	}
	return nil
}

func evaluateTraceabilityDimension(root string, state CurrentState, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "traceability", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	ledger, err := runImplementationLedgerAudit(root, state, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("implementation ledger audit failed: %v", err))
		return res
	}

	dec, err := runDecisionRegistryAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("decision registry audit failed: %v", err))
		return res
	}

	if errs := validateTraceabilityAudit(ledger, dec); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("implementation ledger %d entries verified against spec baseline %s", ledger.Entries, ledger.ProductSpecBaseline),
		fmt.Sprintf("decision registry %d decisions reconciled with work graph", dec.Decisions),
	)
	return res
}

func validateOperationalQualityAudit(comp ComplexityAuditSummary, perf PerformanceRegistryAuditSummary, rec RecoveryJournalAuditSummary) []string {
	var errs []string
	if comp.ReadinessViolations > 0 {
		errs = append(errs, fmt.Sprintf("%d complexity readiness violations detected", comp.ReadinessViolations))
		return errs
	}
	if comp.Evaluated == 0 {
		errs = append(errs, "complexity audit evaluated 0 units")
		return errs
	}
	if perf.RuntimeCoverageStatus != "complete" {
		errs = append(errs, fmt.Sprintf("performance registry runtime coverage status is %s", perf.RuntimeCoverageStatus))
		return errs
	}
	if perf.Targets == 0 || perf.Metrics == 0 {
		errs = append(errs, "performance registry has 0 targets or metrics")
		return errs
	}
	if rec.Status != "valid" && rec.Status != "empty" {
		errs = append(errs, fmt.Sprintf("recovery journal status is %s", rec.Status))
		return errs
	}
	return nil
}

func evaluateOperationalQualityDimension(root string, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "operational_quality", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	comp, err := runComplexityAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("complexity audit failed: %v", err))
		return res
	}

	perf, err := runPerformanceRegistryAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("performance registry audit failed: %v", err))
		return res
	}

	rec, err := runRecoveryJournalAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("recovery journal audit failed: %v", err))
		return res
	}

	if errs := validateOperationalQualityAudit(comp, perf, rec); len(errs) > 0 {
		res.Errors = append(res.Errors, errs...)
		return res
	}

	journalDetail := fmt.Sprintf("recovery journal valid with %d events intact", rec.Events)
	if rec.Status == "empty" {
		journalDetail = "recovery journal quiescent (0 active incidents/events)"
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("complexity within limits across %d evaluated units (%d within budget)", comp.Evaluated, comp.WithinBudget),
		fmt.Sprintf("performance registry complete with %d targets and %d metrics", perf.Targets, perf.Metrics),
		journalDetail,
	)
	return res
}

func runMainlineHealthGate(root string, state CurrentState, graph WorkGraph, stateErr error) (MainlineHealthGateReport, int) {
	report := MainlineHealthGateReport{
		Root:       root,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Dimensions: make(map[string]HealthDimensionResult, 5),
	}

	codeDim := evaluateCodeDimension(root, state, graph, stateErr)
	report.Dimensions["code"] = codeDim

	secDim := evaluateSecurityDimension(root)
	report.Dimensions["security"] = secDim

	gitDim, headCommit := evaluateGitDimension(root)
	report.Dimensions["git"] = gitDim
	report.HeadCommit = headCommit

	traceDim := evaluateTraceabilityDimension(root, state, graph, stateErr)
	report.Dimensions["traceability"] = traceDim

	opDim := evaluateOperationalQualityDimension(root, graph, stateErr)
	report.Dimensions["operational_quality"] = opDim

	aggregateHealthGateReport(&report)

	return report, report.ExitCode
}

// aggregateHealthGateReport aggregates the 5 required dimensions into the final gate status,
// exit code, and human-readable summary. It enforces that all 5 dimensions must exist,
// be marked passed, and have no FAIL status. If any dimension is UNVERIFIED (and none FAIL),
// the status is MAIN_UNVERIFIED.
func aggregateHealthGateReport(report *MainlineHealthGateReport) (HealthGateStatus, int) {
	if report == nil {
		return StatusMainUnhealthy, ExitCodeMainUnhealthy
	}

	hasFail := false
	hasUnverified := false

	requiredDimensions := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, req := range requiredDimensions {
		dim, ok := report.Dimensions[req]
		if !ok || dim.Status == "FAIL" || (!dim.Passed && dim.Status != "UNVERIFIED") {
			hasFail = true
		}
		if ok && dim.Status == "UNVERIFIED" {
			hasUnverified = true
		}
	}

	if hasFail {
		report.Status = StatusMainUnhealthy
		report.ExitCode = ExitCodeMainUnhealthy
		report.Summary = "Mainline health gate FAILED; one or more required health dimensions failed verification."
	} else if hasUnverified {
		report.Status = StatusMainUnverified
		report.ExitCode = ExitCodeMainUnverified
		report.Summary = "Mainline health gate UNVERIFIED due to unresolved environmental or diagnostic conditions."
	} else {
		report.Status = StatusMainHealthy
		report.ExitCode = ExitCodeMainHealthy
		report.Summary = "Mainline health gate PASSED across all 5 dimensions (Code, Security, Git, Traceability, Operational Quality)."
	}

	return report.Status, report.ExitCode
}

func printMainlineHealthGateReport(report MainlineHealthGateReport) {
	fmt.Printf("=== MAINLINE HEALTH GATE REPORT ===\n")
	fmt.Printf("Status:      %s (exit %d)\n", report.Status, report.ExitCode)
	fmt.Printf("Summary:     %s\n", report.Summary)
	fmt.Printf("Timestamp:   %s\n", report.Timestamp)
	fmt.Printf("Root:        %s\n", report.Root)
	fmt.Printf("Head Commit: %s\n", report.HeadCommit)
	fmt.Println("-----------------------------------")
	dimOrder := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, name := range dimOrder {
		dim, ok := report.Dimensions[name]
		if !ok {
			continue
		}
		statusTag := fmt.Sprintf("[%s]", dim.Status)
		fmt.Printf("%-12s %-20s (passed=%t)\n", statusTag, dim.Name, dim.Passed)
		for _, d := range dim.Details {
			fmt.Printf("    ✓ %s\n", d)
		}
		for _, e := range dim.Errors {
			fmt.Printf("    ✗ %s\n", e)
		}
	}
	fmt.Println("===================================")
}
