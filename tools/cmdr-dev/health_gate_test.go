package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestRunMainlineHealthGate_BlocksOnAnyDimensionFailure(t *testing.T) {
	dims := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, failedDim := range dims {
		t.Run("fails_on_"+failedDim, func(t *testing.T) {
			report := MainlineHealthGateReport{
				Root:       ".",
				Dimensions: make(map[string]HealthDimensionResult, 5),
			}
			for _, d := range dims {
				if d == failedDim {
					report.Dimensions[d] = HealthDimensionResult{
						Name:   d,
						Passed: false,
						Status: "FAIL",
						Errors: []string{"simulated dimension failure"},
					}
				} else {
					report.Dimensions[d] = HealthDimensionResult{
						Name:    d,
						Passed:  true,
						Status:  "PASS",
						Details: []string{"ok"},
					}
				}
			}

			status, exitCode := aggregateHealthGateReport(&report)

			if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
				t.Fatalf("dimension %s failure did not return StatusMainUnhealthy: status=%s exit=%d", failedDim, status, exitCode)
			}
			if report.Status != StatusMainUnhealthy || report.ExitCode != ExitCodeMainUnhealthy {
				t.Fatalf("dimension %s failure did not set report fields: status=%s exit=%d", failedDim, report.Status, report.ExitCode)
			}
			if !strings.Contains(report.Summary, "FAILED") {
				t.Fatalf("dimension %s failure did not produce failure summary, got: %s", failedDim, report.Summary)
			}
		})

		t.Run("fails_on_unpassed_"+failedDim, func(t *testing.T) {
			report := MainlineHealthGateReport{
				Root:       ".",
				Dimensions: make(map[string]HealthDimensionResult, 5),
			}
			for _, d := range dims {
				if d == failedDim {
					report.Dimensions[d] = HealthDimensionResult{
						Name:   d,
						Passed: false,
						Status: "UNKNOWN",
						Errors: []string{"simulated unknown unpassed state"},
					}
				} else {
					report.Dimensions[d] = HealthDimensionResult{
						Name:    d,
						Passed:  true,
						Status:  "PASS",
						Details: []string{"ok"},
					}
				}
			}

			status, exitCode := aggregateHealthGateReport(&report)

			if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
				t.Fatalf("unpassed dimension %s did not result in StatusMainUnhealthy: status=%s exit=%d", failedDim, status, exitCode)
			}
			if report.Status != StatusMainUnhealthy || report.ExitCode != ExitCodeMainUnhealthy {
				t.Fatalf("unpassed dimension %s did not set report fields: status=%s exit=%d", failedDim, report.Status, report.ExitCode)
			}
		})
	}
}

func TestRunMainlineHealthGate_BlocksOnUnverifiedDimension(t *testing.T) {
	dims := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, unverifiedDim := range dims {
		t.Run("unverified_on_"+unverifiedDim, func(t *testing.T) {
			report := MainlineHealthGateReport{
				Root:       ".",
				Dimensions: make(map[string]HealthDimensionResult, 5),
			}
			for _, d := range dims {
				if d == unverifiedDim {
					report.Dimensions[d] = HealthDimensionResult{
						Name:   d,
						Passed: false,
						Status: "UNVERIFIED",
						Errors: []string{"simulated unverified condition"},
					}
				} else {
					report.Dimensions[d] = HealthDimensionResult{
						Name:    d,
						Passed:  true,
						Status:  "PASS",
						Details: []string{"ok"},
					}
				}
			}

			status, exitCode := aggregateHealthGateReport(&report)

			if status != StatusMainUnverified || exitCode != ExitCodeMainUnverified {
				t.Fatalf("unverified dimension %s did not result in StatusMainUnverified: status=%s exit=%d", unverifiedDim, status, exitCode)
			}
			if report.Status != StatusMainUnverified || report.ExitCode != ExitCodeMainUnverified {
				t.Fatalf("UNVERIFIED dimension %s was presented as %s (exit %d), must be MAIN_UNVERIFIED", unverifiedDim, report.Status, report.ExitCode)
			}
			if !strings.Contains(report.Summary, "UNVERIFIED") {
				t.Fatalf("unverified dimension %s did not produce UNVERIFIED summary, got: %s", unverifiedDim, report.Summary)
			}
		})
	}

	t.Run("fail_takes_precedence_over_unverified", func(t *testing.T) {
		report := MainlineHealthGateReport{
			Root:       ".",
			Dimensions: make(map[string]HealthDimensionResult, 5),
		}
		for _, d := range dims {
			switch d {
			case "git":
				report.Dimensions[d] = HealthDimensionResult{
					Name:   d,
					Passed: false,
					Status: "UNVERIFIED",
					Errors: []string{"git unverified"},
				}
			case "security":
				report.Dimensions[d] = HealthDimensionResult{
					Name:   d,
					Passed: false,
					Status: "FAIL",
					Errors: []string{"security failed"},
				}
			default:
				report.Dimensions[d] = HealthDimensionResult{
					Name:    d,
					Passed:  true,
					Status:  "PASS",
					Details: []string{"ok"},
				}
			}
		}

		status, exitCode := aggregateHealthGateReport(&report)

		if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
			t.Fatalf("expected failure to take precedence over unverified, got: status=%s exit=%d", status, exitCode)
		}
		if report.Status != StatusMainUnhealthy || report.ExitCode != ExitCodeMainUnhealthy {
			t.Fatalf("expected report StatusMainUnhealthy when failure coexists with unverified, got: %s", report.Status)
		}
	})
}

