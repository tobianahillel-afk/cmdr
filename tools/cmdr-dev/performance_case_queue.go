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

type caseQueueProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareCaseQueueProbe(root string) (caseQueueProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/case-queue", false)
	if err != nil {
		return caseQueueProbe{}, func() {}, fmt.Errorf("Case Queue runtime unavailable: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-case-queue-perf-")
	if err != nil {
		return caseQueueProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- created under verified repository root.
	}
	binaryName := "case-queue-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binaryPath, "./cmd/perf-probe"); err != nil {
		cleanup()
		return caseQueueProbe{}, func() {}, fmt.Errorf("build Case Queue probe: %w", err)
	}
	return caseQueueProbe{moduleRoot: moduleRoot, binaryPath: binaryPath}, cleanup, nil
}

func (probe caseQueueProbe) run(iterations int, operation, mode string, memoryMiB int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 1000000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("invalid Case Queue iteration count")
	}
	if operation != "projection" && operation != "saved-view" && operation != "handoff" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("invalid Case Queue operation")
	}
	if mode != "latency" && mode != "resource" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("invalid Case Queue probe mode")
	}
	args := []string{"-iterations", strconv.Itoa(iterations), "-mode", mode, "-operation", operation}
	if mode == "resource" {
		args = append(args, "-memory-limit-mib", strconv.Itoa(memoryMiB))
	}
	output, err := runPerformanceProcess(probe.moduleRoot, probe.binaryPath, args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var result pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
		return result, fmt.Errorf("decode Case Queue probe: %w", err)
	}
	if result.Operation != operation || result.Mode != mode || result.Operations != iterations ||
		result.ElapsedNS <= 0 || result.NSPerOperation <= 0 || math.IsNaN(result.NSPerOperation) ||
		math.IsInf(result.NSPerOperation, 0) {
		return result, fmt.Errorf("invalid Case Queue performance observation")
	}
	return result, nil
}
