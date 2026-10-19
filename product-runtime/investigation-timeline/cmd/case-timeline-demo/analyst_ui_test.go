package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestLocalAnalystViewerAssetsAreIsolated(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	for _, asset := range []struct {
		route, mime, contains string
	}{
		{"/analyst", "text/html", "Investigation Timeline"},
		{"/analyst.css", "text/css", "color-scheme"},
		{"/analyst.js", "text/javascript", "fetch(url"},
	} {
		t.Run(asset.route, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, asset.route, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, r)
			if response.Code != http.StatusOK ||
				!strings.HasPrefix(response.Header().Get("Content-Type"), asset.mime) ||
				!strings.Contains(response.Body.String(), asset.contains) {
				t.Fatalf("analyst asset unavailable %s: %d", asset.route, response.Code)
			}
			for _, header := range []string{"Cache-Control", "Content-Security-Policy", "Referrer-Policy", "X-Frame-Options", "X-Content-Type-Options"} {
				if response.Header().Get(header) == "" {
					t.Fatalf("missing analyst security header %s", header)
				}
			}
			if strings.Contains(response.Body.String(), testBearer) ||
				strings.Contains(response.Body.String(), "demo-hidden-source") ||
				strings.Contains(response.Body.String(), "demo-tenant-b") {
				t.Fatal("static asset disclosed fixture identity or credential")
			}
			for _, method := range []string{http.MethodHead, http.MethodPost, http.MethodDelete} {
				request := httptest.NewRequest(method, asset.route, nil)
				out := httptest.NewRecorder()
				handler.ServeHTTP(out, request)
				if method == http.MethodHead && (out.Code != http.StatusOK || out.Body.Len() != 0) {
					t.Fatalf("HEAD should not send body: %d %q", out.Code, out.Body.String())
				}
				if method != http.MethodHead && out.Code != http.StatusMethodNotAllowed {
					t.Fatalf("%s accepted on asset: %d", method, out.Code)
				}
			}
		})
	}
}

func TestAnalystViewerUsesAuthenticatedExistingTimeline(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	response, err := server.Client().Get(server.URL + "/analyst")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<16))
	response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK || !strings.Contains(string(body), "/analyst.js") {
		t.Fatalf("viewer not served: %d %v", response.StatusCode, err)
	}
	status, data := requestDemo(t, server, "/case-timeline?case_id=case-a&limit=20", "")
	if status != http.StatusUnauthorized || strings.Contains(data, "event-new") {
		t.Fatal("public viewer bypassed API authentication")
	}
	status, data = requestDemo(t, server, "/case-timeline?case_id=case-a&limit=20", testBearer)
	if status != http.StatusOK || !strings.Contains(data, "event-old") ||
		!strings.Contains(data, "event-new") || strings.Contains(data, "demo-hidden-source") {
		t.Fatalf("authorized viewer would receive wrong source data: %d %s", status, data)
	}
	status, data = requestDemo(t, server, "/case-timeline?case_id=case-b", testBearer)
	if status != http.StatusNotFound || strings.Contains(data, "event-b") {
		t.Fatal("viewer API crossed the tenant boundary")
	}
}

func TestAnalystViewerIsNotCrossOriginOrExecutableFixture(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/analyst", "/analyst.js"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Header().Get("Access-Control-Allow-Origin") != "" ||
			!strings.Contains(response.Header().Get("Content-Security-Policy"), "connect-src 'self'") {
			t.Fatal("viewer accidentally allowed cross-origin fetching")
		}
	}
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, prohibited := range []string{"innerHTML", "localStorage", "sessionStorage", "eval("} {
		if strings.Contains(string(js), prohibited) {
			t.Fatalf("analyst UI uses unsafe DOM or persistent token feature %s", prohibited)
		}
	}
	if !strings.Contains(string(js), "textContent") || !strings.Contains(string(js), "credentials: \"omit\"") {
		t.Fatal("missing safe text rendering or credential discipline")
	}
}

func TestAnalystViewerExposesBoundedTemporalNavigation(t *testing.T) {
	html, err := analystAssets.ReadFile("analyst.html")
	if err != nil { t.Fatal(err) }
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil { t.Fatal(err) }
	for _, mark := range []string{`id="from"`, `id="to"`, `type="datetime-local"`} {
		if !strings.Contains(string(html), mark) { t.Fatalf("time selector absent: %s", mark) }
	}
	for _, mark := range []string{"new Date(fromInput.value)", "new Date(toInput.value)",
		"from.getTime() >= to.getTime()", "toISOString()", `url.searchParams.set("from"`,
		`url.searchParams.set("to"`, "new AbortController()", "currentRevision !== revision",
		`field.addEventListener("input"`} {
		if !strings.Contains(string(js), mark) { t.Fatalf("missing temporal or stale-response behavior: %s", mark) }
	}
}

