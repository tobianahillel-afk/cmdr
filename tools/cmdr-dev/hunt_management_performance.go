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

const huntManagementRuntimeModule = "github.com/tobianahillel-afk/cmdr/product-runtime/hunt-management"

type huntManagementProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareHuntManagementProbe(root string) (huntManagementProbe, func(), error) {
	if _, err := resolveRepoPath(root, "product-runtime/hunt-management", false); err != nil {
		return huntManagementProbe{}, func() {}, fmt.Errorf("hunt-management runtime is not implemented: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-hunt-management-perf-")
	if err != nil {
		return huntManagementProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- tempDir is created beneath the validated repository root.
	}

	goMod := []byte("module cmdr.local/hunt-management-perf\n\ngo " + engineeringGoVersion +
		"\n\nrequire " + huntManagementRuntimeModule + " v0.0.0\n\nreplace " +
		huntManagementRuntimeModule + " => ../product-runtime/hunt-management\n")
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "go.mod"), goMod); err != nil {
		cleanup()
		return huntManagementProbe{}, func() {}, fmt.Errorf("write Hunt Management probe go.mod: %w", err)
	}
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "main.go"), []byte(huntManagementProbeSource)); err != nil {
		cleanup()
		return huntManagementProbe{}, func() {}, fmt.Errorf("write Hunt Management probe source: %w", err)
	}

	binaryName := "hunt-management-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(tempDir, "go", "build", "-trimpath", "-o", binaryPath, "."); err != nil {
		cleanup()
		return huntManagementProbe{}, func() {}, fmt.Errorf("build Hunt Management probe: %w", err)
	}
	return huntManagementProbe{moduleRoot: tempDir, binaryPath: binaryPath}, cleanup, nil
}

func (probe huntManagementProbe) run(iterations int) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 10_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Hunt Management iterations must be within 1..10000000")
	}
	output, err := runPerformanceProcess(
		probe.moduleRoot,
		probe.binaryPath,
		"-iterations", strconv.Itoa(iterations),
	)
	if err != nil {
		return pilotProjectionProbeObservation{}, err
	}
	var observation pilotProjectionProbeObservation
	if err := json.Unmarshal([]byte(output), &observation); err != nil {
		return observation, fmt.Errorf("decode Hunt Management projection probe: %w", err)
	}
	if observation.Operation != "projection" || observation.Mode != "latency" ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Hunt Management projection observation")
	}
	return observation, nil
}

const huntManagementProbeSource = `package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	huntmanagement "github.com/tobianahillel-afk/cmdr/product-runtime/hunt-management"
)

var projectionSink huntmanagement.Result

func main() {
	iterations := flag.Int("iterations", 20000, "operations")
	flag.Parse()
	if *iterations < 1 || *iterations > 10000000 {
		fatal("invalid iterations")
	}

	input := huntmanagement.Input{
		TenantRef: "tenant-a",
		AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "env-prod",
		Question: "What explains the suspicious authentication sequence?",
		Scope: "identity and endpoint telemetry",
		TimeStart: "2026-09-28T08:00:00Z",
		TimeEnd: "2026-09-28T09:00:00Z",
		OwnerRef: "owner-1",
		ContributorRefs: []string{"analyst-a", "analyst-b"},
		References: []huntmanagement.Reference{
			{
				Kind: huntmanagement.ReferenceQuery,
				ID: "qry-perf-001",
				TenantRef: "tenant-a",
				AccessDecision: huntmanagement.AccessAllow,
				State: huntmanagement.StateAvailable,
			},
			{
				Kind: huntmanagement.ReferenceSearchJob,
				ID: "job-perf-001",
				TenantRef: "tenant-a",
				AccessDecision: huntmanagement.AccessAllow,
				QueryRef: "qry-perf-001",
				State: huntmanagement.StateAvailable,
			},
			{
				Kind: huntmanagement.ReferenceCase,
				ID: "case-perf-001",
				TenantRef: "tenant-a",
				AccessDecision: huntmanagement.AccessAllow,
				State: huntmanagement.StatePartial,
			},
		},
		CorrelationID: "corr-hunt-perf-001",
		WorkspaceProvenance: "caller-owned-draft",
	}
	preflight := huntmanagement.Project(input)
	if !preflight.Allowed || preflight.Projection == nil {
		fatal("projection preflight failed")
	}

	started := time.Now()
	for i := 0; i < *iterations; i++ {
		projectionSink = huntmanagement.Project(input)
	}
	elapsed := time.Since(started)
	if !projectionSink.Allowed || projectionSink.Projection == nil || elapsed <= 0 {
		fatal("projection sink invalid")
	}
	result := map[string]any{
		"operation": "projection",
		"mode": "latency",
		"operations": *iterations,
		"elapsed_ns": elapsed.Nanoseconds(),
		"ns_per_operation": float64(elapsed.Nanoseconds()) / float64(*iterations),
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fatal(err.Error())
	}
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(2)
}
`
