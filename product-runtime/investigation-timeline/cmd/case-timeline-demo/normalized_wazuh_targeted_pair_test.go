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
	"time"
)

func TestSelectedWazuhPairBeyondTenIsPrivatelyInspectable(t *testing.T) {
	aOwner,aCase,_,bOwner,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	var firstFeed strings.Builder
	for i:=0;i<12;i++ {
		fmt.Fprintf(&firstFeed,"{\"id\":\"a-%02d\",\"timestamp\":\"2026-10-17T10:00:%02dZ\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE TARGET %02d\",\"level\":9}}\n",i,i,i)
	}
	aFeed:=writeJSONLEvents(t,firstFeed.String())
	bFeed:=writeJSONLEvents(t,"{\"id\":\"b-20\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE SECOND TARGET\",\"level\":8}}\n")
	aPath:=filepath.Join(t.TempDir(),"first.json")
	bPath:=filepath.Join(t.TempDir(),"second.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",aOwner,aCase,[]string{aFeed},aPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",bOwner,bCase,[]string{bFeed},bPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	aBytes,err:=os.ReadFile(aPath);if err!=nil{t.Fatal(err)}
	bBytes,err:=os.ReadFile(bPath);if err!=nil{t.Fatal(err)}
	aHash:=fmt.Sprintf("%x",sha256.Sum256(aBytes))
	bHash:=fmt.Sprintf("%x",sha256.Sum256(bBytes))
	aData,_,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",aPath,aCase,aHash);if err!=nil{t.Fatal(err)}
	bData,_,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",bPath,bCase,bHash);if err!=nil{t.Fatal(err)}
	aGroups,err:=collectAuthorizedWazuhObservations(aData,aCase);if err!=nil{t.Fatal(err)}
	bGroups,err:=collectAuthorizedWazuhObservations(bData,bCase);if err!=nil{t.Fatal(err)}
	key:=authorizedWazuhCohort{rule:"5710",agent:"001"}
	if len(aGroups[key])!=12||len(bGroups[key])!=1{t.Fatal("expected twelve and one source alerts")}
	firstID:=aGroups[key][10].id
	secondID:=bGroups[key][0].id
	foundPair,found:=selectedAuthorizedWazuhPairByEntryIDs(aGroups,bGroups,firstID,secondID)
	if !found||foundPair.OffsetSeconds!=10{t.Fatalf("selected beyond-ten nearest pair invalid: %+v",foundPair)}
	var summary bytes.Buffer
	if err:=compareVerifiedSavedCases(context.Background(),"analyst-a",aPath,aCase,aHash,bPath,bCase,bHash,&summary);err!=nil{t.Fatal(err)}
	var receipt localCaseComparison
	if err:=json.Unmarshal(summary.Bytes(),&receipt);err!=nil{t.Fatal(err)}
	if len(receipt.SharedWazuhEventPairs)!=10{t.Fatal("summary limit changed")}
	for _,p:=range receipt.SharedWazuhEventPairs{if p.FirstEntryID==firstID{t.Fatal("target is inside first ten")}}
	args:=[]string{"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-inspect-saved-pair-first-entry",firstID,
		"-inspect-saved-pair-second-entry",secondID}
	dest:=filepath.Join(t.TempDir(),"private.json")
	var output bytes.Buffer
	if err:=runDemoCLI(context.Background(),append(append([]string{},args...),"-private-report",dest),&output);err!=nil{t.Fatal(err)}
	if output.Len()!=0{t.Fatal("sensitive evidence leaked on stdout")}
	report,err:=os.ReadFile(dest);if err!=nil{t.Fatal(err)}
	if !bytes.Contains(report,[]byte("PRIVATE TARGET 10"))||
		!bytes.Contains(report,[]byte("PRIVATE SECOND TARGET"))||
		bytes.Contains(report,[]byte("Bearer ")){t.Fatal("targeted private evidence not preserved safely")}
	var doc struct{Format string `json:"format"`; Pair sharedAuthorizedWazuhEventPair `json:"pair"`}
	if json.Unmarshal(report,&doc)!=nil||doc.Format!="cmdr-private-verified-wazuh-pair-v2"||
		doc.Pair.FirstEntryID!=firstID||doc.Pair.SecondEntryID!=secondID{t.Fatal("wrong source evidence identities")}
	if os.PathSeparator=='/'{st,err:=os.Stat(dest);if err!=nil||st.Mode().Perm()!=0600{t.Fatal("private report mode not 0600")}}
	output.Reset()
	if err:=runDemoCLI(context.Background(),append(append([]string{},args...),"-private-report",dest),&output);err==nil||output.Len()!=0{t.Fatal("existing artifact overwritten or leaked")}
	again,err:=os.ReadFile(dest);if err!=nil||!bytes.Equal(report,again){t.Fatal("private report changed")}
	for _,tc:=range []struct{name,key,value string}{
		{"wrong peer","-inspect-saved-pair-second-entry","other-id"},
		{"wrong first","-inspect-saved-pair-first-entry","foreign-id"},
		{"wrong digest","-expect-other-snapshot-sha256",strings.Repeat("0",64)},
		{"foreign actor","-principal","foreign-analyst"},
	}{
		t.Run(tc.name,func(t *testing.T){
			x:=append([]string{},args...)
			for i:=0;i+1<len(x);i++{if x[i]==tc.key{x[i+1]=tc.value;break}}
			p:=filepath.Join(t.TempDir(),"denied.json")
			var out bytes.Buffer
			if err:=runDemoCLI(context.Background(),append(x,"-private-report",p),&out);err==nil||out.Len()!=0{t.Fatal("unauthorized private dossier")}
			if _,err:=os.Lstat(p);!os.IsNotExist(err){t.Fatal("denied dossier created")}
		})
	}
}

func TestExactWazuhPairMaintainsNearestCohortSemantics(t *testing.T){
	key:=authorizedWazuhCohort{rule:"5710",agent:"001"}
	now:=time.Date(2026,10,17,10,0,0,0,time.UTC)
	left:=map[authorizedWazuhCohort][]authorizedWazuhObservation{key:{{id:"first",at:now}}}
	right:=map[authorizedWazuhCohort][]authorizedWazuhObservation{key:{
		{id:"near",at:now.Add(2*time.Second)},{id:"far",at:now.Add(5*time.Second)},
	}}
	if _,ok:=selectedAuthorizedWazuhPairByEntryIDs(left,right,"first","far");ok{t.Fatal("non-nearest peer accepted")}
	if p,ok:=selectedAuthorizedWazuhPairByEntryIDs(left,right,"first","near");!ok||p.SecondEntryID!="near"{t.Fatal("valid nearest pair rejected")}
	if _,ok:=selectedAuthorizedWazuhPairByEntryIDs(left,right,"other","near");ok{t.Fatal("foreign first ID accepted")}
}
