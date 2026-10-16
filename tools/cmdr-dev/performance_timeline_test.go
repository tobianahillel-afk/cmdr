package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestTimelineSubprocessTimeoutHelper(t *testing.T) {
	if os.Getenv("CMDR_TIMELINE_TIMEOUT_HELPER") == "sleep" {
		time.Sleep(3 * time.Second)
	}
}

func TestTimelineProcessEnforcesChildDeadline(t *testing.T) {
	t.Setenv("CMDR_TIMELINE_TIMEOUT_HELPER", "sleep")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := runTimelinePerformanceProcess(ctx, t.TempDir(), os.Args[0], "-test.run=^TestTimelineSubprocessTimeoutHelper$")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("child deadline not enforced: %v", err)
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("cancelled child remained running: %s", elapsed)
	}
}

func TestTimelineProcessCancelAndSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := runTimelinePerformanceProcess(ctx, t.TempDir(), "go", "version")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled process was launched: %v", err)
	}
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	output, err := runTimelinePerformanceProcess(ctx, t.TempDir(), "go", "version")
	if err != nil || !strings.HasPrefix(output, "go version") {
		t.Fatalf("real subprocess failed: %q %v", output, err)
	}
}

func TestTimelineProbeContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (timelineProbe{moduleRoot: t.TempDir(), binaryPath: "should-not-run"}).runContext(ctx, 1, "latency", 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation was not propagated: %v", err)
	}
}
