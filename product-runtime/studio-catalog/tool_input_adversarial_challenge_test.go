// Package studiocatalog adversarial challenge test suite for CAP-STD-004.
// Penetration, fuzzing, injection, secret extraction, tenant isolation,
// and execution verb bypass resistance.
package studiocatalog

import (
	"context"
	"encoding/json"
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

// -------------------------------------------------------------------------
// Challenge 7: DraftBindings Secret Non-Disclosure (PR #225 Codex Thread 4564670897)
// -------------------------------------------------------------------------

func TestChallenge_DraftBindings_SecretNonDisclosure(t *testing.T) {
	ctx := context.Background()
	auth := NewAuthContext("worker-test", "tenant-vault-test", nil, PermStudioLibraryRead)
	store := NewInMemoryToolStore()
	service := NewDefaultToolInputService(store)

	const (
		synthRawSecret1 = "synth-token-vault-secret-xyz-987"
		synthRawSecret2 = "synth-password-cleartext-pass-12345"
		synthRawSecret3 = "synth-undeclared-secret-token-777"
	)

	contract := &ToolInputContract{
		ToolID:      "tool-secret-guard",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-vault-test",
		Parameters: []InputParameterDefinition{
			{
				Name:           "api_key",
				FunctionalType: FunctionalTypeSecretRef,
				Required:       true,
			},
			{
				Name:           "sensitive_token",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				IsSensitive:    true,
			},
			{
				Name:           "regular_host",
				FunctionalType: FunctionalTypeString,
				Required:       true,
			},
		},
	}

	// Caller supplies literal secrets in base bindings and corrections
	baseBindings := []InputBinding{
		{
			ParameterName: "api_key",
			SourceType:    BindingSourceLiteral,
			Value:         synthRawSecret1, // ILLEGAL raw secret attempt!
		},
		{
			ParameterName: "regular_host",
			SourceType:    BindingSourceLiteral,
			Value:         "api.example.local",
		},
	}

	corrections := []InputBinding{
		{
			ParameterName: "sensitive_token",
			SourceType:    BindingSourceLiteral,
			Value:         synthRawSecret2, // Sensitive literal value
		},
		{
			ParameterName: "undeclared_sec",
			SourceType:    BindingSourceSecretRef,
			Value:         synthRawSecret3, // Secret-ref binding with raw value
		},
	}

	req := &PrepareCorrectedBindingsRequest{
		TenantRef:    "tenant-vault-test",
		ToolID:       "tool-secret-guard",
		ToolVersion:  "1.0.0",
		Contract:     contract,
		BaseBindings: baseBindings,
		Corrections:  corrections,
	}

	resp, err := service.PrepareCorrectedBindings(ctx, &auth, req)
	if err != nil {
		t.Fatalf("unexpected PrepareCorrectedBindings error: %v", err)
	}

	// 1. Direct inspection of DraftBindings: plaintext secrets must NEVER be present!
	for _, b := range resp.DraftBindings {
		valStr := fmt.Sprintf("%v", b.Value)
		if strings.Contains(valStr, synthRawSecret1) ||
			strings.Contains(valStr, synthRawSecret2) ||
			strings.Contains(valStr, synthRawSecret3) {
			t.Fatalf("CRITICAL SECURITY DEFECT: Raw secret exposed in resp.DraftBindings: %+v", b)
		}
		if b.ParameterName == "api_key" && b.Value != "***" {
			t.Errorf("Expected api_key Value to be redacted to '***', got: %v", b.Value)
		}
		if b.ParameterName == "sensitive_token" && b.Value != "***" {
			t.Errorf("Expected sensitive_token Value to be redacted to '***', got: %v", b.Value)
		}
	}

	// 2. JSON serialization inspection: plaintext secrets must NEVER be present!
	rawJSON, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	jsonStr := string(rawJSON)
	if strings.Contains(jsonStr, synthRawSecret1) {
		t.Fatalf("CRITICAL SECURITY DEFECT: synthRawSecret1 leaked into JSON serialization: %s", jsonStr)
	}
	if strings.Contains(jsonStr, synthRawSecret2) {
		t.Fatalf("CRITICAL SECURITY DEFECT: synthRawSecret2 leaked into JSON serialization: %s", jsonStr)
	}
	if strings.Contains(jsonStr, synthRawSecret3) {
		t.Fatalf("CRITICAL SECURITY DEFECT: synthRawSecret3 leaked into JSON serialization: %s", jsonStr)
	}

	// 3. Structured logging inspection (%v, %+v, String(), GoString()): plaintext secrets must NEVER be present!
	logFmt := fmt.Sprintf("%v %+v %s %#v", resp, resp, resp.String(), resp.GoString())
	if strings.Contains(logFmt, synthRawSecret1) ||
		strings.Contains(logFmt, synthRawSecret2) ||
		strings.Contains(logFmt, synthRawSecret3) {
		t.Fatalf("CRITICAL SECURITY DEFECT: Raw secret exposed in structured log outputs: %s", logFmt)
	}
}

// -------------------------------------------------------------------------
// Challenge 8: Validation Issue Sensitive Non-Interpolation (PR #225 Codex Thread 4564670908)
// -------------------------------------------------------------------------

func TestChallenge_ValidationIssue_SensitiveNonInterpolation(t *testing.T) {
	scope := CallerScopeContext{TenantID: "tenant-sensitive-test"}

	const (
		synthSensitiveString = "synth-secret-token-invalid-allowed-value-12345"
		synthSensitiveInt    = 999999999
		synthSensitiveFloat  = 88888.88
		synthSensitiveEnum   = "synth-illegal-secret-enum-value"
	)

	maxInt := 100.0
	minInt := 10.0
	maxFloat := 50.0
	minFloat := 5.0

	contract := &ToolInputContract{
		ToolID:      "tool-diag-test",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-sensitive-test",
		Parameters: []InputParameterDefinition{
			{
				Name:           "sensitive_str",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				IsSensitive:    true,
				AllowedValues:  []any{"allowed_token_alpha", "allowed_token_beta"},
			},
			{
				Name:           "sensitive_int_bounds",
				FunctionalType: FunctionalTypeInteger,
				Required:       true,
				IsSensitive:    true,
				Minimum:        &minInt,
				Maximum:        &maxInt,
			},
			{
				Name:           "sensitive_int_type",
				FunctionalType: FunctionalTypeInteger,
				Required:       true,
				IsSensitive:    true,
			},
			{
				Name:           "sensitive_float_bounds",
				FunctionalType: FunctionalTypeNumber,
				Required:       true,
				IsSensitive:    true,
				Minimum:        &minFloat,
				Maximum:        &maxFloat,
			},
			{
				Name:           "sensitive_enum",
				FunctionalType: FunctionalTypeEnum,
				Required:       true,
				IsSensitive:    true,
				AllowedValues:  []any{"PROD_TOKEN", "STAGING_TOKEN"},
			},
			{
				Name:           "non_sensitive_str",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				IsSensitive:    false, // Regular parameter for diagnostic comparison
				AllowedValues:  []any{"standard_alpha", "standard_beta"},
			},
		},
	}

	bindings := []InputBinding{
		{ParameterName: "sensitive_str", SourceType: BindingSourceLiteral, Value: synthSensitiveString},
		{ParameterName: "sensitive_int_bounds", SourceType: BindingSourceLiteral, Value: synthSensitiveInt},
		{ParameterName: "sensitive_int_type", SourceType: BindingSourceLiteral, Value: "not-an-int-secret-token-777"},
		{ParameterName: "sensitive_float_bounds", SourceType: BindingSourceLiteral, Value: synthSensitiveFloat},
		{ParameterName: "sensitive_enum", SourceType: BindingSourceLiteral, Value: synthSensitiveEnum},
		{ParameterName: "non_sensitive_str", SourceType: BindingSourceLiteral, Value: "invalid_regular_input"},
	}

	outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-diag-challenge")
	if outcome.Valid {
		t.Fatalf("expected validation failure for invalid sensitive inputs")
	}

	// Verify that none of the sensitive values are reflected in any issue message
	for _, issue := range outcome.Issues {
		if strings.Contains(issue.Message, synthSensitiveString) {
			t.Fatalf("CRITICAL SECURITY DEFECT: Sensitive string reflected in issue message: %s", issue.Message)
		}
		if strings.Contains(issue.Message, fmt.Sprintf("%d", synthSensitiveInt)) {
			t.Fatalf("CRITICAL SECURITY DEFECT: Sensitive integer reflected in issue message: %s", issue.Message)
		}
		if strings.Contains(issue.Message, "not-an-int-secret-token-777") {
			t.Fatalf("CRITICAL SECURITY DEFECT: Sensitive type mismatch value reflected in issue message: %s", issue.Message)
		}
		if strings.Contains(issue.Message, fmt.Sprintf("%v", synthSensitiveFloat)) {
			t.Fatalf("CRITICAL SECURITY DEFECT: Sensitive float reflected in issue message: %s", issue.Message)
		}
		if strings.Contains(issue.Message, synthSensitiveEnum) {
			t.Fatalf("CRITICAL SECURITY DEFECT: Sensitive enum reflected in issue message: %s", issue.Message)
		}
	}

	// Verify that non-sensitive parameter diagnostics DO include the invalid input for debugging
	foundNonSensitiveEcho := false
	for _, issue := range outcome.Issues {
		if issue.ParameterName == "non_sensitive_str" && strings.Contains(issue.Message, "invalid_regular_input") {
			foundNonSensitiveEcho = true
			break
		}
	}
	if !foundNonSensitiveEcho {
		t.Errorf("Expected non-sensitive diagnostic to retain field-level explanation with value")
	}
}

// -------------------------------------------------------------------------
// Challenge 9: Optional Defaults Strict Validation (PR #225 Codex Thread 4564670924)
// -------------------------------------------------------------------------

func TestChallenge_OptionalDefaults_StrictValidation(t *testing.T) {
	scope := CallerScopeContext{
		TenantID:          "tenant-defaults-test",
		AuthorizedTenants: []string{"tenant-defaults-test"},
	}

	maxPort := 65535.0
	minPort := 1.0

	t.Run("Invalid Integer Default Fails Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-int",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "port",
					FunctionalType: FunctionalTypeInteger,
					Required:       false,
					DefaultValue:   "not-an-integer-string", // FORBIDDEN: string default for integer
					Minimum:        &minPort,
					Maximum:        &maxPort,
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-bad-int")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Invalid default value 'not-an-integer-string' accepted as Valid=true")
		}
		if outcome.State != ValidationStateInvalid {
			t.Errorf("Expected state %q, got %q", ValidationStateInvalid, outcome.State)
		}
		if len(outcome.Issues) == 0 {
			t.Fatalf("Expected issues for invalid integer default")
		}
		if outcome.Issues[0].Code != string(ErrorCodeInvalidInput) {
			t.Errorf("Expected ErrorCodeInvalidInput, got %s", outcome.Issues[0].Code)
		}
	})

	t.Run("Out-of-Bounds Number Default Fails Validation", func(t *testing.T) {
		maxBound := 100.0
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-bounds",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "threshold",
					FunctionalType: FunctionalTypeNumber,
					Required:       false,
					DefaultValue:   999.5, // FORBIDDEN: exceeds maximum 100.0
					Maximum:        &maxBound,
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-out-bounds")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Out-of-bounds default value accepted as Valid=true")
		}
		if outcome.State != ValidationStateInvalid {
			t.Errorf("Expected state %q, got %q", ValidationStateInvalid, outcome.State)
		}
	})

	t.Run("Disallowed Enum Member Default Fails Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-enum",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "environment",
					FunctionalType: FunctionalTypeEnum,
					Required:       false,
					DefaultValue:   "STAGE", // Not in allowed values ["PROD", "DEV"]
					AllowedValues:  []any{"PROD", "DEV"},
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-enum")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Disallowed enum member default accepted as Valid=true")
		}
	})

	t.Run("Literal Secret in Optional Default Fails Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-secret",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "api_key",
					FunctionalType: FunctionalTypeSecretRef,
					Required:       false,
					DefaultValue:   "synth-raw-secret-token-in-default", // FORBIDDEN: raw secret
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-raw-secret")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Raw literal secret accepted in DefaultValue as Valid=true")
		}
		if outcome.State != ValidationStateInvalid {
			t.Errorf("Expected state %q, got %q", ValidationStateInvalid, outcome.State)
		}
	})

	t.Run("Stale Secret Reference in Optional Default Fails Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-stale-sec",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "api_key",
					FunctionalType: FunctionalTypeSecretRef,
					Required:       false,
					DefaultValue: &SecretReference{
						RefID:      "sec-def-stale",
						TenantID:   "tenant-defaults-test",
						SecretType: "bearer",
						IsStale:    true, // STALE!
					},
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-stale-sec")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Stale secret reference accepted in DefaultValue as Valid=true")
		}
		if outcome.State != ValidationStateRestrictedInput {
			t.Errorf("Expected state %q, got %q", ValidationStateRestrictedInput, outcome.State)
		}
	})

	t.Run("Cross-Tenant Entity Reference in Optional Default Fails Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-bad-def-cross-tenant",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "case_ref",
					FunctionalType: FunctionalTypeEntityRef,
					Required:       false,
					DefaultValue: &EntityReference{
						RefID:    "case-def-cross",
						Kind:     "case",
						TenantID: "tenant-other-corp", // Cross-tenant!
					},
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-cross-tenant")
		if outcome.Valid {
			t.Fatalf("CRITICAL: Cross-tenant entity reference accepted in DefaultValue as Valid=true")
		}
		if outcome.State != ValidationStateCrossScopeBlocked {
			t.Errorf("Expected state %q, got %q", ValidationStateCrossScopeBlocked, outcome.State)
		}
	})

	t.Run("Valid Optional Default Passes Validation", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:   "tool-valid-def",
			TenantID: "tenant-defaults-test",
			Parameters: []InputParameterDefinition{
				{
					Name:           "port",
					FunctionalType: FunctionalTypeInteger,
					Required:       false,
					DefaultValue:   8080,
					Minimum:        &minPort,
					Maximum:        &maxPort,
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, []InputBinding{}, scope, "corr-def-valid")
		if !outcome.Valid {
			t.Fatalf("Expected valid optional default to succeed, got issues: %+v", outcome.Issues)
		}
		if outcome.State != ValidationStateValid {
			t.Errorf("Expected state %q, got %q", ValidationStateValid, outcome.State)
		}
		if len(outcome.ValidatedBindings) != 1 {
			t.Fatalf("Expected 1 validated binding projected from default, got %d", len(outcome.ValidatedBindings))
		}
		if outcome.ValidatedBindings[0].ValidationState != ValidationStateValid {
			t.Errorf("Expected ValidatedBinding state %q, got %q", ValidationStateValid, outcome.ValidatedBindings[0].ValidationState)
		}
	})
}

