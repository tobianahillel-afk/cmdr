package investigationtimeline

import (
	"bytes"
	"encoding/json"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestCaseHTMLAuthorizedTimelineAndServerPageNavigation(t *testing.T) {
	h, store := httpWorkspaceFixture(t)
	first := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&limit=1&offset=0",
		"Bearer analyst-a", map[string]string{"Accept": "text/html,application/xhtml+xml"})
	if first.Code != http.StatusOK ||
		first.Header().Get("Content-Type") != "text/html; charset=utf-8" ||
		first.Header().Get("Cache-Control") != "no-store" ||
		!strings.Contains(first.Header().Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatalf("browser is not protected: status=%d headers=%v", first.Code, first.Header())
	}
	body := first.Body.String()
	if !strings.Contains(body, "entry-old") ||
		!strings.Contains(body, "Remplacée (une version ultérieure autorisée est visible)") ||
		strings.Contains(body, "<script") ||
		strings.Contains(body, "entry-hidden") ||
		strings.Contains(body, "src-hidden") ||
		strings.Contains(body, "entry-new</th>") {
		t.Fatalf("source-authorized older page invalid: %s", body)
	}
	link := regexp.MustCompile(`rel="next" href="([^"]+)"`).FindStringSubmatch(body)
	if len(link) != 2 {
		t.Fatalf("page link unavailable: %s", body)
	}
	nextURL := html.UnescapeString(link[1])
	if nextURL != "/case-timeline?case_id=case-one&limit=1&offset=1" {
		t.Fatalf("unexpected navigation URL %q", nextURL)
	}
	next := httpWorkspaceGet(h, nextURL, "Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if next.Code != http.StatusOK ||
		!strings.Contains(next.Body.String(), "entry-new") ||
		!strings.Contains(next.Body.String(), "entry-old") ||
		strings.Contains(next.Body.String(), "entry-hidden") ||
		!strings.Contains(next.Body.String(), "Non déterminé (aucun successeur autorisé visible dans cette lecture)") {
		t.Fatalf("newer authorized HTML page incorrect: %d %s", next.Code, next.Body.String())
	}
	// Reading only the older page must still incorporate server-authorized
	// off-page versions; pagination itself cannot change supersession.
	if store.caseChecks < 2 {
		t.Fatalf("Case owner not rechecked on page navigation: %d", store.caseChecks)
	}
}

func TestCaseHTMLFailClosedAndJSONBackwardCompatibility(t *testing.T) {
	for _, tc := range []struct {
		accept string
		isHTML bool
	}{
		{"", false},
		{"application/json", false},
		{"application/json,text/html", false},
		{"text/html;q=0", false},
		{"text/html", true},
		{"text/html,application/xhtml+xml", true},
		{strings.Repeat("x", 1025), false},
	} {
		t.Run(tc.accept[:min(16, len(tc.accept))], func(t *testing.T) {
			h, _ := httpWorkspaceFixture(t)
			response := httpWorkspaceGet(h, "/case-timeline?case_id=case-one",
				"Bearer analyst-a", map[string]string{"Accept": tc.accept})
			if response.Code != http.StatusOK {
				t.Fatalf("status %d: %s", response.Code, response.Body.String())
			}
			if tc.isHTML {
				if !strings.HasPrefix(response.Header().Get("Content-Type"), "text/html") ||
					!strings.Contains(response.Body.String(), "<table>") {
					t.Fatal("browser did not receive HTML")
				}
				return
			}
			if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
				t.Fatal("old JSON response content-type changed")
			}
			var old Result
			if err := json.Unmarshal(response.Body.Bytes(), &old); err != nil ||
				len(old.Entries) != 2 || old.Entries[0].ID != "entry-old" {
				t.Fatalf("original JSON schema changed: %+v, %v", old, err)
			}
		})
	}
	for _, tc := range []struct{ auth, caseID string }{
		{"", "case-one"},
		{"Bearer invalid", "case-one"},
		{"Bearer analyst-a", "case-two"},
		{"Bearer analyst-a", "missing-case"},
	} {
		h, _ := httpWorkspaceFixture(t)
		denied := httpWorkspaceGet(h, "/case-timeline?case_id="+tc.caseID,
			tc.auth, map[string]string{"Accept": "text/html"})
		if denied.Code == http.StatusOK ||
			strings.Contains(denied.Body.String(), "entry-old") ||
			strings.Contains(denied.Body.String(), "entry-new") ||
			strings.Contains(denied.Body.String(), "case-one") ||
			strings.Contains(denied.Body.String(), "<!doctype html>") {
			t.Fatalf("HTML representation bypassed owner: %d %s", denied.Code, denied.Body.String())
		}
	}
}

func TestCaseHTMLTemplateEscapesAllSourceControlledFields(t *testing.T) {
	invalid := "<img src=x onerror=alert(1)>"
	source := Result{Status: "partial", Entries: []ProjectedEntry{{
		ID: invalid, Kind: "observed", SourceTime: "2026-10-21T05:00:00Z",
		Timezone: "+00:00", TimeUncertain: true,
		SourceRef: invalid, Provenance: invalid, Version: 1,
		Supersedes: invalid, Superseded: true,
	}}}
	var w bytes.Buffer
	err := renderTimelineHTML(&w, url.Values{"case_id": {invalid}}, invalid,
		timelineReadSelection{}, source)
	if err != nil {
		t.Fatal(err)
	}
	out := w.String()
	if strings.Contains(out, invalid) ||
		strings.Contains(out, "<script") ||
		strings.Contains(out, "<img src=") ||
		!strings.Contains(out, "&lt;img") ||
		!strings.Contains(out, "Incertitude") {
		t.Fatalf("template interpolated source-controlled HTML: %s", out)
	}
	// Href is fixed-route + Encode, so even a source-shaped query cannot
	// manufacture a second-origin destination or extra hidden authority.
	q := timelineHTMLPageURL(url.Values{
		"case_id": {"case-one"},
		"q": {"a&permission=allow"},
		"limit": {"50"},
		"offset": {"999"},
	}, 50, 50)
	if strings.Contains(q, "&permission=allow") ||
		!strings.HasPrefix(q, "/case-timeline?") ||
		!strings.Contains(q, "q=a%26permission%3Dallow") {
		t.Fatalf("link encoded uncontrolled query incorrectly: %s", q)
	}
}

func TestCaseHTMLRespectsExistingAuthorizedKindFilter(t *testing.T) {
	h, _ := httpWorkspaceFixture(t)
	page := httpWorkspaceGet(h,
		"/case-timeline?case_id=case-one&kind=inferred&limit=1",
		"Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if page.Code != http.StatusOK ||
		!strings.Contains(page.Body.String(), "entry-old") ||
		!strings.Contains(page.Body.String(), "Remplacée") ||
		strings.Contains(page.Body.String(), "entry-new</th>") {
		t.Fatalf("authorized filter lost supersession: %d %s", page.Code, page.Body.String())
	}
}
