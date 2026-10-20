package workqueue

import (
	"context"
	"testing"
)

// CAP-CMD-101 Class-0: the SLA-risk system view must not report a degraded
// filtered cohort merely because another authorized status has missing SLA
// policy. Other-tenant and denied owner items must never influence the state.
func TestSLARiskPartialDiagnosticRespectsAuthorizedStatusFilter(t *testing.T) {
	o:=&fakeOwner{
		items:[]Item{
			{Kind:Incident,ID:"inc-open-at-risk",Tenant:"t1",Environment:"e1",Status:"open",UpdatedAt:"2026-10-19T10:00:00Z",SLAAtRisk:enabled(true)},
			{Kind:Task,ID:"task-closed-no-sla",Tenant:"t1",Environment:"e1",Status:"closed",UpdatedAt:"2026-10-19T10:01:00Z"},
			{Kind:Task,ID:"task-denied-no-sla",Tenant:"t1",Environment:"e1",Status:"open",UpdatedAt:"2026-10-19T10:02:00Z"},
			{Kind:Incident,ID:"foreign-no-sla",Tenant:"other",Environment:"e1",Status:"open",UpdatedAt:"2026-10-19T10:03:00Z"},
		},
		allowed:map[string]bool{"Incident:inc-open-at-risk":true,"Task:task-closed-no-sla":true},
	}
	open:=request(SLARisk)
	open.StatusFilter="open"
	r,err:=Project(context.Background(),o,open)
	if err!=nil { t.Fatalf("authorized filtered queue rejected: %v",err) }
	if r.State!="ready" || len(r.Diagnostics)!=0 || len(r.Rows)!=1 ||
		r.Rows[0].ID!="inc-open-at-risk" {
		t.Fatalf("off-filter SLA gap contaminated the visible source state: %+v",r)
	}
	for _,query:=range o.calls {
		if query=="Incident:foreign-no-sla" {
			t.Fatal("owner was queried for other-tenant ID")
		}
	}
	closed:=request(SLARisk)
	closed.StatusFilter="closed"
	r,err=Project(context.Background(),o,closed)
	if err!=nil || r.State!="partial-sources" || len(r.Rows)!=0 ||
		len(r.Diagnostics)!=1 || r.Diagnostics[0]!="sla-policy-unavailable" {
		t.Fatalf("in-scope missing SLA policy must still report partial data: %v %+v",err,r)
	}
	unfiltered,err:=Project(context.Background(),o,request(SLARisk))
	if err!=nil || unfiltered.State!="partial-sources" ||
		len(unfiltered.Rows)!=1 || len(unfiltered.Diagnostics)!=1 ||
		unfiltered.Diagnostics[0]!="sla-policy-unavailable" {
		t.Fatalf("legacy unfiltered SLA risk partial behavior changed: %v %+v",err,unfiltered)
	}
	// No new exposure or status change for a fully authorized ordinary view.
	tasks:=request(Tasks)
	tasks.StatusFilter="open"
	r,err=Project(context.Background(),o,tasks)
	if err!=nil || r.State!="no-authorized-items" ||
		len(r.Rows)!=0 || len(r.Diagnostics)!=0 {
		t.Fatalf("ordinary authorized Tasks view unexpectedly degraded: %v %+v",err,r)
	}
}
