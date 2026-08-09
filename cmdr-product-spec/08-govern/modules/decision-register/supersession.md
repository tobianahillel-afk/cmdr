---
id: govern-decision-register-supersession
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-015, REQ-PROD-020]
open_decisions: [OPEN-013, OPEN-015]
---
# Decision Supersession and Expiration — GOV-1 Supporting Contract

## Supersession

A Decision is never edited in place to change historical authority. A replacement Decision references the superseded Decision, explains why a new Decision is required, identifies the effective replacement and preserves the complete old request/version/authority/Approval/Policy context.

Supersession may be required by:
- material Action Request/version change;
- target/scope change;
- Policy or authority change;
- expired/revoked Approval or delegation;
- corrected factual context;
- changed Decision conditions/time bounds;
- explicit Decision Maker reconsideration under valid authority.

## Expiration

Expiration means the Decision’s authority/time bound is no longer current. It **does not delete** the Decision, rationale, approvals, provenance or previously prepared handoff package. An expired Decision cannot be treated as current execution authority.

## Handoff effects

When a Decision expires or is superseded:
- any associated Execution Handoff Package becomes `stale`, `expired-with-decision` or `superseded`;
- future GOV-2 must not create/start a Response Run from that package without a current authoritative Decision;
- historical links remain navigable under permission.

## Comparison

Decision Register supports a conceptual version comparison that exposes changed request version, targets/scope, authority, Approvals, Policy context, rationale, conditions, expiry and rollback/verification requirements. Final diff UI belongs to the later screen phase.

## Invariants

- reject ≠ delete;
- defer ≠ reject;
- request-more-information ≠ reject;
- supersede ≠ mutate history;
- expiration ≠ deletion;
- a new Decision requires current authority and prerequisites rather than inheriting authority silently from the old Decision.
