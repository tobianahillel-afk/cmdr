package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type CIObservation struct {
	HeadSHA    string `json:"head_sha"`
	RunID      uint64 `json:"run_id"`
	Event      string `json:"event"`
	Conclusion string `json:"conclusion"`
	ObservedAt string `json:"observed_at"`
}

type RecoveryReconcileRequest struct {
	WorkUnit string
	AsOf     string
	CI       *CIObservation
}

type RecoveryReconciliation struct {
	WorkUnit       string         `json:"work_unit"`
	CurrentHead    string         `json:"current_head"`
	Checkpoint     string         `json:"checkpoint_status"`
	CheckpointHead string         `json:"checkpoint_head,omitempty"`
	GitRelation    string         `json:"git_relation"`
	LeaseStatus    string         `json:"lease_status"`
	CIStatus       string         `json:"ci_status"`
	Outcome        string         `json:"outcome"`
	NextAction     string         `json:"next_action"`
	Reasons        []string       `json:"reasons"`
	CI             *CIObservation `json:"ci,omitempty"`
}

var knownCIConclusions = map[string]bool{
	"success": true, "failure": true, "cancelled": true, "timed_out": true,
	"action_required": true, "skipped": true, "neutral": true, "stale": true,
}

func runRecoveryReconciliation(root string, request RecoveryReconcileRequest, state CurrentState, graph WorkGraph) (RecoveryReconciliation, error) {
	workUnit := request.WorkUnit
	if workUnit == "" {
		workUnit = state.Execution.ActiveWorkUnit
	}
	if _, ok := graphNode(graph, workUnit); !ok {
		return RecoveryReconciliation{}, fmt.Errorf("unknown reconciliation work unit %s", workUnit)
	}
	asOf, err := parseLeaseTimeOrNow(request.AsOf)
	if err != nil {
		return RecoveryReconciliation{}, err
	}
	if request.CI != nil {
		if err := validateCIObservation(*request.CI, asOf); err != nil {
			return RecoveryReconciliation{}, err
		}
	}
	currentHead, err := localGitHead(root)
	if err != nil {
		return RecoveryReconciliation{}, err
	}
	checkpoint, err := runResumeCheckpoint(root, workUnit, state, graph)
	if err != nil {
		return RecoveryReconciliation{}, err
	}
	_, claims, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return RecoveryReconciliation{}, err
	}
	result, err := reconcileRecovery(root, checkpoint, claims, currentHead, request.CI, asOf)
	if err != nil {
		return result, err
	}
	return result, nil
}

func runRecoveryReconciliationAudit(root string, state CurrentState, graph WorkGraph) (RecoveryReconciliation, error) {
	result, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{}, state, graph)
	if err != nil {
		return result, err
	}
	if result.Outcome == "conflict" {
		return result, fmt.Errorf("recovery reconciliation conflict: %s", strings.Join(result.Reasons, ", "))
	}
	return result, nil
}

func reconcileRecovery(root string, checkpoint ResumeCheckpoint, claims WorkLeaseRegistry, currentHead string, ci *CIObservation, asOf time.Time) (RecoveryReconciliation, error) {
	result := RecoveryReconciliation{
		WorkUnit: checkpoint.WorkUnit, CurrentHead: currentHead, Checkpoint: checkpoint.Status,
		GitRelation: "none", LeaseStatus: "none", CIStatus: "not-provided", CI: ci,
	}
	add := func(reason string) {
		if !containsString(result.Reasons, reason) {
			result.Reasons = append(result.Reasons, reason)
		}
	}

	if checkpoint.Status == "empty" {
		result.Outcome = "resume"
		result.NextAction = "claim-or-start"
		add("no-prior-journal-event")
		return result, nil
	}

	result.CheckpointHead = checkpoint.HeadSHA
	relation, err := gitHeadRelation(root, checkpoint.HeadSHA, currentHead)
	if err != nil {
		return result, err
	}
	result.GitRelation = relation
	result.LeaseStatus = checkpointLeaseStatus(checkpoint, claims, asOf)
	if ci != nil {
		result.CIStatus = ci.Conclusion
	}

	if relation == "diverged" {
		result.Outcome = "conflict"
		result.NextAction = "stop-and-review-divergence"
		add("checkpoint-head-is-not-ancestor-of-current-head")
		return finalizeReconciliation(result), nil
	}
	if result.LeaseStatus == "missing" {
		result.Outcome = "conflict"
		result.NextAction = "stop-and-repair-lease-binding"
		add("checkpoint-lease-missing")
		return finalizeReconciliation(result), nil
	}
	if result.LeaseStatus == "expired" || result.LeaseStatus == "released" {
		result.Outcome = "stale-claim"
		result.NextAction = "reacquire-before-mutation"
		add("checkpoint-lease-" + result.LeaseStatus)
		if relation == "ancestor" {
			add("repository-advanced-since-checkpoint")
		}
		return finalizeReconciliation(result), nil
	}
	if relation == "ancestor" {
		result.Outcome = "revalidate"
		result.NextAction = "reconcile-newer-commits-and-rerun-selected-checks"
		add("repository-advanced-since-checkpoint")
		return finalizeReconciliation(result), nil
	}

	switch checkpoint.LastEventKind {
	case "validation-passed":
		if ci == nil {
			result.Outcome = "revalidate"
			result.NextAction = "obtain-ci-evidence-or-rerun-validation"
			add("validation-pass-has-no-ci-observation")
		} else if ci.HeadSHA != currentHead {
			result.Outcome = "revalidate"
			result.NextAction = "rerun-validation-for-current-head"
			add("ci-head-does-not-match-current-head")
		} else if ci.Conclusion != "success" {
			result.Outcome = "revalidate"
			result.NextAction = "fix-or-rerun-failed-ci"
			add("ci-conclusion-is-not-success")
		} else {
			result.Outcome = "resume"
			result.NextAction = "continue-after-validated-head"
			add("matching-successful-ci-evidence")
		}
	case "validation-failed":
		result.Outcome = "resume"
		result.NextAction = "fix-last-validation-failure"
		add("last-validation-failed")
	case "validation-started":
		result.Outcome = "revalidate"
		result.NextAction = "rerun-interrupted-validation"
		add("validation-was-interrupted")
	case "handoff":
		result.Outcome = "resume"
		result.NextAction = "consume-handoff"
		add("explicit-handoff")
	case "lease-released":
		result.Outcome = "stale-claim"
		result.NextAction = "reacquire-before-mutation"
		add("journal-records-lease-release")
	default:
		result.Outcome = "resume"
		result.NextAction = "continue-from-checkpoint"
		add("checkpoint-head-matches-current-head")
	}
	return finalizeReconciliation(result), nil
}

