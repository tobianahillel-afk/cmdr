package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type commandWorkQueueProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareCommandWorkQueueProbe(root string) (commandWorkQueueProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/work-queue", false)
	if err != nil {
		return commandWorkQueueProbe{}, func() {}, fmt.Errorf("Command queue unavailable: %w", err)
	}
	dir, err := os.MkdirTemp(root, ".cmdr-command-work-queue-perf-")
	if err != nil {
		return commandWorkQueueProbe{}, func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) } // #nosec G703 - created below checked repository root
	name := "command-work-queue-perf"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(dir, name)
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binary, "./cmd/perf-probe"); err != nil {
		cleanup()
		return commandWorkQueueProbe{}, func() {}, fmt.Errorf("build Command queue probe: %w", err)
	}
	return commandWorkQueueProbe{moduleRoot: moduleRoot, binaryPath: binary}, cleanup, nil
}

func (p commandWorkQueueProbe) run(iterations int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 100000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("invalid iteration count")
	}
	raw, err := runPerformanceProcess(p.moduleRoot, p.binaryPath, "-iterations", strconv.Itoa(iterations), "-mode", "latency", "-operation", "projection")
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var sample pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &sample); err != nil {
		return sample, err
	}
	if sample.Operation != "projection" || sample.Mode != "latency" || sample.Operations != iterations ||
		sample.ElapsedNS <= 0 || sample.NSPerOperation <= 0 || math.IsNaN(sample.NSPerOperation) || math.IsInf(sample.NSPerOperation, 0) {
		return sample, fmt.Errorf("invalid Command queue probe output")
	}
	return sample, nil
}

type commandQueueResourceSample struct {
	Operation string `json:"operation"`
	Mode string `json:"mode"`
	Operations int `json:"operations"`
	ElapsedNS int64 `json:"elapsed_ns"`
	PeakHeapBytes uint64 `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64 `json:"total_allocation_bytes"`
}

// Deep validation executes the SAME compiled first-party Work Queue module,
// not a metadata benchmark pretending to exercise Command.
func (p commandWorkQueueProbe) runResource(iterations, memoryLimitMiB int) (DeepPerformanceObservation, error) {
	if iterations<1 || iterations>100000 || memoryLimitMiB<64 || memoryLimitMiB>8192 {
		return DeepPerformanceObservation{}, fmt.Errorf("invalid Command resource probe bounds")
	}
	raw,err:=runPerformanceProcess(p.moduleRoot,p.binaryPath,
		"-iterations",strconv.Itoa(iterations),"-mode","resource",
		"-operation","projection","-memory-limit-mib",strconv.Itoa(memoryLimitMiB))
	if err!=nil {return DeepPerformanceObservation{},err}
	var sample commandQueueResourceSample
	if err:=json.Unmarshal([]byte(strings.TrimSpace(raw)),&sample);err!=nil {
		return DeepPerformanceObservation{},err
	}
	if sample.Operation!="projection" || sample.Mode!="resource" || sample.Operations!=iterations ||
		sample.ElapsedNS<=0 || sample.PeakHeapBytes==0 || sample.TotalAllocationBytes==0 ||
		sample.PeakHeapBytes>uint64(memoryLimitMiB)*1024*1024 {
		return DeepPerformanceObservation{},fmt.Errorf("invalid Command resource observation")
	}
	return DeepPerformanceObservation{
		DurationMS: sample.ElapsedNS/int64(1000000),
		PeakHeapBytes: sample.PeakHeapBytes,
		TotalAllocationBytes: sample.TotalAllocationBytes,
		Operations: int64(sample.Operations),
		ThroughputOpsPerSec: float64(sample.Operations)*1e9/float64(sample.ElapsedNS),
	},nil
}
