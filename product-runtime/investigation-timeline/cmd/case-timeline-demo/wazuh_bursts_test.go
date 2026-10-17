package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestWazuhTemporalBurstsSlidingWindowAndRanking(t *testing.T) {
	at := func(sec int) time.Time {
		return time.Date(2026, 10, 17, 10, 0, sec, 0, time.UTC)
	}
	groups := map[string][]wazuhBurstObservation{
		"001/5710": {
			{At: at(59), Feed: 2}, {At: at(0), Feed: 1},
			{At: at(1), Feed: 1}, {At: at(2), Feed: 2},
			{At: at(2).Add(120 * time.Second), Feed: 3},
		},
		"002/5500": {
			{At: at(5), Feed: 0}, {At: at(6), Feed: 0}, {At: at(7), Feed: 0},
		},
		"003/9999": {{At: at(0), Feed: 1}, {At: at(1), Feed: 2}},
		"004/7777": {{At: at(0), Feed: 1}, {At: at(0).Add(61*time.Second), Feed: 2},
			{At: at(0).Add(122*time.Second), Feed: 3}},
	}
	bursts := deriveWazuhTemporalBursts(groups)
	if len(bursts) != 2 {
		t.Fatalf("need only two qualifying 60-second windows, got %+v", bursts)
	}
	if got := bursts[0]; got.AgentID != "001" || got.RuleID != "5710" ||
		got.AlertCount != 4 || got.DistinctFeeds != 2 ||
		got.FirstUTC != "2026-10-17T10:00:00Z" ||
		got.LastUTC != "2026-10-17T10:00:59Z" || got.WindowSeconds != 60 {
		t.Fatalf("sliding window/count/feed result: %+v", got)
	}
	if got := bursts[1]; got.AgentID != "002" || got.AlertCount != 3 ||
		got.DistinctFeeds != 1 {
		t.Fatalf("single-feed window incorrectly handled: %+v", got)
	}
}

func TestWazuhTemporalBurstsExactBoundary(t *testing.T) {
	base := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	got := deriveWazuhTemporalBursts(map[string][]wazuhBurstObservation{
		"01/42": {{At: base, Feed: 1}, {At: base.Add(30*time.Second), Feed: 1},
			{At: base.Add(60*time.Second), Feed: 2}},
		"02/42": {{At: base, Feed: 1}, {At: base.Add(30*time.Second), Feed: 1},
			{At: base.Add(60*time.Second+time.Nanosecond), Feed: 2}},
	})
	if len(got) != 1 || got[0].AgentID != "01" || got[0].AlertCount != 3 ||
		got[0].DistinctFeeds != 2 {
		t.Fatalf("inclusive 60s boundary not respected: %+v", got)
	}
}

func TestWazuhTemporalBurstsAuthorizedMultifeedEndToEnd(t *testing.T) {
	owner := snapshotFixture(t)
	paths := []string{
		writeJSONLEvents(t,
			"{\"id\":\"a.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"RAW-PRIVATE-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":5,\"description\":\"RAW-PRIVATE-DESC\"}}\n"+
			"{\"id\":\"a.2\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":7}}\n"),
		writeJSONLEvents(t,
			"{\"id\":\"b.1\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12}}\n"+
			"{\"id\":\"b.2\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":15}}\n"),
	}
	data, imported, levels, err := loadMultipleWazuhFeedsWithLevels(
		context.Background(), "analyst-a", owner, "forensic-case-on-disk", paths)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithLevels(context.Background(), "analyst-a",
		data, imported, levels, "forensic-case-on-disk", &out); err != nil {
		t.Fatal(err)
	}
	var report wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if imported != 4 || len(report.TemporalBursts) != 1 {
		t.Fatalf("authorized real-file burst missing: %+v", report.TemporalBursts)
	}
	burst := report.TemporalBursts[0]
	if burst.AgentID != "001" || burst.RuleID != "5710" ||
		burst.AlertCount != 3 || burst.DistinctFeeds != 2 ||
		burst.FirstUTC != "2026-10-17T10:00:00Z" ||
		burst.LastUTC != "2026-10-17T10:00:40Z" {
		t.Fatalf("incorrect authorized burst: %+v", burst)
	}
	for _, private := range []string{"RAW-PRIVATE", "a.1", "b.1", "Bearer ", paths[0]} {
		if strings.Contains(out.String(), private) {
			t.Fatalf("sensitive raw source leaked into Wazuh summary: %q", private)
		}
	}
	var denied bytes.Buffer
	if err := writeDemoWazuh4RuleTriageWithLevels(context.Background(), "analyst-b",
		data, imported, levels, "forensic-case-on-disk", &denied); err == nil || denied.Len() != 0 {
		t.Fatalf("foreign principal received burst: err=%v output=%q", err, denied.String())
	}
}

func TestWazuhTemporalBurstsLimit(t *testing.T) {
	base := time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)
	groups := make(map[string][]wazuhBurstObservation)
	for n := 0; n < 40; n++ {
		key := "000/" + strings.Repeat("1", n%5+1) + string(rune('0'+n/5))
		groups[key] = []wazuhBurstObservation{
			{At: base, Feed: 1}, {At: base.Add(time.Second), Feed: 2},
			{At: base.Add(2*time.Second), Feed: 3},
		}
	}
	got := deriveWazuhTemporalBursts(groups)
	if len(got) != 10 {
		t.Fatalf("top-ten summary limit broken: %d", len(got))
	}
}
