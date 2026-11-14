package main

import (
	"errors"
	"testing"
)

func TestMainlineHealthGateConstants(t *testing.T) {
	if StatusMainHealthy != "MAIN_HEALTHY" {
		t.Fatalf("expected MAIN_HEALTHY, got %s", StatusMainHealthy)
	}
	if StatusMainUnhealthy != "MAIN_UNHEALTHY" {
		t.Fatalf("expected MAIN_UNHEALTHY, got %s", StatusMainUnhealthy)
	}
	if StatusMainUnverified != "MAIN_UNVERIFIED" {
		t.Fatalf("expected MAIN_UNVERIFIED, got %s", StatusMainUnverified)
	}

	if ExitCodeMainHealthy != 0 {
		t.Fatalf("expected ExitCodeMainHealthy == 0, got %d", ExitCodeMainHealthy)
	}
	if ExitCodeMainUnhealthy != 1 {
		t.Fatalf("expected ExitCodeMainUnhealthy == 1, got %d", ExitCodeMainUnhealthy)
	}
	if ExitCodeMainUnverified != 2 {
		t.Fatalf("expected ExitCodeMainUnverified == 2, got %d", ExitCodeMainUnverified)
	}
}

func TestEvaluateCodeDimension_RejectsStateError(t *testing.T) {
	dim := evaluateCodeDimension(".", CurrentState{}, WorkGraph{}, errors.New("simulated state error"))
	if dim.Passed {
		t.Fatalf("expected Passed == false on state error, got true")
	}
	if dim.Status != "FAIL" {
		t.Fatalf("expected Status == FAIL, got %s", dim.Status)
	}
	if len(dim.Errors) == 0 {
		t.Fatalf("expected non-empty errors list")
	}
}

func TestEvaluateTraceabilityDimension_RejectsStateError(t *testing.T) {
	dim := evaluateTraceabilityDimension(".", CurrentState{}, WorkGraph{}, errors.New("simulated state error"))
	if dim.Passed {
		t.Fatalf("expected Passed == false on state error, got true")
	}
	if dim.Status != "FAIL" {
		t.Fatalf("expected Status == FAIL, got %s", dim.Status)
	}
}

func TestEvaluateOperationalQualityDimension_RejectsStateError(t *testing.T) {
	dim := evaluateOperationalQualityDimension(".", WorkGraph{}, errors.New("simulated state error"))
	if dim.Passed {
		t.Fatalf("expected Passed == false on state error, got true")
	}
	if dim.Status != "FAIL" {
		t.Fatalf("expected Status == FAIL, got %s", dim.Status)
	}
}

func TestRunMainlineHealthGate_OverallStatusAggregation(t *testing.T) {
	// Case 1: When state error is present, overall status must be MAIN_UNHEALTHY with exit code 1
	report, exitCode := runMainlineHealthGate(".", CurrentState{}, WorkGraph{}, errors.New("state decode failed"))
	if report.Status != StatusMainUnhealthy {
		t.Fatalf("expected StatusMainUnhealthy on state failure, got %s", report.Status)
	}
	if exitCode != ExitCodeMainUnhealthy {
		t.Fatalf("expected exit code %d, got %d", ExitCodeMainUnhealthy, exitCode)
	}
	if len(report.Dimensions) != 5 {
		t.Fatalf("expected 5 evaluated dimensions, got %d", len(report.Dimensions))
	}
}

func TestPrintMainlineHealthGateReport_FormatsWithoutPanic(t *testing.T) {
	report := MainlineHealthGateReport{
		Status:     StatusMainHealthy,
		ExitCode:   ExitCodeMainHealthy,
		Summary:    "All systems healthy",
		Timestamp:  "2026-11-14T22:00:00Z",
		Root:       "/test/root",
		HeadCommit: "1234567890abcdef1234567890abcdef12345678",
		Dimensions: map[string]HealthDimensionResult{
			"code": {
				Name:    "code",
				Passed:  true,
				Status:  "PASS",
				Details: []string{"all checks passed"},
			},
			"security": {
				Name:    "security",
				Passed:  true,
				Status:  "PASS",
				Details: []string{"0 secrets"},
			},
		},
	}
	// Call printMainlineHealthGateReport to ensure formatting does not panic
	printMainlineHealthGateReport(report)
}
