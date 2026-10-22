package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	Capability             string  `json:"capability"`
	WorkUnit               string  `json:"work_unit"`
	RuntimeBoundary        string  `json:"runtime_boundary"`
	RuntimeRoot            string  `json:"runtime_root"`
	ImplementationCommit   string  `json:"implementation_commit"`
	FinalValidatedHead     string  `json:"final_validated_head"`
	SuccessorWorkUnit      string  `json:"successor_work_unit,omitempty"`
	SuccessorValidatedHead string  `json:"successor_validated_head,omitempty"`
	CoveragePercent        float64 `json:"coverage_percent"`
	SASTFindings           int     `json:"sast_findings"`
	SCAActionable          int     `json:"sca_actionable_findings"`
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
		successorWorkUnit, successorHead, err := resolveRuntimeSuccessorEvidence(root, graph, entry, manifest, finalHead)
		if err != nil {
			return summary, fmt.Errorf("implementation handoff %s successor evidence: %w", entry.HandoffPath, err)
		}
		validatedRootHead := finalHead
		if successorHead != "" {
			validatedRootHead = successorHead
		}
		if err := verifyRuntimeClaimGitEvidence(root, implementationCommit, validatedRootHead, entry.RuntimeRoot); err != nil {
			return summary, fmt.Errorf("implementation handoff %s has stale or invalid Git evidence: %w", entry.HandoffPath, err)
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
			ImplementationCommit:   handoff.ImplementationCommit,
			FinalValidatedHead:     handoff.FinalValidatedHead,
			SuccessorWorkUnit:      successorWorkUnit,
			SuccessorValidatedHead: successorHead,
			CoveragePercent:        handoff.Evidence.GlobalCoveragePercent,
			SASTFindings:           handoff.Evidence.SASTFindings,
			SCAActionable:          handoff.Evidence.SCAActionableFindings,
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

// verifyRuntimeClaimGitEvidence rejects syntactically plausible SHA claims
// unless both commits exist, are on the current ancestry, and the recorded
// product-runtime boundary has not changed since the final validated head.
// An old successful run is not evidence for newly added runtime files.
func verifyRuntimeClaimGitEvidence(root, implementationCommit, finalValidatedHead, runtimeRoot string) error {
	baseRoot := strings.TrimSuffix(runtimeRoot, "/**")
	paths, err := normalizeChangedPaths([]string{baseRoot})
	if err != nil || len(paths) != 1 || paths[0] != baseRoot ||
		!strings.HasPrefix(baseRoot, "product-runtime/") || !strings.HasSuffix(runtimeRoot, "/**") {
		return fmt.Errorf("invalid runtime claim root %q", runtimeRoot)
	}
	if err := ensureGitAncestor(root, implementationCommit, finalValidatedHead); err != nil {
		return fmt.Errorf("implementation to validation ancestry: %w", err)
	}
	headOut, err := exec.Command("git", "-C", root, "rev-parse", "--verify", "HEAD").Output() // #nosec G204,G702 -- no shell, fixed Git command and repository root.
	if err != nil {
		return fmt.Errorf("resolve current Git HEAD: %w", err)
	}
	currentHead, err := validateFullCommitID(strings.TrimSpace(string(headOut)))
	if err != nil {
		return fmt.Errorf("current Git HEAD: %w", err)
	}
	if err := ensureGitAncestor(root, finalValidatedHead, currentHead); err != nil {
		return fmt.Errorf("validated head to current HEAD ancestry: %w", err)
	}
	cmd := exec.Command("git", "-C", root, "diff", "--no-ext-diff", "--no-renames", "--quiet",
		finalValidatedHead, currentHead, "--", filepath.ToSlash(baseRoot)) // #nosec G204,G702 -- fixed Git flags; both SHAs were validated through ensureGitAncestor and root is repository-confined.
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("runtime root %s changed after validated head %s; fresh runtime evidence is required", baseRoot, finalValidatedHead)
		}
		return fmt.Errorf("compare validated runtime root: %w", err)
	}
	return nil
}

