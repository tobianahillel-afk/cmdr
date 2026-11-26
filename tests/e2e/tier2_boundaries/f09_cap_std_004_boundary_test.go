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

func setupMockGraphRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "graph-mock")

	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
	}

	for _, rel := range copyFiles {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(mockRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatalf("failed to create dir for %s: %v", dst, err)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("failed to copy %s to mock: %v", rel, err)
		}
	}

	// Also copy work/lots directory for manifest validation
	lotsSrc := filepath.Join(root, "work", "lots")
	lotsDst := filepath.Join(mockRoot, "work", "lots")
	_ = filepath.Walk(lotsSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		rel, rErr := filepath.Rel(lotsSrc, path)
		if rErr != nil {
			return nil
		}
		target := filepath.Join(lotsDst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, rErr := os.ReadFile(path)
		if rErr == nil {
			_ = os.WriteFile(target, data, 0644)
		}
		return nil
	})

	// Copy canonical product spec canonical path
	specDir := filepath.Join(mockRoot, "cmdr-product-spec")
	_ = os.MkdirAll(specDir, 0755)

	return mockRoot
}

// T2.F9.1: Work graph validation fails closed on cyclic dependency involving E10-STD-004A-INPUT.
func TestT2_F9_1_WorkGraphFailsOnCyclicDependency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockGraphRepo(t)
	graphPath := filepath.Join(mockRoot, "work", "graph.json")

	data, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read mock graph: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse graph JSON: %v", err)
	}

	nodes, ok := doc["nodes"].([]interface{})
	if !ok || len(nodes) == 0 {
		t.Fatalf("nodes missing from graph JSON")
	}

	// Introduce a cyclic dependency involving existing E10-STD-004A-INPUT and a child node pointing back
	foundTarget := false
	for _, raw := range nodes {
		m, ok := raw.(map[string]interface{})
		if !ok || m["id"] != "E10-STD-004A-INPUT" {
			continue
		}
		deps, _ := m["depends_on"].([]interface{})
		m["depends_on"] = append(deps, "E10-STD-004-CYCLE-DEP")
		unlocks, _ := m["unlocks"].([]interface{})
		m["unlocks"] = append(unlocks, "E10-STD-004-CYCLE-DEP")
		foundTarget = true
		break
	}
	if !foundTarget {
		t.Fatalf("T2.F9.1 setup error: E10-STD-004A-INPUT missing from work/graph.json")
	}
	cycleNode2 := map[string]interface{}{
		"id":         "E10-STD-004-CYCLE-DEP",
		"title":      "Cyclic child",
		"type":       "task",
		"status":     "ACCEPTED",
		"depends_on": []interface{}{"E10-STD-004A-INPUT"},
		"unlocks":    []interface{}{"E10-STD-004A-INPUT"},
	}

	doc["nodes"] = append(nodes, cycleNode2)
	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated graph: %v", err)
	}
	if err := os.WriteFile(graphPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated graph: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "doctor", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F9.1 FAIL: doctor passed despite cyclic dependency in work graph")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "cyclic") && !strings.Contains(combined, "cycle") {
		t.Fatalf("T2.F9.1 FAIL: expected cyclic dependency error, got: %s", combined)
	}
	t.Log("T2.F9.1 PASS: work graph validation strictly fails closed on cyclic dependency")
}

// T2.F9.2: Work graph validation fails on duplicate node ID.
func TestT2_F9_2_WorkGraphFailsOnDuplicateNodeID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockGraphRepo(t)
	graphPath := filepath.Join(mockRoot, "work", "graph.json")

	data, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read mock graph: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse graph JSON: %v", err)
	}

	nodes, ok := doc["nodes"].([]interface{})
	if !ok || len(nodes) == 0 {
		t.Fatalf("nodes missing from graph JSON")
	}

	// Duplicate the first node
	doc["nodes"] = append(nodes, nodes[0])
	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated graph: %v", err)
	}
	if err := os.WriteFile(graphPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated graph: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "doctor", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F9.2 FAIL: doctor passed despite duplicate node ID")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "duplicate") {
		t.Fatalf("T2.F9.2 FAIL: expected duplicate node id error, got: %s", combined)
	}
	t.Log("T2.F9.2 PASS: work graph validation strictly rejects duplicate work node IDs")
}

