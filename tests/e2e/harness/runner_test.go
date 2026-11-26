package harness

import (
	"errors"
	"strings"
	"testing"
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
}
