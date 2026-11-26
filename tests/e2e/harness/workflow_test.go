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
        continue-on-error: true
        run: go run . mainline-health-gate --root ../..
`
		res := ValidateHealthGateWorkflow(wf)
		if res.Valid || !res.HasHealthGateStep || !res.HasContinueOnError {
			t.Fatalf("expected HasContinueOnError=true and Valid=false, got %+v", res)
		}
	})

	t.Run("RejectsIfFalse", func(t *testing.T) {
		for _, expr := range []string{"false", "${{ false }}", "!true"} {
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
