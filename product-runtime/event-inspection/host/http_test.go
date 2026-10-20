package inspectionhost

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

type testSession struct {
	scope Session
	err error
	calls int
}
func (s *testSession) ResolveInspectionSession(r *http.Request) (context.Context, Session, error) {
	s.calls++
	return r.Context(), s.scope, s.err
}

type testSearch struct {
	allowed bool
	err error
	calls int
	last Request
}
func (s *testSearch) CanInspectSearchResult(_ context.Context, scope Session, req Request) (bool, error) {
	s.calls++
	s.last = req
	if scope.Principal != "analyst" { return false, nil }
	return s.allowed, s.err
}

type testEvent struct {
	owned AuthorizedEvent
	err error
	calls int
}
func (e *testEvent) ReadAuthorizedInspectionEvent(_ context.Context, scope Session, req Request) (AuthorizedEvent, error) {
	e.calls++
	if scope.Tenant != req.TenantRef { return AuthorizedEvent{}, errors.New("scope mismatch") }
	return e.owned, e.err
}

func standardEvent() inspection.Event {
	return inspection.Event{
		Ref: "event-1", TenantRef: "tenant-a", State: inspection.StatePartial,
		Source: inspection.SourceMetadata{SourceRef: "edr", DataSourceRef: "endpoint-a", Kind: "telemetry"},
		Parser: inspection.ParserMetadata{ParserRef: "parser-a", Version: "v1", Status: "degraded"},
		SourceFields: map[string]string{"source.ip": "192.0.2.4"},
		NormalizedFields: map[string]string{"network.client.ip": "192.0.2.4"},
		MissingFields: []string{"network.server.ip"},
		RawPayload: []byte("RAW-PRIVATE-STRING"), RenderedPayload: []byte("rendered authorized"),
		RawDerivedSensitive: map[string]string{"secret.token": "RAW-PRIVATE-STRING"},
		Enrichments: []inspection.EventEnrichment{{
			Producer: "detector-a", Version: "v3", Freshness: "stale",
			Values: map[string]string{"risk.score": "72"},
		}},
	}
}

func standardOwners() (*testSession, *testSearch, *testEvent) {
	return &testSession{scope: Session{Principal: "analyst", Tenant: "tenant-a", Environment: "prod"}},
		&testSearch{allowed: true},
		&testEvent{owned: AuthorizedEvent{
			Allowed: true, Event: standardEvent(),
			RawDecision: inspection.AccessDeny, RenderDecision: inspection.AccessAllow,
		}}
}

func standardRequest() Request {
	return Request{
		TenantRef: "tenant-a", EnvironmentRef: "prod", EventRef: "event-1",
		JobID: "job-1", RunID: "run-1", CorrelationID: "corr-1", SourceHint: "edr",
	}
}
func doRequest(t *testing.T, h http.Handler, query Request) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(query)
	if err != nil { t.Fatal(err) }
	return doRaw(h, http.MethodPost, route, "application/json", data)
}
func doRaw(h http.Handler, method, path, contentType string, data []byte) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, bytes.NewReader(data))
	if contentType != "" { r.Header.Set("Content-Type", contentType) }
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func newTestHandler(t *testing.T, s *testSession, search *testSearch, event *testEvent) http.Handler {
	t.Helper()
	h, err := NewHandler(s, search, event)
	if err != nil { t.Fatal(err) }
	return h
}
func parsed(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil { t.Fatal(err) }
	return out
}

