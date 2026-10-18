package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestWazuhFeedsPersistReopenAndSearchWithoutOriginalInputs(t *testing.T) {
	const caseID="forensic-case-on-disk"
	owner:=snapshotFixture(t)
	feedA:=writeJSONLEvents(t,
		"{\"id\":\"overlapping-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell execution on endpoint\",\"level\":9}}\n")
	feedB:=writeJSONLEvents(t,
		"{\"id\":\"overlapping-id\",\"timestamp\":\"2026-10-17T10:00:45Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"description\":\"POWERSHELL execution second collector\",\"level\":8}}\n"+
		"{\"id\":\"benign-id\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"rule\":{\"id\":\"1000\",\"description\":\"scheduled maintenance\"}}\n")
	dest:=filepath.Join(t.TempDir(),"private-authorized-wazuh-case.json")
	var receipt bytes.Buffer
	if err:=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",owner,"-into-case",caseID,
		"-wazuh-alerts",feedA,"-wazuh-alerts",feedB,
		"-save-wazuh-case-snapshot",dest,
	},&receipt);err!=nil {
		t.Fatalf("Wazuh to private Case persistence CLI failed: %v",err)
	}
	blob,err:=os.ReadFile(dest)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(blob))
	if !strings.Contains(receipt.String(),digest) || strings.Contains(receipt.String(),"PowerShell") ||
		strings.Contains(receipt.String(),"POWERSHELL") || strings.Contains(receipt.String(),"Bearer") {
		t.Fatalf("persistence receipt leaks Wazuh evidence or lacks digest: %q",receipt.String())
	}
	data,err:=loadDemoSnapshotWithExpectedDigest(dest,digest)
	if err!=nil {t.Fatalf("persisted Wazuh owner did not validate by byte hash: %v",err)}
	ids,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(ids)!=3 {
		t.Fatalf("Wazuh alert original-ID collision lost independent Case source refs: %d %v",len(ids),err)
	}
	if len(data.Cases)!=1 || len(data.Principals)!=1 {
		t.Fatalf("Wazuh owner was not saved as a single least-privilege Case")
	}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	base,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,ids,base)
	requestBody:=[]byte("{\"case_id\":\""+caseID+"\",\"value\":\"powershell\",\"from\":\"2026-10-17T09:59:00Z\",\"to\":\"2026-10-17T10:02:00Z\"}")
	req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(requestBody))
	req.Header.Set("Authorization","Bearer "+token)
	req.Header.Set("Content-Type","application/json")
	response:=httptest.NewRecorder()
	handler.ServeHTTP(response,req)
	if response.Code!=http.StatusOK {
		t.Fatalf("persisted Wazuh Case source hunt failed: %d %s",response.Code,response.Body.String())
	}
	var found struct{Job search.SearchJob `json:"job"`}
	if err:=json.Unmarshal(response.Body.Bytes(),&found);err!=nil {t.Fatal(err)}
	if found.Job.State!=search.JobCompleted || len(found.Job.ResultRefs)!=2 {
		t.Fatalf("authorized Search Job failed to correlate two Wazuh collectors: %+v",found.Job)
	}
	seen:=map[string]bool{}
	for _,ref:=range found.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource {t.Fatalf("unexpected source in persisted Wazuh search: %+v",ref)}
		if _,ok:=ids[ref.Ref];!ok {t.Fatalf("cross-Case source reference leaked: %+v",ref)}
		if seen[ref.Ref] {t.Fatalf("distinct Wazuh collectors collapsed: %+v",ref)}
		seen[ref.Ref]=true
		sourceReq:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+ref.Ref,nil)
		sourceReq.Header.Set("Authorization","Bearer "+token)
		sourceResponse:=httptest.NewRecorder()
		base.ServeHTTP(sourceResponse,sourceReq)
		if sourceResponse.Code!=http.StatusOK {
			t.Fatalf("saved Wazuh typed source no longer separately authorized: %d",sourceResponse.Code)
		}
	}
	for _,secret:=range []string{"PowerShell execution on endpoint","POWERSHELL execution second collector",
		"scheduled maintenance",token} {
		if strings.Contains(response.Body.String(),secret) {
			t.Fatalf("SearchJob leaked raw source preview or bearer: %q",secret)
		}
	}
	if err:=os.Remove(feedA);err!=nil {t.Fatal(err)}
	if err:=os.Remove(feedB);err!=nil {t.Fatal(err)}
	var offline bytes.Buffer
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",dest,caseID,digest,&offline);err!=nil {
		t.Fatalf("Wazuh imports could not be verified after original alert files were removed: %v",err)
	}
	if !strings.Contains(offline.String(),"\"authorized_imported_events\":3") ||
		strings.Contains(offline.String(),"PowerShell") {
		t.Fatalf("offline Wazuh Case verification failed or leaked preview: %q",offline.String())
	}
	ctx,cancel:=context.WithTimeout(context.Background(),600*time.Millisecond)
	defer cancel()
	var viewer bytes.Buffer
	if err:=runDemoCLI(ctx,[]string{
		"-principal","analyst-a","-snapshot",dest,"-resume-import-search-case",caseID,
		"-expect-snapshot-sha256",digest,
	},&viewer);err!=nil || !strings.Contains(viewer.String(),"Local analyst viewer: http://127.0.0.1:") {
		t.Fatalf("saved Wazuh Case cannot resume independently: %v %q",err,viewer.String())
	}
}

