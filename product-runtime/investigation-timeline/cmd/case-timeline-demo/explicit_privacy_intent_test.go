package main

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// An empty shell variable must never turn an explicit private report into a
// successful stdout publication of sensitive, authorized Case/source refs.
// This is a real privacy boundary shared by three independent CLI journeys.
func TestExplicitEmptyPrivateReportNeverFallsBackToStdout(t *testing.T) {
	archive,caseID,sha := savedHeadlessMixedCase(t)
	_,_,_,other,otherCase,otherSHA := buildTwoVerifiedLocalCaseArchives(t)
	workflows:=[]struct{name string;args []string}{
		{"inventory",[]string{"-snapshot",archive,"-principal","analyst-a",
			"-verify-import-snapshot-case",caseID,"-expect-snapshot-sha256",sha}},
		{"headless search",[]string{"-snapshot",archive,"-principal","analyst-a",
			"-hunt-saved-case",caseID,"-expect-snapshot-sha256",sha,
			"-hunt-field","source_ref","-hunt-family","all","-hunt-value","local-jsonl-source-",
			"-hunt-from","2026-10-17T09:59:00Z","-hunt-to","2026-10-17T10:02:00Z"}},
		{"two Case comparison",[]string{"-snapshot",archive,"-principal","analyst-a",
			"-compare-saved-case",caseID,"-expect-snapshot-sha256",sha,
			"-compare-saved-with-snapshot",other,"-compare-saved-with-case",otherCase,
			"-expect-other-snapshot-sha256",otherSHA}},
	}
	for _,tc:=range workflows{
		t.Run(tc.name,func(t *testing.T){
			for _,empty:=range []string{"two-argument","equals-notation"}{
				var args []string
				args=append(args,tc.args...)
				if empty=="two-argument"{args=append(args,"-private-report","")}
				if empty=="equals-notation"{args=append(args,"-private-report=")}
				var stdout bytes.Buffer
				if err:=runDemoCLI(context.Background(),args,&stdout);err==nil || stdout.Len()!=0{
					t.Fatalf("explicit but empty private destination disclosed a Case/source report: %v %q",err,stdout.String())
				}
			}
		})
	}
}

// If an operator explicitly requests byte integrity verification, a missing
// environment variable must not silently become an unverified listener.
// Omission of the optional flag remains backwards-compatible, but explicit
// empty values are rejected before any bearer or source receipt is emitted.
func TestExplicitEmptySnapshotDigestCannotDisableVerification(t *testing.T) {
	archive,caseID,_:=savedHeadlessMixedCase(t)
	for _,extra:=range [][]string{
		{"-resume-import-search-case",caseID,"-expect-snapshot-sha256",""},
		{"-resume-import-search-case",caseID,"-expect-snapshot-sha256="},
		{"-verify-import-snapshot-case",caseID,"-expect-snapshot-sha256",""},
		{"-hunt-saved-case",caseID,"-hunt-field","source_ref",
			"-hunt-value","local-jsonl-source-",
			"-hunt-from","2026-10-17T09:59:00Z","-hunt-to","2026-10-17T10:02:00Z",
			"-expect-snapshot-sha256",""},
	} {
		ctx,cancel:=context.WithTimeout(context.Background(),300*time.Millisecond)
		var out bytes.Buffer
		args:=append([]string{"-principal","analyst-a","-snapshot",archive},extra...)
		err:=runDemoCLI(ctx,args,&out)
		cancel()
		if err==nil || out.Len()!=0{
			t.Fatalf("empty operator digest started an unchecked Case workflow: %v %q",err,out.String())
		}
		for _,secret:=range []string{"Bearer","127.0.0.1:","local-jsonl-entry-"}{
			if strings.Contains(out.String(),secret){t.Fatal("sensitive output escaped before digest validation")}
		}
	}
	// An ordinary explicitly private output still works with a real path.
	// This proves the fix protects only missing input, not the capability.
	_,caseID,sha:=savedHeadlessMixedCase(t)
	outfile:=filepath.Join(t.TempDir(),"private-receipt.json")
	var stdout bytes.Buffer
	err:=runDemoCLI(context.Background(),[]string{
		"-principal","analyst-a","-snapshot",archive,
		"-verify-import-snapshot-case",caseID,
		"-expect-snapshot-sha256",sha,
		"-private-report",outfile,
	},&stdout)
	if err!=nil || stdout.Len()!=0{
		t.Fatalf("explicit valid private output path regressed: %v %q",err,stdout.String())
	}
}
