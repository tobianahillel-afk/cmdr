package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func sourceNavigationFixture(t *testing.T) (http.Handler, *integrationStore) {
	t.Helper()
	store, _ := workspaceFixture()
	handler, err := NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	return handler, store
}

func TestSourceNavigationAuthorizedCaseToTypedDestination(t *testing.T) {
	handler, store := sourceNavigationFixture(t)
	for _, tc := range []struct{ entry, kind, object string }{
		{"entry-new", "telemetry-event", "event-1"},
		{"entry-old", "evidence", "evidence-1"},
	} {
		result := httpWorkspaceGet(handler,
			"/source-navigation?case_id=case-one&entry_id="+tc.entry,
			"Bearer analyst-a", map[string]string{
				"X-Tenant-Ref":        "tenant-b",
				"X-Source-Permission": "allow",
				"X-Source-Object":     "evidence-secret",
			})
		if result.Code != http.StatusOK {
			t.Fatalf("%s status=%d: %s", tc.entry, result.Code, result.Body.String())
		}
		var data SourceNavigationResult
		if err := json.Unmarshal(result.Body.Bytes(), &data); err != nil {
			t.Fatal(err)
		}
		if data.Status != "ready" || data.Draft == nil ||
			data.Draft.OwnerKind != tc.kind || data.Draft.OwnerObjectID != tc.object ||
			data.Draft.Return.CaseID != "case-one" || data.Draft.Return.EntryID != tc.entry ||
			data.Draft.Return.Version < 1 {
			t.Fatalf("wrong authorized navigation: %+v", data)
		}
		if result.Header().Get("Cache-Control") != "no-store" ||
			result.Header().Get("Referrer-Policy") != "no-referrer" {
			t.Fatal("source navigation response is cacheable or referrer-exposing")
		}
	}
	if store.caseChecks < 4 || store.sourceReads < 4 {
		t.Fatalf("owner permissions not rechecked: Case=%d source=%d", store.caseChecks, store.sourceReads)
	}
}

func TestSourceNavigationMasksDeniedForeignAndUnknownEntries(t *testing.T) {
	for _, tc := range []struct {
		name, path, auth string
		code             int
	}{
		{"hidden-source", "/source-navigation?case_id=case-one&entry_id=entry-hidden", "Bearer analyst-a", 404},
		{"other-tenant", "/source-navigation?case_id=case-two&entry_id=entry-other", "Bearer analyst-a", 404},
		{"unknown-entry", "/source-navigation?case_id=case-one&entry_id=no-such-id", "Bearer analyst-a", 404},
		{"anonymous", "/source-navigation?case_id=case-one&entry_id=entry-new", "", 401},
		{"bad-session", "/source-navigation?case_id=case-one&entry_id=entry-new", "Bearer invalid", 401},
		{"spoof-tenant", "/source-navigation?case_id=case-one&entry_id=entry-new&tenant=tenant-b", "Bearer analyst-a", 400},
		{"spoof-owner", "/source-navigation?case_id=case-one&entry_id=entry-new&owner_object_id=event-1", "Bearer analyst-a", 400},
		{"duplicate-id", "/source-navigation?case_id=case-one&entry_id=entry-new&entry_id=entry-old", "Bearer analyst-a", 400},
		{"wildcard", "/source-navigation?case_id=case-one&entry_id=%20*%20", "Bearer analyst-a", 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := sourceNavigationFixture(t)
			response := httpWorkspaceGet(h, tc.path, tc.auth, nil)
			if response.Code != tc.code {
				t.Fatalf("status=%d want=%d: %s", response.Code, tc.code, response.Body.String())
			}
			if response.Code != http.StatusOK &&
				(response.Body.String() != "{\"error\":\"request-unavailable\"}\n" ||
					strings.Contains(response.Body.String(), "entry-hidden") ||
					strings.Contains(response.Body.String(), "evidence-secret")) {
				t.Fatalf("denial leaks metadata: %s", response.Body.String())
			}
		})
	}
}

func TestSourceNavigationRechecksCaseAndDestinationRevocation(t *testing.T) {
	h, store := sourceNavigationFixture(t)
	// An authenticated Case access grant is not a capability or an ACL cache.
	store.revokeAfterList = true
	result := httpWorkspaceGet(h, "/source-navigation?case_id=case-one&entry_id=entry-new",
		"Bearer analyst-a", nil)
	if result.Code != http.StatusNotFound {
		t.Fatalf("revoked Case destination opened: %d", result.Code)
	}
	h, store = sourceNavigationFixture(t)
	owner := store.sources["src-new"]
	owner.owner = "analyst-b" // typed destination access was revoked
	store.sources["src-new"] = owner
	result = httpWorkspaceGet(h, "/source-navigation?case_id=case-one&entry_id=entry-new",
		"Bearer analyst-a", nil)
	if result.Code != http.StatusNotFound || strings.Contains(result.Body.String(), "event-1") {
		t.Fatalf("revoked source leaked destination: %d %s", result.Code, result.Body.String())
	}
	h, store = sourceNavigationFixture(t)
	entry := store.entries["entry-new"]
	entry.Version = 3 // owner changed while analyst navigates
	store.entries["entry-new"] = entry
	// The draft is based on the latest owner-resolved value, never a client version.
	result = httpWorkspaceGet(h, "/source-navigation?case_id=case-one&entry_id=entry-new",
		"Bearer analyst-a", nil)
	if result.Code != http.StatusOK {
		t.Fatalf("current authorized owner version unavailable: %d", result.Code)
	}
	var draft SourceNavigationResult
	if err := json.Unmarshal(result.Body.Bytes(), &draft); err != nil ||
		draft.Draft == nil || draft.Draft.Return.Version != 3 {
		t.Fatalf("stale version returned: %+v %v", draft, err)
	}
}

func TestSourceNavigationHTTPBoundaryRefusesMutation(t *testing.T) {
	h, store := sourceNavigationFixture(t)
	for _, method := range []string{http.MethodPost, http.MethodPatch, http.MethodDelete} {
		req := httptest.NewRequest(method, "/source-navigation?case_id=case-one&entry_id=entry-new", nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		if resp.Code != http.StatusMethodNotAllowed || resp.Header().Get("Allow") != "GET" {
			t.Fatalf("non-read method accepted: %s %d", method, resp.Code)
		}
	}
	req := httptest.NewRequest(http.MethodGet,
		"/source-navigation?case_id=case-one&entry_id=entry-new", strings.NewReader("payload"))
	req.Header.Set("Authorization", "Bearer analyst-a")
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest || store.catalogCalls != 0 {
		t.Fatalf("GET body reached owner: %d", resp.Code)
	}
}

func TestSourceNavigationRequiresAllOwnerProviders(t *testing.T) {
	store, _ := workspaceFixture()
	if _, err := NewCaseSourceNavigationHTTPHandler(nil, store, store, store, store); err == nil {
		t.Fatal("missing session provider accepted")
	}
	if _, err := NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{}, store, store, nil, store); err == nil {
		t.Fatal("missing entry owner accepted")
	}
	if _, err := NewCaseSourceNavigationHTTPHandler(timelineHTTPSession{}, store, store, store, nil); err == nil {
		t.Fatal("missing source owner accepted")
	}
}
