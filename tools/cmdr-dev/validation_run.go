package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ExecutedCheck struct {
	ID          string `json:"id"`
	ExecutorKey string `json:"executor_key"`
	Status      string `json:"status"`
	Evidence    string `json:"evidence"`
}

type ValidationExecutionSummary struct {
	WorkUnit           string          `json:"work_unit"`
	Tier               string          `json:"tier"`
	SelectedChecks     int             `json:"selected_checks"`
	ExecutedChecks     int             `json:"executed_checks"`
	PreflightSatisfied int             `json:"preflight_satisfied"`
	Checks             []ExecutedCheck `json:"checks"`
}

func runValidationExecution(root, workUnit, changesFile string, state CurrentState, graph WorkGraph) (ValidationExecutionSummary, error) {
	plan, err := runValidationPlan(root, workUnit, changesFile, state, graph)
	if err != nil {
		return ValidationExecutionSummary{}, err
	}
	catalog, err := loadCheckCatalog(root)
	if err != nil {
		return ValidationExecutionSummary{}, err
	}
	if err := validateValidationPlan(plan, catalog); err != nil {
		return ValidationExecutionSummary{}, err
	}

	tempDir, err := os.MkdirTemp(root, ".cmdr-validation-")
	if err != nil {
		return ValidationExecutionSummary{}, err
	}
	defer os.RemoveAll(tempDir)

	summary := ValidationExecutionSummary{
		WorkUnit:       plan.WorkUnit,
		Tier:           plan.Tier,
		SelectedChecks: len(plan.SelectedChecks),
	}
	for _, planned := range plan.SelectedChecks {
		mode, err := validationExecutorMode(planned.ExecutorKey)
		if err != nil {
			return summary, fmt.Errorf("%s: %w", planned.ID, err)
		}
		if mode == "preflight" {
			summary.PreflightSatisfied++
			summary.Checks = append(summary.Checks, ExecutedCheck{
				ID: planned.ID, ExecutorKey: planned.ExecutorKey,
				Status:   "PREFLIGHT_SATISFIED",
				Evidence: preflightEvidence(planned.ExecutorKey),
			})
			continue
		}
		evidence, err := executeValidationCheck(root, tempDir, changesFile, planned.ExecutorKey, state, graph)
		if err != nil {
			summary.Checks = append(summary.Checks, ExecutedCheck{
				ID: planned.ID, ExecutorKey: planned.ExecutorKey,
				Status: "FAIL", Evidence: err.Error(),
			})
			return summary, fmt.Errorf("selected check %s (%s) failed: %w", planned.ID, planned.ExecutorKey, err)
		}
		summary.ExecutedChecks++
		summary.Checks = append(summary.Checks, ExecutedCheck{
			ID: planned.ID, ExecutorKey: planned.ExecutorKey,
			Status: "PASS", Evidence: evidence,
		})
	}
	if summary.ExecutedChecks+summary.PreflightSatisfied != summary.SelectedChecks {
		return summary, fmt.Errorf("execution accounting mismatch: selected=%d executed=%d preflight=%d",
			summary.SelectedChecks, summary.ExecutedChecks, summary.PreflightSatisfied)
	}
	return summary, nil
}

func validationExecutorMode(key string) (string, error) {
	switch key {
	case "git-changes", "impact", "validation-plan", "validation-run":
		return "preflight", nil
	case "gofmt", "go-vet", "go-unit",
		"spec-index", "spec-baseline", "coverage-graph", "obligations", "coverage-audit",
		"validate-manifests", "architecture-audit", "dependency-audit", "boundary-edge-audit",
		"complexity-audit", "context", "doctor", "next", "check-catalog-audit", "security-gate-audit", "security-test-audit", "deep-security-audit", "decision-registry-audit", "secret-scan", "gosec-go", "govulncheck-go", "sbom":
		return "execute", nil
	default:
		return "", fmt.Errorf("unsupported executor_key %q", key)
	}
}

func preflightEvidence(key string) string {
	switch key {
	case "git-changes":
		return "change set already derived from exact validated Git commit identities"
	case "impact":
		return "impact report already computed while building the selected plan"
	case "validation-plan":
		return "selected plan already computed and structurally validated"
	case "validation-run":
		return "current process is the adaptive validation executor"
	default:
		return "preflight satisfied"
	}
}

