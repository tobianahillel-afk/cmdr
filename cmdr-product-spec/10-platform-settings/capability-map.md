---
id: platform-settings-capability-map
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Platform Settings Capability Map

## Tenant, Environment and Administrative Foundations

| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-001 | Tenant lifecycle and isolation boundary | Tenant | SET-TEN-001 | Security, Administrative Audit |
| CAP-SET-002 | Environment lifecycle and Tenant scope | Environment, Tenant ref | SET-TEN-001 | Security, Experience Architecture |
| CAP-SET-003 | administrative change validation and provenance | Tenant, Environment | SET-TEN-001, SET-AUD-001 | Security, Administrative Audit |
| CAP-SET-004 | cross-product Tenant/Environment context semantics | Tenant/Environment refs | existing product screens | Experience Architecture, Design System, Security |

## Identity Administration — Principals, Roles and Access Reviews

| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-005 | Principal lifecycle and human/service identity-state administration | Principal, Tenant ref | SET-IAM-001 | Security, Administrative Audit |
| CAP-SET-006 | Role lifecycle/constraints and bounded typed Principal/Role relation | Role, Principal ref, Tenant ref | SET-IAM-001 | Security Permission Model, ABAC, SoD, Decision Authority boundary |
| CAP-SET-007 | periodic Access Review, evidence/provenance and keep/revoke disposition/handoff | Principal/Role refs | SET-IAM-001, SET-AUD-001 | Security, Administrative Audit; Govern only when separately required |

## Secrets & Connections — Integration and Secret Reference Administration

| Capability | Canonical responsibility | Primary objects | Primary screen | Dependencies |
|---|---|---|---|---|
| CAP-SET-008 | Integration administrative lifecycle, local validation, connection-test administrative request/result projection and safe disable | Integration, Tenant ref; Environment only if sourced | SET-SEC-001, SET-AUD-001 | Security, Administrative Audit; external probe executor remains separately owned/unassigned until sourced |
| CAP-SET-009 | Secret Reference lifecycle, reference-only registration, rotation/expiry/revocation of the reference and security handoff | Secret Reference, Tenant ref; Environment only if sourced | SET-SEC-001, SET-AUD-001 | Security secret controls, SoD, step-up, Administrative Audit; underlying secret executor remains external/unassigned until sourced |

Settings owns the capability semantics. Dependency/source ownership remains with the canonical domain. No new Screen ID, Permission ID or canonical object is introduced.

`Principal` ≠ Person/Customer/Tenant; `Role` ≠ Permission/Group/Decision Authority. A typed Principal/Role relation is not itself a Permission grant. CAP-SET-007 does not invent direct assignment-removal mechanics; its `revoke` outcome is a disposition/handoff under the current sources.

`Integration` ≠ Connection object ≠ Connector object ≠ Model Provider ≠ Data Source. Integration `capabilities` metadata is external-connection metadata, not the CMDR Capability object and never creates a `CAP-*` ID. A `Test Connection` action does not transfer technical probe execution to Settings; CAP-SET-008 owns the administrative request/status/result projection and provenance only where the executor is not canonically assigned.

`Secret Reference` ≠ raw Secret ≠ Credential/API Key/Token/Certificate. Secret Reference `rotating`/`revoked` states and reference mutation do not prove generation, write, rotation or revocation of underlying external secret material. CAP-SET-009 remains reference-only and provider-neutral.

Current Settings capability structure after the functional build: **9 capabilities / 243 numbered sections / 54 mandatory tables**. `CAP-SET-010+` is not allocated or reserved. Platform Settings Capability Specification remains **PARTIAL**.
