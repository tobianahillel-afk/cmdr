package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// One real operator Case is the SAME independently SHA-verified input to
// source inventory, the first-party mixed Event Search Job, and a two-Case
// comparison. Three separately authorized private reports are produced
// without exposing evidence in stdout, combining source owner graphs, or
// republishing the original Case snapshot.
func TestIntegratedSavedCaseProducesThreeConsistentPrivateInvestigationReports(t *testing.T) {
	file, selectedCase, digest := savedHeadlessMixedCase(t)
	_, _, _, otherFile, otherCase, otherDigest := buildTwoVerifiedLocalCaseArchives(t)
	dir := t.TempDir()
	paths := map[string]string{}
	workflows := []struct {
		name string
		args []string
	}{
		{"inventory", []string{
			"-principal", "analyst-a", "-snapshot", file,
			"-verify-import-snapshot-case", selectedCase,
			"-expect-snapshot-sha256", digest,
		}},
		{"mixed_search", []string{
			"-principal", "analyst-a", "-snapshot", file,
			"-hunt-saved-case", selectedCase, "-expect-snapshot-sha256", digest,
			"-hunt-field", "source_ref", "-hunt-family", "all",
			"-hunt-value", "local-jsonl-source-",
			"-hunt-from", "2026-10-17T09:59:00Z",
			"-hunt-to", "2026-10-17T10:02:00Z",
		}},
		{"comparison", []string{
			"-principal", "analyst-a", "-snapshot", file,
			"-compare-saved-case", selectedCase, "-expect-snapshot-sha256", digest,
			"-compare-saved-with-snapshot", otherFile,
			"-compare-saved-with-case", otherCase,
			"-expect-other-snapshot-sha256", otherDigest,
		}},
	}
	for _, workflow := range workflows {
		dest := filepath.Join(dir, workflow.name+".json")
		var stdout bytes.Buffer
		args := append(append([]string{}, workflow.args...), "-private-report", dest)
		if err := runDemoCLI(context.Background(), args, &stdout); err != nil {
			t.Fatalf("%s: independent owner authorization or report publication failed: %v", workflow.name, err)
		}
		if stdout.Len() != 0 {
			t.Fatalf("%s leaked investigation data to terminal", workflow.name)
		}
		paths[workflow.name] = dest
	}
	read := func(key string) map[string]any {
		t.Helper()
		raw, err := os.ReadFile(paths[key])
		if err != nil { t.Fatal(err) }
		for _, secret := range []string{
			"PRIVATE-WAZUH-ONE", "PRIVATE-WAZUH-TWO",
			"PRIVATE-GENERIC-NEAR", "Bearer token", "127.0.0.1:",
			"\"preview\"", "\"audit_trail\"",
		} {
			if bytes.Contains(raw, []byte(secret)) {
				t.Fatalf("%s report leaked source preview or bearer: %q", key, secret)
			}
		}
		var value map[string]any
		if err := json.Unmarshal(raw, &value); err != nil { t.Fatal(err) }
		return value
	}
	inventory := read("inventory")
	hunt := read("mixed_search")
	comparison := read("comparison")
	if inventory["format"] != "cmdr-offline-authorized-case-verification-v1" ||
		inventory["case_id"] != selectedCase ||
		inventory["authorized_wazuh_alerts"] != float64(2) ||
		inventory["authorized_generic_events"] != float64(1) {
		t.Fatalf("saved mixed Case source inventory inconsistent: %+v", inventory)
	}
	if hunt["format"] != "cmdr-local-headless-authorized-hunt-v1" ||
		hunt["case_id"] != selectedCase ||
		hunt["job_id"] == "" {
		t.Fatalf("private headless Search Job not associated with same saved Case: %+v", hunt)
	}
	refs, ok := hunt["authorized_refs"].([]any)
	if !ok || len(refs) != 3 {
		t.Fatalf("private headless hunt omitted authorized Wazuh/generic hits: %+v", hunt)
	}
	leads, ok := hunt["mixed_temporal_context"].([]any)
	if !ok || len(leads) != 2 {
		t.Fatalf("first-party authorized mixed temporal leads missing: %+v", hunt)
	}
	if comparison["format"] != "cmdr-local-verified-saved-case-comparison-v1" {
		t.Fatalf("private report did not reuse first-party comparison: %+v", comparison)
	}
	first, ok := comparison["first"].(map[string]any)
	if !ok || first["case_id"] != selectedCase {
		t.Fatalf("comparison disagrees with inventoried/searched Case: %+v", comparison)
	}
	second, ok := comparison["second"].(map[string]any)
	if !ok || second["case_id"] != otherCase {
		t.Fatalf("comparison did not independently select second Case: %+v", comparison)
	}
	// No recomputation or mutation is allowed: outputs must not modify
	// the original private archived Case or its operator-held digest.
	for _, archive := range []struct{path, expectedSHA string}{
		{file, digest}, {otherFile, otherDigest},
	} {
		data, err := os.ReadFile(archive.path)
		if err != nil { t.Fatal(err) }
		if fmt.Sprintf("%x", sha256.Sum256(data)) != archive.expectedSHA {
			t.Fatalf("archived Case was modified by one of the three report workflows: %q", archive.path)
		}
	}
	for _, dst := range paths {
		if dst == file || dst == otherFile {
			t.Fatal("private report path overwrote an archived Case")
		}
	}
}
