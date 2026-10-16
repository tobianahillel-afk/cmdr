// Command perf-probe measures the bounded read-only Investigation Timeline projection.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

const maxIterations = 100000

type observation struct {
	Operation            string  `json:"operation"`
	Mode                 string  `json:"mode"`
	Operations           int     `json:"operations"`
	ElapsedNS            int64   `json:"elapsed_ns"`
	NSPerOperation       float64 `json:"ns_per_operation"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
}

func main() {
	if err := runCLI(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func runCLI(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("timeline-perf", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	iterations := flags.Int("iterations", 200, "bounded projection iterations")
	mode := flags.String("mode", "latency", "latency or resource")
	memoryMiB := flags.Int("memory-limit-mib", 0, "resource mode Go heap limit")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	result, err := measure(*iterations, *mode, *memoryMiB)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(result)
}

func syntheticInput() timeline.Input {
	input := timeline.Input{
		Scope: timeline.Scope{
			Tenant: "tenant-benchmark", Environment: "test",
			CasePermission: "allow", TimelinePermission: "allow",
		},
		Entries: make([]timeline.Entry, 5000),
	}
	for i := range input.Entries {
		id := fmt.Sprintf("entry-%05d", i)
		input.Entries[i] = timeline.Entry{
			ID: id, Kind: "observed", Tenant: input.Scope.Tenant,
			SourcePermission: "allow", Timezone: "+00:00",
			SourceTime: "2026-10-01T10:00:00Z",
			SourceRef: "telemetry:" + id, Provenance: "synthetic:" + id,
			Version: 1,
		}
	}
	return input
}

func measure(iterations int, mode string, memoryMiB int) (observation, error) {
	if iterations < 1 || iterations > maxIterations {
		return observation{}, fmt.Errorf("iterations outside 1..%d", maxIterations)
	}
	if mode != "latency" && mode != "resource" {
		return observation{}, fmt.Errorf("unsupported mode")
	}
	if mode == "latency" && memoryMiB != 0 ||
		mode == "resource" && (memoryMiB < 64 || memoryMiB > 2048) {
		return observation{}, fmt.Errorf("invalid memory limit for mode")
	}
	if mode == "resource" {
		previous := debug.SetMemoryLimit(int64(memoryMiB) * 1024 * 1024)
		defer debug.SetMemoryLimit(previous)
	}
	input := syntheticInput()
	execute := func() bool {
		result := timeline.Project(input)
		return result.Status == "allow" && len(result.Entries) == len(input.Entries)
	}
	if !execute() {
		return observation{}, fmt.Errorf("Timeline synthetic preflight rejected")
	}
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	peak := before.HeapAlloc
	started := time.Now()
	for i := 0; i < iterations; i++ {
		if !execute() {
			return observation{}, fmt.Errorf("Timeline projection failed at iteration %d", i)
		}
		if mode == "resource" && i%64 == 0 {
			var sample runtime.MemStats
			runtime.ReadMemStats(&sample)
			if sample.HeapAlloc > peak {
				peak = sample.HeapAlloc
			}
		}
	}
	elapsed := time.Since(started)
	runtime.ReadMemStats(&after)
	if after.HeapAlloc > peak {
		peak = after.HeapAlloc
	}
	if elapsed <= 0 {
		return observation{}, fmt.Errorf("nonpositive elapsed time")
	}
	return observation{
		Operation: "projection", Mode: mode, Operations: iterations,
		ElapsedNS: elapsed.Nanoseconds(),
		NSPerOperation: float64(elapsed.Nanoseconds()) / float64(iterations),
		PeakHeapBytes: peak, TotalAllocationBytes: after.TotalAlloc - before.TotalAlloc,
	}, nil
}
