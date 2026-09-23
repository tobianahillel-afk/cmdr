package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type CoordinationConflict struct {
	LeftLease    string `json:"left_lease"`
	LeftWorkUnit string `json:"left_work_unit"`
	LeftAgent    string `json:"left_agent"`
	LeftPath     string `json:"left_path"`
	RightLease   string `json:"right_lease"`
	RightWorkUnit string `json:"right_work_unit"`
	RightAgent   string `json:"right_agent"`
	RightPath    string `json:"right_path"`
}

type CoordinationAuditSummary struct {
	AsOf          string                 `json:"as_of"`
	ActiveClaims  int                    `json:"active_claims"`
	Mutating      int                    `json:"mutating_claims"`
	ReadOnly      int                    `json:"read_only_claims"`
	Conflicts     int                    `json:"conflicts"`
	Status        string                 `json:"status"`
	ConflictItems []CoordinationConflict `json:"conflict_items,omitempty"`
}

type CoordinationHandoff struct {
	SchemaVersion          int      `json:"schema_version"`
	WorkUnit               string   `json:"work_unit"`
	LeaseID                string   `json:"lease_id"`
	LeaseMode              string   `json:"lease_mode"`
	FromAgent              string   `json:"from_agent"`
	HeadSHA                string   `json:"head_sha"`
	GeneratedAt            string   `json:"generated_at"`
	CheckpointStatus       string   `json:"checkpoint_status"`
	CheckpointDigest       string   `json:"checkpoint_digest,omitempty"`
	RevalidationRequired   bool     `json:"revalidation_required"`
	JournalEventKind       string   `json:"journal_event_kind"`
	ReleaseRequired        bool     `json:"release_required"`
	ReceiverMustReacquire  bool     `json:"receiver_must_reacquire"`
	MutationAfterRelease   bool     `json:"mutation_after_release_allowed"`
	Protocol               []string `json:"protocol"`
	Digest                  string   `json:"digest_sha256"`
}

type coordinationHandoffBody struct {
	SchemaVersion          int      `json:"schema_version"`
	WorkUnit               string   `json:"work_unit"`
	LeaseID                string   `json:"lease_id"`
	LeaseMode              string   `json:"lease_mode"`
	FromAgent              string   `json:"from_agent"`
	HeadSHA                string   `json:"head_sha"`
	GeneratedAt            string   `json:"generated_at"`
	CheckpointStatus       string   `json:"checkpoint_status"`
	CheckpointDigest       string   `json:"checkpoint_digest,omitempty"`
	RevalidationRequired   bool     `json:"revalidation_required"`
	JournalEventKind       string   `json:"journal_event_kind"`
	ReleaseRequired        bool     `json:"release_required"`
	ReceiverMustReacquire  bool     `json:"receiver_must_reacquire"`
	MutationAfterRelease   bool     `json:"mutation_after_release_allowed"`
	Protocol               []string `json:"protocol"`
}

type activeCoordinationClaim struct {
	Lease    WorkLease
	Manifest WorkManifestV2
}

func runCoordinationAudit(root, asOfValue string, graph WorkGraph) (CoordinationAuditSummary, error) {
	policy, registry, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return CoordinationAuditSummary{}, err
	}
	asOf, err := parseLeaseTimeOrNow(asOfValue)
	if err != nil {
		return CoordinationAuditSummary{}, err
	}
	if _, err := auditWorkLeases(policy, registry, graph, asOf); err != nil {
		return CoordinationAuditSummary{}, err
	}
	active, err := loadActiveCoordinationClaims(root, policy, registry, graph, asOf)
	if err != nil {
		return CoordinationAuditSummary{}, err
	}
	summary := evaluateCoordinationConflicts(active, asOf)
	if summary.Conflicts > 0 {
		return summary, fmt.Errorf("coordination conflict: %d overlapping mutable scope pair(s)", summary.Conflicts)
	}
	return summary, nil
}

