package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func TestOperatorSnapshotCustomAnalystCompletesFullSourceJourney(t *testing.T) {
	// An operator-selected disk owner can declare a different synthetic
	// analyst; Case, Timeline Entry and source ACLs remain server-owned.
	data := demoData()
	const principal = "forensic-operator-17"
	const caseID = "case & forensics #17"
	data.Principals[0].ID = principal
	data.Cases[0].ID = caseID
	data.Cases[0].Readers = []string{principal}
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-a" {
			data.Entries[i].CaseID = caseID
			data.Entries[i].Readers = []string{principal}
		}
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.CaseID == "case-a" {
			data.Sources[i].Identity.CaseID = caseID
			if len(data.Sources[i].Readers) > 0 {
				data.Sources[i].Readers = []string{principal}
			}
		}
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "custom-owner.json")
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{"-snapshot=" + path, "-principal=" + principal}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("custom analyst listener failed: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("custom analyst listener not ready")
	}
	var address, token string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			address = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			token = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	parsed, err := url.Parse(address)
	if err != nil || parsed.Hostname() != "127.0.0.1" ||
		parsed.Query().Get("case_id") != caseID || parsed.Query().Get("limit") != "2" || len(token) != 64 {
		t.Fatalf("invalid encoded Case startup URL: %q %v", address, err)
	}
	root := "http://" + parsed.Host
	client := &http.Client{Timeout: 3 * time.Second}
	request := func(path, bearer string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, root+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			t.Fatal(err)
		}
		return resp.StatusCode, string(body)
	}
	for _, check := range []struct {
		path     string
		status   int
		required string
	}{
		{"/demo-cases", http.StatusOK, ""},
		{"/case-timeline?" + url.Values{"case_id": {caseID}}.Encode(), http.StatusOK, "event-new"},
		{"/case-source-navigation?" + url.Values{"case_id": {caseID}, "entry_id": {"event-new"}}.Encode(), http.StatusOK, "demo-event-a"},
		{"/demo-source-open?" + url.Values{"case_id": {caseID}, "entry_id": {"event-new"}}.Encode(), http.StatusOK, "Événement réseau synthétique"},
		{"/case-timeline?case_id=case-b", http.StatusNotFound, ""},
		{"/demo-source-open?" + url.Values{"case_id": {caseID}, "entry_id": {"restricted"}}.Encode(), http.StatusNotFound, ""},
	} {
		status, body := request(check.path, token)
		if status != check.status || (check.required != "" && !strings.Contains(body, check.required)) ||
			strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "demo-secret") ||
			strings.Contains(body, "demo-tenant-b") {
			t.Fatalf("disk-backed owner journey %s: HTTP %d body %q", check.path, status, body)
		}
		if check.path == "/demo-cases" {
			var visible struct {
				Cases []string `json:"cases"`
			}
			if err := json.Unmarshal([]byte(body), &visible); err != nil ||
				len(visible.Cases) != 1 || visible.Cases[0] != caseID {
				t.Fatalf("incorrect JSON-escaped Case discovery: %+v %v", visible, err)
			}
		}
	}
	if status, body := request("/demo-source-open?"+
		url.Values{"case_id": {caseID}, "entry_id": {"event-new"}}.Encode(), ""); status != http.StatusUnauthorized ||
		strings.Contains(body, "demo-event-a") {
		t.Fatalf("anonymous source read succeeded: %d %s", status, body)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("custom analyst listener did not stop")
	}
}
