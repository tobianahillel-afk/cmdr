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

func TestCLICompleteAuthorizedCrossCaseEventSequence(t *testing.T) {
	path := twoCaseComparisonSnapshot(t) // 207 permitted rows and denied source in Case A; 1 permitted row in Case C
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-compare-case=case-a", "-with-case=case-c", "-compare-events",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var report localCaseEventSequence
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-cross-case-event-sequence-v1" || !report.Synthetic ||
		report.First.VisibleEntries != 207 || report.Second.VisibleEntries != 1 ||
		len(report.Events) != 208 || report.CertainEntries != 208 || report.UncertainEntries != 0 {
		t.Fatalf("incomplete actual disk-to-ACL-to-sequence journey: %+v", report)
	}
	foundCaseC, foundInferred := false, false
	var previous time.Time
	for _, item := range report.Events {
		entry := item.Entry
		if item.CaseID != "case-a" && item.CaseID != "case-c" ||
			entry.ID == "" || entry.SourceRef == "" || entry.Provenance == "" || entry.TimeUncertain {
			t.Fatalf("unauthorized/malformed merged event: %+v", item)
		}
		if item.CaseID == "case-c" {
			foundCaseC = entry.ID == "event-c" && entry.SourceRef == "disk-source-c" &&
				entry.Provenance == "snapshot-only-provenance" && entry.Kind == "inferred"
		}
		if item.CaseID == "case-a" && entry.Kind == "inferred" {
			foundInferred = true
		}
		instant, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil || !previous.IsZero() && instant.Before(previous) {
			t.Fatalf("event chronology not globally ordered: %+v", item)
		}
		previous = instant
	}
	if !foundCaseC || !foundInferred ||
		report.Events[len(report.Events)-1].Entry.ID != "event-new" {
		t.Fatal("did not preserve authorized Case origin, inferred state and ordering")
	}
	for _, hidden := range []string{
		"demo-hidden-source", "restricted-synthetic", "SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT",
		"demo-tenant-b", "Bearer ", testBearer,
	} {
		if strings.Contains(out.String(), hidden) {
			t.Fatalf("restricted data was exported: %s", hidden)
		}
	}
}

func TestCLICrossCaseSequenceMarksUnknownTimestampInsteadOfInventingChronology(t *testing.T) {
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "undetermined"
			data.Entries[i].Value.Timezone = "unknown"
		}
	}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "unknown.json")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := writeDemoCaseEventSequence(context.Background(), "analyst-a", path, "case-a", "case-c", &out); err != nil {
		t.Fatal(err)
	}
	var report localCaseEventSequence
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.CertainEntries != 207 || report.UncertainEntries != 1 ||
		report.Events[len(report.Events)-1].CaseID != "case-c" ||
		!report.Events[len(report.Events)-1].Entry.TimeUncertain {
		t.Fatalf("unknown clock was placed into certain chronology: %+v", report)
	}
}

func TestCrossCaseEventModeFailsClosedWithoutPartialReceipts(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-b", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=missing", "-compare-events"},
		{"-principal=analyst-b", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-c", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-a", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-report-case=case-a", "-compare-case=case-a", "-with-case=case-c", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-source-case=case-a", "-source-ref=demo-evidence-a", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=not-absolute.json", "-compare-case=case-a", "-with-case=case-c", "-compare-events"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("unauthorized/conflicting options emitted event data: err=%v out=%s", err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoCaseEventSequence(ctx, "analyst-a", path, "case-a", "case-c", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled sequence emitted data: %v %q", err, out.String())
	}
	if err := writeDemoCaseEventSequence(context.Background(), "analyst-a", path, "case-a", "case-c", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("report writer error was lost: %v", err)
	}
}

func TestCrossCaseSequenceUsesAuthorizedOwnerOnly(t *testing.T) {
	data := demoData()
	// No permitted second Case: same-tenant Case membership and source ACL
	// are necessary, and an inaccessible source never becomes an event row.
	h, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = collectAuthorizedCrossCaseSequence(context.Background(), h, testBearer, "case-a", "case-b")
	if err == nil {
		t.Fatal("foreign Case was accepted")
	}
	_, err = collectAuthorizedCrossCaseSequence(context.Background(), h, testBearer, "case-a", "case-a")
	if err == nil {
		t.Fatal("duplicate Case accepted")
	}
}
