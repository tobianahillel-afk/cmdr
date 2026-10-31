package workqueue

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type taskRelationOwner struct {
	*fakeOwner
	links []IncidentTaskReference
	linkErr error
}

func (o *taskRelationOwner) ReadAuthorizedIncidentTaskReferences(ctx context.Context,
	principal, tenant, environment, incidentID string) ([]IncidentTaskReference, error) {
	if ctx.Err() != nil || o.linkErr != nil {
		return nil, errors.New("private Task relationship failure")
	}
	if principal != "analyst" || tenant != "t1" || environment != "e1" || incidentID != "inc-a" {
		return nil, nil
	}
	return o.links, nil
}

func taskScope() SessionScope {
	return SessionScope{Principal: "analyst", Tenant: "t1", Environment: "e1"}
}

func newTaskRelationOwner() *taskRelationOwner {
	o := &fakeOwner{allowed: map[string]bool{
		"Incident:inc-a": true, "Task:task-visible": true,
	}}
	o.items = []Item{
		{Kind: Incident, ID: "inc-a", Tenant: "t1", Environment: "e1",
			Owner: "analyst", Status: "active", UpdatedAt: "2026-10-21T02:00:00Z"},
		{Kind: Task, ID: "task-visible", Tenant: "t1", Environment: "e1",
			Owner: "responder", Status: "in-progress", UpdatedAt: "2026-10-21T03:00:00Z"},
		{Kind: Task, ID: "task-hidden", Tenant: "t1", Environment: "e1",
			Owner: "PRIVATE-TASK-OWNER", Status: "blocked", UpdatedAt: "2026-10-21T03:01:00Z"},
	}
	o.allowed["Task:task-hidden"] = false
	return &taskRelationOwner{fakeOwner: o, links: []IncidentTaskReference{
		{IncidentID: "inc-a", TaskID: "task-visible", Tenant: "t1", Environment: "e1"},
		{IncidentID: "inc-a", TaskID: "task-hidden", Tenant: "t1", Environment: "e1"},
	}}
}

func TestIncidentTaskLinksUseIndependentTaskACLAndActualHTMLRoute(t *testing.T) {
	o := newTaskRelationOwner()
	out, err := readAuthorizedIncidentTasks(context.Background(), o, taskScope(), "inc-a")
	if err != nil || len(out) != 1 || out[0].ID != "task-visible" ||
		out[0].Owner != "responder" || out[0].Status != "in-progress" {
		t.Fatalf("visible Task not independently authorized: %+v, %v", out, err)
	}
	if !strings.Contains(strings.Join(o.calls, ","), "Task:task-visible") ||
		!strings.Contains(strings.Join(o.calls, ","), "Task:task-hidden") {
		t.Fatalf("linked Task read did not check both objects: %v", o.calls)
	}
	h := NewIncidentDetailHTTPHandler(trustedHost(), o)
	resp := detailGET(h, "/command/incident-detail?incident=inc-a")
	body := resp.Body.String()
	if resp.Code != http.StatusOK || !strings.Contains(body, "Task task-visible") ||
		!strings.Contains(body, "responder") || !strings.Contains(body, "in-progress") ||
		strings.Contains(body, "task-hidden") || strings.Contains(body, "PRIVATE-TASK-OWNER") ||
		strings.Contains(body, "task-visible?") {
		t.Fatalf("denied Task or hidden Task source leaked: %d %s", resp.Code, body)
	}
	o.allowed["Task:task-visible"] = false
	denied := detailGET(h, "/command/incident-detail?incident=inc-a")
	if denied.Code != http.StatusOK ||
		!strings.Contains(denied.Body.String(), "Tasks liées indisponibles ou non autorisées") ||
		strings.Contains(denied.Body.String(), "task-visible") {
		t.Fatalf("revoked Task permission leaked: %d %s", denied.Code, denied.Body.String())
	}
}

