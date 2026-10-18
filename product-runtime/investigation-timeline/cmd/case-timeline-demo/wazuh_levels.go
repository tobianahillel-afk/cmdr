package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Source metadata, not an authorization claim or a derived severity score.
// Wazuh 4.x rule level is an optional numeric field (0..16).
type wazuhRuleLevel struct {
	Level   int
	Known   bool
	RuleID  string
	AgentID string
	Feed    int // Server-assigned import ordinal; zero for the single-feed mode.
}

// Co-occurrence describes an agent/rule pair observed on separate offline
// operator feeds. It is not a causal correlation or evidence of compromise.
type wazuhMultiFeedSignal struct {
	AgentID               string `json:"agent_id"`
	RuleID                string `json:"rule_id"`
	AlertCount            int    `json:"alert_count"`
	DistinctFeeds         int    `json:"distinct_feeds"`
	HighestKnownRuleLevel *int   `json:"highest_known_rule_level,omitempty"`
}

// Host/agent grouping is derived only from already authorized alert entries.
type wazuhAgentCount struct {
	AgentID               string `json:"agent_id"`
	Count                 int    `json:"count"`
	DistinctRuleIDs       int    `json:"distinct_rule_ids"`
	HighestKnownRuleLevel *int   `json:"highest_known_rule_level,omitempty"`
}

type wazuhLevelCount struct {
	Level int `json:"level"`
	Count int `json:"count"`
}

// An actionable aggregate of authorized imported alerts sharing a Wazuh
// detection rule ID. It intentionally omits alert IDs and source previews.
type wazuhRuleIDCount struct {
	RuleID                string `json:"rule_id"`
	Count                 int    `json:"count"`
	HighestKnownRuleLevel *int   `json:"highest_known_rule_level,omitempty"`
	FirstCertainUTC       string `json:"first_certain_utc,omitempty"`
	LastCertainUTC        string `json:"last_certain_utc,omitempty"`
}

// A bounded, descriptive time distribution of visible imported alerts.
// These are alert-volume windows, not inferred incident boundaries.
type wazuhAlertWindow struct {
	StartUTC string `json:"start_utc"`
	Count int `json:"count"`
	HighestKnownRuleLevel *int `json:"highest_known_rule_level,omitempty"`
	UnknownRuleLevels int `json:"unknown_rule_levels"`
}

type wazuhRuleTriageReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	ImportedEntries int `json:"imported_entries"`
	Case caseTriageReceipt `json:"case_triage"`
	RuleLevels []wazuhLevelCount `json:"rule_levels"`
	UnknownRuleLevels int `json:"unknown_rule_levels"`
	DistinctRuleIDs int `json:"distinct_rule_ids"`
	UnknownRuleIDs int `json:"unknown_rule_ids"`
	TopRules []wazuhRuleIDCount `json:"top_rules"`
	DistinctAgentIDs int `json:"distinct_agent_ids"`
	UnknownAgentIDs int `json:"unknown_agent_ids"`
	TopAgents []wazuhAgentCount `json:"top_agents"`
	MultiFeedSignals []wazuhMultiFeedSignal `json:"multi_feed_signals,omitempty"`
	TemporalBursts []wazuhAlertBurst `json:"temporal_bursts,omitempty"`
	DistinctTimeWindows int `json:"distinct_time_windows"`
	PeakWindows []wazuhAlertWindow `json:"peak_windows"`
	UnplacedAlerts int `json:"unplaced_alerts"`
}

// This mapping is computed from the SAME bytes that were normalized and
// imported into the Case; no second file open or independent ACL is introduced.
func wazuhRuleLevels(raw []byte, caseID string) (map[string]wazuhRuleLevel, error) {
	return wazuhRuleLevelsWithPrefix(raw, caseID, "")
}

