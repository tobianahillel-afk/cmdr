package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Measure the actual first-party HTTP host benchmark. This is a synthetic
// in-process load test; it cannot measure source-provider or network SLOs.
func measureEventInspectionHostLoad(ctx context.Context, root string, target DeepPerformanceTarget) (DeepPerformanceObservation, error) {
	if target.Iterations < 1 || target.Iterations > 100000 {
		return DeepPerformanceObservation{}, fmt.Errorf("invalid HTTP host iteration count")
	}
	dir, err := resolveRepoPath(root, "product-runtime/event-inspection-host", false)
	if err != nil {
		return DeepPerformanceObservation{}, err
	}
	// #nosec G204,G702 -- constant executable/benchmark and bounded decimal iterations; no shell.
	cmd := exec.CommandContext(ctx, "go", "test", "-run", "^$", "-bench", "^BenchmarkAuthorizedInspectionHTTP$", "-benchmem", "-count=1", "-benchtime="+strconv.Itoa(target.Iterations)+"x", ".")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return DeepPerformanceObservation{}, fmt.Errorf("authorized HTTP benchmark failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	for _, line := range strings.Split(stdout.String(), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 8 ||
			(fields[0] != "BenchmarkAuthorizedInspectionHTTP" && !strings.HasPrefix(fields[0], "BenchmarkAuthorizedInspectionHTTP-")) ||
			fields[3] != "ns/op" || fields[5] != "B/op" || fields[7] != "allocs/op" {
			continue
		}
		ops, e1 := strconv.ParseInt(fields[1], 10, 64)
		unsignedOps, e5 := strconv.ParseUint(fields[1], 10, 64)
		ns, e2 := strconv.ParseFloat(fields[2], 64)
		bytesPerOp, e3 := strconv.ParseUint(fields[4], 10, 64)
		allocationsPerOp, e4 := strconv.ParseUint(fields[6], 10, 64)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil || ops != int64(target.Iterations) ||
			ns <= 0 || ns > 1e12 || bytesPerOp > 1<<28 || allocationsPerOp > 1<<20 {
			return DeepPerformanceObservation{}, fmt.Errorf("invalid authorized HTTP benchmark measurement")
		}
		durationNS := float64(ops) * ns
		return DeepPerformanceObservation{
			DurationMS:           int64(durationNS / float64(time.Millisecond)),
			Operations:           ops,
			ThroughputOpsPerSec:  1e9 / ns,
			TotalAllocationBytes: unsignedOps * bytesPerOp,
			// Go benchmark reports B/op but not process peak heap; do not fabricate it.
		}, nil
	}
	return DeepPerformanceObservation{}, fmt.Errorf("no matching authorized HTTP benchmark observation")
}
