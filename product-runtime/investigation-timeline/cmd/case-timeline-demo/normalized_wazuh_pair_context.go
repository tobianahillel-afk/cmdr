package main

import (
	"context"
	"errors"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// One generic source-owned, independently selectable lead close to a Wazuh
// pair endpoint. This is contextual co-occurrence, NOT causality or a verdict.
// Raw generic preview never appears in the comparative private report.
type authorizedGenericPairContext struct {
	EntryID       string `json:"entry_id"`
	TimeUTC       string `json:"time_utc"`
	OffsetSeconds int    `json:"offset_seconds"`
}

type authorizedGenericPairContextSummary struct {
	Entries   []authorizedGenericPairContext `json:"entries"`
	Truncated bool                           `json:"truncated"`
}

// Correlate locally owned normalized JSONL evidence with the already selected
// authorized Wazuh event within inclusive +/-60 seconds. Cases are processed
// independently; no source can inherit permissions from its counterpart.
// The input data is the scope-minimized result of loadAuthorizedImportedCase.
func nearbyAuthorizedGenericPairContext(ctx context.Context, principal string,
	data timeline.MemoryTimelineData, caseID string, imported map[string]struct{},
	centerUTC string) (authorizedGenericPairContextSummary, error) {
	if ctx == nil || ctx.Err() != nil {
		return authorizedGenericPairContextSummary{}, errors.New("generic context inspection cancelled")
	}
	center, err := time.Parse(time.RFC3339Nano, centerUTC)
	if err != nil {
		return authorizedGenericPairContextSummary{}, errors.New("paired Wazuh clock unavailable")
	}
	sources := make(map[string]timeline.MemoryTimelineSource, len(data.Sources))
	for _, src := range data.Sources {
		if src.Identity.CaseID == caseID {
			sources[src.Identity.Reference] = src
		}
	}
	type candidate struct {
		authorizedGenericPairContext
		at time.Time
		distance time.Duration
	}
	var candidates []candidate
	const radius = time.Minute
	for _, entry := range data.Entries {
		if entry.CaseID != caseID {
			continue
		}
		if _, ok := imported[entry.Value.ID]; !ok {
			continue
		}
		src, ok := sources[entry.Value.SourceRef]
		if !ok || src.Identity.Kind != "telemetry-event" || src.Wazuh != nil ||
			entry.Value.Kind != "observed" || entry.Value.ClockUncertaintyMS != 0 {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
		if err != nil {
			return authorizedGenericPairContextSummary{}, errors.New("generic Case source has invalid UTC clock")
		}
		offset := at.Sub(center)
		if offset < -radius || offset > radius {
			continue
		}
		distance := offset
		if distance < 0 { distance = -distance }
		candidates = append(candidates, candidate{
			authorizedGenericPairContext: authorizedGenericPairContext{
				EntryID: entry.Value.ID, TimeUTC: at.UTC().Format(time.RFC3339Nano),
				OffsetSeconds: int(offset / time.Second),
			},
			at: at, distance: distance,
		})
	}
	// Prioritize actual proximity, not Case catalogue insertion order; only
	// ten contextual leads are published, with deterministic tie-breaking.
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.distance != b.distance { return a.distance < b.distance }
		if !a.at.Equal(b.at) { return a.at.Before(b.at) }
		return a.EntryID < b.EntryID
	})
	result := authorizedGenericPairContextSummary{
		Entries: make([]authorizedGenericPairContext, 0, 10),
		Truncated: len(candidates) > 10,
	}
	for i, selected := range candidates {
		if i == 10 { break }
		// A canonical imported ID is only a selector. Reopen its owner on
		// the SAME independently authorized Case before disclosing even
		// contextual metadata; discard previews, never export raw events.
		opened, err := inspectDemoEntryWithData(ctx, principal, data, caseID, selected.EntryID)
		if err != nil || !opened.Synthetic || opened.CaseID != caseID ||
			opened.EntryID != selected.EntryID || opened.Kind != "telemetry-event" {
			return authorizedGenericPairContextSummary{}, errors.New("generic contextual source authorization failed")
		}
		result.Entries = append(result.Entries, selected.authorizedGenericPairContext)
	}
	if err := ctx.Err(); err != nil {
		return authorizedGenericPairContextSummary{}, err
	}
	return result, nil
}
