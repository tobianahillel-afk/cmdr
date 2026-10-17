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

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
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

func TestTemporalWindowHonorsOffsetsAndExclusiveEndWithTwoCaseOwner(t *testing.T) {
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil { t.Fatal(err) }
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "2026-10-01T11:02:05+01:00"
			data.Entries[i].Value.Timezone = "+01:00"
		}
	}
	for i := range data.Cases {
		if data.Cases[i].ID == "case-c" {
			data.Cases[i].EntryIDs = append(data.Cases[i].EntryIDs, "case-c-end-exclusive")
		}
	}
	data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
		CaseID: "case-c",
		Value: timeline.Entry{
			ID: "case-c-end-exclusive", Tenant: "demo-tenant-a", Kind: "observed",
			SourceTime: "2026-10-01T10:05:00Z", Timezone: "+00:00", Version: 1,
			SourceRef: "disk-source-c", Provenance: "boundary-from-disk",
		},
		Readers: []string{"analyst-a"},
	})
	blob, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "boundary-owner.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := writeDemoCaseWindowSequence(context.Background(), "analyst-a", path, "case-a", "case-c", "2026-10-01T10:00:00Z", &out); err != nil { t.Fatal(err) }
	var report localCaseWindowSequence
	if err := json.Unmarshal(out.Bytes(), &report); err != nil { t.Fatal(err) }
	if report.FirstEvents != 206 || report.SecondEvents != 1 || len(report.Events) != 207 {
		t.Fatalf("unexpected boundary receipt counts: %+v", report)
	}
	if strings.Contains(out.String(), "case-c-end-exclusive") || strings.Contains(out.String(), "boundary-from-disk") {
		t.Fatal("exclusive upper-bound leaked an out-of-window entry")
	}
	found := false
	for _, item := range report.Events {
		if item.CaseID == "case-c" {
			found = item.Entry.ID == "event-c" &&
				item.Entry.SourceTime == "2026-10-01T11:02:05+01:00"
		}
	}
	if !found {
		t.Fatal("original owner timezone/offset was lost")
	}
}

func TestTemporalWindowFailsClosedWhenSecondSourceOwnerDenied(t *testing.T) {
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil { t.Fatal(err) }
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "2026-10-01T10:02:05Z"
		}
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.CaseID == "case-c" {
			data.Sources[i].Readers = nil
		}
	}
	blob, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "denied-source.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := writeDemoCaseWindowSequence(context.Background(), "analyst-a", path, "case-a", "case-c", "2026-10-01T10:00:00Z", &out); err == nil || out.Len() != 0 {
		t.Fatalf("denied owner produced partial event receipt: err=%v output=%q", err, out.String())
	}
}

func TestTemporalWindowDirectBoundaryValidation(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	for _, tc := range []struct{principal,first,second,window string}{
		{"", "case-a", "case-c", "2026-10-01T10:00:00Z"},
		{"analyst-a", "", "case-c", "2026-10-01T10:00:00Z"},
		{"analyst-a", "case-a", "case-a", "2026-10-01T10:00:00Z"},
		{"analyst-a", strings.Repeat("x",1025), "case-c", "2026-10-01T10:00:00Z"},
		{"analyst-a", "case-a", strings.Repeat("y",1025), "2026-10-01T10:00:00Z"},
		{"analyst-a", "case-a", "case-c", "2026-10-01T10:02:00Z"},
		{"analyst-a", "case-a", "case-c", "2026-10-01T10:00:00+00:00"},
		{"missing-user", "case-a", "case-c", "2026-10-01T10:00:00Z"},
		{"analyst-a", "case-a", "missing", "2026-10-01T10:00:00Z"},
	} {
		var out bytes.Buffer
		if err := writeDemoCaseWindowSequence(context.Background(), tc.principal, path, tc.first, tc.second, tc.window, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid selected authority/window emitted data: %+v output=%q", tc, out.String())
		}
	}
}

