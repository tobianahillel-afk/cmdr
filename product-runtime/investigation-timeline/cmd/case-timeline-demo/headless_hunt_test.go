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

func savedHeadlessMixedCase(t *testing.T) (string, string, string) {
	t.Helper()
	const caseID = "forensic-case-on-disk"
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"w-one\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-ONE\",\"level\":7}}\n"+
		"{\"id\":\"w-two\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-TWO\",\"level\":7}}\n")
	generic := writeJSONLEvents(t,
		"{\"id\":\"g-near\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PRIVATE-GENERIC-NEAR\"}\n")
	file := filepath.Join(t.TempDir(), "private-mixed-case.json")
	if err := writePersistedMixedCase(context.Background(), "analyst-a", snapshotFixture(t),
		caseID, []string{wazuh}, []string{generic}, file, &bytes.Buffer{}); err != nil {
		t.Fatalf("real mixed Case save: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil { t.Fatal(err) }
	return file, caseID, fmt.Sprintf("%x", sha256.Sum256(data))
}

func TestHeadlessMixedHuntExecutesRealOwnerSearchAndTemporalCorrelation(t *testing.T) {
	file, caseID, digest := savedHeadlessMixedCase(t)
	var out bytes.Buffer
	args := []string{
		"-principal", "analyst-a", "-snapshot", file,
		"-hunt-saved-case", caseID, "-expect-snapshot-sha256", digest,
		"-hunt-field", "source_ref", "-hunt-family", "all",
		"-hunt-value", "local-jsonl-source-",
		"-hunt-from", "2026-10-17T09:59:00Z", "-hunt-to", "2026-10-17T10:02:00Z",
	}
	if err := runDemoCLI(context.Background(), args, &out); err != nil {
		t.Fatalf("real CLI headless mixed Case hunt unavailable: %v", err)
	}
	var result struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		JobID string `json:"job_id"`
		Refs []struct{ Source, Ref string } `json:"authorized_refs"`
		WazuhHits []struct {EntryID string `json:"entry_id"`} `json:"wazuh_hits"`
		Mixed []mixedTemporalLead `json:"mixed_temporal_context"`
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("headless hunt has no JSON receipt: %v %q", err, out.String())
	}
	if result.Format != "cmdr-local-headless-authorized-hunt-v1" ||
		!result.Synthetic || result.CaseID != caseID ||
		result.JobID == "" || len(result.Refs) != 3 ||
		len(result.WazuhHits) != 2 || len(result.Mixed) != 2 {
		t.Fatalf("headless hunt did not integrate authorized SearchJob and temporal leads: %+v", result)
	}
	members := make(map[string]bool)
	for _, ref := range result.Refs {
		if ref.Source != localJSONLSearchSource || ref.Ref == "" || members[ref.Ref] {
			t.Fatalf("headless hunt disclosed foreign or repeated source: %+v", ref)
		}
		members[ref.Ref] = true
	}
	for _, lead := range result.Mixed {
		if !members[lead.WazuhEntryID] || len(lead.GenericEvents) == 0 ||
			lead.WindowSeconds != 60 {
			t.Fatalf("unlinked or malformed mixed Case context: %+v", lead)
		}
		for _, neighbor := range lead.GenericEvents {
			if !members[neighbor.EntryID] || neighbor.OffsetSeconds > 60 ||
				neighbor.OffsetSeconds < -60 {
				t.Fatalf("cross-family lead escaped authorized Job: %+v", neighbor)
			}
		}
	}
	for _, secret := range []string{
		"PRIVATE-WAZUH-ONE", "PRIVATE-WAZUH-TWO",
		"PRIVATE-GENERIC-NEAR", "Bearer token", "Local analyst viewer:",
		"127.0.0.1:", "\"audit_trail\"", "\"preview\"",
	} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("headless receipt disclosed source data/owner internals: %q", secret)
		}
	}
	// Same real Case, existing Event Search implementation: restrictive
	// family cannot emit mixed cross-family contexts.
	out.Reset()
	if err := runSavedCaseHeadlessHunt(context.Background(), "analyst-a", file,
		caseID, digest, "source_ref", "wazuh", "local-jsonl-source-",
		"2026-10-17T09:59:00Z", "2026-10-17T10:02:00Z", &out); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "mixed_temporal_context") {
		t.Fatalf("Wazuh-only hunt leaked cross-family correlation: %q", out.String())
	}
}

