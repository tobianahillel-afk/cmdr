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
	pilotClosureHandoffPath      = "engineering/pilot/e2e-handoff.json"
	pilotRuntimeProgressPath     = "work/lots/E9-PILOT-001C-RUNTIME/PROGRESS.json"
	pilotSecurityProgressPath    = "work/lots/E9-PILOT-001C-SECURITY/PROGRESS.json"
	pilotPerformanceProgressPath = "work/lots/E9-PILOT-001C-PERF/PROGRESS.json"
	pilotClosureWorkUnit         = "E9-PILOT-001D"
)

var pilotAdversarialTests = []string{
	"TestProjectRejectsInvalidTenantReferences",
	"TestProjectRejectsUnsafeReturnOrigins",
	"TestProjectUnknownAuthorizationFailsClosed",
}

type PilotClosureHandoff struct {
	SchemaVersion                   int      `json:"schema_version"`
	HandoffKind                     string   `json:"handoff_kind"`
	Readiness                       string   `json:"readiness"`
	SliceID                         string   `json:"slice_id"`
	Capability                      string   `json:"capability"`
	ProductionReadinessClaim        bool     `json:"production_readiness_claim"`
	CanonicalPilotPlanStatus        string   `json:"canonical_pilot_plan_status"`
	SelectedOpenDecisions           []string `json:"selected_open_decisions"`
	GlobalOpenDecisionsOutsidePilot int      `json:"global_open_decisions_outside_pilot"`
	EvidenceSources                 []string `json:"evidence_sources"`
	Limitations                     []string `json:"limitations"`
	ScaleNextConditions             []string `json:"scale_next_conditions"`
}

type PilotRuntimeProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	SecurityValidation struct {
		GlobalCoveragePercent                  float64 `json:"global_coverage_percent"`
		ChangedSecurityCriticalFloorPercent    float64 `json:"changed_security_critical_floor_percent"`
		ChangedSecurityCriticalFloorPassedOnPR bool    `json:"changed_security_critical_floor_passed_on_pull_request"`
		AuthorizationNegativeGate              string  `json:"authorization_negative_gate"`
		TenantIsolationNegativeGate            string  `json:"tenant_isolation_negative_gate"`
		SASTFindings                           int     `json:"sast_findings"`
		SCAActionableFindings                  int     `json:"sca_actionable_findings"`
		RuntimeDependencies                    int     `json:"runtime_dependencies"`
		CrossBoundaryEdges                     int     `json:"cross_boundary_edges"`
	} `json:"security_validation"`
}

type PilotSecurityProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	Validation         struct {
		Result                                   string  `json:"result"`
		CoverageStatus                           string  `json:"coverage_status"`
		GlobalCoveragePercent                    float64 `json:"global_coverage_percent"`
		GlobalFloorPercent                       float64 `json:"global_floor_percent"`
		ChangedFloorPercent                      float64 `json:"changed_floor_percent"`
		PullRequestChangedSecurityCriticalScopes int     `json:"pull_request_changed_security_critical_scopes"`
		AuthorizationNegative                    string  `json:"authorization_negative"`
		TenantIsolationNegative                  string  `json:"tenant_isolation_negative"`
		DeepSecurity                             string  `json:"deep_security"`
	} `json:"validation"`
	Gates struct {
		AuthorizationNegative string `json:"SEC-AUTH-NEG-001"`
		TenantIsolation       string `json:"SEC-TENANT-ISO-001"`
	} `json:"gates"`
	StaticSecurity struct {
		GosecVersion          string   `json:"gosec_version"`
		ScanRoots             []string `json:"scan_roots"`
		Findings              int      `json:"findings"`
		GovulncheckVersion    string   `json:"govulncheck_version"`
		Modules               int      `json:"modules"`
		InformationalFindings int      `json:"informational_findings"`
		ActionableFindings    int      `json:"actionable_findings"`
	} `json:"static_security"`
}

