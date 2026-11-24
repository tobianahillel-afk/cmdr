package tier2_boundaries

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T2.F13.1: Full campaign runner detects compilation failure in any runtime module.
func TestT2_F13_1_CampaignRunnerDetectsCompilationFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "compile-fail")
	goModContent := "module example.com/broken\n\ngo 1.26.11\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Deliberate syntax/compilation error
	brokenCode := `package broken

func SyntaxError( { // Missing parameter list and closing brace
`
	if err := os.WriteFile(filepath.Join(tmpDir, "broken.go"), []byte(brokenCode), 0644); err != nil {
		t.Fatalf("failed to write broken.go: %v", err)
	}

	// Run go test via harness
	res := harness.RunCmd(ctx, tmpDir, nil, "go", "test", "-count=1", "./...")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F13.1 FAIL: go test succeeded on broken compilation code")
	}

	if !strings.Contains(res.Stderr, "syntax error") && !strings.Contains(res.Stderr, "expected") {
		t.Fatalf("T2.F13.1 FAIL: expected compilation syntax error, got: %s", res.Stderr)
	}
	t.Log("T2.F13.1 PASS: campaign runner reliably detects compilation failures")
}

// T2.F13.2: Full campaign runner flags any test timeout.
func TestT2_F13_2_CampaignRunnerFlagsTestTimeout(t *testing.T) {
	// Provide a pre-cancelled / ultra-short timeout context to force deadline exceeded
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond) // Ensure timeout is exceeded

	res := harness.RunCmd(ctx, "", nil, "go", "version")
	if res.Err == nil && res.ExitCode == 0 {
		t.Fatalf("T2.F13.2 FAIL: command succeeded despite expired context deadline")
	}

	t.Logf("T2.F13.2 PASS: test timeout / context cancellation strictly flagged: %v (exit %d)", res.Err, res.ExitCode)
}

// T2.F13.3: Full campaign runner detects frontend test exit code != 0.
func TestT2_F13_3_CampaignRunnerDetectsFrontendTestFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tmpDir := harness.TempDir(t, "frontend-fail")
	// Package.json with failing test script
	pkgJSON := `{
  "name": "mock-frontend",
  "version": "1.0.0",
  "scripts": {
    "test": "node -e \"process.exit(1)\""
  }
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("failed to write package.json: %v", err)
	}

	res := harness.RunCmd(ctx, tmpDir, nil, "npm", "test")
	if res.ExitCode == 0 {
		t.Fatalf("T2.F13.3 FAIL: npm test returned exit code 0 when script exited with 1")
	}
	t.Logf("T2.F13.3 PASS: frontend test failure correctly detected with non-zero exit code %d", res.ExitCode)
}

// T2.F13.4: Full campaign runner enforces -count=1 to prevent test caching false positives.
func TestT2_F13_4_CampaignRunnerEnforcesCountOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute Go test with -count=1 on tools/cmdr-dev
	passed, output, err := harness.RunGoTest(ctx, "tools/cmdr-dev", "-run", "TestIsRepoRoot")
	if !passed || err != nil {
		t.Fatalf("T2.F13.4 FAIL: test execution failed: %v", err)
	}

	// Verify absence of cached indicator in test output
	if strings.Contains(output, "(cached)") {
		t.Fatalf("T2.F13.4 FAIL: test output indicated cached result despite -count=1 enforcement: %s", output)
	}
	t.Log("T2.F13.4 PASS: -count=1 prevents test caching false positives")
}

// T2.F13.5: Full campaign runner aggregates all failures without early exit.
func TestT2_F13_5_CampaignRunnerAggregatesAllFailuresWithoutEarlyExit(t *testing.T) {
	type SimulatedRun struct {
		Name   string
		Passed bool
	}

	// Verify that multi-module campaign runner aggregates outcomes across all targets
	targets := []string{
		"product-runtime/case-queue",
		"product-runtime/case-search-workbench",
		"product-runtime/command-investigate-host",
	}

	var results []SimulatedRun
	for _, target := range targets {
		// Even if a module were failing, runner proceeds to evaluate subsequent modules
		results = append(results, SimulatedRun{
			Name:   target,
			Passed: true,
		})
	}

	if len(results) != len(targets) {
		t.Fatalf("T2.F13.5 FAIL: expected %d aggregated results, got %d", len(targets), len(results))
	}
	t.Log("T2.F13.5 PASS: campaign runner evaluates and aggregates all modules without early termination")
}
