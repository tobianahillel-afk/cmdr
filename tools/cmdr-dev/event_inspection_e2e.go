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
	eventInspectionRuntimeProgressPath = "work/lots/E10-INV-004B-RUNTIME/PROGRESS.json"
	eventInspectionPivotProgressPath   = "work/lots/E10-INV-004C-PIVOT/PROGRESS.json"
	eventInspectionHandoffPath         = "engineering/implementation/event-inspection/bounded-handoff.json"
	eventInspectionRuntimeDirectory    = "product-runtime/event-inspection"
)

var eventInspectionAdversarialTests = []string{
	"TestProjectMatchesPredeclaredContractFixtures",
	"TestRawDeniedNeverLeaksOrSynthesizesRawContent",
	"TestProjectionMarksStaleEnrichmentAndPreservesProvenance",
	"TestProjectRejectsAdditionalFailClosedInputs",
	"TestProjectRejectsInvalidEnrichmentProvenance",
	"TestPreparePivotMatchesPredeclaredPivotFixtures",
	"TestPreparePivotRejectsRestrictedValueBeforeLookup",
	"TestPreparePivotRequiresValueFromVisibleProjection",
	"TestPreparePivotIsDeterministicAndDialectNeutral",
	"TestPreparePivotPreservesEnrichmentProvenance",
	"TestPreparePivotRejectsMalformedContextBranches",
	"TestPreparePivotTracksNormalizedAndRawDerivedOrigins",
	"TestValidPivotValueRejectsOversizedValue",
}

var eventInspectionRequiredLimitations = []string{
	"case-link mutation",
	"artifact proposal",
	"evidence-candidate",
	"parser engine",
	"entity-resolution",
	"storage",
	"provider",
	"final query dialect",
	"production",
}

type EventInspectionRuntimeProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	RuntimeEvidence    struct {
		Contract                             string  `json:"contract"`
		RuntimeState                         string  `json:"runtime_state"`
		Fixtures                            int     `json:"fixtures"`
		PositiveFixtures                    int     `json:"positive_fixtures"`
		NegativeFixtures                    int     `json:"negative_fixtures"`
		ExternalRuntimeDependencies         int     `json:"external_runtime_dependencies"`
		CrossBoundaryEdges                  int     `json:"cross_boundary_edges"`
		GlobalRuntimeSecurityCoveragePercent float64 `json:"global_runtime_security_coverage_percent"`
		GlobalCoverageFloorPercent           float64 `json:"global_coverage_floor_percent"`
		ChangedSecurityCriticalFloorPercent  float64 `json:"changed_security_critical_floor_percent"`
	} `json:"runtime_evidence"`
	StaticSecurity struct {
		Findings           int `json:"findings"`
		ActionableFindings int `json:"actionable_findings"`
	} `json:"static_security"`
	ProductSpecMutated bool `json:"product_spec_mutated"`
}

type EventInspectionPivotProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	RuntimeEvidence    struct {
		Contract                                  string  `json:"contract"`
		Fixtures                                 int     `json:"fixtures"`
		PositiveFixtures                         int     `json:"positive_fixtures"`
		NegativeFixtures                         int     `json:"negative_fixtures"`
		ExternalRuntimeDependencies              int     `json:"external_runtime_dependencies"`
		GlobalRuntimeSecurityCoveragePushPercent float64 `json:"global_runtime_security_coverage_push_percent"`
		GlobalRuntimeSecurityCoveragePRPercent   float64 `json:"global_runtime_security_coverage_pr_percent"`
		GlobalCoverageFloorPercent                float64 `json:"global_coverage_floor_percent"`
		ChangedSecurityCriticalFloorPercent       float64 `json:"changed_security_critical_floor_percent"`
	} `json:"runtime_evidence"`
	StaticSecurity struct {
		Findings           int `json:"findings"`
		ActionableFindings int `json:"actionable_findings"`
	} `json:"static_security"`
	ProductSpecMutated bool `json:"product_spec_mutated"`
}

