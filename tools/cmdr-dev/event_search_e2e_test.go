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

func TestFrontendE2EHandoffRejectsReadinessOverclaim(t *testing.T) {
	head := strings.Repeat("a", 40)
	progress := EventSearchFrontendE2EProgress{
		SchemaVersion:      1,
		WorkUnit:           "E10-INV-002F-E2E",
		Status:             "VERIFIED",
		FinalValidatedHead: head,
	}
	progress.Validation.Result = "PASS"
	progress.FrontendSecurity.GlobalRuntimeCoveragePercent = 91
	progress.FrontendSecurity.AuthorizationNegative = "PASS"
	progress.FrontendSecurity.TenantIsolationNegative = "PASS"
	progress.FrontendSecurity.HostileDOMContent = "PASS"
	progress.FrontendSecurity.PermissionDenialClears = "PASS"
	progress.FrontendSecurity.CancellationStaleBlocking = "PASS"
	progress.FrontendSecurity.OfflineSemanticState = "PASS"
	progress.Performance.BudgetMS = 16
	progress.Performance.PushObservedMS = 0.1
	progress.Performance.PullRequestObservedMS = 0.2
	progress.Performance.AbsolutePass = true

	handoff := EventSearchFrontendE2EHandoff{
		SchemaVersion:      1,
		WorkUnit:           "E10-INV-002F-E2E",
		Result:             "VERIFIED",
		FinalValidatedHead: head,
		PullRequest:        36,
	}
	handoff.Evidence.PushRun = 1
	handoff.Evidence.PullRequestRun = 2
	handoff.Evidence.NodeTestFiles = 5
	handoff.Evidence.ProductionModules = 4
	handoff.Evidence.GlobalRuntimeCoveragePercent = 91
	handoff.Evidence.FrontendP95PushMS = 0.1
	handoff.Evidence.FrontendP95PullRequestMS = 0.2
	handoff.Evidence.FrontendP95BudgetMS = 16

	if err := validateEventSearchProgressIdentity(progress.WorkUnit, progress.Status, progress.FinalValidatedHead, progress.ProductSpecMutated, "E10-INV-002F-E2E"); err != nil {
		t.Fatal(err)
	}
	handoff.FullCapabilityCompletionClaim = true
	if !handoff.FullCapabilityCompletionClaim {
		t.Fatal("expected overclaim fixture")
	}
}
