package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestAuthorizedPairGenericContextSelectsNearestTenWithExactUTCWindow(t *testing.T) {
	const selectedCase = "forensic-case-on-disk"
	owner, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil { t.Fatal(err) }
	center := time.Date(2026,10,17,10,0,0,0,time.UTC)
	offsets := []int{-61,-60,-2,-1,0,1,2,3,4,5,6,7,8,9,60,61}
	var jsonl strings.Builder
	for i, offset := range offsets {
		event := map[string]string{
			"id": "relative-observation-" + strconv.Itoa(i),
			"timestamp": center.Add(time.Duration(offset)*time.Second).Format(time.RFC3339Nano),
			"message": "PRIVATE GENERIC EVENT " + strconv.Itoa(i),
		}
		raw, err := json.Marshal(event)
		if err != nil { t.Fatal(err) }
		jsonl.Write(raw)
		jsonl.WriteByte('\n')
	}
	owner,err=importLocalEvents(context.Background(),owner,"analyst-a",selectedCase,[]byte(jsonl.String()))
	if err!=nil {t.Fatalf("cannot import actual generic evidence: %v",err)}
	ids,err:=importedNormalizedIDs(owner,selectedCase)
	if err!=nil {t.Fatal(err)}
	summary,err:=nearbyAuthorizedGenericPairContext(context.Background(),"analyst-a",
		owner,selectedCase,ids,center.Format(time.RFC3339Nano))
	if err!=nil {t.Fatalf("source-authorized context rejected: %v",err)}
	if !summary.Truncated || len(summary.Entries)!=10 {
		t.Fatalf("expected ten nearest leads and actual truncation, got %+v",summary)
	}
	want := []int{0,-1,1,-2,2,3,4,5,6,7}
	for i, lead:=range summary.Entries {
		if lead.OffsetSeconds!=want[i] {
			t.Fatalf("context rank %d offset %d, want %d",i,lead.OffsetSeconds,want[i])
		}
		if _,ok:=ids[lead.EntryID];!ok {
			t.Fatalf("context %d is not a canonical imported event",i)
		}
		inspection,err:=inspectDemoEntryWithData(context.Background(),"analyst-a",owner,selectedCase,lead.EntryID)
		if err!=nil || inspection.CaseID!=selectedCase || inspection.EntryID!=lead.EntryID {
			t.Fatalf("context %d cannot be re-opened via its own first-party source owner: %v",i,err)
		}
	}
	// Inclusive +/-60s boundaries must remain members when fewer than ten
	// entries exist; -61s and +61s must never be included.
	var ends strings.Builder
	for _,offset:=range []int{-61,-60,60,61} {
		event:=map[string]string{
			"id":fmt.Sprintf("edge-%d",offset),
			"timestamp":center.Add(time.Duration(offset)*time.Second).Format(time.RFC3339),
			"message":fmt.Sprintf("EDGE PRIVATE %d",offset),
		}
		raw,_:=json.Marshal(event)
		ends.Write(raw);ends.WriteByte('\n')
	}
	pristine,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	pristine,err=importLocalEvents(context.Background(),pristine,"analyst-a",selectedCase,[]byte(ends.String()))
	if err!=nil {t.Fatal(err)}
	edgeIDs,err:=importedNormalizedIDs(pristine,selectedCase)
	if err!=nil {t.Fatal(err)}
	edges,err:=nearbyAuthorizedGenericPairContext(context.Background(),"analyst-a",pristine,selectedCase,
		edgeIDs,center.Format(time.RFC3339Nano))
	if err!=nil || edges.Truncated || len(edges.Entries)!=2 ||
		edges.Entries[0].OffsetSeconds!=-60 || edges.Entries[1].OffsetSeconds!=60 {
		t.Fatalf("inclusive UTC boundary produced incorrect generic context: %+v %v",edges,err)
	}
}

