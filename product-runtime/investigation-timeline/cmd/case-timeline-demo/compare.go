package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
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
	SharedCertainWindows []sharedCertainWindow `json:"shared_certain_five_minute_windows"`
}

// This is an observational co-occurrence hint, not causal correlation.
// Source references may be Case-scoped and are not compared as identities.
type sharedCertainWindow struct {
	StartUTC       string `json:"start_utc"`
	FirstEntries   int    `json:"first_entries"`
	SecondEntries  int    `json:"second_entries"`
	FirstObserved  int    `json:"first_observed"`
	SecondObserved int    `json:"second_observed"`
	FirstInferred  int    `json:"first_inferred"`
	SecondInferred int    `json:"second_inferred"`
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
	result, err := collectAuthorizedCaseComparison(ctx, principal, data, firstID, secondID)
	if err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(result)
}

// Reuse this in-memory integration for both existing snapshot comparison and
// operator-local JSONL import. Both Case ACL traversals must succeed before
// the caller may write an aggregate receipt.
func collectAuthorizedCaseComparison(ctx context.Context, principal string, data timeline.MemoryTimelineData, firstID, secondID string) (localCaseComparison, error) {
	if principal == "" || firstID == "" || secondID == "" || firstID == secondID ||
		len(firstID) > 1024 || len(secondID) > 1024 {
		return localCaseComparison{}, errors.New("select two distinct authorized Cases and a principal")
	}
	token, err := randomBearer()
	if err != nil {
		return localCaseComparison{}, err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return localCaseComparison{}, err
	}
	first, err := collectDemoCaseTriage(ctx, handler, token, firstID)
	if err != nil {
		return localCaseComparison{}, err
	}
	second, err := collectDemoCaseTriage(ctx, handler, token, secondID)
	if err != nil {
		return localCaseComparison{}, err
	}
	if err := ctx.Err(); err != nil {
		return localCaseComparison{}, err
	}
	return buildCaseComparison(first, second)
}

// Aggregate-only comparison: no raw Entry, preview or source ACL bypass.
func buildCaseComparison(first, second caseTriageReceipt) (localCaseComparison, error) {
	result := localCaseComparison{
		Format:            "cmdr-local-case-comparison-v1",
		Synthetic:         true,
		First:             first,
		Second:            second,
		VisibleEntryDelta: first.VisibleEntries - second.VisibleEntries,
		ObservedDelta:     first.Kinds["observed"] - second.Kinds["observed"],
		InferredDelta:     first.Kinds["inferred"] - second.Kinds["inferred"],
		SharedCertainWindows: make([]sharedCertainWindow, 0),
	}
	// Compare precise, known five-minute bins, not just enclosing time spans.
	// All bins originate from fully validated, source-authorized Case pages.
	for start, a := range first.certainWindows {
		b, ok := second.certainWindows[start]
		if !ok { continue }
		result.SharedCertainWindows = append(result.SharedCertainWindows, sharedCertainWindow{
			StartUTC: start, FirstEntries: a.Total, SecondEntries: b.Total,
			FirstObserved: a.Observed, SecondObserved: b.Observed,
			FirstInferred: a.Inferred, SecondInferred: b.Inferred,
		})
	}
	sort.Slice(result.SharedCertainWindows, func(i, j int) bool {
		return result.SharedCertainWindows[i].StartUTC < result.SharedCertainWindows[j].StartUTC
	})
	if len(result.SharedCertainWindows) > 5 {
		result.SharedCertainWindows = result.SharedCertainWindows[:5]
	}
	// Unknown clock intervals are never promoted to a known overlap.
	if first.FirstCertainUTC != "" && first.LastCertainUTC != "" &&
		second.FirstCertainUTC != "" && second.LastCertainUTC != "" {
		aStart, aErr := time.Parse(time.RFC3339Nano, first.FirstCertainUTC)
		aEnd, bErr := time.Parse(time.RFC3339Nano, first.LastCertainUTC)
		bStart, cErr := time.Parse(time.RFC3339Nano, second.FirstCertainUTC)
		bEnd, dErr := time.Parse(time.RFC3339Nano, second.LastCertainUTC)
		if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
			return localCaseComparison{}, errors.New("inconsistent authorized time bounds")
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
	return result, nil
}
