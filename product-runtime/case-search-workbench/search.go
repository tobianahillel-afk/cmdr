package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// This is an offline integration composition, not an IAM or storage provider.
// A trusted operator chooses a synthetic principal in a validated local
// snapshot. Every returned entry comes through the owner-authorized Timeline
// HTTP handler, not from raw snapshot arrays or caller-supplied access flags.
type searchOptions struct {
	Principal, CaseID, Query, From, To string
}

type caseSearchReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	CaseID string `json:"case_id"`
	QueryScope string `json:"query_scope"`
	Job search.SearchJob `json:"job"`
}

type ownerPage struct {
	Status string `json:"status"`
	Entries []timeline.ProjectedEntry `json:"entries"`
	Page struct {
		Offset int `json:"offset"`
		Limit int `json:"limit"`
		TotalVisible int `json:"total_visible"`
		HasMore bool `json:"has_more"`
	} `json:"page"`
}

// Execute only reuses references previously released by the server-side
// Case/Entry/typed-source owner. Authorization is re-evaluated on execution;
// no client-side claim about the source is sufficient.
type caseSearchBackend struct {
	owner http.Handler
	caseID, tenant, environment, query string
	sourceKinds map[string]string
}

func (b *caseSearchBackend) Execute(ctx context.Context, req search.BackendRequest) (search.BackendResult, error) {
	if ctx.Err() != nil || req.TenantRef != b.tenant || req.EnvironmentRef != b.environment ||
		req.TimeStart.IsZero() || !req.TimeStart.Before(req.TimeEnd) {
		return search.BackendResult{}, errors.New("search scope unavailable")
	}
	rows, err := readOwnerPages(ctx, b.owner, b.caseID)
	if err != nil {
		return search.BackendResult{}, errors.New("authorized search unavailable")
	}
	visible := make(map[string]bool)
	for _, row := range rows {
		if b.sourceKinds[row.SourceRef] == "telemetry-event" {
			visible[row.SourceRef] = true
		}
	}
	allowed := make(map[string]bool, len(req.Sources))
	for _, src := range req.Sources {
		// Revoked or disappeared source authorization must fail the job,
		// not silently masquerade as an empty successful search.
		if b.sourceKinds[src] != "telemetry-event" || allowed[src] || !visible[src] {
			return search.BackendResult{}, errors.New("authorized search source no longer visible")
		}
		allowed[src] = true
	}
	result := search.BackendResult{
		TenantRef: req.TenantRef,
		CompletedSources: append([]string(nil), req.Sources...),
	}
	for _, entry := range rows {
		if err := ctx.Err(); err != nil {
			return search.BackendResult{}, err
		}
		if !allowed[entry.SourceRef] || entry.TimeUncertain {
			continue // unknown clocks cannot enter a precise time-bounded query
		}
		at, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil {
			return search.BackendResult{}, errors.New("invalid authorized event time")
		}
		if at.Before(req.TimeStart) || !at.Before(req.TimeEnd) {
			continue
		}
		// Search is limited to owner-released reference/provenance fields.
		// No raw preview, denied source, or canonical snapshot is searched.
		if !strings.Contains(strings.ToLower(entry.ID), b.query) &&
			!strings.Contains(strings.ToLower(entry.SourceRef), b.query) &&
			!strings.Contains(strings.ToLower(entry.Provenance), b.query) {
			continue
		}
		result.ResultRefs = append(result.ResultRefs, search.ResultReference{Source: entry.SourceRef, Ref: entry.ID})
	}
	return result, nil
}

// readOwnerPages rejects denied, inconsistent, changing or duplicated
// paginated results before any Search Job result reference is released.
func readOwnerPages(ctx context.Context, handler http.Handler, caseID string) ([]timeline.ProjectedEntry, error) {
	const size = 200
	seen := make(map[string]bool)
	rows := make([]timeline.ProjectedEntry, 0)
	total := -1
	for offset := 0; ; offset += size {
		if ctx.Err() != nil || offset > timeline.MaxEntries {
			return nil, errors.New("Case owner traversal unavailable")
		}
		query := url.Values{"case_id": {caseID}, "limit": {"200"}, "offset": {fmt.Sprint(offset)}}
		req := httptest.NewRequest(http.MethodGet, "/case-timeline?"+query.Encode(), nil).WithContext(ctx)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != http.StatusOK || response.Body.Len() > 8<<20 {
			return nil, errors.New("Case owner authorization unavailable")
		}
		var page ownerPage
		if err := json.Unmarshal(response.Body.Bytes(), &page); err != nil {
			return nil, errors.New("invalid Case owner page")
		}
		if total == -1 {
			total = page.Page.TotalVisible
		}
		if total < 0 || total > timeline.MaxEntries || page.Page.TotalVisible != total ||
			page.Page.Limit != size || page.Page.Offset != offset ||
			len(page.Entries) > size || offset+len(page.Entries) > total ||
			(len(page.Entries) == 0 && offset < total) ||
			page.Page.HasMore != (offset+len(page.Entries) < total) ||
			(page.Status != "allow" && page.Status != "filtered" && page.Status != "partial" &&
				page.Status != "empty" && page.Status != "uncertain") {
			return nil, errors.New("inconsistent authorized pagination")
		}
		for _, row := range page.Entries {
			if row.ID == "" || row.SourceRef == "" || row.Provenance == "" ||
				seen[row.ID] || row.Version < 1 {
				return nil, errors.New("invalid authorized event reference")
			}
			seen[row.ID] = true
			rows = append(rows, row)
		}
		if !page.Page.HasMore {
			break
		}
	}
	if len(rows) != total || ctx.Err() != nil {
		return nil, errors.New("incomplete authorized Case search")
	}
	return rows, nil
}

