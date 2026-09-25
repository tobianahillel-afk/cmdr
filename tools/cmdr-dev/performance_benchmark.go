package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	benchmarkPolicyPath    = "engineering/performance/benchmark-policy.json"
	benchmarkBaselinePath  = "engineering/performance/benchmark-baselines.json"
	defaultCIEnvironment   = "PERF-ENV-GITHUB-UBUNTU-SHARED"
	engineeringNodeVersion = "24.21.0"
)

type BenchmarkWorkloadDefinition struct {
	Key                 string `json:"key"`
	HandlerKey          string `json:"handler_key"`
	WarmupIterations    int    `json:"warmup_iterations"`
	SampleCount         int    `json:"sample_count"`
	MaxSampleDurationMS int    `json:"max_sample_duration_ms"`
}

type BenchmarkPolicy struct {
	SchemaVersion int                           `json:"schema_version"`
	PolicyKind    string                        `json:"policy_kind"`
	Workloads     []BenchmarkWorkloadDefinition `json:"workloads"`
}

type PerformanceBaseline struct {
	ID                string             `json:"id"`
	TargetID          string             `json:"target_id"`
	EnvironmentID     string             `json:"environment_id"`
	SourceSHA         string             `json:"source_sha"`
	Toolchain         string             `json:"toolchain"`
	TargetDigest      string             `json:"target_digest"`
	EnvironmentDigest string             `json:"environment_digest"`
	WorkloadDigest    string             `json:"workload_digest"`
	MetricValues      map[string]float64 `json:"metric_values"`
}

type PerformanceBaselineRegistry struct {
	SchemaVersion int                   `json:"schema_version"`
	RegistryKind  string                `json:"registry_kind"`
	Records       []PerformanceBaseline `json:"records"`
}

type BenchmarkMeasurement struct {
	Kind  string  `json:"kind"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

type BenchmarkSample struct {
	Measurements []BenchmarkMeasurement `json:"measurements"`
}

type BenchmarkMetricResult struct {
	ID                           string  `json:"id"`
	Kind                         string  `json:"kind"`
	Unit                         string  `json:"unit"`
	Aggregation                  string  `json:"aggregation"`
	Observed                     float64 `json:"observed"`
	Comparator                   string  `json:"comparator"`
	Budget                       float64 `json:"budget"`
	AbsolutePass                 bool    `json:"absolute_pass"`
	Baseline                     float64 `json:"baseline,omitempty"`
	RelativeRegressionPercent    float64 `json:"relative_regression_percent,omitempty"`
	MaxRelativeRegressionPercent float64 `json:"max_relative_regression_percent,omitempty"`
	RelativePass                 bool    `json:"relative_pass"`
}

type PerformanceBenchmarkResult struct {
	TargetID          string                  `json:"target_id"`
	EnvironmentID     string                  `json:"environment_id"`
	SourceSHA         string                  `json:"source_sha"`
	Toolchain         string                  `json:"toolchain"`
	TargetDigest      string                  `json:"target_digest"`
	EnvironmentDigest string                  `json:"environment_digest"`
	WorkloadDigest    string                  `json:"workload_digest"`
	Samples           int                     `json:"samples"`
	Metrics           []BenchmarkMetricResult `json:"metrics"`
}

type pilotProjectionProbeObservation struct {
	Operation            string  `json:"operation,omitempty"`
	Mode                 string  `json:"mode"`
	Operations           int     `json:"operations"`
	ElapsedNS            int64   `json:"elapsed_ns"`
	NSPerOperation       float64 `json:"ns_per_operation"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
}

type pilotProjectionProbe struct {
	moduleRoot string
	binaryPath string
}

type PerformanceBenchmarkAuditSummary struct {
	Stage         string                       `json:"stage"`
	EnvironmentID string                       `json:"environment_id"`
	SourceSHA     string                       `json:"source_sha"`
	Registered    int                          `json:"registered_targets"`
	Selected      int                          `json:"selected_targets"`
	Executed      int                          `json:"executed_targets"`
	Reused        int                          `json:"reused_targets"`
	Metrics       int                          `json:"metrics_evaluated"`
	Status        string                       `json:"status"`
	ByScope       map[string]int               `json:"by_scope"`
	Results       []PerformanceBenchmarkResult `json:"results,omitempty"`
}

