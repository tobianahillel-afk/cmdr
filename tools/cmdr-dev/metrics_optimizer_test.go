package main

import (
	"strings"
	"testing"
	"time"
)

func optimizerSnapshot(registry EngineMetricsRegistry, sha, observedAt, tier string, values map[string]float64, unavailable map[string]string) EngineMetricSnapshot {
	var observations []EngineMetricObservation
	var missing []EngineMetricUnavailable
	for _, definition := range registry.Metrics {
		if reason, ok := unavailable[definition.ID]; ok {
			missing = append(missing, EngineMetricUnavailable{MetricID: definition.ID, Reason: reason})
			continue
		}
		observations = append(observations, EngineMetricObservation{
			MetricID: definition.ID, Value: values[definition.ID], SourceKey: definition.SourceKey,
			SourceSHA: sha, ObservedAt: observedAt, Evidence: "test:" + definition.ID,
		})
	}
	body := engineMetricSnapshotBody{
		SchemaVersion: 1, SnapshotKind: "engineering-engine-metrics", WorkUnit: "E8-MET-001C",
		SourceSHA: sha, ObservedAt: observedAt, PlanTier: tier, Observations: observations, Unavailable: missing,
	}
	return EngineMetricSnapshot{
		SchemaVersion: body.SchemaVersion, SnapshotKind: body.SnapshotKind, WorkUnit: body.WorkUnit,
		SourceSHA: body.SourceSHA, ObservedAt: body.ObservedAt, PlanTier: body.PlanTier,
		Observations: body.Observations, Unavailable: body.Unavailable, Digest: digestCanonical(body),
	}
}

func safeOptimizerValues() map[string]float64 {
	return map[string]float64{
		"MET-VALIDATION-SELECTED":    40,
		"MET-VALIDATION-EXECUTED":    36,
		"MET-VALIDATION-COST":        60,
		"MET-MANDATORY-CHECKS":       42,
		"MET-SAFETY-GREEN":           1,
		"MET-CONTEXT-SOURCES":        80,
		"MET-CONTEXT-DEPENDENCIES":   20,
		"MET-CACHE-REUSE-RATE":       80,
		"MET-RECOVERY-REVALIDATE":    0,
		"MET-COORDINATION-CONFLICTS": 0,
	}
}

func TestOptimizerWithholdsWithoutBaseline(t *testing.T) {
	registry := completeSnapshotRegistry()
	current := optimizerSnapshot(registry, strings.Repeat("a", 40), "2026-09-23T21:00:00Z", "standard", safeOptimizerValues(), nil)
	report, err := deriveEngineOptimizations(registry, current, nil)
	if err != nil {
		t.Fatal(err)
	}
	if report.SafetyStatus != "safe-no-baseline" || len(report.Recommendations) != 0 || len(report.Withheld) != 3 {
		t.Fatalf("unexpected no-baseline report: %#v", report)
	}
	if report.Proof.DirectPlanMutationAllowed || report.Proof.MandatoryMutationAllowed ||
		report.Proof.PrerequisiteMutationAllowed || report.Proof.ProductSpecMutationAllowed {
		t.Fatalf("optimizer proof allows forbidden mutation: %#v", report.Proof)
	}
}

func TestOptimizerDerivesOnlyRelativeAdvisoryRecommendations(t *testing.T) {
	registry := completeSnapshotRegistry()
	baseValues := safeOptimizerValues()
	currentValues := safeOptimizerValues()
	currentValues["MET-VALIDATION-SELECTED"] = 44
	currentValues["MET-VALIDATION-COST"] = 70
	currentValues["MET-CACHE-REUSE-RATE"] = 60
	currentValues["MET-CONTEXT-SOURCES"] = 95
	currentValues["MET-CONTEXT-DEPENDENCIES"] = 25

	baseline := optimizerSnapshot(registry, strings.Repeat("b", 40), "2026-09-23T20:00:00Z", "standard", baseValues, nil)
	current := optimizerSnapshot(registry, strings.Repeat("c", 40), "2026-09-23T21:00:00Z", "standard", currentValues, nil)
	report, err := deriveEngineOptimizations(registry, current, &baseline)
	if err != nil {
		t.Fatal(err)
	}
	if report.SafetyStatus != "safe" || len(report.Recommendations) != 3 || len(report.Withheld) != 0 {
		t.Fatalf("unexpected comparative report: %#v", report)
	}
	for _, recommendation := range report.Recommendations {
		if !recommendation.Advisory || len(recommendation.EstimatedBenefit) == 0 || len(recommendation.Constraints) == 0 {
			t.Fatalf("recommendation lacks advisory safety evidence: %#v", recommendation)
		}
	}
}

func TestOptimizerWithholdsUnsafeCurrentSnapshot(t *testing.T) {
	registry := completeSnapshotRegistry()
	values := safeOptimizerValues()
	values["MET-SAFETY-GREEN"] = 0
	current := optimizerSnapshot(registry, strings.Repeat("d", 40), "2026-09-23T21:00:00Z", "standard", values, nil)
	report, err := deriveEngineOptimizations(registry, current, nil)
	if err != nil {
		t.Fatal(err)
	}
	if report.SafetyStatus != "withheld-unsafe" || len(report.Recommendations) != 0 || len(report.Withheld) != 3 {
		t.Fatalf("unsafe snapshot should withhold optimization: %#v", report)
	}
}

func TestOptimizerWithholdsUnavailableCacheAndNoRegression(t *testing.T) {
	registry := completeSnapshotRegistry()
	values := safeOptimizerValues()
	baseline := optimizerSnapshot(registry, strings.Repeat("e", 40), "2026-09-23T20:00:00Z", "standard", values, map[string]string{"MET-CACHE-REUSE-RATE": "not-applicable"})
	current := optimizerSnapshot(registry, strings.Repeat("f", 40), "2026-09-23T21:00:00Z", "standard", values, map[string]string{"MET-CACHE-REUSE-RATE": "not-applicable"})
	report, err := deriveEngineOptimizations(registry, current, &baseline)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Recommendations) != 0 || len(report.Withheld) != 3 {
		t.Fatalf("expected all opportunities withheld: %#v", report)
	}
}

func TestOptimizerRejectsTamperedSnapshotDigest(t *testing.T) {
	registry := completeSnapshotRegistry()
	current := optimizerSnapshot(registry, strings.Repeat("1", 40), time.Date(2026, 9, 23, 21, 0, 0, 0, time.UTC).Format(time.RFC3339), "standard", safeOptimizerValues(), nil)
	current.Digest = strings.Repeat("0", 64)
	if _, err := deriveEngineOptimizations(registry, current, nil); err == nil {
		t.Fatal("expected tampered snapshot rejection")
	}
}
