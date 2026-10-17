package investigationtimeline

import (
	"context"
	"testing"
)

// A failed late Timeline Entry authorization is not allowed to cause earlier
// SourceRef metadata to be dispatched to an external source ACL provider.
func TestLegacyChecksAllTimelineRightsBeforeSourceLookups(t *testing.T) {
	in := allowedFixture()
	second := timelineEntryVariant(in.Entries[0], "later", "owner:later")
	in.Entries = append(in.Entries, second)
	checker := allowedAccessChecker()
	checker.fail = "timeline-entry:later"
	out := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if out.Status != "deny" || len(out.Entries) != 0 || len(checker.calls) != 3 {
		t.Fatalf("late Timeline provider failure was not atomic: %+v calls=%+v", out, checker.calls)
	}
	for _, req := range checker.calls {
		if req.Resource == "source" {
			t.Fatalf("source selector disclosed before complete Timeline authorization: %+v", req)
		}
	}
}

// Prevalidation must cover the last allowed snapshot BEFORE any source
// lookup, not just the first one. Denied entries remain opaque as before.
func TestLegacyPrevalidatesAllEntriesBeforeSourceIO(t *testing.T) {
	for _, edit := range []struct {
		name   string
		mutate func(*Entry)
	}{
		{"invalid-kind", func(e *Entry) { e.Kind = "unrecognized" }},
		{"missing-provenance", func(e *Entry) { e.Provenance = "" }},
		{"invalid-source", func(e *Entry) { e.SourceRef = " * " }},
		{"invalid-version", func(e *Entry) { e.Version = 0 }},
	} {
		t.Run(edit.name, func(t *testing.T) {
			in := allowedFixture()
			second := timelineEntryVariant(in.Entries[0], "later", "owner:later")
			edit.mutate(&second)
			in.Entries = append(in.Entries, second)
			checker := allowedAccessChecker()
			checker.decisions["timeline-entry:later"] = "allow"
			checker.decisions["source:owner:later"] = "allow"
			out := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
			if out.Status != "deny" || len(out.Entries) != 0 || len(checker.calls) != 3 {
				t.Fatalf("invalid late entry reached source provider: %+v calls=%+v", out, checker.calls)
			}
			for _, req := range checker.calls {
				if req.Resource == "source" {
					t.Fatalf("source ACL received metadata despite rejected snapshot: %+v", req)
				}
			}
		})
	}
}

// A successful two-entry projection must still check all independent rights.
func TestLegacyThreePhaseAccessOrder(t *testing.T) {
	in := allowedFixture()
	in.Entries = append(in.Entries, timelineEntryVariant(in.Entries[0], "later", "owner:later"))
	checker := allowedAccessChecker()
	checker.decisions["timeline-entry:later"] = "allow"
	checker.decisions["source:owner:later"] = "allow"
	out := ProjectWithAccessChecker(context.Background(), in, "case-123", checker)
	if out.Status != "allow" || len(out.Entries) != 2 || len(checker.calls) != 5 {
		t.Fatalf("trusted three-phase projection changed outputs: %+v calls=%+v", out, checker.calls)
	}
	for i, kind := range []string{"case", "timeline-entry", "timeline-entry", "source", "source"} {
		if checker.calls[i].Resource != kind {
			t.Fatalf("provider operation %d got %s want %s", i, checker.calls[i].Resource, kind)
		}
	}
}
