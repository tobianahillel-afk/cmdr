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
