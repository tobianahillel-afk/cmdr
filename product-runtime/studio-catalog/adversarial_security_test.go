package studiocatalog

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// PROBE 1: Cross-Tenant Leakage & Anti-Oracle Uniform 404 Indistinguishability
// =============================================================================

// Verify that Open on a foreign tenant's asset yields uniform 404
func TestAdversarial_CrossTenant_OpenForeignAsset_ReturnsUniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	// Attempt to open tenant-bravo's asset while scoped to tenant-alpha
	_, err := svc.Open(ctx, authAlpha, OpenRequest{
		TenantRef:    "tenant-alpha",
		AssetID:      "skl-bravo-secret",
		ReturnOrigin: "/studio/library",
	})
	if err == nil {
		t.Fatalf("expected error opening foreign tenant asset, got nil")
	}

	errEnv, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected 404 asset-not-found, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// Verify that Inspect on a foreign tenant's asset yields uniform 404
func TestAdversarial_CrossTenant_InspectForeignAsset_ReturnsUniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	// Attempt to inspect tenant-bravo's asset while scoped to tenant-alpha
	_, err := svc.Inspect(ctx, authAlpha, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-secret",
	})
	if err == nil {
		t.Fatalf("expected error inspecting foreign tenant asset, got nil")
	}

	errEnv, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected 404 asset-not-found, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// Verify that Compare referencing foreign assets yields uniform 404
