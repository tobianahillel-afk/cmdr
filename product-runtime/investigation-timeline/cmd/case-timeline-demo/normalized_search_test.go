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

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestNormalizedJSONLSearchOwnsRealImportedSources(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil { t.Fatal(err) }
	path := writeJSONLEvents(t,
		"{\"id\":\"entry.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"Unauthorized PowerShell execution\"}\n" +
		"{\"id\":\"entry.2\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"message\":\"Routine service restart\"}\n" +
		"{\"id\":\"entry.3\",\"timestamp\":\"2026-10-17T10:05:00Z\",\"message\":\"POWERSHELL late activity\"}\n")
	content, err := readLocalEventFile(path)
	if err != nil { t.Fatal(err) }
	data, err := importLocalEvents(context.Background(), snapshot, "analyst-a", caseID, content)
	if err != nil { t.Fatal(err) }
	imported := make(map[string]struct{})
	for _, entry := range data.Entries[len(snapshot.Entries):] {
		imported[entry.Value.ID] = struct{}{}
	}
	if len(imported) != 3 { t.Fatalf("unexpected import count %d",len(imported)) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	base, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	handler := newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, imported, base)
	query := map[string]string{
		"case_id":caseID, "value":"powershell",
		"from":"2026-10-17T09:59:00Z", "to":"2026-10-17T10:02:00Z",
	}
	invoke := func(t *testing.T, value map[string]string, bearer string, method string) *httptest.ResponseRecorder {
		t.Helper()
		body, err := json.Marshal(value)
		if err != nil { t.Fatal(err) }
		req := httptest.NewRequest(method, "/demo-local-search", bytes.NewReader(body))
		req.Header.Set("Content-Type","application/json")
		if bearer != "" { req.Header.Set("Authorization", "Bearer "+bearer) }
		response := httptest.NewRecorder()
		handler.ServeHTTP(response,req)
		return response
	}
	ok := invoke(t,query,token,http.MethodPost)
	var got struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		Job search.SearchJob `json:"job"`
	}
	if ok.Code!=http.StatusOK || json.Unmarshal(ok.Body.Bytes(),&got)!=nil ||
		got.Format!="cmdr-local-jsonl-event-search-v1" || !got.Synthetic ||
		got.CaseID!=caseID || got.Job.State!=search.JobCompleted ||
		got.Job.Provenance.QueryVersion!="normalized-jsonl-rendered-literal-v1" ||
		len(got.Job.ResultRefs)!=1 || len(got.Job.AuditTrail)<3 {
		t.Fatalf("real normalized source SearchJob did not complete: %d %s",ok.Code,ok.Body.String())
	}
	for _, secret:=range []string{"Unauthorized PowerShell","Routine service","POWERSHELL late",token} {
		if strings.Contains(ok.Body.String(),secret) {
			t.Fatalf("source content or bearer leaked through Job: %q",secret)
		}
	}
	hit:=got.Job.ResultRefs[0]
	if hit.Source!=localJSONLSearchSource {
		t.Fatalf("unexpected source in Job: %+v",hit)
	}
	if _, authorized:=imported[hit.Ref]; !authorized {
		t.Fatalf("search included pre-existing or unrelated Case source: %+v",hit)
	}
	opened := httptest.NewRecorder()
	sourceReq := httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+hit.Ref,nil)
	sourceReq.Header.Set("Authorization","Bearer "+token)
	base.ServeHTTP(opened,sourceReq)
	if opened.Code!=http.StatusOK {
		t.Fatalf("matched generic JSONL source not independently openable: %d",opened.Code)
	}
	query["value"]="nonexistent-literal"
	empty:=invoke(t,query,token,http.MethodPost)
	var emptyJob struct{ Job search.SearchJob `json:"job"` }
	if empty.Code!=http.StatusOK || json.Unmarshal(empty.Body.Bytes(),&emptyJob)!=nil ||
		emptyJob.Job.State!=search.JobCompleted || len(emptyJob.Job.ResultRefs)!=0 {
		t.Fatalf("authorized no-match should complete empty: %d %s",empty.Code,empty.Body.String())
	}
	for _, tc := range []struct {
		name,key,value,bearer,method string
	}{
		{"no bearer","","","",http.MethodPost},
		{"wrong bearer","","","wrong",http.MethodPost},
		{"foreign Case","case_id","case-b",token,http.MethodPost},
		{"short literal","value","aa",token,http.MethodPost},
		{"control byte","value","bad\nline",token,http.MethodPost},
		{"invalid start","from","not-a-time",token,http.MethodPost},
		{"reverse time","from","2026-10-17T11:00:00Z",token,http.MethodPost},
		{"GET forbidden","","",token,http.MethodGet},
	}{
		t.Run(tc.name,func(t *testing.T) {
			testQuery:=make(map[string]string)
			for k,v:=range query {testQuery[k]=v}
			if tc.key!="" {testQuery[tc.key]=tc.value}
			denied:=invoke(t,testQuery,tc.bearer,tc.method)
			if denied.Code==http.StatusOK || strings.Contains(denied.Body.String(),hit.Ref) {
				t.Fatalf("denied request leaked source identities: %d %s",denied.Code,denied.Body.String())
			}
		})
	}
}

func TestNormalizedJSONLSearchRejectsOversizeWithoutPartialHits(t *testing.T) {
	const caseID="forensic-case-on-disk"
	snapshot,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	var raw strings.Builder
	for i:=0;i<501;i++ {
		fmt.Fprintf(&raw,"{\"id\":\"e-%d\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"matching denied overflow\"}\n",i)
	}
	data,err:=importLocalEvents(context.Background(),snapshot,"analyst-a",caseID,[]byte(raw.String()))
	if err!=nil {t.Fatal(err)}
	imported:=make(map[string]struct{})
	for _,entry:=range data.Entries[len(snapshot.Entries):]{imported[entry.Value.ID]=struct{}{}}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	base,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,base)
	body:=[]byte("{\"case_id\":\""+caseID+"\",\"value\":\"matching\",\"from\":\"2026-10-17T09:59:00Z\",\"to\":\"2026-10-17T10:01:00Z\"}")
	req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
	req.Header.Set("Authorization","Bearer "+token)
	req.Header.Set("Content-Type","application/json")
	out:=httptest.NewRecorder()
	handler.ServeHTTP(out,req)
	if out.Code==http.StatusOK || strings.Contains(out.Body.String(),"local-jsonl-entry-") {
		t.Fatalf("overflow must fail the ENTIRE authorized Job: %d %.100s",out.Code,out.Body.String())
	}
}
