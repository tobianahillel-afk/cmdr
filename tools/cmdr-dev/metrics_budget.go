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
	engineEfficiencyPolicyPath   = "engineering/metrics/efficiency-budget-policy.json"
	engineEfficiencyBaselinePath = "engineering/metrics/efficiency-baselines.json"
)

type EngineEfficiencyFloorRule struct {
	MetricID   string  `json:"metric_id"`
	Comparator string  `json:"comparator"`
	Value      float64 `json:"value"`
	Rationale  string  `json:"rationale"`
}

type EngineEfficiencyBudgetRule struct {
	MetricID                     string  `json:"metric_id"`
	Blocking                     bool    `json:"blocking"`
	MaxRelativeRegressionPercent float64 `json:"max_relative_regression_percent"`
	Rationale                    string  `json:"rationale"`
}

type EngineEfficiencyFloorProvenance struct {
	SourceSHA      string `json:"source_sha"`
	WorkflowRunID  int64  `json:"workflow_run_id"`
	MandatoryFloor int    `json:"mandatory_floor"`
}

type EngineEfficiencyBudgetPolicy struct {
	SchemaVersion            int                             `json:"schema_version"`
	PolicyKind               string                          `json:"policy_kind"`
	DefaultPolicy            string                          `json:"default_policy"`
	ProfileMatch             []string                        `json:"profile_match"`
	MandatoryFloorProvenance EngineEfficiencyFloorProvenance `json:"mandatory_floor_provenance"`
	SafetyFloors             []EngineEfficiencyFloorRule     `json:"safety_floors"`
	EfficiencyBudgets        []EngineEfficiencyBudgetRule    `json:"efficiency_budgets"`
}

type EngineEfficiencyBaselineMetric struct {
	MetricID string  `json:"metric_id"`
	Value    float64 `json:"value"`
}

type EngineEfficiencyBaseline struct {
	ID             string                           `json:"id"`
	PlanTier       string                           `json:"plan_tier"`
	RiskDomains    []string                         `json:"risk_domains"`
	MandatoryFloor int                              `json:"mandatory_floor"`
	SourceSHA      string                           `json:"source_sha"`
	ObservedAt     string                           `json:"observed_at"`
	SnapshotDigest string                           `json:"snapshot_digest"`
	WorkflowRunID  int64                            `json:"workflow_run_id"`
	Metrics        []EngineEfficiencyBaselineMetric `json:"metrics"`
}

type EngineEfficiencyBaselineRegistry struct {
	SchemaVersion int                        `json:"schema_version"`
	BaselineKind  string                     `json:"baseline_kind"`
	Baselines     []EngineEfficiencyBaseline `json:"baselines"`
}

type EngineEfficiencyBudgetFinding struct {
	MetricID                  string   `json:"metric_id"`
	BaselineValue             *float64 `json:"baseline_value,omitempty"`
	CurrentValue              float64  `json:"current_value"`
	AllowedValue              float64  `json:"allowed_value"`
	RelativeRegressionPercent *float64 `json:"relative_regression_percent,omitempty"`
	Blocking                  bool     `json:"blocking"`
	Reason                    string   `json:"reason"`
}

type EngineEfficiencyBudgetWithheld struct {
	MetricID string `json:"metric_id"`
	Reason   string `json:"reason"`
}

type EngineEfficiencyBudgetReport struct {
	SchemaVersion      int                              `json:"schema_version"`
	ReportKind         string                           `json:"report_kind"`
	WorkUnit           string                           `json:"work_unit"`
	SourceSHA          string                           `json:"source_sha"`
	SnapshotDigest     string                           `json:"snapshot_digest"`
	ProfileKey         string                           `json:"profile_key"`
	BaselineID         string                           `json:"baseline_id,omitempty"`
	BaselineStatus     string                           `json:"baseline_status"`
	SafetyStatus       string                           `json:"safety_status"`
	EfficiencyStatus   string                           `json:"efficiency_status"`
	BlockingViolations []EngineEfficiencyBudgetFinding  `json:"blocking_violations,omitempty"`
	Advisories         []EngineEfficiencyBudgetFinding  `json:"advisories,omitempty"`
	Withheld           []EngineEfficiencyBudgetWithheld `json:"withheld,omitempty"`
	Digest             string                           `json:"digest_sha256"`
}

