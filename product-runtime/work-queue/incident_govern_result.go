package workqueue

import (
	"context"
	"errors"
	"time"
)

// IncidentResultReference comes from a trusted, independently authorized
// object-link owner. A link never grants access to its Govern-owned Result.
type IncidentResultReference struct {
	IncidentID  string
	ResultID    string
	Tenant      string
	Environment string
}

// IncidentResultReferencesOwner authorizes the relation for a given Incident.
// It is not the Govern data owner and exposes only stable scoped references.
type IncidentResultReferencesOwner interface {
	ReadAuthorizedIncidentResultReferences(ctx context.Context, principal, tenant, environment, incidentID string) ([]IncidentResultReference, error)
}

// AuthorizedGovernResult is a bounded, already-permission-filtered Govern
// projection, not executor output or a Command Incident state transition.
type AuthorizedGovernResult struct {
	ID            string
	Tenant        string
	Environment   string
	RunID         string
	Lifecycle     string
	Outcome       string
	Verification  string
	ResidualRisk  string
	SourceRef     string
	SourceVersion string
	CreatedAt     string
}

// GovernResultOwner is a separate source-side permission boundary. It must
// evaluate perm.govern.result.read, tenant, environment and each Result ID
// independently; returning nil masks inaccessible results.
type GovernResultOwner interface {
	ReadAuthorizedGovernResult(ctx context.Context, principal, tenant, environment, resultID string) (*AuthorizedGovernResult, error)
}

var errIncidentGovernResult = errors.New("incident result projection unavailable")

func readAuthorizedIncidentGovernResults(
	ctx context.Context, links IncidentResultReferencesOwner, govern GovernResultOwner,
	scope SessionScope, incidentID string,
) ([]AuthorizedGovernResult, error) {
	if links == nil || govern == nil {
		return nil, nil // optional partial projection, never a claim of no Result
	}
	if ctx.Err() != nil {
		return nil, errIncidentGovernResult
	}
	references, err := links.ReadAuthorizedIncidentResultReferences(ctx,
		scope.Principal, scope.Tenant, scope.Environment, incidentID)
	if err != nil || ctx.Err() != nil || len(references) > 8 {
		return nil, errIncidentGovernResult
	}
	seen := make(map[string]bool, len(references))
	out := make([]AuthorizedGovernResult, 0, len(references))
	for _, ref := range references {
		if ctx.Err() != nil ||
			ref.IncidentID != incidentID || ref.Tenant != scope.Tenant ||
			ref.Environment != scope.Environment || invalidScope(ref.ResultID) ||
			!validOverviewText(ref.ResultID, 128) || seen[ref.ResultID] {
			return nil, errIncidentGovernResult
		}
		seen[ref.ResultID] = true
		result, err := govern.ReadAuthorizedGovernResult(ctx,
			scope.Principal, scope.Tenant, scope.Environment, ref.ResultID)
		if err != nil || ctx.Err() != nil {
			return nil, errIncidentGovernResult
		}
		if result == nil {
			continue // denial: never disclose ID, status, existence or count
		}
		if result.ID != ref.ResultID || result.Tenant != scope.Tenant ||
			result.Environment != scope.Environment || invalidScope(result.RunID) ||
			!validOverviewText(result.RunID, 128) ||
			invalidScope(result.SourceRef) || invalidScope(result.SourceVersion) ||
			!validOverviewText(result.SourceRef, 128) ||
			!validOverviewText(result.SourceVersion, 80) ||
			!validOverviewText(result.Lifecycle, 64) ||
			!validOverviewText(result.Outcome, 64) ||
			!validOverviewText(result.Verification, 128) ||
			!validOverviewText(result.ResidualRisk, 256) {
			return nil, errIncidentGovernResult
		}
		created, err := time.Parse(time.RFC3339Nano, result.CreatedAt)
		if err != nil || created.IsZero() {
			return nil, errIncidentGovernResult
		}
		out = append(out, *result)
	}
	return out, nil
}
