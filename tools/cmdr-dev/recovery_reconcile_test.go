package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReconcileEmptyCheckpointStartsCleanly(t *testing.T) {
	root := initReconcileGit(t)
	head := testGitOutput(t, root, "rev-parse", "HEAD")
	result, err := reconcileRecovery(root, ResumeCheckpoint{Status: "empty", WorkUnit: "E7-REC-001C"}, WorkLeaseRegistry{}, head, nil, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "resume" || result.NextAction != "claim-or-start" {
		t.Fatalf("unexpected fresh reconciliation: %#v", result)
	}
}

func TestValidationPassedWithoutCIRequiresRevalidation(t *testing.T) {
	root := initReconcileGit(t)
	head := testGitOutput(t, root, "rev-parse", "HEAD")
	cp, claims := reconcileCheckpointAndClaim(head, "validation-passed")
	result, err := reconcileRecovery(root, cp, claims, head, nil, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "revalidate" || !containsString(result.Reasons, "validation-pass-has-no-ci-observation") {
		t.Fatalf("missing CI was treated as success: %#v", result)
	}
}

func TestMatchingSuccessfulCIAllowsResume(t *testing.T) {
	root := initReconcileGit(t)
	head := testGitOutput(t, root, "rev-parse", "HEAD")
	cp, claims := reconcileCheckpointAndClaim(head, "validation-passed")
	ci := &CIObservation{HeadSHA: head, RunID: 42, Event: "push", Conclusion: "success", ObservedAt: "2026-09-23T11:30:00Z"}
	result, err := reconcileRecovery(root, cp, claims, head, ci, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "resume" || result.NextAction != "continue-after-validated-head" {
		t.Fatalf("valid CI evidence did not permit resume: %#v", result)
	}
}

func TestRepositoryAdvanceForcesRevalidation(t *testing.T) {
	root := initReconcileGit(t)
	oldHead := testGitOutput(t, root, "rev-parse", "HEAD")
	writeTestFile(t, root, "new.txt", "new\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "advance")
	current := testGitOutput(t, root, "rev-parse", "HEAD")
	cp, claims := reconcileCheckpointAndClaim(oldHead, "commit-observed")
	result, err := reconcileRecovery(root, cp, claims, current, nil, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.GitRelation != "ancestor" || result.Outcome != "revalidate" {
		t.Fatalf("repository advance was not revalidated: %#v", result)
	}
}

func TestDivergedHeadFailsClosed(t *testing.T) {
	root := initReconcileGit(t)
	base := testGitOutput(t, root, "rev-parse", "HEAD")
	runTestGit(t, root, "checkout", "-b", "other")
	writeTestFile(t, root, "other.txt", "other\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "other")
	other := testGitOutput(t, root, "rev-parse", "HEAD")
	runTestGit(t, root, "checkout", "-b", "current", base)
	writeTestFile(t, root, "current.txt", "current\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "current")
	current := testGitOutput(t, root, "rev-parse", "HEAD")
	cp, claims := reconcileCheckpointAndClaim(other, "commit-observed")
	result, err := reconcileRecovery(root, cp, claims, current, nil, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.GitRelation != "diverged" || result.Outcome != "conflict" {
		t.Fatalf("divergence was not stopped: %#v", result)
	}
}

func TestExpiredLeaseRequiresReacquisition(t *testing.T) {
	root := initReconcileGit(t)
	head := testGitOutput(t, root, "rev-parse", "HEAD")
	cp, claims := reconcileCheckpointAndClaim(head, "implementation-started")
	claims.Claims[0].ExpiresAt = "2026-09-23T11:00:00Z"
	result, err := reconcileRecovery(root, cp, claims, head, nil, mustLeaseTime(t, "2026-09-23T12:00:00Z"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "stale-claim" || result.NextAction != "reacquire-before-mutation" {
		t.Fatalf("expired lease was not treated as stale: %#v", result)
	}
}

func TestCIObservationMustBeCompleteAndCanonical(t *testing.T) {
	if _, err := ciObservationFromFlags(strings.Repeat("a", 40), "", "", "", ""); err == nil {
		t.Fatal("expected partial CI observation rejection")
	}
	ci, err := ciObservationFromFlags(strings.Repeat("a", 40), "42", "push", "success", "2026-09-23T11:30:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCIObservation(*ci, mustLeaseTime(t, "2026-09-23T12:00:00Z")); err != nil {
		t.Fatal(err)
	}
}

func initReconcileGit(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runTestGit(t, root, "init")
	runTestGit(t, root, "config", "user.name", "CMDR Test")
	runTestGit(t, root, "config", "user.email", "cmdr-test@example.invalid")
	if err := os.WriteFile(filepath.Join(root, "base.txt"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "base")
	return root
}

func reconcileCheckpointAndClaim(head, kind string) (ResumeCheckpoint, WorkLeaseRegistry) {
	return ResumeCheckpoint{
			Status: resumeStatusForEvent(kind), WorkUnit: "E7-REC-001C", Events: 1,
			LastSequence: 1, LastEventID: "EVT-TEST0001", LastEventKind: kind,
			AgentID: "AGENT-test", LeaseID: "LEASE-TEST0001", HeadSHA: head,
		}, WorkLeaseRegistry{Claims: []WorkLease{{
			ID: "LEASE-TEST0001", WorkUnit: "E7-REC-001C", AgentID: "AGENT-test", BaseHeadSHA: head,
			AcquiredAt: "2026-09-23T10:00:00Z", ExpiresAt: "2026-09-23T14:00:00Z",
		}}}
}
