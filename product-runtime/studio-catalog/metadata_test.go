package studiocatalog

import (
	"context"
	"testing"
	"time"
)

func createTestAsset(tenantID, assetID, name string, kind AssetKind, state LifecycleState) *Asset {
	return &Asset{
		ID:                 assetID,
		TenantID:           tenantID,
		EnvironmentID:      "env-prod",
		Type:               kind,
		Name:               name,
		DisplayName:        "Display " + name,
		Description:        "Description for " + name,
		Tags:               []string{"secops", "test"},
		Owner: Owner{
			OwnerID:     "usr-01",
			OwnerName:   "Alice Sec",
			OwnerType:   "user",
			IsUnknown:   false,
			DisplayName: "Alice Sec",
		},
		Contributors: []string{"usr-02"},
		Version: VersionInfo{
			VersionID:   "v-1.0.0",
			SemVer:      "1.0.0",
			CreatedAt:   time.Now().UTC(),
			IsImmutable: true,
		},
		LifecycleState:     state,
		FunctionalState:    FunctionalStateReady,
		DeliveryMode:       DeliveryModeDelivered,
		ExecutionClass:     "read-only",
		RiskClassification: RiskClassificationLow,
		IsAIPowered:        false,
		HasNonAIFallback:   true,
		Parameters: []Parameter{
			{Name: "param1", Type: "string", Required: true},
		},
		Dependencies: []Dependency{
			{Kind: "tool", RefID: "tool-01"},
		},
		Consumers: []Consumer{
			{ConsumerType: "workflow", ConsumerRef: "wf-01"},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// =========================================================================
// 1. Metadata Validation Tests (REQ-PROD-016, CAP-STD-002)
// =========================================================================

func TestValidateAssetMetadata_CompleteAsset(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "skl-100", "test-skill", AssetKindSkill, LifecycleStatePublished)
	res := ValidateAssetMetadata(asset)

	if !res.Valid {
		t.Fatalf("expected valid metadata, got issues: %+v", res.Issues)
	}
	if res.State != FunctionalStateComplete {
		t.Fatalf("expected functional state complete, got %s", res.State)
	}
	if len(res.MissingRequiredFields) != 0 {
		t.Fatalf("expected zero missing required fields, got %v", res.MissingRequiredFields)
	}
}

func TestValidateAssetMetadata_PartialAsset_MissingRiskAndDelivery(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "skl-101", "partial-skill", AssetKindSkill, LifecycleStateDraft)
	asset.RiskClassification = ""
	asset.DeliveryMode = ""

	res := ValidateAssetMetadata(asset)
	if !res.Valid {
		t.Fatalf("expected valid=true for partial asset, got invalid with issues: %+v", res.Issues)
	}
	if res.State != FunctionalStatePartial {
		t.Fatalf("expected functional state partial, got %s", res.State)
	}
}

func TestValidateAssetMetadata_MissingRequiredFields_DraftInvalid(t *testing.T) {
	testCases := []struct {
		name         string
		mutator      func(a *Asset)
		expectedCode string
	}{
		{
			name:         "missing ID",
			mutator:      func(a *Asset) { a.ID = "" },
			expectedCode: "missing_id",
		},
		{
			name:         "missing TenantID",
			mutator:      func(a *Asset) { a.TenantID = "" },
			expectedCode: "missing_tenant_id",
		},
		{
			name:         "invalid Type",
			mutator:      func(a *Asset) { a.Type = "unknown-type" },
			expectedCode: "invalid_type",
		},
		{
			name:         "missing Name",
			mutator:      func(a *Asset) { a.Name = "" },
			expectedCode: "missing_name",
		},
		{
			name:         "missing SemVer",
			mutator:      func(a *Asset) { a.Version.SemVer = "" },
			expectedCode: "missing_semver",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			asset := createTestAsset("tenant-alpha", "skl-err", "err-skill", AssetKindSkill, LifecycleStateDraft)
			tc.mutator(asset)

			res := ValidateAssetMetadata(asset)
			if res.Valid {
				t.Fatalf("expected validation failure for %s", tc.name)
			}
			if res.State != FunctionalStateDraft {
				t.Fatalf("expected functional state draft, got %s", res.State)
			}

			found := false
			for _, issue := range res.Issues {
				if issue.Code == tc.expectedCode {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected issue code %s, got: %+v", tc.expectedCode, res.Issues)
			}
		})
	}
}

func TestValidateAssetMetadata_NilAsset(t *testing.T) {
	res := ValidateAssetMetadata(nil)
	if res.Valid {
		t.Fatalf("nil asset must be invalid")
	}
	if len(res.MissingRequiredFields) == 0 {
		t.Fatalf("expected missing required fields for nil asset")
	}
}

func TestValidateAssetMetadata_LifecycleStatePropagation(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "skl-102", "dep-skill", AssetKindSkill, LifecycleStateDeprecated)
	res := ValidateAssetMetadata(asset)
	if res.State != FunctionalStateDeprecated {
		t.Fatalf("expected functional state deprecated, got %s", res.State)
	}

	asset.LifecycleState = LifecycleStateSuperseded
	asset.Supersession.IsSuperseded = true
	resSup := ValidateAssetMetadata(asset)
	if resSup.State != FunctionalStateSuperseded {
		t.Fatalf("expected functional state superseded, got %s", resSup.State)
	}
}

