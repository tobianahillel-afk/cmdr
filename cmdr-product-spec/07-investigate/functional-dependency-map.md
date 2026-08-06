---
id: investigate-functional-dependency-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
---
# Functional dependency map — Investigate through Phase 4B.3A

The detailed Capability Register shards and governance Dependency Register are canonical. This map summarizes the product-level chain.

| Range | Upstream | Downstream | Owner boundary | Failure behavior |
|---|---|---|---|---|
| CAP-INV-001..114 | Command, Shared, Cases/Evidence | Collection, Analysis, Detection intake | Investigate reasoning; Command runtime retained | missing/stale context explicit |
| CAP-INV-201..215 | Settings/Endpoint/Govern | CAP-INV-301..397 | acquisition owner retained | partial/offline/unsupported explicit |
| CAP-INV-301..397 | Artifact/Collection/Studio/Shared | Evidence/Findings and CAP-INV-401 | technical analysis stays Investigate | limitations and provenance inherited |
| CAP-INV-401..417 | Findings/Hunts/technical handoffs/Settings/Shared/Studio | Review Package | authoring concepts Investigate; runtime absent | no source/field/ground truth invented |
| CAP-INV-418..421 | Review Package, Settings targets, OPEN-017 | Govern Action Request | Investigate prepares; Settings administers targets | target-specific readiness and incompatibility visible |
| CAP-INV-422..424 | candidate, plans and Govern authority | runtime Result projections | Govern/Settings/runtime execute | shadow non-alerting; per-target partiality retained |
| CAP-INV-425..428 | Command/Settings/Endpoint runtime projections | quality, production review and improvements | runtime/Signals remain source-owned | expected/observed divergence explicit |
| CAP-INV-429..434 | quality, drift, performance, coverage and Results | new Draft, Govern request, rollback or retirement | proposal Investigate; active change governed | no silent tuning/suppression/rollback |
| CAP-INV-435 | all lifecycle owners | CAP-INV-401/402/403/406, Settings request, Reporting, future 4B.3B | relation only; no ownership transfer | missing links remain visible |

## Rules
- dependency or projection never transfers ownership;
- source, freshness, tenant/environment/target, permission, immutable version, authority and return origin accompany every transition;
- Shared mechanisms are consumed, never duplicated;
- no model is required for an essential workflow;
- OPEN-017 remains open and selects no runtime or language;
- CAP-INV-5xx and Threat Intelligence remain absent.
