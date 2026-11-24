package tier2_boundaries

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// evaluateFinalCertificationGate evaluates the composite requirements for declaring MAIN_HEALTHY.
func evaluateFinalCertificationGate(
	remoteBranches []string,
	openPRs int,
	unpushedCommits int,
	healthGateStatus string,
	headCommit string,
) (bool, string) {
	for _, b := range remoteBranches {
		trimmed := strings.TrimSpace(b)
		if trimmed != "" && !strings.HasSuffix(trimmed, "/main") && !strings.HasSuffix(trimmed, "/HEAD") {
			return false, fmt.Sprintf("orphaned remote feature branch detected: %s", trimmed)
		}
	}

	if openPRs > 0 {
		return false, fmt.Sprintf("%d open PRs exist on repository", openPRs)
	}

	if unpushedCommits > 0 {
		return false, fmt.Sprintf("%d unpushed commits exist ahead of origin/main", unpushedCommits)
	}

	if healthGateStatus != "MAIN_HEALTHY" {
		return false, fmt.Sprintf("health gate status is %s (must be MAIN_HEALTHY)", healthGateStatus)
	}

	shaPattern := regexp.MustCompile(`^[0-9a-f]{40}$`)
	if !shaPattern.MatchString(headCommit) {
		return false, fmt.Sprintf("invalid HEAD commit SHA format: %s", headCommit)
	}

	return true, ""
}

// T2.F14.1: Final certification fails closed if any remote feature branch exists.
func TestT2_F14_1_CertificationFailsClosedOnRemoteFeatureBranch(t *testing.T) {
	branchesWithFeature := []string{
		"origin/main",
		"origin/feat/unmerged-studio-drafts",
	}

	passed, reason := evaluateFinalCertificationGate(branchesWithFeature, 0, 0, "MAIN_HEALTHY", "c1e75db81662b33a0f1208afdbf8b688fe7977e3")
	if passed {
		t.Fatalf("T2.F14.1 FAIL: certification passed despite presence of remote feature branch")
	}
	if !strings.Contains(reason, "orphaned remote feature branch detected") {
		t.Fatalf("T2.F14.1 FAIL: unexpected failure reason: %s", reason)
	}
	t.Log("T2.F14.1 PASS: certification strictly fails closed when remote feature branches exist")
}

// T2.F14.2: Final certification fails closed if any open PR exists.
func TestT2_F14_2_CertificationFailsClosedOnOpenPR(t *testing.T) {
	branchesClean := []string{"origin/main"}

	passed, reason := evaluateFinalCertificationGate(branchesClean, 1, 0, "MAIN_HEALTHY", "c1e75db81662b33a0f1208afdbf8b688fe7977e3")
	if passed {
		t.Fatalf("T2.F14.2 FAIL: certification passed despite 1 open PR")
	}
	if !strings.Contains(reason, "open PRs exist") {
		t.Fatalf("T2.F14.2 FAIL: unexpected failure reason: %s", reason)
	}
	t.Log("T2.F14.2 PASS: certification strictly fails closed when open PRs exist")
}

// T2.F14.3: Final certification fails closed if unpushed commits exist.
func TestT2_F14_3_CertificationFailsClosedOnUnpushedCommits(t *testing.T) {
	branchesClean := []string{"origin/main"}

	passed, reason := evaluateFinalCertificationGate(branchesClean, 0, 2, "MAIN_HEALTHY", "c1e75db81662b33a0f1208afdbf8b688fe7977e3")
	if passed {
		t.Fatalf("T2.F14.3 FAIL: certification passed despite unpushed commits")
	}
	if !strings.Contains(reason, "unpushed commits exist") {
		t.Fatalf("T2.F14.3 FAIL: unexpected failure reason: %s", reason)
	}
	t.Log("T2.F14.3 PASS: certification strictly fails closed when unpushed commits exist ahead of main")
}

// T2.F14.4: Final certification rejects MAIN_HEALTHY if any health dimension is UNVERIFIED.
func TestT2_F14_4_CertificationRejectsUnverifiedDimension(t *testing.T) {
	branchesClean := []string{"origin/main"}

	passed, reason := evaluateFinalCertificationGate(branchesClean, 0, 0, "MAIN_UNVERIFIED", "c1e75db81662b33a0f1208afdbf8b688fe7977e3")
	if passed {
		t.Fatalf("T2.F14.4 FAIL: certification declared healthy when health gate status is MAIN_UNVERIFIED")
	}
	if !strings.Contains(reason, "MAIN_UNVERIFIED") {
		t.Fatalf("T2.F14.4 FAIL: unexpected failure reason: %s", reason)
	}
	t.Log("T2.F14.4 PASS: certification strictly rejects MAIN_HEALTHY if any dimension is UNVERIFIED")
}

// T2.F14.5: Final certification asserts HEAD commit is verified.
func TestT2_F14_5_CertificationAssertsHEADCommitVerified(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	headSHA, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to get git HEAD commit: %v", err)
	}

	if !harness.GitCommitExists(ctx, headSHA) {
		t.Fatalf("T2.F14.5 FAIL: HEAD commit %s not found in git object store", headSHA)
	}

	// Verify malformed HEAD is rejected
	passed, _ := evaluateFinalCertificationGate([]string{"origin/main"}, 0, 0, "MAIN_HEALTHY", "invalid-sha-xyz")
	if passed {
		t.Fatalf("T2.F14.5 FAIL: certification accepted invalid HEAD commit SHA format")
	}

	t.Logf("T2.F14.5 PASS: HEAD commit %s verified and validated in git object storage", headSHA)
}
