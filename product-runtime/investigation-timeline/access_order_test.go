package investigationtimeline

import (
	"context"
	"reflect"
	"testing"
)

// A requester may choose a different order of Timeline Entry IDs. Both
// server-boundary adapters must still return the same source-backed chronology;
// no source authorization decision is inherited from the caller's flags.
func TestAuthorizedTimelineOrderIgnoresRequestPermutation(t *testing.T) {
	input := allowedFixture()
	additional := input.Entries[0]
	additional.ID = "a"
	additional.Kind = "inferred"
	additional.SourceRef = "owner:a"
	additional.Provenance = "trusted:inference-a"
	additional.SourceTime = "2026-10-01T11:00:00+01:00" // same instant as x
	additional.Timezone = "+01:00"
	additional.SourcePermission = "deny" // owner/source checker must override
	input.Entries = append(input.Entries, additional)
	input.Scope.CasePermission = "deny"
	input.Scope.TimelinePermission = "unknown"

	var ownerReference, legacyReference Result
	for n, ids := range [][]string{{"x", "a"}, {"a", "x"}} {
		checker, owner := ownerFixture()
		owner.entries["a"] = additional
		checker.decisions["timeline-entry:a"] = "allow"
		checker.decisions["source:owner:a"] = "allow"

		ownerResult := ProjectFromOwnerSnapshots(
			context.Background(), input.Scope, "case-123", ids, checker, owner,
		)
		if ownerResult.Status != "allow" || len(ownerResult.Entries) != 2 ||
			ownerResult.Entries[0].ID != "a" ||
			ownerResult.Entries[0].Kind != "inferred" ||
			ownerResult.Entries[1].ID != "x" ||
			ownerResult.Entries[1].Kind != "observed" {
			t.Fatalf("owner-resolved chronological order lost for %v: %+v", ids, ownerResult)
		}
		if len(checker.calls) != 5 || len(owner.calls) != 2 {
			t.Fatalf("not all owner/source rights checked for %v: checker=%+v resolver=%+v", ids, checker.calls, owner.calls)
		}
		if n > 0 && !reflect.DeepEqual(ownerReference, ownerResult) {
			t.Fatalf("owner-authorized result depends on requested ID order: %v vs %v", ownerReference, ownerResult)
		}
		ownerReference = ownerResult

		legacy := input
		legacy.Entries = make([]Entry, 0, len(ids))
		for _, id := range ids {
			if id == "a" {
				legacy.Entries = append(legacy.Entries, additional)
			} else {
				legacy.Entries = append(legacy.Entries, input.Entries[0])
			}
		}
		legacyChecker := allowedAccessChecker()
		legacyChecker.decisions["timeline-entry:a"] = "allow"
		legacyChecker.decisions["source:owner:a"] = "allow"
		legacyResult := ProjectWithAccessChecker(context.Background(), legacy, "case-123", legacyChecker)
		if legacyResult.Status != "allow" || len(legacyResult.Entries) != 2 ||
			legacyResult.Entries[0].ID != "a" || legacyResult.Entries[1].ID != "x" {
			t.Fatalf("legacy adapter chronological order lost for %v: %+v", ids, legacyResult)
		}
		if len(legacyChecker.calls) != 5 {
			t.Fatalf("legacy adapter skipped a provider check for %v: %+v", ids, legacyChecker.calls)
		}
		if n > 0 && !reflect.DeepEqual(legacyReference, legacyResult) {
			t.Fatalf("legacy adapter result depends on request ordering: %v vs %v", legacyReference, legacyResult)
		}
		legacyReference = legacyResult
	}
}
