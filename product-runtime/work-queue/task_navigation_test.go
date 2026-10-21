package workqueue

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func deepTaskOwner() *fakeOwner {
	o := &fakeOwner{allowed: map[string]bool{}}
	for i := 0; i < 123; i++ {
		id := fmt.Sprintf("task-%03d", i)
		o.items = append(o.items, Item{
			Kind: Task, ID: id, Tenant: "t1", Environment: "e1",
			Owner: "analyst", Status: "in-progress",
			UpdatedAt: "2026-10-21T09:00:00Z",
		})
		o.allowed["Task:"+id] = true
	}
	o.items = append(o.items,
		Item{Kind: Task, ID: "hidden-task", Tenant: "t1", Environment: "e1",
			Owner: "private-source", Status: "blocked", UpdatedAt: "2026-10-21T09:00:00Z"},
		Item{Kind: Task, ID: "foreign-task", Tenant: "t2", Environment: "e1",
			Owner: "foreign-source", Status: "blocked", UpdatedAt: "2026-10-21T09:00:00Z"},
	)
	return o
}

func TestLinkedTaskNavigationRechecksOwnerAndRestoresTaskBeyondFirstPage(t *testing.T) {
	owner, session := deepTaskOwner(), trustedHost()
	h := NewHTTPHandler(session, owner)
	selection := "view=tasks&selected=Task%3Atask-100"
	ordinary := returnHTML(h, selection)
	if ordinary.Code != http.StatusNotFound {
		t.Fatalf("ordinary selected Task escaped page limit: %d", ordinary.Code)
	}
	authorized := returnHTML(h, selection+"&restore=1")
	text := authorized.Body.String()
	if authorized.Code != http.StatusOK ||
		!strings.Contains(text, "Task — task-100") ||
		!strings.Contains(text, "Page précédente") ||
		!strings.Contains(text, "cursor=") ||
		strings.Contains(text, "task-000") ||
		strings.Contains(text, "hidden-task") || strings.Contains(text, "foreign-task") ||
		strings.Contains(text, "private-source") || strings.Contains(text, "foreign-source") {
		t.Fatalf("owner-auth Task navigation failed or disclosed denied source: %d %s", authorized.Code, text)
	}
	if session.calls != 2 || len(owner.calls) != 2*124 {
		t.Fatalf("off-page Task bypassed fresh owner authorization: %d calls, %d owners", session.calls, len(owner.calls))
	}
	owner.allowed["Task:task-100"] = false
	denied := returnHTML(h, selection+"&restore=1")
	if denied.Code != http.StatusNotFound || strings.Contains(denied.Body.String(), "task-100") {
		t.Fatalf("revoked Task exposed: %d %s", denied.Code, denied.Body.String())
	}
}

func TestLinkedTaskNavigationRejectsMalformedOrContextConfusingHintsBeforeOwner(t *testing.T) {
	owner, session := deepTaskOwner(), trustedHost()
	h := NewHTTPHandler(session, owner)
	for _, query := range []string{
		"view=incidents&selected=Task%3Atask-100&restore=1",
		"view=tasks&selected=Incident%3Atask-100&restore=1",
		"view=all&selected=Task%3Atask-100&restore=1",
		"view=tasks&selected=Task%3Atask-100&restore=0",
		"view=tasks&selected=Task%3Atask-100&restore=1&restore=1",
		"view=tasks&selected=Task%3Atask-100&restore=1&principal=admin",
		"view=tasks&selected=Task%3Atask-100&restore=1&tenant=t2",
		"view=tasks&selected=Task%3Atask-100&restore=1&cursor=deadbeef",
		"view=tasks&selected=Task%3Atask-100&restore=1&status=open",
		"view=tasks&selected=Task%3A%00&restore=1",
		"view=tasks&restore=1",
	} {
		got := returnHTML(h, query)
		if got.Code != http.StatusBadRequest || strings.Contains(got.Body.String(), "task-100") {
			t.Fatalf("invalid Task navigation hint accepted: %s -> %d %s", query, got.Code, got.Body.String())
		}
	}
	json := testQueueRequest(h, http.MethodGet,
		"/command/work-queue?view=tasks&selected=Task%3Atask-100&restore=1")
	if json.Code != http.StatusBadRequest ||
		strings.Contains(json.Body.String(), "task-100") ||
		session.calls != 0 || len(owner.calls) != 0 {
		t.Fatalf("malformed or non-HTML navigation reached owner: status=%d session=%d acl=%d", json.Code, session.calls, len(owner.calls))
	}
}

func TestIncidentDetailOffersOnlyIndependentlyAuthorizedTaskNavigation(t *testing.T) {
	o := newTaskRelationOwner()
	detail := NewIncidentDetailHTTPHandler(trustedHost(), o)
	visible := detailGET(detail, "/command/incident-detail?incident=inc-a")
	if visible.Code != http.StatusOK ||
		!strings.Contains(visible.Body.String(), "/command/work-queue?view=tasks&amp;selected=Task%3Atask-visible&amp;restore=1") ||
		strings.Contains(visible.Body.String(), "task-hidden") {
		t.Fatalf("Task navigation link must be owner-authorized: %d %s", visible.Code, visible.Body.String())
	}
	// The link is not a grant: fresh Work Queue projection checks Task ACL.
	work := NewHTTPHandler(trustedHost(), o)
	selected := returnHTML(work, "view=tasks&selected=Task%3Atask-visible&restore=1")
	if selected.Code != http.StatusOK || !strings.Contains(selected.Body.String(), "Task — task-visible") {
		t.Fatalf("selected Task owner reread failed: %d %s", selected.Code, selected.Body.String())
	}
	o.allowed["Task:task-visible"] = false
	revoked := returnHTML(work, "view=tasks&selected=Task%3Atask-visible&restore=1")
	if revoked.Code != http.StatusNotFound || strings.Contains(revoked.Body.String(), "task-visible") {
		t.Fatalf("Task link was wrongly treated as a grant: %d %s", revoked.Code, revoked.Body.String())
	}
	after := detailGET(detail, "/command/incident-detail?incident=inc-a")
	if after.Code != http.StatusOK || strings.Contains(after.Body.String(), "/command/work-queue?view=tasks") {
		t.Fatalf("Incident Detail retained revoked Task navigation: %d %s", after.Code, after.Body.String())
	}
}
