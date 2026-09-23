package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ContextSource struct {
	Path      string   `json:"path"`
	Authority string   `json:"authority"`
	Reasons   []string `json:"reasons"`
	SHA256    string   `json:"sha256"`
	Content   string   `json:"content"`
}

type AgentContextBundle struct {
	SchemaVersion       int             `json:"schema_version"`
	WorkUnit            string          `json:"work_unit"`
	ProductSpecBaseline string          `json:"product_spec_baseline"`
	DependencyUnits     []string        `json:"dependency_units"`
	Sources             []ContextSource `json:"sources"`
	BundleDigest        string          `json:"bundle_digest"`
}

type ContextSummary struct {
	WorkUnit           string `json:"work_unit"`
	Sources            int    `json:"sources"`
	ProductSources     int    `json:"product_sources"`
	EngineeringSources int    `json:"engineering_sources"`
	ExecutionSources   int    `json:"execution_sources"`
	DependencyUnits    int    `json:"dependency_units"`
	BundleDigest       string `json:"bundle_digest"`
	Output             string `json:"output"`
	Mode               string `json:"mode"`
}

type contextCollector struct {
	root    string
	sources map[string]*ContextSource
}

func runContextCompiler(root, workUnit, output string, check bool, state CurrentState, graph WorkGraph) (ContextSummary, error) {
	if workUnit == "" {
		workUnit = state.Execution.ActiveWorkUnit
	}
	node, ok := graphNode(graph, workUnit)
	if !ok {
		return ContextSummary{}, fmt.Errorf("context work unit %s is missing from graph", workUnit)
	}
	if node.Status == "BLOCKED" || node.Status == "BLOCKED_DECISION" {
		return ContextSummary{}, fmt.Errorf("context work unit %s is blocked", workUnit)
	}

	manifestPath := filepath.Join(root, "work", "lots", workUnit, "manifest.json")
	header, err := decodeManifestHeader(manifestPath)
	if err != nil {
		return ContextSummary{}, err
	}
	if header.SchemaVersion != 2 {
		return ContextSummary{}, fmt.Errorf("%s: context compilation requires strict manifest v2", workUnit)
	}
	manifest, err := decodeWorkManifestV2(manifestPath)
	if err != nil {
		return ContextSummary{}, err
	}

	dependencies, err := dependencyClosure(workUnit, graph)
	if err != nil {
		return ContextSummary{}, err
	}
	collector := contextCollector{root: root, sources: map[string]*ContextSource{}}

	for _, source := range []struct {
		path   string
		reason string
	}{
		{"AGENTS.md", "repository-wide autonomous engineering rules"},
		{"AI_START_HERE.md", "universal autonomous-agent entry protocol"},
		{"engineering/governance/source-authority.md", "implementation authority and conflict rules"},
		{"engineering/work/complexity-policy.md", "work-unit sizing and split policy"},
		{statePath, "current execution cache for reconciliation"},
		{graphPath, "canonical engineering work graph"},
	} {
		if err := collector.add(source.path, source.reason); err != nil {
			return ContextSummary{}, err
		}
	}
	if err := collector.add(filepath.ToSlash(filepath.Join("work", "lots", workUnit, "manifest.json")), "selected work-unit manifest"); err != nil {
		return ContextSummary{}, err
	}
	if manifest.Parent != "" {
		parentPath := filepath.ToSlash(filepath.Join("work", "lots", manifest.Parent, "manifest.json"))
		if err := collector.add(parentPath, "parent work-unit context"); err != nil {
			return ContextSummary{}, err
		}
	}

	for _, dep := range dependencies {
		base := filepath.ToSlash(filepath.Join("work", "lots", dep))
		if err := collector.add(base+"/manifest.json", "dependency manifest: "+dep); err != nil {
			return ContextSummary{}, err
		}
		for _, name := range []string{"HANDOFF.json", "PROGRESS.json"} {
			rel := base + "/" + name
			if fileExists(filepath.Join(root, filepath.FromSlash(rel))) {
				if err := collector.add(rel, "dependency evidence: "+dep); err != nil {
					return ContextSummary{}, err
				}
			}
		}
	}

	productSources, err := resolveProductContextSources(root, state.ProductSpec.CanonicalPath, manifest.ProductRefs)
	if err != nil {
		return ContextSummary{}, err
	}
	for path, reasons := range productSources {
		for _, reason := range reasons {
			if err := collector.add(path, reason); err != nil {
				return ContextSummary{}, err
			}
		}
	}

	bundle, err := collector.bundle(workUnit, state.ProductSpec.BaselineCommit, dependencies)
	if err != nil {
		return ContextSummary{}, err
	}
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return ContextSummary{}, err
	}
	data = append(data, '\n')

	if output == "" || output == "engineering/spec-index/inventory.json" {
		output = filepath.ToSlash(filepath.Join("engineering", "context", workUnit+".json"))
	}
	outputPath, err := resolveRepoPath(root, output, !check)
	if err != nil {
		return ContextSummary{}, err
	}
	mode := "write"
	if check {
		mode = "check"
		existing, err := readRepoFile(root, outputPath)
		if err != nil {
			return ContextSummary{}, fmt.Errorf("read context bundle: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return ContextSummary{}, fmt.Errorf("context bundle is stale: %s", outputPath)
		}
	} else {
		outputPath, err = writeRepoFile(root, outputPath, data)
		if err != nil {
			return ContextSummary{}, err
		}
	}

	summary := ContextSummary{
		WorkUnit:        workUnit,
		Sources:         len(bundle.Sources),
		DependencyUnits: len(dependencies),
		BundleDigest:    bundle.BundleDigest,
		Output:          outputPath,
		Mode:            mode,
	}
	for _, source := range bundle.Sources {
		switch source.Authority {
		case "product":
			summary.ProductSources++
		case "engineering":
			summary.EngineeringSources++
		case "execution":
			summary.ExecutionSources++
		}
	}
	return summary, nil
}

