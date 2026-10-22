---
id: govern-policy-gates-gate-evaluation
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-020]
---
# Policy Gate Evaluation — GOV-1 Supporting Contract

## Purpose

Provide the module-level supporting semantics consumed by `CAP-GOV-007`. The canonical functional contract remains the capability file; this document prevents historical `pass/fail` shorthand from being mistaken for final GOV-1 semantics.

## Outcomes

GOV-1 uses exactly these functional evaluation outcomes:

- `pass` — explicit evaluated conditions are satisfied for the bound inputs; **not** proof that the action is safe or authorized;
- `warn` — a non-blocking Policy concern requires visibility/review;
- `block` — the Policy context blocks progression unless governance legitimately resolves the condition; **not** an automatic Decision of rejection;
- `unknown` — required input/evaluation cannot be established; never silently passes;
- `not-applicable` — Policy does not apply to this request/scope/version, with reason.

## Bound evaluation context

Every evaluation is bound to:
- Action Request id/version;
- Policy id/version;
- applicability scope and effective period;
- relevant target/scope references;
- input source references/versions;
- evaluation method or Tool/Tool Call if used;
- outcome/reasons/missing inputs;
- timestamp, reviewer and automation provenance.

## Re-evaluation

A material change to request, target/scope, Policy version or dependent input makes the prior evaluation stale. Re-evaluation creates a new attributable record/diff; it never rewrites historical evidence.

## Boundary

No Policy Evaluation creates Approval, Decision, active exception, Response Run or target effect. Conflict and exception review routes to `CAP-GOV-008`; final authority routes through CAP-GOV-009..015.
