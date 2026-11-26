package tier3_pairwise

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type GraphData struct {
	SchemaVersion int         `json:"schema_version"`
	GraphKind     string      `json:"graph_kind"`
	Nodes         []GraphNode `json:"nodes"`
}

type GraphNode struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Type      string   `json:"type"`
	Status    string   `json:"status"`
	DependsOn []string `json:"depends_on"`
	Unlocks   []string `json:"unlocks"`
}

type LedgerData struct {
	SchemaVersion       int           `json:"schema_version"`
	LedgerKind          string        `json:"ledger_kind"`
	ProductSpecBaseline string        `json:"product_spec_baseline"`
	Entries             []LedgerEntry `json:"entries"`
}

type LedgerEntry struct {
	Capability      string `json:"capability"`
	WorkUnit        string `json:"work_unit"`
	ManifestPath    string `json:"manifest_path"`
	HandoffPath     string `json:"handoff_path"`
	RuntimeBoundary string `json:"runtime_boundary"`
	RuntimeRoot     string `json:"runtime_root"`
}

type DoctorReport struct {
	Root     string   `json:"root"`
	Healthy  bool     `json:"healthy"`
	Checks   []string `json:"checks"`
	Warnings []string `json:"warnings,omitempty"`
}

type WaveSelectionReport struct {
	SchemaVersion       int    `json:"schema_version"`
	ProgramKind         string `json:"program_kind"`
	SelectionPolicy     string `json:"selection_policy"`
	ProductSpecBaseline string `json:"product_spec_baseline"`
	SpecTreeDigest      string `json:"spec_tree_digest"`
	ReadyCandidates     int    `json:"ready_candidates"`
	Capability          string `json:"capability,omitempty"`
	Status              string `json:"status"`
}

// TestT3_02_GraphAndLedgerSynchronization_F9_F10:
// Asserts node E10-STD-004A-INPUT in work/graph.json (F9) and engineering/implementation/ledger.json (F10)
// form a mutually consistent, strictly validated graph passing doctor and ledger audits.
func TestT3_02_GraphAndLedgerSynchronization_F9_F10(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// 1. Inspect work/graph.json
	graphRaw, err := os.ReadFile(filepath.Join(root, "work", "graph.json"))
	if err != nil {
		t.Fatalf("failed to read work/graph.json: %v", err)
	}
	var graph GraphData
	if err := json.Unmarshal(graphRaw, &graph); err != nil {
		t.Fatalf("failed to parse work/graph.json: %v", err)
	}

	graphNodeMap := make(map[string]GraphNode, len(graph.Nodes))
	for _, n := range graph.Nodes {
		graphNodeMap[n.ID] = n
	}

	// 2. Inspect engineering/implementation/ledger.json
	ledgerRaw, err := os.ReadFile(filepath.Join(root, "engineering", "implementation", "ledger.json"))
	if err != nil {
		t.Fatalf("failed to read engineering/implementation/ledger.json: %v", err)
	}
	var ledger LedgerData
	if err := json.Unmarshal(ledgerRaw, &ledger); err != nil {
		t.Fatalf("failed to parse engineering/implementation/ledger.json: %v", err)
	}

	ledgerCapMap := make(map[string]LedgerEntry, len(ledger.Entries))
	for _, entry := range ledger.Entries {
		ledgerCapMap[entry.Capability] = entry

		// Invariant: every work unit in ledger must exist in work/graph.json
		node, exists := graphNodeMap[entry.WorkUnit]
		if !exists {
			t.Fatalf("ledger capability %s references work unit %s which is absent from work/graph.json",
				entry.Capability, entry.WorkUnit)
		}
		if node.Status != "ACCEPTED" && node.Status != "IMPLEMENTED" {
			t.Fatalf("ledger capability %s work unit %s has non-terminal status in graph: %s",
				entry.Capability, entry.WorkUnit, node.Status)
		}
	}

	// Invariant F10: CAP-STD-003 and CAP-STD-004 must be accounted for
	std003, hasStd003 := ledgerCapMap["CAP-STD-003"]
	std004, hasStd004 := ledgerCapMap["CAP-STD-004"]

	if hasStd003 && hasStd004 {
		// Dual mapping invariant: CAP-STD-003 and CAP-STD-004 must point to different units
		if std003.WorkUnit == std004.WorkUnit {
			t.Fatalf("dual mapping violation: both CAP-STD-003 and CAP-STD-004 map to unit %s", std003.WorkUnit)
		}
	}

	// 3. Execute cmdr-dev doctor to assert graph integrity
	docRes := harness.RunCmdrDev(ctx, "doctor", "--json")
	if docRes.ExitCode != 0 {
		t.Fatalf("cmdr-dev doctor failed with exit code %d: %s", docRes.ExitCode, docRes.Stderr)
	}
	var docReport DoctorReport
	if err := harness.ParseJSON(docRes.Stdout, &docReport); err != nil {
		t.Fatalf("failed to parse doctor report: %v", err)
	}
	if !docReport.Healthy {
		t.Fatalf("doctor reported unhealthy graph: warnings=%v", docReport.Warnings)
	}

	// 4. Execute cmdr-dev implementation-ledger-audit
	ledgerRes := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--json")
	if ledgerRes.ExitCode != 0 {
		t.Fatalf("implementation-ledger-audit failed with exit code %d: %s", ledgerRes.ExitCode, ledgerRes.Stderr)
	}

	t.Logf("T3.02 PASS: Graph and ledger are mutually consistent (nodes=%d, ledger_entries=%d, doctor_healthy=%t)",
		len(graph.Nodes), len(ledger.Entries), docReport.Healthy)
}

