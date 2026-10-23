package studiocatalog

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// AssetKind defines canonical Studio asset categories.
type AssetKind string

const (
	AssetKindSkill           AssetKind = "skill"
	AssetKindWorkflow        AssetKind = "workflow"
	AssetKindAutomationAgent AssetKind = "automation-agent"
	AssetKindAgentTeam       AssetKind = "agent-team"
	AssetKindTool            AssetKind = "tool"
)

// AllAssetKinds provides the list of valid asset kinds.
var AllAssetKinds = []AssetKind{
	AssetKindSkill,
	AssetKindWorkflow,
	AssetKindAutomationAgent,
	AssetKindAgentTeam,
	AssetKindTool,
}

// IsValidAssetKind checks whether a string is a recognized asset kind.
func IsValidAssetKind(k AssetKind) bool {
	for _, kind := range AllAssetKinds {
		if k == kind {
			return true
		}
	}
	return false
}

// LifecycleState represents canonical lifecycle status.
type LifecycleState string

const (
	LifecycleStateDraft      LifecycleState = "draft"
	LifecycleStateReview     LifecycleState = "review"
	LifecycleStateApproved   LifecycleState = "approved"
	LifecycleStatePublished  LifecycleState = "published"
	LifecycleStateDeprecated LifecycleState = "deprecated"
	LifecycleStateArchived   LifecycleState = "archived"
	LifecycleStateRetired    LifecycleState = "retired"
	LifecycleStateSuperseded LifecycleState = "superseded"
)

// FunctionalState represents runtime projection health status.
type FunctionalState string

const (
	FunctionalStateReady              FunctionalState = "ready"
	FunctionalStatePartial            FunctionalState = "partial"
	FunctionalStateStale              FunctionalState = "stale"
	FunctionalStatePermissionFiltered FunctionalState = "permission-filtered"
)

// RiskClassification defines operational risk level.
type RiskClassification string

const (
	RiskClassificationLow      RiskClassification = "low"
	RiskClassificationMedium   RiskClassification = "medium"
	RiskClassificationHigh     RiskClassification = "high"
	RiskClassificationCritical RiskClassification = "critical"
)

// Owner represents asset ownership with strict fallback under REQ-OBJ-009.
type Owner struct {
	OwnerID     string `json:"owner_id"`
	OwnerName   string `json:"owner_name"`
	OwnerType   string `json:"owner_type"` // "user", "team", "system", "unknown"
	IsUnknown   bool   `json:"is_unknown"`
	DisplayName string `json:"display_name"`
}

// Normalize ensures that missing or blank owners default to "owner unknown" without fabrication.
func (o *Owner) Normalize() {
	if strings.TrimSpace(o.OwnerID) == "" || strings.TrimSpace(o.OwnerName) == "" || o.IsUnknown {
		o.OwnerID = "unknown"
		o.OwnerName = "owner unknown"
		o.OwnerType = "unknown"
		o.IsUnknown = true
		o.DisplayName = "owner unknown"
	} else if o.DisplayName == "" {
		o.DisplayName = o.OwnerName
	}
}

// VersionInfo represents immutable semantic version metadata.
type VersionInfo struct {
	VersionID   string    `json:"version_id"`
	SemVer      string    `json:"semver"`
	CreatedAt   time.Time `json:"created_at"`
	IsImmutable bool      `json:"is_immutable"`
}

