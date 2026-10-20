package workqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// CAP-CMD-109's view-level aggregate includes only independently authorized
// items in the current canonical view/status filter. Counts are never inferred
// from hidden sources or the owner's Incident/Task UpdatedAt timestamp.
func TestViewFreshnessExplainsAuthorizedCohortWithoutScopeLeaks(t *testing.T) {
	owner := fixture()
	owner.items[0].SourceFreshness = &SourceFreshness{Source:"allowed-incident",ObservedAt:freshAt(-5*time.Second),Health:"available",PolicyVersion:"incident-policy-v1",AgingAfterSeconds:60,StaleAfterSeconds:120}
	owner.items[1].SourceFreshness = &SourceFreshness{Source:"allowed-task",ObservedAt:freshAt(-3*time.Minute),Health:"available",PolicyVersion:"task-policy-v2",AgingAfterSeconds:60,StaleAfterSeconds:120}
	owner.items[2].SourceFreshness = &SourceFreshness{Source:"SECRET-DENIED-INCIDENT",Health:"unavailable"}
	owner.items[4].SourceFreshness = &SourceFreshness{Source:"SECRET-OTHER-TENANT",Health:"unavailable"}
	// An authorized item with no optional policy/source data is truly
	// unknown; never silently label it current.
	got,err := Project(context.Background(),owner,request(All))
	if err!=nil {t.Fatal(err)}
	f:=got.ViewFreshness
	if f==nil || f.Total!=3 || f.Current!=1 || f.Stale!=1 || f.Unknown!=1 ||
		f.Aging!=0 || f.Partial!=0 || f.Unavailable!=0 {
		t.Fatalf("incorrect authorized all-view source inventory: %+v",f)
	}
	b,err:=json.Marshal(got)
	if err!=nil {t.Fatal(err)}
	for _,secret:=range []string{"SECRET-DENIED-INCIDENT","SECRET-OTHER-TENANT"} {
		if strings.Contains(string(b),secret) {t.Fatalf("denied source metadata escaped: %s",secret)}
	}
	inc,err:=Project(context.Background(),owner,request(Incidents))
	if err!=nil || inc.ViewFreshness==nil || inc.ViewFreshness.Total!=1 ||
		inc.ViewFreshness.Current!=1 || inc.ViewFreshness.Stale!=0 {
		t.Fatalf("task freshness contaminated Incident view: %v %+v",err,inc.ViewFreshness)
	}
	taskReq:=request(Tasks)
	taskReq.StatusFilter="open"
	task,err:=Project(context.Background(),owner,taskReq)
	if err!=nil || task.ViewFreshness==nil || task.ViewFreshness.Total!=2 ||
		task.ViewFreshness.Stale!=1 || task.ViewFreshness.Unknown!=1 {
		t.Fatalf("filtered view counting incorrect: %v %+v",err,task.ViewFreshness)
	}
}

