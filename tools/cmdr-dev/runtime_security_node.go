package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func measureNodeRuntimeCoverage(root, tempDir string, changedPaths []string, adapter runtimeSecurityAdapter) (RuntimeCoverageMeasurement, error) {
	moduleRoot, err := resolveRepoPath(root, adapter.RuntimeRoot, false)
	if err != nil {
		return RuntimeCoverageMeasurement{}, err
	}
	if strings.TrimSpace(adapter.TestFile) == "" || strings.TrimSpace(adapter.CoverageInclude) == "" {
		return RuntimeCoverageMeasurement{}, fmt.Errorf("%s node adapter is incomplete", adapter.BoundaryID)
	}
	lcovName := strings.ReplaceAll(adapter.BoundaryID, "/", "-") + ".lcov"
	lcovPath, err := resolveRepoPath(root, filepath.Join(tempDir, lcovName), true)
	if err != nil {
		return RuntimeCoverageMeasurement{}, err
	}
	// #nosec G204,G702 -- executable and flags are fixed; adapter values are compiled constants and no shell is used.
	cmd := exec.Command(
		"node",
		"--test",
		"--experimental-test-coverage",
		"--test-coverage-include="+adapter.CoverageInclude,
		"--test-coverage-lines=90",
		"--test-reporter=spec",
		"--test-reporter=lcov",
		"--test-reporter-destination=stdout",
		"--test-reporter-destination="+lcovPath,
		adapter.TestFile,
	)
	cmd.Dir = moduleRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return RuntimeCoverageMeasurement{}, fmt.Errorf("%s node coverage test failed: %s", adapter.BoundaryID, compactProcessDiagnostic(stdout.String(), stderr.String(), err))
	}
	data, err := readRepoFile(root, lcovPath)
	if err != nil {
		return RuntimeCoverageMeasurement{}, err
	}
	return parseNodeLCOV(root, moduleRoot, data, changedPaths, adapter)
}

func runNodeRuntimeNegativeTest(root string, adapter runtimeSecurityAdapter, pattern, label string) error {
	moduleRoot, err := resolveRepoPath(root, adapter.RuntimeRoot, false)
	if err != nil {
		return err
	}
	if strings.TrimSpace(adapter.TestFile) == "" {
		return fmt.Errorf("%s node adapter has no test file", adapter.BoundaryID)
	}
	// #nosec G204,G702 -- executable and flags are fixed; regex/test file are compiled constants and no shell is used.
	cmd := exec.Command("node", "--test", "--test-name-pattern="+pattern, adapter.TestFile)
	cmd.Dir = moduleRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s negative test failed: %s", adapter.BoundaryID, label, compactProcessDiagnostic(stdout.String(), stderr.String(), err))
	}
	return nil
}

func parseNodeLCOV(root, moduleRoot string, data []byte, changedPaths []string, adapter runtimeSecurityAdapter) (RuntimeCoverageMeasurement, error) {
	changed := map[string]bool{}
	for _, path := range changedPaths {
		changed[filepath.ToSlash(path)] = true
	}
	expected := filepath.ToSlash(filepath.Join(adapter.RuntimeRoot, filepath.FromSlash(adapter.CoverageInclude)))
	var current string
	var lf, lh int
	var total, covered, changedTotal, changedCovered int
	flush := func() error {
		if current == "" {
			return nil
		}
		if lf <= 0 || lh < 0 || lh > lf {
			return fmt.Errorf("invalid LCOV line totals for %s", current)
		}
		repoPath, err := normalizeNodeCoveragePath(root, moduleRoot, current)
		if err != nil {
			return err
		}
		if repoPath != expected {
			return fmt.Errorf("node coverage source %s is outside expected source %s", repoPath, expected)
		}
		total += lf
		covered += lh
		if changed[repoPath] {
			changedTotal += lf
			changedCovered += lh
		}
		current, lf, lh = "", 0, 0
		return nil
	}

	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		switch {
		case strings.HasPrefix(line, "SF:"):
			if current != "" {
				return RuntimeCoverageMeasurement{}, fmt.Errorf("LCOV source record missing end_of_record")
			}
			current = strings.TrimSpace(strings.TrimPrefix(line, "SF:"))
		case strings.HasPrefix(line, "LF:"):
			value, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "LF:")))
			if err != nil {
				return RuntimeCoverageMeasurement{}, fmt.Errorf("invalid LCOV LF value")
			}
			lf = value
		case strings.HasPrefix(line, "LH:"):
			value, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "LH:")))
			if err != nil {
				return RuntimeCoverageMeasurement{}, fmt.Errorf("invalid LCOV LH value")
			}
			lh = value
		case line == "end_of_record":
			if err := flush(); err != nil {
				return RuntimeCoverageMeasurement{}, err
			}
		}
	}
	if current != "" {
		if err := flush(); err != nil {
			return RuntimeCoverageMeasurement{}, err
		}
	}
	if total == 0 {
		return RuntimeCoverageMeasurement{}, fmt.Errorf("%s node LCOV contains no executable lines", adapter.BoundaryID)
	}
	result := RuntimeCoverageMeasurement{
		GlobalPercent:     100 * float64(covered) / float64(total),
		TotalStatements:   total,
		CoveredStatements: covered,
	}
	if changedTotal > 0 {
		result.ChangedExecutableStmts = changedTotal
		result.ChangedExecutablePercent = 100 * float64(changedCovered) / float64(changedTotal)
	}
	return result, nil
}

func normalizeNodeCoveragePath(root, moduleRoot, reported string) (string, error) {
	reported = strings.TrimSpace(strings.TrimPrefix(reported, "file://"))
	if reported == "" {
		return "", fmt.Errorf("node LCOV source path is empty")
	}
	path := reported
	if !filepath.IsAbs(path) {
		path = filepath.Join(moduleRoot, filepath.FromSlash(path))
	}
	rel, err := filepath.Rel(root, filepath.Clean(path))
	if err != nil {
		return "", err
	}
	normalized, err := normalizeChangedPaths([]string{filepath.ToSlash(rel)})
	if err != nil || len(normalized) != 1 {
		return "", fmt.Errorf("node LCOV source escapes repository: %s", reported)
	}
	return normalized[0], nil
}

func compactProcessDiagnostic(stdout, stderr string, runErr error) string {
	for _, value := range []string{stderr, stdout} {
		value = strings.TrimSpace(value)
		if value != "" {
			if len(value) > 4000 {
				return value[:4000] + "...[truncated]"
			}
			return value
		}
	}
	return runErr.Error()
}
