package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveRepoPath(root, candidate string, allowMissing bool) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	rootAbs = filepath.Clean(rootAbs)
	rootInfo, err := os.Lstat(rootAbs)
	if err != nil {
		return "", fmt.Errorf("inspect repository root: %w", err)
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("repository root must not be a symlink")
	}

	path := candidate
	if !filepath.IsAbs(path) {
		path = filepath.Join(rootAbs, filepath.FromSlash(path))
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)

	rel, err := filepath.Rel(rootAbs, path)
	if err != nil {
		return "", fmt.Errorf("relativize repository path: %w", err)
	}
	if filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes repository root: %s", candidate)
	}

	current := rootAbs
	if rel != "." {
		for _, part := range strings.Split(rel, string(os.PathSeparator)) {
			if part == "" || part == "." {
				continue
			}
			current = filepath.Join(current, part)
			info, err := os.Lstat(current) // #nosec G703 -- current is built only from a filepath.Rel-confined repository path.
			if os.IsNotExist(err) && allowMissing {
				break
			}
			if err != nil {
				return "", fmt.Errorf("inspect repository path %s: %w", candidate, err)
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return "", fmt.Errorf("repository path traverses symlink: %s", candidate)
			}
		}
	}
	return path, nil
}

func lstatRepoEntry(root, candidate string) (os.FileInfo, error) {
	cleaned := filepath.Clean(filepath.FromSlash(candidate))
	parentCandidate := filepath.Dir(cleaned)
	parent, err := resolveRepoPath(root, parentCandidate, false)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(parent, filepath.Base(cleaned))
	info, err := os.Lstat(path) // #nosec G703 -- parent is repository-confined; final entry is intentionally lstat'ed to detect a symlink without following it.
	if err != nil {
		return nil, err
	}
	return info, nil
}

func readRepoFile(root, candidate string) ([]byte, error) {
	path, err := resolveRepoPath(root, candidate, false)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path) // #nosec G304,G703 -- resolveRepoPath proves repository containment and rejects symlink traversal.
}

func openRepoFile(root, candidate string) (*os.File, error) {
	path, err := resolveRepoPath(root, candidate, false)
	if err != nil {
		return nil, err
	}
	return os.Open(path) // #nosec G304,G703 -- resolveRepoPath proves repository containment and rejects symlink traversal.
}

func statRepoPath(root, candidate string) (os.FileInfo, error) {
	path, err := resolveRepoPath(root, candidate, false)
	if err != nil {
		return nil, err
	}
	return os.Stat(path) // #nosec G703 -- resolveRepoPath proves repository containment and rejects symlink traversal.
}

func writeRepoFile(root, candidate string, data []byte) (string, error) {
	path, err := resolveRepoPath(root, candidate, true)
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(path)
	if _, err := resolveRepoPath(root, parent, true); err != nil {
		return "", err
	}
	if err := os.MkdirAll(parent, 0o750); err != nil { // #nosec G301,G703 -- directory is repository-confined and intentionally owner/group only.
		return "", err
	}
	path, err = resolveRepoPath(root, path, true)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil { // #nosec G306,G703 -- path is repository-confined; generated evidence is intentionally owner-only.
		return "", err
	}
	return path, nil
}
