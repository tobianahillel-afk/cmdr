package casequeue

import (
	"sort"
	"strings"
)

const savedViewFallbackDiagnostic = "saved-view:fallback"

type SavedViewApplication struct {
	Allowed                    bool
	ErrorCode                  ErrorCode
	Applied                    bool
	Fallback                   bool
	FallbackReason             string
	PermissionFilteredElements int
	Diagnostics                []string
	Fields                     []string
	Filters                    []Filter
	Sort                       []SortSpec
	Sanitized                  *SavedView
}

func ApplySavedView(tenantRef string, view *SavedView) SavedViewApplication {
	result := SavedViewApplication{Allowed: true, Diagnostics: []string{}, Fields: []string{}, Filters: []Filter{}, Sort: []SortSpec{}}
	fallback := func() SavedViewApplication {
		result.Fallback = true
		result.FallbackReason = "saved-view-unavailable"
		result.Diagnostics = []string{savedViewFallbackDiagnostic}
		result.Fields = []string{}
		result.Filters = []Filter{}
		result.Sort = []SortSpec{}
		result.Sanitized = nil
		return result
	}
	reject := func(code ErrorCode) SavedViewApplication {
		return SavedViewApplication{
			Allowed: false, ErrorCode: code,
			Diagnostics: []string{}, Fields: []string{}, Filters: []Filter{}, Sort: []SortSpec{},
		}
	}

	if view == nil {
		return fallback()
	}
	if strings.TrimSpace(view.ID) == "" {
		return reject(ErrorMissingSavedViewID)
	}
	if view.TenantRef != tenantRef {
		return reject(ErrorCrossTenantSavedView)
	}
	if strings.TrimSpace(view.Version) == "" {
		return reject(ErrorMissingSavedViewVersion)
	}
	switch view.AccessDecision {
	case "deny":
		return fallback()
	case "allow":
	default:
		return reject(ErrorInvalidSavedViewAccessDecision)
	}

	sanitized := &SavedView{
		ID: view.ID, TenantRef: view.TenantRef, Version: view.Version, AccessDecision: "allow",
		Fields: []SavedViewField{}, Filters: []SavedViewFilter{}, Sort: []SavedViewSort{},
	}
	fieldSeen := map[string]bool{}
	for _, field := range view.Fields {
		switch field.AccessDecision {
		case "deny":
			result.PermissionFilteredElements++
			continue
		case "allow":
		default:
			return reject(ErrorInvalidSavedViewElementAccess)
		}
		name := strings.TrimSpace(field.Name)
		if name == "" {
			return reject(ErrorInvalidSavedViewField)
		}
		if fieldSeen[name] {
			return reject(ErrorDuplicateSavedViewField)
		}
		fieldSeen[name] = true
		result.Fields = append(result.Fields, name)
		sanitized.Fields = append(sanitized.Fields, SavedViewField{Name: name, AccessDecision: "allow"})
	}
	sort.Strings(result.Fields)
	sort.Slice(sanitized.Fields, func(i, j int) bool { return sanitized.Fields[i].Name < sanitized.Fields[j].Name })

	for _, filter := range view.Filters {
		switch filter.AccessDecision {
		case "deny":
			result.PermissionFilteredElements++
			continue
		case "allow":
		default:
			return reject(ErrorInvalidSavedViewElementAccess)
		}
		validated, code := validateFilters([]Filter{{Field: filter.Field, Values: filter.Values}})
		if code != ErrorNone {
			return reject(code)
		}
		result.Filters = append(result.Filters, cloneFilter(validated[0]))
		sanitized.Filters = append(sanitized.Filters, SavedViewFilter{
			Field: validated[0].Field, Values: append([]string{}, validated[0].Values...), AccessDecision: "allow",
		})
	}

	sortSeen := map[string]bool{}
	for _, spec := range view.Sort {
		switch spec.AccessDecision {
		case "deny":
			result.PermissionFilteredElements++
			continue
		case "allow":
		default:
			return reject(ErrorInvalidSavedViewElementAccess)
		}
		if sortSeen[spec.Field] {
			return reject(ErrorDuplicateSortField)
		}
		sortSeen[spec.Field] = true
		validated, code := validateSort([]SortSpec{{Field: spec.Field, Direction: spec.Direction}})
		if code != ErrorNone {
			return reject(code)
		}
		result.Sort = append(result.Sort, validated[0])
		sanitized.Sort = append(sanitized.Sort, SavedViewSort{
			Field: validated[0].Field, Direction: validated[0].Direction, AccessDecision: "allow",
		})
	}

	if result.PermissionFilteredElements > 0 {
		result.Diagnostics = append(result.Diagnostics, "saved-view:permission-filtered")
	}
	result.Applied = true
	result.Sanitized = sanitized
	return result
}

func ProjectWithSavedViewFallback(in Input) Result {
	application := ApplySavedView(in.TenantRef, in.SavedView)
	if !application.Allowed {
		return Result{AuditRequired: true, ErrorCode: application.ErrorCode}
	}

	projectInput := cloneInputForSavedView(in)
	projectInput.SavedView = application.Sanitized
	got := Project(projectInput)
	if !got.Allowed || got.Projection == nil {
		return got
	}

	if application.Fallback {
		got.Projection.Diagnostics = appendUniqueSorted(got.Projection.Diagnostics, savedViewFallbackDiagnostic)
	}
	if application.PermissionFilteredElements > 0 {
		got.Projection.PermissionFilteredElements += application.PermissionFilteredElements
		got.Projection.States = appendUniqueSorted(got.Projection.States, "permission-filtered")
		got.Projection.Diagnostics = appendUniqueSorted(got.Projection.Diagnostics, "saved-view:permission-filtered")
	}
	return got
}

func cloneInputForSavedView(in Input) Input {
	out := in
	out.AuthorizedTenants = append([]string{}, in.AuthorizedTenants...)
	out.Filters = make([]Filter, len(in.Filters))
	for i, filter := range in.Filters {
		out.Filters[i] = cloneFilter(filter)
	}
	out.Sort = append([]SortSpec{}, in.Sort...)
	out.Cases = make([]CaseSnapshot, len(in.Cases))
	for i, item := range in.Cases {
		out.Cases[i] = cloneCase(item)
	}
	return out
}

func cloneFilter(in Filter) Filter {
	return Filter{Field: in.Field, Values: append([]string{}, in.Values...)}
}

func appendUniqueSorted(values []string, value string) []string {
	out := append([]string{}, values...)
	for _, existing := range out {
		if existing == value {
			sort.Strings(out)
			return out
		}
	}
	out = append(out, value)
	sort.Strings(out)
	return out
}
