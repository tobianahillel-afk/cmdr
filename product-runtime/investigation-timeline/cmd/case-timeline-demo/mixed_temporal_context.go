package main

import (
	"errors"
	"sort"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// These are temporal leads, NOT incidents or claims of causal attribution.
// Every member must have come from ONE completed, source-reauthorized
// Event Search Job. No other Case entry is inspected for correlations.
type mixedTemporalNeighbour struct {
	EntryID string `json:"entry_id"`
	SourceTimeUTC string `json:"source_time_utc"`
	OffsetSeconds int `json:"offset_seconds"`
}
type mixedTemporalLead struct {
	WazuhEntryID string `json:"wazuh_entry_id"`
	WazuhTimeUTC string `json:"wazuh_time_utc"`
	WindowSeconds int `json:"window_seconds"`
	GenericEvents []mixedTemporalNeighbour `json:"generic_events"`
}

// A bounded derived view (<=200 authorized refs; <=10 alerts, 5 generic
// neighbours each). Do not query another owner or use a saved raw log.
// A source reference is actionable only via a fresh source-open request.
func deriveAuthorizedMixedTemporalContext(data timeline.MemoryTimelineData, caseID string,
	resultRefs []search.ResultReference, wazuh map[string]timeline.WazuhSourceMetadata) ([]mixedTemporalLead, error) {
	if len(resultRefs) > maxWazuhSearchMatches {
		return nil, errors.New("too many authorized search references")
	}
	type observation struct {
		id string
		at time.Time
	}
	visible := make(map[string]struct{}, len(resultRefs))
	for _, ref := range resultRefs {
		if ref.Source != localJSONLSearchSource || ref.Ref == "" {
			return nil, errors.New("foreign authorized search source")
		}
		if _, repeated := visible[ref.Ref]; repeated {
			return nil, errors.New("duplicate authorized search member")
		}
		visible[ref.Ref] = struct{}{}
	}
	entries := make(map[string]timeline.MemoryTimelineEntry, len(visible))
	for _, entry := range data.Entries {
		if entry.CaseID != caseID {
			continue
		}
		if _, ok := visible[entry.Value.ID]; ok {
			if _, duplicate := entries[entry.Value.ID]; duplicate {
				return nil, errors.New("duplicate owner entry")
			}
			entries[entry.Value.ID] = entry
		}
	}
	alerts := make([]observation, 0)
	generic := make([]observation, 0)
	for id := range visible {
		entry, ok := entries[id]
		if !ok || entry.Value.Kind != "observed" ||
			entry.Value.SourceRef == "" {
			return nil, errors.New("authorized Search Job member not in selected Case")
		}
		at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
		if err != nil {
			return nil, errors.New("authorized event has invalid source time")
		}
		point := observation{id: id, at: at.UTC()}
		if _, typedWazuh := wazuh[id]; typedWazuh {
			alerts = append(alerts, point)
		} else {
			generic = append(generic, point)
		}
	}
	sort.Slice(alerts, func(i,j int) bool {
		if alerts[i].at.Equal(alerts[j].at) {return alerts[i].id < alerts[j].id}
		return alerts[i].at.Before(alerts[j].at)
	})
	sort.Slice(generic, func(i,j int) bool {
		if generic[i].at.Equal(generic[j].at) {return generic[i].id < generic[j].id}
		return generic[i].at.Before(generic[j].at)
	})
	const radius = 60 * time.Second
	leads := make([]mixedTemporalLead, 0, 10)
	for _, alert := range alerts {
		neighbors := make([]mixedTemporalNeighbour, 0, 5)
		for _, event := range generic {
			delta := event.at.Sub(alert.at)
			if delta < -radius || delta > radius { continue }
			neighbors = append(neighbors, mixedTemporalNeighbour{
				EntryID: event.id, SourceTimeUTC: event.at.Format(time.RFC3339Nano),
				OffsetSeconds: int(delta / time.Second),
			})
		}
		if len(neighbors) == 0 {continue}
		sort.Slice(neighbors, func(i,j int) bool {
			abs := func(v int) int { if v < 0 {return -v}; return v }
			di,dj:=abs(neighbors[i].OffsetSeconds),abs(neighbors[j].OffsetSeconds)
			if di!=dj {return di<dj}
			if neighbors[i].SourceTimeUTC!=neighbors[j].SourceTimeUTC {
				return neighbors[i].SourceTimeUTC<neighbors[j].SourceTimeUTC
			}
			return neighbors[i].EntryID<neighbors[j].EntryID
		})
		if len(neighbors)>5 {neighbors=neighbors[:5]}
		leads=append(leads,mixedTemporalLead{
			WazuhEntryID:alert.id,WazuhTimeUTC:alert.at.Format(time.RFC3339Nano),
			WindowSeconds:60,GenericEvents:neighbors,
		})
		if len(leads)==10 {break}
	}
	return leads,nil
}
