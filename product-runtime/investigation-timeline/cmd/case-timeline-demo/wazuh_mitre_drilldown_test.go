package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A technique is an untrusted source label; it never carries Event.read.
// The existing Case -> Timeline -> typed-source owner traversal must run for
// every page before revealing even the authorized synthetic entry IDs.
func TestWazuhMITREDrilldownUsesExistingAuthorizedPagedSourceRoute(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var lines strings.Builder
	for i := 0; i < 53; i++ {
		labels := "[\"T1078\"]"
		if i%2 == 0 {
			labels = "[\"T1059.001\",\"T1059.001\",\"T1078\"]"
		}
		fmt.Fprintf(&lines, "{\"id\":\"secret-%d\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"full_log\":\"SECRET-RAW\",\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-PREVIEW\",\"mitre\":{\"id\":%s}}}\n", i, labels)
	}
	alerts := writeJSONLEvents(t, lines.String())
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(), "analyst-a", owner, "case-a", alerts)
	if err != nil {
		t.Fatal(err)
	}
	token := strings.Repeat("a", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	handler := newDemoWazuhCohortHandler(token, "analyst-a", "case-a", data, imported, raw, base)
	query := func(path, bearer string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w
	}
	seen := make(map[string]bool)
	for _, page := range []struct{ offset, size int; more bool }{
		{0, 25, true}, {25, 2, false},
	} {
		url := fmt.Sprintf("/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&offset=%d", page.offset)
		got := query(url, token)
		if got.Code != http.StatusOK || got.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("authorized page refused: %d %s", got.Code, got.Body.String())
		}
		var receipt wazuhDrilldownReceipt
		if err := json.Unmarshal(got.Body.Bytes(), &receipt); err != nil {
			t.Fatal(err)
		}
		if receipt.Selector != "technique_id" || receipt.Value != "T1059.001" ||
			receipt.TotalMatches != 27 || receipt.Offset != page.offset ||
			len(receipt.Entries) != page.size || receipt.HasMore != page.more {
			t.Fatalf("incorrect technique page: %+v", receipt)
		}
		for _, entry := range receipt.Entries {
			if seen[entry.ID] {
				t.Fatalf("duplicate matching entry: %s", entry.ID)
			}
			seen[entry.ID] = true
		}
		for _, protected := range []string{"SECRET-RAW", "PRIVATE-PREVIEW", "secret-0", "Bearer ", "demo-tenant-b"} {
			if bytes.Contains(got.Body.Bytes(), []byte(protected)) {
				t.Fatalf("sensitive source escaped: %s", protected)
			}
		}
	}
	if len(seen) != 27 {
		t.Fatalf("missing authorized matching source records: %d", len(seen))
	}
	// Top-level technique is independent from rule or agent; duplicate and
	// mixed authority selectors are refused before any source data escapes.
	for _, path := range []string{
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&agent_id=001",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&rule_id=5710",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&technique_id=T1078",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&tenant=tenant-a",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=t1059.001",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.01",
		"/demo-wazuh-drilldown?case_id=case-b&technique_id=T1059.001",
		"/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001&offset=1",
	} {
		got := query(path, token)
		if got.Code != http.StatusNotFound || strings.Contains(got.Body.String(), "T1059") {
			t.Fatalf("forged source selector leaked: %s status=%d body=%s", path, got.Code, got.Body.String())
		}
	}
	if got := query("/demo-wazuh-drilldown?case_id=case-a&technique_id=T1059.001", ""); got.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated technique returned %d", got.Code)
	}
	var denied bytes.Buffer
	if err := writeDemoWazuhDrilldownFilteredPage(context.Background(), "analyst-b",
		data, imported, raw, "case-a", "technique_id", "T1059.001", "", 0, &denied);
		err == nil || denied.Len() != 0 {
		t.Fatalf("denied principal leaked selection: %v %q", err, denied.String())
	}
	var bad bytes.Buffer
	if err := writeDemoWazuhDrilldownFilteredPage(context.Background(), "analyst-a",
		data, imported, raw, "case-a", "technique_id", "T1059.001", "5710", 0, &bad);
		err == nil || bad.Len() != 0 {
		t.Fatalf("malformed cross-selector emitted data: %v %q", err, bad.String())
	}
}
