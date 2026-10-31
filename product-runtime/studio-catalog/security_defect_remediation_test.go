package studiocatalog

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// =============================================================================
// Defect A (CRITICAL): ValidateTenantScope Enforcement
// =============================================================================

func TestDefectA_ValidateTenantScope_EmptyActiveTenant_RejectedAtEntry(t *testing.T) {
	// Empty TenantID must be rejected with 400 missing-tenant
	auth := AuthContext{
		PrincipalID: "usr-1",
		TenantID:    "", // Empty!
	}
	err := ValidateTenantScope(auth, "tenant-alpha", "corr-a1")
	if err == nil {
		t.Fatalf("expected error for empty active tenant, got nil")
	}
	if err.HTTPStatus != 400 || err.Code != ErrorCodeMissingTenant {
		t.Fatalf("expected 400 missing-tenant, got HTTP %d, code %s", err.HTTPStatus, err.Code)
	}
}

func TestDefectA_ValidateTenantScope_ActiveTenantNotInATS_RejectedAtEntry(t *testing.T) {
	// Active tenant inconsistent with ATS must be rejected with 403 unauthorized-tenant-in-set
	auth := AuthContext{
		PrincipalID:         "usr-1",
		TenantID:            "tenant-charlie", // Inconsistent with ATS
		AuthorizedTenantSet: []string{"tenant-alpha", "tenant-bravo"},
	}
	err := ValidateTenantScope(auth, "tenant-alpha", "corr-a2")
	if err == nil {
		t.Fatalf("expected error when active tenant is not in ATS, got nil")
	}
	if err.HTTPStatus != 403 || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got HTTP %d, code %s", err.HTTPStatus, err.Code)
	}
}

func TestDefectA_ValidateTenantScope_EmptyATS_CrossTenantRequest_Rejected(t *testing.T) {
	// Single-tenant caller with nil/empty ATS requesting foreign tenant must be rejected
	auth := NewAuthContext("usr-analyst", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)

	err := ValidateTenantScope(auth, "tenant-bravo", "corr-a3")
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: cross-tenant access allowed with nil ATS!")
	}
	if err.HTTPStatus != 403 || err.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got HTTP %d, code %s", err.HTTPStatus, err.Code)
	}
}

func TestDefectA_ValidateTenantScope_NominalAllowedCases(t *testing.T) {
	// Single tenant accessing own scope
	authSingle := NewAuthContext("usr-analyst", "tenant-alpha", nil, PermStudioLibraryRead)
	if err := ValidateTenantScope(authSingle, "tenant-alpha", "corr-a4"); err != nil {
		t.Fatalf("expected single tenant self-access to succeed, got %v", err)
	}

	// MSSP caller accessing tenant in ATS
	authMSSP := NewAuthContext("usr-mssp", "tenant-alpha", []string{"tenant-alpha", "tenant-bravo"}, PermStudioLibraryRead)
	if err := ValidateTenantScope(authMSSP, "tenant-bravo", "corr-a5"); err != nil {
		t.Fatalf("expected MSSP access to tenant in ATS to succeed, got %v", err)
	}
}

// =============================================================================
// Defect B (CRITICAL): CanManageAssetKind Read Permission Escalation Prevention
// =============================================================================

func TestDefectB_CanManageAssetKind_ReadOnlyPermissions_ReturnFalseForAllKinds(t *testing.T) {
	// LibraryRead alone
	authLibRead := NewAuthContext("usr-reader", "tenant-alpha", nil, PermStudioLibraryRead)
	for _, kind := range AllAssetKinds {
		if authLibRead.CanManageAssetKind(kind) {
			t.Fatalf("CRITICAL SECURITY DEFECT: CanManageAssetKind(%s) returned true with PermStudioLibraryRead alone!", kind)
		}
	}

	// BuilderRead alone
	authBuilderRead := NewAuthContext("usr-builder-reader", "tenant-alpha", nil, PermStudioBuilderRead)
	for _, kind := range AllAssetKinds {
		if authBuilderRead.CanManageAssetKind(kind) {
			t.Fatalf("CRITICAL SECURITY DEFECT: CanManageAssetKind(%s) returned true with PermStudioBuilderRead alone!", kind)
		}
	}

	// Both LibraryRead and BuilderRead combined
	authBothRead := NewAuthContext("usr-reader-both", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioBuilderRead)
	for _, kind := range AllAssetKinds {
		if authBothRead.CanManageAssetKind(kind) {
			t.Fatalf("CRITICAL SECURITY DEFECT: CanManageAssetKind(%s) returned true with read permissions alone!", kind)
		}
	}
}

