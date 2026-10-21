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
	MITREIDs []string
}

// Host/agent grouping is derived only from already authorized alert entries.
type wazuhAgentCount struct {
	AgentID               string `json:"agent_id"`
	Count                 int    `json:"count"`
	DistinctRuleIDs       int    `json:"distinct_rule_ids"`
	HighestKnownRuleLevel *int   `json:"highest_known_rule_level,omitempty"`
}

// An ATT&CK technique label as claimed by authorized Wazuh rule metadata.
// This is a source attribution, never a verified behavioral classification.
type wazuhMITRETechniqueCount struct {
	TechniqueID string `json:"technique_id"`
	Count int `json:"count"`
	HighestKnownRuleLevel *int `json:"highest_known_rule_level,omitempty"`
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
	DistinctTechniqueIDs int `json:"distinct_technique_ids"`
	UnmappedTechniqueAlerts int `json:"unmapped_technique_alerts"`
	TopTechniques []wazuhMITRETechniqueCount `json:"top_techniques"`
	DistinctTimeWindows int `json:"distinct_time_windows"`
	PeakWindows []wazuhAlertWindow `json:"peak_windows"`
	UnplacedAlerts int `json:"unplaced_alerts"`
}

// This mapping is computed from the SAME bytes that were normalized and
// imported into the Case; no second file open or independent ACL is introduced.
func wazuhRuleLevels(raw []byte, caseID string) (map[string]wazuhRuleLevel, error) {
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
			MITRE struct {
				ID []string `json:"id"`
			} `json:"mitre"`
			} `json:"rule"`
		}
		if err := json.Unmarshal(line, &alert); err != nil {
			return nil, errors.New("invalid Wazuh rule level metadata")
		}
		digest := sha256.Sum256([]byte(caseID + "\x00wazuh-" + alert.ID))
		entryID := fmt.Sprintf("local-jsonl-entry-%x", digest)
		if _, exists := result[entryID]; exists {
			return nil, errors.New("duplicate Wazuh rule identity")
		}
		level := wazuhRuleLevel{RuleID: alert.Rule.ID, AgentID: alert.Agent.ID}
		// Source ATT&CK annotations are untrusted metadata, not IAM or
		// proof that a technique occurred. Deduplicate within each alert.
		if len(alert.Rule.MITRE.ID) > 16 {
			return nil, errors.New("too many Wazuh technique labels")
		}
		seen := make(map[string]bool, len(alert.Rule.MITRE.ID))
		for _, technique := range alert.Rule.MITRE.ID {
			if !validWazuhTechniqueID(technique) {
				return nil, errors.New("invalid Wazuh ATT&CK technique label")
			}
			if !seen[technique] {
				seen[technique] = true
				level.MITREIDs = append(level.MITREIDs, technique)
			}
		}
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

// Accept only canonical source labels T#### and T####.###, never arbitrary
// descriptions, URLs, executable expressions or user-supplied permission IDs.
func validWazuhTechniqueID(id string) bool {
	if len(id) != 5 && len(id) != 9 {
		return false
	}
	if id[0] != 'T' {
		return false
	}
	for i := 1; i < 5; i++ {
		if id[i] < '0' || id[i] > '9' {
			return false
		}
	}
	if len(id) == 9 {
		if id[5] != '.' {
			return false
		}
		for i := 6; i < 9; i++ {
			if id[i] < '0' || id[i] > '9' {
				return false
			}
		}
	}
	return true
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
	techniques := make(map[string]*wazuhMITRETechniqueCount)
	agentRules := make(map[string]map[string]bool)
	visibleImported, unknown, unplaced, unknownRules, unknownAgents, unmappedTechniques := 0, 0, 0, 0, 0, 0
	receipt, err := collectDemoCaseTriageWithObserver(ctx, handler, token, intoCase, func(entry timeline.ProjectedEntry) error {
		level, exists := levels[entry.ID]
		if !exists {
			return nil
		}
		visibleImported++
		if len(level.MITREIDs) == 0 {
			unmappedTechniques++
		}
		for _, techniqueID := range level.MITREIDs {
			technique := techniques[techniqueID]
			if technique == nil {
				technique = &wazuhMITRETechniqueCount{TechniqueID: techniqueID}
				techniques[techniqueID] = technique
			}
			technique.Count++
			if level.Known && (technique.HighestKnownRuleLevel == nil || level.Level > *technique.HighestKnownRuleLevel) {
				highest := level.Level
				technique.HighestKnownRuleLevel = &highest
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
	// A Wazuh alert may carry multiple ATT&CK labels, so technique counts
	// are tag occurrences across visible alerts, not an exclusive partition.
	topTechniques := make([]wazuhMITRETechniqueCount, 0, len(techniques))
	for _, technique := range techniques {
		topTechniques = append(topTechniques, *technique)
	}
	sort.Slice(topTechniques, func(i, j int) bool {
		if topTechniques[i].Count != topTechniques[j].Count {
			return topTechniques[i].Count > topTechniques[j].Count
		}
		return topTechniques[i].TechniqueID < topTechniques[j].TechniqueID
	})
	if len(topTechniques) > 10 {
		topTechniques = topTechniques[:10]
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
		DistinctTimeWindows: len(windows), PeakWindows: peaks, UnplacedAlerts: unplaced,
		DistinctTechniqueIDs: len(techniques), UnmappedTechniqueAlerts: unmappedTechniques, TopTechniques: topTechniques,
	})
}
