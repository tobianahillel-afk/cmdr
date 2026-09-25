package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	eventSearchClosureHandoffPath  = "engineering/implementation/event-search/bounded-handoff.json"
	eventSearchValidationProgress  = "work/lots/E10-INV-002B-RUNTIME/PROGRESS.json"
	eventSearchRuntimeProgress     = "work/lots/E10-INV-002C-RUNTIME/PROGRESS.json"
	eventSearchPerformanceProgress = "work/lots/E10-INV-002C-PERF/PROGRESS.json"
	eventSearchFrontendProgress    = "work/lots/E10-INV-002F-E2E/PROGRESS.json"
	eventSearchFrontendHandoff     = "work/lots/E10-INV-002F-E2E/HANDOFF.json"
	eventSearchClosureWorkUnit     = "E10-INV-002D"
	eventSearchRuntimeDirectory    = "product-runtime/event-search"
)

var eventSearchAdversarialTests = []string{
	"TestValidateContractFixtures",
	"TestSearchJobRejectsCrossTenantBackendResultWithoutLeakingReferences",
	"TestSearchJobBackendMustAccountForEveryRequestedSource",
	"TestSearchJobPartialPreservesValidReferencesAndFailedSources",
	"TestSearchJobCancellationDoesNotFabricateCompletion",
	"TestSearchJobContextCancellationProducesCancelledState",
	"TestSearchJobRetryCreatesNewIdentityAndPreservesProvenance",
	"TestSearchJobBackendErrorFailsWithoutPersistingErrorText",
	"TestValidateBackendResultRejectsMalformedEvidence",
}

var eventSearchRequiredLimitations = []string{
	"Saved Search",
	"Case-link",
	"final query language",
	"final search/index/storage/provider",
	"production search SLO",
}

type EventSearchClosureHandoff struct {
	SchemaVersion            int      `json:"schema_version"`
	HandoffKind              string   `json:"handoff_kind"`
	Readiness                string   `json:"readiness"`
	Capability               string   `json:"capability"`
	Screen                   string   `json:"screen"`
	ProductionReadinessClaim bool     `json:"production_readiness_claim"`
	BlockingOpenDecisions    []string `json:"blocking_open_decisions"`
	EvidenceSources          []string `json:"evidence_sources"`
	Limitations              []string `json:"limitations"`
	NextConditions           []string `json:"next_conditions"`
}

type EventSearchValidationEvidence struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	Runtime            struct {
		ExternalRuntimeDependencies int `json:"external_runtime_dependencies"`
	} `json:"runtime"`
	Security struct {
		GlobalRuntimeCoveragePercent float64 `json:"global_runtime_coverage_percent"`
		SASTFindings                 int     `json:"sast_findings"`
		SCAActionableFindings        int     `json:"sca_actionable_findings"`
		ThirdPartyRuntimeComponents  int     `json:"third_party_runtime_components"`
	} `json:"security"`
	Performance struct {
		ObservedMS   float64 `json:"observed_ms"`
		BudgetMS     float64 `json:"budget_ms"`
		AbsolutePass bool    `json:"absolute_pass"`
	} `json:"performance"`
}

type EventSearchRuntimeEvidence struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	Runtime            struct {
		ExternalRuntimeDependencies int      `json:"external_runtime_dependencies"`
		LifecycleStates             []string `json:"lifecycle_states"`
		TerminalRetry               string   `json:"terminal_retry"`
		BackendKind                 string   `json:"backend_kind"`
		ProviderSelected            bool     `json:"provider_selected"`
		IndexSelected               bool     `json:"index_selected"`
		StorageSelected             bool     `json:"storage_selected"`
	} `json:"runtime"`
}

type EventSearchPerformanceEvidence struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	Validation         struct {
		Result string `json:"result"`
	} `json:"validation"`
	Performance struct {
		ObservedMS               float64 `json:"observed_ms"`
		BudgetMS                 float64 `json:"budget_ms"`
		AbsolutePass             bool    `json:"absolute_pass"`
		BackendExecutionIncluded bool    `json:"backend_execution_included"`
		ProductionSLOClaimed     bool    `json:"production_slo_claimed"`
		ValidationP95MS          float64 `json:"validation_p95_ms"`
	} `json:"performance"`
}

