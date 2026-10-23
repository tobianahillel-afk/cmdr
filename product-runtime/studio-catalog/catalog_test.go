package studiocatalog

import (
	"context"
	"strings"
	"testing"
	"time"
)

// helper to create a test asset
func makeTestAsset(tenantID, id string, kind AssetKind, name string, lifecycle LifecycleState, risk RiskClassification) *Asset {
	now := time.Now().UTC()
	return &Asset{
		ID:             id,
		TenantID:       tenantID,
		EnvironmentID:  "env-prod",
		Type:           kind,
		Name:           name,
		DisplayName:    "Display " + name,
		Description:    "Description for " + name,
		Tags:           []string{"security", string(kind)},
		Owner: Owner{
			OwnerID:     "usr-lead",
			OwnerName:   "Lead Analyst",
			OwnerType:   "user",
			IsUnknown:   false,
			DisplayName: "Lead Analyst",
		},
		Contributors: []string{"usr-analyst-1"},
		Version: VersionInfo{
			VersionID:   "ver-1.0.0",
			SemVer:      "1.0.0",
			CreatedAt:   now,
			IsImmutable: true,
		},
		LifecycleState:     lifecycle,
		FunctionalState:    FunctionalStateReady,
		ExecutionClass:     "class_0_inspect",
		RiskClassification: risk,
		IsAIPowered:        true,
		HasNonAIFallback:   true,
		Parameters: []Parameter{
			{Name: "ip_address", Type: "string", Required: true, Description: "Target IP"},
			{Name: "timeout", Type: "int", Required: false, Description: "Timeout seconds"},
		},
		Dependencies: []Dependency{
			{
				Kind:              "tool",
				RefID:             "tool-whois",
				VersionConstraint: "^1.0.0",
				IsOptional:        false,
				HealthStatus:      "available",
			},
			{
				Kind:       "secret_reference",
				RefID:      "sec-api-key",
				Purpose:    "Third-party API key",
				IsOptional: true,
				IsMasked:   true,
			},
		},
		Consumers: []Consumer{
			{ConsumerType: "workflow", ConsumerRef: "wfl-incident-response"},
		},
		Supersession: Supersession{
			IsSuperseded: false,
		},
		Provenance: Provenance{
			CreatedBy:       "usr-lead",
			CreatedAt:       now,
			UpdatedBy:       "usr-lead",
			UpdatedAt:       now,
			ProvenanceChain: []string{id + "@1.0.0"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// setupTestStore seeds assets across two tenants (tenant-alpha and tenant-bravo)
func setupTestStore(t *testing.T) (*InMemoryStore, *DefaultCatalogService) {
	store := NewInMemoryStore()
	ctx := context.Background()

	// Tenant Alpha assets
	assetsAlpha := []*Asset{
		makeTestAsset("tenant-alpha", "skl-01", AssetKindSkill, "triage-alert", LifecycleStatePublished, RiskClassificationLow),
		makeTestAsset("tenant-alpha", "skl-02", AssetKindSkill, "enrich-ioc", LifecycleStatePublished, RiskClassificationMedium),
		makeTestAsset("tenant-alpha", "wfl-01", AssetKindWorkflow, "incident-flow", LifecycleStatePublished, RiskClassificationHigh),
		makeTestAsset("tenant-alpha", "wfl-02", AssetKindWorkflow, "phishing-investigation", LifecycleStateDraft, RiskClassificationMedium),
		makeTestAsset("tenant-alpha", "agn-01", AssetKindAutomationAgent, "sentinel-agent", LifecycleStatePublished, RiskClassificationLow),
		makeTestAsset("tenant-alpha", "agn-02", AssetKindAutomationAgent, "containment-agent", LifecycleStateReview, RiskClassificationCritical),
		makeTestAsset("tenant-alpha", "team-01", AssetKindAgentTeam, "soc-tier1-team", LifecycleStatePublished, RiskClassificationLow),
		makeTestAsset("tenant-alpha", "tool-01", AssetKindTool, "whois-lookup", LifecycleStatePublished, RiskClassificationLow),
		makeTestAsset("tenant-alpha", "tool-02", AssetKindTool, "ip-reputation", LifecycleStatePublished, RiskClassificationLow),
	}

	// Make one asset with missing owner to test REQ-OBJ-009
	unknownOwnerAsset := makeTestAsset("tenant-alpha", "skl-no-owner", AssetKindSkill, "legacy-orphan-skill", LifecycleStateDeprecated, RiskClassificationLow)
	unknownOwnerAsset.Owner = Owner{} // Blank owner
	assetsAlpha = append(assetsAlpha, unknownOwnerAsset)

	for _, a := range assetsAlpha {
		if err := store.Put(ctx, a); err != nil {
			t.Fatalf("failed to seed asset: %v", err)
		}
	}

	// Tenant Bravo assets (Foreign tenant)
	assetsBravo := []*Asset{
		makeTestAsset("tenant-bravo", "skl-bravo-secret", AssetKindSkill, "bravo-confidential", LifecycleStatePublished, RiskClassificationHigh),
		makeTestAsset("tenant-bravo", "wfl-bravo-01", AssetKindWorkflow, "bravo-workflow", LifecycleStatePublished, RiskClassificationMedium),
	}
	for _, a := range assetsBravo {
		if err := store.Put(ctx, a); err != nil {
			t.Fatalf("failed to seed bravo asset: %v", err)
		}
	}

	svc := NewCatalogService(store)
	return store, svc
}

// Full admin auth context with all permissions
func adminAuth(tenantID string) AuthContext {
	return NewAuthContext(
		"usr-admin",
		tenantID,
		[]string{tenantID},
		PermStudioLibraryRead,
		PermStudioSkillRead,
		PermStudioSkillManage,
		PermStudioWorkflowRead,
		PermStudioWorkflowManage,
		PermStudioAutomationAgentRead,
		PermStudioAutomationAgentManage,
		PermStudioAgentTeamRead,
		PermStudioAgentTeamManage,
		PermCmdrStudioVersionRead,
		PermCmdrStudioVersionManage,
	)
}

// ---------------------------------------------------------------------
// 1. Browse Tests
// ---------------------------------------------------------------------

func TestBrowse_NominalPagination(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Page 1: limit 3
	resp1, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     3,
	})
	if err != nil {
		t.Fatalf("Browse page 1 failed: %v", err)
	}
	if len(resp1.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(resp1.Items))
	}
	if !resp1.Pagination.HasMore {
		t.Fatalf("expected HasMore to be true")
	}
	if resp1.Pagination.CursorNext == "" {
		t.Fatalf("expected non-empty CursorNext")
	}
	if resp1.Pagination.TotalCount != 10 {
		t.Fatalf("expected TotalCount 10, got %d", resp1.Pagination.TotalCount)
	}

	// Page 2: use CursorNext
	resp2, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     3,
		Cursor:    resp1.Pagination.CursorNext,
	})
	if err != nil {
		t.Fatalf("Browse page 2 failed: %v", err)
	}
	if len(resp2.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(resp2.Items))
	}

	// Verify items do not overlap between page 1 and page 2
	ids1 := make(map[string]bool)
	for _, it := range resp1.Items {
		ids1[it.ID] = true
	}
	for _, it := range resp2.Items {
		if ids1[it.ID] {
			t.Fatalf("overlapping item %s found across pages", it.ID)
		}
	}
}

