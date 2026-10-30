package studiocatalog

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetadataService defines the Class 0 metadata inspection and Class 2 draft metadata management operations (CAP-STD-002).
type MetadataService interface {
	ValidateMetadata(ctx context.Context, auth AuthContext, req ValidateMetadataRequest) (*ValidateMetadataResponse, error)
	SuggestClassificationAndTags(ctx context.Context, auth AuthContext, req SuggestTagsRequest) (*SuggestTagsResponse, error)
	AssignOwner(ctx context.Context, auth AuthContext, req AssignOwnerRequest) (*AssignOwnerResponse, error)
	UpdateDraftMetadata(ctx context.Context, auth AuthContext, req UpdateDraftMetadataRequest) (*UpdateDraftMetadataResponse, error)
	RecordSupersession(ctx context.Context, auth AuthContext, req RecordSupersessionRequest) (*RecordSupersessionResponse, error)
	GetSupersessionHistory(ctx context.Context, auth AuthContext, req SupersessionHistoryRequest) (*SupersessionHistoryResponse, error)
}

// ValidationSeverity defines severity levels for metadata validation findings.
type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
	SeverityInfo    ValidationSeverity = "info"
)

// ValidationIssue represents a single validation finding on an asset.
type ValidationIssue struct {
	Field    string             `json:"field"`
	Message  string             `json:"message"`
	Severity ValidationSeverity `json:"severity"`
	Code     string             `json:"code"`
}

// MetadataValidationResult details the validation evaluation of an asset's metadata.
type MetadataValidationResult struct {
	Valid                 bool              `json:"valid"`
	State                 FunctionalState   `json:"state"`
	Issues                []ValidationIssue `json:"issues"`
	MissingRequiredFields []string          `json:"missing_required_fields"`
	ValidatedAt           time.Time         `json:"validated_at"`
}

// TagSuggestionResult provides deterministic tag and risk classification suggestions without requiring AI (REQ-PROD-016).
type TagSuggestionResult struct {
	SuggestedTags         []string           `json:"suggested_tags"`
	SuggestedRisk         RiskClassification `json:"suggested_risk"`
	SuggestedDeliveryMode DeliveryMode       `json:"suggested_delivery_mode"`
	Rationale             []string           `json:"rationale"`
	IsAIPowered           bool               `json:"is_ai_powered"` // Strictly false per REQ-PROD-016
}

