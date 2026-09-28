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
	savedQueryAssetsRuntimeProgressPath = "work/lots/E10-INV-006B-RUNTIME/PROGRESS.json"
	savedQueryAssetsHandoffProgressPath = "work/lots/E10-INV-006C-HANDOFF/PROGRESS.json"
	savedQueryAssetsClosureHandoffPath  = "engineering/implementation/saved-query-assets/bounded-handoff.json"
	savedQueryAssetsRuntimeDirectory    = "product-runtime/saved-query-assets"
)

var savedQueryAssetsAdversarialTests = []string{
	"TestProjectMatchesPredeclaredContractFixtures",
	"TestProjectionDeepCopiesCallerOwnedAsset",
	"TestProjectRejectsAdditionalInvalidInputs",
	"TestIncompatibleOverridesStaleAndDiagnosticsAreSorted",
	"TestBuildExecutionHandoffPreservesVersionValuesAndContext",
	"TestBuildExecutionHandoffDeepCopiesProjectionAndParameterValues",
	"TestBuildExecutionHandoffRejectsIneligibleProjection",
	"TestBuildExecutionHandoffRejectsScopeMismatch",
	"TestBuildExecutionHandoffRejectsUnknownAndDuplicateParameterBindings",
	"TestExecutionHandoffDraftContainsNoSearchJobIdentityOrState",
}

var savedQueryAssetsRequiredLimitations = []string{
	"class-2 mutation",
	"canonical Saved Search or Query Asset persistence schema",
	"Saved View",
	"Detection Rule",
	"Search Job creation and execution",
	"storage engine",
	"provider",
	"retention",
	"collaboration backend",
	"final query dialect",
	"dedicated final Saved Search or Query Asset UI",
	"production",
	"full-capability",
}

type SavedQueryAssetsRuntimeProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	Validation         struct {
		PushWorkflowRun        int64  `json:"push_workflow_run"`
		PullRequestWorkflowRun int64  `json:"pull_request_workflow_run"`
		Result                 string `json:"result"`
		SavedQueryAssets       string `json:"saved_query_assets_contract"`
		SecurityTests          string `json:"security_tests"`
		SAST                   string `json:"sast"`
		SCA                    string `json:"sca"`
		RuntimeDependencies    string `json:"runtime_dependencies"`
		BoundaryEdges          string `json:"boundary_edges"`
	} `json:"validation"`
	RuntimeObservation struct {
		Contract                     string  `json:"contract"`
		RuntimeState                 string  `json:"runtime_state"`
		Fixtures                     int     `json:"fixtures"`
		Positive                     int     `json:"positive"`
		Negative                     int     `json:"negative"`
		Incompatible                 int     `json:"incompatible"`
		Stale                        int     `json:"stale"`
		RuntimeDependencies          int     `json:"runtime_dependencies"`
		PushGlobalCoveragePercent    float64 `json:"push_global_coverage_percent"`
		PullRequestCoveragePercent   float64 `json:"pull_request_global_coverage_percent"`
		SASTFindings                 int     `json:"sast_findings"`
		SCAActionableFindings        int     `json:"sca_actionable_findings"`
	} `json:"runtime_observation"`
	ProductSpecMutated bool `json:"product_spec_mutated"`
}

type SavedQueryAssetsHandoffProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	Validation         struct {
		PushWorkflowRun        int64  `json:"push_workflow_run"`
		PullRequestWorkflowRun int64  `json:"pull_request_workflow_run"`
		Result                 string `json:"result"`
	} `json:"validation"`
	HandoffObservation struct {
		Immutable                         bool    `json:"immutable"`
		ExactQueryVersionPreserved        bool    `json:"exact_query_version_preserved"`
		OpaqueParameterValuesPreserved    bool    `json:"opaque_parameter_values_preserved"`
		SourceContextCopied               bool    `json:"source_context_copied"`
		SourceReevaluationRequired        bool    `json:"source_reevaluation_required"`
		PermissionReevaluationRequired    bool    `json:"permission_reevaluation_required"`
		ReturnContextAssetIdentityOnly    bool    `json:"return_context_asset_identity_only"`
		SearchJobCreated                  bool    `json:"search_job_created"`
		SearchJobExecuted                 bool    `json:"search_job_executed"`
		DirectEventSearchRuntimeDependency bool   `json:"direct_event_search_runtime_dependency"`
		PushGlobalCoveragePercent         float64 `json:"push_global_coverage_percent"`
		PullRequestCoveragePercent        float64 `json:"pull_request_global_coverage_percent"`
		SASTFindings                      int     `json:"sast_findings"`
		SCAActionableFindings             int     `json:"sca_actionable_findings"`
	} `json:"handoff_observation"`
	ProductSpecMutated bool `json:"product_spec_mutated"`
}

