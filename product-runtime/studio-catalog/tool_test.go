package studiocatalog

import (
	"context"
	"strings"
	"testing"
	"time"
)

// Helper to construct a test ToolService populated with sample tools
func setupTestToolService(t *testing.T) (*DefaultToolService, *InMemoryToolStore) {
	t.Helper()
	store := NewInMemoryToolStore()
	svc := NewDefaultToolService(store)

	ctx := context.Background()
	now := time.Now().UTC()

	// 1. Nominal draft tool in tenant-alpha
	draftTool := &ToolDefinition{
		ID:                  "tool-text-formatter",
		TenantID:            "tenant-alpha",
		EnvironmentID:       "env-prod",
		Name:                "Text Formatter Tool",
		DisplayName:         "Text Formatter",
		Description:         "Performs deterministic text formatting and case transformations",
		Purpose:             "Format strings into uppercase, lowercase, or trimmed representations",
		Owner:               Owner{OwnerID: "dev-team", OwnerName: "DevOps", OwnerType: "team"},
		Version:             VersionInfo{VersionID: "v1.0.0", SemVer: "1.0.0", CreatedAt: now, IsImmutable: false},
		FunctionalState:     ToolStateDraft,
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"to_upper", "to_lower", "trim_space"},
		Operations: []ToolOperation{
			{
				OperationID: "op-to-upper",
				Name:        "to_upper",
				Description: "Convert text to uppercase",
				Parameters: []Parameter{
					{Name: "input", Type: "string", Required: true, Description: "Target string input"},
				},
				IsIdempotent: true,
				SideEffectContext: SideEffectContext{
					HasSideEffects: false,
					MutationScope:  "read-only",
					IsReversible:   true,
				},
			},
		},
		Parameters: []Parameter{
			{Name: "input", Type: "string", Required: true, Description: "Target string input"},
		},
		Dependencies: []Dependency{
			{Kind: "runtime-host", RefID: "host-app-server", HealthStatus: "available"},
		},
		ActionClass:        ActionClass0Observation,
		RiskClassification: RiskClassificationLow,
		SideEffects: SideEffectContext{
			HasSideEffects: false,
			MutationScope:  "read-only",
			IsReversible:   true,
		},
		Limitations: []string{"Maximum input length 1024 characters"},
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  true,
			HealthStatus: "available",
			CheckedAt:    now,
			Message:      "App server host available",
		},
		Provenance: Provenance{
			CreatedBy: "alice",
			CreatedAt: now,
			UpdatedBy: "alice",
			UpdatedAt: now,
		},
		CreatedAt: now,
		UpdatedAt: now,
		AuditTrail: []AuditRecord{
			{Timestamp: now, Actor: "alice", Reason: "initial draft creation", Action: "create_draft"},
		},
	}
	if err := store.Put(ctx, draftTool); err != nil {
		t.Fatalf("setup: failed to seed draftTool: %v", err)
	}

	// 2. Reference tool in tenant-alpha (immutable, available-reference)
	refTool := &ToolDefinition{
		ID:                  "tool-currency-converter",
		TenantID:            "tenant-alpha",
		Name:                "Currency Converter Tool",
		DisplayName:         "Currency Converter",
		Description:         "Converts monetary amounts using configured exchange rates",
		Purpose:             "Convert currency amounts during invoice generation",
		Owner:               Owner{OwnerID: "finance-team", OwnerName: "Finance Team", OwnerType: "team"},
		Version:             VersionInfo{VersionID: "v1.2.0", SemVer: "1.2.0", CreatedAt: now, IsImmutable: true},
		FunctionalState:     ToolStateAvailableReference,
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"convert_amount"},
		Operations: []ToolOperation{
			{
				OperationID: "op-convert",
				Name:        "convert_amount",
				Description: "Convert amount between currency codes",
				Parameters:  []Parameter{{Name: "amount", Type: "number", Required: true}},
			},
		},
		Parameters: []Parameter{{Name: "amount", Type: "number", Required: true}},
		Dependencies: []Dependency{
			{Kind: "rate-table", RefID: "table-rates-v1", HealthStatus: "available"},
		},
		ActionClass:        ActionClass0Observation,
		RiskClassification: RiskClassificationLow,
		SideEffects:        SideEffectContext{HasSideEffects: false, IsReversible: true},
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  true,
			HealthStatus: "available",
			CheckedAt:    now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := store.Put(ctx, refTool); err != nil {
		t.Fatalf("setup: failed to seed refTool: %v", err)
	}

	// 3. Offline runtime provider tool in tenant-alpha (CAP-STD-003 §25 Given 1)
	offlineTool := &ToolDefinition{
		ID:                  "tool-offline-validator",
		TenantID:            "tenant-alpha",
		Name:                "JSON Schema Validator",
		DisplayName:         "Schema Validator",
		Description:         "Validates JSON documents against schemas",
		Purpose:             "Validate API payloads prior to persistence",
		Owner:               Owner{OwnerID: "qa-team", OwnerName: "QA Team", OwnerType: "team"},
		Version:             VersionInfo{VersionID: "v0.9.0", SemVer: "0.9.0", CreatedAt: now},
		FunctionalState:     ToolStateDraft,
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"validate_schema"},
		Operations: []ToolOperation{
			{
				OperationID: "op-validate",
				Name:        "validate_schema",
				Parameters:  []Parameter{{Name: "payload", Type: "string", Required: true}},
			},
		},
		Parameters: []Parameter{{Name: "payload", Type: "string", Required: true}},
		ActionClass:        ActionClass1Assessment,
		RiskClassification: RiskClassificationMedium,
		SideEffects:        SideEffectContext{HasSideEffects: true, MutationScope: "local-session", IsReversible: true},
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  false, // Runtime daemon offline!
			HealthStatus: "unavailable",
			CheckedAt:    now,
			Message:      "Validator daemon offline",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := store.Put(ctx, offlineTool); err != nil {
		t.Fatalf("setup: failed to seed offlineTool: %v", err)
	}

	// 4. Foreign tool belonging to tenant-bravo (for anti-oracle tests)
	foreignTool := &ToolDefinition{
		ID:                  "tool-bravo-calculator",
		TenantID:            "tenant-bravo",
		Name:                "Tenant Bravo Tax Calculator",
		DisplayName:         "Bravo Tax Calculator",
		Description:         "Calculates local sales tax for regional orders",
		Purpose:             "Compute regional sales taxes for checkout",
		Owner:               Owner{OwnerID: "bravo-admin", OwnerName: "Bravo Admin"},
		FunctionalState:     ToolStateAvailableReference,
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"compute_tax"},
		Operations: []ToolOperation{
			{
				OperationID: "op-bravo-tax",
				Name:        "compute_tax",
				Parameters:  []Parameter{{Name: "subtotal", Type: "number", Required: true}},
			},
		},
		RuntimeAvailability: RuntimeAvailabilityStatus{IsAvailable: true, HealthStatus: "available", CheckedAt: now},
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if err := store.Put(ctx, foreignTool); err != nil {
		t.Fatalf("setup: failed to seed foreignTool: %v", err)
	}

	return svc, store
}

