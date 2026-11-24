package harness

import (
	"context"
	"fmt"
	"strings"
)

// GitStatusPorcelain returns the output of `git status --porcelain`.
func GitStatusPorcelain(ctx context.Context) (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	res := RunCmd(ctx, root, nil, "git", "status", "--porcelain")
	if res.ExitCode != 0 {
		return "", fmt.Errorf("git status failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	return strings.TrimSpace(res.Stdout), nil
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

// GitRemoteBranches returns the list of remote branch references.
func GitRemoteBranches(ctx context.Context) ([]string, error) {
	root, err := RepoRoot()
	if err != nil {
		return nil, err
	}
	res := RunCmd(ctx, root, nil, "git", "branch", "-r")
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

// GitCommitExists checks if a commit SHA exists in git object storage.
func GitCommitExists(ctx context.Context, sha string) bool {
	root, err := RepoRoot()
	if err != nil {
		return false
	}
	trimmed := strings.TrimSpace(sha)
	if trimmed == "" {
		return false
	}
	res := RunCmd(ctx, root, nil, "git", "cat-file", "-e", trimmed+"^{commit}")
	return res.ExitCode == 0
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
