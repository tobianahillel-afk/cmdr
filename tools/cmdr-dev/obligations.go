package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type ImplementationObligation struct {
	ID          string   `json:"id"`
	Family      string   `json:"family"`
	SubjectID   string   `json:"subject_id"`
	SourcePaths []string `json:"source_paths"`
}

type UnresolvedProductReference struct {
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	SourcePaths []string `json:"source_paths"`
}

type ObligationSet struct {
	SchemaVersion        int                          `json:"schema_version"`
	SpecTreeDigest       string                       `json:"spec_tree_digest"`
	Obligations          []ImplementationObligation   `json:"obligations"`
	UnresolvedReferences []UnresolvedProductReference `json:"unresolved_references"`
}

type ObligationSummary struct {
	Obligations            int            `json:"obligations"`
	ByFamily               map[string]int `json:"by_family"`
	UnresolvedReferences              int      `json:"unresolved_references"`
	UnregisteredRequirementReferences []string `json:"unregistered_requirement_references,omitempty"`
	RegisteredCapabilities            int      `json:"registered_capabilities"`
	RegisteredRequirements int            `json:"registered_requirements"`
	RegisteredScreens      int            `json:"registered_screens"`
	RegisteredPermissions  int            `json:"registered_permissions"`
	SpecTreeDigest         string         `json:"spec_tree_digest"`
	Output                 string         `json:"output"`
	Mode                   string         `json:"mode"`
}

func runObligations(root, specRel, output string, check bool) (ObligationSummary, error) {
	inventory, err := buildSpecInventory(root, specRel)
	if err != nil {
		return ObligationSummary{}, err
	}
	screens, err := loadActiveScreenIDs(root, specRel)
	if err != nil {
		return ObligationSummary{}, err
	}
	graph, err := buildProductGraph(inventory, screens)
	if err != nil {
		return ObligationSummary{}, err
	}
	capabilities, err := loadActiveCapabilityIDs(root, specRel)
	if err != nil {
		return ObligationSummary{}, err
	}
	requirements, err := loadActiveRequirementIDs(root, specRel)
	if err != nil {
		return ObligationSummary{}, err
	}
	permissions, err := loadActivePermissionIDs(root, specRel)
	if err != nil {
		return ObligationSummary{}, err
	}
	set, err := buildObligationSet(graph, capabilities, requirements, screens, permissions)
	if err != nil {
		return ObligationSummary{}, err
	}

	data, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return ObligationSummary{}, err
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
			return ObligationSummary{}, fmt.Errorf("read obligations: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return ObligationSummary{}, fmt.Errorf("obligations are stale: %s", outputPath)
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
			return ObligationSummary{}, err
		}
		if err := os.WriteFile(outputPath, data, 0o644); err != nil {
			return ObligationSummary{}, err
		}
	}

	summary := summarizeObligations(set)
	summary.RegisteredCapabilities = len(capabilities)
	summary.RegisteredRequirements = len(requirements)
	summary.RegisteredScreens = len(screens)
	summary.RegisteredPermissions = len(permissions)
	summary.Output = outputPath
	summary.Mode = mode
	return summary, nil
}