func loadActiveCoordinationClaims(root string, policy WorkLeasePolicy, registry WorkLeaseRegistry, graph WorkGraph, asOf time.Time) ([]activeCoordinationClaim, error) {
	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	var active []activeCoordinationClaim
	for _, claim := range registry.Claims {
		acquired, expires, released, err := validateWorkLeaseRecord(policy, claim, nodes, asOf)
		if err != nil {
			return nil, err
		}
		if derivedLeaseStatus(acquired, expires, released, asOf) != "active" {
			continue
		}
		manifestPath := filepath.Join(root, "work", "lots", claim.WorkUnit, "manifest.json")
		manifest, err := decodeWorkManifestV2(root, manifestPath)
		if err != nil {
			return nil, fmt.Errorf("coordination manifest %s: %w", claim.WorkUnit, err)
		}
		if manifest.ID != claim.WorkUnit {
			return nil, fmt.Errorf("coordination manifest mismatch for %s", claim.WorkUnit)
		}
		active = append(active, activeCoordinationClaim{Lease: claim, Manifest: manifest})
	}
	sort.Slice(active, func(i, j int) bool { return active[i].Lease.ID < active[j].Lease.ID })
	return active, nil
}

func evaluateCoordinationConflicts(active []activeCoordinationClaim, asOf time.Time) CoordinationAuditSummary {
	summary := CoordinationAuditSummary{
		AsOf:         asOf.UTC().Format(time.RFC3339),
		ActiveClaims: len(active),
		Status:       "clear",
	}
	for _, item := range active {
		switch normalizedLeaseMode(item.Lease.Mode) {
		case "read-only":
			summary.ReadOnly++
		default:
			summary.Mutating++
		}
	}
	for i := 0; i < len(active); i++ {
		for j := i + 1; j < len(active); j++ {
			left, right := active[i], active[j]
			if normalizedLeaseMode(left.Lease.Mode) == "read-only" || normalizedLeaseMode(right.Lease.Mode) == "read-only" {
				continue
			}
			for _, lp := range left.Manifest.AllowedPaths {
				for _, rp := range right.Manifest.AllowedPaths {
					if !pathPatternsOverlap(lp, rp) {
						continue
					}
					summary.ConflictItems = append(summary.ConflictItems, CoordinationConflict{
						LeftLease: left.Lease.ID, LeftWorkUnit: left.Lease.WorkUnit, LeftAgent: left.Lease.AgentID, LeftPath: lp,
						RightLease: right.Lease.ID, RightWorkUnit: right.Lease.WorkUnit, RightAgent: right.Lease.AgentID, RightPath: rp,
					})
				}
			}
		}
	}
	sort.Slice(summary.ConflictItems, func(i, j int) bool {
		a, b := summary.ConflictItems[i], summary.ConflictItems[j]
		if a.LeftLease != b.LeftLease { return a.LeftLease < b.LeftLease }
		if a.RightLease != b.RightLease { return a.RightLease < b.RightLease }
		if a.LeftPath != b.LeftPath { return a.LeftPath < b.LeftPath }
		return a.RightPath < b.RightPath
	})
	summary.Conflicts = len(summary.ConflictItems)
	if summary.Conflicts > 0 {
		summary.Status = "conflict"
	}
	return summary
}

