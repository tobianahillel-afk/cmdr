package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Local source-focused event sequence built exclusively from the existing
// authenticated Case discovery and typed-source ACL filtered Timeline pages.
// The operator explicitly requests the source; this unsigned output is sensitive.
type localSourceTimeline struct {
	Format         string                    `json:"format"`
	Synthetic      bool                      `json:"synthetic"`
	CaseID         string                    `json:"case_id"`
	SourceRef      string                    `json:"source_reference"`
	Visibility     string                    `json:"case_visibility"`
	MatchedEntries int                       `json:"matched_entries"`
	Observed       int                       `json:"observed"`
	Inferred       int                       `json:"inferred"`
	Entries        []timeline.ProjectedEntry `json:"entries"`
}

// No listener, new provider or permission shortcut. Full authorized paging
// must complete before any byte is written to the local output.
func writeDemoSourceTimeline(ctx context.Context, principal, snapshot, caseID, reference string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || caseID == "" || reference == "" ||
		len(caseID) > 1024 || len(reference) > 1024 {
		return errors.New("explicit principal, Case and bounded source reference required")
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
	report := localSourceTimeline{
		Format: "cmdr-local-source-timeline-v1", Synthetic: true,
		CaseID: caseID, SourceRef: reference,
		Entries: make([]timeline.ProjectedEntry, 0),
	}
	caseSummary, err := collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		if entry.SourceRef != reference {
			return nil
		}
		report.Entries = append(report.Entries, entry)
		if entry.Kind == "observed" {
			report.Observed++
		}
		if entry.Kind == "inferred" {
			report.Inferred++
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(report.Entries) == 0 {
		return errors.New("source not visible in authorized Case")
	}
	report.MatchedEntries = len(report.Entries)
	report.Visibility = caseSummary.Visibility
	return json.NewEncoder(output).Encode(report)
}
