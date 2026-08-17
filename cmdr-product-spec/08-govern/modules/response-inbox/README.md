---
id: govern-response-inbox
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-UX-008]
open_decisions: [OPEN-010, OPEN-013]
---
# Response Inbox — GOV-1

## Mission

Receive submitted Action Requests into a Govern-only review queue, establish intake/lifecycle state, expose blockers/deadlines/source and route the exact request/version into Action Center.

## Owned capabilities

- `CAP-GOV-001` — Govern Intake and Preconditions.
- `CAP-GOV-002` — Govern Response Inbox and Request Queue Management.
- `CAP-GOV-003` — Action Request Lifecycle Management.

## Boundary with Command Work Queue

Response Inbox is **not** Command Work Queue. It contains governed Action Requests/review projections, not general Incidents and Tasks. Shared provides generic search/filter/assignment mechanisms; Govern owns the meaning of request states, deadlines, blockers and Govern review assignment.

## Functional inputs

Action Request/version, source product/return origin, requester, target/scope summary, intake state, review blockers, deadlines/expiration and source projection status.

## Functional outputs

Govern Review Context, queue projection, assignment/reassignment events, information-request routing, lifecycle transitions and exact Action Center navigation context.

## Key states

`received`, `incomplete`, `ready-for-review`, `blocked`, `information-required`, `policy-review`, `authority-review`, `approval-pending`, `decision-ready`, `decided`, `withdrawn`, `expired`, `superseded`.

These are processing/queue projections; the complete persisted object state machine remains future.

## AI and automation

Deterministic validation, filters, deadlines and assignment rules work without AI. AI may summarize blockers or suggest an assignee but never approves or decides.

## Screen

Existing screen `GOV-INB-001` remains active and is not rewritten by GOV-1. Its detailed screen design is deferred; `../../screen-capability-map.md` maps it to these capabilities.

## Dependencies

Action Request, Settings identity/tenant/environment projections, Shared Search/Notifications/Assignments/Inspector/Linking/Trace, source-product handoffs, OPEN-010 and OPEN-013.

## Acceptance

A request can be received, identified as incomplete, returned for information, assigned and reopened in Action Center while preserving exact version/return origin and without appearing as a duplicate general Command work item or creating authority.