func TestViewFreshnessCountsAllAuthorizedPagesAndRevalidatesCursor(t *testing.T) {
	owner:=fixture()
	owner.items=nil
	owner.allowed=map[string]bool{}
	for i:=0;i<81;i++ {
		id:=fmt.Sprintf("task-%03d",i)
		fresh:=&SourceFreshness{Source:"command-owned",ObservedAt:freshAt(-5*time.Second),Health:"available",PolicyVersion:"t1-e1-v1",AgingAfterSeconds:60,StaleAfterSeconds:120}
		if i>=60 {fresh.ObservedAt=freshAt(-5*time.Minute)}
		owner.items=append(owner.items,Item{Kind:Task,ID:id,Tenant:"t1",Environment:"e1",Owner:"analyst",Status:"open",UpdatedAt:"2026-10-19T09:00:00Z",SourceFreshness:fresh})
		owner.allowed["Task:"+id]=true
	}
	first,err:=Project(context.Background(),owner,request(Tasks))
	if err!=nil || first.NextCursor=="" || len(first.Rows)!=50 {t.Fatalf("first-page snapshot: %v %+v",err,first)}
	f:=first.ViewFreshness
	if f==nil || f.Total!=81 || f.Current!=60 || f.Stale!=21 {
		t.Fatalf("view summary only counted displayed page: %+v",f)
	}
	secondRequest:=request(Tasks)
	secondRequest.Cursor=first.NextCursor
	second,err:=Project(context.Background(),owner,secondRequest)
	if err!=nil || second.ViewFreshness==nil || second.ViewFreshness.Total!=81 ||
		second.ViewFreshness.Stale!=21 || len(second.Rows)!=31 {
		t.Fatalf("subsequent page changes authorized cohort summary: %v %+v",err,second.ViewFreshness)
	}
	// Revocation makes a prior-page token stale. Neither old rows nor
	// their formerly authorized freshness counts may be returned.
	owner.allowed["Task:task-060"]=false
	revoked,err:=Project(context.Background(),owner,secondRequest)
	if err!=ErrCursor || revoked.ViewFreshness!=nil || len(revoked.Rows)!=0 {
		t.Fatalf("revoked/stale cursor exposed aggregate: %v %+v",err,revoked)
	}
}

func TestViewFreshnessDegradedAndUnknownStatesRemainDistinct(t *testing.T) {
	owner:=fixture()
	owner.items[0].SourceFreshness=&SourceFreshness{Source:"source-INC",Health:"partial"}
	owner.items[1].SourceFreshness=&SourceFreshness{Source:"source-TASK",Health:"unavailable"}
	owner.items[3].SourceFreshness=&SourceFreshness{Source:"source-NO-POLICY",ObservedAt:freshAt(-time.Hour),Health:"available"}
	all,err:=Project(context.Background(),owner,request(All))
	if err!=nil || all.State!="partial-sources" || all.ViewFreshness==nil ||
		all.ViewFreshness.Partial!=1 || all.ViewFreshness.Unavailable!=1 ||
		all.ViewFreshness.Unknown!=1 {
		t.Fatalf("unknown confused with degraded owner health: %v %+v",err,all)
	}
	handler:=NewHTTPHandler(trustedHost(),owner)
	html:=testQueueRequestHTML(handler,http.MethodGet,"/command/work-queue?view=all")
	// The per-row table MUST legitimately show authorized sources; only
	// the NEW view-level summary must contain aggregate integers.
	prefix,rest:=strings.Cut(html.Body.String(),"<p aria-label=\"Fraîcheur de la vue autorisée\">")
	summary,found:=strings.Cut(rest,"</p>")
	if html.Code!=200 || !prefix || !found ||
		!strings.Contains(summary,"3 éléments") ||
		strings.Contains(summary,"source-INC") ||
		strings.Contains(summary,"source-TASK") {
		t.Fatalf("authorized freshness aggregate missing or exposing owner labels: %d %q",html.Code,summary)
	}
}

func TestViewFreshnessPreservesLegacyJSONAndDeniesInvalidOwners(t *testing.T) {
	legacy,err:=Project(context.Background(),fixture(),request(All))
	if err!=nil || legacy.ViewFreshness!=nil {t.Fatalf("legacy provider unexpectedly received invented source data: %v",err)}
	b,err:=json.Marshal(legacy)
	if err!=nil {t.Fatal(err)}
	if strings.Contains(string(b),"view_freshness") {t.Fatalf("legacy JSON contract changed: %q",b)}
	invalid:=fixture()
	invalid.items[0].SourceFreshness=&SourceFreshness{Source:"bad\nsource",Health:"unavailable"}
	result,err:=Project(context.Background(),invalid,request(All))
	if err!=ErrIncomplete || result.ViewFreshness!=nil || len(result.Rows)!=0 {
		t.Fatalf("malformed authorized metadata returned partial counts: %v %+v",err,result)
	}
}
