package investigationtimeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type timelineHTTPSession struct{}

func (timelineHTTPSession) ResolveCaseTimelineSession(r *http.Request) (context.Context, Scope, error) {
	if r.Header.Get("Authorization") != "Bearer analyst-a" {
		return nil, Scope{}, errors.New("unauthenticated")
	}
	return context.WithValue(r.Context(), integrationPrincipal{}, "analyst-a"),
		Scope{Tenant: "tenant-a", Environment: "prod",
			CasePermission: "deny", TimelinePermission: "unknown"}, nil
}

func httpWorkspaceFixture(t *testing.T) (http.Handler, *integrationStore) {
	t.Helper()
	store, _ := workspaceFixture()
	handler, err := NewCaseTimelineHTTPHandler(timelineHTTPSession{}, store, store, store, store)
	if err != nil {
		t.Fatal(err)
	}
	return handler, store
}

func httpWorkspaceGet(handler http.Handler, uri, auth string, headers map[string]string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, uri, nil)
	if auth != "" {
		request.Header.Set("Authorization", auth)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestCaseTimelineHTTPAuthorizedCaseToSources(t *testing.T) {
	handler, store := httpWorkspaceFixture(t)
	response := httpWorkspaceGet(handler, "/case-timeline?case_id=case-one",
		"Bearer analyst-a", map[string]string{
			"X-Tenant-Ref": "tenant-b", "X-Environment-Ref": "evil",
			"X-Timeline-Permission": "allow",
		})
	if response.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "no-store" ||
		response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("sensitive response cache or MIME guard missing")
	}
	var result Result
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "filtered" || len(result.Entries) != 2 ||
		result.Entries[0].ID != "entry-old" || result.Entries[0].Kind != "inferred" ||
		result.Entries[1].ID != "entry-new" || result.Entries[1].Supersedes != "entry-old" {
		t.Fatalf("integration lost observed/inferred chronology or corrections: %+v", result)
	}
	for _, hidden := range []string{"entry-hidden", "src-hidden", "evidence-secret", "tenant-b"} {
		if strings.Contains(response.Body.String(), hidden) {
			t.Fatalf("cross-source or cross-tenant disclosure: %q", hidden)
		}
	}
	if store.caseChecks != 2 || store.catalogCalls != 1 ||
		store.entryReads != 3 || store.sourceReads != 3 {
		t.Fatalf("missing Case/owner/source check: Case=%d catalog=%d entries=%d sources=%d",
			store.caseChecks, store.catalogCalls, store.entryReads, store.sourceReads)
	}
}

func TestCaseTimelineHTTPDenyWithoutCatalogEnumeration(t *testing.T) {
	cases := []struct {
		name, uri, auth string
		want int
	}{
		{"anonymous", "/case-timeline?case_id=case-one", "", http.StatusUnauthorized},
		{"invalid-session", "/case-timeline?case_id=case-one", "Bearer invalid", http.StatusUnauthorized},
		{"different-tenant", "/case-timeline?case_id=case-two", "Bearer analyst-a", http.StatusNotFound},
		{"absent-case", "/case-timeline?case_id=no-such-case", "Bearer analyst-a", http.StatusNotFound},
		{"spoofed-tenant", "/case-timeline?case_id=case-one&tenant=tenant-b", "Bearer analyst-a", http.StatusBadRequest},
		{"spoofed-permission", "/case-timeline?case_id=case-one&permission=allow", "Bearer analyst-a", http.StatusBadRequest},
		{"duplicate-case", "/case-timeline?case_id=case-one&case_id=case-two", "Bearer analyst-a", http.StatusBadRequest},
		{"no-case", "/case-timeline", "Bearer analyst-a", http.StatusBadRequest},
		{"wildcard", "/case-timeline?case_id=%2A", "Bearer analyst-a", http.StatusBadRequest},
		{"control-char", "/case-timeline?case_id=bad%0Acase", "Bearer analyst-a", http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler, store := httpWorkspaceFixture(t)
			response := httpWorkspaceGet(handler, tc.uri, tc.auth, nil)
			if response.Code != tc.want || store.catalogCalls != 0 {
				t.Fatalf("status=%d catalog calls=%d: %s", response.Code, store.catalogCalls, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "case-one") ||
				strings.Contains(response.Body.String(), "tenant-a") {
				t.Fatal("error disclosed protected identity")
			}
		})
	}
}

func TestCaseTimelineHTTPUnavailableOwnerAndRevokedCase(t *testing.T) {
	for _, condition := range []string{"catalog-unavailable", "case-revoked"} {
		t.Run(condition, func(t *testing.T) {
			handler, store := httpWorkspaceFixture(t)
			if condition == "catalog-unavailable" {
				store.failCatalog = true
			} else {
				store.revokeAfterList = true
			}
			response := httpWorkspaceGet(handler, "/case-timeline?case_id=case-one", "Bearer analyst-a", nil)
			if response.Code != http.StatusNotFound || response.Body.String() != "{\"error\":\"request-unavailable\"}\n" {
				t.Fatalf("provider failure not masked: %d %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestCaseTimelineHTTPRejectsMutationOrBody(t *testing.T) {
	handler, store := httpWorkspaceFixture(t)
	req := httptest.NewRequest(http.MethodPost, "/case-timeline?case_id=case-one", strings.NewReader("payload"))
	req.Header.Set("Authorization", "Bearer analyst-a")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != "GET" {
		t.Fatalf("POST was not rejected: %d", response.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/case-timeline?case_id=case-one", strings.NewReader("payload"))
	req.Header.Set("Authorization", "Bearer analyst-a")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest || store.catalogCalls != 0 {
		t.Fatalf("GET body accepted: %d catalog=%d", response.Code, store.catalogCalls)
	}
}

func TestCaseTimelineHTTPRejectsMissingServices(t *testing.T) {
	store, _ := workspaceFixture()
	if _, err := NewCaseTimelineHTTPHandler(nil, store, store, store, store); err == nil {
		t.Fatal("accepted missing authentication service")
	}
	if _, err := NewCaseTimelineHTTPHandler(timelineHTTPSession{}, nil, store, store, store); err == nil {
		t.Fatal("accepted missing ACL service")
	}
}

func TestCaseTimelineHTTPRejectsMalformedEscape(t *testing.T) {
	handler, store := httpWorkspaceFixture(t)
	req := httptest.NewRequest(http.MethodGet, "/case-timeline?case_id=case-one", nil)
	req.URL.RawQuery = "case_id=%GG"
	req.Header.Set("Authorization", "Bearer analyst-a")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest || store.catalogCalls != 0 {
		t.Fatalf("malformed query accepted: %d", response.Code)
	}
}
