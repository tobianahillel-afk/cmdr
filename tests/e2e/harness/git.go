package harness

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// GitStatusPorcelainInDir returns the output of `git status --porcelain` in the specified directory.
func GitStatusPorcelainInDir(ctx context.Context, dir string) (string, error) {
	if strings.TrimSpace(dir) == "" {
		return "", fmt.Errorf("git status directory must be non-empty")
	}
	res := RunCmd(ctx, dir, nil, "git", "status", "--porcelain")
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git status failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
}

// GitStatusPorcelain returns the output of `git status --porcelain`.
func GitStatusPorcelain(ctx context.Context) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	return GitStatusPorcelainInDir(ctx, root)
}

// GitDiffSpecBaseline checks for differences in cmdr-product-spec against a baseline commit.
func GitDiffSpecBaseline(ctx context.Context, baselineCommit string) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	res := RunCmd(ctx, root, nil, "git", "diff", baselineCommit, "--", "cmdr-product-spec")
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git diff failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
}

// GitRemoteBranchesInDir returns the list of remote branch references in the specified directory.
func GitRemoteBranchesInDir(ctx context.Context, dir string) ([]string, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, fmt.Errorf("git branch directory must be non-empty")
	}
	res := RunCmd(ctx, dir, nil, "git", "branch", "-r")
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("git branch -r failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	lines := strings.Split(res.Stdout, "\n")
	var branches []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.Contains(trimmed, "->") {
			branches = append(branches, trimmed)
		}
	}
	return branches, nil
}

// GitRemoteBranches returns the list of remote branch references.
func GitRemoteBranches(ctx context.Context) ([]string, error) {
	root, err := RepoRoot()
	if err != nil {
		return nil, err
	}
	return GitRemoteBranchesInDir(ctx, root)
}

// GitCommitExistsInDir checks if a commit SHA exists in git object storage in the specified directory.
func GitCommitExistsInDir(ctx context.Context, dir, sha string) bool {
	if strings.TrimSpace(dir) == "" {
		return false
	}
	trimmed := strings.TrimSpace(sha)
	if trimmed == "" {
		return false
	}
	res := RunCmd(ctx, dir, nil, "git", "cat-file", "-e", trimmed+"^{commit}")
	return res.ExitCode == 0
}

// GitCommitExists checks if a commit SHA exists in git object storage.
func GitCommitExists(ctx context.Context, sha string) bool {
	root, err := RepoRoot()
	if err != nil {
		return false
	}
	return GitCommitExistsInDir(ctx, root, sha)
}

// GitAheadCountInDir returns the number of commits reachable from headRef but not from upstreamRef.
func GitAheadCountInDir(ctx context.Context, dir, upstreamRef, headRef string) (int, error) {
	if strings.TrimSpace(dir) == "" {
		return -1, fmt.Errorf("git rev-list directory must be non-empty")
	}
	up := strings.TrimSpace(upstreamRef)
	hd := strings.TrimSpace(headRef)
	if up == "" || hd == "" {
		return -1, fmt.Errorf("upstreamRef and headRef must be non-empty")
	}
	res := RunCmd(ctx, dir, nil, "git", "rev-list", "--count", up+".."+hd)
	if res.ExitCode != 0 {
		return -1, fmt.Errorf("git rev-list --count %s..%s failed (code %d): %s", up, hd, res.ExitCode, res.Stderr)
	}
	var count int
	if _, err := fmt.Sscanf(strings.TrimSpace(res.Stdout), "%d", &count); err != nil {
		return -1, fmt.Errorf("failed to parse git rev-list --count output %q: %w", res.Stdout, err)
	}
	return count, nil
}

// GitHeadCommit returns the current HEAD commit hash.
func GitHeadCommit(ctx context.Context) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	res := RunCmd(ctx, root, nil, "git", "rev-parse", "HEAD")
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git rev-parse HEAD failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
}

// GitDiffPath returns the git diff for a specific directory or file between two revisions.
func GitDiffPath(ctx context.Context, rev1, rev2, relPath string) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	res := RunCmd(ctx, root, nil, "git", "diff", rev1, rev2, "--", relPath)
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git diff failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
}

// GitCurrentBranch returns the current checked-out branch name.
func GitCurrentBranch(ctx context.Context) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	res := RunCmd(ctx, root, nil, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git rev-parse --abbrev-ref HEAD failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
}

// IsMainBranch reports whether test execution is currently running against the main branch.
func IsMainBranch(ctx context.Context) bool {
	if os.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		return false
	}
	if ref := os.Getenv("GITHUB_REF_NAME"); ref != "" && ref != "main" {
		return false
	}
	branch, err := GitCurrentBranch(ctx)
	if err == nil && branch != "" && branch != "main" {
		return false
	}
	return true
}

// CloneIsolatedRepo creates an isolated Git clone under .cmdr-tmp/ using --shared --no-checkout
// and checks out the exact current HEAD commit onto a local 'main' branch.
// Works on both named branches and detached-HEAD CI checkouts without mutating RepoRoot().
func CloneIsolatedRepo(t *testing.T, prefix string) string {
	t.Helper()
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	headCommit, err := GitHeadCommit(ctx)
	if err != nil || strings.TrimSpace(headCommit) == "" {
		t.Fatalf("failed to resolve HEAD commit: %v", err)
	}
	headCommit = strings.TrimSpace(headCommit)

	mockRoot := TempDir(t, prefix)

	cloneRes := RunCmd(ctx, root, nil, "git", "clone", "--shared", "--no-checkout", root, mockRoot)
	if cloneRes.ExitCode != 0 {
		t.Fatalf("git clone --shared failed (exit %d): %s\n%s", cloneRes.ExitCode, cloneRes.Stderr, cloneRes.Stdout)
	}

	for _, cfg := range [][]string{
		{"config", "user.name", "Test Writer"},
		{"config", "user.email", "test@cmdr.local"},
	} {
		cfgRes := RunCmd(ctx, mockRoot, nil, "git", cfg...)
		if cfgRes.ExitCode != 0 {
			t.Fatalf("git %v in mockRoot failed (exit %d): %s", cfg, cfgRes.ExitCode, cfgRes.Stderr)
		}
	}

	coRes := RunCmd(ctx, mockRoot, nil, "git", "checkout", "-B", "main", headCommit)
	if coRes.ExitCode != 0 {
		t.Fatalf("git checkout -B main %s failed (exit %d): %s\n%s", headCommit, coRes.ExitCode, coRes.Stderr, coRes.Stdout)
	}

	return mockRoot
}
