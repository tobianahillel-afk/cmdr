package casequeue

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
)

const (
	// These are local routing markers, not caller-supplied URLs or permissions.
	OriginCaseQueue       = "case-queue"
	OriginCommandIncident = "command-incident"
	CaseWorkspaceScreen   = "INV-CAS-001"
	returnTokenVersion    = 1
	maxReturnTokenBytes   = 8192
	maxReturnScroll       = 10000000

	ErrorInvalidHandoffState ErrorCode = "invalid-handoff-state"
	ErrorCaseNotVisible      ErrorCode = "handoff-case-not-visible"
	ErrorCaseStale           ErrorCode = "handoff-case-stale"
	ErrorReturnViewChanged   ErrorCode = "handoff-view-changed"
)

// ReturnState holds only bounded, caller-visible navigation state. It is not an
// authorization token, a URL, a persisted Saved View, or a Case snapshot.
type ReturnState struct {
	TenantRef        string     `json:"tenant_ref"`
	EnvironmentRef   string     `json:"environment_ref"`
	Origin           string     `json:"origin"`
	SelectedCaseID   string     `json:"selected_case_id"`
	ScrollOffset     int        `json:"scroll_offset"`
	SearchQuery      string     `json:"search_query"`
	Filters          []Filter   `json:"filters"`
	Sort             []SortSpec `json:"sort"`
	SavedViewID      string     `json:"saved_view_id"`
	SavedViewVersion string     `json:"saved_view_version"`
	ViewDirty        bool       `json:"view_dirty"`
	// Optional transient page location, never a grant or persistent snapshot.
	PageLimit        int        `json:"page_limit,omitempty"`
	PageCursor       string     `json:"page_cursor,omitempty"`
}

// WorkspaceHandoff is a read-only transition. Neither it nor its return token
// contains an Incident, Finding, mutable Case, or privileged Saved View field.
type WorkspaceHandoff struct {
	ScreenID       string `json:"screen_id"`
	CaseID         string `json:"case_id"`
	TenantRef      string `json:"tenant_ref"`
	EnvironmentRef string `json:"environment_ref"`
	CaseUpdatedAt  string `json:"case_updated_at"`
	ReturnToken    string `json:"return_token"`
}

type HandoffResult struct {
	Allowed       bool              `json:"allowed"`
	ErrorCode     ErrorCode         `json:"error_code"`
	AuditRequired bool              `json:"audit_required"`
	Handoff       *WorkspaceHandoff `json:"handoff,omitempty"`
}

type returnEnvelope struct {
	Version int `json:"version"`
	State   ReturnState `json:"state"`
}

func validOrigin(origin string) bool {
	return origin == OriginCaseQueue || origin == OriginCommandIncident
}

func safeSearch(query string) bool {
	if len(query) > 256 {
		return false
	}
	for _, r := range query {
		if r < 0x20 {
			return false
		}
	}
	return true
}

func visibleRow(projection *Projection, caseID string) (*Row, ErrorCode) {
	if projection == nil || caseID == "" {
		return nil, ErrorCaseNotVisible
	}
	for i := range projection.Rows {
		row := &projection.Rows[i]
		if row.Case.ID != caseID {
			continue
		}
		if row.Case.TenantRef != projection.TenantRef {
			return nil, ErrorCaseNotVisible
		}
		if row.Stale || row.Case.Freshness == "stale" {
			return nil, ErrorCaseStale
		}
		return row, ErrorNone
	}
	return nil, ErrorCaseNotVisible
}

func viewReference(tenant string, view *SavedView) (string, string, ErrorCode) {
	applied := ApplySavedView(tenant, view)
	if !applied.Allowed {
		return "", "", applied.ErrorCode
	}
	if applied.Applied && applied.Sanitized != nil {
		return applied.Sanitized.ID, applied.Sanitized.Version, ErrorNone
	}
	return "", "", ErrorNone
}

// OpenCaseWorkspace validates the full current projection before emitting any
// transition. Selection is checked AFTER access control, filters and Saved View.
func OpenCaseWorkspace(in Input, selectedCaseID string, scroll int, origin string) HandoffResult {
	return openCaseWorkspaceWithPage(in, selectedCaseID, scroll, origin, 0, "")
}

