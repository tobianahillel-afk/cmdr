package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestSingleCaseJSONLFromRealOwnerDiskAcrossPages(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path, "-export-case=case-a",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(bytes.TrimSpace(out.Bytes()), []byte("\n"))
	if len(lines) != 208 {
		t.Fatalf("wanted metadata + 207 authorized Timeline events, got %d", len(lines))
	}
	var metadata caseJSONLHeader
	if err := json.Unmarshal(lines[0], &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata.RecordType != "metadata" || metadata.Format != "cmdr-local-case-timeline-jsonl-v1" ||
		!metadata.Synthetic || metadata.CaseID != "case-a" || metadata.VisibleEntries != 207 {
		t.Fatalf("invalid JSONL provenance/count: %+v", metadata)
	}
	seen := make(map[string]bool)
	var observed, inferred int
	for _, line := range lines[1:] {
		var row caseJSONLEntry
		if err := json.Unmarshal(line, &row); err != nil {
			t.Fatal(err)
		}
		if row.RecordType != "timeline-entry" || row.CaseID != "case-a" ||
			row.Entry.ID == "" || row.Entry.SourceRef == "" || row.Entry.Provenance == "" ||
			seen[row.Entry.ID] {
			t.Fatalf("invalid or duplicate authorized JSONL row: %+v", row)
		}
		seen[row.Entry.ID] = true
		switch row.Entry.Kind {
		case "observed":
			observed++
		case "inferred":
			inferred++
		default:
			t.Fatalf("unexpected classification: %+v", row)
		}
	}
	if observed == 0 || inferred == 0 || !seen["event-new"] {
		t.Fatalf("real multi-page data missing: observed=%d inferred=%d", observed, inferred)
	}
	for _, hidden := range []string{
		"demo-hidden-source", "restricted-synthetic", "SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT",
		"demo-tenant-b", "disk-source-c", "snapshot-only-provenance", "Bearer ", testBearer,
	} {
		if strings.Contains(out.String(), hidden) {
			t.Fatalf("hidden data leaked into JSONL: %s", hidden)
		}
	}
}

func TestSingleCaseJSONLFailsClosedWithoutPartialOutput(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"cross-tenant", []string{"-principal=analyst-a", "-snapshot=" + path, "-export-case=case-b"}},
		{"unknown-case", []string{"-principal=analyst-a", "-snapshot=" + path, "-export-case=unknown"}},
		{"unknown-principal", []string{"-principal=unrecognized", "-snapshot=" + path, "-export-case=case-a"}},
		{"no-file", []string{"-principal=analyst-a", "-snapshot=/nonexistent/no-data.json", "-export-case=case-a"}},
		{"conflicting-summary", []string{"-export-case=case-a", "-report-case=case-a"}},
		{"conflicting-comparison", []string{"-export-case=case-a", "-compare-case=case-a", "-with-case=case-c"}},
		{"conflicting-inspection", []string{"-export-case=case-a", "-inspect-case=case-a", "-inspect-entry=event-new"}},
		{"conflicting-source", []string{"-export-case=case-a", "-source-case=case-a", "-source-ref=demo-evidence-a"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := runDemoCLI(context.Background(), tc.args, &out); err == nil || out.Len() != 0 {
				t.Fatalf("failed request leaked partial rows: err=%v bytes=%d", err, out.Len())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := runDemoCLI(ctx, []string{
		"-principal=analyst-a", "-snapshot=" + path, "-export-case=case-a",
	}, &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled export emitted data: err=%v bytes=%d", err, out.Len())
	}
}

type jsonlFailWriter struct{}
func (jsonlFailWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestSingleCaseJSONLWriteFailurePropagates(t *testing.T) {
	err := writeDemoCaseJSONL(context.Background(), "analyst-a", "", "case-a", jsonlFailWriter{})
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("output failure was swallowed: %v", err)
	}
}