func TestRunMainlineHealthGate_RejectsMissingRequiredDimension(t *testing.T) {
	requiredDimensions := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, missingDim := range requiredDimensions {
		t.Run("missing_"+missingDim, func(t *testing.T) {
			report := MainlineHealthGateReport{
				Root:       ".",
				Dimensions: make(map[string]HealthDimensionResult),
			}
			for _, d := range requiredDimensions {
				if d == missingDim {
					continue // deliberately omit this required dimension
				}
				report.Dimensions[d] = HealthDimensionResult{
					Name:    d,
					Passed:  true,
					Status:  "PASS",
					Details: []string{"ok"},
				}
			}

			status, exitCode := aggregateHealthGateReport(&report)

			if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
				t.Fatalf("missing required dimension %s was not rejected: status=%s exit=%d", missingDim, status, exitCode)
			}
			if report.Status != StatusMainUnhealthy || report.ExitCode != ExitCodeMainUnhealthy {
				t.Fatalf("missing required dimension %s report field not StatusMainUnhealthy: status=%s exit=%d", missingDim, report.Status, report.ExitCode)
			}
			if !strings.Contains(report.Summary, "FAILED") {
				t.Fatalf("expected failure summary for missing dimension %s, got: %s", missingDim, report.Summary)
			}
		})
	}
}

func TestRunMainlineHealthGate_PassesWhenAllDimensionsPass(t *testing.T) {
	dims := []string{"code", "security", "git", "traceability", "operational_quality"}
	report := MainlineHealthGateReport{
		Root:       ".",
		Dimensions: make(map[string]HealthDimensionResult, 5),
	}
	for _, d := range dims {
		report.Dimensions[d] = HealthDimensionResult{
			Name:    d,
			Passed:  true,
			Status:  "PASS",
			Details: []string{"ok"},
		}
	}

	status, exitCode := aggregateHealthGateReport(&report)

	if status != StatusMainHealthy || exitCode != ExitCodeMainHealthy {
		t.Fatalf("expected MAIN_HEALTHY (exit 0) when all dimensions pass, got status=%s exit=%d", status, exitCode)
	}
	if report.Status != StatusMainHealthy || report.ExitCode != ExitCodeMainHealthy {
		t.Fatalf("expected report StatusMainHealthy, got status=%s exit=%d", report.Status, report.ExitCode)
	}
	if !strings.Contains(report.Summary, "PASSED") {
		t.Fatalf("expected summary to indicate pass, got: %s", report.Summary)
	}
}

