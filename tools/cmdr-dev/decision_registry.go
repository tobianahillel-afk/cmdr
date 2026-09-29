package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	decisionPolicyPath   = "engineering/decisions/decision-policy.json"
	decisionRegistryPath = "engineering/decisions/decision-registry.json"
)

type DecisionClassPolicy struct {
	ID                        string `json:"id"`
	RecordRequired            bool   `json:"record_required"`
	CurrentExternalResearch   string `json:"current_external_research"`
	AlternativesRequired      bool   `json:"alternatives_required"`
	BenchmarkRule             string `json:"benchmark_rule"`
	AdversarialReviewRequired bool   `json:"adversarial_review_required"`
}

type DecisionPolicy struct {
	SchemaVersion     int                   `json:"schema_version"`
	DefaultPolicy     string                `json:"default_policy"`
	Statuses          []string              `json:"statuses"`
	ProductBoundaries []string              `json:"product_boundaries"`
	RiskTags          []string              `json:"risk_tags"`
	Classes           []DecisionClassPolicy `json:"classes"`
}

type EngineeringDecision struct {
	ID                   string   `json:"id"`
	DecisionKey          string   `json:"decision_key"`
	Title                string   `json:"title"`
	Class                string   `json:"class"`
	Status               string   `json:"status"`
	Owner                string   `json:"owner"`
	Question             string   `json:"question"`
	Rationale            string   `json:"rationale"`
	ProductBoundary      string   `json:"product_boundary"`
	AffectedWorkUnits    []string `json:"affected_work_units"`
	Constraints          []string `json:"constraints"`
	RiskTags             []string `json:"risk_tags"`
	PerformanceSensitive bool     `json:"performance_sensitive"`
	CriticalFactors      []string `json:"critical_factors"`
	EvidenceRefs         []string `json:"evidence_refs"`
	BenchmarkRefs        []string `json:"benchmark_refs"`
	AdversarialRefs      []string `json:"adversarial_refs"`
	SupersededBy         string   `json:"superseded_by"`
}

type DecisionRegistry struct {
	SchemaVersion int                   `json:"schema_version"`
	RegistryKind  string                `json:"registry_kind"`
	Decisions     []EngineeringDecision `json:"decisions"`
}

type DecisionRegistryAuditSummary struct {
	Decisions            int            `json:"decisions"`
	ByClass              map[string]int `json:"by_class"`
	ByStatus             map[string]int `json:"by_status"`
	BlockedProduct       int            `json:"blocked_product"`
	Accepted             int            `json:"accepted"`
	PerformanceSensitive int            `json:"performance_sensitive"`
	Critical             int            `json:"critical"`
}

