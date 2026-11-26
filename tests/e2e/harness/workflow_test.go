package harness

import "testing"

const validTestWorkflowYAML = `name: engineering-kernel-full

on:
  workflow_dispatch:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"

jobs:
  full-kernel:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Enforce Mainline Health Gate
        working-directory: tools/cmdr-dev
        run: go run . mainline-health-gate --root ../..
`

func TestValidateHealthGateWorkflow(t *testing.T) {
	t.Run("AcceptsValidWorkflow", func(t *testing.T) {
		res := ValidateHealthGateWorkflow(validTestWorkflowYAML)
		if !res.Valid || !res.HasHealthGateStep || !res.HasPushMain || !res.HasPRMain || res.HasContinueOnError || res.HasDisabledIf {
			t.Fatalf("expected valid workflow, got %+v", res)
		}
	})

	t.Run("AcceptsMultilineRunBlock", func(t *testing.T) {
		wf := `name: full
on:
  push:
    branches: [main]
  pull_request:
    branches: ["main"]
jobs:
  full-kernel:
    steps:
      - name: Health gate multiline
        run: |
          cd tools/cmdr-dev
          go run . mainline-health-gate --root ../..
`
		res := ValidateHealthGateWorkflow(wf)
		if !res.Valid || !res.HasHealthGateStep {
			t.Fatalf("expected multiline run block workflow to be valid, got %+v", res)
		}
	})

	t.Run("RejectsMissingHealthGateStep", func(t *testing.T) {
		wf := `name: missing
on:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Test
        run: go test ./...
`
		res := ValidateHealthGateWorkflow(wf)
		if res.Valid || res.HasHealthGateStep {
			t.Fatalf("expected rejection for missing mainline-health-gate step, got %+v", res)
		}
	})

	t.Run("RejectsCommentedHealthGateStep", func(t *testing.T) {
		wf := `name: commented
on:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: mainline-health-gate decoy name
        # run: go run . mainline-health-gate --root ../..
        run: |
          # go run . mainline-health-gate --root ../..
          go run . doctor --root ../..
`
		res := ValidateHealthGateWorkflow(wf)
		if res.Valid || res.HasHealthGateStep {
			t.Fatalf("expected rejection for commented/doctor-only health gate step, got %+v", res)
		}
	})

	t.Run("RejectsDoctorOnlyFallback", func(t *testing.T) {
		wf := `name: doctor-only
on:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Doctor
        run: go run . doctor --root ../..
`
		res := ValidateHealthGateWorkflow(wf)
		if res.Valid || res.HasHealthGateStep {
			t.Fatalf("expected rejection for doctor-only workflow, got %+v", res)
		}
	})

	t.Run("RejectsContinueOnErrorTrue", func(t *testing.T) {
		for _, val := range []string{"true", "${{ true }}", `"${{ true }}"`, "${{ ! false }}"} {
			wf := `name: non-blocking
on:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Enforce Mainline Health Gate
        continue-on-error: ` + val + `
        run: go run . mainline-health-gate --root ../..
`
			res := ValidateHealthGateWorkflow(wf)
			if res.Valid || !res.HasHealthGateStep || !res.HasContinueOnError {
				t.Fatalf("expected HasContinueOnError=true and Valid=false for continue-on-error: %s, got %+v", val, res)
			}
		}
	})

	t.Run("RejectsIfFalse", func(t *testing.T) {
		for _, expr := range []string{
			"false",
			"${{ false }}",
			"!true",
			"${{ ! true }}",
			"! true",
			"${{ false && github.ref == 'refs/heads/main' }}",
			"github.event_name == 'push' && false",
		} {
			wf := `name: disabled
on:
  push:
    branches:
      - "main"
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Enforce Mainline Health Gate
        if: ` + expr + `
        run: go run . mainline-health-gate --root ../..
`
			res := ValidateHealthGateWorkflow(wf)
			if res.Valid || !res.HasDisabledIf {
				t.Fatalf("expected HasDisabledIf=true for if: %s, got %+v", expr, res)
			}
		}
	})

	t.Run("RejectsJobLevelIfFalseAndContinueOnError", func(t *testing.T) {
		cases := []struct {
			name          string
			jobExtra      string
			wantDisabled  bool
			wantContinueE bool
		}{
			{
				name:         "JobIfFalse",
				jobExtra:     "    if: false",
				wantDisabled: true,
			},
			{
				name:         "JobIfSpacedNotTrue",
				jobExtra:     "    if: ${{ ! true }}",
				wantDisabled: true,
			},
			{
				name:          "JobContinueOnErrorTrue",
				jobExtra:      "    continue-on-error: true",
				wantContinueE: true,
			},
			{
				name:          "JobContinueOnErrorExprTrue",
				jobExtra:      "    continue-on-error: ${{ true }}",
				wantContinueE: true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				wf := `name: job-bypass
on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
jobs:
  full-kernel:
    runs-on: ubuntu-latest
` + tc.jobExtra + `
    steps:
      - name: Enforce Mainline Health Gate
        run: go run . mainline-health-gate --root ../..
`
				res := ValidateHealthGateWorkflow(wf)
				if res.Valid {
					t.Fatalf("expected Valid=false for %s, got %+v", tc.name, res)
				}
				if tc.wantDisabled && !res.HasDisabledIf {
					t.Fatalf("expected HasDisabledIf=true for %s, got %+v", tc.name, res)
				}
				if tc.wantContinueE && !res.HasContinueOnError {
					t.Fatalf("expected HasContinueOnError=true for %s, got %+v", tc.name, res)
				}
			})
		}
	})

	t.Run("RejectsStepInWrongJob", func(t *testing.T) {
		wf := `name: wrong-job
on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
jobs:
  decoy-unrequired-job:
    runs-on: ubuntu-latest
    steps:
      - name: Decoy Gate
        run: go run . mainline-health-gate --root ../..
  full-kernel:
    runs-on: ubuntu-latest
    steps:
      - name: Tests
        run: go test ./...
`
		res := ValidateHealthGateWorkflow(wf)
		if res.Valid || res.HasHealthGateStep {
			t.Fatalf("expected Valid=false and HasHealthGateStep=false when step is in wrong job, got %+v", res)
		}
	})

	t.Run("AcceptsWorkflowWhenUnrelatedJobHasIfFalse", func(t *testing.T) {
		wf := `name: unrelated-disabled
on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
jobs:
  helper-job:
    runs-on: ubuntu-latest
    if: false
    continue-on-error: true
    steps:
      - name: Helper
        run: echo helper
  full-kernel:
    runs-on: ubuntu-latest
    steps:
      - name: Enforce Mainline Health Gate
        run: go run . mainline-health-gate --root ../..
`
		res := ValidateHealthGateWorkflow(wf)
		if !res.Valid || !res.HasHealthGateStep || res.HasDisabledIf || res.HasContinueOnError {
			t.Fatalf("expected valid workflow when only unrelated helper-job is disabled, got %+v", res)
		}
	})

	t.Run("RejectsEchoDecoyInvocation", func(t *testing.T) {
		for _, decoy := range []string{
			`echo "mainline-health-gate"`,
			`echo "go run . mainline-health-gate"`,
		} {
			wf := `name: echo-decoy
on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]
jobs:
  full-kernel:
    runs-on: ubuntu-latest
    steps:
      - name: Decoy
        run: ` + decoy + `
`
			res := ValidateHealthGateWorkflow(wf)
			if res.Valid || res.HasHealthGateStep {
				t.Fatalf("expected echo decoy %s to be rejected, got %+v", decoy, res)
			}
		}
	})

	t.Run("RejectsMissingMainTriggers", func(t *testing.T) {
		wfNoPush := `name: no-push
on:
  pull_request:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Enforce Mainline Health Gate
        run: go run . mainline-health-gate --root ../..
`
		resNoPush := ValidateHealthGateWorkflow(wfNoPush)
		if resNoPush.Valid || resNoPush.HasPushMain {
			t.Fatalf("expected rejection when push main trigger is missing, got %+v", resNoPush)
		}

		wfNoPR := `name: no-pr
on:
  push:
    branches:
      - "main"
jobs:
  full-kernel:
    steps:
      - name: Enforce Mainline Health Gate
        run: go run . mainline-health-gate --root ../..
`
		resNoPR := ValidateHealthGateWorkflow(wfNoPR)
		if resNoPR.Valid || resNoPR.HasPRMain {
			t.Fatalf("expected rejection when pull_request main trigger is missing, got %+v", resNoPR)
		}
	})
}