func (c *contextCollector) add(rel, reason string) error {
	rel = filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
	if rel == "." || strings.HasPrefix(rel, "../") || filepath.IsAbs(rel) {
		return fmt.Errorf("invalid context source path %q", rel)
	}
	data, err := os.ReadFile(filepath.Join(c.root, filepath.FromSlash(rel)))
	if err != nil {
		return fmt.Errorf("read context source %s: %w", rel, err)
	}
	sum := sha256.Sum256(data)
	if existing, ok := c.sources[rel]; ok {
		existing.Reasons = appendUnique(existing.Reasons, reason)
		sort.Strings(existing.Reasons)
		return nil
	}
	c.sources[rel] = &ContextSource{
		Path:      rel,
		Authority: sourceAuthority(rel),
		Reasons:   []string{reason},
		SHA256:    fmt.Sprintf("%x", sum),
		Content:   string(data),
	}
	return nil
}

func (c *contextCollector) bundle(workUnit, baseline string, dependencies []string) (AgentContextBundle, error) {
	sources := make([]ContextSource, 0, len(c.sources))
	for _, source := range c.sources {
		source.Reasons = append([]string(nil), source.Reasons...)
		sort.Strings(source.Reasons)
		sources = append(sources, *source)
	}
	sort.Slice(sources, func(i, j int) bool {
		ri, rj := authorityRank(sources[i].Authority), authorityRank(sources[j].Authority)
		if ri != rj {
			return ri < rj
		}
		return sources[i].Path < sources[j].Path
	})
	dependencies = append([]string(nil), dependencies...)
	sort.Strings(dependencies)
	bundle := AgentContextBundle{
		SchemaVersion:       1,
		WorkUnit:            workUnit,
		ProductSpecBaseline: baseline,
		DependencyUnits:     dependencies,
		Sources:             sources,
	}
	h := sha256.New()
	fmt.Fprintf(h, "schema=1\nwork=%s\nbaseline=%s\n", workUnit, baseline)
	for _, dep := range dependencies {
		fmt.Fprintf(h, "dep=%s\n", dep)
	}
	for _, source := range sources {
		fmt.Fprintf(h, "%s\x00%s\x00%s\x00%s\n",
			source.Authority, source.Path, source.SHA256, strings.Join(source.Reasons, "\x1f"))
	}
	bundle.BundleDigest = fmt.Sprintf("%x", h.Sum(nil))
	return bundle, nil
}

