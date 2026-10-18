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

// Verify the actual Wazuh import -> two private snapshots -> independent
// digest/owner checks -> cross-Case triage -> precise authorized pairs.
// Many matches and misleading cohorts must not create an unbounded or
// ambiguous forensic correlation receipt.
func TestCrossCaseWazuhPairsNearestDeterministicAndBounded(t *testing.T) {
	aOwner,aCase,_,bOwner,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	var firstFeed strings.Builder
	for i:=0;i<12;i++ {
		fmt.Fprintf(&firstFeed,
			"{\"id\":\"a-%02d\",\"timestamp\":\"2026-10-17T10:00:%02dZ\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE FIRST ALERT\",\"level\":9}}\n",i,i)
	}
	secondFeed:=strings.Join([]string{
		`{"id":"b-near-5","timestamp":"2026-10-17T10:00:05Z","agent":{"id":"001"},"rule":{"id":"5710","description":"PRIVATE B FIVE","level":7}}`,
		`{"id":"b-near-7","timestamp":"2026-10-17T10:00:07Z","agent":{"id":"001"},"rule":{"id":"5710","description":"PRIVATE B SEVEN","level":8}}`,
		`{"id":"b-wrong-rule","timestamp":"2026-10-17T10:00:06Z","agent":{"id":"001"},"rule":{"id":"9999","description":"PRIVATE WRONG RULE"}}`,
		`{"id":"b-wrong-agent","timestamp":"2026-10-17T10:00:06Z","agent":{"id":"999"},"rule":{"id":"5710","description":"PRIVATE WRONG AGENT"}}`,
		`{"id":"b-outside","timestamp":"2026-10-17T10:04:00Z","agent":{"id":"001"},"rule":{"id":"5710","description":"PRIVATE TOO LATE"}}`,
	}, "\n")+"\n"
	aFeed:=writeJSONLEvents(t,firstFeed.String())
	bFeed:=writeJSONLEvents(t,secondFeed)
	aPath:=filepath.Join(t.TempDir(),"first-authorized-case.json")
	bPath:=filepath.Join(t.TempDir(),"second-authorized-case.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",aOwner,aCase,[]string{aFeed},aPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",bOwner,bCase,[]string{bFeed},bPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	aBytes,err:=os.ReadFile(aPath)
	if err!=nil {t.Fatal(err)}
	bBytes,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	aDigest,bDigest:=fmt.Sprintf("%x",sha256.Sum256(aBytes)),fmt.Sprintf("%x",sha256.Sum256(bBytes))
	if err:=os.Remove(aFeed);err!=nil {t.Fatal(err)}
	if err:=os.Remove(bFeed);err!=nil {t.Fatal(err)}
	aData,aImported,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",aPath,aCase,aDigest)
	if err!=nil {t.Fatal(err)}
	bData,bImported,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",bPath,bCase,bDigest)
	if err!=nil {t.Fatal(err)}
	if len(aData.Cases)!=1 || len(bData.Cases)!=1 {t.Fatal("owner scope leaked another Case")}
	run:=func() (localCaseComparison,[]byte) {
		t.Helper()
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),[]string{
			"-principal","analyst-a","-snapshot",aPath,
			"-compare-saved-case",aCase,"-expect-snapshot-sha256",aDigest,
			"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
			"-expect-other-snapshot-sha256",bDigest,
		},&out);err!=nil {t.Fatal(err)}
		var decoded localCaseComparison
		if err:=json.Unmarshal(out.Bytes(),&decoded);err!=nil {t.Fatal(err)}
		return decoded,out.Bytes()
	}
	got,raw:=run()
	if got.First.CaseID!=aCase || got.Second.CaseID!=bCase ||
		len(got.SharedWazuhEventPairs)!=10 {
		t.Fatalf("bounded two-owner Wazuh correlation missing: %+v",got.SharedWazuhEventPairs)
	}
	for i,p:=range got.SharedWazuhEventPairs {
		if p.RuleID!="5710" || p.AgentID!="001" ||
			p.FirstTimeUTC!=fmt.Sprintf("2026-10-17T10:00:%02dZ",i) ||
			p.FirstEntryID=="" || p.SecondEntryID=="" {
			t.Fatalf("cohort or deterministic chronological ordering incorrect at %d: %+v",i,p)
		}
		if _,ok:=aImported[p.FirstEntryID];!ok {t.Fatal("first source outside independently authorized Case")}
		if _,ok:=bImported[p.SecondEntryID];!ok {t.Fatal("second source outside independently authorized Case")}
		if p.OffsetSeconds < -60 || p.OffsetSeconds > 60 {t.Fatal("pair exceeded inclusive time window")}
	}
	if got.SharedWazuhEventPairs[0].SecondTimeUTC!="2026-10-17T10:00:05Z" ||
		got.SharedWazuhEventPairs[0].OffsetSeconds!=5 ||
		got.SharedWazuhEventPairs[9].SecondTimeUTC!="2026-10-17T10:00:07Z" ||
		got.SharedWazuhEventPairs[9].OffsetSeconds!=-2 {
		t.Fatalf("nearest event selection not correct at either end: %+v",got.SharedWazuhEventPairs)
	}
	for _,secret:=range []string{
		"PRIVATE FIRST ALERT","PRIVATE B FIVE","PRIVATE B SEVEN",
		"PRIVATE WRONG RULE","PRIVATE WRONG AGENT","PRIVATE TOO LATE","Bearer ",
	} {
		if bytes.Contains(raw,[]byte(secret)) {t.Fatalf("comparison leaked raw evidence: %q",secret)}
	}
	_,again:=run()
	if !bytes.Equal(raw,again) {t.Fatal("owner-authorized pair order varies across identical verified inputs")}
}
