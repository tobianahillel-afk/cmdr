package workqueue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestOwnerFilterRequiresObjectAuthorizationAndScopesSLAWarning(t *testing.T) {
	o:=fixture()
	req:=request(SLARisk)
	req.OwnerFilter="analyst"
	got,err:=Project(context.Background(),o,req)
	if err!=nil || got.State!="ready" || len(got.Rows)!=1 ||
		got.Rows[0].ID!="inc-a" || got.OwnerFilter!="analyst" ||
		len(got.Diagnostics)!=0 {
		t.Fatalf("filtered-owner SLA view contaminated by other-owner SLA gap: %v %+v",err,got)
	}
	// Other owner is readable but has no SLA policy: keep an honest warning.
	req.OwnerFilter="other"
	other,err:=Project(context.Background(),o,req)
	if err!=nil || other.State!="partial-sources" || len(other.Rows)!=0 ||
		len(other.Diagnostics)!=1 || other.Diagnostics[0]!="sla-policy-unavailable" {
		t.Fatalf("relevant owner SLA unknown suppressed: %v %+v",err,other)
	}
	req=request(All)
	req.OwnerFilter="other"
	other,err=Project(context.Background(),o,req)
	if err!=nil || len(other.Rows)!=1 || other.Rows[0].ID!="task-no-sla" {
		t.Fatalf("actual authorized owner did not filter correctly: %v %+v",err,other)
	}
	for _,id:=range o.calls {
		if id=="Incident:foreign" {t.Fatal("cross-tenant ID submitted to ACL")}
	}
	// Authorized row values are never treated as authority over other
	// objects, even if caller requests the string of a denied item's owner.
	req.OwnerFilter="secret"
	denied,err:=Project(context.Background(),fixture(),req)
	if err!=nil || denied.State!="no-authorized-items" || len(denied.Rows)!=0 {
		t.Fatalf("denied owner matched hidden row: %v %+v",err,denied)
	}
	for _,invalid:=range []string{"*","\n",strings.Repeat("a",129)} {
		req.OwnerFilter=invalid
		badOwner:=fixture()
		_,err:=Project(context.Background(),badOwner,req)
		if err!=ErrView || len(badOwner.calls)!=0 {
			t.Fatalf("bad owner reached per-item ACL: %q %v",invalid,err)
		}
	}
}

func TestNativeOwnerFilterRetainsAuthorizedPageAndSelection(t *testing.T) {
	o:=largeAuthorizedQueue()
	for i:=range o.items {
		if strings.HasPrefix(o.items[i].ID,"task-") && o.items[i].ID>="task-075" {
			o.items[i].Owner="other"
		}
	}
	host:=trustedHost()
	h:=NewHTTPHandler(host,o)
	firstURL:="/command/work-queue?view=tasks&status=open&owner=analyst&sort=oldest"
	first:=fetchAuthorizedPage(t,h,firstURL)
	if len(first.Rows)!=50 || first.NextCursor=="" || first.OwnerFilter!="analyst" {
		t.Fatalf("owner-scoped page 1 incorrect: %+v",first)
	}
	nextURL:=firstURL+"&cursor="+url.QueryEscape(first.NextCursor)
	second:=fetchAuthorizedPage(t,h,nextURL)
	if len(second.Rows)!=25 || second.Rows[0].ID!="task-050" ||
		second.NextCursor!="" || !second.HasPrevious {
		t.Fatalf("owner-scoped page 2 incorrect: %+v",second)
	}
	req:=httptest.NewRequest(http.MethodGet,nextURL,nil)
	req.Header.Set("Accept","text/html")
	htmlPage:=httptest.NewRecorder()
	h.ServeHTTP(htmlPage,req)
	if htmlPage.Code!=200 || !strings.Contains(htmlPage.Body.String(), `name="owner" value="analyst"`) ||
		!strings.Contains(htmlPage.Body.String(), "owner=analyst") ||
		strings.Contains(htmlPage.Body.String(),"task-075") ||
		strings.Contains(htmlPage.Body.String(),"task-secret") {
		t.Fatalf("owner filter absent or data leaked in page 2 HTML: %d %q",htmlPage.Code,htmlPage.Body.String())
	}
	query:=pageFormQuery(t,htmlPage.Body.String(),"Task:task-070")
	if query.Get("owner")!="analyst" || query.Get("cursor")!=first.NextCursor ||
		query.Get("status")!="open" || query.Get("sort")!="oldest" {
		t.Fatalf("selected second-page context lost owner: %v",query)
	}
	selected:=testQueueRequest(h,http.MethodGet,"/command/work-queue?"+query.Encode())
	var result Result
	if selected.Code!=200 || json.Unmarshal(selected.Body.Bytes(),&result)!=nil ||
		result.Selected==nil || result.Selected.ID!="task-070" ||
		result.OwnerFilter!="analyst" || len(result.Rows)!=25 {
		t.Fatalf("owner-scoped row selection failed: %d %q",selected.Code,selected.Body.String())
	}
	// A cursor is not a cross-owner capability, even while the principal
	// and tenant remain unchanged. It MUST be bound to the owner filter.
	wrong:=testQueueRequest(h,http.MethodGet,
		"/command/work-queue?view=tasks&status=open&owner=other&sort=oldest&cursor="+url.QueryEscape(first.NextCursor))
	if wrong.Code!=http.StatusNotFound || strings.Contains(wrong.Body.String(),"task-075") {
		t.Fatalf("cursor rebound across owner filters: %d %q",wrong.Code,wrong.Body.String())
	}
	// A revoked object invalidates page two, not merely the selected row.
	o.allowed["Task:task-070"]=false
	stale:=testQueueRequest(h,http.MethodGet,"/command/work-queue?"+query.Encode())
	if stale.Code!=http.StatusNotFound || strings.Contains(stale.Body.String(),"task-070") {
		t.Fatalf("revocation permitted stale page selection: %d %q",stale.Code,stale.Body.String())
	}
	if host.calls<5 {t.Fatalf("host did not reauthenticate per request: %d",host.calls)}
}

func TestOwnerFilterRejectsUntrustedHTTPQueryBeforeSession(t *testing.T) {
	for _,suffix:=range []string{
		"owner=%2A", "owner=%0A", "owner="+strings.Repeat("a",129),
		"owner=alice&owner=bob", "owner=alice&tenant=other",
	} {
		o:=fixture()
		host:=trustedHost()
		out:=testQueueRequest(NewHTTPHandler(host,o),http.MethodGet,"/command/work-queue?"+suffix)
		if out.Code!=http.StatusBadRequest || len(o.calls)!=0 || host.calls!=0 {
			t.Fatalf("invalid owner query reached session/ACL: %s %d",suffix,out.Code)
		}
	}
}
