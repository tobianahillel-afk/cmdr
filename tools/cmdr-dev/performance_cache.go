package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const performanceCachePath = "engineering/performance/performance-cache.json"

type PerformanceCacheRecord struct {
	ID                  string                     `json:"id"`
	TargetID            string                     `json:"target_id"`
	EnvironmentID       string                     `json:"environment_id"`
	SourceSHA           string                     `json:"source_sha"`
	Toolchain           string                     `json:"toolchain"`
	TargetDigest        string                     `json:"target_digest"`
	BudgetDigest        string                     `json:"budget_digest"`
	EnvironmentDigest   string                     `json:"environment_digest"`
	WorkloadDigest      string                     `json:"workload_digest"`
	DecisionBasisDigest string                     `json:"decision_basis_digest"`
	CacheKey            string                     `json:"cache_key"`
	EvidenceDigest      string                     `json:"evidence_digest"`
	Outcome             string                     `json:"outcome"`
	DecisionIDs         []string                   `json:"decision_ids"`
	Result              PerformanceBenchmarkResult `json:"result"`
}

type PerformanceCacheRegistry struct {
	SchemaVersion int                      `json:"schema_version"`
	RegistryKind  string                   `json:"registry_kind"`
	Records       []PerformanceCacheRecord `json:"records"`
}

type PerformanceRevisitSignal struct {
	DecisionID  string `json:"decision_id"`
	Kind        string `json:"kind"`
	State       string `json:"state"`
	EvidenceRef string `json:"evidence_ref"`
	Detail      string `json:"detail"`
}

type PerformanceCacheEntryStatus struct {
	ID       string   `json:"id"`
	Reusable bool     `json:"reusable"`
	Reasons  []string `json:"reasons,omitempty"`
}

type PerformanceCacheAuditSummary struct {
	Records     int                           `json:"records"`
	Reusable    int                           `json:"reusable"`
	Invalid     int                           `json:"invalid"`
	Regressions int                           `json:"regressions"`
	Status      string                        `json:"status"`
	Entries     []PerformanceCacheEntryStatus `json:"entries,omitempty"`
	Signals     []PerformanceRevisitSignal    `json:"revisit_signals,omitempty"`
}

type performanceCacheIdentity struct {
	SourceSHA           string `json:"source_sha"`
	TargetID            string `json:"target_id"`
	TargetDigest        string `json:"target_digest"`
	BudgetDigest        string `json:"budget_digest"`
	EnvironmentID       string `json:"environment_id"`
	EnvironmentDigest   string `json:"environment_digest"`
	WorkloadDigest      string `json:"workload_digest"`
	Toolchain           string `json:"toolchain"`
	DecisionBasisDigest string `json:"decision_basis_digest"`
}

func runPerformanceCacheAudit(root string, state CurrentState, graph WorkGraph) (PerformanceCacheAuditSummary, error) {
	summary, _, err := evaluatePerformanceCache(root, state, graph)
	return summary, err
}

func reusablePerformanceResults(root string) (map[string]PerformanceBenchmarkResult, PerformanceCacheAuditSummary, error) {
	state, graph, err := loadRepositoryState(root)
	if err != nil {
		return nil, PerformanceCacheAuditSummary{}, err
	}
	summary, results, err := evaluatePerformanceCache(root, state, graph)
	return results, summary, err
}

