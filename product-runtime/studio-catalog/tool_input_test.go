// Package studiocatalog unit test suite for CAP-STD-004
// Tool Input, Parameter and Validation Contract.
package studiocatalog

import (
	"context"
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
)

// -------------------------------------------------------------------------
// Test Fixtures & Helpers
// -------------------------------------------------------------------------

func sampleToolInputContract(tenantID, toolID string) *ToolInputContract {
	minL := 3
	maxL := 20
	minV := 1.0
	maxV := 65535.0

	return &ToolInputContract{
		ToolID:      toolID,
		ToolVersion: "1.0.0",
		TenantID:    tenantID,
		Parameters: []InputParameterDefinition{
			{
				Name:           "target_host",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				Description:    "Target hostname or IP address",
				MinLength:      &minL,
				MaxLength:      &maxL,
			},
			{
				Name:           "port",
				FunctionalType: FunctionalTypeInteger,
				Required:       true,
				Description:    "Target port",
				Minimum:        &minV,
				Maximum:        &maxV,
			},
			{
				Name:           "use_tls",
				FunctionalType: FunctionalTypeBoolean,
				Required:       false,
				Description:    "Enable TLS",
				DefaultValue:   true,
			},
			{
				Name:           "api_key",
				FunctionalType: FunctionalTypeSecretRef,
				Required:       true,
				Description:    "API secret reference owned by Settings",
				SecretRefPolicy: &SecretRefPolicy{
					AllowedSecretTypes: []string{"bearer", "api-key"},
				},
			},
			{
				Name:           "incident_case",
				FunctionalType: FunctionalTypeEntityRef,
				Required:       false,
				Description:    "Associated Investigate Case reference",
				EntityScope: &EntityScope{
					TargetKind: "case",
				},
			},
			{
				Name:           "log_query",
				FunctionalType: FunctionalTypeQueryAssetRef,
				Required:       false,
				Description:    "Saved log query asset reference",
			},
			{
				Name:           "scan_mode",
				FunctionalType: FunctionalTypeEnum,
				Required:       false,
				Description:    "Scan operational mode",
				AllowedValues:  []any{"quick", "thorough", "deep"},
				DefaultValue:   "quick",
			},
			{
				Name:           "tags",
				FunctionalType: FunctionalTypeArray,
				Required:       false,
				Description:    "List of tag labels",
			},
			{
				Name:           "config_options",
				FunctionalType: FunctionalTypeObject,
				Required:       false,
				Description:    "Map of additional parameters",
			},
			{
				Name:           "timeout_sec",
				FunctionalType: FunctionalTypeNumber,
				Required:       false,
				Description:    "Floating-point timeout duration",
			},
		},
	}
}

// -------------------------------------------------------------------------
// AC-1: Given required input missing, When validation runs,
// Then it fails explicitly with missing-input and invents NO default (CAP-STD-004 §25)
// -------------------------------------------------------------------------

func TestToolInput_AC1_MissingRequiredInput_NeverInventDefault(t *testing.T) {
	contract := sampleToolInputContract("tenant-alpha", "tool-network-scanner")
	scope := CallerScopeContext{
		PrincipalID: "analyst-1",
		TenantID:    "tenant-alpha",
	}

	// Supply bindings for port and api_key, but omit required "target_host"
	bindings := []InputBinding{
		{
			ParameterName:   "port",
			SourceType:      BindingSourceLiteral,
			Value:           443,
			ValidationState: ValidationStateUnvalidated,
		},
		{
			ParameterName: "api_key",
			SourceType:    BindingSourceSecretRef,
			SecretRef: &SecretReference{
				RefID:      "sec-ref-001",
				Owner:      "Settings",
				TenantID:   "tenant-alpha",
				SecretType: "api-key",
			},
			ValidationState: ValidationStateUnvalidated,
		},
	}

	outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-test-ac1")

	if outcome.Valid {
		t.Fatalf("Expected validation to fail due to missing required parameter, got Valid=true")
	}

	if outcome.State != ValidationStateMissingInput {
		t.Errorf("Expected outcome state %q, got %q", ValidationStateMissingInput, outcome.State)
	}

	// Verify the missing input issue is explicit
	foundMissing := false
	for _, issue := range outcome.Issues {
		if issue.ParameterName == "target_host" && issue.Code == string(ErrorCodeMissingInput) {
			foundMissing = true
			if issue.State != ValidationStateMissingInput {
				t.Errorf("Expected issue state %q, got %q", ValidationStateMissingInput, issue.State)
			}
			if !strings.Contains(issue.Message, "target_host") {
				t.Errorf("Expected issue message to name target_host, got %q", issue.Message)
			}
		}
	}
	if !foundMissing {
		t.Errorf("Expected issue with code 'missing-input' for target_host, got: %+v", outcome.Issues)
	}

	// Invariant: verify NO default was fabricated for target_host
	for _, vb := range outcome.ValidatedBindings {
		if vb.ParameterName == "target_host" {
			t.Errorf("CRITICAL INVARIANT VIOLATION: validated binding found for missing target_host: %+v", vb)
		}
	}

	// Verify optional parameter with explicit default used that default ("use_tls" -> true, "scan_mode" -> "quick")
	foundTLSDefault := false
	for _, vb := range outcome.ValidatedBindings {
		if vb.ParameterName == "use_tls" {
			foundTLSDefault = true
			if vb.MaskedValue != true {
				t.Errorf("Expected default value true for use_tls, got %v", vb.MaskedValue)
			}
		}
	}
	if !foundTLSDefault {
		t.Errorf("Expected explicit default value for optional use_tls")
	}
}

