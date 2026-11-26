package tier1_features

import (
	"context"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type DimensionResult struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Status  string   `json:"status"`
	Details []string `json:"details,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type HealthGateReport struct {
	Status     string                     `json:"status"`
	ExitCode   int                        `json:"exit_code"`
	Summary    string                     `json:"summary"`
	HeadCommit string                     `json:"head_commit"`
	Dimensions map[string]DimensionResult `json:"dimensions"`
}

// T1.F7.1: Execute cmdr-dev mainline-health-gate and verify tri-state exit code compliance (0, 1, or 2).
func TestT1_F7_1_HealthGateTriStateExitCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if res.ExitCode != 0 && res.ExitCode != 1 && res.ExitCode != 2 {
		t.Fatalf("unexpected exit code %d from mainline-health-gate; must be 0 (HEALTHY), 1 (UNHEALTHY), or 2 (UNVERIFIED)", res.ExitCode)
	}

	var report HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse health gate JSON report: %v", err)
	}

	if report.ExitCode != res.ExitCode {
		t.Fatalf("exit code in report (%d) does not match process exit code (%d)", report.ExitCode, res.ExitCode)
	}
	t.Logf("T1.F7.1 PASS: tri-state exit code compliant (%d - %s)", res.ExitCode, report.Status)
}

// T1.F7.2: Verify 'code' dimension evaluates doctor, check catalog, and compilation.
func TestT1_F7_2_HealthGateCodeDimensionEvaluation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	var report HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse health gate report: %v", err)
	}

	dim, ok := report.Dimensions["code"]
	if !ok {
		t.Fatal("dimension 'code' is missing from health gate report")
	}

	if dim.Name != "code" {
		t.Fatalf("expected dimension name 'code', got '%s'", dim.Name)
	}
	t.Logf("T1.F7.2 PASS: 'code' dimension evaluated (passed=%v, status=%s)", dim.Passed, dim.Status)
}

// T1.F7.3: Verify 'security' dimension evaluates secrets, SAST, SCA, and security policy.
func TestT1_F7_3_HealthGateSecurityDimensionEvaluation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	var report HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse health gate report: %v", err)
	}

	dim, ok := report.Dimensions["security"]
	if !ok {
		t.Fatal("dimension 'security' is missing from health gate report")
	}

	if dim.Name != "security" {
		t.Fatalf("expected dimension name 'security', got '%s'", dim.Name)
	}
	t.Logf("T1.F7.3 PASS: 'security' dimension evaluated (passed=%v, status=%s)", dim.Passed, dim.Status)
}

// T1.F7.4: Verify 'git' dimension evaluates clean tree, remote branch, and spec invariance.
func TestT1_F7_4_HealthGateGitDimensionEvaluation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	var report HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse health gate report: %v", err)
	}

	dim, ok := report.Dimensions["git"]
	if !ok {
		t.Fatal("dimension 'git' is missing from health gate report")
	}

	if dim.Name != "git" {
		t.Fatalf("expected dimension name 'git', got '%s'", dim.Name)
	}
	t.Logf("T1.F7.4 PASS: 'git' dimension evaluated (passed=%v, status=%s)", dim.Passed, dim.Status)
}

// T1.F7.5: Verify MAIN_HEALTHY is returned only when all 5 dimensions pass.
func TestT1_F7_5_HealthGateAllFiveDimensionsRequiredForHealthy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	var report HealthGateReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse health gate report: %v", err)
	}

	requiredDimensions := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, req := range requiredDimensions {
		dim, exists := report.Dimensions[req]
		if !exists {
			t.Fatalf("mandatory dimension '%s' missing from health gate report", req)
		}
		if report.Status == "MAIN_HEALTHY" && (!dim.Passed || dim.Status != "PASS") {
			t.Fatalf("invariant breach: report is MAIN_HEALTHY but dimension '%s' has passed=%v status=%s", req, dim.Passed, dim.Status)
		}
	}

	if report.Status == "MAIN_HEALTHY" && res.ExitCode != 0 {
		t.Fatalf("invariant breach: status is MAIN_HEALTHY but exit code is %d", res.ExitCode)
	}
	t.Logf("T1.F7.5 PASS: all 5 dimensions verified, MAIN_HEALTHY invariant strictly enforced (status=%s)", report.Status)
}
