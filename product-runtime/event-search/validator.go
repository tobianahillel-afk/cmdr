// Package eventsearch validates the bounded Event Search execution envelope.
//
// It deliberately does not parse or execute a query language, select an index,
// choose a storage engine, or contact a provider. It only enforces the
// tenant/time/source/permission/provenance invariants required before a later
// execution layer may create a Search Job.
package eventsearch

import (
	"strings"
	"time"
	"unicode/utf8"
)

type ErrorCode string

const (
	ErrorNone                 ErrorCode = ""
	ErrorCaseLinkOpenDecision ErrorCode = "case-link-open-decision"
	ErrorMissingTenant        ErrorCode = "missing-tenant"
	ErrorWildcardTenant       ErrorCode = "wildcard-tenant"
	ErrorCrossTenant          ErrorCode = "cross-tenant"
	ErrorMissingEnvironment   ErrorCode = "missing-environment"
	ErrorInvalidTimeRange     ErrorCode = "invalid-time-range"
	ErrorEmptySources         ErrorCode = "empty-sources"
	ErrorUnauthorizedSource   ErrorCode = "unauthorized-source"
	ErrorInvalidQuery         ErrorCode = "invalid-query"
	ErrorPermissionDenied     ErrorCode = "permission-denied"
	ErrorInvalidCorrelationID ErrorCode = "invalid-correlation-id"
	ErrorInvalidQueryVersion  ErrorCode = "invalid-query-version"
)

const (
	PermissionSearchExecute = "perm.investigate.search.execute"
	PermissionQueryRead     = "perm.shared-capabilities.query.read"
	PermissionSearchManage  = "perm.shared-capabilities.search-job.manage"
	PermissionTelemetryRead = "perm.shared-capabilities.telemetry-event.read"
)

var requiredExecutionPermissions = [...]string{
	PermissionSearchExecute,
	PermissionQueryRead,
	PermissionSearchManage,
	PermissionTelemetryRead,
}

type Input struct {
	TenantRef         string
	AuthorizedTenants []string
	EnvironmentRef    string
	TimeStart         string
	TimeEnd           string
	Sources           []string
	AuthorizedSources []string
	Permissions       []string
	QueryPresent      bool
	QueryValid        bool
	CorrelationID     string
	QueryVersion      string
	CaseLinkRequested bool
}

type Envelope struct {
	TenantRef      string
	EnvironmentRef string
	TimeStart      time.Time
	TimeEnd        time.Time
	Sources        []string
	CorrelationID  string
	QueryVersion   string
}

type Result struct {
	Allowed              bool
	ErrorCode            ErrorCode
	ProtectedDataVisible bool
	AuditRequired        bool
	Envelope             Envelope
}

func Validate(in Input) Result {
	result := Result{AuditRequired: true}
	reject := func(code ErrorCode) Result {
		result.ErrorCode = code
		return result
	}

	if in.CaseLinkRequested {
		return reject(ErrorCaseLinkOpenDecision)
	}

	tenant := in.TenantRef
	if strings.TrimSpace(tenant) == "" {
		return reject(ErrorMissingTenant)
	}
	if tenant == "*" {
		return reject(ErrorWildcardTenant)
	}
	if !contains(in.AuthorizedTenants, tenant) {
		return reject(ErrorCrossTenant)
	}

	environment := in.EnvironmentRef
	if strings.TrimSpace(environment) == "" {
		return reject(ErrorMissingEnvironment)
	}

	start, errStart := time.Parse(time.RFC3339, in.TimeStart)
	end, errEnd := time.Parse(time.RFC3339, in.TimeEnd)
	if errStart != nil || errEnd != nil || !start.Before(end) {
		return reject(ErrorInvalidTimeRange)
	}

	if len(in.Sources) == 0 {
		return reject(ErrorEmptySources)
	}
	authorizedSources := stringSet(in.AuthorizedSources)
	for _, source := range in.Sources {
		if strings.TrimSpace(source) == "" || !authorizedSources[source] {
			return reject(ErrorUnauthorizedSource)
		}
	}

	if !in.QueryPresent || !in.QueryValid {
		return reject(ErrorInvalidQuery)
	}

	permissions := stringSet(in.Permissions)
	for _, permission := range requiredExecutionPermissions {
		if !permissions[permission] {
			return reject(ErrorPermissionDenied)
		}
	}

	if !validOpaqueReference(in.CorrelationID) {
		return reject(ErrorInvalidCorrelationID)
	}
	if !validOpaqueReference(in.QueryVersion) {
		return reject(ErrorInvalidQueryVersion)
	}

	result.Allowed = true
	result.ProtectedDataVisible = true
	result.Envelope = Envelope{
		TenantRef:      tenant,
		EnvironmentRef: environment,
		TimeStart:      start,
		TimeEnd:        end,
		Sources:        append([]string(nil), in.Sources...),
		CorrelationID:  in.CorrelationID,
		QueryVersion:   in.QueryVersion,
	}
	return result
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func stringSet(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}

func validOpaqueReference(value string) bool {
	if value == "" || len(value) > 256 || !utf8.ValidString(value) || strings.TrimSpace(value) != value {
		return false
	}
	for _, r := range value {
		if !allowedReferenceRune(r) {
			return false
		}
	}
	return true
}

func allowedReferenceRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= '0' && r <= '9':
		return true
	case r == '-', r == '_', r == '.', r == ':', r == '/':
		return true
	default:
		return false
	}
}
