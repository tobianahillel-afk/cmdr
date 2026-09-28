package savedqueryassets

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildExecutionHandoffPreservesVersionValuesAndContext(t *testing.T) {
	projected := Project(validInput())
	values := []ParameterValue{
		{Name: "hostname", OpaqueValue: []byte("host-17")},
		{Name: "window", OpaqueValue: []byte("15m")},
	}
	got := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod", ParameterValues: values,
	})
	if !got.Allowed || got.Draft == nil {
		t.Fatalf("handoff rejected: %#v", got)
	}
	draft := got.Draft
	if draft.QueryRef != "qry-001" || draft.QueryTenantRef != "tenant-a" || draft.QueryVersion != "v7" {
		t.Fatalf("query identity changed: %#v", draft)
	}
	if string(draft.ParameterValues[0].OpaqueValue) != "host-17" ||
		string(draft.ParameterValues[1].OpaqueValue) != "15m" {
		t.Fatalf("parameter values changed: %#v", draft.ParameterValues)
	}
	if len(draft.Sources) != 1 || draft.Sources[0].ID != "source-edr" ||
		len(draft.Sources[0].RequiredFields) != 2 ||
		draft.Sources[0].RequiredFields[0] != "host.name" ||
		draft.Sources[0].RequiredFields[1] != "event.action" {
		t.Fatalf("source context changed: %#v", draft.Sources)
	}
	if !draft.SourceReevaluationRequired || !draft.PermissionReevaluationRequired ||
		!draft.OriginalAssetImmutable {
		t.Fatalf("handoff re-evaluation/immutability requirements missing: %#v", draft)
	}
	if draft.ReturnContext.AssetID != "asset-001" || draft.ReturnContext.AssetKind != AssetSavedSearch {
		t.Fatalf("return context changed asset identity: %#v", draft.ReturnContext)
	}
	if got.SourceAssetMutated || got.SearchJobCreated || got.SearchJobExecuted {
		t.Fatalf("handoff performed forbidden side effect: %#v", got)
	}
}

func TestBuildExecutionHandoffDeepCopiesProjectionAndParameterValues(t *testing.T) {
	projected := Project(validInput())
	value := []byte("host-17")
	got := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod",
		ParameterValues: []ParameterValue{{Name: "hostname", OpaqueValue: value}},
	})
	if !got.Allowed || got.Draft == nil {
		t.Fatalf("handoff rejected: %#v", got)
	}

	value[0] = 'X'
	projected.Projection.ParameterNames[0] = "mutated"
	projected.Projection.Sources[0].ID = "mutated"
	projected.Projection.Sources[0].Fields[0].Name = "mutated"
	projected.Projection.QueryRef = "mutated"

	draft := got.Draft
	if string(draft.ParameterValues[0].OpaqueValue) != "host-17" ||
		draft.ParameterNames[0] != "hostname" ||
		draft.Sources[0].ID != "source-edr" ||
		draft.Sources[0].RequiredFields[0] != "host.name" ||
		draft.QueryRef != "qry-001" {
		t.Fatalf("handoff aliases caller-owned state: %#v", draft)
	}
}

func TestBuildExecutionHandoffRejectsIneligibleProjection(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Input)
	}{
		{"stale", func(in *Input) { in.ValidationState = ValidationStale }},
		{"incompatible", func(in *Input) { in.Sources[0].Fields[0].State = FieldRemoved }},
		{"asset-access-denied", func(in *Input) { in.AssetAccessDecision = AccessDeny }},
		{"query-access-denied", func(in *Input) { in.QueryAccessDecision = AccessDeny }},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			in := validInput()
			tt.edit(&in)
			projected := Project(in)
			got := BuildExecutionHandoff(projected, HandoffInput{
				TenantRef: "tenant-a", EnvironmentRef: "env-prod",
			})
			if got.Allowed || got.Draft != nil || !got.AuditRequired ||
				got.SearchJobCreated || got.SearchJobExecuted || got.SourceAssetMutated {
				t.Fatalf("ineligible projection produced handoff: %#v", got)
			}
		})
	}
}

func TestBuildExecutionHandoffRejectsScopeMismatch(t *testing.T) {
	projected := Project(validInput())
	for _, input := range []HandoffInput{
		{TenantRef: "tenant-b", EnvironmentRef: "env-prod"},
		{TenantRef: "tenant-a", EnvironmentRef: "env-other"},
		{TenantRef: "", EnvironmentRef: "env-prod"},
	} {
		got := BuildExecutionHandoff(projected, input)
		if got.Allowed || got.ErrorCode != ErrorHandoffScopeMismatch || got.Draft != nil {
			t.Fatalf("scope mismatch accepted: %#v", got)
		}
	}
}

func TestBuildExecutionHandoffRejectsUnknownAndDuplicateParameterBindings(t *testing.T) {
	projected := Project(validInput())
	unknown := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod",
		ParameterValues: []ParameterValue{{Name: "unknown", OpaqueValue: []byte("x")}},
	})
	if unknown.Allowed || unknown.ErrorCode != ErrorHandoffUnknownParameter {
		t.Fatalf("unknown parameter accepted: %#v", unknown)
	}

	duplicate := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod",
		ParameterValues: []ParameterValue{
			{Name: "hostname", OpaqueValue: []byte("a")},
			{Name: "hostname", OpaqueValue: []byte("b")},
		},
	})
	if duplicate.Allowed || duplicate.ErrorCode != ErrorHandoffDuplicateParameter {
		t.Fatalf("duplicate parameter accepted: %#v", duplicate)
	}
}

func TestExecutionHandoffDraftContainsNoSearchJobIdentityOrState(t *testing.T) {
	projected := Project(validInput())
	got := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod",
	})
	if !got.Allowed || got.Draft == nil {
		t.Fatalf("handoff rejected: %#v", got)
	}
	data, err := json.Marshal(got.Draft)
	if err != nil {
		t.Fatal(err)
	}
	lower := strings.ToLower(string(data))
	if strings.Contains(lower, "searchjob") || strings.Contains(lower, "search_job") ||
		strings.Contains(lower, "job_id") || strings.Contains(lower, "job_state") {
		t.Fatalf("handoff draft minted Search Job identity/state: %s", data)
	}
}

func TestExecutionHandoffAllowsUnboundDeclaredParametersWithoutInventingDefaults(t *testing.T) {
	projected := Project(validInput())
	got := BuildExecutionHandoff(projected, HandoffInput{
		TenantRef: "tenant-a", EnvironmentRef: "env-prod",
		ParameterValues: []ParameterValue{{Name: "hostname", OpaqueValue: []byte("host-17")}},
	})
	if !got.Allowed || got.Draft == nil {
		t.Fatalf("handoff rejected: %#v", got)
	}
	if len(got.Draft.ParameterNames) != 2 || len(got.Draft.ParameterValues) != 1 {
		t.Fatalf("handoff synthesized parameter values/defaults: %#v", got.Draft)
	}
}
