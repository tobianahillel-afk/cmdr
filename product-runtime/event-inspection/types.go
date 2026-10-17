// Package eventinspection implements the bounded read-only Event inspection projection.
package eventinspection

type ErrorCode string

const (
	ErrorNone                        ErrorCode = ""
	ErrorMutationExcluded            ErrorCode = "mutation-excluded"
	ErrorUnsupportedMutation         ErrorCode = "unsupported-mutation"
	ErrorMissingTenant               ErrorCode = "missing-tenant"
	ErrorWildcardTenant              ErrorCode = "wildcard-tenant"
	ErrorCrossTenant                 ErrorCode = "cross-tenant"
	ErrorMissingEnvironment          ErrorCode = "missing-environment"
	ErrorMissingEventRef             ErrorCode = "missing-event-ref"
	ErrorEventTenantMismatch         ErrorCode = "event-tenant-mismatch"
	ErrorPermissionDenied            ErrorCode = "permission-denied"
	ErrorInvalidAccessDecision       ErrorCode = "invalid-access-decision"
	ErrorProjectionDenied            ErrorCode = "projection-denied"
	ErrorInvalidEnrichmentProvenance ErrorCode = "invalid-enrichment-provenance"
	ErrorDerivedAsSource             ErrorCode = "derived-as-source"
	ErrorInvalidCorrelationID        ErrorCode = "invalid-correlation-id"
	ErrorInvalidEventState           ErrorCode = "invalid-event-state"
	ErrorInvalidPivotContext         ErrorCode = "invalid-pivot-context"
	ErrorRestrictedPivotValue        ErrorCode = "restricted-pivot-value"
)

const PermissionTelemetryRead = "perm.shared-capabilities.telemetry-event.read"

type AccessDecision string

const (
	AccessAllow AccessDecision = "allow"
	AccessDeny  AccessDecision = "deny"
)

type EventState string

const (
	StateAvailable         EventState = "available"
	StateSourceUnavailable EventState = "source-unavailable"
	StateTombstone         EventState = "tombstone"
	StatePartial           EventState = "partial"
	StateRestricted        EventState = "restricted"
)

type SourceMetadata struct {
	SourceRef     string
	DataSourceRef string
	Kind          string
}

type ParserMetadata struct {
	ParserRef string
	Version   string
	Status    string
}

type EventEnrichment struct {
	Producer          string
	Version           string
	Freshness         string
	PresentedAsSource bool
	Values            map[string]string
}

type ProjectionEnrichment struct {
	Producer  string
	Version   string
	Freshness string
	Stale     bool
	Values    map[string]string
}

type Event struct {
	Ref                 string
	TenantRef           string
	State               EventState
	Source              SourceMetadata
	Parser              ParserMetadata
	SourceFields        map[string]string
	NormalizedFields    map[string]string
	MissingFields       []string
	RawPayload          []byte
	RenderedPayload     []byte
	RawDerivedSensitive map[string]string
	Enrichments         []EventEnrichment
}

type Input struct {
	TenantRef              string
	AuthorizedTenants      []string
	EnvironmentRef         string
	Permissions            []string
	RawAccessDecision      AccessDecision
	RenderedAccessDecision AccessDecision
	CorrelationID          string
	MutationRequested      string
	Event                  Event
}

type Projection struct {
	TenantRef           string
	EnvironmentRef      string
	EventRef            string
	CorrelationID       string
	State               EventState
	Source              SourceMetadata
	Parser              ParserMetadata
	SourceFields        map[string]string
	NormalizedFields    map[string]string
	MissingFields       []string
	RawPayload          []byte
	RenderedPayload     []byte
	RawDerivedSensitive map[string]string
	Enrichments         []ProjectionEnrichment
}

type Result struct {
	Allowed       bool
	ErrorCode     ErrorCode
	AuditRequired bool
	Projection    *Projection
}


const PivotTargetCapability = "CAP-INV-002"

type PivotSelection struct {
	Field           string
	Value           string
	ValuePresent    bool
	ValueAuthorized bool
	TimeStart       string
	TimeEnd         string
	ReturnOrigin    string
}

type PivotTimeRange struct {
	Start string
	End   string
}

type PivotReturnContext struct {
	TenantRef     string
	EnvironmentRef string
	EventRef      string
	CorrelationID string
	Origin        string
}

type PivotDraft struct {
	TargetCapability string
	TenantRef        string
	EnvironmentRef   string
	EventRef         string
	CorrelationID    string
	Field            string
	Value            string
	FieldOrigin      string
	SourceRef        string
	DataSourceRef    string
	TimeRange        PivotTimeRange
	ReturnContext    PivotReturnContext
}

type PivotResult struct {
	Allowed       bool
	ErrorCode     ErrorCode
	AuditRequired bool
	Draft         *PivotDraft
}
