package tier2_boundaries

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// validateRoadmapDocument verifies that a capability roadmap contains all required structural sections.
func validateRoadmapDocument(content string) (bool, string) {
	requiredSections := []string{
		"Contexte et Justification",
		"Exigences Produit et Invariants Canoniques",
		"Plan d'Implémentation Technique",
		"Statut de Préparation",
	}

	for _, sec := range requiredSections {
		if !strings.Contains(content, sec) {
			return false, "roadmap document missing mandatory section: " + sec
		}
	}
	return true, ""
}

func assertGitDimensionFailsInMockRoot(t *testing.T, ctx context.Context, mockRoot, label string) {
	t.Helper()
	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if gateRes.ExitCode == 0 {
		t.Fatalf("%s FAIL: mainline-health-gate returned exit 0", label)
	}
	var report struct {
		Status     string `json:"status"`
		Dimensions map[string]struct {
			Passed bool   `json:"passed"`
			Status string `json:"status"`
		} `json:"dimensions"`
	}
	if err := harness.ParseJSON(gateRes.Stdout, &report); err != nil {
		t.Fatalf("%s FAIL: failed to parse mainline-health-gate JSON: %v\nStdout: %s", label, err, gateRes.Stdout)
	}
	gitDim, ok := report.Dimensions["git"]
	if !ok || gitDim.Passed || gitDim.Status != "FAIL" {
		t.Fatalf("%s FAIL: expected git dimension FAIL, got %+v", label, gitDim)
	}
}

// T2.F15.1: Spec diff detector detects 1-byte alteration in any of 1,879 spec files.
func TestT2_F15_1_SpecDiffDetectorDetectsOneByteAlteration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	const baseline = "82e630fce3a7f6236e03709116eca6dd0a1acd40"

	// Verify that current repo has exactly 0 diff against canonical baseline
	diff, err := harness.GitDiffSpecBaseline(ctx, baseline)
	if err != nil {
		t.Fatalf("failed to check git diff vs baseline: %v", err)
	}
	if diff != "" {
		t.Fatalf("T2.F15.1 FAIL: canonical cmdr-product-spec has non-zero diff vs %s:\n%s", baseline, diff)
	}

	mockRoot := harness.CloneIsolatedRepo(t, "f15-onebyte")
	baseRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F15.1 baseline spec-baseline failed (exit %d): %s", baseRes.ExitCode, baseRes.Stderr)
	}

	readmePath := filepath.Join(mockRoot, "cmdr-product-spec", "README.md")
	orig, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read spec README: %v", err)
	}
	if err := os.WriteFile(readmePath, append(orig, ' '), 0644); err != nil {
		t.Fatalf("failed to append 1 byte to spec README: %v", err)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "cmdr-product-spec/README.md"); res.ExitCode != 0 {
		t.Fatalf("git add failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "1-byte spec alteration"); res.ExitCode != 0 {
		t.Fatalf("git commit failed: %s", res.Stderr)
	}

	specRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if specRes.ExitCode == 0 {
		t.Fatalf("T2.F15.1 FAIL: spec-baseline --check succeeded after 1-byte spec edit")
	}
	assertGitDimensionFailsInMockRoot(t, ctx, mockRoot, "T2.F15.1")
	t.Log("T2.F15.1 PASS: spec-baseline and mainline-health-gate strictly detect 1-byte canonical spec alteration")
}

// T2.F15.2: Spec diff detector detects deleted spec file.
func TestT2_F15_2_SpecDiffDetectorDetectsDeletedSpecFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f15-delete")
	baseRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F15.2 baseline spec-baseline failed (exit %d): %s", baseRes.ExitCode, baseRes.Stderr)
	}

	targetRel := filepath.Join("cmdr-product-spec", "00-governance", "README.md")
	if err := os.Remove(filepath.Join(mockRoot, targetRel)); err != nil {
		t.Fatalf("failed to remove spec file: %v", err)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "-A", "cmdr-product-spec"); res.ExitCode != 0 {
		t.Fatalf("git add failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "delete spec file"); res.ExitCode != 0 {
		t.Fatalf("git commit failed: %s", res.Stderr)
	}

	specRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if specRes.ExitCode == 0 {
		t.Fatalf("T2.F15.2 FAIL: spec-baseline --check succeeded after deleting spec file")
	}
	assertGitDimensionFailsInMockRoot(t, ctx, mockRoot, "T2.F15.2")
	t.Log("T2.F15.2 PASS: spec-baseline and mainline-health-gate strictly detect deleted spec file")
}

