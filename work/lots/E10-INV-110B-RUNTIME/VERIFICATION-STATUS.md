# E10-INV-110B — verified CI and residual release gates (2026-10-16)

## Exact code snapshot after child-timeout hardening

- Repository `tobianahillel-afk/cmdr`, PR #36 (open draft), implementation branch `engineering/autonomous-development-foundation`.
- Runtime fix SHA `2421f8be0a970a44334efb6babab0fd0cb09cd5c`, base `c48c1f489cd26e90c82a73b00420013b327405a5`; changed Go runtime control-plane files and `CHILD-TIMEOUT.md`.
- **Verified Actions status on this exact SHA**: push engineering-kernel run [38844455739](https://github.com/tobianahillel-afk/cmdr/actions/runs/38844455739) SUCCESS; PR engineering-kernel run [38844459758](https://github.com/tobianahillel-afk/cmdr/actions/runs/38844459758) SUCCESS; PR investigation-timeline-runtime run [38844459765](https://github.com/tobianahillel-afk/cmdr/actions/runs/38844459765) SUCCESS. No Timeline dedicated push run was found for this snapshot; do not count one.
- A separate one-time second-context read-only L2 was scheduled, **pinned exactly to this SHA**, with diff `c48c1f48..2421f8be`; its report is not available to this artifact. **L2_EN_ATTENTE**, not reviewer approval.
- The code does not declare `E10-INV-110B-RUNTIME` VERIFIED yet. The GitHub task graph remains READY pending evidence reconciliation and independent risk decisions.

## Genuine historical performance evidence, not this run

- On-demand Actions run [38842748892](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842748892) on **older source SHA** `dae5fe2c64e1eb1befb0cc0130cc42040ba5bfdf` passed. Synthetic 5,000-entry projection: p95 of five 200-operation batch **means** = **4.70007566 ms**; 1,000 operations sampled Go heap maximum **5,147,088 bytes**. Neither is per-request p95 or peak RSS. These numbers **must not** be represented as benchmarks of `2421f8be`.
- An independent JSON artifact was produced under intermediate SHA `ae91af5e51250fec84f0fc075f372400fcd7b5c0` in run [38843953373](https://github.com/tobianahillel-afk/cmdr/actions/runs/38843953373), not at the new fix SHA.
- The latest change adds cancellation and bounded Go child execution; CI verified compilation and tests, **not** a repeat of the on-demand performance workflow on this exact code.

## Unresolved risk gates

1. **Authoritative permissions and source ownership:** projection consumes caller-provided decisions; tests demonstrate in-process safe masking, not proof of actual server-side Case, Shared Timeline Entry, Evidence/Govern/Command permission provider integration. No production authorization claim.
2. **Timeline Entry / Shared boundary:** no mutating Timeline Engine, source-provider index, persistent storage, final UI, class-2 contribution/correction, export, or cross-owner schema. Canonical Product Spec immutable; `OPEN-013` unresolved.
3. **Perf interpretation:** latest CI cannot establish production SLO, request-p95, peak RSS, transport latency, large-scale steady-state behavior, or Windows/macOS validation.
4. **Independent audit:** L2 task created on code SHA `2421f8be`; do not mark approved or assert findings unavailable from that exact report have been resolved.
5. **Next action:** reconcile independent L2 if accessible, decide bounded runtime handoff gates before promoting B, and only then unlock C-NAV per graph. A one-off `performance-benchmark-audit`/deep on current code may be requested by an authorized workflow dispatch if exposed; do not fabricate fresh runs.

## Technology choice

For a maximum 5,000-entry in-memory read-only projection, reuse **Go stdlib stable sorting + first-party permission-result envelope + existing CMDR benchmark control-plane**, not a new event-store/index/backend. Go `exec.CommandContext` imposes subprocess deadlines with no added dependency. Production persistence/index comparisons remain deferred until requirements, ADR and representative benchmarks exist.
