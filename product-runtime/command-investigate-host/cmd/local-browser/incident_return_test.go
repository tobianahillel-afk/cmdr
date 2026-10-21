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

func localBrowserLink(t *testing.T, body, expression string) string {
	t.Helper()
	matches := regexp.MustCompile(expression).FindStringSubmatch(body)
	if len(matches) != 2 {
		t.Fatalf("expected authorized navigation link missing: %s", body)
	}
	return html.UnescapeString(matches[1])
}

// Real first-party HTTP handlers, independently authorized Command Incident,
// Investigate Case, typed Timeline source and source.read, and fresh return GET.
func TestLocalBrowserIncidentCaseSourceTimelineIncidentRoundTrip(t *testing.T) {
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

	start := demoRequest(t, client, server.URL,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"analyst", secrets.analyst, "text/html")
	command := demoBody(t, start)
	if start.StatusCode != http.StatusOK {
		t.Fatalf("Command owner denied source selection: %d", start.StatusCode)
	}
	caseURL := localBrowserLink(t, command, `href="(/case-timeline[?][^"]+)"`)
	caseParsed, err := url.Parse(caseURL)
	if err != nil || caseParsed.Path != "/case-timeline" ||
		caseParsed.Query().Get("origin_incident") != "incident-a" ||
		caseParsed.Query().Get("case_id") != "case-a" ||
		len(caseParsed.Query()) != 2 {
		t.Fatalf("Command produced unscoped or over-broad link: %q: %v", caseURL, err)
	}
	caseValues := caseParsed.Query()
	caseValues.Set("limit", "1")
	caseValues.Set("offset", "0")
	caseURL = "/case-timeline?" + caseValues.Encode()

	caseResponse := demoRequest(t, client, server.URL, caseURL,
		"analyst", secrets.analyst, "text/html")
	caseHTML := demoBody(t, caseResponse)
	if caseResponse.StatusCode != http.StatusOK || !strings.Contains(caseHTML, "event-v1") {
		t.Fatalf("Case independently refused: %d %s", caseResponse.StatusCode, caseHTML)
	}
	commandReturn := localBrowserLink(t, caseHTML,
		`href="(/command/work-queue[?][^"]+)" rel="noreferrer">Retour`)
	if commandReturn != "/command/work-queue?selected=Incident%3Aincident-a&view=incidents" {
		t.Fatalf("untrusted return origin altered destination: %q", commandReturn)
	}
	sourceURL := localBrowserLink(t, caseHTML, `href="(/demo-source-open[?][^"]+)"`)
	sourceResponse := demoRequest(t, client, server.URL, sourceURL,
		"analyst", secrets.analyst, "text/html")
	sourceHTML := demoBody(t, sourceResponse)
	if sourceResponse.StatusCode != http.StatusOK ||
		!strings.Contains(sourceHTML, "Synthetic source") {
		t.Fatalf("independent typed source access denied: %d %s", sourceResponse.StatusCode, sourceHTML)
	}
	timelineReturn := localBrowserLink(t, sourceHTML, `href="(/case-timeline[?][^"]+)"`)
	timelineURL, err := url.Parse(timelineReturn)
	if err != nil || timelineURL.Path != "/case-timeline" ||
		timelineURL.Query().Get("origin_incident") != "incident-a" ||
		timelineURL.Query().Get("case_id") != "case-a" ||
		timelineURL.Query().Get("limit") != "1" ||
		timelineURL.Query().Get("offset") != "0" {
		t.Fatalf("source lost validated Case+Incident history: %q %v", timelineReturn, err)
	}
	backTimeline := demoRequest(t, client, server.URL, timelineReturn,
		"analyst", secrets.analyst, "text/html")
	backCaseHTML := demoBody(t, backTimeline)
	if backTimeline.StatusCode != http.StatusOK ||
		!strings.Contains(backCaseHTML, "event-v1") ||
		!strings.Contains(backCaseHTML, "Retour à l") {
		t.Fatalf("source return did not reauthorize Case: %d %s", backTimeline.StatusCode, backCaseHTML)
	}
	returnToCommand := localBrowserLink(t, backCaseHTML,
		`href="(/command/work-queue[?][^"]+)" rel="noreferrer">Retour`)
	incidentResponse := demoRequest(t, client, server.URL, returnToCommand,
		"analyst", secrets.analyst, "text/html")
	incidentHTML := demoBody(t, incidentResponse)
	if incidentResponse.StatusCode != http.StatusOK ||
		!strings.Contains(incidentHTML, "Incident — incident-a") {
		t.Fatalf("Command destination did not freshly reauthorize the selected Incident: %d %s",
			incidentResponse.StatusCode, incidentHTML)
	}
	for _, secret := range []string{secrets.analyst, secrets.commander} {
		for _, page := range []string{caseHTML, sourceHTML, backCaseHTML, incidentHTML} {
			if strings.Contains(page, secret) {
				t.Fatal("synthetic credentials leaked through browser navigation")
			}
		}
	}
	// A Command-only principal can still read its original Incident but not
	// gain Case/source access by replaying the navigation-only URL.
	commanderReturn := demoRequest(t, client, server.URL, returnToCommand,
		"commander", secrets.commander, "text/html")
	if commanderReturn.StatusCode != http.StatusOK {
		t.Fatalf("Command did not independently accept commander: %d", commanderReturn.StatusCode)
	}
	_ = demoBody(t, commanderReturn)
	commanderCase := demoRequest(t, client, server.URL, caseURL,
		"commander", secrets.commander, "text/html")
	body := demoBody(t, commanderCase)
	if commanderCase.StatusCode == http.StatusOK ||
		strings.Contains(body, "event-v1") ||
		strings.Contains(body, "incident-a") {
		t.Fatalf("navigation hint promoted Command right to Case: %d %s", commanderCase.StatusCode, body)
	}
}
