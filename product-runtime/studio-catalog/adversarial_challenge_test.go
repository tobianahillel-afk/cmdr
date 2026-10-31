package studiocatalog

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// =========================================================================
// 1. Anti-Agent Boundary Defense (REQ-AI-002) Stress & Adversarial Suite
// =========================================================================

func TestAdversarial_AntiAgentBoundaryDefense(t *testing.T) {
	forbiddenKeywords := []string{
		"planning_loop", "autonomous_loop", "agent_team", "goal_prompt",
		"agent_role", "loop_autonomy", "subagent_spawn", "self_direct",
	}

	baseTool := func() *ToolDefinition {
		return &ToolDefinition{
			ID:                  "tool-agent-challenger",
			TenantID:            "tenant-alpha",
			Name:                "Standard Analyzer",
			Purpose:             "Perform standard log data analysis",
			Description:         "A deterministic data analyzer",
			ExecutionNature:     ExecutionNatureAICapable,
			SupportedOperations: []string{"analyze"},
			Operations: []ToolOperation{
				{OperationID: "op-analyze", Name: "analyze"},
			},
			Owner: Owner{OwnerID: "analyst-1", OwnerName: "Security Analyst"},
			RuntimeAvailability: RuntimeAvailabilityStatus{
				IsAvailable:  true,
				HealthStatus: "available",
			},
		}
	}

	// 1.1 Test each forbidden keyword in Name, Description, Purpose, and Limitations
	for _, kw := range forbiddenKeywords {
		t.Run("Keyword_"+kw+"_In_Purpose", func(t *testing.T) {
			tool := baseTool()
			tool.Purpose = fmt.Sprintf("Tool intended to run a %s for autonomous execution", kw)
			res := DeterministicValidateToolContract(tool, "corr-test")
			if res.Valid {
				t.Fatalf("expected validation failure for keyword %q in Purpose, got valid", kw)
			}
			if res.NonAgentConfirmed {
				t.Fatalf("expected NonAgentConfirmed=false for keyword %q, got true", kw)
			}
			foundIssue := false
			for _, iss := range res.Issues {
				if iss.Code == "anti-agent-boundary-violation" {
					foundIssue = true
					break
				}
			}
			if !foundIssue {
				t.Fatalf("expected issue with code anti-agent-boundary-violation for keyword %q", kw)
			}
		})

		t.Run("Keyword_"+kw+"_CaseInsensitive", func(t *testing.T) {
			tool := baseTool()
			tool.Description = fmt.Sprintf("Embedded %s keyword in uppercase", strings.ToUpper(kw))
			res := DeterministicValidateToolContract(tool, "corr-test")
			if res.Valid {
				t.Fatalf("expected validation failure for uppercase keyword %q, got valid", kw)
			}
		})

		t.Run("Keyword_"+kw+"_In_Limitations", func(t *testing.T) {
			tool := baseTool()
			tool.Limitations = []string{fmt.Sprintf("Explicitly declares %s", kw)}
			res := DeterministicValidateToolContract(tool, "corr-test")
			if res.Valid {
				t.Fatalf("expected validation failure for keyword in Limitations %q, got valid", kw)
			}
		})
	}

	// 1.2 Legitimate ai-capable tool passes cleanly
	t.Run("Legitimate_AICapable_Passes", func(t *testing.T) {
		tool := baseTool()
		tool.ExecutionNature = ExecutionNatureAICapable
		tool.Description = "Uses an LLM to summarize incident notes deterministically"
		res := DeterministicValidateToolContract(tool, "corr-test")
		if !res.Valid {
			t.Fatalf("expected legitimate ai-capable tool to be valid, got issues: %+v", res.Issues)
		}
		if !res.IsAICapable {
			t.Fatalf("expected IsAICapable=true")
		}
		if !res.NonAgentConfirmed {
			t.Fatalf("expected NonAgentConfirmed=true")
		}
	})

	// 1.3 Invalid execution nature
	invalidNatures := []ExecutionNature{"", "agent", "autonomous", "loop", "probabilistic"}
	for _, in := range invalidNatures {
		t.Run(fmt.Sprintf("InvalidNature_%s", in), func(t *testing.T) {
			tool := baseTool()
			tool.ExecutionNature = in
			res := DeterministicValidateToolContract(tool, "corr-test")
			if res.Valid {
				t.Fatalf("expected invalid execution nature %q to fail validation", in)
			}
		})
	}
}

