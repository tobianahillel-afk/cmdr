package tier2_boundaries

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T2.F10.1: Ledger audit fails closed if CAP-STD-003 and CAP-STD-004 point to same unit.
func TestT2_F10_1_LedgerAuditFailsIfBothCapabilitiesPointToSameUnit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")

	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse ledger JSON: %v", err)
	}

	entries, ok := doc["entries"].([]interface{})
	if !ok || len(entries) < 2 {
		t.Fatalf("unexpected entries structure in ledger")
	}

	// Force second entry to share the same work unit as first entry
	firstEntry := entries[0].(map[string]interface{})
	secondEntry := entries[1].(map[string]interface{})
	targetUnit := firstEntry["work_unit"].(string)
	secondEntry["work_unit"] = targetUnit

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F10.1 FAIL: ledger audit passed when two capabilities share the same work unit")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "duplicate implementation work unit") && !strings.Contains(combined, targetUnit) {
		t.Fatalf("T2.F10.1 FAIL: expected 'duplicate implementation work unit' error, got: %s", combined)
	}
	t.Log("T2.F10.1 PASS: ledger audit strictly fails closed when two capabilities share a single work unit")
}

// T2.F10.2: Ledger audit fails closed if CAP-STD-004 points to non-existent unit.
func TestT2_F10_2_LedgerAuditFailsIfCapabilityPointsToNonExistentUnit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")

	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse ledger JSON: %v", err)
	}

	entries := doc["entries"].([]interface{})
	firstEntry := entries[0].(map[string]interface{})
	firstEntry["work_unit"] = "E10-NON-EXISTENT-UNIT-999"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F10.2 FAIL: ledger audit passed with non-existent work unit")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "not VERIFIED or ACCEPTED") && !strings.Contains(combined, "E10-NON-EXISTENT-UNIT-999") {
		t.Fatalf("T2.F10.2 FAIL: expected non-existent work unit error, got: %s", combined)
	}
	t.Log("T2.F10.2 PASS: ledger audit strictly fails closed when pointing to non-existent work unit")
}

// T2.F10.3: Ledger audit fails if total capability count is decreased.
func TestT2_F10_3_LedgerAuditFailsIfCapabilityCountDecreased(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")

	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse ledger JSON: %v", err)
	}

	// Empty entries array (0 entries)
	doc["entries"] = []interface{}{}

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F10.3 FAIL: ledger audit passed with 0 entries")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "must contain at least one verified claim") {
		t.Fatalf("T2.F10.3 FAIL: expected empty ledger error, got: %s", combined)
	}
	t.Log("T2.F10.3 PASS: ledger audit strictly enforces non-decreasing capability count floor")
}

// T2.F10.4: Ledger audit rejects malformed capability ID format.
func TestT2_F10_4_LedgerAuditRejectsMalformedCapabilityIDFormat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")

	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse ledger JSON: %v", err)
	}

	entries := doc["entries"].([]interface{})
	firstEntry := entries[0].(map[string]interface{})
	firstEntry["capability"] = "INVALID-CAPABILITY-FORMAT-NO-NUMBERS"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F10.4 FAIL: ledger audit passed with malformed capability ID")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "has invalid capability") {
		t.Fatalf("T2.F10.4 FAIL: expected 'has invalid capability' error, got: %s", combined)
	}
	t.Log("T2.F10.4 PASS: ledger audit strictly rejects malformed capability ID format")
}

// T2.F10.5: Ledger audit rejects duplicate capability registration.
func TestT2_F10_5_LedgerAuditRejectsDuplicateCapabilityRegistration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")

	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse ledger JSON: %v", err)
	}

	entries := doc["entries"].([]interface{})
	// Duplicate the first capability ID onto the second entry
	firstEntry := entries[0].(map[string]interface{})
	secondEntry := entries[1].(map[string]interface{})
	secondEntry["capability"] = firstEntry["capability"]

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F10.5 FAIL: ledger audit passed with duplicate capability registration")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "duplicate implementation capability") {
		t.Fatalf("T2.F10.5 FAIL: expected 'duplicate implementation capability' error, got: %s", combined)
	}
	t.Log("T2.F10.5 PASS: ledger audit strictly rejects duplicate capability registration")
}
