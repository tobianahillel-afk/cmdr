package main

import (
	"os"
	"path/filepath"
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

func TestEventInspectionProgressDecodersAcceptVerifiedEvidenceShape(t *testing.T) {
	root := t.TempDir()
	runtimePath := filepath.Join(root, filepath.FromSlash(eventInspectionRuntimeProgressPath))
	pivotPath := filepath.Join(root, filepath.FromSlash(eventInspectionPivotProgressPath))
	if err := os.MkdirAll(filepath.Dir(runtimePath), 0o750); err != nil {
		t.Fatal(err)
	}
	runtimeJSON := `{
		"schema_version":1,
		"work_unit":"E10-INV-004B-RUNTIME",
		"status":"VERIFIED",
		"final_validated_head":"5f6300b7a20d04b52c7c14b23fc31f90c55e9827",
		"validation":{"implementation_commit":"dfa2fe9786b375fc625dbf0a86f0913bf8569a54","adapter_test_fix_commit":"4edec0917b418ac569f60fef0c0638aeed521bee","security_scope_fix_commit":"5f6300b7a20d04b52c7c14b23fc31f90c55e9827","push_workflow_run":1,"pull_request_workflow_run":2,"result":"PASS"},
		"runtime_evidence":{"contract":"EVENT-INSPECTION-READONLY-CONTRACT-V1","runtime_state":"implemented","fixtures":17,"positive_fixtures":4,"negative_fixtures":13,"external_runtime_dependencies":0,"cross_boundary_edges":0,"global_runtime_security_coverage_percent":91.47,"global_coverage_floor_percent":80,"changed_security_critical_floor_percent":90,"pr_changed_security_scopes":4,"authorization_required_scopes":4,"tenant_isolation_required_scopes":4},
		"static_security":{"gosec_version":"v2.28.0","scan_roots":["product-runtime/event-inspection"],"findings":0,"govulncheck_version":"v1.8.0","analyzed_modules":8,"informational_findings":0,"actionable_findings":0},
		"invariants":["bounded"],
		"product_spec_mutated":false,
		"next_unlocked":["E10-INV-004C-PIVOT"]
	}`
	pivotJSON := `{
		"schema_version":1,
		"work_unit":"E10-INV-004C-PIVOT",
		"status":"VERIFIED",
		"final_validated_head":"ca95fd13ddf20d2a8937b2ed06929e5de6b4c40c",
		"validation":{"implementation_commit":"13d75b415f1b4f46942766f713f4d9cae0315a80","branch_coverage_fix_commit":"ca95fd13ddf20d2a8937b2ed06929e5de6b4c40c","push_workflow_run":3,"pull_request_workflow_run":4,"result":"PASS"},
		"runtime_evidence":{"contract":"EVENT-INSPECTION-READONLY-CONTRACT-V1","fixtures":17,"positive_fixtures":4,"negative_fixtures":13,"external_runtime_dependencies":0,"global_runtime_security_coverage_push_percent":91.61,"global_runtime_security_coverage_pr_percent":91.66,"global_coverage_floor_percent":80,"changed_security_critical_floor_percent":90},
		"static_security":{"gosec_version":"v2.28.0","findings":0,"govulncheck_version":"v1.8.0","analyzed_modules":8,"actionable_findings":0},
		"invariants":["bounded"],
		"product_spec_mutated":false,
		"next_unlocked":["E10-INV-004D-CLOSURE"]
	}`
	if err := os.WriteFile(runtimePath, []byte(runtimeJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(pivotPath), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pivotPath, []byte(pivotJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	var runtime EventInspectionRuntimeProgress
	if err := decodeStrict(root, eventInspectionRuntimeProgressPath, &runtime); err != nil {
		t.Fatal(err)
	}
	var pivot EventInspectionPivotProgress
	if err := decodeStrict(root, eventInspectionPivotProgressPath, &pivot); err != nil {
		t.Fatal(err)
	}
	if runtime.Validation.Result != "PASS" || pivot.Validation.Result != "PASS" {
		t.Fatalf("unexpected validation evidence: runtime=%#v pivot=%#v", runtime.Validation, pivot.Validation)
	}
}
