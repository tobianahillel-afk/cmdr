package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Each line is a standalone JSON object for local DFIR processing.
// Neither this export nor the operator-selected snapshot is authenticated
// production evidence. Export contents include sensitive authorized metadata.
type caseJSONLHeader struct {
	RecordType string `json:"record_type"`
	Format string `json:"format"`
	CaseID string `json:"case_id"`
	Synthetic bool `json:"synthetic"`
	VisibleEntries int `json:"visible_entries"`
	Visibility string `json:"visibility"`
}

type caseJSONLEntry struct {
	RecordType string `json:"record_type"`
	CaseID string `json:"case_id"`
	Entry timeline.ProjectedEntry `json:"entry"`
}

// Full independent Case, Timeline Entry, and source authorization is
// performed through the existing first-party HTTP handler before ANY
// sensitive JSONL bytes are written to the caller-controlled output.
func collectAuthorizedCaseJSONL(ctx context.Context, handler http.Handler, token, caseID string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rows := make([]timeline.ProjectedEntry, 0)
	receipt, err := collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		if len(rows) >= timeline.MaxEntries {
			return errors.New("JSONL event export exceeds permitted capacity")
		}
		rows = append(rows, entry)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(rows) != receipt.VisibleEntries {
		return nil, errors.New("incomplete authorized JSONL export")
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(caseJSONLHeader{
		RecordType: "metadata", Format: "cmdr-local-case-timeline-jsonl-v1",
		CaseID: caseID, Synthetic: true, VisibleEntries: len(rows),
		Visibility: receipt.Visibility,
	}); err != nil {
		return nil, err
	}
	for _, entry := range rows {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := encoder.Encode(caseJSONLEntry{
			RecordType: "timeline-entry", CaseID: caseID, Entry: entry,
		}); err != nil {
			return nil, err
		}
		if buffer.Len() > 2*maxDemoSnapshotBytes {
			return nil, errors.New("JSONL export exceeds bounded receipt size")
		}
	}
	return buffer.Bytes(), nil
}

func writeDemoCaseJSONL(ctx context.Context, principal, snapshot, caseID string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || caseID == "" || len(caseID) > 1024 {
		return errors.New("explicit principal and bounded Case required for JSONL export")
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
	payload, err := collectAuthorizedCaseJSONL(ctx, handler, token, caseID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err = output.Write(payload)
	return err
}
