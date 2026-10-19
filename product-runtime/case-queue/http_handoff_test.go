package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestNativeQueueHTTPOpensOnlyCurrentAuthorizedVisibleCase(t *testing.T) {
	visible := nativeCase(" case-opaque ", "PowerShell investigation")
	visible.Status = "review"
	visible.Owner = "analyst-a"
	denied := nativeCase("case-denied", "PowerShell denied investigation")
	denied.Status = "review"
	denied.Owner = "analyst-a"
	sessionCalls, caseChecks := 0, []string{}
	session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		sessionCalls++
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("unverified user")
		}
		return r.Context(), QueueScope{Tenant:"tenant-a", Environment:"prod"}, nil
	})
	provider := queueProviderFunc(func(_ context.Context, scope QueueScope) ([]CaseSnapshot, error) {
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {t.Fatal("untrusted scope")}
		return []CaseSnapshot{visible,denied},nil
	})
	access := queueAccessFunc(func(_ context.Context, _ QueueScope, kind, id string)(string,error) {
		if kind!="case" {t.Fatalf("unexpected owner resource %s",kind)}
		caseChecks=append(caseChecks,id)
		if id=="case-denied" {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,access)
	if err!=nil {t.Fatal(err)}
	p := url.Values{
		"q":{"PowerShell"},"status":{"review"},"owner":{"analyst-a"},
		"freshness":{"available"},"sort":{"title:asc"},
		"open_case":{visible.ID},"scroll":{"128"},
	}
	got:=getQueue(h,"/case-queue?"+p.Encode())
	if got.Code!=http.StatusOK {t.Fatalf("authorized selection HTTP %d %q",got.Code,got.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(got.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	if handoff.ScreenID!=CaseWorkspaceScreen || handoff.CaseID!=visible.ID ||
		handoff.TenantRef!="tenant-a" || handoff.EnvironmentRef!="prod" ||
		handoff.ReturnToken=="" {
		t.Fatalf("native handoff does not preserve canonical Case context: %+v",handoff)
	}
	if len(caseChecks)!=2 || caseChecks[0]!=visible.ID || caseChecks[1]!="case-denied" {
		t.Fatalf("did not authorize every Case before selection/filter: %v",caseChecks)
	}
	current:=Input{TenantRef:"tenant-a",EnvironmentRef:"prod",
		AuthorizedTenants:[]string{"tenant-a"},
		Cases:[]CaseSnapshot{visible}}
	state,code:=RestoreCaseQueue(handoff.ReturnToken,current)
	if code!=ErrorNone || state.SelectedCaseID!=visible.ID || state.ScrollOffset!=128 ||
		state.SearchQuery!="PowerShell" || len(state.Filters)!=3 || len(state.Sort)!=1 ||
		state.Sort[0].Field!="title" {
		t.Fatalf("native handoff not restorable by existing Case Queue engine: %v %+v",code,state)
	}
	if sessionCalls!=1 {t.Fatalf("duplicate/different session opened: %d",sessionCalls)}
	for _,secret:=range []string{"case-denied","PowerShell denied investigation","Bearer operator"} {
		if strings.Contains(got.Body.String(),secret) {t.Fatalf("handoff disclosed forbidden Case/token: %q",secret)}
	}
	if got.Header().Get("Cache-Control")!="no-store" {
		t.Fatal("navigation state leaked via cacheable HTTP")
	}
	// Current permission is rechecked at restore, not transferred by token.
	current.Cases=nil
	if _,code:=RestoreCaseQueue(handoff.ReturnToken,current);code==ErrorNone {
		t.Fatal("return token granted Case access after current ACL changed")
	}
}

func TestNativeQueueHTTPHandoffDeniesDeniedStaleOrFilteredSelection(t *testing.T) {
	allowed:=nativeCase("case-allowed","Investigate endpoint")
	stale:=nativeCase("case-stale","Investigate stale")
	stale.Freshness="stale"
	denied:=nativeCase("case-denied","Investigate restricted")
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {return nil,QueueScope{},errors.New("denied")}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		return []CaseSnapshot{allowed,stale,denied},nil
	})
	access:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		if kind!="case" {t.Fatal("unexpected nested authorization")}
		if id=="case-denied" {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,access)
	if err!=nil {t.Fatal(err)}
	for _,tc:=range []struct{name,id,query string}{
		{"denied owner","case-denied",""},
		{"stale permitted Case","case-stale",""},
		{"filtered out Case","case-allowed","?q=unmatched"},
		{"nonexistent Case","missing-case",""},
	}{
		t.Run(tc.name,func(t *testing.T){
			u:="/case-queue?open_case="+url.QueryEscape(tc.id)
			if tc.query!="" {u+="&"+strings.TrimPrefix(tc.query,"?")}
			got:=getQueue(h,u)
			if got.Code!=http.StatusNotFound || strings.Contains(got.Body.String(),tc.id) ||
				strings.Contains(got.Body.String(),"return_token") {
				t.Fatalf("server emitted denied/stale/filtered Case navigation: %d %q",got.Code,got.Body.String())
			}
		})
	}
	// The legacy read-only listing endpoint remains intact with no selection.
	got:=getQueue(h,"/case-queue?sort=title:asc")
	if got.Code!=http.StatusOK || strings.Contains(got.Body.String(),"case-denied") {
		t.Fatalf("existing Queue projection regressed: %d %s",got.Code,got.Body.String())
	}
	var projection Projection
	if err:=json.Unmarshal(got.Body.Bytes(),&projection);err!=nil {t.Fatal(err)}
	if len(projection.Rows)!=2 {t.Fatal("regular Queue no longer includes accessible Cases")}
}

func TestNativeQueueHTTPHandoffRejectsInvalidQueryBeforeAuthorization(t *testing.T) {
	sessionCalls:=0
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		sessionCalls++
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("invalid navigation request accessed provider")
		return nil,nil
	})
	auth:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("invalid navigation request accessed owner authorization")
		return "",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,auth)
	if err!=nil {t.Fatal(err)}
	for _,q:=range []string{
		"scroll=1","open_case=","open_case=%00",
		"open_case=one&open_case=two",
		"open_case=case-ok&scroll=-1",
		"open_case=case-ok&scroll=not-an-integer",
		"open_case=case-ok&scroll="+strconv.Itoa(maxReturnScroll+1),
		"open_case=case-ok&permission=allow",
		"open_case=case-ok&tenant=tenant-a",
		"open_case=case-ok&origin=command-incident",
		"open_case=case-ok&status=review&status=triage",
		"open_case="+strings.Repeat("x",1025),
	}{
		got:=getQueue(h,"/case-queue?"+q)
		if got.Code!=http.StatusBadRequest || strings.Contains(got.Body.String(),"return_token") {
			t.Fatalf("malformed request reached authorization: %q %d",q,got.Code)
		}
	}
	if sessionCalls!=0 {t.Fatalf("invalid request reached session resolver %d times",sessionCalls)}
}
