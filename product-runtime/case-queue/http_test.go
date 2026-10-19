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

type queueSessionFunc func(*http.Request) (context.Context, QueueScope, error)
func (f queueSessionFunc) ResolveQueueSession(r *http.Request) (context.Context, QueueScope, error) { return f(r) }

type queueProviderFunc func(context.Context, QueueScope) ([]CaseSnapshot, error)
func (f queueProviderFunc) ListQueueCandidates(ctx context.Context, scope QueueScope) ([]CaseSnapshot, error) { return f(ctx,scope) }

type queueAccessFunc func(context.Context, QueueScope, string, string) (string, error)
func (f queueAccessFunc) CheckQueueRead(ctx context.Context, scope QueueScope, kind, id string) (string,error) { return f(ctx,scope,kind,id) }

func nativeCase(id, title string) CaseSnapshot {
	return CaseSnapshot{ID:id,TenantRef:"tenant-a",HumanID:id,Title:title,
		Status:"investigating",Owner:"analyst-a",UpdatedAt:"2026-10-19T10:00:00Z",
		Freshness:"available",AccessDecision:"allow"}
}

func getQueue(h http.Handler, url string) *httptest.ResponseRecorder {
	req:=httptest.NewRequest(http.MethodGet,url,nil)
	req.Header.Set("Authorization","Bearer operator")
	got:=httptest.NewRecorder()
	h.ServeHTTP(got,req)
	return got
}

func TestNativeQueueHTTPAuthorizesEveryObjectBeforeFiltering(t *testing.T) {
	first:=nativeCase("case-visible","PowerShell investigation")
	first.AccessDecision="deny" // provider decisions are NOT authorities
	first.Incident=&IncidentContext{
		ID:"incident-secret",TenantRef:"tenant-a",AccessDecision:"allow",
		State:"available",Priority:"P0",Owner:"internal-secret-owner",
	}
	first.Findings=[]FindingContext{
		{ID:"finding-visible",TenantRef:"tenant-a",CaseID:first.ID,AccessDecision:"deny",
		 State:"available",Status:"suspicious"},
		{ID:"finding-secret",TenantRef:"tenant-a",CaseID:first.ID,AccessDecision:"allow",
		 State:"available",Status:"forbidden-secret"},
	}
	second:=nativeCase("case-secret","TOP-SECRET-CASE-TITLE")
	seen:=[]string{}
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error) {
		if r.Header.Get("Authorization")!="Bearer operator" {
			return nil,QueueScope{},errors.New("unauthorized")
		}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(ctx context.Context,scope QueueScope)([]CaseSnapshot,error){
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {t.Fatal("scope was derived from query")}
		return []CaseSnapshot{first,second},nil
	})
	auth:=queueAccessFunc(func(_ context.Context,scope QueueScope,kind,id string)(string,error){
		seen=append(seen,kind+":"+id)
		if scope.Tenant!="tenant-a" { t.Fatal("cross-scope authorization") }
		if id=="case-secret"||id=="incident-secret"||id=="finding-secret" {return "deny",nil}
		return "allow",nil
	})
	handler,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,auth)
	if err!=nil {t.Fatal(err)}
	response:=getQueue(handler,"/case-queue?q=PowerShell")
	if response.Code!=http.StatusOK {t.Fatalf("authorized read rejected: %d %q",response.Code,response.Body.String())}
	var projected Projection
	if err:=json.Unmarshal(response.Body.Bytes(),&projected);err!=nil {t.Fatal(err)}
	if len(projected.Rows)!=1 || projected.Rows[0].Case.ID!="case-visible" ||
		projected.Rows[0].Case.Incident!=nil ||
		len(projected.Rows[0].Case.Findings)!=1 ||
		projected.Rows[0].Case.Findings[0].ID!="finding-visible" {
		t.Fatalf("native queue leaked owner-denied projections or did not search correctly: %+v",projected)
	}
	for _,secret:=range []string{"case-secret","TOP-SECRET","incident-secret","internal-secret-owner","finding-secret","forbidden-secret"} {
		if strings.Contains(response.Body.String(),secret) {t.Fatalf("denied source metadata leaked %q",secret)}
	}
	if projected.PermissionFilteredCases!=0 {
		t.Fatal("denied candidate Case count disclosed to another analyst")
	}
	for _,wanted:=range []string{"case:case-visible","case:case-secret","incident:incident-secret","finding:finding-visible","finding:finding-secret"} {
		found:=false
		for _,s:=range seen { if s==wanted {found=true} }
		if !found {t.Fatalf("missing independent owner authorization of %s: %v",wanted,seen)}
	}
	if response.Header().Get("Cache-Control")!="no-store" ||
		response.Header().Get("Vary")!="Authorization" {
		t.Fatal("sensitive Queue projection is cacheable")
	}
	// No filtering or search is performed over inaccessible candidate Cases.
	noMatch:=getQueue(handler,"/case-queue?q=TOP-SECRET")
	if noMatch.Code!=http.StatusOK || strings.Contains(noMatch.Body.String(),"case-secret") {
		t.Fatalf("search revealed denied Case: %d %q",noMatch.Code,noMatch.Body.String())
	}
	var empty Projection
	if err:=json.Unmarshal(noMatch.Body.Bytes(),&empty);err!=nil {t.Fatal(err)}
	if len(empty.Rows)!=0 {t.Fatal("search across denied Cases returned a row")}
}

