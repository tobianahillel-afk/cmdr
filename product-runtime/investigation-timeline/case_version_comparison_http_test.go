package investigationtimeline

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCaseVersionComparisonReusesAuthoritativeOwnerChain(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/case-version-comparison?case_id=case-a&entry_id=new", nil)
	req.Header.Set("Authorization", "Bearer aa")
	got := rr(handlers.Compare, req)
	if got.Code != http.StatusOK {
		t.Fatalf("authorized comparison failed: %d %s", got.Code, got.Body.String())
	}
	var result caseVersionComparison
	if err := json.Unmarshal(got.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "ready" || result.IsAuthority ||
		result.Current.ID != "new" || result.Current.Kind != "observed" ||
		result.Current.Version != 2 || result.Current.Supersedes != "old" ||
		result.Previous.ID != "old" || result.Previous.Kind != "inferred" ||
		result.Previous.Version != 1 {
		t.Fatalf("wrong Case/Timeline version lineage: %+v", result)
	}
	if got.Header().Get("Cache-Control") != "no-store" ||
		got.Header().Get("Referrer-Policy") != "no-referrer" ||
		got.Header().Get("Vary") != "Authorization" {
		t.Fatal("sensitive response missing non-cache privacy headers")
	}
	for _, secret := range []string{"src-secret", "evidence-secret", "hidden", "tenant-b"} {
		if strings.Contains(got.Body.String(), secret) {
			t.Fatalf("hidden source leaked: %s", secret)
		}
	}
}

func TestCaseVersionComparisonDenyAtHTTPBoundary(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	valid := "/case-version-comparison?case_id=case-a&entry_id=new"
	cases := []struct {
		name, method, path, bearer, body string
		status                          int
	}{
		{"anonymous", "GET", valid, "", "", http.StatusUnauthorized},
		{"bad-token", "GET", valid, "Bearer wrong", "", http.StatusUnauthorized},
		{"foreign-case", "GET", "/case-version-comparison?case_id=case-b&entry_id=foreign", "Bearer aa", "", http.StatusNotFound},
		{"foreign-entry", "GET", "/case-version-comparison?case_id=case-a&entry_id=foreign", "Bearer aa", "", http.StatusNotFound},
		{"not-superseding", "GET", "/case-version-comparison?case_id=case-a&entry_id=old", "Bearer aa", "", http.StatusNotFound},
		{"missing-entry", "GET", "/case-version-comparison?case_id=case-a&entry_id=missing", "Bearer aa", "", http.StatusNotFound},
		{"POST", "POST", valid, "Bearer aa", "", http.StatusMethodNotAllowed},
		{"body", "GET", valid, "Bearer aa", "not-a-read", http.StatusBadRequest},
		{"no-case", "GET", "/case-version-comparison?entry_id=new", "Bearer aa", "", http.StatusBadRequest},
		{"no-entry", "GET", "/case-version-comparison?case_id=case-a", "Bearer aa", "", http.StatusBadRequest},
		{"duplicate", "GET", valid+"&entry_id=old", "Bearer aa", "", http.StatusBadRequest},
		{"tenant-injection", "GET", valid+"&tenant=tenant-b", "Bearer aa", "", http.StatusBadRequest},
		{"permission-injection", "GET", valid+"&permission=allow", "Bearer aa", "", http.StatusBadRequest},
		{"empty", "GET", "/case-version-comparison?case_id=case-a&entry_id=", "Bearer aa", "", http.StatusBadRequest},
		{"wildcard", "GET", "/case-version-comparison?case_id=case-a&entry_id=%2A", "Bearer aa", "", http.StatusBadRequest},
		{"encoded-null", "GET", "/case-version-comparison?case_id=case-a&entry_id=%00", "Bearer aa", "", http.StatusBadRequest},
		{"too-long", "GET", "/case-version-comparison?case_id=case-a&entry_id="+strings.Repeat("x", maxHTTPCaseIDBytes+1), "Bearer aa", "", http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var body *strings.Reader
			var request *http.Request
			if tc.body != "" {
				body = strings.NewReader(tc.body)
				request = httptest.NewRequest(tc.method, tc.path, body)
			} else {
				request = httptest.NewRequest(tc.method, tc.path, nil)
			}
			if tc.bearer != "" {
				request.Header.Set("Authorization", tc.bearer)
			}
			got := rr(handlers.Compare, request)
			if got.Code != tc.status {
				t.Fatalf("HTTP %d, want %d: %s", got.Code, tc.status, got.Body.String())
			}
			for _, secret := range []string{"src-secret", "evidence-secret", "tenant-b", "event-a"} {
				if strings.Contains(got.Body.String(), secret) {
					t.Fatalf("denied comparison disclosed source %s", secret)
				}
			}
		})
	}
	for _, broken := range []string{"case_id=%GG&entry_id=new", "case_id=case-a&entry_id=new&entry_id=old"} {
		request := httptest.NewRequest(http.MethodGet, "/case-version-comparison", nil)
		request.URL.RawQuery = broken
		request.Header.Set("Authorization", "Bearer aa")
		if got := rr(handlers.Compare, request); got.Code != http.StatusBadRequest {
			t.Fatalf("malformed query accepted: HTTP %d", got.Code)
		}
	}
}

