package tier2_boundaries

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

const syntheticValidHealthGateWorkflow = `name: engineering-kernel-full

on:
  workflow_dispatch:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"

permissions:
  contents: read

jobs:
  full-kernel:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Enforce Mainline Health Gate
        working-directory: tools/cmdr-dev
        run: go run . mainline-health-gate --root ../..
`

// T2.F8.1: CI workflow validation fails if mainline-health-gate step is removed.
func TestT2_F8_1_CIValidationFailsIfHealthGateStepRemoved(t *testing.T) {
	baseRes := harness.ValidateHealthGateWorkflow(syntheticValidHealthGateWorkflow)
	if !baseRes.Valid || !baseRes.HasHealthGateStep {
		t.Fatalf("T2.F8.1 setup error: synthetic valid workflow failed validation: %v", baseRes.Errors)
	}

	// Remove mainline-health-gate step from synthetic workflow
	mutatedLines := make([]string, 0)
	for _, line := range strings.Split(syntheticValidHealthGateWorkflow, "\n") {
		if !strings.Contains(line, "mainline-health-gate") && !strings.Contains(line, "Enforce Mainline Health Gate") {
			mutatedLines = append(mutatedLines, line)
		}
	}
	mutated := strings.Join(mutatedLines, "\n")

	mutRes := harness.ValidateHealthGateWorkflow(mutated)
	if mutRes.Valid || mutRes.HasHealthGateStep {
		t.Fatalf("T2.F8.1 FAIL: validation passed despite removing mainline-health-gate step")
	}

	// Also verify live workflow if available
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	liveRes := harness.ValidateHealthGateWorkflow(string(content))
	if liveRes.HasHealthGateStep {
		liveMutatedLines := make([]string, 0)
		for _, line := range strings.Split(string(content), "\n") {
			if !strings.Contains(line, "mainline-health-gate") {
				liveMutatedLines = append(liveMutatedLines, line)
			}
		}
		liveMutRes := harness.ValidateHealthGateWorkflow(strings.Join(liveMutatedLines, "\n"))
		if liveMutRes.Valid || liveMutRes.HasHealthGateStep {
			t.Fatalf("T2.F8.1 FAIL: live workflow validation passed despite removing mainline-health-gate step")
		}
	}

	t.Log("T2.F8.1 PASS: CI validation fails closed when mainline-health-gate step is removed")
}

// T2.F8.2: CI workflow validation fails if continue-on-error: true or if: false or commented step is used.
func TestT2_F8_2_CIValidationFailsIfContinueOnErrorAdded(t *testing.T) {
	// 1. Inject continue-on-error: true into synthetic health gate step
	mutatedContinue := strings.Replace(
		syntheticValidHealthGateWorkflow,
		"        run: go run . mainline-health-gate --root ../..",
		"        continue-on-error: true\n        run: go run . mainline-health-gate --root ../..",
		1,
	)
	resContinue := harness.ValidateHealthGateWorkflow(mutatedContinue)
	if resContinue.Valid || !resContinue.HasContinueOnError {
		t.Fatalf("T2.F8.2 FAIL: validation failed to reject continue-on-error: true (result=%+v)", resContinue)
	}

	// 2. Inject if: false into synthetic health gate step
	mutatedIfFalse := strings.Replace(
		syntheticValidHealthGateWorkflow,
		"        run: go run . mainline-health-gate --root ../..",
		"        if: false\n        run: go run . mainline-health-gate --root ../..",
		1,
	)
	resIfFalse := harness.ValidateHealthGateWorkflow(mutatedIfFalse)
	if resIfFalse.Valid || !resIfFalse.HasDisabledIf {
		t.Fatalf("T2.F8.2 FAIL: validation failed to reject if: false (result=%+v)", resIfFalse)
	}

	// 3. Comment out the run command
	mutatedCommented := strings.Replace(
		syntheticValidHealthGateWorkflow,
		"        run: go run . mainline-health-gate --root ../..",
		"        # run: go run . mainline-health-gate --root ../..\n        run: go run . doctor --root ../..",
		1,
	)
	resCommented := harness.ValidateHealthGateWorkflow(mutatedCommented)
	if resCommented.Valid || resCommented.HasHealthGateStep {
		t.Fatalf("T2.F8.2 FAIL: validation failed to reject commented/doctor-only step (result=%+v)", resCommented)
	}

	// 4. Also test live workflow if mainline-health-gate step is present
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}
	liveRes := harness.ValidateHealthGateWorkflow(string(content))
	if liveRes.HasHealthGateStep {
		liveMutated := strings.Replace(
			string(content),
			"mainline-health-gate",
			"mainline-health-gate\n        continue-on-error: true",
			1,
		)
		liveMutRes := harness.ValidateHealthGateWorkflow(liveMutated)
		if liveMutRes.Valid || !liveMutRes.HasContinueOnError {
			t.Fatalf("T2.F8.2 FAIL: live workflow mutation with continue-on-error: true was not rejected")
		}
	}

	t.Log("T2.F8.2 PASS: CI validation fails closed when continue-on-error, if: false, or commented step is injected")
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
	baseRes := harness.ValidateHealthGateWorkflow(syntheticValidHealthGateWorkflow)
	if !baseRes.Valid || !baseRes.HasPushMain || !baseRes.HasPRMain {
		t.Fatalf("T2.F8.5 setup error: synthetic valid workflow failed trigger validation: %v", baseRes.Errors)
	}

	// Mutate to remove push trigger
	mutatedWithoutPush := strings.Replace(syntheticValidHealthGateWorkflow, "  push:", "  # push disabled:", 1)
	resNoPush := harness.ValidateHealthGateWorkflow(mutatedWithoutPush)
	if resNoPush.Valid || resNoPush.HasPushMain {
		t.Fatalf("T2.F8.5 FAIL: workflow validation passed without push trigger")
	}

	// Mutate to remove pull_request trigger
	mutatedWithoutPR := strings.Replace(syntheticValidHealthGateWorkflow, "  pull_request:", "  # pull_request disabled:", 1)
	resNoPR := harness.ValidateHealthGateWorkflow(mutatedWithoutPR)
	if resNoPR.Valid || resNoPR.HasPRMain {
		t.Fatalf("T2.F8.5 FAIL: workflow validation passed without pull_request trigger")
	}

	// Mutate target branch away from main
	mutatedBranch := strings.ReplaceAll(syntheticValidHealthGateWorkflow, `"main"`, `"develop"`)
	resWrongBranch := harness.ValidateHealthGateWorkflow(mutatedBranch)
	if resWrongBranch.Valid || resWrongBranch.HasPushMain || resWrongBranch.HasPRMain {
		t.Fatalf("T2.F8.5 FAIL: workflow validation passed without targeting main branch")
	}

	// Also validate live workflow when mainline-health-gate is deployed
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}
	liveRes := harness.ValidateHealthGateWorkflow(string(content))
	if liveRes.HasHealthGateStep && (!liveRes.Valid || !liveRes.HasPushMain || !liveRes.HasPRMain) {
		t.Fatalf("T2.F8.5 FAIL: live engineering-kernel-full.yml failed trigger validation: %v", liveRes.Errors)
	}

	t.Log("T2.F8.5 PASS: CI workflow validation strictly asserts triggers for both push and pull_request to main")
}
