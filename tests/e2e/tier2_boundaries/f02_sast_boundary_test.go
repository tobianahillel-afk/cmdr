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

type GosecRawReport struct {
	Issues []struct {
		Severity   string `json:"severity"`
		Confidence string `json:"confidence"`
		RuleID     string `json:"rule_id"`
		Details    string `json:"details"`
		File       string `json:"file"`
		Line       string `json:"line"`
	} `json:"Issues"`
	Stats struct {
		Files int `json:"files"`
		Lines int `json:"lines"`
		Nosec int `json:"nosec"`
		Found int `json:"found"`
	} `json:"Stats"`
}

// T2.F2.1: SAST fails closed if G204/G702 nosec annotation is absent.
func TestT2_F2_1_SASTFailsClosedWithoutNosec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "sast-nosec")
	goModContent := "module example.com/testmod\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create dummy go.mod: %v", err)
	}

	// Deliberate subprocess call without any #nosec comment
	code := `package main

import (
	"os/exec"
)

func RunUserCmd(arg string) error {
	cmd := exec.Command("cmd.exe", "/c", arg)
	return cmd.Run()
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644); err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	reportPath := filepath.Join(tmpDir, "report.json")
	res := harness.RunCmd(ctx, tmpDir, nil, "go", "run", "github.com/securego/gosec/v2/cmd/gosec@v2.28.0",
		"-no-fail", "-fmt=json", "-out="+reportPath, "-severity=medium", "-confidence=medium", "./...")
	if res.ExitCode != 0 && res.ExitCode != 1 {
		t.Fatalf("gosec execution error (exit %d): %s", res.ExitCode, res.Stderr)
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read gosec report: %v\nStdout: %s\nStderr: %s", err, res.Stdout, res.Stderr)
	}

	var report GosecRawReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to parse gosec JSON report: %v", err)
	}

	foundG204 := false
	for _, issue := range report.Issues {
		if issue.RuleID == "G204" {
			foundG204 = true
			break
		}
	}
	if !foundG204 {
		t.Fatalf("T2.F2.1 FAIL: expected G204 finding when #nosec is missing, but found 0 matching issues")
	}
	t.Log("T2.F2.1 PASS: SAST fails closed and detects unannotated subprocess invocation")
}

// T2.F2.2: SAST scans path with spaces without crashing.
func TestT2_F2_2_SASTScansPathWithSpaces(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	baseDir := harness.TempDir(t, "sast-space")
	spacedDir := filepath.Join(baseDir, "module with spaces")
	if err := os.MkdirAll(spacedDir, 0755); err != nil {
		t.Fatalf("failed to create directory with spaces: %v", err)
	}

	goModContent := "module example.com/spacedmod\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(spacedDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod in spaced dir: %v", err)
	}

	code := `package main

import "fmt"

func main() {
	fmt.Println("clean execution")
}
`
	if err := os.WriteFile(filepath.Join(spacedDir, "main.go"), []byte(code), 0644); err != nil {
		t.Fatalf("failed to write main.go in spaced dir: %v", err)
	}

	reportPath := filepath.Join(spacedDir, "report.json")
	res := harness.RunCmd(ctx, spacedDir, nil, "go", "run", "github.com/securego/gosec/v2/cmd/gosec@v2.28.0",
		"-no-fail", "-fmt=json", "-out="+reportPath, "./...")
	if res.ExitCode != 0 {
		t.Fatalf("T2.F2.2 FAIL: gosec crashed on directory with spaces: exit %d\nStderr: %s", res.ExitCode, res.Stderr)
	}

	if _, err := os.Stat(reportPath); err != nil {
		t.Fatalf("T2.F2.2 FAIL: expected report file in spaced dir was not generated: %v", err)
	}
	t.Log("T2.F2.2 PASS: SAST scanner successfully handles paths containing spaces")
}

// T2.F2.3: SAST handles empty file gracefully.
func TestT2_F2_3_SASTHandlesEmptyFileGracefully(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "sast-empty")
	goModContent := "module example.com/emptymod\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	// 0-byte file
	if err := os.WriteFile(filepath.Join(tmpDir, "empty.go"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to write 0-byte go file: %v", err)
	}

	reportPath := filepath.Join(tmpDir, "report.json")
	res := harness.RunCmd(ctx, tmpDir, nil, "go", "run", "github.com/securego/gosec/v2/cmd/gosec@v2.28.0",
		"-no-fail", "-fmt=json", "-out="+reportPath, "./...")
	if res.ExitCode != 0 {
		t.Fatalf("T2.F2.3 FAIL: gosec exited abnormally on 0-byte file: exit %d\nStderr: %s", res.ExitCode, res.Stderr)
	}
	t.Log("T2.F2.3 PASS: SAST scanner gracefully tolerates 0-byte files without panics")
}

// T2.F2.4: SAST detects command injection patterns.
func TestT2_F2_4_SASTDetectsCommandInjectionPattern(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "sast-inject")
	goModContent := "module example.com/injectmod\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	code := `package main

import (
	"os/exec"
)

func ExecuteInjected(untrustedCmd string) error {
	// Obvious command execution vulnerability pattern
	c := exec.Command("sh", "-c", untrustedCmd)
	return c.Run()
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "inject.go"), []byte(code), 0644); err != nil {
		t.Fatalf("failed to create inject.go: %v", err)
	}

	reportPath := filepath.Join(tmpDir, "report.json")
	res := harness.RunCmd(ctx, tmpDir, nil, "go", "run", "github.com/securego/gosec/v2/cmd/gosec@v2.28.0",
		"-no-fail", "-fmt=json", "-out="+reportPath, "-severity=medium", "-confidence=medium", "./...")
	if res.ExitCode != 0 && res.ExitCode != 1 {
		t.Fatalf("gosec failed to execute: %s", res.Stderr)
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read gosec report: %v", err)
	}

	var report GosecRawReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("failed to parse gosec JSON report: %v", err)
	}

	detected := false
	for _, issue := range report.Issues {
		if issue.RuleID == "G204" {
			detected = true
			break
		}
	}
	if !detected {
		t.Fatalf("T2.F2.4 FAIL: command injection pattern G204 was not detected by SAST")
	}
	t.Log("T2.F2.4 PASS: command injection pattern detected and flagged by SAST")
}

// T2.F2.5: SAST JSON output formatting invariant.
func TestT2_F2_5_SASTJSONOutputFormattingInvariant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev sast-go --json failed with exit code %d: %s", res.ExitCode, res.Stderr)
	}

	// Validate JSON fields invariant
	harness.AssertJSONFieldEqual(t, res.Stdout, "tool", "gosec")
	harness.AssertJSONFieldEqual(t, res.Stdout, "version", "v2.28.0")
	harness.AssertJSONFieldExists(t, res.Stdout, "scan_root")
	harness.AssertJSONFieldExists(t, res.Stdout, "findings")
	harness.AssertJSONFieldExists(t, res.Stdout, "by_severity")

	var report map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(res.Stdout)), &report); err != nil {
		t.Fatalf("T2.F2.5 FAIL: invalid JSON formatting produced by sast-go: %v", err)
	}

	findings, ok := report["findings"].(float64)
	if !ok {
		t.Fatalf("T2.F2.5 FAIL: 'findings' field is not a number: %T", report["findings"])
	}
	if findings != 0 {
		t.Fatalf("T2.F2.5 FAIL: expected 0 findings on certified codebase, got %v", findings)
	}
	t.Log("T2.F2.5 PASS: SAST JSON output conforms strictly to invariant schema")
}
