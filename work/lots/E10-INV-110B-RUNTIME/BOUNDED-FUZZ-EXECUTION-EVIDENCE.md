# E10-INV-110B — verified bounded fuzz, platform, and on-demand evidence

## Immutable source and exact GitHub Actions results

- Repository: `tobianahillel-afk/cmdr`; branch `engineering/autonomous-development-foundation`; PR #36 (open draft).
- Source under test: **`0764ee3079eb787c26e04ed6db0639c111323106`**. Parent: `d6c64bd2727478c3f9470c770e4259faa7e45715`. The evidence in this file **must not be attributed to any subsequent documentation commit**.
- GitHub Actions [38848606015](https://github.com/tobianahillel-afk/cmdr/actions/runs/38848606015) — `investigation-timeline-evidence` push — **SUCCESS**. Verified job results:
  - `evidence` job **116804609087**: **SUCCESS**; explicit latency and resource target IDs validated, one selected and executed target each; no cached/reused latency target.
  - `Timeline process tree / macos-latest` job **116804609287**: **SUCCESS**. Executed Go platform smoke tests, including child deadline, descendant cancellation, bounded stdout/stderr and real Go probe.
  - `Timeline process tree / windows-latest` job **116804609336**: **SUCCESS**, Go 1.26.8 windows/amd64. Same targeted tests passed, including descendant termination and diagnostics bounds.
  - `Timeline bounded opaque-source fuzz / Linux` job **116804609348**: **SUCCESS**. Go native fuzz engine, two workers for 25 seconds: **836,692 executions**, five new interesting inputs (nine total), no failing property found. This is a time-bounded campaign, **not** exhaustive validation.
- Engineering kernel push [38848605980](https://github.com/tobianahillel-afk/cmdr/actions/runs/38848605980): **SUCCESS**.
- Engineering kernel PR [38848610466](https://github.com/tobianahillel-afk/cmdr/actions/runs/38848610466): **SUCCESS**.
- Dedicated investigation-timeline-runtime PR [38848610415](https://github.com/tobianahillel-afk/cmdr/actions/runs/38848610415): **SUCCESS**. No dedicated runtime push run claimed for this workflow-only code update.

## Exact-source synthetic benchmark data

The source-`0764ee...` run measured the following in **GitHub Actions hosted Ubuntu**, not production:

- `PERF-TGT-INV-TIMELINE-PROJECTION`: `pass`, **3.524015215 ms**, p95 **across five means of 200 projections each** (not p95 of independent requests), provisional engineering upper bound 5 ms. Five samples, zero reused targets.
- `PERF-DEEP-INV-TIMELINE-PROJECTION`: `pass`, **3,248,576 bytes** maximum periodically sampled Go heap, **1,887,091,600 bytes** total allocations over **1,000** operations, **3,399 ms** execution, about **294.16 operations/s**. The heap value is not process peak RSS.
- Artifact [12024501366](https://github.com/tobianahillel-afk/cmdr/actions/runs/38848606015/artifacts): `investigation-timeline-on-demand-0764ee3079eb787c26e04ed6db0639c111323106-38848606015`, both exact-SHA JSON evidence files, retention until **2026-11-15T04:03:49Z** (unless manually deleted).

## Remaining gates — do not promote B based on this evidence alone

1. Runtime `Project` still consumes **caller-supplied** Case/Timeline/source permission outcomes. No authoritative provider, server endpoint, remote permission evaluation or real Case/Evidence/Govern/Command ownership check is proved.
2. Fuzzing targeted **opaque fields under four access classifications**; it cannot prove freedom from all data leaks, cross-boundary misuse, performance bugs or authorization bypasses.
3. Performance numbers omit network, storage, transport, UI, persistence, index, real user concurrency, request-level latency percentiles and peak RSS. They do not certify a production SLO.
4. The independent L2 was successfully **scheduled** on the actual code commit `0764ee...`, and its report has **not** been retrieved through this repository or task context. State **L2_EN_ATTENTE**, no approval asserted.
5. The work graph still has `E10-INV-110B-RUNTIME` READY, `E10-INV-110C-NAV` BLOCKED; Product Spec and decision `OPEN-013` unchanged. Do not merge to main or mark VERIFIED until the real release/handoff gates have been reconciled.

Reuse-first implementation used Go's native fuzz engine and previously pinned GitHub-hosted CI actions. No third-party fuzz runner, paid API, additional product dependency or storage backend was added.