func TestNativeQueueHTTPFailsClosedOnUntrustedScopeAndOwnerFailure(t *testing.T) {
	valid:=nativeCase("case-ok","Valid title")
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	for _,tc:=range []struct{
		name string
		data []CaseSnapshot
		decision string
		authError error
		code int
	}{
		{"foreign provider Case",[]CaseSnapshot{valid,nativeCase("foreign","Secret")}, "allow",nil,404},
		{"invalid owner decision",[]CaseSnapshot{valid},"unknown-grant",nil,404},
		{"owner outage",[]CaseSnapshot{valid},"allow",errors.New("unavailable"),404},
		{"missing object",[]CaseSnapshot{valid},"missing",nil,200},
		{"empty Case ID",[]CaseSnapshot{nativeCase("","Invalid")},"allow",nil,404},
		{"valid",[]CaseSnapshot{valid},"allow",nil,200},
	} {
		t.Run(tc.name,func(t *testing.T) {
			if tc.name=="foreign provider Case" {tc.data[1].TenantRef="tenant-b"}
			p:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){return tc.data,nil})
			a:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
				return tc.decision,tc.authError
			})
			h,err:=NewAuthorizedQueueHTTPHandler(session,p,a)
			if err!=nil {t.Fatal(err)}
			got:=getQueue(h,"/case-queue")
			if got.Code!=tc.code {t.Fatalf("HTTP %d want %d: %s",got.Code,tc.code,got.Body.String())}
			if strings.Contains(got.Body.String(),"Secret") || strings.Contains(got.Body.String(),"Valid title") && tc.code!=200 {
				t.Fatal("failed query returned candidate data")
			}
		})
	}
	for _,raw:=range []string{
		"/case-queue?tenant=other", "/case-queue?q=one&q=two",
		"/case-queue?q=%00", "/case-queue?q=%0A", "/case-queue?q=abc&access_decision=allow",
		"/case-queue?q="+strings.Repeat("x",257),
	} {
		p:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){t.Fatal("untrusted request reached provider");return nil,nil})
		a:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){t.Fatal("untrusted request reached auth");return "",nil})
		h,_:=NewAuthorizedQueueHTTPHandler(session,p,a)
		if got:=getQueue(h,raw);got.Code!=http.StatusBadRequest {t.Fatalf("unsafe query %q -> %d",raw,got.Code)}
	}
	authBad:=queueSessionFunc(func(*http.Request)(context.Context,QueueScope,error){
		return nil,QueueScope{},errors.New("deny")
	})
	p:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){t.Fatal("unauthenticated provider access");return nil,nil})
	a:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){t.Fatal("unauthenticated auth");return "",nil})
	h,_:=NewAuthorizedQueueHTTPHandler(authBad,p,a)
	if got:=getQueue(h,"/case-queue");got.Code!=http.StatusUnauthorized {t.Fatal("unauthenticated read allowed")}
}

