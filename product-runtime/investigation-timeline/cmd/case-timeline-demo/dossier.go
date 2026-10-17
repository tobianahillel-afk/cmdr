package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A local, unsigned analyst handoff. It is NOT a chain-of-custody artifact,
// production Case repository or assertion that separate entries are causal.
type localCaseDossier struct {
	Format         string                    `json:"format"`
	Synthetic      bool                      `json:"synthetic"`
	CaseID         string                    `json:"case_id"`
	SelectedEntry  timeline.ProjectedEntry   `json:"selected_entry"`
	Source         localEntryInspection      `json:"source"`
	Summary        caseTriageReceipt         `json:"case_summary"`
	RelatedEntries []timeline.ProjectedEntry `json:"same_source_entries"`
}

// Compose only existing owner-authorized Case/Timeline/source HTTP handlers.
// The selected source is revalidated after complete authorized pagination.
// No bytes are written to the caller before all checks succeed.
func writeDemoCaseDossier(ctx context.Context, principal, snapshot, caseID, entryID string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || caseID == "" || entryID == "" ||
		len(principal) > 1024 || len(caseID) > 1024 || len(entryID) > 1024 {
		return errors.New("explicit bounded principal, Case and Timeline Entry required")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	// Collect all visible rows from the existing paginated owner endpoint.
	// Do not read source owner objects from the operator file directly.
	entries := make([]timeline.ProjectedEntry, 0)
	summary, err := collectDemoCaseTriageWithObserver(ctx, handler, token, caseID,
		func(entry timeline.ProjectedEntry) error {
			entries = append(entries, entry)
			return nil
		})
	if err != nil {
		return err
	}
	var selected timeline.ProjectedEntry
	found := false
	for _, entry := range entries {
		if entry.ID == entryID {
			selected = entry
			found = true
			break
		}
	}
	if !found {
		return errors.New("selected Timeline Entry is not authorized in the Case")
	}
	source, err := inspectAuthorizedDemoEntry(ctx, handler, token, caseID, entryID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	report := localCaseDossier{
		Format: "cmdr-local-case-dossier-v1", Synthetic: true,
		CaseID: caseID, SelectedEntry: selected, Summary: summary,
		Source: source,
		RelatedEntries: make([]timeline.ProjectedEntry, 0),
	}
	for _, entry := range entries {
		if entry.SourceRef == selected.SourceRef {
			report.RelatedEntries = append(report.RelatedEntries, entry)
		}
	}
	return json.NewEncoder(output).Encode(report)
}
