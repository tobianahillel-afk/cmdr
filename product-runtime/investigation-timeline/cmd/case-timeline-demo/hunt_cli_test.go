package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A real operator-selected disk snapshot exercises the same in-process
// authenticated Case/Timeline/source boundary as the interactive analyst.
func diskHuntSnapshot(t *testing.T) string {
	t.Helper()
	data, err := loadDemoSnapshot(twoCaseComparisonSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := range data.Entries {
		switch data.Entries[i].Value.ID {
		case "event-old":
			data.Entries[i].Value.Provenance = "dns IOC.EXAMPLE"
		case "event-c":
			data.Entries[i].Value.Provenance = "forensic IOC.EXAMPLE"
		case "restricted":
			data.Entries[i].Value.Provenance = "hidden IOC.EXAMPLE"
		case "event-b":
			data.Entries[i].Value.Provenance = "other-tenant IOC.EXAMPLE"
		}
	}
	blob, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "real-hunt-snapshot.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestOfflineHuntActualDiskTwoCaseAuthorizedText(t *testing.T) {
	var output bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + diskHuntSnapshot(t),
		"-hunt-q=ioc.example",
	}, &output)
	if err != nil {
		t.Fatal(err)
	}
	var got localCrossCaseHunt
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Format != "cmdr-local-cross-case-text-hunt-v1" || !got.Synthetic ||
		got.Query != "ioc.example" || got.SearchedCases != 2 ||
		got.MatchingCases != 2 || got.VisibleMatches != 2 ||
		len(got.Results) != 2 {
		t.Fatalf("real authorized cross-Case search incomplete: %+v", got)
	}
	var found = map[string]localHuntCase{}
	for _, item := range got.Results {
		found[item.CaseID] = item
	}
	if found["case-a"].Matches != 1 || found["case-a"].Inferred != 1 ||
		found["case-c"].Matches != 1 || found["case-c"].Inferred != 1 {
		t.Fatalf("Case-specific search counts incorrect: %+v", found)
	}
	for _, denied := range []string{
		"restricted", "hidden", "event-old", "event-c", "event-b",
		"case-b", "demo-hidden-source", "demo-secret",
		"Bearer", "demo-tenant", "SENSITIVE-PREVIEW",
	} {
		if strings.Contains(output.String(), denied) {
			t.Fatalf("hunt exported hidden or raw metadata %q: %s", denied, output.String())
		}
	}
}

func TestOfflineHuntPagesAllMatchesInsteadOfFirst200(t *testing.T) {
	var output bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + twoCaseComparisonSnapshot(t),
		"-hunt-q=authorized-on-disk",
	}, &output)
	if err != nil {
		t.Fatal(err)
	}
	var got localCrossCaseHunt
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.SearchedCases != 2 || got.MatchingCases != 1 ||
		got.VisibleMatches != 205 || len(got.Results) != 1 ||
		got.Results[0].CaseID != "case-a" ||
		got.Results[0].Observed != 103 || got.Results[0].Inferred != 102 {
		t.Fatalf("multi-page authenticated hunt incorrect: %+v", got)
	}
}

