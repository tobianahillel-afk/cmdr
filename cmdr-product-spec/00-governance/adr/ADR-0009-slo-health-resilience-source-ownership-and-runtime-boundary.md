---
id: ADR-0009
domain: 00-governance
status: validated
owner: Product Architecture
updated: 2026-08-16
source-of-truth: canonical
---
# ADR-0009 — SLO / Health / Resilience Source, Ownership and Runtime Boundary

## Context

Delivery Roadmap Phase 6 retains SLO / resilience as a Platform Scale responsibility. The canonical corpus already exposes Platform Health, neutral Health and Metrics implementation contracts, Shared metric mechanisms, source-attributed provider/source health and freshness, Offline/Retry semantics, Business Service context, tenant isolation and Govern response authority. It did not select a generic SLO target store, canonical SLO object, central SLO calculator, generic failover/recovery runtime, or SLO administration permission.

This ADR records the exact human-approved source-of-truth, ownership and runtime boundaries before any future capability allocation. It does not claim runtime implementation.

## Human approval

**Approval reference:** Explicit project-owner approval in this conversation.  
**Approved decision-set checksum:** `adb8312c2eb5cb65062177c72ee3b23cbe9f3c165593dab514a5aa53d6ad674a`

The following authority input is recorded **verbatim**.

## Approved D1 — SLO source of truth

CMDR uses a hybrid source-attributed SLO model. Platform Architecture owns the neutral Health/Metrics contract semantics used to project SLO metadata, but the initial architecture introduces no generic CMDR SLO target store and no generic SLO configuration. Each effective SLO target remains owned by the authoritative source that supplies it and must carry Tenant scope, authoritative source reference, target value/unit, version, effective period and provenance; Environment is included only when sourced. Platform Settings may present the target and sourced state to authorized consumers but does not become the target owner or runtime calculator.

## Approved D2 — SLO identity model

SLO is not a canonical CMDR object in the initial architecture. It is a non-canonical source-attributed configuration/projection referenced by Tenant plus a typed subject reference plus authoritative source reference plus target version/effective period; Environment is included only when the source defines it. No SLO, Health Observation, Threshold, Service Level, Availability, Reliability, Resilience, RTO or RPO canonical object is created by this decision.

## Approved D3 — Measurement / calculation model

Platform Architecture owns the neutral Health and Metrics contract envelope; Shared retains generic metric definition/version, aggregation, reconciliation and freshness mechanisms; authoritative source/runtime owners retain measurement acquisition and any source-specific availability, reliability or SLO calculation. The initial architecture creates no central CMDR SLO calculator. A projected SLO state or breach is authoritative only when source, target version, measurement window, calculation provenance and freshness are available. Platform Settings may perform bounded presentation derivations such as age, stale, unknown or partial from explicit source metadata, but must not synthesize SLO compliance or breach when authoritative calculation evidence is absent.

## Approved D4 — Health boundary

Platform Health is a projection and bounded deterministic presentation layer, not a probe or service-monitoring runtime. Platform Settings owns Health UI/projection and source, freshness and impact presentation. Source/runtime owners retain probes, measurement acquisition, service monitoring and authoritative source Health/SLO computation; Platform Architecture owns the neutral contract. Missing, stale, conflicting or unsupported observations remain explicit and no healthy state is inferred. Settings may derive only presentation facts from explicit source timestamps/rules and never generate measurements or claim runtime ownership.

## Approved D5 — Resilience semantics

Initial Phase-6 Resilience architecture is limited to source-attributed degradation/Health observability, the existing Offline/Retry contract semantics (retry class, backoff, expiry, queue, conflict and user visibility), and existing provider-routing fallback constraints as configuration rather than execution. Generic runtime redundancy, automatic failover, manual failover, data-source fallback, recovery, disaster recovery, RTO and RPO are outside the initial architecture and require separate source-audited architecture/implementation decisions. Retry, fallback, failover, recovery and rollback remain distinct.

## Approved D6 — Configuration / permission model

Initial architecture creates no generic SLO target, threshold, window, breach-policy or monitoring-state administrative mutation and allocates no new permission. perm.settings.health.read authorizes only Health/SLO projection reads within Tenant and source visibility; it does not authorize configuration, export or runtime actions. Local deterministic validation of received metadata is no-effect. Acknowledge maintenance remains non-executable/disabled until an owner, action class and explicit permission are sourced; enable/disable monitoring remains a source/runtime-owner operation. Any future SLO configuration requires a separate source/permission audit and preserves OPEN-013.