func evaluatePerformanceCache(root string, state CurrentState, graph WorkGraph) (PerformanceCacheAuditSummary, map[string]PerformanceBenchmarkResult, error) {
	_, registry, _, err := loadValidatedPerformanceConfiguration(root)
	if err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	benchmarkPolicy, err := loadBenchmarkPolicy(root)
	if err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	if _, err := loadBenchmarkBaselines(root, registry, benchmarkPolicy); err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}

	var cache PerformanceCacheRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(performanceCachePath)), &cache); err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	if cache.SchemaVersion != 1 || cache.RegistryKind != "performance-evidence-cache" {
		return PerformanceCacheAuditSummary{}, nil, fmt.Errorf("invalid performance cache registry header")
	}

	decisionCache, err := runReusableEvidenceCache(root, "", state, graph)
	if err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	decisionBasisDigest := digestCanonical(decisionCache.Entries)
	reusableDecision := map[string]bool{}
	for _, entry := range decisionCache.Entries {
		reusableDecision[entry.DecisionID] = true
	}

	var decisions DecisionRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionRegistryPath)), &decisions); err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	decisionByID := map[string]EngineeringDecision{}
	for _, decision := range decisions.Decisions {
		decisionByID[decision.ID] = decision
	}

	sourceSHA, err := currentSourceCommit(root)
	if err != nil {
		return PerformanceCacheAuditSummary{}, nil, err
	}
	targets := map[string]PerformanceTarget{}
	for _, target := range registry.Targets {
		targets[target.ID] = target
	}
	environments := map[string]PerformanceEnvironment{}
	for _, environment := range registry.Environments {
		environments[environment.ID] = environment
	}
	workloads := benchmarkWorkloadIndex(benchmarkPolicy)

	summary := PerformanceCacheAuditSummary{Records: len(cache.Records)}
	reusable := map[string]PerformanceBenchmarkResult{}
	seenIDs := map[string]bool{}
	seenKeys := map[string]bool{}

	for _, record := range cache.Records {
		if strings.TrimSpace(record.ID) == "" || seenIDs[record.ID] {
			return summary, nil, fmt.Errorf("invalid or duplicate performance cache id %q", record.ID)
		}
		seenIDs[record.ID] = true
		if record.Outcome != "pass" && record.Outcome != "regression" {
			return summary, nil, fmt.Errorf("performance cache %s has invalid outcome %q", record.ID, record.Outcome)
		}
		if !sha256Pattern.MatchString(record.CacheKey) || !sha256Pattern.MatchString(record.EvidenceDigest) {
			return summary, nil, fmt.Errorf("performance cache %s requires SHA-256 cache/evidence digests", record.ID)
		}

		status, signalItems, result, key := evaluatePerformanceCacheRecord(
			record, sourceSHA, decisionBasisDigest, targets, environments, workloads, decisionByID, reusableDecision,
		)
		summary.Entries = append(summary.Entries, status)
		summary.Signals = append(summary.Signals, signalItems...)
		if status.Reusable {
			if seenKeys[key] {
				return summary, nil, fmt.Errorf("duplicate reusable performance cache key %s", key)
			}
			seenKeys[key] = true
			reusable[benchmarkBaselineKey(record.TargetID, record.EnvironmentID)] = result
			summary.Reusable++
		} else {
			summary.Invalid++
		}
		if record.Outcome == "regression" {
			summary.Regressions++
		}
	}
	sort.Slice(summary.Entries, func(i, j int) bool { return summary.Entries[i].ID < summary.Entries[j].ID })
	sort.Slice(summary.Signals, func(i, j int) bool {
		if summary.Signals[i].DecisionID != summary.Signals[j].DecisionID {
			return summary.Signals[i].DecisionID < summary.Signals[j].DecisionID
		}
		return summary.Signals[i].EvidenceRef < summary.Signals[j].EvidenceRef
	})

	switch {
	case len(registry.Targets) == 0 && len(cache.Records) == 0:
		summary.Status = "not-applicable"
	case summary.Regressions > 0:
		summary.Status = "revisit-required"
	case summary.Invalid > 0:
		summary.Status = "cache-miss"
	default:
		summary.Status = "reusable"
	}
	if summary.Regressions > 0 {
		return summary, reusable, fmt.Errorf("performance cache contains %d benchmark regression(s); E5 revisit required", summary.Regressions)
	}
	return summary, reusable, nil
}

