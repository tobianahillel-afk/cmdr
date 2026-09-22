package main

import "testing"

func validArchitectureRegistry() ArchitectureRegistry {
	return ArchitectureRegistry{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unregistered-path",
		Boundaries: []ArchitectureBoundary{
			{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}, MutableByImplementation: false},
			{ID: "engineering", Kind: "engineering", Roots: []string{"engineering/**", "tools/**"}, MutableByImplementation: true, MayDependOn: []string{"product-spec"}},
			{ID: "execution", Kind: "execution", Roots: []string{"work/**"}, MutableByImplementation: true, MayDependOn: []string{"engineering"}},
		},
	}
}

func TestValidateArchitectureRegistry(t *testing.T) {
	if err := validateArchitectureRegistry(validArchitectureRegistry()); err != nil {
		t.Fatal(err)
	}
}

func TestArchitectureRegistryRejectsOverlappingRoots(t *testing.T) {
	registry := validArchitectureRegistry()
	registry.Boundaries = append(registry.Boundaries,
		ArchitectureBoundary{ID: "nested", Kind: "engineering", Roots: []string{"engineering/nested/**"}, MutableByImplementation: true})
	if err := validateArchitectureRegistry(registry); err == nil {
		t.Fatal("expected overlapping-root rejection")
	}
}

func TestManifestPathOwnershipRejectsUnknownPath(t *testing.T) {
	manifest := validManifestV2()
	manifest.AllowedPaths = []string{"src/product/**"}
	if err := validateManifestPathOwnership(manifest, validArchitectureRegistry()); err == nil {
		t.Fatal("expected unknown-path rejection")
	}
}

func TestManifestPathOwnershipRejectsProductSpecMutation(t *testing.T) {
	manifest := validManifestV2()
	manifest.AllowedPaths = []string{"cmdr-product-spec/06-command/**"}
	if err := validateManifestPathOwnership(manifest, validArchitectureRegistry()); err == nil {
		t.Fatal("expected non-mutable Product Spec rejection")
	}
}

func TestManifestPathOwnershipRequiresProductSpecForbidden(t *testing.T) {
	manifest := validManifestV2()
	manifest.AllowedPaths = []string{"engineering/work/**"}
	manifest.ForbiddenPaths = []string{"tmp/**"}
	if err := validateManifestPathOwnership(manifest, validArchitectureRegistry()); err == nil {
		t.Fatal("expected Product Spec forbidden-path requirement")
	}
}
