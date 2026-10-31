package studiocatalog

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

// =============================================================================
// Challenger Empirical Verification: Defect A (Tenant Isolation)
// =============================================================================

func TestChallenger_DefectA_TenantIsolationEnforcement(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()

	// Case 1: Empty TenantID in AuthContext must be rejected with 400
	authEmptyTenant := AuthContext{
		PrincipalID: "usr-challenger",
		TenantID:    "   ", // whitespace only
	}
	errEnv := ValidateTenantScope(authEmptyTenant, "tenant-alpha", "corr-ch-1")
	if errEnv == nil || errEnv.HTTPStatus != 400 || errEnv.Code != ErrorCodeMissingTenant {
		t.Fatalf("expected 400 missing-tenant for empty active tenant, got: %v", errEnv)
	}

	// Case 2: Active tenant not in caller's ATS must be rejected with 403
	authInconsistentATS := AuthContext{
		PrincipalID:         "usr-challenger",
		TenantID:            "tenant-charlie",
		AuthorizedTenantSet: []string{"tenant-alpha", "tenant-bravo"},
	}
	errEnv = ValidateTenantScope(authInconsistentATS, "tenant-alpha", "corr-ch-2")
	if errEnv == nil || errEnv.HTTPStatus != 403 || errEnv.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set for active tenant not in ATS, got: %v", errEnv)
	}

	// Case 3: Single-tenant caller with nil ATS attempting cross-tenant access must be rejected with 403
	authNilATS := NewAuthContext("usr-analyst", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)
	errEnv = ValidateTenantScope(authNilATS, "tenant-bravo", "corr-ch-3")
	if errEnv == nil || errEnv.HTTPStatus != 403 || errEnv.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set for cross-tenant request with nil ATS, got: %v", errEnv)
	}

	// Case 4: Single-tenant caller with empty slice ATS attempting cross-tenant access must be rejected with 403
	authEmptySliceATS := NewAuthContext("usr-analyst", "tenant-alpha", []string{}, PermStudioLibraryRead, PermStudioSkillRead)
	errEnv = ValidateTenantScope(authEmptySliceATS, "tenant-bravo", "corr-ch-4")
	if errEnv == nil || errEnv.HTTPStatus != 403 || errEnv.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set for cross-tenant request with empty slice ATS, got: %v", errEnv)
	}

	// Case 5: Verify Inspect rejects cross-tenant access end-to-end
	_, err := svc.Inspect(ctx, authNilATS, InspectRequest{
		TenantRef: "tenant-bravo",
		AssetID:   "skl-bravo-secret",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY VULNERABILITY: Inspect allowed cross-tenant access with nil ATS!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Fatalf("expected 403 unauthorized-tenant-in-set on Inspect, got: %v", err)
	}
}

// =============================================================================
// Challenger Empirical Verification: Defect B (Tool RBAC Mutation Prevention)
// =============================================================================

