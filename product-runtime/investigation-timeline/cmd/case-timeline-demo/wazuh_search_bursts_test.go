package main

import (
 "context"
 "encoding/json"
 "io"
 "net/http"
 "net/http/httptest"
 "net/url"
 "strings"
 "testing"
 "time"

 search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
)

func TestSearchBurstsReuseAuthorizedWazuhSearchJob(t *testing.T) {
 const caseID = "forensic-case-on-disk"
 owner := snapshotFixture(t)
 paths := []string{
  writeJSONLEvents(t, "{\"id\":\"a.1\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\",\"name\":\"SECRET-HOST\"},\"rule\":{\"id\":\"5710\",\"level\":5,\"description\":\"SECRET-PAYLOAD\"}}\n"+
   "{\"id\":\"a.2\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":7}}\n"),
  writeJSONLEvents(t, "{\"id\":\"b.1\",\"timestamp\":\"2026-10-17T10:00:40Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":12}}\n"+
   "{\"id\":\"b.2\",\"timestamp\":\"2026-10-17T10:02:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"level\":15}}\n"),
 }
 data, imported, levels, err := loadMultipleWazuhFeedsWithLevels(context.Background(), "analyst-a", owner, caseID, paths)
 if err != nil { t.Fatal(err) }
 token, err := randomBearer()
 if err != nil { t.Fatal(err) }
 base, err := newDemoHandlerWithData(token, "analyst-a", data)
 if err != nil { t.Fatal(err) }
 handler := newDemoWazuhCohortReceiptHandler(token, "analyst-a", caseID, data, imported, nil, levels,
  base, func(context.Context, io.Writer) error { return nil })
 query := url.Values{
  "case_id": {caseID}, "field": {"rule_id"}, "value": {"5710"},
  "from": {"2026-10-17T09:59:00Z"}, "to": {"2026-10-17T10:03:00Z"},
 }
 req := httptest.NewRequest(http.MethodGet, "/demo-wazuh-search?"+query.Encode(), nil)
 req.Header.Set("Authorization", "Bearer "+token)
 response := httptest.NewRecorder()
 handler.ServeHTTP(response, req)
 if response.Code != http.StatusOK { t.Fatalf("authorized real search failed: %d %s",response.Code,response.Body.String()) }
 var receipt struct {
  Job search.SearchJob `json:"job"`
  Hits []wazuhSearchHit `json:"search_hits"`
  Bursts []wazuhSearchBurst `json:"search_bursts"`
 }
 if err := json.Unmarshal(response.Body.Bytes(),&receipt);err!=nil {t.Fatal(err)}
 if receipt.Job.State!=search.JobCompleted || len(receipt.Job.ResultRefs)!=4 ||
  len(receipt.Hits)!=4 || len(receipt.Bursts)!=1 {
  t.Fatalf("real SearchJob, 4 authorized hits and one lead expected: %+v",receipt)
 }
 burst:=receipt.Bursts[0]
 if burst.AgentID!="001" || burst.RuleID!="5710" || burst.AlertCount!=3 ||
  burst.DistinctFeeds!=2 || burst.FirstUTC!="2026-10-17T10:00:00Z" ||
  burst.LastUTC!="2026-10-17T10:00:40Z" || burst.WindowSeconds!=60 ||
  len(burst.EntryIDs)!=3 {t.Fatalf("incorrect authorized 60s burst: %+v",burst)}
 found:=make(map[string]bool)
 for _,ref:=range receipt.Job.ResultRefs {found[ref.Ref]=true}
 for _,id:=range burst.EntryIDs {if !found[id] {t.Fatalf("burst references event outside completed SearchJob: %q",id)}}
 if found[importedWazuhID(caseID,"b.2",2)] && strings.Contains(strings.Join(burst.EntryIDs,","), importedWazuhID(caseID,"b.2",2)) {
  t.Fatal("alert beyond strongest window incorrectly cited as burst member")
 }
 for _,sensitive:=range []string{"SECRET-HOST","SECRET-PAYLOAD",token,"Bearer "} {
  if strings.Contains(response.Body.String(),sensitive) {t.Fatalf("burst receipt leaked %q",sensitive)}
 }
 unauthorized:=httptest.NewRequest(http.MethodGet,"/demo-wazuh-search?"+query.Encode(),nil)
 unauthorized.Header.Set("Authorization","Bearer invalid")
 denied:=httptest.NewRecorder()
 handler.ServeHTTP(denied,unauthorized)
 if denied.Code==http.StatusOK || strings.Contains(denied.Body.String(),burst.AgentID) {
  t.Fatalf("unauthorized burst disclosed: %d %s",denied.Code,denied.Body.String())
 }
}

func TestSearchBurstsSlidingWindowBoundaryRankingAndUncertain(t *testing.T) {
 at:=time.Date(2026,10,17,10,0,0,0,time.UTC)
 makeHit:=func(id,agent,rule string, instant time.Time,feed int) wazuhSearchHit {
  return wazuhSearchHit{EntryID:id,AgentID:agent,RuleID:rule,Feed:feed,SourceTimeUTC:instant.Format(time.RFC3339Nano)}
 }
 hits:=[]wazuhSearchHit{
  makeHit("a","001","5710",at,1),
  makeHit("b","001","5710",at.Add(30*time.Second),1),
  makeHit("c","001","5710",at.Add(60*time.Second),2),
  makeHit("late","001","5710",at.Add(60*time.Second+time.Nanosecond),3),
  makeHit("wrong-rule","001","5500",at,1),
  makeHit("other-agent","002","5710",at,1),
  {EntryID:"unknown-time",AgentID:"001",RuleID:"5710",SourceTimeUTC:"bad"},
 }
 bursts:=deriveWazuhSearchBursts(hits)
 if len(bursts)!=1 || bursts[0].AlertCount!=3 || bursts[0].DistinctFeeds!=2 ||
  strings.Contains(strings.Join(bursts[0].EntryIDs,","),"late") {
  t.Fatalf("inclusive exact 60-second window mishandled: %+v",bursts)
 }
 if len(deriveWazuhSearchBursts(hits[:2]))!=0 {
  t.Fatal("two alerts should not form a burst")
 }
}
