# E10-CMD-101J-DIAGNOSTICS — Authorized SLA Risk view diagnostic correction

## Canonical scope and four owners

- **Product:** Command CAP-CMD-101 Unified Work Queue, six fixed system views and exact status filter. CAP-CMD-105 owns SLA projection meaning and its external SLA policy/contract data. The unimplemented Class-2 mutation/pause/follow-up obligations require OPEN-013 and are **excluded**.
- **Engineering:** Existing E10-CMD-101 epic, child task E10-CMD-101J-DIAGNOSTICS (depends on E10-CMD-101I-SELECTION IMPLEMENTED). Only the child is IMPLEMENTED, not all of CAP-CMD-101 or CAP-CMD-105 ACCEPTED.
- **Code:** first-party Go `product-runtime/work-queue/queue.go`; existing `Project` source and owner. No new permission engine, CLI mode, runtime dependency, source SLA policy or adapter.
- **Git:** child of PR #169 frozen exact HEAD `b78bb5c5a456883e44718debf897789c24aee86a`, different business parent CAP-CMD-109; no change to ancestor L2-audited SHA.

## Reproduced R2 gap and bounded repair

Before this change, `Project` inside `case SLARisk` evaluated `if item.SLAAtRisk == nil { partialSLA = true }` for **every owner-authorized item in the tenant** before the selected exact `req.StatusFilter` was applied. Consequently `view=sla-risk&status=open` could report `State=partial-sources` and `sla-policy-unavailable` **only because an independently readable CLOSED Task had no SLA policy**, despite that Task being outside the filter and unshown.

This is a product-state misleading diagnostic / potential in-tenant filtered-cohort information side channel, not a cross-tenant leak: Owner.CanRead already ran, and foreign items were skipped. The correction adds **one existing status equality guard at the decision that sets partialSLA**. Missing policy for an authorized item *inside* the status filter still produces partial-sources and no invented SLA-risk membership, satisfying canonical CAP-CMD-101 §§14-16/22 and CAP-CMD-105 §§11/14-15/22. Neither permission flow, source truth, membership criteria, sort/page cursor, nor six views change.

## Behavioral proof

A single deterministic first-party Go integration test constructs four real owner candidate types: authorized OPEN at-risk Incident; authorized CLOSED unknown-SLA Task; same-tenant DENIED OPEN unknown-SLA Task; foreign-tenant unknown-SLA Incident. Verifies:
- SLA Risk + status=open => exactly authorized Incident, `ready`, **zero** `sla-policy-unavailable` warnings;
- SLA Risk + status=closed => no inferred risk rows, but **partial-sources** and exact diagnostic for still-relevant missing policy;
- SLA Risk without filter => historical partial-sources preserved;
- Tasks + status=open => no authorized matching rows and no spurious SLA degradation;
- no Owner.CanRead call on foreign tenant identity, denied object cannot affect results.

## Integration limits and release status

This code runs under the synthetic first-party Owner and SessionResolver; real SLA policy providers, contractual policy ownership, automatic refresh, Task follow-up, pause/restart, IAM, live SIEM data, cross-product source links and durable production acceptance are NOT claimed. No changes to `cmdr-product-spec/**` or mandatory CI thresholds. One-sentence behavior correction with one meaningful regression test, not a new engine or broad audit.

Exact-head GitHub `command-work-queue-runtime` Go Linux/macOS/Windows, `engineering-kernel` including CHK-MANIFESTS/COMPLEXITY/ARCHITECTURE/SECURITY-TESTS changed-critical >=90%, and SAST must pass. Until verified and required native approvals obtained: DRAFT, UNMERGED; parent #169 independent L2 already executed, no duplicate L2 for this same-model R2 correction unless materially different risk emerges. Current work cache is non-authoritative; any new READY work must be selected by reconciling canonical obligations and actual frozen Git parents.
