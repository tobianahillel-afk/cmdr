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

// Full-row oracle is deliberately excluded from b.N timing.
func benchmarkProjectionMatchesInput(input Input, result Result) bool {
	if result.Status != "allow" || len(result.Entries) != len(input.Entries) {
		return false
	}
	for i, got := range result.Entries {
		want := input.Entries[len(input.Entries)-1-i]
		if got.ID != want.ID || got.Kind != want.Kind ||
			got.SourceRef != want.SourceRef || got.Provenance != want.Provenance ||
			got.SourceTime != want.SourceTime || got.Timezone != want.Timezone ||
			got.Version != want.Version || got.Supersedes != want.Supersedes ||
			got.TimeUncertain {
			return false
		}
	}
	return true
}

func benchmarkProjectTimeline(b *testing.B, count int) {
	input := benchmarkTimelineInput(count)
	firstID := fmt.Sprintf("entry-%05d", count-1)
	lastID := "entry-00000"
	preflight := Project(input)
	if !benchmarkProjectionMatchesInput(input, preflight) {
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
	b.StopTimer()
	if !benchmarkProjectionMatchesInput(input, Project(input)) {
		b.Fatal("full-row postflight failed")
	}
}

func BenchmarkProject500Entries(b *testing.B) {
	benchmarkProjectTimeline(b, 500)
}

func BenchmarkProject5000Entries(b *testing.B) {
	benchmarkProjectTimeline(b, 5000)
}

func TestBenchmarkFullRowOracleRejectsInteriorMutation(t *testing.T) {
	input := benchmarkTimelineInput(5000)
	result := Project(input)
	if !benchmarkProjectionMatchesInput(input, result) {
		t.Fatal("valid source rejected")
	}
	for _, change := range []func(*ProjectedEntry){
		func(e *ProjectedEntry) { e.ID = "counterfeit" },
		func(e *ProjectedEntry) { e.Kind = "observed" },
		func(e *ProjectedEntry) { e.SourceTime = "2000-01-01T00:00:00Z" },
		func(e *ProjectedEntry) { e.Provenance = "invented" },
	} {
		modified := result
		modified.Entries = append([]ProjectedEntry(nil), result.Entries...)
		change(&modified.Entries[2500])
		if benchmarkProjectionMatchesInput(input, modified) {
			t.Fatal("accepted corrupt interior row")
		}
	}
}