func TestAnalystCaseDiscoveryUiUsesSafeAuthorizedFlow(t *testing.T) {
	html, err := analystAssets.ReadFile("analyst.html")
	if err != nil { t.Fatal(err) }
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil { t.Fatal(err) }
	for _, marker := range []string{`id="discover"`, `id="case-options"`, `id="case-discovery-status"`} {
		if !strings.Contains(string(html), marker) { t.Fatalf("missing case chooser: %s", marker) }
	}
	for _, marker := range []string{`new URL("/demo-cases"`, `headers: { Authorization: "Bearer " + bearer }`,
		`credentials: "omit"`, `option.value = id`,
		`option.textContent = JSON.stringify(id).replace(/ /g, "\\u0020")`, `new AbortController()`,
		`currentRevision !== discoveryRevision`, `if (field === token)`,
		`caseID.dispatchEvent(new Event("input"`} {
		if !strings.Contains(string(js), marker) { t.Fatalf("missing case discovery security behavior: %s", marker) }
	}
}

func TestAnalystViewerSourceNavigationIsServerAuthorized(t *testing.T) {
	html, err := analystAssets.ReadFile("analyst.html")
	if err != nil { t.Fatal(err) }
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil { t.Fatal(err) }
	for _, marker := range []string{`id="source-detail"`, `id="source-status"`, `id="return-to-case"`} {
		if !strings.Contains(string(html), marker) { t.Fatalf("source navigation UI missing %s", marker) }
	}
	for _, marker := range []string{
		`new URL("/case-source-navigation"`, `url.searchParams.set("entry_id", entryID)`,
		`headers: { Authorization: "Bearer " + bearer }`,
		`data.return_is_authority !== false`,
		`data.return_case_id !== selectedCase`,
		`data.return_entry_id !== entryID`,
		`sourceDetail.replaceChildren()`,
		`returnToCase.addEventListener("click"`,
		`clearSource();`,
		`fetchPage(); // reauthorize the Case and Timeline`,
		`currentRevision !== sourceRevision`,
	} {
		if !strings.Contains(string(js), marker) { t.Fatalf("safe source flow missing %s", marker) }
	}
}

func TestAnalystReturnContextIsPresentationOnly(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil { t.Fatal(err) }
	server := httptest.NewServer(handler)
	defer server.Close()
	req := "/case-source-navigation?case_id=case-a&entry_id=event-new&limit=20&offset=0&kind=observed"
	status, data := requestDemo(t, server, req, testBearer)
	if status != http.StatusOK { t.Fatalf("cannot navigate: %d %s", status, data) }
	var draft struct {
		ReturnQuery string `json:"return_query"`
		ReturnIsAuthority bool `json:"return_is_authority"`
	}
	if err := json.Unmarshal([]byte(data), &draft); err != nil { t.Fatal(err) }
	if draft.ReturnIsAuthority || draft.ReturnQuery == "" || strings.Contains(draft.ReturnQuery, "entry_id") ||
		strings.Contains(draft.ReturnQuery, "Bearer") {
		t.Fatalf("unsafe return context: %+v", draft)
	}
	status, data = requestDemo(t, server, "/analyst?" + draft.ReturnQuery, "")
	if status != http.StatusOK || !strings.Contains(data, "Investigation Timeline") ||
		strings.Contains(data, testBearer) { t.Fatalf("return UI not inert: %d %s", status, data) }
	status, data = requestDemo(t, server, "/case-timeline?" + draft.ReturnQuery, "")
	if status != http.StatusUnauthorized || strings.Contains(data, "event-new") {
		t.Fatalf("return context borrowed permissions: %d %s", status, data)
	}
	status, data = requestDemo(t, server, "/case-timeline?" + draft.ReturnQuery, testBearer)
	if status != http.StatusOK || !strings.Contains(data, "event-new") {
		t.Fatalf("authenticated return failed: %d %s", status, data)
	}
}

func TestAnalystReturnContextControlsNeverStoreBearer(t *testing.T) {
	html, err := analystAssets.ReadFile("analyst.html")
	if err != nil { t.Fatal(err) }
	js, err := analystAssets.ReadFile("analyst.js")
	if err != nil { t.Fatal(err) }
	for _, marker := range []string{`id="copy-return-link"`, "Le lien contient les critères de recherche"} {
		if !strings.Contains(string(html), marker) { t.Fatalf("missing explicit return control %q", marker) }
	}
	for _, marker := range []string{
		"function restoreReturnSelection()", "const allowed = new Set(",
		"restoreOffsetOnSubmit", "if (!returnLinkURL) return",
		"navigator.clipboard.writeText(returnLinkURL)", "new URL(\"/analyst\", window.location.origin)",
		"expected.delete(\"entry_id\")", "received.size !== expected.size",
		"copyReturnLink.disabled = true", "restoreReturnSelection();",
	} {
		if !strings.Contains(string(js), marker) { t.Fatalf("unsafe or incomplete return selection: %q", marker) }
	}
	for _, prohibited := range []string{"localStorage", "sessionStorage", "document.cookie", "history.replaceState"} {
		if strings.Contains(string(js), prohibited) { t.Fatalf("persistent authority state introduced: %s", prohibited) }
	}
}


