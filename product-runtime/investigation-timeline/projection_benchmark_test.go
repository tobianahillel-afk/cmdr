package investigationtimeline

import (
	"fmt"
	"testing"
	"time"
)

// Synthetic in-process throughput only: no source provider, DB, UI,
// remote authorization or production latency is exercised.
func benchmarkTimelineInput(count int) Input {
	input := Input{Scope: Scope{
		Tenant: "tenant-benchmark", Environment: "test",
		CasePermission: "allow", TimelinePermission: "allow",
	}}
	input.Entries = make([]Entry, count)
	base := time.Date(2026, time.October, 1, 10, 0, 0, 0, time.UTC)
	for i := range input.Entries {
		id := fmt.Sprintf("entry-%05d", i)
		kind := "observed"
		if i%2 != 0 {
			kind = "inferred"
		}
		input.Entries[i] = Entry{
			ID: id, Kind: kind, Tenant: input.Scope.Tenant,
			SourcePermission: "allow", Timezone: "+00:00",
			SourceTime: base.Add(time.Duration(count-1-i) * time.Second).Format(time.RFC3339),
			SourceRef:  "telemetry:" + id, Provenance: "synthetic:" + id,
			Version: 1,
		}
	}
	return input
}

func benchmarkProjectTimeline(b *testing.B, count int) {
	input := benchmarkTimelineInput(count)
	firstID := fmt.Sprintf("entry-%05d", count-1)
	lastID := "entry-00000"
	preflight := Project(input)
	if preflight.Status != "allow" || len(preflight.Entries) != count ||
		preflight.Entries[0].ID != firstID ||
		preflight.Entries[count-1].ID != lastID {
		b.Fatalf("benchmark fixture rejected: status=%s rows=%d", preflight.Status, len(preflight.Entries))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Project(input)
		if result.Status != "allow" || len(result.Entries) != count ||
			result.Entries[0].ID != firstID ||
			result.Entries[count-1].ID != lastID {
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
