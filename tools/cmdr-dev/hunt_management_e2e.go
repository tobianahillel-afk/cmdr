package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const (
	huntManagementRuntimeProgressPath = "work/lots/E10-INV-005B-RUNTIME/PROGRESS.json"
	huntManagementHandoffPath         = "engineering/implementation/hunt-management/bounded-handoff.json"
	huntManagementRuntimeDirectory    = "product-runtime/hunt-management"
)

var huntManagementAdversarialTests = []string{
	"TestProjectMatchesPredeclaredContractFixtures",
	"TestProjectionDeepCopiesCallerOwnedWorkspace",
	"TestProjectionPreservesCanonicalOwnershipVersionAndProvenance",
	"TestProjectRejectsAdditionalInvalidInputs",
	"TestProjectionDoesNotSynthesizeOptionalCorrelationOrProvenance",
}

var huntManagementRequiredLimitations = []string{
	"Hunt lifecycle",
	"Case promotion",
	"Hypothesis",
	"Query",
	"Search Job execution",
	"storage",
	"provider",
	"retention",
	"collaboration",
	"final query dialect",
	"final Hunt UI",
	"production",
	"full-capability",
}

type HuntManagementRuntimeProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	Validation         struct {
		PushWorkflowRun        int64  `json:"push_workflow_run"`
		PullRequestWorkflowRun int64  `json:"pull_request_workflow_run"`
		Result                 string `json:"result"`
		HuntManagementContract string `json:"hunt_management_contract"`
		SecurityTests          string `json:"security_tests"`
		SASTGo                 string `json:"sast_go"`
		SCAGo                  string `json:"sca_go"`
		RuntimeDependencyAudit string `json:"runtime_dependency_audit"`
		BoundaryEdgeAudit      string `json:"boundary_edge_audit"`
	} `json:"validation"`
	RuntimeObservation struct {
		Contract              string  `json:"contract"`
		RuntimeState          string  `json:"runtime_state"`
		Fixtures              int     `json:"fixtures"`
		Positive              int     `json:"positive"`
		Negative              int     `json:"negative"`
		RuntimeDependencies   int     `json:"runtime_dependencies"`
		RuntimeBoundaries     int     `json:"runtime_boundaries"`
		SecurityScopes        int     `json:"security_scopes"`
		GlobalRuntimeCoverage float64 `json:"global_runtime_coverage_percent"`
		GlobalCoverageFloor   float64 `json:"global_coverage_floor_percent"`
		ChangedSecurityFloor  float64 `json:"changed_security_floor_percent"`
		SASTFindings          int     `json:"sast_findings"`
		SCAActionableFindings int     `json:"sca_actionable_findings"`
	} `json:"runtime_observation"`
	Invariants         []string `json:"invariants"`
	ProductSpecMutated bool     `json:"product_spec_mutated"`
	NextUnlocked       []string `json:"next_unlocked"`
}

type HuntManagementClosureHandoff struct {
	SchemaVersion                 int      `json:"schema_version"`
	HandoffKind                   string   `json:"handoff_kind"`
	Readiness                     string   `json:"readiness"`
	Capability                    string   `json:"capability"`
	ProductionReadinessClaim      bool     `json:"production_readiness_claim"`
	FullCapabilityCompletionClaim bool     `json:"full_capability_completion_claim"`
	Class2MutationEnabled         bool     `json:"class2_mutation_enabled"`
	ProviderSelected              bool     `json:"provider_selected"`
	StorageEngineSelected         bool     `json:"storage_engine_selected"`
	RetentionPolicySelected       bool     `json:"retention_policy_selected"`
	CollaborationBackendSelected  bool     `json:"collaboration_backend_selected"`
	FinalQueryDialectSelected     bool     `json:"final_query_dialect_selected"`
	DedicatedFinalUIClaim         bool     `json:"dedicated_final_ui_claim"`
	BlockingOpenDecisions         []string `json:"blocking_open_decisions"`
	EvidenceSources               []string `json:"evidence_sources"`
	Limitations                   []string `json:"limitations"`
	NextConditions                []string `json:"next_conditions"`
}

type HuntManagementAdversarialSummary struct {
	Expected []string `json:"expected"`
	Passed   []string `json:"passed"`
	Count    int      `json:"count"`
}

type HuntManagementE2EAuditSummary struct {
	Capability                        string  `json:"capability"`
	ContractID                        string  `json:"contract_id"`
	AdversarialTests                  int     `json:"adversarial_tests"`
	RuntimeCoveragePercent            float64 `json:"runtime_coverage_percent"`
	SASTFindings                      int     `json:"sast_findings"`
	SCAActionableFindings             int     `json:"sca_actionable_findings"`
	ProjectionP95MS                   float64 `json:"projection_p95_ms"`
	ProjectionBudgetMS                float64 `json:"projection_budget_ms"`
	RuntimeUnchangedSinceVerification bool    `json:"runtime_unchanged_since_verification"`
	BlockingOpenDecisions             int     `json:"blocking_open_decisions"`
	Limitations                       int     `json:"limitations"`
	ProductionReadinessClaim          bool    `json:"production_readiness_claim"`
	FullCapabilityClaim               bool    `json:"full_capability_claim"`
	Status                            string  `json:"status"`
}

