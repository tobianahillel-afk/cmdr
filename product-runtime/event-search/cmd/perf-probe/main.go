// Command perf-probe executes bounded Event Search validation or Search Job
// orchestration using synthetic tenant-scoped inputs and reports aggregate local performance.
package main

import (
	"context"
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

	eventsearch "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

const maxIterations = 10_000_000

var (
	validationSink    eventsearch.Result
	orchestrationSink eventsearch.SearchJob
)

type observation struct {
	Operation            string  `json:"operation"`
	Mode                 string  `json:"mode"`
	Operations           int     `json:"operations"`
	ElapsedNS            int64   `json:"elapsed_ns"`
	NSPerOperation       float64 `json:"ns_per_operation"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
}

type completedBackend struct{}

func (completedBackend) Execute(_ context.Context, req eventsearch.BackendRequest) (eventsearch.BackendResult, error) {
	return eventsearch.BackendResult{
		TenantRef:        req.TenantRef,
		CompletedSources: append([]string(nil), req.Sources...),
	}, nil
}

func main() {
	if err := runCLI(os.Args[1:], os.Stdout); err != nil {
		exitf("%v", err)
	}
}

func runCLI(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("perf-probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	iterations := fs.Int("iterations", 20_000, "operations to execute")
	mode := fs.String("mode", "latency", "probe mode: latency or resource")
	operation := fs.String("operation", "validation", "operation: validation or orchestration")
	memoryLimitMiB := fs.Int("memory-limit-mib", 0, "optional Go memory limit for resource mode")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	result, err := runProbeOperation(*iterations, *mode, *memoryLimitMiB, *operation)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(result)
}

func runProbe(iterations int, mode string, memoryLimitMiB int) (observation, error) {
	return runProbeOperation(iterations, mode, memoryLimitMiB, "validation")
}

func runProbeOperation(iterations int, mode string, memoryLimitMiB int, operation string) (observation, error) {
	if iterations < 1 || iterations > maxIterations {
		return observation{}, fmt.Errorf("iterations must be within 1..%d", maxIterations)
	}
	if mode != "latency" && mode != "resource" {
		return observation{}, fmt.Errorf("mode must be latency or resource")
	}
	if operation != "validation" && operation != "orchestration" {
		return observation{}, fmt.Errorf("operation must be validation or orchestration")
	}
	if mode == "latency" && memoryLimitMiB != 0 {
		return observation{}, fmt.Errorf("latency mode does not accept memory-limit-mib")
	}
	if mode == "resource" {
		if memoryLimitMiB < 64 || memoryLimitMiB > 8192 {
			return observation{}, fmt.Errorf("resource memory-limit-mib must be within 64..8192")
		}
		old := debug.SetMemoryLimit(int64(memoryLimitMiB) * 1024 * 1024)
		defer debug.SetMemoryLimit(old)
	}

	input := eventsearch.Input{
		TenantRef:         "tenant-a",
		AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef:    "env-prod",
		TimeStart:         "2026-09-24T10:00:00Z",
		TimeEnd:           "2026-09-24T11:00:00Z",
		Sources:           []string{"source-a"},
		AuthorizedSources: []string{"source-a"},
		Permissions: []string{
			eventsearch.PermissionSearchExecute,
			eventsearch.PermissionQueryRead,
			eventsearch.PermissionSearchManage,
			eventsearch.PermissionTelemetryRead,
		},
		QueryPresent:  true,
		QueryValid:    true,
		CorrelationID: "corr-perf-001",
		QueryVersion:  "query-v1",
	}
	if got := eventsearch.Validate(input); !got.Allowed || !got.ProtectedDataVisible ||
		got.Envelope.TenantRef != "tenant-a" || got.Envelope.CorrelationID != "corr-perf-001" {
		return observation{}, fmt.Errorf("synthetic Event Search validation preflight failed")
	}

	var orchestrator *eventsearch.Orchestrator
	var baseJob eventsearch.SearchJob
	audit := eventsearch.AuditContext{
		ActorRef:  "perf-probe",
		Rationale: "deterministic Event Search orchestration benchmark",
	}
	if operation == "orchestration" {
		var err error
		orchestrator, err = eventsearch.NewOrchestrator(completedBackend{})
		if err != nil {
			return observation{}, fmt.Errorf("construct orchestration probe: %w", err)
		}
		baseJob, err = orchestrator.Create(input, audit)
		if err != nil {
			return observation{}, fmt.Errorf("create orchestration probe job: %w", err)
		}
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
					updatePeak(&peak)
				case <-stop:
					return
				}
			}
		}()
	}

	ctx := context.Background()
	started := time.Now()
	for i := 0; i < iterations; i++ {
		switch operation {
		case "validation":
			validationSink = eventsearch.Validate(input)
		case "orchestration":
			var err error
			orchestrationSink, err = orchestrator.Execute(ctx, baseJob, audit)
			if err != nil {
				if stop != nil {
					close(stop)
					sampler.Wait()
				}
				return observation{}, fmt.Errorf("orchestration iteration failed: %w", err)
			}
		}
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
	switch operation {
	case "validation":
		if !validationSink.Allowed {
			return observation{}, fmt.Errorf("validation sink ended in unexpected state")
		}
	case "orchestration":
		if orchestrationSink.State != eventsearch.JobCompleted ||
			len(orchestrationSink.ResultRefs) != 0 || len(orchestrationSink.FailedSources) != 0 {
			return observation{}, fmt.Errorf("orchestration sink ended in unexpected state")
		}
	}
	if elapsed <= 0 {
		return observation{}, fmt.Errorf("non-positive elapsed duration")
	}

	return observation{
		Operation:            operation,
		Mode:                 mode,
		Operations:           iterations,
		ElapsedNS:            elapsed.Nanoseconds(),
		NSPerOperation:       float64(elapsed.Nanoseconds()) / float64(iterations),
		PeakHeapBytes:        peak.Load(),
		TotalAllocationBytes: after.TotalAlloc - before.TotalAlloc,
	}, nil
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
