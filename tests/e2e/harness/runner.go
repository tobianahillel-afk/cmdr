package harness

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// RuntimeGoModules lists all 12 Go product-runtime modules.
var RuntimeGoModules = []string{
	"product-runtime/case-queue",
	"product-runtime/case-search-workbench",
	"product-runtime/command-investigate-host",
	"product-runtime/context-envelope",
	"product-runtime/event-inspection",
	"product-runtime/event-inspection-host",
	"product-runtime/event-search",
	"product-runtime/hunt-management",
	"product-runtime/investigation-timeline",
	"product-runtime/saved-query-assets",
	"product-runtime/studio-catalog",
	"product-runtime/work-queue",
}

// ModuleTestResult records the verification results for an individual module.
type ModuleTestResult struct {
	Module     string
	VetPassed  bool
	VetOutput  string
	TestPassed bool
	TestOutput string
	Duration   time.Duration
	Err        error
}

// VerifiedGoTestResult holds structured `go test -json` execution proof.
type VerifiedGoTestResult struct {
	Passed            bool
	ExitCode          int
	Output            string
	RawStdout         string
	RawStderr         string
	ExecutedTests     []string
	PassedTests       []string
	FailedTests       []string
	SkippedTests      []string
	CompilationFailed bool
}

type goTestJSONEvent struct {
	Action      string `json:"Action"`
	Package     string `json:"Package"`
	ImportPath  string `json:"ImportPath"`
	Test        string `json:"Test"`
	Output      string `json:"Output"`
	FailedBuild string `json:"FailedBuild"`
}

type testExecState struct {
	seenRun  bool
	seenPass bool
	seenFail bool
	seenSkip bool
}

// VerifyGoTestJSON parses `go test -json` NDJSON output and verifies that every expected test/subtest
// in expectedPackage actually started (`run`) and passed (`pass`) without failing or skipping.
// When expectedTests is empty, requires at least one test in expectedPackage (or any package if empty)
// to have executed and passed.
func VerifyGoTestJSON(jsonData []byte, expectedPackage string, expectedTests []string) (VerifiedGoTestResult, error) {
	var res VerifiedGoTestResult
	var outBuf strings.Builder

	type key struct {
		pkg  string
		test string
	}

	states := make(map[key]*testExecState)
	pkgFailed := make(map[string]bool)
	var testOrder []key

	scanner := bufio.NewScanner(bytes.NewReader(jsonData))
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var ev goTestJSONEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			outBuf.WriteString(line)
			outBuf.WriteByte('\n')
			continue
		}

		if ev.Output != "" {
			outBuf.WriteString(ev.Output)
		}

		pkg := ev.Package
		if pkg == "" {
			pkg = ev.ImportPath
		}

		if ev.Action == "build-fail" || ev.FailedBuild != "" {
			res.CompilationFailed = true
			if pkg != "" {
				pkgFailed[pkg] = true
			}
			continue
		}

		if ev.Test == "" {
			if ev.Action == "fail" && pkg != "" {
				pkgFailed[pkg] = true
			}
			continue
		}

		k := key{pkg: pkg, test: ev.Test}
		st, exists := states[k]
		if !exists {
			st = &testExecState{}
			states[k] = st
			testOrder = append(testOrder, k)
		}

		switch ev.Action {
		case "run":
			st.seenRun = true
		case "pass":
			st.seenPass = true
		case "fail":
			st.seenFail = true
		case "skip":
			st.seenSkip = true
		}
	}

	if err := scanner.Err(); err != nil {
		res.Output = strings.TrimSpace(outBuf.String())
		return res, fmt.Errorf("failed to scan go test -json output: %w", err)
	}

	executedSet := make(map[string]bool)
	passedSet := make(map[string]bool)
	failedSet := make(map[string]bool)
	skippedSet := make(map[string]bool)

	for _, k := range testOrder {
		st := states[k]
		if st.seenRun && !executedSet[k.test] {
			executedSet[k.test] = true
			res.ExecutedTests = append(res.ExecutedTests, k.test)
		}
		if st.seenFail && !failedSet[k.test] {
			failedSet[k.test] = true
			res.FailedTests = append(res.FailedTests, k.test)
		}
		if st.seenSkip && !skippedSet[k.test] {
			skippedSet[k.test] = true
			res.SkippedTests = append(res.SkippedTests, k.test)
		}
		if st.seenRun && st.seenPass && !st.seenFail && !st.seenSkip && !passedSet[k.test] {
			passedSet[k.test] = true
			res.PassedTests = append(res.PassedTests, k.test)
		}
	}

	res.Output = strings.TrimSpace(outBuf.String())

	if res.CompilationFailed {
		return res, fmt.Errorf("go test compilation failed")
	}
	if len(res.FailedTests) > 0 || len(pkgFailed) > 0 {
		return res, fmt.Errorf("go test reported failures (failed tests: %v)", res.FailedTests)
	}

	if len(expectedTests) == 0 {
		matchedCount := 0
		for _, k := range testOrder {
			if expectedPackage != "" && k.pkg != expectedPackage {
				continue
			}
			st := states[k]
			if st.seenRun && st.seenPass && !st.seenFail && !st.seenSkip {
				matchedCount++
			}
		}
		if matchedCount == 0 {
			return res, fmt.Errorf("zero tests executed and passed in package %q", expectedPackage)
		}
		res.Passed = true
		return res, nil
	}

	for _, exp := range expectedTests {
		found := false
		for _, k := range testOrder {
			if k.test != exp {
				continue
			}
			if expectedPackage != "" && k.pkg != expectedPackage {
				continue
			}
			st := states[k]
			if st.seenRun && st.seenPass && !st.seenFail && !st.seenSkip {
				found = true
				break
			}
		}
		if !found {
			return res, fmt.Errorf("expected test %q did not run and pass in package %q", exp, expectedPackage)
		}
	}

	res.Passed = true
	return res, nil
}

