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


func TestIncidentOverviewValidationRejectsEachUntrustedFieldIndependently(t *testing.T) {
	scope := trustedHost().scope
	valid, err := validateIncidentOverview(trustedIncidentOverview(), scope, "inc-100")
	if err != nil || valid == nil {
		t.Fatalf("complete authorized overview rejected: %v", err)
	}
	absent, err := validateIncidentOverview(nil, scope, "inc-100")
	if err != nil || absent != nil {
		t.Fatalf("missing optional source should remain unknown: %+v %v", absent, err)
	}
	cases := []struct {
		name string
		change func(*AuthorizedIncidentOverview)
	}{
		{"other-incident", func(s *AuthorizedIncidentOverview) { s.IncidentID = "foreign" }},
		{"other-tenant", func(s *AuthorizedIncidentOverview) { s.Tenant = "other" }},
		{"other-environment", func(s *AuthorizedIncidentOverview) { s.Environment = "prod" }},
		{"empty-source", func(s *AuthorizedIncidentOverview) { s.SourceRef = "" }},
		{"wildcard-source", func(s *AuthorizedIncidentOverview) { s.SourceRef = "*" }},
		{"untrusted-source-controls", func(s *AuthorizedIncidentOverview) { s.SourceRef = "secret\nother" }},
		{"oversized-source", func(s *AuthorizedIncidentOverview) { s.SourceRef = strings.Repeat("s", 129) }},
		{"invalid-source-utf8", func(s *AuthorizedIncidentOverview) { s.SourceRef = string([]byte{0xff}) }},
		{"empty-version", func(s *AuthorizedIncidentOverview) { s.SourceVersion = "" }},
		{"oversized-version", func(s *AuthorizedIncidentOverview) { s.SourceVersion = strings.Repeat("v", 81) }},
		{"version-control", func(s *AuthorizedIncidentOverview) { s.SourceVersion = "v1\rforged" }},
		{"invalid-version-utf8", func(s *AuthorizedIncidentOverview) { s.SourceVersion = string([]byte{0xff}) }},
		{"oversized-title", func(s *AuthorizedIncidentOverview) { s.Title = strings.Repeat("x", 257) }},
		{"title-control", func(s *AuthorizedIncidentOverview) { s.Title = "x\nsecret" }},
		{"invalid-title-utf8", func(s *AuthorizedIncidentOverview) { s.Title = string([]byte{0xff}) }},
		{"oversized-priority", func(s *AuthorizedIncidentOverview) { s.Priority = strings.Repeat("p", 65) }},
		{"invalid-priority-utf8", func(s *AuthorizedIncidentOverview) { s.Priority = string([]byte{0xff}) }},
		{"oversized-severity", func(s *AuthorizedIncidentOverview) { s.Severity = strings.Repeat("s", 65) }},
		{"invalid-severity-utf8", func(s *AuthorizedIncidentOverview) { s.Severity = string([]byte{0xff}) }},
		{"oversized-business-impact", func(s *AuthorizedIncidentOverview) { s.BusinessImpact = strings.Repeat("x", 513) }},
		{"invalid-impact-utf8", func(s *AuthorizedIncidentOverview) { s.BusinessImpact = string([]byte{0xff}) }},
		{"invalid-observed-at", func(s *AuthorizedIncidentOverview) { s.ObservedAt = "untrusted" }},
		{"zero-observed-at", func(s *AuthorizedIncidentOverview) { s.ObservedAt = "0001-01-01T00:00:00Z" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := trustedIncidentOverview()
			tc.change(in)
			out, err := validateIncidentOverview(in, scope, "inc-100")
			if err == nil || out != nil {
				t.Fatalf("untrusted source value accepted: %+v %v", out, err)
			}
		})
	}
	// A freshly verified read projection is a defensive copy; changing it
	// after validation must never rewrite the Command owner's source.
	original := trustedIncidentOverview()
	copyValue, err := validateIncidentOverview(original, scope, "inc-100")
	if err != nil || copyValue == original {
		t.Fatalf("source result not copied: %+v %v", copyValue, err)
	}
	copyValue.Title = "changed"
	if original.Title == "changed" {
		t.Fatal("validated projection mutated source owner")
	}
}

func TestIncidentOverviewReadRejectsCancellationEvenIfSourceIgnoresContext(t *testing.T) {
	owner := &overviewFixtureOwner{fakeOwner: deepReturnOwner(), overview: trustedIncidentOverview()}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	out, err := readAuthorizedIncidentOverview(ctx, owner, trustedHost().scope, "inc-100")
	if err == nil || out != nil || owner.calls != 1 {
		t.Fatalf("canceled owner result escaped: %v, %+v, calls=%d", err, out, owner.calls)
	}
	owner.err = errors.New("SENSITIVE-OWNER-FAILURE")
	_, err = readAuthorizedIncidentOverview(context.Background(), owner, trustedHost().scope, "inc-100")
	if err == nil || strings.Contains(err.Error(), "SENSITIVE") {
		t.Fatalf("internal owner error exposed: %v", err)
	}
}
