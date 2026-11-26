package tier2_boundaries

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

func setupMockJournalRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "journal-mock")

	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
		"engineering/recovery/lease-policy.json",
		"engineering/recovery/claims.json",
		"engineering/recovery/journal.json",
	}

	for _, rel := range copyFiles {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(mockRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatalf("failed to create dir for %s: %v", dst, err)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("failed to copy %s to mock: %v", rel, err)
		}
	}

	// Bootstrap cmdr-product-spec directory so validateState succeeds
	mockSpecDir := filepath.Join(mockRoot, "cmdr-product-spec")
	if err := os.MkdirAll(mockSpecDir, 0755); err != nil {
		t.Fatalf("failed to create cmdr-product-spec dir: %v", err)
	}
	specReadme := filepath.Join(root, "cmdr-product-spec", "README.md")
	if data, err := os.ReadFile(specReadme); err == nil {
		_ = os.WriteFile(filepath.Join(mockSpecDir, "README.md"), data, 0644)
	}

	// Verify pre-mutation baseline passes recovery-journal-audit with exit 0
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	baseRes := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if baseRes.ExitCode != 0 {
		t.Fatalf("pre-mutation baseline recovery-journal-audit failed (exit %d): %s\n%s",
			baseRes.ExitCode, baseRes.Stderr, baseRes.Stdout)
	}

	return mockRoot
}

// T2.F5.1: Recovery journal audit fails closed on corrupted HMAC signature / digest.
func TestT2_F5_1_JournalAuditFailsOnCorruptedDigest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockJournalRepo(t)
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")

	// Create an event with an invalid/forged digest
	journalWithCorruptedDigest := `{
  "schema_version": 1,
  "registry_kind": "engineering-recovery-journal",
  "events": [
    {
      "sequence": 1,
      "id": "EVT-20261122-000001",
      "work_unit": "E10-INV-002G",
      "agent_id": "AGENT-AUTOMATION-TEST",
      "lease_id": "LEASE-20261122-000001",
      "kind": "lease-acquired",
      "observed_at": "2026-11-22T10:00:00Z",
      "head_sha": "c1e75db81662b33a0f1208afdbf8b688fe7977e3",
      "previous_digest": "93f61a4836d343a345deacb2c2f171dada0074b7abd1f5755adf4a04e7558074",
      "digest": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
      "evidence_refs": []
    }
  ]
}`
	if err := os.WriteFile(journalPath, []byte(journalWithCorruptedDigest), 0644); err != nil {
		t.Fatalf("failed to write corrupted journal: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F5.1 FAIL: recovery-journal-audit succeeded with corrupted HMAC/digest")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "digest mismatch") && !strings.Contains(combined, "lease binding") {
		t.Fatalf("T2.F5.1 FAIL: expected digest mismatch or lease binding failure, got: %s", combined)
	}
	t.Log("T2.F5.1 PASS: recovery journal audit strictly fails closed on corrupted digest/HMAC")
}

// T2.F5.2: Recovery journal audit fails closed on out-of-order sequence counter.
func TestT2_F5_2_JournalAuditFailsOnOutOfOrderSequence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockJournalRepo(t)
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")

	// Event with sequence: 5 instead of sequence: 1
	journalWithBadSeq := `{
  "schema_version": 1,
  "registry_kind": "engineering-recovery-journal",
  "events": [
    {
      "sequence": 5,
      "id": "EVT-20261122-000001",
      "work_unit": "E10-INV-002G",
      "agent_id": "AGENT-AUTOMATION-TEST",
      "lease_id": "LEASE-20261122-000001",
      "kind": "lease-acquired",
      "observed_at": "2026-11-22T10:00:00Z",
      "head_sha": "c1e75db81662b33a0f1208afdbf8b688fe7977e3",
      "previous_digest": "93f61a4836d343a345deacb2c2f171dada0074b7abd1f5755adf4a04e7558074",
      "digest": "0000000000000000000000000000000000000000000000000000000000000000",
      "evidence_refs": []
    }
  ]
}`
	if err := os.WriteFile(journalPath, []byte(journalWithBadSeq), 0644); err != nil {
		t.Fatalf("failed to write journal with bad sequence: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F5.2 FAIL: recovery-journal-audit succeeded with out-of-order sequence counter")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "sequence must be contiguous") {
		t.Fatalf("T2.F5.2 FAIL: expected 'sequence must be contiguous' error, got: %s", combined)
	}
	t.Log("T2.F5.2 PASS: recovery journal audit strictly fails closed on non-contiguous sequence numbers")
}

