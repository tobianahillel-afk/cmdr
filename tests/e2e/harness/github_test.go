package harness

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestVerifyCodexPRRecord_RejectsUnattributedPR(t *testing.T) {
	itemMissingFields := map[string]interface{}{
		"pr_number":             float64(225),
		"status":                "MERGED",
		"disposition":           "INCORPORATED_IN_MAIN",
		"disposition_rationale": "Merged commit e4666f586854e02766568ffd2cfe29ef32696b49",
	}

	res, err := verifyCodexPRRecord(itemMissingFields, 225, "e4666f586854e02766568ffd2cfe29ef32696b49", func(string) bool { return true })
	if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved || res.Status != StatusUnverified {
		t.Fatalf("expected ErrUnverified for unattributed prs.json record, got res=%+v err=%v", res, err)
	}

	itemExplicitFalse := map[string]interface{}{
		"pr_number":               float64(225),
		"status":                  "MERGED",
		"disposition":             "INCORPORATED_IN_MAIN",
		"review_threads_resolved": false,
		"merge_commit":            "e4666f586854e02766568ffd2cfe29ef32696b49",
	}
	resFalse, errFalse := verifyCodexPRRecord(itemExplicitFalse, 225, "e4666f586854e02766568ffd2cfe29ef32696b49", func(string) bool { return true })
	if !errors.Is(errFalse, ErrUnverified) || resFalse.Verified || resFalse.Resolved {
		t.Fatalf("expected ErrUnverified when review_threads_resolved=false, got res=%+v err=%v", resFalse, errFalse)
	}
}

func TestVerifyCodexPRRecord_RejectsWrongSHAAndMissingGitCommit(t *testing.T) {
	const expectedSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"
	const wrongSHA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	itemWrongSHA := map[string]interface{}{
		"pr_number":               float64(225),
		"status":                  "MERGED",
		"disposition":             "INCORPORATED_IN_MAIN",
		"review_threads_resolved": true,
		"merge_commit":            wrongSHA,
	}
	resWrong, errWrong := verifyCodexPRRecord(itemWrongSHA, 225, expectedSHA, func(string) bool { return true })
	if !errors.Is(errWrong, ErrUnverified) || resWrong.Verified || resWrong.Resolved {
		t.Fatalf("expected ErrUnverified on wrong SHA attribution, got res=%+v err=%v", resWrong, errWrong)
	}

	itemMissingCommit := map[string]interface{}{
		"pr_number":               float64(225),
		"status":                  "MERGED",
		"disposition":             "INCORPORATED_IN_MAIN",
		"review_threads_resolved": true,
		"merge_commit":            expectedSHA,
	}
	resNoCommit, errNoCommit := verifyCodexPRRecord(itemMissingCommit, 225, expectedSHA, func(string) bool { return false })
	if !errors.Is(errNoCommit, ErrUnverified) || resNoCommit.Verified || resNoCommit.Resolved {
		t.Fatalf("expected ErrUnverified when attributed commit is missing from git, got res=%+v err=%v", resNoCommit, errNoCommit)
	}

	itemShortSHA := map[string]interface{}{
		"pr_number":               float64(225),
		"status":                  "MERGED",
		"disposition":             "INCORPORATED_IN_MAIN",
		"review_threads_resolved": true,
		"merge_commit":            "e4666f58",
	}
	resShort, errShort := verifyCodexPRRecord(itemShortSHA, 225, expectedSHA, func(string) bool { return true })
	if !errors.Is(errShort, ErrUnverified) || resShort.Verified || resShort.Resolved {
		t.Fatalf("expected ErrUnverified on non-40-hex SHA, got res=%+v err=%v", resShort, errShort)
	}
}

