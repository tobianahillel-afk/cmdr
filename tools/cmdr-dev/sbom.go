package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	cycloneDXSchemaURL     = "https://cyclonedx.org/schema/bom-1.7.schema.json"
	cycloneDXSpecVersion   = "1.7"
	sbomGeneratorVersion   = "sbom-v1"
	defaultSBOMOutputPath  = "engineering/security/evidence/cmdr.cdx.json"
)

type CDXProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CDXComponent struct {
	Type       string        `json:"type"`
	BOMRef     string        `json:"bom-ref,omitempty"`
	Name       string        `json:"name"`
	Version    string        `json:"version,omitempty"`
	Scope      string        `json:"scope,omitempty"`
	Properties []CDXProperty `json:"properties,omitempty"`
}

type CDXTools struct {
	Components []CDXComponent `json:"components,omitempty"`
}

type CDXMetadata struct {
	Tools     CDXTools     `json:"tools"`
	Component CDXComponent `json:"component"`
}

type CDXComposition struct {
	Aggregate  string   `json:"aggregate"`
	Assemblies []string `json:"assemblies,omitempty"`
}

type CycloneDXBOM struct {
	Schema       string           `json:"$schema"`
	BOMFormat    string           `json:"bomFormat"`
	SpecVersion  string           `json:"specVersion"`
	Version      int              `json:"version"`
	Metadata     CDXMetadata      `json:"metadata"`
	Components   []CDXComponent   `json:"components,omitempty"`
	Compositions []CDXComposition `json:"compositions,omitempty"`
}

type sbomInputComponent struct {
	Kind      string
	Ecosystem string
	Name      string
	Version   string
	Manifest  string
	Scope     string
	Direct    bool
	Commit    string
	License   string
}

type SBOMSummary struct {
	SchemaVersion         string `json:"schema_version"`
	Generator             string `json:"generator"`
	GeneratorVersion      string `json:"generator_version"`
	SourceSHA             string `json:"source_sha"`
	Components            int    `json:"components"`
	RuntimeComponents     int    `json:"runtime_components"`
	DevelopmentComponents int    `json:"development_components"`
	UnsupportedManifests  int    `json:"unsupported_manifests"`
	DigestSHA256          string `json:"digest_sha256"`
	Output                string `json:"output"`
}

var exactNPMVersionPattern = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

func runSBOM(root, output string) (SBOMSummary, error) {
	sourceSHA, err := currentSourceCommit(root)
	if err != nil {
		return SBOMSummary{}, err
	}
	inputs, unsupported, err := collectSBOMInputs(root)
	if err != nil {
		return SBOMSummary{}, err
	}
	if len(unsupported) > 0 {
		return SBOMSummary{UnsupportedManifests: len(unsupported)}, fmt.Errorf("SBOM refuses unsupported dependency manifests: %v", unsupported)
	}
	for _, input := range inputs {
		if err := validateResolvedSBOMInput(input); err != nil {
			return SBOMSummary{}, err
		}
	}
	bom, runtimeCount, developmentCount, err := buildCycloneDXBOM(sourceSHA, inputs)
	if err != nil {
		return SBOMSummary{}, err
	}
	data, err := json.MarshalIndent(bom, "", "  ")
	if err != nil {
		return SBOMSummary{}, err
	}
	data = append(data, '\n')
	digest := sha256.Sum256(data)

	if strings.TrimSpace(output) == "" {
		output = defaultSBOMOutputPath
	}
	outputPath, err := writeRepoFile(root, output, data)
	if err != nil {
		return SBOMSummary{}, err
	}
	return SBOMSummary{
		SchemaVersion:         cycloneDXSpecVersion,
		Generator:             "cmdr-dev",
		GeneratorVersion:      sbomGeneratorVersion,
		SourceSHA:             sourceSHA,
		Components:            len(bom.Components),
		RuntimeComponents:     runtimeCount,
		DevelopmentComponents: developmentCount,
		UnsupportedManifests:  0,
		DigestSHA256:          fmt.Sprintf("%x", digest[:]),
		Output:                outputPath,
	}, nil
}

func currentSourceCommit(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--verify", "HEAD^{commit}") // #nosec G204,G702 -- executable and arguments are fixed; root is the resolved repository root; no shell is used.
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("resolve SBOM source commit: %w", err)
	}
	sha, err := validateFullCommitID(strings.TrimSpace(string(out)))
	if err != nil {
		return "", fmt.Errorf("invalid SBOM source commit: %w", err)
	}
	return sha, nil
}

