package main

import (
	"fmt"
	"path/filepath"
	"sort"
)

type EngineOptimizationBenefit struct {
	MetricID string  `json:"metric_id"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit"`
}

type EngineOptimizationRecommendation struct {
	ID               string                      `json:"id"`
	Kind             string                      `json:"kind"`
	Advisory         bool                        `json:"advisory"`
	MetricIDs        []string                    `json:"metric_ids"`
	EstimatedBenefit []EngineOptimizationBenefit `json:"estimated_benefit"`
	Rationale        string                      `json:"rationale"`
	Constraints      []string                    `json:"constraints"`
	RevisitTrigger   string                      `json:"revisit_trigger"`
}

type EngineOptimizationWithheld struct {
	OpportunityID string   `json:"opportunity_id"`
	MetricIDs     []string `json:"metric_ids"`
	Reason        string   `json:"reason"`
}

type EngineOptimizationSafetyProof struct {
	CurrentSnapshotDigest       string `json:"current_snapshot_digest"`
	BaselineSnapshotDigest      string `json:"baseline_snapshot_digest,omitempty"`
	SafetyGreen                 bool   `json:"safety_green"`
	CoordinationConflictFree    bool   `json:"coordination_conflict_free"`
	RecoveryRevalidationClear   bool   `json:"recovery_revalidation_clear"`
	MandatoryFloorKnown         bool   `json:"mandatory_floor_known"`
	DirectPlanMutationAllowed   bool   `json:"direct_plan_mutation_allowed"`
	MandatoryMutationAllowed    bool   `json:"mandatory_mutation_allowed"`
	PrerequisiteMutationAllowed bool   `json:"prerequisite_mutation_allowed"`
	ProductSpecMutationAllowed  bool   `json:"product_spec_mutation_allowed"`
}

type EngineOptimizationReport struct {
	SchemaVersion   int                                `json:"schema_version"`
	ReportKind      string                             `json:"report_kind"`
	WorkUnit        string                             `json:"work_unit"`
	SourceSHA       string                             `json:"source_sha"`
	SnapshotDigest  string                             `json:"snapshot_digest"`
	BaselineDigest  string                             `json:"baseline_digest,omitempty"`
	SafetyStatus    string                             `json:"safety_status"`
	Recommendations []EngineOptimizationRecommendation `json:"recommendations"`
	Withheld        []EngineOptimizationWithheld       `json:"withheld,omitempty"`
	Proof           EngineOptimizationSafetyProof      `json:"proof"`
	Digest          string                             `json:"digest_sha256"`
}

type engineOptimizationReportBody struct {
	SchemaVersion   int                                `json:"schema_version"`
	ReportKind      string                             `json:"report_kind"`
	WorkUnit        string                             `json:"work_unit"`
	SourceSHA       string                             `json:"source_sha"`
	SnapshotDigest  string                             `json:"snapshot_digest"`
	BaselineDigest  string                             `json:"baseline_digest,omitempty"`
	SafetyStatus    string                             `json:"safety_status"`
	Recommendations []EngineOptimizationRecommendation `json:"recommendations"`
	Withheld        []EngineOptimizationWithheld       `json:"withheld,omitempty"`
	Proof           EngineOptimizationSafetyProof      `json:"proof"`
}

var optimizationOpportunityMetrics = map[string][]string{
	"OPT-VALIDATION-SELECTIVITY": {"MET-VALIDATION-COST", "MET-VALIDATION-SELECTED", "MET-MANDATORY-CHECKS"},
	"OPT-CACHE-REUSE":            {"MET-CACHE-REUSE-RATE"},
	"OPT-CONTEXT-BOUNDING":       {"MET-CONTEXT-SOURCES", "MET-CONTEXT-DEPENDENCIES"},
}

func compileEngineOptimizationReport(root string, current EngineMetricSnapshot) (EngineOptimizationReport, error) {
	var registry EngineMetricsRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineMetricsRegistryPath)), &registry); err != nil {
		return EngineOptimizationReport{}, err
	}
	if _, err := validateEngineMetricsRegistry(registry); err != nil {
		return EngineOptimizationReport{}, err
	}
	return deriveEngineOptimizations(registry, current, nil)
}

