package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrivatePairBundleReopensBothIndependentlyAuthorizedSources(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash,pair:=makeSavedPairDrilldownOwners(t)
	dest:=filepath.Join(t.TempDir(),"full-pair-private.json")
	args:=[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-inspect-saved-pair","1","-private-report",dest,
	}
	var output bytes.Buffer
	if err:=runDemoCLI(context.Background(),args,&output);err!=nil {
		t.Fatalf("cannot privately inspect actual Wazuh correlated pair: %v",err)
	}
	if output.Len()!=0 {t.Fatalf("source previews leaked to stdout: %q",output.String())}
	file,err:=os.ReadFile(dest)
	if err!=nil {t.Fatal(err)}
	var bundle struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		First localEntryInspection `json:"first"`
		Second localEntryInspection `json:"second"`
	}
	if err:=json.Unmarshal(file,&bundle);err!=nil {t.Fatalf("invalid private bundle JSON: %v",err)}
	if bundle.Format!="cmdr-private-verified-wazuh-pair-v1" || !bundle.Synthetic ||
		bundle.Pair.FirstEntryID!=pair.FirstEntryID || bundle.Pair.SecondEntryID!=pair.SecondEntryID ||
		bundle.Pair.OffsetSeconds!=20 || bundle.Pair.RuleID!="5710" || bundle.Pair.AgentID!="001" ||
		bundle.First.CaseID!=aCase || bundle.First.EntryID!=pair.FirstEntryID ||
		bundle.Second.CaseID!=bCase || bundle.Second.EntryID!=pair.SecondEntryID ||
		!strings.Contains(bundle.First.Preview,"PRIVATE FIRST OWNER PAIR") ||
		!strings.Contains(bundle.Second.Preview,"PRIVATE SECOND OWNER PAIR") {
		t.Fatalf("private bundle did not retain both real authorized source previews and correlation: %+v",bundle)
	}
	if strings.Contains(string(file),"Bearer ") || strings.Contains(string(file),"Local analyst viewer:") {
		t.Fatal("private paired report emitted session credentials or loopback listener")
	}
	info,err:=os.Stat(dest)
	if err!=nil || !info.Mode().IsRegular() {t.Fatal("private pair report missing")}
	if _,err:=os.Stat(dest);err==nil && info.Mode().Perm()!=0600 && os.PathSeparator=='/' {
		t.Fatalf("private source bundle should be 0600 on POSIX: %s",info.Mode().Perm())
	}
	// Existing evidence may never be replaced, including a repeat of the
	// same fully authorized pair and its independently held fingerprints.
	output.Reset()
	if err:=runDemoCLI(context.Background(),args,&output);err==nil || output.Len()!=0 {
		t.Fatal("repeat of private pair report overwrote existing artifact or leaked stdout")
	}
	again,err:=os.ReadFile(dest)
	if err!=nil || !bytes.Equal(again,file) {t.Fatal("existing private paired evidence was modified")}
}

func TestPrivatePairBundleNoPartialOutputForDeniedSecondSource(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash,_:=makeSavedPairDrilldownOwners(t)
	base:=[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-inspect-saved-pair","1",
	}
	for _,test:=range []struct{name string; override []string}{
		{"bad second owner fingerprint",[]string{"-expect-other-snapshot-sha256",strings.Repeat("0",64)}},
		{"bad first owner fingerprint",[]string{"-expect-snapshot-sha256",strings.Repeat("0",64)}},
		{"missing selected pair",[]string{"-inspect-saved-pair","2"}},
		{"zero selected pair",[]string{"-inspect-saved-pair","0"}},
		{"invalid selected pair",[]string{"-inspect-saved-pair","11"}},
		{"foreign source principal",[]string{"-principal","foreign-analyst"}},
		{"other Case missing",[]string{"-compare-saved-with-case","different-case"}},
		{"unauthorized combination",[]string{"-serve-import"}},
	} {
		t.Run(test.name,func(t *testing.T){
			// Replace, don't double-specify, all overrides.
			args:=append([]string(nil),base...)
			for i:=0;i<len(test.override);i++ {
				key:=test.override[i]
				if key=="-serve-import" {args=append(args,key);continue}
				val:=test.override[i+1];i++
				for pos:=0;pos<len(args)-1;pos++ {if args[pos]==key {args[pos+1]=val;break}}
			}
			dest:=filepath.Join(t.TempDir(),"must-not-exist.json")
			args=append(args,"-private-report",dest)
			var stdout bytes.Buffer
			if err:=runDemoCLI(context.Background(),args,&stdout);err==nil || stdout.Len()!=0 {
				t.Fatalf("denied source produced private bundle or stdout: %v %q",err,stdout.String())
			}
			if _,err:=os.Stat(dest);!os.IsNotExist(err) {t.Fatalf("denied result created private artifact: %v",err)}
		})
	}
	// A first-party Case owner authorized for both files will still fail
	// before file publication if the SECOND snapshot changes after save.
	original,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,append(original,'\n'),0600);err!=nil {t.Fatal(err)}
	dest:=filepath.Join(t.TempDir(),"changed-second.json")
	args:=append(append([]string(nil),base...),"-private-report",dest)
	var stdout bytes.Buffer
	if err:=runDemoCLI(context.Background(),args,&stdout);err==nil || stdout.Len()!=0 {
		t.Fatal("modified second owner allowed one-sided report")
	}
	if _,err:=os.Stat(dest);!os.IsNotExist(err) {t.Fatalf("mismatched SHA created partial private report: %v",err)}
}

func TestPrivatePairBundleCannotCrossAuthorizedTenantBoundary(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,_,_:=makeSavedPairDrilldownOwners(t)
	other,err:=loadDemoSnapshot(bPath)
	if err!=nil {t.Fatal(err)}
	other.Principals[0].Tenant="independent-foreign-tenant"
	other.Cases[0].Tenant="independent-foreign-tenant"
	for i:=range other.Entries {other.Entries[i].Value.Tenant="independent-foreign-tenant"}
	for i:=range other.Sources {other.Sources[i].Identity.Tenant="independent-foreign-tenant"}
	if _,err:=newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",other);err!=nil {
		t.Fatalf("other-tenant Case must be internally owner-valid to test cross-tenant refusal: %v",err)
	}
	raw,err:=json.Marshal(other)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,raw,0600);err!=nil {t.Fatal(err)}
	bHash:=fmt.Sprintf("%x",sha256.Sum256(raw))
	dest:=filepath.Join(t.TempDir(),"denied-tenant.json")
	var stdout bytes.Buffer
	err=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-inspect-saved-pair","1","-private-report",dest,
	},&stdout)
	if err==nil || stdout.Len()!=0 {t.Fatal("valid second tenant snapshot bypassed first tenant owner isolation")}
	if _,err:=os.Lstat(dest);!os.IsNotExist(err) {t.Fatal("denied cross-tenant evidence bundle created")}
}
