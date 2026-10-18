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

func TestExplicitPrivatePairIncludesAuthorizedGenericSourcePreviewsOnlyByConsent(t *testing.T) {
	firstOriginal,firstCase,_,secondOriginal,secondCase,_,_:=makeSavedPairDrilldownOwners(t)
	aFeed:=writeJSONLEvents(t,
		"{\"id\":\"near-first\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"SENSITIVE-GENERIC-FIRST-ACCEPTED\"}\n"+
		"{\"id\":\"far-first\",\"timestamp\":\"2026-10-17T10:04:00Z\",\"message\":\"SENSITIVE-GENERIC-OUTSIDE-WINDOW\"}\n")
	bFeed:=writeJSONLEvents(t,
		"{\"id\":\"near-second\",\"timestamp\":\"2026-10-17T10:00:22Z\",\"message\":\"SENSITIVE-GENERIC-SECOND-ACCEPTED\"}\n")
	aPath:=filepath.Join(t.TempDir(),"first-private-owner.json")
	bPath:=filepath.Join(t.TempDir(),"second-private-owner.json")
	for _,x:=range []struct{source,caseID,feed,dest string}{
		{firstOriginal,firstCase,aFeed,aPath},{secondOriginal,secondCase,bFeed,bPath},
	}{
		if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",
			x.source,x.caseID,[]string{x.feed},x.dest,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	}
	hash:=func(path string)string{
		raw,err:=os.ReadFile(path);if err!=nil{t.Fatal(err)}
		return fmt.Sprintf("%x",sha256.Sum256(raw))
	}
	aHash,bHash:=hash(aPath),hash(bPath)
	args:=[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",firstCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",secondCase,
		"-expect-other-snapshot-sha256",bHash,"-inspect-saved-pair","1",
	}
	private:=filepath.Join(t.TempDir(),"explicit-sensitive-evidence.json")
	var stdout bytes.Buffer
	withOptIn:=append(append([]string(nil),args...),
		"-include-pair-generic-evidence","-private-report",private)
	if err:=runDemoCLI(context.Background(),withOptIn,&stdout);err!=nil {
		t.Fatalf("explicit full source-owned context dossier failed: %v",err)
	}
	if stdout.Len()!=0 {t.Fatalf("private evidence copied into stdout: %q",stdout.String())}
	raw,err:=os.ReadFile(private)
	if err!=nil {t.Fatal(err)}
	var report struct {
		Format string `json:"format"`
		FirstContext authorizedGenericPairContextSummary `json:"first_generic_context"`
		SecondContext authorizedGenericPairContextSummary `json:"second_generic_context"`
		FirstEvidence []localEntryInspection `json:"first_generic_evidence"`
		SecondEvidence []localEntryInspection `json:"second_generic_evidence"`
	}
	if err:=json.Unmarshal(raw,&report);err!=nil{t.Fatal(err)}
	if report.Format!="cmdr-private-verified-wazuh-pair-v3" ||
		len(report.FirstEvidence)!=2 || len(report.SecondEvidence)!=1 ||
		len(report.FirstContext.Entries)!=len(report.FirstEvidence) ||
		len(report.SecondContext.Entries)!=len(report.SecondEvidence) {
		t.Fatalf("full context evidence or version missing: %+v",report)
	}
	for _,side:=range []struct{caseID string; contexts []authorizedGenericPairContext; inspections []localEntryInspection}{
		{firstCase,report.FirstContext.Entries,report.FirstEvidence},
		{secondCase,report.SecondContext.Entries,report.SecondEvidence},
	}{
		for i,opened:=range side.inspections {
			if !opened.Synthetic || opened.Kind!="telemetry-event" ||
				opened.CaseID!=side.caseID || opened.EntryID!=side.contexts[i].EntryID {
				t.Fatalf("source preview escaped its independent Case owner: %+v",opened)
			}
		}
	}
	if !bytes.Contains(raw,[]byte("SENSITIVE-GENERIC-FIRST-ACCEPTED")) ||
		!bytes.Contains(raw,[]byte("SENSITIVE-GENERIC-SECOND-ACCEPTED")) ||
		bytes.Contains(raw,[]byte("SENSITIVE-GENERIC-OUTSIDE-WINDOW")) ||
		bytes.Contains(raw,[]byte("Bearer ")) {
		t.Fatal("private opted-in source selection wrong or secrets copied")
	}
	// The original operator command remains a metadata-only v2 dossier.
	ordinary:=filepath.Join(t.TempDir(),"original-behavior.json")
	stdout.Reset()
	if err:=runDemoCLI(context.Background(),
		append(append([]string(nil),args...),"-private-report",ordinary),&stdout);err!=nil {
		t.Fatalf("normal metadata-only pair inspection regressed: %v",err)
	}
	minimal,err:=os.ReadFile(ordinary);if err!=nil{t.Fatal(err)}
	if !bytes.Contains(minimal,[]byte("cmdr-private-verified-wazuh-pair-v2")) ||
		bytes.Contains(minimal,[]byte("SENSITIVE-GENERIC-FIRST-ACCEPTED")) ||
		bytes.Contains(minimal,[]byte("SENSITIVE-GENERIC-SECOND-ACCEPTED")) {
		t.Fatal("generic preview exported without explicit consent")
	}
	// Explicit opt-in never enables a stdout report or a comparison-only
	// workflow; no partial file may be produced if other Case changed.
	for _,bad:=range [][]string{
		append(append([]string(nil),args...),"-include-pair-generic-evidence"),
		append(append([]string(nil),args[:len(args)-2]...),"-include-pair-generic-evidence",
			"-private-report",filepath.Join(t.TempDir(),"reject-no-pair.json")),
	} {
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),bad,&out);err==nil || out.Len()!=0 {
			t.Fatalf("generic preview consent bypassed strict private-pair mode: %v %q",bad,out.String())
		}
	}
	second,err:=os.ReadFile(bPath);if err!=nil{t.Fatal(err)}
	if err:=os.WriteFile(bPath,append(second,'\n'),0600);err!=nil{t.Fatal(err)}
	failed:=filepath.Join(t.TempDir(),"must-not-publish.json")
	var refused bytes.Buffer
	changedArgs:=append(append([]string(nil),args...),
		"-include-pair-generic-evidence","-private-report",failed)
	if err:=runDemoCLI(context.Background(),changedArgs,&refused);err==nil || refused.Len()!=0 {
		t.Fatal("invalid second Case SHA permitted partial sensitive preview")
	}
	if _,err:=os.Lstat(failed);!os.IsNotExist(err) {
		t.Fatal("private dossier published after independent owner verification failure")
	}
}

