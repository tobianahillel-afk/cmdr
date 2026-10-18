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

// Full-diff regression: a valid second Case and correctly recomputed digest
// must NOT turn a revoked second-side generic source into private evidence.
// The selected Wazuh pair and one other still-allowed source remain usable.
func TestConsolidatedPairDossierExcludesRevokedSecondGenericSource(t *testing.T) {
	firstInput, firstCase, _, secondInput, secondCase, _, pair := makeSavedPairDrilldownOwners(t)
	firstFeed := writeJSONLEvents(t,
		"{\"id\":\"first-readable\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"FIRST-GENERIC-ALLOWED-CONTENT\"}\n")
	secondFeed := writeJSONLEvents(t,
		"{\"id\":\"second-readable\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"message\":\"SECOND-GENERIC-ALLOWED-CONTENT\"}\n"+
		"{\"id\":\"second-revoked\",\"timestamp\":\"2026-10-17T10:00:21Z\",\"message\":\"SECOND-GENERIC-REVOKED-CONTENT\"}\n")
	aPath := filepath.Join(t.TempDir(), "first-owner.json")
	bPath := filepath.Join(t.TempDir(), "second-owner.json")
	for _, item := range []struct{src, caseID, feed, dest string}{
		{firstInput, firstCase, firstFeed, aPath},
		{secondInput, secondCase, secondFeed, bPath},
	} {
		if err := writeDemoPersistedNormalizedCase(context.Background(), "analyst-a",
			item.src, item.caseID, []string{item.feed}, item.dest, &bytes.Buffer{}); err != nil {
			t.Fatalf("cannot save independently owned Case %s: %v", item.caseID, err)
		}
	}
	aRaw, err := os.ReadFile(aPath)
	if err != nil { t.Fatal(err) }
	aHash := fmt.Sprintf("%x", sha256.Sum256(aRaw))
	second, err := loadDemoSnapshot(bPath)
	if err != nil { t.Fatal(err) }
	if len(second.Principals) != 1 { t.Fatal("second Case must start with one authorized principal") }
	principal := second.Principals[0]
	const otherReader = "revoked-second-source-reader"
	second.Principals = append(second.Principals, timeline.MemoryTimelinePrincipal{
		ID: otherReader, Tenant: principal.Tenant, Environment: principal.Environment,
	})
	deniedEntry, deniedSource := "", ""
	for _, entry := range second.Entries {
		if entry.CaseID != secondCase || entry.Value.Provenance != "operator-local-jsonl" {
			continue
		}
		for _, source := range second.Sources {
			if source.Identity.Reference == entry.Value.SourceRef &&
				strings.Contains(source.Preview, "SECOND-GENERIC-REVOKED-CONTENT") {
				deniedEntry, deniedSource = entry.Value.ID, entry.Value.SourceRef
			}
		}
	}
	if deniedEntry == "" || deniedSource == "" { t.Fatal("revoked second generic source not present") }
	for i := range second.Sources {
		if second.Sources[i].Identity.Reference == deniedSource {
			second.Sources[i].Readers = []string{otherReader}
		}
	}
	// Changing the owner ACL and *also* updating the operator-held hash
	// proves access rights are checked independently of byte integrity.
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), "analyst-a", second); err != nil {
		t.Fatalf("revoked source must leave other second Case sources valid: %v", err)
	}
	updated, err := json.Marshal(second)
	if err != nil { t.Fatal(err) }
	if err := os.WriteFile(bPath, updated, 0600); err != nil { t.Fatal(err) }
	bHash := fmt.Sprintf("%x", sha256.Sum256(updated))
	destination := filepath.Join(t.TempDir(), "scoped-pair-evidence.json")
	var stdout bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", aPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", aHash,
		"-compare-saved-with-snapshot", bPath, "-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", bHash, "-inspect-saved-pair", "1",
		"-include-pair-generic-evidence", "-private-report", destination,
	}, &stdout); err != nil {
		t.Fatalf("private dossier should survive unrelated second source revocation: %v", err)
	}
	if stdout.Len() != 0 { t.Fatalf("sensitive dossier was written to console: %q", stdout.String()) }
	raw, err := os.ReadFile(destination)
	if err != nil { t.Fatal(err) }
	var report struct {
		Format string `json:"format"`
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		SecondGenericContext authorizedGenericPairContextSummary `json:"second_generic_context"`
		FirstGenericEvidence []localEntryInspection `json:"first_generic_evidence"`
		SecondGenericEvidence []localEntryInspection `json:"second_generic_evidence"`
	}
	if err := json.Unmarshal(raw, &report); err != nil { t.Fatal(err) }
	if report.Format != "cmdr-private-verified-wazuh-pair-v3" ||
		report.Pair.FirstEntryID != pair.FirstEntryID ||
		report.Pair.SecondEntryID != pair.SecondEntryID ||
		len(report.FirstGenericEvidence) != 1 {
		t.Fatalf("revocation corrupted authorized first source or Wazuh pair: %+v", report)
	}
	allowed := false
	for _, opened := range report.SecondGenericEvidence {
		if opened.CaseID != secondCase || !opened.Synthetic || opened.Kind != "telemetry-event" {
			t.Fatalf("generic evidence escaped second owner: %+v", opened)
		}
		if strings.Contains(opened.Preview, "SECOND-GENERIC-ALLOWED-CONTENT") { allowed = true }
		if opened.EntryID == deniedEntry || strings.Contains(opened.Preview, "SECOND-GENERIC-REVOKED-CONTENT") {
			t.Fatal("denied second Case evidence reopened from a valid digest")
		}
	}
	if !allowed { t.Fatal("valid second generic evidence was hidden by unrelated revocation") }
	for _, lead := range report.SecondGenericContext.Entries {
		if lead.EntryID == deniedEntry { t.Fatal("denied second entry leaked even as metadata") }
	}
	for _, denied := range []string{
		"SECOND-GENERIC-REVOKED-CONTENT", deniedEntry, deniedSource, otherReader,
	} {
		if bytes.Contains(raw, []byte(denied)) { t.Fatalf("private dossier disclosed denied second owner data: %q", denied) }
	}
}
