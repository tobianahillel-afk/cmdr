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

func TestPrivateWazuhPairDossierIncludesOnlyOwnerAuthorizedNearbyGenericLeads(t *testing.T) {
	firstWazuh, firstCase, _, secondWazuh, secondCase, _, pair := makeSavedPairDrilldownOwners(t)
	firstFeed := writeJSONLEvents(t,
		"{\"id\":\"generic-near-first\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PRIVATE-GENERIC-FIRST-CONTEXT\"}\n"+
		"{\"id\":\"generic-far-first\",\"timestamp\":\"2026-10-17T10:04:00Z\",\"message\":\"PRIVATE-GENERIC-OUTSIDE-WINDOW\"}\n")
	secondFeed := writeJSONLEvents(t,
		"{\"id\":\"generic-near-second\",\"timestamp\":\"2026-10-17T10:00:22Z\",\"message\":\"PRIVATE-GENERIC-SECOND-CONTEXT\"}\n")
	firstPath := filepath.Join(t.TempDir(), "case-first.json")
	secondPath := filepath.Join(t.TempDir(), "case-second.json")
	for _, side := range []struct{owner, caseID, feed, path string}{
		{firstWazuh,firstCase,firstFeed,firstPath},
		{secondWazuh,secondCase,secondFeed,secondPath},
	} {
		if err := writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",
			side.owner,side.caseID,[]string{side.feed},side.path,&bytes.Buffer{});err != nil {
			t.Fatalf("could not save authorized Case: %v",err)
		}
	}
	hash := func(path string) string {
		t.Helper()
		data,err := os.ReadFile(path)
		if err != nil { t.Fatal(err) }
		return fmt.Sprintf("%x",sha256.Sum256(data))
	}
	firstHash,secondHash := hash(firstPath),hash(secondPath)
	if err := os.Remove(firstFeed);err != nil { t.Fatal(err) }
	if err := os.Remove(secondFeed);err != nil { t.Fatal(err) }

	report := filepath.Join(t.TempDir(),"private-combined-dossier.json")
	if err := writeVerifiedCorrelatedPairPrivate(context.Background(),"analyst-a",
		firstPath,firstCase,firstHash,secondPath,secondCase,secondHash,1,report);err != nil {
		t.Fatalf("cannot build real correlated private dossier: %v",err)
	}
	raw,err := os.ReadFile(report)
	if err != nil { t.Fatal(err) }
	var dossier struct {
		Format string `json:"format"`
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
		First localEntryInspection `json:"first"`
		Second localEntryInspection `json:"second"`
		FirstContext authorizedGenericPairContextSummary `json:"first_generic_context"`
		SecondContext authorizedGenericPairContextSummary `json:"second_generic_context"`
	}
	if err := json.Unmarshal(raw,&dossier);err != nil { t.Fatal(err) }
	if dossier.Format!="cmdr-private-verified-wazuh-pair-v2" ||
		dossier.Pair.FirstEntryID!=pair.FirstEntryID ||
		dossier.Pair.SecondEntryID!=pair.SecondEntryID ||
		!strings.Contains(dossier.First.Preview,"PRIVATE FIRST OWNER PAIR") ||
		!strings.Contains(dossier.Second.Preview,"PRIVATE SECOND OWNER PAIR") {
		t.Fatalf("paired Wazuh source inspection regressed: %+v",dossier)
	}
	if len(dossier.FirstContext.Entries)!=2 || len(dossier.SecondContext.Entries)!=1 ||
		dossier.FirstContext.Truncated || dossier.SecondContext.Truncated ||
		dossier.FirstContext.Entries[0].TimeUTC!="2026-10-17T10:00:10Z" ||
		dossier.SecondContext.Entries[0].TimeUTC!="2026-10-17T10:00:22Z" ||
		dossier.FirstContext.Entries[0].OffsetSeconds!=10 ||
		dossier.FirstContext.Entries[1].TimeUTC!="2026-10-17T10:01:00Z" ||
		dossier.FirstContext.Entries[1].OffsetSeconds!=60 ||
		dossier.SecondContext.Entries[0].OffsetSeconds!=2 {
		t.Fatalf("nearby authorized generic Case context not attached with correct clock/Case: %+v %+v",
			dossier.FirstContext,dossier.SecondContext)
	}
	// The linked generic IDs are real source-owned imported refs in their
	// OWN Cases. They must never be a leaked raw payload or a cross-Case grant.
	for _, item := range []struct{path,caseID,sha,id string}{
		{firstPath,firstCase,firstHash,dossier.FirstContext.Entries[0].EntryID},
		{secondPath,secondCase,secondHash,dossier.SecondContext.Entries[0].EntryID},
	} {
		owner,ids,err := loadAuthorizedImportedCase(context.Background(),"analyst-a",
			item.path,item.caseID,item.sha)
		if err != nil { t.Fatal(err) }
		if _,ok:=ids[item.id];!ok { t.Fatalf("context ref is not a canonical owner import: %q",item.id) }
		opened,err:=inspectDemoEntryWithData(context.Background(),"analyst-a",owner,item.caseID,item.id)
		if err!=nil || opened.CaseID!=item.caseID || opened.EntryID!=item.id || !opened.Synthetic {
			t.Fatalf("context link cannot reopen with independent source authorization: %+v %v",opened,err)
		}
	}
	for _,secret := range []string{
		"PRIVATE-GENERIC-FIRST-CONTEXT","PRIVATE-GENERIC-SECOND-CONTEXT",
		"PRIVATE-GENERIC-OUTSIDE-WINDOW","Bearer ","Local analyst viewer:",
	} {
		if bytes.Contains(raw,[]byte(secret)) {
			t.Fatalf("paired report copied raw generic source or session secret: %q",secret)
		}
	}
	// An unrelated, unreadable or changed second owner may never cause the
	// first Case source preview to be partially published.
	changed,err:=os.ReadFile(secondPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(secondPath,append(changed,'\n'),0600);err!=nil {t.Fatal(err)}
	denied:=filepath.Join(t.TempDir(),"no-partial.json")
	if err:=writeVerifiedCorrelatedPairPrivate(context.Background(),"analyst-a",
		firstPath,firstCase,firstHash,secondPath,secondCase,secondHash,1,denied);err==nil {
		t.Fatal("tampered second owner allowed partial correlated context report")
	}
	if _,err:=os.Lstat(denied);!os.IsNotExist(err) {
		t.Fatal("unauthorized or changed second Case published an output")
	}
}
