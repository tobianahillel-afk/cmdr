package main

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"io"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Reuse the first-party Event Search Job lifecycle and the existing Case,
// Timeline, and typed-source owner for normalized operator-local JSONL.
// No in-band Wazuh metadata, user-supplied permission or file path grants access.
const localJSONLSearchSource = "normalized-jsonl-authorized-import"

type localJSONLSearchBackend struct {
	owner http.Handler
	bearer, caseID, tenant, environment, literal, field string
	imported map[string]struct{}
	// Source-owned metadata is a query filter only AFTER typed-source ACL.
	wazuh map[string]timeline.WazuhSourceMetadata
}

func (b *localJSONLSearchBackend) Execute(ctx context.Context, request search.BackendRequest) (search.BackendResult, error) {
	if request.TenantRef != b.tenant || request.EnvironmentRef != b.environment ||
		len(request.Sources) != 1 || request.Sources[0] != localJSONLSearchSource ||
		!request.TimeStart.Before(request.TimeEnd) {
		return search.BackendResult{}, errors.New("normalized event search scope changed")
	}
	if b.field != "preview" && len(b.wazuh) != len(b.imported) {
		return search.BackendResult{}, errors.New("persistent Wazuh numeric pivots unavailable for this Case")
	}
	result := search.BackendResult{
		TenantRef: request.TenantRef,
		CompletedSources: []string{localJSONLSearchSource},
	}
	visible, candidates := 0, 0
	_, err := collectDemoCaseTriageWithObserver(ctx, b.owner, b.bearer, b.caseID, func(entry timeline.ProjectedEntry) error {
		if _, imported := b.imported[entry.ID]; !imported {
			return nil // unrelated Case entries are never candidates
		}
		visible++
		if entry.Kind != "observed" || entry.TimeUncertain {
			return errors.New("imported source identity or time uncertain")
		}
		at, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil { return errors.New("invalid owner timestamp") }
		if at.Before(request.TimeStart) || !at.Before(request.TimeEnd) {
			return nil
		}
		candidates++
		if candidates > maxWazuhPreviewCandidates {
			return errors.New("normalized search exceeds 500 candidates")
		}
		// This existing owner call is the security boundary: authorize
		// source before examining its rendered text, even on no-match.
		q := url.Values{"case_id": {b.caseID}, "entry_id": {entry.ID}}
		var opened struct {
			Status string `json:"status"`
			Synthetic bool `json:"synthetic"`
			CaseID string `json:"case_id"`
			EntryID string `json:"entry_id"`
			Kind string `json:"kind"`
			ObjectID string `json:"object_id"`
			Preview string `json:"preview"`
		}
		if err := queryDemoOwner(ctx, b.owner, b.bearer,
			"/demo-source-open?"+q.Encode(), &opened); err != nil {
			return errors.New("normalized source owner denied access")
		}
		if opened.Status != "ready" || !opened.Synthetic ||
			opened.CaseID != b.caseID || opened.EntryID != entry.ID ||
			opened.Kind != "telemetry-event" || opened.ObjectID == "" ||
			len(opened.Preview) > 2048 || !utf8.ValidString(opened.Preview) {
			return errors.New("normalized source owner identity or preview invalid")
		}
		var match bool
		switch b.field {
		case "preview":
			match = strings.Contains(strings.ToLower(opened.Preview), strings.ToLower(b.literal))
		case "rule_id":
			meta, ok := b.wazuh[entry.ID]
			match = ok && meta.RuleID != "" && meta.RuleID == b.literal
		case "agent_id":
			meta, ok := b.wazuh[entry.ID]
			match = ok && meta.AgentID != "" && meta.AgentID == b.literal
		default:
			return errors.New("unsupported persisted Wazuh search field")
		}
		if !match { return nil }
		if len(result.ResultRefs) >= maxWazuhSearchMatches {
			return errors.New("normalized search exceeds 200 matching events")
		}
		result.ResultRefs = append(result.ResultRefs, search.ResultReference{
			Source: localJSONLSearchSource, Ref: entry.ID,
		})
		return nil
	})
	if err != nil { return search.BackendResult{}, err }
	if visible != len(b.imported) || ctx.Err() != nil {
		return search.BackendResult{}, errors.New("normalized search authorization or traversal incomplete")
	}
	return result, nil
}

