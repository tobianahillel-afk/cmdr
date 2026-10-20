package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// This test uses real on-disk Wazuh-shaped alerts and the existing authorized
// Case -> Timeline -> typed-source navigation path, not a mocked technique count.
func TestWazuhMITRETechniquesAppearOnlyInAuthorizedCaseTriage(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"private-a\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"full_log\":\"PRIVATE-RAW-A\",\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"private-preview-a\",\"mitre\":{\"id\":[\"T1059\",\"T1059.001\",\"T1059.001\"]}}}\n"+
			"{\"id\":\"private-b\",\"timestamp\":\"2026-10-01T10:00:01Z\",\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"private-preview-b\",\"mitre\":{\"id\":[\"T1059.001\",\"T1078\"]}}}\n"+
			"{\"id\":\"private-c\",\"timestamp\":\"2026-10-01T10:00:02Z\",\"rule\":{\"id\":\"5500\",\"description\":\"private-preview-c\",\"mitre\":{\"id\":[\"T1078\"]}}}\n"+
			"{\"id\":\"private-d\",\"timestamp\":\"2026-10-01T10:00:03Z\",\"rule\":{\"id\":\"5700\",\"level\":0,\"description\":\"private-preview-d\"}}\n",
	)
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner,
		"-wazuh-alerts=" + alerts, "-into-case=case-a", "-triage-import",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ImportedEntries != 4 || got.DistinctTechniqueIDs != 3 ||
		got.UnmappedTechniqueAlerts != 1 || len(got.TopTechniques) != 3 {
		t.Fatalf("source MITRE metadata not faithfully aggregated: %+v", got)
	}
	techniques := got.TopTechniques
	if techniques[0].TechniqueID != "T1059.001" || techniques[0].Count != 2 ||
		techniques[0].HighestKnownRuleLevel == nil || *techniques[0].HighestKnownRuleLevel != 12 ||
		techniques[1].TechniqueID != "T1078" || techniques[1].Count != 2 ||
		techniques[1].HighestKnownRuleLevel == nil || *techniques[1].HighestKnownRuleLevel != 12 ||
		techniques[2].TechniqueID != "T1059" || techniques[2].Count != 1 {
		t.Fatalf("technique multiplicity, count or severity source wrong: %+v", techniques)
	}
	for _, sensitive := range []string{"private-a", "private-preview", "PRIVATE-RAW", "demo-hidden-source", "demo-tenant-b", "Bearer ", "sourceRef"} {
		if strings.Contains(out.String(), sensitive) {
			t.Fatalf("unauthorized or raw source data escaped aggregate: %q", sensitive)
		}
	}
	var denied bytes.Buffer
	if err := writeDemoWazuh4RuleTriage(context.Background(), "analyst-b", owner, "case-a", alerts, &denied);
		err == nil || denied.Len() != 0 {
		t.Fatalf("cross-owner MITRE cohort leaked bytes: %v / %q", err, denied.String())
	}
}

func TestWazuhMITRETechniqueRejectsMalformedOrExcessiveSourceLabels(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	bad := []string{
		"[\"T12\"]", "[\"TA0001\"]", "[\"T1059.12\"]",
		"[\"T1059.0000\"]", "[\"t1059\"]", "[\"T1059<script>\"]",
		"[\"T1059\",17]", "\"T1059\"",
		"[" + strings.Repeat("\"T1059\",", 16) + "\"T1078\"]",
	}
	for _, labels := range bad {
		t.Run(labels, func(t *testing.T) {
			alerts := writeJSONLEvents(t,
				"{\"id\":\"bad-mitre\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"rule\":{\"id\":\"5710\",\"level\":7,\"description\":\"private-preview\",\"mitre\":{\"id\":"+
					labels+"}}}\n",
			)
			var out bytes.Buffer
			err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", alerts, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid ATT&CK labels accepted or partial output: %v / %q", err, out.String())
			}
		})
	}
}

func TestWazuhMITRETechniqueTopTenDeterministicAcrossSourceOrder(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var source strings.Builder
	for i := 12; i >= 0; i-- {
		fmt.Fprintf(&source,
			"{\"id\":\"tech-%d\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"rule\":{\"id\":\"5710\",\"description\":\"seen\",\"mitre\":{\"id\":[\"T%04d\"]}}}\n",
			i, 1000+i)
	}
	alerts := writeJSONLEvents(t, source.String())
	var out bytes.Buffer
	if err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", alerts, &out); err != nil {
		t.Fatal(err)
	}
	var got wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.DistinctTechniqueIDs != 13 || got.UnmappedTechniqueAlerts != 0 || len(got.TopTechniques) != 10 {
		t.Fatalf("top technique bound incorrect: %+v", got)
	}
	for i, technique := range got.TopTechniques {
		if technique.TechniqueID != fmt.Sprintf("T%04d", 1000+i) ||
			technique.Count != 1 || technique.HighestKnownRuleLevel != nil {
			t.Fatalf("missing level invented or technique order unstable: %+v", got.TopTechniques)
		}
	}
}