// =============================================================================
// Task 3.1: Class 0 InspectTool Tests
// =============================================================================

func TestTool_Inspect_Nominal(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	resp, err := svc.InspectTool(ctx, &auth, "tenant-alpha", "tool-text-formatter")
	if err != nil {
		t.Fatalf("expected inspect to succeed, got %v", err)
	}

	if resp == nil || resp.Tool == nil {
		t.Fatalf("expected non-nil inspect response and tool")
	}
	if resp.Tool.ID != "tool-text-formatter" {
		t.Errorf("expected tool ID 'tool-text-formatter', got %q", resp.Tool.ID)
	}
	if resp.Tool.FunctionalState != ToolStateDraft {
		t.Errorf("expected state 'draft', got %q", resp.Tool.FunctionalState)
	}
	if resp.IsExecutable {
		t.Errorf("expected IsExecutable to be false under REQ-PROD-006")
	}
	if resp.CorrelationID == "" {
		t.Errorf("expected non-empty correlation ID")
	}
	if !resp.RuntimeAvailability.IsAvailable {
		t.Errorf("expected runtime availability true")
	}
}

func TestTool_Inspect_NotFound(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	resp, err := svc.InspectTool(ctx, &auth, "tenant-alpha", "tool-non-existent")
	if err == nil {
		t.Fatalf("expected 404 error for non-existent tool, got nil response: %v", resp)
	}

	errEnv, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Errorf("expected 404 asset-not-found, got HTTP %d, code %s", errEnv.HTTPStatus, errEnv.Code)
	}
}

