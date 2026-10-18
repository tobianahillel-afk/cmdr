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

func buildTwoVerifiedLocalCaseArchives(t *testing.T) (string,string,string,string,string,string) {
	t.Helper()
	const firstCase="forensic-case-on-disk"
	const secondCase="forensic-case-on-disk-secondary"
	original:=snapshotFixture(t)
	b,err:=loadDemoSnapshot(original)
	if err!=nil {t.Fatal(err)}
	// Create a second, independent first-party Case owner based on the
	// same fixture, with the Case and its references consistently renamed.
	found:=false
	for i:=range b.Cases {
		if b.Cases[i].ID==firstCase {b.Cases[i].ID=secondCase;found=true}
	}
	if !found {t.Fatal("missing first owner Case in fixture")}
	for i:=range b.Entries {
		if b.Entries[i].CaseID==firstCase {b.Entries[i].CaseID=secondCase}
	}
	for i:=range b.Sources {
		if b.Sources[i].Identity.CaseID==firstCase {b.Sources[i].Identity.CaseID=secondCase}
	}
	if _,err:=newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",b);err!=nil {
		t.Fatalf("second independent Case fixture not valid: %v",err)
	}
	secondInput:=filepath.Join(t.TempDir(),"independent-second-owner.json")
	serialized,err:=json.Marshal(b)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(secondInput,serialized,0600);err!=nil {t.Fatal(err)}
	feedOne:=writeJSONLEvents(t,
		"{\"id\":\"case-one-observation\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"message\":\"PowerShell suspicious endpoint activity\"}\n")
	feedTwo:=writeJSONLEvents(t,
		"{\"id\":\"case-two-observation\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"message\":\"Firewall blocked outbound activity\"}\n")
	aPath:=filepath.Join(t.TempDir(),"a-authorized-case.json")
	bPath:=filepath.Join(t.TempDir(),"b-authorized-case.json")
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",original,firstCase,
		[]string{feedOne},aPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writeDemoPersistedNormalizedCase(context.Background(),"analyst-a",secondInput,secondCase,
		[]string{feedTwo},bPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	aData,err:=os.ReadFile(aPath)
	if err!=nil {t.Fatal(err)}
	bData,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	return aPath,firstCase,fmt.Sprintf("%x",sha256.Sum256(aData)),
		bPath,secondCase,fmt.Sprintf("%x",sha256.Sum256(bData))
}

func TestCrossSavedCasesAggregatesTwoIndependentlyAuthorizedOwners(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash:=buildTwoVerifiedLocalCaseArchives(t)
	var result bytes.Buffer
	if err:=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,
		"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,
		"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
	},&result);err!=nil {t.Fatalf("two independently verified Case owners cannot compare: %v",err)}
	var decoded struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		First struct {
			CaseID string `json:"case_id"`
			VisibleEntries int `json:"visible_entries"`
		} `json:"first"`
		Second struct {
			CaseID string `json:"case_id"`
			VisibleEntries int `json:"visible_entries"`
		} `json:"second"`
		Shared []struct { Start string `json:"start_utc"` } `json:"shared_certain_five_minute_windows"`
	}
	if err:=json.Unmarshal(result.Bytes(),&decoded);err!=nil {t.Fatal(err)}
	if decoded.Format!="cmdr-local-verified-saved-case-comparison-v1" ||
		!decoded.Synthetic || decoded.First.CaseID!=aCase ||
		decoded.Second.CaseID!=bCase || decoded.First.VisibleEntries==0 ||
		decoded.Second.VisibleEntries==0 || len(decoded.Shared)==0 {
		t.Fatalf("Case correlation did not run over both owners: %+v",decoded)
	}
	for _,forbidden:=range []string{
		"PowerShell suspicious endpoint activity","Firewall blocked outbound activity",
		"Bearer ","Local analyst viewer:","127.0.0.1:",
	} {
		if strings.Contains(result.String(),forbidden) {
			t.Fatalf("cross-Case aggregate leaked source content or listener: %q",forbidden)
		}
	}
	// The original input JSONL files are irrelevant to a comparison
	// between two separately saved authorized owner snapshots.
}

