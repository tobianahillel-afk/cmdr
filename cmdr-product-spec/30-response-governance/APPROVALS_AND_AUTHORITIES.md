# Approvals and Authorities

## Objective

Manage who can authorise which actions under which scope and conditions.

## Scope

This specification owns the page-local behaviour of **Approvals and Authorities** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Principal
- Role
- Authority
- Delegation
- Approval chain
- Tenant

## Features

- Authority matrix
- Action and impact thresholds
- Sequential/parallel approval chains
- Delegation and expiry
- Separation-of-duties rules
- Emergency authority
- Conflict and absence handling

## UX and interactions

- Users see why approval is required
- Delegation scope and expiry are prominent
- Authority changes expose effective impact
- Approval identity is strongly authenticated

## Permissions

`authority.manage`; using authority requires `response.approve`; viewing sensitive assignments is ABAC-scoped.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Delegations are scheduled/active/expired/revoked; decisions retain historical authority snapshot.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Identity, RBAC/ABAC, policy, audit, notifications.

## Acceptance criteria

- Requester self-approval rules are enforced
- Expired delegation cannot approve
- Every approval records authentication context
- Authority changes are audited
- Historical decisions remain interpretable

## Open questions

- Is step-up authentication mandatory for critical actions?
- How are on-call rotations integrated?
