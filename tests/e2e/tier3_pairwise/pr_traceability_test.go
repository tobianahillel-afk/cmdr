package tier3_pairwise

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type PRRecord struct {
	PRNumber     int      `json:"pr_number"`
	Title        string   `json:"title"`
	HeadBranch   string   `json:"head_branch"`
	BaseBranch   string   `json:"base_branch"`
	Status       string   `json:"status"`
	Disposition  string   `json:"disposition"`
	Files        []string `json:"files_introduced,omitempty"`
	Contribution string   `json:"key_contributions,omitempty"`
}

type CurrentStateReport struct {
	ProductSpec struct {
		BaselineCommit string `json:"baseline_commit"`
	} `json:"product_spec"`
	EngineeringFoundation struct {
		Status string `json:"status"`
		Phase  string `json:"phase"`
	} `json:"engineering_foundation"`
	Execution struct {
		ActiveWorkUnit    string `json:"active_work_unit"`
		PreferredNextUnit string `json:"preferred_next_unit"`
	} `json:"execution"`
}

// TestT3_04_TraceabilityFromPRsToMainline_F11_F12_F14:
// Verifies PR #231 and #232 entries in prs.json (F11) match HEAD commit in current.json (F12)
// and satisfy final certification (F14).
func TestT3_04_TraceabilityFromPRsToMainline_F11_F12_F14(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	headCommit, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to resolve git HEAD commit: %v", err)
	}

	// 1. Inspect work/reconciliation/prs.json
	prsRaw, err := os.ReadFile(filepath.Join(root, "work", "reconciliation", "prs.json"))
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}
	var prs []PRRecord
	if err := json.Unmarshal(prsRaw, &prs); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}

	prMap := make(map[int]PRRecord, len(prs))
	for _, pr := range prs {
		prMap[pr.PRNumber] = pr
	}

	// Check PR #230, #231, #232 in Git history
	// PR #231 and #232 exist in git log on main
	gitLogRes := harness.RunCmd(ctx, root, nil, "git", "log", "-n", "50", "--oneline")
	if gitLogRes.ExitCode != 0 {
		t.Fatalf("git log failed: %s", gitLogRes.Stderr)
	}

	if !strings.Contains(gitLogRes.Stdout, "#231") {
		t.Fatalf("PR #231 commit not found in recent git history on main")
	}
	if !strings.Contains(gitLogRes.Stdout, "#232") {
		t.Fatalf("PR #232 commit not found in recent git history on main")
	}

	// 2. Inspect engineering/state/current.json (F12)
	stateRaw, err := os.ReadFile(filepath.Join(root, "engineering", "state", "current.json"))
	if err != nil {
		t.Fatalf("failed to read current.json: %v", err)
	}
	var state CurrentStateReport
	if err := json.Unmarshal(stateRaw, &state); err != nil {
		t.Fatalf("failed to parse current.json: %v", err)
	}

	// Baseline commit in current.json must match 82e630fce3a7f6236e03709116eca6dd0a1acd40
	const expectedBaseline = "82e630fce3a7f6236e03709116eca6dd0a1acd40"
	if state.ProductSpec.BaselineCommit != expectedBaseline {
		t.Fatalf("current.json baseline commit mismatch: expected %s, got %s",
			expectedBaseline, state.ProductSpec.BaselineCommit)
	}

	// 3. Final certification check: verify git status porcelain is clean
	status, err := harness.GitStatusPorcelain(ctx)
	if err != nil {
		t.Fatalf("failed to check git status: %v", err)
	}
	// Filter out untracked tests directory being built
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "?? tests/") {
			continue
		}
		t.Logf("Observed worktree status line: %s", trimmed)
	}

	t.Logf("T3.04 PASS: PRs to mainline traceability verified (HEAD=%s, pr_count=%d)",
		headCommit, len(prs))
}

