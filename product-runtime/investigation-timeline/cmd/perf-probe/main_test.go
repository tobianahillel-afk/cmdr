package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestProbeRejectsInvalidParameters(t *testing.T) {
	for _, tc := range []struct{ iterations int; mode string; memory int }{
		{0, "latency", 0}, {-1, "latency", 0}, {maxIterations + 1, "latency", 0},
		{1, "unsupported", 0}, {1, "latency", 128}, {1, "resource", 0},
		{1, "resource", 63}, {1, "resource", 2049},
	} {
		if _, err := measure(tc.iterations, tc.mode, tc.memory); err == nil {
			t.Fatalf("accepted invalid probe %+v", tc)
		}
	}
}

func TestProbeIsBoundToRealProjection(t *testing.T) {
	for _, mode := range []string{"latency", "resource"} {
		memory := 0
		if mode == "resource" {
			memory = 128
		}
		got, err := measure(2, mode, memory)
		if err != nil {
			t.Fatal(err)
		}
		if got.Operation != "projection" || got.Mode != mode || got.Operations != 2 ||
			got.ElapsedNS <= 0 || got.NSPerOperation <= 0 ||
			got.PeakHeapBytes == 0 || got.TotalAllocationBytes == 0 {
			t.Fatalf("invalid measured observation %+v", got)
		}
	}
}

func TestProbeCLISafeOutput(t *testing.T) {
	var out bytes.Buffer
	if err := runCLI([]string{"-iterations=1"}, &out); err != nil {
		t.Fatal(err)
	}
	var got observation
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Operation != "projection" || got.Operations != 1 ||
		strings.Contains(out.String(), "tenant-benchmark") {
		t.Fatalf("unsafe or invalid observation: %s", out.String())
	}
	if err := runCLI([]string{"-mode=bad"}, &out); err == nil {
		t.Fatal("unsupported CLI mode was accepted")
	}
}
