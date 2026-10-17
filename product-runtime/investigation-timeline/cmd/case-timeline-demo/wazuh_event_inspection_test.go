package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

func importedWazuhID(caseID, rawID string, feed int) string {
	normalizedID := "wazuh-" + rawID
	if feed > 0 {
		normalizedID = fmt.Sprintf("feed-%d-%s", feed, normalizedID)
	}
	digest := sha256.Sum256([]byte(caseID + "\x00" + normalizedID))
	return fmt.Sprintf("local-jsonl-entry-%x", digest)
}

func TestWazuhMultiFeedDiskToAuthorizedEventInspectionAndPivot(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":9,\"description\":\"private-operator-first-alert\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5711\",\"level\":11,\"description\":\"private-operator-second-alert\"}}\n")
	for _, tc := range []struct {
		feed int
		agent, rule, preview string
	}{
		{1, "001", "5710", "private-operator-first-alert"},
		{2, "002", "5711", "private-operator-second-alert"},
	} {
		t.Run(tc.agent, func(t *testing.T) {
			entryID := importedWazuhID(caseID, "same.1", tc.feed)
			var output bytes.Buffer
			err := runDemoCLI(context.Background(), []string{
				"-snapshot=" + owner, "-principal=analyst-a", "-into-case=" + caseID,
				"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
				"-inspect-wazuh-entry=" + entryID,
				"-pivot-wazuh-field=wazuh.agent.id",
				"-pivot-start=2026-10-17T09:00:00Z",
				"-pivot-end=2026-10-17T11:00:00Z",
			}, &output)
			if err != nil {
				t.Fatal(err)
			}
			var got wazuhInspectionReceipt
			if err := json.Unmarshal(output.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if !got.Synthetic || got.CaseID != caseID ||
				got.Format != "cmdr-local-wazuh4-event-inspection-v1" ||
				got.Inspection == nil || got.Pivot == nil ||
				got.Inspection.EventRef != entryID ||
				got.Inspection.State != inspection.StateAvailable ||
				got.Inspection.Source.Kind != "telemetry-event" ||
				got.Inspection.NormalizedFields["wazuh.agent.id"] != tc.agent ||
				got.Inspection.NormalizedFields["wazuh.rule.id"] != tc.rule ||
				got.Pivot.Field != "wazuh.agent.id" ||
				got.Pivot.Value != tc.agent ||
				got.Pivot.EventRef != entryID ||
				got.Pivot.TargetCapability != inspection.PivotTargetCapability ||
				got.Pivot.ReturnContext.EventRef != entryID {
				t.Fatalf("first-party inspected Wazuh pivot incorrect: %+v", got)
			}
			if string(got.Inspection.RenderedPayload) != tc.preview ||
				len(got.Inspection.RawPayload) != 0 ||
				len(got.Inspection.RawDerivedSensitive) != 0 {
				t.Fatalf("wrong rendered versus raw access: %+v", got.Inspection)
			}
			other := "private-operator-first-alert"
			if tc.feed == 1 {
				other = "private-operator-second-alert"
			}
			if strings.Contains(output.String(), other) ||
				strings.Contains(output.String(), "demo-tenant-b") ||
				strings.Contains(output.String(), "Bearer ") {
				t.Fatalf("other feed/tenant or bearer leaked: %s", output.String())
			}
		})
	}
}

func TestWazuhSingleFeedInspectionWithoutPivot(t *testing.T) {
	owner := snapshotFixture(t)
	path := writeJSONLEvents(t, "{\"id\":\"one.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"000\"},\"rule\":{\"id\":\"5710\",\"description\":\"offline-one-feed\"}}\n")
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + owner, "-principal=analyst-a",
		"-into-case=forensic-case-on-disk", "-wazuh-alerts=" + path,
		"-inspect-wazuh-entry=" + importedWazuhID("forensic-case-on-disk", "one.1", 0),
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var got wazuhInspectionReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Pivot != nil || got.Inspection == nil || got.Inspection.NormalizedFields["wazuh.agent.id"] != "000" ||
		string(got.Inspection.RenderedPayload) != "offline-one-feed" {
		t.Fatalf("single-feed Event Inspection incorrect: %+v", got)
	}
}

