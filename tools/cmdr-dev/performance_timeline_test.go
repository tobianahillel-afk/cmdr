package main

import (
	"path/filepath"
	"testing"
)

func TestTimelineProbeRealExecutableMeasurement(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	probe, cleanup, err := prepareTimelineProbe(root)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	for _, mode := range []string{"latency", "resource"} {
		memory := 0
		if mode == "resource" {
			memory = 128
		}
		got, err := probe.run(2, mode, memory)
		if err != nil {
			t.Fatal(err)
		}
		if got.Operations != 2 || got.Operation != "projection" || got.PeakHeapBytes == 0 {
			t.Fatalf("invalid bound observation: %+v", got)
		}
	}
}

func TestTimelineExplicitBenchmarkSelectionIsBounded(t *testing.T) {
	targets := []PerformanceTarget{
		{ID: "PERF-TGT-INV-TIMELINE-PROJECTION", EnvironmentID: defaultCIEnvironment, Stages: []string{"on-demand"}},
		{ID: "PERF-TGT-INV-TIMELINE-NAVIGATION", EnvironmentID: defaultCIEnvironment, Stages: []string{"on-demand"}},
	}
	chosen, err := selectPerformanceTargetsByID(targets, "on-demand", defaultCIEnvironment, []string{"PERF-TGT-INV-TIMELINE-PROJECTION"})
	if err != nil || len(chosen) != 1 || chosen[0].ID != "PERF-TGT-INV-TIMELINE-PROJECTION" {
		t.Fatalf("Timeline explicit projection selection failed: %#v %v", chosen, err)
	}
	if _, err := selectPerformanceTargetsByID(targets, "pr", defaultCIEnvironment, []string{"PERF-TGT-INV-TIMELINE-PROJECTION"}); err == nil {
		t.Fatal("on-demand target was permitted in pr stage")
	}
	if _, err := selectPerformanceTargetsByID(targets, "on-demand", defaultCIEnvironment, []string{"not-a-target"}); err == nil {
		t.Fatal("unknown target was accepted")
	}
	if _, err := selectPerformanceTargetsByID(targets, "on-demand", defaultCIEnvironment, []string{"PERF-TGT-INV-TIMELINE-PROJECTION", "PERF-TGT-INV-TIMELINE-PROJECTION"}); err == nil {
		t.Fatal("duplicate target selection was accepted")
	}
}