var knownBenchmarkHandlerKeys = map[string]bool{
	"builtin-cmdr-dev-metadata-audit-v1":          true,
	"builtin-pilot-context-projection-v1":         true,
	"builtin-event-search-validation-v1":          true,
	"builtin-event-search-orchestration-v1":       true,
	"builtin-event-search-frontend-state-v1":      true,
	"builtin-event-inspection-preimplementation-v1": true,
}

func runPerformanceBenchmarkAudit(root, stage, environmentID, changesFile string) (PerformanceBenchmarkAuditSummary, error) {
	policy, registry, architecture, err := loadValidatedPerformanceConfiguration(root)
	if err != nil {
		return PerformanceBenchmarkAuditSummary{}, err
	}
	if _, err := validatePerformanceRegistry(policy, registry, architecture); err != nil {
		return PerformanceBenchmarkAuditSummary{}, err
	}
	benchmarkPolicy, err := loadBenchmarkPolicy(root)
	if err != nil {
		return PerformanceBenchmarkAuditSummary{}, err
	}
	baselines, err := loadBenchmarkBaselines(root, registry, benchmarkPolicy)
	if err != nil {
		return PerformanceBenchmarkAuditSummary{}, err
	}
	if !knownPerformanceStages[stage] {
		return PerformanceBenchmarkAuditSummary{}, fmt.Errorf("unknown performance benchmark stage %q", stage)
	}
	env, ok := performanceEnvironmentByID(registry, environmentID)
	if !ok {
		return PerformanceBenchmarkAuditSummary{}, fmt.Errorf("unknown performance environment %s", environmentID)
	}

	var changed []string
	if strings.TrimSpace(changesFile) != "" {
		changed, err = readChangedPaths(root, changesFile)
		if err != nil {
			return PerformanceBenchmarkAuditSummary{}, err
		}
	}

	sourceSHA, err := currentSourceCommit(root)
	if err != nil {
		return PerformanceBenchmarkAuditSummary{}, err
	}
	summary := PerformanceBenchmarkAuditSummary{
		Stage: stage, EnvironmentID: environmentID, SourceSHA: sourceSHA,
		Registered: len(registry.Targets), ByScope: map[string]int{},
	}
	if len(registry.Targets) == 0 {
		summary.Status = "not-applicable"
		return summary, nil
	}

	workloads := benchmarkWorkloadIndex(benchmarkPolicy)
	selected := selectPerformanceTargets(registry.Targets, stage, environmentID, changed)
	summary.Selected = len(selected)
	if len(selected) == 0 {
		summary.Status = "not-selected"
		return summary, nil
	}
	cachedResults, _, err := reusablePerformanceResults(root)
	if err != nil {
		return summary, err
	}

	for _, target := range selected {
		sourceDigest, err := performanceTargetSourceDigest(root, target)
		if err != nil {
			return summary, fmt.Errorf("performance target %s source identity: %w", target.ID, err)
		}
		if cached, ok := cachedResults[benchmarkBaselineKey(target.ID, environmentID)]; ok {
			summary.Reused++
			summary.Metrics += len(cached.Metrics)
			summary.ByScope[target.ScopeKind]++
			summary.Results = append(summary.Results, cached)
			continue
		}
		definition, ok := workloads[target.WorkloadKey]
		if !ok {
			return summary, fmt.Errorf("performance target %s has no benchmark workload definition for %s", target.ID, target.WorkloadKey)
		}
		result, err := executePerformanceTarget(root, sourceDigest, target, env, definition, baselines)
		if err != nil {
			return summary, err
		}
		summary.Executed++
		summary.Metrics += len(result.Metrics)
		summary.ByScope[target.ScopeKind]++
		summary.Results = append(summary.Results, result)
	}
	summary.Status = "pass"
	return summary, nil
}

