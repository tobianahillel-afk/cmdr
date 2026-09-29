package main

import (
	"fmt"
	"sort"
	"strings"
)

type PlannedCheck struct {
	ID          string   `json:"id"`
	ExecutorKey string   `json:"executor_key"`
	CostTier    string   `json:"cost_tier"`
	Mandatory   bool     `json:"mandatory"`
	Reasons     []string `json:"reasons"`
}

type ValidationPlan struct {
	WorkUnit       string         `json:"work_unit"`
	Tier           string         `json:"tier"`
	HighRisk       bool           `json:"high_risk"`
	RiskDomains    []string       `json:"risk_domains"`
	UnknownPaths   []string       `json:"unknown_paths,omitempty"`
	SelectedChecks []PlannedCheck `json:"selected_checks"`
	ByCostTier     map[string]int `json:"by_cost_tier"`
	CostUnits      int            `json:"cost_units"`
}

func runValidationPlan(root, workUnit, changesFile string, state CurrentState, graph WorkGraph) (ValidationPlan, error) {
	impact, err := runImpactAnalysis(root, workUnit, changesFile, state, graph)
	if err != nil {
		return ValidationPlan{}, err
	}
	catalog, err := loadCheckCatalog(root)
	if err != nil {
		return ValidationPlan{}, err
	}
	plan, err := buildValidationPlan(impact, catalog)
	if err != nil {
		return ValidationPlan{}, err
	}
	if err := validateValidationPlan(plan, catalog); err != nil {
		return ValidationPlan{}, err
	}
	return plan, nil
}

func buildValidationPlan(impact ImpactReport, catalog CheckCatalog) (ValidationPlan, error) {
	index := map[string]ValidationCheck{}
	for _, check := range catalog.Checks {
		index[check.ID] = check
	}
	selected := map[string]map[string]bool{}
	add := func(id, reason string) error {
		if _, ok := index[id]; !ok {
			return fmt.Errorf("cannot select unknown check %s", id)
		}
		if selected[id] == nil {
			selected[id] = map[string]bool{}
		}
		selected[id][reason] = true
		return nil
	}

	for _, check := range catalog.Checks {
		if check.AlwaysOnPR {
			if err := add(check.ID, "pr-safety-floor"); err != nil {
				return ValidationPlan{}, err
			}
		}
		for _, changed := range impact.ChangedPaths {
			for _, pattern := range check.TriggerPaths {
				if triggerPatternMatches(pattern, changed) {
					if err := add(check.ID, "path-trigger:"+changed); err != nil {
						return ValidationPlan{}, err
					}
					break
				}
			}
		}
	}

	manifestDomains := map[string]bool{}
	for _, evidence := range impact.Evidence {
		if strings.HasPrefix(evidence.Reason, "manifest-") {
			manifestDomains[evidence.Domain] = true
		}
	}
	if len(manifestDomains) > 0 {
		for _, check := range catalog.Checks {
			if intersectsDomains(check.RiskDomains, manifestDomains) {
				if err := add(check.ID, "manifest-scope-domain"); err != nil {
					return ValidationPlan{}, err
				}
			}
		}
	}

	tier := normalImpactTier(impact)
	if impact.HighRisk || len(impact.UnknownPaths) > 0 || containsString(impact.RiskDomains, "unknown") {
		tier = "strict"
		for _, check := range catalog.Checks {
			if check.Mandatory {
				if err := add(check.ID, "strict-risk-escalation"); err != nil {
					return ValidationPlan{}, err
				}
			}
		}
	}

	if err := closeCheckPrerequisites(selected, index); err != nil {
		return ValidationPlan{}, err
	}
	order, err := topologicalSelectedChecks(selected, index)
	if err != nil {
		return ValidationPlan{}, err
	}

	plan := ValidationPlan{
		WorkUnit:     impact.WorkUnit,
		Tier:         tier,
		HighRisk:     impact.HighRisk,
		RiskDomains:  append([]string(nil), impact.RiskDomains...),
		UnknownPaths: append([]string(nil), impact.UnknownPaths...),
		ByCostTier:   map[string]int{},
	}
	for _, id := range order {
		check := index[id]
		var reasons []string
		for reason := range selected[id] {
			reasons = append(reasons, reason)
		}
		sort.Strings(reasons)
		plan.SelectedChecks = append(plan.SelectedChecks, PlannedCheck{
			ID: id, ExecutorKey: check.ExecutorKey, CostTier: check.CostTier,
			Mandatory: check.Mandatory, Reasons: reasons,
		})
		plan.ByCostTier[check.CostTier]++
		plan.CostUnits += costTierUnits(check.CostTier)
	}
	return plan, nil
}

