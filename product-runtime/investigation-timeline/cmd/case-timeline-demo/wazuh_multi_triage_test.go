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

func TestMultiWazuhCaseTriageThroughOwnerAuthorization(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"rule\":{\"description\":\"confidential sensor A\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.123\",\"timestamp\":\"2026-10-01T10:02:02.130Z\",\"rule\":{\"description\":\"confidential sensor B\"}}\n")
	before := map[string][]byte{}
	for _, path := range []string{owner, first, second} {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[path] = b
	}
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+owner,
		"-wazuh-alerts="+first, "-wazuh-alerts="+second,
		"-into-case=case-a", "-triage-import",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var receipt localImportedTriage
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.Format != "cmdr-local-multi-wazuh4-case-triage-v1" ||
		!receipt.Synthetic || receipt.ImportedEntries != 2 ||
		receipt.Case.CaseID != "case-a" || receipt.Case.VisibleEntries < 2 ||
		receipt.Case.Kinds["observed"] < 2 || receipt.Case.FirstCertainUTC == "" ||
		receipt.Case.LastCertainUTC == "" || receipt.Case.DistinctSources < 2 {
		t.Fatalf("incorrect authorized multi-feed triage: %+v", receipt)
	}
	for _, secret := range []string{"confidential sensor", "full_log", "Bearer ", "demo-hidden-source", "demo-tenant-b", first, second} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("sensitive raw feed/owner information leaked: %q", secret)
		}
	}
	for path, content := range before {
		after, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(after, content) {
			t.Fatalf("operator file modified: %s: %v", path, err)
		}
	}
}

func TestMultiWazuhTriageFailClosedBeforeOutput(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	first := writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"rule\":{\"description\":\"private\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"b.2\",\"timestamp\":\"2026-10-01T10:02:02.130Z\",\"rule\":{\"description\":\"private 2\"}}\n")
	bad := writeJSONLEvents(t, "{not-json}\n")
	for _, tc := range []struct {
		name, principal, caseID string
		files []string
	}{
		{"cross-tenant-principal", "analyst-b", "case-a", []string{first, second}},
		{"foreign-case", "analyst-a", "case-b", []string{first, second}},
		{"missing-case", "analyst-a", "does-not-exist", []string{first, second}},
		{"malformed-second-feed", "analyst-a", "case-a", []string{first, bad}},
		{"duplicate-file", "analyst-a", "case-a", []string{first, first}},
		{"not-multi", "analyst-a", "case-a", []string{first}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := writeDemoMultiWazuhTriage(context.Background(), tc.principal, owner, tc.caseID, tc.files, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid multi-feed triage emitted output: %v / %q", err, out.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	err := writeDemoMultiWazuhTriage(ctx, "analyst-a", owner, "case-a", []string{first, second}, &out)
	if !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled triage emitted output: %v / %q", err, out.String())
	}
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot="+owner, "-into-case=case-a", "-wazuh-alerts="+first, "-triage-import"},
		{"-principal=analyst-a", "-snapshot="+owner, "-into-case=case-a", "-wazuh-alerts="+first, "-wazuh-alerts="+second, "-triage-import", "-compare-proximity"},
		{"-principal=analyst-a", "-snapshot="+owner, "-into-case=case-a", "-wazuh-alerts="+first, "-wazuh-alerts="+second, "-triage-import", "-with-case=case-c"},
	} {
		out.Reset()
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("ambiguous CLI leaked data: %v / %q", err, out.String())
		}
	}
}
