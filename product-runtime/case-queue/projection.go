package casequeue

import (
	"sort"
	"strings"
	"time"
)

var allowedFilterFields = map[string]bool{
	"status": true, "owner": true, "incident_priority": true, "finding_status": true, "freshness": true,
}

var allowedSortFields = map[string]bool{
	"updated_at": true, "human_id": true, "title": true, "owner": true, "status": true,
}

var defaultSort = []SortSpec{
	{Field: "updated_at", Direction: "desc"},
	{Field: "human_id", Direction: "asc"},
	{Field: "id", Direction: "asc"},
}

var stableTieBreakers = []string{"human_id", "id"}

type evaluatedRow struct {
	row                        Row
	updated                    time.Time
	permissionFilteredElements int
}

func Project(in Input) Result {
	result := Result{AuditRequired: true}
	reject := func(code ErrorCode) Result {
		result.ErrorCode = code
		return result
	}
	if strings.TrimSpace(in.MutationRequested) != "" {
		return reject(ErrorMutationOpenDecision)
	}
	if in.ExportRequested {
		return reject(ErrorExportForbidden)
	}
	if in.WorkQueueAggregationRequested {
		return reject(ErrorWorkQueueAggregationForbidden)
	}
	if strings.TrimSpace(in.TenantRef) == "" {
		return reject(ErrorMissingTenant)
	}
	if in.TenantRef == "*" {
		return reject(ErrorWildcardTenant)
	}
	if !contains(in.AuthorizedTenants, in.TenantRef) {
		return reject(ErrorCrossTenant)
	}
	if strings.TrimSpace(in.EnvironmentRef) == "" {
		return reject(ErrorMissingEnvironment)
	}

	requestFilters, code := validateFilters(in.Filters)
	if code != ErrorNone {
		return reject(code)
	}
	requestSort, code := validateSort(in.Sort)
	if code != ErrorNone {
		return reject(code)
	}

	diagnostics := map[string]bool{}
	permissionFilteredElements := 0
	appliedFields := []string{}
	viewFilters := []Filter{}
	viewSort := []SortSpec{}
	if in.SavedView != nil {
		view := in.SavedView
		if strings.TrimSpace(view.ID) == "" {
			return reject(ErrorMissingSavedViewID)
		}
		if view.TenantRef != in.TenantRef {
			return reject(ErrorCrossTenantSavedView)
		}
		if strings.TrimSpace(view.Version) == "" {
			return reject(ErrorMissingSavedViewVersion)
		}
		switch view.AccessDecision {
		case "allow":
		case "deny":
			return reject(ErrorSavedViewAccessDenied)
		default:
			return reject(ErrorInvalidSavedViewAccessDecision)
		}
		fieldSeen := map[string]bool{}
		for _, field := range view.Fields {
			switch field.AccessDecision {
			case "deny":
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
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
			appliedFields = append(appliedFields, name)
		}
		for _, filter := range view.Filters {
			switch filter.AccessDecision {
			case "deny":
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
				continue
			case "allow":
			default:
				return reject(ErrorInvalidSavedViewElementAccess)
			}
			validated, filterCode := validateFilters([]Filter{{Field: filter.Field, Values: filter.Values}})
			if filterCode != ErrorNone {
				return reject(filterCode)
			}
			viewFilters = append(viewFilters, validated...)
		}
		for _, spec := range view.Sort {
			switch spec.AccessDecision {
			case "deny":
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
				continue
			case "allow":
			default:
				return reject(ErrorInvalidSavedViewElementAccess)
			}
			validated, sortCode := validateSort([]SortSpec{{Field: spec.Field, Direction: spec.Direction}})
			if sortCode != ErrorNone {
				return reject(sortCode)
			}
			viewSort = append(viewSort, validated...)
		}
	}
	sort.Strings(appliedFields)

	effectiveFilters := append([]Filter{}, viewFilters...)
	effectiveFilters = append(effectiveFilters, requestFilters...)
	effectiveSort := append([]SortSpec{}, requestSort...)
	if len(effectiveSort) == 0 {
		effectiveSort = append([]SortSpec{}, viewSort...)
	}
	if len(effectiveSort) == 0 {
		effectiveSort = append([]SortSpec{}, defaultSort...)
	}
	effectiveSort = appendTieBreakers(effectiveSort, stableTieBreakers)

	search := strings.ToLower(strings.TrimSpace(in.SearchQuery))
	permissionFilteredCases := 0
	caseIDs := map[string]bool{}
	rows := []evaluatedRow{}
	for _, item := range in.Cases {
		if item.TenantRef != in.TenantRef {
			return reject(ErrorCrossTenantCase)
		}
		if strings.TrimSpace(item.ID) == "" {
			return reject(ErrorMissingCaseID)
		}
		if caseIDs[item.ID] {
			return reject(ErrorDuplicateCaseID)
		}
		caseIDs[item.ID] = true
		switch item.AccessDecision {
		case "deny":
			permissionFilteredCases++
			diagnostics["case:permission-filtered"] = true
			continue
		case "allow":
		default:
			return reject(ErrorInvalidCaseAccessDecision)
		}
		if strings.TrimSpace(item.HumanID) == "" || strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Status) == "" {
			return reject(ErrorIncompleteCaseSnapshot)
		}
		updated, err := time.Parse(time.RFC3339, item.UpdatedAt)
		if err != nil {
			return reject(ErrorInvalidCaseUpdatedAt)
		}
		visible := cloneCase(item)
		visible.Incident = nil
		visible.Findings = nil
		row := evaluatedRow{row: Row{Case: visible, FindingStatuses: []string{}, Diagnostics: []string{}}, updated: updated}
		switch item.Freshness {
		case "available":
		case "partial":
			row.row.Partial = true
			row.row.Diagnostics = append(row.row.Diagnostics, "case:"+item.ID+":partial")
		case "stale":
			row.row.Stale = true
			row.row.Diagnostics = append(row.row.Diagnostics, "case:"+item.ID+":stale")
		default:
			return reject(ErrorInvalidCaseFreshness)
		}

		if item.Incident != nil {
			incident := item.Incident
			if incident.TenantRef != in.TenantRef {
				return reject(ErrorCrossTenantIncident)
			}
			switch incident.AccessDecision {
			case "deny":
				row.permissionFilteredElements++
				row.row.Diagnostics = append(row.row.Diagnostics, "incident-context:permission-filtered")
			case "allow":
				switch incident.State {
				case "available":
					row.row.IncidentPriority = incident.Priority
					cloned := *incident
					row.row.Case.Incident = &cloned
				case "partial":
					row.row.Partial = true
					row.row.IncidentPriority = incident.Priority
					cloned := *incident
					row.row.Case.Incident = &cloned
					row.row.Diagnostics = append(row.row.Diagnostics, "incident-context:"+item.ID+":partial")
				case "stale":
					row.row.Stale = true
					row.row.IncidentPriority = incident.Priority
					cloned := *incident
					row.row.Case.Incident = &cloned
					row.row.Diagnostics = append(row.row.Diagnostics, "incident-context:"+item.ID+":stale")
				case "unavailable":
					row.row.Partial = true
					row.row.Diagnostics = append(row.row.Diagnostics, "incident-context:"+item.ID+":unavailable")
				default:
					return reject(ErrorInvalidIncidentState)
				}
			default:
				return reject(ErrorInvalidIncidentAccessDecision)
			}
		}

		for _, finding := range item.Findings {
			if finding.TenantRef != in.TenantRef {
				return reject(ErrorCrossTenantFinding)
			}
			if finding.CaseID != item.ID {
				return reject(ErrorFindingCaseMismatch)
			}
			switch finding.AccessDecision {
			case "deny":
				row.permissionFilteredElements++
				row.row.Diagnostics = append(row.row.Diagnostics, "finding-context:permission-filtered")
				continue
			case "allow":
			default:
				return reject(ErrorInvalidFindingAccessDecision)
			}
			switch finding.State {
			case "available":
				row.row.FindingStatuses = append(row.row.FindingStatuses, finding.Status)
				row.row.Case.Findings = append(row.row.Case.Findings, finding)
			case "partial":
				row.row.Partial = true
				row.row.FindingStatuses = append(row.row.FindingStatuses, finding.Status)
				row.row.Case.Findings = append(row.row.Case.Findings, finding)
				row.row.Diagnostics = append(row.row.Diagnostics, "finding-context:"+item.ID+":partial")
			case "stale":
				row.row.Stale = true
				row.row.FindingStatuses = append(row.row.FindingStatuses, finding.Status)
				row.row.Case.Findings = append(row.row.Case.Findings, finding)
				row.row.Diagnostics = append(row.row.Diagnostics, "finding-context:"+item.ID+":stale")
			case "unavailable":
				row.row.Partial = true
				row.row.Diagnostics = append(row.row.Diagnostics, "finding-context:"+item.ID+":unavailable")
			default:
				return reject(ErrorInvalidFindingState)
			}
		}

		switch item.ActivityState {
		case "", "available":
		case "partial", "unavailable":
			row.row.Partial = true
			row.row.Diagnostics = append(row.row.Diagnostics, "activity:"+item.ID+":"+item.ActivityState)
		case "stale":
			row.row.Stale = true
			row.row.Diagnostics = append(row.row.Diagnostics, "activity:"+item.ID+":stale")
		default:
			return reject(ErrorInvalidActivityState)
		}
		if !matchesSearch(row, search) || !matchesFilters(row, effectiveFilters) {
			continue
		}
		rows = append(rows, row)
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rowLess(rows[i], rows[j], effectiveSort)
	})
	states := map[string]bool{}
	if len(rows) == 0 {
		states["empty"] = true
	} else {
		states["available"] = true
	}
	if permissionFilteredCases > 0 || permissionFilteredElements > 0 {
		states["permission-filtered"] = true
	}
	if in.ViewDirty {
		states["view-dirty"] = true
	}
	projectedRows := make([]Row, 0, len(rows))
	for _, row := range rows {
		if row.row.Partial {
			states["partial"] = true
		}
		if row.row.Stale {
			states["stale"] = true
		}
		permissionFilteredElements += row.permissionFilteredElements
		for _, diagnostic := range row.row.Diagnostics {
			diagnostics[diagnostic] = true
		}
		row.row.Diagnostics = append([]string{}, row.row.Diagnostics...)
		row.row.FindingStatuses = append([]string{}, row.row.FindingStatuses...)
		projectedRows = append(projectedRows, row.row)
	}
	if permissionFilteredElements > 0 {
		states["permission-filtered"] = true
	}

	result.Allowed = true
	result.Projection = &Projection{
		TenantRef: in.TenantRef, EnvironmentRef: in.EnvironmentRef,
		States: sortedKeys(states), Diagnostics: sortedKeys(diagnostics), Rows: projectedRows,
		PermissionFilteredCases: permissionFilteredCases, PermissionFilteredElements: permissionFilteredElements,
		AppliedSavedViewFields: append([]string{}, appliedFields...), AppliedSort: append([]SortSpec{}, effectiveSort...),
	}
	return result
}