func TestValidateMetadata_ServiceMethod(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	asset := createTestAsset("tenant-alpha", "skl-srv-01", "srv-skill", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, asset)

	// Nominal
	resp, err := svc.ValidateMetadata(ctx, auth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-srv-01",
	})
	if err != nil {
		t.Fatalf("ValidateMetadata failed: %v", err)
	}
	if !resp.Validation.Valid {
		t.Fatalf("expected valid metadata")
	}

	// 404 Not Found
	_, errNotFound := svc.ValidateMetadata(ctx, auth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "non-existent-id",
	})
	if errNotFound == nil {
		t.Fatalf("expected 404 for non-existent asset")
	}
	env, ok := errNotFound.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 404 || env.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected uniform 404 asset-not-found, got: %v", errNotFound)
	}

	// Cross-tenant uniform 404 oracle elimination
	_, errCross := svc.ValidateMetadata(ctx, auth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-01", // in tenant-bravo
	})
	if errCross == nil {
		t.Fatalf("expected 404 for cross-tenant asset")
	}
	envCross, ok := errCross.(*ErrorEnvelope)
	if !ok || envCross.HTTPStatus != 404 || envCross.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected uniform 404 on cross-tenant metadata inspection, got: %v", errCross)
	}

	// 403 Missing Permission
	unauth := NewAuthContext("usr-unauth", "tenant-alpha", nil) // no permissions
	_, errPerm := svc.ValidateMetadata(ctx, unauth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-srv-01",
	})
	if errPerm == nil {
		t.Fatalf("expected 403 for unauthorized caller")
	}
	envPerm, ok := errPerm.(*ErrorEnvelope)
	if !ok || envPerm.HTTPStatus != 403 {
		t.Fatalf("expected 403, got: %v", errPerm)
	}

	// 400 Missing Asset ID
	_, errMissingID := svc.ValidateMetadata(ctx, auth, ValidateMetadataRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "   ",
	})
	if errMissingID == nil {
		t.Fatalf("expected 400 for empty asset ID")
	}
}

// =========================================================================
// 2. Owner Unknown Fallback Tests (REQ-OBJ-009)
// =========================================================================

func TestOwnerUnknown_FallbackInValidation(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "skl-unk-01", "unknown-owner-skill", AssetKindSkill, LifecycleStateDraft)
	asset.Owner = Owner{
		OwnerID:   "",
		OwnerName: "",
	}
	asset.Owner.Normalize()

	if !asset.Owner.IsUnknown {
		t.Fatalf("expected IsUnknown to be true after normalize")
	}
	if asset.Owner.OwnerName != "owner unknown" {
		t.Fatalf("expected owner name 'owner unknown', got %s", asset.Owner.OwnerName)
	}

	res := ValidateAssetMetadata(asset)
	// Must be valid, but functional state partial because owner is unknown
	if !res.Valid {
		t.Fatalf("asset with unknown owner must still be valid")
	}
	if res.State != FunctionalStatePartial {
		t.Fatalf("expected partial state for unknown owner, got %s", res.State)
	}

	found := false
	for _, issue := range res.Issues {
		if issue.Code == "owner_unknown" && issue.Severity == SeverityInfo {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected owner_unknown info issue, got: %+v", res.Issues)
	}
}