func evaluatePerformanceCacheRecord(
	record PerformanceCacheRecord,
	currentSourceSHA, decisionBasisDigest string,
	targets map[string]PerformanceTarget,
	environments map[string]PerformanceEnvironment,
	workloads map[string]BenchmarkWorkloadDefinition,
	decisions map[string]EngineeringDecision,
	reusableDecision map[string]bool,
) (PerformanceCacheEntryStatus, []PerformanceRevisitSignal, PerformanceBenchmarkResult, string) {
	status := PerformanceCacheEntryStatus{ID: record.ID}
	addReason := func(reason string) {
		if !containsString(status.Reasons, reason) {
			status.Reasons = append(status.Reasons, reason)
		}
	}

	target, targetOK := targets[record.TargetID]
	if !targetOK {
		addReason("target-missing")
	}
	environment, environmentOK := environments[record.EnvironmentID]
	if !environmentOK {
		addReason("environment-missing")
	}
	workload, workloadOK := BenchmarkWorkloadDefinition{}, false
	if targetOK {
		workload, workloadOK = workloads[target.WorkloadKey]
		if !workloadOK {
			addReason("workload-missing")
		}
	}

	if record.SourceSHA != currentSourceSHA {
		addReason("source-change")
	}
	if record.Toolchain != "go"+engineeringGoVersion {
		addReason("toolchain-change")
	}
	if record.DecisionBasisDigest != decisionBasisDigest {
		addReason("research-decision-freshness-change")
	}

	expectedKey := ""
	if targetOK && environmentOK && workloadOK {
		identity := performanceCacheIdentity{
			SourceSHA: currentSourceSHA, TargetID: target.ID, TargetDigest: digestCanonical(target),
			BudgetDigest: digestCanonical(target.Metrics), EnvironmentID: environment.ID,
			EnvironmentDigest: digestCanonical(environment), WorkloadDigest: digestCanonical(workload),
			Toolchain: "go" + engineeringGoVersion, DecisionBasisDigest: decisionBasisDigest,
		}
		expectedKey = digestCanonical(identity)
		if record.TargetDigest != identity.TargetDigest {
			addReason("target-change")
		}
		if record.BudgetDigest != identity.BudgetDigest {
			addReason("budget-change")
		}
		if record.EnvironmentDigest != identity.EnvironmentDigest {
			addReason("environment-change")
		}
		if record.WorkloadDigest != identity.WorkloadDigest {
			addReason("workload-change")
		}
		if record.CacheKey != expectedKey {
			addReason("cache-key-change")
		}
	}

	if digestCanonical(record.Result) != record.EvidenceDigest {
		addReason("evidence-digest-change")
	}
	if record.Result.TargetID != record.TargetID || record.Result.EnvironmentID != record.EnvironmentID ||
		record.Result.SourceSHA != record.SourceSHA || record.Result.Toolchain != record.Toolchain ||
		record.Result.TargetDigest != record.TargetDigest || record.Result.EnvironmentDigest != record.EnvironmentDigest ||
		record.Result.WorkloadDigest != record.WorkloadDigest {
		addReason("result-identity-change")
	}

	regressed := false
	for _, metric := range record.Result.Metrics {
		if !metric.AbsolutePass || !metric.RelativePass {
			regressed = true
			break
		}
	}
	if record.Outcome == "pass" && regressed {
		addReason("pass-record-contains-regression")
	}
	if record.Outcome == "regression" && !regressed {
		addReason("regression-record-has-no-failing-metric")
	}

	var signals []PerformanceRevisitSignal
	seenDecision := map[string]bool{}
	for _, decisionID := range record.DecisionIDs {
		if seenDecision[decisionID] {
			addReason("duplicate-decision-binding")
			continue
		}
		seenDecision[decisionID] = true
		decision, ok := decisions[decisionID]
		if !ok || decision.Status != "accepted" || !decision.PerformanceSensitive {
			addReason("decision-binding-invalid")
			continue
		}
		if !reusableDecision[decisionID] {
			addReason("decision-evidence-not-reusable")
			continue
		}
		if record.Outcome == "regression" {
			signals = append(signals, PerformanceRevisitSignal{
				DecisionID: decisionID, Kind: "benchmark-regression", State: "fired",
				EvidenceRef: record.ID, Detail: "cached benchmark result violates an accepted performance constraint",
			})
		}
	}
	if record.Outcome == "regression" && len(record.DecisionIDs) == 0 {
		addReason("regression-without-decision-binding")
	}

	sort.Strings(status.Reasons)
	status.Reusable = record.Outcome == "pass" && len(status.Reasons) == 0
	return status, signals, record.Result, expectedKey
}
