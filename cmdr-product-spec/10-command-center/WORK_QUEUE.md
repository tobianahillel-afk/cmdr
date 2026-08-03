# Work Queue

## Objective

Coordinate analyst work across incidents, alerts, investigations and decisions without mixing incompatible object lifecycles.

## Scope

This specification owns the page-local behaviour of **Work Queue** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Incident
- Alert
- Case
- Response request
- Task
- Principal

## Features

- Saved views and personal/team queues
- Assignment, ownership and hand-off
- SLA and ageing indicators
- Bulk acknowledgement or assignment where policy allows
- Dependency and blocker badges
- Shift handover summary

## UX and interactions

- Column configuration persists per user
- Bulk actions show exact scope and partial-failure results
- Opening detail does not reset filters
- Queue supports keyboard triage
- Handover captures unresolved decisions and context

## Permissions

`command.read`; assignment and bulk operations require `command.coordinate` or object-specific manage permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Displays canonical states for each object type and never merges them into one generic status.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Task service, identity, SLA policy, notifications.

## Acceptance criteria

- Saved filters are reproducible by URL
- Bulk actions are idempotent
- Users cannot assign across inaccessible tenants
- Queue counts reconcile with detail views
- Shift handover is auditable

## Open questions

- Should queue ranking be configurable or fixed?
- How are AI suggestions displayed without becoming authority?