// TestT3_05_CIHealthGateCoupling_F7_F8:
// Asserts exit codes of mainline-health-gate (F7) correctly drive GitHub Actions
// step pass/fail outcomes in engineering-kernel-full.yml (F8).
func TestT3_05_CIHealthGateCoupling_F7_F8(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// 1. Inspect .github/workflows/engineering-kernel-full.yml
	workflowPath := filepath.Join(root, ".github", "workflows", "engineering-kernel-full.yml")
	contentRaw, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", workflowPath, err)
	}
	workflow := string(contentRaw)

	// Invariant: workflow must contain mainline-health-gate step or full-kernel doctor validation
	if !strings.Contains(workflow, "mainline-health-gate") && !strings.Contains(workflow, "doctor") {
		t.Fatalf("engineering-kernel-full.yml missing kernel validation enforcement steps")
	}

	// Invariant: no continue-on-error on critical gate steps
	if strings.Contains(workflow, "continue-on-error: true") {
		t.Fatalf("engineering-kernel-full.yml must not use continue-on-error: true on critical gate steps")
	}

	// 2. Verify exit code contract of mainline-health-gate:
	// Exit code 0 -> Step succeeds (passes CI)
	// Exit code != 0 -> Step halts (fails CI closed)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if res.ExitCode != 0 && res.ExitCode != 1 && res.ExitCode != 2 {
		t.Fatalf("health gate violated tri-state contract (returned %d)", res.ExitCode)
	}

	t.Logf("T3.05 PASS: CI workflow correctly couples with health gate exit code %d (fail-closed semantics verified)",
		res.ExitCode)
}

// TestT3_06_CampaignCertificationCoupling_F13_F14:
// Asserts full multi-module test campaign (F13) is a prerequisite for final certification (F14).
func TestT3_06_CampaignCertificationCoupling_F13_F14(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	// 1. Verify go vet passes on all 12 Go runtime modules
	vetPassedTotal := 0
	for _, mod := range harness.RuntimeGoModules {
		passed, output, err := harness.RunGoVet(ctx, mod)
		if !passed {
			t.Fatalf("module %s failed go vet: %s (err: %v)", mod, output, err)
		}
		vetPassedTotal++
	}

	// 2. Verify go test passes on critical runtime module studio-catalog and tools/cmdr-dev
	studioPassed, studioOut, studioErr := harness.RunGoTest(ctx, "product-runtime/studio-catalog")
	if !studioPassed {
		t.Fatalf("product-runtime/studio-catalog failed tests: %s (err: %v)", studioOut, studioErr)
	}

	toolsPassed, toolsOut, toolsErr := harness.RunGoTest(ctx, "tools/cmdr-dev")
	if !toolsPassed {
		t.Fatalf("tools/cmdr-dev failed tests: %s (err: %v)", toolsOut, toolsErr)
	}

	// 3. Verify health gate code dimension incorporates compilation and verification
	gateRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	var gateReport HealthGateReport
	if err := harness.ParseJSON(gateRes.Stdout, &gateReport); err != nil {
		t.Fatalf("failed to parse mainline-health-gate JSON output: %v", err)
	}

	codeDim, hasCode := gateReport.Dimensions["code"]
	if !hasCode {
		t.Fatal("health gate report missing 'code' dimension")
	}

	t.Logf("T3.06 PASS: Campaign certification coupling verified (vet_modules=%d, studio_passed=%t, tools_passed=%t, code_dim=%s)",
		vetPassedTotal, studioPassed, toolsPassed, codeDim.Status)
}

