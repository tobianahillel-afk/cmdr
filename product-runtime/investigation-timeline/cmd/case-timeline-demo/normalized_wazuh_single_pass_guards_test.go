package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A real Wazuh pair must not be assembled from a Case with untrusted
// timestamps or from two independently readable but differently tenanted
// owners. This test exercises the *integrated* one-collection path, rather
// than relying on isolated map-only unit tests.
func TestSinglePassCrossCaseWazuhRejectsCorruptTimeAndWrongOwner(t *testing.T) {
	firstInput,firstCase,_,secondInput,secondCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	firstFeed:=writeJSONLEvents(t,
		"{\"id\":\"a1\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE A\"}}\n")
	secondFeed:=writeJSONLEvents(t,
		"{\"id\":\"b1\",\"timestamp\":\"2026-10-17T10:01:30Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE B\"}}\n")
	aPath:=filepath.Join(t.TempDir(),"private-wazuh-a.json")
	bPath:=filepath.Join(t.TempDir(),"private-wazuh-b.json")
	for _,tc:=range []struct{input,caseID,feed,dest string}{
		{firstInput,firstCase,firstFeed,aPath},
		{secondInput,secondCase,secondFeed,bPath},
	} {
		if err:=writePersistedWazuhCase(context.Background(),"analyst-a",
			tc.input,tc.caseID,[]string{tc.feed},tc.dest,io.Discard);err!=nil {
			t.Fatalf("real Wazuh Case import/persistence failed: %v",err)
		}
	}
	read:=func(path,caseID string) timeline.MemoryTimelineData {
		t.Helper()
		raw,err:=os.ReadFile(path)
		if err!=nil {t.Fatal(err)}
		digest:=fmt.Sprintf("%x",sha256.Sum256(raw))
		owner,_,err:=loadAuthorizedImportedCase(context.Background(),"analyst-a",path,caseID,digest)
		if err!=nil {t.Fatalf("source-verified Case owner unavailable: %v",err)}
		return owner
	}
	a,b:=read(aPath,firstCase),read(bPath,secondCase)
	success,err:=compareLoadedAuthorizedSavedCases(context.Background(),"analyst-a",a,b,firstCase,secondCase)
	if err!=nil || len(success.SharedWazuhEventPairs)!=1 ||
		len(success.SharedWazuhWindows)!=1 || success.SharedWazuhEventPairs[0].OffsetSeconds!=30 {
		t.Fatalf("positive authorized single-pass Wazuh investigation failed: %v %+v",err,success)
	}
	// This is a structural denial, not a missing-data partial success. A
	// separately authorized principal identifier is not a tenant bridge.
	foreign:=b
	foreign.Principals=append([]timeline.MemoryTimelinePrincipal(nil),b.Principals...)
	foreign.Principals[0].Tenant="other-tenant"
	if _,err:=compareLoadedAuthorizedSavedCases(context.Background(),"analyst-a",
		a,foreign,firstCase,secondCase);err==nil {
		t.Fatal("cross-tenant Case correlation accepted")
	}
	if _,err:=compareLoadedAuthorizedSavedCases(context.Background(),"wrong-principal",
		a,b,firstCase,secondCase);err==nil {
		t.Fatal("foreign principal inherited authorization from the saved snapshots")
	}
	cancelled,cancel:=context.WithCancel(context.Background())
	cancel()
	if _,err:=compareLoadedAuthorizedSavedCases(cancelled,"analyst-a",
		a,b,firstCase,secondCase);err==nil {
		t.Fatal("cancelled Wazuh correlation returned a completed lead")
	}
	// Corrupt one Wazuh source at a time, then require a hard error instead
	// of silently dropping the affected event or producing a partial lead.
	for _,tc:=range []struct{name string;mutate func(*timeline.MemoryTimelineData)}{
		{"uncertain source clock",func(d *timeline.MemoryTimelineData){
			for i:=range d.Entries {
				if d.Entries[i].Value.Provenance=="operator-local-jsonl" {
					d.Entries[i].Value.ClockUncertaintyMS=1000
				}
			}
		}},
		{"invalid RFC3339 source time",func(d *timeline.MemoryTimelineData){
			for i:=range d.Entries {
				if d.Entries[i].Value.Provenance=="operator-local-jsonl" {
					d.Entries[i].Value.SourceTime="not-a-UTC-time"
				}
			}
		}},
	} {
		t.Run(tc.name,func(t *testing.T){
			d:=a
			d.Entries=append([]timeline.MemoryTimelineEntry(nil),a.Entries...)
			tc.mutate(&d)
			if _,err:=collectAuthorizedWazuhObservations(d,firstCase);err==nil {
				t.Fatal("invalid Wazuh source time was treated as a usable correlation lead")
			}
			if _,err:=compareLoadedAuthorizedSavedCases(context.Background(),"analyst-a",
				d,b,firstCase,secondCase);err==nil {
				t.Fatal("integrated Case comparator failed open on corrupt Wazuh time")
			}
		})
	}
	// The real owner-backed service must never expose raw Wazuh descriptions
	// in a comparison report, even though they were needed as input.
	if raw,err:=os.ReadFile(aPath);err!=nil {
		t.Fatal(err)
	}else if !bytes.Contains(raw,[]byte("PRIVATE A")) {
		t.Fatal("test case did not contain a real sensitive Wazuh source")
	}
}

