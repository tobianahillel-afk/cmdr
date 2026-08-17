---
id: platform-settings-action-classification
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
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

## Models & Providers — Provider Administration and Routing Configuration

| Class | Settings meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Model Provider, configuration, sourced model-availability/health projection, routing configuration and provenance | read-only; observed runtime facts remain projections from their source |
| 1 | strictly local no-effect assessment | validate Model Provider configuration, lifecycle preconditions, routing/fallback constraints and stale version | deterministic local validation only; no external provider operation |
| 2 | reversible/versioned Settings administration | create/update/activate/disable Model Provider configuration; update allowed-use/routing/fallback configuration where source-backed | tenant-scoped, auditable; `OPEN-013` remains open |
| 3 | external authority-bearing/runtime effect | not assigned to Settings merely by `Tester`, routing configuration or a provider switch observation | actual external execution retains its separately sourced runtime/authority owner |
| 4 | destructive/irreversible effect | none defined by this lot | separate source and heightened authority required |

A provider-test request, model availability observation, health observation or effective provider/model selection is not reclassified as a local Settings execution. Fallback configuration is not automatic failover. No silent provider switch is permitted in the Settings projection.

## Sources & Parsers — Data Source and Parser Administration

| Class | Settings meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Data Source/Parser state, configuration, version, freshness, sourced health, error/quality metadata and provenance | read only; runtime observations remain source-attributed projections |
| 1 | strictly local no-effect assessment | validate Tenant scope, lifecycle preconditions, stale version, schema reference, fixture definition or other locally available configuration facts | deterministic local validation only; no connector, collection, probe or parser-engine execution |
| 2 | reversible/versioned Settings administration | create/update/activate/disable Data Source configuration; create/update/version/activate/retire Parser configuration where source-backed | tenant-scoped, auditable; `OPEN-013` remains open; no Source→Parser assignment is implied |
| 3 | external/runtime or authority-bearing effect | `Test Source` or `Test Parser` only when it invokes an external service, connector, collection path, parser engine or other runtime | Settings owns request/preconditions/handoff/observed-result projection only; executor remains with its separately sourced owner |
| 4 | destructive/irreversible effect | none defined by this lot | separate source and heightened authority required |

`Test Source` and `Test Parser` are therefore not automatically Class 1. An administrative test button cannot silently reclassify an outbound probe, collection, ingestion or parser execution as a no-effect Settings operation. The current corpus also defines no automatic parser selection, Source→Parser routing, fallback or precedence.

Ordinary Settings administration does not become Govern Response execution merely because it is important. No new Permission ID, Screen ID, canonical object, Requirement ID or OPEN ID is introduced. If a genuinely required operation cannot be expressed within these existing boundaries, the lot is BLOCKED rather than reclassified locally.
