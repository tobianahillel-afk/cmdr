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

const eventInspectionRuntimeModule = "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"

type eventInspectionProbe struct {
	moduleRoot string
	binaryPath string
}

func prepareEventInspectionProbe(root string) (eventInspectionProbe, func(), error) {
	if _, err := resolveRepoPath(root, "product-runtime/event-inspection", false); err != nil {
		return eventInspectionProbe{}, func() {}, fmt.Errorf("event-inspection runtime is not implemented: %w", err)
	}
	tempDir, err := os.MkdirTemp(root, ".cmdr-event-inspection-perf-")
	if err != nil {
		return eventInspectionProbe{}, func() {}, err
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir) // #nosec G703 -- tempDir is created beneath the validated repository root.
	}

	goMod := []byte("module cmdr.local/event-inspection-perf\n\ngo " + engineeringGoVersion +
		"\n\nrequire " + eventInspectionRuntimeModule + " v0.0.0\n\nreplace " +
		eventInspectionRuntimeModule + " => ../product-runtime/event-inspection\n")
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "go.mod"), goMod); err != nil {
		cleanup()
		return eventInspectionProbe{}, func() {}, fmt.Errorf("write Event Inspection probe go.mod: %w", err)
	}
	if _, err := writeRepoFile(root, filepath.Join(tempDir, "main.go"), []byte(eventInspectionProbeSource)); err != nil {
		cleanup()
		return eventInspectionProbe{}, func() {}, fmt.Errorf("write Event Inspection probe source: %w", err)
	}

	binaryName := "event-inspection-perf"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)
	if _, err := runPerformanceProcess(tempDir, "go", "build", "-trimpath", "-o", binaryPath, "."); err != nil {
		cleanup()
		return eventInspectionProbe{}, func() {}, fmt.Errorf("build Event Inspection probe: %w", err)
	}
	return eventInspectionProbe{moduleRoot: tempDir, binaryPath: binaryPath}, cleanup, nil
}

func (probe eventInspectionProbe) run(iterations int, operation string) (pilotProjectionProbeObservation, error) {
	if iterations < 1 || iterations > 10_000_000 {
		return pilotProjectionProbeObservation{}, fmt.Errorf("Event Inspection iterations must be within 1..10000000")
	}
	if operation != "projection" && operation != "pivot" {
		return pilotProjectionProbeObservation{}, fmt.Errorf("unknown Event Inspection probe operation %q", operation)
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
		return observation, fmt.Errorf("decode Event Inspection %s probe: %w", operation, err)
	}
	if observation.Operation != operation || observation.Mode != "latency" ||
		observation.Operations != iterations || observation.ElapsedNS <= 0 ||
		observation.NSPerOperation <= 0 || math.IsNaN(observation.NSPerOperation) ||
		math.IsInf(observation.NSPerOperation, 0) {
		return observation, fmt.Errorf("invalid Event Inspection %s observation", operation)
	}
	return observation, nil
}

const eventInspectionProbeSource = `package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	eventinspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

var projectionSink eventinspection.Result
var pivotSink eventinspection.PivotResult

func main() {
	iterations := flag.Int("iterations", 20000, "operations")
	operation := flag.String("operation", "projection", "projection or pivot")
	flag.Parse()
	if *iterations < 1 || *iterations > 10000000 {
		fatal("invalid iterations")
	}
	if *operation != "projection" && *operation != "pivot" {
		fatal("invalid operation")
	}

	input := eventinspection.Input{
		TenantRef: "tenant-a",
		AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "env-prod",
		Permissions: []string{eventinspection.PermissionTelemetryRead},
		RawAccessDecision: eventinspection.AccessAllow,
		RenderedAccessDecision: eventinspection.AccessAllow,
		CorrelationID: "corr-event-perf-001",
		Event: eventinspection.Event{
			Ref: "evt-perf-001",
			TenantRef: "tenant-a",
			State: eventinspection.StateAvailable,
			Source: eventinspection.SourceMetadata{SourceRef:"source-a", DataSourceRef:"data-source-a", Kind:"telemetry"},
			Parser: eventinspection.ParserMetadata{ParserRef:"parser-a", Version:"v1", Status:"parsed"},
			SourceFields: map[string]string{"source.ip":"192.0.2.10"},
			NormalizedFields: map[string]string{"network.client.ip":"192.0.2.10"},
			MissingFields: []string{"user.name"},
			RawPayload: []byte("raw"),
			RenderedPayload: []byte("rendered"),
		},
	}
	projected := eventinspection.Project(input)
	if !projected.Allowed || projected.Projection == nil {
		fatal("projection preflight failed")
	}
	selection := eventinspection.PivotSelection{
		Field:"source.ip", Value:"192.0.2.10", ValuePresent:true, ValueAuthorized:true,
		TimeStart:"2026-09-25T10:00:00Z", TimeEnd:"2026-09-25T11:00:00Z",
		ReturnOrigin:"event-search:run-001:row-17",
	}
	if got := eventinspection.PreparePivot(projected.Projection, selection); !got.Allowed || got.Draft == nil {
		fatal("pivot preflight failed")
	}

	started := time.Now()
	for i := 0; i < *iterations; i++ {
		switch *operation {
		case "projection":
			projectionSink = eventinspection.Project(input)
		case "pivot":
			pivotSink = eventinspection.PreparePivot(projected.Projection, selection)
		}
	}
	elapsed := time.Since(started)
	if *operation == "projection" && (!projectionSink.Allowed || projectionSink.Projection == nil) {
		fatal("projection sink invalid")
	}
	if *operation == "pivot" && (!pivotSink.Allowed || pivotSink.Draft == nil) {
		fatal("pivot sink invalid")
	}
	if elapsed <= 0 {
		fatal("non-positive elapsed duration")
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

func fatal(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(2)
}
`