func deriveEngineOptimizations(registry EngineMetricsRegistry, current EngineMetricSnapshot, baseline *EngineMetricSnapshot) (EngineOptimizationReport, error) {
	if err := validateOptimizationSnapshot(registry, current); err != nil {
		return EngineOptimizationReport{}, fmt.Errorf("current optimization snapshot: %w", err)
	}
	proof, safe, reason := optimizationSafetyProof(current)
	body := engineOptimizationReportBody{
		SchemaVersion:  1,
		ReportKind:     "engineering-engine-optimization",
		WorkUnit:       current.WorkUnit,
		SourceSHA:      current.SourceSHA,
		SnapshotDigest: current.Digest,
		SafetyStatus:   "safe",
		Proof:          proof,
	}
	if !safe {
		body.SafetyStatus = "withheld-unsafe"
		withholdAllOptimizationOpportunities(&body, reason)
		return finalizeEngineOptimizationReport(body), nil
	}
	if baseline == nil {
		body.SafetyStatus = "safe-no-baseline"
		withholdAllOptimizationOpportunities(&body, "verified comparison baseline is unavailable")
		return finalizeEngineOptimizationReport(body), nil
	}
	if err := validateOptimizationSnapshot(registry, *baseline); err != nil {
		return EngineOptimizationReport{}, fmt.Errorf("baseline optimization snapshot: %w", err)
	}
	body.BaselineDigest = baseline.Digest
	body.Proof.BaselineSnapshotDigest = baseline.Digest

	if baseline.WorkUnit != current.WorkUnit {
		body.SafetyStatus = "safe-incomparable"
		withholdAllOptimizationOpportunities(&body, "baseline belongs to a different work unit")
		return finalizeEngineOptimizationReport(body), nil
	}
	currentTime, _ := parseCanonicalLeaseTime(current.ObservedAt)
	baselineTime, _ := parseCanonicalLeaseTime(baseline.ObservedAt)
	if !baselineTime.Before(currentTime) {
		body.SafetyStatus = "safe-incomparable"
		withholdAllOptimizationOpportunities(&body, "baseline is not older than current snapshot")
		return finalizeEngineOptimizationReport(body), nil
	}
	_, baselineSafe, baselineReason := optimizationSafetyProof(*baseline)
	if !baselineSafe {
		body.SafetyStatus = "safe-incomparable"
		withholdAllOptimizationOpportunities(&body, "baseline safety proof failed: "+baselineReason)
		return finalizeEngineOptimizationReport(body), nil
	}

	currentValues, currentUnavailable := snapshotMetricMaps(current)
	baselineValues, baselineUnavailable := snapshotMetricMaps(*baseline)
	definitions := metricDefinitionMap(registry)

	deriveValidationOptimization(&body, definitions, current, *baseline, currentValues, baselineValues, currentUnavailable, baselineUnavailable)
	deriveCacheOptimization(&body, definitions, currentValues, baselineValues, currentUnavailable, baselineUnavailable)
	deriveContextOptimization(&body, definitions, currentValues, baselineValues, currentUnavailable, baselineUnavailable)
	sortOptimizationReport(&body)
	return finalizeEngineOptimizationReport(body), nil
}