func TestHeadlessCaseHuntRejectsUnverifiedOrUnauthorizedInputWithNoPartialReceipt(t *testing.T) {
	file, caseID, digest := savedHeadlessMixedCase(t)
	for _, scenario := range []struct {
		name, principal, caseID, sha, field, family, query, from, to string
	}{
		{"wrong hash","analyst-a",caseID,strings.Repeat("0",64),"preview","all","Wazuh","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
		{"missing hash","analyst-a",caseID,"","preview","all","Wazuh","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
		{"wrong Case","analyst-a","foreign-case",digest,"preview","all","Wazuh","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
		{"foreign analyst","foreign-analyst",caseID,digest,"preview","all","Wazuh","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
		{"invalid family","analyst-a",caseID,digest,"preview","foreign","Wazuh","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
		{"out of order time","analyst-a",caseID,digest,"preview","all","Wazuh","2026-10-17T10:02:00Z","2026-10-17T09:59:00Z"},
		{"inapplicable pivot","analyst-a",caseID,digest,"rule_id","generic","5710","2026-10-17T09:59:00Z","2026-10-17T10:02:00Z"},
	} {
		t.Run(scenario.name,func(t *testing.T) {
			var output bytes.Buffer
			err := runSavedCaseHeadlessHunt(context.Background(),scenario.principal,file,
				scenario.caseID,scenario.sha,scenario.field,scenario.family,scenario.query,
				scenario.from,scenario.to,&output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("invalid headless Case hunt leaked partial evidence: %v %q",err,output.String())
			}
		})
	}
	// A parseable modification of the authorized file must not silently
	// change the evidence searched against the operator-held digest.
	original,err:=os.ReadFile(file)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(file,append(original,'\n'),0600);err!=nil {t.Fatal(err)}
	var output bytes.Buffer
	if err:=runSavedCaseHeadlessHunt(context.Background(),"analyst-a",file,
		caseID,digest,"preview","all","Wazuh",
		"2026-10-17T09:59:00Z","2026-10-17T10:02:00Z",&output);
		err==nil || output.Len()!=0 {
		t.Fatalf("tampered but parseable Case produced a search receipt: %q",output.String())
	}
}

func TestHeadlessCaseHuntCLIRejectsMixedModes(t *testing.T) {
	file, caseID, digest := savedHeadlessMixedCase(t)
	args:=[]string{"-principal","analyst-a","-snapshot",file,
		"-hunt-saved-case",caseID,"-expect-snapshot-sha256",digest,
		"-hunt-value","local-jsonl-source-","-hunt-field","source_ref",
		"-hunt-from","2026-10-17T09:59:00Z","-hunt-to","2026-10-17T10:02:00Z"}
	for _,extra:=range [][]string{
		{"-serve-import"}, {"-resume-import-search-case",caseID},
		{"-export-case",caseID}, {"-compare-saved-case",caseID},
		{"-hunt-family","all","-report-case",caseID},
	} {
		var out bytes.Buffer
		parameters:=append(append([]string(nil),args...),extra...)
		if err:=runDemoCLI(context.Background(),parameters,&out);err==nil || out.Len()!=0 {
			t.Fatalf("headless hunt combined with other mode: %v %q",extra,out.String())
		}
	}
	var out bytes.Buffer
	if err:=runDemoCLI(context.Background(),[]string{
		"-snapshot",file,"-hunt-value","sensitive",
	},&out);err==nil || out.Len()!=0 {
		t.Fatalf("orphan headless search selection accepted: %q",out.String())
	}
}