func TestAuthorizedHTTPProjectionIsBrowserCompatibleAndRawDenied(t *testing.T) {
	s, search, event := standardOwners()
	w := doRequest(t, newTestHandler(t, s, search, event), standardRequest())
	if w.Code != http.StatusOK { t.Fatalf("status=%d body=%s", w.Code, w.Body.String()) }
	if s.calls != 1 || search.calls != 1 || event.calls != 1 { t.Fatal("independent authorization chain was bypassed") }
	if search.last.JobID != "job-1" || search.last.RunID != "run-1" || search.last.SourceHint != "edr" { t.Fatal(search.last) }
	j := parsed(t, w.Body.Bytes())
	if j["kind"] != "inspection" || j["eventRef"] != "event-1" ||
		j["tenantRef"] != "tenant-a" || j["environmentRef"] != "prod" { t.Fatalf("bad event identity: %+v", j) }
	source := j["source"].(map[string]any)
	parser := j["parser"].(map[string]any)
	if source["sourceRef"] != "edr" || parser["parserRef"] != "parser-a" ||
		parser["version"] != "v1" { t.Fatalf("bad provenance: %+v %+v", source, parser) }
	if j["raw"].(map[string]any)["access"] != "denied" ||
		j["rendered"].(map[string]any)["text"] != "rendered authorized" { t.Fatal("raw/rendered authorization invalid") }
	if strings.Contains(w.Body.String(), "RAW-PRIVATE-STRING") { t.Fatal("raw denial leaked classified content") }
	enrichment := j["enrichments"].([]any)[0].(map[string]any)
	if enrichment["producer"] != "detector-a" || enrichment["stale"] != true ||
		enrichment["values"].(map[string]any)["risk.score"] != "72" { t.Fatal(enrichment) }
	if j["missingFields"].([]any)[0] != "network.server.ip" { t.Fatal("missing-field provenance") }
	if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("X-Content-Type-Options") != "nosniff" ||
		w.Header().Get("Access-Control-Allow-Origin") != "" { t.Fatal("unsafe response headers") }
}

func TestIndependentRawOnlyDecisionAndSourceData(t *testing.T) {
	s, search, event := standardOwners()
	event.owned.RawDecision = inspection.AccessAllow
	event.owned.RenderDecision = inspection.AccessDeny
	w := doRequest(t, newTestHandler(t, s, search, event), standardRequest())
	if w.Code != http.StatusOK { t.Fatal(w.Code, w.Body.String()) }
	j := parsed(t, w.Body.Bytes())
	if j["raw"].(map[string]any)["text"] != "RAW-PRIVATE-STRING" ||
		j["rendered"].(map[string]any)["access"] != "denied" { t.Fatal(j) }
	if j["rawDerivedSensitive"].(map[string]any)["secret.token"] != "RAW-PRIVATE-STRING" { t.Fatal("allowed raw lost") }
}

func TestDeniedSearchNeverCallsTelemetryAndDenialsAreGeneric(t *testing.T) {
	for _, scenario := range []string{"search-deny", "search-error", "event-deny", "event-error",
		"wrong-event-ref", "foreign-event-tenant", "both-payloads-denied"} {
		t.Run(scenario, func(t *testing.T) {
			s, search, e := standardOwners()
			switch scenario {
			case "search-deny": search.allowed = false
			case "search-error": search.err = errors.New("confidential owner error")
			case "event-deny": e.owned.Allowed = false
			case "event-error": e.err = errors.New("secret protected owner error")
			case "wrong-event-ref": e.owned.Event.Ref = "event-foreign"
			case "foreign-event-tenant": e.owned.Event.TenantRef = "tenant-b"
			case "both-payloads-denied": e.owned.RenderDecision = inspection.AccessDeny
			}
			w := doRequest(t, newTestHandler(t, s, search, e), standardRequest())
			if w.Code != http.StatusNotFound || strings.Contains(w.Body.String(), "secret") ||
				strings.Contains(w.Body.String(), "event-foreign") { t.Fatalf("%s: %d %s", scenario, w.Code, w.Body.String()) }
			if strings.HasPrefix(scenario, "search-") && e.calls != 0 { t.Fatal("telemetry consulted before Search authorization") }
		})
	}
}

func TestSessionScopeIsNeverTakenFromClient(t *testing.T) {
	s, search, e := standardOwners()
	h := newTestHandler(t, s, search, e)
	q := standardRequest()
	q.TenantRef = "tenant-b"
	w := doRequest(t, h, q)
	if w.Code != http.StatusNotFound || search.calls != 0 || e.calls != 0 { t.Fatal(w.Code) }
	s.scope.Principal = ""
	w = doRequest(t, h, standardRequest())
	if w.Code != http.StatusNotFound || search.calls != 0 { t.Fatal("unauthenticated scope allowed") }
	s.scope.Principal = "analyst"
	s.err = errors.New("token expired")
	w = doRequest(t, h, standardRequest())
	if w.Code != http.StatusNotFound || search.calls != 0 { t.Fatal("session error leaked") }
}