// Reuse the same authorized handoff engine for a Case selected from one
// visible Queue page. A paging cursor is presentation state, never authority.
func openCaseWorkspaceWithPage(in Input, selectedCaseID string, scroll int,
	origin string, limit int, cursor string) HandoffResult {
	rejected := func(code ErrorCode) HandoffResult {
		return HandoffResult{ErrorCode: code, AuditRequired: true}
	}
	if !validOrigin(origin) || scroll < 0 || scroll > maxReturnScroll ||
		!safeSearch(in.SearchQuery) || strings.TrimSpace(selectedCaseID) == "" {
		return rejected(ErrorInvalidHandoffState)
	}
	projection := ProjectWithSavedViewFallback(in)
	if !projection.Allowed || projection.Projection == nil {
		return rejected(projection.ErrorCode)
	}
	if limit == 0 && cursor != "" {
		return rejected(ErrorInvalidHandoffState)
	}
	visibleProjection := projection.Projection
	if limit != 0 {
		page, err := paginateQueueProjection(projection.Projection, limit, cursor, in.SavedView)
		if err != nil {
			return rejected(ErrorInvalidHandoffState)
		}
		visibleProjection = page.Projection
	}
	row, code := visibleRow(visibleProjection, selectedCaseID)
	if code != ErrorNone {
		return rejected(code)
	}
	id, version, code := viewReference(in.TenantRef, in.SavedView)
	if code != ErrorNone {
		return rejected(code)
	}
	filters, code := validateFilters(in.Filters)
	if code != ErrorNone {
		return rejected(code)
	}
	sortSpecs, code := validateSort(in.Sort)
	if code != ErrorNone {
		return rejected(code)
	}
	state := ReturnState{
		TenantRef: in.TenantRef, EnvironmentRef: in.EnvironmentRef, Origin: origin,
		SelectedCaseID: row.Case.ID, ScrollOffset: scroll, SearchQuery: in.SearchQuery,
		Filters: filters, Sort: sortSpecs, SavedViewID: id, SavedViewVersion: version,
		ViewDirty: in.ViewDirty, PageLimit: limit, PageCursor: cursor,
	}
	payload, err := json.Marshal(returnEnvelope{Version: returnTokenVersion, State: state})
	if err != nil || len(payload) > maxReturnTokenBytes {
		return rejected(ErrorInvalidHandoffState)
	}
	return HandoffResult{Allowed: true, AuditRequired: true, Handoff: &WorkspaceHandoff{
		ScreenID: CaseWorkspaceScreen, CaseID: row.Case.ID, TenantRef: in.TenantRef,
		EnvironmentRef: in.EnvironmentRef, CaseUpdatedAt: row.Case.UpdatedAt,
		ReturnToken: base64.RawURLEncoding.EncodeToString(payload),
	}}
}

// RestoreCaseQueue rejects stale/denied selections and re-evaluates the current
// caller's tenant and Saved View. Base64 is serialization, NOT authorization.
func RestoreCaseQueue(token string, current Input) (ReturnState, ErrorCode) {
	fail := func(code ErrorCode) (ReturnState, ErrorCode) { return ReturnState{}, code }
	if token == "" || len(token) > maxReturnTokenBytes*2 {
		return fail(ErrorInvalidHandoffState)
	}
	payload, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(payload) > maxReturnTokenBytes {
		return fail(ErrorInvalidHandoffState)
	}
	var envelope returnEnvelope
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return fail(ErrorInvalidHandoffState)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fail(ErrorInvalidHandoffState)
	}
	state := envelope.State
	if envelope.Version != returnTokenVersion || !validOrigin(state.Origin) ||
		state.TenantRef == "" || state.TenantRef == "*" ||
		state.TenantRef != current.TenantRef || state.EnvironmentRef == "" ||
		state.EnvironmentRef != current.EnvironmentRef ||
		state.SelectedCaseID == "" || !safeSearch(state.SearchQuery) ||
		state.ScrollOffset < 0 || state.ScrollOffset > maxReturnScroll {
		return fail(ErrorInvalidHandoffState)
	}
	id, version, code := viewReference(current.TenantRef, current.SavedView)
	if code != ErrorNone {
		return fail(code)
	}
	if id != state.SavedViewID || version != state.SavedViewVersion {
		return fail(ErrorReturnViewChanged)
	}
	filters, code := validateFilters(state.Filters)
	if code != ErrorNone {
		return fail(code)
	}
	sortSpecs, code := validateSort(state.Sort)
	if code != ErrorNone {
		return fail(code)
	}
	// Clone before re-evaluating. Current snapshots and Shared state are immutable.
	query := cloneInputForSavedView(current)
	query.SearchQuery, query.Filters, query.Sort = state.SearchQuery, filters, sortSpecs
	query.ViewDirty = state.ViewDirty
	projection := ProjectWithSavedViewFallback(query)
	if !projection.Allowed || projection.Projection == nil {
		return fail(projection.ErrorCode)
	}
	// For an opened page, revalidate both the complete current authorized
	// projection and the specific original page before restoring selection.
	// Revoked sources, edited row contents or Shared View changes cause
	// pagination to fail closed rather than silently jump to another page.
	visibleProjection := projection.Projection
	if state.PageLimit == 0 && state.PageCursor != "" {
		return fail(ErrorInvalidHandoffState)
	}
	if state.PageLimit != 0 {
		page, err := paginateQueueProjection(projection.Projection,
			state.PageLimit, state.PageCursor, current.SavedView)
		if err != nil {
			return fail(ErrorInvalidHandoffState)
		}
		visibleProjection = page.Projection
	}
	if _, code = visibleRow(visibleProjection, state.SelectedCaseID); code != ErrorNone {
		return fail(code)
	}
	state.Filters, state.Sort = filters, sortSpecs
	return state, ErrorNone
}
