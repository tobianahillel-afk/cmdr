package harness

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

// RaceCheckStatus represents the outcome classification of a race-enabled test run.
type RaceCheckStatus string

const (
	RaceCheckPass       RaceCheckStatus = "PASS"
	RaceCheckFail       RaceCheckStatus = "FAIL"
	RaceCheckUnverified RaceCheckStatus = "UNVERIFIED"
)

// RaceCheckOutcome records the evaluated status and human-readable reason for a race check.
type RaceCheckOutcome struct {
	Status RaceCheckStatus
	Reason string
}

// EvaluateRaceTestOutcome classifies a race-enabled `go test -race` run.
// Fails closed on data races, test/build/0-test failures, and ANY failure on non-Windows platforms (e.g., Linux CI).
// Returns RaceCheckUnverified ONLY on Windows when verified toolchain/cgo/ASR incompatibility prevented execution.
func EvaluateRaceTestOutcome(goos string, res VerifiedGoTestResult, runErr error) RaceCheckOutcome {
	if strings.Contains(res.Output, "WARNING: DATA RACE") ||
		strings.Contains(res.RawStdout, "WARNING: DATA RACE") ||
		strings.Contains(res.RawStderr, "WARNING: DATA RACE") {
		return RaceCheckOutcome{
			Status: RaceCheckFail,
			Reason: "data race detected (WARNING: DATA RACE)",
		}
	}

	if runErr == nil && res.Passed {
		return RaceCheckOutcome{
			Status: RaceCheckPass,
			Reason: "verified execution passed with zero data races",
		}
	}

	errText := ""
	if runErr != nil {
		errText = runErr.Error()
	}
	combined := strings.ToLower(strings.Join([]string{res.Output, res.RawStdout, res.RawStderr, errText}, "\n"))

	canBeWindowsEnvUnverified := strings.EqualFold(strings.TrimSpace(goos), "windows") &&
		!res.CompilationFailed &&
		len(res.ExecutedTests) == 0 &&
		len(res.FailedTests) == 0 &&
		len(res.PassedTests) == 0 &&
		!strings.Contains(combined, "panic:") &&
		!strings.Contains(combined, "fatal error:") &&
		!strings.Contains(combined, "test timed out") &&
		!strings.Contains(combined, "deadline exceeded")

	if canBeWindowsEnvUnverified {
		toolchainSigs := []string{
			"-race is not supported",
			"race detector requires cgo",
			"-race requires cgo",
			"cgo_enabled=0",
			`c compiler "gcc" not found`,
			`exec: "gcc": executable file not found`,
		}
		for _, sig := range toolchainSigs {
			if strings.Contains(combined, sig) {
				return RaceCheckOutcome{
					Status: RaceCheckUnverified,
					Reason: fmt.Sprintf("Windows race toolchain unavailable (%s)", sig),
				}
			}
		}

		asrSigs := []string{
			"access is denied",
			"accès refusé",
			"file contains a virus or potentially unwanted software",
		}
		for _, sig := range asrSigs {
			if strings.Contains(combined, sig) {
				return RaceCheckOutcome{
					Status: RaceCheckUnverified,
					Reason: fmt.Sprintf("Windows Defender ASR blocked test binary execution (%s)", sig),
				}
			}
		}
	}

	reason := "race test execution failed"
	if runErr != nil {
		reason = runErr.Error()
	}
	return RaceCheckOutcome{
		Status: RaceCheckFail,
		Reason: reason,
	}
}

// RunGoTestWithRace runs `go test -json -count=1 -race` for expectedTests in moduleRelPath
// and evaluates the outcome using EvaluateRaceTestOutcome(runtime.GOOS, ...).
func RunGoTestWithRace(ctx context.Context, moduleRelPath, expectedPkg string, expectedTests []string) (VerifiedGoTestResult, RaceCheckOutcome) {
	var extraArgs []string
	extraArgs = append(extraArgs, "-race")
	if len(expectedTests) == 1 {
		extraArgs = append(extraArgs, "-run", "^"+expectedTests[0]+"$")
	} else if len(expectedTests) > 1 {
		extraArgs = append(extraArgs, "-run", "^("+strings.Join(expectedTests, "|")+")$")
	}
	res, err := RunGoTestVerified(ctx, moduleRelPath, expectedPkg, expectedTests, extraArgs...)
	return res, EvaluateRaceTestOutcome(runtime.GOOS, res, err)
}
