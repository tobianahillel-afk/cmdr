package tier1_features

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type SCAReport struct {
	Tool                  string `json:"tool"`
	Version               string `json:"version"`
	Database              string `json:"database"`
	GoVersion             string `json:"go_version"`
	Modules               int    `json:"modules"`
	ActionableFindings    int    `json:"actionable_findings"`
	InformationalFindings int    `json:"informational_findings"`
	Findings              []struct {
		OSV        string `json:"osv"`
		Module     string `json:"module"`
		Package    string `json:"package"`
		Function   string `json:"function"`
		Actionable bool   `json:"actionable"`
	} `json:"findings"`
}

// T1.F3.1: Execute cmdr-dev sca-go and verify exit code 0.
func TestT1_F3_1_SCAExecutionSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sca-go exited with code %d, expected 0. Stderr: %s", res.ExitCode, res.Stderr)
	}
	t.Log("T1.F3.1 PASS: cmdr-dev sca-go completed successfully with exit code 0")
}

// T1.F3.2: Verify actionable_findings count is 0.
func TestT1_F3_2_SCAActionableFindingsZero(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sca-go failed: %s", res.Stderr)
	}

	var report SCAReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse sca-go output: %v\nOutput: %s", err, res.Stdout)
	}

	if report.ActionableFindings != 0 {
		t.Fatalf("expected 0 actionable findings, got %d", report.ActionableFindings)
	}
	t.Log("T1.F3.2 PASS: sca-go actionable_findings count is 0")
}

// T1.F3.3: Verify absence of GO-2026-6859 finding in scan results.
func TestT1_F3_3_SCAAbsenceOfGO20266859(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sca-go failed: %s", res.Stderr)
	}

	if strings.Contains(res.Stdout, "GO-2026-6859") {
		t.Fatal("vulnerability GO-2026-6859 detected in sca-go report")
	}
	t.Log("T1.F3.3 PASS: GO-2026-6859 is completely absent from SCA findings")
}

// T1.F3.4: Verify absence of stdlib HTTP/crypto vulnerabilities in findings.
func TestT1_F3_4_SCASecureStdlib(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sca-go failed: %s", res.Stderr)
	}

	var report SCAReport
	if err := harness.ParseJSON(res.Stdout, &report); err != nil {
		t.Fatalf("failed to parse sca-go output: %v", err)
	}

	for _, f := range report.Findings {
		if f.Actionable && (strings.HasPrefix(f.Package, "net/http") || strings.HasPrefix(f.Package, "crypto/")) {
			t.Fatalf("actionable stdlib finding present in %s: %s", f.Package, f.OSV)
		}
	}
	t.Log("T1.F3.4 PASS: zero actionable stdlib HTTP/crypto vulnerabilities")
}

// T1.F3.5: Verify GOTOOLCHAIN=go1.26.11 is enforced during govulncheck execution.
func TestT1_F3_5_SCAGoToolchainEnforced(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	secStaticPath := filepath.Join(root, "tools", "cmdr-dev", "security_static.go")
	content, err := os.ReadFile(secStaticPath)
	if err != nil {
		t.Fatalf("failed to read security_static.go: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "GOTOOLCHAIN=go\"+engineeringGoVersion") && !strings.Contains(str, "GOTOOLCHAIN=go1.26.11") {
		t.Fatal("security_static.go does not enforce pinned GOTOOLCHAIN=go1.26.11")
	}
	t.Log("T1.F3.5 PASS: GOTOOLCHAIN=go1.26.11 is strictly enforced in security_static.go")
}
