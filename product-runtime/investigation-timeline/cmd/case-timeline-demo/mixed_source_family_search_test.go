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
)

func TestMixedCaseSearchJobSelectsAuthorizedSourceFamily(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"w-a\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE WAZUH A\",\"level\":7}}\n"+
		"{\"id\":\"w-b\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE WAZUH B\",\"level\":8}}\n")
	generic := writeJSONLEvents(t,
		"{\"id\":\"g-one\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PRIVATE GENERIC EVENT\"}\n")
	saved:=filepath.Join(t.TempDir(),"mixed-case.json")
	if err:=writePersistedMixedCase(context.Background(),"analyst-a",snapshotFixture(t),caseID,
		[]string{wazuh},[]string{generic},saved,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	raw,err:=os.ReadFile(saved)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(raw))
	data, imported, err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",saved,caseID,digest)
	if err!=nil {t.Fatal(err)}
	if len(imported)!=3 {t.Fatalf("invalid mixed fixture count %d",len(imported))}
	token,err:=randomBearer()
	if err!=nil {t.Fatal(err)}
	base,err:=newDemoHandlerWithData(token,"analyst-a",data)
	if err!=nil {t.Fatal(err)}
	handler:=newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,data,imported,base)
	query:=func(family,field,term,cid,bearer string) *httptest.ResponseRecorder {
		body,_:=json.Marshal(map[string]string{
			"case_id":cid,"source_family":family,"field":field,"value":term,
			"from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:02:00Z",
		})
		req:=httptest.NewRequest(http.MethodPost,"/demo-local-search",bytes.NewReader(body))
		req.Header.Set("Content-Type","application/json")
		req.Header.Set("Authorization","Bearer "+bearer)
		out:=httptest.NewRecorder()
		handler.ServeHTTP(out,req)
		return out
	}
	for _,tc:=range []struct {
		family,field,term string
		count int
		wazuh bool
	}{
		{"all","source_ref","local-jsonl-source-",3,false},
		{"wazuh","source_ref","local-jsonl-source-",2,true},
		{"generic","source_ref","local-jsonl-source-",1,false},
		{"wazuh","rule_id","5710",2,true},
	} {
		t.Run(tc.family+"-"+tc.field,func(t *testing.T){
			response:=query(tc.family,tc.field,tc.term,caseID,token)
			var receipt struct {
				Job search.SearchJob `json:"job"`
			}
			if response.Code!=http.StatusOK || json.Unmarshal(response.Body.Bytes(),&receipt)!=nil ||
				receipt.Job.State!=search.JobCompleted || len(receipt.Job.ResultRefs)!=tc.count {
				t.Fatalf("real authorized family search failed: %d %s",response.Code,response.Body.String())
			}
			seen:=map[string]bool{}
			for _,ref:=range receipt.Job.ResultRefs {
				if ref.Source!=localJSONLSearchSource || seen[ref.Ref] {
					t.Fatalf("foreign/duplicate family search reference: %+v",ref)
				}
				seen[ref.Ref]=true
				var isWazuh bool
				for _,entry:=range data.Entries {
					if entry.Value.ID!=ref.Ref { continue }
					for _,source:=range data.Sources {
						if source.Identity.Reference==entry.Value.SourceRef {isWazuh=source.Wazuh!=nil}
					}
				}
				if tc.family=="wazuh" && !isWazuh ||
					tc.family=="generic" && isWazuh {
					t.Fatalf("cross-family source exposed in result: %+v",ref)
				}
				req:=httptest.NewRequest(http.MethodGet,
					"/demo-source-open?case_id="+caseID+"&entry_id="+ref.Ref,nil)
				req.Header.Set("Authorization","Bearer "+token)
				sourceResponse:=httptest.NewRecorder()
				base.ServeHTTP(sourceResponse,req)
				if sourceResponse.Code!=http.StatusOK {
					t.Fatalf("family result source owner not reopenable: %d",sourceResponse.Code)
				}
			}
			for _,secret:=range []string{"PRIVATE WAZUH","PRIVATE GENERIC","Bearer "+token} {
				if strings.Contains(response.Body.String(),secret) {
					t.Fatalf("source preview or secret disclosed: %q",secret)
				}
			}
		})
	}
	for _,tc:=range []struct{family,field,term,cid,bearer string}{
		{"other","source_ref","local-jsonl-source-",caseID,token},
		{"generic","rule_id","5710",caseID,token},
		{"wazuh","source_ref","local-jsonl-source-","foreign-case",token},
		{"wazuh","source_ref","local-jsonl-source-",caseID,"bad-token"},
	} {
		response:=query(tc.family,tc.field,tc.term,tc.cid,tc.bearer)
		if response.Code==http.StatusOK || strings.Contains(response.Body.String(),"local-jsonl-entry-") {
			t.Fatalf("invalid/unauthorized family search exposed data: %d %s",response.Code,response.Body.String())
		}
	}
}
