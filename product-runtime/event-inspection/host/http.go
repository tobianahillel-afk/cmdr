// Package inspectionhost binds the existing Event Inspection projector to
// separately authorized, read-only HTTP requests. The embedding server owns
// authentication, TLS, telemetry access, source redaction and durable audit.
package inspectionhost

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

const route = "/investigate/event-inspection"
const maxRequestBytes = 8192

var opaqueRef = regexp.MustCompile(`^[A-Za-z0-9_.:/-]{1,256}$`)

// Session is populated by trusted host authentication, never request JSON.
type Session struct {
	Principal   string
	Tenant      string
	Environment string
}

type SessionResolver interface {
	ResolveInspectionSession(*http.Request) (context.Context, Session, error)
}

// Request is an untrusted selection, not an access grant.
type Request struct {
	TenantRef      string `json:"tenantRef"`
	EnvironmentRef string `json:"environmentRef"`
	EventRef       string `json:"eventRef"`
	JobID          string `json:"jobId"`
	RunID          string `json:"runId"`
	CorrelationID  string `json:"correlationId"`
	SourceHint     string `json:"sourceHint"`
}

// SearchOwner must independently re-authorize the selected Event's presence in
// this principal's Search Job/run at the specified tenant and environment.
type SearchOwner interface {
	CanInspectSearchResult(context.Context, Session, Request) (bool, error)
}

// AuthorizedEvent is a server-supplied result of independent telemetry-event
// permission, source and per-field classification. No client access decision
// or browser-provided source fields are consulted.
type AuthorizedEvent struct {
	Allowed        bool
	Event          inspection.Event
	RawDecision    inspection.AccessDecision
	RenderDecision inspection.AccessDecision
}

// TelemetryOwner must re-authorize Event.read and raw/rendered independently,
// resolve only the exact Event, and redact classified source/normalized/
// enriched fields before returning them. It must enforce environment scope.
type TelemetryOwner interface {
	ReadAuthorizedInspectionEvent(context.Context, Session, Request) (AuthorizedEvent, error)
}

type handler struct {
	sessions SessionResolver
	search   SearchOwner
	events   TelemetryOwner
}

func NewHandler(sessions SessionResolver, search SearchOwner, events TelemetryOwner) (http.Handler, error) {
	if nilAuthority(sessions) || nilAuthority(search) || nilAuthority(events) {
		return nil, errors.New("inspection host requires independent trusted authorities")
	}
	return &handler{sessions: sessions, search: search, events: events}, nil
}

func nilAuthority(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	}
	return false
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Access-Control-Allow-Origin", "") // Never grant cross-origin browser reads.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		unavailable(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil || r.URL.Path != route || r.URL.RawQuery != "" || r.Body == nil {
		unavailable(w, http.StatusBadRequest)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		unavailable(w, http.StatusUnsupportedMediaType)
		return
	}
	var query Request
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&query); err != nil {
		unavailable(w, http.StatusBadRequest)
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		unavailable(w, http.StatusBadRequest)
		return
	}
	if !validRef(query.TenantRef) || query.TenantRef == "*" || !validRef(query.EnvironmentRef) ||
		!validRef(query.EventRef) || !validRef(query.JobID) || !validRef(query.RunID) ||
		!validRef(query.CorrelationID) || !validRef(query.SourceHint) {
		unavailable(w, http.StatusBadRequest)
		return
	}
	if r.Context().Err() != nil {
		unavailable(w, http.StatusServiceUnavailable)
		return
	}
	ctx, session, err := h.sessions.ResolveInspectionSession(r)
	if err != nil || ctx == nil || ctx.Err() != nil || !validRef(session.Principal) ||
		!validRef(session.Tenant) || !validRef(session.Environment) ||
		session.Tenant != query.TenantRef || session.Environment != query.EnvironmentRef {
		unavailable(w, http.StatusNotFound)
		return
	}
	// Authorization is source-owned. The query and sourceHint never grant access.
	allowed, err := h.search.CanInspectSearchResult(ctx, session, query)
	if err != nil || !allowed || ctx.Err() != nil || r.Context().Err() != nil {
		unavailable(w, http.StatusNotFound)
		return
	}
	owned, err := h.events.ReadAuthorizedInspectionEvent(ctx, session, query)
	if err != nil || !owned.Allowed || ctx.Err() != nil || r.Context().Err() != nil ||
		owned.Event.Ref != query.EventRef || owned.Event.TenantRef != session.Tenant {
		unavailable(w, http.StatusNotFound)
		return
	}
	result := inspection.Project(inspection.Input{
		TenantRef: session.Tenant, AuthorizedTenants: []string{session.Tenant},
		EnvironmentRef: session.Environment, Permissions: []string{inspection.PermissionTelemetryRead},
		RawAccessDecision: owned.RawDecision, RenderedAccessDecision: owned.RenderDecision,
		CorrelationID: query.CorrelationID, Event: owned.Event,
	})
	if !result.Allowed || result.Projection == nil {
		unavailable(w, http.StatusNotFound)
		return
	}
	response, ok := toBrowserProjection(*result.Projection, owned.RawDecision, owned.RenderDecision)
	if !ok || ctx.Err() != nil || r.Context().Err() != nil {
		unavailable(w, http.StatusNotFound)
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		unavailable(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(data, '\n'))
}

