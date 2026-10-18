package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Exact Wazuh Entry selectors are evidence locators, NOT persisted access
// grants. Even a previously legitimate pair and a CORRECT snapshot SHA-256
// must fail when the selected canonical typed source's ACL is revoked.
// This tests the integrated new selector, not the legacy ordinal path.
func TestFullReleaseExactWazuhPairIDsFailClosedAfterSourceRevocation(t *testing.T) {
	firstPath, firstCase, _, secondPath, secondCase, secondHash, pair :=
		makeSavedPairDrilldownOwners(t)

	first, err := loadDemoSnapshot(firstPath)
	if err != nil { t.Fatal(err) }
	if len(first.Principals) != 1 { t.Fatal("saved owner must be scoped") }
	current := first.Principals[0]
	const unrelated = "other-reader-no-longer-sharing-wazuh-source"
	first.Principals = append(first.Principals, timeline.MemoryTimelinePrincipal{
		ID: unrelated, Tenant: current.Tenant, Environment: current.Environment,
	})

	var firstSourceRef string
	for _, entry := range first.Entries {
		if entry.CaseID == firstCase && entry.Value.ID == pair.FirstEntryID {
			firstSourceRef = entry.Value.SourceRef
			break
		}
	}
	if firstSourceRef == "" { t.Fatal("first Wazuh correlated source missing") }
	revoked := false
	for i := range first.Sources {
		if first.Sources[i].Identity.Reference == firstSourceRef &&
			first.Sources[i].Identity.CaseID == firstCase {
			first.Sources[i].Readers = []string{unrelated}
			revoked = true
		}
	}
	if !revoked { t.Fatal("first owner's canonical typed source missing") }
	if _,err := newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",first);err != nil {
		t.Fatalf("changed ACL must leave the owner dataset structurally valid: %v",err)
	}
	updated,err := json.Marshal(first)
	if err != nil { t.Fatal(err) }
	if err:=os.WriteFile(firstPath,updated,0600);err!=nil {t.Fatal(err)}
	// Simulate a valid operator-authorized snapshot *after* ACL revocation,
	// not an obsolete digest check that could conceal an ACL bypass.
	newFirstHash := fmt.Sprintf("%x",sha256.Sum256(updated))
	dest := filepath.Join(t.TempDir(),"must-not-write-private-pair.json")
	var stdout bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-principal","analyst-a",
		"-snapshot",firstPath,
		"-compare-saved-case",firstCase,
		"-expect-snapshot-sha256",newFirstHash,
		"-compare-saved-with-snapshot",secondPath,
		"-compare-saved-with-case",secondCase,
		"-expect-other-snapshot-sha256",secondHash,
		"-inspect-saved-pair-first-entry",pair.FirstEntryID,
		"-inspect-saved-pair-second-entry",pair.SecondEntryID,
		"-include-pair-generic-evidence",
		"-private-report",dest,
	},&stdout)
	if err==nil { t.Fatal("revoked exact Wazuh pair Entry selectors yielded private evidence") }
	if stdout.Len()!=0 {
		t.Fatalf("denied exact Entry selector disclosed private evidence: %q",stdout.String())
	}
	if _,statErr:=os.Lstat(dest);!os.IsNotExist(statErr) {
		t.Fatalf("denied exact pair inspection left sensitive file: %v",statErr)
	}
	// The independent second snapshot has not lost its legitimate source
	// because of the first owner's revocation; the denial cannot be
	// explained away as a malformed second Case or changed hash.
	if _,_,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",
		secondPath,secondCase,secondHash);err!=nil {
		t.Fatalf("independently authorized second Case unexpectedly unavailable: %v",err)
	}
}

