package commandinvestigatehost

import (
	"html"
	"net/http"
	"strings"
	"testing"
)

// This exercises BOTH production-shaped Go handlers across real HTTP routing,
// with independent *synthetic* Command and Investigate owner fixtures. A
// source-authorized link is never a transferable Case read permission.
func TestCommandLinkedCaseOpensActualOwnerAuthorizedHTMLTimeline(t *testing.T) {
	host, commandOwner := combinedTestHostWithHistory(t, true)
	origin := request(host,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"Bearer analyst", "text/html")
	link := "/case-timeline?case_id=case-a"
	if origin.Code != http.StatusOK ||
		!strings.Contains(html.UnescapeString(origin.Body.String()), "href=\""+link+"\"") ||
		commandOwner.lookups != 1 {
		t.Fatalf("Command did not expose an authorized linked Case: %d %s", origin.Code, origin.Body.String())
	}
	page := request(host, link+"&limit=1&offset=0", "Bearer analyst", "text/html")
	if page.Code != http.StatusOK || page.Header().Get("Content-Type") != "text/html; charset=utf-8" ||
		!strings.Contains(page.Body.String(), "<h1>Chronologie du Case</h1>") ||
		!strings.Contains(page.Body.String(), "event-a") ||
		!strings.Contains(page.Body.String(), "Remplacée (une version ultérieure autorisée est visible)") ||
		strings.Contains(page.Body.String(), "event-b</th>") ||
		!strings.Contains(page.Body.String(), "rel=\"next\"") {
		t.Fatalf("linked Case not readable as source-authorized HTML: %d %s", page.Code, page.Body.String())
	}
	next := request(host, link+"&limit=1&offset=1", "Bearer analyst", "text/html")
	if next.Code != http.StatusOK ||
		!strings.Contains(next.Body.String(), "event-b") ||
		!strings.Contains(next.Body.String(), "event-a") ||
		!strings.Contains(next.Body.String(), "Non déterminé") {
		t.Fatalf("newer authorized version not displayed without false current assertion: %d %s", next.Code, next.Body.String())
	}
	for _, token := range []string{"Bearer commander", "", "Bearer invalid"} {
		rejected := request(host, link, token, "text/html")
		if rejected.Code == http.StatusOK ||
			strings.Contains(rejected.Body.String(), "<!doctype html>") ||
			strings.Contains(rejected.Body.String(), "event-a") ||
			strings.Contains(rejected.Body.String(), "source-a") ||
			strings.Contains(rejected.Body.String(), "superseded") {
			t.Fatalf("Command token promoted into Investigate read: %d %s", rejected.Code, rejected.Body.String())
		}
	}
	// API clients are still given precisely the legacy JSON representation.
	api := request(host, link, "Bearer analyst", "application/json")
	if api.Code != http.StatusOK ||
		!strings.HasPrefix(api.Header().Get("Content-Type"), "application/json") ||
		!strings.Contains(api.Body.String(), "\"superseded\":true") {
		t.Fatalf("HTML path broke original JSON owner projection: %d %s", api.Code, api.Body.String())
	}
}
