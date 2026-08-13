---
id: capability-register-settings-secrets-and-connections
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-13
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-011
  - REQ-PROD-012
---
# Capability Register — Settings Secrets & Connections

## Lot identity
Execution lot: **Secrets & Connections — Integration and Secret Reference Administration** under **Delivery Roadmap Phase 6 — Platform Scale**. This is an execution lot, not a new roadmap phase.

## Active capability rows
| Capability ID | Name | Owner/module | Status | Delivery | Canonical file | Primary roles | Primary objects | Consumers | Requirement IDs | OPEN | Dependencies | Supersession | Review |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| `CAP-SET-008` | Integration Administrative Lifecycle, Validation and Connection State | Platform Settings Product Lead / Secrets & Connections | draft | defined / planned | `10-platform-settings/capabilities/cap-set-008-integration-administrative-lifecycle-validation-and-connection-state.md` | Platform Administrator; Security Administrator secondary | Integration, Tenant ref; Environment only if sourced | Settings, Studio, Endpoint, Command, Investigate, Govern, Shared | REQ-PROD-006/009/010/011/012; REQ-AI-001/003/004/009/010 | OPEN-008, OPEN-013 | Integration object; SET-SEC-001; Security; Administrative Audit | none | 2026-08-13 |
| `CAP-SET-009` | Secret Reference Administrative Lifecycle, Rotation and Revocation | Platform Settings Product Lead / Secrets & Connections | draft | defined / planned | `10-platform-settings/capabilities/cap-set-009-secret-reference-administrative-lifecycle-rotation-and-revocation.md` | Platform Administrator; Security Administrator secondary | Secret Reference, Tenant ref; Environment only if sourced | Settings, Integrations, Studio, Endpoint, Command, Investigate, Govern, Shared | REQ-PROD-006/009/010/011/012; REQ-AI-001/003/004/009/010 | OPEN-013 | Secret Reference object; SET-SEC-001; Security secret controls; Administrative Audit | none | 2026-08-13 |

## Structural totals
- capabilities: **2**;
- defined / proposed / planned: **2 / 0 / 2**;
- numbered sections: **54**;
- mandatory tables: **12**;
- meaningful Given/When/Then scenarios: **12**;
- owner conflicts: **0**;
- duplicate/recycled IDs: **0**.

Settings cumulative current build content becomes **9 capabilities / 243 sections / 54 mandatory tables**.

## Namespace
`CAP-SET-008` and `CAP-SET-009` are allocated immutably by this lot after execution-time revalidation. No `CAP-SET-010+` ID is allocated or reserved.

## Canonical distinctions
- `Integration` is canonical. `Connection` is module/functional terminology, not an object.
- Integration `capabilities` metadata is not the CMDR Capability object and never allocates `CAP-*` IDs.
- No canonical `Connector` or `Credential` object is created.
- `Secret Reference` is canonical and reference-only; raw secret material is never a Settings object exposed by these contracts.
- `Secret Reference` rotation/revocation is not automatically rotation/revocation of external secret or credential material.
- Model Provider, Data Source and Parser remain separate Settings objects/later lots.

## Executor boundaries
Current canonical sources do **not** assign Platform Settings the technical network/probe executor for `Test Connection`. `CAP-SET-008` therefore owns administrative request/preconditions/status/result projection/provenance and handoff, while external execution remains unassigned until a canonical owner exists.

Current canonical sources do **not** assign Platform Settings the generator/writer/provider executor for underlying secret rotation or external credential revocation. `CAP-SET-009` therefore owns Secret Reference lifecycle/reference mutation and administrative handoff/result projection only.

## Screen, permission and object gates
This lot creates **0 new Screen IDs, 0 new Permission IDs and 0 new canonical objects**. It reuses `SET-SEC-001`, `SET-AUD-001`, `perm.platform-settings.integration.read/manage`, `perm.platform-settings.secret-reference.read/manage` and the existing `perm.settings.secret.read-metadata/manage` surface permissions without bulk normalization.

## Requirements and OPEN
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**; this lot creates/deletes no Requirement ID and makes no state transition merely for PASS. OPEN remains **18**. `OPEN-008`, `OPEN-012` and `OPEN-013` remain unresolved as applicable; no new OPEN is manufactured.

## Delivery boundary
Both capabilities are documentary `defined / planned`. No implementation, API, protocol, physical schema, provider support, connector runtime, Vault/KMS/HSM, final RBAC/ABAC or product release is claimed. Platform Settings Capability Specification and Delivery Roadmap Phase 6 Capability Specification remain **PARTIAL**.