// -------------------------------------------------------------------------
// AC-2: Given Secret Reference, When inspected/validated,
// Then raw secret stays masked, Settings ownership visible, raw secret rejected
// -------------------------------------------------------------------------

func TestToolInput_AC2_SecretReferenceMasking_And_SettingsOwnership(t *testing.T) {
	contract := sampleToolInputContract("tenant-alpha", "tool-network-scanner")
	scope := CallerScopeContext{
		PrincipalID: "analyst-1",
		TenantID:    "tenant-alpha",
	}

	t.Run("Nominal Valid Secret Reference Masking", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceSecretRef,
				SecretRef: &SecretReference{
					RefID:      "sec-ref-alpha-999",
					Owner:      "Settings",
					TenantID:   "tenant-alpha",
					SecretType: "bearer",
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-test-ac2")
		if !outcome.Valid {
			t.Fatalf("Expected valid outcome, got issues: %+v", outcome.Issues)
		}

		foundSecret := false
		for _, vb := range outcome.ValidatedBindings {
			if vb.ParameterName == "api_key" {
				foundSecret = true
				if vb.MaskedValue != "***" {
					t.Errorf("CRITICAL SECURITY VIOLATION: secret masked_value must be '***', got %v", vb.MaskedValue)
				}
				if !vb.IsMasked {
					t.Errorf("Expected IsMasked=true")
				}
				if vb.RefID != "sec-ref-alpha-999" {
					t.Errorf("Expected RefID 'sec-ref-alpha-999', got %q", vb.RefID)
				}
			}
		}
		if !foundSecret {
			t.Errorf("Expected validated binding for api_key")
		}

		// Verify masked parameters list
		hasMaskedParam := false
		for _, mp := range outcome.MaskedParameters {
			if mp == "api_key" {
				hasMaskedParam = true
				break
			}
		}
		if !hasMaskedParam {
			t.Errorf("Expected api_key in MaskedParameters list: %v", outcome.MaskedParameters)
		}
	})

	t.Run("Raw Secret String Passed As Literal Is Strictly Blocked", func(t *testing.T) {
		// Attacker or misconfigured caller tries to pass raw credential directly
		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceLiteral,
				Value:         "SUPER_SECRET_RAW_TOKEN_123456", // RAW SECRET
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-test-raw-sec")
		if outcome.Valid {
			t.Fatalf("Expected validation to fail when raw secret is passed, got Valid=true")
		}

		foundSecIssue := false
		for _, issue := range outcome.Issues {
			if issue.ParameterName == "api_key" {
				foundSecIssue = true
				if issue.Code != string(ErrorCodeInvalidInput) {
					t.Errorf("Expected ErrorCodeInvalidInput, got %q", issue.Code)
				}
				// Verify raw secret is NEVER leaked in the error message!
				if strings.Contains(issue.Message, "SUPER_SECRET_RAW_TOKEN_123456") {
					t.Fatalf("CRITICAL SECURITY LEAK: raw secret reflected in issue message: %q", issue.Message)
				}
			}
		}
		if !foundSecIssue {
			t.Errorf("Expected validation issue for api_key")
		}
	})

	t.Run("Stale Secret Reference Fails with RestrictedInput", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceSecretRef,
				SecretRef: &SecretReference{
					RefID:      "sec-ref-expired",
					Owner:      "Settings",
					TenantID:   "tenant-alpha",
					SecretType: "api-key",
					IsStale:    true, // STALE / EXPIRED
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-test-stale-sec")
		if outcome.Valid {
			t.Fatalf("Expected validation to fail for stale secret reference")
		}
		if outcome.State != ValidationStateRestrictedInput {
			t.Errorf("Expected state %q, got %q", ValidationStateRestrictedInput, outcome.State)
		}

		foundStale := false
		for _, issue := range outcome.Issues {
			if issue.ParameterName == "api_key" && issue.Code == string(ErrorCodeStaleReference) {
				foundStale = true
				if issue.State != ValidationStateRestrictedInput {
					t.Errorf("Expected issue state %q, got %q", ValidationStateRestrictedInput, issue.State)
				}
			}
		}
		if !foundStale {
			t.Errorf("Expected issue for stale reference")
		}
	})
}

// -------------------------------------------------------------------------
// AC-3: Given IA off, When bindings are validated or prepared,
// Then deterministic validation and manual binding remain 100% usable (CAP-STD-004 §25)
// -------------------------------------------------------------------------

func TestToolInput_AC3_DeterministicValidation_WithoutAI(t *testing.T) {
	contract := sampleToolInputContract("tenant-alpha", "tool-network-scanner")
	scope := CallerScopeContext{
		PrincipalID: "analyst-1",
		TenantID:    "tenant-alpha",
	}

	bindings := []InputBinding{
		{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "scan.target.local"},
		{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 8080},
		{
			ParameterName: "api_key",
			SourceType:    BindingSourceSecretRef,
			SecretRef: &SecretReference{
				RefID:      "sec-valid-01",
				Owner:      "Settings",
				TenantID:   "tenant-alpha",
				SecretType: "bearer",
			},
		},
		{ParameterName: "scan_mode", SourceType: BindingSourceLiteral, Value: "thorough"},
	}

	outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-test-ac3")

	if !outcome.Valid {
		t.Fatalf("Expected valid outcome, got issues: %+v", outcome.Issues)
	}

	if !outcome.DeterministicValidated {
		t.Errorf("Expected DeterministicValidated=true")
	}

	if outcome.State != ValidationStateValid {
		t.Errorf("Expected state %q, got %q", ValidationStateValid, outcome.State)
	}
}

