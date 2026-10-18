package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The source import is valid: before the fix the empty save flag silently
// selected the legacy stdout JSONL export and leaked the source message.
// Test both spellings produced by ordinary shell command expansion.
func TestExplicitEmptyPrivateCaseSaveNeverFallsBackToStdout(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot := snapshotFixture(t)
	generic := writeJSONLEvents(t,
		"{\"id\":\"private-event-1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"VERY-PRIVATE-FORENSIC-EVIDENCE-PREVIEW\"}\n")
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"private-wazuh-1\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-ALERT-PREVIEW\",\"level\":9}}\n")
	cases := []struct {
		name string
		args []string
		flag string
	}{
		{"generic normalized save", []string{
			"-principal", "analyst-a", "-snapshot", snapshot,
			"-into-case", caseID, "-import-jsonl", generic,
		}, "-save-import-snapshot"},
		{"Wazuh owner save", []string{
			"-principal", "analyst-a", "-snapshot", snapshot,
			"-into-case", caseID, "-wazuh-alerts", wazuh,
		}, "-save-wazuh-case-snapshot"},
		{"mixed source owner save", []string{
			"-principal", "analyst-a", "-snapshot", snapshot,
			"-into-case", caseID, "-wazuh-alerts", wazuh,
			"-import-jsonl", generic,
		}, "-save-mixed-case-snapshot"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, spelled := range []string{"-empty-arg", "-equals-notation"} {
				args := append([]string(nil), tc.args...)
				if spelled == "-empty-arg" {
					args = append(args, tc.flag, "")
				} else {
					args = append(args, tc.flag+"=")
				}
				var stdout bytes.Buffer
				err := runDemoCLI(context.Background(), args, &stdout)
				if err == nil || stdout.Len() != 0 {
					t.Fatalf("explicit empty private save downgraded to stdout: %v %q", err, stdout.String())
				}
				for _, secret := range []string{
					"VERY-PRIVATE-FORENSIC-EVIDENCE-PREVIEW",
					"PRIVATE-WAZUH-ALERT-PREVIEW",
					"Bearer", "127.0.0.1:",
				} {
					if strings.Contains(stdout.String(), secret) {
						t.Fatalf("private evidence appeared in output: %q", secret)
					}
				}
			}
		})
	}
}

func TestExplicitPrivateCaseSaveStillCreatesOnlyRequestedFile(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	snapshot := snapshotFixture(t)
	generic := writeJSONLEvents(t,
		"{\"id\":\"sensitive-normalized-1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"PRIVATE-SOURCE-MESSAGE-DO-NOT-PRINT\"}\n")
	out := filepath.Join(t.TempDir(), "authorized-private-owner.json")
	var stdout bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", snapshot, "-into-case", caseID,
		"-import-jsonl", generic, "-save-import-snapshot", out,
	}, &stdout)
	if err != nil { t.Fatalf("valid explicit private output rejected: %v", err) }
	if _, err := os.Stat(out); err != nil { t.Fatalf("saved authorized owner missing: %v", err) }
	if !strings.Contains(stdout.String(), "sha256") ||
		strings.Contains(stdout.String(), "PRIVATE-SOURCE-MESSAGE-DO-NOT-PRINT") {
		t.Fatalf("private destination receipt incorrect or leaking source preview: %q", stdout.String())
	}
	if _, err := loadDemoSnapshot(out); err != nil { t.Fatalf("saved owner invalid: %v", err) }
}
