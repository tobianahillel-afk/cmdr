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

type factorsTestHost struct {
	*fakeOwner
	factors AuthorizedFactors
	factorErr error
	calls []string
}
func (h *factorsTestHost) ReadAuthorizedFactors(_ context.Context, principal, tenant, environment string, kind Kind, id string) (AuthorizedFactors,error) {
	h.calls=append(h.calls,principal+":"+tenant+":"+environment+":"+string(kind)+":"+id)
	if h.factorErr != nil {return AuthorizedFactors{},h.factorErr}
	return h.factors,nil
}
func TestSelectedCoordinationDimensionsKeepSixIndependentOwnerClaims(t *testing.T) {
	host:=&factorsTestHost{fakeOwner:fixture(),factors:AuthorizedFactors{
		Priority:&Dimension{Value:"low",Source:"Command",ObservedAt:"2026-10-19T08:00:00Z"},
		Severity:&Dimension{Value:"critical",Source:"Alert provider",ObservedAt:"2026-10-19T07:00:00Z"},
		Urgency:&Dimension{Value:"medium",Source:"Command"},
		Impact:&Dimension{Value:"limited",Source:"Service owner"},
		Confidence:&Dimension{Value:"0.91",Source:"Detection engine"},
		SLARisk:&Dimension{Value:"unknown",Source:"SLA policy"},
	}}
	h:=NewHTTPHandler(trustedHost(),host)
	selected:="/command/work-queue?view=incidents&selected=Incident:inc-a"
	r:=testQueueRequest(h,http.MethodGet,selected)
	if r.Code!=200 {t.Fatalf("authorized factor inspection failed: %d %s",r.Code,r.Body.String())}
	var result Result
	if err:=json.Unmarshal(r.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
	if result.Selected==nil || result.Selected.ID!="inc-a" || result.SelectedFactors==nil ||
		result.SelectedFactors.Priority==nil || result.SelectedFactors.Severity==nil ||
		result.SelectedFactors.Priority.Value!="low" || result.SelectedFactors.Severity.Value!="critical" ||
		result.SelectedFactors.Confidence.Value!="0.91" || result.SelectedFactors.SLARisk.Value!="unknown" {
		t.Fatalf("distinct source dimensions lost or priority inferred from severity: %+v",result)
	}
	if len(host.calls)!=1 || host.calls[0]!="analyst:t1:e1:Incident:inc-a" {
		t.Fatalf("factor source called without exact host-authorized scope: %v",host.calls)
	}
	for _,before:=range []string{"\"severity\"","\"confidence\"","\"sla_risk\""} {
		if !strings.Contains(r.Body.String(),before) {t.Fatalf("missing separate factor: %q",before)}
	}
	// Existing unselected list remains unchanged; no factor provider call or
	// hidden privileged source field is added to its JSON output.
	host.calls=nil
	list:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents")
	if list.Code!=200 || strings.Contains(list.Body.String(),"selected_factors") || len(host.calls)!=0 {
		t.Fatalf("non-selected queue leaked factor projections: %d %q",list.Code,list.Body.String())
	}
}

func TestFactorSourceDenialAndProviderErrorNeverRevealPartialSelectedRow(t *testing.T) {
	for _,failure:=range []struct{
		name string
		payload AuthorizedFactors
		err error
	}{
		{"factor source failed",AuthorizedFactors{},errors.New("foreign source id secret")},
		{"invalid returned source",AuthorizedFactors{Severity:&Dimension{Value:"critical",Source:"danger\nprivate"}},nil},
		{"invalid time",AuthorizedFactors{Priority:&Dimension{Value:"high",Source:"Command",ObservedAt:"not-time"}},nil},
		{"missing provenance",AuthorizedFactors{Impact:&Dimension{Value:"high",Source:""}},nil},
	} {
		t.Run(failure.name,func(t *testing.T){
			host:=&factorsTestHost{fakeOwner:fixture(),factors:failure.payload,factorErr:failure.err}
			rec:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
				"/command/work-queue?selected=Incident:inc-a")
			if rec.Code!=http.StatusServiceUnavailable ||
				strings.Contains(rec.Body.String(),"inc-a") ||
				strings.Contains(rec.Body.String(),"private") ||
				strings.Contains(rec.Body.String(),"critical") ||
				strings.Contains(rec.Body.String(),"foreign source") {
				t.Fatalf("factor denial disclosed partial selection or raw provider error: %d %q",rec.Code,rec.Body.String())
			}
		})
	}
	// Denied, foreign, off-page and non-selected objects must NEVER cause
	// a call to source-factor adapter or disclose any source metadata.
	for _,sel:=range []string{"Incident:inc-hidden","Incident:foreign","Task:missing"} {
		host:=&factorsTestHost{fakeOwner:fixture(),factors:AuthorizedFactors{
			Severity:&Dimension{Value:"private",Source:"unauthorized remote"},
		}}
		rec:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
			"/command/work-queue?selected="+sel)
		if rec.Code!=http.StatusNotFound || len(host.calls)!=0 ||
			strings.Contains(rec.Body.String(),"private") {
			t.Fatalf("foreign/denied factor owner queried: %s: %d %q",sel,rec.Code,rec.Body.String())
		}
	}
}

func TestNativeSelectedFactorsEscapedAndMissingFactorsExplicitlyUnknown(t *testing.T) {
	host:=&factorsTestHost{fakeOwner:fixture(),factors:AuthorizedFactors{
		Severity:&Dimension{Value:"high <script>alert(1)</script>",Source:"sensor & SOC"},
	}}
	request:=http.MethodGet
	rec:=testQueueRequestHTML(NewHTTPHandler(trustedHost(),host),request,
		"/command/work-queue?selected=Incident:inc-a")
	if rec.Code!=200 {t.Fatalf("selected HTML unavailable: %d %q",rec.Code,rec.Body.String())}
	page:=rec.Body.String()
	if strings.Contains(page,"<script>alert(1)</script>") ||
		!strings.Contains(page,"&lt;script&gt;") || !strings.Contains(page,"sensor &amp; SOC") ||
		!strings.Contains(page,"Dimensions de coordination") ||
		!strings.Contains(page,"Confiance</dt><dd>inconnue") ||
		!strings.Contains(page,"Priorité</dt><dd>inconnue") {
		t.Fatalf("HTML escaped presentation or unknown source treatment failed: %q",page)
	}
	// Existing hosts without CAP-CMD-104 adapter still show all six fields
	// as unknown (rather than inventing a priority or SLA state).
	legacy:=testQueueRequestHTML(NewHTTPHandler(trustedHost(),fixture()),request,
		"/command/work-queue?selected=Incident:inc-a")
	if legacy.Code!=200 || !strings.Contains(legacy.Body.String(),"Sévérité</dt><dd>inconnue") {
		t.Fatalf("backwards-compatible host without factors was misrepresented: %d",legacy.Code)
	}
}
func testQueueRequestHTML(h http.Handler, method,path string) *httptest.ResponseRecorder {
	req:=httptest.NewRequest(method,path,nil)
	req.Header.Set("Accept","text/html")
	rec:=httptest.NewRecorder()
	h.ServeHTTP(rec,req)
	return rec
}
