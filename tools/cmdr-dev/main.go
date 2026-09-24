package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	statePath = "engineering/state/current.json"
	graphPath = "work/graph.json"
)

type CurrentState struct {
	SchemaVersion int    `json:"schema_version"`
	StateKind     string `json:"state_kind"`
	ProductSpec   struct {
		CanonicalPath                string `json:"canonical_path"`
		BaselineCommit               string `json:"baseline_commit"`
		DocumentaryStatus            string `json:"documentary_status"`
		RuntimeImplementationClaimed bool   `json:"runtime_implementation_claimed"`
		KnownCapabilities            int    `json:"known_capabilities"`
		KnownRequirements            int    `json:"known_requirements"`
		KnownOpenDecisions           int    `json:"known_open_decisions"`
		KnownActiveScreens           int    `json:"known_active_screens"`
	} `json:"product_spec"`
	EngineeringFoundation struct {
		Status     string `json:"status"`
		Phase      string `json:"phase"`
		PhaseTitle string `json:"phase_title"`
	} `json:"engineering_foundation"`
	Execution struct {
		ActiveWorkUnit    string `json:"active_work_unit"`
		PreferredNextUnit string `json:"preferred_next_unit"`
		SelectionRule     string `json:"selection_rule"`
	} `json:"execution"`
	Reconciliation struct {
		RequiredOnResume             bool `json:"required_on_resume"`
		RecordedStateIsAuthoritative bool `json:"recorded_state_is_authoritative"`
	} `json:"reconciliation"`
}

type WorkGraph struct {
	SchemaVersion int        `json:"schema_version"`
	GraphKind     string     `json:"graph_kind"`
	Nodes         []WorkNode `json:"nodes"`
}

type WorkNode struct {
	ID                      string   `json:"id"`
	Title                   string   `json:"title"`
	Type                    string   `json:"type"`
	Status                  string   `json:"status"`
	DependsOn               []string `json:"depends_on"`
	DependencyTerminalState string   `json:"dependency_terminal_state,omitempty"`
	Unlocks                 []string `json:"unlocks"`
}

type StatusOutput struct {
	Root             string `json:"root"`
	Phase            string `json:"phase"`
	PhaseTitle       string `json:"phase_title"`
	FoundationStatus string `json:"foundation_status"`
	ActiveWorkUnit   string `json:"active_work_unit"`
	PreferredNext    string `json:"preferred_next_unit"`
	ComputedNext     string `json:"computed_next_unit,omitempty"`
	ProductBaseline  string `json:"product_spec_baseline"`
}

type DoctorOutput struct {
	Root     string   `json:"root"`
	Healthy  bool     `json:"healthy"`
	Checks   []string `json:"checks"`
	Warnings []string `json:"warnings,omitempty"`
}

