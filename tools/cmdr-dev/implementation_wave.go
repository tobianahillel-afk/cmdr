package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const implementationWavePolicy = "multi-dimensional-v1"

type WaveCandidateScore struct {
	Capability      string  `json:"capability"`
	ReuseScore      float64 `json:"reuse_score"`
	PriorityScore   float64 `json:"priority_score"`
	SecurityScore   float64 `json:"security_score"`
	RiskScore       float64 `json:"risk_score"`
	DependencyScore float64 `json:"dependency_score"`
	TotalScore      float64 `json:"total_score"`
}

type ImplementationWaveSelection struct {
	SchemaVersion       int                 `json:"schema_version"`
	ProgramKind         string              `json:"program_kind"`
	SelectionPolicy     string              `json:"selection_policy"`
	ProductSpecBaseline string              `json:"product_spec_baseline"`
	SpecTreeDigest      string              `json:"spec_tree_digest"`
	ReadyCandidates     int                 `json:"ready_candidates"`
	PreferredFamily     string              `json:"preferred_family,omitempty"`
	Capability          string              `json:"capability"`
	Name                string              `json:"name"`
	RegisterPath        string              `json:"register_path"`
	CanonicalFile       string              `json:"canonical_file"`
	DeliveryStatus      string              `json:"delivery_status"`
	DeliveryMode        string              `json:"delivery_mode"`
	CandidateScore      *WaveCandidateScore `json:"candidate_score,omitempty"`
	Status              string              `json:"status"`
}

var securityKeywordPattern = regexp.MustCompile(`(?i)(security|tenant|isolation|secret|permission|authorization|auth|credential|integrity|audit|gate)`)

func evaluateWaveCandidateScore(record ImplementationReadinessRecord, implementedFamilies map[string]bool) WaveCandidateScore {
	family, _ := capabilityFamily(record.Capability)
	score := WaveCandidateScore{
		Capability: record.Capability,
	}

	// 1. REUSE DIMENSION (0 - 40 points)
	// Maximizes co-location with existing implemented modules/families
	if implementedFamilies[family] {
		score.ReuseScore = 40.0
	}

	// 2. PRIORITY / PRODUCT VALUE DIMENSION (0 - 30 points)
	// Strategic module hierarchy based on platform foundations and delivery status
	var basePriority float64
	switch family {
	case "CAP-STD":
		basePriority = 25.0 // Studio tooling and contracts
	case "CAP-CMD":
		basePriority = 22.0 // Command and orchestrator center
	case "CAP-SET":
		basePriority = 20.0 // Platform configuration & tenant envelope
	case "CAP-INV":
		basePriority = 18.0 // Forensic investigation & event search
	case "CAP-GOV":
		basePriority = 16.0 // Policy & governance rules
	case "CAP-EPT":
		basePriority = 14.0 // Endpoint telemetry & agent
	default:
		basePriority = 10.0
	}
	if strings.EqualFold(record.DeliveryStatus, "defined") || strings.EqualFold(record.DeliveryStatus, "ready") {
		basePriority += 5.0
	}
	score.PriorityScore = basePriority

	// 3. SECURITY DIMENSION (0 - 20 points)
	// Prioritizes tenant-isolation, secrets protection, and authorization mechanisms
	if securityKeywordPattern.MatchString(record.Capability) ||
		securityKeywordPattern.MatchString(record.Name) ||
		family == "CAP-SET" || family == "CAP-GOV" {
		score.SecurityScore = 20.0
	} else {
		score.SecurityScore = 10.0
	}

	// 4. RISK DIMENSION (0 - 15 points)
	// Lower risk gets higher score: clean specification anchor and no blockers
	var riskPoints float64
	if strings.TrimSpace(record.CanonicalFile) != "" {
		riskPoints += 10.0
	}
	if len(record.Blockers) == 0 {
		riskPoints += 5.0
	}
	score.RiskScore = riskPoints

	// 5. DEPENDENCY POSTURE (0 - 15 points)
	// Foundation capabilities (lower numeric sequence) and unblocked progression
	var depPoints float64 = 10.0
	suffixIdx := strings.LastIndex(record.Capability, "-")
	if suffixIdx > 0 && suffixIdx+1 < len(record.Capability) {
		numPart := strings.TrimLeft(record.Capability[suffixIdx+1:], "0")
		if num, err := strconv.Atoi(numPart); err == nil && num <= 10 {
			depPoints += 5.0 // Foundational sequence within family
		}
	}
	score.DependencyScore = depPoints

	score.TotalScore = score.ReuseScore + score.PriorityScore + score.SecurityScore + score.RiskScore + score.DependencyScore
	return score
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

	type scoredCandidate struct {
		record ImplementationReadinessRecord
		score  WaveCandidateScore
	}
	scored := make([]scoredCandidate, len(ready))
	for i, r := range ready {
		scored[i] = scoredCandidate{
			record: r,
			score:  evaluateWaveCandidateScore(r, implementedFamilies),
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		// 1. Primary order: multi-dimensional total score (descending)
		if scored[i].score.TotalScore != scored[j].score.TotalScore {
			return scored[i].score.TotalScore > scored[j].score.TotalScore
		}
		// 2. Secondary order: preferred implemented family (reuse)
		leftFamily, leftErr := capabilityFamily(scored[i].record.Capability)
		rightFamily, rightErr := capabilityFamily(scored[j].record.Capability)
		leftPreferred := leftErr == nil && implementedFamilies[leftFamily]
		rightPreferred := rightErr == nil && implementedFamilies[rightFamily]
		if leftPreferred != rightPreferred {
			return leftPreferred
		}
		// 3. Stable alphabetical tie-breaking
		return scored[i].record.Capability < scored[j].record.Capability
	})

	selectedScored := scored[0]
	selected := selectedScored.record
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
	selection.CandidateScore = &selectedScored.score
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
