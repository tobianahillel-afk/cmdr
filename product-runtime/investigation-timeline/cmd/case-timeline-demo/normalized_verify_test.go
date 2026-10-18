package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

type rejectedOfflineReceiptWriter struct{}
func (rejectedOfflineReceiptWriter) Write([]byte) (int, error) {
	return 0, errors.New("operator output unavailable")
}

func TestOfflineAuthorizedCaseVerificationRealMultifeedRoundTrip(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	original := snapshotFixture(t)
	one := writeJSONLEvents(t,
		"{\"id\":\"duplicate\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"Unauthorized endpoint PowerShell execution\"}\n")
	two := writeJSONLEvents(t,
		"{\"id\":\"duplicate\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"Firewall blocked PowerShell callback\"}\n")
	file := filepath.Join(t.TempDir(), "private-case.json")
	var saved bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal","analyst-a","-snapshot",original,"-into-case",caseID,
		"-import-jsonl",one,"-import-jsonl-extra",two,
		"-save-import-snapshot",file,
	}, &saved); err != nil {
		t.Fatalf("real Case save failed: %v",err)
	}
	raw, err := os.ReadFile(file)
	if err != nil { t.Fatal(err) }
	expected := fmt.Sprintf("%x",sha256.Sum256(raw))
	if !strings.Contains(saved.String(),expected) {
		t.Fatal("saved Case receipt did not include actual on-disk digest")
	}
	// No HTTP loopback and no new bearer: a CLI read-only validation receipt
	// can finish on an uncancelled context, without any network listener.
	var answer bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal","analyst-a","-snapshot",file,
		"-verify-import-snapshot-case",caseID,"-expect-snapshot-sha256",expected,
	}, &answer); err != nil {
		t.Fatalf("offline verified Case failed: %v",err)
	}
	var result struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		Count int `json:"authorized_imported_events"`
		ByFeed map[string]int `json:"authorized_events_by_import_slot"`
	}
	if err := json.Unmarshal(answer.Bytes(),&result); err != nil {
		t.Fatalf("offline verification receipt is not JSON: %v %q",err,answer.String())
	}
	if result.Format!="cmdr-offline-authorized-case-verification-v1" ||
		!result.Synthetic || result.CaseID!=caseID || result.Count!=2 ||
		result.ByFeed["operator-local-jsonl-feed-1"]!=1 ||
		result.ByFeed["operator-local-jsonl-feed-2"]!=1 {
		t.Fatalf("incorrect Case authorization / feed provenance verification: %+v",result)
	}
	for _,secret := range []string{
		"PowerShell","Firewall blocked","local-jsonl-entry-",
		"Bearer","127.0.0.1:","source-owner",
	} {
		if strings.Contains(answer.String(),secret) {
			t.Fatalf("offline integrity receipt disclosed evidence or listener: %q",secret)
		}
	}
	// The persisted owner no longer requires the original inputs.
	if err:=os.Remove(one);err!=nil {t.Fatal(err)}
	if err:=os.Remove(two);err!=nil {t.Fatal(err)}
	answer.Reset()
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",file,caseID,expected,&answer);err!=nil {
		t.Fatalf("offline verification depends on deleted input JSONL: %v",err)
	}
}

func TestOfflineCaseVerificationFailClosedBeforeAnyReceipt(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot := snapshotFixture(t)
	feed := writeJSONLEvents(t,
		"{\"id\":\"only-event\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"sensitive investigator notes\"}\n")
	file := filepath.Join(t.TempDir(),"saved-owner.json")
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",snapshot,caseID,[]string{feed},file,io.Discard);err!=nil {
		t.Fatalf("prepare owner: %v",err)
	}
	contents,err:=os.ReadFile(file)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(contents))
	ctx,cancel:=context.WithCancel(context.Background())
	cancel()
	for _,tc:=range []struct {
		name string
		ctx context.Context
		principal, snapshot, caseID, digest string
	}{
		{"missing context",nil,"analyst-a",file,caseID,digest},
		{"cancelled context",ctx,"analyst-a",file,caseID,digest},
		{"missing snapshot",context.Background(),"analyst-a","",caseID,digest},
		{"missing Case",context.Background(),"analyst-a",file,"",digest},
		{"missing digest",context.Background(),"analyst-a",file,caseID,""},
		{"incorrect sha256",context.Background(),"analyst-a",file,caseID,strings.Repeat("0",64)},
		{"invalid sha256",context.Background(),"analyst-a",file,caseID,"invalid"},
		{"foreign principal",context.Background(),"foreign-analyst",file,caseID,digest},
		{"foreign Case",context.Background(),"analyst-a",file,"foreign-case",digest},
		{"missing file",context.Background(),"analyst-a",file+"-does-not-exist",caseID,digest},
	} {
		t.Run(tc.name,func(t *testing.T){
			var output bytes.Buffer
			err:=verifySavedNormalizedCase(tc.ctx,tc.principal,tc.snapshot,tc.caseID,tc.digest,&output)
			if err==nil || output.Len()!=0 {
				t.Fatalf("denied verification returned a receipt: %v %q",err,output.String())
			}
		})
	}
	// A no-op JSON change remains parseable, yet invalidates independently
	// held evidence digests. No verification receipt may escape.
	if err:=os.WriteFile(file,append(contents,'\n'),0600);err!=nil {t.Fatal(err)}
	if _,err:=loadDemoSnapshot(file);err!=nil {t.Fatalf("tamper should remain parseable: %v",err)}
	var denied bytes.Buffer
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",file,caseID,digest,&denied);err==nil || denied.Len()!=0 {
		t.Fatalf("tampered owner was verified: %q",denied.String())
	}
	// Restore bytes; a faulty output sink must be reported rather than
	// claiming that the operator received a durable validation receipt.
	if err:=os.WriteFile(file,contents,0600);err!=nil {t.Fatal(err)}
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",file,caseID,digest,rejectedOfflineReceiptWriter{});err==nil {
		t.Fatal("write failure falsely treated as successful offline verification")
	}
}

