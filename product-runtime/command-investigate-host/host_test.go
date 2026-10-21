package commandinvestigatehost

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
)

// These are independent synthetic host authorities, not a production IAM or
// a claim that Command Incident.read confers Investigate Case.read.
type testCommandSession struct{}

func (testCommandSession) ResolveSession(_ context.Context, r *http.Request) (queue.SessionScope, error) {
	switch r.Header.Get("Authorization") {
	case "Bearer analyst":
		return queue.SessionScope{Principal: "analyst-a", Tenant: "t1", Environment: "prod"}, nil
	case "Bearer commander":
		return queue.SessionScope{Principal: "command-only", Tenant: "t1", Environment: "prod"}, nil
	default:
		return queue.SessionScope{}, errors.New("unauthenticated")
	}
}

type testIncidentOwner struct{ lookups int }

func (*testIncidentOwner) List(_ context.Context, tenant, environment string) ([]queue.Item, error) {
	if tenant != "t1" || environment != "prod" {
		return nil, errors.New("out of scope")
	}
	return []queue.Item{{Kind: queue.Incident, ID: "incident-a",
		Tenant: "t1", Environment: "prod", Owner: "incident-commander",
		Status: "active", UpdatedAt: "2026-10-21T02:00:00Z"}}, nil
}

func (*testIncidentOwner) CanRead(_ context.Context, principal, tenant, environment string, kind queue.Kind, id string) (bool, error) {
	return (principal == "analyst-a" || principal == "command-only") &&
		tenant == "t1" && environment == "prod" && kind == queue.Incident && id == "incident-a", nil
}

func (o *testIncidentOwner) ReadAuthorizedLinkedCases(_ context.Context, principal, tenant, environment, incidentID string) (*queue.AuthorizedLinkedCases, error) {
	o.lookups++
	if (principal != "analyst-a" && principal != "command-only") ||
		tenant != "t1" || environment != "prod" || incidentID != "incident-a" {
		return nil, errors.New("link denied")
	}
	return &queue.AuthorizedLinkedCases{Cases: []queue.LinkedCaseSummary{{
		CaseID: "case-a", Status: "investigating", LinkSource: "Shared",
	}}}, nil
}

func combinedTestHost(t *testing.T) (http.Handler, *testIncidentOwner) {
	t.Helper()
	owner := &testIncidentOwner{}
	command := queue.NewHTTPHandler(testCommandSession{}, owner)
	auth := func(r *http.Request) (string, error) {
		switch r.Header.Get("Authorization") {
		case "Bearer analyst":
			return "analyst-a", nil
		case "Bearer commander":
			return "command-only", nil
		default:
			return "", errors.New("no Investigate session")
		}
	}
	data := timeline.MemoryTimelineData{
		Principals: []timeline.MemoryTimelinePrincipal{
			{ID: "analyst-a", Tenant: "t1", Environment: "prod"},
			{ID: "command-only", Tenant: "t1", Environment: "prod"},
		},
		Cases: []timeline.MemoryTimelineCase{{
			ID: "case-a", Tenant: "t1", Environment: "prod",
			Readers: []string{"analyst-a"}, EntryIDs: []string{"event-a"},
		}},
		Entries: []timeline.MemoryTimelineEntry{{
			CaseID: "case-a",
			Value: timeline.Entry{ID: "event-a", Tenant: "t1", Kind: "observed",
				SourceTime: "2026-10-21T02:05:00Z", Timezone: "+00:00",
				SourceRef: "source-a", Provenance: "telemetry", Version: 1,
				SourcePermission: "allow"},
			Readers: []string{"analyst-a"},
		}},
		Sources: []timeline.MemoryTimelineSource{{
			Identity: timeline.SourceOwnerIdentity{Reference: "source-a",
				ObjectID: "event-owner-a", Kind: "telemetry-event",
				Tenant: "t1", Environment: "prod", CaseID: "case-a"},
			Readers: []string{"analyst-a"},
		}},
	}
	receiver, err := timeline.NewMemoryCaseTimelineHTTPHandler(auth, data)
	if err != nil {
		t.Fatal(err)
	}
	host, err := NewHandler(command, receiver)
	if err != nil {
		t.Fatal(err)
	}
	return host, owner
}