type engineEfficiencyBudgetReportBody struct {
	SchemaVersion      int                              `json:"schema_version"`
	ReportKind         string                           `json:"report_kind"`
	WorkUnit           string                           `json:"work_unit"`
	SourceSHA          string                           `json:"source_sha"`
	SnapshotDigest     string                           `json:"snapshot_digest"`
	ProfileKey         string                           `json:"profile_key"`
	BaselineID         string                           `json:"baseline_id,omitempty"`
	BaselineStatus     string                           `json:"baseline_status"`
	SafetyStatus       string                           `json:"safety_status"`
	EfficiencyStatus   string                           `json:"efficiency_status"`
	BlockingViolations []EngineEfficiencyBudgetFinding  `json:"blocking_violations,omitempty"`
	Advisories         []EngineEfficiencyBudgetFinding  `json:"advisories,omitempty"`
	Withheld           []EngineEfficiencyBudgetWithheld `json:"withheld,omitempty"`
}

var engineMetricDigestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

func compileEngineEfficiencyBudgetReport(root string, snapshot EngineMetricSnapshot, plan ValidationPlan) (EngineEfficiencyBudgetReport, error) {
	var registry EngineMetricsRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineMetricsRegistryPath)), &registry); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	if _, err := validateEngineMetricsRegistry(registry); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	if err := validateOptimizationSnapshot(registry, snapshot); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	var policy EngineEfficiencyBudgetPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineEfficiencyPolicyPath)), &policy); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	if err := validateEngineEfficiencyBudgetPolicy(registry, policy); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	var baselines EngineEfficiencyBaselineRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(engineEfficiencyBaselinePath)), &baselines); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	if err := validateEngineEfficiencyBaselines(registry, baselines); err != nil {
		return EngineEfficiencyBudgetReport{}, err
	}
	return evaluateEngineEfficiencyBudget(registry, policy, baselines, snapshot, plan)
}

func validateEngineEfficiencyBudgetPolicy(registry EngineMetricsRegistry, policy EngineEfficiencyBudgetPolicy) error {
	if policy.SchemaVersion != 1 || policy.PolicyKind != "engineering-engine-efficiency-budgets" ||
		policy.DefaultPolicy != "deny-unknown-metric" {
		return fmt.Errorf("invalid engine efficiency budget policy header")
	}
	expectedProfile := []string{"plan_tier", "risk_domains", "mandatory_floor"}
	if len(policy.ProfileMatch) != len(expectedProfile) {
		return fmt.Errorf("efficiency budget profile_match must be %v", expectedProfile)
	}
	for i := range expectedProfile {
		if policy.ProfileMatch[i] != expectedProfile[i] {
			return fmt.Errorf("efficiency budget profile_match must be %v", expectedProfile)
		}
	}
	if sha, err := validateFullCommitID(policy.MandatoryFloorProvenance.SourceSHA); err != nil || sha != policy.MandatoryFloorProvenance.SourceSHA {
		return fmt.Errorf("mandatory floor provenance requires canonical source SHA")
	}
	if policy.MandatoryFloorProvenance.WorkflowRunID <= 0 || policy.MandatoryFloorProvenance.MandatoryFloor <= 0 {
		return fmt.Errorf("mandatory floor provenance is incomplete")
	}
	definitions := metricDefinitionMap(registry)
	seen := map[string]bool{}
	for _, floor := range policy.SafetyFloors {
		definition, ok := definitions[floor.MetricID]
		if !ok || seen[floor.MetricID] {
			return fmt.Errorf("invalid or duplicate safety floor metric %s", floor.MetricID)
		}
		seen[floor.MetricID] = true
		if definition.SafetyClass == "efficiency" {
			return fmt.Errorf("efficiency metric %s cannot define a safety floor", floor.MetricID)
		}
		if floor.Comparator != "equal" && floor.Comparator != "minimum" {
			return fmt.Errorf("safety floor %s has unsupported comparator %s", floor.MetricID, floor.Comparator)
		}
		if math.IsNaN(floor.Value) || math.IsInf(floor.Value, 0) || strings.TrimSpace(floor.Rationale) == "" {
			return fmt.Errorf("safety floor %s is invalid", floor.MetricID)
		}
	}
	for _, budget := range policy.EfficiencyBudgets {
		definition, ok := definitions[budget.MetricID]
		if !ok || seen[budget.MetricID] {
			return fmt.Errorf("invalid or duplicate efficiency budget metric %s", budget.MetricID)
		}
		seen[budget.MetricID] = true
		if definition.SafetyClass != "efficiency" || !definition.OptimizationEligible {
			return fmt.Errorf("metric %s is not eligible for efficiency budgeting", budget.MetricID)
		}
		if definition.Direction != "lower-better" && definition.Direction != "higher-better" {
			return fmt.Errorf("efficiency metric %s has unsupported direction %s", budget.MetricID, definition.Direction)
		}
		if math.IsNaN(budget.MaxRelativeRegressionPercent) || math.IsInf(budget.MaxRelativeRegressionPercent, 0) ||
			budget.MaxRelativeRegressionPercent < 0 || budget.MaxRelativeRegressionPercent > 100 ||
			strings.TrimSpace(budget.Rationale) == "" {
			return fmt.Errorf("efficiency budget %s is invalid", budget.MetricID)
		}
	}
	return nil
}

