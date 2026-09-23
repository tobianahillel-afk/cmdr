package main

import (
	"strings"
	"testing"
)

func testEfficiencyPolicy() EngineEfficiencyBudgetPolicy {
	return EngineEfficiencyBudgetPolicy{
		SchemaVersion: 1,
		PolicyKind:    "engineering-engine-efficiency-budgets",
		DefaultPolicy: "deny-unknown-metric",
		ProfileMatch:  []string{"plan_tier", "risk_domains", "mandatory_floor"},
		MandatoryFloorProvenance: EngineEfficiencyFloorProvenance{
			SourceSHA: strings.Repeat("a", 40), WorkflowRunID: 1, MandatoryFloor: 43,
		},
		SafetyFloors: []EngineEfficiencyFloorRule{
			{MetricID: "MET-MANDATORY-CHECKS", Comparator: "minimum", Value: 43, Rationale: "verified mandatory floor"},
			{MetricID: "MET-SAFETY-GREEN", Comparator: "equal", Value: 1, Rationale: "mandatory checks must pass"},
			{MetricID: "MET-COORDINATION-CONFLICTS", Comparator: "equal", Value: 0, Rationale: "ownership must be conflict free"},
			{MetricID: "MET-RECOVERY-REVALIDATE", Comparator: "equal", Value: 0, Rationale: "revalidation cannot be suppressed"},
		},
		EfficiencyBudgets: []EngineEfficiencyBudgetRule{
			{MetricID: "MET-VALIDATION-SELECTED", Blocking: true, MaxRelativeRegressionPercent: 0, Rationale: "deterministic same-profile selection"},
			{MetricID: "MET-VALIDATION-EXECUTED", Blocking: true, MaxRelativeRegressionPercent: 0, Rationale: "deterministic same-profile execution"},
			{MetricID: "MET-VALIDATION-COST", Blocking: true, MaxRelativeRegressionPercent: 0, Rationale: "deterministic same-profile cost"},
			{MetricID: "MET-CACHE-REUSE-RATE", Blocking: true, MaxRelativeRegressionPercent: 0, Rationale: "reuse must not regress"},
			{MetricID: "MET-CONTEXT-SOURCES", Blocking: false, MaxRelativeRegressionPercent: 0, Rationale: "context growth requires review"},
			{MetricID: "MET-CONTEXT-DEPENDENCIES", Blocking: false, MaxRelativeRegressionPercent: 0, Rationale: "dependency growth requires review"},
		},
	}
}

func budgetSafeValues() map[string]float64 {
	values := safeOptimizerValues()
	values["MET-MANDATORY-CHECKS"] = 43
	return values
}

func baselineForSnapshot(snapshot EngineMetricSnapshot, plan ValidationPlan) EngineEfficiencyBaseline {
	values, unavailable := snapshotMetricMaps(snapshot)
	var metrics []EngineEfficiencyBaselineMetric
	for id, value := range values {
		if !unavailable[id] {
			metrics = append(metrics, EngineEfficiencyBaselineMetric{MetricID: id, Value: value})
		}
	}
	return EngineEfficiencyBaseline{
		ID: "BASE-TEST", PlanTier: plan.Tier, RiskDomains: append([]string(nil), plan.RiskDomains...),
		MandatoryFloor: int(values["MET-MANDATORY-CHECKS"]), SourceSHA: snapshot.SourceSHA,
		ObservedAt: snapshot.ObservedAt, SnapshotDigest: snapshot.Digest, WorkflowRunID: 1, Metrics: metrics,
	}
}

func TestEfficiencyBudgetWithholdsWithoutComparableBaseline(t *testing.T) {
	registry := completeSnapshotRegistry()
	snapshot := optimizerSnapshot(registry, strings.Repeat("b", 40), "2026-09-23T21:00:00Z", "standard", budgetSafeValues(), nil)
	plan := ValidationPlan{WorkUnit: snapshot.WorkUnit, Tier: "standard", RiskDomains: []string{"performance", "security"}}
	report, err := evaluateEngineEfficiencyBudget(registry, testEfficiencyPolicy(),
		EngineEfficiencyBaselineRegistry{SchemaVersion: 1, BaselineKind: "engineering-engine-efficiency-baselines"},
		snapshot, plan)
	if err != nil {
		t.Fatal(err)
	}
	if report.SafetyStatus != "pass" || report.EfficiencyStatus != "not-enforced-no-baseline" ||
		len(report.BlockingViolations) != 0 || len(report.Withheld) != 6 {
		t.Fatalf("unexpected no-baseline budget report: %#v", report)
	}
}

