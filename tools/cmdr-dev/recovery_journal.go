package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const recoveryJournalPath = "engineering/recovery/journal.json"

var (
	recoveryEventIDPattern    = regexp.MustCompile(`^EVT-[A-Z0-9][A-Z0-9._-]{7,63}$`)
	recoveryJournalRootDigest = digestCanonical("cmdr-execution-journal-v1")
	knownRecoveryEventKinds   = map[string]bool{
		"lease-acquired":         true,
		"implementation-started": true,
		"commit-observed":        true,
		"validation-started":     true,
		"validation-passed":      true,
		"validation-failed":      true,
		"handoff":                true,
		"lease-released":         true,
	}
)

type RecoveryJournalEvent struct {
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

type RecoveryJournal struct {
	SchemaVersion int                    `json:"schema_version"`
	RegistryKind  string                 `json:"registry_kind"`
	Events        []RecoveryJournalEvent `json:"events"`
}

type recoveryJournalEventBody struct {
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
}

type RecoveryJournalAuditSummary struct {
	Events     int    `json:"events"`
	WorkUnits  int    `json:"work_units"`
	Agents     int    `json:"agents"`
	LastSeq    int    `json:"last_sequence"`
	LastDigest string `json:"last_digest"`
	Status     string `json:"status"`
}

type ResumeCheckpoint struct {
	Status         string   `json:"status"`
	WorkUnit       string   `json:"work_unit"`
	Events         int      `json:"events"`
	LastSequence   int      `json:"last_sequence,omitempty"`
	LastEventID    string   `json:"last_event_id,omitempty"`
	LastEventKind  string   `json:"last_event_kind,omitempty"`
	LastEventDigest string  `json:"last_event_digest,omitempty"`
	AgentID        string   `json:"agent_id,omitempty"`
	LeaseID        string   `json:"lease_id,omitempty"`
	HeadSHA        string   `json:"head_sha,omitempty"`
	EvidenceRefs   []string `json:"evidence_refs,omitempty"`
	RecentEventIDs []string `json:"recent_event_ids,omitempty"`
}

func runRecoveryJournalAudit(root string, graph WorkGraph) (RecoveryJournalAuditSummary, error) {
	_, claims, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return RecoveryJournalAuditSummary{}, err
	}
	var journal RecoveryJournal
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(recoveryJournalPath)), &journal); err != nil {
		return RecoveryJournalAuditSummary{}, err
	}
	return auditRecoveryJournal(journal, claims, graph)
}

func auditRecoveryJournal(journal RecoveryJournal, claims WorkLeaseRegistry, graph WorkGraph) (RecoveryJournalAuditSummary, error) {
	summary := RecoveryJournalAuditSummary{Events: len(journal.Events), Status: "valid"}
	if journal.SchemaVersion != 1 || journal.RegistryKind != "engineering-recovery-journal" {
		return summary, fmt.Errorf("invalid recovery journal header")
	}
	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	claimByID := map[string]WorkLease{}
	for _, claim := range claims.Claims {
		claimByID[claim.ID] = claim
	}
	seenIDs := map[string]bool{}
	workUnits := map[string]bool{}
	agents := map[string]bool{}
	previousDigest := recoveryJournalRootDigest
	var previousTime time.Time

	for i, event := range journal.Events {
		if event.Sequence != i+1 {
			return summary, fmt.Errorf("journal sequence must be contiguous: expected %d got %d", i+1, event.Sequence)
		}
		if !recoveryEventIDPattern.MatchString(event.ID) || seenIDs[event.ID] {
			return summary, fmt.Errorf("invalid or duplicate recovery event id %q", event.ID)
		}
		seenIDs[event.ID] = true
		if _, ok := nodes[event.WorkUnit]; !ok {
			return summary, fmt.Errorf("event %s references unknown work unit %s", event.ID, event.WorkUnit)
		}
		if !agentIDPattern.MatchString(event.AgentID) || !workLeaseIDPattern.MatchString(event.LeaseID) {
			return summary, fmt.Errorf("event %s has invalid agent or lease identity", event.ID)
		}
		claim, ok := claimByID[event.LeaseID]
		if !ok || claim.AgentID != event.AgentID || claim.WorkUnit != event.WorkUnit {
			return summary, fmt.Errorf("event %s does not match its lease binding", event.ID)
		}
		if !knownRecoveryEventKinds[event.Kind] {
			return summary, fmt.Errorf("event %s has unknown kind %q", event.ID, event.Kind)
		}
		observed, err := parseCanonicalLeaseTime(event.ObservedAt)
		if err != nil {
			return summary, fmt.Errorf("event %s observed_at: %w", event.ID, err)
		}
		if !previousTime.IsZero() && observed.Before(previousTime) {
			return summary, fmt.Errorf("event %s moves journal time backwards", event.ID)
		}
		previousTime = observed
		if _, err := validateFullCommitID(event.HeadSHA); err != nil {
			return summary, fmt.Errorf("event %s head_sha: %w", event.ID, err)
		}
		if event.PreviousDigest != previousDigest {
			return summary, fmt.Errorf("event %s previous_digest breaks journal chain", event.ID)
		}
		if err := validateRecoveryEvidenceRefs(event.ID, event.EvidenceRefs); err != nil {
			return summary, err
		}
		if !sha256Pattern.MatchString(event.Digest) || event.Digest != recoveryEventDigest(event) {
			return summary, fmt.Errorf("event %s digest mismatch", event.ID)
		}
		acquired, _ := parseCanonicalLeaseTime(claim.AcquiredAt)
		expires, _ := parseCanonicalLeaseTime(claim.ExpiresAt)
		if observed.Before(acquired) || observed.After(expires) {
			return summary, fmt.Errorf("event %s is outside its lease interval", event.ID)
		}
		if claim.ReleasedAt != "" {
			released, _ := parseCanonicalLeaseTime(claim.ReleasedAt)
			if observed.After(released) {
				return summary, fmt.Errorf("event %s occurs after lease release", event.ID)
			}
		}
		previousDigest = event.Digest
		workUnits[event.WorkUnit] = true
		agents[event.AgentID] = true
	}
	summary.WorkUnits = len(workUnits)
	summary.Agents = len(agents)
	if len(journal.Events) > 0 {
		last := journal.Events[len(journal.Events)-1]
		summary.LastSeq, summary.LastDigest = last.Sequence, last.Digest
	} else {
		summary.Status = "empty"
		summary.LastDigest = recoveryJournalRootDigest
	}
	return summary, nil
}

