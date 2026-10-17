package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// A real three-file Wazuh Case import becomes an authorized cross-feed summary.
func TestMultiWazuhAuthorizedCrossFeedAgentRuleSignal(t *testing.T) {
	owner := snapshotFixture(t)
	files := []string{
		writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"PRIVATE-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":7,\"description\":\"PRIVATE-FIRST\"}}\n"+
			"{\"id\":\"a.2\",\"timestamp\":\"2026-10-17T10:00:02Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":8,\"description\":\"PRIVATE-FIRST-2\"}}\n"),
		writeJSONLEvents(t, "{\"id\":\"b.1\",\"timestamp\":\"2026-10-17T10:00:03Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-SECOND\"}}\n"),
		writeJSONLEvents(t, "{\"id\":\"c.1\",\"timestamp\":\"2026-10-17T10:00:04Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":15,\"description\":\"PRIVATE-THIRD\"}}\n"),
	}
	data, imported, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(),
		"analyst-a", owner, "forensic-case-on-disk", files)
	if err != nil || imported != 4 || len(levels) != 4 {
		t.Fatalf("three-feed import unavailable: %v count=%d metadata=%d", err, imported, len(levels))
	}
	var output bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithLevels(context.Background(), "analyst-a",
		data, imported, levels, "forensic-case-on-disk", &output); err != nil {
		t.Fatal(err)
	}
	var result wazuhRuleTriageReceipt
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.MultiFeedSignals) != 1 {
		t.Fatalf("incorrect cross-feed group count: %+v", result.MultiFeedSignals)
	}
	signal := result.MultiFeedSignals[0]
	if signal.AgentID != "001" || signal.RuleID != "5710" || signal.AlertCount != 3 ||
		signal.DistinctFeeds != 2 || signal.HighestKnownRuleLevel == nil ||
		*signal.HighestKnownRuleLevel != 12 {
		t.Fatalf("cross-feed grouping fabricated a match: %+v", signal)
	}
	for _, forbidden := range []string{"PRIVATE-", "a.1", "b.1", "c.1",
		"PRIVATE-HOST", "Bearer ", files[0], files[1], files[2]} {
		if strings.Contains(output.String(), forbidden) {
			t.Fatalf("raw source or credential in signal receipt: %q", forbidden)
		}
	}
	var denied bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithLevels(context.Background(), "analyst-b",
		data, imported, levels, "forensic-case-on-disk", &denied); err == nil || denied.Len() != 0 {
		t.Fatalf("foreign principal read cross-feed signal: err=%v output=%q", err, denied.String())
	}
}

func TestSingleFeedDoesNotFabricateMultiFeedSignal(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"x.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":8,\"description\":\"visible single one\"}}\n"+
		"{\"id\":\"x.2\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"visible single two\"}}\n")
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(),
		"analyst-a", owner, "forensic-case-on-disk", first)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithData(context.Background(), "analyst-a",
		data, imported, raw, "forensic-case-on-disk", &output); err != nil {
		t.Fatal(err)
	}
	var result wazuhRuleTriageReceipt
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.ImportedEntries != 2 || len(result.MultiFeedSignals) != 0 {
		t.Fatalf("same source repeated alerts must not be cross-feed: %+v", result.MultiFeedSignals)
	}
}
