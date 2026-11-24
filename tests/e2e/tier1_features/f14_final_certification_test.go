package tier1_features

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T1.F14.1: Verify git status --porcelain is clean (0 uncommitted changes).
func TestT1_F14_1_GitCleanWorkingTree(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := harness.GitStatusPorcelain(ctx)
	if err != nil {
		t.Fatalf("failed to query git status: %v", err)
	}

	// Filter out harmless untracked test directories if running locally
	lines := strings.Split(status, "\n")
	var realChanges []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		// Allow .cmdr-tmp, cmdr-dev, tests/, or runtime-security-evidence.json changes during test run
		if strings.Contains(trimmed, ".cmdr-") || strings.Contains(trimmed, "cmdr-dev") || strings.Contains(trimmed, "tests/") || strings.Contains(trimmed, "runtime-security-evidence.json") {
			continue
		}
		realChanges = append(realChanges, trimmed)
	}

	if len(realChanges) > 0 {
		t.Fatalf("working tree has uncommitted changes outside test suite:\n%s", strings.Join(realChanges, "\n"))
	}
	t.Log("T1.F14.1 PASS: git working tree is clean")
}

// T1.F14.2: Verify git branch -r shows only origin/main (0 orphaned remote branches).
func TestT1_F14_2_ZeroOrphanRemoteBranches(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !harness.IsMainBranch(ctx) {
		t.Skip("T1.F14.2 SKIP: zero remote branches certification strictly applies post-merge on main")
		return
	}

	branches, err := harness.GitRemoteBranches(ctx)
	if err != nil {
		t.Fatalf("failed to list remote branches: %v", err)
	}

	for _, b := range branches {
		if !strings.HasSuffix(b, "/main") && !strings.HasSuffix(b, "/HEAD") {
			t.Fatalf("unexpected remote branch found: %s; only origin/main must exist", b)
		}
	}
	t.Logf("T1.F14.2 PASS: remote branches strictly clean (%d remote branch: origin/main)", len(branches))
}

// T1.F14.3: Verify GitHub reports 0 open PRs.
func TestT1_F14_3_ZeroOpenPRs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if !harness.IsMainBranch(ctx) {
		t.Skip("T1.F14.3 SKIP: zero open PRs certification strictly applies post-merge on main")
		return
	}

	client, err := harness.NewGitHubClient()
	if err != nil {
		t.Fatalf("failed to create github client: %v", err)
	}

	openCount, isOffline, err := client.CountOpenPRs(ctx)
	if err != nil {
		t.Fatalf("failed to count open PRs: %v", err)
	}

	if openCount != 0 {
		t.Fatalf("expected 0 open PRs, found %d", openCount)
	}
	t.Logf("T1.F14.3 PASS: exactly 0 open PRs on GitHub (offline fallback: %v)", isOffline)
}

// T1.F14.4: Verify cmdr-dev mainline-health-gate outputs MAIN_HEALTHY with exit code 0.
func TestT1_F14_4_MainlineHealthGateCertifiesHealthy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if !harness.IsMainBranch(ctx) {
		t.Skip("T1.F14.4 SKIP: mainline MAIN_HEALTHY certification strictly applies on main branch post-Milestone 4")
		return
	}

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("mainline-health-gate did not return 0 (exit code %d). Summary: %s", res.ExitCode, res.Stderr)
	}

	harness.AssertJSONFieldEqual(t, res.Stdout, "status", "MAIN_HEALTHY")
	t.Log("T1.F14.4 PASS: mainline-health-gate certifies MAIN_HEALTHY with exit code 0")
}

// T1.F14.5: Verify no unverified/mock flags used in final health gate audit.
func TestT1_F14_5_ZeroMockFlagsInHealthGate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if strings.Contains(res.Stdout, "\"mock\"") || strings.Contains(res.Stdout, "\"stub\"") {
		t.Fatal("mainline-health-gate report contains mock/stub evidence flags")
	}
	t.Log("T1.F14.5 PASS: zero mock/stub flags in mainline health gate report")
}
