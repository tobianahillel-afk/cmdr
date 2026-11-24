package tier1_features

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type SASTReport struct {
	Findings   int            `json:"findings"`
	BySeverity map[string]int `json:"by_severity"`
	FindingItems []struct {
		RuleID string `json:"rule_id"`
		Path   string `json:"path"`
		Line   int    `json:"line"`
	} `json:"finding_items"`
}

// T1.F2.1: Execute cmdr-dev sast-go and verify exit code 0.
func TestT1_F2_1_SASTExecutionSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go exited with code %d, expected 0. Stderr: %s", res.ExitCode, res.Stderr)
	}
	t.Log("T1.F2.1 PASS: cmdr-dev sast-go exited cleanly with code 0")
}

// T1.F2.2: Verify findings count is 0.
func TestT1_F2_2_SASTTotalFindingsZero(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go failed: %s", res.Stderr)
	}

	var report SASTReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse sast-go output: %v\nOutput: %s", err, res.Stdout)
	}

	if report.Findings != 0 {
		t.Fatalf("expected 0 SAST findings, got %d", report.Findings)
	}
	t.Log("T1.F2.2 PASS: total SAST findings count is exactly 0")
}

// T1.F2.3: Verify by_severity.HIGH count is 0.
func TestT1_F2_3_SASTHighSeverityZero(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go failed: %s", res.Stderr)
	}

	harness.AssertJSONFieldEqual(t, res.Stdout, "by_severity.HIGH", 0)
	t.Log("T1.F2.3 PASS: SAST by_severity.HIGH is exactly 0")
}

// T1.F2.4: Verify by_severity.MEDIUM count is 0.
func TestT1_F2_4_SASTMediumSeverityZero(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go failed: %s", res.Stderr)
	}

	harness.AssertJSONFieldEqual(t, res.Stdout, "by_severity.MEDIUM", 0)
	t.Log("T1.F2.4 PASS: SAST by_severity.MEDIUM is exactly 0")
}

// T1.F2.5: Verify absence of G204/G702 finding items on performance_timeline_process_windows.go.
func TestT1_F2_5_SASTWindowsTaskkillAbsence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go failed: %s", res.Stderr)
	}

	var report SASTReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse sast-go output: %v", err)
	}

	for _, item := range report.FindingItems {
		if strings.Contains(item.Path, "performance_timeline_process_windows.go") {
			if item.RuleID == "G204" || item.RuleID == "G702" {
				t.Fatalf("unexpected finding %s on %s:%d", item.RuleID, item.Path, item.Line)
			}
		}
	}
	t.Log("T1.F2.5 PASS: zero G204/G702 findings on performance_timeline_process_windows.go")
}
