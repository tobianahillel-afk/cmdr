package harness

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestVerifyGoTestJSON_RejectsZeroExecutedTests(t *testing.T) {
	zeroTestsJSON := strings.Join([]string{
		`{"Action":"start","Package":"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev"}`,
		`{"Action":"output","Package":"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev","Output":"testing: warning: no tests to run\n"}`,
		`{"Action":"pass","Package":"github.com/tobianahillel-afk/cmdr/tools/cmdr-dev"}`,
	}, "\n")

	res1, err1 := VerifyGoTestJSON([]byte(zeroTestsJSON), "github.com/tobianahillel-afk/cmdr/tools/cmdr-dev", []string{"TestIsRepoRoot"})
	if err1 == nil || res1.Passed {
		t.Fatalf("expected failure when expected test did not execute, got res=%+v err=%v", res1, err1)
	}

	res2, err2 := VerifyGoTestJSON([]byte(zeroTestsJSON), "github.com/tobianahillel-afk/cmdr/tools/cmdr-dev", nil)
	if err2 == nil || res2.Passed {
		t.Fatalf("expected failure when zero tests executed with empty expectedTests, got res=%+v err=%v", res2, err2)
	}
}

func TestVerifyGoTestJSON_RejectsMissingSubtest(t *testing.T) {
	const pkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"
	parentOnlyJSON := strings.Join([]string{
		`{"Action":"start","Package":"` + pkg + `"}`,
		`{"Action":"run","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation"}`,
		`{"Action":"pass","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation"}`,
		`{"Action":"pass","Package":"` + pkg + `"}`,
	}, "\n")

	res, err := VerifyGoTestJSON([]byte(parentOnlyJSON), pkg, []string{
		"TestChallenge_DeepCopy_NonMutation",
		"TestChallenge_DeepCopy_NonMutation/Modifying_Inspected_Contract_Cannot_Mutate_Stored_Contract",
	})
	if err == nil || res.Passed {
		t.Fatalf("expected failure when required subtest did not run, got res=%+v err=%v", res, err)
	}
}

func TestVerifyGoTestJSON_RejectsFailedOrSkippedTest(t *testing.T) {
	const pkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"

	failedJSON := strings.Join([]string{
		`{"Action":"run","Package":"` + pkg + `","Test":"TestA"}`,
		`{"Action":"fail","Package":"` + pkg + `","Test":"TestA"}`,
		`{"Action":"fail","Package":"` + pkg + `"}`,
	}, "\n")
	resFail, errFail := VerifyGoTestJSON([]byte(failedJSON), pkg, []string{"TestA"})
	if errFail == nil || resFail.Passed || len(resFail.FailedTests) != 1 {
		t.Fatalf("expected failure on Action=fail, got res=%+v err=%v", resFail, errFail)
	}

	skippedJSON := strings.Join([]string{
		`{"Action":"run","Package":"` + pkg + `","Test":"TestA"}`,
		`{"Action":"skip","Package":"` + pkg + `","Test":"TestA"}`,
		`{"Action":"pass","Package":"` + pkg + `"}`,
	}, "\n")
	resSkip, errSkip := VerifyGoTestJSON([]byte(skippedJSON), pkg, []string{"TestA"})
	if errSkip == nil || resSkip.Passed || len(resSkip.SkippedTests) != 1 {
		t.Fatalf("expected rejection on Action=skip for expected test, got res=%+v err=%v", resSkip, errSkip)
	}

	buildFailJSON := strings.Join([]string{
		`{"ImportPath":"` + pkg + `","Action":"build-output","Output":"# ` + pkg + `\nsyntax error\n"}`,
		`{"ImportPath":"` + pkg + `","Action":"build-fail","FailedBuild":"` + pkg + `"}`,
	}, "\n")
	resBuild, errBuild := VerifyGoTestJSON([]byte(buildFailJSON), pkg, []string{"TestA"})
	if errBuild == nil || resBuild.Passed || !resBuild.CompilationFailed {
		t.Fatalf("expected CompilationFailed=true on build-fail, got res=%+v err=%v", resBuild, errBuild)
	}
}

func TestVerifyGoTestJSON_AcceptsVerifiedRunAndPass(t *testing.T) {
	const pkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"
	validJSON := strings.Join([]string{
		`{"Action":"start","Package":"` + pkg + `"}`,
		`{"Action":"run","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation"}`,
		`{"Action":"run","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation/SubA"}`,
		`{"Action":"output","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation/SubA","Output":"--- PASS: TestChallenge_DeepCopy_NonMutation/SubA (0.00s)\n"}`,
		`{"Action":"pass","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation/SubA"}`,
		`{"Action":"pass","Package":"` + pkg + `","Test":"TestChallenge_DeepCopy_NonMutation"}`,
		`{"Action":"pass","Package":"` + pkg + `"}`,
	}, "\n")

	res, err := VerifyGoTestJSON([]byte(validJSON), pkg, []string{
		"TestChallenge_DeepCopy_NonMutation",
		"TestChallenge_DeepCopy_NonMutation/SubA",
	})
	if err != nil || !res.Passed {
		t.Fatalf("expected verified pass, got res=%+v err=%v", res, err)
	}
}

