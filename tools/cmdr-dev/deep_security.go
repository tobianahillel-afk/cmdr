package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const deepSecurityPolicyPath = "engineering/security/deep-security-policy.json"

type DeepSecurityGatePolicy struct {
	GateID      string   `json:"gate_id"`
	TargetKinds []string `json:"target_kinds"`
	PRMode      string   `json:"pr_mode"`
}

type DeepSecurityTarget struct {
	ID                 string `json:"id"`
	GateID             string `json:"gate_id"`
	Kind               string `json:"kind"`
	BoundaryID         string `json:"boundary_id"`
	Owner              string `json:"owner"`
	Target             string `json:"target"`
	Seed               string `json:"seed"`
	Configuration      string `json:"configuration"`
	MinimumDurationSec int    `json:"minimum_duration_seconds"`
	SyntheticOnly      bool   `json:"synthetic_only"`
}

type DeepSecurityPolicy struct {
	SchemaVersion     int                      `json:"schema_version"`
	DefaultPolicy     string                   `json:"default_policy"`
	EvidenceDirectory string                   `json:"evidence_directory"`
	ControlPlanePaths []string                 `json:"control_plane_paths"`
	Gates             []DeepSecurityGatePolicy `json:"gates"`
	Targets           []DeepSecurityTarget     `json:"targets"`
}

type DeepSecurityEvidence struct {
	SchemaVersion   int    `json:"schema_version"`
	GateID         string `json:"gate_id"`
	TargetID       string `json:"target_id"`
	Stage          string `json:"stage"`
	SourceCommit   string `json:"source_commit"`
	Target         string `json:"target"`
	Seed           string `json:"seed"`
	Configuration  string `json:"configuration"`
	DurationSec    int    `json:"duration_seconds"`
	EvidenceKind   string `json:"evidence_kind"`
	SyntheticData  bool   `json:"synthetic_data"`
	Result         string `json:"result"`
	ArtifactSHA256 string `json:"artifact_sha256,omitempty"`
}

type DeepSecurityDecision struct {
	GateID    string   `json:"gate_id"`
	Readiness string   `json:"readiness"`
	Decision  string   `json:"decision"`
	Reason    string   `json:"reason"`
	Targets   []string `json:"targets,omitempty"`
}

type DeepSecurityAuditSummary struct {
	Stage                   string                 `json:"stage"`
	RuntimeBoundaries       int                    `json:"runtime_boundaries"`
	SecuritySensitiveChange bool                   `json:"security_sensitive_change"`
	RegisteredTargets       int                    `json:"registered_targets"`
	SelectedTargets         int                    `json:"selected_targets"`
	DeferredGates           int                    `json:"deferred_gates"`
	ValidatedEvidence       int                    `json:"validated_evidence"`
	Decisions               []DeepSecurityDecision `json:"decisions"`
}

var deepTargetIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
var deepSHA256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

var deepTargetEvidenceKinds = map[string]string{
	"fuzz":      "fuzz-report",
	"dast":      "dast-report",
	"race":      "test-result",
	"sanitizer": "test-result",
}

func runDeepSecurityAudit(root, stage, changesFile, requestedGate string) (DeepSecurityAuditSummary, error) {
	gates, err := loadSecurityGateRegistry(root)
	if err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	var architecture ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &architecture); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(architecture); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	var runtimePolicy RuntimeSecurityPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(runtimeSecurityPolicyPath)), &runtimePolicy); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	if err := validateRuntimeSecurityPolicy(runtimePolicy); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	var policy DeepSecurityPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(deepSecurityPolicyPath)), &policy); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	if err := validateDeepSecurityPolicy(policy, gates, architecture); err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	if !knownSecurityStages[stage] {
		return DeepSecurityAuditSummary{}, fmt.Errorf("unknown deep security stage %q", stage)
	}

	var changed []string
	if stage == "pr" {
		if strings.TrimSpace(changesFile) == "" {
			return DeepSecurityAuditSummary{}, fmt.Errorf("PR deep security scheduling requires changed-path evidence")
		}
		changed, err = readChangedPaths(root, changesFile)
		if err != nil {
			return DeepSecurityAuditSummary{}, err
		}
	}
	head, err := currentRepositoryCommit(root)
	if err != nil {
		return DeepSecurityAuditSummary{}, err
	}
	return evaluateDeepSecurityPolicy(root, stage, requestedGate, changed, head, architecture, gates, runtimePolicy, policy)
}

