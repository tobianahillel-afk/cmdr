package casequeue

import (
	"reflect"
	"testing"
)

func TestApplySavedViewFallsBackWithoutLeakingDeniedConfiguration(t *testing.T) {
	view := &SavedView{
		ID: "secret-view", TenantRef: "tenant-a", Version: "v7", AccessDecision: "deny",
		Fields: []SavedViewField{{Name: "hidden-field", AccessDecision: "allow"}},
		Filters: []SavedViewFilter{{Field: "owner", Values: []string{"secret-owner"}, AccessDecision: "allow"}},
		Sort: []SavedViewSort{{Field: "owner", Direction: "desc", AccessDecision: "allow"}},
	}
	got := ApplySavedView("tenant-a", view)
	if !got.Allowed || !got.Fallback || got.Applied || got.Sanitized != nil ||
		got.FallbackReason != "saved-view-unavailable" ||
		len(got.Fields) != 0 || len(got.Filters) != 0 || len(got.Sort) != 0 ||
		!reflect.DeepEqual(got.Diagnostics, []string{savedViewFallbackDiagnostic}) {
		t.Fatalf("denied view did not fail safely to default: %#v", got)
	}
}

func TestApplySavedViewMissingSnapshotFallsBackExplicitly(t *testing.T) {
	got := ApplySavedView("tenant-a", nil)
	if !got.Allowed || !got.Fallback || got.Applied || got.Sanitized != nil ||
		got.FallbackReason != "saved-view-unavailable" {
		t.Fatalf("missing view fallback mismatch: %#v", got)
	}
}

func TestApplySavedViewRejectsCrossTenantAndMalformedIdentity(t *testing.T) {
	tests := []struct {
		name string
		view *SavedView
		want ErrorCode
	}{
		{"missing-id", &SavedView{TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow"}, ErrorMissingSavedViewID},
		{"cross-tenant", &SavedView{ID: "v", TenantRef: "tenant-b", Version: "v1", AccessDecision: "allow"}, ErrorCrossTenantSavedView},
		{"missing-version", &SavedView{ID: "v", TenantRef: "tenant-a", AccessDecision: "allow"}, ErrorMissingSavedViewVersion},
		{"invalid-access", &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "maybe"}, ErrorInvalidSavedViewAccessDecision},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ApplySavedView("tenant-a", tc.view)
			if got.Allowed || got.ErrorCode != tc.want || got.Sanitized != nil {
				t.Fatalf("got=%#v want=%s", got, tc.want)
			}
		})
	}
}

func TestApplySavedViewReevaluatesEveryElementPermission(t *testing.T) {
	view := &SavedView{
		ID: "view-1", TenantRef: "tenant-a", Version: "v3", AccessDecision: "allow",
		Fields: []SavedViewField{
			{Name: "title", AccessDecision: "allow"},
			{Name: "owner", AccessDecision: "deny"},
		},
		Filters: []SavedViewFilter{
			{Field: "status", Values: []string{" INVESTIGATING "}, AccessDecision: "allow"},
			{Field: "owner", Values: []string{"hidden-owner"}, AccessDecision: "deny"},
		},
		Sort: []SavedViewSort{
			{Field: "updated_at", Direction: "desc", AccessDecision: "allow"},
			{Field: "owner", Direction: "asc", AccessDecision: "deny"},
		},
	}
	got := ApplySavedView("tenant-a", view)
	if !got.Allowed || !got.Applied || got.Fallback || got.Sanitized == nil {
		t.Fatalf("valid view rejected: %#v", got)
	}
	if got.PermissionFilteredElements != 3 ||
		!reflect.DeepEqual(got.Fields, []string{"title"}) ||
		!reflect.DeepEqual(got.Filters, []Filter{{Field: "status", Values: []string{"investigating"}}}) ||
		!reflect.DeepEqual(got.Sort, []SortSpec{{Field: "updated_at", Direction: "desc"}}) ||
		!reflect.DeepEqual(got.Diagnostics, []string{"saved-view:permission-filtered"}) {
		t.Fatalf("permission reevaluation mismatch: %#v", got)
	}
	if len(got.Sanitized.Fields) != 1 || got.Sanitized.Fields[0].Name != "title" ||
		len(got.Sanitized.Filters) != 1 || got.Sanitized.Filters[0].Field != "status" ||
		len(got.Sanitized.Sort) != 1 || got.Sanitized.Sort[0].Field != "updated_at" {
		t.Fatalf("sanitized view leaked denied elements: %#v", got.Sanitized)
	}
}

func TestApplySavedViewDoesNotAliasCallerOwnedConfiguration(t *testing.T) {
	view := &SavedView{
		ID: "view-1", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow",
		Fields: []SavedViewField{{Name: "title", AccessDecision: "allow"}},
		Filters: []SavedViewFilter{{Field: "owner", Values: []string{"alice"}, AccessDecision: "allow"}},
		Sort: []SavedViewSort{{Field: "title", Direction: "asc", AccessDecision: "allow"}},
	}
	got := ApplySavedView("tenant-a", view)
	if !got.Allowed || got.Sanitized == nil {
		t.Fatalf("unexpected result: %#v", got)
	}
	view.Fields[0].Name = "mutated"
	view.Filters[0].Values[0] = "mutated"
	view.Sort[0].Field = "owner"
	if got.Sanitized.Fields[0].Name != "title" ||
		got.Sanitized.Filters[0].Values[0] != "alice" ||
		got.Sanitized.Sort[0].Field != "title" {
		t.Fatalf("adapter aliases caller-owned view: %#v", got.Sanitized)
	}
}

