package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type stubAccessChecker struct {
	decisions map[string]string
	calls     []AccessRequest
	fail      string
	cancel    context.CancelFunc
}

func accessKey(resource, id string) string { return resource + ":" + id }

func (s *stubAccessChecker) CheckAccess(_ context.Context, r AccessRequest) (string, error) {
	s.calls = append(s.calls, r)
	if s.cancel != nil {
		s.cancel()
	}
	if accessKey(r.Resource, r.ObjectID) == s.fail {
		return "", errors.New("provider unavailable")
	}
	decision, ok := s.decisions[accessKey(r.Resource, r.ObjectID)]
	if !ok {
		return "deny", nil
	}
	return decision, nil
}

func allowedAccessChecker() *stubAccessChecker {
	return &stubAccessChecker{decisions: map[string]string{
		"case:case-123":       "allow",
		"timeline-entry:x":   "allow",
		"source:owner:x":      "allow",
	}}
}

func TestAuthorizedProjectionDerivesIndependentPermissions(t *testing.T) {
	in := allowedFixture()
	in.Scope.CasePermission = "deny"
	in.Scope.TimelinePermission = "unknown"
	in.Entries[0].SourcePermission = "deny"
	original := in
	original.Entries = append([]Entry(nil), in.Entries...)
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "allow" || len(got.Entries) != 1 || got.Entries[0].ID != "x" {
		t.Fatalf("server authorization did not override forged flags: %+v", got)
	}
	if !reflect.DeepEqual(in, original) {
		t.Fatal("provider adapter mutated caller input")
	}
	want := []AccessRequest{
		{Tenant: "tenant-a", Environment: "prod", Resource: "case", ObjectID: "case-123", Permission: caseReadPermission},
		{Tenant: "tenant-a", Environment: "prod", Resource: "timeline-entry", ObjectID: "x", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", Resource: "source", ObjectID: "owner:x", Permission: sourceReadPermission},
	}
	if !reflect.DeepEqual(checker.calls, want) {
		t.Fatalf("independent access obligations lost: got %+v want %+v", checker.calls, want)
	}
}

func TestAuthorizedProjectionCannotForgeAllowFlags(t *testing.T) {
	for _, tc := range []struct {
		name       string
		resource   string
		decision   string
		wantStatus string
		wantCalls  int
	}{
		{"case-denied", "case:case-123", "deny", "deny", 1},
		{"case-unknown", "case:case-123", "unknown", "deny", 1},
		{"case-missing", "case:case-123", "missing", "deny", 1},
		{"timeline-denied", "timeline-entry:x", "deny", "filtered", 2},
		{"timeline-missing", "timeline-entry:x", "missing", "partial", 2},
		{"timeline-unknown", "timeline-entry:x", "unknown", "filtered", 2},
		{"source-denied", "source:owner:x", "deny", "filtered", 3},
		{"source-missing", "source:owner:x", "missing", "partial", 3},
		{"source-unknown", "source:owner:x", "unknown", "filtered", 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			checker := allowedAccessChecker()
			checker.decisions[tc.resource] = tc.decision
			got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if got.Status != tc.wantStatus || len(got.Entries) != 0 || len(checker.calls) != tc.wantCalls {
				t.Fatalf("denied object leaked or next checker was invoked: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestAuthorizedProjectionFailsClosedOnProviderErrors(t *testing.T) {
	for _, resource := range []string{"case:case-123", "timeline-entry:x", "source:owner:x"} {
		t.Run(resource, func(t *testing.T) {
			checker := allowedAccessChecker()
			checker.fail = resource
			got := ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-123", checker)
			if got.Status != "deny" || len(got.Entries) != 0 {
				t.Fatalf("provider failure exposed source: %+v", got)
			}
			checker = allowedAccessChecker()
			checker.decisions[resource] = "unexpected"
			got = ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-123", checker)
			if got.Status != "deny" || len(got.Entries) != 0 {
				t.Fatalf("invalid provider outcome exposed source: %+v", got)
			}
		})
	}
}

func TestAuthorizedProjectionIgnoresForeignAndMalformedEntryIDs(t *testing.T) {
	in := allowedFixture()
	foreign := in.Entries[0]
	foreign.Tenant = "tenant-other"
	foreign.ID = "secret"
	foreign.SourceRef = "owner:secret"
	in.Entries = append(in.Entries, foreign)
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "filtered" || len(got.Entries) != 1 || len(checker.calls) != 3 {
		t.Fatalf("foreign entry was checked or disclosed: %+v, calls=%v", got, checker.calls)
	}
	for _, id := range []string{"", "*", "  "} {
		in := allowedFixture()
		in.Entries[0].ID = id
		checker := allowedAccessChecker()
		got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
		if got.Status != "filtered" || len(got.Entries) != 0 || len(checker.calls) != 1 {
			t.Fatalf("malformed entry used as wildcard: %+v; calls=%v", got, checker.calls)
		}
	}
}

func TestAuthorizedProjectionRejectsMissingSourceReferences(t *testing.T) {
	for _, ref := range []string{"", "*", "  "} {
		in := allowedFixture()
		in.Entries[0].SourceRef = ref
		checker := allowedAccessChecker()
		got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
		if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 2 {
			t.Fatalf("invalid linked source widened access: %+v; calls=%v", got, checker.calls)
		}
	}
}

func TestAuthorizedProjectionRejectsInvalidEnvelope(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*Input)
	}{
		{"missing-tenant", func(i *Input) { i.Scope.Tenant = "" }},
		{"wildcard-tenant", func(i *Input) { i.Scope.Tenant = " * " }},
		{"missing-environment", func(i *Input) { i.Scope.Environment = "" }},
		{"too-many-entries", func(i *Input) { i.Entries = make([]Entry, MaxEntries+1) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			tc.edit(&in)
			checker := allowedAccessChecker()
			if got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker); got.Status != "deny" || len(checker.calls) != 0 {
				t.Fatalf("invalid scope reached provider: %+v calls=%v", got, checker.calls)
			}
		})
	}
	for _, caseID := range []string{"", "*", " * "} {
		checker := allowedAccessChecker()
		if got := ProjectWithAccessChecker(context.Background(), allowedFixture(), caseID, checker); got.Status != "deny" || len(checker.calls) != 0 {
			t.Fatalf("invalid Case ID reached provider: %+v calls=%v", got, checker.calls)
		}
	}
	if got := ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-123", nil); got.Status != "deny" {
		t.Fatalf("nil provider allowed access: %+v", got)
	}
	if got := ProjectWithAccessChecker(nil, allowedFixture(), "case-123", allowedAccessChecker()); got.Status != "deny" {
		t.Fatalf("nil context allowed access: %+v", got)
	}
}

func TestAuthorizedProjectionStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	checker := allowedAccessChecker()
	if got := ProjectWithAccessChecker(ctx, allowedFixture(), "case-123", checker); got.Status != "deny" || len(checker.calls) != 0 {
		t.Fatalf("cancelled context reached provider: %+v calls=%v", got, checker.calls)
	}
	ctx, cancel = context.WithCancel(context.Background())
	checker = allowedAccessChecker()
	checker.cancel = cancel
	if got := ProjectWithAccessChecker(ctx, allowedFixture(), "case-123", checker); got.Status != "deny" || len(checker.calls) != 1 {
		t.Fatalf("cancelled provider still used: %+v calls=%v", got, checker.calls)
	}
}