func TestDefectB_CanManageAssetKind_ToolRequiresPermStudioArtifactManage(t *testing.T) {
	authManage := NewAuthContext("usr-manager", "tenant-alpha", nil, PermStudioArtifactManage)
	if !authManage.CanManageAssetKind(AssetKindTool) {
		t.Fatalf("expected CanManageAssetKind(AssetKindTool) to return true with PermStudioArtifactManage")
	}

	// Verify alias perm.cmdr-studio.artifact.manage also works
	authAlias := NewAuthContext("usr-manager-alias", "tenant-alpha", nil, "perm.cmdr-studio.artifact.manage")
	if !authAlias.CanManageAssetKind(AssetKindTool) {
		t.Fatalf("expected CanManageAssetKind(AssetKindTool) to return true with alias perm.cmdr-studio.artifact.manage")
	}
}

func TestDefectB_AssignOwner_ToolAsset_ReadOnlyUser_Forbidden(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Seed a tool asset
	toolAsset := createTestAsset("tenant-alpha", "tool-sec-01", "sec-tool", AssetKindTool, LifecycleStateDraft)
	_ = store.Put(ctx, toolAsset)

	// User with only PermStudioLibraryRead and PermStudioBuilderRead attempts AssignOwner
	readerAuth := NewAuthContext("usr-reader", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioBuilderRead)
	_, err := svc.AssignOwner(ctx, readerAuth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "tool-sec-01",
		NewOwner:  Owner{OwnerID: "usr-hacker", OwnerName: "Hacker", OwnerType: "user"},
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: AssignOwner on Tool asset succeeded with read permissions alone!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}

	// User with PermStudioArtifactManage succeeds
	managerAuth := NewAuthContext("usr-manager", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioArtifactManage)
	resp, err := svc.AssignOwner(ctx, managerAuth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "tool-sec-01",
		NewOwner:  Owner{OwnerID: "usr-legit", OwnerName: "Legit Owner", OwnerType: "user"},
	})
	if err != nil {
		t.Fatalf("expected AssignOwner with PermStudioArtifactManage to succeed, got %v", err)
	}
	if resp.Asset.Owner.OwnerID != "usr-legit" {
		t.Fatalf("expected owner usr-legit, got %s", resp.Asset.Owner.OwnerID)
	}
}

// =============================================================================
// Defect C (CRITICAL): RecordSupersession Permission & Atomic Rollback
// =============================================================================

func TestDefectC_RecordSupersession_ReadOnlyLibraryPermission_Forbidden(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	older := createTestAsset("tenant-alpha", "skl-v1", "skill-v1", AssetKindSkill, LifecycleStateDraft)
	older.Version.SemVer = "1.0.0"
	newer := createTestAsset("tenant-alpha", "skl-v2", "skill-v2", AssetKindSkill, LifecycleStateDraft)
	newer.Version.SemVer = "1.1.0"
	_ = store.Put(ctx, older)
	_ = store.Put(ctx, newer)

	// Caller with PermStudioLibraryRead alone MUST be rejected
	readerAuth := NewAuthContext("usr-reader", "tenant-alpha", nil, PermStudioLibraryRead)
	_, err := svc.RecordSupersession(ctx, readerAuth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-v1",
		NewerAssetID: "skl-v2",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: RecordSupersession succeeded with PermStudioLibraryRead alone!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}
}

// Store wrapper to simulate second Put failure for atomic rollback verification
type failingSecondPutStore struct {
	Store
	putCount int
	failOn   int
}

func (f *failingSecondPutStore) Put(ctx context.Context, asset *Asset) error {
	f.putCount++
	if f.putCount == f.failOn {
		return errors.New("simulated storage failure on second write")
	}
	return f.Store.Put(ctx, asset)
}

