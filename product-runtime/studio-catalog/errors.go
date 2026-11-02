// Package studiocatalog implements the genuine Go runtime for CMDR Studio Library & Asset Catalog
// conforming to CAP-STD-001, CAP-STD-002, tenant-scope-contract.md, error-contract.md, and permission-catalog.md.
package studiocatalog

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

// Sentinel error definitions conforming to CAP-STD-004.
var (
	ErrMissingInput              = errors.New("missing-input")
	ErrInvalidInput              = errors.New("invalid-input")
	ErrRestrictedInput           = errors.New("restricted-input")
	ErrCrossScopeBlocked         = errors.New("cross-scope-blocked")
	ErrStaleReference            = errors.New("stale-reference")
	ErrUnsupportedFunctionalType = errors.New("unsupported-functional-type")
	ErrExecutionBlocked          = errors.New("execution-blocked")
)

// ErrorCode represents deterministic machine-readable error codes.
type ErrorCode string

const (
	// Tenant scoping errors (tenant-scope-contract.md)
	ErrorCodeMissingTenant                ErrorCode = "missing-tenant"
	ErrorCodeWildcardTenant               ErrorCode = "wildcard-tenant"
	ErrorCodeUnauthorizedTenantInSet      ErrorCode = "unauthorized-tenant-in-set"
	ErrorCodeCrossTenantMutationForbidden ErrorCode = "cross-tenant-mutation-forbidden"

	// Resource & Oracle errors (CAP-STD-001, CAP-STD-002)
	ErrorCodeAssetNotFound       ErrorCode = "asset-not-found"
	ErrorCodeMissingAssetID      ErrorCode = "missing-asset-id"
	ErrorCodeInvalidAssetKind    ErrorCode = "invalid-asset-kind"
	ErrorCodeMissingReturnOrigin ErrorCode = "missing-return-origin"

	// Permissions errors (permission-catalog.md)
	ErrorCodePermissionDenied ErrorCode = "permission-denied"

	// Query & Pagination errors (pagination-filtering-sorting.md)
	ErrorCodeInvalidCursor            ErrorCode = "invalid-cursor"
	ErrorCodeInvalidFilterParameter   ErrorCode = "invalid-filter-parameter"
	ErrorCodeUnsupportedGroupingFacet ErrorCode = "unsupported-grouping-facet"
	ErrorCodePageLimitExceeded        ErrorCode = "page-limit-exceeded"

	// Execution policy errors (REQ-PROD-006)
	ErrorCodeExecutionBlocked ErrorCode = "execution-blocked"

	// Tool input validation errors (CAP-STD-004)
	ErrorCodeMissingInput              ErrorCode = "missing-input"
	ErrorCodeInvalidInput              ErrorCode = "invalid-input"
	ErrorCodeRestrictedInput           ErrorCode = "restricted-input"
	ErrorCodeCrossScopeBlocked         ErrorCode = "cross-scope-blocked"
	ErrorCodeStaleReference            ErrorCode = "stale-reference"
	ErrorCodeUnsupportedFunctionalType ErrorCode = "unsupported-functional-type"

	// Metadata & Supersession errors (CAP-STD-002)
	ErrorCodeImmutableAssetMutation ErrorCode = "immutable-asset-mutation"
	ErrorCodeInvalidMetadata        ErrorCode = "invalid-metadata"
	ErrorCodeCircularSupersession   ErrorCode = "circular-supersession"

	// Generic internal error
	ErrorCodeInternalError ErrorCode = "internal-error"
)

// ErrorEnvelope adheres strictly to cmdr-product-spec/17-implementation-contracts/error-contract.md.
type ErrorEnvelope struct {
	Code          ErrorCode `json:"code"`
	Message       string    `json:"message"`
	CorrelationID string    `json:"correlation_id"`
	Retryable     bool      `json:"retryable"`
	PartialResult any       `json:"partial_result"`
	HTTPStatus    int       `json:"-"`
}

// Error implements the standard Go error interface.
func (e *ErrorEnvelope) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s (correlation_id: %s)", e.Code, e.Message, e.CorrelationID)
}

// ErrorResponse wraps the error envelope for standard API JSON responses.
type ErrorResponse struct {
	Error *ErrorEnvelope `json:"error"`
}

// NewError creates a new canonical ErrorEnvelope.
func NewError(code ErrorCode, message string, correlationID string, httpStatus int, retryable bool) *ErrorEnvelope {
	if correlationID == "" {
		correlationID = GenerateCorrelationID()
	}
	return &ErrorEnvelope{
		Code:          code,
		Message:       message,
		CorrelationID: correlationID,
		Retryable:     retryable,
		PartialResult: nil,
		HTTPStatus:    httpStatus,
	}
}

// NewAssetNotFoundError creates a uniform 404 envelope eliminating existence oracles.
func NewAssetNotFoundError(assetID, tenantID, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeAssetNotFound,
		fmt.Sprintf("Asset %q was not found in tenant %q", assetID, tenantID),
		correlationID,
		404,
		false,
	)
}

