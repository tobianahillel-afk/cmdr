package workqueue

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type alertOriginFixtureOwner struct {
	*fakeOwner
	links []IncidentAlertReference
	alerts map[string]*AuthorizedIncidentAlert
	linkError error
	alertError error
	linkCalls int
	alertCalls int
}

func (o *alertOriginFixtureOwner) ReadAuthorizedIncidentAlertReferences(_ context.Context, principal, tenant, environment, incidentID string) ([]IncidentAlertReference, error) {
	o.linkCalls++
	if principal != "analyst" || tenant != "t1" || environment != "e1" || incidentID != "inc-100" {
		return nil, errors.New("PRIVATE_LINK_SCOPE")
	}
	return o.links, o.linkError
}

func (o *alertOriginFixtureOwner) ReadAuthorizedIncidentAlert(_ context.Context, principal, tenant, environment, alertID string) (*AuthorizedIncidentAlert, error) {
	o.alertCalls++
	if principal != "analyst" || tenant != "t1" || environment != "e1" || alertID == "" {
		return nil, errors.New("PRIVATE_ALERT_SCOPE")
	}
	return o.alerts[alertID], o.alertError
}

func alertOriginFixture() *alertOriginFixtureOwner {
	return &alertOriginFixtureOwner{
		fakeOwner: deepReturnOwner(),
		links: []IncidentAlertReference{
			{IncidentID: "inc-100", AlertID: "alert-safe", Tenant: "t1", Environment: "e1"},
			{IncidentID: "inc-100", AlertID: "alert-hidden", Tenant: "t1", Environment: "e1"},
		},
		alerts: map[string]*AuthorizedIncidentAlert{
			"alert-safe": {
				ID: "alert-safe", Tenant: "t1", Environment: "e1",
				Severity: "medium", Status: "triaged",
				Disposition: "<script>alert('owned')</script>",
				SourceRef: "command-alert-source", SourceVersion: "v2",
				ObservedAt: "2026-10-21T10:00:00Z",
			},
			// alert-hidden is independently not authorized: nil lookup.
		},
	}
}

func TestIncidentAlertOriginExistingDetailRendersOnlyAlertReadAuthorizedMetadata(t *testing.T) {
	o := alertOriginFixture()
	res := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
		"/command/incident-detail?incident=inc-100")
	body := res.Body.String()
	if res.Code != http.StatusOK || o.linkCalls != 1 || o.alertCalls != 2 {
		t.Fatalf("independent link and Alert owner not consulted: %d, %d, %d", res.Code, o.linkCalls, o.alertCalls)
	}
	for _, visible := range []string{"alert-safe", "medium", "triaged", "command-alert-source", "v2", "&lt;script&gt;"} {
		if !strings.Contains(body, visible) {
			t.Fatalf("missing authorized source field %q", visible)
		}
	}
	for _, hidden := range []string{"alert-hidden", "<script>", "PRIVATE_", "synthetic-forbidden"} {
		if strings.Contains(body, hidden) {
			t.Fatalf("source disclosure %q", hidden)
		}
	}
	o.allowed["Incident:inc-100"] = false
	denied := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
		"/command/incident-detail?incident=inc-100")
	if denied.Code != http.StatusNotFound || o.linkCalls != 1 || o.alertCalls != 2 {
		t.Fatalf("Incident revocation failed before linked source read: %d", denied.Code)
	}
}

func TestIncidentAlertOriginMalformedLinkOrAlertFailsBeforeAnyHTML(t *testing.T) {
	scenarios := []struct {
		name string
		bad func(*alertOriginFixtureOwner)
	}{
		{"link-owner-error", func(o *alertOriginFixtureOwner) { o.linkError = errors.New("PRIVATE_LINK_FAILURE") }},
		{"wrong-incident-relation", func(o *alertOriginFixtureOwner) { o.links[0].IncidentID = "other" }},
		{"wrong-tenant-link", func(o *alertOriginFixtureOwner) { o.links[0].Tenant = "other" }},
		{"wrong-environment-link", func(o *alertOriginFixtureOwner) { o.links[0].Environment = "other" }},
		{"empty-alert-id", func(o *alertOriginFixtureOwner) { o.links[0].AlertID = "" }},
		{"hostile-alert-id", func(o *alertOriginFixtureOwner) { o.links[0].AlertID = "bad\nalert" }},
		{"duplicate-link", func(o *alertOriginFixtureOwner) { o.links[1].AlertID = "alert-safe" }},
		{"too-many-links", func(o *alertOriginFixtureOwner) {
			o.links = make([]IncidentAlertReference, 21)
		}},
		{"alert-owner-error", func(o *alertOriginFixtureOwner) { o.alertError = errors.New("PRIVATE_ALERT_FAILURE") }},
		{"wrong-source-alert", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].ID = "other" }},
		{"wrong-source-tenant", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].Tenant = "other" }},
		{"wrong-source-environment", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].Environment = "other" }},
		{"missing-source-reference", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].SourceRef = "" }},
		{"bad-source-version", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].SourceVersion = "v1\nspoofed" }},
		{"bad-observation-time", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].ObservedAt = "yesterday" }},
		{"oversized-disposition", func(o *alertOriginFixtureOwner) { o.alerts["alert-safe"].Disposition = strings.Repeat("x", 257) }},
	}
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			o := alertOriginFixture()
			tc.bad(o)
			res := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), o),
				"/command/incident-detail?incident=inc-100")
			body := res.Body.String()
			if res.Code != http.StatusServiceUnavailable ||
				strings.Contains(body, "Incident Detail — aperçu") ||
				strings.Contains(body, "alert-safe") ||
				strings.Contains(body, "PRIVATE_") {
				t.Fatalf("invalid or revoked Alert source disclosed HTML: %d %s", res.Code, body)
			}
		})
	}
}

func TestIncidentAlertOriginOptionalLegacyAndCanceledContexts(t *testing.T) {
	legacy := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), deepReturnOwner()),
		"/command/incident-detail?incident=inc-100")
	if legacy.Code != http.StatusOK || !strings.Contains(legacy.Body.String(), "Origine Alert non disponible") {
		t.Fatalf("legacy owner was not partial: %d", legacy.Code)
	}
	o := alertOriginFixture()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := readAuthorizedIncidentAlertOrigins(ctx, o, trustedHost().scope, "inc-100")
	if err == nil || len(result) != 0 || o.linkCalls != 0 {
		t.Fatalf("canceled scope reached Alert authority: %v %+v", err, result)
	}
}