func TestAggregateHealthGateReport_NilSafety(t *testing.T) {
	// Must not panic on nil pointer
	status, exitCode := aggregateHealthGateReport(nil)
	if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
		t.Fatalf("expected StatusMainUnhealthy on nil pointer, got status=%s exit=%d", status, exitCode)
	}

	// Must fail closed on empty/nil Dimensions map
	emptyReport := &MainlineHealthGateReport{Root: "."}
	aggregateHealthGateReport(emptyReport)
	if emptyReport.Status != StatusMainUnhealthy || emptyReport.ExitCode != ExitCodeMainUnhealthy {
		t.Fatalf("expected StatusMainUnhealthy on nil dimensions map, got status=%s exit=%d", emptyReport.Status, emptyReport.ExitCode)
	}
}

func TestAggregateHealthGateReport_TamperedPassedFalseBlocks(t *testing.T) {
	dims := []string{"code", "security", "git", "traceability", "operational_quality"}
	report := MainlineHealthGateReport{
		Root:       ".",
		Dimensions: make(map[string]HealthDimensionResult, 5),
	}
	for _, d := range dims {
		report.Dimensions[d] = HealthDimensionResult{
			Name:    d,
			Passed:  true,
			Status:  "PASS",
			Details: []string{"ok"},
		}
	}
	// Tamper: Status is "PASS" but Passed is false
	tampered := report.Dimensions["security"]
	tampered.Passed = false
	report.Dimensions["security"] = tampered

	status, exitCode := aggregateHealthGateReport(&report)

	if status != StatusMainUnhealthy || exitCode != ExitCodeMainUnhealthy {
		t.Fatalf("expected StatusMainUnhealthy when Passed is false, got status=%s exit=%d", status, exitCode)
	}
}

func TestEvaluateGitDimension_RejectsInvalidGitRoot(t *testing.T) {
	dim, headCommit := evaluateGitDimension("non-existent-directory-for-test")
	if dim.Passed {
		t.Fatalf("expected evaluateGitDimension to fail on invalid root")
	}
	if dim.Status != "UNVERIFIED" {
		t.Fatalf("expected Status == UNVERIFIED, got %s", dim.Status)
	}
	if headCommit != "" {
		t.Fatalf("expected empty head commit on invalid git root, got %s", headCommit)
	}
}

