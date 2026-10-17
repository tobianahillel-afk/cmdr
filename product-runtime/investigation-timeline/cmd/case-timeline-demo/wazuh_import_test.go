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
	"strings"
	"testing"
	"time"
)

const wazuh4Fixture = "{\"timestamp\":\"2026-10-17T09:30:15.725+0200\",\"id\":\"1792222215.125\",\"rule\":{\"level\":10,\"description\":\"Wazuh SSH authentication alert\",\"id\":\"5710\"},\"agent\":{\"id\":\"001\",\"name\":\"linux-host\"},\"full_log\":\"private raw event\"}\n" +
	"{\"timestamp\":\"2026-10-17T07:31:16Z\",\"id\":\"1792222276.126\",\"full_log\":\"Wazuh detected file change\",\"location\":\"/var/log/auth.log\"}\n"

func TestNormalizeWazuh4AlertsIntoExistingEventContract(t *testing.T) {
	got, err := normalizeWazuh4Alerts(context.Background(), []byte(wazuh4Fixture))
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(bytes.TrimSuffix(got, []byte("\n")), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected exactly two normalized events, got %d", len(lines))
	}
	var first, second normalizedLocalEvent
	if err := json.Unmarshal(lines[0], &first); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(lines[1], &second); err != nil {
		t.Fatal(err)
	}
	if first.ID != "wazuh-1792222215.125" ||
		first.Timestamp != "2026-10-17T09:30:15.725+02:00" ||
		first.Message != "Wazuh SSH authentication alert" ||
		second.ID != "wazuh-1792222276.126" ||
		second.Message != "Wazuh detected file change" {
		t.Fatalf("source fields not faithfully selected: %+v %+v", first, second)
	}
	if strings.Contains(string(got), "private raw event") ||
		strings.Contains(string(got), "linux-host") ||
		strings.Contains(string(got), "/var/log/auth.log") {
		t.Fatal("unselected Wazuh source fields copied into normalized event")
	}
	owner := demoData()
	merged, err := importLocalEvents(context.Background(), owner, "analyst-a", "case-a", got)
	if err != nil || len(merged.Cases[0].EntryIDs) != len(owner.Cases[0].EntryIDs)+2 {
		t.Fatalf("normalized Wazuh events not importable into authorized Case: %v", err)
	}
}

