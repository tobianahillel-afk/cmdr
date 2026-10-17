package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

const wazuhLocalSearchSource = "wazuh4-authorized-local-import"
const maxWazuhSearchMatches = 200

// This backend executes one pre-authorized Event Inspection PivotDraft over
// imported Wazuh metadata. The effective authority remains the existing
// Case/Timeline/typed-source owner, not agent/rule values or caller flags.
type wazuhPivotBackend struct {
	owner http.Handler
	token, caseID, tenant, environment string
	levels map[string]wazuhRuleLevel
	imported int
	field, value string
}

func (b *wazuhPivotBackend) Execute(ctx context.Context, req search.BackendRequest) (search.BackendResult, error) {
	if err := ctx.Err(); err != nil {
		return search.BackendResult{}, err
	}
	if req.TenantRef != b.tenant || req.EnvironmentRef != b.environment ||
		!req.TimeStart.Before(req.TimeEnd) ||
		len(req.Sources) != 1 || req.Sources[0] != wazuhLocalSearchSource {
		return search.BackendResult{}, errors.New("Wazuh search tenant, time or source scope changed")
	}

	result := search.BackendResult{
		TenantRef: req.TenantRef,
		CompletedSources: []string{wazuhLocalSearchSource},
	}
	visibleImported := 0
	_, err := collectDemoCaseTriageWithObserver(ctx, b.owner, b.token, b.caseID, func(entry timeline.ProjectedEntry) error {
		meta, imported := b.levels[entry.ID]
		if !imported {
			return nil
		}
		visibleImported++
		// Unknown times cannot satisfy a precise half-open search window.
		if entry.TimeUncertain {
			return nil
		}
		at, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil {
			return errors.New("invalid owner-authorized Wazuh timestamp")
		}
		if at.Before(req.TimeStart) || !at.Before(req.TimeEnd) {
			return nil
		}
		var actual string
		switch b.field {
		case "wazuh.rule.id":
			actual = meta.RuleID
		case "wazuh.agent.id":
			actual = meta.AgentID
		default:
			return errors.New("unsupported normalized Wazuh pivot")
		}
		if actual == "" || actual != b.value {
			return nil
		}
		if len(result.ResultRefs) >= maxWazuhSearchMatches {
			// Never return a silently truncated "completed" search.
			return errors.New("Wazuh search exceeds bounded result capacity")
		}
		// Second, independent typed-source opening makes matched references
		// actionable; raw source previews are never copied into SearchJob.
		q := url.Values{"case_id": {b.caseID}, "entry_id": {entry.ID}}
		var opened struct {
			Status string `json:"status"`
			Synthetic bool `json:"synthetic"`
			CaseID string `json:"case_id"`
			EntryID string `json:"entry_id"`
			Kind string `json:"kind"`
			ObjectID string `json:"object_id"`
		}
		if err := queryDemoOwner(ctx, b.owner, b.token, "/demo-source-open?"+q.Encode(), &opened); err != nil {
			return errors.New("Wazuh telemetry source authorization revoked")
		}
		if opened.Status != "ready" || !opened.Synthetic || opened.CaseID != b.caseID ||
			opened.EntryID != entry.ID || opened.Kind != "telemetry-event" || opened.ObjectID == "" {
			return errors.New("Wazuh typed source identity changed")
		}
		result.ResultRefs = append(result.ResultRefs, search.ResultReference{
			Source: wazuhLocalSearchSource, Ref: entry.ID,
		})
		return nil
	})
	if err != nil {
		return search.BackendResult{}, err
	}
	if visibleImported != b.imported || visibleImported != len(b.levels) || ctx.Err() != nil {
		return search.BackendResult{}, errors.New("incomplete owner-authorized Wazuh search")
	}
	return result, nil
}

// executeDemoWazuhPivot preserves the original first-party Job lifecycle,
// validated tenant/time/source/permission envelope, stable result references
// and audit trail. The synthetic local operator role is NOT production IAM.
func executeDemoWazuhPivot(ctx context.Context, principal, caseID string, owner http.Handler,
	token, tenant, environment string, levels map[string]wazuhRuleLevel, imported int,
	pivot *inspection.PivotDraft) (search.SearchJob, error) {
	if pivot == nil || pivot.TargetCapability != inspection.PivotTargetCapability ||
		pivot.TenantRef != tenant || pivot.EnvironmentRef != environment ||
		(pivot.Field != "wazuh.rule.id" && pivot.Field != "wazuh.agent.id") ||
		pivot.Value == "" || pivot.TimeRange.Start == "" || pivot.TimeRange.End == "" {
		return search.SearchJob{}, errors.New("invalid authorized Event Inspection pivot")
	}
	return executeDemoWazuhMetadataSearch(ctx, principal, caseID, owner, token, tenant, environment,
		levels, imported, pivot.Field, pivot.Value, pivot.TimeRange.Start, pivot.TimeRange.End)
}

// Both the existing CLI Event Inspection pivot and the authenticated local
// analyst search route use one first-party Event Search Job execution path.
// The route cannot grant itself more access than the per-Case owner.
func executeDemoWazuhMetadataSearch(ctx context.Context, principal, caseID string, owner http.Handler,
	token, tenant, environment string, levels map[string]wazuhRuleLevel, imported int,
	field, value, start, end string) (search.SearchJob, error) {
	if ctx == nil || (field != "wazuh.rule.id" && field != "wazuh.agent.id") ||
		!validWazuhCohortID(value) {
		return search.SearchJob{}, errors.New("invalid bounded Wazuh search selection")
	}
	backend := &wazuhPivotBackend{
		owner: owner, token: token, caseID: caseID,
		tenant: tenant, environment: environment, levels: levels,
		imported: imported, field: field, value: value,
	}
	orchestrator, err := search.NewOrchestrator(backend)
	if err != nil {
		return search.SearchJob{}, err
	}
	input := search.Input{
		TenantRef: tenant, AuthorizedTenants: []string{tenant},
		EnvironmentRef: environment,
		TimeStart: start, TimeEnd: end,
		Sources: []string{wazuhLocalSearchSource},
		AuthorizedSources: []string{wazuhLocalSearchSource},
		Permissions: []string{
			search.PermissionSearchExecute, search.PermissionQueryRead,
			search.PermissionSearchManage, search.PermissionTelemetryRead,
		},
		QueryPresent: true, QueryValid: true,
		CorrelationID: "offline-wazuh-event-inspection",
		QueryVersion: "wazuh-agent-rule-pivot-v1",
	}
	audit := search.AuditContext{ActorRef: principal, Rationale: "offline-authorized-wazuh-pivot"}
	job, err := orchestrator.Create(input, audit)
	if err != nil {
		return search.SearchJob{}, err
	}
	job, err = orchestrator.Execute(ctx, job, audit)
	if err != nil || job.State != search.JobCompleted || ctx.Err() != nil {
		return search.SearchJob{}, errors.New("authorized Wazuh search job did not complete")
	}
	return job, nil
}
