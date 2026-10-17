package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Cross-Case temporal pivots must be based on real operator-selected disk
// snapshots and the existing independently authorized owner HTTP traversal.
func windowPivotSnapshot(t *testing.T, uncertain bool) string {
	t.Helper()
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			if uncertain {
				data.Entries[i].Value.SourceTime = "undetermined"
				data.Entries[i].Value.Timezone = "unknown"
			} else {
				data.Entries[i].Value.SourceTime = "2026-10-01T10:02:05Z"
			}
		}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "window-owner.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLIFiveMinutePivotAcrossAuthorizedDiskCases(t *testing.T) {
	path := windowPivotSnapshot(t, false)
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-compare-case=case-a", "-with-case=case-c",
		"-compare-events", "-pivot-window=2026-10-01T10:00:00Z",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got localCaseEventSequence
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Format != "cmdr-local-cross-case-window-events-v1" || !got.Synthetic ||
		got.WindowStartUTC != "2026-10-01T10:00:00Z" ||
		got.WindowEndUTC != "2026-10-01T10:05:00Z" ||
		got.First.VisibleEntries != 207 || got.Second.VisibleEntries != 1 ||
		got.MatchedFirst != 206 || got.MatchedSecond != 1 ||
		got.CertainEntries != 207 || got.UncertainEntries != 0 ||
		len(got.Events) != 207 {
		t.Fatalf("incorrect authorized temporal intersection: %+v", got)
	}
	foundDiskEvent := false
	for _, row := range got.Events {
		if row.Entry.TimeUncertain || row.Entry.ID == "event-new" ||
			row.Entry.ID == "restricted" ||
			row.CaseID != "case-a" && row.CaseID != "case-c" {
			t.Fatalf("out-of-window or unpermitted entry: %+v", row)
		}
		if row.CaseID == "case-c" {
			foundDiskEvent = row.Entry.ID == "event-c" &&
				row.Entry.SourceRef == "disk-source-c" &&
				row.Entry.Provenance == "snapshot-only-provenance" &&
				row.Entry.Kind == "inferred"
		}
	}
	if !foundDiskEvent {
		t.Fatal("second disk-only Case source absent from pivot")
	}
	for _, forbidden := range []string{
		"demo-hidden-source", "demo-tenant-b",
		"SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT", "Bearer ",
	} {
		if strings.Contains(out.String(), forbidden) {
			t.Fatalf("forbidden metadata in temporal receipt: %q", forbidden)
		}
	}
}

func TestFiveMinutePivotFailsClosedWhenOneCaseCannotContribute(t *testing.T) {
	for _, tc := range []struct{ snapshot, secondCase, window string }{
		{twoCaseComparisonSnapshot(t), "case-c", "2026-10-01T10:00:00Z"},
		{windowPivotSnapshot(t, true), "case-c", "2026-10-01T10:00:00Z"},
		{windowPivotSnapshot(t, false), "case-b", "2026-10-01T10:00:00Z"},
		{windowPivotSnapshot(t, false), "missing", "2026-10-01T10:00:00Z"},
		{windowPivotSnapshot(t, false), "case-c", "2026-10-01T10:05:00Z"},
		{windowPivotSnapshot(t, false), "case-c", "2026-10-01T10:00:00+00:00"},
		{windowPivotSnapshot(t, false), "case-c", "2026-10-01T10:03:00Z"},
		{windowPivotSnapshot(t, false), "case-c", "invalid"},
	} {
		var out bytes.Buffer
		err := writeDemoCaseEventSequenceWithWindow(context.Background(), "analyst-a",
			tc.snapshot, "case-a", tc.secondCase, tc.window, &out)
		if err == nil || out.Len() != 0 {
			t.Fatalf("invalid or unauthorized window emitted data: err=%v out=%q", err, out.String())
		}
	}
}

func TestWindowPivotRejectsModeConflictCancellationAndWriterFailure(t *testing.T) {
	path := windowPivotSnapshot(t, false)
	for _, flags := range [][]string{
		{"-pivot-window=2026-10-01T10:00:00Z"},
		{"-compare-case=case-a", "-with-case=case-c", "-pivot-window=2026-10-01T10:00:00Z"},
		{"-compare-case=case-a", "-with-case=case-c", "-compare-events",
			"-pivot-window=2026-10-01T10:00:00Z", "-report-case=case-a"},
		{"-compare-case=case-a", "-with-case=case-c", "-compare-events",
			"-pivot-window=2026-10-01T10:00:00Z", "-source-case=case-a", "-source-ref=demo-evidence-a"},
	} {
		var out bytes.Buffer
		args := append([]string{"-principal=analyst-a", "-snapshot=" + path}, flags...)
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid mode emitted partial report: err=%v out=%q", err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoCaseEventSequenceWithWindow(ctx, "analyst-a", path,
		"case-a", "case-c", "2026-10-01T10:00:00Z", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled investigation leaked rows: %v %q", err, out.String())
	}
	if err := writeDemoCaseEventSequenceWithWindow(context.Background(), "analyst-a", path,
		"case-a", "case-c", "2026-10-01T10:00:00Z", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("output failure not propagated: %v", err)
	}
}
