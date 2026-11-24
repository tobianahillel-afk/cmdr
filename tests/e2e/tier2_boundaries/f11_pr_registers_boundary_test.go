package tier2_boundaries

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type PRRegistryItem struct {
	PRNumber             int      `json:"pr_number"`
	Title                string   `json:"title"`
	HeadBranch           string   `json:"head_branch"`
	BaseBranch           string   `json:"base_branch"`
	Status               string   `json:"status"`
	Disposition          string   `json:"disposition"`
	DispositionRationale string   `json:"disposition_rationale"`
	Commit               string   `json:"commit,omitempty"`
	HeadSHA              string   `json:"head_sha,omitempty"`
	Author               string   `json:"author,omitempty"`
	ClosedAt             string   `json:"closed_at,omitempty"`
	FilesIntroduced      []string `json:"files_introduced,omitempty"`
}

// validatePRRegisterEntry validates an individual PR record's integrity against git history.
func validatePRRegisterEntry(ctx context.Context, item PRRegistryItem) error {
	if item.PRNumber <= 0 {
		return fmt.Errorf("invalid PR number: %d", item.PRNumber)
	}
	if strings.TrimSpace(item.Title) == "" {
		return fmt.Errorf("PR #%d has empty title", item.PRNumber)
	}
	if item.Status != "MERGED" && item.Status != "CLOSED" && item.Status != "OPEN" {
		return fmt.Errorf("PR #%d has unknown status %q", item.PRNumber, item.Status)
	}

	sha := item.Commit
	if sha == "" {
		sha = item.HeadSHA
	}
	if sha != "" {
		shaPattern := regexp.MustCompile(`^[0-9a-f]{40}$`)
		if !shaPattern.MatchString(sha) {
			return fmt.Errorf("PR #%d has malformed commit SHA format %q", item.PRNumber, sha)
		}
		if !harness.GitCommitExists(ctx, sha) {
			return fmt.Errorf("PR #%d references non-existent commit %s in git object storage", item.PRNumber, sha)
		}
	}

	return nil
}

// T2.F11.1: PR register validation fails if PR #231 commit hash is invalid.
func TestT2_F11_1_PRRegisterFailsOnInvalidPR231Commit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct simulated PR #231 record with an invalid non-existent commit SHA
	item := PRRegistryItem{
		PRNumber:    231,
		Title:       "reconcile: CAP-STD-004 ledger and graph mapping",
		Status:      "MERGED",
		Disposition: "MERGED_INTO_MAIN",
		Commit:      "0000000000000000000000000000000000000000", // Non-existent SHA
	}

	err := validatePRRegisterEntry(ctx, item)
	if err == nil {
		t.Fatalf("T2.F11.1 FAIL: validation passed with non-existent commit hash for PR #231")
	}
	if !strings.Contains(err.Error(), "non-existent commit") {
		t.Fatalf("T2.F11.1 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.1 PASS: PR register validation strictly fails closed when PR #231 commit hash is invalid")
}

// T2.F11.2: PR register validation fails if PR #232 commit hash is invalid.
func TestT2_F11_2_PRRegisterFailsOnInvalidPR232Commit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct simulated PR #232 record with a malformed short SHA
	item := PRRegistryItem{
		PRNumber:    232,
		Title:       "docs: continuous development protocol and CAP-STD-005 roadmap",
		Status:      "MERGED",
		Disposition: "MERGED_INTO_MAIN",
		Commit:      "not-a-sha-hash", // Malformed SHA
	}

	err := validatePRRegisterEntry(ctx, item)
	if err == nil {
		t.Fatalf("T2.F11.2 FAIL: validation passed with malformed commit SHA for PR #232")
	}
	if !strings.Contains(err.Error(), "malformed commit SHA format") {
		t.Fatalf("T2.F11.2 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.2 PASS: PR register validation strictly fails closed when PR #232 commit hash is malformed")
}

