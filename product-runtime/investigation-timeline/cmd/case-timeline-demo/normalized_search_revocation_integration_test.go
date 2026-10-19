package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

// Regression at the actual first-party boundary: the selected Case can have
// one authorized matching event followed by a different imported event whose
// source access is revoked. A completed Search Job must never reveal the
// earlier hit, even when the revoked source would not match the query text.
func TestNormalizedSearchRejectsMidFlightSourceRevocationWithoutPartialJob(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	source, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil { t.Fatal(err) }
	path := writeJSONLEvents(t,
		"{\"id\":\"one\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"PowerShell observed by endpoint\"}\n"+
		"{\"id\":\"two\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"Routine service message that does not match\"}\n")
	raw, err := readLocalEventFile(path)
	if err != nil { t.Fatal(err) }
	data, err := importLocalEvents(context.Background(), source, "analyst-a", caseID, raw)
	if err != nil { t.Fatal(err) }
	imported := make(map[string]struct{})
	for _, entry := range data.Entries[len(source.Entries):] {
		imported[entry.Value.ID] = struct{}{}
	}
	if len(imported) != 2 { t.Fatalf("expected two distinct imported owner entries, got %d", len(imported)) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	payload, err := json.Marshal(map[string]string{
		"case_id": caseID, "value": "powershell",
		"from": "2026-10-17T09:59:00Z", "to": "2026-10-17T10:01:00Z",
	})
	if err != nil { t.Fatal(err) }
	invoke := func(h http.Handler, ctx context.Context) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/demo-local-search", bytes.NewReader(payload)).WithContext(ctx)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		got := httptest.NewRecorder()
		h.ServeHTTP(got, req)
		return got
	}
	// Positive control: the search really can find the earlier authorized
	// event before testing a revocation of the later benign source.
	positive := invoke(newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, imported, base),
		context.Background())
	var completed struct { Job search.SearchJob `json:"job"` }
	if positive.Code != http.StatusOK ||
		json.Unmarshal(positive.Body.Bytes(), &completed) != nil ||
		completed.Job.State != search.JobCompleted || len(completed.Job.ResultRefs) != 1 {
		t.Fatalf("positive first-party authorized search did not run: %d %s", positive.Code, positive.Body.String())
	}
	firstRef := completed.Job.ResultRefs[0].Ref
	for _, mode := range []string{"revoked-source", "request-cancelled"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			opens := 0
			owner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/demo-source-open" {
					opens++
					if opens == 2 {
						if mode == "request-cancelled" {
							cancel()
						} else {
							w.WriteHeader(http.StatusNotFound)
							return
						}
					}
				}
				base.ServeHTTP(w, r)
			})
			h := newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, imported, owner)
			got := invoke(h, ctx)
			if opens != 2 {
				t.Fatalf("search failed to reauthorize the second, nonmatching source: %d source opens", opens)
			}
			if got.Code == http.StatusOK ||
				strings.Contains(got.Body.String(), firstRef) ||
				strings.Contains(got.Body.String(), "PowerShell observed") ||
				strings.Contains(got.Body.String(), token) ||
				strings.Contains(got.Body.String(), "\"result_refs\"") {
				t.Fatalf("incomplete source permission yielded partial Search Job %d: %s", got.Code, got.Body.String())
			}
			if got.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("denied hunt response lost no-store policy")
			}
		})
	}
}