func TestWazuhEventInspectionRejectsUnauthorizedAndAmbiguousSelections(t *testing.T) {
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"protected-feed-one\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5711\",\"description\":\"protected-feed-two\"}}\n")
	good := []string{
		"-snapshot=" + owner, "-principal=analyst-a", "-into-case=forensic-case-on-disk",
		"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
		"-inspect-wazuh-entry=" + importedWazuhID("forensic-case-on-disk", "same.1", 2),
		"-pivot-wazuh-field=wazuh.rule.id",
		"-pivot-start=2026-10-17T09:00:00Z",
		"-pivot-end=2026-10-17T11:00:00Z",
	}
	for _, tc := range []struct {
		name string
		replace int
		value string
		extra []string
	}{
		{"foreign principal", 1, "-principal=analyst-b", nil},
		{"foreign Case", 2, "-into-case=case-b", nil},
		{"unimported entry", 5, "-inspect-wazuh-entry=event-old", nil},
		{"bad interval", 7, "-pivot-start=2026-10-17T12:00:00Z", nil},
		{"unsupported pivot", 6, "-pivot-wazuh-field=hostname", nil},
		{"other mode", -1, "", []string{"-serve-import"}},
		{"other comparison", -1, "", []string{"-with-case=case-b"}},
		{"missing end", 8, "-pivot-end=", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string(nil), good...)
			if tc.replace >= 0 {
				args[tc.replace] = tc.value
			}
			args = append(args, tc.extra...)
			var output bytes.Buffer
			if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() > 0 {
				t.Fatalf("invalid selection leaked inspection: %v %s", err, output.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := runDemoCLI(ctx, good, &output); err == nil || output.Len() != 0 {
		t.Fatalf("cancelled Wazuh inspection emitted data: %v %s", err, output.String())
	}
}

func TestWazuhInspectionCannotInventMissingSearchPivot(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	feed := writeJSONLEvents(t, "{\"id\":\"no-rule.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"000\"},\"rule\":{\"description\":\"missing rule ID\"}}\n")
	args := []string{
		"-snapshot=" + owner, "-principal=analyst-a", "-into-case=" + caseID,
		"-wazuh-alerts=" + feed, "-inspect-wazuh-entry=" + importedWazuhID(caseID, "no-rule.1", 0),
		"-pivot-wazuh-field=wazuh.rule.id",
		"-pivot-start=2026-10-17T09:00:00Z",
		"-pivot-end=2026-10-17T11:00:00Z",
	}
	var denied bytes.Buffer
	if err := runDemoCLI(context.Background(), args, &denied); err == nil || denied.Len() != 0 {
		t.Fatalf("missing rule identity produced fabricated search pivot: %v %s", err, denied.String())
	}
	args[5] = "-pivot-wazuh-field=wazuh.agent.id"
	var allowed bytes.Buffer
	if err := runDemoCLI(context.Background(), args, &allowed); err != nil {
		t.Fatal(err)
	}
	var result wazuhInspectionReceipt
	if err := json.Unmarshal(allowed.Bytes(), &result); err != nil ||
		result.Pivot == nil || result.Pivot.Field != "wazuh.agent.id" ||
		result.Pivot.Value != "000" {
		t.Fatalf("valid agent-only pivot lost: %v %+v", err, result)
	}
	// Operator's explicit sensitive read is unaffected by missing rule IDs.
	args = args[:5]
	var inspected bytes.Buffer
	if err := runDemoCLI(context.Background(), args, &inspected); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(inspected.Bytes(), &result); err != nil || result.Inspection == nil {
		t.Fatalf("source inspection without search pivot failed: %v %+v", err, result)
	}
}

func TestWazuhInspectionRejectsInvalidBoundaryBeforeOpeningSources(t *testing.T) {
	for _, tc := range []struct {
		name, principal, caseID, entryID, pivotField, from, to string
		paths []string
	}{
		{name:"missing principal", caseID:"case-a", entryID:"event-a", paths:[]string{"/nonexistent"}},
		{name:"missing Case", principal:"analyst-a", entryID:"event-a", paths:[]string{"/nonexistent"}},
		{name:"missing entry", principal:"analyst-a", caseID:"case-a", paths:[]string{"/nonexistent"}},
		{name:"missing feeds", principal:"analyst-a", caseID:"case-a", entryID:"event-a"},
		{name:"too many feeds", principal:"analyst-a", caseID:"case-a", entryID:"event-a", paths:[]string{"a","b","c","d","e"}},
		{name:"bounds without pivot", principal:"analyst-a", caseID:"case-a", entryID:"event-a", paths:[]string{"a"}, from:"2026-10-17T10:00:00Z"},
		{name:"pivot without times", principal:"analyst-a", caseID:"case-a", entryID:"event-a", paths:[]string{"a"}, pivotField:"wazuh.rule.id"},
		{name:"unapproved pivot field", principal:"analyst-a", caseID:"case-a", entryID:"event-a", paths:[]string{"a"}, pivotField:"raw.full_log", from:"2026-10-17T09:00:00Z", to:"2026-10-17T11:00:00Z"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := writeDemoWazuhEventInspection(context.Background(), tc.principal, "", tc.caseID,
				tc.paths, tc.entryID, tc.pivotField, tc.from, tc.to, &output); err == nil ||
				output.Len() != 0 {
				t.Fatalf("invalid source selection passed: %v %s", err, output.String())
			}
		})
	}
}

func TestWazuhInspectionRejectsMalformedSourceMetadataAndWriterFailure(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	malformed := writeJSONLEvents(t, "{\"id\":\"malformed.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"not-numeric\",\"description\":\"invalid source rule\"}}\n")
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + owner, "-principal=analyst-a", "-into-case=" + caseID,
		"-wazuh-alerts=" + malformed,
		"-inspect-wazuh-entry=" + importedWazuhID(caseID, "malformed.1", 0),
	}, &output); err == nil || output.Len() != 0 {
		t.Fatalf("unvalidated source rule metadata was exposed: %v %s", err, output.String())
	}
	valid := writeJSONLEvents(t, "{\"id\":\"valid.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"description\":\"valid source\"}}\n")
	if err := writeDemoWazuhEventInspection(context.Background(), "analyst-a", owner, caseID,
		[]string{valid}, importedWazuhID(caseID, "valid.1", 0), "", "", "",
		rejectReportOutput{}); err == nil {
		t.Fatal("failed sensitive receipt writer was silently accepted")
	}
}
