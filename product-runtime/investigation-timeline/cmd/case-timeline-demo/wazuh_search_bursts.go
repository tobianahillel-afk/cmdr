package main

import (
	"sort"
	"strings"
	"time"
)

// Derived solely from the completed Search Job's already authorized Wazuh
// hits, never from unfiltered Case snapshots. A burst is not an incident.
type wazuhSearchBurst struct {
	AgentID string `json:"agent_id"`
	RuleID string `json:"rule_id"`
	AlertCount int `json:"alert_count"`
	DistinctFeeds int `json:"distinct_feeds"`
	FirstUTC string `json:"first_utc"`
	LastUTC string `json:"last_utc"`
	WindowSeconds int `json:"window_seconds"`
	EntryIDs []string `json:"entry_ids"`
}

type wazuhSearchBurstEvent struct {
	At time.Time
	Feed int
	EntryID string
}

// Reuse the proven 60-second sliding-window/two-pointer method from the
// independent E10-INV-120A temporal-burst capability (PR #51), adapting it
// to return the exact subset of authorized Event Search hit identities.
// Bound by Event Search's existing 200-result maximum, not by a new scan.
func deriveWazuhSearchBursts(hits []wazuhSearchHit) []wazuhSearchBurst {
	if len(hits) < 3 || len(hits) > maxWazuhSearchMatches {
		return []wazuhSearchBurst{}
	}
	groups := make(map[string][]wazuhSearchBurstEvent)
	for _, hit := range hits {
		if hit.AgentID == "" || hit.RuleID == "" {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, hit.SourceTimeUTC)
		if err != nil {
			continue
		}
		key := hit.AgentID + "/" + hit.RuleID
		groups[key] = append(groups[key], wazuhSearchBurstEvent{
			At: at, Feed: hit.Feed, EntryID: hit.EntryID,
		})
	}
	bursts := make([]wazuhSearchBurst, 0)
	for key, events := range groups {
		agent, rule, _ := strings.Cut(key, "/")
		if len(events) < 3 || agent == "" || rule == "" {
			continue
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].At.Equal(events[j].At) {
				return events[i].EntryID < events[j].EntryID
			}
			return events[i].At.Before(events[j].At)
		})
		first, bestStart, bestEnd, bestCount, bestFeeds := 0, 0, 0, 0, 0
		feedCount := make(map[int]int)
		for right, event := range events {
			feedCount[event.Feed]++
			for event.At.Sub(events[first].At) > 60*time.Second {
				oldFeed := events[first].Feed
				feedCount[oldFeed]--
				if feedCount[oldFeed] == 0 {
					delete(feedCount, oldFeed)
				}
				first++
			}
			count := right-first+1
			if count >= 3 && (count > bestCount ||
				count == bestCount && len(feedCount) > bestFeeds) {
				bestStart, bestEnd, bestCount, bestFeeds = first, right, count, len(feedCount)
			}
		}
		if bestCount == 0 {
			continue
		}
		burst := wazuhSearchBurst{
			AgentID: agent, RuleID: rule, AlertCount: bestCount,
			DistinctFeeds: bestFeeds, WindowSeconds: 60,
			FirstUTC: events[bestStart].At.UTC().Format(time.RFC3339Nano),
			LastUTC: events[bestEnd].At.UTC().Format(time.RFC3339Nano),
			EntryIDs: make([]string, 0, bestCount),
		}
		for _, event := range events[bestStart:bestEnd+1] {
			burst.EntryIDs = append(burst.EntryIDs, event.EntryID)
		}
		bursts = append(bursts, burst)
	}
	sort.Slice(bursts, func(i,j int) bool {
		a,b:=bursts[i],bursts[j]
		if a.AlertCount!=b.AlertCount { return a.AlertCount>b.AlertCount }
		if a.DistinctFeeds!=b.DistinctFeeds { return a.DistinctFeeds>b.DistinctFeeds }
		if a.FirstUTC!=b.FirstUTC { return a.FirstUTC<b.FirstUTC }
		if a.AgentID!=b.AgentID { return a.AgentID<b.AgentID }
		return a.RuleID<b.RuleID
	})
	if len(bursts)>10 { bursts=bursts[:10] }
	return bursts
}
