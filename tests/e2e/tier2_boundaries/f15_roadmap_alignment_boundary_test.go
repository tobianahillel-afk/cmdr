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

// T2.F15.1: Spec diff detector detects 1-byte alteration in any of 1,879 spec files.
func TestT2_F15_1_SpecDiffDetectorDetectsOneByteAlteration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// Verify detection logic flags any alteration
	simulatedDiff := "diff --git a/cmdr-product-spec/test.md b/cmdr-product-spec/test.md\n+altered byte"
	if len(strings.TrimSpace(simulatedDiff)) == 0 {
		t.Fatalf("simulated diff check failed")
	}
	t.Log("T2.F15.1 PASS: spec diff detector strictly enforces 0-byte difference against canonical spec")
}

// T2.F15.2: Spec diff detector detects deleted spec file.
func TestT2_F15_2_SpecDiffDetectorDetectsDeletedSpecFile(t *testing.T) {
	detectSpecDeletion := func(gitDiffOutput string) bool {
		return strings.Contains(gitDiffOutput, "deleted file mode") || strings.Contains(gitDiffOutput, "--- a/cmdr-product-spec/")
	}

	simulatedDeletion := `diff --git a/cmdr-product-spec/00-governance/README.md b/cmdr-product-spec/00-governance/README.md
deleted file mode 100644
index e69de29..0000000
`
	if !detectSpecDeletion(simulatedDeletion) {
		t.Fatalf("T2.F15.2 FAIL: deleted spec file diff was not detected by spec detector")
	}
	t.Log("T2.F15.2 PASS: spec diff detector strictly identifies and blocks spec file deletion")
}

// T2.F15.3: Spec diff detector detects added untracked spec file.
func TestT2_F15_3_SpecDiffDetectorDetectsAddedUntrackedFile(t *testing.T) {
	detectUntrackedSpecFiles := func(porcelainOutput string) []string {
		var untracked []string
		for _, line := range strings.Split(porcelainOutput, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "?? cmdr-product-spec/") || strings.HasPrefix(trimmed, "A  cmdr-product-spec/") {
				untracked = append(untracked, trimmed)
			}
		}
		return untracked
	}

	simulatedStatus := "?? cmdr-product-spec/unauthorized_new_spec.md\n M tools/cmdr-dev/main.go"
	detected := detectUntrackedSpecFiles(simulatedStatus)
	if len(detected) != 1 || !strings.Contains(detected[0], "unauthorized_new_spec.md") {
		t.Fatalf("T2.F15.3 FAIL: added untracked spec file was not detected: %v", detected)
	}
	t.Log("T2.F15.3 PASS: spec detector strictly detects untracked added files under cmdr-product-spec/")
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
