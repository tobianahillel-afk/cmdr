# E10-INV-110B — aggregate Timeline ID + source authority-selector budget

## Verified reasoning at parent SHA `159af9193631e023582db9aec2f2eec40499fdbf`

Both `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots` already bound:
1. the aggregate requested Timeline Entry IDs with tenant, environment and Case ID; and
2. the aggregate linked `SourceRef` values with the same context.

However, the second counter was initialized **before** the Entry IDs were charged to the first counter. Consequently, a 4 MiB authorized Timeline Entry ID followed by a 4 MiB linked source reference was individually acceptable under each counter but exceeded the overall engineering 8 MiB selector budget. The source authority could receive more total selector bytes than intended, before the pure `Project` eventually rejected the oversized combined projection.

The minimal fix reuses `usedIDs`, `usedSourceSelectors` and the existing `boundedSourceSelector` helper: after the complete ID preflight and before **any authority call**, initialize the source budget from the final `usedIDs` counter. Both adapters now count Case/tenant/environment, all validated Timeline Entry IDs and every authorized source reference together. No new policy engine or external dependency; no change to Shared Timeline ownership or canonical Product Spec.

New adversarial tests exercise a single 4 MiB Timeline Entry ID followed by a 4 MiB linked source ref in **both** adapters. The Case and Timeline Entry checks are allowed, but the source request must be rejected without contacting its authority. Existing negative tests cover cumulative source-only overflow and denied/cross-tenant objects. This hardening is not a production authorization-provider implementation and does not mean every upstream provider enforces correct Case membership. `OPEN-013` remains open. Historical on-demand performance measurements do not apply to this SHA.

**Evidence discipline:** check the exact-SHA GitHub Actions status and the independent L2 result separately. Until they are verified, B stays READY and C stays BLOCKED.

## Exact GitHub CI confirmation — code SHA `4b6caad0fbed0c59e8ac6eb51cfa86a0bc326c23`

The four workflow runs were confirmed COMPLETED / SUCCESS on the exact code SHA (2026-10-16):
- [engineering-kernel push #38901811787](https://github.com/tobianahillel-afk/cmdr/actions/runs/38901811787)
- [engineering-kernel PR #38901819158](https://github.com/tobianahillel-afk/cmdr/actions/runs/38901819158)
- [investigation-timeline-runtime push #38901811725](https://github.com/tobianahillel-afk/cmdr/actions/runs/38901811725): Linux Go race/test, cover, vet and Go formatting **PASS**; full Timeline module 96.4% statements (projection package 98.3%, synthetic perf-probe package 91.3%). macOS and Windows native Timeline + subprocess test jobs **SUCCESS**.
- [investigation-timeline-runtime PR #38901819443](https://github.com/tobianahillel-afk/cmdr/actions/runs/38901819443).

Linux synthetic Go `testing.B` measurements from the code-SHA push run: 500 entry projection 465168–532284 ns/op; 5000 entries 5664194–6084098 ns/op across three 100-operation local repetitions. These are **per-batch means, not per-request p95** and do not establish a production 5ms objective. Existing on-demand run measurements from previous commits are not re-attributed to this code SHA. The updated source-selector budget is code- and mock-test-verified; no real authenticated provider integration is included. `E10-INV-110B-RUNTIME` stays READY, `E10-INV-110C-NAV` stays BLOCKED. An independent second-context read-only L2 task was scheduled exactly for this code SHA; the report is **L2_EN_ATTENTE** and this text is not an L2 approval.
