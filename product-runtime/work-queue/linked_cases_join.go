package workqueue

import (
	"context"
)

// AuthorizedIncidentCaseReferences is the Shared Object Linking owner's
// read-authorized existing relations for ONE current Incident. A nil pointer
// means unknown/unavailable; a nonnil empty slice means authorized no links.
type AuthorizedIncidentCaseReferences struct {
	References []LinkedCaseReference
}

type LinkedCaseReference struct {
	CaseID string
	LinkSource string
	LinkVersion string
}

// IncidentCaseReferencesOwner is owned by the trusted HOST/Shared linking
// service. It MUST independently authorize each returned relation in scope;
// the Command Incident read decision NEVER authorizes the link.
type IncidentCaseReferencesOwner interface {
	ReadAuthorizedIncidentCaseReferences(ctx context.Context, principal, tenant, environment, incidentID string) (*AuthorizedIncidentCaseReferences, error)
}

// AuthorizedCaseIdentity is a trusted Investigate-owned projection, not a
// Command-owned Case copy. The Case owner must recheck read access per ID.
type AuthorizedCaseIdentity struct {
	ID string
	Tenant string
	Environment string
	Status string
}

type InvestigateCaseIdentityOwner interface {
	ReadAuthorizedCaseIdentity(ctx context.Context, principal, tenant, environment, caseID string) (*AuthorizedCaseIdentity, error)
}

// LinkedCaseOwnerJoin supplies the existing Command selected-Incident
// interface with TWO independent source owners. No new route, Case mutation,
// implicit permission or replicated Case status authority is introduced.
//
// Both host interfaces must implement their own permissions and audit.
// A nil Case means not readable/absent, without distinguishing which.
type LinkedCaseOwnerJoin struct {
	Links IncidentCaseReferencesOwner
	Cases InvestigateCaseIdentityOwner
}

func (j LinkedCaseOwnerJoin) ReadAuthorizedLinkedCases(ctx context.Context, principal, tenant, environment, incidentID string) (*AuthorizedLinkedCases, error) {
	if ctx == nil || ctx.Err() != nil || invalidScope(principal) ||
		invalidScope(tenant) || invalidScope(environment) || invalidScope(incidentID) {
		return nil,errLinkedCases
	}
	if j.Links == nil || j.Cases == nil {
		return nil,nil // no trustworthy complete source, NOT an empty Case set
	}
	refs,err := j.Links.ReadAuthorizedIncidentCaseReferences(ctx,principal,tenant,environment,incidentID)
	if err != nil { return nil,errLinkedCases }
	if refs == nil { return nil,nil }
	if len(refs.References)>16 { return nil,errLinkedCases }

	// Validate the complete link cohort BEFORE sending any potentially
	// malformed or repeated Case identifiers to another security authority.
	claimed := &AuthorizedLinkedCases{Cases:make([]LinkedCaseSummary,0,len(refs.References))}
	for _,link:=range refs.References {
		claimed.Cases=append(claimed.Cases,LinkedCaseSummary{
			CaseID:link.CaseID,Status:"unknown",LinkSource:link.LinkSource,LinkVersion:link.LinkVersion,
		})
	}
	if _,err:=validateLinkedCases(claimed);err!=nil { return nil,errLinkedCases }

	authorized:=&AuthorizedLinkedCases{Cases:make([]LinkedCaseSummary,0,len(claimed.Cases))}
	for _,link:=range claimed.Cases {
		if ctx.Err()!=nil { return nil,errLinkedCases }
		c,err:=j.Cases.ReadAuthorizedCaseIdentity(ctx,principal,tenant,environment,link.CaseID)
		if err!=nil { return nil,errLinkedCases }
		if c==nil { continue } // denied and absent are indistinguishable
		if c.ID!=link.CaseID || c.Tenant!=tenant || c.Environment!=environment {
			return nil,errLinkedCases
		}
		link.Status=c.Status // NEVER infer Investigate status from the link
		authorized.Cases=append(authorized.Cases,link)
	}
	// Do not reveal partial authorized results if even one returned source
	// record is inconsistent, malformed or outside its original scope.
	if _,err:=validateLinkedCases(authorized);err!=nil {
		return nil,errLinkedCases
	}
	return authorized,nil
}