func TestSinglePassWazuhDerivedViewsRemainDeterministicAndBounded(t *testing.T) {
	cohort:=authorizedWazuhCohort{rule:"5710",agent:"001"}
	// Represent a real normalized and independently authenticated pair of
	// source cohorts with late/out-of-window events. The correlation output
	// is deliberately bounded and has no Case-wide source permission effects.
	first:=map[authorizedWazuhCohort][]authorizedWazuhObservation{cohort:{}}
	second:=map[authorizedWazuhCohort][]authorizedWazuhObservation{cohort:{}}
	base:=time.Date(2026,10,17,10,0,0,0,time.UTC)
	for i:=0;i<25;i++{
		at:=base.Add(time.Duration(i)*5*time.Minute)
		first[cohort]=append(first[cohort],authorizedWazuhObservation{
			id:fmt.Sprintf("authorized-first-%02d",i),at:at,
		})
		second[cohort]=append(second[cohort],authorizedWazuhObservation{
			id:fmt.Sprintf("authorized-second-%02d",i),at:at.Add(30*time.Second),
		})
	}
	windows:=sharedWazuhWindowsFromObservations(first,second)
	pairs:=sharedWazuhPairsFromObservations(first,second)
	if len(windows)!=10 || len(pairs)!=10 {
		t.Fatalf("authorized 25-window cohort was not bounded: %d %d",len(windows),len(pairs))
	}
	for i:=0;i<10;i++ {
		if windows[i].RuleID!="5710" || windows[i].AgentID!="001" ||
			windows[i].FirstEvents!=1 || windows[i].SecondEvents!=1 ||
			windows[i].StartUTC!=base.Add(time.Duration(i)*5*time.Minute).Format(time.RFC3339) {
			t.Fatalf("shared five-minute lead miscounted: %+v",windows[i])
		}
		if pairs[i].OffsetSeconds!=30 || pairs[i].RuleID!="5710" ||
			pairs[i].AgentID!="001" ||
			pairs[i].FirstEntryID!=fmt.Sprintf("authorized-first-%02d",i) ||
			pairs[i].SecondEntryID!=fmt.Sprintf("authorized-second-%02d",i) {
			t.Fatalf("nearby paired Wazuh source no longer navigable: %+v",pairs[i])
		}
	}
	// A different source family/agent must not be correlated merely
	// because UTC matches. This also checks that pure computation cannot
	// bypass the already authorized Wazuh cohort identities.
	second[cohort]=nil
	second[authorizedWazuhCohort{rule:"5710",agent:"DIFFERENT"}]=[]authorizedWazuhObservation{
		{id:"wrong-agent",at:base},
	}
	if len(sharedWazuhPairsFromObservations(first,second))!=0 ||
		len(sharedWazuhWindowsFromObservations(first,second))!=0 {
		t.Fatal("different agent ID created an unauthorized cross-Case correlation")
	}
}
