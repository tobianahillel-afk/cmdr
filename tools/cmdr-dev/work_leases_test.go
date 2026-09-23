package main

import (
	"strings"
	"testing"
	"time"
)

func leasePolicy() WorkLeasePolicy {
	return WorkLeasePolicy{
		SchemaVersion: 1, RegistryKind: "work-lease-policy", DefaultPolicy: "deny-conflicting-active-lease",
		DefaultLeaseMinutes: 90, MinLeaseMinutes: 15, MaxLeaseMinutes: 240,
		MaxActiveClaimsPerAgent: 1, AcquireStatuses: []string{"READY"},
	}
}

func leaseGraph(status string) WorkGraph {
	return WorkGraph{SchemaVersion: 1, GraphKind: "cmrd-engineering-work-graph", Nodes: []WorkNode{{
		ID: "E7-REC-001A", Title: "lease test", Type: "sublot", Status: status,
	}}}
}

func leaseAt(id, agent, acquired, expires string) WorkLease {
	return WorkLease{
		ID: id, WorkUnit: "E7-REC-001A", AgentID: agent,
		BaseHeadSHA: strings.Repeat("a", 40), AcquiredAt: acquired, ExpiresAt: expires,
	}
}

func mustLeaseTime(t *testing.T, value string) time.Time {
	t.Helper()
	got, err := parseCanonicalLeaseTime(value)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestExpiredLeaseDoesNotBlockAcquire(t *testing.T) {
	asOf := mustLeaseTime(t, "2026-09-23T12:00:00Z")
	registry := WorkLeaseRegistry{SchemaVersion: 1, RegistryKind: "engineering-work-leases", Claims: []WorkLease{
		leaseAt("LEASE-OLD00001", "AGENT-old", "2026-09-23T09:00:00Z", "2026-09-23T10:30:00Z"),
	}}
	summary, err := auditWorkLeases(leasePolicy(), registry, leaseGraph("READY"), asOf)
	if err != nil || summary.Expired != 1 || summary.Active != 0 {
		t.Fatalf("unexpected audit: %#v err=%v", summary, err)
	}
	result := evaluateWorkLeaseAction(leasePolicy(), registry, leaseGraph("READY"), WorkLeaseActionRequest{
		Action: "acquire", LeaseID: "LEASE-NEW00001", AgentID: "AGENT-new",
		WorkUnit: "E7-REC-001A", BaseHeadSHA: strings.Repeat("b", 40),
	}, asOf)
	if !result.Allowed || result.ProposedClaim == nil {
		t.Fatalf("expired claim blocked acquire: %#v", result)
	}
}

func TestActiveLeaseConflictsFailClosed(t *testing.T) {
	asOf := mustLeaseTime(t, "2026-09-23T12:00:00Z")
	registry := WorkLeaseRegistry{SchemaVersion: 1, RegistryKind: "engineering-work-leases", Claims: []WorkLease{
		leaseAt("LEASE-ONE00001", "AGENT-one", "2026-09-23T11:00:00Z", "2026-09-23T13:00:00Z"),
		leaseAt("LEASE-TWO00001", "AGENT-two", "2026-09-23T11:15:00Z", "2026-09-23T13:15:00Z"),
	}}
	if _, err := auditWorkLeases(leasePolicy(), registry, leaseGraph("READY"), asOf); err == nil {
		t.Fatal("expected conflicting active lease rejection")
	}
}

func TestAcquireRejectsBlockedOrTerminalWork(t *testing.T) {
	asOf := mustLeaseTime(t, "2026-09-23T12:00:00Z")
	for _, status := range []string{"BLOCKED", "VERIFIED", "DECOMPOSED"} {
		result := evaluateWorkLeaseAction(leasePolicy(), WorkLeaseRegistry{}, leaseGraph(status), WorkLeaseActionRequest{
			Action: "acquire", LeaseID: "LEASE-NEW00001", AgentID: "AGENT-new",
			WorkUnit: "E7-REC-001A", BaseHeadSHA: strings.Repeat("b", 40),
		}, asOf)
		if result.Allowed || !containsString(result.Reasons, "work-unit-not-ready") {
			t.Fatalf("status %s unexpectedly claimable: %#v", status, result)
		}
	}
}

func TestRenewAndReleaseRequireCurrentOwnerAndActiveLease(t *testing.T) {
	asOf := mustLeaseTime(t, "2026-09-23T12:00:00Z")
	claim := leaseAt("LEASE-ONE00001", "AGENT-one", "2026-09-23T11:00:00Z", "2026-09-23T13:00:00Z")
	registry := WorkLeaseRegistry{Claims: []WorkLease{claim}}
	renew := evaluateWorkLeaseAction(leasePolicy(), registry, leaseGraph("READY"), WorkLeaseActionRequest{
		Action: "renew", LeaseID: claim.ID, AgentID: claim.AgentID, DurationMinutes: 60,
	}, asOf)
	if !renew.Allowed || renew.ProposedClaim == nil || renew.ProposedClaim.ExpiresAt != "2026-09-23T13:00:00Z" {
		t.Fatalf("valid renewal rejected: %#v", renew)
	}
	wrong := evaluateWorkLeaseAction(leasePolicy(), registry, leaseGraph("READY"), WorkLeaseActionRequest{
		Action: "release", LeaseID: claim.ID, AgentID: "AGENT-other",
	}, asOf)
	if wrong.Allowed || !containsString(wrong.Reasons, "agent-does-not-own-lease") {
		t.Fatalf("non-owner release allowed: %#v", wrong)
	}
	release := evaluateWorkLeaseAction(leasePolicy(), registry, leaseGraph("READY"), WorkLeaseActionRequest{
		Action: "release", LeaseID: claim.ID, AgentID: claim.AgentID,
	}, asOf)
	if !release.Allowed || release.ProposedClaim == nil || release.ProposedClaim.ReleasedAt != "2026-09-23T12:00:00Z" {
		t.Fatalf("valid release rejected: %#v", release)
	}
}

func TestLeaseRecordRequiresCanonicalUTCAndBoundedDuration(t *testing.T) {
	asOf := mustLeaseTime(t, "2026-09-23T12:00:00Z")
	nodes := map[string]WorkNode{"E7-REC-001A": {ID: "E7-REC-001A", Status: "READY"}}
	claim := leaseAt("LEASE-ONE00001", "AGENT-one", "2026-09-23T11:00:00+00:00", "2026-09-23T13:00:00Z")
	if _, _, _, err := validateWorkLeaseRecord(leasePolicy(), claim, nodes, asOf); err == nil {
		t.Fatal("expected non-canonical timestamp rejection")
	}
	claim = leaseAt("LEASE-ONE00001", "AGENT-one", "2026-09-23T11:00:00Z", "2026-09-24T00:00:00Z")
	if _, _, _, err := validateWorkLeaseRecord(leasePolicy(), claim, nodes, asOf); err == nil {
		t.Fatal("expected oversized lease rejection")
	}
}