func TestNativeQueueHTTPRejectsUnrelatedMethodsAndMissingAuthorities(t *testing.T) {
	if _,err:=NewAuthorizedQueueHTTPHandler(nil,nil,nil);err==nil {t.Fatal("empty authorities accepted")}
	handler,_:=NewAuthorizedQueueHTTPHandler(
		queueSessionFunc(func(*http.Request)(context.Context,QueueScope,error){return context.Background(),QueueScope{Tenant:"t",Environment:"e"},nil}),
		queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){return nil,nil}),
		queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "deny",nil}),
	)
	for _,method:=range []string{http.MethodPost,http.MethodPut,http.MethodDelete} {
		req:=httptest.NewRequest(method,"/case-queue",strings.NewReader("{\"assign\":\"analyst\"}"))
		got:=httptest.NewRecorder()
		handler.ServeHTTP(got,req)
		if got.Code!=http.StatusMethodNotAllowed || got.Header().Get("Allow")!="GET" {
			t.Fatalf("read-only native Queue accepted %s: %d",method,got.Code)
		}
	}
}

func TestNativeQueueHTTPRejectsPartialOrCrossScopedContextAtomically(t *testing.T) {
	base:=nativeCase("case-a","Permitted Case")
	base.Incident=&IncidentContext{ID:"incident-a",TenantRef:"tenant-a",
		AccessDecision:"allow",State:"available",Priority:"P2"}
	base.Findings=[]FindingContext{{ID:"finding-a",CaseID:"case-a",TenantRef:"tenant-a",
		AccessDecision:"allow",State:"available",Status:"review"}}
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	for _,tc:=range []struct {
		name string
		change func(*CaseSnapshot)
		denyKind string
		ownerErr bool
		expect int
	}{
		{"allow independently authorized Incident and Finding",nil,"",false,200},
		{"Incident belongs to foreign tenant",func(c *CaseSnapshot){c.Incident.TenantRef="tenant-b"},"",false,404},
		{"Incident identifier missing",func(c *CaseSnapshot){c.Incident.ID=""},"",false,404},
		{"Incident permission infrastructure failed",nil,"incident",true,404},
		{"Finding belongs to foreign tenant",func(c *CaseSnapshot){c.Findings[0].TenantRef="tenant-b"},"",false,404},
		{"Finding linked to another Case",func(c *CaseSnapshot){c.Findings[0].CaseID="case-b"},"",false,404},
		{"Finding identifier missing",func(c *CaseSnapshot){c.Findings[0].ID=""},"",false,404},
		{"Finding permission infrastructure failed",nil,"finding",true,404},
		{"Failed native projection returns generic deny",func(c *CaseSnapshot){c.HumanID=""},"",false,404},
	} {
		t.Run(tc.name,func(t *testing.T){
			c:=base
			inc:=*base.Incident
			c.Incident=&inc
			c.Findings=append([]FindingContext(nil),base.Findings...)
			if tc.change!=nil {tc.change(&c)}
			p:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
				return []CaseSnapshot{c},nil
			})
			a:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
				if kind==tc.denyKind && tc.ownerErr {
					return "",errors.New("owner read outage")
				}
				return "allow",nil
			})
			h,err:=NewAuthorizedQueueHTTPHandler(session,p,a)
			if err!=nil {t.Fatal(err)}
			got:=getQueue(h,"/case-queue")
			if got.Code!=tc.expect {t.Fatalf("HTTP %d expected %d: %q",got.Code,tc.expect,got.Body.String())}
			if tc.expect==200 {
				var projected Projection
				if err:=json.Unmarshal(got.Body.Bytes(),&projected);err!=nil {t.Fatal(err)}
				if len(projected.Rows)!=1 || projected.Rows[0].Case.Incident==nil ||
					len(projected.Rows[0].Case.Findings)!=1 {
					t.Fatalf("authorized nested Case context missing: %+v",projected)
				}
			} else if strings.Contains(got.Body.String(),"case-a") ||
				strings.Contains(got.Body.String(),"incident-a") ||
				strings.Contains(got.Body.String(),"finding-a") {
				t.Fatalf("incomplete authorization leaked partial Case: %q",got.Body.String())
			}
		})
	}
}

