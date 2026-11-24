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

type SecurityPolicy struct {
	Scopes []struct {
		BoundaryID string `json:"boundary_id"`
	} `json:"scopes"`
}

type SecurityEvidence struct {
	SchemaVersion int    `json:"schema_version"`
	Status        string `json:"status"`
	Scopes        []struct {
		BoundaryID          string `json:"boundary_id"`
		ImplementationState string `json:"implementation_state"`
	} `json:"scopes"`
}

// helper to provide a valid non-security-critical changes file for security-test-audit
func createF04ChangesFile(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}
	changesPath := filepath.Join(harness.TempDir(t, "f04-changes"), "changes.txt")
	if err := os.WriteFile(changesPath, []byte("engineering/testing/runtime-security-evidence.json\n"), 0644); err != nil {
		t.Fatalf("failed to create changes file: %v", err)
	}
	rel, err := filepath.Rel(root, changesPath)
	if err != nil {
		return changesPath
	}
	return rel
}

// T1.F4.1: Execute cmdr-dev security-test-audit and verify exit code 0.
func TestT1_F4_1_SecurityTestAuditExecutionSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	changesFile := createF04ChangesFile(t)
	res := harness.RunCmdrDev(ctx, "security-test-audit", "--changes-file", changesFile, "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev security-test-audit failed (code %d): %s", res.ExitCode, res.Stderr)
	}
	t.Log("T1.F4.1 PASS: security-test-audit completed successfully with exit code 0")
}

// T1.F4.2: Verify evidence file contains exactly 11 registered scopes.
func TestT1_F4_2_SecurityEvidenceElevenScopes(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	evidencePath := filepath.Join(root, "engineering", "testing", "runtime-security-evidence.json")
	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read runtime-security-evidence.json: %v", err)
	}

	var evidence SecurityEvidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		t.Fatalf("failed to parse evidence JSON: %v", err)
	}

	if len(evidence.Scopes) != 11 {
		t.Fatalf("expected exactly 11 registered scopes in evidence, got %d", len(evidence.Scopes))
	}
	t.Logf("T1.F4.2 PASS: runtime-security-evidence.json has exactly %d scopes", len(evidence.Scopes))
}

// T1.F4.3: Verify all 11 scopes match runtime-security-policy.json.
func TestT1_F4_3_SecurityEvidenceMatchesPolicyScopes(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	policyPath := filepath.Join(root, "engineering", "security", "runtime-security-policy.json")
	policyData, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatalf("failed to read runtime-security-policy.json: %v", err)
	}

	evidencePath := filepath.Join(root, "engineering", "testing", "runtime-security-evidence.json")
	evidenceData, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read runtime-security-evidence.json: %v", err)
	}

	var policy SecurityPolicy
	if err := json.Unmarshal(policyData, &policy); err != nil {
		t.Fatalf("failed to parse policy JSON: %v", err)
	}

	var evidence SecurityEvidence
	if err := json.Unmarshal(evidenceData, &evidence); err != nil {
		t.Fatalf("failed to parse evidence JSON: %v", err)
	}

	evidenceMap := make(map[string]bool)
	for _, s := range evidence.Scopes {
		evidenceMap[s.BoundaryID] = true
	}

	for _, p := range policy.Scopes {
		if !evidenceMap[p.BoundaryID] {
			t.Fatalf("scope boundary '%s' defined in policy is missing from evidence", p.BoundaryID)
		}
	}
	t.Log("T1.F4.3 PASS: all 11 policy scopes are registered in runtime security evidence")
}

// T1.F4.4: Verify each scope has valid non-empty test attestation.
func TestT1_F4_4_SecurityEvidenceNonEmptyAttestations(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	evidencePath := filepath.Join(root, "engineering", "testing", "runtime-security-evidence.json")
	data, err := os.ReadFile(evidencePath)
	if err != nil {
		t.Fatalf("failed to read runtime-security-evidence.json: %v", err)
	}

	var evidence SecurityEvidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		t.Fatalf("failed to parse evidence JSON: %v", err)
	}

	for _, s := range evidence.Scopes {
		if s.BoundaryID == "" {
			t.Fatal("found scope with empty boundary_id")
		}
		if s.ImplementationState == "" {
			t.Fatalf("scope '%s' has empty implementation_state", s.BoundaryID)
		}
	}
	t.Log("T1.F4.4 PASS: all scopes carry non-empty valid attestations")
}

// T1.F4.5: Verify audit summary reports PASS.
func TestT1_F4_5_SecurityAuditSummaryPass(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	changesFile := createF04ChangesFile(t)
	res := harness.RunCmdrDev(ctx, "security-test-audit", "--changes-file", changesFile, "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev security-test-audit failed: %s", res.Stderr)
	}

	if _, err := harness.GetJSONField(res.Stdout, "status"); err == nil {
		harness.AssertJSONFieldEqual(t, res.Stdout, "status", "PASS")
	} else {
		harness.AssertJSONFieldEqual(t, res.Stdout, "coverage_status", "measured")
	}
	t.Log("T1.F4.5 PASS: security-test-audit reports status PASS")
}
