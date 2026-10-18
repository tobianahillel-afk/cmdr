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

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestOneStepMixedCaseSaveReopensBothSourcesInSingleSearchJob(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"a\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell WAZUH-A-PRIVATE\",\"level\":8}}\n"+
		"{\"id\":\"b\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell WAZUH-B-PRIVATE\",\"level\":9}}\n"+
		"{\"id\":\"c\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PowerShell WAZUH-C-PRIVATE\",\"level\":7}}\n")
	generic := writeJSONLEvents(t,
		"{\"id\":\"generic-event\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PowerShell GENERIC-PRIVATE telemetry\"}\n")
	dir := t.TempDir()
	private := filepath.Join(dir, "mixed-private-case.json")
	var console bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal","analyst-a","-snapshot",owner,"-into-case",caseID,
		"-wazuh-alerts",wazuh,"-import-jsonl",generic,
		"-save-mixed-case-snapshot",private,
	}, &console); err != nil {
		t.Fatalf("one-step private mixed Case creation failed: %v",err)
	}
	// The only output must be the final private Case. No transient Wazuh
	// owner archive exists on disk or appears in stdout.
	files,err:=os.ReadDir(dir)
	if err!=nil || len(files)!=1 || files[0].Name()!="mixed-private-case.json" {
		t.Fatalf("mixed import wrote an unexpected intermediate owner: %+v %v",files,err)
	}
	savedBytes,err:=os.ReadFile(private)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(savedBytes))
	if !strings.Contains(console.String(),digest) {
		t.Fatal("mixed Case save did not print digest for exact final output bytes")
	}
	for _,secret:=range []string{"WAZUH-A-PRIVATE","GENERIC-PRIVATE","Bearer "}{
		if strings.Contains(console.String(),secret) {t.Fatalf("console leaked source or bearer %s",secret)}
	}
	data,err:=loadDemoSnapshotWithExpectedDigest(private,digest)
	if err!=nil {t.Fatalf("new mixed Case cannot be reloaded and digest checked: %v",err)}
	ids,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(ids)!=4 {t.Fatalf("expected four owner-validated mixed entries: %d %v",len(ids),err)}
	wazuhIDs:=make(map[string]bool)
	oneGeneric:=""
	bySource:=make(map[string]timeline.MemoryTimelineSource)
	for _,src:=range data.Sources {bySource[src.Identity.Reference]=src}
	for _,entry:=range data.Entries {
		if _,ok:=ids[entry.Value.ID];!ok {continue}
		src,ok:=bySource[entry.Value.SourceRef]
		if !ok {t.Fatalf("persisted Entry lost its typed source %s",entry.Value.ID)}
		if src.Wazuh!=nil {
			if src.Wazuh.RuleID!="5710" || src.Wazuh.AgentID!="001" {
				t.Fatalf("lost canonical Wazuh rule/agent: %+v",src.Wazuh)
			}
			wazuhIDs[entry.Value.ID]=true
		} else {
			oneGeneric=entry.Value.ID
		}
	}
	if len(wazuhIDs)!=3 || oneGeneric=="" {t.Fatalf("one-step mixed Case lost provenance: %d, %s",len(wazuhIDs),oneGeneric)}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	ownerHandler,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,ids,ownerHandler)
	invoke:=func(field,value,bearer string)*httptest.ResponseRecorder {
		body,_:=json.Marshal(map[string]string{
			"case_id":caseID,"field":field,"value":value,
			"from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:02:00Z",
		})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Authorization","Bearer "+bearer)
		req.Header.Set("Content-Type","application/json")
		w:=httptest.NewRecorder()
		handler.ServeHTTP(w,req)
		return w
	}
	type output struct {
		Job search.SearchJob `json:"job"`
		Hits []wazuhSearchHit `json:"search_hits"`
		Bursts []wazuhSearchBurst `json:"search_bursts"`
	}
	all:=invoke("preview","powershell",token)
	var complete output
	if all.Code!=http.StatusOK || json.Unmarshal(all.Body.Bytes(),&complete)!=nil ||
		complete.Job.State!=search.JobCompleted || len(complete.Job.ResultRefs)!=4 ||
		len(complete.Hits)!=3 || len(complete.Bursts)!=1 || complete.Bursts[0].AlertCount!=3 {
		t.Fatalf("reopened one-step mixed Case SearchJob/correlation failed: %d %s",all.Code,all.Body.String())
	}
	reopened:=make(map[string]bool)
	for _,ref:=range complete.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource || reopened[ref.Ref] {t.Fatalf("wrong source reference %+v",ref)}
		reopened[ref.Ref]=true
		open:=httptest.NewRequest(http.MethodGet,
			"/demo-source-open?case_id="+caseID+"&entry_id="+ref.Ref,nil)
		open.Header.Set("Authorization","Bearer "+token)
		rr:=httptest.NewRecorder()
		ownerHandler.ServeHTTP(rr,open)
		if rr.Code!=http.StatusOK {t.Fatalf("SearchJob hit denied by independent owner source %s",ref.Ref)}
	}
	if !reopened[oneGeneric] {t.Fatal("generic event excluded from mixed Search Job")}
	for _,id:=range complete.Bursts[0].EntryIDs {
		if !wazuhIDs[id] {t.Fatalf("generic source misclassified as Wazuh burst %s",id)}
	}
	numeric:=invoke("rule_id","5710",token)
	var pivots output
	if numeric.Code!=http.StatusOK || json.Unmarshal(numeric.Body.Bytes(),&pivots)!=nil ||
		len(pivots.Job.ResultRefs)!=3 || len(pivots.Hits)!=3 {
		t.Fatalf("owner-linked numeric Wazuh pivot was lost: %d %s",numeric.Code,numeric.Body.String())
	}
	onlyGeneric:=invoke("preview","GENERIC-PRIVATE",token)
	var genericHit output
	if onlyGeneric.Code!=http.StatusOK || json.Unmarshal(onlyGeneric.Body.Bytes(),&genericHit)!=nil ||
		len(genericHit.Job.ResultRefs)!=1 || genericHit.Job.ResultRefs[0].Ref!=oneGeneric ||
		len(genericHit.Hits)!=0 || len(genericHit.Bursts)!=0 {
		t.Fatalf("generic-only query manufactured Wazuh results: %d %s",onlyGeneric.Code,onlyGeneric.Body.String())
	}
	for _,secret:=range []string{"WAZUH-A-PRIVATE","GENERIC-PRIVATE",token}{
		if strings.Contains(all.Body.String(),secret) || strings.Contains(onlyGeneric.Body.String(),secret) {
			t.Fatalf("mixed SearchJob leaked raw source or bearer: %s",secret)
		}
	}
	denied:=invoke("rule_id","5710","wrong-token")
	if denied.Code==http.StatusOK || strings.Contains(denied.Body.String(),"local-jsonl-entry-"){
		t.Fatal("unauthorized mixed Case search exposed source references")
	}
	if err:=os.Remove(wazuh);err!=nil {t.Fatal(err)}
	if err:=os.Remove(generic);err!=nil {t.Fatal(err)}
	if _,err:=loadDemoSnapshotWithExpectedDigest(private,digest);err!=nil {
		t.Fatalf("mixed Case unexpectedly depends on original raw evidence: %v",err)
	}
}

