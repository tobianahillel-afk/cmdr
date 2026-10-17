package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// This test goes through CLI parsing, two actual operator files, JSONL import,
// real loopback listener, bearer authentication and typed-source ACL handlers.
// No intermediate owner JSON or fabricated httptest route is used.
func TestLiveJSONLImportAuthenticatedAnalystLoopback(t *testing.T) {
	ownerPath := snapshotFixture(t)
	eventsPath := writeJSONLEvents(t, importedJSONL)
	originalOwner, err := os.ReadFile(ownerPath)
	if err != nil {
		t.Fatal(err)
	}
	originalEvents, err := os.ReadFile(eventsPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
			"-principal=analyst-a", "-snapshot=" + ownerPath,
			"-import-jsonl=" + eventsPath,
			"-into-case=forensic-case-on-disk", "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("live import exited before listener: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("imported analyst viewer did not open loopback")
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
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "127.0.0.1" ||
		parsed.Query().Get("case_id") != "forensic-case-on-disk" || len(bearer) != 64 {
		t.Fatalf("invalid imported listener endpoint: %q", endpoint)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	get := func(path, credential string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "http://"+parsed.Host+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		if credential != "" {
			req.Header.Set("Authorization", "Bearer "+credential)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			t.Fatal(err)
		}
		return resp.StatusCode, string(data)
	}
	digest := sha256.Sum256([]byte("forensic-case-on-disk\x00sysmon-1"))
	importedID := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
	if status, body := get("/case-timeline?case_id=forensic-case-on-disk&limit=20", bearer); status != http.StatusOK ||
		!strings.Contains(body, importedID) || !strings.Contains(body, "operator-local-jsonl") ||
		strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "demo-tenant-b") {
		t.Fatalf("live imported timeline failed authorization or chronology: %d %s", status, body)
	}
	if status, body := get("/demo-source-open?case_id=forensic-case-on-disk&entry_id="+url.QueryEscape(importedID), bearer);
		status != http.StatusOK || !strings.Contains(body, "Offline Sysmon process event") {
		t.Fatalf("live typed source preview not imported: %d %s", status, body)
	}
	for _, path := range []string{
		"/demo-cases", "/case-timeline?case_id=forensic-case-on-disk",
		"/demo-source-open?case_id=forensic-case-on-disk&entry_id=" + url.QueryEscape(importedID),
	} {
		if status, body := get(path, ""); status != http.StatusUnauthorized ||
			strings.Contains(body, importedID) || strings.Contains(body, "Offline Sysmon") {
			t.Fatalf("unauthenticated imported data disclosed: %s %d %s", path, status, body)
		}
	}
	for _, path := range []string{
		"/case-timeline?case_id=case-b",
		"/demo-source-open?case_id=forensic-case-on-disk&entry_id=restricted",
		"/demo-source-open?case_id=case-b&entry_id=event-b",
	} {
		if status, body := get(path, bearer); status != http.StatusNotFound ||
			strings.Contains(body, importedID) || strings.Contains(body, "demo-hidden-source") {
			t.Fatalf("foreign/denied data accessible: %s %d %s", path, status, body)
		}
	}
	if status, body := get("/analyst", ""); status != http.StatusOK ||
		!strings.Contains(body, "analyst") || strings.Contains(body, "Offline Sysmon") {
		t.Fatalf("analyst UI was not served safely: %d %s", status, body)
	}
	if strings.Contains(log.String(), "Offline Sysmon") ||
		strings.Contains(log.String(), "operator-local-jsonl") ||
		strings.Contains(log.String(), "demo-hidden-source") {
		t.Fatal("listener startup output disclosed imported owner data")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("imported listener shutdown failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("imported listener did not stop on cancellation")
	}
	finalOwner, err := os.ReadFile(ownerPath)
	if err != nil || !bytes.Equal(originalOwner, finalOwner) {
		t.Fatal("live imported viewer modified owner file")
	}
	finalEvents, err := os.ReadFile(eventsPath)
	if err != nil || !bytes.Equal(originalEvents, finalEvents) {
		t.Fatal("live imported viewer modified original JSONL input")
	}
}

func TestLiveJSONLImportRejectsUnauthorizedAndAmbiguousModes(t *testing.T) {
	owner := snapshotFixture(t)
	events := writeJSONLEvents(t, importedJSONL)
	invalidEvents := writeJSONLEvents(t, "{invalid json}\n")
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"missing-input", []string{"-serve-import"}},
		{"missing-owner", []string{"-serve-import", "-into-case=forensic-case-on-disk", "-import-jsonl=" + events}},
		{"unauthorized-case", []string{"-principal=analyst-b", "-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + events, "-serve-import"}},
		{"foreign-case", []string{"-principal=analyst-a", "-snapshot=" + owner, "-into-case=case-b", "-import-jsonl=" + events, "-serve-import"}},
		{"invalid-events", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + invalidEvents, "-serve-import"}},
		{"same-input-files", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + owner, "-serve-import"}},
		{"relative-events", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=relative.jsonl", "-serve-import"}},
		{"triage-conflict", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + events, "-triage-import", "-serve-import"}},
		{"comparison-conflict", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + events, "-with-case=case-c", "-serve-import"}},
		{"report-conflict", []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + events, "-report-case=forensic-case-on-disk", "-serve-import"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := runDemoCLI(context.Background(), tc.args, &output); err == nil || output.Len() != 0 {
				t.Fatalf("invalid live import started or disclosed bearer: %v output=%q", err, output.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := runDemoCLI(ctx, []string{"-snapshot=" + owner, "-into-case=forensic-case-on-disk",
		"-import-jsonl=" + events, "-serve-import"}, &output); err != context.Canceled || output.Len() != 0 {
		t.Fatalf("cancelled import produced token/listener: %v %q", err, output.String())
	}
}
