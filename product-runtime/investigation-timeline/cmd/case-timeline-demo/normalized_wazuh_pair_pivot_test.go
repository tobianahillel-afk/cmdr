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
)

// Construct TWO distinct authorized Wazuh cohorts and separate Case owners.
// A pivot must not include pairs belonging to the other cohort.
func savedWazuhPairPivotFixture(t *testing.T) (string,string,string,string,string,string) {
	t.Helper()
	firstOwner, firstCase, _, secondOwner, secondCase, _ := buildTwoVerifiedLocalCaseArchives(t)
	firstFeed := writeJSONLEvents(t,
		"{\"id\":\"a1\",\"timestamp\":\"2026-10-17T10:00:01Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"FIRST PRIVATE A\",\"level\":9}}\n"+
		"{\"id\":\"a2\",\"timestamp\":\"2026-10-17T10:00:02Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"FIRST PRIVATE B\",\"level\":9}}\n"+
		"{\"id\":\"b1\",\"timestamp\":\"2026-10-17T10:00:03Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"8020\",\"description\":\"FIRST PRIVATE C\",\"level\":9}}\n"+
		"{\"id\":\"b2\",\"timestamp\":\"2026-10-17T10:00:04Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"8020\",\"description\":\"FIRST PRIVATE D\",\"level\":9}}\n")
	secondFeed := writeJSONLEvents(t,
		"{\"id\":\"peer-a\",\"timestamp\":\"2026-10-17T10:00:07Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"SECOND PRIVATE A\",\"level\":8}}\n"+
		"{\"id\":\"peer-b\",\"timestamp\":\"2026-10-17T10:00:08Z\",\"agent\":{\"id\":\"002\"},\"rule\":{\"id\":\"8020\",\"description\":\"SECOND PRIVATE B\",\"level\":8}}\n")
	firstPath := filepath.Join(t.TempDir(), "first-pivot-owner.json")
	secondPath := filepath.Join(t.TempDir(), "second-pivot-owner.json")
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",firstOwner,firstCase,[]string{firstFeed},firstPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	if err:=writePersistedWazuhCase(context.Background(),"analyst-a",secondOwner,secondCase,[]string{secondFeed},secondPath,&bytes.Buffer{});err!=nil {t.Fatal(err)}
	firstBytes,err:=os.ReadFile(firstPath);if err!=nil {t.Fatal(err)}
	secondBytes,err:=os.ReadFile(secondPath);if err!=nil {t.Fatal(err)}
	return firstPath,firstCase,fmt.Sprintf("%x",sha256.Sum256(firstBytes)),
		secondPath,secondCase,fmt.Sprintf("%x",sha256.Sum256(secondBytes))
}

func TestSourceAuthorizedWazuhRuleAgentPivotEndToEnd(t *testing.T) {
	firstPath,firstCase,firstSHA,secondPath,secondCase,secondSHA:=savedWazuhPairPivotFixture(t)
	base:=[]string{"-principal","analyst-a","-snapshot",firstPath,
		"-compare-saved-case",firstCase,"-expect-snapshot-sha256",firstSHA,
		"-compare-saved-with-snapshot",secondPath,"-compare-saved-with-case",secondCase,
		"-expect-other-snapshot-sha256",secondSHA,
		"-list-saved-wazuh-pair-offset","0"}
	read:=func(filters ...string) struct {
		RuleID string `json:"rule_id"`
		AgentID string `json:"agent_id"`
		Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
		HasMore bool `json:"has_more"`
	} {
		t.Helper()
		args:=append(append([]string{},base...),filters...)
		var stdout bytes.Buffer
		if err:=runDemoCLI(context.Background(),args,&stdout);err!=nil {
			t.Fatalf("source-authorized pivot rejected: %v",err)
		}
		for _,secret:=range []string{"PRIVATE A","PRIVATE B","PRIVATE C","PRIVATE D","Bearer ","127.0.0.1:"} {
			if strings.Contains(stdout.String(),secret) {t.Fatalf("source preview leaked into pivot metadata: %q",secret)}
		}
		var got struct {
			RuleID string `json:"rule_id"`
			AgentID string `json:"agent_id"`
			Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
			HasMore bool `json:"has_more"`
		}
		if err:=json.Unmarshal(stdout.Bytes(),&got);err!=nil {t.Fatal(err)}
		return got
	}
	global:=read()
	a:=read("-pair-rule-id","5710","-pair-agent-id","001")
	b:=read("-pair-rule-id","8020","-pair-agent-id","002")
	empty:=read("-pair-rule-id","8020","-pair-agent-id","001")
	if len(global.Pairs)!=4 || len(a.Pairs)!=2 || len(b.Pairs)!=2 ||
		len(empty.Pairs)!=0 || global.HasMore || a.HasMore || b.HasMore ||
		a.RuleID!="5710" || a.AgentID!="001" ||
		b.RuleID!="8020" || b.AgentID!="002" {
		t.Fatalf("exact source cohort pivot crossed sources or lost matches: global=%+v a=%+v b=%+v empty=%+v",global,a,b,empty)
	}
	for _,p:=range b.Pairs {
		if p.RuleID!="8020" || p.AgentID!="002" {
			t.Fatalf("pair outside selected cohort escaped: %+v",p)
		}
	}
	// Source-linked pivot output is actionable with the existing separate
	// exact-ID typed-source reauthorization, not a shortcut to raw previews.
	dest:=filepath.Join(t.TempDir(),"authorized-pivot-pair.json")
	cmd:=append(append([]string{},base[:len(base)-2]...),
		"-inspect-saved-pair-first-entry",b.Pairs[0].FirstEntryID,
		"-inspect-saved-pair-second-entry",b.Pairs[0].SecondEntryID,
		"-private-report",dest)
	var stdout bytes.Buffer
	if err:=runDemoCLI(context.Background(),cmd,&stdout);err!=nil {t.Fatal(err)}
	if stdout.Len()!=0 {t.Fatal("private dossier escaped on stdout")}
	private,err:=os.ReadFile(dest)
	if err!=nil {t.Fatal(err)}
	if !bytes.Contains(private,[]byte("FIRST PRIVATE C")) ||
		!bytes.Contains(private,[]byte("SECOND PRIVATE B")) {
		t.Fatal("selected cohort metadata did not reopen canonical Wazuh pair sources")
	}
}

