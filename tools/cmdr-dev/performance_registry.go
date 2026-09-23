package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	performancePolicyPath   = "engineering/performance/performance-policy.json"
	performanceRegistryPath = "engineering/performance/performance-registry.json"
)

type PerformancePolicy struct {
	SchemaVersion     int                 `json:"schema_version"`
	DefaultPolicy     string              `json:"default_policy"`
	ScopeKinds        []string            `json:"scope_kinds"`
	MetricKinds       []string            `json:"metric_kinds"`
	MetricUnits       map[string][]string `json:"metric_units"`
	Aggregations      []string            `json:"aggregations"`
	Comparators       []string            `json:"comparators"`
	EnvironmentKinds []string            `json:"environment_kinds"`
	Stages            []string            `json:"stages"`
	WorkloadKeys      []string            `json:"workload_keys"`
}

type PerformanceEnvironment struct {
	ID                        string `json:"id"`
	Kind                      string `json:"kind"`
	OS                        string `json:"os"`
	Arch                      string `json:"arch"`
	RunnerClass               string `json:"runner_class"`
	RelativeRegressionAllowed bool   `json:"relative_regression_allowed"`
	Rationale                 string `json:"rationale"`
}

type PerformanceMetricBudget struct {
	ID                           string  `json:"id"`
	Kind                         string  `json:"kind"`
	Unit                         string  `json:"unit"`
	Aggregation                  string  `json:"aggregation"`
	Comparator                   string  `json:"comparator"`
	Budget                       float64 `json:"budget"`
	MaxRelativeRegressionPercent float64 `json:"max_relative_regression_percent"`
}

type PerformanceTarget struct {
	ID                  string                    `json:"id"`
	ScopeKind           string                    `json:"scope_kind"`
	BoundaryID          string                    `json:"boundary_id"`
	OwnerPath           string                    `json:"owner_path"`
	WorkloadKey         string                    `json:"workload_key"`
	WorkloadDescription string                    `json:"workload_description"`
	EnvironmentID       string                    `json:"environment_id"`
	Stages              []string                  `json:"stages"`
	TriggerPaths        []string                  `json:"trigger_paths"`
	Metrics             []PerformanceMetricBudget `json:"metrics"`
	Rationale           string                    `json:"rationale"`
}

type PerformanceRegistry struct {
	SchemaVersion int                      `json:"schema_version"`
	RegistryKind  string                   `json:"registry_kind"`
	Environments  []PerformanceEnvironment `json:"environments"`
	Targets       []PerformanceTarget      `json:"targets"`
}

type PerformanceRegistryAuditSummary struct {
	Environments             int            `json:"environments"`
	Targets                  int            `json:"targets"`
	Metrics                  int            `json:"metrics"`
	RuntimeBoundaries        int            `json:"runtime_boundaries"`
	RuntimeTargetedBoundaries int            `json:"runtime_targeted_boundaries"`
	RuntimeCoverageStatus    string         `json:"runtime_coverage_status"`
	ByScope                  map[string]int `json:"by_scope"`
	ByStage                  map[string]int `json:"by_stage"`
}

