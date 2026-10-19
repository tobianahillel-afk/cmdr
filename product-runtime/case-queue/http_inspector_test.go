package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// CAP-INV-101 §17: Inspector reads only the current source-authorized Queue
// row, without transitioning to Case Workspace or granting extra Case access.
func TestNativeCaseInspectorReusesAuthorizedFilteredProjection(t *testing.T) {
	allowed := nativeCase("case-visible", "PowerShell investigation")
	allowed.Freshness = "partial"
	allowed.Incident = &IncidentContext{
		ID:"incident-a",TenantRef:"tenant-a",State:"available",Priority:"P2",
		AccessDecision:"deny",Owner:"hidden-command-owner",
	}
	allowed.Findings = []FindingContext{
		{ID:"finding-visible",CaseID:"case-visible",TenantRef:"tenant-a",
		 State:"available",Status:"review",AccessDecision:"deny"},
		{ID:"finding-secret",CaseID:"case-visible",TenantRef:"tenant-a",
		 State:"available",Status:"TOP-SECRET-FINDING",AccessDecision:"allow"},
	}
	denied:=nativeCase("case-secret","TOP-SECRET-CASE")
	other:=nativeCase("case-other","Unrelated")
	seen:=map[string]int{}
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {return nil,QueueScope{},errors.New("no session")}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		return []CaseSnapshot{denied,other,allowed},nil
	})
	acl:=queueAccessFunc(func(_ context.Context,scope QueueScope,kind,id string)(string,error){
		if scope.Tenant!="tenant-a"||scope.Environment!="prod" {t.Fatal("wrong trusted scope")}
		seen[kind+":"+id]++
		if id=="case-secret"||id=="incident-a"||id=="finding-secret" {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,acl)
	if err!=nil {t.Fatal(err)}
	got:=getQueue(h,"/case-queue?status=investigating&q=PowerShell&inspect_case=case-visible")
	if got.Code!=http.StatusOK {t.Fatalf("authorized inspect failed: %d %s",got.Code,got.Body.String())}
	var row Row
	if err:=json.Unmarshal(got.Body.Bytes(),&row);err!=nil {t.Fatal(err)}
	if row.Case.ID!="case-visible" || !row.Partial || row.Case.Incident!=nil ||
		len(row.Case.Findings)!=1 || row.Case.Findings[0].ID!="finding-visible" {
		t.Fatalf("Inspector does not preserve sanitized Case state: %+v",row)
	}
	for _,secret:=range []string{"case-secret","TOP-SECRET","incident-a","hidden-command-owner","finding-secret","case-other"} {
		if strings.Contains(got.Body.String(),secret) {t.Fatalf("inspector leaked forbidden context %q",secret)}
	}
	if got.Header().Get("Cache-Control")!="no-store" ||
		got.Header().Get("Vary")!="Authorization" {t.Fatal("inspector response is cacheable")}
	for _,want:=range []string{"case:case-visible","case:case-secret","incident:incident-a","finding:finding-visible","finding:finding-secret"} {
		if seen[want]==0 {t.Fatalf("Inspector skipped independent source ACL %s",want)}
	}
	for _,query:=range []string{
		"/case-queue?inspect_case=case-secret",
		"/case-queue?inspect_case=case-other&q=PowerShell",
		"/case-queue?inspect_case=case-visible&status=completed",
	} {
		denial:=getQueue(h,query)
		if denial.Code!=http.StatusNotFound || strings.Contains(denial.Body.String(),"case-") {
			t.Fatalf("not-visible Inspector selection leaked Case: %d %s",denial.Code,denial.Body.String())
		}
	}
}

func TestNativeCaseInspectorRespectsCurrentPageAndRevocation(t *testing.T) {
	cases:=[]CaseSnapshot{nativeCase("case-c","Charlie"),nativeCase("case-a","Alpha"),nativeCase("case-b","Bravo")}
	canReadB:=true
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){return cases,nil})
	acl:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		if kind=="case" && id=="case-b" && !canReadB {return "deny",nil}
		return "allow",nil
	})
	h,_:=NewAuthorizedQueueHTTPHandler(session,provider,acl)
	_,cursor:=decodeCasePage(t,casePageRequest(h,"/case-queue?limit=1&sort=title:asc"))
	if cursor=="" {t.Fatal("missing page cursor")}
	// Case B may only be inspected when it is on the current authorized page.
	first:=getQueue(h,"/case-queue?limit=1&sort=title:asc&inspect_case=case-b")
	if first.Code!=http.StatusNotFound {t.Fatal("inspector selected Case from outside current page")}
	target:="/case-queue?sort=title:asc&limit=1&cursor="+url.QueryEscape(cursor)+"&inspect_case=case-b"
	second:=getQueue(h,target)
	if second.Code!=http.StatusOK {t.Fatalf("cannot inspect current page Case: %d",second.Code)}
	var row Row
	if err:=json.Unmarshal(second.Body.Bytes(),&row);err!=nil || row.Case.ID!="case-b" {
		t.Fatalf("wrong page Inspector row: %+v %v",row,err)
	}
	canReadB=false
	revoked:=getQueue(h,target)
	if revoked.Code!=http.StatusNotFound || strings.Contains(revoked.Body.String(),"case-b") {
		t.Fatalf("revoked Case still inspected: %d %s",revoked.Code,revoked.Body.String())
	}
}

func TestNativeCaseInspectorRejectsInvalidControlsBeforeSession(t *testing.T) {
	calls:=0
	session:=queueSessionFunc(func(*http.Request)(context.Context,QueueScope,error){
		calls++;return context.Background(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("malformed Inspector request reached provider");return nil,nil
	})
	acl:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("malformed Inspector request reached ACL");return "",nil
	})
	h,_:=NewAuthorizedQueueHTTPHandler(session,provider,acl)
	for _,query:=range []string{
		"inspect_case=",
		"inspect_case=%20%20",
		"inspect_case=case-a&inspect_case=case-b",
		"inspect_case=case-a&open_case=case-a",
		"inspect_case=case-a&scroll=1",
		"inspect_case=case-a&tenant=tenant-b",
		"inspect_case=case-a&principal=admin",
		"inspect_case=case-a%00",
		"inspect_case="+strings.Repeat("a",1025),
	} {
		got:=getQueue(h,"/case-queue?"+query)
		if got.Code!=http.StatusBadRequest {t.Fatalf("malformed Inspector query accepted: %s => %d",query,got.Code)}
	}
	if calls!=0 {t.Fatalf("unsafe inspector requests reached session resolver: %d",calls)}
}
