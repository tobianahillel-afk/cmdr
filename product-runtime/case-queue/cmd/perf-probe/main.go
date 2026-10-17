// Command perf-probe measures bounded Case Queue operations with synthetic tenant-scoped data.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	casequeue "github.com/tobianahillel-afk/cmdr/product-runtime/case-queue"
)

const maxIterations = 1000000

var projectionSink casequeue.Result
var handoffSink casequeue.HandoffResult

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
	flags := flag.NewFlagSet("case-queue-perf", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	iterations := flags.Int("iterations", 10000, "bounded operation iterations")
	operation := flags.String("operation", "projection", "projection, saved-view or handoff")
	mode := flags.String("mode", "latency", "latency or resource")
	memory := flags.Int("memory-limit-mib", 0, "resource-mode Go heap limit")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	result, err := runProbe(*iterations, *operation, *mode, *memory)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(result)
}

func fixture() casequeue.Input {
	return casequeue.Input{
		TenantRef: "tenant-perf", AuthorizedTenants: []string{"tenant-perf"},
		EnvironmentRef: "env-perf",
		Cases: []casequeue.CaseSnapshot{{
			ID: "case-perf", TenantRef: "tenant-perf", HumanID: "CASE-PERF-1",
			Title: "Synthetic Case", Status: "open", Owner: "analyst",
			UpdatedAt: "2026-09-29T10:00:00Z", AccessDecision: "allow", Freshness: "available",
		}},
	}
}

func selectedInput(operation string) casequeue.Input {
	input := fixture()
	if operation != "projection" {
		input.SavedView = &casequeue.SavedView{
			ID: "view-perf", TenantRef: "tenant-perf", Version: "v1", AccessDecision: "allow",
			Fields: []casequeue.SavedViewField{{Name: "title", AccessDecision: "allow"}},
		}
	}
	return input
}

func runProbe(iterations int, operation, mode string, memoryMiB int) (observation, error) {
	return runProbeWithInput(iterations, operation, mode, memoryMiB, selectedInput(operation))
}

func runProbeWithInput(iterations int, operation, mode string, memoryMiB int, input casequeue.Input) (observation, error) {
	if iterations < 1 || iterations > maxIterations {
		return observation{}, fmt.Errorf("iterations outside 1..%d", maxIterations)
	}
	if operation != "projection" && operation != "saved-view" && operation != "handoff" {
		return observation{}, fmt.Errorf("unsupported operation")
	}
	if mode != "latency" && mode != "resource" {
		return observation{}, fmt.Errorf("unsupported mode")
	}
	if (mode == "latency" && memoryMiB != 0) ||
		(mode == "resource" && (memoryMiB < 64 || memoryMiB > 2048)) {
		return observation{}, fmt.Errorf("invalid memory limit")
	}
	if mode == "resource" {
		previous := debug.SetMemoryLimit(int64(memoryMiB) * 1024 * 1024)
		defer debug.SetMemoryLimit(previous)
	}
	execute := func() bool {
		switch operation {
		case "projection":
			projectionSink = casequeue.Project(input)
			return projectionSink.Allowed && projectionSink.Projection != nil &&
				len(projectionSink.Projection.Rows) == 1
		case "saved-view":
			projectionSink = casequeue.ProjectWithSavedViewFallback(input)
			return projectionSink.Allowed && projectionSink.Projection != nil &&
				len(projectionSink.Projection.Rows) == 1
		default:
			handoffSink = casequeue.OpenCaseWorkspace(input, "case-perf", 12, casequeue.OriginCaseQueue)
			return handoffSink.Allowed && handoffSink.Handoff != nil &&
				handoffSink.Handoff.CaseID == "case-perf"
		}
	}
	if !execute() {
		return observation{}, fmt.Errorf("Case Queue runtime preflight rejected synthetic input")
	}
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	var peak atomic.Uint64
	peak.Store(before.HeapAlloc)

	var stop chan struct{}
	var sampler sync.WaitGroup
	if mode == "resource" {
		stop = make(chan struct{})
		sampler.Add(1)
		go func() {
			defer sampler.Done()
			ticker := time.NewTicker(100 * time.Microsecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					var stats runtime.MemStats
					runtime.ReadMemStats(&stats)
					for {
						old := peak.Load()
						if stats.HeapAlloc <= old || peak.CompareAndSwap(old, stats.HeapAlloc) {
							break
						}
					}
				case <-stop:
					return
				}
			}
		}()
	}
	started := time.Now()
	allAllowed := true
	for i := 0; i < iterations; i++ {
		if !execute() {
			allAllowed = false
			break
		}
	}
	elapsed := time.Since(started)
	if stop != nil {
		close(stop)
		sampler.Wait()
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	if after.HeapAlloc > peak.Load() {
		peak.Store(after.HeapAlloc)
	}
	if !allAllowed || elapsed <= 0 {
		return observation{}, fmt.Errorf("Case Queue runtime iteration failed")
	}
	return observation{
		Operation: operation, Mode: mode, Operations: iterations,
		ElapsedNS:            elapsed.Nanoseconds(),
		NSPerOperation:       float64(elapsed.Nanoseconds()) / float64(iterations),
		PeakHeapBytes:        peak.Load(),
		TotalAllocationBytes: after.TotalAlloc - before.TotalAlloc,
	}, nil
}