// -------------------------------------------------------------------------
// Functional Types Coverage (All 10 Types)
// -------------------------------------------------------------------------

func TestToolInput_AllFunctionalTypes_Conformance(t *testing.T) {
	contract := &ToolInputContract{
		ToolID:      "tool-all-types",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{Name: "f_string", FunctionalType: FunctionalTypeString, Required: true},
			{Name: "f_integer", FunctionalType: FunctionalTypeInteger, Required: true},
			{Name: "f_boolean", FunctionalType: FunctionalTypeBoolean, Required: true},
			{Name: "f_number", FunctionalType: FunctionalTypeNumber, Required: true},
			{Name: "f_array", FunctionalType: FunctionalTypeArray, Required: true},
			{Name: "f_object", FunctionalType: FunctionalTypeObject, Required: true},
			{Name: "f_secret", FunctionalType: FunctionalTypeSecretRef, Required: true},
			{Name: "f_entity", FunctionalType: FunctionalTypeEntityRef, Required: true},
			{Name: "f_query", FunctionalType: FunctionalTypeQueryAssetRef, Required: true},
			{Name: "f_enum", FunctionalType: FunctionalTypeEnum, Required: true, AllowedValues: []any{"A", "B", "C"}},
		},
	}

	scope := CallerScopeContext{TenantID: "tenant-alpha"}

	t.Run("Nominal Valid All Types", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "f_string", SourceType: BindingSourceLiteral, Value: "hello"},
			{ParameterName: "f_integer", SourceType: BindingSourceLiteral, Value: 42},
			{ParameterName: "f_boolean", SourceType: BindingSourceLiteral, Value: true},
			{ParameterName: "f_number", SourceType: BindingSourceLiteral, Value: 3.14159},
			{ParameterName: "f_array", SourceType: BindingSourceLiteral, Value: []string{"x", "y"}},
			{ParameterName: "f_object", SourceType: BindingSourceLiteral, Value: map[string]any{"k": "v"}},
			{ParameterName: "f_secret", SourceType: BindingSourceSecretRef, SecretRef: &SecretReference{RefID: "s-1", TenantID: "tenant-alpha"}},
			{ParameterName: "f_entity", SourceType: BindingSourceEntityRef, EntityRef: &EntityReference{RefID: "e-1", TenantID: "tenant-alpha"}},
			{ParameterName: "f_query", SourceType: BindingSourceAssetRef, QueryAssetRef: &QueryAssetReference{AssetID: "q-1", TenantID: "tenant-alpha"}},
			{ParameterName: "f_enum", SourceType: BindingSourceLiteral, Value: "B"},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-all-types-valid")
		if !outcome.Valid {
			t.Fatalf("Expected valid, got issues: %+v", outcome.Issues)
		}
		if len(outcome.ValidatedBindings) != 10 {
			t.Errorf("Expected 10 validated bindings, got %d", len(outcome.ValidatedBindings))
		}
	})

	t.Run("Type Mismatches for Every Type", func(t *testing.T) {
		mismatches := []struct {
			paramName string
			badValue  any
		}{
			{"f_string", 12345},
			{"f_integer", "not-an-integer"},
			{"f_boolean", "true"}, // String instead of bool
			{"f_number", "3.14"},  // String instead of float/number
			{"f_array", "not-a-slice"},
			{"f_object", 42},
			{"f_enum", "INVALID_CHOICE"},
		}

		for _, tc := range mismatches {
			t.Run(tc.paramName, func(t *testing.T) {
				singleContract := &ToolInputContract{
					ToolID:     "tool-single",
					TenantID:   "tenant-alpha",
					Parameters: []InputParameterDefinition{},
				}
				for _, p := range contract.Parameters {
					if p.Name == tc.paramName {
						singleContract.Parameters = append(singleContract.Parameters, p)
						break
					}
				}

				bindings := []InputBinding{
					{ParameterName: tc.paramName, SourceType: BindingSourceLiteral, Value: tc.badValue},
				}

				outcome := DeterministicValidateInputBindings(singleContract, bindings, scope, "corr-mismatch")
				if outcome.Valid {
					t.Errorf("Expected validation failure for %s with bad value %v, got valid", tc.paramName, tc.badValue)
				}
				if outcome.State != ValidationStateInvalid {
					t.Errorf("Expected state %q, got %q", ValidationStateInvalid, outcome.State)
				}
			})
		}
	})
}

// -------------------------------------------------------------------------
// Cross-Tenant Isolation Enforcement (REQ-SEC-001)
// -------------------------------------------------------------------------

