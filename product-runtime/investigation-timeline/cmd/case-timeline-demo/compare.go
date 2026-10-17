package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// LocalCaseComparison is an unsigned, synthetic, operator-requested receipt.
// Both sides are collected independently through the same authenticated
// Case/Timeline/source owner HTTP handlers as the analyst viewer.
type localCaseComparison struct {
	Format             string            `json:"format"`
	Synthetic          bool              `json:"synthetic"`
	First              caseTriageReceipt `json:"first"`
	Second             caseTriageReceipt `json:"second"`
	VisibleEntryDelta  int               `json:"visible_entry_delta"`
	ObservedDelta      int               `json:"observed_delta"`
	InferredDelta      int               `json:"inferred_delta"`
	CertainTimeOverlap *timeOverlap      `json:"certain_time_overlap_utc,omitempty"`
}

type timeOverlap struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// writeDemoCaseComparison does not start a listener or accept permissions
// from CLI/HTTP. All source information comes from the already authorized
// owner pagination path. Nothing is written until BOTH Cases validate.
func writeDemoCaseComparison(ctx context.Context, principal, snapshot, firstID, secondID string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || firstID == "" || secondID == "" ||
		firstID == secondID || len(firstID) > 1024 || len(secondID) > 1024 {
		return errors.New("select two distinct authorized Case IDs and a principal")
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
	first, err := collectDemoCaseTriage(ctx, handler, token, firstID)
	if err != nil {
		return err
	}
	second, err := collectDemoCaseTriage(ctx, handler, token, secondID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	result := localCaseComparison{
		Format:            "cmdr-local-case-comparison-v1",
		Synthetic:         true,
		First:             first,
		Second:            second,
		VisibleEntryDelta: first.VisibleEntries - second.VisibleEntries,
		ObservedDelta:     first.Kinds["observed"] - second.Kinds["observed"],
		InferredDelta:     first.Kinds["inferred"] - second.Kinds["inferred"],
	}
	// Unknown clock intervals are never promoted to a known overlap.
	if first.FirstCertainUTC != "" && first.LastCertainUTC != "" &&
		second.FirstCertainUTC != "" && second.LastCertainUTC != "" {
		aStart, aErr := time.Parse(time.RFC3339Nano, first.FirstCertainUTC)
		aEnd, bErr := time.Parse(time.RFC3339Nano, first.LastCertainUTC)
		bStart, cErr := time.Parse(time.RFC3339Nano, second.FirstCertainUTC)
		bEnd, dErr := time.Parse(time.RFC3339Nano, second.LastCertainUTC)
		if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
			return errors.New("inconsistent authorized time bounds")
		}
		if bStart.After(aStart) {
			aStart = bStart
		}
		if bEnd.Before(aEnd) {
			aEnd = bEnd
		}
		if !aStart.After(aEnd) {
			result.CertainTimeOverlap = &timeOverlap{
				Start: aStart.UTC().Format(time.RFC3339Nano),
				End:   aEnd.UTC().Format(time.RFC3339Nano),
			}
		}
	}
	return json.NewEncoder(output).Encode(result)
}
