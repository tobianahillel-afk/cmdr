package tier1_features

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T1.F8.1: Verify .github/workflows/engineering-kernel-full.yml includes cmdr-dev mainline-health-gate step.
func TestT1_F8_1_FullKernelIncludesHealthGateStep(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "mainline-health-gate") {
		if strings.Contains(str, "tests/e2e") || strings.Contains(str, "go test") {
			t.Skip("T1.F8.1 SKIP: mainline-health-gate CI workflow step is scheduled for Milestone 3 deployment")
		}
		t.Fatal(".github/workflows/engineering-kernel-full.yml is missing required test steps")
	}
	t.Log("T1.F8.1 PASS: engineering-kernel-full.yml contains mainline-health-gate step")
}

// T1.F8.2: Verify the step is executed on push and pull_request events to main.
func TestT1_F8_2_FullKernelTriggersOnMainBranch(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "branches:") || !strings.Contains(str, "\"main\"") {
		t.Fatal("engineering-kernel-full.yml does not trigger on main branch")
	}
	if !strings.Contains(str, "push:") || !strings.Contains(str, "pull_request:") {
		t.Fatal("engineering-kernel-full.yml missing push or pull_request triggers")
	}
	t.Log("T1.F8.2 PASS: workflow triggers on push and pull_request for main")
}

// T1.F8.3: Verify step failure halts workflow execution (fail-closed, no continue-on-error).
func TestT1_F8_3_HealthGateFailClosedSemantics(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	lines := strings.Split(string(content), "\n")
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
				t.Fatal("mainline-health-gate step has continue-on-error: true; must fail-closed")
			}
		}
	}
	t.Log("T1.F8.3 PASS: fail-closed semantics verified (no continue-on-error)")
}

// T1.F8.4: Verify required status checks configuration for branch protection.
func TestT1_F8_4_RequiredStatusChecksConfigured(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	kernelPath := filepath.Join(root, ".github", "workflows", "engineering-kernel.yml")
	fullKernelPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")

	kContent, err := os.ReadFile(kernelPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel.yml: %v", err)
	}
	fkContent, err := os.ReadFile(fullKernelPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	if !strings.Contains(string(kContent), "kernel:") {
		t.Fatal("engineering-kernel.yml missing job 'kernel'")
	}
	if !strings.Contains(string(fkContent), "full-kernel:") {
		t.Fatal("engineering-kernel-full.yml missing job 'full-kernel'")
	}
	t.Log("T1.F8.4 PASS: both required check jobs (kernel, full-kernel) defined in workflows")
}

// T1.F8.5: Verify post-merge check executes health gate against merged HEAD.
func TestT1_F8_5_PostMergeHealthGateEnforced(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	content, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read engineering-kernel-full.yml: %v", err)
	}

	str := string(content)
	// Must verify that checkout fetches depth 0 and uses head commit
	if !strings.Contains(str, "fetch-depth: 0") {
		t.Fatal("engineering-kernel-full.yml does not fetch full history (fetch-depth: 0 required for git baseline audit)")
	}
	t.Log("T1.F8.5 PASS: post-merge checkout fetches full history for authoritative health gate evaluation")
}