func TestValidateCodeAudit_BoundaryChecks(t *testing.T) {
	tests := []struct {
		name    string
		doc     DoctorOutput
		cat     CheckCatalogSummary
		wantErr bool
	}{
		{
			name:    "valid nominal",
			doc:     DoctorOutput{Healthy: true, Checks: []string{"chk1"}},
			cat:     CheckCatalogSummary{Checks: 10, Mandatory: 5},
			wantErr: false,
		},
		{
			name:    "doctor unhealthy",
			doc:     DoctorOutput{Healthy: false, Warnings: []string{"doctor warning"}, Checks: []string{"chk1"}},
			cat:     CheckCatalogSummary{Checks: 10, Mandatory: 5},
			wantErr: true,
		},
		{
			name:    "zero doctor checks",
			doc:     DoctorOutput{Healthy: true, Checks: []string{}},
			cat:     CheckCatalogSummary{Checks: 10, Mandatory: 5},
			wantErr: true,
		},
		{
			name:    "zero catalog checks",
			doc:     DoctorOutput{Healthy: true, Checks: []string{"chk1"}},
			cat:     CheckCatalogSummary{Checks: 0, Mandatory: 0},
			wantErr: true,
		},
		{
			name:    "zero mandatory checks",
			doc:     DoctorOutput{Healthy: true, Checks: []string{"chk1"}},
			cat:     CheckCatalogSummary{Checks: 10, Mandatory: 0},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateCodeAudit(tc.doc, tc.cat)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
	}
}

func TestValidateSecurityAudit_BoundaryChecks(t *testing.T) {
	tests := []struct {
		name      string
		secScan   SecretScanSummary
		gates     SecurityGateAuditSummary
		goVersion string
		wantErr   bool
	}{
		{
			name:      "valid nominal",
			secScan:   SecretScanSummary{CandidatePaths: 100, FindingCount: 0},
			gates:     SecurityGateAuditSummary{Gates: 10, Blocking: 5},
			goVersion: "1.26.11",
			wantErr:   false,
		},
		{
			name:      "secrets detected",
			secScan:   SecretScanSummary{CandidatePaths: 100, FindingCount: 1},
			gates:     SecurityGateAuditSummary{Gates: 10, Blocking: 5},
			goVersion: "1.26.11",
			wantErr:   true,
		},
		{
			name:      "zero candidate paths",
			secScan:   SecretScanSummary{CandidatePaths: 0, FindingCount: 0},
			gates:     SecurityGateAuditSummary{Gates: 10, Blocking: 5},
			goVersion: "1.26.11",
			wantErr:   true,
		},
		{
			name:      "zero registered gates",
			secScan:   SecretScanSummary{CandidatePaths: 100, FindingCount: 0},
			gates:     SecurityGateAuditSummary{Gates: 0, Blocking: 0},
			goVersion: "1.26.11",
			wantErr:   true,
		},
		{
			name:      "zero blocking gates",
			secScan:   SecretScanSummary{CandidatePaths: 100, FindingCount: 0},
			gates:     SecurityGateAuditSummary{Gates: 10, Blocking: 0},
			goVersion: "1.26.11",
			wantErr:   true,
		},
		{
			name:      "go toolchain version mismatch",
			secScan:   SecretScanSummary{CandidatePaths: 100, FindingCount: 0},
			gates:     SecurityGateAuditSummary{Gates: 10, Blocking: 5},
			goVersion: "1.25.0",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateSecurityAudit(tc.secScan, tc.gates, tc.goVersion)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
	}
}

func TestValidateTraceabilityAudit_BoundaryChecks(t *testing.T) {
	tests := []struct {
		name    string
		ledger  ImplementationLedgerAuditSummary
		dec     DecisionRegistryAuditSummary
		wantErr bool
	}{
		{
			name:    "valid nominal",
			ledger:  ImplementationLedgerAuditSummary{Status: "valid", Entries: 10, VerifiedClaims: 10},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: false,
		},
		{
			name:    "valid PASS status",
			ledger:  ImplementationLedgerAuditSummary{Status: "PASS", Entries: 10, VerifiedClaims: 10},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: false,
		},
		{
			name:    "invalid ledger status",
			ledger:  ImplementationLedgerAuditSummary{Status: "invalid", Entries: 10, VerifiedClaims: 10},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: true,
		},
		{
			name:    "zero ledger entries",
			ledger:  ImplementationLedgerAuditSummary{Status: "valid", Entries: 0, VerifiedClaims: 0},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: true,
		},
		{
			name:    "zero verified claims",
			ledger:  ImplementationLedgerAuditSummary{Status: "valid", Entries: 10, VerifiedClaims: 0},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: true,
		},
		{
			name:    "claims entries mismatch",
			ledger:  ImplementationLedgerAuditSummary{Status: "valid", Entries: 10, VerifiedClaims: 9},
			dec:     DecisionRegistryAuditSummary{Decisions: 5},
			wantErr: true,
		},
		{
			name:    "zero decisions",
			ledger:  ImplementationLedgerAuditSummary{Status: "valid", Entries: 10, VerifiedClaims: 10},
			dec:     DecisionRegistryAuditSummary{Decisions: 0},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateTraceabilityAudit(tc.ledger, tc.dec)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
	}
}

func TestValidateOperationalQualityAudit_BoundaryChecks(t *testing.T) {
	tests := []struct {
		name    string
		comp    ComplexityAuditSummary
		perf    PerformanceRegistryAuditSummary
		rec     RecoveryJournalAuditSummary
		wantErr bool
	}{
		{
			name:    "valid nominal with quiescent recovery journal",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "empty", Events: 0},
			wantErr: false,
		},
		{
			name:    "valid nominal with active valid recovery journal",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "valid", Events: 3},
			wantErr: false,
		},
		{
			name:    "complexity readiness violations",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 2},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "empty"},
			wantErr: true,
		},
		{
			name:    "zero complexity evaluated units",
			comp:    ComplexityAuditSummary{Evaluated: 0, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "empty"},
			wantErr: true,
		},
		{
			name:    "incomplete performance runtime coverage",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "partial", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "empty"},
			wantErr: true,
		},
		{
			name:    "zero performance targets",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 0, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "empty"},
			wantErr: true,
		},
		{
			name:    "zero performance metrics",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 0},
			rec:     RecoveryJournalAuditSummary{Status: "empty"},
			wantErr: true,
		},
		{
			name:    "invalid recovery journal status",
			comp:    ComplexityAuditSummary{Evaluated: 10, ReadinessViolations: 0},
			perf:    PerformanceRegistryAuditSummary{RuntimeCoverageStatus: "complete", Targets: 5, Metrics: 10},
			rec:     RecoveryJournalAuditSummary{Status: "invalid"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateOperationalQualityAudit(tc.comp, tc.perf, tc.rec)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
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
			"git": {
				Name:    "git",
				Passed:  true,
				Status:  "PASS",
				Details: []string{"clean tree"},
			},
			"traceability": {
				Name:    "traceability",
				Passed:  true,
				Status:  "PASS",
				Details: []string{"ledger valid"},
			},
			"operational_quality": {
				Name:    "operational_quality",
				Passed:  true,
				Status:  "PASS",
				Details: []string{"recovery journal quiescent"},
			},
		},
	}
	// Call printMainlineHealthGateReport to ensure formatting does not panic
	printMainlineHealthGateReport(report)
}