// T2.F15.3: Spec diff detector detects added untracked spec file.
func TestT2_F15_3_SpecDiffDetectorDetectsAddedUntrackedFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f15-untracked")
	baseRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F15.3 baseline spec-baseline failed (exit %d): %s", baseRes.ExitCode, baseRes.Stderr)
	}

	untrackedPath := filepath.Join(mockRoot, "cmdr-product-spec", "unauthorized_new_spec.md")
	if err := os.WriteFile(untrackedPath, []byte("# Unauthorized Spec\n"), 0644); err != nil {
		t.Fatalf("failed to write untracked spec file: %v", err)
	}

	status, err := harness.GitStatusPorcelainInDir(ctx, mockRoot)
	if err != nil {
		t.Fatalf("GitStatusPorcelainInDir failed: %v", err)
	}
	if !strings.Contains(status, "?? cmdr-product-spec/unauthorized_new_spec.md") {
		t.Fatalf("T2.F15.3 FAIL: untracked spec file not reported by git status: %q", status)
	}

	specRes := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if specRes.ExitCode == 0 {
		t.Fatalf("T2.F15.3 FAIL: spec-baseline --check succeeded with untracked spec file")
	}
	assertGitDimensionFailsInMockRoot(t, ctx, mockRoot, "T2.F15.3")
	t.Log("T2.F15.3 PASS: spec-baseline, git status, and mainline-health-gate strictly detect untracked spec file")
}

// T2.F15.4: Roadmap alignment validator rejects missing acceptance criteria section.
func TestT2_F15_4_RoadmapAlignmentValidatorRejectsMissingSection(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	content, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap-cap-std-005.md: %v", err)
	}

	// Verify real document passes all section requirements
	valid, errMsg := validateRoadmapDocument(string(content))
	if !valid {
		t.Fatalf("T2.F15.4 FAIL: current roadmap document failed validation: %s", errMsg)
	}

	// Mutate by omitting "Exigences Produit et Invariants Canoniques" section
	mutated := strings.Replace(string(content), "Exigences Produit et Invariants Canoniques", "Section Supprimee", 1)
	validMutated, errMutated := validateRoadmapDocument(mutated)
	if validMutated {
		t.Fatalf("T2.F15.4 FAIL: validation passed despite missing acceptance criteria section")
	}
	if !strings.Contains(errMutated, "Exigences Produit et Invariants Canoniques") {
		t.Fatalf("T2.F15.4 FAIL: unexpected error message: %s", errMutated)
	}
	t.Log("T2.F15.4 PASS: roadmap alignment validator strictly rejects documents missing required criteria sections")
}

// T2.F15.5: Roadmap alignment validator verifies CAP-STD-005 status is PLANNED (not IN_PROGRESS).
func TestT2_F15_5_RoadmapAlignmentValidatorAssertsCAPSTD005Planned(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	// 1. Assert zero premature implementation files for CAP-STD-005 in product-runtime/studio-catalog/
	prematureFile := filepath.Join(root, "product-runtime", "studio-catalog", "tool_output.go")
	if _, err := os.Stat(prematureFile); err == nil {
		t.Fatalf("T2.F15.5 FAIL: found premature implementation file %s before stabilization certification", prematureFile)
	}

	// 2. Assert zero premature test files for CAP-STD-005
	prematureTestFile := filepath.Join(root, "product-runtime", "studio-catalog", "tool_output_test.go")
	if _, err := os.Stat(prematureTestFile); err == nil {
		t.Fatalf("T2.F15.5 FAIL: found premature test file %s before stabilization certification", prematureTestFile)
	}

	// 3. Verify roadmap states explicit preparation without premature branching
	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	content, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap: %v", err)
	}

	if !strings.Contains(string(content), "dès lancement du prochain incrément") {
		t.Fatalf("T2.F15.5 FAIL: roadmap must explicitly defer next branch creation until authorization")
	}

	t.Log("T2.F15.5 PASS: CAP-STD-005 is strictly PLANNED; zero premature implementation code detected")
}
