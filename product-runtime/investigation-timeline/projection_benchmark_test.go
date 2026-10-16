package investigationtimeline

import (
	"fmt"
	"testing"
)

// Synthetic in-process throughput only: no source provider, DB, UI,
// remote authorization or production latency is exercised.
func benchmarkTimelineInput(count int) Input {
	input := Input{Scope: Scope{
		Tenant: "tenant-benchmark", Environment: "test",
		CasePermission: "allow", TimelinePermission: "allow",
	}}
	input.Entries = make([]Entry, count)
	for i := range input.Entries {
		id := fmt.Sprintf("entry-%05d", i)
		input.Entries[i] = Entry{
			ID: id, Kind: "observed", Tenant: input.Scope.Tenant,
			SourcePermission: "allow", Timezone: "+00:00",
			SourceTime: "2026-10-01T10:00:00Z",
			SourceRef:  "telemetry:" + id, Provenance: "synthetic:" + id,
			Version: 1,
		}
	}
	return input
}

func benchmarkProjectTimeline(b *testing.B, count int) {
	input := benchmarkTimelineInput(count)
	preflight := Project(input)
	if preflight.Status != "allow" || len(preflight.Entries) != count {
		b.Fatalf("benchmark fixture rejected: status=%s rows=%d", preflight.Status, len(preflight.Entries))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Project(input)
		if result.Status != "allow" || len(result.Entries) != count {
			b.Fatalf("benchmark lost output: status=%s rows=%d", result.Status, len(result.Entries))
		}
	}
}

func BenchmarkProject500Entries(b *testing.B) {
	benchmarkProjectTimeline(b, 500)
}

func BenchmarkProject5000Entries(b *testing.B) {
	benchmarkProjectTimeline(b, 5000)
}
