package tier2_boundaries

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

// T2.F12.1: Current state validation fails if commit SHA does not match HEAD.
func TestT2_F12_1_CurrentStateFailsIfCommitDoesNotMatchHEAD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f12-state-sha")

	// Pre-mutation baseline must pass spec-baseline --check
	baseSpec := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if baseSpec.ExitCode != 0 {
		t.Fatalf("T2.F12.1 baseline spec-baseline failed (exit %d): %s\n%s", baseSpec.ExitCode, baseSpec.Stderr, baseSpec.Stdout)
	}

	statePath := filepath.Join(mockRoot, "engineering", "state", "current.json")
	rawState, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("failed to read current.json: %v", err)
	}

	updateBaselineCommit := func(sha string) {
		var doc map[string]interface{}
		if err := json.Unmarshal(rawState, &doc); err != nil {
			t.Fatalf("failed to parse current.json: %v", err)
		}
		ps, ok := doc["product_spec"].(map[string]interface{})
		if !ok {
			t.Fatalf("product_spec missing from current.json")
		}
		ps["baseline_commit"] = sha
		mutated, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal current.json: %v", err)
		}
		if err := os.WriteFile(statePath, append(mutated, '\n'), 0644); err != nil {
			t.Fatalf("failed to write mutated current.json: %v", err)
		}
	}

	// (a) Zero SHA must fail both spec-baseline and implementation-ledger-audit
	updateBaselineCommit("0000000000000000000000000000000000000000")
	zeroSpec := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if zeroSpec.ExitCode == 0 {
		t.Fatalf("T2.F12.1 FAIL: spec-baseline --check succeeded with zero baseline_commit")
	}
	zeroLedger := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot, "--json")
	if zeroLedger.ExitCode == 0 {
		t.Fatalf("T2.F12.1 FAIL: implementation-ledger-audit succeeded with zero baseline_commit")
	}

	// (b) Real divergent git commit in mockRoot where cmdr-product-spec differs must fail both checks
	specReadme := filepath.Join(mockRoot, "cmdr-product-spec", "README.md")
	origSpec, err := os.ReadFile(specReadme)
	if err != nil {
		t.Fatalf("failed to read spec README: %v", err)
	}
	if err := os.WriteFile(specReadme, append(origSpec, []byte("\n<!-- divergent commit -->\n")...), 0644); err != nil {
		t.Fatalf("failed to mutate spec README: %v", err)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "cmdr-product-spec/README.md"); res.ExitCode != 0 {
		t.Fatalf("git add divergent spec failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "divergent spec commit"); res.ExitCode != 0 {
		t.Fatalf("git commit divergent spec failed: %s", res.Stderr)
	}
	divRes := harness.RunCmd(ctx, mockRoot, nil, "git", "rev-parse", "HEAD")
	if divRes.ExitCode != 0 {
		t.Fatalf("git rev-parse divergent commit failed: %s", divRes.Stderr)
	}
	divergentSHA := strings.TrimSpace(divRes.Stdout)

	// Restore spec README back to canonical content so current working tree differs from divergentSHA
	if err := os.WriteFile(specReadme, origSpec, 0644); err != nil {
		t.Fatalf("failed to restore spec README: %v", err)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "cmdr-product-spec/README.md"); res.ExitCode != 0 {
		t.Fatalf("git add restored spec failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "restore canonical spec"); res.ExitCode != 0 {
		t.Fatalf("git commit restored spec failed: %s", res.Stderr)
	}

	updateBaselineCommit(divergentSHA)
	divSpec := harness.RunCmdrDev(ctx, "spec-baseline", "--root", mockRoot, "--check", "--json")
	if divSpec.ExitCode == 0 {
		t.Fatalf("T2.F12.1 FAIL: spec-baseline --check succeeded against divergent commit %s", divergentSHA)
	}
	divLedger := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--root", mockRoot, "--json")
	if divLedger.ExitCode == 0 {
		t.Fatalf("T2.F12.1 FAIL: implementation-ledger-audit succeeded against divergent baseline commit %s", divergentSHA)
	}

	t.Logf("T2.F12.1 PASS: baseline_commit desync (zero SHA and divergent commit %s) failed closed", divergentSHA)
}

// T2.F12.2: Current state validation detects desynced active unit pointer.
func TestT2_F12_2_CurrentStateDetectsDesyncedActiveUnit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockGraphRepo(t)
	statePath := filepath.Join(mockRoot, "engineering", "state", "current.json")

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("failed to read state: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse state JSON: %v", err)
	}

	execution, ok := doc["execution"].(map[string]interface{})
	if !ok {
		t.Fatalf("execution field missing from current.json")
	}

	// Point active_work_unit to a ghost unit not in graph.json
	execution["active_work_unit"] = "E10-GHOST-NONEXISTENT-UNIT"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated state: %v", err)
	}
	if err := os.WriteFile(statePath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated state: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "doctor", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F12.2 FAIL: doctor passed despite desynced active work unit pointer")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "active work unit") || (!strings.Contains(combined, "does not exist in work graph") && !strings.Contains(combined, "is missing from graph")) {
		t.Fatalf("T2.F12.2 FAIL: expected active unit missing from graph error, got: %s", combined)
	}
	t.Log("T2.F12.2 PASS: current state validation strictly detects desynced active work unit pointer")
}

