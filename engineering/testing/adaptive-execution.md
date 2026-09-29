# Adaptive validation executor

`cmdr-dev validation-run` executes the prerequisite-closed plan produced by the E3 planner.

The executor is deliberately closed:

- every `executor_key` is matched by a compiled switch in `cmdr-dev`;
- unknown keys fail;
- no manifest, catalog entry or plan can introduce a shell command;
- external processes are fixed program/argument vectors (`gofmt`, `go vet`, `go test`) and do not use a shell;
- product/spec/architecture checks call the existing in-process engine functions directly.

Four checks are marked `PREFLIGHT_SATISFIED` rather than recursively executed:

- `git-changes`: the validated change file already exists before planning;
- `impact`: computed while planning;
- `validation-plan`: computed and validated before execution;
- `validation-run`: represents the current executor process.

All other selected checks run. Execution accounting must satisfy:

`executed + preflight_satisfied == selected`.

During the first E3-IMPACT-002B commit, adaptive execution runs before the legacy full kernel. The legacy kernel remains authoritative so the two paths can be compared on the same push and PR.


## Parity proof and rollout

Adaptive/full parity was proven on the same branch head `8205a798070aaf3aabf3ca4cce724a5b1356a093`:

- push run `35798358793`: adaptive PASS, then legacy full kernel PASS;
- pull-request run `35798364603`: adaptive PASS, then legacy full kernel PASS;
- adaptive summary: 21 selected checks, 17 executed checks, 4 preflight-satisfied checks;
- strict changes execute every mandatory catalog check.

After that proof, `.github/workflows/engineering-kernel.yml` is the ordinary adaptive push/PR gate.

The independent static fallback is preserved as `.github/workflows/engineering-kernel-full.yml`. It contains the pre-rollout explicit full sequence rather than delegating to the adaptive executor. On this foundation branch it runs when that fallback workflow itself changes; once present on the default integration branch it is also available through `workflow_dispatch`.

Any change to either workflow is classified as CI control-plane/high risk, so the adaptive path escalates to `strict` and executes the complete mandatory catalog before accepting the CI change.
