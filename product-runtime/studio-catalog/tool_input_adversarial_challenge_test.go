// Package studiocatalog adversarial challenge test suite for CAP-STD-004.
// Penetration, fuzzing, injection, secret extraction, tenant isolation,
// and execution verb bypass resistance.
package studiocatalog

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
)

// -------------------------------------------------------------------------
// Challenge 1: Null, Nil, and Malformed Boundary Edge Cases
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_NullAndMalformedEdgeCases(t *testing.T) {
	scope := CallerScopeContext{TenantID: "tenant-alpha"}

	t.Run("Nil Contract Handled Safely", func(t *testing.T) {
		outcome := DeterministicValidateInputBindings(nil, []InputBinding{}, scope, "corr-nil-contract")
		if outcome.Valid {
			t.Errorf("Expected nil contract to fail validation")
		}
		if outcome.State != ValidationStateInvalid {
			t.Errorf("Expected state %q, got %q", ValidationStateInvalid, outcome.State)
		}
	})

	t.Run("Nil Bindings Slice with Required Parameters", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-edge",
			TenantID: "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "req_param", FunctionalType: FunctionalTypeString, Required: true},
			},
		}
		outcome := DeterministicValidateInputBindings(contract, nil, scope, "corr-nil-bindings")
		if outcome.Valid {
			t.Errorf("Expected nil bindings with required parameter to fail")
		}
		if outcome.State != ValidationStateMissingInput {
			t.Errorf("Expected missing-input state, got %q", outcome.State)
		}
	})

	t.Run("Empty Struct Binding with Zero Values", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-edge",
			TenantID: "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "int_param", FunctionalType: FunctionalTypeInteger, Required: true},
			},
		}
		// Empty binding has nil Value
		emptyBinding := []InputBinding{
			{ParameterName: "int_param"},
		}
		outcome := DeterministicValidateInputBindings(contract, emptyBinding, scope, "corr-empty-struct")
		if outcome.Valid {
			t.Errorf("Expected empty binding for required int to fail")
		}
	})

	t.Run("Extreme Integer Values and Floating-point Boundaries", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-edge",
			TenantID: "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "max_int", FunctionalType: FunctionalTypeInteger, Required: true},
				{Name: "min_int", FunctionalType: FunctionalTypeInteger, Required: true},
				{Name: "float_inf", FunctionalType: FunctionalTypeNumber, Required: true},
			},
		}

		bindings := []InputBinding{
			{ParameterName: "max_int", SourceType: BindingSourceLiteral, Value: int64(math.MaxInt64)},
			{ParameterName: "min_int", SourceType: BindingSourceLiteral, Value: int64(math.MinInt64)},
			{ParameterName: "float_inf", SourceType: BindingSourceLiteral, Value: 1e308},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-extreme-values")
		if !outcome.Valid {
			t.Errorf("Expected extreme valid numbers to be handled safely, got issues: %+v", outcome.Issues)
		}
	})
}

// -------------------------------------------------------------------------
// Challenge 2: Injection Payloads & Malicious Syntax Resistance
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_InjectionPayloads(t *testing.T) {
	contract := &ToolInputContract{
		ToolID:   "tool-security",
		TenantID: "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{
				Name:           "filter_expression",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				Pattern:        `^[a-zA-Z0-9_\-\. ]+$`, // Strict alphanumeric and dots
			},
			{
				Name:           "free_text",
				FunctionalType: FunctionalTypeString,
				Required:       false,
			},
		},
	}
	scope := CallerScopeContext{TenantID: "tenant-alpha"}

	payloads := []struct {
		name        string
		payload     string
		expectMatch bool
	}{
		{"SQL Injection (OR 1=1)", "' OR 1=1 --", false},
		{"SQL Injection (DROP TABLE)", "admin'; DROP TABLE users; --", false},
		{"Command Injection ($(whoami))", "$(whoami)", false},
		{"Command Injection (; rm -rf)", "; rm -rf /", false},
		{"Command Injection (| cat)", "| cat /etc/passwd", false},
		{"Path Traversal (../../)", "../../../../etc/shadow", false},
		{"Null Byte Injection", "valid_part\x00malicious", false},
		{"XSS Script Tag", "<script>alert('pwned')</script>", false},
		{"JSON Injection", `{"is_admin": true, "role": "root"}`, false},
		{"Safe Alphanumeric Host", "valid.host.example.com", true},
	}

	for _, tc := range payloads {
		t.Run(tc.name, func(t *testing.T) {
			bindings := []InputBinding{
				{
					ParameterName: "filter_expression",
					SourceType:    BindingSourceLiteral,
					Value:         tc.payload,
				},
				{
					ParameterName: "free_text",
					SourceType:    BindingSourceLiteral,
					Value:         tc.payload,
				},
			}

			outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-injection")
			if tc.expectMatch {
				if !outcome.Valid {
					t.Errorf("Expected valid payload %q to pass, got: %+v", tc.payload, outcome.Issues)
				}
			} else {
				if outcome.Valid {
					t.Errorf("CRITICAL SECURITY DEFECT: Injection payload %q unexpectedly bypassed pattern constraint", tc.payload)
				}
			}
		})
	}
}

