package investigationtimeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// NewCaseSourceNavigationHTTPHandler returns a read-only destination draft.
// The host owns routing and the destination owner's independent access gate.
func NewCaseSourceNavigationHTTPHandler(
	sessions CaseTimelineSessionResolver, checker AccessChecker,
	catalog CaseTimelineCatalog, entries EntryResolver, sources SourceOwnerResolver,
) (http.Handler, error) {
	if missingAuthority(sessions) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) ||
		missingAuthority(sources) {
		return nil, fmt.Errorf("Case source navigation requires authenticated owner services")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Vary", "Authorization")
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeCaseTimelineError(w, http.StatusMethodNotAllowed)
			return
		}
		if r.URL == nil || r.Body != nil && r.Body != http.NoBody ||
			len(r.TransferEncoding) != 0 || len(r.URL.RawQuery) > 2*maxHTTPCaseIDBytes*3+64 {
			writeCaseTimelineError(w, http.StatusBadRequest)
			return
		}
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(query) != 2 || len(query["case_id"]) != 1 ||
			len(query["entry_id"]) != 1 ||
			!validNavigationSelector(query.Get("case_id")) ||
			!validNavigationSelector(query.Get("entry_id")) {
			writeCaseTimelineError(w, http.StatusBadRequest)
			return
		}
		if r.Context().Err() != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		ctx, scope, err := sessions.ResolveCaseTimelineSession(r)
		if err != nil || missingAuthority(ctx) || ctx.Err() != nil ||
			invalidScope(scope.Tenant) || invalidScope(scope.Environment) {
			writeCaseTimelineError(w, http.StatusUnauthorized)
			return
		}
		// No claimed permissions, tenant, source kind or owner ID are accepted.
		result := OpenCaseTimelineSourceNavigation(ctx, Scope{
			Tenant: scope.Tenant, Environment: scope.Environment,
		}, query.Get("case_id"), query.Get("entry_id"),
			checker, catalog, entries, sources)
		if result.Status != "ready" || result.Draft == nil {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		if r.Context().Err() != nil || ctx.Err() != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		payload, err := json.Marshal(result)
		if err != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(append(payload, '\n'))
	}), nil
}
