package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCLIReopensOwnerSourceFromAuthorizedDiskSnapshot(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-open-source-case=case-c", "-open-source-entry=event-c",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got localSourceOwnerReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Format != "cmdr-local-authorized-source-preview-v1" ||
		!got.Synthetic || got.CaseID != "case-c" || got.EntryID != "event-c" ||
		got.Kind != "evidence" || got.ObjectID != "object-c" ||
		got.Preview != "SENSITIVE-PREVIEW-MUST-NOT-BE-IN-REPORT" {
		t.Fatalf("owner-backed source inspection incomplete: %+v", got)
	}
	for _, forbidden := range []string{
		"demo-hidden-source", "restricted-synthetic", "demo-tenant-b",
		"Bearer ", testBearer, "http://",
	} {
		if strings.Contains(out.String(), forbidden) {
			t.Fatalf("owner receipt contained unrelated metadata %q", forbidden)
		}
	}
}

func TestCLISourceOwnerInspectionRejectsUnauthorizedAndMixedModes(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	for _, args := range [][]string{
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-entry=event-c"},
		{"-principal=analyst-b", "-snapshot=" + path, "-open-source-case=case-c", "-open-source-entry=event-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-b", "-open-source-entry=event-b"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-a", "-open-source-entry=restricted"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-a", "-open-source-entry=event-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-c", "-open-source-entry=missing"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-c", "-open-source-entry=event-c", "-compare-events"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-c", "-open-source-entry=event-c", "-source-ref=disk-source-c"},
		{"-principal=analyst-a", "-snapshot=" + path, "-open-source-case=case-c", "-open-source-entry=event-c", "-report-case=case-a"},
		{"-principal=analyst-a", "-snapshot=relative.json", "-open-source-case=case-c", "-open-source-entry=event-c"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("unauthorized or mixed source inspection exposed data: err=%v out=%q", err, out.String())
		}
	}
}

func TestCLISourceOwnerInspectionCancellationAndOutputFailure(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoSourceOwnerReceipt(ctx, "analyst-a", path, "case-c", "event-c", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled source inspection exposed data: %v %q", err, out.String())
	}
	if err := writeDemoSourceOwnerReceipt(context.Background(), "analyst-a", path, "case-c", "event-c", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("source receipt ignored output error: %v", err)
	}
}

func TestCLISourceOwnerInspectionDoesNotTrustEntryFromAnotherCase(t *testing.T) {
	path := twoCaseComparisonSnapshot(t)
	for _, pair := range [][2]string{
		{"case-c", "event-new"},
		{"case-b", "event-c"},
		{"case-a", "event-c"},
	} {
		var out bytes.Buffer
		if err := writeDemoSourceOwnerReceipt(context.Background(), "analyst-a", path, pair[0], pair[1], &out); err == nil || out.Len() != 0 {
			t.Fatalf("source crossed Case ownership boundary %q: %q", pair, out.String())
		}
	}
}
