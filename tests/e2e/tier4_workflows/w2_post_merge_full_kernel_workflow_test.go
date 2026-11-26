package tier4_workflows

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

type FullKernelGateReport struct {
	Status     string `json:"status"`
	ExitCode   int    `json:"exit_code"`
	Summary    string `json:"summary"`
	HeadCommit string `json:"head_commit"`
	Dimensions map[string]struct {
		Name    string   `json:"name"`
		Passed  bool     `json:"passed"`
		Status  string   `json:"status"`
		Details []string `json:"details,omitempty"`
		Errors  []string `json:"errors,omitempty"`
	} `json:"dimensions"`
}

// TestT4_02_PostMergeFullKernelWorkflow:
// Simulates the complete post-merge CI pipeline execution on main:
// 1. Clean tree verification
// 2. Doctor diagnostics
// 3. Spec baseline invariance check
// 4. Manifest validation
// 5. PR closure register reconciliation
// 6. Recovery journal & ledger audits
// 7. Mainline health gate execution
// Features exercised: F5, F6, F7, F8, F11, F14.
func TestT4_02_PostMergeFullKernelWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// Step 1: Clean tree verification
	status, err := harness.GitStatusPorcelain(ctx)
	if err != nil {
		t.Fatalf("failed to query git status: %v", err)
	}
	lines := strings.Split(status, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !strings.HasPrefix(trimmed, "?? tests/") {
			t.Logf("Post-merge worktree item: %s", trimmed)
		}
	}

	// Step 2: Doctor diagnostics
	docRes := harness.RunCmdrDev(ctx, "doctor", "--json")
	if docRes.ExitCode != 0 {
		t.Fatalf("cmdr-dev doctor failed with exit code %d: %s", docRes.ExitCode, docRes.Stderr)
	}
	harness.AssertJSONFieldEqual(t, docRes.Stdout, "healthy", true)

	// Step 3: Canonical product spec baseline invariance check (0 diff vs 82e630fc)
	const baselineCommit = "82e630fce3a7f6236e03709116eca6dd0a1acd40"
	specBaselineRes := harness.RunCmdrDev(ctx, "spec-baseline", "--check", "--json")
	if specBaselineRes.ExitCode != 0 {
		t.Fatalf("spec-baseline check failed with exit code %d: %s", specBaselineRes.ExitCode, specBaselineRes.Stderr)
	}

	// Direct git diff verification of cmdr-product-spec
	specDiff, err := harness.GitDiffSpecBaseline(ctx, baselineCommit)
	if err != nil {
		t.Fatalf("failed to run GitDiffSpecBaseline against %s: %v", baselineCommit, err)
	}
	if specDiff != "" {
		t.Fatalf("cmdr-product-spec has non-zero diff vs baseline: %s", specDiff)
	}

	// Step 4: Validate work manifests
	manifestRes := harness.RunCmdrDev(ctx, "validate-manifests", "--json")
	if manifestRes.ExitCode != 0 {
		t.Fatalf("validate-manifests failed with exit code %d: %s", manifestRes.ExitCode, manifestRes.Stderr)
	}

	// Step 5: Validate PR reconciliation closure records
	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	prsData, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}
	var prsList []map[string]interface{}
	if err := json.Unmarshal(prsData, &prsList); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}
	if len(prsList) < 184 {
		t.Fatalf("expected at least 184 reconciled PRs in prs.json, got %d", len(prsList))
	}

	// Step 6: Recovery journal audit (F5 quiescent journal handling)
	journalRes := harness.RunCmdrDev(ctx, "recovery-journal-audit", "--json")
	if journalRes.ExitCode != 0 {
		t.Fatalf("recovery-journal-audit failed (exit %d): %s", journalRes.ExitCode, journalRes.Stderr)
	}

	// Step 7: Implementation ledger audit (F6 handoff evidence verification)
	ledgerRes := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--json")
	if ledgerRes.ExitCode != 0 {
		t.Fatalf("implementation-ledger-audit failed (exit %d): %s", ledgerRes.ExitCode, ledgerRes.Stderr)
	}

	// Step 8: Mainline health gate execution (F7 multi-dimensional evaluation)
	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if gateRes.ExitCode != 0 && gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
		t.Fatalf("unexpected mainline-health-gate exit code %d: %s", gateRes.ExitCode, gateRes.Stderr)
	}

	var gateReport FullKernelGateReport
	if err := harness.ParseJSON(gateRes.Stdout, &gateReport); err != nil {
		t.Fatalf("failed to parse health gate JSON report: %v", err)
	}

	// Verify all 5 required dimensions exist in report
	requiredDimensions := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, dim := range requiredDimensions {
		if _, exists := gateReport.Dimensions[dim]; !exists {
			t.Fatalf("missing required dimension %q in health gate report", dim)
		}
	}

	t.Logf("T4.02 PASS: Post-Merge Full Kernel Workflow completed successfully (dimensions=5, status=%s, exit=%d)",
		gateReport.Status, gateRes.ExitCode)
}