type EventSearchFrontendE2EProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	Validation         struct {
		Result string `json:"result"`
	} `json:"validation"`
	FrontendSecurity struct {
		GlobalRuntimeCoveragePercent float64 `json:"global_runtime_coverage_percent"`
		AuthorizationNegative        string  `json:"authorization_negative"`
		TenantIsolationNegative      string  `json:"tenant_isolation_negative"`
		HostileDOMContent            string  `json:"hostile_dom_content"`
		PermissionDenialClears       string  `json:"permission_denial_clears_projection"`
		CancellationStaleBlocking    string  `json:"cancellation_stale_result_blocking"`
		OfflineSemanticState         string  `json:"offline_semantic_state"`
	} `json:"frontend_security"`
	Performance struct {
		BudgetMS              float64 `json:"budget_ms"`
		PushObservedMS        float64 `json:"push_observed_ms"`
		PullRequestObservedMS float64 `json:"pull_request_observed_ms"`
		AbsolutePass          bool    `json:"absolute_pass"`
	} `json:"performance"`
	ProductSpecMutated bool `json:"product_spec_mutated"`
}

type EventSearchFrontendE2EHandoff struct {
	SchemaVersion                 int    `json:"schema_version"`
	WorkUnit                      string `json:"work_unit"`
	Result                        string `json:"result"`
	FinalValidatedHead            string `json:"final_validated_head"`
	PullRequest                   int    `json:"pull_request"`
	ProductionReadinessClaim      bool   `json:"production_readiness_claim"`
	FullCapabilityCompletionClaim bool   `json:"full_capability_completion_claim"`
	ProductSpecMutated            bool   `json:"product_spec_mutated"`
	Evidence                      struct {
		PushRun                      int     `json:"push_run"`
		PullRequestRun               int     `json:"pull_request_run"`
		NodeTestFiles                int     `json:"node_test_files"`
		ProductionModules            int     `json:"production_modules"`
		GlobalRuntimeCoveragePercent float64 `json:"global_runtime_coverage_percent"`
		FrontendP95PushMS            float64 `json:"frontend_p95_push_ms"`
		FrontendP95PullRequestMS     float64 `json:"frontend_p95_pull_request_ms"`
		FrontendP95BudgetMS          float64 `json:"frontend_p95_budget_ms"`
		RuntimeDependencies          int     `json:"runtime_dependencies"`
	} `json:"evidence"`
}

type EventSearchAdversarialSummary struct {
	Expected []string `json:"expected"`
	Passed   []string `json:"passed"`
	Count    int      `json:"count"`
}

type EventSearchE2EAuditSummary struct {
	Capability                string  `json:"capability"`
	Screen                    string  `json:"screen"`
	ContractID                string  `json:"contract_id"`
	ContractFixtures          int     `json:"contract_fixtures"`
	NegativeFixtures          int     `json:"negative_fixtures"`
	AdversarialTests          int     `json:"adversarial_tests"`
	RuntimeEvidenceHead       string  `json:"runtime_evidence_head"`
	PerformanceEvidenceHead   string  `json:"performance_evidence_head"`
	CurrentHead               string  `json:"current_head"`
	RuntimeUnchangedSincePerf bool    `json:"runtime_unchanged_since_performance_evidence"`
	RuntimeCoveragePercent    float64 `json:"runtime_coverage_percent"`
	SASTFindings              int     `json:"sast_findings"`
	SCAActionableFindings     int     `json:"sca_actionable_findings"`
	ValidationP95MS           float64 `json:"validation_p95_ms"`
	ValidationBudgetMS        float64 `json:"validation_budget_ms"`
	OrchestrationP95MS        float64 `json:"orchestration_p95_ms"`
	OrchestrationBudgetMS     float64 `json:"orchestration_budget_ms"`
	FrontendEvidenceHead      string  `json:"frontend_evidence_head"`
	FrontendCoveragePercent   float64 `json:"frontend_coverage_percent"`
	FrontendP95MS             float64 `json:"frontend_p95_ms"`
	FrontendP95BudgetMS       float64 `json:"frontend_p95_budget_ms"`
	RecoveryOutcome           string  `json:"recovery_outcome"`
	RecoveryNextAction        string  `json:"recovery_next_action"`
	BlockingOpenDecisions     int     `json:"blocking_open_decisions"`
	ProductionReadinessClaim  bool    `json:"production_readiness_claim"`
	Limitations               int     `json:"limitations"`
	Status                    string  `json:"status"`
}

