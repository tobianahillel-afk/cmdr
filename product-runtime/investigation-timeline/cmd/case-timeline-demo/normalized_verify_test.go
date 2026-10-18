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