func TestTool_Inspect_CrossTenant_AntiOracle404(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	authAlpha := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	// Attempt to inspect tenant-bravo tool scoped to tenant-alpha
	// Must return identical 404 asset-not-found as a non-existent asset (never 403 or existence hint)
	resp, err := svc.InspectTool(ctx, &authAlpha, "tenant-alpha", "tool-bravo-calculator")
	if err == nil {
		t.Fatalf("expected error for cross-tenant inspect, got: %v", resp)
	}

	errEnv, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if errEnv.HTTPStatus != 404 || errEnv.Code != ErrorCodeAssetNotFound {
		t.Errorf("expected uniform 404 asset-not-found, got HTTP %d, code %s", errEnv.HTTPStatus, errEnv.Code)
	}
}

func TestTool_Inspect_MissingOrWildcardTenant(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	// Missing tenant
	_, err := svc.InspectTool(ctx, &auth, "", "tool-text-formatter")
	if err == nil {
		t.Fatalf("expected error on empty tenantRef")
	}
	if env, ok := err.(*ErrorEnvelope); !ok || env.Code != ErrorCodeMissingTenant {
		t.Errorf("expected missing-tenant code, got %v", err)
	}

	// Wildcard tenants (*, %, all)
	wildcards := []string{"*", "%", "ALL", "all"}
	for _, wc := range wildcards {
		_, err := svc.InspectTool(ctx, &auth, wc, "tool-text-formatter")
		if err == nil {
			t.Errorf("expected error on wildcard tenant %q, got nil", wc)
		}
		if env, ok := err.(*ErrorEnvelope); !ok || env.Code != ErrorCodeWildcardTenant {
			t.Errorf("expected wildcard-tenant code for %q, got %v", wc, err)
		}
	}
}

func TestTool_Inspect_MSSP_AuthorizedTenantSet(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// MSSP operator with ATS = ["tenant-alpha", "tenant-bravo"]
	authMSSP := NewAuthContext("mssp-op", "tenant-alpha", []string{"tenant-alpha", "tenant-bravo"}, PermStudioLibraryRead)

	// 1. Can inspect tool in tenant-alpha
	respA, errA := svc.InspectTool(ctx, &authMSSP, "tenant-alpha", "tool-text-formatter")
	if errA != nil || respA == nil {
		t.Fatalf("expected MSSP access to tenant-alpha to succeed, got %v", errA)
	}

	// 2. Can inspect tool in tenant-bravo
	respB, errB := svc.InspectTool(ctx, &authMSSP, "tenant-bravo", "tool-bravo-calculator")
	if errB != nil || respB == nil {
		t.Fatalf("expected MSSP access to tenant-bravo to succeed, got %v", errB)
	}

	// 3. Rejected when attempting to inspect unauthorized tenant-charlie (not in ATS)
	_, errC := svc.InspectTool(ctx, &authMSSP, "tenant-charlie", "some-tool")
	if errC == nil {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got nil")
	}
	if env, ok := errC.(*ErrorEnvelope); !ok || env.HTTPStatus != 403 || env.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Errorf("expected 403 unauthorized-tenant-in-set, got %v", errC)
	}
}

func TestTool_Inspect_PermissionDenied(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Caller without perm.studio.library.read
	authNoPerm := NewAuthContext("unauthorized-user", "tenant-alpha", nil, "perm.some.unrelated.perm")

	_, err := svc.InspectTool(ctx, &authNoPerm, "tenant-alpha", "tool-text-formatter")
	if err == nil {
		t.Fatalf("expected permission denied error, got nil")
	}

	if env, ok := err.(*ErrorEnvelope); !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Errorf("expected 403 permission-denied, got %v", err)
	}
}

