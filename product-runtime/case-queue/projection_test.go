package casequeue

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type fixtureExpected struct {
	Allowed                    bool     `json:"allowed"`
	ErrorCode                  string   `json:"error_code"`
	States                     []string `json:"states"`
	Diagnostics                []string `json:"diagnostics"`
	ProjectedCaseIDs           []string `json:"projected_case_ids"`
	ProjectedCount             int      `json:"projected_count"`
	PermissionFilteredCases    int      `json:"permission_filtered_cases"`
	PermissionFilteredElements int      `json:"permission_filtered_elements"`
	AppliedSavedViewFields     []string `json:"applied_saved_view_fields"`
	AppliedSort                []string `json:"applied_sort"`
	ImmutableProjection        bool     `json:"immutable_projection"`
	DeepCopyIsolated           bool     `json:"deep_copy_isolated"`
	MutationExecuted           bool     `json:"mutation_executed"`
	ExportJobCreated           bool     `json:"export_job_created"`
	WorkQueueAggregated        bool     `json:"work_queue_aggregated"`
	AuditRequired              bool     `json:"audit_required"`
}

type fixtureSet struct {
	Cases []struct {
		ID       string          `json:"id"`
		Input    Input           `json:"input"`
		Expected fixtureExpected `json:"expected"`
	} `json:"cases"`
}

func TestProjectMatchesPredeclaredContractFixtures(t *testing.T) {
	data, err := os.ReadFile("../../engineering/implementation/case-queue/fixtures.json") // #nosec G304 -- fixed repository fixture path.
	if err != nil {
		t.Fatal(err)
	}
	var fixtures fixtureSet
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	if len(fixtures.Cases) != 24 {
		t.Fatalf("fixture count=%d want=24", len(fixtures.Cases))
	}
	for _, fixture := range fixtures.Cases {
		fixture := fixture
		t.Run(fixture.ID, func(t *testing.T) {
			got := Project(fixture.Input)
			assertFixtureResult(t, got, fixture.Expected)
		})
	}
}

func assertFixtureResult(t *testing.T, got Result, want fixtureExpected) {
	t.Helper()
	if got.Allowed != want.Allowed || string(got.ErrorCode) != want.ErrorCode ||
		got.AuditRequired != want.AuditRequired || got.MutationExecuted != want.MutationExecuted ||
		got.ExportJobCreated != want.ExportJobCreated || got.WorkQueueAggregated != want.WorkQueueAggregated {
		t.Fatalf("result envelope mismatch: got=%#v want=%#v", got, want)
	}
	if !got.Allowed {
		if got.Projection != nil {
			t.Fatalf("rejected result leaked projection: %#v", got.Projection)
		}
		return
	}
	if got.Projection == nil {
		t.Fatal("allowed result has nil projection")
	}
	p := got.Projection
	ids := make([]string, 0, len(p.Rows))
	for _, row := range p.Rows {
		ids = append(ids, row.Case.ID)
	}
	sortStrings := make([]string, 0, len(p.AppliedSort))
	for _, spec := range p.AppliedSort {
		sortStrings = append(sortStrings, spec.Field+":"+spec.Direction)
	}
	if !reflect.DeepEqual(p.States, want.States) ||
		!reflect.DeepEqual(p.Diagnostics, want.Diagnostics) ||
		!reflect.DeepEqual(ids, want.ProjectedCaseIDs) ||
		len(p.Rows) != want.ProjectedCount ||
		p.PermissionFilteredCases != want.PermissionFilteredCases ||
		p.PermissionFilteredElements != want.PermissionFilteredElements ||
		!reflect.DeepEqual(p.AppliedSavedViewFields, want.AppliedSavedViewFields) ||
		!reflect.DeepEqual(sortStrings, want.AppliedSort) ||
		!want.ImmutableProjection || !want.DeepCopyIsolated {
		t.Fatalf("projection mismatch: got=%#v want=%#v", p, want)
	}
}