func collectSBOMInputs(root string) ([]sbomInputComponent, []UnsupportedDependencyManifest, error) {
	var registry ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &registry); err != nil {
		return nil, nil, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return nil, nil, err
	}

	var inputs []sbomInputComponent
	var unsupported []UnsupportedDependencyManifest
	for _, boundary := range registry.Boundaries {
		if boundary.Kind != "product-runtime" {
			continue
		}
		for _, rootPattern := range boundary.Roots {
			deps, bad, err := scanRuntimeBoundary(root, boundary.ID, rootPattern)
			if err != nil {
				return nil, nil, err
			}
			unsupported = append(unsupported, bad...)
			for _, dep := range deps {
				inputs = append(inputs, sbomInputComponent{
					Kind: "library", Ecosystem: dep.Ecosystem, Name: dep.Name, Version: dep.Version,
					Manifest: dep.Manifest, Scope: "runtime", Direct: true,
				})
			}
		}
	}

	devDeps, devUnsupported, err := collectEngineeringDependencies(root)
	if err != nil {
		return nil, nil, err
	}
	inputs = append(inputs, devDeps...)
	unsupported = append(unsupported, devUnsupported...)

	tools, err := loadDevelopmentTools(root)
	if err != nil {
		return nil, nil, err
	}
	for _, tool := range tools.Tools {
		inputs = append(inputs, sbomInputComponent{
			Kind: "application", Ecosystem: "go", Name: tool.Module, Version: tool.Version,
			Manifest: developmentToolRegistryPath, Scope: "development", Direct: true,
			Commit: tool.Commit, License: tool.License,
		})
	}
	return dedupeSBOMInputs(inputs), dedupeUnsupportedManifests(unsupported), nil
}

func collectEngineeringDependencies(root string) ([]sbomInputComponent, []UnsupportedDependencyManifest, error) {
	paths, err := trackedRepositoryPaths(root)
	if err != nil {
		return nil, nil, err
	}
	const prefix = "tools/cmdr-dev/"
	var inputs []sbomInputComponent
	var unsupported []UnsupportedDependencyManifest
	for _, rel := range paths {
		if !strings.HasPrefix(rel, prefix) {
			continue
		}
		name := filepath.Base(filepath.FromSlash(rel))
		var deps []RuntimeDependency
		switch name {
		case "go.mod":
			deps, err = parseGoMod(root, rel, "engineering-control-plane", rel)
		case "package.json":
			deps, err = parsePackageJSON(root, rel, "engineering-control-plane", rel)
		case "vcpkg.json":
			deps, err = parseVCPKG(root, rel, "engineering-control-plane", rel)
		default:
			if kind := unsupportedManifestKind(name); kind != "" {
				unsupported = append(unsupported, UnsupportedDependencyManifest{
					Boundary: "engineering-control-plane", Path: rel, Kind: kind,
				})
			}
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		for _, dep := range deps {
			inputs = append(inputs, sbomInputComponent{
				Kind: "library", Ecosystem: dep.Ecosystem, Name: dep.Name, Version: dep.Version,
				Manifest: dep.Manifest, Scope: "development", Direct: true,
			})
		}
	}
	return inputs, unsupported, nil
}

func validateResolvedSBOMInput(input sbomInputComponent) error {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Manifest) == "" {
		return fmt.Errorf("SBOM component has incomplete identity: %+v", input)
	}
	version := strings.TrimSpace(input.Version)
	if version == "" || version == "*" || strings.EqualFold(version, "latest") || strings.EqualFold(version, "unversioned") {
		return fmt.Errorf("SBOM refuses unresolved %s component %s from %s with version %q", input.Ecosystem, input.Name, input.Manifest, input.Version)
	}
	switch input.Ecosystem {
	case "go":
		if !strings.HasPrefix(version, "v") || strings.ContainsAny(version, "*<>=~^| ,") {
			return fmt.Errorf("SBOM refuses unresolved Go version %q for %s", version, input.Name)
		}
	case "npm":
		if !exactNPMVersionPattern.MatchString(version) {
			return fmt.Errorf("SBOM requires exact npm version for %s, got %q", input.Name, version)
		}
	case "vcpkg":
		return fmt.Errorf("SBOM cannot claim exact vcpkg identity for %s from constraint %q without a resolved lock/baseline", input.Name, version)
	default:
		return fmt.Errorf("SBOM does not support ecosystem %q for %s", input.Ecosystem, input.Name)
	}
	if input.Scope != "runtime" && input.Scope != "development" {
		return fmt.Errorf("SBOM component %s has invalid scope %q", input.Name, input.Scope)
	}
	return nil
}

