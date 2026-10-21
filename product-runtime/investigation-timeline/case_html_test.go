package investigationtimeline

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCaseHTMLAuthorizedCorrectionsAndSafePagination(t *testing.T) {
	h, store := httpWorkspaceFixture(t)
	first := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&limit=1&offset=0",
		"Bearer analyst-a", map[string]string{"Accept": "text/html,application/xhtml+xml;q=0.9"})
	if first.Code != http.StatusOK ||
		!strings.HasPrefix(first.Header().Get("Content-Type"), "text/html") ||
		first.Header().Get("Cache-Control") != "no-store" ||
		!strings.Contains(first.Header().Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatalf("authorized browser representation unavailable: %d %s", first.Code, first.Body.String())
	}
	text := first.Body.String()
	for _, want := range []string{
		"Chronologie Investigate", "entry-old", "inferred",
		"Remplacée (version ultérieure autorisée visible)",
		"case_id=case-one", "limit=1", "offset=1",
		"nouvelle autorisation requise",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing authorized HTML semantics %q: %s", want, text)
		}
	}
	for _, hidden := range []string{"entry-hidden", "src-hidden", "evidence-secret", "tenant-b", "entry-new"} {
		if strings.Contains(text, hidden) {
			t.Fatalf("hidden/off-page source data leaked into HTML: %q", hidden)
		}
	}
	second := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&limit=1&offset=1",
		"Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if second.Code != http.StatusOK ||
		!strings.Contains(second.Body.String(), "entry-new") ||
		!strings.Contains(second.Body.String(), "Corrige l’entrée autorisée entry-old") ||
		!strings.Contains(second.Body.String(), "Page précédente") ||
		strings.Contains(second.Body.String(), "Remplacée (version ultérieure autorisée visible)") {
		t.Fatalf("wrong authorized successor page: %d %s", second.Code, second.Body.String())
	}
	if store.catalogCalls != 2 || store.sourceReads != 6 {
		t.Fatalf("HTML bypassed owner chain: catalog=%d source=%d", store.catalogCalls, store.sourceReads)
	}
}

func TestCaseHTMLNeverLeaksOnDeniedOrUnknownRequest(t *testing.T) {
	for _, tc := range []struct{
		name, bearer, caseID string
		want int
	}{
		{"anonymous", "", "case-one", http.StatusUnauthorized},
		{"wrong-case", "Bearer analyst-a", "case-two", http.StatusNotFound},
		{"invalid-session", "Bearer invalid", "case-one", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := httpWorkspaceFixture(t)
			got := httpWorkspaceGet(h, "/case-timeline?case_id="+tc.caseID,
				tc.bearer, map[string]string{"Accept": "text/html"})
			if got.Code != tc.want || strings.Contains(got.Body.String(), "<html") ||
				strings.Contains(got.Body.String(), "entry-old") ||
				strings.Contains(got.Body.String(), "src-old") ||
				strings.Contains(got.Body.String(), "Remplacée") {
				t.Fatalf("denied browser view leaked Case source: %d %s", got.Code, got.Body.String())
			}
		})
	}
}

func TestCaseHTMLInjectionEscapedAndJSONCompatible(t *testing.T) {
	escaped, err := renderAuthorizedTimelineHTML("<img src=x onerror=alert(1)>", Result{
		Status: "allow",
		Entries: []ProjectedEntry{{
			ID: "<script>alert(1)</script>", Kind: "observed",
			SourceRef: "<svg/onload=alert(1)>", Provenance: "<b>untrusted</b>",
			Version: 1, Timezone: "+00:00", SourceTime: "2026-10-21T00:00:00Z",
		}},
	}, timelineReadSelection{}, url.Values{})
	if err != nil || strings.Contains(string(escaped), "<script>") ||
		strings.Contains(string(escaped), "<svg/onload=") ||
		strings.Contains(string(escaped), "<img src=") ||
		!strings.Contains(string(escaped), "&lt;script&gt;") {
		t.Fatalf("source metadata must be escaped: %v %s", err, escaped)
	}
	h, _ := httpWorkspaceFixture(t)
	jsonDefault := httpWorkspaceGet(h, "/case-timeline?case_id=case-one",
		"Bearer analyst-a", nil)
	jsonExplicit := httpWorkspaceGet(h, "/case-timeline?case_id=case-one",
		"Bearer analyst-a", map[string]string{"Accept": "application/json,text/html"})
	if jsonDefault.Code != http.StatusOK || jsonExplicit.Code != http.StatusOK ||
		jsonDefault.Body.String() != jsonExplicit.Body.String() ||
		jsonDefault.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("content negotiation modified original JSON consumers: %s / %s",
			jsonDefault.Body.String(), jsonExplicit.Body.String())
	}
	for _, accept := range []string{"text/html;q=0", "text/html;q=bogus", "application/json"} {
		response := httpWorkspaceGet(h, "/case-timeline?case_id=case-one",
			"Bearer analyst-a", map[string]string{"Accept": accept})
		if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("unsupported content preference served HTML: %s", accept)
		}
	}
}

func TestCaseHTMLPagerLinksPreserveOnlyValidatedFilters(t *testing.T) {
	h, _ := httpWorkspaceFixture(t)
	res := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&kind=inferred&limit=1&offset=0",
		"Bearer analyst-a", map[string]string{"Accept": "text/html"})
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "entry-old") ||
		strings.Contains(res.Body.String(), "entry-new") ||
		!strings.Contains(res.Body.String(), "Remplacée (version ultérieure autorisée visible)") {
		t.Fatalf("filter incorrectly erased authorized supersession: %d %s", res.Code, res.Body.String())
	}
	if strings.Contains(res.Body.String(), "Page suivante") {
		t.Fatal("a non-security filter must not create phantom next pages")
	}
}
