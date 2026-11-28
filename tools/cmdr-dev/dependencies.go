package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
		if boundary.Kind != "product-runtime" && boundary.Kind != "product-module" {
			continue
		}
		runtimeBoundaries++
		for _, rootPattern := range boundary.Roots {
			scanned, bad, err := scanRuntimeBoundary(root, boundary.ID, rootPattern)
			if err != nil {
				return DependencyAuditSummary{}, err
			}
			if err := validateRuntimeBoundaryGoModReplacements(root, registry, boundary, rootPattern); err != nil {
				return DependencyAuditSummary{}, err
			}
			deps = append(deps, scanned...)
			unsupported = append(unsupported, bad...)
			for _, dep := range scanned {
				supportedManifests[dep.Manifest] = true
			}
		}
	}
	workspaceBad, err := scanRepositoryGoWorkspaceManifests(root, registry)
	if err != nil {
		return DependencyAuditSummary{}, err
	}
	unsupported = append(unsupported, workspaceBad...)
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
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected below.
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("runtime boundary contains symlink: %s", path)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if isGoWorkspaceFileName(entry.Name()) {
			unsupported = append(unsupported, UnsupportedDependencyManifest{Boundary: boundaryID, Path: rel, Kind: "go-workspace"})
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if path != scanRoot && ignoredDependencyDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
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

type goModReplacement struct {
	OldModule   string
	OldVersion  string
	Target      string
	NewVersion  string
	IsLocalPath bool
}

func unquoteGoModToken(tok string) string {
	tok = strings.TrimSpace(tok)
	if len(tok) >= 2 && ((tok[0] == '"' && tok[len(tok)-1] == '"') || (tok[0] == '`' && tok[len(tok)-1] == '`')) {
		if unquoted, err := strconv.Unquote(tok); err == nil {
			return unquoted
		}
	}
	return tok
}

func isAbsoluteOrDrivePath(target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	if filepath.IsAbs(target) || strings.HasPrefix(target, "/") || strings.HasPrefix(target, "\\") {
		return true
	}
	if len(target) >= 2 && ((target[0] >= 'a' && target[0] <= 'z') || (target[0] >= 'A' && target[0] <= 'Z')) && target[1] == ':' {
		return true
	}
	return false
}

func isLocalGoReplaceTarget(target string, rhsTokens []string) bool {
	if len(rhsTokens) == 1 {
		return true
	}
	if isAbsoluteOrDrivePath(target) {
		return true
	}
	for _, prefix := range []string{"./", "../", ".\\", "..\\"} {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return target == "." || target == ".."
}

func resolveLocalReplaceRepoRel(root, manifestRel, target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", fmt.Errorf("empty local replacement path")
	}
	if isAbsoluteOrDrivePath(target) {
		return "", fmt.Errorf("local replacement path %q must be relative", target)
	}
	normalizedTarget := strings.ReplaceAll(target, "\\", "/")
	manifestDir := filepath.Dir(filepath.FromSlash(manifestRel))
	joined := filepath.Clean(filepath.Join(manifestDir, filepath.FromSlash(normalizedTarget)))
	if filepath.IsAbs(joined) || joined == ".." || strings.HasPrefix(joined, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("local replacement path %q escapes repository root", target)
	}
	rel := filepath.ToSlash(joined)
	if _, err := resolveRepoPath(root, rel, true); err != nil {
		return "", err
	}
	return rel, nil
}

func boundaryForRuntimeRepoPath(registry ArchitectureRegistry, path string) (ArchitectureBoundary, bool) {
	path = filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	var found *ArchitectureBoundary
	for i := range registry.Boundaries {
		boundary := &registry.Boundaries[i]
		if boundary.Kind != "product-runtime" && boundary.Kind != "product-module" {
			continue
		}
		for _, pattern := range boundary.Roots {
			if repoPathOwnedByPattern(pattern, path) {
				if found != nil && found.ID != boundary.ID {
					return ArchitectureBoundary{}, false
				}
				copy := *boundary
				found = &copy
			}
		}
	}
	if found == nil {
		return ArchitectureBoundary{}, false
	}
	return *found, true
}

func validateRuntimeBoundaryGoModReplacements(root string, registry ArchitectureRegistry, sourceBoundary ArchitectureBoundary, rootPattern string) error {
	scanRoot, err := runtimeScanRoot(root, rootPattern)
	if err != nil {
		return err
	}
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected below.
	return filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("runtime boundary contains symlink: %s", path)
		}
		if entry.IsDir() {
			if path != scanRoot && ignoredDependencyDir(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() != "go.mod" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		_, _, replacements, err := parseGoModDetailed(root, path, sourceBoundary.ID, rel)
		if err != nil {
			return err
		}
		for _, rep := range replacements {
			if !rep.IsLocalPath {
				return fmt.Errorf("%s unapproved non-local Go replacement %s => %s %s", rel, rep.OldModule, rep.Target, rep.NewVersion)
			}
			targetRel, err := resolveLocalReplaceRepoRel(root, rel, rep.Target)
			if err != nil {
				return fmt.Errorf("%s invalid Go replacement %s => %s: %w", rel, rep.OldModule, rep.Target, err)
			}
			if _, err := resolveRepoPath(root, targetRel, false); err != nil {
				return fmt.Errorf("%s Go replacement %s => %s failed containment check: %w", rel, rep.OldModule, rep.Target, err)
			}
			info, err := statRepoPath(root, targetRel)
			if err != nil || !info.IsDir() {
				return fmt.Errorf("%s Go replacement %s targets missing or non-directory path %s", rel, rep.OldModule, targetRel)
			}
			targetBoundary, ok := boundaryForRuntimeRepoPath(registry, targetRel)
			if !ok {
				return fmt.Errorf("%s Go replacement %s targets path %s outside product-runtime/product-module boundaries", rel, rep.OldModule, targetRel)
			}
			if targetBoundary.ID != sourceBoundary.ID {
				allowed := false
				for _, dep := range sourceBoundary.MayDependOn {
					if dep == targetBoundary.ID {
						allowed = true
						break
					}
				}
				if !allowed {
					return fmt.Errorf("%s Go replacement %s targets unauthorized boundary %s (not in %s may_depend_on)", rel, rep.OldModule, targetBoundary.ID, sourceBoundary.ID)
				}
			}
			targetGoModRel := targetRel + "/go.mod"
			targetGoModInfo, err := lstatRepoEntry(root, targetGoModRel)
			if err != nil || targetGoModInfo.IsDir() || targetGoModInfo.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("%s Go replacement %s targets %s without a regular go.mod", rel, rep.OldModule, targetRel)
			}
			targetGoModAbs := filepath.Join(root, filepath.FromSlash(targetGoModRel))
			targetModule, _, _, err := parseGoModDetailed(root, targetGoModAbs, targetBoundary.ID, targetGoModRel)
			if err != nil {
				return fmt.Errorf("%s Go replacement %s failed reading target %s: %w", rel, rep.OldModule, targetGoModRel, err)
			}
			if targetModule == "" || targetModule != rep.OldModule {
				return fmt.Errorf("%s Go replacement %s targets %s whose module is %q (expected %q)", rel, rep.OldModule, targetRel, targetModule, rep.OldModule)
			}
		}
		return nil
	})
}

