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
	if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 3 {
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

func TestLegacySuppressesInvisibleAuthoritySelectors(t *testing.T) {
	for _, selector := range []string{"x\x00bad", "x\nspoof", "x\u200bhidden", "x\u202ebad"} {
		for _, location := range []string{"tenant", "environment", "case"} {
			in := allowedFixture()
			caseID := "case-123"
			switch location {
			case "tenant":
				in.Scope.Tenant = selector
			case "environment":
				in.Scope.Environment = selector
			case "case":
				caseID = selector
			}
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), in, caseID, checker)
			if got.Status != "deny" || len(checker.calls) != 0 {
				t.Fatalf("%s malformed selector was forwarded: %+v calls=%+v", location, got, checker.calls)
			}
		}
		in := allowedFixture()
		in.Entries[0].ID = selector
		checker := allowedAccessChecker()
		got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
		if got.Status != "filtered" || len(got.Entries) != 0 || len(checker.calls) != 1 {
			t.Fatalf("malformed Timeline ID was forwarded: %+v calls=%+v", got, checker.calls)
		}
		in = allowedFixture()
		in.Entries[0].SourceRef = selector
		checker = allowedAccessChecker()
		got = ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
		if got.Status != "deny" || len(checker.calls) != 2 {
			t.Fatalf("malformed source ID was forwarded: %+v calls=%+v", got, checker.calls)
		}
	}
}

func TestLegacyRejectsDuplicateVisibleSelectorsBeforeAnyAuthority(t *testing.T) {
	base := allowedFixture()
	base.Entries = append(base.Entries, base.Entries[0])
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), base, "case-123", checker)
	if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 0 {
		t.Fatalf("duplicate Timeline selector reached authority: %+v calls=%+v", got, checker.calls)
	}
}

func TestLegacyOpaqueDuplicateCannotSuppressVisibleAuthorization(t *testing.T) {
	for _, tc := range []struct {
		name, tenant, id, wantStatus string
	}{
		{"foreign-same-ID", "foreign-tenant", "x", "filtered"},
		{"foreign-oversized-ID", "foreign-tenant", "", "filtered"},
		{"malformed-same-tenant", "tenant-a", " * ", "filtered"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			opaque := in.Entries[0]
			opaque.Tenant = tc.tenant
			opaque.ID = tc.id
			if tc.name == "foreign-oversized-ID" {
				opaque.ID = strings.Repeat("hidden", MaxInputBytes/6+1)
			}
			in.Entries = append(in.Entries, opaque)
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if got.Status != tc.wantStatus || len(got.Entries) != 1 || got.Entries[0].ID != "x" || len(checker.calls) != 3 {
				t.Fatalf("opaque source changed visible access: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestLegacyAggregatesTimelineIDsAndSourceRefsBeforeSourceAuthority(t *testing.T) {
	// Each selector is individually below the 8 MiB limit. Together they
	// breach it. The source-access authority must never receive the ref.
	in := allowedFixture()
	id := strings.Repeat("i", MaxInputBytes/2)
	ref := strings.Repeat("s", MaxInputBytes/2)
	in.Entries[0].ID = id
	in.Entries[0].SourceRef = ref
	checker := allowedAccessChecker()
	checker.decisions["timeline-entry:"+id] = "allow"
	result := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if result.Status != "deny" || len(result.Entries) != 0 ||
		len(checker.calls) != 2 ||
		checker.calls[0].Resource != "case" ||
		checker.calls[1].Resource != "timeline-entry" {
		t.Fatalf("combined selectors bypassed source boundary: %+v calls=%+v", result, checker.calls)
	}
}

func TestLegacyMalformedAuthorizedEntryCannotReachSourceProvider(t *testing.T) {
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
			in := allowedFixture()
			tc.mutate(&in.Entries[0])
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if got.Status != "deny" || len(got.Entries) != 0 ||
				len(checker.calls) != 2 ||
				checker.calls[0].Resource != "case" ||
				checker.calls[1].Resource != "timeline-entry" {
				t.Fatalf("invalid legacy entry reached source authority: %+v calls=%+v", got, checker.calls)
			}
		})
	}
}

func TestLegacyUnknownClockStillChecksSourcePermission(t *testing.T) {
	in := allowedFixture()
	in.Entries[0].SourceTime = "unknown clock"
	checker := allowedAccessChecker()
	got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if got.Status != "uncertain" || len(got.Entries) != 1 ||
		len(checker.calls) != 3 || checker.calls[2].Resource != "source" {
		t.Fatalf("legitimate uncertain timestamp incorrectly rejected: %+v calls=%+v", got, checker.calls)
	}
}

func TestDeniedLongTimelineIDDoesNotExhaustVisibleSourceBudget(t *testing.T) {
	// Both IDs are within the globally enforced 8 MiB request budget.
	// Only the short, permitted ID should count against the linked-source
	// authority budget. A denied 3 MiB ID cannot deny a 6 MiB visible ref.
	deniedID := strings.Repeat("h", 3<<20)
	allowedSource := strings.Repeat("s", 6<<20)
	checker := allowedAccessChecker()
	checker.decisions["timeline-entry:"+deniedID] = "deny"
	checker.decisions["source:"+allowedSource] = "allow"
	in := allowedFixture()
	in.Entries[0].SourceRef = allowedSource
	denied := in.Entries[0]
	denied.ID = deniedID
	denied.SourceRef = "owner:hidden"
	in.Entries = append(in.Entries, denied)
	result := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if result.Status != "filtered" || len(result.Entries) != 1 ||
		result.Entries[0].ID != "x" || len(checker.calls) != 4 {
		t.Fatalf("denied ID reduced valid legacy source budget: status=%s entries=%d calls=%d",
			result.Status, len(result.Entries), len(checker.calls))
	}
	if checker.calls[3].Resource != "source" || checker.calls[3].ObjectID != allowedSource {
		t.Fatal("wrong authorized source checked")
	}
}
