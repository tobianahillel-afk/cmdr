package huntmanagement

import (
	"sort"
	"strings"
	"time"
)

func Project(in Input) Result {
	result := Result{AuditRequired: true}
	reject := func(code ErrorCode) Result {
		result.ErrorCode = code
		return result
	}

	if strings.TrimSpace(in.MutationRequested) != "" {
		return reject(ErrorMutationOpenDecision)
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
	if strings.TrimSpace(in.Question) == "" {
		return reject(ErrorMissingQuestion)
	}
	if strings.TrimSpace(in.Scope) == "" {
		return reject(ErrorMissingScope)
	}
	start, errStart := time.Parse(time.RFC3339, in.TimeStart)
	end, errEnd := time.Parse(time.RFC3339, in.TimeEnd)
	if errStart != nil || errEnd != nil || !start.Before(end) {
		return reject(ErrorInvalidTimeRange)
	}
	if strings.TrimSpace(in.OwnerRef) == "" {
		return reject(ErrorMissingOwner)
	}

	contributors := make([]string, len(in.ContributorRefs))
	for i, contributor := range in.ContributorRefs {
		if strings.TrimSpace(contributor) == "" {
			return reject(ErrorInvalidContributor)
		}
		contributors[i] = contributor
	}

	projected := make([]ProjectedReference, 0, len(in.References))
	limitations := make([]string, 0, len(in.References))
	visibleQueries := map[string]bool{}
	for _, ref := range in.References {
		owner, ok := canonicalOwner(ref.Kind)
		if !ok || strings.TrimSpace(ref.ID) == "" {
			return reject(ErrorInvalidReference)
		}
		if ref.TenantRef != in.TenantRef {
			return reject(ErrorCrossTenantReference)
		}
		if ref.AccessDecision != AccessAllow {
			return reject(ErrorReferenceAccessDenied)
		}
		switch ref.State {
		case ReferenceAvailable:
		case ReferenceStale, ReferencePartial, ReferenceUnavailable:
			limitations = append(limitations, string(ref.Kind)+":"+ref.ID+":"+string(ref.State))
		default:
			return reject(ErrorInvalidReferenceState)
		}
		if ref.Kind == ReferenceQuery {
			visibleQueries[ref.ID] = true
		}
		projected = append(projected, ProjectedReference{
			Kind: ref.Kind, ID: ref.ID, TenantRef: ref.TenantRef, QueryRef: ref.QueryRef,
			State: ref.State, CanonicalOwner: owner, Version: ref.Version, ProvenanceRef: ref.ProvenanceRef,
		})
	}
	for _, ref := range projected {
		if ref.Kind == ReferenceSearchJob &&
			(strings.TrimSpace(ref.QueryRef) == "" || !visibleQueries[ref.QueryRef]) {
			return reject(ErrorSearchJobQueryNotVisible)
		}
	}

	sort.Strings(limitations)
	result.Allowed = true
	result.Projection = &Projection{
		TenantRef: in.TenantRef, EnvironmentRef: in.EnvironmentRef,
		Question: in.Question, Scope: in.Scope,
		TimeRange: TimeRange{Start: in.TimeStart, End: in.TimeEnd},
		OwnerRef: in.OwnerRef, ContributorRefs: contributors,
		References: projected, Limitations: limitations,
		Complete: len(limitations) == 0, CorrelationID: in.CorrelationID,
	}
	return result
}

func canonicalOwner(kind ReferenceKind) (string, bool) {
	switch kind {
	case ReferenceQuery, ReferenceSearchJob:
		return "Shared Capabilities", true
	case ReferenceHypothesis, ReferenceCase:
		return "Investigate", true
	case ReferenceIncident:
		return "Command", true
	default:
		return "", false
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
