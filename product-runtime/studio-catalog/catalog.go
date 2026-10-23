package studiocatalog

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// CatalogService defines the 8 Class 0 operations and permanently blocked execution verbs.
type CatalogService interface {
	// The 8 Class 0 Operations (REQ-PROD-006, CAP-STD-001, CAP-STD-002)
	Browse(ctx context.Context, auth AuthContext, req BrowseRequest) (*BrowseResponse, error)
	Search(ctx context.Context, auth AuthContext, req SearchRequest) (*SearchResponse, error)
	Filter(ctx context.Context, auth AuthContext, req FilterRequest) (*FilterResponse, error)
	Group(ctx context.Context, auth AuthContext, req GroupRequest) (*GroupResponse, error)
	Open(ctx context.Context, auth AuthContext, req OpenRequest) (*OpenResponse, error)
	Inspect(ctx context.Context, auth AuthContext, req InspectRequest) (*InspectResponse, error)
	Compare(ctx context.Context, auth AuthContext, req CompareRequest) (*CompareResponse, error)
	ClonePreparation(ctx context.Context, auth AuthContext, req ClonePreparationRequest) (*ClonePreparationResponse, error)

	// Permanently blocked execution triggers under REQ-PROD-006
	Invoke(ctx context.Context, auth AuthContext, req any) error
	Execute(ctx context.Context, auth AuthContext, req any) error
	Run(ctx context.Context, auth AuthContext, req any) error
	Call(ctx context.Context, auth AuthContext, req any) error
	Trigger(ctx context.Context, auth AuthContext, req any) error
}

// Store defines persistence methods for Studio catalog assets.
type Store interface {
	Get(ctx context.Context, tenantID, assetID string) (*Asset, error)
	List(ctx context.Context, tenantID string) ([]*Asset, error)
	Put(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, tenantID, assetID string) error
	GetAllAssetsAcrossTenants(ctx context.Context) ([]*Asset, error) // For internal oracle checking
}

// InMemoryStore provides a thread-safe in-memory store for assets.
type InMemoryStore struct {
	mu     sync.RWMutex
	assets map[string]*Asset // keyed by tenantID + "/" + assetID
}

// NewInMemoryStore initializes an empty in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		assets: make(map[string]*Asset),
	}
}

// Get retrieves an asset by tenant and ID.
func (s *InMemoryStore) Get(ctx context.Context, tenantID, assetID string) (*Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := tenantID + "/" + assetID
	a, ok := s.assets[key]
	if !ok {
		return nil, nil
	}
	return a.Clone(), nil
}

// List returns all assets for a given tenant.
func (s *InMemoryStore) List(ctx context.Context, tenantID string) ([]*Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prefix := tenantID + "/"
	var result []*Asset
	for k, a := range s.assets {
		if strings.HasPrefix(k, prefix) {
			result = append(result, a.Clone())
		}
	}
	return result, nil
}