func TestValidateWorkflowRequiredJob(t *testing.T) {
	if err := ValidateWorkflowRequiredJob(validTestWorkflowYAML, "full-kernel"); err != nil {
		t.Fatalf("expected validTestWorkflowYAML to contain full-kernel job, got %v", err)
	}
	if err := ValidateWorkflowRequiredJob(validTestWorkflowYAML, "kernel"); err == nil {
		t.Fatalf("expected error when kernel job is absent")
	}
	commented := `name: commented
on:
  push:
    branches: ["main"]
jobs:
  # full-kernel:
  optional-check:
    runs-on: ubuntu-latest
`
	if err := ValidateWorkflowRequiredJob(commented, "full-kernel"); err == nil {
		t.Fatalf("expected error when full-kernel job is commented out")
	}
}

func TestValidateWorkflowSyntax(t *testing.T) {
	if err := ValidateWorkflowSyntax(validTestWorkflowYAML); err != nil {
		t.Fatalf("expected validTestWorkflowYAML to pass syntax check, got %v", err)
	}

	malformed := `name: Broken Workflow
on: [push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Broken Indentation
     run: echo bad
`
	if err := ValidateWorkflowSyntax(malformed); err == nil {
		t.Fatalf("expected unclosed '[' and misaligned indentation to fail syntax validation")
	}

	misalignedOnly := `name: Misaligned
on:
  push:
    branches: ["main"]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Broken Indentation
     run: echo bad
`
	if err := ValidateWorkflowSyntax(misalignedOnly); err == nil {
		t.Fatalf("expected misaligned indentation to fail syntax validation")
	}

	tabYAML := "name: Tab\non:\n\tpush:\n"
	if err := ValidateWorkflowSyntax(tabYAML); err == nil {
		t.Fatalf("expected tab YAML to fail syntax validation")
	}

	unclosedQuote := `name: "unclosed
on: push
`
	if err := ValidateWorkflowSyntax(unclosedQuote); err == nil {
		t.Fatalf("expected unclosed quote to fail syntax validation")
	}
}
