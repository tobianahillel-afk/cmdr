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

type canonicalOwnershipHost struct {
 *fakeOwner
 claimed *AuthorizedCanonicalOwnership
 failure error
 calls []string
}

func (h *canonicalOwnershipHost) ReadAuthorizedCanonicalOwnership(_ context.Context, principal, tenant, environment string, kind Kind, id string) (*AuthorizedCanonicalOwnership,error) {
 h.calls=append(h.calls,principal+":"+tenant+":"+environment+":"+string(kind)+":"+id)
 if h.failure!=nil {return nil,h.failure}
 return h.claimed,nil
}

func TestSelectedCanonicalOwnershipIsDistinctFromOperationalOwner(t *testing.T) {
 source:=&canonicalOwnershipHost{fakeOwner:fixture(),
 claimed:&AuthorizedCanonicalOwnership{
  ProductOwner:"Command Product Lead",RegisterSource:"Ownership Register",
  RegisterVersion:"registry-v3",IdentityStatus:"identity-unresolved",
 }}
 handler:=NewHTTPHandler(trustedHost(),source)
 selected:="/command/work-queue?view=incidents&selected=Incident:inc-a"
 resp:=testQueueRequest(handler,http.MethodGet,selected)
 if resp.Code!=200 {t.Fatalf("selected owner unavailable: %d %s",resp.Code,resp.Body.String())}
 var result Result
 if err:=json.Unmarshal(resp.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
 if result.Selected==nil || result.SelectedCanonicalOwnership==nil ||
  result.SelectedCanonicalOwnership.ProductOwner!="Command Product Lead" ||
  result.SelectedCanonicalOwnership.IdentityStatus!="identity-unresolved" ||
  result.SelectedCanonicalOwnership.RegisterVersion!="registry-v3" {
  t.Fatalf("canonical owner absent or conflated with assignment: %+v",result)
 }
 if len(source.calls)!=1 || source.calls[0]!="analyst:t1:e1:Incident:inc-a" {
  t.Fatalf("canonical source received wrong session/current-row scope: %v",source.calls)
 }
 // No source lookup or product ownership metadata for an unselected view.
 source.calls=nil
 list:=testQueueRequest(handler,http.MethodGet,"/command/work-queue?view=incidents")
 if list.Code!=200 || len(source.calls)!=0 ||
  strings.Contains(list.Body.String(),"selected_canonical_ownership") {
  t.Fatalf("unselected row leaked registry context: %d %s",list.Code,list.Body.String())
 }
 // Product owner is independent of the operational item responsible owner:
 // never transfer owner, assign staff, or resolve identity optimistically.
 if result.Selected.Owner=="Command Product Lead" {
  t.Fatal("fixture accidentally equates canonical registry with operational assignment")
 }
}

func TestCanonicalOwnershipViewEscapesUntrustedRegistryText(t *testing.T) {
 source:=&canonicalOwnershipHost{fakeOwner:fixture(),claimed:&AuthorizedCanonicalOwnership{
  ProductOwner:"Command <script>evil()</script>",
  RegisterSource:"Ownership Register",IdentityStatus:"resolved",
 }}
 r:=httptest.NewRequest(http.MethodGet,"/command/work-queue?selected=Incident:inc-a",nil)
 r.Header.Set("Accept","text/html")
 w:=httptest.NewRecorder()
 NewHTTPHandler(trustedHost(),source).ServeHTTP(w,r)
 if w.Code!=200 || !strings.Contains(w.Body.String(),"&lt;script&gt;") ||
  strings.Contains(w.Body.String(),"<script>evil()</script>") ||
  !strings.Contains(w.Body.String(),"Le responsable opérationnel") {
  t.Fatalf("selected product ownership HTML unsafe or unavailable: %d %s",w.Code,w.Body.String())
 }
}

func TestCanonicalOwnershipFailsClosedOnInvalidOrDeniedRegister(t *testing.T) {
 cases:=[]struct{
  name string
  value *AuthorizedCanonicalOwnership
  err error
 }{
  {"registry denies",nil,errors.New("secret foreign owner")},
  {"owner without source",&AuthorizedCanonicalOwnership{ProductOwner:"secret-owner"},nil},
  {"unknown state",&AuthorizedCanonicalOwnership{RegisterSource:"reg",IdentityStatus:"authorized-elevated"},nil},
  {"injected registry",&AuthorizedCanonicalOwnership{ProductOwner:"private",RegisterSource:"reg\nother"},nil},
  {"invalid owner",&AuthorizedCanonicalOwnership{ProductOwner:"private\r",RegisterSource:"reg"},nil},
  {"version without source",&AuthorizedCanonicalOwnership{RegisterVersion:"version-secret"},nil},
 }
 for _,tc:=range cases {
  t.Run(tc.name,func(t *testing.T) {
   source:=&canonicalOwnershipHost{fakeOwner:fixture(),claimed:tc.value,failure:tc.err}
   rec:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
    "/command/work-queue?selected=Incident:inc-a")
   if rec.Code!=http.StatusServiceUnavailable ||
    strings.Contains(rec.Body.String(),"inc-a") ||
    strings.Contains(rec.Body.String(),"private") ||
    strings.Contains(rec.Body.String(),"secret") {
    t.Fatalf("failed register leaked selected data: %d %q",rec.Code,rec.Body.String())
   }
  })
 }
}

func TestCanonicalOwnershipCannotReadDeniedOrNonselectedItems(t *testing.T) {
 for _,token:=range []string{"Incident:foreign","Incident:inc-hidden","Task:missing"} {
  source:=&canonicalOwnershipHost{fakeOwner:fixture(),claimed:&AuthorizedCanonicalOwnership{
   ProductOwner:"private",RegisterSource:"private-register",
  }}
  resp:=testQueueRequest(NewHTTPHandler(trustedHost(),source),http.MethodGet,
   "/command/work-queue?selected="+token)
  if resp.Code!=http.StatusNotFound || len(source.calls)!=0 {
   t.Fatalf("denied or absent item reached ownership registry: %s %v",token,source.calls)
  }
 }
}

func TestCanonicalOwnershipUnknownWhenRegisterNotConnected(t *testing.T) {
 handler:=NewHTTPHandler(trustedHost(),fixture())
 resp:=testQueueRequest(handler,http.MethodGet,"/command/work-queue?selected=Incident:inc-a")
 if resp.Code!=200 || strings.Contains(resp.Body.String(),"selected_canonical_ownership") {
  t.Fatalf("missing register manufactured canonical owner: %d %q",resp.Code,resp.Body.String())
 }
 r:=httptest.NewRequest(http.MethodGet,"/command/work-queue?selected=Incident:inc-a",nil)
 r.Header.Set("Accept","text/html")
 w:=httptest.NewRecorder()
 handler.ServeHTTP(w,r)
 if w.Code!=200 || !strings.Contains(w.Body.String(),"Produit propriétaire : inconnu") {
  t.Fatalf("missing register not marked unknown: %d %q",w.Code,w.Body.String())
 }
}
