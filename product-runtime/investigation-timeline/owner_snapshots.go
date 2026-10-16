package investigationtimeline

import "context"

// EntryResolver loads a canonical Timeline Entry from a trusted server-side
// source of truth. It must bind the authenticated principal, requested Case,
// tenant and environment; a caller-provided Entry is never an authority.
type EntryResolver interface {
	ResolveEntry(context.Context, AccessRequest) (Entry, error)
}

// ProjectFromOwnerSnapshots takes only selected Timeline Entry IDs from the
// requester. The trusted owner resolver supplies every projected field after
// the Case and Timeline Entry access checks; linked source access is checked
// separately. No concrete provider or persistence layer is bundled.
func ProjectFromOwnerSnapshots(
	ctx context.Context, scope Scope, caseID string, ids []string,
	checker AccessChecker, resolver EntryResolver,
) Result {
	deny := Result{Status: "deny"}
	if ctx == nil || checker == nil || resolver == nil ||
		invalidScope(scope.Tenant) || invalidScope(scope.Environment) ||
		invalidScope(caseID) || len(ids) > MaxEntries {
		return deny
	}
	request := func(resource, id, permission string) AccessRequest {
		return AccessRequest{
			Tenant: scope.Tenant, Environment: scope.Environment,
			CaseID: caseID, Resource: resource, ObjectID: id, Permission: permission,
		}
	}
	check := func(r AccessRequest) (string, bool) {
		if ctx.Err() != nil {
			return "", false
		}
		decision, err := checker.CheckAccess(ctx, r)
		if err != nil || ctx.Err() != nil {
			return "", false
		}
		switch decision {
		case "allow", "deny", "unknown", "missing":
			return decision, true
		default:
			return "", false
		}
	}
	// Validate the entire requested ID set before consulting ANY authority.
	// A malformed/duplicate later selector must not trigger earlier owner reads.
	seen := make(map[string]bool, len(ids))
	usedIDs := len(scope.Tenant) + len(scope.Environment) + len(caseID)
	usedSourceSelectors := usedIDs
	if usedIDs > MaxInputBytes {
		return deny // bound scope and Case ID even if no entries were selected
	}
	for _, id := range ids {
		usedIDs += len(id)
		if invalidScope(id) || seen[id] || usedIDs > MaxInputBytes {
			return deny
		}
		seen[id] = true
	}
	// Source selectors share the envelope with every validated requested ID;
	// keeping a separate starting counter would miss the ID bytes above.
	usedSourceSelectors = usedIDs
	allowed, valid := check(request("case", caseID, caseReadPermission))
	if !valid || allowed != "allow" {
		return deny
	}
	out := Input{
		Scope: Scope{Tenant: scope.Tenant, Environment: scope.Environment,
			CasePermission: "allow", TimelinePermission: "allow"},
		Entries: make([]Entry, 0, len(ids)),
	}
	for _, id := range ids {
		timelineRequest := request("timeline-entry", id, timelineReadPermission)
		decision, valid := check(timelineRequest)
		if !valid {
			return deny // provider error or unexpected response: atomic failure
		}
		if decision != "allow" {
			// A denied ID cannot be resolved or used in a linked-source check.
			out.Entries = append(out.Entries, Entry{
				ID: id, Tenant: scope.Tenant, SourcePermission: decision,
			})
			continue
		}
		entry, err := resolver.ResolveEntry(ctx, timelineRequest)
		if err != nil || ctx.Err() != nil ||
			entry.ID != id || entry.Tenant != scope.Tenant ||
			!validTimelineEntryBeforeSourceLookup(entry) ||
			!boundedSourceSelector(&usedSourceSelectors, entry.SourceRef) {
			return deny
		}
		sourceDecision, valid := check(request("source", entry.SourceRef, sourceReadPermission))
		if !valid {
			return deny
		}
		entry.SourcePermission = sourceDecision // ignore any snapshot-provided permission flag
		out.Entries = append(out.Entries, entry)
	}
	if ctx.Err() != nil {
		return deny
	}
	return Project(out)
}
