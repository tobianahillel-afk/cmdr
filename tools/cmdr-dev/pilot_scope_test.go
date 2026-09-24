package main

import "testing"

func validPilotScopeForTest() PilotScope {
	return PilotScope{
		SchemaVersion: 1,
		SliceID:       "PILOT-CONTEXT-001",
		Title:         "Context",
		Capability:    "CAP-SET-004",
		CapabilitySource: "cmdr-product-spec/cap.md",
		Requirements: []string{"REQ-PROD-008"},
		Permissions: []string{"perm.platform-settings.tenant.read"},
		ImplementationContracts: []string{"cmdr-product-spec/17-implementation-contracts/context.md"},
		CanonicalObjects: []string{"cmdr-product-spec/05-domain-model/objects/tenant.md"},
		DependencySources: []string{"cmdr-product-spec/04-experience-architecture/context-preservation.md"},
		MVPSource: "cmdr-product-spec/18-roadmap-and-releases/mvp-scope.md",
		PilotPlanSource: "cmdr-product-spec/18-roadmap-and-releases/pilot-plan.md",
		DependencyRoadmapSource: "cmdr-product-spec/18-roadmap-and-releases/dependency-roadmap.md",
		ActionClasses: []string{"0","1"},
		IncludedBehaviors: []string{"preserve tenant"},
		ExcludedBehaviors: []string{"no mutation"},
		RuntimeIntent: PilotRuntimeIntent{
			Component:"context-envelope", BoundaryKind:"product-runtime",
			AuthorizationModel:"caller supplied", ExternalRuntimeDependencies:[]string{},
		},
		SelectionRationale: []string{"bounded"},
	}
}

func TestPilotScopeShapeAcceptsReadOnlyBoundedSlice(t *testing.T) {
	if err := validatePilotScopeShape(validPilotScopeForTest()); err != nil {
		t.Fatal(err)
	}
}

func TestPilotScopeShapeRejectsOpenDecision(t *testing.T) {
	scope := validPilotScopeForTest()
	scope.OpenDecisions = []string{"OPEN-013"}
	if err := validatePilotScopeShape(scope); err == nil {
		t.Fatal("expected unresolved open decision rejection")
	}
}

func TestPilotScopeShapeRejectsMutationClass(t *testing.T) {
	scope := validPilotScopeForTest()
	scope.ActionClasses = []string{"0","2"}
	if err := validatePilotScopeShape(scope); err == nil {
		t.Fatal("expected Class 2 rejection")
	}
}

func TestPilotScopeShapeRejectsExternalRuntimeDependency(t *testing.T) {
	scope := validPilotScopeForTest()
	scope.RuntimeIntent.ExternalRuntimeDependencies = []string{"example.org/runtime"}
	if err := validatePilotScopeShape(scope); err == nil {
		t.Fatal("expected external runtime dependency rejection")
	}
}

func TestSameStringSetIsOrderIndependent(t *testing.T) {
	if !sameStringSet([]string{"b","a"}, []string{"a","b"}) {
		t.Fatal("expected equal sets")
	}
	if sameStringSet([]string{"a"}, []string{"a","b"}) {
		t.Fatal("expected unequal sets")
	}
}
