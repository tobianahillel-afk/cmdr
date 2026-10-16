package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestImportedCaseVersionComparisonAcrossPageAndFilters(t *testing.T) {
	data, err := loadDemoDataset(filepath.Join("testdata", "example-owners.json"))
	if err != nil { t.Fatal(err) }
	handler, err := newDemoHandlerWithData(testBearer, "investigator-42", data)
	if err != nil { t.Fatal(err) }
	server := httptest.NewServer(handler)
	defer server.Close()
	route := "/case-version-comparison?case_id=case-incident-42&entry_id=alert-2"
	status, body := requestDemo(t, server, route, testBearer)
	if status != http.StatusOK { t.Fatalf("version comparison inaccessible: %d %s", status, body) }
	var result struct {
		Status string `json:"status"`
		Current timeline.ProjectedEntry `json:"current"`
		Previous timeline.ProjectedEntry `json:"previous"`
		IsAuthority bool `json:"is_authority"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil { t.Fatal(err) }
	if result.Status != "ready" || result.IsAuthority ||
		result.Current.ID != "alert-2" || result.Current.Version != 2 ||
		result.Current.Kind != "inferred" || result.Current.Supersedes != "alert-1" ||
		result.Previous.ID != "alert-1" || result.Previous.Version != 1 ||
		result.Previous.Kind != "observed" {
		t.Fatalf("wrong authorized version comparison: %+v", result)
	}
	// Independent selection does not rely on the current filtered page.
	status, body = requestDemo(t, server,
		"/case-timeline?case_id=case-incident-42&kind=inferred&limit=1", testBearer)
	if status != http.StatusOK || strings.Contains(body, "\"id\":\"alert-1\"") {
		t.Fatalf("comparison is not independent of the filtered page: %d %s", status, body)
	}
	for _, route := range []string{
		"/case-version-comparison?case_id=case-private-99&entry_id=alert-99",
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-99",
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-1",
	} {
		status, body := requestDemo(t, server, route, testBearer)
		if status != http.StatusNotFound || strings.Contains(body, "incident-99") {
			t.Fatalf("inaccessible version accepted: %s %d %s", route, status, body)
		}
	}
	for _, route := range []string{
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-2&tenant=local-tenant-b",
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-2&object_id=evidence-2",
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-2&entry_id=alert-99",
		"/case-version-comparison?case_id=" + url.QueryEscape("*") + "&entry_id=alert-2",
	} {
		status, _ := requestDemo(t, server, route, testBearer)
		if status != http.StatusBadRequest { t.Fatalf("tampered comparison accepted %s: %d", route, status) }
	}
	status, body = requestDemo(t, server, "/case-version-comparison?case_id=case-incident-42&entry_id=alert-2", "")
	if status != http.StatusUnauthorized || strings.Contains(body, "alert-1") {
		t.Fatalf("anonymous version disclosure: %d %s", status, body)
	}
	// Even if the current entry remains authorized, a revoked predecessor
	// MUST disappear; no ID, source ref, or provenance can be returned.
	for i := range data.Sources {
		if data.Sources[i].Identity.Reference == "telemetry-1" {
			data.Sources[i].Readers = nil
		}
	}
	restricted, err := newDemoHandlerWithData(testBearer, "investigator-42", data)
	if err != nil { t.Fatal(err) }
	private := httptest.NewServer(restricted)
	defer private.Close()
	status, body = requestDemo(t, private,
		"/case-version-comparison?case_id=case-incident-42&entry_id=alert-2", testBearer)
	if status != http.StatusNotFound || strings.Contains(body, "telemetry-1") ||
		strings.Contains(body, "alert-1") {
		t.Fatalf("hidden predecessor leaked: %d %s", status, body)
	}
}