## Approved D7 — Failover / recovery authority

CMDR exposes no generic failover or recovery execution in the initial SLO/Health/Resilience architecture. Failover/recovery may be observed only as sourced state. Any future effectful failover/recovery must identify a runtime executor and separate configuration from execution; if modeled as governed response it requires one selected Tenant, Security authorization re-evaluation, Govern Action Request/Decision/Response Run authority, verification, rollback/recovery semantics and Result/audit. Platform administration or Health read access never grants response authority. Automated failover/recovery is excluded initially, so OPEN-015 remains open and is not a definition blocker.

## Approved D8 — Customer / MSSP visibility

Initial Phase-6 MVP SLO/Health visibility is Tenant-local plus MSSP read-only aggregation across the Security-resolved Authorized Tenant Set, preserving Tenant identity, authoritative source, target version and freshness for every projection. Customer remains an external projection. Customer-facing or other external SLO publication is not authorized by this architecture and remains fenced by OPEN-019. Cross-tenant mutation, administration, response and export widening remain prohibited.

## Approved D9 — Search / Reporting / Export

No SLO/Health exception is introduced to ADR-0008. Search remains single-selected-Tenant; Report remains single-Tenant; Export remains single-Tenant and never widens visibility; Shared retains Search/Reporting/Export ownership. An aggregated MSSP Health/SLO view cannot directly create a multi-tenant Search, Report or Export.

## Normative consequences

1. `SLO` is a non-canonical, source-attributed projection/configuration reference in the initial architecture.
2. No canonical `SLO`, `Health Observation`, `Threshold`, `Service Level`, `Availability`, `Reliability`, `Resilience`, `RTO` or `RPO` object is introduced.
3. Platform Architecture owns neutral Health/Metrics contract semantics; Shared retains generic metric definition/version, aggregation, reconciliation and freshness mechanisms.
4. Authoritative source/runtime owners retain measurement acquisition and authoritative source-specific Health/SLO calculations.
5. Platform Settings owns the Platform Health UI/projection and bounded presentation derivations only; it does not become a probe, monitoring, acquisition or universal SLO-calculation runtime.
6. Source-attributed degradation observability, existing Offline/Retry semantics and existing provider fallback constraints are the initial resilience scope. Generic failover, recovery, DR, RTO and RPO remain outside it.
7. `perm.settings.health.read` remains read-only. No generic SLO/Health management permission is created.
8. `Acknowledge maintenance` is non-executable/disabled until an owner, action class and explicit permission are separately sourced.
9. Generic failover/recovery execution is absent. Future governed effectful execution must preserve selected-Tenant context, Security re-evaluation and Govern authority.
10. Tenant-local plus Authorized-Tenant-Set MSSP read-only aggregation is allowed for SLO/Health projections; cross-tenant mutation/admin/response/export widening is prohibited.
11. External/customer SLO publication remains outside this architecture and fenced by `OPEN-019`.
12. Search remains single-selected-Tenant; Report and Export remain single-Tenant; Shared retains Search/Reporting/Export ownership.
13. A sourced SLO breach may be projected and notified, but it does not automatically create a Command Incident/Task, change priority, trigger Govern or trigger automation.
14. Deterministic/manual paths remain mandatory; AI may explain sourced state but may not fabricate measurements/compliance/breaches or execute response.

## Open decisions preserved

- `OPEN-006` remains resolved by ADR-0008.
- `OPEN-008` remains open for platform/source availability and support.
- `OPEN-013` remains open for default governance/authority of reversible class-2 mutations.
- `OPEN-015` remains open for the Automation Run / Response Run provenance bridge.
- `OPEN-019` remains open for dissemination/releasability/sharing/access.

No other OPEN is closed or weakened by this ADR.

## Explicit non-goals

This ADR creates or allocates:
- **0 new capability IDs**;
- **0 `CAP-SET-014` allocation or reservation**;
- **0 new canonical objects**;
- **0 new Permission IDs**;
- **0 new Screen IDs**.

It does not implement probes, monitoring, measurement acquisition, SLO calculation, provider/source runtime changes, failover, recovery, DR, Search, Reporting, Export, Customer publication, product code, API/protocol or physical storage schema.

## Follow-up

After this architecture-recording run is closed and remotely verified, work must stop. A **new source-audited capability preparation** against the resulting FINAL HEAD is required before any functional SLO capability, capability namespace allocation, permission addition, screen addition or runtime implementation.
