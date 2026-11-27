package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePackageJSONExcludesDevAndLocalDependencies(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "package.json")
	content := `{
  "dependencies": {"external": "^1.2.3", "internal": "workspace:*"},
  "optionalDependencies": {"optional": "2.0.0"},
  "peerDependencies": {"peer": "3.0.0"},
  "devDependencies": {"dev-only": "4.0.0"}
}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	deps, err := parsePackageJSON(root, path, "runtime", "src/package.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 runtime dependencies, got %#v", deps)
	}
	for _, dep := range deps {
		if dep.Name == "dev-only" || dep.Name == "internal" {
			t.Fatalf("non-runtime/local dependency leaked into audit: %#v", dep)
		}
	}
}

func TestParseGoModExcludesLocalReplacement(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "go.mod")
	content := `module example.local/app

go 1.25

require (
  example.com/external v1.2.3
  example.local/internal v0.0.0
)

replace example.local/internal => ../internal
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	deps, err := parseGoMod(root, path, "runtime", "src/go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 1 || deps[0].Name != "example.com/external" {
		t.Fatalf("unexpected Go runtime dependencies: %#v", deps)
	}
}

func TestParseVCPKGDependencies(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "vcpkg.json")
	content := `{"dependencies":["openssl",{"name":"zlib","version>=":"1.3.1"}]}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	deps, err := parseVCPKG(root, path, "runtime", "src/vcpkg.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %#v", deps)
	}
}

func TestTrustedBaseIsExactVersionAllowlist(t *testing.T) {
	dep := RuntimeDependency{Ecosystem: "npm", Name: "crypto-safe", Version: "1.2.3"}
	trusted := TrustedBase{RuntimeExceptions: []TrustedRuntimeException{{
		Ecosystem: "npm", Name: "crypto-safe", AllowedVersions: []string{"1.2.3"},
	}}}
	if !trustedDependencyAllowed(dep, trusted.RuntimeExceptions) {
		t.Fatal("expected exact dependency approval")
	}
	dep.Version = "1.2.4"
	if trustedDependencyAllowed(dep, trusted.RuntimeExceptions) {
		t.Fatal("unexpected approval for unlisted version")
	}
}

func TestUnsupportedManifestKindFailsClosed(t *testing.T) {
	if got := unsupportedManifestKind("pyproject.toml"); got != "python" {
		t.Fatalf("unexpected unsupported kind %q", got)
	}
	if got := unsupportedManifestKind("random.txt"); got != "" {
		t.Fatalf("ordinary file incorrectly treated as package manifest: %q", got)
	}
}

func TestRunDependencyAudit_RejectsReplacePathSmuggling(t *testing.T) {
	seedRepo := func(t *testing.T) string {
		t.Helper()
		root := t.TempDir()
		writeTextFixture(t, filepath.Join(root, "cmdr-product-spec", "README.md"), "# Spec\n")
		writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "go.mod"), "module example.com/mod-a\n\ngo 1.26.11\n")
		writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-b", "go.mod"), "module example.com/mod-b\n\ngo 1.26.11\n")
		writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-c", "go.mod"), "module example.com/mod-c\n\ngo 1.26.11\n")
		writeTextFixture(t, filepath.Join(root, "tests", "e2e", "go.mod"), "module example.com/e2e\n\ngo 1.26.11\n")

		arch := ArchitectureRegistry{
			SchemaVersion: 1,
			DefaultPolicy: "deny-unregistered-path",
			Boundaries: []ArchitectureBoundary{
				{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
				{ID: "mod-a-runtime", Kind: "product-runtime", Roots: []string{"product-runtime/mod-a/**"}, MutableByImplementation: true, MayDependOn: []string{"mod-c-runtime"}},
				{ID: "mod-b-runtime", Kind: "product-module", Roots: []string{"product-runtime/mod-b/**"}, MutableByImplementation: true, MayDependOn: []string{"mod-c-runtime"}},
				{ID: "mod-c-runtime", Kind: "product-runtime", Roots: []string{"product-runtime/mod-c/**"}, MutableByImplementation: true},
				{ID: "engineering-control-plane", Kind: "engineering", Roots: []string{"tests/**"}, MutableByImplementation: true, MayDependOn: []string{"product-spec"}},
			},
		}
		writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), arch)
		writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(trustedBasePath)), TrustedBase{
			SchemaVersion: 1,
			DefaultPolicy: "deny-third-party-product-functionality",
		})
		return root
	}

	t.Run("authorized sibling replace passes", func(t *testing.T) {
		root := seedRepo(t)
		writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "go.mod"),
			"module example.com/mod-a\n\ngo 1.26.11\n\nrequire example.com/mod-c v0.0.0\n\nreplace example.com/mod-c => ../mod-c\n")
		writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-b", "go.mod"),
			"module example.com/mod-b\n\ngo 1.26.11\n\nrequire example.com/mod-c v0.0.0\n\nreplace (\n\texample.com/mod-c => ../mod-c // authorized product-module dep\n)\n")
		if _, err := runDependencyAudit(root); err != nil {
			t.Fatalf("expected authorized sibling replace to pass, got: %v", err)
		}
	})

	cases := []struct {
		name      string
		targetMod string
		goModBody string
	}{
		{
			name:      "7c external relative replace escaping root",
			targetMod: "mod-a",
			goModBody: "module example.com/mod-a\n\ngo 1.26.11\n\nrequire github.com/evil/thirdparty v0.0.0\nreplace github.com/evil/thirdparty => ../../../evil_thirdparty\n",
		},
		{
			name:      "7c2 in-repo non-runtime tests/e2e replace",
			targetMod: "mod-a",
			goModBody: "module example.com/mod-a\n\ngo 1.26.11\n\nrequire github.com/evil/smuggled v0.0.0\nreplace github.com/evil/smuggled => ../../tests/e2e\n",
		},
		{
			name:      "7e3 unauthorized sibling boundary replace",
			targetMod: "mod-a",
			goModBody: "module example.com/mod-a\n\ngo 1.26.11\n\nrequire example.com/mod-b v0.0.0\nreplace example.com/mod-b => ../mod-b\n",
		},
		{
			name:      "7f2 product-module external relative replace",
			targetMod: "mod-b",
			goModBody: "module example.com/mod-b\n\ngo 1.26.11\n\nrequire github.com/evil/thirdparty v0.0.0\nreplace github.com/evil/thirdparty => ../../../evil_thirdparty\n",
		},
		{
			name:      "no-whitespace arrow replace smuggling",
			targetMod: "mod-b",
			goModBody: "module example.com/mod-b\n\ngo 1.26.11\n\nrequire github.com/evil/thirdparty v0.0.0\nreplace github.com/evil/thirdparty=>../../tests/e2e\n",
		},
		{
			name:      "replace block module-identity spoofing",
			targetMod: "mod-a",
			goModBody: "module example.com/mod-a\n\ngo 1.26.11\n\nrequire github.com/evil/alias v0.0.0\nreplace (\n\tgithub.com/evil/alias => ../mod-c // spoofed\n)\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := seedRepo(t)
			writeTextFixture(t, filepath.Join(root, "product-runtime", tc.targetMod, "go.mod"), tc.goModBody)
			if _, err := runDependencyAudit(root); err == nil {
				t.Fatalf("expected runDependencyAudit to fail for %s, got nil", tc.name)
			}
		})
	}
}