// prefix is supplied only by the operator-side feed loader, never by
// Wazuh input or an HTTP request. It matches the already-imported entry ID.
func wazuhRuleLevelsWithPrefix(raw []byte, caseID, prefix string) (map[string]wazuhRuleLevel, error) {
	result := make(map[string]wazuhRuleLevel)
	for _, line := range bytes.Split(raw, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var alert struct {
			ID string `json:"id"`
			Agent struct { ID string `json:"id"` } `json:"agent"`
			Rule struct {
				Level *int   `json:"level"`
				ID    string `json:"id"`
			} `json:"rule"`
		}
		if err := json.Unmarshal(line, &alert); err != nil {
			return nil, errors.New("invalid Wazuh rule level metadata")
		}
		digest := sha256.Sum256([]byte(caseID + "\x00" + prefix + "wazuh-" + alert.ID))
		entryID := fmt.Sprintf("local-jsonl-entry-%x", digest)
		if _, exists := result[entryID]; exists {
			return nil, errors.New("duplicate Wazuh rule identity")
		}
		level := wazuhRuleLevel{RuleID: alert.Rule.ID, AgentID: alert.Agent.ID}
		// Wazuh 4.x rule.id is a decimal string, not an ACL or alert ID.
		// Restrict its width and alphabet before including it in a receipt.
		if len(level.RuleID) > 16 {
			return nil, errors.New("invalid Wazuh rule ID")
		}
		for _, digit := range level.RuleID {
			if digit < '0' || digit > '9' {
				return nil, errors.New("invalid Wazuh rule ID")
			}
		}
		// Agent IDs are the stable Wazuh 4.x strings (including 000).
		// Never place an agent name, IP, hostname or full_log in this receipt.
		if len(level.AgentID) > 16 {
			return nil, errors.New("invalid Wazuh agent ID")
		}
		for _, digit := range level.AgentID {
			if digit < '0' || digit > '9' {
				return nil, errors.New("invalid Wazuh agent ID")
			}
		}
		if alert.Rule.Level != nil {
			if *alert.Rule.Level < 0 || *alert.Rule.Level > 16 {
				return nil, errors.New("Wazuh rule level outside 0..16")
			}
			level.Level, level.Known = *alert.Rule.Level, true
		}
		result[entryID] = level
	}
	return result, nil
}

// writeDemoWazuh4RuleTriage adds genuine source rule-level analysis to the
// existing full authenticated Case/Timeline/typed-source pagination path.
// Metadata is counted ONLY when its imported Timeline Entry is authorized.
func writeDemoWazuh4RuleTriage(ctx context.Context, principal, snapshotPath, intoCase, alertsPath string, output io.Writer) error {
	merged, imported, raw, err := prepareWazuh4CaseWithRaw(ctx, principal, snapshotPath, intoCase, alertsPath)
	if err != nil {
		return err
	}
	return writeDemoWazuh4RuleTriageWithData(ctx, principal, merged, imported, raw, intoCase, output)
}

// Reuse the same Case/Timeline/typed-source traversal for offline CLI and
// authenticated loopback UI without opening the operator's files twice.
func writeDemoWazuh4RuleTriageWithData(ctx context.Context, principal string, merged timeline.MemoryTimelineData, imported int, raw []byte, intoCase string, output io.Writer) error {
	levels, err := wazuhRuleLevels(raw, intoCase)
	if err != nil {
		return err
	}
	return writeDemoWazuh4RuleTriageWithLevels(ctx, principal, merged, imported, levels, intoCase, output)
}

