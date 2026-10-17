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

// Real owner snapshot -> real Go HTTP Case/source authorization and text search.
// Distinct Cases can contain the same IOC-like text without sharing source IDs.
func TestOperatorDiskSnapshotCrossCaseTextHuntUsesRealOwnerACLs(t *testing.T) {
	data := demoData()
	data.Entries[0].Value.Provenance = "dns-observed:10.0.0.7"
	data.Entries[2].Value.Provenance = "hidden-observed:10.0.0.7"
	data.Cases = append(data.Cases, timeline.MemoryTimelineCase{
		ID: "case-c", Tenant: "demo-tenant-a", Environment: "lab",
		Readers: []string{"analyst-a"}, EntryIDs: []string{"disk-c-ioc"},
	})
	data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
		CaseID: "case-c", Readers: []string{"analyst-a"},
		Value: timeline.Entry{
			ID: "disk-c-ioc", Kind: "inferred", Tenant: "demo-tenant-a",
			SourceTime: "2026-10-01T12:00:00Z", Timezone: "+00:00",
			Version: 1, SourceRef: "case-c-own-source", Provenance: "incident:10.0.0.7",
		},
	})
	data.Sources = append(data.Sources, timeline.MemoryTimelineSource{
		Identity: timeline.SourceOwnerIdentity{
			Reference: "case-c-own-source", Kind: "incident", ObjectID: "case-c-ioc",
			Tenant: "demo-tenant-a", Environment: "lab", CaseID: "case-c",
		},
		Readers: []string{"analyst-a"}, Preview: "only-local-case-c-authorized-preview",
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "operator-ioc.json")
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	imported, err := loadDemoSnapshot(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := newDemoHandlerWithData(testBearer, "analyst-a", imported)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()

	status, body := requestDemo(t, server, "/demo-cases", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "case-a") ||
		!strings.Contains(body, "case-c") || strings.Contains(body, "case-b") {
		t.Fatalf("owner Case discovery failed: %d %s", status, body)
	}
	found := map[string]string{}
	for _, caseID := range []string{"case-a", "case-c"} {
		status, body = requestDemo(t, server,
			"/case-timeline?case_id="+caseID+"&q=10.0.0.7&limit=200&offset=0",
			testBearer)
		if status != http.StatusOK {
			t.Fatalf("%s denied authorized text search: %d %s", caseID, status, body)
		}
		var response struct {
			Entries []timeline.ProjectedEntry `json:"entries"`
			Page struct {
				TotalVisible int `json:"total_visible"`
			} `json:"page"`
		}
		if err := json.Unmarshal([]byte(body), &response); err != nil {
			t.Fatal(err)
		}
		if len(response.Entries) != 1 || response.Page.TotalVisible != 1 {
			t.Fatalf("%s leaked hidden match or failed to find text: %s", caseID, body)
		}
		found[caseID] = response.Entries[0].SourceRef
		if strings.Contains(body, "demo-hidden-source") ||
			strings.Contains(body, "hidden-observed") ||
			strings.Contains(body, "only-local-case-c-authorized-preview") {
			t.Fatalf("private source details leaked from text hunt: %s", body)
		}
	}
	if found["case-a"] == found["case-c"] ||
		found["case-c"] != "case-c-own-source" {
		t.Fatalf("source references unexpectedly shared across owners: %+v", found)
	}
	status, body = requestDemo(t, server,
		"/case-timeline?case_id=case-b&q=10.0.0.7&limit=200", testBearer)
	if status != http.StatusNotFound || strings.Contains(body, "case-c-own-source") {
		t.Fatalf("cross-tenant text search exposed: %d %s", status, body)
	}
	status, body = requestDemo(t, server,
		"/case-timeline?case_id=case-c&q=10.0.0.7&limit=200", "")
	if status != http.StatusUnauthorized || strings.Contains(body, "case-c-own-source") {
		t.Fatalf("unauthenticated text search exposed: %d %s", status, body)
	}
	status, body = requestDemo(t, server,
		"/demo-source-open?case_id=case-c&entry_id=disk-c-ioc", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "only-local-case-c-authorized-preview") {
		t.Fatalf("authorized text match cannot navigate to its owner: %d %s", status, body)
	}
	status, body = requestDemo(t, server,
		"/demo-source-open?case_id=case-a&entry_id=disk-c-ioc", testBearer)
	if status != http.StatusNotFound || strings.Contains(body, "only-local-case-c-authorized-preview") {
		t.Fatalf("cross-Case source preview exposed: %d %s", status, body)
	}
}