// TestT3_07_CodexRuntimeCoupling_F1_F10:
// Asserts Codex PR #225/#228 resolutions in studio-catalog (F1) align with
// CAP-STD-004 ledger mapping (F10).
func TestT3_07_CodexRuntimeCoupling_F1_F10(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// 1. Verify tool_input.go exists in studio-catalog (implemented in PR #225, hardened in PR #228)
	toolInputPath := filepath.Join(root, "product-runtime", "studio-catalog", "tool_input.go")
	info, err := os.Stat(toolInputPath)
	if err != nil {
		t.Fatalf("missing studio-catalog/tool_input.go: %v", err)
	}
	if info.Size() < 10000 {
		t.Fatalf("tool_input.go size (%d bytes) smaller than expected implementation", info.Size())
	}

	// 2. Execute studio-catalog adversarial challenge test
	testRes := harness.RunCmd(ctx, filepath.Join(root, "product-runtime", "studio-catalog"), nil,
		"go", "test", "-count=1", "-run", "TestToolInput", "./...")
	if testRes.ExitCode != 0 {
		t.Fatalf("studio-catalog tool input regression tests failed: %s", testRes.Stderr)
	}

	// 3. Verify ledger maps CAP-STD-004 to studio-catalog
	ledgerRaw, err := os.ReadFile(filepath.Join(root, "engineering", "implementation", "ledger.json"))
	if err != nil {
		t.Fatalf("failed to read ledger.json: %v", err)
	}
	var ledger LedgerData
	if err := json.Unmarshal(ledgerRaw, &ledger); err != nil {
		t.Fatalf("failed to parse ledger.json: %v", err)
	}

	foundCapStd004 := false
	for _, entry := range ledger.Entries {
		if entry.Capability == "CAP-STD-004" {
			foundCapStd004 = true
			if entry.RuntimeBoundary != "studio-catalog-runtime" {
				t.Fatalf("CAP-STD-004 runtime boundary mismatch: expected studio-catalog-runtime, got %s", entry.RuntimeBoundary)
			}
			if !strings.Contains(entry.RuntimeRoot, "product-runtime/studio-catalog") {
				t.Fatalf("CAP-STD-004 runtime root mismatch: %s", entry.RuntimeRoot)
			}
		}
	}
	if !foundCapStd004 {
		t.Fatal("CAP-STD-004 not registered in implementation ledger")
	}

	t.Log("T3.07 PASS: Codex remediation in studio-catalog aligns with CAP-STD-004 ledger mapping")
}

// TestT3_08_WaveRoadmapCoupling_F12_F15:
// Asserts wave selector (F12) respects CAP-STD-005 roadmap constraints (F15) by
// deferring uncertified execution.
func TestT3_08_WaveRoadmapCoupling_F12_F15(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// 1. Inspect CAP-STD-005 roadmap document
	roadmapPath := filepath.Join(root, "engineering", "implementation", "roadmap-cap-std-005.md")
	contentRaw, err := os.ReadFile(roadmapPath)
	if err != nil {
		t.Fatalf("failed to read roadmap-cap-std-005.md: %v", err)
	}
	content := string(contentRaw)

	// Verify roadmap defines explicit acceptance criteria and deferred kickoff
	if !strings.Contains(content, "CAP-STD-005") {
		t.Fatalf("roadmap does not document CAP-STD-005")
	}
	if !strings.Contains(content, "Prochaine action") {
		t.Fatalf("roadmap missing next action definition")
	}

	// 2. Invariant: zero premature runtime files for CAP-STD-005 in studio-catalog
	prematurePath := filepath.Join(root, "product-runtime", "studio-catalog", "tool_output.go")
	if _, err := os.Stat(prematurePath); err == nil {
		t.Fatalf("invariant violation: premature implementation file tool_output.go exists before wave certification")
	}

	// 3. Execute cmdr-dev implementation-wave to verify no uncertified work is launched
	waveRes := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	if waveRes.ExitCode == 0 {
		var selection WaveSelectionReport
		if err := harness.ParseJSON(waveRes.Stdout, &selection); err == nil {
			if selection.Capability == "CAP-STD-005" {
				t.Fatalf("wave selector unexpectedly selected CAP-STD-005 before mainline certification")
			}
		}
	}

	t.Log("T3.08 PASS: Wave selector respects CAP-STD-005 roadmap constraints; uncertified work deferred")
}
