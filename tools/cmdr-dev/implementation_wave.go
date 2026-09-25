package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const implementationWavePolicy = "reuse-first-v1"

type ImplementationWaveSelection struct {
	SchemaVersion       int    `json:"schema_version"`
	ProgramKind         string `json:"program_kind"`
	SelectionPolicy     string `json:"selection_policy"`
	ProductSpecBaseline string `json:"product_spec_baseline"`
	SpecTreeDigest      string `json:"spec_tree_digest"`
	ReadyCandidates     int    `json:"ready_candidates"`
	PreferredFamily     string `json:"preferred_family,omitempty"`
	Capability          string `json:"capability"`
	Name                string `json:"name"`
	RegisterPath        string `json:"register_path"`
	CanonicalFile       string `json:"canonical_file"`
	DeliveryStatus      string `json:"delivery_status"`
	DeliveryMode        string `json:"delivery_mode"`
	Status              string `json:"status"`
}

func runImplementationWave(root string, state CurrentState, graph WorkGraph) (ImplementationWaveSelection, error) {
	program, err := runImplementationReadiness(root, state, graph)
	if err != nil {
		return ImplementationWaveSelection{}, err
	}
	selection, err := selectImplementationWave(program)
	if err != nil {
		return selection, err
	}
	if selection.ProductSpecBaseline != state.ProductSpec.BaselineCommit {
		return selection, fmt.Errorf("implementation wave Product Spec baseline mismatch")
	}
	inventory, err := buildSpecInventory(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return selection, err
	}
	if inventory.TreeDigest != selection.SpecTreeDigest {
		return selection, fmt.Errorf("implementation wave Product Spec digest changed during selection")
	}
	canonicalFile, err := resolveWaveCanonicalFile(selection, inventory)
	if err != nil {
		return selection, err
	}
	selection.CanonicalFile = canonicalFile
	canonicalRel, err := normalizeCanonicalCapabilityPath(state.ProductSpec.CanonicalPath, selection.CanonicalFile)
	if err != nil {
		return selection, fmt.Errorf("selected capability %s canonical file: %w", selection.Capability, err)
	}
	data, err := readRepoFile(root, canonicalRel)
	if err != nil {
		return selection, fmt.Errorf("selected capability %s canonical file %s is unavailable: %w", selection.Capability, canonicalRel, err)
	}
	if len(data) == 0 {
		return selection, fmt.Errorf("selected capability %s canonical file %s is empty", selection.Capability, canonicalRel)
	}
	selection.Status = "PASS"
	return selection, nil
}

func selectImplementationWave(program ImplementationReadinessProgram) (ImplementationWaveSelection, error) {
	selection := ImplementationWaveSelection{
		SchemaVersion:       1,
		ProgramKind:         "bounded-implementation-wave",
		SelectionPolicy:     implementationWavePolicy,
		ProductSpecBaseline: program.ProductSpecBaseline,
		SpecTreeDigest:      program.SpecTreeDigest,
	}
	if program.Status != "PASS" {
		return selection, fmt.Errorf("implementation readiness program is not PASS")
	}

	implementedFamilies := map[string]bool{}
	var ready []ImplementationReadinessRecord
	for _, record := range program.Records {
		switch record.State {
		case "IMPLEMENTED":
			family, err := capabilityFamily(record.Capability)
			if err != nil {
				return selection, err
			}
			implementedFamilies[family] = true
		case "READY":
			ready = append(ready, record)
		}
	}
	selection.ReadyCandidates = len(ready)
	if len(ready) == 0 {
		return selection, fmt.Errorf("implementation readiness contains no READY capability")
	}

	sort.Slice(ready, func(i, j int) bool {
		leftFamily, leftErr := capabilityFamily(ready[i].Capability)
		rightFamily, rightErr := capabilityFamily(ready[j].Capability)
		leftPreferred := leftErr == nil && implementedFamilies[leftFamily]
		rightPreferred := rightErr == nil && implementedFamilies[rightFamily]
		if leftPreferred != rightPreferred {
			return leftPreferred
		}
		return ready[i].Capability < ready[j].Capability
	})

	selected := ready[0]
	family, err := capabilityFamily(selected.Capability)
	if err != nil {
		return selection, err
	}
	if implementedFamilies[family] {
		selection.PreferredFamily = family
	}
	selection.Capability = selected.Capability
	selection.Name = selected.Name
	selection.RegisterPath = selected.RegisterPath
	selection.CanonicalFile = selected.CanonicalFile
	selection.DeliveryStatus = selected.DeliveryStatus
	selection.DeliveryMode = selected.DeliveryMode
	return selection, nil
}

func resolveWaveCanonicalFile(selection ImplementationWaveSelection, inventory SpecInventory) (string, error) {
	if strings.TrimSpace(selection.CanonicalFile) != "" {
		return selection.CanonicalFile, nil
	}
	var matches []string
	for _, doc := range inventory.Documents {
		if doc.ID == selection.Capability && doc.Active && doc.Canonical {
			path := filepath.ToSlash(doc.Path)
			prefix := strings.TrimSuffix(filepath.ToSlash(inventory.SpecRoot), "/") + "/"
			if !strings.HasPrefix(path, prefix) {
				return "", fmt.Errorf("canonical document %s for %s is outside Product Spec root", path, selection.Capability)
			}
			matches = append(matches, strings.TrimPrefix(path, prefix))
		}
	}
	sort.Strings(matches)
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("selected READY capability %s has no canonical file in its register and no active canonical Product Spec document", selection.Capability)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("selected READY capability %s resolves to multiple active canonical Product Spec documents: %v", selection.Capability, matches)
	}
}

func capabilityFamily(id string) (string, error) {
	if !capabilityIDPattern.MatchString(id) {
		return "", fmt.Errorf("invalid capability id %q", id)
	}
	lastDash := strings.LastIndex(id, "-")
	if lastDash <= 0 {
		return "", fmt.Errorf("capability id %q has no family", id)
	}
	return id[:lastDash], nil
}

func normalizeCanonicalCapabilityPath(specRel, canonicalFile string) (string, error) {
	canonicalFile = filepath.ToSlash(strings.Trim(strings.TrimSpace(canonicalFile), "`"))
	if canonicalFile == "" {
		return "", fmt.Errorf("canonical file is empty")
	}
	if strings.HasPrefix(canonicalFile, filepath.ToSlash(specRel)+"/") {
		values, err := normalizeChangedPaths([]string{canonicalFile})
		if err != nil || len(values) != 1 {
			return "", fmt.Errorf("invalid canonical path %q", canonicalFile)
		}
		return values[0], nil
	}
	joined := filepath.ToSlash(filepath.Join(specRel, filepath.FromSlash(canonicalFile)))
	values, err := normalizeChangedPaths([]string{joined})
	if err != nil || len(values) != 1 {
		return "", fmt.Errorf("invalid canonical path %q", canonicalFile)
	}
	if !strings.HasPrefix(values[0], filepath.ToSlash(specRel)+"/") {
		return "", fmt.Errorf("canonical path escapes Product Spec root: %q", canonicalFile)
	}
	return values[0], nil
}