func evaluateDeepSecurityPolicy(
	root, stage, requestedGate string,
	changedPaths []string,
	currentCommit string,
	architecture ArchitectureRegistry,
	gates SecurityGateRegistry,
	runtimePolicy RuntimeSecurityPolicy,
	policy DeepSecurityPolicy,
) (DeepSecurityAuditSummary, error) {
	gateByID := map[string]SecurityGate{}
	for _, gate := range gates.Gates {
		gateByID[gate.ID] = gate
	}
	var runtimeCount int
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			runtimeCount++
		}
	}
	sensitive := false
	if stage == "pr" {
		sensitive = deepSecuritySensitiveChange(changedPaths, policy, runtimePolicy)
	}
	summary := DeepSecurityAuditSummary{
		Stage:                   stage,
		RuntimeBoundaries:       runtimeCount,
		SecuritySensitiveChange: sensitive,
		RegisteredTargets:       len(policy.Targets),
	}

	if stage == "on-demand" && requestedGate != "" {
		if _, ok := gateByID[requestedGate]; !ok {
			return summary, fmt.Errorf("requested deep security gate %s is unknown", requestedGate)
		}
		found := false
		for _, item := range policy.Gates {
			if item.GateID == requestedGate {
				found = true
				break
			}
		}
		if !found {
			return summary, fmt.Errorf("requested gate %s is not a registered deep security gate", requestedGate)
		}
	}

	targetsByGate := map[string][]DeepSecurityTarget{}
	for _, target := range policy.Targets {
		targetsByGate[target.GateID] = append(targetsByGate[target.GateID], target)
	}
	for gateID := range targetsByGate {
		sort.Slice(targetsByGate[gateID], func(i, j int) bool {
			return targetsByGate[gateID][i].ID < targetsByGate[gateID][j].ID
		})
	}

	for _, gatePolicy := range policy.Gates {
		gate := gateByID[gatePolicy.GateID]
		decision := DeepSecurityDecision{GateID: gate.ID, Readiness: gate.Readiness}
		if stage == "pr" && !sensitive {
			decision.Decision = "not-triggered"
			decision.Reason = "PR does not touch registered security-sensitive paths"
			summary.Decisions = append(summary.Decisions, decision)
			continue
		}
		if stage == "on-demand" && requestedGate != "" && requestedGate != gate.ID {
			decision.Decision = "not-requested"
			decision.Reason = "another deep security gate was requested"
			summary.Decisions = append(summary.Decisions, decision)
			continue
		}
		if !containsString(gate.Stages, stage) {
			decision.Decision = "not-scheduled"
			decision.Reason = "gate does not declare this execution stage"
			summary.Decisions = append(summary.Decisions, decision)
			continue
		}

		targets := targetsByGate[gate.ID]
		if len(targets) == 0 {
			decision.Decision = "deferred"
			decision.Reason = "no concrete runtime target is registered"
			summary.DeferredGates++
			summary.Decisions = append(summary.Decisions, decision)
			continue
		}
		if gate.Readiness != "active" {
			return summary, fmt.Errorf("deep security gate %s has registered targets but readiness is %s", gate.ID, gate.Readiness)
		}
		for _, target := range targets {
			decision.Targets = append(decision.Targets, target.ID)
			evidence, err := loadDeepSecurityEvidence(root, policy.EvidenceDirectory, target.ID)
			if err != nil {
				return summary, fmt.Errorf("target %s: %w", target.ID, err)
			}
			if err := validateDeepSecurityEvidence(evidence, target, gate, stage, currentCommit); err != nil {
				return summary, fmt.Errorf("target %s: %w", target.ID, err)
			}
			summary.SelectedTargets++
			summary.ValidatedEvidence++
		}
		decision.Decision = "evidence-validated"
		decision.Reason = "active gate has registered target(s) and commit-bound passing evidence"
		summary.Decisions = append(summary.Decisions, decision)
	}
	sort.Slice(summary.Decisions, func(i, j int) bool {
		return summary.Decisions[i].GateID < summary.Decisions[j].GateID
	})
	return summary, nil
}

