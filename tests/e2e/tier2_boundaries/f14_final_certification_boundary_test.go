package tier2_boundaries

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

func initBareOriginAndClone(t *testing.T, prefix string) (string, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	base := harness.TempDir(t, prefix)
	bareDir := filepath.Join(base, "origin.git")
	cloneDir := filepath.Join(base, "work")

	if res := harness.RunCmd(ctx, base, nil, "git", "init", "--bare", "-b", "main", bareDir); res.ExitCode != 0 {
		t.Fatalf("git init --bare failed: %s\n%s", res.Stderr, res.Stdout)
	}
	if res := harness.RunCmd(ctx, base, nil, "git", "clone", bareDir, cloneDir); res.ExitCode != 0 {
		t.Fatalf("git clone bare failed: %s\n%s", res.Stderr, res.Stdout)
	}
	for _, cfg := range [][]string{
		{"config", "user.name", "Test Writer"},
		{"config", "user.email", "test@cmdr.local"},
		{"checkout", "-B", "main"},
	} {
		if res := harness.RunCmd(ctx, cloneDir, nil, "git", cfg...); res.ExitCode != 0 {
			t.Fatalf("git %v failed: %s", cfg, res.Stderr)
		}
	}
	if err := os.WriteFile(filepath.Join(cloneDir, "README.md"), []byte("seed\n"), 0644); err != nil {
		t.Fatalf("write seed: %v", err)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "add", "README.md"); res.ExitCode != 0 {
		t.Fatalf("git add failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "commit", "-m", "initial main"); res.ExitCode != 0 {
		t.Fatalf("git commit failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "push", "-u", "origin", "main"); res.ExitCode != 0 {
		t.Fatalf("git push main failed: %s", res.Stderr)
	}
	return bareDir, cloneDir
}

func seedMinimalHealthGateDir(t *testing.T, dir string) {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	for _, rel := range []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
	} {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dst, err)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("write %s: %v", dst, err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "cmdr-product-spec"), 0755); err != nil {
		t.Fatalf("mkdir cmdr-product-spec: %v", err)
	}
}

