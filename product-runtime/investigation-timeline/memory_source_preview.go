package investigationtimeline

import (
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// sourcePreviewHandler opens a synthetic owner-owned source from a selected
// Timeline Entry. It deliberately is NOT a production source API: the owner
// snapshot and its preview are static, in memory, and have no persistence.
func (s *memoryTimelineProvider) sourcePreviewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Vary", "Authorization")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
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
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(query) != 2 || len(query["case_id"]) != 1 || len(query["entry_id"]) != 1 {
			writeCaseTimelineError(w, http.StatusBadRequest)
			return
		}
		caseID, entryID := query.Get("case_id"), query.Get("entry_id")
		if invalidScope(caseID) || invalidScope(entryID) ||
			len(caseID) > maxHTTPCaseIDBytes || len(entryID) > maxHTTPCaseIDBytes ||
			!utf8.ValidString(caseID) || !utf8.ValidString(entryID) ||
			strings.ContainsAny(caseID+entryID, "\x00\r\n") {
			writeCaseTimelineError(w, http.StatusBadRequest)
			return
		}
		if r.Context().Err() != nil {
			writeCaseTimelineError(w, http.StatusServiceUnavailable)
			return
		}
		ctx, scope, err := s.ResolveCaseTimelineSession(r)
		if err != nil || ctx == nil || ctx.Err() != nil {
			writeCaseTimelineError(w, http.StatusUnauthorized)
			return
		}
		// The *fresh* server-derived Case -> Timeline -> owner authorization
		// is evaluated again; no client-provided navigation draft is trusted.
		draft, ok := DraftCaseSourceNavigation(ctx, scope, caseID, entryID, s, s, s, s)
		if !ok {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		// The opened owner object is identified from trusted server records.
		// A draft is not a capability and a sourceRef is never taken from HTTP.
		key := memorySourceObjectKey{caseID: caseID, kind: draft.Kind, objectID: draft.ObjectID}
		owner, present := s.sourceObjects[key]
		if !present || owner.Identity.Tenant != scope.Tenant ||
			owner.Identity.Environment != scope.Environment || owner.Identity.CaseID != caseID ||
			owner.Identity.Kind != draft.Kind || owner.Identity.ObjectID != draft.ObjectID ||
			owner.Identity.Reference != s.entries[entryID].Value.SourceRef {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		permission, valid := canonicalSourceReadPermission(draft.Kind)
		if !valid {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		decision, err := s.CheckAccess(ctx, AccessRequest{
			Tenant: scope.Tenant, Environment: scope.Environment, CaseID: caseID,
			Resource: draft.Kind, ObjectID: draft.ObjectID, Permission: permission,
		})
		if err != nil || decision != "allow" || ctx.Err() != nil || r.Context().Err() != nil {
			writeCaseTimelineError(w, http.StatusNotFound)
			return
		}
		payload := sourcePreviewResponse{
			Status: "ready", CaseID: caseID, EntryID: entryID,
			Kind: draft.Kind, ObjectID: draft.ObjectID,
			Preview: owner.Preview, Synthetic: true,
		}
		// Presentation is separate from the source-owner trust decision.
		writeAuthorizedSourcePreview(w, r, ctx, payload)
	})
}
