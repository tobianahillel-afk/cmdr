package harness

import (
	"context"
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
	args := []string{"test", "-count=1"}
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
