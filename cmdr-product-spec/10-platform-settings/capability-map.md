---
id: platform-settings-capability-map
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Platform Settings Capability Map

## Tenant, Environment and Administrative Foundations

| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-001 | Tenant lifecycle and isolation boundary | Tenant | SET-TEN-001 | Security, Administrative Audit |
| CAP-SET-002 | Environment lifecycle and Tenant scope | Environment, Tenant ref | SET-TEN-001 | Security, Experience Architecture |
| CAP-SET-003 | administrative change validation and provenance | Tenant, Environment | SET-TEN-001, SET-AUD-001 | Security, Administrative Audit |
| CAP-SET-004 | cross-product Tenant/Environment context semantics | Tenant/Environment refs | existing product screens | Experience Architecture, Design System, Security |

Settings owns the capability semantics. Dependency/mechanism ownership remains with the source domain. No new Screen ID, Permission ID or canonical object is introduced.
