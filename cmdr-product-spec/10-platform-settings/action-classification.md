---
id: platform-settings-action-classification
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# Platform Settings Action Classification

This mapping consumes the established cross-product action-class semantics; it does not transfer Studio, Security or Govern ownership.

## Tenant and Environment Foundations

| Class | Settings meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Tenant/Environment, state, type, labels, provenance, context | read only under current permission |
| 1 | no-effect deterministic assessment | validate state transition, tenant scope, context compatibility, stale precondition | no mutation or authority |
| 2 | reversible/versioned Settings administration | provision/create, suspend/freeze, reactivate where source permits, offboard/retire, metadata changes | owner-defined mutation; auditable; OPEN-013 remains open |
| 3 | external authority-bearing effect | not created merely by this lot | only when a separate canonical policy requires Govern/runtime authority |
| 4 | destructive/irreversible effect | none defined by this lot | separate source and heightened authority required |

## Identity Administration — Principals, Roles and Access Reviews

| Class | Settings meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Principal/Role/state/typed relation/Access Review/provenance | read under current permission; no authority transfer |
| 1 | no-effect deterministic assessment | validate Principal or Role transition, Tenant scope, Role expiry condition, explicit Security constraint, preview review disposition | no mutation; Security remains authorization owner |
| 2 | reversible/versioned administration where source-backed | initialize Principal for `Inviter`, suspend/sourced lifecycle mutation, create/mutate Role, start/record review disposition, typed relation mutation only if independently source-proven | auditable; OPEN-013 remains open; no autonomous AI mutation |
| 3 | authority-bearing effect | not created by Identity Administration | preserve Security/Govern/runtime owner when separately required |
| 4 | destructive/irreversible effect | none defined by this lot | separate source and heightened authority required |

Role expiry is a condition/constraint, never a Role state. Access Review revocation is represented as a **revocation disposition/handoff** while no canonical generic assignment-removal mechanic exists. A Principal/Role relation is not a Permission definition or effective authorization.

Ordinary Settings administration does not become Govern Response execution simply because it is important. No Group, AccessAssignment, PermissionAssignment, Effective Access engine, new Permission ID or new Screen ID is introduced. If a genuinely required operation cannot be expressed within these existing boundaries, the lot is BLOCKED rather than reclassified locally.
