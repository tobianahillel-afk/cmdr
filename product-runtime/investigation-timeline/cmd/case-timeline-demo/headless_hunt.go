package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

// Execute the SAME authenticated first-party Event Search Job used by the
// analyst screen, but locally without a browser, port, or exported bearer.
// Each matched typed source is independently authorized BEFORE the query
// predicate; correlation uses only references from this completed Job.
func runSavedCaseHeadlessHunt(ctx context.Context, principal, snapshotPath, caseID,
	expectedDigest, field, family, value, from, to string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || principal == "" || snapshotPath == "" ||
		caseID == "" || expectedDigest == "" || output == nil {
		return errors.New("headless Case hunt requires selected principal, Case and operator-held snapshot SHA-256")
	}
	data, imported, err := loadAuthorizedImportedCase(ctx, principal, snapshotPath, caseID, expectedDigest)
	if err != nil { return err }
	token, err := randomBearer()
	if err != nil { return err }
	owner, err := newDemoHandlerWithData(token, principal, data)
	if err != nil { return err }
	handler := newDemoLocalJSONLSearchHandler(token, principal, caseID, data, imported, owner)
	selection, err := json.Marshal(struct {
		CaseID string `json:"case_id"`
		Field string `json:"field"`
		Family string `json:"source_family"`
		Value string `json:"value"`
		From string `json:"from"`
		To string `json:"to"`
	}{caseID, field, family, value, from, to})
	if err != nil || len(selection) > 4096 {
		return errors.New("headless search selection exceeds the existing bounded search contract")
	}
	req := httptest.NewRequest(http.MethodPost, "/demo-local-search", bytes.NewReader(selection)).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if ctx.Err() != nil || response.Code != http.StatusOK ||
		response.Body.Len() > maxDemoSnapshotBytes+4096 {
		return errors.New("authorized headless Event Search Job unavailable")
	}
	var completed struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		Job search.SearchJob `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
		Bursts []wazuhSearchBurst `json:"search_bursts"`
		MixedContext []mixedTemporalLead `json:"mixed_temporal_context"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &completed); err != nil ||
		completed.Format != "cmdr-local-jsonl-event-search-v1" ||
		!completed.Synthetic || completed.CaseID != caseID ||
		completed.Job.State != search.JobCompleted ||
		len(completed.Job.ResultRefs) > maxWazuhSearchMatches {
		return errors.New("headless Event Search receipt incomplete or invalid")
	}
	// Never echo the operator's literal search value, bearer, raw preview,
	// actor/tenant audit internals, or any unrelated owner records. Keep
	// actionable refs and already validated Wazuh/temporal leads only.
	var receipt bytes.Buffer
	err = json.NewEncoder(&receipt).Encode(struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		JobID string `json:"job_id"`
		SourceFamily string `json:"source_family"`
		QueryField string `json:"query_field"`
		Refs []search.ResultReference `json:"authorized_refs"`
		Hits []wazuhSearchHit `json:"wazuh_hits,omitempty"`
		Bursts []wazuhSearchBurst `json:"wazuh_bursts,omitempty"`
		MixedContext []mixedTemporalLead `json:"mixed_temporal_context,omitempty"`
	}{
		Format:"cmdr-local-headless-authorized-hunt-v1", Synthetic:true,
		CaseID:caseID, JobID:completed.Job.ID, SourceFamily:family,
		QueryField:field, Refs:completed.Job.ResultRefs,
		Hits:completed.Hits, Bursts:completed.Bursts,
		MixedContext:completed.MixedContext,
	})
	if err != nil || ctx.Err() != nil || strings.Contains(receipt.String(), token) {
		return errors.New("headless Event Search output unavailable")
	}
	// io.Copy reports io.ErrShortWrite when the destination silently
	// accepts only part of an authorized forensic receipt.
	_, err = io.Copy(output, &receipt)
	return err
}
