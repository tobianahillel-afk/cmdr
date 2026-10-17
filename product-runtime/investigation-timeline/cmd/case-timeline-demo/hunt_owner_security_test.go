package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Covers the actual HTTP boundary traversed by a direct cross-Case event
// shortcut, including owner reauthorization and malformed client selectors.
// The browser never supplies the owner identity or a source permission flag.
func TestCrossCaseEventPivotReauthorizesCaseAndTypedSource(t *testing.T) {
	h, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(h)
	defer srv.Close()
	status, body := requestDemo(t, srv, "/case-timeline?case_id=case-a&q=event-new&limit=20&offset=0", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "\"id\":\"event-new\"") ||
		strings.Contains(body, "demo-hidden-source") {
		t.Fatalf("exact event search did not authorize a real Case result: %d %s", status, body)
	}
	path := "/case-source-navigation?" + url.Values{
		"case_id": {"case-a"}, "entry_id": {"event-new"},
		"q": {"event-new"}, "limit": {"20"}, "offset": {"0"},
	}.Encode()
	status, body = requestDemo(t, srv, path, testBearer)
	if status != http.StatusOK || !strings.Contains(body, "\"status\":\"ready\"") ||
		!strings.Contains(body, "event-new") {
		t.Fatalf("exact event could not reauthorize navigation: %d %s", status, body)
	}
	status, body = requestDemo(t, srv, "/demo-source-open?case_id=case-a&entry_id=event-new", testBearer)
	if status != http.StatusOK || !strings.Contains(body, "\"synthetic\":true") {
		t.Fatalf("owner did not reauthorize the exact event source: %d %s", status, body)
	}
	for _, path := range []string{
		"/case-source-navigation?case_id=case-b&entry_id=event-new",
		"/demo-source-open?case_id=case-b&entry_id=event-new",
		"/demo-source-open?case_id=case-a&entry_id=restricted",
	} {
		status, body = requestDemo(t, srv, path, testBearer)
		if status != http.StatusNotFound || strings.Contains(body, "demo-telemetry-a") ||
			strings.Contains(body, "demo-hidden-source") {
			t.Fatalf("wrong Case or denied entry borrowed a source permission: %s: %d %s", path, status, body)
		}
	}
	for _, path := range []string{
		"/case-timeline?case_id=case-a&q=event-new",
		"/case-source-navigation?case_id=case-a&entry_id=event-new",
		"/demo-source-open?case_id=case-a&entry_id=event-new",
	} {
		status, body = requestDemo(t, srv, path, "")
		if status != http.StatusUnauthorized || strings.Contains(body, "event-new") {
			t.Fatalf("anonymous pivot disclosed owner identity: %s: %d %s", path, status, body)
		}
	}
}

// Broken query strings must fail at the real typed-owner HTTP boundary,
// rather than trusting a browser-side pivot or a cached search result.
func TestCrossCasePivotOwnerPreviewRejectsUntrustedHTTPSelectors(t *testing.T) {
	h, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, method, target, payload string
		want int
	}{
		{"post-source", http.MethodPost, "/demo-source-open?case_id=case-a&entry_id=event-new", "", http.StatusMethodNotAllowed},
		{"missing-case", http.MethodGet, "/demo-source-open?entry_id=event-new", "", http.StatusBadRequest},
		{"duplicate-entry", http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=event-new&entry_id=event-old", "", http.StatusBadRequest},
		{"unknown-selector", http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=event-new&principal=analyst-b", "", http.StatusBadRequest},
		{"wildcard-case", http.MethodGet, "/demo-source-open?case_id=%2A&entry_id=event-new", "", http.StatusBadRequest},
		{"control-case", http.MethodGet, "/demo-source-open?case_id=case-a%0A&entry_id=event-new", "", http.StatusBadRequest},
		{"oversized-query", http.MethodGet, "/demo-source-open?case_id=case-a&entry_id="+strings.Repeat("x", 4096), "", http.StatusBadRequest},
		{"get-body", http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=event-new", "not-allowed", http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var payload io.Reader
			if tc.payload != "" {
				payload = strings.NewReader(tc.payload)
			}
			r := httptest.NewRequest(tc.method, tc.target, payload)
			r.Header.Set("Authorization", "Bearer "+testBearer)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.want || strings.Contains(w.Body.String(), "demo-event-a") ||
				strings.Contains(w.Body.String(), "demo-hidden-source") {
				t.Fatalf("untrusted selector exposed owner data: got %d want %d: %s", w.Code, tc.want, w.Body.String())
			}
		})
	}
	r := httptest.NewRequest(http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=event-new", nil)
	r.Header.Set("Authorization", "Bearer "+testBearer)
	ctx, cancel := context.WithCancel(r.Context())
	cancel()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	if w.Code != http.StatusServiceUnavailable || strings.Contains(w.Body.String(), "demo-event-a") {
		t.Fatalf("cancelled source lookup returned identity: %d %s", w.Code, w.Body.String())
	}
}
