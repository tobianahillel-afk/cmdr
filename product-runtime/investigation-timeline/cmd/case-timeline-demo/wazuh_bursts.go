package main

import (
	"sort"
	"strings"
	"time"
)

// Both offline Wazuh triage and completed authorized Event Search Jobs use
// ONE 60-second concentration algorithm. This consolidates PR #51's
// independent temporal implementation with PR #58's source-linked search
// bursts, rather than allowing two divergent detection semantics.
const wazuhBurstWindow = 60 * time.Second
const wazuhBurstMinimum = 3

type wazuhBurstObservation struct {
	At time.Time
	Feed int
	EntryID string // set for SearchJob hits; never exported by offline triage
}

type wazuhAlertBurst struct {
	AgentID string `json:"agent_id"`
	RuleID string `json:"rule_id"`
	AlertCount int `json:"alert_count"`
	DistinctFeeds int `json:"distinct_feeds"`
	FirstUTC string `json:"first_certain_utc"`
	LastUTC string `json:"last_certain_utc"`
	WindowSeconds int `json:"window_seconds"`
}

// Sorts only authorized observations and finds the strongest inclusive
// 60-second window in O(n log n). On a count tie prefer the window spanning
// more distinct operator-selected feeds. Events remain sorted so the caller
// may return exact already-authorized members without recomputing a window.
func strongestWazuhBurstWindow(observations []wazuhBurstObservation) (int,int,int,int) {
	if len(observations) < wazuhBurstMinimum {
		return 0,0,0,0
	}
	sort.Slice(observations,func(i,j int) bool {
		if observations[i].At.Equal(observations[j].At) {
			if observations[i].EntryID != observations[j].EntryID {
				return observations[i].EntryID < observations[j].EntryID
			}
			return observations[i].Feed < observations[j].Feed
		}
		return observations[i].At.Before(observations[j].At)
	})
	left, bestStart, bestEnd, bestCount, bestFeeds := 0,0,0,0,0
	feedCounts := make(map[int]int)
	for right,observation := range observations {
		feedCounts[observation.Feed]++
		for observation.At.Sub(observations[left].At)>wazuhBurstWindow {
			old:=observations[left].Feed
			feedCounts[old]--
			if feedCounts[old]==0 {delete(feedCounts,old)}
			left++
		}
		count:=right-left+1
		if count>=wazuhBurstMinimum && (count>bestCount ||
			count==bestCount && len(feedCounts)>bestFeeds) {
			bestStart,bestEnd,bestCount,bestFeeds=left,right,count,len(feedCounts)
		}
	}
	return bestStart,bestEnd,bestCount,bestFeeds
}

// A source-free summary of only owner-authorized, time-certain imported
// Wazuh alerts, useful even without a separately executed Search Job.
func deriveWazuhTemporalBursts(groups map[string][]wazuhBurstObservation) []wazuhAlertBurst {
	bursts:=make([]wazuhAlertBurst,0)
	for key,observations:=range groups {
		agent,rule,ok:=strings.Cut(key,"/")
		if !ok || agent=="" || rule=="" {continue}
		first,last,count,feeds:=strongestWazuhBurstWindow(observations)
		if count==0 {continue}
		bursts=append(bursts,wazuhAlertBurst{
			AgentID:agent,RuleID:rule,AlertCount:count,DistinctFeeds:feeds,
			FirstUTC:observations[first].At.UTC().Format(time.RFC3339Nano),
			LastUTC:observations[last].At.UTC().Format(time.RFC3339Nano),
			WindowSeconds:int(wazuhBurstWindow/time.Second),
		})
	}
	sort.Slice(bursts,func(i,j int) bool {
		a,b:=bursts[i],bursts[j]
		if a.AlertCount!=b.AlertCount {return a.AlertCount>b.AlertCount}
		if a.DistinctFeeds!=b.DistinctFeeds {return a.DistinctFeeds>b.DistinctFeeds}
		if a.FirstUTC!=b.FirstUTC {return a.FirstUTC<b.FirstUTC}
		if a.AgentID!=b.AgentID {return a.AgentID<b.AgentID}
		return a.RuleID<b.RuleID
	})
	if len(bursts)>10 {bursts=bursts[:10]}
	return bursts
}
