package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// An imported analyst may have multiple authorized Cases. The local viewer
// must derive the initial Case from the owner snapshot, not a demo constant.
func TestImportedMultipleCaseDiscoveryAndInitialSelection(t *testing.T) {
	data := importedCaseDataset()
	data.Cases = append(data.Cases, timeline.MemoryTimelineCase{
		ID: "case-000", Tenant: "demo-tenant-a", Environment: "lab",
		Readers: []string{"investigator-42"},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- serveDemoWithData(ctx, "investigator-42", log, data)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("multi-Case server failed before startup: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("multi-Case server startup timed out")
	}
	var endpoint, bearer string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			endpoint = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			bearer = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Hostname() != "127.0.0.1" ||
		u.Query().Get("case_id") != "case-000" || len(bearer) != 64 {
		t.Fatalf("incorrect owner-derived initial Case URL %q", endpoint)
	}
	req, err := http.NewRequest(http.MethodGet, "http://"+u.Host+"/demo-cases", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	response.Body.Close()
	if err != nil || response.StatusCode != 200 || !strings.Contains(string(body), "case-000") ||
		!strings.Contains(string(body), "case-imported") ||
		strings.Contains(string(body), "case-b") {
		t.Fatalf("multi-Case picker disclosed foreign Case or missed permitted Case: %d %s, %v", response.StatusCode, body, err)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("multi-Case server ignored cancellation")
	}
}
