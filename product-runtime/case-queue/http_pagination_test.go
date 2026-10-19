package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func casePageRequest(handler http.Handler, target string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, target, nil)
	r.Header.Set("Authorization", "Bearer operator")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func decodeCasePage(t *testing.T, r *httptest.ResponseRecorder) (Projection, string) {
	t.Helper()
	if r.Code != http.StatusOK {
		t.Fatalf("authorized page unavailable: %d %s", r.Code, r.Body.String())
	}
	var parsed struct {
		Projection
		NextCursor string `json:"next_cursor"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("not a Case Queue projection: %v", err)
	}
	if r.Header().Get("Cache-Control") != "no-store" ||
		r.Header().Get("Referrer-Policy") != "no-referrer" {
		t.Fatal("Case Queue page response has unsafe cache/referrer headers")
	}
	return parsed.Projection, parsed.NextCursor
}

func TestNativeCasePagesReuseOwnerProjectionAndNeverRevealDeniedCases(t *testing.T) {
	a := nativeCase("case-a", "Alpha")
	b := nativeCase("case-b", "Bravo")
	c := nativeCase("case-c", "Charlie")
	secret := nativeCase("private-case", "TOP-SECRET-CASE-TITLE")
	// The client cannot rely on provider-supplied AccessDecision.
	secret.AccessDecision = "allow"
	aclCalls := 0
	session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("unauthorized")
		}
		return r.Context(), QueueScope{Tenant: "tenant-a", Environment: "prod"}, nil
	})
	provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot, error) {
		return []CaseSnapshot{secret, c, a, b}, nil
	})
	access := queueAccessFunc(func(_ context.Context, _ QueueScope, kind, id string) (string, error) {
		aclCalls++
		if kind != "case" {
			t.Fatalf("unexpected owner kind %q", kind)
		}
		if id == "private-case" {
			return "deny", nil
		}
		return "allow", nil
	})
	h, err := NewAuthorizedQueueHTTPHandler(session, provider, access)
	if err != nil { t.Fatal(err) }
	projection, cursor := decodeCasePage(t,
		casePageRequest(h, "/case-queue?limit=2&sort=title:asc"))
	if len(projection.Rows) != 2 ||
		projection.Rows[0].Case.ID != "case-a" ||
		projection.Rows[1].Case.ID != "case-b" ||
		cursor == "" || len(cursor) > 256 {
		t.Fatalf("wrong first authorized page: rows=%+v cursor=%q", projection.Rows, cursor)
	}
	for _, secretValue := range []string{"private-case", "TOP-SECRET", "tenant-b"} {
		if strings.Contains(cursor, secretValue) {
			t.Fatalf("presentation cursor leaked owner data %q", secretValue)
		}
	}
	firstACLCount := aclCalls
	second := casePageRequest(h,
		"/case-queue?sort=title:asc&limit=2&cursor="+url.QueryEscape(cursor))
	after, last := decodeCasePage(t, second)
	if len(after.Rows) != 1 || after.Rows[0].Case.ID != "case-c" || last != "" {
		t.Fatalf("wrong second authorized page: %+v, cursor %q", after.Rows, last)
	}
	if aclCalls-firstACLCount != 4 {
		t.Fatalf("subsequent page skipped full server ACL pass: %d", aclCalls-firstACLCount)
	}
	if strings.Contains(second.Body.String(), "TOP-SECRET") ||
		strings.Contains(second.Body.String(), "private-case") {
		t.Fatal("denied Case content leaked on next page")
	}
	// Existing clients that do not request pagination see the original
	// complete, shape-compatible native Case Queue projection.
	legacy := casePageRequest(h, "/case-queue?sort=title:asc")
	full, _ := decodeCasePage(t, legacy)
	if len(full.Rows) != 3 {
		t.Fatalf("pagination altered the legacy full Case Queue: %d",len(full.Rows))
	}
	if strings.Contains(legacy.Body.String(), "next_cursor") {
		t.Fatal("nonpaginated native Case Queue schema unexpectedly changed")
	}
}

func TestNativeCaseCursorFailsClosedAfterSourceRevocationOrProjectionChange(t *testing.T) {
	a, b, c := nativeCase("case-a","Alpha"), nativeCase("case-b","Bravo"), nativeCase("case-c","Charlie")
	canReadB := true
	revision := "v1"
	session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
		if r.Header.Get("Authorization") != "Bearer operator" {
			return nil, QueueScope{}, errors.New("no session")
		}
		return r.Context(), QueueScope{Tenant:"tenant-a", Environment:"prod"},nil
	})
	provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot,error) {
		return []CaseSnapshot{a,b,c},nil
	})
	access := queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		if kind!="case" { t.Fatal("unexpected owner kind") }
		if id=="case-b" && !canReadB {return "deny",nil}
		return "allow",nil
	})
	views := queueSavedViewFunc(func(_ context.Context,scope QueueScope)(*SavedView,error){
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {t.Fatal("untrusted view scope")}
		return &SavedView{ID:"shared-case-view",TenantRef:"tenant-a",Version:revision,AccessDecision:"allow"},nil
	})
	h,err:=NewAuthorizedQueueHTTPHandlerWithSavedViews(session,provider,access,views)
	if err!=nil { t.Fatal(err) }
	_,cursor:=decodeCasePage(t,casePageRequest(h,"/case-queue?limit=1&sort=title:asc"))
	if cursor=="" {t.Fatal("missing continuation")}
	nextURL:="/case-queue?sort=title:asc&limit=1&cursor="+url.QueryEscape(cursor)
	// The full projection may be unchanged even though Shared version has
	// changed. Cursor MUST still be rejected before exposing any next row.
	revision="v2"
	stale:=casePageRequest(h,nextURL)
	if stale.Code!=http.StatusNotFound ||
		strings.Contains(stale.Body.String(),"case-b") {
		t.Fatalf("stale Shared view cursor accepted: %d %q",stale.Code,stale.Body.String())
	}
	revision="v1"
	canReadB=false
	revoked:=casePageRequest(h,nextURL)
	if revoked.Code!=http.StatusNotFound ||
		strings.Contains(revoked.Body.String(),"case-c") ||
		strings.Contains(revoked.Body.String(),"case-b") {
		t.Fatalf("revoked owner cursor returned data: %d %q",revoked.Code,revoked.Body.String())
	}
	// A new first page only contains the now-authorized Cases.
	current,_:=decodeCasePage(t,casePageRequest(h,"/case-queue?limit=2&sort=title:asc"))
	if len(current.Rows)!=2 || current.Rows[0].Case.ID!="case-a" ||
		current.Rows[1].Case.ID!="case-c" {
		t.Fatalf("fresh owner projection not rebuilt: %+v",current.Rows)
	}
}

func TestNativeCasePagesRejectInvalidOrForgedRequestsBeforeProviders(t *testing.T) {
	session:=queueSessionFunc(func(*http.Request)(context.Context,QueueScope,error){
		t.Fatal("invalid pagination request reached session");return nil,QueueScope{},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("invalid pagination reached Case provider");return nil,nil
	})
	acl:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("invalid pagination reached owner ACL");return "",nil
	})
	h,err:=NewAuthorizedQueueHTTPHandler(session,provider,acl)
	if err!=nil {t.Fatal(err)}
	for _,path:=range []string{
		"/case-queue?limit=0", "/case-queue?limit=101",
		"/case-queue?limit=-1", "/case-queue?limit=NaN",
		"/case-queue?cursor=abc", "/case-queue?limit=1&limit=2",
		"/case-queue?limit=1&cursor=", "/case-queue?limit=1&cursor=%00",
		"/case-queue?cursor=abc&open_case=case-a",
		"/case-queue?limit=1&open_case=case-a&scroll=-1",
		"/case-queue?limit=1&sort=not_allowed:asc",
		"/case-queue?limit=1&tenant=tenant-b",
	} {
		if got:=casePageRequest(h,path);got.Code!=http.StatusBadRequest {
			t.Fatalf("malformed paging accepted %q: %d %q",path,got.Code,got.Body.String())
		}
	}
}

func TestNativeCaseCursorCannotSelectUnscopedOrForeignRows(t *testing.T) {
	a,b := nativeCase("case-a","Alpha"),nativeCase("case-b","Bravo")
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		return []CaseSnapshot{a,b},nil
	})
	auth:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "allow",nil})
	h,_:=NewAuthorizedQueueHTTPHandler(session,provider,auth)
	_,firstCursor:=decodeCasePage(t,casePageRequest(h,"/case-queue?limit=1"))
	for _,value:=range []string{
		"not-base64", firstCursor+"invalid", strings.Repeat("a",257),
	} {
		got:=casePageRequest(h,"/case-queue?limit=1&cursor="+url.QueryEscape(value))
		if got.Code!=http.StatusNotFound && got.Code!=http.StatusBadRequest {
			t.Fatalf("forged cursor accepted: %d %q",got.Code,got.Body.String())
		}
		if strings.Contains(got.Body.String(),"case-b") {
			t.Fatal("forged cursor disclosed authorized row without valid page state")
		}
	}
	// Changing even only a visible Case's freshness/metadata invalidates
	// the cursor, rather than silently delivering a possibly inconsistent page.
	a.Freshness="stale"
	got:=casePageRequest(h,"/case-queue?limit=1&cursor="+url.QueryEscape(firstCursor))
	if got.Code!=http.StatusNotFound {
		t.Fatalf("changed full projection accepted prior cursor: %d",got.Code)
	}
}
