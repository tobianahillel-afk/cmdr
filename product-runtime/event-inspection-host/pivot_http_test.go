package inspectionhost

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/http"
	"strings"
	"testing"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

func pivotTestBody(t *testing.T, alter func(map[string]any)) []byte {
	t.Helper()
	q := map[string]any{
		"tenantRef": "tenant-a", "environmentRef": "prod", "eventRef": "event-1",
		"jobId": "job-1", "runId": "run-1", "correlationId": "corr-1", "sourceHint": "edr",
		"field": "source.ip", "value": "192.0.2.4", "valuePresent": true,
		"timeStart": "2026-10-20T09:00:00Z", "timeEnd": "2026-10-20T12:00:00Z",
	}
	if alter != nil {
		alter(q)
	}
	data, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func pivotRequest(t *testing.T, h http.Handler, alter func(map[string]any)) (int, string) {
	t.Helper()
	w := doRaw(h, http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, alter))
	return w.Code, w.Body.String()
}

func TestHTTPPivotDraftIndependentlyReauthorizesSelectedSearchAndEvent(t *testing.T) {
	s, search, event := standardOwners()
	h := newTestHandler(t, s, search, event)
	status, body := pivotRequest(t, h, nil)
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s", status, body)
	}
	draft := parsed(t, []byte(body))
	if draft["kind"] != "pivot-draft" || draft["targetCapability"] != "CAP-INV-002" ||
		draft["field"] != "source.ip" || draft["value"] != "192.0.2.4" ||
		draft["fieldOrigin"] != "source" || draft["jobId"] != "job-1" ||
		draft["runId"] != "run-1" || draft["auditRequired"] != true {
		t.Fatalf("pivot draft lost provenance: %v", draft)
	}
	back := draft["returnContext"].(map[string]any)
	if back["origin"] != "event-search" || back["eventRef"] != "event-1" ||
		back["environmentRef"] != "prod" {
		t.Fatalf("return origin not server-derived: %v", back)
	}
	if s.calls != 1 || search.calls != 1 || event.calls != 1 {
		t.Fatalf("owner calls session=%d search=%d telemetry=%d", s.calls, search.calls, event.calls)
	}
	if strings.Contains(body, "RAW-PRIVATE-STRING") || strings.Contains(body, "secret.token") {
		t.Fatal("raw-denied data appeared in draft")
	}
	// Existing Inspector response remains unchanged on this same handler.
	if w := doRequest(t, h, standardRequest()); w.Code != http.StatusOK {
		t.Fatalf("prior Inspector route regressed: %d %s", w.Code, w.Body.String())
	}
}

func TestHTTPPivotNeverPromotesMissingEmptyValuedField(t *testing.T) {
	s, search, event := standardOwners()
	h := newTestHandler(t, s, search, event)
	status, body := pivotRequest(t, h, func(q map[string]any) {
		q["field"] = "missing.empty"
		q["value"] = ""
	})
	if status != http.StatusNotFound || strings.Contains(body, "fieldOrigin") {
		t.Fatalf("missing field pivot accepted: %d %s", status, body)
	}
	// Genuine field presence with empty value is still legal.
	event.owned.Event.NormalizedFields["genuinely.empty"] = ""
	status, body = pivotRequest(t, h, func(q map[string]any) {
		q["field"] = "genuinely.empty"
		q["value"] = ""
	})
	if status != http.StatusOK || parsed(t, []byte(body))["fieldOrigin"] != "normalized" {
		t.Fatalf("present empty field rejected: %d %s", status, body)
	}
}

