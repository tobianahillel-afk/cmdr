package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"strings"
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
		"case:case-123":    "allow",
		"timeline-entry:x": "allow",
		"source:owner:x":   "allow",
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
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "case", ObjectID: "case-123", Permission: caseReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "timeline-entry", ObjectID: "x", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "source", ObjectID: "owner:x", Permission: sourceReadPermission},
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

type caseMembershipChecker struct {
	calls []AccessRequest
}

func (s *caseMembershipChecker) CheckAccess(_ context.Context, r AccessRequest) (string, error) {
	s.calls = append(s.calls, r)
	if r.CaseID != "case-123" {
		return "deny", nil
	}
	if r.Resource == "case" {
		if r.ObjectID == r.CaseID && r.Permission == caseReadPermission {
			return "allow", nil
		}
		return "deny", nil
	}
	// Entry x belongs to case-123, and only source owner:x is linked.
	if r.Resource == "timeline-entry" && r.ObjectID == "x" && r.Permission == timelineReadPermission {
		return "allow", nil
	}
	if r.Resource == "source" && r.ObjectID == "owner:x" && r.Permission == sourceReadPermission {
		return "allow", nil
	}
	return "deny", nil
}

func TestAuthorizationRequestsBindAllObjectsToSelectedCase(t *testing.T) {
	checker := &caseMembershipChecker{}
	got := ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-123", checker)
	if got.Status != "allow" || len(got.Entries) != 1 || len(checker.calls) != 3 {
		t.Fatalf("Case-bound provider rejected expected entry: %+v calls=%+v", got, checker.calls)
	}
	for _, request := range checker.calls {
		if request.CaseID != "case-123" {
			t.Fatalf("access obligation lost case membership context: %+v", request)
		}
	}
	// The same entry ID must not be smuggled into another Case by reusing
	// a checker permission obtained for the legitimate Case.
	other := &caseMembershipChecker{}
	out := ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-other", other)
	if out.Status != "deny" || len(out.Entries) != 0 || len(other.calls) != 1 {
		t.Fatalf("cross-case object splice permitted: %+v calls=%+v", out, other.calls)
	}
}

type deniedRelatedCaseChecker struct {
	calls []AccessRequest
}

func (s *deniedRelatedCaseChecker) CheckAccess(_ context.Context, r AccessRequest) (string, error) {
	s.calls = append(s.calls, r)
	if r.Resource == "case" && r.CaseID == "case-123" {
		return "allow", nil
	}
	// A timeline object or its source may be readable in a tenant, but
	// its owner can still reject its association with this Case.
	if r.Resource == "timeline-entry" && r.CaseID == "case-123" {
		return "deny", nil
	}
	return "allow", nil
}

func TestAuthorizationCaseMembershipDenialCannotTriggerSourceLookup(t *testing.T) {
	checker := &deniedRelatedCaseChecker{}
	out := ProjectWithAccessChecker(context.Background(), allowedFixture(), "case-123", checker)
	if out.Status != "filtered" || len(out.Entries) != 0 || len(checker.calls) != 2 {
		t.Fatalf("unrelated Timeline Entry leaked or source checked: %+v calls=%+v", out, checker.calls)
	}
	if checker.calls[1].CaseID != "case-123" || checker.calls[1].Resource != "timeline-entry" {
		t.Fatalf("case membership request not passed to authority: %+v", checker.calls[1])
	}
}

func TestAuthorizedProjectionBoundsSourceReferenceBeforeProvider(t *testing.T) {
	in := allowedFixture()
	in.Entries[0].SourceRef = strings.Repeat("untrusted:", MaxInputBytes/10+1)
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "deny" || len(got.Entries) != 0 {
		t.Fatalf("oversized source reference allowed: %+v", got)
	}
	if len(checker.calls) != 2 || checker.calls[0].Resource != "case" || checker.calls[1].Resource != "timeline-entry" {
		t.Fatalf("oversized SourceRef reached source authority: %+v", checker.calls)
	}
}

