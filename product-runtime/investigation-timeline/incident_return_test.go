package investigationtimeline

import (
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestCaseIncidentReturnContextHonorsHTMLAndSourceSecurity(t *testing.T) {
	h, _ := httpWorkspaceFixture(t)
	target := "/case-timeline?case_id=case-one&kind=inferred&limit=1&offset=0&origin_incident=incident-1"
	page := httpWorkspaceGet(h, target, "Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if page.Code != http.StatusOK {
		t.Fatalf("authorized Case page rejected: %d %s", page.Code, page.Body.String())
	}
	link := regexp.MustCompile(`href="(/command/work-queue[?][^"]+)"`).FindStringSubmatch(page.Body.String())
	if len(link) != 2 {
		t.Fatalf("missing Command navigation in authorized Case HTML: %s", page.Body.String())
	}
	parsed, err := url.Parse(html.UnescapeString(link[1]))
	if err != nil || parsed.Path != "/command/work-queue" ||
		parsed.Query().Get("selected") != "Incident:incident-1" ||
		parsed.Query().Get("view") != "incidents" || len(parsed.Query()) != 2 {
		t.Fatalf("Command destination is not a single fixed owner route: %v, %v", parsed, err)
	}
	// All Case pagination URLs must preserve inert origin, not the credentials.
	page2 := httpWorkspaceGet(h,
		"/case-timeline?case_id=case-one&limit=1&offset=0&origin_incident=incident-1",
		"Bearer analyst-a", map[string]string{"Accept": "text/html"})
	next := regexp.MustCompile(`rel="next" href="([^"]+)"`).FindStringSubmatch(page2.Body.String())
	if page2.Code != http.StatusOK || len(next) != 2 ||
		!strings.Contains(html.UnescapeString(next[1]), "origin_incident=incident-1") {
		t.Fatalf("next authorized Case page lost Command return hint: %d %s", page2.Code, page2.Body.String())
	}
	// Existing JSON API must never accept/reflect the HTML-only hint.
	json := httpWorkspaceGet(h, target, "Bearer analyst-a", map[string]string{"Accept": "application/json"})
	if json.Code != http.StatusBadRequest ||
		strings.Contains(json.Body.String(), "incident-1") {
		t.Fatalf("JSON API silently widened by navigation: %d %s", json.Code, json.Body.String())
	}
}

func TestIncidentReturnRejectsMalformedHintsBeforeCaseAccess(t *testing.T) {
	h, _ := httpWorkspaceFixture(t)
	bad := []string{
		"origin_incident=", "origin_incident=%0A", "origin_incident=Incident%0Aevil",
		"origin_incident=%20incident-a", "origin_incident=incident-a%20",
		"origin_incident=incident-a&origin_incident=incident-b",
		"origin_incident=%FF", "origin_incident=incident-a%00evil",
		"origin_incident=incident-a&tenant=foreign",
		"origin_incident=" + strings.Repeat("x", 129),
	}
	for _, suffix := range bad {
		t.Run(suffix[:min(len(suffix), 48)], func(t *testing.T) {
			page := httpWorkspaceGet(h,
				"/case-timeline?case_id=case-one&"+suffix,
				"Bearer analyst-a", map[string]string{"Accept": "text/html"})
			if page.Code != http.StatusBadRequest || strings.Contains(page.Body.String(), "entry-old") ||
				strings.Contains(page.Body.String(), "/command/work-queue?") {
				t.Fatalf("invalid return hint influenced authorized Case: %d %s", page.Code, page.Body.String())
			}
		})
	}
	denied := httpWorkspaceGet(h, "/case-timeline?case_id=case-two&origin_incident=incident-a",
		"Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if denied.Code == http.StatusOK || strings.Contains(denied.Body.String(), "incident-a") {
		t.Fatalf("foreign Case obtained Incident return context: %d %s", denied.Code, denied.Body.String())
	}
}

func TestSourceReturnAllowsOnlyCanonicalBoundedNestedCommandHint(t *testing.T) {
	valid := "/case-timeline?case_id=case-a&limit=1&offset=0&origin_incident=incident-a"
	got, ok := validatedSourceReturnURL(valid, "case-a")
	if !ok || got != "/case-timeline?case_id=case-a&limit=1&offset=0&origin_incident=incident-a" {
		t.Fatalf("nested authorization-free incident hint lost: %q %v", got, ok)
	}
	for _, bad := range []string{
		"/case-timeline?case_id=case-a&limit=1&origin_incident=",
		"/case-timeline?case_id=case-a&limit=1&origin_incident=incident-a&origin_incident=incident-b",
		"/case-timeline?case_id=case-a&limit=1&origin_incident=incident-a&permission=allow",
		"/case-timeline?case_id=case-a&limit=1&origin_incident=incident-a%0A",
		"/case-timeline?case_id=case-b&limit=1&origin_incident=incident-a",
		"https://evil.invalid/case-timeline?case_id=case-a&limit=1&origin_incident=incident-a",
	} {
		if got, ok := validatedSourceReturnURL(bad, "case-a"); ok {
			t.Fatalf("forged source history admitted: %q -> %q", bad, got)
		}
	}
}
