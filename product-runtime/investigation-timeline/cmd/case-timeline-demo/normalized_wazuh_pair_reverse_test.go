package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type reverseWazuhPairReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	SecondEntryID string `json:"second_entry_id"`
	Offset int `json:"offset"`
	HasMore bool `json:"has_more"`
	Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
}

func TestReverseSecondWazuhPairPivotReopensAuthorizedPrivateDossier(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPageFixture(t)
	base := []string{"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondSHA}
	var original bytes.Buffer
	if err := runDemoCLI(context.Background(), append(append([]string{}, base...),
		"-list-saved-wazuh-pair-offset", "0"), &original); err != nil {
		t.Fatal(err)
	}
	var initial wazuhPairPageReceipt
	if err := json.Unmarshal(original.Bytes(), &initial); err != nil { t.Fatal(err) }
	if len(initial.Pairs) != 10 { t.Fatal("expected ten original leads") }
	peerID := initial.Pairs[0].SecondEntryID

	read := func(offset string) reverseWazuhPairReceipt {
		t.Helper()
		args := append(append([]string{}, base...), "-find-saved-pair-second-entry", peerID,
			"-find-saved-pair-second-offset", offset)
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err != nil { t.Fatal(err) }
		for _, forbidden := range []string{"PRIVATE FIRST", "PRIVATE SECOND", "Bearer ", "127.0.0.1:"} {
			if strings.Contains(out.String(), forbidden) { t.Fatal("private data leaked in reverse discovery") }
		}
		var receipt reverseWazuhPairReceipt
		if err := json.Unmarshal(out.Bytes(), &receipt); err != nil { t.Fatal(err) }
		if receipt.Format != "cmdr-authorized-wazuh-second-entry-pairs-v1" ||
			!receipt.Synthetic || receipt.SecondEntryID != peerID {
			t.Fatalf("wrong reverse selection: %+v", receipt)
		}
		for _, pair := range receipt.Pairs {
			if pair.SecondEntryID != peerID { t.Fatalf("foreign second Entry: %+v", pair) }
		}
		return receipt
	}
	first, second, exhausted := read("0"), read("10"), read("12")
	if first.Offset != 0 || len(first.Pairs) != 10 || !first.HasMore ||
		second.Offset != 10 || len(second.Pairs) != 2 || second.HasMore ||
		exhausted.Offset != 12 || len(exhausted.Pairs) != 0 || exhausted.HasMore {
		t.Fatalf("reverse pair pages do not expose all leads: %+v / %+v / %+v", first, second, exhausted)
	}
	if first.Pairs[9].FirstEntryID == second.Pairs[0].FirstEntryID {
		t.Fatal("same first Entry repeated across reverse pages")
	}
	// A second-side selection is only navigation. Reopen the private
	// source evidence with both independently authorized exact Entry IDs.
	private := filepath.Join(t.TempDir(), "authorized-reverse-pair.json")
	var stdout bytes.Buffer
	args := append(append([]string{}, base...),
		"-inspect-saved-pair-first-entry", second.Pairs[0].FirstEntryID,
		"-inspect-saved-pair-second-entry", peerID, "-private-report", private)
	if err := runDemoCLI(context.Background(), args, &stdout); err != nil { t.Fatal(err) }
	if stdout.Len() != 0 { t.Fatal("private dossier escaped onto stdout") }
	payload, err := os.ReadFile(private)
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(payload, []byte("PRIVATE FIRST 10")) ||
		!bytes.Contains(payload, []byte("PRIVATE SECOND")) {
		t.Fatal("reverse lead did not reopen actual source-authorized evidence")
	}
}

func TestReverseSecondWazuhPairPivotRejectsUnsafeAndMixedModes(t *testing.T) {
	aPath, aCase, aSHA, bPath, bCase, bSHA := savedWazuhPairPageFixture(t)
	base := []string{"-principal", "analyst-a", "-snapshot", aPath,
		"-compare-saved-case", aCase, "-expect-snapshot-sha256", aSHA,
		"-compare-saved-with-snapshot", bPath,
		"-compare-saved-with-case", bCase, "-expect-other-snapshot-sha256", bSHA}
	var page bytes.Buffer
	if err := runDemoCLI(context.Background(), append(append([]string{}, base...),
		"-list-saved-wazuh-pair-offset", "0"), &page); err != nil { t.Fatal(err) }
	var pairs wazuhPairPageReceipt
	if err := json.Unmarshal(page.Bytes(), &pairs); err != nil { t.Fatal(err) }
	good := append(append([]string{}, base...),
		"-find-saved-pair-second-entry", pairs.Pairs[0].SecondEntryID)
	variants := []struct {
		name string
		args []string
	}{
		{"empty second id", append(append([]string{}, base...), "-find-saved-pair-second-entry", "")},
		{"unknown second id", append(append([]string{}, base...), "-find-saved-pair-second-entry", "absent-id")},
		{"negative offset", append(append([]string{}, good...), "-find-saved-pair-second-offset", "-1")},
		{"oversized offset", append(append([]string{}, good...), "-find-saved-pair-second-offset", "5001")},
		{"offset without second id", append(append([]string{}, base...), "-find-saved-pair-second-offset", "0")},
		{"mixed first discovery", append(append([]string{}, good...), "-find-saved-pair-first-entry", pairs.Pairs[0].FirstEntryID)},
		{"mixed original page", append(append([]string{}, good...), "-list-saved-wazuh-pair-offset", "10")},
		{"mixed private selector", append(append([]string{}, good...), "-inspect-saved-pair-first-entry", pairs.Pairs[0].FirstEntryID)},
		{"mixed source previews", append(append([]string{}, good...), "-include-pair-generic-evidence")},
		{"mixed private report", append(append([]string{}, good...), "-private-report", filepath.Join(t.TempDir(), "deny.json"))},
	}
	for _, tc := range variants {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := runDemoCLI(context.Background(), tc.args, &out); err == nil || out.Len() != 0 {
				t.Fatalf("reverse discovery accepted unsafe mode: %v %q", err, out.String())
			}
		})
	}
	for _, mutation := range []struct{name, flag, value string}{
		{"foreign principal", "-principal", "foreign-analyst"},
		{"other digest", "-expect-other-snapshot-sha256", strings.Repeat("0", 64)},
		{"first digest", "-expect-snapshot-sha256", strings.Repeat("0", 64)},
		{"second Case", "-compare-saved-with-case", "foreign-case"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			args := append([]string{}, good...)
			for i := 0; i+1 < len(args); i++ {
				if args[i] == mutation.flag { args[i+1] = mutation.value; break }
			}
			var out bytes.Buffer
			if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
				t.Fatalf("unauthorized reverse discovery emitted metadata: %v %q", err, out.String())
			}
		})
	}
}

