package main

import (
	"encoding/json"
	"net/http"
	"sort"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// newDemoCaseListHandler is ONLY a selector for the immutable local fixture.
// It uses the same verified bearer callback as /case-timeline; no browser
// principal, tenant, permission or list of Case IDs is accepted as authority.
// A production Case catalogue belongs to the host Case owner, not this demo.
func newDemoCaseListHandler(authenticate func(*http.Request) (string, error), data timeline.MemoryTimelineData) http.Handler {
	// Copy the minimal fixture ACL once, independent of request data.
	principals := make(map[string]timeline.MemoryTimelinePrincipal, len(data.Principals))
	for _, p := range data.Principals {
		principals[p.ID] = p
	}
	cases := append([]timeline.MemoryTimelineCase(nil), data.Cases...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Vary", "Authorization")
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, `{"error":"request-unavailable"}`, http.StatusMethodNotAllowed)
			return
		}
		if r.URL == nil || r.URL.RawQuery != "" || r.Body != nil && r.Body != http.NoBody ||
			len(r.TransferEncoding) != 0 {
			http.Error(w, `{"error":"request-unavailable"}`, http.StatusBadRequest)
			return
		}
		if r.Context().Err() != nil {
			http.Error(w, `{"error":"request-unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		principalID, err := authenticate(r)
		p, exists := principals[principalID]
		if err != nil || !exists {
			http.Error(w, `{"error":"request-unavailable"}`, http.StatusUnauthorized)
			return
		}
		allowed := make([]string, 0, len(cases))
		for _, c := range cases {
			if c.Tenant != p.Tenant || c.Environment != p.Environment {
				continue
			}
			for _, reader := range c.Readers {
				if reader == principalID {
					allowed = append(allowed, c.ID)
					break
				}
			}
		}
		sort.Strings(allowed)
		if r.Context().Err() != nil {
			http.Error(w, `{"error":"request-unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(struct {
			Cases []string `json:"cases"`
		}{Cases: allowed})
	})
}
