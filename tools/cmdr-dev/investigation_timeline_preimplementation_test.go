package main

import (
	"context"
	"strings"
	"testing"
)

func TestInvestigationTimelinePreimplementationMustNotEmitMeasurements(t *testing.T) {
	for _, operation := range []string{"navigation"} {
		key := "builtin-investigation-timeline-" + operation + "-preimplementation-v1"
		if !knownBenchmarkHandlerKeys[key] || !knownDeepPerformanceHandlerKeys[key] {
			t.Fatalf("missing explicit handler allowlist %s", key)
		}
		_, cleanup, err := prepareBenchmarkHandler(t.TempDir(), BenchmarkWorkloadDefinition{HandlerKey: key})
		if cleanup != nil {
			cleanup()
		}
		if err == nil || !strings.Contains(err.Error(), "preimplementation") {
			t.Fatalf("benchmark %s must fail closed: %v", key, err)
		}
		_, err = runDeepPerformanceHandler(context.Background(), t.TempDir(), DeepPerformanceTarget{HandlerKey: key})
		if err == nil || !strings.Contains(err.Error(), "preimplementation") {
			t.Fatalf("deep benchmark %s must fail closed: %v", key, err)
		}
	}
}
