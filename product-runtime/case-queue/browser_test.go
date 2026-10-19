package casequeue

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestNativeQueueBrowserServesReadOnlyAssetsAndUsesCurrentSourceACL(t *testing.T) {
 allowed := nativeCase("case-visible", "Authorized investigation")
 denied := nativeCase("case-secret", "TOP-SECRET-CASE-CONTENT")
 session := queueSessionFunc(func(r *http.Request) (context.Context, QueueScope, error) {
  cookie, err := r.Cookie("host-session")
  if err != nil || cookie.Value != "authorized" {
   return nil, QueueScope{}, errors.New("not logged in")
  }
  return r.Context(), QueueScope{Tenant: "tenant-a", Environment: "prod"}, nil
 })
 provider := queueProviderFunc(func(context.Context, QueueScope) ([]CaseSnapshot, error) {
  return []CaseSnapshot{allowed,denied}, nil
 })
 checks := 0
 access := queueAccessFunc(func(_ context.Context, _ QueueScope, resource, id string) (string, error) {
  checks++
  if resource == "case" && id == "case-secret" { return "deny", nil }
  return "allow", nil
 })
 h, err := NewAuthorizedQueueBrowserHandler(session, provider, access, nil)
 if err != nil { t.Fatal(err) }
 for _, path := range []string{"/case-queue/ui", "/case-queue/ui.js", "/case-queue/ui.css"} {
  got := httptest.NewRecorder()
  h.ServeHTTP(got, httptest.NewRequest("GET", path,nil))
  if got.Code != http.StatusOK { t.Fatalf("browser asset %q: %d",path,got.Code) }
  if got.Header().Get("Cache-Control") != "no-store" ||
   !strings.Contains(got.Header().Get("Content-Security-Policy"),"default-src 'none'") {
   t.Fatalf("missing private browser CSP/cache protection on %q",path)
  }
  if strings.Contains(got.Body.String(),"TOP-SECRET-CASE-CONTENT") ||
   strings.Contains(got.Body.String(),"case-visible") {
   t.Fatalf("static asset contains operator Case data: %q",path)
  }
 }
 html := httptest.NewRecorder()
 h.ServeHTTP(html,httptest.NewRequest("GET","/case-queue/ui",nil))
 if !strings.Contains(html.Body.String(),"Case Queue") ||
  !strings.Contains(html.Body.String(),"/case-queue/ui.js") {
  t.Fatal("browser does not load native queue UI")
 }
 unauth := httptest.NewRecorder()
 h.ServeHTTP(unauth,httptest.NewRequest("GET","/case-queue?limit=25",nil))
 if unauth.Code != http.StatusUnauthorized || strings.Contains(unauth.Body.String(),"case-visible") {
  t.Fatalf("browser could access cases without trusted host session: %d",unauth.Code)
 }
 request := httptest.NewRequest("GET","/case-queue?limit=25",nil)
 request.AddCookie(&http.Cookie{Name:"host-session",Value:"authorized"})
 list := httptest.NewRecorder()
 h.ServeHTTP(list,request)
 if list.Code != http.StatusOK || !strings.Contains(list.Body.String(),"case-visible") ||
  strings.Contains(list.Body.String(),"case-secret") || strings.Contains(list.Body.String(),"TOP-SECRET") ||
  checks != 2 {
  t.Fatalf("browser API bypassed source owner authorization: %d %s checks=%d",list.Code,list.Body.String(),checks)
 }
}
func TestNativeQueueBrowserRejectsUnsafePathsAndMethodsWithoutProviders(t *testing.T) {
 calls := 0
 session := queueSessionFunc(func(r *http.Request) (context.Context,QueueScope,error) {
  calls++
  return r.Context(),QueueScope{Tenant:"tenant-a",Environment:"prod"},nil
 })
 provider := queueProviderFunc(func(context.Context,QueueScope)([]CaseSnapshot,error) {
  calls++;return nil,nil
 })
 access := queueAccessFunc(func(context.Context,QueueScope,string,string)(string,error) {
  calls++;return "allow",nil
 })
 h,err := NewAuthorizedQueueBrowserHandler(session,provider,access,nil)
 if err != nil {t.Fatal(err)}
 for _, tc := range []struct{ method,path string; want int }{
  {"POST","/case-queue/ui",http.StatusMethodNotAllowed},
  {"GET","/case-queue/ui?tenant=secret",http.StatusBadRequest},
  {"GET","/case-queue/ui/extra",http.StatusNotFound},
  {"GET","/case-queue/private",http.StatusNotFound},
  {"GET","/favicon.ico",http.StatusNotFound},
 } {
  got:=httptest.NewRecorder()
  h.ServeHTTP(got,httptest.NewRequest(tc.method,tc.path,nil))
  if got.Code!=tc.want || strings.Contains(got.Body.String(),"tenant-a") {
   t.Fatalf("unexpected browser route response %s %s => %d",tc.method,tc.path,got.Code)
  }
 }
 if calls!=0 {t.Fatal("UI assets/invalid routes invoked Case owner or session")}
}
func TestNativeQueueBrowserConstructorRequiresOwner(t *testing.T) {
 if _,err:=NewAuthorizedQueueBrowserHandler(nil,nil,nil,nil);err==nil {
  t.Fatal("browser without first-party Case/tenant/session ACL accepted")
 }
}
