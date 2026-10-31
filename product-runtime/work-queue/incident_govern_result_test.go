package workqueue

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testResultRelations struct {
	refs []IncidentResultReference
	err error
	calls int
}
func (o *testResultRelations) ReadAuthorizedIncidentResultReferences(_ context.Context, principal, tenant, environment, incidentID string) ([]IncidentResultReference,error) {
	o.calls++
	if principal != "analyst" || tenant != "t1" || environment != "e1" || incidentID != "inc-100" {
		return nil, errors.New("PRIVATE_LINK_DENIAL")
	}
	return o.refs,o.err
}
type testGovernResults struct {
	results map[string]*AuthorizedGovernResult
	err error
	calls int
}
func (o *testGovernResults) ReadAuthorizedGovernResult(_ context.Context, principal, tenant, environment, resultID string) (*AuthorizedGovernResult,error) {
	o.calls++
	if principal != "analyst" || tenant != "t1" || environment != "e1" {
		return nil, errors.New("PRIVATE_GOVERN_DENIAL")
	}
	return o.results[resultID],o.err
}
func resultFixture() (*fakeOwner,*testResultRelations,*testGovernResults) {
	owner:=deepReturnOwner()
	links:=&testResultRelations{refs: []IncidentResultReference{
		{IncidentID:"inc-100",ResultID:"result-safe",Tenant:"t1",Environment:"e1"},
		{IncidentID:"inc-100",ResultID:"result-denied",Tenant:"t1",Environment:"e1"},
	}}
	gov:=&testGovernResults{results:map[string]*AuthorizedGovernResult{
		"result-safe": {
			ID:"result-safe",Tenant:"t1",Environment:"e1",RunID:"run-3",
			Lifecycle:"finalized",Outcome:"verification-failed",
			Verification:"failed",ResidualRisk:"<script>residual-risk</script>",
			SourceRef:"govern-authorized-projection",SourceVersion:"v2",
			CreatedAt:"2026-10-21T02:10:00Z",
		},
	}}
	return owner,links,gov
}
func resultDetail(owner Owner,links IncidentResultReferencesOwner,gov GovernResultOwner) *httptest.ResponseRecorder {
	return detailGET(NewIncidentDetailHTTPHandlerWithResults(trustedHost(),owner,links,gov),
		"/command/incident-detail?incident=inc-100")
}
func TestIncidentGovernResultRealDetailKeepsSeparateSourceOwnerBoundary(t *testing.T) {
	owner,links,gov:=resultFixture()
	r:=resultDetail(owner,links,gov)
	body:=r.Body.String()
	if r.Code != http.StatusOK || links.calls!=1 || gov.calls!=2 {
		t.Fatalf("separate owner calls not observed: status=%d relation=%d govern=%d",r.Code,links.calls,gov.calls)
	}
	for _,s:=range []string{"Result result-safe","run-3","finalized","verification-failed","failed","&lt;script&gt;","govern-authorized-projection","v2"} {
		if !strings.Contains(body,s) { t.Fatalf("authorized Govern source omitted %q",s) }
	}
	for _,s:=range []string{"result-denied","<script>","PRIVATE_","success confirmé"} {
		if strings.Contains(body,s) {t.Fatalf("denied/raw Govern information exposed %q",s)}
	}
	owner.allowed["Incident:inc-100"]=false
	denied:=resultDetail(owner,links,gov)
	if denied.Code!=http.StatusNotFound || links.calls!=1 || gov.calls!=2 {
		t.Fatalf("revoked Incident reached Govern source: %d",denied.Code)
	}
}
func TestIncidentGovernResultDeniesAnyUntrustedRelationOrProjectionBeforeHTML(t *testing.T) {
	cases:=[]struct{name string; mutate func(*testResultRelations,*testGovernResults)}{
		{"relation-owner-error",func(l *testResultRelations,_ *testGovernResults){l.err=errors.New("PRIVATE_RELATION")}},
		{"too-many-links",func(l *testResultRelations,_ *testGovernResults){l.refs=make([]IncidentResultReference,9)}},
		{"duplicate-link",func(l *testResultRelations,_ *testGovernResults){l.refs[1].ResultID="result-safe"}},
		{"wrong-incident",func(l *testResultRelations,_ *testGovernResults){l.refs[0].IncidentID="hidden"}},
		{"wrong-tenant",func(l *testResultRelations,_ *testGovernResults){l.refs[0].Tenant="other"}},
		{"wrong-environment",func(l *testResultRelations,_ *testGovernResults){l.refs[0].Environment="other"}},
		{"hostile-result-id",func(l *testResultRelations,_ *testGovernResults){l.refs[0].ResultID="bad\nresult"}},
		{"govern-owner-error",func(_ *testResultRelations,g *testGovernResults){g.err=errors.New("PRIVATE_GOVERN")}},
		{"wrong-source-result",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].ID="foreign"}},
		{"wrong-source-tenant",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].Tenant="foreign"}},
		{"wrong-source-env",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].Environment="foreign"}},
		{"missing-run",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].RunID=""}},
		{"invalid-run",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].RunID="run\rforged"}},
		{"source-wildcard",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].SourceRef="*"}},
		{"source-invalid-version",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].SourceVersion="v\n2"}},
		{"source-invalid-text",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].Outcome=strings.Repeat("x",65)}},
		{"source-control",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].ResidualRisk="r\ninternal"}},
		{"source-invalid-time",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].CreatedAt="tomorrow"}},
		{"source-zero-time",func(_ *testResultRelations,g *testGovernResults){g.results["result-safe"].CreatedAt="0001-01-01T00:00:00Z"}},
	}
	for _,tc:=range cases {
		t.Run(tc.name,func(t *testing.T){
			owner,links,gov:=resultFixture()
			tc.mutate(links,gov)
			res:=resultDetail(owner,links,gov)
			if res.Code!=http.StatusServiceUnavailable ||
				strings.Contains(res.Body.String(),"Incident Detail — aperçu") ||
				strings.Contains(res.Body.String(),"result-safe") ||
				strings.Contains(res.Body.String(),"PRIVATE_") {
				t.Fatalf("unsafe partial Govern result returned: %d %s",res.Code,res.Body.String())
			}
		})
	}
}
func TestIncidentGovernResultLegacyAndContextFailureAreExplicitPartialOrDenied(t *testing.T) {
	owner,links,gov:=resultFixture()
	legacy:=resultDetail(owner,nil,nil)
	if legacy.Code!=http.StatusOK || !strings.Contains(legacy.Body.String(),"Result Govern non disponible") ||
		links.calls!=0 || gov.calls!=0 { t.Fatalf("legacy owner regression: %d",legacy.Code) }
	oneMissing:=resultDetail(owner,links,nil)
	if oneMissing.Code!=http.StatusOK || links.calls!=0 {t.Fatal("missing Govern owner triggered relation lookup")}
	ctx,cancel:=context.WithCancel(context.Background())
	cancel()
	out,err:=readAuthorizedIncidentGovernResults(ctx,links,gov,trustedHost().scope,"inc-100")
	if err==nil||len(out)!=0||links.calls!=0 {t.Fatalf("canceled source traversed relation: %v, %+v",err,out)}
	// A missing authorized Result record does not reveal even its stable ID.
	empty:=&testGovernResults{results:map[string]*AuthorizedGovernResult{}}
	r:=resultDetail(owner,links,empty)
	if r.Code!=http.StatusOK || strings.Contains(r.Body.String(),"result-denied") ||
		strings.Contains(r.Body.String(),"result-safe") ||
		!strings.Contains(r.Body.String(),"Result Govern non disponible") {
		t.Fatalf("denied/unknown result leaked: %d %s",r.Code,r.Body.String())
	}
}