// NewPermissionDeniedError creates a 403 envelope when general library read is missing.
func NewPermissionDeniedError(requiredPermission, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodePermissionDenied,
		fmt.Sprintf("Caller lacks required permission %q to access Studio Library", requiredPermission),
		correlationID,
		403,
		false,
	)
}

// NewExecutionBlockedError creates a 400 envelope permanently blocking execution verbs under REQ-PROD-006.
func NewExecutionBlockedError(verb, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeExecutionBlocked,
		fmt.Sprintf("Execution verb %q is permanently blocked in Studio Library (Class 0 read-only discovery, REQ-PROD-006)", verb),
		correlationID,
		400,
		false,
	)
}

// NewImmutableAssetMutationError creates a 400 envelope when attempting to mutate non-draft metadata.
func NewImmutableAssetMutationError(assetID, state, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeImmutableAssetMutation,
		fmt.Sprintf("Asset %q in state %q cannot be edited (only draft assets are mutable under CAP-STD-002)", assetID, state),
		correlationID,
		400,
		false,
	)
}

// NewInvalidMetadataError creates a 400 envelope when metadata fails required validation.
func NewInvalidMetadataError(details, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeInvalidMetadata,
		fmt.Sprintf("Invalid metadata specification: %s", details),
		correlationID,
		400,
		false,
	)
}

// NewCircularSupersessionError creates a 400 envelope when supersession graph contains cycles.
func NewCircularSupersessionError(assetID, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeCircularSupersession,
		fmt.Sprintf("Circular supersession link detected involving asset %q", assetID),
		correlationID,
		400,
		false,
	)
}

// NewMissingInputError creates a 400 envelope when a required parameter is missing without default.
func NewMissingInputError(paramName, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeMissingInput,
		fmt.Sprintf("Required parameter %q is missing and has no explicit default (CAP-STD-004 §25)", paramName),
		correlationID,
		400,
		false,
	)
}

// NewInvalidInputError creates a 400 envelope when an input violates type or constraints.
func NewInvalidInputError(paramName, reason, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeInvalidInput,
		fmt.Sprintf("Parameter %q is invalid: %s", paramName, reason),
		correlationID,
		400,
		false,
	)
}

// NewRestrictedInputError creates a 403 envelope when an unauthorized or stale secret ref is accessed.
func NewRestrictedInputError(paramName, reason, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeRestrictedInput,
		fmt.Sprintf("Parameter %q access restricted: %s", paramName, reason),
		correlationID,
		403,
		false,
	)
}

// NewCrossScopeBlockedError creates a 403 envelope when an input ref violates tenant/scope isolation.
func NewCrossScopeBlockedError(paramName, targetScope, callerScope, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeCrossScopeBlocked,
		fmt.Sprintf("Parameter %q target scope %q violates caller scope %q (cross-scope-blocked under CAP-STD-004 §22)", paramName, targetScope, callerScope),
		correlationID,
		403,
		false,
	)
}

// NewStaleReferenceError creates a 400 envelope when a referenced entity or secret reference is expired.
func NewStaleReferenceError(refID, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeStaleReference,
		fmt.Sprintf("Referenced entity or secret %q is stale or expired (CAP-STD-004 §22)", refID),
		correlationID,
		400,
		false,
	)
}

// NewUnsupportedFunctionalTypeError creates a 400 envelope when a functional type is unrecognized.
func NewUnsupportedFunctionalTypeError(fType, correlationID string) *ErrorEnvelope {
	return NewError(
		ErrorCodeUnsupportedFunctionalType,
		fmt.Sprintf("Functional type %q is unsupported (CAP-STD-004 §8, §11)", fType),
		correlationID,
		400,
		false,
	)
}

// Is supports standard errors.Is comparison against sentinel errors.
func (e *ErrorEnvelope) Is(target error) bool {
	if e == nil || target == nil {
		return false
	}
	if target.Error() == string(e.Code) {
		return true
	}
	switch e.Code {
	case ErrorCodeMissingInput:
		return target == ErrMissingInput
	case ErrorCodeInvalidInput:
		return target == ErrInvalidInput
	case ErrorCodeRestrictedInput:
		return target == ErrRestrictedInput
	case ErrorCodeCrossScopeBlocked:
		return target == ErrCrossScopeBlocked
	case ErrorCodeStaleReference:
		return target == ErrStaleReference
	case ErrorCodeUnsupportedFunctionalType:
		return target == ErrUnsupportedFunctionalType
	case ErrorCodeExecutionBlocked:
		return target == ErrExecutionBlocked
	}
	return false
}

// GenerateCorrelationID generates a secure random correlation identifier.
func GenerateCorrelationID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "corr-fallback-deterministic"
	}
	return fmt.Sprintf("corr-%s", hex.EncodeToString(b))
}
