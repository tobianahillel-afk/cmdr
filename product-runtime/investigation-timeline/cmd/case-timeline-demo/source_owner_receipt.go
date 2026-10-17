package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"unicode/utf8"
)

// A manually requested, sensitive local receipt, never an authority token,
// tamper-evident evidence package or production source-object export.
type localSourceOwnerReceipt struct {
	Format    string `json:"format"`
	Synthetic bool   `json:"synthetic"`
	CaseID    string `json:"case_id"`
	EntryID   string `json:"entry_id"`
	Kind      string `json:"kind"`
	ObjectID  string `json:"object_id"`
	Preview   string `json:"preview"`
}

// Reuse the existing trusted owner HTTP route, which reauthorizes the
// Case, Timeline Entry and typed source at open time. Never accept an
// object ID or owner permission from an operator flag or prior JSON receipt.
func writeDemoSourceOwnerReceipt(
	ctx context.Context, principal, snapshot, caseID, entryID string, output io.Writer,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || caseID == "" || entryID == "" ||
		len(caseID) > 1024 || len(entryID) > 1024 ||
		!utf8.ValidString(caseID) || !utf8.ValidString(entryID) {
		return errors.New("explicit bounded Case, entry and analyst required")
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
	query := url.Values{"case_id": {caseID}, "entry_id": {entryID}}
	var opened struct {
		Status    string `json:"status"`
		Synthetic bool   `json:"synthetic"`
		CaseID    string `json:"case_id"`
		EntryID   string `json:"entry_id"`
		Kind      string `json:"kind"`
		ObjectID  string `json:"object_id"`
		Preview   string `json:"preview"`
	}
	if err := queryDemoOwner(ctx, handler, token, "/demo-source-open?"+query.Encode(), &opened); err != nil {
		return errors.New("source owner unavailable for selected Case and entry")
	}
	if opened.Status != "ready" || !opened.Synthetic ||
		opened.CaseID != caseID || opened.EntryID != entryID ||
		opened.ObjectID == "" || len(opened.ObjectID) > 1024 ||
		len(opened.Preview) > 2048 || !utf8.ValidString(opened.Preview) {
		return errors.New("invalid authorized source owner response")
	}
	switch opened.Kind {
	case "evidence", "telemetry-event", "incident", "decision":
	default:
		return errors.New("invalid source owner kind")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localSourceOwnerReceipt{
		Format: "cmdr-local-authorized-source-preview-v1",
		Synthetic: true, CaseID: caseID, EntryID: entryID,
		Kind: opened.Kind, ObjectID: opened.ObjectID, Preview: opened.Preview,
	})
}
