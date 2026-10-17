package main

import (
	"strings"
	"testing"
)

func journalClaim() WorkLease {
	return WorkLease{
		ID: "LEASE-TEST0001", WorkUnit: "E7-REC-001B", AgentID: "AGENT-test",
		BaseHeadSHA: strings.Repeat("a", 40),
		AcquiredAt:  "2026-09-23T10:00:00Z", ExpiresAt: "2026-09-23T14:00:00Z",
	}
}

func journalGraph() WorkGraph {
	return WorkGraph{SchemaVersion: 1, GraphKind: "cmrd-engineering-work-graph", Nodes: []WorkNode{
		{ID: "E7-REC-001B", Title: "journal", Type: "sublot", Status: "READY"},
	}}
}

func journalEvent(seq int, id, kind, at, previous string) RecoveryJournalEvent {
	event := RecoveryJournalEvent{
		Sequence: seq, ID: id, WorkUnit: "E7-REC-001B", AgentID: "AGENT-test",
		LeaseID: "LEASE-TEST0001", Kind: kind, ObservedAt: at,
		HeadSHA: strings.Repeat("b", 40), PreviousDigest: previous,
		EvidenceRefs: []string{"commit:" + strings.Repeat("b", 40)},
	}
	event.Digest = recoveryEventDigest(event)
	return event
}

func validJournal() RecoveryJournal {
	first := journalEvent(1, "EVT-TEST0001", "lease-acquired", "2026-09-23T10:00:00Z", recoveryJournalRootDigest)
	second := journalEvent(2, "EVT-TEST0002", "implementation-started", "2026-09-23T10:05:00Z", first.Digest)
	return RecoveryJournal{SchemaVersion: 1, RegistryKind: "engineering-recovery-journal", Events: []RecoveryJournalEvent{first, second}}
}

func TestRecoveryJournalHashChainAndCheckpoint(t *testing.T) {
	journal := validJournal()
	claims := WorkLeaseRegistry{SchemaVersion: 1, RegistryKind: "engineering-work-leases", Claims: []WorkLease{journalClaim()}}
	summary, err := auditRecoveryJournal(journal, claims, journalGraph())
	if err != nil || summary.Events != 2 || summary.Status != "valid" {
		t.Fatalf("unexpected journal audit: %#v err=%v", summary, err)
	}
	checkpoint := compileResumeCheckpoint(journal, "E7-REC-001B")
	if checkpoint.Status != "implementing" || checkpoint.LastEventID != "EVT-TEST0002" || len(checkpoint.RecentEventIDs) != 2 {
		t.Fatalf("unexpected checkpoint: %#v", checkpoint)
	}
}

func TestRecoveryJournalRejectsTamperedEvent(t *testing.T) {
	journal := validJournal()
	journal.Events[0].HeadSHA = strings.Repeat("c", 40)
	claims := WorkLeaseRegistry{Claims: []WorkLease{journalClaim()}}
	if _, err := auditRecoveryJournal(journal, claims, journalGraph()); err == nil {
		t.Fatal("expected tampered digest rejection")
	}
}

func TestRecoveryJournalRejectsSequenceGap(t *testing.T) {
	journal := validJournal()
	journal.Events[1].Sequence = 3
	journal.Events[1].Digest = recoveryEventDigest(journal.Events[1])
	claims := WorkLeaseRegistry{Claims: []WorkLease{journalClaim()}}
	if _, err := auditRecoveryJournal(journal, claims, journalGraph()); err == nil {
		t.Fatal("expected sequence gap rejection")
	}
}

func TestRecoveryJournalRejectsEventOutsideLease(t *testing.T) {
	journal := validJournal()
	journal.Events[1].ObservedAt = "2026-09-23T15:00:00Z"
	journal.Events[1].Digest = recoveryEventDigest(journal.Events[1])
	claims := WorkLeaseRegistry{Claims: []WorkLease{journalClaim()}}
	if _, err := auditRecoveryJournal(journal, claims, journalGraph()); err == nil {
		t.Fatal("expected event outside lease rejection")
	}
}

func TestResumeCheckpointIsBoundedToFiveRecentEvents(t *testing.T) {
	journal := RecoveryJournal{SchemaVersion: 1, RegistryKind: "engineering-recovery-journal"}
	previous := recoveryJournalRootDigest
	for i := 1; i <= 8; i++ {
		event := journalEvent(i, "EVT-TEST000"+string(rune('0'+i)), "commit-observed", "2026-09-23T10:10:00Z", previous)
		journal.Events = append(journal.Events, event)
		previous = event.Digest
	}
	checkpoint := compileResumeCheckpoint(journal, "E7-REC-001B")
	if len(checkpoint.RecentEventIDs) != 5 {
		t.Fatalf("checkpoint is not bounded: %#v", checkpoint.RecentEventIDs)
	}
}

func TestEmptyJournalCheckpointIsExplicit(t *testing.T) {
	checkpoint := compileResumeCheckpoint(RecoveryJournal{}, "E7-REC-001B")
	if checkpoint.Status != "empty" || checkpoint.Events != 0 || checkpoint.WorkUnit != "E7-REC-001B" {
		t.Fatalf("unexpected empty checkpoint: %#v", checkpoint)
	}
}
