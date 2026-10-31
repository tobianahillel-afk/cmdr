package workqueue

import (
	"context"
	"errors"
	"time"
)

// IncidentTaskReference is a scoped relationship already authorized by
// Command's relation owner. It conveys no permission to read the Task.
type IncidentTaskReference struct {
	IncidentID  string
	TaskID      string
	Tenant      string
	Environment string
}

// IncidentTaskReferencesOwner independently authorizes links from an Incident
// to existing Command-owned Tasks for the current principal and scope.
type IncidentTaskReferencesOwner interface {
	ReadAuthorizedIncidentTaskReferences(ctx context.Context, principal, tenant, environment, incidentID string) ([]IncidentTaskReference, error)
}

// AuthorizedIncidentTask is an existing Work Queue Task projection selected by
// Command Task.read, not a new Task schema, editable Task, or assignment right.
type AuthorizedIncidentTask struct {
	ID        string
	Owner     string
	Status    string
	UpdatedAt string
}

var errIncidentTaskLinks = errors.New("incident task projection unavailable")

// readAuthorizedIncidentTasks reuses the existing Command Owner.List/CanRead
// authorization instead of adding a parallel Task access-control engine.
// Missing or individually forbidden Tasks never reveal IDs or counts.
func readAuthorizedIncidentTasks(ctx context.Context, owner Owner, scope SessionScope, incidentID string) ([]AuthorizedIncidentTask, error) {
	links, ok := owner.(IncidentTaskReferencesOwner)
	if !ok {
		return nil, nil // Optional partial source: not evidence of zero Tasks.
	}
	if ctx == nil || ctx.Err() != nil {
		return nil, errIncidentTaskLinks
	}
	refs, err := links.ReadAuthorizedIncidentTaskReferences(ctx,
		scope.Principal, scope.Tenant, scope.Environment, incidentID)
	if err != nil || ctx.Err() != nil || len(refs) > 12 {
		return nil, errIncidentTaskLinks
	}
	if len(refs) == 0 {
		return nil, nil
	}
	seen := make(map[string]bool, len(refs))
	for _, ref := range refs {
		if ref.IncidentID != incidentID || ref.Tenant != scope.Tenant ||
			ref.Environment != scope.Environment || invalidScope(ref.TaskID) ||
			!validOverviewText(ref.TaskID, 128) || seen[ref.TaskID] {
			return nil, errIncidentTaskLinks
		}
		seen[ref.TaskID] = true
	}
	// A second source read is deliberately marked as a new nontransactional
	// owner snapshot; it cannot borrow the previous Incident.read grant.
	items, err := owner.List(ctx, scope.Tenant, scope.Environment)
	if err != nil || ctx.Err() != nil || len(items) > 5000 {
		return nil, errIncidentTaskLinks
	}
	found := make(map[string]Item, len(refs))
	for _, item := range items {
		if item.Kind != Task || !seen[item.ID] {
			continue
		}
		if _, duplicate := found[item.ID]; duplicate ||
			item.Tenant != scope.Tenant || item.Environment != scope.Environment {
			return nil, errIncidentTaskLinks
		}
		found[item.ID] = item
	}
	result := make([]AuthorizedIncidentTask, 0, len(refs))
	for _, ref := range refs {
		item, present := found[ref.TaskID]
		if !present {
			continue // No existence inference from a missing owned Task.
		}
		allowed, err := owner.CanRead(ctx, scope.Principal, scope.Tenant,
			scope.Environment, Task, ref.TaskID)
		if err != nil || ctx.Err() != nil {
			return nil, errIncidentTaskLinks
		}
		if !allowed {
			continue // Incident.read and relation.read do not grant Task.read.
		}
		if !validOverviewText(item.Owner, 128) ||
			!validOverviewText(item.Status, 64) {
			return nil, errIncidentTaskLinks
		}
		updated, err := time.Parse(time.RFC3339Nano, item.UpdatedAt)
		if err != nil || updated.IsZero() {
			return nil, errIncidentTaskLinks
		}
		result = append(result, AuthorizedIncidentTask{
			ID: item.ID, Owner: item.Owner, Status: item.Status,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}