func validateEngineEfficiencyBaselines(registry EngineMetricsRegistry, baselines EngineEfficiencyBaselineRegistry) error {
	if baselines.SchemaVersion != 1 || baselines.BaselineKind != "engineering-engine-efficiency-baselines" {
		return fmt.Errorf("invalid engine efficiency baseline registry header")
	}
	definitions := metricDefinitionMap(registry)
	ids := map[string]bool{}
	profiles := map[string]bool{}
	for _, baseline := range baselines.Baselines {
		if strings.TrimSpace(baseline.ID) == "" || ids[baseline.ID] {
			return fmt.Errorf("invalid or duplicate efficiency baseline id %q", baseline.ID)
		}
		ids[baseline.ID] = true
		if baseline.PlanTier != "fast" && baseline.PlanTier != "standard" && baseline.PlanTier != "strict" {
			return fmt.Errorf("baseline %s has invalid plan tier", baseline.ID)
		}
		if baseline.MandatoryFloor <= 0 || baseline.WorkflowRunID <= 0 {
			return fmt.Errorf("baseline %s has incomplete provenance", baseline.ID)
		}
		if sha, err := validateFullCommitID(baseline.SourceSHA); err != nil || sha != baseline.SourceSHA {
			return fmt.Errorf("baseline %s has invalid source SHA", baseline.ID)
		}
		if _, err := parseCanonicalLeaseTime(baseline.ObservedAt); err != nil {
			return fmt.Errorf("baseline %s observed_at: %w", baseline.ID, err)
		}
		if !engineMetricDigestPattern.MatchString(baseline.SnapshotDigest) {
			return fmt.Errorf("baseline %s has invalid snapshot digest", baseline.ID)
		}
		risks := append([]string(nil), baseline.RiskDomains...)
		sort.Strings(risks)
		for i := range risks {
			if !knownRiskDomains[risks[i]] || (i > 0 && risks[i] == risks[i-1]) {
				return fmt.Errorf("baseline %s has invalid risk domains", baseline.ID)
			}
		}
		if strings.Join(risks, "\x00") != strings.Join(baseline.RiskDomains, "\x00") {
			return fmt.Errorf("baseline %s risk domains must be sorted", baseline.ID)
		}
		metricSeen := map[string]bool{}
		for _, metric := range baseline.Metrics {
			definition, ok := definitions[metric.MetricID]
			if !ok || metricSeen[metric.MetricID] {
				return fmt.Errorf("baseline %s has unknown or duplicate metric %s", baseline.ID, metric.MetricID)
			}
			metricSeen[metric.MetricID] = true
			if math.IsNaN(metric.Value) || math.IsInf(metric.Value, 0) || metric.Value < 0 {
				return fmt.Errorf("baseline %s metric %s has invalid value", baseline.ID, metric.MetricID)
			}
			if definition.Unit == "percent" && metric.Value > 100 {
				return fmt.Errorf("baseline %s percentage metric %s exceeds 100", baseline.ID, metric.MetricID)
			}
		}
		key := efficiencyProfileKey(baseline.PlanTier, baseline.RiskDomains, baseline.MandatoryFloor)
		if profiles[key] {
			return fmt.Errorf("duplicate efficiency baseline profile %s", key)
		}
		profiles[key] = true
	}
	return nil
}

