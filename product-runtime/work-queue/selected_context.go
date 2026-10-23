package workqueue

import (
	"context"
	"errors"
)

// enrichSelectedResult reuses exactly one host-owner authorization chain for
// the existing Work Queue selection and the read-only Incident Detail page.
// Object read never grants SLA, canonical register or linked Case read.
func enrichSelectedResult(ctx context.Context, owner Owner, scope SessionScope, result *Result) error {
	if result.Selected == nil {
		return nil
	}
	row := result.Selected
	factors := AuthorizedFactors{}
	if source, ok := owner.(FactorsOwner); ok {
		fetched, err := source.ReadAuthorizedFactors(ctx, scope.Principal,
			scope.Tenant, scope.Environment, row.Kind, row.ID)
		if err != nil {
			return errors.New("coordination factors unavailable")
		}
		factors, err = validateAuthorizedFactors(fetched)
		if err != nil {
			return errors.New("coordination factors unavailable")
		}
	}
	result.SelectedFactors = &factors

	if source, ok := owner.(SLAOwner); ok {
		fetched, err := source.ReadAuthorizedSLA(ctx, scope.Principal,
			scope.Tenant, scope.Environment, row.Kind, row.ID)
		if err != nil {
			return errors.New("SLA source unavailable")
		}
		result.SelectedSLA, err = validateAuthorizedSLA(fetched)
		if err != nil {
			return errors.New("SLA source unavailable")
		}
	}
	if source, ok := owner.(CanonicalOwnershipOwner); ok {
		fetched, err := source.ReadAuthorizedCanonicalOwnership(ctx,
			scope.Principal, scope.Tenant, scope.Environment, row.Kind, row.ID)
		if err != nil {
			return errors.New("canonical ownership unavailable")
		}
		result.SelectedCanonicalOwnership, err = validateAuthorizedCanonicalOwnership(fetched)
		if err != nil {
			return errors.New("canonical ownership unavailable")
		}
	}
	if row.Kind == Incident {
		if source, ok := owner.(LinkedCasesOwner); ok {
			fetched, err := source.ReadAuthorizedLinkedCases(ctx, scope.Principal,
				scope.Tenant, scope.Environment, row.ID)
			if err != nil {
				return errors.New("linked Cases unavailable")
			}
			result.SelectedLinkedCases, err = validateLinkedCases(fetched)
			if err != nil {
				return errors.New("linked Cases unavailable")
			}
		}
	}
	return nil
}
