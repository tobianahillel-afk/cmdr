package main

import (
	"strings"
	"testing"
	"time"
)

func TestCommandInvestigateDeepParsesActualGoBenchmarkResourceMetrics(t *testing.T) {
	sample := "goos: linux\nBenchmarkSameOriginRouteResource-4 250 1818 ns/op 421 B/op 4 allocs/op 902144 peakheap-B 105250 totalalloc-B\nPASS\n"
	observation, err := parseCommandInvestigateRoutingDeep(sample, 250, 500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if observation.Operations != 250 || observation.PeakHeapBytes != 902144 ||
		observation.TotalAllocationBytes != 105250 || observation.ThroughputOpsPerSec <= 0 {
		t.Fatalf("metric mismatch: %+v", observation)
	}
	for _, tc := range []struct{ name, data string }{
		{"missing", "PASS\n"},
		{"wrong-count", strings.Replace(sample, " 250 ", " 249 ", 1)},
		{"missing-peak", strings.Replace(sample, " 902144 peakheap-B", "", 1)},
		{"missing-total", strings.Replace(sample, " 105250 totalalloc-B", "", 1)},
		{"negative-total", strings.Replace(sample, "105250 totalalloc-B", "-1 totalalloc-B", 1)},
		{"not-a-number", strings.Replace(sample, "902144 peakheap-B", "nan peakheap-B", 1)},
		{"duplicate-metric", strings.Replace(sample, " 902144 peakheap-B", " 902144 peakheap-B 902144 peakheap-B", 1)},
		{"no-latency", strings.Replace(sample, " 1818 ns/op", "", 1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := parseCommandInvestigateRoutingDeep(tc.data, 250, time.Second); err == nil {
				t.Fatalf("unverified resource evidence accepted: %+v", got)
			}
		})
	}
}
