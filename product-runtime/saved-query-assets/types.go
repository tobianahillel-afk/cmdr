// Package savedqueryassets implements the bounded read-only Saved Search / Query Asset projection.
package savedqueryassets

type ErrorCode string

const (
	ErrorNone                        ErrorCode = ""
	ErrorMutationOpenDecision        ErrorCode = "mutation-open-decision"
	ErrorSearchJobExecutionForbidden ErrorCode = "search-job-execution-forbidden"
	ErrorMissingTenant               ErrorCode = "missing-tenant"
	ErrorWildcardTenant              ErrorCode = "wildcard-tenant"
	ErrorCrossTenant                 ErrorCode = "cross-tenant"
	ErrorMissingEnvironment          ErrorCode = "missing-environment"
	ErrorMissingAssetID              ErrorCode = "missing-asset-id"
	ErrorInvalidAssetKind            ErrorCode = "invalid-asset-kind"
	ErrorAssetAccessDenied           ErrorCode = "asset-access-denied"
	ErrorMissingQueryRef             ErrorCode = "missing-query-ref"
	ErrorCrossTenantQuery            ErrorCode = "cross-tenant-query"
	ErrorMissingQueryVersion         ErrorCode = "missing-query-version"
	ErrorQueryAccessDenied           ErrorCode = "query-access-denied"
	ErrorMissingAuthor               ErrorCode = "missing-author"
	ErrorInvalidValidationState      ErrorCode = "invalid-validation-state"
	ErrorInvalidParameter            ErrorCode = "invalid-parameter"
	ErrorMissingSourcePrerequisites  ErrorCode = "missing-source-prerequisites"
	ErrorInvalidSource               ErrorCode = "invalid-source"
	ErrorCrossTenantSource           ErrorCode = "cross-tenant-source"
	ErrorInvalidSourceState          ErrorCode = "invalid-source-state"
	ErrorInvalidField                ErrorCode = "invalid-field"
	ErrorInvalidFieldState           ErrorCode = "invalid-field-state"
	ErrorHandoffProjectionUnavailable ErrorCode = "handoff-projection-unavailable"
	ErrorHandoffIneligible             ErrorCode = "handoff-ineligible"
	ErrorHandoffScopeMismatch          ErrorCode = "handoff-scope-mismatch"
	ErrorHandoffUnknownParameter       ErrorCode = "handoff-unknown-parameter"
	ErrorHandoffDuplicateParameter     ErrorCode = "handoff-duplicate-parameter"
)

type AssetKind string

const (
	AssetSavedSearch AssetKind = "saved-search"
	AssetQueryAsset  AssetKind = "query-asset"
)

type AccessDecision string

const (
	AccessAllow AccessDecision = "allow"
	AccessDeny  AccessDecision = "deny"
)

type ValidationState string

const (
	ValidationValidated    ValidationState = "validated"
	ValidationStale        ValidationState = "stale"
	ValidationIncompatible ValidationState = "incompatible"
)

type SourceState string

const (
	SourceAvailable SourceState = "available"
	SourceStale     SourceState = "stale"
	SourceMissing   SourceState = "missing"
)

type FieldState string

const (
	FieldAvailable FieldState = "available"
	FieldStale     FieldState = "stale"
	FieldMissing   FieldState = "missing"
	FieldRemoved   FieldState = "removed"
)

type Field struct {
	Name  string     `json:"name"`
	State FieldState `json:"state"`
}

type Source struct {
	ID        string      `json:"id"`
	TenantRef string      `json:"tenant_ref"`
	State     SourceState `json:"state"`
	Fields    []Field     `json:"fields"`
}

type Input struct {
	TenantRef                   string         `json:"tenant_ref"`
	AuthorizedTenants           []string       `json:"authorized_tenants"`
	EnvironmentRef              string         `json:"environment_ref"`
	AssetID                     string         `json:"asset_id"`
	AssetKind                   AssetKind      `json:"asset_kind"`
	AssetAccessDecision         AccessDecision `json:"asset_access_decision"`
	QueryRef                    string         `json:"query_ref"`
	QueryTenantRef              string         `json:"query_tenant_ref"`
	QueryVersion                string         `json:"query_version"`
	QueryAccessDecision         AccessDecision `json:"query_access_decision"`
	AuthorRef                   string         `json:"author_ref"`
	ValidationState             ValidationState `json:"validation_state"`
	ParameterNames              []string       `json:"parameter_names"`
	Sources                     []Source       `json:"sources"`
	LineageRef                  string         `json:"lineage_ref"`
	DeprecationReason           string         `json:"deprecation_reason"`
	ReplacementAssetRef         string         `json:"replacement_asset_ref"`
	MutationRequested           string         `json:"mutation_requested"`
	SearchJobExecutionRequested bool           `json:"search_job_execution_requested"`
}

type ProjectedField struct {
	Name  string
	State FieldState
}

type ProjectedSource struct {
	ID        string
	TenantRef string
	State     SourceState
	Fields    []ProjectedField
}

type Projection struct {
	TenantRef           string
	EnvironmentRef      string
	AssetID             string
	AssetKind           AssetKind
	QueryRef            string
	QueryTenantRef      string
	QueryVersion        string
	AuthorRef           string
	ValidationState     ValidationState
	ParameterNames      []string
	Sources             []ProjectedSource
	LineageRef          string
	DeprecationReason   string
	ReplacementAssetRef string
	Compatibility       string
	Diagnostics         []string
	HandoffEligible     bool
}

type Result struct {
	Allowed             bool
	ErrorCode           ErrorCode
	Compatibility       string
	Diagnostics         []string
	HandoffEligible     bool
	ImmutableProjection bool
	DeepCopyIsolated    bool
	MutationExecuted    bool
	SearchJobCreated    bool
	SearchJobExecuted   bool
	AuditRequired       bool
	Projection          *Projection
}

type ParameterValue struct {
	Name        string
	OpaqueValue []byte
}

type HandoffSourceContext struct {
	ID             string
	RequiredFields []string
}

type ReturnContextToken struct {
	AssetID   string
	AssetKind AssetKind
}

type HandoffInput struct {
	TenantRef       string
	EnvironmentRef  string
	ParameterValues []ParameterValue
}

type ExecutionHandoffDraft struct {
	TenantRef                      string
	EnvironmentRef                 string
	QueryRef                       string
	QueryTenantRef                 string
	QueryVersion                   string
	ParameterNames                 []string
	ParameterValues                []ParameterValue
	Sources                        []HandoffSourceContext
	ReturnContext                  ReturnContextToken
	SourceReevaluationRequired     bool
	PermissionReevaluationRequired bool
	OriginalAssetImmutable         bool
}

type HandoffResult struct {
	Allowed            bool
	ErrorCode          ErrorCode
	AuditRequired      bool
	SourceAssetMutated bool
	SearchJobCreated   bool
	SearchJobExecuted  bool
	Draft              *ExecutionHandoffDraft
}