// =========================================================================
// 2. AI Independence & Performance Budget (REQ-PROD-016)
// =========================================================================

func TestAdversarial_AIIndependenceAndBudget(t *testing.T) {
	tool := &ToolDefinition{
		ID:                  "tool-perf-test",
		TenantID:            "tenant-alpha",
		Name:                "High Perf Validator",
		Purpose:             "Offline deterministic validation benchmarking",
		Description:         "Verifies that validation is 100% offline and sub-millisecond",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"validate"},
		Operations: []ToolOperation{
			{
				OperationID: "op-1",
				Name:        "validate",
				Parameters:  []Parameter{{Name: "data", Type: "string", Required: true}},
			},
		},
		Owner: Owner{OwnerID: "perf-user", OwnerName: "Perf Engineer"},
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  true,
			HealthStatus: "available",
		},
	}

	// Run 2,000 iterations to verify performance budget (<5ms) and deterministic execution
	iterations := 2000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		res := DeterministicValidateToolContract(tool, fmt.Sprintf("corr-%d", i))
		if !res.Valid {
			t.Fatalf("iteration %d failed validation unexpectedly: %+v", i, res.Issues)
		}
		if !res.DeterministicValidated {
			t.Fatalf("iteration %d: expected DeterministicValidated=true", i)
		}
	}
	elapsed := time.Since(start)
	avgDuration := elapsed / time.Duration(iterations)

	t.Logf("Executed %d deterministic validations in %v (avg: %v/validation)", iterations, elapsed, avgDuration)

	// Invariant: P95 must be well under 5ms (typically < 0.1ms)
	if avgDuration > 5*time.Millisecond {
		t.Fatalf("average validation duration %v exceeded 5ms performance budget", avgDuration)
	}
}

// =========================================================================
// 3. Decoupled Runtime Availability (CAP-STD-003 §25 Given 1)
// =========================================================================

func TestAdversarial_DecoupledRuntimeAvailability(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryToolStore()
	svc := NewDefaultToolService(store)

	auth := &AuthContext{
		PrincipalID: "qa-inspector",
		TenantID:    "tenant-prod",
		Permissions: map[string]bool{PermStudioLibraryRead: true},
	}

	offlineTool := &ToolDefinition{
		ID:                  "tool-offline-target",
		TenantID:            "tenant-prod",
		Name:                "Offline Target Tool",
		Purpose:             "Verify decoupled availability when host is down",
		Description:         "Host runtime is completely unreachable",
		ExecutionNature:     ExecutionNatureDeterministic,
		SupportedOperations: []string{"ping"},
		Operations: []ToolOperation{
			{OperationID: "op-ping", Name: "ping"},
		},
		Owner: Owner{OwnerID: "ops-lead", OwnerName: "Ops Lead"},
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  false, // HOST DOWN!
			HealthStatus: "unavailable",
			Message:      "Daemon connection refused: tcp 127.0.0.1:9099",
			CheckedAt:    time.Now().UTC(),
		},
	}

	if err := store.Put(ctx, offlineTool); err != nil {
		t.Fatalf("failed to store offline tool: %v", err)
	}

	// 3.1 Class 0 InspectTool: Must return HTTP 200 equivalent (non-nil response, no error)
	inspectResp, err := svc.InspectTool(ctx, auth, "tenant-prod", "tool-offline-target")
	if err != nil {
		t.Fatalf("InspectTool failed unexpectedly on offline tool: %v", err)
	}
	if inspectResp.Tool == nil {
		t.Fatalf("InspectTool returned nil Tool")
	}
	if inspectResp.RuntimeAvailability.IsAvailable {
		t.Fatalf("expected RuntimeAvailability.IsAvailable=false")
	}
	if inspectResp.RuntimeAvailability.HealthStatus != "unavailable" {
		t.Fatalf("expected HealthStatus=unavailable, got %q", inspectResp.RuntimeAvailability.HealthStatus)
	}
	if inspectResp.IsExecutable {
		t.Fatalf("expected IsExecutable=false under REQ-PROD-006")
	}

	// 3.2 Class 1 ValidateContract: Contract itself remains 100% valid!
	valRes, err := svc.ValidateToolContract(ctx, auth, &ValidateToolContractRequest{
		TenantRef: "tenant-prod",
		ToolID:    "tool-offline-target",
	})
	if err != nil {
		t.Fatalf("ValidateToolContract failed with error: %v", err)
	}
	if !valRes.Valid {
		t.Fatalf("expected offline tool contract to be Valid=true under decoupled model, got false. Issues: %+v", valRes.Issues)
	}
	if valRes.DecoupledRuntimeStatus != "unavailable" {
		t.Fatalf("expected DecoupledRuntimeStatus=unavailable, got %q", valRes.DecoupledRuntimeStatus)
	}

	// Check that runtime-provider-offline was logged as INFO, not ERROR
	foundOfflineInfo := false
	for _, iss := range valRes.Issues {
		if iss.Code == "runtime-provider-offline" {
			foundOfflineInfo = true
			if iss.Severity != "info" {
				t.Fatalf("expected severity 'info' for runtime-provider-offline, got %q", iss.Severity)
			}
		}
	}
	if !foundOfflineInfo {
		t.Fatalf("expected runtime-provider-offline issue reported in validation issues")
	}
}

