package main

import (
	"bytes"
	"fmt"
	"os"
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

	cmd := exec.Command("git", "-C", root, "diff", "--name-only", "-z", "--no-renames", "--diff-filter=ACDMRTUXB", baseCommit, headCommit, "--")
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
	if gitCommitExists(root, commit) {
		return nil
	}
	fetch := exec.Command("git", "-C", root, "fetch", "--no-tags", "--depth=1", "origin", commit)
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
	cmd := exec.Command("git", "-C", root, "cat-file", "-e", commit+"^{commit}")
	return cmd.Run() == nil
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
