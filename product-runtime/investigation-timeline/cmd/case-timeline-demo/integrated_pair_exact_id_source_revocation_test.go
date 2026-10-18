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
	if _,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",
		secondPath,secondCase,secondHash);err!=nil {
		t.Fatalf("independently authorized second Case unexpectedly unavailable: %v",err)
	}
}
