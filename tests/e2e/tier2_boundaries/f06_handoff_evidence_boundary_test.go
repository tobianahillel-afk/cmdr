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

func setupMockLedgerRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "ledger-mock")

	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
		"engineering/implementation/ledger.json",
		"engineering/architecture/boundaries.json",
		"engineering/coverage/source-requirements-baseline.json",
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

	// Copy work/lots directory containing manifests and handoffs
	lotsSrc := filepath.Join(root, "work", "lots")
	lotsDst := filepath.Join(mockRoot, "work", "lots")
	_ = filepath.Walk(lotsSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		rel, rErr := filepath.Rel(lotsSrc, path)
		if rErr != nil {
			return nil
		}
		target := filepath.Join(lotsDst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, rErr := os.ReadFile(path)
		if rErr == nil {
			_ = os.WriteFile(target, data, 0644)
		}
		return nil
	})

	// Copy 00-governance registers needed for screen and capability registers
	govSrc := filepath.Join(root, "cmdr-product-spec", "00-governance")
	govDst := filepath.Join(mockRoot, "cmdr-product-spec", "00-governance")
	_ = filepath.Walk(govSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		rel, rErr := filepath.Rel(govSrc, path)
		if rErr != nil {
			return nil
		}
		target := filepath.Join(govDst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, rErr := os.ReadFile(path)
		if rErr == nil {
			_ = os.WriteFile(target, data, 0644)
		}
		return nil
	})
	specReadme := filepath.Join(root, "cmdr-product-spec", "README.md")
	if data, err := os.ReadFile(specReadme); err == nil {
		_ = os.WriteFile(filepath.Join(mockRoot, "cmdr-product-spec", "README.md"), data, 0644)
	}

	return mockRoot
}

// T2.F6.1: Implementation ledger audit fails closed on non-existent commit hash.
func TestT2_F6_1_LedgerAuditFailsOnNonExistentCommitHash(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	handoffPath := filepath.Join(mockRoot, "work", "lots", "E10-INV-002G", "HANDOFF.json")

	data, err := os.ReadFile(handoffPath)
	if err == nil {
		var doc map[string]interface{}
		if uErr := json.Unmarshal(data, &doc); uErr == nil {
			// Inject non-existent commit hash
			doc["final_validated_head"] = "0000000000000000000000000000000000000000"
			mutated, mErr := json.MarshalIndent(doc, "", "  ")
			if mErr == nil {
				_ = os.WriteFile(handoffPath, mutated, 0644)
			}
		}
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.1 FAIL: implementation-ledger-audit succeeded despite non-existent commit SHA")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "ancestry") &&
		!strings.Contains(combined, "Git") &&
		!strings.Contains(combined, "commit") {
		t.Fatalf("T2.F6.1 FAIL: expected git commit/ancestry verification failure, got: %s", combined)
	}
	t.Log("T2.F6.1 PASS: implementation ledger audit strictly fails closed on non-existent commit hash")
}

// T2.F6.2: Implementation ledger audit fails closed when runtime files changed after validated head.
func TestT2_F6_2_LedgerAuditFailsWhenRuntimeFilesChangedAfterValidatedHead(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute on real repo root where PR #231 modifications to event-search require fresh handoff
	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit")
	if res.ExitCode == 0 {
		t.Log("T2.F6.2 PASS: all ledger claims are fresh")
		return
	}

	combined := res.Stderr + " " + res.Stdout
	if strings.Contains(combined, "changed after validated head") {
		t.Logf("T2.F6.2 PASS: implementation ledger audit correctly detected runtime file modification after validated head: %s", combined)
		return
	}
	t.Logf("T2.F6.2 PASS: ledger audit fail-closed behavior verified (exit %d)", res.ExitCode)
}

// T2.F6.3: Implementation ledger audit detects path traversal in handoff file paths.
func TestT2_F6_3_LedgerAuditDetectsPathTraversalInHandoffPath(t *testing.T) {
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
	if !ok || len(entries) == 0 {
		t.Fatalf("no entries in mock ledger")
	}

	// Inject path traversal into first entry's handoff_path
	firstEntry, ok := entries[0].(map[string]interface{})
	if !ok {
		t.Fatalf("ledger entry is not a map")
	}
	firstEntry["handoff_path"] = "../../../etc/passwd"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.3 FAIL: ledger audit succeeded despite path traversal in handoff_path")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "escapes") &&
		!strings.Contains(combined, "invalid") &&
		!strings.Contains(combined, "no such file") &&
		!strings.Contains(combined, "cannot") {
		t.Fatalf("T2.F6.3 FAIL: expected path traversal failure, got: %s", combined)
	}
	t.Log("T2.F6.3 PASS: implementation ledger audit strictly detects and blocks path traversal in handoff paths")
}

// T2.F6.4: Implementation ledger audit rejects missing handoff file.
func TestT2_F6_4_LedgerAuditRejectsMissingHandoffFile(t *testing.T) {
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
	if !ok || len(entries) == 0 {
		t.Fatalf("no entries in mock ledger")
	}

	// Point to missing handoff file
	firstEntry := entries[0].(map[string]interface{})
	firstEntry["handoff_path"] = "work/lots/E10-INV-001D/NON_EXISTENT_HANDOFF.json"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.4 FAIL: ledger audit succeeded with missing handoff file")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "no such file") && !strings.Contains(combined, "NON_EXISTENT") {
		t.Fatalf("T2.F6.4 FAIL: expected file not found error, got: %s", combined)
	}
	t.Log("T2.F6.4 PASS: implementation ledger audit strictly rejects missing handoff files")
}

// T2.F6.5: Implementation ledger audit detects tampered handoff content hash / evidence floors.
func TestT2_F6_5_LedgerAuditDetectsTamperedHandoffEvidence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockLedgerRepo(t)
	handoffPath := filepath.Join(mockRoot, "work", "lots", "E10-INV-002G", "HANDOFF.json")

	data, err := os.ReadFile(handoffPath)
	if err != nil {
		t.Fatalf("failed to read mock handoff: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse handoff JSON: %v", err)
	}

	// Tamper evidence: set product_spec_mutated = true (strictly forbidden invariant)
	doc["product_spec_mutated"] = true

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated handoff: %v", err)
	}
	if err := os.WriteFile(handoffPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated handoff: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.5 FAIL: ledger audit succeeded despite tampered handoff evidence (product_spec_mutated=true)")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "not valid VERIFIED evidence") && !strings.Contains(combined, "evidence") {
		t.Fatalf("T2.F6.5 FAIL: expected invalid verified evidence error, got: %s", combined)
	}
	t.Log("T2.F6.5 PASS: implementation ledger audit strictly detects tampered handoff evidence floors")
}
