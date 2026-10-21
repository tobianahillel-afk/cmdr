package investigationtimeline

import (
	"net/url"
	"strings"
	"testing"
)

func TestIncidentDetailReturnHintIsStrictAndDistinctFromLegacyQueue(t *testing.T) {
	cases := []struct {
		query string
		wantOrigin string
		wantDetail bool
		ok bool
	}{
		{"case_id=case-a&origin_incident=inc-a&origin_view=detail", "inc-a", true, true},
		{"case_id=case-a&origin_incident=inc-a", "inc-a", false, true},
		{"case_id=case-a", "", false, true},
		{"case_id=case-a&origin_view=detail", "", false, false},
		{"case_id=case-a&origin_incident=inc-a&origin_view=", "", false, false},
		{"case_id=case-a&origin_incident=inc-a&origin_view=queue", "", false, false},
		{"case_id=case-a&origin_incident=inc-a&origin_view=detail&origin_view=detail", "", false, false},
		{"case_id=case-a&origin_incident=%0A&origin_view=detail", "", false, false},
	}
	for _, tc := range cases {
		values, err := url.ParseQuery(tc.query)
		if err != nil {
			t.Fatal(err)
		}
		origin, detail, ok := extractIncidentReturn(values)
		if origin != tc.wantOrigin || detail != tc.wantDetail || ok != tc.ok {
			t.Fatalf("untrusted query %q origin=%q detail=%v ok=%v", tc.query, origin, detail, ok)
		}
		if ok && (values.Has("origin_incident") || values.Has("origin_view")) {
			t.Fatalf("untrusted origin keys escaped strict selection parser: %v", values)
		}
	}
	if got := commandIncidentReturnURL("inc-a", true); got != "/command/incident-detail?incident=inc-a" {
		t.Fatalf("untrusted destination: %s", got)
	}
	if got := commandIncidentReturnURL("inc-a"); !strings.HasPrefix(got, "/command/work-queue?") ||
		!strings.Contains(got, "restore=1") {
		t.Fatalf("legacy Case return changed: %s", got)
	}
}

func TestIncidentDetailReturnSurvivesOnlyCanonicalSameCaseSourceURL(t *testing.T) {
	const raw = "/case-timeline?case_id=case-a&origin_incident=inc-a&origin_view=detail&kind=observed&limit=1&offset=0"
	canonical, ok := validatedSourceReturnURL(raw, "case-a")
	if !ok || canonical == "" {
		t.Fatalf("authorized selection rejected: %s", canonical)
	}
	back, err := url.Parse(canonical)
	if err != nil || back.Path != "/case-timeline" ||
		back.Query().Get("origin_incident") != "inc-a" ||
		back.Query().Get("origin_view") != "detail" ||
		back.Query().Get("kind") != "observed" ||
		back.Query().Get("limit") != "1" {
		t.Fatalf("canonical Case/Incident return lost: %s %v", canonical, err)
	}
	for _, bad := range []string{
		"/case-timeline?case_id=case-b&origin_incident=inc-a&origin_view=detail&limit=1&offset=0",
		"//evil.test/case-timeline?case_id=case-a&origin_incident=inc-a&origin_view=detail&limit=1&offset=0",
		"/case-timeline?case_id=case-a&origin_incident=inc-a&origin_view=https://evil.test&limit=1&offset=0",
		"/case-timeline?case_id=case-a&origin_incident=inc-a&origin_view=detail&tenant=other&limit=1&offset=0",
		"/case-timeline?case_id=case-a&origin_incident=inc-a&origin_view=detail&origin_view=detail&limit=1&offset=0",
	} {
		if target, ok := validatedSourceReturnURL(bad, "case-a"); ok {
			t.Fatalf("hostile source return was accepted: %s", target)
		}
	}
}
