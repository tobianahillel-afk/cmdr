package main

import (
	"testing"
	"time"
)

func TestCoordinationDetectsOverlappingMutatingScopes(t *testing.T) {
	now := time.Date(2026, 9, 23, 20, 0, 0, 0, time.UTC)
	active := []activeCoordinationClaim{
		{Lease: WorkLease{ID:"LEASE-AAAAAAAA",WorkUnit:"E7-REC-001D",AgentID:"AGENT-one",Mode:"mutating"}, Manifest: WorkManifestV2{AllowedPaths:[]string{"tools/cmdr-dev/**"}}},
		{Lease: WorkLease{ID:"LEASE-BBBBBBBB",WorkUnit:"E7-REC-001C",AgentID:"AGENT-two",Mode:"mutating"}, Manifest: WorkManifestV2{AllowedPaths:[]string{"tools/cmdr-dev/recovery/**"}}},
	}
	summary := evaluateCoordinationConflicts(active, now)
	if summary.Status != "conflict" || summary.Conflicts == 0 {
		t.Fatalf("expected mutable scope conflict: %#v", summary)
	}
}

func TestCoordinationAllowsExplicitReadOnlyOverlap(t *testing.T) {
	now := time.Date(2026, 9, 23, 20, 0, 0, 0, time.UTC)
	active := []activeCoordinationClaim{
		{Lease: WorkLease{ID:"LEASE-AAAAAAAA",WorkUnit:"E7-REC-001D",AgentID:"AGENT-one",Mode:"mutating"}, Manifest: WorkManifestV2{AllowedPaths:[]string{"engineering/**"}}},
		{Lease: WorkLease{ID:"LEASE-BBBBBBBB",WorkUnit:"E7-REC-001C",AgentID:"AGENT-two",Mode:"read-only"}, Manifest: WorkManifestV2{AllowedPaths:[]string{"engineering/recovery/**"}}},
	}
	summary := evaluateCoordinationConflicts(active, now)
	if summary.Status != "clear" || summary.Conflicts != 0 || summary.ReadOnly != 1 || summary.Mutating != 1 {
		t.Fatalf("explicit read-only overlap should be allowed: %#v", summary)
	}
}

func TestNormalizedLeaseModeDefaultsToMutating(t *testing.T) {
	if got := normalizedLeaseMode(""); got != "mutating" {
		t.Fatalf("empty legacy mode must remain mutating, got %s", got)
	}
	if got := normalizedLeaseMode("read-only"); got != "read-only" {
		t.Fatalf("read-only mode rejected: %s", got)
	}
	if got := normalizedLeaseMode("writer"); got != "invalid" {
		t.Fatalf("unexpected mode acceptance: %s", got)
	}
}

func TestBuildCoordinationHandoffBindsExactHeadAndForbidsPostReleaseMutation(t *testing.T) {
	now := time.Date(2026, 9, 23, 20, 0, 0, 0, time.UTC)
	claim := WorkLease{ID:"LEASE-AAAAAAAA",WorkUnit:"E7-REC-001D",AgentID:"AGENT-one",Mode:"mutating"}
	head := "0123456789abcdef0123456789abcdef01234567"
	checkpoint := ResumeCheckpoint{Status:"ready",LastEventDigest:"abc123"}
	reconcile := RecoveryReconciliation{Outcome:"resume"}
	h := buildCoordinationHandoff(claim, head, now, checkpoint, reconcile)
	if h.HeadSHA != head || h.WorkUnit != claim.WorkUnit || h.LeaseID != claim.ID {
		t.Fatalf("handoff identity mismatch: %#v", h)
	}
	if !h.ReleaseRequired || !h.ReceiverMustReacquire || h.MutationAfterRelease {
		t.Fatalf("unsafe handoff protocol: %#v", h)
	}
	if h.Digest == "" || h.JournalEventKind != "handoff" {
		t.Fatalf("handoff lacks deterministic evidence: %#v", h)
	}
}

func TestBuildCoordinationHandoffMarksRevalidation(t *testing.T) {
	now := time.Date(2026, 9, 23, 20, 0, 0, 0, time.UTC)
	claim := WorkLease{ID:"LEASE-AAAAAAAA",WorkUnit:"E7-REC-001D",AgentID:"AGENT-one"}
	h := buildCoordinationHandoff(claim, "0123456789abcdef0123456789abcdef01234567", now, ResumeCheckpoint{Status:"ready"}, RecoveryReconciliation{Outcome:"revalidate"})
	if !h.RevalidationRequired {
		t.Fatal("handoff must preserve revalidation requirement")
	}
}