// -------------------------------------------------------------------------
// Challenge 3: Secret Extraction and Plaintext Leakage Attempts
// (REQ-OBJ-009, REQ-SEC-001, CAP-SET-009)
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_SecretExtractionAttempts(t *testing.T) {
	contract := &ToolInputContract{
		ToolID:   "tool-vault",
		TenantID: "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{
				Name:           "vault_token",
				FunctionalType: FunctionalTypeSecretRef,
				Required:       true,
				SecretRefPolicy: &SecretRefPolicy{
					AllowedSecretTypes: []string{"bearer"},
				},
			},
			{
				Name:           "sensitive_passphrase",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				IsSensitive:    true, // Sensitive literal field
			},
		},
	}
	scope := CallerScopeContext{TenantID: "tenant-alpha"}

	t.Run("Plaintext Secret Injection Blocked with Zero Echo", func(t *testing.T) {
		rawSecret := "SG.abcdef1234567890_VERY_SECRET_KEY"
		bindings := []InputBinding{
			{
				ParameterName: "vault_token",
				SourceType:    BindingSourceLiteral,
				Value:         rawSecret, // Raw secret attempt
			},
			{
				ParameterName: "sensitive_passphrase",
				SourceType:    BindingSourceLiteral,
				Value:         "MyUltraSecretPassword123!",
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-extract-1")
		if outcome.Valid {
			t.Fatalf("Expected validation failure for raw secret injection")
		}

		// Ensure raw secret is NEVER reflected anywhere in issues, messages, or validated bindings
		for _, issue := range outcome.Issues {
			if strings.Contains(issue.Message, rawSecret) {
				t.Fatalf("CRITICAL SECURITY VULNERABILITY: Raw secret reflected in validation issue: %s", issue.Message)
			}
		}
		for _, vb := range outcome.ValidatedBindings {
			if vb.MaskedValue == rawSecret {
				t.Fatalf("CRITICAL SECURITY VULNERABILITY: Raw secret found unmasked in ValidatedBindings: %v", vb.MaskedValue)
			}
		}
	})

	t.Run("Tampered Secret Owner Overridden to Settings in Validated Binding Without Caller Mutation", func(t *testing.T) {
		sec := &SecretReference{
			RefID:      "sec-tampered",
			Owner:      "AttackerControlledStore", // Tampering attempt
			TenantID:   "tenant-alpha",
			SecretType: "bearer",
			IsMasked:   false,
		}

		bindings := []InputBinding{
			{
				ParameterName: "vault_token",
				SourceType:    BindingSourceSecretRef,
				SecretRef:     sec,
			},
			{
				ParameterName: "sensitive_passphrase",
				SourceType:    BindingSourceLiteral,
				Value:         "Password99!",
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-tamper-owner")
		if !outcome.Valid {
			t.Fatalf("Expected valid outcome, got: %+v", outcome.Issues)
		}

		// Invariant: Caller input MUST NOT be mutated (CAP-STD-004 §10 "no caller mutation", §19)
		if sec.Owner != "AttackerControlledStore" {
			t.Errorf("CRITICAL VIOLATION: Caller input was mutated! sec.Owner changed to %q", sec.Owner)
		}
		if sec.IsMasked {
			t.Errorf("CRITICAL VIOLATION: Caller input was mutated! sec.IsMasked changed to true")
		}

		// Verify Owner was normalized strictly to "Settings" in the derived ValidatedBinding under CAP-SET-009 / §25 AC-2
		foundVaultToken := false
		for _, vb := range outcome.ValidatedBindings {
			if vb.ParameterName == "vault_token" {
				foundVaultToken = true
				if vb.Owner != "Settings" {
					t.Errorf("CRITICAL: Validated binding owner must be 'Settings', got %q", vb.Owner)
				}
				if !vb.IsMasked {
					t.Errorf("CRITICAL: Validated binding must be masked")
				}
				if vb.MaskedValue != "***" {
					t.Errorf("CRITICAL: Validated binding MaskedValue must be '***'")
				}
			}
		}
		if !foundVaultToken {
			t.Fatalf("Expected validated binding for vault_token")
		}
	})

	t.Run("Disallowed Secret Type Blocked by Policy", func(t *testing.T) {
		bindings := []InputBinding{
			{
				ParameterName: "vault_token",
				SourceType:    BindingSourceSecretRef,
				SecretRef: &SecretReference{
					RefID:      "sec-ssh",
					Owner:      "Settings",
					TenantID:   "tenant-alpha",
					SecretType: "ssh-private-key", // Not allowed by policy (only "bearer")
				},
			},
			{
				ParameterName: "sensitive_passphrase",
				SourceType:    BindingSourceLiteral,
				Value:         "Pass123",
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-disallowed-sec-type")
		if outcome.Valid {
			t.Fatalf("Expected validation failure for disallowed secret type")
		}
		if outcome.State != ValidationStateRestrictedInput {
			t.Errorf("Expected state %q, got %q", ValidationStateRestrictedInput, outcome.State)
		}
	})
}

// -------------------------------------------------------------------------
// Challenge 4: Cross-Tenant Bypass Attacks & Wildcard Evasions
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_CrossTenantBypassAttempts(t *testing.T) {
	contract := sampleToolInputContract("tenant-corp-a", "tool-sec-scanner")

	t.Run("Wildcard Tenant In Caller Scope Strictly Rejected", func(t *testing.T) {
		wildcards := []string{"*", "%", "ALL", "all", "All"}
		for _, w := range wildcards {
			scope := CallerScopeContext{TenantID: w}
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-wc")
			if outcome.Valid {
				t.Errorf("Expected wildcard tenant %q to be strictly rejected, got Valid=true", w)
			}
			foundWcIssue := false
			for _, issue := range outcome.Issues {
				if issue.Code == string(ErrorCodeWildcardTenant) {
					foundWcIssue = true
					break
				}
			}
			if !foundWcIssue {
				t.Errorf("Expected ErrorCodeWildcardTenant for wildcard %q", w)
			}
		}
	})

	t.Run("Blank Tenant In Caller Scope Strictly Rejected", func(t *testing.T) {
		scope := CallerScopeContext{TenantID: "   "}
		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-blank")
		if outcome.Valid {
			t.Errorf("Expected blank tenant to fail")
		}
		foundMissing := false
		for _, issue := range outcome.Issues {
			if issue.Code == string(ErrorCodeMissingTenant) {
				foundMissing = true
				break
			}
		}
		if !foundMissing {
			t.Errorf("Expected ErrorCodeMissingTenant for blank tenant")
		}
	})

	t.Run("Substring and Prefix Tenant Spoofing Rejected", func(t *testing.T) {
		// Target tenant is "tenant-corp-a"
		// Attacker attempts prefix "tenant-corp-a-extra" or suffix "corp-a"
		spoofedTenants := []string{
			"tenant-corp-a-extra",
			"tenant-corp-b",
			"corp-a",
			"tenant-corp",
			"tenant-corp-a\x00suffix",
		}

		scope := CallerScopeContext{TenantID: "tenant-corp-a"}

		for _, spoofed := range spoofedTenants {
			t.Run(spoofed, func(t *testing.T) {
				bindings := []InputBinding{
					{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
					{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 443},
					{
						ParameterName: "api_key",
						SourceType:    BindingSourceSecretRef,
						SecretRef:     &SecretReference{RefID: "s-1", TenantID: "tenant-corp-a"},
					},
					{
						ParameterName: "incident_case",
						SourceType:    BindingSourceEntityRef,
						EntityRef: &EntityReference{
							RefID:    "case-target",
							Kind:     "case",
							TenantID: spoofed, // SPOOFED TENANT!
						},
					},
				}

				outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-spoof")
				if outcome.Valid {
					t.Fatalf("CRITICAL SECURITY DEFECT: Spoofed tenant %q bypassed cross-tenant isolation", spoofed)
				}
				if outcome.State != ValidationStateCrossScopeBlocked {
					t.Errorf("Expected state cross-scope-blocked for spoofed tenant %q, got %q", spoofed, outcome.State)
				}
			})
		}
	})
}

// -------------------------------------------------------------------------
// Challenge 5: Execution Verb Bypass Attacks (REQ-PROD-006)
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_ExecutionVerbBypassAttempts(t *testing.T) {
	service := NewDefaultToolInputService(nil)
	ctx := context.Background()

	// Exhaustive auth permutations: admin, system, empty, nil
	authVariants := []*AuthContext{
		nil,
		{},
		{PrincipalID: "admin", TenantID: "tenant-alpha", Permissions: map[string]bool{"*": true, "admin": true, PermStudioLibraryRead: true}},
		{PrincipalID: "system", TenantID: "tenant-alpha", Permissions: map[string]bool{"perm.studio.artifact.manage": true}},
	}

	payloadVariants := []any{
		nil,
		"",
		map[string]any{"command": "exec", "force": true},
		`{"bypass": true, "elevate": true}`,
		[]byte("destructive payload"),
	}

	for _, auth := range authVariants {
		for _, payload := range payloadVariants {
			calls := []struct {
				verb string
				fn   func() error
			}{
				{"Invoke", func() error { return service.Invoke(ctx, auth, "tool-pwn", payload) }},
				{"Execute", func() error { return service.Execute(ctx, auth, "tool-pwn", payload) }},
				{"Run", func() error { return service.Run(ctx, auth, "tool-pwn", payload) }},
				{"Call", func() error { return service.Call(ctx, auth, "tool-pwn", payload) }},
				{"Trigger", func() error { return service.Trigger(ctx, auth, "tool-pwn", payload) }},
			}

			for _, c := range calls {
				err := c.fn()
				if err == nil {
					t.Fatalf("CRITICAL SECURITY FAILURE: Verb %s succeeded with auth=%+v, payload=%+v", c.verb, auth, payload)
				}

				var env *ErrorEnvelope
				if !errors.As(err, &env) {
					t.Fatalf("expected *ErrorEnvelope from verb %s, got %T", c.verb, err)
				}
				if env.HTTPStatus != 400 {
					t.Errorf("expected HTTP 400 for verb %s, got %d", c.verb, env.HTTPStatus)
				}
				if env.Code != ErrorCodeExecutionBlocked {
					t.Errorf("expected ErrorCodeExecutionBlocked, got %q", env.Code)
				}
			}
		}
	}
}

// -------------------------------------------------------------------------
// Challenge 6: High-Contention Concurrency & Deterministic Idempotence Harness
// -------------------------------------------------------------------------

func TestChallenge_ToolInput_HighContentionHarness(t *testing.T) {
	contract := sampleToolInputContract("tenant-stress", "tool-stress-01")
	scope := CallerScopeContext{
		PrincipalID: "stress-worker",
		TenantID:    "tenant-stress",
	}

	baseBindings := []InputBinding{
		{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "stress.host.local"},
		{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 8443},
		{
			ParameterName: "api_key",
			SourceType:    BindingSourceSecretRef,
			SecretRef:     &SecretReference{RefID: "sec-stress-1", TenantID: "tenant-stress", SecretType: "bearer"},
		},
		{ParameterName: "scan_mode", SourceType: BindingSourceLiteral, Value: "deep"},
	}

	const numGoroutines = 50
	const iterationsPerGoroutine = 50

	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines*iterationsPerGoroutine)

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < iterationsPerGoroutine; i++ {
				corrID := fmt.Sprintf("corr-stress-%d-%d", gid, i)
				outcome := DeterministicValidateInputBindings(contract, baseBindings, scope, corrID)

				if !outcome.Valid {
					errCh <- fmt.Errorf("goroutine %d iter %d: expected valid, got issues: %+v", gid, i, outcome.Issues)
					return
				}
				if outcome.State != ValidationStateValid {
					errCh <- fmt.Errorf("goroutine %d iter %d: expected state valid, got %s", gid, i, outcome.State)
					return
				}
				if !outcome.DeterministicValidated {
					errCh <- fmt.Errorf("goroutine %d iter %d: expected deterministic validation flag", gid, i)
					return
				}
				if outcome.CorrelationID != corrID {
					errCh <- fmt.Errorf("goroutine %d iter %d: correlation ID mismatch", gid, i)
					return
				}
			}
		}(g)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("Concurrent stress test failure: %v", err)
	}
}