func TestToolInput_CrossTenantIsolation_Enforced(t *testing.T) {
	contract := sampleToolInputContract("tenant-alpha", "tool-network-scanner")
	scope := CallerScopeContext{
		PrincipalID: "analyst-1",
		TenantID:    "tenant-alpha", // Active caller is tenant-alpha
	}

	t.Run("Cross-Tenant Entity Reference Blocked", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceSecretRef,
				SecretRef:     &SecretReference{RefID: "sec-01", TenantID: "tenant-alpha"},
			},
			{
				ParameterName: "incident_case",
				SourceType:    BindingSourceEntityRef,
				EntityRef: &EntityReference{
					RefID:    "case-9999",
					Kind:     "case",
					TenantID: "tenant-bravo", // CROSS-TENANT TARGET!
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-cross-entity")
		if outcome.Valid {
			t.Fatalf("Expected cross-tenant entity reference to fail, got Valid=true")
		}
		if outcome.State != ValidationStateCrossScopeBlocked {
			t.Errorf("Expected outcome state %q, got %q", ValidationStateCrossScopeBlocked, outcome.State)
		}

		foundIssue := false
		for _, issue := range outcome.Issues {
			if issue.ParameterName == "incident_case" && issue.Code == string(ErrorCodeCrossScopeBlocked) {
				foundIssue = true
				if issue.State != ValidationStateCrossScopeBlocked {
					t.Errorf("Expected issue state %q, got %q", ValidationStateCrossScopeBlocked, issue.State)
				}
			}
		}
		if !foundIssue {
			t.Errorf("Expected cross-scope-blocked issue for incident_case, got: %+v", outcome.Issues)
		}
	})

	t.Run("Cross-Tenant Secret Reference Blocked", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceSecretRef,
				SecretRef: &SecretReference{
					RefID:    "sec-foreign",
					TenantID: "tenant-charlie", // FOREIGN TENANT SECRET!
				},
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-cross-secret")
		if outcome.Valid {
			t.Fatalf("Expected cross-tenant secret reference to fail")
		}
		if outcome.State != ValidationStateCrossScopeBlocked {
			t.Errorf("Expected state %q, got %q", ValidationStateCrossScopeBlocked, outcome.State)
		}
	})

	t.Run("MSSP Authorized Tenant Set (ATS) Permits Authorized Tenancy", func(t *testing.T) {
		msspScope := CallerScopeContext{
			PrincipalID:       "mssp-operator",
			TenantID:          "tenant-alpha",
			AuthorizedTenants: []string{"tenant-alpha", "tenant-bravo", "tenant-charlie"},
		}

		bindings := []InputBinding{
			{ParameterName: "target_host", SourceType: BindingSourceLiteral, Value: "example.com"},
			{ParameterName: "port", SourceType: BindingSourceLiteral, Value: 80},
			{
				ParameterName: "api_key",
				SourceType:    BindingSourceSecretRef,
				SecretRef:     &SecretReference{RefID: "sec-mssp", TenantID: "tenant-bravo"}, // Allowed in ATS
			},
			{
				ParameterName: "incident_case",
				SourceType:    BindingSourceEntityRef,
				EntityRef:     &EntityReference{RefID: "case-mssp", Kind: "case", TenantID: "tenant-charlie"}, // Allowed in ATS
			},
		}

		outcome := DeterministicValidateInputBindings(contract, bindings, msspScope, "corr-mssp")
		if !outcome.Valid {
			t.Fatalf("Expected MSSP with authorized tenant set to succeed, got issues: %+v", outcome.Issues)
		}
	})
}

// -------------------------------------------------------------------------
// Constraints Validation (MinLength, MaxLength, Min, Max, Regex Pattern)
// -------------------------------------------------------------------------

func TestToolInput_ConstraintsValidation(t *testing.T) {
	minL := 5
	maxL := 10
	contract := &ToolInputContract{
		ToolID:   "tool-constraints",
		TenantID: "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{
				Name:           "username",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				MinLength:      &minL,
				MaxLength:      &maxL,
				Pattern:        "^[a-z]+$",
			},
		},
	}
	scope := CallerScopeContext{TenantID: "tenant-alpha"}

	t.Run("String Too Short", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "username", SourceType: BindingSourceLiteral, Value: "abc"}, // len 3 < 5
		}
		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-short")
		if outcome.Valid {
			t.Errorf("Expected failure for string too short")
		}
	})

	t.Run("String Too Long", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "username", SourceType: BindingSourceLiteral, Value: "thisiswaytoolong"}, // len 16 > 10
		}
		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-long")
		if outcome.Valid {
			t.Errorf("Expected failure for string too long")
		}
	})

	t.Run("Regex Pattern Mismatch", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "username", SourceType: BindingSourceLiteral, Value: "User123!"}, // uppercase & numbers fail ^[a-z]+$
		}
		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-regex")
		if outcome.Valid {
			t.Errorf("Expected failure for regex pattern violation")
		}
	})

	t.Run("Constraint Nominal Match", func(t *testing.T) {
		bindings := []InputBinding{
			{ParameterName: "username", SourceType: BindingSourceLiteral, Value: "validuser"}, // 9 chars, lowercase letters
		}
		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-match")
		if !outcome.Valid {
			t.Errorf("Expected nominal constraint match to pass, got: %+v", outcome.Issues)
		}
	})
}

// -------------------------------------------------------------------------
// Class 0: Inspect Tool Input Contract & Anti-Oracle 404
// -------------------------------------------------------------------------

