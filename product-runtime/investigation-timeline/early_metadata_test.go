package investigationtimeline

import (
	"context"
	"strings"
	"testing"
)

func TestMalformedVisibleForensicMetadataStopsBeforeSourceAuthority(t *testing.T) {
	bad := string([]byte{0xff})
	cases := []struct {
		name   string
		mutate func(*Entry)
	}{
		{"source-time", func(e *Entry) { e.SourceTime += bad }},
		{"timezone", func(e *Entry) { e.Timezone += bad }},
		{"provenance", func(e *Entry) { e.Provenance += bad }},
		{"supersedes", func(e *Entry) { e.Supersedes = "prior-" + bad }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The legacy adapter must not call the external source provider
			// when the Timeline Entry metadata itself is already invalid.
			input := allowedFixture()
			tc.mutate(&input.Entries[0])
			checker := allowedAccessChecker()
			got := ProjectWithAccessChecker(context.Background(), input, "case-123", checker)
			if got.Status != "deny" || len(got.Entries) != 0 || len(checker.calls) != 2 {
				t.Fatalf("legacy adapter called source authority or leaked entries: %+v calls=%+v", got, checker.calls)
			}
			if checker.calls[0].Resource != "case" || checker.calls[1].Resource != "timeline-entry" {
				t.Fatalf("legacy adapter unexpectedly called source provider: %+v", checker.calls)
			}

			// The authoritative owner-resolver adapter must also reject
			// before seeking permission from any source owner.
			ownerChecker, owner := ownerFixture()
			entry := owner.entries["x"]
			tc.mutate(&entry)
			owner.entries["x"] = entry
			got = ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
				"case-123", []string{"x"}, ownerChecker, owner)
			if got.Status != "deny" || len(got.Entries) != 0 ||
				len(ownerChecker.calls) != 2 || len(owner.calls) != 1 {
				t.Fatalf("owner adapter called source authority or leaked entries: %+v checker=%+v resolver=%+v",
					got, ownerChecker.calls, owner.calls)
			}
			if ownerChecker.calls[0].Resource != "case" || ownerChecker.calls[1].Resource != "timeline-entry" {
				t.Fatalf("owner adapter unexpectedly called source provider: %+v", ownerChecker.calls)
			}
		})
	}
}

func TestOpaqueMalformedMetadataRemainsUninspectedBeforeSourceLookup(t *testing.T) {
	bad := string([]byte{0xff})
	for _, decision := range []struct {
		value  string
		status string
	}{{"deny", "filtered"}, {"unknown", "filtered"}, {"missing", "partial"}} {
		t.Run(decision.value, func(t *testing.T) {
			in := allowedFixture()
			in.Entries[0].SourceTime = bad
			in.Entries[0].Timezone = bad
			in.Entries[0].Provenance = bad
			in.Entries[0].Supersedes = bad
			checker := allowedAccessChecker()
			checker.decisions["timeline-entry:x"] = decision.value
			got := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if got.Status != decision.status || len(got.Entries) != 0 || len(checker.calls) != 2 {
				t.Fatalf("opaque metadata affected legacy projection: %+v calls=%+v", got, checker.calls)
			}
			ownerChecker, owner := ownerFixture()
			ownerChecker.decisions["timeline-entry:x"] = decision.value
			entry := owner.entries["x"]
			entry.SourceTime = bad
			owner.entries["x"] = entry
			result := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope,
				"case-123", []string{"x"}, ownerChecker, owner)
			if result.Status != decision.status || len(result.Entries) != 0 ||
				len(owner.calls) != 0 || len(ownerChecker.calls) != 2 {
				t.Fatalf("denied owner metadata was resolved: %+v owner=%+v checker=%+v", result, owner.calls, ownerChecker.calls)
			}
		})
	}
}

func TestPreSourceLookupAcceptsValidMultiScriptForensicMetadata(t *testing.T) {
	input := allowedFixture()
	input.Entries[0].Provenance = "témoins / ראיות / أدلة"
	input.Entries[0].SourceTime = "clock-" + strings.Repeat("?", 3)
	if !validTimelineEntryBeforeSourceLookup(input.Entries[0]) {
		t.Fatal("valid Unicode or unknown-clock metadata was rejected")
	}
	checker := allowedAccessChecker()
	result := ProjectWithAccessChecker(context.Background(), input, "case-123", checker)
	if result.Status != "uncertain" || len(result.Entries) != 1 ||
		result.Entries[0].Provenance != input.Entries[0].Provenance ||
		result.Entries[0].SourceTime != input.Entries[0].SourceTime {
		t.Fatalf("legitimate forensic metadata was altered: %+v", result)
	}
}