func loadValidatedPerformanceConfiguration(root string) (PerformancePolicy, PerformanceRegistry, ArchitectureRegistry, error) {
	var policy PerformancePolicy
	var registry PerformanceRegistry
	var architecture ArchitectureRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(performancePolicyPath)), &policy); err != nil {
		return policy, registry, architecture, err
	}
	if err := validatePerformancePolicy(policy); err != nil {
		return policy, registry, architecture, err
	}
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(performanceRegistryPath)), &registry); err != nil {
		return policy, registry, architecture, err
	}
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &architecture); err != nil {
		return policy, registry, architecture, err
	}
	if err := validateArchitectureRegistry(architecture); err != nil {
		return policy, registry, architecture, err
	}
	return policy, registry, architecture, nil
}

func loadBenchmarkPolicy(root string) (BenchmarkPolicy, error) {
	var policy BenchmarkPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(benchmarkPolicyPath)), &policy); err != nil {
		return policy, err
	}
	if policy.SchemaVersion != 1 || policy.PolicyKind != "compiled-performance-workloads" {
		return policy, fmt.Errorf("invalid benchmark policy header")
	}
	if len(policy.Workloads) != len(knownPerformanceWorkloadKeys) {
		return policy, fmt.Errorf("benchmark workload registry mismatch: expected %d, got %d", len(knownPerformanceWorkloadKeys), len(policy.Workloads))
	}
	seen := map[string]bool{}
	for _, workload := range policy.Workloads {
		if !knownPerformanceWorkloadKeys[workload.Key] || seen[workload.Key] {
			return policy, fmt.Errorf("unknown or duplicate benchmark workload %q", workload.Key)
		}
		seen[workload.Key] = true
		if !knownBenchmarkHandlerKeys[workload.HandlerKey] {
			return policy, fmt.Errorf("benchmark workload %s uses unknown handler %s", workload.Key, workload.HandlerKey)
		}
		if workload.WarmupIterations < 0 || workload.WarmupIterations > 20 {
			return policy, fmt.Errorf("benchmark workload %s has invalid warmup_iterations", workload.Key)
		}
		if workload.SampleCount < 3 || workload.SampleCount > 100 {
			return policy, fmt.Errorf("benchmark workload %s has invalid sample_count", workload.Key)
		}
		if workload.MaxSampleDurationMS < 1 || workload.MaxSampleDurationMS > 60000 {
			return policy, fmt.Errorf("benchmark workload %s has invalid max_sample_duration_ms", workload.Key)
		}
	}
	return policy, nil
}

func loadBenchmarkBaselines(root string, registry PerformanceRegistry, policy BenchmarkPolicy) (map[string]PerformanceBaseline, error) {
	var baselines PerformanceBaselineRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(benchmarkBaselinePath)), &baselines); err != nil {
		return nil, err
	}
	if baselines.SchemaVersion != 1 || baselines.RegistryKind != "performance-baselines" {
		return nil, fmt.Errorf("invalid performance baseline registry header")
	}
	targets := map[string]PerformanceTarget{}
	for _, target := range registry.Targets {
		targets[target.ID] = target
	}
	envs := map[string]PerformanceEnvironment{}
	for _, env := range registry.Environments {
		envs[env.ID] = env
	}
	workloads := benchmarkWorkloadIndex(policy)
	index := map[string]PerformanceBaseline{}
	for _, baseline := range baselines.Records {
		target, ok := targets[baseline.TargetID]
		if !ok {
			return nil, fmt.Errorf("performance baseline %s references unknown target %s", baseline.ID, baseline.TargetID)
		}
		env, ok := envs[baseline.EnvironmentID]
		if !ok {
			return nil, fmt.Errorf("performance baseline %s references unknown environment %s", baseline.ID, baseline.EnvironmentID)
		}
		workload, ok := workloads[target.WorkloadKey]
		if !ok {
			return nil, fmt.Errorf("performance baseline %s references target without workload definition", baseline.ID)
		}
		if _, err := validateFullCommitID(baseline.SourceSHA); err != nil {
			return nil, fmt.Errorf("performance baseline %s source_sha: %w", baseline.ID, err)
		}
		expectedToolchain := benchmarkToolchain(workload)
		if baseline.Toolchain != expectedToolchain {
			return nil, fmt.Errorf("performance baseline %s uses stale toolchain %s; expected %s", baseline.ID, baseline.Toolchain, expectedToolchain)
		}
		if baseline.TargetDigest != digestCanonical(target) ||
			baseline.EnvironmentDigest != digestCanonical(env) ||
			baseline.WorkloadDigest != digestCanonical(workload) {
			return nil, fmt.Errorf("performance baseline %s identity no longer matches current target/environment/workload", baseline.ID)
		}
		key := benchmarkBaselineKey(baseline.TargetID, baseline.EnvironmentID)
		if _, exists := index[key]; exists {
			return nil, fmt.Errorf("duplicate performance baseline for %s/%s", baseline.TargetID, baseline.EnvironmentID)
		}
		index[key] = baseline
	}
	return index, nil
}

