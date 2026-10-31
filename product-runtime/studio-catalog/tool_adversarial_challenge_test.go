package studiocatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// CHALLENGE 1: REQ-PROD-006 Permanent Execution Blocking Adversarial Penetration
// =============================================================================

func TestChallenge_ExecutionVerbs_AdversarialPenetrationMatrix(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	verbs := []struct {
		name string
		call func(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	}{
		{"Invoke", svc.Invoke},
		{"Execute", svc.Execute},
		{"Run", svc.Run},
		{"Call", svc.Call},
		{"Trigger", svc.Trigger},
	}

	payloads := []struct {
		desc    string
		payload any
	}{
		{"nil payload", nil},
		{"empty string", ""},
		{"shell injection rm -rf", "; rm -rf / ; cat /etc/passwd"},
		{"nested map with command", map[string]any{"cmd": "powershell.exe -c calc", "args": []string{"-nop", "-w", "hidden"}}},
		{"massive payload (100KB)", strings.Repeat("A", 102400)},
		{"raw JSON", json.RawMessage(`{"exploit": true, "eval": "process.exit(1)"}`)},
		{"channel payload", make(chan int)},
		{"func payload", func() { panic("exploit") }},
	}

	auths := []struct {
		desc string
		auth *AuthContext
	}{
		{"nil auth context", nil},
		{"anonymous caller", &AuthContext{PrincipalID: "anon", Permissions: map[string]bool{}}},
		{"standard reader", &AuthContext{PrincipalID: "user-reader", TenantID: "tenant-alpha", Permissions: map[string]bool{PermStudioLibraryRead: true}}},
		{"builder author", &AuthContext{PrincipalID: "user-builder", TenantID: "tenant-alpha", Permissions: map[string]bool{PermStudioBuilderRead: true}}},
		{"superadmin wildcard", &AuthContext{PrincipalID: "root", TenantID: "tenant-alpha", Permissions: map[string]bool{"*": true}}},
		{"cross-tenant mssp", &AuthContext{PrincipalID: "mssp-admin", TenantID: "tenant-alpha", AuthorizedTenantSet: []string{"tenant-alpha", "tenant-bravo"}, Permissions: map[string]bool{PermStudioLibraryRead: true, PermStudioBuilderRead: true}}},
	}

	toolIDs := []string{
		"tool-text-formatter",
		"tool-currency-converter",
		"tool-offline-validator",
		"tool-bravo-calculator",
		"non-existent-tool-id",
		"",
		"   ",
		"../../etc/passwd",
		"' OR '1'='1",
	}

	for _, v := range verbs {
		for _, p := range payloads {
			for _, a := range auths {
				for _, tid := range toolIDs {
					err := v.call(ctx, a.auth, tid, p.payload)
					if err == nil {
						t.Fatalf("CRITICAL SECURITY VULNERABILITY: Verb %s succeeded with payload %s, auth %s, toolID %q!",
							v.name, p.desc, a.desc, tid)
					}

					env, ok := err.(*ErrorEnvelope)
					if !ok {
						t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
					}
					if env.HTTPStatus != 400 {
						t.Errorf("verb %s: expected HTTP 400, got %d", v.name, env.HTTPStatus)
					}
					if env.Code != ErrorCodeExecutionBlocked {
						t.Errorf("verb %s: expected error code %s, got %s", v.name, ErrorCodeExecutionBlocked, env.Code)
					}
					if env.Retryable {
						t.Errorf("verb %s: expected Retryable=false", v.name)
					}
					if !strings.Contains(env.Message, "REQ-PROD-006") {
						t.Errorf("verb %s: expected REQ-PROD-006 in message, got %q", v.name, env.Message)
					}
					if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(v.name)) {
						t.Errorf("verb %s: expected verb name in message, got %q", v.name, env.Message)
					}
				}
			}
		}
	}
}

// =============================================================================
// CHALLENGE 2: Cross-Tenant Existence Oracle & Isolation Invariant
// =============================================================================

