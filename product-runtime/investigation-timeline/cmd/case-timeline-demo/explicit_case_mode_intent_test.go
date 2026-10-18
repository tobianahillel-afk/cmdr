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
	flags := []string{
		"resume-import-search-case", "verify-import-snapshot-case",
		"export-case", "report-case", "compare-case", "with-case",
		"source-case", "source-ref", "inspect-case", "inspect-entry",
		"inspect-wazuh-entry", "hunt-saved-case", "compare-saved-case",
		"compare-saved-with-case", "compare-saved-with-snapshot",
		"expect-other-snapshot-sha256", "expect-snapshot-sha256",
		"save-import-snapshot", "save-wazuh-case-snapshot",
		"save-mixed-case-snapshot", "private-report",
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