func normalImpactTier(impact ImpactReport) string {
	for _, domain := range impact.RiskDomains {
		switch domain {
		case "product-spec", "coverage-traceability", "architecture", "runtime-dependencies", "security", "performance":
			return "standard"
		}
	}
	return "fast"
}

func closeCheckPrerequisites(selected map[string]map[string]bool, index map[string]ValidationCheck) error {
	var ensure func(string) error
	ensure = func(id string) error {
		check, ok := index[id]
		if !ok {
			return fmt.Errorf("unknown selected check %s", id)
		}
		for _, prerequisite := range check.Prerequisites {
			if _, ok := index[prerequisite]; !ok {
				return fmt.Errorf("check %s depends on unknown check %s", id, prerequisite)
			}
			if selected[prerequisite] == nil {
				selected[prerequisite] = map[string]bool{}
			}
			selected[prerequisite]["prerequisite-of:"+id] = true
			if err := ensure(prerequisite); err != nil {
				return err
			}
		}
		return nil
	}
	var ids []string
	for id := range selected {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := ensure(id); err != nil {
			return err
		}
	}
	return nil
}

func topologicalSelectedChecks(selected map[string]map[string]bool, index map[string]ValidationCheck) ([]string, error) {
	const (
		visiting = 1
		done     = 2
	)
	state := map[string]int{}
	var order []string
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case visiting:
			return fmt.Errorf("selected check cycle at %s", id)
		case done:
			return nil
		}
		state[id] = visiting
		check, ok := index[id]
		if !ok {
			return fmt.Errorf("selected unknown check %s", id)
		}
		deps := append([]string(nil), check.Prerequisites...)
		sort.Strings(deps)
		for _, dep := range deps {
			if selected[dep] == nil {
				continue
			}
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = done
		order = append(order, id)
		return nil
	}
	var ids []string
	for id := range selected {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return nil, err
		}
	}
	return order, nil
}

func validateValidationPlan(plan ValidationPlan, catalog CheckCatalog) error {
	if plan.Tier != "fast" && plan.Tier != "standard" && plan.Tier != "strict" {
		return fmt.Errorf("invalid validation tier %q", plan.Tier)
	}
	index := map[string]ValidationCheck{}
	position := map[string]int{}
	for _, check := range catalog.Checks {
		index[check.ID] = check
	}
	for i, planned := range plan.SelectedChecks {
		check, ok := index[planned.ID]
		if !ok {
			return fmt.Errorf("plan includes unknown check %s", planned.ID)
		}
		if planned.ExecutorKey != check.ExecutorKey || planned.CostTier != check.CostTier || planned.Mandatory != check.Mandatory {
			return fmt.Errorf("plan metadata mismatch for %s", planned.ID)
		}
		if len(planned.Reasons) == 0 {
			return fmt.Errorf("plan check %s has no selection reason", planned.ID)
		}
		if _, exists := position[planned.ID]; exists {
			return fmt.Errorf("plan duplicates check %s", planned.ID)
		}
		position[planned.ID] = i
	}
	for _, check := range catalog.Checks {
		if check.AlwaysOnPR {
			if _, ok := position[check.ID]; !ok {
				return fmt.Errorf("plan omitted PR safety-floor check %s", check.ID)
			}
		}
		if plan.Tier == "strict" && check.Mandatory {
			if _, ok := position[check.ID]; !ok {
				return fmt.Errorf("strict plan omitted mandatory check %s", check.ID)
			}
		}
	}
	for id, pos := range position {
		for _, prerequisite := range index[id].Prerequisites {
			depPos, ok := position[prerequisite]
			if !ok {
				return fmt.Errorf("plan omitted prerequisite %s for %s", prerequisite, id)
			}
			if depPos >= pos {
				return fmt.Errorf("prerequisite %s must precede %s", prerequisite, id)
			}
		}
	}
	return nil
}

func intersectsDomains(domains []string, selected map[string]bool) bool {
	for _, domain := range domains {
		if selected[domain] {
			return true
		}
	}
	return false
}

func costTierUnits(tier string) int {
	switch tier {
	case "fast":
		return 1
	case "standard":
		return 5
	case "heavy":
		return 20
	default:
		return 100
	}
}
