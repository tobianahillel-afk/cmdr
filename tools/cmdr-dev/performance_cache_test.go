package main

import "testing"

func cacheTarget() PerformanceTarget {
	target := validRuntimePerformanceTarget()
	target.Metrics = []PerformanceMetricBudget{{
		ID: "latency-p95", Kind: "latency", Unit: "ms", Aggregation: "p95",
		Comparator: "upper-bound", Budget: 100, MaxRelativeRegressionPercent: 10,
	}}
	return target
}

func cacheEnvironment() PerformanceEnvironment {
	env := validPerformanceEnvironment()
	env.RelativeRegressionAllowed = true
	return env
}

func cacheWorkload() BenchmarkWorkloadDefinition {
	return BenchmarkWorkloadDefinition{
		Key: "cmdr-dev-metadata-audit-v1", HandlerKey: "builtin-cmdr-dev-metadata-audit-v1",
		WarmupIterations: 1, SampleCount: 3, MaxSampleDurationMS: 1000,
	}
}

func cacheResult(target PerformanceTarget, env PerformanceEnvironment, workload BenchmarkWorkloadDefinition, source string, pass bool) PerformanceBenchmarkResult {
	return PerformanceBenchmarkResult{
		TargetID: target.ID, EnvironmentID: env.ID, SourceSHA: source, Toolchain: "go" + engineeringGoVersion,
		TargetDigest: digestCanonical(target), EnvironmentDigest: digestCanonical(env), WorkloadDigest: digestCanonical(workload),
		Samples: 3,
		Metrics: []BenchmarkMetricResult{{
			ID: "latency-p95", Kind: "latency", Unit: "ms", Aggregation: "p95",
			Observed: 90, Comparator: "upper-bound", Budget: 100, AbsolutePass: pass,
			Baseline: 85, RelativePass: pass, MaxRelativeRegressionPercent: 10,
		}},
	}
}

func cacheRecord(source, decisionBasis string, outcome string, decisionIDs []string) PerformanceCacheRecord {
	target := cacheTarget()
	env := cacheEnvironment()
	workload := cacheWorkload()
	result := cacheResult(target, env, workload, source, outcome == "pass")
	identity := performanceCacheIdentity{
		SourceSHA: source, TargetID: target.ID, TargetDigest: digestCanonical(target),
		BudgetDigest: digestCanonical(target.Metrics), EnvironmentID: env.ID,
		EnvironmentDigest: digestCanonical(env), WorkloadDigest: digestCanonical(workload),
		Toolchain: "go" + engineeringGoVersion, DecisionBasisDigest: decisionBasis,
	}
	return PerformanceCacheRecord{
		ID: "PERF-CACHE-0001", TargetID: target.ID, EnvironmentID: env.ID, SourceSHA: source,
		Toolchain: "go" + engineeringGoVersion, TargetDigest: identity.TargetDigest, BudgetDigest: identity.BudgetDigest,
		EnvironmentDigest: identity.EnvironmentDigest, WorkloadDigest: identity.WorkloadDigest,
		DecisionBasisDigest: decisionBasis, CacheKey: digestCanonical(identity), EvidenceDigest: digestCanonical(result),
		Outcome: outcome, DecisionIDs: decisionIDs, Result: result,
	}
}

func TestPerformanceCacheFreshPassIsReusable(t *testing.T) {
	source := "0123456789abcdef0123456789abcdef01234567"
	basis := digestCanonical([]ReusableEvidenceEntry{})
	record := cacheRecord(source, basis, "pass", nil)
	status, signals, _, _ := evaluatePerformanceCacheRecord(
		record, source, basis,
		map[string]PerformanceTarget{record.TargetID: cacheTarget()},
		map[string]PerformanceEnvironment{record.EnvironmentID: cacheEnvironment()},
		map[string]BenchmarkWorkloadDefinition{"cmdr-dev-metadata-audit-v1": cacheWorkload()},
		map[string]EngineeringDecision{}, map[string]bool{},
	)
	if !status.Reusable || len(status.Reasons) != 0 || len(signals) != 0 {
		t.Fatalf("expected reusable cache, got %#v signals=%#v", status, signals)
	}
}

