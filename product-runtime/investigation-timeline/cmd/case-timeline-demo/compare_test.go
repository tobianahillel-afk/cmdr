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

func twoCaseComparisonSnapshot(t *testing.T) string {
	t.Helper()
	data := demoData()
	data.Cases = append(data.Cases, timeline.MemoryTimelineCase{
		ID: "case-c", Tenant: "demo-tenant-a", Environment: "lab",
		Readers: []string{"analyst-a"}, EntryIDs: []string{"event-c"},
	})
	data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
		CaseID: "case-c",
		Value: timeline.Entry{
			ID: "event-c", Tenant: "demo-tenant-a", Kind: "inferred",
			SourceTime: "2026-10-01T10:30:00Z", Timezone: "+00:00",
			Version: 1, SourceRef: "disk-source-c",
			Provenance: "snapshot-only-provenance",
		},
		Readers: []string{"analyst-a"},
	})
	data.Sources = append(data.Sources, timeline.MemoryTimelineSource{
		Identity: timeline.SourceOwnerIdentity{
			Reference: "disk-source-c", Kind: "evidence",
			ObjectID: "object-c", Tenant: "demo-tenant-a",
			Environment: "lab", CaseID: "case-c",
		},
		Readers: []string{"analyst-a"},
		Preview: "SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT",
	})
	// Ensure first Case spans two authorized pages; hidden source remains
	// catalogued but must not enter the receipt.
	for i := 0; i < 205; i++ {
		id := fmt.Sprintf("extra-%03d", i)
		kind := "observed"
		if i%2 != 0 {
			kind = "inferred"
		}
		data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, id)
		data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
			CaseID: "case-a",
			Value: timeline.Entry{
				ID: id, Tenant: "demo-tenant-a", Kind: kind,
				SourceTime: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC).
					Add(time.Duration(i) * time.Second).Format(time.RFC3339),
				Timezone: "+00:00", Version: 1,
				SourceRef: "demo-evidence-a", Provenance: "authorized-on-disk",
			},
			Readers: []string{"analyst-a"},
		})
	}
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "two-case-compare.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLICompareTwoRealDiskCasesAcrossPagination(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-compare-case=case-a", "-with-case=case-c",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var report localCaseComparison
	if err := json.Unmarshal(output.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-case-comparison-v1" || !report.Synthetic ||
		report.First.CaseID != "case-a" || report.Second.CaseID != "case-c" ||
		report.First.VisibleEntries != 207 || report.Second.VisibleEntries != 1 ||
		report.First.Kinds["observed"] != 104 || report.First.Kinds["inferred"] != 103 ||
		report.Second.Kinds["inferred"] != 1 || report.VisibleEntryDelta != 206 ||
		report.ObservedDelta != 104 || report.InferredDelta != 102 ||
		report.CertainTimeOverlap == nil ||
		report.CertainTimeOverlap.Start != "2026-10-01T10:30:00Z" ||
		report.CertainTimeOverlap.End != "2026-10-01T10:30:00Z" {
		t.Fatalf("incomplete cross-Case evidence: %+v", report)
	}
	for _, secret := range []string{
		"SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT", "demo-hidden-source",
		"demo-secret", "demo-tenant-b", "event-c", "event-new",
		"event-old", "extra-000", "Bearer ", testBearer, "http://",
	} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("comparison leaked raw/denied data %q", secret)
		}
	}
	if !strings.Contains(output.String(), "disk-source-c") || report.First.DistinctSources != 2 {
		// The receipt intentionally exports source references, not provenance.
		t.Fatal("source aggregates from disk not observed")
	}
}

func TestCLIComparisonFailsClosedBeforeAnyOutput(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-b"},
		{"-principal=analyst-b", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=missing"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a", "-with-case=case-a"},
		{"-principal=analyst-a", "-snapshot=" + path, "-compare-case=case-a"},
		{"-principal=analyst-a", "-snapshot=" + path, "-with-case=case-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-report-case=case-a", "-compare-case=case-a", "-with-case=case-c"},
		{"-principal=analyst-a", "-snapshot=relative.json", "-compare-case=case-a", "-with-case=case-c"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid comparison emitted data: %v %s", err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoCaseComparison(ctx, "analyst-a", path, "case-a", "case-c", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled comparison leaked data: %v %q", err, out.String())
	}
	if err := writeDemoCaseComparison(context.Background(), "analyst-a", path, "case-a", "case-c", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("output failure not propagated: %v", err)
	}
}

func TestCLIComparisonWithoutCertainOverlap(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	var out bytes.Buffer
	if err := writeDemoCaseComparison(context.Background(), "analyst-a", path, "case-c", "case-a", &out); err != nil {
		t.Fatal(err)
	}
	var got localCaseComparison
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.VisibleEntryDelta != -206 || got.ObservedDelta != -104 ||
		got.CertainTimeOverlap == nil {
		t.Fatalf("comparison order or overlap altered: %+v", got)
	}
}