func TestSourceAuthorizedWazuhPivotConflictsAndOwnerDenials(t *testing.T) {
	firstPath,firstCase,firstSHA,secondPath,secondCase,secondSHA:=savedWazuhPairPivotFixture(t)
	base:=[]string{"-principal","analyst-a","-snapshot",firstPath,
		"-compare-saved-case",firstCase,"-expect-snapshot-sha256",firstSHA,
		"-compare-saved-with-snapshot",secondPath,"-compare-saved-with-case",secondCase,
		"-expect-other-snapshot-sha256",secondSHA,
		"-list-saved-wazuh-pair-offset","0",
		"-pair-rule-id","5710","-pair-agent-id","001"}
	for _,tc:=range []struct{name string;extra []string}{
		{"missing agent",[]string{"-pair-agent-id",""}},
		{"missing rule",[]string{"-pair-rule-id",""}},
		{"oversized rule",[]string{"-pair-rule-id",strings.Repeat("5",129)}},
		{"wrong second digest",[]string{"-expect-other-snapshot-sha256",strings.Repeat("0",64)}},
		{"foreign principal",[]string{"-principal","not-the-owner"}},
		{"foreign second Case",[]string{"-compare-saved-with-case","not-the-owner-Case"}},
		{"finder mixed",[]string{"-find-saved-pair-first-entry","random"}},
		{"pair inspection mixed",[]string{"-inspect-saved-pair","1"}},
		{"preview mixed",[]string{"-include-pair-generic-evidence"}},
		{"private output mixed",[]string{"-private-report",filepath.Join(t.TempDir(),"must-not-exist.json")}},
	} {
		t.Run(tc.name,func(t *testing.T) {
			args:=append([]string{},base...)
			// Replace an existing argument when appropriate, otherwise append.
			flag,value:=tc.extra[0],"";if len(tc.extra)>1 {value=tc.extra[1]}
			replaced:=false
			for i:=0;i+1<len(args);i++ {
				if args[i]==flag {args[i+1]=value;replaced=true;break}
			}
			if !replaced {args=append(args,tc.extra...)}
			var output bytes.Buffer
			if err:=runDemoCLI(context.Background(),args,&output);err==nil || output.Len()!=0 {
				t.Fatalf("denied pivot leaked metadata: %v %q",err,output.String())
			}
		})
	}
	// A changed but valid JSON second snapshot must never emit first-Case
	// Wazuh cohort metadata before integrity of BOTH owners is checked.
	data,err:=os.ReadFile(secondPath);if err!=nil {t.Fatal(err)}
	if err:=os.WriteFile(secondPath,append(data,'\n'),0600);err!=nil {t.Fatal(err)}
	var output bytes.Buffer
	if err:=runDemoCLI(context.Background(),base,&output);err==nil || output.Len()!=0 {
		t.Fatalf("tampered second Case produced first-Case pivot data: %v %q",err,output.String())
	}
}
