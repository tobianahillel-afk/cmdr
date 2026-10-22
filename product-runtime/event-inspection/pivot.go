package eventinspection

import (
	"strings"
	"time"
	"unicode/utf8"
)

func PreparePivot(projection *Projection, selection PivotSelection) PivotResult {
	result := PivotResult{AuditRequired: true}
	reject := func(code ErrorCode) PivotResult {
		result.ErrorCode = code
		return result
	}

	if projection == nil || !validPivotProjectionIdentity(projection) {
		return reject(ErrorInvalidPivotContext)
	}
	if !selection.ValueAuthorized {
		return reject(ErrorRestrictedPivotValue)
	}
	if !validOpaqueReference(selection.Field) || !selection.ValuePresent || !validPivotValue(selection.Value) {
		return reject(ErrorInvalidPivotContext)
	}
	if !validOpaqueReference(selection.ReturnOrigin) {
		return reject(ErrorInvalidPivotContext)
	}
	start, errStart := time.Parse(time.RFC3339, selection.TimeStart)
	end, errEnd := time.Parse(time.RFC3339, selection.TimeEnd)
	if errStart != nil || errEnd != nil || !start.Before(end) {
		return reject(ErrorInvalidPivotContext)
	}
	fieldOrigin, ok := projectedFieldOrigin(projection, selection.Field, selection.Value)
	if !ok {
		return reject(ErrorInvalidPivotContext)
	}

	startText := start.UTC().Format(time.RFC3339Nano)
	endText := end.UTC().Format(time.RFC3339Nano)
	draft := &PivotDraft{
		TargetCapability: PivotTargetCapability,
		TenantRef: projection.TenantRef, EnvironmentRef: projection.EnvironmentRef,
		EventRef: projection.EventRef, CorrelationID: projection.CorrelationID,
		Field: selection.Field, Value: selection.Value, FieldOrigin: fieldOrigin,
		SourceRef: projection.Source.SourceRef, DataSourceRef: projection.Source.DataSourceRef,
		TimeRange: PivotTimeRange{Start: startText, End: endText},
		ReturnContext: PivotReturnContext{
			TenantRef: projection.TenantRef, EnvironmentRef: projection.EnvironmentRef,
			EventRef: projection.EventRef, CorrelationID: projection.CorrelationID,
			Origin: selection.ReturnOrigin,
		},
	}
	result.Allowed = true
	result.Draft = draft
	return result
}

func validPivotProjectionIdentity(projection *Projection) bool {
	return projection != nil &&
		validOpaqueReference(projection.TenantRef) &&
		projection.TenantRef != "*" &&
		validOpaqueReference(projection.EnvironmentRef) &&
		validOpaqueReference(projection.EventRef) &&
		validOpaqueReference(projection.CorrelationID) &&
		validOpaqueReference(projection.Source.SourceRef)
}

func validPivotValue(value string) bool {
	return len(value) <= 64*1024 && utf8.ValidString(value) &&
		!strings.ContainsRune(value, '\x00')
}

func projectedFieldOrigin(projection *Projection, field, value string) (string, bool) {
	if projection.SourceFields[field] == value {
		return "source", true
	}
	if projection.NormalizedFields[field] == value {
		return "normalized", true
	}
	if projection.RawDerivedSensitive[field] == value {
		return "raw-derived-sensitive", true
	}
	for _, enrichment := range projection.Enrichments {
		if enrichment.Values[field] == value {
			return "enrichment:" + enrichment.Producer + "@" + enrichment.Version, true
		}
	}
	return "", false
}
