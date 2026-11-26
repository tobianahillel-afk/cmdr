package tier4_workflows

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type w4HealthGateReport struct {
	Status     string `json:"status"`
	ExitCode   int    `json:"exit_code"`
	Dimensions map[string]struct {
		Name   string   `json:"name"`
		Passed bool     `json:"passed"`
		Status string   `json:"status"`
		Errors []string `json:"errors"`
	} `json:"dimensions"`
}

// Mirrors tools/cmdr-dev/recovery_journal.go:49-60 struct field order for canonical SHA-256 hashing.
type w4RecoveryJournalEventBody struct {
	Sequence       int      `json:"sequence"`
	ID             string   `json:"id"`
	WorkUnit       string   `json:"work_unit"`
	AgentID        string   `json:"agent_id"`
	LeaseID        string   `json:"lease_id"`
	Kind           string   `json:"kind"`
	ObservedAt     string   `json:"observed_at"`
	HeadSHA        string   `json:"head_sha"`
	EvidenceRefs   []string `json:"evidence_refs"`
	PreviousDigest string   `json:"previous_digest"`
}

type w4RecoveryJournalEvent struct {
	Sequence       int      `json:"sequence"`
	ID             string   `json:"id"`
	WorkUnit       string   `json:"work_unit"`
	AgentID        string   `json:"agent_id"`
	LeaseID        string   `json:"lease_id"`
	Kind           string   `json:"kind"`
	ObservedAt     string   `json:"observed_at"`
	HeadSHA        string   `json:"head_sha"`
	EvidenceRefs   []string `json:"evidence_refs"`
	PreviousDigest string   `json:"previous_digest"`
	Digest         string   `json:"digest"`
}

type w4RecoveryJournal struct {
	SchemaVersion int                      `json:"schema_version"`
	RegistryKind  string                   `json:"registry_kind"`
	Events        []w4RecoveryJournalEvent `json:"events"`
}

type w4WorkLease struct {
	ID          string `json:"id"`
	WorkUnit    string `json:"work_unit"`
	AgentID     string `json:"agent_id"`
	Mode        string `json:"mode,omitempty"`
	BaseHeadSHA string `json:"base_head_sha"`
	AcquiredAt  string `json:"acquired_at"`
	ExpiresAt   string `json:"expires_at"`
	ReleasedAt  string `json:"released_at,omitempty"`
}

type w4WorkLeaseRegistry struct {
	SchemaVersion int           `json:"schema_version"`
	RegistryKind  string        `json:"registry_kind"`
	Claims        []w4WorkLease `json:"claims"`
}

func w4DigestCanonical(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func setupW4MockRepo(t *testing.T) string {
	t.Helper()
	mockRoot := harness.CloneIsolatedRepo(t, "w4-tamper")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headRes := harness.RunCmd(ctx, mockRoot, nil, "git", "rev-parse", "HEAD")
	if headRes.ExitCode != 0 {
		t.Fatalf("failed to resolve mockRoot HEAD: %s", headRes.Stderr)
	}
	headSHA := strings.ToLower(strings.TrimSpace(headRes.Stdout))

	// Seed valid released lease matching the baseline journal event
	claims := w4WorkLeaseRegistry{
		SchemaVersion: 1,
		RegistryKind:  "engineering-work-leases",
		Claims: []w4WorkLease{
			{
				ID:          "LEASE-W4-MOCK-0001",
				WorkUnit:    "E10-INV-002G",
				AgentID:     "AGENT-W4-MOCK",
				Mode:        "mutating",
				BaseHeadSHA: headSHA,
				AcquiredAt:  "2026-11-24T10:00:00Z",
				ExpiresAt:   "2026-11-24T11:30:00Z",
				ReleasedAt:  "2026-11-24T11:00:00Z",
			},
		},
	}
	claimsBytes, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock claims: %v", err)
	}
	claimsPath := filepath.Join(mockRoot, "engineering", "recovery", "claims.json")
	if err := os.WriteFile(claimsPath, append(claimsBytes, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mock claims: %v", err)
	}

	// Seed valid recovery journal event with real SHA-256 chain digest
	prevDigest := w4DigestCanonical("cmdr-execution-journal-v1")
	body := w4RecoveryJournalEventBody{
		Sequence:       1,
		ID:             "EVT-W4-MOCK-0001",
		WorkUnit:       "E10-INV-002G",
		AgentID:        "AGENT-W4-MOCK",
		LeaseID:        "LEASE-W4-MOCK-0001",
		Kind:           "handoff",
		ObservedAt:     "2026-11-24T10:30:00Z",
		HeadSHA:        headSHA,
		EvidenceRefs:   []string{"work/lots/E10-INV-002G/HANDOFF.json"},
		PreviousDigest: prevDigest,
	}
	evtDigest := w4DigestCanonical(body)
	journal := w4RecoveryJournal{
		SchemaVersion: 1,
		RegistryKind:  "engineering-recovery-journal",
		Events: []w4RecoveryJournalEvent{
			{
				Sequence:       body.Sequence,
				ID:             body.ID,
				WorkUnit:       body.WorkUnit,
				AgentID:        body.AgentID,
				LeaseID:        body.LeaseID,
				Kind:           body.Kind,
				ObservedAt:     body.ObservedAt,
				HeadSHA:        body.HeadSHA,
				EvidenceRefs:   body.EvidenceRefs,
				PreviousDigest: body.PreviousDigest,
				Digest:         evtDigest,
			},
		},
	}
	journalBytes, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock journal: %v", err)
	}
	journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")
	if err := os.WriteFile(journalPath, append(journalBytes, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mock journal: %v", err)
	}

	addRes := harness.RunCmd(ctx, mockRoot, nil, "git", "add",
		"engineering/recovery/claims.json", "engineering/recovery/journal.json")
	if addRes.ExitCode != 0 {
		t.Fatalf("git add seeded recovery files failed: %s", addRes.Stderr)
	}
	commitRes := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "seed valid W4 recovery journal baseline")
	if commitRes.ExitCode != 0 {
		t.Fatalf("git commit seeded recovery files failed: %s", commitRes.Stderr)
	}

	return mockRoot
}