func request(h http.Handler, path, token, accept string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		r.Header.Set("Authorization", token)
	}
	if accept != "" {
		r.Header.Set("Accept", accept)
	}
	out := httptest.NewRecorder()
	h.ServeHTTP(out, r)
	return out
}

func TestRealSameOriginCommandToInvestigateNavigationIndependentlyReauthorizes(t *testing.T) {
	host, owner := combinedTestHost(t)
	start := request(host, "/command/work-queue?view=incidents&selected=Incident:incident-a",
		"Bearer analyst", "text/html")
	if start.Code != http.StatusOK || start.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Command read unavailable: %d %s", start.Code, start.Body.String())
	}
	link := "/case-timeline?case_id=" + url.QueryEscape("case-a")
	if !strings.Contains(html.UnescapeString(start.Body.String()), "href=\""+link+"\"") {
		t.Fatalf("authorized Case link not rendered in Command: %s", start.Body.String())
	}
	if owner.lookups != 1 {
		t.Fatalf("link owner was not consulted independently: %d", owner.lookups)
	}
	// Follow the actual same-origin link, with a new independent Investigate
	// session and Case/Entry/source authorization chain.
	caseRead := request(host, link, "Bearer analyst", "")
	if caseRead.Code != http.StatusOK || caseRead.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("actual Case destination unavailable: %d %s", caseRead.Code, caseRead.Body.String())
	}
	var body struct {
		Status string `json:"status"`
		Entries []struct { ID string `json:"id"`; Kind string `json:"kind"` } `json:"entries"`
	}
	if err := json.Unmarshal(caseRead.Body.Bytes(), &body); err != nil ||
		len(body.Entries) != 1 || body.Entries[0].ID != "event-a" ||
		body.Entries[0].Kind != "observed" {
		t.Fatalf("Case source owner projection not delivered: %+v, %v", body, err)
	}
	// Command-read alone still permits a link but confers NO Case.read.
	commandOnly := request(host, "/command/work-queue?selected=Incident:incident-a",
		"Bearer commander", "text/html")
	if commandOnly.Code != http.StatusOK || !strings.Contains(commandOnly.Body.String(), "case_id=case-a") {
		t.Fatalf("Command reader unexpectedly refused navigation pointer: %d", commandOnly.Code)
	}
	denied := request(host, link, "Bearer commander", "")
	if denied.Code != http.StatusNotFound || strings.Contains(denied.Body.String(), "event-a") ||
		strings.Contains(denied.Body.String(), "source-a") {
		t.Fatalf("Command privilege promoted to Investigate: %d %s", denied.Code, denied.Body.String())
	}
	for _, bad := range []string{"", "Bearer invalid"} {
		got := request(host, link, bad, "")
		if got.Code == http.StatusOK || strings.Contains(got.Body.String(), "event-a") {
			t.Fatalf("anonymous Case accessed source data: %d %s", got.Code, got.Body.String())
		}
	}
}

func TestSameOriginHostRoutesFailClosedWithoutGrantingOrRedirecting(t *testing.T) {
	if _, err := NewHandler(nil, http.NotFoundHandler()); err == nil {
		t.Fatal("accepted missing Command owner handler")
	}
	if _, err := NewHandler(http.NotFoundHandler(), nil); err == nil {
		t.Fatal("accepted missing Investigate owner handler")
	}
	host, owner := combinedTestHost(t)
	for _, path := range []string{
		"/", "/case-timeline/", "/command/work-queue/", "/command/work-queue/../case-timeline",
		"/case-timeline?case_id=case-a&tenant=t1",
	} {
		got := request(host, path, "Bearer analyst", "")
		if got.Code == http.StatusOK || strings.Contains(got.Body.String(), "event-a") ||
			strings.Contains(got.Body.String(), "source-a") {
			t.Fatalf("host routed unsupported/forged request: %s: %d %s", path, got.Code, got.Body.String())
		}
	}
	if owner.lookups != 0 {
		t.Fatalf("unexpected Command source lookups on unknown routes: %d", owner.lookups)
	}
	req := httptest.NewRequest(http.MethodPost, "/case-timeline?case_id=case-a", strings.NewReader("mutation"))
	out := httptest.NewRecorder()
	host.ServeHTTP(out, req)
	if out.Code != http.StatusMethodNotAllowed {
		t.Fatalf("mutation reached read-only host: %d", out.Code)
	}
}