func validateDeepSecurityPolicy(policy DeepSecurityPolicy, gates SecurityGateRegistry, architecture ArchitectureRegistry) error {
	if policy.SchemaVersion != 1 {
		return fmt.Errorf("unsupported deep security policy schema_version %d", policy.SchemaVersion)
	}
	if policy.DefaultPolicy != "deny-unregistered-deep-security-target" {
		return fmt.Errorf("deep security policy default must be deny-unregistered-deep-security-target")
	}
	if policy.EvidenceDirectory != "engineering/testing/deep-security-evidence" {
		return fmt.Errorf("deep security evidence directory must be engineering/testing/deep-security-evidence")
	}
	if len(policy.ControlPlanePaths) == 0 {
		return fmt.Errorf("deep security policy requires control-plane paths")
	}
	for _, pattern := range policy.ControlPlanePaths {
		if _, err := normalizeArchitecturePattern(pattern); err != nil {
			return fmt.Errorf("deep security control-plane path: %w", err)
		}
	}

	gateByID := map[string]SecurityGate{}
	for _, gate := range gates.Gates {
		gateByID[gate.ID] = gate
	}
	expected := map[string][]string{
		"SEC-FUZZ-001": {"fuzz"},
		"SEC-DAST-001": {"dast"},
		"SEC-RACE-001": {"race", "sanitizer"},
	}
	if len(policy.Gates) != len(expected) {
		return fmt.Errorf("deep security policy must register exactly %d canonical dynamic gates", len(expected))
	}
	seen := map[string]DeepSecurityGatePolicy{}
	for _, item := range policy.Gates {
		wantKinds, ok := expected[item.GateID]
		if !ok {
			return fmt.Errorf("deep security policy uses non-canonical gate %s", item.GateID)
		}
		if _, duplicate := seen[item.GateID]; duplicate {
			return fmt.Errorf("duplicate deep security gate policy %s", item.GateID)
		}
		seen[item.GateID] = item
		if item.PRMode != "security-sensitive-only" {
			return fmt.Errorf("deep security gate %s must use security-sensitive-only PR mode", item.GateID)
		}
		if !sameStringSet(item.TargetKinds, wantKinds) {
			return fmt.Errorf("deep security gate %s target kinds mismatch", item.GateID)
		}
		gate, ok := gateByID[item.GateID]
		if !ok {
			return fmt.Errorf("deep security gate %s is missing from security registry", item.GateID)
		}
		if !containsString(gate.Stages, "pr") || !containsString(gate.Stages, "on-demand") {
			return fmt.Errorf("deep security gate %s must declare PR and on-demand stages", item.GateID)
		}
	}

	runtimeByID := map[string]ArchitectureBoundary{}
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			runtimeByID[boundary.ID] = boundary
		}
	}
	seenTargets := map[string]bool{}
	for _, target := range policy.Targets {
		if !deepTargetIDPattern.MatchString(target.ID) {
			return fmt.Errorf("invalid deep security target id %q", target.ID)
		}
		if seenTargets[target.ID] {
			return fmt.Errorf("duplicate deep security target %s", target.ID)
		}
		seenTargets[target.ID] = true
		gatePolicy, ok := seen[target.GateID]
		if !ok {
			return fmt.Errorf("target %s references unregistered deep security gate %s", target.ID, target.GateID)
		}
		if !containsString(gatePolicy.TargetKinds, target.Kind) {
			return fmt.Errorf("target %s kind %s is not allowed for gate %s", target.ID, target.Kind, target.GateID)
		}
		if _, ok := deepTargetEvidenceKinds[target.Kind]; !ok {
			return fmt.Errorf("target %s uses unknown kind %s", target.ID, target.Kind)
		}
		if _, ok := runtimeByID[target.BoundaryID]; !ok {
			return fmt.Errorf("target %s references non-runtime or unknown boundary %s", target.ID, target.BoundaryID)
		}
		if strings.TrimSpace(target.Owner) == "" || strings.TrimSpace(target.Target) == "" ||
			strings.TrimSpace(target.Seed) == "" || strings.TrimSpace(target.Configuration) == "" {
			return fmt.Errorf("target %s requires owner, target, seed and configuration", target.ID)
		}
		if target.MinimumDurationSec <= 0 {
			return fmt.Errorf("target %s requires a positive minimum duration", target.ID)
		}
		if !target.SyntheticOnly {
			return fmt.Errorf("target %s must be synthetic-data-only", target.ID)
		}
	}
	return nil
}

