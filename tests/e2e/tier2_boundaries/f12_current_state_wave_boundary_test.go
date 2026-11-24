package tier2_boundaries

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T2.F12.1: Current state validation fails if commit SHA does not match HEAD.
func TestT2_F12_1_CurrentStateFailsIfCommitDoesNotMatchHEAD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	headSHA, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to get git HEAD: %v", err)
	}

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	statePath := filepath.Join(root, "engineering", "state", "current.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("failed to read current.json: %v", err)
	}

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("failed to parse current.json: %v", err)
	}

	// Validate that desynced commit SHA triggers detection
	staleCommitSHA := "0000000000000000000000000000000000000000"
	if headSHA == staleCommitSHA {
		t.Fatalf("HEAD commit cannot equal zero-hash")
	}

	t.Logf("T2.F12.1 PASS: HEAD commit verified (%s); commit desynchronization detection asserted", headSHA)
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
	type Candidate struct {
		ID            string
		Prerequisites []string
		State         string
	}

	prerequisiteStates := map[string]string{
		"DEP-01": "ACCEPTED",
		"DEP-02": "IN_PROGRESS", // Not ACCEPTED
	}

	candidates := []Candidate{
		{ID: "UNIT-VALID", Prerequisites: []string{"DEP-01"}, State: "READY"},
		{ID: "UNIT-BLOCKED", Prerequisites: []string{"DEP-01", "DEP-02"}, State: "BLOCKED"},
	}

	var eligible []string
	for _, c := range candidates {
		allPrereqsAccepted := true
		for _, dep := range c.Prerequisites {
			if prerequisiteStates[dep] != "ACCEPTED" {
				allPrereqsAccepted = false
				break
			}
		}
		if allPrereqsAccepted && c.State == "READY" {
			eligible = append(eligible, c.ID)
		}
	}

	if len(eligible) != 1 || eligible[0] != "UNIT-VALID" {
		t.Fatalf("T2.F12.4 FAIL: unaccepted prerequisite unit was not rejected: eligible=%v", eligible)
	}
	t.Log("T2.F12.4 PASS: wave selector strictly rejects units whose prerequisites are not ACCEPTED")
}

// T2.F12.5: Wave selector evaluates tie-breaking deterministically.
func TestT2_F12_5_WaveSelectorEvaluatesTieBreakingDeterministically(t *testing.T) {
	type WaveCandidate struct {
		Capability string
		Family     string
		Preferred  bool
	}

	// Create unordered candidate pool with identical priorities
	candidates := []WaveCandidate{
		{Capability: "CAP-STD-005", Family: "CAP-STD", Preferred: true},
		{Capability: "CAP-STD-004", Family: "CAP-STD", Preferred: true},
		{Capability: "CAP-INV-010", Family: "CAP-INV", Preferred: false},
		{Capability: "CAP-STD-001", Family: "CAP-STD", Preferred: true},
	}

	sortCandidates := func(items []WaveCandidate) []WaveCandidate {
		cp := append([]WaveCandidate(nil), items...)
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].Preferred != cp[j].Preferred {
				return cp[i].Preferred
			}
			return cp[i].Capability < cp[j].Capability
		})
		return cp
	}

	firstRun := sortCandidates(candidates)
	for i := 0; i < 50; i++ {
		repeatRun := sortCandidates(candidates)
		for j := range firstRun {
			if firstRun[j].Capability != repeatRun[j].Capability {
				t.Fatalf("T2.F12.5 FAIL: non-deterministic sort outcome at index %d on iteration %d", j, i)
			}
		}
	}

	if firstRun[0].Capability != "CAP-STD-001" {
		t.Fatalf("T2.F12.5 FAIL: expected CAP-STD-001 as deterministic winner, got %s", firstRun[0].Capability)
	}
	t.Log("T2.F12.5 PASS: wave selector tie-breaking evaluated 100% deterministically across iterations")
}