type PilotPerformanceProgress struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Status             string `json:"status"`
	FinalValidatedHead string `json:"final_validated_head"`
	ProductSpecMutated bool   `json:"product_spec_mutated"`
	Validation         struct {
		Result                       string  `json:"result"`
		PerformanceTarget            string  `json:"performance_target"`
		Metric                       string  `json:"metric"`
		BudgetMS                     float64 `json:"budget_ms"`
		PushObservedP95MS            float64 `json:"push_observed_p95_ms"`
		PullRequestObservedP95MS     float64 `json:"pull_request_observed_p95_ms"`
		AbsoluteBudgetPass           bool    `json:"absolute_budget_pass"`
		RelativeBaselineClaimed      bool    `json:"relative_baseline_claimed"`
		RuntimeGlobalCoveragePercent float64 `json:"runtime_global_coverage_percent"`
		SASTFindings                 int     `json:"sast_findings"`
		SCAActionableFindings        int     `json:"sca_actionable_findings"`
		PerformanceCacheRecords      int     `json:"performance_cache_records"`
	} `json:"validation"`
	Implementation struct {
		BenchmarkHandler    string   `json:"benchmark_handler"`
		BenchmarkStages     []string `json:"benchmark_stages"`
		Probe               string   `json:"probe"`
		MeasurementScope    string   `json:"measurement_scope"`
		DeepHandler         string   `json:"deep_handler"`
		RuntimeDependencies int      `json:"runtime_dependencies"`
	} `json:"implementation"`
}

type PilotAdversarialTestSummary struct {
	Expected []string `json:"expected"`
	Passed   []string `json:"passed"`
	Count    int      `json:"count"`
}

type PilotE2EAuditSummary struct {
	SliceID                  string  `json:"slice_id"`
	Capability               string  `json:"capability"`
	ProductBaselineCommit    string  `json:"product_baseline_commit"`
	ProductTreeDigest        string  `json:"product_tree_digest"`
	Requirements             int     `json:"requirements"`
	Permissions              int     `json:"permissions"`
	ContractID               string  `json:"contract_id"`
	ContractFixtures         int     `json:"contract_fixtures"`
	NegativeFixtures         int     `json:"negative_fixtures"`
	AdversarialTests         int     `json:"adversarial_tests"`
	RuntimeCoveragePercent   float64 `json:"runtime_coverage_percent"`
	AuthorizationNegative    string  `json:"authorization_negative"`
	TenantIsolationNegative  string  `json:"tenant_isolation_negative"`
	SASTFindings             int     `json:"sast_findings"`
	SCAActionableFindings    int     `json:"sca_actionable_findings"`
	PerformanceBudgetMS      float64 `json:"performance_budget_ms"`
	WorstObservedP95MS       float64 `json:"worst_observed_p95_ms"`
	RecoveryOutcome          string  `json:"recovery_outcome"`
	RecoveryNextAction       string  `json:"recovery_next_action"`
	GlobalOpenDecisions      int     `json:"global_open_decisions_outside_pilot"`
	SelectedOpenDecisions    int     `json:"selected_open_decisions"`
	ProductionReadinessClaim bool    `json:"production_readiness_claim"`
	Limitations              int     `json:"limitations"`
	ScaleNextConditions      int     `json:"scale_next_conditions"`
	Status                   string  `json:"status"`
}