func TestExactOwnerCaseIDSurvivesRealPickerTimelineSearchAndSourceOpen(t *testing.T) {
	// This is one complete first-party owner-authorized contract journey,
	// not a second UI parser: Case IDs are opaque, including outer spaces.
	const selected = " forensic-case-on-disk "
	const normal = "forensic-case-on-disk"
	data, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil { t.Fatal(err) }
	found := false
	for _, c := range data.Cases {
		if c.ID != normal { continue }
		copyCase := c
		copyCase.ID = selected
		copyCase.EntryIDs = nil
		data.Cases = append(data.Cases, copyCase)
		found = true
		break
	}
	if !found { t.Fatal("fixture lacks first Case") }
	// Import into ONLY the whitespace-distinct Case. If any layer silently
	// trims the selected ID, the Search Job must NOT find these events.
	raw := []byte("{\"id\":\"opaque-owner-identity\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"Unusual PowerShell execution\"}\\n")
	importedData, err := importLocalEvents(context.Background(), data, "analyst-a", selected, raw)
	if err != nil { t.Fatalf("selected exact-Case import failed: %v", err) }
	ids := make(map[string]struct{})
	for _, entry := range importedData.Entries {
		if entry.CaseID == selected && entry.Value.Provenance == "operator-local-jsonl" {
			ids[entry.Value.ID] = struct{}{}
		}
	}
	if len(ids) != 1 { t.Fatalf("expected one isolated imported source: %d", len(ids)) }
	owner, err := newDemoHandlerWithData(testBearer, "analyst-a", importedData)
	if err != nil { t.Fatalf("owner rejected distinct Case IDs: %v", err) }
	handler := newDemoLocalJSONLSearchHandler(testBearer, "analyst-a", selected, importedData, ids, owner)
	server := httptest.NewServer(handler)
	defer server.Close()

	status, catalog := requestDemo(t, server, "/demo-cases", testBearer)
	var choices struct { Cases []string `json:"cases"` }
	if status != http.StatusOK || json.Unmarshal([]byte(catalog), &choices) != nil {
		t.Fatalf("owner-authorized picker unavailable: %d %s", status, catalog)
	}
	seen := make(map[string]bool)
	for _, id := range choices.Cases { seen[id] = true }
	if !seen[normal] || !seen[selected] {
		t.Fatalf("picker aliased or removed authorized opaque ID: %+v", choices.Cases)
	}
	for _, id := range []string{normal, selected} {
		code, timeline := requestDemo(t, server, "/case-timeline?case_id="+url.QueryEscape(id), testBearer)
		if code != http.StatusOK {
			t.Fatalf("owner Timeline refused selected valid Case %q: %d", id, code)
		}
		hasImported := strings.Contains(timeline, "local-jsonl-entry-")
		if hasImported != (id == selected) {
			t.Fatalf("Timeline crossed exact Case identity %q, event present=%v", id, hasImported)
		}
	}
	query := func(caseID string) *httptest.ResponseRecorder {
		t.Helper()
		body, _ := json.Marshal(map[string]string{
			"case_id": caseID, "value": "powershell",
			"from": "2026-10-17T09:59:00Z",
			"to": "2026-10-17T10:01:00Z",
		})
		req := httptest.NewRequest(http.MethodPost, "/demo-local-search", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testBearer)
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		return out
	}
	out := query(selected)
	var job struct {
		CaseID string `json:"case_id"`
		Job search.SearchJob `json:"job"`
	}
	if out.Code != http.StatusOK || json.Unmarshal(out.Body.Bytes(), &job) != nil ||
		job.CaseID != selected || job.Job.State != search.JobCompleted ||
		len(job.Job.ResultRefs) != 1 {
		t.Fatalf("exact selected Case search did not complete: %d %s", out.Code, out.Body.String())
	}
	ref := job.Job.ResultRefs[0].Ref
	if _, ok := ids[ref]; !ok || strings.Contains(out.Body.String(), "Unusual PowerShell execution") {
		t.Fatalf("Search Job leaked unauthorized source or raw preview: %+v", job.Job)
	}
	if denied := query(normal); denied.Code == http.StatusOK {
		t.Fatal("Search Job silently normalized selected Case to a different accessible Case")
	}
	source := httptest.NewRequest(http.MethodGet,
		"/demo-source-open?case_id="+url.QueryEscape(selected)+"&entry_id="+url.QueryEscape(ref), nil)
	source.Header.Set("Authorization", "Bearer "+testBearer)
	opened := httptest.NewRecorder()
	handler.ServeHTTP(opened, source)
	if opened.Code != http.StatusOK {
		t.Fatalf("exact selected Case source could not be reauthorized: %d", opened.Code)
	}
	source.URL.RawQuery = "case_id="+url.QueryEscape(normal)+"&entry_id="+url.QueryEscape(ref)
	notOwned := httptest.NewRecorder()
	handler.ServeHTTP(notOwned, source)
	if notOwned.Code == http.StatusOK || strings.Contains(notOwned.Body.String(), "Unusual PowerShell") {
		t.Fatal("wrong Case selector reopened source from another authorized Case")
	}
}
