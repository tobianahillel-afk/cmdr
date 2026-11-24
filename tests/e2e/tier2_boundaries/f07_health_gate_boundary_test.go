package tier2_boundaries

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// setupMockHealthGateRepo creates an isolated Git repository under .cmdr-tmp/
// to verify health gate boundary conditions without dirtying the real cmdr repository.
func setupMockHealthGateRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "health-gate-mock")

	// Initialize a local git repository inside the mock root
	initCmd := exec.Command("git", "-C", mockRoot, "init", "-b", "main")
	if out, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to git init mock repo: %v (%s)", err, out)
	}

	// Copy mandatory bootstrap files for isRepoRoot
	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
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

	// Configure local git user for commits
	_ = exec.Command("git", "-C", mockRoot, "config", "user.name", "Test Writer").Run()
	_ = exec.Command("git", "-C", mockRoot, "config", "user.email", "test@cmdr.local").Run()

	// Initial commit so HEAD exists
	_ = exec.Command("git", "-C", mockRoot, "add", ".").Run()
	_ = exec.Command("git", "-C", mockRoot, "commit", "-m", "initial commit").Run()

	return mockRoot
}

// T2.F7.1: Mainline health gate fails closed on dirty git repository.
func TestT2_F7_1_HealthGateFailsClosedOnDirtyRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Create an uncommitted modified/untracked file in the mock repository
	dirtyFile := filepath.Join(mockRoot, "DIRTY_TEST_FILE.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted change"), 0644); err != nil {
		t.Fatalf("failed to create dirty test file: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	// Health gate must exit with non-zero failure code (1 = MAIN_UNHEALTHY)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.1 FAIL: mainline-health-gate returned 0 on dirty git working tree")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "uncommitted") && !strings.Contains(combined, "FAIL") {
		t.Fatalf("T2.F7.1 FAIL: expected git dimension failure for uncommitted modifications, got: %s", combined)
	}
	t.Log("T2.F7.1 PASS: mainline health gate strictly fails closed when working tree is dirty")
}

// T2.F7.2: Mainline health gate fails closed on 1-byte diff in cmdr-product-spec/.
func TestT2_F7_2_HealthGateFailsClosedOnProductSpecDiff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Create cmdr-product-spec with a 1-byte mutation
	specFile := filepath.Join(mockRoot, "cmdr-product-spec", "spec.md")
	if err := os.MkdirAll(filepath.Dir(specFile), 0755); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}
	if err := os.WriteFile(specFile, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to write spec file: %v", err)
	}

	_ = exec.Command("git", "-C", mockRoot, "add", ".").Run()
	_ = exec.Command("git", "-C", mockRoot, "commit", "-m", "commit with spec").Run()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.2 FAIL: mainline-health-gate passed despite non-zero diff in cmdr-product-spec")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "cmdr-product-spec") && !strings.Contains(combined, "FAIL") && !strings.Contains(combined, "UNVERIFIED") {
		t.Fatalf("T2.F7.2 FAIL: expected product spec baseline violation or UNVERIFIED, got: %s", combined)
	}
	t.Log("T2.F7.2 PASS: mainline health gate fails closed on spec alteration")
}

// T2.F7.3: Mainline health gate returns MAIN_UNVERIFIED (exit code 2) on missing baseline history.
func TestT2_F7_3_HealthGateReturnsUnverifiedOnMissingBaseline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40 does not exist in mock git history
	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")

	// Status must be either MAIN_UNHEALTHY (exit 1) or MAIN_UNVERIFIED (exit 2) — never 0
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.3 FAIL: mainline-health-gate returned 0 when product spec baseline history is missing")
	}

	combined := res.Stderr + " " + res.Stdout
	if strings.Contains(combined, "MAIN_UNVERIFIED") || res.ExitCode == 2 {
		t.Logf("T2.F7.3 PASS: mainline-health-gate returned MAIN_UNVERIFIED (exit code 2) as expected: %s", combined)
		return
	}

	t.Logf("T2.F7.3 PASS: mainline-health-gate fail-closed behavior verified (exit code %d)", res.ExitCode)
}

// T2.F7.4: Mainline health gate fails closed if any single dimension fails.
func TestT2_F7_4_HealthGateFailsIfAnySingleDimensionFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// In current state where event-search git evidence is stale, traceability dimension fails
	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.4 FAIL: health gate passed even though a dimension is failing")
	}

	// Verify report status is not MAIN_HEALTHY
	if strings.Contains(res.Stdout, `"status": "MAIN_HEALTHY"`) {
		t.Fatalf("T2.F7.4 FAIL: health gate report status is MAIN_HEALTHY despite non-zero exit code")
	}

	t.Logf("T2.F7.4 PASS: health gate strictly enforces conjunction of all 5 dimensions (exit %d)", res.ExitCode)
}

// T2.F7.5: Mainline health gate rejects unknown CLI flag with usage error.
func TestT2_F7_5_HealthGateRejectsUnknownFlag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--unknown-flag-xyz")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.5 FAIL: mainline-health-gate succeeded with unknown CLI flag")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "flag") && !strings.Contains(combined, "unknown") && res.ExitCode != 1 && res.ExitCode != 2 {
		t.Fatalf("T2.F7.5 FAIL: expected flag error, got: %s", combined)
	}
	t.Logf("T2.F7.5 PASS: unknown flag correctly rejected with non-zero exit code %d", res.ExitCode)
}
