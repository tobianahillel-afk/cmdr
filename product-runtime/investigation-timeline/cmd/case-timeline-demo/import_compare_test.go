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

func TestImportJSONLToAuthorizedTwoCaseComparison(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	before, err := os.ReadFile(owner)
	if err != nil { t.Fatal(err) }
	events := writeJSONLEvents(t, importedJSONL)
	var out bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+owner,
		"-import-jsonl="+events, "-into-case=case-a",
		"-triage-import", "-with-case=case-c",
	}, &out)
	if err != nil { t.Fatal(err) }
	var report localImportedComparison
	if err := json.Unmarshal(out.Bytes(), &report); err != nil { t.Fatal(err) }
	got := report.Comparison
	if report.Format != "cmdr-local-imported-case-comparison-v1" ||
		!report.Synthetic || report.ImportedEntries != 2 ||
		got.First.CaseID != "case-a" || got.Second.CaseID != "case-c" ||
		got.First.VisibleEntries != 209 || got.Second.VisibleEntries != 1 ||
		got.First.Kinds["observed"] != 106 || got.First.Kinds["inferred"] != 103 ||
		got.Second.Kinds["inferred"] != 1 || got.VisibleEntryDelta != 208 ||
		got.ObservedDelta != 106 || got.InferredDelta != 102 {
		t.Fatalf("import was not included in authorized comparison: %+v", report)
	}
	if got.CertainTimeOverlap == nil ||
		got.CertainTimeOverlap.Start != "2026-10-01T10:30:00Z" ||
		got.CertainTimeOverlap.End != "2026-10-01T10:30:00Z" {
		t.Fatalf("certain-time intersection lost: %+v", got.CertainTimeOverlap)
	}
	for _, secret := range []string{
		"Offline Sysmon process event", "Offline Windows log event",
		"SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT", "demo-hidden-source",
		"demo-secret", "demo-tenant-b", "operator-local-jsonl",
		"principals\":", "Bearer ", "local-jsonl-entry-",
	} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("comparison leaked owner, raw event or hidden source %q", secret)
		}
	}
	after, err := os.ReadFile(owner)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("one-shot comparison modified original owner file")
	}
}

func TestImportJSONLComparisonDeniesSecondCaseBeforeWriting(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	good := writeJSONLEvents(t, importedJSONL)
	bad := writeJSONLEvents(t, "{\"id\":\"evil\",\"timestamp\":\"bad\",\"message\":\"x\"}\n")
	for _, tc := range []struct{ name, owner, events, into, other, principal string }{
		{"foreign-second", owner, good, "case-a", "case-b", "analyst-a"},
		{"missing-second", owner, good, "case-a", "absent", "analyst-a"},
		{"same-case", owner, good, "case-a", "case-a", "analyst-a"},
		{"foreign-principal", owner, good, "case-a", "case-c", "analyst-b"},
		{"invalid-jsonl", owner, bad, "case-a", "case-c", "analyst-a"},
		{"no-owner", "", good, "case-a", "case-c", "analyst-a"},
		{"no-events", owner, "", "case-a", "case-c", "analyst-a"},
		{"same-file", owner, owner, "case-a", "case-c", "analyst-a"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := writeDemoImportedJSONLComparison(context.Background(), tc.principal,
				tc.owner, tc.into, tc.other, tc.events, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid comparison emitted output: %v (%d bytes)", err, out.Len())
			}
		})
	}
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot="+owner, "-import-jsonl="+good, "-into-case=case-a", "-with-case=case-c"},
		{"-principal=analyst-a", "-snapshot="+owner, "-import-jsonl="+good, "-into-case=case-a", "-with-case=case-c", "-triage-import", "-compare-events"},
		{"-principal=analyst-a", "-snapshot="+owner, "-import-jsonl="+good, "-into-case=case-a", "-with-case=case-c", "-triage-import", "-compare-case=case-a"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("mixed or ambiguous CLI mode emitted output: %v", err)
		}
	}
}

func TestImportJSONLComparisonCancellationAndWriterError(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	events := writeJSONLEvents(t, importedJSONL)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoImportedJSONLComparison(ctx, "analyst-a", owner, "case-a", "case-c", events, &out);
		!errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled case comparison released data: %v", err)
	}
	if err := writeDemoImportedJSONLComparison(context.Background(), "analyst-a", owner,
		"case-a", "case-c", events, rejectingSnapshotWriter{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("export writer failure not propagated: %v", err)
	}
}
