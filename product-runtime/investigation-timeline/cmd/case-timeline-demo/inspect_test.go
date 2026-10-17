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
)

func TestOfflineEntryInspectionUsesAuthorizedSourcePreview(t *testing.T) {
	var output bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=event-old",
	}, &output)
	if err != nil { t.Fatal(err) }
	var got localEntryInspection
	if err := json.Unmarshal(output.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Format != "cmdr-local-entry-inspection-v1" || !got.Synthetic ||
		got.CaseID != "case-a" || got.EntryID != "event-old" ||
		got.Kind != "evidence" || got.ObjectID != "demo-evidence-a" ||
		!strings.Contains(got.Preview, "Pièce synthétique") {
		t.Fatalf("wrong authorized inspection: %+v", got)
	}
	for _, forbidden := range []string{"demo-hidden-source", "demo-secret", "demo-tenant-b", testBearer, "Bearer token"} {
		if strings.Contains(output.String(), forbidden) {
			t.Fatalf("inspected report leaked unrelated %q", forbidden)
		}
	}
}

func TestOfflineEntryInspectionReadsOperatorDiskSource(t *testing.T) {
	data := demoData()
	data.Sources[1].Preview = "operator-disk-only-authenticated-artifact-preview"
	payload, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "owner-snapshot.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + path, "-principal=analyst-a",
		"-inspect-case=case-a", "-inspect-entry=event-old",
	}, &out); err != nil { t.Fatal(err) }
	if !strings.Contains(out.String(), data.Sources[1].Preview) {
		t.Fatalf("offline inspection did not traverse disk-snapshot source authorization: %s", out.String())
	}
}

func TestOfflineEntryInspectionDenialsProduceNoOutput(t *testing.T) {
	cases := [][]string{
		{"-inspect-case=case-a"},
		{"-inspect-entry=event-old"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-source-case=case-a", "-source-ref=demo-evidence-a"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-report-case=case-a"},
		{"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=restricted"},
		{"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=absent"},
		{"-principal=analyst-a", "-inspect-case=case-b", "-inspect-entry=event-b"},
		{"-principal=analyst-b", "-inspect-case=case-a", "-inspect-entry=event-old"},
		{"-principal=unknown", "-inspect-case=case-a", "-inspect-entry=event-old"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-snapshot=relative.json"},
	}
	for _, args := range cases {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len()!=0 {
			t.Fatalf("invalid inspection produced output: args=%v, err=%v output=%s", args, err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background()); cancel()
	var out bytes.Buffer
	if err := writeDemoEntryInspection(ctx, "analyst-a", "", "case-a", "event-old", &out); !errors.Is(err, context.Canceled) || out.Len()!=0 {
		t.Fatalf("cancelled inspection disclosed data: %v %q", err, out.String())
	}
	if err := writeDemoEntryInspection(context.Background(), "analyst-a", "", "case-a", "event-old", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("writer error lost: %v", err)
	}
}

func TestOfflineEntryInspectionDifferentTypedOwners(t *testing.T) {
	for _, tc := range []struct {
		principal, caseID, entryID, kind, objectID, phrase string
	}{
		{"analyst-a", "case-a", "event-new", "telemetry-event", "demo-event-a", "Événement réseau"},
		{"analyst-b", "case-b", "event-b", "incident", "demo-incident", "Incident synthétique"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			var out bytes.Buffer
			err := runDemoCLI(context.Background(), []string{
				"-principal=" + tc.principal, "-inspect-case=" + tc.caseID,
				"-inspect-entry=" + tc.entryID,
			}, &out)
			if err != nil {
				t.Fatal(err)
			}
			var got localEntryInspection
			if err := json.Unmarshal(out.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.Kind != tc.kind || got.ObjectID != tc.objectID ||
				got.EntryID != tc.entryID || got.CaseID != tc.caseID ||
				!strings.Contains(got.Preview, tc.phrase) {
				t.Fatalf("typed source owner was not resolved: %+v", got)
			}
		})
	}
}

func TestOfflineEntryInspectionEscapedDiskCaseAndNonDefaultAnalyst(t *testing.T) {
	data := demoData()
	const principal = "offline-analyst-only"
	const caseID = "local case & #1"
	const entryID = "evidence-id & #"
	data.Principals[0].ID = principal
	data.Cases[0].ID = caseID
	data.Cases[0].Readers = []string{principal}
	data.Cases[0].EntryIDs = []string{entryID}
	data.Entries = data.Entries[:1]
	data.Entries[0].CaseID = caseID
	data.Entries[0].Value.ID = entryID
	data.Entries[0].Value.Version = 1
	data.Entries[0].Value.Supersedes = ""
	data.Entries[0].Readers = []string{principal}
	data.Entries[0].Value.SourceRef = "offline-owner"
	data.Sources = data.Sources[:1]
	data.Sources[0].Identity.Reference = "offline-owner"
	data.Sources[0].Identity.CaseID = caseID
	data.Sources[0].Identity.ObjectID = "disk-evidence"
	data.Sources[0].Readers = []string{principal}
	data.Sources[0].Preview = "signed-by-no-one-offline-disk-preview"
	// Discard the unrelated tenant: this local file is the authority for
	// one synthetic analyst, never an enterprise identity provider.
	data.Principals = data.Principals[:1]
	data.Cases = data.Cases[:1]
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "offline-owner.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + path, "-principal=" + principal,
		"-inspect-case=" + caseID, "-inspect-entry=" + entryID,
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got localEntryInspection
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.CaseID != caseID || got.EntryID != entryID ||
		got.ObjectID != "disk-evidence" || got.Preview != data.Sources[0].Preview {
		t.Fatalf("disk owner or URL escaping failed: %+v", got)
	}
}

func TestOfflineEntryInspectionRejectsMalformedSnapshotsAndIDs(t *testing.T) {
	invalid := []struct {
		principal, snapshot, caseID, entryID string
	}{
		{"", "", "case-a", "event-old"},
		{"analyst-a", "", "", "event-old"},
		{"analyst-a", "", "case-a", ""},
		{"analyst-a", "", strings.Repeat("x", 1025), "event-old"},
		{"analyst-a", "", "case-a", strings.Repeat("x", 1025)},
		{"analyst-a", filepath.Join(t.TempDir(), "does-not-exist"), "case-a", "event-old"},
	}
	for _, tc := range invalid {
		var out bytes.Buffer
		if err := writeDemoEntryInspection(context.Background(), tc.principal,
			tc.snapshot, tc.caseID, tc.entryID, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid source inspection exposed output: %+v %q", tc, out.String())
		}
	}
}
