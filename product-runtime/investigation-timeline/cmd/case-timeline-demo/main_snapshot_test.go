package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func snapshotForLocalInvestigator() timeline.MemoryTimelineData {
	data := demoData()
	data.Principals[0].ID = "forensic-operator"
	for i := range data.Cases {
		for j, id := range data.Cases[i].Readers {
			if id == "analyst-a" { data.Cases[i].Readers[j] = "forensic-operator" }
		}
	}
	for i := range data.Entries {
		for j, id := range data.Entries[i].Readers {
			if id == "analyst-a" { data.Entries[i].Readers[j] = "forensic-operator" }
		}
	}
	for i := range data.Sources {
		for j, id := range data.Sources[i].Readers {
			if id == "analyst-a" { data.Sources[i].Readers[j] = "forensic-operator" }
		}
	}
	data.Cases[0].ID = "case-imported"
	for i := range data.Entries {
		if data.Entries[i].CaseID=="case-a" { data.Entries[i].CaseID = "case-imported" }
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.CaseID=="case-a" { data.Sources[i].Identity.CaseID = "case-imported" }
	}
	return data
}

func TestSnapshotCLIRealHTTPFlowAndIsolation(t *testing.T) {
	original := snapshotForLocalInvestigator()
	payload, err := json.Marshal(original)
	if err != nil {t.Fatal(err)}
	path := filepath.Join(t.TempDir(),"trusted-snapshot.json")
	if err := os.WriteFile(path,payload,0o600); err != nil {t.Fatal(err)}
	ctx,cancel:=context.WithCancel(context.Background())
	defer cancel()
	log:= &demoCapture{ready:make(chan struct{})}
	done:=make(chan error,1)
	go func(){done<-runDemoCLI(ctx,[]string{"-principal=forensic-operator","-snapshot="+path},log)}()
	select {
	case <-log.ready:
	case err:=<-done: t.Fatalf("snapshot listener failed: %v",err)
	case <-time.After(8*time.Second):t.Fatal("snapshot server not ready")
	}
	var endpoint, bearer string
	for _, line:=range strings.Split(log.String(),"\n") {
		if strings.HasPrefix(line,"URL: ") {endpoint=strings.TrimPrefix(line,"URL: ")}
		if strings.HasPrefix(line,"Bearer token (this run only): ") {
			bearer=strings.TrimPrefix(line,"Bearer token (this run only): ")
		}
	}
	parsed,err:=url.Parse(endpoint)
	if err!=nil || parsed.Hostname()!="127.0.0.1" || parsed.Query().Get("case_id")!="case-imported" || len(bearer)!=64 {
		t.Fatalf("unsafe or incorrect local route %q bearer length=%d err=%v",endpoint,len(bearer),err)
	}
	req,err:=http.NewRequest(http.MethodGet,endpoint,nil)
	if err!=nil{t.Fatal(err)}
	req.Header.Set("Authorization","Bearer "+bearer)
	client:=&http.Client{Timeout:2*time.Second}
	response,err:=client.Do(req)
	if err!=nil{t.Fatal(err)}
	body,err:=io.ReadAll(io.LimitReader(response.Body,1<<20))
	response.Body.Close()
	if err!=nil || response.StatusCode!=http.StatusOK {t.Fatalf("authorized local snapshot request %d %v",response.StatusCode,err)}
	if !bytes.Contains(body,[]byte("event-new")) || !bytes.Contains(body,[]byte("event-old")) {
		t.Fatalf("imported Timeline Entries missing: %s",body)
	}
	for _, secret:=range []string{"restricted","demo-secret","demo-hidden-source","event-b","demo-tenant-b"} {
		if bytes.Contains(body,[]byte(secret)){t.Fatalf("hidden/foreign source leaked: %s",secret)}
	}
	foreign:=strings.Replace(endpoint,"case_id=case-imported","case_id=case-b",1)
	r,err:=http.NewRequest(http.MethodGet,foreign,nil)
	if err!=nil{t.Fatal(err)}
	r.Header.Set("Authorization","Bearer "+bearer)
	resp,err:=client.Do(r)
	if err!=nil{t.Fatal(err)}
	resp.Body.Close()
	if resp.StatusCode!=http.StatusNotFound {t.Fatalf("foreign tenant Case disclosed: HTTP %d",resp.StatusCode)}
	anonymous,err:=client.Get(endpoint)
	if err!=nil{t.Fatal(err)}
	anonymous.Body.Close()
	if anonymous.StatusCode!=http.StatusUnauthorized {t.Fatalf("anonymous snapshot readable HTTP %d",anonymous.StatusCode)}
	cancel()
	select{
	case err:=<-done:if err!=nil{t.Fatal(err)}
	case <-time.After(5*time.Second):t.Fatal("snapshot listener did not shut down")
	}
}

