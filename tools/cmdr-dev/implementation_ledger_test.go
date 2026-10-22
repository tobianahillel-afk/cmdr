package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestValidateImplementationLedgerShape(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion:       1,
		LedgerKind:          "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries: []ImplementationLedgerEntry{{
			Capability:      "CAP-SET-004",
			WorkUnit:        "E9-PILOT-001C-RUNTIME",
			ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
			HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary: "pilot-context-envelope-runtime",
			RuntimeRoot:     "product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err != nil {
		t.Fatal(err)
	}
}

func TestImplementationLedgerRejectsDuplicateCapability(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	entry := ImplementationLedgerEntry{
		Capability:      "CAP-SET-004",
		WorkUnit:        "E9-PILOT-001C-RUNTIME",
		ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
		HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
		RuntimeBoundary: "pilot-context-envelope-runtime",
		RuntimeRoot:     "product-runtime/context-envelope/**",
	}
	ledger := ImplementationLedger{
		SchemaVersion: 1, LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries:             []ImplementationLedgerEntry{entry, entry},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected duplicate capability rejection")
	}
}

func TestImplementationLedgerRejectsPathEscape(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion: 1, LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries: []ImplementationLedgerEntry{{
			Capability: "CAP-SET-004", WorkUnit: "E9-PILOT-001C-RUNTIME",
			ManifestPath:    "../outside.json",
			HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary: "pilot-context-envelope-runtime",
			RuntimeRoot:     "product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected path escape rejection")
	}
}

func TestImplementationLedgerRejectsBaselineDrift(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion: 1, LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: "fedcba9876543210fedcba9876543210fedcba98",
		Entries: []ImplementationLedgerEntry{{
			Capability: "CAP-SET-004", WorkUnit: "E9-PILOT-001C-RUNTIME",
			ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
			HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary: "pilot-context-envelope-runtime",
			RuntimeRoot:     "product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected baseline drift rejection")
	}
}

func TestImplementationLedgerRejectsDuplicateWorkUnit(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion: 1, LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries: []ImplementationLedgerEntry{
			{
				Capability: "CAP-SET-004", WorkUnit: "E9-PILOT-001C-RUNTIME",
				ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
				HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
				RuntimeBoundary: "pilot-context-envelope-runtime",
				RuntimeRoot:     "product-runtime/context-envelope/**",
			},
			{
				Capability: "CAP-SET-005", WorkUnit: "E9-PILOT-001C-RUNTIME",
				ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
				HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
				RuntimeBoundary: "another-runtime",
				RuntimeRoot:     "product-runtime/another/**",
			},
		},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected duplicate work unit rejection")
	}
}

func TestImplementationLedgerRejectsDuplicateRuntimeBoundary(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion: 1, LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries: []ImplementationLedgerEntry{
			{
				Capability: "CAP-SET-004", WorkUnit: "E9-PILOT-001C-RUNTIME",
				ManifestPath:    "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
				HandoffPath:     "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
				RuntimeBoundary: "pilot-context-envelope-runtime",
				RuntimeRoot:     "product-runtime/context-envelope/**",
			},
			{
				Capability: "CAP-SET-005", WorkUnit: "E9-PILOT-OTHER-RUNTIME",
				ManifestPath:    "work/lots/E9-PILOT-OTHER-RUNTIME/manifest.json",
				HandoffPath:     "work/lots/E9-PILOT-OTHER-RUNTIME/HANDOFF.json",
				RuntimeBoundary: "pilot-context-envelope-runtime",
				RuntimeRoot:     "product-runtime/context-envelope/**",
			},
		},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected duplicate runtime boundary rejection")
	}
}

