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

func TestProductModuleBoundaryEdgeAudit(t *testing.T) {
	seed := func(t *testing.T, allow bool, modAGoMod string) string {
		t.Helper()
		root := t.TempDir()
		writePackage(t, root, "cmdr-product-spec/README.md", "# Spec\n")
		writePackage(t, root, "product-runtime/mod-a/go.mod", modAGoMod)
		writePackage(t, root, "product-runtime/mod-b/go.mod", "module example/b\n\ngo 1.26.11\n")
		writePackage(t, root, "tests/e2e/go.mod", "module example/e2e\n\ngo 1.26.11\n")
		mayDependOn := []string{}
		if allow {
			mayDependOn = []string{"mod-b-runtime"}
		}
		arch := ArchitectureRegistry{
			SchemaVersion: 1,
			DefaultPolicy: "deny-unregistered-path",
			Boundaries: []ArchitectureBoundary{
				{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
				{ID: "mod-a-runtime", Kind: "product-module", Roots: []string{"product-runtime/mod-a/**"}, MutableByImplementation: true, MayDependOn: mayDependOn},
				{ID: "mod-b-runtime", Kind: "product-runtime", Roots: []string{"product-runtime/mod-b/**"}, MutableByImplementation: true},
				{ID: "engineering-control-plane", Kind: "engineering", Roots: []string{"tests/**"}, MutableByImplementation: true, MayDependOn: []string{"product-spec"}},
			},
		}
		writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), arch)
		return root
	}

	t.Run("authorized product-module edge passes", func(t *testing.T) {
		root := seed(t, true, "module example/a\n\ngo 1.26.11\n\nrequire example/b v0.0.0\nreplace example/b => ../mod-b\n")
		summary, err := runBoundaryEdgeAudit(root)
		if err != nil {
			t.Fatalf("expected authorized product-module boundary edge to pass, got: %v", err)
		}
		if summary.RuntimeBoundaries != 2 || summary.ObservedEdges != 1 || summary.CrossBoundaryEdges != 1 {
			t.Fatalf("unexpected summary: %+v", summary)
		}
	})

	t.Run("unauthorized product-module sibling edge fails", func(t *testing.T) {
		root := seed(t, false, "module example/a\n\ngo 1.26.11\n\nrequire example/b v0.0.0\nreplace example/b => ../mod-b\n")
		if _, err := runBoundaryEdgeAudit(root); err == nil {
			t.Fatal("expected forbidden product-module cross-boundary edge to fail")
		}
	})

	t.Run("product-module replace to external path fails", func(t *testing.T) {
		root := seed(t, true, "module example/a\n\ngo 1.26.11\n\nrequire github.com/evil/thirdparty v0.0.0\nreplace github.com/evil/thirdparty => ../../../evil_thirdparty\n")
		if _, err := runBoundaryEdgeAudit(root); err == nil {
			t.Fatal("expected external relative replace in product-module to fail")
		}
	})

	t.Run("product-module replace to non-runtime tests/e2e fails", func(t *testing.T) {
		root := seed(t, true, "module example/a\n\ngo 1.26.11\n\nrequire example/e2e v0.0.0\nreplace example/e2e => ../../tests/e2e\n")
		if _, err := runBoundaryEdgeAudit(root); err == nil {
			t.Fatal("expected non-runtime replace in product-module to fail")
		}
	})
}

func TestRunBoundaryEdgeAudit_RejectsGoWorkspace(t *testing.T) {
	seed := func(t *testing.T) string {
		t.Helper()
		root := t.TempDir()
		writePackage(t, root, "cmdr-product-spec/README.md", "# Spec\n")
		writePackage(t, root, "product-runtime/mod-a/go.mod", "module example/a\n\ngo 1.26.11\n")
		writePackage(t, root, "product-runtime/mod-b/go.mod", "module example/b\n\ngo 1.26.11\n")
		arch := ArchitectureRegistry{
			SchemaVersion: 1,
			DefaultPolicy: "deny-unregistered-path",
			Boundaries: []ArchitectureBoundary{
				{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
				{ID: "mod-a-runtime", Kind: "product-runtime", Roots: []string{"product-runtime/mod-a/**"}, MutableByImplementation: true},
				{ID: "mod-b-runtime", Kind: "product-module", Roots: []string{"product-runtime/mod-b/**"}, MutableByImplementation: true},
			},
		}
		writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), arch)
		return root
	}

	for _, rel := range []string{
		"go.work",
		"go.work.sum",
		"product-runtime/mod-a/go.work",
		"product-runtime/mod-b/go.work",
		"product-runtime/mod-b/nested/go.work.sum",
	} {
		t.Run(rel, func(t *testing.T) {
			root := seed(t)
			writePackage(t, root, rel, "go 1.26.11\n\nuse .\n")
			if _, err := runBoundaryEdgeAudit(root); err == nil {
				t.Fatalf("expected runBoundaryEdgeAudit to fail when %s is present", rel)
			}
		})
	}
}
