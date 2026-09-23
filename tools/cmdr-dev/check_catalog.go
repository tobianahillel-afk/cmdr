package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const checkCatalogPath = "engineering/testing/check-catalog.json"

type ValidationCheck struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	ExecutorKey   string   `json:"executor_key"`
	RiskDomains   []string `json:"risk_domains"`
	CostTier      string   `json:"cost_tier"`
	Mandatory     bool     `json:"mandatory"`
	AlwaysOnPR    bool     `json:"always_on_pr"`
	TriggerPaths  []string `json:"trigger_paths"`
	Prerequisites []string `json:"prerequisites"`
	EvidenceKind  string   `json:"evidence_kind"`
}

type CheckCatalog struct {
	SchemaVersion int               `json:"schema_version"`
	DefaultPolicy string            `json:"default_policy"`
	RiskDomains   []string          `json:"risk_domains"`
	CostTiers     []string          `json:"cost_tiers"`
	Checks        []ValidationCheck `json:"checks"`
}

type CheckCatalogSummary struct {
	Checks        int            `json:"checks"`
	Mandatory     int            `json:"mandatory"`
	AlwaysOnPR    int            `json:"always_on_pr"`
	ByCostTier    map[string]int `json:"by_cost_tier"`
	ByRiskDomain  map[string]int `json:"by_risk_domain"`
	Prerequisites int            `json:"prerequisite_edges"`
}

var checkIDPattern = regexp.MustCompile(`^CHK-[A-Z0-9-]+$`)

var knownRiskDomains = map[string]bool{
	"code-quality":          true,
	"unit-correctness":      true,
	"product-spec":          true,
	"coverage-traceability": true,
	"work-governance":       true,
	"architecture":          true,
	"runtime-dependencies":  true,
	"security":              true,
	"agent-context":         true,
	"repository-health":     true,
	"unknown":               true,
}

var knownCostTiers = map[string]bool{
	"fast":     true,
	"standard": true,
	"heavy":    true,
}

var knownExecutorKeys = map[string]bool{
	"gofmt":               true,
	"go-vet":              true,
	"go-unit":             true,
	"spec-index":          true,
	"spec-baseline":       true,
	"coverage-graph":      true,
	"obligations":         true,
	"coverage-audit":      true,
	"validate-manifests":  true,
	"architecture-audit":  true,
	"dependency-audit":    true,
	"boundary-edge-audit": true,
	"check-catalog-audit": true,
	"impact":              true,
	"validation-plan":     true,
	"git-changes":         true,
	"validation-run":      true,
	"security-gate-audit": true,
	"secret-scan":         true,
	"gosec-go":            true,
	"govulncheck-go":      true,
	"complexity-audit":    true,
	"context":             true,
	"doctor":              true,
	"next":                true,
}

func loadCheckCatalog(root string) (CheckCatalog, error) {
	var catalog CheckCatalog
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(checkCatalogPath)), &catalog); err != nil {
		return catalog, err
	}
	if err := validateCheckCatalog(catalog); err != nil {
		return catalog, err
	}
	return catalog, nil
}

func runCheckCatalogAudit(root string) (CheckCatalogSummary, error) {
	catalog, err := loadCheckCatalog(root)
	if err != nil {
		return CheckCatalogSummary{}, err
	}
	summary := CheckCatalogSummary{
		Checks:       len(catalog.Checks),
		ByCostTier:   map[string]int{},
		ByRiskDomain: map[string]int{},
	}
	for _, check := range catalog.Checks {
		if check.Mandatory {
			summary.Mandatory++
		}
		if check.AlwaysOnPR {
			summary.AlwaysOnPR++
		}
		summary.ByCostTier[check.CostTier]++
		for _, domain := range check.RiskDomains {
			summary.ByRiskDomain[domain]++
		}
		summary.Prerequisites += len(check.Prerequisites)
	}
	return summary, nil
}