func TestTool_Inspect_DecoupledRuntimeAvailability(t *testing.T) {
	// CAP-STD-003 §25 Given 1:
	// A Tool is defined but its underlying runtime is absent.
	// When inspected, definition remains 100% valid and availability is reported separately.
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	resp, err := svc.InspectTool(ctx, &auth, "tenant-alpha", "tool-offline-validator")
	if err != nil {
		t.Fatalf("expected inspection to succeed even if runtime is offline, got %v", err)
	}

	if resp.Tool == nil || resp.Tool.ID != "tool-offline-validator" {
		t.Fatalf("expected valid tool definition returned")
	}
	if resp.RuntimeAvailability.IsAvailable {
		t.Errorf("expected runtime availability to be false")
	}
	if resp.RuntimeAvailability.HealthStatus != "unavailable" {
		t.Errorf("expected health status 'unavailable', got %q", resp.RuntimeAvailability.HealthStatus)
	}
}

// =============================================================================
// Task 3.2: Class 1 ValidateToolContract Tests (REQ-PROD-016, REQ-AI-002)
// =============================================================================

func TestTool_ValidateContract_NominalValid(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("designer-1", "tenant-alpha", nil, PermStudioLibraryRead)

	req := &ValidateToolContractRequest{
		TenantRef: "tenant-alpha",
		ToolID:    "tool-text-formatter",
	}

	res, err := svc.ValidateToolContract(ctx, &auth, req)
	if err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}

	if !res.Valid {
		t.Errorf("expected contract to be valid, issues: %+v", res.Issues)
	}
	if !res.DeterministicValidated {
		t.Errorf("expected DeterministicValidated to be true")
	}
	if !res.NonAgentConfirmed {
		t.Errorf("expected NonAgentConfirmed to be true")
	}
	if res.CompletenessScore < 0.8 {
		t.Errorf("expected completeness score >= 0.8, got %f", res.CompletenessScore)
	}
}

func TestTool_ValidateContract_WithoutAI_AlwaysFunctional(t *testing.T) {
	// REQ-PROD-016: Deterministic validation works without AI models or networks
	candidate := &ToolDefinition{
		ID:                  "tool-calculator",
		TenantID:            "tenant-alpha",
		Name:                "Basic Calculator",
		Purpose:             "Perform arithmetic calculations deterministically",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"add", "subtract"},
		Operations: []ToolOperation{
			{
				OperationID: "op-add",
				Name:        "add",
				Parameters: []Parameter{
					{Name: "a", Type: "number", Required: true},
					{Name: "b", Type: "number", Required: true},
				},
			},
		},
		ActionClass:        ActionClass0Observation,
		RiskClassification: RiskClassificationLow,
		SideEffects:        SideEffectContext{HasSideEffects: false, IsReversible: true},
	}

	res := DeterministicValidateToolContract(candidate, "corr-test-deterministic")
	if !res.Valid {
		t.Fatalf("expected 100%% deterministic validation to succeed, issues: %+v", res.Issues)
	}
	if !res.DeterministicValidated {
		t.Errorf("expected deterministic validated flag true")
	}
	if res.CompletenessScore < 0.8 {
		t.Errorf("expected completeness score >= 0.8, got %f", res.CompletenessScore)
	}
}

func TestTool_ValidateContract_AntiAgentBoundaryDefense(t *testing.T) {
	// CAP-STD-003 §25 Given 2 & REQ-AI-002:
	// An AI-capable tool must never morph into an Agent.
	// Contract validation rejects agent-exclusive attributes.

	// 1. Nominal AI-capable tool (valid, non-agent)
	validAITool := &ToolDefinition{
		ID:                  "tool-summarizer-ai",
		TenantID:            "tenant-alpha",
		Name:                "AI Text Summarizer",
		Purpose:             "Summarize document text using model provider",
		ExecutionNature:     ExecutionNatureAICapable,
		SupportedOperations: []string{"summarize"},
		Operations: []ToolOperation{
			{
				OperationID: "op-sum",
				Name:        "summarize",
				Parameters:  []Parameter{{Name: "text", Type: "string", Required: true}},
			},
		},
		ActionClass:        ActionClass0Observation,
		RiskClassification: RiskClassificationLow,
		SideEffects:        SideEffectContext{HasSideEffects: false, IsReversible: true},
	}

	resValid := DeterministicValidateToolContract(validAITool, "corr-ai-valid")
	if !resValid.Valid {
		t.Fatalf("expected valid AI-capable tool to pass, issues: %+v", resValid.Issues)
	}
	if !resValid.IsAICapable {
		t.Errorf("expected IsAICapable = true")
	}
	if !resValid.NonAgentConfirmed {
		t.Errorf("expected NonAgentConfirmed = true")
	}

	// 2. Tool attempting to declare autonomous loop (must be rejected)
	agentMorphingTool := validAITool.Clone()
	agentMorphingTool.Purpose = "Defines an autonomous_loop for background coordination"

	resViolating := DeterministicValidateToolContract(agentMorphingTool, "corr-ai-invalid")
	if resViolating.Valid {
		t.Fatalf("expected agent-morphing tool to be rejected by anti-agent defense")
	}
	if resViolating.NonAgentConfirmed {
		t.Errorf("expected NonAgentConfirmed = false")
	}

	foundViolationIssue := false
	for _, issue := range resViolating.Issues {
		if issue.Code == "anti-agent-boundary-violation" {
			foundViolationIssue = true
			break
		}
	}
	if !foundViolationIssue {
		t.Errorf("expected anti-agent-boundary-violation issue code, got %+v", resViolating.Issues)
	}
}

