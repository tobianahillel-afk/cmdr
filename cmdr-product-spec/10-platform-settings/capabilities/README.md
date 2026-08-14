---
id: platform-settings-capabilities
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
---
# Platform Settings Capabilities

Platform Settings capability contracts use the immutable owner-aligned namespace `CAP-SET-*`. IDs are never recycled. Current allocation is exactly `CAP-SET-001..013`; `CAP-SET-014+` is neither allocated nor reserved.

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

Models & Providers contributes **2 capabilities / 54 numbered sections / 12 mandatory tables / 13 meaningful GWT** and remains **PASS AFTER POST-PUBLICATION VERIFICATION — 204/204 PASS, 0 PENDING, 0 FAIL**.

Canonical boundaries remain: Model Provider is distinct from Integration; model metadata is not a canonical Model object; Policy remains Govern-owned; provider administration/routing configuration is distinct from provider execution; availability, health and effective selection are sourced projections; fallback configuration does not assert automatic failover.

## Sources & Parsers — Data Source Administration and Parser Transformation Administration

| Capability | Title | Status | Delivery |
|---|---|---|---|
| `CAP-SET-012` | Data Source Administrative Lifecycle, Scope, Freshness and Health Projection | draft | defined / planned |
| `CAP-SET-013` | Parser Administrative Lifecycle, Versioned Transformation, Fixtures and Validation | draft | defined / planned |

Sources & Parsers contributes **2 capabilities / 54 numbered sections / 12 mandatory tables / 8 meaningful GWT**. Settings cumulative functional BUILD content becomes **13 capabilities / 351 sections / 78 mandatory tables**.

Canonical boundaries:
- `Data Source` and `Parser` remain distinct canonical objects with no direct canonical relationship;
- no Source→Parser assignment, compatibility object, routing, selection, fallback or precedence is introduced;
- `Integration` remains distinct from Data Source;
- source health/freshness and Parser error/quality/test results are source-attributed projections when runtime-derived;
- Settings does not own acquisition, collection, ingestion, connector/probe execution, parser engine, normalization, stream processing, schema-registry implementation or storage merely by defining these capabilities;
- schema format, initial schema version and SLO/limits remain unresolved in implementation contracts;
- no ECS, OCSF, CIM, OpenTelemetry or other unsourced schema standard is selected.

All thirteen capabilities have one capability owner: **Platform Settings Product Lead**. Security, Shared, Experience Architecture, Design System, Govern, Health, Administrative Audit, Studio/runtime and other product/runtime owners retain their canonical ownership.

The lot reuses `SET-SRC-001`, `SET-HLT-001`, `SET-AUD-001`, `SET-SEC-001`, existing `perm.platform-settings.data-source.*`, `perm.platform-settings.parser.*` and existing `perm.settings.source.*` aliases. It creates **0 new Screen IDs, 0 new Permission IDs and 0 new canonical objects**.

BUILD-time quality for Sources & Parsers is **292 PASS / 24 PENDING-REMOTE / 0 FAIL**. Final `316/316` is forbidden until actual remote publication, CI/applicability inspection and documentary closure complete.

Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**. Documentary capability definition does not claim implementation/runtime availability.

Platform Settings Capability Specification remains **PARTIAL** after this lot; later Platform Scale work remains separate source-audited work. `CAP-SET-014+` remains unallocated and unreserved.
