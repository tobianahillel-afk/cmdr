package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type queueSavedViewFunc func(context.Context, QueueScope) (*SavedView, error)
func (f queueSavedViewFunc) CurrentQueueSavedView(ctx context.Context, scope QueueScope) (*SavedView, error) {
	return f(ctx, scope)
}

func TestNativeHTTPSharedSavedViewAuthorizedQueueAndReturn(t *testing.T) {
	first:=nativeCase("case-review","Review requested")
	first.Status="review"
	second:=nativeCase("case-open","Secret queued")
	second.Status="investigating"
	currentVersion:="v1"
	viewAllowed:=true
	resolverCalls:=0
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {return nil,QueueScope{},errors.New("denied")}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	cases:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error) {
		return []CaseSnapshot{first,second},nil
	})
	authorizer:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		if kind!="case" { t.Fatalf("unexpected object owner %q",kind) }
		return "allow",nil
	})
	views:=queueSavedViewFunc(func(_ context.Context,scope QueueScope)(*SavedView,error){
		resolverCalls++
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {t.Fatal("view resolver used client-owned scope")}
		decision:="allow"
		if !viewAllowed {decision="deny"}
		return &SavedView{
			ID:"queue-analyst-review",TenantRef:scope.Tenant,
			Version:currentVersion,AccessDecision:decision,
			Fields:[]SavedViewField{
				{Name:"status",AccessDecision:"allow"},
				{Name:"FORBIDDEN-SOURCE-FIELD",AccessDecision:"deny"},
			},
			Filters:[]SavedViewFilter{
				{Field:"status",Values:[]string{"review"},AccessDecision:"allow"},
				{Field:"owner",Values:[]string{"SECRET-DENIED-OWNER"},AccessDecision:"deny"},
			},
			Sort:[]SavedViewSort{{Field:"title",Direction:"asc",AccessDecision:"allow"}},
		},nil
	})
	handler,err:=NewAuthorizedQueueHTTPHandlerWithSavedViews(sessions,cases,authorizer,views)
	if err!=nil {t.Fatal(err)}
	list:=getQueue(handler,"/case-queue")
	if list.Code!=http.StatusOK {t.Fatalf("Shared view GET failed %d %q",list.Code,list.Body.String())}
	var projected Projection
	if err:=json.Unmarshal(list.Body.Bytes(),&projected);err!=nil {t.Fatal(err)}
	if len(projected.Rows)!=1 || projected.Rows[0].Case.ID!=first.ID ||
		len(projected.AppliedSavedViewFields)!=1 || projected.AppliedSavedViewFields[0]!="status" ||
		projected.PermissionFilteredElements!=2 {
		t.Fatalf("Shared view source/field authorization not applied: %+v",projected)
	}
	for _,secret:=range []string{"FORBIDDEN-SOURCE-FIELD","SECRET-DENIED-OWNER","case-open","Secret queued"} {
		if strings.Contains(list.Body.String(),secret) {
			t.Fatalf("private Saved View field or filtered Case leaked: %q",secret)
		}
	}
	open:=getQueue(handler,"/case-queue?open_case=case-review&scroll=19")
	if open.Code!=http.StatusOK {t.Fatalf("open with Shared view failed: %d %s",open.Code,open.Body.String())}
	var handoff WorkspaceHandoff
	if err:=json.Unmarshal(open.Body.Bytes(),&handoff);err!=nil {t.Fatal(err)}
	if handoff.ReturnToken=="" {t.Fatal("missing first-party Queue return state")}
	payload,_:=json.Marshal(map[string]string{"return_token":handoff.ReturnToken})
	postQueue:=func() *httptest.ResponseRecorder {
		return postQueueReturn(handler,"application/json",string(payload))
	}
	ok:=postQueue()
	if ok.Code!=http.StatusOK {t.Fatalf("same Shared view return failed: %d %s",ok.Code,ok.Body.String())}
	var restored ReturnState
	if err:=json.Unmarshal(ok.Body.Bytes(),&restored);err!=nil {t.Fatal(err)}
	if restored.SavedViewID!="queue-analyst-review" ||
		restored.SavedViewVersion!="v1" || restored.SelectedCaseID!="case-review" ||
		restored.ScrollOffset!=19 {
		t.Fatalf("return did not preserve current Shared view and scroll: %+v",restored)
	}
	if resolverCalls!=3 {t.Fatalf("GET, open and return must individually re-resolve Shared view: %d",resolverCalls)}
	currentVersion="v2"
	stale:=postQueue()
	if stale.Code!=http.StatusNotFound || strings.Contains(stale.Body.String(),first.ID) {
		t.Fatalf("changed Shared view version authorized stale token: %d %q",stale.Code,stale.Body.String())
	}
	viewAllowed=false
	denied:=postQueue()
	if denied.Code!=http.StatusNotFound || strings.Contains(denied.Body.String(),first.ID) {
		t.Fatalf("revoked Shared view authorized old token: %d %q",denied.Code,denied.Body.String())
	}
	// A denied current Shared view safely falls back to the generic default,
	// not the formerly saved restricted view and never exposes its fields.
	fallback:=getQueue(handler,"/case-queue")
	if fallback.Code!=http.StatusOK {t.Fatalf("denied Shared view fallback fails: %d",fallback.Code)}
	var defaultProjection Projection
	if err:=json.Unmarshal(fallback.Body.Bytes(),&defaultProjection);err!=nil {t.Fatal(err)}
	if len(defaultProjection.Rows)!=2 ||
		!strings.Contains(fallback.Body.String(),"saved-view:fallback") ||
		strings.Contains(fallback.Body.String(),"FORBIDDEN-SOURCE-FIELD") {
		t.Fatalf("denied Shared view fallback leaked or reused view: %+v",defaultProjection)
	}
}

