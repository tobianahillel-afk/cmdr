package workqueue

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func deepReturnOwner() *fakeOwner {
	o := &fakeOwner{allowed: map[string]bool{}}
	for i := 0; i < 123; i++ {
		id := fmt.Sprintf("inc-%03d", i)
		o.items = append(o.items, Item{
			Kind: Incident, ID: id, Tenant: "t1", Environment: "e1",
			Owner: "analyst", Status: "in-progress",
			UpdatedAt: "2026-10-20T09:00:00Z",
		})
		o.allowed["Incident:"+id] = true
	}
	o.items = append(o.items,
		Item{Kind: Incident, ID: "inc-hidden", Tenant: "t1", Environment: "e1",
			Status: "hidden", UpdatedAt: "2026-10-20T09:00:00Z"},
		Item{Kind: Incident, ID: "foreign-secret", Tenant: "t2", Environment: "e1",
			Status: "hidden", UpdatedAt: "2026-10-20T09:00:00Z"},
	)
	return o
}

func returnHTML(h http.Handler, query string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, "/command/work-queue?"+query, nil)
	r.Header.Set("Accept", "text/html")
	out := httptest.NewRecorder()
	h.ServeHTTP(out, r)
	return out
}

func TestDeepIncidentReturnReconstructsFreshAuthorizedThirdPage(t *testing.T) {
	o, session := deepReturnOwner(), trustedHost()
	h := NewHTTPHandler(session, o)
	q := "view=incidents&selected=Incident%3Ainc-100"
	// Ordinary page selection NEVER authorizes an off-page ID.
	old := returnHTML(h, q)
	if old.Code != http.StatusNotFound {
		t.Fatalf("old page-bound selection was widened: %d", old.Code)
	}
	p := returnHTML(h, q+"&restore=1")
	html := p.Body.String()
	if p.Code != http.StatusOK || !strings.Contains(html, "Incident — inc-100") ||
		!strings.Contains(html, "inc-101") || strings.Contains(html, "inc-000") ||
		!strings.Contains(html, "Page précédente") || !strings.Contains(html, "Première page") ||
		!strings.Contains(html, "cursor=") {
		t.Fatalf("Incident failed owner reauthorization or cursor reconstruction: %d %s", p.Code, html)
	}
	if strings.Contains(html, "foreign-secret") || strings.Contains(html, "inc-hidden") {
		t.Fatal("deep Incident return disclosed an unauthorized row")
	}
	if session.calls != 2 || len(o.calls) != 2*124 {
		t.Fatalf("resume skipped full fresh owner authorization: session=%d per-record=%d",
			session.calls, len(o.calls))
	}
	// Owner revocation must be checked again; a once-visible Incident is
	// never an authorization grant or a restored snapshot.
	o.allowed["Incident:inc-100"] = false
	gone := returnHTML(h, q+"&restore=1")
	if gone.Code != http.StatusNotFound ||
		strings.Contains(gone.Body.String(), "inc-100") ||
		strings.Contains(gone.Body.String(), "foreign-secret") {
		t.Fatalf("revoked Incident was restored: %d %s", gone.Code, gone.Body.String())
	}
}

func TestDeepIncidentReturnIsHTMLOnlyAndRejectsForgedSelection(t *testing.T) {
	o, session := deepReturnOwner(), trustedHost()
	h := NewHTTPHandler(session, o)
	for _, query := range []string{
		"view=incidents&selected=Incident%3Ainc-100&restore=0",
		"view=incidents&selected=Incident%3Ainc-100&restore=1&restore=1",
		"view=incidents&restore=1",
		"view=tasks&selected=Incident%3Ainc-100&restore=1",
		"view=incidents&selected=Task%3Atask-100&restore=1",
		"view=incidents&selected=Incident%3Ainc-100&restore=1&status=open",
		"view=incidents&selected=Incident%3Ainc-100&restore=1&cursor=deadbeef",
		"view=incidents&selected=Incident%3Ainc-100&restore=1&principal=admin",
		"view=incidents&selected=Incident%3Ainc-100&restore=1&tenant=t2",
	} {
		got := returnHTML(h, query)
		if got.Code != http.StatusBadRequest || strings.Contains(got.Body.String(), "inc-100") {
			t.Fatalf("forged return accepted/leaked: %q -> %d %s", query, got.Code, got.Body.String())
		}
	}
	r := testQueueRequest(h, http.MethodGet, "/command/work-queue?view=incidents&selected=Incident%3Ainc-100&restore=1")
	if r.Code != http.StatusBadRequest || strings.Contains(r.Body.String(), "inc-100") {
		t.Fatalf("JSON API silently gained deep selection: %d %s", r.Code, r.Body.String())
	}
	// All invalid return hints fail before a privileged session/owner read.
	if session.calls != 0 || len(o.calls) != 0 {
		t.Fatalf("untrusted return hints reached owners: session=%d owner=%d", session.calls, len(o.calls))
	}
}