func validRef(s string) bool { return opaqueRef.MatchString(s) }

type payload struct {
	Access string `json:"access"`
	Text   string `json:"text,omitempty"`
}

type browserProjection struct {
	Kind                string                            `json:"kind"`
	TenantRef           string                            `json:"tenantRef"`
	EnvironmentRef      string                            `json:"environmentRef"`
	EventRef            string                            `json:"eventRef"`
	CorrelationID       string                            `json:"correlationId"`
	State               inspection.EventState             `json:"state"`
	Source              inspection.SourceMetadata        `json:"source"`
	Parser              inspection.ParserMetadata        `json:"parser"`
	SourceFields        map[string]string                 `json:"sourceFields"`
	NormalizedFields    map[string]string                 `json:"normalizedFields"`
	RawDerivedSensitive map[string]string                 `json:"rawDerivedSensitive,omitempty"`
	MissingFields       []string                          `json:"missingFields"`
	Enrichments         []inspection.ProjectionEnrichment `json:"enrichments"`
	Raw                 payload                            `json:"raw"`
	Rendered            payload                            `json:"rendered"`
}

func toBrowserProjection(p inspection.Projection, raw, rendered inspection.AccessDecision) (browserProjection, bool) {
	result := browserProjection{
		Kind: "inspection", TenantRef: p.TenantRef, EnvironmentRef: p.EnvironmentRef,
		EventRef: p.EventRef, CorrelationID: p.CorrelationID, State: p.State,
		Source: p.Source, Parser: p.Parser, SourceFields: p.SourceFields,
		NormalizedFields: p.NormalizedFields, RawDerivedSensitive: p.RawDerivedSensitive,
		MissingFields: p.MissingFields, Enrichments: p.Enrichments,
		Raw: payload{Access: "denied"}, Rendered: payload{Access: "denied"},
	}
	if !validRef(p.Source.SourceRef) || !boundedFields(p.SourceFields) ||
		!boundedFields(p.NormalizedFields) || !boundedFields(p.RawDerivedSensitive) ||
		len(p.MissingFields) > 256 || len(p.Enrichments) > 64 {
		return browserProjection{}, false
	}
	for _, field := range p.MissingFields {
		if !validField(field) { return browserProjection{}, false }
	}
	for _, e := range p.Enrichments {
		if !validRef(e.Producer) || !boundedText(e.Version, 256) ||
			!boundedText(e.Freshness, 256) || !boundedFields(e.Values) {
			return browserProjection{}, false
		}
	}
	if raw == inspection.AccessAllow {
		if !utf8.Valid(p.RawPayload) || len(p.RawPayload) > 65536 { return browserProjection{}, false }
		result.Raw = payload{Access: "allowed", Text: string(p.RawPayload)}
	}
	if rendered == inspection.AccessAllow {
		if !utf8.Valid(p.RenderedPayload) || len(p.RenderedPayload) > 65536 { return browserProjection{}, false }
		result.Rendered = payload{Access: "allowed", Text: string(p.RenderedPayload)}
	}
	return result, true
}

var fieldName = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)

func validField(s string) bool {
	return fieldName.MatchString(s) && s != "__proto__" && s != "constructor" && s != "prototype"
}

func boundedText(s string, max int) bool {
	return len(s) <= max && utf8.ValidString(s) && !strings.ContainsRune(s, '\x00')
}

func boundedFields(fields map[string]string) bool {
	if len(fields) > 256 { return false }
	for k, v := range fields {
		if !validField(k) || !boundedText(v, 4096) { return false }
	}
	return true
}

func unavailable(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = io.WriteString(w, "{\"error\":\"request-unavailable\"}\n")
}
