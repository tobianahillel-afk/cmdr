package workqueue

import (
    "encoding/json"
    "html"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"
)

func TestLinkedCaseNavigationOnlyForIndependentlyAuthorizedSelectedIncident(t *testing.T) {
    // This is a local Command -> Investigate read-only route link. It conveys
    // exactly the Case ID, NOT the Command session or an authorization grant.
    links := &linkReferenceFixture{refs: &AuthorizedIncidentCaseReferences{References: []LinkedCaseReference{
        {CaseID: "case-allowed&tenant=forged#x", LinkSource: "Shared"},
        {CaseID: "case-hidden", LinkSource: "Shared"},
    }}}
    cases := &caseIdentityFixture{records: map[string]*AuthorizedCaseIdentity{
        "case-allowed&tenant=forged#x": {
            ID: "case-allowed&tenant=forged#x",
            Tenant: "t1", Environment: "e1", Status: "investigating",
            StatusSource: "Investigate", StatusObservedAt: "2026-10-20T12:00:00Z",
            StatusFreshness: "stale",
        },
    }}
    source := &independentlyJoinedCommandOwner{fakeOwner: fixture(),
        LinkedCaseOwnerJoin: LinkedCaseOwnerJoin{Links: links, Cases: cases}}
    h := NewHTTPHandler(trustedHost(), source)

    makeRequest := func(path string) *httptest.ResponseRecorder {
        t.Helper()
        req := httptest.NewRequest(http.MethodGet, path, nil)
        req.Header.Set("Accept", "text/html")
        out := httptest.NewRecorder()
        h.ServeHTTP(out, req)
        return out
    }
    selected := makeRequest("/command/work-queue?view=incidents&selected=Incident:inc-a")
    if selected.Code != http.StatusOK {
        t.Fatalf("selected Incident refused: %d %s", selected.Code, selected.Body.String())
    }
    markup := html.UnescapeString(selected.Body.String())
    // The path MUST include URL-encoded query values: an untrusted source
    // Case ID with ampersand/hash cannot create a tenant or environment query.
    wantLink := "/case-timeline?case_id=" + url.QueryEscape("case-allowed&tenant=forged#x")
    if !strings.Contains(markup, "href=\""+wantLink+"\"") ||
        !strings.Contains(markup, "autorisation revérifiée à la destination") {
        t.Fatalf("Case navigation missing or improperly encoded: %s", markup)
    }
    parsed, err := url.Parse(wantLink)
    if err != nil {
        t.Fatal(err)
    }
    if parsed.Path != "/case-timeline" || len(parsed.Query()) != 1 ||
        parsed.Query().Get("case_id") != "case-allowed&tenant=forged#x" {
        t.Fatalf("unsafe Case destination query: %v", parsed)
    }
    if strings.Contains(markup, "case-hidden") ||
        strings.Contains(markup, "case_id=case-hidden") ||
        len(links.calls) != 1 || len(cases.calls) != 2 {
        t.Fatalf("denied Case leaked or skipped independent owners: %d, %d, %s",
            len(links.calls), len(cases.calls), markup)
    }
    var projection Result
    jsonResult := testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
    if err := json.Unmarshal(jsonResult.Body.Bytes(), &projection); err != nil {
        t.Fatal(err)
    }
    if projection.SelectedLinkedCases == nil || len(projection.SelectedLinkedCases.Cases) != 1 {
        t.Fatalf("legacy read-only JSON changed: %+v", projection.SelectedLinkedCases)
    }

    // No selected Incident means the separately protected link/case owners
    // must not be called and no Case navigation must appear.
    before := len(links.calls)
    list := makeRequest("/command/work-queue?view=incidents")
    if list.Code != http.StatusOK || strings.Contains(list.Body.String(), "/case-timeline?") ||
        len(links.calls) != before {
        t.Fatalf("unselected Case link leaked or triggered source lookup: %d %s",
            list.Code, list.Body.String())
    }
    task := makeRequest("/command/work-queue?view=tasks&selected=Task:task-b")
    if task.Code != http.StatusOK || strings.Contains(task.Body.String(), "/case-timeline?") ||
        len(links.calls) != before {
        t.Fatalf("Task inherited linked Case navigation: %d %s",
            task.Code, task.Body.String())
    }
}

func TestLinkedCaseNavigationNoLinkOnUnknownOrRefusedSource(t *testing.T) {
    source := &linkedCasesHost{fakeOwner: fixture(), links: nil}
    handler := NewHTTPHandler(trustedHost(), source)
    req := httptest.NewRequest(http.MethodGet,
        "/command/work-queue?selected=Incident:inc-a", nil)
    req.Header.Set("Accept", "text/html")
    recorder := httptest.NewRecorder()
    handler.ServeHTTP(recorder, req)
    if recorder.Code != http.StatusOK ||
        strings.Contains(recorder.Body.String(), "/case-timeline?") ||
        !strings.Contains(recorder.Body.String(), "Relations inconnues") {
        t.Fatalf("unknown relation was turned into navigation: %d %s", recorder.Code, recorder.Body.String())
    }
}
