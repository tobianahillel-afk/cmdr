package investigationtimeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"
)

// CaseSourceNavigationHTTPHandler returns an authorized, provider-neutral
// destination draft, never an arbitrary URL or redirect.
type CaseSourceNavigationHTTPHandler struct {
	sessions CaseTimelineSessionResolver
	checker AccessChecker
	catalog CaseTimelineCatalog
	entries EntryResolver
	sources SourceOwnerResolver
}

func NewCaseSourceNavigationHTTPHandler(sessions CaseTimelineSessionResolver,
	checker AccessChecker, catalog CaseTimelineCatalog, entries EntryResolver,
	sources SourceOwnerResolver) (http.Handler, error) {
	if missingAuthority(sessions) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) || missingAuthority(sources) {
		return nil, fmt.Errorf("source navigation requires authenticated owner services")
	}
	return &CaseSourceNavigationHTTPHandler{sessions, checker, catalog, entries, sources}, nil
}

type navigationHTTPDraft struct {
	SourceNavigationDraft
	ReturnQuery string `json:"return_query"`
	ReturnIsAuthority bool `json:"return_is_authority"`
}

func (h *CaseSourceNavigationHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeCaseTimelineError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil || r.Body != nil && r.Body != http.NoBody ||
		len(r.TransferEncoding) != 0 || len(r.URL.RawQuery) > 4096 {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(q["case_id"]) != 1 || len(q["entry_id"]) != 1 || len(q) > 8 {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	caseID, entryID := q.Get("case_id"), q.Get("entry_id")
	if invalidScope(caseID) || invalidScope(entryID) ||
		len(caseID) > maxHTTPCaseIDBytes || len(entryID) > maxHTTPCaseIDBytes ||
		!utf8.ValidString(caseID) || !utf8.ValidString(entryID) {
		writeCaseTimelineError(w, http.StatusBadRequest)
		return
	}
	// Return context is a sanitized, presentation-only selection; reject any
	// caller-provided route, tenant, owner object, URL or permission flag.
	selection := make(url.Values, len(q)-1)
	for key, values := range q {
		if key != "entry_id" {
			selection[key] = values
		}
	}
	if _, valid := parseTimelineReadSelection(selection); !valid {
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
	draft, ok := DraftCaseSourceNavigation(ctx,
		Scope{Tenant: scope.Tenant, Environment: scope.Environment},
		caseID, entryID, h.checker, h.catalog, h.entries, h.sources)
	if !ok {
		writeCaseTimelineError(w, http.StatusNotFound)
		return
	}
	if ctx.Err() != nil || r.Context().Err() != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	encoded, err := json.Marshal(navigationHTTPDraft{
		SourceNavigationDraft: draft,
		ReturnQuery: selection.Encode(),
		ReturnIsAuthority: false,
	})
	if err != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(encoded, '\n'))
}

var _ http.Handler = (*CaseSourceNavigationHTTPHandler)(nil)