func executeValidationCheck(root, tempDir, changesFile, key string, state CurrentState, graph WorkGraph) (string, error) {
	switch key {
	case "gofmt":
		cmd := exec.Command("gofmt", "-d", ".")
		cmd.Dir = filepath.Join(root, "tools", "cmdr-dev")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 && len(bytes.TrimSpace(stdout.Bytes())) > 0 {
				return "", fmt.Errorf("gofmt diff is non-empty:\n%s", stdout.String())
			}
			diagnostic := strings.TrimSpace(stderr.String())
			if diagnostic == "" {
				diagnostic = err.Error()
			}
			return "", fmt.Errorf("gofmt execution failed: %s", diagnostic)
		}
		if len(bytes.TrimSpace(stdout.Bytes())) > 0 {
			return "", fmt.Errorf("gofmt diff is non-empty:\n%s", stdout.String())
		}
		return "gofmt produced no diff", nil
	case "go-vet":
		return runFixedGoProcess(filepath.Join(root, "tools", "cmdr-dev"), "vet")
	case "go-unit":
		return runFixedGoProcess(filepath.Join(root, "tools", "cmdr-dev"), "test")
	case "spec-index":
		out := filepath.Join(tempDir, "spec-inventory.json")
		summary, err := runSpecIndex(root, state.ProductSpec.CanonicalPath, out, false)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("files=%d canonical=%d digest=%s", summary.Files, summary.ActiveCanonicalDocuments, summary.TreeDigest), nil
	case "spec-baseline":
		summary, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("baseline=%s digest=%s", summary.ProductSpecBaselineCommit, summary.TreeDigest), nil
	case "coverage-graph":
		out := filepath.Join(tempDir, "product-graph.json")
		summary, err := runCoverageGraph(root, state.ProductSpec.CanonicalPath, out, false)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("entities=%d edges=%d digest=%s", summary.Entities, summary.Edges, summary.SpecTreeDigest), nil
	case "obligations":
		out := filepath.Join(tempDir, "obligations.json")
		summary, err := runObligations(root, state.ProductSpec.CanonicalPath, out, false)
		if err != nil {
			return "", err
		}
		if err := validateKnownObligationCounts(summary, state); err != nil {
			return "", err
		}
		return fmt.Sprintf("obligations=%d unresolved=%d", summary.Obligations, summary.UnresolvedReferences), nil
	case "coverage-audit":
		out := filepath.Join(tempDir, "coverage-audit.json")
		summary, err := runCoverageAudit(root, state.ProductSpec.CanonicalPath, out, false)
		if err != nil {
			return "", err
		}
		if err := validateKnownCoverageAuditCounts(summary, state); err != nil {
			return "", err
		}
		return fmt.Sprintf("gaps=%d digest=%s", summary.Gaps, summary.SpecTreeDigest), nil
	case "validate-manifests":
		summary, err := runManifestValidation(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("manifests=%d strict_v2=%d", summary.Manifests, summary.StrictV2), nil
	case "architecture-audit":
		summary, err := runArchitectureAudit(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("boundaries=%d path_claims=%d", summary.Boundaries, summary.AllowedPathClaims), nil
	case "dependency-audit":
		summary, err := runDependencyAudit(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("runtime_dependencies=%d unapproved=%d", summary.RuntimeDependencies, summary.Unapproved), nil
	case "boundary-edge-audit":
		summary, err := runBoundaryEdgeAudit(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("observed_edges=%d cross_boundary=%d", summary.ObservedEdges, summary.CrossBoundaryEdges), nil
	case "complexity-audit":
		summary, err := runComplexityAudit(root, graph)
		if err != nil {
			return "", err
		}
		if err := validateComplexityReadiness(summary); err != nil {
			return "", err
		}
		return fmt.Sprintf("evaluated=%d split_required=%d readiness_violations=%d", summary.Evaluated, summary.SplitRequired, summary.ReadinessViolations), nil
	case "context":
		out := filepath.Join(tempDir, "context.json")
		summary, err := runContextCompiler(root, state.Execution.ActiveWorkUnit, out, false, state, graph)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("sources=%d dependencies=%d digest=%s", summary.Sources, summary.DependencyUnits, summary.BundleDigest), nil
	case "doctor":
		summary, err := doctor(root, state, graph)
		if err != nil {
			return "", err
		}
		if !summary.Healthy {
			return "", fmt.Errorf("repository doctor is unhealthy")
		}
		return fmt.Sprintf("healthy=true checks=%d", len(summary.Checks)), nil
	case "next":
		next, err := selectNext(graph)
		if err != nil {
			return "", err
		}
		if next == nil {
			return "", fmt.Errorf("no executable READY work unit")
		}
		return "next=" + next.ID, nil
	case "check-catalog-audit":
		summary, err := runCheckCatalogAudit(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("checks=%d mandatory=%d", summary.Checks, summary.Mandatory), nil
	case "security-gate-audit":
		summary, err := runSecurityGateAudit(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("security_gates=%d active=%d specified=%d deferred_runtime=%d",
			summary.Gates,
			summary.ByReadiness["active"],
			summary.ByReadiness["specified"],
			summary.ByReadiness["deferred-runtime"]), nil
	case "security-test-audit":
		summary, err := runSecurityTestAudit(root, changesFile)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("runtime_boundaries=%d scopes=%d coverage_status=%s global_floor=%.0f changed_floor=%.0f auth_required=%d tenant_required=%d",
			summary.RuntimeBoundaries, summary.RegisteredScopes, summary.CoverageStatus,
			summary.GlobalCoverageFloorPercent, summary.ChangedCoverageFloorPercent,
			summary.AuthorizationRequired, summary.TenantIsolationRequired), nil
	case "deep-security-audit":
		summary, err := runDeepSecurityAudit(root, "pr", changesFile, "")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("stage=%s sensitive=%t runtime_boundaries=%d targets=%d selected=%d deferred=%d evidence=%d",
			summary.Stage, summary.SecuritySensitiveChange, summary.RuntimeBoundaries,
			summary.RegisteredTargets, summary.SelectedTargets, summary.DeferredGates, summary.ValidatedEvidence), nil
	case "decision-registry-audit":
		summary, err := runDecisionRegistryAudit(root, graph)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("decisions=%d accepted=%d blocked_product=%d critical=%d performance_sensitive=%d",
			summary.Decisions, summary.Accepted, summary.BlockedProduct, summary.Critical, summary.PerformanceSensitive), nil
	case "secret-scan":
		summary, err := runSecretScan(root, changesFile, false)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("mode=%s candidates=%d scanned=%d findings=%d allowlisted=%d skipped_binary=%d skipped_generated=%d",
			summary.Mode, summary.CandidatePaths, summary.ScannedFiles, summary.FindingCount,
			summary.Allowlisted, summary.SkippedBinary, summary.SkippedGenerated), nil
	case "gosec-go":
		summary, err := runGoSAST(root)
		if err != nil {
			return "", fmt.Errorf("%w; safe_findings=%s", err, formatSASTSafeFindings(summary))
		}
		return fmt.Sprintf("tool=%s version=%s scan_root=%s findings=%d",
			summary.Tool, summary.Version, summary.ScanRoot, summary.FindingCount), nil
	case "govulncheck-go":
		summary, err := runGoSCA(root)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("tool=%s version=%s modules=%d informational=%d actionable=%d db=%s",
			summary.Tool, summary.Version, summary.Modules, summary.InformationalFindings,
			summary.ActionableFindings, summary.Database), nil
	case "sbom":
		out := filepath.Join(tempDir, "cmdr.cdx.json")
		summary, err := runSBOM(root, out)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("schema=cyclonedx-%s generator=%s/%s source=%s components=%d runtime=%d development=%d unsupported=%d digest=%s",
			summary.SchemaVersion, summary.Generator, summary.GeneratorVersion, summary.SourceSHA,
			summary.Components, summary.RuntimeComponents, summary.DevelopmentComponents,
			summary.UnsupportedManifests, summary.DigestSHA256), nil
	default:
		return "", fmt.Errorf("unsupported executable key %q", key)
	}
}

func runFixedGoProcess(dir, action string) (string, error) {
	var cmd *exec.Cmd
	var label string
	switch action {
	case "vet":
		cmd = exec.Command("go", "vet", "./...")
		label = "go vet ./..."
	case "test":
		cmd = exec.Command("go", "test", "./...")
		label = "go test ./..."
	default:
		return "", fmt.Errorf("unsupported fixed Go action %q", action)
	}
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = strings.TrimSpace(stdout.String())
		}
		return "", fmt.Errorf("%s failed: %s", label, message)
	}
	message := strings.TrimSpace(stdout.String())
	if message == "" {
		message = "process exited successfully"
	}
	return message, nil
}
