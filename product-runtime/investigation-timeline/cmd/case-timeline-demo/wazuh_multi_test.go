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

func TestMultipleWazuhFilesThroughAuthorizedLoopback(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"123.456\",\"timestamp\":\"2026-10-17T10:01:00+0200\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":8,\"description\":\"Wazuh first feed alert\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"123.456\",\"timestamp\":\"2026-10-17T10:01:01+0200\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"Wazuh second feed alert\"}}\n")
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
	// The browser cohort endpoint must analyze ALL imported feeds, not just
	// the first alert file, with per-source authorization on every request.
	cohortPath := "/demo-wazuh-cohorts?case_id=forensic-case-on-disk"
	status, cohortJSON := request(cohortPath, bearer)
	if status != http.StatusOK {
		t.Fatalf("multi-feed cohort endpoint unavailable: %d %s", status, cohortJSON)
	}
	var cohort wazuhRuleTriageReceipt
	if err := json.Unmarshal([]byte(cohortJSON), &cohort); err != nil {
		t.Fatal(err)
	}
	if cohort.ImportedEntries != 2 || len(cohort.TopRules) != 1 ||
		cohort.TopRules[0].RuleID != "5710" || cohort.TopRules[0].Count != 2 ||
		len(cohort.TopAgents) != 2 || cohort.TopAgents[0].AgentID != "001" ||
		cohort.TopAgents[1].AgentID != "002" ||
		cohort.TopAgents[1].HighestKnownRuleLevel == nil ||
		*cohort.TopAgents[1].HighestKnownRuleLevel != 12 {
		t.Fatalf("multi-feed analyst cohort does not reflect both feeds: %+v", cohort)
	}
	for _, secret := range []string{"Wazuh first feed alert", "Wazuh second feed alert", "123.456", "Bearer "} {
		if strings.Contains(cohortJSON, secret) {
			t.Fatalf("multi-feed cohort exposed source bytes or bearer: %q", secret)
		}
	}
	for _, tc := range []struct {
		path, token string
		want int
	}{
		{cohortPath, "", http.StatusUnauthorized},
		{cohortPath, strings.Repeat("0", 64), http.StatusUnauthorized},
		{"/demo-wazuh-cohorts?case_id=case-b", bearer, http.StatusNotFound},
		{cohortPath + "&case_id=forensic-case-on-disk", bearer, http.StatusNotFound},
	} {
		code, body := request(tc.path, tc.token)
		if code != tc.want || strings.Contains(body, "5710") {
			t.Fatalf("unauthorized multi-feed cohort leaked: status=%d body=%q", code, body)
		}
	}
	// From the multi-source cohort, pivot to the actual authorized entries,
	// then use the existing typed-source preview path for each selected ID.
	drillPath := "/demo-wazuh-drilldown?case_id=forensic-case-on-disk&rule_id=5710"
	status, drillJSON := request(drillPath, bearer)
	if status != http.StatusOK {
		t.Fatalf("multi-feed cohort-to-entry drilldown unavailable: %d %s", status, drillJSON)
	}
	var drill wazuhDrilldownReceipt
	if err := json.Unmarshal([]byte(drillJSON), &drill); err != nil {
		t.Fatal(err)
	}
	if drill.TotalMatches != 2 || len(drill.Entries) != 2 ||
		drill.Selector != "rule_id" || drill.Value != "5710" {
		t.Fatalf("multi-feed drilldown failed to select both alerts: %+v", drill)
	}
	// Rule drilldown joins both feeds. The agent pivot must isolate each
	// specific feed even when raw Wazuh alert IDs collide.
	for _, tc := range []struct {
		agent string
		expectedFeed int
	}{
		{agent: "001", expectedFeed: 1},
		{agent: "002", expectedFeed: 2},
	} {
		code, body := request("/demo-wazuh-drilldown?case_id=forensic-case-on-disk&agent_id="+tc.agent, bearer)
		var agentDrill wazuhDrilldownReceipt
		if code != http.StatusOK || json.Unmarshal([]byte(body), &agentDrill) != nil ||
			agentDrill.TotalMatches != 1 || len(agentDrill.Entries) != 1 ||
			agentDrill.Selector != "agent_id" || agentDrill.Value != tc.agent {
			t.Fatalf("agent %q multi-feed drilldown mismatch: %d %s", tc.agent, code, body)
		}
		identity := "forensic-case-on-disk\x00feed-" + string(rune('0'+tc.expectedFeed)) + "-wazuh-123.456"
		checksum := sha256.Sum256([]byte(identity))
		wanted := "local-jsonl-entry-" + hex.EncodeToString(checksum[:])
		if agentDrill.Entries[0].ID != wanted || strings.Contains(body, "Wazuh first feed alert") ||
			strings.Contains(body, "Wazuh second feed alert") {
			t.Fatalf("agent %q opened wrong/hidden Wazuh source: %s", tc.agent, body)
		}
	}
	for i := range drill.Entries {
		if drill.Entries[i].ID == "" || strings.Contains(drillJSON, "Wazuh first feed alert") ||
			strings.Contains(drillJSON, "Wazuh second feed alert") {
			t.Fatalf("drilldown leaked or omitted source: %s", drillJSON)
		}
	}
	for _, requestData := range []struct {
		path, bearer string
		want int
	}{
		{drillPath, "", http.StatusUnauthorized},
		{drillPath, strings.Repeat("0", 64), http.StatusUnauthorized},
		{"/demo-wazuh-drilldown?case_id=case-b&rule_id=5710", bearer, http.StatusNotFound},
		{drillPath + "&agent_id=001", bearer, http.StatusNotFound},
	} {
		code, body := request(requestData.path, requestData.bearer)
		if code != requestData.want || strings.Contains(body, "local-jsonl-entry-") {
			t.Fatalf("invalid drilldown leaked identity: %d %s", code, body)
		}
	}
	path := "/case-timeline?case_id=forensic-case-on-disk&limit=20"
	status, page := request(path, bearer)
	if status != http.StatusOK {
		t.Fatalf("authorized Case Timeline unavailable: %d %s", status, page)
	}
	for feed, preview := range []string{"Wazuh first feed alert", "Wazuh second feed alert"} {
		identity := "forensic-case-on-disk\x00feed-" + string(rune('1'+feed)) + "-wazuh-123.456"
		digest := sha256.Sum256([]byte(identity))
		entry := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
		if !strings.Contains(page, entry) || !strings.Contains(page, "operator-local-jsonl") {
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
		identity := "forensic-case-on-disk\x00feed-" + string(rune('1'+i)) + "-wazuh-same.4"
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
