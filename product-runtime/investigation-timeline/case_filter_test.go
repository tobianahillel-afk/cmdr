package investigationtimeline

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCaseTimelineHTTPFilteredInvestigation(t *testing.T) {
	cases := []struct {
		label, query, status string
		ids                  []string
	}{
		{"observed", "&kind=observed", "filtered", []string{"entry-new"}},
		{"inferred", "&kind=inferred", "filtered", []string{"entry-old"}},
		{"before-11", "&to=2026-10-01T11%3A00%3A00Z", "filtered", []string{"entry-old"}},
		{"11-to-12", "&from=2026-10-01T11%3A00%3A00Z&to=2026-10-01T12%3A00%3A00Z", "filtered", []string{"entry-new"}},
		{"empty-window", "&from=2030-01-01T00%3A00%3A00Z", "filtered", nil},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			h, store := httpWorkspaceFixture(t)
			res := httpWorkspaceGet(h, "/case-timeline?case_id=case-one"+tc.query, "Bearer analyst-a", nil)
			if res.Code != http.StatusOK || store.caseChecks != 2 || store.catalogCalls != 1 || store.entryReads != 3 || store.sourceReads != 3 {
				t.Fatalf("not full ACL/owner chain: HTTP=%d Case=%d catalog=%d Entry=%d source=%d",
					res.Code, store.caseChecks, store.catalogCalls, store.entryReads, store.sourceReads)
			}
			var result Result
			if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.Status != tc.status || len(result.Entries) != len(tc.ids) {
				t.Fatalf("unexpected filtered timeline: %+v", result)
			}
			for i, id := range tc.ids {
				if result.Entries[i].ID != id {
					t.Fatalf("index %d: %+v", i, result)
				}
			}
			for _, secret := range []string{"entry-hidden", "src-hidden", "tenant-b"} {
				if strings.Contains(res.Body.String(), secret) {
					t.Fatalf("hidden disclosure %s", secret)
				}
			}
		})
	}
}

func TestCaseTimelineHTTPRejectsInvalidFiltersWithoutProviderIO(t *testing.T) {
	for _, suffix := range []string{
		"&kind=bogus", "&kind=observed&kind=inferred", "&tenant=tenant-b",
		"&source_permission=allow", "&from=not-a-time",
		"&from=2026-10-01T10%3A00%3A00-00%3A00",
		"&from=2026-10-01T10%3A00%3A00Z&to=2026-10-01T10%3A00%3A00Z",
		"&from=2026-10-02T00%3A00%3A00Z&to=2026-10-01T00%3A00%3A00Z",
		"&to=2026-10-01T10%3A00%3A00Z&permission=allow",
	} {
		h, store := httpWorkspaceFixture(t)
		res := httpWorkspaceGet(h, "/case-timeline?case_id=case-one"+suffix, "Bearer analyst-a", nil)
		if res.Code != http.StatusBadRequest || store.catalogCalls != 0 || store.caseChecks != 0 {
			t.Fatalf("invalid filter was forwarded to provider: %q, HTTP=%d", suffix, res.Code)
		}
	}
}

func TestTimelineSelectionPreservesUncertainChronology(t *testing.T) {
	sel, valid := parseTimelineReadSelection(url.Values{
		"case_id": {"case-one"}, "from": {"2026-10-01T10:00:00Z"}, "to": {"2026-10-01T11:00:00Z"},
	})
	if !valid {
		t.Fatal("refused valid interval")
	}
	in := Result{Status: "allow", Entries: []ProjectedEntry{
		{ID: "inside", Kind: "observed", SourceTime: "2026-10-01T10:00:00Z"},
		{ID: "outside", Kind: "observed", SourceTime: "2026-10-01T11:00:00Z"},
		{ID: "uncertain", Kind: "inferred", SourceTime: "not-a-clock", TimeUncertain: true},
	}}
	out := selectAuthorizedTimeline(in, sel)
	if out.Status != "partial" || len(out.Entries) != 2 || out.Entries[0].ID != "inside" ||
		out.Entries[1].ID != "uncertain" || !out.Entries[1].TimeUncertain {
		t.Fatalf("window silently lost uncertainty: %+v", out)
	}
	sel.kind = "observed"
	out = selectAuthorizedTimeline(in, sel)
	if out.Status != "allow" || len(out.Entries) != 1 || out.Entries[0].ID != "inside" {
		t.Fatalf("kind filter changed chronology: %+v", out)
	}
	if len(in.Entries) != 3 || in.Entries[1].ID != "outside" {
		t.Fatal("input mutated")
	}
}
