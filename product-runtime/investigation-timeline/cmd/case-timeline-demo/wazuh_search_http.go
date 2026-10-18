package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// The imported Wazuh lab exposes a GET-only, bearer-guarded search endpoint.
// It reuses the existing Search Job, complete per-Case owner traversal and
// typed-source checks. Wazuh metadata is a filter, never an authorization.
func serveDemoWazuhSearch(w http.ResponseWriter, r *http.Request, token, principal, caseID string,
	data timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, owner http.Handler) {
	query := r.URL.Query()
	if len(query) != 5 || len(query["case_id"]) != 1 || query.Get("case_id") != caseID ||
		len(query["field"]) != 1 || len(query["value"]) != 1 ||
		len(query["from"]) != 1 || len(query["to"]) != 1 {
		http.NotFound(w, r)
		return
	}
	field := ""
	switch query.Get("field") {
	case "rule_id":
		field = "wazuh.rule.id"
	case "agent_id":
		field = "wazuh.agent.id"
	case "preview":
		field = "wazuh.preview.contains"
	default:
		http.NotFound(w, r)
		return
	}
	if (field == "wazuh.preview.contains" && !validWazuhPreviewLiteral(query.Get("value"))) ||
		(field != "wazuh.preview.contains" && !validWazuhCohortID(query.Get("value"))) {
		http.NotFound(w, r)
		return
	}
	from, startErr := time.Parse(time.RFC3339Nano, query.Get("from"))
	to, endErr := time.Parse(time.RFC3339Nano, query.Get("to"))
	if startErr != nil || endErr != nil || !from.Before(to) {
		http.NotFound(w, r)
		return
	}
	var tenant, environment string
	for _, p := range data.Principals {
		if p.ID == principal {
			tenant, environment = p.Tenant, p.Environment
			break
		}
	}
	if tenant == "" || environment == "" || imported < 1 || len(levels) != imported {
		http.NotFound(w, r)
		return
	}
	job, hits, err := executeDemoWazuhMetadataSearchWithHits(r.Context(), principal, caseID, owner, token,
		tenant, environment, levels, imported, field, query.Get("value"), query.Get("from"), query.Get("to"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	result := struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		Selector string `json:"selector"`
		Value string `json:"value"`
		Job any `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
		Bursts []wazuhSearchBurst `json:"search_bursts"`
	}{
		Format: "cmdr-local-wazuh4-executed-search-v1", Synthetic: true,
		CaseID: caseID, Selector: query.Get("field"), Value: query.Get("value"), Job: job,
		Hits: hits,
		Bursts: deriveWazuhSearchBursts(hits),
	}
	// Never send a partial response when authorization or job execution fails.
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(result); err != nil || r.Context().Err() != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(encoded.Bytes())
}
