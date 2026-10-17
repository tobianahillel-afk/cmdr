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
	merged, err := importLocalEvents(ctx, snapshot, principal, caseID, events)
	if err != nil {
		return err
	}
	// The imported owner graph is already fully validated. Each accepted
	// normalized event adds exactly one Timeline Entry and never modifies an
	// existing one; avoid a second Case traversal or permission decision.
	imported := len(merged.Entries) - len(snapshot.Entries)
	if imported < 1 || imported > timeline.MaxEntries {
		return errors.New("inconsistent imported Case entry count")
	}
	receipt, err := collectDemoAuthorizedCaseTriage(ctx, principal, merged, caseID)
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
