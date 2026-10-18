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

func TestMultipleNormalizedJSONLFeedsCompleteAuthorizedSearchJob(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshotPath := snapshotFixture(t)
	feeds := []string{
		writeJSONLEvents(t,
			"{\"id\":\"shared.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"PowerShell suspicious first feed\"}\n"+
			"{\"id\":\"first.only\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"message\":\"Routine unrelated service\"}\n"),
		writeJSONLEvents(t,
			"{\"id\":\"shared.1\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"message\":\"POWERSHELL suspicious second feed\"}\n"+
			"{\"id\":\"late.1\",\"timestamp\":\"2026-10-17T11:00:00Z\",\"message\":\"PowerShell outside window\"}\n"),
	}
	data, imported, err := loadMultipleNormalizedJSONLFeeds(context.Background(),
		"analyst-a",snapshotPath,caseID,feeds)
	if err != nil {t.Fatal(err)}
	if len(imported)!=4 {t.Fatalf("expected 4 separately owned imported events, got %d",len(imported))}
	shared := make([]string,0,2)
	seenFeed := map[string]bool{}
	for _,entry := range data.Entries {
		if _,ok:=imported[entry.Value.ID];!ok {continue}
		seenFeed[entry.Value.Provenance]=true
		if entry.Value.SourceTime=="2026-10-17T10:00:00Z" ||
			entry.Value.SourceTime=="2026-10-17T10:02:00Z" {
			shared=append(shared,entry.Value.ID)
		}
	}
	if len(shared)!=2 || shared[0]==shared[1] ||
		!seenFeed["operator-local-jsonl-feed-1"] ||
		!seenFeed["operator-local-jsonl-feed-2"] {
		t.Fatalf("cross-feed original ID collision was not isolated: %+v",shared)
	}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	base,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,base)
	body:=[]byte("{\"case_id\":\""+caseID+"\",\"value\":\"powershell\",\"from\":\"2026-10-17T09:59:00Z\",\"to\":\"2026-10-17T10:05:00Z\"}")
	request:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
	request.Header.Set("Authorization","Bearer "+token)
	request.Header.Set("Content-Type","application/json")
	response:=httptest.NewRecorder()
	handler.ServeHTTP(response,request)
	if response.Code!=http.StatusOK {t.Fatalf("search denied: %d %s",response.Code,response.Body.String())}
	var receipt struct{
		Job search.SearchJob `json:"job"`
	}
	if err:=json.Unmarshal(response.Body.Bytes(),&receipt);err!=nil {t.Fatal(err)}
	if receipt.Job.State!=search.JobCompleted || len(receipt.Job.ResultRefs)!=2 {
		t.Fatalf("expected two authorized cross-feed matches in Search Job: %+v",receipt.Job)
	}
	authorized:=map[string]bool{}
	for _,ref:=range receipt.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource {t.Fatalf("wrong Source in Search Job: %+v",ref)}
		authorized[ref.Ref]=true
	}
	for _,ref:=range shared {
		if !authorized[ref] {t.Fatalf("missing distinct cross-feed event %s",ref)}
	}
	for _,secret:=range []string{"PowerShell suspicious first feed","POWERSHELL suspicious second feed","Routine unrelated service",token} {
		if strings.Contains(response.Body.String(),secret) {t.Fatalf("unexpected source payload in SearchJob: %q",secret)}
	}
	for _,ref:=range shared {
		open:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+ref,nil)
		open.Header.Set("Authorization","Bearer "+token)
		ownerResponse:=httptest.NewRecorder()
		base.ServeHTTP(ownerResponse,open)
		if ownerResponse.Code!=http.StatusOK {t.Fatalf("authorized search hit no longer typed-source openable: %s",ref)}
	}
	denied:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
	denied.Header.Set("Authorization","Bearer invalid")
	denied.Header.Set("Content-Type","application/json")
	notAllowed:=httptest.NewRecorder()
	handler.ServeHTTP(notAllowed,denied)
	if notAllowed.Code==http.StatusOK || strings.Contains(notAllowed.Body.String(),shared[0]) {
		t.Fatal("unauthorized request disclosed imported event")
	}
}

func TestMultipleNormalizedJSONLRejectsDuplicatedSourcesAndInvalidModes(t *testing.T) {
	const caseID="forensic-case-on-disk"
	snapshotPath:=snapshotFixture(t)
	content:="{\"id\":\"duplicate-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"valid record\"}\n"
	a:=writeJSONLEvents(t,content)
	b:=writeJSONLEvents(t,content)
	if _,_,err:=loadMultipleNormalizedJSONLFeeds(context.Background(),"analyst-a",snapshotPath,caseID,[]string{a,a});err==nil {
		t.Fatal("same file selected twice should not create false independent sources")
	}
	data,imported,err:=loadMultipleNormalizedJSONLFeeds(context.Background(),"analyst-a",snapshotPath,caseID,[]string{a,b})
	if err!=nil || len(imported)!=2 {
		t.Fatalf("independent files with repeated original ID should have stable distinct owner identities: %v %v",len(imported),err)
	}
	_ = data
	if _,_,err:=loadMultipleNormalizedJSONLFeeds(context.Background(),"analyst-a",snapshotPath,"foreign-case",[]string{a,b});err==nil {
		t.Fatal("invalid/denied Case must not import multiple feeds")
	}
	if _,_,err:=loadMultipleNormalizedJSONLFeeds(context.Background(),"analyst-a",snapshotPath,caseID,[]string{a});err==nil {
		t.Fatal("multi-feed entrypoint should not masquerade as single-file import")
	}
	if err:=runDemoCLI(context.Background(),[]string{
		"-import-jsonl-extra",b,"-snapshot",snapshotPath,"-import-jsonl",a,"-into-case",caseID,
	},&bytes.Buffer{});err==nil {
		t.Fatal("additional feed accepted without explicit serve-import mode")
	}
}

func TestSingleNormalizedImporterLegacyIDAndProvenanceUnchanged(t *testing.T) {
	const caseID="forensic-case-on-disk"
	snapshot,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	source:=[]byte("{\"id\":\"event-stable\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"legacy event\"}\n")
	legacy,err:=importLocalEvents(context.Background(),snapshot,"analyst-a",caseID,source)
	if err!=nil {t.Fatal(err)}
	explicit,err:=importLocalEventsFromFeed(context.Background(),snapshot,"analyst-a",caseID,source,0)
	if err!=nil {t.Fatal(err)}
	l:=legacy.Entries[len(snapshot.Entries)]
	e:=explicit.Entries[len(snapshot.Entries)]
	if l.Value.ID!=e.Value.ID || l.Value.SourceRef!=e.Value.SourceRef ||
		l.Value.Provenance!="operator-local-jsonl" || e.Value.Provenance!="operator-local-jsonl" {
		t.Fatalf("single-import ID or provenance changed: %v %v",l.Value,e.Value)
	}
}