// Parameter defines inputs/outputs for executable specifications.
type Parameter struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	Description  string `json:"description,omitempty"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// Dependency represents external dependencies and secret references.
type Dependency struct {
	Kind              string `json:"kind"` // "tool", "skill", "secret_reference"
	RefID             string `json:"ref_id"`
	VersionConstraint string `json:"version_constraint,omitempty"`
	IsOptional        bool   `json:"is_optional"`
	HealthStatus      string `json:"health_status,omitempty"` // "available", "stale", "missing"
	Purpose           string `json:"purpose,omitempty"`
	IsMasked          bool   `json:"is_masked,omitempty"`
}

// Consumer represents registered consumers of an asset.
type Consumer struct {
	ConsumerType string `json:"consumer_type"` // "workflow", "automation-agent"
	ConsumerRef  string `json:"consumer_ref"`
}

// Supersession records historical lineage links.
type Supersession struct {
	IsSuperseded bool   `json:"is_superseded"`
	SupersededBy string `json:"superseded_by,omitempty"`
	Supersedes   string `json:"supersedes,omitempty"`
}

// Provenance tracks deep lineage and creation origins.
type Provenance struct {
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	ParentID        string    `json:"parent_id,omitempty"`
	ClonedFrom      string    `json:"cloned_from,omitempty"`
	ProvenanceChain []string  `json:"provenance_chain,omitempty"`
}

// Asset is the complete canonical Studio asset entity.
type Asset struct {
	ID                 string             `json:"id"`
	TenantID           string             `json:"tenant_id"`
	EnvironmentID      string             `json:"environment_id,omitempty"`
	Type               AssetKind          `json:"type"`
	Name               string             `json:"name"`
	DisplayName        string             `json:"display_name"`
	Description        string             `json:"description"`
	Tags               []string           `json:"tags"`
	Owner              Owner              `json:"owner"`
	Contributors       []string           `json:"contributors"`
	Version            VersionInfo        `json:"version"`
	LifecycleState     LifecycleState     `json:"lifecycle_state"`
	FunctionalState    FunctionalState    `json:"functional_state"`
	ExecutionClass     string             `json:"execution_class"`
	RiskClassification RiskClassification `json:"risk_classification"`
	IsAIPowered        bool               `json:"is_ai_powered"`
	HasNonAIFallback   bool               `json:"has_non_ai_fallback"`
	Parameters         []Parameter        `json:"parameters"`
	Dependencies       []Dependency       `json:"dependencies"`
	Consumers          []Consumer         `json:"consumers"`
	Supersession       Supersession       `json:"supersession"`
	Provenance         Provenance         `json:"provenance"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// Clone creates an immutable deep copy of an Asset.
func (a *Asset) Clone() *Asset {
	if a == nil {
		return nil
	}
	cp := *a
	cp.Owner.Normalize()

	if a.Tags != nil {
		cp.Tags = make([]string, len(a.Tags))
		copy(cp.Tags, a.Tags)
	}
	if a.Contributors != nil {
		cp.Contributors = make([]string, len(a.Contributors))
		copy(cp.Contributors, a.Contributors)
	}
	if a.Parameters != nil {
		cp.Parameters = make([]Parameter, len(a.Parameters))
		copy(cp.Parameters, a.Parameters)
	}
	if a.Dependencies != nil {
		cp.Dependencies = make([]Dependency, len(a.Dependencies))
		copy(cp.Dependencies, a.Dependencies)
	}
	if a.Consumers != nil {
		cp.Consumers = make([]Consumer, len(a.Consumers))
		copy(cp.Consumers, a.Consumers)
	}
	if a.Provenance.ProvenanceChain != nil {
		cp.Provenance.ProvenanceChain = make([]string, len(a.Provenance.ProvenanceChain))
		copy(cp.Provenance.ProvenanceChain, a.Provenance.ProvenanceChain)
	}
	return &cp
}

// ToSummary projects a full asset into a catalog summary item.
func (a *Asset) ToSummary() AssetSummary {
	owner := a.Owner
	owner.Normalize()
	tagsCopy := make([]string, len(a.Tags))
	copy(tagsCopy, a.Tags)

	return AssetSummary{
		ID:                 a.ID,
		TenantID:           a.TenantID,
		Type:               a.Type,
		Name:               a.Name,
		DisplayName:        a.DisplayName,
		Description:        a.Description,
		Version:            a.Version.SemVer,
		LifecycleState:     a.LifecycleState,
		FunctionalState:    a.FunctionalState,
		RiskClassification: a.RiskClassification,
		IsAIPowered:        a.IsAIPowered,
		HasNonAIFallback:   a.HasNonAIFallback,
		Owner:              owner,
		Tags:               tagsCopy,
		UpdatedAt:          a.UpdatedAt,
	}
}