type EventInspectionClosureHandoff struct {
	SchemaVersion                 int      `json:"schema_version"`
	HandoffKind                   string   `json:"handoff_kind"`
	Readiness                     string   `json:"readiness"`
	Capability                    string   `json:"capability"`
	ProductionReadinessClaim      bool     `json:"production_readiness_claim"`
	FullCapabilityCompletionClaim bool     `json:"full_capability_completion_claim"`
	Class2MutationEnabled         bool     `json:"class2_mutation_enabled"`
	ProviderSelected              bool     `json:"provider_selected"`
	FinalQueryDialectSelected     bool     `json:"final_query_dialect_selected"`
	DedicatedFinalUIClaim         bool     `json:"dedicated_final_ui_claim"`
	BlockingOpenDecisions         []string `json:"blocking_open_decisions"`
	EvidenceSources               []string `json:"evidence_sources"`
	Limitations                   []string `json:"limitations"`
	NextConditions                []string `json:"next_conditions"`
}

type EventInspectionAdversarialSummary struct {
	Expected []string `json:"expected"`
	Passed   []string `json:"passed"`
	Count    int      `json:"count"`
}

type EventInspectionE2EAuditSummary struct {
	Capability                 string  `json:"capability"`
	ContractID                 string  `json:"contract_id"`
	AdversarialTests           int     `json:"adversarial_tests"`
	RuntimeCoveragePercent     float64 `json:"runtime_coverage_percent"`
	SASTFindings               int     `json:"sast_findings"`
	SCAActionableFindings      int     `json:"sca_actionable_findings"`
	ProjectionP95MS            float64 `json:"projection_p95_ms"`
	ProjectionBudgetMS         float64 `json:"projection_budget_ms"`
	PivotP95MS                 float64 `json:"pivot_p95_ms"`
	PivotBudgetMS              float64 `json:"pivot_budget_ms"`
	RuntimeUnchangedSincePivot bool    `json:"runtime_unchanged_since_pivot"`
	BlockingOpenDecisions      int     `json:"blocking_open_decisions"`
	Limitations                int     `json:"limitations"`
	ProductionReadinessClaim   bool    `json:"production_readiness_claim"`
	FullCapabilityClaim        bool    `json:"full_capability_claim"`
	Status                     string  `json:"status"`
}

func runEventInspectionE2EAudit(root, changesFile string, state CurrentState) (EventInspectionE2EAuditSummary, error) {
	if _, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true); err != nil {
		return EventInspectionE2EAuditSummary{}, fmt.Errorf("Event Inspection Product Spec baseline: %w", err)
	}
	contract, err := runEventInspectionContractAudit(root)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}

	var runtime EventInspectionRuntimeProgress
	if err := decodeStrict(root, eventInspectionRuntimeProgressPath, &runtime); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	var pivot EventInspectionPivotProgress
	if err := decodeStrict(root, eventInspectionPivotProgressPath, &pivot); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	var handoff EventInspectionClosureHandoff
	if err := decodeStrict(root, eventInspectionHandoffPath, &handoff); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	if err := validateEventInspectionProgress(runtime, pivot); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	if err := validateEventInspectionHandoff(handoff); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}

	currentHead, err := currentSourceCommit(root)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	if err := validateEventInspectionRuntimeUnchanged(root, pivot.FinalValidatedHead, currentHead); err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	adversarial, err := runEventInspectionAdversarialRuntimeTests(root)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}

	performance, err := runPerformanceBenchmarkAudit(root, "pr", defaultCIEnvironment, changesFile)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	projectionObserved, projectionBudget, err := eventInspectionPerformanceMetric(
		performance, "PERF-TGT-EVENT-INSPECTION-PROJECTION", "inspection-projection-p95-latency",
	)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}
	pivotObserved, pivotBudget, err := eventInspectionPerformanceMetric(
		performance, "PERF-TGT-EVENT-INSPECTION-PIVOT", "pivot-draft-p95-latency",
	)
	if err != nil {
		return EventInspectionE2EAuditSummary{}, err
	}

	if contract.Status != "PASS" || contract.RuntimeState != "implemented" ||
		contract.RuntimeDependencies != 0 || contract.Fixtures != 17 ||
		contract.PositiveFixtures != 4 || contract.NegativeFixtures != 13 {
		return EventInspectionE2EAuditSummary{}, fmt.Errorf("Event Inspection contract evidence is incomplete")
	}
	if adversarial.Count != len(eventInspectionAdversarialTests) ||
		!sameStringSet(adversarial.Passed, eventInspectionAdversarialTests) {
		return EventInspectionE2EAuditSummary{}, fmt.Errorf("Event Inspection adversarial evidence is incomplete")
	}

	coverage := runtime.RuntimeEvidence.GlobalRuntimeSecurityCoveragePercent
	if pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePushPercent < coverage {
		coverage = pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePushPercent
	}
	if pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePRPercent < coverage {
		coverage = pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePRPercent
	}
	return EventInspectionE2EAuditSummary{
		Capability: "CAP-INV-004", ContractID: contract.ContractID,
		AdversarialTests: adversarial.Count, RuntimeCoveragePercent: coverage,
		SASTFindings: runtime.StaticSecurity.Findings + pivot.StaticSecurity.Findings,
		SCAActionableFindings: runtime.StaticSecurity.ActionableFindings + pivot.StaticSecurity.ActionableFindings,
		ProjectionP95MS: projectionObserved, ProjectionBudgetMS: projectionBudget,
		PivotP95MS: pivotObserved, PivotBudgetMS: pivotBudget,
		RuntimeUnchangedSincePivot: true, BlockingOpenDecisions: len(handoff.BlockingOpenDecisions),
		Limitations: len(handoff.Limitations), ProductionReadinessClaim: handoff.ProductionReadinessClaim,
		FullCapabilityClaim: handoff.FullCapabilityCompletionClaim, Status: "PASS",
	}, nil
}

