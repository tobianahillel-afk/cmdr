package main

import (
	"context"
	"testing"
)

func TestSyntheticIncidentAlertLinkAndSourceAreSeparatelyAuthorized(t *testing.T) {
	owner := demoCommandOwner{}
	cases := []struct {
		name string
		principal string
		tenant string
		environment string
		incident string
		alert string
		want bool
	}{
		{"allowed-analyst", "analyst", demoTenant, demoEnvironment, "incident-a", "alert-a", true},
		{"incident-only-commander", "commander", demoTenant, demoEnvironment, "incident-a", "alert-a", false},
		{"unrecognized-principal", "unknown", demoTenant, demoEnvironment, "incident-a", "alert-a", false},
		{"other-tenant", "analyst", "other", demoEnvironment, "incident-a", "alert-a", false},
		{"other-environment", "analyst", demoTenant, "prod", "incident-a", "alert-a", false},
		{"other-incident", "analyst", demoTenant, demoEnvironment, "other-incident", "alert-a", false},
		{"other-alert", "analyst", demoTenant, demoEnvironment, "incident-a", "other-alert", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			links, err := owner.ReadAuthorizedIncidentAlertReferences(context.Background(),
				tc.principal, tc.tenant, tc.environment, tc.incident)
			if err != nil {
				t.Fatal(err)
			}
			wantLink := tc.principal == "analyst" && tc.tenant == demoTenant &&
				tc.environment == demoEnvironment && tc.incident == "incident-a"
			if (len(links) != 0) == wantLink {
				t.Fatalf("unauthorized relation exposed: %+v", links)
			}
			// This second check is deliberately independent: even a trusted
			// link never conveys permission to inspect an Alert.
			alert, err := owner.ReadAuthorizedIncidentAlert(context.Background(),
				tc.principal, tc.tenant, tc.environment, tc.alert)
			if err != nil || (alert != nil) != (tc.principal == "analyst" &&
				tc.tenant == demoTenant && tc.environment == demoEnvironment &&
				tc.alert == "alert-a") {
				t.Fatalf("Alert.read relation authority confused: %+v %v", alert, err)
			}
			if tc.want && (len(links) != 1 || alert == nil || alert.ID != "alert-a" ||
				alert.SourceRef != "synthetic-command-alert-owner" || alert.SourceVersion != "v1" ||
				alert.ObservedAt == "") {
				t.Fatalf("incomplete authorized Alert provenance: %+v %+v", links, alert)
			}
		})
	}
}

func TestSyntheticIncidentAlertCanceledOwnerCallsDenyWithoutSource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	owner := demoCommandOwner{}
	links, err := owner.ReadAuthorizedIncidentAlertReferences(ctx,
		"analyst", demoTenant, demoEnvironment, "incident-a")
	if err == nil || len(links) != 0 {
		t.Fatalf("canceled link owner returned source IDs: %+v, %v", links, err)
	}
	alert, err := owner.ReadAuthorizedIncidentAlert(ctx,
		"analyst", demoTenant, demoEnvironment, "alert-a")
	if err == nil || alert != nil {
		t.Fatalf("canceled Alert owner disclosed metadata: %+v, %v", alert, err)
	}
}
