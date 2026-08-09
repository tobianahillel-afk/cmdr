---
id: govern-policy-gates
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
---
# Policy Gates — GOV-1

## Mission

Determine which Policy versions are applicable to an Action Request, evaluate explicit conditions, preserve `pass/warn/block/unknown/not-applicable`, surface conflicts and govern Exception Candidates without turning any Policy output into an automatic Decision.

## Owned capabilities

- `CAP-GOV-007` — Policy Applicability and Evaluation.
- `CAP-GOV-008` — Policy Conflict, Exception and Waiver Assessment.

## Scope

GOV-1 covers candidate discovery, applicability, outcome/reasons, missing inputs, comparison/replay, conflict review and bounded Exception Candidates. It does **not** select a final Policy engine/language, define a production rule-authoring system, silently activate exceptions or execute Policy-driven target changes.

## Functional outcomes

Canonical GOV-1 Policy Evaluation outcomes:
- `pass`;
- `warn`;
- `block`;
- `unknown`;
- `not-applicable`.

`pass` does not mean safe; `block` is not an automatic rejection Decision; `unknown` is never silently coerced to pass; applicability is distinct from satisfaction.

## Conflict / exception rules

A Policy conflict preserves every involved Policy/version/outcome and reason. An Exception Candidate must have justification, exact target/scope, time bounds/expiry, risks and required authority. Candidate ≠ active exception; exception ≠ Policy deletion; emergency ≠ bypass.

## AI/no-AI

AI may suggest candidate Policies, summarize outcomes/conflicts or draft an Exception Candidate. Catalogs, deterministic checks, diffs, matrices and human review provide the complete non-AI path.

## Screen

`GOV-POL-001` remains active; its generic historical references to authoring/activation are not interpreted as GOV-1 implementation ownership. Detailed screen rewrite is deferred.

## Dependencies

Action Request/context/risk/completeness, Security permission/authority rules, Settings scope projections, Shared Versioning/Trace/Linking, Authority/Approval/Decision capabilities, OPEN-013.
