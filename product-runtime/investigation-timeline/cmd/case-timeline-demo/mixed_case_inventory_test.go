package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestOfflineInventoryDistinguishesAuthorizedMixedWazuhAndGenericSources(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	wazuh := writeJSONLEvents(t,
		"{\"id\":\"a\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-ONE\",\"level\":7}}\n"+
		"{\"id\":\"b\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"PRIVATE-WAZUH-TWO\",\"level\":8}}\n")
	generic := writeJSONLEvents(t,
		"{\"id\":\"generic-id\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"PRIVATE-GENERIC-MESSAGE\"}\n")
	file := filepath.Join(t.TempDir(), "authorized-case.json")
	if err := writePersistedMixedCase(context.Background(), "analyst-a", snapshotFixture(t), caseID,
		[]string{wazuh}, []string{generic}, file, &bytes.Buffer{}); err != nil {
		t.Fatalf("real one-step Wazuh + generic save failed: %v", err)
	}
	raw, err := os.ReadFile(file)
	if err != nil { t.Fatal(err) }
	hash := fmt.Sprintf("%x", sha256.Sum256(raw))
	check := func(path, expectedHash string) (struct {
		Format string `json:"format"`
		Count int `json:"authorized_imported_events"`
		Wazuh int `json:"authorized_wazuh_alerts"`
		Generic int `json:"authorized_generic_events"`
		Slots map[string]int `json:"authorized_events_by_import_slot"`
	}, string, error) {
		var decoded struct {
			Format string `json:"format"`
			Count int `json:"authorized_imported_events"`
			Wazuh int `json:"authorized_wazuh_alerts"`
			Generic int `json:"authorized_generic_events"`
			Slots map[string]int `json:"authorized_events_by_import_slot"`
		}
		var output bytes.Buffer
		err := runDemoCLI(context.Background(), []string{
			"-principal", "analyst-a", "-snapshot", path,
			"-verify-import-snapshot-case", caseID, "-expect-snapshot-sha256", expectedHash,
		}, &output)
		if err == nil { err = json.Unmarshal(output.Bytes(), &decoded) }
		return decoded, output.String(), err
	}
	receipt, text, err := check(file, hash)
	if err != nil || receipt.Format != "cmdr-offline-authorized-case-verification-v1" ||
		receipt.Count != 3 || receipt.Wazuh != 2 || receipt.Generic != 1 ||
		receipt.Slots["operator-local-jsonl"] != 2 || receipt.Slots["operator-local-jsonl-feed-1"] != 1 {
		t.Fatalf("mixed Case source inventory wrongly merges source families: %+v %v", receipt, err)
	}
	for _,secret := range []string{"PRIVATE-WAZUH", "PRIVATE-GENERIC", "local-jsonl-entry-",
		"Bearer ", "rule_id", "agent_id", "127.0.0.1:"} {
		if strings.Contains(text, secret) { t.Fatalf("source inventory leaked sensitive metadata: %q", secret) }
	}

	// Revocation in a valid legacy operator snapshot must be honored
	// independently for each typed source, BEFORE counting mixed inputs.
	data, err := loadDemoSnapshotWithExpectedDigest(file, hash)
	if err != nil { t.Fatal(err) }
	var principal timeline.MemoryTimelinePrincipal
	for _, p := range data.Principals {
		if p.ID == "analyst-a" { principal = p }
	}
	if principal.ID == "" { t.Fatal("principal absent from saved owner") }
	data.Principals = append(data.Principals, timeline.MemoryTimelinePrincipal{
		ID: "other-reader", Tenant: principal.Tenant, Environment: principal.Environment,
	})
	revoked := false
	for i := range data.Sources {
		if data.Sources[i].Wazuh != nil && !revoked {
			data.Sources[i].Readers = []string{"other-reader"}
			revoked = true
		}
	}
	if !revoked { t.Fatal("Wazuh source metadata missing") }
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), "analyst-a", data); err != nil {
		t.Fatalf("source-revocation input graph must remain valid: %v", err)
	}
	changed, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	legacy := filepath.Join(t.TempDir(), "old-owner-with-revocation.json")
	if err := os.WriteFile(legacy, changed, 0600); err != nil { t.Fatal(err) }
	newHash := fmt.Sprintf("%x", sha256.Sum256(changed))
	filtered, privateText, err := check(legacy, newHash)
	if err != nil || filtered.Count != 2 || filtered.Wazuh != 1 || filtered.Generic != 1 {
		t.Fatalf("denied Wazuh source incorrectly counted in authorized inventory: %+v %v", filtered, err)
	}
	for _,secret := range []string{"PRIVATE-WAZUH", "PRIVATE-GENERIC", "local-jsonl-source-", "other-reader"} {
		if strings.Contains(privateText, secret) { t.Fatalf("source-denied data disclosed in verification receipt: %q", secret) }
	}
}
