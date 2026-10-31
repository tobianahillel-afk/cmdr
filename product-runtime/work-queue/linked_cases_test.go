package workqueue

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

type linkedCasesHost struct {
 *fakeOwner
 links *AuthorizedLinkedCases
 failure error
 calls []string
}
func (h *linkedCasesHost) ReadAuthorizedLinkedCases(_ context.Context, principal, tenant, environment, incidentID string) (*AuthorizedLinkedCases,error) {
 h.calls=append(h.calls,principal+":"+tenant+":"+environment+":"+incidentID)
 if h.failure!=nil { return nil,h.failure }
 return h.links,nil
}
func sampleLinks()*AuthorizedLinkedCases{
 return &AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{
  {CaseID:"case-001",Status:"investigating",LinkSource:"Object Linking Service",LinkVersion:"link-v3"},
  {CaseID:"case-002",Status:"review",LinkSource:"Object Linking Service",LinkVersion:"link-v4"},
 }}
}
func TestSelectedIncidentHasIndependentAuthorizedInvestigateCaseLinks(t *testing.T) {
 source:=&linkedCasesHost{fakeOwner:fixture(),links:sampleLinks()}
 h:=NewHTTPHandler(trustedHost(),source)
 r:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
 if r.Code!=http.StatusOK {t.Fatalf("selected Incident unavailable: %d %s",r.Code,r.Body.String())}
 var v Result
 if err:=json.Unmarshal(r.Body.Bytes(),&v);err!=nil {t.Fatal(err)}
 if v.Selected==nil || v.SelectedLinkedCases==nil ||
   len(v.SelectedLinkedCases.Cases)!=2 ||
   v.SelectedLinkedCases.Cases[1].CaseID!="case-002" {
  t.Fatalf("first-party linked Investigate projection missing: %+v",v)
 }
 if len(source.calls)!=1 || source.calls[0]!="analyst:t1:e1:inc-a" {
  t.Fatalf("wrong independent scope or duplicate host lookup: %v",source.calls)
 }
 source.calls=nil
 listing:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents")
 if listing.Code!=200 || len(source.calls)!=0 || strings.Contains(listing.Body.String(),"case-001") {
  t.Fatalf("unselected list disclosed Investigate Case: %d %s",listing.Code,listing.Body.String())
 }
 source.calls=nil
 task:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=tasks&selected=Task:task-b")
 if task.Code!=200 || len(source.calls)!=0 || strings.Contains(task.Body.String(),"case-001") {
  t.Fatalf("Task selection called Incident relation owner: %d %s",task.Code,task.Body.String())
 }
}
func TestLinkedInvestigateCasesEscapeInNativeHTMLAndDoNotInventMissingLinks(t *testing.T) {
 source:=&linkedCasesHost{fakeOwner:fixture(),links:&AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{
  {CaseID:"case-<script>alert(1)</script>",Status:"triage",LinkSource:"Service <img src=x onerror=alert(1)>"},
 }}}
 req:=httptest.NewRequest(http.MethodGet,"/command/work-queue?selected=Incident:inc-a",nil)
 req.Header.Set("Accept","text/html")
 w:=httptest.NewRecorder()
 NewHTTPHandler(trustedHost(),source).ServeHTTP(w,req)
 if w.Code!=200 || !strings.Contains(w.Body.String(),"&lt;script&gt;") ||
  strings.Contains(w.Body.String(),"<script>alert(1)</script>") ||
  !strings.Contains(w.Body.String(),"&lt;img") {
  t.Fatalf("linked Case source text was unsafe in browser: %d %s",w.Code,w.Body.String())
 }
 blank:=testQueueRequest(NewHTTPHandler(trustedHost(),fixture()),http.MethodGet,
  "/command/work-queue?selected=Incident:inc-a")
 if blank.Code!=200 || strings.Contains(blank.Body.String(),"selected_linked_cases") {
  t.Fatalf("absent link provider manufactured Case state: %d %q",blank.Code,blank.Body.String())
 }
}
func TestLinkedCaseSourceFailuresCannotDiscloseSelectedIncidentOrCase(t *testing.T) {
 for _,tc:=range []struct{name string; links *AuthorizedLinkedCases; failure error}{
  {"source refusal",nil,errors.New("secret-case-001")},
  {"missing provenance",&AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{{CaseID:"secret-case-001",Status:"review"}}},nil},
  {"invalid status",&AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{{CaseID:"secret-case-001",Status:"malicious",LinkSource:"owner"}}},nil},
  {"duplicate",&AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{{CaseID:"secret-case-001",Status:"review",LinkSource:"owner"},{CaseID:"secret-case-001",Status:"triage",LinkSource:"owner"}}},nil},
  {"controls",&AuthorizedLinkedCases{Cases:[]LinkedCaseSummary{{CaseID:"secret-case-001",Status:"review",LinkSource:"bad\nsource"}}},nil},
 } {
  t.Run(tc.name,func(t *testing.T){
   source:=&linkedCasesHost{fakeOwner:fixture(),links:tc.links,failure:tc.failure}
   r:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
    "/command/work-queue?selected=Incident:inc-a")
   if r.Code!=http.StatusServiceUnavailable || strings.Contains(r.Body.String(),"inc-a") ||
    strings.Contains(r.Body.String(),"secret") || strings.Contains(r.Body.String(),"case-001") {
    t.Fatalf("linked source failure leaked selected record or denied source: %d %q",r.Code,r.Body.String())
   }
  })
 }
}
func TestLinkedCaseHostNeverReceivesDeniedMissingOrForeignIncident(t *testing.T) {
 for _,selection:=range []string{"Incident:inc-hidden","Incident:foreign","Incident:missing","Task:task-no-sla"} {
  source:=&linkedCasesHost{fakeOwner:fixture(),links:sampleLinks()}
  r:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
   "/command/work-queue?selected="+selection)
  if selection=="Task:task-no-sla" {
   if r.Code!=200 || len(source.calls)!=0 {t.Fatalf("Task leaked Case access: %d %v",r.Code,source.calls)}
  } else if r.Code!=http.StatusNotFound || len(source.calls)!=0 {
   t.Fatalf("denied Incident reached external linked Case host: %s %d %v",selection,r.Code,source.calls)
  }
 }
}