func runW4HealthGate(t *testing.T, ctx context.Context, mockRoot string) (harness.CmdResult, w4HealthGateReport) {
	t.Helper()
	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	var report w4HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse mainline-health-gate JSON report: %v\nStdout: %s\nStderr: %s", err, res.Stdout, res.Stderr)
	}
	return res, report
}

func assertOtherDimensionsHealthy(t *testing.T, report w4HealthGateReport, expectedFailDim string) {
	t.Helper()
	allDims := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, dimName := range allDims {
		dim, ok := report.Dimensions[dimName]
		if !ok {
			t.Fatalf("health gate report missing dimension %q", dimName)
		}
		if dimName == expectedFailDim {
			continue
		}
		if !dim.Passed || dim.Status != "PASS" {
			t.Fatalf("expected non-target dimension %q to remain PASS, got status=%s errors=%v", dimName, dim.Status, dim.Errors)
		}
	}
}

// TestT4_04_TamperDetectionAndFailClosed:
// Verifies untampered baseline passes MAIN_HEALTHY (exit 0), then injects each of the
// four tamper vectors into isolated mock repositories and asserts fail-closed rejection
// in the exact expected health-gate dimension.
func TestT4_04_TamperDetectionAndFailClosed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	t.Run("Baseline_Untampered_MainHealthy", func(t *testing.T) {
		mockRoot := setupW4MockRepo(t)
		gateRes, report := runW4HealthGate(t, ctx, mockRoot)
		if gateRes.ExitCode != 0 {
			t.Fatalf("T4.04 baseline FAIL: untampered mock repo exited %d (expected 0 MAIN_HEALTHY)\nStdout: %s\nStderr: %s",
				gateRes.ExitCode, gateRes.Stdout, gateRes.Stderr)
		}
		if report.Status != "MAIN_HEALTHY" || report.ExitCode != 0 {
			t.Fatalf("T4.04 baseline FAIL: expected status=MAIN_HEALTHY exit_code=0, got status=%s exit_code=%d",
				report.Status, report.ExitCode)
		}
		assertOtherDimensionsHealthy(t, report, "")
		t.Log("T4.04 Baseline PASS: untampered W4 mock repo certified MAIN_HEALTHY (exit 0 across all 5 dimensions)")
	})

	t.Run("Vector1_RecoveryJournalDigestTamper", func(t *testing.T) {
		mockRoot := setupW4MockRepo(t)
		journalPath := filepath.Join(mockRoot, "engineering", "recovery", "journal.json")
		journalData, err := os.ReadFile(journalPath)
		if err != nil {
			t.Fatalf("failed to read mock journal: %v", err)
		}
		var journalDoc map[string]interface{}
		if err := json.Unmarshal(journalData, &journalDoc); err != nil {
			t.Fatalf("failed to parse journal: %v", err)
		}
		events, ok := journalDoc["events"].([]interface{})
		if !ok || len(events) == 0 {
			t.Fatalf("T4.04 Vector 1 setup error: expected seeded journal events, got len=0")
		}
		firstEvt, ok := events[0].(map[string]interface{})
		if !ok {
			t.Fatalf("journal event[0] is not a JSON object")
		}
		firstEvt["digest"] = "deadbeefbadc0ffee000000000000000000000000000000000000000000000000"
		corruptedJournal, err := json.MarshalIndent(journalDoc, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal corrupted journal: %v", err)
		}
		if err := os.WriteFile(journalPath, corruptedJournal, 0644); err != nil {
			t.Fatalf("failed to write corrupted journal: %v", err)
		}
		// Commit tampered journal inside mockRoot so git working tree stays clean and operational_quality is isolated
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "engineering/recovery/journal.json"); res.ExitCode != 0 {
			t.Fatalf("git add tampered journal failed: %s", res.Stderr)
		}
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "tamper recovery journal digest"); res.ExitCode != 0 {
			t.Fatalf("git commit tampered journal failed: %s", res.Stderr)
		}

		gateRes, report := runW4HealthGate(t, ctx, mockRoot)
		if gateRes.ExitCode == 0 {
			t.Fatalf("T4.04 Vector 1 FAIL: mainline-health-gate returned exit 0 on tampered journal digest")
		}
		if gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
			t.Fatalf("T4.04 Vector 1 FAIL: unexpected exit code %d; expected 1 or 2", gateRes.ExitCode)
		}
		if report.Status == "MAIN_HEALTHY" {
			t.Fatalf("T4.04 Vector 1 FAIL: status reported MAIN_HEALTHY")
		}
		dim, ok := report.Dimensions["operational_quality"]
		if !ok || dim.Passed || dim.Status != "FAIL" {
			t.Fatalf("T4.04 Vector 1 FAIL: expected operational_quality FAIL, got %+v", dim)
		}
		joinedErrs := strings.Join(dim.Errors, " | ")
		if !strings.Contains(joinedErrs, "digest mismatch") {
			t.Fatalf("T4.04 Vector 1 FAIL: expected 'digest mismatch' in operational_quality errors, got: %s", joinedErrs)
		}
		assertOtherDimensionsHealthy(t, report, "operational_quality")
		t.Logf("T4.04 Vector 1 PASS: operational_quality failed closed on journal digest mismatch (exit=%d)", gateRes.ExitCode)
	})

	t.Run("Vector2_ImplementationLedgerPhantomUnit", func(t *testing.T) {
		mockRoot := setupW4MockRepo(t)
		ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")
		ledgerData, err := os.ReadFile(ledgerPath)
		if err != nil {
			t.Fatalf("failed to read mock ledger: %v", err)
		}
		var ledgerDoc map[string]interface{}
		if err := json.Unmarshal(ledgerData, &ledgerDoc); err != nil {
			t.Fatalf("failed to parse ledger: %v", err)
		}
		entries, ok := ledgerDoc["entries"].([]interface{})
		if !ok {
			t.Fatalf("ledger entries is not an array")
		}
		ledgerDoc["entries"] = append(entries, map[string]interface{}{
			"capability":       "CAP-TAMPER-001",
			"work_unit":        "E10-PHANTOM-UNIT-NONEXISTENT",
			"manifest_path":    "work/lots/E10-PHANTOM/manifest.json",
			"handoff_path":     "work/lots/E10-PHANTOM/HANDOFF.json",
			"runtime_boundary": "event-search-runtime",
			"runtime_root":     "product-runtime/event-search/**",
		})
		corruptedLedger, err := json.MarshalIndent(ledgerDoc, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal corrupted ledger: %v", err)
		}
		if err := os.WriteFile(ledgerPath, corruptedLedger, 0644); err != nil {
			t.Fatalf("failed to write corrupted ledger: %v", err)
		}
		// Commit tampered ledger inside mockRoot so git working tree stays clean and traceability is isolated
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "engineering/implementation/ledger.json"); res.ExitCode != 0 {
			t.Fatalf("git add tampered ledger failed: %s", res.Stderr)
		}
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "tamper implementation ledger phantom unit"); res.ExitCode != 0 {
			t.Fatalf("git commit tampered ledger failed: %s", res.Stderr)
		}

		gateRes, report := runW4HealthGate(t, ctx, mockRoot)
		if gateRes.ExitCode == 0 {
			t.Fatalf("T4.04 Vector 2 FAIL: mainline-health-gate returned exit 0 on tampered ledger")
		}
		if gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
			t.Fatalf("T4.04 Vector 2 FAIL: unexpected exit code %d; expected 1 or 2", gateRes.ExitCode)
		}
		if report.Status == "MAIN_HEALTHY" {
			t.Fatalf("T4.04 Vector 2 FAIL: status reported MAIN_HEALTHY")
		}
		dim, ok := report.Dimensions["traceability"]
		if !ok || dim.Passed || dim.Status != "FAIL" {
			t.Fatalf("T4.04 Vector 2 FAIL: expected traceability FAIL, got %+v", dim)
		}
		joinedErrs := strings.Join(dim.Errors, " | ")
		if !strings.Contains(joinedErrs, "implementation ledger audit failed") {
			t.Fatalf("T4.04 Vector 2 FAIL: expected 'implementation ledger audit failed' in traceability errors, got: %s", joinedErrs)
		}
		assertOtherDimensionsHealthy(t, report, "traceability")
		t.Logf("T4.04 Vector 2 PASS: traceability failed closed on tampered ledger (exit=%d)", gateRes.ExitCode)
	})

	t.Run("Vector3_UntrackedSentinelFile", func(t *testing.T) {
		mockRoot := setupW4MockRepo(t)
		sentinelPath := filepath.Join(mockRoot, ".tamper_sentinel.tmp")
		if err := os.WriteFile(sentinelPath, []byte("adversarial-tamper-probe"), 0644); err != nil {
			t.Fatalf("failed to create sentinel file: %v", err)
		}

		gateRes, report := runW4HealthGate(t, ctx, mockRoot)
		if gateRes.ExitCode == 0 {
			t.Fatalf("T4.04 Vector 3 FAIL: mainline-health-gate returned exit 0 with untracked sentinel file")
		}
		if gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
			t.Fatalf("T4.04 Vector 3 FAIL: unexpected exit code %d; expected 1 or 2", gateRes.ExitCode)
		}
		if report.Status == "MAIN_HEALTHY" {
			t.Fatalf("T4.04 Vector 3 FAIL: status reported MAIN_HEALTHY")
		}
		dim, ok := report.Dimensions["git"]
		if !ok || dim.Passed || dim.Status != "FAIL" {
			t.Fatalf("T4.04 Vector 3 FAIL: expected git FAIL, got %+v", dim)
		}
		joinedErrs := strings.Join(dim.Errors, " | ")
		if !strings.Contains(joinedErrs, "working tree has uncommitted modifications") || !strings.Contains(joinedErrs, ".tamper_sentinel.tmp") {
			t.Fatalf("T4.04 Vector 3 FAIL: expected uncommitted '.tamper_sentinel.tmp' error in git dimension, got: %s", joinedErrs)
		}
		assertOtherDimensionsHealthy(t, report, "git")
		t.Logf("T4.04 Vector 3 PASS: git failed closed on untracked sentinel file (exit=%d)", gateRes.ExitCode)
	})

	t.Run("Vector4_ProductSpecMutation", func(t *testing.T) {
		mockRoot := setupW4MockRepo(t)
		specReadmePath := filepath.Join(mockRoot, "cmdr-product-spec", "README.md")
		if err := os.WriteFile(specReadmePath, []byte("# TAMPERED PRODUCT SPEC INVARIANT\n"), 0644); err != nil {
			t.Fatalf("failed to alter product spec invariant: %v", err)
		}
		// Commit spec tampering inside mockRoot so working tree is clean and git baseline diff invariant is evaluated
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "cmdr-product-spec/README.md"); res.ExitCode != 0 {
			t.Fatalf("git add tampered spec failed: %s", res.Stderr)
		}
		if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "tamper product spec invariant"); res.ExitCode != 0 {
			t.Fatalf("git commit tampered spec failed: %s", res.Stderr)
		}

		gateRes, report := runW4HealthGate(t, ctx, mockRoot)
		if gateRes.ExitCode == 0 {
			t.Fatalf("T4.04 Vector 4 FAIL: mainline-health-gate returned exit 0 on altered product spec")
		}
		if gateRes.ExitCode != 1 && gateRes.ExitCode != 2 {
			t.Fatalf("T4.04 Vector 4 FAIL: unexpected exit code %d; expected 1 or 2", gateRes.ExitCode)
		}
		if report.Status == "MAIN_HEALTHY" {
			t.Fatalf("T4.04 Vector 4 FAIL: status reported MAIN_HEALTHY")
		}
		dim, ok := report.Dimensions["git"]
		if !ok || dim.Passed || dim.Status != "FAIL" {
			t.Fatalf("T4.04 Vector 4 FAIL: expected git FAIL, got %+v", dim)
		}
		joinedErrs := strings.Join(dim.Errors, " | ")
		if !strings.Contains(joinedErrs, "cmdr-product-spec has non-zero diff vs baseline commit") {
			t.Fatalf("T4.04 Vector 4 FAIL: expected product spec baseline diff error in git dimension, got: %s", joinedErrs)
		}
		assertOtherDimensionsHealthy(t, report, "git")
		t.Logf("T4.04 Vector 4 PASS: git failed closed on product spec baseline diff (exit=%d)", gateRes.ExitCode)
	})
}