func TestChallenge_CrossTenant_ExistenceOracleIndistinguishability(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Caller is strictly in tenant-alpha
	authAlpha := NewAuthContext("analyst-1", "tenant-alpha", []string{"tenant-alpha"}, PermStudioLibraryRead)

	// Probe A: Query known tool belonging to tenant-bravo from tenant-alpha
	respA, errA := svc.InspectTool(ctx, &authAlpha, "tenant-alpha", "tool-bravo-calculator")
	if respA != nil || errA == nil {
		t.Fatalf("expected cross-tenant inspect to fail, got resp=%+v err=%v", respA, errA)
	}

	// Probe B: Query completely non-existent tool ID from tenant-alpha
	respB, errB := svc.InspectTool(ctx, &authAlpha, "tenant-alpha", "ghost-tool-id-nonexistent-12345")
	if respB != nil || errB == nil {
		t.Fatalf("expected non-existent inspect to fail, got resp=%+v err=%v", respB, errB)
	}

	envA, okA := errA.(*ErrorEnvelope)
	envB, okB := errB.(*ErrorEnvelope)
	if !okA || !okB {
		t.Fatalf("expected *ErrorEnvelope for both, got %T and %T", errA, errB)
	}

	// Anti-oracle validation:
	// 1. Both HTTP status codes must be identical (404)
	if envA.HTTPStatus != 404 || envB.HTTPStatus != 404 {
		t.Errorf("HTTP status mismatch: A=%d, B=%d (both must be 404)", envA.HTTPStatus, envB.HTTPStatus)
	}
	// 2. Both error codes must be identical ("asset-not-found")
	if envA.Code != ErrorCodeAssetNotFound || envB.Code != ErrorCodeAssetNotFound {
		t.Errorf("Error code mismatch: A=%s, B=%s (both must be %s)", envA.Code, envB.Code, ErrorCodeAssetNotFound)
	}
	// 3. Retryability must be identical (false)
	if envA.Retryable != envB.Retryable || envA.Retryable != false {
		t.Errorf("Retryable mismatch: A=%v, B=%v (must be false)", envA.Retryable, envB.Retryable)
	}
	// 4. Message format must follow the exact same template without leaking existence in another tenant
	expectedMsgA := `Asset "tool-bravo-calculator" was not found in tenant "tenant-alpha"`
	expectedMsgB := `Asset "ghost-tool-id-nonexistent-12345" was not found in tenant "tenant-alpha"`
	if envA.Message != expectedMsgA {
		t.Errorf("leak in message A: expected %q, got %q", expectedMsgA, envA.Message)
	}
	if envB.Message != expectedMsgB {
		t.Errorf("leak in message B: expected %q, got %q", expectedMsgB, envB.Message)
	}
}

func TestChallenge_CrossTenant_DirectScopeAccess_WithoutATS(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// Caller is in tenant-alpha, ATS contains only tenant-alpha
	authAlpha := NewAuthContext("user-alpha", "tenant-alpha", []string{"tenant-alpha"}, PermStudioLibraryRead)

	// Caller attempts to access tenant-bravo scope directly
	_, err := svc.InspectTool(ctx, &authAlpha, "tenant-bravo", "tool-bravo-calculator")
	if err == nil {
		t.Fatalf("expected 403 unauthorized-tenant-in-set, got nil")
	}

	env, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if env.HTTPStatus != 403 || env.Code != ErrorCodeUnauthorizedTenantInSet {
		t.Errorf("expected 403 unauthorized-tenant-in-set, got HTTP %d, code %s", env.HTTPStatus, env.Code)
	}
}

// =============================================================================
// CHALLENGE 3: Cross-Tenant Mutation via MSSP ATS Session
// =============================================================================

func TestChallenge_MSSP_CrossTenantMutation_CreateDraft_Forbidden(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	// MSSP operator with read & builder permissions across multiple tenants in ATS
	// Primary tenant: tenant-alpha. ATS: ["tenant-alpha", "tenant-bravo"]
	authMSSP := NewAuthContext(
		"mssp-lead",
		"tenant-alpha",
		[]string{"tenant-alpha", "tenant-bravo"},
		PermStudioLibraryRead,
		PermStudioBuilderRead,
	)

	// Attempt 1: Create draft tool in foreign tenant-bravo
	reqCreate := &CreateDraftToolRequest{
		TenantRef:           "tenant-bravo", // Target is tenant-bravo!
		ID:                  "tool-mssp-injected",
		Name:                "Injected Tool",
		Purpose:             "Attempt to create tool in tenant-bravo from tenant-alpha MSSP session",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"noop"},
	}

	created, err := svc.CreateDraftTool(ctx, &authMSSP, reqCreate)
	if err == nil || created != nil {
		t.Fatalf("CRITICAL SECURITY VIOLATION: MSSP cross-tenant draft creation succeeded! created=%+v", created)
	}

	env, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if env.HTTPStatus != 403 || env.Code != ErrorCodeCrossTenantMutationForbidden {
		t.Errorf("expected 403 cross-tenant-mutation-forbidden, got HTTP %d, code %s", env.HTTPStatus, env.Code)
	}
}

