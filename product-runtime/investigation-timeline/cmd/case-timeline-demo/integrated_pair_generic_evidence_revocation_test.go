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

// Consolidated end-to-end security boundary: two independently authorized
// saved Cases and sensitive v3 generic previews in one explicitly requested
// private dossier. An imported generic source may be revoked AFTER saving;
// its Entry ID and preview must remain absent even with a correct fresh SHA.
func TestIntegratedOptedInPairOmitsRevokedFirstGenericSource(t *testing.T) {
	firstOriginal, firstCase, _, secondOriginal, secondCase, _, pair := makeSavedPairDrilldownOwners(t)
	firstFeed := writeJSONLEvents(t,
		"{\"id\":\"approved-first\",\"timestamp\":\"2026-10-17T10:00:07Z\",\"message\":\"ALLOW-GENERIC-FIRST-CONTENT\"}\n"+
		"{\"id\":\"denied-first\",\"timestamp\":\"2026-10-17T10:00:08Z\",\"message\":\"DENIED-GENERIC-FIRST-CONTENT\"}\n")
	secondFeed := writeJSONLEvents(t,
		"{\"id\":\"approved-second\",\"timestamp\":\"2026-10-17T10:00:25Z\",\"message\":\"ALLOW-GENERIC-SECOND-CONTENT\"}\n")
	firstPath := filepath.Join(t.TempDir(), "first-private-case.json")
	secondPath := filepath.Join(t.TempDir(), "second-private-case.json")
	for _, p := range []struct{original, caseID, feed, path string}{
		{firstOriginal, firstCase, firstFeed, firstPath},
		{secondOriginal, secondCase, secondFeed, secondPath},
	} {
		if err := writeDemoPersistedNormalizedCase(context.Background(), "analyst-a",
			p.original, p.caseID, []string{p.feed}, p.path, &bytes.Buffer{}); err != nil {
			t.Fatalf("cannot persist an independently owned Case: %v", err)
		}
	}
	first, err := loadDemoSnapshot(firstPath)
	if err != nil { t.Fatal(err) }
	if len(first.Principals) != 1 { t.Fatal("unexpected first Case principal graph") }
	p := first.Principals[0]
	const refusedPrincipal = "generic-source-previous-reader"
	first.Principals = append(first.Principals, timeline.MemoryTimelinePrincipal{
		ID: refusedPrincipal, Tenant: p.Tenant, Environment: p.Environment,
	})
	deniedID := ""
	deniedRef := ""
	for _, e := range first.Entries {
		if e.CaseID != firstCase || e.Value.Provenance != "operator-local-jsonl" {
			continue
		}
		for _, s := range first.Sources {
			if s.Identity.Reference == e.Value.SourceRef &&
				strings.Contains(s.Preview, "DENIED-GENERIC-FIRST-CONTENT") {
				deniedID = e.Value.ID
				deniedRef = e.Value.SourceRef
			}
		}
	}
	if deniedID == "" || deniedRef == "" { t.Fatal("imported denied source absent") }
	for i := range first.Sources {
		if first.Sources[i].Identity.Reference == deniedRef {
			first.Sources[i].Readers = []string{refusedPrincipal}
		}
	}
	// The owner remains valid and the operator retains its NEW correct digest:
	// a valid hash must not override the typed-source read ACL.
	if _, err := newDemoHandlerWithData(strings.Repeat("0",64), "analyst-a", first); err != nil {
		t.Fatalf("revoked source made input owner invalid: %v", err)
	}
	revised, err := json.Marshal(first)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(firstPath, revised, 0600); err != nil { t.Fatal(err) }
	firstHash := fmt.Sprintf("%x", sha256.Sum256(revised))
	secondRaw, err := os.ReadFile(secondPath)
	if err != nil { t.Fatal(err) }
	secondHash := fmt.Sprintf("%x", sha256.Sum256(secondRaw))
	dest := filepath.Join(t.TempDir(), "only-owned-private-evidence.json")
	var stdout bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstHash,
		"-compare-saved-with-snapshot", secondPath, "-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondHash,
		"-inspect-saved-pair", "1", "-include-pair-generic-evidence",
		"-private-report", dest,
	}, &stdout); err != nil {
		t.Fatalf("independent Wazuh correlation should survive irrelevant generic revocation: %v", err)
	}
	if stdout.Len() != 0 { t.Fatalf("private report emitted to console: %q", stdout.String()) }
	raw, err := os.ReadFile(dest)
	if err != nil { t.Fatal(err) }
	var report struct {
		Format string `json:"format"`
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		FirstGenericContext authorizedGenericPairContextSummary `json:"first_generic_context"`
		FirstGenericEvidence []localEntryInspection `json:"first_generic_evidence"`
		SecondGenericEvidence []localEntryInspection `json:"second_generic_evidence"`
	}
	if err := json.Unmarshal(raw, &report); err != nil { t.Fatal(err) }
	if report.Format != "cmdr-private-verified-wazuh-pair-v3" ||
		report.Pair.FirstEntryID != pair.FirstEntryID ||
		report.Pair.SecondEntryID != pair.SecondEntryID ||
		len(report.SecondGenericEvidence) != 1 {
		t.Fatalf("two-Case dossier or independent second source lost: %+v", report)
	}
	foundAllowedFirst := false
	for _, e := range report.FirstGenericEvidence {
		if strings.Contains(e.Preview, "ALLOW-GENERIC-FIRST-CONTENT") {
			foundAllowedFirst = true
		}
		if e.EntryID == deniedID || strings.Contains(e.Preview, "DENIED-GENERIC-FIRST-CONTENT") {
			t.Fatal("revoked first source passed private expansion")
		}
	}
	if !foundAllowedFirst { t.Fatalf("authorized first source missing: %+v", report.FirstGenericEvidence) }
	if !strings.Contains(report.SecondGenericEvidence[0].Preview,"ALLOW-GENERIC-SECOND-CONTENT") {
		t.Fatal("second Case authorized generic source missing")
	}
	if bytes.Contains(raw, []byte("DENIED-GENERIC-FIRST-CONTENT")) ||
		bytes.Contains(raw, []byte(deniedID)) ||
		bytes.Contains(raw, []byte(deniedRef)) ||
		bytes.Contains(raw, []byte(refusedPrincipal)) {
		t.Fatal("private dossier disclosed revoked source metadata or preview")
	}
	for _, item := range report.FirstGenericContext.Entries {
		if item.EntryID == deniedID { t.Fatal("revoked source leaked via metadata-only contextual lead") }
	}
}