func runEventSearchE2EAudit(root string, state CurrentState, graph WorkGraph) (EventSearchE2EAuditSummary, error) {
	contract, err := runEventSearchContractAudit(root)
	if err != nil {
		return EventSearchE2EAuditSummary{}, fmt.Errorf("event search contract: %w", err)
	}
	baseline, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true)
	if err != nil {
		return EventSearchE2EAuditSummary{}, fmt.Errorf("product baseline: %w", err)
	}
	if baseline.ProductSpecBaselineCommit != state.ProductSpec.BaselineCommit {
		return EventSearchE2EAuditSummary{}, fmt.Errorf("product baseline commit mismatch")
	}

	var handoff EventSearchClosureHandoff
	if err := decodeStrict(root, eventSearchClosureHandoffPath, &handoff); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	var validation EventSearchValidationEvidence
	if err := decodeEventSearchEvidence(root, eventSearchValidationProgress, &validation); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	var runtime EventSearchRuntimeEvidence
	if err := decodeEventSearchEvidence(root, eventSearchRuntimeProgress, &runtime); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	var performance EventSearchPerformanceEvidence
	if err := decodeEventSearchEvidence(root, eventSearchPerformanceProgress, &performance); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}

	var frontendProgress EventSearchFrontendE2EProgress
	if err := decodeEventSearchEvidence(root, eventSearchFrontendProgress, &frontendProgress); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	var frontendHandoff EventSearchFrontendE2EHandoff
	if err := decodeEventSearchEvidence(root, eventSearchFrontendHandoff, &frontendHandoff); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}

	currentHead, err := currentSourceCommit(root)
	if err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	if err := validateEventSearchRuntimeUnchanged(root, performance.FinalValidatedHead, currentHead); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	adversarial, err := runEventSearchAdversarialRuntimeTests(root)
	if err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	recovery, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{WorkUnit: eventSearchClosureWorkUnit}, state, graph)
	if err != nil {
		return EventSearchE2EAuditSummary{}, fmt.Errorf("recovery reconciliation: %w", err)
	}
	if recovery.Outcome == "conflict" {
		return EventSearchE2EAuditSummary{}, fmt.Errorf("recovery reconciliation conflict: %s", strings.Join(recovery.Reasons, ", "))
	}

	if err := validateEventSearchClosureEvidence(handoff, validation, runtime, performance, frontendProgress, frontendHandoff, contract, adversarial, recovery); err != nil {
		return EventSearchE2EAuditSummary{}, err
	}
	return EventSearchE2EAuditSummary{
		Capability: "CAP-INV-002", Screen: "INV-EVS-001",
		ContractID: contract.ContractID, ContractFixtures: contract.Fixtures, NegativeFixtures: contract.NegativeFixtures,
		AdversarialTests:    adversarial.Count,
		RuntimeEvidenceHead: runtime.FinalValidatedHead, PerformanceEvidenceHead: performance.FinalValidatedHead,
		CurrentHead: currentHead, RuntimeUnchangedSincePerf: true,
		RuntimeCoveragePercent: validation.Security.GlobalRuntimeCoveragePercent,
		SASTFindings:           validation.Security.SASTFindings, SCAActionableFindings: validation.Security.SCAActionableFindings,
		ValidationP95MS: performance.Performance.ValidationP95MS, ValidationBudgetMS: 1,
		OrchestrationP95MS: performance.Performance.ObservedMS, OrchestrationBudgetMS: performance.Performance.BudgetMS,
		FrontendEvidenceHead: frontendProgress.FinalValidatedHead,
		FrontendCoveragePercent: frontendProgress.FrontendSecurity.GlobalRuntimeCoveragePercent,
		FrontendP95MS: frontendProgress.Performance.PullRequestObservedMS,
		FrontendP95BudgetMS: frontendProgress.Performance.BudgetMS,
		RecoveryOutcome: recovery.Outcome, RecoveryNextAction: recovery.NextAction,
		BlockingOpenDecisions:    len(handoff.BlockingOpenDecisions),
		ProductionReadinessClaim: handoff.ProductionReadinessClaim,
		Limitations:              len(handoff.Limitations), Status: "PASS",
	}, nil
}

