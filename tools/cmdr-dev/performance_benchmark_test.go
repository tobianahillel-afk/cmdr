package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func benchmarkTargetForTest() PerformanceTarget {
	target := validRuntimePerformanceTarget()
	target.EnvironmentID = "PERF-ENV-TEST"
	target.Stages = []string{"pr"}
	target.TriggerPaths = []string{"runtime/a/**"}
	target.Metrics = []PerformanceMetricBudget{{
		ID: "latency-p95", Kind: "latency", Unit: "ms", Aggregation: "p95",
		Comparator: "upper-bound", Budget: 20, MaxRelativeRegressionPercent: 0,
	}}
	return target
}

func latencySamplesMS(values ...float64) []BenchmarkSample {
	var samples []BenchmarkSample
	for _, value := range values {
		samples = append(samples, BenchmarkSample{Measurements: []BenchmarkMeasurement{{
			Kind: "latency", Unit: "ms", Value: value,
		}}})
	}
	return samples
}

func TestAggregateBenchmarkValues(t *testing.T) {
	values := []float64{1, 2, 3, 4, 100}
	cases := map[string]float64{"median": 3, "p95": 100, "p99": 100, "max": 100, "mean": 22}
	for aggregation, want := range cases {
		got, err := aggregateBenchmarkValues(values, aggregation)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%s: got %v want %v", aggregation, got, want)
		}
	}
}

func TestAbsoluteUpperBoundRegressionFails(t *testing.T) {
	target := benchmarkTargetForTest()
	target.Metrics[0].Budget = 5
	_, err := evaluatePerformanceMetrics(target, validPerformanceEnvironment(), latencySamplesMS(6, 7, 8), nil)
	if err == nil || !strings.Contains(err.Error(), "absolute") {
		t.Fatalf("expected absolute budget failure, got %v", err)
	}
}

func TestAbsoluteLowerBoundRegressionFails(t *testing.T) {
	target := benchmarkTargetForTest()
	target.Metrics[0] = PerformanceMetricBudget{
		ID: "throughput", Kind: "throughput", Unit: "ops/s", Aggregation: "median",
		Comparator: "lower-bound", Budget: 100,
	}
	samples := []BenchmarkSample{
		{Measurements: []BenchmarkMeasurement{{Kind: "throughput", Unit: "ops/s", Value: 80}}},
		{Measurements: []BenchmarkMeasurement{{Kind: "throughput", Unit: "ops/s", Value: 90}}},
		{Measurements: []BenchmarkMeasurement{{Kind: "throughput", Unit: "ops/s", Value: 95}}},
	}
	_, err := evaluatePerformanceMetrics(target, validPerformanceEnvironment(), samples, nil)
	if err == nil || !strings.Contains(err.Error(), "absolute") {
		t.Fatalf("expected lower-bound failure, got %v", err)
	}
}

func TestRelativeRegressionFails(t *testing.T) {
	target := benchmarkTargetForTest()
	target.Metrics[0].Budget = 100
	target.Metrics[0].Aggregation = "median"
	target.Metrics[0].MaxRelativeRegressionPercent = 10
	baseline := &PerformanceBaseline{MetricValues: map[string]float64{"latency-p95": 10}}
	_, err := evaluatePerformanceMetrics(target, validPerformanceEnvironment(), latencySamplesMS(12, 12, 12), baseline)
	if err == nil || !strings.Contains(err.Error(), "regressed") {
		t.Fatalf("expected relative regression failure, got %v", err)
	}
}

func TestRelativeImprovementPasses(t *testing.T) {
	target := benchmarkTargetForTest()
	target.Metrics[0].Budget = 100
	target.Metrics[0].Aggregation = "median"
	target.Metrics[0].MaxRelativeRegressionPercent = 10
	baseline := &PerformanceBaseline{MetricValues: map[string]float64{"latency-p95": 10}}
	results, err := evaluatePerformanceMetrics(target, validPerformanceEnvironment(), latencySamplesMS(8, 8, 9), baseline)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].RelativePass || results[0].RelativeRegressionPercent != 0 {
		t.Fatalf("unexpected relative result: %#v", results)
	}
}

func TestTargetSelectionUsesStageEnvironmentAndChangedPath(t *testing.T) {
	target := benchmarkTargetForTest()
	selected := selectPerformanceTargets([]PerformanceTarget{target}, "pr", "PERF-ENV-TEST", []string{"runtime/a/file.go"})
	if len(selected) != 1 {
		t.Fatalf("expected target selection, got %#v", selected)
	}
	if got := selectPerformanceTargets([]PerformanceTarget{target}, "pr", "PERF-ENV-TEST", []string{"other/file.go"}); len(got) != 0 {
		t.Fatalf("unrelated change selected target: %#v", got)
	}
	if got := selectPerformanceTargets([]PerformanceTarget{target}, "release", "PERF-ENV-TEST", nil); len(got) != 0 {
		t.Fatalf("wrong stage selected target: %#v", got)
	}
}

func TestBenchmarkResultIdentityIsDigestBound(t *testing.T) {
	target := benchmarkTargetForTest()
	env := validPerformanceEnvironment()
	workload := BenchmarkWorkloadDefinition{
		Key: "cmdr-dev-metadata-audit-v1", HandlerKey: "builtin-cmdr-dev-metadata-audit-v1",
		WarmupIterations: 1, SampleCount: 3, MaxSampleDurationMS: 1000,
	}
	result := PerformanceBenchmarkResult{
		TargetID: target.ID, EnvironmentID: env.ID,
		SourceSHA: strings.Repeat("a", 40), Toolchain: "go" + engineeringGoVersion,
		TargetDigest: digestCanonical(target), EnvironmentDigest: digestCanonical(env),
		WorkloadDigest: digestCanonical(workload),
	}
	if result.TargetDigest == "" || result.EnvironmentDigest == "" || result.WorkloadDigest == "" ||
		result.SourceSHA == "" || result.Toolchain == "" {
		t.Fatalf("incomplete result identity: %#v", result)
	}
}

func TestPilotProjectionProbeExecutesRealRuntime(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	probe, cleanup, err := preparePilotProjectionProbe(root)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	observation, err := probe.run(1000, "latency", 0)
	if err != nil {
		t.Fatal(err)
	}
	if observation.Operations != 1000 || observation.NSPerOperation <= 0 {
		t.Fatalf("unexpected pilot projection observation: %#v", observation)
	}
}

func TestEventSearchBenchmarkHandlerFailsClosedBeforeRuntime(t *testing.T) {
	root := t.TempDir()
	workload := BenchmarkWorkloadDefinition{
		Key:                 "event-search-validation-v1",
		HandlerKey:          "builtin-event-search-validation-v1",
		WarmupIterations:    1,
		SampleCount:         1,
		MaxSampleDurationMS: 100,
	}
	_, cleanup, err := prepareBenchmarkHandler(root, workload)
	cleanup()
	if err == nil {
		t.Fatal("expected Event Search benchmark handler to fail closed before runtime materialization")
	}
}
