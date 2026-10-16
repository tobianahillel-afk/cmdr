package investigationtimeline

import (
	"context"
	"strings"
)

// AccessRequest is a server-side authorization obligation. Resource and
// permission are distinct so Case, Shared Timeline Entry and linked sources
// never inherit each other's rights.
type AccessRequest struct {
	Tenant      string
	Environment string
	Resource    string
	ObjectID    string
	Permission  string
}

// AccessChecker is supplied by a trusted server integration, NOT the caller
// providing Input. Its implementation must resolve authoritative decisions.
type AccessChecker interface {
	CheckAccess(context.Context, AccessRequest) (string, error)
}

const (
	caseReadPermission     = "perm.investigate.case.read"
	timelineReadPermission = "perm.shared-capabilities.timeline-entry.read"
	sourceReadPermission   = "source-owner-read"
)

// ProjectWithAccessChecker is the opt-in trusted integration boundary for
// read-only Timeline projections. It deliberately ignores ALL permission
// decisions supplied by the Input; only checker-derived results are used.
// This is not a production provider: none is bundled or configured here.
func ProjectWithAccessChecker(ctx context.Context, in Input, caseID string, checker AccessChecker) Result {
	deny := Result{Status: "deny"}
	if ctx == nil || checker == nil ||
		invalidScope(in.Scope.Tenant) || invalidScope(in.Scope.Environment) ||
		strings.TrimSpace(caseID) == "" || strings.TrimSpace(caseID) == "*" ||
		len(in.Entries) > MaxEntries {
		return deny
	}
	checked := func(resource, objectID, permission string) (string, bool) {
		if ctx.Err() != nil {
			return "", false
		}
		decision, err := checker.CheckAccess(ctx, AccessRequest{
			Tenant: in.Scope.Tenant, Environment: in.Scope.Environment,
			Resource: resource, ObjectID: objectID, Permission: permission,
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
	for _, source := range in.Entries {
		if ctx.Err() != nil {
			return deny
		}
		entry := source                 // never mutate the caller's backing array
		entry.SourcePermission = "deny" // never trust caller-provided allow
		if entry.Tenant != in.Scope.Tenant {
			derived.Entries = append(derived.Entries, entry)
			continue
		}
		// Missing entry IDs cannot be authorized against Shared ownership;
		// keep them invisible rather than using a wildcard authorization.
		if strings.TrimSpace(entry.ID) == "" || strings.TrimSpace(entry.ID) == "*" {
			derived.Entries = append(derived.Entries, entry)
			continue
		}
		timelineAccess, valid := checked("timeline-entry", entry.ID, timelineReadPermission)
		if !valid {
			return deny
		}
		entry.SourcePermission = timelineAccess
		if timelineAccess == "allow" {
			// Source checks are independent. They MUST be resolved by their
			// owner; this generic resource label does not grant any right.
			if strings.TrimSpace(entry.SourceRef) == "" || strings.TrimSpace(entry.SourceRef) == "*" {
				return deny
			}
			sourceAccess, valid := checked("source", entry.SourceRef, sourceReadPermission)
			if !valid {
				return deny
			}
			entry.SourcePermission = sourceAccess
		}
		derived.Entries = append(derived.Entries, entry)
	}
	return Project(derived)
}
