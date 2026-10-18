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

// A second Case being authorized never compensates for a REVOKED source in
// the first Case. In particular, the correlated pair ordinal is a selector,
// NOT a grant to read a source whose ACL was changed after an older receipt.
// Test the whole operator-controlled snapshot -> new digest -> CLI -> private
// file pipeline, not merely an ACL helper or invented policy.
func TestIntegratedPrivatePairRejectsRevokedFirstSourceWithValidBothDigests(t *testing.T) {
	firstPath, firstCase, _, secondPath, secondCase, secondDigest, pair :=
		makeSavedPairDrilldownOwners(t)
	first, err := loadDemoSnapshot(firstPath)
	if err != nil { t.Fatal(err) }
	if len(first.Principals) != 1 { t.Fatal("expected first saved owner scoped to one principal") }
	owner := first.Principals[0]
	const otherReader = "independently-authorized-but-not-first-source-reader"
	first.Principals = append(first.Principals, timeline.MemoryTimelinePrincipal{
		ID: otherReader, Tenant: owner.Tenant, Environment: owner.Environment,
	})
	var sourceRef string
	for _, event := range first.Entries {
		if event.CaseID == firstCase && event.Value.ID == pair.FirstEntryID {
			sourceRef = event.Value.SourceRef
			break
		}
	}
	if sourceRef == "" { t.Fatal("correlated first event absent in first saved Case") }
	changed := false
	for i := range first.Sources {
		if first.Sources[i].Identity.CaseID == firstCase &&
			first.Sources[i].Identity.Reference == sourceRef {
			first.Sources[i].Readers = []string{otherReader}
			changed = true
		}
	}
	if !changed { t.Fatal("typed source for correlated first event absent") }
	// The fixture remains an internally valid multi-reader owner, and the
	// operator records the CORRECT digest for its updated snapshot.
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), "analyst-a", first); err != nil {
		t.Fatalf("source-revoked owner must remain structurally valid: %v", err)
	}
	modified, err := json.Marshal(first)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(firstPath, modified, 0600); err != nil { t.Fatal(err) }
	firstDigest := fmt.Sprintf("%x", sha256.Sum256(modified))
	// This independent second owner remains fully readable. It cannot
	// authorize a source belonging to the first Case.
	second, err := loadDemoSnapshot(secondPath)
	if err != nil { t.Fatal(err) }
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), "analyst-a", second); err != nil {
		t.Fatalf("positive second owner fixture invalid: %v", err)
	}
	destination := filepath.Join(t.TempDir(), "must-not-expose-first-source.json")
	args := []string{
		"-principal", "analyst-a",
		"-snapshot", firstPath, "-compare-saved-case", firstCase,
		"-expect-snapshot-sha256", firstDigest,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondDigest,
		"-inspect-saved-pair", "1",
		"-private-report", destination,
	}
	var stdout bytes.Buffer
	if err := runDemoCLI(context.Background(), args, &stdout); err == nil {
		t.Fatal("private pair source obtained first Case access from authorized second Case")
	}
	if stdout.Len() != 0 { t.Fatalf("denied paired inspection wrote stdout: %q", stdout.String()) }
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		t.Fatalf("first source denial left a private artifact: %v", err)
	}
	// The explicit resume mode must also reject the private-report flag;
	// it may never silently start a listener or emit credentials instead.
	unauthorizedReport := filepath.Join(t.TempDir(), "invalid-resume-output.json")
	stdout.Reset()
	if err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a",
		"-snapshot", secondPath, "-resume-import-search-case", secondCase,
		"-expect-snapshot-sha256", secondDigest,
		"-private-report", unauthorizedReport,
	}, &stdout); err == nil || stdout.Len() != 0 {
		t.Fatalf("resume/private report ambiguity started a session or leaked token: %q", stdout.String())
	}
	if _, err := os.Lstat(unauthorizedReport); !os.IsNotExist(err) {
		t.Fatal("invalid resume mode wrote an unauthorized file")
	}
}