func TestDefectC_RecordSupersession_SecondPutFailure_RollsBackOlderAsset(t *testing.T) {
	underlyingStore := NewInMemoryStore()
	wrappedStore := &failingSecondPutStore{
		Store:  underlyingStore,
		failOn: 2, // Fails on the 2nd Put call (when writing newerAsset)
	}
	svc := NewCatalogService(wrappedStore)
	ctx := context.Background()

	older := createTestAsset("tenant-alpha", "skl-older", "older-skill", AssetKindSkill, LifecycleStateDraft)
	older.Version.SemVer = "1.0.0"
	older.LifecycleState = LifecycleStateDraft
	older.Supersession.IsSuperseded = false
	older.Supersession.SupersededBy = ""

	newer := createTestAsset("tenant-alpha", "skl-newer", "newer-skill", AssetKindSkill, LifecycleStateDraft)
	newer.Version.SemVer = "2.0.0"

	_ = underlyingStore.Put(ctx, older)
	_ = underlyingStore.Put(ctx, newer)

	auth := adminAuth("tenant-alpha")

	// Execute RecordSupersession - 1st Put (older) succeeds, 2nd Put (newer) fails, rollback restores older
	_, err := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-older",
		NewerAssetID: "skl-newer",
	})
	if err == nil {
		t.Fatalf("expected error from simulated second Put failure, got nil")
	}

	// Verify older asset in store has been rolled back to its original state!
	olderAfter, getErr := underlyingStore.Get(ctx, "tenant-alpha", "skl-older")
	if getErr != nil {
		t.Fatalf("failed to retrieve older asset after rollback: %v", getErr)
	}
	if olderAfter.Supersession.IsSuperseded {
		t.Fatalf("ROLLBACK FAILURE: older asset remained superseded despite second Put failure!")
	}
	if olderAfter.Supersession.SupersededBy != "" {
		t.Fatalf("ROLLBACK FAILURE: older asset still has SupersededBy = %q", olderAfter.Supersession.SupersededBy)
	}
	if olderAfter.LifecycleState != LifecycleStateDraft {
		t.Fatalf("ROLLBACK FAILURE: older asset lifecycle state was not restored, got %s", olderAfter.LifecycleState)
	}
}

// =============================================================================
// Defect D (HIGH): Asset-Specific Reauthorization & Supersession Lineage Filtering
// =============================================================================

func TestDefectD_ValidateMetadata_UnauthorizedAssetKind_Returns403(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Seed a workflow asset
	wf := createTestAsset("tenant-alpha", "wf-secret", "confidential-workflow", AssetKindWorkflow, LifecycleStateDraft)
	_ = store.Put(ctx, wf)

	// Caller has library.read and skill.read, but lacks workflow.read
	auth := NewAuthContext("usr-analyst", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)

	_, err := svc.ValidateMetadata(ctx, auth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-secret",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: ValidateMetadata leaked metadata of unauthorized workflow!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}
}

func TestDefectD_SuggestClassificationAndTags_UnauthorizedAssetKind_Returns403(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Seed an automation-agent asset
	agent := createTestAsset("tenant-alpha", "agent-secret", "secret-agent", AssetKindAutomationAgent, LifecycleStateDraft)
	_ = store.Put(ctx, agent)

	// Caller lacks automation-agent.read
	auth := NewAuthContext("usr-analyst", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)

	_, err := svc.SuggestClassificationAndTags(ctx, auth, SuggestTagsRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "agent-secret",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: SuggestClassificationAndTags allowed unauthorized agent inspection!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}
}

func TestDefectD_GetSupersessionHistory_FiltersOutUnauthorizedAssetsInChain(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Lineage: Skill v1 (skl-1) -> superseded by Workflow v2 (wf-2) -> superseded by Workflow v3 (wf-3)
	skl1 := createTestAsset("tenant-alpha", "skl-1", "skill-v1", AssetKindSkill, LifecycleStateSuperseded)
	skl1.Version.SemVer = "1.0.0"
	skl1.Supersession.IsSuperseded = true
	skl1.Supersession.SupersededBy = "wf-2"

	wf2 := createTestAsset("tenant-alpha", "wf-2", "workflow-v2", AssetKindWorkflow, LifecycleStateSuperseded)
	wf2.Version.SemVer = "2.0.0"
	wf2.Supersession.IsSuperseded = true
	wf2.Supersession.Supersedes = "skl-1"
	wf2.Supersession.SupersededBy = "wf-3"

	wf3 := createTestAsset("tenant-alpha", "wf-3", "workflow-v3", AssetKindWorkflow, LifecycleStatePublished)
	wf3.Version.SemVer = "3.0.0"
	wf3.Supersession.IsSuperseded = false
	wf3.Supersession.Supersedes = "wf-2"

	_ = store.Put(ctx, skl1)
	_ = store.Put(ctx, wf2)
	_ = store.Put(ctx, wf3)

	// Caller has permission to read Workflows, but NOT Skills
	wfOnlyAuth := NewAuthContext(
		"usr-wf-engineer",
		"tenant-alpha",
		nil,
		PermStudioLibraryRead,
		PermStudioWorkflowRead, // PermStudioSkillRead is OMITTED!
	)

	// Caller requests history of wf-2
	resp, err := svc.GetSupersessionHistory(ctx, wfOnlyAuth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-2",
	})
	if err != nil {
		t.Fatalf("GetSupersessionHistory failed: %v", err)
	}

	// Verify that skl-1 (Skill) is NOT included in the Chain
	for _, node := range resp.Chain {
		if node.AssetID == "skl-1" {
			t.Fatalf("CRITICAL SECURITY DEFECT: Unauthorized ancestor skl-1 was leaked in supersession chain!")
		}
	}

	// Verify that authorized workflows (wf-2, wf-3) ARE present
	foundWF2 := false
	foundWF3 := false
	for _, node := range resp.Chain {
		if node.AssetID == "wf-2" {
			foundWF2 = true
		}
		if node.AssetID == "wf-3" {
			foundWF3 = true
		}
	}
	if !foundWF2 || !foundWF3 {
		t.Fatalf("expected wf-2 and wf-3 in chain, got chain: %+v", resp.Chain)
	}
}

