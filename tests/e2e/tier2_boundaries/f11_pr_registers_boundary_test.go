package tier2_boundaries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

func setupIsolatedPRRegisterRepo(t *testing.T, ctx context.Context, prefix string) string {
	t.Helper()
	mockRoot := harness.CloneIsolatedRepo(t, prefix)
	if err := harness.ValidatePRRegisterConsistency(ctx, mockRoot); err != nil {
		t.Fatalf("pre-mutation PR register baseline failed in %s: %v", prefix, err)
	}
	return mockRoot
}

func mutatePRRegisterJSON(t *testing.T, root string, mutate func(records []harness.PRMetadata) []harness.PRMetadata) {
	t.Helper()
	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}
	var records []harness.PRMetadata
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}
	updated := mutate(records)
	out, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal prs.json: %v", err)
	}
	if err := os.WriteFile(prsPath, append(out, '\n'), 0644); err != nil {
		t.Fatalf("failed to write mutated prs.json: %v", err)
	}
}

func appendPRToClosureMarkdown(t *testing.T, root string, prNumber int) {
	t.Helper()
	mdPath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read PR_CLOSURE_REGISTER.md: %v", err)
	}
	line := fmt.Sprintf("\n| #%d | test | MERGED | INCORPORATED_IN_MAIN |\n", prNumber)
	if err := os.WriteFile(mdPath, append(data, []byte(line)...), 0644); err != nil {
		t.Fatalf("failed to update PR_CLOSURE_REGISTER.md: %v", err)
	}
}

// T2.F11.1: PR register validation fails if PR #231 commit hash is invalid.
func TestT2_F11_1_PRRegisterFailsOnInvalidPR231Commit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupIsolatedPRRegisterRepo(t, ctx, "f11-pr231")
	mutatePRRegisterJSON(t, mockRoot, func(records []harness.PRMetadata) []harness.PRMetadata {
		for i := range records {
			if records[i].PRNumber == 231 {
				records[i].Commit = "0000000000000000000000000000000000000000"
				records[i].MergeCommit = ""
				records[i].HeadSHA = ""
				records[i].SHA = ""
				return records
			}
		}
		return append(records, harness.PRMetadata{
			PRNumber:             231,
			Title:                "reconcile: CAP-STD-004 ledger and graph mapping",
			Status:               "MERGED",
			Disposition:          "INCORPORATED_IN_MAIN",
			DispositionRationale: "Reconciled CAP-STD-004",
			Commit:               "0000000000000000000000000000000000000000",
		})
	})
	appendPRToClosureMarkdown(t, mockRoot, 231)

	err := harness.ValidatePRRegisterConsistency(ctx, mockRoot)
	if err == nil {
		t.Fatalf("T2.F11.1 FAIL: ValidatePRRegisterConsistency passed with non-existent commit hash for PR #231")
	}
	if !strings.Contains(err.Error(), "non-existent commit") {
		t.Fatalf("T2.F11.1 FAIL: expected non-existent commit error, got: %v", err)
	}

	ghClient, ghErr := harness.NewGitHubClient()
	if ghErr != nil {
		t.Fatalf("NewGitHubClient failed: %v", ghErr)
	}
	emptyBin := harness.TempDir(t, "f11-empty-gh")
	t.Setenv("PATH", emptyBin)
	codexRes, codexErr := ghClient.VerifyCodexPRResolution(ctx, 231, "0000000000000000000000000000000000000000")
	if !errors.Is(codexErr, harness.ErrUnverified) || codexRes.Verified {
		t.Fatalf("T2.F11.1 FAIL: expected VerifyCodexPRResolution to reject unattributed PR #231, got res=%+v err=%v", codexRes, codexErr)
	}

	t.Log("T2.F11.1 PASS: PR register validation strictly fails closed when PR #231 commit hash is invalid")
}

// T2.F11.2: PR register validation fails if PR #232 commit hash is invalid.
func TestT2_F11_2_PRRegisterFailsOnInvalidPR232Commit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupIsolatedPRRegisterRepo(t, ctx, "f11-pr232")
	mutatePRRegisterJSON(t, mockRoot, func(records []harness.PRMetadata) []harness.PRMetadata {
		for i := range records {
			if records[i].PRNumber == 232 {
				records[i].Commit = "not-a-sha-hash"
				records[i].MergeCommit = ""
				records[i].HeadSHA = ""
				records[i].SHA = ""
				return records
			}
		}
		return append(records, harness.PRMetadata{
			PRNumber:             232,
			Title:                "docs: continuous development protocol and CAP-STD-005 roadmap",
			Status:               "MERGED",
			Disposition:          "INCORPORATED_IN_MAIN",
			DispositionRationale: "Documented roadmap",
			Commit:               "not-a-sha-hash",
		})
	})
	appendPRToClosureMarkdown(t, mockRoot, 232)

	err := harness.ValidatePRRegisterConsistency(ctx, mockRoot)
	if err == nil {
		t.Fatalf("T2.F11.2 FAIL: ValidatePRRegisterConsistency passed with malformed commit SHA for PR #232")
	}
	if !strings.Contains(err.Error(), "malformed commit SHA format") {
		t.Fatalf("T2.F11.2 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.2 PASS: PR register validation strictly fails closed when PR #232 commit hash is malformed")
}