func TestTool_ValidateContract_MissingMandatoryFields(t *testing.T) {
	badTool := &ToolDefinition{
		ID:      "", // Missing ID
		Name:    "", // Missing Name
		Purpose: "", // Missing Purpose
	}

	res := DeterministicValidateToolContract(badTool, "corr-bad")
	if res.Valid {
		t.Fatalf("expected invalid result for tool missing ID, Name, Purpose")
	}

	hasMissingID := false
	hasMissingName := false
	hasMissingPurpose := false

	for _, issue := range res.Issues {
		switch issue.Code {
		case "missing-tool-id":
			hasMissingID = true
		case "missing-tool-name":
			hasMissingName = true
		case "missing-tool-purpose":
			hasMissingPurpose = true
		}
	}

	if !hasMissingID || !hasMissingName || !hasMissingPurpose {
		t.Errorf("expected missing field issues, got %+v", res.Issues)
	}
}

func TestTool_ValidateContract_InvalidExecutionNature(t *testing.T) {
	tool := &ToolDefinition{
		ID:                  "tool-custom",
		Name:                "Custom Tool",
		Purpose:             "Some purpose",
		ExecutionNature:     "quantum-ai", // Invalid!
		SupportedOperations: []string{"run"},
	}

	res := DeterministicValidateToolContract(tool, "corr-invalid-exec")
	if res.Valid {
		t.Fatalf("expected validation failure for invalid execution nature")
	}

	found := false
	for _, issue := range res.Issues {
		if issue.Code == "invalid-execution-nature" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected invalid-execution-nature issue, got %+v", res.Issues)
	}
}

func TestTool_ValidateContract_EmptyOperations_Rejected(t *testing.T) {
	tool := &ToolDefinition{
		ID:                  "tool-no-ops",
		Name:                "Empty Tool",
		Purpose:             "Test tool without operations",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: nil, // No operations!
	}

	res := DeterministicValidateToolContract(tool, "corr-no-ops")
	if res.Valid {
		t.Fatalf("expected failure for tool with zero operations")
	}

	found := false
	for _, issue := range res.Issues {
		if issue.Code == "unsupported-operation" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unsupported-operation issue, got %+v", res.Issues)
	}
}

// =============================================================================
// Task 3.3: Class 2 CreateDraftTool and UpdateDraftTool Tests
// =============================================================================

func TestTool_CreateDraft_Nominal(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("designer-1", "tenant-alpha", nil, PermStudioBuilderRead, PermStudioArtifactManage)

	req := &CreateDraftToolRequest{
		TenantRef:           "tenant-alpha",
		ID:                  "tool-created-draft",
		Name:                "Timestamp Normalizer",
		Purpose:             "Normalize RFC3339 timestamps to UTC format",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"normalize"},
		ActionClass:         ActionClass0Observation,
		RiskClassification:  RiskClassificationLow,
		SideEffects:         SideEffectContext{HasSideEffects: false, IsReversible: true},
		Owner:               Owner{OwnerID: "designer-1", OwnerName: "Designer", OwnerType: "user"},
		Reason:              "Registering new timestamp normalization capability",
	}

	created, err := svc.CreateDraftTool(ctx, &auth, req)
	if err != nil {
		t.Fatalf("expected CreateDraftTool to succeed, got %v", err)
	}

	if created == nil {
		t.Fatalf("expected non-nil tool returned")
	}
	if created.ID != "tool-created-draft" {
		t.Errorf("expected ID 'tool-created-draft', got %q", created.ID)
	}
	if created.FunctionalState != ToolStateDraft {
		t.Errorf("expected initial state 'draft', got %q", created.FunctionalState)
	}
	if len(created.AuditTrail) == 0 {
		t.Errorf("expected initial audit trail entry")
	}
	if created.Owner.IsUnknown {
		t.Errorf("expected normalized owner")
	}
}