func TestVerifyGoTestJSON_RejectsSkippedChildSubtestWhenParentOnlyListed(t *testing.T) {
	const pkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"
	skippedSubJSON := strings.Join([]string{
		`{"Action":"start","Package":"` + pkg + `"}`,
		`{"Action":"run","Package":"` + pkg + `","Test":"TestParent"}`,
		`{"Action":"run","Package":"` + pkg + `","Test":"TestParent/Sub1"}`,
		`{"Action":"skip","Package":"` + pkg + `","Test":"TestParent/Sub1"}`,
		`{"Action":"pass","Package":"` + pkg + `","Test":"TestParent"}`,
		`{"Action":"pass","Package":"` + pkg + `"}`,
	}, "\n")

	res, err := VerifyGoTestJSON([]byte(skippedSubJSON), pkg, []string{"TestParent"})
	if err == nil || res.Passed || !strings.Contains(err.Error(), "TestParent/Sub1") {
		t.Fatalf("expected rejection citing TestParent/Sub1 when child subtest is skipped, got res=%+v err=%v", res, err)
	}
}

func TestEvaluateRaceTestOutcome(t *testing.T) {
	t.Run("DataRaceDetectedFails", func(t *testing.T) {
		res := VerifiedGoTestResult{Output: "WARNING: DATA RACE\nWrite at 0x00c0001"}
		out := EvaluateRaceTestOutcome("linux", res, errors.New("exit status 66"))
		if out.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail on WARNING: DATA RACE, got %+v", out)
		}
	})

	t.Run("NonRaceCommandOrZeroTestFailureFails", func(t *testing.T) {
		res := VerifiedGoTestResult{Passed: false, Output: "zero tests executed"}
		for _, goos := range []string{"linux", "windows"} {
			out := EvaluateRaceTestOutcome(goos, res, errors.New("zero tests executed"))
			if out.Status != RaceCheckFail {
				t.Fatalf("expected RaceCheckFail on %s for non-race failure, got %+v", goos, out)
			}
		}
	})

	t.Run("CompilationErrorOnWindowsWithAccessDeniedTextStillFails", func(t *testing.T) {
		res := VerifiedGoTestResult{
			Passed:            false,
			CompilationFailed: true,
			Output:            "syntax error: access is denied string in comment",
		}
		out := EvaluateRaceTestOutcome("windows", res, errors.New("compile failed"))
		if out.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail when CompilationFailed=true, got %+v", out)
		}
	})

	t.Run("LinuxFailsClosedEvenOnCgoError", func(t *testing.T) {
		res := VerifiedGoTestResult{Passed: false, RawStderr: "go: -race requires cgo; enable cgo by setting CGO_ENABLED=1"}
		out := EvaluateRaceTestOutcome("linux", res, errors.New("exit status 2"))
		if out.Status != RaceCheckFail {
			t.Fatalf("expected linux to fail closed on cgo error, got %+v", out)
		}
	})

	t.Run("WindowsCgoOrASRMarksUnverifiedNeverPass", func(t *testing.T) {
		resCgo := VerifiedGoTestResult{Passed: false, RawStderr: "go: -race requires cgo; enable cgo by setting CGO_ENABLED=1"}
		outCgo := EvaluateRaceTestOutcome("windows", resCgo, errors.New("exit status 2"))
		if outCgo.Status != RaceCheckUnverified {
			t.Fatalf("expected RaceCheckUnverified on Windows without cgo, got %+v", outCgo)
		}

		resASR := VerifiedGoTestResult{Passed: false, RawStderr: "fork/exec C:\\tmp\\pkg.test.exe: Access is denied."}
		outASR := EvaluateRaceTestOutcome("windows", resASR, errors.New("exit status 1"))
		if outASR.Status != RaceCheckUnverified {
			t.Fatalf("expected RaceCheckUnverified on Windows ASR block, got %+v", outASR)
		}
	})

	t.Run("VerifiedSuccessPasses", func(t *testing.T) {
		res := VerifiedGoTestResult{Passed: true, PassedTests: []string{"TestAdversarial_ConcurrentStoreStress"}}
		out := EvaluateRaceTestOutcome("linux", res, nil)
		if out.Status != RaceCheckPass {
			t.Fatalf("expected RaceCheckPass, got %+v", out)
		}
	})

	t.Run("WindowsFailedOrExecutedTestWithAccessDeniedFails", func(t *testing.T) {
		resFailed := VerifiedGoTestResult{
			Passed:        false,
			ExecutedTests: []string{"TestAuthReject"},
			FailedTests:   []string{"TestAuthReject"},
			Output:        "FAIL TestAuthReject: access is denied",
		}
		outFailed := EvaluateRaceTestOutcome("windows", resFailed, errors.New("exit status 1"))
		if outFailed.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail when FailedTests non-empty despite access is denied, got %+v", outFailed)
		}

		resPassed := VerifiedGoTestResult{
			Passed:      false,
			PassedTests: []string{"TestA"},
			RawStderr:   "access is denied",
		}
		outPassed := EvaluateRaceTestOutcome("windows", resPassed, errors.New("exit status 1"))
		if outPassed.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail when PassedTests non-empty despite access is denied, got %+v", outPassed)
		}
	})

	t.Run("WindowsPanicWithAccessDeniedFails", func(t *testing.T) {
		resWithTests := VerifiedGoTestResult{
			Passed:        false,
			ExecutedTests: []string{"TestPanic"},
			Output:        "panic: access is denied [recovered]",
		}
		if out := EvaluateRaceTestOutcome("windows", resWithTests, errors.New("exit status 2")); out.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail on panic with ExecutedTests, got %+v", out)
		}

		resNoTests := VerifiedGoTestResult{
			Passed: false,
			Output: "panic: access is denied [recovered]",
		}
		if out := EvaluateRaceTestOutcome("windows", resNoTests, errors.New("exit status 2")); out.Status != RaceCheckFail {
			t.Fatalf("expected RaceCheckFail on panic with empty ExecutedTests, got %+v", out)
		}
	})

	t.Run("WindowsTimeoutWithAccessDeniedFails", func(t *testing.T) {
		for _, timeoutText := range []string{"test timed out after 10m0s: access is denied", "context deadline exceeded: access is denied"} {
			res := VerifiedGoTestResult{
				Passed: false,
				Output: timeoutText,
			}
			out := EvaluateRaceTestOutcome("windows", res, errors.New("exit status 1"))
			if out.Status != RaceCheckFail {
				t.Fatalf("expected RaceCheckFail for %q, got %+v", timeoutText, out)
			}
		}
	})
}

