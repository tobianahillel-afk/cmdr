package commandinvestigatehost

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Route composition is not an authorization upgrade. The optional detail
// handler must be explicitly supplied and remains a separate host owner.
func TestIncidentDetailOptionalHostDispatchAndLegacyRoutes(t *testing.T) {
	commandCalls, timelineCalls, detailCalls := 0, 0, 0
	command := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commandCalls++
		w.WriteHeader(http.StatusAccepted)
	})
	timeline := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timelineCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	detail := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		detailCalls++
		if r.Method != http.MethodGet {
			http.Error(w, "method denied by detail owner", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	for _, tc := range []struct {
		name string
		handler http.Handler
		want int
	}{
		{"no-detail", nil, http.StatusNotFound},
		{"with-detail", detail, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var h http.Handler
			var err error
			if tc.handler == nil {
				h, err = NewHandler(command, timeline)
			} else {
				h, err = NewHandler(command, timeline, tc.handler)
			}
			if err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
				"/command/incident-detail?incident=event-a", nil))
			if rec.Code != tc.want || rec.Header().Get("Cache-Control") != "no-store" ||
				rec.Header().Get("Referrer-Policy") != "no-referrer" ||
				rec.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatalf("optional detail routing leaked or bypassed host: %d %+v", rec.Code, rec.Header())
			}
			oldCommand := httptest.NewRecorder()
			h.ServeHTTP(oldCommand, httptest.NewRequest(http.MethodGet, "/command/work-queue", nil))
			oldCase := httptest.NewRecorder()
			h.ServeHTTP(oldCase, httptest.NewRequest(http.MethodGet, "/case-timeline", nil))
			if oldCommand.Code != http.StatusAccepted || oldCase.Code != http.StatusNoContent {
				t.Fatal("existing owner routes were redirected or replaced")
			}
			unmatched := httptest.NewRecorder()
			h.ServeHTTP(unmatched, httptest.NewRequest(http.MethodGet, "/command/incident-detail/extra", nil))
			if unmatched.Code != http.StatusNotFound {
				t.Fatalf("near-miss route accepted: %d", unmatched.Code)
			}
		})
	}
	if detailCalls != 1 || commandCalls != 2 || timelineCalls != 2 {
		t.Fatalf("wrong owner received a request: command=%d case=%d detail=%d", commandCalls, timelineCalls, detailCalls)
	}
}

func TestIncidentDetailHandlerConstructorRejectsAmbiguousDelegations(t *testing.T) {
	ok := http.NotFoundHandler()
	for _, args := range []struct {
		command, timeline http.Handler
		detail []http.Handler
	}{
		{nil, ok, []http.Handler{ok}},
		{ok, nil, []http.Handler{ok}},
		{ok, ok, []http.Handler{nil}},
		{ok, ok, []http.Handler{ok, ok}},
	} {
		if h, err := NewHandler(args.command, args.timeline, args.detail...); err == nil || h != nil {
			t.Fatal("host accepted missing or ambiguous owner handler")
		}
	}
	h, err := NewHandler(ok, ok, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "denied by detail", http.StatusForbidden)
	}))
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/command/incident-detail?incident=private", nil))
	if rec.Code != http.StatusForbidden {
		t.Fatal("host bypassed detail owner rejection")
	}
}
