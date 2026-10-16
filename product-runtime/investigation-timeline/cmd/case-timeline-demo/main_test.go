package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"strings"
	"testing"
	"time"

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

type demoCapture struct {
	mu sync.Mutex
	output bytes.Buffer
	ready chan struct{}
	once sync.Once
}

func (d *demoCapture) Write(data []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	count, err := d.output.Write(data)
	if strings.Contains(d.output.String(), "Ctrl-C exits.") {
		d.once.Do(func() { close(d.ready) })
	}
	return count, err
}

func (d *demoCapture) String() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.output.String()
}

func TestDemoRealLoopbackListenerAndShutdown(t *testing.T) {
	for _, principal := range []string{"analyst-a", "analyst-b"} {
		t.Run(principal, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			log := &demoCapture{ready: make(chan struct{})}
			done := make(chan error, 1)
			go func() { done <- runDemoCLI(ctx, []string{"-principal="+principal}, log) }()
			select {
			case <-log.ready:
			case err := <-done:
				t.Fatalf("demo exited before serving: %v", err)
			case <-time.After(10 * time.Second):
				t.Fatal("demo did not announce its endpoint")
			}
			var endpoint, token string
			for _, line := range strings.Split(log.String(), "\n") {
				if strings.HasPrefix(line, "URL: ") {
					endpoint = strings.TrimPrefix(line, "URL: ")
				}
				if strings.HasPrefix(line, "Bearer token (this run only): ") {
					token = strings.TrimPrefix(line, "Bearer token (this run only): ")
				}
			}
			parsed, err := url.Parse(endpoint)
			if err != nil || parsed.Hostname() != "127.0.0.1" || len(token) != 64 {
				t.Fatalf("unsafe or invalid loopback endpoint: %q credential len=%d: %v", endpoint, len(token), err)
			}
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			client := &http.Client{Timeout: 2 * time.Second}
			response, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
			response.Body.Close()
			if err != nil || response.StatusCode != 200 {
				t.Fatalf("real HTTP route failed: %d %s %v", response.StatusCode, body, err)
			}
			wantCase := "case-a"
			if principal == "analyst-b" {
				wantCase = "case-b"
			}
			if parsed.Query().Get("case_id") != wantCase || !strings.Contains(string(body), "entries") {
				t.Fatalf("wrong case routing %q: %s", endpoint, body)
			}
			anonymous, err := http.Get(endpoint)
			if err != nil {
				t.Fatal(err)
			}
			anonymous.Body.Close()
			if anonymous.StatusCode != http.StatusUnauthorized {
				t.Fatalf("anonymous loopback access permitted: %d", anonymous.StatusCode)
			}
			cancel()
			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("server shutdown failed: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("loopback server did not stop after cancellation")
			}
		})
	}
}

func TestDemoRefusesCancelledStartupAndInvalidPrincipal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := serveDemo(ctx, "analyst-a", &out); err != context.Canceled || out.Len() != 0 {
		t.Fatalf("cancelled startup exposed a token or listener: %v %s", err, out.String())
	}
	if err := serveDemo(context.Background(), "unrecognized", &out); err == nil || out.Len() != 0 {
		t.Fatalf("unknown principal was accepted: %v %s", err, out.String())
	}
}

func TestDemoCLIRejectsUnsupportedOptions(t *testing.T) {
	for _, args := range [][]string{
		{"-unknown"},
		{"-principal=analyst-a", "unexpected"},
		{"-principal=admin"},
	} {
		var output bytes.Buffer
		ctx, cancel := context.WithCancel(context.Background())
		if err := runDemoCLI(ctx, args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("bad CLI parameters accepted: %v, output=%q", args, output.String())
		}
		cancel()
	}
	var output bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := runDemoCLI(ctx, nil, &output); err != context.Canceled || output.Len() != 0 {
		t.Fatalf("cancelled default CLI start accepted: %v", err)
	}
}