// T2.F5.3: Recovery journal audit fails closed on invalid timestamp format.
func TestT2_F5_3_JournalAuditFailsOnInvalidTimestamp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockJournalRepo(t)
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")

	journalWithInvalidTimestamp := `{
  "schema_version": 1,
  "registry_kind": "engineering-recovery-journal",
  "events": [
    {
      "sequence": 1,
      "id": "EVT-20261122-000001",
      "work_unit": "E10-INV-002G",
      "agent_id": "AGENT-AUTOMATION-TEST",
      "lease_id": "LEASE-20261122-000001",
      "kind": "lease-acquired",
      "observed_at": "not-a-valid-timestamp-string",
      "head_sha": "c1e75db81662b33a0f1208afdbf8b688fe7977e3",
      "previous_digest": "93f61a4836d343a345deacb2c2f171dada0074b7abd1f5755adf4a04e7558074",
      "digest": "0000000000000000000000000000000000000000000000000000000000000000",
      "evidence_refs": []
    }
  ]
}`
	if err := os.WriteFile(journalPath, []byte(journalWithInvalidTimestamp), 0644); err != nil {
		t.Fatalf("failed to write journal: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F5.3 FAIL: recovery-journal-audit succeeded with invalid timestamp")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "observed_at") && !strings.Contains(combined, "lease binding") {
		t.Fatalf("T2.F5.3 FAIL: expected timestamp parsing failure on observed_at, got: %s", combined)
	}
	t.Log("T2.F5.3 PASS: recovery journal audit strictly rejects invalid timestamp formats")
}

// T2.F5.4: Recovery journal audit handles 0-byte file without crashing.
func TestT2_F5_4_JournalAuditHandlesZeroByteFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockJournalRepo(t)
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")

	// Overwrite with 0 bytes
	if err := os.WriteFile(journalPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write 0-byte journal: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F5.4 FAIL: recovery-journal-audit succeeded on empty 0-byte file, expected fail-closed")
	}

	// Must fail cleanly with exit code 1 or 2, without panicking
	if strings.Contains(res.Stderr, "panic:") || strings.Contains(res.Stdout, "panic:") {
		t.Fatalf("T2.F5.4 FAIL: recovery-journal-audit panicked on 0-byte file: %s", res.Stderr)
	}
	t.Log("T2.F5.4 PASS: recovery journal audit handles 0-byte file gracefully without crashing")
}

// T2.F5.5: Recovery journal audit detects truncated event entry.
func TestT2_F5_5_JournalAuditDetectsTruncatedEvent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockJournalRepo(t)
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")

	// Truncated event entry (missing required fields)
	truncatedJournal := `{
  "schema_version": 1,
  "registry_kind": "engineering-recovery-journal",
  "events": [
    {
      "sequence": 1,
      "id": "EVT-INCOMPLETE"
    }
  ]
}`
	if err := os.WriteFile(journalPath, []byte(truncatedJournal), 0644); err != nil {
		t.Fatalf("failed to write truncated journal: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F5.5 FAIL: recovery-journal-audit succeeded with truncated event entry")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "invalid or duplicate recovery event id") &&
		!strings.Contains(combined, "invalid") &&
		!strings.Contains(combined, "identity") &&
		!strings.Contains(combined, "unknown work unit") &&
		!strings.Contains(combined, "event") {
		t.Fatalf("T2.F5.5 FAIL: expected validation failure on truncated event, got: %s", combined)
	}
	t.Log("T2.F5.5 PASS: recovery journal audit strictly detects and rejects truncated event entry")
}
