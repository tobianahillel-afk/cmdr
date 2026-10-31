package workqueue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type slaTestHost struct {
	*fakeOwner
	sla *AuthorizedSLA
	slaErr error
	calls []string
}
func (h *slaTestHost) ReadAuthorizedSLA(_ context.Context, principal,tenant,environment string,kind Kind,id string) (*AuthorizedSLA,error) {
	h.calls=append(h.calls,principal+":"+tenant+":"+environment+":"+string(kind)+":"+id)
	return h.sla,h.slaErr
}

func TestSelectedAuthorizedSLAProvidesPolicyStateAndDueWithoutInventedRisk(t *testing.T) {
	remaining:=int64(0)
	host:=&slaTestHost{fakeOwner:fixture(),sla:&AuthorizedSLA{
		State:"at-risk",PolicySource:"Command SLA policy",PolicyVersion:"v2",
		ObservedAt:"2026-10-20T10:00:00Z",DueAt:"2026-10-20T11:00:00Z",
		RemainingSeconds:&remaining,
	}}
	h:=NewHTTPHandler(trustedHost(),host)
	out:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
	if out.Code!=200 {t.Fatalf("valid SLA inspection failed: %d %s",out.Code,out.Body.String())}
	var result Result
	if err:=json.Unmarshal(out.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
	if result.Selected==nil || result.Selected.ID!="inc-a" ||
		result.SelectedSLA==nil || result.SelectedSLA.State!="at-risk" ||
		result.SelectedSLA.PolicyVersion!="v2" ||
		result.SelectedSLA.RemainingSeconds==nil || *result.SelectedSLA.RemainingSeconds!=0 {
		t.Fatalf("source-owned SLA state/provenance unavailable: %+v",result)
	}
	if len(host.calls)!=1 || host.calls[0]!="analyst:t1:e1:Incident:inc-a" {
		t.Fatalf("SLA owner not called in exact authorized scope: %+v",host.calls)
	}
	// No client-side SLA value or policy is accepted; the selected source
	// lookup never runs for the ordinary unselected six-view queue.
	host.calls=nil
	list:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents")
	if list.Code!=200 || len(host.calls)!=0 || strings.Contains(list.Body.String(),"selected_sla") ||
		strings.Contains(list.Body.String(),"Command SLA policy") {
		t.Fatalf("unselected queue accessed privileged SLA: %d %q",list.Code,list.Body.String())
	}
	html:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
	if html.Code!=200 {t.Fatal(html.Body.String())}
	// A real HTML representation of the same selected object is separately
	// checked; it remains the same single Work Queue route and native template.
}

func TestSelectedSLADenialFailsClosedWithoutSelectedRowOrPolicy(t *testing.T) {
	cases:=[]struct{name string;value *AuthorizedSLA;err error}{
		{"provider unavailable",nil,errors.New("external customer secret")},
		{"invalid state",&AuthorizedSLA{State:"force-healthy",PolicySource:"private"},nil},
		{"missing policy evidence",&AuthorizedSLA{State:"healthy",PolicySource:"Provider"},nil},
		{"malformed due",&AuthorizedSLA{State:"at-risk",PolicySource:"Provider",PolicyVersion:"v1",
			ObservedAt:"2026-10-20T10:00:00Z",DueAt:"bad-time"},nil},
		{"newline injection",&AuthorizedSLA{State:"at-risk",PolicySource:"Provider\nprivate",
			PolicyVersion:"v1",ObservedAt:"2026-10-20T10:00:00Z"},nil},
	}
	for _,tt:=range cases {
		t.Run(tt.name,func(t *testing.T){
			host:=&slaTestHost{fakeOwner:fixture(),sla:tt.value,slaErr:tt.err}
			rec:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
				"/command/work-queue?selected=Incident:inc-a")
			if rec.Code!=http.StatusServiceUnavailable ||
				strings.Contains(rec.Body.String(),"inc-a") ||
				strings.Contains(rec.Body.String(),"Provider") ||
				strings.Contains(rec.Body.String(),"customer secret") {
				t.Fatalf("invalid SLA disclosed owner content: %d %q",rec.Code,rec.Body.String())
			}
		})
	}
	for _,selected:=range []string{"Incident:inc-hidden","Incident:foreign","Task:missing"} {
		host:=&slaTestHost{fakeOwner:fixture(),sla:&AuthorizedSLA{State:"at-risk",PolicySource:"sensitive"}}
		rec:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
			"/command/work-queue?selected="+selected)
		if rec.Code!=http.StatusNotFound || len(host.calls)!=0 {
			t.Fatalf("forbidden selection queried SLA owner: %q %d %+v",selected,rec.Code,host.calls)
		}
	}
}

func TestMissingSLAPolicyRemainsUnknownAndNoNewServerRights(t *testing.T) {
	h:=NewHTTPHandler(trustedHost(),fixture()) // No SLAOwner available
	res:=testQueueRequest(h,http.MethodGet,"/command/work-queue?selected=Incident:inc-a")
	if res.Code!=200 {t.Fatal(res.Body.String())}
	if strings.Contains(res.Body.String(),"selected_sla") || strings.Contains(res.Body.String(),"remaining_seconds") {
		t.Fatalf("missing SLA source fabricated a state: %q",res.Body.String())
	}
	host:=&slaTestHost{fakeOwner:fixture()}
	res=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
		"/command/work-queue?selected=Incident:inc-a")
	if res.Code!=200 || strings.Contains(res.Body.String(),"selected_sla") {
		t.Fatalf("nil source SLA leaked or invented claim: %d %q",res.Code,res.Body.String())
	}
	unknown:=&AuthorizedSLA{State:"unknown"}
	if _,err:=validateAuthorizedSLA(unknown);err!=nil {t.Fatalf("unknown must be representable: %v",err)}
	remaining:=int64(11)
	if _,err:=validateAuthorizedSLA(&AuthorizedSLA{State:"unknown",RemainingSeconds:&remaining});err==nil {
		t.Fatal("unknown SLA cannot contain invented remaining time")
	}
}
