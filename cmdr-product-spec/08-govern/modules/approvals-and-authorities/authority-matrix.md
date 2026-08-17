---
id: govern-approvals-and-authorities-authority-matrix
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013]
---
# Functional Authority Matrix Semantics — GOV-1

## Purpose

Define the dimensions that GOV-1 authority assessment must consider without pretending to be the final RBAC/ABAC/Decision Authority implementation.

## Authority dimensions

| Dimension | Functional use | Must not be interpreted as |
|---|---|---|
| action type / class | identifies authority family and possible step-up | permission grant |
| exact target + included scope | bounds where authority applies | ownership of target |
| tenant / environment | isolation and contextual boundary | cross-tenant access by default |
| business/security impact | may increase required authority | opaque universal score |
| reversibility / rollback context | may alter review/authority requirement | guarantee that rollback exists |
| requester relation | SoD/conflict input | automatic ineligibility without applicable rule |
| configured Role/group | candidate authority source | contextual authority itself |
| explicit delegation | scoped/time-bound authority source | permanent Role/privilege |
| authority scope/restrictions | limits candidate/approver use | technical capability to execute |
| effective period / expiry | temporal validity | deletion after expiry |
| Policy/Exception context | additional constraints/requirements | Decision |
| emergency authority | bounded shortened-path authority | ungoverned bypass |

## Assessment output

`CAP-GOV-009` creates a request-specific Authority Requirement and Authority Context. `CAP-GOV-010` then determines candidate eligibility and SoD. `CAP-GOV-011` records the actual Approval. `CAP-GOV-015` records Decision authority. None of these steps mutates a Settings Role or grants a runtime permission.

## N-of-M

GOV-1 may state that an approval requirement requires multiple independent approvers, but does not define the final quorum algorithm, ordering, cryptographic signing or storage representation.

## Failure behavior

Missing, stale, conflicting, expired or out-of-scope authority is `unknown`, `partial`, `conflicted`, `expired` or `blocked`; it never falls back to “admin role means authorized.”

## Boundary

Platform Settings retains users/Roles/groups administration. Security retains the canonical permission, SoD, step-up and authority policy. Govern owns request-specific authority assessment and authority-bearing Approval/Decision records.
