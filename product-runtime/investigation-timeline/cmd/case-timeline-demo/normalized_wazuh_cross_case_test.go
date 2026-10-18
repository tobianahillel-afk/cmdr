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

func TestCrossSnapshotWazuhRuleAgentTemporalLeadsFromRealPersistedSources(t *testing.T) {
	// Reuse two separately saved/owner-validated Cases as inputs. These may
	// contain unrelated generic JSONL sources; only Wazuh source metadata
	// from the existing first-party normalizer can form a Wazuh cohort.
	aOwner,aCase,_,bOwner,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	feedA:=writeJSONLEvents(t,
		"{\"id\":\"a-1\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"CONFIDENTIAL FIRST WAZUH SOURCE\",\"level\":9}}\n"+
		"{\"id\":\"a-2\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"CONFIDENTIAL SECOND WAZUH SOURCE\",\"level\":10}}\n"+
		"{\"id\":\"a-3\",\"timestamp\":\"2026-10-17T10:06:00Z\",\"agent\":{\"id\":\"004\"},\"rule\":{\"id\":\"5720\",\"description\":\"unrelated A\"}}\n")
	feedB:=writeJSONLEvents(t,
		"{\"id\":\"b-1\",\"timestamp\":\"2026-10-17T10:03:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"CONFIDENTIAL THIRD WAZUH SOURCE\",\"level\":8}}\n"+
		"{\"id\":\"b-2\",\"timestamp\":\"2026-10-17T10:04:00Z\",\"agent\":{\"id\":\"099\"},\"rule\":{\"id\":\"5710\",\"description\":\"different agent not a cohort\"}}\n"+
		"{\"id\":\"b-3\",\"timestamp\":\"2026-10-17T10:06:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"later window no matching A\"}}\n")
	aPath:=filepath.Join(t.TempDir(),"saved-Wazuh-case-A.json")
	bPath:=filepath.Join(t.TempDir(),"saved-Wazuh-case-B.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",aOwner,aCase,[]string{feedA},aPath,&bytes.Buffer{});err!=nil {
		t.Fatalf("first actual Wazuh Case persistence failed: %v",err)
	}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",bOwner,bCase,[]string{feedB},bPath,&bytes.Buffer{});err!=nil {
		t.Fatalf("second actual Wazuh Case persistence failed: %v",err)
	}
	firstRaw,err:=os.ReadFile(aPath)
	if err!=nil {t.Fatal(err)}
	secondRaw,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	aDigest:=fmt.Sprintf("%x",sha256.Sum256(firstRaw))
	bDigest:=fmt.Sprintf("%x",sha256.Sum256(secondRaw))
	// The original Wazuh alert files can disappear after their owner
	// projections have been safely saved.
	if err:=os.Remove(feedA);err!=nil {t.Fatal(err)}
	if err:=os.Remove(feedB);err!=nil {t.Fatal(err)}
	var output bytes.Buffer
	err=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",aPath,"-compare-saved-case",aCase,
		"-expect-snapshot-sha256",aDigest,"-compare-saved-with-snapshot",bPath,
		"-compare-saved-with-case",bCase,"-expect-other-snapshot-sha256",bDigest,
	},&output)
	if err!=nil {t.Fatalf("real dual verified Wazuh Case comparison failed: %v",err)}
	var result localCaseComparison
	if err:=json.Unmarshal(output.Bytes(),&result);err!=nil {t.Fatal(err)}
	if result.First.CaseID!=aCase || result.Second.CaseID!=bCase ||
		result.Format!="cmdr-local-verified-saved-case-comparison-v1" ||
		len(result.SharedWazuhWindows)!=1 {
		t.Fatalf("did not correlate only authorized Wazuh temporal cohorts: %+v",result)
	}
	// A 10:02:00 Wazuh source in the first authorized Case and the
	// independently stored 10:03:00 source in the second Case are EXACTLY
	// +60s apart. Link the two Case-scoped Entry IDs; do not invent shared
	// global source identities or reveal either raw source preview.
	if len(result.SharedWazuhEventPairs)!=1 {
		t.Fatalf("inclusive +60s precise Wazuh lead absent, or another cohort leaked: %+v",result.SharedWazuhEventPairs)
	}
	pair:=result.SharedWazuhEventPairs[0]
	if pair.RuleID!="5710" || pair.AgentID!="001" || pair.OffsetSeconds!=60 ||
		pair.FirstTimeUTC!="2026-10-17T10:02:00Z" ||
		pair.SecondTimeUTC!="2026-10-17T10:03:00Z" ||
		pair.FirstEntryID=="" || pair.SecondEntryID=="" {
		t.Fatalf("incorrect independently authorized Wazuh source pair: %+v",pair)
	}
	// IDs must actually belong to different saved Case owners. A mere
	// matching Wazuh rule ID or fabricated ref is not a navigable event.
	firstOwner,err:=loadDemoSnapshot(aPath)
	if err!=nil {t.Fatal(err)}
	secondOwner,err:=loadDemoSnapshot(bPath)
	if err!=nil {t.Fatal(err)}
	contains:=func(owner []string,target string) bool {
		for _,id:=range owner {if id==target {return true}}
		return false
	}
	var firstCatalog,secondCatalog []string
	for _,c:=range firstOwner.Cases {if c.ID==aCase {firstCatalog=c.EntryIDs}}
	for _,c:=range secondOwner.Cases {if c.ID==bCase {secondCatalog=c.EntryIDs}}
	if !contains(firstCatalog,pair.FirstEntryID) || !contains(secondCatalog,pair.SecondEntryID) {
		t.Fatalf("event pair not navigable in its independently verified Case: %+v",pair)
	}
	pairJSON,err:=json.Marshal(pair)
	if err!=nil {t.Fatal(err)}
	for _,secret:=range []string{"CONFIDENTIAL","local-jsonl-source-","Bearer"} {
		if strings.Contains(string(pairJSON),secret) {
			t.Fatalf("precise event link disclosed source preview or bearer: %q",secret)
		}
	}
	lead:=result.SharedWazuhWindows[0]
	if lead.StartUTC!="2026-10-17T10:00:00Z" || lead.RuleID!="5710" ||
		lead.AgentID!="001" || lead.FirstEvents!=2 || lead.SecondEvents!=1 {
		t.Fatalf("wrong independently owner-authorized Wazuh overlap: %+v",lead)
	}
	for _,forbidden:=range []string{
		"CONFIDENTIAL FIRST WAZUH SOURCE","CONFIDENTIAL SECOND WAZUH SOURCE",
		"CONFIDENTIAL THIRD WAZUH SOURCE","Bearer token",
		"local-jsonl-entry-","local-jsonl-source-","127.0.0.1:",
	} {
		// Legacy Case aggregate may contain a small number of already
		// authorized source references under top_sources. The new Wazuh
		// cohort itself must never include source IDs, so check that
		// separately below instead of claiming old response fields absent.
		if strings.HasPrefix(forbidden,"local-jsonl-") { continue }
		if strings.Contains(output.String(),forbidden) {
			t.Fatalf("Wazuh comparison revealed raw source material: %q",forbidden)
		}
	}
	serializedLead,err:=json.Marshal(lead)
	if err!=nil {t.Fatal(err)}
	if strings.Contains(string(serializedLead),"local-jsonl-") ||
		strings.Contains(string(serializedLead),"CONFIDENTIAL") {
		t.Fatal("Wazuh lead contains evidence references or previews")
	}
}

