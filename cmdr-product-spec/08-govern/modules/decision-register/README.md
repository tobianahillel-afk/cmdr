---
id: govern-decision-register
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Decision Register — GOV-1

## Mission

Maintain the immutable/supersedable authority history of Govern Decisions and the provenance of their execution-handoff packages, without turning an approved Decision into proof of execution.

## Owned capabilities

- `CAP-GOV-015` — Decision Recording, Disposition, Conditions and Expiration.
- `CAP-GOV-016` — Govern Decision Provenance and Execution Handoff.

`CAP-GOV-014` Decision Preparation is owned by Action Center and feeds this module.

## Decision semantics

Decision dispositions in GOV-1 are `approve`, `approve-with-conditions`, `reject`, `defer`, `request-more-information`, `cancel` and `supersede`. Historical object-state wording remains subject to future Objects reconciliation; GOV-1 does not silently rewrite the canonical object schema.

A Decision records Decision Maker, authority context, request version, target/scope, rationale, uncertainty, conditions, permitted/prohibited scope, time/start constraints, expiry, rollback/verification requirements, Approvals and dissent.

## Invariants

- Approval ≠ Decision.
- Decision ≠ recommendation, Policy outcome, risk score or AI suggestion.
- Decision ≠ Response Run or Result.
- `approve` ≠ execute.
- `approve-with-conditions` ≠ unrestricted approval.
- Decision condition ≠ Policy.
- expiration ≠ deletion.
- supersession preserves the prior Decision and provenance.

## Execution handoff

`CAP-GOV-016` prepares a no-effect **Execution Handoff Package** containing the exact authorized action, targets, scope, limits, conditions, expiry, Approval/Decision context and rollback/verification requirements. It is not a Response Run and cannot mutate a target. Future GOV-2 is the first capability-specification lot allowed to define execution/run/verification/rollback behavior.

## Screen

`GOV-DEC-001` remains active and is not rewritten at detailed-screen level. Search/filter/export use Shared mechanisms and preserve permission-aware provenance.

## Dependencies

Action Request and all GOV-1 review records; Security Decision Authority/SoD/step-up; Shared Versioning/Trace/Linking/Reporting/Export; Studio provenance; future GOV-2/GOV-3; OPEN-007/013/015.