func TestOfflineCaseVerificationCLIExclusiveAndMandatoryFingerprint(t *testing.T) {
	const caseID="forensic-case-on-disk"
	snapshot:=snapshotFixture(t)
	for _,args:=range [][]string{
		{"-snapshot",snapshot,"-verify-import-snapshot-case",caseID},
		{"-snapshot",snapshot,"-verify-import-snapshot-case",caseID,"-resume-import-search-case",caseID,
		 "-expect-snapshot-sha256",strings.Repeat("0",64)},
		{"-snapshot",snapshot,"-verify-import-snapshot-case",caseID,"-serve-import",
		 "-expect-snapshot-sha256",strings.Repeat("0",64)},
		{"-snapshot",snapshot,"-verify-import-snapshot-case",caseID,"-into-case",caseID,
		 "-expect-snapshot-sha256",strings.Repeat("0",64)},
		{"-snapshot",snapshot,"-verify-import-snapshot-case",caseID,"-expect-snapshot-sha256","invalid"},
	} {
		var output bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&output);err==nil || output.Len()!=0 {
			t.Fatalf("invalid/no-digest offline verification mode accepted: %v %q",args,output.String())
		}
	}
}

func TestOfflineVerifierRefusesUnimportedButReadableCase(t *testing.T) {
	const caseID="forensic-case-on-disk"
	// Original owner snapshot is a valid and readable Case. It is not a
	// saved normalized import, so verification must refuse a misleading
	// forensic completeness receipt even when its SHA-256 is correct.
	path:=snapshotFixture(t)
	raw,err:=os.ReadFile(path)
	if err!=nil {t.Fatal(err)}
	digest:=fmt.Sprintf("%x",sha256.Sum256(raw))
	var out bytes.Buffer
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",path,caseID,digest,&out);err==nil || out.Len()!=0 {
		t.Fatalf("ordinary Case snapshot falsely verified as normalized import: %q",out.String())
	}
}

func TestPrivateCaseAtomicPublicationDeniesUnsafeFilesystemDestinations(t *testing.T) {
	dir:=t.TempDir()
	existing:=filepath.Join(dir,"original-forensic-evidence.json")
	const sentinel="DO NOT OVERWRITE EXISTING FORENSIC EVIDENCE"
	if err:=os.WriteFile(existing,[]byte(sentinel),0600);err!=nil {t.Fatal(err)}
	// A hardlink to the same evidence is not a new writable destination.
	link:=filepath.Join(dir,"hardlink.json")
	if err:=os.Link(existing,link);err!=nil {t.Fatalf("test hardlinks unavailable: %v",err)}
	for _,dst:=range []string{
		existing,
		link,
		dir,
		filepath.Join(dir,"directory-that-does-not-exist","new-case.json"),
	} {
		if err:=writePrivateSnapshotAtomic(context.Background(),dst,[]byte("new sensitive owner"));err==nil {
			t.Fatalf("unsafe private Case destination accepted: %q",dst)
		}
	}
	got,err:=os.ReadFile(existing)
	if err!=nil || string(got)!=sentinel {t.Fatal("existing forensic evidence modified")}
	for _,ent:=range []string{existing,link} {
		contents,err:=os.ReadFile(ent)
		if err!=nil || string(contents)!=sentinel {t.Fatalf("existing evidence hardlink modified %s: %v",ent,err)}
	}
	for _,entry:=range []string{
		"directory-that-does-not-exist",
	} {
		if _,err:=os.Stat(filepath.Join(dir,entry));!os.IsNotExist(err) {
			t.Fatal("failed write created a missing directory")
		}
	}
	// The caller cannot circumvent the bounded snapshot or choose relative
	// evidence paths. These checks run before any output is created.
	for _,test:=range []struct{ dest string; data []byte; ctx context.Context }{
		{filepath.Join(dir,"too-large.json"),bytes.Repeat([]byte("x"),maxDemoSnapshotBytes+1),context.Background()},
		{filepath.Join(dir,"empty.json"),nil,context.Background()},
		{"relative-owner.json",[]byte("sensitive"),context.Background()},
		{filepath.Join(dir,"nil-context.json"),[]byte("sensitive"),nil},
	} {
		if err:=writePrivateSnapshotAtomic(test.ctx,test.dest,test.data);err==nil {
			t.Fatalf("invalid private Case export accepted: %q",test.dest)
		}
		if filepath.IsAbs(test.dest) {
			if _,err:=os.Lstat(test.dest);!os.IsNotExist(err) {
				t.Fatalf("rejected path was created: %q",test.dest)
			}
		}
	}
	// Failed attempts must never leave a .partial file with sensitive bytes.
	entries,err:=os.ReadDir(dir)
	if err!=nil {t.Fatal(err)}
	for _,entry:=range entries {
		if strings.HasPrefix(entry.Name(),".cmdr-private-case-") {
			t.Fatalf("private owner tempfile leaked after a denied write: %q",entry.Name())
		}
	}
}

