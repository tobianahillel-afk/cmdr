package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"strings"
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

func TestOwnerSnapshotRejectsOversizedLinkedSourceBeforeAuthority(t *testing.T) {
	checker, owner := ownerFixture()
	e := owner.entries["x"]
	e.SourceRef = strings.Repeat("A", MaxInputBytes+1)
	owner.entries["x"] = e
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
	if got.Status != "deny" || len(got.Entries) != 0 {
		t.Fatalf("oversized owner-controlled source reference reached output: %+v", got)
	}
	if len(checker.calls) != 2 || len(owner.calls) != 1 {
		t.Fatalf("oversized owner-controlled source reference reached source authorization: checker=%+v owner=%+v", checker.calls, owner.calls)
	}
}

func TestOwnerSnapshotRejectsOversizedScopeBeforeProviderCalls(t *testing.T) {
	for _, variant := range []string{"case", "tenant", "environment"} {
		t.Run(variant, func(t *testing.T) {
			checker, owner := ownerFixture()
			scope := allowedFixture().Scope
			caseID := "case-123"
			switch variant {
			case "case":
				caseID = strings.Repeat("C", MaxInputBytes+1)
			case "tenant":
				scope.Tenant = strings.Repeat("T", MaxInputBytes+1)
			case "environment":
				scope.Environment = strings.Repeat("E", MaxInputBytes+1)
			}
			got := ProjectFromOwnerSnapshots(context.Background(), scope, caseID, nil, checker, owner)
			if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 0 || len(owner.calls) != 0 {
				t.Fatalf("oversized authority selector allowed provider I/O: %+v checker=%+v owner=%+v", got, checker.calls, owner.calls)
			}
		})
	}
}

func TestOwnerSnapshotBoundsCumulativeSourceSelectorsBeforeSecondLookup(t *testing.T) {
	checker, owner := ownerFixture()
	first := owner.entries["x"]
	first.SourceRef = strings.Repeat("a", MaxInputBytes/2)
	owner.entries["x"] = first
	second := first
	second.ID = "y"
	second.SourceRef = strings.Repeat("b", MaxInputBytes/2)
	owner.entries["y"] = second
	checker.decisions["timeline-entry:y"] = "allow"
	checker.decisions["source:"+first.SourceRef] = "allow"
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x", "y"}, checker, owner)
	if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 4 || len(owner.calls) != 2 {
		t.Fatalf("cumulative owner selectors reached source authority: %+v checker=%d resolver=%d", got, len(checker.calls), len(owner.calls))
	}
	for _, call := range checker.calls {
		if call.Resource == "source" && call.ObjectID == second.SourceRef {
			t.Fatal("second oversized cumulative owner source selector reached provider")
		}
	}
}

func TestOwnerSnapshotSourceSelectorBudgetDoesNotChargeDeniedTimelineID(t *testing.T) {
	checker, owner := ownerFixture()
	first := owner.entries["x"]
	second := first
	second.ID = "y"
	second.SourceRef = strings.Repeat("secret", MaxInputBytes/6+1)
	owner.entries["y"] = second
	checker.decisions["timeline-entry:y"] = "deny"
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x", "y"}, checker, owner)
	if got.Status != "filtered" || len(got.Entries) != 1 || len(owner.calls) != 1 || len(checker.calls) != 4 {
		t.Fatalf("denied Timeline Entry consumed source budget or was resolved: %+v resolver=%d checker=%d", got, len(owner.calls), len(checker.calls))
	}
}

func TestOwnerRejectsInvisibleSelectorsBeforeAnyAuthority(t *testing.T) {
	for _, tc := range []struct{ name, value string }{
		{"embedded-null", "x\x00hidden"},
		{"embedded-linebreak", "x\nsecond"},
		{"zero-width", "x\u200bhidden"},
		{"bidi-override", "x\u202ehidden"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, position := range []string{"tenant", "environment", "case", "id"} {
				checker, owner := ownerFixture()
				scope := allowedFixture().Scope
				caseID := "case-123"
				ids := []string{"x"}
				switch position {
				case "tenant":
					scope.Tenant = tc.value
				case "environment":
					scope.Environment = tc.value
				case "case":
					caseID = tc.value
				case "id":
					ids = []string{tc.value}
				}
				got := ProjectFromOwnerSnapshots(context.Background(), scope, caseID, ids, checker, owner)
				if got.Status != "deny" || len(checker.calls) != 0 || len(owner.calls) != 0 {
					t.Fatalf("%s reached provider for %q: %+v checker=%+v owner=%+v", position, tc.value, got, checker.calls, owner.calls)
				}
			}
		})
	}
}

func TestOwnerRejectsInvisibleSourceReferenceBeforeSourceProvider(t *testing.T) {
	for _, sourceRef := range []string{"owner:x\x00hidden", "owner:x\nspoof", "owner:\u202eabcd", "owner:\u200bsecret"} {
		checker, owner := ownerFixture()
		entry := owner.entries["x"]
		entry.SourceRef = sourceRef
		owner.entries["x"] = entry
		got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
		if got.Status != "deny" || len(checker.calls) != 2 || len(owner.calls) != 1 {
			t.Fatalf("malformed source reached authority: %+v checker=%+v", got, checker.calls)
		}
	}
}