// RunGoTestVerified executes `go test -json -count=1 <extraArgs...> ./...` and verifies that
// every expected test/subtest actually ran and passed without failing or skipping.
func RunGoTestVerified(ctx context.Context, moduleRelPath, expectedPackage string, expectedTests []string, extraArgs ...string) (VerifiedGoTestResult, error) {
	root, err := RepoRoot()
	if err != nil {
		return VerifiedGoTestResult{ExitCode: -1}, err
	}
	dir := filepath.Join(root, moduleRelPath)
	args := []string{"test", "-json", "-count=1"}
	args = append(args, extraArgs...)
	args = append(args, "./...")

	cmdRes := RunCmd(ctx, dir, nil, "go", args...)
	verified, verifyErr := VerifyGoTestJSON([]byte(cmdRes.Stdout), expectedPackage, expectedTests)
	verified.ExitCode = cmdRes.ExitCode
	verified.RawStdout = cmdRes.Stdout
	verified.RawStderr = cmdRes.Stderr

	if strings.TrimSpace(cmdRes.Stderr) != "" {
		if verified.Output == "" {
			verified.Output = strings.TrimSpace(cmdRes.Stderr)
		} else {
			verified.Output = verified.Output + "\n" + strings.TrimSpace(cmdRes.Stderr)
		}
	}

	if cmdRes.Err != nil || cmdRes.ExitCode != 0 {
		verified.Passed = false
		if verifyErr != nil {
			return verified, fmt.Errorf("go test failed (exit %d): %v (%w)", cmdRes.ExitCode, cmdRes.Err, verifyErr)
		}
		return verified, fmt.Errorf("go test failed (exit %d): %v", cmdRes.ExitCode, cmdRes.Err)
	}
	if verifyErr != nil {
		verified.Passed = false
		return verified, verifyErr
	}

	return verified, nil
}

// RunGoVet runs `go vet ./...` in the specified module directory.
func RunGoVet(ctx context.Context, moduleRelPath string) (bool, string, error) {
	root, err := RepoRoot()
	if err != nil {
		return false, "", err
	}
	dir := filepath.Join(root, moduleRelPath)
	res := RunCmd(ctx, dir, nil, "go", "vet", "./...")
	output := strings.TrimSpace(res.Stdout + "\n" + res.Stderr)
	return res.ExitCode == 0, output, res.Err
}

// RunGoTest runs `go test -count=1 ./...` in the specified module directory.
func RunGoTest(ctx context.Context, moduleRelPath string, extraArgs ...string) (bool, string, error) {
	root, err := RepoRoot()
	if err != nil {
		return false, "", err
	}
	dir := filepath.Join(root, moduleRelPath)
	args := []string{"test", "-v", "-count=1"}
	args = append(args, extraArgs...)
	args = append(args, "./...")
	res := RunCmd(ctx, dir, nil, "go", args...)
	output := strings.TrimSpace(res.Stdout + "\n" + res.Stderr)
	return res.ExitCode == 0, output, res.Err
}

// RunFrontendTest executes `npm test` or `node --test` in product-runtime/event-search-frontend.
func RunFrontendTest(ctx context.Context) (bool, string, error) {
	root, err := RepoRoot()
	if err != nil {
		return false, "", err
	}
	dir := filepath.Join(root, "product-runtime", "event-search-frontend")
	res := RunCmd(ctx, dir, nil, "npm", "test")
	output := strings.TrimSpace(res.Stdout + "\n" + res.Stderr)
	return res.ExitCode == 0, output, res.Err
}

// RunCampaign runs vet and tests across all runtime Go modules and cmdr-dev.
func RunCampaign(ctx context.Context) ([]ModuleTestResult, bool) {
	var results []ModuleTestResult
	allPassed := true

	modules := append([]string{}, RuntimeGoModules...)
	modules = append(modules, "tools/cmdr-dev")

	for _, mod := range modules {
		start := time.Now()
		vPass, vOut, _ := RunGoVet(ctx, mod)
		tPass, tOut, err := RunGoTest(ctx, mod)
		dur := time.Since(start)

		if !vPass || !tPass {
			allPassed = false
		}

		results = append(results, ModuleTestResult{
			Module:     mod,
			VetPassed:  vPass,
			VetOutput:  vOut,
			TestPassed: tPass,
			TestOutput: tOut,
			Duration:   dur,
			Err:        err,
		})
	}

	return results, allPassed
}