func TestAssignOwner_ServiceNominalAndFallback(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	asset := createTestAsset("tenant-alpha", "skl-own-01", "assign-owner-skill", AssetKindSkill, LifecycleStateDraft)
	_ = store.Put(ctx, asset)

	// 1. Assign known owner
	resp, err := svc.AssignOwner(ctx, auth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-own-01",
		NewOwner: Owner{
			OwnerID:   "usr-bob",
			OwnerName: "Bob Incident",
			OwnerType: "user",
		},
		Actor:  "lead-architect",
		Reason: "Handover to Bob",
	})
	if err != nil {
		t.Fatalf("AssignOwner failed: %v", err)
	}
	if resp.Asset.Owner.OwnerID != "usr-bob" || resp.Asset.Owner.IsUnknown {
		t.Fatalf("expected owner usr-bob, got %+v", resp.Asset.Owner)
	}
	if resp.PreviousOwner.OwnerID != "usr-01" {
		t.Fatalf("expected previous owner usr-01, got %s", resp.PreviousOwner.OwnerID)
	}
	if len(resp.Asset.AuditTrail) == 0 {
		t.Fatalf("expected audit trail record")
	}
	lastAudit := resp.Asset.AuditTrail[len(resp.Asset.AuditTrail)-1]
	if lastAudit.Action != "assign_owner" || lastAudit.Actor != "usr-admin" {
		t.Fatalf("audit trail mismatch: %+v", lastAudit)
	}

	// 2. Assign blank owner -> explicit owner unknown fallback (REQ-OBJ-009)
	respUnknown, errUnk := svc.AssignOwner(ctx, auth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-own-01",
		NewOwner: Owner{
			OwnerID:   "",
			OwnerName: "",
		},
		Actor:  "lead-architect",
		Reason: "Owner departure",
	})
	if errUnk != nil {
		t.Fatalf("AssignOwner unknown failed: %v", errUnk)
	}
	if !respUnknown.Asset.Owner.IsUnknown || respUnknown.Asset.Owner.OwnerName != "owner unknown" {
		t.Fatalf("expected explicit owner unknown, got %+v", respUnknown.Asset.Owner)
	}

	// 3. Permission denied: caller lacks manage permission for asset kind
	readerAuth := NewAuthContext("usr-reader", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioSkillRead)
	_, errPerm := svc.AssignOwner(ctx, readerAuth, AssignOwnerRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-own-01",
		NewOwner: Owner{
			OwnerID:   "usr-charlie",
			OwnerName: "Charlie",
		},
	})
	if errPerm == nil {
		t.Fatalf("expected 403 for read-only user attempting owner assignment")
	}
	envPerm, ok := errPerm.(*ErrorEnvelope)
	if !ok || envPerm.HTTPStatus != 403 {
		t.Fatalf("expected 403 permission denied, got %v", errPerm)
	}
}

// =========================================================================
// 3. Deterministic Tag and Risk Suggestion Tests (REQ-PROD-016)
// =========================================================================

func TestSuggestClassificationAndTags_Deterministic(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "skl-sug-01", "incident-detection-skill", AssetKindSkill, LifecycleStatePublished)
	asset.Description = "Detect and report active SIEM threat incidents"

	sug1 := SuggestAssetTagsAndClassification(asset)
	sug2 := SuggestAssetTagsAndClassification(asset)

	// Strict determinism check
	if sug1.SuggestedRisk != sug2.SuggestedRisk {
		t.Fatalf("non-deterministic risk suggestion: %v vs %v", sug1.SuggestedRisk, sug2.SuggestedRisk)
	}
	if len(sug1.SuggestedTags) != len(sug2.SuggestedTags) {
		t.Fatalf("non-deterministic tag count: %d vs %d", len(sug1.SuggestedTags), len(sug2.SuggestedTags))
	}
	if sug1.IsAIPowered || sug2.IsAIPowered {
		t.Fatalf("IsAIPowered must be strictly false per REQ-PROD-016")
	}

	// Verify keywords matched
	foundDetection := false
	for _, t := range sug1.SuggestedTags {
		if t == "detection" {
			foundDetection = true
			break
		}
	}
	if !foundDetection {
		t.Fatalf("expected tag 'detection' from keywords, got: %v", sug1.SuggestedTags)
	}
}