func TestBrowse_CompoundSorting_NameASC(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef:     "tenant-alpha",
		Limit:         50,
		SortField:     "name",
		SortDirection: "asc",
	})
	if err != nil {
		t.Fatalf("Browse sorted by name failed: %v", err)
	}

	for i := 1; i < len(resp.Items); i++ {
		prev := strings.ToLower(resp.Items[i-1].Name)
		curr := strings.ToLower(resp.Items[i].Name)
		if prev > curr {
			t.Fatalf("items not sorted by name asc: %s > %s", prev, curr)
		}
	}
}

func TestBrowse_CompoundSorting_IDAsc(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef:     "tenant-alpha",
		Limit:         50,
		SortField:     "id",
		SortDirection: "asc",
	})
	if err != nil {
		t.Fatalf("Browse sorted by id failed: %v", err)
	}

	for i := 1; i < len(resp.Items); i++ {
		if resp.Items[i-1].ID > resp.Items[i].ID {
			t.Fatalf("items not sorted by id asc: %s > %s", resp.Items[i-1].ID, resp.Items[i].ID)
		}
	}
}

func TestBrowse_LimitClamping(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Limit 0 should default to 25
	respDefault, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     0,
	})
	if err != nil {
		t.Fatalf("Browse default limit failed: %v", err)
	}
	if respDefault.Pagination.Limit != 25 {
		t.Fatalf("expected default limit 25, got %d", respDefault.Pagination.Limit)
	}

	// Limit 500 should be clamped to 100
	respClamped, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     500,
	})
	if err != nil {
		t.Fatalf("Browse clamped limit failed: %v", err)
	}
	if respClamped.Pagination.Limit != 100 {
		t.Fatalf("expected clamped limit 100, got %d", respClamped.Pagination.Limit)
	}
}

