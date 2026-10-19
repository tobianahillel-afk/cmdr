package casequeue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func postQueueReturn(h http.Handler, contentType, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/case-queue/return", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer operator")
	req.Header.Set("Content-Type", contentType)
	out := httptest.NewRecorder()
	h.ServeHTTP(out, req)
	return out
}

func TestNativeQueueHTTPFullReturnRoundTripReauthorizesCurrentCase(t *testing.T) {
	visible := nativeCase(" case-opaque ", "PowerShell investigation")
	visible.Status, visible.Owner = "review", "analyst-a"
	denied := nativeCase("case-denied", "PowerShell forbidden")
	denied.Status, denied.Owner = "review", "analyst-a"
	permitted := true
	sessionCalls, caseChecks := 0, 0
	sessions := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		sessionCalls++
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("unauthenticated")
		}
		return r.Context(), QueueScope{Tenant:"tenant-a",Environment:"prod"}, nil
	})
	provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot, error) {
		return []CaseSnapshot{visible, denied}, nil
	})
	access := queueAccessFunc(func(_ context.Context, _ QueueScope, kind, id string) (string,error) {
		if kind != "case" {t.Fatalf("unexpected object authority %s",kind)}
		caseChecks++
		if id=="case-denied" || (id==visible.ID && !permitted) {return "deny",nil}
		return "allow",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,access)
	if err!=nil {t.Fatal(err)}
	query:=url.Values{
		"q":{"PowerShell"}, "status":{"review"}, "owner":{"analyst-a"},
		"freshness":{"available"}, "sort":{"title:asc"},
		"open_case":{visible.ID}, "scroll":{"127"},
	}
	opened:=getQueue(h,"/case-queue?"+query.Encode())
	if opened.Code!=http.StatusOK {t.Fatalf("first party open HTTP %d: %s",opened.Code,opened.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(opened.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	if handoff.ReturnToken=="" {t.Fatal("Case owner did not issue return state")}
	input,err:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	if err!=nil {t.Fatal(err)}
	returned:=postQueueReturn(h,"application/json",string(input))
	if returned.Code!=http.StatusOK {t.Fatalf("return reauthorization HTTP %d: %s",returned.Code,returned.Body.String())}
	var restored ReturnState
	if err:=json.Unmarshal(returned.Body.Bytes(),&restored);err!=nil {t.Fatal(err)}
	if restored.SelectedCaseID!=visible.ID || restored.ScrollOffset!=127 ||
		restored.SearchQuery!="PowerShell" || len(restored.Filters)!=3 ||
		len(restored.Sort)!=1 || restored.Sort[0].Field!="title" ||
		restored.TenantRef!="tenant-a" || restored.EnvironmentRef!="prod" {
		t.Fatalf("native return did not restore selection/filter/sort/scroll: %+v",restored)
	}
	if strings.Contains(returned.Body.String(),"PowerShell forbidden") ||
		strings.Contains(returned.Body.String(),"case-denied") ||
		strings.Contains(returned.Body.String(),handoff.ReturnToken) {
		t.Fatalf("return leaked denied Case or token: %s",returned.Body.String())
	}
	if caseChecks!=4 || sessionCalls!=2 {t.Fatalf("both requests must freshly authorize all candidates: checks=%d sessions=%d",caseChecks,sessionCalls)}
	if returned.Header().Get("Cache-Control")!="no-store" ||
		returned.Header().Get("Referrer-Policy")!="no-referrer" {
		t.Fatal("returned Case navigation is cacheable/referrer-leaky")
	}
	// Revocation between open and return never silently grants access.
	permitted=false
	revoked:=postQueueReturn(h,"application/json",string(input))
	if revoked.Code!=http.StatusNotFound || strings.Contains(revoked.Body.String(),visible.ID) ||
		strings.Contains(revoked.Body.String(),"return_token") {
		t.Fatalf("revoked owner could restore Case navigation: %d %q",revoked.Code,revoked.Body.String())
	}
}

func TestNativeQueueHTTPReturnDeniesChangedTenantAndStaleCase(t *testing.T) {
	snapshot:=nativeCase("case-a","Investigate endpoint")
	currentTenant:="tenant-a"
	stale:=false
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {return nil,QueueScope{},errors.New("denied")}
		return r.Context(),QueueScope{Tenant:currentTenant,Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		s:=snapshot
		s.TenantRef=currentTenant // independent provider offers an in-scope row on new tenant
		if stale {s.Freshness="stale"}
		return []CaseSnapshot{s},nil
	})
	access:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "allow",nil})
	h,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,access)
	if err!=nil {t.Fatal(err)}
	first:=getQueue(h,"/case-queue?open_case=case-a&scroll=3")
	if first.Code!=http.StatusOK {t.Fatalf("open %d %s",first.Code,first.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(first.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	payload,_:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	currentTenant="another-tenant"
	second:=postQueueReturn(h,"application/json",string(payload))
	if second.Code!=http.StatusNotFound || strings.Contains(second.Body.String(),"case-a") {
		t.Fatalf("old-tenant token was accepted in a new session: %d %q",second.Code,second.Body.String())
	}
	currentTenant="tenant-a"
	stale=true
	second=postQueueReturn(h,"application/json",string(payload))
	if second.Code!=http.StatusNotFound || strings.Contains(second.Body.String(),"case-a") {
		t.Fatalf("stale Case returned to live queue: %d %q",second.Code,second.Body.String())
	}
}

func TestNativeQueueHTTPReturnRejectsMalformedBodyBeforeProviderOrSession(t *testing.T) {
	called:=0
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		called++;return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("malformed return reached Case owner");return nil,nil
	})
	access:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("malformed return reached authority");return "",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,access)
	if err!=nil {t.Fatal(err)}
	for _,tc:=range []struct{name,contentType,body string}{
		{"empty","application/json",""},
		{"not JSON","application/json","not-json"},
		{"missing token","application/json","{}"},
		{"not string","application/json",`{"return_token":42}`},
		{"duplicate token","application/json",`{"return_token":"one","return_token":"two"}`},
		{"forged tenant","application/json",`{"return_token":"one","tenant":"another"}`},
		{"trailing document","application/json",`{"return_token":"one"} {"return_token":"two"}`},
		{"oversize token","application/json",`{"return_token":"`+strings.Repeat("a",maxReturnTokenBytes*2+1)+`"}`},
		{"wrong media type","text/plain",`{"return_token":"token"}`},
	} {
		t.Run(tc.name,func(t *testing.T){
			out:=postQueueReturn(h,tc.contentType,tc.body)
			if out.Code!=http.StatusBadRequest || strings.Contains(out.Body.String(),"return_token") {
				t.Fatalf("malformed return reached owner %d %q",out.Code,out.Body.String())
			}
		})
	}
	for _,u:=range []string{"/case-queue/return?tenant=tenant-a","/case-queue/return?return_token=token"} {
		req:=httptest.NewRequest(http.MethodPost,u,bytes.NewBufferString(`{"return_token":"opaque"}`))
		req.Header.Set("Content-Type","application/json")
		out:=httptest.NewRecorder()
		h.ServeHTTP(out,req)
		if out.Code!=http.StatusBadRequest {t.Fatalf("return URL query accepted %s",u)}
	}
	req:=httptest.NewRequest(http.MethodPost,"/case-queue",bytes.NewBufferString(`{"return_token":"opaque"}`))
	out:=httptest.NewRecorder()
	h.ServeHTTP(out,req)
	if out.Code!=http.StatusMethodNotAllowed {t.Fatal("POST allowed on normal Queue endpoint")}
	req=httptest.NewRequest(http.MethodGet,"/case-queue/return?return_token=token",nil)
	out=httptest.NewRecorder()
	h.ServeHTTP(out,req)
	if out.Code!=http.StatusMethodNotAllowed {t.Fatal("GET allowed to leak return token into URL")}
	if called!=0 {t.Fatalf("malformed return reached session: %d",called)}
}
