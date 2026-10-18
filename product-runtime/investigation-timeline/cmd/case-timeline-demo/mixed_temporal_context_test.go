package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"path/filepath"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestMixedSavedCaseSearchBuildsOnlyAuthorizedCrossFamilyTemporalLeads(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	wazuhFile := writeJSONLEvents(t,
		"{\"id\":\"w-one\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-FIRST\",\"level\":7}}\n"+
		"{\"id\":\"w-two\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-SECOND\",\"level\":7}}\n")
	genericFile := writeJSONLEvents(t,
		"{\"id\":\"g-near\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PRIVATE-GENERIC-NEAR\"}\n"+
		"{\"id\":\"g-boundary\",\"timestamp\":\"2026-10-17T09:59:00Z\",\"message\":\"PRIVATE-GENERIC-BOUNDARY\"}\n"+
		"{\"id\":\"g-far\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"message\":\"PRIVATE-GENERIC-FAR\"}\n")
	saved := filepath.Join(t.TempDir(),"mixed-private-case.json")
	if err := writePersistedMixedCase(context.Background(),"analyst-a",snapshotFixture(t),caseID,
		[]string{wazuhFile},[]string{genericFile},saved,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	data,imported,err := loadAuthorizedImportedCase(context.Background(),"analyst-a",saved,caseID,"")
	if err!=nil {t.Fatal(err)}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	owner,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,owner)
	run:=func(family,caseName,bearer string)*httptest.ResponseRecorder{
		t.Helper()
		body,_:=json.Marshal(map[string]string{
			"case_id":caseName,"field":"source_ref","source_family":family,
			"value":"local-jsonl-source-","from":"2026-10-17T09:58:30Z",
			"to":"2026-10-17T10:03:00Z",
		})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Authorization","Bearer "+bearer)
		req.Header.Set("Content-Type","application/json")
		result:=httptest.NewRecorder()
		handler.ServeHTTP(result,req)
		return result
	}
	reply:=run("all",caseID,token)
	var decoded struct {
		Job search.SearchJob `json:"job"`
		SearchHits []struct{ EntryID string `json:"entry_id"` } `json:"search_hits"`
		Leads []mixedTemporalLead `json:"mixed_temporal_context"`
	}
	if reply.Code!=http.StatusOK || json.Unmarshal(reply.Body.Bytes(),&decoded)!=nil ||
		decoded.Job.State!=search.JobCompleted || len(decoded.Job.ResultRefs)!=5 ||
		len(decoded.SearchHits)!=2 || len(decoded.Leads)!=2 {
		t.Fatalf("real mixed Case SearchJob and temporal context failed: %d %s",reply.Code,reply.Body.String())
	}
	jobMembers:=map[string]bool{}
	ownerWazuh:=map[string]bool{}
	for _,ref:=range decoded.Job.ResultRefs {jobMembers[ref.Ref]=true}
	for _,hit:=range decoded.SearchHits {ownerWazuh[hit.EntryID]=true}
	if len(jobMembers)!=5 || len(ownerWazuh)!=2 {t.Fatal("duplicate or invalid Search Job members")}
	boundarySeen:=false
	for _,lead:=range decoded.Leads {
		if !jobMembers[lead.WazuhEntryID] || !ownerWazuh[lead.WazuhEntryID] ||
			lead.WindowSeconds!=60 || len(lead.GenericEvents)==0 {
			t.Fatalf("lead not anchored to authorized Wazuh SearchJob hit: %+v",lead)
		}
		for _,neighbor:=range lead.GenericEvents {
			if !jobMembers[neighbor.EntryID] || ownerWazuh[neighbor.EntryID] ||
				neighbor.OffsetSeconds < -60 || neighbor.OffsetSeconds > 60 {
				t.Fatalf("generic lead member not authorized: %+v",neighbor)
			}
			if neighbor.OffsetSeconds==-60 {boundarySeen=true}
			var open = httptest.NewRecorder()
			req:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+neighbor.EntryID,nil)
			req.Header.Set("Authorization","Bearer "+token)
			owner.ServeHTTP(open,req)
			if open.Code!=http.StatusOK {t.Fatalf("generic lead source cannot be independently reopened: %d",open.Code)}
		}
	}
	if !boundarySeen {t.Fatal("inclusive 60-second context boundary not observed")}
	for _,s:=range []string{"PRIVATE-WAZUH-FIRST","PRIVATE-WAZUH-SECOND",
		"PRIVATE-GENERIC-NEAR","PRIVATE-GENERIC-FAR",token}{
		if strings.Contains(reply.Body.String(),s){
			t.Fatalf("matched event payload or bearer leaked in contextual Search Job: %q",s)
		}
	}
	for _,family:=range []string{"wazuh","generic"}{
		single:=run(family,caseID,token)
		if single.Code!=http.StatusOK || strings.Contains(single.Body.String(),"mixed_temporal_context") {
			t.Fatalf("single-source-family query leaked cross-family context: %d %s",single.Code,single.Body.String())
		}
	}
	for _,bad:=range []struct{owner,bearer string}{
		{"foreign-case",token},{caseID,"wrong"},
	}{
		denied:=run("all",bad.owner,bad.bearer)
		if denied.Code==http.StatusOK || strings.Contains(denied.Body.String(),"mixed_temporal_context") {
			t.Fatalf("denied source owner exposed cross-family context: %d",denied.Code)
		}
	}
}

func TestMixedContextRefusesForeignOrDuplicateSearchJobMembership(t *testing.T) {
	const caseID="forensic-case-on-disk"
	source,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	for _,refs:=range [][]search.ResultReference{
		{{Source:"external",Ref:"unauthorized"}},
		{{Source:localJSONLSearchSource,Ref:"missing-from-owner"}},
		{{Source:localJSONLSearchSource,Ref:"duplicate"},
			{Source:localJSONLSearchSource,Ref:"duplicate"}},
	} {
		_,err:=deriveAuthorizedMixedTemporalContext(source,caseID,refs,
			map[string]timeline.WazuhSourceMetadata{})
		if err==nil {t.Fatalf("unsafe job membership accepted: %+v",refs)}
	}
}