func TestBrowse_MissingTenant(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "",
	})
	if err == nil {
		t.Fatalf("expected error for missing tenant")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeMissingTenant || errEnv.HTTPStatus != 400 {
		t.Fatalf("expected 400 missing-tenant, got %v", err)
	}
}

func TestBrowse_WildcardTenant(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	wildcards := []string{"*", "%", "all", "ALL"}
	for _, wc := range wildcards {
		_, err := svc.Browse(ctx, auth, BrowseRequest{
			TenantRef: wc,
		})
		if err == nil {
			t.Fatalf("expected error for wildcard tenant %q", wc)
		}
		errEnv, ok := err.(*ErrorEnvelope)
		if !ok || errEnv.Code != ErrorCodeWildcardTenant || errEnv.HTTPStatus != 400 {
			t.Fatalf("expected 400 wildcard-tenant for %q, got %v", wc, err)
		}
	}
}

func TestBrowse_MissingLibraryReadPermission(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Caller lacks perm.studio.library.read
	auth := NewAuthContext("usr-restricted", "tenant-alpha", []string{"tenant-alpha"}, PermStudioSkillRead)

	_, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
	})
	if err == nil {
		t.Fatalf("expected permission-denied error")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodePermissionDenied || errEnv.HTTPStatus != 403 {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}
}

func TestBrowse_SilentOmissionOfRestrictedKinds(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Caller has library.read and skill.read only (no workflow, agent, agent-team)
	auth := NewAuthContext(
		"usr-skill-only",
		"tenant-alpha",
		[]string{"tenant-alpha"},
		PermStudioLibraryRead,
		PermStudioSkillRead,
	)

	resp, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     50,
	})
	if err != nil {
		t.Fatalf("Browse failed: %v", err)
	}

	// Should only see skills and tools (tools are governed by library.read)
	for _, it := range resp.Items {
		if it.Type != AssetKindSkill && it.Type != AssetKindTool {
			t.Fatalf("unpermitted asset kind %s leaked in browse list", it.Type)
		}
	}
}

func TestBrowse_InvalidCursorTampered(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Cursor:    "tampered-invalid-base64",
	})
	if err == nil {
		t.Fatalf("expected error for invalid cursor")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeInvalidCursor || errEnv.HTTPStatus != 400 {
		t.Fatalf("expected 400 invalid-cursor, got %v", err)
	}
}

func TestBrowse_CursorTenantMismatch(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Create cursor for tenant-bravo
	cur := EncodeCursor(CursorPayload{
		Version:       1,
		TenantID:      "tenant-bravo",
		SortField:     "updated_at",
		SortDirection: "desc",
		TiebreakerID:  "skl-01",
	})

	_, err := svc.Browse(ctx, auth, BrowseRequest{
		TenantRef: "tenant-alpha",
		Cursor:    cur,
	})
	if err == nil {
		t.Fatalf("expected error for tenant mismatch in cursor")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeInvalidCursor {
		t.Fatalf("expected invalid-cursor, got %v", err)
	}
}

// ---------------------------------------------------------------------
// 2. Search Tests
// ---------------------------------------------------------------------

func TestSearch_NominalDeterministic(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Search(ctx, auth, SearchRequest{
		TenantRef:             "tenant-alpha",
		QueryText:             "triage",
		SharedSearchAvailable: true,
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(resp.Items) == 0 {
		t.Fatalf("expected search matches for 'triage', got 0")
	}
	for _, it := range resp.Items {
		hasMatch := strings.Contains(strings.ToLower(it.Name), "triage") ||
			strings.Contains(strings.ToLower(it.DisplayName), "triage") ||
			strings.Contains(strings.ToLower(it.Description), "triage")
		if !hasMatch {
			for _, tag := range it.Tags {
				if strings.Contains(strings.ToLower(tag), "triage") {
					hasMatch = true
					break
				}
			}
		}
		if !hasMatch {
			t.Fatalf("item %s returned but does not match query", it.Name)
		}
	}
}

func TestSearch_EmptyQueryRejected(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Search(ctx, auth, SearchRequest{
		TenantRef: "tenant-alpha",
		QueryText: "   ",
	})
	if err == nil {
		t.Fatalf("expected error for empty search query")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeInvalidFilterParameter {
		t.Fatalf("expected invalid-filter-parameter, got %v", err)
	}
}

func TestSearch_DegradedFallbackFunctionalStatePartial(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Search(ctx, auth, SearchRequest{
		TenantRef:             "tenant-alpha",
		QueryText:             "phishing",
		SharedSearchAvailable: false, // Simulates shared search outage
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(resp.Items) == 0 {
		t.Fatalf("expected matches for 'phishing'")
	}
	for _, it := range resp.Items {
		if it.FunctionalState != FunctionalStatePartial {
			t.Fatalf("expected functional_state 'partial' during search degradation, got %s", it.FunctionalState)
		}
	}
}

// ---------------------------------------------------------------------
// 3. Filter Tests
// ---------------------------------------------------------------------

func TestFilter_ByKind(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			Kinds: []AssetKind{AssetKindWorkflow},
		},
	})
	if err != nil {
		t.Fatalf("Filter by kind failed: %v", err)
	}
	if len(resp.Items) == 0 {
		t.Fatalf("expected workflow items")
	}
	for _, it := range resp.Items {
		if it.Type != AssetKindWorkflow {
			t.Fatalf("expected only workflows, got %s", it.Type)
		}
	}
}

func TestFilter_ByLifecycleState(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			LifecycleStates: []LifecycleState{LifecycleStateDraft},
		},
	})
	if err != nil {
		t.Fatalf("Filter by lifecycle failed: %v", err)
	}
	if len(resp.Items) == 0 {
		t.Fatalf("expected draft items")
	}
	for _, it := range resp.Items {
		if it.LifecycleState != LifecycleStateDraft {
			t.Fatalf("expected only draft, got %s", it.LifecycleState)
		}
	}
}