func TestWazuhPrivatePersistenceProtectsExistingEvidenceAndCLIIsolation(t *testing.T) {
	const caseID="forensic-case-on-disk"
	owner:=snapshotFixture(t)
	feed:=writeJSONLEvents(t,
		"{\"id\":\"valid\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"rule\":{\"description\":\"legitimate test alert\"}}\n")
	dest:=filepath.Join(t.TempDir(),"new-private-case.json")
	var receipt bytes.Buffer
	if err:=runDemoCLI(context.Background(),[]string{
		"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,
		"-save-wazuh-case-snapshot",dest,
	},&receipt);err!=nil {t.Fatalf("single Wazuh save must use same owner path: %v",err)}
	data,err:=loadDemoSnapshot(dest)
	if err!=nil {t.Fatal(err)}
	ids,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(ids)!=1 {t.Fatalf("single Wazuh source not retained: %v",err)}
	sentinel:=[]byte("existing sensitive chain of evidence")
	existing:=filepath.Join(t.TempDir(),"not-to-overwrite.json")
	if err:=os.WriteFile(existing,sentinel,0600);err!=nil {t.Fatal(err)}
	for _,mode:=range [][]string{
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot",existing},
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot","relative.json"},
		{"-snapshot",owner,"-into-case","missing-case","-wazuh-alerts",feed,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"denied.json")},
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"mixed.json"),"-serve-import"},
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"mixed.json"),"-save-import-snapshot","other.json"},
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"mixed.json"),"-import-jsonl",feed},
		{"-snapshot",owner,"-into-case",caseID,"-wazuh-alerts",feed,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"mixed.json"),"-triage-import"},
		{"-snapshot",owner,"-into-case",caseID,"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"no-alerts.json")},
	} {
		var output bytes.Buffer
		if err:=runDemoCLI(context.Background(),mode,&output);err==nil || output.Len()!=0 {
			t.Fatalf("invalid/unsafe Wazuh save returned receipt or success: %v %q",mode,output.String())
		}
	}
	unchanged,err:=os.ReadFile(existing)
	if err!=nil || !bytes.Equal(unchanged,sentinel) {
		t.Fatal("operator evidence was overwritten by Wazuh Case save")
	}
}