func buildObligationSet(
	graph ProductGraph,
	capabilities map[string][]string,
	requirements map[string][]string,
	screens map[string]struct{},
	permissions map[string][]string,
) (ObligationSet, error) {
	entityByID := make(map[string]ProductEntity, len(graph.Entities))
	referenceSources := map[string][]string{}
	for _, entity := range graph.Entities {
		entityByID[entity.ID] = entity
	}
	for _, edge := range graph.Edges {
		referenceSources[edge.To] = append(referenceSources[edge.To], edge.SourcePath)
	}
	normalizeRegistryEvidence(referenceSources)

	var obligations []ImplementationObligation
	for id, registrySources := range capabilities {
		entity, ok := entityByID[id]
		if !ok || entity.ReferenceOnly || entity.Kind != "capability" {
			continue
		}
		obligations = append(obligations, newObligation("capability", id, append([]string{entity.SourcePath}, registrySources...)))
	}
	for id, sources := range requirements {
		obligations = append(obligations, newObligation("requirement", id, sources))
	}
	for id := range screens {
		entity, ok := entityByID[id]
		if !ok || entity.ReferenceOnly || entity.Kind != "screen" || entity.RegistryStatus != "registered-active" {
			continue
		}
		obligations = append(obligations, newObligation("screen", id, []string{entity.SourcePath, filepath.ToSlash(filepath.Join("cmdr-product-spec", "00-governance", "registers", "screen-register.md"))}))
	}
	for id, sources := range permissions {
		obligations = append(obligations, newObligation("permission", id, sources))
	}
	for _, entity := range graph.Entities {
		if entity.ReferenceOnly {
			continue
		}
		switch entity.Kind {
		case "implementation-contract", "canonical-object":
			obligations = append(obligations, newObligation(entity.Kind, entity.ID, []string{entity.SourcePath}))
		}
	}

	sort.Slice(obligations, func(i, j int) bool {
		if obligations[i].Family != obligations[j].Family {
			return obligations[i].Family < obligations[j].Family
		}
		return obligations[i].SubjectID < obligations[j].SubjectID
	})
	if err := validateObligations(obligations); err != nil {
		return ObligationSet{}, err
	}

	active := map[string]struct{}{}
	for id := range capabilities {
		active[id] = struct{}{}
	}
	for id := range requirements {
		active[id] = struct{}{}
	}
	for id := range screens {
		active[id] = struct{}{}
	}
	for id := range permissions {
		active[id] = struct{}{}
	}

	var unresolved []UnresolvedProductReference
	for _, entity := range graph.Entities {
		if !entity.ReferenceOnly {
			continue
		}
		if _, ok := active[entity.ID]; ok {
			continue
		}
		sources := append([]string(nil), referenceSources[entity.ID]...)
		if len(sources) == 0 {
			continue
		}
		unresolved = append(unresolved, UnresolvedProductReference{
			ID:          entity.ID,
			Kind:        entity.Kind,
			SourcePaths: sources,
		})
	}
	sort.Slice(unresolved, func(i, j int) bool {
		if unresolved[i].ID != unresolved[j].ID {
			return unresolved[i].ID < unresolved[j].ID
		}
		return unresolved[i].Kind < unresolved[j].Kind
	})

	return ObligationSet{
		SchemaVersion:        1,
		SpecTreeDigest:       graph.SpecTreeDigest,
		Obligations:          obligations,
		UnresolvedReferences: unresolved,
	}, nil
}

func newObligation(family, subject string, sources []string) ImplementationObligation {
	sources = append([]string(nil), sources...)
	sort.Strings(sources)
	unique := sources[:0]
	var previous string
	for i, source := range sources {
		if source == "" {
			continue
		}
		if i == 0 || source != previous {
			unique = append(unique, source)
			previous = source
		}
	}
	return ImplementationObligation{
		ID:          "OBL:" + family + ":" + subject,
		Family:      family,
		SubjectID:   subject,
		SourcePaths: unique,
	}
}

func validateObligations(obligations []ImplementationObligation) error {
	seen := map[string]struct{}{}
	for _, obligation := range obligations {
		if obligation.ID == "" || obligation.Family == "" || obligation.SubjectID == "" {
			return fmt.Errorf("invalid empty obligation identity")
		}
		if len(obligation.SourcePaths) == 0 {
			return fmt.Errorf("obligation %s has no source provenance", obligation.ID)
		}
		if _, ok := seen[obligation.ID]; ok {
			return fmt.Errorf("duplicate obligation %s", obligation.ID)
		}
		seen[obligation.ID] = struct{}{}
	}
	return nil
}

func summarizeObligations(set ObligationSet) ObligationSummary {
	summary := ObligationSummary{
		Obligations:          len(set.Obligations),
		ByFamily:             map[string]int{},
		UnresolvedReferences: len(set.UnresolvedReferences),
		SpecTreeDigest:       set.SpecTreeDigest,
	}
	for _, obligation := range set.Obligations {
		summary.ByFamily[obligation.Family]++
	}
	for _, unresolved := range set.UnresolvedReferences {
		if unresolved.Kind == "requirement" {
			summary.UnregisteredRequirementReferences = append(summary.UnregisteredRequirementReferences, unresolved.ID)
		}
	}
	sort.Strings(summary.UnregisteredRequirementReferences)
	return summary
}
