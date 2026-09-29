// Package contextenvelope implements the bounded PILOT-CONTEXT-ENVELOPE-V1
// projection semantics. It is intentionally pure and in-process: no persistence,
// network, authorization engine, or secret-bearing payload belongs here.
package contextenvelope

import (
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxReferenceBytes    = 128
	maxReturnOriginBytes = 1024
)

// AuthorizationOutcome is a fresh caller-supplied destination authorization result.
type AuthorizationOutcome string

const (
	AuthorizationAllow AuthorizationOutcome = "allow"
	AuthorizationDeny  AuthorizationOutcome = "deny"
)

// State is the functional projection state defined by the pilot contract.
type State string

const (
	StateContextValid        State = "context-valid"
	StateContextIncompatible State = "context-incompatible"
	StatePermissionDenied    State = "permission-denied"
	StateExpired             State = "expired"
	StateSourceMissing       State = "source-missing"
)

// ErrorCode is a contract-defined input-validation error.
type ErrorCode string

const (
	ErrorInvalidTenantRef  ErrorCode = "invalid-tenant-ref"
	ErrorUnsafeReturnOrigin ErrorCode = "unsafe-return-origin"
)

// Input contains only the references and caller-supplied facts needed to project
// navigation context. It deliberately has no payload, secret, token, or policy field.
type Input struct {
	SourceProduct                  string               `json:"source_product"`
	DestinationProduct             string               `json:"destination_product"`
	SourceTenantRef                string               `json:"source_tenant_ref"`
	DestinationTenantRef           string               `json:"destination_tenant_ref"`
	EnvironmentRef                 string               `json:"environment_ref"`
	EnvironmentTenantRef           string               `json:"environment_tenant_ref"`
	DestinationRequiresEnvironment bool                 `json:"destination_requires_environment"`
	EnvironmentCompatible          bool                 `json:"environment_compatible"`
	DestinationAuthorization       AuthorizationOutcome `json:"destination_authorization"`
	ContextFresh                   bool                 `json:"context_fresh"`
	ReturnOrigin                   string               `json:"return_origin"`
}

// Result is the deterministic projected context envelope.
type Result struct {
	State                     State     `json:"state"`
	TenantRef                 string    `json:"tenant_ref"`
	EnvironmentRef            string    `json:"environment_ref"`
	ReturnOrigin              string    `json:"return_origin"`
	EnvironmentChoiceRequired bool      `json:"environment_choice_required"`
	Cleared                   bool      `json:"cleared"`
	ClearedReason             string    `json:"cleared_reason"`
	ProtectedRefsVisible      bool      `json:"protected_refs_visible"`
	TraceMinimalOnly          bool      `json:"trace_minimal_only"`
	Error                     ErrorCode `json:"error"`
}

// Project applies PILOT-CONTEXT-ENVELOPE-V1 without performing authorization,
// persistence, network I/O, or implicit tenant/environment lookup.
func Project(in Input) Result {
	if !validReturnOrigin(in.ReturnOrigin) {
		return invalid(ErrorUnsafeReturnOrigin)
	}

	if in.SourceTenantRef == "" {
		return masked(StateSourceMissing, "source-missing", in.ReturnOrigin)
	}
	if !validReference(in.SourceTenantRef) || !validReference(in.DestinationTenantRef) {
		return invalid(ErrorInvalidTenantRef)
	}

	if in.DestinationAuthorization != AuthorizationAllow {
		return masked(StatePermissionDenied, "permission-denied", in.ReturnOrigin)
	}
	if !in.ContextFresh {
		return masked(StateExpired, "expired", in.ReturnOrigin)
	}

	if in.SourceTenantRef != in.DestinationTenantRef {
		return visibleIncompatible(
			in.DestinationTenantRef,
			in.ReturnOrigin,
			in.DestinationRequiresEnvironment,
			true,
			"tenant-changed",
		)
	}

	if in.EnvironmentRef == "" {
		if in.DestinationRequiresEnvironment {
			return visibleIncompatible(
				in.DestinationTenantRef,
				in.ReturnOrigin,
				true,
				false,
				"environment-required",
			)
		}
		return visibleValid(in.DestinationTenantRef, "", in.ReturnOrigin)
	}

	if !validReference(in.EnvironmentRef) ||
		!validReference(in.EnvironmentTenantRef) ||
		in.EnvironmentTenantRef != in.DestinationTenantRef {
		return visibleIncompatible(
			in.DestinationTenantRef,
			in.ReturnOrigin,
			in.DestinationRequiresEnvironment,
			true,
			"environment-tenant-mismatch",
		)
	}

	if !in.EnvironmentCompatible {
		return visibleIncompatible(
			in.DestinationTenantRef,
			in.ReturnOrigin,
			in.DestinationRequiresEnvironment,
			true,
			"environment-incompatible",
		)
	}

	return visibleValid(in.DestinationTenantRef, in.EnvironmentRef, in.ReturnOrigin)
}

func invalid(code ErrorCode) Result {
	return Result{
		ProtectedRefsVisible: false,
		TraceMinimalOnly:     true,
		Error:                code,
	}
}

func masked(state State, reason, returnOrigin string) Result {
	return Result{
		State:                state,
		ReturnOrigin:         returnOrigin,
		Cleared:              true,
		ClearedReason:        reason,
		ProtectedRefsVisible: false,
		TraceMinimalOnly:     true,
	}
}

func visibleValid(tenantRef, environmentRef, returnOrigin string) Result {
	return Result{
		State:                StateContextValid,
		TenantRef:            tenantRef,
		EnvironmentRef:       environmentRef,
		ReturnOrigin:         returnOrigin,
		ProtectedRefsVisible: true,
	}
}

func visibleIncompatible(tenantRef, returnOrigin string, choiceRequired, cleared bool, reason string) Result {
	return Result{
		State:                     StateContextIncompatible,
		TenantRef:                 tenantRef,
		ReturnOrigin:              returnOrigin,
		EnvironmentChoiceRequired: choiceRequired,
		Cleared:                   cleared,
		ClearedReason:             reason,
		ProtectedRefsVisible:      true,
	}
}

func validReference(value string) bool {
	if len(value) < 1 || len(value) > maxReferenceBytes || !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if r == '*' || unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func validReturnOrigin(value string) bool {
	if value == "" {
		return true
	}
	if len(value) > maxReturnOriginBytes || !utf8.ValidString(value) ||
		!strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") ||
		strings.Contains(value, "\\") {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return false
	}
	return !parsed.IsAbs() &&
		parsed.Host == "" &&
		parsed.Scheme == "" &&
		parsed.Opaque == "" &&
		parsed.RawQuery == "" &&
		parsed.Fragment == "" &&
		!parsed.ForceQuery
}
