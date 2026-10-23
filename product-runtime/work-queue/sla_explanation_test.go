package workqueue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSelectedSLAExplainsSourceCalendarAndHistoricalPauseReadOnly(t *testing.T) {
	// Source metadata belongs to its own independently authorized SLA host,
	// and is visible ONLY after fresh Work Queue item authorization.
	reason := "Change window <script>globalThis.pwn=1</script>"
	host := &slaTestHost{fakeOwner:fixture(), sla:&AuthorizedSLA{
		State:"paused", PolicySource:"tenant-SLA-authority",
		PolicyVersion:"policy-v7", ObservedAt:"2026-10-20T11:30:00Z",
		DueAt:"2026-10-21T11:30:00Z", CalendarSource:"24x7-contract-calendar",
		PausedAt:"2026-10-20T11:25:00Z", PauseReason:reason,
	}}
	handler:=NewHTTPHandler(trustedHost(),host)
	request:=httptest.NewRequest(http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a",nil)
	request.Header.Set("Accept","text/html")
	html:=httptest.NewRecorder()
	handler.ServeHTTP(html,request)
	if html.Code!=http.StatusOK || !strings.Contains(html.Body.String(),"24x7-contract-calendar") ||
		!strings.Contains(html.Body.String(),"2026-10-20T11:25:00Z") ||
		!strings.Contains(html.Body.String(),"Change window") ||
		strings.Contains(html.Body.String(),"<script>globalThis.pwn=1</script>") {
		t.Fatalf("independently authorized pause metadata unavailable or executable: %d %s",html.Code,html.Body.String())
	}
	if !strings.Contains(html.Body.String(),"&lt;script&gt;") {
		t.Fatalf("source-supplied pause explanation not context-escaped: %q",html.Body.String())
	}
	if len(host.calls)!=1 || host.calls[0]!="analyst:t1:e1:Incident:inc-a" {
		t.Fatalf("unexpected SLA source lookup scope: %+v",host.calls)
	}
	var value struct{ SelectedSLA *AuthorizedSLA `json:"selected_sla"` }
	resp:=testQueueRequest(handler,http.MethodGet,"/command/work-queue?view=incidents&selected=Incident:inc-a")
	if resp.Code!=http.StatusOK { t.Fatal(resp.Body.String()) }
	if err:=json.Unmarshal(resp.Body.Bytes(),&value);err!=nil {t.Fatal(err)}
	if value.SelectedSLA==nil || value.SelectedSLA.PausedAt!="2026-10-20T11:25:00Z" ||
		value.SelectedSLA.PauseReason!=reason ||
		value.SelectedSLA.CalendarSource!="24x7-contract-calendar" {
		t.Fatalf("JSON paused source provenance unexpectedly lost: %+v",value.SelectedSLA)
	}
	host.calls=nil
	unselected:=testQueueRequest(handler,http.MethodGet,"/command/work-queue?view=incidents")
	if unselected.Code!=http.StatusOK || len(host.calls)!=0 ||
		strings.Contains(unselected.Body.String(),"pause_reason") ||
		strings.Contains(unselected.Body.String(),"tenant-SLA-authority") {
		t.Fatalf("unselected queue exposes privileged SLA: %d %s",unselected.Code,unselected.Body.String())
	}
}

func TestSLAExplainPausesWithUnavailableSourceFieldsWithoutInventingEvidence(t *testing.T) {
	host:=&slaTestHost{fakeOwner:fixture(),sla:&AuthorizedSLA{
		State:"paused",PolicySource:"independent-authority",PolicyVersion:"v3",
		ObservedAt:"2026-10-20T11:30:00Z",
	}}
	rec:=httptest.NewRecorder()
	request:=httptest.NewRequest(http.MethodGet,"/command/work-queue?selected=Incident:inc-a",nil)
	request.Header.Set("Accept","text/html")
	NewHTTPHandler(trustedHost(),host).ServeHTTP(rec,request)
	if rec.Code!=http.StatusOK || !strings.Contains(rec.Body.String(),"horodatage inconnu") ||
		!strings.Contains(rec.Body.String(),"raison : non fournie") ||
		!strings.Contains(rec.Body.String(),"calendrier : non fourni") {
		t.Fatalf("incomplete but authorized pause falsely fabricated details: %d %s",rec.Code,rec.Body.String())
	}
}

func TestSLAPauseClaimsRejectContradictoryOrUnsafeSourceValues(t *testing.T) {
	base:=AuthorizedSLA{State:"paused",PolicySource:"tenant-policy",
		PolicyVersion:"v3",ObservedAt:"2026-10-20T11:30:00Z",
		PausedAt:"2026-10-20T11:25:00Z",PauseReason:"change window"}
	tests:=[]struct{name string;change func(*AuthorizedSLA)}{
		{"pause after observed",func(a *AuthorizedSLA){a.PausedAt="2026-10-20T12:00:00Z"}},
		{"bad pause timestamp",func(a *AuthorizedSLA){a.PausedAt="Monday"}},
		{"pause on healthy",func(a *AuthorizedSLA){a.State="healthy"}},
		{"pause on unknown",func(a *AuthorizedSLA){a.State="unknown"}},
		{"reason without paused state",func(a *AuthorizedSLA){a.State="breached"}},
		{"source calendar control",func(a *AuthorizedSLA){a.CalendarSource="raw\nprivate-source"}},
		{"pause reason control",func(a *AuthorizedSLA){a.PauseReason="explanation\rsecret"}},
		{"oversized pause reason",func(a *AuthorizedSLA){a.PauseReason=strings.Repeat("x",257)}},
	}
	for _,tc:=range tests{
		t.Run(tc.name,func(t *testing.T){
			payload:=base
			tc.change(&payload)
			host:=&slaTestHost{fakeOwner:fixture(),sla:&payload}
			response:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
				"/command/work-queue?selected=Incident:inc-a")
			if response.Code!=http.StatusServiceUnavailable ||
				strings.Contains(response.Body.String(),"inc-a") ||
				strings.Contains(response.Body.String(),"tenant-policy") ||
				strings.Contains(response.Body.String(),"private-source") {
				t.Fatalf("invalid source-owned SLA leaked selected data: %d %q",response.Code,response.Body.String())
			}
		})
	}
	for _,selected:=range []string{"Incident:foreign","Incident:inc-hidden","Task:missing"} {
		host:=&slaTestHost{fakeOwner:fixture(),sla:&base}
		response:=testQueueRequest(NewHTTPHandler(trustedHost(),host),http.MethodGet,
			"/command/work-queue?selected="+selected)
		if response.Code!=http.StatusNotFound || len(host.calls)!=0 {
			t.Fatalf("inaccessible item caused SLA source lookup: %q %+v",selected,host.calls)
		}
	}
}