func TestVerifyCodexPRRecord_RejectsInvalidStatusDispositionAndEmptySHA(t *testing.T) {
	const expectedSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"

	cases := []struct {
		name        string
		expectedSHA string
		item        map[string]interface{}
	}{
		{
			name:        "ClosedStatus",
			expectedSHA: expectedSHA,
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"status":                  "CLOSED",
				"disposition":             "INCORPORATED_IN_MAIN",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
		{
			name:        "MissingStatus",
			expectedSHA: expectedSHA,
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"disposition":             "INCORPORATED_IN_MAIN",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
		{
			name:        "RejectedDisposition",
			expectedSHA: expectedSHA,
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"status":                  "MERGED",
				"disposition":             "REJECTED",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
		{
			name:        "SupersededDisposition",
			expectedSHA: expectedSHA,
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"status":                  "MERGED",
				"disposition":             "SUPERSEDED",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
		{
			name:        "ClosedDisposition",
			expectedSHA: expectedSHA,
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"status":                  "MERGED",
				"disposition":             "CLOSED",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
		{
			name:        "EmptyExpectedSHA",
			expectedSHA: "",
			item: map[string]interface{}{
				"pr_number":               float64(225),
				"status":                  "MERGED",
				"disposition":             "INCORPORATED_IN_MAIN",
				"review_threads_resolved": true,
				"merge_commit":            expectedSHA,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := verifyCodexPRRecord(tc.item, 225, tc.expectedSHA, func(string) bool { return true })
			if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved || res.Status != StatusUnverified {
				t.Fatalf("expected ErrUnverified for %s, got res=%+v err=%v", tc.name, res, err)
			}
		})
	}
}

func TestVerifyCodexPRRecord_AcceptsExplicitSHAAttributedProof(t *testing.T) {
	const expectedSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"

	item := map[string]interface{}{
		"pr_number":               float64(225),
		"status":                  "MERGED",
		"disposition":             "INCORPORATED_IN_MAIN",
		"review_threads_resolved": true,
		"merge_commit":            expectedSHA,
	}
	res, err := verifyCodexPRRecord(item, 225, expectedSHA, func(sha string) bool { return sha == expectedSHA })
	if err != nil || !res.Verified || !res.Resolved || res.Status != StatusVerified || res.MatchedSHA != expectedSHA {
		t.Fatalf("expected verified resolution for explicit SHA-attributed proof, got res=%+v err=%v", res, err)
	}
}

func TestVerifyGraphQLPRResolution(t *testing.T) {
	const expectedSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"

	t.Run("RejectsGraphQLErrors", func(t *testing.T) {
		payload := `{"errors":[{"message":"API rate limit exceeded"}]}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if !errors.Is(err, ErrUnverified) || res.Verified {
			t.Fatalf("expected ErrUnverified on GraphQL errors, got res=%+v err=%v", res, err)
		}
	})

	t.Run("RejectsEmptyNodes", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"` + expectedSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[]}}}}}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if !errors.Is(err, ErrUnverified) || res.Verified {
			t.Fatalf("expected ErrUnverified on empty reviewThreads nodes, got res=%+v err=%v", res, err)
		}
	})

	t.Run("RejectsHasNextPage", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"` + expectedSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":true},"nodes":[{"isResolved":true}]}}}}}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if !errors.Is(err, ErrUnverified) || res.Verified {
			t.Fatalf("expected ErrUnverified when hasNextPage=true, got res=%+v err=%v", res, err)
		}
	})

	t.Run("RejectsNonMergedState", func(t *testing.T) {
		for _, state := range []string{"CLOSED", "OPEN", "UNKNOWN"} {
			payload := `{"data":{"repository":{"pullRequest":{"state":"` + state + `","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"` + expectedSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true}]}}}}}`
			res, err := verifyGraphQLPRResolution(payload, expectedSHA)
			if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved || res.Status != StatusUnverified {
				t.Fatalf("expected ErrUnverified for state %s, got res=%+v err=%v", state, res, err)
			}
		}
	})

	t.Run("RejectsEmptyExpectedSHA", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"` + expectedSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true}]}}}}}`
		for _, emptySHA := range []string{"", "   "} {
			res, err := verifyGraphQLPRResolution(payload, emptySHA)
			if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved {
				t.Fatalf("expected ErrUnverified for empty expectedSHA %q, got res=%+v err=%v", emptySHA, res, err)
			}
		}
	})

	t.Run("RejectsSHAMismatch", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"headRefOid":"cccccccccccccccccccccccccccccccccccccccc","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true}]}}}}}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if err != nil || !res.Verified || res.Resolved || res.Status != "SHA_MISMATCH" {
			t.Fatalf("expected Verified=true Resolved=false Status=SHA_MISMATCH, got res=%+v err=%v", res, err)
		}
	})

	t.Run("RejectsUnresolvedThread", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"` + expectedSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true},{"isResolved":false}]}}}}}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if err != nil || !res.Verified || res.Resolved || res.Status != StatusVerified {
			t.Fatalf("expected Verified=true Resolved=false, got res=%+v err=%v", res, err)
		}
	})

	t.Run("AcceptsResolvedWithMatchingMergeCommit", func(t *testing.T) {
		payload := `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + expectedSHA + `"},"headRefOid":"cccccccccccccccccccccccccccccccccccccccc","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true}]}}}}}`
		res, err := verifyGraphQLPRResolution(payload, expectedSHA)
		if err != nil || !res.Verified || !res.Resolved || res.Status != StatusVerified || res.MatchedSHA != expectedSHA {
			t.Fatalf("expected verified resolved result, got res=%+v err=%v", res, err)
		}
	})
}

func TestVerifyCodexPRResolution_RejectsEmptySHAAndNoFallbackOnLiveGraphQLFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tmpRepo := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.name", "Test"},
		{"config", "user.email", "test@cmdr.local"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", tmpRepo}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v (%s)", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(tmpRepo, "README.md"), []byte("seed\n"), 0644); err != nil {
		t.Fatalf("write seed: %v", err)
	}
	if out, err := exec.Command("git", "-C", tmpRepo, "add", "README.md").CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v (%s)", err, out)
	}
	if out, err := exec.Command("git", "-C", tmpRepo, "commit", "-m", "seed").CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v (%s)", err, out)
	}
	shaOut, err := exec.Command("git", "-C", tmpRepo, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v", err)
	}
	headSHA := strings.TrimSpace(string(shaOut))

	// Populate a valid prs.json so any accidental fallback would succeed instead of failing closed.
	recDir := filepath.Join(tmpRepo, "work", "reconciliation")
	if err := os.MkdirAll(recDir, 0755); err != nil {
		t.Fatalf("mkdir reconciliation: %v", err)
	}
	prsJSON, err := json.Marshal([]map[string]interface{}{
		{
			"pr_number":               225,
			"title":                   "PR 225",
			"status":                  "MERGED",
			"disposition":             "INCORPORATED_IN_MAIN",
			"disposition_rationale":   "Merged",
			"review_threads_resolved": true,
			"merge_commit":            headSHA,
		},
	})
	if err != nil {
		t.Fatalf("marshal prs.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(recDir, "prs.json"), prsJSON, 0644); err != nil {
		t.Fatalf("write prs.json: %v", err)
	}

	client := &GitHubClient{repoRoot: tmpRepo}

	// 1. Empty expectedSHA must fail immediately.
	resEmpty, errEmpty := client.VerifyCodexPRResolution(ctx, 225, "")
	if !errors.Is(errEmpty, ErrUnverified) || resEmpty.Verified || resEmpty.Resolved {
		t.Fatalf("expected ErrUnverified for empty expectedSHA, got res=%+v err=%v", resEmpty, errEmpty)
	}

	// 2. When live `gh` exits 0 with invalid/unverified GraphQL payloads, must NOT fall back to prs.json.
	binDir := t.TempDir()
	payloadFile := filepath.Join(binDir, "gql.json")
	if runtime.GOOS == "windows" {
		script := "@echo off\r\ntype \"" + payloadFile + "\"\r\nexit /b 0\r\n"
		if err := os.WriteFile(filepath.Join(binDir, "gh.cmd"), []byte(script), 0755); err != nil {
			t.Fatalf("write gh.cmd: %v", err)
		}
	} else {
		script := "#!/bin/sh\ncat \"" + payloadFile + "\"\nexit 0\n"
		if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(script), 0755); err != nil {
			t.Fatalf("write gh: %v", err)
		}
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	badPayloads := map[string]string{
		"GraphQLErrors": `{"errors":[{"message":"boom"}]}`,
		"HasNextPage":   `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + headSHA + `"},"headRefOid":"` + headSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":true},"nodes":[{"isResolved":true}]}}}}}`,
		"EmptyNodes":    `{"data":{"repository":{"pullRequest":{"state":"MERGED","mergeCommit":{"oid":"` + headSHA + `"},"headRefOid":"` + headSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[]}}}}}`,
		"ClosedState":   `{"data":{"repository":{"pullRequest":{"state":"CLOSED","mergeCommit":{"oid":"` + headSHA + `"},"headRefOid":"` + headSHA + `","reviewThreads":{"pageInfo":{"hasNextPage":false},"nodes":[{"isResolved":true}]}}}}}`,
	}

	for name, payload := range badPayloads {
		t.Run(name, func(t *testing.T) {
			if err := os.WriteFile(payloadFile, []byte(payload), 0644); err != nil {
				t.Fatalf("write payload: %v", err)
			}
			res, err := client.VerifyCodexPRResolution(ctx, 225, headSHA)
			if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved || res.Source != "github_graphql" {
				t.Fatalf("expected ErrUnverified from github_graphql without prs.json fallback for %s, got res=%+v err=%v", name, res, err)
			}
		})
	}
}

func TestParseOpenPRsJSON(t *testing.T) {
	invalidCases := []string{
		"",
		"   ",
		"{not-json",
		"null",
		`{"number":1}`,
		`[{"wrong":"obj"}]`,
		`[{}]`,
		`[null]`,
		`[{"number":0}]`,
		`[{"number":-1}]`,
		`[{"number":"234"}]`,
		`[{"number":1.5}]`,
	}
	for _, invalid := range invalidCases {
		res, err := ParseOpenPRsJSON(invalid)
		if !errors.Is(err, ErrUnverified) || res.Verified || res.Status != StatusUnverified {
			t.Fatalf("expected ErrUnverified for %q, got res=%+v err=%v", invalid, res, err)
		}
	}

	resZero, errZero := ParseOpenPRsJSON("[]")
	if errZero != nil || !resZero.Verified || resZero.Count != 0 || resZero.Status != StatusVerified {
		t.Fatalf("expected Verified Count=0 for [], got res=%+v err=%v", resZero, errZero)
	}

	resOne, errOne := ParseOpenPRsJSON(`[{"number":234}]`)
	if errOne != nil || !resOne.Verified || resOne.Count != 1 || resOne.Status != StatusVerified {
		t.Fatalf("expected Verified Count=1, got res=%+v err=%v", resOne, errOne)
	}
}

func TestValidatePRRegisterConsistency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tmpRepo := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.name", "Test"},
		{"config", "user.email", "test@cmdr.local"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", tmpRepo}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v (%s)", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(tmpRepo, "README.md"), []byte("seed\n"), 0644); err != nil {
		t.Fatalf("write seed: %v", err)
	}
	if out, err := exec.Command("git", "-C", tmpRepo, "add", "README.md").CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v (%s)", err, out)
	}
	if out, err := exec.Command("git", "-C", tmpRepo, "commit", "-m", "seed").CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v (%s)", err, out)
	}
	shaOut, err := exec.Command("git", "-C", tmpRepo, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v", err)
	}
	validSHA := strings.TrimSpace(string(shaOut))

	recDir := filepath.Join(tmpRepo, "work", "reconciliation")
	if err := os.MkdirAll(recDir, 0755); err != nil {
		t.Fatalf("mkdir reconciliation: %v", err)
	}

	writeRegister := func(records []PRMetadata, md string) {
		b, mErr := json.MarshalIndent(records, "", "  ")
		if mErr != nil {
			t.Fatalf("marshal prs.json: %v", mErr)
		}
		if err := os.WriteFile(filepath.Join(recDir, "prs.json"), b, 0644); err != nil {
			t.Fatalf("write prs.json: %v", err)
		}
		if err := os.WriteFile(filepath.Join(recDir, "PR_CLOSURE_REGISTER.md"), []byte(md), 0644); err != nil {
			t.Fatalf("write PR_CLOSURE_REGISTER.md: %v", err)
		}
	}

	validRecord := PRMetadata{
		PRNumber:             230,
		Title:                "feat: health gate",
		Status:               "MERGED",
		Disposition:          "INCORPORATED_IN_MAIN",
		DispositionRationale: "Merged into main",
		Commit:               validSHA,
	}

	writeRegister([]PRMetadata{validRecord}, "# PR Closure Register\n- #230 MERGED\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err != nil {
		t.Fatalf("expected valid register to pass, got %v", err)
	}

	badZero := validRecord
	badZero.Commit = "0000000000000000000000000000000000000000"
	writeRegister([]PRMetadata{badZero}, "- #230\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err == nil {
		t.Fatalf("expected non-existent commit SHA to fail")
	}

	badFormat := validRecord
	badFormat.Commit = "not-a-sha"
	writeRegister([]PRMetadata{badFormat}, "- #230\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err == nil {
		t.Fatalf("expected malformed SHA to fail")
	}

	badDisp := validRecord
	badDisp.Disposition = ""
	writeRegister([]PRMetadata{badDisp}, "- #230\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err == nil {
		t.Fatalf("expected missing disposition to fail")
	}

	badContradiction := validRecord
	badContradiction.Disposition = "REJECTED"
	writeRegister([]PRMetadata{badContradiction}, "- #230\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err == nil {
		t.Fatalf("expected MERGED+REJECTED contradiction to fail")
	}

	writeRegister([]PRMetadata{validRecord}, "# Empty register\n")
	if err := ValidatePRRegisterConsistency(ctx, tmpRepo); err == nil {
		t.Fatalf("expected missing PR #230 in markdown register to fail")
	}
}
