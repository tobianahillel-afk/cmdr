package investigationtimeline

import (
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// validatedSourceReturnURL accepts only a canonical, same-Case Timeline
// read selection. It is a navigation hint, never a delegated capability:
// clicking it always reruns the destination's authorization chain.
func validatedSourceReturnURL(raw, caseID string) (string, bool) {
	if raw == "" || len(raw) > 2048 || strings.Contains(raw, "#") {
		return "", false
	}
	destination, err := url.Parse(raw)
	if err != nil || destination.IsAbs() || destination.Host != "" ||
		destination.User != nil || destination.Opaque != "" ||
		destination.Path != "/case-timeline" ||
		destination.EscapedPath() != "/case-timeline" ||
		destination.Fragment != "" || destination.ForceQuery {
		return "", false
	}
	values, err := url.ParseQuery(destination.RawQuery)
	if err != nil || len(values) == 0 || len(values) > 8 ||
		len(values["case_id"]) != 1 || values.Get("case_id") != caseID {
		return "", false
	}
	origin, valid := extractIncidentOrigin(values)
	if !valid || len(values) > 7 {
		return "", false
	}
	selection, ok := parseTimelineReadSelection(values)
	if !ok {
		return "", false
	}
	if !selection.paged {
		// The producing Timeline view always supplies a bounded page.
		return "", false
	}
	if origin != "" {
		values.Set("origin_incident", origin)
	}
	return timelineHTMLPageURL(values, selection.limit, selection.offset), true
}

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
		if err != nil || (len(query) != 2 && len(query) != 3) ||
			len(query["case_id"]) != 1 || len(query["entry_id"]) != 1 ||
			(len(query) == 3 && len(query["return_to"]) != 1) {
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
		// Navigation data only. No principal, scope, source reference or
		// permission is ever taken from the optional return target.
		returnURL := "/case-timeline?" + (url.Values{"case_id": {caseID}}).Encode()
		if values := query["return_to"]; len(values) != 0 {
			var ok bool
			returnURL, ok = validatedSourceReturnURL(values[0], caseID)
			if !ok {
				writeCaseTimelineError(w, http.StatusBadRequest)
				return
			}
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
			Preview: owner.Preview, Synthetic: true, ReturnURL: returnURL,
		}
		// Presentation is separate from the source-owner trust decision.
		writeAuthorizedSourcePreview(w, r, ctx, payload)
	})
}