func TestNativeQueueHTTPFailsClosedOnProviderLimitsAndCancellation(t *testing.T) {
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	auth:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){return "allow",nil})
	for _,tc:=range []struct {
		name string
		provider queueProviderFunc
		want int
	}{
		{"provider failure",func(context.Context,QueueScope)([]CaseSnapshot,error){return nil,errors.New("db unavailable")},404},
		{"unbounded candidate response",func(context.Context,QueueScope)([]CaseSnapshot,error){return make([]CaseSnapshot,5001),nil},404},
		{"empty queue",func(context.Context,QueueScope)([]CaseSnapshot,error){return nil,nil},200},
		{"response limit",func(context.Context,QueueScope)([]CaseSnapshot,error){
			c:=nativeCase("oversized","Very large Case")
			c.Title=strings.Repeat("private",650000)
			return []CaseSnapshot{c},nil
		},503},
	} {
		t.Run(tc.name,func(t *testing.T) {
			h,_:=NewAuthorizedQueueHTTPHandler(session,tc.provider,auth)
			got:=getQueue(h,"/case-queue")
			if got.Code!=tc.want {t.Fatalf("HTTP %d want %d",got.Code,tc.want)}
			if tc.want!=200 && (strings.Contains(got.Body.String(),"private") ||
				strings.Contains(got.Body.String(),"db unavailable")) {
				t.Fatal("provider failure leaked internal or oversized data")
			}
		})
	}
	// The request may be cancelled before the session is even resolved;
	// neither the Case provider nor any permission check may execute.
	p:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){t.Fatal("cancelled request reached provider");return nil,nil})
	h,_:=NewAuthorizedQueueHTTPHandler(session,p,auth)
	req:=httptest.NewRequest(http.MethodGet,"/case-queue",nil)
	ctx,cancel:=context.WithCancel(req.Context())
	cancel()
	got:=httptest.NewRecorder()
	h.ServeHTTP(got,req.WithContext(ctx))
	if got.Code!=http.StatusServiceUnavailable {t.Fatalf("cancelled read returned HTTP %d",got.Code)}
}