func TestIncidentTaskLinksInvalidRelationOrSourceFailsWithoutPartialTaskOutput(t *testing.T) {
	baseline := IncidentTaskReference{IncidentID:"inc-a", TaskID:"task-visible", Tenant:"t1", Environment:"e1"}
	tests := []struct {
		name string
		alter func(*taskRelationOwner)
	}{
		{"wrong-incident", func(o *taskRelationOwner){ o.links[0].IncidentID = "foreign-incident" }},
		{"wrong-tenant", func(o *taskRelationOwner){ o.links[0].Tenant = "tenant-other" }},
		{"wrong-environment", func(o *taskRelationOwner){ o.links[0].Environment = "secret-env" }},
		{"malformed-task-id", func(o *taskRelationOwner){ o.links[0].TaskID = "task\n<script>" }},
		{"duplicate-ref", func(o *taskRelationOwner){ o.links=append(o.links, baseline) }},
		{"too-many-links", func(o *taskRelationOwner){
			for len(o.links) <= 12 { o.links=append(o.links, IncidentTaskReference{
				IncidentID:"inc-a",TaskID:"extra",Tenant:"t1",Environment:"e1",
			}) }
		}},
		{"links-fail", func(o *taskRelationOwner){ o.linkErr=errors.New("private Task relationship failure") }},
		{"task-owner-read-fails", func(o *taskRelationOwner){ o.errRead=true }},
		{"task-list-fails", func(o *taskRelationOwner){ o.errList=true }},
		{"duplicate-task-owner-row", func(o *taskRelationOwner){ o.items=append(o.items,o.items[1]) }},
		{"foreign-task-in-owner", func(o *taskRelationOwner){ o.items[1].Tenant="other-tenant" }},
		{"unsafe-owner-text", func(o *taskRelationOwner){ o.items[1].Owner="<script>\nprivate" }},
		{"unsafe-task-status", func(o *taskRelationOwner){ o.items[1].Status="active\rprivate" }},
		{"invalid-timestamp", func(o *taskRelationOwner){ o.items[1].UpdatedAt="tomorrow" }},
	}
	for _, tc := range tests {
		t.Run(tc.name,func(t *testing.T){
			o:=newTaskRelationOwner()
			tc.alter(o)
			resp:=detailGET(NewIncidentDetailHTTPHandler(trustedHost(),o),
				"/command/incident-detail?incident=inc-a")
			if resp.Code!=http.StatusServiceUnavailable ||
				strings.Contains(resp.Body.String(),"task-visible") ||
				strings.Contains(resp.Body.String(),"task-hidden") ||
				strings.Contains(resp.Body.String(),"responder") ||
				strings.Contains(resp.Body.String(),"PRIVATE-TASK-OWNER") {
				t.Fatalf("invalid link/task emitted HTML: %d %s",resp.Code,resp.Body.String())
			}
		})
	}
}

func TestIncidentTaskLinksOptionalMissingRevokedAndBounded(t *testing.T) {
	o:=newTaskRelationOwner()
	o.links=nil
	out,err:=readAuthorizedIncidentTasks(context.Background(),o,taskScope(),"inc-a")
	if err!=nil||len(out)!=0 { t.Fatalf("unknown links should be partial: %v %+v",err,out) }
	o=newTaskRelationOwner()
	o.links[0].TaskID="missing-task"
	out,err=readAuthorizedIncidentTasks(context.Background(),o,taskScope(),"inc-a")
	if err!=nil||len(out)!=0 { t.Fatalf("missing linked owner Task became visible: %v %+v",err,out) }
	o=newTaskRelationOwner()
	ctx,cancel:=context.WithCancel(context.Background())
	cancel()
	out,err=readAuthorizedIncidentTasks(ctx,o,taskScope(),"inc-a")
	if err==nil || len(out)!=0 { t.Fatalf("canceled authority accepted: %v %+v",err,out) }
	o=newTaskRelationOwner()
	o.items=append(o.items,make([]Item,5000)...)
	_,err=readAuthorizedIncidentTasks(context.Background(),o,taskScope(),"inc-a")
	if err==nil { t.Fatal("oversized Task owner snapshot accepted") }
	legacy:=detailGET(NewIncidentDetailHTTPHandler(trustedHost(),deepReturnOwner()),
		"/command/incident-detail?incident=inc-100")
	if legacy.Code!=http.StatusOK ||
		!strings.Contains(legacy.Body.String(),"Tasks liées indisponibles ou non autorisées") {
		t.Fatalf("legacy owner misinterpreted absence as Tasks: %d",legacy.Code)
	}
}
