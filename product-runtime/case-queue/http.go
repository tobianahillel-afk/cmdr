package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// QueueScope is derived from an authenticated server session, never from an
// HTTP parameter or a provider's AccessDecision field.
type QueueScope struct {
	Tenant string
	Environment string
}

// QueueSessionResolver supplies verified session scope. The host owns
// authentication, TLS, throttling and durable permission/audit records.
type QueueSessionResolver interface {
	ResolveQueueSession(*http.Request) (context.Context, QueueScope, error)
}

// QueueCaseProvider supplies candidates, not authorization decisions. Its
// records may include Incident and Finding context, which are reauthorized.
type QueueCaseProvider interface {
	ListQueueCandidates(context.Context, QueueScope) ([]CaseSnapshot, error)
}

// QueueReadAuthorizer is an injected first-party object permission checker.
// Valid resources are "case", "incident" and "finding". Missing is denial.
type QueueReadAuthorizer interface {
	CheckQueueRead(context.Context, QueueScope, string, string) (string, error)
}

type nativeQueueHTTP struct {
	sessions QueueSessionResolver
	cases QueueCaseProvider
	access QueueReadAuthorizer
}

// NewAuthorizedQueueHTTPHandler exposes the existing CAP-INV-101 read-only
// projection through real per-request owner authorization. It deliberately
// does not implement a provider, authentication, Saved View persistence or
// Case mutation while those product integrations are not selected.
func NewAuthorizedQueueHTTPHandler(sessions QueueSessionResolver,
	cases QueueCaseProvider, access QueueReadAuthorizer) (http.Handler, error) {
	if sessions == nil || cases == nil || access == nil {
		return nil, errors.New("Case Queue requires session, candidate provider and read authorization")
	}
	return nativeQueueHTTP{sessions: sessions, cases: cases, access: access}, nil
}

func (h nativeQueueHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Vary", "Authorization")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		queueError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil || r.Body != nil && r.Body != http.NoBody ||
		len(r.TransferEncoding) != 0 || len(r.URL.RawQuery) > 1024 {
		queueError(w, http.StatusBadRequest)
		return
	}
	parameters, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(parameters) > 1 {
		queueError(w, http.StatusBadRequest)
		return
	}
	query := ""
	if len(parameters) != 0 {
		values, ok := parameters["q"]
		if !ok || len(values) != 1 {
			queueError(w, http.StatusBadRequest)
			return
		}
		query = values[0]
	}
	if len(query) > 256 || !utf8.ValidString(query) {
		queueError(w, http.StatusBadRequest)
		return
	}
	for _, c := range query {
		if c < 0x20 || c == 0x7f {
			queueError(w, http.StatusBadRequest)
			return
		}
	}
	if r.Context().Err() != nil {
		queueError(w, http.StatusServiceUnavailable)
		return
	}
	ctx, scope, err := h.sessions.ResolveQueueSession(r)
	if err != nil || ctx == nil || ctx.Err() != nil ||
		strings.TrimSpace(scope.Tenant) == "" || scope.Tenant == "*" ||
		strings.TrimSpace(scope.Environment) == "" {
		queueError(w, http.StatusUnauthorized)
		return
	}
	candidates, err := h.cases.ListQueueCandidates(ctx, scope)
	if err != nil || ctx.Err() != nil || len(candidates) > 5000 {
		queueError(w, http.StatusNotFound)
		return
	}
	// Build only an authorized owner graph. The candidate's embedded
	// AccessDecision must NOT grant rights, including nested projections.
	authorized := make([]CaseSnapshot, 0, len(candidates))
	for _, candidate := range candidates {
		if ctx.Err() != nil || r.Context().Err() != nil {
			queueError(w, http.StatusServiceUnavailable)
			return
		}
		if candidate.TenantRef != scope.Tenant || strings.TrimSpace(candidate.ID) == "" {
			queueError(w, http.StatusNotFound)
			return
		}
		visible, err := h.allowed(ctx, scope, "case", candidate.ID)
		if err != nil {
			queueError(w, http.StatusNotFound)
			return
		}
		if !visible { continue }
		if candidate.Incident != nil {
			if candidate.Incident.TenantRef != scope.Tenant ||
				strings.TrimSpace(candidate.Incident.ID) == "" {
				queueError(w, http.StatusNotFound)
				return
			}
			allowed, err := h.allowed(ctx, scope, "incident", candidate.Incident.ID)
			if err != nil {
				queueError(w, http.StatusNotFound)
				return
			}
			if !allowed {
				candidate.Incident = nil
			} else {
				copyIncident := *candidate.Incident
				copyIncident.AccessDecision = "allow"
				candidate.Incident = &copyIncident
			}
		}
		visibleFindings := make([]FindingContext, 0, len(candidate.Findings))
		for _, finding := range candidate.Findings {
			if finding.TenantRef != scope.Tenant || finding.CaseID != candidate.ID ||
				strings.TrimSpace(finding.ID) == "" {
				queueError(w, http.StatusNotFound)
				return
			}
			allowed, err := h.allowed(ctx, scope, "finding", finding.ID)
			if err != nil {
				queueError(w, http.StatusNotFound)
				return
			}
			if allowed {
				finding.AccessDecision = "allow"
				visibleFindings = append(visibleFindings, finding)
			}
		}
		candidate.Findings = visibleFindings
		candidate.AccessDecision = "allow"
		authorized = append(authorized, candidate)
	}
	// Reuse the one existing native Case Queue projection; no second
	// search/sort engine and no client supplied tenant or Saved View ACL.
	result := ProjectWithSavedViewFallback(Input{
		TenantRef: scope.Tenant, EnvironmentRef: scope.Environment,
		AuthorizedTenants: []string{scope.Tenant},
		SearchQuery: query, Cases: authorized,
	})
	if !result.Allowed || result.Projection == nil || ctx.Err() != nil ||
		r.Context().Err() != nil {
		queueError(w, http.StatusNotFound)
		return
	}
	response, err := json.Marshal(result.Projection)
	if err != nil || len(response) > 4<<20 {
		queueError(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(response, '\n'))
}

func (h nativeQueueHTTP) allowed(ctx context.Context, scope QueueScope,
	resource, id string) (bool, error) {
	decision, err := h.access.CheckQueueRead(ctx, scope, resource, id)
	if err != nil || ctx.Err() != nil { return false, errors.New("owner authorization failed") }
	switch decision {
	case "allow": return true, nil
	case "deny", "missing", "unknown": return false, nil
	default: return false, errors.New("invalid owner authorization decision")
	}
}

func queueError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte("{\"error\":\"request-unavailable\"}\n"))
}

var _ http.Handler = nativeQueueHTTP{}
