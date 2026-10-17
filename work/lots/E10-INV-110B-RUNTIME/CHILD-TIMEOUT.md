# E10-INV-110B — bounded Timeline child execution

## Verified code change under review

A previous L2 identified that the generic `runPerformanceProcess` uses `exec.Command` without `CommandContext`; a deep-audit parent timeout could return while the Timeline child continued using resources. This increment addresses **only the Timeline first-party probe path**, reusing Go `context`, `exec.CommandContext` and the existing CMDR probe format. No vendor/paid API, storage backend, new authorization provider or policy exception is introduced.

- Timeline Go build subprocess: 90-second timeout, cleans temporary directory on failure.
- Timeline sample subprocess: 25-second hard child deadline, bounded by the parent deep-audit `ctx` when applicable.
- Cancelled/expired contexts return wrapped `context.Canceled` or `context.DeadlineExceeded`; no synthetic result is returned.
- Tests exercise a blocked child binary, an already-cancelled context, normal `go version`, and cancelled probe dispatch. Windows/macOS/Linux `exec.CommandContext` is supported by the Go standard library, though this is not yet evidence that all OS variants were exercised in CI.
- Generic subprocess helpers for other runtime boundaries remain unchanged. B remains not VERIFIED until independent review, exact-HEAD CI and remaining scope questions are addressed.

The already stored 4.70007566 ms batch-mean p95 / sampled heap result belongs solely to older SHA `dae5fe2c64e1eb1befb0cc0130cc42040ba5bfdf`, run 38842748892. This patch is **not** a new measured benchmark or a real ACL-provider check.
