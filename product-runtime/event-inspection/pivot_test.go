package eventinspection

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type pivotFixtureSet struct {
	Cases []struct {
		ID string `json:"id"`
		Input struct {
			TenantRef              string   `json:"tenant_ref"`
			AuthorizedTenants      []string `json:"authorized_tenants"`
			EnvironmentRef         string   `json:"environment_ref"`
			EventRef               string   `json:"event_ref"`
			EventTenantRef         string   `json:"event_tenant_ref"`
			Permissions            []string `json:"permissions"`
			RawAccessDecision      string   `json:"raw_access_decision"`
			RenderedAccessDecision string   `json:"rendered_access_decision"`
			PivotRequested         bool     `json:"pivot_requested"`
			PivotField             string   `json:"pivot_field"`
			PivotValuePresent      bool     `json:"pivot_value_present"`
			PivotTimeStart         string   `json:"pivot_time_start"`
			PivotTimeEnd           string   `json:"pivot_time_end"`
			ReturnOrigin           string   `json:"return_origin"`
		} `json:"input"`
		Expected struct {
			Allowed                bool   `json:"allowed"`
			ErrorCode              string `json:"error_code"`
			PivotDraftCreated      bool   `json:"pivot_draft_created"`
			PivotDialectSelected   bool   `json:"pivot_dialect_selected"`
			ReturnContextPreserved bool   `json:"return_context_preserved"`
			AuditRequired          bool   `json:"audit_required"`
		} `json:"expected"`
	} `json:"cases"`
}

func TestPreparePivotMatchesPredeclaredPivotFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/implementation/event-inspection/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures pivotFixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}

	consumed := 0
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		if !fixture.Input.PivotRequested {
			continue
		}
		consumed++
		t.Run(fixture.ID, func(t *testing.T) {
			in := fixtureInput(
				fixture.Input.TenantRef, fixture.Input.AuthorizedTenants,
				fixture.Input.EnvironmentRef, fixture.Input.EventRef, fixture.Input.EventTenantRef,
				fixture.Input.Permissions, AccessDecision(fixture.Input.RawAccessDecision),
				AccessDecision(fixture.Input.RenderedAccessDecision),
			)
			in.Event.SourceFields = map[string]string{fixture.Input.PivotField: "192.0.2.10"}
			projected := Project(in)
			if !projected.Allowed || projected.Projection == nil {
				t.Fatalf("inspection prerequisite rejected: %#v", projected)
			}
			got := PreparePivot(projected.Projection, PivotSelection{
				Field: fixture.Input.PivotField, Value: "192.0.2.10",
				ValuePresent: fixture.Input.PivotValuePresent, ValueAuthorized: true,
				TimeStart: fixture.Input.PivotTimeStart, TimeEnd: fixture.Input.PivotTimeEnd,
				ReturnOrigin: fixture.Input.ReturnOrigin,
			})
			if got.Allowed != fixture.Expected.Allowed || string(got.ErrorCode) != fixture.Expected.ErrorCode {
				t.Fatalf("pivot result mismatch: got %#v expected allowed=%t error=%q", got, fixture.Expected.Allowed, fixture.Expected.ErrorCode)
			}
			if got.AuditRequired != fixture.Expected.AuditRequired {
				t.Fatalf("audit mismatch: got %t want %t", got.AuditRequired, fixture.Expected.AuditRequired)
			}
			if (got.Draft != nil) != fixture.Expected.PivotDraftCreated {
				t.Fatalf("draft presence mismatch: %#v", got)
			}
			if fixture.Expected.PivotDialectSelected {
				t.Fatal("bounded runtime must never select a pivot dialect")
			}
			if got.Draft != nil {
				assertPivotContext(t, got.Draft, projected.Projection, fixture.Input.ReturnOrigin)
			}
		})
	}
	if consumed != 2 {
		t.Fatalf("consumed %d pivot fixtures, want 2", consumed)
	}
}

func TestPreparePivotRejectsRestrictedValueBeforeLookup(t *testing.T) {
	projected := mustProject(t)
	got := PreparePivot(projected, PivotSelection{
		Field: "not.present", Value: "classified", ValuePresent: true, ValueAuthorized: false,
		TimeStart: "2026-09-25T10:00:00Z", TimeEnd: "2026-09-25T11:00:00Z",
		ReturnOrigin: "event-search:run-001:row-17",
	})
	if got.Allowed || got.ErrorCode != ErrorRestrictedPivotValue || got.Draft != nil || !got.AuditRequired {
		t.Fatalf("restricted value did not fail closed: %#v", got)
	}
}

