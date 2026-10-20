// Package workqueue implements a bounded, read-only Command-owned Work Queue
// projection. It does not create Incidents, Tasks, Saved Views or permissions.
package workqueue

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type Kind string

const (
	Incident Kind = "Incident"
	Task     Kind = "Task"
)

// Exactly the six Command-owned system view keys in the canonical catalogue.
// These are configurations of ONE queue, not six separately routed screens.
type View string

const (
	All        View = "all"
	Incidents  View = "incidents"
	Tasks      View = "tasks"
	Unassigned View = "unassigned"
	SLARisk    View = "sla-risk"
	MyWork     View = "my-work"
)

var SystemViews = [...]View{All, Incidents, Tasks, Unassigned, SLARisk, MyWork}

// Item is the minimal, already owner-resolved Command Incident/Task projection.
// Provider owns these attributes. Linked Case, Decision, Run and raw evidence
// are intentionally absent: relationships never transfer authorization.
type Item struct {
	Kind         Kind
	ID           string
	Tenant       string
	Environment  string
	Owner        string
	Followers    []string
	Status       string
	UpdatedAt    string
	SLAAtRisk    *bool // nil: unavailable from the authoritative policy source
	SourceFreshness *SourceFreshness // optional, host-owned metadata; never client scope
}

// Owner is an authoritative host adapter, NOT a client authorization flag.
// CanRead MUST evaluate the current principal, tenant, environment AND object.
type Owner interface {
	List(ctx context.Context, tenant, environment string) ([]Item, error)
	CanRead(ctx context.Context, principal, tenant, environment string, kind Kind, id string) (bool, error)
}

type Request struct {
	Principal   string
	Tenant      string
	Environment string
	View        View
	// Optional exact owner-provided status filter (Class-0 read-only).
	StatusFilter string
	// Exact responsible owner identifier (Class-0), applied after CanRead.
	OwnerFilter string
	// Read-only chronological order. Blank preserves the original newest-first
	// presentation; no arbitrary client sort column or Shared view schema.
	SortOrder string
	// Cursor is an untrusted, non-authorizing hint for a later authorized page.
	Cursor string
}

type Row struct {
	Kind      Kind   `json:"kind"`
	ID        string `json:"id"`
	Owner     string `json:"owner,omitempty"`
	Status    string `json:"status,omitempty"`
	UpdatedAt string `json:"updated_at"`
	Freshness *FreshnessProjection `json:"freshness,omitempty"`
}