func TestChallenge_MSSP_CrossTenantMutation_UpdateDraft_Forbidden(t *testing.T) {
	svc, store := setupTestToolService(t)
	ctx := context.Background()

	// Seed a draft tool in tenant-bravo
	bravoDraft := &ToolDefinition{
		ID:                  "tool-bravo-draft",
		TenantID:            "tenant-bravo",
		Name:                "Bravo Draft Tool",
		Purpose:             "Draft tool in tenant-bravo",
		FunctionalState:     ToolStateDraft,
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"calculate"},
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}
	if err := store.Put(ctx, bravoDraft); err != nil {
		t.Fatalf("failed to seed bravoDraft: %v", err)
	}

	// MSSP operator with primary tenant-alpha and ATS ["tenant-alpha", "tenant-bravo"]
	authMSSP := NewAuthContext(
		"mssp-lead",
		"tenant-alpha",
		[]string{"tenant-alpha", "tenant-bravo"},
		PermStudioLibraryRead,
		PermStudioBuilderRead,
	)

	// Attempt 2: Update draft tool in foreign tenant-bravo
	reqUpdate := &UpdateDraftToolRequest{
		TenantRef:   "tenant-bravo",
		ToolID:      "tool-bravo-draft",
		Description: "Malicious update from tenant-alpha session",
	}

	updated, err := svc.UpdateDraftTool(ctx, &authMSSP, reqUpdate)
	if err == nil || updated != nil {
		t.Fatalf("CRITICAL SECURITY VIOLATION: MSSP cross-tenant draft update succeeded! updated=%+v", updated)
	}

	env, ok := err.(*ErrorEnvelope)
	if !ok {
		t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
	}
	if env.HTTPStatus != 403 || env.Code != ErrorCodeCrossTenantMutationForbidden {
		t.Errorf("expected 403 cross-tenant-mutation-forbidden, got HTTP %d, code %s", env.HTTPStatus, env.Code)
	}

	// Verify tool in tenant-bravo was NOT modified
	toolAfter, _ := store.Get(ctx, "tenant-bravo", "tool-bravo-draft")
	if toolAfter.Description == "Malicious update from tenant-alpha session" {
		t.Fatalf("CRITICAL: store was mutated despite error return!")
	}
}

// =============================================================================
// CHALLENGE 4: Immutability of All Non-Draft Tools
// =============================================================================

func TestChallenge_Immutability_AllNonDraftStates_Rejected(t *testing.T) {
	svc, store := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("builder-1", "tenant-alpha", []string{"tenant-alpha"}, PermStudioBuilderRead)

	nonDraftStates := []struct {
		state  ToolFunctionalState
		toolID string
	}{
		{ToolStateAvailableReference, "tool-state-avail-ref"},
		{ToolStateDeprecated, "tool-state-deprecated"},
		{ToolStateDisabled, "tool-state-disabled"},
		{ToolStateSuperseded, "tool-state-superseded"},
		{ToolStateReviewReady, "tool-state-review-ready"},
	}

	for _, tc := range nonDraftStates {
		t.Run("State_"+string(tc.state), func(t *testing.T) {
			tool := &ToolDefinition{
				ID:                  tc.toolID,
				TenantID:            "tenant-alpha",
				Name:                "Immutable Tool " + string(tc.state),
				Purpose:             "Testing immutability for state " + string(tc.state),
				FunctionalState:     tc.state,
				ExecutionNature:     ExecutionNatureDeterministic,
				SupportedOperations: []string{"test"},
				CreatedAt:           time.Now().UTC(),
				UpdatedAt:           time.Now().UTC(),
			}
			if err := store.Put(ctx, tool); err != nil {
				t.Fatalf("failed to seed tool: %v", err)
			}

			// Attempt to update
			req := &UpdateDraftToolRequest{
				TenantRef:   "tenant-alpha",
				ToolID:      tc.toolID,
				Description: "Tampered description",
				Purpose:     "Tampered purpose",
			}

			updated, err := svc.UpdateDraftTool(ctx, &auth, req)
			if err == nil || updated != nil {
				t.Fatalf("IMMUTABILITY VIOLATION: Tool in state %q was updated!", tc.state)
			}

			env, ok := err.(*ErrorEnvelope)
			if !ok {
				t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
			}
			if env.HTTPStatus != 400 {
				t.Errorf("state %s: expected HTTP 400, got %d", tc.state, env.HTTPStatus)
			}
			if env.Code != ErrorCodeImmutableAssetMutation {
				t.Errorf("state %s: expected error code %s, got %s", tc.state, ErrorCodeImmutableAssetMutation, env.Code)
			}

			// Verify tool remains intact
			toolAfter, _ := store.Get(ctx, "tenant-alpha", tc.toolID)
			if toolAfter.Description == "Tampered description" {
				t.Fatalf("IMMUTABILITY VIOLATION: state %s tool was modified in store!", tc.state)
			}
		})
	}
}