func TestChallenger_DefectB_ToolRBACMutations_ReadOnlyRejected(t *testing.T) {
	catStore, catalogSvc := setupTestStore(t)
	toolSvc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Seed existing tool asset in catalog store for AssignOwner test
	toolAsset := createTestAsset("tenant-alpha", "tool-target-01", "target-tool", AssetKindTool, LifecycleStateDraft)
	_ = catStore.Put(ctx, toolAsset)

	// User identities with read-only permissions
	readAuthLibrary := NewAuthContext("usr-reader-lib", "tenant-alpha", nil, PermStudioLibraryRead)
	readAuthBuilder := NewAuthContext("usr-reader-bld", "tenant-alpha", nil, PermStudioBuilderRead)
	readAuthCombined := NewAuthContext("usr-reader-both", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioBuilderRead)

	readContexts := []*AuthContext{&readAuthLibrary, &readAuthBuilder, &readAuthCombined}

	for i, ac := range readContexts {
		// 1. CanManageAssetKind(AssetKindTool) must be false
		if ac.CanManageAssetKind(AssetKindTool) {
			t.Fatalf("context %d: CanManageAssetKind(AssetKindTool) returned true for read permissions!", i)
		}

		// 2. AssignOwner via CatalogService must return 403
		_, errAssign := catalogSvc.AssignOwner(ctx, *ac, AssignOwnerRequest{
			TenantRef: "tenant-alpha",
			AssetID:   "tool-target-01",
			NewOwner:  Owner{OwnerID: "usr-hacker", OwnerName: "Hacker", OwnerType: "user"},
		})
		if errAssign == nil {
			t.Fatalf("context %d: AssignOwner succeeded with read-only permissions!", i)
		}
		envAssign, ok := errAssign.(*ErrorEnvelope)
		if !ok || envAssign.HTTPStatus != 403 || envAssign.Code != ErrorCodePermissionDenied {
			t.Fatalf("context %d: expected 403 permission-denied on AssignOwner, got %v", i, errAssign)
		}

		// 3. CreateDraftTool must return 403
		_, errCreate := toolSvc.CreateDraftTool(ctx, ac, &CreateDraftToolRequest{
			TenantRef:       "tenant-alpha",
			ID:              "tool-unauthorized-new",
			Name:            "Unauthorized Tool",
			Purpose:         "Attempting creation without artifact manage permission",
			ExecutionNature: ExecutionNatureDeterministic,
		})
		if errCreate == nil {
			t.Fatalf("context %d: CreateDraftTool succeeded with read-only permissions!", i)
		}
		envCreate, ok := errCreate.(*ErrorEnvelope)
		if !ok || envCreate.HTTPStatus != 403 || envCreate.Code != ErrorCodePermissionDenied {
			t.Fatalf("context %d: expected 403 permission-denied on CreateDraftTool, got %v", i, errCreate)
		}

		// 4. UpdateDraftTool must return 403
		_, errUpdate := toolSvc.UpdateDraftTool(ctx, ac, &UpdateDraftToolRequest{
			TenantRef: "tenant-alpha",
			ToolID:    "tool-text-formatter",
			Name:      "Hacked Name",
			Purpose:   "Attempting update without artifact manage permission",
		})
		if errUpdate == nil {
			t.Fatalf("context %d: UpdateDraftTool succeeded with read-only permissions!", i)
		}
		envUpdate, ok := errUpdate.(*ErrorEnvelope)
		if !ok || envUpdate.HTTPStatus != 403 || envUpdate.Code != ErrorCodePermissionDenied {
			t.Fatalf("context %d: expected 403 permission-denied on UpdateDraftTool, got %v", i, errUpdate)
		}
	}
}

// =============================================================================
// Challenger Empirical Verification: Defect C (Supersession RBAC & Rollback)
// =============================================================================

func TestChallenger_DefectC_RecordSupersession_PermissionsAndRollback(t *testing.T) {
	ctx := context.Background()

	// 1. Permission checks
	store := NewInMemoryStore()
	svc := NewCatalogService(store)

	sklOlder := createTestAsset("tenant-alpha", "skl-old", "skill-old", AssetKindSkill, LifecycleStateDraft)
	sklOlder.Version.SemVer = "1.0.0"
	sklNewer := createTestAsset("tenant-alpha", "skl-new", "skill-new", AssetKindSkill, LifecycleStateDraft)
	sklNewer.Version.SemVer = "1.1.0"
	_ = store.Put(ctx, sklOlder)
	_ = store.Put(ctx, sklNewer)

	// Caller with ONLY PermStudioLibraryRead MUST receive 403
	authLibRead := NewAuthContext("usr-read", "tenant-alpha", nil, PermStudioLibraryRead)
	_, err := svc.RecordSupersession(ctx, authLibRead, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-old",
		NewerAssetID: "skl-new",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: RecordSupersession succeeded with PermStudioLibraryRead alone!")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied, got %v", err)
	}

	// Caller with ONLY PermStudioBuilderRead MUST receive 403
	authBuilderRead := NewAuthContext("usr-builder", "tenant-alpha", nil, PermStudioBuilderRead)
	_, err = svc.RecordSupersession(ctx, authBuilderRead, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-old",
		NewerAssetID: "skl-new",
	})
	if err == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: RecordSupersession succeeded with PermStudioBuilderRead alone!")
	}

	// Caller has manage for Skill, but older is Skill and newer is Workflow -> must fail without workflow manage
	wfNewer := createTestAsset("tenant-alpha", "wf-new", "workflow-new", AssetKindWorkflow, LifecycleStateDraft)
	_ = store.Put(ctx, wfNewer)
	authSkillManageOnly := NewAuthContext("usr-skill-mgr", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillManage)
	_, err = svc.RecordSupersession(ctx, authSkillManageOnly, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-old",
		NewerAssetID: "wf-new",
	})
	if err == nil {
		t.Fatalf("expected error when caller lacks manage permission for newer asset kind, got nil")
	}

	// 2. Atomic Rollback Verification
	failingStore := &failingSecondPutStore{
		Store:  NewInMemoryStore(),
		failOn: 2,
	}
	failingSvc := NewCatalogService(failingStore)
	_ = failingStore.Put(ctx, sklOlder)
	_ = failingStore.Put(ctx, sklNewer)

	authAdmin := adminAuth("tenant-alpha")
	_, err = failingSvc.RecordSupersession(ctx, authAdmin, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-old",
		NewerAssetID: "skl-new",
	})
	if err == nil {
		t.Fatalf("expected error on simulated second put failure, got nil")
	}

	// Check that sklOlder was restored to its un-superseded state
	restored, getErr := failingStore.Get(ctx, "tenant-alpha", "skl-old")
	if getErr != nil || restored == nil {
		t.Fatalf("failed to retrieve older asset: %v", getErr)
	}
	if restored.Supersession.IsSuperseded || restored.Supersession.SupersededBy != "" {
		t.Fatalf("ROLLBACK FAILURE: older asset remained marked superseded after failure!")
	}
}

