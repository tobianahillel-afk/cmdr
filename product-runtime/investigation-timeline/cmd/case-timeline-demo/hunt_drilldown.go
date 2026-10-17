package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A deliberately explicit, sensitive LOCAL report for one authorized Case.
// Unlike the aggregate cross-Case hunt, it contains already authorized event
// identifiers and source references for an analyst to inspect individually.
const maxHuntDrilldownEvents = 25

type localHuntCaseEvents struct {
	Format        string                    `json:"format"`
	Synthetic     bool                      `json:"synthetic"`
	Query         string                    `json:"query"`
	CaseID        string                    `json:"case_id"`
	TotalMatches  int                       `json:"total_matches"`
	Returned      int                       `json:"returned"`
	Truncated     bool                      `json:"truncated"`
	Events        []timeline.ProjectedEntry `json:"events"`
}

// collectAuthorizedHuntCaseEvents uses the same paged, authenticated
// Case/Timeline/typed-source HTTP owner path as normal triage. It cannot
// inspect the operator-owned raw snapshot directly or bypass source ACLs.
func collectAuthorizedHuntCaseEvents(
	ctx context.Context, handler http.Handler, token, caseID, query string,
) (localHuntCaseEvents, error) {
	if err := ctx.Err(); err != nil {
		return localHuntCaseEvents{}, err
	}
	if caseID == "" || len(caseID) > 1024 || !validHuntText(query) {
		return localHuntCaseEvents{}, errors.New("explicit authorized Case and printable 2-96 byte query required")
	}
	events := make([]timeline.ProjectedEntry, 0, maxHuntDrilldownEvents)
	summary, err := collectDemoCaseTriageWithObserverAndQuery(
		ctx, handler, token, caseID, query,
		func(entry timeline.ProjectedEntry) error {
			if len(events) < maxHuntDrilldownEvents {
				events = append(events, entry)
			}
			return nil
		},
	)
	if err != nil {
		return localHuntCaseEvents{}, err
	}
	if err := ctx.Err(); err != nil {
		return localHuntCaseEvents{}, err
	}
	return localHuntCaseEvents{
		Format: "cmdr-local-selected-case-hunt-events-v1",
		Synthetic: true, Query: query, CaseID: caseID,
		TotalMatches: summary.VisibleEntries, Returned: len(events),
		Truncated: summary.VisibleEntries > len(events), Events: events,
	}, nil
}

func writeDemoHuntCaseEvents(
	ctx context.Context, principal, snapshot, caseID, query string, output io.Writer,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || principal == "*" || caseID == "" || !validHuntText(query) {
		return errors.New("explicit synthetic principal, authorized Case and text query required")
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
	result, err := collectAuthorizedHuntCaseEvents(ctx, handler, token, caseID, query)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// No partial output if a later page is unavailable or source ACL changes.
	return json.NewEncoder(output).Encode(result)
}
