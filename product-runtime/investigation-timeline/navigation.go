package investigationtimeline

import "context"

// SourceNavigationDraft is a read-only routing draft, NOT a URL or permission
// token. Hosts resolve and authorize the destination independently at open.
type SourceNavigationDraft struct {
	Status        string `json:"status"`
	Kind          string `json:"kind"`
	ObjectID      string `json:"object_id"`
	ReturnCaseID  string `json:"return_case_id"`
	ReturnEntryID string `json:"return_entry_id"`
	SourceTime    string `json:"source_time"`
	EntryVersion  int    `json:"entry_version"`
}

// DraftCaseSourceNavigation reuses the full Case -> Timeline -> source owner
// projection, then RECHECKS Case, Entry and canonical destination permissions.
// A browser-supplied sourceRef, destination kind, or permission is never used.
func DraftCaseSourceNavigation(ctx context.Context, scope Scope, caseID, entryID string,
	checker AccessChecker, catalog CaseTimelineCatalog, entries EntryResolver,
	sources SourceOwnerResolver) (SourceNavigationDraft, bool) {
	deny := SourceNavigationDraft{}
	if missingAuthority(ctx) || missingAuthority(checker) ||
		missingAuthority(catalog) || missingAuthority(entries) || missingAuthority(sources) ||
		invalidScope(caseID) || invalidScope(entryID) || len(entryID) > 1024 ||
		ctx.Err() != nil {
		return deny, false
	}
	timeline := OpenCaseTimeline(ctx, scope, caseID, checker, catalog, entries, sources)
	if timeline.Status == "deny" || ctx.Err() != nil {
		return deny, false
	}
	var selected *ProjectedEntry
	for i := range timeline.Entries {
		if timeline.Entries[i].ID == entryID {
			selected = &timeline.Entries[i]
			break
		}
	}
	if selected == nil {
		return deny, false
	}
	check := func(resource, object, permission string) bool {
		if ctx.Err() != nil {
			return false
		}
		result, err := checker.CheckAccess(ctx, AccessRequest{
			Tenant: scope.Tenant, Environment: scope.Environment, CaseID: caseID,
			Resource: resource, ObjectID: object, Permission: permission,
		})
		return err == nil && ctx.Err() == nil && result == "allow"
	}
	// Permission may change between the timeline listing and the click.
	if !check("case", caseID, caseReadPermission) ||
		!check("timeline-entry", entryID, timelineReadPermission) {
		return deny, false
	}
	// Re-resolve the owner snapshot to reject a stale version/source/time.
	owner, err := entries.ResolveEntry(ctx, AccessRequest{
		Tenant: scope.Tenant, Environment: scope.Environment, CaseID: caseID,
		Resource: "timeline-entry", ObjectID: entryID, Permission: timelineReadPermission,
	})
	if err != nil || ctx.Err() != nil || owner.ID != entryID ||
		owner.Tenant != scope.Tenant || owner.Version != selected.Version ||
		owner.SourceRef != selected.SourceRef || owner.SourceTime != selected.SourceTime ||
		owner.Kind != selected.Kind || !validTimelineEntryBeforeSourceLookup(owner) {
		return deny, false
	}
	// SourceRef is trusted only as an opaque selector for the owner resolver.
	identity, err := sources.ResolveSource(ctx, AccessRequest{
		Tenant: scope.Tenant, Environment: scope.Environment, CaseID: caseID,
		Resource: "source-reference", ObjectID: owner.SourceRef,
	})
	if err != nil || ctx.Err() != nil || identity.Reference != owner.SourceRef ||
		identity.Tenant != scope.Tenant || identity.Environment != scope.Environment ||
		identity.CaseID != caseID || invalidScope(identity.ObjectID) ||
		len(identity.ObjectID) > 1024 {
		return deny, false
	}
	permission, supported := canonicalSourceReadPermission(identity.Kind)
	if !supported || !check(identity.Kind, identity.ObjectID, permission) {
		return deny, false
	}
	return SourceNavigationDraft{
		Status: "ready", Kind: identity.Kind, ObjectID: identity.ObjectID,
		ReturnCaseID: caseID, ReturnEntryID: entryID,
		SourceTime: selected.SourceTime, EntryVersion: selected.Version,
	}, true
}