func TestApplySavedViewRejectsInvalidAllowedElements(t *testing.T) {
	tests := []struct {
		name string
		edit func(*SavedView)
		want ErrorCode
	}{
		{"invalid-field-access", func(v *SavedView) { v.Fields = []SavedViewField{{Name: "title", AccessDecision: "maybe"}} }, ErrorInvalidSavedViewElementAccess},
		{"blank-field", func(v *SavedView) { v.Fields = []SavedViewField{{Name: " ", AccessDecision: "allow"}} }, ErrorInvalidSavedViewField},
		{"duplicate-field", func(v *SavedView) { v.Fields = []SavedViewField{{Name: "title", AccessDecision: "allow"}, {Name: "title", AccessDecision: "allow"}} }, ErrorDuplicateSavedViewField},
		{"invalid-filter-access", func(v *SavedView) { v.Filters = []SavedViewFilter{{Field: "owner", Values: []string{"a"}, AccessDecision: "maybe"}} }, ErrorInvalidSavedViewElementAccess},
		{"invalid-filter", func(v *SavedView) { v.Filters = []SavedViewFilter{{Field: "secret", Values: []string{"a"}, AccessDecision: "allow"}} }, ErrorInvalidFilterField},
		{"invalid-sort-access", func(v *SavedView) { v.Sort = []SavedViewSort{{Field: "owner", Direction: "asc", AccessDecision: "maybe"}} }, ErrorInvalidSavedViewElementAccess},
		{"invalid-sort", func(v *SavedView) { v.Sort = []SavedViewSort{{Field: "owner", Direction: "sideways", AccessDecision: "allow"}} }, ErrorInvalidSortDirection},
		{"duplicate-sort", func(v *SavedView) { v.Sort = []SavedViewSort{{Field: "owner", Direction: "asc", AccessDecision: "allow"}, {Field: "owner", Direction: "desc", AccessDecision: "allow"}} }, ErrorDuplicateSortField},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			view := &SavedView{ID: "v", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow"}
			tc.edit(view)
			got := ApplySavedView("tenant-a", view)
			if got.Allowed || got.ErrorCode != tc.want || got.Sanitized != nil {
				t.Fatalf("got=%#v want=%s", got, tc.want)
			}
		})
	}
}

func TestProjectWithSavedViewFallbackUsesDefaultViewWithoutLeaks(t *testing.T) {
	in := baseInput()
	in.SavedView = &SavedView{
		ID: "hidden", TenantRef: "tenant-a", Version: "v1", AccessDecision: "deny",
		Fields: []SavedViewField{{Name: "hidden-field", AccessDecision: "allow"}},
		Filters: []SavedViewFilter{{Field: "owner", Values: []string{"nobody"}, AccessDecision: "allow"}},
		Sort: []SavedViewSort{{Field: "owner", Direction: "desc", AccessDecision: "allow"}},
	}
	got := ProjectWithSavedViewFallback(in)
	if !got.Allowed || got.Projection == nil || len(got.Projection.Rows) != 1 {
		t.Fatalf("fallback projection rejected: %#v", got)
	}
	if len(got.Projection.AppliedSavedViewFields) != 0 ||
		!contains(got.Projection.Diagnostics, savedViewFallbackDiagnostic) {
		t.Fatalf("fallback exposed hidden Saved View state: %#v", got.Projection)
	}
	wantSort := appendTieBreakers(defaultSort, stableTieBreakers)
	if !reflect.DeepEqual(got.Projection.AppliedSort, wantSort) {
		t.Fatalf("fallback did not use deterministic default sort: got=%v want=%v", got.Projection.AppliedSort, wantSort)
	}
}

func TestProjectWithSavedViewFallbackPreservesPermissionFilteredAndDirtyState(t *testing.T) {
	in := baseInput()
	in.ViewDirty = true
	in.SavedView = &SavedView{
		ID: "view-1", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow",
		Fields: []SavedViewField{{Name: "title", AccessDecision: "allow"}, {Name: "secret", AccessDecision: "deny"}},
	}
	got := ProjectWithSavedViewFallback(in)
	if !got.Allowed || got.Projection == nil {
		t.Fatalf("projection rejected: %#v", got)
	}
	if got.Projection.PermissionFilteredElements != 1 ||
		!contains(got.Projection.States, "permission-filtered") ||
		!contains(got.Projection.States, "view-dirty") ||
		!contains(got.Projection.Diagnostics, "saved-view:permission-filtered") {
		t.Fatalf("saved-view diagnostics lost: %#v", got.Projection)
	}
}

func TestProjectWithSavedViewFallbackRejectsCrossTenantView(t *testing.T) {
	in := baseInput()
	in.SavedView = &SavedView{ID: "view-1", TenantRef: "tenant-b", Version: "v1", AccessDecision: "allow"}
	got := ProjectWithSavedViewFallback(in)
	if got.Allowed || got.ErrorCode != ErrorCrossTenantSavedView || got.Projection != nil {
		t.Fatalf("cross-tenant view did not fail closed: %#v", got)
	}
}
