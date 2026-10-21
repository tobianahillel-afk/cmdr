package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A real localhost-style HTTP server drives both first-party handlers under
// one ephemeral Basic credential: the Incident Detail is a separate page,
// not the legacy Work Queue selected-card disguised as an accepted screen.
func TestLocalBrowserIncidentDetailUsesFreshCommandAndIndependentCaseOwners(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(app)
	defer s.Close()
	client := s.Client()

	first := demoRequest(t, client, s.URL,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"analyst", secrets.analyst, "text/html")
	home := demoBody(t, first)
	if first.StatusCode != http.StatusOK ||
		!strings.Contains(home, "/command/incident-detail?incident=incident-a") {
		t.Fatalf("authorized Work Queue did not link to actual Incident Detail: %d %s", first.StatusCode, home)
	}

	detail := demoRequest(t, client, s.URL,
		"/command/incident-detail?incident=incident-a",
		"analyst", secrets.analyst, "text/html")
	page := demoBody(t, detail)
	if detail.StatusCode != http.StatusOK || !strings.Contains(page, "Incident Detail") ||
		!strings.Contains(page, "incident-a") || !strings.Contains(page, "case_id=case-a") ||
		!strings.Contains(page, "restore=1") || strings.Contains(page, secrets.analyst) ||
		strings.Contains(page, secrets.commander) ||
		!strings.Contains(detail.Header.Get("Content-Security-Policy"), "form-action 'none'") ||
		!strings.Contains(page, "Suspicious authentication sequence (synthetic)") ||
		!strings.Contains(page, "high") || !strings.Contains(page, "medium") ||
		!strings.Contains(page, "Identity triage required") ||
		!strings.Contains(page, "synthetic-command-incident-owner") ||
		!strings.Contains(page, "Alert alert-a") ||
		!strings.Contains(page, "synthetic-command-alert-owner") ||
		!strings.Contains(page, "triaged") {
		t.Fatalf("real selected Incident Detail refused/unsafe: %d %s", detail.StatusCode, page)
	}
	back := demoRequest(t, client, s.URL,
		"/command/work-queue?restore=1&selected=Incident%3Aincident-a&view=incidents",
		"analyst", secrets.analyst, "text/html")
	if back.StatusCode != http.StatusOK || !strings.Contains(demoBody(t, back), "incident-a") {
		t.Fatal("read-only Incident Detail cannot reauthorize its Work Queue return")
	}
	casePage := demoRequest(t, client, s.URL,
		"/case-timeline?case_id=case-a&origin_incident=incident-a",
		"analyst", secrets.analyst, "text/html")
	if casePage.StatusCode != http.StatusOK || !strings.Contains(demoBody(t, casePage), "event-v1") {
		t.Fatal("Incident Detail link lost the independent Case owner")
	}
	commandOnly := demoRequest(t, client, s.URL,
		"/command/incident-detail?incident=incident-a",
		"commander", secrets.commander, "text/html")
	body := demoBody(t, commandOnly)
	if commandOnly.StatusCode != http.StatusOK || strings.Contains(body, "case_id=case-a") ||
		strings.Contains(body, "event-v1") || strings.Contains(body, "source-v1") ||
		strings.Contains(body, "Suspicious authentication sequence") ||
		strings.Contains(body, "alert-a") || strings.Contains(body, "synthetic-command-alert-owner") ||
		!strings.Contains(body, "Informations métier indisponibles") {
		t.Fatalf("Command-only reader was promoted to Case permission: %d %s", commandOnly.StatusCode, body)
	}
	for _, tc := range []struct{ path, role, secret string; want int }{
		{"/command/incident-detail?incident=incident-a", "", "", http.StatusUnauthorized},
		{"/command/incident-detail?incident=incident-a", "analyst", "bad-password", http.StatusUnauthorized},
		{"/command/incident-detail?incident=unknown", "analyst", secrets.analyst, http.StatusNotFound},
		{"/command/incident-detail?incident=incident-a&tenant=synthetic-tenant", "analyst", secrets.analyst, http.StatusBadRequest},
	} {
		res := demoRequest(t, client, s.URL, tc.path, tc.role, tc.secret, "text/html")
		body := demoBody(t, res)
		if res.StatusCode != tc.want || strings.Contains(body, "source-v1") {
			t.Fatalf("unauthorized detail permitted or exposed source: %s %d %s", tc.path, res.StatusCode, body)
		}
	}
}