func cloneCase(in CaseSnapshot) CaseSnapshot {
	out := in
	out.Findings = append([]FindingContext{}, in.Findings...)
	if in.Incident != nil {
		cloned := *in.Incident
		out.Incident = &cloned
	}
	return out
}

func validateFilters(filters []Filter) ([]Filter, ErrorCode) {
	out := make([]Filter, 0, len(filters))
	for _, filter := range filters {
		if !allowedFilterFields[filter.Field] {
			return nil, ErrorInvalidFilterField
		}
		if len(filter.Values) == 0 {
			return nil, ErrorEmptyFilterValue
		}
		copyFilter := Filter{Field: filter.Field, Values: make([]string, 0, len(filter.Values))}
		for _, value := range filter.Values {
			value = strings.TrimSpace(value)
			if value == "" {
				return nil, ErrorEmptyFilterValue
			}
			copyFilter.Values = append(copyFilter.Values, strings.ToLower(value))
		}
		out = append(out, copyFilter)
	}
	return out, ErrorNone
}

func validateSort(specs []SortSpec) ([]SortSpec, ErrorCode) {
	out := make([]SortSpec, 0, len(specs))
	seen := map[string]bool{}
	for _, spec := range specs {
		if !allowedSortFields[spec.Field] {
			return nil, ErrorInvalidSortField
		}
		if spec.Direction != "asc" && spec.Direction != "desc" {
			return nil, ErrorInvalidSortDirection
		}
		if seen[spec.Field] {
			return nil, ErrorDuplicateSortField
		}
		seen[spec.Field] = true
		out = append(out, spec)
	}
	return out, ErrorNone
}

