package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestWazuh4OneCaseTriageAndTwoCaseAggregateFromDisk(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"wazuh-investigation-777\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"description\":\"SENSITIVE Wazuh rule preview\"},\"full_log\":\"PRIVATE RAW LOG\"}\n")
	originalOwner, err := os.ReadFile(owner)
	if err != nil { t.Fatal(err) }
	originalAlerts, err := os.ReadFile(alerts)
	if err != nil { t.Fatal(err) }

	base := []string{"-principal=analyst-a", "-snapshot="+owner,
		"-wazuh-alerts="+alerts, "-into-case=case-a", "-triage-import"}
	var single bytes.Buffer
	if err := runDemoCLI(context.Background(), base, &single); err != nil {
		t.Fatal(err)
	}
	var triage localImportedTriage
	if err := json.Unmarshal(single.Bytes(), &triage); err != nil { t.Fatal(err) }
	if triage.Format != "cmdr-local-wazuh4-case-triage-v1" || !triage.Synthetic ||
		triage.ImportedEntries != 1 || triage.Case.CaseID != "case-a" ||
		triage.Case.VisibleEntries != 208 || triage.Case.Kinds["observed"] != 105 {
		t.Fatalf("Wazuh single Case did not reach authorized triage: %+v", triage)
	}

	var joint bytes.Buffer
	args := append(append([]string{}, base...), "-with-case=case-c")
	if err := runDemoCLI(context.Background(), args, &joint); err != nil { t.Fatal(err) }
	var comparison localImportedComparison
	if err := json.Unmarshal(joint.Bytes(), &comparison); err != nil { t.Fatal(err) }
	if comparison.Format != "cmdr-local-wazuh4-case-comparison-v1" || !comparison.Synthetic ||
		comparison.ImportedEntries != 1 ||
		comparison.Comparison.First.CaseID != "case-a" ||
		comparison.Comparison.Second.CaseID != "case-c" ||
		comparison.Comparison.First.VisibleEntries != 208 ||
		comparison.Comparison.Second.VisibleEntries != 1 ||
		comparison.Comparison.VisibleEntryDelta != 207 ||
		comparison.Comparison.CertainTimeOverlap == nil {
		t.Fatalf("two independently authorized Case aggregates missing: %+v", comparison)
	}
	for _, output := range []string{single.String(),joint.String()} {
		for _, secret := range []string{
			"SENSITIVE Wazuh rule preview", "PRIVATE RAW LOG",
			"SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT", "demo-hidden-source",
			"demo-tenant-b", "Bearer ", "principals", "wazuh-investigation-777",
		} {
			if strings.Contains(output,secret) {t.Fatalf("Wazuh report disclosed %q",secret)}
		}
	}
	if after,err:=os.ReadFile(owner); err!=nil || !bytes.Equal(originalOwner,after) {
		t.Fatal("owner file changed during forensic triage")
	}
	if after,err:=os.ReadFile(alerts); err!=nil || !bytes.Equal(originalAlerts,after) {
		t.Fatal("original Wazuh alert file was modified")
	}
}

func TestWazuh4OfflineTriageAuthorizationIsAtomic(t *testing.T) {
	owner:=twoCaseComparisonSnapshot(t)
	alerts:=writeJSONLEvents(t,
		"{\"id\":\"wazuh-investigation-777\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"description\":\"hidden Wazuh telemetry\"}}\n")
	bad:=writeJSONLEvents(t,"{\"@timestamp\":\"2026-10-01T10:30:00Z\",\"message\":\"unsupported ECS\"}\n")
	base:=[]string{"-principal=analyst-a","-snapshot="+owner,
		"-wazuh-alerts="+alerts,"-into-case=case-a","-triage-import"}
	cases:=[]struct{name string; mutate func([]string)[]string}{
		{"foreign-second",func(a []string)[]string{return append(a,"-with-case=case-b")}},
		{"missing-second",func(a []string)[]string{return append(a,"-with-case=absent")}},
		{"same-case",func(a []string)[]string{return append(a,"-with-case=case-a")}},
		{"window-without-proximity",func(a []string)[]string{return append(a,"-with-case=case-c","-proximity-seconds=3")}},
		{"proximity-without-second",func(a []string)[]string{return append(a,"-compare-proximity")}},
		{"denied-principal",func(a []string)[]string{return append(a,"-principal=analyst-b")}},
		{"invalid-alert",func(a []string)[]string{return append(a,"-wazuh-alerts="+bad)}},
		{"ambiguous-viewer",func(a []string)[]string{return append(a,"-serve-import")}},
	}
	for _,tc:=range cases {
		t.Run(tc.name,func(t *testing.T){
			var out bytes.Buffer
			if err:=runDemoCLI(context.Background(),tc.mutate(append([]string{},base...)),&out); err==nil || out.Len()!=0 {
				t.Fatalf("denied mode emitted output or passed: %v / %q",err,out.String())
			}
		})
	}
	ctx,cancel:=context.WithCancel(context.Background());cancel()
	for _,second:=range []string{"","case-c"}{
		var out bytes.Buffer
		var err error
		if second=="" {err=writeDemoWazuh4Triage(ctx,"analyst-a",owner,"case-a",alerts,&out)
		}else{err=writeDemoWazuh4Comparison(ctx,"analyst-a",owner,"case-a",second,alerts,&out)}
		if !errors.Is(err,context.Canceled)||out.Len()!=0 {t.Fatalf("cancelled Wazuh analysis emitted: %v",err)}
	}
}
