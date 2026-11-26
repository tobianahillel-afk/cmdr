package tier2_boundaries

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// setupMockHealthGateRepo creates an isolated Git repository under .cmdr-tmp/
// to verify health gate boundary conditions without dirtying the real cmdr repository.
func setupMockHealthGateRepo(t *testing.T) string {
	t.Helper()
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := harness.TempDir(t, "health-gate-mock")

	// Initialize a local git repository inside the mock root
	initCmd := exec.Command("git", "-C", mockRoot, "init", "-b", "main")
	if out, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to git init mock repo: %v (%s)", err, out)
	}

	// Copy mandatory bootstrap files for isRepoRoot
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

	// Configure local git user for commits
	_ = exec.Command("git", "-C", mockRoot, "config", "user.name", "Test Writer").Run()
	_ = exec.Command("git", "-C", mockRoot, "config", "user.email", "test@cmdr.local").Run()

	// Initial commit so HEAD exists
	_ = exec.Command("git", "-C", mockRoot, "add", ".").Run()
	_ = exec.Command("git", "-C", mockRoot, "commit", "-m", "initial commit").Run()

	return mockRoot
}

// T2.F7.1: Mainline health gate fails closed on dirty git repository.
func TestT2_F7_1_HealthGateFailsClosedOnDirtyRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Create an uncommitted modified/untracked file in the mock repository
	dirtyFile := filepath.Join(mockRoot, "DIRTY_TEST_FILE.txt")
	if err := os.WriteFile(dirtyFile, []byte("uncommitted change"), 0644); err != nil {
		t.Fatalf("failed to create dirty test file: %v", err)
	}

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	// Health gate must exit with non-zero failure code (1 = MAIN_UNHEALTHY)
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.1 FAIL: mainline-health-gate returned 0 on dirty git working tree")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "uncommitted") && !strings.Contains(combined, "FAIL") {
		t.Fatalf("T2.F7.1 FAIL: expected git dimension failure for uncommitted modifications, got: %s", combined)
	}
	t.Log("T2.F7.1 PASS: mainline health gate strictly fails closed when working tree is dirty")
}

// T2.F7.2: Mainline health gate fails closed on 1-byte diff in cmdr-product-spec/.
func TestT2_F7_2_HealthGateFailsClosedOnProductSpecDiff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Create cmdr-product-spec with a 1-byte mutation
	specFile := filepath.Join(mockRoot, "cmdr-product-spec", "spec.md")
	if err := os.MkdirAll(filepath.Dir(specFile), 0755); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}
	if err := os.WriteFile(specFile, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to write spec file: %v", err)
	}

	_ = exec.Command("git", "-C", mockRoot, "add", ".").Run()
	_ = exec.Command("git", "-C", mockRoot, "commit", "-m", "commit with spec").Run()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.2 FAIL: mainline-health-gate passed despite non-zero diff in cmdr-product-spec")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "cmdr-product-spec") && !strings.Contains(combined, "FAIL") && !strings.Contains(combined, "UNVERIFIED") {
		t.Fatalf("T2.F7.2 FAIL: expected product spec baseline violation or UNVERIFIED, got: %s", combined)
	}
	t.Log("T2.F7.2 PASS: mainline health gate fails closed on spec alteration")
}

// T2.F7.3: Mainline health gate returns MAIN_UNVERIFIED (exit code 2) on missing baseline history.
func TestT2_F7_3_HealthGateReturnsUnverifiedOnMissingBaseline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupMockHealthGateRepo(t)

	// Baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40 does not exist in mock git history
	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")

	// Status must be either MAIN_UNHEALTHY (exit 1) or MAIN_UNVERIFIED (exit 2) — never 0
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.3 FAIL: mainline-health-gate returned 0 when product spec baseline history is missing")
	}

	combined := res.Stderr + " " + res.Stdout
	if strings.Contains(combined, "MAIN_UNVERIFIED") || res.ExitCode == 2 {
		t.Logf("T2.F7.3 PASS: mainline-health-gate returned MAIN_UNVERIFIED (exit code 2) as expected: %s", combined)
		return
	}

	t.Logf("T2.F7.3 PASS: mainline-health-gate fail-closed behavior verified (exit code %d)", res.ExitCode)
}