func TestCrossSavedCasesFailsBeforeAnyReceiptOnSecondOwnerDenial(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash:=buildTwoVerifiedLocalCaseArchives(t)
	run:=func(principal,firstHash,secondHash,caseTwo string) (string,error) {
		var output bytes.Buffer
		err:=compareVerifiedSavedCases(context.Background(),principal,aPath,aCase,firstHash,
			bPath,caseTwo,secondHash,&output)
		return output.String(),err
	}
	for _,tc:=range []struct{name,principal,firstHash,secondHash,caseTwo string}{
		{"wrong second SHA","analyst-a",aHash,strings.Repeat("0",64),bCase},
		{"wrong first SHA","analyst-a",strings.Repeat("0",64),bHash,bCase},
		{"second Case absent","analyst-a",aHash,bHash,"missing-case"},
		{"principal denied","foreign-analyst",aHash,bHash,bCase},
		{"same Case IDs","analyst-a",aHash,bHash,aCase},
		{"missing second digest","analyst-a",aHash,"",bCase},
	} {
		t.Run(tc.name,func(t *testing.T) {
			output,err:=run(tc.principal,tc.firstHash,tc.secondHash,tc.caseTwo)
			if err==nil || output!="" {
				t.Fatalf("unauthorized/incomplete cross-snapshot comparison produced output: %v %q",err,output)
			}
		})
	}
	// Change the second file while keeping it parseable; the second
	// digest check must fail before writing any aggregates from the first.
	raw,err:=os.ReadFile(bPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,append(raw,'\n'),0600);err!=nil {t.Fatal(err)}
	receipt,err:=run("analyst-a",aHash,bHash,bCase)
	if err==nil || receipt!="" {
		t.Fatalf("tampered second owner emitted partial first owner report: %v %q",err,receipt)
	}
}

func TestCrossSavedCasesNeverComparesDifferentTenantScopes(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,_:=buildTwoVerifiedLocalCaseArchives(t)
	second,err:=loadDemoSnapshot(bPath)
	if err!=nil {t.Fatal(err)}
	second.Principals[0].Tenant="foreign-synthetic-tenant"
	second.Cases[0].Tenant="foreign-synthetic-tenant"
	for i:=range second.Entries {second.Entries[i].Value.Tenant="foreign-synthetic-tenant"}
	for i:=range second.Sources {second.Sources[i].Identity.Tenant="foreign-synthetic-tenant"}
	if _,err:=newDemoHandlerWithData(strings.Repeat("0",64),"analyst-a",second);err!=nil {
		t.Fatalf("other-tenant input must remain internally authorized: %v",err)
	}
	content,err:=json.Marshal(second)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(bPath,content,0600);err!=nil {t.Fatal(err)}
	bHash:=fmt.Sprintf("%x",sha256.Sum256(content))
	var result bytes.Buffer
	err=compareVerifiedSavedCases(context.Background(),"analyst-a",aPath,aCase,aHash,
		bPath,bCase,bHash,&result)
	if err==nil || result.Len()!=0 {
		t.Fatalf("two differently tenanted valid owners compared: %v %q",err,result.String())
	}
}

func TestCrossSavedCaseComparisonCLIRejectsUnrelatedModes(t *testing.T) {
	aPath,aCase,aHash,bPath,bCase,bHash:=buildTwoVerifiedLocalCaseArchives(t)
	base:=[]string{
		"-principal","analyst-a","-snapshot",aPath,
		"-compare-saved-case",aCase,"-expect-snapshot-sha256",aHash,
		"-compare-saved-with-snapshot",bPath,"-compare-saved-with-case",bCase,
		"-expect-other-snapshot-sha256",bHash,
	}
	for _,bad:=range [][]string{
		{"-serve-import"}, {"-report-case",aCase}, {"-compare-case",aCase},
		{"-save-import-snapshot",filepath.Join(t.TempDir(),"bad.json")},
		{"-import-jsonl",bPath}, {"-resume-import-search-case",aCase},
	} {
		var output bytes.Buffer
		args:=append(append([]string(nil),base...),bad...)
		if err:=runDemoCLI(context.Background(),args,&output);err==nil || output.Len()!=0 {
			t.Fatalf("unrelated CLI mode accepted in verified comparison: %v %q",bad,output.String())
		}
	}
}