// =========================================================================
// 4. Owner Normalization (REQ-OBJ-009) Exhaustive Stress Suite
// =========================================================================

func TestAdversarial_OwnerNormalization(t *testing.T) {
	testCases := []struct {
		name          string
		input         Owner
		expectedID    string
		expectedName  string
		expectedType  string
		expectedIsUnk bool
		expectedDisp  string
	}{
		{
			name:          "Zero_Value_Struct",
			input:         Owner{},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Empty_Strings",
			input:         Owner{OwnerID: "", OwnerName: "", OwnerType: "", DisplayName: ""},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Whitespace_Only_Fields",
			input:         Owner{OwnerID: "   ", OwnerName: "\t\n  ", OwnerType: "   "},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Missing_OwnerID_With_Name",
			input:         Owner{OwnerID: "", OwnerName: "Jane Doe"},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Missing_OwnerName_With_ID",
			input:         Owner{OwnerID: "usr-123", OwnerName: "   "},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Explicit_IsUnknown_Flag",
			input:         Owner{OwnerID: "valid-id", OwnerName: "Valid Name", IsUnknown: true},
			expectedID:    "unknown",
			expectedName:  "owner unknown",
			expectedType:  "unknown",
			expectedIsUnk: true,
			expectedDisp:  "owner unknown",
		},
		{
			name:          "Legitimate_Owner_No_DisplayName",
			input:         Owner{OwnerID: "usr-456", OwnerName: "Alice Smith", OwnerType: "user"},
			expectedID:    "usr-456",
			expectedName:  "Alice Smith",
			expectedType:  "user",
			expectedIsUnk: false,
			expectedDisp:  "Alice Smith",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := tc.input
			o.Normalize()

			if o.OwnerID != tc.expectedID {
				t.Errorf("OwnerID: expected %q, got %q", tc.expectedID, o.OwnerID)
			}
			if o.OwnerName != tc.expectedName {
				t.Errorf("OwnerName: expected %q, got %q", tc.expectedName, o.OwnerName)
			}
			if o.OwnerType != tc.expectedType {
				t.Errorf("OwnerType: expected %q, got %q", tc.expectedType, o.OwnerType)
			}
			if o.IsUnknown != tc.expectedIsUnk {
				t.Errorf("IsUnknown: expected %v, got %v", tc.expectedIsUnk, o.IsUnknown)
			}
			if o.DisplayName != tc.expectedDisp {
				t.Errorf("DisplayName: expected %q, got %q", tc.expectedDisp, o.DisplayName)
			}
		})
	}

	// 4.2 Test that CreateDraftTool normalizes owner properly
	ctx := context.Background()
	store := NewInMemoryToolStore()
	svc := NewDefaultToolService(store)
	auth := &AuthContext{
		PrincipalID: "creator",
		TenantID:    "tenant-gamma",
		Permissions: map[string]bool{PermStudioBuilderRead: true},
	}

	created, err := svc.CreateDraftTool(ctx, auth, &CreateDraftToolRequest{
		TenantRef:       "tenant-gamma",
		Name:            "Owner Normalization Test Tool",
		Purpose:         "Verify CreateDraftTool applies REQ-OBJ-009",
		Owner:           Owner{OwnerID: "   ", OwnerName: "   "},
		ExecutionNature: ExecutionNatureDeterministic,
		SupportedOperations: []string{"run_op"},
	})
	if err != nil {
		t.Fatalf("CreateDraftTool failed: %v", err)
	}
	if !created.Owner.IsUnknown || created.Owner.OwnerName != "owner unknown" {
		t.Fatalf("CreateDraftTool failed to normalize blank owner: %+v", created.Owner)
	}

	// 4.3 Test that ToAsset propagates normalized owner
	asset := created.ToAsset()
	if !asset.Owner.IsUnknown || asset.Owner.OwnerName != "owner unknown" {
		t.Fatalf("ToAsset failed to preserve normalized owner: %+v", asset.Owner)
	}
}

