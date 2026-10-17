package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type scopedOwnerChecker struct {
	calls      []AccessRequest
	decisions  map[string]string
	failObject string
}

func (s *scopedOwnerChecker) CheckAccess(_ context.Context, r AccessRequest) (string, error) {
	s.calls = append(s.calls, r)
	if r.CaseID != "case-123" || r.Tenant != "tenant-a" || r.Environment != "prod" {
		return "deny", nil
	}
	if r.Resource == "case" && r.ObjectID == "case-123" && r.Permission == caseReadPermission {
		return "allow", nil
	}
	if r.Resource == "timeline-entry" && r.Permission != timelineReadPermission {
		return "deny", nil
	}
	if r.Resource == "source" && r.Permission != sourceReadPermission {
		return "deny", nil
	}
	key := r.Resource + ":" + r.ObjectID
	if key == s.failObject {
		return "", errors.New("source provider unavailable")
	}
	return s.decisions[key], nil
}

func timelineEntryVariant(original Entry, id, source string) Entry {
	e := original
	e.ID = id
	e.SourceRef = source
	e.Provenance = "forensic:" + id
	return e
}

func TestOwnerChecksAcrossMultipleIndependentTimelineSources(t *testing.T) {
	base := allowedFixture()
	allowed := base.Entries[0]
	unrelated := timelineEntryVariant(allowed, "unrelated", "owner:secret")
	missing := timelineEntryVariant(allowed, "missing-source", "owner:missing")
	foreign := timelineEntryVariant(allowed, "other-tenant", "owner:foreign")
	foreign.Tenant = "foreign"
	base.Entries = []Entry{allowed, unrelated, missing, foreign}
	checker := &scopedOwnerChecker{decisions: map[string]string{
		"timeline-entry:x":              "allow",
		"source:owner:x":                "allow",
		"timeline-entry:unrelated":      "deny",
		"timeline-entry:missing-source": "allow",
		"source:owner:missing":          "missing",
	}}
	result := ProjectWithAccessChecker(context.Background(), base, "case-123", checker)
	if result.Status != "partial" || len(result.Entries) != 1 || result.Entries[0].ID != "x" {
		t.Fatalf("unrelated/opaque owner data disclosed: %+v", result)
	}
	want := []AccessRequest{
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "case", ObjectID: "case-123", Permission: caseReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "timeline-entry", ObjectID: "x", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "timeline-entry", ObjectID: "unrelated", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "timeline-entry", ObjectID: "missing-source", Permission: timelineReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "source", ObjectID: "owner:x", Permission: sourceReadPermission},
		{Tenant: "tenant-a", Environment: "prod", CaseID: "case-123", Resource: "source", ObjectID: "owner:missing", Permission: sourceReadPermission},
	}
	if !reflect.DeepEqual(checker.calls, want) {
		t.Fatalf("wrong provider call order, source lookup after denial, or cross-tenant call: got %+v want %+v", checker.calls, want)
	}
}

func TestOwnerProviderFailureIsAtomicAfterEarlierAllowedEntries(t *testing.T) {
	base := allowedFixture()
	base.Entries = append(base.Entries, timelineEntryVariant(base.Entries[0], "second", "owner:second"))
	checker := &scopedOwnerChecker{decisions: map[string]string{
		"timeline-entry:x":      "allow",
		"source:owner:x":        "allow",
		"timeline-entry:second": "allow",
	}, failObject: "source:owner:second"}
	result := ProjectWithAccessChecker(context.Background(), base, "case-123", checker)
	if result.Status != "deny" || len(result.Entries) != 0 {
		t.Fatalf("partial output escaped after failure in second source: %+v", result)
	}
	if got := len(checker.calls); got != 5 {
		t.Fatalf("provider failure did not occur after first complete authorized entry: calls=%+v", checker.calls)
	}
}