func TestProjectionDeepCopyIsolation(t *testing.T) {
	in := baseInput()
	in.Cases[0].Incident = &IncidentContext{ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "allow", State: "available", Priority: "P1"}
	in.Cases[0].Findings = []FindingContext{{ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a", AccessDecision: "allow", State: "available", Status: "confirmed"}}
	got := Project(in)
	if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 1 {
		t.Fatalf("unexpected projection: %#v", got)
	}
	in.Cases[0].Title = "MUTATED"
	in.Cases[0].Incident.Priority = "P9"
	in.Cases[0].Findings[0].Status = "mutated"
	if got.Projection.Rows[0].Case.Title == "MUTATED" ||
		got.Projection.Rows[0].Case.Incident.Priority == "P9" ||
		got.Projection.Rows[0].Case.Findings[0].Status == "mutated" {
		t.Fatal("projection aliases caller-owned input")
	}
	got.Projection.Rows[0].Case.Title = "OUTPUT-MUTATION"
	if in.Cases[0].Title == "OUTPUT-MUTATION" {
		t.Fatal("caller input aliases projected output")
	}
}

func TestProjectRejectsMalformedSnapshots(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*Input)
		want ErrorCode
	}{
		{"unauthorized-tenant", func(in *Input) { in.AuthorizedTenants = []string{"tenant-b"} }, ErrorCrossTenant},
		{"duplicate-case", func(in *Input) { in.Cases = append(in.Cases, in.Cases[0]) }, ErrorDuplicateCaseID},
		{"invalid-case-access", func(in *Input) { in.Cases[0].AccessDecision = "maybe" }, ErrorInvalidCaseAccessDecision},
		{"incomplete-case", func(in *Input) { in.Cases[0].Title = "" }, ErrorIncompleteCaseSnapshot},
		{"invalid-freshness", func(in *Input) { in.Cases[0].Freshness = "future" }, ErrorInvalidCaseFreshness},
		{"invalid-incident-access", func(in *Input) {
			in.Cases[0].Incident = &IncidentContext{ID: "i", TenantRef: "tenant-a", AccessDecision: "maybe", State: "available"}
		}, ErrorInvalidIncidentAccessDecision},
		{"invalid-incident-state", func(in *Input) {
			in.Cases[0].Incident = &IncidentContext{ID: "i", TenantRef: "tenant-a", AccessDecision: "allow", State: "future"}
		}, ErrorInvalidIncidentState},
		{"finding-case-mismatch", func(in *Input) {
			in.Cases[0].Findings = []FindingContext{{ID: "f", TenantRef: "tenant-a", CaseID: "other", AccessDecision: "allow", State: "available"}}
		}, ErrorFindingCaseMismatch},
		{"invalid-finding-access", func(in *Input) {
			in.Cases[0].Findings = []FindingContext{{ID: "f", TenantRef: "tenant-a", CaseID: "case-a", AccessDecision: "maybe", State: "available"}}
		}, ErrorInvalidFindingAccessDecision},
		{"invalid-finding-state", func(in *Input) {
			in.Cases[0].Findings = []FindingContext{{ID: "f", TenantRef: "tenant-a", CaseID: "case-a", AccessDecision: "allow", State: "future"}}
		}, ErrorInvalidFindingState},
		{"invalid-activity", func(in *Input) { in.Cases[0].ActivityState = "future" }, ErrorInvalidActivityState},
		{"empty-filter-values", func(in *Input) { in.Filters = []Filter{{Field: "owner"}} }, ErrorEmptyFilterValue},
		{"blank-filter-value", func(in *Input) { in.Filters = []Filter{{Field: "owner", Values: []string{" "}}} }, ErrorEmptyFilterValue},
		{"invalid-sort-field", func(in *Input) { in.Sort = []SortSpec{{Field: "secret", Direction: "asc"}} }, ErrorInvalidSortField},
		{"duplicate-sort", func(in *Input) { in.Sort = []SortSpec{{Field: "owner", Direction: "asc"}, {Field: "owner", Direction: "desc"}} }, ErrorDuplicateSortField},
		{"missing-view-id", func(in *Input) { in.SavedView = &SavedView{TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow"} }, ErrorMissingSavedViewID},
		{"invalid-view-access", func(in *Input) { in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "maybe"} }, ErrorInvalidSavedViewAccessDecision},
		{"invalid-view-field-access", func(in *Input) {
			in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow", Fields: []SavedViewField{{Name: "title", AccessDecision: "maybe"}}}
		}, ErrorInvalidSavedViewElementAccess},
		{"blank-view-field", func(in *Input) {
			in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow", Fields: []SavedViewField{{Name: " ", AccessDecision: "allow"}}}
		}, ErrorInvalidSavedViewField},
		{"duplicate-view-field", func(in *Input) {
			in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow", Fields: []SavedViewField{{Name: "title", AccessDecision: "allow"}, {Name: "title", AccessDecision: "allow"}}}
		}, ErrorDuplicateSavedViewField},
		{"invalid-view-filter-access", func(in *Input) {
			in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow", Filters: []SavedViewFilter{{Field: "owner", Values: []string{"alice"}, AccessDecision: "maybe"}}}
		}, ErrorInvalidSavedViewElementAccess},
		{"invalid-view-sort-access", func(in *Input) {
			in.SavedView = &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow", Sort: []SavedViewSort{{Field: "owner", Direction: "asc", AccessDecision: "maybe"}}}
		}, ErrorInvalidSavedViewElementAccess},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := baseInput()
			tc.mut(&in)
			got := Project(in)
			if got.Allowed || got.ErrorCode != tc.want || got.Projection != nil {
				t.Fatalf("got=%#v want error=%s", got, tc.want)
			}
		})
	}
}

