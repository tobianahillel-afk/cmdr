package main

import (
	"strings"
	"testing"
)

func TestParsePilotAdversarialTestEventsRequiresEveryExpectedTest(t *testing.T) {
	var lines []string
	for _, name := range pilotAdversarialTests {
		lines = append(lines, `{"Action":"pass","Test":"`+name+`"}`)
	}
	summary, err := parsePilotAdversarialTestEvents([]byte(strings.Join(lines, "\n")))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != len(pilotAdversarialTests) {
		t.Fatalf("unexpected adversarial count: %#v", summary)
	}

	if _, err := parsePilotAdversarialTestEvents([]byte(lines[0])); err == nil {
		t.Fatal("expected incomplete adversarial evidence rejection")
	}
}

func TestValidatePilotClosureEvidenceRejectsProductionReadinessClaim(t *testing.T) {
	state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary := validPilotClosureEvidence()
	handoff.ProductionReadinessClaim = true
	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err == nil {
		t.Fatal("expected production-readiness claim rejection")
	}
}

func TestValidatePilotClosureEvidenceRejectsSecurityRegression(t *testing.T) {
	state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary := validPilotClosureEvidence()
	security.StaticSecurity.Findings = 1
	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err == nil {
		t.Fatal("expected SAST regression rejection")
	}
}

func TestValidatePilotClosureEvidenceRejectsPerformanceRegression(t *testing.T) {
	state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary := validPilotClosureEvidence()
	performance.Validation.PullRequestObservedP95MS = 2
	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err == nil {
		t.Fatal("expected performance regression rejection")
	}
}

func TestValidatePilotClosureEvidenceRejectsMissingScopeLimitation(t *testing.T) {
	state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary := validPilotClosureEvidence()
	handoff.Limitations = handoff.Limitations[1:]
	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err == nil {
		t.Fatal("expected missing limitation rejection")
	}
}

func TestValidatePilotClosureEvidenceRejectsRecoveryConflict(t *testing.T) {
	state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary := validPilotClosureEvidence()
	recovery.Outcome = "conflict"
	if err := validatePilotClosureEvidence(state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary); err == nil {
		t.Fatal("expected recovery conflict rejection")
	}
}

func validPilotClosureEvidence() (
	CurrentState,
	PilotScope,
	PilotClosureHandoff,
	PilotRuntimeProgress,
	PilotSecurityProgress,
	PilotPerformanceProgress,
	RecoveryReconciliation,
	PilotAdversarialTestSummary,
	PilotScopeAuditSummary,
	PilotContractAuditSummary,
) {
	state := CurrentState{}
	state.ProductSpec.KnownOpenDecisions = 17
	scope := PilotScope{
		SliceID: "PILOT-CONTEXT-001", Capability: "CAP-SET-004",
		ExcludedBehaviors: []string{"excluded-one", "excluded-two"},
	}
	handoff := PilotClosureHandoff{
		SchemaVersion: 1, HandoffKind: "e9-bounded-pilot-closure",
		Readiness: "BOUNDED_PILOT_VALIDATED_NOT_PRODUCTION_READY",
		SliceID: scope.SliceID, Capability: scope.Capability,
		CanonicalPilotPlanStatus: "draft", GlobalOpenDecisionsOutsidePilot: 17,
		EvidenceSources: []string{pilotRuntimeProgressPath, pilotSecurityProgressPath, pilotPerformanceProgressPath},
		Limitations: append([]string(nil), scope.ExcludedBehaviors...),
		ScaleNextConditions: []string{"one", "two", "three"},
	}
	const head = "0123456789abcdef0123456789abcdef01234567"
	runtime := PilotRuntimeProgress{SchemaVersion: 1, WorkUnit: "E9-PILOT-001C-RUNTIME", Status: "VERIFIED", FinalValidatedHead: head}
	runtime.SecurityValidation.GlobalCoveragePercent = 95
	runtime.SecurityValidation.ChangedSecurityCriticalFloorPercent = 90
	runtime.SecurityValidation.ChangedSecurityCriticalFloorPassedOnPR = true
	runtime.SecurityValidation.AuthorizationNegativeGate = "PASS"
	runtime.SecurityValidation.TenantIsolationNegativeGate = "PASS"

	security := PilotSecurityProgress{SchemaVersion: 1, WorkUnit: "E9-PILOT-001C-SECURITY", Status: "VERIFIED", FinalValidatedHead: head}
	security.Validation.Result = "PASS"
	security.Validation.CoverageStatus = "measured"
	security.Validation.GlobalCoveragePercent = 95
	security.Validation.GlobalFloorPercent = 80
	security.Validation.ChangedFloorPercent = 90
	security.Validation.PullRequestChangedSecurityCriticalScopes = 1
	security.Validation.AuthorizationNegative = "PASS"
	security.Validation.TenantIsolationNegative = "PASS"
	security.Validation.DeepSecurity = "PASS"
	security.Gates.AuthorizationNegative = "active:test"
	security.Gates.TenantIsolation = "active:test"

	performance := PilotPerformanceProgress{SchemaVersion: 1, WorkUnit: "E9-PILOT-001C-PERF", Status: "VERIFIED", FinalValidatedHead: head}
	performance.Validation.Result = "PASS"
	performance.Validation.PerformanceTarget = "PERF-TGT-PILOT-CONTEXT-PROJECTION"
	performance.Validation.Metric = "projection-p95-latency"
	performance.Validation.BudgetMS = 1
	performance.Validation.PushObservedP95MS = 0.1
	performance.Validation.PullRequestObservedP95MS = 0.2
	performance.Validation.AbsoluteBudgetPass = true
	performance.Validation.RuntimeGlobalCoveragePercent = 90
	performance.Implementation.BenchmarkHandler = "builtin-test"
	performance.Implementation.Probe = "probe"

	recovery := RecoveryReconciliation{Checkpoint: "empty", Outcome: "resume", NextAction: "claim-or-start"}
	adversarial := PilotAdversarialTestSummary{Expected: append([]string(nil), pilotAdversarialTests...), Passed: append([]string(nil), pilotAdversarialTests...), Count: len(pilotAdversarialTests)}
	scopeSummary := PilotScopeAuditSummary{SliceID: scope.SliceID, Capability: scope.Capability, Status: "PASS"}
	contractSummary := PilotContractAuditSummary{ContractID: "PILOT-CONTEXT-ENVELOPE-V1", Fixtures: 11, NegativeFixtures: 9, RuntimeDependencies: 0, Status: "PASS"}
	return state, scope, handoff, runtime, security, performance, recovery, adversarial, scopeSummary, contractSummary
}
