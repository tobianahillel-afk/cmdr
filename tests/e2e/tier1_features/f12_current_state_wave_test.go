package tier1_features

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

type StatusResult struct {
	Root             string `json:"root"`
	Phase            string `json:"phase"`
	PhaseTitle       string `json:"phase_title"`
	FoundationStatus string `json:"foundation_status"`
	ActiveWorkUnit   string `json:"active_work_unit"`
	PreferredNext    string `json:"preferred_next_unit"`
	ProductBaseline  string `json:"product_spec_baseline"`
}

type WaveResult struct {
	SchemaVersion   int    `json:"schema_version"`
	SelectionPolicy string `json:"selection_policy"`
	ReadyCandidates int    `json:"ready_candidates"`
	Capability      string `json:"capability"`
	Status          string `json:"status"`
}

// T1.F12.1: Verify engineering/state/current.json references current git HEAD commit or valid baseline.
func TestT1_F12_1_CurrentStateReferencesValidCommit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	statePath := filepath.Join(root, "engineering", "state", "current.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("failed to read current.json: %v", err)
	}

	var state struct {
		ProductSpec struct {
			BaselineCommit string `json:"baseline_commit"`
		} `json:"product_spec"`
		EngineeringFoundation struct {
			Status string `json:"status"`
		} `json:"engineering_foundation"`
	}
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("failed to parse current.json: %v", err)
	}

	if !harness.GitCommitExists(ctx, state.ProductSpec.BaselineCommit) {
		t.Fatalf("baseline_commit '%s' in current.json does not exist in git", state.ProductSpec.BaselineCommit)
	}
	t.Logf("T1.F12.1 PASS: current.json references existing commit baseline %s", state.ProductSpec.BaselineCommit)
}

// T1.F12.2: Verify cmdr-dev status executes without errors.
func TestT1_F12_2_CmdrDevStatusExecution(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "status", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev status failed (code %d): %s", res.ExitCode, res.Stderr)
	}

	var status StatusResult
	if err := harness.ParseJSON(res.Stdout, &status); err != nil {
		t.Fatalf("failed to parse cmdr-dev status output: %v", err)
	}

	if status.FoundationStatus != "ACCEPTED" {
		t.Fatalf("expected foundation_status 'ACCEPTED', got '%s'", status.FoundationStatus)
	}
	t.Logf("T1.F12.2 PASS: cmdr-dev status succeeded (phase=%s, foundation=%s)", status.Phase, status.FoundationStatus)
}

// T1.F12.3: Verify cmdr-dev implementation-wave executes and reports selection policy.
func TestT1_F12_3_ImplementationWavePolicy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev implementation-wave failed (code %d): %s", res.ExitCode, res.Stderr)
	}

	var wave WaveResult
	if err := harness.ParseJSON(res.Stdout, &wave); err != nil {
		t.Fatalf("failed to parse implementation-wave output: %v", err)
	}

	if wave.SelectionPolicy == "" {
		t.Fatal("implementation-wave output has empty selection_policy")
	}
	t.Logf("T1.F12.3 PASS: implementation-wave policy verified: %s", wave.SelectionPolicy)
}

// T1.F12.4: Verify wave selector returns a valid candidate or clean completed/quiescent state.
func TestT1_F12_4_WaveSelectorValidCandidateOrQuiescent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev implementation-wave failed: %s", res.Stderr)
	}

	var wave WaveResult
	if err := harness.ParseJSON(res.Stdout, &wave); err != nil {
		t.Fatalf("failed to parse implementation-wave output: %v", err)
	}

	if wave.Capability != "" {
		if !strings.HasPrefix(wave.Capability, "CAP-") {
			t.Fatalf("invalid candidate capability format '%s'", wave.Capability)
		}
	}
	t.Logf("T1.F12.4 PASS: wave selector returned candidate=%s, status=%s", wave.Capability, wave.Status)
}

type ReadinessSummary struct {
	Total             int `json:"total"`
	Implemented       int `json:"implemented"`
	Ready             int `json:"ready"`
	Blocked           int `json:"blocked"`
	Proposed          int `json:"proposed"`
	DecisionBlocked   int `json:"decision_blocked"`
	DependencyBlocked int `json:"dependency_blocked"`
}

type ReadinessRecord struct {
	Capability string `json:"capability"`
	State      string `json:"state"`
}

type ReadinessResult struct {
	SchemaVersion int               `json:"schema_version"`
	ProgramKind   string            `json:"program_kind"`
	Records       []ReadinessRecord `json:"records"`
	Summary       ReadinessSummary  `json:"summary"`
	Status        string            `json:"status"`
}

// T1.F12.5: Verify wave selector does not blindly select unready units.
func TestT1_F12_5_WaveSelectorRespectsReadiness(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "implementation-readiness", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("T1.F12.5 FAIL: cmdr-dev implementation-readiness failed (code %d): stdout=%s stderr=%s",
			res.ExitCode, res.Stdout, res.Stderr)
	}

	harness.AssertJSONFieldExists(t, res.Stdout, "program_kind")
	harness.AssertJSONFieldEqual(t, res.Stdout, "program_kind", "capability-implementation-readiness")
	harness.AssertJSONFieldEqual(t, res.Stdout, "status", "PASS")

	var readiness ReadinessResult
	if err := harness.ParseJSON(res.Stdout, &readiness); err != nil {
		t.Fatalf("T1.F12.5 FAIL: failed to parse implementation-readiness JSON: %v", err)
	}
	if readiness.SchemaVersion != 1 {
		t.Fatalf("T1.F12.5 FAIL: expected schema_version=1, got %d", readiness.SchemaVersion)
	}
	if readiness.Summary.Total <= 0 {
		t.Fatalf("T1.F12.5 FAIL: expected summary.total > 0, got %d", readiness.Summary.Total)
	}
	accounted := readiness.Summary.Implemented + readiness.Summary.Ready + readiness.Summary.Blocked + readiness.Summary.Proposed
	if readiness.Summary.Total != accounted {
		t.Fatalf("T1.F12.5 FAIL: readiness summary accounting mismatch: total=%d, sum=%d", readiness.Summary.Total, accounted)
	}

	waveRes := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	if waveRes.ExitCode != 0 {
		t.Fatalf("T1.F12.5 FAIL: cmdr-dev implementation-wave failed (code %d): stdout=%s stderr=%s",
			waveRes.ExitCode, waveRes.Stdout, waveRes.Stderr)
	}
	var wave WaveResult
	if err := harness.ParseJSON(waveRes.Stdout, &wave); err != nil {
		t.Fatalf("T1.F12.5 FAIL: failed to parse implementation-wave JSON: %v", err)
	}
	if wave.Capability == "" {
		t.Fatalf("T1.F12.5 FAIL: implementation-wave returned empty capability")
	}

	foundReady := false
	for _, rec := range readiness.Records {
		if rec.Capability == wave.Capability {
			if rec.State != "READY" {
				t.Fatalf("T1.F12.5 FAIL: wave selector chose %s in non-READY state %q", wave.Capability, rec.State)
			}
			foundReady = true
			break
		}
	}
	if !foundReady {
		t.Fatalf("T1.F12.5 FAIL: selected capability %s not found in implementation-readiness records", wave.Capability)
	}

	t.Logf("T1.F12.5 PASS: readiness verified (%s, total=%d, ready=%d) and wave selector chose READY candidate %s",
		readiness.ProgramKind, readiness.Summary.Total, readiness.Summary.Ready, wave.Capability)
}