func TestDemoAnalystTimeWindowAndKindEndToEnd(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	cases := []struct {
		name string
		parameters url.Values
		wantIDs []string
		wantCode int
	}{
		{"before-upper-exclusive", url.Values{"from":{"2026-10-01T09:30:00Z"}, "to":{"2026-10-01T11:00:00Z"}}, []string{"event-old"}, 200},
		{"after-lower-inclusive", url.Values{"from":{"2026-10-01T11:00:00Z"}, "to":{"2026-10-01T12:00:00Z"}}, []string{"event-new"}, 200},
		{"normalized-zone", url.Values{"from":{"2026-10-01T12:00:00+01:00"}, "to":{"2026-10-01T13:00:00+01:00"}, "kind":{"observed"}}, []string{"event-new"}, 200},
		{"kind-plus-time", url.Values{"from":{"2026-10-01T09:30:00Z"}, "to":{"2026-10-01T11:30:00Z"}, "kind":{"inferred"}}, []string{"event-old"}, 200},
		{"empty-valid-window", url.Values{"from":{"2030-01-01T00:00:00Z"}}, nil, 200},
		{"reversed-window", url.Values{"from":{"2026-10-01T12:00:00Z"}, "to":{"2026-10-01T10:00:00Z"}}, nil, 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parameters := tc.parameters
			parameters.Set("case_id", "case-a")
			parameters.Set("limit", "20")
			status, body := requestDemo(t, server, "/case-timeline?"+parameters.Encode(), testBearer)
			if status != tc.wantCode {
				t.Fatalf("time-window request HTTP %d want %d: %s", status, tc.wantCode, body)
			}
			if tc.wantCode != 200 {
				if strings.Contains(body, "event-old") || strings.Contains(body, "event-new") {
					t.Fatal("invalid time window disclosed timeline metadata")
				}
				return
			}
			var response struct{
				Entries []timeline.ProjectedEntry `json:"entries"`
				Page struct{ TotalVisible int `json:"total_visible"` } `json:"page"`
			}
			if err := json.Unmarshal([]byte(body), &response); err != nil {
				t.Fatal(err)
			}
			if len(response.Entries) != len(tc.wantIDs) || response.Page.TotalVisible != len(tc.wantIDs) {
				t.Fatalf("wrong authorized page: %+v", response)
			}
			for i, id := range tc.wantIDs {
				if response.Entries[i].ID != id {
					t.Fatalf("incorrect chronology after time filtering: %+v", response.Entries)
				}
			}
			if strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "demo-tenant-b") {
				t.Fatal("time filter disclosed hidden or foreign source")
			}
		})
	}
	for _, tc := range []struct{caseID, bearer string; want int}{
		{"case-a","",401}, {"case-b",testBearer,404},
	}{
		p := url.Values{"case_id":{tc.caseID},"from":{"2026-10-01T09:30:00Z"},"to":{"2026-10-01T12:00:00Z"}}
		status, body := requestDemo(t, server, "/case-timeline?"+p.Encode(), tc.bearer)
		if status != tc.want || strings.Contains(body, "event-new") || strings.Contains(body, "demo-incident-b") {
			t.Fatalf("time selection bypassed authorization: status %d body %s", status, body)
		}
	}
}

func TestDemoCaseSourceNavigationEndToEnd(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil { t.Fatal(err) }
	server := httptest.NewServer(handler)
	defer server.Close()
	status, body := requestDemo(t, server,
		"/case-source-navigation?case_id=case-a&entry_id=event-new&kind=observed&limit=20&offset=0", testBearer)
	if status != http.StatusOK {
		t.Fatalf("navigation route did not resolve: %d %s", status, body)
	}
	var result struct {
		Status string `json:"status"`
		Kind string `json:"kind"`
		ObjectID string `json:"object_id"`
		ReturnCaseID string `json:"return_case_id"`
		ReturnEntryID string `json:"return_entry_id"`
		ReturnQuery string `json:"return_query"`
		ReturnIsAuthority bool `json:"return_is_authority"`
	}
	if err := json.Unmarshal([]byte(body), &result); err != nil { t.Fatal(err) }
	if result.Status != "ready" || result.Kind != "telemetry-event" ||
		result.ObjectID != "demo-event-a" || result.ReturnCaseID != "case-a" ||
		result.ReturnEntryID != "event-new" || result.ReturnIsAuthority ||
		result.ReturnQuery != "case_id=case-a&kind=observed&limit=20&offset=0" {
		t.Fatalf("Case -> Timeline -> owner destination failed: %+v", result)
	}
	for _, tc := range []struct{path, token string; want int}{
		{"/case-source-navigation?case_id=case-a&entry_id=event-new","",401},
		{"/case-source-navigation?case_id=case-a&entry_id=restricted",testBearer,404},
		{"/case-source-navigation?case_id=case-b&entry_id=event-b",testBearer,404},
		{"/case-source-navigation?case_id=case-a&entry_id=event-b",testBearer,404},
		{"/case-source-navigation?case_id=case-a&entry_id=event-new&object_id=demo-event-a",testBearer,400},
		{"/case-source-navigation?case_id=case-a&entry_id=event-new&return_url=https://evil.invalid",testBearer,400},
	} {
		status, body := requestDemo(t, server, tc.path, tc.token)
		if status != tc.want || strings.Contains(body, "demo-secret") ||
			strings.Contains(body, "demo-event-a") {
			t.Fatalf("denial leaked owner details: %d %s", status, body)
		}
	}
}
