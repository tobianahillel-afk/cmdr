package main

import "testing"

func impactTestRegistry() ArchitectureRegistry {
	return ArchitectureRegistry{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unregistered-path",
		Boundaries: []ArchitectureBoundary{
			{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
			{ID: "engineering", Kind: "engineering", Roots: []string{"engineering/**", "tools/**", "AGENTS.md", "AI_START_HERE.md"}, MutableByImplementation: true, MayDependOn: []string{"product-spec"}},
			{ID: "execution", Kind: "execution", Roots: []string{"work/**"}, MutableByImplementation: true, MayDependOn: []string{"engineering", "product-spec"}},
			{ID: "ci", Kind: "ci", Roots: []string{".github/**"}, MutableByImplementation: true, MayDependOn: []string{"engineering", "execution"}},
		},
	}
}

func impactTestCatalog() CheckCatalog {
	catalog := minimalCheckCatalog()
	catalog.Checks = append(catalog.Checks,
		ValidationCheck{
			ID: "CHK-MANIFESTS", Title: "Manifests", ExecutorKey: "validate-manifests",
			RiskDomains: []string{"work-governance"}, CostTier: "fast", Mandatory: true,
			TriggerPaths: []string{"work/**"}, EvidenceKind: "invariant",
		},
		ValidationCheck{
			ID: "CHK-SPEC-INVENTORY", Title: "Spec", ExecutorKey: "spec-index",
			RiskDomains: []string{"product-spec", "coverage-traceability"}, CostTier: "standard", Mandatory: true,
			TriggerPaths: []string{"cmdr-product-spec/**"}, EvidenceKind: "generated-artifact",
		},
	)
	return catalog
}

func TestComputeImpactProductSpecChange(t *testing.T) {
	report, err := computeImpact("E3-X-001", []string{"cmdr-product-spec/06-command/example.md"}, validManifestV2(), impactTestCatalog(), impactTestRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if !containsString(report.RiskDomains, "product-spec") || !containsString(report.RiskDomains, "coverage-traceability") {
		t.Fatalf("missing product domains: %#v", report)
	}
	if len(report.UnknownPaths) != 0 {
		t.Fatalf("product spec path should be known: %#v", report.UnknownPaths)
	}
}

func TestComputeImpactUnknownPathEscalates(t *testing.T) {
	report, err := computeImpact("E3-X-001", []string{"mystery/file.xyz"}, validManifestV2(), impactTestCatalog(), impactTestRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if !report.HighRisk || !containsString(report.RiskDomains, "unknown") || !containsString(report.RiskDomains, "security") {
		t.Fatalf("unknown path was not escalated: %#v", report)
	}
	if len(report.UnknownPaths) != 1 || report.UnknownPaths[0] != "mystery/file.xyz" {
		t.Fatalf("unexpected unknown paths: %#v", report.UnknownPaths)
	}
}

func TestComputeImpactCIChangeIsHighRisk(t *testing.T) {
	report, err := computeImpact("E3-X-001", []string{".github/workflows/engineering-kernel.yml"}, validManifestV2(), impactTestCatalog(), impactTestRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if !report.HighRisk || !containsString(report.RiskDomains, "security") || !containsString(report.RiskDomains, "work-governance") {
		t.Fatalf("CI change was not high risk: %#v", report)
	}
}

func TestComputeImpactManifestSecurityRefsEscalate(t *testing.T) {
	manifest := validManifestV2()
	manifest.ProductRefs.Permissions = []string{"perm.example"}
	report, err := computeImpact("E3-X-001", []string{"work/graph.json"}, manifest, impactTestCatalog(), impactTestRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if !report.HighRisk || !containsString(report.RiskDomains, "security") || !containsString(report.RiskDomains, "product-spec") {
		t.Fatalf("security product refs were not escalated: %#v", report)
	}
}

func TestNormalizeChangedPathsRejectsEscape(t *testing.T) {
	if _, err := normalizeChangedPaths([]string{"../outside"}); err == nil {
		t.Fatal("expected repository escape rejection")
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
