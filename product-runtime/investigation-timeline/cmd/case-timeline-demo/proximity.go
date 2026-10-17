package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"
)

// A local, unsigned, synthetic analyst lead. A close timestamp is NOT a
// correlation, shared source identity, incident link, or causality finding.
type proximityLead struct {
	FirstEntryID   string `json:"first_entry_id"`
	SecondEntryID  string `json:"second_entry_id"`
	FirstKind      string `json:"first_kind"`
	SecondKind     string `json:"second_kind"`
	FirstTime      string `json:"first_source_time"`
	SecondTime     string `json:"second_source_time"`
	SeparationMS   int64  `json:"separation_ms"`
}

type localCaseProximity struct {
	Format                  string          `json:"format"`
	Synthetic               bool            `json:"synthetic"`
	FirstCase               string          `json:"first_case"`
	SecondCase              string          `json:"second_case"`
	WindowSeconds           int             `json:"window_seconds"`
	FirstCertainEntries     int             `json:"first_certain_entries"`
	SecondCertainEntries    int             `json:"second_certain_entries"`
	FirstUncertainExcluded  int             `json:"first_uncertain_excluded"`
	SecondUncertainExcluded int             `json:"second_uncertain_excluded"`
	Selection               string          `json:"selection"`
	Leads                   []proximityLead `json:"leads"`
}

// writeDemoCaseProximity authenticates and completely loads BOTH Cases and
// typed sources via existing HTTP owner handlers. Never output partial data.
func writeDemoCaseProximity(ctx context.Context, principal, snapshot, firstID, secondID string, windowSeconds int, out io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || firstID == "" || secondID == "" || firstID == secondID ||
		len(firstID) > 1024 || len(secondID) > 1024 || windowSeconds < 1 || windowSeconds > 300 {
		return errors.New("select an authorized principal, two distinct Cases and a 1..300 second window")
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
	events, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, firstID, secondID)
	if err != nil {
		return err
	}
	result, err := selectCaseProximity(ctx, events, windowSeconds)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(result)
}

type proximityEvent struct {
	id string
	kind string
	sourceTime string
	time time.Time
}

func selectCaseProximity(ctx context.Context, data localCaseEventSequence, seconds int) (localCaseProximity, error) {
	if seconds < 1 || seconds > 300 || data.First.CaseID == "" ||
		data.Second.CaseID == "" || data.First.CaseID == data.Second.CaseID ||
		len(data.Events) > 10000 {
		return localCaseProximity{}, errors.New("invalid bounded authorized Case sequence")
	}
	result := localCaseProximity{
		Format: "cmdr-local-cross-case-proximity-v1", Synthetic: true,
		FirstCase: data.First.CaseID, SecondCase: data.Second.CaseID,
		WindowSeconds: seconds,
		Selection: "nearest certain event in second Case for each certain first-Case event; at most 25 globally closest leads, not exhaustive pairs",
		Leads: make([]proximityLead, 0),
	}
	var first, second []proximityEvent
	for _, row := range data.Events {
		if err := ctx.Err(); err != nil {
			return localCaseProximity{}, err
		}
		id := row.Entry.ID
		if id == "" || (row.CaseID != result.FirstCase && row.CaseID != result.SecondCase) {
			return localCaseProximity{}, errors.New("inconsistent authorized Case event")
		}
		if row.Entry.TimeUncertain {
			if row.CaseID == result.FirstCase { result.FirstUncertainExcluded++ } else { result.SecondUncertainExcluded++ }
			continue
		}
		instant, err := time.Parse(time.RFC3339Nano, row.Entry.SourceTime)
		if err != nil {
			return localCaseProximity{}, errors.New("invalid certain event clock")
		}
		event := proximityEvent{id: id, kind: row.Entry.Kind, sourceTime: row.Entry.SourceTime, time: instant}
		if row.CaseID == result.FirstCase {
			first = append(first, event)
		} else {
			second = append(second, event)
		}
	}
	result.FirstCertainEntries, result.SecondCertainEntries = len(first), len(second)
	// The upstream combined sequence is sorted, but keep local slices ordered
	// independently so the search remains correct for independent source data.
	sort.Slice(first, func(i,j int) bool { return first[i].time.Before(first[j].time) ||
		first[i].time.Equal(first[j].time) && first[i].id < first[j].id })
	sort.Slice(second, func(i,j int) bool { return second[i].time.Before(second[j].time) ||
		second[i].time.Equal(second[j].time) && second[i].id < second[j].id })
	maxGap := time.Duration(seconds)*time.Second
	candidates := make([]proximityLead, 0, len(first))
	for _, left := range first {
		if err := ctx.Err(); err != nil { return localCaseProximity{}, err }
		if len(second) == 0 { break }
		at := sort.Search(len(second), func(i int) bool { return !second[i].time.Before(left.time) })
		previous := at - 1
		if previous >= 0 {
			// For equal prior timestamps pick the first ID, not the last.
			previousInstant := second[previous].time
			previous = sort.Search(len(second), func(i int) bool {
				return !second[i].time.Before(previousInstant)
			})
		}
		var best proximityLead
		found := false
		for _, index := range []int{previous, at} {
			if index < 0 || index >= len(second) { continue }
			right := second[index]
			// Avoid signed Duration overflow for events centuries apart.
			var delta time.Duration
			if left.time.Before(right.time) {
				delta = right.time.Sub(left.time)
			} else {
				delta = left.time.Sub(right.time)
			}
			if delta > maxGap { continue }
			candidate := proximityLead{
				FirstEntryID: left.id, SecondEntryID: right.id,
				FirstKind: left.kind, SecondKind: right.kind,
				FirstTime: left.sourceTime, SecondTime: right.sourceTime,
				SeparationMS: delta.Milliseconds(),
			}
			if !found || candidate.SeparationMS < best.SeparationMS ||
				candidate.SeparationMS == best.SeparationMS && candidate.SecondEntryID < best.SecondEntryID {
				best, found = candidate, true
			}
		}
		if found { candidates = append(candidates, best) }
	}
	sort.Slice(candidates,func(i,j int)bool{
		a,b:=candidates[i],candidates[j]
		if a.SeparationMS!=b.SeparationMS{return a.SeparationMS<b.SeparationMS}
		if a.FirstEntryID!=b.FirstEntryID{return a.FirstEntryID<b.FirstEntryID}
		return a.SecondEntryID<b.SecondEntryID
	})
	if len(candidates)>25{candidates=candidates[:25]}
	result.Leads=append(result.Leads,candidates...)
	return result,nil
}
