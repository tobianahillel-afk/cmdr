// A bounded local performance probe for the actual Command Work Queue runtime.
// It measures only in-process projection with one trusted synthetic owner.
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
	"sync/atomic"
	"time"

	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
)

type syntheticOwner struct{}

func (syntheticOwner) List(_ context.Context, tenant, env string) ([]queue.Item, error) {
	return []queue.Item{{Kind: queue.Incident, ID: "incident-perf", Tenant: tenant, Environment: env, Owner: "analyst", Status: "active", UpdatedAt: "2026-10-19T10:00:00Z"}}, nil
}
func (syntheticOwner) CanRead(_ context.Context, principal, tenant, env string, kind queue.Kind, id string) (bool, error) {
	return principal == "analyst" && tenant == "perf-tenant" && env == "perf-env" && kind == queue.Incident && id == "incident-perf", nil
}

const maxMeasuredOperations = 1_000_000

type observation struct {
	Operation            string  `json:"operation"`
	Mode                 string  `json:"mode"`
	Operations           int     `json:"operations"`
	ElapsedNS            int64   `json:"elapsed_ns"`
	NSPerOperation       float64 `json:"ns_per_operation"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes,omitempty"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes,omitempty"`
}

var sink queue.Result

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("command-work-queue-perf", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	n := fs.Int("iterations", 10000, "operations")
	mode := fs.String("mode", "latency", "mode")
	op := fs.String("operation", "projection", "operation")
	memoryMiB := fs.Int("memory-limit-mib", 0, "resource Go heap limit MiB")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 || *n < 1 || *n > 100000 || (*mode != "latency" && *mode != "resource") || *op != "projection" {
		return fmt.Errorf("unsupported performance operation, mode or iteration count")
	}
	if *mode == "latency" && *memoryMiB != 0 {
		return fmt.Errorf("latency mode forbids memory limit")
	}
	if *mode == "resource" {
		if *memoryMiB < 64 || *memoryMiB > 8192 {
			return fmt.Errorf("invalid resource memory limit")
		}
		old := debug.SetMemoryLimit(int64(*memoryMiB) * 1024 * 1024)
		defer debug.SetMemoryLimit(old)
	}
	req := queue.Request{Principal: "analyst", Tenant: "perf-tenant", Environment: "perf-env", View: queue.All}
	owner := syntheticOwner{}
	perform := func() error {
		result, err := queue.Project(context.Background(), owner, req)
		if err != nil {
			return err
		}
		if len(result.Rows) != 1 || result.Rows[0].ID != "incident-perf" {
			return fmt.Errorf("unexpected Command projection")
		}
		sink = result
		return nil
	}
	if err := perform(); err != nil {
		return err
	}
	var before runtime.MemStats
	var peak atomic.Uint64
	var stop, done chan struct{}
	if *mode == "resource" {
		runtime.GC()
		runtime.ReadMemStats(&before)
		peak.Store(before.HeapAlloc)
		stop, done = make(chan struct{}), make(chan struct{})
		go func() {
			defer close(done)
			tick := time.NewTicker(100 * time.Microsecond)
			defer tick.Stop()
			for {
				select {
				case <-tick.C:
					var sample runtime.MemStats
					runtime.ReadMemStats(&sample)
					for {
						old := peak.Load()
						if sample.HeapAlloc <= old || peak.CompareAndSwap(old, sample.HeapAlloc) {
							break
						}
					}
				case <-stop:
					return
				}
			}
		}()
		defer func() { close(stop); <-done }()
	}
	// Never substitute a fabricated nanosecond for a timer that has not ticked.
	// Include every additional projection in the reported operation count.
	started := time.Now()
	actual := 0
	var elapsed time.Duration
	for actual < maxMeasuredOperations {
		batch := *n
		if batch > maxMeasuredOperations-actual {
			batch = maxMeasuredOperations - actual
		}
		for i := 0; i < batch; i++ {
			if err := perform(); err != nil {
				return err
			}
		}
		actual += batch
		elapsed = time.Since(started)
		if elapsed > 0 {
			break
		}
	}
	if elapsed <= 0 {
		return fmt.Errorf("Command projection timer unresolved after %d actual operations", actual)
	}
	result := observation{Operation: *op, Mode: *mode, Operations: actual, ElapsedNS: elapsed.Nanoseconds(),
		NSPerOperation: float64(elapsed.Nanoseconds()) / float64(actual)}
	if *mode == "resource" {
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		if after.HeapAlloc > peak.Load() {
			peak.Store(after.HeapAlloc)
		}
		result.PeakHeapBytes = peak.Load() // sampled Go heap, not true RSS or absolute peak
		result.TotalAllocationBytes = after.TotalAlloc - before.TotalAlloc
	}
	return json.NewEncoder(out).Encode(result)
}
func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
