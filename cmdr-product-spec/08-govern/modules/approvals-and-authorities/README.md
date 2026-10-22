---
id: govern-approvals-and-authorities
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013]
---
# Approvals & Authorities — GOV-1

## Mission

Establish the contextual authority required for an Action Request, determine eligible approvers under scope/expiry/SoD, manage Approval Requests/Approvals and govern delegation, substitution, escalation and emergency paths without administering platform identities or executing the requested action.

## Owned capabilities

- `CAP-GOV-009` — Authority Requirement and Authorization Context Assessment.
- `CAP-GOV-010` — Approver Eligibility and Separation-of-Duties Assessment.
- `CAP-GOV-011` — Approval Request and Approval Record Management.
- `CAP-GOV-012` — Delegation, Substitution and Escalation Governance.
- `CAP-GOV-013` — Emergency and Time-Bounded Approval Governance.

## Ownership boundary

Platform Settings retains administration of users, Roles, groups, tenants/environments and configured identity/authority sources. Security retains the permission, Decision Authority, SoD and step-up policy. Govern evaluates those sources **in the request context** and owns Approval/Decision authority records.

`role ≠ authority`; `configured authority ≠ contextual authority`; `approver candidate ≠ eligible approver ≠ actual approver`.

## Approval boundary

Approval Request ≠ Approval. Approval ≠ Decision. Approval ≠ execution. A Human Gate remains Studio-owned and is neither an Approval nor a Decision. An Automation Run cannot approve itself or become a Response Run.

## SoD

SoD evaluates requester, action owner, candidate approver, authority scope and policy context. Requester self-approval is prohibited whenever applicable SoD requires separation. A simple Role mismatch is not the definition of SoD.

## Delegation/emergency

Delegation is scoped, time-bound and sourced; it does not create a permanent Role. Escalation routes unresolved authority but is not an Approval. Emergency governance requires explicit justification, bounded scope/target/duration, authority, audit and future retrospective-review requirement; urgency alone is not approval.

## Screen

`GOV-AUT-001` remains active. Its detailed screen specification is not rewritten; GOV-1 maps the existing surface to these capabilities and leaves final controls/columns/filters to the later screen phase.

## Dependencies

Security authority/SoD/step-up/emergency sources, Settings identities, Shared Notifications/Trace/Versioning/Linking, Action Request/Policy/Risk context, OPEN-007 and OPEN-013.