func TestOneStepMixedCaseSaveDoesNotPublishPartialEvidence(t *testing.T) {
	const caseID="forensic-case-on-disk"
	owner:=snapshotFixture(t)
	wazuh:=writeJSONLEvents(t,
		"{\"id\":\"alert\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE Wazuh\"}}\n")
	valid:=writeJSONLEvents(t,
		"{\"id\":\"generic\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"PRIVATE generic\"}\n")
	invalid:=writeJSONLEvents(t,
		"{\"id\":\"broken\",\"timestamp\":\"not-a-time\",\"message\":\"untrusted\"}\n")
	for _,tc:=range []struct{name string; wazuh, generic string}{
		{"bad generic normalization",wazuh,invalid},
		{"duplicate physical input",wazuh,wazuh},
		{"wazuh invalid JSONL",invalid,valid},
		{"owner alias",owner,valid},
	}{
		t.Run(tc.name,func(t *testing.T){
			dir:=t.TempDir()
			dest:=filepath.Join(dir,"must-not-exist.json")
			var out bytes.Buffer
			err:=runDemoCLI(context.Background(),[]string{
				"-principal","analyst-a","-snapshot",owner,"-into-case",caseID,
				"-wazuh-alerts",tc.wazuh,"-import-jsonl",tc.generic,
				"-save-mixed-case-snapshot",dest,
			},&out)
			if err==nil || out.Len()!=0 {t.Fatalf("invalid mixed Case published output: %v %q",err,out.String())}
			if _,err:=os.Lstat(dest);!os.IsNotExist(err){
				t.Fatal("invalid mixed import created a partially populated Case archive")
			}
			files,err:=os.ReadDir(dir)
			if err!=nil || len(files)!=0 {t.Fatalf("failed import left private temporary artifacts: %+v %v",files,err)}
		})
	}
	// Existing archives and evidence must never be overwritten even when
	// a valid mixed Case is presented to the atomic no-replace publisher.
	dest:=filepath.Join(t.TempDir(),"existing-evidence.json")
	const sentinel="EXISTING SENSITIVE EVIDENCE"
	if err:=os.WriteFile(dest,[]byte(sentinel),0600);err!=nil {t.Fatal(err)}
	err:=writePersistedMixedCase(context.Background(),"analyst-a",owner,caseID,
		[]string{wazuh},[]string{valid},dest,&bytes.Buffer{})
	if err==nil {t.Fatal("mixed snapshot overwrote an existing evidence file")}
	got,err:=os.ReadFile(dest)
	if err!=nil || string(got)!=sentinel {t.Fatal("failed mixed save altered original evidence")}
	// CLI must reject mixing this explicitly selected save mode with
	// other export, investigation, Wazuh or listener operations.
	for _,bad:=range [][]string{
		{"-serve-import"},{"-save-wazuh-case-snapshot",filepath.Join(t.TempDir(),"bad")},
		{"-save-import-snapshot",filepath.Join(t.TempDir(),"bad")},
		{"-report-case",caseID},{"-resume-import-search-case",caseID},
	}{
		dest:=filepath.Join(t.TempDir(),"not-created.json")
		args:=append([]string{
			"-principal","analyst-a","-snapshot",owner,"-into-case",caseID,
			"-wazuh-alerts",wazuh,"-import-jsonl",valid,
			"-save-mixed-case-snapshot",dest,
		},bad...)
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&out);err==nil || out.Len()!=0 {
			t.Fatalf("mixed save accepted unrelated mode %v: %q",bad,out.String())
		}
		if _,err:=os.Lstat(dest);!os.IsNotExist(err){
			t.Fatal("rejected CLI mode created a mixed Case snapshot")
		}
	}
}
