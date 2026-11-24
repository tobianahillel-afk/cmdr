package tier1_features

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type WorkGraph struct {
	Nodes []struct {
		ID        string   `json:"id"`
		Title     string   `json:"title"`
		Status    string   `json:"status"`
		DependsOn []string `json:"depends_on"`
	} `json:"nodes"`
}

// T1.F9.1: Verify work/graph.json contains node E10-STD-004A-INPUT.
func TestT1_F9_1_WorkGraphContainsCAPSTD004Node(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	graphPath := filepath.Join(root, "work", "graph.json")
	data, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read work/graph.json: %v", err)
	}

	var graph WorkGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		t.Fatalf("failed to parse work/graph.json: %v", err)
	}

	found := false
	for _, node := range graph.Nodes {
		if node.ID == "E10-STD-004A-INPUT" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("node 'E10-STD-004A-INPUT' not found in work/graph.json")
	}
	t.Log("T1.F9.1 PASS: node E10-STD-004A-INPUT is registered in work/graph.json")
}

// T1.F9.2: Verify node status is ACCEPTED.
func TestT1_F9_2_WorkGraphNodeStatusAccepted(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	graphPath := filepath.Join(root, "work", "graph.json")
	data, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read work/graph.json: %v", err)
	}

	var graph WorkGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		t.Fatalf("failed to parse work/graph.json: %v", err)
	}

	for _, node := range graph.Nodes {
		if node.ID == "E10-STD-004A-INPUT" {
			if node.Status != "ACCEPTED" {
				t.Fatalf("expected node E10-STD-004A-INPUT status 'ACCEPTED', got '%s'", node.Status)
			}
			t.Log("T1.F9.2 PASS: node E10-STD-004A-INPUT status is ACCEPTED")
			return
		}
	}
	t.Fatal("node E10-STD-004A-INPUT not found in work/graph.json")
}

// T1.F9.3: Verify work/lots/E10-STD-004A-INPUT/manifest.json exists and validates.
func TestT1_F9_3_CAPSTD004ManifestValidation(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	manifestPath := filepath.Join(root, "work", "lots", "E10-STD-004A-INPUT", "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("manifest.json missing for E10-STD-004A-INPUT: %v", err)
	}

	var manifest struct {
		SchemaVersion int    `json:"schema_version"`
		ID            string `json:"id"`
		Title         string `json:"title"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("failed to parse manifest.json: %v", err)
	}

	if manifest.ID != "E10-STD-004A-INPUT" {
		t.Fatalf("expected manifest id 'E10-STD-004A-INPUT', got '%s'", manifest.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "validate-manifests", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("validate-manifests failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	t.Log("T1.F9.3 PASS: manifest.json exists, correctly formatted, and validates")
}

// T1.F9.4: Verify work/lots/E10-STD-004A-INPUT/HANDOFF.json exists with valid evidence.
func TestT1_F9_4_CAPSTD004HandoffValidation(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	handoffPath := filepath.Join(root, "work", "lots", "E10-STD-004A-INPUT", "HANDOFF.json")
	data, err := os.ReadFile(handoffPath)
	if err != nil {
		t.Fatalf("HANDOFF.json missing for E10-STD-004A-INPUT: %v", err)
	}

	var handoff struct {
		WorkUnit           string `json:"work_unit"`
		Result             string `json:"result"`
		FinalValidatedHead string `json:"final_validated_head"`
	}
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatalf("failed to parse HANDOFF.json: %v", err)
	}

	if handoff.WorkUnit != "E10-STD-004A-INPUT" {
		t.Fatalf("expected work_unit 'E10-STD-004A-INPUT', got '%s'", handoff.WorkUnit)
	}
	if handoff.Result != "ACCEPTED" && handoff.Result != "VERIFIED" {
		t.Fatalf("expected result 'ACCEPTED' or 'VERIFIED', got '%s'", handoff.Result)
	}
	t.Log("T1.F9.4 PASS: HANDOFF.json exists and contains valid verified evidence")
}

// T1.F9.5: Verify graph acyclicity with cmdr-dev doctor.
func TestT1_F9_5_DoctorAcyclicWorkGraph(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "doctor", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev doctor failed with exit code %d: %s", res.ExitCode, res.Stderr)
	}

	var doc struct {
		Healthy bool     `json:"healthy"`
		Checks  []string `json:"checks"`
	}
	if err := harness.ParseJSON(res.Stdout, &doc); err != nil {
		t.Fatalf("failed to parse doctor output: %v", err)
	}

	if !doc.Healthy {
		t.Fatal("cmdr-dev doctor reported unhealthy state")
	}

	foundCheck := false
	for _, check := range doc.Checks {
		if check == "work graph is acyclic" {
			foundCheck = true
			break
		}
	}
	if !foundCheck {
		t.Fatal("doctor checks do not include 'work graph is acyclic'")
	}
	t.Log("T1.F9.5 PASS: doctor confirms work graph is strictly acyclic and healthy")
}
