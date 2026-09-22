package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePackageJSONExcludesDevAndLocalDependencies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "package.json")
	content := `{
  "dependencies": {"external": "^1.2.3", "internal": "workspace:*"},
  "optionalDependencies": {"optional": "2.0.0"},
  "peerDependencies": {"peer": "3.0.0"},
  "devDependencies": {"dev-only": "4.0.0"}
}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	deps, err := parsePackageJSON(path, "runtime", "src/package.json")
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
	path := filepath.Join(t.TempDir(), "go.mod")
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
	deps, err := parseGoMod(path, "runtime", "src/go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 1 || deps[0].Name != "example.com/external" {
		t.Fatalf("unexpected Go runtime dependencies: %#v", deps)
	}
}

func TestParseVCPKGDependencies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "vcpkg.json")
	content := `{"dependencies":["openssl",{"name":"zlib","version>=":"1.3.1"}]}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	deps, err := parseVCPKG(path, "runtime", "src/vcpkg.json")
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