func TestRunCampaignInRoot_AggregatesFailuresWithoutEarlyExit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tmpRoot := t.TempDir()
	mods := []struct {
		rel  string
		body string
	}{
		{
			rel: "mod-a-pass",
			body: `package modapass
import "testing"
func TestPassA(t *testing.T) {}
`,
		},
		{
			rel: "mod-b-fail",
			body: `package modbfail
import "testing"
func TestFailB(t *testing.T) { t.Fatal("intentional failure") }
`,
		},
		{
			rel: "mod-c-pass",
			body: `package modcpass
import "testing"
func TestPassC(t *testing.T) {}
`,
		},
	}

	var modNames []string
	for _, m := range mods {
		modDir := filepath.Join(tmpRoot, m.rel)
		if err := os.MkdirAll(modDir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", m.rel, err)
		}
		if err := os.WriteFile(filepath.Join(modDir, "go.mod"), []byte("module example.com/"+m.rel+"\n\ngo 1.26.11\n"), 0644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}
		if err := os.WriteFile(filepath.Join(modDir, "mod_test.go"), []byte(m.body), 0644); err != nil {
			t.Fatalf("write mod_test.go: %v", err)
		}
		modNames = append(modNames, m.rel)
	}

	results, allPassed := RunCampaignInRoot(ctx, tmpRoot, modNames)
	if len(results) != 3 {
		t.Fatalf("expected 3 aggregated module results without early exit, got %d", len(results))
	}
	if !results[0].VetPassed || !results[1].VetPassed || !results[2].VetPassed {
		t.Fatalf("expected go vet to pass on all 3 modules, got %+v", results)
	}
	if runtime.GOOS == "windows" {
		low := strings.ToLower(results[0].TestOutput)
		if strings.Contains(low, "access is denied") || strings.Contains(low, "accès refusé") {
			t.Skipf("Windows Defender ASR blocked go test child binary execution (vet and 3-module aggregation verified): %s", results[0].TestOutput)
		}
	}
	if allPassed {
		t.Fatalf("expected allPassed=false when mod-b-fail fails")
	}
	if !results[0].TestPassed {
		t.Fatalf("expected mod-a-pass to pass test, got %+v", results[0])
	}
	if results[1].TestPassed {
		t.Fatalf("expected mod-b-fail to fail tests, got %+v", results[1])
	}
	if !results[2].TestPassed {
		t.Fatalf("expected mod-c-pass to be evaluated and pass after mod-b-fail failure, got %+v", results[2])
	}
}
