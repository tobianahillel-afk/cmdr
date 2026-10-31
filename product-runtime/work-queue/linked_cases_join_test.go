package workqueue

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "strings"
 "testing"
)

type linkReferenceFixture struct {
 refs *AuthorizedIncidentCaseReferences
 err error
 calls []string
}
func (f *linkReferenceFixture) ReadAuthorizedIncidentCaseReferences(_ context.Context,p,t,e,i string)(*AuthorizedIncidentCaseReferences,error){
 f.calls=append(f.calls,p+":"+t+":"+e+":"+i)
 return f.refs,f.err
}

type caseIdentityFixture struct {
 records map[string]*AuthorizedCaseIdentity
 errAt string
 calls []string
}
func (f *caseIdentityFixture) ReadAuthorizedCaseIdentity(_ context.Context,p,t,e,id string)(*AuthorizedCaseIdentity,error){
 f.calls=append(f.calls,p+":"+t+":"+e+":"+id)
 if f.errAt==id{return nil,errors.New("PRIVATE SOURCE ERROR")}
 return f.records[id],nil
}

type independentlyJoinedCommandOwner struct {
 *fakeOwner
 LinkedCaseOwnerJoin
}

func TestLinkedCaseOwnerJoinSelectedIncidentFullHTTP(t *testing.T) {
 links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
  {CaseID:"case-001",LinkSource:"Shared Object Linking",LinkVersion:"v1"},
  {CaseID:"case-denied",LinkSource:"Shared Object Linking",LinkVersion:"v2"},
  {CaseID:"case-002",LinkSource:"Shared Object Linking",LinkVersion:"v3"},
 }}}
 cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{
  "case-001":{ID:"case-001",Tenant:"t1",Environment:"e1",Status:"investigating"},
  "case-002":{ID:"case-002",Tenant:"t1",Environment:"e1",Status:"review"},
 }}
 source:=&independentlyJoinedCommandOwner{fakeOwner:fixture(),
  LinkedCaseOwnerJoin:LinkedCaseOwnerJoin{Links:links,Cases:cases}}
 h:=NewHTTPHandler(trustedHost(),source)
 r:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
 if r.Code!=200 {t.Fatalf("selected Incident join failed: %d %s",r.Code,r.Body.String())}
 var projected Result
 if err:=json.Unmarshal(r.Body.Bytes(),&projected);err!=nil{t.Fatal(err)}
 if projected.SelectedLinkedCases==nil || len(projected.SelectedLinkedCases.Cases)!=2 ||
  projected.SelectedLinkedCases.Cases[0].CaseID!="case-001" ||
  projected.SelectedLinkedCases.Cases[0].Status!="investigating" ||
  projected.SelectedLinkedCases.Cases[1].CaseID!="case-002" {
  t.Fatalf("independent sources were not joined under selected Incident: %+v",projected.SelectedLinkedCases)
 }
 if len(links.calls)!=1 || links.calls[0]!="analyst:t1:e1:inc-a" ||
  len(cases.calls)!=3 || cases.calls[0]!="analyst:t1:e1:case-001" ||
  cases.calls[1]!="analyst:t1:e1:case-denied" {
  t.Fatalf("independent source authorization not invoked in same scope: %v %v",links.calls,cases.calls)
 }
 if strings.Contains(r.Body.String(),"case-denied") || strings.Contains(r.Body.String(),"PRIVATE") {
  t.Fatal("denied linked Case disclosed")
 }
 links.calls=nil;cases.calls=nil
 for _,route:=range []string{
  "/command/work-queue?view=incidents",
  "/command/work-queue?selected=Task:task-b",
 }{
  rr:=testQueueRequest(h,http.MethodGet,route)
  if rr.Code!=200 || len(links.calls)!=0 || len(cases.calls)!=0 {
   t.Fatalf("unselected/Task triggered forbidden relation read: %s %d %v %v",route,rr.Code,links.calls,cases.calls)
  }
 }
}

