package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOfflineEntryInspectionUsesAuthorizedSourcePreview(t *testing.T) {
	var output bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=event-old",
	}, &output)
	if err != nil { t.Fatal(err) }
	var got localEntryInspection
	if err := json.Unmarshal(output.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.Format != "cmdr-local-entry-inspection-v1" || !got.Synthetic ||
		got.CaseID != "case-a" || got.EntryID != "event-old" ||
		got.Kind != "evidence" || got.ObjectID != "demo-evidence-a" ||
		!strings.Contains(got.Preview, "Pièce synthétique") {
		t.Fatalf("wrong authorized inspection: %+v", got)
	}
	for _, forbidden := range []string{"demo-hidden-source", "demo-secret", "demo-tenant-b", testBearer, "Bearer token"} {
		if strings.Contains(output.String(), forbidden) {
			t.Fatalf("inspected report leaked unrelated %q", forbidden)
		}
	}
}

func TestOfflineEntryInspectionReadsOperatorDiskSource(t *testing.T) {
	data := demoData()
	data.Sources[1].Preview = "operator-disk-only-authenticated-artifact-preview"
	payload, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "owner-snapshot.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + path, "-principal=analyst-a",
		"-inspect-case=case-a", "-inspect-entry=event-old",
	}, &out); err != nil { t.Fatal(err) }
	if !strings.Contains(out.String(), data.Sources[1].Preview) {
		t.Fatalf("offline inspection did not traverse disk-snapshot source authorization: %s", out.String())
	}
}

func TestOfflineEntryInspectionDenialsProduceNoOutput(t *testing.T) {
	cases := [][]string{
		{"-inspect-case=case-a"},
		{"-inspect-entry=event-old"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-source-case=case-a", "-source-ref=demo-evidence-a"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-report-case=case-a"},
		{"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=restricted"},
		{"-principal=analyst-a", "-inspect-case=case-a", "-inspect-entry=absent"},
		{"-principal=analyst-a", "-inspect-case=case-b", "-inspect-entry=event-b"},
		{"-principal=analyst-b", "-inspect-case=case-a", "-inspect-entry=event-old"},
		{"-principal=unknown", "-inspect-case=case-a", "-inspect-entry=event-old"},
		{"-inspect-case=case-a", "-inspect-entry=event-old", "-snapshot=relative.json"},
	}
	for _, args := range cases {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len()!=0 {
			t.Fatalf("invalid inspection produced output: args=%v, err=%v output=%s", args, err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background()); cancel()
	var out bytes.Buffer
	if err := writeDemoEntryInspection(ctx, "analyst-a", "", "case-a", "event-old", &out); !errors.Is(err, context.Canceled) || out.Len()!=0 {
		t.Fatalf("cancelled inspection disclosed data: %v %q", err, out.String())
	}
	if err := writeDemoEntryInspection(context.Background(), "analyst-a", "", "case-a", "event-old", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("writer error lost: %v", err)
	}
}