// AssetSummary is the lightweight projection for browse/search/filter listings.
type AssetSummary struct {
	ID                 string             `json:"id"`
	TenantID           string             `json:"tenant_id"`
	Type               AssetKind          `json:"type"`
	Name               string             `json:"name"`
	DisplayName        string             `json:"display_name"`
	Description        string             `json:"description"`
	Version            string             `json:"version"`
	LifecycleState     LifecycleState     `json:"lifecycle_state"`
	FunctionalState    FunctionalState    `json:"functional_state"`
	RiskClassification RiskClassification `json:"risk_classification"`
	IsAIPowered        bool               `json:"is_ai_powered"`
	HasNonAIFallback   bool               `json:"has_non_ai_fallback"`
	Owner              Owner              `json:"owner"`
	Tags               []string           `json:"tags"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// PaginationInfo details pagination state.
type PaginationInfo struct {
	Limit        int    `json:"limit"`
	HasMore      bool   `json:"has_more"`
	CursorNext   string `json:"cursor_next,omitempty"`
	CursorPrev   string `json:"cursor_prev,omitempty"`
	TotalCount   int    `json:"total_count"`
	IsTotalExact bool   `json:"is_total_exact"`
}

// CursorPayload represents decoded internal cursor structure.
type CursorPayload struct {
	Version       int    `json:"v"`
	TenantID      string `json:"t"`
	SortField     string `json:"sf"`
	SortDirection string `json:"sd"` // "asc" or "desc"
	LastValue     string `json:"lv"`
	TiebreakerID  string `json:"id"`
	Fingerprint   string `json:"fp"`
}

// EncodeCursor serializes cursor payload into an opaque URL-safe base64 string.
func EncodeCursor(p CursorPayload) string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeCursor parses an opaque cursor string.
func DecodeCursor(s string) (*CursorPayload, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encoding: %w", err)
	}
	var p CursorPayload
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("corrupted cursor payload: %w", err)
	}
	if p.Version != 1 || p.TiebreakerID == "" {
		return nil, fmt.Errorf("invalid cursor schema")
	}
	return &p, nil
}

// Request and Response types for the 8 Class 0 Operations:

// BrowseRequest specifies parameters for browse query.
type BrowseRequest struct {
	TenantRef         string    `json:"tenant_ref"`
	AuthorizedTenants []string  `json:"authorized_tenants,omitempty"`
	Kind              AssetKind `json:"kind,omitempty"`
	Limit             int       `json:"limit,omitempty"`
	Cursor            string    `json:"cursor,omitempty"`
	SortField         string    `json:"sort_field,omitempty"`     // "updated_at", "name", "id", "created_at"
	SortDirection     string    `json:"sort_direction,omitempty"` // "asc" or "desc"
	CorrelationID     string    `json:"correlation_id,omitempty"`
}

// BrowseResponse returns paginated asset summaries.
type BrowseResponse struct {
	Items      []AssetSummary `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// SearchRequest specifies text search parameters.
type SearchRequest struct {
	TenantRef              string    `json:"tenant_ref"`
	AuthorizedTenants      []string  `json:"authorized_tenants,omitempty"`
	QueryText              string    `json:"query_text"`
	Kind                   AssetKind `json:"kind,omitempty"`
	Limit                  int       `json:"limit,omitempty"`
	Cursor                 string    `json:"cursor,omitempty"`
	SharedSearchAvailable  bool      `json:"shared_search_available"` // If false, fallback to local search with "partial" state
	CorrelationID          string    `json:"correlation_id,omitempty"`
}

// SearchResponse returns matching asset summaries.
type SearchResponse struct {
	Items      []AssetSummary `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// FilterCriteria specifies multi-attribute predicates.
type FilterCriteria struct {
	Kinds            []AssetKind          `json:"kinds,omitempty"`
	LifecycleStates  []LifecycleState     `json:"lifecycle_states,omitempty"`
	FunctionalStates []FunctionalState    `json:"functional_states,omitempty"`
	OwnerID          string               `json:"owner_id,omitempty"`
	RiskLevels       []RiskClassification `json:"risk_levels,omitempty"`
	IsAIPowered      *bool                `json:"is_ai_powered,omitempty"`
	Tags             []string             `json:"tags,omitempty"`
}

// FilterRequest specifies filtering parameters.
type FilterRequest struct {
	TenantRef         string         `json:"tenant_ref"`
	AuthorizedTenants []string       `json:"authorized_tenants,omitempty"`
	Criteria          FilterCriteria `json:"criteria"`
	Limit             int            `json:"limit,omitempty"`
	Cursor            string         `json:"cursor,omitempty"`
	SortField         string         `json:"sort_field,omitempty"`
	SortDirection     string         `json:"sort_direction,omitempty"`
	CorrelationID     string         `json:"correlation_id,omitempty"`
}

// FilterResponse returns filtered asset summaries.
type FilterResponse struct {
	Items      []AssetSummary `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

// GroupByFacet defines supported grouping dimensions.
type GroupByFacet string

const (
	FacetKind               GroupByFacet = "by_kind"
	FacetLifecycleState     GroupByFacet = "by_lifecycle_state"
	FacetOwner              GroupByFacet = "by_owner"
	FacetRiskClassification GroupByFacet = "by_risk_classification"
)

// GroupBucket represents aggregated category counts.
type GroupBucket struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// GroupRequest specifies grouping parameters.
type GroupRequest struct {
	TenantRef         string          `json:"tenant_ref"`
	AuthorizedTenants []string        `json:"authorized_tenants,omitempty"`
	GroupBy           GroupByFacet    `json:"group_by"`
	Criteria          *FilterCriteria `json:"criteria,omitempty"`
	CorrelationID     string          `json:"correlation_id,omitempty"`
}

// GroupResponse returns aggregated bucket statistics.
type GroupResponse struct {
	Facet      GroupByFacet  `json:"facet"`
	Groups     []GroupBucket `json:"groups"`
	TotalItems int           `json:"total_items"`
}

// OpenRequest specifies parameters to open an asset in authoring workspace.
type OpenRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	VersionID         string   `json:"version_id,omitempty"`
	ReturnOrigin      string   `json:"return_origin"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// OpenResponse provides client navigation descriptors preserving origin.
type OpenResponse struct {
	TargetSurface string            `json:"target_surface"`
	SurfaceRoute  string            `json:"surface_route"`
	AssetID       string            `json:"asset_id"`
	VersionID     string            `json:"version_id,omitempty"`
	ReturnOrigin  string            `json:"return_origin"`
	ContextTokens map[string]string `json:"context_tokens"`
}

// InspectRequest specifies parameters for deep schema inspection.
type InspectRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	VersionID         string   `json:"version_id,omitempty"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// InspectResponse returns full metadata or sanitized permission-filtered stub.
type InspectResponse struct {
	Asset *Asset `json:"asset"`
}

// CompareRequest specifies two assets or versions to compare.
type CompareRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	BaseAssetID       string   `json:"base_asset_id"`
	TargetAssetID     string   `json:"target_asset_id"`
	BaseVersionID     string   `json:"base_version_id,omitempty"`
	TargetVersionID   string   `json:"target_version_id,omitempty"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// DependencyChange details changed dependency constraints.
type DependencyChange struct {
	RefID         string `json:"ref_id"`
	OldConstraint string `json:"old_constraint"`
	NewConstraint string `json:"new_constraint"`
}

// StructuralDiff details side-by-side differences.
type StructuralDiff struct {
	ParametersAdded            []Parameter        `json:"parameters_added"`
	ParametersRemoved          []Parameter        `json:"parameters_removed"`
	ParametersModified         []Parameter        `json:"parameters_modified"`
	DependenciesChanged        []DependencyChange `json:"dependencies_changed"`
	RiskClassificationChanged  bool               `json:"risk_classification_changed"`
	OldRisk                    RiskClassification `json:"old_risk"`
	NewRisk                    RiskClassification `json:"new_risk"`
	AIFlagChanged              bool               `json:"ai_flag_changed"`
	LifecycleChanged           bool               `json:"lifecycle_changed"`
	OldLifecycle               LifecycleState     `json:"old_lifecycle"`
	NewLifecycle               LifecycleState     `json:"new_lifecycle"`
}

// CompareResponse returns structural diff between two assets.
type CompareResponse struct {
	Base  AssetSummary   `json:"base"`
	Target AssetSummary   `json:"target"`
	Diff  StructuralDiff `json:"diff"`
}

// ClonePreparationRequest specifies clone preparation parameters.
type ClonePreparationRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	TargetTenantRef   string   `json:"target_tenant_ref"` // Must match TenantRef (cross-tenant clone is strictly forbidden!)
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	SourceAssetID     string   `json:"source_asset_id"`
	SourceVersionID   string   `json:"source_version_id,omitempty"`
	ClonedBy          string   `json:"cloned_by"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// DraftCloneDescriptor contains prepared draft clone metadata.
type DraftCloneDescriptor struct {
	DraftAssetID      string    `json:"draft_asset_id"`
	ClonedFromID      string    `json:"cloned_from_id"`
	SourceVersion     string    `json:"source_version"`
	SourceTenantID    string    `json:"source_tenant_id"`
	TargetTenantID    string    `json:"target_tenant_id"`
	ClonedBy          string    `json:"cloned_by"`
	ClonedAt          time.Time `json:"cloned_at"`
	ProvenanceChain   []string  `json:"provenance_chain"`
	PrefilledMetadata *Asset    `json:"prefilled_metadata"`
}

// ClonePreparationResponse returns prepared draft clone descriptor.
type ClonePreparationResponse struct {
	DraftClone DraftCloneDescriptor `json:"draft_clone"`
}
