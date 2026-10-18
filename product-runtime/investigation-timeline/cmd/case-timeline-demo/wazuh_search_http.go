package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// The Wazuh lab uses GET only for numeric pivots. Preview text uses POST JSON
// so sensitive investigation phrases never enter query strings or URL logs.
// It reuses the existing Search Job, complete per-Case owner traversal and
// typed-source checks. Wazuh metadata is a filter, never an authorization.
func serveDemoWazuhSearch(w http.ResponseWriter, r *http.Request, token, principal, caseID string,
	data timeline.MemoryTimelineData, imported int, levels map[string]wazuhRuleLevel, owner http.Handler) {
	// Parse precisely one existing Wazuh Search Job input. POST is used
	// only for textual source hunts; numeric legacy GET pivots stay stable.
	var selection struct {
		CaseID string `json:"case_id"`
		Field string `json:"field"`
		Value string `json:"value"`
		From string `json:"from"`
		To string `json:"to"`
	}
	if r.Method == http.MethodPost {
		if r.URL.RawQuery != "" || r.Header.Get("Content-Type") != "application/json" {
			http.NotFound(w, r)
			return
		}
		// Strict global 4KiB JSON body ceiling; even trailing whitespace is
		// bounded. Unknown fields/trailing objects are refused.
		body, err := io.ReadAll(io.LimitReader(r.Body, 4097))
		if err != nil || len(body) == 0 || len(body) > 4096 {
			http.NotFound(w, r)
			return
		}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&selection); err != nil {
			http.NotFound(w, r)
			return
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			http.NotFound(w, r)
			return
		}
	} else {
		query := r.URL.Query()
		if len(query) != 5 || len(query["case_id"]) != 1 ||
			len(query["field"]) != 1 || len(query["value"]) != 1 ||
			len(query["from"]) != 1 || len(query["to"]) != 1 {
			http.NotFound(w, r)
			return
		}
		selection.CaseID, selection.Field, selection.Value =
			query.Get("case_id"), query.Get("field"), query.Get("value")
		selection.From, selection.To = query.Get("from"), query.Get("to")
	}
	if selection.CaseID != caseID {
		http.NotFound(w, r)
		return
	}
	field := ""
	switch selection.Field {
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
	if field == "wazuh.preview.contains" {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method unavailable", http.StatusMethodNotAllowed)
			return
		}
		if !validWazuhPreviewLiteral(selection.Value) {
			http.NotFound(w, r)
			return
		}
	} else {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method unavailable", http.StatusMethodNotAllowed)
			return
		}
		if !validWazuhCohortID(selection.Value) {
			http.NotFound(w, r)
			return
		}
	}
	from, startErr := time.Parse(time.RFC3339Nano, selection.From)
	to, endErr := time.Parse(time.RFC3339Nano, selection.To)
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
		tenant, environment, levels, imported, field, selection.Value, selection.From, selection.To)
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
		CaseID: caseID, Selector: selection.Field, Job: job,
		Hits: hits,
		Bursts: deriveWazuhSearchBursts(hits),
	}
	// Do NOT echo potentially sensitive hunt text through JSON receipts,
	// reports or locally saved replay files. Numeric pivots preserve their
	// existing metadata receipt contract.
	if field != "wazuh.preview.contains" {
		result.Value = selection.Value
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