// RuntimeSuccessorEvidence is an explicit additive proof beside the FROZEN
// original handoff. An old validated SHA is NEVER silently moved or waived.
type RuntimeSuccessorEvidence struct {
	SchemaVersion       int      `json:"schema_version"`
	PredecessorWorkUnit string   `json:"predecessor_work_unit"`
	SuccessorWorkUnit   string   `json:"successor_work_unit"`
	ManifestPath        string   `json:"manifest_path"`
	HandoffPath         string   `json:"handoff_path"`
	AddedRuntimePaths   []string `json:"added_runtime_paths"`
}

type RuntimeSuccessorHandoff struct {
	SchemaVersion      int    `json:"schema_version"`
	WorkUnit           string `json:"work_unit"`
	Result             string `json:"result"`
	FinalValidatedHead string `json:"final_validated_head"`
	PullRequest        int    `json:"pull_request"`
	Evidence           struct {
		PushRun                      int     `json:"push_run"`
		PullRequestRun               int     `json:"pull_request_run"`
		PushP95MS                    float64 `json:"push_p95_ms"`
		PullRequestP95MS             float64 `json:"pull_request_p95_ms"`
		AbsoluteBudgetMS             float64 `json:"absolute_budget_ms"`
		RuntimeGlobalCoveragePercent float64 `json:"runtime_global_coverage_percent"`
		SASTFindings                 int     `json:"sast_findings"`
		SCAActionableFindings        int     `json:"sca_actionable_findings"`
		CacheRecords                 int     `json:"cache_records"`
	} `json:"evidence"`
	Rules              []string `json:"rules"`
	ProductSpecMutated bool     `json:"product_spec_mutated"`
	NextUnlocked       []string `json:"next_unlocked"`
}