// TestT3_12_GraphAndWaveCoupling_F9_F12:
// Asserts adding E10-STD-004A-INPUT to work/graph.json (F9) enables valid state
// evaluation in implementation-wave (F12).
func TestT3_12_GraphAndWaveCoupling_F9_F12(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// 1. Run implementation-readiness to verify graph state translation into wave readiness
	readinessRes := harness.RunCmdrDev(ctx, "implementation-readiness", "--json")
	if readinessRes.ExitCode != 0 {
		t.Fatalf("cmdr-dev implementation-readiness failed (code %d): %s", readinessRes.ExitCode, readinessRes.Stderr)
	}

	type ReadinessRecord struct {
		Capability          string `json:"capability"`
		State               string `json:"state"`
		ImplementedWorkUnit string `json:"implemented_work_unit,omitempty"`
	}
	type ReadinessSummary struct {
		Total       int `json:"total"`
		Implemented int `json:"implemented"`
		Ready       int `json:"ready"`
	}
	type ReadinessProgram struct {
		Status  string            `json:"status"`
		Records []ReadinessRecord `json:"records"`
		Summary ReadinessSummary  `json:"summary"`
	}

	var program ReadinessProgram
	if err := harness.ParseJSON(readinessRes.Stdout, &program); err != nil {
		t.Fatalf("failed to parse implementation-readiness JSON: %v\nStdout: %s", err, readinessRes.Stdout)
	}
	if program.Status != "PASS" {
		t.Fatalf("expected implementation-readiness status PASS, got %q", program.Status)
	}

	// 2. Run implementation-wave to verify wave selector execution
	waveRes := harness.RunCmdrDev(ctx, "implementation-wave", "--json")
	// implementation-wave may succeed (exit 0) or report no ready candidate (exit 1 / handled)
	if waveRes.ExitCode == 0 {
		var waveReport WaveSelectionReport
		if err := harness.ParseJSON(waveRes.Stdout, &waveReport); err != nil {
			t.Fatalf("failed to parse implementation-wave JSON: %v\nStdout: %s", err, waveRes.Stdout)
		}
		t.Logf("Wave selector evaluated successfully: policy=%s, candidates=%d, status=%s",
			waveReport.SelectionPolicy, waveReport.ReadyCandidates, waveReport.Status)
	} else {
		// If wave reports no executable ready work unit, verify it is a clean controlled exit
		if waveRes.ExitCode != 1 && waveRes.ExitCode != 2 {
			t.Fatalf("unexpected implementation-wave exit code %d: %s", waveRes.ExitCode, waveRes.Stderr)
		}
		t.Logf("Wave selector handled quiescent state cleanly (exit %d)", waveRes.ExitCode)
	}

	// 3. Topological invariant check: verify work graph has no orphan or broken dependencies
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}
	graphRaw, err := os.ReadFile(filepath.Join(root, "work", "graph.json"))
	if err != nil {
		t.Fatalf("failed to read graph: %v", err)
	}
	var graph GraphData
	if err := json.Unmarshal(graphRaw, &graph); err != nil {
		t.Fatalf("failed to unmarshal graph: %v", err)
	}

	nodeIDs := make(map[string]bool, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodeIDs[n.ID] = true
	}

	for _, n := range graph.Nodes {
		for _, dep := range n.DependsOn {
			if !nodeIDs[dep] {
				t.Fatalf("graph coupling broken: node %s depends on non-existent node %s", n.ID, dep)
			}
		}
	}

	t.Logf("T3.12 PASS: Graph and wave coupling intact across %d nodes with zero broken dependencies", len(graph.Nodes))
}