func runPilotE2EAudit(root string, state CurrentState, graph WorkGraph) (PilotE2EAuditSummary, error) {
	scopeSummary, err := runPilotScopeAudit(root, state)
	if err != nil {
		return PilotE2EAuditSummary{}, fmt.Errorf("pilot scope: %w", err)
	}
	contractSummary, err := runPilotContractAudit(root)
	if err != nil {
		return PilotE2EAuditSummary{}, fmt.Errorf("pilot contract: %w", err)
	}
	baseline, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, "engineering/spec-index/baseline.json", true)
	if err != nil {
		return PilotE2EAuditSummary{}, fmt.Errorf("product baseline: %w", err)
	}

	var scope PilotScope
	if err := decodeStrict(root, pilotScopePath, &scope); err != nil {
		return PilotE2EAuditSummary{}, err
	}
	var handoff PilotClosureHandoff
	if err := decodeStrict(root, pilotClosureHandoffPath, &handoff); err != nil {
		return PilotE2EAuditSummary{}, err
	}
	var runtime PilotRuntimeProgress
	if err := decodePilotEvidence(root, pilotRuntimeProgressPath, &runtime); err != nil {
		return PilotE2EAuditSummary{}, err
	}
	var security PilotSecurityProgress
	if err := decodePilotEvidence(root, pilotSecurityProgressPath, &security); err != nil {
		return PilotE2EAuditSummary{}, err
	}
	var performance PilotPerformanceProgress
	if err := decodePilotEvidence(root, pilotPerformanceProgressPath, &performance); err != nil {
		return PilotE2EAuditSummary{}, err
	}

	recovery, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{WorkUnit: pilotClosureWorkUnit}, state, graph)
	if err != nil {
		return PilotE2EAuditSummary{}, fmt.Errorf("recovery reconciliation: %w", err)
	}
	if recovery.Outcome == "conflict" {
		return PilotE2EAuditSummary{}, fmt.Errorf("recovery reconciliation conflict: %s", strings.Join(recovery.Reasons, ", "))
	}
	adversarial, err := runPilotAdversarialRuntimeTests(root)
	if err != nil {
		return PilotE2EAuditSummary{}, err
	}

	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err != nil {
		return PilotE2EAuditSummary{}, err
	}

	worstP95 := performance.Validation.PushObservedP95MS
	if performance.Validation.PullRequestObservedP95MS > worstP95 {
		worstP95 = performance.Validation.PullRequestObservedP95MS
	}
	return PilotE2EAuditSummary{
		SliceID:                  scopeSummary.SliceID,
		Capability:               scopeSummary.Capability,
		ProductBaselineCommit:    baseline.ProductSpecBaselineCommit,
		ProductTreeDigest:        baseline.TreeDigest,
		Requirements:             scopeSummary.Requirements,
		Permissions:              scopeSummary.Permissions,
		ContractID:               contractSummary.ContractID,
		ContractFixtures:         contractSummary.Fixtures,
		NegativeFixtures:         contractSummary.NegativeFixtures,
		AdversarialTests:         adversarial.Count,
		RuntimeCoveragePercent:   security.Validation.GlobalCoveragePercent,
		AuthorizationNegative:    security.Validation.AuthorizationNegative,
		TenantIsolationNegative:  security.Validation.TenantIsolationNegative,
		SASTFindings:             security.StaticSecurity.Findings,
		SCAActionableFindings:    security.StaticSecurity.ActionableFindings,
		PerformanceBudgetMS:      performance.Validation.BudgetMS,
		WorstObservedP95MS:       worstP95,
		RecoveryOutcome:          recovery.Outcome,
		RecoveryNextAction:       recovery.NextAction,
		GlobalOpenDecisions:      handoff.GlobalOpenDecisionsOutsidePilot,
		SelectedOpenDecisions:    len(handoff.SelectedOpenDecisions),
		ProductionReadinessClaim: handoff.ProductionReadinessClaim,
		Limitations:              len(handoff.Limitations),
		ScaleNextConditions:      len(handoff.ScaleNextConditions),
		Status:                   "PASS",
	}, nil
}

