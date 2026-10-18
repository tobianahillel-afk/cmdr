package main

import (
	"sort"
	"time"

)

// Entry references remain scoped to their separately authorized Case owners.
// Pairing is a chronological navigation lead, not event identity or causality.
type sharedAuthorizedWazuhEventPair struct {
	RuleID string `json:"rule_id"`
	AgentID string `json:"agent_id"`
	FirstEntryID string `json:"first_entry_id"`
	SecondEntryID string `json:"second_entry_id"`
	FirstTimeUTC string `json:"first_time_utc"`
	SecondTimeUTC string `json:"second_time_utc"`
	OffsetSeconds int `json:"offset_seconds"`
}

// Uses the SAME imported-ID, typed-source metadata and certain UTC checks as
// five-minute correlation. Both owners were independently authenticated and
// fully traversed before the caller can emit a comparison receipt.
// Derive event-linked 60s leads from the SAME separately authorized
// Wazuh source observations already used for five-minute counts.
func sharedWazuhPairsFromObservations(left,right map[authorizedWazuhCohort][]authorizedWazuhObservation)[]sharedAuthorizedWazuhEventPair{
	pairs, _ := sharedWazuhPairsPageFromObservations(left, right, 0)
	return pairs
}

// Bounded, deterministic navigation of the SAME source-authorized pair leads
// used by the original ten-item receipt. The offset counts only eligible
// matched pairs, never denied or out-of-window observations. One extra item
// establishes whether another page exists without exposing raw events.
func sharedWazuhPairsPageFromObservations(left,right map[authorizedWazuhCohort][]authorizedWazuhObservation,offset int)([]sharedAuthorizedWazuhEventPair,bool){
	if offset<0 || offset>5000 {return nil,false}
	type selection struct { cohort authorizedWazuhCohort; event authorizedWazuhObservation }
	var ordered []selection
	for key,observations:=range left {
		if len(right[key])==0{continue}
		for _,event:=range observations {
			ordered=append(ordered,selection{cohort:key,event:event})
		}
	}
	sort.Slice(ordered,func(i,j int)bool{
		a,b:=ordered[i],ordered[j]
		if !a.event.at.Equal(b.event.at){return a.event.at.Before(b.event.at)}
		if a.cohort.rule!=b.cohort.rule{return a.cohort.rule<b.cohort.rule}
		if a.cohort.agent!=b.cohort.agent{return a.cohort.agent<b.cohort.agent}
		return a.event.id<b.event.id
	})
	const radius=60*time.Second
	pairs:=make([]sharedAuthorizedWazuhEventPair,0,11)
	eligible:=0
	for _,item:=range ordered {
		candidates:=right[item.cohort]
		idx:=sort.Search(len(candidates),func(i int)bool{
			return !candidates[i].at.Before(item.event.at)
		})
		best:=-1
		for _,i:=range []int{idx-1,idx} {
			if i<0 || i>=len(candidates){continue}
			delta:=candidates[i].at.Sub(item.event.at)
			if delta < -radius || delta > radius {continue}
			if best<0 || absTimeOffset(delta)<absTimeOffset(candidates[best].at.Sub(item.event.at)) ||
				(absTimeOffset(delta)==absTimeOffset(candidates[best].at.Sub(item.event.at)) &&
					candidates[i].id<candidates[best].id){
				best=i
			}
		}
		if best<0{continue}
		eligible++
		if eligible<=offset {continue}
		// The preceding group can contain multiple equal-time events.
		// Choose the smallest authorized Entry ID at the winning instant,
		// not the last ID before idx: stable pair navigation after reorder.
		chosenTime:=candidates[best].at
		best=sort.Search(len(candidates),func(i int)bool{
			return !candidates[i].at.Before(chosenTime)
		})
		other:=candidates[best]
		pairs=append(pairs,sharedAuthorizedWazuhEventPair{
			RuleID:item.cohort.rule,AgentID:item.cohort.agent,
			FirstEntryID:item.event.id,SecondEntryID:other.id,
			FirstTimeUTC:item.event.at.Format(time.RFC3339Nano),
			SecondTimeUTC:other.at.Format(time.RFC3339Nano),
			OffsetSeconds:int(other.at.Sub(item.event.at)/time.Second),
		})
		if len(pairs)==11{return pairs[:10],true}
	}
	return pairs,false
}

func absTimeOffset(d time.Duration)time.Duration{
	if d<0{return -d}
	return d
}
