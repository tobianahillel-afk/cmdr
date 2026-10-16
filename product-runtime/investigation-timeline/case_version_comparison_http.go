package investigationtimeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"
)

// CaseVersionComparisonHTTPHandler compares two already-authorized Timeline
// snapshots. It never resolves owner content or authorizes via browser fields.
type CaseVersionComparisonHTTPHandler struct {
	sessions CaseTimelineSessionResolver
	checker AccessChecker
	catalog CaseTimelineCatalog
	entries EntryResolver
	sources SourceOwnerResolver
}

func NewCaseVersionComparisonHTTPHandler(sessions CaseTimelineSessionResolver,
	checker AccessChecker, catalog CaseTimelineCatalog, entries EntryResolver,
	sources SourceOwnerResolver) (http.Handler, error) {
	if missingAuthority(sessions) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) || missingAuthority(sources) {
		return nil, fmt.Errorf("version comparison requires authenticated owner services")
	}
	return &CaseVersionComparisonHTTPHandler{sessions, checker, catalog, entries, sources}, nil
}

type caseVersionComparison struct {
	Status string `json:"status"`
	Current ProjectedEntry `json:"current"`
	Previous ProjectedEntry `json:"previous"`
	// Neither row is a redirect, authorization, or a durable evidence export.
	IsAuthority bool `json:"is_authority"`
}

func (h *CaseVersionComparisonHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Vary", "Authorization")
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
	if err != nil || len(q) != 2 || len(q["case_id"]) != 1 || len(q["entry_id"]) != 1 {
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
	// Reuse the full Case -> Entry -> canonical source ACL chain before
	// selecting either side. A hidden predecessor is redacted by Project.
	result := OpenCaseTimeline(ctx, Scope{Tenant: scope.Tenant, Environment: scope.Environment},
		caseID, h.checker, h.catalog, h.entries, h.sources)
	if result.Status == "deny" {
		writeCaseTimelineError(w, http.StatusNotFound)
		return
	}
	var current *ProjectedEntry
	for i := range result.Entries {
		if result.Entries[i].ID == entryID {
			current = &result.Entries[i]
			break
		}
	}
	if current == nil || current.Supersedes == "" {
		writeCaseTimelineError(w, http.StatusNotFound)
		return
	}
	var previous *ProjectedEntry
	for i := range result.Entries {
		if result.Entries[i].ID == current.Supersedes {
			previous = &result.Entries[i]
			break
		}
	}
	if previous == nil || previous.Version >= current.Version {
		writeCaseTimelineError(w, http.StatusNotFound)
		return
	}
	if ctx.Err() != nil || r.Context().Err() != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	encoded, err := json.Marshal(caseVersionComparison{
		Status: "ready", Current: *current, Previous: *previous, IsAuthority: false,
	})
	if err != nil {
		writeCaseTimelineError(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(encoded, '\n'))
}

var _ http.Handler = (*CaseVersionComparisonHTTPHandler)(nil)
