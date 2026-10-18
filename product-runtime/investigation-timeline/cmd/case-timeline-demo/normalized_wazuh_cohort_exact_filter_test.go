package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// One authorized analyst journey: narrow a large source-authorized cohort
// catalog by a known rule OR agent before following existing exact pair pivots.
func TestSavedWazuhCohortExactDiscoveryAndPairPivot(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	base := []string{
		"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase, "-expect-other-snapshot-sha256", secondSHA,
	}
	call := func(extra ...string) (bytes.Buffer, error) {
		t.Helper()
		var out bytes.Buffer
		args := append(append([]string{}, base...), extra...)
		return out, runDemoCLI(context.Background(), args, &out)
	}
	// A single exact rule or agent suffices; returning a source label does not
	// authorize reading its private payload.
	for _, tc := range []struct {
		args []string
		rule, agent string
	}{
		{[]string{"-catalog-rule-id", "5710"}, "5710", "001"},
		{[]string{"-catalog-agent-id", "002"}, "8020", "002"},
		{[]string{"-catalog-rule-id", "8020", "-catalog-agent-id", "002"}, "8020", "002"},
	} {
		args := append([]string{"-list-saved-wazuh-cohorts-offset", "0"}, tc.args...)
		out, err := call(args...)
		if err != nil { t.Fatalf("catalog filter %v: %v", tc.args, err) }
		var got struct {
			Format string `json:"format"`
			Offset int `json:"offset"`
			HasMore bool `json:"has_more"`
			Cohorts []verifiedWazuhCohortCount `json:"cohorts"`
		}
		if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
		if got.Format != "cmdr-authorized-wazuh-cohort-catalog-v1" ||
			got.HasMore || got.Offset != 0 || len(got.Cohorts) != 1 ||
			got.Cohorts[0].RuleID != tc.rule || got.Cohorts[0].AgentID != tc.agent ||
			got.Cohorts[0].MatchedPairs != 2 {
			t.Fatalf("unexpected authorized exact filter %v: %+v", tc.args, got)
		}
		for _, secret := range []string{"FIRST PRIVATE", "SECOND PRIVATE", "Bearer "} {
			if strings.Contains(out.String(), secret) { t.Fatal("private source preview disclosed") }
		}
		pair, err := call("-list-saved-wazuh-pair-offset", "0",
			"-pair-rule-id", tc.rule, "-pair-agent-id", tc.agent)
		if err != nil { t.Fatal(err) }
		var verified wazuhPairPageReceipt
		if err := json.Unmarshal(pair.Bytes(), &verified); err != nil { t.Fatal(err) }
		if len(verified.Pairs) != got.Cohorts[0].MatchedPairs { t.Fatal("catalog/pair count mismatch") }
	}
	none, err := call("-list-saved-wazuh-cohorts-offset", "0", "-catalog-rule-id", "absent")
	if err != nil { t.Fatal(err) }
	var empty struct{ Cohorts []verifiedWazuhCohortCount `json:"cohorts"` }
	if err := json.Unmarshal(none.Bytes(), &empty); err != nil { t.Fatal(err) }
	if len(empty.Cohorts) != 0 { t.Fatal("nonexistent authorized cohort was fabricated") }

	for _, extra := range [][]string{
		{"-catalog-rule-id", "5710"},
		{"-list-saved-wazuh-cohorts-offset", "0", "-catalog-agent-id", ""},
		{"-list-saved-wazuh-cohorts-offset", "0", "-catalog-rule-id", strings.Repeat("x", 129)},
		{"-list-saved-wazuh-cohorts-offset", "0", "-catalog-agent-id", "001", "-list-saved-wazuh-pair-offset", "0"},
		{"-list-saved-wazuh-cohorts-offset", "0", "-catalog-rule-id", "5710", "-private-report", "/tmp/not-created.json"},
	} {
		out, err := call(extra...)
		if err == nil || out.Len() != 0 { t.Fatalf("mixed/bad mode disclosed data %q: %v", extra, err) }
	}
	denied := append(append([]string{}, base...), "-list-saved-wazuh-cohorts-offset", "0",
		"-catalog-rule-id", "5710", "-principal", "unauthorized")
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), denied, &out); err == nil || out.Len() != 0 {
		t.Fatal("unauthorized catalog filter disclosed source metadata")
	}
}

func TestAuthorizedWazuhCohortExactFiltersPrecedePagination(t *testing.T) {
	left := map[authorizedWazuhCohort][]authorizedWazuhObservation{}
	right := map[authorizedWazuhCohort][]authorizedWazuhObservation{}
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 25; i++ {
		key := authorizedWazuhCohort{rule:fmt.Sprintf("r-%02d",i),agent:fmt.Sprintf("a-%d",i%3)}
		left[key] = []authorizedWazuhObservation{{id:"l",at:now}}
		right[key] = []authorizedWazuhObservation{{id:"r",at:now.Add(time.Second)}}
	}
	filtered, more, err := authorizedWazuhCohortCatalog(left, right, 0, "", "a-2")
	if err != nil || more || len(filtered) != 8 { t.Fatalf("filtered cohort catalog: %d %v %v",len(filtered),more,err) }
	for _, x := range filtered { if x.AgentID != "a-2" { t.Fatal("filter leak") } }
	// Offset applies to matches rather than to the global ranking.
	next, more, err := authorizedWazuhCohortCatalog(left, right, 5, "", "a-2")
	if err != nil || more || len(next) != 3 || next[0].RuleID != "r-17" {
		t.Fatalf("filtered pagination: %+v %v %v", next,more,err)
	}
	both, more, err := authorizedWazuhCohortCatalog(left, right, 0, "r-11", "a-2")
	if err != nil || more || len(both) != 1 || both[0].RuleID != "r-11" { t.Fatal("exact pair filter") }
}
