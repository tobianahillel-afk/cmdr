// Package workqueue implements a bounded, read-only Command-owned Work Queue
// projection. It does not create Incidents, Tasks, Saved Views or permissions.
package workqueue

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
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
	Kind        Kind
	ID          string
	Tenant      string
	Environment string
	Owner       string
	Followers   []string
	Status      string
	UpdatedAt   string
	SLAAtRisk   *bool // nil: unavailable from the authoritative policy source
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
}

type Row struct {
	Kind      Kind   `json:"kind"`
	ID        string `json:"id"`
	Owner     string `json:"owner,omitempty"`
	Status    string `json:"status,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

type Result struct {
	View        View     `json:"view"`
	Tenant      string   `json:"tenant"`
	Environment string   `json:"environment"`
	State       string   `json:"state"`
	Diagnostics []string `json:"diagnostics"`
	Rows        []Row    `json:"rows"`
}

var (
	ErrScope      = errors.New("work queue requires explicit principal, tenant and environment")
	ErrView       = errors.New("unknown canonical Command system view")
	ErrIncomplete = errors.New("untrusted or incomplete work item owner data")
	ErrOwner      = errors.New("work queue owner authorization unavailable")
	ErrTooMany    = errors.New("work queue owner snapshot exceeds local bound")
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
		row     Row
		updated time.Time
	}
	var selected []candidate
	seen := map[string]bool{}
	partialSLA := false
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
				partialSLA = true
			} else {
				accept = *item.SLAAtRisk
			}
		}
		if accept {
			selected = append(selected, candidate{
				row:     Row{Kind: item.Kind, ID: item.ID, Owner: item.Owner, Status: item.Status, UpdatedAt: item.UpdatedAt},
				updated: t,
			})
		}
	}
	sort.Slice(selected, func(i, j int) bool {
		if !selected[i].updated.Equal(selected[j].updated) {
			return selected[i].updated.After(selected[j].updated)
		}
		if selected[i].row.Kind != selected[j].row.Kind {
			return selected[i].row.Kind < selected[j].row.Kind
		}
		return selected[i].row.ID < selected[j].row.ID
	})
	result := Result{View: req.View, Tenant: req.Tenant, Environment: req.Environment,
		State: "ready", Rows: make([]Row, 0, len(selected)), Diagnostics: []string{}}
	if len(selected) == 0 {
		result.State = "no-authorized-items"
	}
	if partialSLA {
		result.State = "partial-sources"
		result.Diagnostics = []string{"sla-policy-unavailable"}
	}
	// Initial bounded slice. Stable cursor/page retrieval and custom Shared
	// Saved View evaluation are separately owned future integration work.
	if len(selected) > 50 {
		selected = selected[:50]
		result.Diagnostics = append(result.Diagnostics, "first-page-only")
	}
	for _, item := range selected {
		result.Rows = append(result.Rows, item.row)
	}
	return result, nil
}

func invalidScope(s string) bool {
	return s == "" || s == "*" || strings.TrimSpace(s) == "" ||
		strings.IndexFunc(s, func(r rune) bool { return r < 0x20 || r == 0x7f }) >= 0
}
