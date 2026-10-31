package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// These samples measure the Go test benchmark of the actual same-origin
// router. They exclude all external owner/provider work and are not a
// production latency or memory service-level claim.
func runCommandInvestigateRoutingDeep(ctx context.Context, root string, target DeepPerformanceTarget) (DeepPerformanceObservation, error) {
	if target.Kind != "resource" || target.Concurrency != 1 || target.DurationSeconds != 0 {
		return DeepPerformanceObservation{}, fmt.Errorf("Command/Investigate routing deep target supports only a bounded single-process resource probe")
	}
	dir, err := resolveRepoPath(root, "product-runtime/command-investigate-host", false)
	if err != nil {
		return DeepPerformanceObservation{}, err
	}
	cmd := exec.CommandContext(ctx, "go", "test", "-run", "^$", "-bench", "^BenchmarkSameOriginRouteResource$", "-benchtime="+strconv.Itoa(target.Iterations)+"x", "-benchmem", "-count=1", ".") // #nosec G204,G702 -- no shell; validated numeric iterations and fixed module/benchmark
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOMEMLIMIT="+strconv.Itoa(target.MemoryLimitMiB)+"MiB")
	started := time.Now()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return DeepPerformanceObservation{}, fmt.Errorf("real command-investigate resource benchmark failed: %s: %w", strings.TrimSpace(string(out)), err)
	}
	return parseCommandInvestigateRoutingDeep(string(out), target.Iterations, time.Since(started))
}

func parseCommandInvestigateRoutingDeep(out string, expectedIterations int, elapsed time.Duration) (DeepPerformanceObservation, error) {
	const name = "BenchmarkSameOriginRouteResource-"
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 || !strings.HasPrefix(fields[0], name) {
			continue
		}
		operations, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil || operations != int64(expectedIterations) || operations <= 0 {
			return DeepPerformanceObservation{}, fmt.Errorf("unexpected resource benchmark iteration count")
		}
		values := make(map[string]float64)
		for i := 2; i+1 < len(fields); i += 2 {
			value, err := strconv.ParseFloat(fields[i], 64)
			if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
				return DeepPerformanceObservation{}, fmt.Errorf("invalid resource benchmark value")
			}
			if _, exists := values[fields[i+1]]; exists {
				return DeepPerformanceObservation{}, fmt.Errorf("duplicate resource benchmark metric")
			}
			values[fields[i+1]] = value
		}
		ns, okLatency := values["ns/op"]
		peak, okPeak := values["peakheap-B"]
		total, okTotal := values["totalalloc-B"]
		if !okLatency || !okPeak || !okTotal || ns <= 0 || peak <= 0 || total <= 0 ||
			peak >= math.Exp2(53) || total >= math.Exp2(53) {
			return DeepPerformanceObservation{}, fmt.Errorf("missing or invalid real deep routing resource metrics")
		}
		return DeepPerformanceObservation{
			DurationMS: elapsed.Milliseconds(), PeakHeapBytes: uint64(peak),
			TotalAllocationBytes: uint64(total), Operations: operations,
			ThroughputOpsPerSec: 1e9 / ns,
		}, nil
	}
	return DeepPerformanceObservation{}, fmt.Errorf("real routing deep benchmark output missing")
}