// -------------------------------------------------------------------------
// Challenge 10: Deep Defensive Copies & Store Mutation Resistance (PR #225 Codex Thread 4564670935)
// -------------------------------------------------------------------------

func TestChallenge_DeepCopy_NonMutation(t *testing.T) {
	ctx := context.Background()
	auth := NewAuthContext("worker-inspect", "tenant-alpha", nil, PermStudioLibraryRead)
	store := NewInMemoryToolStore()
	service := NewDefaultToolInputService(store)

	origMax := 100.0
	origMin := 10.0
	origAllowed := []any{"prod", "staging"}
	origDefault := map[string]any{"timeout": 30}

	toolDef := &ToolDefinition{
		ID:          "tool-immutable",
		TenantID:    "tenant-alpha",
		Name:        "Immutable Tool",
		Description: "Tool for non-mutation testing",
		Version:     VersionInfo{SemVer: "1.0.0"},
		InputContract: &ToolInputContract{
			ToolID:      "tool-immutable",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{
					Name:           "rate_limit",
					FunctionalType: FunctionalTypeInteger,
					Required:       true,
					Minimum:        &origMin,
					Maximum:        &origMax,
					AllowedValues:  origAllowed,
					DefaultValue:   origDefault,
					Constraints: &InputConstraint{
						Minimum:       &origMin,
						Maximum:       &origMax,
						AllowedValues: origAllowed,
					},
				},
			},
		},
	}

	if err := store.Put(ctx, toolDef); err != nil {
		t.Fatalf("store.Put failed: %v", err)
	}

	t.Run("Modifying Inspected Contract Cannot Mutate Stored Contract", func(t *testing.T) {
		inspected, err := service.InspectToolInputContract(ctx, &auth, "tenant-alpha", "tool-immutable")
		if err != nil {
			t.Fatalf("InspectToolInputContract failed: %v", err)
		}

		// Malicious or accidental in-place mutation of inspected contract
		*inspected.Parameters[0].Constraints.Maximum = 999999.0
		*inspected.Parameters[0].Constraints.Minimum = -999999.0
		*inspected.Parameters[0].Maximum = 888888.0
		*inspected.Parameters[0].Minimum = -888888.0
		inspected.Parameters[0].Constraints.AllowedValues[0] = "MUTATED_CONSTRAINT_VALUE"
		inspected.Parameters[0].AllowedValues[0] = "MUTATED_PARAM_VALUE"
		if defMap, ok := inspected.Parameters[0].DefaultValue.(map[string]any); ok {
			defMap["timeout"] = 9999
		}

		// Read back from the store afresh
		freshTool, err := store.Get(ctx, "tenant-alpha", "tool-immutable")
		if err != nil {
			t.Fatalf("store.Get failed: %v", err)
		}
		if freshTool == nil || freshTool.InputContract == nil {
			t.Fatalf("expected non-nil tool and input contract")
		}

		freshParam := freshTool.InputContract.Parameters[0]

		// Assert that the stored contract remained 100% untouched
		if *freshParam.Constraints.Maximum != origMax {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Constraints.Maximum mutated in store! Expected %v, got %v", origMax, *freshParam.Constraints.Maximum)
		}
		if *freshParam.Constraints.Minimum != origMin {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Constraints.Minimum mutated in store! Expected %v, got %v", origMin, *freshParam.Constraints.Minimum)
		}
		if *freshParam.Maximum != origMax {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Parameter.Maximum mutated in store! Expected %v, got %v", origMax, *freshParam.Maximum)
		}
		if *freshParam.Minimum != origMin {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Parameter.Minimum mutated in store! Expected %v, got %v", origMin, *freshParam.Minimum)
		}
		if freshParam.Constraints.AllowedValues[0] != "prod" {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Constraints.AllowedValues mutated in store! Expected 'prod', got %v", freshParam.Constraints.AllowedValues[0])
		}
		if freshParam.AllowedValues[0] != "prod" {
			t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: Parameter.AllowedValues mutated in store! Expected 'prod', got %v", freshParam.AllowedValues[0])
		}
		if defMap, ok := freshParam.DefaultValue.(map[string]any); ok {
			if defMap["timeout"] != 30 {
				t.Fatalf("CRITICAL STORE INTEGRITY VIOLATION: DefaultValue map mutated in store! Expected 30, got %v", defMap["timeout"])
			}
		}
	})

	t.Run("Concurrent Inspection and Mutation Harness", func(t *testing.T) {
		const numWorkers = 50
		var wg sync.WaitGroup
		errCh := make(chan error, numWorkers)

		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				inspected, err := service.InspectToolInputContract(ctx, &auth, "tenant-alpha", "tool-immutable")
				if err != nil {
					errCh <- err
					return
				}
				// Mutate local clone heavily
				newVal := float64(workerID * 1000)
				*inspected.Parameters[0].Constraints.Maximum = newVal
				*inspected.Parameters[0].Maximum = newVal

				// Verify clone has local value
				if *inspected.Parameters[0].Constraints.Maximum != newVal {
					errCh <- fmt.Errorf("worker %d: local mutation lost", workerID)
					return
				}
			}(w)
		}

		wg.Wait()
		close(errCh)
		for err := range errCh {
			t.Fatalf("Concurrent inspection failure: %v", err)
		}

		// Final check on store
		finalTool, _ := store.Get(ctx, "tenant-alpha", "tool-immutable")
		if *finalTool.InputContract.Parameters[0].Constraints.Maximum != origMax {
			t.Fatalf("CRITICAL: Stored contract corrupted during concurrent access")
		}
	})
}

