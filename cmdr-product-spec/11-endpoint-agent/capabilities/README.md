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

---

## EPT-3 — build-time capability layer
The preceding EPT-3 NOT STARTED line is historical pre-EPT-3 evidence. EPT-3 allocates exactly `CAP-EPT-031..046`, all `draft / defined / planned`.

- EPT-3 structure: **16 capabilities / 432 sections / 96 mandatory tables / at least 48 GWT**;
- Endpoint cumulative: **46 capabilities / 1242 sections / 276 mandatory tables**;
- EPT-1 190/190 and EPT-2 200/200 remain preserved; `CAP-EPT-001..030` remain intact;
- local Detection Content consumption/evaluation/match/local signal-candidate does not replace Investigate Detection Engineering or canonical Command Detection/Signal;
- local Endpoint Investigation uses existing context only and creates no Case/Evidence/Finding;
- `OPEN-008` and `OPEN-017` remain open; no platform/source support or final detection runtime/language/model is selected;
- EPT-4..EPT-6 remain **NOT STARTED**;
- Endpoint Screen IDs remain **0**; no Collection, Live Response, containment, API/protocol, physical schema, final RBAC or implementation is introduced.

EPT-3 build-time quality remains pending publication-dependent gates until all 210 gates close.

---

## EPT-4 — build-time capability layer
The preceding EPT-4 NOT STARTED line is historical pre-EPT-4 evidence. EPT-4 allocates exactly `CAP-EPT-047..064`, all `draft / defined / planned`.

- EPT-4 structure: **18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**;
- Endpoint cumulative: **64 capabilities / 1728 sections / 384 mandatory tables**;
- EPT-1 **190/190**, EPT-2 **200/200**, EPT-3 **210/210** remain preserved and `CAP-EPT-001..046` remain intact;
- canonical Investigate Collection Request remains Investigate-owned; Endpoint owns local technical intake/operations/outputs only;
- neutral Collection Item/Collected Technical Output do not become Artifact/Evidence under OPEN-014;
- Endpoint Technical Session/Execution remain distinct from Studio Tool Call/Automation Run and Govern Response Run/Result under OPEN-015;
- OPEN-008/014/015/017 remain open;
- Endpoint Screen IDs remain 0;
- EPT-5/EPT-6 remain NOT STARTED; containment/remediation/update/resilience are not started;
- no API/protocol/transport/final shell/runtime/physical schema/storage/final RBAC/implementation.

EPT-4 historical build-time text above is superseded by its final 220/220 post-publication record.

---

## EPT-5 — build-time capability layer
The preceding EPT-5 NOT STARTED lines are historical pre-EPT-5 snapshots. EPT-5 allocates the **source-driven** set `CAP-EPT-065..081`, all `draft / defined / planned`.

- EPT-5 structure: **17 capabilities / 459 sections / 102 mandatory tables / at least 51 GWT**;
- Endpoint cumulative: **81 capabilities / 2187 sections / 486 mandatory tables**;
- EPT-1 **190/190**, EPT-2 **200/200**, EPT-3 **210/210**, EPT-4 **220/220** remain preserved; `CAP-EPT-001..064` remain intact;
- candidate `065..080` was expanded by `CAP-EPT-081` only because `containment/account-containment.md` independently sources local session lock/termination; directory actions remain external;
- Govern retains Action Request/Approval/Decision/Response Run/response verification/rollback/Result; Endpoint owns only technical primitive/readiness/execution/observation/reversal facts;
- technical verification != Govern verification; technical reversal != Govern rollback; Technical Outcome != Result;
- OPEN-007/008/013/014/015/017 remain open; OPEN total remains 18;
- Endpoint Screen IDs remain 0; EPT-6 remains **NOT STARTED**;
- no API/protocol/native command/PowerShell/shell/firewall syntax/physical schema/final RBAC/final response engine/implementation.

EPT-5 build-time quality is **224 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL** until all 230 gates close remotely.