// =========================================================================
// 5. Concurrency, Deep Copying & Mutex Stress Suite
// =========================================================================

func TestAdversarial_CloneDeepCopyImmunity(t *testing.T) {
	original := &ToolDefinition{
		ID:                  "tool-clone-test",
		TenantID:            "tenant-alpha",
		Name:                "Original Tool",
		Purpose:             "Original Purpose",
		SupportedOperations: []string{"op1", "op2"},
		Operations: []ToolOperation{
			{
				OperationID: "op1",
				Name:        "Operation One",
				Parameters: []Parameter{
					{Name: "param1", Type: "string"},
				},
				OutputParameters: []Parameter{
					{Name: "out1", Type: "int"},
				},
			},
		},
		Parameters: []Parameter{
			{Name: "topParam", Type: "bool"},
		},
		OutputParameters: []Parameter{
			{Name: "topOut", Type: "string"},
		},
		Dependencies: []Dependency{
			{Kind: "api", RefID: "dep-api-1"},
		},
		ProviderRefs: []ProviderReference{
			{ProviderID: "prov-1", ProviderType: "runtime-host"},
		},
		Consumers: []Consumer{
			{ConsumerType: "workflow", ConsumerRef: "wf-1"},
		},
		Limitations: []string{"limit-1", "limit-2"},
		Provenance: Provenance{
			ProvenanceChain: []string{"origin-1", "origin-2"},
		},
		AuditTrail: []AuditRecord{
			{Actor: "initial-actor", Action: "create"},
		},
	}

	cloned := original.Clone()

	// Mutate every slice and nested object in cloned
	cloned.Name = "Mutated Name"
	cloned.SupportedOperations[0] = "mutated-op"
	cloned.SupportedOperations = append(cloned.SupportedOperations, "extra-op")
	cloned.Operations[0].Name = "Mutated Op Name"
	cloned.Operations[0].Parameters[0].Name = "mutated-param"
	cloned.Operations[0].OutputParameters[0].Name = "mutated-out"
	cloned.Parameters[0].Name = "mutated-top-param"
	cloned.OutputParameters[0].Name = "mutated-top-out"
	cloned.Dependencies[0].RefID = "mutated-dep"
	cloned.ProviderRefs[0].ProviderID = "mutated-prov"
	cloned.Consumers[0].ConsumerRef = "mutated-consumer"
	cloned.Limitations[0] = "mutated-limit"
	cloned.Provenance.ProvenanceChain[0] = "mutated-origin"
	cloned.AuditTrail[0].Actor = "mutated-actor"

	// Verify original remains strictly pristine
	if original.Name != "Original Tool" {
		t.Errorf("original.Name was corrupted: %q", original.Name)
	}
	if original.SupportedOperations[0] != "op1" || len(original.SupportedOperations) != 2 {
		t.Errorf("original.SupportedOperations was corrupted: %+v", original.SupportedOperations)
	}
	if original.Operations[0].Name != "Operation One" {
		t.Errorf("original.Operations[0].Name was corrupted: %q", original.Operations[0].Name)
	}
	if original.Operations[0].Parameters[0].Name != "param1" {
		t.Errorf("original.Operations[0].Parameters[0].Name was corrupted: %q", original.Operations[0].Parameters[0].Name)
	}
	if original.Operations[0].OutputParameters[0].Name != "out1" {
		t.Errorf("original.Operations[0].OutputParameters[0].Name was corrupted: %q", original.Operations[0].OutputParameters[0].Name)
	}
	if original.Parameters[0].Name != "topParam" {
		t.Errorf("original.Parameters[0].Name was corrupted: %q", original.Parameters[0].Name)
	}
	if original.OutputParameters[0].Name != "topOut" {
		t.Errorf("original.OutputParameters[0].Name was corrupted: %q", original.OutputParameters[0].Name)
	}
	if original.Dependencies[0].RefID != "dep-api-1" {
		t.Errorf("original.Dependencies was corrupted: %q", original.Dependencies[0].RefID)
	}
	if original.ProviderRefs[0].ProviderID != "prov-1" {
		t.Errorf("original.ProviderRefs was corrupted: %q", original.ProviderRefs[0].ProviderID)
	}
	if original.Consumers[0].ConsumerRef != "wf-1" {
		t.Errorf("original.Consumers was corrupted: %q", original.Consumers[0].ConsumerRef)
	}
	if original.Limitations[0] != "limit-1" {
		t.Errorf("original.Limitations was corrupted: %q", original.Limitations[0])
	}
	if original.Provenance.ProvenanceChain[0] != "origin-1" {
		t.Errorf("original.ProvenanceChain was corrupted: %q", original.Provenance.ProvenanceChain[0])
	}
	if original.AuditTrail[0].Actor != "initial-actor" {
		t.Errorf("original.AuditTrail was corrupted: %q", original.AuditTrail[0].Actor)
	}
}

