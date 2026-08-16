---
id: CAP-SET-014
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-16
source-of-truth: canonical
delivery_status: defined
delivery_mode: planned
primary_roles:
  - Platform Administrator
primary_objects:
  - Tenant
  - Environment
  - Integration
  - Data Source
  - Model Provider
  - Endpoint Agent Fleet
consumers:
  - platform-settings
  - command
  - shared
  - security
requirements:
  - REQ-PROD-003
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-PROD-021
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-019
dependencies:
  - CAP-SET-001
  - CAP-SET-002
  - CAP-SET-010
  - CAP-SET-012
  - CAP-CMD-401
---
# CAP-SET-014 — Platform Health and Source-Attributed SLO Projection

## S1. Identity and owner
`CAP-SET-014` defines the Platform Settings-owned read/project/handoff capability for Platform Health and source-attributed SLO visibility. The capability owner is **Platform Settings Product Lead**. Platform Architecture retains neutral Health/Metrics contract semantics; Shared retains generic Metrics, Search, Reporting, Export and Notification mechanisms; authoritative source/runtime owners retain measurement acquisition and authoritative source-specific calculation.

## S2. User goal
Allow an authorized user to understand Platform Health and source-attributed SLO state through a deterministic read-only projection that preserves source, Tenant, freshness, provenance, uncertainty, impact and controlled handoffs.

## S3. Scope
The capability covers Tenant-local Health/SLO projection and MSSP read-only aggregation over the Security-resolved Authorized Tenant Set. It reuses `SET-HLT-001` and existing context/consumer surfaces. It does not configure SLOs, administer Health, run monitoring, execute failover/recovery, create Incidents/Tasks automatically or publish SLOs externally.

## S4. Canonical architecture boundary
`ADR-0009` is normative. SLO is source-attributed and non-canonical. There is no generic CMDR SLO target store, generic SLO configuration, central SLO calculator, generic failover/recovery runtime or new SLO/Health canonical object. `ADR-0008` remains normative for Tenant/MSSP isolation and single-Tenant Search/Report/Export.

## S5. Roles and consumers
Primary use is by an authorized Platform Administrator or other consumer that already satisfies `perm.settings.health.read` and Tenant/source visibility. Command consumes explicit selected-Tenant handoffs; Shared supplies generic metric/notification/report/export mechanisms; Security remains authoritative for authorization and Authorized Tenant Set resolution.

## S6. Preconditions
A Tenant is selected or attributable. The user is authorized by Security for the Tenant and source. A typed subject reference, authoritative source reference and freshness information are available. Environment is carried only when sourced. Missing authoritative evidence must remain explicit rather than synthesized.

## S7. Functional model
The capability projects source-backed Health, degradation, SLO target metadata, source-calculated SLO state/breach metadata, freshness, provenance and impact context. It may derive bounded presentation facts such as age/stale/unknown/partial from explicit source metadata. It never manufactures a measurement, compliance state or breach.

## S8. Inputs
| Input | Requirement | Semantics |
|---|---|---|
| Tenant | required | Isolation and attribution boundary for every projection. |
| Typed subject reference | required | Identifies the sourced subject without creating a new canonical object. |
| Authoritative source reference | required | Identifies the source that owns the projected state/target. |
| Freshness | required | Preserves observation/metadata recency and stale/unknown interpretation. |
| Environment | conditional | Included only when the source defines it. |
| Source Health state | conditional | Source-backed Health/degradation state. |
| SLO target value/unit/version/effective period | conditional | Required as applicable when a target is displayed. |
| Measurement window | conditional | Required when authoritative SLO state/breach is displayed. |
| Source-calculated SLO state/breach | conditional | Never synthesized by Settings. |
| Calculation provenance | conditional | Required for authoritative state/breach. |
| Impact / Business Service context | conditional | Read-only context; ownership remains with its canonical source. |
| Authorized Tenant Set | conditional MSSP | Security-resolved non-canonical authorization projection. |

## S9. Objects Read
| Object / reference | Read purpose | Ownership boundary |
|---|---|---|
| Tenant | Scope and attribution | Platform Settings administrative context; Security enforces isolation. |
| Environment | Optional sourced context | Platform Settings. |
| Integration | Health subject/context | Existing canonical object only. |
| Data Source | Health/freshness subject | Existing canonical object only. |
| Model Provider | Provider Health context | Existing canonical object only. |
| Endpoint Agent Fleet | Fleet Health context | Existing canonical object only. |
| Policy | Read only where explicitly sourced | Govern-owned; no Policy mutation. |
| Incident / Task | Existing downstream Command context only | Command-owned; never created by this capability. |
| Business Service identity/context | Impact/service context | Shared catalog context; no canonical Service object is created. |
| SLO/source metadata | Non-canonical projection | Authoritative source remains owner. |

