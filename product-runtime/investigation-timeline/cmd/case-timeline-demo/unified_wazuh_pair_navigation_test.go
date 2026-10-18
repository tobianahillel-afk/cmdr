package main

import (
 "bytes"
 "context"
 "encoding/json"
 "path/filepath"
 "strings"
 "testing"
)

// Exercise BOTH previously separate source-authorized navigation paths
// against the same real local saved Case archives, and forbid mixed intents.
func TestUnifiedWazuhPivotAndReverseNavigation(t *testing.T) {
 firstPath,firstCase,firstSHA,secondPath,secondCase,secondSHA:=savedWazuhPairPivotFixture(t)
 base:=[]string{"-principal","analyst-a","-snapshot",firstPath,
   "-compare-saved-case",firstCase,"-expect-snapshot-sha256",firstSHA,
   "-compare-saved-with-snapshot",secondPath,"-compare-saved-with-case",secondCase,
   "-expect-other-snapshot-sha256",secondSHA}
 run:=func(extra ...string) []byte {
  t.Helper()
  var out bytes.Buffer
  if err:=runDemoCLI(context.Background(),append(append([]string{},base...),extra...),&out);err!=nil {t.Fatal(err)}
  if strings.Contains(out.String(),"PRIVATE") || strings.Contains(out.String(),"Bearer ") {t.Fatal("private preview reached metadata-only stdout")}
  return out.Bytes()
 }
 var filtered wazuhPairPageReceipt
 if err:=json.Unmarshal(run("-list-saved-wazuh-pair-offset","0","-pair-rule-id","5710","-pair-agent-id","001"),&filtered);err!=nil {t.Fatal(err)}
 if len(filtered.Pairs)!=2 {t.Fatalf("expected two selected Wazuh pairs: %+v",filtered)}
 var backward reverseWazuhPairReceipt
 if err:=json.Unmarshal(run("-find-saved-pair-second-entry",filtered.Pairs[0].SecondEntryID),&backward);err!=nil {t.Fatal(err)}
 if backward.Format!="cmdr-authorized-wazuh-second-entry-pairs-v1" || len(backward.Pairs)!=2 {
  t.Fatalf("reverse Wazuh navigation did not preserve both original pairs: %+v",backward)
 }
 expected:=map[string]bool{}
 for _,p:=range filtered.Pairs {expected[p.FirstEntryID]=true}
 for _,p:=range backward.Pairs {
  if !expected[p.FirstEntryID] || p.SecondEntryID!=filtered.Pairs[0].SecondEntryID {t.Fatalf("unrelated source exposed in reverse navigation: %+v",p)}
 }
 for _,extra:=range [][]string{
  {"-find-saved-pair-second-entry",filtered.Pairs[0].SecondEntryID,"-list-saved-wazuh-pair-offset","0"},
  {"-find-saved-pair-second-entry",filtered.Pairs[0].SecondEntryID,"-pair-rule-id","5710","-pair-agent-id","001"},
  {"-find-saved-pair-second-offset","0"},
  {"-find-saved-pair-second-entry",filtered.Pairs[0].SecondEntryID,"-find-saved-pair-first-entry",filtered.Pairs[0].FirstEntryID},
  {"-find-saved-pair-second-entry",filtered.Pairs[0].SecondEntryID,"-private-report",filepath.Join(t.TempDir(),"forbidden.json")},
 } {
  var out bytes.Buffer
  if err:=runDemoCLI(context.Background(),append(append([]string{},base...),extra...),&out);err==nil || out.Len()!=0 {
   t.Fatalf("mixed discovery escaped intent guard: %v %q",err,out.String())
  }
 }
}
