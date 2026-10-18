package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

type wazuhHTTPJobReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	CaseID string `json:"case_id"`
	Selector string `json:"selector"`
	Value string `json:"value"`
	Job search.SearchJob `json:"job"`
	Hits []wazuhSearchHit `json:"search_hits"`
}

// Real disk imports, independent case-owner HTTP authentication and a full
// first-party Search Job execute inside the same analyst lab request.
func TestWazuhAnalystHTTPExecutesSearchAndReopensSources(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t,
		"{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"private-first-feed\"}}\n"+
		"{\"id\":\"late.1\",\"timestamp\":\"2026-10-17T10:05:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"private-late-feed\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":6,\"description\":\"private-second-feed\"}}\n")
	data, imported, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(),
		"analyst-a", owner, caseID, []string{first, second})
	if err != nil { t.Fatal(err) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoWazuhCohortReceiptHandler(token, "analyst-a", caseID, data, imported, nil,
		levels, base, func(context.Context, io.Writer) error { return nil })
	query := url.Values{
		"case_id": {caseID}, "field": {"rule_id"}, "value": {"5710"},
		"from": {"2026-10-17T09:59:00Z"}, "to": {"2026-10-17T10:03:00Z"},
	}
	invoke := func(t *testing.T, q url.Values, bearer string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, "/demo-wazuh-search?"+q.Encode(), nil)
		if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		return resp
	}
	ok := invoke(t, query, token)
	if ok.Code != http.StatusOK { t.Fatalf("authorized search returned %d: %s", ok.Code, ok.Body.String()) }
	var receipt wazuhHTTPJobReceipt
	if err := json.Unmarshal(ok.Body.Bytes(), &receipt); err != nil { t.Fatal(err) }
	if receipt.Format != "cmdr-local-wazuh4-executed-search-v1" ||
		!receipt.Synthetic || receipt.CaseID != caseID ||
		receipt.Job.State != search.JobCompleted ||
		len(receipt.Job.ResultRefs) != 2 || len(receipt.Job.AuditTrail) < 3 ||
		receipt.Job.Provenance.TimeStart.IsZero() ||
		receipt.Job.Provenance.TimeEnd.IsZero() {
		t.Fatalf("search result has no authoritative completed job: %+v", receipt)
	}
	// The same authorized execution now yields a source-time-ordered triage
	// view with metadata from the exact bytes of the imported Wazuh feeds.
	if len(receipt.Hits) != len(receipt.Job.ResultRefs) ||
		receipt.Hits[0].EntryID != importedWazuhID(caseID, "same.1", 1) ||
		receipt.Hits[1].EntryID != importedWazuhID(caseID, "same.1", 2) ||
		receipt.Hits[0].SourceTimeUTC != "2026-10-17T10:00:00Z" ||
		receipt.Hits[1].SourceTimeUTC != "2026-10-17T10:01:00Z" ||
		receipt.Hits[0].RuleID != "5710" || receipt.Hits[0].AgentID != "001" ||
		receipt.Hits[0].Feed != 1 || receipt.Hits[1].Feed != 2 ||
		receipt.Hits[0].RuleLevel == nil || *receipt.Hits[0].RuleLevel != 12 ||
		receipt.Hits[1].RuleLevel == nil || *receipt.Hits[1].RuleLevel != 6 {
		t.Fatalf("authorized chronological Wazuh triage context mismatches completed SearchJob: %+v", receipt.Hits)
	}
	matched := map[string]bool{}
	for _, result := range receipt.Job.ResultRefs {
		if result.Source != wazuhLocalSearchSource { t.Fatalf("untrusted source: %+v", result) }
		matched[result.Ref] = true
	}
	if !matched[importedWazuhID(caseID, "same.1", 1)] ||
		!matched[importedWazuhID(caseID, "same.1", 2)] ||
		matched[importedWazuhID(caseID, "late.1", 1)] {
		t.Fatalf("wrong authorized two-feed [from,to) match: %+v", matched)
	}
	for _, forbidden := range []string{"private-first-feed", "private-late-feed",
		"private-second-feed", "demo-tenant-b", token} {
		if strings.Contains(ok.Body.String(), forbidden) {
			t.Fatalf("sensitive payload in Search Job receipt: %s", forbidden)
		}
	}
	// Every returned entry remains independently navigable with the same
	// existing typed-source owner; the Search Job contains no raw preview.
	for ref := range matched {
		source := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet,
			"/demo-source-open?"+url.Values{"case_id": {caseID}, "entry_id": {ref}}.Encode(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		base.ServeHTTP(source, req)
		if source.Code != http.StatusOK { t.Fatalf("search hit not navigable: %s: %d", ref, source.Code) }
	}

	// The same browser route can search an agent; no new Case or owner exists.
	agent := url.Values{}
	for k, v := range query { agent[k] = append([]string(nil), v...) }
	agent.Set("field", "agent_id")
	agent.Set("value", "001")
	narrow := invoke(t, agent, token)
	var selected wazuhHTTPJobReceipt
	if narrow.Code != http.StatusOK || json.Unmarshal(narrow.Body.Bytes(), &selected) != nil ||
		len(selected.Job.ResultRefs) != 1 || len(selected.Hits) != 1 ||
		selected.Hits[0].AgentID != "001" ||
		selected.Job.ResultRefs[0].Ref != importedWazuhID(caseID, "same.1", 1) {
		t.Fatalf("agent-only Search Job unexpected: %d %s", narrow.Code, narrow.Body.String())
	}
	// An authorized no-match is a complete, empty job, not an invented finding.
	agent.Set("value", "999")
	empty := invoke(t, agent, token)
	var noMatch wazuhHTTPJobReceipt
	if empty.Code != http.StatusOK || json.Unmarshal(empty.Body.Bytes(), &noMatch) != nil ||
		noMatch.Job.State != search.JobCompleted || len(noMatch.Job.ResultRefs) != 0 ||
		len(noMatch.Hits) != 0 {
		t.Fatalf("no-match query did not complete safely: %d %s", empty.Code, empty.Body.String())
	}

	for _, tc := range []struct{
		name, key, value, bearer string
	}{
		{"no bearer", "", "", ""},
		{"wrong bearer", "", "", "bad-token"},
		{"foreign Case", "case_id", "case-b", token},
		{"unrecognized field", "field", "hostname", token},
		{"bad identifier", "value", "1%20OR%201", token},
		{"invalid time", "from", "not-a-time", token},
		{"reversed bounds", "from", "2026-10-17T11:00:00Z", token},
		{"unknown parameter", "anything", "1", token},
	} {
		t.Run(tc.name,func(t *testing.T){
			candidate := url.Values{}
			for k, v := range query { candidate[k] = append([]string(nil),v...) }
			if tc.key != "" { candidate.Set(tc.key,tc.value) }
			response := invoke(t,candidate,tc.bearer)
			if response.Code == http.StatusOK ||
				strings.Contains(response.Body.String(),"local-jsonl-entry") ||
				strings.Contains(response.Body.String(),"private-first-feed") {
				t.Fatalf("invalid search exposed events: %d %s",response.Code,response.Body.String())
			}
		})
	}
	duplicate := url.Values{}
	for k,v := range query { duplicate[k] = append([]string(nil),v...) }
	duplicate.Add("value","5710")
	if resp:=invoke(t,duplicate,token);resp.Code==http.StatusOK {t.Fatal("duplicate selector accepted")}

	post := httptest.NewRequest(http.MethodPost, "/demo-wazuh-search?"+query.Encode(), bytes.NewReader(nil))
	post.Header.Set("Authorization","Bearer "+token)
	postResult:=httptest.NewRecorder()
	handler.ServeHTTP(postResult,post)
	if postResult.Code != http.StatusMethodNotAllowed {t.Fatalf("unsafe HTTP method: %d",postResult.Code)}
}

