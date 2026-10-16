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

type timelineProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareTimelineProbe(root string) (timelineProbe, func(), error) {
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
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binaryPath, "./cmd/perf-probe"); err != nil {
		cleanup()
		return timelineProbe{}, func() {}, fmt.Errorf("build Timeline probe: %w", err)
	}
	return timelineProbe{moduleRoot: moduleRoot, binaryPath: binaryPath}, cleanup, nil
}

func (probe timelineProbe) run(iterations int, mode string, memoryMiB int) (pilotProjectionProbeObservation, error) {
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
	output, err := runPerformanceProcess(probe.moduleRoot, probe.binaryPath, args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &observation); err != nil {
		return observation, fmt.Errorf("decode Timeline probe: %w", err)
	}
	if observation.Operation != "projection" || observation.Mode != mode ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Timeline probe observation")
	}
	return observation, nil
}
