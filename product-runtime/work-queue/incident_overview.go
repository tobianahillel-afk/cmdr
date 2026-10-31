package workqueue

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// AuthorizedIncidentOverview is a Command-owned read projection, not an
// Incident mutation or a client-supplied permission. Every non-nil response
// must independently match the authenticated Incident and its source scope.
type AuthorizedIncidentOverview struct {
	IncidentID      string
	Tenant          string
	Environment     string
	Title           string
	Priority        string
	Severity        string
	BusinessImpact  string
	NextAction      string
	NextActionDueAt string
	SourceRef       string
	SourceVersion   string
	ObservedAt      string
}

// IncidentOverviewOwner is optional. Absence means "source unavailable", never
// a fabricated title, priority, severity, impact or next action. It does not grant
// Case, Evidence, Alert, Response Run or export rights.
type IncidentOverviewOwner interface {
	ReadAuthorizedIncidentOverview(ctx context.Context, principal, tenant, environment, incidentID string) (*AuthorizedIncidentOverview, error)
}

var errIncidentOverview = errors.New("incident overview source unavailable")

func validOverviewText(s string, maxBytes int) bool {
	if len(s) > maxBytes || !utf8.ValidString(s) {
		return false
	}
	return strings.IndexFunc(s, func(r rune) bool {
		return r < 0x20 || r == 0x7f
	}) < 0
}

func validateIncidentOverview(raw *AuthorizedIncidentOverview, scope SessionScope, incidentID string) (*AuthorizedIncidentOverview, error) {
	if raw == nil {
		return nil, nil
	}
	if raw.IncidentID != incidentID || raw.Tenant != scope.Tenant ||
		raw.Environment != scope.Environment || invalidScope(raw.SourceRef) ||
		invalidScope(raw.SourceVersion) || !validOverviewText(raw.SourceRef, 128) ||
		!validOverviewText(raw.SourceVersion, 80) ||
		!validOverviewText(raw.Title, 256) ||
		!validOverviewText(raw.Priority, 64) ||
		!validOverviewText(raw.Severity, 64) ||
		!validOverviewText(raw.BusinessImpact, 512) ||
		!validOverviewText(raw.NextAction, 256) {
		return nil, errIncidentOverview
	}
	observed, err := time.Parse(time.RFC3339Nano, raw.ObservedAt)
	if err != nil || observed.IsZero() {
		return nil, errIncidentOverview
	}
	// Due time belongs to the independently authorized Command source. A
	// missing due time is unknown, never a promise that no deadline exists.
	// A due date without an action is an inconsistent partial source.
	if raw.NextActionDueAt != "" {
		if raw.NextAction == "" || !validOverviewText(raw.NextActionDueAt, 64) {
			return nil, errIncidentOverview
		}
		due, err := time.Parse(time.RFC3339Nano, raw.NextActionDueAt)
		if err != nil || due.IsZero() {
			return nil, errIncidentOverview
		}
	}
	copyValue := *raw
	return &copyValue, nil
}

func readAuthorizedIncidentOverview(ctx context.Context, owner Owner, scope SessionScope, incidentID string) (*AuthorizedIncidentOverview, error) {
	source, ok := owner.(IncidentOverviewOwner)
	if !ok {
		return nil, nil
	}
	raw, err := source.ReadAuthorizedIncidentOverview(ctx, scope.Principal,
		scope.Tenant, scope.Environment, incidentID)
	if err != nil || ctx.Err() != nil {
		return nil, errIncidentOverview
	}
	return validateIncidentOverview(raw, scope, incidentID)
}
