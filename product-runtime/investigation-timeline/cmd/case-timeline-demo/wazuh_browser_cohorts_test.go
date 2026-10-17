package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWazuhBrowserCohortsUseAuthenticatedOwnerAndTypedSources(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"real-1\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"PRIVATE-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-PREVIEW\"}}\n"+
			"{\"id\":\"real-2\",\"timestamp\":\"2026-10-01T10:00:01Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-PREVIEW\"}}\n")
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(), "analyst-a", owner, "case-a", alerts)
	if err != nil {
		t.Fatal(err)
	}
	token := strings.Repeat("a", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	handler := newDemoWazuhCohortHandler(token, "analyst-a", "case-a", data, imported, raw, base)
	fetch := func(method, path, bearer string, ctx context.Context) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(method, path, nil).WithContext(ctx)
		if bearer != "" {
			request.Header.Set("Authorization", "Bearer "+bearer)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, request)
		return w
	}
	good := fetch(http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a", token, context.Background())
	if good.Code != http.StatusOK || good.Header().Get("Cache-Control") != "no-store" ||
		!strings.HasPrefix(good.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("authorized cohort request failed: %d %q", good.Code, good.Body.String())
	}
	var receipt wazuhRuleTriageReceipt
	if err := json.Unmarshal(good.Body.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.ImportedEntries != 2 || receipt.Case.VisibleEntries != 209 ||
		len(receipt.TopAgents) != 1 || receipt.TopAgents[0].AgentID != "001" ||
		receipt.TopAgents[0].Count != 2 || receipt.TopAgents[0].HighestKnownRuleLevel == nil ||
		*receipt.TopAgents[0].HighestKnownRuleLevel != 12 ||
		len(receipt.TopRules) != 1 || receipt.TopRules[0].RuleID != "5710" ||
		receipt.TopRules[0].Count != 2 {
		t.Fatalf("cohort calculation or owner scope mismatch: %+v", receipt)
	}
	for _, protected := range []string{"PRIVATE-HOST", "PRIVATE-PREVIEW", "real-1", "demo-tenant-b", "Bearer "} {
		if bytes.Contains(good.Body.Bytes(), []byte(protected)) {
			t.Fatalf("Wazuh browser disclosed unauthorized/raw material %q", protected)
		}
	}
	for _, candidate := range []struct {
		method, path, bearer string
		want int
	}{
		{http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a", "", http.StatusUnauthorized},
		{http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a", strings.Repeat("b", 64), http.StatusUnauthorized},
		{http.MethodGet, "/demo-wazuh-cohorts?case_id=case-b", token, http.StatusNotFound},
		{http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a&case_id=case-a", token, http.StatusNotFound},
		{http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a&tenant=tenant-a", token, http.StatusNotFound},
		{http.MethodPost, "/demo-wazuh-cohorts?case_id=case-a", token, http.StatusMethodNotAllowed},
	} {
		result := fetch(candidate.method, candidate.path, candidate.bearer, context.Background())
		if result.Code != candidate.want ||
			strings.Contains(result.Body.String(), "top_agents") || strings.Contains(result.Body.String(), "5710") {
			t.Fatalf("bad request disclosed cohort: %q %d %q", candidate.path, result.Code, result.Body.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	denied := fetch(http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a", token, ctx)
	if denied.Code != http.StatusNotFound || strings.Contains(denied.Body.String(), "5710") {
		t.Fatalf("cancelled request leaked data: %d %s", denied.Code, denied.Body.String())
	}
}

func TestWazuhBrowserDoesNotInstallCohortsForOrdinarySnapshot(t *testing.T) {
	data := demoData()
	token := strings.Repeat("f", 64)
	handler, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/demo-wazuh-cohorts?case_id=case-a", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	out := httptest.NewRecorder()
	handler.ServeHTTP(out, req)
	if out.Code != http.StatusNotFound || strings.Contains(out.Body.String(), "5710") {
		t.Fatalf("Wazuh route was installed without an authorized import: %d %s", out.Code, out.Body.String())
	}
	// No changed data authority: the existing trusted memory owner remains
	// the sole Case/Timeline and typed-source reader.
}
