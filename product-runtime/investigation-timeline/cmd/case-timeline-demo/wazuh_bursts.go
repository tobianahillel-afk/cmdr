package main

import (
	"sort"
	"strings"
	"time"
)

// A burst is a descriptive concentration of already-authorized, time-certain
// Wazuh alerts for the same agent and rule. It is NOT a security verdict.
const wazuhBurstWindow = 60 * time.Second
const wazuhBurstMinimum = 3

type wazuhBurstObservation struct {
	At   time.Time
	Feed int
}

type wazuhAlertBurst struct {
	AgentID       string `json:"agent_id"`
	RuleID        string `json:"rule_id"`
	AlertCount    int    `json:"alert_count"`
	DistinctFeeds int    `json:"distinct_feeds"`
	FirstUTC      string `json:"first_certain_utc"`
	LastUTC       string `json:"last_certain_utc"`
	WindowSeconds int    `json:"window_seconds"`
}

// For each agent/rule, select its strongest inclusive 60-second window.
// Two pointers and an in-window feed frequency map keep this O(n log n)
// per cohort; never infer a burst from unknown or uncertain timestamps.
// Wazuh agent/rule strings are validated decimal metadata, not permissions.
func deriveWazuhTemporalBursts(groups map[string][]wazuhBurstObservation) []wazuhAlertBurst {
	bursts := make([]wazuhAlertBurst, 0)
	for key, observations := range groups {
		agentID, ruleID, ok := strings.Cut(key, "/")
		if !ok || agentID == "" || ruleID == "" || len(observations) < wazuhBurstMinimum {
			continue
		}
		sort.Slice(observations, func(i, j int) bool {
			if observations[i].At.Equal(observations[j].At) {
				return observations[i].Feed < observations[j].Feed
			}
			return observations[i].At.Before(observations[j].At)
		})
		left := 0
		feedCounts := make(map[int]int)
		bestCount, bestFeeds := 0, 0
		var bestFirst, bestLast time.Time
		for right, observation := range observations {
			feedCounts[observation.Feed]++
			for observation.At.Sub(observations[left].At) > wazuhBurstWindow {
				old := observations[left].Feed
				feedCounts[old]--
				if feedCounts[old] == 0 {
					delete(feedCounts, old)
				}
				left++
			}
			count := right - left + 1
			if count < wazuhBurstMinimum {
				continue
			}
			if count > bestCount || (count == bestCount && len(feedCounts) > bestFeeds) {
				bestCount, bestFeeds = count, len(feedCounts)
				bestFirst, bestLast = observations[left].At, observation.At
			}
		}
		if bestCount == 0 {
			continue
		}
		bursts = append(bursts, wazuhAlertBurst{
			AgentID: agentID, RuleID: ruleID, AlertCount: bestCount,
			DistinctFeeds: bestFeeds,
			FirstUTC: bestFirst.UTC().Format(time.RFC3339Nano),
			LastUTC: bestLast.UTC().Format(time.RFC3339Nano),
			WindowSeconds: int(wazuhBurstWindow / time.Second),
		})
	}
	sort.Slice(bursts, func(i, j int) bool {
		a, b := bursts[i], bursts[j]
		if a.AlertCount != b.AlertCount {
			return a.AlertCount > b.AlertCount
		}
		if a.DistinctFeeds != b.DistinctFeeds {
			return a.DistinctFeeds > b.DistinctFeeds
		}
		if a.FirstUTC != b.FirstUTC {
			return a.FirstUTC < b.FirstUTC
		}
		if a.AgentID != b.AgentID {
			return a.AgentID < b.AgentID
		}
		return a.RuleID < b.RuleID
	})
	if len(bursts) > 10 {
		bursts = bursts[:10]
	}
	return bursts
}
