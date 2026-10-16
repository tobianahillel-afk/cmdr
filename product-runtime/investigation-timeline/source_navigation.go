package investigationtimeline

import (
	"context"
	"unicode"
	"unicode/utf8"
)

// TimelineReturnContext is presentation state, NOT an authorization token.
// A host must re-open the Case and re-check ACLs before restoring it.
type TimelineReturnContext struct {
	CaseID     string `json:"case_id"`
	EntryID    string `json:"entry_id"`
	SourceTime string `json:"source_time"`
	Kind       string `json:"kind"`
	Version    int    `json:"version"`
}

// SourceNavigationDraft is an authorized destination reference, not a URL.
// The owning product remains responsible for its route and further ACL checks.
type SourceNavigationDraft struct {
	OwnerKind     string                `json:"owner_kind"`
	OwnerObjectID string                `json:"owner_object_id"`
	Return        TimelineReturnContext `json:"return_context"`
}

type SourceNavigationResult struct {
	Status string                 `json:"status"`
	Draft  *SourceNavigationDraft `json:"draft,omitempty"`
}

// OpenCaseTimelineSourceNavigation creates a read-only source destination
// draft. It NEVER trusts caller-supplied source IDs, kinds, rights or tenant.
// It reuses the canonical Case -> owner -> typed source authorization path,
// then rechecks the selected destination before handing it to its owner.
func OpenCaseTimelineSourceNavigation(
	ctx context.Context, scope Scope, caseID, entryID string,
	checker AccessChecker, catalog CaseTimelineCatalog,
	entries EntryResolver, sources SourceOwnerResolver,
) SourceNavigationResult {
	deny := SourceNavigationResult{Status: "deny"}
	// Reject untrusted selectors before any provider I/O. Reuse
	// OpenCaseTimeline's existing scope, nil-provider and cancellation gates
	// instead of duplicating that exact authority boundary here.
	if missingAuthority(ctx) || !validNavigationSelector(caseID) ||
		!validNavigationSelector(entryID) {
		return deny
	}
	// A selected entry must have been authorized by the exact same Case
	// catalog and source-owner integration as the Timeline the analyst saw.
	timeline := OpenCaseTimeline(ctx, scope, caseID, checker, catalog, entries, sources)
	if timeline.Status == "deny" || ctx.Err() != nil {
		return deny
	}
	var selected *ProjectedEntry
	for i := range timeline.Entries {
		if timeline.Entries[i].ID == entryID {
			selected = &timeline.Entries[i]
			break
		}
	}
	if selected == nil {
		return deny // no distinction between absent, denied and foreign
	}
	req := func(resource, id, permission string) AccessRequest {
		return AccessRequest{
			Tenant: scope.Tenant, Environment: scope.Environment,
			CaseID: caseID, Resource: resource, ObjectID: id,
			Permission: permission,
		}
	}
	// Resolve destination identity from the trusted owner, not the browser
	// or an inferred SourceRef prefix. Owner metadata is not authority.
	identity, err := sources.ResolveSource(ctx, req("source-reference", selected.SourceRef, ""))
	if err != nil || ctx.Err() != nil ||
		identity.Reference != selected.SourceRef || identity.Tenant != scope.Tenant ||
		identity.Environment != scope.Environment || identity.CaseID != caseID ||
		!validNavigationSelector(identity.ObjectID) {
		return deny
	}
	permission, valid := canonicalSourceReadPermission(identity.Kind)
	if !valid {
		return deny
	}
	// Recheck the destination's own right, then Case and Timeline Entry
	// access after the owner lookup (including revocations during navigation).
	for _, obligation := range []AccessRequest{
		req(identity.Kind, identity.ObjectID, permission),
		req("case", caseID, caseReadPermission),
		req("timeline-entry", entryID, timelineReadPermission),
	} {
		decision, checkErr := checker.CheckAccess(ctx, obligation)
		if checkErr != nil || decision != "allow" || ctx.Err() != nil {
			return deny
		}
	}
	// An entry/source can change between a Timeline page and this read.
	// Refuse a stale link instead of silently navigating to a new owner.
	current, err := entries.ResolveEntry(ctx, req("timeline-entry", entryID, timelineReadPermission))
	if err != nil || ctx.Err() != nil || current.ID != selected.ID ||
		current.Tenant != scope.Tenant || current.SourceRef != selected.SourceRef ||
		current.Version != selected.Version || current.Kind != selected.Kind ||
		current.SourceTime != selected.SourceTime {
		return deny
	}
	draft := &SourceNavigationDraft{
		OwnerKind: identity.Kind, OwnerObjectID: identity.ObjectID,
		Return: TimelineReturnContext{CaseID: caseID, EntryID: entryID,
			SourceTime: selected.SourceTime, Kind: selected.Kind, Version: selected.Version},
	}
	return SourceNavigationResult{Status: "ready", Draft: draft}
}

func validNavigationSelector(value string) bool {
	if invalidScope(value) || len(value) > maxHTTPCaseIDBytes || !utf8.ValidString(value) {
		return false
	}
	for _, c := range value {
		if unicode.IsControl(c) {
			return false
		}
	}
	return true
}
