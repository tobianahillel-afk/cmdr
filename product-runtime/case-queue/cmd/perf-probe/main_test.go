package main

import (
	"io"
	"testing"
)

func TestRejectInvalidProbeArguments(t *testing.T) {
	cases := []struct {
		count int
		operation string
		mode string
		memory int
	}{
		{0, "projection", "latency", 0},
		{maxIterations + 1, "projection", "latency", 0},
		{1, "invalid", "latency", 0},
		{1, "projection", "invalid", 0},
		{1, "projection", "latency", 128},
		{1, "projection", "resource", 2},
	}
	for _, tc := range cases {
		if _, err := runProbe(tc.count, tc.operation, tc.mode, tc.memory); err == nil {
			t.Fatalf("invalid probe input accepted: %+v", tc)
		}
	}
}

func TestMeasureActualRuntimeOperations(t *testing.T) {
	for _, operation := range []string{"projection", "saved-view", "handoff"} {
		for _, mode := range []string{"latency", "resource"} {
			memory := 0
			if mode == "resource" {
				memory = 128
			}
			got, err := runProbe(5, operation, mode, memory)
			if err != nil || got.Operation != operation || got.Mode != mode ||
				got.Operations != 5 || got.ElapsedNS <= 0 || got.NSPerOperation <= 0 {
				t.Fatalf("%s/%s got %+v, err=%v", operation, mode, got, err)
			}
		}
	}
	if err := runCLI([]string{"-iterations", "2", "-operation", "handoff"}, io.Discard); err != nil {
		t.Fatal(err)
	}
}