func TestPerformanceCacheInvalidatesEveryIdentityDimension(t *testing.T) {
	source := "0123456789abcdef0123456789abcdef01234567"
	basis := digestCanonical([]ReusableEvidenceEntry{})
	base := cacheRecord(source, basis, "pass", nil)
	tests := []struct {
		name string
		mutate func(*PerformanceCacheRecord)
		currentSource string
		currentBasis string
		want string
	}{
		{"source", func(*PerformanceCacheRecord){}, "1123456789abcdef0123456789abcdef01234567", basis, "source-change"},
		{"toolchain", func(r *PerformanceCacheRecord){ r.Toolchain = "go0.0.0" }, source, basis, "toolchain-change"},
		{"target", func(r *PerformanceCacheRecord){ r.TargetDigest = digestCanonical("changed") }, source, basis, "target-change"},
		{"budget", func(r *PerformanceCacheRecord){ r.BudgetDigest = digestCanonical("changed") }, source, basis, "budget-change"},
		{"environment", func(r *PerformanceCacheRecord){ r.EnvironmentDigest = digestCanonical("changed") }, source, basis, "environment-change"},
		{"workload", func(r *PerformanceCacheRecord){ r.WorkloadDigest = digestCanonical("changed") }, source, basis, "workload-change"},
		{"decision", func(*PerformanceCacheRecord){}, source, digestCanonical("changed"), "research-decision-freshness-change"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			record := base
			tc.mutate(&record)
			status, _, _, _ := evaluatePerformanceCacheRecord(
				record, tc.currentSource, tc.currentBasis,
				map[string]PerformanceTarget{base.TargetID: cacheTarget()},
				map[string]PerformanceEnvironment{base.EnvironmentID: cacheEnvironment()},
				map[string]BenchmarkWorkloadDefinition{"cmdr-dev-metadata-audit-v1": cacheWorkload()},
				map[string]EngineeringDecision{}, map[string]bool{},
			)
			if status.Reusable || !containsString(status.Reasons, tc.want) {
				t.Fatalf("expected %s invalidation, got %#v", tc.want, status)
			}
		})
	}
}

func TestBenchmarkRegressionFiresE5Signal(t *testing.T) {
	source := "0123456789abcdef0123456789abcdef01234567"
	decisionID := "ENG-DEC-0001"
	entry := ReusableEvidenceEntry{DecisionID: decisionID}
	basis := digestCanonical([]ReusableEvidenceEntry{entry})
	record := cacheRecord(source, basis, "regression", []string{decisionID})
	decision := EngineeringDecision{ID: decisionID, Status: "accepted", PerformanceSensitive: true}
	status, signals, _, _ := evaluatePerformanceCacheRecord(
		record, source, basis,
		map[string]PerformanceTarget{record.TargetID: cacheTarget()},
		map[string]PerformanceEnvironment{record.EnvironmentID: cacheEnvironment()},
		map[string]BenchmarkWorkloadDefinition{"cmdr-dev-metadata-audit-v1": cacheWorkload()},
		map[string]EngineeringDecision{decisionID: decision}, map[string]bool{decisionID:true},
	)
	if status.Reusable {
		t.Fatal("regression must never be reusable")
	}
	if len(signals) != 1 || signals[0].DecisionID != decisionID || signals[0].Kind != "benchmark-regression" || signals[0].State != "fired" {
		t.Fatalf("unexpected revisit signal: %#v", signals)
	}
}

func TestRegressionRequiresValidPerformanceDecisionBinding(t *testing.T) {
	source := "0123456789abcdef0123456789abcdef01234567"
	basis := digestCanonical([]ReusableEvidenceEntry{})
	record := cacheRecord(source, basis, "regression", []string{"ENG-DEC-9999"})
	status, _, _, _ := evaluatePerformanceCacheRecord(
		record, source, basis,
		map[string]PerformanceTarget{record.TargetID: cacheTarget()},
		map[string]PerformanceEnvironment{record.EnvironmentID: cacheEnvironment()},
		map[string]BenchmarkWorkloadDefinition{"cmdr-dev-metadata-audit-v1": cacheWorkload()},
		map[string]EngineeringDecision{}, map[string]bool{},
	)
	if !containsString(status.Reasons, "decision-binding-invalid") {
		t.Fatalf("expected invalid decision binding, got %#v", status)
	}
}