func TestLinkedCaseJoinRejectsInvalidLinkCohortBeforeCaseReads(t *testing.T){
 samples:=[]struct{name string;refs []LinkedCaseReference}{
  {"duplicate",[]LinkedCaseReference{{CaseID:"case-a",LinkSource:"Shared"},{CaseID:"case-a",LinkSource:"Shared"}}},
  {"invalid identifier",[]LinkedCaseReference{{CaseID:"case\nb",LinkSource:"Shared"}}},
  {"missing provenance",[]LinkedCaseReference{{CaseID:"case-a"}}},
  {"over limit",make([]LinkedCaseReference,17)},
 }
 for _,tc:=range samples{
  t.Run(tc.name,func(t *testing.T){
   links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:tc.refs}}
   cases:=&caseIdentityFixture{}
   got,err:=(LinkedCaseOwnerJoin{Links:links,Cases:cases}).ReadAuthorizedLinkedCases(
    context.Background(),"analyst","t1","e1","inc-a")
   if err==nil || got!=nil || len(cases.calls)!=0 {
    t.Fatalf("malformed link cohort reached Investigate owner: %+v %v %v",got,err,cases.calls)
   }
  })
 }
}

func TestLinkedCaseJoinRefusesForeignOrInconsistentInvestigateResults(t *testing.T){
 for _,c:=range []AuthorizedCaseIdentity{
  {ID:"case-other",Tenant:"t1",Environment:"e1",Status:"review"},
  {ID:"case-a",Tenant:"other-tenant",Environment:"e1",Status:"review"},
  {ID:"case-a",Tenant:"t1",Environment:"other-env",Status:"review"},
  {ID:"case-a",Tenant:"t1",Environment:"e1",Status:"unknown-bad"},
 }{
  links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
   {CaseID:"case-a",LinkSource:"Shared"},
  }}}
  cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{"case-a":&c}}
  got,err:=(LinkedCaseOwnerJoin{Links:links,Cases:cases}).ReadAuthorizedLinkedCases(context.Background(),"analyst","t1","e1","inc-a")
  if err==nil || got!=nil {t.Fatalf("foreign or invalid Investigate source accepted: %+v %v",got,err)}
 }
}

func TestLinkedCaseJoinFailClosedWhenSecondSourceFails(t *testing.T){
 links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{References:[]LinkedCaseReference{
  {CaseID:"case-1",LinkSource:"Shared"},
  {CaseID:"case-2",LinkSource:"Shared"},
 }}}
 cases:=&caseIdentityFixture{records:map[string]*AuthorizedCaseIdentity{
  "case-1":{ID:"case-1",Tenant:"t1",Environment:"e1",Status:"review"},
 },errAt:"case-2"}
 source:=&independentlyJoinedCommandOwner{fakeOwner:fixture(),
  LinkedCaseOwnerJoin:LinkedCaseOwnerJoin{Links:links,Cases:cases}}
 r:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
  "/command/work-queue?selected=Incident:inc-a")
 if r.Code!=http.StatusServiceUnavailable || strings.Contains(r.Body.String(),"case-1") ||
  strings.Contains(r.Body.String(),"PRIVATE SOURCE ERROR") || strings.Contains(r.Body.String(),"inc-a"){
  t.Fatalf("partial/unauthorized Case leaked after second provider failure: %d %q",r.Code,r.Body.String())
 }
}

func TestLinkedCaseJoinUnknownVersusAuthorizedEmptyAndCancelled(t *testing.T){
 j:=LinkedCaseOwnerJoin{}
 if got,err:=j.ReadAuthorizedLinkedCases(context.Background(),"analyst","t1","e1","inc-a");err!=nil || got!=nil{
  t.Fatalf("unavailable source must be unknown not empty: %+v %v",got,err)
 }
 links:=&linkReferenceFixture{refs:&AuthorizedIncidentCaseReferences{}}
 cases:=&caseIdentityFixture{}
 j=LinkedCaseOwnerJoin{Links:links,Cases:cases}
 got,err:=j.ReadAuthorizedLinkedCases(context.Background(),"analyst","t1","e1","inc-a")
 if err!=nil || got==nil || len(got.Cases)!=0 || len(cases.calls)!=0 {
  t.Fatalf("authorized empty links mistaken for unavailable: %+v %v",got,err)
 }
 ctx,cancel:=context.WithCancel(context.Background());cancel()
 if out,err:=j.ReadAuthorizedLinkedCases(ctx,"analyst","t1","e1","inc-a");err==nil || out!=nil {
  t.Fatal("cancelled request should not query source owners")
 }
 if out,err:=j.ReadAuthorizedLinkedCases(context.Background(),"analyst","*","e1","inc-a");err==nil || out!=nil{
  t.Fatal("wildcard tenant must never reach independent owners")
 }
}
