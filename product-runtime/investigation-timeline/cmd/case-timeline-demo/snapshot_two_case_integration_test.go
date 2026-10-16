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

// Prove the browser's comparative source is a real operator disk snapshot,
// not two fabricated HTTP pages: two Cases share the same permitted analyst,
// while their typed owner source records remain independent and Case-scoped.
func TestOperatorSnapshotTwoAuthorizedCasesHaveIndependentTimelineAndSourceACLs(t *testing.T) {
	data := demoData()
	data.Cases = append(data.Cases, timeline.MemoryTimelineCase{
		ID: "case-c", Tenant: "demo-tenant-a", Environment: "lab",
		Readers: []string{"analyst-a"}, EntryIDs: []string{"disk-event-c"},
	})
	data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
		CaseID: "case-c",
		Value: timeline.Entry{
			ID: "disk-event-c", Tenant: "demo-tenant-a", Kind: "inferred",
			SourceTime: "2026-10-01T10:30:00Z", Timezone: "+00:00",
			Version: 1, SourceRef: "disk-source-c", Provenance: "only-on-disk-case-c",
		},
		Readers: []string{"analyst-a"},
	})
	data.Sources = append(data.Sources, timeline.MemoryTimelineSource{
		Identity: timeline.SourceOwnerIdentity{
			Reference: "disk-source-c", Kind: "incident", ObjectID: "incident-c",
			Tenant: "demo-tenant-a", Environment: "lab", CaseID: "case-c",
		},
		Readers: []string{"analyst-a"}, Preview: "only-on-disk-preview-c",
	})
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "two-case-owner.json")
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
		t.Fatalf("two authorized Cases not discoverable: %d %s", status, body)
	}
	var kinds = map[string]map[string]int{}
	for _, id := range []string{"case-a", "case-c"} {
		status, body = requestDemo(t, server, "/case-timeline?case_id="+id+"&limit=200", testBearer)
		if status != http.StatusOK {
			t.Fatalf("Case %s failed HTTP authorization: %d %s", id, status, body)
		}
		var page struct {
			Entries []timeline.ProjectedEntry `json:"entries"`
		}
		if err := json.Unmarshal([]byte(body), &page); err != nil {
			t.Fatal(err)
		}
		kinds[id] = map[string]int{}
		for _, row := range page.Entries {
			kinds[id][row.Kind]++
			if id == "case-a" && row.SourceRef == "disk-source-c" ||
				id == "case-c" && row.SourceRef != "disk-source-c" {
				t.Fatalf("cross-Case source contamination in %s: %+v", id, row)
			}
		}
		if strings.Contains(body, "only-on-disk-preview-c") ||
			strings.Contains(body, "demo-hidden-source") ||
			strings.Contains(body, "demo-tenant-b") {
			t.Fatalf("source preview or forbidden metadata leaked into %s: %s", id, body)
		}
	}
	if kinds["case-a"]["observed"] != 1 || kinds["case-a"]["inferred"] != 1 ||
		kinds["case-c"]["observed"] != 0 || kinds["case-c"]["inferred"] != 1 {
		t.Fatalf("incorrect independent observed/inferred populations: %+v", kinds)
	}
	status, body = requestDemo(t, server, "/demo-source-open?case_id=case-c&entry_id=disk-event-c", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "only-on-disk-preview-c") {
		t.Fatalf("authorized disk source not opened through owner: %d %s", status, body)
	}
	for _, target := range []string{
		"/case-timeline?case_id=case-b",
		"/demo-source-open?case_id=case-a&entry_id=disk-event-c",
	} {
		status, body := requestDemo(t, server, target, testBearer)
		if status != http.StatusNotFound || strings.Contains(body, "only-on-disk-preview-c") {
			t.Fatalf("foreign Case/source accepted: %s %d %s", target, status, body)
		}
	}
	status, body = requestDemo(t, server, "/case-timeline?case_id=case-c", "")
	if status != http.StatusUnauthorized || strings.Contains(body, "only-on-disk-case-c") {
		t.Fatalf("anonymous operator source disclosed: %d %s", status, body)
	}
}
