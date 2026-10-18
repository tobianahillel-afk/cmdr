package main

import (
	"errors"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A retrospective lead: same authorized Wazuh rule and agent observed in
// BOTH independently verified Cases during the same fixed UTC five-minute
// window. This is co-occurrence, not an incident verdict or causal relation.
// No raw source text, host, IP, original Wazuh ID or cross-Case source ref
// is exposed through this limited receipt.
type sharedAuthorizedWazuhWindow struct {
	StartUTC     string `json:"start_utc"`
	RuleID       string `json:"rule_id"`
	AgentID      string `json:"agent_id"`
	FirstEvents  int    `json:"first_events"`
	SecondEvents int    `json:"second_events"`
}

// Counts are derived only from the SAME first-party owner graph that the
// Case/Timeline/source ACL chain and its source-open requests validated.
// A generic JSONL Case has no Wazuh owner metadata and yields no cohorts.
func sharedAuthorizedWazuhWindows(first, second timeline.MemoryTimelineData,
	firstCase, secondCase string) ([]sharedAuthorizedWazuhWindow, error) {
	type cohort struct { window, rule, agent string }
	count := func(data timeline.MemoryTimelineData, caseID string) (map[cohort]int, error) {
		imported, err := importedNormalizedIDs(data, caseID)
		if err != nil { return nil, err }
		sources := make(map[string]timeline.MemoryTimelineSource, len(data.Sources))
		for _, source := range data.Sources {
			if source.Identity.CaseID == caseID && source.Identity.Kind == "telemetry-event" {
				sources[source.Identity.Reference] = source
			}
		}
		counts := make(map[cohort]int)
		for _, entry := range data.Entries {
			if entry.CaseID != caseID { continue }
			if _, allowed := imported[entry.Value.ID]; !allowed { continue }
			source, exists := sources[entry.Value.SourceRef]
			if !exists || source.Wazuh == nil { continue }
			meta := source.Wazuh
			if meta.RuleID == "" || meta.AgentID == "" { continue }
			if entry.Value.Kind != "observed" || entry.Value.ClockUncertaintyMS != 0 {
				return nil, errors.New("Wazuh Case entry has uncertain or non-observed timestamp")
			}
			at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
			if err != nil { return nil, errors.New("Wazuh Case source timestamp invalid") }
			key := cohort{
				window: at.UTC().Truncate(5*time.Minute).Format(time.RFC3339),
				rule: meta.RuleID, agent: meta.AgentID,
			}
			counts[key]++
		}
		return counts, nil
	}
	a, err := count(first, firstCase)
	if err != nil { return nil, err }
	b, err := count(second, secondCase)
	if err != nil { return nil, err }
	result := make([]sharedAuthorizedWazuhWindow, 0)
	for key, firstCount := range a {
		secondCount := b[key]
		if secondCount == 0 { continue }
		result = append(result, sharedAuthorizedWazuhWindow{
			StartUTC: key.window, RuleID: key.rule, AgentID: key.agent,
			FirstEvents: firstCount, SecondEvents: secondCount,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.StartUTC != b.StartUTC { return a.StartUTC < b.StartUTC }
		if a.RuleID != b.RuleID { return a.RuleID < b.RuleID }
		return a.AgentID < b.AgentID
	})
	if len(result) > 10 { result = result[:10] }
	return result, nil
}