var decisionIDPattern = regexp.MustCompile(`^ENG-DEC-[0-9]{4,}$`)
var decisionKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{2,127}$`)

var knownDecisionStatuses = map[string]bool{
	"proposed": true, "researching": true, "accepted": true, "rejected": true,
	"blocked-product": true, "revisit-required": true, "superseded": true,
}
var knownProductDecisionBoundaries = map[string]bool{
	"engineering-only":           true,
	"product-open-decision":      true,
	"product-behavior":           true,
	"product-security-invariant": true,
}
var knownDecisionRiskTags = map[string]bool{
	"local-refactor": true, "representation": true, "serialization": true, "cache": true,
	"dependency": true, "storage": true, "indexing": true, "performance": true,
	"concurrency": true, "architecture": true, "security": true, "algorithm": true,
	"durability": true, "reliability": true,
}
var knownCriticalFactors = map[string]bool{
	"security-critical":      true,
	"architecture-critical":  true,
	"algorithm-critical":     true,
	"durability-critical":    true,
	"high-scale-concurrency": true,
}

func runDecisionRegistryAudit(root string, graph WorkGraph) (DecisionRegistryAuditSummary, error) {
	var policy DecisionPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionPolicyPath)), &policy); err != nil {
		return DecisionRegistryAuditSummary{}, err
	}
	if err := validateDecisionPolicy(policy); err != nil {
		return DecisionRegistryAuditSummary{}, err
	}
	var registry DecisionRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionRegistryPath)), &registry); err != nil {
		return DecisionRegistryAuditSummary{}, err
	}
	if err := validateDecisionRegistry(policy, registry, graph); err != nil {
		return DecisionRegistryAuditSummary{}, err
	}
	summary := DecisionRegistryAuditSummary{
		Decisions: len(registry.Decisions), ByClass: map[string]int{}, ByStatus: map[string]int{},
	}
	for _, d := range registry.Decisions {
		summary.ByClass[d.Class]++
		summary.ByStatus[d.Status]++
		if d.Status == "blocked-product" {
			summary.BlockedProduct++
		}
		if d.Status == "accepted" {
			summary.Accepted++
		}
		if d.PerformanceSensitive {
			summary.PerformanceSensitive++
		}
		if d.Class == "C" {
			summary.Critical++
		}
	}
	return summary, nil
}

func validateDecisionPolicy(policy DecisionPolicy) error {
	if policy.SchemaVersion != 1 {
		return fmt.Errorf("unsupported decision policy schema_version %d", policy.SchemaVersion)
	}
	if policy.DefaultPolicy != "deny-unclassified-significant-decision" {
		return fmt.Errorf("decision policy default must be deny-unclassified-significant-decision")
	}
	if err := exactStringSet("decision status", policy.Statuses, knownDecisionStatuses); err != nil {
		return err
	}
	if err := exactStringSet("product decision boundary", policy.ProductBoundaries, knownProductDecisionBoundaries); err != nil {
		return err
	}
	if err := exactStringSet("decision risk tag", policy.RiskTags, knownDecisionRiskTags); err != nil {
		return err
	}
	expected := map[string]DecisionClassPolicy{
		"A": {ID: "A", RecordRequired: false, CurrentExternalResearch: "not-required", AlternativesRequired: false, BenchmarkRule: "not-required", AdversarialReviewRequired: false},
		"B": {ID: "B", RecordRequired: true, CurrentExternalResearch: "if-material-uncertainty", AlternativesRequired: true, BenchmarkRule: "if-performance-sensitive", AdversarialReviewRequired: false},
		"C": {ID: "C", RecordRequired: true, CurrentExternalResearch: "required", AlternativesRequired: true, BenchmarkRule: "if-performance-sensitive", AdversarialReviewRequired: true},
	}
	if len(policy.Classes) != len(expected) {
		return fmt.Errorf("decision policy must define exactly classes A, B and C")
	}
	seen := map[string]bool{}
	for _, class := range policy.Classes {
		want, ok := expected[class.ID]
		if !ok || seen[class.ID] {
			return fmt.Errorf("invalid or duplicate decision class %q", class.ID)
		}
		seen[class.ID] = true
		if class != want {
			return fmt.Errorf("decision class %s policy differs from canonical autonomous decision policy", class.ID)
		}
	}
	return nil
}

func validateDecisionRegistry(policy DecisionPolicy, registry DecisionRegistry, graph WorkGraph) error {
	if registry.SchemaVersion != 1 {
		return fmt.Errorf("unsupported decision registry schema_version %d", registry.SchemaVersion)
	}
	if registry.RegistryKind != "engineering-technical-decisions" {
		return fmt.Errorf("unexpected decision registry kind %q", registry.RegistryKind)
	}
	classPolicy := map[string]DecisionClassPolicy{}
	for _, class := range policy.Classes {
		classPolicy[class.ID] = class
	}
	workUnits := map[string]bool{}
	for _, node := range graph.Nodes {
		workUnits[node.ID] = true
	}

	byID := map[string]EngineeringDecision{}
	liveByKey := map[string]string{}
	for _, d := range registry.Decisions {
		if !decisionIDPattern.MatchString(d.ID) {
			return fmt.Errorf("invalid decision id %q", d.ID)
		}
		if _, duplicate := byID[d.ID]; duplicate {
			return fmt.Errorf("duplicate decision id %s", d.ID)
		}
		if !decisionKeyPattern.MatchString(d.DecisionKey) {
			return fmt.Errorf("decision %s has invalid decision_key %q", d.ID, d.DecisionKey)
		}
		if _, ok := classPolicy[d.Class]; !ok {
			return fmt.Errorf("decision %s has unknown class %q", d.ID, d.Class)
		}
		if !knownDecisionStatuses[d.Status] {
			return fmt.Errorf("decision %s has unknown status %q", d.ID, d.Status)
		}
		if !knownProductDecisionBoundaries[d.ProductBoundary] {
			return fmt.Errorf("decision %s has unknown product boundary %q", d.ID, d.ProductBoundary)
		}
		if strings.TrimSpace(d.Title) == "" || strings.TrimSpace(d.Owner) == "" ||
			strings.TrimSpace(d.Question) == "" || strings.TrimSpace(d.Rationale) == "" {
			return fmt.Errorf("decision %s requires title owner question and rationale", d.ID)
		}
		if len(d.AffectedWorkUnits) == 0 || len(d.Constraints) == 0 {
			return fmt.Errorf("decision %s requires affected_work_units and constraints", d.ID)
		}
		if err := uniqueNonEmptyStrings(d.ID+" affected_work_units", d.AffectedWorkUnits); err != nil {
			return err
		}
		for _, workUnit := range d.AffectedWorkUnits {
			if !workUnits[workUnit] {
				return fmt.Errorf("decision %s references unknown work unit %s", d.ID, workUnit)
			}
		}
		if err := uniqueNonEmptyStrings(d.ID+" constraints", d.Constraints); err != nil {
			return err
		}
		if err := validateKnownUniqueStrings(d.ID+" risk tag", d.RiskTags, knownDecisionRiskTags); err != nil {
			return err
		}
		if err := validateKnownUniqueStrings(d.ID+" critical factor", d.CriticalFactors, knownCriticalFactors); err != nil {
			return err
		}
		if len(d.CriticalFactors) > 0 && d.Class != "C" {
			return fmt.Errorf("decision %s has critical factors but is not Class C", d.ID)
		}
		if d.PerformanceSensitive && d.Class == "A" {
			return fmt.Errorf("performance-sensitive decision %s cannot be Class A", d.ID)
		}
		if d.ProductBoundary != "engineering-only" {
			switch d.Status {
			case "blocked-product", "rejected", "superseded":
			default:
				return fmt.Errorf("decision %s crosses product boundary and must be blocked-product, rejected or superseded", d.ID)
			}
		}
		if d.Status == "accepted" {
			if d.ProductBoundary != "engineering-only" {
				return fmt.Errorf("decision %s cannot be accepted outside engineering-only boundary", d.ID)
			}
			if d.Class != "A" && len(d.EvidenceRefs) == 0 {
				return fmt.Errorf("accepted Class %s decision %s requires evidence_refs", d.Class, d.ID)
			}
			if d.PerformanceSensitive && len(d.BenchmarkRefs) == 0 {
				return fmt.Errorf("accepted performance-sensitive decision %s requires benchmark_refs", d.ID)
			}
			if d.Class == "C" && len(d.AdversarialRefs) == 0 {
				return fmt.Errorf("accepted Class C decision %s requires adversarial_refs", d.ID)
			}
		}
		for label, refs := range map[string][]string{
			"evidence_refs": d.EvidenceRefs, "benchmark_refs": d.BenchmarkRefs, "adversarial_refs": d.AdversarialRefs,
		} {
			if err := uniqueNonEmptyStrings(d.ID+" "+label, refs); err != nil {
				return err
			}
		}
		if d.Status == "superseded" {
			if strings.TrimSpace(d.SupersededBy) == "" {
				return fmt.Errorf("superseded decision %s requires superseded_by", d.ID)
			}
		} else if d.SupersededBy != "" {
			return fmt.Errorf("decision %s has superseded_by but status is %s", d.ID, d.Status)
		}
		if decisionStatusIsLive(d.Status) {
			if previous, exists := liveByKey[d.DecisionKey]; exists {
				return fmt.Errorf("decisions %s and %s are both live for decision_key %s", previous, d.ID, d.DecisionKey)
			}
			liveByKey[d.DecisionKey] = d.ID
		}
		byID[d.ID] = d
	}
	for _, d := range registry.Decisions {
		if d.Status != "superseded" {
			continue
		}
		next, ok := byID[d.SupersededBy]
		if !ok {
			return fmt.Errorf("superseded decision %s references unknown decision %s", d.ID, d.SupersededBy)
		}
		if next.DecisionKey != d.DecisionKey {
			return fmt.Errorf("superseded decision %s and replacement %s do not share decision_key", d.ID, next.ID)
		}
		if next.Status == "superseded" {
			return fmt.Errorf("superseded decision %s points to another superseded decision %s", d.ID, next.ID)
		}
	}
	return nil
}

func decisionStatusIsLive(status string) bool {
	switch status {
	case "proposed", "researching", "accepted", "blocked-product", "revisit-required":
		return true
	default:
		return false
	}
}

func exactStringSet(label string, values []string, expected map[string]bool) error {
	if len(values) != len(expected) {
		return fmt.Errorf("%s enum mismatch: expected %d values got %d", label, len(expected), len(values))
	}
	return validateKnownUniqueStrings(label, values, expected)
}

func validateKnownUniqueStrings(label string, values []string, known map[string]bool) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !known[value] {
			return fmt.Errorf("%s contains unknown value %q", label, value)
		}
		if seen[value] {
			return fmt.Errorf("%s duplicates value %q", label, value)
		}
		seen[value] = true
	}
	return nil
}

func uniqueNonEmptyStrings(label string, values []string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s contains an empty value", label)
		}
		if seen[value] {
			return fmt.Errorf("%s duplicates %q", label, value)
		}
		seen[value] = true
	}
	return nil
}

func sortedDecisionIDs(registry DecisionRegistry) []string {
	ids := make([]string, 0, len(registry.Decisions))
	for _, d := range registry.Decisions {
		ids = append(ids, d.ID)
	}
	sort.Strings(ids)
	return ids
}
