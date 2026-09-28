// Package huntmanagement implements the bounded read-only Hunt workspace projection.
package huntmanagement

type ErrorCode string

const (
	ErrorNone                     ErrorCode = ""
	ErrorMutationOpenDecision     ErrorCode = "mutation-open-decision"
	ErrorMissingTenant            ErrorCode = "missing-tenant"
	ErrorWildcardTenant           ErrorCode = "wildcard-tenant"
	ErrorCrossTenant              ErrorCode = "cross-tenant"
	ErrorMissingEnvironment       ErrorCode = "missing-environment"
	ErrorMissingQuestion          ErrorCode = "missing-question"
	ErrorMissingScope             ErrorCode = "missing-scope"
	ErrorInvalidTimeRange         ErrorCode = "invalid-time-range"
	ErrorMissingOwner             ErrorCode = "missing-owner"
	ErrorInvalidContributor       ErrorCode = "invalid-contributor"
	ErrorInvalidReference         ErrorCode = "invalid-reference"
	ErrorCrossTenantReference     ErrorCode = "cross-tenant-reference"
	ErrorReferenceAccessDenied    ErrorCode = "reference-access-denied"
	ErrorInvalidReferenceState    ErrorCode = "invalid-reference-state"
	ErrorSearchJobQueryNotVisible ErrorCode = "search-job-query-not-visible"
)

type AccessDecision string

const (
	AccessAllow AccessDecision = "allow"
	AccessDeny  AccessDecision = "deny"
)

type ReferenceKind string

const (
	ReferenceQuery      ReferenceKind = "query"
	ReferenceSearchJob  ReferenceKind = "search-job"
	ReferenceHypothesis ReferenceKind = "hypothesis"
	ReferenceCase       ReferenceKind = "case"
	ReferenceIncident   ReferenceKind = "incident"
)

type ReferenceState string

const (
	ReferenceAvailable   ReferenceState = "available"
	ReferenceStale       ReferenceState = "stale"
	ReferencePartial     ReferenceState = "partial"
	ReferenceUnavailable ReferenceState = "unavailable"
)

type Reference struct {
	Kind           ReferenceKind   `json:"kind"`
	ID             string          `json:"id"`
	TenantRef      string          `json:"tenant_ref"`
	AccessDecision AccessDecision  `json:"access_decision"`
	QueryRef       string          `json:"query_ref"`
	State          ReferenceState  `json:"state"`
	Version        string          `json:"version,omitempty"`
	ProvenanceRef  string          `json:"provenance_ref,omitempty"`
}

type Input struct {
	TenantRef         string      `json:"tenant_ref"`
	AuthorizedTenants []string    `json:"authorized_tenants"`
	EnvironmentRef    string      `json:"environment_ref"`
	Question          string      `json:"question"`
	Scope             string      `json:"scope"`
	TimeStart         string      `json:"time_start"`
	TimeEnd           string      `json:"time_end"`
	OwnerRef          string      `json:"owner_ref"`
	ContributorRefs   []string    `json:"contributor_refs"`
	References        []Reference `json:"references"`
	MutationRequested string      `json:"mutation_requested"`
	CorrelationID     string      `json:"correlation_id,omitempty"`
}

type TimeRange struct {
	Start string
	End   string
}

type ProjectedReference struct {
	Kind           ReferenceKind
	ID             string
	TenantRef      string
	QueryRef       string
	State          ReferenceState
	CanonicalOwner string
	Version        string
	ProvenanceRef  string
}

type Projection struct {
	TenantRef       string
	EnvironmentRef  string
	Question        string
	Scope           string
	TimeRange       TimeRange
	OwnerRef        string
	ContributorRefs []string
	References      []ProjectedReference
	Limitations     []string
	Complete        bool
	CorrelationID   string
}

type Result struct {
	Allowed          bool
	ErrorCode        ErrorCode
	AuditRequired    bool
	MutationExecuted bool
	Projection       *Projection
}