func TestCaseVersionComparisonCancelsAndDeniesInvisiblePredecessor(t *testing.T) {
	for _, mutation := range []struct {
		name string
		change func(*MemoryTimelineData)
	}{
		{"predecessor-source-revoked", func(d *MemoryTimelineData) { d.Sources[1].Readers = nil }},
		{"current-source-revoked", func(d *MemoryTimelineData) { d.Sources[0].Readers = nil }},
		{"predecessor-entry-revoked", func(d *MemoryTimelineData) { d.Entries[1].Readers = nil }},
		{"other-tenant", func(d *MemoryTimelineData) { d.Entries[1].Readers = []string{"analyst-secret"} }},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			candidate := memoryDemoFixture()
			mutation.change(&candidate)
			h, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, candidate)
			if err != nil { t.Fatal(err) }
			req := httptest.NewRequest(http.MethodGet, "/case-version-comparison?case_id=case-a&entry_id=new", nil)
			req.Header.Set("Authorization", "Bearer aa")
			got := rr(h.Compare, req)
			if got.Code != http.StatusNotFound || strings.Contains(got.Body.String(), "old") ||
				strings.Contains(got.Body.String(), "src-old") {
				t.Fatalf("masked predecessor comparison leaked: %d %s", got.Code, got.Body.String())
			}
		})
	}
	h, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memoryDemoFixture())
	if err != nil { t.Fatal(err) }
	request := httptest.NewRequest(http.MethodGet, "/case-version-comparison?case_id=case-a&entry_id=new", nil)
	request.Header.Set("Authorization", "Bearer aa")
	ctx, cancel := context.WithCancel(request.Context())
	cancel()
	got := rr(h.Compare, request.WithContext(ctx))
	if got.Code != http.StatusServiceUnavailable {
		t.Fatalf("cancelled comparison accessed owners: %d", got.Code)
	}
	request = httptest.NewRequest(http.MethodGet, "/case-version-comparison?case_id=case-a&entry_id=new", nil)
	request.Header.Set("Authorization", "Bearer aa")
	request.TransferEncoding = []string{"chunked"}
	if got := rr(h.Compare, request); got.Code != http.StatusBadRequest {
		t.Fatalf("transfer encoding accepted: %d", got.Code)
	}
	request = httptest.NewRequest(http.MethodGet, "/case-version-comparison?case_id=case-a&entry_id=new", nil)
	request.URL = nil
	if got := rr(h.Compare, request); got.Code != http.StatusBadRequest {
		t.Fatalf("nil URL accepted: %d", got.Code)
	}
}

func TestCaseVersionComparisonConstructorRequiresAuthority(t *testing.T) {
	store, _ := workspaceFixture()
	for _, h := range []func() error{
		func() error { _, err := NewCaseVersionComparisonHTTPHandler(nil, store, store, store, store); return err },
		func() error { _, err := NewCaseVersionComparisonHTTPHandler(timelineHTTPSession{}, nil, store, store, store); return err },
		func() error { _, err := NewCaseVersionComparisonHTTPHandler(timelineHTTPSession{}, store, nil, store, store); return err },
		func() error { _, err := NewCaseVersionComparisonHTTPHandler(timelineHTTPSession{}, store, store, nil, store); return err },
		func() error { _, err := NewCaseVersionComparisonHTTPHandler(timelineHTTPSession{}, store, store, store, nil); return err },
	} {
		if err := h(); err == nil { t.Fatal("accepted missing owner authority") }
	}
}
