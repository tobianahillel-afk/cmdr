package workqueue

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type fakeOwner struct {
	items []Item
	allowed map[string]bool
	errList bool
	errRead bool
	calls []string
}
func (o *fakeOwner) List(_ context.Context, _, _ string) ([]Item,error) {
	if o.errList {return nil,errors.New("provider failed")}
	return o.items,nil
}
func (o *fakeOwner) CanRead(_ context.Context, principal, tenant, environment string, kind Kind,id string)(bool,error){
	o.calls=append(o.calls,string(kind)+":"+id)
	if o.errRead {return false,errors.New("policy unavailable")}
	return principal=="analyst" && tenant=="t1" && environment=="e1" && o.allowed[string(kind)+":"+id],nil
}
func enabled(b bool)*bool{return &b}
func fixture()*fakeOwner{return &fakeOwner{
	items:[]Item{
		{Kind:Incident,ID:"inc-a",Tenant:"t1",Environment:"e1",Owner:"analyst",Status:"active",UpdatedAt:"2026-10-19T08:00:00Z",SLAAtRisk:enabled(true)},
		{Kind:Task,ID:"task-b",Tenant:"t1",Environment:"e1",Followers:[]string{"analyst"},Status:"open",UpdatedAt:"2026-10-19T09:00:00Z",SLAAtRisk:enabled(false)},
		{Kind:Incident,ID:"inc-hidden",Tenant:"t1",Environment:"e1",Owner:"secret",Status:"secret",UpdatedAt:"2026-10-19T10:00:00Z",SLAAtRisk:enabled(true)},
		{Kind:Task,ID:"task-no-sla",Tenant:"t1",Environment:"e1",Owner:"other",Status:"open",UpdatedAt:"2026-10-19T07:00:00Z"},
		{Kind:Incident,ID:"foreign",Tenant:"t2",Environment:"e1",Status:"private",UpdatedAt:"2026-10-19T12:00:00Z"},
	},
	allowed:map[string]bool{"Incident:inc-a":true,"Task:task-b":true,"Task:task-no-sla":true},
}}
func request(view View)Request{return Request{Principal:"analyst",Tenant:"t1",Environment:"e1",View:view}}
func TestCanonicalSixViewsWithRealAuthorization(t *testing.T){
	want:=map[View][]string{
		All:{"task-b","inc-a","task-no-sla"},
		Incidents:{"inc-a"},
		Tasks:{"task-b","task-no-sla"},
		Unassigned:{"task-b"},
		MyWork:{"task-b","inc-a"},
		SLARisk:{"inc-a"},
	}
	if len(SystemViews)!=6 {t.Fatal("canonical catalogue no longer has six views")}
	for _,view:=range SystemViews{
		t.Run(string(view),func(t *testing.T){
			provider:=fixture()
			result,err:=Project(context.Background(),provider,request(view))
			if err!=nil {t.Fatal(err)}
			var got []string
			for _,row:=range result.Rows{
				got=append(got,row.ID)
				if row.ID=="inc-hidden"||row.ID=="foreign"||row.Status=="secret"{t.Fatalf("permission leak: %+v",row)}
			}
			if !reflect.DeepEqual(got,want[view]){t.Fatalf("expected %v, got %v",want[view],got)}
			if result.View!=view||result.Tenant!="t1"||result.Environment!="e1"{t.Fatal("changed routed view or context")}
			for _,call:=range provider.calls{if call=="Incident:foreign"{t.Fatal("queried foreign ID authorization")}}
			if view==SLARisk{
				if result.State!="partial-sources"||!reflect.DeepEqual(result.Diagnostics,[]string{"sla-policy-unavailable"}){t.Fatal("unknown SLA policy was treated as healthy")}
			}else if result.State!="ready"{t.Fatalf("unexpected state %q",result.State)}
		})
	}
}
func TestReadOnlyQueueFailsClosedForMissingScopeViewAndProvider(t *testing.T){
	for _,req:=range []Request{
		{Principal:"",Tenant:"t1",Environment:"e1",View:All},
		{Principal:"analyst",Tenant:"*",Environment:"e1",View:All},
		{Principal:"analyst",Tenant:"t1",Environment:"",View:All},
		{Principal:"analyst",Tenant:"t1\n",Environment:"e1",View:All},
		request("team-load"),request(""),
	}{
		o:=fixture();r,err:=Project(context.Background(),o,req)
		if err==nil||len(r.Rows)!=0||len(o.calls)!=0{t.Fatalf("invalid request exposed data: %+v %+v",req,r)}
	}
	for _,o:=range []*fakeOwner{{errList:true},{errRead:true,items:fixture().items}}{
		o.allowed=fixture().allowed
		r,err:=Project(context.Background(),o,request(All))
		if err==nil||len(r.Rows)!=0{t.Fatalf("partial provider failure returned data: %v %+v",err,r)}
	}
	o:=fixture()
	if r,err:=Project(nil,o,request(All));err==nil||len(r.Rows)!=0{t.Fatal("nil context authorized")}
	c,cancel:=context.WithCancel(context.Background());cancel()
	if r,err:=Project(c,o,request(All));err==nil||len(r.Rows)!=0{t.Fatal("cancelled context authorized")}
	if r,err:=Project(context.Background(),nil,request(All));err==nil||len(r.Rows)!=0{t.Fatal("nil owner authorized")}
}
func TestNoLeakOnIncompleteDuplicateOrOversizedOwnedProjection(t *testing.T){
	for _,modify:=range []func(*fakeOwner){
		func(o *fakeOwner){o.items[0].UpdatedAt="not-utc"},
		func(o *fakeOwner){o.items[0].Kind="Case"},
		func(o *fakeOwner){o.items[0].ID=""},
		func(o *fakeOwner){o.items=append(o.items,o.items[0])},
		func(o *fakeOwner){o.items=make([]Item,5001)},
	}{
		o:=fixture();modify(o);result,err:=Project(context.Background(),o,request(All))
		if err==nil||len(result.Rows)>0{t.Fatalf("invalid owner snapshot disclosed partial rows: %+v",result)}
	}
}
func TestDeniedMetadataAndDuplicatesDoNotChangeAuthorizedProjection(t *testing.T) {
	baselineOwner := fixture()
	baseline, err := Project(context.Background(), baselineOwner, request(All))
	if err != nil {
		t.Fatalf("baseline: %v", err)
	}
	for _, tc := range []struct {
		name   string
		mutate func(*fakeOwner)
	}{
		{"denied-empty-updated-at", func(o *fakeOwner) { o.items[2].UpdatedAt = "" }},
		{"denied-invalid-timestamp", func(o *fakeOwner) { o.items[2].UpdatedAt = "not-utc" }},
		{"denied-duplicate", func(o *fakeOwner) { o.items = append(o.items, o.items[2]) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := fixture()
			tc.mutate(o)
			actual, err := Project(context.Background(), o, request(All))
			if err != nil {
				t.Fatalf("denied record affected valid data: %v", err)
			}
			if !reflect.DeepEqual(actual, baseline) {
				t.Fatalf("denied record changed authorized results: %+v vs %+v", actual, baseline)
			}
		})
	}
}