func TestFilter_ByRiskClassification(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			RiskLevels: []RiskClassification{RiskClassificationCritical},
		},
	})
	if err != nil {
		t.Fatalf("Filter by risk failed: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].ID != "agn-02" {
		t.Fatalf("expected 1 critical item (agn-02), got %d", len(resp.Items))
	}
}

func TestFilter_ByOwnerUnknown(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			OwnerID: "unknown",
		},
	})
	if err != nil {
		t.Fatalf("Filter by unknown owner failed: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].ID != "skl-no-owner" {
		t.Fatalf("expected 1 orphan skill item, got %d", len(resp.Items))
	}
	if !resp.Items[0].Owner.IsUnknown || resp.Items[0].Owner.DisplayName != "owner unknown" {
		t.Fatalf("expected owner unknown representation, got %+v", resp.Items[0].Owner)
	}
}

func TestFilter_MultiCriteriaConjunction(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	aiTrue := true
	resp, err := svc.Filter(ctx, auth, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			Kinds:           []AssetKind{AssetKindSkill},
			LifecycleStates: []LifecycleState{LifecycleStatePublished},
			RiskLevels:      []RiskClassification{RiskClassificationLow},
			IsAIPowered:     &aiTrue,
			Tags:            []string{"security"},
		},
	})
	if err != nil {
		t.Fatalf("Multi-criteria filter failed: %v", err)
	}
	if len(resp.Items) == 0 {
		t.Fatalf("expected matching skill items")
	}
	for _, it := range resp.Items {
		if it.Type != AssetKindSkill || it.LifecycleState != LifecycleStatePublished || it.RiskClassification != RiskClassificationLow || !it.IsAIPowered {
			t.Fatalf("item does not satisfy all filter conjunctions: %+v", it)
		}
	}
}

// ---------------------------------------------------------------------
// 4. Group Tests
// ---------------------------------------------------------------------

func TestGroup_ByKind(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Group(ctx, auth, GroupRequest{
		TenantRef: "tenant-alpha",
		GroupBy:   FacetKind,
	})
	if err != nil {
		t.Fatalf("Group by kind failed: %v", err)
	}
	if resp.Facet != FacetKind {
		t.Fatalf("expected facet by_kind, got %s", resp.Facet)
	}
	if resp.TotalItems != 10 {
		t.Fatalf("expected total 10 items, got %d", resp.TotalItems)
	}

	bucketMap := make(map[string]int)
	for _, b := range resp.Groups {
		bucketMap[b.Key] = b.Count
	}
	if bucketMap["skill"] != 3 {
		t.Fatalf("expected 3 skills, got %d", bucketMap["skill"])
	}
	if bucketMap["workflow"] != 2 {
		t.Fatalf("expected 2 workflows, got %d", bucketMap["workflow"])
	}
	if bucketMap["tool"] != 2 {
		t.Fatalf("expected 2 tools, got %d", bucketMap["tool"])
	}
}