func runHuntManagementE2EAudit(root, _ string, state CurrentState) (HuntManagementE2EAuditSummary, error) {
	if _, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true); err != nil {
		return HuntManagementE2EAuditSummary{}, fmt.Errorf("Hunt Management Product Spec baseline: %w", err)
	}
	contract, err := runHuntManagementContractAudit(root)
	if err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	var runtime HuntManagementRuntimeProgress
	if err := decodeStrict(root, huntManagementRuntimeProgressPath, &runtime); err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	var handoff HuntManagementClosureHandoff
	if err := decodeStrict(root, huntManagementHandoffPath, &handoff); err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	if err := validateHuntManagementRuntimeProgress(runtime); err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	if err := validateHuntManagementHandoff(handoff); err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	currentHead, err := currentSourceCommit(root)
	if err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	if err := validateHuntManagementRuntimeUnchanged(root, runtime.FinalValidatedHead, currentHead); err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	adversarial, err := runHuntManagementAdversarialRuntimeTests(root)
	if err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	performance, err := runPerformanceBenchmarkTargets(root, "pr", defaultCIEnvironment, []string{
		"PERF-TGT-HUNT-MANAGEMENT-PROJECTION",
	})
	if err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	observed, budget, err := huntManagementPerformanceMetric(
		performance,
		"PERF-TGT-HUNT-MANAGEMENT-PROJECTION",
		"hunt-workspace-projection-p95-latency",
	)
	if err != nil {
		return HuntManagementE2EAuditSummary{}, err
	}
	if contract.Status != "PASS" || contract.RuntimeState != "implemented" ||
		contract.RuntimeDependencies != 0 || contract.Fixtures != 15 ||
		contract.PositiveFixtures != 3 || contract.NegativeFixtures != 12 {
		return HuntManagementE2EAuditSummary{}, fmt.Errorf("Hunt Management contract evidence is incomplete")
	}
	if adversarial.Count != len(huntManagementAdversarialTests) ||
		!sameStringSet(adversarial.Passed, huntManagementAdversarialTests) {
		return HuntManagementE2EAuditSummary{}, fmt.Errorf("Hunt Management adversarial evidence is incomplete")
	}
	return HuntManagementE2EAuditSummary{
		Capability:                        "CAP-INV-005",
		ContractID:                        contract.ContractID,
		AdversarialTests:                  adversarial.Count,
		RuntimeCoveragePercent:            runtime.RuntimeObservation.GlobalRuntimeCoverage,
		SASTFindings:                      runtime.RuntimeObservation.SASTFindings,
		SCAActionableFindings:             runtime.RuntimeObservation.SCAActionableFindings,
		ProjectionP95MS:                   observed,
		ProjectionBudgetMS:                budget,
		RuntimeUnchangedSinceVerification: true,
		BlockingOpenDecisions:             len(handoff.BlockingOpenDecisions),
		Limitations:                       len(handoff.Limitations),
		ProductionReadinessClaim:          handoff.ProductionReadinessClaim,
		FullCapabilityClaim:               handoff.FullCapabilityCompletionClaim,
		Status:                            "PASS",
	}, nil
}

func validateHuntManagementRuntimeProgress(progress HuntManagementRuntimeProgress) error {
	if progress.SchemaVersion != 1 || progress.WorkUnit != "E10-INV-005B-RUNTIME" ||
		progress.Status != "VERIFIED" || progress.Validation.Result != "PASS" ||
		progress.ProductSpecMutated {
		return fmt.Errorf("Hunt Management runtime progress is not verified")
	}
	normalized, err := validateFullCommitID(progress.FinalValidatedHead)
	if err != nil || normalized != progress.FinalValidatedHead {
		return fmt.Errorf("Hunt Management runtime progress has invalid final validated head")
	}
	if progress.RuntimeObservation.Contract != "HUNT-WORKSPACE-READONLY-CONTRACT-V1" ||
		progress.RuntimeObservation.RuntimeState != "implemented" ||
		progress.RuntimeObservation.Fixtures != 15 || progress.RuntimeObservation.Positive != 3 ||
		progress.RuntimeObservation.Negative != 12 || progress.RuntimeObservation.RuntimeDependencies != 0 ||
		progress.RuntimeObservation.GlobalRuntimeCoverage < progress.RuntimeObservation.ChangedSecurityFloor ||
		progress.RuntimeObservation.GlobalRuntimeCoverage < 90 ||
		progress.RuntimeObservation.SASTFindings != 0 || progress.RuntimeObservation.SCAActionableFindings != 0 {
		return fmt.Errorf("Hunt Management runtime progress does not satisfy bounded security evidence")
	}
	return nil
}