func executePerformanceTarget(root, sourceSHA string, target PerformanceTarget, env PerformanceEnvironment, workload BenchmarkWorkloadDefinition, baselines map[string]PerformanceBaseline) (PerformanceBenchmarkResult, error) {
	handler, cleanup, err := prepareBenchmarkHandler(root, workload)
	if err != nil {
		return PerformanceBenchmarkResult{}, fmt.Errorf("performance target %s prepare workload: %w", target.ID, err)
	}
	defer cleanup()

	for i := 0; i < workload.WarmupIterations; i++ {
		if _, err := handler(); err != nil {
			return PerformanceBenchmarkResult{}, fmt.Errorf("performance target %s warmup %d: %w", target.ID, i+1, err)
		}
	}
	samples := make([]BenchmarkSample, 0, workload.SampleCount)
	for i := 0; i < workload.SampleCount; i++ {
		started := time.Now()
		sample, err := handler()
		if err != nil {
			return PerformanceBenchmarkResult{}, fmt.Errorf("performance target %s sample %d: %w", target.ID, i+1, err)
		}
		elapsed := time.Since(started)
		if elapsed > time.Duration(workload.MaxSampleDurationMS)*time.Millisecond {
			return PerformanceBenchmarkResult{}, fmt.Errorf("performance target %s sample %d exceeded %dms cap", target.ID, i+1, workload.MaxSampleDurationMS)
		}
		samples = append(samples, sample)
	}

	var baseline *PerformanceBaseline
	if value, ok := baselines[benchmarkBaselineKey(target.ID, env.ID)]; ok {
		copy := value
		baseline = &copy
	}
	metrics, err := evaluatePerformanceMetrics(target, env, samples, baseline)
	if err != nil {
		return PerformanceBenchmarkResult{}, err
	}
	return PerformanceBenchmarkResult{
		TargetID: target.ID, EnvironmentID: env.ID, SourceSHA: sourceSHA,
		Toolchain:    benchmarkToolchain(workload),
		TargetDigest: digestCanonical(target), EnvironmentDigest: digestCanonical(env),
		WorkloadDigest: digestCanonical(workload), Samples: len(samples), Metrics: metrics,
	}, nil
}

type benchmarkSampleHandler func() (BenchmarkSample, error)