func TestProjectionContextStatesAndPermissionFiltering(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Input)
		wantStates []string
		wantDiag   string
	}{
		{
			name: "incident-denied",
			mutate: func(in *Input) {
				in.Cases[0].Incident = &IncidentContext{
					ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "deny", State: "available",
				}
			},
			wantStates: []string{"available", "permission-filtered"},
			wantDiag:   "incident-context:permission-filtered",
		},
		{
			name: "incident-stale",
			mutate: func(in *Input) {
				in.Cases[0].Incident = &IncidentContext{
					ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "allow", State: "stale", Priority: "P1",
				}
			},
			wantStates: []string{"available", "stale"},
			wantDiag:   "incident-context:case-a:stale",
		},
		{
			name: "incident-unavailable",
			mutate: func(in *Input) {
				in.Cases[0].Incident = &IncidentContext{
					ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "allow", State: "unavailable",
				}
			},
			wantStates: []string{"available", "partial"},
			wantDiag:   "incident-context:case-a:unavailable",
		},
		{
			name: "finding-denied",
			mutate: func(in *Input) {
				in.Cases[0].Findings = []FindingContext{{
					ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a",
					AccessDecision: "deny", State: "available", Status: "confirmed",
				}}
			},
			wantStates: []string{"available", "permission-filtered"},
			wantDiag:   "finding-context:permission-filtered",
		},
		{
			name: "finding-partial",
			mutate: func(in *Input) {
				in.Cases[0].Findings = []FindingContext{{
					ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a",
					AccessDecision: "allow", State: "partial", Status: "confirmed",
				}}
			},
			wantStates: []string{"available", "partial"},
			wantDiag:   "finding-context:case-a:partial",
		},
		{
			name: "finding-stale",
			mutate: func(in *Input) {
				in.Cases[0].Findings = []FindingContext{{
					ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a",
					AccessDecision: "allow", State: "stale", Status: "confirmed",
				}}
			},
			wantStates: []string{"available", "stale"},
			wantDiag:   "finding-context:case-a:stale",
		},
		{
			name: "finding-unavailable",
			mutate: func(in *Input) {
				in.Cases[0].Findings = []FindingContext{{
					ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a",
					AccessDecision: "allow", State: "unavailable", Status: "confirmed",
				}}
			},
			wantStates: []string{"available", "partial"},
			wantDiag:   "finding-context:case-a:unavailable",
		},
		{
			name: "activity-partial",
			mutate: func(in *Input) {
				in.Cases[0].ActivityState = "partial"
			},
			wantStates: []string{"available", "partial"},
			wantDiag:   "activity:case-a:partial",
		},
		{
			name: "activity-unavailable",
			mutate: func(in *Input) {
				in.Cases[0].ActivityState = "unavailable"
			},
			wantStates: []string{"available", "partial"},
			wantDiag:   "activity:case-a:unavailable",
		},
		{
			name: "activity-stale",
			mutate: func(in *Input) {
				in.Cases[0].ActivityState = "stale"
			},
			wantStates: []string{"available", "stale"},
			wantDiag:   "activity:case-a:stale",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := baseInput()
			tc.mutate(&in)
			got := Project(in)
			if !got.Allowed || got.Projection == nil {
				t.Fatalf("unexpected rejection: %#v", got)
			}
			if !reflect.DeepEqual(got.Projection.States, tc.wantStates) {
				t.Fatalf("states=%v want=%v", got.Projection.States, tc.wantStates)
			}
			if !contains(got.Projection.Diagnostics, tc.wantDiag) {
				t.Fatalf("diagnostics=%v missing %q", got.Projection.Diagnostics, tc.wantDiag)
			}
		})
	}
}

func TestProjectionSavedViewElementValidationBranches(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SavedView)
		want   ErrorCode
	}{
		{
			name: "invalid-view-filter-field",
			mutate: func(v *SavedView) {
				v.Filters = []SavedViewFilter{{Field: "secret", Values: []string{"x"}, AccessDecision: "allow"}}
			},
			want: ErrorInvalidFilterField,
		},
		{
			name: "empty-view-filter-value",
			mutate: func(v *SavedView) {
				v.Filters = []SavedViewFilter{{Field: "owner", Values: []string{" "}, AccessDecision: "allow"}}
			},
			want: ErrorEmptyFilterValue,
		},
		{
			name: "invalid-view-sort-field",
			mutate: func(v *SavedView) {
				v.Sort = []SavedViewSort{{Field: "secret", Direction: "asc", AccessDecision: "allow"}}
			},
			want: ErrorInvalidSortField,
		},
		{
			name: "invalid-view-sort-direction",
			mutate: func(v *SavedView) {
				v.Sort = []SavedViewSort{{Field: "owner", Direction: "sideways", AccessDecision: "allow"}}
			},
			want: ErrorInvalidSortDirection,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := baseInput()
			view := &SavedView{ID: "view-1", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow"}
			tc.mutate(view)
			in.SavedView = view
			got := Project(in)
			if got.Allowed || got.ErrorCode != tc.want || got.Projection != nil {
				t.Fatalf("got=%#v want=%s", got, tc.want)
			}
		})
	}
}

