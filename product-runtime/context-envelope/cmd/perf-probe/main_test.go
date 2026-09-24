package main

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

func TestRunCLILatencyJSON(t *testing.T) {
	var out bytes.Buffer
	if err := runCLI([]string{"-iterations", "1000", "-mode", "latency"}, &out); err != nil {
		t.Fatal(err)
	}
	var got observation
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Mode != "latency" || got.Operations != 1000 {
		t.Fatalf("unexpected observation identity: %#v", got)
	}
	if got.ElapsedNS <= 0 || got.NSPerOperation <= 0 || math.IsNaN(got.NSPerOperation) || math.IsInf(got.NSPerOperation, 0) {
		t.Fatalf("invalid latency observation: %#v", got)
	}
	if got.PeakHeapBytes == 0 {
		t.Fatalf("expected non-zero heap evidence: %#v", got)
	}
}

func TestRunProbeResource(t *testing.T) {
	got, err := runProbe(1000, "resource", 64)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != "resource" || got.Operations != 1000 {
		t.Fatalf("unexpected resource identity: %#v", got)
	}
	if got.ElapsedNS <= 0 || got.NSPerOperation <= 0 || got.PeakHeapBytes == 0 {
		t.Fatalf("invalid resource observation: %#v", got)
	}
}

func TestRunProbeRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
		mode       string
		memoryMiB  int
		want       string
	}{
		{name: "zero iterations", iterations: 0, mode: "latency", want: "iterations"},
		{name: "too many iterations", iterations: maxIterations + 1, mode: "latency", want: "iterations"},
		{name: "unknown mode", iterations: 1, mode: "other", want: "mode"},
		{name: "latency memory limit", iterations: 1, mode: "latency", memoryMiB: 64, want: "does not accept"},
		{name: "resource memory too low", iterations: 1, mode: "resource", memoryMiB: 63, want: "memory-limit-mib"},
		{name: "resource memory too high", iterations: 1, mode: "resource", memoryMiB: 8193, want: "memory-limit-mib"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := runProbe(tt.iterations, tt.mode, tt.memoryMiB)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestRunCLIRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"-unknown"}},
		{name: "positional", args: []string{"unexpected"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := runCLI(tt.args, &bytes.Buffer{}); err == nil {
				t.Fatal("expected CLI argument rejection")
			}
		})
	}
}

func TestUpdatePeakPreservesHigherValueAndRaisesLowerValue(t *testing.T) {
	var peak atomic.Uint64
	peak.Store(^uint64(0))
	updatePeak(&peak)
	if got := peak.Load(); got != ^uint64(0) {
		t.Fatalf("higher existing peak changed: %d", got)
	}

	peak.Store(0)
	updatePeak(&peak)
	if got := peak.Load(); got == 0 {
		t.Fatal("expected peak to be updated from current runtime heap")
	}
}

func TestMainSuccessPath(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Args = []string{"perf-probe", "-iterations", "100", "-mode", "latency"}
	os.Stdout = writer
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		_ = reader.Close()
		_ = writer.Close()
	}()

	main()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	var got observation
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Mode != "latency" || got.Operations != 100 || got.NSPerOperation <= 0 {
		t.Fatalf("unexpected main observation: %#v", got)
	}
}

func TestRunCLIDefaults(t *testing.T) {
	var out bytes.Buffer
	if err := runCLI(nil, &out); err != nil {
		t.Fatal(err)
	}
	var got observation
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Mode != "latency" || got.Operations != 20_000 || got.NSPerOperation <= 0 {
		t.Fatalf("unexpected default CLI observation: %#v", got)
	}
}
