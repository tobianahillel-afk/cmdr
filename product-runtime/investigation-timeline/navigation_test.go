package investigationtimeline

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCaseSourceNavigationReauthorizesDestinationAndPreservesReturn(t *testing.T) {
	store, ctx := workspaceFixture()
	got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant: "tenant-a", Environment: "prod"},
		"case-one", "entry-new", store, store, store, store)
	if !ok || got.Status != "ready" || got.Kind != "telemetry-event" ||
		got.ObjectID != "event-1" || got.ReturnCaseID != "case-one" ||
		got.ReturnEntryID != "entry-new" || got.EntryVersion != 2 {
		t.Fatalf("authorized navigation not available: %+v", got)
	}
	// The owner resolution and canonical destination permission MUST occur
	// after the existing Timeline integration has made its checks.
	if store.caseChecks < 3 || store.catalogCalls != 1 ||
		store.entryReads != 4 || store.sourceReads != 4 {
		t.Fatalf("navigation did not recheck owner chain: %+v", store)
	}
}

func TestCaseSourceNavigationCannotOpenDeniedOrForeignSource(t *testing.T) {
	for _, id := range []string{"entry-hidden", "entry-other", "unknown", "*"} {
		store, ctx := workspaceFixture()
		got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant: "tenant-a", Environment: "prod"},
			"case-one", id, store, store, store, store)
		if ok || got.ObjectID != "" || got.Kind != "" {
			t.Fatalf("hidden source %q leaked: %+v", id, got)
		}
	}
	store, ctx := workspaceFixture()
	got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant: "tenant-a", Environment: "prod"},
		"case-two", "entry-other", store, store, store, store)
	if ok || got.ObjectID != "" {
		t.Fatalf("foreign Case leaked: %+v", got)
	}
}

type revokeNavigationSourceChecker struct {
	*integrationStore
	count int
}

func (c *revokeNavigationSourceChecker) CheckAccess(ctx context.Context, r AccessRequest) (string, error) {
	if r.Resource == "telemetry-event" && r.ObjectID == "event-1" {
		c.count++
		if c.count > 1 {
			return "deny", nil
		}
	}
	return c.integrationStore.CheckAccess(ctx, r)
}
func TestCaseSourceNavigationRejectsRevokedDestination(t *testing.T) {
	store, ctx := workspaceFixture()
	checker := &revokeNavigationSourceChecker{integrationStore: store}
	got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant: "tenant-a", Environment: "prod"},
		"case-one", "entry-new", checker, store, store, store)
	if ok || got.ObjectID != "" || checker.count != 2 {
		t.Fatalf("revoked source returned destination: %+v count=%d", got, checker.count)
	}
}

type changedEntryOnSecondRead struct {
	*integrationStore
	count int
}

func (s *changedEntryOnSecondRead) ResolveEntry(ctx context.Context, r AccessRequest) (Entry, error) {
	value, err := s.integrationStore.ResolveEntry(ctx, r)
	if r.ObjectID == "entry-new" {
		s.count++
		if s.count > 1 {
			value.Version++
		}
	}
	return value, err
}
func TestCaseSourceNavigationRejectsStaleOwnerVersion(t *testing.T) {
	store, ctx := workspaceFixture()
	owner := &changedEntryOnSecondRead{integrationStore: store}
	got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant: "tenant-a", Environment: "prod"},
		"case-one", "entry-new", store, store, owner, store)
	if ok || got.ObjectID != "" || owner.count != 2 {
		t.Fatalf("stale owner allowed: %+v reads=%d", got, owner.count)
	}
}

