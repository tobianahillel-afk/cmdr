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

## Identity Administration — Principals, Roles and Access Reviews

| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-005 | Principal lifecycle and human/service identity-state administration | Principal, Tenant ref | SET-IAM-001 | Security, Administrative Audit |
| CAP-SET-006 | Role lifecycle/constraints and bounded typed Principal/Role relation | Role, Principal ref, Tenant ref | SET-IAM-001 | Security Permission Model, ABAC, SoD, Decision Authority boundary |
| CAP-SET-007 | periodic Access Review, evidence/provenance and keep/revoke disposition/handoff | Principal/Role refs | SET-IAM-001, SET-AUD-001 | Security, Administrative Audit; Govern only when separately required |

Settings owns the capability semantics. Dependency/source ownership remains with the canonical domain. No new Screen ID, Permission ID or canonical object is introduced.

`Principal` ≠ Person/Customer/Tenant; `Role` ≠ Permission/Group/Decision Authority. A typed Principal/Role relation is not itself a Permission grant. CAP-SET-007 does not invent direct assignment-removal mechanics; its `revoke` outcome is a disposition/handoff under the current sources.