func TestToolInput_Class0_InspectToolInputContract(t *testing.T) {
	store := NewInMemoryToolStore()
	ctx := context.Background()

	tool := &ToolDefinition{
		ID:       "net-tool-01",
		TenantID: "tenant-alpha",
		Name:     "Network Scanner",
		Purpose:  "Security network scanner",
		Version:  VersionInfo{SemVer: "1.0.0"},
		Parameters: []Parameter{
			{Name: "host", Type: "string", Required: true, Description: "Host"},
			{Name: "token", Type: "secret-ref", Required: true, Description: "Token", DefaultValue: "should-be-masked"},
		},
	}
	if err := store.Put(ctx, tool); err != nil {
		t.Fatalf("failed to put tool: %v", err)
	}

	service := NewDefaultToolInputService(store)
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	t.Run("Nominal Inspect", func(t *testing.T) {
		c, err := service.InspectToolInputContract(ctx, &auth, "tenant-alpha", "net-tool-01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.ToolID != "net-tool-01" {
			t.Errorf("expected tool ID 'net-tool-01', got %q", c.ToolID)
		}
		if len(c.Parameters) != 2 {
			t.Fatalf("expected 2 parameters, got %d", len(c.Parameters))
		}
		// Invariant: sensitive default value is masked as "***"
		for _, p := range c.Parameters {
			if p.Name == "token" {
				if p.DefaultValue != "***" {
					t.Errorf("CRITICAL LEAK: expected secret parameter default to be masked as '***', got %v", p.DefaultValue)
				}
			}
		}
	})

	t.Run("Anti-Oracle 404 on Missing or Cross-Tenant Tool", func(t *testing.T) {
		_, err := service.InspectToolInputContract(ctx, &auth, "tenant-alpha", "non-existent-tool")
		if err == nil {
			t.Fatalf("expected 404 error for non-existent tool")
		}
		var env *ErrorEnvelope
		if errors.As(err, &env) {
			if env.Code != ErrorCodeAssetNotFound {
				t.Errorf("expected ErrorCodeAssetNotFound, got %q", env.Code)
			}
			if env.HTTPStatus != 404 {
				t.Errorf("expected HTTP 404, got %d", env.HTTPStatus)
			}
		} else {
			t.Errorf("expected *ErrorEnvelope, got %T", err)
		}
	})

	t.Run("Permission Denied when Library Read is Missing", func(t *testing.T) {
		noAuth := NewAuthContext("user-no-perm", "tenant-alpha", nil)
		_, err := service.InspectToolInputContract(ctx, &noAuth, "tenant-alpha", "net-tool-01")
		if err == nil {
			t.Fatalf("expected permission denied error")
		}
		var env *ErrorEnvelope
		if errors.As(err, &env) {
			if env.Code != ErrorCodePermissionDenied {
				t.Errorf("expected ErrorCodePermissionDenied, got %q", env.Code)
			}
		}
	})
}

// -------------------------------------------------------------------------
// Class 2: Prepare Corrected Bindings (Draft authoring, reversible, no execution)
// -------------------------------------------------------------------------