func validateOptimizationSnapshot(registry EngineMetricsRegistry, snapshot EngineMetricSnapshot) error {
	if snapshot.SchemaVersion != 1 || snapshot.SnapshotKind != "engineering-engine-metrics" {
		return fmt.Errorf("unsupported metric snapshot header")
	}
	if _, err := validateFullCommitID(snapshot.SourceSHA); err != nil {
		return err
	}
	if _, err := parseCanonicalLeaseTime(snapshot.ObservedAt); err != nil {
		return err
	}
	for _, observation := range snapshot.Observations {
		if err := validateEngineMetricObservation(registry, observation); err != nil {
			return err
		}
		if observation.SourceSHA != snapshot.SourceSHA || observation.ObservedAt != snapshot.ObservedAt {
			return fmt.Errorf("metric %s does not bind snapshot source identity", observation.MetricID)
		}
	}
	if err := validateEngineMetricSnapshotCoverage(registry, snapshot.Observations, snapshot.Unavailable); err != nil {
		return err
	}
	expected := digestCanonical(engineMetricSnapshotBody{
		SchemaVersion: snapshot.SchemaVersion, SnapshotKind: snapshot.SnapshotKind,
		WorkUnit: snapshot.WorkUnit, SourceSHA: snapshot.SourceSHA, ObservedAt: snapshot.ObservedAt,
		PlanTier: snapshot.PlanTier, Observations: snapshot.Observations, Unavailable: snapshot.Unavailable,
	})
	if snapshot.Digest == "" || snapshot.Digest != expected {
		return fmt.Errorf("metric snapshot digest mismatch")
	}
	return nil
}

func optimizationSafetyProof(snapshot EngineMetricSnapshot) (EngineOptimizationSafetyProof, bool, string) {
	values, unavailable := snapshotMetricMaps(snapshot)
	proof := EngineOptimizationSafetyProof{
		CurrentSnapshotDigest:       snapshot.Digest,
		DirectPlanMutationAllowed:   false,
		MandatoryMutationAllowed:    false,
		PrerequisiteMutationAllowed: false,
		ProductSpecMutationAllowed:  false,
	}
	required := []string{"MET-SAFETY-GREEN", "MET-COORDINATION-CONFLICTS", "MET-RECOVERY-REVALIDATE", "MET-MANDATORY-CHECKS"}
	for _, id := range required {
		if unavailable[id] {
			return proof, false, "required safety metric is unavailable: " + id
		}
		if _, ok := values[id]; !ok {
			return proof, false, "required safety metric is missing: " + id
		}
	}
	proof.SafetyGreen = values["MET-SAFETY-GREEN"] == 1
	proof.CoordinationConflictFree = values["MET-COORDINATION-CONFLICTS"] == 0
	proof.RecoveryRevalidationClear = values["MET-RECOVERY-REVALIDATE"] == 0
	proof.MandatoryFloorKnown = values["MET-MANDATORY-CHECKS"] >= 0
	switch {
	case !proof.SafetyGreen:
		return proof, false, "mandatory safety floor is not green"
	case !proof.CoordinationConflictFree:
		return proof, false, "coordination conflicts are present"
	case !proof.RecoveryRevalidationClear:
		return proof, false, "recovery requires revalidation"
	case !proof.MandatoryFloorKnown:
		return proof, false, "mandatory check floor is unknown"
	default:
		return proof, true, ""
	}
}

