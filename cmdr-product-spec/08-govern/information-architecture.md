---
id: 08-govern-information-architecture
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-015, REQ-UX-001, REQ-UX-008, REQ-UX-009, REQ-UX-010]
open_decisions: [OPEN-010]
---
# Govern Information Architecture

## Canonical modules

Govern retains exactly nine existing product modules:

1. `response-inbox`
2. `action-center`
3. `decision-register`
4. `policy-gates`
5. `approvals-and-authorities`
6. `playbooks`
7. `runs-and-rollback`
8. `audit-trail`
9. `response-metrics`

GOV-1 specifies modules 1–5 functionally. Modules 6–7 are read-only boundaries for future GOV-2. Modules 8–9 are read-only boundaries for future GOV-3.

## Workspace responsibilities

### Response Inbox
Govern-only queue of received/submitted Action Requests and their review states. It is not Command Work Queue. It owns queue semantics such as Govern assignment, expiry/deadline context, blockers and information-required state, while Shared provides generic filtering/search/assignment mechanisms.

### Action Center
Focused review workspace for one Action Request/version. It composes context/scope/target, impact/risk/reversibility, completeness, Policy/authority/Approval projections and Decision preparation without taking ownership of source objects.

### Decision Register
Immutable/supersedable authority history: Decision disposition, rationale, conditions, time bounds, authority/Approval references and future execution-handoff status. Rejection/expiry never deletes history.

### Policy Gates
Govern view of applicable Policy versions, evaluation outcomes, conflicts, warnings, blocks, unknowns and Exception Candidates. It does not define a final Policy authoring engine or automatic Decision engine.

### Approvals & Authorities
Contextual authority, approver eligibility, SoD, Approval Requests/records, delegation/escalation and emergency/time-bound authorization. Platform Settings still administers users/roles/groups and configured authority sources.

## Cross-module flow

`Response Inbox → Action Center → Policy Gates / Approvals & Authorities → Action Center Decision Preparation → Decision Register → Execution Handoff Package → future GOV-2`.

A reviewer may move backward for information, conflict resolution or authority clarification. Context preservation must return to the exact Action Request/version and previous queue selection.

## Shared capability use

Govern consumes Shared Search, Linking, Notifications, Jobs, Activity, Trace, Reporting, Collaboration, Comments, Assignments, Inspector, Versioning and Recovery. These are mechanisms, not extra Govern modules.

## Screen boundary

Existing screen IDs remain stable. GOV-1 does not define final visual composition, columns, filters, buttons, animations or shortcuts. `screen-capability-map.md` is the only GOV-1 mapping update required at this level.

## Acceptance

A user can identify which of the nine modules owns each Govern activity, which modules are in GOV-1, and where the execution boundary begins without inferring a new roadmap phase or a new screen.