func TestHTTPPivotNeverCopiesDeniedRawDerivedOrForgedScope(t *testing.T) {
	s, search, event := standardOwners()
	h := newTestHandler(t, s, search, event)
	cases := []struct {
		name string
		alter func(map[string]any)
	}{
		{"raw-restricted", func(q map[string]any) {
			q["field"], q["value"] = "secret.token", "RAW-PRIVATE-STRING"
		}},
		{"value-does-not-match", func(q map[string]any) { q["value"] = "203.0.113.3" }},
		{"wrong-tenant", func(q map[string]any) { q["tenantRef"] = "tenant-b" }},
		{"wrong-environment", func(q map[string]any) { q["environmentRef"] = "dev" }},
		{"search-job-denied", func(q map[string]any) { q["jobId"] = "non-member" }},
		{"missing-presence-bit", func(q map[string]any) { q["valuePresent"] = false }},
		{"unknown-return-authority", func(q map[string]any) { q["returnOrigin"] = "evil" }},
		{"invalid-time-window", func(q map[string]any) { q["timeStart"] = "2030-01-01T00:00:00Z" }},
		{"class2-mutation", func(q map[string]any) { q["mutationRequested"] = "case-link" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The default synthetic searchOwner grants any job, so only this
			// branch injects a hostile independent owner refusal.
			if tc.name == "search-job-denied" {
				search.allowed = false
				defer func() { search.allowed = true }()
			}
			status, body := pivotRequest(t, h, tc.alter)
			if status == http.StatusOK || strings.Contains(body, "RAW-PRIVATE-STRING") ||
				strings.Contains(body, "192.0.2.4") {
				t.Fatalf("denied pivot leaked or succeeded: %d %s", status, body)
			}
		})
	}
	if s.calls == 0 || event.calls == 0 {
		t.Fatal("authorized test paths did not exercise owner chain")
	}
}

func TestHTTPPivotClassifiedFieldsAndSourceProvenance(t *testing.T) {
	s, search, event := standardOwners()
	h := newTestHandler(t, s, search, event)
	status, body := pivotRequest(t, h, func(q map[string]any) {
		q["field"], q["value"] = "risk.score", "72"
	})
	if status != http.StatusOK || parsed(t, []byte(body))["fieldOrigin"] != "enrichment:detector-a@v3" {
		t.Fatalf("enrichment origin lost: %d %s", status, body)
	}
	// The same sensitive field may be selected ONLY after trusted owner
	// authorizes its raw projection, never because the request says so.
	event.owned.RawDecision = inspection.AccessAllow
	status, body = pivotRequest(t, h, func(q map[string]any) {
		q["field"], q["value"] = "secret.token", "RAW-PRIVATE-STRING"
	})
	if status != http.StatusOK || parsed(t, []byte(body))["fieldOrigin"] != "raw-derived-sensitive" {
		t.Fatalf("raw owner-approved field denied: %d %s", status, body)
	}
}

func TestHTTPPivotMalformedAndCrossScopeRequestsFailWithoutOwnerOrSensitiveBytes(t *testing.T) {
	cases := []struct {
		name string
		method string
		path string
		contentType string
		body []byte
		want int
	}{
		{"get-denied", http.MethodGet, pivotRoute, "application/json", pivotTestBody(t, nil), http.StatusMethodNotAllowed},
		{"wrong-content-type", http.MethodPost, pivotRoute, "text/plain", pivotTestBody(t, nil), http.StatusUnsupportedMediaType},
		{"missing-content-type", http.MethodPost, pivotRoute, "", pivotTestBody(t, nil), http.StatusUnsupportedMediaType},
		{"malformed-media-type", http.MethodPost, pivotRoute, "application/json; bad", pivotTestBody(t, nil), http.StatusUnsupportedMediaType},
		{"query-not-allowed", http.MethodPost, pivotRoute + "?trusted=true", "application/json", pivotTestBody(t, nil), http.StatusBadRequest},
		{"invalid-json", http.MethodPost, pivotRoute, "application/json", []byte("{not-json"), http.StatusBadRequest},
		{"empty-json", http.MethodPost, pivotRoute, "application/json", nil, http.StatusBadRequest},
		{"json-trailing-data", http.MethodPost, pivotRoute, "application/json", append(pivotTestBody(t, nil), []byte(" false")...), http.StatusBadRequest},
		{"too-large", http.MethodPost, pivotRoute, "application/json", []byte(strings.Repeat("x", maxRequestBytes+1)), http.StatusBadRequest},
		{"wrong-tenant-pattern", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { q["tenantRef"] = "*" }), http.StatusBadRequest},
		{"bad-event-reference", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { q["eventRef"] = "/bad+ref" }), http.StatusBadRequest},
		{"bad-source-hint", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { q["sourceHint"] = "bad input" }), http.StatusBadRequest},
		{"bad-field-pattern", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { q["field"] = "__proto__" }), http.StatusBadRequest},
		{"field-omitted", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { delete(q, "field") }), http.StatusBadRequest},
		{"unknown-source-field", http.MethodPost, pivotRoute, "application/json", pivotTestBody(t, func(q map[string]any) { q["permissions"] = []string{"admin"} }), http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, search, event := standardOwners()
			h := newTestHandler(t, s, search, event)
			w := doRaw(h, tc.method, tc.path, tc.contentType, tc.body)
			if w.Code != tc.want || s.calls != 0 || search.calls != 0 || event.calls != 0 ||
				strings.Contains(w.Body.String(), "RAW-PRIVATE-STRING") {
				t.Fatalf("bad input crossed authority: code=%d want=%d session=%d search=%d telemetry=%d body=%s",
					w.Code, tc.want, s.calls, search.calls, event.calls, w.Body.String())
			}
		})
	}
}

