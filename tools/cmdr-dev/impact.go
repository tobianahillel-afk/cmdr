package main

import (
	"bufio"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type ImpactEvidence struct {
	Domain string   `json:"domain"`
	Reason string   `json:"reason"`
	Paths  []string `json:"paths"`
}

type ImpactReport struct {
	WorkUnit        string           `json:"work_unit"`
	ChangedPaths    []string         `json:"changed_paths"`
	RiskDomains     []string         `json:"risk_domains"`
	HighRisk        bool             `json:"high_risk"`
	HighRiskReasons []string         `json:"high_risk_reasons,omitempty"`
	UnknownPaths    []string         `json:"unknown_paths,omitempty"`
	Evidence        []ImpactEvidence `json:"evidence"`
}

type impactAccumulator struct {
	domains         map[string]bool
	evidence        map[string]map[string]bool
	highRiskReasons map[string]bool
	unknownPaths    map[string]bool
}

func runImpactAnalysis(root, workUnit, changesFile string, state CurrentState, graph WorkGraph) (ImpactReport, error) {
	if strings.TrimSpace(changesFile) == "" {
		return ImpactReport{}, fmt.Errorf("impact analysis requires --changes-file")
	}
	changedPaths, err := readChangedPaths(root, changesFile)
	if err != nil {
		return ImpactReport{}, err
	}
	if len(changedPaths) == 0 {
		return ImpactReport{}, fmt.Errorf("impact analysis received no changed paths")
	}
	if workUnit == "" {
		workUnit = state.Execution.ActiveWorkUnit
	}
	node, ok := graphNode(graph, workUnit)
	if !ok {
		return ImpactReport{}, fmt.Errorf("impact work unit %s is missing from graph", workUnit)
	}
	if node.Status == "BLOCKED" || node.Status == "BLOCKED_DECISION" {
		return ImpactReport{}, fmt.Errorf("impact work unit %s is blocked", workUnit)
	}
	manifestPath := filepath.Join(root, "work", "lots", workUnit, "manifest.json")
	header, err := decodeManifestHeader(root, manifestPath)
	if err != nil {
		return ImpactReport{}, err
	}
	if header.SchemaVersion != 2 {
		return ImpactReport{}, fmt.Errorf("impact analysis requires strict manifest v2")
	}
	manifest, err := decodeWorkManifestV2(root, manifestPath)
	if err != nil {
		return ImpactReport{}, err
	}
	catalog, err := loadCheckCatalog(root)
	if err != nil {
		return ImpactReport{}, err
	}
	var registry ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &registry); err != nil {
		return ImpactReport{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return ImpactReport{}, err
	}
	return computeImpact(workUnit, changedPaths, manifest, catalog, registry)
}

func computeImpact(workUnit string, changedPaths []string, manifest WorkManifestV2, catalog CheckCatalog, registry ArchitectureRegistry) (ImpactReport, error) {
	acc := impactAccumulator{
		domains:         map[string]bool{},
		evidence:        map[string]map[string]bool{},
		highRiskReasons: map[string]bool{},
		unknownPaths:    map[string]bool{},
	}
	normalized, err := normalizeChangedPaths(changedPaths)
	if err != nil {
		return ImpactReport{}, err
	}

	for _, path := range normalized {
		matchedCheck := false
		for _, check := range catalog.Checks {
			for _, pattern := range check.TriggerPaths {
				if triggerPatternMatches(pattern, path) {
					matchedCheck = true
					for _, domain := range check.RiskDomains {
						acc.add(domain, "check-trigger:"+check.ID, path)
					}
					break
				}
			}
		}

		boundary, owned := architectureBoundaryForPath(registry, path)
		if !owned {
			acc.add("unknown", "unowned-repository-path", path)
			acc.add("security", "unowned-repository-path", path)
			acc.unknownPaths[path] = true
			acc.highRiskReasons["unowned repository path: "+path] = true
			continue
		}
		switch boundary.Kind {
		case "product-documentation":
			acc.add("product-spec", "architecture-boundary:"+boundary.ID, path)
			acc.add("coverage-traceability", "architecture-boundary:"+boundary.ID, path)
		case "product-runtime":
			acc.add("architecture", "architecture-boundary:"+boundary.ID, path)
			acc.add("runtime-dependencies", "architecture-boundary:"+boundary.ID, path)
			acc.add("security", "architecture-boundary:"+boundary.ID, path)
			acc.highRiskReasons["product-runtime change: "+path] = true
		case "execution":
			acc.add("work-governance", "architecture-boundary:"+boundary.ID, path)
		case "ci":
			acc.add("work-governance", "architecture-boundary:"+boundary.ID, path)
			acc.add("repository-health", "architecture-boundary:"+boundary.ID, path)
			acc.add("security", "architecture-boundary:"+boundary.ID, path)
			acc.highRiskReasons["CI control-plane change: "+path] = true
		}
		if strings.HasPrefix(path, "engineering/security/") {
			acc.add("security", "security-control-plane-change", path)
			acc.highRiskReasons["security control-plane change: "+path] = true
		}
		if strings.HasPrefix(path, "engineering/architecture/") {
			acc.add("architecture", "architecture-control-plane-change", path)
			acc.add("security", "architecture-control-plane-change", path)
			acc.highRiskReasons["architecture control-plane change: "+path] = true
		}
		if !matchedCheck {
			acc.add("unknown", "no-check-trigger-matched", path)
			acc.unknownPaths[path] = true
			acc.highRiskReasons["no validation trigger matched: "+path] = true
		}
	}

	acc.addManifestScope(manifest)
	return acc.report(workUnit, normalized), nil
}

func (a *impactAccumulator) add(domain, reason, path string) {
	if !knownRiskDomains[domain] {
		return
	}
	a.domains[domain] = true
	key := domain + "\x00" + reason
	if a.evidence[key] == nil {
		a.evidence[key] = map[string]bool{}
	}
	if path != "" {
		a.evidence[key][path] = true
	}
}

func (a *impactAccumulator) addManifestScope(manifest WorkManifestV2) {
	refs := manifest.ProductRefs
	if len(refs.Capabilities)+len(refs.Requirements)+len(refs.Screens)+len(refs.ImplementationContracts)+len(refs.CanonicalObjects) > 0 {
		a.add("product-spec", "manifest-product-references", manifest.ID)
		a.add("coverage-traceability", "manifest-product-references", manifest.ID)
	}
	if len(refs.Permissions)+len(refs.OpenDecisions) > 0 {
		a.add("product-spec", "manifest-security-product-references", manifest.ID)
		a.add("security", "manifest-security-product-references", manifest.ID)
		a.highRiskReasons["manifest contains security-sensitive product references"] = true
	}
}

func (a *impactAccumulator) report(workUnit string, changedPaths []string) ImpactReport {
	report := ImpactReport{
		WorkUnit:     workUnit,
		ChangedPaths: append([]string(nil), changedPaths...),
		HighRisk:     len(a.highRiskReasons) > 0,
	}
	for domain := range a.domains {
		report.RiskDomains = append(report.RiskDomains, domain)
	}
	sort.Strings(report.RiskDomains)
	for reason := range a.highRiskReasons {
		report.HighRiskReasons = append(report.HighRiskReasons, reason)
	}
	sort.Strings(report.HighRiskReasons)
	for path := range a.unknownPaths {
		report.UnknownPaths = append(report.UnknownPaths, path)
	}
	sort.Strings(report.UnknownPaths)

	keys := make([]string, 0, len(a.evidence))
	for key := range a.evidence {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts := strings.SplitN(key, "\x00", 2)
		evidence := ImpactEvidence{Domain: parts[0], Reason: parts[1]}
		for path := range a.evidence[key] {
			evidence.Paths = append(evidence.Paths, path)
		}
		sort.Strings(evidence.Paths)
		report.Evidence = append(report.Evidence, evidence)
	}
	return report
}

func readChangedPaths(root, path string) ([]string, error) {
	f, err := openRepoFile(root, path)
	if err != nil {
		return nil, fmt.Errorf("open changes file: %w", err)
	}
	defer f.Close()
	var paths []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value != "" {
			paths = append(paths, value)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read changes file: %w", err)
	}
	return normalizeChangedPaths(paths)
}

func normalizeChangedPaths(paths []string) ([]string, error) {
	seen := map[string]bool{}
	var out []string
	for _, path := range paths {
		path = filepath.ToSlash(strings.TrimSpace(path))
		if path == "" {
			continue
		}
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
		if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../") {
			return nil, fmt.Errorf("invalid changed path %q", path)
		}
		if !seen[clean] {
			seen[clean] = true
			out = append(out, clean)
		}
	}
	sort.Strings(out)
	return out, nil
}

func triggerPatternMatches(pattern, path string) bool {
	pattern = filepath.ToSlash(strings.TrimSpace(pattern))
	path = filepath.ToSlash(strings.TrimSpace(path))
	if strings.HasSuffix(pattern, "/**") {
		base := strings.TrimSuffix(pattern, "/**")
		return path == base || strings.HasPrefix(path, base+"/")
	}
	return path == pattern
}

func architectureBoundaryForPath(registry ArchitectureRegistry, path string) (ArchitectureBoundary, bool) {
	var found *ArchitectureBoundary
	for i := range registry.Boundaries {
		boundary := &registry.Boundaries[i]
		for _, root := range boundary.Roots {
			if triggerPatternMatches(root, path) {
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