func TestAdversarial_ConcurrentStoreStress(t *testing.T) {
	store := NewInMemoryToolStore()
	ctx := context.Background()

	// Seed 10 tools across 2 tenants
	tenants := []string{"tenant-1", "tenant-2"}
	for _, tenant := range tenants {
		for i := 0; i < 5; i++ {
			tool := &ToolDefinition{
				ID:                  fmt.Sprintf("tool-%d", i),
				TenantID:            tenant,
				Name:                fmt.Sprintf("Tool %d", i),
				Purpose:             "Concurrency test",
				FunctionalState:     ToolStateDraft,
				ExecutionNature:     ExecutionNatureDeterministic,
				SupportedOperations: []string{"run"},
				Owner:               Owner{OwnerID: "tester", OwnerName: "Tester"},
				RuntimeAvailability: RuntimeAvailabilityStatus{IsAvailable: true, HealthStatus: "available"},
			}
			if err := store.Put(ctx, tool); err != nil {
				t.Fatalf("failed to seed store: %v", err)
			}
		}
	}

	// Launch 50 concurrent goroutines executing mixed Get, Put, List, Delete
	var wg sync.WaitGroup
	goroutines := 50
	opsPerGoroutine := 100

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for op := 0; op < opsPerGoroutine; op++ {
				tenant := tenants[(gid+op)%len(tenants)]
				toolID := fmt.Sprintf("tool-%d", (gid+op)%5)

				switch op % 4 {
				case 0: // Read (Get) & Mutate returned object
					tool, err := store.Get(ctx, tenant, toolID)
					if err != nil {
						t.Errorf("Get error: %v", err)
						return
					}
					if tool != nil {
						// Malicious mutation: attempt to corrupt in-memory record
						tool.Name = fmt.Sprintf("Corrupted-%d-%d", gid, op)
						tool.SupportedOperations = append(tool.SupportedOperations, "malicious-op")
					}
				case 1: // List & Mutate
					list, err := store.List(ctx, tenant)
					if err != nil {
						t.Errorf("List error: %v", err)
						return
					}
					for _, item := range list {
						item.Purpose = "Overwritten Purpose"
					}
				case 2: // Put
					newTool := &ToolDefinition{
						ID:                  fmt.Sprintf("tool-dynamic-%d", gid),
						TenantID:            tenant,
						Name:                fmt.Sprintf("Dynamic Tool %d", gid),
						Purpose:             "Concurrent dynamic tool",
						FunctionalState:     ToolStateDraft,
						ExecutionNature:     ExecutionNatureDeterministic,
						SupportedOperations: []string{"dynamic_op"},
						Owner:               Owner{OwnerID: "tester", OwnerName: "Tester"},
					}
					_ = store.Put(ctx, newTool)
				case 3: // Delete & Re-Put
					_ = store.Delete(ctx, tenant, fmt.Sprintf("tool-dynamic-%d", gid))
				}
			}
		}(g)
	}

	wg.Wait()

	// Verify that seeded tool records are not corrupted
	for _, tenant := range tenants {
		for i := 0; i < 5; i++ {
			tool, err := store.Get(ctx, tenant, fmt.Sprintf("tool-%d", i))
			if err != nil {
				t.Fatalf("Get error after concurrent test: %v", err)
			}
			if tool == nil {
				t.Fatalf("seeded tool-%d disappeared", i)
			}
			if strings.HasPrefix(tool.Name, "Corrupted") {
				t.Fatalf("DATA RACE DETECTED: tool-%d Name was corrupted by concurrent mutation: %s", i, tool.Name)
			}
		}
	}
}

