package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const engineMetricsRegistryPath = "engineering/metrics/metrics-registry.json"

type EngineMetricDefinition struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Category             string `json:"category"`
	Unit                 string `json:"unit"`
	Aggregation          string `json:"aggregation"`
	SourceKey            string `json:"source_key"`
	SafetyClass          string `json:"safety_class"`
	OptimizationEligible bool   `json:"optimization_eligible"`
	Direction            string `json:"direction"`
	Rationale            string `json:"rationale"`
}

type EngineMetricsRegistry struct {
	SchemaVersion int                      `json:"schema_version"`
	RegistryKind  string                   `json:"registry_kind"`
	DefaultPolicy string                   `json:"default_policy"`
	Categories    []string                 `json:"categories"`
	Units         []string                 `json:"units"`
	Aggregations  []string                 `json:"aggregations"`
	SourceKeys    []string                 `json:"source_keys"`
	SafetyClasses []string                 `json:"safety_classes"`
	Directions    []string                 `json:"directions"`
	Metrics       []EngineMetricDefinition `json:"metrics"`
}

type EngineMetricObservation struct {
	MetricID   string  `json:"metric_id"`
	Value      float64 `json:"value"`
	SourceKey  string  `json:"source_key"`
	SourceSHA  string  `json:"source_sha"`
	ObservedAt string  `json:"observed_at"`
	Evidence   string  `json:"evidence"`
}

type EngineMetricsRegistrySummary struct {
	Metrics              int            `json:"metrics"`
	OptimizationEligible int            `json:"optimization_eligible"`
	NonOptimizable       int            `json:"non_optimizable"`
	ByCategory           map[string]int `json:"by_category"`
	BySafetyClass        map[string]int `json:"by_safety_class"`
	BySource             map[string]int `json:"by_source"`
}

var engineMetricIDPattern = regexp.MustCompile("^MET-[A-Z0-9][A-Z0-9-]{2,63}$")

var knownEngineMetricCategories = map[string]bool{
	"validation": true,
	"context":    true,
	"cache":      true,
	"recovery":   true,
	"quality":    true,
}

var knownEngineMetricUnits = map[string]bool{
	"count":     true,
	"cost-unit": true,
	"percent":   true,
	"boolean":   true,
}

var knownEngineMetricAggregations = map[string]bool{
	"last":  true,
	"sum":   true,
	"max":   true,
	"ratio": true,
}

var knownEngineMetricSourceKeys = map[string]bool{
	"validation-plan":      true,
	"validation-execution": true,
	"check-catalog":        true,
	"context-bundle":       true,
	"performance-cache":    true,
	"recovery-reconcile":   true,
	"coordination":         true,
}

var knownEngineMetricSafetyClasses = map[string]bool{
	"safety-floor":  true,
	"quality-floor": true,
	"efficiency":    true,
}

var knownEngineMetricDirections = map[string]bool{
	"lower-better":  true,
	"higher-better": true,
	"invariant":     true,
}

func runEngineMetricsRegistryAudit(root string) (EngineMetricsRegistrySummary, error) {
	var registry EngineMetricsRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineMetricsRegistryPath)), &registry); err != nil {
		return EngineMetricsRegistrySummary{}, err
	}
	return validateEngineMetricsRegistry(registry)
}

