package main

import (
	"context"
	"testing"
)

// Incident.read, Shared relation and Govern Result.read are three independent
// checks even when the same synthetic user can open the Incident web page.
func TestSyntheticGovernResultOwnerDoesNotInheritCommandIncidentRead(t *testing.T) {
	relations:=demoSharedLinks{}
	govern:=demoGovernResultOwner{}
	for _,tc:=range []struct {
		name,principal,tenant,environment,incident,result string
		allowed bool
	}{
		{"analyst","analyst",demoTenant,demoEnvironment,"incident-a","result-a",true},
		{"incident-only-commander","commander",demoTenant,demoEnvironment,"incident-a","result-a",false},
		{"other-tenant","analyst","other",demoEnvironment,"incident-a","result-a",false},
		{"other-environment","analyst",demoTenant,"prod","incident-a","result-a",false},
		{"other-incident","analyst",demoTenant,demoEnvironment,"other","result-a",false},
		{"other-result","analyst",demoTenant,demoEnvironment,"incident-a","other",false},
	} {
		t.Run(tc.name,func(t *testing.T){
			refs,err:=relations.ReadAuthorizedIncidentResultReferences(context.Background(),
				tc.principal,tc.tenant,tc.environment,tc.incident)
			if err!=nil {t.Fatal(err)}
			wantLink:=tc.principal=="analyst"&&tc.tenant==demoTenant&&
				tc.environment==demoEnvironment&&tc.incident=="incident-a"
			if (len(refs)==1)!=wantLink {t.Fatalf("link scope leaked: %+v",refs)}
			result,err:=govern.ReadAuthorizedGovernResult(context.Background(),
				tc.principal,tc.tenant,tc.environment,tc.result)
			wantGovern:=tc.principal=="analyst"&&tc.tenant==demoTenant&&
				tc.environment==demoEnvironment&&tc.result=="result-a"
			if err!=nil||(result!=nil)!=wantGovern {t.Fatalf("Govern read scope leaked: %+v, %v",result,err)}
			if tc.allowed&&(len(refs)!=1||refs[0].ResultID!=result.ID||
				result.RunID!="run-a"||result.Outcome!="verification-failed"||
				result.SourceRef!="synthetic-govern-result-owner"){
				t.Fatalf("authorized Result provenance is incomplete: %+v, %+v",refs,result)
			}
		})
	}
	ctx,cancel:=context.WithCancel(context.Background())
	cancel()
	if refs,err:=relations.ReadAuthorizedIncidentResultReferences(ctx,"analyst",
		demoTenant,demoEnvironment,"incident-a");err==nil||len(refs)!=0 {
		t.Fatalf("canceled relation revealed source: %+v, %v",refs,err)
	}
	if result,err:=govern.ReadAuthorizedGovernResult(ctx,"analyst",
		demoTenant,demoEnvironment,"result-a");err==nil||result!=nil {
		t.Fatalf("canceled Govern read revealed source: %+v, %v",result,err)
	}
}
