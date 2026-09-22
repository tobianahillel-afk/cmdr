package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CoverageGap struct {
	Code        string   `json:"code"`
	SubjectID   string   `json:"subject_id"`
	SubjectKind string   `json:"subject_kind"`
	SourcePaths []string `json:"source_paths"`
	Detail      string   `json:"detail"`
}

type CoverageAuditReport struct {
	SchemaVersion  int           `json:"schema_version"`
	SpecTreeDigest string        `json:"spec_tree_digest"`
	Gaps           []CoverageGap `json:"gaps"`
}

type CoverageAuditSummary struct {
	Gaps                   int            `json:"gaps"`
	ByCode                 map[string]int      `json:"by_code"`
	SubjectsByCode         map[string][]string `json:"subjects_by_code"`
	RegisteredCapabilities int                 `json:"registered_capabilities"`
	RegisteredRequirements int            `json:"registered_requirements"`
	RegisteredScreens      int            `json:"registered_screens"`
	RegisteredPermissions  int            `json:"registered_permissions"`
	ActiveOpenDecisions    int            `json:"active_open_decisions"`
	SpecTreeDigest         string         `json:"spec_tree_digest"`
	Output                 string         `json:"output"`
	Mode                   string         `json:"mode"`
}

func runCoverageAudit(root, specRel, output string, check bool) (CoverageAuditSummary, error) {
	inventory, err := buildSpecInventory(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	screens, err := loadActiveScreenIDs(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	graph, err := buildProductGraph(inventory, screens)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	capabilities, err := loadActiveCapabilityIDs(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	requirements, err := loadActiveRequirementIDs(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	permissions, err := loadActivePermissionIDs(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	openDecisions, err := loadActiveOpenDecisionIDs(root, specRel)
	if err != nil {
		return CoverageAuditSummary{}, err
	}

	report := buildCoverageAudit(graph, capabilities, requirements, screens, permissions, openDecisions, specRel)
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return CoverageAuditSummary{}, err
	}
	data = append(data, '\n')

	outputPath := output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(root, filepath.FromSlash(outputPath))
	}
	mode := "write"
	if check {
		mode = "check"
		existing, err := os.ReadFile(outputPath)
		if err != nil {
			return CoverageAuditSummary{}, fmt.Errorf("read coverage audit: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return CoverageAuditSummary{}, fmt.Errorf("coverage audit is stale: %s", outputPath)
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
			return CoverageAuditSummary{}, err
		}
		if err := os.WriteFile(outputPath, data, 0o644); err != nil {
			return CoverageAuditSummary{}, err
		}
	}

	summary := CoverageAuditSummary{
		Gaps:                   len(report.Gaps),
		ByCode:                 map[string]int{},
		SubjectsByCode:         map[string][]string{},
		RegisteredCapabilities: len(capabilities),
		RegisteredRequirements: len(requirements),
		RegisteredScreens:      len(screens),
		RegisteredPermissions:  len(permissions),
		ActiveOpenDecisions:    len(openDecisions),
		SpecTreeDigest:         graph.SpecTreeDigest,
		Output:                 outputPath,
		Mode:                   mode,
	}
	for _, gap := range report.Gaps {
		summary.ByCode[gap.Code]++
		summary.SubjectsByCode[gap.Code] = append(summary.SubjectsByCode[gap.Code], gap.SubjectID)
	}
	for code := range summary.SubjectsByCode {
		sort.Strings(summary.SubjectsByCode[code])
	}
	return summary, nil
}

func buildCoverageAudit(
	graph ProductGraph,
	capabilities map[string][]string,
	requirements map[string][]string,
	screens map[string]struct{},
	permissions map[string][]string,
	openDecisions map[string][]string,
	specRel string,
) CoverageAuditReport {
	entityByID := map[string]ProductEntity{}
	referenceSources := map[string][]string{}
	for _, entity := range graph.Entities {
		entityByID[entity.ID] = entity
	}
	for _, edge := range graph.Edges {
		referenceSources[edge.To] = append(referenceSources[edge.To], edge.SourcePath)
	}
	normalizeRegistryEvidence(referenceSources)

	var gaps []CoverageGap
	add := func(code, subjectID, subjectKind, detail string, sources []string) {
		sources = normalizedPaths(sources)
		if len(sources) == 0 {
			return
		}
		gaps = append(gaps, CoverageGap{
			Code:        code,
			SubjectID:   subjectID,
			SubjectKind: subjectKind,
			SourcePaths: sources,
			Detail:      detail,
		})
	}

	for id, sources := range capabilities {
		entity, ok := entityByID[id]
		if !ok || entity.ReferenceOnly || entity.Kind != "capability" {
			add("registered-capability-missing-owner", id, "capability", "Active Capability Register entry has no owned canonical capability entity.", sources)
		}
	}
	for _, entity := range graph.Entities {
		if entity.ReferenceOnly || entity.Kind != "capability" {
			continue
		}
		if _, ok := capabilities[entity.ID]; !ok {
			add("owned-capability-unregistered", entity.ID, "capability", "Owned canonical capability is absent from active Capability Register shards.", []string{entity.SourcePath})
		}
	}
	for id, sources := range referenceSources {
		if !capabilityIDPattern.MatchString(id) {
			continue
		}
		if _, ok := capabilities[id]; !ok {
			add("capability-reference-outside-register", id, "capability", "Capability-shaped reference is not an active Capability Register entry.", sources)
		}
	}

	screenRegister := filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers", "screen-register.md"))
	for id := range screens {
		entity, ok := entityByID[id]
		if !ok || entity.ReferenceOnly || entity.Kind != "screen" || entity.RegistryStatus != "registered-active" {
			add("registered-screen-missing-owner", id, "screen", "Active Screen Register entry has no owned canonical screen entity.", []string{screenRegister})
		}
	}
	for _, entity := range graph.Entities {
		if entity.ReferenceOnly || entity.Kind != "screen" {
			continue
		}
		if _, ok := screens[entity.ID]; !ok {
			add("owned-screen-unregistered", entity.ID, "screen", "Owned canonical screen is absent from the active Screen Register.", []string{entity.SourcePath})
		}
	}

	for id, sources := range requirements {
		if len(referenceSources[id]) == 0 {
			add("active-requirement-unreferenced", id, "requirement", "Source Requirement is active in the 122-ID baseline but is never referenced by the active Product Spec corpus.", sources)
		}
	}
	for id, sources := range referenceSources {
		if !requirementIDPattern.MatchString(id) {
			continue
		}
		if _, ok := requirements[id]; !ok {
			add("requirement-reference-outside-catalog", id, "requirement", "Requirement-shaped reference is outside the canonical 122-ID source Requirement Catalog.", sources)
		}
	}

	for id, sources := range permissions {
		if len(referenceSources[id]) == 0 {
			add("registered-permission-unreferenced", id, "permission", "Registered permission is never referenced by the active Product Spec corpus.", sources)
		}
	}
	for id, sources := range referenceSources {
		if !strings.HasPrefix(id, "perm.") {
			continue
		}
		if _, ok := permissions[id]; !ok {
			add("permission-reference-outside-register", id, "permission", "Permission reference is absent from the governance Permission Register.", sources)
		}
	}

	for id, sources := range openDecisions {
		if len(referenceSources[id]) == 0 {
			add("active-open-decision-unreferenced", id, "open-decision", "Active OPEN decision is never referenced by the active Product Spec corpus.", sources)
		}
	}
	for id, sources := range referenceSources {
		if !openDecisionIDPattern.MatchString(id) {
			continue
		}
		if _, ok := openDecisions[id]; !ok {
			add("open-decision-reference-outside-active-register", id, "open-decision", "OPEN-shaped reference is not one of the currently active open decisions.", sources)
		}
	}

	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].Code != gaps[j].Code {
			return gaps[i].Code < gaps[j].Code
		}
		return gaps[i].SubjectID < gaps[j].SubjectID
	})
	return CoverageAuditReport{
		SchemaVersion:  1,
		SpecTreeDigest: graph.SpecTreeDigest,
		Gaps:           gaps,
	}
}

func normalizedPaths(paths []string) []string {
	out := append([]string(nil), paths...)
	sort.Strings(out)
	unique := out[:0]
	var previous string
	for _, path := range out {
		if path == "" || path == previous {
			continue
		}
		unique = append(unique, path)
		previous = path
	}
	return unique
}