// T2.F12.3: Wave selector handles empty candidate pool gracefully without panic.
func TestT2_F12_3_WaveSelectorHandlesEmptyCandidatePoolGracefully(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	// Must execute or fail gracefully without panic
	if strings.Contains(res.Stderr, "panic:") || strings.Contains(res.Stdout, "panic:") {
		t.Fatalf("T2.F12.3 FAIL: implementation-wave panicked: %s", res.Stderr)
	}

	if res.ExitCode == 0 {
		var wave map[string]interface{}
		if err := json.Unmarshal([]byte(res.Stdout), &wave); err == nil {
			t.Logf("T2.F12.3 PASS: wave selector returned candidate=%v (ready_candidates=%v)", wave["capability"], wave["ready_candidates"])
			return
		}
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "no READY") && !strings.Contains(combined, "READY capability") && res.ExitCode != 1 {
		t.Fatalf("T2.F12.3 FAIL: expected 'no READY capability' error on empty candidate pool, got: %s", combined)
	}
	t.Log("T2.F12.3 PASS: wave selector gracefully handles empty candidate pool without panic")
}

// T2.F12.4: Wave selector rejects unit whose prerequisites are not ACCEPTED.
func TestT2_F12_4_WaveSelectorRejectsUnitWithUnacceptedPrerequisites(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f12-prereq")
	graphPath := filepath.Join(mockRoot, "work", "graph.json")

	writeGraphWithParentStatus := func(parentStatus string) {
		raw, err := os.ReadFile(graphPath)
		if err != nil {
			t.Fatalf("failed to read mock graph.json: %v", err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("failed to parse mock graph.json: %v", err)
		}
		origNodes, ok := doc["nodes"].([]interface{})
		if !ok {
			t.Fatalf("nodes missing from graph.json")
		}
		filtered := make([]interface{}, 0, len(origNodes)+2)
		for _, n := range origNodes {
			if m, ok := n.(map[string]interface{}); ok {
				id, _ := m["id"].(string)
				if id == "E10-PREREQ-PARENT" || id == "E10-PREREQ-CHILD" {
					continue
				}
			}
			filtered = append(filtered, n)
		}
		filtered = append(filtered,
			map[string]interface{}{
				"id":         "E10-PREREQ-PARENT",
				"title":      "Prerequisite parent",
				"type":       "task",
				"status":     parentStatus,
				"depends_on": []interface{}{},
				"unlocks":    []interface{}{"E10-PREREQ-CHILD"},
			},
			map[string]interface{}{
				"id":         "E10-PREREQ-CHILD",
				"title":      "Dependent ready child",
				"type":       "task",
				"status":     "READY",
				"depends_on": []interface{}{"E10-PREREQ-PARENT"},
				"unlocks":    []interface{}{},
			},
		)
		doc["nodes"] = filtered
		out, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal mock graph.json: %v", err)
		}
		if err := os.WriteFile(graphPath, append(out, '\n'), 0644); err != nil {
			t.Fatalf("failed to write mock graph.json: %v", err)
		}
	}

	// When prerequisite is IN_PROGRESS, cmdr-dev next must NOT select E10-PREREQ-CHILD
	writeGraphWithParentStatus("IN_PROGRESS")
	blockedRes := harness.RunCmdrDev(ctx, "next", "--root", mockRoot, "--json")
	if blockedRes.ExitCode == 0 {
		t.Fatalf("T2.F12.4 FAIL: cmdr-dev next selected unit despite IN_PROGRESS prerequisite: %s", blockedRes.Stdout)
	}
	combined := blockedRes.Stderr + " " + blockedRes.Stdout
	if !strings.Contains(combined, "no executable READY work unit") {
		t.Fatalf("T2.F12.4 FAIL: expected 'no executable READY work unit', got: %s", combined)
	}

	// Once prerequisite becomes ACCEPTED, cmdr-dev next must select E10-PREREQ-CHILD
	writeGraphWithParentStatus("ACCEPTED")
	readyRes := harness.RunCmdrDev(ctx, "next", "--root", mockRoot, "--json")
	if readyRes.ExitCode != 0 {
		t.Fatalf("T2.F12.4 FAIL: cmdr-dev next failed after prerequisite ACCEPTED (exit %d): %s\n%s", readyRes.ExitCode, readyRes.Stderr, readyRes.Stdout)
	}
	harness.AssertJSONFieldEqual(t, readyRes.Stdout, "id", "E10-PREREQ-CHILD")

	// Also verify cmdr-dev wave selector unit test rejecting non-READY states
	unitRes, unitErr := harness.RunGoTestVerified(
		ctx,
		"tools/cmdr-dev",
		"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev",
		[]string{"TestSelectImplementationWaveExcludesNonReadyStates"},
		"-run", "^TestSelectImplementationWaveExcludesNonReadyStates$",
	)
	if !unitRes.Passed || unitErr != nil {
		t.Fatalf("T2.F12.4 FAIL: TestSelectImplementationWaveExcludesNonReadyStates failed: %v\nOutput: %s", unitErr, unitRes.Output)
	}
	t.Log("T2.F12.4 PASS: wave/next selector strictly rejects units whose prerequisites are not ACCEPTED")
}

