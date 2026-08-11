---
id: endpoint-agent-readme
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
requirements:
  - REQ-PROD-018
  - REQ-PROD-012
  - REQ-SEC-004
  - REQ-SEC-005
---
# Endpoint Agent

## Mission
Servir de composant local sécurisé pour les capacités endpoint visées : télémétrie, détection, inspection, collecte, Live Response, containment, résilience et audit.

## Possède
- exécution locale selon contrats, policy, permission et autorité ;
- records locaux et résultats d'opérations ;
- EPT-1 local identity/registration, binding, observed platform/version/inventory, health/connectivity/state, capability availability and provenance semantics.

## Consomme
Fleet et policies Settings, workflow Investigate, Decisions/Response Runs Govern, contexte Command.

## Exclusions
- ne possède pas la flotte administrative ;
- ne possède pas Endpoint Policy ou enrollment administration ;
- ne possède pas Incident, Case, Decision ou workflow métier ;
- ne constitue pas à lui seul une preuve de delivery EDR complet.

## Capability Specification
EPT-1 — **Enrollment, Inventory, Health and Platform Foundations** allocates `CAP-EPT-001..014`, all `draft / defined / planned`, with **14 capabilities / 378 sections / 84 mandatory tables**. EPT-1 is an execution lot under Delivery Roadmap Phase 5, not a Roadmap Phase.

Future lots EPT-2 Telemetry, EPT-3 Detection/Investigation, EPT-4 Collection/Live Response, EPT-5 Containment/Verification and EPT-6 Updates/Resilience/Security remain **NOT STARTED**.

See [`capabilities/README.md`](capabilities/README.md), [`capability-map.md`](capability-map.md), [`object-consumption-map.md`](object-consumption-map.md), [`action-classification.md`](action-classification.md), [`automation-and-ai-model.md`](automation-and-ai-model.md), [`cross-product-links.md`](cross-product-links.md) and [`information-architecture.md`](information-architecture.md).

## Transitions principales
Télémétrie → détection ; collecte → Artifact/Evidence ; commande autorisée → résultat vérifiable. These future paths are not started by EPT-1.

## Place de l'IA
Les fonctions essentielles de l'agent ne dépendent pas d'un modèle. Toute inference future reste classifiée, explicable et gouvernée. EPT-1 requires deterministic/manual alternatives.

## Delivery classification
**Target :** EDR natif complet. **Current evidence :** `planned`. EPT-1 defines documentary foundations only; it proves no engine, supported platform or release. `OPEN-008` remains open.

## Sources
- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- [`information-architecture.md`](information-architecture.md)
- [`product-definition.md`](product-definition.md)

## Critère d'acceptation
Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.

---

## EPT-2 build addendum — Telemetry, Observation and Technical Capability Declaration
`CAP-EPT-015..030` add **16 capabilities / 432 sections / 96 mandatory tables** for source/observation semantics, Shared telemetry-event projection, process/file/network/auth/system/sensor observations, normalization, quality/loss/rate/backpressure, privacy, detailed capability declaration/availability, cross-product handoff and provenance.

Endpoint cumulative content becomes **30 capabilities / 810 sections / 180 mandatory tables**. EPT-1 `CAP-EPT-001..014` remain intact. Shared retains `telemetry-event` and generic normalization; Settings retains source/Fleet/Policy administration; Investigate retains Evidence/Finding/Case; Studio retains Tool/Run; Govern retains Decision/Response Run/Result. `OPEN-008` remains open. No Endpoint Screen ID, implementation, API/protocol, physical schema or EPT-3+ work is introduced.

---

## EPT-3 build addendum — Local Detection and Endpoint Investigation
The earlier EPT-3 NOT STARTED wording is preserved as historical pre-EPT-3 evidence. EPT-3 now allocates `CAP-EPT-031..046`, all `draft / defined / planned`, for Detection Content consumption/eligibility, local evaluation/match/local signal-candidate semantics, detection context/grouping/coverage, local process/file/network/user-session/system investigation, local timeline/correlation, pivots, detection-to-investigation expansion, local summary/handoff and provenance.

Endpoint cumulative content becomes **46 capabilities / 1242 sections / 276 mandatory tables**. EPT-1 remains 190/190 PASS and EPT-2 remains 200/200 PASS; `CAP-EPT-001..030` are unchanged. Investigate retains Detection Engineering/Case/Evidence/Finding, Command retains canonical Detection/Signal/Alert/Incident, Shared retains generic Search/Timeline/Linking/Correlation, Settings retains administration and Govern retains response authority. `OPEN-008` and `OPEN-017` remain open. Endpoint Screen IDs remain 0. EPT-4..EPT-6 remain NOT STARTED. No acquisition, Live Response, containment, API/protocol, physical schema, final detection runtime/language/model, final RBAC or implementation is introduced.