func executeDemoLocalJSONLSearch(ctx context.Context, principal, caseID string,
	owner http.Handler, bearer, tenant, environment string,
	imported map[string]struct{}, wazuh map[string]timeline.WazuhSourceMetadata,
	field, literal, start, end string) (search.SearchJob, error) {
	if ctx == nil || len(imported) == 0 ||
		!((field == "preview" && validWazuhPreviewLiteral(literal)) ||
			((field == "rule_id" || field == "agent_id") && validWazuhCohortID(literal))) {
		return search.SearchJob{}, errors.New("select bounded normalized text search")
	}
	backend := &localJSONLSearchBackend{
		owner: owner, bearer: bearer, caseID: caseID,
		tenant: tenant, environment: environment,
		imported: imported, literal: literal, field: field, wazuh: wazuh,
	}
	orchestrator, err := search.NewOrchestrator(backend)
	if err != nil { return search.SearchJob{}, err }
	queryVersion := "normalized-jsonl-rendered-literal-v1"
	if field != "preview" { queryVersion = "persisted-wazuh-rule-agent-pivot-v1" }
	job, err := orchestrator.Create(search.Input{
		TenantRef: tenant, AuthorizedTenants: []string{tenant},
		EnvironmentRef: environment,
		TimeStart: start, TimeEnd: end,
		Sources: []string{localJSONLSearchSource},
		AuthorizedSources: []string{localJSONLSearchSource},
		Permissions: []string{
			search.PermissionSearchExecute, search.PermissionQueryRead,
			search.PermissionSearchManage, search.PermissionTelemetryRead,
		},
		QueryPresent: true, QueryValid: true,
		CorrelationID: "offline-normalized-jsonl-hunt",
		QueryVersion: queryVersion,
	}, search.AuditContext{ActorRef: principal, Rationale: "offline-normalized-source-hunt"})
	if err != nil { return search.SearchJob{}, err }
	job, err = orchestrator.Execute(ctx, job, search.AuditContext{
		ActorRef: principal, Rationale: "offline-normalized-source-hunt",
	})
	if err != nil || job.State != search.JobCompleted || ctx.Err() != nil {
		return search.SearchJob{}, errors.New("normalized authorized search did not complete")
	}
	return job, nil
}

// This endpoint is installed exclusively for an operator-approved generic
// JSONL import, never the ordinary snapshot or Wazuh-specific viewer.
func newDemoLocalJSONLSearchHandler(token, principal, caseID string,
	data timeline.MemoryTimelineData, imported map[string]struct{}, owner http.Handler) http.Handler {
	// Build the immutable authorized source metadata index ONCE from the
	// already validated owner graph. Numeric pivots never read raw Wazuh.
	sources := make(map[string]timeline.MemoryTimelineSource, len(data.Sources))
	for _, src := range data.Sources { sources[src.Identity.Reference] = src }
	wazuh := make(map[string]timeline.WazuhSourceMetadata)
	for _, entry := range data.Entries {
		if entry.CaseID != caseID { continue }
		if _, ok := imported[entry.Value.ID]; !ok { continue }
		src, found := sources[entry.Value.SourceRef]
		if found && src.Identity.CaseID == caseID &&
			src.Identity.Kind == "telemetry-event" && src.Wazuh != nil {
			wazuh[entry.Value.ID] = *src.Wazuh
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/demo-local-search" {
			owner.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		// Use POST for sensitive hunt terms: no literal source text, token,
		// or Case selector is placed in a URL or browser history.
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method unavailable", http.StatusMethodNotAllowed)
			return
		}
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")),
			[]byte("Bearer "+token)) != 1 {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		if r.URL.RawQuery != "" || r.Header.Get("Content-Type") != "application/json" {
			http.NotFound(w, r)
			return
		}
		var selection struct {
			CaseID string `json:"case_id"`
			Field string `json:"field,omitempty"`
			Value string `json:"value"`
			From string `json:"from"`
			To string `json:"to"`
		}
		// Enforce a total POST body ceiling, including trailing whitespace
		// after a valid JSON object; do not silently accept a truncated body.
		body, readErr := io.ReadAll(io.LimitReader(r.Body, 4097))
		if readErr != nil || len(body) == 0 || len(body) > 4096 {
			http.NotFound(w, r)
			return
		}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if selection.Field == "" { selection.Field = "preview" }
		if err := decoder.Decode(&selection); err != nil || selection.CaseID != caseID ||
			!((selection.Field == "preview" && validWazuhPreviewLiteral(selection.Value)) ||
				((selection.Field == "rule_id" || selection.Field == "agent_id") &&
					validWazuhCohortID(selection.Value))) {
			http.NotFound(w, r)
			return
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			http.NotFound(w, r)
			return
		}
		from, startErr := time.Parse(time.RFC3339Nano, selection.From)
		to, endErr := time.Parse(time.RFC3339Nano, selection.To)
		if startErr != nil || endErr != nil || !from.Before(to) {
			http.NotFound(w, r)
			return
		}
		var tenant, environment string
		for _, p := range data.Principals {
			if p.ID == principal {
				tenant, environment = p.Tenant, p.Environment
				break
			}
		}
		if tenant == "" || environment == "" {
			http.NotFound(w, r)
			return
		}
		job, err := executeDemoLocalJSONLSearch(r.Context(), principal, caseID,
			owner, token, tenant, environment, imported, wazuh,
			selection.Field, selection.Value, selection.From, selection.To)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		receipt := struct {
			Format string `json:"format"`
			Synthetic bool `json:"synthetic"`
			CaseID string `json:"case_id"`
			Job search.SearchJob `json:"job"`
		}{
			Format: "cmdr-local-jsonl-event-search-v1", Synthetic: true,
			CaseID: caseID, Job: job,
		}
		// Serialize only after all owner authorizations and Job completion;
		// no partial result or preview is written to the HTTP response.
		var payload bytes.Buffer
		if err := json.NewEncoder(&payload).Encode(receipt); err != nil || r.Context().Err() != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(payload.Bytes())
	})
}
