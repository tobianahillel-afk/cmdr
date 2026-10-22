package savedqueryassets

import (
	"sort"
	"strings"
)

func Project(in Input) Result {
	result := Result{AuditRequired: true, Diagnostics: []string{}}
	reject := func(code ErrorCode) Result {
		result.ErrorCode = code
		result.Diagnostics = []string{}
		return result
	}

	if strings.TrimSpace(in.MutationRequested) != "" {
		return reject(ErrorMutationOpenDecision)
	}
	if in.SearchJobExecutionRequested {
		return reject(ErrorSearchJobExecutionForbidden)
	}
	if strings.TrimSpace(in.TenantRef) == "" {
		return reject(ErrorMissingTenant)
	}
	if in.TenantRef == "*" {
		return reject(ErrorWildcardTenant)
	}
	if !containsString(in.AuthorizedTenants, in.TenantRef) {
		return reject(ErrorCrossTenant)
	}
	if strings.TrimSpace(in.EnvironmentRef) == "" {
		return reject(ErrorMissingEnvironment)
	}
	if strings.TrimSpace(in.AssetID) == "" {
		return reject(ErrorMissingAssetID)
	}
	if in.AssetKind != AssetSavedSearch && in.AssetKind != AssetQueryAsset {
		return reject(ErrorInvalidAssetKind)
	}
	if in.AssetAccessDecision != AccessAllow {
		return reject(ErrorAssetAccessDenied)
	}
	if strings.TrimSpace(in.QueryRef) == "" {
		return reject(ErrorMissingQueryRef)
	}
	if in.QueryTenantRef != in.TenantRef {
		return reject(ErrorCrossTenantQuery)
	}
	if strings.TrimSpace(in.QueryVersion) == "" {
		return reject(ErrorMissingQueryVersion)
	}
	if in.QueryAccessDecision != AccessAllow {
		return reject(ErrorQueryAccessDenied)
	}
	if strings.TrimSpace(in.AuthorRef) == "" {
		return reject(ErrorMissingAuthor)
	}

	incompatible := false
	stale := false
	switch in.ValidationState {
	case ValidationValidated:
	case ValidationStale:
		stale = true
		result.Diagnostics = append(result.Diagnostics, "validation:stale")
	case ValidationIncompatible:
		incompatible = true
		result.Diagnostics = append(result.Diagnostics, "validation:incompatible")
	default:
		return reject(ErrorInvalidValidationState)
	}

	parameters := make([]string, len(in.ParameterNames))
	for i, parameter := range in.ParameterNames {
		if strings.TrimSpace(parameter) == "" {
			return reject(ErrorInvalidParameter)
		}
		parameters[i] = parameter
	}
	if len(in.Sources) == 0 {
		return reject(ErrorMissingSourcePrerequisites)
	}

	sources := make([]ProjectedSource, 0, len(in.Sources))
	for _, source := range in.Sources {
		if strings.TrimSpace(source.ID) == "" {
			return reject(ErrorInvalidSource)
		}
		if source.TenantRef != in.TenantRef {
			return reject(ErrorCrossTenantSource)
		}
		projected := ProjectedSource{
			ID: source.ID, TenantRef: source.TenantRef, State: source.State,
			Fields: make([]ProjectedField, 0, len(source.Fields)),
		}
		switch source.State {
		case SourceAvailable:
		case SourceStale:
			stale = true
			result.Diagnostics = append(result.Diagnostics, "source:"+source.ID+":stale")
		case SourceMissing:
			incompatible = true
			result.Diagnostics = append(result.Diagnostics, "source:"+source.ID+":missing")
			sources = append(sources, projected)
			continue
		default:
			return reject(ErrorInvalidSourceState)
		}
		for _, field := range source.Fields {
			if strings.TrimSpace(field.Name) == "" {
				return reject(ErrorInvalidField)
			}
			switch field.State {
			case FieldAvailable:
			case FieldStale:
				stale = true
				result.Diagnostics = append(result.Diagnostics, "field:"+source.ID+"/"+field.Name+":stale")
			case FieldMissing, FieldRemoved:
				incompatible = true
				result.Diagnostics = append(result.Diagnostics, "field:"+source.ID+"/"+field.Name+":"+string(field.State))
			default:
				return reject(ErrorInvalidFieldState)
			}
			projected.Fields = append(projected.Fields, ProjectedField{Name: field.Name, State: field.State})
		}
		sources = append(sources, projected)
	}

	sort.Strings(result.Diagnostics)
	result.Compatibility = "compatible"
	if incompatible {
		result.Compatibility = "incompatible"
	} else if stale {
		result.Compatibility = "stale"
	}
	result.Allowed = true
	result.HandoffEligible = result.Compatibility == "compatible"
	result.ImmutableProjection = true
	result.DeepCopyIsolated = true
	result.Projection = &Projection{
		TenantRef: in.TenantRef, EnvironmentRef: in.EnvironmentRef,
		AssetID: in.AssetID, AssetKind: in.AssetKind,
		QueryRef: in.QueryRef, QueryTenantRef: in.QueryTenantRef, QueryVersion: in.QueryVersion,
		AuthorRef: in.AuthorRef, ValidationState: in.ValidationState,
		ParameterNames: parameters, Sources: sources,
		LineageRef: in.LineageRef, DeprecationReason: in.DeprecationReason,
		ReplacementAssetRef: in.ReplacementAssetRef,
		Compatibility: result.Compatibility,
		Diagnostics: append([]string(nil), result.Diagnostics...),
		HandoffEligible: result.HandoffEligible,
	}
	return result
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
