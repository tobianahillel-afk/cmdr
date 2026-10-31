package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMultiWazuhCrossFeedSignalOpensOnlyAuthorizedMatchingEvents(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t,
		"{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-FIRST\"}}\n"+
		"{\"id\":\"other.1\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":8,\"description\":\"PRIVATE-OTHER\"}}\n")
	second := writeJSONLEvents(t,
		"{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:02Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-SECOND\"}}\n")
	data, count, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(), "analyst-a", owner, "forensic-case-on-disk", []string{first, second})
	if err != nil || count != 3 {
		t.Fatalf("load two distinct feeds: count=%d err=%v", count, err)
	}
	token := strings.Repeat("d", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	handler := newDemoWazuhCohortReceiptHandler(token, "analyst-a", "forensic-case-on-disk", data, count, nil, levels, base,
		func(context.Context, io.Writer) error { return nil })
	request := func(path, bearer string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest("GET", path, nil)
		if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }
		got := httptest.NewRecorder()
		handler.ServeHTTP(got, req)
		return got
	}
	path := "/demo-wazuh-drilldown?case_id=forensic-case-on-disk&agent_id=001&rule_id=5710"
	got := request(path, token)
	var receipt wazuhDrilldownReceipt
	if got.Code != http.StatusOK || json.Unmarshal(got.Body.Bytes(), &receipt) != nil ||
		receipt.Selector != "agent_rule" || receipt.Value != "001" ||
		receipt.SecondaryValue != "5710" || receipt.TotalMatches != 2 || len(receipt.Entries) != 2 {
		t.Fatalf("two-feed owner-authorized signal not actionable: %d %s", got.Code, got.Body.String())
	}
	for _, forbidden := range []string{"same.1", "PRIVATE-FIRST", "PRIVATE-OTHER", "PRIVATE-SECOND", token, "PRIVATE-HOST"} {
		if bytes.Contains(got.Body.Bytes(), []byte(forbidden)) {
			t.Fatalf("raw source or bearer exposed: %q", forbidden)
		}
	}
	// Same rule on a different agent is not part of the chosen intersection.
	other := request("/demo-wazuh-drilldown?case_id=forensic-case-on-disk&agent_id=002&rule_id=5710", token)
	var otherReceipt wazuhDrilldownReceipt
	if other.Code != http.StatusOK || json.Unmarshal(other.Body.Bytes(), &otherReceipt) != nil ||
		otherReceipt.TotalMatches != 1 || len(otherReceipt.Entries) != 1 {
		t.Fatalf("agent/rule AND semantics incorrect: %d %s", other.Code, other.Body.String())
	}
	for _, tc := range []struct {
		path, bearer string
		want int
	}{
		{path, "", http.StatusUnauthorized},
		{path, strings.Repeat("e", 64), http.StatusUnauthorized},
		{"/demo-wazuh-drilldown?case_id=foreign-case&agent_id=001&rule_id=5710", token, http.StatusNotFound},
		{path + "&agent_id=002", token, http.StatusNotFound},
		{path + "&rule_id=5500", token, http.StatusNotFound},
		{path + "&unknown=value", token, http.StatusNotFound},
		{"/demo-wazuh-drilldown?case_id=forensic-case-on-disk&agent_id=001&rule_id=%2A", token, http.StatusNotFound},
	}{
		got := request(tc.path, tc.bearer)
		if got.Code != tc.want || bytes.Contains(got.Body.Bytes(), []byte("local-wazuh4-drilldown")) {
			t.Fatalf("invalid selection revealed data: %s %d %s", tc.path, got.Code, got.Body.String())
		}
	}
	// Revocation of one metadata/owner entry fails closed with no partial result.
	delete(levels, receipt.Entries[0].ID)
	revoked := request(path, token)
	if revoked.Code != http.StatusNotFound || bytes.Contains(revoked.Body.Bytes(), []byte("local-wazuh4-drilldown")) {
		t.Fatalf("metadata mismatch returned partial source details: %d %s", revoked.Code, revoked.Body.String())
	}
}

func TestMultiWazuhCrossFeedSignalRejectsInvalidSelectorWithoutOutput(t *testing.T) {
	data := demoData()
	for _, tc := range []struct { selector, first, second string }{
		{"agent_rule", "001", "*"},
		{"agent_rule", "001", ""},
		{"agent_rule", "<host>", "5710"},
		{"agent_id", "001", "5710"},
		{"unknown", "001", ""},
	}{
		var out bytes.Buffer
		err := writeDemoWazuhDrilldownWithLevelsFilter(context.Background(), "analyst-a", data, 1, nil, "case-a", tc.selector, tc.first, tc.second, &out)
		if err == nil || out.Len() != 0 {
			t.Fatalf("invalid signal accepted: %q %s", tc.selector, out.String())
		}
	}
}
