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

func makeSavedPairDrilldownOwners(t *testing.T) (string,string,string,string,string,string,sharedAuthorizedWazuhEventPair) {
	t.Helper()
	aOwner,aCase,_,bOwner,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	firstFile:=writeJSONLEvents(t,
		"{\"id\":\"a-pair\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE FIRST OWNER PAIR\",\"level\":9}}\n")
	secondFile:=writeJSONLEvents(t,
		"{\"id\":\"b-pair\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE SECOND OWNER PAIR\",\"level\":9}}\n")
	firstPath:=filepath.Join(t.TempDir(),"first-saved-case.json")
	secondPath:=filepath.Join(t.TempDir(),"second-saved-case.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",aOwner,aCase,[]string{firstFile},firstPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",bOwner,bCase,[]string{secondFile},secondPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	aBytes,err:=os.ReadFile(firstPath)
	if err!=nil {t.Fatal(err)}
	bBytes,err:=os.ReadFile(secondPath)
	if err!=nil {t.Fatal(err)}
	aHash:=fmt.Sprintf("%x",sha256.Sum256(aBytes))
	bHash:=fmt.Sprintf("%x",sha256.Sum256(bBytes))
	if err:=os.Remove(firstFile);err!=nil {t.Fatal(err)}
	if err:=os.Remove(secondFile);err!=nil {t.Fatal(err)}
	var output bytes.Buffer
	if err:=compareVerifiedSavedCases(context.Background(),"analyst-a",firstPath,aCase,aHash,secondPath,bCase,bHash,&output);err!=nil {t.Fatal(err)}
	var result localCaseComparison
	if err:=json.Unmarshal(output.Bytes(),&result);err!=nil {t.Fatal(err)}
	if len(result.SharedWazuhEventPairs)!=1 {t.Fatalf("expected one correlated independently authorized Wazuh event pair: %+v",result.SharedWazuhEventPairs)}
	return firstPath,aCase,aHash,secondPath,bCase,bHash,result.SharedWazuhEventPairs[0]
}

func TestPrivateCrossCasePairOpensBothActualOwnerSourcesIndependently(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash,pair:=makeSavedPairDrilldownOwners(t)
	if pair.FirstEntryID=="" || pair.SecondEntryID=="" || pair.OffsetSeconds!=20 {t.Fatalf("invalid pair: %+v",pair)}
	for _,side:=range []struct{name,path,caseID,hash,entryID,ownSecret,foreignSecret string}{
		{"first",aPath,aCase,aHash,pair.FirstEntryID,"PRIVATE FIRST OWNER PAIR","PRIVATE SECOND OWNER PAIR"},
		{"second",bPath,bCase,bHash,pair.SecondEntryID,"PRIVATE SECOND OWNER PAIR","PRIVATE FIRST OWNER PAIR"},
	} {
		t.Run(side.name,func(t *testing.T){
			destination:=filepath.Join(t.TempDir(),"authorized-source.json")
			var stdout bytes.Buffer
			if err:=runDemoCLI(context.Background(),[]string{
				"-principal","analyst-a","-snapshot",side.path,
				"-inspect-saved-case",side.caseID,
				"-inspect-saved-entry",side.entryID,
				"-expect-snapshot-sha256",side.hash,
				"-private-report",destination,
			},&stdout);err!=nil {t.Fatalf("authorized pair source private open failed: %v",err)}
			if stdout.Len()!=0 {t.Fatalf("sensitive private receipt leaked to stdout: %q",stdout.String())}
			info,err:=os.Stat(destination)
			if err!=nil || !info.Mode().IsRegular() {t.Fatal("private source report not produced")}
			data,err:=os.ReadFile(destination)
			if err!=nil {t.Fatal(err)}
			var inspected localEntryInspection
			if err:=json.Unmarshal(data,&inspected);err!=nil {t.Fatal(err)}
			if inspected.CaseID!=side.caseID || inspected.EntryID!=side.entryID ||
				!strings.Contains(inspected.Preview,side.ownSecret) ||
				strings.Contains(string(data),side.foreignSecret) ||
				strings.Contains(string(data),"Bearer ") {
				t.Fatalf("pair inspected outside authorized source boundary: %+v",inspected)
			}
			// No overwrite or temp leftovers when the operator repeats the
			// exact command for an existing evidence file.
			var repeat bytes.Buffer
			err=runDemoCLI(context.Background(),[]string{
				"-principal","analyst-a","-snapshot",side.path,
				"-inspect-saved-case",side.caseID,
				"-inspect-saved-entry",side.entryID,
				"-expect-snapshot-sha256",side.hash,
				"-private-report",destination,
			},&repeat)
			if err==nil || repeat.Len()!=0 {t.Fatal("existing private source report overwritten or echoed")}
			again,err:=os.ReadFile(destination)
			if err!=nil || !bytes.Equal(again,data) {t.Fatal("evidence file changed after denied repeat")}
		})
	}
}

func TestCrossCaseEventIDNeverGrantsForeignSnapshotAccess(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash,pair:=makeSavedPairDrilldownOwners(t)
	for _,tc:=range []struct{name,principal,path,caseID,hash,entryID string}{
		{"first ID on second owner","analyst-a",bPath,bCase,bHash,pair.FirstEntryID},
		{"second ID on first owner","analyst-a",aPath,aCase,aHash,pair.SecondEntryID},
		{"wrong Case","analyst-a",aPath,bCase,aHash,pair.FirstEntryID},
		{"wrong principal","foreign-analyst",aPath,aCase,aHash,pair.FirstEntryID},
		{"wrong digest","analyst-a",aPath,aCase,strings.Repeat("0",64),pair.FirstEntryID},
		{"invalid digest","analyst-a",aPath,aCase,"not-a-sha256",pair.FirstEntryID},
	} {
		t.Run(tc.name,func(t *testing.T){
			dest:=filepath.Join(t.TempDir(),"should-not-exist.json")
			var stdout bytes.Buffer
			err:=runDemoCLI(context.Background(),[]string{
				"-principal",tc.principal,"-snapshot",tc.path,
				"-inspect-saved-case",tc.caseID,
				"-inspect-saved-entry",tc.entryID,
				"-expect-snapshot-sha256",tc.hash,
				"-private-report",dest,
			},&stdout)
			if err==nil || stdout.Len()!=0 {t.Fatalf("unauthorized entry succeeded or disclosed source: %v %q",err,stdout.String())}
			if _,err:=os.Stat(dest);!os.IsNotExist(err) {t.Fatalf("failed source open left output: %v",err)}
		})
	}
}

func TestPrivateSavedEntryRejectsUnsafeCombinedCLIAndInvalidDestination(t *testing.T) {
	aPath,aCase,aHash,_,_,_,pair:=makeSavedPairDrilldownOwners(t)
	good:=[]string{"-principal","analyst-a","-snapshot",aPath,
		"-inspect-saved-case",aCase,"-inspect-saved-entry",pair.FirstEntryID,
		"-expect-snapshot-sha256",aHash}
	for _,extra:=range [][]string{
		{"-compare-saved-case",aCase},
		{"-verify-import-snapshot-case",aCase},
		{"-serve-import"},
		{"-private-report",""},
		{"-private-report","relative.json"},
	} {
		args:=append(append([]string{},good...),extra...)
		if len(extra)==0 || extra[0]!="-private-report" {
			args=append(args,"-private-report",filepath.Join(t.TempDir(),"denied.json"))
		}
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&out);err==nil || out.Len()!=0 {
			t.Fatalf("unsafe combined/source export mode accepted: %v %q",args,out.String())
		}
	}
	// Absent -private-report must never downgrade a sensitive preview to stdout.
	var out bytes.Buffer
	if err:=runDemoCLI(context.Background(),good,&out);err==nil || out.Len()!=0 {
		t.Fatal("missing private output path leaked source into CLI stdout")
	}
}
