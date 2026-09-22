package main

import "testing"

func minimalCheckCatalog() CheckCatalog {
	return CheckCatalog{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unknown-check",
		RiskDomains: []string{
			"code-quality", "unit-correctness", "product-spec", "coverage-traceability",
			"work-governance", "architecture", "runtime-dependencies", "security",
			"agent-context", "repository-health", "unknown",
		},
		CostTiers: []string{"fast", "standard", "heavy"},
		Checks: []ValidationCheck{
			{
				ID: "CHK-FORMAT", Title: "Format", ExecutorKey: "gofmt",
				RiskDomains: []string{"code-quality"}, CostTier: "fast",
				Mandatory: true, AlwaysOnPR: true,
				TriggerPaths: []string{"tools/cmdr-dev/**"}, EvidenceKind: "process-exit",
			},
			{
				ID: "CHK-VET", Title: "Vet", ExecutorKey: "go-vet",
				RiskDomains: []string{"code-quality"}, CostTier: "fast",
				Mandatory: true, AlwaysOnPR: true,
				TriggerPaths: []string{"tools/cmdr-dev/**"}, Prerequisites: []string{"CHK-FORMAT"},
				EvidenceKind: "process-exit",
			},
		},
	}
}

func TestValidateCheckCatalog(t *testing.T) {
	if err := validateCheckCatalog(minimalCheckCatalog()); err != nil {
		t.Fatal(err)
	}
}

func TestCheckCatalogRejectsUnknownExecutor(t *testing.T) {
	catalog := minimalCheckCatalog()
	catalog.Checks[0].ExecutorKey = "shell:anything"
	if err := validateCheckCatalog(catalog); err == nil {
		t.Fatal("expected unknown executor rejection")
	}
}

func TestCheckCatalogRejectsUnknownRiskDomain(t *testing.T) {
	catalog := minimalCheckCatalog()
	catalog.Checks[0].RiskDomains = []string{"made-up"}
	if err := validateCheckCatalog(catalog); err == nil {
		t.Fatal("expected unknown risk-domain rejection")
	}
}

func TestCheckCatalogRejectsPrerequisiteCycle(t *testing.T) {
	catalog := minimalCheckCatalog()
	catalog.Checks[0].Prerequisites = []string{"CHK-VET"}
	if err := validateCheckCatalog(catalog); err == nil {
		t.Fatal("expected prerequisite cycle rejection")
	}
}

func TestCheckCatalogRejectsDuplicateExecutor(t *testing.T) {
	catalog := minimalCheckCatalog()
	catalog.Checks[1].ExecutorKey = "gofmt"
	if err := validateCheckCatalog(catalog); err == nil {
		t.Fatal("expected duplicate executor rejection")
	}
}
