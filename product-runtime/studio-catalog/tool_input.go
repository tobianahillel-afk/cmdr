// Package studiocatalog implements the runtime contract for CMDR Studio Tool Inputs,
// Parameters, Bindings, and Deterministic Validation conforming to CAP-STD-004,
// REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, and REQ-SEC-001.
package studiocatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// -------------------------------------------------------------------------
// Functional Types and Constants (CAP-STD-004 §8, §11)
// -------------------------------------------------------------------------

// FunctionalType defines the recognized functional types for tool parameters.
type FunctionalType string

const (
	FunctionalTypeString        FunctionalType = "string"
	FunctionalTypeInteger       FunctionalType = "integer"
	FunctionalTypeBoolean       FunctionalType = "boolean"
	FunctionalTypeNumber        FunctionalType = "number"
	FunctionalTypeArray         FunctionalType = "array"
	FunctionalTypeObject        FunctionalType = "object"
	FunctionalTypeSecretRef     FunctionalType = "secret-ref"
	FunctionalTypeEntityRef     FunctionalType = "entity-ref"
	FunctionalTypeQueryAssetRef FunctionalType = "query-asset-ref"
	FunctionalTypeEnum          FunctionalType = "enum"
)

// IsValid checks if the functional type is among the canonical recognized types.
func (f FunctionalType) IsValid() bool {
	switch f {
	case FunctionalTypeString,
		FunctionalTypeInteger,
		FunctionalTypeBoolean,
		FunctionalTypeNumber,
		FunctionalTypeArray,
		FunctionalTypeObject,
		FunctionalTypeSecretRef,
		FunctionalTypeEntityRef,
		FunctionalTypeQueryAssetRef,
		FunctionalTypeEnum:
		return true
	default:
		return false
	}
}

// NormalizeFunctionalType handles common variations into canonical types.
func NormalizeFunctionalType(raw string) FunctionalType {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	switch normalized {
	case "string", "str", "text":
		return FunctionalTypeString
	case "int", "integer":
		return FunctionalTypeInteger
	case "bool", "boolean":
		return FunctionalTypeBoolean
	case "number", "float", "double":
		return FunctionalTypeNumber
	case "array", "list", "slice":
		return FunctionalTypeArray
	case "object", "map", "dict", "struct":
		return FunctionalTypeObject
	case "secret-ref", "secret_ref", "secret_reference", "secretref":
		return FunctionalTypeSecretRef
	case "entity-ref", "entity_ref", "entityref", "ref":
		return FunctionalTypeEntityRef
	case "query-asset-ref", "query_asset_ref", "asset-ref", "asset_ref":
		return FunctionalTypeQueryAssetRef
	case "enum", "enumeration":
		return FunctionalTypeEnum
	default:
		return FunctionalType(raw)
	}
}

// -------------------------------------------------------------------------
// Validation State Machine (CAP-STD-004 §14)
// -------------------------------------------------------------------------

// ValidationState represents the 6 canonical functional states.
type ValidationState string

const (
	ValidationStateUnvalidated       ValidationState = "unvalidated"
	ValidationStateValid             ValidationState = "valid"
	ValidationStateInvalid           ValidationState = "invalid"
	ValidationStateMissingInput      ValidationState = "missing-input"
	ValidationStateRestrictedInput   ValidationState = "restricted-input"
	ValidationStateCrossScopeBlocked ValidationState = "cross-scope-blocked"
)

// IsValid checks if the validation state is valid.
func (s ValidationState) IsValid() bool {
	switch s {
	case ValidationStateUnvalidated,
		ValidationStateValid,
		ValidationStateInvalid,
		ValidationStateMissingInput,
		ValidationStateRestrictedInput,
		ValidationStateCrossScopeBlocked:
		return true
	default:
		return false
	}
}

// -------------------------------------------------------------------------
// Binding Source Types
// -------------------------------------------------------------------------

// BindingSourceType specifies the origin of an input binding.
type BindingSourceType string

const (
	BindingSourceLiteral   BindingSourceType = "literal"
	BindingSourceSecretRef BindingSourceType = "secret-ref"
	BindingSourceEntityRef BindingSourceType = "entity-ref"
	BindingSourceAssetRef  BindingSourceType = "asset-ref"
)

// -------------------------------------------------------------------------
// References & Policy Structures (REQ-OBJ-009, REQ-SEC-001)
// -------------------------------------------------------------------------

// SecretReference represents an opaque reference to a Platform Settings secret (CAP-SET-009).
// In CMDR Studio, raw secret values are never accepted, stored, or revealed.
type SecretReference struct {
	RefID         string `json:"ref_id"`
	Owner         string `json:"owner"` // Strictly "Settings" (CAP-SET-009)
	TenantID      string `json:"tenant_id"`
	EnvironmentID string `json:"environment_id,omitempty"`
	SecretType    string `json:"secret_type,omitempty"`
	IsMasked      bool   `json:"is_masked"`    // Always true
	MaskedValue   string `json:"masked_value"` // Always "***"
	IsStale       bool   `json:"is_stale,omitempty"`
}

// Mask ensures secret reference properties conform strictly to masking requirements.
func (s *SecretReference) Mask() {
	if s == nil {
		return
	}
	s.Owner = "Settings"
	s.IsMasked = true
	s.MaskedValue = "***"
}

// NewSecretReference creates a properly masked SecretReference.
func NewSecretReference(refID, tenantID, envID, secretType string) *SecretReference {
	return &SecretReference{
		RefID:         refID,
		Owner:         "Settings",
		TenantID:      tenantID,
		EnvironmentID: envID,
		SecretType:    secretType,
		IsMasked:      true,
		MaskedValue:   "***",
		IsStale:       false,
	}
}

// SecretRefPolicy defines policy constraints on secret references.
type SecretRefPolicy struct {
	AllowedSecretTypes []string `json:"allowed_secret_types,omitempty"`
	RequiredScope      string   `json:"required_scope,omitempty"`
	AllowPlaintext     bool     `json:"allow_plaintext"` // ALWAYS false; true triggers validation rejection
}

// EntityReference represents a reference to a platform entity (e.g. case, host, alert).
type EntityReference struct {
	RefID         string `json:"ref_id"`
	Kind          string `json:"kind"`
	TenantID      string `json:"tenant_id"`
	EnvironmentID string `json:"environment_id,omitempty"`
}

// QueryAssetReference represents a reference to a saved query asset.
type QueryAssetReference struct {
	AssetID  string `json:"asset_id"`
	TenantID string `json:"tenant_id"`
}

