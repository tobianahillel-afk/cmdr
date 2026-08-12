---
id: platform-settings-action-classification
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# Platform Settings Action Classification — Tenant and Environment Foundations

This mapping consumes the established cross-product action-class semantics; it does not transfer Studio or Govern ownership.

| Class | Settings meaning | Representative first-lot actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Tenant/Environment, state, type, labels, provenance, context | read only under current permission |
| 1 | no-effect deterministic assessment | validate state transition, tenant scope, context compatibility, stale precondition | no mutation or authority |
| 2 | reversible/versioned Settings administration | provision/create, suspend/freeze, reactivate where source permits, offboard/retire, metadata changes | owner-defined mutation; auditable; OPEN-013 remains open |
| 3 | external authority-bearing effect | not created merely by this lot | only when a separate canonical policy requires Govern/runtime authority |
| 4 | destructive/irreversible effect | none defined by this lot | separate source and heightened authority required |

Ordinary Settings administration does not become Govern Response execution simply because it is important. Suspension/offboarding/retirement do not prove physical deletion. If execution reveals a required new Permission ID or Screen ID, this lot becomes BLOCKED; no exception is permitted.
