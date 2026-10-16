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
	if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 3 || len(owner.calls) != 2 {
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

func TestOwnerAggregatesTimelineIDsAndSourceRefsBeforeSourceAuthority(t *testing.T) {
	// Legacy and owner-resolved paths must share the same bounded byte
	// budget across submitted Timeline Entry IDs and resolved source refs.
	id := strings.Repeat("i", MaxInputBytes/2)
	ref := strings.Repeat("s", MaxInputBytes/2)
	checker, owner := ownerFixture()
	entry := owner.entries["x"]
	entry.ID = id
	entry.SourceRef = ref
	owner.entries = map[string]Entry{id: entry}
	checker.decisions["timeline-entry:"+id] = "allow"
	result := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
		"case-123", []string{id}, checker, owner)
	if result.Status != "deny" || len(result.Entries) != 0 ||
		len(checker.calls) != 2 || len(owner.calls) != 1 ||
		checker.calls[1].Resource != "timeline-entry" ||
		owner.calls[0].ObjectID != id {
		t.Fatalf("combined selectors reached source authority: %+v checker=%d owner=%d",
			result, len(checker.calls), len(owner.calls))
	}
}

func TestOwnerMalformedAuthorizedEntryCannotReachSourceProvider(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*Entry)
	}{
		{"missing-provenance", func(e *Entry) { e.Provenance = "" }},
		{"invalid-kind", func(e *Entry) { e.Kind = "unrecognized" }},
		{"zero-version", func(e *Entry) { e.Version = 0 }},
		{"negative-clock-uncertainty", func(e *Entry) { e.ClockUncertaintyMS = -1 }},
		{"self-supersession", func(e *Entry) { e.Supersedes = e.ID }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checker, owner := ownerFixture()
			entry := owner.entries["x"]
			tc.mutate(&entry)
			owner.entries["x"] = entry
			got := ProjectFromOwnerSnapshots(
				context.Background(), allowedFixture().Scope, "case-123",
				[]string{"x"}, checker, owner,
			)
			if got.Status != "deny" || len(got.Entries) != 0 ||
				len(owner.calls) != 1 || len(checker.calls) != 2 ||
				checker.calls[0].Resource != "case" ||
				checker.calls[1].Resource != "timeline-entry" {
				t.Fatalf("invalid owner entry reached source authority: %+v checker=%+v owner=%+v",
					got, checker.calls, owner.calls)
			}
		})
	}
}

func TestOwnerUnknownClockStillChecksSourcePermission(t *testing.T) {
	checker, owner := ownerFixture()
	entry := owner.entries["x"]
	entry.SourceTime = "unknown clock"
	owner.entries["x"] = entry
	got := ProjectFromOwnerSnapshots(
		context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner,
	)
	if got.Status != "uncertain" || len(got.Entries) != 1 ||
		len(checker.calls) != 3 || checker.calls[2].Resource != "source" {
		t.Fatalf("legitimate uncertain timestamp incorrectly rejected: %+v checker=%+v",
			got, checker.calls)
	}
}

func TestOwnerChecksAllTimelineRightsBeforeAnyOwnerOrSourceCalls(t *testing.T) {
	checker, owner := ownerFixture()
	second := owner.entries["x"]
	second.ID = "y"
	second.SourceRef = "owner:y"
	owner.entries["y"] = second
	checker.decisions["timeline-entry:y"] = "allow"
	checker.decisions["source:owner:y"] = "allow"

	got := ProjectFromOwnerSnapshots(
		context.Background(), allowedFixture().Scope, "case-123",
		[]string{"x", "y"}, checker, owner,
	)
	if got.Status != "allow" || len(got.Entries) != 2 ||
		len(owner.calls) != 2 || len(checker.calls) != 5 {
		t.Fatalf("two authorized snapshots not projected: %+v owner=%+v calls=%+v", got, owner.calls, checker.calls)
	}
	resources := []string{"case", "timeline-entry", "timeline-entry", "source", "source"}
	for i, call := range checker.calls {
		if call.Resource != resources[i] {
			t.Fatalf("source checked before all Timeline rights at call %d: %+v", i, checker.calls)
		}
	}

	checker, owner = ownerFixture()
	owner.entries["y"] = second
	checker.fail = "timeline-entry:y"
	got = ProjectFromOwnerSnapshots(
		context.Background(), allowedFixture().Scope, "case-123",
		[]string{"x", "y"}, checker, owner,
	)
	if got.Status != "deny" || len(got.Entries) != 0 ||
		len(owner.calls) != 0 || len(checker.calls) != 3 {
		t.Fatalf("later Timeline provider error permitted owner/source I/O: %+v owner=%+v checker=%+v", got, owner.calls, checker.calls)
	}
	for _, call := range checker.calls {
		if call.Resource == "source" {
			t.Fatal("source authority reached after a later Timeline permission failure")
		}
	}
}

func TestOwnerDeniedLaterTimelineNeverResolvesDeniedEntry(t *testing.T) {
	checker, owner := ownerFixture()
	denied := owner.entries["x"]
	denied.ID = "hidden"
	denied.SourceRef = "owner:secret"
	owner.entries["hidden"] = denied
	checker.decisions["timeline-entry:hidden"] = "deny"
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
		"case-123", []string{"x", "hidden"}, checker, owner)
	if got.Status != "filtered" || len(got.Entries) != 1 || got.Entries[0].ID != "x" ||
		len(owner.calls) != 1 || owner.calls[0].ObjectID != "x" ||
		len(checker.calls) != 4 {
		t.Fatalf("denied later Timeline metadata leaked: %+v owner=%+v checker=%+v", got, owner.calls, checker.calls)
	}
	for _, call := range checker.calls {
		if call.Resource == "source" && call.ObjectID == "owner:secret" {
			t.Fatal("unauthorized source selector reached provider")
		}
	}
}