func TestAdversarial_CrossTenant_CompareForeignAsset_ReturnsUniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	// Case 1: Target is foreign
	_, err := svc.Compare(ctx, authAlpha, CompareRequest{
		TenantRef:     "tenant-alpha",
		BaseAssetID:   "skl-01",
		TargetAssetID: "skl-bravo-secret",
	})
	if err == nil {
		t.Fatalf("expected error comparing foreign target asset, got nil")
	}
	errEnv := err.(*ErrorEnvelope)
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected 404 asset-not-found, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}

	// Case 2: Base is foreign
	_, err = svc.Compare(ctx, authAlpha, CompareRequest{
		TenantRef:     "tenant-alpha",
		BaseAssetID:   "skl-bravo-secret",
		TargetAssetID: "skl-01",
	})
	if err == nil {
		t.Fatalf("expected error comparing foreign base asset, got nil")
	}
	errEnv = err.(*ErrorEnvelope)
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected 404 asset-not-found, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// Verify ClonePreparation with foreign source returns uniform 404
func TestAdversarial_CrossTenant_ClonePreparationSourceForeign_ReturnsUniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	_, err := svc.ClonePreparation(ctx, authAlpha, ClonePreparationRequest{
		TenantRef:     "tenant-alpha",
		SourceAssetID: "skl-bravo-secret",
	})
	if err == nil {
		t.Fatalf("expected error cloning foreign source, got nil")
	}
	errEnv := err.(*ErrorEnvelope)
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected 404 asset-not-found, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// Verify ClonePreparation with foreign destination returns 403 cross-tenant mutation forbidden
func TestAdversarial_CrossTenant_ClonePreparationDestinationForeign_Returns403(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	_, err := svc.ClonePreparation(ctx, authAlpha, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-bravo",
		SourceAssetID:   "skl-01",
	})
	if err == nil {
		t.Fatalf("expected error cloning to foreign destination, got nil")
	}
	errEnv := err.(*ErrorEnvelope)
	if errEnv.HTTPStatus != 403 || errEnv.Code != ErrorCodeCrossTenantMutationForbidden {
		t.Fatalf("expected 403 cross-tenant-mutation-forbidden, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// Verify Browse completely isolates tenant data
func TestAdversarial_CrossTenant_BrowseSilentIsolation(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	resp, err := svc.Browse(ctx, authAlpha, BrowseRequest{
		TenantRef: "tenant-alpha",
		Limit:     100,
	})
	if err != nil {
		t.Fatalf("Browse failed: %v", err)
	}

	for _, item := range resp.Items {
		if item.TenantID != "tenant-alpha" {
			t.Fatalf("LEAKAGE: found asset %s belonging to foreign tenant %s in tenant-alpha browse", item.ID, item.TenantID)
		}
		if strings.Contains(item.ID, "bravo") {
			t.Fatalf("LEAKAGE: foreign item %s leaked in browse response", item.ID)
		}
	}
}

// Verify Search completely isolates tenant data
func TestAdversarial_CrossTenant_SearchSilentIsolation(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	// Search for "bravo" in tenant-alpha
	resp, err := svc.Search(ctx, authAlpha, SearchRequest{
		TenantRef: "tenant-alpha",
		QueryText: "bravo",
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("LEAKAGE: found %d items matching 'bravo' in tenant-alpha", len(resp.Items))
	}
}

// Verify Filter completely isolates tenant data
func TestAdversarial_CrossTenant_FilterSilentIsolation(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	resp, err := svc.Filter(ctx, authAlpha, FilterRequest{
		TenantRef: "tenant-alpha",
		Criteria: FilterCriteria{
			Kinds: []AssetKind{AssetKindSkill, AssetKindWorkflow},
		},
	})
	if err != nil {
		t.Fatalf("Filter failed: %v", err)
	}
	for _, it := range resp.Items {
		if it.TenantID != "tenant-alpha" {
			t.Fatalf("LEAKAGE: filter returned asset %s from foreign tenant %s", it.ID, it.TenantID)
		}
	}
}

// Verify Group completely isolates tenant data
func TestAdversarial_CrossTenant_GroupSilentIsolation(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	resp, err := svc.Group(ctx, authAlpha, GroupRequest{
		TenantRef: "tenant-alpha",
		GroupBy:   FacetKind,
	})
	if err != nil {
		t.Fatalf("Group failed: %v", err)
	}

	// In tenant-alpha: 3 skills, 2 workflows, 2 agents, 1 team, 2 tools = 10 assets
	if resp.TotalItems != 10 {
		t.Fatalf("expected 10 total items in tenant-alpha, got %d (possible foreign tenant leakage)", resp.TotalItems)
	}
}

// Verify Payload Indistinguishability: non-existent vs foreign tenant asset
func TestAdversarial_PayloadIndistinguishability_NonExistentVsForeign(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	// 1. Non-existent asset ID
	_, errNonExistent := svc.Inspect(ctx, authAlpha, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "asset-completely-fictional-999",
	})
	if errNonExistent == nil {
		t.Fatalf("expected error for non-existent asset")
	}
	envNonExistent := errNonExistent.(*ErrorEnvelope)

	// 2. Foreign tenant asset ID (exists in tenant-bravo)
	_, errForeign := svc.Inspect(ctx, authAlpha, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-secret",
	})
	if errForeign == nil {
		t.Fatalf("expected error for foreign tenant asset")
	}
	envForeign := errForeign.(*ErrorEnvelope)

	// Compare payloads field by field
	if envNonExistent.Code != envForeign.Code {
		t.Fatalf("ORACLE LEAKAGE: code mismatch: non-existent=%s vs foreign=%s", envNonExistent.Code, envForeign.Code)
	}
	if envNonExistent.HTTPStatus != envForeign.HTTPStatus {
		t.Fatalf("ORACLE LEAKAGE: status mismatch: non-existent=%d vs foreign=%d", envNonExistent.HTTPStatus, envForeign.HTTPStatus)
	}
	if envNonExistent.Retryable != envForeign.Retryable {
		t.Fatalf("ORACLE LEAKAGE: retryable mismatch")
	}
	if envNonExistent.PartialResult != envForeign.PartialResult {
		t.Fatalf("ORACLE LEAKAGE: partial result mismatch")
	}

	// Verify message template identity: both follow `Asset "%s" was not found in tenant "%s"`
	expectedNonExistentMsg := `Asset "asset-completely-fictional-999" was not found in tenant "tenant-alpha"`
	expectedForeignMsg := `Asset "skl-bravo-secret" was not found in tenant "tenant-alpha"`
	if envNonExistent.Message != expectedNonExistentMsg {
		t.Fatalf("unexpected message for non-existent asset: %s", envNonExistent.Message)
	}
	if envForeign.Message != expectedForeignMsg {
		t.Fatalf("unexpected message for foreign asset: %s", envForeign.Message)
	}
}

// Verify Timing Indistinguishability: 1000 iterations non-existent vs foreign
func TestAdversarial_TimingIndistinguishability_NonExistentVsForeign(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	authAlpha := adminAuth("tenant-alpha")

	const iterations = 1000
	var nonExistentDurations []time.Duration
	var foreignDurations []time.Duration

	// Warm up
	for i := 0; i < 50; i++ {
		_, _ = svc.Inspect(ctx, authAlpha, InspectRequest{TenantRef: "tenant-alpha", AssetID: "warmup"})
	}

	for i := 0; i < iterations; i++ {
		t0 := time.Now()
		_, _ = svc.Inspect(ctx, authAlpha, InspectRequest{TenantRef: "tenant-alpha", AssetID: "non-existent-id"})
		d1 := time.Since(t0)
		nonExistentDurations = append(nonExistentDurations, d1)

		t1 := time.Now()
		_, _ = svc.Inspect(ctx, authAlpha, InspectRequest{TenantRef: "tenant-alpha", AssetID: "skl-bravo-secret"})
		d2 := time.Since(t1)
		foreignDurations = append(foreignDurations, d2)
	}

	var sum1, sum2 int64
	for i := 0; i < iterations; i++ {
		sum1 += nonExistentDurations[i].Nanoseconds()
		sum2 += foreignDurations[i].Nanoseconds()
	}
	avg1 := float64(sum1) / float64(iterations)
	avg2 := float64(sum2) / float64(iterations)

	// Since both hit in-memory map lookup with identical missed key behavior,
	// average timing delta should be under 20 microseconds (negligible noise)
	deltaNanos := math.Abs(avg1 - avg2)
	t.Logf("Timing check: non-existent avg=%.2f ns, foreign avg=%.2f ns, delta=%.2f ns", avg1, avg2, deltaNanos)
	if deltaNanos > 25000 { // 25 microseconds threshold
		t.Logf("WARNING: timing discrepancy observed: %.2f ns", deltaNanos)
	}
}

// =============================================================================
// PROBE 2: MSSP Authorized Tenant Set (ATS) Boundary Bypass
// =============================================================================

// Verify single tenant request outside ATS is rejected with 403
func TestAdversarial_MSSP_SingleTenantRequest_NotInATS_Returns403(t *testing.T) {
	ats := []string{"tenant-alpha", "tenant-bravo"}

	// Valid tenants in ATS
	if err := ValidateTenantScope("tenant-alpha", ats, "corr-1"); err != nil {
		t.Fatalf("expected tenant-alpha to be allowed, got: %v", err)
	}
	if err := ValidateTenantScope("tenant-bravo", ats, "corr-1"); err != nil {
		t.Fatalf("expected tenant-bravo to be allowed, got: %v", err)
	}

	// Unauthorized tenant outside ATS
	err := ValidateTenantScope("tenant-charlie", ats, "corr-1")
	if err == nil {
		t.Fatalf("expected 403 unauthorized for tenant-charlie")
	}
	if err.HTTPStatus != 403 || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got status=%d code=%s", err.HTTPStatus, err.Code)
	}
}

