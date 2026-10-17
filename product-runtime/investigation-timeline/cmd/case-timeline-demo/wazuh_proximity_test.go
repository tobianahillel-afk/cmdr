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

func TestWazuh4AlertsToTwoAuthorizedCasesTemporalLeads(t *testing.T) {
	ownerPath := proximityDiskSnapshot(t)
	alertsPath := writeJSONLEvents(t,
		"{\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"id\":\"1792222215.125\",\"rule\":{\"description\":\"sensitive Wazuh SSH authentication\"},\"full_log\":\"private raw event\"}\n")
	ownerBefore, err := os.ReadFile(ownerPath)
	if err != nil { t.Fatal(err) }
	alertsBefore, err := os.ReadFile(alertsPath)
	if err != nil { t.Fatal(err) }
	var output bytes.Buffer
	err = runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+ownerPath, "-wazuh-alerts="+alertsPath,
		"-into-case=case-a", "-triage-import", "-with-case=case-c",
		"-compare-proximity", "-proximity-seconds=1",
	}, &output)
	if err != nil { t.Fatal(err) }
	var got localImportedProximity
	if err := json.Unmarshal(output.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Format != "cmdr-local-wazuh4-case-proximity-v1" ||
		!got.Synthetic || got.ImportedEntries != 1 ||
		got.Proximity.FirstCase != "case-a" || got.Proximity.SecondCase != "case-c" ||
		got.Proximity.WindowSeconds != 1 { t.Fatalf("incorrect Wazuh proximity receipt %+v", got) }
	found := false
	for _, lead := range got.Proximity.Leads {
		if strings.HasPrefix(lead.FirstEntryID, "local-jsonl-entry-") &&
			lead.SecondEntryID == "event-c" && lead.SeparationMS == 125 &&
			lead.FirstKind == "observed" && lead.SecondKind == "inferred" {
			found = true
		}
	}
	if !found { t.Fatalf("no Wazuh alert proximity lead in %+v", got.Proximity.Leads) }
	for _, secret := range []string{"sensitive Wazuh SSH", "private raw event",
		"demo-hidden-source", "demo-tenant-b", "Bearer ", "principals\""} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("sensitive data leaked: %q", secret)
		}
	}
	ownerAfter, err := os.ReadFile(ownerPath)
	if err != nil || !bytes.Equal(ownerBefore, ownerAfter) {
		t.Fatal("owner snapshot was modified")
	}
	alertsAfter, err := os.ReadFile(alertsPath)
	if err != nil || !bytes.Equal(alertsBefore, alertsAfter) {
		t.Fatal("source Wazuh alert file was modified")
	}
}

func TestWazuh4TemporalComparisonFailsClosed(t *testing.T) {
	owner := proximityDiskSnapshot(t)
	valid := writeJSONLEvents(t,
		"{\"timestamp\":\"2026-10-01T10:02:02.125Z\",\"id\":\"123.456\",\"rule\":{\"description\":\"restricted\"}}\n")
	invalid := writeJSONLEvents(t, "{\"@timestamp\":\"2026-10-01T10:02:02Z\",\"event\":{\"original\":\"unsupported Wazuh 5\"}}\n")
	for _, tc := range []struct{
		name, principal, other, alerts string
		seconds int
	}{
		{"foreign-second-Case", "analyst-a", "case-b", valid, 1},
		{"missing-second-Case", "analyst-a", "missing", valid, 1},
		{"same-Case", "analyst-a", "case-a", valid, 1},
		{"denied-principal", "analyst-b", "case-c", valid, 1},
		{"unsupported-source-format", "analyst-a", "case-c", invalid, 1},
		{"zero-window", "analyst-a", "case-c", valid, 0},
		{"overlarge-window", "analyst-a", "case-c", valid, 301},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			err := writeDemoWazuh4Proximity(context.Background(), tc.principal, owner,
				"case-a", tc.other, tc.alerts, tc.seconds, &output)
			if err == nil || output.Len() != 0 {
				t.Fatalf("unauthorized Wazuh analysis leaked %d bytes: %v", output.Len(), err)
			}
		})
	}
	for _, flags := range [][]string{
		{"-wazuh-alerts="+valid,"-snapshot="+owner,"-into-case=case-a","-triage-import","-with-case=case-a"},
		{"-wazuh-alerts="+valid,"-snapshot="+owner,"-into-case=case-a","-compare-proximity"},
		{"-wazuh-alerts="+valid,"-snapshot="+owner,"-into-case=case-a","-serve-import","-triage-import","-with-case=case-c","-compare-proximity"},
		{"-wazuh-alerts="+valid,"-snapshot="+owner,"-into-case=case-a","-triage-import","-with-case=case-c","-compare-proximity","-compare-events"},
		{"-wazuh-alerts="+valid,"-snapshot="+owner,"-into-case=case-a","-triage-import","-with-case=case-c","-compare-proximity","-import-jsonl="+valid},
	} {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), flags, &output); err == nil || output.Len()!=0 {
			t.Fatalf("ambiguous Wazuh options permitted or output leaked: %s", output.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := writeDemoWazuh4Proximity(ctx, "analyst-a", owner, "case-a",
		"case-c", valid, 1, &output); !errors.Is(err, context.Canceled) || output.Len()!=0 {
		t.Fatalf("cancelled analysis leaked output: %v", err)
	}
}