// EntityScope enforces tenant and kind isolation for references.
type EntityScope struct {
	ExpectedTenant string `json:"expected_tenant,omitempty"`
	TargetKind     string `json:"target_kind,omitempty"`
}

// InputConstraint defines detailed validation constraints.
type InputConstraint struct {
	MinLength        *int     `json:"min_length,omitempty"`
	MaxLength        *int     `json:"max_length,omitempty"`
	Pattern          string   `json:"pattern,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusive_minimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusive_maximum,omitempty"`
	MinItems         *int     `json:"min_items,omitempty"`
	MaxItems         *int     `json:"max_items,omitempty"`
	AllowedValues    []any    `json:"allowed_values,omitempty"`
}

// InputParameterDefinition defines a parameter specification under CAP-STD-004.
type InputParameterDefinition struct {
	Name            string           `json:"name"`
	FunctionalType  FunctionalType   `json:"functional_type"`
	Description     string           `json:"description,omitempty"`
	Required        bool             `json:"required"`
	AllowedValues   []any            `json:"allowed_values,omitempty"`
	EntityScope     *EntityScope     `json:"entity_scope,omitempty"`
	SecretRefPolicy *SecretRefPolicy `json:"secret_ref_policy,omitempty"`
	DefaultValue    any              `json:"default_value,omitempty"` // Explicit default; nil if none
	Constraints     *InputConstraint `json:"constraints,omitempty"`
	IsSensitive     bool             `json:"is_sensitive,omitempty"`
	MinLength       *int             `json:"min_length,omitempty"`
	MaxLength       *int             `json:"max_length,omitempty"`
	Minimum         *float64         `json:"minimum,omitempty"`
	Maximum         *float64         `json:"maximum,omitempty"`
	Pattern         string           `json:"pattern,omitempty"`
}

// Clone creates a deep copy of InputParameterDefinition.
func (p *InputParameterDefinition) Clone() InputParameterDefinition {
	cp := *p
	if p.AllowedValues != nil {
		cp.AllowedValues = make([]any, len(p.AllowedValues))
		copy(cp.AllowedValues, p.AllowedValues)
	}
	if p.EntityScope != nil {
		es := *p.EntityScope
		cp.EntityScope = &es
	}
	if p.SecretRefPolicy != nil {
		srp := *p.SecretRefPolicy
		if p.SecretRefPolicy.AllowedSecretTypes != nil {
			srp.AllowedSecretTypes = make([]string, len(p.SecretRefPolicy.AllowedSecretTypes))
			copy(srp.AllowedSecretTypes, p.SecretRefPolicy.AllowedSecretTypes)
		}
		cp.SecretRefPolicy = &srp
	}
	if p.Constraints != nil {
		c := *p.Constraints
		cp.Constraints = &c
	}
	return cp
}

// InputBinding represents an argument bound to a parameter by the caller.
type InputBinding struct {
	ParameterName   string               `json:"parameter_name"`
	SourceType      BindingSourceType    `json:"source_type"`
	Value           any                  `json:"value,omitempty"`
	SecretRef       *SecretReference     `json:"secret_ref,omitempty"`
	EntityRef       *EntityReference     `json:"entity_ref,omitempty"`
	QueryAssetRef   *QueryAssetReference `json:"query_asset_ref,omitempty"`
	ValidationState ValidationState      `json:"validation_state"`
	TenantID        string               `json:"tenant_id,omitempty"`
	EnvironmentID   string               `json:"environment_id,omitempty"`
}

// Clone creates a deep copy of an InputBinding.
func (b *InputBinding) Clone() InputBinding {
	cp := *b
	if b.SecretRef != nil {
		s := *b.SecretRef
		cp.SecretRef = &s
	}
	if b.EntityRef != nil {
		e := *b.EntityRef
		cp.EntityRef = &e
	}
	if b.QueryAssetRef != nil {
		q := *b.QueryAssetRef
		cp.QueryAssetRef = &q
	}
	return cp
}

// ToolInputContract specifies the collection of input parameters for a Tool or Operation.
type ToolInputContract struct {
	ToolID      string                     `json:"tool_id"`
	ToolVersion string                     `json:"tool_version"`
	TenantID    string                     `json:"tenant_id"`
	OperationID string                     `json:"operation_id,omitempty"`
	Parameters  []InputParameterDefinition `json:"parameters"`
}

// Clone creates a deep copy of ToolInputContract.
func (c *ToolInputContract) Clone() *ToolInputContract {
	if c == nil {
		return nil
	}
	cp := *c
	if c.Parameters != nil {
		cp.Parameters = make([]InputParameterDefinition, len(c.Parameters))
		for i, p := range c.Parameters {
			cp.Parameters[i] = p.Clone()
		}
	}
	return &cp
}

// -------------------------------------------------------------------------
// Validation Outcome & Issues
// -------------------------------------------------------------------------

// InputValidationIssue describes a specific validation problem detected on a parameter.
type InputValidationIssue struct {
	ParameterName string          `json:"parameter_name"`
	Code          string          `json:"code"`
	Message       string          `json:"message"`
	Severity      string          `json:"severity"` // "error", "warning", "info"
	State         ValidationState `json:"state"`
}

// ValidatedBinding represents a safe, validated, and normalized parameter binding.
type ValidatedBinding struct {
	ParameterName   string            `json:"parameter_name"`
	FunctionalType  FunctionalType    `json:"functional_type"`
	SourceType      BindingSourceType `json:"source_type"`
	MaskedValue     any               `json:"masked_value"`
	RefID           string            `json:"ref_id,omitempty"`
	Owner           string            `json:"owner,omitempty"` // Strictly "Settings" for secret-ref (CAP-SET-009, §25 AC-2)
	TenantID        string            `json:"tenant_id,omitempty"`
	ValidationState ValidationState   `json:"validation_state"`
	IsMasked        bool              `json:"is_masked"`
}

// InputValidationOutcome contains the complete deterministic assessment result.
type InputValidationOutcome struct {
	State                  ValidationState        `json:"state"`
	Valid                  bool                   `json:"valid"`
	ToolID                 string                 `json:"tool_id"`
	ToolVersion            string                 `json:"tool_version"`
	CallerTenantID         string                 `json:"caller_tenant_id"`
	Issues                 []InputValidationIssue `json:"issues"`
	ValidatedBindings      []ValidatedBinding     `json:"validated_bindings"`
	MaskedParameters       []string               `json:"masked_parameters"`
	DeterministicValidated bool                   `json:"deterministic_validated"` // Always true
	CorrelationID          string                 `json:"correlation_id"`
	ValidatedAt            time.Time              `json:"validated_at"`
}

