package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestSavedMixedWazuhAndGenericSourcesShareAuthorizedSearchAndBursts(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"alert-A\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell suspicious WAZUH-SECRET-A\",\"level\":9}}\n"+
		"{\"id\":\"alert-B\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell suspicious WAZUH-SECRET-B\",\"level\":8}}\n"+
		"{\"id\":\"alert-C\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell suspicious WAZUH-SECRET-C\",\"level\":7}}\n")
	generic := writeJSONLEvents(t,
		"{\"id\":\"generic-auth-event\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PowerShell generic GENERIC-SECRET source evidence\"}\n")
	first := filepath.Join(t.TempDir(), "saved-wazuh-owner.json")
	var firstReceipt bytes.Buffer
	if err:=runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", owner, "-into-case", caseID,
		"-wazuh-alerts", wazuh, "-save-wazuh-case-snapshot", first,
	}, &firstReceipt);err!=nil {
		t.Fatalf("existing Wazuh Case save failed: %v",err)
	}
	mixed := filepath.Join(t.TempDir(), "saved-mixed-case.json")
	if err:=runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", first, "-into-case", caseID,
		"-import-jsonl", generic, "-save-import-snapshot", mixed,
	}, &bytes.Buffer{});err!=nil {
		t.Fatalf("existing generic JSONL append to Wazuh owner failed: %v",err)
	}
	data,err:=loadDemoSnapshot(mixed)
	if err!=nil {t.Fatal(err)}
	ids,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(ids)!=4 {
		t.Fatalf("mixed validated owner must contain 3 Wazuh and 1 generic source: %d %v",len(ids),err)
	}
	var genericRef string
	wazuhRefIDs:=map[string]bool{}
	for _,entry:=range data.Entries {
		if _,ok:=ids[entry.Value.ID];!ok {continue}
		hasWazuh:=false
		for _,src:=range data.Sources {
			if src.Identity.Reference==entry.Value.SourceRef && src.Wazuh!=nil {
				hasWazuh=true
				break
			}
		}
		if hasWazuh {wazuhRefIDs[entry.Value.ID]=true} else {genericRef=entry.Value.ID}
	}
	if len(wazuhRefIDs)!=3 || genericRef=="" {t.Fatalf("mixed source-owned metadata not preserved: %d %s",len(wazuhRefIDs),genericRef)}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	ownerHandler,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,ids,ownerHandler)
	invoke:=func(field,value,bearer,selectedCase string)*httptest.ResponseRecorder {
		body,_:=json.Marshal(map[string]string{
			"case_id":selectedCase, "field":field, "value":value,
			"from":"2026-10-17T09:59:00Z", "to":"2026-10-17T10:02:00Z",
		})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Authorization","Bearer "+bearer)
		req.Header.Set("Content-Type","application/json")
		rec:=httptest.NewRecorder()
		handler.ServeHTTP(rec,req)
		return rec
	}
	type receipt struct {
		Job search.SearchJob `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
		Bursts []wazuhSearchBurst `json:"search_bursts"`
	}
	previewResponse:=invoke("preview","powershell",token,caseID)
	var preview receipt
	if previewResponse.Code!=http.StatusOK || json.Unmarshal(previewResponse.Body.Bytes(),&preview)!=nil ||
		preview.Job.State!=search.JobCompleted || len(preview.Job.ResultRefs)!=4 ||
		len(preview.Hits)!=3 || len(preview.Bursts)!=1 ||
		preview.Bursts[0].AlertCount!=3 || len(preview.Bursts[0].EntryIDs)!=3 {
		t.Fatalf("mixed preview hunt lost generic SearchJob results or Wazuh chronology: %d %s",previewResponse.Code,previewResponse.Body.String())
	}
	seen:=make(map[string]bool)
	for _,ref:=range preview.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource || seen[ref.Ref] {t.Fatalf("bad mixed source ref: %+v",ref)}
		seen[ref.Ref]=true
		sourceReq:=httptest.NewRequest(http.MethodGet,
			"/demo-source-open?case_id="+caseID+"&entry_id="+ref.Ref,nil)
		sourceReq.Header.Set("Authorization","Bearer "+token)
		opened:=httptest.NewRecorder()
		ownerHandler.ServeHTTP(opened,sourceReq)
		if opened.Code!=http.StatusOK {t.Fatalf("Search Job ref cannot independently re-open: %s",ref.Ref)}
	}
	if !seen[genericRef] {t.Fatal("generic event missing from mixed owner-authorized Search Job")}
	for _,id:=range preview.Bursts[0].EntryIDs {
		if !wazuhRefIDs[id] {t.Fatalf("generic source incorrectly counted as Wazuh burst member: %q",id)}
	}
	for _,field:=range []string{"rule_id","agent_id"} {
		query:="5710"
		if field=="agent_id" {query="001"}
		got:=invoke(field,query,token,caseID)
		var result receipt
		if got.Code!=http.StatusOK || json.Unmarshal(got.Body.Bytes(),&result)!=nil ||
			result.Job.State!=search.JobCompleted || len(result.Job.ResultRefs)!=3 ||
			len(result.Hits)!=3 || len(result.Bursts)!=1 {
			t.Fatalf("Wazuh numeric pivot must remain usable in mixed Case: %s %d %s",field,got.Code,got.Body.String())
		}
		for _,ref:=range result.Job.ResultRefs {
			if !wazuhRefIDs[ref.Ref] || ref.Ref==genericRef {
				t.Fatalf("numeric Wazuh pivot matched generic event: %+v",ref)
			}
		}
	}
	// A search exclusively matching generic JSONL emits zero Wazuh hits
	// and no speculative temporal bursts, while retaining its source ref.
	onlyGeneric:=invoke("preview","GENERIC-SECRET",token,caseID)
	var genericReceipt receipt
	if onlyGeneric.Code!=http.StatusOK || json.Unmarshal(onlyGeneric.Body.Bytes(),&genericReceipt)!=nil ||
		len(genericReceipt.Job.ResultRefs)!=1 || genericReceipt.Job.ResultRefs[0].Ref!=genericRef ||
		len(genericReceipt.Hits)!=0 || len(genericReceipt.Bursts)!=0 {
		t.Fatalf("generic-only match manufactured Wazuh metadata: %d %s",onlyGeneric.Code,onlyGeneric.Body.String())
	}
	for _,secret:=range []string{"WAZUH-SECRET","GENERIC-SECRET",token} {
		if strings.Contains(previewResponse.Body.String(),secret) || strings.Contains(onlyGeneric.Body.String(),secret) {
			t.Fatalf("completed mixed Search Job disclosed raw preview/bearer: %q",secret)
		}
	}
	for _,denied:=range []struct{token,caseID string}{
		{"wrong-bearer",caseID}, {token,"other-case"},
	} {
		reply:=invoke("provenance","operator-local-jsonl",denied.token,denied.caseID)
		if reply.Code==http.StatusOK || strings.Contains(reply.Body.String(),genericRef) {
			t.Fatalf("mixed source scope leaked to unauthorized principal/Case: %d %s",reply.Code,reply.Body.String())
		}
	}
	// Material ACL revocation case: even when three Wazuh records are
	// normalized and a preview or numeric pivot otherwise matches, denying
	// one typed source must reject the ENTIRE Search Job before chronology
	// or burst leads are disclosed. A cached Wazuh metadata map is never an
	// alternative authorization authority.
	restricted := data
	restricted.Principals = append([]timeline.MemoryTimelinePrincipal(nil), data.Principals...)
	for _,p := range data.Principals {
		if p.ID == "analyst-a" {
			restricted.Principals = append(restricted.Principals,
				timeline.MemoryTimelinePrincipal{
					ID:"restricted-source-reader",Tenant:p.Tenant,Environment:p.Environment,
				})
			break
		}
	}
	restricted.Sources = append([]timeline.MemoryTimelineSource(nil),data.Sources...)
	revoked := false
	for i:=range restricted.Sources {
		if restricted.Sources[i].Wazuh != nil {
			restricted.Sources[i].Readers = []string{"restricted-source-reader"}
			revoked = true
			break
		}
	}
	if !revoked {t.Fatal("cannot exercise Wazuh source revocation")}
	limitedOwner,err:=newDemoHandlerWithData(token,"analyst-a",restricted)
	if err!=nil {t.Fatalf("valid owner with independently denied source refused: %v",err)}
	limitedHandler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,restricted,ids,limitedOwner)
	for _,filter:=range []struct {field,value string}{
		{"preview","powershell"},{"rule_id","5710"},
	} {
		body,_:=json.Marshal(map[string]string{
			"case_id":caseID,"field":filter.field,"value":filter.value,
			"from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:02:00Z",
		})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Authorization","Bearer "+token)
		req.Header.Set("Content-Type","application/json")
		rejected:=httptest.NewRecorder()
		limitedHandler.ServeHTTP(rejected,req)
		if rejected.Code==http.StatusOK ||
			strings.Contains(rejected.Body.String(),"local-jsonl-entry-") ||
			strings.Contains(rejected.Body.String(),"WAZUH-SECRET") ||
			strings.Contains(rejected.Body.String(),"search_bursts") {
			t.Fatalf("revoked Wazuh source leaked authorized results or burst leads: %s %d %s",
				filter.field,rejected.Code,rejected.Body.String())
		}
	}
	// The saved combined owner is independent of both original input files.
	if err:=os.Remove(wazuh);err!=nil {t.Fatal(err)}
	if err:=os.Remove(generic);err!=nil {t.Fatal(err)}
	if _,err:=loadDemoSnapshot(mixed);err!=nil {
		t.Fatalf("reopenable Case unexpectedly requires original input feeds: %v",err)
	}
}