func TestOfflineHuntDoesNotInferMatchesFromInaccessibleSource(t *testing.T) {
	path := diskHuntSnapshot(t)
	var out bytes.Buffer
	err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path,
		"-hunt-q=hidden",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	var got localCrossCaseHunt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.MatchingCases != 0 || got.VisibleMatches != 0 ||
		len(got.Results) != 0 || strings.Contains(out.String(), "restricted") {
		t.Fatalf("protected source was searchable: %+v", got)
	}
	// A different principal is explicitly resolved by the server, not by
	// exposing the first principal's Case list.
	out.Reset()
	err = runDemoCLI(context.Background(), []string{
		"-principal=analyst-b", "-snapshot=" + path,
		"-hunt-q=ioc.example",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.SearchedCases != 1 || got.MatchingCases != 1 ||
		got.VisibleMatches != 1 || len(got.Results) != 1 ||
		got.Results[0].CaseID != "case-b" ||
		strings.Contains(out.String(), "case-a") || strings.Contains(out.String(), "case-c") {
		t.Fatalf("foreign tenant leaked through hunt: %+v", got)
	}
}

func TestOfflineHuntRejectsMixedModesAndInvalidQueries(t *testing.T) {
	for _, args := range [][]string{
		{"-hunt-q="},
		{"-hunt-q=x"},
		{"-hunt-q= "},
		{"-hunt-q=line\nbreak"},
		{"-hunt-q=" + strings.Repeat("a", 97)},
		{"-hunt-q=valid", "-report-case=case-a"},
		{"-hunt-q=valid", "-compare-events"},
		{"-hunt-q=valid", "-inspect-entry=event-old"},
		{"-hunt-q=valid", "-source-case=case-a"},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("bad mode/query accepted or leaked output: %q %v", args, err)
		}
	}
}

func TestOfflineHuntDoesNotPublishPartialResultAfterSecondCaseDenial(t *testing.T) {
	data, err := loadDemoSnapshot(diskHuntSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	blocked := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/case-timeline" && r.URL.Query().Get("case_id") == "case-c" {
			http.Error(w, "unavailable", http.StatusNotFound)
			return
		}
		handler.ServeHTTP(w, r)
	})
	got, err := collectAuthorizedCrossCaseHunt(context.Background(), blocked, testBearer, "ioc.example")
	if err == nil || got.Format != "" || got.VisibleMatches != 0 || len(got.Results) != 0 {
		t.Fatalf("second Case denial produced partial cross-Case report: %+v %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := collectAuthorizedCrossCaseHunt(ctx, handler, testBearer, "ioc.example"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled hunt was not rejected: %v", err)
	}
}

func TestOfflineHuntBoundsAuthorizedCatalogBeforeAnyTimelineRead(t *testing.T) {
	requests := 0
	fake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/demo-cases" {
			t.Fatalf("oversized catalog should not trigger timeline request: %s", r.URL)
		}
		cases := make([]string, maxHuntCases+1)
		for i := range cases {
			cases[i] = fmt.Sprintf("authorized-case-%d", i)
		}
		_ = json.NewEncoder(w).Encode(struct {
			Cases []string `json:"cases"`
		}{cases})
	})
	result, err := collectAuthorizedCrossCaseHunt(context.Background(), fake, "ephemeral", "ioc.example")
	if err == nil || result.Format != "" || requests != 1 {
		t.Fatalf("oversized authorized Case catalog accepted: %+v %v requests=%d", result, err, requests)
	}
}

func TestOfflineHuntEmptyResultsAndOutputFailure(t *testing.T) {
	path := diskHuntSnapshot(t)
	var buffer bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path, "-hunt-q=not-present",
	}, &buffer); err != nil {
		t.Fatal(err)
	}
	var result localCrossCaseHunt
	if err := json.Unmarshal(buffer.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.SearchedCases != 2 || result.VisibleMatches != 0 ||
		result.MatchingCases != 0 || len(result.Results) != 0 {
		t.Fatalf("empty authorized result misreported: %+v", result)
	}
	if err := writeDemoCrossCaseHunt(context.Background(), "analyst-a", path,
		"ioc.example", failingHuntWriter{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("failed output writer was treated as success: %v", err)
	}
}

type failingHuntWriter struct{}

func (failingHuntWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestOfflineHuntRejectsMalformedDiscoveredCaseCatalog(t *testing.T) {
	data, err := loadDemoSnapshot(diskHuntSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	realHandler, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	for _, list := range [][]string{
		nil,
		{"case-a", "case-a"},
		{"case-a", ""},
		{"case-a", strings.Repeat("x", 1025)},
	} {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/demo-cases" {
				_ = json.NewEncoder(w).Encode(struct {
					Cases []string `json:"cases"`
				}{list})
				return
			}
			calls++
			realHandler.ServeHTTP(w, r)
		})
		got, err := collectAuthorizedCrossCaseHunt(context.Background(), handler, testBearer, "ioc.example")
		if err == nil || got.VisibleMatches != 0 || len(got.Results) != 0 {
			t.Fatalf("bad discovered Case list was accepted: %v %+v", err, got)
		}
		if calls > 1 {
			t.Fatalf("invalid Case catalog caused extra owner reads: %d", calls)
		}
	}
}

