package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	runtimeSecurityPolicyPath   = "engineering/security/runtime-security-policy.json"
	runtimeSecurityEvidencePath = "engineering/testing/runtime-security-evidence.json"
)

type RuntimeSecurityScope struct {
	BoundaryID               string   `json:"boundary_id"`
	Owner                    string   `json:"owner"`
	SecurityCriticalPaths    []string `json:"security_critical_paths"`
	AuthorizationRequired    bool     `json:"authorization_required"`
	AuthorizationRationale   string   `json:"authorization_rationale"`
	TenantIsolationRequired  bool     `json:"tenant_isolation_required"`
	TenantIsolationRationale string   `json:"tenant_isolation_rationale"`
}

type RuntimeSecurityPolicy struct {
	SchemaVersion                        int                    `json:"schema_version"`
	DefaultPolicy                        string                 `json:"default_policy"`
	GlobalCoverageFloorPercent           float64                `json:"global_coverage_floor_percent"`
	ChangedSecurityCriticalCoverageFloor float64                `json:"changed_security_critical_coverage_floor_percent"`
	AuthorizationGateID                  string                 `json:"authorization_gate_id"`
	TenantIsolationGateID                string                 `json:"tenant_isolation_gate_id"`
	Scopes                               []RuntimeSecurityScope `json:"scopes"`
}

type RuntimeSecurityEvidenceScope struct {
	BoundaryID                     string   `json:"boundary_id"`
	ImplementationState            string   `json:"implementation_state"`
	ChangedSecurityCritical        bool     `json:"changed_security_critical"`
	ChangedSecurityCriticalPercent *float64 `json:"changed_security_critical_coverage_percent"`
	AuthorizationNegativePassed    *bool    `json:"authorization_negative_passed"`
	TenantIsolationNegativePassed  *bool    `json:"tenant_isolation_negative_passed"`
}

type RuntimeSecurityEvidence struct {
	SchemaVersion         int                            `json:"schema_version"`
	Status                string                         `json:"status"`
	Reason                string                         `json:"reason"`
	SourceCommit          string                         `json:"source_commit"`
	GlobalCoveragePercent *float64                       `json:"global_coverage_percent"`
	Scopes                []RuntimeSecurityEvidenceScope `json:"scopes"`
}

type SecurityTestAuditSummary struct {
	RuntimeBoundaries             int      `json:"runtime_boundaries"`
	RegisteredScopes              int      `json:"registered_scopes"`
	PreimplementationScopes       int      `json:"preimplementation_scopes"`
	ImplementedScopes             int      `json:"implemented_scopes"`
	AuthorizationRequired         int      `json:"authorization_required"`
	TenantIsolationRequired       int      `json:"tenant_isolation_required"`
	ChangedSecurityCriticalScopes int      `json:"changed_security_critical_scopes"`
	CoverageStatus                string   `json:"coverage_status"`
	GlobalCoveragePercent         *float64 `json:"global_coverage_percent,omitempty"`
	GlobalCoverageFloorPercent    float64  `json:"global_coverage_floor_percent"`
	ChangedCoverageFloorPercent   float64  `json:"changed_security_critical_coverage_floor_percent"`
	NotApplicableReason           string   `json:"not_applicable_reason,omitempty"`
}

func runSecurityTestAudit(root, changesFile string) (SecurityTestAuditSummary, error) {
	var architecture ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &architecture); err != nil {
		return SecurityTestAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(architecture); err != nil {
		return SecurityTestAuditSummary{}, err
	}
	gates, err := loadSecurityGateRegistry(root)
	if err != nil {
		return SecurityTestAuditSummary{}, err
	}
	var policy RuntimeSecurityPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(runtimeSecurityPolicyPath)), &policy); err != nil {
		return SecurityTestAuditSummary{}, err
	}
	var evidence RuntimeSecurityEvidence
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(runtimeSecurityEvidencePath)), &evidence); err != nil {
		return SecurityTestAuditSummary{}, err
	}
	var changed []string
	if strings.TrimSpace(changesFile) != "" {
		changed, err = readChangedPaths(root, changesFile)
		if err != nil {
			return SecurityTestAuditSummary{}, err
		}
	}
	head, err := currentRepositoryCommit(root)
	if err != nil {
		return SecurityTestAuditSummary{}, err
	}
	implementationStates, err := detectRuntimeImplementationStates(root, architecture)
	if err != nil {
		return SecurityTestAuditSummary{}, err
	}
	return evaluateSecurityTestPolicyWithImplementationStates(
		architecture, gates, policy, evidence, changed, head, implementationStates,
	)
}

