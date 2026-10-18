package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// The three real analyst journeys MUST publish the complete authorized
// receipt into a private new-only file without streaming any source refs or
// token into stdout. They deliberately reuse the same source owner and CLI.
func TestPrivateInvestigationReportsIntegrateAllThreeOfflineJourneys(t *testing.T) {
	caseFile, caseID, fingerprint := savedHeadlessMixedCase(t)
	aFile, aCase, aHash, bFile, bCase, bHash := buildTwoVerifiedLocalCaseArchives(t)
	runs := []struct {
		name string
		args []string
		required string
	}{
		{"offline inventory", []string{
			"-principal", "analyst-a", "-snapshot", caseFile,
			"-verify-import-snapshot-case", caseID,
			"-expect-snapshot-sha256", fingerprint,
		}, "cmdr-offline-authorized-case-verification-v1"},
		{"headless mixed Event Search", []string{
			"-principal", "analyst-a", "-snapshot", caseFile,
			"-hunt-saved-case", caseID, "-expect-snapshot-sha256", fingerprint,
			"-hunt-field", "source_ref", "-hunt-family", "all",
			"-hunt-value", "local-jsonl-source-",
			"-hunt-from", "2026-10-17T09:59:00Z", "-hunt-to", "2026-10-17T10:02:00Z",
		}, "cmdr-local-headless-authorized-hunt-v1"},
		{"independent Case comparison", []string{
			"-principal", "analyst-a", "-snapshot", aFile,
			"-compare-saved-case", aCase, "-expect-snapshot-sha256", aHash,
			"-compare-saved-with-snapshot", bFile,
			"-compare-saved-with-case", bCase,
			"-expect-other-snapshot-sha256", bHash,
		}, "cmdr-local-verified-saved-case-comparison-v1"},
	}
	for _, tc := range runs {
		t.Run(tc.name, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "private-investigation-report.json")
			args := append(append([]string(nil), tc.args...), "-private-report", dest)
			var stdout bytes.Buffer
			if err := runDemoCLI(context.Background(), args, &stdout); err != nil {
				t.Fatalf("private authorized receipt failed: %v", err)
			}
			if stdout.Len() != 0 {
				t.Fatalf("private report streamed sensitive output to stdout: %q", stdout.String())
			}
			file, err := os.ReadFile(dest)
			if err != nil { t.Fatal(err) }
			if !json.Valid(file) || !bytes.Contains(file, []byte(tc.required)) {
				t.Fatalf("private report is not the complete authorized receipt: %q", file)
			}
			if runtime.GOOS != "windows" {
				info, err := os.Stat(dest)
				if err != nil { t.Fatal(err) }
				if info.Mode().Perm() != 0600 {
					t.Fatalf("private report permissions are not 0600: %s", info.Mode().Perm())
				}
			}
			for _, leaked := range []string{
				"PRIVATE-WAZUH-ONE", "PRIVATE-WAZUH-TWO",
				"PRIVATE-GENERIC-NEAR", "PowerShell suspicious endpoint activity",
				"Bearer token", "Local analyst viewer:", "127.0.0.1:",
				"\"preview\"", "\"audit_trail\"",
			} {
				if bytes.Contains(file, []byte(leaked)) {
					t.Fatalf("private result included unrequested source data: %s", leaked)
				}
			}
			// The publisher is no-overwrite, even for a valid authorized
			// report: existing forensic output may never be truncated.
			again := bytes.Buffer{}
			if err := runDemoCLI(context.Background(), args, &again); err == nil || again.Len() != 0 {
				t.Fatalf("private report overwritten or duplicate success: %q", again.String())
			}
			unchanged, err := os.ReadFile(dest)
			if err != nil || !bytes.Equal(unchanged, file) {
				t.Fatalf("existing report changed on repeated publication: %v", err)
			}
		})
	}
}

func TestPrivateInvestigationReportsRejectDeniedOwnerBeforeFilesystemSideEffect(t *testing.T) {
	caseFile, caseID, fingerprint := savedHeadlessMixedCase(t)
	aFile, aCase, aHash, bFile, bCase, bHash := buildTwoVerifiedLocalCaseArchives(t)
	denied := []struct{
		name string
		args []string
	}{
		{"inventory wrong principal", []string{
			"-principal", "foreign-analyst", "-snapshot", caseFile,
			"-verify-import-snapshot-case", caseID,
			"-expect-snapshot-sha256", fingerprint,
		}},
		{"hunt changed fingerprint", []string{
			"-principal", "analyst-a", "-snapshot", caseFile,
			"-hunt-saved-case", caseID,
			"-expect-snapshot-sha256", strings.Repeat("0", 64),
			"-hunt-value", "local-jsonl-source-", "-hunt-field", "source_ref",
			"-hunt-from", "2026-10-17T09:59:00Z", "-hunt-to", "2026-10-17T10:02:00Z",
		}},
		{"comparison denied second Case", []string{
			"-principal", "analyst-a", "-snapshot", aFile,
			"-compare-saved-case", aCase, "-expect-snapshot-sha256", aHash,
			"-compare-saved-with-snapshot", bFile,
			"-compare-saved-with-case", "not-authorized",
			"-expect-other-snapshot-sha256", bHash,
		}},
		{"comparison denied second digest", []string{
			"-principal", "analyst-a", "-snapshot", aFile,
			"-compare-saved-case", aCase, "-expect-snapshot-sha256", aHash,
			"-compare-saved-with-snapshot", bFile,
			"-compare-saved-with-case", bCase,
			"-expect-other-snapshot-sha256", strings.Repeat("0",64),
		}},
	}
	for _,tc:=range denied {
		t.Run(tc.name,func(t *testing.T) {
			dir:=t.TempDir()
			dest:=filepath.Join(dir,"must-not-exist.json")
			args:=append(append([]string(nil),tc.args...),"-private-report",dest)
			var stdout bytes.Buffer
			if err:=runDemoCLI(context.Background(),args,&stdout);err==nil || stdout.Len()!=0 {
				t.Fatalf("denied Case produced private report or stdout: %q",stdout.String())
			}
			if _,err:=os.Lstat(dest);!os.IsNotExist(err) {
				t.Fatal("denied Case created an output file")
			}
			files,err:=os.ReadDir(dir)
			if err!=nil || len(files)!=0 {t.Fatalf("denied request left partial forensic files: %+v %v",files,err)}
		})
	}
}

func TestPrivateInvestigationReportCannotBeSmuggledIntoOtherModes(t *testing.T) {
	file, caseID, fingerprint := savedHeadlessMixedCase(t)
	snapshot := snapshotFixture(t)
	for _,args:=range [][]string{
		{"-snapshot", snapshot, "-report-case", caseID},
		{"-snapshot", file, "-resume-import-search-case", caseID},
		{"-snapshot", file},
		{"-snapshot", file, "-verify-import-snapshot-case", caseID},
		{"-principal", "analyst-a", "-snapshot", file,
		 "-verify-import-snapshot-case", caseID, "-expect-snapshot-sha256", fingerprint,
		 "-serve-import"},
	} {
		dir:=t.TempDir()
		dest:=filepath.Join(dir,"should-never-exist.json")
		params:=append(append([]string(nil),args...),"-private-report",dest)
		var stdout bytes.Buffer
		if err:=runDemoCLI(context.Background(),params,&stdout);err==nil || stdout.Len()!=0 {
			t.Fatalf("unrelated private report mode accepted: %v",args)
		}
		if _,err:=os.Lstat(dest);!os.IsNotExist(err) {
			t.Fatal("unexpected report created by invalid CLI combination")
		}
	}
}
