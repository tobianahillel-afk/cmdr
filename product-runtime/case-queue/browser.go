package casequeue

import (
 "embed"
 "errors"
 "net/http"
)

//go:embed browser.html browser.js browser.css
var browserAssets embed.FS

// NewAuthorizedQueueBrowserHandler is an opt-in, same-origin, host-mounted
// read-only consumer for the existing first-party Queue HTTP projection.
// Authentication/session, owner providers, TLS and CSRF remain with the host;
// static files NEVER embed Case data, access tokens, or permission decisions.
func NewAuthorizedQueueBrowserHandler(sessions QueueSessionResolver,
 cases QueueCaseProvider, access QueueReadAuthorizer,
 views QueueSavedViewProvider) (http.Handler, error) {
 var api http.Handler
 var err error
 if views == nil {
  api, err = NewAuthorizedQueueHTTPHandler(sessions, cases, access)
 } else {
  api, err = NewAuthorizedQueueHTTPHandlerWithSavedViews(sessions, cases, access, views)
 }
 if err != nil { return nil, err }
 if api == nil { return nil, errors.New("authorized Queue handler required") }
 return queueBrowser{api: api}, nil
}

type queueBrowser struct { api http.Handler }

func (h queueBrowser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 if r.URL == nil { http.Error(w, "not found", http.StatusNotFound); return }
 if r.URL.Path == "/case-queue" || r.URL.Path == "/case-queue/return" {
  h.api.ServeHTTP(w, r)
  return
 }
 var name, mimeType string
 switch r.URL.Path {
 case "/case-queue/ui": name, mimeType = "browser.html", "text/html; charset=utf-8"
 case "/case-queue/ui.js": name, mimeType = "browser.js", "text/javascript; charset=utf-8"
 case "/case-queue/ui.css": name, mimeType = "browser.css", "text/css; charset=utf-8"
 default: http.Error(w, "not found", http.StatusNotFound); return
 }
 w.Header().Set("Cache-Control", "no-store")
 w.Header().Set("Referrer-Policy", "no-referrer")
 w.Header().Set("X-Content-Type-Options", "nosniff")
 w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
 if r.Method != http.MethodGet {
  w.Header().Set("Allow", http.MethodGet)
  http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  return
 }
 if r.URL.RawQuery != "" {
  http.Error(w, "invalid UI request", http.StatusBadRequest)
  return
 }
 body, err := browserAssets.ReadFile(name)
 if err != nil { http.Error(w, "unavailable", http.StatusServiceUnavailable); return }
 w.Header().Set("Content-Type", mimeType)
 w.WriteHeader(http.StatusOK)
 _, _ = w.Write(body)
}
