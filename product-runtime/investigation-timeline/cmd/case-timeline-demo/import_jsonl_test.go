package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

const importedJSONL = "{\"id\":\"sysmon-1\",\"timestamp\":\"2026-10-15T10:20:30+02:00\",\"message\":\"Offline Sysmon process event\"}\n" +
	"{\"id\":\"sysmon-2\",\"timestamp\":\"2026-10-15T10:20:31Z\",\"message\":\"Offline Windows log event\"}\n"

func writeJSONLEvents(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "telemetry.jsonl")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLocalJSONLImportDiskToAuthenticatedTimelineSource(t *testing.T) {
	ownerPath := snapshotFixture(t)
	original, err := os.ReadFile(ownerPath)
	if err != nil {
		t.Fatal(err)
	}
	eventsPath := writeJSONLEvents(t, importedJSONL)
	var out bytes.Buffer
	if err := writeDemoImportedJSONL(context.Background(), "analyst-a", ownerPath, "forensic-case-on-disk", eventsPath, &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() == 0 || out.Len() > maxDemoSnapshotBytes {
		t.Fatal("missing or oversized owner snapshot output")
	}
	check, err := os.ReadFile(ownerPath)
	if err != nil || !bytes.Equal(check, original) {
		t.Fatal("import modified operator input snapshot")
	}
	derivedPath := filepath.Join(t.TempDir(), "derived.json")
	if err := os.WriteFile(derivedPath, out.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	derived, err := loadDemoSnapshot(derivedPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(derived.Cases[0].EntryIDs) != 5 {
		t.Fatalf("events not integrated into Case: %+v", derived.Cases[0])
	}
	handler, err := newDemoHandlerWithData(testBearer, "analyst-a", derived)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	status, body := requestDemo(t, server, "/case-timeline?case_id=forensic-case-on-disk&limit=20", testBearer)
	if status != http.StatusOK || strings.Contains(body, "restricted") || strings.Contains(body, "demo-hidden-source") ||
		strings.Contains(body, "event-b") {
		t.Fatalf("imported timeline isolation failed: %d %s", status, body)
	}
	digest := sha256.Sum256([]byte("forensic-case-on-disk\x00sysmon-1"))
	id := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
	if !strings.Contains(body, id) || !strings.Contains(body, "operator-local-jsonl") ||
		!strings.Contains(body, "2026-10-15T10:20:30+02:00") {
		t.Fatalf("disk event was not projected: %s", body)
	}
	for _, route := range []string{
		"/demo-source-open?case_id=forensic-case-on-disk&entry_id=" + url.QueryEscape(id),
		"/case-source-navigation?case_id=forensic-case-on-disk&entry_id=" + url.QueryEscape(id),
	} {
		code, source := requestDemo(t, server, route, testBearer)
		if code != http.StatusOK {
			t.Fatalf("typed source from imported event inaccessible: %d %s", code, source)
		}
		if strings.Contains(route, "demo-source-open") && !strings.Contains(source, "Offline Sysmon process event") {
			t.Fatalf("source preview did not come from JSONL: %s", source)
		}
		code, leaked := requestDemo(t, server, route, "")
		if code != http.StatusUnauthorized || strings.Contains(leaked, "Sysmon") {
			t.Fatalf("anonymous source disclosure: %d %s", code, leaked)
		}
	}
	forbidden, err := newDemoHandlerWithData(testBearer, "analyst-b", derived)
	if err != nil {
		t.Fatal(err)
	}
	foreign := httptest.NewServer(forbidden)
	defer foreign.Close()
	code, leaked := requestDemo(t, foreign, "/case-timeline?case_id=forensic-case-on-disk", testBearer)
	if code != http.StatusNotFound || strings.Contains(leaked, id) {
		t.Fatalf("foreign tenant saw imported records: %d %s", code, leaked)
	}
}

func TestLocalJSONLImportFailsClosedWithoutOutput(t *testing.T) {
	owner := snapshotFixture(t)
	inputs := []struct{ name, text string }{
		{"duplicate", importedJSONL + "{\"id\":\"sysmon-1\",\"timestamp\":\"2026-10-15T10:20:32Z\",\"message\":\"bad\"}\n"},
		{"invalid-time", "{\"id\":\"x\",\"timestamp\":\"yesterday\",\"message\":\"bad\"}\n"},
		{"forged-access", "{\"id\":\"x\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"bad\",\"readers\":[\"analyst-b\"]}\n"},
		{"missing-message", "{\"id\":\"x\",\"timestamp\":\"2026-10-15T10:20:30Z\"}\n"},
		{"invalid-id", "{\"id\":\"../event\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"bad\"}\n"},
		{"blank-line", "\n" + importedJSONL},
		{"extra-json", "{\"id\":\"x\"} {}"},
		{"empty", ""},
	}
	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := writeDemoImportedJSONL(context.Background(), "analyst-a", owner, "forensic-case-on-disk",
				writeJSONLEvents(t, tc.text), &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid import emitted owner records: %v", err)
			}
		})
	}
	for _, principal := range []string{"analyst-b", "missing"} {
		var out bytes.Buffer
		if err := writeDemoImportedJSONL(context.Background(), principal, owner, "forensic-case-on-disk",
			writeJSONLEvents(t, importedJSONL), &out); err == nil || out.Len() > 0 {
			t.Fatalf("unauthorized %s import accepted", principal)
		}
	}
	snapshot, err := loadDemoSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Cases[0].EntryIDs = append(snapshot.Cases[0].EntryIDs, "event-b")
	if _, err := importLocalEvents(context.Background(), snapshot, "analyst-a", "forensic-case-on-disk", []byte(importedJSONL)); err == nil {
		t.Fatal("foreign owner graph was promoted to import authority")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := importLocalEvents(ctx, demoData(), "analyst-a", "case-a", []byte(importedJSONL)); err == nil {
		t.Fatal("cancelled import accepted")
	}
}

func TestLocalJSONLImportPathsAndCaseCapacity(t *testing.T) {
	owner := snapshotFixture(t)
	path := writeJSONLEvents(t, importedJSONL)
	linked := filepath.Join(t.TempDir(), "events-link.jsonl")
	if err := os.Symlink(path, linked); err == nil {
		if _, err := readLocalEventFile(linked); err == nil {
			t.Fatal("symlink accepted")
		}
	}
	for _, invalid := range []string{"relative.jsonl", t.TempDir(), filepath.Join(t.TempDir(), "missing")} {
		if _, err := readLocalEventFile(invalid); err == nil {
			t.Fatalf("invalid event path accepted: %q", invalid)
		}
	}
	snapshot, err := loadDemoSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Cases[0].EntryIDs = make([]string, timeline.MaxEntries)
	if _, err := importLocalEvents(context.Background(), snapshot, "analyst-a", "forensic-case-on-disk", []byte(importedJSONL)); err == nil {
		t.Fatal("over-capacity Case import accepted")
	}
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner, "-import-jsonl=" + path,
		"-into-case=forensic-case-on-disk",
	}, &out); err != nil || out.Len() == 0 {
		t.Fatalf("CLI imported no snapshot: %v", err)
	}
	var got timeline.MemoryTimelineData
	if err := json.Unmarshal(out.Bytes(), &got); err != nil || len(got.Entries) != len(demoData().Entries)+2 {
		t.Fatalf("CLI output not reusable: %v", err)
	}
}
