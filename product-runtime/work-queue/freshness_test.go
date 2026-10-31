package workqueue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func freshAt(delta time.Duration) string {
	return time.Now().UTC().Add(delta).Format(time.RFC3339Nano)
}

func TestFreshnessCanonicalStatesNeedExplicitOwnerPolicy(t *testing.T) {
	for _, tc := range []struct {
		name string
		source *SourceFreshness
		want string
	}{
		{"no-metadata", nil, ""},
		{"no-policy", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-time.Hour), Health:"available"}, "unknown"},
		{"no-time", &SourceFreshness{Source:"sensor-A", Health:"available", PolicyVersion:"tenant-rule-v1", AgingAfterSeconds:60, StaleAfterSeconds:120}, "unknown"},
		{"future-clock", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(time.Hour), Health:"available", PolicyVersion:"tenant-rule-v1", AgingAfterSeconds:60, StaleAfterSeconds:120}, "unknown"},
		{"current", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-5*time.Second), Health:"available", PolicyVersion:"tenant-rule-v1", AgingAfterSeconds:60, StaleAfterSeconds:120}, "current"},
		{"aging", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-90*time.Second), Health:"available", PolicyVersion:"tenant-rule-v1", AgingAfterSeconds:60, StaleAfterSeconds:120}, "aging"},
		{"stale", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-3*time.Minute), Health:"available", PolicyVersion:"tenant-rule-v1", AgingAfterSeconds:60, StaleAfterSeconds:120}, "stale"},
		{"down", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-3*time.Hour), Health:"unavailable"}, "source-unavailable"},
		{"partial", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-3*time.Hour), Health:"partial"}, "partial"},
		{"unversioned-threshold", &SourceFreshness{Source:"sensor-A", ObservedAt:freshAt(-time.Hour), Health:"available", AgingAfterSeconds:60, StaleAfterSeconds:120}, "unknown"},
	} {
		t.Run(tc.name, func(t *testing.T){
			owner:=fixture()
			owner.items[0].SourceFreshness=tc.source
			result,err:=Project(context.Background(),owner,request(Incidents))
			if err!=nil || len(result.Rows)!=1 {t.Fatalf("authorized queue failed: %v %+v",err,result)}
			got:=result.Rows[0].Freshness
			if tc.source==nil {
				if got!=nil {t.Fatalf("legacy unconfigured metadata changed: %+v",got)}
				return
			}
			if got==nil || got.State!=tc.want || got.Source!="sensor-A" {
				t.Fatalf("wrong source freshness: want %s got %+v",tc.want,got)
			}
			if tc.source.PolicyVersion!="" && got.PolicyVersion!=tc.source.PolicyVersion {
				t.Fatalf("versioned policy provenance lost: %+v",got)
			}
			if tc.want=="partial" || tc.want=="source-unavailable" {
				if result.State!="partial-sources" || len(result.Diagnostics)!=1 ||
					result.Diagnostics[0]!="work-source-degraded" {
					t.Fatalf("source degradation hidden or misclassified: %+v",result)
				}
			} else if result.State!="ready" {
				t.Fatalf("unknown/stale do not indicate source outage: %+v",result)
			}
		})
	}
}

func TestFreshnessPermissionIsolationAndFilteredSourceHealth(t *testing.T) {
	owner:=fixture()
	owner.items[0].SourceFreshness=&SourceFreshness{Source:"sensor-A",ObservedAt:freshAt(-time.Minute),Health:"available"}
	owner.items[1].SourceFreshness=&SourceFreshness{Source:"sensor-B",Health:"unavailable"}
	owner.items[2].SourceFreshness=&SourceFreshness{Source:"TOP-SECRET-SENSOR",Health:"unavailable"}
	owner.items[4].SourceFreshness=&SourceFreshness{Source:"CROSS-TENANT-SENSOR",Health:"unavailable"}
	result,err:=Project(context.Background(),owner,request(Incidents))
	if err!=nil || len(result.Rows)!=1 || result.State!="ready" {
		t.Fatalf("other-view degradation leaked into selected Incident view: %v %+v",err,result)
	}
	if result.Rows[0].Freshness==nil || result.Rows[0].Freshness.State!="unknown" {
		t.Fatalf("missing policy not declared unknown: %+v",result.Rows[0])
	}
	all,err:=Project(context.Background(),owner,request(All))
	if err!=nil || all.State!="partial-sources" {t.Fatalf("visible Task source outage invisible: %v %+v",err,all)}
	raw,err:=json.Marshal(all)
	if err!=nil {t.Fatal(err)}
	for _,secret:=range []string{"TOP-SECRET-SENSOR","CROSS-TENANT-SENSOR"} {
		if strings.Contains(string(raw),secret) {t.Fatalf("denied source leaked into output: %s",secret)}
	}
	// Filtering away the degraded Task must also filter away its diagnostics.
	statusReq:=request(All)
	statusReq.StatusFilter="active"
	filtered,err:=Project(context.Background(),owner,statusReq)
	if err!=nil || filtered.State!="ready" || len(filtered.Diagnostics)!=0 {
		t.Fatalf("inaccessible to this filter source degraded current page: %v %+v",err,filtered)
	}
}

