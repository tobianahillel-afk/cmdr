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

func TestWazuhDrilldownReusesAuthorizedImportedTimeline(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"real-1\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"PRIVATE-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-PREVIEW\"}}\n"+
		"{\"id\":\"real-2\",\"timestamp\":\"2026-10-01T10:00:01Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-PREVIEW\"}}\n")
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(), "analyst-a", owner, "case-a", alerts)
	if err != nil { t.Fatal(err) }
	token := strings.Repeat("a", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoWazuhCohortHandler(token, "analyst-a", "case-a", data, imported, raw, base)
	request := func(method, path, bearer string, ctx context.Context) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, nil).WithContext(ctx)
		if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		return out
	}
	for _, tc := range []struct { path string; count int }{
		{"/demo-wazuh-drilldown?case_id=case-a&agent_id=001", 1},
		{"/demo-wazuh-drilldown?case_id=case-a&rule_id=5710", 2},
		{"/demo-wazuh-drilldown?case_id=case-a&agent_id=003", 0},
	} {
		got := request(http.MethodGet, tc.path, token, context.Background())
		if got.Code != 200 || got.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("drilldown rejected: %s %d %s", tc.path, got.Code, got.Body.String())
		}
		var receipt wazuhDrilldownReceipt
		if err := json.Unmarshal(got.Body.Bytes(), &receipt); err != nil { t.Fatal(err) }
		if receipt.Format != "cmdr-local-wazuh4-drilldown-v1" || !receipt.Synthetic ||
			receipt.TotalMatches != tc.count || len(receipt.Entries) != tc.count {
			t.Fatalf("unexpected Wazuh drilldown %s: %+v", tc.path, receipt)
		}
		for _, protected := range []string{"PRIVATE-HOST", "PRIVATE-PREVIEW", "real-1", "real-2", "demo-tenant-b"} {
			if bytes.Contains(got.Body.Bytes(), []byte(protected)) { t.Fatalf("raw or foreign data leaked: %s", protected) }
		}
	}
	denials := []struct { method,path,bearer string; want int }{
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001","",401},
		{"GET","/demo-wazuh-drilldown?case_id=case-b&agent_id=001",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&agent_id=002",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=%2A",token,404},
		{"POST","/demo-wazuh-drilldown?case_id=case-a&agent_id=001",token,405},
	}
	for _, tc := range denials {
		got := request(tc.method,tc.path,tc.bearer,context.Background())
		if got.Code != tc.want || bytes.Contains(got.Body.Bytes(),[]byte("local-wazuh4-drilldown")) {
			t.Fatalf("request leaked data: %s %d %s",tc.path,got.Code,got.Body.String())
		}
	}
	ctx,cancel := context.WithCancel(context.Background())
	cancel()
	got := request("GET","/demo-wazuh-drilldown?case_id=case-a&rule_id=5710",token,ctx)
	if got.Code != 404 || bytes.Contains(got.Body.Bytes(),[]byte("local-wazuh4-drilldown")) {
		t.Fatalf("cancelled response disclosed data: %d %s",got.Code,got.Body.String())
	}
}

func TestWazuhDrilldownUnavailableOutsideWazuhSession(t *testing.T) {
	handler,err:=newDemoHandlerWithData(strings.Repeat("a",64),"analyst-a",demoData())
	if err!=nil {t.Fatal(err)}
	req:=httptest.NewRequest("GET","/demo-wazuh-drilldown?case_id=case-a&rule_id=5710",nil)
	req.Header.Set("Authorization","Bearer "+strings.Repeat("a",64))
	out:=httptest.NewRecorder();handler.ServeHTTP(out,req)
	if out.Code!=404 {t.Fatalf("ordinary snapshot exposed Wazuh drilldown: %d",out.Code)}
}