func TestHTTPPivotCancellationAndOwnerFailureNeverReturnDraft(t *testing.T) {
	t.Run("canceled-before-authority", func(t *testing.T) {
		s, search, event := standardOwners()
		h := newTestHandler(t, s, search, event)
		r := httptest.NewRequest(http.MethodPost, pivotRoute, bytes.NewReader(pivotTestBody(t, nil)))
		r.Header.Set("Content-Type", "application/json")
		ctx, cancel := context.WithCancel(r.Context())
		cancel()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r.WithContext(ctx))
		if w.Code != http.StatusServiceUnavailable || s.calls != 0 || search.calls != 0 || event.calls != 0 {
			t.Fatalf("canceled scope traversed authority: %d", w.Code)
		}
	})
	t.Run("no-body", func(t *testing.T) {
		s, search, event := standardOwners()
		h := newTestHandler(t, s, search, event)
		r := httptest.NewRequest(http.MethodPost, pivotRoute, nil)
		r.Header.Set("Content-Type", "application/json")
		r.Body = nil
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest || s.calls != 0 || search.calls != 0 || event.calls != 0 {
			t.Fatalf("nil body crossed authority: %d", w.Code)
		}
	})
	for _, tc := range []struct {
		name string
		breakSource func(*testSession, *testSearch, *testEvent)
	}{
		{"denied-session", func(s *testSession, _ *testSearch, _ *testEvent) { s.err = errors.New("secret-session-info") }},
		{"denied-search", func(_ *testSession, s *testSearch, _ *testEvent) { s.allowed = false }},
		{"errored-search", func(_ *testSession, s *testSearch, _ *testEvent) { s.err = errors.New("secret-search-info") }},
		{"denied-telemetry", func(_ *testSession, _ *testSearch, e *testEvent) { e.owned.Allowed = false }},
		{"errored-telemetry", func(_ *testSession, _ *testSearch, e *testEvent) { e.err = errors.New("secret-telemetry-info") }},
		{"wrong-event-ref", func(_ *testSession, _ *testSearch, e *testEvent) { e.owned.Event.Ref = "other-event" }},
		{"both-variants-denied", func(_ *testSession, _ *testSearch, e *testEvent) { e.owned.RenderDecision = inspection.AccessDeny }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, search, event := standardOwners()
			tc.breakSource(s, search, event)
			status, body := pivotRequest(t, newTestHandler(t, s, search, event), nil)
			if status != http.StatusNotFound || strings.Contains(body, "secret") ||
				strings.Contains(body, "192.0.2.4") || strings.Contains(body, "draft") {
				t.Fatalf("source failure exposed information: %d %s", status, body)
			}
		})
	}
}
