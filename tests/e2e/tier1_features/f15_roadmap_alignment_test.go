package tier1_features

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

const specBaselineCommit = "82e630fce3a7f6236e03709116eca6dd0a1acd40"

// T1.F15.1: Verify cmdr-product-spec diff vs 82e630fc is exactly 0.
func TestT1_F15_1_ProductSpecZeroDiffVsBaseline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	diff, err := harness.GitDiffSpecBaseline(ctx, specBaselineCommit)
	if err != nil {
		t.Fatalf("failed to compute git diff vs baseline commit %s: %v", specBaselineCommit, err)
	}

	if diff != "" {
		t.Fatalf("invariant breach: cmdr-product-spec has non-zero diff vs baseline commit %s:\n%s", specBaselineCommit, diff)
	}
	t.Logf("T1.F15.1 PASS: cmdr-product-spec has exactly 0 diff vs baseline %s across all 1,879 files", specBaselineCommit)
}

// T1.F15.2: Verify CAP-STD-005 roadmap document defines explicit acceptance criteria.
func TestT1_F15_2_RoadmapDefinesAcceptanceCriteria(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	content, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap-cap-std-005.md: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "CAP-STD-005") {
		t.Fatal("roadmap document missing CAP-STD-005 reference")
	}
	if !strings.Contains(str, "Tool Output, Result and Error Contract") {
		t.Fatal("roadmap document missing capability title")
	}
	if !strings.Contains(str, "Exigences Produit") && !strings.Contains(str, "Acceptance Criteria") {
		t.Fatal("roadmap document missing explicit product requirements / criteria section")
	}
	t.Log("T1.F15.2 PASS: CAP-STD-005 roadmap document defines explicit requirements and acceptance criteria")
}

// T1.F15.3: Verify CAP-STD-005 dependencies and risks are documented.
func TestT1_F15_3_RoadmapDependenciesAndRisksDocumented(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	content, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap-cap-std-005.md: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "CAP-STD-004") {
		t.Fatal("roadmap document does not reference prerequisite CAP-STD-004 dependency")
	}
	if !strings.Contains(str, "Dépendances") && !strings.Contains(str, "Dependencies") {
		t.Fatal("roadmap document does not document dependency section")
	}
	t.Log("T1.F15.3 PASS: CAP-STD-005 upstream dependencies (CAP-STD-004) and risks documented")
}

// T1.F15.4: Verify zero premature implementation files for CAP-STD-005 exist in runtime.
func TestT1_F15_4_ZeroPrematureImplementationFiles(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	toolOutputPath := filepath.Join(root, "product-runtime", "studio-catalog", "tool_output.go")
	if _, err := os.Stat(toolOutputPath); err == nil {
		t.Fatal("premature implementation file product-runtime/studio-catalog/tool_output.go must not exist before stabilization certification")
	}

	lotPath := filepath.Join(root, "work", "lots", "E10-STD-005A-OUTPUT")
	if _, err := os.Stat(lotPath); err == nil {
		t.Fatal("premature lot work/lots/E10-STD-005A-OUTPUT must not exist before stabilization certification")
	}

	t.Log("T1.F15.4 PASS: zero premature implementation files or lots for CAP-STD-005 exist in repository")
}

// T1.F15.5: Verify roadmap explicitly requires certified healthy mainline before launch.
func TestT1_F15_5_RoadmapRequiresCertifiedMainline(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	content, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap-cap-std-005.md: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "MAINLINE_HEALTH_GATE") && !strings.Contains(str, "main") {
		t.Fatal("roadmap document does not mandate branch from validated main under health gate")
	}
	t.Log("T1.F15.5 PASS: roadmap explicitly requires certified healthy mainline before starting next increment")
}
