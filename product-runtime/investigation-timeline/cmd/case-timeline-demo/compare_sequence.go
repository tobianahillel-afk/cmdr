package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// An explicitly requested sensitive local receipt; no ownership or custody claim.
type localCaseEventSequence struct {
	Format           string            `json:"format"`
	WindowStartUTC   string            `json:"window_start_utc,omitempty"`
	WindowEndUTC     string            `json:"window_end_utc,omitempty"`
	MatchedFirst     int               `json:"matched_first,omitempty"`
	MatchedSecond    int               `json:"matched_second,omitempty"`
	Synthetic        bool              `json:"synthetic"`
	First            caseTriageReceipt `json:"first"`
	Second           caseTriageReceipt `json:"second"`
	CertainEntries   int               `json:"certain_entries"`
	UncertainEntries int               `json:"uncertain_entries"`
	Events           []localCaseEvent  `json:"events"`
}

type localCaseEvent struct {
	CaseID string                  `json:"case_id"`
	Entry  timeline.ProjectedEntry `json:"entry"`
}

type orderedCaseEvent struct {
	value   localCaseEvent
	instant time.Time
	known   bool
}

// Both independent Case/Timeline/source ACL traversals must complete before
// any combined event is returned. The owner provider, not these rows, grants
// access. Event matches and causal links are never inferred.
func collectAuthorizedCrossCaseSequence(ctx context.Context, handler http.Handler, token, firstID, secondID string) (localCaseEventSequence, error) {
	if firstID == "" || secondID == "" || firstID == secondID ||
		len(firstID) > 1024 || len(secondID) > 1024 {
		return localCaseEventSequence{}, errors.New("two distinct bounded Case IDs required")
	}
	rows := make([]orderedCaseEvent, 0)
	appendCase := func(caseID string) func(timeline.ProjectedEntry) error {
		return func(entry timeline.ProjectedEntry) error {
			instant, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
			if !entry.TimeUncertain && err != nil {
				return errors.New("invalid certain Timeline timestamp")
			}
			rows = append(rows, orderedCaseEvent{
				value: localCaseEvent{CaseID: caseID, Entry: entry},
				instant: instant, known: !entry.TimeUncertain && err == nil,
			})
			return nil
		}
	}
	first, err := collectDemoCaseTriageWithObserver(ctx, handler, token, firstID, appendCase(firstID))
	if err != nil {
		return localCaseEventSequence{}, err
	}
	second, err := collectDemoCaseTriageWithObserver(ctx, handler, token, secondID, appendCase(secondID))
	if err != nil {
		return localCaseEventSequence{}, err
	}
	if err := ctx.Err(); err != nil {
		return localCaseEventSequence{}, err
	}
	if len(rows) != first.VisibleEntries+second.VisibleEntries || len(rows) > 2*timeline.MaxEntries {
		return localCaseEventSequence{}, errors.New("incomplete authorized cross-Case sequence")
	}
	// Known instants first; unknown chronology never masquerades as UTC order.
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.known != b.known {
			return a.known
		}
		if a.known && !a.instant.Equal(b.instant) {
			return a.instant.Before(b.instant)
		}
		if a.value.CaseID != b.value.CaseID {
			return a.value.CaseID < b.value.CaseID
		}
		return a.value.Entry.ID < b.value.Entry.ID
	})
	result := localCaseEventSequence{
		Format: "cmdr-local-cross-case-event-sequence-v1", Synthetic: true,
		First: first, Second: second, Events: make([]localCaseEvent, 0, len(rows)),
	}
	for _, row := range rows {
		result.Events = append(result.Events, row.value)
		if row.known {
			result.CertainEntries++
		} else {
			result.UncertainEntries++
		}
	}
	return result, nil
}

// selectAuthorizedFiveMinuteWindow narrows an already completely authorized
// two-Case sequence. It does not infer a shared owner, cause or identity.
func selectAuthorizedFiveMinuteWindow(report localCaseEventSequence, selected string) (localCaseEventSequence, error) {
	start, err := time.Parse(time.RFC3339, selected)
	if err != nil || selected != start.UTC().Format(time.RFC3339) ||
		start.Minute()%5 != 0 || start.Second() != 0 {
		return localCaseEventSequence{}, errors.New("select a canonical UTC five-minute boundary")
	}
	end := start.Add(5 * time.Minute)
	filtered := make([]localCaseEvent, 0)
	firstCount, secondCount := 0, 0
	for _, item := range report.Events {
		if item.Entry.TimeUncertain {
			continue // uncertainty must never be promoted to a certain window
		}
		instant, err := time.Parse(time.RFC3339Nano, item.Entry.SourceTime)
		if err != nil {
			return localCaseEventSequence{}, errors.New("invalid certain authorized Timeline timestamp")
		}
		if instant.Before(start) || !instant.Before(end) {
			continue
		}
		filtered = append(filtered, item)
		if item.CaseID == report.First.CaseID {
			firstCount++
		} else if item.CaseID == report.Second.CaseID {
			secondCount++
		} else {
			return localCaseEventSequence{}, errors.New("unexpected Case in authorized sequence")
		}
	}
	if firstCount == 0 || secondCount == 0 {
		return localCaseEventSequence{}, errors.New("no certain events from both Cases in selected window")
	}
	report.Format = "cmdr-local-cross-case-window-events-v1"
	report.WindowStartUTC = start.Format(time.RFC3339)
	report.WindowEndUTC = end.Format(time.RFC3339)
	report.MatchedFirst, report.MatchedSecond = firstCount, secondCount
	report.Events = filtered
	report.CertainEntries = len(filtered)
	report.UncertainEntries = 0
	return report, nil
}

func writeDemoCaseEventSequence(ctx context.Context, principal, snapshot, firstID, secondID string, out io.Writer) error {
	return writeDemoCaseEventSequenceWithWindow(ctx, principal, snapshot, firstID, secondID, "", out)
}

func writeDemoCaseEventSequenceWithWindow(ctx context.Context, principal, snapshot, firstID, secondID, window string, out io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || firstID == "" || secondID == "" ||
		firstID == secondID || len(firstID) > 1024 || len(secondID) > 1024 {
		return errors.New("explicit principal and two distinct authorized Cases required")
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
	report, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, firstID, secondID)
	if err != nil {
		return err
	}
	if window != "" {
		report, err = selectAuthorizedFiveMinuteWindow(report, window)
		if err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(report)
}