func TestToolInput_Class2_PrepareCorrectedBindings(t *testing.T) {
	store := NewInMemoryToolStore()
	ctx := context.Background()
	contract := sampleToolInputContract("tenant-alpha", "net-tool-02")

	service := NewDefaultToolInputService(store)
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	// Initial incomplete bindings (missing required target_host and port)
	baseBindings := []InputBinding{
		{
			ParameterName: "api_key",
			SourceType:    BindingSourceSecretRef,
			SecretRef:     &SecretReference{RefID: "sec-1", TenantID: "tenant-alpha"},
		},
	}

	// Corrections providing the missing parameters
	corrections := []InputBinding{
		{
			ParameterName: "target_host",
			SourceType:    BindingSourceLiteral,
			Value:         "corrected.target.com",
		},
		{
			ParameterName: "port",
			SourceType:    BindingSourceLiteral,
			Value:         443,
		},
	}

	req := &PrepareCorrectedBindingsRequest{
		TenantRef:    "tenant-alpha",
		ToolID:       "net-tool-02",
		Contract:     contract,
		BaseBindings: baseBindings,
		Corrections:  corrections,
	}

	resp, err := service.PrepareCorrectedBindings(ctx, &auth, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsDraft {
		t.Errorf("expected IsDraft=true")
	}
	if !resp.Reversible {
		t.Errorf("expected Reversible=true")
	}
	if !resp.Outcome.Valid {
		t.Fatalf("expected corrected draft outcome to be valid, got issues: %+v", resp.Outcome.Issues)
	}
	if len(resp.DraftBindings) != 3 {
		t.Errorf("expected 3 merged draft bindings, got %d", len(resp.DraftBindings))
	}
}

// -------------------------------------------------------------------------
// Mapping Suggestions (Deterministic Heuristic with Manual Fallback)
// -------------------------------------------------------------------------

func TestToolInput_SuggestInputMapping_Deterministic(t *testing.T) {
	contract := sampleToolInputContract("tenant-alpha", "tool-suggest")
	store := NewInMemoryToolStore()
	service := NewDefaultToolInputService(store)
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	req := &SuggestInputMappingRequest{
		TenantRef: "tenant-alpha",
		ToolID:    "tool-suggest",
		Contract:  contract,
		AvailableFields: []string{
			"target_host", // Exact match
			"TARGETHOST",  // Case variant
			"Target-Port", // Normalized variant for port
			"random_unrelated_field",
		},
	}

	resp, err := service.SuggestInputMapping(context.Background(), &auth, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsDeterministic {
		t.Errorf("expected IsDeterministic=true")
	}
	if len(resp.Suggestions) == 0 {
		t.Errorf("expected suggestions, got none")
	}

	foundHostMatch := false
	for _, s := range resp.Suggestions {
		if s.ParameterName == "target_host" {
			foundHostMatch = true
			if s.Confidence < 0.8 {
				t.Errorf("expected high confidence for exact match, got %v", s.Confidence)
			}
		}
	}
	if !foundHostMatch {
		t.Errorf("expected suggestion for target_host")
	}
}

// -------------------------------------------------------------------------
// REQ-PROD-006: Permanent Execution Verbs Blocking (HTTP 400 execution-blocked)
// -------------------------------------------------------------------------

func TestToolInput_ExecutionVerbs_PermanentlyBlocked(t *testing.T) {
	service := NewDefaultToolInputService(nil)
	ctx := context.Background()
	auth := NewAuthContext("user-1", "tenant-alpha", nil, PermStudioLibraryRead)

	verbs := []struct {
		name string
		call func() error
	}{
		{"Invoke", func() error { return service.Invoke(ctx, &auth, "tool-1", map[string]any{}) }},
		{"Execute", func() error { return service.Execute(ctx, &auth, "tool-1", map[string]any{}) }},
		{"Run", func() error { return service.Run(ctx, &auth, "tool-1", map[string]any{}) }},
		{"Call", func() error { return service.Call(ctx, &auth, "tool-1", map[string]any{}) }},
		{"Trigger", func() error { return service.Trigger(ctx, &auth, "tool-1", map[string]any{}) }},
	}

	for _, tc := range verbs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatalf("CRITICAL REQ-PROD-006 BREACH: verb %s did not return error", tc.name)
			}

			var env *ErrorEnvelope
			if !errors.As(err, &env) {
				t.Fatalf("expected *ErrorEnvelope, got %T (%v)", err, err)
			}

			if env.Code != ErrorCodeExecutionBlocked {
				t.Errorf("expected error code %q, got %q", ErrorCodeExecutionBlocked, env.Code)
			}
			if env.HTTPStatus != 400 {
				t.Errorf("expected HTTP 400, got %d", env.HTTPStatus)
			}
			if !errors.Is(err, ErrExecutionBlocked) {
				t.Errorf("expected errors.Is(err, ErrExecutionBlocked) to be true")
			}
		})
	}
}

// -------------------------------------------------------------------------
// Remediation Verification Tests: AC-1, Sentinel Errors, Bounds, and Sorting
// -------------------------------------------------------------------------

func TestToolInput_AC1_RequiredWithDefaultStillFailsWhenMissing(t *testing.T) {
	scope := CallerScopeContext{PrincipalID: "p1", TenantID: "tenant-alpha"}

	// 1. Missing required parameter with default must FAIL closed (CAP-STD-004 §2, §25 AC-1)
	singleReqContract := &ToolInputContract{
		ToolID:      "test-tool",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{
				Name:           "req_with_default",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				DefaultValue:   "default-val",
			},
		},
	}

	outcome := DeterministicValidateInputBindings(singleReqContract, []InputBinding{}, scope, "corr-req-def-1")
	if outcome.Valid {
		t.Fatalf("CRITICAL BREACH: Expected validation to fail for missing required parameter with default value")
	}
	if outcome.State != ValidationStateMissingInput {
		t.Errorf("Expected state missing-input, got %s", outcome.State)
	}
	if len(outcome.Issues) != 1 || outcome.Issues[0].Code != string(ErrorCodeMissingInput) {
		t.Fatalf("Expected ErrorCodeMissingInput, got issues: %+v", outcome.Issues)
	}
	if len(outcome.ValidatedBindings) != 0 {
		t.Errorf("Expected zero validated bindings, got: %+v", outcome.ValidatedBindings)
	}

	// 2. Providing required parameter succeeds and populates optional default
	twoParamContract := &ToolInputContract{
		ToolID:      "test-tool",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{
				Name:           "req_with_default",
				FunctionalType: FunctionalTypeString,
				Required:       true,
				DefaultValue:   "default-val",
			},
			{
				Name:           "opt_with_default",
				FunctionalType: FunctionalTypeString,
				Required:       false,
				DefaultValue:   "opt-default-val",
			},
		},
	}

	outcomeWithReq := DeterministicValidateInputBindings(twoParamContract, []InputBinding{
		{ParameterName: "req_with_default", SourceType: BindingSourceLiteral, Value: "user-supplied"},
	}, scope, "corr-req-def-2")
	if !outcomeWithReq.Valid {
		t.Fatalf("Expected valid outcome, got issues: %+v", outcomeWithReq.Issues)
	}
	if len(outcomeWithReq.ValidatedBindings) != 2 {
		t.Fatalf("Expected 2 validated bindings, got %d", len(outcomeWithReq.ValidatedBindings))
	}
	var foundReq, foundOpt bool
	for _, vb := range outcomeWithReq.ValidatedBindings {
		if vb.ParameterName == "req_with_default" && vb.MaskedValue == "user-supplied" {
			foundReq = true
		}
		if vb.ParameterName == "opt_with_default" && vb.MaskedValue == "opt-default-val" {
			foundOpt = true
		}
	}
	if !foundReq || !foundOpt {
		t.Errorf("Expected both required user-supplied and optional default, got: %+v", outcomeWithReq.ValidatedBindings)
	}
}

