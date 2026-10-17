package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// An operator-owned JSONL file must reach the same authenticated HTTP
// timeline/source/analyst path without serializing another owner snapshot.
func TestLiveImportedJSONLBrowserJourney(t *testing.T) {
	owner := snapshotFixture(t)
	before, err := os.ReadFile(owner)
	if err != nil {
		t.Fatal(err)
	}
	events := writeJSONLEvents(t, importedJSONL)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
		"-principal=analyst-a", "-snapshot=" + owner, "-import-jsonl=" + events,
		"-into-case=forensic-case-on-disk", "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("live JSONL analyst listener did not start: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("live JSONL analyst listener was not ready")
	}
	var link, token string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			link = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			token = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	parsed, err := url.Parse(link)
	if err != nil || parsed.Hostname() != "127.0.0.1" || len(token) != 64 ||
		parsed.Query().Get("case_id") != "forensic-case-on-disk" {
		t.Fatalf("unsafe live JSONL listener or token: %q", link)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	get := func(path, bearer string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "http://"+parsed.Host+path, nil)
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
	caseID := "forensic-case-on-disk"
	baseQuery := url.Values{"case_id": {caseID}}
	sum := sha256.Sum256([]byte(caseID + "\x00sysmon-1"))
	entryID := "local-jsonl-entry-" + hex.EncodeToString(sum[:])
	for _, tc := range []struct {
		path string
		bearer string
		want int
		contains string
	}{
		{"/analyst", "", http.StatusOK, "CMDR"},
		{"/demo-cases", token, http.StatusOK, caseID},
		{"/case-timeline?" + baseQuery.Encode() + "&limit=200", token, http.StatusOK, entryID},
		{"/demo-source-open?" + url.Values{"case_id": {caseID}, "entry_id": {entryID}}.Encode(),
			token, http.StatusOK, "Offline Sysmon process event"},
		{"/case-source-navigation?" + url.Values{"case_id": {caseID}, "entry_id": {entryID}}.Encode(),
			token, http.StatusOK, entryID},
		{"/case-timeline?" + baseQuery.Encode(), "", http.StatusUnauthorized, ""},
		{"/demo-source-open?" + url.Values{"case_id": {caseID}, "entry_id": {entryID}}.Encode(),
			"", http.StatusUnauthorized, ""},
		{"/case-timeline?case_id=case-b", token, http.StatusNotFound, ""},
	} {
		status, body := get(tc.path, tc.bearer)
		if status != tc.want || (tc.contains != "" && !strings.Contains(body, tc.contains)) ||
			strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "DO NOT DISPLAY") ||
			strings.Contains(body, "demo-tenant-b") {
			t.Fatalf("live import route %s status=%d want=%d body=%s", tc.path, status, tc.want, body)
		}
	}
	current, err := os.ReadFile(owner)
	if err != nil || !bytes.Equal(current, before) {
		t.Fatal("live JSONL import changed the original owner file")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("live JSONL analyst listener did not stop")
	}
}

func TestLiveImportedJSONLRejectsUntrustedInputsBeforeListener(t *testing.T) {
	owner := snapshotFixture(t)
	events := writeJSONLEvents(t, importedJSONL)
	bad := writeJSONLEvents(t, "{\"id\":\"x\",\"timestamp\":\"wrong\",\"message\":\"bad\"}\n")
	cases := []struct {
		name string
		args []string
	}{
		{"missing-owner", []string{"-serve-import", "-import-jsonl=" + events, "-into-case=case-a"}},
		{"missing-events", []string{"-serve-import", "-snapshot=" + owner, "-into-case=forensic-case-on-disk"}},
		{"missing-case", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + events}},
		{"same-file", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + owner, "-into-case=forensic-case-on-disk"}},
		{"invalid-events", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + bad, "-into-case=forensic-case-on-disk"}},
		{"foreign-principal", []string{"-serve-import", "-principal=analyst-b", "-snapshot=" + owner, "-import-jsonl=" + events, "-into-case=forensic-case-on-disk"}},
		{"foreign-case", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + events, "-into-case=case-b"}},
		{"mixed-triage", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + events, "-into-case=forensic-case-on-disk", "-triage-import"}},
		{"mixed-export", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=" + events, "-into-case=forensic-case-on-disk", "-export-case=forensic-case-on-disk"}},
		{"relative-events", []string{"-serve-import", "-snapshot=" + owner, "-import-jsonl=relative.jsonl", "-into-case=forensic-case-on-disk"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := runDemoCLI(context.Background(), tc.args, &output); err == nil || output.Len() != 0 {
				t.Fatalf("unsafe live import accepted or listener started: %v %q", err, output.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := serveDemoWithImportedJSONL(ctx, "analyst-a", owner, "forensic-case-on-disk", events, &out); err != context.Canceled || out.Len() != 0 {
		t.Fatalf("cancelled live import published listener: %v %q", err, out.String())
	}
}

func TestLiveImportStartupSelectsImportedCaseAmongReadableCases(t *testing.T) {
	owner := snapshotFixture(t)
	data, err := loadDemoSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	// The operator may have an alphabetically earlier, readable Case.
	// The requested imported Case, not that unrelated Case, should open.
	other := data.Cases[0]
	other.ID = "aaa-another-readable-case"
	other.EntryIDs = nil
	data.Cases = append(data.Cases, other)
	if first, err := firstDemoAuthorizedCase(data, "analyst-a"); err != nil || first != other.ID {
		t.Fatalf("fixture has no competing authorized Case: %q %v", first, err)
	}
	content, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "multi-case-owner.json")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	events := writeJSONLEvents(t, importedJSONL)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
			"-principal=analyst-a", "-snapshot=" + path, "-import-jsonl=" + events,
			"-into-case=forensic-case-on-disk", "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("live import of nonfirst Case failed: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("live multi-Case listener did not start")
	}
	var link string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			link = strings.TrimPrefix(line, "URL: ")
		}
	}
	parsed, err := url.Parse(link)
	if err != nil || parsed.Query().Get("case_id") != "forensic-case-on-disk" {
		t.Fatalf("imported Case is not initial analyst selection: %q", link)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("multi-Case import listener did not stop")
	}
}

func TestLiveImportRejectsInvalidOwnerBeforeBearerCreation(t *testing.T) {
	data := demoData()
	data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, "nonexistent-entry")
	var output bytes.Buffer
	if err := serveDemoWithData(context.Background(), "analyst-a", data, "", &output); err == nil || output.Len() != 0 {
		t.Fatalf("untrusted owner state started a listener: %v %q", err, output.String())
	}
	if err := serveDemoWithData(context.Background(), "", demoData(), "", &output); err == nil || output.Len() != 0 {
		t.Fatal("empty analyst identity accepted by local server")
	}
}
