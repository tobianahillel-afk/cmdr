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