func TestSuggestClassificationAndTags_CriticalDestructiveActions(t *testing.T) {
	asset := createTestAsset("tenant-alpha", "agn-sug-02", "endpoint-isolation-agent", AssetKindAutomationAgent, LifecycleStateDraft)
	asset.Description = "Contains and terminates malicious host processes"
	asset.Dependencies = append(asset.Dependencies, Dependency{
		Kind:  "secret_reference",
		RefID: "vault://keys/containment-token",
	})

	sug := SuggestAssetTagsAndClassification(asset)
	if sug.SuggestedRisk != RiskClassificationCritical {
		t.Fatalf("expected Critical risk for containment/terminate capability, got %s", sug.SuggestedRisk)
	}
	if sug.IsAIPowered {
		t.Fatalf("IsAIPowered must be false")
	}
}

func TestSuggestClassificationAndTags_Service(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	asset := createTestAsset("tenant-alpha", "skl-sug-03", "forensic-evidence-collector", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, asset)

	resp, err := svc.SuggestClassificationAndTags(ctx, auth, SuggestTagsRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-sug-03",
	})
	if err != nil {
		t.Fatalf("SuggestClassificationAndTags failed: %v", err)
	}
	if resp.Suggestions.IsAIPowered {
		t.Fatalf("expected IsAIPowered=false")
	}
	if len(resp.Suggestions.SuggestedTags) == 0 {
		t.Fatalf("expected suggested tags")
	}
}

// =========================================================================
// 4. Draft Metadata Update Tests (Class 2)
// =========================================================================

func TestUpdateDraftMetadata_Nominal(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	asset := createTestAsset("tenant-alpha", "skl-draft-01", "draft-skill", AssetKindSkill, LifecycleStateDraft)
	_ = store.Put(ctx, asset)

	aiFlag := true
	fallbackFlag := true
	resp, err := svc.UpdateDraftMetadata(ctx, auth, UpdateDraftMetadataRequest{
		TenantRef:          "tenant-alpha",
		AssetID:            "skl-draft-01",
		DisplayName:        "Updated Draft Display",
		Description:        "Updated draft description",
		Tags:               []string{"new-tag", "incident"},
		RiskClassification: RiskClassificationMedium,
		DeliveryMode:       DeliveryModePreview,
		IsAIPowered:        &aiFlag,
		HasNonAIFallback:   &fallbackFlag,
		Actor:              "designer-01",
		Reason:             "Refining draft parameters",
	})
	if err != nil {
		t.Fatalf("UpdateDraftMetadata failed: %v", err)
	}
	if resp.Asset.DisplayName != "Updated Draft Display" {
		t.Fatalf("expected updated display name, got %s", resp.Asset.DisplayName)
	}
	if resp.Asset.RiskClassification != RiskClassificationMedium {
		t.Fatalf("expected updated risk medium, got %s", resp.Asset.RiskClassification)
	}
	if resp.Asset.DeliveryMode != DeliveryModePreview {
		t.Fatalf("expected updated delivery preview, got %s", resp.Asset.DeliveryMode)
	}
	if !resp.Asset.IsAIPowered || !resp.Asset.HasNonAIFallback {
		t.Fatalf("expected updated AI flags")
	}
	if len(resp.Asset.AuditTrail) == 0 {
		t.Fatalf("expected audit trail entry for draft update")
	}
}

