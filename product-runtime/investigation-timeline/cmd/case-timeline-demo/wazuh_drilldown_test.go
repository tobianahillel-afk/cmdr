package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func TestWazuhDrilldownPaginatesAuthorizedEvents(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var alertLines strings.Builder
	for i := 0; i < 53; i++ {
		fmt.Fprintf(&alertLines, "{\"id\":\"page-%d\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-DATA\"}}\n", i)
	}
	alerts := writeJSONLEvents(t, alertLines.String())
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(), "analyst-a", owner, "case-a", alerts)
	if err != nil { t.Fatal(err) }
	token := strings.Repeat("a", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoWazuhCohortHandler(token, "analyst-a", "case-a", data, imported, raw, base)
	seen := make(map[string]bool)
	for _, page := range []struct{ offset, count int; more bool }{
		{0,25,true},{25,25,true},{50,3,false},
	} {
		path := fmt.Sprintf("/demo-wazuh-drilldown?case_id=case-a&agent_id=001&offset=%d", page.offset)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization","Bearer "+token)
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		if out.Code != http.StatusOK || out.Header().Get("Cache-Control") != "no-store" { t.Fatalf("paged request rejected: %d %s", out.Code, out.Body.String()) }
		var receipt wazuhDrilldownReceipt
		if err := json.Unmarshal(out.Body.Bytes(), &receipt); err != nil { t.Fatal(err) }
		if receipt.Offset != page.offset || receipt.TotalMatches != 53 || receipt.HasMore != page.more || len(receipt.Entries) != page.count {
			t.Fatalf("wrong authorized page: %+v", receipt)
		}
		if strings.Contains(out.Body.String(),"PRIVATE-DATA") || strings.Contains(out.Body.String(),"Bearer ") { t.Fatal("raw source or credential leaked") }
		for _, entry := range receipt.Entries {
			if seen[entry.ID] { t.Fatalf("duplicate entry across pages: %s", entry.ID) }
			seen[entry.ID]=true
		}
	}
	if len(seen) != 53 { t.Fatalf("missing authorized matching rows: %d",len(seen)) }
	for _, query := range []string{
		"case_id=case-a&rule_id=5710&offset=-25",
		"case_id=case-a&rule_id=5710&offset=1",
		"case_id=case-a&rule_id=5710&offset=25&offset=50",
		"case_id=case-a&rule_id=5710&offset=999999999",
		"case_id=case-a&rule_id=5710&offset=0&tenant=tenant-a",
		"case_id=case-b&rule_id=5710&offset=25",
	} {
		req := httptest.NewRequest("GET","/demo-wazuh-drilldown?"+query,nil)
		req.Header.Set("Authorization","Bearer "+token)
		out := httptest.NewRecorder()
		handler.ServeHTTP(out,req)
		if out.Code != http.StatusNotFound || strings.Contains(out.Body.String(),"page-") {
			t.Fatalf("invalid page leaked: %q status=%d",query,out.Code)
		}
	}
	ctx,cancel := context.WithCancel(context.Background());cancel()
	req := httptest.NewRequest("GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&offset=25",nil).WithContext(ctx)
	req.Header.Set("Authorization","Bearer "+token)
	out:=httptest.NewRecorder();handler.ServeHTTP(out,req)
	if out.Code != http.StatusNotFound || strings.Contains(out.Body.String(),"entries") { t.Fatalf("cancelled page leaked: %s",out.Body.String()) }
}
