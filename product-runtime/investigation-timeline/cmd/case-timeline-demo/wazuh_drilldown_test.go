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
		{"/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710", 1},
		{"/demo-wazuh-drilldown?case_id=case-a&agent_id=002&rule_id=5710", 1},
		{"/demo-wazuh-drilldown?case_id=case-a&agent_id=003&rule_id=5710", 0},
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
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&tenant=tenant-a",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&rule_id=100",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=%2A",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-b&agent_id=001&rule_id=5710",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&agent_id=002",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&agent_id=002",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&offset=1",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&offset=25&offset=50",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&unknown=1",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=",token,404},
		{"GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=%2A",token,404},
		{"POST","/demo-wazuh-drilldown?case_id=case-a&agent_id=001",token,405},
	}
	for _, tc := range denials {
		got := request(tc.method,tc.path,tc.bearer,context.Background())
		if got.Code != tc.want || bytes.Contains(got.Body.Bytes(),[]byte("local-wazuh4-drilldown")) {
			t.Fatalf("request leaked data: %s %d %s",tc.path,got.Code,got.Body.String())
		}
	}
	// The browser controls must never be an authority. A direct in-process
	// caller providing invalid combined filters is also rejected atomically.
	for _, bad := range []struct {
		selector, value, secondary string
		offset int
	}{
		{"agent_rule", "001", "", 0},
		{"agent_rule", "001", "*", 0},
		{"agent_rule", "001", "5710", 1},
		{"agent_rule", "001", "5710", -25},
		{"agent_id", "001", "5710", 0},
		{"rule_id", "5710", "001", 0},
		{"unknown", "001", "", 0},
		{"agent_rule", "*", "5710", 0},
	} {
		var output bytes.Buffer
		err := writeDemoWazuhDrilldownFilteredPage(context.Background(), "analyst-a",
			data, imported, raw, "case-a", bad.selector, bad.value, bad.secondary, bad.offset, &output)
		if err == nil || output.Len() != 0 {
			t.Fatalf("invalid combined selector accepted: %+v, output %s", bad, output.String())
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

func TestWazuhCombinedDrilldownUsesAuthorizedPagedIntersection(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var alertsBuilder strings.Builder
	for i := 0; i < 61; i++ {
		agent, rule := "001", "5710"
		if i%3 == 0 { agent = "002" }
		if i%4 == 0 { rule = "1001" }
		fmt.Fprintf(&alertsBuilder, "{\"id\":\"pair-%d\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"%s\"},\"rule\":{\"id\":\"%s\",\"level\":9,\"description\":\"PRIVATE-DATA\"}}\n", i, agent, rule)
	}
	alerts := writeJSONLEvents(t, alertsBuilder.String())
	data, imported, raw, err := prepareWazuh4CaseWithRaw(context.Background(), "analyst-a", owner, "case-a", alerts)
	if err != nil { t.Fatal(err) }
	token := strings.Repeat("a", 64)
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoWazuhCohortHandler(token, "analyst-a", "case-a", data, imported, raw, base)
	matched := 0
	for i := 0; i < 61; i++ { if i%3 != 0 && i%4 != 0 { matched++ } }
	seen := map[string]bool{}
	for offset := 0; offset < matched; offset += wazuhDrilldownPageSize {
		path := fmt.Sprintf("/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710&offset=%d",offset)
		req := httptest.NewRequest("GET", path, nil); req.Header.Set("Authorization", "Bearer "+token)
		out := httptest.NewRecorder();handler.ServeHTTP(out, req)
		if out.Code != 200 { t.Fatalf("combined filter rejected: %d %s",out.Code,out.Body.String()) }
		var result wazuhDrilldownReceipt
		if err := json.Unmarshal(out.Body.Bytes(), &result);err != nil {t.Fatal(err)}
		if result.Selector != "agent_rule" || result.Value != "001" || result.SecondaryValue != "5710" ||
			result.TotalMatches != matched || result.Offset != offset ||
			result.HasMore != (offset+len(result.Entries)<matched) {
			t.Fatalf("bad intersection page: %+v",result)
		}
		for _, e := range result.Entries {
			if seen[e.ID] { t.Fatal("repeated event identity") }
			seen[e.ID] = true
		}
		if strings.Contains(out.Body.String(),"PRIVATE-DATA") {t.Fatal("source preview leaked in combined result")}
	}
	if len(seen) != matched {t.Fatalf("incomplete intersected result: %d want %d",len(seen),matched)}
	req := httptest.NewRequest("GET","/demo-wazuh-drilldown?case_id=case-a&agent_id=001&rule_id=5710",nil)
	out := httptest.NewRecorder();handler.ServeHTTP(out,req)
	if out.Code != 401 || strings.Contains(out.Body.String(),"agent_rule") {t.Fatal("unauthenticated intersection disclosed results")}
}