// The single- and multi-feed paths share the same authenticated,
// per-entry typed-source pagination and aggregate evaluator.
func writeDemoWazuh4RuleTriageWithLevels(ctx context.Context, principal string, merged timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, intoCase string, output io.Writer) error {
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, merged)
	if err != nil {
		return err
	}
	counts := make(map[int]int)
	windows := make(map[string]*wazuhAlertWindow)
	rules := make(map[string]*wazuhRuleIDCount)
	agents := make(map[string]*wazuhAgentCount)
	agentRules := make(map[string]map[string]bool)
	multiPairs := make(map[string]*wazuhMultiFeedSignal)
	multiFeeds := make(map[string]map[int]bool)
	burstGroups := make(map[string][]wazuhBurstObservation)
	visibleImported, unknown, unplaced, unknownRules, unknownAgents := 0, 0, 0, 0, 0
	receipt, err := collectDemoCaseTriageWithObserver(ctx, handler, token, intoCase, func(entry timeline.ProjectedEntry) error {
		level, exists := levels[entry.ID]
		if !exists {
			return nil
		}
		visibleImported++
		// Join feed identity only after server-authorized typed-source traversal.
		// Both IDs are validated decimal strings in wazuhRuleLevelsWithPrefix.
		if level.Feed > 0 && level.AgentID != "" && level.RuleID != "" {
			key := level.AgentID + "/" + level.RuleID
			signal := multiPairs[key]
			if signal == nil {
				signal = &wazuhMultiFeedSignal{AgentID: level.AgentID, RuleID: level.RuleID}
				multiPairs[key] = signal
				multiFeeds[key] = make(map[int]bool)
			}
			signal.AlertCount++
			multiFeeds[key][level.Feed] = true
			if level.Known && (signal.HighestKnownRuleLevel == nil || level.Level > *signal.HighestKnownRuleLevel) {
				maximum := level.Level
				signal.HighestKnownRuleLevel = &maximum
			}
		}
		if level.Known {
			counts[level.Level]++
		} else {
			unknown++
		}
		if level.AgentID == "" {
			unknownAgents++
		} else {
			agent := agents[level.AgentID]
			if agent == nil {
				agent = &wazuhAgentCount{AgentID: level.AgentID}
				agents[level.AgentID] = agent
				agentRules[level.AgentID] = make(map[string]bool)
			}
			agent.Count++
			if level.RuleID != "" && !agentRules[level.AgentID][level.RuleID] {
				agentRules[level.AgentID][level.RuleID] = true
				agent.DistinctRuleIDs++
			}
			if level.Known && (agent.HighestKnownRuleLevel == nil || level.Level > *agent.HighestKnownRuleLevel) {
				highest := level.Level
				agent.HighestKnownRuleLevel = &highest
			}
		}
		var rule *wazuhRuleIDCount
		if level.RuleID == "" {
			unknownRules++
		} else {
			rule = rules[level.RuleID]
			if rule == nil {
				rule = &wazuhRuleIDCount{RuleID: level.RuleID}
				rules[level.RuleID] = rule
			}
			rule.Count++
			if level.Known && (rule.HighestKnownRuleLevel == nil || level.Level > *rule.HighestKnownRuleLevel) {
				highest := level.Level
				rule.HighestKnownRuleLevel = &highest
			}
		}
		if entry.TimeUncertain {
			unplaced++
			return nil
		}
		when, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil {
			return errors.New("invalid authorized Wazuh Timeline timestamp")
		}
		// The first-party observer runs only after Case/Timeline and typed
		// source authorization. Uncertain timestamps never contribute.
		if level.AgentID != "" && level.RuleID != "" {
			key := level.AgentID + "/" + level.RuleID
			burstGroups[key] = append(burstGroups[key], wazuhBurstObservation{
				At: when.UTC(), Feed: level.Feed, EntryID: entry.ID,
			})
		}
		// Fixed-width UTC precision keeps lexicographic min/max chronological.
		utc := when.UTC().Format("2006-01-02T15:04:05.000000000Z")
		if rule != nil {
			if rule.FirstCertainUTC == "" || utc < rule.FirstCertainUTC {
				rule.FirstCertainUTC = utc
			}
			if rule.LastCertainUTC == "" || utc > rule.LastCertainUTC {
				rule.LastCertainUTC = utc
			}
		}
		key := when.UTC().Truncate(5 * time.Minute).Format(time.RFC3339)
		window := windows[key]
		if window == nil {
			window = &wazuhAlertWindow{StartUTC: key}
			windows[key] = window
		}
		window.Count++
		if level.Known {
			if window.HighestKnownRuleLevel == nil || level.Level > *window.HighestKnownRuleLevel {
				highest := level.Level
				window.HighestKnownRuleLevel = &highest
			}
		} else {
			window.UnknownRuleLevels++
		}
		return nil
	})
	if err != nil {
		return err
	}
	if visibleImported != imported || visibleImported != len(levels) {
		return errors.New("not every imported Wazuh alert is visible to the authorized Case reader")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	keys := make([]int, 0, len(counts))
	for level := range counts {
		keys = append(keys, level)
	}
	sort.Ints(keys)
	peaks := make([]wazuhAlertWindow, 0, len(windows))
	for _, window := range windows {
		peaks = append(peaks, *window)
	}
	sort.Slice(peaks, func(i, j int) bool {
		if peaks[i].Count != peaks[j].Count {
			return peaks[i].Count > peaks[j].Count
		}
		return peaks[i].StartUTC < peaks[j].StartUTC
	})
	// Export at most ten active windows; all windows are still counted.
	if len(peaks) > 10 {
		peaks = peaks[:10]
	}
	// Expose the busiest detection rule cohorts, never raw Wazuh alert IDs.
	// Exact seconds are supplied only for known authorized timestamps.
	topRules := make([]wazuhRuleIDCount, 0, len(rules))
	for _, rule := range rules {
		topRules = append(topRules, *rule)
	}
	sort.Slice(topRules, func(i, j int) bool {
		if topRules[i].Count != topRules[j].Count {
			return topRules[i].Count > topRules[j].Count
		}
		return topRules[i].RuleID < topRules[j].RuleID
	})
	if len(topRules) > 10 {
		topRules = topRules[:10]
	}
	// Rank affected Wazuh agents for an operator-visible host pivot.
	topAgents := make([]wazuhAgentCount, 0, len(agents))
	for _, agent := range agents {
		topAgents = append(topAgents, *agent)
	}
	sort.Slice(topAgents, func(i, j int) bool {
		if topAgents[i].Count != topAgents[j].Count {
			return topAgents[i].Count > topAgents[j].Count
		}
		return topAgents[i].AgentID < topAgents[j].AgentID
	})
	if len(topAgents) > 10 {
		topAgents = topAgents[:10]
	}
	// A signal requires independently imported observations in at least two
	// distinct feeds; repeated alerts within one feed are not cross-feed.
	multiSignals := make([]wazuhMultiFeedSignal, 0)
	for key, signal := range multiPairs {
		if len(multiFeeds[key]) < 2 {
			continue
		}
		signal.DistinctFeeds = len(multiFeeds[key])
		multiSignals = append(multiSignals, *signal)
	}
	sort.Slice(multiSignals, func(i, j int) bool {
		if multiSignals[i].DistinctFeeds != multiSignals[j].DistinctFeeds {
			return multiSignals[i].DistinctFeeds > multiSignals[j].DistinctFeeds
		}
		if multiSignals[i].AlertCount != multiSignals[j].AlertCount {
			return multiSignals[i].AlertCount > multiSignals[j].AlertCount
		}
		if multiSignals[i].AgentID != multiSignals[j].AgentID {
			return multiSignals[i].AgentID < multiSignals[j].AgentID
		}
		return multiSignals[i].RuleID < multiSignals[j].RuleID
	})
	if len(multiSignals) > 10 {
		multiSignals = multiSignals[:10]
	}
	histogram := make([]wazuhLevelCount, 0, len(keys))
	for _, level := range keys {
		histogram = append(histogram, wazuhLevelCount{Level: level, Count: counts[level]})
	}
	return json.NewEncoder(output).Encode(wazuhRuleTriageReceipt{
		Format: "cmdr-local-wazuh4-case-triage-v1",
		Synthetic: true, ImportedEntries: imported, Case: receipt,
		RuleLevels: histogram, UnknownRuleLevels: unknown,
		DistinctRuleIDs: len(rules), UnknownRuleIDs: unknownRules, TopRules: topRules,
		DistinctAgentIDs: len(agents), UnknownAgentIDs: unknownAgents, TopAgents: topAgents,
		MultiFeedSignals: multiSignals,
		TemporalBursts: deriveWazuhTemporalBursts(burstGroups),
		DistinctTimeWindows: len(windows), PeakWindows: peaks, UnplacedAlerts: unplaced,
	})
}
