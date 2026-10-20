package workqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func largeAuthorizedQueue() *fakeOwner {
	o := &fakeOwner{allowed:map[string]bool{}}
	for i:=0;i<123;i++{
		id:=fmt.Sprintf("task-%03d",i)
		o.items=append(o.items,Item{Kind:Task, ID:id, Tenant:"t1",
			Environment:"e1",Owner:"analyst",Status:"open",
			UpdatedAt:"2026-10-20T09:00:00Z"})
		o.allowed["Task:"+id]=true
	}
	// The denied in-scope record is never returned or used in the digest.
	o.items=append(o.items,Item{Kind:Task,ID:"task-secret",Tenant:"t1",
		Environment:"e1",Owner:"secret",Status:"hidden",
		UpdatedAt:"2026-10-20T09:00:00Z"})
	// A foreign tenant ID must never be submitted to current owner.CanRead.
	o.items=append(o.items,Item{Kind:Task,ID:"foreign-secret",Tenant:"t2",
		Environment:"e1",Status:"hidden",UpdatedAt:"2026-10-20T09:00:00Z"})
	return o
}

func TestCommandWorkQueueRealThreePageAuthorizedJourney(t *testing.T){
	o:=largeAuthorizedQueue()
	host:=trustedHost()
	h:=NewHTTPHandler(host,o)
	query:="/command/work-queue?view=tasks&status=open&sort=oldest"
	want:=[]int{50,50,23}
	for pageIndex,expected:=range want {
		former:=len(o.calls)
		r:=testQueueRequest(h,http.MethodGet,query)
		if r.Code!=200 {t.Fatalf("page %d failed %d %s",pageIndex,r.Code,r.Body.String())}
		var got Result
		if err:=json.Unmarshal(r.Body.Bytes(),&got);err!=nil {t.Fatal(err)}
		if len(got.Rows)!=expected || got.HasPrevious!=(pageIndex>0) {
			t.Fatalf("page %d returned wrong slice / previous marker: %+v",pageIndex,got)
		}
		if len(o.calls)-former!=124 {
			t.Fatalf("page %d did not reauthorize every in-scope Incident/Task: %d",pageIndex,len(o.calls)-former)
		}
		for _,call:=range o.calls{
			if call=="Task:foreign-secret"{t.Fatal("cross-tenant identifier queried")}
		}
		for i,row:=range got.Rows{
			wantID:=fmt.Sprintf("task-%03d",pageIndex*queuePageSize+i)
			if row.ID!=wantID || strings.Contains(row.ID,"secret"){
				t.Fatalf("page %d leaked / reordered %q, expected %q",pageIndex,row.ID,wantID)
			}
		}
		if pageIndex==0 && (len(got.Diagnostics)!=1||got.Diagnostics[0]!="first-page-only"){
			t.Fatalf("previous first-page diagnostic was lost: %+v",got.Diagnostics)
		}
		if pageIndex==2 && got.NextCursor!="" {t.Fatal("last page advertises non-existent page")}
		if got.NextCursor!=""{
			if len(got.NextCursor)>maxQueueCursorBytes {t.Fatal("unbounded cursor")}
			query="/command/work-queue?view=tasks&status=open&sort=oldest&cursor="+url.QueryEscape(got.NextCursor)
		}
	}
	if host.calls!=3 {t.Fatalf("host SessionResolver not re-invoked for each page: %d",host.calls)}
}

