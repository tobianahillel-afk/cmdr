package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const savedQueryAssetsRuntimeModule = "github.com/tobianahillel-afk/cmdr/product-runtime/saved-query-assets"

type savedQueryAssetsProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareSavedQueryAssetsProbe(root string) (savedQueryAssetsProbe, func(), error) {
	if _, err := resolveRepoPath(root, "product-runtime/saved-query-assets", false); err != nil {
		return savedQueryAssetsProbe{}, func() {}, fmt.Errorf("saved-query-assets runtime is not implemented: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-saved-query-assets-perf-")
	if err != nil {
		return savedQueryAssetsProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- tempDir is created beneath the validated repository root.
	}

	goMod := []byte("module cmdr.local/saved-query-assets-perf\n\ngo " + engineeringGoVersion +
		"\n\nrequire " + savedQueryAssetsRuntimeModule + " v0.0.0\n\nreplace " +
		savedQueryAssetsRuntimeModule + " => ../product-runtime/saved-query-assets\n")
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "go.mod"), goMod); err != nil {
		cleanup()
		return savedQueryAssetsProbe{}, func() {}, fmt.Errorf("write Saved Query Assets probe go.mod: %w", err)
	}
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "main.go"), []byte(savedQueryAssetsProbeSource)); err != nil {
		cleanup()
		return savedQueryAssetsProbe{}, func() {}, fmt.Errorf("write Saved Query Assets probe source: %w", err)
	}

	binaryName := "saved-query-assets-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(tempDir, "go", "build", "-trimpath", "-o", binaryPath, "."); err != nil {
		cleanup()
		return savedQueryAssetsProbe{}, func() {}, fmt.Errorf("build Saved Query Assets probe: %w", err)
	}
	return savedQueryAssetsProbe{moduleRoot: tempDir, binaryPath: binaryPath}, cleanup, nil
}

func (probe savedQueryAssetsProbe) run(iterations int, operation string) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 10_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Saved Query Assets iterations must be within 1..10000000")
	}
	if operation != "projection" && operation != "handoff" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("unknown Saved Query Assets probe operation %q", operation)
	}
	output, err := runPerformanceProcess(
		probe.moduleRoot,
		probe.binaryPath,
		"-iterations", strconv.Itoa(iterations),
		"-operation", operation,
	)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(output), &observation); err != nil {
		return observation, fmt.Errorf("decode Saved Query Assets %s probe: %w", operation, err)
	}
	if observation.Operation != operation || observation.Mode != "latency" ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Saved Query Assets %s observation", operation)
	}
	return observation, nil
}

const savedQueryAssetsProbeSource = `package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	savedqueryassets "github.com/tobianahillel-afk/cmdr/product-runtime/saved-query-assets"
)

var projectionSink savedqueryassets.Result
var handoffSink savedqueryassets.HandoffResult

func main() {
	iterations := flag.Int("iterations", 20000, "operations")
	operation := flag.String("operation", "projection", "projection or handoff")
	flag.Parse()
	if *iterations < 1 || *iterations > 10000000 {
		fatal("invalid iterations")
	}
	if *operation != "projection" && *operation != "handoff" {
		fatal("invalid operation")
	}

	input := validInput()
	projected := savedqueryassets.Project(input)
	if !projected.Allowed || projected.Projection == nil || !projected.HandoffEligible {
		fatal("projection preflight failed")
	}
	handoffInput := savedqueryassets.HandoffInput{
		TenantRef: "tenant-a",
		EnvironmentRef: "env-prod",
		ParameterValues: []savedqueryassets.ParameterValue{
			{Name: "hostname", OpaqueValue: []byte("host-17")},
			{Name: "window", OpaqueValue: []byte("15m")},
		},
	}
	preflightHandoff := savedqueryassets.BuildExecutionHandoff(projected, handoffInput)
	if !preflightHandoff.Allowed || preflightHandoff.Draft == nil {
		fatal("handoff preflight failed")
	}

	started := time.Now()
	switch *operation {
	case "projection":
		for i := 0; i < *iterations; i++ {
			projectionSink = savedqueryassets.Project(input)
		}
		if !projectionSink.Allowed || projectionSink.Projection == nil {
			fatal("projection sink invalid")
		}
	case "handoff":
		for i := 0; i < *iterations; i++ {
			handoffSink = savedqueryassets.BuildExecutionHandoff(projected, handoffInput)
		}
		if !handoffSink.Allowed || handoffSink.Draft == nil ||
			handoffSink.SearchJobCreated || handoffSink.SearchJobExecuted {
			fatal("handoff sink invalid")
		}
	}
	elapsed := time.Since(started)
	if elapsed <= 0 {
		fatal("invalid elapsed time")
	}
	result := map[string]any{
		"operation": *operation,
		"mode": "latency",
		"operations": *iterations,
		"elapsed_ns": elapsed.Nanoseconds(),
		"ns_per_operation": float64(elapsed.Nanoseconds()) / float64(*iterations),
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fatal(err.Error())
	}
}

func validInput() savedqueryassets.Input {
	return savedqueryassets.Input{
		TenantRef: "tenant-a",
		AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "env-prod",
		AssetID: "asset-perf-001",
		AssetKind: savedqueryassets.AssetSavedSearch,
		AssetAccessDecision: savedqueryassets.AccessAllow,
		QueryRef: "qry-perf-001",
		QueryTenantRef: "tenant-a",
		QueryVersion: "version-7",
		QueryAccessDecision: savedqueryassets.AccessAllow,
		AuthorRef: "analyst-perf-1",
		ValidationState: savedqueryassets.ValidationValidated,
		ParameterNames: []string{"hostname", "window"},
		Sources: []savedqueryassets.Source{{
			ID: "source-edr",
			TenantRef: "tenant-a",
			State: savedqueryassets.SourceAvailable,
			Fields: []savedqueryassets.Field{
				{Name: "host.name", State: savedqueryassets.FieldAvailable},
				{Name: "event.action", State: savedqueryassets.FieldAvailable},
			},
		}},
	}
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(2)
}
`