func TestSnapshotCLIRefusesMalformedOrUnauthorizedInput(t *testing.T) {
	valid,err:=json.Marshal(snapshotForLocalInvestigator())
	if err!=nil{t.Fatal(err)}
	bad:=[]struct{name string;content string;args []string}{
		{"corrupt-json","{bad",[]string{"-principal=forensic-operator"}},
		{"unknown-principal",string(valid),[]string{"-principal=not-authorized"}},
		{"oversized",strings.Repeat("x",timeline.MaxTrustedSnapshotBytes+1),[]string{"-principal=forensic-operator"}},
	}
	for _,tc:=range bad{
		t.Run(tc.name,func(t *testing.T){
			p:=filepath.Join(t.TempDir(),"input.json")
			if err:=os.WriteFile(p,[]byte(tc.content),0o600);err!=nil{t.Fatal(err)}
			ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
			defer cancel()
			var output bytes.Buffer
			args:=append(append([]string(nil),tc.args...),"-snapshot="+p)
			if err:=runDemoCLI(ctx,args,&output);err==nil || output.Len()!=0 {
				t.Fatalf("invalid snapshot activated listener: err=%v bytes=%d",err,output.Len())
			}
		})
	}
	var output bytes.Buffer
	if err:=runDemoCLI(context.Background(),[]string{"-snapshot=/unavailable/snapshot.json"},&output);err==nil || output.Len()!=0{
		t.Fatalf("missing snapshot opened: %v",err)
	}
	if err:=runDemoCLI(context.Background(),[]string{"-snapshot=/unavailable","unexpected"},&output);err==nil || output.Len()!=0 {
		t.Fatalf("invalid CLI options accepted: %v",err)
	}
}

func TestSnapshotLocalAuthorityRejectsInconsistentOwnerOrUnknownPrincipal(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*timeline.MemoryTimelineData)
		principal string
	}{
		{"no-such-reader",func(*timeline.MemoryTimelineData){},"unknown"},
		{"reader-without-case",func(d *timeline.MemoryTimelineData){
			d.Cases[0].Readers=nil
		},"forensic-operator"},
		{"source-owned-by-other-case",func(d *timeline.MemoryTimelineData){
			d.Sources[0].Identity.CaseID="case-b"
		},"forensic-operator"},
		{"entry-foreign-tenant",func(d *timeline.MemoryTimelineData){
			d.Entries[0].Value.Tenant="demo-tenant-b"
		},"forensic-operator"},
		{"case-reader-not-in-tenant",func(d *timeline.MemoryTimelineData){
			d.Cases[0].Readers=[]string{"analyst-b"}
		},"forensic-operator"},
	} {
		t.Run(tc.name,func(t *testing.T){
			data:=snapshotForLocalInvestigator()
			tc.edit(&data)
			var output bytes.Buffer
			ctx,cancel:=context.WithTimeout(context.Background(),2*time.Second)
			defer cancel()
			if err:=serveDemoWithData(ctx,tc.principal,&output,data);err==nil || output.Len()!=0 {
				t.Fatalf("inconsistent offline owner data opened a listener: %v",err)
			}
		})
	}
}

func TestSnapshotCannotTurnUntrustedSourceFlagIntoAuthority(t *testing.T) {
	data:=snapshotForLocalInvestigator()
	data.Entries[2].Value.SourcePermission="allow"
	handler,err:=newDemoHandlerWithData(testBearer,"forensic-operator",data)
	if err!=nil{t.Fatal(err)}
	server:=httptest.NewServer(handler)
	defer server.Close()
	request,err:=http.NewRequest(http.MethodGet,server.URL+"/case-timeline?case_id=case-imported",nil)
	if err!=nil{t.Fatal(err)}
	request.Header.Set("Authorization","Bearer "+testBearer)
	response,err:=server.Client().Do(request)
	if err!=nil{t.Fatal(err)}
	defer response.Body.Close()
	body,err:=io.ReadAll(io.LimitReader(response.Body,1<<20))
	if err!=nil{t.Fatal(err)}
	if response.StatusCode!=http.StatusOK || !bytes.Contains(body,[]byte("event-old")) ||
		bytes.Contains(body,[]byte("restricted")) || bytes.Contains(body,[]byte("demo-hidden-source")){
		t.Fatalf("untrusted snapshot source claim bypassed owner policy: %d %s",response.StatusCode,body)
	}
}
