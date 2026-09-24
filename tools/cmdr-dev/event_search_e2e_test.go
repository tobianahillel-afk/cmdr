package main

import (
	"strings"
	"testing"
)

func TestParseEventSearchAdversarialTestEvents(t *testing.T) {
	var lines []string
	for _, name := range eventSearchAdversarialTests {
		lines = append(lines, `{"Action":"pass","Test":"`+name+`"}`)
	}
	summary, err := parseEventSearchAdversarialTestEvents([]byte(strings.Join(lines, "\n") + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != len(eventSearchAdversarialTests) {
		t.Fatalf("unexpected adversarial count: %#v", summary)
	}
}

func TestParseEventSearchAdversarialTestEventsFailsClosedWhenMissing(t *testing.T) {
	data := []byte(`{"Action":"pass","Test":"TestValidateContractFixtures"}` + "\n")
	if _, err := parseEventSearchAdversarialTestEvents(data); err == nil {
		t.Fatal("expected incomplete adversarial evidence rejection")
	}
}

func TestEventSearchHandoffRequiredLimitations(t *testing.T) {
	handoff := EventSearchClosureHandoff{
		Limitations: []string{
			"Event Search UI integration remains outside this bounded core",
			"Saved Search behavior is not implemented",
			"Case-link mutation remains excluded while OPEN-013 is unresolved",
			"the final query language and dialect are not selected",
			"the final search/index/storage/provider stack is not selected",
			"no production search SLO is claimed",
		},
	}
	for _, fragment := range eventSearchRequiredLimitations {
		if !containsFoldFragment(handoff.Limitations, fragment) {
			t.Fatalf("missing limitation fragment %q", fragment)
		}
	}
}

func TestValidateEventSearchProgressIdentity(t *testing.T) {
	head := strings.Repeat("a", 40)
	if err := validateEventSearchProgressIdentity("E10-INV-002C-PERF", "VERIFIED", head, false, "E10-INV-002C-PERF"); err != nil {
		t.Fatal(err)
	}
	if err := validateEventSearchProgressIdentity("E10-INV-002C-PERF", "VERIFIED", head, true, "E10-INV-002C-PERF"); err == nil {
		t.Fatal("expected Product Spec mutation rejection")
	}
}