func TestGroup_ByLifecycleState(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Group(ctx, auth, GroupRequest{
		TenantRef: "tenant-alpha",
		GroupBy:   FacetLifecycleState,
	})
	if err != nil {
		t.Fatalf("Group by lifecycle failed: %v", err)
	}
	if len(resp.Groups) == 0 {
		t.Fatalf("expected groups")
	}
}

func TestGroup_ByRiskClassification(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Group(ctx, auth, GroupRequest{
		TenantRef: "tenant-alpha",
		GroupBy:   FacetRiskClassification,
	})
	if err != nil {
		t.Fatalf("Group by risk failed: %v", err)
	}
	bucketMap := make(map[string]int)
	for _, b := range resp.Groups {
		bucketMap[b.Key] = b.Count
	}
	if bucketMap["critical"] != 1 {
		t.Fatalf("expected 1 critical item, got %d", bucketMap["critical"])
	}
}

func TestGroup_UnsupportedFacet(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Group(ctx, auth, GroupRequest{
		TenantRef: "tenant-alpha",
		GroupBy:   "by_invalid_facet",
	})
	if err == nil {
		t.Fatalf("expected error for unsupported grouping facet")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeUnsupportedGroupingFacet || errEnv.HTTPStatus != 400 {
		t.Fatalf("expected 400 unsupported-grouping-facet, got %v", err)
	}
}

// ---------------------------------------------------------------------
// 5. Open Tests
// ---------------------------------------------------------------------

func TestOpen_Nominal(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Open(ctx, auth, OpenRequest{
		TenantRef:    "tenant-alpha",
		AssetID:      "skl-01",
		ReturnOrigin: "/studio/library?tab=skills",
	})
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if resp.TargetSurface != "builder" {
		t.Fatalf("expected builder surface, got %s", resp.TargetSurface)
	}
	if resp.SurfaceRoute != "/studio/builder/skill/skl-01" {
		t.Fatalf("unexpected surface route: %s", resp.SurfaceRoute)
	}
	if resp.ReturnOrigin != "/studio/library?tab=skills" {
		t.Fatalf("return origin corrupted: %s", resp.ReturnOrigin)
	}
}

func TestOpen_MissingReturnOrigin(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Open(ctx, auth, OpenRequest{
		TenantRef:    "tenant-alpha",
		AssetID:      "skl-01",
		ReturnOrigin: "",
	})
	if err == nil {
		t.Fatalf("expected error for missing return origin")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeMissingReturnOrigin {
		t.Fatalf("expected missing-return-origin, got %v", err)
	}
}

func TestOpen_NonExistentAsset_404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Open(ctx, auth, OpenRequest{
		TenantRef:    "tenant-alpha",
		AssetID:      "skl-does-not-exist",
		ReturnOrigin: "/studio/library",
	})
	if err == nil {
		t.Fatalf("expected 404 for non-existent asset")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeAssetNotFound || errEnv.HTTPStatus != 404 {
		t.Fatalf("expected 404 asset-not-found, got %v", err)
	}
}

func TestOpen_CrossTenant_AntiOracle_Uniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Scoped to tenant-alpha, attempts to open skl-bravo-secret from tenant-bravo
	auth := adminAuth("tenant-alpha")

	_, err := svc.Open(ctx, auth, OpenRequest{
		TenantRef:    "tenant-alpha",
		AssetID:      "skl-bravo-secret",
		ReturnOrigin: "/studio/library",
	})
	if err == nil {
		t.Fatalf("expected error for cross-tenant open")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeAssetNotFound || errEnv.HTTPStatus != 404 {
		t.Fatalf("expected uniform 404 asset-not-found (no 403 oracle), got %v", err)
	}
}

// ---------------------------------------------------------------------
// 6. Inspect Tests (Existence Oracle Elimination & Sanitized Stub)
// ---------------------------------------------------------------------

func TestInspect_NominalFullSchema(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}
	if resp.Asset == nil || resp.Asset.ID != "skl-01" {
		t.Fatalf("expected asset skl-01, got nil or wrong id")
	}
	if len(resp.Asset.Parameters) == 0 {
		t.Fatalf("expected parameters to be populated")
	}
	if len(resp.Asset.Dependencies) == 0 {
		t.Fatalf("expected dependencies to be populated")
	}
	if resp.Asset.FunctionalState != FunctionalStateReady {
		t.Fatalf("expected ready functional state, got %s", resp.Asset.FunctionalState)
	}
}

func TestInspect_NonExistentAsset_404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-nonexistent-id",
	})
	if err == nil {
		t.Fatalf("expected 404 for non-existent asset")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeAssetNotFound || errEnv.HTTPStatus != 404 {
		t.Fatalf("expected 404 asset-not-found, got %v", err)
	}
}

