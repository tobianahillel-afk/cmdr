package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestWazuhRuleIDCohortsAreVisibleOnlyAfterCaseAuthorization(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"secret-alert-A\",\"timestamp\":\"2026-10-01T10:00:00.500Z\",\"rule\":{\"id\":\"5710\",\"level\":4,\"description\":\"PRIVATE-PREVIEW-A\"}}\n"+
			"{\"id\":\"secret-alert-B\",\"timestamp\":\"2026-10-01T10:00:00Z\",\"rule\":{\"id\":\"5710\",\"level\":12,\"description\":\"PRIVATE-PREVIEW-B\"}}\n"+
			"{\"id\":\"secret-alert-C\",\"timestamp\":\"2026-10-01T10:05:00Z\",\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-PREVIEW-C\"}}\n"+
			"{\"id\":\"secret-alert-D\",\"timestamp\":\"2026-10-01T10:04:00+01:00\",\"rule\":{\"id\":\"5500\",\"level\":7,\"description\":\"PRIVATE-PREVIEW-D\"}}\n"+
			"{\"id\":\"secret-alert-E\",\"timestamp\":\"2026-10-01T10:03:00Z\",\"rule\":{\"level\":0,\"description\":\"PRIVATE-PREVIEW-E\"}}\n")
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + owner,
		"-wazuh-alerts=" + alerts, "-into-case=case-a", "-triage-import",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var receipt wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.ImportedEntries != 5 || receipt.Case.VisibleEntries != 212 ||
		receipt.DistinctRuleIDs != 2 || receipt.UnknownRuleIDs != 1 ||
		len(receipt.TopRules) != 2 {
		t.Fatalf("authorized Wazuh cohorts incomplete: %+v", receipt)
	}
	first, second := receipt.TopRules[0], receipt.TopRules[1]
	if first.RuleID != "5710" || first.Count != 3 ||
		first.HighestKnownRuleLevel == nil || *first.HighestKnownRuleLevel != 12 ||
		first.FirstCertainUTC != "2026-10-01T10:00:00.000000000Z" ||
		first.LastCertainUTC != "2026-10-01T10:05:00.000000000Z" ||
		second.RuleID != "5500" || second.Count != 1 ||
		second.HighestKnownRuleLevel == nil || *second.HighestKnownRuleLevel != 7 ||
		second.FirstCertainUTC != "2026-10-01T09:04:00.000000000Z" {
		t.Fatalf("wrong rule cohort counts/time-ordering: %+v", receipt.TopRules)
	}
	for _, protected := range []string{
		"secret-alert", "PRIVATE-PREVIEW", "demo-hidden-source",
		"demo-tenant-b", "Bearer ", "full_log",
	} {
		if strings.Contains(out.String(), protected) {
			t.Fatalf("Wazuh raw/unauthorized data leaked: %q", protected)
		}
	}
	var denied bytes.Buffer
	if err := writeDemoWazuh4RuleTriage(context.Background(), "analyst-b", owner, "case-a", alerts, &denied);
		err == nil || denied.Len() != 0 {
		t.Fatalf("unauthorized rule cohort was exposed: %v / %q", err, denied.String())
	}
}

func TestWazuhRuleIDRequiresBoundedDecimalStringBeforeOutput(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	for _, candidate := range []string{
		"\"57x0\"", "5710", "\"12345678901234567\"", "\"57 10\"",
	} {
		t.Run(candidate, func(t *testing.T) {
			alerts := writeJSONLEvents(t,
				"{\"id\":\"source-alert\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"id\":"+
					candidate+",\"level\":9,\"description\":\"private-alert\"}}\n")
			var out bytes.Buffer
			err := writeDemoWazuh4RuleTriage(context.Background(), "analyst-a", owner, "case-a", alerts, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("malformed rule ID accepted or leaked output: %v / %q", err, out.String())
			}
		})
	}
}

func TestWazuhRuleIDCohortTopTenIsDeterministic(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	var lines strings.Builder
	for i := 12; i >= 0; i-- {
		// Deliberately reverse input IDs; all rule frequencies are tied.
		fmt.Fprintf(&lines,
			"{\"id\":\"alert-%d\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"id\":\"%04d\",\"level\":0,\"description\":\"visible\"}}\n", i, 1000+i)
	}
	path := writeJSONLEvents(t, lines.String())
	var out bytes.Buffer
	if err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", path, &out); err != nil {
		t.Fatal(err)
	}
	var receipt wazuhRuleTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.DistinctRuleIDs != 13 || len(receipt.TopRules) != 10 || receipt.UnknownRuleIDs != 0 {
		t.Fatalf("missing rule cohort totals: %+v", receipt)
	}
	for i, rule := range receipt.TopRules {
		if rule.RuleID != fmt.Sprintf("%04d", 1000+i) || rule.Count != 1 ||
			rule.HighestKnownRuleLevel == nil || *rule.HighestKnownRuleLevel != 0 {
			t.Fatalf("rule rank or numeric level zero mishandled: %+v", rule)
		}
	}
}
