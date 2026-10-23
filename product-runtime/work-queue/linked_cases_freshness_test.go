package workqueue

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestLinkedCaseStatusFreshnessFromInvestigateNotShared(t *testing.T) {
 links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
  {CaseID:"case-fresh",LinkSource:"Shared Linking",LinkVersion:"L-1"},
  {CaseID:"case-denied",LinkSource:"Shared Linking",LinkVersion:"L-2"},
  {CaseID:"case-stale",LinkSource:"Shared Linking",LinkVersion:"L-3"},
 }}}
 cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{
  "case-fresh":{ID:"case-fresh",Tenant:"t1",Environment:"e1",Status:"investigating",
   StatusSource:"Investigate authorized Case owner",StatusVersion:"v17",
   StatusObservedAt:"2026-10-20T09:00:00Z",StatusFreshness:"current"},
  "case-stale":{ID:"case-stale",Tenant:"t1",Environment:"e1",Status:"review",
   StatusSource:"Investigate authorized Case owner",StatusVersion:"v12",
   StatusObservedAt:"2026-10-18T08:00:00Z",StatusFreshness:"stale"},
 }}
 source:=&independentlyJoinedCommandOwner{fakeOwner:fixture(),
  LinkedCaseOwnerJoin:LinkedCaseOwnerJoin{Links:links,Cases:cases}}
 handler:=NewHTTPHandler(trustedHost(),source)
 r:=testQueueRequest(handler,http.MethodGet,
  "/command/work-queue?view=incidents&selected=Incident:inc-a")
 if r.Code!=200 {t.Fatalf("source-owned status freshness rejected: %d %s",r.Code,r.Body.String())}
 var got Result
 if err:=json.Unmarshal(r.Body.Bytes(),&got);err!=nil {t.Fatal(err)}
 if got.SelectedLinkedCases==nil || len(got.SelectedLinkedCases.Cases)!=2 {
  t.Fatalf("unexpected selected link count: %+v",got.SelectedLinkedCases)
 }
 a,b:=got.SelectedLinkedCases.Cases[0],got.SelectedLinkedCases.Cases[1]
 if a.CaseID!="case-fresh" || a.StatusFreshness!="current" ||
  a.StatusObservedAt!="2026-10-20T09:00:00Z" || a.StatusSource!="Investigate authorized Case owner" ||
  b.CaseID!="case-stale" || b.StatusFreshness!="stale" ||
  b.StatusVersion!="v12" || b.StatusObservedAt!="2026-10-18T08:00:00Z" {
  t.Fatalf("owner status observation not faithfully projected: %+v %+v",a,b)
 }
 if strings.Contains(r.Body.String(),"case-denied") ||
  strings.Contains(r.Body.String(),"L-2") ||
  len(cases.calls)!=3 || len(links.calls)!=1 {
  t.Fatal("denied Case disclosed or owner authority bypassed")
 }
 // The native HTML must expose date/source for a stale Case, and still
 // autoescape text owned by the independent Investigate source.
 cases.records["case-stale"].StatusSource="Investigate <img src=x onerror=alert(1)>"
 request:=httptest.NewRequest(http.MethodGet,
  "/command/work-queue?selected=Incident:inc-a",nil)
 request.Header.Set("Accept","text/html")
 html:=httptest.NewRecorder()
 handler.ServeHTTP(html,request)
 if html.Code!=200 || !strings.Contains(html.Body.String(),"fraîcheur du statut stale") ||
  !strings.Contains(html.Body.String(),"observé 2026-10-18T08:00:00Z") ||
  !strings.Contains(html.Body.String(),"&lt;img") ||
  strings.Contains(html.Body.String(),"<img src=x onerror=") ||
  strings.Contains(html.Body.String(),"case-denied") {
  t.Fatalf("unsafe or missing source status freshness HTML: %d %s",html.Code,html.Body.String())
 }
}

func TestLinkedCaseSourceStatusFreshnessAbsentRemainsUnknown(t *testing.T) {
 links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
  {CaseID:"case-a",LinkSource:"Shared"},
 }}}
 cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{
  "case-a":{ID:"case-a",Tenant:"t1",Environment:"e1",Status:"review"},
 }}
 source:=&independentlyJoinedCommandOwner{fakeOwner:fixture(),
  LinkedCaseOwnerJoin:LinkedCaseOwnerJoin{Links:links,Cases:cases}}
 handler:=NewHTTPHandler(trustedHost(),source)
 req:=httptest.NewRequest(http.MethodGet,"/command/work-queue?selected=Incident:inc-a",nil)
 req.Header.Set("Accept","text/html")
 response:=httptest.NewRecorder()
 handler.ServeHTTP(response,req)
 if response.Code!=200 || !strings.Contains(response.Body.String(),"fraîcheur du statut unknown") ||
  !strings.Contains(response.Body.String(),"date inconnue") {
  t.Fatalf("status without observed source was incorrectly presented as current: %d %s",
   response.Code,response.Body.String())
 }
}

func TestLinkedCaseStatusOwnerInvalidFreshnessNeverLeaksPartialSelectedData(t *testing.T) {
 tests:=[]struct{name string; status *AuthorizedCaseIdentity}{
  {"missing observed at",&AuthorizedCaseIdentity{ID:"case-2",Tenant:"t1",Environment:"e1",Status:"completed",
   StatusSource:"Investigate owner",StatusFreshness:"current"}},
  {"missing status source",&AuthorizedCaseIdentity{ID:"case-2",Tenant:"t1",Environment:"e1",Status:"completed",
   StatusObservedAt:"2026-10-20T09:00:00Z",StatusFreshness:"stale"}},
  {"invalid time",&AuthorizedCaseIdentity{ID:"case-2",Tenant:"t1",Environment:"e1",Status:"completed",
   StatusSource:"Investigate owner",StatusObservedAt:"not-a-time",StatusFreshness:"stale"}},
  {"unknown status freshness",&AuthorizedCaseIdentity{ID:"case-2",Tenant:"t1",Environment:"e1",Status:"completed",
   StatusSource:"Investigate owner",StatusObservedAt:"2026-10-20T09:00:00Z",StatusFreshness:"definitely-fresh"}},
  {"control source",&AuthorizedCaseIdentity{ID:"case-2",Tenant:"t1",Environment:"e1",Status:"completed",
   StatusSource:"Investigate\nprivate",StatusObservedAt:"2026-10-20T09:00:00Z",StatusFreshness:"current"}},
 }
 for _,tc:=range tests {
  t.Run(tc.name,func(t *testing.T){
   links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
    {CaseID:"case-1",LinkSource:"Shared"},
    {CaseID:"case-2",LinkSource:"Shared"},
   }}}
   cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{
    "case-1":{ID:"case-1",Tenant:"t1",Environment:"e1",Status:"review",
      StatusSource:"Investigate",StatusObservedAt:"2026-10-20T09:00:00Z",StatusFreshness:"current"},
    "case-2":tc.status,
   }}
   source:=&independentlyJoinedCommandOwner{fakeOwner:fixture(),
    LinkedCaseOwnerJoin:LinkedCaseOwnerJoin{Links:links,Cases:cases}}
   r:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
    "/command/work-queue?selected=Incident:inc-a")
   if r.Code!=http.StatusServiceUnavailable ||
    strings.Contains(r.Body.String(),"case-1") ||
    strings.Contains(r.Body.String(),"case-2") ||
    strings.Contains(r.Body.String(),"inc-a") ||
    strings.Contains(r.Body.String(),"private") {
     t.Fatalf("malformed Investigate status leaked partial Incident/Case: %d %q",r.Code,r.Body.String())
   }
  })
 }
}