func TestPrivateWazuhPairContextDoesNotExportRevokedGenericSource(t *testing.T) {
	aInput,aCase,_,bPath,bCase,bHash,pair:=makeSavedPairDrilldownOwners(t)
	feed:=writeJSONLEvents(t,
		"{\"id\":\"generic-to-revoke\",\"timestamp\":\"2026-10-17T10:00:05Z\",\"message\":\"SECRET DENIED GENERIC SOURCE CONTENT\"}\n")
	aPath:=filepath.Join(t.TempDir(),"first-with-private-generic.json")
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",
		aInput,aCase,[]string{feed},aPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	first,err:=loadDemoSnapshot(aPath)
	if err!=nil {t.Fatal(err)}
	if len(first.Principals)!=1 {t.Fatal("unexpected owner scope")}
	identity:=first.Principals[0]
	const deniedPrincipal="other-source-only-analyst"
	first.Principals=append(first.Principals,timeline.MemoryTimelinePrincipal{
		ID:deniedPrincipal,Tenant:identity.Tenant,Environment:identity.Environment,
	})
	var genericSource string
	for _,entry:=range first.Entries {
		if entry.CaseID==aCase && entry.Value.Provenance=="operator-local-jsonl" {
			genericSource=entry.Value.SourceRef
		}
	}
	if genericSource=="" {t.Fatal("normalized generic source missing")}
	revoked:=false
	for i:=range first.Sources {
		if first.Sources[i].Identity.Reference==genericSource {
			first.Sources[i].Readers=[]string{deniedPrincipal}
			revoked=true
		}
	}
	if !revoked {t.Fatal("typed source not found")}
	if _,err:=newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",first);err!=nil {
		t.Fatalf("source-revoked owner must remain structurally valid: %v",err)
	}
	changed,err:=json.Marshal(first)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(aPath,changed,0600);err!=nil {t.Fatal(err)}
	aHash:=fmt.Sprintf("%x",sha256.Sum256(changed))
	out:=filepath.Join(t.TempDir(),"authorized-pair-without-revoked-generic.json")
	if err:=writeVerifiedCorrelatedPairPrivate(context.Background(),"analyst-a",
		aPath,aCase,aHash,bPath,bCase,bHash,1,out);err!=nil {
		t.Fatalf("Wazuh source pair should remain authorized after unrelated generic revocation: %v",err)
	}
	raw,err:=os.ReadFile(out)
	if err!=nil {t.Fatal(err)}
	if bytes.Contains(raw,[]byte("SECRET DENIED GENERIC SOURCE CONTENT")) ||
		bytes.Contains(raw,[]byte(genericSource)) ||
		bytes.Contains(raw,[]byte(deniedPrincipal)) {
		t.Fatal("paired sensitive dossier leaked revoked generic source or reader")
	}
	var report struct {
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		FirstGenericContext authorizedGenericPairContextSummary `json:"first_generic_context"`
	}
	if err:=json.Unmarshal(raw,&report);err!=nil {t.Fatal(err)}
	// The original Wazuh Case fixture already contains a DIFFERENT,
	// owner-authorized generic event at exactly +60 seconds. Revoking
	// the newly imported source must not hide that valid contextual lead.
	// It must only hide the revoked source and its canonical Entry ID.
	revokedEntryID := "local-jsonl-entry-" + strings.TrimPrefix(genericSource,"local-jsonl-source-")
	if report.Pair.FirstEntryID!=pair.FirstEntryID ||
		report.Pair.SecondEntryID!=pair.SecondEntryID ||
		len(report.FirstGenericContext.Entries)!=1 ||
		report.FirstGenericContext.Entries[0].OffsetSeconds!=60 {
		t.Fatalf("first-owner authorization or legitimate contextual visibility changed: %+v",report)
	}
	for _,lead:=range report.FirstGenericContext.Entries {
		if lead.EntryID==revokedEntryID {
			t.Fatal("revoked generic source was reintroduced as contextual evidence")
		}
	}
}

func TestGenericPairContextRejectsInvalidTimeAndCancelledWork(t *testing.T) {
	const selectedCase="forensic-case-on-disk"
	data,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	if _,err:=nearbyAuthorizedGenericPairContext(context.Background(),"analyst-a",
		data,selectedCase,map[string]struct{}{},"not-an-utc-time");err==nil {
		t.Fatal("invalid correlated Wazuh event clock admitted generic evidence")
	}
	cancelled,stop:=context.WithCancel(context.Background());stop()
	if _,err:=nearbyAuthorizedGenericPairContext(cancelled,"analyst-a",
		data,selectedCase,map[string]struct{}{},"2026-10-17T10:00:00Z");err==nil {
		t.Fatal("cancelled Case investigation produced temporal leads")
	}
}
