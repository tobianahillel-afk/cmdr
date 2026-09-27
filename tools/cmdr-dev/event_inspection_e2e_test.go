package main

import (
	"strings"
	"testing"
)

func validEventInspectionHandoff() EventInspectionClosureHandoff {
	return EventInspectionClosureHandoff{
		SchemaVersion:         1,
		HandoffKind:           "event-inspection-bounded-core",
		Readiness:             "BOUNDED_EVENT_INSPECTION_VALIDATED_NOT_FULL_CAPABILITY",
		Capability:            "CAP-INV-004",
		BlockingOpenDecisions: []string{"OPEN-013", "OPEN-014"},
		EvidenceSources:       []string{eventInspectionRuntimeProgressPath, eventInspectionPivotProgressPath},
		Limitations: []string{
			"Case-link mutation remains excluded while OPEN-013 is unresolved.",
			"Artifact proposal and Evidence-candidate mutations remain excluded while OPEN-013 and OPEN-014 are unresolved.",
			"No parser engine or entity-resolution engine is selected.",
			"No storage engine or provider is selected.",
			"No final query dialect is selected.",
			"No production readiness or production SLO is claimed.",
		},
		NextConditions: []string{"one", "two", "three", "four"},
	}
}

func TestValidateEventInspectionHandoff(t *testing.T) {
	if err := validateEventInspectionHandoff(validEventInspectionHandoff()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateEventInspectionHandoffRejectsOverclaim(t *testing.T) {
	handoff := validEventInspectionHandoff()
	handoff.ProductionReadinessClaim = true
	if err := validateEventInspectionHandoff(handoff); err == nil {
		t.Fatal("expected production readiness overclaim rejection")
	}
	handoff = validEventInspectionHandoff()
	handoff.BlockingOpenDecisions = []string{"OPEN-013"}
	if err := validateEventInspectionHandoff(handoff); err == nil {
		t.Fatal("expected missing OPEN-014 rejection")
	}
}

func TestParseEventInspectionAdversarialEvents(t *testing.T) {
	var lines []string
	for _, name := range eventInspectionAdversarialTests {
		lines = append(lines, `{"Action":"pass","Test":"`+name+`"}`)
	}
	summary, err := parseEventInspectionAdversarialTestEvents([]byte(strings.Join(lines, "\n") + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != len(eventInspectionAdversarialTests) {
		t.Fatalf("unexpected adversarial count: %#v", summary)
	}
}

func TestParseEventInspectionAdversarialEventsFailsClosedWhenMissing(t *testing.T) {
	if _, err := parseEventInspectionAdversarialTestEvents([]byte(`{"Action":"pass","Test":"TestRawDeniedNeverLeaksOrSynthesizesRawContent"}` + "\n")); err == nil {
		t.Fatal("expected incomplete Event Inspection adversarial evidence rejection")
	}
}

func TestEventInspectionPerformanceMetric(t *testing.T) {
	summary := PerformanceBenchmarkAuditSummary{Status: "pass", Results: []PerformanceBenchmarkResult{{
		TargetID: "target",
		Metrics:  []BenchmarkMetricResult{{ID: "metric", Observed: 0.5, Budget: 1, AbsolutePass: true, RelativePass: true}},
	}}}
	observed, budget, err := eventInspectionPerformanceMetric(summary, "target", "metric")
	if err != nil {
		t.Fatal(err)
	}
	if observed != 0.5 || budget != 1 {
		t.Fatalf("unexpected metric: %v/%v", observed, budget)
	}
	summary.Results[0].Metrics[0].AbsolutePass = false
	if _, _, err := eventInspectionPerformanceMetric(summary, "target", "metric"); err == nil {
		t.Fatal("expected failed performance metric rejection")
	}
}