// =============================================================================
// CHALLENGE 5: Invalid State Transitions (published, deployed)
// =============================================================================

func TestChallenge_StateTransitions_PublishedDeployedAndInvalid(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()
	auth := NewAuthContext("builder-1", "tenant-alpha", []string{"tenant-alpha"}, PermStudioBuilderRead)

	forbiddenTransitions := []struct {
		desc  string
		state ToolFunctionalState
	}{
		{"lowercase published", "published"},
		{"lowercase deployed", "deployed"},
		{"uppercase PUBLISHED", "PUBLISHED"},
		{"uppercase DEPLOYED", "DEPLOYED"},
		{"completely invalid state", "active-in-production"},
		{"sql injection state", "draft'; DROP TABLE tools; --"},
	}

	for _, tc := range forbiddenTransitions {
		t.Run(tc.desc, func(t *testing.T) {
			req := &UpdateDraftToolRequest{
				TenantRef:       "tenant-alpha",
				ToolID:          "tool-text-formatter",
				FunctionalState: tc.state,
			}

			updated, err := svc.UpdateDraftTool(ctx, &auth, req)
			if err == nil || updated != nil {
				t.Fatalf("STATE TRANSITION VIOLATION: Transitioned to forbidden state %q!", tc.state)
			}

			env, ok := err.(*ErrorEnvelope)
			if !ok {
				t.Fatalf("expected *ErrorEnvelope, got %T: %v", err, err)
			}
			if env.HTTPStatus != 400 {
				t.Errorf("state %q: expected HTTP 400, got %d", tc.state, env.HTTPStatus)
			}
			if env.Code != ErrorCodeInvalidMetadata {
				t.Errorf("state %q: expected error code %s, got %s", tc.state, ErrorCodeInvalidMetadata, env.Code)
			}
		})
	}
}

// =============================================================================
// CHALLENGE 6: Anti-Agent Boundary Defense Stress Testing (REQ-AI-002)
// =============================================================================

