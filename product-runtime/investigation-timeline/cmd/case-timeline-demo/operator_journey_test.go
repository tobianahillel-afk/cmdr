package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Exercise the actual loopback CLI, not just a synthetic handler or fixtures.
// The operator-selected JSON comes from the earlier reusable offline-owner PR.
func TestOperatorSnapshotActualLoopbackCaseToSourceJourney(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("testdata", "example-owners.json"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := loadDemoSnapshot(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := newDemoHandlerWithData(testBearer, "investigator-42", data); err != nil {
		t.Fatalf("operator principal rejected: %v", err)
	}
	otherTenant, err := newDemoHandlerWithData(testBearer, "investigator-99", data)
	if err != nil {
		t.Fatalf("second tenant principal rejected: %v", err)
	}
	isolated := httptest.NewServer(otherTenant)
	defer isolated.Close()
	statusB, bodyB := requestDemo(t, isolated, "/demo-cases", testBearer)
	if statusB != http.StatusOK || !strings.Contains(bodyB, "case-private-99") ||
		strings.Contains(bodyB, "case-incident-42") {
		t.Fatalf("second tenant discovery not isolated: %d %s", statusB, bodyB)
	}
	statusB, bodyB = requestDemo(t, isolated, "/case-timeline?case_id=case-incident-42", testBearer)
	if statusB != http.StatusNotFound || strings.Contains(bodyB, "alert-1") {
		t.Fatalf("second tenant reached first tenant Case: %d %s", statusB, bodyB)
	}
	if _, err := newDemoHandlerWithData(testBearer, "analyst-a", data); err == nil {
		t.Fatal("embedded principal was incorrectly granted access to imported owner data")
	}
	if _, err := newDemoHandler(testBearer, "investigator-42"); err == nil {
		t.Fatal("embedded synthetic dataset accepted an imported principal")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	capture := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{"-principal=investigator-42", "-snapshot=" + path}, capture)
	}()
	select {
	case <-capture.ready:
	case err := <-done:
		t.Fatalf("local Case server exited before ready: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("local Case server did not start")
	}
	var endpoint, bearer string
	for _, line := range strings.Split(capture.String(), "\n") {
		switch {
		case strings.HasPrefix(line, "URL: "):
			endpoint = strings.TrimPrefix(line, "URL: ")
		case strings.HasPrefix(line, "Bearer token (this run only): "):
			bearer = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Hostname() != "127.0.0.1" || len(bearer) != 64 ||
		u.Query().Get("case_id") != "case-incident-42" {
		t.Fatalf("invalid authenticated operator endpoint %q: %v", endpoint, err)
	}
	client := &http.Client{Timeout: 2 * time.Second}
	get := func(path, token string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, u.Scheme+"://"+u.Host+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
		if err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, string(body)
	}
	for _, tc := range []struct {
		path, bearer string
		code int
		must, forbidden string
	}{
		{"/demo-cases", bearer, 200, "case-incident-43", "case-private-99"},
		{"/case-timeline?case_id=case-incident-43&limit=20", bearer, 200, "alert-3", "alert-99"},
		{"/case-source-navigation?case_id=case-incident-43&entry_id=alert-3&limit=20&offset=0", bearer, 200, "event-3", "incident-99"},
		{"/demo-source-open?case_id=case-incident-43&entry_id=alert-3", bearer, 200, "Synthetic malware triage observation", "Synthetic tenant-B incident"},
		{"/case-timeline?case_id=case-incident-42&limit=20", bearer, 200, "alert-2", "alert-99"},
		{"/case-source-navigation?case_id=case-incident-42&entry_id=alert-1&limit=20&offset=0", bearer, 200, "event-1", "incident-99"},
		{"/demo-source-open?case_id=case-incident-42&entry_id=alert-1", bearer, 200, "Synthetic DNS observation", "Synthetic tenant-B incident"},
		{"/case-timeline?case_id=case-private-99", bearer, 404, "", "alert-99"},
		{"/case-source-navigation?case_id=case-private-99&entry_id=alert-99", bearer, 404, "", "incident-99"},
		{"/demo-cases", "", 401, "", "case-incident-42"},
		{"/case-timeline?case_id=case-incident-42", "", 401, "", "alert-1"},
	} {
		t.Run(tc.path+"/"+tc.bearer[:min(4,len(tc.bearer))], func(t *testing.T) {
			status, body := get(tc.path, tc.bearer)
			if status != tc.code ||
				tc.must != "" && !strings.Contains(body, tc.must) ||
				tc.forbidden != "" && strings.Contains(body, tc.forbidden) {
				t.Fatalf("unexpected operator-backed HTTP response: %d %s", status, body)
			}
		})
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("local Case server shutdown failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("local Case server did not stop")
	}
}