func decodeEventSearchEvidence(root, path string, dst any) error {
	data, err := readRepoFile(root, path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func validateEventSearchClosureEvidence(
	handoff EventSearchClosureHandoff,
	validation EventSearchValidationEvidence,
	runtime EventSearchRuntimeEvidence,
	performance EventSearchPerformanceEvidence,
	frontendProgress EventSearchFrontendE2EProgress,
	frontendHandoff EventSearchFrontendE2EHandoff,
	contract EventSearchContractAuditSummary,
	adversarial EventSearchAdversarialSummary,
	recovery RecoveryReconciliation,
) error {
	if handoff.SchemaVersion != 1 || handoff.HandoffKind != "event-search-bounded-core" ||
		handoff.Readiness != "BOUNDED_CORE_AND_FRONTEND_VALIDATED_NOT_FULL_CAPABILITY" ||
		handoff.Capability != "CAP-INV-002" || handoff.Screen != "INV-EVS-001" {
		return fmt.Errorf("invalid Event Search bounded handoff identity/readiness")
	}
	if handoff.ProductionReadinessClaim {
		return fmt.Errorf("bounded Event Search handoff must not claim production readiness")
	}
	if !sameStringSet(handoff.BlockingOpenDecisions, []string{"OPEN-013"}) {
		return fmt.Errorf("Event Search handoff must preserve OPEN-013 as the explicit Case-link blocker")
	}
	expectedSources := []string{
		eventSearchValidationProgress,
		eventSearchRuntimeProgress,
		eventSearchPerformanceProgress,
		eventSearchFrontendProgress,
		eventSearchFrontendHandoff,
	}
	if !sameStringSet(handoff.EvidenceSources, expectedSources) {
		return fmt.Errorf("Event Search handoff evidence sources mismatch")
	}
	for _, fragment := range eventSearchRequiredLimitations {
		if !containsFoldFragment(handoff.Limitations, fragment) {
			return fmt.Errorf("Event Search handoff omits limitation containing %q", fragment)
		}
	}
	if len(handoff.NextConditions) < 4 {
		return fmt.Errorf("Event Search handoff requires at least four explicit next conditions")
	}

	if err := validateEventSearchProgressIdentity(validation.WorkUnit, validation.Status, validation.FinalValidatedHead, validation.ProductSpecMutated, "E10-INV-002B-RUNTIME"); err != nil {
		return err
	}
	if validation.Runtime.ExternalRuntimeDependencies != 0 ||
		validation.Security.GlobalRuntimeCoveragePercent < 80 ||
		validation.Security.SASTFindings != 0 ||
		validation.Security.SCAActionableFindings != 0 ||
		validation.Security.ThirdPartyRuntimeComponents != 0 ||
		!validation.Performance.AbsolutePass ||
		validation.Performance.BudgetMS <= 0 ||
		validation.Performance.ObservedMS > validation.Performance.BudgetMS {
		return fmt.Errorf("Event Search validation evidence does not satisfy bounded security/performance floors")
	}

	if err := validateEventSearchProgressIdentity(runtime.WorkUnit, runtime.Status, runtime.FinalValidatedHead, runtime.ProductSpecMutated, "E10-INV-002C-RUNTIME"); err != nil {
		return err
	}
	wantStates := []string{"queued", "running", "completed", "partial", "failed", "cancelled"}
	if runtime.Runtime.ExternalRuntimeDependencies != 0 ||
		!sameStringSet(runtime.Runtime.LifecycleStates, wantStates) ||
		runtime.Runtime.TerminalRetry != "new job and run identity" ||
		runtime.Runtime.BackendKind != "interface-only" ||
		runtime.Runtime.ProviderSelected || runtime.Runtime.IndexSelected || runtime.Runtime.StorageSelected {
		return fmt.Errorf("Event Search runtime evidence violates bounded backend-neutral lifecycle invariants")
	}

	if err := validateEventSearchProgressIdentity(performance.WorkUnit, performance.Status, performance.FinalValidatedHead, performance.ProductSpecMutated, "E10-INV-002C-PERF"); err != nil {
		return err
	}
	if performance.Validation.Result != "PASS" ||
		!performance.Performance.AbsolutePass ||
		performance.Performance.BudgetMS <= 0 ||
		performance.Performance.ObservedMS > performance.Performance.BudgetMS ||
		performance.Performance.BackendExecutionIncluded ||
		performance.Performance.ProductionSLOClaimed ||
		performance.Performance.ValidationP95MS > 1 {
		return fmt.Errorf("Event Search performance evidence is incomplete or outside bounded budgets")
	}

	if err := validateEventSearchProgressIdentity(
		frontendProgress.WorkUnit,
		frontendProgress.Status,
		frontendProgress.FinalValidatedHead,
		frontendProgress.ProductSpecMutated,
		"E10-INV-002F-E2E",
	); err != nil {
		return err
	}
	if frontendProgress.Validation.Result != "PASS" ||
		frontendProgress.FrontendSecurity.GlobalRuntimeCoveragePercent < 80 ||
		frontendProgress.FrontendSecurity.AuthorizationNegative != "PASS" ||
		frontendProgress.FrontendSecurity.TenantIsolationNegative != "PASS" ||
		frontendProgress.FrontendSecurity.HostileDOMContent != "PASS" ||
		frontendProgress.FrontendSecurity.PermissionDenialClears != "PASS" ||
		frontendProgress.FrontendSecurity.CancellationStaleBlocking != "PASS" ||
		frontendProgress.FrontendSecurity.OfflineSemanticState != "PASS" ||
		!frontendProgress.Performance.AbsolutePass ||
		frontendProgress.Performance.BudgetMS <= 0 ||
		frontendProgress.Performance.PushObservedMS > frontendProgress.Performance.BudgetMS ||
		frontendProgress.Performance.PullRequestObservedMS > frontendProgress.Performance.BudgetMS {
		return fmt.Errorf("Event Search frontend E2E progress does not satisfy security/performance floors")
	}
	if err := validateEventSearchProgressIdentity(
		frontendHandoff.WorkUnit,
		frontendHandoff.Result,
		frontendHandoff.FinalValidatedHead,
		frontendHandoff.ProductSpecMutated,
		"E10-INV-002F-E2E",
	); err != nil {
		return err
	}
	if frontendHandoff.PullRequest < 1 ||
		frontendHandoff.ProductionReadinessClaim ||
		frontendHandoff.FullCapabilityCompletionClaim ||
		frontendHandoff.Evidence.PushRun < 1 ||
		frontendHandoff.Evidence.PullRequestRun < 1 ||
		frontendHandoff.Evidence.NodeTestFiles < 5 ||
		frontendHandoff.Evidence.ProductionModules != 4 ||
		frontendHandoff.Evidence.GlobalRuntimeCoveragePercent < 80 ||
		frontendHandoff.Evidence.FrontendP95BudgetMS <= 0 ||
		frontendHandoff.Evidence.FrontendP95PushMS > frontendHandoff.Evidence.FrontendP95BudgetMS ||
		frontendHandoff.Evidence.FrontendP95PullRequestMS > frontendHandoff.Evidence.FrontendP95BudgetMS ||
		frontendHandoff.Evidence.RuntimeDependencies != 0 {
		return fmt.Errorf("Event Search frontend E2E handoff is incomplete or overclaims readiness")
	}
	if contract.Status != "PASS" || contract.RuntimeState != "implemented" ||
		contract.RuntimeDependencies != 0 || contract.Fixtures < 1 || contract.NegativeFixtures < 1 {
		return fmt.Errorf("Event Search executable contract evidence is incomplete")
	}
	if adversarial.Count != len(eventSearchAdversarialTests) || !sameStringSet(adversarial.Passed, eventSearchAdversarialTests) {
		return fmt.Errorf("Event Search adversarial runtime evidence incomplete")
	}
	if recovery.Outcome == "conflict" || strings.TrimSpace(recovery.NextAction) == "" {
		return fmt.Errorf("Event Search recovery evidence is unsafe")
	}
	return nil
}

func validateEventSearchProgressIdentity(workUnit, status, head string, productSpecMutated bool, expected string) error {
	if workUnit != expected || status != "VERIFIED" {
		return fmt.Errorf("%s progress is not VERIFIED", expected)
	}
	normalized, err := validateFullCommitID(head)
	if err != nil || normalized != head {
		return fmt.Errorf("%s progress has invalid final validated head", expected)
	}
	if productSpecMutated {
		return fmt.Errorf("%s claims Product Spec mutation", expected)
	}
	return nil
}

func validateEventSearchRuntimeUnchanged(root, evidenceHead, currentHead string) error {
	evidenceHead, err := validateFullCommitID(evidenceHead)
	if err != nil {
		return fmt.Errorf("Event Search performance evidence head: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return fmt.Errorf("current source head: %w", err)
	}
	if err := ensureGitAncestor(root, evidenceHead, currentHead); err != nil {
		return fmt.Errorf("Event Search performance evidence ancestry: %w", err)
	}
	diff := exec.Command("git", "-C", root, "diff", "--quiet", evidenceHead, currentHead, "--", eventSearchRuntimeDirectory) // #nosec G204,G702 -- executable/path are fixed and commit arguments are validated full hex IDs.
	if err := diff.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("Event Search runtime changed after measured performance evidence; fresh performance proof is required")
		}
		return fmt.Errorf("compare Event Search runtime against performance evidence: %w", err)
	}
	return nil
}

func runEventSearchAdversarialRuntimeTests(root string) (EventSearchAdversarialSummary, error) {
	pattern := "^(" + strings.Join(eventSearchAdversarialTests, "|") + ")$"
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", pattern, ".") // #nosec G204 -- executable is fixed and pattern is assembled only from compile-time test names.
	cmd.Dir = filepath.Join(root, eventSearchRuntimeDirectory)
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return EventSearchAdversarialSummary{}, fmt.Errorf("Event Search adversarial runtime tests failed: %w", err)
	}
	return parseEventSearchAdversarialTestEvents(stdout.Bytes())
}

func parseEventSearchAdversarialTestEvents(data []byte) (EventSearchAdversarialSummary, error) {
	type testEvent struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
	}
	passed := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return EventSearchAdversarialSummary{}, fmt.Errorf("decode Event Search go test event: %w", err)
		}
		if event.Action == "pass" && containsString(eventSearchAdversarialTests, event.Test) {
			passed[event.Test] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return EventSearchAdversarialSummary{}, err
	}
	var values []string
	for test := range passed {
		values = append(values, test)
	}
	sort.Strings(values)
	expected := append([]string(nil), eventSearchAdversarialTests...)
	sort.Strings(expected)
	summary := EventSearchAdversarialSummary{Expected: expected, Passed: values, Count: len(values)}
	if !sameStringSet(values, expected) {
		return summary, fmt.Errorf("Event Search adversarial runtime tests incomplete: passed=%v expected=%v", values, expected)
	}
	return summary, nil
}

func containsFoldFragment(values []string, fragment string) bool {
	fragment = strings.ToLower(fragment)
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), fragment) {
			return true
		}
	}
	return false
}