func evaluateEngineEfficiencyBudget(
	registry EngineMetricsRegistry,
	policy EngineEfficiencyBudgetPolicy,
	baselines EngineEfficiencyBaselineRegistry,
	snapshot EngineMetricSnapshot,
	plan ValidationPlan,
) (EngineEfficiencyBudgetReport, error) {
	values, unavailable := snapshotMetricMaps(snapshot)
	mandatoryValue, ok := values["MET-MANDATORY-CHECKS"]
	if !ok || unavailable["MET-MANDATORY-CHECKS"] || mandatoryValue != math.Trunc(mandatoryValue) {
		return EngineEfficiencyBudgetReport{}, fmt.Errorf("mandatory check floor is unavailable or non-integral")
	}
	risks := append([]string(nil), plan.RiskDomains...)
	sort.Strings(risks)
	profileKey := efficiencyProfileKey(plan.Tier, risks, int(mandatoryValue))
	body := engineEfficiencyBudgetReportBody{
		SchemaVersion:    1,
		ReportKind:       "engineering-engine-efficiency-budget",
		WorkUnit:         snapshot.WorkUnit,
		SourceSHA:        snapshot.SourceSHA,
		SnapshotDigest:   snapshot.Digest,
		ProfileKey:       profileKey,
		BaselineStatus:   "unavailable",
		SafetyStatus:     "pass",
		EfficiencyStatus: "not-enforced-no-baseline",
	}
	for _, floor := range policy.SafetyFloors {
		current, exists := values[floor.MetricID]
		reason := ""
		violated := false
		if unavailable[floor.MetricID] || !exists {
			violated = true
			reason = "required safety/quality floor metric is unavailable"
		} else {
			switch floor.Comparator {
			case "equal":
				violated = current != floor.Value
				reason = "current value must equal verified floor"
			case "minimum":
				violated = current < floor.Value
				reason = "current value must not fall below verified floor"
			}
		}
		if violated {
			value := current
			body.BlockingViolations = append(body.BlockingViolations, EngineEfficiencyBudgetFinding{
				MetricID: floor.MetricID, CurrentValue: value, AllowedValue: floor.Value,
				Blocking: true, Reason: reason,
			})
		}
	}
	if len(body.BlockingViolations) > 0 {
		body.SafetyStatus = "fail"
	}
	var baseline *EngineEfficiencyBaseline
	for i := range baselines.Baselines {
		candidate := &baselines.Baselines[i]
		if efficiencyProfileKey(candidate.PlanTier, candidate.RiskDomains, candidate.MandatoryFloor) == profileKey {
			baseline = candidate
			break
		}
	}
	if baseline == nil {
		for _, budget := range policy.EfficiencyBudgets {
			body.Withheld = append(body.Withheld, EngineEfficiencyBudgetWithheld{
				MetricID: budget.MetricID, Reason: "verified comparable efficiency baseline is unavailable",
			})
		}
	} else {
		body.BaselineID = baseline.ID
		body.BaselineStatus = "matched"
		body.EfficiencyStatus = "pass"
		baselineValues := map[string]float64{}
		for _, metric := range baseline.Metrics {
			baselineValues[metric.MetricID] = metric.Value
		}
		definitions := metricDefinitionMap(registry)
		for _, budget := range policy.EfficiencyBudgets {
			current, currentOK := values[budget.MetricID]
			base, baselineOK := baselineValues[budget.MetricID]
			if unavailable[budget.MetricID] || !currentOK {
				body.Withheld = append(body.Withheld, EngineEfficiencyBudgetWithheld{
					MetricID: budget.MetricID, Reason: "current efficiency metric is unavailable",
				})
				continue
			}
			if !baselineOK {
				body.Withheld = append(body.Withheld, EngineEfficiencyBudgetWithheld{
					MetricID: budget.MetricID, Reason: "baseline does not contain this efficiency metric",
				})
				continue
			}
			definition := definitions[budget.MetricID]
			allowed := base
			if definition.Direction == "lower-better" {
				allowed = base * (1 + budget.MaxRelativeRegressionPercent/100)
			} else {
				allowed = base * (1 - budget.MaxRelativeRegressionPercent/100)
			}
			regressed := (definition.Direction == "lower-better" && current > allowed) ||
				(definition.Direction == "higher-better" && current < allowed)
			if !regressed {
				continue
			}
			baseCopy := base
			var percent *float64
			if base != 0 {
				p := 100 * math.Abs(current-base) / math.Abs(base)
				percent = &p
			}
			finding := EngineEfficiencyBudgetFinding{
				MetricID: budget.MetricID, BaselineValue: &baseCopy, CurrentValue: current,
				AllowedValue: allowed, RelativeRegressionPercent: percent,
				Blocking: budget.Blocking, Reason: "same-profile efficiency metric regressed beyond explicit budget",
			}
			if budget.Blocking {
				body.BlockingViolations = append(body.BlockingViolations, finding)
			} else {
				body.Advisories = append(body.Advisories, finding)
			}
		}
		if len(body.BlockingViolations) > 0 {
			body.EfficiencyStatus = "fail"
		} else if len(body.Advisories) > 0 {
			body.EfficiencyStatus = "pass-with-advisories"
		}
	}
	sortBudgetReport(&body)
	report := finalizeEngineEfficiencyBudgetReport(body)
	if len(report.BlockingViolations) > 0 {
		return report, fmt.Errorf("engine efficiency budget blocked: %d blocking violation(s)", len(report.BlockingViolations))
	}
	return report, nil
}