func TestErrorEnvelope_Is_SentinelAndEnvelopeComparisons(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{"MissingInput Sentinel", NewMissingInputError("f1", "c1"), ErrMissingInput, true},
		{"InvalidInput Sentinel", NewInvalidInputError("f1", "msg", "c1"), ErrInvalidInput, true},
		{"RestrictedInput Sentinel", NewRestrictedInputError("f1", "reason", "c1"), ErrRestrictedInput, true},
		{"CrossScopeBlocked Sentinel", NewCrossScopeBlockedError("f1", "tenant-b", "tenant-a", "c1"), ErrCrossScopeBlocked, true},
		{"StaleReference Sentinel", NewStaleReferenceError("f1", "c1"), ErrStaleReference, true},
		{"UnsupportedFunctionalType Sentinel", NewUnsupportedFunctionalTypeError("custom", "c1"), ErrUnsupportedFunctionalType, true},
		{"ExecutionBlocked Sentinel", NewExecutionBlockedError("Invoke", "c1"), ErrExecutionBlocked, true},
		{"Envelope to Envelope Match", NewMissingInputError("f1", "c1"), NewMissingInputError("f2", "c2"), true},
		{"Envelope to Envelope Mismatch", NewMissingInputError("f1", "c1"), NewInvalidInputError("f1", "msg", "c2"), false},
		{"Foreign error mismatch", NewMissingInputError("f1", "c1"), errors.New("missing-input-foreign"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.expected {
				t.Errorf("errors.Is(%v, %v) = %v; expected %v", tc.err, tc.target, got, tc.expected)
			}
		})
	}
}

func TestToolInput_NumericalBounds_OverflowAndInfNaN(t *testing.T) {
	scope := CallerScopeContext{PrincipalID: "p1", TenantID: "tenant-alpha"}

	t.Run("Integer Overflow Rejection", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-num",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "int_param", FunctionalType: FunctionalTypeInteger, Required: true},
			},
		}

		overflowValues := []any{
			1e25,
			-1e25,
			uint64(math.MaxUint64),
			math.NaN(),
			math.Inf(1),
			math.Inf(-1),
		}

		for _, ov := range overflowValues {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "int_param", SourceType: BindingSourceLiteral, Value: ov},
			}, scope, "corr-overflow")
			if outcome.Valid {
				t.Errorf("Expected integer validation to reject value %v (%T), but passed", ov, ov)
			}
		}
	})

	t.Run("Integer Full Go Types Support", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-num",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "int_param", FunctionalType: FunctionalTypeInteger, Required: true},
			},
		}

		validInts := []any{
			int8(12),
			int16(300),
			int32(70000),
			int64(9000000000),
			uint8(250),
			uint16(65000),
			uint32(1000000),
			uint(42),
			float64(42.0),
		}

		for _, vi := range validInts {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "int_param", SourceType: BindingSourceLiteral, Value: vi},
			}, scope, "corr-valid-int")
			if !outcome.Valid {
				t.Errorf("Expected integer validation to accept %v (%T), got: %+v", vi, vi, outcome.Issues)
			}
		}
	})

	t.Run("Number Rejection of NaN and Inf", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-num",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "num_param", FunctionalType: FunctionalTypeNumber, Required: true},
			},
		}

		invalidFloats := []any{
			math.NaN(),
			math.Inf(1),
			math.Inf(-1),
			float32(math.NaN()),
		}

		for _, inv := range invalidFloats {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "num_param", SourceType: BindingSourceLiteral, Value: inv},
			}, scope, "corr-num-nan")
			if outcome.Valid {
				t.Errorf("Expected number validation to reject %v (%T), but passed", inv, inv)
			}
		}
	})

	t.Run("Number Full Go Types Support", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-num",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "num_param", FunctionalType: FunctionalTypeNumber, Required: true},
			},
		}

		validNumbers := []any{
			float32(3.14),
			float64(2.71828),
			int8(10),
			int16(200),
			int32(3000),
			int64(40000),
			uint8(5),
			uint16(50),
			uint32(500),
			uint64(5000),
		}

		for _, vn := range validNumbers {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "num_param", SourceType: BindingSourceLiteral, Value: vn},
			}, scope, "corr-valid-num")
			if !outcome.Valid {
				t.Errorf("Expected number validation to accept %v (%T), got: %+v", vn, vn, outcome.Issues)
			}
		}
	})
}

