package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// The Case picker is a security boundary: invalid owner data must fail closed
// even when the host passes a syntactically valid synthetic principal.
func TestDemoCaseDiscoveryFailClosedOnInvalidOwner(t *testing.T) {
	auth := func(*http.Request) (string, error) { return "analyst-a", nil }
	for _, tc := range []struct {
		name string
		data timeline.MemoryTimelineData
	}{
		{"no-owner", timeline.MemoryTimelineData{}},
		{"duplicate-case", func() timeline.MemoryTimelineData {
			data := demoData()
			data.Cases = append(data.Cases, data.Cases[0])
			return data
		}()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newDemoCaseListHandler(auth, tc.data)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
			if rec.Code != http.StatusServiceUnavailable ||
				rec.Header().Get("Cache-Control") != "no-store" ||
				rec.Header().Get("Content-Type") != "application/json; charset=utf-8" ||
				strings.Contains(rec.Body.String(), "case-a") ||
				strings.Contains(rec.Body.String(), "demo-tenant") {
				t.Fatalf("invalid Case catalog disclosed details: %d %q", rec.Code, rec.Body.String())
			}
		})
	}
	h := newDemoCaseListHandler(auth, demoData())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "case-a") ||
		strings.Contains(rec.Body.String(), "case-b") {
		t.Fatalf("Case discovery lost tenant isolation: %d %s", rec.Code, rec.Body.String())
	}
}

// The JSONL exporter traverses authenticated pages before returning ANY bytes.
// If the owner revokes access after page 1, it must not serialize partial rows.
func TestCaseJSONLExportRejectsSecondPageRevocationWithoutPartialOutput(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	data, err := loadDemoSnapshot(owner)
	if err != nil {
		t.Fatal(err)
	}
	h, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	visitedSecondPage := false
	revoked := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/case-timeline" && r.URL.Query().Get("offset") == "200" {
			visitedSecondPage = true
			w.Header().Set("Cache-Control", "no-store")
			http.Error(w, "request-unavailable", http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
	payload, err := collectAuthorizedCaseJSONL(context.Background(), revoked, testBearer, "case-a")
	if !visitedSecondPage || err == nil || len(payload) != 0 {
		t.Fatalf("revoked Case continued export: visited=%v err=%v bytes=%d", visitedSecondPage, err, len(payload))
	}
	anonymous := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	if data, err := collectAuthorizedCaseJSONL(context.Background(), anonymous, "", "case-a"); err == nil || len(data) != 0 {
		t.Fatalf("anonymous JSONL export yielded bytes: err=%v bytes=%d", err, len(data))
	}
}

// A corrupt on-disk owner snapshot must not yield a bearer token, port or
// reconstructed Case even if the companion normalized JSONL is valid.
func TestLiveImportInvalidOwnerHasNoListenerOrOutput(t *testing.T) {
	events := writeJSONLEvents(t, importedJSONL)
	owner := filepath.Join(t.TempDir(), "owner.json")
	if err := os.WriteFile(owner, []byte("{broken json"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	args := []string{
		"-principal=analyst-a",
		"-snapshot=" + owner, "-import-jsonl=" + events,
		"-into-case=case-a", "-serve-import",
	}
	if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() > 0 {
		t.Fatalf("invalid owner started loopback or emitted data: err=%v output=%q", err, out.String())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := runDemoCLI(ctx, args, &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled loader disclosed data: err=%v output=%q", err, out.String())
	}
}