func evaluateSecurityTestPolicy(
	architecture ArchitectureRegistry,
	gates SecurityGateRegistry,
	policy RuntimeSecurityPolicy,
	evidence RuntimeSecurityEvidence,
	changedPaths []string,
	currentCommit string,
) (SecurityTestAuditSummary, error) {
	states := map[string]string{}
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			states[boundary.ID] = "implemented"
		}
	}
	return evaluateSecurityTestPolicyWithImplementationStates(
		architecture, gates, policy, evidence, changedPaths, currentCommit, states,
	)
}

func evaluateSecurityTestPolicyWithImplementationStates(
	architecture ArchitectureRegistry,
	gates SecurityGateRegistry,
	policy RuntimeSecurityPolicy,
	evidence RuntimeSecurityEvidence,
	changedPaths []string,
	currentCommit string,
	implementationStates map[string]string,
) (SecurityTestAuditSummary, error) {
	if err := validateRuntimeSecurityPolicy(policy); err != nil {
		return SecurityTestAuditSummary{}, err
	}
	if evidence.SchemaVersion != 1 {
		return SecurityTestAuditSummary{}, fmt.Errorf("unsupported runtime security evidence schema_version %d", evidence.SchemaVersion)
	}

	var runtime []ArchitectureBoundary
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			runtime = append(runtime, boundary)
		}
	}
	sort.Slice(runtime, func(i, j int) bool { return runtime[i].ID < runtime[j].ID })

	summary := SecurityTestAuditSummary{
		RuntimeBoundaries:           len(runtime),
		RegisteredScopes:            len(policy.Scopes),
		GlobalCoverageFloorPercent:  policy.GlobalCoverageFloorPercent,
		ChangedCoverageFloorPercent: policy.ChangedSecurityCriticalCoverageFloor,
	}

	scopeByBoundary := map[string]RuntimeSecurityScope{}
	runtimeByID := map[string]ArchitectureBoundary{}
	for _, boundary := range runtime {
		runtimeByID[boundary.ID] = boundary
	}
	for _, scope := range policy.Scopes {
		if _, duplicate := scopeByBoundary[scope.BoundaryID]; duplicate {
			return summary, fmt.Errorf("duplicate runtime security scope for boundary %s", scope.BoundaryID)
		}
		boundary, ok := runtimeByID[scope.BoundaryID]
		if !ok {
			return summary, fmt.Errorf("runtime security scope %s references non-runtime or unknown boundary", scope.BoundaryID)
		}
		if err := validateRuntimeSecurityScope(scope, boundary); err != nil {
			return summary, err
		}
		scopeByBoundary[scope.BoundaryID] = scope
		if scope.AuthorizationRequired {
			summary.AuthorizationRequired++
		}
		if scope.TenantIsolationRequired {
			summary.TenantIsolationRequired++
		}
	}
	for _, boundary := range runtime {
		if _, ok := scopeByBoundary[boundary.ID]; !ok {
			return summary, fmt.Errorf("product-runtime boundary %s has no explicit security ownership scope", boundary.ID)
		}
	}

	if len(runtime) == 0 {
		if len(policy.Scopes) != 0 || len(implementationStates) != 0 {
			return summary, fmt.Errorf("runtime security state exists while no product-runtime boundary exists")
		}
		if evidence.Status != "not-applicable" {
			return summary, fmt.Errorf("runtime security evidence must be not-applicable while no product runtime exists")
		}
		if strings.TrimSpace(evidence.Reason) == "" || evidence.SourceCommit != "" || evidence.GlobalCoveragePercent != nil || len(evidence.Scopes) != 0 {
			return summary, fmt.Errorf("not-applicable runtime security evidence must contain only a reason")
		}
		summary.CoverageStatus = "not-applicable"
		summary.NotApplicableReason = evidence.Reason
		return summary, nil
	}

	if len(implementationStates) != len(runtime) {
		return summary, fmt.Errorf("runtime implementation-state count does not match product-runtime boundaries")
	}
	for _, boundary := range runtime {
		state, ok := implementationStates[boundary.ID]
		if !ok {
			return summary, fmt.Errorf("runtime boundary %s has no detected implementation state", boundary.ID)
		}
		switch state {
		case "preimplementation":
			summary.PreimplementationScopes++
		case "implemented":
			summary.ImplementedScopes++
		default:
			return summary, fmt.Errorf("runtime boundary %s has unknown implementation state %q", boundary.ID, state)
		}
	}

	evidenceByBoundary := map[string]RuntimeSecurityEvidenceScope{}
	for _, item := range evidence.Scopes {
		if _, duplicate := evidenceByBoundary[item.BoundaryID]; duplicate {
			return summary, fmt.Errorf("duplicate runtime security evidence for boundary %s", item.BoundaryID)
		}
		if _, ok := scopeByBoundary[item.BoundaryID]; !ok {
			return summary, fmt.Errorf("runtime security evidence references unregistered scope %s", item.BoundaryID)
		}
		evidenceByBoundary[item.BoundaryID] = item
	}
	if len(evidenceByBoundary) != len(scopeByBoundary) {
		return summary, fmt.Errorf("runtime security evidence scope count does not match registered runtime security scopes")
	}
	for _, boundary := range runtime {
		item := evidenceByBoundary[boundary.ID]
		state := implementationStates[boundary.ID]
		if item.ImplementationState != state {
			return summary, fmt.Errorf("scope %s implementation-state evidence mismatch: got %q want %q",
				boundary.ID, item.ImplementationState, state)
		}
		if state == "preimplementation" {
			if item.ChangedSecurityCritical || item.ChangedSecurityCriticalPercent != nil ||
				item.AuthorizationNegativePassed != nil || item.TenantIsolationNegativePassed != nil {
				return summary, fmt.Errorf("preimplementation scope %s must not claim runtime coverage or negative-test evidence", boundary.ID)
			}
		}
	}

	expectedStatus := "measured"
	if summary.ImplementedScopes == 0 {
		expectedStatus = "preimplementation"
	} else if summary.PreimplementationScopes > 0 {
		expectedStatus = "mixed"
	}
	if evidence.Status != expectedStatus {
		return summary, fmt.Errorf("runtime security evidence status %q does not match detected state %q", evidence.Status, expectedStatus)
	}
	if expectedStatus == "preimplementation" {
		if strings.TrimSpace(evidence.Reason) == "" || evidence.SourceCommit != "" || evidence.GlobalCoveragePercent != nil {
			return summary, fmt.Errorf("preimplementation runtime evidence requires a reason and must not claim commit-bound coverage")
		}
		summary.CoverageStatus = "preimplementation"
		summary.NotApplicableReason = evidence.Reason
		return summary, nil
	}

	if len(changedPaths) == 0 {
		return summary, fmt.Errorf("runtime security audit requires changed-path evidence once executable product runtime exists")
	}
	if !fullCommitIDPattern.MatchString(evidence.SourceCommit) || evidence.SourceCommit != currentCommit {
		return summary, fmt.Errorf("runtime security evidence source_commit does not match checked-out commit")
	}
	if evidence.GlobalCoveragePercent == nil || !validCoveragePercent(*evidence.GlobalCoveragePercent) {
		return summary, fmt.Errorf("runtime security evidence requires a valid global coverage percent")
	}
	if *evidence.GlobalCoveragePercent < policy.GlobalCoverageFloorPercent {
		return summary, fmt.Errorf("global runtime coverage %.2f is below %.2f floor", *evidence.GlobalCoveragePercent, policy.GlobalCoverageFloorPercent)
	}
	summary.CoverageStatus = expectedStatus
	summary.GlobalCoveragePercent = evidence.GlobalCoveragePercent

	gateByID := map[string]SecurityGate{}
	for _, gate := range gates.Gates {
		gateByID[gate.ID] = gate
	}
	for _, boundary := range runtime {
		scope := scopeByBoundary[boundary.ID]
		item := evidenceByBoundary[boundary.ID]
		if implementationStates[boundary.ID] == "preimplementation" {
			continue
		}
		changed := scopeChanged(scope, changedPaths)
		if item.ChangedSecurityCritical != changed {
			return summary, fmt.Errorf("scope %s changed-security-critical evidence mismatch", scope.BoundaryID)
		}
		if changed {
			summary.ChangedSecurityCriticalScopes++
			if item.ChangedSecurityCriticalPercent == nil || !validCoveragePercent(*item.ChangedSecurityCriticalPercent) {
				return summary, fmt.Errorf("scope %s requires changed security-critical coverage evidence", scope.BoundaryID)
			}
			if *item.ChangedSecurityCriticalPercent < policy.ChangedSecurityCriticalCoverageFloor {
				return summary, fmt.Errorf("scope %s changed security-critical coverage %.2f is below %.2f floor",
					scope.BoundaryID, *item.ChangedSecurityCriticalPercent, policy.ChangedSecurityCriticalCoverageFloor)
			}
		} else if item.ChangedSecurityCriticalPercent != nil {
			return summary, fmt.Errorf("scope %s must not claim changed coverage when no security-critical path changed", scope.BoundaryID)
		}

		if scope.AuthorizationRequired {
			if err := requireActiveNegativeGate(gateByID, policy.AuthorizationGateID, "authorization"); err != nil {
				return summary, err
			}
			if item.AuthorizationNegativePassed == nil || !*item.AuthorizationNegativePassed {
				return summary, fmt.Errorf("scope %s requires passing authorization negative tests", scope.BoundaryID)
			}
		} else if item.AuthorizationNegativePassed != nil {
			return summary, fmt.Errorf("scope %s must not claim authorization negative-test evidence when not required", scope.BoundaryID)
		}
		if scope.TenantIsolationRequired {
			if err := requireActiveNegativeGate(gateByID, policy.TenantIsolationGateID, "tenant isolation"); err != nil {
				return summary, err
			}
			if item.TenantIsolationNegativePassed == nil || !*item.TenantIsolationNegativePassed {
				return summary, fmt.Errorf("scope %s requires passing tenant-isolation negative tests", scope.BoundaryID)
			}
		} else if item.TenantIsolationNegativePassed != nil {
			return summary, fmt.Errorf("scope %s must not claim tenant-isolation evidence when not required", scope.BoundaryID)
		}
	}
	return summary, nil
}


