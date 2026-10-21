package workqueue

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type overviewFixtureOwner struct {
	*fakeOwner
	overview *AuthorizedIncidentOverview
	err error
	calls int
	principal string
}

func (o *overviewFixtureOwner) ReadAuthorizedIncidentOverview(
	_ context.Context, principal, tenant, environment, incidentID string,
) (*AuthorizedIncidentOverview, error) {
	o.calls++
	o.principal = principal
	if tenant != "t1" || environment != "e1" || incidentID != "inc-100" {
		return nil, errors.New("unauthorized source argument")
	}
	return o.overview, o.err
}

func trustedIncidentOverview() *AuthorizedIncidentOverview {
	return &AuthorizedIncidentOverview{
		IncidentID: "inc-100", Tenant: "t1", Environment: "e1",
		Title: "<img src=x onerror=alert(1)>",
		Priority: "high", Severity: "medium",
		BusinessImpact: "Service availability unknown",
		SourceRef: "Command/owner", SourceVersion: "v1",
		ObservedAt: "2026-10-21T02:00:00Z",
	}
}

func TestIncidentDetailCanonicalHeaderIsAuthorizedAndHTMLEscaped(t *testing.T) {
	o := &overviewFixtureOwner{fakeOwner: deepReturnOwner(), overview: trustedIncidentOverview()}
	h := NewIncidentDetailHTTPHandler(trustedHost(), o)
	res := detailGET(h, "/command/incident-detail?incident=inc-100")
	body := res.Body.String()
	if res.Code != http.StatusOK || o.calls != 1 || o.principal != "analyst" {
		t.Fatalf("header was not independently requested from Command owner: %d calls=%d", res.Code, o.calls)
	}
	for _, field := range []string{"high", "medium", "Service availability unknown", "Command/owner", "v1", "02:00:00Z", "&lt;img"} {
		if !strings.Contains(body, field) {
			t.Fatalf("authorized source field absent: %s", field)
		}
	}
	if strings.Contains(body, "<img src=x") || strings.Contains(body, "onerror=alert(1)>") {
		t.Fatal("Command-provided title was interpreted as HTML")
	}
	if !strings.Contains(body, "Statut du dossier : partiel") {
		t.Fatal("reading the Command header must not imply a fully delivered Incident Detail")
	}
	o.allowed["Incident:inc-100"] = false
	denied := detailGET(h, "/command/incident-detail?incident=inc-100")
	if denied.Code != http.StatusNotFound || o.calls != 1 ||
		strings.Contains(denied.Body.String(), "Service availability") {
		t.Fatal("revoked Incident reached its independently authorized metadata source")
	}
}

func TestIncidentDetailOptionalCanonicalHeaderNeverInventsMissingFields(t *testing.T) {
	o := &overviewFixtureOwner{fakeOwner: deepReturnOwner()}
	partial := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
		"/command/incident-detail?incident=inc-100")
	if partial.Code != http.StatusOK ||
		!strings.Contains(partial.Body.String(), "Informations métier indisponibles") ||
		strings.Contains(partial.Body.String(), "Command/owner") {
		t.Fatalf("missing source was invented: %d %s", partial.Code, partial.Body.String())
	}
	o.overview = trustedIncidentOverview()
	o.overview.Title = ""
	o.overview.Priority = ""
	o.overview.Severity = ""
	o.overview.BusinessImpact = ""
	present := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
		"/command/incident-detail?incident=inc-100")
	if present.Code != http.StatusOK || !strings.Contains(present.Body.String(), "inconnu — non fourni") ||
		!strings.Contains(present.Body.String(), "inconnue — non fournie") {
		t.Fatalf("partial owner data was incorrectly filled in: %d %s", present.Code, present.Body.String())
	}
	// An old owner without the optional capability is still supported.
	legacy := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), deepReturnOwner()),
		"/command/incident-detail?incident=inc-100")
	if legacy.Code != http.StatusOK || !strings.Contains(legacy.Body.String(), "Informations métier indisponibles") {
		t.Fatalf("legacy source could not serve partial Incident Detail: %d %s", legacy.Code, legacy.Body.String())
	}
}

func TestIncidentDetailRejectedSourceNeverLeaksOrReturnsPartialHTML(t *testing.T) {
	testCases := []struct {
		name string
		change func(*overviewFixtureOwner)
	}{
		{"wrong-tenant", func(o *overviewFixtureOwner) { o.overview.Tenant = "t2" }},
		{"wrong-environment", func(o *overviewFixtureOwner) { o.overview.Environment = "e2" }},
		{"wrong-incident", func(o *overviewFixtureOwner) { o.overview.IncidentID = "hidden-incident" }},
		{"malformed-time", func(o *overviewFixtureOwner) { o.overview.ObservedAt = "yesterday" }},
		{"no-source-ref", func(o *overviewFixtureOwner) { o.overview.SourceRef = "" }},
		{"no-source-version", func(o *overviewFixtureOwner) { o.overview.SourceVersion = "" }},
		{"oversized-title", func(o *overviewFixtureOwner) { o.overview.Title = strings.Repeat("a", 257) }},
		{"control-in-impact", func(o *overviewFixtureOwner) { o.overview.BusinessImpact = "secret\nforged field" }},
		{"provider-error", func(o *overviewFixtureOwner) { o.err = errors.New("SECRET-PROVIDER-INTERNAL") }},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := &overviewFixtureOwner{fakeOwner: deepReturnOwner(), overview: trustedIncidentOverview()}
			tc.change(o)
			res := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
				"/command/incident-detail?incident=inc-100")
			if res.Code != http.StatusServiceUnavailable || o.calls != 1 ||
				strings.Contains(res.Body.String(), "Incident Detail — aperçu") ||
				strings.Contains(res.Body.String(), "SECRET") ||
				strings.Contains(res.Body.String(), "Service availability") {
				t.Fatalf("invalid data produced partial/secret HTML: %d %s", res.Code, res.Body.String())
			}
		})
	}
}
