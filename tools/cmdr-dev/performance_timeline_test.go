package main

import (
	"path/filepath"
	"testing"
)

func TestTimelineProbeRealExecutableMeasurement(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	probe, cleanup, err := prepareTimelineProbe(root)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	for _, mode := range []string{"latency", "resource"} {
		memory := 0
		if mode == "resource" {
			memory = 128
		}
		got, err := probe.run(2, mode, memory)
		if err != nil {
			t.Fatal(err)
		}
		if got.Operations != 2 || got.Operation != "projection" || got.PeakHeapBytes == 0 {
			t.Fatalf("invalid bound observation: %+v", got)
		}
	}
}
