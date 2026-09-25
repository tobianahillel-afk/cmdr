package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var fullCommitIDPattern = regexp.MustCompile(`^(?:[0-9a-fA-F]{40}|[0-9a-fA-F]{64})$`)

type GitChangesSummary struct {
	BaseCommit   string   `json:"base_commit"`
	HeadCommit   string   `json:"head_commit"`
	ChangedPaths []string `json:"changed_paths"`
	Count        int      `json:"count"`
	Output       string   `json:"output"`
}

func runGitChanges(root, baseCommit, headCommit, output string) (GitChangesSummary, error) {
	baseCommit, err := validateFullCommitID(baseCommit)
	if err != nil {
		return GitChangesSummary{}, fmt.Errorf("base commit: %w", err)
	}
	headCommit, err = validateFullCommitID(headCommit)
	if err != nil {
		return GitChangesSummary{}, fmt.Errorf("head commit: %w", err)
	}
	if baseCommit == headCommit {
		return GitChangesSummary{}, fmt.Errorf("base and head commits are identical")
	}

	if err := ensureGitCommit(root, baseCommit); err != nil {
		return GitChangesSummary{}, fmt.Errorf("base commit %s: %w", baseCommit, err)
	}
	if err := ensureGitCommit(root, headCommit); err != nil {
		return GitChangesSummary{}, fmt.Errorf("head commit %s: %w", headCommit, err)
	}

	cmd := exec.Command("git", "-C", root, "diff", "--name-only", "-z", "--no-renames", "--diff-filter=ACDMRTUXB", baseCommit, headCommit, "--") // #nosec G204,G702 -- executable is fixed, root is the resolved repository root, and both commit ids are full-hex validated before argument passing; no shell is used.
	raw, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if ok := errorAs(err, &exitErr); ok {
			return GitChangesSummary{}, fmt.Errorf("git diff failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return GitChangesSummary{}, fmt.Errorf("git diff failed: %w", err)
	}

	paths, err := parseNULPaths(raw)
	if err != nil {
		return GitChangesSummary{}, err
	}
	if len(paths) == 0 {
		return GitChangesSummary{}, fmt.Errorf("commit range contains no changed repository paths")
	}

	if output == "" || output == "engineering/spec-index/inventory.json" {
		output = filepath.ToSlash(filepath.Join("engineering", "testing", "changes.txt"))
	}
	outputPath, err := resolveRepoPath(root, output, true)
	if err != nil {
		return GitChangesSummary{}, err
	}
	data := []byte(strings.Join(paths, "\n") + "\n")
	outputPath, err = writeRepoFile(root, outputPath, data)
	if err != nil {
		return GitChangesSummary{}, fmt.Errorf("write change set: %w", err)
	}

	return GitChangesSummary{
		BaseCommit:   baseCommit,
		HeadCommit:   headCommit,
		ChangedPaths: paths,
		Count:        len(paths),
		Output:       outputPath,
	}, nil
}

func validateFullCommitID(value string) (string, error) {
	value = strings.TrimSpace(value)
	if !fullCommitIDPattern.MatchString(value) {
		return "", fmt.Errorf("expected a full 40- or 64-hex Git commit id")
	}
	allZero := true
	for _, ch := range value {
		if ch != '0' {
			allZero = false
			break
		}
	}
	if allZero {
		return "", fmt.Errorf("all-zero commit id is not valid change evidence")
	}
	return strings.ToLower(value), nil
}

func ensureGitCommit(root, commit string) error {
	normalized, err := validateFullCommitID(commit)
	if err != nil {
		return err
	}
	commit = normalized
	if gitCommitExists(root, commit) {
		return nil
	}
	fetch := exec.Command("git", "-C", root, "fetch", "--no-tags", "--depth=1", "origin", commit) // #nosec G204,G702 -- executable/flags/remote are fixed and commit is revalidated as full hexadecimal; no shell is used.
	var stderr bytes.Buffer
	fetch.Stderr = &stderr
	if err := fetch.Run(); err != nil {
		return fmt.Errorf("commit unavailable and exact fetch failed: %s", strings.TrimSpace(stderr.String()))
	}
	if !gitCommitExists(root, commit) {
		return fmt.Errorf("commit is still unavailable after exact fetch")
	}
	return nil
}

func gitCommitExists(root, commit string) bool {
	normalized, err := validateFullCommitID(commit)
	if err != nil {
		return false
	}
	commit = normalized
	cmd := exec.Command("git", "-C", root, "cat-file", "-e", commit+"^{commit}") // #nosec G204,G702 -- executable/flags are fixed and commit is full-hex validated; no shell is used.
	return cmd.Run() == nil
}

func ensureGitAncestor(root, ancestor, descendant string) error {
	ancestor, err := validateFullCommitID(ancestor)
	if err != nil {
		return fmt.Errorf("ancestor commit: %w", err)
	}
	descendant, err = validateFullCommitID(descendant)
	if err != nil {
		return fmt.Errorf("descendant commit: %w", err)
	}
	if err := ensureGitCommit(root, ancestor); err != nil {
		return fmt.Errorf("load ancestor commit: %w", err)
	}
	if err := ensureGitCommit(root, descendant); err != nil {
		return fmt.Errorf("load descendant commit: %w", err)
	}
	ok, err := gitIsAncestor(root, ancestor, descendant)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	shallow, err := gitRepositoryIsShallow(root)
	if err != nil {
		return err
	}
	if !shallow {
		return fmt.Errorf("%s is not an ancestor of %s", ancestor, descendant)
	}
	for _, depth := range []int{32, 128, 512, 2048} {
		if err := fetchGitCommitHistory(root, descendant, depth); err != nil {
			return fmt.Errorf("hydrate descendant history to depth %d: %w", depth, err)
		}
		ok, err := gitIsAncestor(root, ancestor, descendant)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return fmt.Errorf("unable to establish ancestry after bounded shallow-history hydration: ancestor=%s descendant=%s", ancestor, descendant)
}

func gitIsAncestor(root, ancestor, descendant string) (bool, error) {
	cmd := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", ancestor, descendant) // #nosec G204,G702 -- executable/flags are fixed and both commits are validated full hex IDs by the caller.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("check Git ancestry: %w", err)
	}
	return true, nil
}

func gitRepositoryIsShallow(root string) (bool, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--is-shallow-repository") // #nosec G204,G702 -- executable/arguments are fixed and root is the resolved repository root; no shell is used.
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("detect shallow Git repository: %w", err)
	}
	switch strings.TrimSpace(string(out)) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected Git shallow-repository response %q", strings.TrimSpace(string(out)))
	}
}

func fetchGitCommitHistory(root, commit string, depth int) error {
	commit, err := validateFullCommitID(commit)
	if err != nil {
		return err
	}
	if depth < 2 {
		return fmt.Errorf("invalid history depth %d", depth)
	}
	fetch := exec.Command("git", "-C", root, "fetch", "--no-tags", fmt.Sprintf("--depth=%d", depth), "origin", commit) // #nosec G204,G702 -- executable/flags/remote are fixed; commit is full-hex validated and depth is an integer selected by compiled code.
	var stderr bytes.Buffer
	fetch.Stderr = &stderr
	if err := fetch.Run(); err != nil {
		return fmt.Errorf("exact history fetch failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

func parseNULPaths(raw []byte) ([]string, error) {
	parts := bytes.Split(raw, []byte{0})
	var values []string
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		values = append(values, string(part))
	}
	paths, err := normalizeChangedPaths(values)
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

func errorAs(err error, target any) bool {
	switch t := target.(type) {
	case **exec.ExitError:
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			*t = exitErr
		}
		return ok
	default:
		return false
	}
}
