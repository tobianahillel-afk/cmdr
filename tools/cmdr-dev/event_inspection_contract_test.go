package main

import "testing"

func testEventInspectionContract() EventInspectionExecutableContract {
	return EventInspectionExecutableContract{
		PermissionPolicy: EventInspectionPermissionPolicy{
			RequiredProjectionPermissions: []string{"perm.shared-capabilities.telemetry-event.read"},
		},
	}
}

func validEventInspectionInput() EventInspectionFixtureInput {
	return EventInspectionFixtureInput{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "env-prod", EventRef: "evt-001", EventTenantRef: "tenant-a",
		Permissions: []string{"perm.shared-capabilities.telemetry-event.read"},
		RawAccessDecision: "allow", RenderedAccessDecision: "allow",
		RawPayloadPresent: true, RenderedPayloadPresent: true,
	}
}

func TestEventInspectionRawDeniedDoesNotLeakRawDerivedSensitiveValues(t *testing.T) {
	in := validEventInspectionInput()
	in.RawAccessDecision = "deny"
	in.RawDerivedSensitivePresent = true
	got := evaluateEventInspectionFixture(testEventInspectionContract(), in)
	if !got.Allowed || got.RawVisible || got.RawDerivedSensitiveVisible || !got.RenderedVisible {
		t.Fatalf("raw-denied projection leaked or hid rendered data: %#v", got)
	}
}

func TestEventInspectionRejectsCrossTenantProjection(t *testing.T) {
	in := validEventInspectionInput()
	in.TenantRef = "tenant-b"
	got := evaluateEventInspectionFixture(testEventInspectionContract(), in)
	if got.Allowed || got.ErrorCode != "cross-tenant" || got.ProvenanceVisible {
		t.Fatalf("cross-tenant projection was not rejected safely: %#v", got)
	}
}

func TestEventInspectionRejectsDerivedAsSource(t *testing.T) {
	in := validEventInspectionInput()
	in.EnrichmentsPresent = true
	in.EnrichmentProducer = "producer"
	in.EnrichmentVersion = "v1"
	in.EnrichmentFreshness = "fresh"
	in.EnrichmentPresentedAsSource = true
	got := evaluateEventInspectionFixture(testEventInspectionContract(), in)
	if got.Allowed || got.ErrorCode != "derived-as-source" {
		t.Fatalf("derived enrichment masquerading as source was accepted: %#v", got)
	}
}

func TestEventInspectionPivotIsDialectNeutralDraft(t *testing.T) {
	in := validEventInspectionInput()
	in.PivotRequested = true
	in.PivotField = "source.ip"
	in.PivotValuePresent = true
	in.PivotTimeStart = "2026-09-25T10:00:00Z"
	in.PivotTimeEnd = "2026-09-25T11:00:00Z"
	in.ReturnOrigin = "event-search:run-001:row-17"
	got := evaluateEventInspectionFixture(testEventInspectionContract(), in)
	if !got.Allowed || !got.PivotDraftCreated || got.PivotDialectSelected || !got.ReturnContextPreserved {
		t.Fatalf("pivot draft semantics are invalid: %#v", got)
	}
}

func TestEventInspectionMutationsRemainExcluded(t *testing.T) {
	for _, mutation := range []string{"case-link", "artifact-proposal", "evidence-candidate"} {
		in := validEventInspectionInput()
		in.MutationRequested = mutation
		got := evaluateEventInspectionFixture(testEventInspectionContract(), in)
		if got.Allowed || got.ErrorCode != "mutation-excluded" || got.MutationExecuted {
			t.Fatalf("%s mutation escaped read-only contract: %#v", mutation, got)
		}
	}
}
