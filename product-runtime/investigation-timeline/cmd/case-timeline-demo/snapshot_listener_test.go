package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestOperatorSnapshotRealLoopbackAnalystJourney(t *testing.T) {
	path := snapshotFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{"-principal=analyst-a", "-snapshot=" + path}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("operator Case snapshot listener did not start: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("operator snapshot listener not ready")
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
	if err != nil || parsed.Hostname() != "127.0.0.1" ||
		parsed.Query().Get("case_id") != "forensic-case-on-disk" || len(token) != 64 {
		t.Fatalf("operator snapshot advertised an incorrect or unsafe route: %q", endpoint)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	for _, tc := range []struct {
		path string
		bearer string
		want int
		contains string
	}{
		{endpoint, token, 200, "event-new"},
		{"http://" + parsed.Host + "/demo-cases", token, 200, "forensic-case-on-disk"},
		{"http://" + parsed.Host + "/case-timeline?case_id=case-b", token, 404, ""},
		{endpoint, "", 401, ""},
	} {
		req, err := http.NewRequest(http.MethodGet, tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if tc.bearer != "" {
			req.Header.Set("Authorization", "Bearer " + tc.bearer)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()
		if err != nil || resp.StatusCode != tc.want ||
			(tc.contains != "" && !strings.Contains(string(body), tc.contains)) ||
			strings.Contains(string(body), "demo-hidden-source") ||
			strings.Contains(string(body), "demo-tenant-b") {
			t.Fatalf("live snapshot HTTP %d want %d: %s: %v", resp.StatusCode, tc.want, body, err)
		}
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("snapshot-backed listener did not stop")
	}
}
