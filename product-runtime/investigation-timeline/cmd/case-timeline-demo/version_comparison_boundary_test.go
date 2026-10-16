package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCaseVersionComparisonRejectsUnauthenticatedAndMalformedBoundary(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil { t.Fatal(err) }
	valid := "/case-version-comparison?case_id=case-a&entry_id=event-new"
	for _, tc := range []struct {
		name string
		method string
		path string
		body string
		bearer string
		want int
	}{
		{"unauthenticated", "GET", valid, "", "", http.StatusUnauthorized},
		{"wrong-bearer", "GET", valid, "", "forged", http.StatusUnauthorized},
		{"method", "POST", valid, "", testBearer, http.StatusMethodNotAllowed},
		{"body", "GET", valid, "unsafe", testBearer, http.StatusBadRequest},
		{"no-case", "GET", "/case-version-comparison?entry_id=event-new", "", testBearer, http.StatusBadRequest},
		{"no-entry", "GET", "/case-version-comparison?case_id=case-a", "", testBearer, http.StatusBadRequest},
		{"duplicate-entry", "GET", valid+"&entry_id=event-old", "", testBearer, http.StatusBadRequest},
		{"caller-tenant", "GET", valid+"&tenant=demo-tenant-b", "", testBearer, http.StatusBadRequest},
		{"caller-permission", "GET", valid+"&permission=allow", "", testBearer, http.StatusBadRequest},
		{"empty-entry", "GET", "/case-version-comparison?case_id=case-a&entry_id=", "", testBearer, http.StatusBadRequest},
		{"wildcard-entry", "GET", "/case-version-comparison?case_id=case-a&entry_id=%2A", "", testBearer, http.StatusBadRequest},
		{"foreign-case", "GET", "/case-version-comparison?case_id=case-b&entry_id=event-b", "", testBearer, http.StatusNotFound},
		{"not-superseding", "GET", "/case-version-comparison?case_id=case-a&entry_id=event-old", "", testBearer, http.StatusNotFound},
		{"missing-entry", "GET", "/case-version-comparison?case_id=case-a&entry_id=missing", "", testBearer, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			if tc.body != "" { body = strings.NewReader(tc.body) }
			req := httptest.NewRequest(tc.method, "http://localhost"+tc.path, body)
			if tc.bearer != "" { req.Header.Set("Authorization", "Bearer "+tc.bearer) }
			out := httptest.NewRecorder()
			handler.ServeHTTP(out, req)
			if out.Code != tc.want { t.Fatalf("unexpected status %d wanted %d: %s",out.Code,tc.want,out.Body.String()) }
			if strings.Contains(out.Body.String(), "demo-hidden-source") || strings.Contains(out.Body.String(), "event-b") {
				t.Fatal("unauthorized source metadata disclosed")
			}
			if out.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("comparison response may be cached")
			}
		})
	}
	req := httptest.NewRequest("GET", "http://localhost"+valid, nil)
	req.Header.Set("Authorization", "Bearer "+testBearer)
	req.TransferEncoding = []string{"chunked"}
	out := httptest.NewRecorder()
	handler.ServeHTTP(out, req)
	if out.Code != http.StatusBadRequest { t.Fatalf("transfer encoding accepted: %d",out.Code) }

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = httptest.NewRequest("GET", "http://localhost"+valid, nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+testBearer)
	out = httptest.NewRecorder()
	handler.ServeHTTP(out, req)
	if out.Code != http.StatusServiceUnavailable { t.Fatalf("cancelled comparison accepted: %d",out.Code) }
}