func TestReverseSecondWazuhPairHonorsNearestPeerAndAmbiguousIDs(t *testing.T) {
	key := authorizedWazuhCohort{rule:"5710", agent:"001"}
	other := authorizedWazuhCohort{rule:"9999", agent:"001"}
	now := time.Date(2026, 10, 18, 10, 0, 0, 0, time.UTC)
	left := map[authorizedWazuhCohort][]authorizedWazuhObservation{
		key: {{id:"left", at:now}},
	}
	right := map[authorizedWazuhCohort][]authorizedWazuhObservation{
		key: {{id:"target", at:now.Add(20*time.Second)}, {id:"closer", at:now.Add(time.Second)}},
	}
	// Filtering second events before nearest-match computation would create
	// a false pair with "target". The full authorized cohort must still win.
	pairs, more, err := verifiedWazuhPairsForSecondEntry(context.Background(), left, right, "target", 0)
	if err != nil || more || len(pairs) != 0 {
		t.Fatalf("non-nearest second ID manufactured a pair: %+v %v %v", pairs, more, err)
	}
	right[other] = []authorizedWazuhObservation{{id:"target", at:now}}
	if _, _, err := verifiedWazuhPairsForSecondEntry(context.Background(), left, right, "target", 0); err == nil {
		t.Fatal("ambiguous right-side Entry ID accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := verifiedWazuhPairsForSecondEntry(ctx, left, right, "target", 0); err == nil {
		t.Fatal("cancelled investigation succeeded")
	}
}
