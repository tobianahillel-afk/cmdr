package main

import (
	"strings"
	"testing"
)

func validHuntManagementHandoff() HuntManagementClosureHandoff {
	return HuntManagementClosureHandoff{
		SchemaVersion: 1,
		HandoffKind: "hunt-management-bounded-core",
		Readiness: "BOUNDED_HUNT_MANAGEMENT_VALIDATED_NOT_FULL_CAPABILITY",
		Capability: "CAP-INV-005",
		BlockingOpenDecisions: []string{"OPEN-013"},
		EvidenceSources: []string{huntManagementRuntimeProgressPath},
		Limitations: []string{
			"Hunt lifecycle and all class-2 Hunt mutations remain excluded while OPEN-013 is unresolved.",
			"Case promotion and Case-link mutation remain excluded.",
			"Hypothesis creation, linking and status mutation remain excluded.",
			"Query and Saved Search mutation remain excluded.",
			"Search Job execution remains excluded.",
			"No storage engine, provider or retention policy is selected.",
			"No collaboration backend is selected.",
			"No final query dialect is selected.",
			"No dedicated final Hunt UI is claimed.",
			"No production readiness or production Hunt SLO is claimed.",
			"No full-capability completion is claimed.",
		},
		NextConditions: []string{"one", "two", "three", "four"},
	}
}

func TestValidateHuntManagementHandoff(t *testing.T) {
	if err := validateHuntManagementHandoff(validHuntManagementHandoff()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateHuntManagementHandoffRejectsOverclaim(t *testing.T) {
	handoff := validHuntManagementHandoff()
	handoff.ProductionReadinessClaim = true
	if err := validateHuntManagementHandoff(handoff); err == nil {
		t.Fatal("expected production readiness overclaim rejection")
	}
	handoff = validHuntManagementHandoff()
	handoff.Class2MutationEnabled = true
	if err := validateHuntManagementHandoff(handoff); err == nil {
		t.Fatal("expected class-2 mutation overclaim rejection")
	}
	handoff = validHuntManagementHandoff()
	handoff.BlockingOpenDecisions = nil
	if err := validateHuntManagementHandoff(handoff); err == nil {
		t.Fatal("expected missing OPEN-013 rejection")
	}
}

func TestParseHuntManagementAdversarialEvents(t *testing.T) {
	var lines []string
	for _, name := range huntManagementAdversarialTests {
		lines = append(lines, `{"Action":"pass","Test":"`+name+`"}`)
	}
	summary, err := parseHuntManagementAdversarialTestEvents([]byte(strings.Join(lines, "\n") + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != len(huntManagementAdversarialTests) {
		t.Fatalf("unexpected adversarial count: %#v", summary)
	}
}

func TestParseHuntManagementAdversarialEventsFailsClosedWhenMissing(t *testing.T) {
	if _, err := parseHuntManagementAdversarialTestEvents([]byte(`{"Action":"pass","Test":"TestProjectionDeepCopiesCallerOwnedWorkspace"}` + "\n")); err == nil {
		t.Fatal("expected incomplete Hunt Management adversarial evidence rejection")
	}
}

func TestHuntManagementPerformanceMetric(t *testing.T) {
	summary := PerformanceBenchmarkAuditSummary{Status: "pass", Results: []PerformanceBenchmarkResult{{
		TargetID: "target",
		Metrics: []BenchmarkMetricResult{{ID: "metric", Observed: 0.5, Budget: 1, AbsolutePass: true, RelativePass: true}},
	}}}
	observed, budget, err := huntManagementPerformanceMetric(summary, "target", "metric")
	if err != nil {
		t.Fatal(err)
	}
	if observed != 0.5 || budget != 1 {
		t.Fatalf("unexpected metric: %v/%v", observed, budget)
	}
	summary.Results[0].Metrics[0].AbsolutePass = false
	if _, _, err := huntManagementPerformanceMetric(summary, "target", "metric"); err == nil {
		t.Fatal("expected failed performance metric rejection")
	}
}

func TestValidateHuntManagementRuntimeProgress(t *testing.T) {
	progress := HuntManagementRuntimeProgress{
		SchemaVersion: 1,
		WorkUnit: "E10-INV-005B-RUNTIME",
		Status: "VERIFIED",
		FinalValidatedHead: "65df48f8466367438fbf3bbb933025869317fa88",
	}
	progress.Validation.Result = "PASS"
	progress.RuntimeObservation.Contract = "HUNT-WORKSPACE-READONLY-CONTRACT-V1"
	progress.RuntimeObservation.RuntimeState = "implemented"
	progress.RuntimeObservation.Fixtures = 15
	progress.RuntimeObservation.Positive = 3
	progress.RuntimeObservation.Negative = 12
	progress.RuntimeObservation.GlobalRuntimeCoverage = 91.84
	progress.RuntimeObservation.GlobalCoverageFloor = 80
	progress.RuntimeObservation.ChangedSecurityFloor = 90
	if err := validateHuntManagementRuntimeProgress(progress); err != nil {
		t.Fatal(err)
	}
	progress.RuntimeObservation.SASTFindings = 1
	if err := validateHuntManagementRuntimeProgress(progress); err == nil {
		t.Fatal("expected SAST evidence rejection")
	}
}