func TestUpdateDraftMetadata_PublishedAsset_Blocked(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Published assets are immutable under CAP-STD-002
	asset := createTestAsset("tenant-alpha", "skl-pub-01", "published-skill", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, asset)

	_, err := svc.UpdateDraftMetadata(ctx, auth, UpdateDraftMetadataRequest{
		TenantRef:   "tenant-alpha",
		AssetID:     "skl-pub-01",
		DisplayName: "Attempted Mutate",
	})
	if err == nil {
		t.Fatalf("expected error mutating published asset")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 400 || env.Code != ErrorCodeImmutableAssetMutation {
		t.Fatalf("expected 400 immutable-asset-mutation, got: %v", err)
	}
}

// =========================================================================
// 5. Supersession Linking & Lineage Navigation Tests (CAP-STD-002)
// =========================================================================

func TestRecordSupersession_NominalLink(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	v1 := createTestAsset("tenant-alpha", "skl-v1", "skill-v1", AssetKindSkill, LifecycleStatePublished)
	v1.Version.SemVer = "1.0.0"
	_ = store.Put(ctx, v1)

	v2 := createTestAsset("tenant-alpha", "skl-v2", "skill-v2", AssetKindSkill, LifecycleStatePublished)
	v2.Version.SemVer = "2.0.0"
	_ = store.Put(ctx, v2)

	resp, err := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-v1",
		NewerAssetID: "skl-v2",
		Actor:        "version-manager",
		Reason:       "Version 2.0.0 replaces 1.0.0 with enhanced telemetry",
	})
	if err != nil {
		t.Fatalf("RecordSupersession failed: %v", err)
	}

	if resp.OlderAsset.LifecycleState != LifecycleStateSuperseded {
		t.Fatalf("expected older asset state superseded, got %s", resp.OlderAsset.LifecycleState)
	}

	// Verify persistence in store
	updatedV1, _ := store.Get(ctx, "tenant-alpha", "skl-v1")
	if !updatedV1.Supersession.IsSuperseded || updatedV1.Supersession.SupersededBy != "skl-v2" {
		t.Fatalf("v1 supersession not persisted correctly: %+v", updatedV1.Supersession)
	}
	if len(updatedV1.AuditTrail) == 0 {
		t.Fatalf("v1 missing audit trail")
	}

	updatedV2, _ := store.Get(ctx, "tenant-alpha", "skl-v2")
	if updatedV2.Supersession.Supersedes != "skl-v1" {
		t.Fatalf("v2 supersession link not persisted correctly: %+v", updatedV2.Supersession)
	}
	if len(updatedV2.AuditTrail) == 0 {
		t.Fatalf("v2 missing audit trail")
	}
}

func TestRecordSupersession_ValidationErrors(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	v1 := createTestAsset("tenant-alpha", "skl-err-v1", "v1", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, v1)

	// Self-supersession
	_, errSelf := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-err-v1",
		NewerAssetID: "skl-err-v1",
	})
	if errSelf == nil {
		t.Fatalf("expected error for self-supersession")
	}

	// Non-existent newer asset -> 404
	_, errNotFound := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "skl-err-v1",
		NewerAssetID: "skl-non-existent",
	})
	if errNotFound == nil {
		t.Fatalf("expected 404 for non-existent newer asset")
	}
}

func TestGetSupersessionHistory_ThreeNodeLineage(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Setup 3-generation lineage: v1 -> v2 -> v3
	v1 := createTestAsset("tenant-alpha", "wf-v1", "workflow-v1", AssetKindWorkflow, LifecycleStateSuperseded)
	v1.Version.SemVer = "1.0.0"
	v1.Supersession = Supersession{IsSuperseded: true, SupersededBy: "wf-v2"}
	_ = store.Put(ctx, v1)

	v2 := createTestAsset("tenant-alpha", "wf-v2", "workflow-v2", AssetKindWorkflow, LifecycleStateSuperseded)
	v2.Version.SemVer = "2.0.0"
	v2.Supersession = Supersession{IsSuperseded: true, Supersedes: "wf-v1", SupersededBy: "wf-v3"}
	_ = store.Put(ctx, v2)

	v3 := createTestAsset("tenant-alpha", "wf-v3", "workflow-v3", AssetKindWorkflow, LifecycleStatePublished)
	v3.Version.SemVer = "3.0.0"
	v3.Supersession = Supersession{IsSuperseded: false, Supersedes: "wf-v2"}
	_ = store.Put(ctx, v3)

	// Test navigating history starting from middle node (v2)
	respV2, err := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-v2",
	})
	if err != nil {
		t.Fatalf("GetSupersessionHistory from v2 failed: %v", err)
	}

	if respV2.RootAssetID != "wf-v1" {
		t.Fatalf("expected RootAssetID wf-v1, got %s", respV2.RootAssetID)
	}
	if respV2.LatestAssetID != "wf-v3" {
		t.Fatalf("expected LatestAssetID wf-v3, got %s", respV2.LatestAssetID)
	}
	if respV2.CurrentAssetID != "wf-v2" {
		t.Fatalf("expected CurrentAssetID wf-v2, got %s", respV2.CurrentAssetID)
	}
	if len(respV2.Chain) != 3 {
		t.Fatalf("expected 3 nodes in chain, got %d", len(respV2.Chain))
	}

	// Verify order: wf-v1, wf-v2, wf-v3
	expectedOrder := []string{"wf-v1", "wf-v2", "wf-v3"}
	for i, node := range respV2.Chain {
		if node.AssetID != expectedOrder[i] {
			t.Fatalf("chain order mismatch at %d: expected %s, got %s", i, expectedOrder[i], node.AssetID)
		}
	}

	// Verify active and current flags
	if !respV2.Chain[1].IsCurrent {
		t.Fatalf("wf-v2 must have IsCurrent=true")
	}
	if respV2.Chain[0].IsActive || respV2.Chain[1].IsActive {
		t.Fatalf("v1 and v2 should not be active")
	}
	if !respV2.Chain[2].IsActive {
		t.Fatalf("v3 must be active")
	}

	// Test navigating history starting from root (v1)
	respV1, errV1 := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-v1",
	})
	if errV1 != nil || len(respV1.Chain) != 3 {
		t.Fatalf("GetSupersessionHistory from v1 failed: %v", errV1)
	}

	// Test navigating history starting from tip (v3)
	respV3, errV3 := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "wf-v3",
	})
	if errV3 != nil || len(respV3.Chain) != 3 {
		t.Fatalf("GetSupersessionHistory from v3 failed: %v", errV3)
	}
}

