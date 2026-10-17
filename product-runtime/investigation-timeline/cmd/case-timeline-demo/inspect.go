package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// localEntryInspection is an operator-requested, sensitive synthetic preview,
// never a custody-certified evidence export or production authorization result.
type localEntryInspection struct {
	Format    string `json:"format"`
	Synthetic bool   `json:"synthetic"`
	CaseID    string `json:"case_id"`
	EntryID   string `json:"entry_id"`
	Kind      string `json:"source_kind"`
	ObjectID  string `json:"object_id"`
	Preview   string `json:"preview"`
}

// This path intentionally goes through the existing authenticated HTTP owner
// handler rather than reading a source directly from the operator JSON file.
func writeDemoEntryInspection(ctx context.Context, principal, snapshot, caseID, entryID string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || caseID == "" || entryID == "" ||
		len(caseID) > 1024 || len(entryID) > 1024 {
		return errors.New("explicit principal, Case and bounded Timeline Entry required")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	receipt, err := inspectDemoEntryWithData(ctx, principal, data, caseID, entryID)
	if err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(receipt)
}

// Reuses the existing typed source owner for both ordinary Case inspection
// and selected imported Wazuh Event Inspection. This is the only source-open
// authorization boundary in the local CLI.
func inspectDemoEntryWithData(ctx context.Context, principal string, data timeline.MemoryTimelineData, caseID, entryID string) (localEntryInspection, error) {
	token, err := randomBearer()
	if err != nil {
		return localEntryInspection{}, err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return localEntryInspection{}, err
	}
	q := url.Values{"case_id": {caseID}, "entry_id": {entryID}}
	var source struct {
		Status    string `json:"status"`
		Synthetic bool   `json:"synthetic"`
		CaseID    string `json:"case_id"`
		EntryID   string `json:"entry_id"`
		Kind      string `json:"kind"`
		ObjectID  string `json:"object_id"`
		Preview   string `json:"preview"`
	}
	if err := queryDemoOwner(ctx, handler, token, "/demo-source-open?"+q.Encode(), &source); err != nil {
		return localEntryInspection{}, err
	}
	if err := ctx.Err(); err != nil {
		return localEntryInspection{}, err
	}
	if source.Status != "ready" || !source.Synthetic ||
		source.CaseID != caseID || source.EntryID != entryID ||
		source.Kind == "" || source.ObjectID == "" {
		return localEntryInspection{}, errors.New("source was not authorized or returned an inconsistent identity")
	}
	receipt := localEntryInspection{
		Format: "cmdr-local-entry-inspection-v1", Synthetic: true,
		CaseID: caseID, EntryID: entryID,
		Kind: source.Kind, ObjectID: source.ObjectID, Preview: source.Preview,
	}
	return receipt, nil
}