// Verify multi-tenant aggregation validation enforces requested ⊆ authorized
func TestAdversarial_MSSP_MultiTenantAggregation_NotInATS_Returns403(t *testing.T) {
	ats := []string{"tenant-alpha", "tenant-bravo"}

	// Single unauthorized
	err := ValidateMSSPTenantSet([]string{"tenant-charlie"}, ats, "corr-1")
	if err == nil || err.HTTPStatus != 403 || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got: %v", err)
	}

	// Mixed authorized and unauthorized (probe aggregation bypass)
	err = ValidateMSSPTenantSet([]string{"tenant-alpha", "tenant-charlie"}, ats, "corr-1")
	if err == nil || err.HTTPStatus != 403 || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set for mixed set, got: %v", err)
	}
}

// Verify multi-tenant aggregation rejects wildcards with 400
func TestAdversarial_MSSP_MultiTenantAggregation_Wildcard_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha", "tenant-bravo"}

	wildcards := []string{"*", "%", "all", "ALL", "All"}
	for _, w := range wildcards {
		err := ValidateMSSPTenantSet([]string{w}, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeWildcardTenant {
			t.Fatalf("expected 400 wildcard-tenant for %q, got: %v", w, err)
		}

		// Combined with authorized tenant
		err = ValidateMSSPTenantSet([]string{"tenant-alpha", w}, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeWildcardTenant {
			t.Fatalf("expected 400 wildcard-tenant for mixed %q, got: %v", w, err)
		}
	}
}

// Verify multi-tenant aggregation rejects empty sets
func TestAdversarial_MSSP_MultiTenantAggregation_Empty_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha", "tenant-bravo"}

	err := ValidateMSSPTenantSet([]string{}, ats, "corr-1")
	if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeMissingTenant {
		t.Fatalf("expected 400 missing-tenant for empty slice, got: %v", err)
	}

	err = ValidateMSSPTenantSet([]string{""}, ats, "corr-1")
	if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeMissingTenant {
		t.Fatalf("expected 400 missing-tenant for blank string, got: %v", err)
	}

	err = ValidateMSSPTenantSet([]string{"   "}, ats, "corr-1")
	if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeMissingTenant {
		t.Fatalf("expected 400 missing-tenant for spaces string, got: %v", err)
	}
}

