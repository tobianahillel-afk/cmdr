package workqueue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestExactStatusFilterIsAppliedOnlyAfterCurrentOwnerAuthorization(t *testing.T) {
	cases := []struct{ view View; status string; expected []string; state string }{
		{All,"open",[]string{"task-b","task-no-sla"},"ready"},
		{Tasks,"open",[]string{"task-b","task-no-sla"},"ready"},
		{Incidents,"active",[]string{"inc-a"},"ready"},
		{All,"secret",nil,"no-authorized-items"},
		{All,"missing",nil,"no-authorized-items"},
	}
	for _,tc := range cases {
		owner:=fixture()
		req:=request(tc.view)
		req.StatusFilter=tc.status
		result,err:=Project(context.Background(),owner,req)
		if err!=nil {t.Fatal(err)}
		var got []string
		for _,row:=range result.Rows { got=append(got,row.ID) }
		if !reflect.DeepEqual(got,tc.expected) || result.StatusFilter!=tc.status || result.State!=tc.state {
			t.Fatalf("status %q view %q unexpectedly returned %+v",tc.status,tc.view,result)
		}
		if len(owner.calls)!=4 {t.Fatalf("status filtering bypassed current owner checks: %+v",owner.calls)}
		for _,id:=range owner.calls {
			if id=="Incident:foreign" {t.Fatal("cross-tenant ID submitted to owner ACL")}
		}
	}
	noFilter,err:=Project(context.Background(),fixture(),request(All))
	if err!=nil || noFilter.StatusFilter!="" || len(noFilter.Rows)!=3 {
		t.Fatalf("legacy unfiltered projection changed: %+v %v",noFilter,err)
	}
}

func TestCommandStatusFilterOneRouteJSONAndHTMLNavigation(t *testing.T) {
	route:="/command/work-queue?view=tasks&status=open"
	owner:=fixture()
	req:=httptest.NewRequest(http.MethodGet,route,nil)
	req.Header.Set("Accept","application/json")
	out:=httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),owner).ServeHTTP(out,req)
	var jsonResult Result
	if out.Code!=200 || json.Unmarshal(out.Body.Bytes(),&jsonResult)!=nil ||
		jsonResult.View!=Tasks || jsonResult.StatusFilter!="open" ||
		len(jsonResult.Rows)!=2 || jsonResult.Rows[0].ID!="task-b" {
		t.Fatalf("filtered JSON did not use canonical projection: %d %q",out.Code,out.Body.String())
	}
	if out.Header().Get("Cache-Control")!="no-store" {t.Fatal("filtered response may be cached")}
	hreq:=httptest.NewRequest(http.MethodGet,route,nil)
	hreq.Header.Set("Accept","text/html")
	html:=httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),fixture()).ServeHTTP(html,hreq)
	body:=html.Body.String()
	if html.Code!=200 || strings.Count(body,"/command/work-queue?view=")!=6 ||
		strings.Count(body,"status=open")!=6 ||
		!strings.Contains(body,"value=\"open\"") ||
		!strings.Contains(body,"<form method=\"get\" action=\"/command/work-queue\"") {
		t.Fatalf("filter not preserved in single-route browser: %d %q",html.Code,body)
	}
	if !strings.Contains(html.Header().Get("Content-Security-Policy"),"form-action 'self'") ||
		strings.Contains(body,"inc-hidden") || strings.Contains(body,"foreign") {
		t.Fatalf("unsafe or unauthorized filtered browser result: %q",body)
	}
	// An authorized owner field can contain punctuation or HTML markup;
	// status is always URL-encoded on links and context-escaped in text.
	const status=`"><img src=x onerror=alert(1)>`
	malicious:=fixture()
	malicious.items[0].Status=status
	path:="/command/work-queue?view=incidents&status="+url.QueryEscape(status)
	m:=httptest.NewRequest(http.MethodGet,path,nil)
	m.Header.Set("Accept","text/html")
	safe:=httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),malicious).ServeHTTP(safe,m)
	if safe.Code!=200 || !strings.Contains(safe.Body.String(),"inc-a") ||
		strings.Contains(safe.Body.String(),"<img ") || strings.Contains(safe.Body.String(),"onerror=alert(1)>") {
		t.Fatalf("unsafe HTML from status filter: %d %q",safe.Code,safe.Body.String())
	}
}

func TestCommandStatusFilterFailsClosedAtUntrustedBoundary(t *testing.T) {
	bad:=[]string{
		"/command/work-queue?status=x&status=y",
		"/command/work-queue?view=all&status=%00",
		"/command/work-queue?view=all&status=%0A",
		"/command/work-queue?view=all&status="+strings.Repeat("x",65),
		"/command/work-queue?view=all&status=*",
		"/command/work-queue?view=all&status=%FF",
	}
	for _,path:=range bad {
		owner:=fixture()
		out:=testQueueRequest(NewHTTPHandler(trustedHost(),owner),http.MethodGet,path)
		if out.Code!=http.StatusBadRequest || len(owner.calls)!=0 ||
			strings.Contains(out.Body.String(),"inc-a") {
			t.Fatalf("invalid filter reached owner or leaked data: %q %d %q",path,out.Code,out.Body.String())
		}
	}
	filtered:=httptest.NewRequest(http.MethodGet,"/command/work-queue?view=all&status=open",nil)
	filtered.Header.Set("Accept","text/html")
	sessions:=&fixedHostSession{err:context.Canceled}
	out:=httptest.NewRecorder()
	NewHTTPHandler(sessions,fixture()).ServeHTTP(out,filtered)
	if out.Code!=http.StatusUnauthorized || strings.Contains(out.Body.String(),"<table") {
		t.Fatalf("revoked authenticated session produced filtered owner HTML: %d",out.Code)
	}
}
