package investigationtimeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

const maxHTTPCaseIDBytes = 1024

// The Case picker and Case Timeline must apply identical selector rules.
// Candidate Case IDs come from an owner provider, not from an HTTP request,
// but must remain safe, reopenable HTTP selectors after authorization.
// ValidHTTPCaseSelector is the single first-party predicate shared by
// Case picker, Timeline HTTP and the local Case viewer bootstrap.
// It preserves opaque, whitespace-distinct owner IDs without normalization.
func ValidHTTPCaseSelector(id string) bool {
	return !invalidScope(id) && len(id) <= maxHTTPCaseIDBytes &&
		utf8.ValidString(id) && !strings.ContainsAny(id, "\x00\r\n")
}

func validHTTPCaseSelector(id string) bool { return ValidHTTPCaseSelector(id) }

// CaseTimelineSessionResolver authenticates a request at the server boundary.
// It must return a context derived from the request carrying the verified
// principal, plus the server-derived tenant/environment; not browser headers.
type CaseTimelineSessionResolver interface {
	ResolveCaseTimelineSession(*http.Request) (context.Context, Scope, error)
}

// CaseTimelineHTTPHandler adapts the already integrated read-only Case flow.
// The host application owns TLS, canonical routing and durable audit logging.
type CaseTimelineHTTPHandler struct {
	sessions CaseTimelineSessionResolver
	checker  AccessChecker
	catalog  CaseTimelineCatalog
	entries  EntryResolver
	sources  SourceOwnerResolver
}

func NewCaseTimelineHTTPHandler(sessions CaseTimelineSessionResolver, checker AccessChecker,
	catalog CaseTimelineCatalog, entries EntryResolver, sources SourceOwnerResolver) (http.Handler, error) {
	if missingAuthority(sessions) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) || missingAuthority(sources) {
		return nil, fmt.Errorf("Case Timeline HTTP requires authenticated owner services")
	}
	return &CaseTimelineHTTPHandler{sessions: sessions, checker: checker,
		catalog: catalog, entries: entries, sources: sources}, nil
}

// ServeHTTP accepts only a Case selector. Neither scope nor permissions nor
// Timeline Entry IDs may be supplied as request parameters. Denials and
// owner-provider failures return an indistinguishable generic 404.
func (h *CaseTimelineHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeCaseTimelineError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil || r.Body != nil && r.Body != http.NoBody ||
		len(r.TransferEncoding) != 0 || len(r.URL.RawQuery) > maxHTTPCaseIDBytes*3+640 {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(query["case_id"]) != 1 || len(query) > 7 {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	caseID := query.Get("case_id")
	selection, valid := parseTimelineReadSelection(query)
	if !valid {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	if !validHTTPCaseSelector(caseID) {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	if r.Context().Err() != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	ctx, scope, err := h.sessions.ResolveCaseTimelineSession(r)
	if err != nil || missingAuthority(ctx) || ctx.Err() != nil ||
		invalidScope(scope.Tenant) || invalidScope(scope.Environment) {
		writeCaseTimelineError(w, http.StatusUnauthorized)
		return
	}
	// Discard all permission flags: OpenCaseTimeline consults the trusted
	// AccessChecker for Case, Timeline Entry and each canonical source.
	result := OpenCaseTimeline(ctx, Scope{
		Tenant: scope.Tenant, Environment: scope.Environment,
	}, caseID, h.checker, h.catalog, h.entries, h.sources)
	if result.Status == "deny" {
		writeCaseTimelineError(w, http.StatusNotFound)
		return
	}
	// Apply view selection after the complete Case/Entry/source ACL chain.
	result = selectAuthorizedTimeline(result, selection)
	if r.Context().Err() != nil || ctx.Err() != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	// Only authorized, already filtered entries may be counted or paged.
	// The non-paged response shape stays backward compatible.
	var response any = result
	if selection.paged {
		response = paginateAuthorizedTimeline(result, selection)
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(encoded, '\n'))
}

func writeCaseTimelineError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"request-unavailable"}` + "\n"))
}

var _ http.Handler = (*CaseTimelineHTTPHandler)(nil)
