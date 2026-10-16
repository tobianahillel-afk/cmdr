# E10-INV-110B — aggregate Timeline ID + source authority-selector budget

## Verified reasoning at parent SHA `159af9193631e023582db9aec2f2eec40499fdbf`

Both `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots` already bound:
1. the aggregate requested Timeline Entry IDs with tenant, environment and Case ID; and
2. the aggregate linked `SourceRef` values with the same context.

However, the second counter was initialized **before** the Entry IDs were charged to the first counter. Consequently, a 4 MiB authorized Timeline Entry ID followed by a 4 MiB linked source reference was individually acceptable under each counter but exceeded the overall engineering 8 MiB selector budget. The source authority could receive more total selector bytes than intended, before the pure `Project` eventually rejected the oversized combined projection.

The minimal fix reuses `usedIDs`, `usedSourceSelectors` and the existing `boundedSourceSelector` helper: after the complete ID preflight and before **any authority call**, initialize the source budget from the final `usedIDs` counter. Both adapters now count Case/tenant/environment, all validated Timeline Entry IDs and every authorized source reference together. No new policy engine or external dependency; no change to Shared Timeline ownership or canonical Product Spec.

New adversarial tests exercise a single 4 MiB Timeline Entry ID followed by a 4 MiB linked source ref in **both** adapters. The Case and Timeline Entry checks are allowed, but the source request must be rejected without contacting its authority. Existing negative tests cover cumulative source-only overflow and denied/cross-tenant objects. This hardening is not a production authorization-provider implementation and does not mean every upstream provider enforces correct Case membership. `OPEN-013` remains open. Historical on-demand performance measurements do not apply to this SHA.

**Evidence discipline:** check the exact-SHA GitHub Actions status and the independent L2 result separately. Until they are verified, B stays READY and C stays BLOCKED.