func TestPreparePivotRequiresValueFromVisibleProjection(t *testing.T) {
	projected := mustProject(t)
	got := PreparePivot(projected, PivotSelection{
		Field: "source.ip", Value: "198.51.100.77", ValuePresent: true, ValueAuthorized: true,
		TimeStart: "2026-09-25T10:00:00Z", TimeEnd: "2026-09-25T11:00:00Z",
		ReturnOrigin: "event-search:run-001:row-17",
	})
	if got.Allowed || got.ErrorCode != ErrorInvalidPivotContext || got.Draft != nil {
		t.Fatalf("forged field/value entered pivot: %#v", got)
	}
}

func TestPreparePivotIsDeterministicAndDialectNeutral(t *testing.T) {
	projected := mustProject(t)
	selection := PivotSelection{
		Field: "source.ip", Value: "192.0.2.10", ValuePresent: true, ValueAuthorized: true,
		TimeStart: "2026-09-25T12:00:00+02:00", TimeEnd: "2026-09-25T13:00:00+02:00",
		ReturnOrigin: "event-search:run-001:row-17",
	}
	first := PreparePivot(projected, selection)
	second := PreparePivot(projected, selection)
	if !reflect.DeepEqual(first, second) || first.Draft == nil {
		t.Fatalf("pivot draft is not deterministic: first=%#v second=%#v", first, second)
	}
	if first.Draft.TimeRange.Start != "2026-09-25T10:00:00Z" ||
		first.Draft.TimeRange.End != "2026-09-25T11:00:00Z" {
		t.Fatalf("time range was not deterministically normalized: %#v", first.Draft.TimeRange)
	}
	if first.Draft.TargetCapability != PivotTargetCapability || first.Draft.FieldOrigin != "source" {
		t.Fatalf("unexpected bounded target/provenance: %#v", first.Draft)
	}
}

func TestPreparePivotPreservesEnrichmentProvenance(t *testing.T) {
	in := validInput()
	in.Event.Enrichments = []EventEnrichment{{
		Producer: "intel-engine", Version: "2026.09", Freshness: "fresh",
		Values: map[string]string{"threat.score": "95"},
	}}
	projected := Project(in)
	if !projected.Allowed || projected.Projection == nil {
		t.Fatalf("inspection prerequisite rejected: %#v", projected)
	}
	got := PreparePivot(projected.Projection, PivotSelection{
		Field: "threat.score", Value: "95", ValuePresent: true, ValueAuthorized: true,
		TimeStart: "2026-09-25T10:00:00Z", TimeEnd: "2026-09-25T11:00:00Z",
		ReturnOrigin: "event-search:run-001:row-17",
	})
	if !got.Allowed || got.Draft == nil ||
		got.Draft.FieldOrigin != "enrichment:intel-engine@2026.09" {
		t.Fatalf("enrichment provenance lost: %#v", got)
	}
}

func mustProject(t *testing.T) *Projection {
	t.Helper()
	got := Project(validInput())
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("valid projection rejected: %#v", got)
	}
	return got.Projection
}

func assertPivotContext(t *testing.T, draft *PivotDraft, projection *Projection, origin string) {
	t.Helper()
	if draft.TargetCapability != PivotTargetCapability ||
		draft.TenantRef != projection.TenantRef ||
		draft.EnvironmentRef != projection.EnvironmentRef ||
		draft.EventRef != projection.EventRef ||
		draft.CorrelationID != projection.CorrelationID ||
		draft.SourceRef != projection.Source.SourceRef ||
		draft.DataSourceRef != projection.Source.DataSourceRef ||
		draft.ReturnContext.TenantRef != projection.TenantRef ||
		draft.ReturnContext.EnvironmentRef != projection.EnvironmentRef ||
		draft.ReturnContext.EventRef != projection.EventRef ||
		draft.ReturnContext.CorrelationID != projection.CorrelationID ||
		draft.ReturnContext.Origin != origin {
		t.Fatalf("pivot return context not preserved: %#v", draft)
	}
}
