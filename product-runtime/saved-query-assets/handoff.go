package savedqueryassets

import "strings"

// BuildExecutionHandoff creates a backend-neutral immutable draft for Event Search.
// It never creates or executes a Search Job and never calls another runtime.
func BuildExecutionHandoff(projected Result, in HandoffInput) HandoffResult {
	result := HandoffResult{AuditRequired: true}
	reject := func(code ErrorCode) HandoffResult {
		result.ErrorCode = code
		return result
	}

	if !projected.Allowed || projected.Projection == nil {
		return reject(ErrorHandoffProjectionUnavailable)
	}
	projection := projected.Projection
	if !projected.HandoffEligible || !projection.HandoffEligible ||
		projected.Compatibility != "compatible" || projection.Compatibility != "compatible" {
		return reject(ErrorHandoffIneligible)
	}
	if strings.TrimSpace(in.TenantRef) == "" || strings.TrimSpace(in.EnvironmentRef) == "" ||
		in.TenantRef != projection.TenantRef || in.EnvironmentRef != projection.EnvironmentRef ||
		projection.QueryTenantRef != projection.TenantRef ||
		strings.TrimSpace(projection.QueryRef) == "" || strings.TrimSpace(projection.QueryVersion) == "" {
		return reject(ErrorHandoffScopeMismatch)
	}

	declaredParameters := make(map[string]struct{}, len(projection.ParameterNames))
	parameterNames := make([]string, len(projection.ParameterNames))
	for i, name := range projection.ParameterNames {
		if strings.TrimSpace(name) == "" {
			return reject(ErrorHandoffIneligible)
		}
		if _, exists := declaredParameters[name]; exists {
			return reject(ErrorHandoffIneligible)
		}
		declaredParameters[name] = struct{}{}
		parameterNames[i] = name
	}

	seenValues := make(map[string]struct{}, len(in.ParameterValues))
	parameterValues := make([]ParameterValue, 0, len(in.ParameterValues))
	for _, binding := range in.ParameterValues {
		if _, ok := declaredParameters[binding.Name]; !ok {
			return reject(ErrorHandoffUnknownParameter)
		}
		if _, duplicate := seenValues[binding.Name]; duplicate {
			return reject(ErrorHandoffDuplicateParameter)
		}
		seenValues[binding.Name] = struct{}{}
		parameterValues = append(parameterValues, ParameterValue{
			Name:        binding.Name,
			OpaqueValue: append([]byte(nil), binding.OpaqueValue...),
		})
	}

	sources := make([]HandoffSourceContext, 0, len(projection.Sources))
	for _, source := range projection.Sources {
		if strings.TrimSpace(source.ID) == "" || source.TenantRef != projection.TenantRef {
			return reject(ErrorHandoffIneligible)
		}
		requiredFields := make([]string, len(source.Fields))
		for i, field := range source.Fields {
			if strings.TrimSpace(field.Name) == "" {
				return reject(ErrorHandoffIneligible)
			}
			requiredFields[i] = field.Name
		}
		sources = append(sources, HandoffSourceContext{
			ID:             source.ID,
			RequiredFields: requiredFields,
		})
	}

	result.Allowed = true
	result.Draft = &ExecutionHandoffDraft{
		TenantRef:                      projection.TenantRef,
		EnvironmentRef:                 projection.EnvironmentRef,
		QueryRef:                       projection.QueryRef,
		QueryTenantRef:                 projection.QueryTenantRef,
		QueryVersion:                   projection.QueryVersion,
		ParameterNames:                 parameterNames,
		ParameterValues:                parameterValues,
		Sources:                        sources,
		ReturnContext:                  ReturnContextToken{AssetID: projection.AssetID, AssetKind: projection.AssetKind},
		SourceReevaluationRequired:     true,
		PermissionReevaluationRequired: true,
		OriginalAssetImmutable:         true,
	}
	return result
}