func TestOfflineHuntInvalidQueryBoundaryIsRejectedBeforeOwnerReads(t *testing.T) {
	requests := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		http.Error(w, "not expected", http.StatusInternalServerError)
	})
	for _, query := range []string{"", "a", " bad", "bad ", "\tX", "x\u0000y", strings.Repeat("a", 97), string([]byte{0xff, 0xfe})} {
		if result, err := collectAuthorizedCrossCaseHunt(context.Background(), handler, "token", query); err == nil || result.Format != "" {
			t.Fatalf("invalid query %q accepted: %+v %v", query, result, err)
		}
	}
	if requests != 0 {
		t.Fatalf("invalid queries sent %d owner HTTP requests", requests)
	}
	if !validHuntText(strings.Repeat("x", 96)) || !validHuntText("ioc.example") {
		t.Fatal("valid bounded indicator text rejected")
	}
}

func TestOfflineHuntBadSecondCasePageDiscardsAllCounts(t *testing.T) {
	data, err := loadDemoSnapshot(diskHuntSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	h, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	broken := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/case-timeline" && r.URL.Query().Get("case_id") == "case-c" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"status":"allow","entries":[],"page":{"offset":0,"limit":200,"total_visible":1,"has_more":false}}`)
			return
		}
		h.ServeHTTP(w, r)
	})
	got, err := collectAuthorizedCrossCaseHunt(context.Background(), broken, testBearer, "ioc.example")
	if err == nil || got.SearchedCases != 0 || got.VisibleMatches != 0 {
		t.Fatalf("incomplete authorized second Case page produced partial report: %+v %v", got, err)
	}
}

func TestOfflineHuntUsesBuiltInAuthorizedOwnerWhenNoSnapshotSelected(t *testing.T) {
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-hunt-q=synthetic",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var result localCrossCaseHunt
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.SearchedCases != 1 || result.MatchingCases != 1 ||
		result.VisibleMatches != 2 || len(result.Results) != 1 ||
		result.Results[0].CaseID != "case-a" ||
		result.Results[0].Observed != 1 || result.Results[0].Inferred != 1 {
		t.Fatalf("fallback Case authorization and filtering failed: %+v", result)
	}
	if strings.Contains(out.String(), "demo-hidden-source") ||
		strings.Contains(out.String(), "case-b") {
		t.Fatalf("built-in hunt exported inaccessible metadata: %s", out.String())
	}
}

func TestOfflineHuntUnavailablePrincipalNeverEmitsReceipt(t *testing.T) {
	path := diskHuntSnapshot(t)
	for _, principal := range []string{"not-a-principal", "analyst-elsewhere", "*"} {
		var out bytes.Buffer
		err := runDemoCLI(context.Background(), []string{
			"-principal=" + principal, "-snapshot=" + path, "-hunt-q=ioc.example",
		}, &out)
		if err == nil || out.Len() != 0 {
			t.Fatalf("unknown synthetic principal generated report: %q %v", principal, err)
		}
	}
}