func TestPersistedWazuhRuleAndAgentSearchJobsAfterRestart(t *testing.T) {
	const caseID="forensic-case-on-disk"
	owner:=snapshotFixture(t)
	one:=writeJSONLEvents(t,
		"{\"id\":\"same-alert\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell activity\" ,\"level\":9}}\n")
	two:=writeJSONLEvents(t,
		"{\"id\":\"same-alert\",\"timestamp\":\"2026-10-17T10:00:15Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"description\":\"Shell alert on other feed\",\"level\":6}}\n"+
		"{\"id\":\"other-alert\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"1000\",\"description\":\"routine maintenance\"}}\n")
	dest:=filepath.Join(t.TempDir(),"authorized-wazuh-owner.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",owner,caseID,[]string{one,two},dest,&bytes.Buffer{});err!=nil {
		t.Fatalf("Wazuh authorized save failed: %v",err)
	}
	raw,err:=os.ReadFile(dest)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(raw))
	if err:=os.Remove(one);err!=nil {t.Fatal(err)}
	if err:=os.Remove(two);err!=nil {t.Fatal(err)}
	data,ids,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",dest,caseID,digest)
	if err!=nil || len(ids)!=3 {t.Fatalf("Wazuh saved Case not independently restorable: %d %v",len(ids),err)}
	metaCount:=0
	for _,src:=range data.Sources {
		if src.Wazuh==nil {continue}
		metaCount++
		if src.Wazuh.RuleID=="5710" && (src.Wazuh.RuleLevel==nil ||
			*src.Wazuh.RuleLevel<6 || *src.Wazuh.RuleLevel>9) {
			t.Fatalf("saved canonical Wazuh severity lost: %+v",src.Wazuh)
		}
		if src.Wazuh.Feed<1 || src.Wazuh.Feed>2 {
			t.Fatalf("saved source feed discriminator missing: %+v",src.Wazuh)
		}
	}
	if metaCount!=3 {t.Fatalf("expected 3 owner-attached Wazuh source metadata records, got %d",metaCount)}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	ownerHandler,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	searchHandler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,ids,ownerHandler)
	invoke:=func(field,value,bearer string) *httptest.ResponseRecorder {
		t.Helper()
		requestBody,_:=json.Marshal(map[string]string{
			"case_id":caseID,"field":field,"value":value,
			"from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:02:00Z"})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(requestBody))
		req.Header.Set("Authorization","Bearer "+bearer)
		req.Header.Set("Content-Type","application/json")
		out:=httptest.NewRecorder()
		searchHandler.ServeHTTP(out,req)
		return out
	}
	for _,tc:=range []struct{field,value string;want int}{
		{"rule_id","5710",2},
		{"rule_id","1000",1},
		{"agent_id","001",1},
		{"agent_id","002",2},
		{"rule_id","9999",0},
		{"preview","powershell",1},
	}{
		out:=invoke(tc.field,tc.value,token)
		var receipt struct { Job search.SearchJob `json:"job"` }
		if out.Code!=http.StatusOK || json.Unmarshal(out.Body.Bytes(),&receipt)!=nil ||
			receipt.Job.State!=search.JobCompleted || len(receipt.Job.ResultRefs)!=tc.want {
			t.Fatalf("saved Wazuh field=%s value=%s pivot failed: %d %s",tc.field,tc.value,out.Code,out.Body.String())
		}
		for _,result:=range receipt.Job.ResultRefs {
			if _,authorized:=ids[result.Ref];!authorized {
				t.Fatalf("pivot returned unauthorized source: %+v",result)
			}
			srcReq:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+result.Ref,nil)
			srcReq.Header.Set("Authorization","Bearer "+token)
			opened:=httptest.NewRecorder()
			ownerHandler.ServeHTTP(opened,srcReq)
			if opened.Code!=http.StatusOK {t.Fatalf("returned source not reopenable by owner: %s",result.Ref)}
		}
		for _,secret:=range []string{"PowerShell activity","Shell alert on other feed","routine maintenance",token} {
			if strings.Contains(out.Body.String(),secret) {t.Fatalf("SearchJob leaked raw source or bearer: %q",secret)}
		}
	}
	for _,bad:=range []struct{field,value,bearer string}{
		{"rule_id","x5710",token},
		{"agent_id","10.1.1.1",token},
		{"rule_id","5710","bad-bearer"},
		{"agent_id","",token},
		{"full_log","PowerShell",token},
	}{
		out:=invoke(bad.field,bad.value,bad.bearer)
		if out.Code==http.StatusOK || strings.Contains(out.Body.String(),"local-jsonl-entry-") {
			t.Fatalf("invalid or unauthorized Wazuh pivot returned source: %s %d",bad.field,out.Code)
		}
	}
	// A saved owner with malformed numerical fields must fail validation
	// on reload, before it can issue any Case or source access.
	for _,change:=range []func(*timeline.WazuhSourceMetadata){
		func(m *timeline.WazuhSourceMetadata){m.AgentID="other-tenant-host"},
		func(m *timeline.WazuhSourceMetadata){v:=17;m.RuleLevel=&v},
		func(m *timeline.WazuhSourceMetadata){m.Feed=5},
	}{
		copyData:=data
		copyData.Sources=append([]timeline.MemoryTimelineSource(nil),data.Sources...)
		for i:=range copyData.Sources {
			if copyData.Sources[i].Wazuh==nil {continue}
			meta:=*copyData.Sources[i].Wazuh
			change(&meta)
			copyData.Sources[i].Wazuh=&meta
			break
		}
		if _,err:=newDemoHandlerWithData(token,"analyst-a",copyData);err==nil {
			t.Fatal("invalid or spoofed Wazuh source metadata trusted by owner")
		}
	}
}