func parseGoModDetailed(root, path, boundary, manifest string) (string, []RuntimeDependency, []goModReplacement, error) {
	data, err := readRepoFile(root, path)
	if err != nil {
		return "", nil, nil, err
	}
	moduleName := ""
	var replacements []goModReplacement
	lines := strings.Split(string(data), "\n")
	inReplace := false
	for _, raw := range lines {
		line := strings.TrimSpace(strings.SplitN(raw, "//", 2)[0])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			fields := strings.Fields(strings.TrimPrefix(line, "module "))
			if len(fields) > 0 {
				moduleName = unquoteGoModToken(fields[0])
			}
			continue
		}
		if strings.HasPrefix(line, "replace") && strings.TrimSpace(strings.TrimPrefix(line, "replace")) == "(" {
			inReplace = true
			continue
		}
		if inReplace && line == ")" {
			inReplace = false
			continue
		}
		isSingleReplace := false
		if strings.HasPrefix(line, "replace ") || strings.HasPrefix(line, "replace\t") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "replace"))
			isSingleReplace = true
		} else if !inReplace {
			continue
		}
		if !isSingleReplace && !inReplace {
			continue
		}
		parts := strings.SplitN(line, "=>", 2)
		if len(parts) != 2 {
			continue
		}
		left := strings.Fields(parts[0])
		right := strings.Fields(parts[1])
		if len(left) == 0 || len(right) == 0 {
			continue
		}
		oldMod := unquoteGoModToken(left[0])
		oldVer := ""
		if len(left) >= 2 {
			oldVer = unquoteGoModToken(left[1])
		}
		target := unquoteGoModToken(right[0])
		newVer := ""
		if len(right) >= 2 {
			newVer = unquoteGoModToken(right[1])
		}
		replacements = append(replacements, goModReplacement{
			OldModule:   oldMod,
			OldVersion:  oldVer,
			Target:      target,
			NewVersion:  newVer,
			IsLocalPath: isLocalGoReplaceTarget(target, right),
		})
	}

	var requires []RuntimeDependency
	inRequire := false
	for _, raw := range lines {
		line := strings.TrimSpace(strings.SplitN(raw, "//", 2)[0])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "require") && strings.TrimSpace(strings.TrimPrefix(line, "require")) == "(" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") || strings.HasPrefix(line, "require\t") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "require"))
		} else if !inRequire {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		requires = append(requires, RuntimeDependency{
			Boundary:  boundary,
			Ecosystem: "go",
			Name:      unquoteGoModToken(fields[0]),
			Version:   unquoteGoModToken(fields[1]),
			Manifest:  manifest,
		})
	}
	return moduleName, requires, replacements, nil
}