func efficiencyProfileKey(tier string, risks []string, mandatoryFloor int) string {
	normalized := append([]string(nil), risks...)
	sort.Strings(normalized)
	return digestCanonical(struct {
		Tier           string   `json:"tier"`
		RiskDomains    []string `json:"risk_domains"`
		MandatoryFloor int      `json:"mandatory_floor"`
	}{Tier: tier, RiskDomains: normalized, MandatoryFloor: mandatoryFloor})
}

func sortBudgetReport(body *engineEfficiencyBudgetReportBody) {
	sort.Slice(body.BlockingViolations, func(i, j int) bool {
		return body.BlockingViolations[i].MetricID < body.BlockingViolations[j].MetricID
	})
	sort.Slice(body.Advisories, func(i, j int) bool {
		return body.Advisories[i].MetricID < body.Advisories[j].MetricID
	})
	sort.Slice(body.Withheld, func(i, j int) bool {
		return body.Withheld[i].MetricID < body.Withheld[j].MetricID
	})
}

func finalizeEngineEfficiencyBudgetReport(body engineEfficiencyBudgetReportBody) EngineEfficiencyBudgetReport {
	sortBudgetReport(&body)
	return EngineEfficiencyBudgetReport{
		SchemaVersion: body.SchemaVersion, ReportKind: body.ReportKind, WorkUnit: body.WorkUnit,
		SourceSHA: body.SourceSHA, SnapshotDigest: body.SnapshotDigest, ProfileKey: body.ProfileKey,
		BaselineID: body.BaselineID, BaselineStatus: body.BaselineStatus,
		SafetyStatus: body.SafetyStatus, EfficiencyStatus: body.EfficiencyStatus,
		BlockingViolations: body.BlockingViolations, Advisories: body.Advisories, Withheld: body.Withheld,
		Digest: digestCanonical(body),
	}
}
