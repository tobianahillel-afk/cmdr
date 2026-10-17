package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func testOwner() timeline.MemoryTimelineData {
	return timeline.MemoryTimelineData{
		Principals: []timeline.MemoryTimelinePrincipal{
			{ID: "analyst-a", Tenant: "tenant-a", Environment: "lab"},
			{ID: "analyst-b", Tenant: "tenant-b", Environment: "lab"},
		},
		Cases: []timeline.MemoryTimelineCase{
			{ID: "case-a", Tenant: "tenant-a", Environment: "lab",
				Readers: []string{"analyst-a"},
				EntryIDs: []string{"allowed-late", "restricted", "evidence-entry", "allowed-early"}},
			{ID: "case-b", Tenant: "tenant-b", Environment: "lab",
				Readers: []string{"analyst-b"}, EntryIDs: []string{"foreign"}},
		},
		Entries: []timeline.MemoryTimelineEntry{
			{CaseID: "case-a", Readers: []string{"analyst-a"}, Value: timeline.Entry{
				ID: "allowed-early", Tenant: "tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00",
				SourceRef: "telemetry-early", Provenance: "sensor-dns", Version: 1}},
			{CaseID: "case-a", Readers: []string{"analyst-a"}, Value: timeline.Entry{
				ID: "allowed-late", Tenant: "tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T10:10:00Z", Timezone: "+00:00",
				SourceRef: "telemetry-late", Provenance: "sensor-auth", Version: 1}},
			{CaseID: "case-a", Readers: []string{"analyst-a"}, Value: timeline.Entry{
				ID: "restricted", Tenant: "tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T10:02:00Z", Timezone: "+00:00",
				SourceRef: "private-secret", Provenance: "sensor-secret", Version: 1}},
			{CaseID: "case-a", Readers: []string{"analyst-a"}, Value: timeline.Entry{
				ID: "evidence-entry", Tenant: "tenant-a", Kind: "inferred",
				SourceTime: "2026-10-01T10:03:00Z", Timezone: "+00:00",
				SourceRef: "evidence-1", Provenance: "sensor-evidence", Version: 1}},
			{CaseID: "case-b", Readers: []string{"analyst-b"}, Value: timeline.Entry{
				ID: "foreign", Tenant: "tenant-b", Kind: "observed",
				SourceTime: "2026-10-01T10:01:00Z", Timezone: "+00:00",
				SourceRef: "private-foreign", Provenance: "sensor-foreign", Version: 1}},
		},
		Sources: []timeline.MemoryTimelineSource{
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "telemetry-early", Kind: "telemetry-event", ObjectID: "event-early",
				Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"},
				Readers: []string{"analyst-a"}, Preview: "private permitted telemetry payload"},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "telemetry-late", Kind: "telemetry-event", ObjectID: "event-late",
				Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"},
				Readers: []string{"analyst-a"}, Preview: "private permitted late payload"},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "private-secret", Kind: "telemetry-event", ObjectID: "event-secret",
				Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"},
				Readers: []string{}, Preview: "DO-NOT-REVEAL-TOP-SECRET"},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "evidence-1", Kind: "evidence", ObjectID: "object-evidence",
				Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"},
				Readers: []string{"analyst-a"}, Preview: "non-telemetry evidence"},
			{Identity: timeline.SourceOwnerIdentity{
				Reference: "private-foreign", Kind: "telemetry-event", ObjectID: "event-foreign",
				Tenant: "tenant-b", Environment: "lab", CaseID: "case-b"},
				Readers: []string{"analyst-b"}, Preview: "foreign tenant private payload"},
		},
	}
}

func testOptions() searchOptions {
	return searchOptions{
		Principal: "analyst-a", CaseID: "case-a",
		Query: "sensor", From: "2026-10-01T10:00:00Z", To: "2026-10-01T10:15:00Z",
	}
}

