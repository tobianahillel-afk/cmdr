package main

import (
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
