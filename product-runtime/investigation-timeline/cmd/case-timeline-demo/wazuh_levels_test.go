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

func TestWazuhRuleLevelsRequireAuthorizedImportedCase(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"alert-1\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":3,\"description\":\"PRIVATE-WAZUH-DESCRIPTION\"}}\n"+
			"{\"id\":\"alert-2\",\"timestamp\":\"2026-10-01T10:30:01Z\",\"rule\":{\"level\":12,\"description\":\"suspicious ssh\"}}\n"+
			"{\"id\":\"alert-3\",\"timestamp\":\"2026-10-01T10:30:02Z\",\"rule\":{\"level\":12,\"description\":\"other ssh\"}}\n"+
			"{\"id\":\"alert-4\",\"timestamp\":\"2026-10-01T10:30:03Z\",\"rule\":{\"description\":\"unknown level\"}}\n")
	beforeOwner, err := os.ReadFile(owner)
	if err != nil {
		t.Fatal(err)
	}
	beforeAlerts, err := os.ReadFile(alerts)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-principal=analyst-a", "-snapshot=" + owner,
		"-into-case=case-a", "-wazuh-alerts=" + alerts, "-triage-import"}
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), args, &output); err != nil {
		t.Fatal(err)
	}
	var result wazuhRuleTriageReceipt
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Format != "cmdr-local-wazuh4-case-triage-v1" ||
		!result.Synthetic || result.ImportedEntries != 4 ||
		result.Case.CaseID != "case-a" || result.Case.VisibleEntries != 211 ||
		result.UnknownRuleLevels != 1 || len(result.RuleLevels) != 2 ||
		result.RuleLevels[0] != (wazuhLevelCount{Level: 3, Count: 1}) ||
		result.RuleLevels[1] != (wazuhLevelCount{Level: 12, Count: 2}) {
		t.Fatalf("authorized Wazuh rule-level analysis incorrect: %+v", result)
	}
	for _, secret := range []string{"PRIVATE-WAZUH-DESCRIPTION", "demo-hidden-source",
		"demo-tenant-b", "alert-1", "Bearer ", "private raw event"} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("untrusted alert or unauthorized metadata leaked: %s", secret)
		}
	}
	for _, invalid := range [][]string{
		append(append([]string{}, args...), "-principal=analyst-b"),
		append(append([]string{}, args...), "-into-case=case-b"),
		append(append([]string{}, args...), "-with-case=case-b"),
	} {
		var denied bytes.Buffer
		if err := runDemoCLI(context.Background(), invalid, &denied); err == nil || denied.Len() != 0 {
			t.Fatalf("denied Wazuh triage wrote data: %v / %q", err, denied.String())
		}
	}
	ownerAfter, err := os.ReadFile(owner)
	if err != nil || !bytes.Equal(beforeOwner, ownerAfter) {
		t.Fatal("Wazuh Case triage modified owner snapshot")
	}
	alertsAfter, err := os.ReadFile(alerts)
	if err != nil || !bytes.Equal(beforeAlerts, alertsAfter) {
		t.Fatal("Wazuh Case triage modified source alerts")
	}
}

func TestWazuhRuleLevelBadSourceOrCancellationEmitsNothing(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	cases := []struct {
		name, raw string
	}{
		{"too-high", "{\"id\":\"hi\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":17,\"description\":\"alert\"}}\n"},
		{"negative", "{\"id\":\"lo\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":-1,\"description\":\"alert\"}}\n"},
		{"wrong-type", "{\"id\":\"str\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":\"high\",\"description\":\"alert\"}}\n"},
		{"duplicate", "{\"id\":\"twice\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":3,\"description\":\"alert\"}}\n{\"id\":\"twice\",\"timestamp\":\"2026-10-01T10:30:02Z\",\"rule\":{\"level\":5,\"description\":\"alert\"}}\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeJSONLEvents(t, tc.raw)
			var out bytes.Buffer
			err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", path, &out)
			if err == nil || out.Len() != 0 {
				t.Fatalf("invalid severity accepted or emitted data: %v %q", err, out.String())
			}
		})
	}
	path := writeJSONLEvents(t, "{\"id\":\"one\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":16,\"description\":\"alert\"}}\n")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoWazuh4Triage(ctx, "analyst-a", owner, "case-a", path, &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled Wazuh analysis emitted output: %v", err)
	}
}

func TestWazuhRuleLevelHistogramsIncludeBoundaryAndAbsentValues(t *testing.T) {
	owner := twoCaseComparisonSnapshot(t)
	alerts := writeJSONLEvents(t,
		"{\"id\":\"low-boundary\",\"timestamp\":\"2026-10-01T10:30:00Z\",\"rule\":{\"level\":0,\"description\":\"low\"}}\n"+
			"{\"id\":\"high-boundary\",\"timestamp\":\"2026-10-01T10:30:01Z\",\"rule\":{\"level\":16,\"description\":\"high\"}}\n"+
			"{\"id\":\"no-level\",\"timestamp\":\"2026-10-01T10:30:02Z\",\"rule\":{\"description\":\"not classified\"}}\n")
	var output bytes.Buffer
	if err := writeDemoWazuh4Triage(context.Background(), "analyst-a", owner, "case-a", alerts, &output); err != nil {
		t.Fatal(err)
	}
	var receipt wazuhRuleTriageReceipt
	if err := json.Unmarshal(output.Bytes(), &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.ImportedEntries != 3 || receipt.UnknownRuleLevels != 1 ||
		len(receipt.RuleLevels) != 2 ||
		receipt.RuleLevels[0] != (wazuhLevelCount{Level: 0, Count: 1}) ||
		receipt.RuleLevels[1] != (wazuhLevelCount{Level: 16, Count: 1}) {
		t.Fatalf("Wazuh boundary levels lost: %+v", receipt)
	}
}

func TestWazuhRuleIdentifiersAreBoundToAuthorizedCase(t *testing.T) {
	raw := []byte("{\"id\":\"same-alert\",\"rule\":{\"level\":12}}\n" +
		"{\"id\":\"missing\",\"rule\":{}}\n")
	first, err := wazuhRuleLevels(raw, "case-a")
	if err != nil { t.Fatal(err) }
	second, err := wazuhRuleLevels(raw, "case-c")
	if err != nil { t.Fatal(err) }
	if len(first) != 2 || len(second) != 2 {
		t.Fatalf("wrong Wazuh identity cardinality: %d %d", len(first), len(second))
	}
	var known, unknown int
	for id, level := range first {
		if _, collides := second[id]; collides {
			t.Fatal("Wazuh identity was shared across independent Cases")
		}
		if level.Known {
			if level.Level != 12 { t.Fatalf("wrong source level %d", level.Level) }
			known++
		} else {
			unknown++
		}
	}
	if known != 1 || unknown != 1 {
		t.Fatalf("known and unknown Wazuh rule levels not preserved: %d %d", known, unknown)
	}
}
