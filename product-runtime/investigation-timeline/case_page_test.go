package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func decodeTimelinePage(t *testing.T, body []byte) timelinePagedHTTPResult {
	t.Helper()
	var response timelinePagedHTTPResult
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	return response
}

func TestCaseTimelineHTTPPagesAuthorizedChronology(t *testing.T) {
	handler, store := httpWorkspaceFixture(t)
	want := []string{"entry-old", "entry-new"}
	for offset := 0; offset < 3; offset++ {
		response := httpWorkspaceGet(handler,
			"/case-timeline?case_id=case-one&limit=1&offset="+strconv.Itoa(offset),
			"Bearer analyst-a", nil)
		if response.Code != http.StatusOK {
			t.Fatalf("page %d: HTTP %d %s", offset, response.Code, response.Body.String())
		}
		page := decodeTimelinePage(t, response.Body.Bytes())
		if page.Status != "filtered" || page.Page.Offset != offset ||
			page.Page.Limit != 1 || page.Page.TotalVisible != 2 ||
			page.Page.HasMore != (offset == 0) {
			t.Fatalf("page %d metadata: %+v", offset, page)
		}
		if offset < 2 && (len(page.Entries) != 1 || page.Entries[0].ID != want[offset]) {
			t.Fatalf("wrong chronology at page %d: %+v", offset, page)
		}
		if offset == 2 && len(page.Entries) != 0 {
			t.Fatalf("out-of-range page should be empty: %+v", page)
		}
		for _, hidden := range []string{"entry-hidden", "src-hidden", "tenant-b", "evidence-secret"} {
			if strings.Contains(response.Body.String(), hidden) {
				t.Fatalf("source identity %q leaked in paginated HTTP", hidden)
			}
		}
	}
	if store.caseChecks != 6 || store.catalogCalls != 3 ||
		store.entryReads != 9 || store.sourceReads != 9 {
		t.Fatalf("paging bypassed complete owner/ACL chain: %+v", store)
	}
	legacy := httpWorkspaceGet(handler, "/case-timeline?case_id=case-one", "Bearer analyst-a", nil)
	if legacy.Code != http.StatusOK || strings.Contains(legacy.Body.String(), "\"page\"") {
		t.Fatalf("unpaged response contract changed: %s", legacy.Body.String())
	}
	var result Result
	if err := json.Unmarshal(legacy.Body.Bytes(), &result); err != nil ||
		len(result.Entries) != 2 || result.Status != "filtered" {
		t.Fatalf("backward incompatible default response: %+v (%v)", result, err)
	}
}

func TestCaseTimelineHTTPPageAfterAuthorizationAndFilters(t *testing.T) {
	handler, _ := httpWorkspaceFixture(t)
	for _, suffix := range []string{
		"&kind=observed&limit=1",
		"&from=2026-10-01T11%3A00%3A00Z&limit=1&offset=0",
		"&kind=observed&from=2026-10-01T10%3A00%3A00Z&to=2026-10-02T00%3A00%3A00Z&limit=1",
	} {
		response := httpWorkspaceGet(handler, "/case-timeline?case_id=case-one"+suffix,
			"Bearer analyst-a", nil)
		page := decodeTimelinePage(t, response.Body.Bytes())
		if response.Code != http.StatusOK || page.Page.TotalVisible != 1 ||
			page.Page.HasMore || len(page.Entries) != 1 ||
			page.Entries[0].ID != "entry-new" {
			t.Fatalf("invalid combined view %q: %+v HTTP=%d", suffix, page, response.Code)
		}
	}
	denied := httpWorkspaceGet(handler, "/case-timeline?case_id=case-two&limit=1",
		"Bearer analyst-a", nil)
	if denied.Code != http.StatusNotFound ||
		strings.Contains(denied.Body.String(), "total_visible") {
		t.Fatalf("cross-tenant page metadata disclosed: %d %s", denied.Code, denied.Body.String())
	}
}

func TestMemoryTimelinePagesShowOnlyAuthorizedSourceCounts(t *testing.T) {
	handler, err := NewMemoryCaseTimelineHTTPHandler(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		bearer string
		total  int
		status string
	}{
		{"Bearer aa", 2, "filtered"},
		{"Bearer ab", 3, "allow"},
	} {
		request := httpWorkspaceGet(handler, "/case-timeline?case_id=case-a&limit=1",
			tc.bearer, nil)
		page := decodeTimelinePage(t, request.Body.Bytes())
		if request.Code != http.StatusOK || page.Page.TotalVisible != tc.total ||
			page.Status != tc.status || !page.Page.HasMore || len(page.Entries) != 1 {
			t.Fatalf("source ACL pagination %s: %+v", tc.bearer, page)
		}
		if tc.total == 2 && strings.Contains(request.Body.String(), "src-secret") {
			t.Fatal("denied source leaked")
		}
	}
}

func TestCaseTimelineHTTPRejectsInvalidPagingBeforeProviderIO(t *testing.T) {
	for _, suffix := range []string{
		"&offset=0", "&limit=0", "&limit=201", "&limit=-1", "&limit=+1",
		"&limit=01", "&limit=abc", "&limit=9999999999999999999999",
		"&limit=1&limit=2", "&limit=1&offset=-1",
		"&limit=1&offset=5001", "&limit=1&offset=01", "&limit=1&offset=abc",
		"&limit=1&offset=999999999999999999999999",
		"&kind=observed&from=2026-10-01T00%3A00%3A00Z&to=2026-10-02T00%3A00%3A00Z&limit=1&offset=0&permission=allow",
	} {
		handler, store := httpWorkspaceFixture(t)
		response := httpWorkspaceGet(handler,
			"/case-timeline?case_id=case-one"+suffix, "Bearer analyst-a", nil)
		if response.Code != http.StatusBadRequest || store.caseChecks != 0 ||
			store.catalogCalls != 0 {
			t.Fatalf("invalid paging %s contacted owners: HTTP=%d checks=%d",
				suffix, response.Code, store.caseChecks)
		}
	}
}

func TestPageProjectionDoesNotModifyInput(t *testing.T) {
	in := Result{Status: "partial", Entries: []ProjectedEntry{
		{ID: "a", Kind: "observed"}, {ID: "b", Kind: "inferred"}, {ID: "c", Kind: "observed"},
	}}
	copyEntries := append([]ProjectedEntry(nil), in.Entries...)
	out := paginateAuthorizedTimeline(in, timelineReadSelection{paged: true, limit: 2, offset: 1})
	if out.Page.TotalVisible != 3 || out.Page.HasMore ||
		len(out.Entries) != 2 || out.Entries[0].ID != "b" {
		t.Fatalf("incorrect bounded page: %+v", out)
	}
	out.Entries[0].ID = "mutated"
	if !reflect.DeepEqual(in.Entries, copyEntries) {
		t.Fatal("page mutated authorized result")
	}
	selection, ok := parseTimelineReadSelection(url.Values{"case_id": {"case-one"}, "limit": {"200"}, "offset": {"5000"}})
	if !ok || !selection.paged || selection.limit != 200 || selection.offset != 5000 {
		t.Fatal("max-page policy incorrectly rejected")
	}
}
