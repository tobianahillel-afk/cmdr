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

func groupedFixture() *fakeOwner {
	o := largeAuthorizedQueue()
	for i := range o.items {
		if o.items[i].Tenant != "t1" || o.items[i].ID == "task-secret" { continue }
		// The first group has 45 authorized items and the second 78,
		// crossing two page boundaries. All items retain first-party ACLs.
		if o.items[i].ID < "task-045" { o.items[i].Status = "closed" } else { o.items[i].Status = "open" }
		o.items[i].UpdatedAt = fmt.Sprintf("2026-10-20T09:%02d:00Z",i%60)
	}
	return o
}

func TestCanonicalManualGroupingAcrossAuthorizedPagesAndSelection(t *testing.T) {
	o:=groupedFixture()
	host:=trustedHost()
	h:=NewHTTPHandler(host,o)
	base:="/command/work-queue?view=tasks&group=status&sort=oldest"
	first:=fetchAuthorizedPage(t,h,base)
	if first.GroupBy!="status" || len(first.Rows)!=50 || first.NextCursor=="" ||
		first.Rows[0].Status!="closed" || first.Rows[44].Status!="closed" ||
		first.Rows[45].Status!="open" {
		t.Fatalf("grouping did not apply to all 123 authorized rows before pagination: %+v",first)
	}
	for _,row:=range first.Rows{
		if row.ID=="task-secret" || row.ID=="foreign-secret" {t.Fatalf("denied/foreign item grouped: %+v",row)}
	}
	secondURL:=base+"&cursor="+url.QueryEscape(first.NextCursor)
	second:=fetchAuthorizedPage(t,h,secondURL)
	if len(second.Rows)!=50 || second.Rows[0].Status!="open" ||
		second.NextCursor=="" || !second.HasPrevious {
		t.Fatalf("group cursor page 2 incorrect: %+v",second)
	}
	third:=fetchAuthorizedPage(t,h,base+"&cursor="+url.QueryEscape(second.NextCursor))
	if len(third.Rows)!=23 || third.NextCursor!="" {
		t.Fatalf("group cursor page 3 incorrect: %+v",third)
	}
	seen:=map[string]bool{}
	for _,page:=range []Result{first,second,third} {
		var lastGroup string
		var lastID string
		for _,row:=range page.Rows {
			if seen[row.ID] {t.Fatalf("row duplicated across pages: %s",row.ID)}
			seen[row.ID]=true
			if lastGroup!="" && row.Status < lastGroup {t.Fatal("status grouping not sorted")}
			if lastGroup==row.Status && lastID==row.ID {t.Fatal("duplicate group row")}
			lastGroup,lastID=row.Status,row.ID
		}
	}
	if len(seen)!=123 {t.Fatalf("grouping dropped authorized rows: %d",len(seen))}
	req:=httptest.NewRequest(http.MethodGet,secondURL,nil)
	req.Header.Set("Accept","text/html")
	browser:=httptest.NewRecorder(); h.ServeHTTP(browser,req)
	html:=browser.Body.String()
	if browser.Code!=http.StatusOK ||
		!strings.Contains(html,`name="group" value="status"`) ||
		!strings.Contains(html,"group=status") ||
		!strings.Contains(html,`scope="colgroup"`) ||
		strings.Contains(html,"task-secret") || strings.Contains(html,"foreign-secret") {
		t.Fatalf("group context/ACL not preserved in browser: %d %s",browser.Code,html)
	}
	token:=string(second.Rows[20].Kind)+":"+second.Rows[20].ID
	selectedQuery:=pageFormQuery(t,html,token)
	if selectedQuery.Get("group")!="status" || selectedQuery.Get("cursor")!=first.NextCursor {
		t.Fatalf("selected owner page lost group context: %v",selectedQuery)
	}
	selected:=testQueueRequest(h,http.MethodGet,"/command/work-queue?"+selectedQuery.Encode())
	var selectedResult Result
	if selected.Code!=200 || json.Unmarshal(selected.Body.Bytes(),&selectedResult)!=nil ||
		selectedResult.Selected==nil || selectedResult.Selected.ID!=second.Rows[20].ID {
		t.Fatalf("grouped page item could not be selected after reauthorization: %d %s",selected.Code,selected.Body.String())
	}
	// A cursor created under one grouping cannot be used for another
	// grouping or for an ungrouped projection, even with same principal.
	for _,q:=range []string{
		"/command/work-queue?view=tasks&group=type&sort=oldest&cursor="+url.QueryEscape(first.NextCursor),
		"/command/work-queue?view=tasks&sort=oldest&cursor="+url.QueryEscape(first.NextCursor),
	} {
		out:=testQueueRequest(h,http.MethodGet,q)
		if out.Code!=http.StatusNotFound {
			t.Fatalf("cursor reused outside its authorized grouping: %d %s",out.Code,out.Body.String())
		}
	}
	// Per-request owner revocation must invalidate preexisting group cursors.
	o.allowed["Task:"+second.Rows[20].ID]=false
	revoked:=testQueueRequest(h,http.MethodGet,secondURL)
	if revoked.Code!=http.StatusNotFound ||
		strings.Contains(revoked.Body.String(),second.Rows[20].ID){
		t.Fatalf("grouping cursor bypassed revoked owner: %d %s",revoked.Code,revoked.Body.String())
	}
	if host.calls<7 {t.Fatalf("HTTP host was not consulted on every grouped request: %d",host.calls)}
}

func TestManualGroupingTypeUsesExistingKindWithoutNewOwnerSource(t *testing.T) {
	o:=fixture()
	req:=request(All)
	req.GroupBy="type"
	result,err:=Project(context.Background(),o,req)
	if err!=nil || result.GroupBy!="type" {t.Fatalf("known grouping rejected: %v %+v",err,result)}
	wasTask:=false
	for _,row:=range result.Rows {
		if row.Kind==Task {wasTask=true}
		if wasTask && row.Kind==Incident {t.Fatal("Incident appeared after Task in type grouping")}
		if row.ID=="inc-hidden" || row.Status=="secret" || row.ID=="foreign" {
			t.Fatalf("grouped result disclosed denied owner item: %+v",row)
		}
	}
	raw,_:=json.Marshal(result)
	if !strings.Contains(string(raw),`"group_by":"type"`) {t.Fatalf("group selection missing JSON receipt: %s",raw)}
}

func TestManualGroupingRejectsBadInputBeforeSessionAndKeepsLegacy(t *testing.T) {
	for _,query:=range []string{
		"group=severity", "group=priority", "group=", "group=%2A",
		"group=type&group=status", "group=type&tenant=other",
		"group="+strings.Repeat("x",200),
	} {
		o:=fixture();session:=trustedHost()
		out:=testQueueRequest(NewHTTPHandler(session,o),http.MethodGet,
			"/command/work-queue?"+query)
		if out.Code!=http.StatusBadRequest || session.calls!=0 || len(o.calls)!=0 {
			t.Fatalf("invalid grouped query accessed session/ACL: %s %d",query,out.Code)
		}
	}
	o:=fixture()
	req:=request(All)
	req.GroupBy="severity"
	if _,err:=Project(context.Background(),o,req);err!=ErrView || len(o.calls)!=0 {
		t.Fatalf("unsupported direct projection group accessed owner: %v",err)
	}
	legacy:=fetchAuthorizedPage(t,NewHTTPHandler(trustedHost(),fixture()),"/command/work-queue?view=all")
	if legacy.GroupBy!="" {
		t.Fatal("legacy ungrouped route unexpectedly altered")
	}
}