func runCoordinationHandoff(root, workUnit, leaseID, agentID, asOfValue string, state CurrentState, graph WorkGraph) (CoordinationHandoff, error) {
	if workUnit == "" {
		workUnit = state.Execution.ActiveWorkUnit
	}
	if !agentIDPattern.MatchString(agentID) || !workLeaseIDPattern.MatchString(leaseID) {
		return CoordinationHandoff{}, fmt.Errorf("handoff requires valid agent and lease identities")
	}
	asOf, err := parseLeaseTimeOrNow(asOfValue)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	if _, err := runCoordinationAudit(root, asOf.Format(time.RFC3339), graph); err != nil {
		return CoordinationHandoff{}, err
	}
	policy, registry, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes { nodes[node.ID] = node }
	var claim *WorkLease
	for i := range registry.Claims {
		if registry.Claims[i].ID == leaseID {
			value := registry.Claims[i]
			claim = &value
			break
		}
	}
	if claim == nil || claim.AgentID != agentID || claim.WorkUnit != workUnit {
		return CoordinationHandoff{}, fmt.Errorf("handoff lease binding does not match work unit and agent")
	}
	acquired, expires, released, err := validateWorkLeaseRecord(policy, *claim, nodes, asOf)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	if derivedLeaseStatus(acquired, expires, released, asOf) != "active" {
		return CoordinationHandoff{}, fmt.Errorf("handoff requires an active lease; reacquire before mutation or handoff")
	}
	head, err := localGitHead(root)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	relation, err := gitHeadRelation(root, claim.BaseHeadSHA, head)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	if relation == "diverged" {
		return CoordinationHandoff{}, fmt.Errorf("handoff lease base diverges from current Git HEAD")
	}
	checkpoint, err := runResumeCheckpoint(root, workUnit, state, graph)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	reconciliation, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{WorkUnit: workUnit, AsOf: asOf.Format(time.RFC3339)}, state, graph)
	if err != nil {
		return CoordinationHandoff{}, err
	}
	if reconciliation.Outcome == "conflict" || reconciliation.Outcome == "stale-claim" {
		return CoordinationHandoff{}, fmt.Errorf("handoff blocked by recovery outcome %s", reconciliation.Outcome)
	}
	return buildCoordinationHandoff(*claim, head, asOf, checkpoint, reconciliation), nil
}

func buildCoordinationHandoff(claim WorkLease, head string, asOf time.Time, checkpoint ResumeCheckpoint, reconciliation RecoveryReconciliation) CoordinationHandoff {
	body := coordinationHandoffBody{
		SchemaVersion: 1,
		WorkUnit: claim.WorkUnit, LeaseID: claim.ID, LeaseMode: normalizedLeaseMode(claim.Mode),
		FromAgent: claim.AgentID, HeadSHA: head, GeneratedAt: asOf.UTC().Format(time.RFC3339),
		CheckpointStatus: checkpoint.Status, CheckpointDigest: checkpoint.LastEventDigest,
		RevalidationRequired: reconciliation.Outcome == "revalidate",
		JournalEventKind: "handoff", ReleaseRequired: true, ReceiverMustReacquire: true,
		MutationAfterRelease: false,
		Protocol: []string{
			"append a handoff journal event bound to this exact work unit, lease and Git head",
			"persist release of the current lease through an ordinary compare-and-swap Git commit",
			"receiver acquires a new lease before any mutation",
			"receiver runs recovery reconciliation against current Git and CI reality before continuing",
		},
	}
	return CoordinationHandoff{
		SchemaVersion: body.SchemaVersion, WorkUnit: body.WorkUnit, LeaseID: body.LeaseID, LeaseMode: body.LeaseMode,
		FromAgent: body.FromAgent, HeadSHA: body.HeadSHA, GeneratedAt: body.GeneratedAt,
		CheckpointStatus: body.CheckpointStatus, CheckpointDigest: body.CheckpointDigest,
		RevalidationRequired: body.RevalidationRequired, JournalEventKind: body.JournalEventKind,
		ReleaseRequired: body.ReleaseRequired, ReceiverMustReacquire: body.ReceiverMustReacquire,
		MutationAfterRelease: body.MutationAfterRelease, Protocol: append([]string(nil), body.Protocol...),
		Digest: digestCanonical(body),
	}
}

func normalizedLeaseMode(mode string) string {
	switch strings.TrimSpace(mode) {
	case "", "mutating":
		return "mutating"
	case "read-only":
		return "read-only"
	default:
		return "invalid"
	}
}
