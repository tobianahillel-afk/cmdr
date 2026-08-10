---
id: roadmap-phase-5-studio-and-endpoint
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-10
source-of-truth: canonical
---
# Phase 5 Studio And Endpoint

## Objectif
Définir Phase 5 Studio and Endpoint pour CMDR tout en conservant Studio et Endpoint comme deux capability domains séparés.

## Périmètre
Document canonique Roadmap. Il ne remplace pas les sources propriétaires et ne transforme aucun execution lot en sous-phase numérotée.

## Propriétaire fonctionnel
Product Operations Lead.

## Objets concernés
- Concepts du document
- Références canoniques liées

## Fonctionnalités historiques prévues
- Skills/Agents/workflows.
- Assurance/control room.
- EDR telemetry/detection/response.
- Fleet administration comme dépendance Settings, sans transfert d'ownership vers Endpoint.

## UX et interactions
- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions
Les modifications suivent `../14-security-permissions-and-trust/permission-model.md`; aucun execution lot ne finalise implicitement RBAC/ABAC.

## Dépendances
- `../00-governance/source-of-truth-policy.md`
- `dependency-roadmap.md`
- Studio/Settings/Govern/Shared/Endpoint ownership sources.

## STD-1 — verified
- STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts: **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**.
- `CAP-STD-001..016`: 16 capabilities / 432 sections / 96 mandatory tables.

## STD-2 — current execution
- STD-2 — Workflow Builder & Orchestration: 17 capabilities `CAP-STD-017..033` / 459 sections / 102 mandatory tables.
- scope: Workflow definition/version, Builder, data/graph/composition, deterministic control, errors/retry/partial/compensation, Human Gate boundary, readiness/pre-publish/provenance.
- status: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**; fifth functional SHA `655e9ce0ade2d64a7738a6a479572fef9b6f0e2f` verified at 5 ahead / 0 behind from baseline.
- new Screen IDs / detailed rewrites: 0 / 0.
- runtime scheduler / Automation Run lifecycle / publishing-deployment / Endpoint implementation: not included.

## Future lots
- STD-3 — Agents, Human Gates & Runtime Control: NOT STARTED.
- STD-4 — Assurance & Lifecycle: NOT STARTED.
- Endpoint capability specification: NOT STARTED.
- Studio capability specification: PARTIAL.
- Delivery Roadmap Phase 5 overall: PARTIAL.

STD-1/2/3/4 are execution-lot labels only. No Phase 5A/5B/5C/5D exists.

## Implementation boundary
No API, protocol, product code, orchestration language/runtime, final graph/Tool/Tool Call physical schema, final RBAC/ABAC, provider/runtime selection, detailed screen rewrite or Endpoint implementation is introduced by STD-2.

## Questions ouvertes
OPEN-007/013/015 remain open where consumed. OPEN-008 remains Endpoint/platform support. The permission namespace anomaly remains unresolved.

## Next candidate
STD-3 — Agents, Human Gates & Runtime Control. **Do not begin STD-3 implicitly.**

## STD-2 post-publication verification evidence
Baseline `04dcdb43fd7f944a700bf936eebef003546095eb`; fifth functional/build SHA `655e9ce0ade2d64a7738a6a479572fef9b6f0e2f`; five functional commits reachable in order; PR #2 remains Draft/open/unmerged; branch/main README unchanged; CI N/A; `CAP-STD-001..016` intact; Endpoint 0; STD-3/4 not started.

Final STD-2 documentary verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**.

## STD-3 current execution addendum
The full STD-1/STD-2 roadmap text above is preserved as the exact pre-STD-3 snapshot; its STD-3 `NOT STARTED` / `Next candidate` text is historical evidence only.

- execution lot: **STD-3 — Agents, Human Gates & Runtime Control**;
- `CAP-STD-034..051`: **18 capabilities / 486 sections / 108 mandatory tables**;
- fifth functional/build SHA: `c658168c9de6bd803941116989bc3aaedf154260`;
- baseline `c472b055ce00fd33efd96ac920b0add5f65ab8f7` → build: **5 ahead / 0 behind**, same merge base;
- scope: Agents/Teams, bounded planning/oversight, Human Gate runtime, Automation Runs/control, queue/scheduling semantics, failures/retries/partial, transient context, Control Room, outcome/handoff and provenance;
- ownership boundaries: Govern retains Approval/Decision/Response Run/Result; Settings retains provider/secret/runtime administration; Shared retains generic Jobs/queue/scheduling/Trace; Endpoint retains technical primitives;
- no new Screen ID/detailed rewrite, agent framework/model/provider/runtime/scheduler/API/protocol/code/final JSON Schema/RBAC;
- post-publication history-preservation correction is being recorded before the final 210-gate verdict.

STD-3 remains **PENDING POST-PUBLICATION VERIFICATION**. STD-4 and Endpoint remain **NOT STARTED**; Delivery Roadmap Phase 5 and global/repository maturity remain PARTIAL.