package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
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
	for _, forbidden := range []string{"demo-hidden-source", "demo-tenant-b", "demo-secret", "restricted", "\"case-b\""} {
		if bytes.Contains(out.Bytes(), []byte(forbidden)) {
			t.Fatalf("legacy import exported inaccessible owner records: %q", forbidden)
		}
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
	if len(derived.Cases) != 1 || len(derived.Principals) != 1 ||
		derived.Principals[0].ID != "analyst-a" ||
		derived.Cases[0].ID != "forensic-case-on-disk" ||
		len(derived.Cases[0].EntryIDs) != 4 || len(derived.Entries) != 4 {
		t.Fatalf("import failed to preserve exactly the readable Case projection: %+v", derived)
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
	if err := json.Unmarshal(out.Bytes(), &got); err != nil || len(got.Cases) != 1 ||
		len(got.Principals) != 1 || len(got.Entries) != 4 ||
		len(got.Cases[0].EntryIDs) != 4 ||
		bytes.Contains(out.Bytes(), []byte("demo-tenant-b")) {
		t.Fatalf("CLI exported data beyond the selected authorized Case: %v", err)
	}
}

type rejectingSnapshotWriter struct{}

func (rejectingSnapshotWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestLocalJSONLImportStructuralBoundsAndIdentityConflicts(t *testing.T) {
	data := demoData()
	t.Run("existing-entry-identity", func(t *testing.T) {
		clone := data
		clone.Cases = append([]timeline.MemoryTimelineCase(nil), data.Cases...)
		clone.Entries = append([]timeline.MemoryTimelineEntry(nil), data.Entries...)
		clone.Cases[0].EntryIDs = append([]string(nil), data.Cases[0].EntryIDs...)
		sum := sha256.Sum256([]byte("case-a\x00sysmon-1"))
		id := "local-jsonl-entry-" + hex.EncodeToString(sum[:])
		clone.Cases[0].EntryIDs[0] = id
		clone.Entries[0].Value.ID = id
		if _, err := importLocalEvents(context.Background(), clone, "analyst-a", "case-a", []byte(importedJSONL)); err == nil {
			t.Fatal("import overwrote a pre-existing Timeline Entry")
		}
	})
	t.Run("existing-source-identity", func(t *testing.T) {
		clone := data
		clone.Sources = append([]timeline.MemoryTimelineSource(nil), data.Sources...)
		clone.Entries = append([]timeline.MemoryTimelineEntry(nil), data.Entries...)
		sum := sha256.Sum256([]byte("case-a\x00sysmon-1"))
		ref := "local-jsonl-source-" + hex.EncodeToString(sum[:])
		clone.Sources[0].Identity.Reference = ref
		clone.Entries[0].Value.SourceRef = ref
		if _, err := importLocalEvents(context.Background(), clone, "analyst-a", "case-a", []byte(importedJSONL)); err == nil {
			t.Fatal("import overwrote an existing typed source")
		}
	})
	for _, input := range []struct {
		name string
		data []byte
	}{
		{"oversized-content", bytes.Repeat([]byte("x"), maxDemoSnapshotBytes+1)},
		{"too-many-lines", bytes.Repeat([]byte("\n"), timeline.MaxEntries+2)},
		{"oversized-line", []byte(strings.Repeat("x", 4097))},
		{"invalid-unicode", []byte("{\"id\":\"x\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"\xff\"}")},
		{"overlong-message", []byte("{\"id\":\"x\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"" + strings.Repeat("a", 513) + "\"}")},
		{"invalid-control-character", []byte("{\"id\":\"x\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"bad\\nmessage\"}")},
	} {
		t.Run(input.name, func(t *testing.T) {
			if _, err := importLocalEvents(context.Background(), data, "analyst-a", "case-a", input.data); err == nil {
				t.Fatal("invalid input was promoted to owner snapshot")
			}
		})
	}
	for _, tc := range []struct{ principal, caseID string }{
		{"analyst-a", "case-b"}, {"analyst-b", "case-a"}, {"analyst-a", "missing"},
		{"", "case-a"}, {"analyst-a", ""}, {"analyst-a", strings.Repeat("c", 1025)},
	} {
		if _, err := importLocalEvents(context.Background(), data, tc.principal, tc.caseID, []byte(importedJSONL)); err == nil {
			t.Fatalf("unauthorized Case import accepted: %+v", tc)
		}
	}
}

func TestLocalJSONLImportMaintainsOwnerSnapshotAndWriterFailures(t *testing.T) {
	original := demoData()
	before := append([]string(nil), original.Cases[0].EntryIDs...)
	merged, err := importLocalEvents(context.Background(), original, "analyst-a", "case-a", []byte(importedJSONL))
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Cases[0].EntryIDs) != len(before)+2 || len(original.Cases[0].EntryIDs) != len(before) ||
		strings.Join(before, ",") != strings.Join(original.Cases[0].EntryIDs, ",") {
		t.Fatal("import changed caller owned snapshot or omitted new events")
	}
	short := strings.TrimSuffix(importedJSONL, "\n")
	if next, err := importLocalEvents(context.Background(), original, "analyst-a", "case-a", []byte(short)); err != nil ||
		len(next.Entries) != len(original.Entries)+2 {
		t.Fatalf("valid JSONL without trailing newline was rejected: %v", err)
	}
	ownerPath := snapshotFixture(t)
	eventPath := writeJSONLEvents(t, importedJSONL)
	for _, tc := range []struct{ snapshot, events string }{
		{"", eventPath}, {ownerPath, ""}, {ownerPath, ownerPath},
		{ownerPath, "relative.jsonl"},
	} {
		var output bytes.Buffer
		if err := writeDemoImportedJSONL(context.Background(), "analyst-a", tc.snapshot,
			"forensic-case-on-disk", tc.events, &output); err == nil || output.Len() > 0 {
			t.Fatalf("untrusted path returned owner records: %v", err)
		}
	}
	err = writeDemoImportedJSONL(context.Background(), "analyst-a", ownerPath,
		"forensic-case-on-disk", eventPath, rejectingSnapshotWriter{})
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("failure of operator output target was swallowed: %v", err)
	}
}
