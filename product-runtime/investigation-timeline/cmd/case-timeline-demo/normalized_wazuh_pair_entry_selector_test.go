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

func TestPrivateWazuhPairSelectExactEntryIDs(t *testing.T) {
 aPath, aCase, aHash, bPath, bCase, bHash, pair := makeSavedPairDrilldownOwners(t)
 args := []string{
  "-principal", "analyst-a", "-snapshot", aPath,
  "-compare-saved-case", aCase, "-expect-snapshot-sha256", aHash,
  "-compare-saved-with-snapshot", bPath, "-compare-saved-with-case", bCase,
  "-expect-other-snapshot-sha256", bHash,
  "-inspect-saved-pair-first-entry", pair.FirstEntryID,
  "-inspect-saved-pair-second-entry", pair.SecondEntryID,
 }
 dest := filepath.Join(t.TempDir(), "pair-private.json")
 call := append(append([]string(nil), args...), "-private-report", dest)
 var stdout bytes.Buffer
 if err := runDemoCLI(context.Background(), call, &stdout); err != nil {
  t.Fatalf("authorized Entry pair should produce private dossier: %v", err)
 }
 if stdout.Len() != 0 { t.Fatal("sensitive pair disclosed to stdout") }
 raw, err := os.ReadFile(dest)
 if err != nil { t.Fatal(err) }
 var doc map[string]any
 if err := json.Unmarshal(raw, &doc); err != nil { t.Fatal(err) }
 pairJSON, ok := doc["pair"].(map[string]any)
 if !ok || doc["format"] != "cmdr-private-verified-wazuh-pair-v2" ||
   pairJSON["first_entry_id"] != pair.FirstEntryID ||
   pairJSON["second_entry_id"] != pair.SecondEntryID ||
   !bytes.Contains(raw, []byte("PRIVATE FIRST OWNER PAIR")) ||
   !bytes.Contains(raw, []byte("PRIVATE SECOND OWNER PAIR")) {
  t.Fatal("selected pair private evidence not retained with original identities")
 }
 if os.PathSeparator == '/' {
  info, err := os.Stat(dest)
  if err != nil || info.Mode().Perm() != 0600 { t.Fatalf("private mode not 0600: %v", err) }
 }
 stdout.Reset()
 if err := runDemoCLI(context.Background(), call, &stdout); err == nil || stdout.Len()!=0 {
  t.Fatal("existing artifact replaced or leaked to stdout")
 }
 savedAgain, err := os.ReadFile(dest)
 if err != nil || !bytes.Equal(raw, savedAgain) { t.Fatal("private artifact mutated") }

 v3dest := filepath.Join(t.TempDir(), "pair-private-v3.json")
 v3args := append(append([]string(nil), args...), "-include-pair-generic-evidence",
  "-private-report", v3dest)
 stdout.Reset()
 if err := runDemoCLI(context.Background(), v3args, &stdout); err != nil || stdout.Len()!=0 {
  t.Fatalf("opt-in private v3 pair failed or printed sensitive data: %v", err)
 }
 v3raw, err := os.ReadFile(v3dest)
 if err != nil { t.Fatal(err) }
 var v3 map[string]any
 if json.Unmarshal(v3raw, &v3) != nil || v3["format"] != "cmdr-private-verified-wazuh-pair-v3" {
  t.Fatal("opt-in private v3 format lost")
 }
}

func TestPrivateWazuhPairExactEntryIDsRejectUnsafeSelectors(t *testing.T) {
 aPath, aCase, aHash, bPath, bCase, bHash, pair := makeSavedPairDrilldownOwners(t)
 base := []string{
  "-principal", "analyst-a", "-snapshot", aPath,
  "-compare-saved-case", aCase, "-expect-snapshot-sha256", aHash,
  "-compare-saved-with-snapshot", bPath, "-compare-saved-with-case", bCase,
  "-expect-other-snapshot-sha256", bHash,
  "-inspect-saved-pair-first-entry", pair.FirstEntryID,
  "-inspect-saved-pair-second-entry", pair.SecondEntryID,
 }
 cases := []struct {name,key,value string; drop bool; extra []string}{
  {name:"wrong first",key:"-inspect-saved-pair-first-entry",value:"foreign-entry"},
  {name:"wrong second",key:"-inspect-saved-pair-second-entry",value:"foreign-entry"},
  {name:"missing first",key:"-inspect-saved-pair-first-entry",drop:true},
  {name:"missing second",key:"-inspect-saved-pair-second-entry",drop:true},
  {name:"empty first",key:"-inspect-saved-pair-first-entry",value:""},
  {name:"incorrect other owner SHA",key:"-expect-other-snapshot-sha256",value:strings.Repeat("0",64)},
  {name:"foreign principal",key:"-principal",value:"foreign-analyst"},
  {name:"mixed ordinal and ids",extra:[]string{"-inspect-saved-pair","1"}},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   selected:=append([]string(nil),base...)
   if tc.key!="" {
    for i:=0;i+1<len(selected);i++ {
     if selected[i]==tc.key {
      if tc.drop { selected=append(selected[:i],selected[i+2:]...) } else { selected[i+1]=tc.value }
      break
     }
    }
   }
   selected=append(selected,tc.extra...)
   dest:=filepath.Join(t.TempDir(),"must-not-exist.json")
   selected=append(selected,"-private-report",dest)
   var stdout bytes.Buffer
   if err:=runDemoCLI(context.Background(),selected,&stdout);err==nil || stdout.Len()!=0 {
    t.Fatalf("denied selector succeeded or printed private data: %v",err)
   }
   if _,err:=os.Lstat(dest);!os.IsNotExist(err) {
    t.Fatal("denied selector created private artifact")
   }
  })
 }
}
