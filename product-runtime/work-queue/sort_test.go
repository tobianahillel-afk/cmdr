package workqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestManualTimestampSortKeepsAuthorizationBeforeOrdering(t *testing.T) {
	for _, tc := range []struct {
		order string
		want []string
	}{
		{"", []string{"task-b", "inc-a", "task-no-sla"}},
		{"newest", []string{"task-b", "inc-a", "task-no-sla"}},
		{"oldest", []string{"task-no-sla", "inc-a", "task-b"}},
	} {
		t.Run(tc.order, func(t *testing.T) {
			owner := fixture()
			req := request(All)
			req.SortOrder = tc.order
			res, err := Project(context.Background(), owner, req)
			if err != nil { t.Fatal(err) }
			ids := make([]string, 0, len(res.Rows))
			for _, item := range res.Rows { ids = append(ids, item.ID) }
			if !reflect.DeepEqual(ids, tc.want) || res.SortOrder != tc.order {
				t.Fatalf("sort=%q result=%+v, want %v", tc.order, res, tc.want)
			}
			if len(owner.calls) != 4 {
				t.Fatalf("sort bypassed per-entity authorization: %+v", owner.calls)
			}
			for _, called := range owner.calls {
				if called == "Incident:foreign" { t.Fatal("foreign identifier queried") }
			}
		})
	}
	owner := fixture()
	req := request(Tasks)
	req.StatusFilter = "open"
	req.SortOrder = "oldest"
	got, err := Project(context.Background(), owner, req)
	if err != nil || len(got.Rows) != 2 ||
		got.Rows[0].ID != "task-no-sla" || got.Rows[1].ID != "task-b" {
		t.Fatalf("view/status intersection changed by chronological sort: %+v %v", got, err)
	}
	// Retain old JSON contract when no sort is selected.
	unsorted, err := Project(context.Background(), fixture(), request(All))
	if err != nil { t.Fatal(err) }
	b, err := json.Marshal(unsorted)
	if err != nil || strings.Contains(string(b), "\"sort\"") {
		t.Fatalf("backward-compatible JSON now includes forced sort: %q %v", b, err)
	}
}

func TestWorkQueueSortOneRouteJSONHTMLAndPreservedFilters(t *testing.T) {
	const path = "/command/work-queue?view=tasks&status=open&sort=oldest"
	jsonReq := httptest.NewRequest(http.MethodGet, path, nil)
	jsonReq.Header.Set("Accept", "application/json")
	jsonOut := httptest.NewRecorder()
	NewHTTPHandler(trustedHost(), fixture()).ServeHTTP(jsonOut, jsonReq)
	var result Result
	if jsonOut.Code != 200 || json.Unmarshal(jsonOut.Body.Bytes(), &result) != nil ||
		result.View != Tasks || result.StatusFilter != "open" || result.SortOrder != "oldest" ||
		len(result.Rows) != 2 || result.Rows[0].ID != "task-no-sla" {
		t.Fatalf("one-route authorized sort failed: %d %s", jsonOut.Code, jsonOut.Body.String())
	}
	for _, view := range SystemViews {
		path := "/command/work-queue?view=" + string(view) + "&status=open&sort=oldest"
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Accept", "text/html")
		out := httptest.NewRecorder()
		NewHTTPHandler(trustedHost(), fixture()).ServeHTTP(out, req)
		body := out.Body.String()
		if out.Code != 200 || strings.Count(body, "href=") != 6 ||
			strings.Count(body, "&amp;sort=oldest") != 6 ||
			strings.Count(body, "&amp;status=open") != 6 ||
			!strings.Contains(body, "name=\"sort\"") ||
			!strings.Contains(body, "<option value=\"oldest\" selected>") ||
			!strings.Contains(body, "name=\"status\"") {
			t.Fatalf("six-view sort/filter browser navigation broken: %s %d %s", view, out.Code, body)
		}
		if !strings.Contains(out.Header().Get("Content-Security-Policy"), "form-action 'self'") ||
			strings.Contains(body, "inc-hidden") || strings.Contains(body, "foreign") {
			t.Fatalf("unsafe HTML or unauthorized sorted row: %s", body)
		}
	}
}

func TestWorkQueueSortRejectsInjectedColumnsAndRevokedOwner(t *testing.T) {
	bad := []string{
		"/command/work-queue?sort=",
		"/command/work-queue?sort=random",
		"/command/work-queue?sort=owner",
		"/command/work-queue?sort=oldest&sort=newest",
		"/command/work-queue?sort=%FF",
		"/command/work-queue?sort=%0A",
		"/command/work-queue?sort=oldest&tenant=t2",
	}
	for _, path := range bad {
		o := fixture()
		res := testQueueRequest(NewHTTPHandler(trustedHost(), o), http.MethodGet, path)
		if res.Code != http.StatusBadRequest || len(o.calls) != 0 {
			t.Fatalf("untrusted sort reached provider: %s %d", path, res.Code)
		}
	}
	o := fixture()
	req := request(All)
	req.SortOrder = "custom-owner"
	if out, err := Project(context.Background(), o, req); err == nil || len(out.Rows) != 0 || len(o.calls) != 0 {
		t.Fatalf("unknown server sort allowed: %+v %v", out, err)
	}
	o = fixture()
	o.errRead = true
	r := testQueueRequest(NewHTTPHandler(trustedHost(), o), http.MethodGet, "/command/work-queue?sort=oldest")
	if r.Code != http.StatusServiceUnavailable || strings.Contains(r.Body.String(), "task-b") {
		t.Fatalf("revoked entity authorization leaked sorted rows: %d %q", r.Code, r.Body.String())
	}
}

func TestWorkQueueOldestFirstPageAppliesSortBeforeBound(t *testing.T) {
	owner := fixture()
	owner.items = nil
	owner.allowed = map[string]bool{}
	for i := 0; i < 65; i++ {
		id := fmt.Sprintf("task-%03d", i)
		owner.items = append(owner.items, Item{
			Kind: Task, ID: id, Tenant: "t1", Environment: "e1",
			Status: "open", UpdatedAt: "2026-10-19T09:00:00Z",
		})
		owner.allowed["Task:"+id] = true
	}
	req := request(Tasks)
	req.SortOrder = "oldest"
	result, err := Project(context.Background(), owner, req)
	if err != nil || len(result.Rows) != 50 ||
		!reflect.DeepEqual(result.Diagnostics, []string{"first-page-only"}) {
		t.Fatalf("bounded first page lost during sorting: %+v %v", result, err)
	}
	for i := 1; i < len(result.Rows); i++ {
		if result.Rows[i-1].ID >= result.Rows[i].ID {
			t.Fatalf("equal-time deterministic order changed: %+v", result.Rows)
		}
	}
}
