package investigationtimeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type revokeNavigationCaseAccess struct {
	*integrationStore
	checks int
}

func (s *revokeNavigationCaseAccess) CheckAccess(ctx context.Context, r AccessRequest) (string, error) {
	if r.Resource == "case" {
		s.checks++
		// Two existing Case checks protect the Timeline read; the third
		// must protect the destination handoff against a later revocation.
		if s.checks >= 3 {
			return "deny", nil
		}
	}
	return s.integrationStore.CheckAccess(ctx, r)
}

type driftingNavigationEntryResolver struct {
	*integrationStore
	reads int
}

func (s *driftingNavigationEntryResolver) ResolveEntry(ctx context.Context, r AccessRequest) (Entry, error) {
	entry, err := s.integrationStore.ResolveEntry(ctx, r)
	s.reads++
	if s.reads > 3 && r.ObjectID == "entry-new" {
		entry.SourceRef = "changed-after-first-read"
	}
	return entry, err
}

func TestSourceNavigationDeniesCaseRevokedDuringHandoff(t *testing.T) {
	store, ctx := workspaceFixture()
	revoking := &revokeNavigationCaseAccess{integrationStore: store}
	got := OpenCaseTimelineSourceNavigation(ctx, Scope{
		Tenant: "tenant-a", Environment: "prod",
	}, "case-one", "entry-new", revoking, store, store, store)
	if got.Status != "deny" || got.Draft != nil || revoking.checks < 3 {
		t.Fatalf("late Case revocation leaked destination: %+v, checks=%d", got, revoking.checks)
	}
}

func TestSourceNavigationDeniesEntryChangedAfterProjection(t *testing.T) {
	store, ctx := workspaceFixture()
	drifting := &driftingNavigationEntryResolver{integrationStore: store}
	got := OpenCaseTimelineSourceNavigation(ctx, Scope{
		Tenant: "tenant-a", Environment: "prod",
	}, "case-one", "entry-new", store, store, drifting, store)
	if got.Status != "deny" || got.Draft != nil || drifting.reads < 4 {
		t.Fatalf("stale entry owner version/source navigation allowed: %+v, reads=%d", got, drifting.reads)
	}
}

func TestSourceNavigationDoesNotCallOwnersAfterCancellation(t *testing.T) {
	store, ctx := workspaceFixture()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	got := OpenCaseTimelineSourceNavigation(ctx, Scope{
		Tenant: "tenant-a", Environment: "prod",
	}, "case-one", "entry-new", store, store, store, store)
	if got.Status != "deny" || store.catalogCalls != 0 || store.sourceReads != 0 {
		t.Fatalf("cancelled context accessed protected owners: %+v", got)
	}
}

func TestSourceNavigationRejectsInvalidHostAndSelector(t *testing.T) {
	store, ctx := workspaceFixture()
	for _, tc := range []struct {
		name, tenant, caseID, entryID string
		checker                       AccessChecker
	}{
		{"no-case-checker", "tenant-a", "case-one", "entry-new", nil},
		{"wrong-tenant", "tenant-b", "case-one", "entry-new", store},
		{"empty-case", "tenant-a", "", "entry-new", store},
		{"wildcard-entry", "tenant-a", "case-one", "*", store},
		{"invalid-control", "tenant-a", "case-one", "entry\nnew", store},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := OpenCaseTimelineSourceNavigation(ctx, Scope{
				Tenant: tc.tenant, Environment: "prod",
			}, tc.caseID, tc.entryID, tc.checker, store, store, store)
			if got.Status != "deny" || got.Draft != nil {
				t.Fatalf("bad scope/authority produced a navigation draft: %+v", got)
			}
		})
	}
}

type sourceUnavailableSession struct{}

func (sourceUnavailableSession) ResolveCaseTimelineSession(*http.Request) (context.Context, Scope, error) {
	return nil, Scope{}, errors.New("unavailable")
}

func TestSourceNavigationHTTPRejectsFailedSessionAndBody(t *testing.T) {
	store, _ := workspaceFixture()
	h, err := NewCaseSourceNavigationHTTPHandler(sourceUnavailableSession{}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodGet, "/source-navigation?case_id=case-one&entry_id=entry-new", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized || store.catalogCalls != 0 {
		t.Fatalf("untrusted session reached owner: %d", w.Code)
	}
}

func TestMemoryCaseSourceNavigationRemainsTenantScoped(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		uri, auth string
		code      int
		kind      string
	}{
		{"/source-navigation?case_id=case-a&entry_id=new", "Bearer aa", 200, "telemetry-event"},
		{"/source-navigation?case_id=case-a&entry_id=old", "Bearer aa", 200, "evidence"},
		{"/source-navigation?case_id=case-a&entry_id=hidden", "Bearer aa", 404, ""},
		{"/source-navigation?case_id=case-b&entry_id=foreign", "Bearer aa", 404, ""},
		{"/source-navigation?case_id=case-b&entry_id=foreign", "Bearer bb", 200, "incident"},
	} {
		t.Run(tc.uri+tc.auth, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.uri, nil)
			req.Header.Set("Authorization", tc.auth)
			r := httptest.NewRecorder()
			handlers.Source.ServeHTTP(r, req)
			if r.Code != tc.code {
				t.Fatalf("HTTP %d, expected %d: %s", r.Code, tc.code, r.Body.String())
			}
			if tc.code == 200 {
				var result SourceNavigationResult
				if err := json.Unmarshal(r.Body.Bytes(), &result); err != nil ||
					result.Status != "ready" || result.Draft == nil ||
					result.Draft.OwnerKind != tc.kind {
					t.Fatalf("invalid source owner navigation: %+v %v", result, err)
				}
			} else if r.Body.String() != "{\"error\":\"request-unavailable\"}\n" {
				t.Fatalf("denial revealed source: %s", r.Body.String())
			}
		})
	}
}