type SavedQueryAssetsClosureHandoff struct {
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

type SavedQueryAssetsAdversarialSummary struct {
	Expected []string `json:"expected"`
	Passed   []string `json:"passed"`
	Count    int      `json:"count"`
}

type SavedQueryAssetsE2EAuditSummary struct {
	Capability                        string  `json:"capability"`
	ContractID                        string  `json:"contract_id"`
	AdversarialTests                  int     `json:"adversarial_tests"`
	RuntimeCoveragePercent            float64 `json:"runtime_coverage_percent"`
	SASTFindings                      int     `json:"sast_findings"`
	SCAActionableFindings             int     `json:"sca_actionable_findings"`
	ProjectionP95MS                   float64 `json:"projection_p95_ms"`
	ProjectionBudgetMS                float64 `json:"projection_budget_ms"`
	HandoffP95MS                      float64 `json:"handoff_p95_ms"`
	HandoffBudgetMS                   float64 `json:"handoff_budget_ms"`
	RuntimeUnchangedSinceVerification bool    `json:"runtime_unchanged_since_verification"`
	BlockingOpenDecisions             int     `json:"blocking_open_decisions"`
	Limitations                       int     `json:"limitations"`
	ProductionReadinessClaim          bool    `json:"production_readiness_claim"`
	FullCapabilityClaim               bool    `json:"full_capability_claim"`
	Status                            string  `json:"status"`
}

func runSavedQueryAssetsE2EAudit(root, _ string, state CurrentState) (SavedQueryAssetsE2EAuditSummary, error) {
	if _, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, fmt.Errorf("Saved Query Assets Product Spec baseline: %w", err)
	}
	contract, err := runSavedQueryAssetsContractAudit(root)
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	var runtime SavedQueryAssetsRuntimeProgress
	if err := decodeStrict(root, savedQueryAssetsRuntimeProgressPath, &runtime); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	var handoffProgress SavedQueryAssetsHandoffProgress
	if err := decodeStrict(root, savedQueryAssetsHandoffProgressPath, &handoffProgress); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	var closure SavedQueryAssetsClosureHandoff
	if err := decodeStrict(root, savedQueryAssetsClosureHandoffPath, &closure); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	if err := validateSavedQueryAssetsRuntimeProgress(runtime); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	if err := validateSavedQueryAssetsHandoffProgress(handoffProgress); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	if err := validateSavedQueryAssetsClosureHandoff(closure); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}

	currentHead, err := currentSourceCommit(root)
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	if err := validateSavedQueryAssetsRuntimeUnchanged(root, handoffProgress.FinalValidatedHead, currentHead); err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	adversarial, err := runSavedQueryAssetsAdversarialRuntimeTests(root)
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	performance, err := runPerformanceBenchmarkTargets(root, "pr", defaultCIEnvironment, []string{
		"PERF-TGT-SAVED-QUERY-ASSETS-PROJECTION",
		"PERF-TGT-SAVED-QUERY-ASSETS-HANDOFF",
	})
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	projectionObserved, projectionBudget, err := savedQueryAssetsPerformanceMetric(
		performance, "PERF-TGT-SAVED-QUERY-ASSETS-PROJECTION", "query-asset-projection-p95-latency",
	)
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	handoffObserved, handoffBudget, err := savedQueryAssetsPerformanceMetric(
		performance, "PERF-TGT-SAVED-QUERY-ASSETS-HANDOFF", "execution-handoff-p95-latency",
	)
	if err != nil {
		return SavedQueryAssetsE2EAuditSummary{}, err
	}
	if contract.Status != "PASS" || contract.RuntimeState != "implemented" ||
		contract.RuntimeDependencies != 0 || contract.Fixtures != 18 ||
		contract.PositiveFixtures != 5 || contract.NegativeFixtures != 13 ||
		contract.IncompatibleCases != 2 || contract.StaleCases != 1 {
		return SavedQueryAssetsE2EAuditSummary{}, fmt.Errorf("Saved Query Assets contract evidence is incomplete")
	}
	if adversarial.Count != len(savedQueryAssetsAdversarialTests) ||
		!sameStringSet(adversarial.Passed, savedQueryAssetsAdversarialTests) {
		return SavedQueryAssetsE2EAuditSummary{}, fmt.Errorf("Saved Query Assets adversarial evidence is incomplete")
	}

	coverage := runtime.RuntimeObservation.PushGlobalCoveragePercent
	for _, value := range []float64{
		runtime.RuntimeObservation.PullRequestCoveragePercent,
		handoffProgress.HandoffObservation.PushGlobalCoveragePercent,
		handoffProgress.HandoffObservation.PullRequestCoveragePercent,
	} {
		if value < coverage {
			coverage = value
		}
	}
	return SavedQueryAssetsE2EAuditSummary{
		Capability: "CAP-INV-006", ContractID: contract.ContractID,
		AdversarialTests: adversarial.Count, RuntimeCoveragePercent: coverage,
		SASTFindings: runtime.RuntimeObservation.SASTFindings + handoffProgress.HandoffObservation.SASTFindings,
		SCAActionableFindings: runtime.RuntimeObservation.SCAActionableFindings + handoffProgress.HandoffObservation.SCAActionableFindings,
		ProjectionP95MS: projectionObserved, ProjectionBudgetMS: projectionBudget,
		HandoffP95MS: handoffObserved, HandoffBudgetMS: handoffBudget,
		RuntimeUnchangedSinceVerification: true,
		BlockingOpenDecisions: len(closure.BlockingOpenDecisions),
		Limitations: len(closure.Limitations),
		ProductionReadinessClaim: closure.ProductionReadinessClaim,
		FullCapabilityClaim: closure.FullCapabilityCompletionClaim,
		Status: "PASS",
	}, nil
}

