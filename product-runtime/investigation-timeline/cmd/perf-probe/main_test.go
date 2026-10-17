package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
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

func TestSyntheticProbeRejectsTrivialOrderingOrClassification(t *testing.T) {
	input := syntheticInput()
	if len(input.Entries) != 5000 || input.Entries[0].Kind != "observed" ||
		input.Entries[4999].Kind != "inferred" ||
		input.Entries[0].SourceTime <= input.Entries[4999].SourceTime {
		t.Fatal("benchmark must contain reverse-chronological mixed-kind inputs")
	}
	projected := timeline.Project(input)
	if projected.Status != "allow" || len(projected.Entries) != len(input.Entries) {
		t.Fatal("mixed benchmark fixture was rejected")
	}
	if projected.Entries[0].ID != "entry-04999" ||
		projected.Entries[0].Kind != "inferred" ||
		projected.Entries[len(projected.Entries)-1].ID != "entry-00000" ||
		projected.Entries[len(projected.Entries)-1].Kind != "observed" {
		t.Fatalf("benchmark oracle did not verify chronology and classification: first=%+v last=%+v",
			projected.Entries[0], projected.Entries[len(projected.Entries)-1])
	}
}

type rejectingWriter struct{}

func (rejectingWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestProbeCLIFailClosedErrorPaths(t *testing.T) {
	for _, args := range [][]string{
		{"-iterations=not-a-number"},
		{"-unknown-option"},
		{"unexpected-positional-argument"},
	} {
		var buf bytes.Buffer
		if err := runCLI(args, &buf); err == nil {
			t.Fatalf("invalid CLI accepted: %v", args)
		}
	}
	if err := runCLI([]string{"-iterations=1"}, rejectingWriter{}); err == nil {
		t.Fatal("CLI silently discarded JSON writer failure")
	}
}

func TestProbeMainSuccessPath(t *testing.T) {
	originalArgs := os.Args
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Args = originalArgs
		os.Stdout = originalStdout
		_ = reader.Close()
		_ = writer.Close()
	}()
	os.Args = []string{"timeline-perf", "-iterations=1"}
	os.Stdout = writer
	main()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var result observation
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Operation != "projection" || result.Mode != "latency" || result.Operations != 1 {
		t.Fatalf("main emitted invalid observation: %+v", result)
	}
}

func TestFullRowProjectionOracleRejectsInteriorCorruption(t *testing.T) {
	input := syntheticInput()
	result := timeline.Project(input)
	if !validSyntheticProjection(input, result) {
		t.Fatal("valid 5000-row source rejected")
	}
	mutations := []struct {
		name string
		apply func(*timeline.ProjectedEntry)
	}{
		{"wrong-interior-ID", func(e *timeline.ProjectedEntry) { e.ID = "counterfeit" }},
		{"wrong-observed-inferred", func(e *timeline.ProjectedEntry) { e.Kind = "observed" }},
		{"rewritten-provenance", func(e *timeline.ProjectedEntry) { e.Provenance = "invented" }},
		{"rewritten-source", func(e *timeline.ProjectedEntry) { e.SourceRef = "invented" }},
		{"rewritten-timestamp", func(e *timeline.ProjectedEntry) { e.SourceTime = "2020-01-01T00:00:00Z" }},
		{"silently-uncertain", func(e *timeline.ProjectedEntry) { e.TimeUncertain = true }},
	}
	for _, tc := range mutations {
		t.Run(tc.name, func(t *testing.T) {
			altered := result
			altered.Entries = append([]timeline.ProjectedEntry(nil), result.Entries...)
			tc.apply(&altered.Entries[2500])
			if validSyntheticProjection(input, altered) {
				t.Fatalf("full-row oracle accepted %s", tc.name)
			}
		})
	}
}

func TestFullRowProjectionOracleRejectsInvalidStatusOrLength(t *testing.T) {
	input := syntheticInput()
	valid := timeline.Project(input)
	if !validSyntheticProjection(input, valid) {
		t.Fatal("valid projection rejected")
	}
	wrongStatus := valid
	wrongStatus.Status = "filtered"
	if validSyntheticProjection(input, wrongStatus) {
		t.Fatal("oracle accepted non-allow status")
	}
	missingRow := valid
	missingRow.Entries = valid.Entries[:len(valid.Entries)-1]
	if validSyntheticProjection(input, missingRow) {
		t.Fatal("oracle accepted truncated projection")
	}
	emptyRows := valid
	emptyRows.Entries = nil
	if validSyntheticProjection(input, emptyRows) {
		t.Fatal("oracle accepted missing projection")
	}
	modified := valid
	modified.Entries = append([]timeline.ProjectedEntry(nil), valid.Entries...)
	modified.Entries[2500].Version = 2
	if validSyntheticProjection(input, modified) {
		t.Fatal("oracle accepted rewritten interior version")
	}
	modified.Entries[2500] = valid.Entries[2500]
	modified.Entries[2500].Timezone = "unknown"
	if validSyntheticProjection(input, modified) {
		t.Fatal("oracle accepted altered interior timezone")
	}
}
