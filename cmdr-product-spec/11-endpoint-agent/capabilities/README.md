---
id: endpoint-capabilities-readme
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint Agent Capabilities

Endpoint Capability Specification is executed as lots under Delivery Roadmap Phase 5 — Studio and Endpoint. EPT-1 is the first lot and is not a Roadmap Phase.

## EPT-1 — Enrollment, Inventory, Health and Platform Foundations
`CAP-EPT-001..014` define provider/platform-neutral foundations only: identity/registration, local enrollment, tenant/environment binding, observed platform/OS/architecture, Agent version/build/compatibility context, inventory/freshness, health/self-check, heartbeat/connectivity/last-seen, operational state, technical capability availability, Fleet/Policy boundaries and provenance.

Every capability is `draft / defined / planned`, uses 27 sections and the six mandatory capability tables. No capability claims implementation or supported-platform delivery.

## Future lots — NOT STARTED
- EPT-2 — Telemetry, Observation and Technical Capability Declaration.
- EPT-3 — Local Detection and Endpoint Investigation.
- EPT-4 — Collection and Live Response Technical Execution.
- EPT-5 — Containment, Verification and Governed Response Primitives.
- EPT-6 — Updates, Resilience, Security and Endpoint Provenance.

## Mandatory boundaries
Fleet, enrollment administration and Endpoint Policy remain Platform Settings-owned; Endpoint Agent owns individual local technical state. Endpoint Agent != Studio Automation Agent; technical capability != Studio Tool; technical output != Govern Result; Endpoint execution != Response Run/Tool Call. OPEN-008 remains open.

---

## EPT-2 — build-time capability layer
The preceding EPT-2 NOT STARTED line is preserved as the pre-EPT-2 snapshot. Current EPT-2 content allocates exactly `CAP-EPT-015..030`, all `draft / defined / planned`.

- EPT-2 structure: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**;
- Endpoint cumulative: **30 capabilities / 810 sections / 180 mandatory tables**;
- EPT-1 `CAP-EPT-001..014` remain intact;
- `CAP-EPT-011` remains the foundation availability summary; `CAP-EPT-027/028` provide detailed declaration/dynamic availability without duplication;
- Shared retains `telemetry-event` and generic normalization;
- `OPEN-008` remains open; no platform/source support is declared delivered;
- EPT-3..EPT-6 remain **NOT STARTED**;
- Endpoint Screen IDs remain **0** and no implementation/API/protocol/physical schema/final RBAC is introduced.

EPT-2 remains pending final post-publication verification until all 200 gates close.
