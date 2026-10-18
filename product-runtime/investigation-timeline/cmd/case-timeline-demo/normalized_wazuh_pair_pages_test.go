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

func savedWazuhPairPageFixture(t *testing.T) (string,string,string,string,string,string) {
	t.Helper()
	aOwner,aCase,_,bOwner,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	var firstFeed strings.Builder
	for i:=0;i<12;i++ {
		fmt.Fprintf(&firstFeed,
			"{\"id\":\"event-%02d\",\"timestamp\":\"2026-10-17T10:00:%02dZ\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE FIRST %02d\",\"level\":9}}\n",
			i,i,i)
	}
	aSource:=writeJSONLEvents(t,firstFeed.String())
	bSource:=writeJSONLEvents(t,
		"{\"id\":\"peer\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE SECOND\",\"level\":8}}\n")
	aPath:=filepath.Join(t.TempDir(),"saved-first.json")
	bPath:=filepath.Join(t.TempDir(),"saved-second.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",aOwner,aCase,[]string{aSource},aPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",bOwner,bCase,[]string{bSource},bPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	aBytes,err:=os.ReadFile(aPath);if err!=nil {t.Fatal(err)}
	bBytes,err:=os.ReadFile(bPath);if err!=nil {t.Fatal(err)}
	return aPath,aCase,fmt.Sprintf("%x",sha256.Sum256(aBytes)),
		bPath,bCase,fmt.Sprintf("%x",sha256.Sum256(bBytes))
}

type wazuhPairPageReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	FirstCase string `json:"first_case_id"`
	SecondCase string `json:"second_case_id"`
	Offset int `json:"offset"`
	HasMore bool `json:"has_more"`
	Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
}

func TestAuthorizedWazuhPairPaginationReachesBeyondFirstTen(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash:=savedWazuhPairPageFixture(t)
	base:=[]string{"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash}
	read:=func(offset string) wazuhPairPageReceipt {
		t.Helper()
		args:=append(append([]string{},base...),"-list-saved-wazuh-pair-offset",offset)
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&out);err!=nil {t.Fatalf("pair page %s unavailable: %v",offset,err)}
		for _,secret:=range []string{"PRIVATE FIRST","PRIVATE SECOND","Bearer ","Local analyst viewer:","127.0.0.1:"} {
			if strings.Contains(out.String(),secret) {t.Fatalf("sensitive preview or listener leaked: %q",secret)}
		}
		var doc wazuhPairPageReceipt
		if err:=json.Unmarshal(out.Bytes(),&doc);err!=nil {t.Fatal(err)}
		if doc.Format!="cmdr-local-authorized-wazuh-pair-page-v1" ||
			!doc.Synthetic || doc.FirstCase!=aCase || doc.SecondCase!=bCase {
			t.Fatalf("wrong authorized Case pairs: %+v",doc)
		}
		return doc
	}
	one:=read("0")
	two:=read("10")
	three:=read("12")
	if one.Offset!=0 || len(one.Pairs)!=10 || !one.HasMore ||
		two.Offset!=10 || len(two.Pairs)!=2 || two.HasMore ||
		three.Offset!=12 || len(three.Pairs)!=0 || three.HasMore {
		t.Fatalf("bounded pair pagination incorrect: page0=%+v page10=%+v page12=%+v",one,two,three)
	}
	if len(one.Pairs)>0 && one.Pairs[9].FirstEntryID==two.Pairs[0].FirstEntryID {
		t.Fatal("page boundary repeats an authorized Entry")
	}
	if two.Pairs[0].OffsetSeconds!=10 || two.Pairs[1].OffsetSeconds!=9 {
		t.Fatalf("stable nearest-peer offsets corrupted: %+v",two.Pairs)
	}
	// This integration combines two previously conflicting navigation
	// branches. An explicit page may not also find a specific Entry or
	// inspect private evidence: reject before metadata or file publication.
	for _,extra:=range [][]string{
		{"-find-saved-pair-first-entry",two.Pairs[0].FirstEntryID},
		{"-inspect-saved-pair-first-entry",two.Pairs[0].FirstEntryID,
			"-inspect-saved-pair-second-entry",two.Pairs[0].SecondEntryID,
			"-private-report",filepath.Join(t.TempDir(),"must-not-write.json")},
	} {
		args:=append(append(append([]string{},base...),
			"-list-saved-wazuh-pair-offset","10"),extra...)
		var denied bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&denied);err==nil || denied.Len()!=0 {
			t.Fatalf("conflicting pair selection leaked or executed: %v, %q",err,denied.String())
		}
		for i:=0;i<len(extra);i++ {
			if extra[i]=="-private-report" {
				if _,err:=os.Stat(extra[i+1]);!os.IsNotExist(err) {
					t.Fatal("conflicting metadata/page selector produced sensitive private evidence")
				}
			}
		}
	}
	// Both newly integrated selectors remain usable individually:
	// a paginated Entry can be resolved by exact first-ID discovery.
	var byExact bytes.Buffer
	findArgs:=append(append([]string{},base...),
		"-find-saved-pair-first-entry",two.Pairs[0].FirstEntryID)
	if err:=runDemoCLI(context.Background(),findArgs,&byExact);err!=nil {
		t.Fatalf("exact Entry discovery regressed after pagination integration: %v",err)
	}
	var lead struct{ Pair sharedAuthorizedWazuhEventPair `json:"pair"` }
	if err:=json.Unmarshal(byExact.Bytes(),&lead);err!=nil {t.Fatal(err)}
	if lead.Pair.FirstEntryID!=two.Pairs[0].FirstEntryID ||
		lead.Pair.SecondEntryID!=two.Pairs[0].SecondEntryID ||
		bytes.Contains(byExact.Bytes(),[]byte("PRIVATE FIRST")) {
		t.Fatalf("exact discovery mismatched paginated authorized pair: %q",byExact.String())
	}
	// A discovered lead beyond ten must lead to the SAME private, typed-
	// source re-authorized dossier as direct exact-ID inspection.
	dest:=filepath.Join(t.TempDir(),"selected-11th-pair.json")
	cmd:=append(append([]string{},base...),
		"-inspect-saved-pair-first-entry",two.Pairs[0].FirstEntryID,
		"-inspect-saved-pair-second-entry",two.Pairs[0].SecondEntryID,
		"-private-report",dest)
	var stdout bytes.Buffer
	if err:=runDemoCLI(context.Background(),cmd,&stdout);err!=nil {t.Fatal(err)}
	if stdout.Len()!=0 {t.Fatal("sensitive dossier written to stdout")}
	private,err:=os.ReadFile(dest)
	if err!=nil {t.Fatal(err)}
	if !bytes.Contains(private,[]byte("PRIVATE FIRST 10")) ||
		!bytes.Contains(private,[]byte("PRIVATE SECOND")) {
		t.Fatal("pagination lead did not reopen independently authorized source evidence")
	}
}

