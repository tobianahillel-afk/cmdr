package savedqueryassets

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"
)

type fixtureExpected struct {
	Allowed             bool     `json:"allowed"`
	ErrorCode           string   `json:"error_code"`
	Compatibility       string   `json:"compatibility"`
	Diagnostics         []string `json:"diagnostics"`
	HandoffEligible     bool     `json:"handoff_eligible"`
	ImmutableProjection bool     `json:"immutable_projection"`
	DeepCopyIsolated    bool     `json:"deep_copy_isolated"`
	MutationExecuted    bool     `json:"mutation_executed"`
	SearchJobCreated    bool     `json:"search_job_created"`
	SearchJobExecuted   bool     `json:"search_job_executed"`
	AuditRequired       bool     `json:"audit_required"`
}

type fixtureSet struct {
	Cases []struct {
		ID       string          `json:"id"`
		Input    Input           `json:"input"`
		Expected fixtureExpected `json:"expected"`
	} `json:"cases"`
}

func TestProjectMatchesPredeclaredContractFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/implementation/saved-query-assets/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures fixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Cases) != 18 {
		t.Fatalf("fixture count=%d want=18", len(fixtures.Cases))
	}
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		t.Run(fixture.ID, func(t *testing.T) {
			got := Project(fixture.Input)
			want := fixture.Expected
			if got.Allowed != want.Allowed || string(got.ErrorCode) != want.ErrorCode ||
				got.Compatibility != want.Compatibility || got.HandoffEligible != want.HandoffEligible ||
				got.ImmutableProjection != want.ImmutableProjection || got.DeepCopyIsolated != want.DeepCopyIsolated ||
				got.MutationExecuted != want.MutationExecuted || got.SearchJobCreated != want.SearchJobCreated ||
				got.SearchJobExecuted != want.SearchJobExecuted || got.AuditRequired != want.AuditRequired ||
				!sameStrings(got.Diagnostics, want.Diagnostics) {
				t.Fatalf("result mismatch:\n got=%#v\nwant=%#v", got, want)
			}
			if !got.Allowed {
				if got.Projection != nil {
					t.Fatalf("denied request leaked projection: %#v", got.Projection)
				}
				return
			}
			if got.Projection == nil {
				t.Fatal("allowed request has no projection")
			}
			if got.Projection.Compatibility != want.Compatibility ||
				got.Projection.HandoffEligible != want.HandoffEligible ||
				!sameStrings(got.Projection.Diagnostics, want.Diagnostics) {
				t.Fatalf("projection compatibility mismatch: %#v", got.Projection)
			}
		})
	}
}

func TestProjectionDeepCopiesCallerOwnedAsset(t *testing.T) {
	in := validInput()
	in.ParameterNames = []string{"hostname", "window"}
	in.Sources = []Source{{
		ID: "source-edr", TenantRef: "tenant-a", State: SourceAvailable,
		Fields: []Field{{Name: "host.name", State: FieldAvailable}},
	}}
	in.LineageRef = "asset-parent-1"
	in.DeprecationReason = "legacy query"
	in.ReplacementAssetRef = "asset-2"

	got := Project(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("valid asset rejected: %#v", got)
	}
	in.ParameterNames[0] = "mutated"
	in.Sources[0].ID = "mutated"
	in.Sources[0].Fields[0].Name = "mutated"

	p := got.Projection
	if p.ParameterNames[0] != "hostname" || p.Sources[0].ID != "source-edr" ||
		p.Sources[0].Fields[0].Name != "host.name" || p.LineageRef != "asset-parent-1" ||
		p.DeprecationReason != "legacy query" || p.ReplacementAssetRef != "asset-2" {
		t.Fatalf("projection aliases caller input: %#v", p)
	}
}

func TestProjectionPreservesQueryVersionAuthorValidationAndMetadata(t *testing.T) {
	in := validInput()
	in.QueryRef = "qry-77"
	in.QueryVersion = "version-19"
	in.AuthorRef = "author-42"
	in.ValidationState = ValidationStale
	in.LineageRef = "parent-7"
	in.DeprecationReason = "superseded"
	in.ReplacementAssetRef = "asset-next"
	in.Sources[0].State = SourceStale

	got := Project(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("projection rejected: %#v", got)
	}
	p := got.Projection
	if p.QueryRef != "qry-77" || p.QueryVersion != "version-19" || p.AuthorRef != "author-42" ||
		p.ValidationState != ValidationStale || p.LineageRef != "parent-7" ||
		p.DeprecationReason != "superseded" || p.ReplacementAssetRef != "asset-next" ||
		p.Compatibility != "stale" || p.HandoffEligible {
		t.Fatalf("identity/provenance lost: %#v", p)
	}
}

