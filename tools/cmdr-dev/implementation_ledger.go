package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const implementationLedgerPath = "engineering/implementation/ledger.json"

type ImplementationLedgerEntry struct {
	Capability      string `json:"capability"`
	WorkUnit        string `json:"work_unit"`
	ManifestPath    string `json:"manifest_path"`
	HandoffPath     string `json:"handoff_path"`
	RuntimeBoundary string `json:"runtime_boundary"`
	RuntimeRoot     string `json:"runtime_root"`
}

type ImplementationLedger struct {
	SchemaVersion       int                         `json:"schema_version"`
	LedgerKind          string                      `json:"ledger_kind"`
	ProductSpecBaseline string                      `json:"product_spec_baseline"`
	Entries             []ImplementationLedgerEntry `json:"entries"`
}

type RuntimeImplementationHandoff struct {
	SchemaVersion        int    `json:"schema_version"`
	WorkUnit             string `json:"work_unit"`
	Result               string `json:"result"`
	ImplementationCommit string `json:"implementation_commit"`
	FinalValidatedHead   string `json:"final_validated_head"`
	PullRequest          int    `json:"pull_request"`
	Evidence             struct {
		PushRun               int     `json:"push_run"`
		PullRequestRun        int     `json:"pull_request_run"`
		Contract              string  `json:"contract"`
		Fixtures              int     `json:"fixtures"`
		NegativeFixtures      int     `json:"negative_fixtures"`
		RuntimeFiles          int     `json:"runtime_files"`
		RuntimeDependencies   int     `json:"runtime_dependencies"`
		CrossBoundaryEdges    int     `json:"cross_boundary_edges"`
		GlobalCoveragePercent float64 `json:"global_coverage_percent"`
		SASTFindings          int     `json:"sast_findings"`
		SCAActionableFindings int     `json:"sca_actionable_findings"`
	} `json:"evidence"`
	Invariants         []string `json:"invariants"`
	ProductSpecMutated bool     `json:"product_spec_mutated"`
	NextUnlocked       []string `json:"next_unlocked"`
}

type ImplementationLedgerClaim struct {
	Capability           string  `json:"capability"`
	WorkUnit             string  `json:"work_unit"`
	RuntimeBoundary      string  `json:"runtime_boundary"`
	RuntimeRoot          string  `json:"runtime_root"`
	ImplementationCommit string  `json:"implementation_commit"`
	FinalValidatedHead   string  `json:"final_validated_head"`
	CoveragePercent      float64 `json:"coverage_percent"`
	SASTFindings         int     `json:"sast_findings"`
	SCAActionable        int     `json:"sca_actionable_findings"`
}

type ImplementationLedgerAuditSummary struct {
	Entries             int                         `json:"entries"`
	VerifiedClaims      int                         `json:"verified_claims"`
	Capabilities        int                         `json:"capabilities"`
	RuntimeBoundaries   int                         `json:"runtime_boundaries"`
	ProductSpecBaseline string                      `json:"product_spec_baseline"`
	SpecTreeDigest      string                      `json:"spec_tree_digest"`
	Claims              []ImplementationLedgerClaim `json:"claims"`
	Status              string                      `json:"status"`
}

