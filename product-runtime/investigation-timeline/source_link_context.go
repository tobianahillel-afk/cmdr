package investigationtimeline

import (
	"context"
	"net/http"
)

// sourceLinkContextKey cannot be set by an HTTP query/header or another
// package. Only a trusted host which actually mounts the independently
// authorized source-open endpoint opts a Timeline handler into local links.
type sourceLinkContextKey struct{}

// WithLocalSourceLinks opt-ins a trusted host to the existing synthetic
// /demo-source-open destination. It grants no source permission: the
// destination independently authenticates and checks Case, Entry and source.
// Unwrapped Timeline handlers preserve their former HTML and JSON shape.
func WithLocalSourceLinks(timeline http.Handler) http.Handler {
	if timeline == nil {
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			http.Error(w, "request unavailable", http.StatusBadRequest)
			return
		}
		timeline.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), sourceLinkContextKey{}, true)))
	})
}
