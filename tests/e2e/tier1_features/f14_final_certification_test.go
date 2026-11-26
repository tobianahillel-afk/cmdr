package tier1_features

import (
	"context"
	"errors"
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

	if strings.TrimSpace(status) != "" {
		t.Fatalf("working tree has uncommitted changes:\n%s", strings.TrimSpace(status))
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

	prRes, err := client.CountOpenPRs(ctx)
	if errors.Is(err, harness.ErrUnverified) || !prRes.Verified {
		t.Fatalf("T1.F14.3 UNVERIFIED: cannot certify 0 open PRs without verified GitHub API response: %v", err)
	}
	if err != nil {
		t.Fatalf("failed to count open PRs: %v", err)
	}

	if prRes.Count != 0 {
		t.Fatalf("expected 0 open PRs, found %d", prRes.Count)
	}
	t.Logf("T1.F14.3 PASS: exactly 0 open PRs on GitHub (status=%s)", prRes.Status)
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
	if res.ExitCode != 0 && res.ExitCode != 1 && res.ExitCode != 2 {
		t.Fatalf("mainline-health-gate failed unexpectedly (exit %d): %s\n%s", res.ExitCode, res.Stderr, res.Stdout)
	}

	var report struct {
		Status     string                 `json:"status"`
		Dimensions map[string]interface{} `json:"dimensions"`
	}
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse mainline-health-gate JSON report: %v\nStdout: %s", err, res.Stdout)
	}
	for _, dim := range []string{"code", "security", "git", "traceability", "operational_quality"} {
		if _, ok := report.Dimensions[dim]; !ok {
			t.Fatalf("mainline-health-gate report missing dimension %q", dim)
		}
	}

	if strings.Contains(res.Stdout, "\"mock\"") || strings.Contains(res.Stdout, "\"stub\"") {
		t.Fatal("mainline-health-gate report contains mock/stub evidence flags")
	}
	t.Log("T1.F14.5 PASS: zero mock/stub flags in mainline health gate report")
}