func runImplementationLedgerAudit(root string, state CurrentState, graph WorkGraph) (ImplementationLedgerAuditSummary, error) {
	var ledger ImplementationLedger
	if err := decodeStrict(root, implementationLedgerPath, &ledger); err != nil {
		return ImplementationLedgerAuditSummary{}, err
	}
	if err := validateImplementationLedgerShape(ledger, state); err != nil {
		return ImplementationLedgerAuditSummary{}, err
	}

	obligations, err := buildCurrentObligationSet(root, state)
	if err != nil {
		return ImplementationLedgerAuditSummary{}, err
	}
	activeCapabilities := map[string]bool{}
	for _, obligation := range obligations.Obligations {
		if obligation.Family == "capability" {
			activeCapabilities[obligation.SubjectID] = true
		}
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return ImplementationLedgerAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return ImplementationLedgerAuditSummary{}, err
	}
	boundaries := map[string]ArchitectureBoundary{}
	for _, boundary := range registry.Boundaries {
		boundaries[boundary.ID] = boundary
	}

	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}

	summary := ImplementationLedgerAuditSummary{
		Entries:             len(ledger.Entries),
		ProductSpecBaseline: ledger.ProductSpecBaseline,
		SpecTreeDigest:      obligations.SpecTreeDigest,
		Status:              "PASS",
	}
	capabilities := map[string]bool{}
	runtimeBoundaries := map[string]bool{}

	for _, entry := range ledger.Entries {
		if !activeCapabilities[entry.Capability] {
			return summary, fmt.Errorf("implementation claim capability %s is not an active capability obligation", entry.Capability)
		}
		node, ok := nodes[entry.WorkUnit]
		if !ok || node.Status != "VERIFIED" {
			return summary, fmt.Errorf("implementation claim %s work unit %s is not VERIFIED", entry.Capability, entry.WorkUnit)
		}
		manifest, err := decodeWorkManifestV2(root, entry.ManifestPath)
		if err != nil {
			return summary, err
		}
		if manifest.ID != entry.WorkUnit || !containsString(manifest.ProductRefs.Capabilities, entry.Capability) {
			return summary, fmt.Errorf("implementation claim %s is not owned by manifest %s", entry.Capability, entry.ManifestPath)
		}

		var handoff RuntimeImplementationHandoff
		if err := decodeStrict(root, entry.HandoffPath, &handoff); err != nil {
			return summary, err
		}
		if handoff.WorkUnit != entry.WorkUnit || handoff.Result != "VERIFIED" || handoff.ProductSpecMutated {
			return summary, fmt.Errorf("implementation handoff %s is not valid VERIFIED evidence", entry.HandoffPath)
		}
		implementationCommit, err := validateFullCommitID(handoff.ImplementationCommit)
		if err != nil || implementationCommit != handoff.ImplementationCommit {
			return summary, fmt.Errorf("implementation handoff %s has invalid implementation commit", entry.HandoffPath)
		}
		finalHead, err := validateFullCommitID(handoff.FinalValidatedHead)
		if err != nil || finalHead != handoff.FinalValidatedHead {
			return summary, fmt.Errorf("implementation handoff %s has invalid final validated head", entry.HandoffPath)
		}
		if handoff.Evidence.PushRun < 1 || handoff.Evidence.PullRequestRun < 1 ||
			handoff.Evidence.RuntimeFiles < 1 || handoff.Evidence.Fixtures < 1 ||
			handoff.Evidence.NegativeFixtures < 1 ||
			handoff.Evidence.RuntimeDependencies != 0 ||
			handoff.Evidence.CrossBoundaryEdges != 0 ||
			handoff.Evidence.GlobalCoveragePercent < 80 ||
			handoff.Evidence.SASTFindings != 0 ||
			handoff.Evidence.SCAActionableFindings != 0 {
			return summary, fmt.Errorf("implementation handoff %s does not meet verified runtime evidence floors", entry.HandoffPath)
		}

		boundary, ok := boundaries[entry.RuntimeBoundary]
		if !ok || boundary.Kind != "product-runtime" || !boundary.MutableByImplementation {
			return summary, fmt.Errorf("implementation claim %s runtime boundary %s is not an implementation-owned product-runtime boundary", entry.Capability, entry.RuntimeBoundary)
		}
		if len(boundary.Roots) != 1 || boundary.Roots[0] != entry.RuntimeRoot {
			return summary, fmt.Errorf("implementation claim %s runtime root mismatch for boundary %s", entry.Capability, entry.RuntimeBoundary)
		}
		baseRoot := strings.TrimSuffix(entry.RuntimeRoot, "/**")
		info, err := statRepoPath(root, filepath.FromSlash(baseRoot))
		if err != nil || !info.IsDir() {
			return summary, fmt.Errorf("implementation claim %s runtime root %s does not exist", entry.Capability, baseRoot)
		}

		capabilities[entry.Capability] = true
		runtimeBoundaries[entry.RuntimeBoundary] = true
		summary.Claims = append(summary.Claims, ImplementationLedgerClaim{
			Capability: entry.Capability, WorkUnit: entry.WorkUnit,
			RuntimeBoundary: entry.RuntimeBoundary, RuntimeRoot: entry.RuntimeRoot,
			ImplementationCommit: handoff.ImplementationCommit,
			FinalValidatedHead:   handoff.FinalValidatedHead,
			CoveragePercent:      handoff.Evidence.GlobalCoveragePercent,
			SASTFindings:         handoff.Evidence.SASTFindings,
			SCAActionable:        handoff.Evidence.SCAActionableFindings,
		})
	}
	sort.Slice(summary.Claims, func(i, j int) bool { return summary.Claims[i].Capability < summary.Claims[j].Capability })
	summary.VerifiedClaims = len(summary.Claims)
	summary.Capabilities = len(capabilities)
	summary.RuntimeBoundaries = len(runtimeBoundaries)
	if summary.VerifiedClaims != summary.Entries {
		return summary, fmt.Errorf("implementation ledger verification accounting mismatch")
	}
	return summary, nil
}

