package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSavedWazuhCohortCatalogTwoAuthorizedCases(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	base := []string{
		"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondSHA,
	}
	read := func(offset string) struct {
		Format string `json:"format"`
		HasMore bool `json:"has_more"`
		Cohorts []verifiedWazuhCohortCount `json:"cohorts"`
	} {
		t.Helper()
		var out bytes.Buffer
		args := append(append([]string{}, base...), "-list-saved-wazuh-cohorts-offset", offset)
		if err := runDemoCLI(context.Background(), args, &out); err != nil {
			t.Fatal(err)
		}
		for _, secret := range []string{"FIRST PRIVATE", "SECOND PRIVATE", "Bearer ", "127.0.0.1:"} {
			if strings.Contains(out.String(), secret) {
				t.Fatalf("source preview disclosed in catalog: %q", secret)
			}
		}
		var got struct {
			Format string `json:"format"`
			HasMore bool `json:"has_more"`
			Cohorts []verifiedWazuhCohortCount `json:"cohorts"`
		}
		if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
		return got
	}
	page := read("0")
	if page.Format != "cmdr-authorized-wazuh-cohort-catalog-v1" ||
		page.HasMore || len(page.Cohorts) != 2 ||
		page.Cohorts[0].RuleID != "5710" || page.Cohorts[0].AgentID != "001" ||
		page.Cohorts[1].RuleID != "8020" || page.Cohorts[1].AgentID != "002" ||
		page.Cohorts[0].MatchedPairs != 2 || page.Cohorts[1].MatchedPairs != 2 {
		t.Fatalf("wrong authorized cohort catalog: %+v", page)
	}
	var drill bytes.Buffer
	args := append(append([]string{}, base...),
		"-list-saved-wazuh-pair-offset", "0",
		"-pair-rule-id", page.Cohorts[1].RuleID,
		"-pair-agent-id", page.Cohorts[1].AgentID)
	if err := runDemoCLI(context.Background(), args, &drill); err != nil { t.Fatal(err) }
	var pairs wazuhPairPageReceipt
	if err := json.Unmarshal(drill.Bytes(), &pairs); err != nil { t.Fatal(err) }
	if len(pairs.Pairs) != page.Cohorts[1].MatchedPairs {
		t.Fatalf("catalog count and existing exact cohort pagination disagree: %d != %d",
			page.Cohorts[1].MatchedPairs, len(pairs.Pairs))
	}
	after := read("1")
	if len(after.Cohorts) != 1 || after.Cohorts[0].RuleID != "8020" || after.HasMore {
		t.Fatalf("unstable catalog offset: %+v", after)
	}
	if beyond := read("2"); len(beyond.Cohorts) != 0 || beyond.HasMore {
		t.Fatalf("exhausted catalog returned extra records: %+v", beyond)
	}
}

func TestSavedWazuhCohortCatalogDeniesUnauthorizedAndMixedModes(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	base := []string{"-principal", "analyst-a", "-snapshot", firstPath,
		"-compare-saved-case", firstCase, "-expect-snapshot-sha256", firstSHA,
		"-compare-saved-with-snapshot", secondPath,
		"-compare-saved-with-case", secondCase,
		"-expect-other-snapshot-sha256", secondSHA,
		"-list-saved-wazuh-cohorts-offset", "0"}
	for _, extra := range [][]string{
		{"-principal", "not-authorized"},
		{"-expect-other-snapshot-sha256", strings.Repeat("0",64)},
		{"-compare-saved-with-case", "other-case"},
		{"-list-saved-wazuh-cohorts-offset", "-1"},
		{"-list-saved-wazuh-cohorts-offset", "5001"},
		{"-list-saved-wazuh-pair-offset", "0"},
		{"-pair-rule-id", "5710", "-pair-agent-id", "001"},
		{"-find-saved-pair-second-entry", "peer-a"},
		{"-inspect-saved-pair", "1"},
		{"-include-pair-generic-evidence"},
		{"-private-report", filepath.Join(t.TempDir(), "forbidden.json")},
	} {
		args := append([]string{}, base...)
		for j:=0;j<len(extra);j++ {
			if j+1<len(extra) && strings.HasPrefix(extra[j], "-") &&
				extra[j+1]!="-include-pair-generic-evidence" &&
				!strings.HasPrefix(extra[j+1], "-") {
				updated:=false
				for k:=0;k+1<len(args);k++ {
					if args[k]==extra[j] {args[k+1]=extra[j+1];updated=true;break}
				}
				if !updated {args=append(args,extra[j],extra[j+1])}
				j++
			} else {args=append(args,extra[j])}
		}
		var out bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&out);err==nil||out.Len()!=0 {
			t.Fatalf("unauthorized or mixed catalog emitted output: args=%q err=%v out=%q",extra,err,out.String())
		}
	}
}

func TestWazuhCohortCatalogIsDeterministicAndPaginatesTwelve(t *testing.T) {
	left:=map[authorizedWazuhCohort][]authorizedWazuhObservation{}
	right:=map[authorizedWazuhCohort][]authorizedWazuhObservation{}
	start:=time.Date(2026,10,17,10,0,0,0,time.UTC)
	for i:=0;i<12;i++ {
		key:=authorizedWazuhCohort{rule:fmt.Sprintf("%04d",i),agent:"001"}
		left[key]=[]authorizedWazuhObservation{{id:"first",at:start}}
		right[key]=[]authorizedWazuhObservation{{id:"second",at:start.Add(time.Second)}}
	}
	first,more,err:=authorizedWazuhCohortCatalog(left,right,0)
	if err!=nil||len(first)!=10||!more||first[0].RuleID!="0000" {t.Fatalf("first: %+v %v %v",first,more,err)}
	second,more,err:=authorizedWazuhCohortCatalog(left,right,10)
	if err!=nil||len(second)!=2||more||second[0].RuleID!="0010" {t.Fatalf("second: %+v %v %v",second,more,err)}
	for i:=0;i<12;i++ {
		left[authorizedWazuhCohort{rule:fmt.Sprintf("%04d",i),agent:"001"}]=[]authorizedWazuhObservation{{id:"first",at:start.Add(-2*time.Hour)}}
	}
	none,more,err:=authorizedWazuhCohortCatalog(left,right,0)
	if err!=nil||len(none)!=0||more {t.Fatal("unmatched cohorts should not appear")}
}

func TestWazuhCohortCatalogChangedSecondSnapshotNeverEmitsMetadata(t *testing.T) {
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA := savedWazuhPairPivotFixture(t)
	data,err:=os.ReadFile(secondPath)
	if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(secondPath,append(data, byte(10)),0600);err!=nil {t.Fatal(err)}
	var out bytes.Buffer
	err=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",firstPath,
		"-compare-saved-case",firstCase,"-expect-snapshot-sha256",firstSHA,
		"-compare-saved-with-snapshot",secondPath,
		"-compare-saved-with-case",secondCase,
		"-expect-other-snapshot-sha256",secondSHA,
		"-list-saved-wazuh-cohorts-offset","0",
	},&out)
	if err==nil||out.Len()!=0 {t.Fatalf("tampered source leaked metadata: %v %q",err,out.String())}
}
