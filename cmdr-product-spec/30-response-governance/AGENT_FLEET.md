# Agent Fleet

## Objective

Manage the execution and collection agents used by governed response and forensic workflows.

## Scope

This specification owns the page-local behaviour of **Agent Fleet** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Agent
- Tenant
- Capability
- Certificate
- Run
- Health signal

## Features

- Enrollment and revocation
- Connectivity and version
- Capability inventory
- Tenant and environment assignment
- Certificate rotation
- Queued work and recent runs
- Upgrade waves

## UX and interactions

- Agent health does not imply target health
- Dangerous management actions show blast radius
- Offline agents expose queued/expired work
- Capabilities are evidence-backed, not inferred from labels

## Permissions

`agent.manage`; viewing may be limited by tenant and environment; executing work still requires response or analysis permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Agent states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

PKI, updater, execution engine, observability, tenancy.

## Acceptance criteria

- Enrollment is mutually authenticated
- Revoked agents cannot receive work
- Capabilities reconcile with executed checks
- Upgrade rollback is available
- Cross-tenant assignment is blocked

## Open questions

- Which agent deployment models are supported?
- What offline execution policy is allowed?
