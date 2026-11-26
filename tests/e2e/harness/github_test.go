package harness

import (
	"errors"
	"testing"
)

func TestVerifyCodexPRRecord_RejectsUnattributedPR(t *testing.T) {
	itemMissingFields := map[string]interface{}{
		"pr_number":             float64(225),
		"status":                "MERGED",
		"disposition_rationale": "Merged commit e4666f586854e02766568ffd2cfe29ef32696b49",
	}

	res, err := verifyCodexPRRecord(itemMissingFields, 225, "e4666f586854e02766568ffd2cfe29ef32696b49", func(string) bool { return true })
	if !errors.Is(err, ErrUnverified) || res.Verified || res.Resolved || res.Status != StatusUnverified {
		t.Fatalf("expected ErrUnverified for unattributed prs.json record, got res=%+v err=%v", res, err)
	}

	itemExplicitFalse := map[string]interface{}{
		"pr_number":               float64(225),
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
		"review_threads_resolved": true,
		"merge_commit":            wrongSHA,
	}
	resWrong, errWrong := verifyCodexPRRecord(itemWrongSHA, 225, expectedSHA, func(string) bool { return true })
	if !errors.Is(errWrong, ErrUnverified) || resWrong.Verified || resWrong.Resolved {
		t.Fatalf("expected ErrUnverified on wrong SHA attribution, got res=%+v err=%v", resWrong, errWrong)
	}

	itemMissingCommit := map[string]interface{}{
		"pr_number":               float64(225),
		"review_threads_resolved": true,
		"merge_commit":            expectedSHA,
	}
	resNoCommit, errNoCommit := verifyCodexPRRecord(itemMissingCommit, 225, expectedSHA, func(string) bool { return false })
	if !errors.Is(errNoCommit, ErrUnverified) || resNoCommit.Verified || resNoCommit.Resolved {
		t.Fatalf("expected ErrUnverified when attributed commit is missing from git, got res=%+v err=%v", resNoCommit, errNoCommit)
	}

	itemShortSHA := map[string]interface{}{
		"pr_number":               float64(225),
		"review_threads_resolved": true,
		"merge_commit":            "e4666f58",
	}
	resShort, errShort := verifyCodexPRRecord(itemShortSHA, 225, expectedSHA, func(string) bool { return true })
	if !errors.Is(errShort, ErrUnverified) || resShort.Verified || resShort.Resolved {
		t.Fatalf("expected ErrUnverified on non-40-hex SHA, got res=%+v err=%v", resShort, errShort)
	}
}

func TestVerifyCodexPRRecord_AcceptsExplicitSHAAttributedProof(t *testing.T) {
	const expectedSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"

	item := map[string]interface{}{
		"pr_number":               float64(225),
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

func TestParseOpenPRsJSON(t *testing.T) {
	for _, invalid := range []string{"", "   ", "{not-json", "null", `{"number":1}`} {
		res, err := parseOpenPRsJSON(invalid)
		if !errors.Is(err, ErrUnverified) || res.Verified || res.Status != StatusUnverified {
			t.Fatalf("expected ErrUnverified for %q, got res=%+v err=%v", invalid, res, err)
		}
	}

	resZero, errZero := parseOpenPRsJSON("[]")
	if errZero != nil || !resZero.Verified || resZero.Count != 0 || resZero.Status != StatusVerified {
		t.Fatalf("expected Verified Count=0 for [], got res=%+v err=%v", resZero, errZero)
	}

	resOne, errOne := parseOpenPRsJSON(`[{"number":234}]`)
	if errOne != nil || !resOne.Verified || resOne.Count != 1 || resOne.Status != StatusVerified {
		t.Fatalf("expected Verified Count=1, got res=%+v err=%v", resOne, errOne)
	}
}
