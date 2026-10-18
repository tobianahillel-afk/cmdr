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

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

// One cross-product acceptance: a user can save and hunt heterogeneous local
// evidence in two separately authorized Cases, correlate real Wazuh events
// across the Cases, and inspect BOTH original paired sources privately.
// All three journeys reuse the SAME existing Case/Timeline/typed-source
// authorization and never merge the two owner graphs.
func TestMixedCasesHeadlessHuntCorrelateAndOpenRealPrivateWazuhPair(t *testing.T) {
	aWazuh, aCase, _, bWazuh, bCase, _, originalPair := makeSavedPairDrilldownOwners(t)
	firstGeneric := writeJSONLEvents(t,
		"{\"id\":\"generic-a\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"GENERIC-CONTEXT-FIRST-NOT-WAZUH\"}\n")
	secondGeneric := writeJSONLEvents(t,
		"{\"id\":\"generic-b\",\"timestamp\":\"2026-10-17T10:00:22Z\",\"message\":\"GENERIC-CONTEXT-SECOND-NOT-WAZUH\"}\n")
	first := filepath.Join(t.TempDir(), "mixed-first-case.json")
	second := filepath.Join(t.TempDir(), "mixed-second-case.json")
	for _,item := range []struct{owner,caseID,feed,dest string}{
		{aWazuh,aCase,firstGeneric,first},
		{bWazuh,bCase,secondGeneric,second},
	} {
		if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",
			item.owner,item.caseID,[]string{item.feed},item.dest,&bytes.Buffer{});err!=nil {
			t.Fatalf("unable to save mixed Wazuh and generic Case: %v",err)
		}
	}
	readDigest:=func(path string) string {
		t.Helper()
		raw,err:=os.ReadFile(path)
		if err!=nil {t.Fatal(err)}
		return fmt.Sprintf("%x",sha256.Sum256(raw))
	}
	firstHash,secondHash:=readDigest(first),readDigest(second)
	if err:=os.Remove(firstGeneric);err!=nil {t.Fatal(err)}
	if err:=os.Remove(secondGeneric);err!=nil {t.Fatal(err)}
	// Both mixed local owners remain fully usable after the raw JSONL
	// sources are removed. Search Job is the existing authenticated handler,
	// invoked in-process; do not ship another predicate/search engine.
	for _,side:=range []struct{path,caseID,hash,literal string}{
		{first,aCase,firstHash,"GENERIC-CONTEXT-FIRST-NOT-WAZUH"},
		{second,bCase,secondHash,"GENERIC-CONTEXT-SECOND-NOT-WAZUH"},
	} {
		var hunt bytes.Buffer
		if err:=runSavedCaseHeadlessHunt(context.Background(),"analyst-a",
			side.path,side.caseID,side.hash,
			"preview","generic",side.literal,"2026-10-17T09:59:00Z",
			"2026-10-17T10:02:00Z",&hunt);err!=nil {
			t.Fatalf("generic member of mixed saved Case cannot be hunted: %v",err)
		}
		var receipt struct{
			CaseID string `json:"case_id"`
			Refs []search.ResultReference `json:"authorized_refs"`
			Bursts []wazuhSearchBurst `json:"wazuh_bursts"`
		}
		if err:=json.Unmarshal(hunt.Bytes(),&receipt);err!=nil {t.Fatal(err)}
		if receipt.CaseID!=side.caseID || len(receipt.Refs)!=1 || len(receipt.Bursts)!=0 {
			t.Fatalf("generic-only authenticated search invented Wazuh results: %s",hunt.String())
		}
		if strings.Contains(hunt.String(),side.literal) ||
			strings.Contains(hunt.String(),"Bearer ") {
			t.Fatal("generic search receipt exposed original raw evidence or bearer")
		}
	}
	var comparison bytes.Buffer
	if err:=compareVerifiedSavedCases(context.Background(),"analyst-a",
		first,aCase,firstHash,second,bCase,secondHash,&comparison);err!=nil {
		t.Fatalf("mixed-source Case comparison failed: %v",err)
	}
	var compared localCaseComparison
	if err:=json.Unmarshal(comparison.Bytes(),&compared);err!=nil {t.Fatal(err)}
	if compared.First.CaseID!=aCase || compared.Second.CaseID!=bCase ||
		len(compared.SharedWazuhEventPairs)!=1 ||
		compared.SharedWazuhEventPairs[0].FirstEntryID!=originalPair.FirstEntryID ||
		compared.SharedWazuhEventPairs[0].SecondEntryID!=originalPair.SecondEntryID ||
		compared.SharedWazuhEventPairs[0].OffsetSeconds!=20 ||
		compared.First.VisibleEntries<2 || compared.Second.VisibleEntries<2 {
		t.Fatalf("mixed Case comparison dropped authorized events or forged cross-Case pairs: %+v",compared)
	}
	for _,s:=range []string{"GENERIC-CONTEXT-FIRST","GENERIC-CONTEXT-SECOND","PRIVATE FIRST OWNER PAIR","PRIVATE SECOND OWNER PAIR","Bearer "} {
		if strings.Contains(comparison.String(),s) {
			t.Fatalf("cross-Case comparison leaked raw source data: %q",s)
		}
	}
	destination:=filepath.Join(t.TempDir(),"private-correlated-sources.json")
	if err:=writeVerifiedCorrelatedPairPrivate(context.Background(),"analyst-a",
		first,aCase,firstHash,second,bCase,secondHash,1,destination);err!=nil {
		t.Fatalf("cannot reopen two real Wazuh paired sources from mixed Case owners: %v",err)
	}
	raw,err:=os.ReadFile(destination)
	if err!=nil {t.Fatal(err)}
	var result struct{
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		First localEntryInspection `json:"first"`
		Second localEntryInspection `json:"second"`
	}
	if err:=json.Unmarshal(raw,&result);err!=nil {t.Fatal(err)}
	if result.Pair.FirstEntryID!=originalPair.FirstEntryID ||
		result.Pair.SecondEntryID!=originalPair.SecondEntryID ||
		!strings.Contains(result.First.Preview,"PRIVATE FIRST OWNER PAIR") ||
		!strings.Contains(result.Second.Preview,"PRIVATE SECOND OWNER PAIR") ||
		result.First.CaseID!=aCase || result.Second.CaseID!=bCase {
		t.Fatalf("private two-source report did not reopen independently authorized Wazuh events: %+v",result)
	}
	for _,forbidden:=range []string{"GENERIC-CONTEXT-FIRST","GENERIC-CONTEXT-SECOND","Bearer "} {
		if bytes.Contains(raw,[]byte(forbidden)) {
			t.Fatalf("pair preview bundle mixed in unrelated generic source: %q",forbidden)
		}
	}
	// Even with a valid first Case, a changed second owner must not write
	// only the first private preview or overwrite an existing report.
	changed,err:=os.ReadFile(second)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(second,append(changed,'\n'),0600);err!=nil {t.Fatal(err)}
	denied:=filepath.Join(t.TempDir(),"no-partial-after-second-integrity-failure.json")
	if err:=writeVerifiedCorrelatedPairPrivate(context.Background(),"analyst-a",
		first,aCase,firstHash,second,bCase,secondHash,1,denied);err==nil {
		t.Fatal("changed second owner allowed partial private Wazuh report")
	}
	if _,err:=os.Lstat(denied);!os.IsNotExist(err) {
		t.Fatal("second Case fingerprint failure published sensitive first-Case preview")
	}
}
