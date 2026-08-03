# Tenancy and Scoping

## Objective

Enforce strict customer and organisational isolation while supporting authorised MSSP operations.

## Scope

This specification owns the page-local behaviour of **Tenancy and Scoping** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Tenant
- Workspace
- Principal
- Role assignment
- Data object

## Features

- Tenant lifecycle
- User and service assignment
- Environment and business-unit scopes
- MSSP portfolio views using safe summaries
- Data residency and jurisdiction attributes
- Tenant export and deletion workflows

## UX and interactions

- Current tenant is always visible
- Cross-tenant views use explicit aggregate mode
- Switching tenant clears incompatible selection and cache
- No hidden fallback tenant

## Permissions

`tenant.manage` for administration; every other permission is evaluated inside tenant scope.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Tenant states: provisioning, active, suspended, offboarding, archived.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Identity, storage, search, integrations, retention.

## Acceptance criteria

- Object queries require tenant key
- Caches and exports are tenant-scoped
- Suspension blocks mutations according to policy
- Offboarding is auditable
- Aggregate views cannot expose unauthorised raw data

## Open questions

- What cross-tenant analytics are permitted?
- How are shared global threat-intel objects represented?
