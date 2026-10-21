package main

import (
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

// Real first-party Go localhost HTTP round trip. All Command, Investigate,
// Timeline, typed-source and destination Incident owner checks are exercised
// separately; the navigation hint itself never carries an authorization.
func TestLocalBrowserIncidentDetailCaseSourceReturnToDetail(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(app)
	defer server.Close()
	client := server.Client()
	get := func(path, role, secret string) (int, string) {
		t.Helper()
		resp := demoRequest(t, client, server.URL, path, role, secret, "text/html")
		return resp.StatusCode, demoBody(t, resp)
	}
	link := func(body, path string) string {
		t.Helper()
		pattern := "href=\"(" + regexp.QuoteMeta(path) + "[?][^\"]+)\""
		found := regexp.MustCompile(pattern).FindStringSubmatch(body)
		if len(found) != 2 {
			t.Fatalf("expected %s link: %s", path, body)
		}
		return html.UnescapeString(found[1])
	}
	status, detail := get("/command/incident-detail?incident=incident-a", "analyst", secrets.analyst)
	if status != http.StatusOK || !strings.Contains(detail, "Result result-a") {
		t.Fatalf("original Incident Detail unavailable: %d %s", status, detail)
	}
	caseURL := link(detail, "/case-timeline")
	parsed, err := url.Parse(caseURL)
	if err != nil || parsed.Query().Get("origin_incident") != "incident-a" ||
		parsed.Query().Get("origin_view") != "detail" {
		t.Fatalf("Incident Detail did not preserve fixed return target: %s %v", caseURL, err)
	}
	status, casePage := get(caseURL, "analyst", secrets.analyst)
	if status != http.StatusOK || !strings.Contains(casePage, "Retour à Incident Detail") {
		t.Fatalf("Case denied or wrong return destination: %d %s", status, casePage)
	}
	sourceURL := link(casePage, "/demo-source-open")
	status, sourcePage := get(sourceURL, "analyst", secrets.analyst)
	if status != http.StatusOK || !strings.Contains(sourcePage, "Prévisualisation autorisée") {
		t.Fatalf("typed source denied: %d %s", status, sourcePage)
	}
	backToCase := link(sourcePage, "/case-timeline")
	restored, err := url.Parse(backToCase)
	if err != nil || restored.Query().Get("origin_view") != "detail" ||
		restored.Query().Get("origin_incident") != "incident-a" {
		t.Fatalf("source return dropped Incident Detail origin: %s %v", backToCase, err)
	}
	status, returnedCase := get(backToCase, "analyst", secrets.analyst)
	if status != http.StatusOK {
		t.Fatalf("Case/source return failed authorization: %d %s", status, returnedCase)
	}
	backToIncident := link(returnedCase, "/command/incident-detail")
	if backToIncident != "/command/incident-detail?incident=incident-a" {
		t.Fatalf("return was not fixed to Incident Detail: %s", backToIncident)
	}
	status, incidentAgain := get(backToIncident, "analyst", secrets.analyst)
	if status != http.StatusOK || !strings.Contains(incidentAgain, "Result result-a") ||
		!strings.Contains(incidentAgain, "Incident Detail") {
		t.Fatalf("destination Command Incident was not independently reauthorized: %d %s", status, incidentAgain)
	}
	// The legacy Work Queue origin still returns to the Work Queue; existing
	// callers are not secretly redirected to another product screen.
	status, legacy := get("/case-timeline?case_id=case-a&origin_incident=incident-a", "analyst", secrets.analyst)
	if status != http.StatusOK || !strings.Contains(legacy, "/command/work-queue?") ||
		strings.Contains(legacy, "Retour à Incident Detail") {
		t.Fatalf("legacy Work Queue return changed: %d %s", status, legacy)
	}
	for _, path := range []string{
		"/case-timeline?case_id=case-a&origin_view=detail",
		"/case-timeline?case_id=case-a&origin_incident=incident-a&origin_view=work-queue",
		"/case-timeline?case_id=case-a&origin_incident=incident-a&origin_view=detail&origin_view=detail",
		"/case-timeline?case_id=case-a&origin_incident=incident-a&origin_view=https:%2F%2Fevil.test",
	} {
		status, body := get(path, "analyst", secrets.analyst)
		if status != http.StatusBadRequest || strings.Contains(body, "incident-a") ||
			strings.Contains(body, "event-v1") {
			t.Fatalf("untrusted return target accepted: %s status=%d body=%s", path, status, body)
		}
	}
	status, body := get(caseURL, "commander", secrets.commander)
	if status == http.StatusOK || strings.Contains(body, "Incident Detail") ||
		strings.Contains(body, "event-v1") {
		t.Fatalf("Command-only viewer gained Case access: %d %s", status, body)
	}
	status, body = get("/command/incident-detail?incident=not-authorized", "analyst", secrets.analyst)
	if status != http.StatusNotFound || strings.Contains(body, "incident-a") {
		t.Fatalf("unknown origin was promoted to Incident grant: %d %s", status, body)
	}
}