func validateSavedQueryAssetsRuntimeProgress(progress SavedQueryAssetsRuntimeProgress) error {
	if progress.SchemaVersion != 1 || progress.WorkUnit != "E10-INV-006B-RUNTIME" ||
		progress.Status != "VERIFIED" || progress.Validation.Result != "PASS" || progress.ProductSpecMutated {
		return fmt.Errorf("Saved Query Assets runtime progress is not verified")
	}
	normalized, err := validateFullCommitID(progress.FinalValidatedHead)
	if err != nil || normalized != progress.FinalValidatedHead {
		return fmt.Errorf("Saved Query Assets runtime progress has invalid final validated head")
	}
	if progress.RuntimeObservation.Contract != "QUERY-ASSET-READONLY-CONTRACT-V1" ||
		progress.RuntimeObservation.RuntimeState != "implemented" ||
		progress.RuntimeObservation.Fixtures != 18 || progress.RuntimeObservation.Positive != 5 ||
		progress.RuntimeObservation.Negative != 13 || progress.RuntimeObservation.Incompatible != 2 ||
		progress.RuntimeObservation.Stale != 1 || progress.RuntimeObservation.RuntimeDependencies != 0 ||
		progress.RuntimeObservation.PushGlobalCoveragePercent < 90 ||
		progress.RuntimeObservation.PullRequestCoveragePercent < 90 ||
		progress.RuntimeObservation.SASTFindings != 0 || progress.RuntimeObservation.SCAActionableFindings != 0 {
		return fmt.Errorf("Saved Query Assets runtime progress does not satisfy bounded security evidence")
	}
	return nil
}

func validateSavedQueryAssetsHandoffProgress(progress SavedQueryAssetsHandoffProgress) error {
	if progress.SchemaVersion != 1 || progress.WorkUnit != "E10-INV-006C-HANDOFF" ||
		progress.Status != "VERIFIED" || progress.Validation.Result != "PASS" || progress.ProductSpecMutated {
		return fmt.Errorf("Saved Query Assets handoff progress is not verified")
	}
	normalized, err := validateFullCommitID(progress.FinalValidatedHead)
	if err != nil || normalized != progress.FinalValidatedHead {
		return fmt.Errorf("Saved Query Assets handoff progress has invalid final validated head")
	}
	h := progress.HandoffObservation
	if !h.Immutable || !h.ExactQueryVersionPreserved || !h.OpaqueParameterValuesPreserved ||
		!h.SourceContextCopied || !h.SourceReevaluationRequired || !h.PermissionReevaluationRequired ||
		!h.ReturnContextAssetIdentityOnly || h.SearchJobCreated || h.SearchJobExecuted ||
		h.DirectEventSearchRuntimeDependency || h.PushGlobalCoveragePercent < 90 ||
		h.PullRequestCoveragePercent < 90 || h.SASTFindings != 0 || h.SCAActionableFindings != 0 {
		return fmt.Errorf("Saved Query Assets handoff progress does not satisfy bounded evidence")
	}
	return nil
}

