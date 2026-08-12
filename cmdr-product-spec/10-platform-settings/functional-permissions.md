---
id: platform-settings-functional-permissions
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Platform Settings Functional Permissions — Tenant and Environment Foundations

Security Architecture owns the Permission Model and canonical permission catalogue. This document maps the first Settings capabilities to existing permissions only.

| Capability | Read permission | Manage permission | Additional constraint |
|---|---|---|---|
| CAP-SET-001 | `perm.platform-settings.tenant.read` | `perm.platform-settings.tenant.manage` | Tenant ABAC/isolation; step-up/SoD when Security requires |
| CAP-SET-002 | `perm.platform-settings.environment.read` | `perm.platform-settings.environment.manage` | Tenant + Environment ABAC/isolation |
| CAP-SET-003 | same as target object | same as target object | re-evaluate permission at apply time; audit deny/cross-scope |
| CAP-SET-004 | source read permissions only | none | destination re-evaluates access; context never grants permission |

`perm.settings.tenant.read/manage` remain historical/UI aliases where already referenced. This lot performs no bulk normalization and creates **zero new Permission IDs**. If a capability cannot be expressed safely with the existing permission families, the execution lot is BLOCKED and a separate permission-governance run is required.
