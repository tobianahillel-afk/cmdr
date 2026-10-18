package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Demonstrates ONE investigation journey only possible when the frozen
// reverse-Entry and rule/agent cohort branches are truly integrated:
// exact authorized Wazuh cohort -> second Entry ID -> all authorized first
// counterparts -> independently source-authorized private paired dossier.
func TestIntegratedCohortToReversePairToPrivateEvidence(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	base := []string{"-principal", "analyst-a",
		"-snapshot", firstPath, "-compare-saved-case", firstCase,
		"-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondSHA}
	cohortCmd := append(append([]string{}, base...),
		"-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", "5710", "-pair-agent-id", "001")
	var cohortOut bytes.Buffer
	if err := runDemoCLI(context.Background(), cohortCmd, &cohortOut); err != nil {
		t.Fatalf("authorized rule+agent pivot rejected: %v", err)
	}
	var cohort struct {
		RuleID string `json:"rule_id"`
		AgentID string `json:"agent_id"`
		Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
	}
	if err := json.Unmarshal(cohortOut.Bytes(), &cohort); err != nil { t.Fatal(err) }
	if cohort.RuleID != "5710" || cohort.AgentID != "001" || len(cohort.Pairs) != 2 {
		t.Fatalf("wrong authorized pivot cohort: %+v", cohort)
	}
	peerID := cohort.Pairs[0].SecondEntryID
	for _, pair := range cohort.Pairs {
		if pair.SecondEntryID != peerID || pair.RuleID != "5710" || pair.AgentID != "001" {
			t.Fatalf("source cohort mixed or missing second-side peer: %+v", pair)
		}
	}

	reverseCmd := append(append([]string{}, base...),
		"-find-saved-pair-second-entry", peerID,
		"-find-saved-pair-second-offset", "0")
	var reverseOut bytes.Buffer
	if err := runDemoCLI(context.Background(), reverseCmd, &reverseOut); err != nil {
		t.Fatalf("authorized second-entry reverse lookup rejected: %v", err)
	}
	var reverse reverseWazuhPairReceipt
	if err := json.Unmarshal(reverseOut.Bytes(), &reverse); err != nil { t.Fatal(err) }
	if reverse.Format != "cmdr-authorized-wazuh-second-entry-pairs-v1" ||
		len(reverse.Pairs) != 2 || reverse.HasMore {
		t.Fatalf("reverse navigation failed for chosen filtered second Entry: %+v", reverse)
	}
	expected := make(map[string]bool, len(cohort.Pairs))
	for _, pair := range cohort.Pairs { expected[pair.FirstEntryID] = true }
	for _, pair := range reverse.Pairs {
		if pair.SecondEntryID != peerID || !expected[pair.FirstEntryID] {
			t.Fatalf("reverse navigation does not reconcile with cohort search: %+v", pair)
		}
		delete(expected, pair.FirstEntryID)
	}
	if len(expected) != 0 { t.Fatalf("cohort leads absent in reverse navigation: %+v", expected) }
	for _, text := range []string{"FIRST PRIVATE", "SECOND PRIVATE", "Bearer ", "127.0.0.1:"} {
		if strings.Contains(cohortOut.String(), text) || strings.Contains(reverseOut.String(), text) {
			t.Fatalf("metadata pivot leaked source evidence: %q", text)
		}
	}
	privatePath := filepath.Join(t.TempDir(), "reopened-source-pair.json")
	inspect := append(append([]string{}, base...),
		"-inspect-saved-pair-first-entry", reverse.Pairs[0].FirstEntryID,
		"-inspect-saved-pair-second-entry", peerID,
		"-private-report", privatePath)
	var stdout bytes.Buffer
	if err := runDemoCLI(context.Background(), inspect, &stdout); err != nil {
		t.Fatalf("correlated source owner refused private re-open: %v", err)
	}
	if stdout.Len() != 0 { t.Fatalf("private evidence escaped into stdout: %q", stdout.String()) }
	raw, err := os.ReadFile(privatePath)
	if err != nil { t.Fatal(err) }
	if !bytes.Contains(raw, []byte("FIRST PRIVATE")) || !bytes.Contains(raw, []byte("SECOND PRIVATE")) {
		t.Fatal("private paired dossier did not include independently authorized source evidence")
	}
}

func TestUnifiedWazuhPairNavigationCannotCombineSourcePivotWithReverseOrPrivateMode(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	base := []string{"-principal", "analyst-a",
		"-snapshot", firstPath, "-compare-saved-case", firstCase,
		"-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondSHA}
	var page bytes.Buffer
	if err := runDemoCLI(context.Background(), append(append([]string{},base...),
		"-list-saved-wazuh-pair-offset","0","-pair-rule-id","5710","-pair-agent-id","001"), &page);err != nil {t.Fatal(err)}
	var selected wazuhPairPageReceipt
	if err := json.Unmarshal(page.Bytes(), &selected);err != nil {t.Fatal(err)}
	if len(selected.Pairs)==0 {t.Fatal("missing authorized Wazuh lead")}
	for _,flags:= range [][]string{
		{"-find-saved-pair-second-entry",selected.Pairs[0].SecondEntryID,
			"-pair-rule-id","5710","-pair-agent-id","001"},
		{"-find-saved-pair-second-entry",selected.Pairs[0].SecondEntryID,
			"-list-saved-wazuh-pair-offset","0"},
		{"-find-saved-pair-first-entry",selected.Pairs[0].FirstEntryID,
			"-find-saved-pair-second-entry",selected.Pairs[0].SecondEntryID},
		{"-list-saved-wazuh-pair-offset","0","-pair-rule-id","5710",
			"-pair-agent-id","001","-private-report",filepath.Join(t.TempDir(),"unsafe.json")},
	} {
		args:=append(append([]string{},base...),flags...)
		var output bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&output);err==nil || output.Len()!=0 {
			t.Fatalf("unsafe mixed navigation mode accepted or leaked output: %v %q",err,output.String())
		}
	}
}
