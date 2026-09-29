// Package casequeue implements the bounded read-only CAP-INV-101 Case Queue projection.
package casequeue

type ErrorCode string

const (
	ErrorNone                           ErrorCode = ""
	ErrorMutationOpenDecision           ErrorCode = "mutation-open-decision"
	ErrorExportForbidden                ErrorCode = "export-forbidden"
	ErrorWorkQueueAggregationForbidden  ErrorCode = "work-queue-aggregation-forbidden"
	ErrorMissingTenant                  ErrorCode = "missing-tenant"
	ErrorWildcardTenant                 ErrorCode = "wildcard-tenant"
	ErrorCrossTenant                    ErrorCode = "cross-tenant"
	ErrorMissingEnvironment             ErrorCode = "missing-environment"
	ErrorInvalidFilterField             ErrorCode = "invalid-filter-field"
	ErrorEmptyFilterValue               ErrorCode = "empty-filter-value"
	ErrorInvalidSortField               ErrorCode = "invalid-sort-field"
	ErrorInvalidSortDirection           ErrorCode = "invalid-sort-direction"
	ErrorDuplicateSortField             ErrorCode = "duplicate-sort-field"
	ErrorMissingSavedViewID             ErrorCode = "missing-saved-view-id"
	ErrorCrossTenantSavedView           ErrorCode = "cross-tenant-saved-view"
	ErrorMissingSavedViewVersion        ErrorCode = "missing-saved-view-version"
	ErrorSavedViewAccessDenied          ErrorCode = "saved-view-access-denied"
	ErrorInvalidSavedViewAccessDecision ErrorCode = "invalid-saved-view-access-decision"
	ErrorInvalidSavedViewElementAccess  ErrorCode = "invalid-saved-view-element-access-decision"
	ErrorInvalidSavedViewField          ErrorCode = "invalid-saved-view-field"
	ErrorDuplicateSavedViewField        ErrorCode = "duplicate-saved-view-field"
	ErrorCrossTenantCase                ErrorCode = "cross-tenant-case"
	ErrorMissingCaseID                  ErrorCode = "missing-case-id"
	ErrorDuplicateCaseID                ErrorCode = "duplicate-case-id"
	ErrorInvalidCaseAccessDecision      ErrorCode = "invalid-case-access-decision"
	ErrorIncompleteCaseSnapshot         ErrorCode = "incomplete-case-snapshot"
	ErrorInvalidCaseUpdatedAt           ErrorCode = "invalid-case-updated-at"
	ErrorInvalidCaseFreshness           ErrorCode = "invalid-case-freshness"
	ErrorCrossTenantIncident            ErrorCode = "cross-tenant-incident"
	ErrorInvalidIncidentAccessDecision  ErrorCode = "invalid-incident-access-decision"
	ErrorInvalidIncidentState           ErrorCode = "invalid-incident-state"
	ErrorCrossTenantFinding             ErrorCode = "cross-tenant-finding"
	ErrorFindingCaseMismatch            ErrorCode = "finding-case-mismatch"
	ErrorInvalidFindingAccessDecision   ErrorCode = "invalid-finding-access-decision"
	ErrorInvalidFindingState            ErrorCode = "invalid-finding-state"
	ErrorInvalidActivityState           ErrorCode = "invalid-activity-state"
)

type Filter struct {
	Field  string   `json:"field"`
	Values []string `json:"values"`
}

type SortSpec struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type SavedViewField struct {
	Name           string `json:"name"`
	AccessDecision string `json:"access_decision"`
}

type SavedViewFilter struct {
	Field          string   `json:"field"`
	Values         []string `json:"values"`
	AccessDecision string   `json:"access_decision"`
}

type SavedViewSort struct {
	Field          string `json:"field"`
	Direction      string `json:"direction"`
	AccessDecision string `json:"access_decision"`
}

type SavedView struct {
	ID             string            `json:"id"`
	TenantRef      string            `json:"tenant_ref"`
	Version        string            `json:"version"`
	AccessDecision string            `json:"access_decision"`
	Fields         []SavedViewField  `json:"fields"`
	Filters        []SavedViewFilter `json:"filters"`
	Sort           []SavedViewSort   `json:"sort"`
}

type IncidentContext struct {
	ID             string `json:"id"`
	TenantRef      string `json:"tenant_ref"`
	AccessDecision string `json:"access_decision"`
	State          string `json:"state"`
	Priority       string `json:"priority"`
	Impact         string `json:"impact"`
	Owner          string `json:"owner"`
	UpdatedAt      string `json:"updated_at"`
}

type FindingContext struct {
	ID             string `json:"id"`
	TenantRef      string `json:"tenant_ref"`
	CaseID         string `json:"case_id"`
	AccessDecision string `json:"access_decision"`
	State          string `json:"state"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at"`
}

type CaseSnapshot struct {
	ID             string           `json:"id"`
	TenantRef      string           `json:"tenant_ref"`
	HumanID        string           `json:"human_id"`
	Title          string           `json:"title"`
	Status         string           `json:"status"`
	Owner          string           `json:"owner"`
	NextAction     string           `json:"next_action"`
	UpdatedAt      string           `json:"updated_at"`
	AccessDecision string           `json:"access_decision"`
	Freshness      string           `json:"freshness"`
	Incident       *IncidentContext `json:"incident"`
	Findings       []FindingContext `json:"findings"`
	ActivityState  string           `json:"activity_state"`
}

type Input struct {
	TenantRef                     string         `json:"tenant_ref"`
	AuthorizedTenants             []string       `json:"authorized_tenants"`
	EnvironmentRef                string         `json:"environment_ref"`
	SearchQuery                   string         `json:"search_query"`
	Filters                       []Filter       `json:"filters"`
	Sort                          []SortSpec     `json:"sort"`
	ViewDirty                     bool           `json:"view_dirty"`
	SavedView                     *SavedView     `json:"saved_view"`
	Cases                         []CaseSnapshot `json:"cases"`
	MutationRequested             string         `json:"mutation_requested"`
	ExportRequested               bool           `json:"export_requested"`
	WorkQueueAggregationRequested bool           `json:"work_queue_aggregation_requested"`
}

type Row struct {
	Case             CaseSnapshot `json:"case"`
	IncidentPriority string       `json:"incident_priority,omitempty"`
	FindingStatuses  []string     `json:"finding_statuses,omitempty"`
	Partial          bool         `json:"partial"`
	Stale            bool         `json:"stale"`
	Diagnostics      []string     `json:"diagnostics,omitempty"`
}

type Projection struct {
	TenantRef                  string     `json:"tenant_ref"`
	EnvironmentRef             string     `json:"environment_ref"`
	States                     []string   `json:"states"`
	Diagnostics                []string   `json:"diagnostics"`
	Rows                       []Row      `json:"rows"`
	PermissionFilteredCases    int        `json:"permission_filtered_cases"`
	PermissionFilteredElements int        `json:"permission_filtered_elements"`
	AppliedSavedViewFields     []string   `json:"applied_saved_view_fields"`
	AppliedSort                []SortSpec `json:"applied_sort"`
}

type Result struct {
	Allowed             bool        `json:"allowed"`
	ErrorCode           ErrorCode   `json:"error_code"`
	AuditRequired       bool        `json:"audit_required"`
	MutationExecuted    bool        `json:"mutation_executed"`
	ExportJobCreated    bool        `json:"export_job_created"`
	WorkQueueAggregated bool        `json:"work_queue_aggregated"`
	Projection          *Projection `json:"projection,omitempty"`
}
