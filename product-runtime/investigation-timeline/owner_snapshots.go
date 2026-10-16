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
	return projectFromOwnerSnapshots(ctx, scope, caseID, ids, checker, resolver, nil)
}

// ProjectFromOwnerSnapshotsWithTypedSources is an opt-in, owner-routed source
// authorization boundary. The source resolver and checker must both be
// authoritative, trusted server-side integrations; no production provider is
// bundled. The legacy generic-source adapter above remains unchanged.
func ProjectFromOwnerSnapshotsWithTypedSources(
	ctx context.Context, scope Scope, caseID string, ids []string,
	checker AccessChecker, resolver EntryResolver, sourceResolver SourceOwnerResolver,
) Result {
	if missingAuthority(sourceResolver) {
		return Result{Status: "deny"}
	}
	return projectFromOwnerSnapshots(ctx, scope, caseID, ids, checker, resolver, sourceResolver)
}

func projectFromOwnerSnapshots(
	ctx context.Context, scope Scope, caseID string, ids []string,
	checker AccessChecker, resolver EntryResolver, sourceResolver SourceOwnerResolver,
) Result {
	deny := Result{Status: "deny"}
	if missingAuthority(ctx) || missingAuthority(checker) || missingAuthority(resolver) ||
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
	// The entire requester-supplied selector set is bounded independently.
	// Only IDs confirmed allowed by the Timeline Entry authority may later
	// consume the linked-source selector budget.
	allowed, valid := check(request("case", caseID, caseReadPermission))
	if !valid || allowed != "allow" {
		return deny
	}
	out := Input{
		Scope: Scope{Tenant: scope.Tenant, Environment: scope.Environment,
			CasePermission: "allow", TimelinePermission: "allow"},
		Entries: make([]Entry, 0, len(ids)),
	}
	// Complete every Timeline Entry permission check before any owner
	// resolution or linked-source check. A later provider error must not
	// cause earlier source metadata to be fetched or forwarded needlessly.
	// This is a bounded two-pass check, NOT an atomic external ACL snapshot.
	decisions := make([]string, len(ids))
	for i, id := range ids {
		decision, valid := check(request("timeline-entry", id, timelineReadPermission))
		if !valid {
			return deny
		}
		decisions[i] = decision
	}
	// Denied Timeline Entry IDs must never reduce the source budget of an
	// authorized entry. This counter is distinct from the strict all-ID
	// request bound applied BEFORE any authority I/O.
	usedSourceSelectors = len(scope.Tenant) + len(scope.Environment) + len(caseID)
	for i, id := range ids {
		if decisions[i] == "allow" {
			usedSourceSelectors += len(id)
		}
	}
	// Resolve and validate ALL permitted owner snapshots before any linked
	// source access checks. An invalid late owner response must not cause
	// earlier source metadata to be sent to an external ACL provider.
	for i, id := range ids {
		// An ACL provider can return its final allow immediately before the
		// request is cancelled. Never dispatch owner I/O after cancellation.
		if ctx.Err() != nil {
			return deny
		}
		decision := decisions[i]
		if decision != "allow" {
			out.Entries = append(out.Entries, Entry{
				ID: id, Tenant: scope.Tenant, SourcePermission: decision,
			})
			continue
		}
		timelineRequest := request("timeline-entry", id, timelineReadPermission)
		entry, err := resolver.ResolveEntry(ctx, timelineRequest)
		if err != nil || ctx.Err() != nil ||
			entry.ID != id || entry.Tenant != scope.Tenant ||
			!validTimelineEntryBeforeSourceLookup(entry) ||
			!boundedSourceSelector(&usedSourceSelectors, entry.SourceRef) {
			return deny
		}
		// Permission returned by the owner is untrusted. Only the
		// external authoritative source check below can set it.
		entry.SourcePermission = "unknown"
		out.Entries = append(out.Entries, entry)
	}
	// For strict typed integration, resolve ALL trusted source identities
	// before any source ACL check, so a malformed later source cannot lead
	// to partial release of earlier source metadata to an ACL authority.
	sourceRequests := make([]AccessRequest, len(out.Entries))
	// A trusted source reference must resolve to one canonical owner identity
	// throughout this projection. A changing result is not safe to authorize
	// as two different objects under one displayed SourceRef. Compare all
	// observations before dispatching any linked-source ACL checks.
	resolvedSourceIdentities := make(map[string]SourceOwnerIdentity)
	// Canonical IDs returned by owner resolvers are new selectors and must
	// have their own cumulative bound before any typed source ACL I/O.
	usedTypedIDs := len(scope.Tenant) + len(scope.Environment) + len(caseID)
	for i := range out.Entries {
		if decisions[i] != "allow" {
			continue
		}
		if ctx.Err() != nil {
			return deny
		}
		if sourceResolver == nil {
			sourceRequests[i] = request("source", out.Entries[i].SourceRef, sourceReadPermission)
			continue // transitional generic source authority
		}
		lookup := request("source-reference", out.Entries[i].SourceRef, "")
		identity, err := sourceResolver.ResolveSource(ctx, lookup)
		if err != nil || ctx.Err() != nil ||
			identity.Reference != out.Entries[i].SourceRef ||
			identity.CaseID != caseID || identity.Tenant != scope.Tenant ||
			identity.Environment != scope.Environment ||
			invalidScope(identity.ObjectID) {
			return deny
		}
		permission, ok := canonicalSourceReadPermission(identity.Kind)
		if !ok || len(identity.ObjectID) > MaxInputBytes-usedTypedIDs {
			return deny
		}
		if previous, exists := resolvedSourceIdentities[identity.Reference]; exists && previous != identity {
			return deny
		}
		resolvedSourceIdentities[identity.Reference] = identity
		usedTypedIDs += len(identity.ObjectID)
		sourceRequests[i] = request(identity.Kind, identity.ObjectID, permission)
	}
	for i := range out.Entries {
		if decisions[i] != "allow" {
			continue
		}
		sourceDecision, valid := check(sourceRequests[i])
		if !valid {
			return deny
		}
		out.Entries[i].SourcePermission = sourceDecision
	}
	if ctx.Err() != nil {
		return deny
	}
	return Project(out)
}