func TestTool_CreateDraft_PermissionDenied(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Caller without library read or builder read
	auth := NewAuthContext("reader-only", "tenant-alpha", nil, "perm.studio.workflow.read")

	req := &CreateDraftToolRequest{
		TenantRef:       "tenant-alpha",
		Name:            "Disallowed Tool",
		Purpose:         "Should fail",
		ExecutionNature: ExecutionNatureDeterministic,
	}

	_, err := svc.CreateDraftTool(ctx, &auth, req)
	if err == nil {
		t.Fatalf("expected permission denied for CreateDraftTool, got nil")
	}
	if env, ok := err.(*ErrorEnvelope); !ok || env.HTTPStatus != 403 || env.Code != ErrorCodePermissionDenied {
		t.Errorf("expected 403 permission-denied, got %v", err)
	}
}

func TestTool_CreateDraft_CrossTenantMutation_Forbidden(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Caller is scoped to tenant-alpha with builder read and artifact manage, with ATS including both tenants
	auth := NewAuthContext("builder-alpha", "tenant-alpha", []string{"tenant-alpha", "tenant-bravo"}, PermStudioBuilderRead, PermStudioArtifactManage)

	// Attempt to create draft tool in tenant-bravo
	req := &CreateDraftToolRequest{
		TenantRef:       "tenant-bravo", // Foreign tenant!
		Name:            "Cross Tenant Tool",
		Purpose:         "Cross tenant attempt",
		ExecutionNature: ExecutionNatureDeterministic,
	}

	_, err := svc.CreateDraftTool(ctx, &auth, req)
	if err == nil {
		t.Fatalf("expected cross-tenant draft creation to fail, got nil")
	}
	if env, ok := err.(*ErrorEnvelope); !ok || env.HTTPStatus != 403 || env.Code != ErrorCodeCrossTenantMutationForbidden {
		t.Errorf("expected 403 cross-tenant-mutation-forbidden, got %v", err)
	}
}

func TestTool_UpdateDraft_Nominal(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("builder-1", "tenant-alpha", nil, PermStudioBuilderRead, PermStudioArtifactManage)

	req := &UpdateDraftToolRequest{
		TenantRef:   "tenant-alpha",
		ToolID:      "tool-text-formatter",
		Description: "Updated description for Text Formatter Tool with enhanced padding",
		Limitations: []string{"Maximum input length 4096 characters"},
		Reason:      "Increased limit after server upgrade",
	}

	updated, err := svc.UpdateDraftTool(ctx, &auth, req)
	if err != nil {
		t.Fatalf("expected UpdateDraftTool to succeed, got %v", err)
	}

	if updated == nil {
		t.Fatalf("expected non-nil updated tool")
	}
	if updated.Description != req.Description {
		t.Errorf("expected description updated, got %q", updated.Description)
	}
	if len(updated.Limitations) != 1 || updated.Limitations[0] != req.Limitations[0] {
		t.Errorf("expected limitations updated")
	}
	if len(updated.AuditTrail) < 2 {
		t.Errorf("expected audit trail entry added for update")
	}
}