// T2.F9.3: Manifest validation / architecture audit fails on unregistered root path ownership.
func TestT2_F9_3_ManifestValidationFailsOnUnregisteredRootPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f09-arch-bounds")
	baseRes := harness.RunCmdrDev(ctx, "architecture-audit", "--root", mockRoot)
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F9.3 baseline architecture-audit failed (exit %d): %s\n%s", baseRes.ExitCode, baseRes.Stderr, baseRes.Stdout)
	}

	manifestPath := filepath.Join(mockRoot, "work", "lots", "E10-INV-002A", "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to parse manifest JSON: %v", err)
	}

	// Mutate allowed_paths to an unauthorized path outside registered boundaries
	doc["allowed_paths"] = []interface{}{
		"outside/unregistered/path/**",
	}

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated manifest: %v", err)
	}
	if err := os.WriteFile(manifestPath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated manifest: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "architecture-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F9.3 FAIL: architecture-audit passed with unregistered allowed_paths")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "outside registered architecture boundaries") {
		t.Fatalf("T2.F9.3 FAIL: expected 'outside registered architecture boundaries' error, got: %s", combined)
	}
	t.Log("T2.F9.3 PASS: architecture-audit strictly fails closed on unregistered allowed_paths")
}

// T2.F9.4: Manifest validation fails on missing manifest.json.
func TestT2_F9_4_ManifestValidationFailsOnMissingManifest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f09-missing-manifest")
	baseRes := harness.RunCmdrDev(ctx, "validate-manifests", "--root", mockRoot)
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F9.4 baseline validate-manifests failed (exit %d): %s\n%s", baseRes.ExitCode, baseRes.Stderr, baseRes.Stdout)
	}

	// Remove E0-COVERAGE-001C manifest which E1-WORK-001 depends_on
	missingManifest := filepath.Join(mockRoot, "work", "lots", "E0-COVERAGE-001C", "manifest.json")
	if err := os.Remove(missingManifest); err != nil {
		t.Fatalf("failed to remove E0-COVERAGE-001C manifest.json: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "validate-manifests", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F9.4 FAIL: validate-manifests succeeded despite missing dependency manifest E0-COVERAGE-001C")
	}
	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "has no manifest") {
		t.Fatalf("T2.F9.4 FAIL: expected 'has no manifest' error, got: %s", combined)
	}
	t.Logf("T2.F9.4 PASS: validate-manifests failed closed on missing manifest (exit %d)", res.ExitCode)
}

// T2.F9.5: Manifest validation fails on schema type mismatch.
func TestT2_F9_5_ManifestValidationFailsOnSchemaTypeMismatch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mockRoot := harness.CloneIsolatedRepo(t, "f09-schema-mismatch")
	baseRes := harness.RunCmdrDev(ctx, "validate-manifests", "--root", mockRoot)
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F9.5 baseline validate-manifests failed (exit %d): %s\n%s", baseRes.ExitCode, baseRes.Stderr, baseRes.Stdout)
	}

	manifestPath := filepath.Join(mockRoot, "work", "lots", "E10-INV-002A", "manifest.json")

	// Mutate schema_version to string instead of int
	corrupted := []byte(`{
  "schema_version": "two",
  "id": "E10-INV-002A"
}`)
	if err := os.WriteFile(manifestPath, corrupted, 0644); err != nil {
		t.Fatalf("failed to write corrupted manifest: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "validate-manifests", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F9.5 FAIL: validate-manifests succeeded despite schema type mismatch")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "schema_version") &&
		!strings.Contains(combined, "cannot unmarshal") &&
		!strings.Contains(combined, "invalid") {
		t.Fatalf("T2.F9.5 FAIL: expected type mismatch error, got: %s", combined)
	}
	t.Log("T2.F9.5 PASS: manifest validation strictly rejects schema type mismatches")
}
