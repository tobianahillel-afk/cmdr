package main

import (
	"errors"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A limited navigation lead between TWO independently authorized Case owners.
// Entry IDs remain scoped by First/Second Case; they are not global event
// identities, nor proof that one alert caused the other.
type sharedAuthorizedWazuhEventPair struct {
	RuleID        string `json:"rule_id"`
	AgentID       string `json:"agent_id"`
	FirstEntryID  string `json:"first_entry_id"`
	SecondEntryID string `json:"second_entry_id"`
	FirstTimeUTC  string `json:"first_time_utc"`
	SecondTimeUTC string `json:"second_time_utc"`
	OffsetSeconds int    `json:"offset_seconds"`
}

// Inputs are already narrowed and independently verified by
// loadAuthorizedImportedCase, then completely traversed via the existing
// first-party owner in compareVerifiedSavedCases. Reuse the scoped source
// metadata; NEVER join raw JSONL bytes or the two source owner graphs.
// For each first-Case Wazuh event choose the closest second-Case Wazuh event
// with the same rule+agent within inclusive 60s; return at most 10 pairs.
func sharedAuthorizedWazuhEventPairs(first, second timeline.MemoryTimelineData,
	firstCase, secondCase string) ([]sharedAuthorizedWazuhEventPair, error) {
	type cohort struct { rule, agent string }
	type observation struct { id string; at time.Time }
	collect := func(data timeline.MemoryTimelineData, caseID string) (map[cohort][]observation, error) {
		imported, err := importedNormalizedIDs(data, caseID)
		if err != nil { return nil, err }
		sources := make(map[string]timeline.MemoryTimelineSource, len(data.Sources))
		for _, source := range data.Sources {
			if source.Identity.CaseID == caseID && source.Identity.Kind == "telemetry-event" {
				sources[source.Identity.Reference] = source
			}
		}
		groups := make(map[cohort][]observation)
		for _, entry := range data.Entries {
			if entry.CaseID != caseID { continue }
			if _, ok := imported[entry.Value.ID]; !ok { continue }
			source, ok := sources[entry.Value.SourceRef]
			if !ok || source.Wazuh == nil { continue }
			meta := source.Wazuh
			if meta.RuleID == "" || meta.AgentID == "" { continue }
			if entry.Value.Kind != "observed" || entry.Value.ClockUncertaintyMS != 0 {
				return nil, errors.New("Wazuh Case entry time not certain")
			}
			at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
			if err != nil { return nil, errors.New("Wazuh source event time invalid") }
			key := cohort{rule:meta.RuleID, agent:meta.AgentID}
			groups[key] = append(groups[key], observation{id:entry.Value.ID,at:at.UTC()})
		}
		for key := range groups {
			sort.Slice(groups[key],func(i,j int) bool {
				a,b:=groups[key][i],groups[key][j]
				if a.at.Equal(b.at) { return a.id<b.id }
				return a.at.Before(b.at)
			})
		}
		return groups,nil
	}
	left,err:=collect(first,firstCase)
	if err!=nil {return nil,err}
	right,err:=collect(second,secondCase)
	if err!=nil {return nil,err}
	const radius=60*time.Second
	pairs:=make([]sharedAuthorizedWazuhEventPair,0,10)
	// Iterate in deterministic chronology regardless of original file
	// ordering, feed numbering or Go map iteration.
	type selection struct { key cohort; first observation }
	var ordered []selection
	for key,events:=range left {
		if len(right[key])==0 {continue}
		for _,event:=range events {ordered=append(ordered,selection{key,event})}
	}
	sort.Slice(ordered,func(i,j int) bool {
		a,b:=ordered[i],ordered[j]
		if !a.first.at.Equal(b.first.at) {return a.first.at.Before(b.first.at)}
		if a.key.rule!=b.key.rule {return a.key.rule<b.key.rule}
		if a.key.agent!=b.key.agent {return a.key.agent<b.key.agent}
		return a.first.id<b.first.id
	})
	for _,s:=range ordered {
		candidates:=right[s.key]
		idx:=sort.Search(len(candidates),func(i int)bool {
			return !candidates[i].at.Before(s.first.at)
		})
		best:=-1
		for _,i:=range []int{idx-1,idx} {
			if i<0 || i>=len(candidates) {continue}
			delta:=candidates[i].at.Sub(s.first.at)
			if delta < -radius || delta > radius {continue}
			if best<0 || absTimeOffset(delta)<absTimeOffset(candidates[best].at.Sub(s.first.at)) ||
				(absTimeOffset(delta)==absTimeOffset(candidates[best].at.Sub(s.first.at)) &&
					candidates[i].id<candidates[best].id) {
				best=i
			}
		}
		if best<0 {continue}
		b:=candidates[best]
		pairs=append(pairs,sharedAuthorizedWazuhEventPair{
			RuleID:s.key.rule,AgentID:s.key.agent,
			FirstEntryID:s.first.id, SecondEntryID:b.id,
			FirstTimeUTC:s.first.at.Format(time.RFC3339Nano),
			SecondTimeUTC:b.at.Format(time.RFC3339Nano),
			OffsetSeconds:int(b.at.Sub(s.first.at)/time.Second),
		})
		if len(pairs)==10 {break}
	}
	return pairs,nil
}

func absTimeOffset(d time.Duration) time.Duration {
	if d<0 {return -d}
	return d
}
