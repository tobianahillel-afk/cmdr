package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	if (selector != "agent_id" && selector != "rule_id") || !validWazuhCohortID(value) ||
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
		Selector: selector, Value: value, Offset: offset,
		Entries: make([]wazuhDrilldownEntry, 0),
	}
	visibleImported := 0
	_, err = collectDemoCaseTriageWithObserver(ctx, handler, token, caseID, func(entry timeline.ProjectedEntry) error {
		metadata, exists := levels[entry.ID]
		if !exists { return nil }
		visibleImported++
		match := (selector == "agent_id" && metadata.AgentID == value) ||
			(selector == "rule_id" && metadata.RuleID == value)
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
	query := r.URL.Query()
	if (len(query) != 2 && len(query) != 3) || len(query["case_id"]) != 1 || query.Get("case_id") != caseID {
		http.NotFound(w, r)
		return
	}
	selector, value := "", ""
	for _, key := range []string{"agent_id", "rule_id"} {
		if ids, ok := query[key]; ok {
			if len(ids) != 1 || selector != "" {
				http.NotFound(w, r)
				return
			}
			selector, value = key, ids[0]
		}
	}
	if selector == "" || !validWazuhCohortID(value) {
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
	// Unknown selectors are refused, even if the correct Case and cohort were supplied.
	if len(query) == 3 && len(query["offset"]) != 1 {
		http.NotFound(w, r)
		return
	}
	var encoded bytes.Buffer
	if err := writeDemoWazuhDrilldownPage(r.Context(), principal, data, imported, raw, caseID, selector, value, offset, &encoded); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(encoded.Bytes())
}
