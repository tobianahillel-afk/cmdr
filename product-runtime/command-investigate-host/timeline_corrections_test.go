package commandinvestigatehost

import (
	"encoding/json"
	"html"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// This exercises the real Command and Investigate first-party HTTP handlers
// under a single origin. Their distinct owners remain synthetic fixtures;
// no permission is inferred from the Command navigation pointer.
func TestCommandCaseNavigationRendersOnlyAuthorizedTimelineCorrectionState(t *testing.T) {
	host, commandOwner := combinedTestHostWithHistory(t, true)
	start := request(host,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"Bearer analyst", "text/html")
	link := "/case-timeline?case_id=" + url.QueryEscape("case-a")
	originLink := link + "&origin_incident=incident-a"
	if start.Code != http.StatusOK ||
		!strings.Contains(html.UnescapeString(start.Body.String()), "href=\""+originLink+"\"") ||
		commandOwner.lookups != 1 {
		t.Fatalf("source-authorized Command link unavailable: code=%d lookups=%d", start.Code, commandOwner.lookups)
	}
	type correctionEntry struct {
		ID         string `json:"id"`
		Version    int    `json:"version"`
		Supersedes string `json:"supersedes"`
		Superseded bool   `json:"superseded"`
		SourceRef  string `json:"source_ref"`
	}
	type correctionPage struct {
		Entries []correctionEntry `json:"entries"`
	}
	read := func(page string) (correctionPage, int, string) {
		t.Helper()
		out := request(host, link+"&limit=1&offset="+page, "Bearer analyst", "")
		var value correctionPage
		if out.Code == http.StatusOK {
			if err := json.Unmarshal(out.Body.Bytes(), &value); err != nil {
				t.Fatalf("unreadable Case response: %v", err)
			}
		}
		return value, out.Code, out.Body.String()
	}
	first, code, data := read("0")
	if code != http.StatusOK || len(first.Entries) != 1 ||
		first.Entries[0].ID != "event-a" || first.Entries[0].Version != 1 ||
		!first.Entries[0].Superseded || first.Entries[0].SourceRef != "source-a" {
		t.Fatalf("authorized older version not marked despite off-page successor: %d %s", code, data)
	}
	second, code, data := read("1")
	if code != http.StatusOK || len(second.Entries) != 1 ||
		second.Entries[0].ID != "event-b" || second.Entries[0].Version != 2 ||
		second.Entries[0].Supersedes != "event-a" || second.Entries[0].Superseded ||
		second.Entries[0].SourceRef != "source-b" {
		t.Fatalf("authorized newer source/version not preserved: %d %s", code, data)
	}
	// Reading Command does not grant Case or Timeline permission, even
	// after the browser follows the link to an older/corrected Event.
	commandOnly := request(host,
		"/command/work-queue?view=incidents&selected=Incident:incident-a",
		"Bearer commander", "text/html")
	if commandOnly.Code != http.StatusOK ||
		!strings.Contains(commandOnly.Body.String(), "case_id=case-a") {
		t.Fatal("independently authorized Command link missing")
	}
	for _, token := range []string{"Bearer commander", "", "Bearer invalid"} {
		for _, page := range []string{"0", "1"} {
			out := request(host, link+"&limit=1&offset="+page, token, "")
			if (out.Code != http.StatusNotFound && out.Code != http.StatusUnauthorized) ||
				strings.Contains(out.Body.String(), "superseded") ||
				strings.Contains(out.Body.String(), "event-a") ||
				strings.Contains(out.Body.String(), "event-b") ||
				strings.Contains(out.Body.String(), "source-a") ||
				strings.Contains(out.Body.String(), "source-b") {
				t.Fatalf("denied Case leaked correction via page %s: %d %s", page, out.Code, out.Body.String())
			}
		}
	}
}