func resolveProductContextSources(root, specRel string, refs ProductRefsV2) (map[string][]string, error) {
	out := map[string][]string{}
	add := func(path, reason string) {
		out[path] = appendUnique(out[path], reason)
	}
	if len(refs.Capabilities) > 0 {
		active, err := loadActiveCapabilityIDs(root, specRel)
		if err != nil {
			return nil, err
		}
		for _, id := range refs.Capabilities {
			if _, ok := active[id]; !ok {
				return nil, fmt.Errorf("capability %s is not active in the Capability Register", id)
			}
			path, err := capabilityCanonicalPath(root, specRel, id)
			if err != nil {
				return nil, err
			}
			add(path, "canonical capability source: "+id)
		}
	}
	if len(refs.Screens) > 0 {
		active, err := loadActiveScreenIDs(root, specRel)
		if err != nil {
			return nil, err
		}
		for _, id := range refs.Screens {
			if _, ok := active[id]; !ok {
				return nil, fmt.Errorf("screen %s is not active in the Screen Register", id)
			}
			path, err := registryCanonicalPath(root, filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers", "screen-register.md")), id, "Source canonique")
			if err != nil {
				return nil, err
			}
			add(path, "canonical screen source: "+id)
		}
	}
	if len(refs.Requirements) > 0 {
		active, err := loadActiveRequirementIDs(root, specRel)
		if err != nil {
			return nil, err
		}
		for _, id := range refs.Requirements {
			if _, ok := active[id]; !ok {
				return nil, fmt.Errorf("requirement %s is outside the canonical 122-ID catalog", id)
			}
		}
		add(filepath.ToSlash(filepath.Join(specRel, "00-governance", "source-material", "requirements-traceability-matrix.md")), "current product requirement traceability source")
		add("engineering/coverage/source-requirements-baseline.json", "source requirement identity baseline")
	}
	if len(refs.Permissions) > 0 {
		active, err := loadActivePermissionIDs(root, specRel)
		if err != nil {
			return nil, err
		}
		for _, id := range refs.Permissions {
			if _, ok := active[id]; !ok {
				return nil, fmt.Errorf("permission %s is not active in the Permission Register", id)
			}
		}
		add(filepath.ToSlash(filepath.Join(specRel, "14-security-permissions-and-trust", "permission-model.md")), "canonical permission model")
		add(filepath.ToSlash(filepath.Join(specRel, "14-security-permissions-and-trust", "permission-catalog.md")), "canonical permission identifiers")
	}
	if len(refs.OpenDecisions) > 0 {
		active, err := loadActiveOpenDecisionIDs(root, specRel)
		if err != nil {
			return nil, err
		}
		for _, id := range refs.OpenDecisions {
			if _, ok := active[id]; !ok {
				return nil, fmt.Errorf("open decision %s is not active", id)
			}
		}
		add(filepath.ToSlash(filepath.Join(specRel, "00-governance", "source-material", "unresolved-decisions.md")), "active open-decision authority")
	}
	for _, path := range refs.ImplementationContracts {
		if err := validateExplicitProductPath(root, specRel, path); err != nil {
			return nil, err
		}
		add(path, "explicit implementation contract")
	}
	for _, path := range refs.CanonicalObjects {
		if err := validateExplicitProductPath(root, specRel, path); err != nil {
			return nil, err
		}
		add(path, "explicit canonical object source")
	}
	for path := range out {
		sort.Strings(out[path])
	}
	return out, nil
}

func capabilityCanonicalPath(root, specRel, id string) (string, error) {
	dirRel := filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers"))
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dirRel)))
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "capability-register-") || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		rel := filepath.ToSlash(filepath.Join(dirRel, entry.Name()))
		path, err := registryCanonicalPath(root, rel, id, "Canonical file")
		if err == nil {
			return path, nil
		}
		if !strings.Contains(err.Error(), "not found in registry") {
			return "", err
		}
	}
	return "", fmt.Errorf("canonical capability source for %s not found", id)
}