func validateCheckCatalog(catalog CheckCatalog) error {
	if catalog.SchemaVersion != 1 {
		return fmt.Errorf("unsupported check catalog schema_version %d", catalog.SchemaVersion)
	}
	if catalog.DefaultPolicy != "deny-unknown-check" {
		return fmt.Errorf("check catalog default policy must be deny-unknown-check")
	}
	if len(catalog.Checks) == 0 {
		return fmt.Errorf("check catalog is empty")
	}
	if err := validateDeclaredEnumSet("risk domain", catalog.RiskDomains, knownRiskDomains); err != nil {
		return err
	}
	if err := validateDeclaredEnumSet("cost tier", catalog.CostTiers, knownCostTiers); err != nil {
		return err
	}

	index := map[string]ValidationCheck{}
	executorOwners := map[string]string{}
	for _, check := range catalog.Checks {
		if !checkIDPattern.MatchString(check.ID) {
			return fmt.Errorf("invalid check id %q", check.ID)
		}
		if _, exists := index[check.ID]; exists {
			return fmt.Errorf("duplicate check id %s", check.ID)
		}
		if strings.TrimSpace(check.Title) == "" || strings.TrimSpace(check.EvidenceKind) == "" {
			return fmt.Errorf("check %s requires title and evidence_kind", check.ID)
		}
		if !knownExecutorKeys[check.ExecutorKey] {
			return fmt.Errorf("check %s uses unknown executor_key %q", check.ID, check.ExecutorKey)
		}
		if owner, exists := executorOwners[check.ExecutorKey]; exists {
			return fmt.Errorf("executor_key %s is shared by %s and %s", check.ExecutorKey, owner, check.ID)
		}
		executorOwners[check.ExecutorKey] = check.ID
		if !knownCostTiers[check.CostTier] {
			return fmt.Errorf("check %s uses unknown cost tier %q", check.ID, check.CostTier)
		}
		if len(check.RiskDomains) == 0 {
			return fmt.Errorf("check %s has no risk domains", check.ID)
		}
		seenDomains := map[string]bool{}
		for _, domain := range check.RiskDomains {
			if !knownRiskDomains[domain] {
				return fmt.Errorf("check %s uses unknown risk domain %q", check.ID, domain)
			}
			if seenDomains[domain] {
				return fmt.Errorf("check %s duplicates risk domain %s", check.ID, domain)
			}
			seenDomains[domain] = true
		}
		if check.Mandatory && len(check.RiskDomains) == 0 {
			return fmt.Errorf("mandatory check %s has no trigger domains", check.ID)
		}
		if len(check.TriggerPaths) == 0 {
			return fmt.Errorf("check %s has no trigger paths", check.ID)
		}
		seenPaths := map[string]bool{}
		for _, pattern := range check.TriggerPaths {
			if _, err := normalizeArchitecturePattern(pattern); err != nil {
				return fmt.Errorf("check %s trigger path: %w", check.ID, err)
			}
			if seenPaths[pattern] {
				return fmt.Errorf("check %s duplicates trigger path %s", check.ID, pattern)
			}
			seenPaths[pattern] = true
		}
		index[check.ID] = check
	}
	for _, check := range catalog.Checks {
		seen := map[string]bool{}
		for _, prerequisite := range check.Prerequisites {
			if prerequisite == check.ID {
				return fmt.Errorf("check %s cannot depend on itself", check.ID)
			}
			if _, ok := index[prerequisite]; !ok {
				return fmt.Errorf("check %s depends on unknown check %s", check.ID, prerequisite)
			}
			if seen[prerequisite] {
				return fmt.Errorf("check %s duplicates prerequisite %s", check.ID, prerequisite)
			}
			seen[prerequisite] = true
		}
	}
	return validateCheckPrerequisiteAcyclic(index)
}

func validateDeclaredEnumSet(name string, declared []string, known map[string]bool) error {
	if len(declared) != len(known) {
		return fmt.Errorf("%s enum mismatch: expected %d values, got %d", name, len(known), len(declared))
	}
	seen := map[string]bool{}
	for _, value := range declared {
		if !known[value] {
			return fmt.Errorf("unknown %s %q", name, value)
		}
		if seen[value] {
			return fmt.Errorf("duplicate %s %q", name, value)
		}
		seen[value] = true
	}
	for value := range known {
		if !seen[value] {
			return fmt.Errorf("missing %s %q", name, value)
		}
	}
	return nil
}

func validateCheckPrerequisiteAcyclic(index map[string]ValidationCheck) error {
	const (
		visiting = 1
		done     = 2
	)
	state := map[string]int{}
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case visiting:
			return fmt.Errorf("check prerequisite cycle at %s", id)
		case done:
			return nil
		}
		state[id] = visiting
		deps := append([]string(nil), index[id].Prerequisites...)
		sort.Strings(deps)
		for _, dep := range deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = done
		return nil
	}
	ids := make([]string, 0, len(index))
	for id := range index {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
