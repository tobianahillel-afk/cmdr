package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
	"net/http"
	"net/http/httptest"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestPersistedNormalizedMultiFeedReopensAuthorizedSearch(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot := snapshotFixture(t)
	feedA := writeJSONLEvents(t,
		"{\"id\":\"same-event\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"PowerShell first source\"}\n"+
		"{\"id\":\"routine\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"message\":\"benign maintenance\"}\n")
	feedB := writeJSONLEvents(t,
		"{\"id\":\"same-event\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"message\":\"POWERSHELL second source\"}\n")
	saved := filepath.Join(t.TempDir(), "private-case.json")
	var receipt bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", snapshot, "-into-case", caseID,
		"-import-jsonl", feedA, "-import-jsonl-extra", feedB,
		"-save-import-snapshot", saved,
	}, &receipt); err != nil {
		t.Fatalf("explicit private snapshot export failed: %v", err)
	}
	if !strings.Contains(receipt.String(), "sha256") ||
		strings.Contains(receipt.String(), "PowerShell first source") ||
		strings.Contains(receipt.String(), "POWERSHELL second source") ||
		strings.Contains(receipt.String(), "Bearer ") {
		t.Fatalf("private owner save receipt leaks source data or lacks digest: %q", receipt.String())
	}
	stat, err := os.Stat(saved)
	if err != nil {t.Fatal(err)}
	if !stat.Mode().IsRegular() || stat.Size() == 0 {
		t.Fatalf("saved owner snapshot not a regular file: %v",stat.Mode())
	}
	if runtime.GOOS != "windows" && stat.Mode().Perm() != 0600 {
		t.Fatalf("private snapshot must be 0600 on POSIX: %v",stat.Mode().Perm())
	}
	data,err:=loadDemoSnapshot(saved)
	if err!=nil {t.Fatalf("saved owner snapshot does not load: %v",err)}
	ids,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(ids)!=3 {
		t.Fatalf("reopened Case did not identify exactly the imported events: %d, %v",len(ids),err)
	}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	owner,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,ids,owner)
	body:=[]byte("{\"case_id\":\""+caseID+"\",\"value\":\"powershell\",\"from\":\"2026-10-17T09:59:00Z\",\"to\":\"2026-10-17T10:03:00Z\"}")
	request:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
	request.Header.Set("Authorization","Bearer "+token)
	request.Header.Set("Content-Type","application/json")
	response:=httptest.NewRecorder()
	handler.ServeHTTP(response,request)
	if response.Code!=http.StatusOK {
		t.Fatalf("reopened first-party Event Search Job rejected: %d %s",response.Code,response.Body.String())
	}
	var found struct {Job search.SearchJob `json:"job"`}
	if err:=json.Unmarshal(response.Body.Bytes(),&found);err!=nil {t.Fatal(err)}
	if found.Job.State!=search.JobCompleted || len(found.Job.ResultRefs)!=2 {
		t.Fatalf("expected two persisted cross-feed matches with same original ID: %+v",found.Job)
	}
	seen:=map[string]bool{}
	for _,ref:=range found.Job.ResultRefs {
		if ref.Source!=localJSONLSearchSource {t.Fatalf("unexpected persisted source: %+v",ref)}
		if !strings.HasPrefix(ref.Ref,"local-jsonl-entry-") {t.Fatalf("unexpected ref: %+v",ref)}
		if _,exists:=ids[ref.Ref];!exists {t.Fatalf("result not in authorized persisted Case: %+v",ref)}
		seen[ref.Ref]=true
		source:=httptest.NewRequest(http.MethodGet,"/demo-source-open?case_id="+caseID+"&entry_id="+ref.Ref,nil)
		source.Header.Set("Authorization","Bearer "+token)
		open:=httptest.NewRecorder()
		owner.ServeHTTP(open,source)
		if open.Code!=http.StatusOK {t.Fatalf("persisted source could not be reopened independently: %d",open.Code)}
	}
	if len(seen)!=2 {t.Fatal("duplicate original alert IDs collapsed after snapshot reload")}
	for _,secret:=range []string{"PowerShell first source","POWERSHELL second source","benign maintenance",token} {
		if strings.Contains(response.Body.String(),secret) {t.Fatalf("SearchJob disclosed sensitive preview or bearer: %q",secret)}
	}
	bad:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
	bad.Header.Set("Authorization","Bearer wrong")
	bad.Header.Set("Content-Type","application/json")
	denied:=httptest.NewRecorder()
	handler.ServeHTTP(denied,bad)
	if denied.Code==http.StatusOK {t.Fatal("reopened owner accepted former or incorrect bearer")}
	// Reopening from saved owner must not require original selected JSONL
	// files to remain on disk: only the bounded private snapshot is used.
	if err:=os.Remove(feedA);err!=nil {t.Fatal(err)}
	if err:=os.Remove(feedB);err!=nil {t.Fatal(err)}
	if _,err:=loadDemoSnapshot(saved);err!=nil {
		t.Fatalf("persisted owner unexpectedly depends on raw input feeds: %v",err)
	}
	// Exercise the actual reopen CLI and listener lifecycle after deleting
	// original input files, not only its internal owner/search helpers.
	resumeContext, stop := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer stop()
	var viewer bytes.Buffer
	if err:=runDemoCLI(resumeContext, []string{
		"-principal", "analyst-a", "-snapshot", saved,
		"-resume-import-search-case", caseID,
	}, &viewer);err!=nil {
		t.Fatalf("persisted search viewer CLI failed: %v",err)
	}
	if !strings.Contains(viewer.String(), "Local analyst viewer: http://127.0.0.1:") ||
		!strings.Contains(viewer.String(), "Bearer token (this run only):") ||
		strings.Contains(viewer.String(), "PowerShell first source") {
		t.Fatalf("reopened viewer failed to start independently or leaked payload: %q",viewer.String())
	}
}

