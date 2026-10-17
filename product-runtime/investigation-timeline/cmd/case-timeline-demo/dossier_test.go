package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaseDossierUsesRealOwnerChain(t *testing.T) {
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-dossier-case=case-a", "-dossier-entry=event-old",
	}, &out)
	if err != nil { t.Fatal(err) }
	var got localCaseDossier
	if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Format != "cmdr-local-case-dossier-v1" || !got.Synthetic ||
		got.CaseID != "case-a" || got.SelectedEntry.ID != "event-old" ||
		got.SelectedEntry.Kind != "inferred" || got.Summary.VisibleEntries != 2 ||
		got.Summary.Visibility != "filtered" || got.Source.Kind != "evidence" ||
		got.Source.EntryID != "event-old" || got.Source.ObjectID != "demo-evidence-a" ||
		!strings.Contains(got.Source.Preview, "Pièce synthétique") ||
		len(got.RelatedEntries) != 1 || got.RelatedEntries[0].ID != "event-old" {
		t.Fatalf("incomplete integrated dossier: %+v", got)
	}
	for _, forbidden := range []string{"demo-hidden-source","restricted","demo-secret","demo-tenant-b","event-b",testBearer} {
		if strings.Contains(out.String(), forbidden) { t.Fatalf("unauthorized %q in dossier", forbidden) }
	}
}

func TestCaseDossierReadsOperatorDiskAndGroupsSource(t *testing.T) {
	data := demoData()
	const marker = "unique-disk-source-preview-2026"
	data.Sources[1].Preview = marker
	// One additional allowed event points to the same source; one hidden
	// source still exists and must not be exported.
	entry := data.Entries[1]
	entry.Value.ID = "event-extra"
	entry.Value.SourceTime = "2026-10-01T09:30:00Z"
	data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, "event-extra")
	data.Entries = append(data.Entries, entry)
	body, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(),"owners.json")
	if err := os.WriteFile(path, body, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-snapshot="+path, "-principal=analyst-a",
		"-dossier-case=case-a", "-dossier-entry=event-old",
	}, &out)
	if err != nil { t.Fatal(err) }
	var got localCaseDossier
	if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Source.Preview != marker || got.Summary.VisibleEntries != 3 ||
		len(got.RelatedEntries) != 2 || got.RelatedEntries[0].ID != "event-extra" ||
		got.RelatedEntries[1].ID != "event-old" {
		t.Fatalf("disk-to-ACL-to-source dossier failed: %+v", got)
	}
	if strings.Contains(out.String(), "demo-hidden-source") || strings.Contains(out.String(), "restricted") {
		t.Fatal("denied source leaked")
	}
}

func TestCaseDossierDeniedAndAmbiguousModesProduceNoOutput(t *testing.T) {
	cases := [][]string{
		{"-dossier-case=case-a"},
		{"-dossier-entry=event-old"},
		{"-dossier-case=case-a","-dossier-entry=restricted"},
		{"-dossier-case=case-a","-dossier-entry=not-visible"},
		{"-dossier-case=case-b","-dossier-entry=event-b"},
		{"-principal=analyst-b","-dossier-case=case-a","-dossier-entry=event-old"},
		{"-principal=unknown","-dossier-case=case-a","-dossier-entry=event-old"},
		{"-dossier-case=case-a","-dossier-entry=event-old","-report-case=case-a"},
		{"-dossier-case=case-a","-dossier-entry=event-old","-inspect-case=case-a","-inspect-entry=event-old"},
		{"-dossier-case=case-a","-dossier-entry=event-old","-compare-case=case-a","-with-case=case-b"},
		{"-dossier-case=case-a","-dossier-entry=event-old","-source-case=case-a","-source-ref=demo-evidence-a"},
		{"-dossier-case=case-a","-dossier-entry=event-old","-snapshot=relative.json"},
	}
	for _, args := range cases {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("denied dossier emitted data: %v %v %q",args,err,out.String())
		}
	}
}

func TestCaseDossierCancellationAndWriterFailure(t *testing.T) {
	ctx,cancel := context.WithCancel(context.Background()); cancel()
	var out bytes.Buffer
	if err := writeDemoCaseDossier(ctx,"analyst-a","","case-a","event-old",&out); !errors.Is(err,context.Canceled) || out.Len()!=0 {
		t.Fatalf("cancelled dossier leaked: %v %q",err,out.String())
	}
	if err := writeDemoCaseDossier(context.Background(),"analyst-a","","case-a","event-old",rejectReportOutput{}); !errors.Is(err,io.ErrClosedPipe) {
		t.Fatalf("writer failure swallowed: %v",err)
	}
}

