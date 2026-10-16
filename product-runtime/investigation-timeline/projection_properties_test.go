package investigationtimeline

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

// Fixed-seed permutations must not change a same-instant Timeline's content,
// ordering or observed/inferred classification. This is independent of the
// declared contract fixtures and does not need a new runtime dependency.
func TestProjectionPermutationInvariantWithMixedKinds(t *testing.T) {
	in := allowedFixture()
	in.Entries = make([]Entry, 40)
	for i := range in.Entries {
		entry := allowedFixture().Entries[0]
		entry.ID = fmt.Sprintf("entry-%03d", i)
		entry.SourceRef = "telemetry:" + entry.ID
		entry.Provenance = "source:" + entry.ID
		if i%2 != 0 {
			entry.Kind = "inferred"
		}
		in.Entries[i] = entry
	}
	reference := Project(in)
	if reference.Status != "allow" || len(reference.Entries) != len(in.Entries) {
		t.Fatalf("invalid reference projection: %+v", reference)
	}
	for i, item := range reference.Entries {
		wantID := fmt.Sprintf("entry-%03d", i)
		wantKind := "observed"
		if i%2 != 0 {
			wantKind = "inferred"
		}
		if item.ID != wantID || item.Kind != wantKind {
			t.Fatalf("unexpected reference chronology: %+v", item)
		}
	}
	rng := rand.New(rand.NewSource(1102026))
	for attempt := 0; attempt < 30; attempt++ {
		entries := append([]Entry(nil), in.Entries...)
		rng.Shuffle(len(entries), func(i, j int) {
			entries[i], entries[j] = entries[j], entries[i]
		})
		before := append([]Entry(nil), entries...)
		got := Project(Input{Scope: in.Scope, Entries: entries})
		if !reflect.DeepEqual(got, reference) {
			t.Fatalf("shuffled order changed projection at attempt %d: %+v", attempt, got)
		}
		if !reflect.DeepEqual(entries, before) {
			t.Fatalf("projection mutated caller-owned entries at attempt %d", attempt)
		}
	}
}

// Hidden source content must not influence visible content. Keep its
// presence/status fixed, and vary only attributes the caller may not read.
func TestProjectionOpaqueSourceContentDoesNotInfluenceVisibleRows(t *testing.T) {
	for _, tc := range []struct {
		name       string
		decision   string
		tenant     string
		wantStatus string
	}{
		{"denied", "deny", "tenant-a", "filtered"},
		{"unknown", "unknown", "tenant-a", "filtered"},
		{"missing", "missing", "tenant-a", "partial"},
		{"foreign", "allow", "tenant-other", "filtered"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := allowedFixture()
			secret := in.Entries[0]
			secret.Tenant, secret.SourcePermission = tc.tenant, tc.decision
			secret.ID = "protected-object"
			in.Entries = append(in.Entries, secret)
			expected := Project(in)
			if expected.Status != tc.wantStatus || len(expected.Entries) != 1 {
				t.Fatalf("unexpected baseline: %+v", expected)
			}
			for _, mutation := range []func(*Entry){
				func(e *Entry) { e.ID = "" },
				func(e *Entry) { e.ID = "x" },
				func(e *Entry) { e.SourceRef = strings.Repeat("private", 1<<20) },
				func(e *Entry) { e.Provenance = "confidential" },
				func(e *Entry) { e.Kind = "invalid"; e.Version = -1 },
				func(e *Entry) { e.SourceTime = "invalid"; e.ClockUncertaintyMS = -4 },
				func(e *Entry) { e.Supersedes = "hidden-predecessor" },
			} {
				copyInput := allowedFixture()
				opaque := secret
				mutation(&opaque)
				copyInput.Entries = append(copyInput.Entries, opaque)
				got := Project(copyInput)
				if !reflect.DeepEqual(got, expected) {
					t.Fatalf("masked source changed authorized output: %+v vs %+v", got, expected)
				}
			}
		})
	}
}

func TestProjectionSortTotalOrderAcrossClockAndTimestampTies(t *testing.T) {
	base := allowedFixture()
	newEntry := func(id, when, zone string, skew int64) Entry {
		entry := base.Entries[0]
		entry.ID = id
		entry.SourceRef = "source:" + id
		entry.Provenance = "provenance:" + id
		entry.SourceTime = when
		entry.Timezone = zone
		entry.ClockUncertaintyMS = skew
		return entry
	}
	// z-known and a-known denote the exact same instant across two UTC
	// offsets. Explicit ID tie-breaking must not rely on input order.
	base.Entries = []Entry{
		newEntry("late", "2026-10-01T12:00:00Z", "+00:00", 0),
		newEntry("z-known", "2026-10-01T10:00:00Z", "+00:00", 0),
		newEntry("clock-uncertain", "2026-10-01T08:00:00Z", "+00:00", 250),
		newEntry("a-known", "2026-10-01T11:00:00+01:00", "+01:00", 0),
		newEntry("unknown-zone", "2026-10-01T07:00:00Z", "unknown", 0),
		newEntry("early", "2026-10-01T09:00:00Z", "+00:00", 0),
	}
	wantIDs := []string{"early", "a-known", "z-known", "late", "clock-uncertain", "unknown-zone"}
	reference := Project(base)
	if reference.Status != "uncertain" || len(reference.Entries) != len(wantIDs) {
		t.Fatalf("unexpected reference status or size: %+v", reference)
	}
	for i, id := range wantIDs {
		if reference.Entries[i].ID != id {
			t.Fatalf("total chronology at %d: got %s want %s", i, reference.Entries[i].ID, id)
		}
		if reference.Entries[i].TimeUncertain != (i >= 4) {
			t.Fatalf("known/uncertain precedence wrong for %s", id)
		}
	}
	rng := rand.New(rand.NewSource(1102027))
	for trial := 0; trial < 60; trial++ {
		shuffled := append([]Entry(nil), base.Entries...)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		original := append([]Entry(nil), shuffled...)
		if got := Project(Input{Scope: base.Scope, Entries: shuffled}); !reflect.DeepEqual(got, reference) {
			t.Fatalf("permutation %d changed total sort order: %+v", trial, got)
		}
		if !reflect.DeepEqual(shuffled, original) {
			t.Fatalf("permutation %d mutated caller snapshots", trial)
		}
	}
	duplicate := base
	duplicate.Entries = append(append([]Entry(nil), base.Entries...), base.Entries[0])
	if got := Project(duplicate); got.Status != "deny" || len(got.Entries) != 0 {
		t.Fatalf("duplicate authorized ID destroyed comparator uniqueness but was accepted: %+v", got)
	}
}