// SupersessionHistoryNode represents a single linked node in an asset's supersession lineage.
type SupersessionHistoryNode struct {
	AssetID        string         `json:"asset_id"`
	Version        string         `json:"version"`
	LifecycleState LifecycleState `json:"lifecycle_state"`
	IsCurrent      bool           `json:"is_current"`
	IsActive       bool           `json:"is_active"`
	Supersedes     string         `json:"supersedes,omitempty"`
	SupersededBy   string         `json:"superseded_by,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Owner          Owner          `json:"owner"`
	Summary        AssetSummary   `json:"summary"`
}

// ValidateMetadataRequest defines parameters for metadata validation inspection.
type ValidateMetadataRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// ValidateMetadataResponse returns validation results.
type ValidateMetadataResponse struct {
	AssetID    string                   `json:"asset_id"`
	Validation MetadataValidationResult `json:"validation"`
}

// SuggestTagsRequest defines parameters for deterministic tag and classification suggestion.
type SuggestTagsRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// SuggestTagsResponse returns deterministic tag suggestions.
type SuggestTagsResponse struct {
	AssetID     string              `json:"asset_id"`
	Suggestions TagSuggestionResult `json:"suggestions"`
}

// AssignOwnerRequest defines parameters for assigning or updating an asset's owner (Class 2).
type AssignOwnerRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	NewOwner          Owner    `json:"new_owner"`
	Actor             string   `json:"actor"`
	Reason            string   `json:"reason"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// AssignOwnerResponse returns the updated asset after owner assignment.
type AssignOwnerResponse struct {
	Asset         *Asset    `json:"asset"`
	PreviousOwner Owner     `json:"previous_owner"`
	AssignedAt    time.Time `json:"assigned_at"`
}

// UpdateDraftMetadataRequest defines parameters for updating draft metadata (Class 2).
type UpdateDraftMetadataRequest struct {
	TenantRef          string             `json:"tenant_ref"`
	AuthorizedTenants  []string           `json:"authorized_tenants,omitempty"`
	AssetID            string             `json:"asset_id"`
	DisplayName        string             `json:"display_name,omitempty"`
	Description        string             `json:"description,omitempty"`
	Tags               []string           `json:"tags,omitempty"`
	RiskClassification RiskClassification `json:"risk_classification,omitempty"`
	DeliveryMode       DeliveryMode       `json:"delivery_mode,omitempty"`
	ExecutionClass     string             `json:"execution_class,omitempty"`
	IsAIPowered        *bool              `json:"is_ai_powered,omitempty"`
	HasNonAIFallback   *bool              `json:"has_non_ai_fallback,omitempty"`
	Actor              string             `json:"actor"`
	Reason             string             `json:"reason"`
	CorrelationID      string             `json:"correlation_id,omitempty"`
}

// UpdateDraftMetadataResponse returns updated asset and validation status.
type UpdateDraftMetadataResponse struct {
	Asset      *Asset                   `json:"asset"`
	Validation MetadataValidationResult `json:"validation"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

// RecordSupersessionRequest defines parameters to link an older asset to a newer superseding asset (Class 2).
type RecordSupersessionRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	OlderAssetID      string   `json:"older_asset_id"`
	NewerAssetID      string   `json:"newer_asset_id"`
	Actor             string   `json:"actor"`
	Reason            string   `json:"reason"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// RecordSupersessionResponse returns updated asset summaries after supersession.
type RecordSupersessionResponse struct {
	OlderAsset AssetSummary `json:"older_asset"`
	NewerAsset AssetSummary `json:"newer_asset"`
	RecordedAt time.Time    `json:"recorded_at"`
}

// SupersessionHistoryRequest defines parameters to navigate historical supersession lineage.
type SupersessionHistoryRequest struct {
	TenantRef         string   `json:"tenant_ref"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
	AssetID           string   `json:"asset_id"`
	CorrelationID     string   `json:"correlation_id,omitempty"`
}

// SupersessionHistoryResponse returns complete ordered lineage from root ancestor to latest active.
type SupersessionHistoryResponse struct {
	RootAssetID    string                    `json:"root_asset_id"`
	LatestAssetID  string                    `json:"latest_asset_id"`
	CurrentAssetID string                    `json:"current_asset_id"`
	Chain          []SupersessionHistoryNode `json:"chain"`
}

// ValidateAssetMetadata performs deterministic validation of asset metadata without requiring AI (REQ-PROD-016).
// It verifies required fields, checks owner normalization, validates classifications, and determines functional state.
func ValidateAssetMetadata(asset *Asset) MetadataValidationResult {
	now := time.Now().UTC()
	if asset == nil {
		return MetadataValidationResult{
			Valid:                 false,
			State:                 FunctionalStateDraft,
			Issues:                []ValidationIssue{{Field: "asset", Message: "Asset is nil", Severity: SeverityError, Code: "nil_asset"}},
			MissingRequiredFields: []string{"asset"},
			ValidatedAt:           now,
		}
	}

	var issues []ValidationIssue
	var missingRequired []string

	// 1. Identity validation
	if strings.TrimSpace(asset.ID) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "id",
			Message:  "Asset identifier is required and cannot be empty",
			Severity: SeverityError,
			Code:     "missing_id",
		})
		missingRequired = append(missingRequired, "id")
	}

	// 2. Tenant validation
	if strings.TrimSpace(asset.TenantID) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "tenant_id",
			Message:  "Tenant identifier is required and cannot be empty",
			Severity: SeverityError,
			Code:     "missing_tenant_id",
		})
		missingRequired = append(missingRequired, "tenant_id")
	}

	// 3. Type validation
	if !IsValidAssetKind(asset.Type) {
		issues = append(issues, ValidationIssue{
			Field:    "type",
			Message:  fmt.Sprintf("Invalid or unassigned asset kind %q", asset.Type),
			Severity: SeverityError,
			Code:     "invalid_type",
		})
		missingRequired = append(missingRequired, "type")
	}

	// 4. Name validation
	if strings.TrimSpace(asset.Name) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "name",
			Message:  "Asset name is required and cannot be empty",
			Severity: SeverityError,
			Code:     "missing_name",
		})
		missingRequired = append(missingRequired, "name")
	}

	// 5. Version validation
	if strings.TrimSpace(asset.Version.SemVer) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "version.semver",
			Message:  "Semantic version is required and cannot be empty",
			Severity: SeverityError,
			Code:     "missing_semver",
		})
		missingRequired = append(missingRequired, "version.semver")
	}

	// 6. Owner validation (REQ-OBJ-009)
	// Missing or unrecognized owner must normalize to explicit "owner unknown" without fabrication.
	if asset.Owner.IsUnknown || strings.TrimSpace(asset.Owner.OwnerID) == "" || strings.TrimSpace(asset.Owner.OwnerName) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "owner",
			Message:  "Asset owner is unknown or absent; normalized to explicit 'owner unknown' per REQ-OBJ-009",
			Severity: SeverityInfo,
			Code:     "owner_unknown",
		})
	}

	// 7. Risk classification check
	if strings.TrimSpace(string(asset.RiskClassification)) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "risk_classification",
			Message:  "Risk classification is not specified; defaulting to partial state",
			Severity: SeverityWarning,
			Code:     "missing_risk_classification",
		})
	}

	// 8. Delivery mode check
	if strings.TrimSpace(string(asset.DeliveryMode)) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "delivery_mode",
			Message:  "Delivery mode is not specified",
			Severity: SeverityInfo,
			Code:     "missing_delivery_mode",
		})
	}

	// 9. Functional state resolution
	var resolvedState FunctionalState
	hasErrors := false
	for _, issue := range issues {
		if issue.Severity == SeverityError {
			hasErrors = true
			break
		}
	}

	if hasErrors {
		resolvedState = FunctionalStateDraft
	} else if asset.LifecycleState == LifecycleStateSuperseded || asset.Supersession.IsSuperseded {
		resolvedState = FunctionalStateSuperseded
	} else if asset.LifecycleState == LifecycleStateDeprecated {
		resolvedState = FunctionalStateDeprecated
	} else if len(issues) > 0 {
		// Contains warnings or info (e.g. missing risk classification or owner unknown)
		resolvedState = FunctionalStatePartial
	} else {
		resolvedState = FunctionalStateComplete
	}

	return MetadataValidationResult{
		Valid:                 !hasErrors,
		State:                 resolvedState,
		Issues:                issues,
		MissingRequiredFields: missingRequired,
		ValidatedAt:           now,
	}
}

// SuggestAssetTagsAndClassification generates deterministic tag and risk classification suggestions without calling AI (REQ-PROD-016).
func SuggestAssetTagsAndClassification(asset *Asset) TagSuggestionResult {
	if asset == nil {
		return TagSuggestionResult{
			SuggestedTags:         nil,
			SuggestedRisk:         RiskClassificationLow,
			SuggestedDeliveryMode: DeliveryModePlanned,
			Rationale:             []string{"Asset is nil; default low risk applied"},
			IsAIPowered:           false,
		}
	}

	tagSet := make(map[string]bool)
	var rationale []string

	// Tag by asset kind
	tagSet[string(asset.Type)] = true
	rationale = append(rationale, fmt.Sprintf("Added tag %q from asset kind", asset.Type))

	// Keyword-based deterministic tag discovery
	textCorpus := strings.ToLower(fmt.Sprintf("%s %s %s", asset.Name, asset.DisplayName, asset.Description))

	securityKeywords := map[string]string{
		"detect":      "detection",
		"alert":       "alerting",
		"threat":      "threat-intel",
		"incident":    "incident-response",
		"siem":        "siem",
		"edr":         "edr",
		"wazuh":       "wazuh-agent",
		"forensic":    "forensics",
		"report":      "reporting",
		"audit":       "compliance",
		"remediate":   "remediation",
		"containment": "containment",
		"isolate":     "containment",
		"secret":      "credentials",
	}

	for kw, tag := range securityKeywords {
		if strings.Contains(textCorpus, kw) {
			tagSet[tag] = true
			rationale = append(rationale, fmt.Sprintf("Added tag %q from keyword match %q", tag, kw))
		}
	}

	// Preserve existing tags
	for _, t := range asset.Tags {
		trimmed := strings.ToLower(strings.TrimSpace(t))
		if trimmed != "" {
			tagSet[trimmed] = true
		}
	}

	var suggestedTags []string
	for t := range tagSet {
		suggestedTags = append(suggestedTags, t)
	}

	// Risk classification rules (deterministic)
	suggestedRisk := RiskClassificationLow
	hasCriticalActions := strings.Contains(textCorpus, "terminate") ||
		strings.Contains(textCorpus, "isolate") ||
		strings.Contains(textCorpus, "delete") ||
		strings.Contains(textCorpus, "containment")

	hasSecretDeps := false
	for _, dep := range asset.Dependencies {
		if dep.Kind == "secret_reference" {
			hasSecretDeps = true
			break
		}
	}

	if hasCriticalActions {
		suggestedRisk = RiskClassificationCritical
		rationale = append(rationale, "Critical risk suggested due to containment, isolation, or destructive capabilities")
	} else if hasSecretDeps || asset.Type == AssetKindAutomationAgent || asset.Type == AssetKindAgentTeam {
		suggestedRisk = RiskClassificationHigh
		rationale = append(rationale, "High risk suggested due to credential dependencies or autonomous agent execution")
	} else if asset.Type == AssetKindWorkflow || len(asset.Parameters) > 3 {
		suggestedRisk = RiskClassificationMedium
		rationale = append(rationale, "Medium risk suggested due to complex multi-step workflow structure")
	} else {
		suggestedRisk = RiskClassificationLow
		rationale = append(rationale, "Low risk assigned for baseline read-only or inspection capability")
	}

	// Delivery mode suggestion
	suggestedDelivery := DeliveryModePlanned
	if asset.LifecycleState == LifecycleStatePublished {
		suggestedDelivery = DeliveryModeDelivered
	} else if asset.LifecycleState == LifecycleStateReview || asset.LifecycleState == LifecycleStateApproved {
		suggestedDelivery = DeliveryModePreview
	} else if asset.LifecycleState == LifecycleStateDeprecated || asset.LifecycleState == LifecycleStateSuperseded {
		suggestedDelivery = DeliveryModeDeprecated
	}

	return TagSuggestionResult{
		SuggestedTags:         suggestedTags,
		SuggestedRisk:         suggestedRisk,
		SuggestedDeliveryMode: suggestedDelivery,
		Rationale:             rationale,
		IsAIPowered:           false,
	}
}

// ValidateMetadata implements Class 0 metadata inspection validation.
func (s *DefaultCatalogService) ValidateMetadata(ctx context.Context, auth AuthContext, req ValidateMetadataRequest) (*ValidateMetadataResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Base library permission check
	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Asset ID check
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required", corrID, 400, false)
	}

	// 4. Retrieve asset from store
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve asset: %v", err), corrID, 500, false)
	}

	if asset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// 5. Evaluate deterministic validation
	valResult := ValidateAssetMetadata(asset)

	return &ValidateMetadataResponse{
		AssetID:    asset.ID,
		Validation: valResult,
	}, nil
}

// SuggestClassificationAndTags implements deterministic tag and classification suggestion without AI.
func (s *DefaultCatalogService) SuggestClassificationAndTags(ctx context.Context, auth AuthContext, req SuggestTagsRequest) (*SuggestTagsResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Base library permission check
	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Asset ID check
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required", corrID, 400, false)
	}

	// 4. Retrieve asset
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve asset: %v", err), corrID, 500, false)
	}

	if asset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// 5. Compute suggestions
	suggestions := SuggestAssetTagsAndClassification(asset)

	return &SuggestTagsResponse{
		AssetID:     asset.ID,
		Suggestions: suggestions,
	}, nil
}

// AssignOwner implements Class 2 assignment of owner reference with explicit unknown fallback (REQ-OBJ-009).
func (s *DefaultCatalogService) AssignOwner(ctx context.Context, auth AuthContext, req AssignOwnerRequest) (*AssignOwnerResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Base library permission check
	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Asset ID check
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required", corrID, 400, false)
	}

	// 4. Retrieve asset
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve asset: %v", err), corrID, 500, false)
	}

	if asset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// 5. Manage permission check for asset kind
	if !auth.CanManageAssetKind(asset.Type) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("manage permission for %s", asset.Type), corrID)
	}

	// 6. Normalize new owner per REQ-OBJ-009
	prevOwner := asset.Owner
	prevOwner.Normalize()

	newOwner := req.NewOwner
	newOwner.Normalize()

	// 7. Audit trail entry (CAP-STD-002 Section 20)
	actor := strings.TrimSpace(req.Actor)
	if actor == "" {
		actor = auth.PrincipalID
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "Asset owner reference assignment"
	}

	now := time.Now().UTC()
	auditEntry := AuditRecord{
		Timestamp:     now,
		Actor:         actor,
		Reason:        reason,
		OldMetadata:   fmt.Sprintf("owner_id=%s, owner_name=%s, is_unknown=%t", prevOwner.OwnerID, prevOwner.OwnerName, prevOwner.IsUnknown),
		NewMetadata:   fmt.Sprintf("owner_id=%s, owner_name=%s, is_unknown=%t", newOwner.OwnerID, newOwner.OwnerName, newOwner.IsUnknown),
		Version:       asset.Version.SemVer,
		TenantID:      asset.TenantID,
		EnvironmentID: asset.EnvironmentID,
		Action:        "assign_owner",
	}

	asset.Owner = newOwner
	asset.UpdatedAt = now
	asset.AuditTrail = append(asset.AuditTrail, auditEntry)

	// 8. Persist updated asset
	if err := s.store.Put(ctx, asset); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to persist asset owner: %v", err), corrID, 500, false)
	}

	return &AssignOwnerResponse{
		Asset:         asset.Clone(),
		PreviousOwner: prevOwner,
		AssignedAt:    now,
	}, nil
}

// UpdateDraftMetadata implements Class 2 editing of draft metadata.
func (s *DefaultCatalogService) UpdateDraftMetadata(ctx context.Context, auth AuthContext, req UpdateDraftMetadataRequest) (*UpdateDraftMetadataResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Base library permission check
	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Asset ID check
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required", corrID, 400, false)
	}

	// 4. Retrieve asset
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve asset: %v", err), corrID, 500, false)
	}

	if asset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// 5. Manage permission check
	if !auth.CanManageAssetKind(asset.Type) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("manage permission for %s", asset.Type), corrID)
	}

	// 6. Immutability check: only draft assets can have metadata edited
	if asset.LifecycleState != LifecycleStateDraft {
		return nil, NewImmutableAssetMutationError(asset.ID, string(asset.LifecycleState), corrID)
	}

	// 7. Apply updates
	if strings.TrimSpace(req.DisplayName) != "" {
		asset.DisplayName = req.DisplayName
	}
	if strings.TrimSpace(req.Description) != "" {
		asset.Description = req.Description
	}
	if req.Tags != nil {
		asset.Tags = req.Tags
	}
	if strings.TrimSpace(string(req.RiskClassification)) != "" {
		asset.RiskClassification = req.RiskClassification
	}
	if strings.TrimSpace(string(req.DeliveryMode)) != "" {
		asset.DeliveryMode = req.DeliveryMode
	}
	if strings.TrimSpace(req.ExecutionClass) != "" {
		asset.ExecutionClass = req.ExecutionClass
	}
	if req.IsAIPowered != nil {
		asset.IsAIPowered = *req.IsAIPowered
	}
	if req.HasNonAIFallback != nil {
		asset.HasNonAIFallback = *req.HasNonAIFallback
	}

	// 8. Re-evaluate validation
	valResult := ValidateAssetMetadata(asset)
	asset.FunctionalState = valResult.State

	// 9. Audit trail
	actor := strings.TrimSpace(req.Actor)
	if actor == "" {
		actor = auth.PrincipalID
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "Draft metadata update"
	}

	now := time.Now().UTC()
	auditEntry := AuditRecord{
		Timestamp:     now,
		Actor:         actor,
		Reason:        reason,
		OldMetadata:   "draft-prior",
		NewMetadata:   fmt.Sprintf("displayName=%s, risk=%s, delivery=%s", asset.DisplayName, asset.RiskClassification, asset.DeliveryMode),
		Version:       asset.Version.SemVer,
		TenantID:      asset.TenantID,
		EnvironmentID: asset.EnvironmentID,
		Action:        "update_draft",
	}
	asset.UpdatedAt = now
	asset.AuditTrail = append(asset.AuditTrail, auditEntry)

	// 10. Persist
	if err := s.store.Put(ctx, asset); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to persist updated asset: %v", err), corrID, 500, false)
	}

	return &UpdateDraftMetadataResponse{
		Asset:      asset.Clone(),
		Validation: valResult,
		UpdatedAt:  now,
	}, nil
}

// RecordSupersession implements Class 2 linking of prior version to newer version (CAP-STD-002 Section 11, 12, 16).
func (s *DefaultCatalogService) RecordSupersession(ctx context.Context, auth AuthContext, req RecordSupersessionRequest) (*RecordSupersessionResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Permission check: requires version manage or library manage
	if !auth.HasPermission(PermCmdrStudioVersionManage) && !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermCmdrStudioVersionManage, corrID)
	}

	// 3. ID checks
	olderID := strings.TrimSpace(req.OlderAssetID)
	newerID := strings.TrimSpace(req.NewerAssetID)
	if olderID == "" || newerID == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Both older_asset_id and newer_asset_id are required", corrID, 400, false)
	}
	if olderID == newerID {
		return nil, NewError(ErrorCodeInvalidMetadata, "An asset cannot supersede itself", corrID, 400, false)
	}

	// 4. Retrieve both assets
	olderAsset, err := s.store.Get(ctx, req.TenantRef, olderID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve older asset: %v", err), corrID, 500, false)
	}
	if olderAsset == nil {
		return nil, NewAssetNotFoundError(olderID, req.TenantRef, corrID)
	}

	newerAsset, err := s.store.Get(ctx, req.TenantRef, newerID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve newer asset: %v", err), corrID, 500, false)
	}
	if newerAsset == nil {
		return nil, NewAssetNotFoundError(newerID, req.TenantRef, corrID)
	}

	// Verify manage permission for both asset kinds
	if !auth.CanManageAssetKind(olderAsset.Type) && !auth.HasPermission(PermCmdrStudioVersionManage) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("manage permission for older asset kind %s", olderAsset.Type), corrID)
	}
	if !auth.CanManageAssetKind(newerAsset.Type) && !auth.HasPermission(PermCmdrStudioVersionManage) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("manage permission for newer asset kind %s", newerAsset.Type), corrID)
	}

	// 5. Check if older asset is already superseded by a different asset
	if olderAsset.Supersession.IsSuperseded && olderAsset.Supersession.SupersededBy != "" && olderAsset.Supersession.SupersededBy != newerID {
		return nil, NewError(
			ErrorCodeInvalidMetadata,
			fmt.Sprintf("Asset %q is already superseded by %q", olderID, olderAsset.Supersession.SupersededBy),
			corrID,
			400,
			false,
		)
	}

	// 6. Check for circular supersession
	if newerAsset.Supersession.SupersededBy == olderID || newerAsset.Supersession.Supersedes == newerID {
		return nil, NewCircularSupersessionError(olderID, corrID)
	}

	now := time.Now().UTC()
	actor := strings.TrimSpace(req.Actor)
	if actor == "" {
		actor = auth.PrincipalID
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = fmt.Sprintf("Asset %s superseded by %s", olderID, newerID)
	}

	// 7. Update older asset
	olderAsset.Supersession.IsSuperseded = true
	olderAsset.Supersession.SupersededBy = newerID
	olderAsset.LifecycleState = LifecycleStateSuperseded
	olderAsset.FunctionalState = FunctionalStateSuperseded
	olderAsset.UpdatedAt = now
	olderAsset.AuditTrail = append(olderAsset.AuditTrail, AuditRecord{
		Timestamp:       now,
		Actor:           actor,
		Reason:          reason,
		OldMetadata:     fmt.Sprintf("superseded_by=%s", olderAsset.Supersession.SupersededBy),
		NewMetadata:     fmt.Sprintf("superseded_by=%s", newerID),
		Version:         olderAsset.Version.SemVer,
		TenantID:        olderAsset.TenantID,
		EnvironmentID:   olderAsset.EnvironmentID,
		Action:          "superseded",
		SupersessionRef: newerID,
	})

	// 8. Update newer asset
	newerAsset.Supersession.Supersedes = olderID
	newerAsset.UpdatedAt = now
	newerAsset.AuditTrail = append(newerAsset.AuditTrail, AuditRecord{
		Timestamp:       now,
		Actor:           actor,
		Reason:          reason,
		OldMetadata:     fmt.Sprintf("supersedes=%s", newerAsset.Supersession.Supersedes),
		NewMetadata:     fmt.Sprintf("supersedes=%s", olderID),
		Version:         newerAsset.Version.SemVer,
		TenantID:        newerAsset.TenantID,
		EnvironmentID:   newerAsset.EnvironmentID,
		Action:          "supersedes",
		SupersessionRef: olderID,
	})

	// 9. Persist both
	if err := s.store.Put(ctx, olderAsset); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to update older asset: %v", err), corrID, 500, false)
	}
	if err := s.store.Put(ctx, newerAsset); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to update newer asset: %v", err), corrID, 500, false)
	}

	return &RecordSupersessionResponse{
		OlderAsset: olderAsset.ToSummary(),
		NewerAsset: newerAsset.ToSummary(),
		RecordedAt: now,
	}, nil
}

// GetSupersessionHistory implements historical lineage navigation (CAP-STD-002 Section 11, 16, 25).
func (s *DefaultCatalogService) GetSupersessionHistory(ctx context.Context, auth AuthContext, req SupersessionHistoryRequest) (*SupersessionHistoryResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	// 1. Tenant validation
	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Base library permission check
	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Asset ID check
	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required", corrID, 400, false)
	}

	// 4. Retrieve target asset
	targetAsset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve asset: %v", err), corrID, 500, false)
	}
	if targetAsset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// 5. Read permission check for asset kind
	if !auth.CanReadAssetKind(targetAsset.Type) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("read permission for %s", targetAsset.Type), corrID)
	}

	// 6. Trace backward (ancestors via Supersedes)
	var ancestors []*Asset
	visitedBack := make(map[string]bool)
	visitedBack[targetAsset.ID] = true

	currentOlderID := targetAsset.Supersession.Supersedes
	for currentOlderID != "" {
		if visitedBack[currentOlderID] {
			return nil, NewCircularSupersessionError(currentOlderID, corrID)
		}
		visitedBack[currentOlderID] = true

		ancestor, err := s.store.Get(ctx, req.TenantRef, currentOlderID)
		if err != nil || ancestor == nil {
			// Ancestor not in current store; stop tracing back
			break
		}
		ancestors = append([]*Asset{ancestor}, ancestors...) // Prepend so oldest is first
		currentOlderID = ancestor.Supersession.Supersedes
	}

	// 7. Trace forward (descendants via SupersededBy)
	var descendants []*Asset
	visitedFwd := make(map[string]bool)
	visitedFwd[targetAsset.ID] = true

	currentNewerID := targetAsset.Supersession.SupersededBy
	for currentNewerID != "" {
		if visitedFwd[currentNewerID] || visitedBack[currentNewerID] {
			return nil, NewCircularSupersessionError(currentNewerID, corrID)
		}
		visitedFwd[currentNewerID] = true

		descendant, err := s.store.Get(ctx, req.TenantRef, currentNewerID)
		if err != nil || descendant == nil {
			// Descendant not in current store; stop tracing forward
			break
		}
		descendants = append(descendants, descendant)
		currentNewerID = descendant.Supersession.SupersededBy
	}

	// 8. Assemble full ordered chain
	var fullChain []*Asset
	fullChain = append(fullChain, ancestors...)
	fullChain = append(fullChain, targetAsset)
	fullChain = append(fullChain, descendants...)

	chainNodes := make([]SupersessionHistoryNode, len(fullChain))
	for i, a := range fullChain {
		isActive := !a.Supersession.IsSuperseded && a.LifecycleState != LifecycleStateRetired && a.LifecycleState != LifecycleStateArchived
		chainNodes[i] = SupersessionHistoryNode{
			AssetID:        a.ID,
			Version:        a.Version.SemVer,
			LifecycleState: a.LifecycleState,
			IsCurrent:      a.ID == targetAsset.ID,
			IsActive:       isActive,
			Supersedes:     a.Supersession.Supersedes,
			SupersededBy:   a.Supersession.SupersededBy,
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
			Owner:          a.Owner,
			Summary:        a.ToSummary(),
		}
	}

	rootID := fullChain[0].ID
	latestID := fullChain[len(fullChain)-1].ID

	return &SupersessionHistoryResponse{
		RootAssetID:    rootID,
		LatestAssetID:  latestID,
		CurrentAssetID: targetAsset.ID,
		Chain:          chainNodes,
	}, nil
}
