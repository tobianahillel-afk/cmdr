package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type ownedSnapshotResolver struct {
	entries map[string]Entry
	calls   []AccessRequest
	fail    bool
}

func (r *ownedSnapshotResolver) ResolveEntry(_ context.Context, obligation AccessRequest) (Entry, error) {
	r.calls = append(r.calls, obligation)
	if r.fail {
		return Entry{}, errors.New("owner unavailable")
	}
	return r.entries[obligation.ObjectID], nil
}

func ownerFixture() (*stubAccessChecker, *ownedSnapshotResolver) {
	checker := allowedAccessChecker()
	entry := allowedFixture().Entries[0]
	entry.Provenance = "trusted-owner-forensic-provenance"
	entry.SourcePermission = "deny" // resolver-supplied permission is not an authority
	return checker, &ownedSnapshotResolver{entries: map[string]Entry{"x": entry}}
}

func TestOwnerSnapshotProjectionDerivesAllDataFromTrustedResolver(t *testing.T) {
	checker, owner := ownerFixture()
	forged := Scope{Tenant: "tenant-a", Environment: "prod", CasePermission: "deny", TimelinePermission: "unknown"}
	result := ProjectFromOwnerSnapshots(context.Background(), forged, "case-123", []string{"x"}, checker, owner)
	if result.Status != "allow" || len(result.Entries) != 1 || result.Entries[0].Provenance != "trusted-owner-forensic-provenance" {
		t.Fatalf("caller control of snapshot data or flags escaped into output: %+v", result)
	}
	want := []AccessRequest{
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "case", ObjectID: "case-123", Permission: caseReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "timeline-entry", ObjectID: "x", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "source", ObjectID: "owner:x", Permission: sourceReadPermission},
	}
	if !reflect.DeepEqual(checker.calls, want) || len(owner.calls) != 1 || owner.calls[0] != want[1] {
		t.Fatalf("wrong Case, timeline, owner resolver or source sequence: checker=%+v resolver=%+v", checker.calls, owner.calls)
	}
}

func TestOwnerSnapshotProjectionNeverResolvesDeniedEntries(t *testing.T) {
	for _, decision := range []struct{ value, status string }{
		{"deny", "filtered"}, {"unknown", "filtered"}, {"missing", "partial"},
	} {
		t.Run(decision.value, func(t *testing.T) {
			checker, owner := ownerFixture()
			checker.decisions["timeline-entry:x"] = decision.value
			got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
			if got.Status != decision.status || len(got.Entries) != 0 || len(owner.calls) != 0 || len(checker.calls) != 2 {
				t.Fatalf("denied metadata resolved or disclosed: got=%+v resolver=%+v checker=%+v", got, owner.calls, checker.calls)
			}
		})
	}
}

func TestOwnerSnapshotProjectionFailsAtomicallyOnInvalidOwnerData(t *testing.T) {
	mutations := []struct {
		name   string
		mutate func(*ownedSnapshotResolver)
	}{
		{"mismatched-id", func(r *ownedSnapshotResolver) { e := r.entries["x"]; e.ID = "other"; r.entries["x"] = e }},
		{"cross-tenant", func(r *ownedSnapshotResolver) { e := r.entries["x"]; e.Tenant = "tenant-b"; r.entries["x"] = e }},
		{"no-linked-source", func(r *ownedSnapshotResolver) { e := r.entries["x"]; e.SourceRef = ""; r.entries["x"] = e }},
		{"wildcard-linked-source", func(r *ownedSnapshotResolver) { e := r.entries["x"]; e.SourceRef = " * "; r.entries["x"] = e }},
		{"owner-error", func(r *ownedSnapshotResolver) { r.fail = true }},
	}
	for _, tc := range mutations {
		t.Run(tc.name, func(t *testing.T) {
			checker, owner := ownerFixture()
			tc.mutate(owner)
			got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
			if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 2 {
				t.Fatalf("owner mismatch used as authority: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestOwnerSnapshotProjectionStopsOnProviderFailureAndInvalidSelection(t *testing.T) {
	for _, selected := range [][]string{{""}, {" * "}, {"x", "x"}} {
		checker, owner := ownerFixture()
		got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", selected, checker, owner)
		if got.Status != "deny" || len(got.Entries) != 0 || len(owner.calls) != 0 {
			t.Fatalf("invalid selector reached owner: %+v owner=%+v", got, owner.calls)
		}
	}
	for _, failed := range []string{"case:case-123", "timeline-entry:x", "source:owner:x"} {
		checker, owner := ownerFixture()
		checker.fail = failed
		got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
		if got.Status != "deny" || len(got.Entries) != 0 {
			t.Fatalf("provider error exposed data at %s: %+v", failed, got)
		}
	}
}
func TestOwnerSnapshotProjectionRespectsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	checker, owner := ownerFixture()
	got := ProjectFromOwnerSnapshots(ctx, allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
	if got.Status != "deny" || len(checker.calls) != 0 || len(owner.calls) != 0 {
		t.Fatalf("cancelled request called authority: %+v", got)
	}
}