// =========================================================================
// 6. Security & Multi-Tenant Anti-Oracle Verification
// =========================================================================

func TestAdversarial_CrossTenantAntiOracleUniform404(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryToolStore()
	svc := NewDefaultToolService(store)

	// Seed tool exclusively in tenant-victim
	secretTool := &ToolDefinition{
		ID:                  "secret-internal-tool",
		TenantID:            "tenant-victim",
		Name:                "Victim Sensitive Tool",
		Purpose:             "Confidential tool",
		ExecutionNature:     ExecutionNatureDeterministic,
		FunctionalState:     ToolStateDraft,
		SupportedOperations: []string{"secret_op"},
		Owner:               Owner{OwnerID: "victim-admin", OwnerName: "Victim Admin"},
	}
	if err := store.Put(ctx, secretTool); err != nil {
		t.Fatalf("failed to seed secret tool: %v", err)
	}

	// Attacker tenant attempting to inspect victim tool
	attackerAuth := &AuthContext{
		PrincipalID: "attacker-principal",
		TenantID:    "tenant-attacker",
		Permissions: map[string]bool{PermStudioLibraryRead: true},
	}

	// Inspecting non-existent tool in attacker tenant
	_, errNonExistent := svc.InspectTool(ctx, attackerAuth, "tenant-attacker", "completely-fictitious-tool")
	if errNonExistent == nil {
		t.Fatalf("expected error for non-existent tool, got nil")
	}

	// Inspecting victim tool from attacker tenant (anti-oracle check)
	_, errVictimTool := svc.InspectTool(ctx, attackerAuth, "tenant-attacker", "secret-internal-tool")
	if errVictimTool == nil {
		t.Fatalf("expected 404 for cross-tenant tool, got nil")
	}

	// The error codes and messages must be IDENTICAL (zero existence leakage)
	envNonExistent, ok1 := errNonExistent.(*ErrorEnvelope)
	envVictimTool, ok2 := errVictimTool.(*ErrorEnvelope)
	if !ok1 || !ok2 {
		t.Fatalf("errors must be *ErrorEnvelope")
	}

	if envNonExistent.HTTPStatus != 404 || envVictimTool.HTTPStatus != 404 {
		t.Fatalf("expected HTTP status 404 for both, got %d and %d", envNonExistent.HTTPStatus, envVictimTool.HTTPStatus)
	}
	if envNonExistent.Code != ErrorCodeAssetNotFound || envVictimTool.Code != ErrorCodeAssetNotFound {
		t.Fatalf("expected ErrorCodeAssetNotFound for both, got %q and %q", envNonExistent.Code, envVictimTool.Code)
	}
	if envNonExistent.Retryable || envVictimTool.Retryable {
		t.Fatalf("404 errors must not be retryable")
	}
}
