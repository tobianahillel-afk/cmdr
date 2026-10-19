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

// Exercise the actual HTTP Queue -> Workspace -> Return round trip from a
// second authorized page, not just the pagination and handoff helpers alone.
func TestNativePagedQueueOpenAndReturnKeepsCurrentAuthorizedPage(t *testing.T) {
	a, b, c := nativeCase("case-a", "Alpha"), nativeCase("case-b", "Bravo"), nativeCase("case-c", "Charlie")
	hidden := nativeCase("private-case", "HIDDEN-DENIED-CASE")
	ownerAllowed := true
	aclCalls := 0
	session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("not authenticated")
		}
		return r.Context(), QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot,error) {
		return []CaseSnapshot{hidden,c,a,b},nil
	})
	access := queueAccessFunc(func(_ context.Context, _ QueueScope, kind,id string)(string,error){
		aclCalls++
		if kind!="case" {t.Fatalf("unexpected source kind: %s",kind)}
		if id=="private-case" || (id=="case-c" && !ownerAllowed) {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,access)
	if err!=nil {t.Fatal(err)}
	first,cursor:=decodeCasePage(t,casePageRequest(h,"/case-queue?limit=2&sort=title:asc"))
	if len(first.Rows)!=2 || cursor=="" {t.Fatal("missing authorized second page")}
	second,last:=decodeCasePage(t,casePageRequest(h,
		"/case-queue?limit=2&sort=title:asc&cursor="+url.QueryEscape(cursor)))
	if len(second.Rows)!=1 || second.Rows[0].Case.ID!="case-c" || last!="" {
		t.Fatalf("wrong subsequent page: %+v",second.Rows)
	}
	openUrl:=url.Values{
		"limit":{"2"},"sort":{"title:asc"},"cursor":{cursor},
		"open_case":{"case-c"},"scroll":{"91"},
	}
	opened:=getQueue(h,"/case-queue?"+openUrl.Encode())
	if opened.Code!=http.StatusOK {t.Fatalf("cannot open Case from its authorized second page: %d %q",opened.Code,opened.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(opened.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	if handoff.CaseID!="case-c" || handoff.ReturnToken=="" || handoff.ScreenID!=CaseWorkspaceScreen {
		t.Fatalf("paged handoff incorrect: %+v",handoff)
	}
	raw,_:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	returned:=postQueueReturn(h,"application/json",string(raw))
	if returned.Code!=http.StatusOK {t.Fatalf("cannot restore page and Case: %d %q",returned.Code,returned.Body.String())}
	var state ReturnState
	if err:=json.Unmarshal(returned.Body.Bytes(),&state);err!=nil {t.Fatal(err)}
	if state.PageLimit!=2 || state.PageCursor!=cursor ||
		state.SelectedCaseID!="case-c" || state.ScrollOffset!=91 ||
		len(state.Sort)!=1 || state.Sort[0].Field!="title" {
		t.Fatalf("round trip lost source page, selection, filters or scroll: %+v",state)
	}
	if aclCalls!=20 { // five actual HTTP operations * all four source ACL checks
		t.Fatalf("some paged open/return skipped a current Case ACL pass: %d",aclCalls)
	}
	for _,v:=range []string{"private-case","HIDDEN-DENIED-CASE","Bearer operator"} {
		if strings.Contains(opened.Body.String(),v)||strings.Contains(returned.Body.String(),v) {
			t.Fatalf("paged navigation leaked private Case %q",v)
		}
	}
	// Old cursor must not restore after source revocation, even if a
	// previously authorized return token still names the selected Case.
	ownerAllowed=false
	revoked:=postQueueReturn(h,"application/json",string(raw))
	if revoked.Code!=http.StatusNotFound ||
		strings.Contains(revoked.Body.String(),"case-c") ||
		strings.Contains(revoked.Body.String(),"return_token") {
		t.Fatalf("revocation accepted a stale paged return: %d %q",revoked.Code,revoked.Body.String())
	}
}

func TestNativePagedHandoffRejectsOutOfPageOrChangedProjection(t *testing.T) {
	a,b,c:=nativeCase("case-a","Alpha"),nativeCase("case-b","Bravo"),nativeCase("case-c","Charlie")
	change:=false
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {return nil,QueueScope{},errors.New("invalid session")}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		altered:=b
		if change {altered.Title="Bravo updated, cursor stale"}
		return []CaseSnapshot{a,altered,c},nil
	})
	access:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "allow",nil})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,access)
	if err!=nil {t.Fatal(err)}
	_,cursor:=decodeCasePage(t,casePageRequest(h,"/case-queue?sort=title:asc&limit=2"))
	// Case C is not visible on the first page: first-page selection must be
	// rejected despite a permitted Case existing in the full source graph.
	firstOpen:=getQueue(h,"/case-queue?limit=2&sort=title:asc&open_case=case-c")
	if firstOpen.Code!=http.StatusNotFound || strings.Contains(firstOpen.Body.String(),"return_token") {
		t.Fatal("first page was able to open a Case from another page")
	}
	// Case A is not in second page. Cursor does not grant authority to
	// select a different page's Case.
	secondOpen:=getQueue(h,"/case-queue?limit=2&sort=title:asc&cursor="+url.QueryEscape(cursor)+"&open_case=case-a")
	if secondOpen.Code!=http.StatusNotFound || strings.Contains(secondOpen.Body.String(),"return_token") {
		t.Fatal("second page selected first-page Case")
	}
	good:=getQueue(h,"/case-queue?limit=2&sort=title:asc&cursor="+url.QueryEscape(cursor)+"&open_case=case-c")
	if good.Code!=http.StatusOK {t.Fatalf("valid Case C handoff rejected: %d %q",good.Code,good.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(good.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	raw,_:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	change=true
	denied:=postQueueReturn(h,"application/json",string(raw))
	if denied.Code!=http.StatusNotFound || strings.Contains(denied.Body.String(),"case-c") {
		t.Fatalf("stale Case Queue page restored after projection changed: %d %q",denied.Code,denied.Body.String())
	}
}

func TestNativePagedHandoffBadControlsDeniedBeforeSession(t *testing.T) {
	calls:=0
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		calls++;return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("invalid page navigation reached owner");return nil,nil
	})
	auth:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("invalid page navigation reached ACL");return "",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,auth)
	if err!=nil {t.Fatal(err)}
	for _,q:=range []string{
		"open_case=case-a&cursor=arbitrary",
		"open_case=case-a&limit=0",
		"open_case=case-a&limit=101",
		"open_case=case-a&limit=2&cursor=",
		"open_case=case-a&limit=2&cursor=%00",
		"open_case=case-a&limit=2&tenant=tenant-a",
		"open_case=case-a&limit=2&limit=3",
		"open_case=case-a&limit=2&cursor="+strings.Repeat("a",257),
	} {
		resp:=getQueue(h,"/case-queue?"+q)
		if resp.Code!=http.StatusBadRequest || strings.Contains(resp.Body.String(),"return_token") {
			t.Fatalf("invalid paged selection reached owner: %q %d",q,resp.Code)
		}
	}
	if calls!=0 {t.Fatalf("invalid controls reached the session resolver %d times",calls)}
}
