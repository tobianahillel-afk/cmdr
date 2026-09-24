package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

const deepPerformancePolicyPath = "engineering/performance/deep-performance-policy.json"

type DeepPerformanceCaps struct {
	MaxTargetDurationSeconds int `json:"max_target_duration_seconds"`
	MaxTargetMemoryMiB       int `json:"max_target_memory_mib"`
	MaxConcurrency           int `json:"max_concurrency"`
	MaxTargetsPerRun         int `json:"max_targets_per_run"`
}

type DeepPerformanceTarget struct {
	ID                  string   `json:"id"`
	PerformanceTargetID string   `json:"performance_target_id"`
	Kind                string   `json:"kind"`
	HandlerKey          string   `json:"handler_key"`
	Stages              []string `json:"stages"`
	TriggerPaths        []string `json:"trigger_paths"`
	DurationSeconds     int      `json:"duration_seconds"`
	TimeoutSeconds      int      `json:"timeout_seconds"`
	MemoryLimitMiB      int      `json:"memory_limit_mib"`
	Concurrency         int      `json:"concurrency"`
	Iterations          int      `json:"iterations"`
	Rationale           string   `json:"rationale"`
}

type DeepPerformancePolicy struct {
	SchemaVersion     int                     `json:"schema_version"`
	DefaultPolicy     string                  `json:"default_policy"`
	EvidenceDirectory string                  `json:"evidence_directory"`
	ControlPlanePaths []string                `json:"control_plane_paths"`
	Caps              DeepPerformanceCaps     `json:"caps"`
	Targets           []DeepPerformanceTarget `json:"targets"`
}