func deepSecuritySensitiveChange(changed []string, policy DeepSecurityPolicy, runtimePolicy RuntimeSecurityPolicy) bool {
	for _, path := range changed {
		for _, pattern := range policy.ControlPlanePaths {
			if triggerPatternMatches(pattern, path) {
				return true
			}
		}
		for _, scope := range runtimePolicy.Scopes {
			for _, pattern := range scope.SecurityCriticalPaths {
				if triggerPatternMatches(pattern, path) {
					return true
				}
			}
		}
	}
	return false
}

func loadDeepSecurityEvidence(root, directory, targetID string) (DeepSecurityEvidence, error) {
	var evidence DeepSecurityEvidence
	if !deepTargetIDPattern.MatchString(targetID) {
		return evidence, fmt.Errorf("invalid target id")
	}
	path := filepath.Join(root, filepath.FromSlash(directory), targetID+".json")
	if err := decodeStrict(root, path, &evidence); err != nil {
		return evidence, fmt.Errorf("load deep security evidence: %w", err)
	}
	return evidence, nil
}

func validateDeepSecurityEvidence(evidence DeepSecurityEvidence, target DeepSecurityTarget, gate SecurityGate, stage, currentCommit string) error {
	if evidence.SchemaVersion != 1 {
		return fmt.Errorf("unsupported evidence schema_version %d", evidence.SchemaVersion)
	}
	if evidence.GateID != target.GateID || evidence.TargetID != target.ID || evidence.Stage != stage {
		return fmt.Errorf("evidence gate/target/stage identity mismatch")
	}
	if evidence.SourceCommit != currentCommit || !fullCommitIDPattern.MatchString(evidence.SourceCommit) {
		return fmt.Errorf("evidence source_commit does not match checked-out commit")
	}
	if evidence.Target != target.Target || evidence.Seed != target.Seed || evidence.Configuration != target.Configuration {
		return fmt.Errorf("evidence target/seed/configuration mismatch")
	}
	if evidence.DurationSec < target.MinimumDurationSec {
		return fmt.Errorf("evidence duration %d is below target minimum %d", evidence.DurationSec, target.MinimumDurationSec)
	}
	wantKind := deepTargetEvidenceKinds[target.Kind]
	if evidence.EvidenceKind != wantKind || !containsString(gate.EvidenceKinds, wantKind) {
		return fmt.Errorf("evidence kind %q is invalid for target %s", evidence.EvidenceKind, target.ID)
	}
	if !evidence.SyntheticData {
		return fmt.Errorf("deep security evidence must use synthetic data")
	}
	if evidence.Result != "pass" {
		return fmt.Errorf("deep security evidence result must be pass")
	}
	if stage == "release" {
		if !deepSHA256Pattern.MatchString(evidence.ArtifactSHA256) {
			return fmt.Errorf("release deep security evidence requires exact artifact SHA-256")
		}
	} else if evidence.ArtifactSHA256 != "" && !deepSHA256Pattern.MatchString(evidence.ArtifactSHA256) {
		return fmt.Errorf("invalid artifact SHA-256")
	}
	return nil
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[string]bool{}
	for _, value := range a {
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	for _, value := range b {
		if !seen[value] {
			return false
		}
	}
	return true
}
