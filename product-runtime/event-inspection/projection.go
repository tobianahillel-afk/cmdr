package eventinspection

import (
	"strings"
	"unicode/utf8"
)

func Project(in Input) Result {
	result := Result{AuditRequired: true}
	reject := func(code ErrorCode) Result {
		result.ErrorCode = code
		return result
	}

	if strings.TrimSpace(in.MutationRequested) != "" {
		switch in.MutationRequested {
		case "case-link", "artifact-proposal", "evidence-candidate":
			return reject(ErrorMutationExcluded)
		default:
			return reject(ErrorUnsupportedMutation)
		}
	}
	if strings.TrimSpace(in.TenantRef) == "" {
		return reject(ErrorMissingTenant)
	}
	if in.TenantRef == "*" {
		return reject(ErrorWildcardTenant)
	}
	if !contains(in.AuthorizedTenants, in.TenantRef) {
		return reject(ErrorCrossTenant)
	}
	if strings.TrimSpace(in.EnvironmentRef) == "" {
		return reject(ErrorMissingEnvironment)
	}
	if strings.TrimSpace(in.Event.Ref) == "" {
		return reject(ErrorMissingEventRef)
	}
	if strings.TrimSpace(in.Event.TenantRef) == "" || in.Event.TenantRef != in.TenantRef {
		return reject(ErrorEventTenantMismatch)
	}
	if !contains(in.Permissions, PermissionTelemetryRead) {
		return reject(ErrorPermissionDenied)
	}
	if !validAccessDecision(in.RawAccessDecision) || !validAccessDecision(in.RenderedAccessDecision) {
		return reject(ErrorInvalidAccessDecision)
	}
	if in.RawAccessDecision == AccessDeny && in.RenderedAccessDecision == AccessDeny {
		return reject(ErrorProjectionDenied)
	}
	if !validOpaqueReference(in.CorrelationID) {
		return reject(ErrorInvalidCorrelationID)
	}
	if !validEventState(in.Event.State) {
		return reject(ErrorInvalidEventState)
	}

	enrichments := make([]ProjectionEnrichment, 0, len(in.Event.Enrichments))
	for _, enrichment := range in.Event.Enrichments {
		if strings.TrimSpace(enrichment.Producer) == "" ||
			strings.TrimSpace(enrichment.Version) == "" ||
			strings.TrimSpace(enrichment.Freshness) == "" {
			return reject(ErrorInvalidEnrichmentProvenance)
		}
		if enrichment.PresentedAsSource {
			return reject(ErrorDerivedAsSource)
		}
		enrichments = append(enrichments, ProjectionEnrichment{
			Producer: enrichment.Producer, Version: enrichment.Version, Freshness: enrichment.Freshness,
			Stale: strings.EqualFold(enrichment.Freshness, "stale"), Values: cloneStringMap(enrichment.Values),
		})
	}

	projection := &Projection{
		TenantRef: in.TenantRef, EnvironmentRef: in.EnvironmentRef, EventRef: in.Event.Ref,
		CorrelationID: in.CorrelationID, State: in.Event.State, Source: in.Event.Source, Parser: in.Event.Parser,
		SourceFields: cloneStringMap(in.Event.SourceFields), NormalizedFields: cloneStringMap(in.Event.NormalizedFields),
		MissingFields: append([]string(nil), in.Event.MissingFields...), Enrichments: enrichments,
	}
	if in.RawAccessDecision == AccessAllow {
		projection.RawPayload = append([]byte(nil), in.Event.RawPayload...)
		projection.RawDerivedSensitive = cloneStringMap(in.Event.RawDerivedSensitive)
	}
	if in.RenderedAccessDecision == AccessAllow {
		projection.RenderedPayload = append([]byte(nil), in.Event.RenderedPayload...)
	}
	result.Allowed = true
	result.Projection = projection
	return result
}

func validAccessDecision(value AccessDecision) bool {
	return value == AccessAllow || value == AccessDeny
}

func validEventState(value EventState) bool {
	switch value {
	case StateAvailable, StateSourceUnavailable, StateTombstone, StatePartial, StateRestricted:
		return true
	default:
		return false
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	out := make(map[string]string, len(values))
	for key, value := range values {
		out[key] = value
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