func TestOfflineVerificationCountsOnlyReadableImportedSources(t *testing.T) {
	const caseID="forensic-case-on-disk"
	original,err:=loadDemoSnapshot(snapshotFixture(t))
	if err!=nil {t.Fatal(err)}
	first,err:=importLocalEvents(context.Background(),original,"analyst-a",caseID,[]byte(
		"{\"id\":\"allowed\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"readable normalized evidence\"}\n"+
		"{\"id\":\"denied\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"message\":\"CONFIDENTIAL DENIED SOURCE CONTENT\"}\n"))
	if err!=nil {t.Fatal(err)}
	var owner timeline.MemoryTimelineCase
	for _,c:=range first.Cases {if c.ID==caseID {owner=c;break}}
	if owner.ID=="" {t.Fatal("missing selected Case fixture")}
	blockedPrincipal:="read-only-other-analyst"
	first.Principals=append(first.Principals,timeline.MemoryTimelinePrincipal{
		ID:blockedPrincipal,Tenant:owner.Tenant,Environment:owner.Environment,
	})
	deniedID:=""
	for i:=range first.Entries {
		if first.Entries[i].CaseID!=caseID ||
			first.Entries[i].Value.Provenance!="operator-local-jsonl" {
			continue
		}
		if first.Entries[i].Value.SourceTime=="2026-10-17T10:01:00Z" {
			deniedID=first.Entries[i].Value.ID
			first.Entries[i].Readers=[]string{blockedPrincipal}
		}
	}
	if deniedID=="" {t.Fatal("missing imported denied source")}
	for i:=range first.Sources {
		if first.Sources[i].Identity.Reference=="local-jsonl-source-"+strings.TrimPrefix(deniedID,"local-jsonl-entry-") {
			first.Sources[i].Readers=[]string{blockedPrincipal}
		}
	}
	if _,err:=newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",first);err!=nil {
		t.Fatalf("input Case with source-specific deny should remain owner-valid: %v",err)
	}
	raw,err:=json.Marshal(first)
	if err!=nil {t.Fatal(err)}
	file:=filepath.Join(t.TempDir(),"old-unminimized-snapshot.json")
	if err:=os.WriteFile(file,raw,0600);err!=nil {t.Fatal(err)}
	hash:=fmt.Sprintf("%x",sha256.Sum256(raw))
	var receipt bytes.Buffer
	if err:=verifySavedNormalizedCase(context.Background(),"analyst-a",file,caseID,hash,&receipt);err!=nil {
		t.Fatalf("valid readable projection should verify: %v",err)
	}
	var data struct {
		Imported int `json:"authorized_imported_events"`
		Slots map[string]int `json:"authorized_events_by_import_slot"`
	}
	if err:=json.Unmarshal(receipt.Bytes(),&data);err!=nil {t.Fatal(err)}
	if data.Imported!=1 || data.Slots["operator-local-jsonl"]!=1 {
		t.Fatalf("unreadable source affected offline authorized inventory: %+v",data)
	}
	for _,disallowed:=range []string{
		blockedPrincipal,deniedID,"CONFIDENTIAL DENIED SOURCE CONTENT","local-jsonl-source-",
	} {
		if strings.Contains(receipt.String(),disallowed) {
			t.Fatalf("offline receipt disclosed inaccessible source info: %q",disallowed)
		}
	}
}