func TestWazuh4ImportFailsClosed(t *testing.T) {
	valid := strings.TrimSpace(strings.Split(wazuh4Fixture, "\n")[0])
	cases := []struct{ name, input string }{
		{"empty", ""},
		{"blank-line", "\n" + valid},
		{"missing-id", `{"timestamp":"2026-10-17T09:30:15+0200","rule":{"description":"alert"}}`},
		{"invalid-time", `{"timestamp":"yesterday","id":"123","rule":{"description":"alert"}}`},
		{"invalid-id", `{"timestamp":"2026-10-17T09:30:15+0200","id":"../123","rule":{"description":"alert"}}`},
		{"invalid-message", `{"timestamp":"2026-10-17T09:30:15+0200","id":"123","rule":{"description":"bad\nalert"}}`},
		{"no-preview", `{"timestamp":"2026-10-17T09:30:15+0200","id":"123","agent":{"name":"host"}}`},
		{"ecs5-not-supported", `{"@timestamp":"2026-10-17T09:30:15Z","event":{"original":"ECS event"}}`},
		{"invalid-json", "{invalid}"},
		{"oversized-line", strings.Repeat("x", 8193)},
		{"too-many-lines", strings.Repeat(valid+"\n", 5001)},
		{"overlong-preview", `{"timestamp":"2026-10-17T09:30:15+0200","id":"123","rule":{"description":"` + strings.Repeat("a", 513) + `"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := normalizeWazuh4Alerts(context.Background(), []byte(tc.input))
			if err == nil || len(result) != 0 {
				t.Fatalf("unsupported input returned normalized data: %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if result, err := normalizeWazuh4Alerts(ctx, []byte(wazuh4Fixture)); err != context.Canceled || result != nil {
		t.Fatalf("cancelled Wazuh normalization returned data: %v", err)
	}
	owner := snapshotFixture(t)
	events := writeJSONLEvents(t, wazuh4Fixture)
	for _, tc := range [][]string{
		{"-wazuh-alerts=" + events},
		{"-wazuh-alerts=" + events, "-serve-import"},
		{"-wazuh-alerts=" + events, "-serve-import", "-snapshot=" + owner},
		{"-wazuh-alerts=" + events, "-serve-import", "-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-triage-import"},
		{"-wazuh-alerts=" + events, "-serve-import", "-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-import-jsonl=" + events},
		{"-wazuh-alerts=" + events, "-serve-import", "-snapshot=" + owner, "-into-case=case-b"},
		{"-wazuh-alerts=" + events, "-serve-import", "-snapshot=" + owner, "-into-case=forensic-case-on-disk", "-principal=analyst-b"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), tc, &out); err == nil || out.Len() != 0 {
			t.Fatalf("invalid Wazuh import accepted or emitted credential: %v %q", err, out.String())
		}
	}
}

func TestWazuh4DiskAlertsReachBearerAuthorizedCaseAndSource(t *testing.T) {
	ownerPath := snapshotFixture(t)
	alertsPath := writeJSONLEvents(t, wazuh4Fixture)
	ownerBefore, err := os.ReadFile(ownerPath)
	if err != nil {
		t.Fatal(err)
	}
	alertsBefore, err := os.ReadFile(alertsPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
			"-snapshot=" + ownerPath, "-principal=analyst-a",
			"-wazuh-alerts=" + alertsPath, "-into-case=forensic-case-on-disk", "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("Wazuh viewer refused to start: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("Wazuh viewer did not open")
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
	if err != nil || parsed.Hostname() != "127.0.0.1" || len(bearer) != 64 {
		t.Fatalf("viewer not local or no bearer: %q", endpoint)
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
	sum := sha256.Sum256([]byte("forensic-case-on-disk\x00wazuh-1792222215.125"))
	entryID := "local-jsonl-entry-" + hex.EncodeToString(sum[:])
	status, body := get("/case-timeline?case_id=forensic-case-on-disk&limit=20", bearer)
	if status != http.StatusOK || !strings.Contains(body, entryID) ||
		strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "demo-tenant-b") {
		t.Fatalf("Wazuh alert not visible to authorized Case: %d %s", status, body)
	}
	status, body = get("/demo-source-open?case_id=forensic-case-on-disk&entry_id="+url.QueryEscape(entryID), bearer)
	if status != http.StatusOK || !strings.Contains(body, "Wazuh SSH authentication alert") ||
		strings.Contains(body, "private raw event") {
		t.Fatalf("Wazuh source preview failed: %d %s", status, body)
	}
	for _, route := range []string{"/case-timeline?case_id=forensic-case-on-disk",
		"/demo-source-open?case_id=forensic-case-on-disk&entry_id="+url.QueryEscape(entryID)} {
		if status, body := get(route, ""); status != http.StatusUnauthorized || strings.Contains(body, entryID) {
			t.Fatalf("unauthenticated Wazuh data leaked: %d %s", status, body)
		}
	}
	if status, body := get("/case-timeline?case_id=case-b", bearer); status != http.StatusNotFound || strings.Contains(body, entryID) {
		t.Fatalf("foreign Case disclosed Wazuh source: %d %s", status, body)
	}
	if strings.Contains(log.String(), "Wazuh SSH") || strings.Contains(log.String(), entryID) {
		t.Fatal("Wazuh event disclosed in startup output")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Wazuh loopback shutdown: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Wazuh viewer did not stop")
	}
	ownerAfter, err := os.ReadFile(ownerPath)
	if err != nil || !bytes.Equal(ownerBefore, ownerAfter) {
		t.Fatal("Wazuh integration modified owner snapshot")
	}
	alertsAfter, err := os.ReadFile(alertsPath)
	if err != nil || !bytes.Equal(alertsBefore, alertsAfter) {
		t.Fatal("Wazuh integration modified original alerts file")
	}
}

func TestWazuh4ImportBoundaryFailuresCannotExposeViewer(t *testing.T) {
	owner := snapshotFixture(t)
	alerts := writeJSONLEvents(t, wazuh4Fixture)
	invalidAlerts := writeJSONLEvents(t, "{invalid json}\n")
	for _, tc := range []struct {
		name, ownerPath, alertPath, principal, caseID string
	}{
		{"missing-owner", "", alerts, "analyst-a", "forensic-case-on-disk"},
		{"same-file", owner, owner, "analyst-a", "forensic-case-on-disk"},
		{"invalid-owner", alerts, alerts + ".missing", "analyst-a", "forensic-case-on-disk"},
		{"missing-events", owner, alerts + ".missing", "analyst-a", "forensic-case-on-disk"},
		{"invalid-alerts", owner, invalidAlerts, "analyst-a", "forensic-case-on-disk"},
		{"foreign-case", owner, alerts, "analyst-a", "case-b"},
		{"denied-principal", owner, alerts, "analyst-b", "forensic-case-on-disk"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := serveDemoWithWazuhAlerts(context.Background(), tc.principal, tc.ownerPath, tc.caseID, tc.alertPath, &output); err == nil || output.Len() != 0 {
				t.Fatalf("invalid source started listener or exposed bearer: %v / %q", err, output.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := serveDemoWithWazuhAlerts(ctx, "analyst-a", owner, "forensic-case-on-disk", alerts, &output); err != context.Canceled || output.Len() != 0 {
		t.Fatalf("cancelled source started viewer: %v / %q", err, output.String())
	}
}

func TestWazuh4RejectsOversizedAndInvalidUTF8Input(t *testing.T) {
	for _, input := range [][]byte{
		bytes.Repeat([]byte{'x'}, maxDemoSnapshotBytes+1),
		[]byte{0xff},
		[]byte("{\"timestamp\":\"2026-10-17T09:30:15+0200\",\"id\":\"one\",\"rule\":{\"description\":\"ok\"}}\n\n"),
	} {
		got, err := normalizeWazuh4Alerts(context.Background(), input)
		if err == nil || len(got) != 0 {
			t.Fatalf("invalid Wazuh source returned events: %v", err)
		}
	}
}