func TestAuthorizedProjectionBoundsEnvelopeBeforeAnyProvider(t *testing.T) {
	for _, tc := range []string{"case", "tenant", "environment"} {
		t.Run(tc, func(t *testing.T) {
			in := allowedFixture()
			caseID := "case-123"
			switch tc {
			case "case":
				caseID = strings.Repeat("C", MaxInputBytes+1)
			case "tenant":
				in.Scope.Tenant = strings.Repeat("T", MaxInputBytes+1)
			case "environment":
				in.Scope.Environment = strings.Repeat("E", MaxInputBytes+1)
			}
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), in, caseID, checker)
			if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 0 {
				t.Fatalf("unbounded envelope reached provider: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestLegacyAccessCheckerBoundsTimelineSelectorsBeforeAuthority(t *testing.T) {
	// Even the deprecated caller-snapshot entrypoint must not send unbounded
	// object IDs to Case/Timeline permission providers before final projection.
	for _, tc := range []struct {
		name string
		edit func(*Input)
	}{
		{"single-oversized-id", func(i *Input) { i.Entries[0].ID = strings.Repeat("A", MaxInputBytes+1) }},
		{"cumulative-oversized-ids", func(i *Input) {
			second := i.Entries[0]
			i.Entries[0].ID = strings.Repeat("a", MaxInputBytes/2)
			second.ID = strings.Repeat("b", MaxInputBytes/2)
			i.Entries = append(i.Entries, second)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			tc.edit(&in)
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if got.Status != "deny" || len(got.Entries) != 0 ||
				len(checker.calls) != 0 {
				t.Fatalf("oversized Timeline ID reached authority: result=%+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestLegacyAccessCheckerIgnoresOpaqueEntrySizeDuringSelectorBudget(t *testing.T) {
	// Foreign entries cannot consume the permitted tenant's selector budget,
	// nor may their identifiers be sent to the target tenant's authority.
	in := allowedFixture()
	foreign := in.Entries[0]
	foreign.Tenant = "tenant-foreign"
	foreign.ID = strings.Repeat("secret", MaxInputBytes/6+1)
	in.Entries = append(in.Entries, foreign)
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "filtered" || len(got.Entries) != 1 ||
		len(checker.calls) != 3 {
		t.Fatalf("foreign entry size interfered with accessible projection: result=%+v calls=%+v", got, checker.calls)
	}
}

func TestLegacyAccessRejectsCumulativeSourceSelectorsBeforeSecondAuthorityCall(t *testing.T) {
	in := allowedFixture()
	second := in.Entries[0]
	second.ID = "y"
	firstRef := strings.Repeat("a", MaxInputBytes/2)
	secondRef := strings.Repeat("b", MaxInputBytes/2)
	in.Entries[0].SourceRef = firstRef
	second.SourceRef = secondRef
	in.Entries = append(in.Entries, second)
	checker := allowedAccessChecker()
	checker.decisions["source:"+firstRef] = "allow"
	checker.decisions["timeline-entry:y"] = "allow"
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 4 {
		t.Fatalf("cumulative source selectors reached authority or yielded entries: %+v calls=%d", got, len(checker.calls))
	}
	for _, call := range checker.calls {
		if call.Resource == "source" && call.ObjectID == secondRef {
			t.Fatal("second oversized cumulative source selector reached provider")
		}
	}
}

func TestLegacySourceSelectorBudgetIgnoresUnauthorizedTimelineEntries(t *testing.T) {
	in := allowedFixture()
	other := in.Entries[0]
	other.ID = "denied-y"
	other.SourceRef = strings.Repeat("protected", MaxInputBytes/9+1)
	in.Entries = append(in.Entries, other)
	checker := allowedAccessChecker()
	checker.decisions["timeline-entry:denied-y"] = "deny"
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "filtered" || len(got.Entries) != 1 || len(checker.calls) != 4 {
		t.Fatalf("denied Timeline Entry affected visible source budget: %+v calls=%d", got, len(checker.calls))
	}
}
