// Package commandinvestigatehost composes the existing, independently
// authorized Command Work Queue and Investigate Case Timeline HTTP handlers.
// A navigation link does not grant the right to read the destination Case.
package commandinvestigatehost

import (
	"errors"
	"net/http"
)

// NewHandler makes the existing Class-0 Command -> Investigate navigation
// operable under one origin. The caller MUST construct both handlers with
// independently trusted sessions and owner adapters. The router neither
// authenticates nor promotes an Incident permission into a Case permission.
func NewHandler(command, timeline http.Handler, detail ...http.Handler) (http.Handler, error) {
	if command == nil || timeline == nil || len(detail) > 1 ||
		(len(detail) == 1 && detail[0] == nil) {
		return nil, errors.New("independently authorized read handlers are required")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if r == nil || r.URL == nil {
			http.Error(w, "request unavailable", http.StatusBadRequest)
			return
		}
		switch r.URL.Path {
		case "/command/work-queue":
			command.ServeHTTP(w, r)
		case "/command/incident-detail":
			if len(detail) != 1 {
				http.NotFound(w, r)
				return
			}
			detail[0].ServeHTTP(w, r)
		case "/case-timeline":
			timeline.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}), nil
}
