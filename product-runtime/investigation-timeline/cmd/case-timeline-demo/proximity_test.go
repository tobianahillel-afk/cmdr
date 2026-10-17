package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func proximityDiskSnapshot(t *testing.T) string {
	t.Helper()
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil { t.Fatal(err) }
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "2026-10-01T10:02:02.250Z"
			data.Entries[i].Value.Timezone = "+00:00"
		}
	}
	raw, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "temporal-pairs.json")
	if err := os.WriteFile(path, raw, 0o600); err != nil { t.Fatal(err) }
	return path
}

func TestCLILocalProximityTraversesRealAuthorizedSnapshot(t *testing.T) {
	path := proximityDiskSnapshot(t)
	var output bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+path,
		"-compare-case=case-a", "-with-case=case-c",
		"-compare-proximity", "-proximity-seconds=3",
	}, &output)
	if err != nil { t.Fatal(err) }
	var report localCaseProximity
	if err := json.Unmarshal(output.Bytes(), &report); err != nil { t.Fatal(err) }
	if report.Format != "cmdr-local-cross-case-proximity-v1" || !report.Synthetic ||
		report.FirstCase != "case-a" || report.SecondCase != "case-c" ||
		report.FirstCertainEntries != 207 || report.SecondCertainEntries != 1 ||
		report.WindowSeconds != 3 || len(report.Leads) == 0 {
		t.Fatalf("missing real disk-to-ACL temporal candidates: %+v", report)
	}
	lead := report.Leads[0]
	if lead.FirstEntryID != "extra-122" || lead.SecondEntryID != "event-c" ||
		lead.SeparationMS != 250 || lead.SecondKind != "inferred" ||
		lead.FirstKind != "observed" {
		t.Fatalf("incorrect nearest timestamp/classification lead: %+v", lead)
	}
	for _, hidden := range []string{
		"SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT", "demo-hidden-source",
		"restricted-synthetic", "demo-tenant-b", "Bearer ", testBearer,
		"snapshot-only-provenance", "disk-source-c",
	} {
		if strings.Contains(output.String(), hidden) {
			t.Fatalf("owner-only contents/metadata leaked in lead output: %q", hidden)
		}
	}
}

func TestCLIProximityFailsClosedWithoutPartialReport(t *testing.T) {
	path := proximityDiskSnapshot(t)
	cases := [][]string{
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-b", "-compare-proximity"},
		{"-principal=analyst-b", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-c", "-compare-proximity"},
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=missing", "-compare-proximity"},
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-c", "-compare-proximity", "-compare-events"},
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-a", "-compare-proximity"},
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-c", "-compare-proximity", "-proximity-seconds=301"},
		{"-principal=analyst-a", "-snapshot="+path, "-compare-case=case-a", "-with-case=case-c", "-compare-proximity", "-proximity-seconds=0"},
		{"-principal=analyst-a", "-snapshot="+path, "-proximity-seconds=5"},
		{"-principal=analyst-a", "-snapshot="+path, "-report-case=case-a", "-compare-case=case-a", "-with-case=case-c", "-compare-proximity"},
	}
	for _, args := range cases {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("unauthorized/conflicting proximity request emitted output args=%v error=%v output=%q", args, err, output.String())
		}
	}
	var output bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := writeDemoCaseProximity(ctx, "analyst-a", path, "case-a", "case-c", 30, &output); !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("cancelled proximity emitted output: %v %q", err, output.String())
	}
}

func TestLocalProximityExcludesUncertainTimes(t *testing.T) {
	raw, err := loadDemoSnapshot(proximityDiskSnapshot(t))
	if err != nil { t.Fatal(err) }
	for i := range raw.Entries {
		if raw.Entries[i].CaseID == "case-c" {
			raw.Entries[i].Value.Timezone = "unknown"
			raw.Entries[i].Value.SourceTime = "undetermined"
		}
	}
	payload, err := json.Marshal(raw)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "uncertain.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := writeDemoCaseProximity(context.Background(), "analyst-a", path, "case-a", "case-c", 30, &out); err != nil { t.Fatal(err) }
	var report localCaseProximity
	if err := json.Unmarshal(out.Bytes(), &report); err != nil { t.Fatal(err) }
	if report.SecondUncertainExcluded != 1 || report.SecondCertainEntries != 0 || len(report.Leads) != 0 {
		t.Fatalf("uncertain source was assigned invented proximity: %+v", report)
	}
}

func TestNearestTemporalLeadDoesNotInferIdentityOrCausation(t *testing.T) {
	first := localCaseEventSequence{
		First: caseTriageReceipt{CaseID:"case-a"},
		Second: caseTriageReceipt{CaseID:"case-c"},
		Events: []localCaseEvent{
			{CaseID:"case-a", Entry: projectedForProximity("a", "observed", "2026-10-01T10:00:00Z")},
			{CaseID:"case-c", Entry: projectedForProximity("z", "inferred", "2026-10-01T09:59:59Z")},
			{CaseID:"case-c", Entry: projectedForProximity("b", "observed", "2026-10-01T10:00:01Z")},
		},
	}
	result, err := selectCaseProximity(context.Background(), first, 60)
	if err != nil { t.Fatal(err) }
	if len(result.Leads) != 1 || result.Leads[0].SeparationMS != 1000 ||
		result.Leads[0].SecondEntryID != "b" || result.Leads[0].SecondKind != "observed" {
		t.Fatalf("nearest equidistant lead not deterministic: %+v", result.Leads)
	}
	if result.Selection == "" || strings.Contains(result.Selection,"causality established") {
		t.Fatalf("misleading semantics: %+v", result)
	}
	if _, err := selectCaseProximity(context.Background(), first, 301); err == nil {
		t.Fatal("oversized time window accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := selectCaseProximity(ctx, first, 60); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled selection accepted: %v", err)
	}
}

func projectedForProximity(id, kind, instant string) timeline.ProjectedEntry {
	return timeline.ProjectedEntry{ID: id, Kind: kind, SourceTime: instant, Timezone: "+00:00"}
}