func decodePilotEvidence(root, path string, dst any) error {
	data, err := readRepoFile(root, path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func validatePilotClosureEvidence(
	state CurrentState,
	scope PilotScope,
	handoff PilotClosureHandoff,
	runtime PilotRuntimeProgress,
	security PilotSecurityProgress,
	performance PilotPerformanceProgress,
	recovery RecoveryReconciliation,
	adversarial PilotAdversarialTestSummary,
	scopeSummary PilotScopeAuditSummary,
	contractSummary PilotContractAuditSummary,
) error {
	if handoff.SchemaVersion != 1 || handoff.HandoffKind != "e9-bounded-pilot-closure" {
		return fmt.Errorf("invalid pilot closure handoff header")
	}
	if handoff.Readiness != "BOUNDED_PILOT_VALIDATED_NOT_PRODUCTION_READY" {
		return fmt.Errorf("pilot handoff readiness must remain bounded and non-production")
	}
	if handoff.SliceID != scope.SliceID || handoff.Capability != scope.Capability ||
		handoff.SliceID != scopeSummary.SliceID || handoff.Capability != scopeSummary.Capability {
		return fmt.Errorf("pilot handoff identity does not match canonical scope")
	}
	if handoff.ProductionReadinessClaim {
		return fmt.Errorf("bounded pilot must not claim production readiness")
	}
	if handoff.CanonicalPilotPlanStatus != "draft" {
		return fmt.Errorf("canonical pilot plan status must remain explicit draft")
	}
	if len(handoff.SelectedOpenDecisions) != 0 || len(scope.OpenDecisions) != 0 {
		return fmt.Errorf("selected pilot must not silently consume unresolved open decisions")
	}
	if handoff.GlobalOpenDecisionsOutsidePilot != state.ProductSpec.KnownOpenDecisions {
		return fmt.Errorf("global open decision count mismatch: handoff=%d state=%d", handoff.GlobalOpenDecisionsOutsidePilot, state.ProductSpec.KnownOpenDecisions)
	}
	expectedSources := []string{pilotRuntimeProgressPath, pilotSecurityProgressPath, pilotPerformanceProgressPath}
	if !sameStringSet(handoff.EvidenceSources, expectedSources) {
		return fmt.Errorf("pilot handoff evidence sources mismatch: got=%v want=%v", handoff.EvidenceSources, expectedSources)
	}
	for _, excluded := range scope.ExcludedBehaviors {
		if !containsString(handoff.Limitations, excluded) {
			return fmt.Errorf("pilot handoff omitted scope limitation %q", excluded)
		}
	}
	if len(handoff.ScaleNextConditions) < 3 {
		return fmt.Errorf("pilot handoff must state at least three scale-next conditions")
	}

	if err := validateVerifiedPilotProgress(runtime.WorkUnit, runtime.Status, runtime.FinalValidatedHead, runtime.ProductSpecMutated, "E9-PILOT-001C-RUNTIME"); err != nil {
		return err
	}
	if runtime.SecurityValidation.GlobalCoveragePercent < 80 ||
		runtime.SecurityValidation.ChangedSecurityCriticalFloorPercent < 90 ||
		!runtime.SecurityValidation.ChangedSecurityCriticalFloorPassedOnPR ||
		runtime.SecurityValidation.AuthorizationNegativeGate != "PASS" ||
		runtime.SecurityValidation.TenantIsolationNegativeGate != "PASS" ||
		runtime.SecurityValidation.SASTFindings != 0 ||
		runtime.SecurityValidation.SCAActionableFindings != 0 ||
		runtime.SecurityValidation.RuntimeDependencies != 0 ||
		runtime.SecurityValidation.CrossBoundaryEdges != 0 {
		return fmt.Errorf("runtime progress does not satisfy pilot security/correctness floors")
	}

	if err := validateVerifiedPilotProgress(security.WorkUnit, security.Status, security.FinalValidatedHead, security.ProductSpecMutated, "E9-PILOT-001C-SECURITY"); err != nil {
		return err
	}
	if security.Validation.Result != "PASS" || security.Validation.CoverageStatus != "measured" ||
		security.Validation.GlobalCoveragePercent < security.Validation.GlobalFloorPercent ||
		security.Validation.GlobalFloorPercent < 80 ||
		security.Validation.ChangedFloorPercent < 90 ||
		security.Validation.PullRequestChangedSecurityCriticalScopes < 1 ||
		security.Validation.AuthorizationNegative != "PASS" ||
		security.Validation.TenantIsolationNegative != "PASS" ||
		security.Validation.DeepSecurity != "PASS" ||
		!strings.HasPrefix(security.Gates.AuthorizationNegative, "active:") ||
		!strings.HasPrefix(security.Gates.TenantIsolation, "active:") ||
		security.StaticSecurity.Findings != 0 ||
		security.StaticSecurity.ActionableFindings != 0 {
		return fmt.Errorf("security progress does not satisfy measured pilot gates")
	}

	if err := validateVerifiedPilotProgress(performance.WorkUnit, performance.Status, performance.FinalValidatedHead, performance.ProductSpecMutated, "E9-PILOT-001C-PERF"); err != nil {
		return err
	}
	if performance.Validation.Result != "PASS" ||
		performance.Validation.PerformanceTarget != "PERF-TGT-PILOT-CONTEXT-PROJECTION" ||
		performance.Validation.Metric != "projection-p95-latency" ||
		performance.Validation.BudgetMS <= 0 ||
		!performance.Validation.AbsoluteBudgetPass ||
		performance.Validation.RelativeBaselineClaimed ||
		performance.Validation.PushObservedP95MS > performance.Validation.BudgetMS ||
		performance.Validation.PullRequestObservedP95MS > performance.Validation.BudgetMS ||
		performance.Validation.RuntimeGlobalCoveragePercent < 80 ||
		performance.Validation.SASTFindings != 0 ||
		performance.Validation.SCAActionableFindings != 0 ||
		performance.Implementation.RuntimeDependencies != 0 ||
		strings.TrimSpace(performance.Implementation.BenchmarkHandler) == "" ||
		strings.TrimSpace(performance.Implementation.Probe) == "" {
		return fmt.Errorf("performance progress does not satisfy bounded pilot budget")
	}

	if contractSummary.Status != "PASS" || scopeSummary.Status != "PASS" ||
		contractSummary.RuntimeDependencies != 0 || contractSummary.Fixtures < 1 || contractSummary.NegativeFixtures < 1 {
		return fmt.Errorf("pilot scope/contract audit evidence is incomplete")
	}
	if recovery.Outcome == "conflict" || strings.TrimSpace(recovery.NextAction) == "" {
		return fmt.Errorf("recovery evidence is not safely resumable")
	}
	if recovery.Checkpoint == "empty" && (recovery.Outcome != "resume" || recovery.NextAction != "claim-or-start") {
		return fmt.Errorf("empty recovery checkpoint must resume via claim-or-start")
	}
	if adversarial.Count != len(pilotAdversarialTests) || !sameStringSet(adversarial.Passed, pilotAdversarialTests) {
		return fmt.Errorf("adversarial runtime test evidence incomplete: got=%v", adversarial.Passed)
	}
	return nil
}

func validateVerifiedPilotProgress(workUnit, status, head string, productSpecMutated bool, expected string) error {
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

func runPilotAdversarialRuntimeTests(root string) (PilotAdversarialTestSummary, error) {
	pattern := "^(" + strings.Join(pilotAdversarialTests, "|") + ")$"
	cmd := exec.Command("go", "test", "-json", "-count=1", "-run", pattern, ".")
	cmd.Dir = filepath.Join(root, "product-runtime", "context-envelope")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return PilotAdversarialTestSummary{}, fmt.Errorf("pilot adversarial runtime tests failed: %w", err)
	}
	summary, err := parsePilotAdversarialTestEvents(stdout.Bytes())
	if err != nil {
		return summary, err
	}
	return summary, nil
}

func parsePilotAdversarialTestEvents(data []byte) (PilotAdversarialTestSummary, error) {
	type testEvent struct {
		Action string `json:"Action"`
		Test   string `json:"Test"`
	}
	passed := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return PilotAdversarialTestSummary{}, fmt.Errorf("decode go test event: %w", err)
		}
		if event.Action == "pass" && containsString(pilotAdversarialTests, event.Test) {
			passed[event.Test] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return PilotAdversarialTestSummary{}, err
	}
	var values []string
	for test := range passed {
		values = append(values, test)
	}
	sort.Strings(values)
	expected := append([]string(nil), pilotAdversarialTests...)
	sort.Strings(expected)
	summary := PilotAdversarialTestSummary{Expected: expected, Passed: values, Count: len(values)}
	if !sameStringSet(values, expected) {
		return summary, fmt.Errorf("pilot adversarial runtime tests incomplete: passed=%v expected=%v", values, expected)
	}
	return summary, nil
}
