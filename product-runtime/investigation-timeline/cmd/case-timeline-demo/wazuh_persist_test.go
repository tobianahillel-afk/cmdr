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