func TestNativeHTTPSharedSavedViewRejectsForeignOrUnavailableSources(t *testing.T) {
	sessionCalled,caseCalled,viewCalled:=0,0,0
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		sessionCalled++
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	cases:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		caseCalled++
		return []CaseSnapshot{nativeCase("case-1","Sensitive Case")},nil
	})
	auth:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "allow",nil})
	mode:="foreign"
	views:=queueSavedViewFunc(func(_ context.Context,scope QueueScope)(*SavedView,error){
		viewCalled++
		switch mode {
		case "foreign": return &SavedView{ID:"foreign-view",TenantRef:"tenant-foreign",Version:"v1",AccessDecision:"allow"},nil
		case "error": return nil,errors.New("Shared unavailable")
		case "invalid": return &SavedView{ID:"current",TenantRef:scope.Tenant,Version:"v1",AccessDecision:"not-approved"},nil
		}
		return nil,nil
	})
	if _,err:=NewAuthorizedQueueHTTPHandlerWithSavedViews(sessions,cases,auth,nil);err==nil {
		t.Fatal("nil Shared resolver unexpectedly accepted")
	}
	h,err:=NewAuthorizedQueueHTTPHandlerWithSavedViews(sessions,cases,auth,views)
	if err!=nil {t.Fatal(err)}
	for _,m:=range []string{"foreign","error","invalid"} {
		mode=m
		got:=getQueue(h,"/case-queue")
		if got.Code!=http.StatusNotFound ||
			strings.Contains(got.Body.String(),"case-1") ||
			strings.Contains(got.Body.String(),"Sensitive Case") ||
			strings.Contains(got.Body.String(),"foreign-view") {
			t.Fatalf("bad Shared owner response leaked Case/source in mode %s: %d %q",m,got.Code,got.Body.String())
		}
	}
	// Untrusted client-supplied scope and view ID are still rejected before
	// opening sessions/providers/Shared, even when host has a Shared provider.
	prevSessions,prevCases,prevViews:=sessionCalled,caseCalled,viewCalled
	got:=getQueue(h,"/case-queue?view_id=foreign-view")
	if got.Code!=http.StatusBadRequest || prevSessions!=sessionCalled ||
		prevCases!=caseCalled || prevViews!=viewCalled {
		t.Fatal("client-supplied Shared view ID reached server-side owners")
	}
	// Legacy no-Shared hosts still show a normal authorized default Queue.
	legacy,err:=NewAuthorizedQueueHTTPHandler(sessions,cases,auth)
	if err!=nil {t.Fatal(err)}
	plain:=getQueue(legacy,"/case-queue")
	if plain.Code!=http.StatusOK || !strings.Contains(plain.Body.String(),"case-1") {
		t.Fatalf("Shared extension broke backward compatibility: %d %q",plain.Code,plain.Body.String())
	}
}
