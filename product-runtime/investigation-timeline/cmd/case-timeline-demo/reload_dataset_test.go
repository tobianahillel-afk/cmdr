package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func saveReloadableOwnerSnapshot(t *testing.T, path string, data timeline.MemoryTimelineData) {
	t.Helper()
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestOfflineOwnerReloadRechecksLiveSourceAndCasePermissions(t *testing.T) {
	data, err := loadDemoDataset(filepath.Join("testdata", "example-owners.json"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "owners.json")
	saveReloadableOwnerSnapshot(t, path, data)
	handler, err := newReloadingDemoHandler(testBearer, "investigator-42", path)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	nav := "/case-source-navigation?case_id=case-incident-42&entry_id=alert-1"
	check := func(route, token string, want int, forbidden string) {
		t.Helper()
		status, body := requestDemo(t, server, route, token)
		if status != want || (forbidden != "" && strings.Contains(body, forbidden)) {
			t.Fatalf("HTTP %d want %d, leaked %q in %s", status, want, forbidden, body)
		}
	}
	check(nav, testBearer, http.StatusOK, "")
	check("/case-timeline?case_id=case-incident-42", testBearer, http.StatusOK, "")
	check(nav, "", http.StatusUnauthorized, "event-1")
	check("/case-timeline?case_id=case-private-99", testBearer, http.StatusNotFound, "alert-99")

	// The operator updates the trusted owner snapshot: revoke the linked
	// source while keeping Case/Timeline Entry visibility unchanged.
	for i := range data.Sources {
		if data.Sources[i].Identity.Reference == "telemetry-1" {
			data.Sources[i].Readers = nil
		}
	}
	saveReloadableOwnerSnapshot(t, path, data)
	check(nav, testBearer, http.StatusNotFound, "event-1")
	check("/case-timeline?case_id=case-incident-42", testBearer, http.StatusOK, "event-1")

	// After restoring the source, revoke the entire Case. Existing URL,
	// return context and bearer must not preserve the previous access.
	for i := range data.Cases {
		if data.Cases[i].ID == "case-incident-42" {
			data.Cases[i].Readers = nil
		}
	}
	saveReloadableOwnerSnapshot(t, path, data)
	check(nav, testBearer, http.StatusNotFound, "event-1")
	check("/case-timeline?case_id=case-incident-42", testBearer, http.StatusNotFound, "alert-1")

	// Corruption of the host-supplied snapshot cannot resurrect stale ACLs.
	if err := os.WriteFile(path, []byte("{not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	check(nav, testBearer, http.StatusServiceUnavailable, "event-1")
	check("/demo-cases", testBearer, http.StatusServiceUnavailable, "case-incident-42")
	// Anonymous traffic cannot force parsing of a corrupted owner file,
	// even while the authorized caller sees fail-closed unavailability.
	check(nav, "", http.StatusUnauthorized, "event-1")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	check("/case-timeline?case_id=case-incident-42", testBearer, http.StatusServiceUnavailable, "alert-1")

	// Freshly restored permissions are effective on the next request.
	for i := range data.Cases {
		if data.Cases[i].ID == "case-incident-42" {
			data.Cases[i].Readers = []string{"investigator-42"}
		}
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.Reference == "telemetry-1" {
			data.Sources[i].Readers = []string{"investigator-42"}
		}
	}
	saveReloadableOwnerSnapshot(t, path, data)
	check(nav, testBearer, http.StatusOK, "")
}

func TestOfflineOwnerReloadRequiresExplicitTrustedLocalDataset(t *testing.T) {
	if _, err := newReloadingDemoHandler(testBearer, "investigator-42", filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing host dataset accepted")
	}
	if _, err := newReloadingDemoHandler("invalid", "investigator-42", filepath.Join("testdata", "example-owners.json")); err == nil {
		t.Fatal("invalid local bearer accepted")
	}
}
