package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

type EngineMetricUnavailable struct {
	MetricID string `json:"metric_id"`
	Reason   string `json:"reason"`
}

type EngineMetricSnapshot struct {
	SchemaVersion int                       `json:"schema_version"`
	SnapshotKind  string                    `json:"snapshot_kind"`
	WorkUnit      string                    `json:"work_unit"`
	SourceSHA     string                    `json:"source_sha"`
	ObservedAt    string                    `json:"observed_at"`
	PlanTier      string                    `json:"plan_tier"`
	Observations  []EngineMetricObservation `json:"observations"`
	Unavailable   []EngineMetricUnavailable `json:"unavailable,omitempty"`
	Digest        string                    `json:"digest_sha256"`
}

type engineMetricSnapshotBody struct {
	SchemaVersion int                       `json:"schema_version"`
	SnapshotKind  string                    `json:"snapshot_kind"`
	WorkUnit      string                    `json:"work_unit"`
	SourceSHA     string                    `json:"source_sha"`
	ObservedAt    string                    `json:"observed_at"`
	PlanTier      string                    `json:"plan_tier"`
	Observations  []EngineMetricObservation `json:"observations"`
	Unavailable   []EngineMetricUnavailable `json:"unavailable,omitempty"`
}

func compileEngineMetricSnapshot(root, tempDir string, plan ValidationPlan, execution ValidationExecutionSummary, state CurrentState, graph WorkGraph) (EngineMetricSnapshot, error) {
	var registry EngineMetricsRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineMetricsRegistryPath)), &registry); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if _, err := validateEngineMetricsRegistry(registry); err != nil {
		return EngineMetricSnapshot{}, err
	}
	head, err := localGitHead(root)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	observedAt := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)

	contextOutput := filepath.Join(tempDir, "metrics-context.json")
	contextSummary, err := runContextCompiler(root, plan.WorkUnit, contextOutput, false, state, graph)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	cacheSummary, err := runPerformanceCacheAudit(root, state, graph)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	recovery, err := runRecoveryReconciliation(root, RecoveryReconcileRequest{WorkUnit: plan.WorkUnit, AsOf: observedAt}, state, graph)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	coordination, err := runCoordinationAudit(root, observedAt, graph)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	catalog, err := loadCheckCatalog(root)
	if err != nil {
		return EngineMetricSnapshot{}, err
	}
	mandatory := 0
	for _, check := range catalog.Checks {
		if check.Mandatory {
			mandatory++
		}
	}
	return buildEngineMetricSnapshot(registry, head, observedAt, plan, execution, contextSummary, cacheSummary, recovery, coordination, mandatory)
}

