# Response Inbox

## Objective

Prioritise submitted response requests and route them to the correct authority.

## Scope

This specification owns the page-local behaviour of **Response Inbox** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Response request
- Incident
- Case
- Finding
- Principal
- Authority

## Features

- Pending decision queue
- Priority, age and deadline
- Requested action and target
- Requester and evidence completeness
- Required authority and policy status
- Quick review panel

## UX and interactions

- Opening detail preserves queue state
- Missing evidence and policy blockers are explicit
- Queue distinguishes request urgency from incident severity
- Review never executes the action

## Permissions

`response.approve` to decide within authority; `response.request` does not grant approval.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Response request and Decision states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Findings, incidents, policy gates, authority model, notifications.

## Acceptance criteria

- Every item identifies required authority
- Self-approval restrictions are enforced
- Queue counts reconcile with decision register
- Expired requests are handled explicitly
- Tenant context is immutable during review

## Open questions

- How are emergency requests prioritised?
- Which requests may be auto-cancelled after expiry?