func finalizeReconciliation(result RecoveryReconciliation) RecoveryReconciliation {
	sortStrings(result.Reasons)
	return result
}

func validateCIObservation(ci CIObservation, asOf time.Time) error {
	normalized, err := validateFullCommitID(ci.HeadSHA)
	if err != nil {
		return fmt.Errorf("ci observation head_sha: %w", err)
	}
	if normalized != ci.HeadSHA {
		return fmt.Errorf("ci observation head_sha must be canonical lowercase")
	}
	if ci.RunID == 0 {
		return fmt.Errorf("ci observation run_id must be positive")
	}
	if ci.Event != "push" && ci.Event != "pull_request" {
		return fmt.Errorf("ci observation event must be push or pull_request")
	}
	if !knownCIConclusions[ci.Conclusion] {
		return fmt.Errorf("ci observation has unsupported conclusion %q", ci.Conclusion)
	}
	observed, err := parseCanonicalLeaseTime(ci.ObservedAt)
	if err != nil {
		return fmt.Errorf("ci observation observed_at: %w", err)
	}
	if observed.After(asOf) {
		return fmt.Errorf("ci observation is in the future")
	}
	return nil
}

func localGitHead(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "HEAD") // #nosec G204,G702 -- executable and arguments are fixed; root is the resolved repository root.
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("resolve current git head: %w", err)
	}
	head, err := validateFullCommitID(strings.TrimSpace(string(out)))
	if err != nil {
		return "", fmt.Errorf("current git head: %w", err)
	}
	return head, nil
}

func gitHeadRelation(root, checkpointHead, currentHead string) (string, error) {
	checkpointHead, err := validateFullCommitID(checkpointHead)
	if err != nil {
		return "", fmt.Errorf("checkpoint head: %w", err)
	}
	currentHead, err = validateFullCommitID(currentHead)
	if err != nil {
		return "", fmt.Errorf("current head: %w", err)
	}
	if checkpointHead == currentHead {
		return "same", nil
	}
	if !gitCommitExists(root, checkpointHead) || !gitCommitExists(root, currentHead) {
		return "", fmt.Errorf("cannot reconcile unavailable git commit")
	}
	cmd := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", checkpointHead, currentHead) // #nosec G204,G702 -- executable/flags are fixed and both SHAs are full-hex validated; no shell is used.
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err == nil {
		return "ancestor", nil
	}
	var exitErr *exec.ExitError
	if errorAs(err, &exitErr) && exitErr.ExitCode() == 1 {
		return "diverged", nil
	}
	return "", fmt.Errorf("git ancestry check failed: %s", strings.TrimSpace(stderr.String()))
}

func checkpointLeaseStatus(checkpoint ResumeCheckpoint, claims WorkLeaseRegistry, asOf time.Time) string {
	for _, claim := range claims.Claims {
		if claim.ID != checkpoint.LeaseID {
			continue
		}
		acquired, errA := parseCanonicalLeaseTime(claim.AcquiredAt)
		expires, errE := parseCanonicalLeaseTime(claim.ExpiresAt)
		if errA != nil || errE != nil {
			return "invalid"
		}
		var released *time.Time
		if claim.ReleasedAt != "" {
			value, err := parseCanonicalLeaseTime(claim.ReleasedAt)
			if err != nil {
				return "invalid"
			}
			released = &value
		}
		return derivedLeaseStatus(acquired, expires, released, asOf)
	}
	return "missing"
}

func ciObservationFromFlags(head, runID, event, conclusion, observedAt string) (*CIObservation, error) {
	values := []string{head, runID, event, conclusion, observedAt}
	provided := 0
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			provided++
		}
	}
	if provided == 0 {
		return nil, nil
	}
	if provided != len(values) {
		return nil, fmt.Errorf("CI observation requires head, run-id, event, conclusion and observed-at together")
	}
	id, err := strconv.ParseUint(runID, 10, 64)
	if err != nil || id == 0 {
		return nil, fmt.Errorf("invalid CI run id")
	}
	return &CIObservation{
		HeadSHA: head, RunID: id, Event: event, Conclusion: conclusion, ObservedAt: observedAt,
	}, nil
}

func sortStrings(values []string) {
	if len(values) < 2 {
		return
	}
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