// =============================================================================
// Defect E (HIGH): Audit Trail & Provenance Spoofing Prevention
// =============================================================================

func TestDefectE_AssignOwner_ClientSuppliedActor_IgnoredInAuditTrail(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	asset := createTestAsset("tenant-alpha", "skl-audit-01", "audit-skill", AssetKindSkill, LifecycleStateDraft)
	_ = store.Put(ctx, asset)

	auth := adminAuth("tenant-alpha") // PrincipalID = "usr-admin"

	resp, err := svc.AssignOwner(ctx, auth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-audit-01",
		NewOwner:  Owner{OwnerID: "usr-new", OwnerName: "New Owner", OwnerType: "user"},
		Actor:     "spoofed-ciso-actor", // Spoofed actor in request body!
		Reason:    "Testing audit actor authenticity",
	})
	if err != nil {
		t.Fatalf("AssignOwner failed: %v", err)
	}

	lastRecord := resp.Asset.AuditTrail[len(resp.Asset.AuditTrail)-1]
	if lastRecord.Actor != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: Audit record accepted spoofed actor %q instead of authenticated principal %q", lastRecord.Actor, "usr-admin")
	}
}

func TestDefectE_EditDraftMetadata_ClientSuppliedActor_IgnoredInAuditTrail(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	asset := createTestAsset("tenant-alpha", "skl-audit-02", "draft-skill", AssetKindSkill, LifecycleStateDraft)
	_ = store.Put(ctx, asset)

	auth := adminAuth("tenant-alpha") // PrincipalID = "usr-admin"

	resp, err := svc.UpdateDraftMetadata(ctx, auth, UpdateDraftMetadataRequest{
		TenantRef:   "tenant-alpha",
		AssetID:     "skl-audit-02",
		DisplayName: "Updated Display Name",
		Actor:       "spoofed-superadmin", // Spoofed actor!
	})
	if err != nil {
		t.Fatalf("UpdateDraftMetadata failed: %v", err)
	}

	lastRecord := resp.Asset.AuditTrail[len(resp.Asset.AuditTrail)-1]
	if lastRecord.Actor != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: UpdateDraftMetadata accepted spoofed actor %q instead of %q", lastRecord.Actor, "usr-admin")
	}
}

func TestDefectE_RecordSupersession_ClientSuppliedActor_IgnoredInAuditTrail(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	older := createTestAsset("tenant-alpha", "skl-old-audit", "old", AssetKindSkill, LifecycleStateDraft)
	older.Version.SemVer = "1.0.0"
	newer := createTestAsset("tenant-alpha", "skl-new-audit", "new", AssetKindSkill, LifecycleStateDraft)
	newer.Version.SemVer = "2.0.0"
	_ = store.Put(ctx, older)
	_ = store.Put(ctx, newer)

	auth := adminAuth("tenant-alpha") // PrincipalID = "usr-admin"

	_, err := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-old-audit",
		NewerAssetID: "skl-new-audit",
		Actor:        "spoofed-ciso", // Spoofed actor!
	})
	if err != nil {
		t.Fatalf("RecordSupersession failed: %v", err)
	}

	updatedOlder, _ := store.Get(ctx, "tenant-alpha", "skl-old-audit")
	lastRecord := updatedOlder.AuditTrail[len(updatedOlder.AuditTrail)-1]
	if lastRecord.Actor != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: RecordSupersession accepted spoofed actor %q instead of %q", lastRecord.Actor, "usr-admin")
	}
}

