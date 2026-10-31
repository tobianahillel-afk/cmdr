package workqueue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func fetchAuthorizedPage(t *testing.T, h http.Handler, query string) Result {
	t.Helper()
	r:=testQueueRequest(h,http.MethodGet,query)
	if r.Code!=http.StatusOK {t.Fatalf("authorized page unavailable: %s %d %s",query,r.Code,r.Body.String())}
	var data Result
	if err:=json.Unmarshal(r.Body.Bytes(),&data);err!=nil {t.Fatal(err)}
	return data
}

func TestCommandWorkQueueBackwardNavigationReauthorizesAndPreservesContext(t *testing.T) {
	owner:=largeAuthorizedQueue()
	host:=trustedHost()
	handler:=NewHTTPHandler(host,owner)
	path:="/command/work-queue?view=tasks&status=open&sort=oldest"
	first:=fetchAuthorizedPage(t,handler,path)
	if len(first.Rows)!=50 || first.HasPrevious || first.PreviousCursor!="" || first.NextCursor=="" {
		t.Fatalf("first page exposes incorrect backward cursor: %+v",first)
	}
	secondPath:=path+"&cursor="+url.QueryEscape(first.NextCursor)
	second:=fetchAuthorizedPage(t,handler,secondPath)
	if len(second.Rows)!=50 || !second.HasPrevious || second.PreviousCursor!="" ||
		second.Rows[0].ID!="task-050" || second.NextCursor=="" {
		t.Fatalf("second page cannot return to first without a cursor: %+v",second)
	}
	third:=fetchAuthorizedPage(t,handler,path+"&cursor="+url.QueryEscape(second.NextCursor))
	if len(third.Rows)!=23 || !third.HasPrevious || third.NextCursor!="" ||
		third.PreviousCursor=="" || third.Rows[0].ID!="task-100" {
		t.Fatalf("third page has no usable preceding page: %+v",third)
	}
	// Same owner graph, *new* host resolution and CanRead for each step.
	backToSecond:=fetchAuthorizedPage(t,handler,path+"&cursor="+url.QueryEscape(third.PreviousCursor))
	if backToSecond.Rows[0].ID!="task-050" || backToSecond.Rows[49].ID!="task-099" ||
		backToSecond.NextCursor!=second.NextCursor || backToSecond.PreviousCursor!="" {
		t.Fatalf("preceding cursor did not restore authorized second page: %+v",backToSecond)
	}
	backToFirst:=fetchAuthorizedPage(t,handler,path)
	if len(backToFirst.Rows)!=50 || backToFirst.Rows[0].ID!="task-000" ||
		backToFirst.Rows[49].ID!="task-049" || backToFirst.HasPrevious {
		t.Fatalf("previous page did not restore exact first page: %+v",backToFirst)
	}
	if host.calls!=5 || len(owner.calls)!=5*124 {
		t.Fatalf("back navigation bypassed fresh session or per-record authorization: host=%d owner=%d",host.calls,len(owner.calls))
	}
	for _,call:=range owner.calls {
		if call=="Task:foreign-secret" {t.Fatal("back navigation queried foreign tenant owner")}
	}
	// Previously accepted prior-page cursor MUST now fail rather than
	// disclosing rows when current owner access has been revoked.
	owner.allowed["Task:task-012"]=false
	response:=testQueueRequest(handler,http.MethodGet,path+"&cursor="+url.QueryEscape(third.PreviousCursor))
	if response.Code!=http.StatusNotFound || strings.Contains(response.Body.String(),"task-050") {
		t.Fatalf("revoked owner left prior page accessible: %d %q",response.Code,response.Body.String())
	}
}

func TestCommandWorkQueueBrowserPreviousLinkMaintainsSameRouteAndEscaping(t *testing.T) {
	owner:=largeAuthorizedQueue()
	h:=NewHTTPHandler(trustedHost(),owner)
	base:="/command/work-queue?view=tasks&status=open&sort=oldest"
	first:=fetchAuthorizedPage(t,h,base)
	second:=fetchAuthorizedPage(t,h,base+"&cursor="+url.QueryEscape(first.NextCursor))
	thirdQuery:=base+"&cursor="+url.QueryEscape(second.NextCursor)
	request:=httptest.NewRequest(http.MethodGet,thirdQuery,nil)
	request.Header.Set("Accept","text/html")
	html:=httptest.NewRecorder()
	h.ServeHTTP(html,request)
	page:=html.Body.String()
	expectedPrevious:=base+"&cursor="+url.QueryEscape(second.NextCursor)
	_ = expectedPrevious // The next-page token cannot double as previous.
	previousHref:=base+"&cursor="+url.QueryEscape(thirdPreviousToken(t,h,thirdQuery))
	escapedHref:=strings.ReplaceAll(previousHref,"&","&amp;")
	if html.Code!=http.StatusOK || !strings.Contains(page,">Page précédente</a>") ||
		!strings.Contains(page,"href=\""+escapedHref+"\"") ||
		!strings.Contains(page,">Première page</a>") ||
		strings.Contains(page,"task-secret") || strings.Contains(page,"foreign-secret") {
		t.Fatalf("browser previous link lost filters or leaked data: %d %s",html.Code,page)
	}
	// Page 2 has exactly one-step return to page 1, without emitting
	// a fictitious token with start=0, which the cursor codec disallows.
	page2Request:=httptest.NewRequest(http.MethodGet,base+"&cursor="+url.QueryEscape(first.NextCursor),nil)
	page2Request.Header.Set("Accept","text/html")
	page2:=httptest.NewRecorder();h.ServeHTTP(page2,page2Request)
	if page2.Code!=200 || !strings.Contains(page2.Body.String(),">Page précédente</a>") ||
		!strings.Contains(page2.Body.String(),"href=\""+strings.ReplaceAll(base,"&","&amp;")+"\"") {
		t.Fatalf("second page missing direct first-page return: %d %s",page2.Code,page2.Body.String())
	}
}

func thirdPreviousToken(t *testing.T, h http.Handler, thirdQuery string) string {
	t.Helper()
	page:=fetchAuthorizedPage(t,h,thirdQuery)
	if page.PreviousCursor=="" {t.Fatal("missing preceding-page cursor")}
	return page.PreviousCursor
}
