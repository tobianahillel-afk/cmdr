package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Integration: the catalog and exact pair page must agree after a previously
// authorized first-Case Wazuh source is revoked. An Entry ID is never an ACL.
func TestIntegratedWazuhCatalogRecountsAfterSourceRevocation(t *testing.T) {
	firstPath, firstCase, firstHash, secondPath, secondCase, secondHash := savedWazuhPairPivotFixture(t)
	base := []string{
		"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstHash,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondHash,
	}
	run := func(extra ...string) []byte {
		t.Helper()
		var out bytes.Buffer
		args := append(append([]string{}, base...), extra...)
		if err := runDemoCLI(context.Background(), args, &out); err != nil {
			t.Fatalf("authorized catalog/pair discovery: %v", err)
		}
		for _, secret := range []string{"FIRST PRIVATE", "SECOND PRIVATE", "Bearer "} {
			if strings.Contains(out.String(), secret) {
				t.Fatalf("catalog/pair navigation leaked preview: %s", secret)
			}
		}
		return out.Bytes()
	}
	type catalog struct {
		Cohorts []verifiedWazuhCohortCount `json:"cohorts"`
	}
	readCatalog := func() []verifiedWazuhCohortCount {
		t.Helper()
		var result catalog
		if err := json.Unmarshal(run("-list-saved-wazuh-cohorts-offset", "0"), &result); err != nil {
			t.Fatal(err)
		}
		return result.Cohorts
	}
	before := readCatalog()
	if len(before) != 2 || before[0].RuleID != "5710" || before[0].MatchedPairs != 2 {
		t.Fatalf("unexpected initial cohort inventory: %+v", before)
	}
	var firstPairs wazuhPairPageReceipt
	if err := json.Unmarshal(run("-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", "5710", "-pair-agent-id", "001"), &firstPairs); err != nil {
		t.Fatal(err)
	}
	if len(firstPairs.Pairs) != 2 {
		t.Fatalf("expected two first-Case leads: %+v", firstPairs)
	}
	revokedID := firstPairs.Pairs[0].FirstEntryID
	owner, err := loadDemoSnapshot(firstPath)
	if err != nil { t.Fatal(err) }
	current := owner.Principals[0]
	const excluded = "cohort-source-revoked-other-reader"
	owner.Principals = append(owner.Principals, timeline.MemoryTimelinePrincipal{
		ID: excluded, Tenant: current.Tenant, Environment: current.Environment,
	})
	ref := ""
	for _, event := range owner.Entries {
		if event.CaseID == firstCase && event.Value.ID == revokedID {
			ref = event.Value.SourceRef
			break
		}
	}
	if ref == "" { t.Fatal("selected original first-Case source missing") }
	revoked := false
	for i := range owner.Sources {
		if owner.Sources[i].Identity.CaseID == firstCase &&
			owner.Sources[i].Identity.Reference == ref {
			owner.Sources[i].Readers = []string{excluded}
			revoked = true
		}
	}
	if !revoked { t.Fatal("canonical typed source missing") }
	updated, err := json.Marshal(owner)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(firstPath, updated, 0600); err != nil { t.Fatal(err) }
	for i := 0; i+1 < len(base); i++ {
		if base[i] == "-expect-snapshot-sha256" {
			base[i+1] = fmt.Sprintf("%x", sha256.Sum256(updated))
		}
	}
	after := readCatalog()
	if len(after) != 2 || after[0].RuleID != "8020" ||
		after[0].MatchedPairs != 2 || after[1].RuleID != "5710" ||
		after[1].MatchedPairs != 1 {
		t.Fatalf("revoked source still counted in catalog: %+v", after)
	}
	var filtered wazuhPairPageReceipt
	if err := json.Unmarshal(run("-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", "5710", "-pair-agent-id", "001"), &filtered); err != nil {
		t.Fatal(err)
	}
	if len(filtered.Pairs) != after[1].MatchedPairs ||
		filtered.Pairs[0].FirstEntryID == revokedID {
		t.Fatalf("catalog and pair navigation diverge after revocation: %+v", filtered)
	}
}
