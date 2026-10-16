package investigationtimeline

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func memorySourcePreviewFixture() MemoryTimelineData {
	return MemoryTimelineData{
		Principals: []MemoryTimelinePrincipal{
			{ID: "analyst-a", Tenant: "tenant-a", Environment: "lab"},
			{ID: "analyst-b", Tenant: "tenant-b", Environment: "lab"},
		},
		Cases: []MemoryTimelineCase{
			{ID: "case-a", Tenant: "tenant-a", Environment: "lab", Readers: []string{"analyst-a"}, EntryIDs: []string{"entry-a", "hidden-a"}},
			{ID: "case-b", Tenant: "tenant-b", Environment: "lab", Readers: []string{"analyst-b"}, EntryIDs: []string{"entry-b"}},
		},
		Entries: []MemoryTimelineEntry{
			{CaseID: "case-a", Value: Entry{ID: "entry-a", Tenant: "tenant-a", Kind: "observed", SourceRef: "src-a", Provenance: "synthetic-a", SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00", Version: 1}, Readers: []string{"analyst-a"}},
			{CaseID: "case-a", Value: Entry{ID: "hidden-a", Tenant: "tenant-a", Kind: "observed", SourceRef: "src-hidden", Provenance: "synthetic-hidden", SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00", Version: 1}, Readers: []string{"analyst-a"}},
			{CaseID: "case-b", Value: Entry{ID: "entry-b", Tenant: "tenant-b", Kind: "inferred", SourceRef: "src-b", Provenance: "synthetic-b", SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00", Version: 1}, Readers: []string{"analyst-b"}},
		},
		Sources: []MemoryTimelineSource{
			{Identity: SourceOwnerIdentity{Reference: "src-a", Kind: "telemetry-event", ObjectID: "object-a", Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"}, Readers: []string{"analyst-a"}, Preview: "Authorized synthetic DNS record"},
			{Identity: SourceOwnerIdentity{Reference: "src-hidden", Kind: "evidence", ObjectID: "secret-a", Tenant: "tenant-a", Environment: "lab", CaseID: "case-a"}, Preview: "DO NOT DISPLAY"},
			{Identity: SourceOwnerIdentity{Reference: "src-b", Kind: "incident", ObjectID: "object-b", Tenant: "tenant-b", Environment: "lab", CaseID: "case-b"}, Readers: []string{"analyst-b"}, Preview: "Tenant B private synthetic incident"},
		},
	}
}

func TestMemorySourcePreviewOpeningRechecksOwnerAndTenant(t *testing.T) {
	data := memorySourcePreviewFixture()
	for _, tc := range []struct {
		name, principal, path string
		status                int
		preview               string
	}{
		{"allowed", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=entry-a", 200, "Authorized synthetic DNS record"},
		{"other-tenant-own", "analyst-b", "/demo-source-open?case_id=case-b&entry_id=entry-b", 200, "Tenant B private synthetic incident"},
		{"cross-tenant", "analyst-a", "/demo-source-open?case_id=case-b&entry_id=entry-b", 404, ""},
		{"denied-source", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=hidden-a", 404, ""},
		{"wrong-case", "analyst-a", "/demo-source-open?case_id=case-b&entry_id=entry-a", 404, ""},
		{"not-present", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=not-found", 404, ""},
		{"injected-object", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=entry-a&object_id=object-b", 400, ""},
		{"injected-tenant", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=entry-a&tenant=tenant-b", 400, ""},
		{"duplicate-entry", "analyst-a", "/demo-source-open?case_id=case-a&entry_id=entry-a&entry_id=entry-b", 400, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bound, err := NewMemoryCaseTimelineHTTPHandlers(func(*http.Request) (string, error) { return tc.principal, nil }, data)
			if err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()
			bound.OpenSource.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status %d want %d: %s", w.Code, tc.status, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Vary") != "Authorization" {
				t.Fatal("source preview was cacheable")
			}
			if tc.status != 200 {
				for _, forbidden := range []string{"secret-a", "src-hidden", "DO NOT DISPLAY", "object-b", "tenant-b"} {
					if strings.Contains(w.Body.String(), forbidden) {
						t.Fatalf("denial disclosed %s", forbidden)
					}
				}
				return
			}
			var payload struct {
				Status, CaseID, EntryID, Kind, ObjectID, Preview string
				Synthetic                                        bool
			}
			if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.Status != "ready" || !payload.Synthetic || payload.Preview != tc.preview {
				t.Fatalf("incorrect authorized preview: %+v", payload)
			}
			if strings.Contains(w.Body.String(), "DO NOT DISPLAY") {
				t.Fatal("hidden source disclosed")
			}
		})
	}
	bound, err := NewMemoryCaseTimelineHTTPHandlers(func(*http.Request) (string, error) { return "", nil }, data)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	bound.OpenSource.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=entry-a", nil))
	if w.Code != 401 {
		t.Fatalf("anonymous access passed: %d", w.Code)
	}
	w = httptest.NewRecorder()
	bound.OpenSource.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/demo-source-open?case_id=case-a&entry_id=entry-a", nil))
	if w.Code != 405 {
		t.Fatalf("write accepted: %d", w.Code)
	}
}

func TestMemorySourcePreviewIsBoundedAtConstruction(t *testing.T) {
	store, _ := workspaceFixture()
	_ = store // prove no change to the owner-generated runtime fixture.
	data := MemoryTimelineData{
		Principals: []MemoryTimelinePrincipal{{ID: "analyst", Tenant: "t", Environment: "lab"}},
		Cases:      []MemoryTimelineCase{{ID: "case", Tenant: "t", Environment: "lab", Readers: []string{"analyst"}, EntryIDs: []string{"entry"}}},
		Entries:    []MemoryTimelineEntry{{CaseID: "case", Value: Entry{ID: "entry", Tenant: "t", Kind: "observed", SourceRef: "ref", Provenance: "p", Version: 1, SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00"}, Readers: []string{"analyst"}}},
		Sources:    []MemoryTimelineSource{{Identity: SourceOwnerIdentity{Reference: "ref", Kind: "evidence", ObjectID: "object", Tenant: "t", Environment: "lab", CaseID: "case"}, Readers: []string{"analyst"}, Preview: strings.Repeat("x", 2049)}},
	}
	if _, err := NewMemoryCaseTimelineHTTPHandlers(func(*http.Request) (string, error) { return "analyst", nil }, data); err == nil {
		t.Fatal("oversized source preview accepted")
	}
}

func TestMemorySourcePreviewRejectsMalformedOrCancelledInput(t *testing.T) {
	bound, err := NewMemoryCaseTimelineHTTPHandlers(
		func(*http.Request) (string, error) { return "analyst-a", nil }, memorySourcePreviewFixture())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, query, body string
		cancel, transfer  bool
		want              int
	}{
		{"get-body", "case_id=case-a&entry_id=entry-a", "payload", false, false, 400},
		{"transfer-encoding", "case_id=case-a&entry_id=entry-a", "", false, true, 400},
		{"invalid-percent-escape", "case_id=case-a&entry_id=%XX", "", false, false, 400},
		{"invalid-utf8", "case_id=case-a&entry_id=%FF", "", false, false, 400},
		{"wildcard-case", "case_id=*&entry_id=entry-a", "", false, false, 400},
		{"wildcard-entry", "case_id=case-a&entry_id=*", "", false, false, 400},
		{"empty-case", "case_id=&entry_id=entry-a", "", false, false, 400},
		{"missing-entry", "case_id=case-a", "", false, false, 400},
		{"duplicate-case", "case_id=case-a&case_id=case-b&entry_id=entry-a", "", false, false, 400},
		{"oversized-entry", "case_id=case-a&entry_id=" + strings.Repeat("x", 1025), "", false, false, 400},
		{"oversized-query", "case_id=case-a&entry_id=" + strings.Repeat("x", 4097), "", false, false, 400},
		{"cancelled-request", "case_id=case-a&entry_id=entry-a", "", true, false, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			if tc.body != "" {
				body = strings.NewReader(tc.body)
			}
			r := httptest.NewRequest(http.MethodGet, "/demo-source-open?"+tc.query, body)
			if tc.transfer {
				r.TransferEncoding = []string{"chunked"}
			}
			if tc.cancel {
				ctx, cancel := context.WithCancel(r.Context())
				cancel()
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()
			bound.OpenSource.ServeHTTP(w, r)
			if w.Code != tc.want || strings.Contains(w.Body.String(), "Authorized synthetic DNS record") ||
				strings.Contains(w.Body.String(), "object-a") {
				t.Fatalf("invalid request was not rejected safely: status %d body %q", w.Code, w.Body.String())
			}
		})
	}
}
