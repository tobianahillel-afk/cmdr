package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// One synthetic operator-selected Wazuh alert is inspected through the real
// typed-source owner and the existing Event Inspection projection engine.
// No Wazuh input, caller flag or normalized metadata grants permissions.
type wazuhInspectionReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	CaseID string `json:"case_id"`
	Inspection *inspection.Projection `json:"inspection"`
	Pivot *inspection.PivotDraft `json:"pivot_draft,omitempty"`
	SearchJob *search.SearchJob `json:"search_job,omitempty"`
}

func writeDemoWazuhEventInspection(ctx context.Context, principal, snapshot, caseID string, paths []string, entryID, pivotField, pivotFrom, pivotTo string, output io.Writer) error {
	return writeDemoWazuhEventInspectionWithSearch(ctx, principal, snapshot, caseID, paths,
		entryID, pivotField, pivotFrom, pivotTo, false, output)
}

func writeDemoWazuhEventInspectionWithSearch(ctx context.Context, principal, snapshot, caseID string, paths []string, entryID, pivotField, pivotFrom, pivotTo string, execute bool, output io.Writer) error {
	if ctx == nil {
		return errors.New("request context required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if entryID == "" || len(paths) == 0 || len(paths) > maxWazuhAlertFiles {
		return errors.New("select a local Wazuh alert in an authorized Case")
	}
	if execute && pivotField == "" {
		return errors.New("executing a Wazuh search requires an Event Inspection pivot")
	}
	if pivotField == "" && (pivotFrom != "" || pivotTo != "") {
		return errors.New("time bounds require an explicitly selected Wazuh pivot field")
	}
	if pivotField != "" && (pivotFrom == "" || pivotTo == "") {
		return errors.New("Wazuh pivot requires an explicit bounded time interval")
	}
	if pivotField != "" && pivotField != "wazuh.rule.id" && pivotField != "wazuh.agent.id" {
		return errors.New("unsupported Wazuh normalized pivot field")
	}

	var data timeline.MemoryTimelineData
	var imported int
	var levels map[string]wazuhRuleLevel
	var err error
	if len(paths) == 1 {
		var raw []byte
		data, imported, raw, err = prepareWazuh4CaseWithRaw(ctx, principal, snapshot, caseID, paths[0])
		if err == nil {
			levels, err = wazuhRuleLevels(raw, caseID)
		}
	} else {
		data, imported, levels, err = loadMultipleWazuhFeedsWithLevels(ctx, principal, snapshot, caseID, paths)
	}
	if err != nil || imported == 0 || len(levels) != imported {
		return errors.New("Wazuh owner import or metadata unavailable")
	}
	metadata, belongsToImport := levels[entryID]
	if !belongsToImport {
		return errors.New("entry does not belong to selected Wazuh import")
	}

	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	var selected timeline.ProjectedEntry
	visibleImported := 0
	_, err = collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		if _, ok := levels[entry.ID]; ok {
			visibleImported++
			if entry.ID == entryID {
				selected = entry
			}
		}
		return nil
	})
	if err != nil || visibleImported != imported || selected.ID != entryID ||
		selected.Kind != "observed" || selected.SourceRef == "" || selected.TimeUncertain {
		return errors.New("selected Wazuh event not fully owner-authorized")
	}
	// Open again through the already existing Case typed-source owner.
	// This is a fresh authorization decision, not a reuse of a preview.
	opened, err := inspectDemoEntryWithData(ctx, principal, data, caseID, entryID)
	if err != nil {
		return err
	}
	if opened.Kind != "telemetry-event" {
		return errors.New("selected source is not telemetry")
	}
	var tenant, environment string
	for _, p := range data.Principals {
		if p.ID == principal {
			tenant, environment = p.Tenant, p.Environment
			break
		}
	}
	// Case/Tenant membership was already enforced by the Case discovery,
	// complete Timeline traversal and fresh typed-source owner checks above.
	// Do not reimplement another Case authorization layer here.
	if tenant == "" || environment == "" {
		return errors.New("principal identity unavailable")
	}

	normalized := make(map[string]string)
	if metadata.RuleID != "" {
		normalized["wazuh.rule.id"] = metadata.RuleID
	}
	if metadata.AgentID != "" {
		normalized["wazuh.agent.id"] = metadata.AgentID
	}
	datasource := "wazuh4-local-operator-feed"
	if metadata.Feed > 0 {
		datasource = fmt.Sprintf("wazuh4-local-feed-%d", metadata.Feed)
	}
	result := inspection.Project(inspection.Input{
		TenantRef: tenant, AuthorizedTenants: []string{tenant}, EnvironmentRef: environment,
		Permissions: []string{inspection.PermissionTelemetryRead}, // synthetic in-process exercise only
		RawAccessDecision: inspection.AccessDeny, RenderedAccessDecision: inspection.AccessAllow,
		CorrelationID: "offline-wazuh-event-inspection",
		Event: inspection.Event{
			Ref: selected.ID, TenantRef: tenant, State: inspection.StateAvailable,
			Source: inspection.SourceMetadata{
				SourceRef: selected.SourceRef, DataSourceRef: datasource, Kind: opened.Kind,
			},
			Parser: inspection.ParserMetadata{
				ParserRef: "wazuh4-local-normalizer", Version: "1", Status: "normalized",
			},
			NormalizedFields: normalized,
			RenderedPayload: []byte(opened.Preview),
		},
	})
	if !result.Allowed || result.Projection == nil {
		return errors.New("Event Inspection projection denied")
	}
	receipt := wazuhInspectionReceipt{
		Format: "cmdr-local-wazuh4-event-inspection-v1", Synthetic: true,
		CaseID: caseID, Inspection: result.Projection,
	}
	if pivotField != "" {
		value, present := normalized[pivotField]
		pivot := inspection.PreparePivot(result.Projection, inspection.PivotSelection{
			Field: pivotField, Value: value, ValuePresent: present,
			ValueAuthorized: present, TimeStart: pivotFrom, TimeEnd: pivotTo,
			ReturnOrigin: "offline-wazuh-case",
		})
		if !pivot.Allowed || pivot.Draft == nil {
			return errors.New("Wazuh search pivot unavailable")
		}
		receipt.Pivot = pivot.Draft
		if execute {
			job, err := executeDemoWazuhPivot(ctx, principal, caseID, handler, token,
				tenant, environment, levels, imported, receipt.Pivot)
			if err != nil {
				return err
			}
			receipt.SearchJob = &job
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// No partial output: serialize only after complete Case/source rechecks
	// and a successful first-party Event Inspection / Pivot evaluation.
	return json.NewEncoder(output).Encode(receipt)
}