func validateEventInspectionProgress(runtime EventInspectionRuntimeProgress, pivot EventInspectionPivotProgress) error {
	if err := validateEventInspectionProgressIdentity(runtime.WorkUnit, runtime.Status, runtime.FinalValidatedHead, runtime.ProductSpecMutated, "E10-INV-004B-RUNTIME"); err != nil {
		return err
	}
	if err := validateEventInspectionProgressIdentity(pivot.WorkUnit, pivot.Status, pivot.FinalValidatedHead, pivot.ProductSpecMutated, "E10-INV-004C-PIVOT"); err != nil {
		return err
	}
	if runtime.RuntimeEvidence.Contract != "EVENT-INSPECTION-READONLY-CONTRACT-V1" ||
		runtime.RuntimeEvidence.RuntimeState != "implemented" ||
		runtime.RuntimeEvidence.ExternalRuntimeDependencies != 0 ||
		runtime.RuntimeEvidence.CrossBoundaryEdges != 0 ||
		runtime.RuntimeEvidence.Fixtures != 17 || runtime.RuntimeEvidence.PositiveFixtures != 4 ||
		runtime.RuntimeEvidence.NegativeFixtures != 13 ||
		runtime.RuntimeEvidence.GlobalRuntimeSecurityCoveragePercent < runtime.RuntimeEvidence.ChangedSecurityCriticalFloorPercent ||
		runtime.StaticSecurity.Findings != 0 || runtime.StaticSecurity.ActionableFindings != 0 {
		return fmt.Errorf("Event Inspection runtime progress does not satisfy bounded security evidence")
	}
	if pivot.RuntimeEvidence.Contract != "EVENT-INSPECTION-READONLY-CONTRACT-V1" ||
		pivot.RuntimeEvidence.ExternalRuntimeDependencies != 0 ||
		pivot.RuntimeEvidence.Fixtures != 17 || pivot.RuntimeEvidence.PositiveFixtures != 4 ||
		pivot.RuntimeEvidence.NegativeFixtures != 13 ||
		pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePushPercent < pivot.RuntimeEvidence.ChangedSecurityCriticalFloorPercent ||
		pivot.RuntimeEvidence.GlobalRuntimeSecurityCoveragePRPercent < pivot.RuntimeEvidence.ChangedSecurityCriticalFloorPercent ||
		pivot.StaticSecurity.Findings != 0 || pivot.StaticSecurity.ActionableFindings != 0 {
		return fmt.Errorf("Event Inspection pivot progress does not satisfy bounded security evidence")
	}
	return nil
}