// =============================================================================
// Challenger Empirical Verification: Defect D (Metadata & Lineage Authorization)
// =============================================================================

func TestChallenger_DefectD_MetadataAndLineageAuthorization(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	// Seed assets of diverse types
	agent := createTestAsset("tenant-alpha", "agent-top-secret", "Confidential Agent", AssetKindAutomationAgent, LifecycleStateDraft)
	_ = store.Put(ctx, agent)

	// Caller has library read and skill read, but lacks automation-agent.read
	authRestricted := NewAuthContext("usr-restricted", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)

	// 1. ValidateMetadata must return 403
	_, errVM := svc.ValidateMetadata(ctx, authRestricted, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "agent-top-secret",
	})
	if errVM == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: ValidateMetadata succeeded for unauthorized asset kind!")
	}
	envVM, ok := errVM.(*ErrorEnvelope)
	if !ok || envVM.HTTPStatus != 403 || envVM.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied on ValidateMetadata, got %v", errVM)
	}

	// 2. SuggestClassificationAndTags must return 403
	_, errSuggest := svc.SuggestClassificationAndTags(ctx, authRestricted, SuggestTagsRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "agent-top-secret",
	})
	if errSuggest == nil {
		t.Fatalf("CRITICAL SECURITY DEFECT: SuggestClassificationAndTags succeeded for unauthorized asset kind!")
	}
	envSuggest, ok := errSuggest.(*ErrorEnvelope)
	if !ok || envSuggest.HTTPStatus != 403 || envSuggest.Code != ErrorCodePermissionDenied {
		t.Fatalf("expected 403 permission-denied on SuggestClassificationAndTags, got %v", errSuggest)
	}

	// 3. GetSupersessionHistory lineage filtering
	// Build a 4-asset chain: S1 (Skill) -> A2 (Agent) -> W3 (Workflow) -> W4 (Workflow)
	s1 := createTestAsset("tenant-alpha", "skl-ch-1", "skill-1", AssetKindSkill, LifecycleStateSuperseded)
	s1.Supersession.IsSuperseded = true
	s1.Supersession.SupersededBy = "agent-ch-2"

	a2 := createTestAsset("tenant-alpha", "agent-ch-2", "agent-2", AssetKindAutomationAgent, LifecycleStateSuperseded)
	a2.Supersession.IsSuperseded = true
	a2.Supersession.Supersedes = "skl-ch-1"
	a2.Supersession.SupersededBy = "wf-ch-3"

	w3 := createTestAsset("tenant-alpha", "wf-ch-3", "workflow-3", AssetKindWorkflow, LifecycleStateSuperseded)
	w3.Supersession.IsSuperseded = true
	w3.Supersession.Supersedes = "agent-ch-2"
	w3.Supersession.SupersededBy = "wf-ch-4"

	w4 := createTestAsset("tenant-alpha", "wf-ch-4", "workflow-4", AssetKindWorkflow, LifecycleStatePublished)
	w4.Supersession.IsSuperseded = false
	w4.Supersession.Supersedes = "wf-ch-3"

	_ = store.Put(ctx, s1)
	_ = store.Put(ctx, a2)
	_ = store.Put(ctx, w3)
	_ = store.Put(ctx, w4)

	// Caller only has Workflow read permission (lacks Skill and Agent permissions)
	authWFOnly := NewAuthContext("usr-wf-dev", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioWorkflowRead)

	histResp, errHist := svc.GetSupersessionHistory(ctx, authWFOnly, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-ch-3",
	})
	if errHist != nil {
		t.Fatalf("GetSupersessionHistory failed: %v", errHist)
	}

	// Chain must contain ONLY wf-ch-3 and wf-ch-4. Must NEVER leak skl-ch-1 or agent-ch-2!
	for _, node := range histResp.Chain {
		if node.AssetID == "skl-ch-1" || node.AssetID == "agent-ch-2" {
			t.Fatalf("CRITICAL SECURITY DEFECT: Unauthorized node %s leaked in supersession chain!", node.AssetID)
		}
	}
	if len(histResp.Chain) != 2 {
		t.Fatalf("expected exactly 2 nodes in filtered chain, got %d", len(histResp.Chain))
	}
	if histResp.RootAssetID != "wf-ch-3" || histResp.LatestAssetID != "wf-ch-4" {
		t.Fatalf("RootAssetID or LatestAssetID leaked unauthorized boundaries: root=%s latest=%s", histResp.RootAssetID, histResp.LatestAssetID)
	}
}

