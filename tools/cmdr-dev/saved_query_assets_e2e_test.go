package main

import (
	"strings"
	"testing"
)

func validSavedQueryAssetsClosureHandoff() SavedQueryAssetsClosureHandoff {
	return SavedQueryAssetsClosureHandoff{
		SchemaVersion:         1,
		HandoffKind:           "saved-query-assets-bounded-core",
		Readiness:             "BOUNDED_SAVED_QUERY_ASSETS_VALIDATED_NOT_FULL_CAPABILITY",
		Capability:            "CAP-INV-006",
		BlockingOpenDecisions: []string{"OPEN-013"},
		EvidenceSources:       []string{savedQueryAssetsRuntimeProgressPath, savedQueryAssetsHandoffProgressPath},
		Limitations: []string{
			"All class-2 mutation remains excluded while OPEN-013 is unresolved.",
			"No canonical Saved Search or Query Asset persistence schema is selected.",
			"Saved View remains separate.",
			"Detection Rule lifecycle remains outside scope.",
			"Search Job creation and execution remain outside scope.",
			"No storage engine is selected.",
			"No provider is selected.",
			"No retention policy is selected.",
			"No collaboration backend is selected.",
			"No final query dialect is selected.",
			"No dedicated final Saved Search or Query Asset UI is claimed.",
			"No production readiness is claimed.",
			"No full-capability completion is claimed.",
		},
		NextConditions: []string{"one", "two", "three", "four", "five"},
	}
}

