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

// setupMockSecurityRepo copies the architecture, policy, state, and graph files
// into an isolated temp dir under .cmdr-tmp/ to safely test security audit boundary mutations.
func setupMockSecurityRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "sec-evidence-mock")

	// Bootstrap files for isRepoRoot
	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
		"engineering/architecture/boundaries.json",
		"engineering/security/security-gates.json",
		"engineering/security/runtime-security-policy.json",
		"engineering/testing/runtime-security-evidence.json",
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

	// Create directories and stub files for all 12 runtime modules so detectRuntimeImplementationStates resolves
	for _, mod := range harness.RuntimeGoModules {
		dir := filepath.Join(mockRoot, filepath.FromSlash(mod))
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create runtime stub dir %s: %v", dir, err)
		}
		_ = os.WriteFile(filepath.Join(dir, "stub.go"), []byte("package stub\n"), 0644)
	}

	// Create event-search-frontend runtime boundary directory and stub file
	frontendDir := filepath.Join(mockRoot, "product-runtime", "event-search-frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		t.Fatalf("failed to create event-search-frontend dir: %v", err)
	}
	_ = os.WriteFile(filepath.Join(frontendDir, "stub.js"), []byte("// stub\n"), 0644)

	// Create cmdr-product-spec directory so validateState succeeds
	mockSpecDir := filepath.Join(mockRoot, "cmdr-product-spec")
	if err := os.MkdirAll(mockSpecDir, 0755); err != nil {
		t.Fatalf("failed to create cmdr-product-spec dir: %v", err)
	}
	specReadme := filepath.Join(root, "cmdr-product-spec", "README.md")
	if data, err := os.ReadFile(specReadme); err == nil {
		_ = os.WriteFile(filepath.Join(mockSpecDir, "README.md"), data, 0644)
	}

	return mockRoot
}

// T2.F4.1: Security test audit fails closed if scope count is 10 instead of 11.
func TestT2_F4_1_SecurityAuditFailsOnTenScopes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockSecurityRepo(t)
	evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")

	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read mock evidence: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to unmarshal evidence: %v", err)
	}

	scopes, ok := doc["scopes"].([]interface{})
	if !ok || len(scopes) < 2 {
		t.Fatalf("unexpected scopes structure in evidence JSON")
	}

	// Truncate to 10 scopes
	doc["scopes"] = scopes[:10]
	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated evidence: %v", err)
	}
	if err := os.WriteFile(evidencePath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated evidence: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "security-test-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F4.1 FAIL: security-test-audit succeeded on 10 scopes, expected failure")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "scope count does not match") {
		t.Fatalf("T2.F4.1 FAIL: expected 'scope count does not match' error, got: %s", combined)
	}
	t.Log("T2.F4.1 PASS: security test audit strictly fails closed when scope count is 10 instead of 11")
}

// T2.F4.2: Security test audit fails closed if scope count is 12 instead of 11.
func TestT2_F4_2_SecurityAuditFailsOnTwelveScopes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockSecurityRepo(t)
	evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")

	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read mock evidence: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to unmarshal evidence: %v", err)
	}

	scopes, ok := doc["scopes"].([]interface{})
	if !ok || len(scopes) == 0 {
		t.Fatalf("unexpected scopes in evidence JSON")
	}

	// Append a 12th synthetic scope
	extraScope := map[string]interface{}{
		"boundary_id":          "product-runtime/extra-unregistered-boundary",
		"implementation_state": "implemented",
	}
	doc["scopes"] = append(scopes, extraScope)

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated evidence: %v", err)
	}
	if err := os.WriteFile(evidencePath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated evidence: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "security-test-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F4.2 FAIL: security-test-audit succeeded with 12 scopes, expected failure")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "unregistered scope") && !strings.Contains(combined, "scope count does not match") {
		t.Fatalf("T2.F4.2 FAIL: expected unregistered scope or count mismatch error, got: %s", combined)
	}
	t.Log("T2.F4.2 PASS: security test audit strictly fails closed when scope count is 12 instead of 11")
}

// T2.F4.3: Security test audit detects renamed or mismatched scope key.
func TestT2_F4_3_SecurityAuditDetectsRenamedScopeKey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockSecurityRepo(t)
	evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")

	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read mock evidence: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to unmarshal evidence: %v", err)
	}

	scopes, ok := doc["scopes"].([]interface{})
	if !ok || len(scopes) == 0 {
		t.Fatalf("unexpected scopes structure")
	}

	// Mutate the first scope boundary ID to a mismatched key
	firstScope, ok := scopes[0].(map[string]interface{})
	if !ok {
		t.Fatalf("scope entry is not a map")
	}
	firstScope["boundary_id"] = "product-runtime/mismatched-renamed-key"

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated evidence: %v", err)
	}
	if err := os.WriteFile(evidencePath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated evidence: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "security-test-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F4.3 FAIL: security-test-audit succeeded despite renamed scope key")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "references unregistered scope") && !strings.Contains(combined, "mismatched") {
		t.Fatalf("T2.F4.3 FAIL: expected unregistered scope error, got: %s", combined)
	}
	t.Log("T2.F4.3 PASS: security test audit successfully detected renamed or mismatched scope key")
}

// T2.F4.4: Security test audit detects missing test attestation hash / source_commit.
func TestT2_F4_4_SecurityAuditDetectsMissingAttestationHash(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockSecurityRepo(t)
	evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")

	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read mock evidence: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to unmarshal evidence: %v", err)
	}

	// Blank out source_commit attestation
	doc["source_commit"] = ""

	mutated, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated evidence: %v", err)
	}
	if err := os.WriteFile(evidencePath, mutated, 0644); err != nil {
		t.Fatalf("failed to write mutated evidence: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "security-test-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F4.4 FAIL: security-test-audit succeeded despite missing attestation hash")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "source_commit") {
		t.Fatalf("T2.F4.4 FAIL: expected source_commit attestation error, got: %s", combined)
	}
	t.Log("T2.F4.4 PASS: security test audit successfully detected missing source_commit attestation hash")
}

// T2.F4.5: Security test audit rejects corrupted JSON in evidence file.
func TestT2_F4_5_SecurityAuditRejectsCorruptedJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockSecurityRepo(t)
	evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")

	// Write syntactically corrupted JSON
	corrupted := []byte(`{ "schema_version": 1, "status": "measured", "scopes": [ { "boundary_id": "truncated... `)
	if err := os.WriteFile(evidencePath, corrupted, 0644); err != nil {
		t.Fatalf("failed to write corrupted JSON: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "security-test-audit", "--root", mockRoot)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F4.5 FAIL: security-test-audit succeeded on corrupted JSON")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "JSON") && !strings.Contains(combined, "unexpected") && !strings.Contains(combined, "EOF") {
		t.Fatalf("T2.F4.5 FAIL: expected JSON syntax error, got: %s", combined)
	}
	t.Log("T2.F4.5 PASS: security test audit strictly rejects corrupted JSON evidence file")
}
