package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
