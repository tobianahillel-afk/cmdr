package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const securityGateRegistryPath = "engineering/security/security-gates.json"

type SecurityGate struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	ControlFamily       string   `json:"control_family"`
	Readiness           string   `json:"readiness"`
	Stages              []string `json:"stages"`
	BlockingPolicy      string   `json:"blocking_policy"`
	RiskDomains         []string `json:"risk_domains"`
	TriggerPaths        []string `json:"trigger_paths"`
	DataPolicy          string   `json:"data_policy"`
	EvidenceKinds       []string `json:"evidence_kinds"`
	ImplementationKey   string   `json:"implementation_key"`
	ActivationCondition string   `json:"activation_condition"`
}

type SecurityGateRegistry struct {
	SchemaVersion   int            `json:"schema_version"`
	DefaultPolicy   string         `json:"default_policy"`
	ReadinessStates []string       `json:"readiness_states"`
	ExecutionStages []string       `json:"execution_stages"`
	DataPolicies    []string       `json:"data_policies"`
	EvidenceKinds   []string       `json:"evidence_kinds"`
	Gates           []SecurityGate `json:"gates"`
}

type SecurityGateAuditSummary struct {
	Gates       int            `json:"gates"`
	ByReadiness map[string]int `json:"by_readiness"`
	ByFamily    map[string]int `json:"by_family"`
	ByStage     map[string]int `json:"by_stage"`
	Blocking    int            `json:"blocking"`
	Conditional int            `json:"conditional_blocking"`
}

var securityGateIDPattern = regexp.MustCompile(`^SEC-[A-Z0-9-]+$`)

var knownSecurityFamilies = map[string]bool{
	"secrets":                true,
	"sast":                   true,
	"sca":                    true,
	"sbom":                   true,
	"iac":                    true,
	"container":              true,
	"authorization-negative": true,
	"tenant-isolation":       true,
	"fuzz":                   true,
	"dast":                   true,
	"race-concurrency":       true,
	"supply-chain":           true,
}

var knownSecurityReadiness = map[string]bool{
	"specified":        true,
	"active":           true,
	"deferred-runtime": true,
}

var knownSecurityStages = map[string]bool{
	"pr":        true,
	"nightly":   true,
	"release":   true,
	"on-demand": true,
}

var knownSecurityDataPolicies = map[string]bool{
	"source-local-only":      true,
	"artifact-local-only":    true,
	"synthetic-runtime-only": true,
}

var knownSecurityEvidenceKinds = map[string]bool{
	"process-exit":   true,
	"finding-report": true,
	"sbom":           true,
	"test-result":    true,
	"fuzz-report":    true,
	"dast-report":    true,
	"provenance":     true,
}

var knownSecurityImplementationKeys = map[string]bool{
	"builtin-secret-scan-v1":      true,
	"gosec-go-v1":                 true,
	"govulncheck-go-v1":           true,
	"builtin-cyclonedx17-sbom-v1":          true,
	"go-context-envelope-auth-negative-v1": true,
	"go-context-envelope-tenant-iso-v1":    true,
}

func loadSecurityGateRegistry(root string) (SecurityGateRegistry, error) {
	var registry SecurityGateRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(securityGateRegistryPath)), &registry); err != nil {
		return registry, err
	}
	if err := validateSecurityGateRegistry(registry); err != nil {
		return registry, err
	}
	return registry, nil
}

func runSecurityGateAudit(root string) (SecurityGateAuditSummary, error) {
	registry, err := loadSecurityGateRegistry(root)
	if err != nil {
		return SecurityGateAuditSummary{}, err
	}
	summary := SecurityGateAuditSummary{
		Gates:       len(registry.Gates),
		ByReadiness: map[string]int{},
		ByFamily:    map[string]int{},
		ByStage:     map[string]int{},
	}
	for _, gate := range registry.Gates {
		summary.ByReadiness[gate.Readiness]++
		summary.ByFamily[gate.ControlFamily]++
		for _, stage := range gate.Stages {
			summary.ByStage[stage]++
		}
		switch gate.BlockingPolicy {
		case "blocking":
			summary.Blocking++
		case "conditional-blocking":
			summary.Conditional++
		}
	}
	return summary, nil
}