func TestProjectionFiltersAllSupportedFields(t *testing.T) {
	in := baseInput()
	in.Cases[0].Freshness = "stale"
	in.Cases[0].Incident = &IncidentContext{
		ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "allow", State: "available", Priority: "P1",
	}
	in.Cases[0].Findings = []FindingContext{{
		ID: "f-1", TenantRef: "tenant-a", CaseID: "case-a",
		AccessDecision: "allow", State: "available", Status: "confirmed",
	}}
	tests := []Filter{
		{Field: "status", Values: []string{"INVESTIGATING"}},
		{Field: "owner", Values: []string{" ALICE "}},
		{Field: "incident_priority", Values: []string{"p1"}},
		{Field: "finding_status", Values: []string{"confirmed"}},
		{Field: "freshness", Values: []string{"STALE"}},
	}
	for _, filter := range tests {
		t.Run(filter.Field, func(t *testing.T) {
			candidate := in
			candidate.Filters = []Filter{filter}
			got := Project(candidate)
			if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 1 {
				t.Fatalf("supported filter did not match: %#v", got)
			}
			candidate.Filters = []Filter{{Field: filter.Field, Values: []string{"does-not-match"}}}
			got = Project(candidate)
			if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 0 {
				t.Fatalf("nonmatching filter should yield empty queue: %#v", got)
			}
		})
	}
}

func TestProjectionSearchesEveryAllowedTextField(t *testing.T) {
	in := baseInput()
	in.Cases[0].NextAction = "Escalate evidence"
	queries := []string{"case-a", "case-001", "alpha", "alice", "evidence"}
	for _, query := range queries {
		candidate := in
		candidate.SearchQuery = query
		got := Project(candidate)
		if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 1 {
			t.Fatalf("query %q did not match allowed fields: %#v", query, got)
		}
	}
	in.SearchQuery = "not-present"
	got := Project(in)
	if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 0 {
		t.Fatalf("nonmatching search should be empty: %#v", got)
	}
}

func TestProjectionSortsEverySupportedKeyAndDirection(t *testing.T) {
	base := baseInput()
	second := base.Cases[0]
	second.ID = "case-b"
	second.HumanID = "CASE-002"
	second.Title = "Beta"
	second.Owner = "bob"
	second.Status = "closed"
	second.UpdatedAt = "2026-09-29T11:00:00Z"
	base.Cases = append(base.Cases, second)

	fields := []string{"updated_at", "human_id", "title", "owner", "status"}
	for _, field := range fields {
		for _, direction := range []string{"asc", "desc"} {
			t.Run(field+"-"+direction, func(t *testing.T) {
				in := base
				in.Sort = []SortSpec{{Field: field, Direction: direction}}
				got := Project(in)
				if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 2 {
					t.Fatalf("sort failed: %#v", got)
				}
				first := got.Projection.Rows[0].Case.ID
				want := "case-a"
				if direction == "desc" {
					want = "case-b"
				}
				if first != want {
					t.Fatalf("%s %s first=%s want=%s", field, direction, first, want)
				}
			})
		}
	}
}

func TestProjectionStableTieBreakerUsesCaseID(t *testing.T) {
	in := baseInput()
	second := in.Cases[0]
	second.ID = "case-b"
	in.Cases = append(in.Cases, second)
	in.Sort = []SortSpec{{Field: "status", Direction: "asc"}, {Field: "human_id", Direction: "asc"}}
	got := Project(in)
	if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 2 {
		t.Fatalf("unexpected projection: %#v", got)
	}
	if got.Projection.Rows[0].Case.ID != "case-a" || got.Projection.Rows[1].Case.ID != "case-b" {
		t.Fatalf("tie-break ordering is not deterministic: %#v", got.Projection.Rows)
	}
}

func baseInput() Input {
	return Input{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"}, EnvironmentRef: "env-prod",
		Cases: []CaseSnapshot{{
			ID: "case-a", TenantRef: "tenant-a", HumanID: "CASE-001", Title: "Alpha",
			Status: "investigating", Owner: "alice", UpdatedAt: "2026-09-29T10:00:00Z",
			AccessDecision: "allow", Freshness: "available", ActivityState: "available",
		}},
	}
}
