package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// localImportedTriage is a synthetic, unsigned, analyst-requested summary,
// never a serialized owner snapshot or a provider-authoritative finding.
type localImportedTriage struct {
	Format          string            `json:"format"`
	Synthetic       bool              `json:"synthetic"`
	ImportedEntries int               `json:"imported_entries"`
	Case            caseTriageReceipt `json:"case_triage"`
}

// writeDemoImportedJSONLTriage bridges three existing functions:
// trusted operator-local JSONL -> in-memory owner -> authenticated Case
// discovery/projection/typed-source ACL -> bounded triage. It deliberately
// does NOT serialize the merged owner dataset, source previews or a bearer.
func writeDemoImportedJSONLTriage(ctx context.Context, principal, snapshotPath, caseID, eventsPath string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if snapshotPath == "" || eventsPath == "" || snapshotPath == eventsPath {
		return errors.New("select separate existing owner and normalized event files")
	}
	snapshot, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	events, err := readLocalEventFile(eventsPath)
	if err != nil {
		return err
	}
	before := -1
	for _, c := range snapshot.Cases {
		if c.ID == caseID {
			before = len(c.EntryIDs)
			break
		}
	}
	merged, err := importLocalEvents(ctx, snapshot, principal, caseID, events)
	if err != nil {
		return err
	}
	imported := -1
	for _, c := range merged.Cases {
		if c.ID == caseID {
			imported = len(c.EntryIDs) - before
			break
		}
	}
	if imported < 1 || imported > timeline.MaxEntries {
		return errors.New("inconsistent imported Case entry count")
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, merged)
	if err != nil {
		return err
	}
	receipt, err := collectDemoCaseTriage(ctx, handler, token, caseID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localImportedTriage{
		Format: "cmdr-local-imported-case-triage-v1",
		Synthetic: true,
		ImportedEntries: imported,
		Case: receipt,
	})
}
