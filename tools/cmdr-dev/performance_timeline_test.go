package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
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

func TestTimelineSubprocessDescendantHelper(t *testing.T) {
	role := os.Getenv("CMDR_TIMELINE_TREE_HELPER")
	if role == "" {
		return
	}
	if role == "grandchild" {
		if err := os.WriteFile(os.Getenv("CMDR_TIMELINE_TREE_STARTED"), []byte("started"), 0o600); err != nil {
			t.Fatal(err)
		}
		time.Sleep(1500 * time.Millisecond)
		_ = os.WriteFile(os.Getenv("CMDR_TIMELINE_TREE_FINISHED"), []byte("escaped"), 0o600)
		return
	}
	if role != "parent" {
		t.Fatal("invalid helper role")
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestTimelineSubprocessDescendantHelper$")
	cmd.Env = append(os.Environ(), "CMDR_TIMELINE_TREE_HELPER=grandchild")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestTimelineProcessCancellationKillsDescendants(t *testing.T) {
	if os.Getenv("CMDR_TIMELINE_TREE_HELPER") != "" {
		t.Skip("helper mode")
	}
	root := t.TempDir()
	started := filepath.Join(root, "started")
	finished := filepath.Join(root, "finished")
	t.Setenv("CMDR_TIMELINE_TREE_HELPER", "parent")
	t.Setenv("CMDR_TIMELINE_TREE_STARTED", started)
	t.Setenv("CMDR_TIMELINE_TREE_FINISHED", finished)
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()
	_, err := runTimelinePerformanceProcess(ctx, root, os.Args[0], "-test.run=^TestTimelineSubprocessDescendantHelper$")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected child timeout: %v", err)
	}
	if _, err := os.Stat(started); err != nil {
		t.Fatalf("grandchild never ran, test not falsifiable: %v", err)
	}
	time.Sleep(1600 * time.Millisecond)
	if _, err := os.Stat(finished); !os.IsNotExist(err) {
		t.Fatalf("grandchild survived cancellation: %v", err)
	}
}

func TestTimelineOutputHelper(t *testing.T) {
	mode := os.Getenv("CMDR_TIMELINE_OUTPUT_HELPER")
	if mode == "" {
		return
	}
	if mode == "stdout" {
		_, _ = os.Stdout.Write([]byte(strings.Repeat("X", timelineMaxOutputBytes+1)))
		return
	}
	if mode == "stderr" {
		_, _ = os.Stderr.Write([]byte(strings.Repeat("Y", timelineMaxOutputBytes+1)))
		return
	}
	t.Fatal("unrecognized helper mode")
}

func TestTimelineProcessBoundsStdoutAndStderr(t *testing.T) {
	for _, mode := range []string{"stdout", "stderr"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("CMDR_TIMELINE_OUTPUT_HELPER", mode)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := runTimelinePerformanceProcess(ctx, t.TempDir(), os.Args[0], "-test.run=^TestTimelineOutputHelper$")
			if err == nil || !strings.Contains(err.Error(), "output exceeded") {
				t.Fatalf("unbounded %s stream accepted: %v", mode, err)
			}
		})
	}
}

func TestTimelineBoundedWriterDoesNotBlockOnOverflow(t *testing.T) {
	var out timelineBoundedOutput
	a := []byte(strings.Repeat("a", timelineMaxOutputBytes))
	if n, err := out.Write(a); err != nil || n != len(a) || out.exceeded {
		t.Fatalf("exact bound rejected: n=%d err=%v exceeded=%v", n, err, out.exceeded)
	}
	if n, err := out.Write([]byte("secret")); err != nil || n != 6 || !out.exceeded {
		t.Fatalf("overflow not reported without blocking: n=%d err=%v", n, err)
	}
	if len(out.String()) != timelineMaxOutputBytes || strings.Contains(out.String(), "secret") {
		t.Fatal("buffer grew beyond hard limit or exposed overflow")
	}
}

func TestTimelineFailedProcessOutputIsNotReflectedInError(t *testing.T) {
	if os.Getenv("CMDR_TIMELINE_FAILURE_HELPER") == "emit" {
		_, _ = os.Stdout.WriteString("SENSITIVE-TIMELINE-STDOUT-MARKER")
		_, _ = os.Stderr.WriteString("SENSITIVE-TIMELINE-STDERR-MARKER")
		t.Fatal("expected child failure")
	}
}

func TestTimelineSubprocessFailureRedactsChildStreams(t *testing.T) {
	t.Setenv("CMDR_TIMELINE_FAILURE_HELPER", "emit")
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	out, err := runTimelinePerformanceProcess(
		ctx, t.TempDir(), os.Args[0],
		"-test.run=^TestTimelineFailedProcessOutputIsNotReflectedInError$",
	)
	if err == nil || out != "" || !strings.Contains(err.Error(), "output withheld") {
		t.Fatalf("failed Timeline child did not fail closed: output=%q err=%v", out, err)
	}
	if strings.Contains(err.Error(), "SENSITIVE-TIMELINE-STDOUT-MARKER") ||
		strings.Contains(err.Error(), "SENSITIVE-TIMELINE-STDERR-MARKER") {
		t.Fatalf("Timeline child diagnostic leaked in error: %v", err)
	}
}
