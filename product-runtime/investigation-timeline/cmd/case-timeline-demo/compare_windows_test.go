package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Prove temporal co-occurrence using a real operator disk snapshot and
// source-authorized paginated HTTP owners, not invented cross-Case matches.
func TestCLIComparisonCorrelatesOnlyCertainAuthorizedFiveMinuteWindows(t *testing.T) {
	original := twoCaseComparisonSnapshot(t)
	data, err := loadDemoSnapshot(original)
	if err != nil { t.Fatal(err) }
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "2026-10-01T10:02:05Z"
		}
	}
	blob, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "authorized-correlation.json")
	if err := os.WriteFile(path, blob, 0o600); err != nil { t.Fatal(err) }
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot="+path,
		"-compare-case=case-a", "-with-case=case-c",
	}, &output); err != nil { t.Fatal(err) }
	var result localCaseComparison
	if err := json.Unmarshal(output.Bytes(), &result); err != nil { t.Fatal(err) }
	if len(result.SharedCertainWindows) != 1 { t.Fatalf("expected one precise shared window: %+v",result.SharedCertainWindows) }
	window := result.SharedCertainWindows[0]
	if window.StartUTC != "2026-10-01T10:00:00Z" ||
		window.FirstEntries != 206 || window.SecondEntries != 1 ||
		window.FirstObserved != 103 || window.FirstInferred != 103 ||
		window.SecondObserved != 0 || window.SecondInferred != 1 {
		t.Fatalf("incorrect authorized certain-window correlation: %+v",window)
	}
	if strings.Contains(output.String(),"SENSITIVE-PREVIEW") ||
		strings.Contains(output.String(),"demo-hidden-source") ||
		strings.Contains(output.String(),"extra-000") ||
		strings.Contains(output.String(),"Bearer ") {
		t.Fatalf("comparison included raw or forbidden data: %s", output.String())
	}
}

// Unknown timestamps cannot yield a certain five-minute co-occurrence.
func TestCLIComparisonDoesNotCorrelateUncertainSourceTime(t *testing.T) {
	original := twoCaseComparisonSnapshot(t)
	data, err := loadDemoSnapshot(original)
	if err != nil { t.Fatal(err) }
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-c" {
			data.Entries[i].Value.SourceTime = "unknown"
			data.Entries[i].Value.Timezone = "unknown"
		}
	}
	// The owner must still have a valid source and reader; only its clock is uncertain.
	data.Sources = append([]timeline.MemoryTimelineSource(nil),data.Sources...)
	blob, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(),"uncertain-compare.json")
	if err := os.WriteFile(path,blob,0o600); err != nil { t.Fatal(err) }
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(),[]string{
		"-principal=analyst-a","-snapshot="+path,
		"-compare-case=case-a","-with-case=case-c",
	},&out); err != nil {t.Fatal(err)}
	var got localCaseComparison
	if err := json.Unmarshal(out.Bytes(),&got); err != nil {t.Fatal(err)}
	if len(got.SharedCertainWindows)!=0 || got.Second.UncertainEntries != 1 {
		t.Fatalf("uncertain clock promoted to shared timeline window: %+v",got)
	}
}
