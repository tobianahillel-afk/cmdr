package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestOfflineSourceTimelineFiltersAuthorizedSource(t *testing.T) {
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-source-case=case-a", "-source-ref=demo-evidence-a",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var report localSourceTimeline
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-source-timeline-v1" || !report.Synthetic ||
		report.CaseID != "case-a" || report.SourceRef != "demo-evidence-a" ||
		report.Visibility != "filtered" || report.MatchedEntries != 1 ||
		report.Observed != 0 || report.Inferred != 1 || len(report.Entries) != 1 ||
		report.Entries[0].ID != "event-old" || report.Entries[0].Kind != "inferred" ||
		report.Entries[0].Provenance != "synthetic-evidence" {
		t.Fatalf("incorrect authorized exact-source report: %+v", report)
	}
	for _, forbidden := range []string{
		"demo-hidden-source", "restricted", "demo-secret", "demo-tenant-b",
		"event-new", "demo-telemetry-a", "Bearer", testBearer,
	} {
		if strings.Contains(out.String(), forbidden) {
			t.Fatalf("report leaked unrelated or inaccessible %q", forbidden)
		}
	}
}

func TestOfflineSourceTimelineActualDiskMultipage(t *testing.T) {
	data := demoData()
	data.Cases[0].EntryIDs = []string{"restricted"}
	for i := 0; i < 205; i++ {
		id := fmt.Sprintf("analyst-source-%03d", i)
		kind := "observed"
		if i%2 == 1 {
			kind = "inferred"
		}
		data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, id)
		data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
			CaseID: "case-a",
			Value: timeline.Entry{
				ID: id, Tenant: "demo-tenant-a", Kind: kind,
				SourceTime: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC).
					Add(time.Duration(i) * time.Second).Format(time.RFC3339),
				Timezone: "+00:00", Version: 1, SourceRef: "demo-evidence-a",
				Provenance: "disk-authorized-origin",
			},
			Readers: []string{"analyst-a"},
		})
	}
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "operator-owner.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + path, "-principal=analyst-a",
		"-source-case=case-a", "-source-ref=demo-evidence-a",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var got localSourceTimeline
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.MatchedEntries != 205 || len(got.Entries) != 205 ||
		got.Observed != 103 || got.Inferred != 102 ||
		got.Visibility != "filtered" ||
		got.Entries[0].ID != "analyst-source-000" ||
		got.Entries[204].ID != "analyst-source-204" {
		t.Fatalf("source sequence lost pagination, order or provenance: %+v", got)
	}
	for _, e := range got.Entries {
		if e.SourceRef != "demo-evidence-a" || e.Provenance != "disk-authorized-origin" {
			t.Fatalf("unauthorized or non-disk entry: %+v", e)
		}
	}
	if strings.Contains(output.String(), "restricted") ||
		strings.Contains(output.String(), "demo-hidden-source") ||
		strings.Contains(output.String(), "demo-tenant-b") ||
		strings.Contains(output.String(), "demo-telemetry-a") {
		t.Fatal("source-focused receipt disclosed non-selected owner data")
	}
}

func TestOfflineSourceTimelineInvalidAndDeniedFailWithoutOutput(t *testing.T) {
	for _, args := range [][]string{
		{"-source-case=case-a"},
		{"-source-ref=demo-evidence-a"},
		{"-source-case=case-a", "-source-ref=demo-evidence-a", "-report-case=case-a"},
		{"-source-case=case-a", "-source-ref=demo-evidence-a", "-compare-case=case-a", "-with-case=case-b"},
		{"-principal=analyst-a", "-source-case=case-a", "-source-ref=demo-hidden-source"},
		{"-principal=analyst-a", "-source-case=case-a", "-source-ref=missing-source"},
		{"-principal=analyst-a", "-source-case=case-b", "-source-ref=demo-incident-b"},
		{"-principal=analyst-b", "-source-case=case-a", "-source-ref=demo-evidence-a"},
		{"-principal=invalid", "-source-case=case-a", "-source-ref=demo-evidence-a"},
		{"-source-case=case-a", "-source-ref=" + strings.Repeat("x", 1025)},
		{"-source-case=case-a", "-source-ref=demo-evidence-a", "-snapshot=relative.json"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid source request emitted a receipt: args=%v err=%v output=%q", args, err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := writeDemoSourceTimeline(ctx, "analyst-a", "", "case-a", "demo-evidence-a", &output); !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("cancelled source report leaked output: %v %s", err, output.String())
	}
	if err := writeDemoSourceTimeline(context.Background(), "analyst-a", "", "case-a", "demo-evidence-a", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("writer error was lost: %v", err)
	}
}

func TestSourceTimelineObserverUsesValidatedOwnerEntries(t *testing.T) {
	data := demoData()
	handler, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	_, err = collectDemoCaseTriageWithObserver(context.Background(), handler, testBearer, "case-a", func(e timeline.ProjectedEntry) error {
		ids = append(ids, e.ID)
		return nil
	})
	if err != nil || len(ids) != 2 || ids[0] != "event-old" || ids[1] != "event-new" {
		t.Fatalf("observer omitted authorized owner entries: %v %v", ids, err)
	}
	denial := errors.New("stop aggregation")
	_, err = collectDemoCaseTriageWithObserver(context.Background(), handler, testBearer, "case-a", func(timeline.ProjectedEntry) error {
		return denial
	})
	if !errors.Is(err, denial) {
		t.Fatalf("observer failure not propagated: %v", err)
	}
}