func buildCurrentObligationSet(root string, state CurrentState) (ObligationSet, error) {
	inventory, err := buildSpecInventory(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ObligationSet{}, err
	}
	screens, err := loadActiveScreenIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ObligationSet{}, err
	}
	productGraph, err := buildProductGraph(inventory, screens)
	if err != nil {
		return ObligationSet{}, err
	}
	capabilities, err := loadActiveCapabilityIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ObligationSet{}, err
	}
	requirements, err := loadActiveRequirementIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ObligationSet{}, err
	}
	permissions, err := loadActivePermissionIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ObligationSet{}, err
	}
	return buildObligationSet(productGraph, capabilities, requirements, screens, permissions)
}

func validateImplementationLedgerShape(ledger ImplementationLedger, state CurrentState) error {
	if ledger.SchemaVersion != 1 {
		return fmt.Errorf("unsupported implementation ledger schema_version %d", ledger.SchemaVersion)
	}
	if ledger.LedgerKind != "verified-runtime-implementation-evidence" {
		return fmt.Errorf("unexpected implementation ledger kind %q", ledger.LedgerKind)
	}
	if ledger.ProductSpecBaseline != state.ProductSpec.BaselineCommit {
		return fmt.Errorf("implementation ledger Product Spec baseline mismatch")
	}
	if len(ledger.Entries) == 0 {
		return fmt.Errorf("implementation ledger must contain at least one verified claim")
	}
	seenCapabilities := map[string]bool{}
	seenWorkUnits := map[string]bool{}
	seenBoundaries := map[string]bool{}
	for i, entry := range ledger.Entries {
		if !capabilityIDPattern.MatchString(entry.Capability) {
			return fmt.Errorf("implementation ledger entry %d has invalid capability %q", i, entry.Capability)
		}
		if strings.TrimSpace(entry.WorkUnit) == "" || strings.TrimSpace(entry.RuntimeBoundary) == "" {
			return fmt.Errorf("implementation ledger entry %d requires work unit and runtime boundary", i)
		}
		if seenCapabilities[entry.Capability] {
			return fmt.Errorf("duplicate implementation capability %s", entry.Capability)
		}
		if seenWorkUnits[entry.WorkUnit] {
			return fmt.Errorf("duplicate implementation work unit %s", entry.WorkUnit)
		}
		if seenBoundaries[entry.RuntimeBoundary] {
			return fmt.Errorf("duplicate implementation runtime boundary %s", entry.RuntimeBoundary)
		}
		seenCapabilities[entry.Capability] = true
		seenWorkUnits[entry.WorkUnit] = true
		seenBoundaries[entry.RuntimeBoundary] = true
		for label, path := range map[string]string{"manifest": entry.ManifestPath, "handoff": entry.HandoffPath} {
			normalized, err := normalizeChangedPaths([]string{path})
			if err != nil || len(normalized) != 1 || normalized[0] != path || !strings.HasPrefix(path, "work/lots/") {
				return fmt.Errorf("implementation ledger entry %d has invalid %s path %q", i, label, path)
			}
		}
		if !strings.HasPrefix(entry.RuntimeRoot, "product-runtime/") || !strings.HasSuffix(entry.RuntimeRoot, "/**") {
			return fmt.Errorf("implementation ledger entry %d has invalid runtime root %q", i, entry.RuntimeRoot)
		}
	}
	return nil
}
