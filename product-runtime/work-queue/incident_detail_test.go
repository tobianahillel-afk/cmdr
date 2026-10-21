package workqueue

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func detailGET(h http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestIncidentDetailDeepLinkFindsOnlyFreshOwnerAuthorizedIncident(t *testing.T) {
	owner, sessions := deepReturnOwner(), trustedHost()
	detail := NewIncidentDetailHTTPHandler(sessions, owner)
	ok := detailGET(detail, "/command/incident-detail?incident=inc-100")
	body := ok.Body.String()
	if ok.Code != http.StatusOK || !strings.Contains(body, "Incident Detail") ||
		!strings.Contains(body, "inc-100") || !strings.Contains(body, "in-progress") ||
		!strings.Contains(body, "Statut du dossier : partiel") ||
		!strings.Contains(body, "restore=1") ||
		strings.Contains(body, "inc-hidden") || strings.Contains(body, "foreign-secret") {
		t.Fatalf("authorized detail incorrect: status=%d body=%s", ok.Code, body)
	}
	if sessions.calls != 1 || len(owner.calls) != 124 {
		t.Fatalf("detail did not reauthorize all in-scope owner objects: %d %d", sessions.calls, len(owner.calls))
	}
	for k, v := range map[string]string{
		"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff",
		"Referrer-Policy": "no-referrer",
	} {
		if ok.Header().Get(k) != v {
			t.Fatalf("missing safe response header %s", k)
		}
	}
	if !strings.Contains(ok.Header().Get("Content-Security-Policy"), "form-action 'none'") {
		t.Fatal("detail includes unwanted executable form permissions")
	}
	owner.allowed["Incident:inc-100"] = false
	denied := detailGET(detail, "/command/incident-detail?incident=inc-100")
	if denied.Code != http.StatusNotFound || strings.Contains(denied.Body.String(), "inc-100") {
		t.Fatalf("revoked Incident revealed: %d %s", denied.Code, denied.Body.String())
	}
}

func TestIncidentDetailStrictInputNeverContactedOwner(t *testing.T) {
	for _, target := range []string{
		"/command/incident-detail",
		"/command/incident-detail?incident=",
		"/command/incident-detail?incident=inc-100&incident=inc-101",
		"/command/incident-detail?incident=inc-100&tenant=t2",
		"/command/incident-detail?incident=Task%3Ainc-100",
		"/command/incident-detail?incident=%00",
		"/command/incident-detail?incident=%ZZ",
		"/command/incident-detail?incident=" + strings.Repeat("x", 129),
	} {
		owner, sessions := deepReturnOwner(), trustedHost()
		rec := detailGET(NewIncidentDetailHTTPHandler(sessions, owner), target)
		if rec.Code != http.StatusBadRequest || sessions.calls != 0 ||
			len(owner.calls) != 0 || strings.Contains(rec.Body.String(), "inc-100") {
			t.Fatalf("invalid detail target reached authorization: %q status=%d session=%d owner=%d", target, rec.Code, sessions.calls, len(owner.calls))
		}
	}
}

func TestIncidentDetailRejectsNonHTMLMethodAndUnavailableSession(t *testing.T) {
	owner, sessions := deepReturnOwner(), trustedHost()
	handler := NewIncidentDetailHTTPHandler(sessions, owner)
	nonHTML := testQueueRequest(handler, http.MethodGet, "/command/incident-detail?incident=inc-100")
	if nonHTML.Code != http.StatusNotAcceptable || sessions.calls != 0 {
		t.Fatalf("non-browser unexpectedly returned incident: %d", nonHTML.Code)
	}
	post := httptest.NewRequest(http.MethodPost, "/command/incident-detail?incident=inc-100", nil)
	post.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, post)
	if rec.Code != http.StatusMethodNotAllowed || sessions.calls != 0 {
		t.Fatalf("mutation allowed: %d", rec.Code)
	}
	for _, ownerSession := range []*fixedHostSession{
		{err: errors.New("private session error")},
		{scope: SessionScope{Principal: "analyst", Tenant: "*", Environment: "e1"}},
	} {
		response := detailGET(NewIncidentDetailHTTPHandler(ownerSession, owner), "/command/incident-detail?incident=inc-100")
		if response.Code != http.StatusUnauthorized ||
			strings.Contains(response.Body.String(), "private") ||
			strings.Contains(response.Body.String(), "inc-100") {
			t.Fatalf("session failure leaked Incident: %d %s", response.Code, response.Body.String())
		}
	}
}

func TestIncidentDetailOwnerFailureAndOptionalSourcesStayFailClosed(t *testing.T) {
	failing := &fakeOwner{errList: true}
	if r := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), failing),
		"/command/incident-detail?incident=inc-100"); r.Code != http.StatusServiceUnavailable {
		t.Fatalf("provider error accepted: %d %s", r.Code, r.Body.String())
	}
	owner := deepReturnOwner()
	h := NewIncidentDetailHTTPHandler(trustedHost(), owner)
	for _, target := range []string{"inc-hidden", "foreign-secret", "unknown-id"} {
		got := detailGET(h, "/command/incident-detail?incident="+target)
		if got.Code != http.StatusNotFound || strings.Contains(got.Body.String(), target) {
			t.Fatalf("denied/nonexistent Incident revealed: %d %s", got.Code, got.Body.String())
		}
	}
	if got := detailGET(NewIncidentDetailHTTPHandler(nil, owner),
		"/command/incident-detail?incident=inc-100"); got.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing resolver accepted: %d", got.Code)
	}
	if got := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), nil),
		"/command/incident-detail?incident=inc-100"); got.Code != http.StatusServiceUnavailable {
		t.Fatalf("missing trusted owner accepted: %d", got.Code)
	}
	_ = context.Background()
}
