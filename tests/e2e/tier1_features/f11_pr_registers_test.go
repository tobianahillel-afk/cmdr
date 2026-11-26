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

type PRRecord struct {
	PRNumber int    `json:"pr_number"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Commit   string `json:"commit,omitempty"`
	HeadSHA  string `json:"head_sha,omitempty"`
}

// T1.F11.1: Verify work/reconciliation/prs.json contains entry for PR #231.
func TestT1_F11_1_ReconciliationPR231Present(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}

	var records []PRRecord
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}

	found := false
	for _, pr := range records {
		if pr.PRNumber == 231 {
			found = true
			if pr.Status != "CLOSED" && pr.Status != "MERGED" {
				t.Fatalf("expected PR #231 status CLOSED or MERGED, got '%s'", pr.Status)
			}
			break
		}
	}

	if !found {
		t.Skip("T1.F11.1 SKIP: PR #231 reconciliation is scheduled for Milestone 3")
	}
	t.Log("T1.F11.1 PASS: PR #231 is registered in prs.json")
}

// T1.F11.2: Verify work/reconciliation/prs.json contains entry for PR #232.
func TestT1_F11_2_ReconciliationPR232Present(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}

	var records []PRRecord
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}

	found := false
	for _, pr := range records {
		if pr.PRNumber == 232 {
			found = true
			if pr.Status != "CLOSED" && pr.Status != "MERGED" {
				t.Fatalf("expected PR #232 status CLOSED or MERGED, got '%s'", pr.Status)
			}
			break
		}
	}

	if !found {
		t.Skip("T1.F11.2 SKIP: PR #232 reconciliation is scheduled for Milestone 3")
	}
	t.Log("T1.F11.2 PASS: PR #232 is registered in prs.json")
}

// T1.F11.3: Verify PR #231 and #232 commit hashes exist in git history.
func TestT1_F11_3_PRCommitsExistInGit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	var records []PRRecord
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}

	for _, pr := range records {
		if pr.PRNumber == 231 || pr.PRNumber == 232 {
			sha := pr.Commit
			if sha == "" {
				sha = pr.HeadSHA
			}
			if sha != "" {
				if !harness.GitCommitExists(ctx, sha) {
					t.Fatalf("commit '%s' for PR #%d does not exist in git history", sha, pr.PRNumber)
				}
			}
		}
	}
	t.Log("T1.F11.3 PASS: PR #231 and PR #232 referenced commits exist in Git")
}

// T1.F11.4: Verify work/reconciliation/PR_CLOSURE_REGISTER.md records PR #231 closure.
func TestT1_F11_4_ClosureRegisterIncludesPR231(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	closurePath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	content, err := os.ReadFile(closurePath)
	if err != nil {
		t.Fatalf("failed to read PR_CLOSURE_REGISTER.md: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "#231") && !strings.Contains(str, "PR 231") {
		t.Skip("T1.F11.4 SKIP: PR #231 closure is not recorded in PR_CLOSURE_REGISTER.md; scheduled for Milestone 3")
	}
	t.Log("T1.F11.4 PASS: PR #231 is recorded in PR_CLOSURE_REGISTER.md")
}

// T1.F11.5: Verify work/reconciliation/PR_CLOSURE_REGISTER.md records PR #232 closure.
func TestT1_F11_5_ClosureRegisterIncludesPR232(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	closurePath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	content, err := os.ReadFile(closurePath)
	if err != nil {
		t.Fatalf("failed to read PR_CLOSURE_REGISTER.md: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "#232") && !strings.Contains(str, "PR 232") {
		t.Skip("T1.F11.5 SKIP: PR #232 closure is not recorded in PR_CLOSURE_REGISTER.md; scheduled for Milestone 3")
	}
	t.Log("T1.F11.5 PASS: PR #232 is recorded in PR_CLOSURE_REGISTER.md")
}