// RunCaseSearch composes existing Event Search job lifecycle with the
// existing authenticated-in-the-fixture Case/Timeline/Source HTTP owner.
// The only granted Search Job permissions are LOCAL SYNTHETIC exercise
// permissions; never reuse this method as a production IAM adapter.
func runCaseSearch(ctx context.Context, data timeline.MemoryTimelineData, opts searchOptions) (caseSearchReceipt, error) {
	if ctx == nil || ctx.Err() != nil || opts.Principal == "" || opts.CaseID == "" ||
		opts.Query == "" || len(opts.Query) > 96 || !utf8.ValidString(opts.Query) ||
		strings.TrimSpace(opts.Query) != opts.Query ||
		strings.IndexFunc(opts.Query, unicode.IsControl) >= 0 {
		return caseSearchReceipt{}, errors.New("invalid local Case search selection")
	}
	var principal timeline.MemoryTimelinePrincipal
	for _, p := range data.Principals {
		if p.ID == opts.Principal {
			principal = p
			break
		}
	}
	if principal.ID == "" {
		return caseSearchReceipt{}, errors.New("synthetic principal unavailable")
	}
	// Only the operator supplies this in-process principal. There is no
	// network endpoint, caller-supplied Tenant header or reusable credential.
	handlers, err := timeline.NewMemoryCaseTimelineHTTPHandlers(func(*http.Request) (string, error) {
		return opts.Principal, nil
	}, data)
	if err != nil {
		return caseSearchReceipt{}, errors.New("invalid owner dataset")
	}
	rows, err := readOwnerPages(ctx, handlers.Timeline, opts.CaseID)
	if err != nil {
		return caseSearchReceipt{}, err
	}
	// Derive telemetry source kind from the trusted owner dataset AFTER owner
	// projection; denied/cross-tenant sources never participate in this set.
	sourceKinds := make(map[string]string)
	for _, s := range data.Sources {
		sourceKinds[s.Identity.Reference] = s.Identity.Kind
	}
	authorized := make(map[string]bool)
	for _, entry := range rows {
		if sourceKinds[entry.SourceRef] == "telemetry-event" {
			authorized[entry.SourceRef] = true
		}
	}
	if len(authorized) == 0 {
		return caseSearchReceipt{}, errors.New("no authorized telemetry event source in Case")
	}
	sources := make([]string, 0, len(authorized))
	for source := range authorized {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	input := search.Input{
		TenantRef: principal.Tenant, AuthorizedTenants: []string{principal.Tenant},
		EnvironmentRef: principal.Environment, TimeStart: opts.From, TimeEnd: opts.To,
		Sources: sources, AuthorizedSources: append([]string(nil), sources...),
		// A synthetic offline operator is intentionally granted all four
		// demonstration job roles; real deployment MUST inject IAM decisions.
		Permissions: []string{
			search.PermissionSearchExecute, search.PermissionQueryRead,
			search.PermissionSearchManage, search.PermissionTelemetryRead,
		},
		QueryPresent: true, QueryValid: true,
		CorrelationID: "local-case-search", QueryVersion: "local-text-search-v1",
	}
	backend := &caseSearchBackend{
		owner: handlers.Timeline, caseID: opts.CaseID,
		tenant: principal.Tenant, environment: principal.Environment,
		query: strings.ToLower(opts.Query), sourceKinds: sourceKinds,
	}
	orchestrator, err := search.NewOrchestrator(backend)
	if err != nil {
		return caseSearchReceipt{}, err
	}
	audit := search.AuditContext{ActorRef: opts.Principal, Rationale: "offline-read-only-case-search"}
	job, err := orchestrator.Create(input, audit)
	if err != nil {
		return caseSearchReceipt{}, errors.New("search job input rejected")
	}
	job, err = orchestrator.Execute(ctx, job, audit)
	if err != nil {
		return caseSearchReceipt{}, errors.New("authorized search job failed")
	}
	if job.State != search.JobCompleted || ctx.Err() != nil {
		return caseSearchReceipt{}, errors.New("search job did not fully complete")
	}
	return caseSearchReceipt{
		Format: "cmdr-local-authorized-case-search-job-v1",
		Synthetic: true, CaseID: opts.CaseID,
		QueryScope: "authorized-telemetry-timeline-reference-and-provenance-only",
		Job: job,
	}, nil
}