## S10. Objects Created/Modified
| Created/modified canonical object | Result | Reason |
|---|---|---|
| None | **No canonical object is created or modified by this capability.** | The capability is intentionally read/project/handoff oriented under ADR-0009; manufacturing a write would violate the architecture. |

## S11. Action classes
C0 actions are: view Health, target, sourced SLO state/breach, freshness, source/provenance and impact; navigate Tenant/context; read an Authorized-Tenant-Set aggregate; open canonical source/context or Business Service context. Handoffs are human Command handoff, Shared Report request, Shared Export request and optional Shared Notification handoff. `Acknowledge maintenance` remains disabled. Generic SLO/threshold/window/breach-policy configuration, generic monitoring enable/disable, failover, recovery, DR and RTO/RPO administration are unsupported/current-architecture-out.

## S12. SLO target, state and breach semantics
A target is displayed only with authoritative attribution and, as applicable, target value/unit, version and effective-period/source context. An SLO state or breach is authoritative only when Tenant, subject, authoritative source, target/version, measurement window, calculation provenance, freshness and source-calculated state/breach evidence are available. Insufficient evidence remains `unknown`, `partial`, `stale`, `conflicting` or `unsupported`. Missing is never inferred as healthy; stale is never inferred as compliant; a measurement is never inferred as compliance/breach; a breach never automatically becomes Incident, Task, contractual SLA breach, priority change, response or Govern invocation.

## S13. Automation/AI
| AI/automation behavior | Allowed | Constraint |
|---|---|---|
| Explain sourced Health/SLO state, target/source/version, degradation and freshness | yes | Must cite/preserve source and uncertainty. |
| Summarize available measurements or sourced historical trend | yes | No fabricated measurement, availability, reliability or compliance. |
| Suggest investigation or explicit human handoff | yes | Suggestion only; deterministic/manual path remains available. |
| Modify target, hide breach, grant permission/authority | no | No mutation or authority creation. |
| Trigger failover/recovery, Response Run or Automation Run | no | No autonomous effectful action. |
| Mix Tenants or widen Search/Report/Export | no | Tenant and ADR-0008 boundaries remain authoritative. |

## S14. Tenant and MSSP semantics
Tenant is retained on every projection. MSSP aggregation is read-only and limited to the Security-resolved Authorized Tenant Set; every row/result preserves its own Tenant, source, target version, freshness and provenance. Tenant switching/navigation does not grant access. Cross-tenant mutation, administration, response and export widening are forbidden. Search remains single-selected-Tenant; Report and Export remain single-Tenant.

## S15. Security and permissions
`perm.settings.health.read` is the sole capability-specific permission used here and remains read-only. It does not imply export, configuration, administration, monitoring execution, failover, recovery, approval or response authority. Report/Export use Shared authorization semantics. No new Permission ID is created. `OPEN-013` remains open because this capability defines no generic write.

## S16. Outputs
| Output | Semantics | Consumer |
|---|---|---|
| Tenant-local Health projection | Source-backed status/degradation/freshness with uncertainty | Platform Health / authorized user |
| Source-attributed SLO target projection | Target value/unit/version/effective-period context where available | Platform Health |
| Source-attributed SLO state/breach projection | Only authoritative with required window/calculation provenance/freshness | Platform Health |
| Impact / Business Service context | Read-only sourced context | Platform Health / Command handoff |
| MSSP authorized aggregate | Read-only multi-Tenant projection preserving Tenant per result | Authorized MSSP user |
| Deep link / source navigation | Stable return/source context | Existing canonical surfaces |
| Handoff package/context | Selected-Tenant context only; no new canonical object | Command or Shared |

## S17. Transitions/Handoffs
| From | To | Handoff | Constraint |
|---|---|---|---|
| Authoritative source/runtime | Platform Health | Sourced Health/SLO state, target, freshness, provenance | Source/runtime retains acquisition/calculation ownership. |
| Platform Health | Shared Notification | Optional notification request/context | Notification mechanism remains Shared-owned. |
| Platform Health | Command | Explicit selected-Tenant human handoff | No automatic Incident, Task, priority change or response. |
| Platform Health | Shared Reporting | Single-Tenant Report request | `health.read` alone does not authorize Report. |
| Platform Health | Shared Export | Single-Tenant Export request | Export never widens visibility. |
| MSSP aggregate | Tenant-local context | Explicit Tenant navigation | Security authorization is re-evaluated; no cross-Tenant effect. |

## S18. Screen and UX model
Primary screen is existing `SET-HLT-001`; new Screen IDs are zero. Context/consumer screens are `SET-TEN-001`, `SET-AUD-001`, `CMD-MC-001`, `CMD-IWQ-006`, `CMD-RBI-001` and `CMD-CRP-001`. `SET-HLT-001`/CAP-SET-014 owns Health/SLO projection; `CMD-CRP-001`/CAP-CMD-401 remains Customers & Delivery overview/context. Existing loading, empty, partial, error, offline and permission-denied states remain mandatory.

