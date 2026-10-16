package investigationtimeline

import (
	"context"
	"unicode/utf8"
)

// AccessRequest is a server-side authorization obligation. Resource and
// permission are distinct so Case, Shared Timeline Entry and linked sources
// never inherit each other's rights.
type AccessRequest struct {
	Tenant      string
	Environment string
	CaseID      string // bind every object authorization to the selected Case
	Resource    string
	ObjectID    string
	Permission  string
}

// AccessChecker is supplied by a trusted server integration, NOT the caller
// providing Input. Its implementation MUST verify authenticated principal,
// tenant/environment, object permission AND that Timeline Entries and linked
// sources belong to the CaseID in AccessRequest. CaseID is a context binding,
// not itself proof that the checked object belongs to that Case.
type AccessChecker interface {
	CheckAccess(context.Context, AccessRequest) (string, error)
}

const (
	caseReadPermission     = "perm.investigate.case.read"
	timelineReadPermission = "perm.shared-capabilities.timeline-entry.read"
	sourceReadPermission   = "source-owner-read"
)

// boundedSourceSelector validates the aggregate bytes forwarded to the
// external source-access authority. It runs only for Timeline Entries that
// passed the independent Timeline Entry permission check. A denied/foreign
// entry never contributes to this external selector budget.
func boundedSourceSelector(used *int, ref string) bool {
	if invalidScope(ref) || len(ref) > MaxInputBytes-*used {
		return false
	}
	*used += len(ref)
	return true
}

// validTimelineEntryBeforeSourceLookup validates the allowed snapshot's
// structural invariants before sending a linked-source selector to an
// external authority. Time parsing is intentionally NOT required here:
// unknown source clocks are legitimate uncertain entries in Project.
func validTimelineEntryBeforeSourceLookup(entry Entry) bool {
	return validKind(entry.Kind) && entry.Provenance != "" &&
		entry.Version >= 1 && entry.ClockUncertaintyMS >= 0 &&
		(entry.Supersedes == "" || entry.Supersedes != entry.ID) &&
		// The final Project stage must not be the first place to discover
		// lossy UTF-8 forensic metadata. Reject before source-provider I/O.
		utf8.ValidString(entry.SourceTime) && utf8.ValidString(entry.Timezone) &&
		utf8.ValidString(entry.Provenance) && utf8.ValidString(entry.Supersedes)
}

// ProjectWithAccessChecker is the opt-in trusted integration boundary for
// read-only Timeline projections. It deliberately ignores ALL permission
// decisions supplied by the Input; only checker-derived results are used.
// This is not a production provider: none is bundled or configured here.
func ProjectWithAccessChecker(ctx context.Context, in Input, caseID string, checker AccessChecker) Result {
	deny := Result{Status: "deny"}
	if ctx == nil || checker == nil ||
		invalidScope(in.Scope.Tenant) || invalidScope(in.Scope.Environment) ||
		invalidScope(caseID) ||
		len(in.Scope.Tenant)+len(in.Scope.Environment)+len(caseID) > MaxInputBytes ||
		len(in.Entries) > MaxEntries {
		return deny
	}
	// This legacy caller-snapshot adapter is not a source of truth. As on
	// the owner-resolved path, bound every object ID before invoking ANY
	// authority: oversized requester-controlled strings must not reach an
	// authorization provider merely because they will later fail Project.
	// Foreign and malformed entries remain opaque; they are not selectors.
	usedIDs := len(in.Scope.Tenant) + len(in.Scope.Environment) + len(caseID)
	usedSourceSelectors := usedIDs
	// A duplicate visible selector must be rejected BEFORE any external
	// authority receives a request. Opaque/foreign/malformed entries must
	// not influence this decision or reveal their own identifiers.
	seenIDs := make(map[string]bool, len(in.Entries))
	for _, entry := range in.Entries {
		if entry.Tenant != in.Scope.Tenant ||
			invalidScope(entry.ID) {
			continue
		}
		if seenIDs[entry.ID] || len(entry.ID) > MaxInputBytes-usedIDs {
			return deny
		}
		seenIDs[entry.ID] = true
		usedIDs += len(entry.ID)
	}
	// All visible Timeline Entry IDs were validated. Charge them to the
	// same aggregate external-selector budget as linked source references.
	// Otherwise individually legal ID and source budgets can exceed the
	// total bound before a provider receives CheckAccess(source, ...).
	usedSourceSelectors = usedIDs
	checked := func(resource, objectID, permission string) (string, bool) {
		if ctx.Err() != nil {
			return "", false
		}
		decision, err := checker.CheckAccess(ctx, AccessRequest{
			Tenant: in.Scope.Tenant, Environment: in.Scope.Environment,
			CaseID: caseID, Resource: resource, ObjectID: objectID, Permission: permission,
		})
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
	caseAccess, ok := checked("case", caseID, caseReadPermission)
	if !ok || caseAccess != "allow" {
		return deny
	}
	derived := Input{
		Scope: Scope{
			Tenant: in.Scope.Tenant, Environment: in.Scope.Environment,
			CasePermission: "allow", TimelinePermission: "allow",
		},
		Entries: make([]Entry, 0, len(in.Entries)),
	}
	// Finish all Timeline Entry access decisions before contacting any linked
	// source authority. A later provider failure must not leak earlier source
	// selectors to that authority when the projection will be denied anyway.
	decisions := make([]string, len(in.Entries))
	for i, entry := range in.Entries {
		if ctx.Err() != nil {
			return deny
		}
		if entry.Tenant != in.Scope.Tenant || invalidScope(entry.ID) {
			continue // opaque entries are not selectors for this Case
		}
		decision, valid := checked("timeline-entry", entry.ID, timelineReadPermission)
		if !valid {
			return deny
		}
		decisions[i] = decision
	}

	// Validate ALL potentially visible snapshots and the cumulative source
	// selector budget before making even one external source permission call.
	// This reuses the owner-resolved adapter's three-pass trust ordering.
	for i, source := range in.Entries {
		if ctx.Err() != nil {
			return deny
		}
		entry := source
		entry.SourcePermission = "deny" // disregard caller-supplied allow
		if entry.Tenant != in.Scope.Tenant || invalidScope(entry.ID) {
			derived.Entries = append(derived.Entries, entry)
			continue
		}
		entry.SourcePermission = decisions[i]
		if decisions[i] == "allow" {
			if !validTimelineEntryBeforeSourceLookup(entry) ||
				!boundedSourceSelector(&usedSourceSelectors, entry.SourceRef) {
				return deny
			}
			entry.SourcePermission = "unknown"
		}
		derived.Entries = append(derived.Entries, entry)
	}

	for i := range derived.Entries {
		if decisions[i] != "allow" {
			continue
		}
		sourceAccess, valid := checked("source", derived.Entries[i].SourceRef, sourceReadPermission)
		if !valid {
			return deny
		}
		derived.Entries[i].SourcePermission = sourceAccess
	}
	if ctx.Err() != nil {
		return deny
	}
	return Project(derived)
}