---

## EPT-4 build addendum — Collection and Live Response Technical Execution
The historical EPT-4 NOT STARTED wording above is the pre-EPT-4 snapshot. EPT-4 allocates `CAP-EPT-047..064`, all `draft / defined / planned`, covering Collection technical intake/eligibility/authority refs, bounded planning, file/process-system/memory/network acquisition, operation progress, neutral packaging/integrity metadata, transfer/handoff, technical Live Response sessions, command/shell/script/file execution boundaries, technical outputs/errors/operator controls and provenance.

Endpoint cumulative content becomes **64 capabilities / 1728 sections / 384 mandatory tables**. EPT-1/2/3 remain PASS and `CAP-EPT-001..046` are unchanged. Investigate retains Collection Request/Case/Evidence/Finding/Artifact qualification; Govern retains Decision/Response Run/Result/authority; Studio retains Tool Call/Automation Run; Settings retains Fleet/Policy/Secret References; Shared retains Jobs/Trace/Activity/Export. `OPEN-008`, `OPEN-014`, `OPEN-015` and `OPEN-017` remain open. Endpoint Screen IDs remain 0. EPT-5/EPT-6 remain NOT STARTED. No containment, API/protocol/transport/final shell/runtime/physical schema/storage/final RBAC or implementation is introduced.

---

## EPT-5 build addendum — Containment, Verification and Governed Response Primitives
The historical EPT-5 NOT STARTED wording above is the pre-EPT-5 snapshot. EPT-5 allocates `CAP-EPT-065..081`, all `draft / defined / planned`, after complete Containment/effectful/Govern source audit.

Endpoint now has **81 capabilities / 2187 sections / 486 mandatory tables**. EPT-5 itself is **17 / 459 / 102 / at least 51 GWT**. EPT-1/2/3/4 remain PASS and `CAP-EPT-001..064` are unchanged.

EPT-5 covers primitive request/eligibility/authority boundary, readiness/precheck/target state, process/isolation/network/quarantine/delete-restore/service/local-session primitives, normalized technical outcomes, technical verification, partial/unknown/drift, technical reversal, containment release, Govern reconciliation input and provenance. `CAP-EPT-081` is source-driven local-session containment only; directory account actions remain external.

Govern retains Approval/Decision/Response Run/response verification/rollback/Result; Settings retains Fleet/Policy/admin/secrets; Studio retains Human Gate/Tool Call/Automation Run; Shared retains generic mechanisms; Investigate retains Case/Finding/Evidence. OPEN remains 18 and OPEN-007/008/013/014/015/017 stay open. Endpoint Screen IDs remain 0. EPT-6 remains **NOT STARTED**. No implementation, API/protocol/native command/final engine/schema/RBAC is introduced.

---

## EPT-6 build addendum — Updates, Resilience, Security and Endpoint Provenance — 2026-08-12
The earlier EPT-6 `NOT STARTED` wording is preserved above as the exact pre-EPT-6 snapshot. EPT-6 now allocates the source-driven range `CAP-EPT-082..099`, all `draft / defined / planned`.

EPT-6 contributes **18 capabilities / 486 sections / 108 mandatory tables / at least 54 GWT**: update assignment/package/download/readiness/install/activation/progress/retry/health/reversion; offline buffering/replay/restart/crash/resource/dependency resilience; self-protection/anti-tamper, privilege/security context, Secret Reference handling, Local Audit Event, security-state handoff and Endpoint closure.

Endpoint cumulative content becomes **99 capabilities / 2673 sections / 594 mandatory tables**. Global content becomes **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**. Requirements remain conservatively **122 = 99/20/3/0** and OPEN remains **18**.

Settings retains Fleet/update target/channel/wave/Policy and Secret administration; Studio retains Studio Deployment/Reversion; Govern retains response authority/rollback/Result; Shared retains generic Jobs/Retry/Recovery/Trace/Activity; Security retains global permission/privacy/audit-integrity policy. Endpoint Screen IDs remain **0**. No API/protocol/code/package/crypto/physical schema/final RBAC or Phase 6 capability is introduced.

Build-time quality is **233 PASS / 7 PENDING-REMOTE / 0 FAIL**. EPT-6, Endpoint Capability Specification and Delivery Roadmap Phase 5 remain PENDING/PARTIAL until actual post-publication verification closes gates 234–240.