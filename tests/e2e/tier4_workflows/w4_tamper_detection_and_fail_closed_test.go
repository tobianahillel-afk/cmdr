package tier4_workflows

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

func setupW4MockRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "w4-tamper")

	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
		"engineering/architecture/boundaries.json",
		"engineering/security/security-gates.json",
		"engineering/security/runtime-security-policy.json",
		"engineering/testing/runtime-security-evidence.json",
		"engineering/recovery/lease-policy.json",
		"engineering/recovery/claims.json",
		"engineering/recovery/journal.json",
		"engineering/implementation/ledger.json",
		"engineering/decisions/decision-registry.json",
		"engineering/performance/performance-registry.json",
		"cmdr-product-spec/README.md",
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

	for _, mod := range harness.RuntimeGoModules {
		dir := filepath.Join(mockRoot, filepath.FromSlash(mod))
		_ = os.MkdirAll(dir, 0755)
	}

	return mockRoot
}

// TestT4_04_TamperDetectionAndFailClosed:
// Adversarial workflow injecting four active tamper vectors into an isolated target repository:
// 1. Invalid HMAC/digest in recovery journal
// 2. Unregistered/invalid work unit entry in implementation ledger
// 3. Untracked file in repository root
// 4. Altered product spec invariant
// Asserts that cmdr-dev mainline-health-gate evaluates the target and strictly fails closed (ExitCode != 0).
func TestT4_04_TamperDetectionAndFailClosed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	mockRoot := setupW4MockRepo(t)

	// Vector 1: Recovery Journal with Corrupted HMAC Digest
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")
	journalData, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("failed to read mock journal: %v", err)
	}
	var journalDoc map[string]interface{}
	if err := json.Unmarshal(journalData, &journalDoc); err != nil {
		t.Fatalf("failed to parse journal: %v", err)
	}
	if events, ok := journalDoc["events"].([]interface{}); ok && len(events) > 0 {
		if firstEvt, ok := events[0].(map[string]interface{}); ok {
			firstEvt["digest"] = "deadbeefbadc0ffee000000000000000000000000000000000000000000000000"
		}
	}
	corruptedJournal, _ := json.MarshalIndent(journalDoc, "", "  ")
	if err := os.WriteFile(journalPath, corruptedJournal, 0644); err != nil {
		t.Fatalf("failed to write corrupted journal: %v", err)
	}

	// Vector 2: Stale / Non-Existent Work Unit in Implementation Ledger
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")
	ledgerData, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger: %v", err)
	}
	var ledgerDoc map[string]interface{}
	if err := json.Unmarshal(ledgerData, &ledgerDoc); err != nil {
		t.Fatalf("failed to parse ledger: %v", err)
	}
	if entries, ok := ledgerDoc["entries"].([]interface{}); ok {
		tamperedEntry := map[string]interface{}{
			"capability":       "CAP-TAMPER-001",
			"work_unit":        "E10-PHANTOM-UNIT-NONEXISTENT",
			"manifest_path":    "work/lots/E10-PHANTOM/manifest.json",
			"handoff_path":     "work/lots/E10-PHANTOM/HANDOFF.json",
			"runtime_boundary": "event-search-runtime",
			"runtime_root":     "product-runtime/event-search/**",
		}
		ledgerDoc["entries"] = append(entries, tamperedEntry)
	}
	corruptedLedger, _ := json.MarshalIndent(ledgerDoc, "", "  ")
	if err := os.WriteFile(ledgerPath, corruptedLedger, 0644); err != nil {
		t.Fatalf("failed to write corrupted ledger: %v", err)
	}

	// Vector 3: Untracked Sentinel File in Repository Root
	sentinelPath := filepath.Join(mockRoot, ".tamper_sentinel.tmp")
	if err := os.WriteFile(sentinelPath, []byte("adversarial-tamper-probe"), 0644); err != nil {
		t.Fatalf("failed to create sentinel file: %v", err)
	}

	// Vector 4: Product Spec Invariant Alteration
	specReadmePath := filepath.Join(mockRoot, "cmdr-product-spec", "README.md")
	if err := os.WriteFile(specReadmePath, []byte("# TAMPERED PRODUCT SPEC INVARIANT\n"), 0644); err != nil {
		t.Fatalf("failed to alter product spec invariant: %v", err)
	}

	// Execute mainline-health-gate against the tampered mock repository
	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")

	// Strict Fail-Closed Invariant: Exit code 0 is an absolute failure
	if gateRes.ExitCode == 0 {
		t.Fatalf("T4.04 FAIL: mainline-health-gate returned exit code 0 on tampered repository; expected non-zero fail-closed exit code")
	}

	if gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
		t.Fatalf("T4.04 FAIL: unexpected exit code %d; expected 1 (MAIN_UNHEALTHY) or 2 (MAIN_UNVERIFIED)", gateRes.ExitCode)
	}

	var report struct {
		Status     string `json:"status"`
		ExitCode   int    `json:"exit_code"`
		Dimensions map[string]struct {
			Name   string   `json:"name"`
			Passed bool     `json:"passed"`
			Status string   `json:"status"`
			Errors []string `json:"errors"`
		} `json:"dimensions"`
	}

	if err := harness.ParseJSON(gateRes.Stdout, &report); err != nil {
		t.Fatalf("failed to parse mainline-health-gate JSON report: %v\nStdout: %s\nStderr: %s", err, gateRes.Stdout, gateRes.Stderr)
	}

	if report.Status == "MAIN_HEALTHY" {
		t.Fatalf("T4.04 FAIL: health gate status reported MAIN_HEALTHY on tampered repository")
	}

	hasFailedDimension := false
	for dimName, dim := range report.Dimensions {
		if !dim.Passed || dim.Status == "FAIL" || dim.Status == "UNVERIFIED" {
			hasFailedDimension = true
			t.Logf("Tamper vector detected in dimension %s: status=%s, errors=%v", dimName, dim.Status, dim.Errors)
		}
	}

	if !hasFailedDimension {
		t.Fatalf("T4.04 FAIL: health gate failed closed with exit %d but no dimensions reported failure", gateRes.ExitCode)
	}

	t.Logf("T4.04 PASS: Tamper detection and fail-closed invariants strictly enforced (gate status=%s, exit=%d)",
		report.Status, gateRes.ExitCode)
}