// An existing Case's correlated pair is never proof that either source can be
// read later. In particular an owner can revoke the selected Wazuh source
// while both signed-off snapshot byte digests remain internally consistent.
func TestIntegratedOptedInPairFailsBeforeFileWhenSecondWazuhSourceRevoked(t *testing.T) {
	aPath, aCase, aHash, bPath, bCase, _, pair := makeSavedPairDrilldownOwners(t)
	second, err := loadDemoSnapshot(bPath)
	if err != nil { t.Fatal(err) }
	identity := second.Principals[0]
	const denied = "independent-source-revoked-principal"
	second.Principals = append(second.Principals, timeline.MemoryTimelinePrincipal{
		ID:denied, Tenant:identity.Tenant, Environment:identity.Environment,
	})
	ref := ""
	for _, entry := range second.Entries {
		if entry.CaseID == bCase && entry.Value.ID == pair.SecondEntryID {
			ref = entry.Value.SourceRef
		}
	}
	if ref == "" { t.Fatal("selected Wazuh owner ref missing") }
	for i := range second.Sources {
		if second.Sources[i].Identity.Reference == ref {
			second.Sources[i].Readers = []string{denied}
		}
	}
	if _, err := newDemoHandlerWithData(strings.Repeat("0",64), "analyst-a", second); err != nil {
		t.Fatalf("independently valid revoked Case rejected as fixture: %v", err)
	}
	updated, err := json.Marshal(second)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(bPath,updated,0600); err != nil { t.Fatal(err) }
	bHash := fmt.Sprintf("%x",sha256.Sum256(updated))
	destination := filepath.Join(t.TempDir(), "must-not-create.json")
	var stdout bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-inspect-saved-pair","1","-include-pair-generic-evidence",
		"-private-report",destination,
	},&stdout)
	if err == nil || stdout.Len() != 0 {
		t.Fatalf("revoked correlated source authorized private dossier: %v, output %q",err,stdout.String())
	}
	if _, err:=os.Lstat(destination);!os.IsNotExist(err) {
		t.Fatalf("revoked source left partial private file: %v",err)
	}
}