func TestAuthorizedWazuhRenderedPreviewLiteralSearch(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t,
		"{\"id\":\"x.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"Suspicious PowerShell use on host\"}}\n"+
		"{\"id\":\"x.2\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"benign scheduled maintenance\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"x.1\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5730\",\"description\":\"POWERSHELL executed on second feed\"}}\n")
	data, imported, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(),
		"analyst-a", owner, caseID, []string{first, second})
	if err != nil { t.Fatal(err) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoWazuhCohortReceiptHandler(token, "analyst-a", caseID, data, imported, nil,
		levels, base, func(context.Context, io.Writer) error { return nil })
	q := url.Values{
		"case_id": {caseID}, "field": {"preview"}, "value": {"powershell"},
		"from": {"2026-10-17T09:59:00Z"}, "to": {"2026-10-17T10:03:00Z"},
	}
	send := func(q url.Values, bearer string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, "/demo-wazuh-search?"+q.Encode(), nil)
		if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		return out
	}
	out := send(q, token)
	var receipt wazuhHTTPJobReceipt
	if out.Code != http.StatusOK || json.Unmarshal(out.Body.Bytes(), &receipt) != nil ||
		receipt.Job.State != search.JobCompleted ||
		receipt.Job.Provenance.QueryVersion != "wazuh-rendered-preview-literal-v1" ||
		len(receipt.Job.ResultRefs) != 2 || len(receipt.Hits) != 2 {
		t.Fatalf("owner-authorized rendered literal search failed: %d %s", out.Code, out.Body.String())
	}
	refs := map[string]bool{}
	for _, hit := range receipt.Hits { refs[hit.EntryID] = true }
	if !refs[importedWazuhID(caseID, "x.1", 1)] ||
		!refs[importedWazuhID(caseID, "x.1", 2)] ||
		refs[importedWazuhID(caseID, "x.2", 1)] {
		t.Fatalf("literal search results not isolated: %+v", receipt.Hits)
	}
	for _, secret := range []string{"Suspicious PowerShell use", "POWERSHELL executed", "benign scheduled", token} {
		if strings.Contains(out.Body.String(), secret) {
			t.Fatalf("rendered source payload leaked through SearchJob: %q", secret)
		}
	}
	// No match produces a complete empty Job, not an unfiltered fallback.
	q.Set("value", "no-such-literal")
	empty := send(q, token)
	var noHit wazuhHTTPJobReceipt
	if empty.Code != http.StatusOK || json.Unmarshal(empty.Body.Bytes(), &noHit) != nil ||
		noHit.Job.State != search.JobCompleted || len(noHit.Job.ResultRefs) != 0 {
		t.Fatalf("literal no-match should be a completed empty Job: %d %s", empty.Code, empty.Body.String())
	}
	// The renderer is NOT a raw full_log search: rule.description wins when
	// it exists, even if full_log contains the requested term.
	for _, tc := range []struct { name,key,value,bearer string }{
		{"no bearer","","",""},
		{"foreign Case","case_id","foreign-case",token},
		{"missing literal","value","xx",token},
		{"overlong literal","value",strings.Repeat("x",129),token},
		{"newline injection","value","bad\ntext",token},
		{"unsupported DSL selector","field","full_log",token},
		{"invalid time","from","not-a-time",token},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := url.Values{}
			for k,v := range q { args[k]=append([]string(nil),v...) }
			if tc.key!="" { args.Set(tc.key,tc.value) }
			failed:=send(args,tc.bearer)
			if failed.Code == http.StatusOK || strings.Contains(failed.Body.String(),"local-jsonl-entry") {
				t.Fatalf("unsafe literal query returned data: %d %s",failed.Code,failed.Body.String())
			}
		})
	}
}