func TestRuntimeLedgerGitProofRejectsUnverifiedRuntimeChange(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	runTestGit(t, root, "init")
	runTestGit(t, root, "config", "user.name", "CMDR Test")
	runTestGit(t, root, "config", "user.email", "cmdr-test@example.invalid")
	writeTestFile(t, root, "product-runtime/context-envelope/envelope.go", "package envelope\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "bounded runtime implementation")
	implementation := testGitOutput(t, root, "rev-parse", "HEAD")
	writeTestFile(t, root, "engineering/unrelated.md", "unchanged runtime\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "validated runtime plus engineering metadata")
	validated := testGitOutput(t, root, "rev-parse", "HEAD")
	const runtimeRoot = "product-runtime/context-envelope/**"
	if err := verifyRuntimeClaimGitEvidence(root, implementation, validated, runtimeRoot); err != nil {
		t.Fatalf("unchanged runtime should be verified: %v", err)
	}
	writeTestFile(t, root, "engineering/another.md", "still unchanged runtime\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "non-runtime change")
	if err := verifyRuntimeClaimGitEvidence(root, implementation, validated, runtimeRoot); err != nil {
		t.Fatalf("unrelated changes must not revoke runtime evidence: %v", err)
	}
	writeTestFile(t, root, "product-runtime/context-envelope/cmd/perf-probe/main.go", "package main\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "new executable runtime after validation")
	if err := verifyRuntimeClaimGitEvidence(root, implementation, validated, runtimeRoot); err == nil || !strings.Contains(err.Error(), "fresh runtime evidence") {
		t.Fatalf("post-validation runtime addition must invalidate old proof: %v", err)
	}
	if err := verifyRuntimeClaimGitEvidence(root, strings.Repeat("0", 39)+"a", validated, runtimeRoot); err == nil {
		t.Fatal("fabricated existing-looking implementation SHA must fail")
	}
	if err := verifyRuntimeClaimGitEvidence(root, implementation, strings.Repeat("f", 40), runtimeRoot); err == nil {
		t.Fatal("nonexistent final validated head must fail")
	}
	if err := verifyRuntimeClaimGitEvidence(root, validated, implementation, runtimeRoot); err == nil {
		t.Fatal("reversed implementation/validation ancestry must fail")
	}
	if err := verifyRuntimeClaimGitEvidence(root, implementation, validated, "product-runtime/../escape/**"); err == nil {
		t.Fatal("runtime root path escape must fail")
	}
}

func TestRuntimeSuccessorAdditiveProofRejectsAnyOldRuntimeMutation(t *testing.T) {
 if _,err:=exec.LookPath("git");err!=nil{t.Skip("git unavailable")}
 root:=t.TempDir()
 runTestGit(t,root,"init")
 runTestGit(t,root,"config","user.name","CMDR Test")
 runTestGit(t,root,"config","user.email","cmdr-test@example.invalid")
 writeTestFile(t,root,"product-runtime/context-envelope/envelope.go","package envelope\n")
 runTestGit(t,root,"add",".")
 runTestGit(t,root,"commit","-m","validated runtime")
 before:=testGitOutput(t,root,"rev-parse","HEAD")
 paths:=[]string{"product-runtime/context-envelope/cmd/perf-probe/main.go",
  "product-runtime/context-envelope/cmd/perf-probe/main_test.go"}
 for _,p:=range paths{writeTestFile(t,root,p,"package main\n")}
 runTestGit(t,root,"add",".")
 runTestGit(t,root,"commit","-m","verified additive performance probe")
 successor:=testGitOutput(t,root,"rev-parse","HEAD")
 const runtime="product-runtime/context-envelope/**"
 if err:=verifyOnlyRecordedRuntimeAdditions(root,before,successor,runtime,paths);err!=nil{t.Fatal(err)}
 if err:=verifyRuntimeClaimGitEvidence(root,before,successor,runtime);err!=nil{t.Fatal(err)}
 if err:=verifyOnlyRecordedRuntimeAdditions(root,before,successor,runtime,paths[:1]);err==nil{t.Fatal("incomplete evidence passed")}
 if err:=verifyOnlyRecordedRuntimeAdditions(root,before,successor,runtime,append(paths,paths[0]));err==nil{t.Fatal("duplicate path passed")}
 if err:=verifyOnlyRecordedRuntimeAdditions(root,before,successor,runtime,[]string{"product-runtime/other/main.go"});err==nil{t.Fatal("external path passed")}
 writeTestFile(t,root,"product-runtime/context-envelope/envelope.go","package envelope\n// tampered\n")
 runTestGit(t,root,"add",".")
 runTestGit(t,root,"commit","-m","unverified root change")
 if err:=verifyRuntimeClaimGitEvidence(root,before,successor,runtime);err==nil{t.Fatal("post-successor mutation passed")}
 if err:=verifyOnlyRecordedRuntimeAdditions(root,before,testGitOutput(t,root,"rev-parse","HEAD"),runtime,paths);err==nil{t.Fatal("original file mutation passed")}
}