func TestTool_UpdateDraft_Immutability_NonDraftRejected(t *testing.T) {
	// CAP-STD-003 §10, §14:
	// Only tools in state 'draft' can be updated.
	// Updating tool in 'available-reference' must return HTTP 400 immutable-asset-mutation.
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("builder-1", "tenant-alpha", nil, PermStudioBuilderRead, PermStudioArtifactManage)

	req := &UpdateDraftToolRequest{
		TenantRef:   "tenant-alpha",
		ToolID:      "tool-currency-converter", // State: available-reference
		Description: "Illegal attempt to mutate immutable reference tool",
	}

	_, err := svc.UpdateDraftTool(ctx, &auth, req)
	if err == nil {
		t.Fatalf("expected immutable-asset-mutation error for non-draft tool, got nil")
	}

	env, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if env.HTTPStatus != 400 || env.Code != ErrorCodeImmutableAssetMutation {
		t.Errorf("expected 400 immutable-asset-mutation, got HTTP %d, code %s", env.HTTPStatus, env.Code)
	}
}

func TestTool_UpdateDraft_ForbiddenStates_PublishedDeployed(t *testing.T) {
	// CAP-STD-003 §14:
	// States 'published' and 'deployed' are strictly forbidden under STD-1.
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("builder-1", "tenant-alpha", nil, PermStudioBuilderRead, PermStudioArtifactManage)

	forbiddenStates := []ToolFunctionalState{"published", "deployed"}
	for _, st := range forbiddenStates {
		req := &UpdateDraftToolRequest{
			TenantRef:       "tenant-alpha",
			ToolID:          "tool-text-formatter",
			FunctionalState: st,
		}

		_, err := svc.UpdateDraftTool(ctx, &auth, req)
		if err == nil {
			t.Errorf("expected transition to forbidden state %q to fail, got nil", st)
		}
		if env, ok := err.(*ErrorEnvelope); !ok || env.HTTPStatus != 400 || env.Code != ErrorCodeInvalidMetadata {
			t.Errorf("expected 400 invalid-metadata for state %q, got %v", st, err)
		}
	}
}

// =============================================================================
// Task 3.4: REQ-PROD-006 Execution Verb Blocking Tests
// =============================================================================

func TestTool_ExecutionVerbs_PermanentlyBlocked(t *testing.T) {
	// REQ-PROD-006: Verbs Invoke, Execute, Run, Call, Trigger are permanently blocked
	// in Studio Tool Definition and return HTTP 400 with code "execution-blocked".
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("admin-user", "tenant-alpha", nil, PermStudioLibraryRead, PermStudioBuilderRead)

	type execCall struct {
		name string
		fn   func() error
	}

	calls := []execCall{
		{"Invoke", func() error { return svc.Invoke(ctx, &auth, "tool-text-formatter", nil) }},
		{"Execute", func() error { return svc.Execute(ctx, &auth, "tool-text-formatter", nil) }},
		{"Run", func() error { return svc.Run(ctx, &auth, "tool-text-formatter", nil) }},
		{"Call", func() error { return svc.Call(ctx, &auth, "tool-text-formatter", nil) }},
		{"Trigger", func() error { return svc.Trigger(ctx, &auth, "tool-text-formatter", nil) }},
	}

	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			if err == nil {
				t.Fatalf("expected %s to be permanently blocked, got nil", tc.name)
			}

			env, ok := err.(*ErrorEnvelope)
			if !ok {
				t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
			}
			if env.HTTPStatus != 400 {
				t.Errorf("expected HTTP 400 for %s, got %d", tc.name, env.HTTPStatus)
			}
			if env.Code != ErrorCodeExecutionBlocked {
				t.Errorf("expected code 'execution-blocked' for %s, got %s", tc.name, env.Code)
			}
			if !strings.Contains(env.Message, "REQ-PROD-006") {
				t.Errorf("expected citation of REQ-PROD-006 in error message, got: %s", env.Message)
			}
			if env.Retryable {
				t.Errorf("expected retryable = false for blocked execution verb")
			}
		})
	}
}

// =============================================================================
// Task 3.5: REQ-OBJ-009 Owner Unknown Normalization Tests
// =============================================================================

