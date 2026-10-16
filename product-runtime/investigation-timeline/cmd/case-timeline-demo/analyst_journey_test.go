package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// This test crosses actual HTTP handler boundaries with the SAME authenticated
// synthetic owner dataset. No production provider or browser is implied.
func TestAnalystFullInvestigationJourneyAcrossHTTPRoutes(t *testing.T) {
	for _, tc := range []struct {
		principal, ownCase, otherCase, entry, kind, ownerKind, ownerObject string
	}{
		{"analyst-a", "case-a", "case-b", "event-new", "observed", "telemetry-event", "demo-event-a"},
		{"analyst-b", "case-b", "case-a", "event-b", "observed", "incident", "demo-incident"},
	} {
		t.Run(tc.principal, func(t *testing.T) {
			handler, err := newDemoHandler(testBearer, tc.principal)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(handler)
			defer server.Close()

			status, body := requestDemo(t, server, "/demo-cases", testBearer)
			var cases struct { Cases []string `json:"cases"` }
			if status != http.StatusOK || json.Unmarshal([]byte(body), &cases) != nil ||
				len(cases.Cases) != 1 || cases.Cases[0] != tc.ownCase {
				t.Fatalf("Case owner discovery broken: %d %s", status, body)
			}
			for _, forbidden := range []string{tc.otherCase, "demo-secret", "demo-hidden-source"} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("Case discovery exposed %q", forbidden)
				}
			}

			selection := url.Values{"case_id": {tc.ownCase}, "kind": {tc.kind}, "limit": {"20"}, "offset": {"0"}}
			status, body = requestDemo(t, server, "/case-timeline?"+selection.Encode(), testBearer)
			var page struct { Entries []timeline.ProjectedEntry `json:"entries"` }
			if status != http.StatusOK || json.Unmarshal([]byte(body), &page) != nil {
				t.Fatalf("Case timeline inaccessible: %d %s", status, body)
			}
			found := false
			for _, entry := range page.Entries {
				if entry.ID == tc.entry && entry.Kind == tc.kind {
					found = true
				}
			}
			if !found {
				t.Fatalf("selected Case entry missing from live HTTP result: %s", body)
			}
			for _, forbidden := range []string{"demo-hidden-source", "demo-secret", tc.otherCase} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("timeline disclosed %q", forbidden)
				}
			}
			navigation := url.Values{"case_id": {tc.ownCase}, "entry_id": {tc.entry},
				"kind": {tc.kind}, "limit": {"20"}, "offset": {"0"}}
			status, body = requestDemo(t, server, "/case-source-navigation?"+navigation.Encode(), testBearer)
			var draft struct {
				Status, Kind, ObjectID, ReturnCaseID, ReturnEntryID, ReturnQuery string
				ReturnIsAuthority bool `json:"return_is_authority"`
			}
			if status != http.StatusOK || json.Unmarshal([]byte(body), &draft) != nil ||
				draft.Status != "ready" || draft.Kind != tc.ownerKind ||
				draft.ObjectID != tc.ownerObject || draft.ReturnIsAuthority ||
				draft.ReturnCaseID != tc.ownCase || draft.ReturnEntryID != tc.entry {
				t.Fatalf("owner reauthorization navigation failed: %d %s", status, body)
			}
			if draft.ReturnQuery != selection.Encode() {
				t.Fatalf("return context was not the selected authorized Case: %q", draft.ReturnQuery)
			}

			open := url.Values{"case_id": {tc.ownCase}, "entry_id": {tc.entry}}
			status, body = requestDemo(t, server, "/demo-source-open?"+open.Encode(), testBearer)
			var preview struct {
				Status, CaseID, EntryID, Kind, ObjectID, Preview string
				Synthetic bool
			}
			if status != http.StatusOK || json.Unmarshal([]byte(body), &preview) != nil ||
				preview.Status != "ready" || !preview.Synthetic ||
				preview.Kind != tc.ownerKind || preview.ObjectID != tc.ownerObject ||
				preview.Preview == "" {
				t.Fatalf("typed source owner preview unavailable: %d %s", status, body)
			}
			if strings.Contains(body, "DO NOT DISPLAY") || strings.Contains(body, "demo-secret") {
				t.Fatal("denied source metadata disclosed")
			}

			status, body = requestDemo(t, server, "/case-timeline?"+draft.ReturnQuery, testBearer)
			if status != http.StatusOK || !strings.Contains(body, tc.entry) {
				t.Fatalf("authenticated Case return failed: %d %s", status, body)
			}
			for _, endpoint := range []string{
				"/demo-cases",
				"/case-timeline?"+selection.Encode(),
				"/case-source-navigation?"+navigation.Encode(),
				"/demo-source-open?"+open.Encode(),
			} {
				status, body = requestDemo(t, server, endpoint, "")
				if status != http.StatusUnauthorized || strings.Contains(body, tc.entry) {
					t.Fatalf("anonymous request borrowed journey authority: %d %s", status, body)
				}
			}
			for _, endpoint := range []string{
				"/case-timeline?case_id="+tc.otherCase,
				"/case-source-navigation?case_id="+tc.otherCase+"&entry_id="+tc.entry,
				"/demo-source-open?case_id="+tc.otherCase+"&entry_id="+tc.entry,
				"/case-source-navigation?case_id="+tc.ownCase+"&entry_id=restricted",
				"/demo-source-open?case_id="+tc.ownCase+"&entry_id=restricted",
			} {
				status, body = requestDemo(t, server, endpoint, testBearer)
				if status != http.StatusNotFound || strings.Contains(body, tc.ownerObject) ||
					strings.Contains(body, "demo-secret") || strings.Contains(body, "DO NOT DISPLAY") {
					t.Fatalf("cross-tenant/denied source leak: %d %s", status, body)
				}
			}
		})
	}
}
