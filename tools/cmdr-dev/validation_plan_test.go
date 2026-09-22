package main

import "testing"

func planningCatalog() CheckCatalog {
	catalog := minimalCheckCatalog()
	catalog.Checks = append(catalog.Checks,
		ValidationCheck{
			ID: "CHK-SPEC", Title: "Spec", ExecutorKey: "spec-index",
			RiskDomains: []string{"product-spec"}, CostTier: "standard",
			Mandatory: true, TriggerPaths: []string{"cmdr-product-spec/**"},
			Prerequisites: []string{"CHK-VET"}, EvidenceKind: "artifact",
		},
	)
	return catalog
}

func TestValidationPlanIncludesSafetyFloorAndDirectTriggers(t *testing.T) {
	impact := ImpactReport{
		WorkUnit: "E3-X-001",
		ChangedPaths: []string{"cmdr-product-spec/x.md"},
		RiskDomains: []string{"product-spec"},
	}
	plan, err := buildValidationPlan(impact, planningCatalog())
	if err != nil {
		t.Fatal(err)
	}
	if plan.Tier != "standard" {
		t.Fatalf("expected standard tier, got %s", plan.Tier)
	}
	for _, id := range []string{"CHK-FORMAT", "CHK-VET", "CHK-SPEC"} {
		if !planContainsCheck(plan, id) {
			t.Fatalf("plan omitted %s: %#v", id, plan)
		}
	}
	if err := validateValidationPlan(plan, planningCatalog()); err != nil {
		t.Fatal(err)
	}
}

func TestValidationPlanStrictSelectsEveryMandatoryCheck(t *testing.T) {
	impact := ImpactReport{
		WorkUnit: "E3-X-001",
		ChangedPaths: []string{"mystery/file"},
		RiskDomains: []string{"unknown", "security"},
		HighRisk: true,
		UnknownPaths: []string{"mystery/file"},
	}
	catalog := planningCatalog()
	plan, err := buildValidationPlan(impact, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Tier != "strict" {
		t.Fatalf("expected strict tier, got %s", plan.Tier)
	}
	for _, check := range catalog.Checks {
		if check.Mandatory && !planContainsCheck(plan, check.ID) {
			t.Fatalf("strict plan omitted mandatory %s", check.ID)
		}
	}
}

func TestValidationPlanClosesPrerequisites(t *testing.T) {
	impact := ImpactReport{
		WorkUnit: "E3-X-001",
		ChangedPaths: []string{"cmdr-product-spec/x.md"},
		RiskDomains: []string{"product-spec"},
	}
	plan, err := buildValidationPlan(impact, planningCatalog())
	if err != nil {
		t.Fatal(err)
	}
	spec := planCheckPosition(plan, "CHK-SPEC")
	vet := planCheckPosition(plan, "CHK-VET")
	format := planCheckPosition(plan, "CHK-FORMAT")
	if format < 0 || vet < 0 || spec < 0 || !(format < vet && vet < spec) {
		t.Fatalf("prerequisite order is wrong: %#v", plan.SelectedChecks)
	}
}

func TestManifestScopeSelectsDomainChecks(t *testing.T) {
	impact := ImpactReport{
		WorkUnit: "E3-X-001",
		ChangedPaths: []string{"work/graph.json"},
		RiskDomains: []string{"product-spec"},
		Evidence: []ImpactEvidence{{Domain:"product-spec",Reason:"manifest-product-references",Paths:[]string{"E3-X-001"}}},
	}
	plan, err := buildValidationPlan(impact, planningCatalog())
	if err != nil {
		t.Fatal(err)
	}
	if !planContainsCheck(plan, "CHK-SPEC") {
		t.Fatalf("manifest-scope product check not selected: %#v", plan)
	}
}

func planContainsCheck(plan ValidationPlan, id string) bool {
	return planCheckPosition(plan, id) >= 0
}

func planCheckPosition(plan ValidationPlan, id string) int {
	for i, check := range plan.SelectedChecks {
		if check.ID == id {
			return i
		}
	}
	return -1
}