func TestDefectE_ClonePreparation_ClientSuppliedClonedBy_IgnoredInProvenance(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	auth := adminAuth("tenant-alpha") // PrincipalID = "usr-admin"

	resp, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-alpha",
		SourceAssetID:   "skl-01",
		ClonedBy:        "spoofed-lead-engineer", // Spoofed!
	})
	if err != nil {
		t.Fatalf("ClonePreparation failed: %v", err)
	}

	if resp.DraftClone.ClonedBy != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: DraftClone.ClonedBy accepted spoofed %q instead of %q", resp.DraftClone.ClonedBy, "usr-admin")
	}
	if resp.DraftClone.PrefilledMetadata.Provenance.CreatedBy != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: Provenance.CreatedBy accepted spoofed %q instead of %q", resp.DraftClone.PrefilledMetadata.Provenance.CreatedBy, "usr-admin")
	}
	if resp.DraftClone.PrefilledMetadata.Provenance.UpdatedBy != "usr-admin" {
		t.Fatalf("CRITICAL SECURITY DEFECT: Provenance.UpdatedBy accepted spoofed %q instead of %q", resp.DraftClone.PrefilledMetadata.Provenance.UpdatedBy, "usr-admin")
	}
}

// =============================================================================
// Defect F (MEDIUM): Cryptographic HMAC Signed Cursor Integrity
// =============================================================================

func TestDefectF_CursorHMAC_ValidEncodingAndDecoding(t *testing.T) {
	p := CursorPayload{
		Version:       1,
		TenantID:      "tenant-alpha",
		SortField:     "updated_at",
		SortDirection: "desc",
		LastValue:     "2026-10-31T00:00:00Z",
		TiebreakerID:  "skl-01",
		Fingerprint:   "fp-test-123",
	}

	encoded := EncodeCursor(p)
	if !strings.Contains(encoded, ".") {
		t.Fatalf("expected signed cursor with dot separator, got %q", encoded)
	}

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("failed to decode valid signed cursor: %v", err)
	}
	if decoded.TenantID != p.TenantID || decoded.TiebreakerID != p.TiebreakerID || decoded.SortField != p.SortField {
		t.Fatalf("decoded cursor does not match original payload: %+v vs %+v", decoded, p)
	}
}

func TestDefectF_CursorHMAC_TamperingAndBitFlip_Rejected(t *testing.T) {
	p := CursorPayload{
		Version:       1,
		TenantID:      "tenant-alpha",
		SortField:     "updated_at",
		SortDirection: "desc",
		LastValue:     "2026-10-31T00:00:00Z",
		TiebreakerID:  "skl-01",
		Fingerprint:   "fp-test-123",
	}

	validToken := EncodeCursor(p)
	parts := strings.Split(validToken, ".")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	payloadB64, sigB64 := parts[0], parts[1]

	// 1. Bit-flip in signature
	tamperedSigBytes := []byte(sigB64)
	if tamperedSigBytes[0] == 'A' {
		tamperedSigBytes[0] = 'B'
	} else {
		tamperedSigBytes[0] = 'A'
	}
	tamperedSigToken := payloadB64 + "." + string(tamperedSigBytes)
	_, errSig := DecodeCursor(tamperedSigToken)
	if errSig == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: Cursor with flipped signature bit was accepted!")
	}

	// 2. Bit-flip in payload
	tamperedPayloadBytes := []byte(payloadB64)
	if tamperedPayloadBytes[0] == 'A' {
		tamperedPayloadBytes[0] = 'B'
	} else {
		tamperedPayloadBytes[0] = 'A'
	}
	tamperedPayloadToken := string(tamperedPayloadBytes) + "." + sigB64
	_, errPayload := DecodeCursor(tamperedPayloadToken)
	if errPayload == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: Cursor with flipped payload bit was accepted!")
	}

	// 3. Modifying JSON field inside payload without re-signing
	b, _ := base64.RawURLEncoding.DecodeString(payloadB64)
	var tamperedPayload CursorPayload
	_ = json.Unmarshal(b, &tamperedPayload)
	tamperedPayload.TenantID = "tenant-bravo" // Cross-tenant forgery attempt
	newB, _ := json.Marshal(tamperedPayload)
	forgedPayloadB64 := base64.RawURLEncoding.EncodeToString(newB)
	forgedToken := forgedPayloadB64 + "." + sigB64
	_, errForged := DecodeCursor(forgedToken)
	if errForged == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: Forged cursor with modified tenant was accepted with old signature!")
	}

	// 4. Missing signature (unsigned legacy cursor)
	_, errUnsigned := DecodeCursor(payloadB64)
	if errUnsigned == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: Unsigned cursor was accepted!")
	}
}
