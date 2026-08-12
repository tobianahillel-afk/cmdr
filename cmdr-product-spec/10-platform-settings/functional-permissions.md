---
id: platform-settings-functional-permissions
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Platform Settings Functional Permissions

Security Architecture owns the Permission Model, canonical permission catalogue, RBAC/ABAC, SoD, step-up and authorization evaluation. This document maps Settings capabilities to existing permissions only.

## Tenant and Environment Foundations

| Capability | Read permission | Manage permission | Additional constraint |
|---|---|---|---|
| CAP-SET-001 | `perm.platform-settings.tenant.read` | `perm.platform-settings.tenant.manage` | Tenant ABAC/isolation; step-up/SoD when Security requires |
| CAP-SET-002 | `perm.platform-settings.environment.read` | `perm.platform-settings.environment.manage` | Tenant + Environment ABAC/isolation |
| CAP-SET-003 | same as target object | same as target object | re-evaluate permission at apply time; audit deny/cross-scope |
| CAP-SET-004 | source read permissions only | none | destination re-evaluates access; context never grants permission |

`perm.settings.tenant.read/manage` remain historical/UI aliases where already referenced.

## Identity Administration — Principals, Roles and Access Reviews

| Capability | Surface permission | Object-facing permission | Additional constraint |
|---|---|---|---|
| CAP-SET-005 | `perm.settings.identity.read/manage` | `perm.platform-settings.principal.read/manage` | Tenant first; RBAC/ABAC; SoD/step-up when Security requires; lifecycle mutation never grants Govern authority |
| CAP-SET-006 | `perm.settings.identity.read/manage` | `perm.platform-settings.role.read/manage`; Principal read when relation context is required | Role relation ≠ permission grant; Security computes authorization; Tenant-scope validation mandatory |
| CAP-SET-007 | `perm.settings.identity.read/manage` | Principal/Role read permissions according to review scope | review disposition does not create/revoke Permission; revocation is handoff absent source-proven mechanics |

The `perm.settings.identity.*` screen family and `perm.platform-settings.principal.*` / `perm.platform-settings.role.*` object-facing families are existing identifiers with an unresolved alias/normalization relationship. This lot documents that coexistence and performs **no bulk normalization**.

Across all Settings capability lots, **zero new Permission IDs** are created. Identifiers such as `principal.invite`, `principal.suspend`, `role.assign` and `review.execute` are explicitly not introduced. If a capability cannot be expressed safely with the existing permission families, the relevant lot is BLOCKED and a separate Security-owned permission-design run is required.