func TestOfflineHuntRejectsBadDiscoveryAndRevocationMidRead(t *testing.T) {
	data, err := loadDemoSnapshot(diskHuntSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	realOwner, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	invalidCatalog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/demo-cases" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, "not-json")
			return
		}
		realOwner.ServeHTTP(w, r)
	})
	if got, err := collectAuthorizedCrossCaseHunt(context.Background(), invalidCatalog,
		testBearer, "ioc.example"); err == nil || got.Format != "" {
		t.Fatalf("malformed owner discovery accepted: %+v %v", got, err)
	}
	if got, err := collectAuthorizedCrossCaseHunt(context.Background(), realOwner,
		"not-the-authorized-bearer", "ioc.example"); err == nil || got.Format != "" {
		t.Fatalf("wrong bearer accepted: %+v %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	visited := 0
	revoked := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/case-timeline" {
			visited++
			if r.URL.Query().Get("case_id") == "case-c" {
				cancel()
				return
			}
		}
		realOwner.ServeHTTP(w, r)
	})
	got, err := collectAuthorizedCrossCaseHunt(ctx, revoked, testBearer, "ioc.example")
	if !errors.Is(err, context.Canceled) || got.VisibleMatches != 0 ||
		len(got.Results) != 0 || visited != 2 {
		t.Fatalf("revocation leaked partial results: %+v %v visited=%d", got, err, visited)
	}
}

func TestOfflineHuntRanksBoundedAuthorizedFiveMinuteWindows(t *testing.T) {
	windows := map[string]certainWindowCount{
		"2026-10-01T10:00:00Z": {Total: 5, Observed: 3, Inferred: 2},
		"2026-10-01T10:05:00Z": {Total: 9, Observed: 8, Inferred: 1},
		"2026-10-01T10:10:00Z": {Total: 9, Observed: 4, Inferred: 5},
		"2026-10-01T10:15:00Z": {Total: 1, Observed: 1},
	}
	got := topAuthorizedHuntWindows(windows)
	if len(got) != 3 || got[0].StartUTC != "2026-10-01T10:05:00Z" ||
		got[1].StartUTC != "2026-10-01T10:10:00Z" ||
		got[2].StartUTC != "2026-10-01T10:00:00Z" ||
		got[0].Matches != 9 || got[0].Observed != 8 || got[1].Inferred != 5 {
		t.Fatalf("top three Case windows were not ranked deterministically: %+v", got)
	}
	if result := topAuthorizedHuntWindows(nil); len(result) != 0 {
		t.Fatalf("uncertain-only timeline created certain windows: %+v", result)
	}
}

func TestOfflineHuntProvidesAuthorizedTemporalPivotFromDisk(t *testing.T) {
	data, err := loadDemoSnapshot(diskHuntSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := range data.Entries {
		switch data.Entries[i].Value.ID {
		case "event-old":
			data.Entries[i].Value.SourceTime = "2026-10-01T10:01:00Z"
			data.Entries[i].Value.Timezone = "+00:00"
			data.Entries[i].Value.ClockUncertaintyMS = 0
		case "event-c":
			data.Entries[i].Value.SourceTime = "2026-10-01T10:02:00Z"
			data.Entries[i].Value.Timezone = "+00:00"
			data.Entries[i].Value.ClockUncertaintyMS = 0
		}
	}
	blob, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "temporal-hunt.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path, "-hunt-q=ioc.example",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var got localCrossCaseHunt
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.VisibleMatches != 2 || len(got.Results) != 2 {
		t.Fatalf("authorized Case hunt did not complete: %+v", got)
	}
	for _, item := range got.Results {
		if item.Matches != 1 || item.FirstCertainUTC == "" ||
			item.LastCertainUTC != item.FirstCertainUTC || len(item.TopWindows) != 1 {
			t.Fatalf("Case did not have actionable certain-time pivot: %+v", item)
		}
		instant, err := time.Parse(time.RFC3339Nano, item.FirstCertainUTC)
		if err != nil {
			t.Fatal(err)
		}
		wantWindow := instant.UTC().Truncate(5*time.Minute).Format(time.RFC3339)
		if item.TopWindows[0].StartUTC != wantWindow ||
			item.TopWindows[0].Matches != 1 || item.TopWindows[0].Inferred != 1 {
			t.Fatalf("wrong authenticated window aggregate: %+v", item)
		}
	}
	for _, secret := range []string{"restricted", "hidden IOC.EXAMPLE", "event-old", "event-c", "demo-secret", "demo-hidden-source", "case-b"} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("unapproved raw/source/foreign metadata in pivot: %s", output.String())
		}
	}
}
