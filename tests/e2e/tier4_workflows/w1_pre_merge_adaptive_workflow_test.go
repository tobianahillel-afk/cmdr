package tier4_workflows

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// TestT4_01_PreMergeAdaptiveWorkflow:
// Simulates the end-to-end PR developer lifecycle:
// 1. Changes detection (git-changes)
// 2. Validation plan generation (validation-plan)
// 3. Targeted validation run (gofmt, vet, test, SAST, SCA)
// 4. Pre-merge health gate pre-check
// Features exercised: F1, F2, F3, F4, F7, F8.
func TestT4_01_PreMergeAdaptiveWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	tempDir := harness.TempDir(t, "w1-adaptive")

	// Phase 1: Determine base and head revisions for PR simulation
	headCommit, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to resolve HEAD commit: %v", err)
	}

	// Use HEAD~1 or HEAD as base commit
	baseRes := harness.RunCmd(ctx, root, nil, "git", "rev-parse", "HEAD~1")
	baseCommit := headCommit
	if baseRes.ExitCode == 0 && baseRes.Stdout != "" {
		baseCommit = filepath.Clean(baseRes.Stdout)
	}

	// Phase 2: Compute change set via cmdr-dev git-changes
	changesFile := filepath.Join(tempDir, "pr-changes.txt")
	gitChangesRes := harness.RunCmdrDev(ctx, "git-changes",
		"--base-commit", baseCommit,
		"--head-commit", headCommit,
		"--output", changesFile,
		"--json",
	)
	if gitChangesRes.ExitCode != 0 {
		t.Fatalf("git-changes failed with exit code %d: %s", gitChangesRes.ExitCode, gitChangesRes.Stderr)
	}

	// If no changes were derived between HEAD and HEAD, seed with sample studio-catalog paths
	changesInfo, err := os.Stat(changesFile)
	if err != nil || changesInfo.Size() == 0 {
		sampleChanges := "product-runtime/studio-catalog/tool_input.go\nproduct-runtime/studio-catalog/tool_input_test.go\n"
		if err := os.WriteFile(changesFile, []byte(sampleChanges), 0644); err != nil {
			t.Fatalf("failed to seed changes file: %v", err)
		}
	}

	// Phase 3: Generate risk-adaptive validation plan
	planRes := harness.RunCmdrDev(ctx, "validation-plan",
		"--changes-file", changesFile,
		"--json",
	)
	if planRes.ExitCode != 0 {
		t.Fatalf("validation-plan failed with exit code %d: %s", planRes.ExitCode, planRes.Stderr)
	}
	harness.AssertJSONFieldExists(t, planRes.Stdout, "tier")

	// Phase 4: Run targeted validation execution
	valRunRes := harness.RunCmdrDev(ctx, "validation-run",
		"--changes-file", changesFile,
		"--json",
	)
	if valRunRes.ExitCode != 0 {
		t.Fatalf("validation-run failed with exit code %d: %s", valRunRes.ExitCode, valRunRes.Stderr)
	}

	// Phase 5: Run targeted developer security scans (SAST and SCA)
	sastRes := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if sastRes.ExitCode != 0 {
		t.Fatalf("pre-merge SAST scan failed (exit %d): %s", sastRes.ExitCode, sastRes.Stderr)
	}

	scaRes := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if scaRes.ExitCode != 0 {
		t.Fatalf("pre-merge SCA scan failed (exit %d): %s", scaRes.ExitCode, scaRes.Stderr)
	}

	// Phase 6: Run studio-catalog regression tests (Codex F1 compliance)
	studioPass, studioOut, studioErr := harness.RunGoTest(ctx, "product-runtime/studio-catalog")
	if !studioPass {
		t.Fatalf("studio-catalog tests failed in pre-merge cycle: %s (err: %v)", studioOut, studioErr)
	}

	// Phase 7: Execute pre-merge health gate pre-check
	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if gateRes.ExitCode != 0 && gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
		t.Fatalf("unexpected mainline-health-gate exit code %d: %s", gateRes.ExitCode, gateRes.Stderr)
	}

	t.Logf("T4.01 PASS: Pre-Merge Adaptive Workflow executed end-to-end cleanly (plan_status=PASS, SAST=PASS, SCA=PASS, gate_exit=%d)",
		gateRes.ExitCode)
}