// T2.F14.1: Final certification fails closed if any remote feature branch exists.
func TestT2_F14_1_CertificationFailsClosedOnRemoteFeatureBranch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, cloneDir := initBareOriginAndClone(t, "f14-branches")

	cleanBranches, err := harness.GitRemoteBranchesInDir(ctx, cloneDir)
	if err != nil {
		t.Fatalf("T2.F14.1 setup failed: %v", err)
	}
	for _, b := range cleanBranches {
		if !strings.HasSuffix(b, "/main") {
			t.Fatalf("T2.F14.1 setup expected only origin/main, got %v", cleanBranches)
		}
	}

	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "checkout", "-b", "feat/unmerged-studio-drafts"); res.ExitCode != 0 {
		t.Fatalf("git checkout feature branch failed: %s", res.Stderr)
	}
	if err := os.WriteFile(filepath.Join(cloneDir, "draft.txt"), []byte("draft\n"), 0644); err != nil {
		t.Fatalf("write draft: %v", err)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "add", "draft.txt"); res.ExitCode != 0 {
		t.Fatalf("git add draft failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "commit", "-m", "draft commit"); res.ExitCode != 0 {
		t.Fatalf("git commit draft failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "push", "-u", "origin", "feat/unmerged-studio-drafts"); res.ExitCode != 0 {
		t.Fatalf("git push feature branch failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "fetch", "origin"); res.ExitCode != 0 {
		t.Fatalf("git fetch origin failed: %s", res.Stderr)
	}

	dirtyBranches, err := harness.GitRemoteBranchesInDir(ctx, cloneDir)
	if err != nil {
		t.Fatalf("GitRemoteBranchesInDir failed: %v", err)
	}
	foundOrphan := false
	for _, b := range dirtyBranches {
		if b == "origin/feat/unmerged-studio-drafts" {
			foundOrphan = true
		}
	}
	if !foundOrphan {
		t.Fatalf("T2.F14.1 FAIL: expected origin/feat/unmerged-studio-drafts to be detected, got %v", dirtyBranches)
	}
	t.Log("T2.F14.1 PASS: real remote feature branch detected via GitRemoteBranchesInDir")
}

// T2.F14.2: Final certification fails closed if any open PR exists or gh CLI is unavailable.
func TestT2_F14_2_CertificationFailsClosedOnOpenPR(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resOne, errOne := harness.ParseOpenPRsJSON(`[{"number":234}]`)
	if errOne != nil || !resOne.Verified || resOne.Count != 1 {
		t.Fatalf("T2.F14.2 FAIL: expected Verified Count=1, got res=%+v err=%v", resOne, errOne)
	}

	resZero, errZero := harness.ParseOpenPRsJSON(`[]`)
	if errZero != nil || !resZero.Verified || resZero.Count != 0 {
		t.Fatalf("T2.F14.2 FAIL: expected Verified Count=0, got res=%+v err=%v", resZero, errZero)
	}

	for _, bad := range []string{"", "null", `[{"wrong":"obj"}]`} {
		resBad, errBad := harness.ParseOpenPRsJSON(bad)
		if !errors.Is(errBad, harness.ErrUnverified) || resBad.Verified {
			t.Fatalf("T2.F14.2 FAIL: expected ErrUnverified for %q, got res=%+v err=%v", bad, resBad, errBad)
		}
	}

	client, err := harness.NewGitHubClient()
	if err != nil {
		t.Fatalf("NewGitHubClient failed: %v", err)
	}
	emptyBin := harness.TempDir(t, "f14-empty-path")
	t.Setenv("PATH", emptyBin)
	prRes, prErr := client.CountOpenPRs(ctx)
	if !errors.Is(prErr, harness.ErrUnverified) || prRes.Verified || prRes.Status != harness.StatusUnverified {
		t.Fatalf("T2.F14.2 FAIL: expected CountOpenPRs to fail closed with ErrUnverified when gh unavailable, got res=%+v err=%v", prRes, prErr)
	}
	t.Log("T2.F14.2 PASS: ParseOpenPRsJSON and CountOpenPRs fail closed on open PRs, invalid payloads, and unavailable gh CLI")
}

// T2.F14.3: Final certification fails closed if unpushed commits exist.
func TestT2_F14_3_CertificationFailsClosedOnUnpushedCommits(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, cloneDir := initBareOriginAndClone(t, "f14-ahead")

	aheadZero, err := harness.GitAheadCountInDir(ctx, cloneDir, "origin/main", "HEAD")
	if err != nil || aheadZero != 0 {
		t.Fatalf("T2.F14.3 setup failed: expected ahead=0, got %d (err=%v)", aheadZero, err)
	}

	if err := os.WriteFile(filepath.Join(cloneDir, "unpushed.txt"), []byte("local only\n"), 0644); err != nil {
		t.Fatalf("write unpushed.txt: %v", err)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "add", "unpushed.txt"); res.ExitCode != 0 {
		t.Fatalf("git add failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, cloneDir, nil, "git", "commit", "-m", "unpushed commit"); res.ExitCode != 0 {
		t.Fatalf("git commit failed: %s", res.Stderr)
	}

	aheadOne, err := harness.GitAheadCountInDir(ctx, cloneDir, "origin/main", "HEAD")
	if err != nil || aheadOne != 1 {
		t.Fatalf("T2.F14.3 FAIL: expected ahead=1 after local unpushed commit, got %d (err=%v)", aheadOne, err)
	}
	t.Log("T2.F14.3 PASS: unpushed commits ahead of origin/main strictly detected via GitAheadCountInDir")
}

// T2.F14.4: Final certification rejects MAIN_HEALTHY if any health dimension is UNVERIFIED.
func TestT2_F14_4_CertificationRejectsUnverifiedDimension(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.TempDir(t, "f14-unverified")
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "init", "-b", "main"); res.ExitCode != 0 {
		t.Fatalf("git init failed: %s", res.Stderr)
	}
	seedMinimalHealthGateDir(t, mockRoot)
	for _, cfg := range [][]string{
		{"config", "user.name", "Test Writer"},
		{"config", "user.email", "test@cmdr.local"},
		{"add", "."},
		{"commit", "-m", "seed without baseline history"},
	} {
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", cfg...); res.ExitCode != 0 {
			t.Fatalf("git %v failed: %s", cfg, res.Stderr)
		}
	}

	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if gateRes.ExitCode == 0 {
		t.Fatalf("T2.F14.4 FAIL: mainline-health-gate exited 0 despite missing baseline history")
	}
	var report struct {
		Status     string `json:"status"`
		Dimensions map[string]struct {
			Status string `json:"status"`
			Passed bool   `json:"passed"`
		} `json:"dimensions"`
	}
	if err := harness.ParseJSON(gateRes.Stdout, &report); err != nil {
		t.Fatalf("T2.F14.4 FAIL: failed to parse mainline-health-gate output: %v\nStdout: %s", err, gateRes.Stdout)
	}
	if report.Status == "MAIN_HEALTHY" {
		t.Fatalf("T2.F14.4 FAIL: status reported MAIN_HEALTHY despite unverified git dimension")
	}
	if dim, ok := report.Dimensions["git"]; !ok || dim.Status != "UNVERIFIED" {
		t.Fatalf("T2.F14.4 FAIL: expected git dimension UNVERIFIED, got %+v", dim)
	}

	unitRes, unitErr := harness.RunGoTestVerified(
		ctx,
		"tools/cmdr-dev",
		"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev",
		[]string{"TestRunMainlineHealthGate_OverallStatusAggregation"},
		"-run", "^TestRunMainlineHealthGate_OverallStatusAggregation$",
	)
	if !unitRes.Passed || unitErr != nil {
		t.Fatalf("T2.F14.4 FAIL: TestRunMainlineHealthGate_OverallStatusAggregation failed: %v\nOutput: %s", unitErr, unitRes.Output)
	}
	t.Log("T2.F14.4 PASS: certification strictly rejects MAIN_HEALTHY when git dimension is UNVERIFIED")
}

// T2.F14.5: Final certification asserts HEAD commit is verified.
func TestT2_F14_5_CertificationAssertsHEADCommitVerified(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	headSHA, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to get git HEAD commit: %v", err)
	}
	if !harness.GitCommitExists(ctx, headSHA) {
		t.Fatalf("T2.F14.5 FAIL: HEAD commit %s not found in git object store", headSHA)
	}

	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if gateRes.ExitCode != 0 && gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
		t.Fatalf("T2.F14.5 FAIL: mainline-health-gate exited %d: %s", gateRes.ExitCode, gateRes.Stderr)
	}
	var liveReport struct {
		HeadCommit string `json:"head_commit"`
	}
	if err := harness.ParseJSON(gateRes.Stdout, &liveReport); err != nil {
		t.Fatalf("T2.F14.5 FAIL: parse mainline-health-gate output: %v", err)
	}
	if !strings.EqualFold(liveReport.HeadCommit, headSHA) {
		t.Fatalf("T2.F14.5 FAIL: mainline-health-gate head_commit %q != git HEAD %q", liveReport.HeadCommit, headSHA)
	}

	if harness.GitCommitExists(ctx, "invalid-sha-xyz") || harness.GitCommitExists(ctx, "0000000000000000000000000000000000000000") {
		t.Fatalf("T2.F14.5 FAIL: GitCommitExists accepted invalid or zero SHA")
	}

	unbornDir := harness.TempDir(t, "f14-unborn")
	if res := harness.RunCmd(ctx, unbornDir, nil, "git", "init", "-b", "main"); res.ExitCode != 0 {
		t.Fatalf("git init unborn failed: %s", res.Stderr)
	}
	seedMinimalHealthGateDir(t, unbornDir)
	unbornGate := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", unbornDir, "--json")
	if unbornGate.ExitCode == 0 {
		t.Fatalf("T2.F14.5 FAIL: mainline-health-gate exited 0 on unborn HEAD")
	}
	var unbornReport struct {
		Status     string `json:"status"`
		Dimensions map[string]struct {
			Status string `json:"status"`
		} `json:"dimensions"`
	}
	if err := harness.ParseJSON(unbornGate.Stdout, &unbornReport); err != nil {
		t.Fatalf("T2.F14.5 FAIL: parse unborn mainline-health-gate output: %v", err)
	}
	if unbornReport.Status == "MAIN_HEALTHY" || unbornReport.Dimensions["git"].Status != "UNVERIFIED" {
		t.Fatalf("T2.F14.5 FAIL: expected unborn repo git=UNVERIFIED and non-healthy status, got %+v", unbornReport)
	}

	t.Logf("T2.F14.5 PASS: HEAD commit %s verified in mainline-health-gate and unborn/invalid SHAs rejected", headSHA)
}
