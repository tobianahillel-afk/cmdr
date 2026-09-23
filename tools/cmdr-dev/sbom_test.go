package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCycloneDXSBOMIsDeterministicAndSourceBound(t *testing.T) {
	sha := strings.Repeat("a", 40)
	inputs := []sbomInputComponent{
		{Kind: "application", Ecosystem: "go", Name: "golang.org/x/vuln/cmd/govulncheck", Version: "v1.8.0", Manifest: developmentToolRegistryPath, Scope: "development", Direct: true, Commit: strings.Repeat("b", 40), License: "BSD-3-Clause"},
		{Kind: "application", Ecosystem: "go", Name: "github.com/securego/gosec/v2/cmd/gosec", Version: "v2.28.0", Manifest: developmentToolRegistryPath, Scope: "development", Direct: true, Commit: strings.Repeat("c", 40), License: "Apache-2.0"},
	}
	first, runtimeCount, developmentCount, err := buildCycloneDXBOM(sha, inputs)
	if err != nil {
		t.Fatal(err)
	}
	second, _, _, err := buildCycloneDXBOM(sha, inputs)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Fatal("SBOM output is not deterministic")
	}
	if first.Schema != cycloneDXSchemaURL || first.SpecVersion != "1.7" || first.BOMFormat != "CycloneDX" {
		t.Fatalf("unexpected CycloneDX identity: %#v", first)
	}
	if first.Metadata.Component.Version != sha || first.Metadata.Component.BOMRef != "cmdr:source:"+sha {
		t.Fatalf("SBOM root is not source-SHA bound: %#v", first.Metadata.Component)
	}
	if runtimeCount != 0 || developmentCount != 2 {
		t.Fatalf("unexpected scope counts runtime=%d dev=%d", runtimeCount, developmentCount)
	}
	if strings.Contains(string(a), "timestamp") || strings.Contains(string(a), "serialNumber") {
		t.Fatal("deterministic SBOM must not contain timestamp or random serial number")
	}
}

func TestCycloneDXDevelopmentToolsAreExcludedAndScoped(t *testing.T) {
	sha := strings.Repeat("d", 40)
	input := sbomInputComponent{
		Kind: "application", Ecosystem: "go", Name: "example.org/tool", Version: "v1.2.3",
		Manifest: developmentToolRegistryPath, Scope: "development", Direct: true,
	}
	bom, _, _, err := buildCycloneDXBOM(sha, []sbomInputComponent{input})
	if err != nil {
		t.Fatal(err)
	}
	if len(bom.Components) != 1 || bom.Components[0].Scope != "excluded" {
		t.Fatalf("development component must be CycloneDX excluded: %#v", bom.Components)
	}
	props := map[string]string{}
	for _, prop := range bom.Components[0].Properties {
		props[prop.Name] = prop.Value
	}
	if props["cmdr:dependency:scope"] != "development" || props["cmdr:dependency:direct"] != "true" {
		t.Fatalf("development/direct scope is not explicit: %#v", props)
	}
}

func TestSBOMRejectsUnresolvedVersions(t *testing.T) {
	cases := []sbomInputComponent{
		{Kind: "library", Ecosystem: "npm", Name: "pkg", Version: "^1.2.3", Manifest: "package.json", Scope: "development", Direct: true},
		{Kind: "library", Ecosystem: "vcpkg", Name: "lib", Version: "1.0.0", Manifest: "vcpkg.json", Scope: "runtime", Direct: true},
		{Kind: "library", Ecosystem: "go", Name: "example.org/x", Version: "latest", Manifest: "go.mod", Scope: "development", Direct: true},
	}
	for _, input := range cases {
		if err := validateResolvedSBOMInput(input); err == nil {
			t.Fatalf("expected unresolved identity rejection for %#v", input)
		}
	}
}

func TestSBOMAcceptsExactGoAndNPMVersions(t *testing.T) {
	for _, input := range []sbomInputComponent{
		{Kind: "library", Ecosystem: "go", Name: "example.org/x", Version: "v1.2.3", Manifest: "go.mod", Scope: "runtime", Direct: true},
		{Kind: "library", Ecosystem: "npm", Name: "pkg", Version: "1.2.3-beta.1", Manifest: "package.json", Scope: "development", Direct: true},
	} {
		if err := validateResolvedSBOMInput(input); err != nil {
			t.Fatalf("unexpected exact-version rejection: %v", err)
		}
	}
}

func TestSBOMComponentReferenceChangesWithIdentity(t *testing.T) {
	base := sbomInputComponent{Kind: "library", Ecosystem: "go", Name: "example.org/x", Version: "v1.0.0", Manifest: "go.mod", Scope: "runtime", Direct: true}
	other := base
	other.Version = "v1.0.1"
	if sbomComponentRef(base) == sbomComponentRef(other) {
		t.Fatal("component reference must change when resolved identity changes")
	}
}
