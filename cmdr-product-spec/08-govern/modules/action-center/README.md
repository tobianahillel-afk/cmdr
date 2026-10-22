---
id: govern-action-center
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020]
open_decisions: [OPEN-013]
---
# Action Center — GOV-1

## Mission

Provide the focused review workspace for an exact Action Request/version. Action Center assembles source context, scope/target, impact/risk/reversibility, completeness, Policy/authority/Approval projections and Decision preparation while preserving source ownership and performing no target execution.

## Owned GOV-1 capabilities

- `CAP-GOV-004` — Action Request Context, Scope and Target Review.
- `CAP-GOV-005` — Action Impact, Risk and Reversibility Assessment.
- `CAP-GOV-006` — Action Request Completeness and Evidence Context Review.
- `CAP-GOV-014` — Decision Preparation and Review (defined in the Decision lot).

## Consumed review capabilities

- `CAP-GOV-007..008` — Policy outcomes/conflicts/exceptions.
- `CAP-GOV-009..013` — authority/eligibility/Approval/delegation/emergency.
- `CAP-GOV-015..016` — final Decision and execution handoff projections.

## Inputs

Current Action Request/version, Incident/Case/Finding/Evidence references, exact target/scope/time bounds, impact and risk context, Policy outcomes, authority requirements, Approvals, unresolved questions and provenance.

## Outputs

Context/Scope Review, Impact/Risk/Reversibility Assessments, Completeness Review, information requests and a Decision Draft/decision-ready package. None of these outputs is an Approval, Decision, Response Run or Result.

## Distinctions

- target selected ≠ target verified;
- complete ≠ authorized;
- risk score/summary ≠ Decision;
- Evidence context ≠ Evidence qualification;
- Decision Draft ≠ Decision;
- approved authority ≠ executed action.

## AI/no-AI

AI may draft summaries/options only. Deterministic checks, forms, matrices, diffs, Policy/authority viewers and human review provide the full path.

## Screen

`GOV-ACT-001` remains active and unchanged at detailed-screen level. GOV-1 creates no new Screen ID, column, button, filter, wireframe, shortcut or animation specification.

## GOV-2 boundary

Rollback information is reviewed as context/requirement only. Action Center does not execute rollback, start a Response Run or create Result. The final GOV-1 output is an Execution Handoff Package consumed later by GOV-2.