func TestNativeQueueHTTPFiltersAndSortOnlyAuthorizedCases(t *testing.T) {
	first:=nativeCase("c-20","Zeta response case")
	first.Status="review"
	first.Owner="analyst-a"
	first.UpdatedAt="2026-10-19T08:00:00Z"
	second:=nativeCase("c-10","Alpha response case")
	second.Status="review"
	second.Owner="analyst-a"
	second.UpdatedAt="2026-10-19T09:00:00Z"
	hidden:=nativeCase("c-secret","TOP-SECRET response case")
	hidden.Status="review"
	hidden.Owner="analyst-a"
	other:=nativeCase("c-30","Unrelated case")
	other.Status="closed"
	other.Owner="analyst-b"
	session:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		if r.Header.Get("Authorization")!="Bearer operator" {
			return nil,QueueScope{},errors.New("no session")
		}
		return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
	})
	provider:=queueProviderFunc(func(_ context.Context,scope QueueScope)([]CaseSnapshot,error){
		if scope.Tenant!="tenant-a" || scope.Environment!="prod" {
			t.Fatal("query changed authenticated scope")
		}
		return []CaseSnapshot{first,second,hidden,other},nil
	})
	seen:=map[string]int{}
	authorizer:=queueAccessFunc(func(_ context.Context,_ QueueScope,kind,id string)(string,error){
		seen[kind+":"+id]++
		if id=="c-secret" { return "deny",nil }
		return "allow",nil
	})
	handler,err:=NewAuthorizedQueueHTTPHandler(session,provider,authorizer)
	if err!=nil {t.Fatal(err)}
	resp:=getQueue(handler,"/case-queue?q=response&status=review&owner=analyst-a&freshness=available&sort=title:asc")
	if resp.Code!=http.StatusOK {t.Fatalf("filtered authorized queue: HTTP %d %q",resp.Code,resp.Body.String())}
	var result Projection
	if err:=json.Unmarshal(resp.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
	if len(result.Rows)!=2 || result.Rows[0].Case.ID!="c-10" || result.Rows[1].Case.ID!="c-20" {
		t.Fatalf("native queue did not filter and stably sort allowed cases: %+v",result.Rows)
	}
	if len(result.AppliedSort)<1 || result.AppliedSort[0].Field!="title" ||
		result.AppliedSort[0].Direction!="asc" {
		t.Fatalf("native sort was not applied: %+v",result.AppliedSort)
	}
	if result.PermissionFilteredCases!=0 || strings.Contains(resp.Body.String(),"TOP-SECRET") ||
		strings.Contains(resp.Body.String(),"c-secret") {
		t.Fatal("client filter leaked denied Case or its count")
	}
	for _,id:=range []string{"c-10","c-20","c-30","c-secret"} {
		if seen["case:"+id]!=1 {t.Fatalf("candidate %s was not independently authorized before filtering",id)}
	}
	// The same native projection supports descending updated-at and does
	// not mutate any Case or save a Shared Saved View.
	recent:=getQueue(handler,"/case-queue?status=review&sort=updated_at:desc")
	if recent.Code!=http.StatusOK {t.Fatalf("native sort unavailable: %d",recent.Code)}
	var sorted Projection
	if err:=json.Unmarshal(recent.Body.Bytes(),&sorted);err!=nil {t.Fatal(err)}
	if len(sorted.Rows)!=2 || sorted.Rows[0].Case.ID!="c-10" ||
		sorted.Rows[1].Case.ID!="c-20" {
		t.Fatalf("sort by Case update freshness failed: %+v",sorted.Rows)
	}
}

func TestNativeQueueHTTPRejectsUnsafeFilterSortBeforeAuthorityCalls(t *testing.T) {
	sessions:=queueSessionFunc(func(r *http.Request)(context.Context,QueueScope,error){
		t.Fatal("invalid query reached session resolver")
		return nil,QueueScope{},errors.New("not expected")
	})
	provider:=queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error){
		t.Fatal("invalid query reached provider")
		return nil,nil
	})
	access:=queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error){
		t.Fatal("invalid query reached authorizer")
		return "",nil
	})
	handler,err:=NewAuthorizedQueueHTTPHandler(sessions,provider,access)
	if err!=nil {t.Fatal(err)}
	for _,path:=range []string{
		"/case-queue?tenant=tenant-b&status=review",
		"/case-queue?status=review&status=closed",
		"/case-queue?owner=",
		"/case-queue?status=%00",
		"/case-queue?owner="+strings.Repeat("x",129),
		"/case-queue?sort=tenant:asc",
		"/case-queue?sort=title:up",
		"/case-queue?sort=title",
		"/case-queue?sort=title:asc&sort=title:desc",
		"/case-queue?status=review&owner=analyst-a&freshness=available&sort=title:asc&q=x&extra=unsafe",
	} {
		got:=getQueue(handler,path)
		if got.Code!=http.StatusBadRequest || strings.Contains(got.Body.String(),"tenant-b") {
			t.Fatalf("unsafe queue query %q accepted: %d %q",path,got.Code,got.Body.String())
		}
	}
}
