package main

import (
	"net/http"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Retained only for fixture compatibility in demonstration tests.
// The runnable viewer uses the shared owner snapshot created in newDemoHandler;
// all discovery now delegates to the same server-side Case permission boundary.
func newDemoCaseListHandler(authenticate func(*http.Request) (string, error), data timeline.MemoryTimelineData) http.Handler {
	handlers, err := timeline.NewMemoryCaseTimelineHTTPHandlers(authenticate, data)
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("{\"error\":\"request-unavailable\"}\n"))
		})
	}
	return handlers.Cases
}