type Result struct {
	View         View     `json:"view"`
	Tenant       string   `json:"tenant"`
	Environment  string   `json:"environment"`
	State        string   `json:"state"`
	Diagnostics  []string `json:"diagnostics"`
	Rows         []Row    `json:"rows"`
	StatusFilter string   `json:"status_filter,omitempty"`
	OwnerFilter string `json:"owner_filter,omitempty"`
	SortOrder string `json:"sort,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	PreviousCursor string `json:"previous_cursor,omitempty"`
	HasPrevious bool `json:"has_previous,omitempty"`
	// Internal HTML-only current-page position. Populated exclusively after
	// the full owner-authorized cursor hash check, never from raw HTTP input.
	// It must not extend the public JSON projection contract.
	CurrentCursor string `json:"-"`
	// Selected is returned ONLY by the authenticated HTTP route, and only
	// for an item reauthorized and present in its currently visible page.
	// No selected field is serialized for legacy/unselected projections.
	Selected *Row `json:"selected,omitempty"`
	// Selected-only CAP-CMD-104 class-0 coordination factors.
	SelectedFactors *AuthorizedFactors `json:"selected_factors,omitempty"`
	// Present only if at least one authorized matching item has host freshness.
	ViewFreshness *ViewFreshness `json:"view_freshness,omitempty"`
}

var (
	ErrScope          = errors.New("work queue requires explicit principal, tenant and environment")
	ErrView           = errors.New("unknown canonical Command system view")
	ErrIncomplete     = errors.New("untrusted or incomplete work item owner data")
	ErrOwner          = errors.New("work queue owner authorization unavailable")
	ErrTooMany        = errors.New("work queue owner snapshot exceeds local bound")
)

// Project is deliberately bounded and all-or-nothing for owner failures.
// Filtering is applied ONLY after an individual first-party authorization for
// each in-scope object. The caller cannot assert Read=true in an input row.
func Project(ctx context.Context, owner Owner, req Request) (Result, error) {
	var empty Result
	if ctx == nil || ctx.Err() != nil || owner == nil {
		return empty, ErrOwner
	}
	if invalidScope(req.Principal) || invalidScope(req.Tenant) || invalidScope(req.Environment) {
		return empty, ErrScope
	}
	// Only exact, bounded owner status values are accepted. No query DSL,
	// wildcard, column permissions, or client-supplied authorization scope.
	if len(req.StatusFilter) > 64 ||
		(req.StatusFilter != "" && (invalidScope(req.StatusFilter) || !utf8.ValidString(req.StatusFilter))) {
		return empty, ErrView
	}
	if len(req.OwnerFilter) > 128 || (req.OwnerFilter != "" &&
		(invalidScope(req.OwnerFilter) || !utf8.ValidString(req.OwnerFilter))) {
		return empty, ErrView
	}
	if req.SortOrder != "" && req.SortOrder != "newest" && req.SortOrder != "oldest" {
		return empty, ErrView
	}
	if _, err := decodeQueueCursor(req.Cursor); err != nil {
		return empty, ErrCursor
	}
	switch req.View {
	case All, Incidents, Tasks, Unassigned, SLARisk, MyWork:
	default:
		return empty, ErrView
	}
	items, err := owner.List(ctx, req.Tenant, req.Environment)
	if err != nil {
		return empty, ErrOwner
	}
	if len(items) > 5000 {
		return empty, ErrTooMany
	}
	type candidate struct {
		row Row
		updated time.Time
	}
	var selected []candidate
	seen := map[string]bool{}
	partialSLA := false
	partialSourceFreshness := false
	now := time.Now().UTC()
	for _, item := range items {
		if ctx.Err() != nil {
			return empty, ErrOwner
		}
		// A host might batch owner records across scopes. Do not query ACL
		// for a cross-tenant identifier or expose even its presence.
		if item.Tenant != req.Tenant || item.Environment != req.Environment {
			continue
		}
		if (item.Kind != Incident && item.Kind != Task) || invalidScope(item.ID) ||
			strings.TrimSpace(item.UpdatedAt) == "" {
			return empty, ErrIncomplete
		}
		key := string(item.Kind) + "\x00" + item.ID
		if seen[key] {
			return empty, ErrIncomplete
		}
		seen[key] = true
		allowed, err := owner.CanRead(ctx, req.Principal, req.Tenant, req.Environment, item.Kind, item.ID)
		if err != nil {
			return empty, ErrOwner
		}
		if !allowed {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, item.UpdatedAt)
		if err != nil || t.IsZero() {
			return empty, ErrIncomplete
		}
		accept := false
		switch req.View {
		case All:
			accept = true
		case Incidents:
			accept = item.Kind == Incident
		case Tasks:
			accept = item.Kind == Task
		case Unassigned:
			accept = item.Owner == ""
		case MyWork:
			accept = item.Owner == req.Principal
			for _, follower := range item.Followers {
				if follower == req.Principal {
					accept = true
				}
			}
		case SLARisk:
			if item.SLAAtRisk == nil {
				// Missing SLA policy degrades only the CURRENT status-filtered
				// cohort. A filtered-out item cannot make this view partial.
				if (req.StatusFilter == "" || item.Status == req.StatusFilter) &&
					(req.OwnerFilter == "" || item.Owner == req.OwnerFilter) {
					partialSLA = true
				}
			} else {
				accept = *item.SLAAtRisk
			}
		}
		if accept && (req.StatusFilter == "" || item.Status == req.StatusFilter) &&
			(req.OwnerFilter == "" || item.Owner == req.OwnerFilter) {
			// An excluded owner cannot degrade the visible freshness or
			// cause a fatal source-metadata validation error.
			freshness, err := sourceFreshness(item.SourceFreshness, now)
			if err != nil { return empty, ErrIncomplete }
			if freshness != nil && (freshness.State == "source-unavailable" || freshness.State == "partial") {
				partialSourceFreshness = true
			}
			selected = append(selected, candidate{
				row: Row{Kind: item.Kind, ID: item.ID, Owner: item.Owner, Status: item.Status, UpdatedAt: item.UpdatedAt, Freshness: freshness},
				updated: t,
			})
		}
	}
	sort.Slice(selected, func(i, j int) bool {
		if !selected[i].updated.Equal(selected[j].updated) {
			if req.SortOrder == "oldest" {
				return selected[i].updated.Before(selected[j].updated)
			}
			return selected[i].updated.After(selected[j].updated)
		}
		if selected[i].row.Kind != selected[j].row.Kind {
			return selected[i].row.Kind < selected[j].row.Kind
		}
		return selected[i].row.ID < selected[j].row.ID
	})
	result := Result{View: req.View, Tenant: req.Tenant, Environment: req.Environment,
		State: "ready", Rows: make([]Row, 0, len(selected)), Diagnostics: []string{},
		StatusFilter: req.StatusFilter, OwnerFilter: req.OwnerFilter, SortOrder: req.SortOrder}
	if len(selected) == 0 {
		result.State = "no-authorized-items"
	}
	if partialSLA || partialSourceFreshness {
		result.State = "partial-sources"
	}
	// Summarize the full freshly authorized, post-view/filter cohort,
	// not merely the visible 50-row page. The current-page cursor must
	// still validate before any JSON/HTML result is returned.
	viewFreshness := ViewFreshness{Total:len(selected)}
	hasOwnerMetadata := false
	for _, item := range selected {
		if item.row.Freshness == nil {
			viewFreshness.Unknown++
			continue
		}
		hasOwnerMetadata = true
		switch item.row.Freshness.State {
		case "current": viewFreshness.Current++
		case "aging": viewFreshness.Aging++
		case "stale": viewFreshness.Stale++
		case "partial": viewFreshness.Partial++
		case "source-unavailable": viewFreshness.Unavailable++
		default: viewFreshness.Unknown++
		}
	}
	if partialSLA {
		result.Diagnostics = append(result.Diagnostics, "sla-policy-unavailable")
	}
	if partialSourceFreshness {
		result.Diagnostics = append(result.Diagnostics, "work-source-degraded")
	}
	// A later page must be reconstructed from the newly authorized full
	// bounded projection. A stale/revoked cursor fails without partial rows.
	authorized := make([]Row, 0, len(selected))
	for _, item := range selected { authorized = append(authorized, item.row) }
	page, next, previousCursor, hasPrevious, err := authorizedQueuePage(req, authorized)
	if err != nil { return empty, err }
	result.Rows = page
	if hasOwnerMetadata {
		result.ViewFreshness = &viewFreshness
	}
	// req.Cursor is merely a hint until authorizedQueuePage returns nil.
	// Only after this point may the opaque token be used in native HTML
	// navigation controls; stale/revoked/foreign cursors returned ErrCursor.
	result.CurrentCursor = req.Cursor
	result.NextCursor = next
	result.PreviousCursor = previousCursor
	result.HasPrevious = hasPrevious
	if next != "" {
		if hasPrevious {
			result.Diagnostics = append(result.Diagnostics, "more-authorized-pages")
		} else {
			// Preserve the established first-page diagnostic for existing
			// clients; the opt-in cursor now makes later pages available.
			result.Diagnostics = append(result.Diagnostics, "first-page-only")
		}
	}
	return result, nil
}

func invalidScope(s string) bool {
	return s == "" || s == "*" || strings.TrimSpace(s) == "" ||
		strings.IndexFunc(s, func(r rune) bool { return r < 0x20 || r == 0x7f }) >= 0
}