// CallerScopeContext provides caller tenant identity and MSSP permissions context.
type CallerScopeContext struct {
	PrincipalID       string   `json:"principal_id"`
	TenantID          string   `json:"tenant_id"`
	EnvironmentID     string   `json:"environment_id,omitempty"`
	AuthorizedTenants []string `json:"authorized_tenants,omitempty"`
}

// ScopeFromAuth builds CallerScopeContext from an AuthContext.
func ScopeFromAuth(auth *AuthContext) CallerScopeContext {
	if auth == nil {
		return CallerScopeContext{}
	}
	return CallerScopeContext{
		PrincipalID:       auth.PrincipalID,
		TenantID:          auth.TenantID,
		AuthorizedTenants: auth.AuthorizedTenantSet,
	}
}

// -------------------------------------------------------------------------
// Deterministic Validation Engine (CAP-STD-004 §11, §13, §25)
// -------------------------------------------------------------------------

// DeterministicValidateInputBindings executes 100% local, fail-closed deterministic validation.
// It never relies on AI models and strictly obeys the invariant:
// NEVER fabricate or invent default values for missing required inputs.
func DeterministicValidateInputBindings(
	contract *ToolInputContract,
	bindings []InputBinding,
	scope CallerScopeContext,
	correlationID string,
) *InputValidationOutcome {
	if correlationID == "" {
		correlationID = GenerateCorrelationID()
	}

	outcome := &InputValidationOutcome{
		State:                  ValidationStateUnvalidated,
		Valid:                  false,
		Issues:                 []InputValidationIssue{},
		ValidatedBindings:      []ValidatedBinding{},
		MaskedParameters:       []string{},
		DeterministicValidated: true,
		CorrelationID:          correlationID,
		ValidatedAt:            time.Now().UTC(),
	}

	if contract != nil {
		outcome.ToolID = contract.ToolID
		outcome.ToolVersion = contract.ToolVersion
	}
	outcome.CallerTenantID = scope.TenantID

	// 1. Caller Tenant Scope Check
	trimmedCallerTenant := strings.TrimSpace(scope.TenantID)
	if trimmedCallerTenant == "" {
		outcome.State = ValidationStateInvalid
		outcome.Issues = append(outcome.Issues, InputValidationIssue{
			Code:     string(ErrorCodeMissingTenant),
			Message:  "Tenant identifier is required for input validation",
			Severity: "error",
			State:    ValidationStateInvalid,
		})
		return outcome
	}
	if trimmedCallerTenant == "*" || trimmedCallerTenant == "%" || strings.EqualFold(trimmedCallerTenant, "all") {
		outcome.State = ValidationStateInvalid
		outcome.Issues = append(outcome.Issues, InputValidationIssue{
			Code:     string(ErrorCodeWildcardTenant),
			Message:  "Wildcard tenant selectors are strictly forbidden in input validation",
			Severity: "error",
			State:    ValidationStateInvalid,
		})
		return outcome
	}

	// 2. Contract Presence Check
	if contract == nil {
		outcome.State = ValidationStateInvalid
		outcome.Issues = append(outcome.Issues, InputValidationIssue{
			Code:     "nil-contract",
			Message:  "Tool input contract cannot be nil",
			Severity: "error",
			State:    ValidationStateInvalid,
		})
		return outcome
	}

	// 3. Contract Tenant Isolation Check
	if contract.TenantID != "" && contract.TenantID != trimmedCallerTenant {
		if !containsTenant(scope.AuthorizedTenants, contract.TenantID) {
			outcome.State = ValidationStateCrossScopeBlocked
			outcome.Issues = append(outcome.Issues, InputValidationIssue{
				Code:     string(ErrorCodeCrossScopeBlocked),
				Message:  fmt.Sprintf("Tool contract tenant %q does not match caller tenant %q", contract.TenantID, trimmedCallerTenant),
				Severity: "error",
				State:    ValidationStateCrossScopeBlocked,
			})
			return outcome
		}
	}

	// Index caller bindings by parameter name (deep copy to preserve "no caller mutation" invariant)
	bindingMap := make(map[string]InputBinding, len(bindings))
	for _, b := range bindings {
		bindingMap[strings.TrimSpace(b.ParameterName)] = b.Clone()
	}

	// Track parameters processed
	maskedSet := make(map[string]bool)

	// Evaluate every parameter in the pinned contract
	for _, param := range contract.Parameters {
		paramName := strings.TrimSpace(param.Name)
		normType := NormalizeFunctionalType(string(param.FunctionalType))

		// Check functional type validity
		if !normType.IsValid() {
			outcome.Issues = append(outcome.Issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeUnsupportedFunctionalType),
				Message:       fmt.Sprintf("Functional type %q is unsupported under CAP-STD-004", param.FunctionalType),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			continue
		}

		// Security check: SecretRefPolicy cannot allow plaintext
		if param.SecretRefPolicy != nil && param.SecretRefPolicy.AllowPlaintext {
			outcome.Issues = append(outcome.Issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q policy cannot permit plaintext secrets under REQ-SEC-001", paramName),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			continue
		}

		binding, exists := bindingMap[paramName]

		// -------------------------------------------------------------
		// CASE A: Binding missing or empty
		// -------------------------------------------------------------
		if !exists || isBindingEmpty(binding) {
			if param.Required {
				// FAIL CLOSED: Required inputs must ALWAYS fail with missing-input when omitted.
				// NEVER inject or fall back to DefaultValue for required parameters! (CAP-STD-004 §2, §25 AC-1)
				outcome.Issues = append(outcome.Issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeMissingInput),
					Message:       fmt.Sprintf("Required parameter %q is missing (CAP-STD-004 §25)", paramName),
					Severity:      "error",
					State:         ValidationStateMissingInput,
				})
			} else if param.DefaultValue != nil {
				// Optional parameter with explicit default: populate default
				defVal := param.DefaultValue
				isMasked := param.IsSensitive || normType == FunctionalTypeSecretRef
				if isMasked {
					maskedSet[paramName] = true
					defVal = "***"
				}
				vb := ValidatedBinding{
					ParameterName:   paramName,
					FunctionalType:  normType,
					SourceType:      BindingSourceLiteral,
					MaskedValue:     defVal,
					TenantID:        trimmedCallerTenant,
					ValidationState: ValidationStateValid,
					IsMasked:        isMasked,
				}
				if normType == FunctionalTypeSecretRef {
					vb.Owner = "Settings"
				}
				outcome.ValidatedBindings = append(outcome.ValidatedBindings, vb)
			}
			continue
		}

		// -------------------------------------------------------------
		// CASE B: Binding is present -> Validate types, references, constraints
		// -------------------------------------------------------------
		paramIssues := validateBindingAgainstParameter(param, normType, binding, trimmedCallerTenant, scope.AuthorizedTenants)
		paramState := ValidationStateValid

		if len(paramIssues) > 0 {
			// Find most critical state for this parameter
			paramState = determineWorstState(paramIssues)
			outcome.Issues = append(outcome.Issues, paramIssues...)
		}

		// Prepare normalized validated binding
		valBinding := ValidatedBinding{
			ParameterName:   paramName,
			FunctionalType:  normType,
			SourceType:      binding.SourceType,
			TenantID:        trimmedCallerTenant,
			ValidationState: paramState,
		}

		if normType == FunctionalTypeSecretRef || param.IsSensitive {
			maskedSet[paramName] = true
			valBinding.IsMasked = true
			valBinding.MaskedValue = "***"
			if binding.SecretRef != nil {
				valBinding.RefID = binding.SecretRef.RefID
				valBinding.Owner = "Settings"
			}
		} else if binding.EntityRef != nil {
			valBinding.RefID = binding.EntityRef.RefID
			valBinding.MaskedValue = binding.EntityRef.RefID
		} else if binding.QueryAssetRef != nil {
			valBinding.RefID = binding.QueryAssetRef.AssetID
			valBinding.MaskedValue = binding.QueryAssetRef.AssetID
		} else {
			valBinding.MaskedValue = binding.Value
		}

		outcome.ValidatedBindings = append(outcome.ValidatedBindings, valBinding)
	}

	// Populate masked parameters list with deterministic sort order
	outcome.MaskedParameters = make([]string, 0, len(maskedSet))
	for p := range maskedSet {
		outcome.MaskedParameters = append(outcome.MaskedParameters, p)
	}
	sort.Strings(outcome.MaskedParameters)

	// 4. Overall Outcome State Calculation
	if len(outcome.Issues) == 0 {
		outcome.State = ValidationStateValid
		outcome.Valid = true
	} else {
		outcome.Valid = false
		outcome.State = determineOverallOutcomeState(outcome.Issues)
	}

	return outcome
}

