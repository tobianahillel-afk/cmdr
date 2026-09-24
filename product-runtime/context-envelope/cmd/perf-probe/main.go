// Command perf-probe executes the bounded context-envelope projection using
// synthetic references and reports only aggregate local performance evidence.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	contextenvelope "github.com/tobianahillel-afk/cmdr/product-runtime/context-envelope"
)

const maxIterations = 10_000_000

var projectionSink contextenvelope.Result

type observation struct {
	Mode                 string  `json:"mode"`
	Operations           int     `json:"operations"`
	ElapsedNS            int64   `json:"elapsed_ns"`
	NSPerOperation       float64 `json:"ns_per_operation"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
}

func main() {
	iterations := flag.Int("iterations", 20_000, "projection operations to execute")
	mode := flag.String("mode", "latency", "probe mode: latency or resource")
	memoryLimitMiB := flag.Int("memory-limit-mib", 0, "optional Go memory limit for resource mode")
	flag.Parse()

	if *iterations < 1 || *iterations > maxIterations {
		exitf("iterations must be within 1..%d", maxIterations)
	}
	if *mode != "latency" && *mode != "resource" {
		exitf("mode must be latency or resource")
	}
	if *mode == "latency" && *memoryLimitMiB != 0 {
		exitf("latency mode does not accept memory-limit-mib")
	}
	if *mode == "resource" {
		if *memoryLimitMiB < 64 || *memoryLimitMiB > 8192 {
			exitf("resource memory-limit-mib must be within 64..8192")
		}
		old := debug.SetMemoryLimit(int64(*memoryLimitMiB) * 1024 * 1024)
		defer debug.SetMemoryLimit(old)
	}

	input := contextenvelope.Input{
		SourceProduct:                  "command",
		DestinationProduct:             "investigate",
		SourceTenantRef:                "tenant-a",
		DestinationTenantRef:           "tenant-a",
		EnvironmentRef:                 "env-a",
		EnvironmentTenantRef:           "tenant-a",
		DestinationRequiresEnvironment: true,
		EnvironmentCompatible:          true,
		DestinationAuthorization:       contextenvelope.AuthorizationAllow,
		ContextFresh:                   true,
		ReturnOrigin:                   "/investigate",
	}
	if got := contextenvelope.Project(input); got.State != contextenvelope.StateContextValid ||
		got.TenantRef != "tenant-a" || got.EnvironmentRef != "env-a" || !got.ProtectedRefsVisible {
		exitf("synthetic projection preflight failed")
	}

	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	var peak atomic.Uint64
	peak.Store(before.HeapAlloc)

	var stop chan struct{}
	var sampler sync.WaitGroup
	if *mode == "resource" {
		stop = make(chan struct{})
		sampler.Add(1)
		go func() {
			defer sampler.Done()
			ticker := time.NewTicker(100 * time.Microsecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					updatePeak(&peak)
				case <-stop:
					return
				}
			}
		}()
	}

	started := time.Now()
	for i := 0; i < *iterations; i++ {
		projectionSink = contextenvelope.Project(input)
	}
	elapsed := time.Since(started)

	if stop != nil {
		close(stop)
		sampler.Wait()
	}
	updatePeak(&peak)
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	if after.HeapAlloc > peak.Load() {
		peak.Store(after.HeapAlloc)
	}
	if projectionSink.State != contextenvelope.StateContextValid {
		exitf("projection sink ended in unexpected state")
	}
	if elapsed <= 0 {
		exitf("non-positive elapsed duration")
	}

	result := observation{
		Mode:                 *mode,
		Operations:           *iterations,
		ElapsedNS:            elapsed.Nanoseconds(),
		NSPerOperation:       float64(elapsed.Nanoseconds()) / float64(*iterations),
		PeakHeapBytes:        peak.Load(),
		TotalAllocationBytes: after.TotalAlloc - before.TotalAlloc,
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(result); err != nil {
		exitf("encode observation: %v", err)
	}
}

func updatePeak(peak *atomic.Uint64) {
	var current runtime.MemStats
	runtime.ReadMemStats(&current)
	for {
		old := peak.Load()
		if current.HeapAlloc <= old || peak.CompareAndSwap(old, current.HeapAlloc) {
			return
		}
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}
