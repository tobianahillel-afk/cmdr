package investigationtimeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type caseDiscoveryCatalogFunc func(context.Context, Scope) ([]string, error)

func (f caseDiscoveryCatalogFunc) ListCandidateCaseIDs(ctx context.Context, scope Scope) ([]string, error) {
	return f(ctx, scope)
}

func caseDiscoveryGet(h http.Handler, bearer string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/cases", nil)
	if bearer != "" {
		req.Header.Set("Authorization", bearer)
	}
	return rr(h, req)
}

func TestMemoryCaseDiscoveryOpensSameAuthorizedTimeline(t *testing.T) {
	data := memoryDemoFixture()
	// A Case is visible to another same-tenant principal, not to analyst-a.
	data.Cases = append(data.Cases, MemoryTimelineCase{
		ID: "restricted-case", Tenant: "tenant-a", Environment: "prod",
		Readers: []string{"analyst-secret"},
	})
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, data)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		bearer string
		want   []string
	}{
		{"Bearer aa", []string{"case-a"}},
		{"Bearer ab", []string{"case-a", "restricted-case"}},
		{"Bearer bb", []string{"case-b"}},
	} {
		resp := caseDiscoveryGet(handlers.Cases, tc.bearer)
		if resp.Code != http.StatusOK || resp.Header().Get("Cache-Control") != "no-store" ||
			resp.Header().Get("Vary") != "Authorization" {
			t.Fatalf("catalog HTTP=%d for %q: %s", resp.Code, tc.bearer, resp.Body.String())
		}
		var listed struct {
			Cases []string `json:"cases"`
		}
		if err := json.Unmarshal(resp.Body.Bytes(), &listed); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(listed.Cases, tc.want) {
			t.Fatalf("unauthorized or unordered catalog for %q: %v want %v", tc.bearer, listed.Cases, tc.want)
		}
		for _, id := range tc.want {
			got := memoryRequest(handlers.Timeline, id, tc.bearer)
			if got.Code != http.StatusOK {
				t.Fatalf("discovered Case %q not openable for %q: %d", id, tc.bearer, got.Code)
			}
		}
	}
	for _, bearer := range []string{"", "Bearer bad"} {
		resp := caseDiscoveryGet(handlers.Cases, bearer)
		if resp.Code != http.StatusUnauthorized || strings.Contains(resp.Body.String(), "case-a") {
			t.Fatalf("unauthenticated catalog leaked: HTTP %d", resp.Code)
		}
	}
	if resp := memoryRequest(handlers.Timeline, "restricted-case", "Bearer aa"); resp.Code != http.StatusNotFound {
		t.Fatalf("discovery permissions not consistent with timeline HTTP: %d", resp.Code)
	}
}

func TestCaseDiscoveryRejectsInvalidRequestAndCatalogFailure(t *testing.T) {
	catalog := caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
		return []string{"case-one"}, nil
	})
	store, _ := workspaceFixture()
	h, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store, catalog)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method, rawQuery, body string
		code                   int
	}{
		{http.MethodPost, "", "", http.StatusMethodNotAllowed},
		{http.MethodGet, "tenant=tenant-b", "", http.StatusBadRequest},
		{http.MethodGet, "", "unexpected", http.StatusBadRequest},
	} {
		req := httptest.NewRequest(tc.method, "/cases", strings.NewReader(tc.body))
		req.URL.RawQuery = tc.rawQuery
		req.Header.Set("Authorization", "Bearer analyst-a")
		got := rr(h, req)
		if got.Code != tc.code {
			t.Fatalf("%s query %q: HTTP %d", tc.method, tc.rawQuery, got.Code)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/cases", nil)
	req.Header.Set("Authorization", "Bearer analyst-a")
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	got := rr(h, req.WithContext(ctx))
	if got.Code != http.StatusServiceUnavailable {
		t.Fatalf("cancelled discovery opened a catalog: %d", got.Code)
	}
	fault := caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
		return nil, errors.New("owner failed")
	})
	h, err = NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store, fault)
	if err != nil {
		t.Fatal(err)
	}
	got = caseDiscoveryGet(h, "Bearer analyst-a")
	if got.Code != http.StatusNotFound || strings.Contains(got.Body.String(), "owner failed") {
		t.Fatalf("provider failure disclosed sensitive internal error: %d %s", got.Code, got.Body.String())
	}
	for _, malformed := range [][]string{
		{"case-one", "case-one"},
		{""},
		{"*"},
		{strings.Repeat("x", maxHTTPCaseIDBytes+1)},
	} {
		h, err = NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store,
			caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
				return malformed, nil
			}))
		if err != nil {
		t.Fatal(err)
	}
		got = caseDiscoveryGet(h, "Bearer analyst-a")
		if got.Code != http.StatusNotFound || strings.Contains(got.Body.String(), "case-one") {
			t.Fatalf("malformed Case catalog not denied: %#v -> %d %s", malformed, got.Code, got.Body.String())
		}
	}
	if _, err := NewCaseDiscoveryHTTPHandler(nil, store, catalog); err == nil {
		t.Fatal("missing session service accepted")
	}
	if _, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, nil, catalog); err == nil {
		t.Fatal("missing access authority accepted")
	}
}

func TestCaseDiscoveryFiltersDeniedCandidateBeforeResponse(t *testing.T) {
	store, _ := workspaceFixture()
	h, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store,
		caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
			return []string{"case-two", "case-one"}, nil
		}))
	if err != nil {
		t.Fatal(err)
	}
	r := caseDiscoveryGet(h, "Bearer analyst-a")
	if r.Code != http.StatusOK || strings.Contains(r.Body.String(), "case-two") {
		t.Fatalf("unreadable Case leaked: %d %s", r.Code, r.Body.String())
	}
	var got struct {
		Cases []string `json:"cases"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Cases, []string{"case-one"}) {
		t.Fatalf("Case access decision ignored: %v", got.Cases)
	}
}