func TestAuthorizedWazuhPairPaginationDenialsReturnNoPartialReceipt(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash:=savedWazuhPairPageFixture(t)
	base:=[]string{"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
		"-list-saved-wazuh-pair-offset","10"}
	for _,tc:=range []struct{name,flag,value string}{
		{"wrong second digest","-expect-other-snapshot-sha256",strings.Repeat("0",64)},
		{"unauthorized principal","-principal","other-person"},
		{"other Case","-compare-saved-with-case","foreign-case"},
		{"negative offset","-list-saved-wazuh-pair-offset","-1"},
		{"oversized offset","-list-saved-wazuh-pair-offset","5001"},
	} {
		t.Run(tc.name,func(t *testing.T){
			args:=append([]string{},base...)
			for i:=0;i+1<len(args);i++{if args[i]==tc.flag{args[i+1]=tc.value;break}}
			var out bytes.Buffer
			if err:=runDemoCLI(context.Background(),args,&out);err==nil || out.Len()!=0 {
				t.Fatalf("denied pair page emitted metadata: %v %q",err,out.String())
			}
		})
	}
	for _,additional:=range [][]string{
		{"-inspect-saved-pair","1"},
		{"-inspect-saved-pair-first-entry","some-id","-inspect-saved-pair-second-entry","other-id"},
		{"-include-pair-generic-evidence"},
		{"-private-report",filepath.Join(t.TempDir(),"unexpected-file.json")},
		{"-serve-import"},
	} {
		args:=append(append([]string{},base...),additional...)
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&out);err==nil || out.Len()!=0 {
			t.Fatalf("mixed sensitive pair modes accepted: %v %q",additional,out.String())
		}
	}
	// A syntactically valid byte change in the SECOND Case fails integrity
	// before the first Case's authorized data can be printed.
	original,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,append(original,'\n'),0600);err!=nil {t.Fatal(err)}
	var out bytes.Buffer
	if err:=runDemoCLI(context.Background(),base,&out);err==nil || out.Len()!=0 {
		t.Fatalf("changed source returned cross-Case data: %v %q",err,out.String())
	}
}
