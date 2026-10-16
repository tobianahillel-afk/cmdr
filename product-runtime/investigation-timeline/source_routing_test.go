package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type stubSourceOwnerResolver struct {
	values  map[string]SourceOwnerIdentity
	calls   []AccessRequest
	failure bool
}

func (s *stubSourceOwnerResolver) ResolveSource(_ context.Context, req AccessRequest) (SourceOwnerIdentity, error) {
	s.calls = append(s.calls, req)
	if s.failure {
		return SourceOwnerIdentity{}, errors.New("source owner unavailable")
	}
	return s.values[req.ObjectID], nil
}
func ownerSourceFixture(kind string) (*stubAccessChecker, *ownedSnapshotResolver, *stubSourceOwnerResolver) {
	checker, owner := ownerFixture()
	identity := SourceOwnerIdentity{
		Reference: "owner:x", Kind: kind, ObjectID: "resolved-123",
		Tenant: "tenant-a", Environment: "prod", CaseID: "case-123",
	}
	source := &stubSourceOwnerResolver{values: map[string]SourceOwnerIdentity{"owner:x": identity}}
	perm, _ := canonicalSourceReadPermission(kind)
	checker.decisions[kind+":resolved-123"] = "allow"
	_ = perm
	return checker, owner, source
}

func TestOwnerSnapshotTypedSourcesRouteAllCanonicalPermissions(t *testing.T) {
	for _, kind := range []string{"evidence", "telemetry-event", "decision", "incident"} {
		t.Run(kind, func(t *testing.T) {
			checker, owner, source := ownerSourceFixture(kind)
			got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner, source)
			if got.Status != "allow" || len(got.Entries) != 1 || len(checker.calls) != 3 || len(source.calls) != 1 {
				t.Fatalf("route failed: %+v checker=%+v source=%+v", got, checker.calls, source.calls)
			}
			permission, _ := canonicalSourceReadPermission(kind)
			want := AccessRequest{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: kind, ObjectID: "resolved-123", Permission: permission}
			if checker.calls[2] != want {
				t.Fatalf("wrong canonical owner permission: got=%+v want=%+v", checker.calls[2], want)
			}
			if source.calls[0].Resource != "source-reference" || source.calls[0].ObjectID != "owner:x" || source.calls[0].CaseID != "case-123" {
				t.Fatal("resolver lost Case-scoped source identity")
			}
		})
	}
}

func TestOwnerSnapshotTypedSourcesFailClosedForUntrustedOwnerResponses(t *testing.T) {
	for _, mutation := range []struct {
		name   string
		change func(*SourceOwnerIdentity)
	}{
		{"bad-kind", func(x *SourceOwnerIdentity) { x.Kind = "source" }},
		{"empty-type", func(x *SourceOwnerIdentity) { x.Kind = "" }},
		{"foreign-tenant", func(x *SourceOwnerIdentity) { x.Tenant = "tenant-b" }},
		{"foreign-case", func(x *SourceOwnerIdentity) { x.CaseID = "case-other" }},
		{"foreign-environment", func(x *SourceOwnerIdentity) { x.Environment = "test" }},
		{"forged-reference", func(x *SourceOwnerIdentity) { x.Reference = "owner:other" }},
		{"missing-id", func(x *SourceOwnerIdentity) { x.ObjectID = "" }},
		{"wildcard-id", func(x *SourceOwnerIdentity) { x.ObjectID = " * " }},
		{"invisible-id", func(x *SourceOwnerIdentity) { x.ObjectID = "\u200bsecret" }},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			checker, owner, source := ownerSourceFixture("evidence")
			v := source.values["owner:x"]
			mutation.change(&v)
			source.values["owner:x"] = v
			got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner, source)
			if got.Status != "deny" || len(got.Entries) != 0 || len(source.calls) != 1 || len(checker.calls) != 2 {
				t.Fatalf("untrusted source accepted/leaked: %+v checker=%+v source=%+v", got, checker.calls, source.calls)
			}
		})
	}
	checker, owner, source := ownerSourceFixture("evidence")
	source.failure = true
	if got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner, source); got.Status != "deny" || len(checker.calls) != 2 {
		t.Fatal("source resolver exception did not fail closed")
	}
	checker, owner, _ = ownerSourceFixture("evidence")
	if got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner, nil); got.Status != "deny" || len(checker.calls) != 0 {
		t.Fatal("nil typed resolver allowed")
	}
}

func TestOwnerSnapshotTypedSourcesNoEarlyACLWhenLaterIdentityInvalid(t *testing.T) {
	checker, owner, source := ownerSourceFixture("evidence")
	y := owner.entries["x"]
	y.ID = "y"
	y.SourceRef = "owner:y"
	owner.entries["y"] = y
	checker.decisions["timeline-entry:y"] = "allow"
	source.values["owner:y"] = SourceOwnerIdentity{Reference: "owner:y", Kind: "incident", ObjectID: "inc-1", Tenant: "tenant-other", Environment: "prod", CaseID: "case-123"}
	got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x", "y"}, checker, owner, source)
	if got.Status != "deny" || len(checker.calls) != 3 || len(source.calls) != 2 || len(owner.calls) != 2 {
		t.Fatalf("early typed source ACL leaked: %+v checker=%+v", got, checker.calls)
	}
}

func TestOwnerSnapshotTypedSourcesNeverResolveDeniedEntries(t *testing.T) {
	checker, owner, source := ownerSourceFixture("decision")
	checker.decisions["timeline-entry:x"] = "deny"
	got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner, source)
	if got.Status != "filtered" || len(got.Entries) != 0 || len(source.calls) != 0 || len(owner.calls) != 0 || len(checker.calls) != 2 {
		t.Fatalf("denied source queried: %+v", got)
	}
}

func TestTypedSourceRoutingIsNotCallerPermission(t *testing.T) {
	checker, owner, source := ownerSourceFixture("incident")
	checker.decisions["incident:resolved-123"] = "deny"
	forged := allowedFixture().Scope
	forged.CasePermission = "allow"
	forged.TimelinePermission = "allow"
	got := ProjectFromOwnerSnapshotsWithTypedSources(context.Background(), forged, "case-123", []string{"x"}, checker, owner, source)
	if got.Status != "filtered" || len(got.Entries) != 0 {
		t.Fatalf("caller flags bypassed source ACL: %+v", got)
	}
	if !reflect.DeepEqual(source.calls, []AccessRequest{{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "source-reference", ObjectID: "owner:x", Permission: ""}}) {
		t.Fatalf("unbound resolver context: %+v", source.calls)
	}
}