func registryCanonicalPath(root, registryRel, id, column string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(registryRel)))
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	columnIndex := -1
	for _, line := range lines {
		cells := markdownCells(line)
		if len(cells) == 0 {
			continue
		}
		if columnIndex < 0 && strings.EqualFold(strings.TrimSpace(cells[0]), "ID") || columnIndex < 0 && strings.EqualFold(strings.TrimSpace(cells[0]), "Permission") {
			for i, cell := range cells {
				if strings.EqualFold(strings.TrimSpace(cell), column) {
					columnIndex = i
					break
				}
			}
			continue
		}
		if columnIndex >= 0 && cleanCell(cells[0]) == id && columnIndex < len(cells) {
			path, err := canonicalPathFromCell(registryRel, cells[columnIndex])
			if err != nil {
				return "", err
			}
			if err := validateExplicitProductPath(root, strings.Split(registryRel, "/")[0], path); err != nil {
				return "", err
			}
			return path, nil
		}
	}
	return "", fmt.Errorf("%s not found in registry %s", id, registryRel)
}

func markdownCells(line string) []string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") {
		return nil
	}
	parts := strings.Split(strings.Trim(line, "|"), "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func cleanCell(value string) string {
	return strings.Trim(strings.TrimSpace(value), "`*")
}

func canonicalPathFromCell(registryRel, cell string) (string, error) {
	cell = strings.TrimSpace(cell)
	if open := strings.Index(cell, "]("); open >= 0 {
		end := strings.Index(cell[open+2:], ")")
		if end < 0 {
			return "", fmt.Errorf("invalid markdown link in %s", registryRel)
		}
		cell = cell[open+2 : open+2+end]
	} else {
		cell = cleanCell(cell)
	}
	if cell == "" || strings.Contains(cell, "://") {
		return "", fmt.Errorf("invalid canonical path %q in %s", cell, registryRel)
	}
	if !strings.HasPrefix(cell, "../") && !strings.HasPrefix(cell, "./") && !strings.HasPrefix(cell, "cmdr-product-spec/") {
		cell = filepath.ToSlash(filepath.Join(strings.Split(registryRel, "/")[0], cell))
	} else if !strings.HasPrefix(cell, "cmdr-product-spec/") {
		cell = filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(registryRel), filepath.FromSlash(cell))))
	}
	return filepath.ToSlash(cell), nil
}

func validateExplicitProductPath(root, specRel, rel string) error {
	rel = filepath.ToSlash(filepath.Clean(filepath.FromSlash(rel)))
	if !strings.HasPrefix(rel, strings.TrimSuffix(specRel, "/")+"/") {
		return fmt.Errorf("product source %s is outside %s", rel, specRel)
	}
	info, err := statRepoPath(root, rel)
	if err != nil || info.IsDir() {
		return fmt.Errorf("product source %s does not exist as a file", rel)
	}
	return nil
}

func dependencyClosure(rootID string, graph WorkGraph) ([]string, error) {
	index := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		index[node.ID] = node
	}
	if _, ok := index[rootID]; !ok {
		return nil, fmt.Errorf("work unit %s is missing from graph", rootID)
	}
	seen := map[string]bool{}
	var visit func(string) error
	visit = func(id string) error {
		for _, dep := range index[id].DependsOn {
			if _, ok := index[dep]; !ok {
				return fmt.Errorf("work unit %s depends on missing %s", id, dep)
			}
			if seen[dep] {
				continue
			}
			seen[dep] = true
			if err := visit(dep); err != nil {
				return err
			}
		}
		return nil
	}
	if err := visit(rootID); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	sort.Strings(out)
	return out, nil
}

func graphNode(graph WorkGraph, id string) (WorkNode, bool) {
	for _, node := range graph.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return WorkNode{}, false
}

func sourceAuthority(path string) string {
	switch {
	case strings.HasPrefix(path, "cmdr-product-spec/"):
		return "product"
	case strings.HasPrefix(path, "engineering/"), path == "AGENTS.md", path == "AI_START_HERE.md":
		return "engineering"
	default:
		return "execution"
	}
}

func authorityRank(authority string) int {
	switch authority {
	case "product":
		return 1
	case "engineering":
		return 2
	default:
		return 3
	}
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