// The saved owner and Search Job index must be immutable after validation.
// This is a real isolation boundary: RuleLevel is a *int; copying the struct
// without copying the pointed severity lets an unrelated caller mutate
// source triage displayed in already-authorized completed SearchJobs.
func TestSavedWazuhSearchMetadataIgnoresCallerMutationsAfterStartup(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	input := writeJSONLEvents(t,
		"{\"id\":\"mutable-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"authorized immutable Wazuh source\",\"level\":9}}\n")
	out := filepath.Join(t.TempDir(), "persisted-metadata-owner.json")
	if err := writePersistedWazuhCase(context.Background(), "analyst-a", owner, caseID,
		[]string{input}, out, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	data, err := loadDemoSnapshot(out)
	if err != nil { t.Fatal(err) }
	ids, err := importedNormalizedIDs(data, caseID)
	if err != nil { t.Fatal(err) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	typedOwner, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	searchHandler := newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, ids, typedOwner)

	// Deliberately mutate the caller's original Wazuh pointer after both
	// immutable owner and search index have been constructed.
	for i := range data.Sources {
		if data.Sources[i].Wazuh != nil {
			data.Sources[i].Wazuh.RuleID = "9999"
			if data.Sources[i].Wazuh.RuleLevel == nil { t.Fatal("fixture has no rule level") }
			*data.Sources[i].Wazuh.RuleLevel = 17 // invalid post-validation severity
		}
	}
	request := func(value string) *httptest.ResponseRecorder {
		t.Helper()
		query, _ := json.Marshal(map[string]string{
			"case_id": caseID, "field": "rule_id", "value": value,
			"from": "2026-10-17T09:59:00Z", "to": "2026-10-17T10:01:00Z",
		})
		req := httptest.NewRequest(http.MethodPost, "/demo-local-search", bytes.NewReader(query))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		searchHandler.ServeHTTP(resp, req)
		return resp
	}
	matched := request("5710")
	var receipt struct {
		Job search.SearchJob `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
	}
	if matched.Code != http.StatusOK ||
		json.Unmarshal(matched.Body.Bytes(), &receipt) != nil ||
		receipt.Job.State != search.JobCompleted ||
		len(receipt.Job.ResultRefs) != 1 || len(receipt.Hits) != 1 ||
		receipt.Hits[0].RuleID != "5710" ||
		receipt.Hits[0].RuleLevel == nil || *receipt.Hits[0].RuleLevel != 9 {
		t.Fatalf("validated Wazuh source metadata was mutable after startup: %d %s", matched.Code, matched.Body.String())
	}
	stale := request("9999")
	var denied struct{ Job search.SearchJob `json:"job"` }
	if stale.Code != http.StatusOK || json.Unmarshal(stale.Body.Bytes(), &denied) != nil ||
		len(denied.Job.ResultRefs) != 0 {
		t.Fatalf("post-validation caller metadata mutation altered authorized query: %d %s", stale.Code, stale.Body.String())
	}
}

func TestPersistedWazuhSearchBurstsFromCompletedAuthorizedJob(t *testing.T) {
	const caseID="forensic-case-on-disk"
	owner:=snapshotFixture(t)
	first:=writeJSONLEvents(t,
		"{\"id\":\"same-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell first alert\",\"level\":9}}\n"+
		"{\"id\":\"second\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell follow up\",\"level\":8}}\n")
	second:=writeJSONLEvents(t,
		"{\"id\":\"same-id\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell second feed\",\"level\":10}}\n"+
		"{\"id\":\"late\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"Later separate alert\"}}\n"+
		"{\"id\":\"other\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"1000\",\"description\":\"Unrelated maintenance\"}}\n")
	saved:=filepath.Join(t.TempDir(),"scoped-wazuh-owner.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",owner,caseID,
		[]string{first,second},saved,&bytes.Buffer{});err!=nil {t.Fatalf("persist Wazuh: %v",err)}
	blob,err:=os.ReadFile(saved)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(blob))
	if err:=os.Remove(first);err!=nil {t.Fatal(err)}
	if err:=os.Remove(second);err!=nil {t.Fatal(err)}
	data,imported,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",saved,caseID,digest)
	if err!=nil || len(imported)!=5 {t.Fatalf("owner lost independent feeds after restart: %d %v",len(imported),err)}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	ownerHTTP,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,ownerHTTP)
	query:=[]byte("{\"case_id\":\""+caseID+"\",\"field\":\"rule_id\",\"value\":\"5710\",\"from\":\"2026-10-17T09:59:00Z\",\"to\":\"2026-10-17T10:03:00Z\"}")
	send:=func(bearer string)*httptest.ResponseRecorder {
		t.Helper()
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(query))
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Authorization","Bearer "+bearer)
		out:=httptest.NewRecorder()
		handler.ServeHTTP(out,req)
		return out
	}
	ok:=send(token)
	var receipt struct{
		Job search.SearchJob `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
		Bursts []wazuhSearchBurst `json:"search_bursts"`
	}
	if ok.Code!=http.StatusOK || json.Unmarshal(ok.Body.Bytes(),&receipt)!=nil ||
		receipt.Job.State!=search.JobCompleted || len(receipt.Job.ResultRefs)!=4 ||
		len(receipt.Hits)!=4 || len(receipt.Bursts)!=1 {
		t.Fatalf("persisted authorized Wazuh correlation failed: %d %s",ok.Code,ok.Body.String())
	}
	burst:=receipt.Bursts[0]
	if burst.AgentID!="001" || burst.RuleID!="5710" ||
		burst.AlertCount!=3 || burst.DistinctFeeds!=2 || burst.WindowSeconds!=60 ||
		burst.FirstUTC!="2026-10-17T10:00:00Z" ||
		burst.LastUTC!="2026-10-17T10:00:40Z" || len(burst.EntryIDs)!=3 {
		t.Fatalf("wrong saved Wazuh 60-second leads: %+v",burst)
	}
	refs:=make(map[string]bool)
	for _,ref:=range receipt.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource {t.Fatalf("foreign source ref: %+v",ref)}
		refs[ref.Ref]=true
	}
	for _,hit:=range receipt.Hits {
		if !refs[hit.EntryID] {t.Fatalf("hit not present in authorized completed Search Job: %+v",hit)}
	}
	for _,ref:=range burst.EntryIDs {
		if !refs[ref] {t.Fatalf("burst linked foreign source: %s",ref)}
		req:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+ref,nil)
		req.Header.Set("Authorization","Bearer "+token)
		opened:=httptest.NewRecorder()
		ownerHTTP.ServeHTTP(opened,req)
		if opened.Code!=http.StatusOK {t.Fatalf("burst member not independently owner-openable: %s",ref)}
	}
	for _,secret:=range []string{"PowerShell first alert","PowerShell follow up",
		"PowerShell second feed","Later separate alert","Unrelated maintenance",token} {
		if strings.Contains(ok.Body.String(),secret) {t.Fatalf("source preview/bearer leaked through burst: %q",secret)}
	}
	denied:=send("not-this-owner")
	if denied.Code==http.StatusOK || strings.Contains(denied.Body.String(),"local-jsonl-entry-") {
		t.Fatalf("unauthorized caller received burst or source identities: %d %s",denied.Code,denied.Body.String())
	}
}