func prepareBenchmarkHandler(root string, workload BenchmarkWorkloadDefinition) (benchmarkSampleHandler, func(), error) {
	switch workload.HandlerKey {
	case "builtin-cmdr-dev-metadata-audit-v1":
		return func() (BenchmarkSample, error) {
			started := time.Now()
			if _, err := runPerformanceRegistryAudit(root); err != nil {
				return BenchmarkSample{}, err
			}
			return BenchmarkSample{Measurements: []BenchmarkMeasurement{{
				Kind: "latency", Unit: "ns", Value: float64(time.Since(started).Nanoseconds()),
			}}}, nil
		}, func() {}, nil
	case "builtin-pilot-context-projection-v1":
		probe, cleanup, err := preparePilotProjectionProbe(root)
		if err != nil {
			return nil, func() {}, err
		}
		return func() (BenchmarkSample, error) {
			observation, err := probe.run(20_000, "latency", 0)
			if err != nil {
				return BenchmarkSample{}, err
			}
			return BenchmarkSample{Measurements: []BenchmarkMeasurement{{
				Kind: "latency", Unit: "ns", Value: observation.NSPerOperation,
			}}}, nil
		}, cleanup, nil
	case "builtin-event-search-validation-v1":
		probe, cleanup, err := prepareEventSearchValidationProbe(root)
		if err != nil {
			return nil, func() {}, err
		}
		return func() (BenchmarkSample, error) {
			observation, err := probe.run(20_000)
			if err != nil {
				return BenchmarkSample{}, err
			}
			return BenchmarkSample{Measurements: []BenchmarkMeasurement{{
				Kind: "latency", Unit: "ns", Value: observation.NSPerOperation,
			}}}, nil
		}, cleanup, nil
	case "builtin-event-search-orchestration-v1":
		probe, cleanup, err := prepareEventSearchValidationProbe(root)
		if err != nil {
			return nil, func() {}, err
		}
		return func() (BenchmarkSample, error) {
			observation, err := probe.runOperation(20_000, "orchestration", "latency", 0)
			if err != nil {
				return BenchmarkSample{}, err
			}
			return BenchmarkSample{Measurements: []BenchmarkMeasurement{{
				Kind: "latency", Unit: "ns", Value: observation.NSPerOperation,
			}}}, nil
		}, cleanup, nil
	case "builtin-event-search-frontend-state-v1":
		probe, cleanup, err := prepareEventSearchFrontendProbe(root)
		if err != nil {
			return nil, func() {}, err
		}
		return func() (BenchmarkSample, error) {
			observation, err := probe.run(2_000, "latency", 0)
			if err != nil {
				return BenchmarkSample{}, err
			}
			return BenchmarkSample{Measurements: []BenchmarkMeasurement{{
				Kind: "latency", Unit: "ns", Value: observation.NSPerOperation,
			}}}, nil
		}, cleanup, nil
	case "builtin-event-inspection-preimplementation-v1":
		return nil, func() {}, fmt.Errorf("event-inspection benchmark is blocked while event-inspection-runtime is preimplementation; implement E10-INV-004B-RUNTIME before executable measurement")
	default:
		return nil, func() {}, fmt.Errorf("unsupported benchmark handler %s", workload.HandlerKey)
	}
}

func benchmarkToolchain(workload BenchmarkWorkloadDefinition) string {
	if workload.HandlerKey == "builtin-event-search-frontend-state-v1" {
		return "node" + engineeringNodeVersion
	}
	return "go" + engineeringGoVersion
}

type eventSearchFrontendProbe struct {
	moduleRoot string
	scriptPath string
}

func prepareEventSearchFrontendProbe(root string) (eventSearchFrontendProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/event-search-frontend", false)
	if err != nil {
		return eventSearchFrontendProbe{}, func() {}, fmt.Errorf("event-search frontend runtime is not implemented: %w", err)
	}
	scriptPath, err := resolveRepoPath(root, "engineering/performance/event-search/frontend-runtime-probe.mjs", false)
	if err != nil {
		return eventSearchFrontendProbe{}, func() {}, fmt.Errorf("event-search frontend performance probe is unavailable: %w", err)
	}
	version, err := runPerformanceProcess(moduleRoot, "node", "--version")
	if err != nil {
		return eventSearchFrontendProbe{}, func() {}, err
	}
	if strings.TrimSpace(version) != "v"+engineeringNodeVersion {
		return eventSearchFrontendProbe{}, func() {}, fmt.Errorf("event-search frontend benchmark requires node v%s, got %s", engineeringNodeVersion, version)
	}
	return eventSearchFrontendProbe{moduleRoot: moduleRoot, scriptPath: scriptPath}, func() {}, nil
}

