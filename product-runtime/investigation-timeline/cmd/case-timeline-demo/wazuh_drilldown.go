package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

const wazuhDrilldownPageSize = 25

type wazuhDrilldownReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	Selector string `json:"selector"`
	Value string `json:"value"`
	SecondaryValue string `json:"secondary_value,omitempty"`
	TotalMatches int `json:"total_matches"`
	Offset int `json:"offset"`
	HasMore bool `json:"has_more"`
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
	return writeDemoWazuhDrilldownPage(ctx, principal, data, imported, raw, caseID, selector, value, 0, output)
}

func writeDemoWazuhDrilldownPage(ctx context.Context, principal string, data timeline.MemoryTimelineData, imported int, raw []byte, caseID, selector, value string, offset int, output io.Writer) error {
	return writeDemoWazuhDrilldownFilteredPage(ctx, principal, data, imported, raw, caseID, selector, value, "", offset, output)
}

// The intersection reuses exactly the same Case, Timeline and typed-source
// authorization traversal. Wazuh agent/rule labels never grant permissions.
func writeDemoWazuhDrilldownFilteredPage(ctx context.Context, principal string, data timeline.MemoryTimelineData, imported int, raw []byte, caseID, selector, value, secondValue string, offset int, output io.Writer) error {
	if (selector != "agent_id" && selector != "rule_id" && selector != "agent_rule" && selector != "technique_id") ||
		(selector == "technique_id" && !validWazuhTechniqueID(value)) ||
		(selector != "technique_id" && !validWazuhCohortID(value)) ||
		(selector == "agent_rule" && !validWazuhCohortID(secondValue)) ||
		(selector != "agent_rule" && secondValue != "") ||
		offset < 0 || offset > 5000 || offset%wazuhDrilldownPageSize != 0 {
		return errors.New("invalid Wazuh cohort selector")
	}
	levels, err := wazuhRuleLevels(raw, caseID)
	if err != nil { return err }
	token, err := randomBearer()
	if err != nil { return err }
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil { return err }
	result := wazuhDrilldownReceipt{
		Format: "cmdr-local-wazuh4-drilldown-v1", Synthetic: true,
		Selector: selector, Value: value, SecondaryValue: secondValue, Offset: offset,
		Entries: make([]wazuhDrilldownEntry, 0),
	}
	visibleImported := 0
	_, err = collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		metadata, exists := levels[entry.ID]
		if !exists { return nil }
		visibleImported++
		match := (selector == "agent_id" && metadata.AgentID == value) ||
			(selector == "rule_id" && metadata.RuleID == value) ||
			(selector == "agent_rule" && metadata.AgentID == value && metadata.RuleID == secondValue)
		if selector == "technique_id" {
			// Reuse the same source-level parser and already-authorized Timeline
			// entry; an ATT&CK tag is merely a filter, never an entitlement.
			for _, id := range metadata.MITREIDs {
				if id == value {
					match = true
					break
				}
			}
		}
		if !match { return nil }
		result.TotalMatches++
		if result.TotalMatches > offset && len(result.Entries) < wazuhDrilldownPageSize {
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
	result.HasMore = offset+len(result.Entries) < result.TotalMatches
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(result); err != nil { return err }
	_, err = output.Write(encoded.Bytes())
	return err
}

func serveDemoWazuhDrilldown(w http.ResponseWriter, r *http.Request, principal string, data timeline.MemoryTimelineData, imported int, raw []byte, caseID string) {
	// URL.Query() silently drops invalid query pairs (e.g. unknown=%ZZ).
	// Reject malformed syntax before trusting the remaining selection.
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if len(query["case_id"]) != 1 || query.Get("case_id") != caseID {
		http.NotFound(w, r)
		return
	}
	agentValues, hasAgent := query["agent_id"]
	ruleValues, hasRule := query["rule_id"]
	techniqueValues, hasTechnique := query["technique_id"]
	if (!hasAgent && !hasRule && !hasTechnique) ||
		(hasTechnique && (hasAgent || hasRule || len(techniqueValues) != 1 ||
			!validWazuhTechniqueID(techniqueValues[0]))) ||
		(hasAgent && (len(agentValues) != 1 || !validWazuhCohortID(agentValues[0]))) ||
		(hasRule && (len(ruleValues) != 1 || !validWazuhCohortID(ruleValues[0]))) {
		http.NotFound(w, r)
		return
	}
	selector, value, secondValue := "", "", ""
	if hasTechnique {
		selector, value = "technique_id", techniqueValues[0]
	}
	if hasAgent {
		selector, value = "agent_id", agentValues[0]
	}
	if hasRule {
		if hasAgent {
			selector, secondValue = "agent_rule", ruleValues[0]
		} else {
			selector, value = "rule_id", ruleValues[0]
		}
	}
	expectedKeys := 2
	if hasAgent && hasRule {
		expectedKeys++
	}
	if _, hasOffset := query["offset"]; hasOffset {
		expectedKeys++
	}
	// Reject unrecognized query parameters; they can never affect access.
	if len(query) != expectedKeys {
		http.NotFound(w, r)
		return
	}
	offset := 0
	if values, ok := query["offset"]; ok {
		if len(values) != 1 || len(values[0]) == 0 || len(values[0]) > 4 {
			http.NotFound(w, r)
			return
		}
		for _, char := range values[0] {
			if char < '0' || char > '9' {
				http.NotFound(w, r)
				return
			}
		}
		parsed, err := strconv.Atoi(values[0])
		if err != nil || parsed < 0 || parsed > 5000 || parsed%wazuhDrilldownPageSize != 0 {
			http.NotFound(w, r)
			return
		}
		offset = parsed
	}
	var encoded bytes.Buffer
	if err := writeDemoWazuhDrilldownFilteredPage(r.Context(), principal, data, imported, raw, caseID, selector, value, secondValue, offset, &encoded); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(encoded.Bytes())
}