func TestMyWorkViewIncludesContributorsAndWatchers(t *testing.T) {
	o := fixture()
	o.items = append(o.items,
		Item{Kind: Task, ID: "task-contrib", Tenant: "t1", Environment: "e1", Owner: "other", Contributors: []string{"analyst"}, Status: "open", UpdatedAt: "2026-10-19T08:30:00Z", SLAAtRisk: enabled(false)},
		Item{Kind: Task, ID: "task-watch", Tenant: "t1", Environment: "e1", Owner: "other", Watchers: []string{"analyst"}, Status: "open", UpdatedAt: "2026-10-19T08:15:00Z", SLAAtRisk: enabled(false)},
	)
	o.allowed["Task:task-contrib"] = true
	o.allowed["Task:task-watch"] = true
	res, err := Project(context.Background(), o, request(MyWork))
	if err != nil {
		t.Fatalf("MyWork query failed: %v", err)
	}
	ids := map[string]bool{}
	for _, row := range res.Rows {
		ids[row.ID] = true
	}
	for _, id := range []string{"task-b", "inc-a", "task-contrib", "task-watch"} {
		if !ids[id] {
			t.Fatalf("missing expected MyWork item %s: %+v", id, res.Rows)
		}
	}
}

func TestEmptyAndStableBoundedFirstPage(t *testing.T) {
	o := fixture()
	o.allowed = map[string]bool{}
	r, err := Project(context.Background(), o, request(All))
	if err != nil || r.State != "no-authorized-items" || len(r.Rows) != 0 {
		t.Fatalf("empty failed: %+v %v", r, err)
	}
	o = fixture()
	o.items = o.items[:0]
	o.allowed = map[string]bool{}
	for i := 0; i < 80; i++ {
		id := fmt.Sprintf("task-%03d", i)
		o.items = append(o.items, Item{Kind: Task, ID: id, Tenant: "t1", Environment: "e1", UpdatedAt: "2026-10-19T09:00:00Z"})
		o.allowed["Task:"+id] = true
	}
	a, err := Project(context.Background(), o, request(Tasks))
	if err != nil || len(a.Rows) != 50 || len(a.Diagnostics) != 1 || a.Diagnostics[0] != "first-page-only" {
		t.Fatalf("first-page constraint broken %+v %v", a, err)
	}
	b, err := Project(context.Background(), o, request(Tasks))
	if err != nil || !reflect.DeepEqual(a.Rows, b.Rows) {
		t.Fatal("sort is not deterministic")
	}
}
