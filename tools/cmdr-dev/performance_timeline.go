package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	timelineBuildTimeout = 90 * time.Second
	timelineProbeTimeout = 25 * time.Second
)

const timelineMaxOutputBytes = 64 << 10

// timelineBoundedOutput limits allocations from malformed/misbehaving probes.
// Write reports all input consumed so the child is not blocked by the limit.
type timelineBoundedOutput struct {
	data     bytes.Buffer
	exceeded bool
}

func (out *timelineBoundedOutput) Write(p []byte) (int, error) {
	n := len(p)
	remaining := timelineMaxOutputBytes - out.data.Len()
	if remaining > 0 {
		keep := min(n, remaining)
		_, _ = out.data.Write(p[:keep])
	}
	if n > remaining {
		out.exceeded = true
	}
	return n, nil
}

func (out *timelineBoundedOutput) String() string {
	return out.data.String()
}

type timelineProbe struct {
	moduleRoot string
	binaryPath string
}

// The ordinary performance runner has no process deadline. Timeline probes
// are executed with a bounded child lifetime, including on deep-audit timeout.
func runTimelinePerformanceProcess(ctx context.Context, dir, name string, args ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("Timeline process cancelled: %w", err)
	}
	cmd := exec.CommandContext(ctx, name, args...) // #nosec G204,G702 -- constant Go tool or repo-compiled probe, no shell.
	configureTimelineProcessTree(cmd)
	cmd.WaitDelay = 500 * time.Millisecond
	cmd.Dir = dir
	var stdout, stderr timelineBoundedOutput
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	runErr := cmd.Run()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", fmt.Errorf("Timeline process timed out or cancelled: %w", ctxErr)
	}
	if stdout.exceeded || stderr.exceeded {
		return "", fmt.Errorf("Timeline process output exceeded %d bytes per stream", timelineMaxOutputBytes)
	}
	if runErr != nil {
		// Child stdout/stderr may contain local paths, access tokens or
		// sensitive source metadata. Keep the diagnostic in-memory only
		// and propagate the bounded exit status, never the child text.
		return "", fmt.Errorf("Timeline process failed (output withheld): %w", runErr)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func prepareTimelineProbe(root string) (timelineProbe, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), timelineBuildTimeout)
	defer cancel()
	return prepareTimelineProbeContext(ctx, root)
}

func prepareTimelineProbeContext(ctx context.Context, root string) (timelineProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/investigation-timeline", false)
	if err != nil {
		return timelineProbe{}, func() {}, fmt.Errorf("Timeline runtime unavailable: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-timeline-perf-")
	if err != nil {
		return timelineProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- directory created within verified repository root.
	}
	binaryName := "timeline-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runTimelinePerformanceProcess(ctx, moduleRoot, "go", "build", "-trimpath", "-o", binaryPath, "./cmd/perf-probe"); err != nil {
		cleanup()
		return timelineProbe{}, func() {}, fmt.Errorf("build Timeline probe: %w", err)
	}
	return timelineProbe{moduleRoot: moduleRoot, binaryPath: binaryPath}, cleanup, nil
}

func (probe timelineProbe) run(iterations int, mode string, memoryMiB int) (pilotProjectionProbeObservation, error) {
	return probe.runContext(context.Background(), iterations, mode, memoryMiB)
}

func (probe timelineProbe) runContext(ctx context.Context, iterations int, mode string, memoryMiB int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 100000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Timeline iterations outside 1..100000")
	}
	if mode != "latency" && mode != "resource" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("unsupported Timeline probe mode")
	}
	args := []string{"-iterations", strconv.Itoa(iterations), "-mode", mode}
	if mode == "resource" {
		args = append(args, "-memory-limit-mib", strconv.Itoa(memoryMiB))
	}
	runCtx, cancel := context.WithTimeout(ctx, timelineProbeTimeout)
	defer cancel()
	output, err := runTimelinePerformanceProcess(runCtx, probe.moduleRoot, probe.binaryPath, args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	return decodeTimelineProbeObservation([]byte(output), iterations, mode)
}

// decodeTimelineProbeObservation fails closed on mismatched measurement
// identities or impossible arithmetic. No untrusted child JSON or JSON
// decoder errors are reflected in CI logs.
func decodeTimelineProbeObservation(data []byte, iterations int, mode string) (pilotProjectionProbeObservation, error) {
	invalid := pilotProjectionProbeObservation{}
	if iterations < 1 || (mode != "latency" && mode != "resource") {
		return invalid, fmt.Errorf("invalid Timeline probe request identity")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var observation pilotProjectionProbeObservation
	if decoder.Decode(&observation) != nil {
		return invalid, fmt.Errorf("invalid Timeline probe JSON")
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF {
		return invalid, fmt.Errorf("Timeline probe output must contain exactly one JSON object")
	}
	if observation.Operation != "projection" || observation.Mode != mode ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) ||
		observation.PeakHeapBytes == 0 || observation.TotalAllocationBytes == 0 {
		return invalid, fmt.Errorf("invalid Timeline probe observation")
	}
	perOperation := float64(observation.ElapsedNS) / float64(iterations)
	if math.Abs(observation.NSPerOperation-perOperation) >
		math.Max(1e-6, perOperation*1e-9) {
		return invalid, fmt.Errorf("inconsistent Timeline probe elapsed time and sample mean")
	}
	return observation, nil
}