func TestUnknownUnboundedAndNonPostRequestsAreRejectedPreAuth(t *testing.T) {
	s, search, e := standardOwners()
	h := newTestHandler(t, s, search, e)
	cases := []struct{method, path, media string; data []byte; status int}{
		{http.MethodGet, route, "application/json", nil, http.StatusMethodNotAllowed},
		{http.MethodPost, route + "?tenant=tenant-b", "application/json", []byte("{}"), http.StatusBadRequest},
		{http.MethodPost, "/wrong", "application/json", []byte("{}"), http.StatusBadRequest},
		{http.MethodPost, route, "text/plain", []byte("{}"), http.StatusUnsupportedMediaType},
		{http.MethodPost, route, "application/json", []byte(`{"tenantRef":"tenant-a","unknown":"secret"}`), http.StatusBadRequest},
		{http.MethodPost, route, "application/json", []byte("[]"), http.StatusBadRequest},
		{http.MethodPost, route, "application/json", bytes.Repeat([]byte("a"), maxRequestBytes+1), http.StatusBadRequest},
		{http.MethodPost, route, "application/json", []byte(`{} {}`), http.StatusBadRequest},
	}
	for i, c := range cases {
		w := doRaw(h, c.method, c.path, c.media, c.data)
		if w.Code != c.status { t.Errorf("case %d: got %d want %d", i, w.Code, c.status) }
	}
	if s.calls != 0 || search.calls != 0 || e.calls != 0 { t.Fatal("bad request reached authority") }
}

func TestInvalidProviderFieldsAndPayloadCannotEscapeInError(t *testing.T) {
	for _, scenario := range []string{"unsafe-source", "prototype-field", "huge-field", "invalid-utf8",
		"too-large-raw", "derived-as-source", "missing-parser-provenance"} {
		t.Run(scenario, func(t *testing.T) {
			s, search, e := standardOwners()
			switch scenario {
			case "unsafe-source": e.owned.Event.Source.SourceRef = "unsafe source"
			case "prototype-field": e.owned.Event.SourceFields["__proto__"] = "hidden"
			case "huge-field": e.owned.Event.NormalizedFields["field.x"] = strings.Repeat("x", 4097)
			case "invalid-utf8": e.owned.Event.RenderedPayload = []byte{0xff}
			case "too-large-raw":
				e.owned.RawDecision = inspection.AccessAllow
				e.owned.Event.RawPayload = bytes.Repeat([]byte("x"), 65537)
			case "derived-as-source": e.owned.Event.Enrichments[0].PresentedAsSource = true
			case "missing-parser-provenance": e.owned.Event.Enrichments[0].Producer = ""
			}
			w := doRequest(t, newTestHandler(t, s, search, e), standardRequest())
			if w.Code != http.StatusNotFound || strings.Contains(w.Body.String(), "RAW-PRIVATE-STRING") {
				t.Fatalf("%s unexpected status=%d body=%s", scenario, w.Code, w.Body.String())
			}
		})
	}
}

func TestNewHandlerRejectsMissingAndTypedNilServices(t *testing.T) {
	s, search, e := standardOwners()
	if _, err := NewHandler(nil, search, e); err == nil { t.Fatal("accepted missing session") }
	if _, err := NewHandler(s, nil, e); err == nil { t.Fatal("accepted missing search") }
	if _, err := NewHandler(s, search, nil); err == nil { t.Fatal("accepted missing Event owner") }
	var typedNil *testEvent
	if _, err := NewHandler(s, search, typedNil); err == nil { t.Fatal("accepted typed nil") }
}

func TestMalformedAndRevokedSourceNeverPublishesPreviousData(t *testing.T) {
	s, search, e := standardOwners()
	h := newTestHandler(t, s, search, e)
	if doRequest(t, h, standardRequest()).Code != http.StatusOK { t.Fatal("expected initial readable projection") }
	search.allowed = false
	w := doRequest(t, h, standardRequest())
	if w.Code != http.StatusNotFound || strings.Contains(w.Body.String(), "192.0.2.4") { t.Fatal(w.Code, w.Body.String()) }
}
