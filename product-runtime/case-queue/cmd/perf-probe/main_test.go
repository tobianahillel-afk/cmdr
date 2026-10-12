package main

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestRejectInvalidProbeArguments(t *testing.T) {
	cases := []struct {
		count     int
		operation string
		mode      string
		memory    int
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

type rejectingWriter struct{}

func (rejectingWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated output failure")
}

func TestCLIRefusesMalformedOrUndeliverableResults(t *testing.T) {
	for _, args := range [][]string{
		{"-iterations", "not-a-number"},
		{"-iterations", "0"},
		{"-mode", "invalid"},
		{"-operation", "invalid"},
		{"unexpected-positional-argument"},
	} {
		if err := runCLI(args, io.Discard); err == nil {
			t.Fatalf("invalid command accepted: %v", args)
		}
	}
	if err := runCLI([]string{"-iterations", "1"}, rejectingWriter{}); err == nil {
		t.Fatal("output failure was ignored")
	}
}

func TestSyntheticInputRefusesUnauthorizedOperations(t *testing.T) {
	for _, operation := range []string{"projection", "saved-view", "handoff"} {
		input := selectedInput(operation)
		input.AuthorizedTenants = nil
		if _, err := runProbeWithInput(1, operation, "latency", 0, input); err == nil {
			t.Fatalf("%s: unauthorized synthetic fixture produced measurement", operation)
		}
	}
	input := selectedInput("handoff")
	input.Cases[0].Freshness = "stale"
	if _, err := runProbeWithInput(1, "handoff", "latency", 0, input); err == nil {
		t.Fatal("stale Case handoff was benchmarked as available")
	}
}

func TestMainProducesValidOutput(t *testing.T) {
	oldArgs, oldStdout := os.Args, os.Stdout
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer devnull.Close()
	defer func() {
		os.Args, os.Stdout = oldArgs, oldStdout
	}()
	os.Args = []string{"case-queue-perf", "-iterations", "1", "-operation", "projection"}
	os.Stdout = devnull
	main()
}

func TestResourceProbeSamplesHeapDuringRealWork(t *testing.T) {
	got, err := runProbe(10000, "handoff", "resource", 128)
	if err != nil {
		t.Fatal(err)
	}
	if got.Operations != 10000 || got.ElapsedNS <= 0 ||
		got.PeakHeapBytes == 0 || got.TotalAllocationBytes == 0 {
		t.Fatalf("invalid deep resource observation: %+v", got)
	}
}
