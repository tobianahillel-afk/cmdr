package tier2_boundaries

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return os.WriteFile(path, data, perm)
	}
	tmpName := tmpFile.Name()
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	_ = os.Chmod(tmpName, perm)
	_ = os.Remove(path)
	if err := os.Rename(tmpName, path); err != nil {
		return os.WriteFile(path, data, perm)
	}
	return nil
}

func setupMockLedgerRepo(t *testing.T) string {
	t.Helper()
	return harness.CloneIsolatedRepo(t, "ledger-mock")
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
				_ = atomicWriteFile(handoffPath, mutated, 0644)
			}
		}
	}

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.1 FAIL: implementation-ledger-audit succeeded despite non-existent commit SHA")
	}

	combined := strings.ToLower(res.Stderr + " " + res.Stdout)
	if !strings.Contains(combined, "ancestry") &&
		!strings.Contains(combined, "git") &&
		!strings.Contains(combined, "commit") &&
		!strings.Contains(combined, "invalid final validated head") &&
		!strings.Contains(combined, "head") {
		t.Fatalf("T2.F6.1 FAIL: expected git commit/ancestry verification failure, got: %s", combined)
	}
	t.Log("T2.F6.1 PASS: implementation ledger audit strictly fails closed on non-existent commit hash")
}

func runGitMockCmd(t *testing.T, mockRoot string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", mockRoot}, args...)
	out, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed in mock repo %s: %v (%s)", strings.Join(args, " "), mockRoot, err, string(out))
	}
	return strings.TrimSpace(string(out))
}

func setupIsolatedGitLedgerRepo(t *testing.T) (string, string) {
	t.Helper()
	mockRoot := harness.CloneIsolatedRepo(t, "ledger-git-mock")
	validatedHead := runGitMockCmd(t, mockRoot, "rev-parse", "HEAD")
	return mockRoot, validatedHead
}

// T2.F6.2: Implementation ledger audit fails closed when runtime files changed after validated head.
func TestT2_F6_2_LedgerAuditFailsWhenRuntimeFilesChangedAfterValidatedHead(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot, validatedHead := setupIsolatedGitLedgerRepo(t)

	// Baseline sanity check: before any post-validated runtime commit, audit must pass (exit 0)
	baselineRes := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if baselineRes.ExitCode != 0 {
		t.Fatalf("T2.F6.2 FAIL: pre-injection baseline ledger audit failed (exit %d): stdout=%s stderr=%s",
			baselineRes.ExitCode, baselineRes.Stdout, baselineRes.Stderr)
	}

	// Inject fault: commit a runtime file modification AFTER final_validated_head
	mutatedRuntimeFile := filepath.Join(mockRoot, "product-runtime", "event-search", "post_validated_change.go")
	if err := os.WriteFile(mutatedRuntimeFile, []byte("package eventsearch\n\n// Post-validation runtime modification\n"), 0644); err != nil {
		t.Fatalf("failed to write post-validated runtime file: %v", err)
	}
	runGitMockCmd(t, mockRoot, "add", "product-runtime/event-search/post_validated_change.go")
	runGitMockCmd(t, mockRoot, "commit", "-m", "modify event-search runtime after validated head")

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F6.2 FAIL: implementation-ledger-audit succeeded despite runtime file change after validated head %s", validatedHead)
	}
	if res.ExitCode != 1 {
		t.Fatalf("T2.F6.2 FAIL: expected exit code 1, got %d (stdout=%q stderr=%q)", res.ExitCode, res.Stdout, res.Stderr)
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "changed after validated head") {
		t.Fatalf("T2.F6.2 FAIL: expected 'changed after validated head' error, got: %s", combined)
	}
	t.Logf("T2.F6.2 PASS: implementation ledger audit failed closed (exit %d) on post-validated-head runtime change: %s",
		res.ExitCode, strings.TrimSpace(combined))
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
	if err := atomicWriteFile(ledgerPath, mutated, 0644); err != nil {
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
	if err := atomicWriteFile(ledgerPath, mutated, 0644); err != nil {
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
	if err := atomicWriteFile(handoffPath, mutated, 0644); err != nil {
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