func validateSavedQueryAssetsClosureHandoff(handoff SavedQueryAssetsClosureHandoff) error {
	if handoff.SchemaVersion != 1 || handoff.HandoffKind != "saved-query-assets-bounded-core" ||
		handoff.Readiness != "BOUNDED_SAVED_QUERY_ASSETS_VALIDATED_NOT_FULL_CAPABILITY" ||
		handoff.Capability != "CAP-INV-006" {
		return fmt.Errorf("invalid Saved Query Assets bounded handoff identity/readiness")
	}
	if handoff.ProductionReadinessClaim || handoff.FullCapabilityCompletionClaim ||
		handoff.Class2MutationEnabled || handoff.ProviderSelected || handoff.StorageEngineSelected ||
		handoff.RetentionPolicySelected || handoff.CollaborationBackendSelected ||
		handoff.FinalQueryDialectSelected || handoff.DedicatedFinalUIClaim {
		return fmt.Errorf("Saved Query Assets bounded handoff overclaims scope or readiness")
	}
	if !sameStringSet(handoff.BlockingOpenDecisions, []string{"OPEN-013"}) {
		return fmt.Errorf("Saved Query Assets handoff must preserve OPEN-013")
	}
	if !sameStringSet(handoff.EvidenceSources, []string{
		savedQueryAssetsRuntimeProgressPath, savedQueryAssetsHandoffProgressPath,
	}) {
		return fmt.Errorf("Saved Query Assets handoff evidence sources mismatch")
	}
	for _, fragment := range savedQueryAssetsRequiredLimitations {
		if !containsFoldFragment(handoff.Limitations, fragment) {
			return fmt.Errorf("Saved Query Assets handoff omits limitation containing %q", fragment)
		}
	}
	if len(handoff.NextConditions) < 5 {
		return fmt.Errorf("Saved Query Assets handoff requires at least five explicit next conditions")
	}
	return nil
}

func validateSavedQueryAssetsRuntimeUnchanged(root, evidenceHead, currentHead string) error {
	evidenceHead, err := validateFullCommitID(evidenceHead)
	if err != nil {
		return fmt.Errorf("Saved Query Assets handoff evidence head: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return fmt.Errorf("current source head: %w", err)
	}
	if err := ensureGitAncestor(root, evidenceHead, currentHead); err != nil {
		return fmt.Errorf("Saved Query Assets runtime evidence ancestry: %w", err)
	}
	cmd := exec.Command("git", "-C", root, "diff", "--quiet", evidenceHead, currentHead, "--", savedQueryAssetsRuntimeDirectory) // #nosec G204,G702 -- executable/path fixed; commit args validated full hex IDs.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("Saved Query Assets runtime changed after verified handoff evidence; fresh runtime proof is required")
		}
		return fmt.Errorf("compare Saved Query Assets runtime against verified handoff evidence: %w", err)
	}
	return nil
}

func runSavedQueryAssetsAdversarialRuntimeTests(root string) (SavedQueryAssetsAdversarialSummary, error) {
	dir, err := resolveRepoPath(root, savedQueryAssetsRuntimeDirectory, false)
	if err != nil {
		return SavedQueryAssetsAdversarialSummary{}, err
	}
	pattern := "^(" + strings.Join(savedQueryAssetsAdversarialTests, "|") + ")$"
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", pattern, ".") // #nosec G204 -- executable fixed; regex assembled from compile-time test names.
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return SavedQueryAssetsAdversarialSummary{}, fmt.Errorf("Saved Query Assets adversarial runtime tests failed: %w", err)
	}
	return parseSavedQueryAssetsAdversarialTestEvents(stdout.Bytes())
}

func parseSavedQueryAssetsAdversarialTestEvents(data []byte) (SavedQueryAssetsAdversarialSummary, error) {
	type testEvent struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
	}
	passed := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return SavedQueryAssetsAdversarialSummary{}, fmt.Errorf("decode Saved Query Assets go test event: %w", err)
		}
		if event.Action == "pass" && containsString(savedQueryAssetsAdversarialTests, event.Test) {
			passed[event.Test] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return SavedQueryAssetsAdversarialSummary{}, err
	}
	var values []string
	for test := range passed {
		values = append(values, test)
	}
	sort.Strings(values)
	expected := append([]string(nil), savedQueryAssetsAdversarialTests...)
	sort.Strings(expected)
	summary := SavedQueryAssetsAdversarialSummary{Expected: expected, Passed: values, Count: len(values)}
	if !sameStringSet(values, expected) {
		return summary, fmt.Errorf("Saved Query Assets adversarial runtime tests incomplete: passed=%v expected=%v", values, expected)
	}
	return summary, nil
}

func savedQueryAssetsPerformanceMetric(summary PerformanceBenchmarkAuditSummary, targetID, metricID string) (float64, float64, error) {
	if summary.Status != "pass" {
		return 0, 0, fmt.Errorf("Saved Query Assets performance benchmark status is %q", summary.Status)
	}
	for _, result := range summary.Results {
		if result.TargetID != targetID {
			continue
		}
		for _, metric := range result.Metrics {
			if metric.ID == metricID {
				if !metric.AbsolutePass || !metric.RelativePass || metric.Budget <= 0 || metric.Observed > metric.Budget {
					return metric.Observed, metric.Budget, fmt.Errorf("Saved Query Assets performance metric %s failed budget", metricID)
				}
				return metric.Observed, metric.Budget, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("Saved Query Assets performance result missing %s/%s", targetID, metricID)
}