// =============================================================================
// Challenger Empirical Verification: Defect E (Audit Trail & Provenance Spoofing)
// =============================================================================

func TestChallenger_DefectE_AuditTrail_ProvenanceSpoofingDefense(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()

	auth := adminAuth("tenant-alpha") // PrincipalID = "usr-admin"

	// 1. AssignOwner spoof attempt
	asset1 := createTestAsset("tenant-alpha", "skl-sp-1", "skill-1", AssetKindSkill, LifecycleStateDraft)
	_ = store.Put(ctx, asset1)

	respAssign, err := svc.AssignOwner(ctx, auth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-sp-1",
		NewOwner:  Owner{OwnerID: "usr-target", OwnerName: "Target", OwnerType: "user"},
		Actor:     "malicious-injected-actor",
	})
	if err != nil {
		t.Fatalf("AssignOwner failed: %v", err)
	}
	lastAudit := respAssign.Asset.AuditTrail[len(respAssign.Asset.AuditTrail)-1]
	if lastAudit.Actor != "usr-admin" {
		t.Fatalf("AUDIT SPOOF DEFECT: AssignOwner recorded actor %q instead of authenticated %q", lastAudit.Actor, "usr-admin")
	}

	// 2. UpdateDraftMetadata spoof attempt
	respUpdate, err := svc.UpdateDraftMetadata(ctx, auth, UpdateDraftMetadataRequest{
		TenantRef:   "tenant-alpha",
		AssetID:     "skl-sp-1",
		DisplayName: "New Display Name",
		Actor:       "adversary-root-identity",
	})
	if err != nil {
		t.Fatalf("UpdateDraftMetadata failed: %v", err)
	}
	lastAudit = respUpdate.Asset.AuditTrail[len(respUpdate.Asset.AuditTrail)-1]
	if lastAudit.Actor != "usr-admin" {
		t.Fatalf("AUDIT SPOOF DEFECT: UpdateDraftMetadata recorded actor %q instead of authenticated %q", lastAudit.Actor, "usr-admin")
	}

	// 3. ClonePreparation spoof attempt
	respClone, err := svc.ClonePreparation(ctx, auth, ClonePreparationRequest{
		TenantRef:       "tenant-alpha",
		TargetTenantRef: "tenant-alpha",
		SourceAssetID:   "skl-01",
		ClonedBy:        "spoofed-lead-architect",
	})
	if err != nil {
		t.Fatalf("ClonePreparation failed: %v", err)
	}
	if respClone.DraftClone.ClonedBy != "usr-admin" {
		t.Fatalf("PROVENANCE SPOOF DEFECT: DraftClone.ClonedBy = %q", respClone.DraftClone.ClonedBy)
	}
	if respClone.DraftClone.PrefilledMetadata.Provenance.CreatedBy != "usr-admin" {
		t.Fatalf("PROVENANCE SPOOF DEFECT: CreatedBy = %q", respClone.DraftClone.PrefilledMetadata.Provenance.CreatedBy)
	}
	if respClone.DraftClone.PrefilledMetadata.Provenance.UpdatedBy != "usr-admin" {
		t.Fatalf("PROVENANCE SPOOF DEFECT: UpdatedBy = %q", respClone.DraftClone.PrefilledMetadata.Provenance.UpdatedBy)
	}

	// 4. Empty PrincipalID in AuthContext should default to "unknown", never arbitrary input
	authEmptyPrincipal := AuthContext{
		PrincipalID: "",
		TenantID:    "tenant-alpha",
		Permissions: auth.Permissions,
	}
	respAssign2, err := svc.AssignOwner(ctx, authEmptyPrincipal, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-sp-1",
		NewOwner:  Owner{OwnerID: "usr-target2", OwnerName: "Target2", OwnerType: "user"},
		Actor:     "attacker-controlled-string",
	})
	if err != nil {
		t.Fatalf("AssignOwner with empty principal failed: %v", err)
	}
	lastAudit2 := respAssign2.Asset.AuditTrail[len(respAssign2.Asset.AuditTrail)-1]
	if lastAudit2.Actor != "unknown" {
		t.Fatalf("AUDIT SPOOF DEFECT: Expected fallback 'unknown' for empty principal, got %q", lastAudit2.Actor)
	}
}

