package casequeue

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
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

// QueueSavedViewProvider is an optional trusted Shared Capabilities adapter.
// It resolves the current session's selected Saved View on EVERY request,
// including return. Returned AccessDecision values for the view and each
// element MUST come from server-side, permission-rechecked Shared policy.
// The HTTP client cannot submit a Saved View object or choose tenant/scope.
type QueueSavedViewProvider interface {
	CurrentQueueSavedView(context.Context, QueueScope) (*SavedView, error)
}

type nativeQueueHTTP struct {
	sessions QueueSessionResolver
	cases QueueCaseProvider
	access QueueReadAuthorizer
	savedViews QueueSavedViewProvider
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

// NewAuthorizedQueueHTTPHandlerWithSavedViews opts into the existing native
// permission-safe Saved View projector using the host's authoritative Shared
// selection. No Saved View is created, updated, stored or selected by a URL.
// Legacy hosts without a Shared adapter retain the default-view behavior.
func NewAuthorizedQueueHTTPHandlerWithSavedViews(sessions QueueSessionResolver,
	cases QueueCaseProvider, access QueueReadAuthorizer,
	savedViews QueueSavedViewProvider) (http.Handler, error) {
	if savedViews == nil {
		return nil, errors.New("Shared Saved View resolver required")
	}
	base, err := NewAuthorizedQueueHTTPHandler(sessions, cases, access)
	if err != nil { return nil, err }
	handler := base.(nativeQueueHTTP)
	handler.savedViews = savedViews
	return handler, nil
}

func (h nativeQueueHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Vary", "Authorization")
	// Return state contains Case IDs, tenant and analyst filters. Receive it
	// via POST JSON, never a GET URL, browser history or server query logs.
	restoring := r.URL != nil && r.URL.Path == "/case-queue/return"
	if (restoring && r.Method != http.MethodPost) ||
		(!restoring && r.Method != http.MethodGet) {
		if restoring { w.Header().Set("Allow", http.MethodPost) } else { w.Header().Set("Allow", http.MethodGet) }
		queueError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil {
		queueError(w, http.StatusBadRequest)
		return
	}
	var returnToken string
	var parameters url.Values
	if restoring {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" || r.URL.RawQuery != "" ||
			r.Body == nil {
			queueError(w, http.StatusBadRequest)
			return
		}
		const maxReturnBody = maxReturnTokenBytes*2 + 128
		body, err := io.ReadAll(io.LimitReader(r.Body, maxReturnBody+1))
		if err != nil || len(body) == 0 || len(body) > maxReturnBody {
			queueError(w, http.StatusBadRequest)
			return
		}
		// A token must be the ONLY key, exactly once. In particular no
		// client-supplied owner, tenant or permission field is allowed.
		decoder := json.NewDecoder(strings.NewReader(string(body)))
		start, err := decoder.Token()
		if err != nil || start != json.Delim('{') || !decoder.More() {
			queueError(w, http.StatusBadRequest)
			return
		}
		key, err := decoder.Token()
		if err != nil || key != "return_token" ||
			decoder.Decode(&returnToken) != nil || decoder.More() {
			queueError(w, http.StatusBadRequest)
			return
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim('}') || decoder.Decode(new(any)) != io.EOF ||
			returnToken == "" || len(returnToken) > maxReturnTokenBytes*2 {
			queueError(w, http.StatusBadRequest)
			return
		}
	} else {
		if r.Body != nil && r.Body != http.NoBody ||
			len(r.TransferEncoding) != 0 || len(r.URL.RawQuery) > 4096 {
			queueError(w, http.StatusBadRequest)
			return
		}
		var err error
		parameters, err = url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(parameters) > 10 {
			queueError(w, http.StatusBadRequest)
			return
		}
	}
	query := ""
	openCase := ""
	inspectCase := ""
	scrollOffset := 0
	scrollProvided := false
	pageLimit := 0
	pageCursor := ""
	var filters []Filter
	var sortSpecs []SortSpec
	// Only literal, single-valued queue controls reach the existing
	// projector. Authorization of ALL candidates and nested sources still
	// happens first; nothing here selects a tenant or confers permissions.
	for name, values := range parameters {
		literalLimit := 256
		if name == "open_case" || name == "inspect_case" || name == "incident_id" { literalLimit = 1024 }
		if len(values) != 1 || !validQueueLiteral(values[0], literalLimit) {
			queueError(w, http.StatusBadRequest)
			return
		}
		value := values[0]
		switch name {
		case "open_case":
			if strings.TrimSpace(value) == "" {
				queueError(w, http.StatusBadRequest)
				return
			}
			openCase = value // opaque ID; leading/trailing spaces are significant
		case "inspect_case":
			if strings.TrimSpace(value) == "" {
				queueError(w, http.StatusBadRequest)
				return
			}
			inspectCase = value // exact opaque ID; selected only from current authorized Queue rows
		case "scroll":
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 0 || parsed > maxReturnScroll {
				queueError(w, http.StatusBadRequest)
				return
			}
			scrollProvided = true
			scrollOffset = parsed
		case "q":
			query = value
		case "limit":
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 || parsed > 100 {
				queueError(w, http.StatusBadRequest)
				return
			}
			pageLimit = parsed
		case "cursor":
			if value == "" {
				queueError(w, http.StatusBadRequest)
				return
			}
			pageCursor = value
		case "incident_id":
			// Navigation from Command Incident Detail is an exact link
			// selector, not a scope or a new permission. The projector
			// sees only independently authorized Incident snapshots.
			if strings.TrimSpace(value) == "" {
				queueError(w, http.StatusBadRequest)
				return
			}
			filters = append(filters, Filter{Field: name, Values: []string{value}})
		case "status", "owner", "freshness":
			if len(value) > 128 || strings.TrimSpace(value) == "" {
				queueError(w, http.StatusBadRequest)
				return
			}
			filters = append(filters, Filter{Field: name, Values: []string{value}})
		case "sort":
			field, direction, ok := strings.Cut(value, ":")
			if !ok || !allowedSortFields[field] ||
				(direction != "asc" && direction != "desc") {
				queueError(w, http.StatusBadRequest)
				return
			}
			sortSpecs = []SortSpec{{Field: field, Direction: direction}}
		default:
			queueError(w, http.StatusBadRequest)
			return
		}
	}
	if (scrollProvided && openCase == "") ||
		(openCase != "" && inspectCase != "") ||
		(pageCursor != "" && pageLimit == 0) {
		queueError(w, http.StatusBadRequest)
		return
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
	// Shared Capabilities, not the client and not the Case candidate records,
	// owns the current Saved View and its per-element permission decisions.
	// On return the same provider runs again, so a revoked/changed version
	// cannot be silently restored from an old presentation token.
	var currentView *SavedView
	if h.savedViews != nil {
		currentView, err = h.savedViews.CurrentQueueSavedView(ctx, scope)
		if err != nil || ctx.Err() != nil || r.Context().Err() != nil {
			queueError(w, http.StatusNotFound)
			return
		}
		if currentView != nil && currentView.TenantRef != scope.Tenant {
			queueError(w, http.StatusNotFound)
			return
		}
	}
	// Reuse the ONE existing first-party Saved View sanitizer and projector.
	// Forbidden elements are stripped before any Case Queue serialization.
	in := Input{
		TenantRef: scope.Tenant, EnvironmentRef: scope.Environment,
		AuthorizedTenants: []string{scope.Tenant},
		SearchQuery: query, Filters: filters, Sort: sortSpecs,
		SavedView: currentView, Cases: authorized,
	}
	result := ProjectWithSavedViewFallback(in)
	if !result.Allowed || result.Projection == nil || ctx.Err() != nil ||
		r.Context().Err() != nil {
		queueError(w, http.StatusNotFound)
		return
	}
	var response []byte
	if restoring {
		// This routine checks the *current* scoped and independently
		// authorized owner, tenant, environment, view and Case visibility.
		// The decoded token NEVER grants permission to a Case.
		restored, code := RestoreCaseQueue(returnToken, in)
		if code != ErrorNone {
			queueError(w, http.StatusNotFound)
			return
		}
		response, err = json.Marshal(restored)
	} else if inspectCase != "" {
		// Shared Inspector selection consumes the SAME currently authorized
		// Queue row as listing. It is read-only: no new owner, ACL, source
		// traversal, Case mutation or UI implementation.
		visible := result.Projection
		if pageLimit != 0 {
			page, pageErr := paginateQueueProjection(visible, pageLimit, pageCursor, currentView)
			if pageErr != nil {
				queueError(w, http.StatusNotFound)
				return
			}
			visible = page.Projection
		}
		found := false
		for _, row := range visible.Rows {
			if row.Case.ID == inspectCase {
				response, err = json.Marshal(row)
				found = true
				break
			}
		}
		if !found {
			queueError(w, http.StatusNotFound)
			return
		}
	} else if openCase != "" {
		// Reuse the existing CAP-INV-101 Class-0 handoff engine. It checks
		// selection visibility AFTER current server ACL, filters and freshness.
		// The token is presentation state, not a credential or a signed grant.
		handoff := openCaseWorkspaceWithPage(in, openCase, scrollOffset,
			OriginCaseQueue, pageLimit, pageCursor)
		if !handoff.Allowed || handoff.Handoff == nil {
			queueError(w, http.StatusNotFound)
			return
		}
		response, err = json.Marshal(handoff.Handoff)
	} else if pageLimit != 0 {
		page, pageErr := paginateQueueProjection(result.Projection, pageLimit, pageCursor, currentView)
		if pageErr != nil {
			queueError(w, http.StatusNotFound)
			return
		}
		response, err = json.Marshal(page)
	} else {
		response, err = json.Marshal(result.Projection)
	}
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

func validQueueLiteral(value string, maxBytes int) bool {
	if len(value) > maxBytes || !utf8.ValidString(value) {
		return false
	}
	for _, c := range value {
		if c < 0x20 || c == 0x7f {
			return false
		}
	}
	return true
}

func queueError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte("{\"error\":\"request-unavailable\"}\n"))
}

var _ http.Handler = nativeQueueHTTP{}