func appendTieBreakers(specs []SortSpec, tie []string) []SortSpec {
	out := append([]SortSpec{}, specs...)
	seen := map[string]bool{}
	for _, spec := range out {
		seen[spec.Field] = true
	}
	for _, field := range tie {
		if !seen[field] {
			out = append(out, SortSpec{Field: field, Direction: "asc"})
			seen[field] = true
		}
	}
	return out
}

func matchesSearch(row evaluatedRow, query string) bool {
	if query == "" {
		return true
	}
	values := []string{row.row.Case.ID, row.row.Case.HumanID, row.row.Case.Title, row.row.Case.Owner, row.row.Case.NextAction}
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}
	return false
}

func matchesFilters(row evaluatedRow, filters []Filter) bool {
	for _, filter := range filters {
		matched := false
		switch filter.Field {
		case "status":
			matched = containsNormalized(filter.Values, row.row.Case.Status)
		case "owner":
			matched = containsNormalized(filter.Values, row.row.Case.Owner)
		case "incident_priority":
			matched = containsNormalized(filter.Values, row.row.IncidentPriority)
		case "finding_status":
			for _, status := range row.row.FindingStatuses {
				if containsNormalized(filter.Values, status) {
					matched = true
					break
				}
			}
		case "freshness":
			matched = containsNormalized(filter.Values, row.row.Case.Freshness)
		}
		if !matched {
			return false
		}
	}
	return true
}

