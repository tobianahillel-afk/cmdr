package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func authorizedTestTimeline(t *testing.T) http.Handler {
	t.Helper()
	handlers, err := timeline.NewMemoryCaseTimelineHTTPHandlers(func(*http.Request) (string, error) {
		return "analyst-a", nil
	}, testOwner())
	if err != nil {
		t.Fatal(err)
	}
	return handlers.Timeline
}

// Simulate a real previously authorized source losing permission between
// Search Job discovery and execution. A source label is never an ACL.
func TestCaseSearchBackendRevocationAndTrustBoundaries(t *testing.T) {
	handler := authorizedTestTimeline(t)
	from, _ := time.Parse(time.RFC3339, "2026-10-01T10:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2026-10-01T10:15:00Z")
	req := search.BackendRequest{
		TenantRef: "tenant-a", EnvironmentRef: "lab",
		TimeStart: from, TimeEnd: to, Sources: []string{"telemetry-early"},
	}
	backend := &caseSearchBackend{
		owner: handler, caseID: "case-a", tenant: "tenant-a",
		environment: "lab", query: "sensor-dns",
		sourceKinds: map[string]string{
			"telemetry-early": "telemetry-event",
			"telemetry-late":  "telemetry-event",
			"private-secret":  "telemetry-event",
			"evidence-1":      "evidence",
		},
	}
	result, err := backend.Execute(context.Background(), req)
	if err != nil || len(result.ResultRefs) != 1 ||
		result.ResultRefs[0].Ref != "allowed-early" {
		t.Fatalf("authorized owner source search failed: %+v %v", result, err)
	}
	// The owner now refuses the entire Case. The previously recognized
	// source kind must not make stale Search Job privileges sufficient.
	backend.owner = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	if result, err := backend.Execute(context.Background(), req); err == nil ||
		len(result.ResultRefs) != 0 {
		t.Fatalf("revoked owner emitted results: %+v %v", result, err)
	}
	backend.owner = handler
	for _, change := range []struct {
		name   string
		mutate func(*search.BackendRequest)
	}{
		{"foreign tenant", func(r *search.BackendRequest) { r.TenantRef = "tenant-b" }},
		{"wrong environment", func(r *search.BackendRequest) { r.EnvironmentRef = "prod" }},
		{"missing interval", func(r *search.BackendRequest) { r.TimeStart = time.Time{} }},
		{"inverted interval", func(r *search.BackendRequest) { r.TimeEnd = r.TimeStart }},
		{"denied source", func(r *search.BackendRequest) { r.Sources = []string{"private-secret"} }},
		{"wrong source kind", func(r *search.BackendRequest) { r.Sources = []string{"evidence-1"} }},
		{"duplicate source", func(r *search.BackendRequest) {
			r.Sources = []string{"telemetry-early", "telemetry-early"}
		}},
	} {
		t.Run(change.name, func(t *testing.T) {
			input := req
			change.mutate(&input)
			got, err := backend.Execute(context.Background(), input)
			if err == nil || len(got.ResultRefs) != 0 {
				t.Fatalf("untrusted search scope yielded data: %+v %v", got, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := backend.Execute(ctx, req); err == nil {
		t.Fatalf("cancelled search must fail closed: %v", err)
	}
}

// Invalid/untrusted HTTP owner page data must never become a partial
// search result, even if a caller has a synthetic Case identifier.
func TestCaseSearchRejectsInconsistentOwnerPages(t *testing.T) {
	handler := authorizedTestTimeline(t)
	request := httptest.NewRequest(http.MethodGet,
		"/case-timeline?case_id=case-a&limit=200&offset=0", nil)
	original := httptest.NewRecorder()
	handler.ServeHTTP(original, request)
	if original.Code != http.StatusOK {
		t.Fatalf("owner fixture unavailable: %d", original.Code)
	}
	var page ownerPage
	if err := json.Unmarshal(original.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if len(page.Entries) < 2 {
		t.Fatalf("need multiple visible owner rows: %+v", page)
	}
	for _, tc := range []struct {
		name   string
		change func(*ownerPage)
	}{
		{"untrusted offset", func(p *ownerPage) { p.Page.Offset = 200 }},
		{"untrusted total", func(p *ownerPage) { p.Page.TotalVisible++ }},
		{"missing owner row", func(p *ownerPage) { p.Entries = nil }},
		{"duplicate owner row", func(p *ownerPage) { p.Entries[1] = p.Entries[0] }},
		{"missing source identity", func(p *ownerPage) { p.Entries[0].SourceRef = "" }},
		{"missing provenance", func(p *ownerPage) { p.Entries[0].Provenance = "" }},
		{"invalid entry version", func(p *ownerPage) { p.Entries[0].Version = 0 }},
		{"invalid response status", func(p *ownerPage) { p.Status = "not-authorized" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			copyPage := page
			copyPage.Entries = append([]timeline.ProjectedEntry(nil), page.Entries...)
			tc.change(&copyPage)
			data, err := json.Marshal(copyPage)
			if err != nil {
				t.Fatal(err)
			}
			untrusted := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(data)
			})
			got, err := readOwnerPages(context.Background(), untrusted, "case-a")
			if err == nil || len(got) != 0 {
				t.Fatalf("inconsistent owner page was accepted: %+v %v", got, err)
			}
		})
	}
	for _, response := range []struct {
		name string
		code int
		body string
	}{
		{"denied", http.StatusNotFound, ""},
		{"malformed", http.StatusOK, "{bad-json"},
		{"oversized", http.StatusOK, strings.Repeat("x", 8<<20+1)},
	} {
		t.Run(response.name, func(t *testing.T) {
			owner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(response.code)
				_, _ = w.Write([]byte(response.body))
			})
			got, err := readOwnerPages(context.Background(), owner, "case-a")
			if err == nil || len(got) != 0 {
				t.Fatalf("invalid owner response was accepted: %+v %v", got, err)
			}
		})
	}
}

// The CLI's file path and JSON parser are security boundaries. These
// checks use real operator-style disk files rather than mock decoders.
func TestCaseSearchOwnerSnapshotRejectsUnsafeFiles(t *testing.T) {
	directory := t.TempDir()
	good, err := json.Marshal(testOwner())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"unsupported key", append(append([]byte(nil), good[:len(good)-1]...), []byte(",\"unexpected\":1}")...)},
		{"second document", append(append([]byte(nil), good...), good...)},
		{"bad syntax", []byte("{bad-json")},
		{"too large", bytes.Repeat([]byte("x"), 8<<20+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(directory, strings.ReplaceAll(tc.name, " ", "-")+".json")
			if err := os.WriteFile(path, tc.data, 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := readSnapshot(path); err == nil {
				t.Fatal("invalid local source file was accepted")
			}
		})
	}
	validPath := filepath.Join(directory, "valid-owner.json")
	if err := os.WriteFile(validPath, good, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readSnapshot("relative.json"); err == nil {
		t.Fatal("relative owner paths must be refused")
	}
	if _, err := readSnapshot(directory); err == nil {
		t.Fatal("directory is not a regular owner snapshot")
	}
	link := filepath.Join(directory, "owner-link.json")
	if err := os.Symlink(validPath, link); err == nil {
		if _, err := readSnapshot(link); err == nil {
			t.Fatal("symlink owner path accepted")
		}
	}
	args := []string{
		"-snapshot", validPath, "-principal", "analyst-a", "-case", "case-a",
		"-query", "sensor", "-from", "2026-10-01T09:00:00Z", "-to", "2026-10-01T11:00:00Z",
	}
	for _, bad := range [][]string{
		{"-unknown-option"},
		args[:len(args)-1],
		append(append([]string(nil), args...), "unexpected"),
	} {
		var output bytes.Buffer
		if err := runCLI(context.Background(), bad, &output); err == nil || output.Len() != 0 {
			t.Fatalf("CLI accepted malformed selection: %q %v", bad, err)
		}
	}
}