func TestCommandWorkQueueCurrentPageSelectsOnlyAuthorizedRow(t *testing.T){
	o:=largeAuthorizedQueue()
	h:=NewHTTPHandler(trustedHost(),o)
	first:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=tasks")
	var a Result
	if json.Unmarshal(first.Body.Bytes(),&a)!=nil || a.NextCursor=="" {t.Fatal("missing first-page cursor")}
	nextQuery:="/command/work-queue?view=tasks&cursor="+url.QueryEscape(a.NextCursor)
	selected:=testQueueRequest(h,http.MethodGet,nextQuery+"&selected=Task:task-070")
	var b Result
	if selected.Code!=200 || json.Unmarshal(selected.Body.Bytes(),&b)!=nil ||
		b.Selected==nil || b.Selected.ID!="task-070" || len(b.Rows)!=50 ||
		b.Rows[0].ID!="task-050" {
		t.Fatalf("second page did not authorize selected row: %d %s",selected.Code,selected.Body.String())
	}
	denied:=testQueueRequest(h,http.MethodGet,nextQuery+"&selected=Task:task-020")
	if denied.Code!=http.StatusNotFound || strings.Contains(denied.Body.String(),"task-020"){
		t.Fatalf("selected first-page ID gained authority on page two: %d %s",denied.Code,denied.Body.String())
	}
	denied=testQueueRequest(h,http.MethodGet,nextQuery+"&selected=Task:task-secret")
	if denied.Code!=http.StatusNotFound {t.Fatal("denied source selected on page two")}
}

func TestCommandWorkQueueCursorReauthDeniesRevokedOrChangedOwnerState(t *testing.T) {
	o:=largeAuthorizedQueue()
	first,err:=Project(context.Background(),o,request(Tasks))
	if err!=nil || first.NextCursor=="" {t.Fatalf("page token missing: %v",err)}
	cursor:=first.NextCursor
	// A new owner authorization snapshot must make a prior cursor stale,
	// even if its numeric position would still contain authorized rows.
	o.allowed["Task:task-070"]=false
	o.calls=nil
	var req=Request
	req=request(Tasks)
	req.Cursor=cursor
	result,err:=Project(context.Background(),o,req)
	if err!=ErrCursor || len(result.Rows)!=0 || result.Selected!=nil ||
		len(o.calls)!=124 {
		t.Fatalf("revoked owner allowed stale cursor: %v %+v calls=%d",err,result,len(o.calls))
	}
	o.allowed["Task:task-070"]=true
	req.SortOrder="oldest"
	result,err=Project(context.Background(),o,req)
	if err!=ErrCursor || len(result.Rows)!=0 {
		t.Fatalf("cursor rebound to another sort filter: %v %+v",err,result)
	}
	req=request(Incidents)
	req.Cursor=cursor
	result,err=Project(context.Background(),o,req)
	if err!=ErrCursor || len(result.Rows)!=0 {
		t.Fatalf("cursor rebound to another system view: %v %+v",err,result)
	}
	req=request(Tasks)
	req.Principal="other-principal"
	req.Cursor=cursor
	result,err=Project(context.Background(),o,req)
	if err!=ErrCursor || len(result.Rows)!=0 {
		t.Fatalf("cursor rebound across users: %v %+v",err,result)
	}
}

func TestCommandQueueCursorValidationBeforeOwnerAndNoHTMLLeak(t *testing.T) {
	for _,query:=range []string{
		"cursor=","cursor=not-base64$","cursor=a&cursor=b",
		"cursor="+strings.Repeat("a",400),
		"cursor=invalid&tenant=t2",
	}{
		o:=largeAuthorizedQueue();session:=trustedHost()
		resp:=testQueueRequest(NewHTTPHandler(session,o),http.MethodGet,
			"/command/work-queue?"+query)
		if resp.Code!=400 || session.calls!=0 || len(o.calls)!=0 {
			t.Fatalf("bad token reached owner or host: %s %d %d %d",query,resp.Code,session.calls,len(o.calls))
		}
	}
	o:=largeAuthorizedQueue();h:=NewHTTPHandler(trustedHost(),o)
	req:=httptest.NewRequest(http.MethodGet,"/command/work-queue?view=tasks&status=open&sort=oldest",nil)
	req.Header.Set("Accept","text/html")
	out:=httptest.NewRecorder()
	h.ServeHTTP(out,req)
	if out.Code!=200 || !strings.Contains(out.Body.String(),">Page suivante</a>") ||
		!strings.Contains(out.Body.String(),"cursor=") ||
		strings.Contains(out.Body.String(),"task-secret") ||
		strings.Contains(out.Body.String(),"foreign-secret") {
		t.Fatalf("next-page HTML link not permission-safe: %d %s",out.Code,out.Body.String())
	}
}