func TestToolInput_MaskedParameters_DeterministicOrder(t *testing.T) {
	contract := &ToolInputContract{
		ToolID:      "tool-mask",
		ToolVersion: "1.0.0",
		TenantID:    "tenant-alpha",
		Parameters: []InputParameterDefinition{
			{Name: "z_secret", FunctionalType: FunctionalTypeSecretRef, Required: true},
			{Name: "a_token", FunctionalType: FunctionalTypeSecretRef, Required: true},
			{Name: "m_passphrase", FunctionalType: FunctionalTypeString, Required: true, IsSensitive: true},
			{Name: "d_key", FunctionalType: FunctionalTypeSecretRef, Required: true},
		},
	}
	scope := CallerScopeContext{PrincipalID: "p1", TenantID: "tenant-alpha"}

	bindings := []InputBinding{
		{ParameterName: "z_secret", SourceType: BindingSourceSecretRef, SecretRef: &SecretReference{RefID: "s1"}},
		{ParameterName: "a_token", SourceType: BindingSourceSecretRef, SecretRef: &SecretReference{RefID: "s2"}},
		{ParameterName: "m_passphrase", SourceType: BindingSourceLiteral, Value: "supersecret"},
		{ParameterName: "d_key", SourceType: BindingSourceSecretRef, SecretRef: &SecretReference{RefID: "s3"}},
	}

	expectedOrder := []string{"a_token", "d_key", "m_passphrase", "z_secret"}

	for i := 0; i < 50; i++ {
		outcome := DeterministicValidateInputBindings(contract, bindings, scope, "corr-order")
		if !outcome.Valid {
			t.Fatalf("Expected valid outcome, got: %+v", outcome.Issues)
		}
		if !reflect.DeepEqual(outcome.MaskedParameters, expectedOrder) {
			t.Fatalf("Iteration %d: Non-deterministic order! Expected %v, got %v", i, expectedOrder, outcome.MaskedParameters)
		}
	}
}

func TestToolInput_ExclusiveBoundsAndPointerObject(t *testing.T) {
	scope := CallerScopeContext{PrincipalID: "p1", TenantID: "tenant-alpha"}
	exMin := 10.0
	exMax := 100.0

	t.Run("ExclusiveMinimum and ExclusiveMaximum on Integer", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-bounds",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{
					Name:           "bounded_int",
					FunctionalType: FunctionalTypeInteger,
					Required:       true,
					Constraints: &InputConstraint{
						ExclusiveMinimum: &exMin,
						ExclusiveMaximum: &exMax,
					},
				},
			},
		}

		// Boundary tests: 10 fails, 100 fails, 9 fails, 101 fails, 11 passes, 99 passes
		failValues := []int{10, 100, 9, 101}
		for _, fv := range failValues {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "bounded_int", SourceType: BindingSourceLiteral, Value: fv},
			}, scope, "corr-bounds-int-fail")
			if outcome.Valid {
				t.Errorf("Expected integer %d to fail exclusive bounds (10, 100)", fv)
			}
		}

		passValues := []int{11, 50, 99}
		for _, pv := range passValues {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "bounded_int", SourceType: BindingSourceLiteral, Value: pv},
			}, scope, "corr-bounds-int-pass")
			if !outcome.Valid {
				t.Errorf("Expected integer %d to pass exclusive bounds (10, 100), got: %+v", pv, outcome.Issues)
			}
		}
	})

	t.Run("ExclusiveMinimum and ExclusiveMaximum on Number", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-bounds-num",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{
					Name:           "bounded_num",
					FunctionalType: FunctionalTypeNumber,
					Required:       true,
					Constraints: &InputConstraint{
						ExclusiveMinimum: &exMin,
						ExclusiveMaximum: &exMax,
					},
				},
			},
		}

		// Boundary tests: 10.0 fails, 100.0 fails, 10.0001 passes, 99.999 passes
		failValues := []float64{10.0, 100.0, 9.999, 100.001}
		for _, fv := range failValues {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "bounded_num", SourceType: BindingSourceLiteral, Value: fv},
			}, scope, "corr-bounds-num-fail")
			if outcome.Valid {
				t.Errorf("Expected number %f to fail exclusive bounds (10, 100)", fv)
			}
		}

		passValues := []float64{10.0001, 50.5, 99.999}
		for _, pv := range passValues {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "bounded_num", SourceType: BindingSourceLiteral, Value: pv},
			}, scope, "corr-bounds-num-pass")
			if !outcome.Valid {
				t.Errorf("Expected number %f to pass exclusive bounds (10, 100), got: %+v", pv, outcome.Issues)
			}
		}
	})

	t.Run("FunctionalTypeObject Pointer Dereference", func(t *testing.T) {
		contract := &ToolInputContract{
			ToolID:      "tool-obj",
			ToolVersion: "1.0.0",
			TenantID:    "tenant-alpha",
			Parameters: []InputParameterDefinition{
				{Name: "obj_param", FunctionalType: FunctionalTypeObject, Required: true},
			},
		}

		type SampleStruct struct {
			Field string
		}

		// Valid objects: struct value, map value, struct pointer, map pointer
		validObjs := []any{
			SampleStruct{Field: "direct"},
			&SampleStruct{Field: "pointer"},
			map[string]any{"key": "val"},
			&map[string]any{"key": "val"},
		}

		for _, vo := range validObjs {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "obj_param", SourceType: BindingSourceLiteral, Value: vo},
			}, scope, "corr-obj-pass")
			if !outcome.Valid {
				t.Errorf("Expected object %T to pass validation, got: %+v", vo, outcome.Issues)
			}
		}

		// Invalid objects: pointer to scalar, nil pointer, string
		intVal := 42
		invalidObjs := []any{
			&intVal,
			(*SampleStruct)(nil),
			"not-an-object",
		}

		for _, io := range invalidObjs {
			outcome := DeterministicValidateInputBindings(contract, []InputBinding{
				{ParameterName: "obj_param", SourceType: BindingSourceLiteral, Value: io},
			}, scope, "corr-obj-fail")
			if outcome.Valid {
				t.Errorf("Expected invalid object %T (%v) to fail validation", io, io)
			}
		}
	})
}