func containsNormalized(values []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, candidate := range values {
		if strings.ToLower(strings.TrimSpace(candidate)) == value {
			return true
		}
	}
	return false
}

func rowLess(a, b evaluatedRow, specs []SortSpec) bool {
	for _, spec := range specs {
		cmp := compareField(a, b, spec.Field)
		if cmp == 0 {
			continue
		}
		if spec.Direction == "desc" {
			return cmp > 0
		}
		return cmp < 0
	}
	return false
}

func compareField(a, b evaluatedRow, field string) int {
	switch field {
	case "updated_at":
		if a.updated.Before(b.updated) {
			return -1
		}
		if a.updated.After(b.updated) {
			return 1
		}
		return 0
	case "id":
		return strings.Compare(strings.ToLower(a.row.Case.ID), strings.ToLower(b.row.Case.ID))
	case "human_id":
		return strings.Compare(strings.ToLower(a.row.Case.HumanID), strings.ToLower(b.row.Case.HumanID))
	case "title":
		return strings.Compare(strings.ToLower(a.row.Case.Title), strings.ToLower(b.row.Case.Title))
	case "owner":
		return strings.Compare(strings.ToLower(a.row.Case.Owner), strings.ToLower(b.row.Case.Owner))
	case "status":
		return strings.Compare(strings.ToLower(a.row.Case.Status), strings.ToLower(b.row.Case.Status))
	default:
		return 0
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func sortedKeys(values map[string]bool) []string {
	out := make([]string, 0, len(values))
	for value, present := range values {
		if present {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}