func TestCrossSnapshotWazuhCohortsFailClosedOnChangedSecondCase(t *testing.T) {
	// A valid generic JSONL owner has no Wazuh labels: do not invent
	// rule/agent identity merely because Cases share an overlapping time.
	aPath,aCase,aDigest,bPath,bCase,bDigest:=buildTwoVerifiedLocalCaseArchives(t)
	var output bytes.Buffer
	if err:=compareVerifiedSavedCases(context.Background(),"analyst-a",
		aPath,aCase,aDigest,bPath,bCase,bDigest,&output);err!=nil {
		t.Fatal(err)
	}
	var result localCaseComparison
	if err:=json.Unmarshal(output.Bytes(),&result);err!=nil {t.Fatal(err)}
	if len(result.SharedWazuhWindows)!=0 || len(result.SharedWazuhEventPairs)!=0 {
		t.Fatal("ordinary normalized JSONL cases produced invented Wazuh agent/rule correlation")
	}
	// A parseable single-byte modification to the second Case must be
	// detected before any first-Case comparison receipt is emitted.
	content,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,append(content,'\n'),0600);err!=nil {t.Fatal(err)}
	output.Reset()
	if err:=compareVerifiedSavedCases(context.Background(),"analyst-a",
		aPath,aCase,aDigest,bPath,bCase,bDigest,&output);err==nil || output.Len()!=0 {
		t.Fatalf("second Case tampering yielded partial or unauthorized comparison: %q",output.String())
	}
}
