package huntmanagement

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
	ProjectedReferences int      `json:"projected_references"`
	Limitations         []string `json:"limitations"`
	ImmutableProjection bool     `json:"immutable_projection"`
	DeepCopyIsolated    bool     `json:"deep_copy_isolated"`
	MutationExecuted    bool     `json:"mutation_executed"`
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
	data, err := os.ReadFile("../../engineering/implementation/hunt-management/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures fixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Cases) != 15 {
		t.Fatalf("fixture count=%d want=15", len(fixtures.Cases))
	}
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		t.Run(fixture.ID, func(t *testing.T) {
			got := Project(fixture.Input)
			if got.Allowed != fixture.Expected.Allowed {
				t.Fatalf("allowed=%t want=%t", got.Allowed, fixture.Expected.Allowed)
			}
			if string(got.ErrorCode) != fixture.Expected.ErrorCode {
				t.Fatalf("error=%q want=%q", got.ErrorCode, fixture.Expected.ErrorCode)
			}
			if got.AuditRequired != fixture.Expected.AuditRequired {
				t.Fatalf("audit=%t want=%t", got.AuditRequired, fixture.Expected.AuditRequired)
			}
			if got.MutationExecuted != fixture.Expected.MutationExecuted {
				t.Fatalf("mutation=%t want=%t", got.MutationExecuted, fixture.Expected.MutationExecuted)
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
			if len(got.Projection.References) != fixture.Expected.ProjectedReferences {
				t.Fatalf("refs=%d want=%d", len(got.Projection.References), fixture.Expected.ProjectedReferences)
			}
			if !sameStrings(got.Projection.Limitations, fixture.Expected.Limitations) {
				t.Fatalf("limitations=%v want=%v", got.Projection.Limitations, fixture.Expected.Limitations)
			}
			if got.Projection.Complete != (len(fixture.Expected.Limitations) == 0) {
				t.Fatalf("complete=%t limitations=%v", got.Projection.Complete, got.Projection.Limitations)
			}
		})
	}
}

func TestProjectionDeepCopiesCallerOwnedWorkspace(t *testing.T) {
	in := validInput()
	in.ContributorRefs = []string{"hunter-2"}
	in.References = []Reference{
		{Kind: ReferenceQuery, ID: "qry-1", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable, Version: "v7", ProvenanceRef: "prov-q"},
		{Kind: ReferenceSearchJob, ID: "job-1", TenantRef: "tenant-a", AccessDecision: AccessAllow, QueryRef: "qry-1", State: ReferenceStale, Version: "run-9", ProvenanceRef: "prov-job"},
	}
	in.CorrelationID = "corr-hunt-1"
	got := Project(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("valid workspace rejected: %#v", got)
	}

	in.ContributorRefs[0] = "mutated"
	in.References[0].ID = "mutated"
	in.References[0].Version = "mutated"
	in.References[1].QueryRef = "mutated"

	p := got.Projection
	if p.ContributorRefs[0] != "hunter-2" ||
		p.References[0].ID != "qry-1" || p.References[0].Version != "v7" ||
		p.References[1].QueryRef != "qry-1" || p.CorrelationID != "corr-hunt-1" {
		t.Fatalf("projection aliases caller input: %#v", p)
	}
}

func TestProjectionPreservesCanonicalOwnershipVersionAndProvenance(t *testing.T) {
	in := validInput()
	in.References = []Reference{
		{Kind: ReferenceQuery, ID: "q", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable, Version: "qv1", ProvenanceRef: "qp"},
		{Kind: ReferenceSearchJob, ID: "j", TenantRef: "tenant-a", AccessDecision: AccessAllow, QueryRef: "q", State: ReferenceAvailable, Version: "jv1", ProvenanceRef: "jp"},
		{Kind: ReferenceHypothesis, ID: "h", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable, Version: "hv1", ProvenanceRef: "hp"},
		{Kind: ReferenceCase, ID: "c", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable, Version: "cv1", ProvenanceRef: "cp"},
		{Kind: ReferenceIncident, ID: "i", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable, Version: "iv1", ProvenanceRef: "ip"},
	}
	got := Project(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("projection rejected: %#v", got)
	}
	wantOwners := []string{"Shared Capabilities", "Shared Capabilities", "Investigate", "Investigate", "Command"}
	for i, ref := range got.Projection.References {
		if ref.CanonicalOwner != wantOwners[i] || ref.Version == "" || ref.ProvenanceRef == "" {
			t.Fatalf("identity/provenance lost at %d: %#v", i, ref)
		}
	}
}

func TestProjectRejectsAdditionalInvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Input)
		code ErrorCode
	}{
		{"blank-contributor", func(in *Input) { in.ContributorRefs = []string{""} }, ErrorInvalidContributor},
		{"invalid-reference-kind", func(in *Input) {
			in.References = []Reference{{Kind: "unknown", ID: "x", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable}}
		}, ErrorInvalidReference},
		{"empty-reference-id", func(in *Input) {
			in.References = []Reference{{Kind: ReferenceCase, ID: "", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable}}
		}, ErrorInvalidReference},
		{"invalid-reference-state", func(in *Input) {
			in.References = []Reference{{Kind: ReferenceCase, ID: "c", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: "invented"}}
		}, ErrorInvalidReferenceState},
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

func TestProjectionDoesNotSynthesizeOptionalCorrelationOrProvenance(t *testing.T) {
	in := validInput()
	in.References = []Reference{{Kind: ReferenceCase, ID: "c", TenantRef: "tenant-a", AccessDecision: AccessAllow, State: ReferenceAvailable}}
	got := Project(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("projection rejected: %#v", got)
	}
	if got.Projection.CorrelationID != "" || got.Projection.References[0].Version != "" || got.Projection.References[0].ProvenanceRef != "" {
		t.Fatalf("missing facts were synthesized: %#v", got.Projection)
	}
}

func validInput() Input {
	return Input{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "env-prod", Question: "Investigate anomalous authentication activity",
		Scope: "identity/authentication", TimeStart: "2026-09-28T08:00:00Z", TimeEnd: "2026-09-28T10:00:00Z",
		OwnerRef: "user-hunter-1",
	}
}

func sameStrings(a, b []string) bool {
	left := append([]string(nil), a...)
	right := append([]string(nil), b...)
	sort.Strings(left)
	sort.Strings(right)
	return reflect.DeepEqual(left, right)
}