func TestCaseSourceNavigationHTTPFunctionalAndAdversarial(t *testing.T) {
	store, _ := workspaceFixture()
	h, err := NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	request := func(uri, auth string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, uri, nil)
		if auth != "" {
			req.Header.Set("Authorization", auth)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		return w
	}
	response := request("/case-source-navigation?case_id=case-one&entry_id=entry-new&kind=observed&limit=2&offset=0",
		"Bearer analyst-a")
	if response.Code != http.StatusOK {
		t.Fatalf("navigation returned %d %s", response.Code, response.Body.String())
	}
	var draft navigationHTTPDraft
	if err := json.Unmarshal(response.Body.Bytes(), &draft); err != nil {
		t.Fatal(err)
	}
	if draft.Status != "ready" || draft.Kind != "telemetry-event" ||
		draft.ObjectID != "event-1" || draft.ReturnIsAuthority ||
		draft.ReturnQuery != "case_id=case-one&kind=observed&limit=2&offset=0" ||
		response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("incorrect safe route draft: %+v", draft)
	}
	for _, forbidden := range []string{"tenant-b", "src-hidden", "evidence-secret", "entry-hidden", "Bearer"} {
		if strings.Contains(response.Body.String(), forbidden) {
			t.Fatalf("private source/authority leaked: %q", forbidden)
		}
	}
	for _, tc := range []struct {
		uri, auth string
		code      int
	}{
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new", "", 401},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-hidden", "Bearer analyst-a", 404},
		{"/case-source-navigation?case_id=case-two&entry_id=entry-other", "Bearer analyst-a", 404},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-other", "Bearer analyst-a", 404},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new&tenant=tenant-b", "Bearer analyst-a", 400},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new&return_url=https://evil.invalid", "Bearer analyst-a", 400},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new&object_id=event-1", "Bearer analyst-a", 400},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new&entry_id=entry-old", "Bearer analyst-a", 400},
		{"/case-source-navigation?case_id=case-one&entry_id=entry-new&offset=7", "Bearer analyst-a", 400},
	} {
		res := request(tc.uri, tc.auth)
		if res.Code != tc.code || strings.Contains(res.Body.String(), "event-1") {
			t.Fatalf("%q returned %d: %s", tc.uri, res.Code, res.Body.String())
		}
	}
}

func TestCaseNavigationRejectsMissingAuthoritiesAndCancelledRequest(t *testing.T) {
	store, _ := workspaceFixture()
	for _, services := range []struct {
		sessions CaseTimelineSessionResolver
		checker AccessChecker
		catalog CaseTimelineCatalog
		entries EntryResolver
		sources SourceOwnerResolver
	}{
		{nil, store, store, store, store},
		{timelineHTTPSession{}, nil, store, store, store},
		{timelineHTTPSession{}, store, nil, store, store},
		{timelineHTTPSession{}, store, store, nil, store},
		{timelineHTTPSession{}, store, store, store, nil},
	} {
		if _, err := NewCaseSourceNavigationHTTPHandler(
			services.sessions, services.checker, services.catalog,
			services.entries, services.sources); err == nil {
			t.Fatal("accepted missing source navigation authority")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, ok := DraftCaseSourceNavigation(ctx, Scope{Tenant:"tenant-a",Environment:"prod"},
		"case-one", "entry-new", store, store, store, store)
	if ok || got.Status != "" || store.catalogCalls != 0 {
		t.Fatalf("cancelled source navigation reached provider: %+v", got)
	}
}

func TestCaseNavigationHTTPRejectsMalformedRequestsBeforeOwnerIO(t *testing.T) {
	for _, tc := range []struct {
		name string
		method string
		query string
		body string
		status int
	}{
		{"write-method", http.MethodPost, "case_id=case-one&entry_id=entry-new", "", 405},
		{"body-on-get", http.MethodGet, "case_id=case-one&entry_id=entry-new", "payload", 400},
		{"no-entry",http.MethodGet,"case_id=case-one","",400},
		{"wildcard-entry",http.MethodGet,"case_id=case-one&entry_id=+%2A+","",400},
		{"too-long-entry",http.MethodGet,"case_id=case-one&entry_id="+strings.Repeat("x",1025),"",400},
		{"duplicate-kind",http.MethodGet,"case_id=case-one&entry_id=entry-new&kind=observed&kind=inferred","",400},
		{"invalid-route",http.MethodGet,"case_id=case-one&entry_id=entry-new&redirect=https://evil.invalid","",400},
		{"invalid-time",http.MethodGet,"case_id=case-one&entry_id=entry-new&from=garbage","",400},
		{"invalid-escape",http.MethodGet,"case_id=case-one&entry_id=%GG","",400},
	} {
		t.Run(tc.name,func(t *testing.T){
			store,_:=workspaceFixture()
			h,err:=NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{},store,store,store,store)
			if err!=nil {t.Fatal(err)}
			req:=httptest.NewRequest(tc.method,"/case-source-navigation?case_id=case-one",nil)
			req.URL.RawQuery=tc.query
			req.Header.Set("Authorization","Bearer analyst-a")
			if tc.body!="" {req.Body=http.NoBody;req.Body=io.NopCloser(strings.NewReader(tc.body))}
			w:=httptest.NewRecorder()
			h.ServeHTTP(w,req)
			if w.Code!=tc.status || store.catalogCalls!=0 {
				t.Fatalf("invalid route reached owner: HTTP=%d calls=%d",w.Code,store.catalogCalls)
			}
		})
	}
}
