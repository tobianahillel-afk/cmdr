---
id: platform-settings-capabilities
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Platform Settings Capabilities

Platform Settings capability contracts use the immutable owner-aligned namespace `CAP-SET-*`. IDs are never recycled. Current allocation is exactly `CAP-SET-001..011`; `CAP-SET-012+` is neither allocated nor reserved.

## Tenant, Environment and Administrative Foundations

| Capability | Title | Status | Delivery |
|---|---|---|---|
| `CAP-SET-001` | Tenant Administrative Lifecycle and Isolation Boundary | draft | defined / planned |
| `CAP-SET-002` | Environment Administrative Lifecycle and Tenant Scope | draft | defined / planned |
| `CAP-SET-003` | Tenant and Environment Administrative Change Validation and Provenance | draft | defined / planned |
| `CAP-SET-004` | Cross-Product Tenant and Environment Context Preservation | draft | defined / planned |

## Identity Administration — Principals, Roles and Access Reviews

| Capability | Title | Status | Delivery |
|---|---|---|---|
| `CAP-SET-005` | Principal Administrative Lifecycle and Identity State | draft | defined / planned |
| `CAP-SET-006` | Role Administrative Lifecycle, Constraints and Principal-Relation Boundary | draft | defined / planned |
| `CAP-SET-007` | Access Review, Evidence and Revocation Disposition | draft | defined / planned |

Identity Administration remains **PASS AFTER POST-PUBLICATION VERIFICATION — 174/174 PASS, 0 PENDING, 0 FAIL**. Groups, generic Access Assignment and Effective Access remain outside that lot. Role expiry is a condition/constraint, not a state. CAP-SET-007 uses a revocation disposition/handoff because the current canonical corpus does not define generic assignment-removal mechanics.

## Secrets & Connections — Integration and Secret Reference Administration

| Capability | Title | Status | Delivery |
|---|---|---|---|
| `CAP-SET-008` | Integration Administrative Lifecycle, Validation and Connection State | draft | defined / planned |
| `CAP-SET-009` | Secret Reference Administrative Lifecycle, Rotation and Revocation | draft | defined / planned |

Secrets & Connections contributes **2 capabilities / 54 numbered sections / 12 mandatory tables / 12 meaningful GWT** and remains **PASS AFTER POST-PUBLICATION VERIFICATION — 172/172 PASS, 0 PENDING, 0 FAIL**.

`Integration` remains canonical; `Connection` is functional terminology. `Secret Reference` remains reference-only. CAP-SET-008 does not claim the technical external connection-probe executor; CAP-SET-009 does not claim the underlying external value-management executor.

## Models & Providers — Model Provider Administration and Model Routing

| Capability | Title | Status | Delivery |
|---|---|---|---|
| `CAP-SET-010` | Model Provider Administrative Lifecycle, Validation, Model Availability and Health Projection | draft | defined / planned |
| `CAP-SET-011` | Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance | draft | defined / planned |

Models & Providers contributes **2 capabilities / 54 numbered sections / 12 mandatory tables / 13 meaningful GWT**. Settings cumulative functional BUILD content is **11 capabilities / 297 sections / 66 mandatory tables**.

Canonical boundaries:
- Model Provider remains distinct from Integration;
- model metadata and availability projections do not create a canonical Model object;
- provider administration and routing configuration do not transfer provider/runtime execution to Settings;
- fallback configuration does not assert automatic failover;
- effective provider/model selection is a source-attributed observation when supplied by a runtime owner;
- Policy remains Govern-owned;
- optional source-backed Secret Reference use does not create an invented mandatory relation.

All eleven capabilities have one capability owner: **Platform Settings Product Lead**. Security, Shared, Experience Architecture, Design System, Govern, Health, Administrative Audit and Studio/runtime consumers retain their canonical dependency/source ownership.

The Models & Providers lot reuses `SET-MDL-001`, `SET-HLT-001`, `SET-AUD-001`, existing Model Provider permissions and current UI aliases. It creates **0 new Screen IDs, 0 new Permission IDs and 0 new canonical objects**.

Build-time quality for Models & Providers is **198 PASS / 6 PENDING-REMOTE / 0 FAIL**. Final `204/204` is forbidden until remote publication verification is complete.

Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**, including `OPEN-008`, `OPEN-012` and `OPEN-013` unresolved. Documentary capability definition does not claim implementation/runtime availability.

Platform Settings Capability Specification remains **PARTIAL** after this lot; Sources & Parsers and other later Platform Scale work remain separate source-audited lots. `CAP-SET-012+` remains unallocated and unreserved.