package investigationtimeline

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Exercise the host-owned interfaces without using the memory demo constructor.
// This is a real GET Case discovery -> ACL -> Timeline projection flow with
// independently injected owner services; not a deployed identity/persistence layer.
func TestHostInjectedInvestigationDiscoveryOpensAuthorizedTimeline(t *testing.T) {
	store, _ := workspaceFixture()
	candidates := caseDiscoveryCatalogFunc(func(_ context.Context, scope Scope) ([]string, error) {
		if scope.Tenant != "tenant-a" || scope.Environment != "prod" {
			t.Fatalf("untrusted scope at Case owner: %+v", scope)
		}
		return []string{"case-two", "case-one"}, nil
	})
	handlers, err := NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store, candidates, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/cases", handlers.Cases)
	mux.Handle("/case-timeline", handlers.Timeline)

	get := func(path, credential string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if credential != "" {
			req.Header.Set("Authorization", credential)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}
	listed := get("/cases", "Bearer analyst-a")
	var catalog struct {
		Cases []string `json:"cases"`
	}
	if listed.Code != http.StatusOK || json.Unmarshal(listed.Body.Bytes(), &catalog) != nil ||
		!reflect.DeepEqual(catalog.Cases, []string{"case-one"}) {
		t.Fatalf("host Case discovery: HTTP %d %s", listed.Code, listed.Body.String())
	}
	page := get("/case-timeline?case_id="+catalog.Cases[0]+"&kind=observed&q=entry-new&limit=1&offset=0",
		"Bearer analyst-a")
	var result timelinePagedHTTPResult
	if page.Code != http.StatusOK || json.Unmarshal(page.Body.Bytes(), &result) != nil ||
		len(result.Entries) != 1 || result.Entries[0].ID != "entry-new" ||
		result.Entries[0].Kind != "observed" || result.Page.TotalVisible != 1 {
		t.Fatalf("host Timeline projection: HTTP %d %s", page.Code, page.Body.String())
	}
	for _, secret := range []string{"entry-hidden", "src-hidden", "evidence-secret", "case-two", "tenant-b"} {
		if strings.Contains(page.Body.String(), secret) || strings.Contains(listed.Body.String(), secret) {
			t.Fatalf("host discovery/Timeline leaked denied data: %q", secret)
		}
	}
	for _, path := range []string{"/cases", "/case-timeline?case_id=case-one"} {
		if response := get(path, ""); response.Code != http.StatusUnauthorized {
			t.Fatalf("anonymous host route %s: HTTP %d", path, response.Code)
		}
	}
	if resp := get("/case-timeline?case_id=case-two", "Bearer analyst-a"); resp.Code != http.StatusNotFound {
		t.Fatalf("hidden Case opened: %d", resp.Code)
	}
	if resp := get("/cases?tenant=tenant-b", "Bearer analyst-a"); resp.Code != http.StatusBadRequest {
		t.Fatalf("browser tenant selector accepted: %d", resp.Code)
	}
	// Discovery cannot become a capability token: access is rechecked at
	// Timeline open after any owner-side revocation.
	store.revokeAfterList = true
	if resp := get("/case-timeline?case_id=case-one", "Bearer analyst-a"); resp.Code != http.StatusNotFound ||
		strings.Contains(resp.Body.String(), "entry-new") {
		t.Fatalf("Case revocation bypassed at Timeline open: %d %s", resp.Code, resp.Body.String())
	}
}

func TestHostInjectedInvestigationRejectsMissingOwnerServices(t *testing.T) {
	store, _ := workspaceFixture()
	candidate := caseDiscoveryCatalogFunc(func(_ context.Context, _ Scope) ([]string, error) { return nil, nil })
	for _, tc := range []struct {
		name  string
		build func() (CaseInvestigationHTTPHandlers, error)
	}{
		{"session", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(nil, store, candidate, store, store, store)
		}},
		{"case-access", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, nil, candidate, store, store, store)
		}},
		{"case-list", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store, nil, store, store, store)
		}},
		{"timeline-catalog", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store, candidate, nil, store, store)
		}},
		{"entry-owner", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store, candidate, store, nil, store)
		}},
		{"source-owner", func() (CaseInvestigationHTTPHandlers, error) {
			return NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store, candidate, store, store, nil)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, err := tc.build()
			if err == nil || h.Cases != nil || h.Timeline != nil {
				t.Fatalf("missing %s accepted: %#v, %v", tc.name, h, err)
			}
		})
	}
}

func TestHostCaseDiscoveryTimelineToTypedSourceNavigation(t *testing.T) {
	store, _ := workspaceFixture()
	candidates := caseDiscoveryCatalogFunc(func(_ context.Context, _ Scope) ([]string, error) {
		return []string{"case-one", "case-two"}, nil
	})
	h, err := NewCaseInvestigationHTTPHandlers(timelineHTTPSession{}, store,
		candidates, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/cases", h.Cases)
	mux.Handle("/case-timeline", h.Timeline)
	mux.Handle("/source-navigation", h.Source)
	get := func(path string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		r.Header.Set("Authorization", "Bearer analyst-a")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	list := get("/cases")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "case-one") ||
		strings.Contains(list.Body.String(), "case-two") {
		t.Fatalf("host Case picker bypassed tenant permission: %s", list.Body.String())
	}
	page := get("/case-timeline?case_id=case-one&kind=observed&limit=1")
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), "entry-new") {
		t.Fatalf("host authorized Timeline unavailable: %s", page.Body.String())
	}
	navigation := get("/source-navigation?case_id=case-one&entry_id=entry-new")
	var result SourceNavigationResult
	if navigation.Code != http.StatusOK || json.Unmarshal(navigation.Body.Bytes(), &result) != nil ||
		result.Status != "ready" || result.Draft == nil ||
		result.Draft.OwnerKind != "telemetry-event" || result.Draft.OwnerObjectID != "event-1" ||
		result.Draft.Return.EntryID != "entry-new" {
		t.Fatalf("Case-to-Timeline-to-source handoff failed: %d %s", navigation.Code, navigation.Body.String())
	}
	for _, path := range []string{
		"/source-navigation?case_id=case-one&entry_id=entry-hidden",
		"/source-navigation?case_id=case-two&entry_id=entry-other",
	} {
		response := get(path)
		if response.Code != http.StatusNotFound ||
			strings.Contains(response.Body.String(), "src-hidden") ||
			strings.Contains(response.Body.String(), "evidence-secret") {
			t.Fatalf("host source navigation leaked protected owner metadata: %s", response.Body.String())
		}
	}
}
