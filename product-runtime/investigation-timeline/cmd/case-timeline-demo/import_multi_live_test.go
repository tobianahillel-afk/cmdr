package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMultipleLocalJSONLFilesReachAuthorizedAnalystViewer(t *testing.T) {
	owner := snapshotFixture(t)
	original, err := os.ReadFile(owner)
	if err != nil { t.Fatal(err) }
	fileA := writeJSONLEvents(t, "{\"id\":\"sysmon-1\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"Observed process creation\"}\n")
	fileB := writeJSONLEvents(t, "{\"id\":\"windows-2\",\"timestamp\":\"2026-10-15T10:20:31Z\",\"message\":\"Observed Windows authentication\"}\n")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{
			"-principal=analyst-a", "-snapshot="+owner, "-into-case=forensic-case-on-disk",
			"-import-jsonl="+fileA, "-import-jsonl="+fileB, "-serve-import",
		}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done: t.Fatalf("multi-source listener failed before binding: %v",err)
	case <-time.After(10*time.Second): t.Fatal("multi-source listener not ready")
	}
	var link, token, viewer string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line,"URL: ") {link=strings.TrimPrefix(line,"URL: ")}
		if strings.HasPrefix(line,"Bearer token (this run only): ") {token=strings.TrimPrefix(line,"Bearer token (this run only): ")}
		if strings.HasPrefix(line,"Local analyst viewer: ") {viewer=strings.TrimPrefix(line,"Local analyst viewer: ")}
	}
	u,err:=url.Parse(link)
	if err!=nil || u.Hostname()!="127.0.0.1" || len(token)!=64 ||
		!strings.Contains(viewer,"case_id=forensic-case-on-disk") { t.Fatalf("unsafe local service: %q",link) }
	client:=&http.Client{Timeout:3*time.Second}
	get:=func(path,bearer string)(int,string){
		t.Helper()
		req,err:=http.NewRequest(http.MethodGet,"http://"+u.Host+path,nil)
		if err!=nil {t.Fatal(err)}
		if bearer!="" {req.Header.Set("Authorization","Bearer "+bearer)}
		resp,err:=client.Do(req)
		if err!=nil {t.Fatal(err)}
		defer resp.Body.Close()
		b,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20));if err!=nil {t.Fatal(err)}
		return resp.StatusCode,string(b)
	}
	caseID:="forensic-case-on-disk"
	status,body:=get("/case-timeline?case_id="+caseID+"&limit=200",token)
	if status!=200 || strings.Contains(body,"demo-hidden-source") || strings.Contains(body,"demo-tenant-b"){
		t.Fatalf("unsafe multi-file timeline: %d %s",status,body)
	}
	for _, event:=range []struct{id,message string}{
		{"sysmon-1","Observed process creation"},{"windows-2","Observed Windows authentication"},
	} {
		feed := "feed-1"
		if event.id == "windows-2" {feed = "feed-2"}
		digest:=sha256.Sum256([]byte(caseID+"\x00"+feed+"\x00"+event.id))
		entry:="local-jsonl-entry-"+hex.EncodeToString(digest[:])
		if !strings.Contains(body,entry) {t.Fatalf("imported %s not present in timeline",event.id)}
		p:="/demo-source-open?"+url.Values{"case_id":{caseID},"entry_id":{entry}}.Encode()
		code,preview:=get(p,token)
		if code!=200 || !strings.Contains(preview,event.message) {t.Fatalf("source %s not authorized: %d %s",event.id,code,preview)}
		if code,denied:=get(p,""); code!=http.StatusUnauthorized || strings.Contains(denied,event.message){
			t.Fatalf("anonymous source disclosed: %d %s",code,denied)
		}
	}
	if code,foreign:=get("/case-timeline?case_id=case-b",token);code!=http.StatusNotFound || strings.Contains(foreign,"windows-2"){
		t.Fatalf("foreign Case disclosed: %d %s",code,foreign)
	}
	after,err:=os.ReadFile(owner)
	if err!=nil || !bytes.Equal(original,after) {t.Fatal("operator snapshot was modified")}
	cancel()
	select {
	case err:=<-done: if err!=nil {t.Fatal(err)}
	case <-time.After(5*time.Second):t.Fatal("multi-source local server did not shut down")
	}
}

func TestMultipleLocalJSONLFilesFailClosedBeforeBearerOrListener(t *testing.T){
	owner:=snapshotFixture(t)
	ok:=writeJSONLEvents(t,"{\"id\":\"a\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"event a\"}\n")
	duplicate:=writeJSONLEvents(t,"{\"id\":\"a\",\"timestamp\":\"2026-10-15T10:20:31Z\",\"message\":\"duplicate event\"}\n{\"id\":\"a\",\"timestamp\":\"2026-10-15T10:20:32Z\",\"message\":\"duplicate within same feed\"}\n")
	bad:=writeJSONLEvents(t,"{\"id\":\"b\",\"timestamp\":\"invalid\",\"message\":\"invalid\"}\n")
	five:=[]string{"-serve-import","-snapshot="+owner,"-into-case=forensic-case-on-disk"}
	for i:=0;i<5;i++{five=append(five,"-import-jsonl="+ok)}
	cases:=[]struct{name string; args []string}{
		{"duplicate-event-id",[]string{"-serve-import","-snapshot="+owner,"-into-case=forensic-case-on-disk","-import-jsonl="+ok,"-import-jsonl="+duplicate}},
		{"second-file-invalid",[]string{"-serve-import","-snapshot="+owner,"-into-case=forensic-case-on-disk","-import-jsonl="+ok,"-import-jsonl="+bad}},
		{"same-file",[]string{"-serve-import","-snapshot="+owner,"-into-case=forensic-case-on-disk","-import-jsonl="+ok,"-import-jsonl="+ok}},
		{"owner-as-event-file",[]string{"-serve-import","-snapshot="+owner,"-into-case=forensic-case-on-disk","-import-jsonl="+ok,"-import-jsonl="+owner}},
		{"offline-mode-repeated",[]string{"-snapshot="+owner,"-into-case=forensic-case-on-disk","-import-jsonl="+ok,"-import-jsonl="+duplicate}},
		{"more-than-four",five},
		{"foreign-case",[]string{"-serve-import","-snapshot="+owner,"-into-case=case-b","-import-jsonl="+ok,"-import-jsonl="+duplicate}},
	}
	for _,tc:=range cases {
		t.Run(tc.name,func(t *testing.T){
			var output bytes.Buffer
			if err:=runDemoCLI(context.Background(),tc.args,&output);err==nil || output.Len()!=0 {
				t.Fatalf("untrusted multi-file import emitted bearer or started server: err=%v output=%q",err,output.String())
			}
		})
	}
	alias := filepath.Join(t.TempDir(), "events-hardlink.jsonl")
	if err := os.Link(ok, alias); err == nil {
		var output bytes.Buffer
		err := runDemoCLI(context.Background(), []string{
			"-serve-import", "-snapshot=" + owner, "-into-case=forensic-case-on-disk",
			"-import-jsonl=" + ok, "-import-jsonl=" + alias,
		}, &output)
		if err == nil || output.Len() != 0 {
			t.Fatalf("hardlinked duplicate input was accepted: %v %q", err, output.String())
		}
	}
}