// T2.F11.3: PR register validation detects missing closure metadata.
func TestT2_F11_3_PRRegisterDetectsMissingClosureMetadata(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupIsolatedPRRegisterRepo(t, ctx, "f11-metadata")
	mutatePRRegisterJSON(t, mockRoot, func(records []harness.PRMetadata) []harness.PRMetadata {
		for i := range records {
			if records[i].PRNumber == 230 {
				records[i].Disposition = ""
				records[i].DispositionRationale = ""
			}
		}
		return records
	})

	err := harness.ValidatePRRegisterConsistency(ctx, mockRoot)
	if err == nil {
		t.Fatalf("T2.F11.3 FAIL: ValidatePRRegisterConsistency passed despite missing closure disposition on PR #230")
	}
	if !strings.Contains(err.Error(), "missing disposition") {
		t.Fatalf("T2.F11.3 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.3 PASS: PR register validation strictly flags missing closure metadata")
}

// T2.F11.4: PR register validation rejects contradictory merge/closure statuses.
func TestT2_F11_4_PRRegisterRejectsUnmergedPRRegisteredAsClosed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mockRoot := setupIsolatedPRRegisterRepo(t, ctx, "f11-contradict")
	mutatePRRegisterJSON(t, mockRoot, func(records []harness.PRMetadata) []harness.PRMetadata {
		for i := range records {
			if records[i].PRNumber == 230 {
				records[i].Status = "MERGED"
				records[i].Disposition = "REJECTED"
			}
		}
		return records
	})

	err := harness.ValidatePRRegisterConsistency(ctx, mockRoot)
	if err == nil {
		t.Fatalf("T2.F11.4 FAIL: ValidatePRRegisterConsistency passed with MERGED+REJECTED contradiction")
	}
	if !strings.Contains(err.Error(), "contradictory status") {
		t.Fatalf("T2.F11.4 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.4 PASS: PR register validation strictly rejects contradictory merge/closure statuses")
}

// T2.F11.5: PR register validation verifies consistency between prs.json and PR_CLOSURE_REGISTER.md.
func TestT2_F11_5_PRRegisterConsistencyBetweenJSONAndMarkdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}

	var records []harness.PRMetadata
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("T2.F11.5 FAIL: prs.json contains zero records")
	}

	closurePath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	closureContent, err := os.ReadFile(closurePath)
	if err != nil {
		t.Fatalf("failed to read PR_CLOSURE_REGISTER.md: %v", err)
	}
	closureText := string(closureContent)

	missingInMarkdown := 0
	for _, pr := range records {
		if !harness.ClosureRegisterContainsPR(closureText, pr.PRNumber) {
			missingInMarkdown++
		}
	}
	if missingInMarkdown != 0 {
		t.Fatalf("T2.F11.5 FAIL: %d PR record(s) from prs.json missing in PR_CLOSURE_REGISTER.md", missingInMarkdown)
	}

	if err := harness.ValidatePRRegisterConsistency(ctx, root); err != nil {
		t.Fatalf("T2.F11.5 FAIL: ValidatePRRegisterConsistency failed on real repository: %v", err)
	}

	// Verify fault injection: removing #230 from PR_CLOSURE_REGISTER.md in an isolated copy fails closed
	mockRoot := setupIsolatedPRRegisterRepo(t, ctx, "f11-md-desync")
	mockMDPath := filepath.Join(mockRoot, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	mockMDBytes, err := os.ReadFile(mockMDPath)
	if err != nil {
		t.Fatalf("failed to read mock closure register: %v", err)
	}
	desyncedMD := strings.ReplaceAll(strings.ReplaceAll(string(mockMDBytes), "#230", "#000"), "PR 230", "PR 000")
	if err := os.WriteFile(mockMDPath, []byte(desyncedMD), 0644); err != nil {
		t.Fatalf("failed to write desynced closure register: %v", err)
	}
	if err := harness.ValidatePRRegisterConsistency(ctx, mockRoot); err == nil {
		t.Fatalf("T2.F11.5 FAIL: ValidatePRRegisterConsistency passed despite missing PR #230 in PR_CLOSURE_REGISTER.md")
	}

	t.Logf("T2.F11.5 PASS: evaluated %d PR records across prs.json and PR_CLOSURE_REGISTER.md (missing in register: %d)",
		len(records), missingInMarkdown)
}