func writeJSONFixture(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("marshal %s: %v", path, err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeTextFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestValidateDependencyAudit_BoundaryChecks(t *testing.T) {
	tests := []struct {
		name    string
		summary DependencyAuditSummary
		wantErr bool
	}{
		{
			name:    "valid nominal",
			summary: DependencyAuditSummary{RuntimeBoundaries: 13, Unapproved: 0, UnsupportedManifests: 0},
			wantErr: false,
		},
		{
			name:    "unapproved third-party dependency",
			summary: DependencyAuditSummary{RuntimeBoundaries: 13, Unapproved: 1, UnsupportedManifests: 0},
			wantErr: true,
		},
		{
			name:    "unsupported manifest",
			summary: DependencyAuditSummary{RuntimeBoundaries: 13, Unapproved: 0, UnsupportedManifests: 1},
			wantErr: true,
		},
		{
			name:    "zero runtime boundaries",
			summary: DependencyAuditSummary{RuntimeBoundaries: 0, Unapproved: 0, UnsupportedManifests: 0},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := validateDependencyAudit(tc.summary)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
	}
}

func TestValidateCommittedRuntimeSecurityEvidence_BoundaryChecks(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"product-runtime/api/**"}}
	arch := securityTestArchitecture(runtime)
	gates := securityTestGates(true, true)
	scope := RuntimeSecurityScope{
		BoundaryID:               "api",
		Owner:                    "security-platform",
		SecurityCriticalPaths:    []string{"product-runtime/api/auth/**"},
		AuthorizationRequired:    true,
		AuthorizationRationale:   "API enforces authorization",
		TenantIsolationRequired:  true,
		TenantIsolationRationale: "API enforces tenant isolation",
	}
	policy := baseRuntimePolicy(scope)
	states := map[string]string{"api": "implemented"}
	historicalCommit := "e4ea9a315cfa655520d1bc3e308ed9c45c326041"

	validEvidence := func() RuntimeSecurityEvidence {
		return RuntimeSecurityEvidence{
			SchemaVersion:         1,
			Status:                "measured",
			Reason:                "CI-generated measured runtime security evidence",
			SourceCommit:          historicalCommit,
			GlobalCoveragePercent: ptrFloat(91.5),
			Scopes: []RuntimeSecurityEvidenceScope{{
				BoundaryID:                    "api",
				ImplementationState:           "implemented",
				ChangedSecurityCritical:       false,
				AuthorizationNegativePassed:   ptrBool(true),
				TenantIsolationNegativePassed: ptrBool(true),
			}},
		}
	}

	if summary, err := validateCommittedRuntimeSecurityEvidence(arch, gates, policy, validEvidence(), states); err != nil {
		t.Fatalf("expected valid historical commit evidence to pass, got: %v", err)
	} else if summary.CoverageStatus != "measured" || summary.RegisteredScopes != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	cases := []struct {
		name   string
		mutate func(*RuntimeSecurityEvidence)
	}{
		{
			name:   "wrong schema version",
			mutate: func(e *RuntimeSecurityEvidence) { e.SchemaVersion = 2 },
		},
		{
			name:   "empty reason",
			mutate: func(e *RuntimeSecurityEvidence) { e.Reason = "   " },
		},
		{
			name:   "invalid source commit",
			mutate: func(e *RuntimeSecurityEvidence) { e.SourceCommit = "not-a-sha" },
		},
		{
			name:   "global coverage below floor",
			mutate: func(e *RuntimeSecurityEvidence) { e.GlobalCoveragePercent = ptrFloat(79.5) },
		},
		{
			name:   "missing scope",
			mutate: func(e *RuntimeSecurityEvidence) { e.Scopes = nil },
		},
		{
			name: "extra unregistered scope",
			mutate: func(e *RuntimeSecurityEvidence) {
				e.Scopes = append(e.Scopes, RuntimeSecurityEvidenceScope{BoundaryID: "extra", ImplementationState: "implemented"})
			},
		},
		{
			name: "duplicate scope",
			mutate: func(e *RuntimeSecurityEvidence) {
				e.Scopes = append(e.Scopes, e.Scopes[0])
			},
		},
		{
			name: "mismatched implementation state",
			mutate: func(e *RuntimeSecurityEvidence) {
				e.Scopes[0].ImplementationState = "preimplementation"
			},
		},
		{
			name: "authorization negative test false",
			mutate: func(e *RuntimeSecurityEvidence) {
				e.Scopes[0].AuthorizationNegativePassed = ptrBool(false)
			},
		},
		{
			name: "tenant isolation negative test nil",
			mutate: func(e *RuntimeSecurityEvidence) {
				e.Scopes[0].TenantIsolationNegativePassed = nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ev := validEvidence()
			tc.mutate(&ev)
			if _, err := validateCommittedRuntimeSecurityEvidence(arch, gates, policy, ev, states); err == nil {
				t.Fatalf("expected error for %s, got nil", tc.name)
			}
		})
	}
}

func seedSyntheticSecurityRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	writeTextFixture(t, filepath.Join(root, "cmdr-product-spec", "README.md"), "# Spec\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "go.mod"), "module example.com/mod-a\n\ngo 1.26.11\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "mod.go"), "package moda\n\nfunc Value() int { return 1 }\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-b", "go.mod"), "module example.com/mod-b\n\ngo 1.26.11\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-b", "mod.go"), "package modb\n\nfunc Value() int { return 2 }\n")

	arch := ArchitectureRegistry{
		SchemaVersion: 1,
		DefaultPolicy: "deny-unregistered-path",
		Boundaries: []ArchitectureBoundary{
			{ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"}},
			{ID: "mod-a-runtime", Kind: "product-runtime", Roots: []string{"product-runtime/mod-a/**"}},
			{ID: "mod-b-runtime", Kind: "product-module", Roots: []string{"product-runtime/mod-b/**"}},
		},
	}
	writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), arch)

	writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(trustedBasePath)), TrustedBase{
		SchemaVersion: 1,
		DefaultPolicy: "deny-third-party-product-functionality",
	})

	writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(secretAllowlistPath)), SecretAllowlist{
		SchemaVersion: 1,
		Entries:       []SecretAllowlistEntry{},
	})

	realGatesBytes, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(securityGateRegistryPath)))
	if err != nil {
		t.Fatalf("read canonical security-gates.json: %v", err)
	}
	writeTextFixture(t, filepath.Join(root, filepath.FromSlash(securityGateRegistryPath)), string(realGatesBytes))

	writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(runtimeSecurityPolicyPath)), RuntimeSecurityPolicy{
		SchemaVersion:                        1,
		DefaultPolicy:                        "deny-unregistered-runtime-boundary",
		GlobalCoverageFloorPercent:           80,
		ChangedSecurityCriticalCoverageFloor: 90,
		AuthorizationGateID:                  "SEC-AUTH-NEG-001",
		TenantIsolationGateID:                "SEC-TENANT-ISO-001",
		Scopes: []RuntimeSecurityScope{{
			BoundaryID:               "mod-a-runtime",
			Owner:                    "security-platform",
			SecurityCriticalPaths:    []string{"product-runtime/mod-a/**"},
			AuthorizationRequired:    true,
			AuthorizationRationale:   "mod-a enforces auth",
			TenantIsolationRequired:  true,
			TenantIsolationRationale: "mod-a enforces tenant",
		}},
	})

	writeJSONFixture(t, filepath.Join(root, filepath.FromSlash(runtimeSecurityEvidencePath)), RuntimeSecurityEvidence{
		SchemaVersion:         1,
		Status:                "measured",
		Reason:                "measured runtime security evidence",
		SourceCommit:          "e4ea9a315cfa655520d1bc3e308ed9c45c326041",
		GlobalCoveragePercent: ptrFloat(92.0),
		Scopes: []RuntimeSecurityEvidenceScope{{
			BoundaryID:                    "mod-a-runtime",
			ImplementationState:           "implemented",
			ChangedSecurityCritical:       false,
			AuthorizationNegativePassed:   ptrBool(true),
			TenantIsolationNegativePassed: ptrBool(true),
		}},
	})

	return root
}

