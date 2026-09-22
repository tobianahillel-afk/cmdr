package main

import (
	"os"
	"path/filepath"
	"testing"
)

func runtimeRegistry(allow bool) ArchitectureRegistry {
	deps := []string{}
	if allow {
		deps = []string{"runtime-b"}
	}
	return ArchitectureRegistry{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unregistered-path",
		Boundaries: []ArchitectureBoundary{
			{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
			{ID: "runtime-a", Kind: "product-runtime", Roots: []string{"apps/a/**"}, MutableByImplementation: true, MayDependOn: deps},
			{ID: "runtime-b", Kind: "product-runtime", Roots: []string{"libs/b/**"}, MutableByImplementation: true},
		},
	}
}

func writePackage(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceDependencyCreatesAllowedBoundaryEdge(t *testing.T) {
	root := t.TempDir()
	writePackage(t, root, "apps/a/package.json", `{"name":"a","dependencies":{"b":"workspace:*"}}`)
	writePackage(t, root, "libs/b/package.json", `{"name":"b"}`)
	registry := runtimeRegistry(true)
	packages, _, err := discoverLocalPackages(root, registry)
	if err != nil {
		t.Fatal(err)
	}
	edges, err := discoverBoundaryEdges(root, registry, packages)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 || edges[0].From != "runtime-a" || edges[0].To != "runtime-b" {
		t.Fatalf("unexpected edges: %#v", edges)
	}
	if err := validateBoundaryEdges(registry, edges); err != nil {
		t.Fatal(err)
	}
}

func TestForbiddenBoundaryEdgeFails(t *testing.T) {
	root := t.TempDir()
	writePackage(t, root, "apps/a/package.json", `{"name":"a","dependencies":{"b":"workspace:*"}}`)
	writePackage(t, root, "libs/b/package.json", `{"name":"b"}`)
	registry := runtimeRegistry(false)
	packages, _, err := discoverLocalPackages(root, registry)
	if err != nil {
		t.Fatal(err)
	}
	edges, err := discoverBoundaryEdges(root, registry, packages)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateBoundaryEdges(registry, edges); err == nil {
		t.Fatal("expected forbidden cross-boundary edge")
	}
}

func TestLocalPathOutsideRuntimeBoundaryFails(t *testing.T) {
	root := t.TempDir()
	writePackage(t, root, "apps/a/package.json", `{"name":"a","dependencies":{"outside":"file:../../outside"}}`)
	writePackage(t, root, "libs/b/package.json", `{"name":"b"}`)
	registry := runtimeRegistry(true)
	packages, _, err := discoverLocalPackages(root, registry)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := discoverBoundaryEdges(root, registry, packages); err == nil {
		t.Fatal("expected unowned local target rejection")
	}
}

func TestGoLocalReplacementCreatesBoundaryEdge(t *testing.T) {
	root := t.TempDir()
	writePackage(t, root, "apps/a/go.mod", "module example/a\n\ngo 1.25\n\nrequire example/b v0.0.0\nreplace example/b => ../../libs/b\n")
	writePackage(t, root, "libs/b/go.mod", "module example/b\n\ngo 1.25\n")
	registry := runtimeRegistry(true)
	packages, _, err := discoverLocalPackages(root, registry)
	if err != nil {
		t.Fatal(err)
	}
	edges, err := discoverBoundaryEdges(root, registry, packages)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 || edges[0].To != "runtime-b" {
		t.Fatalf("unexpected Go edges: %#v", edges)
	}
}