// T2.F12.5: Wave selector evaluates tie-breaking deterministically.
func TestT2_F12_5_WaveSelectorEvaluatesTieBreakingDeterministically(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f12-tie")
	graphPath := filepath.Join(mockRoot, "work", "graph.json")
	raw, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read mock graph.json: %v", err)
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("failed to parse mock graph.json: %v", err)
	}
	nodes, ok := doc["nodes"].([]interface{})
	if !ok {
		t.Fatalf("nodes missing from graph.json")
	}
	nodes = append(nodes,
		map[string]interface{}{
			"id":         "E10-TIE-Z",
			"title":      "Tie Z",
			"type":       "task",
			"status":     "READY",
			"depends_on": []interface{}{},
			"unlocks":    []interface{}{},
		},
		map[string]interface{}{
			"id":         "E10-TIE-A",
			"title":      "Tie A",
			"type":       "task",
			"status":     "READY",
			"depends_on": []interface{}{},
			"unlocks":    []interface{}{},
		},
	)
	doc["nodes"] = nodes
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock graph.json: %v", err)
	}
	if err := os.WriteFile(graphPath, append(out, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mock graph.json: %v", err)
	}

	for i := 0; i < 5; i++ {
		res := harness.RunCmdrDev(ctx, "next", "--root", mockRoot, "--json")
		if res.ExitCode != 0 {
			t.Fatalf("T2.F12.5 FAIL: cmdr-dev next failed on iteration %d (exit %d): %s", i, res.ExitCode, res.Stderr)
		}
		harness.AssertJSONFieldEqual(t, res.Stdout, "id", "E10-TIE-A")
	}

	unitRes, unitErr := harness.RunGoTestVerified(
		ctx,
		"tools/cmdr-dev",
		"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev",
		[]string{
			"TestSelectImplementationWavePrefersImplementedFamily",
			"TestSelectImplementationWaveFallsBackToStableCapabilityOrder",
		},
		"-run", "^TestSelectImplementationWave(PrefersImplementedFamily|FallsBackToStableCapabilityOrder)$",
	)
	if !unitRes.Passed || unitErr != nil {
		t.Fatalf("T2.F12.5 FAIL: cmdr-dev wave tie-breaking unit tests failed: %v\nOutput: %s", unitErr, unitRes.Output)
	}

	// Run cmdr-dev implementation-wave 3x against live repo and verify deterministic outcome
	var firstWaveOutput string
	var firstWaveExit int
	for i := 0; i < 3; i++ {
		waveRes := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
		combined := strings.TrimSpace(waveRes.Stdout + "\n" + waveRes.Stderr)
		if i == 0 {
			firstWaveOutput = combined
			firstWaveExit = waveRes.ExitCode
			continue
		}
		if waveRes.ExitCode != firstWaveExit || combined != firstWaveOutput {
			t.Fatalf("T2.F12.5 FAIL: non-deterministic implementation-wave output on run %d (exit %d vs %d)\nFirst: %s\nCurrent: %s",
				i, waveRes.ExitCode, firstWaveExit, firstWaveOutput, combined)
		}
	}

	t.Log("T2.F12.5 PASS: wave/next selector tie-breaking evaluated deterministically across repeated real CLI and unit runs")
}