// =============================================================================
// Challenger Empirical Verification: Defect F (HMAC Signed Pagination Cursors)
// =============================================================================

func TestChallenger_DefectF_SignedCursor_TamperResistanceFailsClosed(t *testing.T) {
	payload := CursorPayload{
		Version:       1,
		TenantID:      "tenant-alpha",
		SortField:     "updated_at",
		SortDirection: "desc",
		LastValue:     "2026-10-31T12:00:00Z",
		TiebreakerID:  "skl-01",
		Fingerprint:   "fp-integrity-check",
	}

	token := EncodeCursor(payload)
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		t.Fatalf("expected signed cursor token to have format '<payload>.<signature>', got: %s", token)
	}

	// 1. Bit flip across all bytes of signature
	sigBytes := []byte(parts[1])
	for i := 0; i < len(sigBytes); i++ {
		tamperedSig := make([]byte, len(sigBytes))
		copy(tamperedSig, sigBytes)
		tamperedSig[i] ^= 0x01
		tamperedToken := parts[0] + "." + string(tamperedSig)
		if _, err := DecodeCursor(tamperedToken); err == nil {
			t.Fatalf("SECURITY FLAW: DecodeCursor accepted cursor with bit flipped at byte %d of signature!", i)
		}
	}

	// 2. Bit flip in payload
	payloadBytes := []byte(parts[0])
	for i := 0; i < len(payloadBytes); i++ {
		tamperedPayload := make([]byte, len(payloadBytes))
		copy(tamperedPayload, payloadBytes)
		tamperedPayload[i] ^= 0x01
		tamperedToken := string(tamperedPayload) + "." + parts[1]
		if _, err := DecodeCursor(tamperedToken); err == nil {
			t.Fatalf("SECURITY FLAW: DecodeCursor accepted cursor with bit flipped at byte %d of payload!", i)
		}
	}

	// 3. Forged payload fields: change TiebreakerID or TenantID
	b, _ := base64.RawURLEncoding.DecodeString(parts[0])
	var decodedPayload CursorPayload
	_ = json.Unmarshal(b, &decodedPayload)
	decodedPayload.TiebreakerID = "skl-tampered"
	newB, _ := json.Marshal(decodedPayload)
	forgedPayloadB64 := base64.RawURLEncoding.EncodeToString(newB)
	forgedToken := forgedPayloadB64 + "." + parts[1]
	if _, err := DecodeCursor(forgedToken); err == nil {
		t.Fatalf("SECURITY FLAW: DecodeCursor accepted forged payload with mismatched signature!")
	}

	// 4. Unsigned legacy token rejection
	if _, err := DecodeCursor(parts[0]); err == nil {
		t.Fatalf("SECURITY FLAW: DecodeCursor accepted unsigned token without dot separator!")
	}

	// 5. Extra dot / malformed parts rejection
	if _, err := DecodeCursor(parts[0] + "." + parts[1] + ".extra"); err == nil {
		t.Fatalf("SECURITY FLAW: DecodeCursor accepted token with 3 parts!")
	}

	// 6. Schema check rejection (empty TiebreakerID or Version != 1)
	invalidSchemaPayload := CursorPayload{
		Version:      2, // Invalid version
		TenantID:     "tenant-alpha",
		TiebreakerID: "", // Missing tiebreaker
	}
	signedInvalidSchemaToken := EncodeCursor(invalidSchemaPayload)
	if _, err := DecodeCursor(signedInvalidSchemaToken); err == nil {
		t.Fatalf("SECURITY FLAW: DecodeCursor accepted invalid schema payload!")
	}
}
