package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestMultiWazuhFeedNamespaceRejectsInvalidSourceIdentity(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw string
		feed int
	}{
		{"feed-out-of-range", "{\"id\":\"wazuh-1\"}\n", maxWazuhAlertFiles + 1},
		{"feed-zero", "{\"id\":\"wazuh-1\"}\n", 0},
		{"malformed-normalized", "{invalid-json}\n", 1},
		{"oversized-source-identity", "{\"id\":\"" + strings.Repeat("a", 100) + "\"}\n", 1},
		{"unsupported-source-identity", "{\"id\":\"bad/identity\"}\n", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := namespaceWazuhEvents(context.Background(), []byte(tc.raw), tc.feed)
			if err == nil || len(result) != 0 {
				t.Fatalf("invalid namespaced feed accepted: %v / %s", err, result)
			}
		})
	}
	valid := []byte("{\"id\":\"wazuh-shared.1\",\"timestamp\":\"2026-10-01T10:02:02Z\",\"message\":\"authorized test\"}\n")
	a, err := namespaceWazuhEvents(context.Background(), valid, 1)
	if err != nil { t.Fatal(err) }
	b, err := namespaceWazuhEvents(context.Background(), valid, 2)
	if err != nil { t.Fatal(err) }
	var first, second normalizedLocalEvent
	if err := json.Unmarshal(bytes.TrimSpace(a), &first); err != nil { t.Fatal(err) }
	if err := json.Unmarshal(bytes.TrimSpace(b), &second); err != nil { t.Fatal(err) }
	if first.ID != "feed-1-wazuh-shared.1" || second.ID != "feed-2-wazuh-shared.1" ||
		first.ID == second.ID || first.Message != second.Message {
		t.Fatalf("feed identity collision or altered normalized payload: %+v %+v", first, second)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if data, err := namespaceWazuhEvents(ctx, valid, 1); !errors.Is(err, context.Canceled) || data != nil {
		t.Fatalf("cancelled feed transformation returned bytes: %v", err)
	}
}

func TestMultiWazuhLoaderFailsBeforeAnyOutputOnUnavailableFiles(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	good := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02Z\",\"rule\":{\"description\":\"private\"}}\n")
	other := writeJSONLEvents(t, "{\"id\":\"b.1\",\"timestamp\":\"2026-10-01T10:02:03Z\",\"rule\":{\"description\":\"private\"}}\n")
	missing := good + ".does-not-exist"
	for _, tc := range []struct {
		name, principal, caseID, snapshot string
		files []string
	}{
		{"no-snapshot", "analyst-a", "case-a", "", []string{good, other}},
		{"unknown-principal", "not-a-reader", "case-a", owner, []string{good, other}},
		{"missing-Case", "analyst-a", "unknown-case", owner, []string{good, other}},
		{"missing-second-file", "analyst-a", "case-a", owner, []string{good, missing}},
		{"empty-file-selection", "analyst-a", "case-a", owner, nil},
		{"single-file-selection", "analyst-a", "case-a", owner, []string{good}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, count, err := loadMultipleWazuhFeeds(context.Background(), tc.principal, tc.snapshot, tc.caseID, tc.files)
			if err == nil || count != 0 || len(data.Entries) != 0 {
				t.Fatalf("invalid files or tenant accepted: %v / count=%d", err, count)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, count, err := loadMultipleWazuhFeeds(ctx, "analyst-a", owner, "case-a", []string{good, other}); !errors.Is(err, context.Canceled) || count != 0 {
		t.Fatalf("cancelled multi-file loading continued: %v", err)
	}
}

func TestMultiWazuhComparisonRejectsInvalidSelectionBeforeSourceRead(t *testing.T) {
	for _, tc := range []struct{ first, second string }{
		{"", "case-c"},
		{"case-a", ""},
		{"case-a", "case-a"},
		{"case-a", strings.Repeat("x", 1025)},
	} {
		var out bytes.Buffer
		err := writeDemoMultiWazuhComparison(context.Background(), "analyst-a", "missing-owner", tc.first, tc.second, nil, &out)
		if err == nil || out.Len() > 0 {
			t.Fatalf("invalid comparison selector disclosed output: %v", err)
		}
	}
}

func TestMultiWazuhLoaderRejectsHardlinkedFeeds(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	first := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02Z\",\"rule\":{\"description\":\"private\"}}\n")
	alias := first + "-hardlink"
	if err := os.Link(first, alias); err != nil {
		t.Skipf("hardlinks unavailable on CI filesystem: %v", err)
	}
	_, count, err := loadMultipleWazuhFeeds(context.Background(), "analyst-a", owner, "case-a", []string{first, alias})
	if err == nil || count != 0 {
		t.Fatalf("same underlying source file counted as independent feed: %v", err)
	}
}
