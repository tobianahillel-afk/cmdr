package investigationtimeline

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTimelineSourceLinksRequireExplicitTrustedHostOptIn(t *testing.T) {
	data := memorySourcePreviewFixture()
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(
		func(_ *http.Request) (string, error) { return "analyst-a", nil }, data)
	if err != nil {
		t.Fatal(err)
	}
	render := func(handler http.Handler, target string, spoof bool) (int, string) {
		t.Helper()
		r := httptest.NewRequest(http.MethodGet, target, nil)
		r.Header.Set("Accept", "text/html")
		if spoof {
			r.Header.Set("X-CMDR-Enable-Source-Links", "true")
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code, w.Body.String()
	}
	url := "/case-timeline?case_id=case-a&limit=1"
	status, ordinary := render(handlers.Timeline, url, true)
	if status != http.StatusOK || strings.Contains(ordinary, "/demo-source-open") {
		t.Fatalf("untrusted header exposed opt-in source links: %d %s", status, ordinary)
	}
	wrapped := WithLocalSourceLinks(handlers.Timeline)
	status, authorized := render(wrapped, url, false)
	if status != http.StatusOK || !strings.Contains(authorized,
		"/demo-source-open?case_id=case-a&amp;entry_id=entry-a") ||
		strings.Contains(authorized, "src-a&amp;") ||
		strings.Contains(authorized, "hidden-a") {
		t.Fatalf("trusted host links not scoped to visible entries: %d %s", status, authorized)
	}
	jsonReq := httptest.NewRequest(http.MethodGet, "/case-timeline?case_id=case-a", nil)
	jsonReq.Header.Set("Accept", "application/json")
	out := httptest.NewRecorder()
	wrapped.ServeHTTP(out, jsonReq)
	if out.Code != http.StatusOK || strings.Contains(out.Body.String(), "/demo-source-open") {
		t.Fatalf("opt-in polluted legacy JSON: %d %s", out.Code, out.Body.String())
	}
	if WithLocalSourceLinks(nil) != nil {
		t.Fatal("nil owner handler was promoted to an authorized route")
	}
}

func TestTimelineSourceLinkEncoderNeverCreatesUntrustedDestination(t *testing.T) {
	url := localSourceOpenURL("case a&role=admin", "event<broken>&fake=1")
	if url != "/demo-source-open?case_id=case+a%26role%3Dadmin&entry_id=event%3Cbroken%3E%26fake%3D1" {
		t.Fatalf("source link injection was not encoded: %q", url)
	}
}