func parseGoMod(root, path, boundary, manifest string) ([]RuntimeDependency, error) {
	_, requires, replacements, err := parseGoModDetailed(root, path, boundary, manifest)
	if err != nil {
		return nil, err
	}
	localReplacements := map[string]bool{}
	for _, rep := range replacements {
		if !rep.IsLocalPath {
			continue
		}
		if _, err := resolveLocalReplaceRepoRel(root, manifest, rep.Target); err == nil {
			localReplacements[rep.OldModule] = true
		}
	}
	var out []RuntimeDependency
	for _, req := range requires {
		if localReplacements[req.Name] {
			continue
		}
		out = append(out, req)
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

func isGoWorkspaceFileName(name string) bool {
	return name == "go.work" || name == "go.work.sum"
}

func boundaryForAnyRepoPath(registry ArchitectureRegistry, path string) string {
	if boundary, ok := boundaryForRuntimeRepoPath(registry, path); ok {
		return boundary.ID
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	for _, boundary := range registry.Boundaries {
		for _, pattern := range boundary.Roots {
			if repoPathOwnedByPattern(pattern, cleaned) {
				return boundary.ID
			}
		}
	}
	return "repository"
}

func scanRepositoryGoWorkspaceManifests(root string, registry ArchitectureRegistry) ([]UnsupportedDependencyManifest, error) {
	rootAbs, err := resolveRepoPath(root, ".", false)
	if err != nil {
		return nil, err
	}
	var found []UnsupportedDependencyManifest
	// #nosec G703 -- rootAbs is repository-confined by resolveRepoPath; symlink go.work entries are rejected below.
	err = filepath.WalkDir(rootAbs, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(rootAbs, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if isGoWorkspaceFileName(entry.Name()) {
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("repository contains forbidden Go workspace symlink: %s", rel)
			}
			found = append(found, UnsupportedDependencyManifest{
				Boundary: boundaryForAnyRepoPath(registry, rel),
				Path:     rel,
				Kind:     "go-workspace",
			})
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if path != rootAbs && (ignoredDependencyDir(entry.Name()) || entry.Name() == "testdata" || entry.Name() == ".cmdr-tmp") {
				return filepath.SkipDir
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dedupeUnsupportedManifests(found), nil
}

func unsupportedManifestKind(name string) string {
	switch {
	case isGoWorkspaceFileName(name):
		return "go-workspace"
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
