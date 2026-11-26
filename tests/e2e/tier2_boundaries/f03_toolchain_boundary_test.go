package tier2_boundaries

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type SCAFindingItem struct {
	OSV          string `json:"osv"`
	Ecosystem    string `json:"ecosystem"`
	Module       string `json:"module"`
	Version      string `json:"version"`
	Package      string `json:"package,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	FixedVersion string `json:"fixed_version,omitempty"`
}

type SCASummaryReport struct {
	Tool                  string           `json:"tool"`
	Version               string           `json:"version"`
	Database              string           `json:"database,omitempty"`
	GoVersion             string           `json:"go_version,omitempty"`
	Modules               int              `json:"modules"`
	InformationalFindings int              `json:"informational_findings"`
	ActionableFindings    int              `json:"actionable_findings"`
	Findings              []SCAFindingItem `json:"finding_items,omitempty"`
}

// T2.F3.1: SCA fails closed when local unpatched toolchain is forced.
func TestT2_F3_1_SCAFailsClosedOnUnpatchedToolchain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "sca-unpatched")
	goModContent := "module example.com/unpatched\n\ngo 1.26.10\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	mainCode := `package main

import "net/http"

func main() {
	_ = http.DefaultClient
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainCode), 0644); err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	// Force local unpatched toolchain explicitly
	env := []string{
		"GOTOOLCHAIN=go1.26.10",
	}

	res := harness.RunCmd(ctx, tmpDir, env, "go", "run", "golang.org/x/vuln/cmd/govulncheck@v1.8.0", "-json", "./...")
	// If unpatched Go 1.26.10 is forced, govulncheck either detects vulnerabilities (actionable findings)
	// or rejects execution if go1.26.10 cannot satisfy patched security requirements
	t.Logf("T2.F3.1: govulncheck with GOTOOLCHAIN=go1.26.10 finished with exit code %d", res.ExitCode)

	// Verify fail-closed behavior: either non-zero exit or findings detected in output
	failedClosed := res.ExitCode != 0 || strings.Contains(res.Stdout, "GO-2026-") || strings.Contains(res.Stdout, "VULN")
	if !failedClosed {
		t.Fatalf("T2.F3.1 FAIL: SCA did not fail closed when forced onto unpatched toolchain")
	}
	t.Log("T2.F3.1 PASS: SCA fails closed when unpatched toolchain is forced")
}

// T2.F3.2: SCA handles network disconnect gracefully during database lookup.
func TestT2_F3_2_SCAHandlesNetworkDisconnectGracefully(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "sca-offline")
	goModContent := "module example.com/offlinemod\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644); err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	// Disable network and point to non-routable vulnerability DB
	env := []string{
		"GOPROXY=off",
		"GOPRIVATE=*",
		"GOVULNDB=http://127.0.0.1:9/nonexistent-vulndb",
	}

	res := harness.RunCmd(ctx, tmpDir, env, "go", "run", "golang.org/x/vuln/cmd/govulncheck@v1.8.0", "-json", "./...")
	// Must fail gracefully with clean error message, not panic
	if strings.Contains(res.Stderr, "panic:") || strings.Contains(res.Stdout, "panic:") {
		t.Fatalf("T2.F3.2 FAIL: govulncheck panicked during network disconnection: %s", res.Stderr)
	}
	t.Logf("T2.F3.2 PASS: network disconnection handled gracefully (exit %d)", res.ExitCode)
}

// T2.F3.3: SCA correctly detects vulnerable dependency and rejects tampered tool provenance.
func TestT2_F3_3_SCADetectsVulnerableDependency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Run real cmdr-dev govulncheck stream-parser and fail-closed unit tests
	expectedTests := []string{
		"TestParseGovulncheckReachableFindingHasExactIdentity",
		"TestGoSCAReportsOnlyValidatedPublicAdvisoryIDsAndStillBlocks",
		"TestParseGovulncheckUsesSBOMVersionAndDeduplicates",
		"TestParseGovulncheckFailsClosedWithoutExactVersion",
	}
	res, err := harness.RunGoTestVerified(
		ctx,
		"tools/cmdr-dev",
		"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev",
		expectedTests,
		"-run", "^("+strings.Join(expectedTests, "|")+")$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F3.3 FAIL: cmdr-dev govulncheck parser tests failed: %v\nOutput: %s", err, res.Output)
	}

	// 2. Verify cmdr-dev sca-go fails closed on tampered engineering/security/development-tools.json
	mockRoot := harness.CloneIsolatedRepo(t, "f03-tamper-tools")
	toolsPath := filepath.Join(mockRoot, "engineering", "security", "development-tools.json")
	raw, err := os.ReadFile(toolsPath)
	if err != nil {
		t.Fatalf("failed to read development-tools.json: %v", err)
	}
	var reg map[string]interface{}
	if err := json.Unmarshal(raw, &reg); err != nil {
		t.Fatalf("failed to parse development-tools.json: %v", err)
	}
	tools, ok := reg["tools"].([]interface{})
	if !ok || len(tools) == 0 {
		t.Fatalf("development-tools.json missing tools array")
	}
	for _, entry := range tools {
		if m, ok := entry.(map[string]interface{}); ok && m["id"] == "govulncheck" {
			m["version"] = "v0.0.0-tampered"
		}
	}
	mutated, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal tampered development-tools.json: %v", err)
	}
	if err := os.WriteFile(toolsPath, append(mutated, '\n'), 0644); err != nil {
		t.Fatalf("failed to write tampered development-tools.json: %v", err)
	}

	scaRes := harness.RunCmdrDev(ctx, "sca-go", "--root", mockRoot, "--json")
	if scaRes.ExitCode == 0 {
		t.Fatalf("T2.F3.3 FAIL: cmdr-dev sca-go succeeded despite tampered govulncheck provenance")
	}
	combined := scaRes.Stderr + " " + scaRes.Stdout
	if !strings.Contains(combined, "provenance mismatch") {
		t.Fatalf("T2.F3.3 FAIL: expected provenance mismatch error, got: %s", combined)
	}
	t.Log("T2.F3.3 PASS: real govulncheck vulnerability parser tests and tampered provenance rejection verified")
}