func TestEfficiencyBudgetBlocksSafetyFloorWithoutBaseline(t *testing.T) {
	registry := completeSnapshotRegistry()
	values := budgetSafeValues()
	values["MET-SAFETY-GREEN"] = 0
	snapshot := optimizerSnapshot(registry, strings.Repeat("c", 40), "2026-09-23T21:00:00Z", "standard", values, nil)
	plan := ValidationPlan{WorkUnit: snapshot.WorkUnit, Tier: "standard", RiskDomains: []string{"security"}}
	report, err := evaluateEngineEfficiencyBudget(registry, testEfficiencyPolicy(),
		EngineEfficiencyBaselineRegistry{SchemaVersion: 1, BaselineKind: "engineering-engine-efficiency-baselines"},
		snapshot, plan)
	if err == nil || report.SafetyStatus != "fail" || len(report.BlockingViolations) == 0 {
		t.Fatalf("unsafe floor did not block: report=%#v err=%v", report, err)
	}
}

func TestEfficiencyBudgetBlocksSameProfileRegression(t *testing.T) {
	registry := completeSnapshotRegistry()
	baseValues := budgetSafeValues()
	currentValues := budgetSafeValues()
	currentValues["MET-VALIDATION-SELECTED"] = baseValues["MET-VALIDATION-SELECTED"] + 1
	currentValues["MET-VALIDATION-COST"] = baseValues["MET-VALIDATION-COST"] + 5
	baselineSnapshot := optimizerSnapshot(registry, strings.Repeat("d", 40), "2026-09-23T20:00:00Z", "standard", baseValues, nil)
	currentSnapshot := optimizerSnapshot(registry, strings.Repeat("e", 40), "2026-09-23T21:00:00Z", "standard", currentValues, nil)
	plan := ValidationPlan{WorkUnit: currentSnapshot.WorkUnit, Tier: "standard", RiskDomains: []string{"performance", "security"}}
	baseline := baselineForSnapshot(baselineSnapshot, plan)
	report, err := evaluateEngineEfficiencyBudget(registry, testEfficiencyPolicy(),
		EngineEfficiencyBaselineRegistry{SchemaVersion: 1, BaselineKind: "engineering-engine-efficiency-baselines", Baselines: []EngineEfficiencyBaseline{baseline}},
		currentSnapshot, plan)
	if err == nil || report.EfficiencyStatus != "fail" || len(report.BlockingViolations) < 2 {
		t.Fatalf("same-profile regression did not block: report=%#v err=%v", report, err)
	}
}

func TestEfficiencyBudgetContextGrowthIsAdvisory(t *testing.T) {
	registry := completeSnapshotRegistry()
	baseValues := budgetSafeValues()
	currentValues := budgetSafeValues()
	currentValues["MET-CONTEXT-SOURCES"]++
	baselineSnapshot := optimizerSnapshot(registry, strings.Repeat("f", 40), "2026-09-23T20:00:00Z", "standard", baseValues, nil)
	currentSnapshot := optimizerSnapshot(registry, strings.Repeat("1", 40), "2026-09-23T21:00:00Z", "standard", currentValues, nil)
	plan := ValidationPlan{WorkUnit: currentSnapshot.WorkUnit, Tier: "standard", RiskDomains: []string{"performance"}}
	baseline := baselineForSnapshot(baselineSnapshot, plan)
	report, err := evaluateEngineEfficiencyBudget(registry, testEfficiencyPolicy(),
		EngineEfficiencyBaselineRegistry{SchemaVersion: 1, BaselineKind: "engineering-engine-efficiency-baselines", Baselines: []EngineEfficiencyBaseline{baseline}},
		currentSnapshot, plan)
	if err != nil || report.EfficiencyStatus != "pass-with-advisories" || len(report.Advisories) != 1 {
		t.Fatalf("context growth should be advisory: report=%#v err=%v", report, err)
	}
}

func TestEfficiencyBudgetMandatoryIncreaseMakesProfileIncomparable(t *testing.T) {
	registry := completeSnapshotRegistry()
	baseValues := budgetSafeValues()
	currentValues := budgetSafeValues()
	currentValues["MET-MANDATORY-CHECKS"] = baseValues["MET-MANDATORY-CHECKS"] + 1
	baselineSnapshot := optimizerSnapshot(registry, strings.Repeat("2", 40), "2026-09-23T20:00:00Z", "standard", baseValues, nil)
	currentSnapshot := optimizerSnapshot(registry, strings.Repeat("3", 40), "2026-09-23T21:00:00Z", "standard", currentValues, nil)
	plan := ValidationPlan{WorkUnit: currentSnapshot.WorkUnit, Tier: "standard", RiskDomains: []string{"security"}}
	baseline := baselineForSnapshot(baselineSnapshot, plan)
	report, err := evaluateEngineEfficiencyBudget(registry, testEfficiencyPolicy(),
		EngineEfficiencyBaselineRegistry{SchemaVersion: 1, BaselineKind: "engineering-engine-efficiency-baselines", Baselines: []EngineEfficiencyBaseline{baseline}},
		currentSnapshot, plan)
	if err != nil || report.BaselineStatus != "unavailable" || report.EfficiencyStatus != "not-enforced-no-baseline" {
		t.Fatalf("mandatory-floor increase should require a new baseline: report=%#v err=%v", report, err)
	}
}