func TestValidateSavedQueryAssetsClosureHandoff(t *testing.T) {
	if err := validateSavedQueryAssetsClosureHandoff(validSavedQueryAssetsClosureHandoff()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSavedQueryAssetsClosureHandoffRejectsOverclaim(t *testing.T) {
	h := validSavedQueryAssetsClosureHandoff()
	h.ProductionReadinessClaim = true
	if err := validateSavedQueryAssetsClosureHandoff(h); err == nil {
		t.Fatal("expected production readiness overclaim rejection")
	}
	h = validSavedQueryAssetsClosureHandoff()
	h.Class2MutationEnabled = true
	if err := validateSavedQueryAssetsClosureHandoff(h); err == nil {
		t.Fatal("expected class-2 mutation overclaim rejection")
	}
	h = validSavedQueryAssetsClosureHandoff()
	h.BlockingOpenDecisions = nil
	if err := validateSavedQueryAssetsClosureHandoff(h); err == nil {
		t.Fatal("expected missing OPEN-013 rejection")
	}
}

func TestParseSavedQueryAssetsAdversarialEvents(t *testing.T) {
	var lines []string
	for _, name := range savedQueryAssetsAdversarialTests {
		lines = append(lines, `{"Action":"pass","Test":"`+name+`"}`)
	}
	summary, err := parseSavedQueryAssetsAdversarialTestEvents([]byte(strings.Join(lines, "\n") + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != len(savedQueryAssetsAdversarialTests) {
		t.Fatalf("unexpected adversarial count: %#v", summary)
	}
}

func TestParseSavedQueryAssetsAdversarialEventsFailsClosedWhenMissing(t *testing.T) {
	if _, err := parseSavedQueryAssetsAdversarialTestEvents([]byte(`{"Action":"pass","Test":"TestProjectionDeepCopiesCallerOwnedAsset"}` + "\n")); err == nil {
		t.Fatal("expected incomplete adversarial evidence rejection")
	}
}

func TestSavedQueryAssetsPerformanceMetric(t *testing.T) {
	summary := PerformanceBenchmarkAuditSummary{Status: "pass", Results: []PerformanceBenchmarkResult{{
		TargetID: "target",
		Metrics:  []BenchmarkMetricResult{{ID: "metric", Observed: 0.5, Budget: 1, AbsolutePass: true, RelativePass: true}},
	}}}
	observed, budget, err := savedQueryAssetsPerformanceMetric(summary, "target", "metric")
	if err != nil {
		t.Fatal(err)
	}
	if observed != 0.5 || budget != 1 {
		t.Fatalf("unexpected metric: %v/%v", observed, budget)
	}
	summary.Results[0].Metrics[0].AbsolutePass = false
	if _, _, err := savedQueryAssetsPerformanceMetric(summary, "target", "metric"); err == nil {
		t.Fatal("expected failed performance metric rejection")
	}
}

func TestValidateSavedQueryAssetsRuntimeProgress(t *testing.T) {
	progress := SavedQueryAssetsRuntimeProgress{
		SchemaVersion:      1,
		WorkUnit:           "E10-INV-006B-RUNTIME",
		Status:             "VERIFIED",
		FinalValidatedHead: "f262c1a8ff0d1ecc1ebe5f15c15dad25809f3334",
	}
	progress.Validation.PushWorkflowRun = 36488064659
	progress.Validation.PullRequestWorkflowRun = 36488071274
	progress.Validation.Result = "PASS"
	progress.Validation.SavedQueryAssets = "PASS"
	progress.Validation.SecurityTests = "PASS"
	progress.Validation.SAST = "PASS"
	progress.Validation.SCA = "PASS"
	progress.Validation.PerformanceRegistry = "PASS"
	progress.Validation.Architecture = "PASS"
	progress.Validation.RuntimeDependencies = "PASS"
	progress.Validation.BoundaryEdges = "PASS"
	progress.RuntimeObservation.Contract = "QUERY-ASSET-READONLY-CONTRACT-V1"
	progress.RuntimeObservation.RuntimeState = "implemented"
	progress.RuntimeObservation.Fixtures = 18
	progress.RuntimeObservation.Positive = 5
	progress.RuntimeObservation.Negative = 13
	progress.RuntimeObservation.Incompatible = 2
	progress.RuntimeObservation.Stale = 1
	progress.RuntimeObservation.PushGlobalCoveragePercent = 92.06
	progress.RuntimeObservation.PullRequestCoveragePercent = 92.10
	progress.SecurityObservation.RuntimeBoundaries = 6
	progress.SecurityObservation.RegisteredScopes = 6
	progress.SecurityObservation.AuthorizationRequired = 6
	progress.SecurityObservation.TenantIsolationRequired = 6
	progress.SecurityObservation.CoverageFloorPercent = 80
	progress.SecurityObservation.ChangedSecurityFloorPercent = 90
	progress.SecurityObservation.PullRequestChangedScopes = 6
	progress.Invariants = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	progress.NextUnlocked = []string{"E10-INV-006C-HANDOFF"}
	if err := validateSavedQueryAssetsRuntimeProgress(progress); err != nil {
		t.Fatal(err)
	}
	progress.RuntimeObservation.SASTFindings = 1
	if err := validateSavedQueryAssetsRuntimeProgress(progress); err == nil {
		t.Fatal("expected SAST evidence rejection")
	}
}

func TestValidateSavedQueryAssetsHandoffProgress(t *testing.T) {
	progress := SavedQueryAssetsHandoffProgress{
		SchemaVersion:      1,
		WorkUnit:           "E10-INV-006C-HANDOFF",
		Status:             "VERIFIED",
		FinalValidatedHead: "8af494e31839cce5ef65c55c09ea043642ec27d2",
	}
	progress.Validation.PushWorkflowRun = 36488802670
	progress.Validation.PullRequestWorkflowRun = 36488807591
	progress.Validation.Result = "PASS"
	progress.Validation.SavedQueryAssets = "PASS"
	progress.Validation.SecurityTests = "PASS"
	progress.Validation.SAST = "PASS"
	progress.Validation.SCA = "PASS"
	progress.Validation.RuntimeDependencies = "PASS"
	progress.Validation.BoundaryEdges = "PASS"
	progress.Validation.PerformanceRegistry = "PASS"
	progress.Invariants = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	progress.NextUnlocked = []string{"E10-INV-006D-CLOSURE"}
	h := &progress.HandoffObservation
	h.Immutable = true
	h.ExactQueryVersionPreserved = true
	h.OpaqueParameterValuesPreserved = true
	h.SourceContextCopied = true
	h.SourceReevaluationRequired = true
	h.PermissionReevaluationRequired = true
	h.ReturnContextAssetIdentityOnly = true
	h.PushGlobalCoveragePercent = 92.03
	h.PullRequestCoveragePercent = 92.03
	if err := validateSavedQueryAssetsHandoffProgress(progress); err != nil {
		t.Fatal(err)
	}
	h.SearchJobCreated = true
	if err := validateSavedQueryAssetsHandoffProgress(progress); err == nil {
		t.Fatal("expected Search Job overclaim rejection")
	}
}
