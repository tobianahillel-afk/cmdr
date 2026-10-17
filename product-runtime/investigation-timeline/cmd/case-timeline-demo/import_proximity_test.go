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

func TestImportJSONLToAuthorizedCrossCaseTemporalLeads(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	before, err := os.ReadFile(owner)
	if err != nil {
		t.Fatal(err)
	}
	// The second Case contains an inferred event at 10:02:02.250Z.
	// An imported observed event 125ms earlier must appear as a lead.
	events := writeJSONLEvents(t,
		"{\"id\":\"sysmon-close\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"message\":\"private process command\"}\n")
	var output bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner,
		"-import-jsonl=" + events, "-into-case=case-a",
		"-triage-import", "-with-case=case-c", "-compare-proximity",
		"-proximity-seconds=1",
	}, &output)
	if err != nil {
		t.Fatal(err)
	}
	var receipt localImportedProximity
	if err := json.Unmarshal(output.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.Format != "cmdr-local-imported-case-proximity-v1" ||
		!receipt.Synthetic || receipt.ImportedEntries != 1 ||
		receipt.Proximity.FirstCase != "case-a" ||
		receipt.Proximity.SecondCase != "case-c" ||
		receipt.Proximity.WindowSeconds != 1 ||
		receipt.Proximity.FirstCertainEntries != 208 ||
		receipt.Proximity.SecondCertainEntries != 1 {
		t.Fatalf("import-to-proximity did not traverse both Cases: %+v", receipt)
	}
	found := false
	for _, lead := range receipt.Proximity.Leads {
		if strings.HasPrefix(lead.FirstEntryID, "local-jsonl-entry-") &&
			lead.SecondEntryID == "event-c" &&
			lead.FirstKind == "observed" && lead.SecondKind == "inferred" &&
			lead.SeparationMS == 125 {
			found = true
		}
	}
	if !found {
		t.Fatalf("imported event was not correlated to independently authorized Case: %+v", receipt.Proximity.Leads)
	}
	for _, hidden := range []string{
		"private process command", "SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT",
		"demo-hidden-source", "demo-secret", "demo-tenant-b",
		"Bearer ", "principals\":", "readers\":",
	} {
		if strings.Contains(output.String(), hidden) {
			t.Fatalf("sensitive owner or preview leaked: %s", hidden)
		}
	}
	after, err := os.ReadFile(owner)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("operator owner snapshot was modified by triage")
	}
}

func TestImportJSONLProximityFailsClosed(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	good := writeJSONLEvents(t,
		"{\"id\":\"event-close\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"message\":\"sensitive\"}\n")
	bad := writeJSONLEvents(t, "{\"id\":\"evil\",\"timestamp\":\"bad\",\"message\":\"x\"}\n")
	for _, tc := range []struct {
		name, principal, other, events string
		seconds int
	}{
		{"cross-tenant", "analyst-a", "case-b", good, 1},
		{"unknown-Case", "analyst-a", "missing", good, 1},
		{"same-Case", "analyst-a", "case-a", good, 1},
		{"foreign-reader", "analyst-b", "case-c", good, 1},
		{"malformed-JSONL", "analyst-a", "case-c", bad, 1},
		{"zero-window", "analyst-a", "case-c", good, 0},
		{"oversized-window", "analyst-a", "case-c", good, 301},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			err := writeDemoImportedJSONLProximity(context.Background(), tc.principal,
				owner, "case-a", tc.other, tc.events, tc.seconds, &output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("invalid triage emitted %d bytes: %v", output.Len(), err)
			}
		})
	}
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot=" + owner, "-import-jsonl=" + good,
			"-into-case=case-a", "-triage-import", "-compare-proximity"},
		{"-principal=analyst-a", "-snapshot=" + owner, "-import-jsonl=" + good,
			"-into-case=case-a", "-triage-import", "-with-case=case-c",
			"-compare-proximity", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + owner, "-import-jsonl=" + good,
			"-into-case=case-a", "-triage-import", "-with-case=case-c",
			"-proximity-seconds=3"},
	} {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("ambiguous input emitted partial data: %v (%s)", err, output.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := writeDemoImportedJSONLProximity(ctx, "analyst-a", owner,
		"case-a", "case-c", good, 1, &output); !errors.Is(err, context.Canceled) || output.Len() != 0 {
		t.Fatalf("cancelled analysis released data: %v", err)
	}
	if err := writeDemoImportedJSONLProximity(context.Background(), "analyst-a", owner,
		"case-a", "case-c", good, 1, rejectingSnapshotWriter{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("output error was swallowed: %v", err)
	}
}

func TestImportedProximityRejectsMissingOrInvalidOperatorFiles(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	events := writeJSONLEvents(t,
		"{\"id\":\"event-close\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"message\":\"sensitive\"}\n")
	invalidOwner := writeJSONLEvents(t, "{invalid-json}\n")
	for _, tc := range []struct {
		name, owner, events string
	}{
		{"nonexistent-owner", owner + ".missing", events},
		{"nonexistent-events", owner, events + ".missing"},
		{"invalid-owner-data", invalidOwner, events},
		{"relative-owner", "relative-owner.json", events},
		{"relative-events", owner, "relative-events.jsonl"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			err := writeDemoImportedJSONLProximity(context.Background(), "analyst-a",
				tc.owner, "case-a", "case-c", tc.events, 3, &output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("unsafe local source produced %d bytes: %v", output.Len(), err)
			}
		})
	}
}
