package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// End-to-end four-feed behavior checks deterministic ranking, bounded output,
// unknown-vs-zero severity, and the *distinct-feed* criterion.
func TestMultiWazuhCrossFeedSignalRankingAndTopTen(t *testing.T) {
	owner := snapshotFixture(t)
	var first, second, third, fourth strings.Builder
	for i := 100; i < 112; i++ {
		fmt.Fprintf(&first, "{\"id\":\"first-%d\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"%d\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-SOURCE-A\"}}\n", i, i)
		fmt.Fprintf(&second, "{\"id\":\"second-%d\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"%d\"},\"rule\":{\"id\":\"5710\",\"level\":0,\"description\":\"PRIVATE-SOURCE-B\"}}\n", i, i)
	}
	// Same pair, same feed: increases the alert count, never the feed count.
	first.WriteString("{\"id\":\"first-100-extra\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"100\"},\"rule\":{\"id\":\"5710\",\"level\":6,\"description\":\"PRIVATE-REPEAT\"}}\n")
	third.WriteString("{\"id\":\"third-special\",\"timestamp\":\"2026-10-17T10:00:02Z\",\"agent\":{\"id\":\"105\"},\"rule\":{\"id\":\"5710\",\"level\":15,\"description\":\"PRIVATE-THIRD\"}}\n")
	fourth.WriteString("{\"id\":\"fourth-special\",\"timestamp\":\"2026-10-17T10:00:03Z\",\"agent\":{\"id\":\"105\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-FOURTH\"}}\n")
	// Different rule on the same agent appears in only one feed.
	fourth.WriteString("{\"id\":\"fourth-other-rule\",\"timestamp\":\"2026-10-17T10:00:04Z\",\"agent\":{\"id\":\"105\"},\"rule\":{\"id\":\"5500\",\"description\":\"PRIVATE-UNRELATED\"}}\n")
	paths := []string{
		writeJSONLEvents(t, first.String()),
		writeJSONLEvents(t, second.String()),
		writeJSONLEvents(t, third.String()),
		writeJSONLEvents(t, fourth.String()),
	}
	data, count, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(),
		"analyst-a", owner, "forensic-case-on-disk", paths)
	if err != nil || count != 28 || len(levels) != 28 {
		t.Fatalf("four-feed import count: %d levels %d error %v", count, len(levels), err)
	}
	var out bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithLevels(context.Background(), "analyst-a",
		data, count, levels, "forensic-case-on-disk", &out); err != nil {
		t.Fatal(err)
	}
	var got wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.MultiFeedSignals) != 10 {
		t.Fatalf("top-ten bound not respected: %+v", got.MultiFeedSignals)
	}
	strongest := got.MultiFeedSignals[0]
	if strongest.AgentID != "105" || strongest.RuleID != "5710" ||
		strongest.DistinctFeeds != 4 || strongest.AlertCount != 4 ||
		strongest.HighestKnownRuleLevel == nil || *strongest.HighestKnownRuleLevel != 15 {
		t.Fatalf("four-feed priority signal is wrong: %+v", strongest)
	}
	secondSignal := got.MultiFeedSignals[1]
	if secondSignal.AgentID != "100" || secondSignal.DistinctFeeds != 2 ||
		secondSignal.AlertCount != 3 || secondSignal.HighestKnownRuleLevel == nil ||
		*secondSignal.HighestKnownRuleLevel != 6 {
		t.Fatalf("alert-frequency tie break failed: %+v", secondSignal)
	}
	for i, signal := range got.MultiFeedSignals[2:] {
		if signal.AgentID == "105" || signal.RuleID != "5710" ||
			signal.DistinctFeeds != 2 || signal.AlertCount != 2 ||
			signal.HighestKnownRuleLevel == nil || *signal.HighestKnownRuleLevel != 0 {
			t.Fatalf("unexpected ranked signal %d: %+v", i+2, signal)
		}
	}
	if strings.Contains(out.String(), "PRIVATE-") || strings.Contains(out.String(), "first-100-extra") {
		t.Fatal("private source content leaked in ranked Wazuh signal")
	}
}