func detectRuntimeImplementationStates(root string, architecture ArchitectureRegistry) (map[string]string, error) {
	states := map[string]string{}
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind != "product-runtime" {
			continue
		}
		state, err := detectRuntimeBoundaryImplementationState(root, boundary)
		if err != nil {
			return nil, err
		}
		states[boundary.ID] = state
	}
	return states, nil
}

func detectRuntimeBoundaryImplementationState(root string, boundary ArchitectureBoundary) (string, error) {
	implemented := false
	for _, pattern := range boundary.Roots {
		prefix := strings.TrimSuffix(patternPrefix(pattern), "/")
		if prefix == "" {
			return "", fmt.Errorf("runtime boundary %s has non-concrete root %q", boundary.ID, pattern)
		}
		path, err := resolveRepoPath(root, prefix, false)
		if err != nil {
			return "", fmt.Errorf("runtime boundary %s root: %w", boundary.ID, err)
		}
		info, err := os.Stat(path) // #nosec G703 -- path is repository-confined and symlink-free.
		if err != nil {
			return "", fmt.Errorf("runtime boundary %s root: %w", boundary.ID, err)
		}
		if !info.IsDir() {
			return "", fmt.Errorf("runtime boundary %s root must resolve to a directory", boundary.ID)
		}
		err = filepath.WalkDir(path, func(current string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("runtime boundary %s contains symlink %s", boundary.ID, current)
			}
			if entry.IsDir() {
				return nil
			}
			name := strings.ToLower(entry.Name())
			if name == "readme.md" || name == ".gitkeep" {
				return nil
			}
			implemented = true
			return filepath.SkipAll
		})
		if err != nil {
			return "", err
		}
		if implemented {
			break
		}
	}
	if implemented {
		return "implemented", nil
	}
	return "preimplementation", nil
}