func validateEngineMetricsRegistry(registry EngineMetricsRegistry) (EngineMetricsRegistrySummary, error) {
	summary := EngineMetricsRegistrySummary{
		ByCategory:    map[string]int{},
		BySafetyClass: map[string]int{},
		BySource:      map[string]int{},
	}
	if registry.SchemaVersion != 1 || registry.RegistryKind != "engineering-engine-metrics" ||
		registry.DefaultPolicy != "deny-unregistered-metric" {
		return summary, fmt.Errorf("invalid engine metrics registry header")
	}
	if err := exactStringSet("engine metric category", registry.Categories, knownEngineMetricCategories); err != nil {
		return summary, err
	}
	if err := exactStringSet("engine metric unit", registry.Units, knownEngineMetricUnits); err != nil {
		return summary, err
	}
	if err := exactStringSet("engine metric aggregation", registry.Aggregations, knownEngineMetricAggregations); err != nil {
		return summary, err
	}
	if err := exactStringSet("engine metric source", registry.SourceKeys, knownEngineMetricSourceKeys); err != nil {
		return summary, err
	}
	if err := exactStringSet("engine metric safety class", registry.SafetyClasses, knownEngineMetricSafetyClasses); err != nil {
		return summary, err
	}
	if err := exactStringSet("engine metric direction", registry.Directions, knownEngineMetricDirections); err != nil {
		return summary, err
	}
	if len(registry.Metrics) == 0 {
		return summary, fmt.Errorf("engine metrics registry is empty")
	}

	ids := map[string]bool{}
	for _, metric := range registry.Metrics {
		if !engineMetricIDPattern.MatchString(metric.ID) || ids[metric.ID] {
			return summary, fmt.Errorf("invalid or duplicate engine metric id %q", metric.ID)
		}
		ids[metric.ID] = true
		if strings.TrimSpace(metric.Title) == "" || strings.TrimSpace(metric.Rationale) == "" {
			return summary, fmt.Errorf("engine metric %s requires title and rationale", metric.ID)
		}
		if !knownEngineMetricCategories[metric.Category] ||
			!knownEngineMetricUnits[metric.Unit] ||
			!knownEngineMetricAggregations[metric.Aggregation] ||
			!knownEngineMetricSourceKeys[metric.SourceKey] ||
			!knownEngineMetricSafetyClasses[metric.SafetyClass] ||
			!knownEngineMetricDirections[metric.Direction] {
			return summary, fmt.Errorf("engine metric %s uses an unknown closed-enum value", metric.ID)
		}
		if metric.SafetyClass != "efficiency" && metric.OptimizationEligible {
			return summary, fmt.Errorf("engine metric %s cannot optimize a %s metric", metric.ID, metric.SafetyClass)
		}
		if metric.SafetyClass == "safety-floor" && metric.Direction != "invariant" {
			return summary, fmt.Errorf("safety-floor metric %s must use invariant direction", metric.ID)
		}
		if metric.Unit == "percent" && metric.Aggregation != "last" && metric.Aggregation != "ratio" {
			return summary, fmt.Errorf("percentage metric %s must use last or ratio aggregation", metric.ID)
		}
		summary.Metrics++
		if metric.OptimizationEligible {
			summary.OptimizationEligible++
		} else {
			summary.NonOptimizable++
		}
		summary.ByCategory[metric.Category]++
		summary.BySafetyClass[metric.SafetyClass]++
		summary.BySource[metric.SourceKey]++
	}
	return summary, nil
}

func validateEngineMetricObservation(registry EngineMetricsRegistry, observation EngineMetricObservation) error {
	var definition *EngineMetricDefinition
	for i := range registry.Metrics {
		if registry.Metrics[i].ID == observation.MetricID {
			value := registry.Metrics[i]
			definition = &value
			break
		}
	}
	if definition == nil {
		return fmt.Errorf("observation references unknown metric %s", observation.MetricID)
	}
	if observation.SourceKey != definition.SourceKey {
		return fmt.Errorf("observation %s source mismatch: expected %s got %s", observation.MetricID, definition.SourceKey, observation.SourceKey)
	}
	if math.IsNaN(observation.Value) || math.IsInf(observation.Value, 0) {
		return fmt.Errorf("observation %s value must be finite", observation.MetricID)
	}
	if definition.Unit == "percent" && (observation.Value < 0 || observation.Value > 100) {
		return fmt.Errorf("observation %s percentage is outside 0..100", observation.MetricID)
	}
	if definition.Unit == "boolean" && observation.Value != 0 && observation.Value != 1 {
		return fmt.Errorf("observation %s boolean value must be 0 or 1", observation.MetricID)
	}
	if observation.Value < 0 && definition.Unit != "boolean" {
		return fmt.Errorf("observation %s cannot be negative", observation.MetricID)
	}
	sha, err := validateFullCommitID(observation.SourceSHA)
	if err != nil || sha != observation.SourceSHA {
		return fmt.Errorf("observation %s requires canonical lowercase full source SHA", observation.MetricID)
	}
	if _, err := parseCanonicalLeaseTime(observation.ObservedAt); err != nil {
		return fmt.Errorf("observation %s observed_at: %w", observation.MetricID, err)
	}
	if strings.TrimSpace(observation.Evidence) == "" || len(observation.Evidence) > 256 {
		return fmt.Errorf("observation %s requires a bounded evidence reference", observation.MetricID)
	}
	return nil
}

func engineMetricIDs(registry EngineMetricsRegistry) []string {
	out := make([]string, 0, len(registry.Metrics))
	for _, metric := range registry.Metrics {
		out = append(out, metric.ID)
	}
	sort.Strings(out)
	return out
}