func resolveRuntimeSuccessorEvidence(root string, graph WorkGraph, entry ImplementationLedgerEntry, original WorkManifestV2, validatedHead string) (string, string, error) {
	sidecar := filepath.ToSlash(filepath.Join(filepath.Dir(entry.ManifestPath), "RUNTIME-EVIDENCE-SUCCESSOR.json"))
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(sidecar))); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("inspect successor sidecar: %w", err)
	}
	var proof RuntimeSuccessorEvidence
	if err := decodeStrict(root, sidecar, &proof); err != nil {
		return "", "", err
	}
	if proof.SchemaVersion != 1 || proof.PredecessorWorkUnit != entry.WorkUnit ||
		proof.SuccessorWorkUnit == "" || proof.SuccessorWorkUnit == entry.WorkUnit {
		return "", "", fmt.Errorf("invalid predecessor/successor for %s", sidecar)
	}
	for _, p := range []string{proof.ManifestPath, proof.HandoffPath} {
		norm, err := normalizeChangedPaths([]string{p})
		if err != nil || len(norm) != 1 || norm[0] != p || !strings.HasPrefix(p, "work/lots/"+proof.SuccessorWorkUnit+"/") {
			return "", "", fmt.Errorf("invalid successor evidence path %q", p)
		}
	}
	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	node, ok := nodes[proof.SuccessorWorkUnit]
	if !ok || node.Status != "VERIFIED" {
		return "", "", fmt.Errorf("successor %s is not VERIFIED", proof.SuccessorWorkUnit)
	}
	stack := append([]string(nil), node.DependsOn...)
	visited := map[string]bool{}
	found := false
	for len(stack) > 0 {
		i := len(stack) - 1
		id := stack[i]
		stack = stack[:i]
		if id == entry.WorkUnit {
			found = true
			break
		}
		if visited[id] {
			continue
		}
		visited[id] = true
		dep, ok := nodes[id]
		if !ok {
			return "", "", fmt.Errorf("unknown successor dependency %s", id)
		}
		stack = append(stack, dep.DependsOn...)
	}
	if !found {
		return "", "", fmt.Errorf("successor lacks recorded runtime dependency")
	}
	manifest, err := decodeWorkManifestV2(root, proof.ManifestPath)
	if err != nil {
		return "", "", err
	}
	if manifest.ID != proof.SuccessorWorkUnit || manifest.Parent != original.Parent ||
		!containsString(manifest.ProductRefs.Capabilities, entry.Capability) ||
		!containsString(manifest.AllowedPaths, entry.RuntimeRoot) {
		return "", "", fmt.Errorf("successor capability, parent or runtime ownership mismatch")
	}
	var handoff RuntimeSuccessorHandoff
	if err := decodeStrict(root, proof.HandoffPath, &handoff); err != nil {
		return "", "", err
	}
	head, err := validateFullCommitID(handoff.FinalValidatedHead)
	if err != nil || head != handoff.FinalValidatedHead || handoff.SchemaVersion != 1 ||
		handoff.WorkUnit != proof.SuccessorWorkUnit || handoff.Result != "VERIFIED" ||
		handoff.ProductSpecMutated || handoff.PullRequest < 1 ||
		handoff.Evidence.PushRun < 1 || handoff.Evidence.PullRequestRun < 1 ||
		handoff.Evidence.RuntimeGlobalCoveragePercent < 80 ||
		handoff.Evidence.AbsoluteBudgetMS <= 0 ||
		handoff.Evidence.PushP95MS <= 0 || handoff.Evidence.PushP95MS > handoff.Evidence.AbsoluteBudgetMS ||
		handoff.Evidence.PullRequestP95MS <= 0 || handoff.Evidence.PullRequestP95MS > handoff.Evidence.AbsoluteBudgetMS ||
		handoff.Evidence.SASTFindings != 0 || handoff.Evidence.SCAActionableFindings != 0 {
		return "", "", fmt.Errorf("successor handoff has insufficient verified run evidence")
	}
	if err := ensureGitAncestor(root, validatedHead, head); err != nil {
		return "", "", fmt.Errorf("successor is not a descendant of original validation: %w", err)
	}
	if err := verifyOnlyRecordedRuntimeAdditions(root, validatedHead, head, entry.RuntimeRoot, proof.AddedRuntimePaths); err != nil {
		return "", "", err
	}
	return proof.SuccessorWorkUnit, head, nil
}

// No changed/deleted original runtime file can be excused by this proof.
// Only precisely listed new paths may be observed between the two Git trees.
func verifyOnlyRecordedRuntimeAdditions(root, before, after, runtimeRoot string, paths []string) error {
	base := strings.TrimSuffix(runtimeRoot, "/**")
	if len(paths) == 0 || !strings.HasPrefix(base, "product-runtime/") {
		return fmt.Errorf("successor additions missing or unbounded")
	}
	expected := map[string]bool{}
	for _, p := range paths {
		norm, err := normalizeChangedPaths([]string{p})
		if err != nil || len(norm) != 1 || norm[0] != p || !strings.HasPrefix(p, base+"/") || expected[p] {
			return fmt.Errorf("invalid or duplicate successor path %q", p)
		}
		expected[p] = true
	}
	out, err := exec.Command("git", "-C", root, "diff", "--no-ext-diff", "--no-renames",
		"--name-status", before, after, "--", filepath.ToSlash(base)).Output() // #nosec G204,G702 -- no shell, bounded Git SHAs and root.
	if err != nil {
		return fmt.Errorf("inspect successor Git diff: %w", err)
	}
	actual := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 || parts[0] != "A" || !expected[parts[1]] || actual[parts[1]] {
			return fmt.Errorf("successor mutated or added an unrecorded runtime path: %q", line)
		}
		actual[parts[1]] = true
	}
	if len(actual) != len(expected) {
		return fmt.Errorf("successor added %d paths; proof listed %d", len(actual), len(expected))
	}
	return nil
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