func TestActualAuthorizedCaseSearchJob(t *testing.T) {
	report, err := runCaseSearch(context.Background(), testOwner(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-authorized-case-search-job-v1" ||
		!report.Synthetic || report.CaseID != "case-a" ||
		report.Job.State != search.JobCompleted ||
		report.Job.Provenance.TenantRef != "tenant-a" ||
		len(report.Job.ResultRefs) != 2 || len(report.Job.AuditTrail) != 3 {
		t.Fatalf("wrong integrated Event Search job: %+v", report)
	}
	if report.Job.ResultRefs[0].Ref != "allowed-early" ||
		report.Job.ResultRefs[1].Ref != "allowed-late" {
		t.Fatalf("expected two independently authorized telemetry results: %+v", report.Job.ResultRefs)
	}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"private permitted", "DO-NOT-REVEAL",
		"private-secret", "private-foreign", "evidence-entry", "tenant-b"} {
		if strings.Contains(string(raw), forbidden) {
			t.Fatalf("source/tenant data escaped the authorized search: %s", forbidden)
		}
	}
}

func TestAuthorizedCaseSearchBoundsAndLeastPrivilege(t *testing.T) {
	for _, tc := range []struct {
		name string
		change func(*searchOptions)
	}{
		{"foreign Case", func(o *searchOptions) { o.CaseID = "case-b" }},
		{"foreign principal", func(o *searchOptions) { o.Principal = "analyst-b" }},
		{"empty query", func(o *searchOptions) { o.Query = "" }},
		{"unbounded query", func(o *searchOptions) { o.Query = strings.Repeat("x", 97) }},
		{"bad interval", func(o *searchOptions) { o.To = o.From }},
		{"bad date", func(o *searchOptions) { o.From = "not-a-time" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := testOptions()
			tc.change(&opts)
			got, err := runCaseSearch(context.Background(), testOwner(), opts)
			if err == nil || got.Job.ID != "" || len(got.Job.ResultRefs) != 0 {
				t.Fatalf("invalid search emitted job data: %v %+v", err, got)
			}
		})
	}
	// Search excludes records outside the bounded time window.
	opts := testOptions()
	opts.To = "2026-10-01T10:05:00Z"
	got, err := runCaseSearch(context.Background(), testOwner(), opts)
	if err != nil || len(got.Job.ResultRefs) != 1 ||
		got.Job.ResultRefs[0].Ref != "allowed-early" {
		t.Fatalf("time filter leaked unrelated result: %v %+v", err, got.Job.ResultRefs)
	}
	// A fully authorized but unmatched query completes without fabricating refs.
	opts.Query = "nothing-will-match"
	got, err = runCaseSearch(context.Background(), testOwner(), opts)
	if err != nil || got.Job.State != search.JobCompleted || len(got.Job.ResultRefs) != 0 {
		t.Fatalf("empty match should be complete: %v %+v", err, got.Job)
	}
}

func TestCaseSearchJobCLIReadsActualSnapshotAndWritesNoPartialDenials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "owner.json")
	encoded, err := json.Marshal(testOwner())
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(path, encoded, 0600); err != nil { t.Fatal(err) }
	args := []string{
		"-snapshot", path, "-principal", "analyst-a", "-case", "case-a",
		"-query", "sensor-dns", "-from", "2026-10-01T09:00:00Z", "-to", "2026-10-01T11:00:00Z",
	}
	var out bytes.Buffer
	if err := runCLI(context.Background(), args, &out); err != nil { t.Fatal(err) }
	var got caseSearchReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Job.State != search.JobCompleted || len(got.Job.ResultRefs) != 1 ||
		got.Job.ResultRefs[0].Ref != "allowed-early" {
		t.Fatalf("offline disk-to-job search failed: %+v", got)
	}
	args[6] = "case-b" // caller cannot select another tenant's Case
	out.Reset()
	if err := runCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
		t.Fatalf("foreign Case leaked partial receipt: err=%v output=%s", err, out.String())
	}
	// A cancelled request cannot expose any search result.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	out.Reset()
	if err := runCLI(ctx, args, &out); err == nil || out.Len() != 0 {
		t.Fatalf("cancelled search emitted partial receipt: %v %s", err, out.String())
	}
}
