---
id: platform-settings-functional-permissions
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
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

## Sources & Parsers — existing permission mapping

The Sources & Parsers execution lot uses existing permission families only. At this boundary-preparation step the capability identifiers are not yet claimed as canonically allocated; their object-facing permission mappings are fixed by the canonical objects and existing `SET-SRC-001` surface.

| Administrative scope | Surface permission | Object-facing permission | Additional constraint |
|---|---|---|---|
| Data Source administration | `perm.settings.source.read/manage` | `perm.platform-settings.data-source.read/manage` | Tenant first; RBAC/ABAC; local validation does not grant connector/collection/ingestion/probe execution |
| Parser administration | `perm.settings.source.read/manage` | `perm.platform-settings.parser.read/manage` | Tenant first; RBAC/ABAC; manage does not grant parser-engine execution, Source→Parser assignment or schema-standard selection |

The `perm.settings.source.*` screen family and `perm.platform-settings.data-source.*` / `perm.platform-settings.parser.*` object-facing families already coexist. This run preserves that coexistence and performs **no bulk normalization**. `read` never implies export, execute or approve. An existing `.manage` permission cannot be expanded into an unsourced runtime authority.

Across all Settings capability lots, **zero new Permission IDs** are created. Identifiers such as `principal.invite`, `principal.suspend`, `role.assign`, `review.execute`, `source.test.execute`, `parser.execute` or `parser.assign` are explicitly not introduced. If a capability cannot be expressed safely with the existing permission families, the relevant lot is BLOCKED and a separate Security-owned permission-design run is required.