// Verify the formerly divergent PR #113 source-revocation safety boundary
// still governs BOTH integrated PR #120 navigation paths on the SAME Cases.
// Entry IDs are mere locators, not authorization after a source ACL change.
func TestIntegratedWazuhCohortAndReverseDiscoveryAfterFirstSourceRevocation(t *testing.T) {
	firstPath, firstCase, firstHash, secondPath, secondCase, secondHash :=
		savedWazuhPairPivotFixture(t)
	base := []string{
		"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstHash,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondHash,
	}
	run := func(args ...string) []byte {
		t.Helper()
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), append(append([]string{}, base...), args...), &output); err != nil {
			t.Fatalf("authorized Case navigation failed: %v", err)
		}
		if bytes.Contains(output.Bytes(), []byte("PRIVATE")) || bytes.Contains(output.Bytes(), []byte("Bearer ")) {
			t.Fatal("Wazuh navigation disclosed private preview or token")
		}
		return output.Bytes()
	}
	var before wazuhPairPageReceipt
	if err := json.Unmarshal(run(
		"-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", "5710", "-pair-agent-id", "001",
	), &before); err != nil { t.Fatal(err) }
	if len(before.Pairs) != 2 { t.Fatalf("expected two authorized initial pairs: %+v", before) }
	revokedID := before.Pairs[0].FirstEntryID
	secondID := before.Pairs[0].SecondEntryID

	owner, err := loadDemoSnapshot(firstPath)
	if err != nil { t.Fatal(err) }
	current := owner.Principals[0]
	const otherPrincipal = "source-revocation-owner-only"
	owner.Principals = append(owner.Principals, timeline.MemoryTimelinePrincipal{
		ID: otherPrincipal, Tenant: current.Tenant, Environment: current.Environment,
	})
	ref := ""
	for _, entry := range owner.Entries {
		if entry.CaseID == firstCase && entry.Value.ID == revokedID {
			ref = entry.Value.SourceRef
		}
	}
	if ref == "" { t.Fatal("selected original owner Entry missing") }
	revoked := false
	for i := range owner.Sources {
		if owner.Sources[i].Identity.CaseID == firstCase &&
			owner.Sources[i].Identity.Reference == ref {
			owner.Sources[i].Readers = []string{otherPrincipal}
			revoked = true
		}
	}
	if !revoked { t.Fatal("selected typed source missing") }
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), "analyst-a", owner); err != nil {
		t.Fatalf("revoked input owner must remain structurally valid: %v", err)
	}
	encoded, err := json.Marshal(owner)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(firstPath, encoded, 0600); err != nil { t.Fatal(err) }
	for i := 0; i+1 < len(base); i++ {
		if base[i] == "-expect-snapshot-sha256" {
			base[i+1] = fmt.Sprintf("%x", sha256.Sum256(encoded))
			break
		}
	}

	var filtered wazuhPairPageReceipt
	if err := json.Unmarshal(run(
		"-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", "5710", "-pair-agent-id", "001",
	), &filtered); err != nil { t.Fatal(err) }
	if len(filtered.Pairs) != 1 || filtered.Pairs[0].FirstEntryID == revokedID {
		t.Fatalf("revoked source still in authorized rule/agent cohort: %+v", filtered)
	}
	var reverse reverseWazuhPairReceipt
	if err := json.Unmarshal(run("-find-saved-pair-second-entry", secondID), &reverse); err != nil {
		t.Fatal(err)
	}
	if len(reverse.Pairs) != 1 || reverse.Pairs[0].FirstEntryID == revokedID ||
		reverse.Pairs[0].FirstEntryID != filtered.Pairs[0].FirstEntryID {
		t.Fatalf("revoked pair rediscovered via second-Case reverse Entry selector: %+v", reverse)
	}
	// A previously obtained pair ID cannot reopen evidence following
	// revocation even when the operator separately rehashed the owner file.
	destination := filepath.Join(t.TempDir(), "forbidden-revoked-pair.json")
	var stdout bytes.Buffer
	err = runDemoCLI(context.Background(), append(append([]string{}, base...),
		"-inspect-saved-pair-first-entry", revokedID,
		"-inspect-saved-pair-second-entry", secondID,
		"-private-report", destination), &stdout)
	if err == nil || stdout.Len() != 0 {
		t.Fatalf("revoked exact Entry ID reopened private evidence: %v %q", err, stdout.String())
	}
	if _, statErr := os.Lstat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("revoked Entry left private dossier on disk: %v", statErr)
	}
}
