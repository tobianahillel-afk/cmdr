package main

import "testing"

func validPerformancePolicy() PerformancePolicy {
	return PerformancePolicy{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unregistered-performance-target",
		ScopeKinds:  []string{"engineering-control-plane", "product-runtime"},
		MetricKinds: []string{"latency", "throughput", "cpu-time", "memory", "allocation", "utilization"},
		MetricUnits: map[string][]string{
			"latency":     {"ns", "us", "ms", "s"},
			"throughput":  {"ops/s", "events/s", "bytes/s"},
			"cpu-time":    {"ns/op", "us/op", "ms/op"},
			"memory":      {"bytes", "MiB"},
			"allocation":  {"bytes/op", "allocs/op"},
			"utilization": {"percent"},
		},
		Aggregations:     []string{"median", "p95", "p99", "max", "mean"},
		Comparators:      []string{"upper-bound", "lower-bound"},
		EnvironmentKinds: []string{"shared-ci", "dedicated", "local-calibrated"},
		Stages:           []string{"pr", "nightly", "release", "on-demand"},
		WorkloadKeys:     []string{"cmdr-dev-metadata-audit-v1"},
	}
}

func validPerformanceEnvironment() PerformanceEnvironment {
	return PerformanceEnvironment{
		ID: "PERF-ENV-TEST", Kind: "dedicated", OS: "linux", Arch: "amd64",
		RunnerClass: "test-runner", RelativeRegressionAllowed: true,
		Rationale: "stable synthetic unit-test environment",
	}
}

func performanceArchitecture(withRuntime bool) ArchitectureRegistry {
	registry := validArchitectureRegistry()
	if withRuntime {
		registry.Boundaries = append(registry.Boundaries, ArchitectureBoundary{
			ID: "runtime-a", Kind: "product-runtime", Roots: []string{"runtime/a/**"},
			MutableByImplementation: true, MayDependOn: []string{"engineering"},
		})
	}
	return registry
}

func validRuntimePerformanceTarget() PerformanceTarget {
	return PerformanceTarget{
		ID: "PERF-TGT-RUNTIME-A", ScopeKind: "product-runtime", BoundaryID: "runtime-a",
		OwnerPath: "runtime/a/**", WorkloadKey: "cmdr-dev-metadata-audit-v1",
		WorkloadDescription: "representative deterministic workload", EnvironmentID: "PERF-ENV-TEST",
		Stages: []string{"pr", "release"}, TriggerPaths: []string{"runtime/a/**"},
		Metrics: []PerformanceMetricBudget{{
			ID: "latency-p95", Kind: "latency", Unit: "ms", Aggregation: "p95",
			Comparator: "upper-bound", Budget: 100, MaxRelativeRegressionPercent: 10,
		}},
		Rationale: "protect runtime latency before release",
	}
}

func TestPerformanceRegistryNoRuntimeIsNotApplicable(t *testing.T) {
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
	}
	summary, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(false))
	if err != nil {
		t.Fatal(err)
	}
	if summary.RuntimeCoverageStatus != "not-applicable" || summary.RuntimeBoundaries != 0 {
		t.Fatalf("unexpected no-runtime summary: %#v", summary)
	}
}

func TestPerformanceRegistryRequiresTargetForRuntimeBoundary(t *testing.T) {
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
	}
	if _, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(true)); err == nil {
		t.Fatal("expected missing runtime target rejection")
	}
}

func TestPerformanceRegistryAcceptsOwnedRuntimeTarget(t *testing.T) {
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
		Targets:      []PerformanceTarget{validRuntimePerformanceTarget()},
	}
	summary, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(true))
	if err != nil {
		t.Fatal(err)
	}
	if summary.RuntimeCoverageStatus != "complete" || summary.RuntimeTargetedBoundaries != 1 {
		t.Fatalf("unexpected runtime coverage: %#v", summary)
	}
}

func TestPerformanceRegistryRejectsUnknownRuntimeBoundary(t *testing.T) {
	target := validRuntimePerformanceTarget()
	target.BoundaryID = "missing"
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
		Targets:      []PerformanceTarget{target},
	}
	if _, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(false)); err == nil {
		t.Fatal("expected unknown runtime boundary rejection")
	}
}

func TestPerformanceRegistryRejectsOwnerOutsideBoundary(t *testing.T) {
	target := validRuntimePerformanceTarget()
	target.OwnerPath = "elsewhere/**"
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
		Targets:      []PerformanceTarget{target},
	}
	if _, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(true)); err == nil {
		t.Fatal("expected owner-path rejection")
	}
}

func TestPerformanceRegistryRejectsInvalidMetricUnitAndBudget(t *testing.T) {
	for _, mutate := range []func(*PerformanceTarget){
		func(target *PerformanceTarget) { target.Metrics[0].Unit = "bananas" },
		func(target *PerformanceTarget) { target.Metrics[0].Budget = 0 },
	} {
		target := validRuntimePerformanceTarget()
		mutate(&target)
		registry := PerformanceRegistry{
			SchemaVersion: 1, RegistryKind: "performance-targets",
			Environments: []PerformanceEnvironment{validPerformanceEnvironment()},
			Targets: []PerformanceTarget{target},
		}
		if _, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(true)); err == nil {
			t.Fatal("expected invalid metric rejection")
		}
	}
}

func TestSharedCIRejectsRelativeRegressionBudget(t *testing.T) {
	env := validPerformanceEnvironment()
	env.Kind = "shared-ci"
	env.RelativeRegressionAllowed = false
	target := validRuntimePerformanceTarget()
	registry := PerformanceRegistry{
		SchemaVersion: 1, RegistryKind: "performance-targets",
		Environments: []PerformanceEnvironment{env},
		Targets:      []PerformanceTarget{target},
	}
	if _, err := validatePerformanceRegistry(validPerformancePolicy(), registry, performanceArchitecture(true)); err == nil {
		t.Fatal("expected relative-regression rejection on shared CI")
	}
}
