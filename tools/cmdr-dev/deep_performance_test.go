package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

func validDeepPerformanceCaps() DeepPerformanceCaps {
	return DeepPerformanceCaps{
		MaxTargetDurationSeconds: 600,
		MaxTargetMemoryMiB:       2048,
		MaxConcurrency:           64,
		MaxTargetsPerRun:         2,
	}
}

func validDeepPerformanceTarget() DeepPerformanceTarget {
	return DeepPerformanceTarget{
		ID:                  "PERF-DEEP-RUNTIME-A-PROFILE",
		PerformanceTargetID: "PERF-TGT-RUNTIME-A",
		Kind:                "profile",
		HandlerKey:          "builtin-cmdr-dev-deep-metadata-v1",
		Stages:              []string{"nightly", "release", "on-demand"},
		TriggerPaths:        []string{"runtime/a/**"},
		TimeoutSeconds:      30,
		MemoryLimitMiB:      256,
		Concurrency:         1,
		Iterations:          3,
		Rationale:           "synthetic representative deep performance target",
	}
}

func validDeepPerformancePolicy(targets ...DeepPerformanceTarget) DeepPerformancePolicy {
	return DeepPerformancePolicy{
		SchemaVersion:     1,
		DefaultPolicy:     "deny-unregistered-deep-performance-target",
		EvidenceDirectory: "engineering/testing/deep-performance-evidence",
		ControlPlanePaths: []string{"engineering/performance/**", "engineering/architecture/**"},
		Caps:              validDeepPerformanceCaps(),
		Targets:           targets,
	}
}

func runtimePerformanceRegistry() PerformanceRegistry {
	return PerformanceRegistry{
		SchemaVersion: 1,
		RegistryKind:  "performance-targets",
		Environments:  []PerformanceEnvironment{validPerformanceEnvironment()},
		Targets:       []PerformanceTarget{validRuntimePerformanceTarget()},
	}
}

func TestDeepPerformanceNoRuntimeNoTargetsIsValid(t *testing.T) {
	if err := validateDeepPerformancePolicy(
		validDeepPerformancePolicy(),
		PerformanceRegistry{SchemaVersion: 1, RegistryKind: "performance-targets", Environments: []PerformanceEnvironment{validPerformanceEnvironment()}},
		performanceArchitecture(false),
	); err != nil {
		t.Fatal(err)
	}
}

func TestDeepPerformanceRuntimeTargetRequiresDeepCoverage(t *testing.T) {
	err := validateDeepPerformancePolicy(validDeepPerformancePolicy(), runtimePerformanceRegistry(), performanceArchitecture(true))
	if err == nil || !strings.Contains(err.Error(), "has no deep performance target") {
		t.Fatalf("expected missing deep coverage rejection, got %v", err)
	}
}

func TestDeepPerformanceAcceptsBoundedRuntimeTarget(t *testing.T) {
	if err := validateDeepPerformancePolicy(
		validDeepPerformancePolicy(validDeepPerformanceTarget()),
		runtimePerformanceRegistry(),
		performanceArchitecture(true),
	); err != nil {
		t.Fatal(err)
	}
}

func TestDeepPerformanceRejectsResourceCapViolations(t *testing.T) {
	for _, mutate := range []func(*DeepPerformanceTarget){
		func(target *DeepPerformanceTarget) { target.TimeoutSeconds = 601 },
		func(target *DeepPerformanceTarget) { target.MemoryLimitMiB = 2049 },
		func(target *DeepPerformanceTarget) { target.Concurrency = 65 },
	} {
		target := validDeepPerformanceTarget()
		mutate(&target)
		if err := validateDeepPerformancePolicy(
			validDeepPerformancePolicy(target),
			runtimePerformanceRegistry(),
			performanceArchitecture(true),
		); err == nil {
			t.Fatal("expected resource cap rejection")
		}
	}
}

func TestDeepPerformanceLoadAndSoakSemantics(t *testing.T) {
	load := validDeepPerformanceTarget()
	load.Kind = "load"
	load.Concurrency = 1
	if err := validateDeepPerformancePolicy(validDeepPerformancePolicy(load), runtimePerformanceRegistry(), performanceArchitecture(true)); err == nil {
		t.Fatal("expected load concurrency rejection")
	}

	soak := validDeepPerformanceTarget()
	soak.Kind = "soak"
	soak.DurationSeconds = 0
	if err := validateDeepPerformancePolicy(validDeepPerformancePolicy(soak), runtimePerformanceRegistry(), performanceArchitecture(true)); err == nil {
		t.Fatal("expected soak duration rejection")
	}
}

func TestSelectDeepPerformanceTargetsByStageAndRequest(t *testing.T) {
	target := validDeepPerformanceTarget()
	policy := validDeepPerformancePolicy(target)

	selected, err := selectDeepPerformanceTargets(policy, "nightly", "")
	if err != nil || len(selected) != 1 || selected[0].ID != target.ID {
		t.Fatalf("unexpected nightly selection: %#v %v", selected, err)
	}
	selected, err = selectDeepPerformanceTargets(policy, "on-demand", target.ID)
	if err != nil || len(selected) != 1 || selected[0].ID != target.ID {
		t.Fatalf("unexpected on-demand selection: %#v %v", selected, err)
	}
	if _, err := selectDeepPerformanceTargets(policy, "release", target.ID); err == nil {
		t.Fatal("expected requested target outside on-demand to fail")
	}
}

func TestDeepPerformanceSensitiveChange(t *testing.T) {
	target := validDeepPerformanceTarget()
	policy := validDeepPerformancePolicy(target)
	if !deepPerformanceSensitiveChange([]string{"runtime/a/service.go"}, policy) {
		t.Fatal("expected runtime trigger to be sensitive")
	}
	if deepPerformanceSensitiveChange([]string{"docs/readme.md"}, policy) {
		t.Fatal("unrelated path must not trigger deep performance")
	}
}

func TestDeepPerformanceTimeoutFailsClosed(t *testing.T) {
	target := validDeepPerformanceTarget()
	handler := func(context.Context, string, DeepPerformanceTarget) (DeepPerformanceObservation, error) {
		time.Sleep(100 * time.Millisecond)
		return DeepPerformanceObservation{DurationMS: 100, PeakHeapBytes: 1, Operations: 1}, nil
	}
	_, err := executeDeepPerformanceTargetWithHandler("", validDeepPerformanceCaps(), target, 10*time.Millisecond, handler)
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("expected timeout failure, got %v", err)
	}
}

func TestDeepPerformanceMemoryExhaustionFailsClosed(t *testing.T) {
	target := validDeepPerformanceTarget()
	observation := DeepPerformanceObservation{
		DurationMS:    10,
		PeakHeapBytes: uint64(target.MemoryLimitMiB+1) * 1024 * 1024,
		Operations:    1,
	}
	if err := validateDeepPerformanceObservation(target, observation); err == nil {
		t.Fatal("expected memory limit rejection")
	}
}

func TestDeepPerformanceProfileRequiresArtifact(t *testing.T) {
	target := validDeepPerformanceTarget()
	observation := DeepPerformanceObservation{DurationMS: 10, PeakHeapBytes: 1, Operations: 1}
	if err := validateDeepPerformanceObservation(target, observation); err == nil {
		t.Fatal("expected missing profile artifact rejection")
	}
}
