package main

import (
	"strings"
	"testing"
	"time"
)

func completeSnapshotRegistry() EngineMetricsRegistry {
	registry := validEngineMetricsRegistry()
	registry.Metrics = []EngineMetricDefinition{
		{ID:"MET-VALIDATION-SELECTED",Title:"selected",Category:"validation",Unit:"count",Aggregation:"last",SourceKey:"validation-plan",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"test"},
		{ID:"MET-VALIDATION-EXECUTED",Title:"executed",Category:"validation",Unit:"count",Aggregation:"last",SourceKey:"validation-execution",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"test"},
		{ID:"MET-VALIDATION-COST",Title:"cost",Category:"validation",Unit:"cost-unit",Aggregation:"last",SourceKey:"validation-plan",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"test"},
		{ID:"MET-MANDATORY-CHECKS",Title:"mandatory",Category:"quality",Unit:"count",Aggregation:"last",SourceKey:"check-catalog",SafetyClass:"safety-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"test"},
		{ID:"MET-SAFETY-GREEN",Title:"green",Category:"quality",Unit:"boolean",Aggregation:"last",SourceKey:"validation-execution",SafetyClass:"safety-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"test"},
		{ID:"MET-CONTEXT-SOURCES",Title:"sources",Category:"context",Unit:"count",Aggregation:"last",SourceKey:"context-bundle",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"test"},
		{ID:"MET-CONTEXT-DEPENDENCIES",Title:"deps",Category:"context",Unit:"count",Aggregation:"last",SourceKey:"context-bundle",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"test"},
		{ID:"MET-CACHE-REUSE-RATE",Title:"cache",Category:"cache",Unit:"percent",Aggregation:"ratio",SourceKey:"performance-cache",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"higher-better",Rationale:"test"},
		{ID:"MET-RECOVERY-REVALIDATE",Title:"revalidate",Category:"recovery",Unit:"boolean",Aggregation:"last",SourceKey:"recovery-reconcile",SafetyClass:"quality-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"test"},
		{ID:"MET-COORDINATION-CONFLICTS",Title:"conflicts",Category:"recovery",Unit:"count",Aggregation:"last",SourceKey:"coordination",SafetyClass:"safety-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"test"},
	}
	return registry
}

func TestBuildEngineMetricSnapshotAccountsForEveryMetric(t *testing.T) {
	registry := completeSnapshotRegistry()
	plan := ValidationPlan{
		WorkUnit:"E8-MET-001B",Tier:"standard",CostUnits:12,
		SelectedChecks:[]PlannedCheck{{ID:"CHK-A",Mandatory:true},{ID:"CHK-B",Mandatory:true}},
	}
	execution := ValidationExecutionSummary{
		WorkUnit:"E8-MET-001B",Tier:"standard",SelectedChecks:2,ExecutedChecks:1,PreflightSatisfied:1,
		Checks:[]ExecutedCheck{{ID:"CHK-A",Status:"PASS"},{ID:"CHK-B",Status:"PREFLIGHT_SATISFIED"}},
	}
	snapshot, err := buildEngineMetricSnapshot(
		registry, strings.Repeat("a",40), "2026-09-23T20:00:00Z",
		plan, execution,
		ContextSummary{Sources:14,DependencyUnits:3,BundleDigest:"ctx"},
		PerformanceCacheAuditSummary{Records:0,Status:"not-applicable"},
		RecoveryReconciliation{Outcome:"resume"},
		CoordinationAuditSummary{Conflicts:0,Status:"clear",AsOf:time.Date(2026,9,23,20,0,0,0,time.UTC).Format(time.RFC3339)},
		42,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Observations) != 9 || len(snapshot.Unavailable) != 1 ||
		snapshot.Unavailable[0].MetricID != "MET-CACHE-REUSE-RATE" || snapshot.Digest == "" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
}

func TestBuildEngineMetricSnapshotComputesCacheRate(t *testing.T) {
	registry := completeSnapshotRegistry()
	plan := ValidationPlan{WorkUnit:"E8-MET-001B",Tier:"fast",SelectedChecks:[]PlannedCheck{},CostUnits:0}
	execution := ValidationExecutionSummary{}
	snapshot, err := buildEngineMetricSnapshot(
		registry, strings.Repeat("b",40), "2026-09-23T20:00:00Z",
		plan, execution, ContextSummary{BundleDigest:"ctx"},
		PerformanceCacheAuditSummary{Records:4,Reusable:3,Status:"reusable"},
		RecoveryReconciliation{Outcome:"revalidate"}, CoordinationAuditSummary{Status:"clear"}, 0,
	)
	if err != nil {
		t.Fatal(err)
	}
	var cache, recovery float64 = -1, -1
	for _, o := range snapshot.Observations {
		if o.MetricID=="MET-CACHE-REUSE-RATE" { cache=o.Value }
		if o.MetricID=="MET-RECOVERY-REVALIDATE" { recovery=o.Value }
	}
	if cache != 75 || recovery != 1 {
		t.Fatalf("unexpected computed values cache=%v recovery=%v", cache, recovery)
	}
}

func TestSnapshotCoverageRejectsMissingAndDuplicateMetrics(t *testing.T) {
	registry := completeSnapshotRegistry()
	if err := validateEngineMetricSnapshotCoverage(registry,nil,nil); err == nil {
		t.Fatal("expected missing metric rejection")
	}
	obs := EngineMetricObservation{MetricID:"MET-VALIDATION-COST"}
	if err := validateEngineMetricSnapshotCoverage(registry,[]EngineMetricObservation{obs,obs},nil); err == nil {
		t.Fatal("expected duplicate metric rejection")
	}
}

func TestSafetyGreenDropsWhenMandatoryExecutionMissing(t *testing.T) {
	registry := completeSnapshotRegistry()
	plan := ValidationPlan{WorkUnit:"E8-MET-001B",Tier:"standard",SelectedChecks:[]PlannedCheck{{ID:"CHK-A",Mandatory:true}}}
	execution := ValidationExecutionSummary{Checks:nil}
	snapshot, err := buildEngineMetricSnapshot(
		registry, strings.Repeat("c",40), "2026-09-23T20:00:00Z", plan, execution,
		ContextSummary{BundleDigest:"ctx"}, PerformanceCacheAuditSummary{Records:0,Status:"not-applicable"},
		RecoveryReconciliation{Outcome:"resume"}, CoordinationAuditSummary{Status:"clear"}, 1,
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range snapshot.Observations {
		if o.MetricID=="MET-SAFETY-GREEN" && o.Value != 0 {
			t.Fatalf("missing mandatory execution must make safety green false: %#v", o)
		}
	}
}