// validateBindingAgainstParameter evaluates a single parameter binding against its definition.
func validateBindingAgainstParameter(
	param InputParameterDefinition,
	normType FunctionalType,
	binding InputBinding,
	callerTenant string,
	authorizedTenants []string,
) []InputValidationIssue {
	var issues []InputValidationIssue
	paramName := param.Name

	// 1. Secret Reference Handling (REQ-OBJ-009, REQ-SEC-001, CAP-SET-009)
	if normType == FunctionalTypeSecretRef {
		// Strict check: Literal or plaintext values are strictly forbidden!
		if binding.SourceType == BindingSourceLiteral || (binding.Value != nil && !isNilLiteral(binding.Value)) {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q requires an opaque Secret Reference from Platform Settings (CAP-SET-009); raw or literal secrets are strictly forbidden", paramName),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		if binding.SecretRef == nil || strings.TrimSpace(binding.SecretRef.RefID) == "" {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeMissingInput),
				Message:       fmt.Sprintf("Parameter %q requires a valid Secret Reference ID", paramName),
				Severity:      "error",
				State:         ValidationStateMissingInput,
			})
			return issues
		}

		sec := binding.SecretRef
		// Note: Under CAP-STD-004 §10 ("no caller mutation"), sec is never modified in-place.
		// Settings ownership and masking are projected strictly onto the derived ValidatedBinding.

		// Staleness check
		if sec.IsStale {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeStaleReference),
				Message:       fmt.Sprintf("Secret reference %q is stale or expired in Platform Settings", sec.RefID),
				Severity:      "error",
				State:         ValidationStateRestrictedInput,
			})
		}

		// Tenant scope check for secret reference
		if sec.TenantID != "" && sec.TenantID != callerTenant {
			if !containsTenant(authorizedTenants, sec.TenantID) {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeCrossScopeBlocked),
					Message:       fmt.Sprintf("Secret reference %q belongs to tenant %q, violating caller tenant %q (cross-scope-blocked)", sec.RefID, sec.TenantID, callerTenant),
					Severity:      "error",
					State:         ValidationStateCrossScopeBlocked,
				})
			}
		}

		// Allowed secret types policy check
		if param.SecretRefPolicy != nil && len(param.SecretRefPolicy.AllowedSecretTypes) > 0 && sec.SecretType != "" {
			allowed := false
			for _, t := range param.SecretRefPolicy.AllowedSecretTypes {
				if strings.EqualFold(t, sec.SecretType) {
					allowed = true
					break
				}
			}
			if !allowed {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeRestrictedInput),
					Message:       fmt.Sprintf("Secret type %q is not permitted by policy for parameter %q", sec.SecretType, paramName),
					Severity:      "error",
					State:         ValidationStateRestrictedInput,
				})
			}
		}

		return issues
	}

	// 2. Entity Reference Handling (`entity-ref`)
	if normType == FunctionalTypeEntityRef {
		if binding.EntityRef == nil || strings.TrimSpace(binding.EntityRef.RefID) == "" {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expects an entity-ref reference", paramName),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		ref := binding.EntityRef
		// Tenant isolation check
		if ref.TenantID != "" && ref.TenantID != callerTenant {
			if !containsTenant(authorizedTenants, ref.TenantID) {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeCrossScopeBlocked),
					Message:       fmt.Sprintf("Entity reference %q belongs to tenant %q, violating caller tenant %q (cross-scope-blocked)", ref.RefID, ref.TenantID, callerTenant),
					Severity:      "error",
					State:         ValidationStateCrossScopeBlocked,
				})
			}
		}

		if param.EntityScope != nil {
			if param.EntityScope.ExpectedTenant != "" && ref.TenantID != "" && param.EntityScope.ExpectedTenant != ref.TenantID {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeCrossScopeBlocked),
					Message:       fmt.Sprintf("Entity reference tenant %q does not match expected entity scope %q", ref.TenantID, param.EntityScope.ExpectedTenant),
					Severity:      "error",
					State:         ValidationStateCrossScopeBlocked,
				})
			}
			if param.EntityScope.TargetKind != "" && ref.Kind != "" && !strings.EqualFold(param.EntityScope.TargetKind, ref.Kind) {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeInvalidInput),
					Message:       fmt.Sprintf("Entity reference kind %q does not match expected kind %q", ref.Kind, param.EntityScope.TargetKind),
					Severity:      "error",
					State:         ValidationStateInvalid,
				})
			}
		}
		return issues
	}

	// 3. Query Asset Reference Handling (`query-asset-ref`)
	if normType == FunctionalTypeQueryAssetRef {
		var assetID, assetTenant string
		if binding.QueryAssetRef != nil && strings.TrimSpace(binding.QueryAssetRef.AssetID) != "" {
			assetID = binding.QueryAssetRef.AssetID
			assetTenant = binding.QueryAssetRef.TenantID
		} else if binding.EntityRef != nil && strings.TrimSpace(binding.EntityRef.RefID) != "" {
			assetID = binding.EntityRef.RefID
			assetTenant = binding.EntityRef.TenantID
		} else {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expects a query-asset-ref reference", paramName),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		if assetTenant != "" && assetTenant != callerTenant {
			if !containsTenant(authorizedTenants, assetTenant) {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeCrossScopeBlocked),
					Message:       fmt.Sprintf("Query asset reference %q belongs to tenant %q, violating caller tenant %q (cross-scope-blocked)", assetID, assetTenant, callerTenant),
					Severity:      "error",
					State:         ValidationStateCrossScopeBlocked,
				})
			}
		}
		return issues
	}

	// 4. Literal Value Type Checking & Constraints
	val := binding.Value

	// Check if secret-ref was incorrectly passed to literal parameter
	if binding.SecretRef != nil {
		issues = append(issues, InputValidationIssue{
			ParameterName: paramName,
			Code:          string(ErrorCodeInvalidInput),
			Message:       fmt.Sprintf("Parameter %q expects literal type %q, not a secret reference", paramName, normType),
			Severity:      "error",
			State:         ValidationStateInvalid,
		})
		return issues
	}

	switch normType {
	case FunctionalTypeString:
		strVal, ok := val.(string)
		if !ok {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected string, got %T", paramName, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		// Length checks
		minL := param.MinLength
		if minL == nil && param.Constraints != nil {
			minL = param.Constraints.MinLength
		}
		if minL != nil && len(strVal) < *minL {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q length %d is less than minimum %d", paramName, len(strVal), *minL),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		maxL := param.MaxLength
		if maxL == nil && param.Constraints != nil {
			maxL = param.Constraints.MaxLength
		}
		if maxL != nil && len(strVal) > *maxL {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q length %d exceeds maximum %d", paramName, len(strVal), *maxL),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		// Pattern check
		pat := param.Pattern
		if pat == "" && param.Constraints != nil {
			pat = param.Constraints.Pattern
		}
		if pat != "" {
			matched, err := regexp.MatchString(pat, strVal)
			if err != nil || !matched {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeInvalidInput),
					Message:       fmt.Sprintf("Parameter %q does not match required regex pattern %s", paramName, pat),
					Severity:      "error",
					State:         ValidationStateInvalid,
				})
			}
		}

		// Allowed values
		allowedVals := param.AllowedValues
		if len(allowedVals) == 0 && param.Constraints != nil {
			allowedVals = param.Constraints.AllowedValues
		}
		if len(allowedVals) > 0 {
			match := false
			for _, av := range allowedVals {
				if fmt.Sprintf("%v", av) == strVal {
					match = true
					break
				}
			}
			if !match {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeInvalidInput),
					Message:       fmt.Sprintf("Parameter %q value %q is not in allowed values", paramName, strVal),
					Severity:      "error",
					State:         ValidationStateInvalid,
				})
			}
		}

	case FunctionalTypeInteger:
		intVal, ok := toInt64(val)
		if !ok {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected integer, got %T (%v)", paramName, val, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		minV := param.Minimum
		if minV == nil && param.Constraints != nil {
			minV = param.Constraints.Minimum
		}
		if minV != nil && float64(intVal) < *minV {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %d is less than minimum %v", paramName, intVal, *minV),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		maxV := param.Maximum
		if maxV == nil && param.Constraints != nil {
			maxV = param.Constraints.Maximum
		}
		if maxV != nil && float64(intVal) > *maxV {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %d exceeds maximum %v", paramName, intVal, *maxV),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		if param.Constraints != nil && param.Constraints.ExclusiveMinimum != nil && float64(intVal) <= *param.Constraints.ExclusiveMinimum {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %d must be strictly greater than exclusive minimum %v", paramName, intVal, *param.Constraints.ExclusiveMinimum),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}
		if param.Constraints != nil && param.Constraints.ExclusiveMaximum != nil && float64(intVal) >= *param.Constraints.ExclusiveMaximum {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %d must be strictly less than exclusive maximum %v", paramName, intVal, *param.Constraints.ExclusiveMaximum),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

	case FunctionalTypeNumber:
		numVal, ok := toFloat64(val)
		if !ok {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected number, got %T (%v)", paramName, val, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}

		minV := param.Minimum
		if minV == nil && param.Constraints != nil {
			minV = param.Constraints.Minimum
		}
		if minV != nil && numVal < *minV {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %v is less than minimum %v", paramName, numVal, *minV),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		maxV := param.Maximum
		if maxV == nil && param.Constraints != nil {
			maxV = param.Constraints.Maximum
		}
		if maxV != nil && numVal > *maxV {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %v exceeds maximum %v", paramName, numVal, *maxV),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

		if param.Constraints != nil && param.Constraints.ExclusiveMinimum != nil && numVal <= *param.Constraints.ExclusiveMinimum {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %v must be strictly greater than exclusive minimum %v", paramName, numVal, *param.Constraints.ExclusiveMinimum),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}
		if param.Constraints != nil && param.Constraints.ExclusiveMaximum != nil && numVal >= *param.Constraints.ExclusiveMaximum {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %v must be strictly less than exclusive maximum %v", paramName, numVal, *param.Constraints.ExclusiveMaximum),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

	case FunctionalTypeBoolean:
		_, ok := val.(bool)
		if !ok {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected boolean, got %T (%v)", paramName, val, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

	case FunctionalTypeArray:
		rv := reflect.ValueOf(val)
		if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected array/slice, got %T", paramName, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
			return issues
		}
		if param.Constraints != nil {
			if param.Constraints.MinItems != nil && rv.Len() < *param.Constraints.MinItems {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeInvalidInput),
					Message:       fmt.Sprintf("Parameter %q items count %d is less than min_items %d", paramName, rv.Len(), *param.Constraints.MinItems),
					Severity:      "error",
					State:         ValidationStateInvalid,
				})
			}
			if param.Constraints.MaxItems != nil && rv.Len() > *param.Constraints.MaxItems {
				issues = append(issues, InputValidationIssue{
					ParameterName: paramName,
					Code:          string(ErrorCodeInvalidInput),
					Message:       fmt.Sprintf("Parameter %q items count %d exceeds max_items %d", paramName, rv.Len(), *param.Constraints.MaxItems),
					Severity:      "error",
					State:         ValidationStateInvalid,
				})
			}
		}

	case FunctionalTypeObject:
		rv := reflect.ValueOf(val)
		if rv.IsValid() && rv.Kind() == reflect.Pointer && !rv.IsNil() {
			rv = rv.Elem()
		}
		if !rv.IsValid() || (rv.Kind() != reflect.Map && rv.Kind() != reflect.Struct) {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q expected object/map, got %T", paramName, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}

	case FunctionalTypeEnum:
		allowedVals := param.AllowedValues
		if len(allowedVals) == 0 && param.Constraints != nil {
			allowedVals = param.Constraints.AllowedValues
		}
		valStr := fmt.Sprintf("%v", val)
		match := false
		for _, av := range allowedVals {
			if fmt.Sprintf("%v", av) == valStr {
				match = true
				break
			}
		}
		if !match {
			issues = append(issues, InputValidationIssue{
				ParameterName: paramName,
				Code:          string(ErrorCodeInvalidInput),
				Message:       fmt.Sprintf("Parameter %q value %v is not an allowed enum member", paramName, val),
				Severity:      "error",
				State:         ValidationStateInvalid,
			})
		}
	}

	return issues
}

// -------------------------------------------------------------------------
// Helper Functions for Validation
// -------------------------------------------------------------------------

func isBindingEmpty(b InputBinding) bool {
	if b.SecretRef != nil && strings.TrimSpace(b.SecretRef.RefID) != "" {
		return false
	}
	if b.EntityRef != nil && strings.TrimSpace(b.EntityRef.RefID) != "" {
		return false
	}
	if b.QueryAssetRef != nil && strings.TrimSpace(b.QueryAssetRef.AssetID) != "" {
		return false
	}
	return isNilLiteral(b.Value)
}

func isNilLiteral(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func containsTenant(set []string, target string) bool {
	tTarget := strings.TrimSpace(target)
	for _, item := range set {
		if strings.TrimSpace(item) == tTarget {
			return true
		}
	}
	return false
}

func determineWorstState(issues []InputValidationIssue) ValidationState {
	worst := ValidationStateInvalid
	for _, issue := range issues {
		if issue.State == ValidationStateCrossScopeBlocked {
			return ValidationStateCrossScopeBlocked
		}
		if issue.State == ValidationStateRestrictedInput {
			worst = ValidationStateRestrictedInput
		} else if issue.State == ValidationStateMissingInput && worst != ValidationStateRestrictedInput {
			worst = ValidationStateMissingInput
		}
	}
	return worst
}

func determineOverallOutcomeState(issues []InputValidationIssue) ValidationState {
	hasCrossScope := false
	hasRestricted := false
	hasMissing := false

	for _, issue := range issues {
		switch issue.State {
		case ValidationStateCrossScopeBlocked:
			hasCrossScope = true
		case ValidationStateRestrictedInput:
			hasRestricted = true
		case ValidationStateMissingInput:
			hasMissing = true
		}
	}

	if hasCrossScope {
		return ValidationStateCrossScopeBlocked
	}
	if hasRestricted {
		return ValidationStateRestrictedInput
	}
	if hasMissing {
		return ValidationStateMissingInput
	}
	return ValidationStateInvalid
}

func toInt64(val any) (int64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case float32:
		f := float64(v)
		if math.IsNaN(f) || math.IsInf(f, 0) || f < math.MinInt64 || f > math.MaxInt64 {
			return 0, false
		}
		if float32(int64(v)) == v {
			return int64(v), true
		}
		return 0, false
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) || v < math.MinInt64 || v > math.MaxInt64 {
			return 0, false
		}
		if math.Floor(v) == v {
			return int64(v), true
		}
		return 0, false
	case json.Number:
		i, err := v.Int64()
		return i, err == nil
	default:
		return 0, false
	}
}

func toFloat64(val any) (float64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, false
		}
		return v, true
	case float32:
		f := float64(v)
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return 0, false
		}
		return f, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// -------------------------------------------------------------------------
// Service Contracts & DefaultToolInputService (Class 0, 1, 2 & Permanent Execution Blocking)
// -------------------------------------------------------------------------

// ValidateInputBindingsRequest carries input bindings for deterministic validation.
type ValidateInputBindingsRequest struct {
	TenantRef     TenantRef          `json:"tenant_ref"`
	ToolID        string             `json:"tool_id"`
	ToolVersion   string             `json:"tool_version,omitempty"`
	Contract      *ToolInputContract `json:"contract,omitempty"`
	Bindings      []InputBinding     `json:"bindings"`
	CorrelationID string             `json:"correlation_id,omitempty"`
}

// PrepareCorrectedBindingsRequest carries base bindings and draft corrections.
type PrepareCorrectedBindingsRequest struct {
	TenantRef     TenantRef          `json:"tenant_ref"`
	ToolID        string             `json:"tool_id"`
	ToolVersion   string             `json:"tool_version,omitempty"`
	Contract      *ToolInputContract `json:"contract,omitempty"`
	BaseBindings  []InputBinding     `json:"base_bindings"`
	Corrections   []InputBinding     `json:"corrections"`
	CorrelationID string             `json:"correlation_id,omitempty"`
}

// PrepareCorrectedBindingsResponse is the result of Class 2 draft authoring.
type PrepareCorrectedBindingsResponse struct {
	DraftBindings []InputBinding          `json:"draft_bindings"`
	Outcome       *InputValidationOutcome `json:"outcome"`
	CorrelationID string                  `json:"correlation_id"`
	IsDraft       bool                    `json:"is_draft"`   // Always true
	Reversible    bool                    `json:"reversible"` // Always true
}

// SuggestInputMappingRequest requests deterministic mapping suggestions.
type SuggestInputMappingRequest struct {
	TenantRef       TenantRef          `json:"tenant_ref"`
	ToolID          string             `json:"tool_id"`
	Contract        *ToolInputContract `json:"contract,omitempty"`
	AvailableFields []string           `json:"available_fields"`
	CorrelationID   string             `json:"correlation_id,omitempty"`
}

// SuggestedMapping represents a suggested parameter association.
type SuggestedMapping struct {
	ParameterName  string         `json:"parameter_name"`
	SuggestedField string         `json:"suggested_field"`
	FunctionalType FunctionalType `json:"functional_type"`
	Confidence     float64        `json:"confidence"`
	Reason         string         `json:"reason"`
}

// SuggestInputMappingResponse returns mapping suggestions with deterministic fallback.
type SuggestInputMappingResponse struct {
	Suggestions         []SuggestedMapping `json:"suggestions"`
	DeterministicMethod string             `json:"deterministic_method"`
	IsDeterministic     bool               `json:"is_deterministic"` // Always true
	CorrelationID       string             `json:"correlation_id"`
}

// ToolInputService defines the canonical operations for CAP-STD-004.
type ToolInputService interface {
	// Class 0: Inspect Tool Input Contract
	InspectToolInputContract(ctx context.Context, auth *AuthContext, tenantRef TenantRef, toolID string) (*ToolInputContract, error)

	// Class 1: Deterministic validation of input bindings
	ValidateInputBindings(ctx context.Context, auth *AuthContext, req *ValidateInputBindingsRequest) (*InputValidationOutcome, error)

	// Class 2: Prepare corrected bindings (draft authoring, reversible, no execution)
	PrepareCorrectedBindings(ctx context.Context, auth *AuthContext, req *PrepareCorrectedBindingsRequest) (*PrepareCorrectedBindingsResponse, error)

	// SuggestInputMapping provides deterministic mapping suggestions (AI-optional, manual/deterministic alternative)
	SuggestInputMapping(ctx context.Context, auth *AuthContext, req *SuggestInputMappingRequest) (*SuggestInputMappingResponse, error)

	// REQ-PROD-006: Permanent Execution Blocking (all 5 verbs return 400 execution-blocked)
	Invoke(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Execute(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Run(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Call(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Trigger(ctx context.Context, auth *AuthContext, toolID string, payload any) error
}

// DefaultToolInputService implements ToolInputService.
type DefaultToolInputService struct {
	mu    sync.RWMutex
	store ToolStore
}

// NewDefaultToolInputService constructs a new DefaultToolInputService.
func NewDefaultToolInputService(store ToolStore) *DefaultToolInputService {
	return &DefaultToolInputService{
		store: store,
	}
}

// InspectToolInputContract provides Class 0 read-only inspection of tool parameter schemas.
func (s *DefaultToolInputService) InspectToolInputContract(
	ctx context.Context,
	auth *AuthContext,
	tenantRef TenantRef,
	toolID string,
) (*ToolInputContract, error) {
	corrID := GenerateCorrelationID()
	ac := ensureAuth(auth)

	// 1. Scoping enforcement (uniform 404 anti-oracle)
	if err := ValidateTenantScope(ac, string(tenantRef), corrID); err != nil {
		return nil, err
	}

	// 2. Permission check: dual namespace coexistence (perm.studio.library.read / perm.cmdr-studio.library.read)
	if !ac.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Fetch tool from store
	if s.store == nil {
		return nil, NewAssetNotFoundError(toolID, tenantRef, corrID)
	}
	tool, err := s.store.Get(ctx, tenantRef, toolID)
	if err != nil {
		return nil, err
	}
	if tool == nil || tool.TenantID != tenantRef {
		return nil, NewAssetNotFoundError(toolID, tenantRef, corrID)
	}

	// 4. Extract or project input contract
	var contract *ToolInputContract
	if tool.InputContract != nil {
		contract = tool.InputContract.Clone()
	} else {
		// Project from tool parameters
		contract = &ToolInputContract{
			ToolID:      tool.ID,
			ToolVersion: tool.Version.SemVer,
			TenantID:    tool.TenantID,
			Parameters:  make([]InputParameterDefinition, 0, len(tool.Parameters)),
		}
		for _, p := range tool.Parameters {
			contract.Parameters = append(contract.Parameters, InputParameterDefinition{
				Name:           p.Name,
				FunctionalType: NormalizeFunctionalType(p.Type),
				Required:       p.Required,
				Description:    p.Description,
				DefaultValue:   p.DefaultValue,
			})
		}
	}

	// 5. Mask sensitive defaults and ensure Settings ownership on secret references (REQ-OBJ-009, REQ-SEC-001)
	for i := range contract.Parameters {
		p := &contract.Parameters[i]
		if p.FunctionalType == FunctionalTypeSecretRef || p.IsSensitive {
			if p.DefaultValue != nil {
				p.DefaultValue = "***"
			}
		}
	}

	return contract, nil
}

// ValidateInputBindings performs Class 1 deterministic validation.
func (s *DefaultToolInputService) ValidateInputBindings(
	ctx context.Context,
	auth *AuthContext,
	req *ValidateInputBindingsRequest,
) (*InputValidationOutcome, error) {
	corrID := GenerateCorrelationID()
	if req != nil && req.CorrelationID != "" {
		corrID = req.CorrelationID
	}
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewInvalidInputError("request", "request payload cannot be nil", corrID)
	}

	if err := ValidateTenantScope(ac, string(req.TenantRef), corrID); err != nil {
		return nil, err
	}

	if !ac.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	contract := req.Contract
	if contract == nil {
		if strings.TrimSpace(req.ToolID) == "" {
			return nil, NewInvalidInputError("tool_id", "tool_id or contract is required", corrID)
		}
		if s.store == nil {
			return nil, NewAssetNotFoundError(req.ToolID, string(req.TenantRef), corrID)
		}
		tool, err := s.store.Get(ctx, string(req.TenantRef), req.ToolID)
		if err != nil {
			return nil, err
		}
		if tool == nil || tool.TenantID != string(req.TenantRef) {
			return nil, NewAssetNotFoundError(req.ToolID, string(req.TenantRef), corrID)
		}
		if tool.InputContract != nil {
			contract = tool.InputContract.Clone()
		} else {
			contract = &ToolInputContract{
				ToolID:      tool.ID,
				ToolVersion: tool.Version.SemVer,
				TenantID:    tool.TenantID,
				Parameters:  make([]InputParameterDefinition, 0, len(tool.Parameters)),
			}
			for _, p := range tool.Parameters {
				contract.Parameters = append(contract.Parameters, InputParameterDefinition{
					Name:           p.Name,
					FunctionalType: NormalizeFunctionalType(p.Type),
					Required:       p.Required,
					Description:    p.Description,
					DefaultValue:   p.DefaultValue,
				})
			}
		}
	} else {
		contract = contract.Clone()
	}

	if req.ToolID != "" {
		reqToolID := strings.TrimSpace(req.ToolID)
		if contract.ToolID != "" && reqToolID != contract.ToolID {
			return nil, NewInvalidInputError("tool_id", fmt.Sprintf("requested tool_id %q does not match contract tool_id %q", req.ToolID, contract.ToolID), corrID)
		}
		if contract.ToolID == "" {
			contract.ToolID = reqToolID
		}
	}

	if req.ToolVersion != "" {
		reqVer := strings.TrimSpace(req.ToolVersion)
		if contract.ToolVersion != "" && reqVer != contract.ToolVersion {
			return nil, NewInvalidInputError("tool_version", fmt.Sprintf("requested tool_version %q does not match contract tool_version %q", req.ToolVersion, contract.ToolVersion), corrID)
		}
		if contract.ToolVersion == "" {
			contract.ToolVersion = reqVer
		}
	}

	scope := ScopeFromAuth(&ac)
	scope.TenantID = string(req.TenantRef)
	outcome := DeterministicValidateInputBindings(contract, req.Bindings, scope, corrID)
	return outcome, nil
}

// PrepareCorrectedBindings performs Class 2 reversible draft binding updates.
func (s *DefaultToolInputService) PrepareCorrectedBindings(
	ctx context.Context,
	auth *AuthContext,
	req *PrepareCorrectedBindingsRequest,
) (*PrepareCorrectedBindingsResponse, error) {
	corrID := GenerateCorrelationID()
	if req != nil && req.CorrelationID != "" {
		corrID = req.CorrelationID
	}
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewInvalidInputError("request", "request payload cannot be nil", corrID)
	}

	if err := ValidateTenantScope(ac, string(req.TenantRef), corrID); err != nil {
		return nil, err
	}

	if !ac.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// Merge corrections onto base bindings
	bindingMap := make(map[string]InputBinding)
	var order []string

	for _, b := range req.BaseBindings {
		name := strings.TrimSpace(b.ParameterName)
		bindingMap[name] = b.Clone()
		order = append(order, name)
	}

	for _, c := range req.Corrections {
		name := strings.TrimSpace(c.ParameterName)
		if _, exists := bindingMap[name]; !exists {
			order = append(order, name)
		}
		bindingMap[name] = c.Clone()
	}

	mergedBindings := make([]InputBinding, 0, len(order))
	for _, name := range order {
		mergedBindings = append(mergedBindings, bindingMap[name])
	}

	// Validate the draft bindings
	vReq := &ValidateInputBindingsRequest{
		TenantRef:     req.TenantRef,
		ToolID:        req.ToolID,
		ToolVersion:   req.ToolVersion,
		Contract:      req.Contract,
		Bindings:      mergedBindings,
		CorrelationID: corrID,
	}

	outcome, err := s.ValidateInputBindings(ctx, auth, vReq)
	if err != nil {
		return nil, err
	}

	return &PrepareCorrectedBindingsResponse{
		DraftBindings: mergedBindings,
		Outcome:       outcome,
		CorrelationID: corrID,
		IsDraft:       true,
		Reversible:    true,
	}, nil
}

// SuggestInputMapping provides deterministic parameter mapping suggestions with manual fallback.
func (s *DefaultToolInputService) SuggestInputMapping(
	ctx context.Context,
	auth *AuthContext,
	req *SuggestInputMappingRequest,
) (*SuggestInputMappingResponse, error) {
	corrID := GenerateCorrelationID()
	if req != nil && req.CorrelationID != "" {
		corrID = req.CorrelationID
	}
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewInvalidInputError("request", "request payload cannot be nil", corrID)
	}

	if err := ValidateTenantScope(ac, string(req.TenantRef), corrID); err != nil {
		return nil, err
	}

	if !ac.HasPermission(PermStudioLibraryRead) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	var params []InputParameterDefinition
	if req.Contract != nil {
		params = req.Contract.Parameters
	} else if req.ToolID != "" && s.store != nil {
		tool, err := s.store.Get(ctx, req.TenantRef, req.ToolID)
		if err == nil && tool != nil {
			if tool.InputContract != nil {
				params = tool.InputContract.Parameters
			} else {
				for _, p := range tool.Parameters {
					params = append(params, InputParameterDefinition{
						Name:           p.Name,
						FunctionalType: NormalizeFunctionalType(p.Type),
						Required:       p.Required,
					})
				}
			}
		}
	}

	var suggestions []SuggestedMapping
	cleanStr := func(s string) string {
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "_", "")
		s = strings.ReplaceAll(s, "-", "")
		return s
	}

	for _, p := range params {
		pClean := cleanStr(p.Name)
		for _, f := range req.AvailableFields {
			fClean := cleanStr(f)
			if strings.EqualFold(p.Name, f) {
				suggestions = append(suggestions, SuggestedMapping{
					ParameterName:  p.Name,
					SuggestedField: f,
					FunctionalType: p.FunctionalType,
					Confidence:     1.0,
					Reason:         "Exact name match",
				})
				break
			} else if pClean == fClean {
				suggestions = append(suggestions, SuggestedMapping{
					ParameterName:  p.Name,
					SuggestedField: f,
					FunctionalType: p.FunctionalType,
					Confidence:     0.85,
					Reason:         "Normalized name match",
				})
				break
			}
		}
	}

	return &SuggestInputMappingResponse{
		Suggestions:         suggestions,
		DeterministicMethod: "normalized-name-matching",
		IsDeterministic:     true,
		CorrelationID:       corrID,
	}, nil
}

// -------------------------------------------------------------------------
// REQ-PROD-006: Permanent Execution Blocking (all return 400 execution-blocked)
// -------------------------------------------------------------------------

// Invoke is permanently blocked under REQ-PROD-006.
func (s *DefaultToolInputService) Invoke(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Invoke", GenerateCorrelationID())
}

// Execute is permanently blocked under REQ-PROD-006.
func (s *DefaultToolInputService) Execute(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Execute", GenerateCorrelationID())
}

// Run is permanently blocked under REQ-PROD-006.
func (s *DefaultToolInputService) Run(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Run", GenerateCorrelationID())
}

// Call is permanently blocked under REQ-PROD-006.
func (s *DefaultToolInputService) Call(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Call", GenerateCorrelationID())
}

// Trigger is permanently blocked under REQ-PROD-006.
func (s *DefaultToolInputService) Trigger(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Trigger", GenerateCorrelationID())
}
