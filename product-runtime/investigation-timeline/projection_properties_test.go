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
