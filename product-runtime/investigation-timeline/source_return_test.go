package investigationtimeline

import (
	"encoding/json"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestSourcePreviewRestoresOnlyValidatedSameCaseTimelineSelection(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memorySourcePreviewFixture())
	if err != nil {
		t.Fatal(err)
	}
	returnURL := "/case-timeline?case_id=case-a&kind=observed&limit=1&offset=0&q=entry-a"
	link := "/demo-source-open?" + (url.Values{
		"case_id": {"case-a"}, "entry_id": {"entry-a"}, "return_to": {returnURL},
	}).Encode()
	w := sourcePreviewHTMLRequest(t, handlers.OpenSource, link, "Bearer aa", "text/html")
	if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" ||
		w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("authorized source rejected: %d %s", w.Code, w.Body.String())
	}
	matched := regexp.MustCompile(`href="(/case-timeline[?][^"]+)"`).FindStringSubmatch(w.Body.String())
	if len(matched) != 2 {
		t.Fatalf("no return link in escaped HTML: %s", w.Body.String())
	}
	got := html.UnescapeString(matched[1])
	if got != returnURL {
		t.Fatalf("filtered and paged return context altered: got=%q want=%q", got, returnURL)
	}
	returned := sourcePreviewHTMLRequest(t, handlers.Timeline, got, "Bearer aa", "text/html")
	if returned.Code != http.StatusOK {
		t.Fatalf("source return did not reauthorize Case selection: %d %s", returned.Code, returned.Body.String())
	}
	jsonResponse := sourcePreviewHTMLRequest(t, handlers.OpenSource, link, "Bearer aa", "application/json")
	var decoded map[string]any
	if jsonResponse.Code != http.StatusOK ||
		json.Unmarshal(jsonResponse.Body.Bytes(), &decoded) != nil ||
		decoded["entry_id"] != "entry-a" || decoded["status"] != "ready" {
		t.Fatalf("legacy source JSON changed: %d %s", jsonResponse.Code, jsonResponse.Body.String())
	}
	if _, leaked := decoded["return_to"]; leaked {
		t.Fatal("navigation hint is not part of the source JSON API")
	}
	if _, leaked := decoded["ReturnURL"]; leaked {
		t.Fatal("server navigation hint leaked in JSON")
	}
}

func TestSourceReturnContextNeverRedirectsOrGrantsAuthority(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memorySourcePreviewFixture())
	if err != nil {
		t.Fatal(err)
	}
	cases := []string{
		"https://evil.invalid/case-timeline?case_id=case-a&limit=1",
		"//evil.invalid/case-timeline?case_id=case-a&limit=1",
		"/case-timeline/../case-timeline?case_id=case-a&limit=1",
		"/%63ase-timeline?case_id=case-a&limit=1",
		"/case-timeline?case_id=case-b&limit=1&offset=0",
		"/case-timeline?case_id=case-a&limit=1&tenant=foreign",
		"/case-timeline?case_id=case-a&case_id=case-a&limit=1",
		"/case-timeline?case_id=case-a&limit=0",
		"/case-timeline?case_id=case-a&offset=1",
		"/case-timeline?case_id=case-a&limit=1#exfiltrate",
		"/case-timeline?case_id=case-a&limit=1&from=not-a-time",
		strings.Repeat("/case-timeline?case_id=case-a&limit=1", 100),
	}
	for _, back := range cases {
		t.Run(back[:min(len(back), 60)], func(t *testing.T) {
			if safe, ok := validatedSourceReturnURL(back, "case-a"); ok {
				t.Fatalf("unsafe return accepted: %q", safe)
			}
			link := "/demo-source-open?" + (url.Values{
				"case_id": {"case-a"}, "entry_id": {"entry-a"}, "return_to": {back},
			}).Encode()
			for _, accept := range []string{"text/html", "application/json"} {
				result := sourcePreviewHTMLRequest(t, handlers.OpenSource, link, "Bearer aa", accept)
				if result.Code != http.StatusBadRequest ||
					strings.Contains(result.Body.String(), "<html") ||
					strings.Contains(result.Body.String(), "object-a") ||
					strings.Contains(result.Body.String(), "source-a") {
					t.Fatalf("unsafe return reached source: accept=%q status=%d body=%s",
						accept, result.Code, result.Body.String())
				}
			}
		})
	}
}

func TestSourceReturnCannotPromoteDeniedSourceOrForeignCase(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memorySourcePreviewFixture())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ bearer, caseID, entryID string }{
		{"", "case-a", "entry-a"},
		{"Bearer aa", "case-b", "entry-b"},
		{"Bearer aa", "case-a", "hidden-a"},
	} {
		link := "/demo-source-open?" + (url.Values{
			"case_id": {tc.caseID}, "entry_id": {tc.entryID},
			"return_to": {"/case-timeline?case_id=" + tc.caseID + "&limit=1&offset=0"},
		}).Encode()
		for _, accept := range []string{"text/html", "application/json"} {
			got := sourcePreviewHTMLRequest(t, handlers.OpenSource, link, tc.bearer, accept)
			if got.Code == http.StatusOK || strings.Contains(got.Body.String(), "<html") ||
				strings.Contains(got.Body.String(), "DO NOT DISPLAY") {
				t.Fatalf("unsafe return promoted denied source: %d %s", got.Code, got.Body.String())
			}
		}
	}
}