func deriveValidationOptimization(
	body *engineOptimizationReportBody,
	definitions map[string]EngineMetricDefinition,
	current, baseline EngineMetricSnapshot,
	currentValues, baselineValues map[string]float64,
	currentUnavailable, baselineUnavailable map[string]bool,
) {
	id := "OPT-VALIDATION-SELECTIVITY"
	metrics := optimizationOpportunityMetrics[id]
	if unavailableForComparison(metrics, currentUnavailable, baselineUnavailable) {
		addOptimizationWithheld(body, id, metrics, "validation comparison metric is unavailable")
		return
	}
	for _, metricID := range metrics {
		if !metricOptimizationEligibleOrFloor(definitions[metricID]) {
			addOptimizationWithheld(body, id, metrics, "validation metric registry policy is incompatible")
			return
		}
	}
	if current.PlanTier != baseline.PlanTier {
		addOptimizationWithheld(body, id, metrics, "validation tiers differ and are not comparable")
		return
	}
	if currentValues["MET-MANDATORY-CHECKS"] != baselineValues["MET-MANDATORY-CHECKS"] {
		addOptimizationWithheld(body, id, metrics, "mandatory check floor changed")
		return
	}
	costDelta := currentValues["MET-VALIDATION-COST"] - baselineValues["MET-VALIDATION-COST"]
	selectedDelta := currentValues["MET-VALIDATION-SELECTED"] - baselineValues["MET-VALIDATION-SELECTED"]
	if costDelta <= 0 || selectedDelta <= 0 {
		addOptimizationWithheld(body, id, metrics, "no validation-work regression versus baseline")
		return
	}
	body.Recommendations = append(body.Recommendations, EngineOptimizationRecommendation{
		ID: id, Kind: "investigate-validation-selectivity", Advisory: true,
		MetricIDs: append([]string(nil), metrics...),
		EstimatedBenefit: []EngineOptimizationBenefit{
			{MetricID: "MET-VALIDATION-COST", Value: costDelta, Unit: "cost-unit"},
			{MetricID: "MET-VALIDATION-SELECTED", Value: selectedDelta, Unit: "count"},
		},
		Rationale: "current comparable validation execution selects more work at higher relative cost than the verified baseline",
		Constraints: []string{
			"mandatory check set must remain unchanged",
			"prerequisite closure must remain unchanged",
			"strict-risk escalation must remain unchanged",
			"any selection-policy change requires a separate verified lot",
		},
		RevisitTrigger: "new verified baseline or validation topology change",
	})
}

func deriveCacheOptimization(
	body *engineOptimizationReportBody,
	definitions map[string]EngineMetricDefinition,
	currentValues, baselineValues map[string]float64,
	currentUnavailable, baselineUnavailable map[string]bool,
) {
	id := "OPT-CACHE-REUSE"
	metrics := optimizationOpportunityMetrics[id]
	if unavailableForComparison(metrics, currentUnavailable, baselineUnavailable) {
		addOptimizationWithheld(body, id, metrics, "cache reuse evidence is unavailable")
		return
	}
	if !definitions["MET-CACHE-REUSE-RATE"].OptimizationEligible {
		addOptimizationWithheld(body, id, metrics, "cache reuse metric is not optimization eligible")
		return
	}
	delta := baselineValues["MET-CACHE-REUSE-RATE"] - currentValues["MET-CACHE-REUSE-RATE"]
	if delta <= 0 {
		addOptimizationWithheld(body, id, metrics, "cache reuse did not regress versus baseline")
		return
	}
	body.Recommendations = append(body.Recommendations, EngineOptimizationRecommendation{
		ID: id, Kind: "investigate-cache-reuse", Advisory: true, MetricIDs: append([]string(nil), metrics...),
		EstimatedBenefit: []EngineOptimizationBenefit{{MetricID: "MET-CACHE-REUSE-RATE", Value: delta, Unit: "percent"}},
		Rationale:        "verified reusable-evidence rate regressed versus the comparable baseline",
		Constraints: []string{
			"evidence freshness and decision revisit rules must remain unchanged",
			"invalid or stale cache entries must never become reusable",
			"any cache-policy change requires a separate verified lot",
		},
		RevisitTrigger: "new verified performance-cache evidence or cache policy change",
	})
}

