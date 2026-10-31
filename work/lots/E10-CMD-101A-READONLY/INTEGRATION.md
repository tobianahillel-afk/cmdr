# E10-CMD-101A-READONLY — Command Work Queue class-0 runtime boundary

## Canonical scope and owners

CAP-CMD-101 §§1–27, Command / Incidents and Work Queue, REQ-OBJ-001/012, REQ-PROD-013, REQ-UX-008/009. Product Canonical `cmdr-product-spec/06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md`, exact six system keys in `../saved-views.md`. Command owns Incident and Task. Shared owns saved-view storage, versioning and permission enforcement; other product references never grant access. The new work unit `E10-CMD-101A-READONLY` belongs to epic `E10-CMD-101`. The **Git parent**, frozen PR #36 HEAD `ea1770b48a4bc1b2e850a7ffab15d5839f9e9e49`, is not the engineering parent.

## What actually works

First-party dependency-free Go function `Project(ctx,owner,Request)` accepts an authoritative host-owned `Owner.List` and **per-object** `Owner.CanRead` callback. Each request requires an explicit single Tenant, Environment and principal; unknown view, wildcard/malformed scope, missing provider, excessive dataset (>5000), duplicate IDs or malformed authorized item fail closed. For each candidate in the requested scope, source authorization occurs **before** any user-visible owner/status/ID is projected. A provider ACL error aborts all output, never returning already processed rows. Foreign scoped objects are skipped without asking the provider to authorize a foreign ID.

Supports exactly the catalog's six `view` keys, as different configurations of **one** projected queue: `all`, `incidents`, `tasks`, `unassigned`, `sla-risk`, `my-work`. A missing authoritative SLA policy cannot silently mark an item low-risk: `sla-risk` returns `partial-sources` and only known-at-risk authorized items. My Work uses only authorized owner/follower projections. Results sort by actual RFC3339 timestamp descending, Kind/ID tiebreakers, maximum first 50 rows, no raw source, no external reference or mutation.

## Threat model and exclusions

The caller-supplied provider is part of the trusted host. There is **no** independent IAM implementation, HTTP listener or production data adapter. The caller MUST enforce current tenant, environment and object-level read permission through `Owner.CanRead`; returning true blindly is a hostile/insecure integration. Work Queue is not a Case Queue or full Incident Detail/Task coordination. No linked Case/Decision/Run source projection, Shared Saved View persisted custom configurations, final columns, pagination cursor, return/Inspector frontend, audit sink or action mutation. SLA truth must come from the separate Command SLA owner. The open mutation policy OPEN-013 is not resolved; only class-0 read-only work is included.

### Acceptance

- Six views select the exact expected authorized Incident/Task IDs, stable same-route view projection and deterministic first 50 rows.
- Host authorization rejects denied Case/Task without leaking ID or source, no cross-tenant/environment query to authorization provider.
- Unknown view (`team-load`), missing/wildcard Tenant or invalid scope, owner errors, duplicate/malformed data or canceled context all fail closed.
- Missing SLA risk source is visibly partial, no AI required for any system view.
- Exact HEAD GitHub Go test/vet Ubuntu, macOS and Windows, engineering-kernel manifest/architecture/90%-security gate. No API payante or third-party product code.

This is **CONTRAT/PROTOTYPE with runnable tests**, not an integrated operator-facing Work Queue, accepted capability, native GitHub approval or production deployment. Next material increments: independently authorized actual Incident/Task providers; Shared Saved View permission-aware application with canonical columns; native single-route Workspace with selection/Inspector and exact return. Each must have its own owner and READY dependency.


## E10-CMD-101A exact CI reconciliation — deep-performance prerequisite

On the frozen parent PR #153 HEAD `573a233707de5f70c95fbde95672392f8a1c740c`, both `command-work-queue-runtime` and `investigation-timeline-runtime` succeeded, but `engineering-kernel` #39341894165 FAILED at `CHK-DEEP-PERFORMANCE`: the newly registered `PERF-TGT-COMMAND-WORK-QUEUE-PROJECTION` had no corresponding deep target. This is an actual integration/quality omission, **not** a security coverage finding; do not change the mandatory 90% critical security floor or misreport the CI as green.

An independent child branch fixes the same existing **E10-CMD-101A-READONLY** obligation, not a new canonical capability or imaginary product parent: the real Command Go runtime probe now supports an opt-in, bounded `resource` mode recording sampled Go heap and total runtime allocations. It preserves the original latency benchmark exactly. The existing engineering runner executes the compiled Command probe through `builtin-command-work-queue-projection-v1`, independently validates operations, mode, elapsed time and heap against a 128MiB target, and computes actual throughput. `engineering/performance/deep-performance-policy.json` registers one real `on-demand` resource target referring to the pre-existing performance target. **PR validation checks the deep policy; it does not execute a deep workload at every push.** Go behavioral tests execute the real resource probe and deep handler. No proxy measurement from Case Queue or generic cmdr-dev metadata, no third-party dependencies, no production SLO claim.

The PR #153 parent HEAD stays immutable for its existing review. The child is NOT mergeable until exact-head GitHub workflows and native reviews succeed, and does not by itself deliver the host Incident/Task adapter, Shared Saved Views, browser workspace, durable queue, or full CAP-CMD-101 acceptance. The epic remains DECOMPOSED and this sublot is not promoted to VERIFIED without complete evidence. No duplicate L2 is required solely because of this targeted CI prerequisite repair; a materially changed/risks-bearing later integration receives its own L2.
