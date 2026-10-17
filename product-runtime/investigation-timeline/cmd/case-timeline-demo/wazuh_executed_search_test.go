package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestExecutedWazuhPivotSearchRunsActualJobOnAuthorizedDiskFeeds(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t,
		"{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"first-private-alert\"}}\n"+
		"{\"id\":\"late.1\",\"timestamp\":\"2026-10-17T10:05:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"late-private-alert\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"description\":\"second-private-alert\"}}\n")
	original := map[string][]byte{}
	for _, p := range []string{owner, first, second} {
		contents, err := os.ReadFile(p)
		if err != nil { t.Fatal(err) }
		original[p] = contents
	}
	start := []string{
		"-snapshot=" + owner, "-principal=analyst-a", "-into-case=" + caseID,
		"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
		"-inspect-wazuh-entry=" + importedWazuhID(caseID, "same.1", 1),
		"-pivot-wazuh-field=wazuh.rule.id",
		"-pivot-start=2026-10-17T09:59:00Z",
		"-pivot-end=2026-10-17T10:03:00Z",
		"-execute-wazuh-pivot",
	}
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), start, &out); err != nil { t.Fatal(err) }
	var receipt wazuhInspectionReceipt
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil { t.Fatal(err) }
	if receipt.Pivot == nil || receipt.SearchJob == nil || receipt.SearchJob.State != search.JobCompleted ||
		receipt.SearchJob.Version < 2 || len(receipt.SearchJob.ResultRefs) != 2 ||
		receipt.SearchJob.Provenance.TenantRef != receipt.Pivot.TenantRef ||
		receipt.SearchJob.Provenance.QueryVersion != "wazuh-agent-rule-pivot-v1" ||
		len(receipt.SearchJob.AuditTrail) < 3 {
		t.Fatalf("missing real first-party completed and audited Job: %+v", receipt)
	}
	seen := map[string]bool{}
	for _, ref := range receipt.SearchJob.ResultRefs {
		if ref.Source != wazuhLocalSearchSource {
			t.Fatalf("unexpected source authority: %+v", ref)
		}
		seen[ref.Ref] = true
	}
	if !seen[importedWazuhID(caseID, "same.1", 1)] ||
		!seen[importedWazuhID(caseID, "same.1", 2)] ||
		seen[importedWazuhID(caseID, "late.1", 1)] {
		t.Fatalf("half-open temporal SearchJob matches incorrect: %+v", receipt.SearchJob.ResultRefs)
	}
	for _, forbidden := range []string{"second-private-alert", "late-private-alert", "demo-tenant-b", "Bearer "} {
		if strings.Contains(out.String(), forbidden) {
			t.Fatalf("SearchJob leaked another event's sensitive content %q", forbidden)
		}
	}
	for path, before := range original {
		after, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatalf("search changed operator-owned evidence at %s: %v", path, err)
		}
	}
	// The selected normalized agent ID executes a narrower valid query;
	// the existing draft-only API still never executes a backend implicitly.
	start[6] = "-pivot-wazuh-field=wazuh.agent.id"
	out.Reset()
	if err := runDemoCLI(context.Background(), start, &out); err != nil { t.Fatal(err) }
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil ||
		receipt.SearchJob == nil || len(receipt.SearchJob.ResultRefs) != 1 {
		t.Fatalf("agent pivot search should return one owner-authorized hit: %v %+v", err, receipt)
	}
	start = start[:len(start)-1]
	out.Reset()
	if err := runDemoCLI(context.Background(), start, &out); err != nil { t.Fatal(err) }
	receipt = wazuhInspectionReceipt{}
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil ||
		receipt.SearchJob != nil || receipt.Pivot == nil {
		t.Fatalf("prior draft-only behavior regressed: %v %+v", err, receipt)
	}
}

func TestExecutedWazuhPivotSearchFailsClosedOnScopeAndRevocation(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"secret-first\"}}\n")
	second := writeJSONLEvents(t, "{\"id\":\"same.1\",\"timestamp\":\"2026-10-17T10:01:00Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"5710\",\"description\":\"secret-second\"}}\n")
	good := []string{
		"-snapshot=" + owner, "-principal=analyst-a", "-into-case=" + caseID,
		"-wazuh-alerts=" + first, "-wazuh-alerts=" + second,
		"-inspect-wazuh-entry=" + importedWazuhID(caseID, "same.1", 1),
		"-pivot-wazuh-field=wazuh.rule.id",
		"-pivot-start=2026-10-17T09:59:00Z",
		"-pivot-end=2026-10-17T10:03:00Z",
		"-execute-wazuh-pivot",
	}
	for _, tc := range []struct {
		name string
		slot int
		replacement string
		extra string
	}{
		{"foreign principal", 1, "-principal=analyst-b", ""},
		{"foreign Case", 2, "-into-case=case-b", ""},
		{"entry from unrelated Case", 5, "-inspect-wazuh-entry=event-old", ""},
		{"unsupported normalized pivot", 6, "-pivot-wazuh-field=full_log", ""},
		{"inverted interval", 7, "-pivot-start=2026-10-17T12:00:00Z", ""},
		{"no explicit pivot", 6, "-pivot-wazuh-field=", ""},
		{"conflicting listener", -1, "", "-serve-import"},
		{"conflicting Case comparison", -1, "", "-with-case=case-b"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string(nil), good...)
			if tc.slot >= 0 { args[tc.slot] = tc.replacement }
			if tc.extra != "" { args = append(args, tc.extra) }
			var out bytes.Buffer
			if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
				t.Fatalf("unauthorized executed job disclosed output: %v %s", err, out.String())
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var cancelled bytes.Buffer
	if err := runDemoCLI(ctx, good, &cancelled); err == nil || cancelled.Len() != 0 {
		t.Fatalf("cancelled job wrote data: %v %s", err, cancelled.String())
	}

	data, count, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(), "analyst-a", owner, caseID, []string{first, second})
	if err != nil { t.Fatal(err) }
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	handler, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	for _, tc := range []struct {
		name, token, tenant string
		imported int
	}{
		{"revoked owner bearer", "invalid-token", "demo-tenant-a", count},
		{"foreign tenant", token, "demo-tenant-b", count},
		{"incomplete owner import", token, "demo-tenant-a", count + 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			backend := &wazuhPivotBackend{
				owner:handler, token:tc.token, caseID:caseID, tenant:tc.tenant,
				environment:"lab", levels:levels, imported:tc.imported,
				field:"wazuh.rule.id", value:"5710",
			}
			var request search.BackendRequest
			request.TenantRef = "demo-tenant-a"
			request.EnvironmentRef = "lab"
			request.TimeStart = time.Date(2026,10,17,9,0,0,0,time.UTC)
			request.TimeEnd = time.Date(2026,10,17,12,0,0,0,time.UTC)
			request.Sources = []string{wazuhLocalSearchSource}
			result, err := backend.Execute(context.Background(), request)
			if err == nil || len(result.ResultRefs) != 0 {
				t.Fatalf("revoked owner or tenant admitted results: %+v %v", result, err)
			}
		})
	}
}
