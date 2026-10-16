package investigationtimeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type contractFixtures struct {
	Fixtures []struct {
		ID string `json:"id"`
		Input Input `json:"input"`
		Expected struct {
			Status string `json:"status"`
			VisibleIDs []string `json:"visible_ids"`
			VisibleKinds []string `json:"visible_kinds"`
		} `json:"expected"`
	} `json:"fixtures"`
}

func TestProjectMatchesContractFixtures(t *testing.T) {
	data,err:=os.ReadFile(filepath.Join("..","..","engineering","implementation","investigation-timeline","fixtures.json"))
	if err!=nil { t.Fatal(err) }
	var fixtures contractFixtures
	if err:=json.Unmarshal(data,&fixtures);err!=nil { t.Fatal(err) }
	if len(fixtures.Fixtures)!=18 { t.Fatalf("expected 18 contract cases, got %d",len(fixtures.Fixtures)) }
	for _,f:=range fixtures.Fixtures {
		t.Run(f.ID,func(t *testing.T){
			got:=Project(f.Input)
			if got.Status!=f.Expected.Status {t.Fatalf("status got %s want %s",got.Status,f.Expected.Status)}
			ids,kinds:=[]string{},[]string{}
			for _,e:=range got.Entries {
				ids=append(ids,e.ID)
				kinds=append(kinds,e.Kind)
			}
			if !reflect.DeepEqual(ids,f.Expected.VisibleIDs) || !reflect.DeepEqual(kinds,f.Expected.VisibleKinds) {
				t.Fatalf("projection got ids=%v kinds=%v, expected ids=%v kinds=%v",ids,kinds,f.Expected.VisibleIDs,f.Expected.VisibleKinds)
			}
		})
	}
}

func allowedFixture() Input {
	return Input{Scope:Scope{Tenant:"tenant-a",Environment:"prod",CasePermission:"allow",TimelinePermission:"allow"},
		Entries:[]Entry{{ID:"x",Kind:"observed",Tenant:"tenant-a",SourcePermission:"allow",
		Timezone:"+00:00",SourceTime:"2026-10-01T10:00:00Z",SourceRef:"owner:x",
		Provenance:"telemetry:x",Version:1}}}
}

func TestProjectRejectsInvalidFullScope(t *testing.T) {
	for _,bad:=range []string{"","*","  "} {
		in:=allowedFixture();in.Scope.Tenant=bad
		if got:=Project(in);got.Status!="deny" || len(got.Entries)!=0 {t.Fatalf("scope %q leaked",bad)}
		in=allowedFixture();in.Scope.Environment=bad
		if got:=Project(in);got.Status!="deny" || len(got.Entries)!=0 {t.Fatalf("environment %q leaked",bad)}
	}
	for _,p:=range []string{"deny","unknown",""}{
		in:=allowedFixture();in.Scope.CasePermission=p
		if Project(in).Status!="deny" {t.Fatal("case permission must fail closed")}
		in=allowedFixture();in.Scope.TimelinePermission=p
		if Project(in).Status!="deny" {t.Fatal("Timeline Entry permission must fail closed")}
	}
}

func TestProjectSuppressesSourceIdentifiers(t *testing.T) {
	for _,p:=range []string{"deny","unknown","missing"} {
		in:=allowedFixture();in.Entries[0].ID="TOP-SECRET";in.Entries[0].SourceRef="secret-ref";in.Entries[0].SourcePermission=p
		got:=Project(in);data,err:=json.Marshal(got)
		if err!=nil {t.Fatal(err)}
		if strings.Contains(string(data),"TOP-SECRET") || strings.Contains(string(data),"secret-ref") {t.Fatalf("denied source disclosure: %s",data)}
	}
	in:=allowedFixture();in.Entries[0].Tenant="another-tenant"
	if len(Project(in).Entries)!=0 {t.Fatal("cross-tenant source visible")}
}

func TestProjectRejectsDuplicateOrInvalidEntries(t *testing.T) {
	mutations:=[]func(*Input){
		func(i *Input){i.Entries=append(i.Entries,i.Entries[0])},
		func(i *Input){i.Entries[0].ID=""},
		func(i *Input){i.Entries[0].Kind="invented"},
		func(i *Input){i.Entries[0].SourcePermission="garbage"},
		func(i *Input){i.Entries[0].Provenance=""},
		func(i *Input){i.Entries[0].SourceRef=""},
		func(i *Input){i.Entries[0].Version=0},
		func(i *Input){i.Entries[0].Supersedes="x"},
		func(i *Input){i.Entries[0].ClockUncertaintyMS=-1},
	}
	for _,change:=range mutations {
		in:=allowedFixture();change(&in)
		if r:=Project(in);r.Status!="deny" || len(r.Entries)!=0 {t.Fatalf("malformed entry leaked: %+v",r)}
	}
}

func TestProjectBoundsAndCallerImmutability(t *testing.T) {
	in:=allowedFixture()
	in.Entries=make([]Entry,MaxEntries+1)
	if got:=Project(in);got.Status!="deny" {t.Fatal("entry bound not enforced")}
	in=allowedFixture();in.Entries[0].Provenance=strings.Repeat("x",MaxInputBytes+1)
	if got:=Project(in);got.Status!="deny" {t.Fatal("input byte bound not enforced")}
	in=allowedFixture()
	out:=Project(in)
	if out.Status!="allow" || len(out.Entries)!=1 {t.Fatal("expected authorized result")}
	out.Entries[0].ID="changed"
	if in.Entries[0].ID!="x" {t.Fatal("mutation of caller state")}
}

func TestProjectTimeOrderingAndUncertainty(t *testing.T) {
	in:=allowedFixture()
	second:=in.Entries[0];second.ID="a";second.SourceTime="2026-10-01T11:00:00+01:00"
	third:=in.Entries[0];third.ID="z";third.SourceTime="2026-10-01T09:00:00Z"
	unknown:=in.Entries[0];unknown.ID="unknown";unknown.Timezone="unknown"
	in.Entries=[]Entry{second,unknown,third,in.Entries[0]}
	out:=Project(in)
	want:=[]string{"z","a","x","unknown"}
	var got []string
	for _,e:=range out.Entries {got=append(got,e.ID)}
	if !reflect.DeepEqual(got,want) {t.Fatalf("order=%v want=%v",got,want)}
	if out.Status!="uncertain" || !out.Entries[len(out.Entries)-1].TimeUncertain {t.Fatal("unknown timezone promoted to known UTC")}
	in=allowedFixture();in.Entries[0].SourceTime="not-a-date"
	if got:=Project(in);got.Status!="uncertain" || !got.Entries[0].TimeUncertain {t.Fatal("malformed timestamp promoted")}
	in=allowedFixture();in.Entries[0].ClockUncertaintyMS=500
	if got:=Project(in);got.Status!="uncertain" {t.Fatal("clock skew not marked")}
}

func TestProjectMissingAndSuperseded(t *testing.T) {
	in:=allowedFixture();old:=in.Entries[0];old.ID="old"
	newer:=in.Entries[0];newer.ID="new";newer.Version=2;newer.Supersedes="old"
	in.Entries=[]Entry{old,newer}
	got:=Project(in)
	if len(got.Entries)!=2 || got.Entries[0].ID!="new" || got.Entries[1].ID!="old" || got.Entries[0].Supersedes!="old" {
		t.Fatal("supersession or stable order not preserved")
	}
}
