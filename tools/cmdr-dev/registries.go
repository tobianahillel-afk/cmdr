package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func loadActiveScreenIDs(root, specRel string) (map[string]struct{}, error) {
	path := filepath.Join(root, filepath.FromSlash(specRel), "00-governance", "registers", "screen-register.md")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open screen register: %w", err)
	}
	defer f.Close()

	active := map[string]struct{}{}
	inActiveSection := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			switch {
			case strings.HasPrefix(trimmed, "## Écrans actifs"):
				inActiveSection = true
				continue
			case inActiveSection:
				inActiveSection = false
			}
		}
		if !inActiveSection {
			continue
		}
		if id, ok := firstTableCell(trimmed); ok && screenIDPattern.MatchString(id) {
			active[id] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read screen register: %w", err)
	}
	if len(active) == 0 {
		return nil, fmt.Errorf("screen register yielded no active screen ids")
	}
	return active, nil
}

func firstTableCell(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") {
		return "", false
	}
	rest := strings.TrimPrefix(line, "|")
	end := strings.Index(rest, "|")
	if end < 0 {
		return "", false
	}
	value := strings.TrimSpace(rest[:end])
	value = strings.Trim(value, "`*")
	if value == "" || strings.EqualFold(value, "id") || strings.EqualFold(value, "capability id") || strings.EqualFold(value, "permission") {
		return "", false
	}
	return value, true
}

func loadActiveCapabilityIDs(root, specRel string) (map[string][]string, error) {
	dir := filepath.Join(root, filepath.FromSlash(specRel), "00-governance", "registers")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read capability register directory: %w", err)
	}
	out := map[string][]string{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "capability-register-") || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		rel := filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers", entry.Name()))
		if err := collectFirstColumnIDs(filepath.Join(root, filepath.FromSlash(rel)), rel, capabilityIDPattern, out); err != nil {
			return nil, err
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("capability registers yielded no active capability ids")
	}
	normalizeRegistryEvidence(out)
	return out, nil
}

func loadActivePermissionIDs(root, specRel string) (map[string][]string, error) {
	rel := filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers", "permission-register.md"))
	out := map[string][]string{}
	path := filepath.Join(root, filepath.FromSlash(rel))
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open permission register: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id, ok := firstTableCell(scanner.Text())
		if ok && strings.HasPrefix(id, "perm.") {
			out[id] = append(out[id], rel)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read permission register: %w", err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("permission register yielded no active permission ids")
	}
	normalizeRegistryEvidence(out)
	return out, nil
}

type SourceRequirementBaseline struct {
	SchemaVersion  int      `json:"schema_version"`
	BaselineKind   string   `json:"baseline_kind"`
	SourceCommit   string   `json:"source_commit"`
	SourcePath     string   `json:"source_path"`
	SourceBlobSHA  string   `json:"source_blob_sha"`
	ExpectedCount  int      `json:"expected_count"`
	RequirementIDs []string `json:"requirement_ids"`
}

func loadActiveRequirementIDs(root, specRel string) (map[string][]string, error) {
	_ = specRel
	rel := "engineering/coverage/source-requirements-baseline.json"
	var baseline SourceRequirementBaseline
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(rel)), &baseline); err != nil {
		return nil, fmt.Errorf("decode source requirement baseline: %w", err)
	}
	if baseline.SchemaVersion != 1 || baseline.BaselineKind != "source-requirement-identity" {
		return nil, fmt.Errorf("unsupported source requirement baseline")
	}
	if baseline.SourceCommit == "" || baseline.SourcePath == "" || baseline.SourceBlobSHA == "" {
		return nil, fmt.Errorf("source requirement baseline lacks provenance")
	}
	if baseline.ExpectedCount <= 0 || len(baseline.RequirementIDs) != baseline.ExpectedCount {
		return nil, fmt.Errorf("source requirement baseline count mismatch: expected %d, got %d", baseline.ExpectedCount, len(baseline.RequirementIDs))
	}

	out := make(map[string][]string, len(baseline.RequirementIDs))
	for _, id := range baseline.RequirementIDs {
		if !requirementIDPattern.MatchString(id) {
			return nil, fmt.Errorf("invalid source requirement id %q", id)
		}
		if _, exists := out[id]; exists {
			return nil, fmt.Errorf("duplicate source requirement id %q", id)
		}
		out[id] = []string{rel}
	}
	return out, nil
}

func collectFirstColumnIDs(path, rel string, matcher interface{ MatchString(string) bool }, out map[string][]string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open registry %s: %w", rel, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		id, ok := firstTableCell(scanner.Text())
		if ok && matcher.MatchString(id) {
			out[id] = append(out[id], rel)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read registry %s: %w", rel, err)
	}
	return nil
}

func normalizeRegistryEvidence(values map[string][]string) {
	for id, paths := range values {
		sort.Strings(paths)
		unique := paths[:0]
		var previous string
		for i, path := range paths {
			if i == 0 || path != previous {
				unique = append(unique, path)
				previous = path
			}
		}
		values[id] = unique
	}
}

func loadActiveOpenDecisionIDs(root, specRel string) (map[string][]string, error) {
	rel := filepath.ToSlash(filepath.Join(specRel, "00-governance", "source-material", "unresolved-decisions.md"))
	path := filepath.Join(root, filepath.FromSlash(rel))
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open unresolved decisions: %w", err)
	}
	defer f.Close()

	out := map[string][]string{}
	inAudit := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(trimmed, "## ") {
			switch {
			case strings.HasPrefix(trimmed, "## Complete open-decision audit"):
				inAudit = true
				continue
			case inAudit:
				inAudit = false
			}
		}
		if !inAudit {
			continue
		}
		if id, ok := firstTableCell(trimmed); ok && openDecisionIDPattern.MatchString(id) {
			out[id] = append(out[id], rel)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read unresolved decisions: %w", err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("unresolved decision register yielded no active ids")
	}
	normalizeRegistryEvidence(out)
	return out, nil
}