func TestInspect_CrossTenant_AntiOracle_Uniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Probing asset 'skl-bravo-secret' from tenant-bravo while in tenant-alpha
	_, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-secret",
	})
	if err == nil {
		t.Fatalf("expected 404 for foreign tenant asset")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	// Must return 404 asset-not-found, NEVER 403 cross-tenant
	if !ok || errEnv.Code != ErrorCodeAssetNotFound || errEnv.HTTPStatus != 404 {
		t.Fatalf("expected uniform 404 asset-not-found, got %v", err)
	}
	if !strings.Contains(errEnv.Message, "was not found in tenant") {
		t.Fatalf("unexpected message format: %s", errEnv.Message)
	}
}

func TestInspect_SameTenant_Unpermitted_SanitizedStub200(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Caller has library.read and skill.read, but lacks workflow.read
	auth := NewAuthContext(
		"usr-analyst",
		"tenant-alpha",
		[]string{"tenant-alpha"},
		PermStudioLibraryRead,
		PermStudioSkillRead,
	)

	// Inspect workflow wfl-01 in same tenant
	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wfl-01",
	})
	if err != nil {
		t.Fatalf("Inspect unpermitted same-tenant should succeed with 200 sanitized stub, got: %v", err)
	}
	if resp.Asset == nil {
		t.Fatalf("expected sanitized asset stub, got nil")
	}
	// Must be permission-filtered
	if resp.Asset.FunctionalState != FunctionalStatePermissionFiltered {
		t.Fatalf("expected functional_state 'permission-filtered', got %s", resp.Asset.FunctionalState)
	}
	// Description, parameters, dependencies must be redacted
	if resp.Asset.Description != "" {
		t.Fatalf("description was not redacted")
	}
	if len(resp.Asset.Parameters) != 0 {
		t.Fatalf("parameters were not redacted")
	}
	if len(resp.Asset.Dependencies) != 0 {
		t.Fatalf("dependencies were not redacted")
	}
	if !resp.Asset.Owner.IsUnknown {
		t.Fatalf("owner was not sanitized to unknown")
	}
}

// ---------------------------------------------------------------------
// 7. Compare Tests
// ---------------------------------------------------------------------

func TestCompare_NominalStructuralDiff(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Create updated version of skl-01
	sklV2 := makeTestAsset("tenant-alpha", "skl-01-v2", AssetKindSkill, "triage-alert-v2", LifecycleStatePublished, RiskClassificationMedium)
	sklV2.Version = VersionInfo{VersionID: "ver-2.0.0", SemVer: "2.0.0", CreatedAt: time.Now().UTC(), IsImmutable: true}
	sklV2.Parameters = append(sklV2.Parameters, Parameter{Name: "enrich_geo", Type: "bool", Required: false})
	sklV2.Dependencies = []Dependency{
		{Kind: "tool", RefID: "tool-whois", VersionConstraint: "^2.0.0", IsOptional: false},
	}
	if err := store.Put(ctx, sklV2); err != nil {
		t.Fatalf("failed to put sklV2: %v", err)
	}

	resp, err := svc.Compare(ctx, auth, CompareRequest{
		TenantRef:     "tenant-alpha",
		BaseAssetID:   "skl-01",
		TargetAssetID: "skl-01-v2",
	})
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	diff := resp.Diff
	if !diff.RiskClassificationChanged || diff.OldRisk != RiskClassificationLow || diff.NewRisk != RiskClassificationMedium {
		t.Fatalf("expected risk classification change from low to medium")
	}
	if len(diff.ParametersAdded) != 1 || diff.ParametersAdded[0].Name != "enrich_geo" {
		t.Fatalf("expected enrich_geo parameter added, got %+v", diff.ParametersAdded)
	}
	if len(diff.DependenciesChanged) != 1 || diff.DependenciesChanged[0].OldConstraint != "^1.0.0" || diff.DependenciesChanged[0].NewConstraint != "^2.0.0" {
		t.Fatalf("expected tool-whois constraint change from ^1.0.0 to ^2.0.0, got %+v", diff.DependenciesChanged)
	}
}

