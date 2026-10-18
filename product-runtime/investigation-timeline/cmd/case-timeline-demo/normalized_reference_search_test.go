package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestSourceReferenceAndProvenancePivotsUseActualAuthorizedSearchJob(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	files := []string{
		writeJSONLEvents(t,
			"{\"id\":\"same-original-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"SECRET-FIRST-CASE-PREVIEW\"}\n"),
		writeJSONLEvents(t,
			"{\"id\":\"same-original-id\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"SECRET-SECOND-CASE-PREVIEW\"}\n"+
			"{\"id\":\"later\",\"timestamp\":\"2026-10-17T12:00:00Z\",\"message\":\"SECRET-LATE-PREVIEW\"}\n"),
	}
	data, imported, err := loadMultipleNormalizedJSONLFeeds(
		context.Background(), "analyst-a", owner, caseID, files)
	if err != nil { t.Fatal(err) }
	if len(imported) != 3 { t.Fatalf("expected 3 unique owned sources: %d",len(imported)) }
	byFeed := make(map[string]string)
	for _, entry := range data.Entries {
		if _, exists := imported[entry.Value.ID]; exists &&
			entry.Value.SourceTime != "2026-10-17T12:00:00Z" {
			byFeed[entry.Value.Provenance] = entry.Value.ID
		}
	}
	first := byFeed["operator-local-jsonl-feed-1"]
	second := byFeed["operator-local-jsonl-feed-2"]
	if first == "" || second == "" || first == second { t.Fatalf("colliding original events lost feed identity: %+v",byFeed) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	base, err := newDemoHandlerWithData(token,"analyst-a",data)
	if err != nil { t.Fatal(err) }
	handler := newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,base)
	run := func(field,value,bearer,caseSelector string) *httptest.ResponseRecorder {
		body, _ := json.Marshal(map[string]string{
			"case_id":caseSelector,"field":field,"value":value,
			"from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:01:00Z",
		})
		req := httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Authorization","Bearer "+bearer)
		got := httptest.NewRecorder()
		handler.ServeHTTP(got,req)
		return got
	}
	for _, test := range []struct {
		field,value,wantRef string
	}{
		{"provenance","feed-2",second},
		{"provenance","feed-1",first},
		{"entry_id",second[len(second)-25:],second},
		{"source_ref","local-jsonl-source-"+second[len("local-jsonl-entry-"):],second},
	} {
		t.Run(test.field+"-"+test.value,func(t *testing.T) {
			response:=run(test.field,test.value,token,caseID)
			var result struct {
				Job search.SearchJob `json:"job"`
			}
			if response.Code != http.StatusOK ||
				json.Unmarshal(response.Body.Bytes(),&result) != nil ||
				result.Job.State != search.JobCompleted ||
				result.Job.Provenance.QueryVersion != "authorized-normalized-owner-reference-v1" ||
				len(result.Job.ResultRefs) != 1 ||
				result.Job.ResultRefs[0].Ref != test.wantRef ||
				result.Job.ResultRefs[0].Source != localJSONLSearchSource {
				t.Fatalf("real owner authorized source pivot failed: code=%d body=%s",response.Code,response.Body.String())
			}
			for _, secret := range []string{
				"SECRET-FIRST-CASE-PREVIEW","SECRET-SECOND-CASE-PREVIEW",
				"SECRET-LATE-PREVIEW",token,
			} {
				if strings.Contains(response.Body.String(),secret) { t.Fatalf("raw source leaked into search receipt: %q",secret) }
			}
			source := httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+test.wantRef,nil)
			source.Header.Set("Authorization","Bearer "+token)
			opened:=httptest.NewRecorder()
			base.ServeHTTP(opened,source)
			if opened.Code!=http.StatusOK { t.Fatalf("returned reference is not independently owner-openable: %d",opened.Code) }
		})
	}
	// Search fields are constrained to a predefined metadata enum. No raw
	// JSON property name, foreign Case or bearer may widen source scope.
	for _, denied := range []struct {field,value,bearer,caseSelector string}{
		{"raw_json","host.ip",token,caseID},
		{"provenance","feed-2","wrong",caseID},
		{"source_ref","local-jsonl-source-",token,"foreign-case"},
		{"entry_id","no",token,caseID},
	} {
		out:=run(denied.field,denied.value,denied.bearer,denied.caseSelector)
		if out.Code==http.StatusOK || strings.Contains(out.Body.String(),second) {
			t.Fatalf("unauthorized metadata selector returned owner references: %d %s",out.Code,out.Body.String())
		}
	}
	// A provenance selector for a real feed with no in-window events
	// completes with zero matches instead of silently reading other times.
	out:=run("provenance","feed-4",token,caseID)
	var empty struct{ Job search.SearchJob `json:"job"` }
	if out.Code!=http.StatusOK || json.Unmarshal(out.Body.Bytes(),&empty)!=nil ||
		empty.Job.State!=search.JobCompleted || len(empty.Job.ResultRefs)!=0 {
		t.Fatalf("in-bounds, no-match metadata search incorrectly failed: %d %s",out.Code,out.Body.String())
	}
}