func recoveryEventDigest(event RecoveryJournalEvent) string {
	return digestCanonical(recoveryJournalEventBody{
		Sequence: event.Sequence, ID: event.ID, WorkUnit: event.WorkUnit,
		AgentID: event.AgentID, LeaseID: event.LeaseID, Kind: event.Kind,
		ObservedAt: event.ObservedAt, HeadSHA: strings.ToLower(event.HeadSHA),
		EvidenceRefs: append([]string(nil), event.EvidenceRefs...),
		PreviousDigest: event.PreviousDigest,
	})
}

func validateRecoveryEvidenceRefs(label string, refs []string) error {
	if len(refs) > 16 {
		return fmt.Errorf("%s has too many evidence refs", label)
	}
	seen := map[string]bool{}
	for _, ref := range refs {
		if strings.TrimSpace(ref) == "" || len(ref) > 256 || seen[ref] {
			return fmt.Errorf("%s has invalid or duplicate evidence ref", label)
		}
		seen[ref] = true
	}
	return nil
}

func runResumeCheckpoint(root, workUnit string, state CurrentState, graph WorkGraph) (ResumeCheckpoint, error) {
	if workUnit == "" {
		workUnit = state.Execution.ActiveWorkUnit
	}
	if _, ok := graphNode(graph, workUnit); !ok {
		return ResumeCheckpoint{}, fmt.Errorf("unknown checkpoint work unit %s", workUnit)
	}
	_, claims, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return ResumeCheckpoint{}, err
	}
	var journal RecoveryJournal
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(recoveryJournalPath)), &journal); err != nil {
		return ResumeCheckpoint{}, err
	}
	if _, err := auditRecoveryJournal(journal, claims, graph); err != nil {
		return ResumeCheckpoint{}, err
	}
	return compileResumeCheckpoint(journal, workUnit), nil
}

func compileResumeCheckpoint(journal RecoveryJournal, workUnit string) ResumeCheckpoint {
	checkpoint := ResumeCheckpoint{Status: "empty", WorkUnit: workUnit}
	var events []RecoveryJournalEvent
	for _, event := range journal.Events {
		if event.WorkUnit == workUnit {
			events = append(events, event)
		}
	}
	checkpoint.Events = len(events)
	if len(events) == 0 {
		return checkpoint
	}
	last := events[len(events)-1]
	checkpoint.Status = resumeStatusForEvent(last.Kind)
	checkpoint.LastSequence = last.Sequence
	checkpoint.LastEventID = last.ID
	checkpoint.LastEventKind = last.Kind
	checkpoint.LastEventDigest = last.Digest
	checkpoint.AgentID = last.AgentID
	checkpoint.LeaseID = last.LeaseID
	checkpoint.HeadSHA = last.HeadSHA
	checkpoint.EvidenceRefs = append([]string(nil), last.EvidenceRefs...)
	start := len(events) - 5
	if start < 0 {
		start = 0
	}
	for _, event := range events[start:] {
		checkpoint.RecentEventIDs = append(checkpoint.RecentEventIDs, event.ID)
	}
	return checkpoint
}

func resumeStatusForEvent(kind string) string {
	switch kind {
	case "lease-acquired":
		return "claimed"
	case "implementation-started":
		return "implementing"
	case "commit-observed":
		return "implemented"
	case "validation-started":
		return "testing"
	case "validation-passed":
		return "validated"
	case "validation-failed":
		return "needs-fix"
	case "handoff":
		return "handoff"
	case "lease-released":
		return "released"
	default:
		return "unknown"
	}
}

func sortRecoveryEvidenceRefs(refs []string) []string {
	out := append([]string(nil), refs...)
	sort.Strings(out)
	return out
}
