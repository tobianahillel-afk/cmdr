package main

import (
	"context"
	"testing"
)

func TestSyntheticIncidentOverviewOwnerStrictlyRequiresIndependentScope(t *testing.T) {
	owner := demoCommandOwner{}
	for _, tc := range []struct {
		name, principal, tenant, environment, incident string
		wantVisible bool
	}{
		{"authorized-analyst", "analyst", demoTenant, demoEnvironment, "incident-a", true},
		{"command-only-reader", "commander", demoTenant, demoEnvironment, "incident-a", false},
		{"unknown-principal", "intruder", demoTenant, demoEnvironment, "incident-a", false},
		{"different-tenant", "analyst", "other-tenant", demoEnvironment, "incident-a", false},
		{"different-environment", "analyst", demoTenant, "prod", "incident-a", false},
		{"different-incident", "analyst", demoTenant, demoEnvironment, "foreign-incident", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, err := owner.ReadAuthorizedIncidentOverview(context.Background(),
				tc.principal, tc.tenant, tc.environment, tc.incident)
			if err != nil || (res != nil) != tc.wantVisible {
				t.Fatalf("invalid source projection disclosure, visible=%v, err=%v", res != nil, err)
			}
			if res != nil && (res.IncidentID != "incident-a" ||
				res.Tenant != demoTenant || res.Environment != demoEnvironment ||
				res.SourceRef != "synthetic-command-incident-owner" ||
				res.Title == "" || res.Priority == "" || res.Severity == "" ||
				res.BusinessImpact == "" || res.SourceVersion == "" || res.ObservedAt == "") {
				t.Fatalf("authorized canonical source identity incomplete: %+v", res)
			}
		})
	}
}

func TestSyntheticIncidentOverviewOwnerCanceledContextNeverReturnsFields(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := (demoCommandOwner{}).ReadAuthorizedIncidentOverview(ctx,
		"analyst", demoTenant, demoEnvironment, "incident-a")
	if err == nil || got != nil {
		t.Fatalf("canceled source lookup returned business metadata: %+v, %v", got, err)
	}
}
