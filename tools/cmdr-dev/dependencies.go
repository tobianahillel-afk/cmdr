package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const trustedBasePath = "engineering/security/trusted-base.json"

type TrustedRuntimeException struct {
	ID                  string   `json:"id"`
	Ecosystem           string   `json:"ecosystem"`
	Name                string   `json:"name"`
	AllowedVersions     []string `json:"allowed_versions"`
	Purpose             string   `json:"purpose"`
	Provenance          string   `json:"provenance"`
	License             string   `json:"license"`
	SecurityReview      string   `json:"security_review"`
	Telemetry           string   `json:"telemetry"`
	ReplacementBoundary string   `json:"replacement_boundary"`
}

type TrustedBase struct {
	SchemaVersion     int                       `json:"schema_version"`
	DefaultPolicy     string                    `json:"default_policy"`
	RuntimeExceptions []TrustedRuntimeException `json:"runtime_exceptions"`
	Notes             []string                  `json:"notes"`
}

type RuntimeDependency struct {
	Boundary  string `json:"boundary"`
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Manifest  string `json:"manifest"`
}

type UnsupportedDependencyManifest struct {
	Boundary string `json:"boundary"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
}

type DependencyAuditSummary struct {
	RuntimeBoundaries    int                             `json:"runtime_boundaries"`
	SupportedManifests   int                             `json:"supported_manifests"`
	RuntimeDependencies  int                             `json:"runtime_dependencies"`
	Approved             int                             `json:"approved"`
	Unapproved           int                             `json:"unapproved"`
	UnsupportedManifests int                             `json:"unsupported_manifests"`
	UnapprovedItems      []RuntimeDependency             `json:"unapproved_items,omitempty"`
	UnsupportedItems     []UnsupportedDependencyManifest `json:"unsupported_items,omitempty"`
}

func runDependencyAudit(root string) (DependencyAuditSummary, error) {
	var registry ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &registry); err != nil {
		return DependencyAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return DependencyAuditSummary{}, err
	}
	var trusted TrustedBase
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(trustedBasePath)), &trusted); err != nil {
		return DependencyAuditSummary{}, err
	}
	if err := validateTrustedBase(trusted); err != nil {
		return DependencyAuditSummary{}, err
	}

	var deps []RuntimeDependency
	var unsupported []UnsupportedDependencyManifest
	supportedManifests := map[string]bool{}
	runtimeBoundaries := 0
	for _, boundary := range registry.Boundaries {
		if boundary.Kind != "product-runtime" {
			continue
		}
		runtimeBoundaries++
		for _, rootPattern := range boundary.Roots {
			scanned, bad, err := scanRuntimeBoundary(root, boundary.ID, rootPattern)
			if err != nil {
				return DependencyAuditSummary{}, err
			}
			deps = append(deps, scanned...)
			unsupported = append(unsupported, bad...)
			for _, dep := range scanned {
				supportedManifests[dep.Manifest] = true
			}
		}
	}
	deps = dedupeRuntimeDependencies(deps)
	unsupported = dedupeUnsupportedManifests(unsupported)

	approved, unapproved := evaluateTrustedBase(deps, trusted)
	summary := DependencyAuditSummary{
		RuntimeBoundaries:    runtimeBoundaries,
		SupportedManifests:   len(supportedManifests),
		RuntimeDependencies:  len(deps),
		Approved:             approved,
		Unapproved:           len(unapproved),
		UnsupportedManifests: len(unsupported),
		UnapprovedItems:      unapproved,
		UnsupportedItems:     unsupported,
	}
	if summary.UnsupportedManifests > 0 {
		return summary, fmt.Errorf("unsupported product-runtime dependency manifests detected: %v", unsupported)
	}
	if summary.Unapproved > 0 {
		return summary, fmt.Errorf("unapproved third-party product-runtime dependencies detected: %v", unapproved)
	}
	return summary, nil
}

func validateTrustedBase(trusted TrustedBase) error {
	if trusted.SchemaVersion != 1 {
		return fmt.Errorf("unsupported Trusted Base schema_version %d", trusted.SchemaVersion)
	}
	if trusted.DefaultPolicy != "deny-third-party-product-functionality" {
		return fmt.Errorf("Trusted Base default policy must deny third-party product functionality")
	}
	seen := map[string]bool{}
	for _, exception := range trusted.RuntimeExceptions {
		if exception.ID == "" || exception.Ecosystem == "" || exception.Name == "" || exception.Purpose == "" ||
			exception.Provenance == "" || exception.License == "" || exception.SecurityReview == "" ||
			exception.Telemetry == "" || exception.ReplacementBoundary == "" {
			return fmt.Errorf("Trusted Base exception %q is incomplete", exception.ID)
		}
		if seen[exception.ID] {
			return fmt.Errorf("duplicate Trusted Base exception %s", exception.ID)
		}
		seen[exception.ID] = true
		if len(exception.AllowedVersions) == 0 {
			return fmt.Errorf("Trusted Base exception %s has no allowed versions", exception.ID)
		}
		for _, version := range exception.AllowedVersions {
			if strings.TrimSpace(version) == "" || version == "*" {
				return fmt.Errorf("Trusted Base exception %s has unsafe version allowance %q", exception.ID, version)
			}
		}
	}
	return nil
}

func scanRuntimeBoundary(root, boundaryID, rootPattern string) ([]RuntimeDependency, []UnsupportedDependencyManifest, error) {
	scanRoot, err := runtimeScanRoot(root, rootPattern)
	if err != nil {
		return nil, nil, err
	}
	var deps []RuntimeDependency
	var unsupported []UnsupportedDependencyManifest
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != scanRoot && ignoredDependencyDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		switch entry.Name() {
		case "go.mod":
			found, err := parseGoMod(root, path, boundaryID, rel)
			if err != nil {
				return err
			}
			deps = append(deps, found...)
		case "package.json":
			found, err := parsePackageJSON(root, path, boundaryID, rel)
			if err != nil {
				return err
			}
			deps = append(deps, found...)
		case "vcpkg.json":
			found, err := parseVCPKG(root, path, boundaryID, rel)
			if err != nil {
				return err
			}
			deps = append(deps, found...)
		default:
			if kind := unsupportedManifestKind(entry.Name()); kind != "" {
				unsupported = append(unsupported, UnsupportedDependencyManifest{Boundary: boundaryID, Path: rel, Kind: kind})
			}
		}
		return nil
	})
	return deps, unsupported, err
}

func runtimeScanRoot(root, pattern string) (string, error) {
	pattern, err := normalizeArchitecturePattern(pattern)
	if err != nil {
		return "", err
	}
	prefix := strings.TrimSuffix(pattern, "**")
	prefix = strings.TrimSuffix(prefix, "/")
	if prefix == "" {
		return "", fmt.Errorf("product-runtime boundary root cannot be repository root")
	}
	path, err := resolveRepoPath(root, prefix, false)
	if err != nil {
		return "", fmt.Errorf("product-runtime boundary root %s is invalid: %w", pattern, err)
	}
	info, err := statRepoPath(root, path)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("product-runtime boundary root %s is not an existing directory", pattern)
	}
	return path, nil
}

func parseGoMod(root, path, boundary, manifest string) ([]RuntimeDependency, error) {
	data, err := readRepoFile(root, path)
	if err != nil {
		return nil, err
	}
	localReplacements := map[string]bool{}
	lines := strings.Split(string(data), "\n")
	inReplace := false
	for _, line := range lines {
		line = strings.TrimSpace(strings.SplitN(line, "//", 2)[0])
		if line == "" {
			continue
		}
		if line == "replace (" {
			inReplace = true
			continue
		}
		if inReplace && line == ")" {
			inReplace = false
			continue
		}
		if strings.HasPrefix(line, "replace ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		}
		if inReplace || strings.Contains(line, " => ") {
			parts := strings.Split(line, " => ")
			if len(parts) == 2 {
				left := strings.Fields(parts[0])
				right := strings.Fields(parts[1])
				if len(left) > 0 && len(right) > 0 && (strings.HasPrefix(right[0], "./") || strings.HasPrefix(right[0], "../")) {
					localReplacements[left[0]] = true
				}
			}
		}
	}

	var out []RuntimeDependency
	inRequire := false
	for _, line := range lines {
		line = strings.TrimSpace(strings.SplitN(line, "//", 2)[0])
		if line == "" {
			continue
		}
		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "require "))
		} else if !inRequire {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 || localReplacements[fields[0]] {
			continue
		}
		out = append(out, RuntimeDependency{Boundary: boundary, Ecosystem: "go", Name: fields[0], Version: fields[1], Manifest: manifest})
	}
	return out, nil
}

func parsePackageJSON(root, path, boundary, manifest string) ([]RuntimeDependency, error) {
	var pkg struct {
		Dependencies         map[string]string `json:"dependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
		PeerDependencies     map[string]string `json:"peerDependencies"`
	}
	data, err := readRepoFile(root, path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("decode %s: %w", manifest, err)
	}
	merged := map[string]string{}
	for _, values := range []map[string]string{pkg.Dependencies, pkg.OptionalDependencies, pkg.PeerDependencies} {
		for name, version := range values {
			merged[name] = version
		}
	}
	var out []RuntimeDependency
	for name, version := range merged {
		if isLocalNodeDependency(version) {
			continue
		}
		out = append(out, RuntimeDependency{Boundary: boundary, Ecosystem: "npm", Name: name, Version: version, Manifest: manifest})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func parseVCPKG(root, path, boundary, manifest string) ([]RuntimeDependency, error) {
	var doc struct {
		Dependencies []json.RawMessage `json:"dependencies"`
	}
	data, err := readRepoFile(root, path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("decode %s: %w", manifest, err)
	}
	var out []RuntimeDependency
	for _, raw := range doc.Dependencies {
		var name string
		if err := json.Unmarshal(raw, &name); err == nil {
			out = append(out, RuntimeDependency{Boundary: boundary, Ecosystem: "vcpkg", Name: name, Version: "unversioned", Manifest: manifest})
			continue
		}
		var item struct {
			Name    string `json:"name"`
			Version string `json:"version>="`
		}
		if err := json.Unmarshal(raw, &item); err != nil || item.Name == "" {
			return nil, fmt.Errorf("unsupported vcpkg dependency entry in %s", manifest)
		}
		version := item.Version
		if version == "" {
			version = "unversioned"
		}
		out = append(out, RuntimeDependency{Boundary: boundary, Ecosystem: "vcpkg", Name: item.Name, Version: version, Manifest: manifest})
	}
	return out, nil
}

func evaluateTrustedBase(deps []RuntimeDependency, trusted TrustedBase) (int, []RuntimeDependency) {
	approved := 0
	var unapproved []RuntimeDependency
	for _, dep := range deps {
		if trustedDependencyAllowed(dep, trusted.RuntimeExceptions) {
			approved++
		} else {
			unapproved = append(unapproved, dep)
		}
	}
	return approved, unapproved
}

func trustedDependencyAllowed(dep RuntimeDependency, exceptions []TrustedRuntimeException) bool {
	for _, exception := range exceptions {
		if exception.Ecosystem != dep.Ecosystem || exception.Name != dep.Name {
			continue
		}
		for _, version := range exception.AllowedVersions {
			if version == dep.Version {
				return true
			}
		}
	}
	return false
}

func dedupeRuntimeDependencies(values []RuntimeDependency) []RuntimeDependency {
	sort.Slice(values, func(i, j int) bool {
		a := values[i]
		b := values[j]
		if a.Boundary != b.Boundary {
			return a.Boundary < b.Boundary
		}
		if a.Ecosystem != b.Ecosystem {
			return a.Ecosystem < b.Ecosystem
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		if a.Version != b.Version {
			return a.Version < b.Version
		}
		return a.Manifest < b.Manifest
	})
	if len(values) == 0 {
		return nil
	}
	out := values[:0]
	var previous RuntimeDependency
	for i, value := range values {
		if i == 0 || value != previous {
			out = append(out, value)
			previous = value
		}
	}
	return out
}

func dedupeUnsupportedManifests(values []UnsupportedDependencyManifest) []UnsupportedDependencyManifest {
	sort.Slice(values, func(i, j int) bool {
		if values[i].Boundary != values[j].Boundary {
			return values[i].Boundary < values[j].Boundary
		}
		return values[i].Path < values[j].Path
	})
	if len(values) == 0 {
		return nil
	}
	out := values[:0]
	var previous UnsupportedDependencyManifest
	for i, value := range values {
		if i == 0 || value != previous {
			out = append(out, value)
			previous = value
		}
	}
	return out
}

func isLocalNodeDependency(version string) bool {
	version = strings.TrimSpace(version)
	for _, prefix := range []string{"workspace:", "file:", "link:", "./", "../"} {
		if strings.HasPrefix(version, prefix) {
			return true
		}
	}
	return false
}

func ignoredDependencyDir(name string) bool {
	switch name {
	case ".git", "node_modules", "vendor", "target", ".venv", "venv", "__pycache__":
		return true
	default:
		return false
	}
}

func unsupportedManifestKind(name string) string {
	switch {
	case name == "Cargo.toml":
		return "cargo"
	case name == "pyproject.toml", name == "requirements.txt", strings.HasPrefix(name, "requirements-"):
		return "python"
	case name == "pom.xml", name == "build.gradle", name == "build.gradle.kts":
		return "jvm"
	case name == "conanfile.txt", name == "conanfile.py":
		return "conan"
	case name == "composer.json":
		return "composer"
	case name == "Gemfile":
		return "ruby"
	case name == "Package.swift":
		return "swift"
	case name == "pubspec.yaml":
		return "dart"
	case name == "mix.exs":
		return "elixir"
	default:
		return ""
	}
}