// T2.F3.4: SCA report format maintains schema when 0 actionable findings.
func TestT2_F3_4_SCAReportMaintainsSchemaOnZeroFindings(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("T2.F3.4 FAIL: cmdr-dev sca-go exited %d: %s\n%s", res.ExitCode, res.Stderr, res.Stdout)
	}

	var parsed SCASummaryReport
	if err := harness.StrictUnmarshal([]byte(strings.TrimSpace(res.Stdout)), &parsed); err != nil {
		t.Fatalf("T2.F3.4 FAIL: strict unmarshal of live sca-go output failed: %v\nStdout: %s", err, res.Stdout)
	}
	if parsed.Tool != "govulncheck" {
		t.Fatalf("T2.F3.4 FAIL: expected tool=govulncheck, got %q", parsed.Tool)
	}
	if parsed.Version != "v1.8.0" {
		t.Fatalf("T2.F3.4 FAIL: expected version=v1.8.0, got %q", parsed.Version)
	}
	if parsed.GoVersion != "go1.26.11" {
		t.Fatalf("T2.F3.4 FAIL: expected go_version=go1.26.11, got %q", parsed.GoVersion)
	}
	if parsed.Modules <= 0 {
		t.Fatalf("T2.F3.4 FAIL: expected modules > 0, got %d", parsed.Modules)
	}
	if parsed.ActionableFindings != 0 || len(parsed.Findings) != 0 {
		t.Fatalf("T2.F3.4 FAIL: expected 0 actionable findings and 0 finding_items, got actionable=%d items=%d",
			parsed.ActionableFindings, len(parsed.Findings))
	}
	t.Logf("T2.F3.4 PASS: live cmdr-dev sca-go --json strictly conforms to SCASummary schema (modules=%d, go_version=%s)",
		parsed.Modules, parsed.GoVersion)
}

// T2.F3.5: SCA rejects invalid toolchain version string format.
func TestT2_F3_5_SCARejectsInvalidToolchainVersionFormat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Real `go version` CLI rejects malformed GOTOOLCHAIN values under GOPROXY=off
	invalidToolchains := []string{
		"1.26.11",
		"v1.26.11",
		"go1.26",
		"go1.26.11;rm",
		"go1.26.11-rc1",
		"go1.26.11.1",
	}
	for _, tc := range invalidToolchains {
		res := harness.RunCmd(ctx, "", []string{"GOPROXY=off", "GOTOOLCHAIN=" + tc}, "go", "version")
		if res.ExitCode == 0 {
			t.Fatalf("T2.F3.5 FAIL: `go version` accepted invalid GOTOOLCHAIN=%q (stdout=%s)", tc, res.Stdout)
		}
	}

	// 2. Real cmdr-dev module version match test passes
	unitRes, unitErr := harness.RunGoTestVerified(
		ctx,
		"tools/cmdr-dev",
		"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev",
		[]string{"TestGoSecurityEvidenceVersionMatchesModule"},
		"-run", "^TestGoSecurityEvidenceVersionMatchesModule$",
	)
	if !unitRes.Passed || unitErr != nil {
		t.Fatalf("T2.F3.5 FAIL: TestGoSecurityEvidenceVersionMatchesModule failed: %v\nOutput: %s", unitErr, unitRes.Output)
	}

	// 3. Real cmdr-dev sca-go fails closed when go.mod contains invalid toolchain directive `go 1.26.11-rc1`
	mockRoot := harness.CloneIsolatedRepo(t, "f03-bad-gomod")
	goModPath := filepath.Join(mockRoot, "tools", "cmdr-dev", "go.mod")
	modBytes, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("failed to read mock tools/cmdr-dev/go.mod: %v", err)
	}
	mutatedMod := strings.Replace(string(modBytes), "go 1.26.11", "go 1.26.11-rc1", 1)
	if mutatedMod == string(modBytes) {
		t.Fatalf("T2.F3.5 setup error: could not locate 'go 1.26.11' in tools/cmdr-dev/go.mod")
	}
	if err := os.WriteFile(goModPath, []byte(mutatedMod), 0644); err != nil {
		t.Fatalf("failed to write mutated go.mod: %v", err)
	}
	scaRes := harness.RunCmdrDev(ctx, "sca-go", "--root", mockRoot, "--json")
	if scaRes.ExitCode == 0 {
		t.Fatalf("T2.F3.5 FAIL: cmdr-dev sca-go succeeded despite invalid go 1.26.11-rc1 directive in go.mod")
	}

	t.Log("T2.F3.5 PASS: real Go CLI and cmdr-dev sca-go strictly reject malformed toolchain versions")
}
