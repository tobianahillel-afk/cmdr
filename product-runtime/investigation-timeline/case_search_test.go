package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCaseTimelineSearchAfterFullOwnerAuthorization(t *testing.T) {
	cases := []struct {
		name, query string
		ids []string
	}{
		{"case-insensitive-id", "ENTRY-NEW", []string{"entry-new"}},
		{"source-reference", "SRC-OLD", []string{"entry-old"}},
		{"provenance", "TRUSTED:EVIDENCE", []string{"entry-old"}},
		{"common-provenance", "trusted:", []string{"entry-old", "entry-new"}},
		{"unmatched", "no-match", nil},
		{"hidden-source", "src-hidden", nil},
		{"hidden-entry", "entry-hidden", nil},
		{"hidden-provenance", "restricted", nil},
		{"foreign-source", "src-other", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, store := httpWorkspaceFixture(t)
			path := "/case-timeline?case_id=case-one&limit=1&offset=0&q=" + url.QueryEscape(tc.query)
			resp := httpWorkspaceGet(h, path, "Bearer analyst-a", nil)
			if resp.Code != http.StatusOK || store.caseChecks != 2 ||
				store.catalogCalls != 1 || store.entryReads != 3 || store.sourceReads != 3 {
				t.Fatalf("search bypassed full authorization: status=%d checks=%d/%d/%d/%d",
					resp.Code, store.caseChecks, store.catalogCalls, store.entryReads, store.sourceReads)
			}
			var page timelinePagedHTTPResult
			if err := json.Unmarshal(resp.Body.Bytes(), &page); err != nil {
				t.Fatal(err)
			}
			if page.Page.TotalVisible != len(tc.ids) || len(page.Entries) != min(1, len(tc.ids)) ||
				page.Page.HasMore != (len(tc.ids) > 1) {
				t.Fatalf("unexpected search result %+v for %q", page, tc.query)
			}
			if len(page.Entries) != 0 && page.Entries[0].ID != tc.ids[0] {
				t.Fatalf("unexpected first search result %+v", page)
			}
			for _, hidden := range []string{"entry-hidden", "src-hidden", "restricted", "src-other", "tenant-b"} {
				if strings.Contains(resp.Body.String(), hidden) {
					t.Fatalf("search leaked inaccessible metadata %q", hidden)
				}
			}
		})
	}
}

func TestCaseTimelineSearchComposesWithTimeKindPaginationAndCaseBoundary(t *testing.T) {
	h, _ := httpWorkspaceFixture(t)
	params := url.Values{
		"case_id": {"case-one"}, "q": {"trusted:"},
		"kind": {"observed"}, "from": {"2026-10-01T11:00:00Z"},
		"to": {"2026-10-01T12:00:00Z"}, "limit": {"1"}, "offset": {"0"},
	}
	good := httpWorkspaceGet(h, "/case-timeline?"+params.Encode(), "Bearer analyst-a", nil)
	var page timelinePagedHTTPResult
	if good.Code != http.StatusOK || json.Unmarshal(good.Body.Bytes(), &page) != nil ||
		page.Page.TotalVisible != 1 || len(page.Entries) != 1 || page.Entries[0].ID != "entry-new" {
		t.Fatalf("combined authorized search failed: %d %s", good.Code, good.Body.String())
	}
	params.Set("case_id", "case-two")
	denied := httpWorkspaceGet(h, "/case-timeline?"+params.Encode(), "Bearer analyst-a", nil)
	if denied.Code != http.StatusNotFound || strings.Contains(denied.Body.String(), "entry-other") {
		t.Fatalf("foreign Case search disclosed metadata: %d %s", denied.Code, denied.Body.String())
	}
	anonymous := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&q=trusted:", "", nil)
	if anonymous.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous search accepted: %d", anonymous.Code)
	}
}

func TestCaseTimelineSearchRejectsInvalidInputBeforeProviderCalls(t *testing.T) {
	for _, query := range []string{
		"q=", "q=%20%20", "q=term&q=other", "q=%00", "q=a%0Ab", "q=%FF",
		"q=" + strings.Repeat("x", 97), "q=source&tenant=tenant-b",
	} {
		h, store := httpWorkspaceFixture(t)
		resp := httpWorkspaceGet(h, "/case-timeline?case_id=case-one&"+query, "Bearer analyst-a", nil)
		if resp.Code != http.StatusBadRequest || store.caseChecks != 0 ||
			store.catalogCalls != 0 || store.entryReads != 0 || store.sourceReads != 0 {
			t.Fatalf("invalid search reached providers: query=%q HTTP=%d", query, resp.Code)
		}
	}
}