func TestCaseDossierDistinctAuthorizedTypedSources(t *testing.T) {
	cases := []struct {
		principal, caseID, entryID, kind, object string
	}{
		{"analyst-a", "case-a", "event-new", "telemetry-event", "demo-event-a"},
		{"analyst-b", "case-b", "event-b", "incident", "demo-incident"},
	}
	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			var out bytes.Buffer
			err := runDemoCLI(context.Background(), []string{
				"-principal="+tc.principal, "-dossier-case="+tc.caseID,
				"-dossier-entry="+tc.entryID,
			}, &out)
			if err != nil { t.Fatal(err) }
			var got localCaseDossier
			if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
			if got.Source.Kind != tc.kind || got.Source.ObjectID != tc.object ||
				got.Source.EntryID != tc.entryID || got.SelectedEntry.ID != tc.entryID ||
				len(got.RelatedEntries) != 1 || got.RelatedEntries[0].ID != tc.entryID {
				t.Fatalf("Case -> Timeline -> owner source is not connected: %+v",got)
			}
		})
	}
}

func TestCaseDossierRejectsMissingAndOversizedIdentifiers(t *testing.T) {
	for _, args := range []struct{ principal, caseID, entryID string }{
		{"", "case-a", "event-old"},
		{"analyst-a", "", "event-old"},
		{"analyst-a", "case-a", ""},
		{"analyst-a", strings.Repeat("a", 1025), "event-old"},
		{"analyst-a", "case-a", strings.Repeat("b", 1025)},
		{strings.Repeat("c", 1025), "case-a", "event-old"},
	} {
		var out bytes.Buffer
		err := writeDemoCaseDossier(context.Background(), args.principal, "", args.caseID, args.entryID, &out)
		if err == nil || out.Len() != 0 {
			t.Fatalf("invalid dossier accepted: %v output length=%d", err, out.Len())
		}
	}
	var out bytes.Buffer
	if err := writeDemoCaseDossier(context.Background(), "analyst-a", filepath.Join(t.TempDir(), "absent.json"), "case-a", "event-old", &out); err == nil || out.Len() != 0 {
		t.Fatalf("missing operator snapshot did not fail closed: %v", err)
	}
}

func TestDossierSourceOwnerRevalidationRejectsContradictoryReplies(t *testing.T) {
	const valid = "{\"status\":\"ready\",\"synthetic\":true,\"case_id\":\"case-a\",\"entry_id\":\"event-old\",\"kind\":\"evidence\",\"object_id\":\"object-a\",\"preview\":\"owner-approved\"}"
	for _, tc := range []struct {
		name string
		status int
		body string
	}{
		{"denied-owner", http.StatusForbidden, valid},
		{"invalid-owner-json", http.StatusOK, "{not-json"},
		{"not-ready", http.StatusOK, strings.Replace(valid, "\"status\":\"ready\"", "\"status\":\"deny\"", 1)},
		{"non-synthetic", http.StatusOK, strings.Replace(valid, "\"synthetic\":true", "\"synthetic\":false", 1)},
		{"foreign-case", http.StatusOK, strings.Replace(valid, "\"case_id\":\"case-a\"", "\"case_id\":\"case-b\"", 1)},
		{"different-entry", http.StatusOK, strings.Replace(valid, "\"entry_id\":\"event-old\"", "\"entry_id\":\"restricted\"", 1)},
		{"missing-source-kind", http.StatusOK, strings.Replace(valid, "\"kind\":\"evidence\"", "\"kind\":\"\"", 1)},
		{"missing-owner-object", http.StatusOK, strings.Replace(valid, "\"object_id\":\"object-a\"", "\"object_id\":\"\"", 1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/demo-source-open" || r.URL.Query().Get("case_id") != "case-a" ||
					r.URL.Query().Get("entry_id") != "event-old" ||
					r.Header.Get("Authorization") != "Bearer "+testBearer {
					t.Error("source-owner query was not bound to selected authorized entry")
				}
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			})
			if _, err := inspectAuthorizedDemoEntry(context.Background(), handler, testBearer, "case-a", "event-old"); err == nil {
				t.Fatal("contradictory owner reply was accepted")
			}
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, valid) })
	got, err := inspectAuthorizedDemoEntry(context.Background(), handler, testBearer, "case-a", "event-old")
	if err != nil || got.Preview != "owner-approved" || got.ObjectID != "object-a" {
		t.Fatalf("authorized typed owner response rejected: %+v %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := inspectAuthorizedDemoEntry(ctx, handler, testBearer, "case-a", "event-old"); !errors.Is(err, context.Canceled) {
		t.Fatalf("source owner cancellation not propagated: %v", err)
	}
}

func TestCaseDossierRejectsImportExportModeConflicts(t *testing.T) {
	for _, args := range [][]string{
		{"-dossier-case=case-a", "-dossier-entry=event-old", "-export-case=case-a"},
		{"-dossier-case=case-a", "-dossier-entry=event-old", "-import-jsonl=/tmp/events.jsonl", "-into-case=case-a"},
		{"-dossier-case=case-a", "-dossier-entry=event-old", "-triage-import"},
	} {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("ambiguous dossier mode accepted (%v): err=%v bytes=%d", args, err, output.Len())
		}
	}
}
