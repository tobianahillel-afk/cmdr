package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func supersessionFixture() Input {
	in := allowedFixture()
	old := in.Entries[0]
	old.ID = "historical-v1"
	old.Version = 1
	latest := old
	latest.ID = "current-v2"
	latest.Version = 2
	latest.Supersedes = old.ID
	in.Entries = []Entry{latest, old}
	return in
}

func TestSupersessionVisibilityOnlyFromAuthorizedSuccessors(t *testing.T) {
	for _, tc := range []struct {
		name, permission, tenant string
		wantHistorical bool
	}{
		{"authorized", "allow", "tenant-a", true},
		{"denied", "deny", "tenant-a", false},
		{"unknown", "unknown", "tenant-a", false},
		{"missing", "missing", "tenant-a", false},
		{"foreign", "allow", "tenant-b", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := supersessionFixture()
			in.Entries[0].SourcePermission = tc.permission
			in.Entries[0].Tenant = tc.tenant
			out := Project(in)
			if out.Status == "deny" {
				t.Fatalf("hidden successor denied otherwise valid Case: %+v", out)
			}
			var foundHistorical bool
			for _, row := range out.Entries {
				if row.ID == "historical-v1" {
					foundHistorical = true
					if row.Superseded != tc.wantHistorical {
						t.Fatalf("source visibility mismatch: %+v", out)
					}
				}
				if row.ID == "current-v2" && (tc.permission != "allow" || tc.tenant != "tenant-a") {
					t.Fatalf("inaccessible successor leaked: %+v", out)
				}
			}
			if !foundHistorical {
				t.Fatalf("visible history removed: %+v", out)
			}
			encoded, err := json.Marshal(out)
			if err != nil {
				t.Fatal(err)
			}
			if !tc.wantHistorical && (strings.Contains(string(encoded), "current-v2") ||
				strings.Contains(string(encoded), "\"superseded\":true")) {
				t.Fatalf("hidden successor changed observable history: %s", encoded)
			}
		})
	}
}

func TestSupersessionVisibilityMarksWholeAuthorizedCorrectionChain(t *testing.T) {
	in := supersessionFixture()
	next := in.Entries[0]
	next.ID = "current-v3"
	next.Version = 3
	next.Supersedes = "current-v2"
	in.Entries = append(in.Entries, next)
	result := Project(in)
	if result.Status != "allow" || len(result.Entries) != 3 {
		t.Fatalf("authorized chain discarded: %+v", result)
	}
	status := make(map[string]bool, 3)
	for _, row := range result.Entries {
		status[row.ID] = row.Superseded
	}
	if !status["historical-v1"] || !status["current-v2"] || status["current-v3"] {
		t.Fatalf("incorrect chain status: %+v", status)
	}
	if in.Entries[1].SourcePermission != "allow" {
		t.Fatal("caller owner snapshot was modified")
	}
}

func TestSupersessionVisibilityPagedAuthenticatedCaseHTTP(t *testing.T) {
	handler, _ := httpWorkspaceFixture(t)
	first := httpWorkspaceGet(handler,
		"/case-timeline?case_id=case-one&limit=1&offset=0",
		"Bearer analyst-a", nil)
	if first.Code != http.StatusOK {
		t.Fatalf("authorized page status: %d %s", first.Code, first.Body.String())
	}
	p0 := decodeTimelinePage(t, first.Body.Bytes())
	if p0.Page.TotalVisible != 2 || len(p0.Entries) != 1 ||
		p0.Entries[0].ID != "entry-old" || !p0.Entries[0].Superseded {
		t.Fatalf("version on another page not reflected in permitted history: %+v", p0)
	}
	second := httpWorkspaceGet(handler,
		"/case-timeline?case_id=case-one&limit=1&offset=1",
		"Bearer analyst-a", nil)
	p1 := decodeTimelinePage(t, second.Body.Bytes())
	if second.Code != http.StatusOK || len(p1.Entries) != 1 ||
		p1.Entries[0].ID != "entry-new" || p1.Entries[0].Superseded ||
		p1.Entries[0].Supersedes != "entry-old" {
		t.Fatalf("source-owned successor not preserved: %+v", p1)
	}
	filtered := httpWorkspaceGet(handler,
		"/case-timeline?case_id=case-one&kind=inferred&limit=1",
		"Bearer analyst-a", nil)
	fp := decodeTimelinePage(t, filtered.Body.Bytes())
	if filtered.Code != http.StatusOK || len(fp.Entries) != 1 ||
		fp.Entries[0].ID != "entry-old" || !fp.Entries[0].Superseded {
		t.Fatalf("source-visible replacement lost to non-security filter: %+v", fp)
	}
	for _, tc := range []struct{ caseID, bearer string }{
		{"case-one", ""}, {"case-two", "Bearer analyst-a"},
	} {
		resp := httpWorkspaceGet(handler,
			"/case-timeline?case_id="+tc.caseID+"&limit=1", tc.bearer, nil)
		if resp.Code == http.StatusOK || strings.Contains(resp.Body.String(), "superseded") ||
			strings.Contains(resp.Body.String(), "entry-new") {
			t.Fatalf("denied Case leaked correction data: %d %s", resp.Code, resp.Body.String())
		}
	}
}
