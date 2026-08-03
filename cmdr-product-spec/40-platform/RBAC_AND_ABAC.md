# RBAC and ABAC

## Objective

Implement the canonical permission model consistently across UI, API, search, export and execution.

## Scope

This specification owns the page-local behaviour of **RBAC and ABAC** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Principal
- Role
- Permission
- Attribute
- Authority
- Policy decision

## Features

- Role catalogue
- Custom role constraints
- Attribute policies
- Permission simulation
- Effective-access explanation
- Separation of duties
- Access review and expiry

## UX and interactions

- Users can understand missing access without seeing sensitive rules
- Administrators preview impact before changes
- UI visibility follows but does not replace backend enforcement

## Permissions

`authority.manage`, `tenant.manage` and `platform.admin` are themselves policy-controlled.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Role assignments are pending/active/expired/revoked; policy decisions are allow/deny/indeterminate.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Identity, tenancy, policy engine, audit.

## Acceptance criteria

- Backend enforces every operation
- Search and export use same policy context
- Deny overrides are consistent
- Effective access is explainable
- Privilege changes are audited

## Open questions

- Which policy engine is selected?
- How are emergency privileges time-boxed?
