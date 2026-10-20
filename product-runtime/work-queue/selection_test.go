package workqueue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVisibleCommandItemSelectionPreservesAuthorizedViewFilterAndOrder(t *testing.T) {
	owner := fixture()
	sessions := trustedHost()
	route := "/command/work-queue?view=tasks&status=open&sort=oldest&selected=Task:task-b"
	resp := testQueueRequest(NewHTTPHandler(sessions, owner), http.MethodGet, route)
	var got Result
	if resp.Code != http.StatusOK || json.Unmarshal(resp.Body.Bytes(), &got) != nil ||
		got.Selected == nil || got.Selected.Kind != Task || got.Selected.ID != "task-b" ||
		got.View != Tasks || got.SortOrder != "oldest" || got.StatusFilter != "open" ||
		len(got.Rows) != 2 || got.Rows[0].ID != "task-no-sla" {
		t.Fatalf("selected row did not retain authorized queue context: %d %s", resp.Code, resp.Body.String())
	}
	if sessions.calls != 1 || len(owner.calls) != 4 {
		t.Fatalf("selection re-fetched owner or skipped per-object ACL: sessions=%d calls=%v", sessions.calls, owner.calls)
	}
	for _, call := range owner.calls {
		if call == "Incident:foreign" {
			t.Fatal("cross-tenant ID queried through item selection")
		}
	}
	// Default response remains compatible if the user has selected nothing.
	defaultResult := testQueueRequest(NewHTTPHandler(trustedHost(), fixture()),
		http.MethodGet, "/command/work-queue?view=tasks")
	if defaultResult.Code != 200 || strings.Contains(defaultResult.Body.String(), "\"selected\"") {
		t.Fatalf("unselected legacy JSON response changed: %d %s", defaultResult.Code, defaultResult.Body.String())
	}
}

func TestNativeSelectedCommandRowUsesSameReadOnlyRouteWithoutNewInspector(t *testing.T) {
	owner := fixture()
	sessions := trustedHost()
	request := httptest.NewRequest(http.MethodGet,
		"/command/work-queue?view=tasks&status=open&sort=oldest&selected=Task:task-b", nil)
	request.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	NewHTTPHandler(sessions, owner).ServeHTTP(rec, request)
	body := rec.Body.String()
	if rec.Code != http.StatusOK || sessions.calls != 1 ||
		!strings.Contains(body, "aria-label=\"Contexte sélectionné\"") ||
		!strings.Contains(body, "Task — task-b") ||
		!strings.Contains(body, "name=\"selected\" value=\"Task:task-b\"") ||
		!strings.Contains(body, "name=\"status\"") ||
		!strings.Contains(body, "value=\"open\"") ||
		!strings.Contains(body, "<option value=\"oldest\" selected>") ||
		strings.Count(body, "href=") != 6 ||
		!strings.Contains(body, "form-action 'self'") && !strings.Contains(rec.Header().Get("Content-Security-Policy"), "form-action 'self'") {
		t.Fatalf("single-route selected browser lost form/navigation or CSP: %d %s", rec.Code, body)
	}
	if !strings.Contains(body, "action=\"/command/work-queue\"") ||
		strings.Contains(body, "/command/work-queue/tasks") ||
		strings.Contains(body, "inc-hidden") ||
		strings.Contains(body, "foreign") {
		t.Fatalf("browser selection exposed forbidden objects or extra route: %s", body)
	}
	// Owner-provided ID and context data remain html/template escaped even
	// when the ID is reused as a form-submitted selection token.
	xss := fixture()
	xssID := "<img src=x onerror=alert(1)>"
	xss.items[0].ID = xssID
	xss.allowed["Incident:"+xssID] = true
	req := httptest.NewRequest(http.MethodGet,"/command/work-queue?view=incidents",nil)
	req.Header.Set("Accept","text/html")
	out := httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),xss).ServeHTTP(out,req)
	if out.Code != 200 || !strings.Contains(out.Body.String(), "&lt;img") ||
		strings.Contains(out.Body.String(), "<img") {
		t.Fatalf("untrusted ID escaped incorrectly in selection form: %d %s", out.Code, out.Body.String())
	}
}

func TestSelectedCommandRowFailsClosedForInvalidDeniedOrOffPageIDs(t *testing.T) {
	for _, query := range []string{
		"selected=", "selected=Incident", "selected=Case:inc-a",
		"selected=Task:", "selected=Task:*",
		"selected=Task:task-b&selected=Task:task-no-sla",
		"selected=Task:%FF", "selected=Task:%0Asecret",
		"selected=Incident:"+strings.Repeat("a",257),
		"selected=Task:task-b&tenant=t2",
	} {
		session := trustedHost()
		owner := fixture()
		got := testQueueRequest(NewHTTPHandler(session,owner),http.MethodGet,
			"/command/work-queue?"+query)
		if got.Code != http.StatusBadRequest || session.calls != 0 || len(owner.calls) != 0 {
			t.Fatalf("invalid selection reached trusted owner: %q status=%d", query,got.Code)
		}
	}
	for _, query := range []string{
		"view=all&selected=Incident:inc-hidden",
		"view=all&selected=Incident:foreign",
		"view=tasks&status=closed&selected=Task:task-b",
		"view=incidents&selected=Task:task-b",
		"view=all&selected=Incident:nonexistent",
	} {
		r := testQueueRequest(NewHTTPHandler(trustedHost(),fixture()),http.MethodGet,
			"/command/work-queue?"+query)
		if r.Code != http.StatusNotFound || strings.Contains(r.Body.String(),"inc-hidden") ||
			strings.Contains(r.Body.String(),"foreign") ||
			strings.Contains(r.Body.String(),"task-b") {
			t.Fatalf("unavailable item selection disclosed data: %q %d %s", query,r.Code,r.Body.String())
		}
	}
	owner := fixture()
	owner.items = nil
	owner.allowed = map[string]bool{}
	for i:=0;i<60;i++ {
		id := fmt.Sprintf("task-%03d",i)
		owner.items=append(owner.items,Item{
			Kind:Task,ID:id,Tenant:"t1",Environment:"e1",
			Status:"open",UpdatedAt:"2026-10-20T09:00:00Z",
		})
		owner.allowed["Task:"+id]=true
	}
	// Authorizing an object outside the 50-row visible page must NOT make
	// that object available via a forged selection token.
	page := testQueueRequest(NewHTTPHandler(trustedHost(),owner),http.MethodGet,
		"/command/work-queue?view=tasks&selected=Task:task-059")
	if page.Code != http.StatusNotFound {
		t.Fatalf("selected off-page object leaked beyond visible bound: %d %s",page.Code,page.Body.String())
	}
	revoked:=fixture()
	revoked.errRead=true
	denied:=testQueueRequest(NewHTTPHandler(trustedHost(),revoked),http.MethodGet,
		"/command/work-queue?selected=Incident:inc-a")
	if denied.Code != http.StatusServiceUnavailable || strings.Contains(denied.Body.String(),"inc-a") {
		t.Fatalf("revoked owner returned selected row: %d %s",denied.Code,denied.Body.String())
	}
}
