package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAnalystSearchControlReusesAuthorizedQueryFlow(t *testing.T) {
	page, err := analystAssets.ReadFile("analyst.html")
	if err != nil {
		t.Fatal(err)
	}
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(page), `id="search"`) ||
		!strings.Contains(string(page), `maxlength="96"`) ||
		!strings.Contains(string(js), `url.searchParams.set("q", search.value.trim())`) ||
		!strings.Contains(string(js), "kind, search, fromInput") {
		t.Fatal("search filter is not wired into same authenticated cancellation-safe viewer request")
	}
	h, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet,
		"/case-timeline?case_id=case-a&limit=1&q="+url.QueryEscape("synthetic-evidence"), nil)
	req.Header.Set("Authorization", "Bearer "+testBearer)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"total_visible":1`) ||
		!strings.Contains(resp.Body.String(), `"id":"event-old"`) ||
		strings.Contains(resp.Body.String(), "demo-hidden-source") {
		t.Fatalf("demo full-stack source/provenance search incorrect: %d %s", resp.Code, resp.Body.String())
	}
}