func validateSecurityGateRegistry(registry SecurityGateRegistry) error {
	if registry.SchemaVersion != 1 {
		return fmt.Errorf("unsupported security gate schema_version %d", registry.SchemaVersion)
	}
	if registry.DefaultPolicy != "deny-unregistered-security-control" {
		return fmt.Errorf("security gate default policy must be deny-unregistered-security-control")
	}
	if err := validateSecurityEnumSet("readiness", registry.ReadinessStates, knownSecurityReadiness); err != nil {
		return err
	}
	if err := validateSecurityEnumSet("execution stage", registry.ExecutionStages, knownSecurityStages); err != nil {
		return err
	}
	if err := validateSecurityEnumSet("data policy", registry.DataPolicies, knownSecurityDataPolicies); err != nil {
		return err
	}
	if err := validateSecurityEnumSet("evidence kind", registry.EvidenceKinds, knownSecurityEvidenceKinds); err != nil {
		return err
	}
	if len(registry.Gates) == 0 {
		return fmt.Errorf("security gate registry is empty")
	}

	seenIDs := map[string]bool{}
	seenFamilies := map[string]bool{}
	for _, gate := range registry.Gates {
		if !securityGateIDPattern.MatchString(gate.ID) {
			return fmt.Errorf("invalid security gate id %q", gate.ID)
		}
		if seenIDs[gate.ID] {
			return fmt.Errorf("duplicate security gate id %s", gate.ID)
		}
		seenIDs[gate.ID] = true
		if strings.TrimSpace(gate.Title) == "" || strings.TrimSpace(gate.ActivationCondition) == "" {
			return fmt.Errorf("security gate %s requires title and activation_condition", gate.ID)
		}
		if !knownSecurityFamilies[gate.ControlFamily] {
			return fmt.Errorf("security gate %s uses unknown control family %q", gate.ID, gate.ControlFamily)
		}
		if seenFamilies[gate.ControlFamily] {
			return fmt.Errorf("security control family %s has more than one canonical gate", gate.ControlFamily)
		}
		seenFamilies[gate.ControlFamily] = true
		if !knownSecurityReadiness[gate.Readiness] {
			return fmt.Errorf("security gate %s uses unknown readiness %q", gate.ID, gate.Readiness)
		}
		if gate.Readiness == "active" {
			if gate.ImplementationKey == "" || !knownSecurityImplementationKeys[gate.ImplementationKey] {
				return fmt.Errorf("active security gate %s has unknown implementation key %q", gate.ID, gate.ImplementationKey)
			}
		} else if gate.ImplementationKey != "" {
			return fmt.Errorf("non-active security gate %s must not claim implementation key %q", gate.ID, gate.ImplementationKey)
		}
		if gate.BlockingPolicy != "blocking" && gate.BlockingPolicy != "conditional-blocking" {
			return fmt.Errorf("security gate %s has invalid blocking policy %q", gate.ID, gate.BlockingPolicy)
		}
		if !knownSecurityDataPolicies[gate.DataPolicy] {
			return fmt.Errorf("security gate %s uses unknown data policy %q", gate.ID, gate.DataPolicy)
		}
		if len(gate.Stages) == 0 || len(gate.RiskDomains) == 0 || len(gate.TriggerPaths) == 0 || len(gate.EvidenceKinds) == 0 {
			return fmt.Errorf("security gate %s has incomplete trigger/evidence metadata", gate.ID)
		}
		if !containsString(gate.RiskDomains, "security") {
			return fmt.Errorf("security gate %s must include the security risk domain", gate.ID)
		}
		if err := validateSecurityValues(gate.ID, "stage", gate.Stages, knownSecurityStages); err != nil {
			return err
		}
		if err := validateSecurityValues(gate.ID, "risk domain", gate.RiskDomains, knownRiskDomains); err != nil {
			return err
		}
		if err := validateSecurityValues(gate.ID, "evidence kind", gate.EvidenceKinds, knownSecurityEvidenceKinds); err != nil {
			return err
		}
		seenPaths := map[string]bool{}
		for _, pattern := range gate.TriggerPaths {
			if strings.TrimSpace(pattern) == "" {
				return fmt.Errorf("security gate %s has empty trigger path", gate.ID)
			}
			if _, err := normalizeArchitecturePattern(pattern); err != nil {
				return fmt.Errorf("security gate %s trigger path: %w", gate.ID, err)
			}
			if seenPaths[pattern] {
				return fmt.Errorf("security gate %s duplicates trigger path %s", gate.ID, pattern)
			}
			seenPaths[pattern] = true
		}
	}
	for family := range knownSecurityFamilies {
		if !seenFamilies[family] {
			return fmt.Errorf("security control family %s has no canonical gate", family)
		}
	}
	return nil
}

func validateSecurityEnumSet(name string, declared []string, known map[string]bool) error {
	if len(declared) != len(known) {
		return fmt.Errorf("%s enum mismatch: expected %d values, got %d", name, len(known), len(declared))
	}
	return validateSecurityValues("registry", name, declared, known)
}

func validateSecurityValues(owner, name string, values []string, known map[string]bool) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !known[value] {
			return fmt.Errorf("%s uses unknown %s %q", owner, name, value)
		}
		if seen[value] {
			return fmt.Errorf("%s duplicates %s %q", owner, name, value)
		}
		seen[value] = true
	}
	if owner == "registry" {
		keys := make([]string, 0, len(known))
		for value := range known {
			keys = append(keys, value)
		}
		sort.Strings(keys)
		for _, value := range keys {
			if !seen[value] {
				return fmt.Errorf("registry is missing %s %q", name, value)
			}
		}
	}
	return nil
}