func validateRuntimeSecurityPolicy(policy RuntimeSecurityPolicy) error {
	if policy.SchemaVersion != 1 {
		return fmt.Errorf("unsupported runtime security policy schema_version %d", policy.SchemaVersion)
	}
	if policy.DefaultPolicy != "deny-unregistered-runtime-boundary" {
		return fmt.Errorf("runtime security policy default must be deny-unregistered-runtime-boundary")
	}
	if policy.GlobalCoverageFloorPercent < 80 || !validCoveragePercent(policy.GlobalCoverageFloorPercent) {
		return fmt.Errorf("global runtime coverage floor must be between 80 and 100")
	}
	if policy.ChangedSecurityCriticalCoverageFloor < 90 || !validCoveragePercent(policy.ChangedSecurityCriticalCoverageFloor) {
		return fmt.Errorf("changed security-critical coverage floor must be between 90 and 100")
	}
	if policy.AuthorizationGateID != "SEC-AUTH-NEG-001" || policy.TenantIsolationGateID != "SEC-TENANT-ISO-001" {
		return fmt.Errorf("runtime security policy must reference canonical authorization and tenant-isolation gates")
	}
	return nil
}

func validateRuntimeSecurityScope(scope RuntimeSecurityScope, boundary ArchitectureBoundary) error {
	if strings.TrimSpace(scope.Owner) == "" {
		return fmt.Errorf("runtime security scope %s requires an owner", scope.BoundaryID)
	}
	if len(scope.SecurityCriticalPaths) == 0 {
		return fmt.Errorf("runtime security scope %s requires security-critical paths", scope.BoundaryID)
	}
	for _, pattern := range scope.SecurityCriticalPaths {
		normalized, err := normalizeArchitecturePattern(pattern)
		if err != nil {
			return fmt.Errorf("scope %s security-critical path: %w", scope.BoundaryID, err)
		}
		covered := false
		for _, root := range boundary.Roots {
			if patternCovers(root, normalized) {
				covered = true
				break
			}
		}
		if !covered {
			return fmt.Errorf("scope %s security-critical path %q is outside boundary roots", scope.BoundaryID, pattern)
		}
	}
	if strings.TrimSpace(scope.AuthorizationRationale) == "" || strings.TrimSpace(scope.TenantIsolationRationale) == "" {
		return fmt.Errorf("runtime security scope %s requires explicit authorization and tenant-isolation rationales", scope.BoundaryID)
	}
	return nil
}

