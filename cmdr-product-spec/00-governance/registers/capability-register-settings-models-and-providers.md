---
id: capability-register-settings-models-and-providers
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-13
source-of-truth: registry
---
# Capability Register — Settings Models & Providers

Execution lot: **Models & Providers — Model Provider Administration and Model Routing** under **Delivery Roadmap Phase 6 — Platform Scale**.

| Capability ID | Name | Owner/module | Status | Delivery | Canonical file | Primary objects | Requirement IDs | OPEN |
|---|---|---|---|---|---|---|---|---|
| `CAP-SET-010` | Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection | Platform Settings Product Lead / Models & Providers | draft | defined / planned | `10-platform-settings/capabilities/cap-set-010-model-provider-administrative-lifecycle-validation-availability-and-health-projection.md` | Model Provider, Tenant ref; optional sourced refs only | REQ-PROD-004/005/006/009/010/011/012; REQ-AI-001/003/004/009/010/011 | OPEN-008, OPEN-012, OPEN-013 |
| `CAP-SET-011` | Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance | Platform Settings Product Lead / Models & Providers | draft | defined / planned | `10-platform-settings/capabilities/cap-set-011-model-routing-configuration-eligibility-fallback-and-provider-switch-provenance.md` | Model Provider configuration, Tenant ref | REQ-PROD-004/005/006/009/010/011/012; REQ-AI-001/003/004/009/010/011 | OPEN-008, OPEN-012, OPEN-013 |

## Structural totals
- capabilities: **2**;
- defined / proposed / planned: **2 / 0 / 2**;
- numbered sections: **54**;
- mandatory tables: **12**;
- meaningful Given/When/Then scenarios: **13**;
- owner conflicts: **0**;
- duplicate/recycled IDs: **0**.

Settings cumulative functional BUILD content becomes **11 capabilities / 297 sections / 66 mandatory tables**.

## Namespace
`CAP-SET-010` and `CAP-SET-011` are allocated by this lot. **No `CAP-SET-012+` ID or range is allocated or reserved.**

## Boundaries
- `Model Provider` remains the canonical Platform Settings object for this lot.
- `Integration` remains distinct.
- Model metadata does not create a canonical `Model` object.
- `Policy` remains Govern-owned; no canonical routing-policy object is created.
- `Secret Reference` remains a source-dependent reference when applicable; no mandatory relation is invented.
- Settings owns administrative configuration and provenance, not provider execution.
- model availability, health and effective selection remain sourced projections when supplied externally.
- fallback configuration does not assert automatic failover.

## Existing IDs reused
This lot creates **0 new Screen IDs, 0 new Permission IDs and 0 new canonical objects**. It reuses `SET-MDL-001`, `SET-HLT-001`, `SET-AUD-001`, `perm.platform-settings.model-provider.read/manage` and existing `perm.settings.model.read/manage` aliases.

## Requirements and OPEN
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**. `OPEN-008`, `OPEN-012` and `OPEN-013` remain open. No Requirement or OPEN ID is created, deleted or resolved by this lot.

## Delivery boundary
Both capabilities remain documentary `defined / planned`. No runtime availability or implementation is claimed. Platform Settings Capability Specification and Delivery Roadmap Phase 6 Capability Specification remain **PARTIAL**.