func TestRunCommittedRuntimeSecurityEvidenceAudit_MissingOrCorrupt(t *testing.T) {
	root := seedSyntheticSecurityRepo(t)
	if summary, err := runCommittedRuntimeSecurityEvidenceAudit(root); err != nil {
		t.Fatalf("expected clean synthetic repo to pass committed evidence audit, got: %v", err)
	} else if summary.CoverageStatus != "measured" || summary.RegisteredScopes != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	evPath := filepath.Join(root, filepath.FromSlash(runtimeSecurityEvidencePath))
	if err := os.WriteFile(evPath, []byte("{corrupt json"), 0600); err != nil {
		t.Fatalf("corrupt evidence file: %v", err)
	}
	if _, err := runCommittedRuntimeSecurityEvidenceAudit(root); err == nil {
		t.Fatalf("expected corrupt runtime-security-evidence.json to fail")
	}

	if err := os.Remove(evPath); err != nil {
		t.Fatalf("remove evidence file: %v", err)
	}
	if _, err := runCommittedRuntimeSecurityEvidenceAudit(root); err == nil {
		t.Fatalf("expected missing runtime-security-evidence.json to fail")
	}
}

func TestEvaluateSecurityDimension_RejectsExternalDepAndBadSecurityEvidence(t *testing.T) {
	root := seedSyntheticSecurityRepo(t)
	for _, args := range [][]string{
		{"init"},
		{"config", "user.name", "Test"},
		{"config", "user.email", "test@example.com"},
		{"add", "."},
		{"commit", "-m", "baseline"},
	} {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...) // #nosec G204 -- fixed git command on test TempDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v (%s)", args, err, string(out))
		}
	}

	// Update evidence source_commit to the synthetic repo HEAD commit so gitCommitExists succeeds too.
	headBytes, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output() // #nosec G204 -- fixed git command on test TempDir
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v", err)
	}
	headCommit := strings.TrimSpace(string(headBytes))
	evPath := filepath.Join(root, filepath.FromSlash(runtimeSecurityEvidencePath))
	validEv := RuntimeSecurityEvidence{
		SchemaVersion:         1,
		Status:                "measured",
		Reason:                "measured runtime security evidence",
		SourceCommit:          headCommit,
		GlobalCoveragePercent: ptrFloat(92.0),
		Scopes: []RuntimeSecurityEvidenceScope{{
			BoundaryID:                    "mod-a-runtime",
			ImplementationState:           "implemented",
			ChangedSecurityCritical:       false,
			AuthorizationNegativePassed:   ptrBool(true),
			TenantIsolationNegativePassed: ptrBool(true),
		}},
	}
	writeJSONFixture(t, evPath, validEv)

	clean := evaluateSecurityDimension(root)
	if !clean.Passed || clean.Status != "PASS" {
		t.Fatalf("expected clean security dimension PASS, got status=%s errors=%v", clean.Status, clean.Errors)
	}

	// Inject external dependency into product-module boundary (mod-b).
	modBGoMod := filepath.Join(root, "product-runtime", "mod-b", "go.mod")
	writeTextFixture(t, modBGoMod, "module example.com/mod-b\n\ngo 1.26.11\n\nrequire github.com/example/evil v1.0.0\n")
	depFail := evaluateSecurityDimension(root)
	if depFail.Passed || depFail.Status != "FAIL" {
		t.Fatalf("expected FAIL on external dep in product-module, got status=%s", depFail.Status)
	}

	// Restore mod-b/go.mod and corrupt runtime-security-evidence.json.
	writeTextFixture(t, modBGoMod, "module example.com/mod-b\n\ngo 1.26.11\n")
	if err := os.WriteFile(evPath, []byte("{corrupt json"), 0600); err != nil {
		t.Fatalf("corrupt evidence file: %v", err)
	}
	corruptFail := evaluateSecurityDimension(root)
	if corruptFail.Passed || corruptFail.Status != "FAIL" {
		t.Fatalf("expected FAIL on corrupt runtime-security-evidence.json, got status=%s", corruptFail.Status)
	}

	// Remove runtime-security-evidence.json.
	if err := os.Remove(evPath); err != nil {
		t.Fatalf("remove evidence file: %v", err)
	}
	missingFail := evaluateSecurityDimension(root)
	if missingFail.Passed || missingFail.Status != "FAIL" {
		t.Fatalf("expected FAIL on missing runtime-security-evidence.json, got status=%s", missingFail.Status)
	}
}