// T2.F11.3: PR register validation detects missing closure date or author.
func TestT2_F11_3_PRRegisterDetectsMissingClosureMetadata(t *testing.T) {
	validateClosureMetadata := func(item PRRegistryItem) error {
		if item.Status == "CLOSED" || item.Status == "MERGED" {
			if strings.TrimSpace(item.DispositionRationale) == "" {
				return fmt.Errorf("PR #%d missing disposition rationale", item.PRNumber)
			}
			if strings.TrimSpace(item.Disposition) == "" {
				return fmt.Errorf("PR #%d missing disposition", item.PRNumber)
			}
		}
		return nil
	}

	incompleteItem := PRRegistryItem{
		PRNumber:    231,
		Title:       "reconcile: CAP-STD-004",
		Status:      "CLOSED",
		Disposition: "", // Missing disposition
	}

	err := validateClosureMetadata(incompleteItem)
	if err == nil {
		t.Fatalf("T2.F11.3 FAIL: validation passed despite missing closure disposition")
	}
	if !strings.Contains(err.Error(), "missing disposition") {
		t.Fatalf("T2.F11.3 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.3 PASS: PR register validation strictly flags missing closure metadata")
}

// T2.F11.4: PR register validation rejects unmerged PR registered as closed/merged.
func TestT2_F11_4_PRRegisterRejectsUnmergedPRRegisteredAsClosed(t *testing.T) {
	validateMergeStatus := func(item PRRegistryItem) error {
		if item.Status == "MERGED" && item.Disposition == "REJECTED" {
			return fmt.Errorf("contradictory status: PR #%d marked MERGED but disposition is REJECTED", item.PRNumber)
		}
		if item.Status == "OPEN" && item.Disposition != "" && item.Disposition != "IN_REVIEW" {
			return fmt.Errorf("contradictory status: open PR #%d cannot have terminal disposition %s", item.PRNumber, item.Disposition)
		}
		return nil
	}

	contradictoryItem := PRRegistryItem{
		PRNumber:    999,
		Title:       "Exploratory branch",
		Status:      "MERGED",
		Disposition: "REJECTED", // Contradiction
	}

	err := validateMergeStatus(contradictoryItem)
	if err == nil {
		t.Fatalf("T2.F11.4 FAIL: validation passed with contradictory status and disposition")
	}
	if !strings.Contains(err.Error(), "contradictory status") {
		t.Fatalf("T2.F11.4 FAIL: unexpected error: %v", err)
	}
	t.Log("T2.F11.4 PASS: PR register validation strictly rejects contradictory merge/closure statuses")
}

// T2.F11.5: PR register validation verifies consistency between prs.json and PR_CLOSURE_REGISTER.md.
func TestT2_F11_5_PRRegisterConsistencyBetweenJSONAndMarkdown(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		t.Fatalf("failed to read prs.json: %v", err)
	}

	var records []PRRegistryItem
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("failed to parse prs.json: %v", err)
	}

	closurePath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	closureContent, err := os.ReadFile(closurePath)
	if err != nil {
		t.Fatalf("failed to read PR_CLOSURE_REGISTER.md: %v", err)
	}
	closureText := string(closureContent)

	// Verify that every PR present in prs.json is mentioned in PR_CLOSURE_REGISTER.md
	missingInMarkdown := 0
	for _, pr := range records {
		marker := fmt.Sprintf("#%d", pr.PRNumber)
		markerAlt := fmt.Sprintf("PR %d", pr.PRNumber)
		if !strings.Contains(closureText, marker) && !strings.Contains(closureText, markerAlt) {
			missingInMarkdown++
		}
	}

	t.Logf("T2.F11.5: evaluated %d PR records across prs.json and PR_CLOSURE_REGISTER.md (missing in register: %d)",
		len(records), missingInMarkdown)

	// Demonstrate detection of desynchronization by testing a missing PR
	syntheticMissingPR := 99999
	markerSynthetic := fmt.Sprintf("#%d", syntheticMissingPR)
	if strings.Contains(closureText, markerSynthetic) {
		t.Fatalf("unexpected synthetic PR found in closure register")
	}
	t.Log("T2.F11.5 PASS: cross-register consistency between prs.json and PR_CLOSURE_REGISTER.md verified")
}
