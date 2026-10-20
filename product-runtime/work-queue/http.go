package workqueue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// SessionScope is trusted identity resolved by the host. HTTP query strings,
// headers, cookies and client-selected Case/Incident IDs MUST NOT be accepted
// as principal, tenant or environment by the Work Queue itself.
type SessionScope struct {
	Principal   string
	Tenant      string
	Environment string
}

// SessionResolver is a host-owned authenticated session boundary, not an
// authorization engine. The host supplies a fresh scope for every request.
// The Owner independently authorizes each Incident/Task in Project.
type SessionResolver interface {
	ResolveSession(ctx context.Context, request *http.Request) (SessionScope, error)
}

// NewHTTPHandler is a single read-only Command-owned route for six system
// Work Queue views. There are deliberately no view-specific routes, no saved
// view persistence and no client-selected authorization scope.
func NewHTTPHandler(sessions SessionResolver, owner Owner) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		if r.URL.Path != "/command/work-queue" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if sessions == nil || owner == nil {
			http.Error(w, "work queue unavailable", http.StatusServiceUnavailable)
			return
		}
		// Do not accept client-supplied principal, tenant, environment or
		// arbitrary filter DSL at this first bounded system-view boundary.
		values, queryErr := url.ParseQuery(r.URL.RawQuery)
		if queryErr != nil {
			http.Error(w, "invalid work queue query", http.StatusBadRequest)
			return
		}
		for key, entries := range values {
			if (key != "view" && key != "status" && key != "sort" && key != "selected" && key != "cursor") || len(entries) != 1 {
				http.Error(w, "unsupported work queue query", http.StatusBadRequest)
				return
			}
		}
		view := All
		if entries, ok := values["view"]; ok {
			if entries[0] == "" {
				http.Error(w, "invalid work queue view", http.StatusBadRequest)
				return
			}
			view = View(entries[0])
		}
		// Status is only a literal read-only projection filter; it never
		// selects principal, tenant, environment, or a Saved View owner.
		status := ""
		if entries, ok := values["status"]; ok {
			status = entries[0]
			if len(status) > 64 || (status != "" && invalidScope(status)) {
				http.Error(w, "invalid work queue status", http.StatusBadRequest)
				return
			}
		}
		// Only the two explicitly supported stable timestamp orders. Reject
		// untrusted sort expressions before resolving an owner or its scope.
		sortOrder := ""
		if entries, ok := values["sort"]; ok {
			sortOrder = entries[0]
			if sortOrder != "newest" && sortOrder != "oldest" {
				http.Error(w, "invalid work queue sort", http.StatusBadRequest)
				return
			}
		}
		// The selection token is an object ID hint, NEVER authority. Validate
		// it before session/owner access and resolve it only against the
		// current, fully authorized, post-filter bounded visible result.
		selection := ""
		if entries, ok := values["selected"]; ok {
			selection = entries[0]
			kind, id, parsed := strings.Cut(selection, ":")
			if !parsed || (kind != string(Incident) && kind != string(Task)) ||
				len(id) > 256 || !utf8.ValidString(id) || invalidScope(id) {
				http.Error(w, "invalid work queue selection", http.StatusBadRequest)
				return
			}
		}
		// Reject malformed cursors before host sessions or owner data are read.
		// Correctly formatted but stale/foreign cursors get one generic 404
		// only AFTER the current full owner authorization is recomputed.
		pageCursor := ""
		if entries, ok := values["cursor"]; ok {
			pageCursor = entries[0]
			if pageCursor == "" {
				http.Error(w, "invalid work queue cursor", http.StatusBadRequest)
				return
			}
			if _, err := decodeQueueCursor(pageCursor); err != nil {
				http.Error(w, "invalid work queue cursor", http.StatusBadRequest)
				return
			}
		}
		scope, err := sessions.ResolveSession(r.Context(), r)
		if err != nil || invalidScope(scope.Principal) ||
			invalidScope(scope.Tenant) || invalidScope(scope.Environment) {
			http.Error(w, "work queue session unavailable", http.StatusUnauthorized)
			return
		}
		result, err := Project(r.Context(), owner, Request{
			Principal: scope.Principal, Tenant: scope.Tenant,
			Environment: scope.Environment, View: view, StatusFilter: status,
			SortOrder: sortOrder, Cursor: pageCursor,
		})
		if err != nil {
			if err == ErrCursor {
				http.NotFound(w, r)
			} else if err == ErrView || err == ErrScope {
				http.Error(w, "invalid work queue selection", http.StatusBadRequest)
			} else {
				// Provider failures, missing metadata or revoked ACLs never
				// disclose any partial permitted/denied rows or IDs.
				http.Error(w, "work queue unavailable", http.StatusServiceUnavailable)
			}
			return
		}
		if selection != "" {
			for i := range result.Rows {
				row := &result.Rows[i]
				if string(row.Kind)+":"+row.ID == selection {
					result.Selected = row
					break
				}
			}
			if result.Selected == nil {
				// Identical response whether an ID is nonexistent, denied,
				// outside the current filter or beyond the visible page.
				http.NotFound(w, r)
				return
			}
		}
		if result.Selected != nil {
			// Only AFTER the session, the full owner projection and the
			// current-page selection have been independently authorized.
			// The optional host adapter must authorize each source factor;
			// no client-selected factor or scope is trusted.
			factors := AuthorizedFactors{}
			if source, ok := owner.(FactorsOwner); ok {
				fetched, err := source.ReadAuthorizedFactors(r.Context(), scope.Principal,
					scope.Tenant, scope.Environment, result.Selected.Kind, result.Selected.ID)
				if err != nil {
					http.Error(w, "coordination factors unavailable", http.StatusServiceUnavailable)
					return
				}
				factors, err = validateAuthorizedFactors(fetched)
				if err != nil {
					http.Error(w, "coordination factors unavailable", http.StatusServiceUnavailable)
					return
				}
			}
			result.SelectedFactors = &factors
		}
		w.Header().Set("Vary", "Accept")
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			page, err := renderQueuePage(result)
			if err != nil {
				http.Error(w, "work queue unavailable", http.StatusServiceUnavailable)
				return
			}
			// Navigation is a fresh request to this SAME canonical route:
			// the trusted host session and each object ACL are re-evaluated.
			// No client-side script, externally loaded resource or unsafe HTML.
			w.Header().Set("Content-Security-Policy", "default-src 'none'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(page)
			return
		}
		encoded, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "work queue unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(append(encoded, '\n'))
	})
}
