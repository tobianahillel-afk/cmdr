package workqueue

import (
	"context"
	"errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fixedHostSession struct {
	scope SessionScope
	err error
	calls int
}
func (s *fixedHostSession) ResolveSession(_ context.Context, _ *http.Request) (SessionScope,error) {
	s.calls++
	return s.scope,s.err
}
func trustedHost() *fixedHostSession {
	return &fixedHostSession{scope: SessionScope{Principal:"analyst",Tenant:"t1",Environment:"e1"}}
}
func testQueueRequest(h http.Handler, method,path string) *httptest.ResponseRecorder {
	req:=httptest.NewRequest(method,path,nil)
	rec:=httptest.NewRecorder()
	h.ServeHTTP(rec,req)
	return rec
}

func TestSingleCommandRouteRendersCanonicalSixViewsWithIndependentOwnerACL(t *testing.T) {
	for _,view:=range SystemViews {
		t.Run(string(view),func(t *testing.T){
			sessions:=trustedHost()
			owner:=fixture()
			handler:=NewHTTPHandler(sessions,owner)
			r:=testQueueRequest(handler,http.MethodGet,"/command/work-queue?view="+string(view))
			if r.Code!=200 {t.Fatalf("route rejected canonical view: %d %s",r.Code,r.Body.String())}
			if sessions.calls!=1 {t.Fatal("request did not resolve current authenticated host session")}
			if r.Header().Get("Cache-Control")!="no-store" ||
				r.Header().Get("X-Content-Type-Options")!="nosniff" ||
				!strings.HasPrefix(r.Header().Get("Content-Type"),"application/json") {
				t.Fatalf("response caching/content protection incorrect: %+v",r.Header())
			}
			var result Result
			if err:=json.Unmarshal(r.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
			if result.View!=view||result.Tenant!="t1"||result.Environment!="e1" {
				t.Fatalf("owner scope/view missing: %+v",result)
			}
			for _,row:=range result.Rows {
				if row.ID=="inc-hidden" || row.ID=="foreign" || row.Status=="secret" {
					t.Fatalf("forbidden incident exposed: %+v",row)
				}
			}
			for _,called:=range owner.calls {
				if called=="Incident:foreign" {t.Fatal("foreign source identifier submitted to host ACL")}
			}
		})
	}
	sessions:=trustedHost()
	rec:=testQueueRequest(NewHTTPHandler(sessions,fixture()),http.MethodGet,"/command/work-queue")
	if rec.Code!=200 || !strings.Contains(rec.Body.String(),`"view":"all"`) {
		t.Fatalf("no selected Saved View should default to canonical All: %d %q",rec.Code,rec.Body.String())
	}
}

func TestCommandRouteNeverTrustsClientScopeOrUnrecognizedView(t *testing.T) {
	for _,raw:=range []string{
		"/command/work-queue?view=team-load",
		"/command/work-queue?view=",
		"/command/work-queue?view=all&view=tasks",
		"/command/work-queue?tenant=t2",
		"/command/work-queue?principal=admin",
		"/command/work-queue?environment=e2",
		"/command/work-queue?filter=secret",
	} {
		sessions:=trustedHost()
		owner:=fixture()
		r:=testQueueRequest(NewHTTPHandler(sessions,owner),http.MethodGet,raw)
		if r.Code!=http.StatusBadRequest || len(owner.calls)!=0 {
			t.Fatalf("untrusted client scope or unsupported view accepted: %s status=%d body=%q",raw,r.Code,r.Body.String())
		}
		if strings.Contains(r.Body.String(),"inc-hidden") {t.Fatalf("error leaked private Case: %s",raw)}
	}
	for _,tc:=range []struct {method,route string; status int}{
		{http.MethodPost,"/command/work-queue",http.StatusMethodNotAllowed},
		{http.MethodGet,"/command/work-queue/incidents",http.StatusNotFound},
		{http.MethodGet,"/command/work-queue/tasks",http.StatusNotFound},
	} {
		r:=testQueueRequest(NewHTTPHandler(trustedHost(),fixture()),tc.method,tc.route)
		if r.Code!=tc.status {t.Fatalf("route/method access was not restricted: %s %s %d",tc.method,tc.route,r.Code)}
	}
}

func TestCommandHTTPRejectsMissingHostAuthRevocationAndOwnerFailure(t *testing.T) {
	for _,session:=range []*fixedHostSession{
		{err:errors.New("expired")},
		{scope:SessionScope{Principal:"",Tenant:"t1",Environment:"e1"}},
		{scope:SessionScope{Principal:"analyst",Tenant:"*",Environment:"e1"}},
		{scope:SessionScope{Principal:"analyst",Tenant:"t1",Environment:""}},
	} {
		owner:=fixture()
		r:=testQueueRequest(NewHTTPHandler(session,owner),http.MethodGet,"/command/work-queue?view=all")
		if r.Code!=http.StatusUnauthorized || len(owner.calls)!=0 {
			t.Fatalf("unresolved host scope reached Incident owner: %d",r.Code)
		}
	}
	for _,o:=range []*fakeOwner{
		{errList:true},
		{errRead:true,items:fixture().items,allowed:fixture().allowed},
	} {
		r:=testQueueRequest(NewHTTPHandler(trustedHost(),o),http.MethodGet,"/command/work-queue")
		if r.Code!=http.StatusServiceUnavailable || strings.Contains(r.Body.String(),"inc-a") {
			t.Fatalf("provider failure disclosed partial rows: %d %q",r.Code,r.Body.String())
		}
	}
	if rec:=testQueueRequest(NewHTTPHandler(nil,fixture()),http.MethodGet,"/command/work-queue");rec.Code!=http.StatusServiceUnavailable {
		t.Fatal("missing host session resolver was allowed")
	}
	if rec:=testQueueRequest(NewHTTPHandler(trustedHost(),nil),http.MethodGet,"/command/work-queue");rec.Code!=http.StatusServiceUnavailable {
		t.Fatal("missing owner allowed access")
	}
}