func TestDiscoverCodeTestModules_AndRunRepositoryGoUnitTests(t *testing.T) {
	root := t.TempDir()
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "go.mod"), "module example.com/mod-a\n\ngo 1.26.11\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "mod_test.go"), "package moda\n\nimport \"testing\"\n\nfunc TestPass(t *testing.T) {}\n")
	writeTextFixture(t, filepath.Join(root, "product-runtime", "frontend", "index.js"), "export const ok = true;\n")
	writeTextFixture(t, filepath.Join(root, "tools", "cmdr-dev", "go.mod"), "module example.com/cmdr-dev\n\ngo 1.26.11\n")
	writeTextFixture(t, filepath.Join(root, "tools", "cmdr-dev", "dev_test.go"), "package main\n\nimport \"testing\"\n\nfunc TestDevPass(t *testing.T) {}\n")
	writeTextFixture(t, filepath.Join(root, "tests", "e2e", "go.mod"), "module example.com/e2e\n\ngo 1.26.11\n")
	writeTextFixture(t, filepath.Join(root, "tests", "e2e", "e2e_test.go"), "package e2e\n\nimport \"testing\"\n\nfunc TestShouldNotRun(t *testing.T) { t.Fatal(\"must be excluded\") }\n")

	mods, err := discoverCodeTestModules(root)
	if err != nil {
		t.Fatalf("discoverCodeTestModules failed: %v", err)
	}
	if len(mods) != 2 || mods[0] != "tools/cmdr-dev" || mods[1] != "product-runtime/mod-a" {
		t.Fatalf("unexpected discovered modules: %v", mods)
	}

	tested, unverifiedErr, err := runRepositoryGoUnitTests(root)
	if unverifiedErr != nil || err != nil {
		t.Fatalf("expected runRepositoryGoUnitTests pass, got unverified=%v err=%v", unverifiedErr, err)
	}
	if len(tested) != 2 {
		t.Fatalf("expected 2 tested modules, got %d", len(tested))
	}

	writeTextFixture(t, filepath.Join(root, "product-runtime", "mod-a", "fail_test.go"), "package moda\n\nimport \"testing\"\n\nfunc TestIntentionalFail(t *testing.T) { t.Fatal(\"boom\") }\n")
	_, unverifiedErr, err = runRepositoryGoUnitTests(root)
	if unverifiedErr != nil {
		t.Fatalf("unexpected unverifiedErr on test failure: %v", unverifiedErr)
	}
	if err == nil || !strings.Contains(err.Error(), "product-runtime/mod-a") {
		t.Fatalf("expected failure mentioning product-runtime/mod-a, got: %v", err)
	}
}