func TestPersistedNormalizedSnapshotNeverOverwritesOrLeaksOnFailure(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot:=snapshotFixture(t)
	feed:=writeJSONLEvents(t,
		"{\"id\":\"only\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"sensitive private message\"}\n")
	dir:=t.TempDir()
	dest:=filepath.Join(dir,"case.json")
	sentinel:=[]byte("existing evidence file must not be overwritten")
	if err:=os.WriteFile(dest,sentinel,0600);err!=nil {t.Fatal(err)}
	var out bytes.Buffer
	err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",snapshot,caseID,[]string{feed},dest,&out)
	if err==nil {t.Fatal("existing evidence file was accepted as output target")}
	after,err:=os.ReadFile(dest)
	if err!=nil || !bytes.Equal(after,sentinel) {t.Fatal("existing evidence file changed")}
	if out.Len()!=0 {t.Fatal("failed snapshot save disclosed a receipt")}
	if err:=os.Remove(dest);err!=nil {t.Fatal(err)}
	// A dangling symlink must also fail rather than being followed or replaced.
	link:=filepath.Join(dir,"symlink-case.json")
	if err:=os.Symlink("does-not-exist.json",link);err==nil {
		if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",snapshot,caseID,[]string{feed},link,&out);err==nil {
			t.Fatal("symlink target accepted for evidence")
		}
	}
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",snapshot,caseID,[]string{feed},"relative.json",&out);err==nil {
		t.Fatal("relative evidence output path accepted")
	}
	cancelled,cancel:=context.WithCancel(context.Background())
	cancel()
	if err:=writeDemoPersistedNormalizedCase(cancelled,"analyst-a",snapshot,caseID,[]string{feed},
		filepath.Join(dir,"cancelled.json"),&out);err==nil {
		t.Fatal("cancelled save accepted")
	}
	if _,err:=os.Stat(filepath.Join(dir,"cancelled.json"));!os.IsNotExist(err) {
		t.Fatal("cancelled save created a destination")
	}
	for _,mode:=range [][]string{
		{"-snapshot",snapshot,"-import-jsonl",feed,"-into-case",caseID,
		 "-save-import-snapshot",filepath.Join(dir,"invalid.json"),"-serve-import"},
		{"-snapshot",snapshot,"-import-jsonl",feed,"-into-case",caseID,
		 "-save-import-snapshot",filepath.Join(dir,"invalid.json"),"-report-case",caseID},
		{"-snapshot",snapshot,"-resume-import-search-case",caseID,"-import-jsonl",feed},
		{"-snapshot",snapshot,"-resume-import-search-case",caseID,"-save-import-snapshot",filepath.Join(dir,"invalid.json")},
	} {
		var result bytes.Buffer
		if err:=runDemoCLI(context.Background(),mode,&result);err==nil || result.Len()!=0 {
			t.Fatalf("invalid combined persistence modes accepted or printed private data: %v",mode)
		}
	}
}

func TestReopenedNormalizedEligibilityNeverWidensToOtherCases(t *testing.T) {
	const caseID="forensic-case-on-disk"
	snapshot:=snapshotFixture(t)
	file:=writeJSONLEvents(t,
		"{\"id\":\"owned\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"authenticated\"}\n")
	dest:=filepath.Join(t.TempDir(),"owner.json")
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",snapshot,caseID,[]string{file},dest,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	data,err:=loadDemoSnapshot(dest)
	if err!=nil {t.Fatal(err)}
	if _,err:=importedNormalizedIDs(data,"different-case");err==nil {
		t.Fatal("reopening a foreign Case widened imported search scope")
	}
	idset,err:=importedNormalizedIDs(data,caseID)
	if err!=nil || len(idset)!=1 {t.Fatalf("reopened selected Case missing import: %v",err)}
	for id:=range idset {
		for i:=range data.Entries {
			if data.Entries[i].Value.ID==id {
				data.Entries[i].Value.SourceRef="some-other-owner-source"
			}
		}
	}
	if _,err:=importedNormalizedIDs(data,caseID);err==nil {
		t.Fatal("inconsistent owner source identity accepted after tampering")
	}
}