func deriveContextOptimization(
	body *engineOptimizationReportBody,
	definitions map[string]EngineMetricDefinition,
	currentValues, baselineValues map[string]float64,
	currentUnavailable, baselineUnavailable map[string]bool,
) {
	id := "OPT-CONTEXT-BOUNDING"
	metrics := optimizationOpportunityMetrics[id]
	if unavailableForComparison(metrics, currentUnavailable, baselineUnavailable) {
		addOptimizationWithheld(body, id, metrics, "context comparison evidence is unavailable")
		return
	}
	for _, metricID := range metrics {
		if !definitions[metricID].OptimizationEligible {
			addOptimizationWithheld(body, id, metrics, "context metric is not optimization eligible")
			return
		}
	}
	sourceDelta := currentValues["MET-CONTEXT-SOURCES"] - baselineValues["MET-CONTEXT-SOURCES"]
	dependencyDelta := currentValues["MET-CONTEXT-DEPENDENCIES"] - baselineValues["MET-CONTEXT-DEPENDENCIES"]
	if sourceDelta <= 0 && dependencyDelta <= 0 {
		addOptimizationWithheld(body, id, metrics, "context bundle did not grow versus baseline")
		return
	}
	var benefits []EngineOptimizationBenefit
	if sourceDelta > 0 {
		benefits = append(benefits, EngineOptimizationBenefit{MetricID: "MET-CONTEXT-SOURCES", Value: sourceDelta, Unit: "count"})
	}
	if dependencyDelta > 0 {
		benefits = append(benefits, EngineOptimizationBenefit{MetricID: "MET-CONTEXT-DEPENDENCIES", Value: dependencyDelta, Unit: "count"})
	}
	body.Recommendations = append(body.Recommendations, EngineOptimizationRecommendation{
		ID: id, Kind: "investigate-context-bounding", Advisory: true, MetricIDs: append([]string(nil), metrics...),
		EstimatedBenefit: benefits,
		Rationale:        "the deterministic agent context grew versus the verified comparable baseline",
		Constraints: []string{
			"required product and dependency sources must remain present",
			"context digest and traceability requirements must remain intact",
			"any context-compilation change requires a separate verified lot",
		},
		RevisitTrigger: "new verified context baseline or dependency graph change",
	})
}

func snapshotMetricMaps(snapshot EngineMetricSnapshot) (map[string]float64, map[string]bool) {
	values := map[string]float64{}
	unavailable := map[string]bool{}
	for _, observation := range snapshot.Observations {
		values[observation.MetricID] = observation.Value
	}
	for _, item := range snapshot.Unavailable {
		unavailable[item.MetricID] = true
	}
	return values, unavailable
}

func metricDefinitionMap(registry EngineMetricsRegistry) map[string]EngineMetricDefinition {
	out := map[string]EngineMetricDefinition{}
	for _, definition := range registry.Metrics {
		out[definition.ID] = definition
	}
	return out
}

func metricOptimizationEligibleOrFloor(definition EngineMetricDefinition) bool {
	return definition.OptimizationEligible || definition.SafetyClass == "safety-floor" || definition.SafetyClass == "quality-floor"
}

func unavailableForComparison(metrics []string, current, baseline map[string]bool) bool {
	for _, id := range metrics {
		if current[id] || baseline[id] {
			return true
		}
	}
	return false
}

func withholdAllOptimizationOpportunities(body *engineOptimizationReportBody, reason string) {
	var ids []string
	for id := range optimizationOpportunityMetrics {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		addOptimizationWithheld(body, id, optimizationOpportunityMetrics[id], reason)
	}
}

func addOptimizationWithheld(body *engineOptimizationReportBody, id string, metrics []string, reason string) {
	body.Withheld = append(body.Withheld, EngineOptimizationWithheld{
		OpportunityID: id, MetricIDs: append([]string(nil), metrics...), Reason: reason,
	})
}

func sortOptimizationReport(body *engineOptimizationReportBody) {
	sort.Slice(body.Recommendations, func(i, j int) bool { return body.Recommendations[i].ID < body.Recommendations[j].ID })
	sort.Slice(body.Withheld, func(i, j int) bool { return body.Withheld[i].OpportunityID < body.Withheld[j].OpportunityID })
}

func finalizeEngineOptimizationReport(body engineOptimizationReportBody) EngineOptimizationReport {
	sortOptimizationReport(&body)
	return EngineOptimizationReport{
		SchemaVersion: body.SchemaVersion, ReportKind: body.ReportKind, WorkUnit: body.WorkUnit,
		SourceSHA: body.SourceSHA, SnapshotDigest: body.SnapshotDigest, BaselineDigest: body.BaselineDigest,
		SafetyStatus: body.SafetyStatus, Recommendations: body.Recommendations, Withheld: body.Withheld,
		Proof: body.Proof, Digest: digestCanonical(body),
	}
}
