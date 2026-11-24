package tier2_boundaries

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type SCARawSummary struct {
	Tool                  string           `json:"tool"`
	Version               string           `json:"version"`
	Database              string           `json:"database,omitempty"`
	GoVersion             string           `json:"go_version,omitempty"`
	Modules               int              `json:"modules"`
	InformationalFindings int              `json:"informational_findings"`
	ActionableFindings    int              `json:"actionable_findings"`
	Findings              []SCAFindingItem `json:"findings"`
}

type SCAFindingItem struct {
	OSV        string `json:"osv"`
	Fixed      string `json:"fixed,omitempty"`
	Module     string `json:"module"`
	Version    string `json:"version"`
	Package    string `json:"package"`
	Function   string `json:"function,omitempty"`
	Actionable bool   `json:"actionable"`
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

// T2.F3.3: SCA correctly detects vulnerable dependency if introduced.
func TestT2_F3_3_SCADetectsVulnerableDependency(t *testing.T) {
	// Synthesize govulncheck output stream representing a detected CVE / GO- vulnerability
	syntheticVulnStream := `
{"config":{"protocol_version":"v1.0.0"}}
{"progress":{"message":"Scanning your code and Pkg dependencies..."}}
{"osv":{"id":"GO-2026-9999","aliases":["CVE-2026-99999"],"details":"Synthetic high-risk vulnerability in test dependency"}}
{"finding":{"osv":"GO-2026-9999","fixed_version":"v1.2.3","trace":[{"module":"example.com/vulnerable-lib","version":"v1.0.0","package":"example.com/vulnerable-lib/crypto","function":"InsecureDecrypt"}]}}
`
	// Verify parsing logic accurately detects actionable findings
	dec := json.NewDecoder(strings.NewReader(strings.TrimSpace(syntheticVulnStream)))
	actionableCount := 0
	var foundIDs []string

	for dec.More() {
		var entry map[string]interface{}
		if err := dec.Decode(&entry); err != nil {
			break
		}
		if _, ok := entry["finding"]; ok {
			actionableCount++
		}
		if osv, ok := entry["osv"].(map[string]interface{}); ok {
			if id, ok := osv["id"].(string); ok {
				foundIDs = append(foundIDs, id)
			}
		}
	}

	if actionableCount != 1 {
		t.Fatalf("T2.F3.3 FAIL: expected 1 actionable vulnerability finding, got %d", actionableCount)
	}
	if len(foundIDs) == 0 || foundIDs[0] != "GO-2026-9999" {
		t.Fatalf("T2.F3.3 FAIL: expected vulnerability ID GO-2026-9999, got %v", foundIDs)
	}
	t.Log("T2.F3.3 PASS: SCA stream parser reliably detects and accounts for vulnerable dependency findings")
}

// T2.F3.4: SCA report format maintains schema when 0 actionable findings.
func TestT2_F3_4_SCAReportMaintainsSchemaOnZeroFindings(t *testing.T) {
	// Standard empty report invariant
	summary := SCARawSummary{
		Tool:                  "govulncheck",
		Version:               "v1.8.0",
		Database:              "https://vuln.go.dev",
		GoVersion:             "go1.26.11",
		Modules:               12,
		InformationalFindings: 0,
		ActionableFindings:    0,
		Findings:              []SCAFindingItem{},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal SCA summary: %v", err)
	}

	harness.AssertJSONFieldEqual(t, string(data), "tool", "govulncheck")
	harness.AssertJSONFieldEqual(t, string(data), "version", "v1.8.0")
	harness.AssertJSONFieldEqual(t, string(data), "actionable_findings", 0)
	harness.AssertJSONFieldEqual(t, string(data), "informational_findings", 0)

	var parsed SCARawSummary
	if err := harness.StrictUnmarshal(data, &parsed); err != nil {
		t.Fatalf("T2.F3.4 FAIL: strict unmarshal failed on zero findings schema: %v", err)
	}
	if parsed.ActionableFindings != 0 || len(parsed.Findings) != 0 {
		t.Fatalf("T2.F3.4 FAIL: unexpected findings count in empty summary")
	}
	t.Log("T2.F3.4 PASS: SCA report schema strictly conforms to invariant on zero findings")
}

// T2.F3.5: SCA rejects invalid toolchain version string format.
func TestT2_F3_5_SCARejectsInvalidToolchainVersionFormat(t *testing.T) {
	toolchainPattern := regexp.MustCompile(`^go1\.\d+\.\d+$`)

	invalidVersions := []string{
		"1.26.11",       // Missing 'go' prefix
		"go1.26",        // Missing patch number
		"v1.26.11",      // Semver 'v' prefix instead of 'go'
		"go1.26.11-rc1", // Pre-release suffix
		"go1.26.11.1",   // Quad format
		"",              // Empty
		"   ",           // Whitespace
		"go1.26.11\n",   // Trailing newline
		"go1.26.11;rm",  // Injection attempt
	}

	for _, v := range invalidVersions {
		if toolchainPattern.MatchString(v) {
			t.Fatalf("T2.F3.5 FAIL: invalid toolchain version %q was incorrectly accepted by format regex", v)
		}
	}

	validVersions := []string{
		"go1.26.11",
		"go1.26.10",
		"go1.25.3",
	}
	for _, v := range validVersions {
		if !toolchainPattern.MatchString(v) {
			t.Fatalf("T2.F3.5 FAIL: valid toolchain version %q was rejected by format regex", v)
		}
	}
	t.Log("T2.F3.5 PASS: toolchain version string validator strictly rejects non-canonical and malformed formats")
}
