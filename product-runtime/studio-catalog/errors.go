// Package studiocatalog implements the genuine Go runtime for CMDR Studio Library & Asset Catalog
// conforming to CAP-STD-001, CAP-STD-002, tenant-scope-contract.md, error-contract.md, and permission-catalog.md.
package studiocatalog

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// ErrorCode represents deterministic machine-readable error codes.
type ErrorCode string

const (
	// Tenant scoping errors (tenant-scope-contract.md)
	ErrorCodeMissingTenant            ErrorCode = "missing-tenant"
	ErrorCodeWildcardTenant           ErrorCode = "wildcard-tenant"
	ErrorCodeUnauthorizedTenantInSet  ErrorCode = "unauthorized-tenant-in-set"
	ErrorCodeCrossTenantMutationForbidden ErrorCode = "cross-tenant-mutation-forbidden"

	// Resource & Oracle errors (CAP-STD-001, CAP-STD-002)
	ErrorCodeAssetNotFound            ErrorCode = "asset-not-found"
	ErrorCodeMissingAssetID           ErrorCode = "missing-asset-id"
	ErrorCodeInvalidAssetKind         ErrorCode = "invalid-asset-kind"
	ErrorCodeMissingReturnOrigin      ErrorCode = "missing-return-origin"

	// Permissions errors (permission-catalog.md)
	ErrorCodePermissionDenied         ErrorCode = "permission-denied"

	// Query & Pagination errors (pagination-filtering-sorting.md)
	ErrorCodeInvalidCursor            ErrorCode = "invalid-cursor"
	ErrorCodeInvalidFilterParameter   ErrorCode = "invalid-filter-parameter"
	ErrorCodeUnsupportedGroupingFacet ErrorCode = "unsupported-grouping-facet"
	ErrorCodePageLimitExceeded        ErrorCode = "page-limit-exceeded"

	// Execution policy errors (REQ-PROD-006)
	ErrorCodeExecutionBlocked         ErrorCode = "execution-blocked"

	// Generic internal error
	ErrorCodeInternalError            ErrorCode = "internal-error"
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

// GenerateCorrelationID generates a secure random correlation identifier.
func GenerateCorrelationID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "corr-fallback-deterministic"
	}
	return fmt.Sprintf("corr-%s", hex.EncodeToString(b))
}