var performanceEnvironmentIDPattern = regexp.MustCompile(`^PERF-ENV-[A-Z0-9-]+$`)
var performanceTargetIDPattern = regexp.MustCompile(`^PERF-TGT-[A-Z0-9-]+$`)
var performanceMetricIDPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,63}$`)

var knownPerformanceScopeKinds = map[string]bool{
	"engineering-control-plane": true,
	"product-runtime":           true,
}

var knownPerformanceMetricUnits = map[string][]string{
	"latency":     {"ns", "us", "ms", "s"},
	"throughput":  {"ops/s", "events/s", "bytes/s"},
	"cpu-time":    {"ns/op", "us/op", "ms/op"},
	"memory":      {"bytes", "MiB"},
	"allocation":  {"bytes/op", "allocs/op"},
	"utilization": {"percent"},
}

var knownPerformanceAggregations = map[string]bool{
	"median": true,
	"p95":    true,
	"p99":    true,
	"max":    true,
	"mean":   true,
}

var knownPerformanceComparators = map[string]bool{
	"upper-bound": true,
	"lower-bound": true,
}

var knownPerformanceEnvironmentKinds = map[string]bool{
	"shared-ci":        true,
	"dedicated":        true,
	"local-calibrated": true,
}

var knownPerformanceStages = map[string]bool{
	"pr":        true,
	"nightly":   true,
	"release":   true,
	"on-demand": true,
}

var knownPerformanceWorkloadKeys = map[string]bool{
	"cmdr-dev-metadata-audit-v1": true,
}

func runPerformanceRegistryAudit(root string) (PerformanceRegistryAuditSummary, error) {
	var policy PerformancePolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(performancePolicyPath)), &policy); err != nil {
		return PerformanceRegistryAuditSummary{}, err
	}
	if err := validatePerformancePolicy(policy); err != nil {
		return PerformanceRegistryAuditSummary{}, err
	}

	var registry PerformanceRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(performanceRegistryPath)), &registry); err != nil {
		return PerformanceRegistryAuditSummary{}, err
	}

	var architecture ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &architecture); err != nil {
		return PerformanceRegistryAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(architecture); err != nil {
		return PerformanceRegistryAuditSummary{}, err
	}
	return validatePerformanceRegistry(policy, registry, architecture)
}

func validatePerformancePolicy(policy PerformancePolicy) error {
	if policy.SchemaVersion != 1 || policy.DefaultPolicy != "deny-unregistered-performance-target" {
		return fmt.Errorf("invalid performance policy header")
	}
	if err := exactStringSet("performance scope kind", policy.ScopeKinds, knownPerformanceScopeKinds); err != nil {
		return err
	}
	metricKinds := map[string]bool{}
	for kind := range knownPerformanceMetricUnits {
		metricKinds[kind] = true
	}
	if err := exactStringSet("performance metric kind", policy.MetricKinds, metricKinds); err != nil {
		return err
	}
	if len(policy.MetricUnits) != len(knownPerformanceMetricUnits) {
		return fmt.Errorf("performance metric_units set mismatch")
	}
	for kind, knownUnits := range knownPerformanceMetricUnits {
		got, ok := policy.MetricUnits[kind]
		if !ok {
			return fmt.Errorf("performance policy missing metric units for %s", kind)
		}
		expected := map[string]bool{}
		for _, unit := range knownUnits {
			expected[unit] = true
		}
		if err := exactStringSet("performance "+kind+" unit", got, expected); err != nil {
			return err
		}
	}
	if err := exactStringSet("performance aggregation", policy.Aggregations, knownPerformanceAggregations); err != nil {
		return err
	}
	if err := exactStringSet("performance comparator", policy.Comparators, knownPerformanceComparators); err != nil {
		return err
	}
	if err := exactStringSet("performance environment kind", policy.EnvironmentKinds, knownPerformanceEnvironmentKinds); err != nil {
		return err
	}
	if err := exactStringSet("performance stage", policy.Stages, knownPerformanceStages); err != nil {
		return err
	}
	return exactStringSet("performance workload key", policy.WorkloadKeys, knownPerformanceWorkloadKeys)
}

func validatePerformanceRegistry(policy PerformancePolicy, registry PerformanceRegistry, architecture ArchitectureRegistry) (PerformanceRegistryAuditSummary, error) {
	summary := PerformanceRegistryAuditSummary{
		ByScope: map[string]int{},
		ByStage: map[string]int{},
	}
	if registry.SchemaVersion != 1 || registry.RegistryKind != "performance-targets" {
		return summary, fmt.Errorf("invalid performance registry header")
	}

	boundaries := map[string]ArchitectureBoundary{}
	runtimeBoundaries := map[string]bool{}
	for _, boundary := range architecture.Boundaries {
		boundaries[boundary.ID] = boundary
		if boundary.Kind == "product-runtime" {
			runtimeBoundaries[boundary.ID] = true
		}
	}
	summary.RuntimeBoundaries = len(runtimeBoundaries)

	envByID := map[string]PerformanceEnvironment{}
	for _, env := range registry.Environments {
		if !performanceEnvironmentIDPattern.MatchString(env.ID) {
			return summary, fmt.Errorf("invalid performance environment id %q", env.ID)
		}
		if _, exists := envByID[env.ID]; exists {
			return summary, fmt.Errorf("duplicate performance environment %s", env.ID)
		}
		if !knownPerformanceEnvironmentKinds[env.Kind] {
			return summary, fmt.Errorf("performance environment %s uses unknown kind %s", env.ID, env.Kind)
		}
		if strings.TrimSpace(env.OS) == "" || strings.TrimSpace(env.Arch) == "" ||
			strings.TrimSpace(env.RunnerClass) == "" || strings.TrimSpace(env.Rationale) == "" {
			return summary, fmt.Errorf("performance environment %s has incomplete identity", env.ID)
		}
		envByID[env.ID] = env
	}
	summary.Environments = len(envByID)

	targetIDs := map[string]bool{}
	runtimeTargeted := map[string]bool{}
	for _, target := range registry.Targets {
		if !performanceTargetIDPattern.MatchString(target.ID) || targetIDs[target.ID] {
			return summary, fmt.Errorf("invalid or duplicate performance target id %q", target.ID)
		}
		targetIDs[target.ID] = true
		if !knownPerformanceScopeKinds[target.ScopeKind] {
			return summary, fmt.Errorf("performance target %s uses unknown scope %s", target.ID, target.ScopeKind)
		}
		boundary, ok := boundaries[target.BoundaryID]
		if !ok {
			return summary, fmt.Errorf("performance target %s references unknown boundary %s", target.ID, target.BoundaryID)
		}
		switch target.ScopeKind {
		case "product-runtime":
			if boundary.Kind != "product-runtime" {
				return summary, fmt.Errorf("performance target %s requires a product-runtime boundary", target.ID)
			}
			runtimeTargeted[boundary.ID] = true
		case "engineering-control-plane":
			if boundary.Kind != "engineering" {
				return summary, fmt.Errorf("performance target %s requires an engineering boundary", target.ID)
			}
		}
		owner, err := normalizeArchitecturePattern(target.OwnerPath)
		if err != nil {
			return summary, fmt.Errorf("performance target %s owner_path: %w", target.ID, err)
		}
		owned := false
		for _, root := range boundary.Roots {
			if patternCovers(root, owner) {
				owned = true
				break
			}
		}
		if !owned {
			return summary, fmt.Errorf("performance target %s owner_path %s is outside boundary %s", target.ID, target.OwnerPath, boundary.ID)
		}
		if !knownPerformanceWorkloadKeys[target.WorkloadKey] || strings.TrimSpace(target.WorkloadDescription) == "" {
			return summary, fmt.Errorf("performance target %s has invalid workload", target.ID)
		}
		env, ok := envByID[target.EnvironmentID]
		if !ok {
			return summary, fmt.Errorf("performance target %s references unknown environment %s", target.ID, target.EnvironmentID)
		}
		if strings.TrimSpace(target.Rationale) == "" {
			return summary, fmt.Errorf("performance target %s requires rationale", target.ID)
		}
		if err := validateKnownUniqueStrings(target.ID+" stage", target.Stages, knownPerformanceStages); err != nil {
			return summary, err
		}
		if len(target.Stages) == 0 {
			return summary, fmt.Errorf("performance target %s requires at least one stage", target.ID)
		}
		for _, stage := range target.Stages {
			summary.ByStage[stage]++
		}
		if len(target.TriggerPaths) == 0 {
			return summary, fmt.Errorf("performance target %s requires trigger_paths", target.ID)
		}
		seenTriggers := map[string]bool{}
		for _, trigger := range target.TriggerPaths {
			if _, err := normalizeArchitecturePattern(trigger); err != nil {
				return summary, fmt.Errorf("performance target %s trigger: %w", target.ID, err)
			}
			if seenTriggers[trigger] {
				return summary, fmt.Errorf("performance target %s duplicates trigger %s", target.ID, trigger)
			}
			seenTriggers[trigger] = true
		}
		if len(target.Metrics) == 0 {
			return summary, fmt.Errorf("performance target %s requires metrics", target.ID)
		}
		metricIDs := map[string]bool{}
		for _, metric := range target.Metrics {
			if err := validatePerformanceMetric(target.ID, metric, env); err != nil {
				return summary, err
			}
			if !performanceMetricIDPattern.MatchString(metric.ID) || metricIDs[metric.ID] {
				return summary, fmt.Errorf("performance target %s has invalid/duplicate metric id %q", target.ID, metric.ID)
			}
			metricIDs[metric.ID] = true
			summary.Metrics++
		}
		summary.Targets++
		summary.ByScope[target.ScopeKind]++
	}

	summary.RuntimeTargetedBoundaries = len(runtimeTargeted)
	if len(runtimeBoundaries) == 0 {
		summary.RuntimeCoverageStatus = "not-applicable"
		return summary, nil
	}
	var missing []string
	for id := range runtimeBoundaries {
		if !runtimeTargeted[id] {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		summary.RuntimeCoverageStatus = "incomplete"
		return summary, fmt.Errorf("product-runtime boundaries missing performance targets: %s", strings.Join(missing, ", "))
	}
	summary.RuntimeCoverageStatus = "complete"
	return summary, nil
}

func validatePerformanceMetric(targetID string, metric PerformanceMetricBudget, env PerformanceEnvironment) error {
	units, ok := knownPerformanceMetricUnits[metric.Kind]
	if !ok || !containsString(units, metric.Unit) {
		return fmt.Errorf("performance target %s metric %s has invalid kind/unit %s/%s", targetID, metric.ID, metric.Kind, metric.Unit)
	}
	if !knownPerformanceAggregations[metric.Aggregation] {
		return fmt.Errorf("performance target %s metric %s has unknown aggregation %s", targetID, metric.ID, metric.Aggregation)
	}
	if !knownPerformanceComparators[metric.Comparator] {
		return fmt.Errorf("performance target %s metric %s has unknown comparator %s", targetID, metric.ID, metric.Comparator)
	}
	if metric.Kind == "throughput" && metric.Comparator != "lower-bound" {
		return fmt.Errorf("performance target %s throughput metric %s must use lower-bound", targetID, metric.ID)
	}
	if metric.Kind != "throughput" && metric.Comparator != "upper-bound" {
		return fmt.Errorf("performance target %s metric %s must use upper-bound", targetID, metric.ID)
	}
	if metric.Budget <= 0 || math.IsNaN(metric.Budget) || math.IsInf(metric.Budget, 0) {
		return fmt.Errorf("performance target %s metric %s has invalid budget", targetID, metric.ID)
	}
	if metric.MaxRelativeRegressionPercent < 0 || metric.MaxRelativeRegressionPercent > 100 ||
		math.IsNaN(metric.MaxRelativeRegressionPercent) || math.IsInf(metric.MaxRelativeRegressionPercent, 0) {
		return fmt.Errorf("performance target %s metric %s has invalid relative regression limit", targetID, metric.ID)
	}
	if metric.MaxRelativeRegressionPercent > 0 && !env.RelativeRegressionAllowed {
		return fmt.Errorf("performance target %s metric %s cannot use relative regression in environment %s", targetID, metric.ID, env.ID)
	}
	return nil
}