func TestGetSupersessionHistory_StandaloneAsset(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	asset := createTestAsset("tenant-alpha", "tool-solo", "solo-tool", AssetKindTool, LifecycleStatePublished)
	_ = store.Put(ctx, asset)

	resp, err := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "tool-solo",
	})
	if err != nil {
		t.Fatalf("GetSupersessionHistory for standalone asset failed: %v", err)
	}
	if len(resp.Chain) != 1 {
		t.Fatalf("expected chain length 1, got %d", len(resp.Chain))
	}
	if resp.RootAssetID != "tool-solo" || resp.LatestAssetID != "tool-solo" {
		t.Fatalf("expected root and latest to equal tool-solo")
	}
	if !resp.Chain[0].IsActive || !resp.Chain[0].IsCurrent {
		t.Fatalf("standalone asset must be active and current")
	}
}

func TestGetSupersessionHistory_CrossTenantUniform404(t *testing.T) {
	_, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Attempt to inspect supersession of an asset in tenant-bravo
	_, err := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "skl-bravo-01",
	})
	if err == nil {
		t.Fatalf("expected 404 for cross-tenant asset")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 404 || env.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected uniform 404 asset-not-found, got: %v", err)
	}
}

// =========================================================================
// 6. Adversarial Stress Probes (Cycle Detection, Owner Fallback, Immutability)
// =========================================================================

func TestRecordSupersession_CircularDetection(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	a := createTestAsset("tenant-alpha", "node-A", "asset-A", AssetKindSkill, LifecycleStatePublished)
	b := createTestAsset("tenant-alpha", "node-B", "asset-B", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, a)
	_ = store.Put(ctx, b)

	// Step 1: A is superseded by B (A -> B)
	_, err := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "node-A",
		NewerAssetID: "node-B",
	})
	if err != nil {
		t.Fatalf("unexpected error recording supersession: %v", err)
	}

	// Step 2: Adversarial reverse attempt (B -> A). Must fail with circular supersession error.
	_, errCycle := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "node-B",
		NewerAssetID: "node-A",
	})
	if errCycle == nil {
		t.Fatalf("expected circular supersession error for reverse link")
	}
	env, ok := errCycle.(*ErrorEnvelope)
	if !ok || env.Code != ErrorCodeCircularSupersession {
		t.Fatalf("expected ErrorCodeCircularSupersession, got: %v", errCycle)
	}
}

func TestRecordSupersession_AlreadySupersededConflict(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	a := createTestAsset("tenant-alpha", "node-orig", "asset-orig", AssetKindSkill, LifecycleStatePublished)
	b := createTestAsset("tenant-alpha", "node-rep1", "asset-rep1", AssetKindSkill, LifecycleStatePublished)
	c := createTestAsset("tenant-alpha", "node-rep2", "asset-rep2", AssetKindSkill, LifecycleStatePublished)
	_ = store.Put(ctx, a)
	_ = store.Put(ctx, b)
	_ = store.Put(ctx, c)

	// A is superseded by B
	_, err := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "node-orig",
		NewerAssetID: "node-rep1",
	})
	if err != nil {
		t.Fatalf("initial supersession failed: %v", err)
	}

	// Attempt to supersede A by C
	_, errConflict := svc.RecordSupersession(ctx, auth, RecordSupersessionRequest{
		TenantRef:    "tenant-alpha",
		OlderAssetID: "node-orig",
		NewerAssetID: "node-rep2",
	})
	if errConflict == nil {
		t.Fatalf("expected conflict error for already superseded asset")
	}
	env, ok := errConflict.(*ErrorEnvelope)
	if !ok || env.HTTPStatus != 400 || env.Code != ErrorCodeInvalidMetadata {
		t.Fatalf("expected 400 invalid-metadata, got: %v", errConflict)
	}
}

