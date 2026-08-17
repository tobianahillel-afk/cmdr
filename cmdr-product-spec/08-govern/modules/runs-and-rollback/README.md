---
id: govern-runs-and-rollback
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
---
# Runs & Rollback — GOV-2

## Mission

Own the governed execution envelope from an Execution Plan through Response Run, technical execution coordination, verification, rollback/recovery when required and canonical Result. The module controls **authority, bounds, lifecycle and reconciliation** while technical owners execute their own primitives.

## GOV-2 capability families

Planning/readiness/control:
- `CAP-GOV-020` Target Resolution and Execution Readiness Assessment;
- `CAP-GOV-021` Authorization, Decision and Condition Reconciliation;
- `CAP-GOV-022` Response Run Creation and Lifecycle Management;
- `CAP-GOV-023` Response Run Scheduling, Start, Pause, Stop and Cancellation Control.

Execution/runtime/verification:
- `CAP-GOV-024..029` cover step coordination, Studio/Endpoint handoff, runtime status/error reconciliation and verification.

Rollback/result/provenance:
- `CAP-GOV-030..033` cover rollback eligibility/execution/recovery, canonical Result and cross-product provenance/handoff.

## Core distinctions

- Execution Plan ≠ Response Run.
- Response Run created ≠ started.
- scheduled ≠ started.
- start requested ≠ start confirmed.
- running ≠ successful.
- Step success ≠ Run success.
- technical output ≠ canonical Result.
- Automation Run ≠ Response Run.
- Shared Job ≠ Response Run.
- cancel ≠ rollback.
- compensation ≠ rollback.
- rollback plan ≠ rollback execution.
- rollback successful ≠ full recovery guaranteed.
- runtime success ≠ verification success.

## Execution-owner boundary

Govern owns the Response Run and authorized execution intent. Studio, Endpoint Agent and configured provider/integration owners retain their technical execution primitives and raw technical responses. GOV-2 correlates those technical records back to the Run instead of duplicating them.

## Target and scope safety

Every effectful start/resume/retry/rollback path preserves exact Decision version, exact Playbook version, exact target set, allowed/prohibited scope, conditions, expiry, Approval/Exception validity and fresh-enough readiness. Drift never expands the approved target set.

## Secrets

Runs and plans may carry Secret References and restricted parameter metadata only. Raw secret values are resolved only for an authorized executor through Platform Settings-owned mechanisms and are never stored in Govern Run/Result/log/report semantics.

## Verification and rollback

Runtime completion is only a technical input. GOV-2 separately defines expected outcome, verification sources, observed outcome, residual risk and rollback/recovery needs. Rollback has its own eligibility, authority, exact scope, execution and post-rollback verification.

## Screen

`GOV-RUN-001` remains the canonical Runs & Rollback screen. GOV-2 rewrites no detailed screen specification and creates no new Screen ID, wireframe, final button/column/filter, animation or shortcut.

## GOV-3 boundary

GOV-3 later owns detailed Audit Trail and Response Metrics capabilities. GOV-2 emits complete provenance and conceptual metrics inputs but creates no GOV-3 capability contract.