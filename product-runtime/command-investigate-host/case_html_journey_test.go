package commandinvestigatehost

import (
	"html"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCommandIncidentLinkedCaseOpensOwnerAuthorizedHumanTimeline(t *testing.T) {
	host, commandOwner := combinedTestHostWithHistory(t, true)
	work := request(host,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"Bearer analyst", "text/html")
	link := "/case-timeline?case_id=" + url.QueryEscape("case-a")
	if work.Code != http.StatusOK ||
		!strings.Contains(html.UnescapeString(work.Body.String()), "href=\""+link+"\"") ||
		commandOwner.lookups != 1 {
		t.Fatalf("authorized work-queue Case link missing: HTTP %d calls=%d", work.Code, commandOwner.lookups)
	}
	first := request(host, link+"&limit=1&offset=0", "Bearer analyst", "text/html")
	if first.Code != http.StatusOK ||
		!strings.HasPrefix(first.Header().Get("Content-Type"), "text/html") ||
		!strings.Contains(first.Body.String(), "event-a") ||
		!strings.Contains(first.Body.String(), "source-a") ||
		!strings.Contains(first.Body.String(), "Remplacée (version ultérieure autorisée visible)") ||
		!strings.Contains(first.Body.String(), "Page suivante") ||
		strings.Contains(first.Body.String(), "source-b") ||
		strings.Contains(first.Body.String(), "event-b") {
		t.Fatalf("cross-product old-version HTML leaked or omitted authorized history: %d %s", first.Code, first.Body.String())
	}
	next := request(host, link+"&limit=1&offset=1", "Bearer analyst", "text/html")
	if next.Code != http.StatusOK ||
		!strings.Contains(next.Body.String(), "event-b") ||
		!strings.Contains(next.Body.String(), "source-b") ||
		!strings.Contains(next.Body.String(), "Corrige l’entrée autorisée event-a") ||
		strings.Contains(next.Body.String(), "Remplacée (version ultérieure autorisée visible)") {
		t.Fatalf("cross-product correction successor display incorrect: %d %s", next.Code, next.Body.String())
	}
	for _, token := range []string{"Bearer commander", "", "Bearer invalid"} {
		browser := request(host, link+"&limit=1&offset=0", token, "text/html")
		if browser.Code == http.StatusOK ||
			strings.Contains(browser.Body.String(), "<html") ||
			strings.Contains(browser.Body.String(), "event-a") ||
			strings.Contains(browser.Body.String(), "source-a") ||
			strings.Contains(browser.Body.String(), "Remplacée") {
			t.Fatalf("Command-only/guest rights promoted by HTML navigation: %d %s", browser.Code, browser.Body.String())
		}
	}
	// The same remote Case route remains a JSON-only API for explicit API
	// consumers; human-readable presentation must never mutate the authority.
	api := request(host, link+"&limit=1&offset=0", "Bearer analyst", "application/json")
	if api.Code != http.StatusOK ||
		!strings.HasPrefix(api.Header().Get("Content-Type"), "application/json") ||
		!strings.Contains(api.Body.String(), "\"superseded\":true") ||
		strings.Contains(api.Body.String(), "<html") {
		t.Fatalf("JSON contract changed by human navigation: %d %s", api.Code, api.Body.String())
	}
}
