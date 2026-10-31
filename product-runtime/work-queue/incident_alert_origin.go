package workqueue

import (
	"context"
	"errors"
	"time"
)

// IncidentAlertReference is an Incident-to-Alert link resolved by the Command
// relation owner, not by browser-selected IDs or a child Alert read permission.
type IncidentAlertReference struct {
	IncidentID string
	AlertID string
	Tenant string
	Environment string
}

// IncidentAlertReferencesOwner must authorize the Incident's relation set.
// The returned references alone NEVER confer permission to read an Alert.
type IncidentAlertReferencesOwner interface {
	ReadAuthorizedIncidentAlertReferences(ctx context.Context, principal, tenant, environment, incidentID string) ([]IncidentAlertReference, error)
}

// AuthorizedIncidentAlert is a narrow, owner-provided, Alert.read-authorized
// projection: NOT raw Event, Signal, Evidence or an incident mutation.
type AuthorizedIncidentAlert struct {
	ID string
	Tenant string
	Environment string
	Severity string
	Status string
	Disposition string
	SourceRef string
	SourceVersion string
	ObservedAt string
}

// IncidentAlertOwner independently checks perm.command.alert.read for each
// referenced Alert and may return nil for an inaccessible Alert.
type IncidentAlertOwner interface {
	ReadAuthorizedIncidentAlert(ctx context.Context, principal, tenant, environment, alertID string) (*AuthorizedIncidentAlert, error)
}

var errIncidentAlertOrigin = errors.New("incident alert origin unavailable")

func readAuthorizedIncidentAlertOrigins(ctx context.Context, owner Owner, scope SessionScope, incidentID string) ([]AuthorizedIncidentAlert, error) {
	links, okLinks := owner.(IncidentAlertReferencesOwner)
	alerts, okAlerts := owner.(IncidentAlertOwner)
	// Partial is explicit, not an inferred claim that the Incident has no Alerts.
	if !okLinks || !okAlerts {
		return nil, nil
	}
	if ctx.Err() != nil {
		return nil, errIncidentAlertOrigin
	}
	references, err := links.ReadAuthorizedIncidentAlertReferences(ctx, scope.Principal, scope.Tenant, scope.Environment, incidentID)
	if err != nil || ctx.Err() != nil || len(references) > 20 {
		return nil, errIncidentAlertOrigin
	}
	seen := make(map[string]bool, len(references))
	out := make([]AuthorizedIncidentAlert, 0, len(references))
	for _, ref := range references {
		// This verifies the relation identity before the independently guarded
		// Alert owner is contacted. No foreign, duplicate or hostile IDs.
		if ref.IncidentID != incidentID || ref.Tenant != scope.Tenant ||
			ref.Environment != scope.Environment || invalidScope(ref.AlertID) ||
			!validOverviewText(ref.AlertID, 128) || seen[ref.AlertID] {
			return nil, errIncidentAlertOrigin
		}
		seen[ref.AlertID] = true
		record, err := alerts.ReadAuthorizedIncidentAlert(ctx, scope.Principal, scope.Tenant, scope.Environment, ref.AlertID)
		if err != nil || ctx.Err() != nil {
			return nil, errIncidentAlertOrigin
		}
		if record == nil {
			continue // Alert.read denied: neither ID nor existence is disclosed.
		}
		if record.ID != ref.AlertID || record.Tenant != scope.Tenant ||
			record.Environment != scope.Environment ||
			invalidScope(record.SourceRef) || invalidScope(record.SourceVersion) ||
			!validOverviewText(record.SourceRef, 128) ||
			!validOverviewText(record.SourceVersion, 80) ||
			!validOverviewText(record.Severity, 64) ||
			!validOverviewText(record.Status, 64) ||
			!validOverviewText(record.Disposition, 256) {
			return nil, errIncidentAlertOrigin
		}
		observed, err := time.Parse(time.RFC3339Nano, record.ObservedAt)
		if err != nil || observed.IsZero() {
			return nil, errIncidentAlertOrigin
		}
		out = append(out, *record)
	}
	return out, nil
}
