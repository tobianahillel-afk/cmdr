package investigationtimeline

import (
	"context"
	"testing"
)

// sequenceSourceResolver deliberately simulates owner resolution changing
// between two uses of the same source-reference inside one projection.
type sequenceSourceResolver struct {
	responses []SourceOwnerIdentity
	calls     []AccessRequest
}

func (r *sequenceSourceResolver) ResolveSource(_ context.Context, req AccessRequest) (SourceOwnerIdentity, error) {
	r.calls = append(r.calls, req)
	if len(r.calls) > len(r.responses) {
		return SourceOwnerIdentity{}, nil
	}
	return r.responses[len(r.calls)-1], nil
}

func sameSourceRefFixture() (*stubAccessChecker, *ownedSnapshotResolver, SourceOwnerIdentity) {
	checker, owner, source := ownerSourceFixture("evidence")
	second := owner.entries["x"]
	second.ID = "y"
	owner.entries["y"] = second
	checker.decisions["timeline-entry:y"] = "allow"
	return checker, owner, source.values["owner:x"]
}

func TestTypedSourceReferenceMustResolveConsistentlyWithinOneProjection(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*SourceOwnerIdentity)
		allow  bool
	}{
		{name: "stable", change: func(*SourceOwnerIdentity) {}, allow: true},
		{name: "object-changed", change: func(id *SourceOwnerIdentity) { id.ObjectID = "different-id" }},
		{name: "owner-kind-changed", change: func(id *SourceOwnerIdentity) { id.Kind = "decision" }},
		{name: "case-changed", change: func(id *SourceOwnerIdentity) { id.CaseID = "other-case" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			checker, owner, canonical := sameSourceRefFixture()
			other := canonical
			test.change(&other)
			resolver := &sequenceSourceResolver{responses: []SourceOwnerIdentity{canonical, other}}
			got := ProjectFromOwnerSnapshotsWithTypedSources(
				context.Background(), allowedFixture().Scope, "case-123",
				[]string{"x", "y"}, checker, owner, resolver,
			)
			if len(resolver.calls) != 2 {
				t.Fatalf("both identities were not queried: %+v", resolver.calls)
			}
			for _, call := range resolver.calls {
				if call.Resource != "source-reference" || call.ObjectID != "owner:x" || call.CaseID != "case-123" {
					t.Fatalf("source request escaped Case-bound reference: %+v", call)
				}
			}
			if test.allow {
				if got.Status != "allow" || len(got.Entries) != 2 || len(checker.calls) != 5 {
					t.Fatalf("stable canonical reference rejected: %+v checker=%+v", got, checker.calls)
				}
			} else if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 3 {
				t.Fatalf("inconsistent canonical reference reached linked ACL or was exposed: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}