func TestOwnerPrevalidatesAllSnapshotsBeforeSourceChecks(t *testing.T) {
	for _, malformed := range []struct {
		name   string
		mutate func(*Entry)
	}{
		{"wrong-id", func(e *Entry) { e.ID = "wrong" }},
		{"missing-provenance", func(e *Entry) { e.Provenance = "" }},
		{"invalid-kind", func(e *Entry) { e.Kind = "unrecognized" }},
		{"bad-source-selector", func(e *Entry) { e.SourceRef = " * " }},
		{"oversized-source-selector", func(e *Entry) { e.SourceRef = strings.Repeat("x", MaxInputBytes+1) }},
	} {
		t.Run(malformed.name, func(t *testing.T) {
			checker, owner := ownerFixture()
			second := owner.entries["x"]
			second.ID = "y"
			second.SourceRef = "owner:y"
			malformed.mutate(&second)
			owner.entries["y"] = second
			checker.decisions["timeline-entry:y"] = "allow"
			got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
				"case-123", []string{"x", "y"}, checker, owner)
			if got.Status != "deny" || len(got.Entries) != 0 ||
				len(owner.calls) != 2 || len(checker.calls) != 3 {
				t.Fatalf("invalid second owner leaked first source access: %+v, checker=%+v owner=%+v",
					got, checker.calls, owner.calls)
			}
			for _, call := range checker.calls {
				if call.Resource == "source" {
					t.Fatalf("source checked before all owners validated: %+v", checker.calls)
				}
			}
		})
	}
}

func TestOwnerPrevalidatesBeforeSourceButPreservesSourceDenials(t *testing.T) {
	checker, owner := ownerFixture()
	second := owner.entries["x"]
	second.ID = "y"
	second.SourceRef = "owner:y"
	owner.entries["y"] = second
	checker.decisions["timeline-entry:y"] = "allow"
	checker.decisions["source:owner:y"] = "deny"
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
		"case-123", []string{"x", "y"}, checker, owner)
	if got.Status != "filtered" || len(got.Entries) != 1 ||
		got.Entries[0].ID != "x" || len(owner.calls) != 2 || len(checker.calls) != 5 {
		t.Fatalf("source-denied snapshot incorrectly exposed/blocked: %+v, owner=%+v checker=%+v",
			got, owner.calls, checker.calls)
	}
	if checker.calls[3].Resource != "source" || checker.calls[4].Resource != "source" {
		t.Fatalf("unexpected authorization ordering: %+v", checker.calls)
	}
}

type cancelBetweenTimelineAndOwner struct {
	context.Context
	errReads int
}

// Simulates cancellation immediately after the last timeline-entry ACL
// result, before the owner-resolver phase begins. No scheduling or sleeps.
func (c *cancelBetweenTimelineAndOwner) Err() error {
	c.errReads++
	if c.errReads >= 5 {
		return context.Canceled
	}
	return nil
}

func TestOwnerCancellationBetweenAuthorizationAndResolutionDoesNotReadOwner(t *testing.T) {
	checker, owner := ownerFixture()
	ctx := &cancelBetweenTimelineAndOwner{Context: context.Background()}
	result := ProjectFromOwnerSnapshots(ctx, allowedFixture().Scope,
		"case-123", []string{"x"}, checker, owner)
	if result.Status != "deny" || len(result.Entries) != 0 {
		t.Fatalf("cancelled request did not fail closed: %+v", result)
	}
	if len(checker.calls) != 2 || len(owner.calls) != 0 {
		t.Fatalf("owner read after cancellation: checker=%+v owner=%+v", checker.calls, owner.calls)
	}
	if ctx.errReads < 5 {
		t.Fatalf("cancellation at owner boundary was never examined: reads=%d", ctx.errReads)
	}
}

func TestOwnerDeniedLongTimelineIDDoesNotExhaustVisibleSourceBudget(t *testing.T) {
	// Both IDs are within the globally enforced 8 MiB request budget.
	// Only the short, permitted ID should count against the linked-source
	// authority budget. A denied 3 MiB ID cannot deny a 6 MiB visible ref.
	deniedID := strings.Repeat("h", 3<<20)
	allowedSource := strings.Repeat("s", 6<<20)
	checker := allowedAccessChecker()
	checker.decisions["timeline-entry:"+deniedID] = "deny"
	checker.decisions["source:"+allowedSource] = "allow"
	owner := &ownedSnapshotResolver{entries: map[string]Entry{}}
	entry := allowedFixture().Entries[0]
	entry.SourceRef = allowedSource
	owner.entries["x"] = entry
	got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
		"case-123", []string{"x", deniedID}, checker, owner)
	if got.Status != "filtered" || len(got.Entries) != 1 ||
		got.Entries[0].ID != "x" || len(owner.calls) != 1 || len(checker.calls) != 4 {
		t.Fatalf("denied ID reduced trusted owner source budget: status=%s entries=%d owner=%d calls=%d",
			got.Status, len(got.Entries), len(owner.calls), len(checker.calls))
	}
	if checker.calls[3].Resource != "source" || checker.calls[3].ObjectID != allowedSource {
		t.Fatal("wrong authorized owner source checked")
	}
}
