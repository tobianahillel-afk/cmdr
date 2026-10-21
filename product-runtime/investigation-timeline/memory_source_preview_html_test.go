package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func sourcePreviewHTMLRequest(t *testing.T, h http.Handler, path, bearer, accept string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Accept", accept)
	if bearer != "" {
		req.Header.Set("Authorization", bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestSourcePreviewHTMLUsesSameOwnerBoundaryAndEscapesUntrustedPreview(t *testing.T) {
	data := memorySourcePreviewFixture()
	data.Sources[0].Preview = "<script>alert(1)</script> & owner-controlled content"
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, data)
	if err != nil {
		t.Fatal(err)
	}
	path := "/demo-source-open?case_id=case-a&entry_id=entry-a"
	page := sourcePreviewHTMLRequest(t, handlers.OpenSource, path, "Bearer aa",
		"text/html,application/xhtml+xml,application/json;q=0.9")
	if page.Code != http.StatusOK || page.Header().Get("Content-Type") != "text/html; charset=utf-8" ||
		page.Header().Get("Cache-Control") != "no-store" ||
		page.Header().Get("Referrer-Policy") != "no-referrer" ||
		page.Header().Get("Vary") != "Authorization, Accept" {
		t.Fatalf("unsafe HTML response: status=%d headers=%v body=%s",
			page.Code, page.Header(), page.Body.String())
	}
	body := page.Body.String()
	for _, token := range []string{
		"Source de la chronologie", "case-a", "entry-a", "telemetry-event",
		"object-a", "&lt;script&gt;alert(1)&lt;/script&gt;",
		"owner-controlled content", "Prévisualisation synthétique locale",
		"case_id=case-a", "bouton Retour du navigateur",
	} {
		if !strings.Contains(body, token) {
			t.Fatalf("HTML missing authorized field %q: %s", token, body)
		}
	}
	for _, denied := range []string{
		"<script>", "src-hidden", "secret-a", "DO NOT DISPLAY", "object-b",
		"tenant-b", "Authorization:", "Bearer aa",
	} {
		if strings.Contains(body, denied) {
			t.Fatalf("HTML leaked or interpreted untrusted source %q: %s", denied, body)
		}
	}
	csp := page.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'none'") ||
		!strings.Contains(csp, "frame-ancestors 'none'") ||
		!strings.Contains(csp, "form-action 'none'") ||
		strings.Contains(csp, "unsafe-inline") {
		t.Fatalf("HTML CSP unexpectedly permissive: %q", csp)
	}
	jsonView := sourcePreviewHTMLRequest(t, handlers.OpenSource, path, "Bearer aa", "application/json")
	var original struct {
		Status    string `json:"status"`
		CaseID    string `json:"case_id"`
		EntryID   string `json:"entry_id"`
		Kind      string `json:"kind"`
		ObjectID  string `json:"object_id"`
		Preview   string `json:"preview"`
		Synthetic bool   `json:"synthetic"`
	}
	if jsonView.Code != http.StatusOK ||
		jsonView.Header().Get("Content-Type") != "application/json; charset=utf-8" ||
		jsonView.Header().Get("Vary") != "Authorization" ||
		json.Unmarshal(jsonView.Body.Bytes(), &original) != nil ||
		original.Status != "ready" || original.CaseID != "case-a" ||
		original.EntryID != "entry-a" || original.ObjectID != "object-a" ||
		!original.Synthetic || original.Preview != data.Sources[0].Preview {
		t.Fatalf("source JSON compatibility regressed: status=%d receipt=%+v body=%s",
			jsonView.Code, original, jsonView.Body.String())
	}
}

func TestSourcePreviewHTMLNeverRevealsDeniedOrForeignOwner(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memorySourcePreviewFixture())
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, bearer, path string
		status             int
	}{
		{"anonymous", "", "/demo-source-open?case_id=case-a&entry_id=entry-a", http.StatusUnauthorized},
		{"foreign-tenant", "Bearer aa", "/demo-source-open?case_id=case-b&entry_id=entry-b", http.StatusNotFound},
		{"denied-source", "Bearer aa", "/demo-source-open?case_id=case-a&entry_id=hidden-a", http.StatusNotFound},
		{"foreign-entry", "Bearer aa", "/demo-source-open?case_id=case-b&entry_id=entry-a", http.StatusNotFound},
		{"injected-source-ref", "Bearer aa", "/demo-source-open?case_id=case-a&entry_id=entry-a&source_ref=src-hidden", http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := sourcePreviewHTMLRequest(t, handlers.OpenSource, tc.path, tc.bearer, "text/html")
			if w.Code != tc.status || strings.Contains(w.Body.String(), "<html") ||
				strings.Contains(w.Body.String(), "DO NOT DISPLAY") ||
				strings.Contains(w.Body.String(), "object-a") ||
				strings.Contains(w.Body.String(), "tenant-b") ||
				strings.Contains(w.Body.String(), "src-hidden") {
				t.Fatalf("denied HTML leaked source: status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestSourcePreviewNegotiationDoesNotPromoteDisabledHTML(t *testing.T) {
	cases := map[string]bool{
		"text/html,application/json":       true,
		"text/html;q=0.8,application/json": true,
		"text/html; q=0":                   false,
		"text/html; q=NaN":                 false,
		"text/html; q=2":                   false,
		"application/json, text/html":      false,
		"application/xhtml+xml":            false,
		"text/htmlmalformed":               false,
		"":                                 false,
	}
	for header, expected := range cases {
		if got := acceptsSourcePreviewHTML(header); got != expected {
			t.Errorf("Accept %q selected HTML=%v, expected %v", header, got, expected)
		}
	}
}

func sourceOwnerWithIndependentIdentityCopies() *memoryTimelineProvider {
	data := memorySourcePreviewFixture()
	s := &memoryTimelineProvider{
		authenticate: func(*http.Request) (string, error) { return "analyst-a", nil },
		principals:    make(map[string]MemoryTimelinePrincipal),
		cases:         make(map[string]MemoryTimelineCase),
		entries:       make(map[string]MemoryTimelineEntry),
		sources:       make(map[string]MemoryTimelineSource),
		sourceObjects: make(map[memorySourceObjectKey]MemoryTimelineSource),
	}
	for _, p := range data.Principals {
		s.principals[p.ID] = p
	}
	for _, item := range data.Cases {
		s.cases[item.ID] = item
	}
	for _, item := range data.Entries {
		s.entries[item.Value.ID] = item
	}
	for _, src := range data.Sources {
		s.sources[src.Identity.Reference] = src
		key := memorySourceObjectKey{caseID: src.Identity.CaseID,
			kind: src.Identity.Kind, objectID: src.Identity.ObjectID}
		s.sourceObjects[key] = src
	}
	return s
}

func TestSourceHTMLRechecksCanonicalOwnerIdentityAfterSourceNavigation(t *testing.T) {
	key := memorySourceObjectKey{caseID: "case-a", kind: "telemetry-event", objectID: "object-a"}
	for _, tc := range []struct {
		name string
		change func(*SourceOwnerIdentity)
	}{
		{"tenant-swapped", func(v *SourceOwnerIdentity) { v.Tenant = "tenant-b" }},
		{"environment-swapped", func(v *SourceOwnerIdentity) { v.Environment = "foreign-env" }},
		{"case-swapped", func(v *SourceOwnerIdentity) { v.CaseID = "case-b" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate an inconsistent owner snapshot, while its original
			// source-reference mapping still resolves and is readable.
			// The destination's second identity check must deny HTML.
			s := sourceOwnerWithIndependentIdentityCopies()
			forged := s.sourceObjects[key]
			tc.change(&forged.Identity)
			s.sourceObjects[key] = forged
			r := sourcePreviewHTMLRequest(t, s.sourcePreviewHandler(),
				"/demo-source-open?case_id=case-a&entry_id=entry-a",
				"", "text/html")
			if r.Code != http.StatusNotFound || strings.Contains(r.Body.String(), "<html") ||
				strings.Contains(r.Body.String(), "Authorized synthetic DNS record") ||
				strings.Contains(r.Body.String(), "object-a") {
				t.Fatalf("inconsistent source owner identity was disclosed: %d %s",
					r.Code, r.Body.String())
			}
		})
	}
}

func TestSourceHTMLRejectsRequestWithoutURLAndKeepsSourceSecret(t *testing.T) {
	s := sourceOwnerWithIndependentIdentityCopies()
	req := httptest.NewRequest(http.MethodGet, "/demo-source-open?case_id=case-a&entry_id=entry-a", nil)
	req.URL = nil
	req.Header.Set("Accept", "text/html")
	w := httptest.NewRecorder()
	s.sourcePreviewHandler().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest || strings.Contains(w.Body.String(), "object-a") ||
		strings.Contains(w.Body.String(), "<html") {
		t.Fatalf("missing-URL HTML source accepted or disclosed: %d %s", w.Code, w.Body.String())
	}
}
