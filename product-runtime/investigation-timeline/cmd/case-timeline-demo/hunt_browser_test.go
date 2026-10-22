package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

// Browser behavior is checked with the same standalone Node harness used by
// the existing analyst tests; the Go workflow executes it with go test -race.
func TestAnalystCrossCaseHuntBrowser(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("Node.js unavailable: browser harness not executed")
	}
	cmd := exec.Command("node", "--test", "hunt.test.cjs")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cross-Case hunt browser tests failed: %v\n%s", err, output)
	}
}

func TestDemoServesCrossCaseHuntAssets(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/analyst", "/hunt.js"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("hunt asset unavailable: %s %d", path, w.Code)
		}
		if path == "/analyst" && !strings.Contains(w.Body.String(), "id=\"hunt-run\"") {
			t.Fatal("hunt control missing from analyst page")
		}
		if path == "/hunt.js" && !strings.Contains(w.Body.String(), "/demo-cases") {
			t.Fatal("hunt must rely on authorized server Case discovery")
		}
	}
	req := httptest.NewRequest(http.MethodPost, "/hunt.js", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatal("hunt static asset accepted mutation request")
	}
}
