package workqueue

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func browserRequest(h http.Handler, view string) *httptest.ResponseRecorder {
	path := "/command/work-queue"
	if view != "" { path += "?view="+view }
	req := httptest.NewRequest(http.MethodGet,path,nil)
	req.Header.Set("Accept","text/html")
	result := httptest.NewRecorder()
	h.ServeHTTP(result,req)
	return result
}

func TestNativeWorkQueueBrowserReauthorizesAllSixViewsOnSingleRoute(t *testing.T) {
	for _,view := range SystemViews {
		t.Run(string(view),func(t *testing.T){
			sessions := trustedHost()
			owner := fixture()
			page := browserRequest(NewHTTPHandler(sessions,owner),string(view))
			body := page.Body.String()
			if page.Code != http.StatusOK || sessions.calls != 1 ||
				!strings.HasPrefix(page.Header().Get("Content-Type"),"text/html") ||
				page.Header().Get("Cache-Control")!="no-store" ||
				page.Header().Get("Vary")!="Accept" ||
				!strings.Contains(page.Header().Get("Content-Security-Policy"),"default-src 'none'") ||
				page.Header().Get("X-Frame-Options")!="DENY" {
				t.Fatalf("browser failed fresh host authentication or safe headers for %s: %d %+v",view,page.Code,page.Header())
			}
			if !strings.Contains(body,"aria-current=") ||
				!strings.Contains(body,"Tenant : t1") ||
				!strings.Contains(body,"Environnement : e1") ||
				!strings.Contains(body,"Vue : "+string(view)) {
				t.Fatalf("browser lost scoped view or navigation: %s",body)
			}
			if strings.Count(body,"href=")!=6 || strings.Count(body,"/command/work-queue?view=")!=6 {
				t.Fatalf("expected exactly six links to ONE canonical route: %s",body)
			}
			for _,forbidden := range []string{"inc-hidden","foreign","secret","/command/work-queue/incidents","<script"} {
				if strings.Contains(body,forbidden) {t.Fatalf("unauthorized data or false route in HTML: %q",forbidden)}
			}
			if view==SLARisk && (!strings.Contains(body,"partial-sources") ||
				!strings.Contains(body,"sla-policy-unavailable")) {
				t.Fatal("SLA source-unavailability was silently hidden")
			}
			for _,call:=range owner.calls {
				if call=="Incident:foreign" { t.Fatal("foreign ID used during browser authorization") }
			}
		})
	}
	all := browserRequest(NewHTTPHandler(trustedHost(),fixture()),"")
	if all.Code!=200 || !strings.Contains(all.Body.String(),"Vue : all") {
		t.Fatalf("browser default failed: %d %s",all.Code,all.Body.String())
	}
}

func TestWorkQueueBrowserEscapesOwnerControlledValuesAndKeepsJSON(t *testing.T) {
	owner:=fixture()
	attackerID:="<img src=x onerror=alert(1)>"
	owner.items[0].ID=attackerID
	owner.items[0].Status="<script>alert(2)</script>"
	owner.items[0].Owner="<svg onload=alert(3)>"
	owner.allowed["Incident:"+attackerID]=true
	page:=browserRequest(NewHTTPHandler(trustedHost(),owner),"incidents")
	body:=page.Body.String()
	if page.Code!=200 || !strings.Contains(body,"&lt;img") ||
		!strings.Contains(body,"&lt;script&gt;") ||
		!strings.Contains(body,"&lt;svg") {
		t.Fatalf("safe HTML encoding or authorized row missing: %d %q",page.Code,body)
	}
	for _,needle:=range []string{"<img","<script","<svg","href=javascript:"} {
		if strings.Contains(body,needle) {t.Fatalf("owner data injected executable HTML: %q",needle)}
	}
	req:=httptest.NewRequest(http.MethodGet,"/command/work-queue?view=incidents",nil)
	req.Header.Set("Accept","application/json")
	rec:=httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),fixture()).ServeHTTP(rec,req)
	var decoded Result
	if err:=json.Unmarshal(rec.Body.Bytes(),&decoded);err!=nil {t.Fatal(err)}
	if rec.Code!=200 || !strings.HasPrefix(rec.Header().Get("Content-Type"),"application/json") ||
		decoded.View!=Incidents || strings.Contains(rec.Body.String(),"<!doctype html>") {
		t.Fatalf("existing first-party JSON contract changed: %d %q",rec.Code,rec.Body.String())
	}
}

func TestWorkQueueBrowserDeniedOwnersAndBadViewsNeverRenderRows(t *testing.T) {
	for _,tc:=range []struct {
		view string
		session *fixedHostSession
		owner *fakeOwner
		status int
	}{
		{"all",&fixedHostSession{err:errors.New("session revoked")},fixture(),401},
		{"all",trustedHost(),&fakeOwner{errList:true},503},
		{"all",trustedHost(),&fakeOwner{items:fixture().items,allowed:fixture().allowed,errRead:true},503},
		{"team-load",trustedHost(),fixture(),400},
		{"incidents&tenant=foreign",trustedHost(),fixture(),400},
	} {
		page:=browserRequest(NewHTTPHandler(tc.session,tc.owner),tc.view)
		if page.Code!=tc.status || strings.Contains(page.Body.String(),"inc-a") ||
			strings.Contains(page.Body.String(),"<table") {
			t.Fatalf("browser disclosed rows on denial: %s %d %q",tc.view,page.Code,page.Body.String())
		}
	}
}