func TestProjectionDoesNotSynthesizeOptionalMetadata(t *testing.T) {
	got := Project(validInput())
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("projection rejected: %#v", got)
	}
	p := got.Projection
	if p.LineageRef != "" || p.DeprecationReason != "" || p.ReplacementAssetRef != "" {
		t.Fatalf("missing metadata was synthesized: %#v", p)
	}
}

func TestProjectRejectsAdditionalInvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Input)
		code ErrorCode
	}{
		{"missing-asset-id", func(in *Input) { in.AssetID = "" }, ErrorMissingAssetID},
		{"invalid-asset-kind", func(in *Input) { in.AssetKind = "view" }, ErrorInvalidAssetKind},
		{"invalid-validation-state", func(in *Input) { in.ValidationState = "invented" }, ErrorInvalidValidationState},
		{"invalid-parameter", func(in *Input) { in.ParameterNames = []string{""} }, ErrorInvalidParameter},
		{"missing-sources", func(in *Input) { in.Sources = nil }, ErrorMissingSourcePrerequisites},
		{"invalid-source", func(in *Input) { in.Sources[0].ID = "" }, ErrorInvalidSource},
		{"invalid-source-state", func(in *Input) { in.Sources[0].State = "invented" }, ErrorInvalidSourceState},
		{"invalid-field", func(in *Input) { in.Sources[0].Fields[0].Name = "" }, ErrorInvalidField},
		{"invalid-field-state", func(in *Input) { in.Sources[0].Fields[0].State = "invented" }, ErrorInvalidFieldState},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			in := validInput()
			tt.edit(&in)
			got := Project(in)
			if got.Allowed || got.ErrorCode != tt.code || got.Projection != nil || !got.AuditRequired {
				t.Fatalf("unexpected rejection: %#v", got)
			}
		})
	}
}

func TestIncompatibleOverridesStaleAndDiagnosticsAreSorted(t *testing.T) {
	in := validInput()
	in.ValidationState = ValidationStale
	in.Sources = []Source{
		{
			ID: "z-source", TenantRef: "tenant-a", State: SourceStale,
			Fields: []Field{{Name: "b", State: FieldStale}},
		},
		{
			ID: "a-source", TenantRef: "tenant-a", State: SourceAvailable,
			Fields: []Field{{Name: "a", State: FieldRemoved}},
		},
	}
	got := Project(in)
	if !got.Allowed || got.Compatibility != "incompatible" || got.HandoffEligible {
		t.Fatalf("unexpected compatibility: %#v", got)
	}
	sorted := append([]string(nil), got.Diagnostics...)
	sort.Strings(sorted)
	if !reflect.DeepEqual(got.Diagnostics, sorted) {
		t.Fatalf("diagnostics are not sorted: %v", got.Diagnostics)
	}
}

func validInput() Input {
	return Input{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"}, EnvironmentRef: "env-prod",
		AssetID: "asset-001", AssetKind: AssetSavedSearch, AssetAccessDecision: AccessAllow,
		QueryRef: "qry-001", QueryTenantRef: "tenant-a", QueryVersion: "v7", QueryAccessDecision: AccessAllow,
		AuthorRef: "user-hunter-1", ValidationState: ValidationValidated,
		ParameterNames: []string{"hostname", "window"},
		Sources: []Source{{
			ID: "source-edr", TenantRef: "tenant-a", State: SourceAvailable,
			Fields: []Field{{Name: "host.name", State: FieldAvailable}, {Name: "event.action", State: FieldAvailable}},
		}},
	}
}

func sameStrings(a, b []string) bool {
	left := append([]string(nil), a...)
	right := append([]string(nil), b...)
	sort.Strings(left)
	sort.Strings(right)
	return reflect.DeepEqual(left, right)
}