func TestGenericEvidenceReopeningRequiresSelectedOwnerAndCase(t *testing.T) {
	owner,err:=loadDemoSnapshot(snapshotFixture(t));if err!=nil{t.Fatal(err)}
	const caseID="forensic-case-on-disk"
	owner,err=importLocalEvents(context.Background(),owner,"analyst-a",caseID,[]byte(
		"{\"id\":\"owned\",\"timestamp\":\"2026-10-17T10:00:05Z\",\"message\":\"PRIVATE SOURCE\"}\n"))
	if err!=nil{t.Fatal(err)}
	ids,err:=importedNormalizedIDs(owner,caseID);if err!=nil{t.Fatal(err)}
	ctx,err:=nearbyAuthorizedGenericPairContext(context.Background(),"analyst-a",
		owner,caseID,ids,"2026-10-17T10:00:00Z")
	if err!=nil || len(ctx.Entries)!=1 {t.Fatalf("missing real owned source: %v",err)}
	if _,err:=reopenGenericContextEvidence(context.Background(),"foreign-analyst",
		owner,caseID,ctx);err==nil {t.Fatal("foreign principal reopened private generic source")}
	if _,err:=reopenGenericContextEvidence(context.Background(),"analyst-a",
		owner,"wrong-case",ctx);err==nil {t.Fatal("foreign Case reopened private generic source")}
	// A context ID is a navigation hint; no source is opened without
	// repeating its own first-party typed-source owner authorization.
	invalid:=ctx
	invalid.Entries=[]authorizedGenericPairContext{{EntryID:"unrelated-foreign-event"}}
	if _,err:=reopenGenericContextEvidence(context.Background(),"analyst-a",
		owner,caseID,invalid);err==nil{t.Fatal("arbitrary user-selected context ref opened private source")}
	_ = strings.Contains
}
