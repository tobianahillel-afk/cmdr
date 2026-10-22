package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestWazuhAgentAndRuleCorrelationFromAuthorizedDiskImport(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"a1\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"PRIVATE-AGENT-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"PRIVATE-EVENT\"}}\n"+
			"{\"id\":\"a2\",\"timestamp\":\"2026-10-01T10:00:01Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-EVENT\"}}\n"+
			"{\"id\":\"a3\",\"timestamp\":\"2026-10-01T10:00:02Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5500\",\"description\":\"PRIVATE-EVENT\"}}\n"+
			"{\"id\":\"a4\",\"timestamp\":\"2026-10-01T10:00:03Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"level\":3,\"description\":\"PRIVATE-EVENT\"}}\n"+
			"{\"id\":\"a5\",\"timestamp\":\"2026-10-01T10:00:04Z\",\"rule\":{\"id\":\"5700\",\"level\":0,\"description\":\"PRIVATE-EVENT\"}}\n")
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+owner, "-wazuh-alerts="+alerts,
		"-into-case=case-a", "-triage-import",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ImportedEntries != 5 || got.Case.VisibleEntries != 212 ||
		got.DistinctAgentIDs != 2 || got.UnknownAgentIDs != 1 || len(got.TopAgents) != 2 {
		t.Fatalf("agent cohort missing or miscounted: %+v", got)
	}
	a, b := got.TopAgents[0], got.TopAgents[1]
	if a.AgentID != "001" || a.Count != 3 || a.DistinctRuleIDs != 2 ||
		a.HighestKnownRuleLevel == nil || *a.HighestKnownRuleLevel != 12 ||
		b.AgentID != "002" || b.Count != 1 || b.DistinctRuleIDs != 1 ||
		b.HighestKnownRuleLevel == nil || *b.HighestKnownRuleLevel != 3 {
		t.Fatalf("wrong agent ID or rule correlation: %+v", got.TopAgents)
	}
	for _, secret := range []string{"PRIVATE-AGENT-HOST", "PRIVATE-EVENT", "demo-hidden-source", "demo-tenant-b", "Bearer ", "\"a1\""} {
		if strings.Contains(out.String(), secret) {
			t.Fatalf("source or unauthorized data leaked: %q", secret)
		}
	}
	var denied bytes.Buffer
	if err := writeDemoWazuh4RuleTriage(context.Background(), "analyst-b", owner, "case-a", alerts, &denied);
		err == nil || denied.Len() != 0 {
		t.Fatalf("denied agent grouping emitted data: %v / %q", err, denied.String())
	}
}

func TestWazuhAgentIDRejectsMalformedMetadataBeforeReport(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	for _, id := range []string{`"00x"`, `123`, `"000 1"`, `"12345678901234567"`} {
		t.Run(id, func(t *testing.T) {
			alerts := writeJSONLEvents(t,
				"{\"id\":\"bad-agent\",\"timestamp\":\"2026-10-01T10:01:00Z\",\"agent\":{\"id\":"+id+"},\"rule\":{\"id\":\"5710\",\"description\":\"private\"}}\n")
			var out bytes.Buffer
			err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", alerts, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid Wazuh agent identity accepted or output leaked: %v / %q", err, out.String())
			}
		})
	}
}

func TestWazuhAgentTopTenDeterministicWhenFrequenciesTie(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var lines strings.Builder
	for i := 12; i >= 0; i-- {
		fmt.Fprintf(&lines, "{\"id\":\"agent-alert-%d\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"agent\":{\"id\":\"%03d\"},\"rule\":{\"id\":\"5710\",\"level\":0,\"description\":\"synthetic\"}}\n", i, i)
	}
	alerts := writeJSONLEvents(t, lines.String())
	var out bytes.Buffer
	if err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", alerts, &out); err != nil {
		t.Fatal(err)
	}
	var got wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.DistinctAgentIDs != 13 || len(got.TopAgents) != 10 || got.UnknownAgentIDs != 0 {
		t.Fatalf("agent truncation or counts incorrect: %+v", got)
	}
	for i, agent := range got.TopAgents {
		if agent.AgentID != fmt.Sprintf("%03d", i) || agent.Count != 1 || agent.DistinctRuleIDs != 1 ||
			agent.HighestKnownRuleLevel == nil || *agent.HighestKnownRuleLevel != 0 {
			t.Fatalf("agent ordering or severity zero lost: %+v", agent)
		}
	}
}
