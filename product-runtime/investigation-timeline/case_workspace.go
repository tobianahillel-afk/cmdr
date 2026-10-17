package investigationtimeline

import "context"

// CaseTimelineCatalog lists Timeline Entry IDs belonging to a selected Case.
// It must be a trusted owner-side service. The caller supplies only the Case
// selection, never canonical Timeline Entry snapshots or source identities.
// The catalog must bind the authenticated principal to tenant, environment
// and Case membership. It is not a permission oracle.
type CaseTimelineCatalog interface {
	ListCaseEntryIDs(context.Context, AccessRequest) ([]string, error)
}

// OpenCaseTimeline is an in-process, read-only Case -> Timeline integration.
// It verifies Case access BEFORE consulting the trusted catalog, then delegates
// to the existing owner-snapshot/typed-source pipeline. The second Case check
// there is intentional: access can be revoked between catalog and projection.
// No browser flags, untrusted IDs, persistence, or HTTP authentication are
// treated as authority by this function.
func OpenCaseTimeline(
	ctx context.Context, scope Scope, caseID string,
	checker AccessChecker, catalog CaseTimelineCatalog,
	entries EntryResolver, sources SourceOwnerResolver,
) Result {
	deny := Result{Status: "deny"}
	if missingAuthority(ctx) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) ||
		missingAuthority(sources) || invalidScope(scope.Tenant) ||
		invalidScope(scope.Environment) || invalidScope(caseID) ||
		len(scope.Tenant)+len(scope.Environment)+len(caseID) > MaxInputBytes ||
		ctx.Err() != nil {
		return deny
	}
	caseRequest := AccessRequest{
		Tenant: scope.Tenant, Environment: scope.Environment,
		CaseID: caseID, Resource: "case", ObjectID: caseID,
		Permission: caseReadPermission,
	}
	decision, err := checker.CheckAccess(ctx, caseRequest)
	if err != nil || ctx.Err() != nil || decision != "allow" {
		return deny
	}
	ids, err := catalog.ListCaseEntryIDs(ctx, caseRequest)
	if err != nil || ctx.Err() != nil || len(ids) > MaxEntries {
		return deny
	}
	// Own the returned selection before handing it to the existing
	// bounded, fail-closed entry/owner/source authorization pipeline.
	return ProjectFromOwnerSnapshotsWithTypedSources(
		ctx, scope, caseID, append([]string(nil), ids...),
		checker, entries, sources,
	)
}
