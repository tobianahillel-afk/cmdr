package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestMultipleWazuhFeedsToTwoAuthorizedCaseLeads(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"rule\":{\"description\":\"secret from sensor A\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.130Z\",\"rule\":{\"description\":\"secret from sensor B\"}}\n")
	before := make(map[string][]byte)
	for _, path := range []string{owner, first, second} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[path] = data
	}
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner,
		"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
		"-into-case=case-a", "-triage-import", "-with-case=case-c",
		"-compare-proximity", "-proximity-seconds=1",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var receipt localImportedProximity
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.Format != "cmdr-local-multi-wazuh4-case-proximity-v1" ||
		!receipt.Synthetic || receipt.ImportedEntries != 2 ||
		receipt.Proximity.FirstCase != "case-a" || receipt.Proximity.SecondCase != "case-c" ||
		receipt.Proximity.WindowSeconds != 1 {
		t.Fatalf("invalid multi-feed Case comparison: %+v", receipt)
	}
	for feed, gap := range []int64{125, 120} {
		identity := "case-a\x00feed-" + string(rune('1'+feed)) + "-wazuh-same.123"
		digest := sha256.Sum256([]byte(identity))
		entry := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
		found := false
		for _, lead := range receipt.Proximity.Leads {
			if lead.FirstEntryID == entry && lead.SecondEntryID == "event-c" &&
				lead.SeparationMS == gap && lead.FirstKind == "observed" &&
				lead.SecondKind == "inferred" {
				found = true
			}
		}
		if !found {
			t.Fatalf("Wazuh feed %d absent from authorized nearest-time leads: %+v", feed+1, receipt.Proximity.Leads)
		}
	}
	for _, secret := range []string{"secret from sensor", "demo-hidden-source", "demo-tenant-b",
		"Bearer ", "principals\"", "readers\"", "source_ref\""} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("sensitive source/owner data leaked to cross-Case receipt: %s", secret)
		}
	}
	for path, original := range before {
		after, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(after, original) {
			t.Fatalf("operator-owned file changed: %v", err)
		}
	}
}

func TestMultipleWazuhTemporalComparisonFailsClosed(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	first := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"rule\":{\"description\":\"confidential\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02.130Z\",\"rule\":{\"description\":\"confidential 2\"}}\n")
	bad := writeJSONLEvents(t, "{invalid-json}\n")
	for _, tc := range []struct {
		name, principal, secondCase string
		files                       []string
		seconds                     int
	}{
		{"foreign-second-Case", "analyst-a", "case-b", []string{first, second}, 1},
		{"unknown-second-Case", "analyst-a", "missing", []string{first, second}, 1},
		{"same-Case", "analyst-a", "case-a", []string{first, second}, 1},
		{"denied-principal", "analyst-b", "case-c", []string{first, second}, 1},
		{"malformed-second-feed", "analyst-a", "case-c", []string{first, bad}, 1},
		{"duplicate-file", "analyst-a", "case-c", []string{first, first}, 1},
		{"owner-is-alert-file", "analyst-a", "case-c", []string{first, owner}, 1},
		{"empty-input", "analyst-a", "case-c", nil, 1},
		{"single-file-to-multi-helper", "analyst-a", "case-c", []string{first}, 1},
		{"too-many-feeds", "analyst-a", "case-c", []string{first, second, first, second, bad}, 1},
		{"zero-window", "analyst-a", "case-c", []string{first, second}, 0},
		{"overlarge-window", "analyst-a", "case-c", []string{first, second}, 301},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := writeDemoMultiWazuhProximity(context.Background(), tc.principal, owner,
				"case-a", tc.secondCase, tc.files, tc.seconds, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid multi-feed comparison disclosed report: %v / %q", err, out.String())
			}
		})
	}
	for _, args := range [][]string{
		{"-snapshot=" + owner, "-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
			"-into-case=case-a", "-triage-import", "-compare-proximity"},
		{"-snapshot=" + owner, "-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
			"-into-case=case-a", "-triage-import", "-with-case=case-c", "-compare-proximity", "-serve-import"},
		{"-snapshot=" + owner, "-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
			"-into-case=case-a", "-triage-import", "-with-case=case-c", "-compare-proximity", "-compare-events"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("ambiguous multi-feed comparison produced output: %v / %q", err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	err := writeDemoMultiWazuhProximity(ctx, "analyst-a", owner, "case-a",
		"case-c", []string{first, second}, 1, &out)
	if !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled multi-feed comparison exposed output: %v", err)
	}
}

func TestWazuhFeedNamespaceRejectsInvalidImport(t *testing.T) {
	if _, err := namespaceWazuhEvents(context.Background(),
		[]byte("{\"id\":\"good\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"message\":\"ok\"}\n"), 0); err == nil {
		t.Fatal("feed zero accepted")
	}
	if _, err := namespaceWazuhEvents(context.Background(), []byte("{invalid}\n"), 1); err == nil {
		t.Fatal("malformed normalized event accepted")
	}
	if _, err := namespaceWazuhEvents(context.Background(), []byte("{\"id\":\"!bad\"}\n"), 1); err == nil {
		t.Fatal("unsupported source ID accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := namespaceWazuhEvents(ctx, []byte("{}\n"), 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled feed normalization continued: %v", err)
	}
}
