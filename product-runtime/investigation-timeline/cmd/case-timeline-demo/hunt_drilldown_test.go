package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestOfflineHuntSelectedCaseProvidesAuthorizedEventDrilldown(t *testing.T) {
	path := diskHuntSnapshot(t)
	var report bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-hunt-q=ioc.example", "-hunt-events-case=case-a",
	}, &report); err != nil {
		t.Fatal(err)
	}
	var result localHuntCaseEvents
	if err := json.Unmarshal(report.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Format != "cmdr-local-selected-case-hunt-events-v1" ||
		!result.Synthetic || result.Query != "ioc.example" ||
		result.CaseID != "case-a" || result.TotalMatches != 1 ||
		result.Returned != 1 || result.Truncated || len(result.Events) != 1 {
		t.Fatalf("selected authorized Case drilldown failed: %+v", result)
	}
	if result.Events[0].ID != "event-old" ||
		result.Events[0].Kind != "inferred" ||
		result.Events[0].SourceRef == "" || result.Events[0].Provenance == "" {
		t.Fatalf("missing owner-authorized source lead: %+v", result.Events)
	}
	for _, forbidden := range []string{
		"restricted", "event-b", "case-b", "case-c", "hidden IOC.EXAMPLE",
		"other-tenant", "Bearer ", "demo-secret",
	} {
		if strings.Contains(report.String(), forbidden) {
			t.Fatalf("disclosed non-selected or denied source %q: %s", forbidden, report.String())
		}
	}
	// A second authorized Case must be explicitly selected independently;
	// search does not silently mix events from different Cases.
	report.Reset()
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-hunt-q=ioc.example", "-hunt-events-case=case-c",
	}, &report); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(report.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.CaseID != "case-c" || len(result.Events) != 1 ||
		result.Events[0].ID != "event-c" || strings.Contains(report.String(), "event-old") {
		t.Fatalf("selected Case isolation broken: %+v", result)
	}
}

func TestOfflineHuntSelectedCasePaginatesAndCapsSensitiveReceipt(t *testing.T) {
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + twoCaseComparisonSnapshot(t),
		"-hunt-q=authorized-on-disk", "-hunt-events-case=case-a",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var result localHuntCaseEvents
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.TotalMatches != 205 || result.Returned != maxHuntDrilldownEvents ||
		!result.Truncated || len(result.Events) != maxHuntDrilldownEvents {
		t.Fatalf("multi-page case hunt truncated incorrectly: %+v", result)
	}
	seen := make(map[string]bool)
	for _, event := range result.Events {
		if event.ID == "" || seen[event.ID] || event.SourceRef == "" {
			t.Fatalf("invalid or repeated authorized event: %+v", event)
		}
		seen[event.ID] = true
	}
}

func TestOfflineHuntSelectedCaseRejectsForeignAndDeniedAccess(t *testing.T) {
	path := diskHuntSnapshot(t)
	for _, tc := range []struct{ principal, caseID string }{
		{"analyst-b", "case-a"},
		{"analyst-a", "case-b"},
		{"analyst-a", "missing-case"},
		{"unknown", "case-a"},
	} {
		var out bytes.Buffer
		err := runDemoCLI(context.Background(), []string{
			"-principal=" + tc.principal, "-snapshot=" + path,
			"-hunt-q=ioc.example", "-hunt-events-case=" + tc.caseID,
		}, &out)
		if err == nil || out.Len() != 0 {
			t.Fatalf("unauthorized selected Case produced events: principal=%q case=%q err=%v output=%s",
				tc.principal, tc.caseID, err, out.String())
		}
	}
}

func TestOfflineHuntSelectedCaseFailsWithoutPartialOutputOnLaterPage(t *testing.T) {
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	deniedSecondPage := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/case-timeline" && r.URL.Query().Get("offset") == "200" {
			http.Error(w, "revoked", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
	got, err := collectAuthorizedHuntCaseEvents(context.Background(), deniedSecondPage,
		testBearer, "case-a", "authorized-on-disk")
	if err == nil || got.Format != "" || len(got.Events) != 0 || got.TotalMatches != 0 {
		t.Fatalf("partial sensitive drilldown escaped after source/page revocation: %+v %v", got, err)
	}
}

func TestOfflineHuntSelectedCaseRequiresExplicitSafeModeAndQuery(t *testing.T) {
	for _, args := range [][]string{
		{"-hunt-events-case=case-a"},
		{"-hunt-q=ioc.example", "-hunt-events-case="},
		{"-hunt-q=x", "-hunt-events-case=case-a"},
		{"-hunt-q=ioc.example", "-hunt-events-case=case-a", "-report-case=case-a"},
		{"-hunt-q=ioc.example", "-hunt-events-case=case-a", "-compare-events"},
		{"-hunt-q=ioc.example", "-hunt-events-case=case-a", "-source-ref=source-x"},
	} {
		var output bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("ambiguous or invalid hunt mode accepted: %v output=%s", args, output.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := runDemoCLI(ctx, []string{
		"-principal=analyst-a", "-snapshot=" + diskHuntSnapshot(t),
		"-hunt-q=ioc.example", "-hunt-events-case=case-a",
	}, &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled hunt exported a sensitive report: %v output=%s", err, out.String())
	}
}