// Put stores or updates an asset.
func (s *InMemoryStore) Put(ctx context.Context, asset *Asset) error {
	if asset == nil {
		return fmt.Errorf("cannot put nil asset")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := asset.Clone()
	key := cp.TenantID + "/" + cp.ID
	s.assets[key] = cp
	return nil
}

// Delete removes an asset by tenant and ID.
func (s *InMemoryStore) Delete(ctx context.Context, tenantID, assetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := tenantID + "/" + assetID
	delete(s.assets, key)
	return nil
}

// GetAllAssetsAcrossTenants returns all assets across all tenants.
func (s *InMemoryStore) GetAllAssetsAcrossTenants(ctx context.Context) ([]*Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Asset
	for _, a := range s.assets {
		result = append(result, a.Clone())
	}
	return result, nil
}

// DefaultCatalogService implements the CatalogService interface.
type DefaultCatalogService struct {
	store Store
}

// NewCatalogService creates a new CatalogService instance backed by the given store.
func NewCatalogService(store Store) *DefaultCatalogService {
	return &DefaultCatalogService{
		store: store,
	}
}

// clampLimit validates and bounds page limit between 1 and 100 (default 25).
func clampLimit(limit int) int {
	if limit <= 0 {
		return 25
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// Browse implements Class 0 Operation 1: Browse.
func (s *DefaultCatalogService) Browse(ctx context.Context, auth AuthContext, req BrowseRequest) (*BrowseResponse, error) {
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

	// 3. Retrieve assets for tenant
	assets, err := s.store.List(ctx, req.TenantRef)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve catalog assets", corrID, 500, true)
	}

	// 4. Filter visible assets
	var visible []*Asset
	for _, a := range assets {
		// Silent omission if caller cannot read asset kind
		if !auth.CanReadAssetKind(a.Type) {
			continue
		}
		if req.Kind != "" && a.Type != req.Kind {
			continue
		}
		visible = append(visible, a)
	}

	// 5. Compound Sort: (sort_field, id ASC)
	sortField := req.SortField
	if sortField == "" {
		sortField = "updated_at"
	}
	sortDir := strings.ToLower(req.SortDirection)
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	sortAssets(visible, sortField, sortDir)

	// 6. Paginate using opaque cursor
	limit := clampLimit(req.Limit)
	startIndex := 0
	if req.Cursor != "" {
		payload, err := DecodeCursor(req.Cursor)
		if err != nil || payload.TenantID != req.TenantRef || payload.SortField != sortField || payload.SortDirection != sortDir {
			return nil, NewError(ErrorCodeInvalidCursor, "Invalid or expired cursor token", corrID, 400, false)
		}
		// Find index of tiebreaker
		found := false
		for i, a := range visible {
			if a.ID == payload.TiebreakerID {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			// If cursor item disappeared, cursor is invalid
			return nil, NewError(ErrorCodeInvalidCursor, "Cursor anchor point not found", corrID, 400, false)
		}
	}

	endIndex := startIndex + limit
	hasMore := false
	if endIndex < len(visible) {
		hasMore = true
	} else {
		endIndex = len(visible)
	}

	var pageItems []AssetSummary
	if startIndex < len(visible) {
		for _, a := range visible[startIndex:endIndex] {
			pageItems = append(pageItems, a.ToSummary())
		}
	}

	var cursorNext string
	if hasMore && len(pageItems) > 0 {
		lastItem := pageItems[len(pageItems)-1]
		cursorNext = EncodeCursor(CursorPayload{
			Version:       1,
			TenantID:      req.TenantRef,
			SortField:     sortField,
			SortDirection: sortDir,
			LastValue:     getSortFieldValue(lastItem, sortField),
			TiebreakerID:  lastItem.ID,
			Fingerprint:   fmt.Sprintf("fp-%d", len(visible)),
		})
	}

	return &BrowseResponse{
		Items: pageItems,
		Pagination: PaginationInfo{
			Limit:        limit,
			HasMore:      hasMore,
			CursorNext:   cursorNext,
			TotalCount:   len(visible),
			IsTotalExact: true,
		},
	}, nil
}

// Search implements Class 0 Operation 2: Search.
func (s *DefaultCatalogService) Search(ctx context.Context, auth AuthContext, req SearchRequest) (*SearchResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	query := strings.TrimSpace(req.QueryText)
	if query == "" {
		return nil, NewError(ErrorCodeInvalidFilterParameter, "Search query text cannot be empty", corrID, 400, false)
	}

	assets, err := s.store.List(ctx, req.TenantRef)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve search assets", corrID, 500, true)
	}

	queryLower := strings.ToLower(query)
	var matched []*Asset
	for _, a := range assets {
		if !auth.CanReadAssetKind(a.Type) {
			continue
		}
		if req.Kind != "" && a.Type != req.Kind {
			continue
		}

		// Deterministic token & substring matching
		nameMatch := strings.Contains(strings.ToLower(a.Name), queryLower)
		displayMatch := strings.Contains(strings.ToLower(a.DisplayName), queryLower)
		descMatch := strings.Contains(strings.ToLower(a.Description), queryLower)
		tagMatch := false
		for _, tag := range a.Tags {
			if strings.Contains(strings.ToLower(tag), queryLower) {
				tagMatch = true
				break
			}
		}

		if nameMatch || displayMatch || descMatch || tagMatch {
			matched = append(matched, a)
		}
	}

	// Sort matching items
	sortAssets(matched, "updated_at", "desc")

	limit := clampLimit(req.Limit)
	startIndex := 0
	if req.Cursor != "" {
		payload, err := DecodeCursor(req.Cursor)
		if err != nil || payload.TenantID != req.TenantRef {
			return nil, NewError(ErrorCodeInvalidCursor, "Invalid search cursor", corrID, 400, false)
		}
		found := false
		for i, a := range matched {
			if a.ID == payload.TiebreakerID {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, NewError(ErrorCodeInvalidCursor, "Search cursor item not found", corrID, 400, false)
		}
	}

	endIndex := startIndex + limit
	hasMore := false
	if endIndex < len(matched) {
		hasMore = true
	} else {
		endIndex = len(matched)
	}

	var pageItems []AssetSummary
	if startIndex < len(matched) {
		for _, a := range matched[startIndex:endIndex] {
			summary := a.ToSummary()
			// If Shared Search is unavailable, mark functional state as partial degradation
			if !req.SharedSearchAvailable {
				summary.FunctionalState = FunctionalStatePartial
			}
			pageItems = append(pageItems, summary)
		}
	}

	var cursorNext string
	if hasMore && len(pageItems) > 0 {
		lastItem := pageItems[len(pageItems)-1]
		cursorNext = EncodeCursor(CursorPayload{
			Version:       1,
			TenantID:      req.TenantRef,
			SortField:     "updated_at",
			SortDirection: "desc",
			LastValue:     lastItem.UpdatedAt.Format(time.RFC3339Nano),
			TiebreakerID:  lastItem.ID,
			Fingerprint:   fmt.Sprintf("search-%d", len(matched)),
		})
	}

	return &SearchResponse{
		Items: pageItems,
		Pagination: PaginationInfo{
			Limit:        limit,
			HasMore:      hasMore,
			CursorNext:   cursorNext,
			TotalCount:   len(matched),
			IsTotalExact: true,
		},
	}, nil
}

// Filter implements Class 0 Operation 3: Filter.
func (s *DefaultCatalogService) Filter(ctx context.Context, auth AuthContext, req FilterRequest) (*FilterResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	assets, err := s.store.List(ctx, req.TenantRef)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve filter assets", corrID, 500, true)
	}

	var filtered []*Asset
	for _, a := range assets {
		if !auth.CanReadAssetKind(a.Type) {
			continue
		}

		// Kinds filter
		if len(req.Criteria.Kinds) > 0 {
			match := false
			for _, k := range req.Criteria.Kinds {
				if a.Type == k {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// Lifecycle states filter
		if len(req.Criteria.LifecycleStates) > 0 {
			match := false
			for _, ls := range req.Criteria.LifecycleStates {
				if a.LifecycleState == ls {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// Functional states filter
		if len(req.Criteria.FunctionalStates) > 0 {
			match := false
			for _, fs := range req.Criteria.FunctionalStates {
				if a.FunctionalState == fs {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// Owner filter (including "unknown")
		if req.Criteria.OwnerID != "" {
			if req.Criteria.OwnerID == "unknown" {
				if !a.Owner.IsUnknown {
					continue
				}
			} else {
				if a.Owner.OwnerID != req.Criteria.OwnerID {
					continue
				}
			}
		}

		// Risk levels filter
		if len(req.Criteria.RiskLevels) > 0 {
			match := false
			for _, r := range req.Criteria.RiskLevels {
				if a.RiskClassification == r {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		// AI powered flag filter
		if req.Criteria.IsAIPowered != nil {
			if a.IsAIPowered != *req.Criteria.IsAIPowered {
				continue
			}
		}

		// Tags filter (AND match: all requested tags must be present)
		if len(req.Criteria.Tags) > 0 {
			allTagsFound := true
			for _, reqTag := range req.Criteria.Tags {
				found := false
				for _, aTag := range a.Tags {
					if strings.EqualFold(aTag, reqTag) {
						found = true
						break
					}
				}
				if !found {
					allTagsFound = false
					break
				}
			}
			if !allTagsFound {
				continue
			}
		}

		filtered = append(filtered, a)
	}

	sortField := req.SortField
	if sortField == "" {
		sortField = "updated_at"
	}
	sortDir := strings.ToLower(req.SortDirection)
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	sortAssets(filtered, sortField, sortDir)

	limit := clampLimit(req.Limit)
	startIndex := 0
	if req.Cursor != "" {
		payload, err := DecodeCursor(req.Cursor)
		if err != nil || payload.TenantID != req.TenantRef {
			return nil, NewError(ErrorCodeInvalidCursor, "Invalid filter cursor", corrID, 400, false)
		}
		found := false
		for i, a := range filtered {
			if a.ID == payload.TiebreakerID {
				startIndex = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, NewError(ErrorCodeInvalidCursor, "Filter cursor item not found", corrID, 400, false)
		}
	}

	endIndex := startIndex + limit
	hasMore := false
	if endIndex < len(filtered) {
		hasMore = true
	} else {
		endIndex = len(filtered)
	}

	var pageItems []AssetSummary
	if startIndex < len(filtered) {
		for _, a := range filtered[startIndex:endIndex] {
			pageItems = append(pageItems, a.ToSummary())
		}
	}

	var cursorNext string
	if hasMore && len(pageItems) > 0 {
		lastItem := pageItems[len(pageItems)-1]
		cursorNext = EncodeCursor(CursorPayload{
			Version:       1,
			TenantID:      req.TenantRef,
			SortField:     sortField,
			SortDirection: sortDir,
			LastValue:     getSortFieldValue(lastItem, sortField),
			TiebreakerID:  lastItem.ID,
			Fingerprint:   fmt.Sprintf("filter-%d", len(filtered)),
		})
	}

	return &FilterResponse{
		Items: pageItems,
		Pagination: PaginationInfo{
			Limit:        limit,
			HasMore:      hasMore,
			CursorNext:   cursorNext,
			TotalCount:   len(filtered),
			IsTotalExact: true,
		},
	}, nil
}

// Group implements Class 0 Operation 4: Group.
func (s *DefaultCatalogService) Group(ctx context.Context, auth AuthContext, req GroupRequest) (*GroupResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// Validate grouping facet
	switch req.GroupBy {
	case FacetKind, FacetLifecycleState, FacetOwner, FacetRiskClassification:
	default:
		return nil, NewError(
			ErrorCodeUnsupportedGroupingFacet,
			fmt.Sprintf("Facet %q is not supported. Must be by_kind, by_lifecycle_state, by_owner, or by_risk_classification", req.GroupBy),
			corrID,
			400,
			false,
		)
	}

	assets, err := s.store.List(ctx, req.TenantRef)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve grouping assets", corrID, 500, true)
	}

	counts := make(map[string]int)
	total := 0

	for _, a := range assets {
		if !auth.CanReadAssetKind(a.Type) {
			continue
		}

		// Apply optional pre-filter if provided
		if req.Criteria != nil {
			if len(req.Criteria.Kinds) > 0 {
				match := false
				for _, k := range req.Criteria.Kinds {
					if a.Type == k {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			if len(req.Criteria.LifecycleStates) > 0 {
				match := false
				for _, ls := range req.Criteria.LifecycleStates {
					if a.LifecycleState == ls {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
		}

		total++
		var key string
		switch req.GroupBy {
		case FacetKind:
			key = string(a.Type)
		case FacetLifecycleState:
			key = string(a.LifecycleState)
		case FacetOwner:
			a.Owner.Normalize()
			key = a.Owner.DisplayName
		case FacetRiskClassification:
			key = string(a.RiskClassification)
		}

		counts[key]++
	}

	// Sort buckets deterministically by key ASC
	var keys []string
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buckets []GroupBucket
	for _, k := range keys {
		buckets = append(buckets, GroupBucket{Key: k, Count: counts[k]})
	}

	return &GroupResponse{
		Facet:      req.GroupBy,
		Groups:     buckets,
		TotalItems: total,
	}, nil
}

// Open implements Class 0 Operation 5: Open.
func (s *DefaultCatalogService) Open(ctx context.Context, auth AuthContext, req OpenRequest) (*OpenResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required to open an asset", corrID, 400, false)
	}

	if strings.TrimSpace(req.ReturnOrigin) == "" {
		return nil, NewError(ErrorCodeMissingReturnOrigin, "Return origin is required to open an asset", corrID, 400, false)
	}

	// Existence Oracle Elimination: Look up in active tenant
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to inspect asset", corrID, 500, true)
	}

	if asset == nil {
		// Uniform 404 whether non-existent or foreign tenant
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// Caller must have read permission for this kind
	if !auth.CanReadAssetKind(asset.Type) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("perm.studio.%s.read", asset.Type), corrID)
	}

	targetSurface := "builder"
	surfaceRoute := fmt.Sprintf("/studio/builder/%s/%s", asset.Type, asset.ID)

	return &OpenResponse{
		TargetSurface: targetSurface,
		SurfaceRoute:  surfaceRoute,
		AssetID:       asset.ID,
		VersionID:     asset.Version.VersionID,
		ReturnOrigin:  req.ReturnOrigin,
		ContextTokens: map[string]string{
			"tenant_id": req.TenantRef,
			"type":      string(asset.Type),
		},
	}, nil
}

// Inspect implements Class 0 Operation 6: Inspect.
func (s *DefaultCatalogService) Inspect(ctx context.Context, auth AuthContext, req InspectRequest) (*InspectResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	if strings.TrimSpace(req.AssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Asset ID is required to inspect an asset", corrID, 400, false)
	}

	// Retrieve asset in tenant
	asset, err := s.store.Get(ctx, req.TenantRef, req.AssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve asset", corrID, 500, true)
	}

	// Existence Oracle Elimination:
	// If asset not found in active tenant scope, check across all tenants to ensure uniform 404!
	if asset == nil {
		return nil, NewAssetNotFoundError(req.AssetID, req.TenantRef, corrID)
	}

	// Pre-authorization check: If caller lacks asset-specific read permission,
	// return HTTP 200 sanitized stub with functional_state: "permission-filtered" (Existence Oracle Elimination)
	if !auth.CanReadAssetKind(asset.Type) {
		stub := SanitizeAssetForInspect(asset)
		return &InspectResponse{Asset: stub}, nil
	}

	return &InspectResponse{Asset: asset.Clone()}, nil
}

// Compare implements Class 0 Operation 7: Compare.
func (s *DefaultCatalogService) Compare(ctx context.Context, auth AuthContext, req CompareRequest) (*CompareResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	if strings.TrimSpace(req.BaseAssetID) == "" || strings.TrimSpace(req.TargetAssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Both base_asset_id and target_asset_id are required for comparison", corrID, 400, false)
	}

	baseAsset, err := s.store.Get(ctx, req.TenantRef, req.BaseAssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve base asset", corrID, 500, true)
	}
	if baseAsset == nil {
		return nil, NewAssetNotFoundError(req.BaseAssetID, req.TenantRef, corrID)
	}

	targetAsset, err := s.store.Get(ctx, req.TenantRef, req.TargetAssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve target asset", corrID, 500, true)
	}
	if targetAsset == nil {
		return nil, NewAssetNotFoundError(req.TargetAssetID, req.TenantRef, corrID)
	}

	// Permissions check on both kinds
	if !auth.CanReadAssetKind(baseAsset.Type) || !auth.CanReadAssetKind(targetAsset.Type) {
		return nil, NewPermissionDeniedError("insufficient permissions to compare assets", corrID)
	}

	diff := computeStructuralDiff(baseAsset, targetAsset)

	return &CompareResponse{
		Base:   baseAsset.ToSummary(),
		Target: targetAsset.ToSummary(),
		Diff:   diff,
	}, nil
}

// computeStructuralDiff calculates parameter and dependency differences between two assets.
func computeStructuralDiff(base, target *Asset) StructuralDiff {
	diff := StructuralDiff{
		OldRisk:                   base.RiskClassification,
		NewRisk:                   target.RiskClassification,
		RiskClassificationChanged: base.RiskClassification != target.RiskClassification,
		AIFlagChanged:             base.IsAIPowered != target.IsAIPowered,
		OldLifecycle:              base.LifecycleState,
		NewLifecycle:              target.LifecycleState,
		LifecycleChanged:          base.LifecycleState != target.LifecycleState,
	}

	baseParams := make(map[string]Parameter)
	for _, p := range base.Parameters {
		baseParams[p.Name] = p
	}

	targetParams := make(map[string]Parameter)
	for _, p := range target.Parameters {
		targetParams[p.Name] = p
		bp, exists := baseParams[p.Name]
		if !exists {
			diff.ParametersAdded = append(diff.ParametersAdded, p)
		} else if bp.Type != p.Type || bp.Required != p.Required {
			diff.ParametersModified = append(diff.ParametersModified, p)
		}
	}

	for _, bp := range base.Parameters {
		if _, exists := targetParams[bp.Name]; !exists {
			diff.ParametersRemoved = append(diff.ParametersRemoved, bp)
		}
	}

	// Dependencies comparison
	baseDeps := make(map[string]Dependency)
	for _, d := range base.Dependencies {
		baseDeps[d.RefID] = d
	}

	for _, td := range target.Dependencies {
		bd, exists := baseDeps[td.RefID]
		if exists && bd.VersionConstraint != td.VersionConstraint {
			diff.DependenciesChanged = append(diff.DependenciesChanged, DependencyChange{
				RefID:         td.RefID,
				OldConstraint: bd.VersionConstraint,
				NewConstraint: td.VersionConstraint,
			})
		}
	}

	return diff
}

// ClonePreparation implements Class 0 Operation 8: Clone Preparation.
func (s *DefaultCatalogService) ClonePreparation(ctx context.Context, auth AuthContext, req ClonePreparationRequest) (*ClonePreparationResponse, error) {
	corrID := req.CorrelationID
	if corrID == "" {
		corrID = GenerateCorrelationID()
	}

	if errEnv := ValidateTenantScope(req.TenantRef, auth.AuthorizedTenantSet, corrID); errEnv != nil {
		return nil, errEnv
	}

	if !auth.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	if strings.TrimSpace(req.SourceAssetID) == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Source asset ID is required for clone preparation", corrID, 400, false)
	}

	// Strictly prohibit cross-tenant cloning under ADR-0008
	if req.TargetTenantRef != "" && req.TargetTenantRef != req.TenantRef {
		return nil, NewError(
			ErrorCodeCrossTenantMutationForbidden,
			fmt.Sprintf("Cross-tenant cloning from %q to %q is strictly forbidden", req.TenantRef, req.TargetTenantRef),
			corrID,
			403,
			false,
		)
	}

	source, err := s.store.Get(ctx, req.TenantRef, req.SourceAssetID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, "Failed to retrieve source asset", corrID, 500, true)
	}
	if source == nil {
		return nil, NewAssetNotFoundError(req.SourceAssetID, req.TenantRef, corrID)
	}

	if !auth.CanReadAssetKind(source.Type) {
		return nil, NewPermissionDeniedError(fmt.Sprintf("perm.studio.%s.read", source.Type), corrID)
	}

	now := time.Now().UTC()
	draftID := fmt.Sprintf("%s-clone-%s", source.Type, GenerateRandomSuffix())

	// Build deep copied clone
	clone := source.Clone()
	clone.ID = draftID
	clone.Name = source.Name + "-copy"
	clone.DisplayName = "Copy of " + source.DisplayName
	clone.LifecycleState = LifecycleStateDraft
	clone.FunctionalState = FunctionalStateReady
	clone.Version = VersionInfo{
		VersionID:   fmt.Sprintf("ver-%s-draft", draftID),
		SemVer:      "0.1.0-draft",
		CreatedAt:   now,
		IsImmutable: false,
	}

	// Strip credentials & secret values from dependencies
	var sanitizedDeps []Dependency
	for _, d := range clone.Dependencies {
		if d.Kind == "secret_reference" {
			// Mask or omit raw secret reference
			d.IsMasked = true
			d.Purpose = "Opaque secret reference preserved (unbound)"
		}
		sanitizedDeps = append(sanitizedDeps, d)
	}
	clone.Dependencies = sanitizedDeps

	// Record provenance
	chain := make([]string, len(source.Provenance.ProvenanceChain))
	copy(chain, source.Provenance.ProvenanceChain)
	sourceRef := fmt.Sprintf("%s@%s", source.ID, source.Version.SemVer)
	chain = append(chain, sourceRef)

	clonedBy := req.ClonedBy
	if clonedBy == "" {
		clonedBy = auth.PrincipalID
	}
	if clonedBy == "" {
		clonedBy = "unknown"
	}

	clone.Provenance = Provenance{
		CreatedBy:       clonedBy,
		CreatedAt:       now,
		UpdatedBy:       clonedBy,
		UpdatedAt:       now,
		ParentID:        source.ID,
		ClonedFrom:      source.ID,
		ProvenanceChain: chain,
	}
	clone.CreatedAt = now
	clone.UpdatedAt = now

	return &ClonePreparationResponse{
		DraftClone: DraftCloneDescriptor{
			DraftAssetID:      draftID,
			ClonedFromID:      source.ID,
			SourceVersion:     source.Version.SemVer,
			SourceTenantID:    req.TenantRef,
			TargetTenantID:    req.TenantRef,
			ClonedBy:          clonedBy,
			ClonedAt:          now,
			ProvenanceChain:   chain,
			PrefilledMetadata: clone,
		},
	}, nil
}

// Execution verbs permanently blocked under REQ-PROD-006:

func (s *DefaultCatalogService) Invoke(ctx context.Context, auth AuthContext, req any) error {
	return NewExecutionBlockedError("invoke", GenerateCorrelationID())
}

func (s *DefaultCatalogService) Execute(ctx context.Context, auth AuthContext, req any) error {
	return NewExecutionBlockedError("execute", GenerateCorrelationID())
}

func (s *DefaultCatalogService) Run(ctx context.Context, auth AuthContext, req any) error {
	return NewExecutionBlockedError("run", GenerateCorrelationID())
}

func (s *DefaultCatalogService) Call(ctx context.Context, auth AuthContext, req any) error {
	return NewExecutionBlockedError("call", GenerateCorrelationID())
}

func (s *DefaultCatalogService) Trigger(ctx context.Context, auth AuthContext, req any) error {
	return NewExecutionBlockedError("trigger", GenerateCorrelationID())
}

// Helper sorting functions:

func sortAssets(assets []*Asset, sortField, sortDir string) {
	sort.Slice(assets, func(i, j int) bool {
		valI := getAssetSortValue(assets[i], sortField)
		valJ := getAssetSortValue(assets[j], sortField)

		if valI != valJ {
			if sortDir == "asc" {
				return valI < valJ
			}
			return valI > valJ
		}
		// Deterministic tiebreaker: id ASC
		return assets[i].ID < assets[j].ID
	})
}

func getAssetSortValue(a *Asset, field string) string {
	switch strings.ToLower(field) {
	case "name":
		return strings.ToLower(a.Name)
	case "id":
		return a.ID
	case "created_at":
		return a.CreatedAt.Format(time.RFC3339Nano)
	case "updated_at":
		return a.UpdatedAt.Format(time.RFC3339Nano)
	default:
		return a.UpdatedAt.Format(time.RFC3339Nano)
	}
}

func getSortFieldValue(s AssetSummary, field string) string {
	switch strings.ToLower(field) {
	case "name":
		return strings.ToLower(s.Name)
	case "id":
		return s.ID
	case "updated_at":
		return s.UpdatedAt.Format(time.RFC3339Nano)
	default:
		return s.UpdatedAt.Format(time.RFC3339Nano)
	}
}

// GenerateRandomSuffix returns a short hex string.
func GenerateRandomSuffix() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "abcd"
	}
	return hex.EncodeToString(b)
}