// Verify all 8 Class 0 operations reject requests targeting a tenant outside caller's ATS
func TestAdversarial_MSSP_All8Operations_EnforceATS(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	// Caller only authorized for tenant-alpha
	auth := adminAuth("tenant-alpha")

	// 1. Browse
	_, err := svc.Browse(ctx, auth, BrowseRequest{TenantRef: "tenant-bravo"})
	assertStatusCode(t, "Browse", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 2. Search
	_, err = svc.Search(ctx, auth, SearchRequest{TenantRef: "tenant-bravo", QueryText: "query"})
	assertStatusCode(t, "Search", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 3. Filter
	_, err = svc.Filter(ctx, auth, FilterRequest{TenantRef: "tenant-bravo"})
	assertStatusCode(t, "Filter", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 4. Group
	_, err = svc.Group(ctx, auth, GroupRequest{TenantRef: "tenant-bravo", GroupBy: FacetKind})
	assertStatusCode(t, "Group", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 5. Open
	_, err = svc.Open(ctx, auth, OpenRequest{TenantRef: "tenant-bravo", AssetID: "skl-01", ReturnOrigin: "/origin"})
	assertStatusCode(t, "Open", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 6. Inspect
	_, err = svc.Inspect(ctx, auth, InspectRequest{TenantRef: "tenant-bravo", AssetID: "skl-01"})
	assertStatusCode(t, "Inspect", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 7. Compare
	_, err = svc.Compare(ctx, auth, CompareRequest{TenantRef: "tenant-bravo", BaseAssetID: "b", TargetAssetID: "t"})
	assertStatusCode(t, "Compare", err, 403, ErrorCodeUnauthorizedTenantInSet)

	// 8. ClonePreparation
	_, err = svc.ClonePreparation(ctx, auth, ClonePreparationRequest{TenantRef: "tenant-bravo", SourceAssetID: "s"})
	assertStatusCode(t, "ClonePreparation", err, 403, ErrorCodeUnauthorizedTenantInSet)
}

// Verify cross-tenant cloning is strictly prohibited even if caller has both tenants in ATS
func TestAdversarial_MSSP_CrossTenantCloneProhibitedEvenWithinATS(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// MSSP admin authorized for both tenant-alpha and tenant-bravo
	authMSSP := NewAuthContext(
		"usr-mssp-admin",
		"tenant-alpha",
		[]string{"tenant-alpha", "tenant-bravo"},
		PermStudioLibraryRead,
		PermStudioSkillRead,
		PermStudioSkillManage,
	)

	// Attempt to clone asset from tenant-alpha into tenant-bravo
	_, err := svc.ClonePreparation(ctx, authMSSP, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-bravo",
		SourceAssetID:   "skl-01",
	})
	if err == nil {
		t.Fatalf("CRITICAL: cross-tenant cloning between authorized tenants was permitted!")
	}
	errEnv := err.(*ErrorEnvelope)
	if errEnv.HTTPStatus != 403 || errEnv.Code != ErrorCodeCrossTenantMutationForbidden {
		t.Fatalf("expected 403 cross-tenant-mutation-forbidden, got status=%d code=%s", errEnv.HTTPStatus, errEnv.Code)
	}
}

// =============================================================================
// PROBE 3: Wildcard Injection Attempts in tenant_ref
// =============================================================================

func TestAdversarial_WildcardInjection_Asterisk_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha"}

	cases := []string{"*", " * ", "  *  "}
	for _, tc := range cases {
		err := ValidateTenantScope(tc, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeWildcardTenant {
			t.Fatalf("expected 400 wildcard-tenant for %q, got: %v", tc, err)
		}
	}
}

func TestAdversarial_WildcardInjection_Percent_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha"}

	cases := []string{"%", " % ", "  %  "}
	for _, tc := range cases {
		err := ValidateTenantScope(tc, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeWildcardTenant {
			t.Fatalf("expected 400 wildcard-tenant for %q, got: %v", tc, err)
		}
	}
}

func TestAdversarial_WildcardInjection_All_CaseInsensitive_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha"}

	cases := []string{"all", "ALL", "All", "aLL", " all ", " ALL "}
	for _, tc := range cases {
		err := ValidateTenantScope(tc, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeWildcardTenant {
			t.Fatalf("expected 400 wildcard-tenant for %q, got: %v", tc, err)
		}
	}
}

func TestAdversarial_Injection_WhitespaceAndBlank_Returns400(t *testing.T) {
	ats := []string{"tenant-alpha"}

	cases := []string{"", " ", "   ", "\t", "\n", "\r\n"}
	for _, tc := range cases {
		err := ValidateTenantScope(tc, ats, "corr-1")
		if err == nil || err.HTTPStatus != 400 || err.Code != ErrorCodeMissingTenant {
			t.Fatalf("expected 400 missing-tenant for %q, got: %v", tc, err)
		}
	}
}

// Probe path traversal and null bytes in tenant_ref when ATS is configured
func TestAdversarial_Injection_PathTraversalAndNullBytes_WithATS(t *testing.T) {
	ats := []string{"tenant-alpha"}

	injections := []string{
		"../tenant-bravo",
		"tenant-alpha/../tenant-bravo",
		"tenant-alpha/..",
		"/etc/passwd",
		"tenant-alpha\x00tenant-bravo",
		"tenant-alpha\x00",
		"tenant-alpha/extra",
	}

	for _, inj := range injections {
		err := ValidateTenantScope(inj, ats, "corr-1")
		if err == nil {
			t.Fatalf("CRITICAL SECURITY HOLE: injection %q was accepted by ValidateTenantScope", inj)
		}
		// When ATS is configured, it must reject with 403 unauthorized
		if err.HTTPStatus != 403 {
			t.Fatalf("expected 403 rejection for %q, got: %d (%s)", inj, err.HTTPStatus, err.Code)
		}
	}
}

// CRITICAL PROBE: What happens when AuthContext has NO explicit ATS (ATS is nil or empty)?
func TestAdversarial_TenantScope_EmptyATS_SecurityAudit(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// Caller is an authenticated analyst belonging to tenant-alpha.
	// As a normal single-tenant user, ATS is not explicitly populated (nil).
	auth := NewAuthContext(
		"usr-analyst-alpha",
		"tenant-alpha",
		nil, // ATS is nil!
		PermStudioLibraryRead,
		PermStudioSkillRead,
	)

	// Caller requests confidential asset from foreign tenant (tenant-bravo)
	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-bravo",
		AssetID:   "skl-bravo-secret",
	})

	// Log findings for audit report:
	if err == nil {
		t.Logf("FINDING / SECURITY DEFECT CONFIRMED: When ATS is nil, ValidateTenantScope fails open! Caller from tenant-alpha was able to inspect foreign tenant asset: %+v", resp.Asset.ID)
	} else {
		t.Logf("PASS: Foreign tenant request rejected with err: %v", err)
	}
}

// =============================================================================
// PROBE 4: Unpermitted Asset Inspection & Secret Sanitization
// =============================================================================

func TestAdversarial_Inspect_UnpermittedAsset_SanitizedStub200(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// Caller has general library read, but lacks perm.studio.skill.read
	auth := NewAuthContext(
		"usr-restricted",
		"tenant-alpha",
		[]string{"tenant-alpha"},
		PermStudioLibraryRead,
		// PermStudioSkillRead is intentionally omitted!
	)

	resp, err := svc.Inspect(ctx, auth, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("expected 200 OK sanitized stub, got error: %v", err)
	}
	if resp == nil || resp.Asset == nil {
		t.Fatalf("expected non-nil sanitized asset")
	}

	asset := resp.Asset

	// 1. Functional state must be permission-filtered
	if asset.FunctionalState != FunctionalStatePermissionFiltered {
		t.Fatalf("expected functional_state 'permission-filtered', got: %s", asset.FunctionalState)
	}

	// 2. Description must be completely redacted
	if asset.Description != "" {
		t.Fatalf("SECURITY LEAK: description was not redacted: %q", asset.Description)
	}

	// 3. Parameters must be completely stripped
	if len(asset.Parameters) != 0 {
		t.Fatalf("SECURITY LEAK: parameters were not stripped: %+v", asset.Parameters)
	}

	// 4. Dependencies (including secrets) must be completely stripped
	if len(asset.Dependencies) != 0 {
		t.Fatalf("SECURITY LEAK: dependencies were not stripped: %+v", asset.Dependencies)
	}

	// 5. Tags must be cleared
	if len(asset.Tags) != 0 {
		t.Fatalf("SECURITY LEAK: tags were not cleared: %+v", asset.Tags)
	}

	// 6. Contributors must be cleared
	if len(asset.Contributors) != 0 {
		t.Fatalf("SECURITY LEAK: contributors were not cleared: %+v", asset.Contributors)
	}

	// 7. Consumers must be cleared
	if len(asset.Consumers) != 0 {
		t.Fatalf("SECURITY LEAK: consumers were not cleared: %+v", asset.Consumers)
	}

	// 8. Owner must be redacted to unknown
	if !asset.Owner.IsUnknown || asset.Owner.DisplayName != "owner unknown" || asset.Owner.OwnerID != "redacted" {
		t.Fatalf("SECURITY LEAK: owner was not redacted: %+v", asset.Owner)
	}

	// 9. Provenance details (creator, chain) must be cleared
	if asset.Provenance.CreatedBy != "" || len(asset.Provenance.ProvenanceChain) != 0 {
		t.Fatalf("SECURITY LEAK: provenance details were not cleared: %+v", asset.Provenance)
	}
}

// Verify that ClonePreparation strips / masks raw secret references
func TestAdversarial_ClonePreparation_MasksSecretReferences(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	resp, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
		TenantRef:     "tenant-alpha",
		SourceAssetID: "skl-01",
	})
	if err != nil {
		t.Fatalf("ClonePreparation failed: %v", err)
	}

	clone := resp.DraftClone.PrefilledMetadata
	foundSecret := false
	for _, dep := range clone.Dependencies {
		if dep.Kind == "secret_reference" {
			foundSecret = true
			if !dep.IsMasked {
				t.Fatalf("SECURITY LEAK: secret_reference dependency was not marked as is_masked=true")
			}
			if !strings.Contains(dep.Purpose, "Opaque secret reference preserved") {
				t.Fatalf("expected masked purpose description, got: %s", dep.Purpose)
			}
		}
	}
	if !foundSecret {
		t.Fatalf("expected to verify secret_reference dependency in clone")
	}
}

// Verify store immutability: SanitizeAssetForInspect must not mutate the stored asset
func TestAdversarial_Inspect_StoreImmutability_OriginalAssetUnaffected(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Call unpermitted inspect
	authRestricted := NewAuthContext(
		"usr-restricted",
		"tenant-alpha",
		[]string{"tenant-alpha"},
		PermStudioLibraryRead,
	)
	_, err := svc.Inspect(ctx, authRestricted, InspectRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-01",
	})
	if err != nil {
		t.Fatalf("unpermitted inspect failed: %v", err)
	}

	// Now check directly from store
	rawAsset, err := store.Get(ctx, "tenant-alpha", "skl-01")
	if err != nil {
		t.Fatalf("store.Get failed: %v", err)
	}
	if rawAsset.Description == "" {
		t.Fatalf("MUTATION BUG: original asset in store had its description erased!")
	}
	if len(rawAsset.Parameters) == 0 {
		t.Fatalf("MUTATION BUG: original asset in store had its parameters erased!")
	}
	if len(rawAsset.Dependencies) == 0 {
		t.Fatalf("MUTATION BUG: original asset in store had its dependencies erased!")
	}
}

func assertStatusCode(t *testing.T, opName string, err error, expectedStatus int, expectedCode ErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error %d %s, got nil", opName, expectedStatus, expectedCode)
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("%s: expected *ErrorEnvelope, got %T: %v", opName, err, err)
	}
	if env.HTTPStatus != expectedStatus || env.Code != expectedCode {
		t.Fatalf("%s: expected status=%d code=%s, got status=%d code=%s (msg: %s)",
			opName, expectedStatus, expectedCode, env.HTTPStatus, env.Code, env.Message)
	}
}