func TestCompare_CrossTenant_404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.Compare(ctx, auth, CompareRequest{
		TenantRef:     "tenant-alpha",
		BaseAssetID:   "skl-01",
		TargetAssetID: "skl-bravo-secret", // Belongs to tenant-bravo
	})
	if err == nil {
		t.Fatalf("expected 404 comparing cross-tenant asset")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeAssetNotFound || errEnv.HTTPStatus != 404 {
		t.Fatalf("expected 404 asset-not-found, got %v", err)
	}
}

// ---------------------------------------------------------------------
// 8. Clone Preparation Tests
// ---------------------------------------------------------------------

func TestClonePreparation_NominalDraft(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-alpha",
		SourceAssetID:   "skl-01",
		ClonedBy:        "usr-analyst",
	})
	if err != nil {
		t.Fatalf("ClonePreparation failed: %v", err)
	}

	clone := resp.DraftClone
	if clone.DraftAssetID == "" || clone.DraftAssetID == "skl-01" {
		t.Fatalf("expected new unique draft ID, got %s", clone.DraftAssetID)
	}
	if clone.ClonedFromID != "skl-01" {
		t.Fatalf("expected ClonedFromID skl-01, got %s", clone.ClonedFromID)
	}
	if clone.PrefilledMetadata.LifecycleState != LifecycleStateDraft {
		t.Fatalf("expected draft lifecycle state, got %s", clone.PrefilledMetadata.LifecycleState)
	}
	if clone.PrefilledMetadata.Version.SemVer != "0.1.0-draft" {
		t.Fatalf("expected draft semver 0.1.0-draft, got %s", clone.PrefilledMetadata.Version.SemVer)
	}
	if !strings.Contains(clone.PrefilledMetadata.Name, "-copy") {
		t.Fatalf("expected copy suffix in name: %s", clone.PrefilledMetadata.Name)
	}
	// Verify secret dependencies are stripped or masked
	for _, dep := range clone.PrefilledMetadata.Dependencies {
		if dep.Kind == "secret_reference" && !dep.IsMasked {
			t.Fatalf("secret reference was not masked in clone descriptor")
		}
	}
	// Verify provenance chain includes original
	if len(clone.ProvenanceChain) == 0 || !strings.Contains(clone.ProvenanceChain[len(clone.ProvenanceChain)-1], "skl-01@1.0.0") {
		t.Fatalf("provenance chain missing source lineage: %+v", clone.ProvenanceChain)
	}
}

func TestClonePreparation_CrossTenantMutationForbidden_403(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	_, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-bravo", // Cross-tenant target!
		SourceAssetID:   "skl-01",
	})
	if err == nil {
		t.Fatalf("expected error for cross-tenant clone attempt")
	}
	errEnv, ok := err.(*ErrorEnvelope)
	if !ok || errEnv.Code != ErrorCodeCrossTenantMutationForbidden || errEnv.HTTPStatus != 403 {
		t.Fatalf("expected 403 cross-tenant-mutation-forbidden, got %v", err)
	}
}

// ---------------------------------------------------------------------
// 9. Execution Blocking (REQ-PROD-006) Tests
// ---------------------------------------------------------------------

func TestExecutionBlocking_AllVerbsBlocked(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	verbs := []struct {
		name string
		fn   func() error
	}{
		{"Invoke", func() error { return svc.Invoke(ctx, auth, nil) }},
		{"Execute", func() error { return svc.Execute(ctx, auth, nil) }},
		{"Run", func() error { return svc.Run(ctx, auth, nil) }},
		{"Call", func() error { return svc.Call(ctx, auth, nil) }},
		{"Trigger", func() error { return svc.Trigger(ctx, auth, nil) }},
	}

	for _, v := range verbs {
		err := v.fn()
		if err == nil {
			t.Fatalf("expected execution verb %s to be blocked", v.name)
		}
		errEnv, ok := err.(*ErrorEnvelope)
		if !ok || errEnv.Code != ErrorCodeExecutionBlocked || errEnv.HTTPStatus != 400 {
			t.Fatalf("expected 400 execution-blocked for %s, got %v", v.name, err)
		}
	}
}

// ---------------------------------------------------------------------
// 10. Multi-Tenant MSSP ATS Validation Tests
// ---------------------------------------------------------------------

func TestValidateMSSPTenantSet_NominalSubset(t *testing.T) {
	authorized := []string{"tenant-1", "tenant-2", "tenant-3"}
	requested := []string{"tenant-1", "tenant-3"}

	err := ValidateMSSPTenantSet(requested, authorized, "test-corr")
	if err != nil {
		t.Fatalf("expected valid subset, got %v", err)
	}
}

