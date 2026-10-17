package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// The fixture's actual disk owner snapshot contains two authorized rows,
// a third denied row and a Case belonging to another tenant.
func TestImportedTriageJSONLDiskToAuthorizedCaseReport(t *testing.T) {
	owner := snapshotFixture(t)
	original, err := os.ReadFile(owner)
	if err != nil {
		t.Fatal(err)
	}
	events := writeJSONLEvents(t, importedJSONL)
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner,
		"-import-jsonl=" + events, "-into-case=forensic-case-on-disk",
		"-triage-import",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var got localImportedTriage
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Format != "cmdr-local-imported-case-triage-v1" ||
		!got.Synthetic || got.ImportedEntries != 2 ||
		got.Case.CaseID != "forensic-case-on-disk" ||
		got.Case.VisibleEntries != 4 ||
		got.Case.Kinds["observed"] < 2 ||
		got.Case.Visibility != "filtered" {
		t.Fatalf("import did not reach authorized Case triage: %+v", got)
	}
	for _, forbidden := range []string{
		"demo-hidden-source", "DO NOT DISPLAY", "demo-secret",
		"Offline Sysmon process event", "Offline Windows log event",
		"reader", "principals", "sources", "entries",
		"demo-tenant-b", "Bearer", "http://",
	} {
		if forbidden == "entries" || forbidden == "sources" {
			// The report intentionally uses visible_entries and top_sources.
			continue
		}
		if strings.Contains(output.String(), forbidden) {
			t.Fatalf("summary leaked unfiltered owner data or raw preview %q", forbidden)
		}
	}
	current, err := os.ReadFile(owner)
	if err != nil || !bytes.Equal(current, original) {
		t.Fatal("imported triage mutated disk owner")
	}
}

func TestImportedTriageFailsWithoutOutput(t *testing.T) {
	owner := snapshotFixture(t)
	valid := writeJSONLEvents(t, importedJSONL)
	bad := writeJSONLEvents(t, "{\"id\":\"evil\",\"timestamp\":\"yesterday\",\"message\":\"bad\"}\n")
	cases := []struct{ name string; args []string }{
		{"without-mode", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-triage-import"}},
		{"without-events", []string{"-snapshot=" + owner, "-triage-import"}},
		{"unknown-case", []string{"-snapshot=" + owner, "-import-jsonl=" + valid, "-into-case=other", "-triage-import"}},
		{"foreign-principal", []string{"-principal=analyst-b", "-snapshot=" + owner, "-import-jsonl=" + valid, "-into-case=forensic-case-on-disk", "-triage-import"}},
		{"malformed-event", []string{"-snapshot=" + owner, "-import-jsonl=" + bad, "-into-case=forensic-case-on-disk", "-triage-import"}},
		{"mixed-report", []string{"-snapshot=" + owner, "-import-jsonl=" + valid, "-into-case=forensic-case-on-disk", "-triage-import", "-report-case=forensic-case-on-disk"}},
		{"missing-snapshot", []string{"-import-jsonl=" + valid, "-into-case=case-a", "-triage-import"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := runDemoCLI(context.Background(), tc.args, &output); err == nil || output.Len() != 0 {
				t.Fatalf("invalid import-to-triage returned output: %v %q", err, output.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := writeDemoImportedJSONLTriage(ctx, "analyst-a", owner, "forensic-case-on-disk", valid, &output); !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("cancelled triage leaked data: %v %q", err, output.String())
	}
	if err := writeDemoImportedJSONLTriage(context.Background(), "analyst-a", owner, "forensic-case-on-disk", valid, rejectingSnapshotWriter{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("write failure not propagated: %v", err)
	}
}

func TestImportedTriageRejectsInvalidLocalFilesBeforeOutput(t *testing.T) {
	owner := snapshotFixture(t)
	events := writeJSONLEvents(t, importedJSONL)
	for _, tc := range []struct {
		name string
		snapshot string
		events string
	}{
		{name: "same-file", snapshot: owner, events: owner},
		{name: "nonexistent-owner", snapshot: owner + ".missing", events: events},
		{name: "nonexistent-events", snapshot: owner, events: events + ".missing"},
		{name: "relative-owner", snapshot: "relative-owner.json", events: events},
		{name: "relative-events", snapshot: owner, events: "relative-events.jsonl"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := writeDemoImportedJSONLTriage(context.Background(), "analyst-a", tc.snapshot,
				"forensic-case-on-disk", tc.events, &output); err == nil || output.Len() != 0 {
				t.Fatalf("unsafe local import produced output: err=%v bytes=%d", err, output.Len())
			}
		})
	}
	// The common import phase must reject crafted access-control fields
	// without creating a triage receipt or a full owner snapshot.
	forged := writeJSONLEvents(t,
		"{\"id\":\"evil\",\"timestamp\":\"2026-10-15T10:20:31Z\",\"message\":\"event\",\"readers\":[\"analyst-b\"]}\n")
	var output bytes.Buffer
	if err := writeDemoImportedJSONLTriage(context.Background(), "analyst-a", owner,
		"forensic-case-on-disk", forged, &output); err == nil || output.Len() != 0 {
		t.Fatalf("forged owner access survived normalized ingestion: %v", err)
	}
}
