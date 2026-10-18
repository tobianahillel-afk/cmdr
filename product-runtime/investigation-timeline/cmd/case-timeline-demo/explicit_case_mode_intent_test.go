package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// A shell-expanded empty Case ID must never silently turn an explicitly
// selected confidential offline operation into the default loopback viewer
// (which would print a bearer) or the legacy stdout evidence exporter.
func TestExplicitEmptyCaseModesNeverFallbackToListenerOrEvidenceOutput(t *testing.T) {
	snapshot := snapshotFixture(t)
	// Representative sensitive fall-through routes; all selectors are
	// protected by one shared entry check, not separate validators.
	flags := []string{
		"resume-import-search-case", "export-case", "report-case",
		"source-case", "inspect-case", "compare-case", "with-case",
		"expect-other-snapshot-sha256",
	}
	for _, selector := range flags {
		for _, syntax := range []string{"separate", "equals"} {
			t.Run(selector+"/"+syntax, func(t *testing.T) {
				arguments := []string{"-principal", "analyst-a", "-snapshot", snapshot}
				if syntax == "equals" {
					arguments = append(arguments, "-"+selector+"=")
				} else {
					arguments = append(arguments, "-"+selector, "")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
				defer cancel()
				var output bytes.Buffer
				err := runDemoCLI(ctx, arguments, &output)
				if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
					t.Fatalf("empty %s mode should fail at the CLI boundary, not launch a viewer: %v", selector, err)
				}
				if output.Len() != 0 {
					t.Fatalf("explicit-empty selector exposed a bearer or evidence in stdout: %q", output.String())
				}
			})
		}
	}
}

// Positive CLI behavior is intentionally unchanged. An actually omitted
// selector still chooses the historical default synthetic loopback viewer,
// and nonempty direct Case reporting remains a first-party real workflow.
func TestNonemptyCaseReportRemainsAvailableAfterSelectorHardening(t *testing.T) {
	snapshot := snapshotFixture(t)
	var report bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", snapshot,
		"-report-case", "forensic-case-on-disk",
	}, &report)
	if err != nil { t.Fatalf("nonempty Case report was blocked: %v", err) }
	if !strings.Contains(report.String(), "forensic-case-on-disk") ||
		strings.Contains(report.String(), "Bearer token") {
		t.Fatalf("authorized offline Case report missing or viewer leaked a token: %q", report.String())
	}
}

func TestExplicitEmptySnapshotNeverSelectsDefaultFixtureForOfflineReport(t *testing.T) {
	for _,args := range [][]string{
		{"-snapshot=", "-report-case", "forensic-case-on-disk"},
		{"-snapshot", "", "-report-case", "forensic-case-on-disk"},
	} {
		var output bytes.Buffer
		err:=runDemoCLI(context.Background(),args,&output)
		if err==nil || !strings.Contains(err.Error(),"cannot be empty") || output.Len()!=0 {
			t.Fatalf("explicitly empty Case owner snapshot fell back to synthetic default: %v %q",err,output.String())
		}
	}
}
