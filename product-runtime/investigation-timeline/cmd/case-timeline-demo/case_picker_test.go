package main

import (
	"context"
	"errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestDemoCaseCatalogueToAuthorizedTimelineJourney(t *testing.T) {
	for _, tc := range []struct{ principal, expected, foreign string }{
		{"analyst-a", "case-a", "case-b"},
		{"analyst-b", "case-b", "case-a"},
	} {
		t.Run(tc.principal, func(t *testing.T) {
			handler, err := newDemoHandler(testBearer, tc.principal)
			if err != nil { t.Fatal(err) }
			server := httptest.NewServer(handler)
			defer server.Close()
			status, body := requestDemo(t, server, "/demo-cases", testBearer)
			if status != http.StatusOK || strings.Contains(body, tc.foreign) ||
				strings.Contains(body, "event-") || strings.Contains(body, "demo-hidden-source") ||
				strings.Contains(body, "demo-tenant") || strings.Contains(body, testBearer) {
				t.Fatalf("Case catalogue disclosed foreign source or identity: %d %s", status, body)
			}
			var result struct { Cases []string `json:"cases"` }
			if err := json.Unmarshal([]byte(body), &result); err != nil { t.Fatal(err) }
			if !reflect.DeepEqual(result.Cases, []string{tc.expected}) {
				t.Fatalf("wrong authorized Cases: %+v", result)
			}
			status, body = requestDemo(t, server, "/case-timeline?case_id="+result.Cases[0]+"&limit=20", testBearer)
			if status != http.StatusOK || !strings.Contains(body, `"entries"`) ||
				strings.Contains(body, "demo-hidden-source") { 
				t.Fatalf("discovered Case could not open through ACL chain: %d %s", status, body)
			}
			status, body = requestDemo(t, server, "/case-timeline?case_id="+tc.foreign, testBearer)
			if status != http.StatusNotFound || strings.Contains(body, "event-") {
				t.Fatalf("discovery bypassed Case ACL: %d %s", status, body)
			}
		})
	}
}

func TestDemoCaseCatalogueDeniesUnauthenticatedAndInvalidRequests(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil { t.Fatal(err) }
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, bearer := range []string{"", "wrong", testBearer[:len(testBearer)-1]} {
		status, body := requestDemo(t, server, "/demo-cases", bearer)
		if status != http.StatusUnauthorized || strings.Contains(body, "case-a") {
			t.Fatalf("unauthorized Case list leaked: %d %s", status, body)
		}
	}
	for _, path := range []string{
		"/demo-cases?tenant=demo-tenant-b",
		"/demo-cases?case_id=case-a",
		"/demo-cases?principal=analyst-b",
	} {
		status, body := requestDemo(t, server, path, testBearer)
		if status != http.StatusBadRequest || strings.Contains(body, "case-a") {
			t.Fatalf("caller selection accepted: %d %s", status, body)
		}
	}
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead} {
		r := httptest.NewRequest(method, "/demo-cases", nil)
		r.Header.Set("Authorization", "Bearer "+testBearer)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != http.StatusMethodNotAllowed || strings.Contains(w.Body.String(), "case-a") {
			t.Fatalf("write or HEAD accepted: %s %d", method, w.Code)
		}
	}
	request := httptest.NewRequest(http.MethodGet, "/demo-cases", nil)
	request.Header.Set("Authorization", "Bearer "+testBearer)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	for _, header := range []string{"Cache-Control", "Content-Type", "X-Content-Type-Options", "Referrer-Policy", "Vary"} {
		if response.Header().Get(header) == "" {t.Fatalf("missing %s", header)}
	}
}

func TestDemoCaseCatalogueUsesFixtureReadersAndTenant(t *testing.T) {
	data := demoData()
	data.Cases = append(data.Cases, timeline.MemoryTimelineCase{
		ID: "private-case", Tenant: "demo-tenant-a", Environment: "lab",
		Readers: nil,
	})
	// The Case is synthetic only; this route itself must never consult
	// caller-chosen tenant or accept the record's foreign reader.
	authenticate := func(r *http.Request) (string, error) {return "analyst-a", nil}
	handler := newDemoCaseListHandler(authenticate, data)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "private-case") {
		t.Fatalf("non-reader Case listed: %d %s", w.Code, w.Body.String())
	}
}

func TestDemoCaseCatalogueRejectsCancelledBodyAndUnknownPrincipal(t *testing.T) {
	handler, err := newDemoHandler(testBearer, "analyst-a")
	if err != nil { t.Fatal(err) }
	cancelled := httptest.NewRequest(http.MethodGet, "/demo-cases", nil)
	cancelled.Header.Set("Authorization", "Bearer "+testBearer)
	ctx, cancel := context.WithCancel(cancelled.Context())
	cancel()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, cancelled.WithContext(ctx))
	if w.Code != http.StatusServiceUnavailable || strings.Contains(w.Body.String(), "case-a") {
		t.Fatalf("cancelled request disclosed a Case: %d %s", w.Code, w.Body.String())
	}
	body := httptest.NewRequest(http.MethodGet, "/demo-cases", strings.NewReader("ignored-body"))
	body.Header.Set("Authorization", "Bearer "+testBearer)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, body)
	if w.Code != http.StatusBadRequest || strings.Contains(w.Body.String(), "case-a") {
		t.Fatalf("GET body was accepted: %d %s", w.Code, w.Body.String())
	}
	unknown := newDemoCaseListHandler(func(*http.Request) (string, error) {
		return "unregistered-principal", nil
	}, demoData())
	w = httptest.NewRecorder()
	unknown.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
	if w.Code != http.StatusUnauthorized || strings.Contains(w.Body.String(), "case-a") {
		t.Fatalf("unregistered principal listed a Case: %d %s", w.Code, w.Body.String())
	}
	bad := newDemoCaseListHandler(func(*http.Request) (string, error) {
		return "analyst-a", errors.New("credential revoked")
	}, demoData())
	w = httptest.NewRecorder()
	bad.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
	if w.Code != http.StatusUnauthorized { t.Fatalf("revoked session accepted: %d", w.Code) }
}

func TestDemoCaseCatalogueEmptyAndSortedFixtureCases(t *testing.T) {
	data := demoData()
	for i := range data.Cases {
		data.Cases[i].Readers = nil
	}
	read := func(data timeline.MemoryTimelineData) (int, []string) {
		t.Helper()
		h := newDemoCaseListHandler(func(*http.Request) (string, error) {
			return "analyst-a", nil
		}, data)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/demo-cases", nil))
		var result struct { Cases []string `json:"cases"` }
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil { t.Fatal(err) }
		return w.Code, result.Cases
	}
	code, cases := read(data)
	if code != http.StatusOK || cases == nil || len(cases) != 0 {
		t.Fatalf("empty authorized catalogue not represented as []: %d %v", code, cases)
	}
	data.Cases = append(data.Cases,
		timeline.MemoryTimelineCase{ID:"z-case",Tenant:"demo-tenant-a",Environment:"lab",Readers:[]string{"analyst-a"}},
		timeline.MemoryTimelineCase{ID:"c-case",Tenant:"demo-tenant-a",Environment:"lab",Readers:[]string{"analyst-a"}},
	)
	code, cases = read(data)
	if code != http.StatusOK || !reflect.DeepEqual(cases, []string{"c-case","z-case"}) {
		t.Fatalf("catalogue result order or visibility: %d %v", code, cases)
	}
}