func dedupeSBOMInputs(values []sbomInputComponent) []sbomInputComponent {
	sort.Slice(values, func(i, j int) bool {
		a, b := values[i], values[j]
		ka := strings.Join([]string{a.Scope, a.Ecosystem, a.Name, a.Version, a.Manifest}, "\x00")
		kb := strings.Join([]string{b.Scope, b.Ecosystem, b.Name, b.Version, b.Manifest}, "\x00")
		return ka < kb
	})
	if len(values) == 0 {
		return nil
	}
	out := values[:0]
	var previous string
	for i, value := range values {
		key := strings.Join([]string{value.Scope, value.Ecosystem, value.Name, value.Version, value.Manifest}, "\x00")
		if i == 0 || key != previous {
			out = append(out, value)
			previous = key
		}
	}
	return out
}

func buildCycloneDXBOM(sourceSHA string, inputs []sbomInputComponent) (CycloneDXBOM, int, int, error) {
	sourceSHA, err := validateFullCommitID(sourceSHA)
	if err != nil {
		return CycloneDXBOM{}, 0, 0, err
	}
	rootRef := "cmdr:source:" + sourceSHA
	bom := CycloneDXBOM{
		Schema:      cycloneDXSchemaURL,
		BOMFormat:   "CycloneDX",
		SpecVersion: cycloneDXSpecVersion,
		Version:     1,
		Metadata: CDXMetadata{
			Tools: CDXTools{Components: []CDXComponent{{
				Type: "application", BOMRef: "cmdr:tool:cmdr-dev:" + sbomGeneratorVersion,
				Name: "cmdr-dev", Version: sbomGeneratorVersion,
				Properties: sortedCDXProperties(map[string]string{
					"cmdr:toolchain:go": engineeringGoVersion,
				}),
			}}},
			Component: CDXComponent{
				Type: "application", BOMRef: rootRef, Name: "cmdr", Version: sourceSHA,
				Properties: sortedCDXProperties(map[string]string{
					"cmdr:source:git-sha":  sourceSHA,
					"cmdr:sbom:generator":  "cmdr-dev",
					"cmdr:sbom:schema":     cycloneDXSpecVersion,
				}),
			},
		},
		Compositions: []CDXComposition{{Aggregate: "incomplete", Assemblies: []string{rootRef}}},
	}

	var runtimeCount, developmentCount int
	for _, input := range inputs {
		scope := "required"
		if input.Scope == "development" {
			scope = "excluded"
			developmentCount++
		} else {
			runtimeCount++
		}
		properties := map[string]string{
			"cmdr:dependency:scope":  input.Scope,
			"cmdr:dependency:direct": fmt.Sprintf("%t", input.Direct),
			"cmdr:dependency:ecosystem": input.Ecosystem,
			"cmdr:dependency:manifest": input.Manifest,
		}
		if input.Commit != "" {
			properties["cmdr:tool:commit"] = input.Commit
		}
		if input.License != "" {
			properties["cmdr:tool:license"] = input.License
		}
		bom.Components = append(bom.Components, CDXComponent{
			Type: input.Kind, BOMRef: sbomComponentRef(input),
			Name: input.Name, Version: input.Version, Scope: scope,
			Properties: sortedCDXProperties(properties),
		})
	}
	sort.Slice(bom.Components, func(i, j int) bool { return bom.Components[i].BOMRef < bom.Components[j].BOMRef })
	return bom, runtimeCount, developmentCount, nil
}

func sbomComponentRef(input sbomInputComponent) string {
	raw := strings.Join([]string{input.Scope, input.Ecosystem, input.Name, input.Version, input.Manifest}, "\x00")
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("cmdr:component:%x", sum[:])
}

func sortedCDXProperties(values map[string]string) []CDXProperty {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]CDXProperty, 0, len(keys))
	for _, key := range keys {
		out = append(out, CDXProperty{Name: key, Value: values[key]})
	}
	return out
}
