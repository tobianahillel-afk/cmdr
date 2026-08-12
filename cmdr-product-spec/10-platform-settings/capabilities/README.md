---
id: platform-settings-capabilities
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Platform Settings Capabilities

Platform Settings capability contracts use the immutable owner-aligned namespace `CAP-SET-*`. IDs are never recycled. Current allocation is exactly `CAP-SET-001..007`; `CAP-SET-008+` is neither allocated nor reserved by the Identity Administration lot.

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

All seven capabilities have one canonical capability owner: **Platform Settings Product Lead**. Security, Shared, Experience Architecture, Design System and Govern remain dependency/source owners where applicable. Identity Administration reuses `SET-IAM-001` and `SET-AUD-001`, creates **0 new Screen IDs**, **0 new Permission IDs** and **0 new canonical objects**.

Groups, generic Access Assignment and Effective Access are outside the Identity lot. Role expiry is a condition/constraint, not a state. CAP-SET-007 uses a revocation disposition/handoff because the current canonical corpus does not define generic assignment-removal mechanics.

Platform Settings Capability Specification remains **PARTIAL** after this lot; later Settings/Shared Platform Scale lots remain separately source-audited work.