func buildEngineMetricSnapshot(
	registry EngineMetricsRegistry,
	sourceSHA, observedAt string,
	plan ValidationPlan,
	execution ValidationExecutionSummary,
	contextSummary ContextSummary,
	cacheSummary PerformanceCacheAuditSummary,
	recovery RecoveryReconciliation,
	coordination CoordinationAuditSummary,
	mandatoryChecks int,
) (EngineMetricSnapshot, error) {
	body := engineMetricSnapshotBody{
		SchemaVersion: 1,
		SnapshotKind:  "engineering-engine-metrics",
		WorkUnit:      plan.WorkUnit,
		SourceSHA:     sourceSHA,
		ObservedAt:    observedAt,
		PlanTier:      plan.Tier,
	}
	add := func(id string, value float64, source, evidence string) error {
		observation := EngineMetricObservation{
			MetricID: id, Value: value, SourceKey: source, SourceSHA: sourceSHA,
			ObservedAt: observedAt, Evidence: evidence,
		}
		if err := validateEngineMetricObservation(registry, observation); err != nil {
			return err
		}
		body.Observations = append(body.Observations, observation)
		return nil
	}
	unavailable := func(id, reason string) {
		body.Unavailable = append(body.Unavailable, EngineMetricUnavailable{MetricID: id, Reason: reason})
	}

	if err := add("MET-VALIDATION-SELECTED", float64(len(plan.SelectedChecks)), "validation-plan", "validation-plan:tier="+plan.Tier); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if err := add("MET-VALIDATION-EXECUTED", float64(execution.ExecutedChecks), "validation-execution", fmt.Sprintf("validation-execution:preflight=%d", execution.PreflightSatisfied)); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if err := add("MET-VALIDATION-COST", float64(plan.CostUnits), "validation-plan", "validation-plan:cost-units"); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if err := add("MET-MANDATORY-CHECKS", float64(mandatoryChecks), "check-catalog", "check-catalog:mandatory"); err != nil {
		return EngineMetricSnapshot{}, err
	}
	safetyGreen := 1.0
	statusByID := map[string]string{}
	for _, check := range execution.Checks {
		statusByID[check.ID] = check.Status
	}
	for _, planned := range plan.SelectedChecks {
		if !planned.Mandatory {
			continue
		}
		status := statusByID[planned.ID]
		if status != "PASS" && status != "PREFLIGHT_SATISFIED" {
			safetyGreen = 0
			break
		}
	}
	if err := add("MET-SAFETY-GREEN", safetyGreen, "validation-execution", "validation-execution:mandatory-floor"); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if err := add("MET-CONTEXT-SOURCES", float64(contextSummary.Sources), "context-bundle", "context:"+contextSummary.BundleDigest); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if err := add("MET-CONTEXT-DEPENDENCIES", float64(contextSummary.DependencyUnits), "context-bundle", "context:"+contextSummary.BundleDigest); err != nil {
		return EngineMetricSnapshot{}, err
	}
	if cacheSummary.Records == 0 {
		unavailable("MET-CACHE-REUSE-RATE", "performance-cache:"+cacheSummary.Status)
	} else {
		rate := 100 * float64(cacheSummary.Reusable) / float64(cacheSummary.Records)
		if err := add("MET-CACHE-REUSE-RATE", rate, "performance-cache", "performance-cache:"+cacheSummary.Status); err != nil {
			return EngineMetricSnapshot{}, err
		}
	}
	switch recovery.Outcome {
	case "resume", "revalidate":
		value := 0.0
		if recovery.Outcome == "revalidate" {
			value = 1
		}
		if err := add("MET-RECOVERY-REVALIDATE", value, "recovery-reconcile", "recovery:"+recovery.Outcome); err != nil {
			return EngineMetricSnapshot{}, err
		}
	default:
		unavailable("MET-RECOVERY-REVALIDATE", "recovery-outcome:"+recovery.Outcome)
	}
	if err := add("MET-COORDINATION-CONFLICTS", float64(coordination.Conflicts), "coordination", "coordination:"+coordination.Status); err != nil {
		return EngineMetricSnapshot{}, err
	}

	sort.Slice(body.Observations, func(i, j int) bool { return body.Observations[i].MetricID < body.Observations[j].MetricID })
	sort.Slice(body.Unavailable, func(i, j int) bool { return body.Unavailable[i].MetricID < body.Unavailable[j].MetricID })
	if err := validateEngineMetricSnapshotCoverage(registry, body.Observations, body.Unavailable); err != nil {
		return EngineMetricSnapshot{}, err
	}
	return EngineMetricSnapshot{
		SchemaVersion: body.SchemaVersion, SnapshotKind: body.SnapshotKind, WorkUnit: body.WorkUnit,
		SourceSHA: body.SourceSHA, ObservedAt: body.ObservedAt, PlanTier: body.PlanTier,
		Observations: body.Observations, Unavailable: body.Unavailable, Digest: digestCanonical(body),
	}, nil
}

func validateEngineMetricSnapshotCoverage(registry EngineMetricsRegistry, observations []EngineMetricObservation, unavailable []EngineMetricUnavailable) error {
	registered := map[string]bool{}
	for _, metric := range registry.Metrics {
		registered[metric.ID] = true
	}
	seen := map[string]bool{}
	for _, observation := range observations {
		if !registered[observation.MetricID] || seen[observation.MetricID] {
			return fmt.Errorf("snapshot has unknown or duplicate observation %s", observation.MetricID)
		}
		seen[observation.MetricID] = true
	}
	for _, item := range unavailable {
		if !registered[item.MetricID] || seen[item.MetricID] {
			return fmt.Errorf("snapshot has unknown or duplicate unavailable metric %s", item.MetricID)
		}
		if item.Reason == "" {
			return fmt.Errorf("unavailable metric %s requires reason", item.MetricID)
		}
		seen[item.MetricID] = true
	}
	var missing []string
	for id := range registered {
		if !seen[id] {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return fmt.Errorf("snapshot omits registered metrics: %v", missing)
	}
	return nil
}
