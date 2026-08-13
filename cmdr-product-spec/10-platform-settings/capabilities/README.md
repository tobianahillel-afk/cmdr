---
id: platform-settings-capabilities
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Platform Settings Capabilities

Platform Settings capability contracts use the immutable owner-aligned namespace `CAP-SET-*`. IDs are never recycled. Current allocation is exactly `CAP-SET-001..009`; `CAP-SET-010+` is neither allocated nor reserved.

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

This lot contributes **2 capabilities / 54 numbered sections / 12 mandatory tables / 12 meaningful GWT**. Settings cumulative build content is **9 capabilities / 243 sections / 54 mandatory tables**.

`Integration` remains the canonical object; `Connection` is functional/module terminology and Integration `capabilities` metadata is not the CMDR Capability object or a CAP-* namespace. `Secret Reference` remains reference-only. No Connection, Connector, Credential or raw Secret object is created.

`CAP-SET-008` does not claim ownership of the technical external connection-probe engine: current sources support administrative request/preconditions/status/result projection/provenance and handoff only. `CAP-SET-009` does not claim generation, write, rotation or revocation of underlying external secret/credential material: its lifecycle and mutations are Secret Reference administrative semantics only unless a canonical external owner supplies an observed outcome.

All nine capabilities have one canonical capability owner: **Platform Settings Product Lead**. Security, Shared, Experience Architecture, Design System and Govern remain dependency/source owners where applicable.

The Secrets & Connections lot reuses `SET-SEC-001` and `SET-AUD-001`, creates **0 new Screen IDs**, **0 new Permission IDs** and **0 new canonical objects**. Existing Integration/Secret Reference object-facing permissions and current Secrets & Connections surface permissions are reused without bulk namespace normalization.

Build-time quality for Secrets & Connections is **166 PASS / 6 PENDING-REMOTE / 0 FAIL**. Final PASS is forbidden until the remote-dependent gates are actually executed after publication.

Platform Settings Capability Specification remains **PARTIAL** after this lot; Models & Providers, Sources & Parsers and other later Platform Scale work remain separate source-audited lots.