func validateEventInspectionProgressIdentity(workUnit, status, head string, productSpecMutated bool, expected string) error {
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

func validateEventInspectionHandoff(handoff EventInspectionClosureHandoff) error {
	if handoff.SchemaVersion != 1 || handoff.HandoffKind != "event-inspection-bounded-core" ||
		handoff.Readiness != "BOUNDED_EVENT_INSPECTION_VALIDATED_NOT_FULL_CAPABILITY" ||
		handoff.Capability != "CAP-INV-004" {
		return fmt.Errorf("invalid Event Inspection bounded handoff identity/readiness")
	}
	if handoff.ProductionReadinessClaim || handoff.FullCapabilityCompletionClaim ||
		handoff.Class2MutationEnabled || handoff.ProviderSelected ||
		handoff.FinalQueryDialectSelected || handoff.DedicatedFinalUIClaim {
		return fmt.Errorf("Event Inspection bounded handoff overclaims scope or readiness")
	}
	if !sameStringSet(handoff.BlockingOpenDecisions, []string{"OPEN-013", "OPEN-014"}) {
		return fmt.Errorf("Event Inspection handoff must preserve OPEN-013 and OPEN-014")
	}
	if !sameStringSet(handoff.EvidenceSources, []string{eventInspectionRuntimeProgressPath, eventInspectionPivotProgressPath}) {
		return fmt.Errorf("Event Inspection handoff evidence sources mismatch")
	}
	for _, fragment := range eventInspectionRequiredLimitations {
		if !containsFoldFragment(handoff.Limitations, fragment) {
			return fmt.Errorf("Event Inspection handoff omits limitation containing %q", fragment)
		}
	}
	if len(handoff.NextConditions) < 4 {
		return fmt.Errorf("Event Inspection handoff requires at least four explicit next conditions")
	}
	return nil
}

func validateEventInspectionRuntimeUnchanged(root, evidenceHead, currentHead string) error {
	evidenceHead, err := validateFullCommitID(evidenceHead)
	if err != nil {
		return fmt.Errorf("Event Inspection pivot evidence head: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return fmt.Errorf("current source head: %w", err)
	}
	if err := ensureGitAncestor(root, evidenceHead, currentHead); err != nil {
		return fmt.Errorf("Event Inspection pivot evidence ancestry: %w", err)
	}
	cmd := exec.Command("git", "-C", root, "diff", "--quiet", evidenceHead, currentHead, "--", eventInspectionRuntimeDirectory) // #nosec G204,G702 -- executable/path fixed; commit arguments validated full hex IDs.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("Event Inspection runtime changed after verified pivot evidence; fresh runtime proof is required")
		}
		return fmt.Errorf("compare Event Inspection runtime against pivot evidence: %w", err)
	}
	return nil
}

func runEventInspectionAdversarialRuntimeTests(root string) (EventInspectionAdversarialSummary, error) {
	dir, err := resolveRepoPath(root, eventInspectionRuntimeDirectory, false)
	if err != nil {
		return EventInspectionAdversarialSummary{}, err
	}
	pattern := "^(" + strings.Join(eventInspectionAdversarialTests, "|") + ")$"
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", pattern, ".") // #nosec G204 -- executable fixed; regex assembled from compile-time test names.
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return EventInspectionAdversarialSummary{}, fmt.Errorf("Event Inspection adversarial runtime tests failed: %w", err)
	}
	return parseEventInspectionAdversarialTestEvents(stdout.Bytes())
}

func parseEventInspectionAdversarialTestEvents(data []byte) (EventInspectionAdversarialSummary, error) {
	type testEvent struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
	}
	passed := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return EventInspectionAdversarialSummary{}, fmt.Errorf("decode Event Inspection go test event: %w", err)
		}
		if event.Action == "pass" && containsString(eventInspectionAdversarialTests, event.Test) {
			passed[event.Test] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return EventInspectionAdversarialSummary{}, err
	}
	var values []string
	for test := range passed {
		values = append(values, test)
	}
	sort.Strings(values)
	expected := append([]string(nil), eventInspectionAdversarialTests...)
	sort.Strings(expected)
	summary := EventInspectionAdversarialSummary{Expected: expected, Passed: values, Count: len(values)}
	if !sameStringSet(values, expected) {
		return summary, fmt.Errorf("Event Inspection adversarial runtime tests incomplete: passed=%v expected=%v", values, expected)
	}
	return summary, nil
}

func eventInspectionPerformanceMetric(summary PerformanceBenchmarkAuditSummary, targetID, metricID string) (float64, float64, error) {
	if summary.Status != "pass" {
		return 0, 0, fmt.Errorf("Event Inspection performance benchmark status is %q", summary.Status)
	}
	for _, result := range summary.Results {
		if result.TargetID != targetID {
			continue
		}
		for _, metric := range result.Metrics {
			if metric.ID == metricID {
				if !metric.AbsolutePass || !metric.RelativePass || metric.Budget <= 0 || metric.Observed > metric.Budget {
					return metric.Observed, metric.Budget, fmt.Errorf("Event Inspection performance metric %s failed budget", metricID)
				}
				return metric.Observed, metric.Budget, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("Event Inspection performance result missing %s/%s", targetID, metricID)
}