func scopeChanged(scope RuntimeSecurityScope, changedPaths []string) bool {
	for _, changed := range changedPaths {
		for _, pattern := range scope.SecurityCriticalPaths {
			if triggerPatternMatches(pattern, changed) {
				return true
			}
		}
	}
	return false
}

func requireActiveNegativeGate(gates map[string]SecurityGate, id, label string) error {
	gate, ok := gates[id]
	if !ok {
		return fmt.Errorf("canonical %s negative-test gate %s is missing", label, id)
	}
	if gate.Readiness != "active" || gate.ImplementationKey == "" {
		return fmt.Errorf("%s negative-test gate %s must be active before runtime scope requires it", label, id)
	}
	return nil
}

func validCoveragePercent(value float64) bool {
	return value >= 0 && value <= 100
}

func currentRepositoryCommit(root string) (string, error) {
	cmd := exec.Command("git", "-C", root, "rev-parse", "HEAD") // #nosec G204,G702 -- executable/arguments are fixed and root is a validated CMDR repository root.
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("resolve current repository commit: %w", err)
	}
	commit := strings.TrimSpace(stdout.String())
	if !fullCommitIDPattern.MatchString(commit) {
		return "", fmt.Errorf("current repository commit is not a full Git object id")
	}
	return strings.ToLower(commit), nil
}
