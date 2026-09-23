package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const architectureRegistryPath = "engineering/architecture/boundaries.json"

type ArchitectureBoundary struct {
	ID                      string   `json:"id"`
	Kind                    string   `json:"kind"`
	Roots                   []string `json:"roots"`
	MutableByImplementation bool     `json:"mutable_by_implementation"`
	MayDependOn             []string `json:"may_depend_on"`
}

type ArchitectureRegistry struct {
	SchemaVersion int                    `json:"schema_version"`
	DefaultPolicy string                 `json:"default_policy"`
	Boundaries    []ArchitectureBoundary `json:"boundaries"`
}

type ArchitectureAuditSummary struct {
	Boundaries          int  `json:"boundaries"`
	StrictV2Manifests   int  `json:"strict_v2_manifests"`
	AllowedPathClaims   int  `json:"allowed_path_claims"`
	ProductSpecReadOnly bool `json:"product_spec_read_only"`
}

func runArchitectureAudit(root string) (ArchitectureAuditSummary, error) {
	var registry ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &registry); err != nil {
		return ArchitectureAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return ArchitectureAuditSummary{}, err
	}

	paths, err := manifestPaths(root)
	if err != nil {
		return ArchitectureAuditSummary{}, err
	}
	summary := ArchitectureAuditSummary{
		Boundaries:          len(registry.Boundaries),
		ProductSpecReadOnly: true,
	}
	for _, path := range paths {
		header, err := decodeManifestHeader(path)
		if err != nil {
			return ArchitectureAuditSummary{}, err
		}
		if header.SchemaVersion != 2 {
			continue
		}
		manifest, err := decodeWorkManifestV2(path)
		if err != nil {
			return ArchitectureAuditSummary{}, err
		}
		if err := validateManifestPathOwnership(manifest, registry); err != nil {
			return ArchitectureAuditSummary{}, fmt.Errorf("%s: %w", manifest.ID, err)
		}
		summary.StrictV2Manifests++
		summary.AllowedPathClaims += len(manifest.AllowedPaths)
	}
	return summary, nil
}

func validateArchitectureRegistry(registry ArchitectureRegistry) error {
	if registry.SchemaVersion != 1 {
		return fmt.Errorf("unsupported architecture registry schema_version %d", registry.SchemaVersion)
	}
	if registry.DefaultPolicy != "deny-unregistered-path" {
		return fmt.Errorf("architecture default policy must be deny-unregistered-path")
	}
	if len(registry.Boundaries) == 0 {
		return fmt.Errorf("architecture registry has no boundaries")
	}

	ids := map[string]ArchitectureBoundary{}
	rootOwners := map[string]string{}
	for _, boundary := range registry.Boundaries {
		if strings.TrimSpace(boundary.ID) == "" || strings.TrimSpace(boundary.Kind) == "" {
			return fmt.Errorf("architecture boundary id and kind are required")
		}
		if _, exists := ids[boundary.ID]; exists {
			return fmt.Errorf("duplicate architecture boundary %s", boundary.ID)
		}
		if len(boundary.Roots) == 0 {
			return fmt.Errorf("architecture boundary %s has no roots", boundary.ID)
		}
		ids[boundary.ID] = boundary
		for _, root := range boundary.Roots {
			normalized, err := normalizeArchitecturePattern(root)
			if err != nil {
				return fmt.Errorf("boundary %s: %w", boundary.ID, err)
			}
			for existing, owner := range rootOwners {
				if patternsOverlap(normalized, existing) {
					return fmt.Errorf("architecture roots %q (%s) and %q (%s) overlap", normalized, boundary.ID, existing, owner)
				}
			}
			rootOwners[normalized] = boundary.ID
		}
	}

	product, ok := ids["product-spec"]
	if !ok {
		return fmt.Errorf("product-spec boundary is required")
	}
	if product.MutableByImplementation {
		return fmt.Errorf("product-spec boundary must be non-mutable")
	}
	if len(product.Roots) != 1 || product.Roots[0] != "cmdr-product-spec/**" {
		return fmt.Errorf("product-spec boundary must own exactly cmdr-product-spec/**")
	}

	for _, boundary := range registry.Boundaries {
		for _, dep := range boundary.MayDependOn {
			if dep == boundary.ID {
				return fmt.Errorf("boundary %s cannot depend on itself", boundary.ID)
			}
			if _, ok := ids[dep]; !ok {
				return fmt.Errorf("boundary %s depends on unknown boundary %s", boundary.ID, dep)
			}
		}
	}
	if err := validateBoundaryDependencyAcyclic(registry.Boundaries); err != nil {
		return err
	}
	return nil
}

func validateManifestPathOwnership(manifest WorkManifestV2, registry ArchitectureRegistry) error {
	for _, claim := range manifest.AllowedPaths {
		normalized, err := normalizeArchitecturePattern(claim)
		if err != nil {
			return err
		}
		var owners []ArchitectureBoundary
		for _, boundary := range registry.Boundaries {
			for _, root := range boundary.Roots {
				if patternCovers(root, normalized) {
					owners = append(owners, boundary)
					break
				}
			}
		}
		if len(owners) == 0 {
			return fmt.Errorf("allowed path %q is outside registered architecture boundaries", claim)
		}
		if len(owners) > 1 {
			return fmt.Errorf("allowed path %q has ambiguous architecture ownership", claim)
		}
		if !owners[0].MutableByImplementation {
			return fmt.Errorf("allowed path %q belongs to non-mutable boundary %s", claim, owners[0].ID)
		}
	}

	productForbidden := false
	for _, forbidden := range manifest.ForbiddenPaths {
		normalized, err := normalizeArchitecturePattern(forbidden)
		if err != nil {
			return err
		}
		if patternCovers(normalized, "cmdr-product-spec/**") {
			productForbidden = true
			break
		}
	}
	if !productForbidden {
		return fmt.Errorf("manifest must explicitly forbid cmdr-product-spec/**")
	}
	return nil
}

func normalizeArchitecturePattern(pattern string) (string, error) {
	pattern = filepath.ToSlash(strings.TrimSpace(pattern))
	if pattern == "" || filepath.IsAbs(pattern) || strings.HasPrefix(pattern, "../") || strings.Contains(pattern, "/../") {
		return "", fmt.Errorf("invalid architecture path pattern %q", pattern)
	}
	return pattern, nil
}

func patternPrefix(pattern string) string {
	return strings.TrimSuffix(pattern, "**")
}

func patternCovers(root, claim string) bool {
	root = filepath.ToSlash(strings.TrimSpace(root))
	claim = filepath.ToSlash(strings.TrimSpace(claim))
	if root == claim {
		return true
	}
	rootPrefix := patternPrefix(root)
	claimPrefix := patternPrefix(claim)
	return rootPrefix != "" && strings.HasPrefix(claimPrefix, rootPrefix)
}

func patternsOverlap(a, b string) bool {
	return patternCovers(a, b) || patternCovers(b, a)
}

func validateBoundaryDependencyAcyclic(boundaries []ArchitectureBoundary) error {
	index := map[string]ArchitectureBoundary{}
	for _, boundary := range boundaries {
		index[boundary.ID] = boundary
	}
	const (
		visiting = 1
		done     = 2
	)
	state := map[string]int{}
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case visiting:
			return fmt.Errorf("architecture boundary dependency cycle at %s", id)
		case done:
			return nil
		}
		state[id] = visiting
		deps := append([]string(nil), index[id].MayDependOn...)
		sort.Strings(deps)
		for _, dep := range deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = done
		return nil
	}
	ids := make([]string, 0, len(index))
	for id := range index {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
