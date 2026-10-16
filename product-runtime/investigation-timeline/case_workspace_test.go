package investigationtimeline

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type integrationPrincipal struct{}
type catalogCase struct {
	tenant string
	environment string
	owner string
	ids []string
}
type catalogSource struct {
	owner string
	identity SourceOwnerIdentity
}

// integrationStore simulates one trusted server-side provider across a
// selected Case, canonical Timeline Entry snapshots and typed source ACLs.
// No caller-provided permission flags are consulted. A deployed service must
// replace this fixture with its real authenticated owner/ACL repositories.
type integrationStore struct {
	cases map[string]catalogCase
	entries map[string]Entry
	sources map[string]catalogSource
	caseChecks int
	catalogCalls int
	entryReads int
	sourceReads int
	failCatalog bool
	revokeAfterList bool
}
func workspaceFixture() (*integrationStore, context.Context) {
	store:=&integrationStore{
		cases:map[string]catalogCase{
			"case-one":{tenant:"tenant-a",environment:"prod",owner:"analyst-a",ids:[]string{"entry-new","entry-old","entry-hidden"}},
			"case-two":{tenant:"tenant-b",environment:"prod",owner:"analyst-b",ids:[]string{"entry-other"}},
		},
		entries:map[string]Entry{
			"entry-new":{ID:"entry-new",Kind:"observed",Tenant:"tenant-a",Timezone:"+00:00",SourceTime:"2026-10-01T11:00:00Z",SourceRef:"src-new",Provenance:"trusted:telemetry",Version:2,Supersedes:"entry-old"},
			"entry-old":{ID:"entry-old",Kind:"inferred",Tenant:"tenant-a",Timezone:"+00:00",SourceTime:"2026-10-01T10:00:00Z",SourceRef:"src-old",Provenance:"trusted:evidence",Version:1},
			"entry-hidden":{ID:"entry-hidden",Kind:"observed",Tenant:"tenant-a",Timezone:"+00:00",SourceTime:"2026-10-01T09:00:00Z",SourceRef:"src-hidden",Provenance:"restricted",Version:1},
			"entry-other":{ID:"entry-other",Kind:"observed",Tenant:"tenant-b",Timezone:"+00:00",SourceTime:"2026-10-01T08:00:00Z",SourceRef:"src-other",Provenance:"foreign",Version:1},
		},
		sources:map[string]catalogSource{
			"src-new":{owner:"analyst-a",identity:SourceOwnerIdentity{Reference:"src-new",Kind:"telemetry-event",ObjectID:"event-1",Tenant:"tenant-a",Environment:"prod",CaseID:"case-one"}},
			"src-old":{owner:"analyst-a",identity:SourceOwnerIdentity{Reference:"src-old",Kind:"evidence",ObjectID:"evidence-1",Tenant:"tenant-a",Environment:"prod",CaseID:"case-one"}},
			"src-hidden":{owner:"analyst-b",identity:SourceOwnerIdentity{Reference:"src-hidden",Kind:"evidence",ObjectID:"evidence-secret",Tenant:"tenant-a",Environment:"prod",CaseID:"case-one"}},
			"src-other":{owner:"analyst-b",identity:SourceOwnerIdentity{Reference:"src-other",Kind:"incident",ObjectID:"inc-1",Tenant:"tenant-b",Environment:"prod",CaseID:"case-two"}},
		},
	}
	return store,context.WithValue(context.Background(),integrationPrincipal{},"analyst-a")
}
func(s *integrationStore) allowed(ctx context.Context,r AccessRequest)bool {
	principal,_:=ctx.Value(integrationPrincipal{}).(string)
	c,exists:=s.cases[r.CaseID]
	return exists && principal!="" && c.owner==principal &&
		c.tenant==r.Tenant && c.environment==r.Environment
}
func(s *integrationStore) CheckAccess(ctx context.Context,r AccessRequest)(string,error){
	if ctx.Err()!=nil {return "",ctx.Err()}
	if !s.allowed(ctx,r) {return "deny",nil}
	switch r.Resource{
	case "case":
		s.caseChecks++
		if s.revokeAfterList && s.catalogCalls>0 {return "deny",nil}
		if r.ObjectID==r.CaseID && r.Permission==caseReadPermission {return "allow",nil}
	case "timeline-entry":
		if r.Permission!=timelineReadPermission {return "deny",nil}
		for _,id:=range s.cases[r.CaseID].ids {if id==r.ObjectID{return "allow",nil}}
	case "telemetry-event","evidence","incident","decision":
		permission,valid:=canonicalSourceReadPermission(r.Resource)
		if !valid||r.Permission!=permission{return "deny",nil}
		for _,src:=range s.sources{
			x:=src.identity
			if x.CaseID==r.CaseID && x.Tenant==r.Tenant &&
				x.Environment==r.Environment && x.Kind==r.Resource &&
				x.ObjectID==r.ObjectID {
				if src.owner=="analyst-a" {return "allow",nil}
				return "deny",nil
			}
		}
	}
	return "deny",nil
}
func(s *integrationStore) ListCaseEntryIDs(ctx context.Context,r AccessRequest)([]string,error){
	s.catalogCalls++
	if s.failCatalog {return nil,errors.New("catalog offline")}
	if !s.allowed(ctx,r)||r.Resource!="case"||r.Permission!=caseReadPermission{return nil,errors.New("case not found")}
	return append([]string(nil),s.cases[r.CaseID].ids...),nil
}
func(s *integrationStore) ResolveEntry(ctx context.Context,r AccessRequest)(Entry,error){
	s.entryReads++
	if !s.allowed(ctx,r) || r.Resource!="timeline-entry" {return Entry{},errors.New("entry not found")}
	entry,ok:=s.entries[r.ObjectID]
	if !ok{return Entry{},errors.New("entry not found")}
	return entry,nil
}
func(s *integrationStore) ResolveSource(ctx context.Context,r AccessRequest)(SourceOwnerIdentity,error){
	s.sourceReads++
	if !s.allowed(ctx,r)||r.Resource!="source-reference" {return SourceOwnerIdentity{},errors.New("source not found")}
	source,ok:=s.sources[r.ObjectID]
	if !ok {return SourceOwnerIdentity{},errors.New("source not found")}
	return source.identity,nil
}
func TestCaseTimelineIntegratedOwnerFlow(t *testing.T){
	s,ctx:=workspaceFixture()
	scope:=Scope{Tenant:"tenant-a",Environment:"prod",CasePermission:"deny",TimelinePermission:"unknown"}
	got:=OpenCaseTimeline(ctx,scope,"case-one",s,s,s,s)
	if got.Status!="filtered"||len(got.Entries)!=2 {t.Fatalf("not an authorized integrated projection: %+v",got)}
	wantIDs:=[]string{"entry-old","entry-new"}
	ids:=[]string{got.Entries[0].ID,got.Entries[1].ID}
	if !reflect.DeepEqual(ids,wantIDs)||got.Entries[0].Kind!="inferred"||got.Entries[1].Supersedes!="entry-old" {
		t.Fatalf("forensic ordering/provenance lost: %+v",got)
	}
	if got.Entries[0].Provenance!="trusted:evidence"||got.Entries[1].SourceRef!="src-new" {t.Fatal("owner metadata missing")}
	if s.caseChecks!=2||s.catalogCalls!=1||s.entryReads!=3||s.sourceReads!=3 {
		t.Fatalf("unexpected provider calls case=%d catalog=%d entry=%d source=%d",s.caseChecks,s.catalogCalls,s.entryReads,s.sourceReads)
	}
}
func TestCaseTimelineNeverListsUnauthorizedCases(t *testing.T){
	for _,tc:=range []struct{name,caseID,tenant,environment string}{
		{"other-case","case-two","tenant-a","prod"},
		{"foreign-tenant","case-one","tenant-b","prod"},
		{"foreign-env","case-one","tenant-a","stage"},
		{"invalid-case"," * ","tenant-a","prod"},
	}{
		t.Run(tc.name,func(t *testing.T){
			s,ctx:=workspaceFixture()
			got:=OpenCaseTimeline(ctx,Scope{Tenant:tc.tenant,Environment:tc.environment},tc.caseID,s,s,s,s)
			if got.Status!="deny"||len(got.Entries)!=0||s.catalogCalls!=0||s.entryReads!=0||s.sourceReads!=0 {
				t.Fatalf("unauthorized listing or disclosure: %+v calls=%d",got,s.catalogCalls)
			}
		})
	}
}
func TestCaseTimelineFailsOnCatalogErrorOrRevocation(t *testing.T){
	for _,property:=range []string{"catalog-failure","revoke","canceled"}{
		t.Run(property,func(t *testing.T){
			s,ctx:=workspaceFixture()
			if property=="catalog-failure"{s.failCatalog=true}
			if property=="revoke"{s.revokeAfterList=true}
			if property=="canceled"{var cancel context.CancelFunc;ctx,cancel=context.WithCancel(ctx);cancel()}
			got:=OpenCaseTimeline(ctx,Scope{Tenant:"tenant-a",Environment:"prod"},"case-one",s,s,s,s)
			if got.Status!="deny"||len(got.Entries)!=0||s.entryReads!=0||s.sourceReads!=0 {
				t.Fatalf("failure exposed entry/source metadata: %+v",got)
			}
		})
	}
}
func TestCaseTimelineRejectsMisconfiguredIntegrations(t *testing.T){
	s,ctx:=workspaceFixture()
	for _,tc:=range []struct{name string;check AccessChecker;catalog CaseTimelineCatalog;entries EntryResolver;sources SourceOwnerResolver}{
		{"no-checker",nil,s,s,s},{"no-catalog",s,nil,s,s},
		{"no-entries",s,s,nil,s},{"no-sources",s,s,s,nil},
	}{
		t.Run(tc.name,func(t *testing.T){
			got:=OpenCaseTimeline(ctx,Scope{Tenant:"tenant-a",Environment:"prod"},"case-one",tc.check,tc.catalog,tc.entries,tc.sources)
			if got.Status!="deny"||len(got.Entries)!=0 {t.Fatalf("missing authority accepted: %+v",got)}
		})
	}
}
