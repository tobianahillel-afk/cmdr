package workqueue

import (
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

var nativeHiddenInput = regexp.MustCompile(`<input type="hidden" name="([^"]+)" value="([^"]*)">`)

// Reconstruct an ordinary browser GET from one independently rendered form.
// This checks the actual HTML form contract, rather than assembling a cursor
// manually and thereby missing the original page-reset defect.
func pageFormQuery(t *testing.T, page, selected string) url.Values {
	t.Helper()
	needle := `name="selected" value="` + html.EscapeString(selected) + `"`
	index := strings.Index(page, needle)
	if index < 0 { t.Fatalf("selected row control missing: %q", selected) }
	start := strings.LastIndex(page[:index], `<form method="get" action="/command/work-queue">`)
	if start < 0 { t.Fatal("selected row is not within its own GET form") }
	end := strings.Index(page[index:], "</form>")
	if end < 0 { t.Fatal("selected-row GET form unclosed") }
	form := page[start : index+end]
	values := url.Values{}
	for _, hit := range nativeHiddenInput.FindAllStringSubmatch(form, -1) {
		values.Set(hit[1], html.UnescapeString(hit[2]))
	}
	values.Set("selected", selected)
	return values
}

func TestNativeWorkQueueSelectsActualPageThreeRowWithFreshOwnerAuthorization(t *testing.T) {
	owner := largeAuthorizedQueue()
	host := trustedHost()
	route := NewHTTPHandler(host, owner)
	base := "/command/work-queue?view=tasks&status=open&sort=oldest"
	first := fetchAuthorizedPage(t, route, base)
	second := fetchAuthorizedPage(t, route, base+"&cursor="+url.QueryEscape(first.NextCursor))
	pageThree := base+"&cursor="+url.QueryEscape(second.NextCursor)
	req := httptest.NewRequest(http.MethodGet, pageThree, nil)
	req.Header.Set("Accept", "text/html")
	thirdHTML := httptest.NewRecorder()
	route.ServeHTTP(thirdHTML, req)
	if thirdHTML.Code != http.StatusOK { t.Fatalf("third page unavailable: %d", thirdHTML.Code) }
	body := thirdHTML.Body.String()
	if strings.Contains(body, "task-foreign") || strings.Contains(body, "foreign-secret") {
		t.Fatal("third page disclosed a foreign-tenant row")
	}
	// Only the row's selection form carries the third-page cursor. The
	// editable filter form above the table must reset to page one.
	filterEnd := strings.Index(body, "</form>")
	if filterEnd < 0 || strings.Contains(body[:filterEnd], `name="cursor"`) ||
		strings.Contains(body[:filterEnd], `name="selected"`) {
		t.Fatal("manual filters unexpectedly retain a cursor or a selection")
	}
	values := pageFormQuery(t, body, "Task:task-100")
	if values.Get("cursor") != second.NextCursor ||
		values.Get("view") != "tasks" || values.Get("status") != "open" ||
		values.Get("sort") != "oldest" {
		t.Fatalf("third-page row click did not preserve current context: %v", values)
	}
	selection := testQueueRequest(route, http.MethodGet, "/command/work-queue?"+values.Encode())
	var got Result
	if selection.Code != http.StatusOK || json.Unmarshal(selection.Body.Bytes(), &got) != nil ||
		got.Selected == nil || got.Selected.ID != "task-100" ||
		len(got.Rows) != 23 || got.Rows[0].ID != "task-100" ||
		got.View != Tasks || got.SortOrder != "oldest" || got.StatusFilter != "open" {
		t.Fatalf("row selection from page three did not reopen same authorized page: %d %q",selection.Code,selection.Body.String())
	}
	if host.calls != 4 || len(owner.calls) != 4*124 {
		t.Fatalf("row selection bypassed normal per-request owner auth: %d %d",host.calls,len(owner.calls))
	}
	// Editing manual filters remains a fresh page-one query, as the current
	// cursor is not part of the filter form.
	reset := fetchAuthorizedPage(t,route,base)
	if reset.HasPrevious || reset.Rows[0].ID != "task-000" ||
		len(reset.Rows) != 50 {
		t.Fatalf("manual form context would not return to page one: %+v",reset)
	}
	// The first-page route has no authority to select an off-page ID.
	offPage := testQueueRequest(route,http.MethodGet,base+"&selected=Task:task-100")
	if offPage.Code != http.StatusNotFound ||
		strings.Contains(offPage.Body.String(),"task-100") {
		t.Fatalf("selection escaped page bound: %d %q",offPage.Code,offPage.Body.String())
	}
	// Changing a source ACL changes the authorized page digest. A previously
	// issued third-page selection fails closed; no obsolete row is returned.
	owner.allowed["Task:task-012"]=false
	stale := testQueueRequest(route,http.MethodGet,"/command/work-queue?"+values.Encode())
	if stale.Code != http.StatusNotFound ||
		strings.Contains(stale.Body.String(),"task-100") {
		t.Fatalf("revoked owner accepted a stale selected-page token: %d %q",stale.Code,stale.Body.String())
	}
}

func TestNativeWorkQueuePageTwoSelectionAndFreshSortReset(t *testing.T) {
	owner := largeAuthorizedQueue()
	handler := NewHTTPHandler(trustedHost(),owner)
	base := "/command/work-queue?view=tasks&status=open&sort=oldest"
	first := fetchAuthorizedPage(t,handler,base)
	secondPagePath := base+"&cursor="+url.QueryEscape(first.NextCursor)
	request := httptest.NewRequest(http.MethodGet,secondPagePath,nil)
	request.Header.Set("Accept","text/html")
	page := httptest.NewRecorder()
	handler.ServeHTTP(page,request)
	if page.Code != http.StatusOK {t.Fatalf("page 2 HTML failed: %d",page.Code)}
	values := pageFormQuery(t,page.Body.String(),"Task:task-050")
	if values.Get("cursor") != first.NextCursor {t.Fatal("page two button lost its authorized cursor")}
	back := testQueueRequest(handler,http.MethodGet,"/command/work-queue?"+values.Encode())
	var result Result
	if back.Code != http.StatusOK || json.Unmarshal(back.Body.Bytes(),&result)!=nil ||
		result.Selected==nil || result.Selected.ID!="task-050" ||
		len(result.Rows)!=50 || result.Rows[0].ID!="task-050" {
		t.Fatalf("page two native selection failed: %d %s",back.Code,back.Body.String())
	}
	filterForm:=page.Body.String()[:strings.Index(page.Body.String(),"</form>")]
	if strings.Contains(filterForm,`name="cursor"`) {
		t.Fatal("filter edit is incorrectly pinned to an old page")
	}
	firstNewest:=fetchAuthorizedPage(t,handler,"/command/work-queue?view=tasks&status=open&sort=newest")
	if firstNewest.HasPrevious || len(firstNewest.Rows)!=50 ||
		firstNewest.Rows[0].ID != "task-122" {
		t.Fatalf("edited sort failed to restart at newly authorized first page: %+v",firstNewest)
	}
}
