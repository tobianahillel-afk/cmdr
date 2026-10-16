package investigationtimeline

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type changingSourceOwner struct {
	*integrationStore
	lookups int
	mode    string
}

func (s *changingSourceOwner) ResolveSource(ctx context.Context, req AccessRequest) (SourceOwnerIdentity, error) {
	identity, err := s.integrationStore.ResolveSource(ctx, req)
	s.lookups++
	if err != nil || s.lookups <= 3 {
		return identity, err
	}
	switch s.mode {
	case "provider-fails":
		return SourceOwnerIdentity{}, errors.New("source unavailable")
	case "foreign-tenant":
		identity.Tenant = "tenant-b"
	case "wrong-case":
		identity.CaseID = "case-two"
	case "wrong-reference":
		identity.Reference = "new-reference"
	case "unknown-kind":
		identity.Kind = "unknown"
	case "missing-object":
		identity.ObjectID = ""
	}
	return identity, nil
}

func TestSourceNavigationRejectsOwnerIdentityDriftBeforeHandoff(t *testing.T) {
	for _, mode := range []string{
		"provider-fails", "foreign-tenant", "wrong-case",
		"wrong-reference", "unknown-kind", "missing-object",
	} {
		t.Run(mode, func(t *testing.T) {
			store, ctx := workspaceFixture()
			owner := &changingSourceOwner{integrationStore: store, mode: mode}
			result := OpenCaseTimelineSourceNavigation(ctx, Scope{
				Tenant: "tenant-a", Environment: "prod",
			}, "case-one", "entry-new", store, store, store, owner)
			if result.Status != "deny" || result.Draft != nil || owner.lookups < 4 {
				t.Fatalf("rebound destination owner leaked: %+v, lookups=%d", result, owner.lookups)
			}
		})
	}
}

type revokingNavigationAccessChecker struct {
	*integrationStore
	typedChecks int
	entryChecks int
	mode        string
}

func (s *revokingNavigationAccessChecker) CheckAccess(ctx context.Context, req AccessRequest) (string, error) {
	switch req.Resource {
	case "telemetry-event", "evidence", "incident", "decision":
		s.typedChecks++
		if s.mode == "source" && s.typedChecks >= 4 {
			return "deny", nil
		}
	case "timeline-entry":
		s.entryChecks++
		if s.mode == "entry" && s.entryChecks >= 4 {
			return "deny", nil
		}
	}
	return s.integrationStore.CheckAccess(ctx, req)
}

func TestSourceNavigationRejectsDestinationAndEntryRevocation(t *testing.T) {
	for _, mode := range []string{"source", "entry"} {
		t.Run(mode, func(t *testing.T) {
			store, ctx := workspaceFixture()
			checker := &revokingNavigationAccessChecker{integrationStore: store, mode: mode}
			result := OpenCaseTimelineSourceNavigation(ctx, Scope{
				Tenant: "tenant-a", Environment: "prod",
			}, "case-one", "entry-new", checker, store, store, store)
			if result.Status != "deny" || result.Draft != nil ||
				checker.typedChecks < 3 || checker.entryChecks < 3 {
				t.Fatalf("post-projection permission revocation leaked: %+v, checks=%+v", result, checker)
			}
		})
	}
}

func TestSourceNavigationHTTPRejectsInvalidScopeAndCanceledRequest(t *testing.T) {
	store, _ := workspaceFixture()
	h, err := NewCaseSourceNavigationHTTPHandler(invalidTimelineHTTPSession{"wildcard-tenant"}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet,
		"/source-navigation?case_id=case-one&entry_id=entry-new", nil)
	reply := httptest.NewRecorder()
	h.ServeHTTP(reply, req)
	if reply.Code != http.StatusUnauthorized || store.catalogCalls != 0 {
		t.Fatalf("malformed identity reached owner: %d", reply.Code)
	}
	h, err = NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	req = httptest.NewRequest(http.MethodGet,
		"/source-navigation?case_id=case-one&entry_id=entry-new", nil)
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)
	reply = httptest.NewRecorder()
	h.ServeHTTP(reply, req)
	if reply.Code != http.StatusServiceUnavailable || store.catalogCalls != 0 {
		t.Fatalf("canceled request reached owner: %d", reply.Code)
	}
}

func TestSourceNavigationHTTPRejectsMalformedSelectorEncoding(t *testing.T) {
	h, store := sourceNavigationFixture(t)
	for _, raw := range []string{
		"case_id=case-one&entry_id=%GG",
		"case_id=case-one&entry_id=bad%00entry",
		"case_id=case-one&entry_id=",
		"case_id=case-one",
	} {
		req := httptest.NewRequest(http.MethodGet, "/source-navigation", nil)
		req.URL.RawQuery = raw
		req.Header.Set("Authorization", "Bearer analyst-a")
		reply := httptest.NewRecorder()
		h.ServeHTTP(reply, req)
		if reply.Code != http.StatusBadRequest || store.catalogCalls != 0 {
			t.Fatalf("malformed selector reached owner: %q => %d", raw, reply.Code)
		}
	}
}
