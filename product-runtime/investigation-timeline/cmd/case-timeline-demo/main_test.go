package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

const testBearer = "1111111111111111111111111111111111111111111111111111111111111111"

func requestDemo(t *testing.T, server *httptest.Server, path, bearer string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	response, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, string(body)
}

func TestDemoActualHTTPCaseInvestigationFlow(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	status, body := requestDemo(t, server, "/case-timeline?case_id=case-a&limit=1&offset=0", testBearer)
	if status != http.StatusOK {
		t.Fatalf("authorized page HTTP %d: %s", status, body)
	}
	var page struct {
		Status  string                    `json:"status"`
		Entries []timeline.ProjectedEntry `json:"entries"`
		Page    struct {
			TotalVisible int  `json:"total_visible"`
			HasMore      bool `json:"has_more"`
		} `json:"page"`
	}
	if err := json.Unmarshal([]byte(body), &page); err != nil {
		t.Fatal(err)
	}
	if page.Status != "filtered" || page.Page.TotalVisible != 2 || !page.Page.HasMore ||
		len(page.Entries) != 1 || page.Entries[0].ID != "event-old" ||
		page.Entries[0].Kind != "inferred" {
		t.Fatalf("authorized page wrong: %+v", page)
	}
	status, body = requestDemo(t, server, "/case-timeline?case_id=case-a&limit=1&offset=1", testBearer)
	if status != http.StatusOK || strings.Contains(body, "demo-hidden-source") {
		t.Fatalf("next page HTTP %d: %s", status, body)
	}
	if err := json.Unmarshal([]byte(body), &page); err != nil {
		t.Fatal(err)
	}
	if page.Page.HasMore || page.Page.TotalVisible != 2 || len(page.Entries) != 1 ||
		page.Entries[0].ID != "event-new" || page.Entries[0].Kind != "observed" ||
		page.Entries[0].Supersedes != "event-old" {
		t.Fatalf("chronology and supersession lost: %+v", page)
	}
	for _, secret := range []string{"restricted", "demo-hidden-source", "demo-secret", "demo-tenant-b"} {
		if strings.Contains(body, secret) {
			t.Fatalf("denied or foreign metadata leaked: %s", secret)
		}
	}
	status, body = requestDemo(t, server,
		"/case-timeline?case_id=case-a&kind=observed&from="+url.QueryEscape("2026-10-01T10:30:00Z"),
		testBearer)
	if status != http.StatusOK || !strings.Contains(body, `"id":"event-new"`) ||
		strings.Contains(body, `"id":"event-old"`) {
		t.Fatalf("investigation filter broken: %d %s", status, body)
	}
}

func TestDemoPrincipalIsolationAndDeniedRequests(t *testing.T) {
	for _, tc := range []struct {
		name, principal, uri, token string
		expected int
	}{
		{"no-token", "analyst-a", "/case-timeline?case_id=case-a", "", 401},
		{"bad-token", "analyst-a", "/case-timeline?case_id=case-a", "not-the-token", 401},
		{"cross-tenant-a", "analyst-a", "/case-timeline?case_id=case-b", testBearer, 404},
		{"cross-tenant-b", "analyst-b", "/case-timeline?case_id=case-a", testBearer, 404},
		{"missing-case", "analyst-b", "/case-timeline?case_id=absent", testBearer, 404},
		{"injected-scope", "analyst-a", "/case-timeline?case_id=case-a&tenant=demo-tenant-b", testBearer, 400},
		{"valid-b", "analyst-b", "/case-timeline?case_id=case-b", testBearer, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := newDemoHandler(testBearer, tc.principal)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(handler)
			defer server.Close()
			status, body := requestDemo(t, server, tc.uri, tc.token)
			if status != tc.expected {
				t.Fatalf("HTTP %d wanted %d: %s", status, tc.expected, body)
			}
			if tc.expected != 200 && (strings.Contains(body, "event-new") ||
				strings.Contains(body, "demo-tenant-a") || strings.Contains(body, "demo-hidden-source")) {
				t.Fatalf("error disclosed protected metadata: %s", body)
			}
		})
	}
}

func TestDemoConstructionAndBearerCredential(t *testing.T) {
	for _, tc := range []struct{ token, principal string }{
		{"", "analyst-a"}, {"short", "analyst-a"},
		{strings.Repeat("z", 64), "analyst-a"},
		{testBearer, "admin"}, {testBearer, ""},
	} {
		if _, err := newDemoHandler(tc.token, tc.principal); err == nil {
			t.Fatalf("accepted invalid fixture authority %+v", tc)
		}
	}
	a, err := randomBearer()
	if err != nil {
		t.Fatal(err)
	}
	b, err := randomBearer()
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 64 || a == b {
		t.Fatal("demo ephemeral token generation failed")
	}
	if _, err := hex.DecodeString(a); err != nil {
		t.Fatal(err)
	}
	if _, err := newDemoHandler(a, "analyst-b"); err != nil {
		t.Fatal(err)
	}
}

func TestDemoRejectsNonReadOperations(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/case-timeline?case_id=case-a", strings.NewReader("payload"))
		req.Header.Set("Authorization", "Bearer "+testBearer)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s unexpectedly accepted: %d", method, response.Code)
		}
	}
}