func TestGetSupersessionHistory_CycleGuard(t *testing.T) {
	store, svc := setupTestStore(t)
	ctx := context.Background()
	auth := adminAuth("tenant-alpha")

	// Inject a 2-node cycle directly into store: A <-> B
	a := createTestAsset("tenant-alpha", "node-c1", "asset-c1", AssetKindSkill, LifecycleStateSuperseded)
	a.Supersession = Supersession{IsSuperseded: true, Supersedes: "node-c2", SupersededBy: "node-c2"}
	_ = store.Put(ctx, a)

	b := createTestAsset("tenant-alpha", "node-c2", "asset-c2", AssetKindSkill, LifecycleStateSuperseded)
	b.Supersession = Supersession{IsSuperseded: true, Supersedes: "node-c1", SupersededBy: "node-c1"}
	_ = store.Put(ctx, b)

	// History navigation must safely detect cycle without infinite loop
	_, err := svc.GetSupersessionHistory(ctx, auth, SupersessionHistoryRequest{
		TenantRef: "tenant-alpha",
		AssetID:   "node-c1",
	})
	if err == nil {
		t.Fatalf("expected circular supersession error in history navigation")
	}
	env, ok := err.(*ErrorEnvelope)
	if !ok || env.Code != ErrorCodeCircularSupersession {
		t.Fatalf("expected ErrorCodeCircularSupersession, got: %v", err)
	}
}

func TestUpdateDraftMetadata_AllNonDraftStates_Blocked(t *testing.T) {
	nonDraftStates := []LifecycleState{
		LifecycleStateReview,
		LifecycleStateApproved,
		LifecycleStatePublished,
		LifecycleStateDeprecated,
		LifecycleStateSuperseded,
		LifecycleStateRetired,
		LifecycleStateArchived,
	}

	for _, state := range nonDraftStates {
		t.Run(string(state), func(t *testing.T) {
			store, svc := setupTestStore(t)
			ctx := context.Background()
			auth := adminAuth("tenant-alpha")

			assetID := "asset-" + string(state)
			asset := createTestAsset("tenant-alpha", assetID, "name-"+string(state), AssetKindSkill, state)
			_ = store.Put(ctx, asset)

			_, err := svc.UpdateDraftMetadata(ctx, auth, UpdateDraftMetadataRequest{
				TenantRef:   "tenant-alpha",
				AssetID:     assetID,
				DisplayName: "Tampered Display Name",
			})
			if err == nil {
				t.Fatalf("expected mutation error for non-draft state %q", state)
			}
			env, ok := err.(*ErrorEnvelope)
			if !ok || env.HTTPStatus != 400 || env.Code != ErrorCodeImmutableAssetMutation {
				t.Fatalf("expected 400 immutable-asset-mutation for state %q, got: %v", state, err)
			}
		})
	}
}

func TestOwnerUnknown_WhitespaceAndPartialMatrix(t *testing.T) {
	matrix := []struct {
		name        string
		id          string
		ownerName   string
		wantUnknown bool
	}{
		{"empty strings", "", "", true},
		{"whitespace only", "   ", " \t \n ", true},
		{"valid id empty name", "usr-1", "   ", true},
		{"empty id valid name", "   ", "Alice", true},
		{"both valid", "usr-1", "Alice", false},
	}

	for _, tc := range matrix {
		t.Run(tc.name, func(t *testing.T) {
			o := Owner{
				OwnerID:   tc.id,
				OwnerName: tc.ownerName,
			}
			o.Normalize()

			if o.IsUnknown != tc.wantUnknown {
				t.Fatalf("expected IsUnknown=%t, got %t", tc.wantUnknown, o.IsUnknown)
			}
			if tc.wantUnknown && (o.OwnerID != "unknown" || o.OwnerName != "owner unknown") {
				t.Fatalf("expected normalized unknown owner, got %+v", o)
			}
		})
	}
}