var lifecycleRank = map[string]int{
	"DRAFT":             0,
	"DECOMPOSED":        1,
	"READY":             2,
	"CLAIMED":           3,
	"IMPLEMENTING":      4,
	"IMPLEMENTED":       5,
	"TESTING":           6,
	"REVIEWING":         7,
	"SECURITY_REVIEW":   8,
	"MERGE_READY":       9,
	"MERGED":            10,
	"POST_MERGE_VERIFY": 11,
	"VERIFIED":          12,
	"BLOCKED":           -1,
	"BLOCKED_DECISION":  -1,
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	command := os.Args[1]
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	rootFlag := fs.String("root", "", "repository root")
	jsonFlag := fs.Bool("json", false, "machine-readable output")
	checkFlag := fs.Bool("check", false, "check generated output instead of writing it")
	outputFlag := fs.String("output", "engineering/spec-index/inventory.json", "generated output path")
	workUnitFlag := fs.String("work-unit", "", "work unit id; defaults to active work unit")
	packetIDFlag := fs.String("packet-id", "", "research packet id")
	decisionIDFlag := fs.String("decision-id", "", "engineering decision id")
	asOfDateFlag := fs.String("as-of-date", "", "freshness evaluation date YYYY-MM-DD; defaults to current UTC date")
	asOfTimeFlag := fs.String("as-of-time", "", "lease evaluation time RFC3339 UTC; defaults to current UTC time")
	leaseActionFlag := fs.String("lease-action", "", "lease action: acquire, renew, or release")
	leaseModeFlag := fs.String("lease-mode", "", "lease mode for acquire: mutating or read-only; blank defaults to mutating")
	leaseIDFlag := fs.String("lease-id", "", "opaque work lease id")
	agentIDFlag := fs.String("agent-id", "", "opaque non-secret agent id")
	leaseBaseHeadFlag := fs.String("lease-base-head", "", "full Git commit id used as lease base")
	leaseDurationFlag := fs.Int("lease-duration-minutes", 0, "lease duration; 0 uses policy default")
	ciHeadFlag := fs.String("ci-head", "", "full commit SHA for observed CI run")
	ciRunIDFlag := fs.String("ci-run-id", "", "positive source CI run id")
	ciEventFlag := fs.String("ci-event", "", "source CI event: push or pull_request")
	ciConclusionFlag := fs.String("ci-conclusion", "", "source CI conclusion")
	ciObservedAtFlag := fs.String("ci-observed-at", "", "source CI observation time RFC3339 UTC")
	changesFileFlag := fs.String("changes-file", "", "newline-delimited repository-relative changed paths")
	baseCommitFlag := fs.String("base-commit", "", "full Git base commit id")
	headCommitFlag := fs.String("head-commit", "", "full Git head commit id")
	fullScanFlag := fs.Bool("full-scan", false, "scan all Git-tracked repository files")
	stageFlag := fs.String("stage", "pr", "execution stage: pr, nightly, release, or on-demand")
	requestedGateFlag := fs.String("requested-gate", "", "optional exact deep security gate id for on-demand execution")
	requestedPerformanceTargetFlag := fs.String("requested-performance-target", "", "optional exact PERF-DEEP-* target for on-demand execution")
	performanceEnvironmentFlag := fs.String("performance-environment", defaultCIEnvironment, "performance benchmark environment id")
	if err := fs.Parse(os.Args[2:]); err != nil {
		fail(err)
	}

	root, err := resolveRoot(*rootFlag)
	if err != nil {
		fail(err)
	}

	state, graph, err := loadRepositoryState(root)
	if err != nil {
		fail(err)
	}

	switch command {
	case "doctor":
		out, err := doctor(root, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(out, *jsonFlag)
	case "status":
		next, _ := selectNext(graph)
		out := StatusOutput{
			Root:             root,
			Phase:            state.EngineeringFoundation.Phase,
			PhaseTitle:       state.EngineeringFoundation.PhaseTitle,
			FoundationStatus: state.EngineeringFoundation.Status,
			ActiveWorkUnit:   state.Execution.ActiveWorkUnit,
			PreferredNext:    state.Execution.PreferredNextUnit,
			ProductBaseline:  state.ProductSpec.BaselineCommit,
		}
		if next != nil {
			out.ComputedNext = next.ID
		}
		printValue(out, *jsonFlag)
	case "next":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		next, err := selectNext(graph)
		if err != nil {
			fail(err)
		}
		if next == nil {
			fail(errors.New("no executable READY work unit"))
		}
		printValue(next, *jsonFlag)
	case "spec-index":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runSpecIndex(root, state.ProductSpec.CanonicalPath, *outputFlag, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "spec-baseline":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/spec-index/baseline.json"
		}
		summary, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "coverage-graph":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/coverage/product-graph.json"
		}
		summary, err := runCoverageGraph(root, state.ProductSpec.CanonicalPath, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "obligations":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/coverage/obligations.json"
		}
		summary, err := runObligations(root, state.ProductSpec.CanonicalPath, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		if err := validateKnownObligationCounts(summary, state); err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "coverage-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/coverage/audit.json"
		}
		summary, err := runCoverageAudit(root, state.ProductSpec.CanonicalPath, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		if err := validateKnownCoverageAuditCounts(summary, state); err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "validate-manifests":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runManifestValidation(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "complexity-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runComplexityAudit(root, graph)
		if err != nil {
			fail(err)
		}
		if err := validateComplexityReadiness(summary); err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "context":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runContextCompiler(root, *workUnitFlag, *outputFlag, *checkFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "architecture-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runArchitectureAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "dependency-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDependencyAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "boundary-edge-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runBoundaryEdgeAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "check-catalog-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runCheckCatalogAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "impact":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		report, err := runImpactAnalysis(root, *workUnitFlag, *changesFileFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(report, *jsonFlag)
	case "validation-plan":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		plan, err := runValidationPlan(root, *workUnitFlag, *changesFileFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(plan, *jsonFlag)
	case "security-gate-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runSecurityGateAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "security-test-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runSecurityTestAudit(root, *changesFileFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "deep-security-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDeepSecurityAudit(root, *stageFlag, *changesFileFlag, *requestedGateFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "decision-registry-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDecisionRegistryAudit(root, graph)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "research-packet-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runResearchPacketAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "decision-gate-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDecisionGateAudit(root)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "decision-freshness-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDecisionFreshnessAudit(root, *asOfDateFlag, state, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "performance-registry-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPerformanceRegistryAudit(root)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "performance-benchmark-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPerformanceBenchmarkAudit(root, *stageFlag, *performanceEnvironmentFlag, *changesFileFlag)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "deep-performance-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runDeepPerformanceAudit(root, *stageFlag, *changesFileFlag, *requestedPerformanceTargetFlag)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "performance-cache-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPerformanceCacheAudit(root, state, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "lease-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runWorkLeaseAudit(root, *asOfTimeFlag, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "lease-evaluate":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		workUnit := *workUnitFlag
		if workUnit == "" {
			workUnit = state.Execution.ActiveWorkUnit
		}
		result, err := runWorkLeaseEvaluation(root, WorkLeaseActionRequest{
			Action: *leaseActionFlag, LeaseID: *leaseIDFlag, AgentID: *agentIDFlag, Mode: *leaseModeFlag,
			WorkUnit: workUnit, BaseHeadSHA: *leaseBaseHeadFlag,
			AsOf: *asOfTimeFlag, DurationMinutes: *leaseDurationFlag,
		}, graph)
		if err != nil {
			fail(err)
		}
		printValue(result, *jsonFlag)
	case "recovery-journal-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runRecoveryJournalAudit(root, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "resume-checkpoint":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		checkpoint, err := runResumeCheckpoint(root, *workUnitFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(checkpoint, *jsonFlag)
	case "recovery-reconcile":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		ci, err := ciObservationFromFlags(*ciHeadFlag, *ciRunIDFlag, *ciEventFlag, *ciConclusionFlag, *ciObservedAtFlag)
		if err != nil {
			fail(err)
		}
		result, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{
			WorkUnit: *workUnitFlag, AsOf: *asOfTimeFlag, CI: ci,
		}, state, graph)
		printValue(result, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "recovery-reconcile-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		result, err := runRecoveryReconciliationAudit(root, state, graph)
		printValue(result, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "coordination-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runCoordinationAudit(root, *asOfTimeFlag, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "metrics-registry-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runEngineMetricsRegistryAudit(root)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "coordination-handoff":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		handoff, err := runCoordinationHandoff(root, *workUnitFlag, *leaseIDFlag, *agentIDFlag, *asOfTimeFlag, state, graph)
		printValue(handoff, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "decision-cache":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		cache, err := runReusableEvidenceCache(root, *asOfDateFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(cache, *jsonFlag)
	case "decision-freshness-snapshot":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		snapshot, err := runDecisionFreshnessSnapshot(root, *decisionIDFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(snapshot, *jsonFlag)
	case "research-context":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		bundle, err := compileResearchContext(root, *packetIDFlag)
		if err != nil {
			fail(err)
		}
		printValue(bundle, *jsonFlag)
	case "secret-scan":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runSecretScan(root, *changesFileFlag, *fullScanFlag)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "sast-go":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runGoSAST(root)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "sca-go":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runGoSCA(root)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "sbom":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = defaultSBOMOutputPath
		}
		summary, err := runSBOM(root, output)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "pilot-scope-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPilotScopeAudit(root, state)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "pilot-contract-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPilotContractAudit(root)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "pilot-e2e-audit":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runPilotE2EAudit(root, state, graph)
		printValue(summary, *jsonFlag)
		if err != nil {
			fail(err)
		}
	case "git-changes":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runGitChanges(root, *baseCommitFlag, *headCommitFlag, *outputFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "validation-run":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runValidationExecution(root, *workUnitFlag, *changesFileFlag, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	default:
		usage(os.Stderr)
		fail(fmt.Errorf("unknown command %q", command))
	}
}

func resolveRoot(explicit string) (string, error) {
	if explicit != "" {
		abs, err := filepath.Abs(explicit)
		if err != nil {
			return "", err
		}
		if isRepoRoot(abs) {
			return abs, nil
		}
		return "", fmt.Errorf("%s is not a CMDR repository root", abs)
	}

	if env := os.Getenv("CMDR_ROOT"); env != "" {
		return resolveRoot(env)
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if isRepoRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("CMDR repository root not found")
		}
		dir = parent
	}
}

func isRepoRoot(root string) bool {
	required := []string{"AGENTS.md", "AI_START_HERE.md", statePath, graphPath}
	for _, rel := range required {
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

func loadRepositoryState(root string) (CurrentState, WorkGraph, error) {
	var state CurrentState
	var graph WorkGraph
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(statePath)), &state); err != nil {
		return state, graph, err
	}
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(graphPath)), &graph); err != nil {
		return state, graph, err
	}
	return state, graph, nil
}

func decodeStrict(root, path string, dst any) error {
	f, err := openRepoFile(root, path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode %s: trailing JSON value", path)
		}
		return fmt.Errorf("decode %s: trailing content: %w", path, err)
	}
	return nil
}

func doctor(root string, state CurrentState, graph WorkGraph) (DoctorOutput, error) {
	out := DoctorOutput{Root: root, Healthy: false}
	if err := validateState(root, state, graph); err != nil {
		return out, err
	}
	out.Checks = append(out.Checks,
		"canonical bootstrap files present",
		"state JSON decodes strictly",
		"work graph JSON decodes strictly",
		"work graph node ids and dependencies are valid",
		"work graph is acyclic",
		"active/preferred work units exist",
		"canonical product-spec path exists",
	)

	next, err := selectNext(graph)
	if err != nil {
		return out, err
	}
	if next == nil {
		out.Warnings = append(out.Warnings, "no executable READY work unit")
	} else if state.Execution.PreferredNextUnit != "" && state.Execution.PreferredNextUnit != next.ID {
		out.Warnings = append(out.Warnings,
			fmt.Sprintf("preferred next unit %s differs from computed next %s", state.Execution.PreferredNextUnit, next.ID))
	}
	out.Healthy = true
	return out, nil
}

func validateState(root string, state CurrentState, graph WorkGraph) error {
	if state.SchemaVersion != 1 || graph.SchemaVersion != 1 {
		return errors.New("unsupported state or graph schema version")
	}
	if state.StateKind == "" || graph.GraphKind == "" {
		return errors.New("state_kind and graph_kind are required")
	}
	if state.ProductSpec.CanonicalPath == "" {
		return errors.New("canonical product-spec path is empty")
	}
	productPath, err := resolveRepoPath(root, state.ProductSpec.CanonicalPath, false)
	if err != nil {
		return fmt.Errorf("canonical product-spec path %q is invalid: %w", state.ProductSpec.CanonicalPath, err)
	}
	if info, err := os.Stat(productPath); err != nil || !info.IsDir() { // #nosec G703 -- productPath is repository-confined and symlink-free.
		return fmt.Errorf("canonical product-spec path %q does not exist", state.ProductSpec.CanonicalPath)
	}

	index := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		if strings.TrimSpace(node.ID) == "" {
			return errors.New("work graph contains empty node id")
		}
		if _, exists := index[node.ID]; exists {
			return fmt.Errorf("duplicate work node id %s", node.ID)
		}
		if _, known := lifecycleRank[node.Status]; !known {
			return fmt.Errorf("work node %s has unknown status %s", node.ID, node.Status)
		}
		index[node.ID] = node
	}

	for _, node := range graph.Nodes {
		for _, dep := range node.DependsOn {
			if _, ok := index[dep]; !ok {
				return fmt.Errorf("work node %s depends on missing node %s", node.ID, dep)
			}
		}
		for _, unlocked := range node.Unlocks {
			if _, ok := index[unlocked]; !ok {
				return fmt.Errorf("work node %s unlocks missing node %s", node.ID, unlocked)
			}
		}
	}
	if err := validateAcyclic(graph.Nodes, index); err != nil {
		return err
	}
	if state.Execution.ActiveWorkUnit != "" {
		if _, ok := index[state.Execution.ActiveWorkUnit]; !ok {
			return fmt.Errorf("active work unit %s is missing from graph", state.Execution.ActiveWorkUnit)
		}
	}
	if state.Execution.PreferredNextUnit != "" {
		if _, ok := index[state.Execution.PreferredNextUnit]; !ok {
			return fmt.Errorf("preferred next unit %s is missing from graph", state.Execution.PreferredNextUnit)
		}
	}
	return nil
}

func validateAcyclic(nodes []WorkNode, index map[string]WorkNode) error {
	const (
		gray  = 1
		black = 2
	)
	state := make(map[string]int, len(nodes))
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case gray:
			return fmt.Errorf("work graph cycle detected at %s", id)
		case black:
			return nil
		}
		state[id] = gray
		for _, dep := range index[id].DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = black
		return nil
	}
	for _, node := range nodes {
		if err := visit(node.ID); err != nil {
			return err
		}
	}
	return nil
}

func selectNext(graph WorkGraph) (*WorkNode, error) {
	index := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		index[node.ID] = node
	}
	candidates := make([]WorkNode, 0)
	for _, node := range graph.Nodes {
		if node.Status != "READY" {
			continue
		}
		ok, err := dependenciesSatisfied(node, index)
		if err != nil {
			return nil, err
		}
		if ok {
			candidates = append(candidates, node)
		}
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })
	return &candidates[0], nil
}

func dependenciesSatisfied(node WorkNode, index map[string]WorkNode) (bool, error) {
	required := node.DependencyTerminalState
	if required == "" {
		required = "VERIFIED"
	}
	requiredRank, ok := lifecycleRank[required]
	if !ok || requiredRank < 0 {
		return false, fmt.Errorf("node %s has invalid dependency terminal state %s", node.ID, required)
	}
	for _, depID := range node.DependsOn {
		dep, ok := index[depID]
		if !ok {
			return false, fmt.Errorf("node %s depends on missing node %s", node.ID, depID)
		}
		rank, known := lifecycleRank[dep.Status]
		if !known || rank < requiredRank {
			return false, nil
		}
	}
	return true, nil
}

func printValue(v any, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(v); err != nil {
			fail(err)
		}
		return
	}

	switch x := v.(type) {
	case DoctorOutput:
		fmt.Printf("CMDR engineering doctor: healthy=%t\n", x.Healthy)
		fmt.Printf("root: %s\n", x.Root)
		for _, check := range x.Checks {
			fmt.Printf("PASS  %s\n", check)
		}
		for _, warning := range x.Warnings {
			fmt.Printf("WARN  %s\n", warning)
		}
	case StatusOutput:
		fmt.Printf("phase: %s — %s\n", x.Phase, x.PhaseTitle)
		fmt.Printf("foundation: %s\n", x.FoundationStatus)
		fmt.Printf("active: %s\n", x.ActiveWorkUnit)
		fmt.Printf("preferred next: %s\n", x.PreferredNext)
		fmt.Printf("computed next: %s\n", x.ComputedNext)
		fmt.Printf("product baseline: %s\n", x.ProductBaseline)
	case *WorkNode:
		fmt.Printf("%s — %s [%s]\n", x.ID, x.Title, x.Status)
	case SpecIndexSummary:
		fmt.Printf("spec files: %d\n", x.Files)
		fmt.Printf("active canonical documents: %d\n", x.ActiveCanonicalDocuments)
		fmt.Printf("tree digest: %s\n", x.TreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case SpecBaselineSummary:
		fmt.Printf("spec files: %d\n", x.Files)
		fmt.Printf("active canonical documents: %d\n", x.ActiveCanonicalDocuments)
		fmt.Printf("tree digest: %s\n", x.TreeDigest)
		fmt.Printf("baseline commit: %s\n", x.ProductSpecBaselineCommit)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case CoverageGraphSummary:
		fmt.Printf("entities: %d\n", x.Entities)
		fmt.Printf("edges: %d\n", x.Edges)
		fmt.Printf("owned capabilities: %d\n", x.OwnedCapabilities)
		fmt.Printf("reference-only capabilities: %d\n", x.ReferenceOnlyCapabilities)
		fmt.Printf("distinct requirement references: %d\n", x.DistinctRequirementReferences)
		fmt.Printf("distinct open-decision references: %d\n", x.DistinctOpenDecisionReferences)
		fmt.Printf("distinct permission references: %d\n", x.DistinctPermissionReferences)
		fmt.Printf("owned screens: %d\n", x.OwnedScreens)
		fmt.Printf("registered active screens: %d\n", x.RegisteredActiveScreens)
		fmt.Printf("unregistered owned screens: %v\n", x.UnregisteredOwnedScreens)
		fmt.Printf("reference-only screens: %d\n", x.ReferenceOnlyScreens)
		fmt.Printf("tree digest: %s\n", x.SpecTreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case ObligationSummary:
		fmt.Printf("obligations: %d\n", x.Obligations)
		fmt.Printf("by family: %v\n", x.ByFamily)
		fmt.Printf("unresolved references: %d\n", x.UnresolvedReferences)
		fmt.Printf("registered capabilities: %d\n", x.RegisteredCapabilities)
		fmt.Printf("registered requirements: %d\n", x.RegisteredRequirements)
		fmt.Printf("registered screens: %d\n", x.RegisteredScreens)
		fmt.Printf("registered permissions: %d\n", x.RegisteredPermissions)
		fmt.Printf("tree digest: %s\n", x.SpecTreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case CoverageAuditSummary:
		fmt.Printf("coverage gaps: %d\n", x.Gaps)
		fmt.Printf("by code: %v\n", x.ByCode)
		fmt.Printf("registered capabilities: %d\n", x.RegisteredCapabilities)
		fmt.Printf("registered requirements: %d\n", x.RegisteredRequirements)
		fmt.Printf("registered screens: %d\n", x.RegisteredScreens)
		fmt.Printf("registered permissions: %d\n", x.RegisteredPermissions)
		fmt.Printf("active open decisions: %d\n", x.ActiveOpenDecisions)
		fmt.Printf("tree digest: %s\n", x.SpecTreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case ManifestValidationSummary:
		fmt.Printf("manifests: %d\n", x.Manifests)
		fmt.Printf("legacy v1: %d\n", x.LegacyV1)
		fmt.Printf("strict v2: %d\n", x.StrictV2)
		fmt.Printf("mode: %s\n", x.Mode)
	case ComplexityAuditSummary:
		fmt.Printf("evaluated: %d\n", x.Evaluated)
		fmt.Printf("within budget: %d\n", x.WithinBudget)
		fmt.Printf("with warnings: %d\n", x.WithWarnings)
		fmt.Printf("split required: %d\n", x.SplitRequired)
		fmt.Printf("readiness violations: %d\n", x.ReadinessViolations)
	case ContextSummary:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("sources: %d\n", x.Sources)
		fmt.Printf("product sources: %d\n", x.ProductSources)
		fmt.Printf("engineering sources: %d\n", x.EngineeringSources)
		fmt.Printf("execution sources: %d\n", x.ExecutionSources)
		fmt.Printf("dependency units: %d\n", x.DependencyUnits)
		fmt.Printf("bundle digest: %s\n", x.BundleDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case ArchitectureAuditSummary:
		fmt.Printf("boundaries: %d\n", x.Boundaries)
		fmt.Printf("strict v2 manifests: %d\n", x.StrictV2Manifests)
		fmt.Printf("allowed path claims: %d\n", x.AllowedPathClaims)
		fmt.Printf("product spec read only: %t\n", x.ProductSpecReadOnly)
	case DependencyAuditSummary:
		fmt.Printf("runtime boundaries: %d\n", x.RuntimeBoundaries)
		fmt.Printf("supported manifests: %d\n", x.SupportedManifests)
		fmt.Printf("runtime dependencies: %d\n", x.RuntimeDependencies)
		fmt.Printf("approved: %d\n", x.Approved)
		fmt.Printf("unapproved: %d\n", x.Unapproved)
		fmt.Printf("unsupported manifests: %d\n", x.UnsupportedManifests)
	case BoundaryEdgeAuditSummary:
		fmt.Printf("runtime boundaries: %d\n", x.RuntimeBoundaries)
		fmt.Printf("local packages: %d\n", x.LocalPackages)
		fmt.Printf("observed edges: %d\n", x.ObservedEdges)
		fmt.Printf("cross-boundary edges: %d\n", x.CrossBoundaryEdges)
	case CheckCatalogSummary:
		fmt.Printf("checks: %d\n", x.Checks)
		fmt.Printf("mandatory: %d\n", x.Mandatory)
		fmt.Printf("always on PR: %d\n", x.AlwaysOnPR)
		fmt.Printf("by cost tier: %v\n", x.ByCostTier)
		fmt.Printf("by risk domain: %v\n", x.ByRiskDomain)
		fmt.Printf("prerequisite edges: %d\n", x.Prerequisites)
	case ImpactReport:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("changed paths: %v\n", x.ChangedPaths)
		fmt.Printf("risk domains: %v\n", x.RiskDomains)
		fmt.Printf("high risk: %t\n", x.HighRisk)
		fmt.Printf("high risk reasons: %v\n", x.HighRiskReasons)
		fmt.Printf("unknown paths: %v\n", x.UnknownPaths)
	case ValidationPlan:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("tier: %s\n", x.Tier)
		fmt.Printf("high risk: %t\n", x.HighRisk)
		fmt.Printf("risk domains: %v\n", x.RiskDomains)
		fmt.Printf("selected checks: %d\n", len(x.SelectedChecks))
		fmt.Printf("by cost tier: %v\n", x.ByCostTier)
		fmt.Printf("cost units: %d\n", x.CostUnits)
	case SecurityGateAuditSummary:
		fmt.Printf("security gates: %d\n", x.Gates)
		fmt.Printf("by readiness: %v\n", x.ByReadiness)
		fmt.Printf("by family: %v\n", x.ByFamily)
		fmt.Printf("by stage: %v\n", x.ByStage)
		fmt.Printf("blocking: %d\n", x.Blocking)
		fmt.Printf("conditional blocking: %d\n", x.Conditional)
	case SecurityTestAuditSummary:
		fmt.Printf("runtime boundaries: %d\n", x.RuntimeBoundaries)
		fmt.Printf("registered scopes: %d\n", x.RegisteredScopes)
		fmt.Printf("authorization required: %d\n", x.AuthorizationRequired)
		fmt.Printf("tenant isolation required: %d\n", x.TenantIsolationRequired)
		fmt.Printf("coverage status: %s\n", x.CoverageStatus)
		fmt.Printf("global coverage floor: %.2f\n", x.GlobalCoverageFloorPercent)
		fmt.Printf("changed security-critical coverage floor: %.2f\n", x.ChangedCoverageFloorPercent)
	case DeepSecurityAuditSummary:
		fmt.Printf("stage: %s\n", x.Stage)
		fmt.Printf("runtime boundaries: %d\n", x.RuntimeBoundaries)
		fmt.Printf("security-sensitive change: %t\n", x.SecuritySensitiveChange)
		fmt.Printf("registered targets: %d\n", x.RegisteredTargets)
		fmt.Printf("selected targets: %d\n", x.SelectedTargets)
		fmt.Printf("deferred gates: %d\n", x.DeferredGates)
		fmt.Printf("validated evidence: %d\n", x.ValidatedEvidence)
	case DecisionRegistryAuditSummary:
		fmt.Printf("decisions: %d\n", x.Decisions)
		fmt.Printf("by class: %v\n", x.ByClass)
		fmt.Printf("by status: %v\n", x.ByStatus)
		fmt.Printf("accepted: %d\n", x.Accepted)
		fmt.Printf("blocked product: %d\n", x.BlockedProduct)
		fmt.Printf("performance-sensitive: %d\n", x.PerformanceSensitive)
		fmt.Printf("critical: %d\n", x.Critical)
	case PerformanceRegistryAuditSummary:
		fmt.Printf("performance environments: %d\n", x.Environments)
		fmt.Printf("performance targets: %d\n", x.Targets)
		fmt.Printf("performance metrics: %d\n", x.Metrics)
		fmt.Printf("runtime boundaries: %d\n", x.RuntimeBoundaries)
		fmt.Printf("runtime targeted boundaries: %d\n", x.RuntimeTargetedBoundaries)
		fmt.Printf("runtime coverage status: %s\n", x.RuntimeCoverageStatus)
		fmt.Printf("by scope: %v\n", x.ByScope)
		fmt.Printf("by stage: %v\n", x.ByStage)
	case PerformanceBenchmarkAuditSummary:
		fmt.Printf("stage: %s\n", x.Stage)
		fmt.Printf("environment: %s\n", x.EnvironmentID)
		fmt.Printf("source SHA: %s\n", x.SourceSHA)
		fmt.Printf("registered targets: %d\n", x.Registered)
		fmt.Printf("selected targets: %d\n", x.Selected)
		fmt.Printf("executed targets: %d\n", x.Executed)
		fmt.Printf("metrics evaluated: %d\n", x.Metrics)
		fmt.Printf("status: %s\n", x.Status)
	case PerformanceCacheAuditSummary:
		fmt.Printf("cache records: %d\n", x.Records)
		fmt.Printf("reusable records: %d\n", x.Reusable)
		fmt.Printf("invalid records: %d\n", x.Invalid)
		fmt.Printf("regressions: %d\n", x.Regressions)
		fmt.Printf("status: %s\n", x.Status)
		fmt.Printf("revisit signals: %d\n", len(x.Signals))
	case WorkLeaseAuditSummary:
		fmt.Printf("as of: %s\n", x.AsOf)
		fmt.Printf("claims: %d active=%d expired=%d released=%d\n", x.Claims, x.Active, x.Expired, x.Released)
	case WorkLeaseEvaluation:
		fmt.Printf("action: %s\n", x.Action)
		fmt.Printf("allowed: %t\n", x.Allowed)
		fmt.Printf("reasons: %v\n", x.Reasons)
	case RecoveryJournalAuditSummary:
		fmt.Printf("events: %d work_units=%d agents=%d\n", x.Events, x.WorkUnits, x.Agents)
		fmt.Printf("last sequence: %d\n", x.LastSeq)
		fmt.Printf("last digest: %s\n", x.LastDigest)
		fmt.Printf("status: %s\n", x.Status)
	case ResumeCheckpoint:
		fmt.Printf("status: %s\n", x.Status)
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("events: %d\n", x.Events)
		fmt.Printf("head: %s\n", x.HeadSHA)
	case RecoveryReconciliation:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("current head: %s\n", x.CurrentHead)
		fmt.Printf("checkpoint: %s (%s)\n", x.Checkpoint, x.CheckpointHead)
		fmt.Printf("git relation: %s\n", x.GitRelation)
		fmt.Printf("lease status: %s\n", x.LeaseStatus)
		fmt.Printf("ci status: %s\n", x.CIStatus)
		fmt.Printf("outcome: %s\n", x.Outcome)
		fmt.Printf("next action: %s\n", x.NextAction)
		fmt.Printf("reasons: %v\n", x.Reasons)
	case CoordinationAuditSummary:
		fmt.Printf("as of: %s\n", x.AsOf)
		fmt.Printf("active claims: %d mutating=%d read-only=%d\n", x.ActiveClaims, x.Mutating, x.ReadOnly)
		fmt.Printf("conflicts: %d\n", x.Conflicts)
		fmt.Printf("status: %s\n", x.Status)
	case EngineMetricsRegistrySummary:
		fmt.Printf("metrics: %d\n", x.Metrics)
		fmt.Printf("optimization eligible: %d\n", x.OptimizationEligible)
		fmt.Printf("non-optimizable: %d\n", x.NonOptimizable)
		fmt.Printf("by category: %v\n", x.ByCategory)
		fmt.Printf("by safety class: %v\n", x.BySafetyClass)
	case CoordinationHandoff:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("lease: %s (%s)\n", x.LeaseID, x.LeaseMode)
		fmt.Printf("from agent: %s\n", x.FromAgent)
		fmt.Printf("head: %s\n", x.HeadSHA)
		fmt.Printf("revalidation required: %t\n", x.RevalidationRequired)
		fmt.Printf("digest: %s\n", x.Digest)
	case SecretScanSummary:
		fmt.Printf("mode: %s\n", x.Mode)
		fmt.Printf("candidate paths: %d\n", x.CandidatePaths)
		fmt.Printf("scanned files: %d\n", x.ScannedFiles)
		fmt.Printf("allowlisted: %d\n", x.Allowlisted)
		fmt.Printf("findings: %d\n", x.FindingCount)
		fmt.Printf("by rule: %v\n", x.ByRule)
	case SASTSummary:
		fmt.Printf("tool: %s %s\n", x.Tool, x.Version)
		fmt.Printf("scan root: %s\n", x.ScanRoot)
		fmt.Printf("findings: %d\n", x.FindingCount)
		fmt.Printf("by severity: %v\n", x.BySeverity)
	case SCASummary:
		fmt.Printf("tool: %s %s\n", x.Tool, x.Version)
		fmt.Printf("database: %s\n", x.Database)
		fmt.Printf("go version: %s\n", x.GoVersion)
		fmt.Printf("modules: %d\n", x.Modules)
		fmt.Printf("informational findings: %d\n", x.InformationalFindings)
		fmt.Printf("actionable findings: %d\n", x.ActionableFindings)
	case SBOMSummary:
		fmt.Printf("schema: CycloneDX %s\n", x.SchemaVersion)
		fmt.Printf("generator: %s %s\n", x.Generator, x.GeneratorVersion)
		fmt.Printf("source SHA: %s\n", x.SourceSHA)
		fmt.Printf("components: %d (runtime=%d development=%d)\n", x.Components, x.RuntimeComponents, x.DevelopmentComponents)
		fmt.Printf("unsupported manifests: %d\n", x.UnsupportedManifests)
		fmt.Printf("digest: %s\n", x.DigestSHA256)
		fmt.Printf("output: %s\n", x.Output)
	case PilotContractAuditSummary:
		fmt.Printf("pilot contract: %s\n", x.ContractID)
		fmt.Printf("pilot slice: %s\n", x.SliceID)
		fmt.Printf("runtime boundary: %s\n", x.RuntimeBoundary)
		fmt.Printf("fixtures: %d\n", x.Fixtures)
		fmt.Printf("negative fixtures: %d\n", x.NegativeFixtures)
		fmt.Printf("runtime dependencies: %d\n", x.RuntimeDependencies)
		fmt.Printf("runtime files: %d\n", x.RuntimeFiles)
		fmt.Printf("status: %s\n", x.Status)
	case PilotE2EAuditSummary:
		fmt.Printf("pilot slice: %s\n", x.SliceID)
		fmt.Printf("capability: %s\n", x.Capability)
		fmt.Printf("contract: %s\n", x.ContractID)
		fmt.Printf("adversarial tests: %d\n", x.AdversarialTests)
		fmt.Printf("coverage: %.2f%%\n", x.RuntimeCoveragePercent)
		fmt.Printf("security: auth=%s tenant=%s sast=%d sca=%d\n", x.AuthorizationNegative, x.TenantIsolationNegative, x.SASTFindings, x.SCAActionableFindings)
		fmt.Printf("performance: p95=%.9fms budget=%.3fms\n", x.WorstObservedP95MS, x.PerformanceBudgetMS)
		fmt.Printf("recovery: %s / %s\n", x.RecoveryOutcome, x.RecoveryNextAction)
		fmt.Printf("production readiness claim: %t\n", x.ProductionReadinessClaim)
		fmt.Printf("limitations: %d\n", x.Limitations)
		fmt.Printf("status: %s\n", x.Status)
	case PilotScopeAuditSummary:
		fmt.Printf("pilot slice: %s\n", x.SliceID)
		fmt.Printf("capability: %s\n", x.Capability)
		fmt.Printf("requirements: %d\n", x.Requirements)
		fmt.Printf("permissions: %d\n", x.Permissions)
		fmt.Printf("open decisions: %d\n", x.OpenDecisions)
		fmt.Printf("contracts: %d\n", x.Contracts)
		fmt.Printf("canonical objects: %d\n", x.CanonicalObjects)
		fmt.Printf("product graph digest: %s\n", x.ProductGraphDigest)
		fmt.Printf("status: %s\n", x.Status)
	case GitChangesSummary:
		fmt.Printf("base commit: %s\n", x.BaseCommit)
		fmt.Printf("head commit: %s\n", x.HeadCommit)
		fmt.Printf("changed paths: %d\n", x.Count)
		fmt.Printf("output: %s\n", x.Output)
	case ValidationExecutionSummary:
		fmt.Printf("work unit: %s\n", x.WorkUnit)
		fmt.Printf("tier: %s\n", x.Tier)
		fmt.Printf("selected checks: %d\n", x.SelectedChecks)
		fmt.Printf("executed checks: %d\n", x.ExecutedChecks)
		fmt.Printf("preflight satisfied: %d\n", x.PreflightSatisfied)
	default:
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: cmdr-dev <doctor|status|next|spec-index|spec-baseline|coverage-graph|obligations|coverage-audit|validate-manifests|complexity-audit|context|architecture-audit|dependency-audit|boundary-edge-audit|check-catalog-audit|impact|validation-plan|security-gate-audit|security-test-audit|deep-security-audit|decision-registry-audit|research-packet-audit|decision-gate-audit|decision-freshness-audit|performance-registry-audit|performance-benchmark-audit|deep-performance-audit|performance-cache-audit|lease-audit|lease-evaluate|recovery-journal-audit|resume-checkpoint|recovery-reconcile|recovery-reconcile-audit|coordination-audit|coordination-handoff|metrics-registry-audit|decision-cache|decision-freshness-snapshot|research-context|secret-scan|sast-go|sca-go|sbom|pilot-scope-audit|pilot-contract-audit|pilot-e2e-audit|git-changes|validation-run> [--root PATH] [--json] [--check] [--output PATH]")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "cmdr-dev:", err)
	os.Exit(1)
}

func validateKnownObligationCounts(summary ObligationSummary, state CurrentState) error {
	if summary.RegisteredCapabilities != state.ProductSpec.KnownCapabilities {
		return fmt.Errorf("registered capability count mismatch: expected %d, got %d", state.ProductSpec.KnownCapabilities, summary.RegisteredCapabilities)
	}
	if summary.RegisteredRequirements != state.ProductSpec.KnownRequirements {
		return fmt.Errorf("registered requirement count mismatch: expected %d, got %d", state.ProductSpec.KnownRequirements, summary.RegisteredRequirements)
	}
	if summary.RegisteredScreens != state.ProductSpec.KnownActiveScreens {
		return fmt.Errorf("registered active screen count mismatch: expected %d, got %d", state.ProductSpec.KnownActiveScreens, summary.RegisteredScreens)
	}
	return nil
}

func validateKnownCoverageAuditCounts(summary CoverageAuditSummary, state CurrentState) error {
	if summary.RegisteredCapabilities != state.ProductSpec.KnownCapabilities {
		return fmt.Errorf("coverage audit capability count mismatch: expected %d, got %d", state.ProductSpec.KnownCapabilities, summary.RegisteredCapabilities)
	}
	if summary.RegisteredRequirements != state.ProductSpec.KnownRequirements {
		return fmt.Errorf("coverage audit requirement count mismatch: expected %d, got %d", state.ProductSpec.KnownRequirements, summary.RegisteredRequirements)
	}
	if summary.RegisteredScreens != state.ProductSpec.KnownActiveScreens {
		return fmt.Errorf("coverage audit active screen count mismatch: expected %d, got %d", state.ProductSpec.KnownActiveScreens, summary.RegisteredScreens)
	}
	if summary.ActiveOpenDecisions != state.ProductSpec.KnownOpenDecisions {
		return fmt.Errorf("coverage audit open-decision count mismatch: expected %d, got %d", state.ProductSpec.KnownOpenDecisions, summary.ActiveOpenDecisions)
	}
	return nil
}
