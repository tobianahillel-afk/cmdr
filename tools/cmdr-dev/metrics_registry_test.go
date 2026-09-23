package main

import (
	"math"
	"strings"
	"testing"
)

func validEngineMetricsRegistry() EngineMetricsRegistry {
	return EngineMetricsRegistry{
		SchemaVersion: 1,
		RegistryKind: "engineering-engine-metrics",
		DefaultPolicy: "deny-unregistered-metric",
		Categories: []string{"validation","context","cache","recovery","quality"},
		Units: []string{"count","cost-unit","percent","boolean"},
		Aggregations: []string{"last","sum","max","ratio"},
		SourceKeys: []string{"validation-plan","validation-execution","check-catalog","context-bundle","performance-cache","recovery-reconcile","coordination"},
		SafetyClasses: []string{"safety-floor","quality-floor","efficiency"},
		Directions: []string{"lower-better","higher-better","invariant"},
		Metrics: []EngineMetricDefinition{
			{ID:"MET-VALIDATION-COST",Title:"Validation cost",Category:"validation",Unit:"cost-unit",Aggregation:"sum",SourceKey:"validation-plan",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"lower-better",Rationale:"reduce unnecessary validation work"},
			{ID:"MET-MANDATORY-CHECKS",Title:"Mandatory checks",Category:"quality",Unit:"count",Aggregation:"last",SourceKey:"check-catalog",SafetyClass:"safety-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"mandatory safety floor must not be optimized away"},
		},
	}
}

func TestEngineMetricsRegistryAcceptsClosedModel(t *testing.T) {
	summary, err := validateEngineMetricsRegistry(validEngineMetricsRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if summary.Metrics != 2 || summary.OptimizationEligible != 1 || summary.NonOptimizable != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestEngineMetricsRegistryRejectsSafetyOptimization(t *testing.T) {
	registry := validEngineMetricsRegistry()
	registry.Metrics[1].OptimizationEligible = true
	if _, err := validateEngineMetricsRegistry(registry); err == nil {
		t.Fatal("expected safety-floor optimization rejection")
	}
}

func TestEngineMetricsRegistryRejectsUnknownEnumsAndDuplicates(t *testing.T) {
	for _, mutate := range []func(*EngineMetricsRegistry){
		func(r *EngineMetricsRegistry){ r.Metrics[0].Unit = "tokens" },
		func(r *EngineMetricsRegistry){ r.Metrics = append(r.Metrics, r.Metrics[0]) },
	} {
		registry := validEngineMetricsRegistry()
		mutate(&registry)
		if _, err := validateEngineMetricsRegistry(registry); err == nil {
			t.Fatal("expected invalid registry rejection")
		}
	}
}

func TestEngineMetricObservationRequiresExactSourceAndBoundedValue(t *testing.T) {
	registry := validEngineMetricsRegistry()
	base := EngineMetricObservation{
		MetricID:"MET-VALIDATION-COST",Value:12,SourceKey:"validation-plan",
		SourceSHA:strings.Repeat("a",40),ObservedAt:"2026-09-23T20:00:00Z",Evidence:"validation-plan/current",
	}
	if err := validateEngineMetricObservation(registry, base); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*EngineMetricObservation){
		func(o *EngineMetricObservation){ o.SourceKey = "check-catalog" },
		func(o *EngineMetricObservation){ o.SourceSHA = "deadbeef" },
		func(o *EngineMetricObservation){ o.Value = math.NaN() },
		func(o *EngineMetricObservation){ o.Evidence = "" },
	} {
		value := base
		mutate(&value)
		if err := validateEngineMetricObservation(registry, value); err == nil {
			t.Fatalf("expected observation rejection: %#v", value)
		}
	}
}

func TestBooleanAndPercentageObservationBounds(t *testing.T) {
	registry := validEngineMetricsRegistry()
	registry.Metrics = append(registry.Metrics,
		EngineMetricDefinition{ID:"MET-CACHE-RATE",Title:"Cache reuse",Category:"cache",Unit:"percent",Aggregation:"ratio",SourceKey:"performance-cache",SafetyClass:"efficiency",OptimizationEligible:true,Direction:"higher-better",Rationale:"maximize valid evidence reuse"},
		EngineMetricDefinition{ID:"MET-SAFETY-GREEN",Title:"Safety green",Category:"quality",Unit:"boolean",Aggregation:"last",SourceKey:"validation-execution",SafetyClass:"safety-floor",OptimizationEligible:false,Direction:"invariant",Rationale:"safety gate state is invariant"},
	)
	if _, err := validateEngineMetricsRegistry(registry); err != nil {
		t.Fatal(err)
	}
	base := EngineMetricObservation{SourceSHA:strings.Repeat("b",40),ObservedAt:"2026-09-23T20:00:00Z",Evidence:"test"}
	for _, tc := range []struct{ id, source string; value float64 }{
		{"MET-CACHE-RATE","performance-cache",101},
		{"MET-SAFETY-GREEN","validation-execution",2},
	} {
		o := base
		o.MetricID, o.SourceKey, o.Value = tc.id, tc.source, tc.value
		if err := validateEngineMetricObservation(registry, o); err == nil {
			t.Fatalf("expected bounded value rejection for %s", tc.id)
		}
	}
}
