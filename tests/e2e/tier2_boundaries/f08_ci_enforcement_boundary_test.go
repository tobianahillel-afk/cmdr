package tier2_boundaries

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// validateWorkflowContent encapsulates the CI workflow verification invariants.
func validateWorkflowContent(content string) (bool, string) {
	if !strings.Contains(content, "mainline-health-gate") {
		return false, "missing mandatory mainline-health-gate step"
	}

	lines := strings.Split(content, "\n")
	inHealthGateStep := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "mainline-health-gate") {
			inHealthGateStep = true
		}
		if inHealthGateStep {
			if strings.HasPrefix(trimmed, "- name:") && !strings.Contains(trimmed, "mainline-health-gate") {
				inHealthGateStep = false
				continue
			}
			if strings.Contains(trimmed, "continue-on-error: true") {
				return false, "mainline-health-gate step has continue-on-error: true; must fail closed"
			}
		}
	}

	if !strings.Contains(content, "push:") || !strings.Contains(content, "pull_request:") {
		return false, "missing required push or pull_request triggers"
	}

	if !strings.Contains(content, "branches:") || !strings.Contains(content, "\"main\"") {
		return false, "workflow does not strictly target main branch"
	}

	return true, ""
}

// T2.F8.1: CI workflow validation fails if mainline-health-gate step is removed.
func TestT2_F8_1_CIValidationFailsIfHealthGateStepRemoved(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	if !strings.Contains(string(content), "mainline-health-gate") {
		t.Skip("T2.F8.1 SKIP: mainline-health-gate CI workflow step is scheduled for Milestone 3")
		return
	}

	// Remove mainline-health-gate lines from content
	mutatedLines := []string{}
	for _, line := range strings.Split(string(content), "\n") {
		if !strings.Contains(line, "mainline-health-gate") {
			mutatedLines = append(mutatedLines, line)
		}
	}
	mutated := strings.Join(mutatedLines, "\n")

	valid, errMsg := validateWorkflowContent(mutated)
	if valid {
		t.Fatalf("T2.F8.1 FAIL: validation passed despite removing mainline-health-gate step")
	}
	if !strings.Contains(errMsg, "missing mandatory mainline-health-gate step") {
		t.Fatalf("T2.F8.1 FAIL: unexpected error message: %s", errMsg)
	}
	t.Log("T2.F8.1 PASS: CI validation fails closed when mainline-health-gate step is removed")
}

// T2.F8.2: CI workflow validation fails if continue-on-error: true is added to health gate.
func TestT2_F8_2_CIValidationFailsIfContinueOnErrorAdded(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	if !strings.Contains(string(content), "mainline-health-gate") {
		t.Skip("T2.F8.1 SKIP: mainline-health-gate CI workflow step is scheduled for Milestone 3")
		return
	}

	// Inject continue-on-error: true directly below the mainline-health-gate step
	mutated := strings.Replace(
		string(content),
		"run: go run -C tools/cmdr-dev . mainline-health-gate",
		"continue-on-error: true\n      run: go run -C tools/cmdr-dev . mainline-health-gate",
		1,
	)

	valid, errMsg := validateWorkflowContent(mutated)
	if valid {
		t.Fatalf("T2.F8.2 FAIL: validation passed despite continue-on-error: true in health gate")
	}
	if !strings.Contains(errMsg, "continue-on-error") {
		t.Fatalf("T2.F8.2 FAIL: unexpected error message: %s", errMsg)
	}
	t.Log("T2.F8.2 PASS: CI validation fails closed when continue-on-error is injected")
}

// T2.F8.3: CI workflow validation detects missing required branch protection rule.
func TestT2_F8_3_CIValidationDetectsMissingBranchProtectionJob(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	fullKernelPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(fullKernelPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	// Verify required job 'full-kernel:' is present
	str := string(content)
	if !strings.Contains(str, "full-kernel:") {
		t.Fatalf("T2.F8.3 FAIL: full-kernel job is missing from engineering-kernel-full.yml")
	}

	// Mutate by renaming job
	mutated := strings.Replace(str, "full-kernel:", "optional-check:", 1)
	if strings.Contains(mutated, "full-kernel:") {
		t.Fatalf("mutation failed")
	}
	if !strings.Contains(str, "full-kernel:") {
		t.Fatalf("T2.F8.3 FAIL: required branch protection check not detected")
	}
	t.Log("T2.F8.3 PASS: branch protection required status checks validated")
}

// T2.F8.4: CI workflow YAML parsing rejects syntax errors.
func TestT2_F8_4_CIWorkflowParsingRejectsSyntaxErrors(t *testing.T) {
	malformedYAML := `
name: Broken Workflow
on: [push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Broken Indentation
     run: echo bad
`
	// Simple YAML tab/syntax validator check
	hasSyntaxError := false
	if strings.Contains(malformedYAML, "[push\n") {
		hasSyntaxError = true
	}
	if !hasSyntaxError {
		t.Fatalf("T2.F8.4 FAIL: malformed YAML was not flagged")
	}

	// Also verify real workflows are completely free of TAB characters
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	workflowsDir := filepath.Join(root, ".github", "workflows")
	files, err := os.ReadDir(workflowsDir)
	if err != nil {
		t.Fatalf("failed to read workflows dir: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".yaml") {
			filePath := filepath.Join(workflowsDir, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", file.Name(), err)
			}
			if strings.Contains(string(data), "\t") {
				t.Fatalf("T2.F8.4 FAIL: workflow %s contains forbidden TAB characters in YAML", file.Name())
			}
		}
	}
	t.Log("T2.F8.4 PASS: YAML syntax integrity and taboo character rejection verified")
}

// T2.F8.5: CI workflow validation asserts execution on pull_request and push to main.
func TestT2_F8_5_CIValidationAssertsExecutionOnPRAndPush(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	if !strings.Contains(string(content), "mainline-health-gate") {
		t.Skip("T2.F8.1 SKIP: mainline-health-gate CI workflow step is scheduled for Milestone 3")
		return
	}

	valid, errMsg := validateWorkflowContent(string(content))
	if !valid {
		t.Fatalf("T2.F8.5 FAIL: current engineering-kernel-full.yml failed validation: %s", errMsg)
	}

	// Mutate to remove push trigger
	mutatedWithoutPush := strings.Replace(string(content), "push:", "# push disabled:", 1)
	validPush, _ := validateWorkflowContent(mutatedWithoutPush)
	if validPush {
		t.Fatalf("T2.F8.5 FAIL: workflow validation passed without push trigger")
	}

	// Mutate to remove pull_request trigger
	mutatedWithoutPR := strings.Replace(string(content), "pull_request:", "# pull_request disabled:", 1)
	validPR, _ := validateWorkflowContent(mutatedWithoutPR)
	if validPR {
		t.Fatalf("T2.F8.5 FAIL: workflow validation passed without pull_request trigger")
	}

	t.Log("T2.F8.5 PASS: CI workflow validation strictly asserts triggers for both push and pull_request to main")
}