func TestAuthorizedWindowSelectionRejectsTamperedChronologyAndCaseOrigin(t *testing.T) {
	start := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	entry := func(id, when string, uncertain bool) localCaseEvent {
		return localCaseEvent{
			CaseID: id,
			Entry: timeline.ProjectedEntry{
				ID: id + "-event", Kind: "observed", SourceTime: when,
				Timezone: "+00:00", TimeUncertain: uncertain,
				SourceRef: id + "-source", Provenance: "trusted-owner", Version: 1,
			},
		}
	}
	valid := localCaseEventSequence{Events: []localCaseEvent{
		entry("case-a", "2026-10-01T10:00:00Z", false),
		entry("case-c", "2026-10-01T10:04:59Z", false),
		entry("case-c", "2026-10-01T10:01:00Z", true),
		entry("case-a", "2026-10-01T10:05:00Z", false),
	}}
	got, err := selectAuthorizedCaseWindow(context.Background(), valid, "case-a", "case-c", start)
	if err != nil || got.FirstEvents != 1 || got.SecondEvents != 1 ||
		len(got.Events) != 2 || got.StartUTC != "2026-10-01T10:00:00Z" ||
		got.EndUTC != "2026-10-01T10:05:00Z" {
		t.Fatalf("valid bounded selection failed: %+v %v", got, err)
	}
	for _, tc := range []struct {
		name string
		change func(*localCaseEventSequence)
	}{
		{"invalid-certain-source-time", func(x *localCaseEventSequence) { x.Events[0].Entry.SourceTime = "not-a-time" }},
		{"unexpected-third-case", func(x *localCaseEventSequence) { x.Events[0].CaseID = "case-b" }},
		{"missing-second-case", func(x *localCaseEventSequence) { x.Events[1].Entry.TimeUncertain = true }},
		{"both-outside-window", func(x *localCaseEventSequence) { x.Events[0].Entry.SourceTime = "2026-10-01T09:59:59Z"; x.Events[1].Entry.SourceTime = "2026-10-01T10:05:00Z" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			modified := localCaseEventSequence{Events: append([]localCaseEvent(nil), valid.Events...)}
			tc.change(&modified)
			result, err := selectAuthorizedCaseWindow(context.Background(), modified, "case-a", "case-c", start)
			if err == nil || len(result.Events) != 0 {
				t.Fatalf("tampered authorized sequence returned a partial result: %+v %v", result, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if result, err := selectAuthorizedCaseWindow(ctx, valid, "case-a", "case-c", start);
		!errors.Is(err, context.Canceled) || len(result.Events) != 0 {
		t.Fatalf("cancelled selection returned data: %+v %v", result, err)
	}
}

func TestCrossCaseWindowPivotsToReauthorizedSourceTimeline(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	const window = "2026-10-01T10:00:00Z"
	var combined bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+path,
		"-compare-case=case-a", "-with-case=case-c",
		"-compare-events", "-compare-window="+window,
	}, &combined); err != nil { t.Fatal(err) }
	var comparison localCaseWindowSequence
	if err := json.Unmarshal(combined.Bytes(), &comparison); err != nil { t.Fatal(err) }
	var source string
	for _, event := range comparison.Events {
		if event.CaseID == "case-c" && event.Entry.ID == "event-c" {
			source = event.Entry.SourceRef
		}
	}
	if source != "disk-source-c" { t.Fatalf("typed source missing in authorized cross-Case result: %q", source) }
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+path,
		"-source-case=case-c", "-source-ref="+source, "-source-window="+window,
	}, &out); err != nil { t.Fatal(err) }
	var report localSourceTimeline
	if err := json.Unmarshal(out.Bytes(), &report); err != nil { t.Fatal(err) }
	if report.Format != "cmdr-local-source-timeline-window-v1" ||
		report.CaseID != "case-c" || report.SourceRef != source ||
		report.StartUTC != window || report.EndUTC != "2026-10-01T10:05:00Z" ||
		report.MatchedEntries != 1 || report.Inferred != 1 ||
		len(report.Entries) != 1 || report.Entries[0].ID != "event-c" ||
		report.Entries[0].Provenance != "snapshot-only-provenance" {
		t.Fatalf("functional temporal Case-to-source pivot lost ownership/chronology: %+v", report)
	}
	for _, sensitive := range []string{"demo-hidden-source", "SENSITIVE-PREVIEW", "demo-tenant-b", "Bearer ", testBearer} {
		if strings.Contains(out.String(), sensitive) { t.Fatalf("pivot leaked %q", sensitive) }
	}
}

func TestSourceWindowRechecksOwnerAndFailsWithoutPartialReceipt(t *testing.T) {
	path := localSharedWindowSnapshot(t, "2026-10-01T10:02:05Z", "+00:00")
	common := []string{"-principal=analyst-a", "-snapshot="+path,
		"-source-case=case-c", "-source-ref=disk-source-c"}
	for _, delta := range [][]string{
		{"-source-window=2026-10-01T10:01:00Z"},
		{"-source-window=2026-10-01T10:00:00+00:00"},
		{"-source-window=2026-10-01T10:05:00Z"},
		{"-compare-window=2026-10-01T10:00:00Z"},
		{"-compare-events", "-source-window=2026-10-01T10:00:00Z"},
	} {
		var out bytes.Buffer
		args := append(append([]string(nil), common...), delta...)
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("bad source window accepted: %+v %v %q", delta, err, out.String())
		}
	}
	for _, tc := range []struct{principal,caseID,ref string}{
		{"analyst-b", "case-c", "disk-source-c"},
		{"analyst-a", "case-b", "demo-incident-b"},
		{"analyst-a", "case-c", "forbidden-source"},
		{"analyst-a", "case-a", "disk-source-c"},
	} {
		var out bytes.Buffer
		if err := writeDemoSourceWindowTimeline(context.Background(), tc.principal, path, tc.caseID, tc.ref,
			"2026-10-01T10:00:00Z", &out); err == nil || out.Len() != 0 {
			t.Fatalf("unauthorized source selected: %+v %v %q", tc, err, out.String())
		}
	}
	for _, badTime := range []struct{time,zone string}{
		{"undetermined", "unknown"},
		{"2026-10-01T10:05:00Z", "+00:00"},
		{"2026-10-01T09:59:59Z", "+00:00"},
	} {
		invalidPath := localSharedWindowSnapshot(t, badTime.time, badTime.zone)
		var out bytes.Buffer
		if err := writeDemoSourceWindowTimeline(context.Background(), "analyst-a", invalidPath,
			"case-c", "disk-source-c", "2026-10-01T10:00:00Z", &out); err == nil || out.Len()!=0 {
			t.Fatalf("out-of-window/uncertain source was elevated: %+v %v %q", badTime, err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoSourceWindowTimeline(ctx, "analyst-a", path, "case-c",
		"disk-source-c", "2026-10-01T10:00:00Z", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled source pivot leaked: %v %q", err, out.String())
	}
	if err := writeDemoSourceWindowTimeline(context.Background(), "analyst-a", path, "case-c",
		"disk-source-c", "2026-10-01T10:00:00Z", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("source pivot ignored writer failure: %v", err)
	}
}
