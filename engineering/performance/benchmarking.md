# Representative benchmark and regression engine

`performance-benchmark-audit` executes only registered targets selected for the current stage, environment and changed paths.

## No executable configuration

A performance target contains a `workload_key`, not a command. `benchmark-policy.json` maps each workload key to a compiled `handler_key`. Unknown handlers fail closed.

The current compiled workload is an engineering-control-plane metadata audit. There are currently no registered performance targets, so the real repository result is deterministic `not-applicable`.

## Measurement identity

Every executed result binds:

- exact Git source SHA;
- target ID and canonical target digest;
- environment ID and canonical environment digest;
- workload policy digest;
- Go toolchain version;
- explicit warmup count and sample count through the workload policy.

A relative baseline is reusable only when target, environment, workload and toolchain identity still match exactly.

## Regression semantics

Absolute budgets always apply. Upper-bound metrics fail above budget. Throughput lower-bound metrics fail below budget.

Relative regression is evaluated only when the target requests it and its environment permits it. Missing or stale baselines fail closed. Improvements clamp relative regression to zero rather than appearing as negative regressions.

Shared GitHub CI is intentionally configured as non-authoritative for relative regression; it can support coarse absolute smoke budgets only.
