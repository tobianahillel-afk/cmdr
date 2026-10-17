package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestMultiWazuhFilesProduceAuthorizedTwoCaseComparison(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	a := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"rule\":{\"description\":\"sensitive feed a\"}}\n")
	b := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.130Z\",\"rule\":{\"description\":\"sensitive feed b\"}}\n")
	original := make(map[string][]byte)
	for _, path := range []string{owner, a, b} {
		raw, err := os.ReadFile(path)
		if err != nil { t.Fatal(err) }
		original[path] = raw
	}
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+owner, "-into-case=case-a",
		"-wazuh-alerts="+a, "-wazuh-alerts="+b,
		"-triage-import", "-with-case=case-c",
	}, &output); err != nil { t.Fatal(err) }
	var receipt localImportedComparison
	if err := json.Unmarshal(output.Bytes(), &receipt); err != nil { t.Fatal(err) }
	if receipt.Format != "cmdr-local-multi-wazuh4-case-comparison-v1" ||
		!receipt.Synthetic || receipt.ImportedEntries != 2 ||
		receipt.Comparison.First.CaseID != "case-a" ||
		receipt.Comparison.Second.CaseID != "case-c" ||
		receipt.Comparison.First.Kinds["observed"] < 2 ||
		receipt.Comparison.First.DistinctSources < 2 {
		t.Fatalf("missing authorized multi-feed comparison: %+v", receipt)
	}
	for _, forbidden := range []string{"sensitive feed", "Bearer ", "demo-tenant-b", a, b} {
		if strings.Contains(output.String(), forbidden) {
			t.Fatalf("private source data in comparison receipt: %q", forbidden)
		}
	}
	for path, data := range original {
		after, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(data, after) { t.Fatalf("input changed: %q %v", path, err) }
	}
}

func TestMultiWazuhComparisonDeniedOrInvalidEmitsNothing(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	a := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02Z\",\"rule\":{\"description\":\"a\"}}\n")
	b := writeJSONLEvents(t, "{\"id\":\"b.2\",\"timestamp\":\"2026-10-01T10:02:03Z\",\"rule\":{\"description\":\"b\"}}\n")
	invalid := writeJSONLEvents(t, "{bad-json}\n")
	for _, tc := range []struct{
		name, principal, caseID, withCase string
		paths []string
	}{
		{"foreign-principal", "analyst-b", "case-a", "case-c", []string{a,b}},
		{"foreign-target", "analyst-a", "case-a", "case-b", []string{a,b}},
		{"missing-target", "analyst-a", "case-a", "absent", []string{a,b}},
		{"same-case", "analyst-a", "case-a", "case-a", []string{a,b}},
		{"malformed-second-feed", "analyst-a", "case-a", "case-c", []string{a,invalid}},
		{"duplicate-file", "analyst-a", "case-a", "case-c", []string{a,a}},
	}{
		t.Run(tc.name,func(t *testing.T){
			var output bytes.Buffer
			err := writeDemoMultiWazuhComparison(context.Background(),tc.principal,owner,tc.caseID,tc.withCase,tc.paths,&output)
			if err == nil || output.Len()!=0 { t.Fatalf("denied report leaked output: %v / %q",err,output.String()) }
		})
	}
}
