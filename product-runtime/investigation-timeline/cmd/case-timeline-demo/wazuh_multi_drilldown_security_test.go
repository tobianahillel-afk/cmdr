package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestMultiWazuhDrilldownEnforcesOwnerAndMetadataIntegrity(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"private first\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"private second\"}}\n")
	data, count, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(), "analyst-a", owner, "forensic-case-on-disk", []string{first, second})
	if err != nil || count != 2 || len(levels) != 2 {
		t.Fatalf("valid Wazuh feed setup: count=%d levels=%d err=%v", count, len(levels), err)
	}
	var out bytes.Buffer
	if err := writeDemoWazuhDrilldownWithLevels(context.Background(), "analyst-a", data, count, levels, "forensic-case-on-disk", "agent_id", "001", &out); err != nil {
		t.Fatal(err)
	}
	var result wazuhDrilldownReceipt
	if err := json.Unmarshal(out.Bytes(), &result); err != nil || result.TotalMatches != 1 || len(result.Entries) != 1 {
		t.Fatalf("authorized agent pivot: %s err=%v", out.String(), err)
	}
	if strings.Contains(out.String(), "private first") || strings.Contains(out.String(), "private second") ||
		strings.Contains(out.String(), "same.1") {
		t.Fatalf("raw source content leaked in drilldown: %s", out.String())
	}
	for _, tc := range []struct {
		name, principal, caseID, selector, value string
		count int
		levels map[string]wazuhRuleLevel
		cancel bool
	}{
		{name: "foreign-principal", principal: "analyst-b", caseID: "forensic-case-on-disk", selector: "agent_id", value: "001", count: count, levels: levels},
		{name: "foreign-case", principal: "analyst-a", caseID: "case-b", selector: "agent_id", value: "001", count: count, levels: levels},
		{name: "missing-metadata", principal: "analyst-a", caseID: "forensic-case-on-disk", selector: "agent_id", value: "001", count: count, levels: map[string]wazuhRuleLevel{}},
		{name: "mismatched-count", principal: "analyst-a", caseID: "forensic-case-on-disk", selector: "agent_id", value: "001", count: count+1, levels: levels},
		{name: "invalid-selector", principal: "analyst-a", caseID: "forensic-case-on-disk", selector: "hostname", value: "001", count: count, levels: levels},
		{name: "invalid-identifier", principal: "analyst-a", caseID: "forensic-case-on-disk", selector: "rule_id", value: "not-numeric", count: count, levels: levels},
		{name: "cancelled", principal: "analyst-a", caseID: "forensic-case-on-disk", selector: "agent_id", value: "001", count: count, levels: levels, cancel: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}
			var response bytes.Buffer
			err := writeDemoWazuhDrilldownWithLevels(ctx, tc.principal, data, tc.count, tc.levels, tc.caseID, tc.selector, tc.value, &response)
			if err == nil || response.Len() != 0 {
				t.Fatalf("owner/metadata denial was not atomic: err=%v body=%s", err, response.String())
			}
			if tc.cancel && !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation did not propagate: %v", err)
			}
		})
	}
}