func TestTool_OwnerUnknown_Normalization(t *testing.T) {
	// REQ-OBJ-009: Missing, empty, or whitespace owner normalizes to 'owner unknown', is_unknown: true
	tests := []struct {
		name     string
		input    Owner
		expected Owner
	}{
		{
			name:  "empty struct",
			input: Owner{},
			expected: Owner{
				OwnerID:     "unknown",
				OwnerName:   "owner unknown",
				OwnerType:   "unknown",
				IsUnknown:   true,
				DisplayName: "owner unknown",
			},
		},
		{
			name:  "whitespace only",
			input: Owner{OwnerID: "   ", OwnerName: "   "},
			expected: Owner{
				OwnerID:     "unknown",
				OwnerName:   "owner unknown",
				OwnerType:   "unknown",
				IsUnknown:   true,
				DisplayName: "owner unknown",
			},
		},
		{
			name:  "explicit is_unknown flag",
			input: Owner{OwnerID: "user-123", OwnerName: "Some User", IsUnknown: true},
			expected: Owner{
				OwnerID:     "unknown",
				OwnerName:   "owner unknown",
				OwnerType:   "unknown",
				IsUnknown:   true,
				DisplayName: "owner unknown",
			},
		},
		{
			name:  "valid owner without display name defaults to owner name",
			input: Owner{OwnerID: "team-dev", OwnerName: "Dev Team", OwnerType: "team"},
			expected: Owner{
				OwnerID:     "team-dev",
				OwnerName:   "Dev Team",
				OwnerType:   "team",
				IsUnknown:   false,
				DisplayName: "Dev Team",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := tc.input
			o.Normalize()
			if o.OwnerID != tc.expected.OwnerID || o.OwnerName != tc.expected.OwnerName ||
				o.IsUnknown != tc.expected.IsUnknown || o.DisplayName != tc.expected.DisplayName {
				t.Errorf("normalization mismatch: got %+v, want %+v", o, tc.expected)
			}
		})
	}
}

// =============================================================================
// Task 3.6: Operational & Utility Method Coverage Tests
// =============================================================================

func TestTool_ToAsset_Conversion(t *testing.T) {
	svc, store := setupTestToolService(t)
	ctx := context.Background()

	tool, err := store.Get(ctx, "tenant-alpha", "tool-text-formatter")
	if err != nil || tool == nil {
		t.Fatalf("failed to fetch seeded tool: %v", err)
	}

	asset := tool.ToAsset()
	if asset == nil {
		t.Fatalf("expected non-nil asset projection")
	}
	if asset.Type != AssetKindTool {
		t.Errorf("expected asset kind 'tool', got %q", asset.Type)
	}
	if asset.ID != tool.ID || asset.TenantID != tool.TenantID || asset.Name != tool.Name {
		t.Errorf("asset fields mismatch with tool")
	}

	// Verify SetRuntimeAvailability
	err = svc.SetRuntimeAvailability(ctx, "tenant-alpha", "tool-text-formatter", RuntimeAvailabilityStatus{
		IsAvailable:  false,
		HealthStatus: "degraded",
		CheckedAt:    time.Now().UTC(),
		Message:      "High latency warning",
	})
	if err != nil {
		t.Fatalf("SetRuntimeAvailability failed: %v", err)
	}

	updatedTool, _ := store.Get(ctx, "tenant-alpha", "tool-text-formatter")
	if updatedTool.RuntimeAvailability.IsAvailable {
		t.Errorf("expected runtime availability false after update")
	}
	if updatedTool.RuntimeAvailability.HealthStatus != "degraded" {
		t.Errorf("expected health status 'degraded'")
	}
}

func TestTool_StoreOperations(t *testing.T) {
	store := NewInMemoryToolStore()
	ctx := context.Background()

	tool := &ToolDefinition{
		ID:       "tool-store-test",
		TenantID: "tenant-store",
		Name:     "Store Test Tool",
		Purpose:  "Testing store CRUD",
	}

	// Put
	if err := store.Put(ctx, tool); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// List
	list, err := store.List(ctx, "tenant-store")
	if err != nil || len(list) != 1 {
		t.Fatalf("List failed: expected 1 item, got %d (err: %v)", len(list), err)
	}

	// Get
	got, err := store.Get(ctx, "tenant-store", "tool-store-test")
	if err != nil || got == nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name != tool.Name {
		t.Errorf("name mismatch: got %q, want %q", got.Name, tool.Name)
	}

	// Delete
	if err := store.Delete(ctx, "tenant-store", "tool-store-test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	gotAfter, _ := store.Get(ctx, "tenant-store", "tool-store-test")
	if gotAfter != nil {
		t.Errorf("expected nil after delete, got %+v", gotAfter)
	}
}