func TestValidateMSSPTenantSet_UnauthorizedTenantInSet_403(t *testing.T) {
	authorized := []string{"tenant-1", "tenant-2"}
	requested := []string{"tenant-1", "tenant-999"} // tenant-999 is unauthorized

	err := ValidateMSSPTenantSet(requested, authorized, "test-corr")
	if err == nil {
		t.Fatalf("expected fail-closed error for unauthorized tenant in set")
	}
	if err.Code != ErrorCodeUnauthorizedTenantInSet || err.HTTPStatus != 403 {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got %v", err)
	}
}

func TestValidateTenantScope_WithATS(t *testing.T) {
	ats := []string{"tenant-1", "tenant-2"}

	// Allowed
	if err := ValidateTenantScope("tenant-1", ats, "test-corr"); err != nil {
		t.Fatalf("expected tenant-1 to be allowed in ATS, got %v", err)
	}

	// Disallowed
	err := ValidateTenantScope("tenant-3", ats, "test-corr")
	if err == nil || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected unauthorized-tenant-in-set for tenant-3, got %v", err)
	}
}

// ---------------------------------------------------------------------
// 11. Canonical Owner Unknown Fallback (REQ-OBJ-009) Tests
// ---------------------------------------------------------------------

func TestOwnerUnknown_Normalization(t *testing.T) {
	tests := []struct {
		name  string
		input Owner
	}{
		{"empty struct", Owner{}},
		{"empty strings", Owner{OwnerID: "  ", OwnerName: "  "}},
		{"is_unknown set", Owner{OwnerID: "usr-1", OwnerName: "Alice", IsUnknown: true}},
	}

	for _, tc := range tests {
		o := tc.input
		o.Normalize()
		if !o.IsUnknown {
			t.Fatalf("[%s] expected IsUnknown true", tc.name)
		}
		if o.DisplayName != "owner unknown" {
			t.Fatalf("[%s] expected DisplayName 'owner unknown', got %s", tc.name, o.DisplayName)
		}
		if o.OwnerType != "unknown" {
			t.Fatalf("[%s] expected OwnerType 'unknown', got %s", tc.name, o.OwnerType)
		}
	}
}

// ---------------------------------------------------------------------
// 12. Tool Asset Governance by Library Read Tests
// ---------------------------------------------------------------------

func TestToolAsset_GovernedByLibraryRead(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Caller only has library.read (no skill, workflow, or tool permission)
	auth := NewAuthContext("usr-tool-user", "tenant-alpha", []string{"tenant-alpha"}, PermStudioLibraryRead)

	// Can inspect tool-01
	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "tool-01",
	})
	if err != nil {
		t.Fatalf("expected tool inspection to succeed under library.read: %v", err)
	}
	if resp.Asset.ID != "tool-01" || resp.Asset.FunctionalState != FunctionalStateReady {
		t.Fatalf("expected tool-01 ready, got %+v", resp.Asset)
	}
}

// ---------------------------------------------------------------------
// 13. Dual Namespace Permission Compatibility Tests
// ---------------------------------------------------------------------

func TestPermissionAliases_DualNamespace(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// Caller holds legacy cmdr-studio namespace: perm.cmdr-studio.skill.read
	auth := NewAuthContext(
		"usr-legacy",
		"tenant-alpha",
		[]string{"tenant-alpha"},
		PermStudioLibraryRead,
		"perm.cmdr-studio.skill.read",
	)

	// Inspecting skill should succeed via alias
	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("Inspect failed with legacy namespace alias: %v", err)
	}
	if resp.Asset.FunctionalState != FunctionalStateReady {
		t.Fatalf("expected ready state via alias, got %s", resp.Asset.FunctionalState)
	}
}

// ---------------------------------------------------------------------
// 14. Projection Immutability Tests
// ---------------------------------------------------------------------

func TestProjectionImmutability_DeepCopyIsolation(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp1, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	// Mutate returned asset in-place
	resp1.Asset.Name = "maliciously-mutated-name"
	resp1.Asset.Tags = append(resp1.Asset.Tags, "injected-tag")

	// Re-inspect should return pristine original from store
	resp2, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("Re-inspect failed: %v", err)
	}
	if resp2.Asset.Name == "maliciously-mutated-name" {
		t.Fatalf("deep copy failed: store asset was mutated by caller")
	}
	for _, tag := range resp2.Asset.Tags {
		if tag == "injected-tag" {
			t.Fatalf("injected tag persisted in store")
		}
	}
}
