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

func withParent(caseID, title, parent string) CaseSnapshot {
	c := nativeCase(caseID, title)
	if parent != "" {
		c.Incident = &IncidentContext{
			ID: parent, TenantRef: "tenant-a", AccessDecision: "allow",
			State: "available", Priority: "P2", Owner: "command-analyst",
		}
	}
	return c
}

func TestNativeIncidentLinkedQueueRequiresIndependentSourceAuthorization(t *testing.T) {
	visible := withParent("case-a", "Allowed linked Case", "INC-CaseSensitive")
	casedenied := withParent("case-b", "HIDDEN-DENIED-CASE", "INC-CaseSensitive")
	other := withParent("case-c", "Different casefold parent", "inc-casesensitive")
	parentdenied := withParent("case-d", "HIDDEN-DENIED-INCIDENT", "INC-Secret")
	without := withParent("case-e", "No Incident", "")
	calls := map[string]int{}
	session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("unauthenticated")
		}
		return r.Context(), QueueScope{Tenant:"tenant-a", Environment:"prod"},nil
	})
	provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot,error) {
		return []CaseSnapshot{visible,casedenied,other,parentdenied,without},nil
	})
	acl := queueAccessFunc(func(_ context.Context,scope QueueScope,kind,id string)(string,error){
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {t.Fatal("scope injection")}
		calls[kind+":"+id]++
		if (kind=="case" && id=="case-b") || (kind=="incident" && id=="INC-Secret") {
			return "deny",nil
		}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,acl)
	if err!=nil {t.Fatal(err)}
	matched:=getQueue(h,"/case-queue?incident_id=INC-CaseSensitive")
	if matched.Code!=http.StatusOK {t.Fatalf("linked filter returned %d: %s",matched.Code,matched.Body.String())}
	var result Projection
	if err:=json.Unmarshal(matched.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
	if len(result.Rows)!=1 || result.Rows[0].Case.ID!="case-a" ||
		result.Rows[0].Case.Incident==nil || result.Rows[0].Case.Incident.ID!="INC-CaseSensitive" {
		t.Fatalf("exact parent filter wrongly matched denied or lowercased relation: %+v",result.Rows)
	}
	for _,secret:=range []string{"HIDDEN-DENIED","case-b","case-d","INC-Secret","inc-casesensitive"} {
		if strings.Contains(matched.Body.String(),secret) {t.Fatalf("Incident navigation disclosed denied/other context: %q",secret)}
	}
	if calls["case:case-b"]==0 || calls["incident:INC-Secret"]==0 ||
		calls["incident:INC-CaseSensitive"]==0 {
		t.Fatalf("Case and Incident must be independently authorized before projection: %+v",calls)
	}
	noAccess:=getQueue(h,"/case-queue?incident_id=INC-Secret")
	if noAccess.Code!=http.StatusOK || strings.Contains(noAccess.Body.String(),"INC-Secret") {
		t.Fatalf("denied Command Incident source exposed: %d %s",noAccess.Code,noAccess.Body.String())
	}
	var empty Projection
	if err:=json.Unmarshal(noAccess.Body.Bytes(),&empty);err!=nil {t.Fatal(err)}
	if len(empty.Rows)!=0 {t.Fatal("Case with denied Command Incident unexpectedly matched")}
	notAlias:=getQueue(h,"/case-queue?incident_id=inc-casesensitive")
	if notAlias.Code!=http.StatusOK {t.Fatal("valid lowercased opaque ID failed")}
	var exact Projection
	if err:=json.Unmarshal(notAlias.Body.Bytes(),&exact);err!=nil {t.Fatal(err)}
	if len(exact.Rows)!=1 || exact.Rows[0].Case.ID!="case-c" {
		t.Fatalf("opaque Incident identities incorrectly casefolded: %+v",exact.Rows)
	}
}

func TestNativeIncidentLinkedPagingWorkspaceReturnRechecksOwnerAndSelector(t *testing.T) {
	a:=withParent("case-a","Alpha","INC-77")
	b:=withParent("case-b","Bravo","INC-77")
	c:=withParent("case-c","Charlie","INC-77")
	notLinked:=withParent("case-hidden","HIDDEN-UNRELATED","INC-OTHER")
	incidentAllowed:=true
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error) {
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error) {
		return []CaseSnapshot{notLinked,c,a,b},nil
	})
	acl:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		if kind=="incident" && id=="INC-77" && !incidentAllowed {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,acl)
	if err!=nil {t.Fatal(err)}
	first,cursor:=decodeCasePage(t,casePageRequest(h,
		"/case-queue?incident_id=INC-77&sort=title:asc&limit=2"))
	if len(first.Rows)!=2 || cursor=="" {t.Fatalf("missing Incident-linked first page: %+v",first)}
	second,_:=decodeCasePage(t,casePageRequest(h,
		"/case-queue?incident_id=INC-77&sort=title:asc&limit=2&cursor="+url.QueryEscape(cursor)))
	if len(second.Rows)!=1 || second.Rows[0].Case.ID!="case-c" {
		t.Fatalf("wrong Incident-linked second page: %+v",second.Rows)
	}
	open:=url.Values{
		"incident_id":{"INC-77"}, "sort":{"title:asc"},
		"limit":{"2"},"cursor":{cursor}, "open_case":{"case-c"},"scroll":{"73"},
	}
	selected:=getQueue(h,"/case-queue?"+open.Encode())
	if selected.Code!=http.StatusOK {t.Fatalf("linked Case could not open: %d %s",selected.Code,selected.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(selected.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	body,_:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	back:=postQueueReturn(h,"application/json",string(body))
	if back.Code!=http.StatusOK {t.Fatalf("linked Case could not return: %d %s",back.Code,back.Body.String())}
	var state ReturnState
	if err:=json.Unmarshal(back.Body.Bytes(),&state);err!=nil {t.Fatal(err)}
	if state.SelectedCaseID!="case-c" || state.PageLimit!=2 ||
		state.PageCursor!=cursor || state.ScrollOffset!=73 ||
		len(state.Filters)!=1 || state.Filters[0].Field!="incident_id" ||
		len(state.Filters[0].Values)!=1 || state.Filters[0].Values[0]!="INC-77" {
		t.Fatalf("Command-to-Investigate selector/page lost in workspace transition: %+v",state)
	}
	incidentAllowed=false
	revoked:=postQueueReturn(h,"application/json",string(body))
	if revoked.Code!=http.StatusNotFound || strings.Contains(revoked.Body.String(),"case-c") ||
		strings.Contains(revoked.Body.String(),"INC-77") {
		t.Fatalf("Incident source access revoked but returned stale selected Case: %d %s",revoked.Code,revoked.Body.String())
	}
}

func TestNativeIncidentLinkedQueueRejectsUntrustedControlsBeforeSession(t *testing.T) {
	calls:=0
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		calls++; return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	cases:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("invalid incident query reached Case provider");return nil,nil
	})
	acl:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("invalid incident query reached owner auth");return "",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,cases,acl)
	if err!=nil {t.Fatal(err)}
	for _,q:=range []string{
		"incident_id=",
		"incident_id=%20%20",
		"incident_id=INC-77&incident_id=INC-88",
		"incident_id=INC-77&tenant=tenant-b",
		"incident_id=INC-77&principal=admin",
		"incident_id="+strings.Repeat("a",1025),
		"incident_id=INC-77%00",
	} {
		response:=getQueue(h,"/case-queue?"+q)
		if response.Code!=http.StatusBadRequest {
			t.Fatalf("invalid Command incident selector accepted: %q => %d",q,response.Code)
		}
	}
	if calls!=0 {t.Fatalf("untrusted Command Incident selectors reached session resolver %d times",calls)}
}
