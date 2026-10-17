package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// The drill-down returns only server-authorized Timeline identities and basic
// time/kind fields. Source previews are obtained separately through the
// existing Case/Entry/typed-source navigation authorization chain.
type wazuhDrilldownEntry struct {
	ID string `json:"id"`
	Kind string `json:"kind"`
	SourceTime string `json:"source_time"`
	TimeUncertain bool `json:"time_uncertain"`
}

type wazuhDrilldownReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	Selector string `json:"selector"`
	Value string `json:"value"`
	SecondaryValue string `json:"secondary_value,omitempty"`
	TotalMatches int `json:"total_matches"`
	Entries []wazuhDrilldownEntry `json:"entries"`
}

func validWazuhCohortID(value string) bool {
	if len(value) < 1 || len(value) > 16 {
		return false
	}
	for _, c := range value {
		if c < '0' || c > '9' { return false }
	}
	return true
}

// Every request replays the established authorized Case/Timeline/source path,
// including the complete imported Wazuh set. Nothing is written to the
// response until the entire permission traversal succeeds.
func writeDemoWazuhDrilldown(ctx context.Context, principal string, data timeline.MemoryTimelineData, imported int, raw []byte, caseID, selector, value string, output io.Writer) error {
	if (selector != "agent_id" && selector != "rule_id") || !validWazuhCohortID(value) {
		return errors.New("invalid Wazuh cohort selector")
	}
	levels, err := wazuhRuleLevels(raw, caseID)
	if err != nil { return err }
	return writeDemoWazuhDrilldownWithLevels(ctx, principal, data, imported, levels, caseID, selector, value, output)
}

// Wazuh metadata is joined ONLY after the trusted owner traversal. The
// caller supplies feed-qualified IDs from the exact validated import bytes.
func writeDemoWazuhDrilldownWithLevels(ctx context.Context, principal string, data timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, caseID, selector, value string, output io.Writer) error {
	return writeDemoWazuhDrilldownWithLevelsFilter(ctx, principal, data, imported, levels, caseID, selector, value, "", output)
}

// Cross-feed signals are filters over already-authorized sources, not access rights.
func writeDemoWazuhDrilldownWithLevelsFilter(ctx context.Context, principal string, data timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, caseID, selector, value, secondValue string, output io.Writer) error {
	if (selector != "agent_id" && selector != "rule_id" && selector != "agent_rule") ||
		!validWazuhCohortID(value) ||
		(selector == "agent_rule" && !validWazuhCohortID(secondValue)) ||
		(selector != "agent_rule" && secondValue != "") {
		return errors.New("invalid Wazuh cohort selector")
	}
	token, err := randomBearer()
	if err != nil { return err }
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil { return err }
	result := wazuhDrilldownReceipt{
		Format: "cmdr-local-wazuh4-drilldown-v1", Synthetic: true,
		Selector: selector, Value: value, SecondaryValue: secondValue, Entries: make([]wazuhDrilldownEntry, 0),
	}
	visibleImported := 0
	_, err = collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		metadata, exists := levels[entry.ID]
		if !exists { return nil }
		visibleImported++
		match := (selector == "agent_id" && metadata.AgentID == value) ||
			(selector == "rule_id" && metadata.RuleID == value) ||
			(selector == "agent_rule" && metadata.AgentID == value && metadata.RuleID == secondValue)
		if !match { return nil }
		result.TotalMatches++
		if len(result.Entries) < 25 {
			result.Entries = append(result.Entries, wazuhDrilldownEntry{
				ID: entry.ID, Kind: entry.Kind, SourceTime: entry.SourceTime, TimeUncertain: entry.TimeUncertain,
			})
		}
		return nil
	})
	if err != nil { return err }
	if visibleImported != imported || visibleImported != len(levels) {
		return errors.New("incomplete authorized Wazuh traversal")
	}
	if err := ctx.Err(); err != nil { return err }
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(result); err != nil { return err }
	_, err = output.Write(encoded.Bytes())
	return err
}

// Existing single-feed Wazuh mode shares the same strict HTTP selector
// and typed-source authorization path as the multi-feed mode.
func serveDemoWazuhDrilldown(w http.ResponseWriter, r *http.Request, principal string, data timeline.MemoryTimelineData, imported int, raw []byte, caseID string) {
	levels, err := wazuhRuleLevels(raw, caseID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	serveDemoWazuhDrilldownWithLevels(w, r, principal, data, imported, levels, caseID)
}

func serveDemoWazuhDrilldownWithLevels(w http.ResponseWriter, r *http.Request, principal string, data timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, caseID string) {
	query := r.URL.Query()
	if len(query["case_id"]) != 1 || query.Get("case_id") != caseID {
		http.NotFound(w, r)
		return
	}
	agents, hasAgent := query["agent_id"]
	rules, hasRule := query["rule_id"]
	expectedKeys := 1
	if hasAgent { expectedKeys++ }
	if hasRule { expectedKeys++ }
	if (!hasAgent && !hasRule) || len(query) != expectedKeys ||
		(hasAgent && (len(agents) != 1 || !validWazuhCohortID(agents[0]))) ||
		(hasRule && (len(rules) != 1 || !validWazuhCohortID(rules[0]))) {
		http.NotFound(w, r)
		return
	}
	selector, value, secondValue := "", "", ""
	if hasAgent {
		selector, value = "agent_id", agents[0]
	}
	if hasRule {
		if hasAgent {
			selector, secondValue = "agent_rule", rules[0]
		} else {
			selector, value = "rule_id", rules[0]
		}
	}
	var encoded bytes.Buffer
	if err := writeDemoWazuhDrilldownWithLevelsFilter(r.Context(), principal, data, imported, levels, caseID, selector, value, secondValue, &encoded); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(encoded.Bytes())
}