func TestChallenge_AntiAgentBoundaryDefense_EvasionAttempts(t *testing.T) {
	// Adversarial evasion attempts using capitalization, padding, and field injection
	evasionAttempts := []struct {
		desc           string
		mutateTool     func(t *ToolDefinition)
		expectRejected bool
	}{
		{
			desc: "planning_loop in mixed case in Name",
			mutateTool: func(tool *ToolDefinition) {
				tool.Name = "Tool With PlAnNiNg_LoOp Capability"
			},
			expectRejected: true,
		},
		{
			desc: "autonomous_loop in Description",
			mutateTool: func(tool *ToolDefinition) {
				tool.Description = "This tool executes an autonomous_loop in background"
			},
			expectRejected: true,
		},
		{
			desc: "agent_team in Limitations",
			mutateTool: func(tool *ToolDefinition) {
				tool.Limitations = []string{"Requires agent_team coordination"}
			},
			expectRejected: true,
		},
		{
			desc: "subagent_spawn in Purpose",
			mutateTool: func(tool *ToolDefinition) {
				tool.Purpose = "Provides subagent_spawn functionality"
			},
			expectRejected: true,
		},
		{
			desc: "goal_prompt in uppercase in Purpose",
			mutateTool: func(tool *ToolDefinition) {
				tool.Purpose = "Executes with GOAL_PROMPT directives"
			},
			expectRejected: true,
		},
		{
			desc: "Valid AI-capable tool with safe purpose",
			mutateTool: func(tool *ToolDefinition) {
				tool.ExecutionNature = ExecutionNatureAICapable
				tool.Purpose = "Extract entities from text documents using LLM model provider"
				tool.Description = "Single-pass extraction without agentic loop"
			},
			expectRejected: false,
		},
	}

	for _, tc := range evasionAttempts {
		t.Run(tc.desc, func(t *testing.T) {
			baseTool := &ToolDefinition{
				ID:                  "tool-anti-agent-test",
				TenantID:            "tenant-alpha",
				Name:                "Standard Test Tool",
				Purpose:             "Standard purpose statement",
				ExecutionNature:     ExecutionNatureAICapable,
				SupportedOperations: []string{"execute_op"},
				Operations: []ToolOperation{
					{
						OperationID: "op-1",
						Name:        "execute_op",
						Parameters:  []Parameter{{Name: "input", Type: "string"}},
					},
				},
				ActionClass:        ActionClass0Observation,
				RiskClassification: RiskClassificationLow,
			}

			tc.mutateTool(baseTool)
			res := DeterministicValidateToolContract(baseTool, "corr-anti-agent")

			if tc.expectRejected {
				if res.Valid {
					t.Fatalf("EVASION SUCCESS: %s bypassed anti-agent boundary defense!", tc.desc)
				}
				if res.NonAgentConfirmed {
					t.Errorf("expected NonAgentConfirmed=false for violating tool")
				}
			} else {
				if !res.Valid {
					t.Fatalf("FALSE POSITIVE: %s was rejected: %+v", tc.desc, res.Issues)
				}
				if !res.NonAgentConfirmed {
					t.Errorf("expected NonAgentConfirmed=true for safe AI tool")
				}
			}
		})
	}
}

// =============================================================================
// CHALLENGE 7: Concurrency & Thread-Safety Stress Harness
// =============================================================================

func TestChallenge_Concurrency_HighContentionHarness(t *testing.T) {
	svc, _ := setupTestToolService(t)
	ctx := context.Background()

	const numWorkers = 30
	const iterationsPerWorker = 20

	authAlpha := NewAuthContext("worker-alpha", "tenant-alpha", []string{"tenant-alpha"}, PermStudioLibraryRead, PermStudioBuilderRead)
	authMSSP := NewAuthContext("worker-mssp", "tenant-alpha", []string{"tenant-alpha", "tenant-bravo"}, PermStudioLibraryRead, PermStudioBuilderRead)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		workerID := w
		go func() {
			defer wg.Done()
			for i := 0; i < iterationsPerWorker; i++ {
				switch (workerID + i) % 5 {
				case 0:
					// Concurrent Inspect
					resp, err := svc.InspectTool(ctx, &authAlpha, "tenant-alpha", "tool-text-formatter")
					if err != nil || resp == nil {
						t.Errorf("worker %d inspect failed: %v", workerID, err)
					}
				case 1:
					// Concurrent Validate
					req := &ValidateToolContractRequest{
						TenantRef: "tenant-alpha",
						ToolID:    "tool-text-formatter",
					}
					vRes, err := svc.ValidateToolContract(ctx, &authAlpha, req)
					if err != nil || !vRes.Valid {
						t.Errorf("worker %d validate failed: %v", workerID, err)
					}
				case 2:
					// Concurrent Execution Verb Blocking
					err := svc.Invoke(ctx, &authAlpha, "tool-text-formatter", "payload")
					if err == nil {
						t.Errorf("worker %d invoke succeeded unexpectedly!", workerID)
					}
				case 3:
					// Concurrent Cross-Tenant Oracle Check
					_, err := svc.InspectTool(ctx, &authAlpha, "tenant-alpha", "tool-bravo-calculator")
					if err == nil {
						t.Errorf("worker %d cross-tenant query leaked!", workerID)
					}
				case 4:
					// Concurrent MSSP Cross-Tenant Mutation Attempt (Must always fail 403)
					reqCreate := &CreateDraftToolRequest{
						TenantRef:       "tenant-bravo",
						Name:            fmt.Sprintf("tool-concurrent-%d-%d", workerID, i),
						Purpose:         "testing concurrency",
						ExecutionNature: ExecutionNatureDeterministic,
					}
					_, err := svc.CreateDraftTool(ctx, &authMSSP, reqCreate)
					if err == nil {
						t.Errorf("worker %d MSSP cross-tenant mutation succeeded unexpectedly!", workerID)
					}
				}
			}
		}()
	}

	wg.Wait()
}
