package workqueue

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// Source-owned contributor/watcher membership determines only whether an
// already authorized Incident/Task appears in the My Work view. It never
// grants the underlying Owner.CanRead permission.
func TestMyWorkOperationalRelationsAreAuthorizationGated(t *testing.T) {
	o := fixture()
	o.items = append(o.items,
		Item{Kind: Task, ID: "contrib-only", Tenant: "t1", Environment: "e1",
			Owner: "manager", Contributors: []string{"analyst"},
			Status: "open", UpdatedAt: "2026-10-19T12:00:00Z"},
		Item{Kind: Incident, ID: "watch-only", Tenant: "t1", Environment: "e1",
			Owner: "manager", Watchers: []string{"analyst"},
			Status: "active", UpdatedAt: "2026-10-19T11:00:00Z"},
		Item{Kind: Task, ID: "denied-contribution", Tenant: "t1", Environment: "e1",
			Owner: "manager", Contributors: []string{"analyst"},
			Status: "confidential", UpdatedAt: "2026-10-19T13:00:00Z"},
		Item{Kind: Task, ID: "foreign-observer", Tenant: "other", Environment: "e1",
			Owner: "manager", Watchers: []string{"analyst"},
			Status: "private", UpdatedAt: "2026-10-19T14:00:00Z"},
	)
	o.allowed["Task:contrib-only"] = true
	o.allowed["Incident:watch-only"] = true
	// Explicitly denied and foreign even though the supplied relationship
	// appears to contain this principal.
	o.allowed["Task:denied-contribution"] = false
	o.allowed["Task:foreign-observer"] = true
	result,err := Project(context.Background(),o,request(MyWork))
	if err != nil {t.Fatal(err)}
	ids := map[string]bool{}
	for _,row := range result.Rows {
		ids[row.ID]=true
		if row.ID=="denied-contribution" || row.ID=="foreign-observer" ||
			row.Status=="confidential" || row.Status=="private" {
			t.Fatalf("membership escalated source permission: %+v",row)
		}
	}
	for _,want := range []string{"inc-a","task-b","contrib-only","watch-only"} {
		if !ids[want] {t.Fatalf("authorized My Work membership missing %q: %+v",want,result.Rows)}
	}
	for _,called := range o.calls {
		if called=="Task:foreign-observer" {t.Fatal("foreign owner accessed through watcher relation")}
	}
	raw,err:=json.Marshal(result)
	if err!=nil {t.Fatal(err)}
	if strings.Contains(string(raw),"Contributors") || strings.Contains(string(raw),"Watchers") ||
		strings.Contains(string(raw),"denied-contribution") {
		t.Fatal("My Work leaked raw source membership or denied identity")
	}
	// Applying an explicit responsible-owner filter still uses the true
	// owner, not contributor/watcher membership or the host principal.
	onlyOwned:=request(MyWork)
	onlyOwned.OwnerFilter="analyst"
	filtered,err:=Project(context.Background(),o,onlyOwned)
	if err!=nil || len(filtered.Rows)!=1 || filtered.Rows[0].ID!="inc-a" {
		t.Fatalf("operational membership impersonated item owner: %+v %v",filtered,err)
	}
}

func TestMyWorkOperationalRelationsPreserveExistingRouteAndRevocation(t *testing.T) {
	o:=fixture()
	o.items=append(o.items,Item{
		Kind:Task,ID:"observed-work",Tenant:"t1",Environment:"e1",
		Owner:"someone-else",Watchers:[]string{"analyst"},
		Status:"open",UpdatedAt:"2026-10-19T11:15:00Z",
	})
	o.allowed["Task:observed-work"]=true
	h:=NewHTTPHandler(trustedHost(),o)
	r:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=my-work")
	if r.Code!=http.StatusOK || !strings.Contains(r.Body.String(),"observed-work") {
		t.Fatalf("authorized watcher not visible in current Work Queue: %d %s",r.Code,r.Body.String())
	}
	o.allowed["Task:observed-work"]=false
	r=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=my-work")
	if r.Code!=http.StatusOK || strings.Contains(r.Body.String(),"observed-work") {
		t.Fatalf("revoked watcher remained visible: %d %s",r.Code,r.Body.String())
	}
	// No new query parameter, no contributor/watcher impersonation as the
	// SessionResolver; existing route guards still reject scope injection.
	bad:=testQueueRequest(h,http.MethodGet,"/command/work-queue?view=my-work&principal=analyst")
	if bad.Code!=http.StatusBadRequest {t.Fatalf("client principal injection accepted: %d",bad.Code)}
}