func TestMalformedFreshnessOwnerFailsClosedWithoutPartialRows(t *testing.T) {
	bad:=[]*SourceFreshness{
		{Source:"", Health:"available"},
		{Source:"\ninvalid", Health:"available"},
		{Source:"sensor-A", Health:"unknown-new-state"},
		{Source:"sensor-A", ObservedAt:"not-a-date",Health:"available"},
		{Source:"sensor-A", PolicyVersion:"bad\nversion",Health:"available"},
	}
	for _,meta:=range bad {
		owner:=fixture()
		owner.items[0].SourceFreshness=meta
		denied,err:=Project(context.Background(),owner,request(All))
		if err!=ErrIncomplete || len(denied.Rows)!=0 {
			t.Fatalf("invalid owner freshness returned data: %v %+v",err,denied)
		}
	}
}

func TestFreshnessHTMLAndJSONUseSameAuthorizedProjection(t *testing.T) {
	owner:=fixture()
	owner.items[0].SourceFreshness=&SourceFreshness{
		Source:"<img src=x onerror=alert(1)>",ObservedAt:freshAt(-3*time.Minute),
		Health:"available",PolicyVersion:"tenant-policy-v42",
		AgingAfterSeconds:60,StaleAfterSeconds:120,
	}
	handler:=NewHTTPHandler(trustedHost(),owner)
	url:="/command/work-queue?view=incidents&selected=Incident:inc-a"
	req:=httptest.NewRequest(http.MethodGet,url,nil)
	req.Header.Set("Accept","text/html")
	page:=httptest.NewRecorder()
	handler.ServeHTTP(page,req)
	html:=page.Body.String()
	if page.Code!=200 || !strings.Contains(html,"Fraîcheur source") ||
		!strings.Contains(html,"tenant-policy-v42") ||
		!strings.Contains(html,"stale") ||
		!strings.Contains(html,"&lt;img") ||
		strings.Contains(html,"<img src=") {
		t.Fatalf("freshness HTML not authorized/escaped: %d %q",page.Code,html)
	}
	if page.Header().Get("Cache-Control")!="no-store" {t.Fatal("freshness cached to browser")}
	req=httptest.NewRequest(http.MethodGet,url,nil)
	page=httptest.NewRecorder()
	handler.ServeHTTP(page,req)
	if page.Code!=200 {t.Fatalf("JSON viewer failed %d",page.Code)}
	var response Result
	if err:=json.Unmarshal(page.Body.Bytes(),&response);err!=nil {t.Fatal(err)}
	if response.Selected==nil || response.Selected.Freshness==nil ||
		response.Selected.Freshness.State!="stale" ||
		len(response.Rows)!=1 || response.Rows[0].Freshness==nil {
		t.Fatalf("server owner JSON differs from browser: %+v",response)
	}
	legacy:=httptest.NewRecorder()
	NewHTTPHandler(trustedHost(),fixture()).ServeHTTP(legacy,httptest.NewRequest(http.MethodGet,"/command/work-queue?view=incidents",nil))
	if legacy.Code!=200 || strings.Contains(legacy.Body.String(),"\"freshness\"") {
		t.Fatalf("legacy metadata-free JSON contract changed: %d %q",legacy.Code,legacy.Body.String())
	}
}