type DeepPerformanceObservation struct {
	DurationMS           int64   `json:"duration_ms"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
	Operations           int64   `json:"operations"`
	ThroughputOpsPerSec  float64 `json:"throughput_ops_per_second"`
	ArtifactKind         string  `json:"artifact_kind,omitempty"`
	ArtifactSHA256       string  `json:"artifact_sha256,omitempty"`
	artifact             []byte
}

type DeepPerformanceResult struct {
	TargetID             string  `json:"target_id"`
	PerformanceTargetID  string  `json:"performance_target_id"`
	Kind                 string  `json:"kind"`
	Stage                string  `json:"stage"`
	SourceSHA            string  `json:"source_sha"`
	DurationMS           int64   `json:"duration_ms"`
	PeakHeapBytes        uint64  `json:"peak_heap_bytes"`
	TotalAllocationBytes uint64  `json:"total_allocation_bytes"`
	Operations           int64   `json:"operations"`
	ThroughputOpsPerSec  float64 `json:"throughput_ops_per_second"`
	ArtifactKind         string  `json:"artifact_kind,omitempty"`
	ArtifactSHA256       string  `json:"artifact_sha256,omitempty"`
	EvidencePath         string  `json:"evidence_path,omitempty"`
	ArtifactPath         string  `json:"artifact_path,omitempty"`
}

type DeepPerformanceAuditSummary struct {
	Stage               string                  `json:"stage"`
	SourceSHA           string                  `json:"source_sha"`
	RuntimeBoundaries   int                     `json:"runtime_boundaries"`
	RegisteredTargets   int                     `json:"registered_targets"`
	SelectedTargets     int                     `json:"selected_targets"`
	ExecutedTargets     int                     `json:"executed_targets"`
	RiskSensitiveChange bool                    `json:"risk_sensitive_change"`
	Status              string                  `json:"status"`
	Results             []DeepPerformanceResult `json:"results,omitempty"`
}

var knownDeepPerformanceKinds = map[string]bool{
	"profile":  true,
	"load":     true,
	"soak":     true,
	"resource": true,
}

var knownDeepPerformanceStages = map[string]bool{
	"nightly":   true,
	"release":   true,
	"on-demand": true,
}

var knownDeepPerformanceHandlerKeys = map[string]bool{
	"builtin-cmdr-dev-deep-metadata-v1":   true,
	"builtin-pilot-context-projection-v1": true,
}

func runDeepPerformanceAudit(root, stage, changesFile, requestedTarget string) (DeepPerformanceAuditSummary, error) {
	_, registry, architecture, err := loadValidatedPerformanceConfiguration(root)
	if err != nil {
		return DeepPerformanceAuditSummary{}, err
	}
	var policy DeepPerformancePolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(deepPerformancePolicyPath)), &policy); err != nil {
		return DeepPerformanceAuditSummary{}, err
	}
	if err := validateDeepPerformancePolicy(policy, registry, architecture); err != nil {
		return DeepPerformanceAuditSummary{}, err
	}
	if stage != "pr" && !knownDeepPerformanceStages[stage] {
		return DeepPerformanceAuditSummary{}, fmt.Errorf("unknown deep performance stage %q", stage)
	}

	sourceSHA, err := currentSourceCommit(root)
	if err != nil {
		return DeepPerformanceAuditSummary{}, err
	}
	summary := DeepPerformanceAuditSummary{
		Stage:             stage,
		SourceSHA:         sourceSHA,
		RegisteredTargets: len(policy.Targets),
	}
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			summary.RuntimeBoundaries++
		}
	}

	var changed []string
	if strings.TrimSpace(changesFile) != "" {
		changed, err = readChangedPaths(root, changesFile)
		if err != nil {
			return summary, err
		}
	}
	summary.RiskSensitiveChange = deepPerformanceSensitiveChange(changed, policy)

	if stage == "pr" {
		if len(policy.Targets) == 0 && summary.RuntimeBoundaries == 0 {
			summary.Status = "not-applicable"
		} else if summary.RiskSensitiveChange {
			summary.Status = "risk-triggered-deep-validation"
		} else {
			summary.Status = "policy-valid-not-triggered"
		}
		return summary, nil
	}

	selected, err := selectDeepPerformanceTargets(policy, stage, requestedTarget)
	if err != nil {
		return summary, err
	}
	summary.SelectedTargets = len(selected)
	if len(policy.Targets) == 0 && summary.RuntimeBoundaries == 0 {
		summary.Status = "not-applicable"
		return summary, nil
	}
	if len(selected) == 0 {
		summary.Status = "not-selected"
		return summary, nil
	}
	if len(selected) > policy.Caps.MaxTargetsPerRun {
		return summary, fmt.Errorf("deep performance selected %d targets, exceeding cap %d", len(selected), policy.Caps.MaxTargetsPerRun)
	}

	for _, target := range selected {
		observation, err := executeDeepPerformanceTarget(root, policy.Caps, target)
		if err != nil {
			return summary, fmt.Errorf("deep performance target %s: %w", target.ID, err)
		}
		result, err := persistDeepPerformanceEvidence(root, policy.EvidenceDirectory, sourceSHA, stage, target, observation)
		if err != nil {
			return summary, err
		}
		summary.ExecutedTargets++
		summary.Results = append(summary.Results, result)
	}
	summary.Status = "pass"
	return summary, nil
}

func validateDeepPerformancePolicy(policy DeepPerformancePolicy, registry PerformanceRegistry, architecture ArchitectureRegistry) error {
	if policy.SchemaVersion != 1 || policy.DefaultPolicy != "deny-unregistered-deep-performance-target" {
		return fmt.Errorf("invalid deep performance policy header")
	}
	if policy.EvidenceDirectory != "engineering/testing/deep-performance-evidence" {
		return fmt.Errorf("deep performance evidence directory must be engineering/testing/deep-performance-evidence")
	}
	if policy.Caps.MaxTargetDurationSeconds < 1 || policy.Caps.MaxTargetDurationSeconds > 1800 {
		return fmt.Errorf("deep performance max target duration must be within 1..1800 seconds")
	}
	if policy.Caps.MaxTargetMemoryMiB < 64 || policy.Caps.MaxTargetMemoryMiB > 8192 {
		return fmt.Errorf("deep performance max target memory must be within 64..8192 MiB")
	}
	if policy.Caps.MaxConcurrency < 1 || policy.Caps.MaxConcurrency > 256 {
		return fmt.Errorf("deep performance max concurrency must be within 1..256")
	}
	if policy.Caps.MaxTargetsPerRun < 1 || policy.Caps.MaxTargetsPerRun > 8 {
		return fmt.Errorf("deep performance max targets per run must be within 1..8")
	}
	if len(policy.ControlPlanePaths) == 0 {
		return fmt.Errorf("deep performance policy requires control-plane paths")
	}
	for _, pattern := range policy.ControlPlanePaths {
		if _, err := normalizeArchitecturePattern(pattern); err != nil {
			return fmt.Errorf("deep performance control-plane path: %w", err)
		}
	}

	performanceTargets := map[string]PerformanceTarget{}
	for _, target := range registry.Targets {
		performanceTargets[target.ID] = target
	}
	runtimeBoundaries := map[string]bool{}
	for _, boundary := range architecture.Boundaries {
		if boundary.Kind == "product-runtime" {
			runtimeBoundaries[boundary.ID] = true
		}
	}

	seen := map[string]bool{}
	deepByPerformanceTarget := map[string]int{}
	for _, target := range policy.Targets {
		if !strings.HasPrefix(target.ID, "PERF-DEEP-") || seen[target.ID] {
			return fmt.Errorf("invalid or duplicate deep performance target id %q", target.ID)
		}
		seen[target.ID] = true
		performanceTarget, ok := performanceTargets[target.PerformanceTargetID]
		if !ok {
			return fmt.Errorf("deep performance target %s references unknown performance target %s", target.ID, target.PerformanceTargetID)
		}
		if !knownDeepPerformanceKinds[target.Kind] {
			return fmt.Errorf("deep performance target %s uses unknown kind %s", target.ID, target.Kind)
		}
		if !knownDeepPerformanceHandlerKeys[target.HandlerKey] {
			return fmt.Errorf("deep performance target %s uses unknown handler %s", target.ID, target.HandlerKey)
		}
		if len(target.Stages) == 0 {
			return fmt.Errorf("deep performance target %s requires stages", target.ID)
		}
		seenStages := map[string]bool{}
		for _, stage := range target.Stages {
			if !knownDeepPerformanceStages[stage] || seenStages[stage] {
				return fmt.Errorf("deep performance target %s has invalid/duplicate stage %s", target.ID, stage)
			}
			seenStages[stage] = true
		}
		if len(target.TriggerPaths) == 0 {
			return fmt.Errorf("deep performance target %s requires trigger_paths", target.ID)
		}
		for _, pattern := range target.TriggerPaths {
			if _, err := normalizeArchitecturePattern(pattern); err != nil {
				return fmt.Errorf("deep performance target %s trigger: %w", target.ID, err)
			}
		}
		if target.DurationSeconds < 0 || target.DurationSeconds > target.TimeoutSeconds {
			return fmt.Errorf("deep performance target %s duration must be between zero and timeout", target.ID)
		}
		if target.TimeoutSeconds < 1 || target.TimeoutSeconds > policy.Caps.MaxTargetDurationSeconds {
			return fmt.Errorf("deep performance target %s timeout exceeds policy cap", target.ID)
		}
		if target.MemoryLimitMiB < 64 || target.MemoryLimitMiB > policy.Caps.MaxTargetMemoryMiB {
			return fmt.Errorf("deep performance target %s memory limit exceeds policy cap", target.ID)
		}
		if target.Concurrency < 1 || target.Concurrency > policy.Caps.MaxConcurrency {
			return fmt.Errorf("deep performance target %s concurrency exceeds policy cap", target.ID)
		}
		if target.Iterations < 1 || target.Iterations > 1000000 {
			return fmt.Errorf("deep performance target %s iterations must be within 1..1000000", target.ID)
		}
		if target.Kind == "load" && target.Concurrency < 2 {
			return fmt.Errorf("deep performance load target %s requires concurrency >= 2", target.ID)
		}
		if target.Kind == "soak" && target.DurationSeconds < 1 {
			return fmt.Errorf("deep performance soak target %s requires positive duration", target.ID)
		}
		if strings.TrimSpace(target.Rationale) == "" {
			return fmt.Errorf("deep performance target %s requires rationale", target.ID)
		}
		deepByPerformanceTarget[target.PerformanceTargetID]++
		if performanceTarget.ScopeKind == "product-runtime" && !runtimeBoundaries[performanceTarget.BoundaryID] {
			return fmt.Errorf("deep performance target %s references performance target outside current runtime boundaries", target.ID)
		}
	}

	for _, target := range registry.Targets {
		if target.ScopeKind == "product-runtime" && deepByPerformanceTarget[target.ID] == 0 {
			return fmt.Errorf("product-runtime performance target %s has no deep performance target", target.ID)
		}
	}
	return nil
}

func selectDeepPerformanceTargets(policy DeepPerformancePolicy, stage, requestedTarget string) ([]DeepPerformanceTarget, error) {
	if stage == "on-demand" && strings.TrimSpace(requestedTarget) != "" {
		for _, target := range policy.Targets {
			if target.ID == requestedTarget {
				if !containsString(target.Stages, "on-demand") {
					return nil, fmt.Errorf("requested deep performance target %s does not allow on-demand execution", requestedTarget)
				}
				return []DeepPerformanceTarget{target}, nil
			}
		}
		return nil, fmt.Errorf("requested deep performance target %s is unknown", requestedTarget)
	}
	if stage != "on-demand" && strings.TrimSpace(requestedTarget) != "" {
		return nil, fmt.Errorf("requested deep performance target is allowed only for on-demand stage")
	}
	var selected []DeepPerformanceTarget
	for _, target := range policy.Targets {
		if containsString(target.Stages, stage) {
			selected = append(selected, target)
		}
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].ID < selected[j].ID })
	return selected, nil
}

func deepPerformanceSensitiveChange(changed []string, policy DeepPerformancePolicy) bool {
	for _, path := range changed {
		for _, pattern := range policy.ControlPlanePaths {
			if triggerPatternMatches(pattern, path) {
				return true
			}
		}
		for _, target := range policy.Targets {
			for _, pattern := range target.TriggerPaths {
				if triggerPatternMatches(pattern, path) {
					return true
				}
			}
		}
	}
	return false
}

type deepPerformanceHandler func(context.Context, string, DeepPerformanceTarget) (DeepPerformanceObservation, error)

func executeDeepPerformanceTarget(root string, caps DeepPerformanceCaps, target DeepPerformanceTarget) (DeepPerformanceObservation, error) {
	return executeDeepPerformanceTargetWithHandler(
		root,
		caps,
		target,
		time.Duration(target.TimeoutSeconds)*time.Second,
		runDeepPerformanceHandler,
	)
}

func executeDeepPerformanceTargetWithHandler(
	root string,
	caps DeepPerformanceCaps,
	target DeepPerformanceTarget,
	timeout time.Duration,
	handler deepPerformanceHandler,
) (DeepPerformanceObservation, error) {
	if target.MemoryLimitMiB > caps.MaxTargetMemoryMiB || target.Concurrency > caps.MaxConcurrency {
		return DeepPerformanceObservation{}, fmt.Errorf("target resource request exceeds deep performance caps")
	}
	if timeout <= 0 || timeout > time.Duration(caps.MaxTargetDurationSeconds)*time.Second {
		return DeepPerformanceObservation{}, fmt.Errorf("target timeout exceeds deep performance cap")
	}
	oldMemoryLimit := debug.SetMemoryLimit(int64(target.MemoryLimitMiB) * 1024 * 1024)
	defer debug.SetMemoryLimit(oldMemoryLimit)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	type outcome struct {
		observation DeepPerformanceObservation
		err         error
	}
	done := make(chan outcome, 1)
	go func() {
		observation, err := handler(ctx, root, target)
		done <- outcome{observation: observation, err: err}
	}()
	select {
	case <-ctx.Done():
		return DeepPerformanceObservation{}, fmt.Errorf("target exceeded %s timeout", timeout)
	case result := <-done:
		if result.err != nil {
			return DeepPerformanceObservation{}, result.err
		}
		if err := validateDeepPerformanceObservation(target, result.observation); err != nil {
			return DeepPerformanceObservation{}, err
		}
		return result.observation, nil
	}
}

func runDeepPerformanceHandler(ctx context.Context, root string, target DeepPerformanceTarget) (DeepPerformanceObservation, error) {
	if target.HandlerKey == "builtin-pilot-context-projection-v1" {
		if err := ctx.Err(); err != nil {
			return DeepPerformanceObservation{}, err
		}
		probe, cleanup, err := preparePilotProjectionProbe(root)
		if err != nil {
			return DeepPerformanceObservation{}, err
		}
		defer cleanup()
		observation, err := probe.run(target.Iterations, "resource", target.MemoryLimitMiB)
		if err != nil {
			return DeepPerformanceObservation{}, err
		}
		if err := ctx.Err(); err != nil {
			return DeepPerformanceObservation{}, err
		}
		result := DeepPerformanceObservation{
			DurationMS:           observation.ElapsedNS / int64(time.Millisecond),
			PeakHeapBytes:        observation.PeakHeapBytes,
			TotalAllocationBytes: observation.TotalAllocationBytes,
			Operations:           int64(observation.Operations),
		}
		if observation.ElapsedNS > 0 {
			result.ThroughputOpsPerSec = float64(observation.Operations) / (float64(observation.ElapsedNS) / float64(time.Second))
		}
		return result, nil
	}
	if target.HandlerKey != "builtin-cmdr-dev-deep-metadata-v1" {
		return DeepPerformanceObservation{}, fmt.Errorf("unsupported deep performance handler %s", target.HandlerKey)
	}
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	peak := before.HeapAlloc
	started := time.Now()
	deadline := started.Add(time.Duration(target.DurationSeconds) * time.Second)
	var operations int64

	for {
		if err := ctx.Err(); err != nil {
			return DeepPerformanceObservation{}, err
		}
		for i := 0; i < target.Concurrency && operations < int64(target.Iterations); i++ {
			if _, err := runPerformanceRegistryAudit(root); err != nil {
				return DeepPerformanceObservation{}, err
			}
			operations++
			var current runtime.MemStats
			runtime.ReadMemStats(&current)
			if current.HeapAlloc > peak {
				peak = current.HeapAlloc
			}
		}
		if operations >= int64(target.Iterations) && (target.DurationSeconds == 0 || time.Now().After(deadline)) {
			break
		}
	}
	elapsed := time.Since(started)
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	observation := DeepPerformanceObservation{
		DurationMS:           elapsed.Milliseconds(),
		PeakHeapBytes:        peak,
		TotalAllocationBytes: after.TotalAlloc - before.TotalAlloc,
		Operations:           operations,
	}
	if elapsed > 0 {
		observation.ThroughputOpsPerSec = float64(operations) / elapsed.Seconds()
	}
	if target.Kind == "profile" {
		var profile bytes.Buffer
		if err := pprof.Lookup("heap").WriteTo(&profile, 0); err != nil {
			return DeepPerformanceObservation{}, fmt.Errorf("write heap profile: %w", err)
		}
		sum := sha256.Sum256(profile.Bytes())
		observation.ArtifactKind = "heap-profile"
		observation.ArtifactSHA256 = hex.EncodeToString(sum[:])
		observation.artifact = append([]byte(nil), profile.Bytes()...)
	}
	return observation, nil
}

func validateDeepPerformanceObservation(target DeepPerformanceTarget, observation DeepPerformanceObservation) error {
	if observation.DurationMS < 0 {
		return fmt.Errorf("negative deep performance duration")
	}
	if observation.DurationMS > int64(target.TimeoutSeconds)*1000 {
		return fmt.Errorf("deep performance duration %dms exceeds timeout %ds", observation.DurationMS, target.TimeoutSeconds)
	}
	memoryMiB, err := strconv.ParseUint(strconv.Itoa(target.MemoryLimitMiB), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid deep performance memory limit: %w", err)
	}
	memoryLimit := memoryMiB * 1024 * 1024
	if observation.PeakHeapBytes > memoryLimit {
		return fmt.Errorf("deep performance peak heap %d exceeds memory limit %d", observation.PeakHeapBytes, memoryLimit)
	}
	if observation.Operations < 1 {
		return fmt.Errorf("deep performance observation contains no operations")
	}
	if target.Kind == "soak" && observation.DurationMS < int64(target.DurationSeconds)*1000 {
		return fmt.Errorf("deep performance soak duration %dms is below required %ds", observation.DurationMS, target.DurationSeconds)
	}
	if target.Kind == "profile" {
		if observation.ArtifactKind != "heap-profile" || len(observation.ArtifactSHA256) != 64 || len(observation.artifact) == 0 {
			return fmt.Errorf("profile target requires a hashed heap-profile artifact")
		}
	}
	return nil
}

func persistDeepPerformanceEvidence(
	root, directory, sourceSHA, stage string,
	target DeepPerformanceTarget,
	observation DeepPerformanceObservation,
) (DeepPerformanceResult, error) {
	result := DeepPerformanceResult{
		TargetID: target.ID, PerformanceTargetID: target.PerformanceTargetID,
		Kind: target.Kind, Stage: stage, SourceSHA: sourceSHA,
		DurationMS: observation.DurationMS, PeakHeapBytes: observation.PeakHeapBytes,
		TotalAllocationBytes: observation.TotalAllocationBytes, Operations: observation.Operations,
		ThroughputOpsPerSec: observation.ThroughputOpsPerSec,
		ArtifactKind:        observation.ArtifactKind, ArtifactSHA256: observation.ArtifactSHA256,
	}
	if len(observation.artifact) > 0 {
		artifactPath := filepath.ToSlash(filepath.Join(directory, target.ID+".pprof"))
		if _, err := writeRepoFile(root, artifactPath, observation.artifact); err != nil {
			return result, fmt.Errorf("write deep performance profile: %w", err)
		}
		result.ArtifactPath = artifactPath
	}
	evidencePath := filepath.ToSlash(filepath.Join(directory, target.ID+".json"))
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return result, err
	}
	data = append(data, '\n')
	if _, err := writeRepoFile(root, evidencePath, data); err != nil {
		return result, fmt.Errorf("write deep performance evidence: %w", err)
	}
	result.EvidencePath = evidencePath
	return result, nil
}
