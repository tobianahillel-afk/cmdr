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
	"time"
)

// The second Case is sourced exclusively from operator-controlled disk JSON,
// while the first Case has 207 permitted entries across two HTTP pages.
func localSharedWindowSnapshot(t *testing.T, sourceTime, timezone string) string {
	t.Helper()
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = sourceTime
			data.Entries[i].Value.Timezone = timezone
		}
	}
	blob, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "shared-temporal-window.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLICrossCaseTemporalWindowActualDiskToAuthorizedHTTP(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-compare-case=case-a", "-with-case=case-c", "-compare-events",
		"-compare-window=2026-10-01T10:00:00Z",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var report localCaseWindowSequence
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-cross-case-time-window-v1" ||
		!report.Synthetic || report.FirstCase != "case-a" || report.SecondCase != "case-c" ||
		report.StartUTC != "2026-10-01T10:00:00Z" ||
		report.EndUTC != "2026-10-01T10:05:00Z" ||
		report.FirstEvents != 206 || report.SecondEvents != 1 || len(report.Events) != 207 {
		t.Fatalf("incomplete authorized temporal drilldown: %+v", report)
	}
	var previous time.Time
	foundDisk, observed, inferred := false, 0, 0
	for _, item := range report.Events {
		entry := item.Entry
		if item.CaseID != "case-a" && item.CaseID != "case-c" ||
			entry.TimeUncertain || entry.ID == "" || entry.SourceRef == "" || entry.Provenance == "" {
			t.Fatalf("malformed or unauthorized row in receipt: %+v", item)
		}
		when, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil || when.Before(time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)) ||
			!when.Before(time.Date(2026, 10, 1, 10, 5, 0, 0, time.UTC)) ||
			!previous.IsZero() && when.Before(previous) {
			t.Fatalf("window boundary or chronology corrupted: %+v", item)
		}
		previous = when
		if entry.Kind == "observed" {
			observed++
		} else if entry.Kind == "inferred" {
			inferred++
		}
		if item.CaseID == "case-c" {
			foundDisk = entry.ID == "event-c" && entry.Kind == "inferred" &&
				entry.Provenance == "snapshot-only-provenance" &&
				entry.SourceRef == "disk-source-c"
		}
	}
	if !foundDisk || observed != 103 || inferred != 104 {
		t.Fatalf("disk-origin event or type counts missing: disk=%v observed=%d inferred=%d", foundDisk, observed, inferred)
	}
	for _, restricted := range []string{
		"demo-hidden-source", "restricted-synthetic", "SENSITIVE-PREVIEW",
		"demo-tenant-b", "event-new", "Bearer ", testBearer,
	} {
		if strings.Contains(out.String(), restricted) {
			t.Fatalf("receipt disclosed restricted/out-of-window %q", restricted)
		}
	}
}

func TestCLITemporalWindowRejectsIncompleteAndUnauthorizedPaths(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	base := []string{"-principal=analyst-a", "-snapshot=" + path,
		"-compare-case=case-a", "-with-case=case-c", "-compare-events"}
	for _, change := range [][]string{
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-c", "-compare-window=2026-10-01T10:00:00Z"},
		{"-compare-window=2026-10-01T10:00:00Z"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-b", "-compare-events", "-compare-window=2026-10-01T10:00:00Z"},
		{"-principal=analyst-b", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-c", "-compare-events", "-compare-window=2026-10-01T10:00:00Z"},
		{"-principal=analyst-a", "-snapshot=" + path, "-report-case=case-a", "-compare-events", "-compare-window=2026-10-01T10:00:00Z"},
		{"-principal=analyst-a", "-snapshot=" + path, "-source-case=case-a", "-source-ref=demo-evidence-a", "-compare-window=2026-10-01T10:00:00Z"},
		{"-principal=analyst-a", "-snapshot=relative.json", "-compare-case=case-a", "-with-case=case-c", "-compare-events", "-compare-window=2026-10-01T10:00:00Z"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), change, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid selection emitted data: err=%v data=%q", err, out.String())
		}
	}
	for _, invalid := range []string{
		"not-a-window", "2026-10-01T10:00:00+00:00",
		"2026-10-01T10:01:00Z", "2026-10-01T10:00:01Z",
		"2026-10-01T10:00:00.000Z", "2026-10-01T10:00:00z",
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), append(append([]string(nil), base...), "-compare-window="+invalid), &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid window %q disclosed data: %q", invalid, out.String())
		}
	}
}

func TestCrossCaseWindowExcludesUncertainAndExclusiveBoundary(t *testing.T) {
	for _, tc := range []struct{time,zone string}{
		{"undetermined", "unknown"},
		{"2026-10-01T10:05:00Z", "+00:00"},
		{"2026-10-01T09:59:59Z", "+00:00"},
	} {
		path := localSharedWindowSnapshot(t, tc.time, tc.zone)
		var out bytes.Buffer
		if err := writeDemoCaseWindowSequence(context.Background(), "analyst-a", path, "case-a", "case-c", "2026-10-01T10:00:00Z", &out); err == nil || out.Len() != 0 {
			t.Fatalf("unshared/uncertain chronology fabricated window: %s %v %q", tc.time, err, out.String())
		}
	}
}

func TestCrossCaseWindowCancellationAndWriterFailure(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoCaseWindowSequence(ctx, "analyst-a", path, "case-a", "case-c", "2026-10-01T10:00:00Z", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled window emitted receipt: %v %q", err, out.String())
	}
	if err := writeDemoCaseWindowSequence(context.Background(), "analyst-a", path, "case-a", "case-c", "2026-10-01T10:00:00Z", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("window export ignored output error: %v", err)
	}
}