func (probe eventSearchFrontendProbe) run(iterations int, mode string, memoryLimitMiB int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 1_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Event Search frontend iterations must be within 1..1000000")
	}
	if mode != "latency" && mode != "resource" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("unknown Event Search frontend probe mode %q", mode)
	}
	args := []string{probe.scriptPath, "--iterations", strconv.Itoa(iterations), "--mode", mode}
	if mode == "resource" && memoryLimitMiB > 0 {
		args = append([]string{"--max-old-space-size=" + strconv.Itoa(memoryLimitMiB)}, args...)
	}
	output, err := runPerformanceProcess(probe.moduleRoot, "node", args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(output), &observation); err != nil {
		return observation, fmt.Errorf("decode Event Search frontend probe: %w", err)
	}
	if observation.Operation != "frontend-state-update" || observation.Mode != mode ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Event Search frontend observation")
	}
	return observation, nil
}

type eventSearchValidationProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareEventSearchValidationProbe(root string) (eventSearchValidationProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/event-search", false)
	if err != nil {
		return eventSearchValidationProbe{}, func() {}, fmt.Errorf("event-search validation runtime is not implemented: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-event-search-perf-")
	if err != nil {
		return eventSearchValidationProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- tempDir is created by os.MkdirTemp beneath the validated repository root.
	}
	binaryName := "event-search-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binaryPath, "./cmd/perf-probe"); err != nil {
		cleanup()
		return eventSearchValidationProbe{}, func() {}, fmt.Errorf("build Event Search validation probe: %w", err)
	}
	return eventSearchValidationProbe{moduleRoot: moduleRoot, binaryPath: binaryPath}, cleanup, nil
}

func (probe eventSearchValidationProbe) run(iterations int) (pilotProjectionProbeObservation, error) {
	return probe.runOperation(iterations, "validation", "latency", 0)
}

func (probe eventSearchValidationProbe) runOperation(iterations int, operation, mode string, memoryLimitMiB int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 10_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Event Search iterations must be within 1..10000000")
	}
	if operation != "validation" && operation != "orchestration" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("unknown Event Search probe operation %q", operation)
	}
	args := []string{
		"-iterations", strconv.Itoa(iterations),
		"-mode", mode,
		"-operation", operation,
	}
	if mode == "resource" {
		args = append(args, "-memory-limit-mib", strconv.Itoa(memoryLimitMiB))
	}
	output, err := runPerformanceProcess(probe.moduleRoot, probe.binaryPath, args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(output), &observation); err != nil {
		return observation, fmt.Errorf("decode Event Search %s probe: %w", operation, err)
	}
	if observation.Operation != operation || observation.Mode != mode || observation.Operations != iterations ||
		observation.ElapsedNS <= 0 || observation.NSPerOperation <= 0 ||
		math.IsNaN(observation.NSPerOperation) || math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Event Search %s observation", operation)
	}
	return observation, nil
}

func preparePilotProjectionProbe(root string) (pilotProjectionProbe, func(), error) {
	moduleRoot, err := resolveRepoPath(root, "product-runtime/context-envelope", false)
	if err != nil {
		return pilotProjectionProbe{}, func() {}, err
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-pilot-perf-")
	if err != nil {
		return pilotProjectionProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- tempDir is created by os.MkdirTemp beneath the validated repository root.
	}
	binaryName := "context-envelope-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(moduleRoot, "go", "build", "-trimpath", "-o", binaryPath, "./cmd/perf-probe"); err != nil {
		cleanup()
		return pilotProjectionProbe{}, func() {}, fmt.Errorf("build pilot projection probe: %w", err)
	}
	return pilotProjectionProbe{moduleRoot: moduleRoot, binaryPath: binaryPath}, cleanup, nil
}

func runPerformanceProcess(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...) // #nosec G204,G702 -- callers use either the constant Go executable or a probe binary compiled under the validated repository root; no shell is used.
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		diagnostic := strings.TrimSpace(stderr.String())
		if diagnostic == "" {
			diagnostic = strings.TrimSpace(stdout.String())
		}
		if diagnostic == "" {
			diagnostic = err.Error()
		}
		return "", fmt.Errorf("performance process failed: %s", diagnostic)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (probe pilotProjectionProbe) run(iterations int, mode string, memoryLimitMiB int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 10_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("pilot projection iterations must be within 1..10000000")
	}
	args := []string{"-iterations", strconv.Itoa(iterations), "-mode", mode}
	if mode == "resource" {
		args = append(args, "-memory-limit-mib", strconv.Itoa(memoryLimitMiB))
	}
	output, err := runPerformanceProcess(probe.moduleRoot, probe.binaryPath, args...)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(output), &observation); err != nil {
		return observation, fmt.Errorf("decode pilot projection probe: %w", err)
	}
	if observation.Mode != mode || observation.Operations != iterations ||
		observation.ElapsedNS <= 0 || observation.NSPerOperation <= 0 ||
		math.IsNaN(observation.NSPerOperation) || math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid pilot projection observation")
	}
	return observation, nil
}