## S19. Explicit state handling
`unknown`, `partial`, `stale`, `conflicting` and `unsupported` remain visible source/quality conditions. Loading/Empty/Error/Offline/Permission denied remain UX states. Empty does not mean healthy. Offline/Retry does not imply failover. Valid source-backed information may remain visible alongside partial/error indicators where the existing UX contract permits it.

## S20. Freshness and provenance
Every projected fact preserves its authoritative source and freshness. SLO target/state/breach additionally preserves target version/effective period and calculation/window provenance as applicable. Settings may derive presentation age/stale/unknown/partial from explicit source metadata but cannot synthesize authoritative Health/SLO computation.

## S21. Resilience boundary
The initial resilience scope is observation of sourced degradation, existing Offline/Retry semantics and existing provider fallback constraints as configuration. Retry, fallback, failover, recovery and rollback remain distinct. Generic runtime redundancy, automatic/manual failover, generic Data Source fallback, recovery execution, DR, RTO and RPO remain outside this capability.

## S22. Search, Reporting and Export
Search remains single-selected-Tenant. Report and Export remain Shared-owned and single-Tenant. CAP-SET-014 may hand off selected-Tenant context but does not implement or authorize Reporting/Export and cannot turn an MSSP aggregate into multi-Tenant Search, Report or Export.

## S23. Auditability
Auditable reads/navigation/handoffs preserve actor, Tenant, authoritative source, subject, action, timestamp and correlation context through existing mechanisms. This capability does not create a new audit object. Source-calculated states retain calculation provenance rather than being re-authored by Settings.

## S24. Failure and edge cases
Missing source, missing target version, stale data, conflicting sources, unsupported subject/source, unavailable calculation provenance, permission denial and offline conditions are represented explicitly. No failure path silently upgrades uncertainty to healthy/compliant. External/customer SLO publication remains blocked by `OPEN-019`.

## S25. GWT scenarios
**Given** an authorized user and a Tenant-local source-backed Health observation with freshness, **When** Platform Health is opened, **Then** the source, Tenant, freshness, status and uncertainty are visible and no mutation control is enabled.

**Given** a sourced SLO breach with target version, measurement window, calculation provenance and freshness, **When** the breach is displayed, **Then** it is attributed to its authoritative source and may offer an explicit selected-Tenant human Command handoff without creating Incident, Task, Govern flow or automation.

**Given** an MSSP user with a Security-resolved Authorized Tenant Set, **When** the aggregate Health view is opened, **Then** only authorized Tenants are included read-only and every result preserves Tenant/source/version/freshness without multi-Tenant Report/Export widening.

**Given** a source lacks calculation provenance or freshness required for an authoritative SLO state, **When** Platform Health renders the projection, **Then** it shows an explicit partial/stale/unknown condition and does not synthesize compliance or breach.

**Given** a user has `perm.settings.health.read`, **When** Report, Export, SLO configuration, failover or recovery is considered, **Then** Health read access alone grants none of those effects and only source-owned or Shared handoffs remain available.

## S26. Dependencies and open decisions
Normative dependencies include ADR-0009, ADR-0008, Platform Health, neutral Health/Metrics contracts, Shared Metrics, Business Service context, Permission Model, Tenant isolation, `CAP-SET-001`, `CAP-SET-002`, `CAP-SET-010`, `CAP-SET-012` and `CAP-CMD-401`. `OPEN-008`, `OPEN-013`, `OPEN-015` and `OPEN-019` remain open and unchanged. Source/platform availability, future mutation authority, Automation Run↔Response Run provenance and external publication are implementation/future-architecture concerns, not reasons to invent current behavior.

## S27. Acceptance criteria
- Exactly one Settings capability owns this user goal: `CAP-SET-014`.
- Exactly 27 numbered sections and the six mandatory substantive tables are present.
- `SET-HLT-001` is reused; no new Screen ID is created.
- `perm.settings.health.read` remains read-only; no new Permission ID is created.
- No canonical object is created or modified.
- SLO remains non-canonical and source-attributed; no target store/configuration/calculator is invented.
- Breach authority requires source, target/version, measurement window, calculation provenance and freshness.
- MSSP aggregation remains Authorized-Tenant-Set read-only with Tenant preserved per result.
- Search stays single-selected-Tenant; Report/Export stay Shared-owned, single-Tenant and non-widening.
- No automatic Incident/Task/priority/Govern/Response Run/Automation Run occurs.
- AI remains explanatory/suggestive only with a deterministic/manual path.
- `OPEN-008`, `OPEN-013`, `OPEN-015`, `OPEN-019` remain open; Requirements state remains unchanged.
- Documentary `defined / planned` status does not claim runtime implementation or production availability.
