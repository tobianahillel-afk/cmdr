package investigationtimeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

// CaseListCatalog provides candidate Case IDs from a trusted owner provider.
// Membership is NOT authority: every candidate must separately pass CheckAccess.
type CaseListCatalog interface {
	ListCandidateCaseIDs(context.Context, Scope) ([]string, error)
}

// NewCaseDiscoveryHTTPHandler exposes a read-only, authenticated Case picker
// using the same server-derived principal and independent Case permission gate
// as OpenCaseTimeline. Host applications own TLS, routing, audit and throttling.
func NewCaseDiscoveryHTTPHandler(sessions CaseTimelineSessionResolver, checker AccessChecker,
	catalog CaseListCatalog) (http.Handler, error) {
	if missingAuthority(sessions) || missingAuthority(checker) || missingAuthority(catalog) {
		return nil, fmt.Errorf("Case discovery requires session, owner catalog and Case permission services")
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
		if r.URL == nil || r.URL.RawQuery != "" ||
			r.Body != nil && r.Body != http.NoBody || len(r.TransferEncoding) != 0 {
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
		// Candidate enumeration is not a permission decision.
		ids, err := catalog.ListCandidateCaseIDs(ctx, Scope{
			Tenant: scope.Tenant, Environment: scope.Environment,
		})
		if err != nil || ctx.Err() != nil || len(ids) > MaxEntries {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		seen := make(map[string]bool, len(ids))
		used := len(scope.Tenant) + len(scope.Environment)
		for _, id := range ids {
			if !validHTTPCaseSelector(id) ||
				seen[id] || len(id) > MaxInputBytes-used {
				writeCaseTimelineError(w, http.StatusNotFound)
				return
			}
			seen[id] = true
			used += len(id)
		}
		allowed := make([]string, 0, len(ids))
		for _, id := range ids {
			if ctx.Err() != nil || r.Context().Err() != nil {
				writeCaseTimelineError(w, http.StatusServiceUnavailable)
				return
			}
			decision, checkErr := checker.CheckAccess(ctx, AccessRequest{
				Tenant: scope.Tenant, Environment: scope.Environment,
				CaseID: id, Resource: "case", ObjectID: id,
				Permission: caseReadPermission,
			})
			if checkErr != nil || ctx.Err() != nil {
				writeCaseTimelineError(w, http.StatusNotFound)
				return
			}
			if decision == "allow" {
				allowed = append(allowed, id)
			} else if decision != "deny" && decision != "missing" && decision != "unknown" {
				writeCaseTimelineError(w, http.StatusNotFound)
				return
			}
		}
		sort.Strings(allowed)
		if ctx.Err() != nil || r.Context().Err() != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		encoded, err := json.Marshal(struct {
			Cases []string `json:"cases"`
		}{Cases: allowed})
		if err != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(append(encoded, '\n'))
	}), nil
}