func evaluatePerformanceMetrics(target PerformanceTarget, env PerformanceEnvironment, samples []BenchmarkSample, baseline *PerformanceBaseline) ([]BenchmarkMetricResult, error) {
	if len(samples) == 0 {
		return nil, fmt.Errorf("performance target %s has no samples", target.ID)
	}
	var results []BenchmarkMetricResult
	for _, metric := range target.Metrics {
		values := make([]float64, 0, len(samples))
		for i, sample := range samples {
			value, err := measurementForMetric(sample, metric)
			if err != nil {
				return nil, fmt.Errorf("performance target %s metric %s sample %d: %w", target.ID, metric.ID, i+1, err)
			}
			values = append(values, value)
		}
		observed, err := aggregateBenchmarkValues(values, metric.Aggregation)
		if err != nil {
			return nil, err
		}
		result := BenchmarkMetricResult{
			ID: metric.ID, Kind: metric.Kind, Unit: metric.Unit, Aggregation: metric.Aggregation,
			Observed: observed, Comparator: metric.Comparator, Budget: metric.Budget,
			MaxRelativeRegressionPercent: metric.MaxRelativeRegressionPercent, RelativePass: true,
		}
		switch metric.Comparator {
		case "upper-bound":
			result.AbsolutePass = observed <= metric.Budget
		case "lower-bound":
			result.AbsolutePass = observed >= metric.Budget
		default:
			return nil, fmt.Errorf("performance target %s metric %s has unknown comparator %s", target.ID, metric.ID, metric.Comparator)
		}
		if !result.AbsolutePass {
			return nil, fmt.Errorf("performance target %s metric %s violated absolute %s budget: observed %.6g %s, budget %.6g %s", target.ID, metric.ID, metric.Comparator, observed, metric.Unit, metric.Budget, metric.Unit)
		}
		if metric.MaxRelativeRegressionPercent > 0 {
			if !env.RelativeRegressionAllowed {
				return nil, fmt.Errorf("performance target %s metric %s cannot evaluate relative regression in environment %s", target.ID, metric.ID, env.ID)
			}
			if baseline == nil {
				return nil, fmt.Errorf("performance target %s metric %s requires a matching baseline", target.ID, metric.ID)
			}
			base, ok := baseline.MetricValues[metric.ID]
			if !ok || base <= 0 || math.IsNaN(base) || math.IsInf(base, 0) {
				return nil, fmt.Errorf("performance baseline for %s is missing valid metric %s", target.ID, metric.ID)
			}
			result.Baseline = base
			switch metric.Comparator {
			case "upper-bound":
				result.RelativeRegressionPercent = (observed - base) / base * 100
			case "lower-bound":
				result.RelativeRegressionPercent = (base - observed) / base * 100
			}
			if result.RelativeRegressionPercent < 0 {
				result.RelativeRegressionPercent = 0
			}
			result.RelativePass = result.RelativeRegressionPercent <= metric.MaxRelativeRegressionPercent
			if !result.RelativePass {
				return nil, fmt.Errorf("performance target %s metric %s regressed %.3f%% over baseline; allowed %.3f%%", target.ID, metric.ID, result.RelativeRegressionPercent, metric.MaxRelativeRegressionPercent)
			}
		}
		results = append(results, result)
	}
	return results, nil
}

func measurementForMetric(sample BenchmarkSample, metric PerformanceMetricBudget) (float64, error) {
	for _, measurement := range sample.Measurements {
		if measurement.Kind != metric.Kind {
			continue
		}
		value, ok := convertPerformanceMeasurement(metric.Kind, measurement.Unit, metric.Unit, measurement.Value)
		if ok {
			return value, nil
		}
	}
	return 0, fmt.Errorf("no measurement convertible to %s/%s", metric.Kind, metric.Unit)
}

