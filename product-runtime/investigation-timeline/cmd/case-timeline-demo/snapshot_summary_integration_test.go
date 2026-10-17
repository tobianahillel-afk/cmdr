package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Exercise the same paged, owner-authorized HTTP feed used by the browser
// summary, but with 205 rows loaded from an actual operator JSON file.
func TestDiskSnapshotFeedsAuthorizedMultiPageInvestigationSummary(t *testing.T) {
	data := demoData()
	data.Cases[0].EntryIDs = []string{"restricted"} // denied source remains in catalog
	for i := 0; i < 205; i++ {
		id := fmt.Sprintf("disk-entry-%03d", i)
		kind := "observed"
		if i%2 != 0 {
			kind = "inferred"
		}
		data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, id)
		data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
			CaseID: "case-a",
			Value: timeline.Entry{
				ID: id, Tenant: "demo-tenant-a", Kind: kind,
				SourceTime: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC).
					Add(time.Duration(i) * time.Second).Format(time.RFC3339),
				Timezone: "+00:00", Version: 1,
				SourceRef: "demo-evidence-a", Provenance: "disk-summary-unique-provenance",
			},
			Readers: []string{"analyst-a"},
		})
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "operator-case-summary.json")
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

	var seen = make(map[string]bool, 205)
	var byKind = map[string]int{"observed": 0, "inferred": 0}
	for _, offset := range []int{0, 200} {
		target := fmt.Sprintf("/case-timeline?case_id=case-a&limit=200&offset=%d", offset)
		status, body := requestDemo(t, server, target, testBearer)
		if status != http.StatusOK {
			t.Fatalf("disk-backed page %d returned %d: %s", offset, status, body)
		}
		var page struct {
			Entries []timeline.ProjectedEntry `json:"entries"`
			Page struct {
				Offset int `json:"offset"`
				Limit int `json:"limit"`
				TotalVisible int `json:"total_visible"`
				HasMore bool `json:"has_more"`
			} `json:"page"`
		}
		if err := json.Unmarshal([]byte(body), &page); err != nil {
			t.Fatal(err)
		}
		wantLen := 200
		if offset == 200 {
			wantLen = 5
		}
		if page.Page.Offset != offset || page.Page.Limit != 200 ||
			page.Page.TotalVisible != 205 || page.Page.HasMore != (offset == 0) ||
			len(page.Entries) != wantLen {
			t.Fatalf("inconsistent disk Case page at %d: %+v", offset, page.Page)
		}
		for _, e := range page.Entries {
			if seen[e.ID] || !strings.HasPrefix(e.ID, "disk-entry-") ||
				e.Provenance != "disk-summary-unique-provenance" ||
				e.SourceRef != "demo-evidence-a" || e.TimeUncertain {
				t.Fatalf("non-disk, duplicated or unauthorized row: %+v", e)
			}
			seen[e.ID] = true
			byKind[e.Kind]++
		}
		for _, private := range []string{"demo-hidden-source", "demo-secret", "restricted", "demo-tenant-b"} {
			if strings.Contains(body, private) {
				t.Fatalf("inaccessible disk owner leaked %q", private)
			}
		}
	}
	if len(seen) != 205 || byKind["observed"] != 103 || byKind["inferred"] != 102 {
		t.Fatalf("wrong authorized aggregation: %d rows, kinds=%v", len(seen), byKind)
	}
	for _, path := range []string{
		"/case-timeline?case_id=case-a&limit=200&offset=0",
		"/case-timeline?case_id=case-b&limit=200&offset=0",
	} {
		status, body := requestDemo(t, server, path, "")
		if status != http.StatusUnauthorized || strings.Contains(body, "disk-entry-") {
			t.Fatalf("anonymous page leaked Case: %d %s", status, body)
		}
	}
	status, body := requestDemo(t, server, "/case-timeline?case_id=case-b&limit=200", testBearer)
	if status != http.StatusNotFound || strings.Contains(body, "disk-entry-") {
		t.Fatalf("cross-tenant Case leaked: %d %s", status, body)
	}
	status, body = requestDemo(t, server, "/demo-source-open?case_id=case-a&entry_id=disk-entry-000", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "demo-evidence-a") {
		t.Fatalf("disk-sourced typed preview was not reauthorized: %d %s", status, body)
	}
}
