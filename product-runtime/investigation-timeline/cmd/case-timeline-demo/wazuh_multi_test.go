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

func TestMultipleWazuhFilesThroughAuthorizedLoopback(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"123.456\",\"timestamp\":\"2026-10-17T10:01:00+0200\",\"rule\":{\"description\":\"Wazuh first feed alert\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"123.456\",\"timestamp\":\"2026-10-17T10:01:01+0200\",\"rule\":{\"description\":\"Wazuh second feed alert\"}}\n")
	before := make(map[string][]byte)
	for _, path := range []string{owner, first, second} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[path] = data
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
			"-principal=analyst-a", "-snapshot=" + owner,
			"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
			"-into-case=forensic-case-on-disk", "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("multiple Wazuh import exited early: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("multiple Wazuh import never started loopback listener")
	}
	var address, bearer string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			address = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			bearer = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	endpoint, err := url.Parse(address)
	if err != nil || endpoint.Hostname() != "127.0.0.1" || len(bearer) != 64 {
		t.Fatalf("expected authenticated loopback only, got %q", address)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	request := func(route, credential string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "http://"+endpoint.Host+route, nil)
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
	path := "/case-timeline?case_id=forensic-case-on-disk&limit=20"
	status, page := request(path, bearer)
	if status != http.StatusOK {
		t.Fatalf("authorized Case Timeline unavailable: %d %s", status, page)
	}
	for feed, preview := range []string{"Wazuh first feed alert", "Wazuh second feed alert"} {
		identity := "forensic-case-on-disk\x00feed-" + string(rune('1'+feed)) + "\x00wazuh-123.456"
		digest := sha256.Sum256([]byte(identity))
		entry := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
		if !strings.Contains(page, entry) || !strings.Contains(page, "operator-local-wazuh4/feed-") {
			t.Fatalf("feed %d entry absent from authorized timeline: %s", feed+1, page)
		}
		uri := "/demo-source-open?case_id=forensic-case-on-disk&entry_id=" + url.QueryEscape(entry)
		code, source := request(uri, bearer)
		if code != http.StatusOK || !strings.Contains(source, preview) {
			t.Fatalf("feed %d source preview absent: %d %s", feed+1, code, source)
		}
		code, body := request(uri, "")
		if code != http.StatusUnauthorized || strings.Contains(body, preview) {
			t.Fatalf("anonymous feed %d access leaked: %d %s", feed+1, code, body)
		}
	}
	for _, path := range []string{"/case-timeline?case_id=case-b", "/demo-source-open?case_id=case-b&entry_id=event-b"} {
		code, body := request(path, bearer)
		if code != http.StatusNotFound || strings.Contains(body, "Wazuh") {
			t.Fatalf("cross-tenant Wazuh data leaked: %d %s", code, body)
		}
	}
	if strings.Contains(log.String(), "Wazuh first feed alert") || strings.Contains(log.String(), "Wazuh second feed alert") {
		t.Fatal("Wazuh source content disclosed in startup diagnostics")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Wazuh loopback failed to shut down")
	}
	for path, prior := range before {
		current, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(current, prior) {
			t.Fatalf("operator input modified: %v", err)
		}
	}
}

func TestMultipleWazuhSourcesFailClosedBeforeListener(t *testing.T) {
	owner := snapshotFixture(t)
	good := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"rule\":{\"description\":\"good\"}}\n")
	bad := writeJSONLEvents(t, "{invalid}\n")
	other := writeJSONLEvents(t, "{\"id\":\"b.2\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"rule\":{\"description\":\"good other\"}}\n")
	cases := []struct {
		name, principal string
		paths []string
	}{
		{"duplicate-input", "analyst-a", []string{good, good}},
		{"same-as-owner", "analyst-a", []string{good, owner}},
		{"malformed-second-feed", "analyst-a", []string{good, bad}},
		{"unauthorized-principal", "analyst-b", []string{good, other}},
		{"too-many-files", "analyst-a", []string{good, other, good, other, bad}},
		{"missing-input", "analyst-a", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			err := serveDemoWithWazuhAlertsFiles(context.Background(), tc.principal, owner,
				"forensic-case-on-disk", tc.paths, &output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("bad feed set opened listener or disclosed bearer: %v / %q", err, output.String())
			}
		})
	}
	for _, args := range [][]string{
		{"-snapshot="+owner, "-into-case=forensic-case-on-disk", "-wazuh-alerts="+good, "-wazuh-alerts="+other,
			"-triage-import", "-with-case=case-c", "-compare-proximity"},
		{"-snapshot="+owner, "-into-case=forensic-case-on-disk", "-wazuh-alerts="+good, "-wazuh-alerts="+other,
			"-serve-import", "-import-jsonl="+other},
		{"-snapshot="+owner, "-into-case=forensic-case-on-disk", "-serve-import",
			"-wazuh-alerts="+good, "-wazuh-alerts="+other, "-wazuh-alerts="+good,
			"-wazuh-alerts="+other, "-wazuh-alerts="+bad},
	} {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("ambiguous multi-feed CLI emitted token: %v / %q", err, output.String())
		}
	}
}

func TestFourWazuhFeedsPreserveDistinctSources(t *testing.T) {
	owner := snapshotFixture(t)
	files := make([]string, maxWazuhAlertFiles)
	for i := range files {
		// Reused raw alert IDs are normal across separately exported files.
		files[i] = writeJSONLEvents(t, "{\"id\":\"same.4\",\"timestamp\":\"2026-10-17T11:00:00+0200\",\"rule\":{\"description\":\"separate Wazuh feed\"}}\n")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- serveDemoWithWazuhAlertsFiles(ctx, "analyst-a", owner,
			"forensic-case-on-disk", files, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("four-file viewer did not start: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("four-file viewer start timeout")
	}
	var rawURL, bearer string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			rawURL = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			bearer = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	endpoint, err := url.Parse(rawURL)
	if err != nil || endpoint.Hostname() != "127.0.0.1" || len(bearer) != 64 {
		t.Fatalf("invalid four-feed viewer endpoint %q", rawURL)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "http://"+endpoint.Host+
		"/case-timeline?case_id=forensic-case-on-disk&limit=20", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	resp.Body.Close()
	if readErr != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("four-feed authorized Timeline: %d %v", resp.StatusCode, readErr)
	}
	for i := range files {
		identity := "forensic-case-on-disk\x00feed-" + string(rune('1'+i)) + "\x00wazuh-same.4"
		sum := sha256.Sum256([]byte(identity))
		entry := "local-jsonl-entry-" + hex.EncodeToString(sum[:])
		if !bytes.Contains(body, []byte(entry)) {
			t.Fatalf("four-feed projection missing authorized feed %d", i+1)
		}
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("four-feed listener did not stop")
	}
}