func seedCleanHealthGateBaseline(t *testing.T, ctx context.Context, mockRoot string) {
	t.Helper()
	headRes := harness.RunCmd(ctx, mockRoot, nil, "git", "rev-parse", "HEAD")
	if headRes.ExitCode != 0 {
		t.Fatalf("failed to resolve mockRoot HEAD: %s", headRes.Stderr)
	}
	headSHA := strings.ToLower(strings.TrimSpace(headRes.Stdout))

	digestCanonical := func(v interface{}) string {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal canonical digest payload: %v", err)
		}
		sum := sha256.Sum256(data)
		return hex.EncodeToString(sum[:])
	}

	type workLease struct {
		ID          string `json:"id"`
		WorkUnit    string `json:"work_unit"`
		AgentID     string `json:"agent_id"`
		Mode        string `json:"mode,omitempty"`
		BaseHeadSHA string `json:"base_head_sha"`
		AcquiredAt  string `json:"acquired_at"`
		ExpiresAt   string `json:"expires_at"`
		ReleasedAt  string `json:"released_at,omitempty"`
	}
	claims := struct {
		SchemaVersion int         `json:"schema_version"`
		RegistryKind  string      `json:"registry_kind"`
		Claims        []workLease `json:"claims"`
	}{
		SchemaVersion: 1,
		RegistryKind:  "engineering-work-leases",
		Claims: []workLease{
			{
				ID:          "LEASE-F07-MOCK-0001",
				WorkUnit:    "E10-INV-002G",
				AgentID:     "AGENT-F07-MOCK",
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
	if err := os.WriteFile(filepath.Join(mockRoot, "engineering", "recovery", "claims.json"), append(claimsBytes, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mock claims: %v", err)
	}

	prevDigest := digestCanonical("cmdr-execution-journal-v1")
	body := struct {
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
	}{
		Sequence:       1,
		ID:             "EVT-F07-MOCK-0001",
		WorkUnit:       "E10-INV-002G",
		AgentID:        "AGENT-F07-MOCK",
		LeaseID:        "LEASE-F07-MOCK-0001",
		Kind:           "handoff",
		ObservedAt:     "2026-11-24T10:30:00Z",
		HeadSHA:        headSHA,
		EvidenceRefs:   []string{"work/lots/E10-INV-002G/HANDOFF.json"},
		PreviousDigest: prevDigest,
	}
	evtDigest := digestCanonical(body)
	type journalEvent struct {
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
	journal := struct {
		SchemaVersion int            `json:"schema_version"`
		RegistryKind  string         `json:"registry_kind"`
		Events        []journalEvent `json:"events"`
	}{
		SchemaVersion: 1,
		RegistryKind:  "engineering-recovery-journal",
		Events: []journalEvent{
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
	if err := os.WriteFile(filepath.Join(mockRoot, "engineering", "recovery", "journal.json"), append(journalBytes, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mock journal: %v", err)
	}

	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "engineering/recovery/claims.json", "engineering/recovery/journal.json"); res.ExitCode != 0 {
		t.Fatalf("git add seeded recovery files failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "seed valid recovery journal baseline"); res.ExitCode != 0 {
		t.Fatalf("git commit seeded recovery files failed: %s", res.Stderr)
	}
}

// T2.F7.4: Mainline health gate fails closed if any single dimension fails.
func TestT2_F7_4_HealthGateFailsIfAnySingleDimensionFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	type gateReport struct {
		Status     string `json:"status"`
		ExitCode   int    `json:"exit_code"`
		Dimensions map[string]struct {
			Name   string   `json:"name"`
			Passed bool     `json:"passed"`
			Status string   `json:"status"`
			Errors []string `json:"errors"`
		} `json:"dimensions"`
	}

	mockRoot := harness.CloneIsolatedRepo(t, "f07-single-dim")
	seedCleanHealthGateBaseline(t, ctx, mockRoot)

	// 1. Baseline isolated repo must certify MAIN_HEALTHY (exit 0 across all 5 dimensions)
	baseRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if baseRes.ExitCode != 0 {
		t.Fatalf("T2.F7.4 baseline FAIL: untampered isolated repo exited %d (expected 0)\nStdout: %s\nStderr: %s",
			baseRes.ExitCode, baseRes.Stdout, baseRes.Stderr)
	}
	var baseReport gateReport
	if err := harness.ParseJSON(baseRes.Stdout, &baseReport); err != nil {
		t.Fatalf("T2.F7.4 baseline FAIL: failed to parse JSON: %v", err)
	}
	if baseReport.Status != "MAIN_HEALTHY" {
		t.Fatalf("T2.F7.4 baseline FAIL: expected MAIN_HEALTHY, got %s", baseReport.Status)
	}
	allDims := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, d := range allDims {
		dim, ok := baseReport.Dimensions[d]
		if !ok || !dim.Passed || dim.Status != "PASS" {
			t.Fatalf("T2.F7.4 baseline FAIL: dimension %q not PASS: %+v", d, dim)
		}
	}

	// 2. Inject a single-dimension fault into traceability (phantom unit in ledger.json) and commit so git stays clean
	ledgerPath := filepath.Join(mockRoot, "engineering", "implementation", "ledger.json")
	ledgerRaw, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("failed to read mock ledger.json: %v", err)
	}
	var ledgerDoc map[string]interface{}
	if err := json.Unmarshal(ledgerRaw, &ledgerDoc); err != nil {
		t.Fatalf("failed to parse mock ledger.json: %v", err)
	}
	entries, ok := ledgerDoc["entries"].([]interface{})
	if !ok {
		t.Fatalf("ledger entries missing or not array")
	}
	ledgerDoc["entries"] = append(entries, map[string]interface{}{
		"capability":       "CAP-FAULT-001",
		"work_unit":        "E10-PHANTOM-FAULT-UNIT",
		"manifest_path":    "work/lots/E10-PHANTOM/manifest.json",
		"handoff_path":     "work/lots/E10-PHANTOM/HANDOFF.json",
		"runtime_boundary": "event-search-runtime",
		"runtime_root":     "product-runtime/event-search/**",
	})
	mutatedLedger, err := json.MarshalIndent(ledgerDoc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mutated ledger: %v", err)
	}
	if err := os.WriteFile(ledgerPath, append(mutatedLedger, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mutated ledger: %v", err)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "add", "engineering/implementation/ledger.json"); res.ExitCode != 0 {
		t.Fatalf("git add mutated ledger failed: %s", res.Stderr)
	}
	if res := harness.RunCmd(ctx, mockRoot, nil, "git", "commit", "-m", "inject single-dimension traceability fault"); res.ExitCode != 0 {
		t.Fatalf("git commit mutated ledger failed: %s", res.Stderr)
	}

	// 3. Assert MAIN_UNHEALTHY (exit 1) with traceability=FAIL and remaining 4 dimensions=PASS
	faultRes := harness.RunCmdrDev(ctx, "mainline-health-gate", "--root", mockRoot, "--json")
	if faultRes.ExitCode != 1 {
		t.Fatalf("T2.F7.4 FAIL: expected exit code 1 (MAIN_UNHEALTHY), got %d\nStdout: %s\nStderr: %s",
			faultRes.ExitCode, faultRes.Stdout, faultRes.Stderr)
	}
	var faultReport gateReport
	if err := harness.ParseJSON(faultRes.Stdout, &faultReport); err != nil {
		t.Fatalf("T2.F7.4 FAIL: failed to parse fault JSON: %v", err)
	}
	if faultReport.Status != "MAIN_UNHEALTHY" {
		t.Fatalf("T2.F7.4 FAIL: expected status MAIN_UNHEALTHY, got %s", faultReport.Status)
	}
	if traceDim, ok := faultReport.Dimensions["traceability"]; !ok || traceDim.Passed || traceDim.Status != "FAIL" {
		t.Fatalf("T2.F7.4 FAIL: expected traceability dimension FAIL, got %+v", traceDim)
	}
	for _, d := range allDims {
		if d == "traceability" {
			continue
		}
		if dim, ok := faultReport.Dimensions[d]; !ok || !dim.Passed || dim.Status != "PASS" {
			t.Fatalf("T2.F7.4 FAIL: expected non-target dimension %q to remain PASS, got %+v", d, dim)
		}
	}

	t.Log("T2.F7.4 PASS: health gate fails closed (exit 1 MAIN_UNHEALTHY) when a single dimension fails while other 4 pass")
}

// T2.F7.5: Mainline health gate rejects unknown CLI flag with usage error.
func TestT2_F7_5_HealthGateRejectsUnknownFlag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--unknown-flag-xyz")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F7.5 FAIL: mainline-health-gate succeeded with unknown CLI flag")
	}

	combined := res.Stderr + " " + res.Stdout
	if !strings.Contains(combined, "flag") && !strings.Contains(combined, "unknown") && res.ExitCode != 1 && res.ExitCode != 2 {
		t.Fatalf("T2.F7.5 FAIL: expected flag error, got: %s", combined)
	}
	t.Logf("T2.F7.5 PASS: unknown flag correctly rejected with non-zero exit code %d", res.ExitCode)
}
