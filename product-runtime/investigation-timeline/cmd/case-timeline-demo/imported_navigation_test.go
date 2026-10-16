package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

// Real HTTP integration: a trusted local owner file drives Case discovery,
// Timeline projection and an independently reauthorized source destination.
func TestImportedOwnerSnapshotToSourceNavigationAndSafeReturn(t *testing.T) {
	data, err := loadDemoDataset(filepath.Join("testdata", "example-owners.json"))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "investigator-42", data)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	status, body := requestDemo(t, server, "/demo-cases", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "case-incident-42") ||
		strings.Contains(body, "case-private-99") {
		t.Fatalf("imported Case discovery leaked foreign Case: %d %s", status, body)
	}
	status, body = requestDemo(t, server, "/case-timeline?case_id=case-incident-42&limit=20", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "\"id\":\"alert-1\"") ||
		!strings.Contains(body, "\"id\":\"alert-2\"") || strings.Contains(body, "alert-99") {
		t.Fatalf("imported Timeline integration failed: %d %s", status, body)
	}
	navigation := "/case-source-navigation?case_id=case-incident-42&entry_id=alert-1&limit=20&offset=0"
	status, body = requestDemo(t, server, navigation, testBearer)
	if status != http.StatusOK {
		t.Fatalf("imported source navigation denied: %d %s", status, body)
	}
	var draft struct {
		Status            string `json:"status"`
		Kind              string `json:"kind"`
		ObjectID          string `json:"object_id"`
		ReturnCaseID      string `json:"return_case_id"`
		ReturnEntryID     string `json:"return_entry_id"`
		ReturnQuery       string `json:"return_query"`
		ReturnIsAuthority bool   `json:"return_is_authority"`
	}
	if err := json.Unmarshal([]byte(body), &draft); err != nil {
		t.Fatal(err)
	}
	if draft.Status != "ready" || draft.Kind != "telemetry-event" ||
		draft.ObjectID != "event-1" || draft.ReturnCaseID != "case-incident-42" ||
		draft.ReturnEntryID != "alert-1" || draft.ReturnIsAuthority ||
		draft.ReturnQuery != "case_id=case-incident-42&limit=20&offset=0" {
		t.Fatalf("wrong source-owner navigation/return contract: %+v", draft)
	}
	for _, hidden := range []string{"alert-99", "case-private-99", "incident-99", "investigator-99", testBearer} {
		if strings.Contains(body, hidden) {
			t.Fatalf("imported owner navigation disclosed %s", hidden)
		}
	}
	status, body = requestDemo(t, server, "/case-timeline?"+draft.ReturnQuery, "")
	if status != http.StatusUnauthorized || strings.Contains(body, "alert-1") {
		t.Fatalf("return context borrowed source authority: %d %s", status, body)
	}
	status, body = requestDemo(t, server, "/case-timeline?"+draft.ReturnQuery, testBearer)
	if status != http.StatusOK || !strings.Contains(body, "alert-1") {
		t.Fatalf("return failed to reauthorize Case: %d %s", status, body)
	}
	for _, path := range []string{
		"/case-source-navigation?case_id=case-private-99&entry_id=alert-99",
		"/case-source-navigation?case_id=case-incident-42&entry_id=alert-99",
		navigation + "&object_id=incident-99",
		navigation + "&tenant=local-tenant-b",
	} {
		status, body := requestDemo(t, server, path, testBearer)
		if status == http.StatusOK || strings.Contains(body, "incident-99") {
			t.Fatalf("untrusted or foreign navigation accepted: %s %d %s", path, status, body)
		}
	}
	status, body = requestDemo(t, server, navigation, "")
	if status != http.StatusUnauthorized || strings.Contains(body, "event-1") {
		t.Fatalf("anonymous source navigation leaked: %d %s", status, body)
	}

	// Source-owner readers are authoritative in this offline snapshot: a
	// permitted Timeline Entry must not become a navigable denied source.
	for i := range data.Sources {
		if data.Sources[i].Identity.Reference == "telemetry-1" {
			data.Sources[i].Readers = nil
		}
	}
	restricted, err := newDemoHandlerWithData(testBearer, "investigator-42", data)
	if err != nil {
		t.Fatal(err)
	}
	privateServer := httptest.NewServer(restricted)
	defer privateServer.Close()
	req := "/case-source-navigation?case_id="+url.QueryEscape("case-incident-42")+"&entry_id=alert-1"
	status, body = requestDemo(t, privateServer, req, testBearer)
	if status != http.StatusNotFound || strings.Contains(body, "event-1") {
		t.Fatalf("source-owner revocation bypassed: %d %s", status, body)
	}
}