func convertPerformanceMeasurement(kind, from, to string, value float64) (float64, bool) {
	if from == to {
		return value, true
	}
	switch kind {
	case "latency":
		factors := map[string]float64{"ns": 1, "us": 1e3, "ms": 1e6, "s": 1e9}
		fromFactor, okFrom := factors[from]
		toFactor, okTo := factors[to]
		if okFrom && okTo {
			return value * fromFactor / toFactor, true
		}
	case "cpu-time":
		factors := map[string]float64{"ns/op": 1, "us/op": 1e3, "ms/op": 1e6}
		fromFactor, okFrom := factors[from]
		toFactor, okTo := factors[to]
		if okFrom && okTo {
			return value * fromFactor / toFactor, true
		}
	case "memory":
		factors := map[string]float64{"bytes": 1, "MiB": 1024 * 1024}
		fromFactor, okFrom := factors[from]
		toFactor, okTo := factors[to]
		if okFrom && okTo {
			return value * fromFactor / toFactor, true
		}
	}
	return 0, false
}

func aggregateBenchmarkValues(values []float64, aggregation string) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("cannot aggregate empty benchmark values")
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	switch aggregation {
	case "median":
		mid := len(sorted) / 2
		if len(sorted)%2 == 1 {
			return sorted[mid], nil
		}
		return (sorted[mid-1] + sorted[mid]) / 2, nil
	case "p95":
		return nearestRank(sorted, 0.95), nil
	case "p99":
		return nearestRank(sorted, 0.99), nil
	case "max":
		return sorted[len(sorted)-1], nil
	case "mean":
		var total float64
		for _, value := range sorted {
			total += value
		}
		return total / float64(len(sorted)), nil
	default:
		return 0, fmt.Errorf("unknown benchmark aggregation %s", aggregation)
	}
}

func nearestRank(sorted []float64, percentile float64) float64 {
	index := int(math.Ceil(percentile*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func selectPerformanceTargets(targets []PerformanceTarget, stage, environmentID string, changed []string) []PerformanceTarget {
	var selected []PerformanceTarget
	for _, target := range targets {
		if target.EnvironmentID != environmentID || !containsString(target.Stages, stage) {
			continue
		}
		if len(changed) > 0 {
			matched := false
			for _, path := range changed {
				for _, pattern := range target.TriggerPaths {
					if triggerPatternMatches(pattern, path) {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				continue
			}
		}
		selected = append(selected, target)
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].ID < selected[j].ID })
	return selected
}

func formatPerformanceBenchmarkEvidence(summary PerformanceBenchmarkAuditSummary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "stage=%s environment=%s registered=%d selected=%d executed=%d metrics=%d status=%s source=%s",
		summary.Stage, summary.EnvironmentID, summary.Registered, summary.Selected,
		summary.Executed, summary.Metrics, summary.Status, summary.SourceSHA)
	for _, result := range summary.Results {
		for _, metric := range result.Metrics {
			fmt.Fprintf(&b, " target=%s metric=%s observed=%.9g%s budget=%.9g%s absolute_pass=%t relative_pass=%t",
				result.TargetID, metric.ID, metric.Observed, metric.Unit, metric.Budget, metric.Unit,
				metric.AbsolutePass, metric.RelativePass)
		}
	}
	return b.String()
}

func performanceEnvironmentByID(registry PerformanceRegistry, id string) (PerformanceEnvironment, bool) {
	for _, env := range registry.Environments {
		if env.ID == id {
			return env, true
		}
	}
	return PerformanceEnvironment{}, false
}

func benchmarkWorkloadIndex(policy BenchmarkPolicy) map[string]BenchmarkWorkloadDefinition {
	index := map[string]BenchmarkWorkloadDefinition{}
	for _, workload := range policy.Workloads {
		index[workload.Key] = workload
	}
	return index
}

func benchmarkBaselineKey(targetID, environmentID string) string {
	return targetID + "\x00" + environmentID
}