func validateHuntManagementHandoff(handoff HuntManagementClosureHandoff) error {
	if handoff.SchemaVersion != 1 || handoff.HandoffKind != "hunt-management-bounded-core" ||
		handoff.Readiness != "BOUNDED_HUNT_MANAGEMENT_VALIDATED_NOT_FULL_CAPABILITY" ||
		handoff.Capability != "CAP-INV-005" {
		return fmt.Errorf("invalid Hunt Management bounded handoff identity/readiness")
	}
	if handoff.ProductionReadinessClaim || handoff.FullCapabilityCompletionClaim ||
		handoff.Class2MutationEnabled || handoff.ProviderSelected || handoff.StorageEngineSelected ||
		handoff.RetentionPolicySelected || handoff.CollaborationBackendSelected ||
		handoff.FinalQueryDialectSelected || handoff.DedicatedFinalUIClaim {
		return fmt.Errorf("Hunt Management bounded handoff overclaims scope or readiness")
	}
	if !sameStringSet(handoff.BlockingOpenDecisions, []string{"OPEN-013"}) {
		return fmt.Errorf("Hunt Management handoff must preserve OPEN-013")
	}
	if !sameStringSet(handoff.EvidenceSources, []string{huntManagementRuntimeProgressPath}) {
		return fmt.Errorf("Hunt Management handoff evidence sources mismatch")
	}
	for _, fragment := range huntManagementRequiredLimitations {
		if !containsFoldFragment(handoff.Limitations, fragment) {
			return fmt.Errorf("Hunt Management handoff omits limitation containing %q", fragment)
		}
	}
	if len(handoff.NextConditions) < 4 {
		return fmt.Errorf("Hunt Management handoff requires at least four explicit next conditions")
	}
	return nil
}

func validateHuntManagementRuntimeUnchanged(root, evidenceHead, currentHead string) error {
	evidenceHead, err := validateFullCommitID(evidenceHead)
	if err != nil {
		return fmt.Errorf("Hunt Management runtime evidence head: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return fmt.Errorf("current source head: %w", err)
	}
	if err := ensureGitAncestor(root, evidenceHead, currentHead); err != nil {
		return fmt.Errorf("Hunt Management runtime evidence ancestry: %w", err)
	}
	cmd := exec.Command("git", "-C", root, "diff", "--quiet", evidenceHead, currentHead, "--", huntManagementRuntimeDirectory) // #nosec G204,G702 -- executable/path fixed; commit arguments validated full hex IDs.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("Hunt Management runtime changed after verified runtime evidence; fresh runtime proof is required")
		}
		return fmt.Errorf("compare Hunt Management runtime against verified evidence: %w", err)
	}
	return nil
}

func runHuntManagementAdversarialRuntimeTests(root string) (HuntManagementAdversarialSummary, error) {
	dir, err := resolveRepoPath(root, huntManagementRuntimeDirectory, false)
	if err != nil {
		return HuntManagementAdversarialSummary{}, err
	}
	pattern := "^(" + strings.Join(huntManagementAdversarialTests, "|") + ")$"
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", pattern, ".") // #nosec G204 -- executable fixed; regex assembled from compile-time test names.
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return HuntManagementAdversarialSummary{}, fmt.Errorf("Hunt Management adversarial runtime tests failed: %w", err)
	}
	return parseHuntManagementAdversarialTestEvents(stdout.Bytes())
}

func parseHuntManagementAdversarialTestEvents(data []byte) (HuntManagementAdversarialSummary, error) {
	type testEvent struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
	}
	passed := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return HuntManagementAdversarialSummary{}, fmt.Errorf("decode Hunt Management go test event: %w", err)
		}
		if event.Action == "pass" && containsString(huntManagementAdversarialTests, event.Test) {
			passed[event.Test] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return HuntManagementAdversarialSummary{}, err
	}
	var values []string
	for test := range passed {
		values = append(values, test)
	}
	sort.Strings(values)
	expected := append([]string(nil), huntManagementAdversarialTests...)
	sort.Strings(expected)
	summary := HuntManagementAdversarialSummary{Expected: expected, Passed: values, Count: len(values)}
	if !sameStringSet(values, expected) {
		return summary, fmt.Errorf("Hunt Management adversarial runtime tests incomplete: passed=%v expected=%v", values, expected)
	}
	return summary, nil
}

func huntManagementPerformanceMetric(summary PerformanceBenchmarkAuditSummary, targetID, metricID string) (float64, float64, error) {
	if summary.Status != "pass" {
		return 0, 0, fmt.Errorf("Hunt Management performance benchmark status is %q", summary.Status)
	}
	for _, result := range summary.Results {
		if result.TargetID != targetID {
			continue
		}
		for _, metric := range result.Metrics {
			if metric.ID == metricID {
				if !metric.AbsolutePass || !metric.RelativePass || metric.Budget <= 0 || metric.Observed > metric.Budget {
					return metric.Observed, metric.Budget, fmt.Errorf("Hunt Management performance metric %s failed budget", metricID)
				}
				return metric.Observed, metric.Budget, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("Hunt Management performance result missing %s/%s", targetID, metricID)
}
